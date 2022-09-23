package main

import (
	"context"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	cognito "github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	log "github.com/sirupsen/logrus"
)

var (
	cognitoClient *cognito.CognitoIdentityProvider

	// environment variables
	userPoolID = os.Getenv("COGNITO_USER_POOL_ID")
)

func init() {
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String("ap-south-1"),
	}))

	cognitoClient = cognito.New(sess)

	log.SetFormatter(&log.TextFormatter{
		DisableColors:    true,
		DisableTimestamp: true,
	})

	log.SetLevel(log.InfoLevel)
}

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return router(ctx, request), nil
}

func main() {
	lambda.Start(handler)
}
