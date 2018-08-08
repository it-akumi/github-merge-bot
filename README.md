github-merge-bot
====
Try to merge pull request automatically and notify its result to slack when a review is requested

## Usage
After a deployment, add Webhooks to your repositories and set the Payload URL to endpoint of your API Gateway.

Now this bot is ready.
When you add reviewers to your pull request, it will be automatically merged.
Then the result will be notified to your slack channel.

## How to deploy
You can deploy github-merge-bot as follows.

### Set these properties in template.yml

* Role
  * An arn of role who execute the Lambda function.

* Environment Variables
  * GITHUB\_ACCESS\_TOKEN
    * A pull request is merged by a user whose GitHub Access Token is set.
  * SLACK\_INCOMING\_WEBHOOK\_URL
    * A notification is posted to this URL.

### Build

```
$ go build -o github-merge-bot/github-merge-bot github-merge-bot/github-merge-bot.go
```

### Packaging and Deployment

```
$ sam package \
     --template-file ./template.yml \
     --s3-bucket your-s3-bucket-name \
     --output-template-file your-output-template.yml
$ sam deploy \
     --template-file your-output-template.yml \
     --stack-name your-stack-name
```

See [AWS document](https://docs.aws.amazon.com/lambda/latest/dg/serverless-deploy-wt.html#serverless-deploy) for more details.

## Author
[Takumi Ishii](https://github.com/it-akumi)

## License
[MIT](https://github.com/it-akumi/github-merge-bot/blob/master/LICENSE)
