package slack

import (
	"encoding/json"
	"net/http"
	"net/url"
	"os"
)

type Message struct {
	Text string `json:"text"`
}

func Notify(msg Message) {
	req, err := http.NewRequest{
		"POST",
		os.Getenv("SLACK_INCOMING_WEBHOOK_URL"),
		json.Mershal(msg),
	}
	if err != nil {
		println("Failed to create new request...")
		os.Exit(1)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := &http.client{}.Do(req)
	if err != nil {
		println("Failed to notify to slack...")
		os.Exit(1)
	}
	defer resp.Body.Close()
	return
}
