package model

// UploadResponse response upload image to S3
type UploadResponse struct {
	URL      string `json:"url"`
	FileName string `json:"file_name"`
	Type     string `json:"type"`
}

// CoverUploadResponse response upload image to S3
type CoverUploadResponse struct {
	CoverURL string `json:"cover_url"`
	ThumbURL string `json:"thumb_url"`
	MSG      string `json:"msg"`
}
