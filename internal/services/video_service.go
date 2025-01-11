package services

import (
	"context"
	"io"
	"os"
	"strings"
	"time"

	"github.com/ESSantana/streaming-test/internal/domain"
	"github.com/ESSantana/streaming-test/internal/services/interfaces"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/rs/zerolog/log"

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

func (s *VideoService) ProcessVideoWithOptions(ctx context.Context, bucket, videoKey string, options domain.VideoOptions) (err error) {
	videoObject, err := s.s3Client.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(videoKey),
	})
	if err != nil {
		return err
	}
	defer videoObject.Body.Close()

	videoData, err := io.ReadAll(videoObject.Body)
	if err != nil {
		return err
	}

	videoKeyParts := strings.Split(videoKey, "/")
	videoName := strings.ReplaceAll(videoKeyParts[len(videoKeyParts)-1], ".mp4", "")
	tmpDirRaw := os.TempDir() + "/raw"
	tmpDirProcessed := os.TempDir() + "/processed/" + videoName + "/"

	err = os.MkdirAll(tmpDirRaw, os.ModePerm)
	if err != nil {
		return err
	}
	err = os.MkdirAll(tmpDirProcessed, os.ModePerm)
	if err != nil {
		return err
	}

	tempFilePath := os.TempDir() + "/" + videoKey
	err = os.WriteFile(tempFilePath, videoData, 0666)
	if err != nil {
		return err
	}

	manifestFilePath := tmpDirProcessed + "index.m3u8"
	segmentFilePath := tmpDirProcessed + options.SegmentPrefix + "%03d.ts"

	video := ffmpeg.Input(tempFilePath).Output(manifestFilePath, ffmpeg.KwArgs{
		"vcodec":               options.VideoEncoder,
		"acodec":               options.AudioEncoder,
		"codec":                "copy",
		"start_number":         0,
		"hls_time":             options.HLSFileSize,
		"hls_playlist_type":    "vod",
		"hls_segment_filename": segmentFilePath,
		"hls_list_size":        0,
	}).ErrorToStdOut()

	thumbnail := ffmpeg.Input(tempFilePath).Output(tmpDirProcessed+"thumbnail.jpg", ffmpeg.KwArgs{
		"ss":       options.ThumbnailRefTime,
		"frames:v": "1",
	})

	ffmpeg.MergeOutputs(video, thumbnail).OverWriteOutput().ErrorToStdOut().Run()

	entries, err := os.ReadDir(tmpDirProcessed)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		data, err := os.OpenFile(tmpDirProcessed+entry.Name(), os.O_RDWR, 0666)
		if err != nil {
			return err
		}

		_, err = s.s3Client.PutObject(&s3.PutObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String("processed/" + videoName + "/" + entry.Name()),
			Body:   data,
		})
		if err != nil {
			return err
		}

		err = os.Remove(tmpDirProcessed + entry.Name())
		if err != nil {
			log.Error().Msg(err.Error())
		}
	}

	err = os.Remove(tempFilePath)
	if err != nil {
		log.Error().Msg(err.Error())
	}

	return nil
}

func (s *VideoService) ListAvailableVideos(ctx context.Context, bucket string) (availableVideos []string, err error) {
	out, err := s.s3Client.ListObjects(
		&s3.ListObjectsInput{
			Bucket: aws.String(bucket),
			Prefix: aws.String("processed/"),
		},
	)

	if err != nil {
		return nil, err
	}

	availableVideos = make([]string, 0)
	for _, content := range out.Contents {
		if !strings.HasSuffix(*content.Key, ".m3u8") {
			continue
		}
		videoName := strings.ReplaceAll(*content.Key, "/index.m3u8", "")
		videoDistributionURL := "https://" + os.Getenv("CLOUDFRONT_DIST") + "/" + videoName
		availableVideos = append(availableVideos, videoDistributionURL)
	}

	return availableVideos, nil
}
