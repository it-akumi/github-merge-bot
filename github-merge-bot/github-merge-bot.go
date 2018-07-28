package main

import (
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
)

func GitHubMergeHandler() {
	fmt.Println("Hello, github-merge-bot!")
}

func main() {
	lambda.Start(GitHubMergeHandler)
}
