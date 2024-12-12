package routes

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/google/uuid"
	ffmpeg "github.com/u2takey/ffmpeg-go"
)

const (
	tempDir   = "./tmp/upload_video"
	outputDir = "./videos/output"
)


func UploadVideo(w http.ResponseWriter, r *http.Request) {
	fileID := uuid.New().String()

	err := r.ParseMultipartForm(20 << 20)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	videoFile, videoHeader, err := r.FormFile("upload_video")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer videoFile.Close()

	fmt.Printf("File Name: %s, Size: %v", videoHeader.Filename, videoHeader.Size)

	data, err := io.ReadAll(videoFile)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if len(data) < 1 {
		http.Error(w, "Error at upload video", http.StatusBadRequest)
		return
	}

	err = os.MkdirAll(tempDir, os.ModePerm)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tempFilePath := fmt.Sprintf("%s/%s.mp4", tempDir, fileID)
	err = os.WriteFile(tempFilePath, data, os.ModeAppend)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func SplitVideo(w http.ResponseWriter, r *http.Request) {
	d, err := os.ReadDir(tempDir)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	for _, entry := range d {
		tempFilePath := fmt.Sprintf("%s/%s", tempDir, entry.Name())
		fileName := strings.ReplaceAll(entry.Name(), ".mp4", "")

		outputDir := fmt.Sprintf("%s/%s", outputDir, fileName)
		err = os.MkdirAll(outputDir, os.ModePerm)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		manifestFilePath := outputDir + "/index.m3u8"
		segmentFilePath := outputDir + "/segment%03d.ts"

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
	}
}

