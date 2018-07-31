package main

import (
	"encoding/json"
	"errors"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/google/go-github/github"
)

func HandleRequest(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	pr := new(github.PullRequest)
	if err := json.Unmarshal([]byte(request.Body), &pr); err != nil {
		return events.APIGatewayProxyResponse{}, errors.New("Failed to unmarshal JSON")
	}

	// Try to merge PullRequest

	// Notify result to slack

	return events.APIGatewayProxyResponse{StatusCode: 200, Body: "OK"}, nil
}

func main() {
	lambda.Start(HandleRequest)
}
