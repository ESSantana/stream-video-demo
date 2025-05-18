package models

type Video struct {
	VideoId   string `json:"video_id" dynamodbav:"video_id"`
	VideoName string `json:"video_name" dynamodbav:"video_name"`
	Manifest  string `json:"manifest" dynamodbav:"manifest"`
	Thumbnail string `json:"thumbnail" dynamodbav:"thumbnail"`
	CreatedAt uint64 `json:"-" dynamodbav:"created_at"`
}
