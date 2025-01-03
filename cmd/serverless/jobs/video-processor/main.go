package main

import (
	"context"

	"github.com/ESSantana/streaming-test/cmd/serverless/jobs/video-processor/handler"
	"github.com/ESSantana/streaming-test/internal/services"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var (
	videoProcessorHandler *handler.VideoProcessorHandler
)

func init() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	session, err := session.NewSession()
	if err != nil {
		panic(err)
	}
	s3Client := s3.New(session, aws.NewConfig().WithRegion("sa-east-1"))
	videoService := services.NewVideoService(s3Client)
	
	videoProcessorHandler = handler.NewVideoProcessorHandler(videoService)
}

func main() {
	lambda.Start(Handler)
}

func Handler(ctx context.Context, event events.S3Event) error {
	for _, record := range event.Records {
		err := videoProcessorHandler.ProcessVideo(ctx, record.S3.Object.Key)
		if err != nil { 
			log.Info().Msgf("Error processing video: %s\n", err.Error())
            continue
		}
		log.Info().Msgf("Video processed: %s\n", record.S3.Object.Key)
	}

	return nil
}
