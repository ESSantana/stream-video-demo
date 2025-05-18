package interfaces

import "os"

type StorageManager interface {
	UploadRawVideo(filename, contentType string) (uploadURL string, err error)
	UploadProcessedVideo(basePath, videoID string, processedFiles []os.DirEntry) (err error)
	RetrieveRawVideo(objectKey string) (objectData []byte, err error)
}
