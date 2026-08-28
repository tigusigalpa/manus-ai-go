package manusai

import (
	"encoding/base64"
	"os"
	"path/filepath"
	"testing"
)

const attachmentContent = "attachment content"

func TestAttachmentHelpers(t *testing.T) {
	fromID := NewAttachmentFromFileID("file_123")
	if fromID["type"] != "file" || fromID["file_id"] != "file_123" {
		t.Fatalf("NewAttachmentFromFileID() = %#v", fromID)
	}

	fromURL := NewAttachmentFromURL("https://example.com/report.pdf")
	if fromURL["file_url"] != "https://example.com/report.pdf" {
		t.Fatalf("NewAttachmentFromURL() = %#v", fromURL)
	}

	fromBase64 := NewAttachmentFromBase64("YQ==", "text/plain")
	if fromBase64["file_data"] != "YQ==" || fromBase64["mime_type"] != "text/plain" {
		t.Fatalf("NewAttachmentFromBase64() = %#v", fromBase64)
	}
}

func TestNewAttachmentFromFilePath(t *testing.T) {
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "report.unknown")
	if err := os.WriteFile(filePath, []byte(attachmentContent), 0o600); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	attachment, err := NewAttachmentFromFilePath(filePath)
	if err != nil {
		t.Fatalf("NewAttachmentFromFilePath() error = %v", err)
	}
	if attachment["mime_type"] != defaultContentType {
		t.Fatalf("mime_type = %v, want %s", attachment["mime_type"], defaultContentType)
	}
	if attachment["file_data"] != base64.StdEncoding.EncodeToString([]byte(attachmentContent)) {
		t.Fatalf("file_data = %v", attachment["file_data"])
	}

	if attachment, err := NewAttachmentFromFilePath(filepath.Join(tempDir, "missing.pdf")); err == nil || attachment != nil {
		t.Fatalf("NewAttachmentFromFilePath() = %#v, %v; want error", attachment, err)
	}
}
