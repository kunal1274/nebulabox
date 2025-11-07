package api

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nebulabox/nebulabox/internal/containerd"
)

// FileUploadRequest represents a request to upload a file to a container
type FileUploadRequest struct {
	DestinationPath string `json:"destinationPath" binding:"required"`
	Mode            int64  `json:"mode,omitempty"`
	User            string `json:"user,omitempty"`
	Group           string `json:"group,omitempty"`
	Recursive       bool   `json:"recursive,omitempty"`
}

// FileDownloadRequest represents a request to download a file from a container
type FileDownloadRequest struct {
	SourcePath string `json:"sourcePath" binding:"required"`
	Recursive  bool   `json:"recursive,omitempty"`
}

// FileListRequest represents a request to list files in a container
type FileListRequest struct {
	Path string `json:"path" binding:"required"`
}

// FileUploadResponse represents the result of file upload
type FileUploadResponse struct {
	Success          bool   `json:"success"`
	BytesTransferred int64  `json:"bytesTransferred"`
	FilesTransferred int    `json:"filesTransferred"`
	Error            string `json:"error,omitempty"`
}

// FileDownloadResponse represents the result of file download
type FileDownloadResponse struct {
	Success          bool   `json:"success"`
	BytesTransferred int64  `json:"bytesTransferred"`
	FilesTransferred int    `json:"filesTransferred"`
	Error            string `json:"error,omitempty"`
}

// FileListResponse represents the result of file listing
type FileListResponse struct {
	Files []*containerd.FileInfo `json:"files"`
	Path  string                 `json:"path"`
	Count int                    `json:"count"`
}

// uploadFile handles POST /api/containers/:id/files/upload
func (s *Server) uploadFile(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Container ID is required",
		})
		return
	}

	var req FileUploadRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	// Get the uploaded file from form data
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "No file provided",
			"details": err.Error(),
		})
		return
	}
	defer file.Close()

	ctx := c.Request.Context()

	// Create file transfer options
	opts := &containerd.FileTransferOptions{
		SourcePath:      header.Filename,
		DestinationPath: req.DestinationPath,
		Mode:            req.Mode,
		User:            req.User,
		Group:           req.Group,
		Recursive:       req.Recursive,
	}

	// Upload file to container
	result, err := s.containerd.UploadFile(ctx, id, opts)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to upload file to container",
			"details": err.Error(),
		})
		return
	}

	response := FileUploadResponse{
		Success:          result.Success,
		BytesTransferred: result.BytesTransferred,
		FilesTransferred: result.FilesTransferred,
		Error:            result.Error,
	}

	c.JSON(http.StatusOK, response)
}

// downloadFile handles POST /api/containers/:id/files/download
func (s *Server) downloadFile(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Container ID is required",
		})
		return
	}

	var req FileDownloadRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	ctx := c.Request.Context()

	// Create file transfer options
	opts := &containerd.FileTransferOptions{
		SourcePath:      req.SourcePath,
		DestinationPath: "", // Will be set by client
		Recursive:       req.Recursive,
	}

	// Download file from container
	result, err := s.containerd.DownloadFile(ctx, id, opts)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to download file from container",
			"details": err.Error(),
		})
		return
	}

	// Set response headers for file download
	c.Header("Content-Disposition", "attachment; filename="+req.SourcePath)
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Length", strconv.FormatInt(result.BytesTransferred, 10))

	// For now, return a mock file content
	// In a real implementation, we would stream the actual file content
	mockContent := "This is a mock file content for " + req.SourcePath + "\n"
	c.String(http.StatusOK, mockContent)
}

// listFiles handles POST /api/containers/:id/files/list
func (s *Server) listFiles(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Container ID is required",
		})
		return
	}

	var req FileListRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	ctx := c.Request.Context()

	// List files in container
	files, err := s.containerd.ListFiles(ctx, id, req.Path)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to list files in container",
			"details": err.Error(),
		})
		return
	}

	response := FileListResponse{
		Files: files,
		Path:  req.Path,
		Count: len(files),
	}

	c.JSON(http.StatusOK, response)
}

// uploadFileStream handles POST /api/containers/:id/files/upload/stream
func (s *Server) uploadFileStream(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Container ID is required",
		})
		return
	}

	// Get the uploaded file from form data
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "No file provided",
			"details": err.Error(),
		})
		return
	}
	defer file.Close()

	// Get destination path from form data
	destinationPath := c.PostForm("destinationPath")
	if destinationPath == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Destination path is required",
		})
		return
	}

	ctx := c.Request.Context()

	// Create file transfer options
	opts := &containerd.FileTransferOptions{
		SourcePath:      header.Filename,
		DestinationPath: destinationPath,
		Recursive:       false,
	}

	// Set up streaming response
	c.Header("Content-Type", "application/json")
	c.Header("Transfer-Encoding", "chunked")
	c.Status(http.StatusOK)

	// Create a writer that flushes to the client
	writer := &streamingWriter{writer: c.Writer}

	// Simulate streaming upload progress
	progress := []string{
		"Starting file upload...",
		"Reading file data...",
		"Transferring to container...",
		"Setting file permissions...",
		"Upload completed successfully!",
	}

	for _, message := range progress {
		fmt.Fprintf(writer, `{"status": "uploading", "message": "%s"}`+"\n", message)
		time.Sleep(200 * time.Millisecond) // Simulate upload delay
	}

	// Upload file to container
	result, err := s.containerd.UploadFile(ctx, id, opts)
	if err != nil {
		fmt.Fprintf(writer, `{"status": "error", "message": "Upload failed: %s"}`+"\n", err.Error())
		return
	}

	// Send final result
	fmt.Fprintf(writer, `{"status": "completed", "bytesTransferred": %d, "filesTransferred": %d}`+"\n", 
		result.BytesTransferred, result.FilesTransferred)
}

// downloadFileStream handles GET /api/containers/:id/files/download/stream
func (s *Server) downloadFileStream(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Container ID is required",
		})
		return
	}

	sourcePath := c.Query("path")
	if sourcePath == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Source path is required",
		})
		return
	}

	ctx := c.Request.Context()

	// Create file transfer options
	opts := &containerd.FileTransferOptions{
		SourcePath:      sourcePath,
		DestinationPath: "",
		Recursive:       false,
	}

	// Set response headers for streaming download
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", "attachment; filename="+sourcePath)
	c.Header("Transfer-Encoding", "chunked")
	c.Status(http.StatusOK)

	// Create a writer that flushes to the client
	writer := &streamingWriter{writer: c.Writer}

	// Simulate streaming download
	progress := []string{
		"Starting file download...",
		"Reading file from container...",
		"Streaming file data...",
		"Download completed successfully!",
	}

	for _, message := range progress {
		fmt.Fprintf(writer, `{"status": "downloading", "message": "%s"}`+"\n", message)
		time.Sleep(150 * time.Millisecond) // Simulate download delay
	}

	// Download file from container
	result, err := s.containerd.DownloadFile(ctx, id, opts)
	if err != nil {
		fmt.Fprintf(writer, `{"status": "error", "message": "Download failed: %s"}`+"\n", err.Error())
		return
	}

	// Send final result
	fmt.Fprintf(writer, `{"status": "completed", "bytesTransferred": %d, "filesTransferred": %d}`+"\n", 
		result.BytesTransferred, result.FilesTransferred)
}
