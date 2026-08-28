package manusai

import (
	"encoding/base64"
	"fmt"
	"mime"
	"os"
	"path/filepath"
)

// NewAttachmentFromFileID creates a file attachment that refers to an uploaded Manus file.
func NewAttachmentFromFileID(fileID string) map[string]interface{} {
	return map[string]interface{}{
		"type":    "file",
		"file_id": fileID,
	}
}

// NewAttachmentFromURL creates a file attachment that refers to a publicly accessible URL.
func NewAttachmentFromURL(url string) map[string]interface{} {
	return map[string]interface{}{
		"type":     "file",
		"file_url": url,
	}
}

// NewAttachmentFromBase64 creates a file attachment from base64-encoded data and its MIME type.
func NewAttachmentFromBase64(base64Data, mimeType string) map[string]interface{} {
	return map[string]interface{}{
		"type":      "file",
		"file_data": base64Data,
		"mime_type": mimeType,
	}
}

// NewAttachmentFromFilePath reads a local file and creates a base64-encoded attachment.
func NewAttachmentFromFilePath(filePath string) (map[string]interface{}, error) {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("file not found: %s", filePath)
	}

	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	mimeType := mime.TypeByExtension(filepath.Ext(filePath))
	if mimeType == "" {
		mimeType = defaultContentType
	}

	base64Data := base64.StdEncoding.EncodeToString(content)

	return NewAttachmentFromBase64(base64Data, mimeType), nil
}
