package domain

type VideoOptions struct {
	// See more info about in https://www.ffmpeg.org/ffmpeg-codecs.html#Audio-Encoders
	AudioEncoder string
	// See more info about in https://www.ffmpeg.org/ffmpeg-codecs.html#Video-Encoders
	VideoEncoder     string
	HLSFileSize      int
	SegmentPrefix    string
	ThumbnailRefTime string
}

func DefaultVideoOptions() VideoOptions {
	options := VideoOptions{
		AudioEncoder:     "aac",
		VideoEncoder:     "libx264",
		HLSFileSize:      30,
		SegmentPrefix:    "segment",
		ThumbnailRefTime: "00:00:10",
	}
	return options
}
