package services

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/ESSantana/streaming-test/internal/domain"
	"github.com/ESSantana/streaming-test/internal/domain/models"
	irepository "github.com/ESSantana/streaming-test/internal/repositories/interfaces"
	iservice "github.com/ESSantana/streaming-test/internal/services/interfaces"
	istorage "github.com/ESSantana/streaming-test/internal/storage/interfaces"
	"github.com/google/uuid"

	ffmpeg "github.com/u2takey/ffmpeg-go"
)

type videoService struct {
	storageManager    istorage.StorageManager
	repositoryManager irepository.RepositoryManager
}

func newVideoService(storageManager istorage.StorageManager, repositoryManager irepository.RepositoryManager) iservice.VideoService {
	return &videoService{
		storageManager:    storageManager,
		repositoryManager: repositoryManager,
	}
}

func (s *videoService) UploadRawVideo(ctx context.Context, filename, contentType string) (presignedURL string, err error) {
	return s.storageManager.UploadRawVideo(filename, contentType)
}

func (s *videoService) ProcessVideoWithOptions(ctx context.Context, videoKey string, options domain.VideoOptions) (err error) {
	videoData, err := s.storageManager.RetrieveRawVideo(videoKey)
	if err != nil {
		return errors.New("error at retrieve raw video: " + err.Error())
	}

	videoName, tmpDirRaw, tmpDirProcessed, err := s.setupProcessEnvironment(videoKey, videoData)
	if err != nil {
		return errors.New("error at setup environment: " + err.Error())
	}

	manifestFilePath := fmt.Sprintf("%s/index.m3u8", tmpDirProcessed)
	segmentFilePath := fmt.Sprintf("%s/%s.ts", tmpDirProcessed, options.SegmentPrefix+"%03d")

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
		return errors.New("error at process request: " + err.Error())
	}

	entries, err := os.ReadDir(tmpDirProcessed)
	if err != nil {
		return errors.New("error at read temp processed dir: " + err.Error())
	}

	id, err := uuid.NewV7()
	if err != nil {
		return errors.New("error at create uuid v7:" + err.Error())
	}

	err = s.storageManager.UploadProcessedVideo(tmpDirProcessed, id.String(), entries)
	if err != nil {
		return errors.New("error at upload processed video: " + err.Error())
	}

	videoRepository := s.repositoryManager.NewVideoRepository()
	err = videoRepository.Save(ctx, models.Video{
		VideoId:   id.String(),
		VideoName: videoName,
		Manifest:  fmt.Sprintf("processed/%s/index.m3u8", id.String()),
		Thumbnail: fmt.Sprintf("processed/%s/thumbnail.jpg", id.String()),
		CreatedAt: uint64(time.Now().Unix()),
	})

	if err != nil {
		return errors.New("error at video repository save batch: " + err.Error())
	}

	return os.RemoveAll(fmt.Sprintf("%s/%s", os.TempDir(), videoName))
}

func (s *videoService) ListAvailableVideos(ctx context.Context) (availableVideos []models.Video, err error) {
	videoRepository := s.repositoryManager.NewVideoRepository()
	return videoRepository.ListAvailableVideos(ctx)
}

func (s *videoService) setupProcessEnvironment(videoKey string, videoData []byte) (videoName, tempRawVideo, tmpDirProcessed string, err error) {
	_, videoName = filepath.Split(videoKey)
	videoExtension := filepath.Ext(videoName)
	videoName = strings.ReplaceAll(videoName, videoExtension, "")

	tmpDirRaw := fmt.Sprintf("%s/%s/raw", os.TempDir(), videoName)
	err = os.MkdirAll(tmpDirRaw, os.ModePerm)
	if err != nil {
		return videoName, tempRawVideo, tmpDirProcessed, err
	}

	tmpDirProcessed = fmt.Sprintf("%s/%s/processed", os.TempDir(), videoName)
	err = os.MkdirAll(tmpDirProcessed, os.ModePerm)
	if err != nil {
		return videoName, tempRawVideo, tmpDirProcessed, err
	}

	tempRawVideo = fmt.Sprintf("%s/raw%s", tmpDirRaw, videoExtension)
	err = os.WriteFile(tempRawVideo, videoData, 0666)
	if err != nil {
		return videoName, tempRawVideo, tmpDirProcessed, err
	}

	return videoName, tempRawVideo, tmpDirProcessed, nil
}
