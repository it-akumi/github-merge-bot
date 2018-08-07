package main

import (
	"./slack"
	"context"
	"encoding/json"
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
		var resultMessage string
		if resp.StatusCode == 405 {
			resultMessage = "Pull Request is not mergeable"
		} else if resp.StatusCode == 409 {
			resultMessage = "Head branch was modified. Review and try the merge again."
		}
		return resultMessage, resp.StatusCode, err
	}
	return *result.Message, resp.StatusCode, nil
}

func HandleRequest(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	prEvent := new(github.PullRequestEvent)
	if err := json.Unmarshal([]byte(request.Body), &prEvent); err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       "Failed to unmarshal JSON\nRequest body may be invalid",
		}, nil
	}

	if prEvent.GetAction() != "review_requested" {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       "Request parameter 'action' must be 'review_requested'",
		}, nil
	}

	owner := *prEvent.GetRepo().Owner.Login
	repo := *prEvent.GetRepo().Name
	number := prEvent.GetNumber()
	resultMessage, statusCode, err := mergePullRequest(owner, repo, number)
	slack.Notify(fmt.Sprint(resultMessage))

	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: statusCode, Body: resultMessage}, nil
	}
	return events.APIGatewayProxyResponse{StatusCode: 200, Body: resultMessage}, nil
}

func main() {
	lambda.Start(HandleRequest)
}
