package dto

type VideoUploadRequest struct {
	Filename    string `json:"filename"`
	ContentType string `json:"content_type"`
}

type VideoUploadResponse struct {
	UploadURL string `json:"upload_url"`
}

type ListVideosResponse struct {
	Videos []string `json:"videos"`
}

type VideoDistributionResponse struct {
	VideoURL string `json:"video_url"`
}
