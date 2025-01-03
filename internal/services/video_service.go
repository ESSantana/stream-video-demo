package services

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/ESSantana/streaming-test/internal/services/interfaces"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"

	ffmpeg "github.com/u2takey/ffmpeg-go"
)

type VideoService struct {
	s3Client *s3.S3
}

func NewVideoService(s3Client *s3.S3) interfaces.VideoService {
	return &VideoService{
		s3Client: s3Client,
	}
}

func (s *VideoService) CreateS3PresignedPutURL(ctx context.Context, bucket, filename, contentType string) (presignedURL string, err error) {
	objectKey := "raw/" + filename

	req, _ := s.s3Client.PutObjectRequest(
		&s3.PutObjectInput{
			Bucket:      aws.String(bucket),
			Key:         aws.String(objectKey),
			ContentType: aws.String(contentType),
		},
	)

	presignedURL, err = req.Presign(time.Minute * 15)
	if err != nil {
		return "", err
	}

	return presignedURL, nil
}

func (s *VideoService) ProcessVideoWithOptions(ctx context.Context, bucket, videoKey string, options interface{}) (err error) {
	videoObject, err := s.s3Client.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(videoKey),
	})
	if err != nil {
		return errors.New("read from s3: " + err.Error())
	}
	defer videoObject.Body.Close()

	videoData, err := io.ReadAll(videoObject.Body)
	if err != nil {
		return errors.New("read object data: " + err.Error())
	}

	tempFilePath := os.TempDir() + "/" + videoKey
	f, err := os.Create(tempFilePath)
	if err != nil {
		return errors.New("create tmp file: " + err.Error())
	}

	_, err = f.Write(videoData)
	if err != nil {
		return errors.New("write to tmp file: " + err.Error())
	}

	videoKeyParts := strings.Split(videoKey, "/")
	videoName := strings.ReplaceAll(videoKeyParts[len(videoKeyParts)-1], ".mp4", "")

	manifestFilePath := os.TempDir() + "/processed/" + videoName + "/index.m3u8"
	segmentFilePath := os.TempDir() + "/processed/" + videoName + "/segment%03d.ts"

	_ = ffmpeg.Input(tempFilePath).Output(manifestFilePath, ffmpeg.KwArgs{
		"vcodec":               "libx264",
		"acodec":               "acc",
		"codec":                "copy",
		"start_number":         0,
		"hls_time":             10,
		"hls_playlist_type":    "vod",
		"hls_segment_filename": segmentFilePath,
		"hls_list_size":        0,
	}).ErrorToStdOut().Run()

	entries, err := os.ReadDir(os.TempDir() + "/processed/" + videoName)
	if err != nil {
		return errors.New("list processed files: " + err.Error())
	}

	for _, entry := range entries {
		data, err := os.OpenFile(entry.Name(), os.O_RDWR, 0666)
		if err != nil {
			return errors.New("open entry file: " + err.Error())
		}

		_, err = s.s3Client.PutObject(&s3.PutObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String("processed/" + videoName + "/" + entry.Name()),
			Body:   data,
		})
		if err != nil {
			return errors.New("save processed to s3: " + err.Error())
		}

		err = os.Remove(entry.Name())
		if err != nil {
			fmt.Println(err.Error())
		}
	}

	return nil
}
