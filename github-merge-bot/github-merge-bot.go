package main

import (
	"./slack"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
	"os"
)

func mergePullRequest(owner, repo string, number int) (string, int, error) {
	ctx := context.Background()
	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_ACCESS_TOKEN")},
	)
	client := github.NewClient(oauth2.NewClient(ctx, src))

	result, resp, err := client.PullRequests.Merge(
		context.Background(), owner, repo, number, "", nil,
	)
	if err != nil {
		return *result.Message, resp.StatusCode, err
	}
	return *result.Message, resp.StatusCode, nil
}

func HandleRequest(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	prEvent := new(github.PullRequestEvent)
	if err := json.Unmarshal([]byte(request.Body), &prEvent); err != nil {
		return events.APIGatewayProxyResponse{}, errors.New("Failed to unmarshal JSON")
	}

	if prEvent.GetAction() != "review_requested" {
		return events.APIGatewayProxyResponse{StatusCode: 400}, errors.New("action must be review_requested")
	}

	owner := *prEvent.GetRepo().Owner.Login
	repo := *prEvent.GetRepo().Name
	number := prEvent.GetNumber()
	resultMessage, statusCode, err := mergePullRequest(owner, repo, number)
	slack.Notify(fmt.Sprint(resultMessage))

	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: statusCode}, err
	}
	return events.APIGatewayProxyResponse{StatusCode: 200}, nil
}

func main() {
	lambda.Start(HandleRequest)
}
