package handlers

import (
	"io"
	"net/http"
	"terabox-api/models"

	"github.com/gin-gonic/gin"
)

// ProxyDownload godoc
// @Summary      Proxy download file
// @Description  Download file through proxy with range support
// @Tags         proxy
// @Produce      octet-stream
// @Param        url        query  string  true  "Direct download URL"
// @Param        file_name  query  string  false "File name"
// @Success      200  {file}    binary
// @Success      206  {file}    binary
// @Failure      400  {object}  models.ErrorResponse
// @Failure      502  {object}  models.ErrorResponse
// @Router       /proxy [get]
func ProxyDownload(c *gin.Context) {
	downloadURL := c.Query("url")
	fileName := c.DefaultQuery("file_name", "download")

	if downloadURL == "" {
		c.JSON(400, models.ErrorResponse{Error: "No URL provided for proxy"})
		return
	}

	client := &http.Client{}
	req, err := http.NewRequest("GET", downloadURL, nil)
	if err != nil {
		c.JSON(400, models.ErrorResponse{Error: "Invalid URL"})
		return
	}

	// Copy headers
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
	req.Header.Set("Referer", "https://terabox.com/")
	req.Header.Set("Cookie", COOKIE)

	// Handle Range requests
	if rangeHeader := c.GetHeader("Range"); rangeHeader != "" {
		req.Header.Set("Range", rangeHeader)
	}

	resp, err := client.Do(req)
	if err != nil {
		c.JSON(502, models.ErrorResponse{Error: "Failed to fetch download"})
		return
	}
	defer resp.Body.Close()

	// Set response headers
	c.Header("Content-Type", resp.Header.Get("Content-Type"))
	c.Header("Content-Disposition", `inline; filename="`+fileName+`"`)
	c.Header("Accept-Ranges", "bytes")
	c.Header("Cache-Control", "public, max-age=3600")

	if resp.Header.Get("Content-Range") != "" {
		c.Header("Content-Range", resp.Header.Get("Content-Range"))
	}
	if resp.Header.Get("Content-Length") != "" {
		c.Header("Content-Length", resp.Header.Get("Content-Length"))
	}

	// Stream the response
	c.Status(resp.StatusCode)
	io.Copy(c.Writer, resp.Body)
}
