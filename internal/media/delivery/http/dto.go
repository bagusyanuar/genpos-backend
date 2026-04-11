package http

type UploadResponse struct {
	URL string `json:"url"`
}

func ToUploadResponse(url string) UploadResponse {
	return UploadResponse{
		URL: url,
	}
}
