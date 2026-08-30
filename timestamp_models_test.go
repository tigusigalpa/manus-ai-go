package manusai

import (
	"encoding/json"
	"net/http"
	"testing"
)

func TestResponseModelsDecodeStringTimestamps(t *testing.T) {
	tests := []struct {
		name string
		data string
		out  interface{}
	}{
		{"project", `{"id":"project-1","created_at":"1725033841"}`, &Project{}},
		{"skill", `{"id":"skill-1","created_at":"2026-08-30T16:04:01.419Z","updated_at":"1725033841419"}`, &Skill{}},
		{"webhook", `{"id":"webhook-1","created_at":"1725033841419"}`, &Webhook{}},
		{"usage", `{"task_id":"task-1","created_at":"1725033841"}`, &UsageRecord{}},
		{"credits", `{"next_refresh_time":"1725033841","refresh_interval":"3600","current_period_end":"1725033841419"}`, &AvailableCredits{}},
		{"website status", `{"status_updated_at":"2026-08-30T16:04:01.419Z"}`, &WebsiteStatusResponse{}},
		{"website checkpoint", `{"version_id":"version-1","created_at":"1725033841419"}`, &WebsiteCheckpoint{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := json.Unmarshal([]byte(tt.data), tt.out); err != nil {
				t.Fatalf("json.Unmarshal() error = %v", err)
			}
		})
	}

	project := &Project{}
	if err := json.Unmarshal([]byte(`{"created_at":"invalid"}`), project); err == nil {
		t.Fatal("Project accepted an invalid created_at timestamp")
	}
}

func TestFileEndpointsDecodeStringTimestamps(t *testing.T) {
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/v2/file.upload":
			_, _ = w.Write([]byte(`{"ok":true,"file":{"id":"file-1","created_at":"1725033841"},"upload_url":"https://upload.example.com","upload_expires_at":"1725033841419"}`))
		case "/v2/file.detail":
			_, _ = w.Write([]byte(`{"ok":true,"file":{"id":"file-1","created_at":"2026-08-30T16:04:01.419Z","expires_at":"1725033841419"}}`))
		default:
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
	})

	file, err := client.CreateFile("report.pdf")
	if err != nil || file.ExpiresAt != 1725033841419 {
		t.Fatalf("CreateFile() = (%#v, %v)", file, err)
	}
	detail, err := client.GetFile("file-1")
	if err != nil || detail.CreatedAt != 1788105841419 || detail.ExpiresAt != 1725033841419 {
		t.Fatalf("GetFile() = (%#v, %v)", detail, err)
	}
}
