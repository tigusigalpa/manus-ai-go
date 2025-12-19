package manusai

import (
	"encoding/base64"
	"fmt"
	"mime"
	"os"
	"path/filepath"
)

func NewAttachmentFromFileID(fileID string) map[string]interface{} {
	return map[string]interface{}{
		"type":    "file_id",
		"file_id": fileID,
	}
}

func NewAttachmentFromURL(url string) map[string]interface{} {
	return map[string]interface{}{
		"type": "url",
		"url":  url,
	}
}

func NewAttachmentFromBase64(base64Data, mimeType string) map[string]interface{} {
	return map[string]interface{}{
		"type":      "data",
		"data":      base64Data,
		"mime_type": mimeType,
	}
}

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
		mimeType = "application/octet-stream"
	}

	base64Data := base64.StdEncoding.EncodeToString(content)

	return NewAttachmentFromBase64(base64Data, mimeType), nil
}
