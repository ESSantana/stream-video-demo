package repositories

import (
	"context"

	irepository "github.com/ESSantana/streaming-test/internal/repositories/interfaces"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type repositoryManager struct {
	conn *dynamodb.Client
}

func NewRepositoryManager() (manager irepository.RepositoryManager, err error) {
	conn, err := connectDynamodb()
	if err != nil {
		return nil, err
	}

	return &repositoryManager{
		conn: conn,
	}, nil
}

func connectDynamodb() (conn *dynamodb.Client, err error) {
	cfg, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithRegion("sa-east-1"),
		config.WithRetryer(func() aws.Retryer {
			return aws.NopRetryer{}
		}),
	)

	if err != nil {
		return nil, err
	}

	return dynamodb.NewFromConfig(cfg), nil
}

func (manager *repositoryManager) NewVideoRepository() irepository.VideoRepository {
	return newVideoRepository(manager.conn)
}
