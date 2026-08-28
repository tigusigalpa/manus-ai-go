package manusai

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"
)

func newTestClient(t *testing.T, handler http.HandlerFunc) *Client {
	t.Helper()
	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	client, err := NewClient("test-api-key", WithBaseURL(server.URL))
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}
	return client
}

func decodeBody(t *testing.T, r *http.Request) map[string]interface{} {
	t.Helper()
	defer r.Body.Close()
	var payload map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		t.Fatalf("decode request body: %v", err)
	}
	return payload
}

func TestNewClient(t *testing.T) {
	t.Run("uses defaults", func(t *testing.T) {
		client, err := NewClient("test-api-key")
		if err != nil {
			t.Fatalf("NewClient() error = %v", err)
		}
		if client.apiKey != "test-api-key" || client.baseURL != DefaultBaseURL {
			t.Fatalf("unexpected client configuration: %#v", client)
		}
	})

	tests := []struct {
		name string
		opts []ClientOption
	}{
		{"empty API key", nil},
		{"nil option", []ClientOption{nil}},
		{"nil HTTP client", []ClientOption{WithHTTPClient(nil)}},
		{"invalid base URL", []ClientOption{WithBaseURL("not a URL")}},
		{"unsupported base URL scheme", []ClientOption{WithBaseURL("ftp://example.com")}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			key := "test-api-key"
			if tt.name == "empty API key" {
				key = ""
			}
			client, err := NewClient(key, tt.opts...)
			if err == nil || client != nil {
				t.Fatalf("NewClient() = (%v, %v), want nil client and an error", client, err)
			}
		})
	}

	t.Run("applies timeout after a nil HTTP client option", func(t *testing.T) {
		client, err := NewClient("test-api-key", WithHTTPClient(nil), WithTimeout(time.Second))
		if err != nil {
			t.Fatalf("NewClient() error = %v", err)
		}
		if client.httpClient.Timeout != time.Second {
			t.Fatalf("timeout = %v, want %v", client.httpClient.Timeout, time.Second)
		}
	})
}

func TestCreateTaskBuildsV2Request(t *testing.T) {
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/v2/task.create" {
			t.Fatalf("request = %s %s", r.Method, r.URL.Path)
		}
		if r.Header.Get("x-manus-api-key") != "test-api-key" {
			t.Fatal("API key header is missing")
		}

		payload := decodeBody(t, r)
		if payload["agent_profile"] != AgentProfileManus16 || payload["title"] != "Review" {
			t.Fatalf("unexpected task payload: %#v", payload)
		}
		message := payload["message"].(map[string]interface{})
		content := message["content"].([]interface{})
		if len(content) != 2 || content[0].(map[string]interface{})["text"] != "Review this" {
			t.Fatalf("unexpected message content: %#v", content)
		}

		_, _ = w.Write([]byte(`{"ok":true,"request_id":"req_1","task_id":"task_123","task_title":"Review","task_url":"https://manus.ai/task/123"}`))
	})

	task, err := client.CreateTask("Review this", &TaskOptions{
		AgentProfile: AgentProfileManus16,
		Title:        "Review",
		Attachments:  []interface{}{NewAttachmentFromFileID("file_123")},
	})
	if err != nil {
		t.Fatalf("CreateTask() error = %v", err)
	}
	if task.TaskID != "task_123" || !task.OK {
		t.Fatalf("unexpected task: %#v", task)
	}
}

func TestCreateTaskValidatesInput(t *testing.T) {
	client, err := NewClient("test-api-key")
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.CreateTask("", nil)
	if err == nil {
		t.Fatal("CreateTask() with an empty prompt returned no error")
	}

	_, err = client.CreateTask("valid prompt", &TaskOptions{Attachments: []interface{}{"not an attachment"}})
	if err == nil {
		t.Fatal("CreateTask() accepted an invalid attachment")
	}

	_, err = client.SendMessage("task_123", "valid message", []interface{}{"not an attachment"})
	if err == nil {
		t.Fatal("SendMessage() accepted an invalid attachment")
	}
}

