package services

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ESSantana/streaming-test/internal/domain"
	iservice "github.com/ESSantana/streaming-test/internal/services/interfaces"
	istorage "github.com/ESSantana/streaming-test/internal/storage/interfaces"
	"github.com/google/uuid"

	// "github.com/rs/zerolog/log"

	ffmpeg "github.com/u2takey/ffmpeg-go"
)

type videoService struct {
	storageManager istorage.StorageManager
}

func newVideoService(storageManager istorage.StorageManager) iservice.VideoService {
	return &videoService{
		storageManager: storageManager,
	}
}

func (s *videoService) UploadRawVideo(ctx context.Context, filename, contentType string) (presignedURL string, err error) {
	return s.storageManager.UploadRawVideo(filename, contentType)
}

func (s *videoService) ProcessVideoWithOptions(ctx context.Context, videoKey string, options domain.VideoOptions) (err error) {
	videoData, err := s.storageManager.RetrieveRawVideo(videoKey)
	if err != nil {
		return err
	}

	_, tmpDirRaw, tmpDirProcessed, err := s.setupProcessEnvironment(videoKey, videoData)
	if err != nil {
		return err
	}

	manifestFilePath := fmt.Sprintf("%s/index.m3u8", tmpDirProcessed)
	segmentFilePath := fmt.Sprintf("%s/%s%03d.ts", tmpDirProcessed, options.SegmentPrefix)

	video := ffmpeg.Input(tmpDirRaw).Output(manifestFilePath, ffmpeg.KwArgs{
		"vcodec":               options.VideoEncoder,
		"acodec":               options.AudioEncoder,
		"codec":                "copy",
		"start_number":         0,
		"hls_time":             options.HLSFileSize,
		"hls_playlist_type":    "vod",
		"hls_segment_filename": segmentFilePath,
		"hls_list_size":        0,
	})

	thumbnail := ffmpeg.Input(tmpDirRaw).Output(fmt.Sprintf("%s/thumbnail.jpg", tmpDirProcessed),
		ffmpeg.KwArgs{
			"ss":       options.ThumbnailRefTime,
			"frames:v": "1",
		},
	)

	err = ffmpeg.MergeOutputs(video, thumbnail).OverWriteOutput().ErrorToStdOut().Run()
	if err != nil {
		return err
	}

	entries, err := os.ReadDir(tmpDirProcessed)
	if err != nil {
		return err
	}

	id, err := uuid.NewV7()
	if err != nil {
		return err
	}

	err = s.storageManager.UploadProcessedVideo(tmpDirProcessed, id.String(), entries)
	if err != nil {
		return err
	}

	// TODO: after process video and save segments, save filename in dynamo db
	// use id, video name and s3 path

	return os.RemoveAll(os.TempDir())
}

// TODO: get it from dynamo db instead of list from s3
func (s *videoService) ListAvailableVideos(ctx context.Context, bucket string) (availableVideos []string, err error) {
	// out, err := s.s3Client.ListObjects(
	// 	&s3.ListObjectsInput{
	// 		Bucket: aws.String(bucket),
	// 		Prefix: aws.String("processed/"),
	// 	},
	// )

	// if err != nil {
	// 	return nil, err
	// }

	// availableVideos = make([]string, 0)
	// for _, content := range out.Contents {
	// 	if !strings.HasSuffix(*content.Key, ".m3u8") {
	// 		continue
	// 	}
	// 	videoName := strings.ReplaceAll(*content.Key, "/index.m3u8", "")
	// 	videoDistributionURL := "https://" + os.Getenv("CLOUDFRONT_DIST") + "/" + videoName
	// 	availableVideos = append(availableVideos, videoDistributionURL)
	// }

	return availableVideos, nil
}

func (s *videoService) setupProcessEnvironment(videoKey string, videoData []byte) (videoName, tmpDirRaw, tmpDirProcessed string, err error) {
	_, videoName = filepath.Split(videoKey)
	videoName = strings.ReplaceAll(videoName, filepath.Ext(videoName), "")

	tmpDirRaw = fmt.Sprintf("%s/raw/%s", os.TempDir(), videoName)
	err = os.MkdirAll(tmpDirRaw, os.ModePerm)
	if err != nil {
		return videoName, tmpDirRaw, tmpDirProcessed, err
	}

	tmpDirProcessed = fmt.Sprintf("%s/processed/%s", os.TempDir(), videoName)
	err = os.MkdirAll(tmpDirProcessed, os.ModePerm)
	if err != nil {
		return videoName, tmpDirRaw, tmpDirProcessed, err
	}

	err = os.WriteFile(tmpDirRaw, videoData, 0666)
	if err != nil {
		return videoName, tmpDirRaw, tmpDirProcessed, err
	}

	return videoName, tmpDirRaw, tmpDirProcessed, nil
}
