package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type DiscordWebhookPayload struct {
	Content string `json:"content"`
}

func postDiscordWebhook(webhookURL string, content string) error {
	if webhookURL == "" {
		return nil
	}

	body, err := json.Marshal(DiscordWebhookPayload{Content: content})
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, webhookURL, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("content-type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode >= 300 {
		b, _ := io.ReadAll(res.Body)
		if len(b) > 0 {
			return fmt.Errorf("discord webhook status=%s body=%s", res.Status, string(b))
		}
		return fmt.Errorf("discord webhook status=%s", res.Status)
	}

	return nil
}
