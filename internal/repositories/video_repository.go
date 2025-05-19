package repositories

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/ESSantana/streaming-test/internal/domain/models"
	"github.com/ESSantana/streaming-test/internal/repositories/interfaces"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type videoRepository struct {
	conn      *dynamodb.Client
	tableName string
}

func newVideoRepository(conn *dynamodb.Client) interfaces.VideoRepository {
	stage := os.Getenv("STAGE")
	return &videoRepository{
		conn:      conn,
		tableName: fmt.Sprintf("video-stream-demo-%s", stage),
	}
}

func (repo *videoRepository) Save(ctx context.Context, video models.Video) (err error) {
	start := time.Now()
	fmt.Printf("save data to %s\n", repo.tableName)
	fmt.Printf("connection established: %v - data: %+v\n", repo.conn != nil, video)

	item, err := attributevalue.MarshalMap(video)
	if err != nil {
		return err
	}
	putRequest := &dynamodb.PutItemInput{
		TableName: aws.String(repo.tableName),
		Item:      item,
	}

	_, err = repo.conn.PutItem(ctx, putRequest)
	fmt.Printf("time elapssed: %v\n", time.Since(start).Milliseconds())
	return err
}

func (repo *videoRepository) ListAvailableVideos(ctx context.Context) (videos []models.Video, err error) {
	scanRequest := &dynamodb.ScanInput{
		TableName: aws.String(repo.tableName),
	}

	for {
		output, err := repo.conn.Scan(ctx, scanRequest)
		if err != nil {
			return videos, err
		}
		var partial []models.Video
		err = attributevalue.UnmarshalListOfMaps(output.Items, &partial)
		if err != nil {
			return videos, err
		}

		videos = append(videos, partial...)

		if output.LastEvaluatedKey != nil {
			scanRequest.ExclusiveStartKey = output.LastEvaluatedKey
			continue
		}

		break
	}

	return videos, nil
}
