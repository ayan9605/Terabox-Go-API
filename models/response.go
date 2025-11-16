package models

// FileInfoResponse represents the file information response
type FileInfoResponse struct {
	FileName     string `json:"file_name" example:"example.pdf"`
	DownloadLink string `json:"download_link" example:"https://example.com/file"`
	Thumbnail    string `json:"thumbnail" example:"https://example.com/thumb.jpg"`
	FileSize     string `json:"file_size" example:"1.50 GB"`
	SizeBytes    int64  `json:"size_bytes" example:"1610612736"`
	ProxyURL     string `json:"proxy_url" example:"http://localhost:8080/proxy?url=..."`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error string `json:"error" example:"Invalid link"`
}

// LinkRequest represents the request body for POST /
type LinkRequest struct {
	Link string `json:"link" binding:"required" example:"https://terabox.com/s/1abc123"`
}

// TeraBoxAPIResponse represents TeraBox API response
type TeraBoxAPIResponse struct {
	Errno  int    `json:"errno"`
	Errmsg string `json:"errmsg"`
	List   []struct {
		ServerFilename string `json:"server_filename"`
		Dlink          string `json:"dlink"`
		Size           int64  `json:"size"`
		Thumbs         struct {
			URL3 string `json:"url3"`
		} `json:"thumbs"`
	} `json:"list"`
}
