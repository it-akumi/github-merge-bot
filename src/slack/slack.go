package slack

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
)

type Message struct {
	Text string `json:"text"`
}

func Notify(text string) error {
	body, err := json.Marshal(Message{Text: text})
	if err != nil {
		return err
	}
	req, err := http.NewRequest(
		"POST",
		os.Getenv("SLACK_INCOMING_WEBHOOK_URL"),
		bytes.NewBuffer(body),
	)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}
