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
	// "github.com/rs/zerolog/log"

	ffmpeg "github.com/u2takey/ffmpeg-go"
)

type videoService struct {
	storage istorage.StorageManager
}

func newVideoService(storage istorage.StorageManager) iservice.VideoService {
	return &videoService{
		storage: storage,
	}
}

func (s *videoService) CreateS3PresignedPutURL(ctx context.Context, filename, contentType string) (presignedURL string, err error) {
	return s.storage.UploadRawVideo(filename, contentType)
}

// TODO: after process video and save segments, save filename in dynamo db
func (s *videoService) ProcessVideoWithOptions(ctx context.Context, videoKey string, options domain.VideoOptions) (err error) {
	videoData, err := s.storage.RetrieveRawVideo(videoKey)
	if err != nil {
		return err
	}

	videoName, rawTmpDir, processedTmpDir, err := s.setupProcessEnvironment(videoKey)
	if err != nil {
		return err
	}

	tempFilePath := fmt.Sprintf("%s/%s", rawTmpDir, videoName)
	err = os.WriteFile(tempFilePath, videoData, 0666)
	if err != nil {
		return err
	}

	manifestFilePath := processedTmpDir + "index.m3u8"
	segmentFilePath := processedTmpDir + options.SegmentPrefix + "%03d.ts"

	video := ffmpeg.Input(tempFilePath).Output(manifestFilePath, ffmpeg.KwArgs{
		"vcodec":               options.VideoEncoder,
		"acodec":               options.AudioEncoder,
		"codec":                "copy",
		"start_number":         0,
		"hls_time":             options.HLSFileSize,
		"hls_playlist_type":    "vod",
		"hls_segment_filename": segmentFilePath,
		"hls_list_size":        0,
	})

	thumbnail := ffmpeg.Input(tempFilePath).Output(processedTmpDir+"thumbnail.jpg", ffmpeg.KwArgs{
		"ss":       options.ThumbnailRefTime,
		"frames:v": "1",
	})

	ffmpeg.MergeOutputs(video, thumbnail).OverWriteOutput().ErrorToStdOut().Run()

	entries, err := os.ReadDir(processedTmpDir)
	if err != nil {
		return err
	}

	err = s.storage.UploadProcessedVideo(processedTmpDir, videoName, entries)
	if err != nil {
		return err
	}

	return s.cleanupProcessEnvironment([]string{rawTmpDir, processedTmpDir})
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

func (s *videoService) setupProcessEnvironment(videoKey string) (videoName, rawTempDir, processedTempDir string, err error) {
	_, videoName = filepath.Split(videoKey)
	videoName = strings.ReplaceAll(videoName, filepath.Ext(videoName), "")

	tmpDirRaw := fmt.Sprintf("%s/raw", os.TempDir())
	err = os.MkdirAll(tmpDirRaw, os.ModePerm)
	if err != nil {
		return videoName, rawTempDir, processedTempDir, err
	}

	tmpDirProcessed := fmt.Sprintf("%s/processed/%s/", os.TempDir(), videoName)
	err = os.MkdirAll(tmpDirProcessed, os.ModePerm)
	if err != nil {
		return videoName, rawTempDir, processedTempDir, err
	}

	return videoName, rawTempDir, processedTempDir, nil
}

func (s *videoService) cleanupProcessEnvironment(paths []string) (err error) {
	for _, p := range paths {
		err = os.RemoveAll(p)
		if err != nil {
			return err
		}
	}
	return nil
}
