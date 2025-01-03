package routes

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	ffmpeg "github.com/u2takey/ffmpeg-go"
)

const (
	tempDir   = "./tmp/upload_video"
	outputDir = "./videos/output"
)

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

