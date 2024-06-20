package controllers

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func DownloadFile(c *gin.Context) {
	// Get file location from request parameter (adjust if needed)
	path := c.Query("path")

	// Get base path (consider using a configurable base path)
	// basePath, err := os.Getwd()
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": true, "message": "Error getting base path"})
	// 	return
	// }

	// // Construct full file path
	// fileLocation := filepath.Join(basePath, path)

	log.Printf("Extracted file path: %s\n", path)

	// Check if file exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"error": true, "message": "File not found"})
		return
	}

	// Open the file for reading
	file, err := os.Open(path)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": true, "message": "Error opening file"})
		return
	}
	defer file.Close() // Ensure file is closed even on errors

	// Get file information for content type and size
	fileInfo, err := file.Stat()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": true, "message": "Error getting file information"})
		return
	}

	buffer := make([]byte, 512)
	_, err = file.Read(buffer)
	if err != nil {
		fmt.Println(err)
	}

	contentType := http.DetectContentType(buffer)

	// // Set content type based on file extension (improve with more extensions)
	// contentType := "application/octet-stream" // Default for unknown types
	// switch filepath.Ext(path) {
	// case ".jpg", ".jpeg":
	// 	contentType = "image/jpeg"
	// case ".png":
	// 	contentType = "image/png"
	// case ".pdf":
	// 	contentType = "application/pdf"
	// case ".txt":
	// 	contentType = "text/plain"
	// case ".docx":
	// 	contentType = "application/vnd.openxmlformats-officedocument.wordprocessingml.document"
	// case ".doc":
	// 	contentType = "application/msword"
	// case ".ppt":
	// 	contentType = "application/vnd.ms-powerpoint"
	// case ".pptx":
	// 	contentType = "application/vnd.openxmlformats-officedocument.presentationml.presentation"
	// case ".xls":
	// 	contentType = "application/vnd.ms-excel"
	// case ".xlsx":
	// 	contentType = "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
	// case ".zip":
	// 	contentType = "application/zip"
	// case ".rar":
	// 	contentType = "application/x-rar-compressed"
	// }

	// Set content disposition header
	contentDisposition := fmt.Sprintf("attachment; filename=%s", filepath.Base(path)) // Use base filename
	c.Writer.Header().Set("Content-Disposition", contentDisposition)
	c.Writer.Header().Set("Content-Type", contentType)
	c.Writer.Header().Set("Content-Length", fmt.Sprintf("%d", fileInfo.Size())) // Set content length for progress

	// Stream the file to the client
	_, err = io.Copy(c.Writer, file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": true, "message": "Error streaming file"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"error":    false,
		"messages": "File downloaded successfully",
	})
}
