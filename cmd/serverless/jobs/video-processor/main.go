package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

var (
	s3Client *s3.S3
)

func init() {
	session, err := session.NewSession()
	if err != nil {
		panic(err)
	}
	s3Client = s3.New(session, aws.NewConfig().WithRegion("sa-east-1"))

}

func main() {
	lambda.Start(Handler)
}

func Handler(ctx context.Context, event events.S3Event) error {
	out, err := json.Marshal(event)
	if err != nil {
		fmt.Printf("Error marshalling event: %s\n", err.Error())
		return nil
	}

	fmt.Printf("Processing event: %s\n", out)

	return nil
}
