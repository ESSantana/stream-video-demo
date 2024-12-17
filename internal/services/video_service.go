package services

import (
	"context"
	"fmt"
	"strconv"

	"github.com/ESSantana/streaming-test/internal/services/interfaces"
)

type VideoService struct {
}

func NewVideoService() interfaces.VideoService {
	return &VideoService{}
}

func (s *VideoService) UploadVideo(ctx context.Context, filename string, extension string, data []byte) (err error) {
	fmt.Println(filename + extension)
	fmt.Println("Video length: " + strconv.Itoa(len(data)))
	return nil
}
