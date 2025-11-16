package handlers

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	// "strings" // REMOVE THIS LINE - not used
	"terabox-api/models"
	"terabox-api/utils"

	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
)

const (
	COOKIE = "ndus=Y4Y9Sg7teHuiSOxGuGKbIV2ymZJJIKC4GYTf7PHe" // Replace with your cookie
)

// GetFileInfo godoc
// @Summary      Get file information (GET)
// @Description  Retrieve TeraBox file information using query parameter
// @Tags         files
// @Accept       json
// @Produce      json
// @Param        url  query  string  true  "TeraBox share URL"
// @Success      200  {object}  models.FileInfoResponse
// @Failure      400  {object}  models.ErrorResponse
// @Failure      500  {object}  models.ErrorResponse
// @Router       /api [get]
func GetFileInfo(c *gin.Context) {
	link := c.Query("url")
	if link == "" {
		c.JSON(400, models.ErrorResponse{Error: "No URL provided. Use ?url=your_terabox_link"})
		return
	}

	processFileInfo(c, link)
}

// PostFileInfo godoc
// @Summary      Get file information (POST)
// @Description  Retrieve TeraBox file information using JSON body
// @Tags         files
// @Accept       json
// @Produce      json
// @Param        request  body  models.LinkRequest  true  "TeraBox link"
// @Success      200  {object}  models.FileInfoResponse
// @Failure      400  {object}  models.ErrorResponse
// @Failure      500  {object}  models.ErrorResponse
// @Router       / [post]
func PostFileInfo(c *gin.Context) {
	var req models.LinkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, models.ErrorResponse{Error: "No link provided in the request body"})
		return
	}

	processFileInfo(c, req.Link)
}

func processFileInfo(c *gin.Context, link string) {
	// Check cache first
	cacheKey := "file:" + link
	if cached, found := utils.GetFromCache(cacheKey); found {
		c.Header("X-Cache-Status", "HIT")
		c.JSON(200, cached)
		return
	}

	// Fetch file info
	fileInfo, err := fetchFileInfo(link, c.Request)
	if err != nil {
		c.JSON(400, models.ErrorResponse{Error: err.Error()})
		return
	}

	// Add proxy URL
	fileInfo.ProxyURL = fmt.Sprintf("%s://%s/proxy?url=%s&file_name=%s",
		getScheme(c),
		c.Request.Host,
		url.QueryEscape(fileInfo.DownloadLink),
		url.QueryEscape(fileInfo.FileName),
	)

	// Cache the result
	utils.SetCache(cacheKey, fileInfo)

	c.Header("X-Cache-Status", "MISS")
	c.JSON(200, fileInfo)
}

func fetchFileInfo(link string, r *http.Request) (*models.FileInfoResponse, error) {
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return nil
		},
	}

	// Initial request to get final URL
	req, err := http.NewRequest("GET", link, nil)
	if err != nil {
		return nil, fmt.Errorf("invalid link: %v", err)
	}

	setHeaders(req)
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch link: %v", err)
	}
	defer resp.Body.Close()

	finalURL := resp.Request.URL.String()
	parsedURL, err := url.Parse(finalURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse URL: %v", err)
	}

	surl := parsedURL.Query().Get("surl")
	if surl == "" {
		return nil, fmt.Errorf("invalid link. Please check the link")
	}

	// Fetch page to extract tokens
	req, _ = http.NewRequest("GET", finalURL, nil)
	setHeaders(req)
	resp, err = client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch page: %v", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	text := string(body)

	jsToken := utils.FindBetween(text, "fn%28%22", "%22%29")
	logid := utils.FindBetween(text, "dp-logid=", "&")
	bdstoken := utils.FindBetween(text, `bdstoken":"`, `"`)

	if jsToken == "" || logid == "" || bdstoken == "" {
		return nil, fmt.Errorf("failed to extract required tokens")
	}

	// Build API URL
	params := url.Values{}
	params.Set("app_id", "250528")
	params.Set("web", "1")
	params.Set("channel", "dubox")
	params.Set("clienttype", "0")
	params.Set("jsToken", jsToken)
	params.Set("dp-logid", logid)
	params.Set("page", "1")
	params.Set("num", "20")
	params.Set("by", "name")
	params.Set("order", "asc")
	params.Set("site_referer", finalURL)
	params.Set("shorturl", surl)
	params.Set("root", "1")

	apiURL := "https://www.terabox.com/share/list?" + params.Encode()
	req, _ = http.NewRequest("GET", apiURL, nil)
	setHeaders(req)

	resp, err = client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch file list: %v", err)
	}
	defer resp.Body.Close()

	var apiResp models.TeraBoxAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %v", err)
	}

	if apiResp.Errno != 0 || len(apiResp.List) == 0 {
		return nil, fmt.Errorf("failed to retrieve file list: %s", apiResp.Errmsg)
	}

	file := apiResp.List[0]
	return &models.FileInfoResponse{
		FileName:     file.ServerFilename,
		DownloadLink: file.Dlink,
		Thumbnail:    file.Thumbs.URL3,
		FileSize:     utils.GetSize(file.Size),
		SizeBytes:    file.Size,
	}, nil
}

func setHeaders(req *http.Request) {
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9,hi;q=0.8")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")
	req.Header.Set("Cookie", COOKIE)
}

func getScheme(c *gin.Context) string {
	if c.Request.TLS != nil || c.GetHeader("X-Forwarded-Proto") == "https" {
		return "https"
	}
	return "http"
}
