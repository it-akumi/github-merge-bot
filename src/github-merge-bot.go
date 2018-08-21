package main

import (
	"./slack"
	"context"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
	"os"
)

func mergePullRequest(pullRequestEvent *github.PullRequestEvent) (string, int, error) {
	ctx := context.Background()
	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_ACCESS_TOKEN")},
	)
	client := github.NewClient(oauth2.NewClient(ctx, src))

	owner := pullRequestEvent.GetRepo().GetOwner().GetLogin()
	repo := pullRequestEvent.GetRepo().GetName()
	number := pullRequestEvent.GetNumber()
	result, resp, err := client.PullRequests.Merge(
		context.Background(), owner, repo, number, "", nil,
	)
	if err != nil {
		var resultMessage string
		if resp.StatusCode == 405 {
			resultMessage = "Pull Request is not mergeable\n"
		} else if resp.StatusCode == 409 {
			resultMessage = "Head branch was modified. Review and try the merge again.\n"
		}
		pullRequestHTMLURL := pullRequestEvent.GetPullRequest().GetHTMLURL()
		return resultMessage + pullRequestHTMLURL, resp.StatusCode, err
	}

	// Delete merged branch
	ref := "heads/" + pullRequestEvent.GetPullRequest().GetHead().GetRef()
	resp, err = client.Git.DeleteRef(context.Background(), owner, repo, ref)
	if err != nil {
		return "Failed to delete merged branch", resp.StatusCode, err
	}

	return result.GetMessage(), resp.StatusCode, nil
}

func HandleRequest(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	pullRequestEvent := new(github.PullRequestEvent)
	if err := json.Unmarshal([]byte(request.Body), &pullRequestEvent); err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       "Failed to unmarshal JSON\nRequest body may be invalid",
		}, nil
	}

	if pullRequestEvent.GetAction() != "review_requested" {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       "Request parameter 'action' must be 'review_requested'",
		}, nil
	}

	resultMessage, statusCode, err := mergePullRequest(pullRequestEvent)
	slack.Notify(resultMessage)

	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: statusCode, Body: resultMessage}, nil
	}
	return events.APIGatewayProxyResponse{StatusCode: 200, Body: resultMessage}, nil
}

func main() {
	lambda.Start(HandleRequest)
}