func TestTaskAndFileRequests(t *testing.T) {
	tests := []struct {
		name     string
		call     func(*Client) error
		method   string
		path     string
		query    url.Values
		response string
	}{
		{"list tasks", func(c *Client) error {
			_, err := c.GetTasks(&TaskFilters{Limit: 10, Order: "desc", Scope: "all"})
			return err
		}, http.MethodGet, "/v2/task.list", url.Values{"limit": {"10"}, "order": {"desc"}, "scope": {"all"}}, `{"ok":true,"tasks":[]}`},
		{"get task", func(c *Client) error { _, err := c.GetTask("task_123"); return err }, http.MethodGet, "/v2/task.detail", url.Values{"task_id": {"task_123"}}, `{"ok":true,"id":"task_123","agent_status":"running"}`},
		{"update task", func(c *Client) error {
			title := "New title"
			_, err := c.UpdateTask("task_123", &TaskUpdate{Title: &title})
			return err
		}, http.MethodPost, "/v2/task.update", nil, `{"ok":true,"id":"task_123"}`},
		{"delete task", func(c *Client) error { _, err := c.DeleteTask("task_123"); return err }, http.MethodPost, "/v2/task.delete", nil, `{"ok":true,"deleted":true}`},
		{"create file", func(c *Client) error { _, err := c.CreateFile("report.pdf"); return err }, http.MethodPost, "/v2/file.upload", nil, `{"ok":true,"file_id":"file_123"}`},
		{"list files", func(c *Client) error { _, err := c.ListFiles(5, "next"); return err }, http.MethodGet, "/v2/file.list", url.Values{"limit": {"5"}, "cursor": {"next"}}, `{"ok":true,"files":[]}`},
		{"get file", func(c *Client) error { _, err := c.GetFile("file_123"); return err }, http.MethodGet, "/v2/file.detail", url.Values{"file_id": {"file_123"}}, `{"ok":true,"file_id":"file_123"}`},
		{"delete file", func(c *Client) error { _, err := c.DeleteFile("file_123"); return err }, http.MethodPost, "/v2/file.delete", nil, `{"ok":true,"deleted":true}`},
		{"create webhook", func(c *Client) error {
			_, err := c.CreateWebhook(&WebhookConfig{URL: "https://example.com/hook", Events: []string{"task_stopped"}})
			return err
		}, http.MethodPost, "/v2/webhook.create", nil, `{"ok":true,"webhook_id":"webhook_123"}`},
		{"delete webhook", func(c *Client) error { return c.DeleteWebhook("webhook_123") }, http.MethodPost, "/v2/webhook.delete", nil, ""},
		{"list messages", func(c *Client) error { _, err := c.ListMessages("task_123", 10, "next", "desc", true); return err }, http.MethodGet, "/v2/task.listMessages", url.Values{"task_id": {"task_123"}, "limit": {"10"}, "cursor": {"next"}, "order": {"desc"}, "verbose": {"true"}}, `{"ok":true,"messages":[]}`},
		{"send message", func(c *Client) error { _, err := c.SendMessage("task_123", "Continue", nil); return err }, http.MethodPost, "/v2/task.sendMessage", nil, `{"ok":true}`},
		{"stop task", func(c *Client) error { _, err := c.StopTask("task_123"); return err }, http.MethodPost, "/v2/task.stop", nil, `{"ok":true}`},
		{"confirm action", func(c *Client) error {
			_, err := c.ConfirmAction("task_123", "event_123", map[string]interface{}{"confirmed": true})
			return err
		}, http.MethodPost, "/v2/task.confirmAction", nil, `{"ok":true}`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
				if r.Method != tt.method || r.URL.Path != tt.path {
					t.Fatalf("request = %s %s, want %s %s", r.Method, r.URL.Path, tt.method, tt.path)
				}
				if got := r.URL.Query(); tt.query != nil && got.Encode() != tt.query.Encode() {
					t.Fatalf("query = %s, want %s", got.Encode(), tt.query.Encode())
				}
				w.WriteHeader(http.StatusOK)
				_, _ = io.WriteString(w, tt.response)
			})
			if err := tt.call(client); err != nil {
				t.Fatalf("request error = %v", err)
			}
		})
	}
}

func TestUploadFileContentAndErrors(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut || r.Header.Get("Content-Type") != "text/plain" {
			t.Fatalf("unexpected upload request: %s, %s", r.Method, r.Header.Get("Content-Type"))
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()
	client, _ := NewClient("test-api-key")
	if err := client.UploadFileContent(server.URL, []byte("content"), "text/plain"); err != nil {
		t.Fatalf("UploadFileContent() error = %v", err)
	}

	for status, wantType := range map[int]string{http.StatusUnauthorized: "authentication", http.StatusBadRequest: "validation", http.StatusInternalServerError: "manus-ai"} {
		t.Run(http.StatusText(status), func(t *testing.T) {
			failingClient := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(status)
				_, _ = w.Write([]byte("failure"))
			})
			_, err := failingClient.GetTask("task_123")
			if err == nil || !strings.Contains(err.Error(), wantType) {
				t.Fatalf("error = %v, want %s error", err, wantType)
			}
		})
	}
}
