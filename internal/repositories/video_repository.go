package repositories

import (
	"context"
	"os"
	"slices"

	"github.com/ESSantana/streaming-test/internal/domain/models"
	"github.com/ESSantana/streaming-test/internal/repositories/interfaces"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type videoRepository struct {
	conn      *dynamodb.Client
	tableName string
}

func newVideoRepository(conn *dynamodb.Client) interfaces.VideoRepository {
	stage := os.Getenv("stage")
	return &videoRepository{
		conn:      conn,
		tableName: "video-stream-demo-" + stage,
	}
}

func (repo *videoRepository) SaveBatch(ctx context.Context, videos []models.Video) (err error) {

	var writeRequestBatches [][]types.WriteRequest

	for chunk := range slices.Chunk(videos, 25) {
		var writeRequest []types.WriteRequest
		for _, video := range chunk {
			item, err := attributevalue.MarshalMap(video)
			if err != nil {
				return nil
			}
			writeRequest = append(writeRequest, types.WriteRequest{
				PutRequest: &types.PutRequest{
					Item: item,
				},
			})
		}
		writeRequestBatches = append(writeRequestBatches, writeRequest)
	}

	for _, batch := range writeRequestBatches {
		output, err := repo.conn.BatchWriteItem(ctx, &dynamodb.BatchWriteItemInput{
			RequestItems: map[string][]types.WriteRequest{
				repo.tableName: batch,
			},
		})
		if err != nil {
			return err
		}

		if len(output.UnprocessedItems[repo.tableName]) > 0 {
			err := repo.reprocessItem(ctx, output.UnprocessedItems[repo.tableName])
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (repo *videoRepository) reprocessItem(ctx context.Context, items []types.WriteRequest) (err error) {
	for _, item := range items {
		putRequest := &dynamodb.PutItemInput{
			TableName: &repo.tableName,
			Item:      item.PutRequest.Item,
		}

		_, err := repo.conn.PutItem(ctx, putRequest)
		if err != nil {
			return err
		}
	}
	return nil
}

func (repo *videoRepository) ListAvailableVideos(ctx context.Context) (videos []models.Video, err error) {
	scanRequest := &dynamodb.ScanInput{
		TableName: &repo.tableName,
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
