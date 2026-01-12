package main

import (
	"fmt"
	"time"

	http "github.com/bogdanfinn/fhttp"
)

func (e *engine) getCharCount(ncode string, headers [][2]string) (int64, error) {
	req, err := e.makeReq(http.MethodGet, getCharCountEndpoint(ncode), headers, nil)
	if err != nil {
		return 0, fmt.Errorf("[%s] Error creating request to Narou: %v", ncode, err)
	}

	res, err := e.doReq(req, true)
	if err != nil {
		return 0, fmt.Errorf("[%s] Error performing request to Narou: %v", ncode, err)
	}
	defer res.Body.Close()

	if res.StatusCode >= 400 {
		return 0, fmt.Errorf("[%s] Error from Narou: status=%s", ncode, res.Status)
	}

	resp, err := readJsonRespBody[[]NarouCharCountResp](res)
	if err != nil {
		return 0, fmt.Errorf("[%s] Error reading JSON body from Narou: %v", ncode, err)
	}

	if len(resp) < 2 {
		return 0, fmt.Errorf("[%s] Unexpected Narou JSON shape (expected 2 array items)", ncode)
	}

	return resp[1].CharCount, nil
}

func (e *engine) getLatestChapterNum(ncode string, headers [][2]string) (int, error) {
	req, err := e.makeReq(http.MethodGet, getLatestChapterNumEndpoint(ncode), headers, nil)
	if err != nil {
		return 0, fmt.Errorf("[%s] Error creating request to Narou: %v", ncode, err)
	}

	res, err := e.doReq(req, true)
	if err != nil {
		return 0, fmt.Errorf("[%s] Error performing request to Narou: %v", ncode, err)
	}
	defer res.Body.Close()

	if res.StatusCode >= 400 {
		return 0, fmt.Errorf("[%s] Error from Narou: status=%s", ncode, res.Status)
	}

	resp, err := readJsonRespBody[[]NarouChapterNumResp](res)
	if err != nil {
		return 0, fmt.Errorf("[%s] Error reading JSON body from Narou: %v", ncode, err)
	}

	if len(resp) < 2 {
		return 0, fmt.Errorf("[%s] Unexpected Narou JSON shape (expected 2 array items)", ncode)
	}

	return resp[1].LatestChapter, nil
}

func (e *engine) sendWatchNotification(ncode string, currentLength, lastLength int64, headers [][2]string) error {
	if err := postDiscordWebhook(e.Webhook, fmt.Sprintf("[%s] A new chapter has been uploaded to Narou!", ncode)); err != nil {
		return fmt.Errorf("[%s] Error posting Discord webhook: %v", ncode, err)
	}

	latestChapterNum, err := e.getLatestChapterNum(ncode, headers)
	if err != nil {
		return err
	}

	latestChapterLink := fmt.Sprintf("https://ncode.syosetu.com/n2267be/%d/", latestChapterNum)

	if err := postDiscordWebhook(e.Webhook, fmt.Sprintf("[%s] The chapter will be available (usually at 1AM JST) at %s", ncode, latestChapterLink)); err != nil {
		return fmt.Errorf("[%s] Error posting Discord webhook: %v", ncode, err)
	}
	if err := postDiscordWebhook(e.Webhook, fmt.Sprintf("[%s] Getting character count of the new chapter...", ncode)); err != nil {
		return fmt.Errorf("[%s] Error posting Discord webhook: %v", ncode, err)
	}

	newCharCount := currentLength - lastLength
	if err := postDiscordWebhook(e.Webhook, fmt.Sprintf("[%s] New chapter character count: %d", ncode, newCharCount)); err != nil {
		return fmt.Errorf("[%s] Error posting Discord webhook: %v", ncode, err)
	}

	return nil
}

func (e *engine) watch(ncode string) {
	var lastLength int64
	var hasLastLength bool

	defaultHeaders := [][2]string{
		{"accept-language", "en-US,en;q=0.9"},
		{"user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/137.0.0.0 Safari/537.36"},
		{"content-type", "application/json"},
	}

	pingNarou := func() {
		currentLength, err := e.getCharCount(ncode, defaultHeaders)
		if err != nil {
			fmt.Println(err)
			return
		}

		if !hasLastLength {
			lastLength = currentLength
			hasLastLength = true
		} else if currentLength != lastLength {
			err := e.sendWatchNotification(ncode, currentLength, lastLength, defaultHeaders)
			if err != nil {
				fmt.Println(err)
			}

			lastLength = currentLength
		}

		fmt.Printf("[%s] Current character count: %d\n", ncode, currentLength)
	}

	pingNarou()

	ticker := time.NewTicker(time.Duration(e.Cooldown) * time.Second)
	defer ticker.Stop()
	for range ticker.C {
		pingNarou()
	}
}
