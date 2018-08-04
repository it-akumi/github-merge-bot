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

func Notify(text string) {
	body, err := json.Marshal(Message{Text: text})
	if err != nil {
		println("Failed to marshal message")
		os.Exit(1)
	}
	req, err := http.NewRequest(
		"POST",
		os.Getenv("SLACK_INCOMING_WEBHOOK_URL"),
		bytes.NewBuffer(body),
	)
	if err != nil {
		println("Failed to create new request")
		os.Exit(1)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		println("Failed to notify to slack")
		os.Exit(1)
	}
	defer resp.Body.Close()
	return
}
