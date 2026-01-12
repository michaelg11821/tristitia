package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"time"

	http "github.com/bogdanfinn/fhttp"
)

func (e *engine) makeReq(method, url string, headers [][2]string, body []byte) (*http.Request, error) {
	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("Failed to create request: %v", err)
	}

	for _, header := range headers {
		req.Header.Add(header[0], header[1])
	}

	return req, nil
}

func (e *engine) doReq(req *http.Request, useDefaultResponseHandling bool) (*http.Response, error) {
	maxRetries := 3

	for i := 0; i < maxRetries; i++ {
		resp, err := e.Client.Do(req)
		if err != nil {
			return nil, fmt.Errorf("Failed to perform request: %v", err)
		}

		if useDefaultResponseHandling {
			if resp.StatusCode >= 400 {
				discardResp(resp)
				time.Sleep(time.Second * time.Duration(e.Cooldown))
				continue
			}
		}
		return resp, nil
	}

	return nil, fmt.Errorf("Failed after %d retries.", maxRetries)
}

func discardResp(resp *http.Response) {
	if resp != nil && resp.Body != nil {
		io.Copy(io.Discard, resp.Body)
		defer resp.Body.Close()
	}
}

func readJsonRespBody[R []NarouChapterNumResp | []NarouCharCountResp](res *http.Response) (R, error) {
	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("Could not read response: %v", err)
	}

	var resp R
	if err := json.Unmarshal(resBody, &resp); err != nil {
		return nil, fmt.Errorf("Could not decode JSON: %v", err)
	}

	return resp, nil
}
