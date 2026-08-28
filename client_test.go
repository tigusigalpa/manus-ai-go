package manusai

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"
)

const testAPIKey = "test-api-key"

type roundTripperFunc func(*http.Request) (*http.Response, error)

func (f roundTripperFunc) RoundTrip(request *http.Request) (*http.Response, error) {
	return f(request)
}

type failingReader struct{}

func (failingReader) Read([]byte) (int, error) {
	return 0, errors.New("read failure")
}

func newTestClient(t *testing.T, handler http.HandlerFunc) *Client {
	t.Helper()
	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	client, err := NewClient(testAPIKey, WithBaseURL(server.URL))
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}
	return client
}

func decodeBody(t *testing.T, r *http.Request) map[string]interface{} {
	t.Helper()
	defer func() { _ = r.Body.Close() }()
	var payload map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		t.Fatalf("decode request body: %v", err)
	}
	return payload
}

func TestNewClient(t *testing.T) {
	t.Run("uses defaults", func(t *testing.T) {
		client, err := NewClient(testAPIKey)
		if err != nil {
			t.Fatalf("NewClient() error = %v", err)
		}
		if client.apiKey != testAPIKey || client.baseURL != DefaultBaseURL {
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
			key := testAPIKey
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
		client, err := NewClient(testAPIKey, WithHTTPClient(nil), WithTimeout(time.Second))
		if err != nil {
			t.Fatalf("NewClient() error = %v", err)
		}
		if client.httpClient.Timeout != time.Second {
			t.Fatalf("timeout = %v, want %v", client.httpClient.Timeout, time.Second)
		}
	})

	t.Run("accepts an OAuth bearer token", func(t *testing.T) {
		client, err := NewClient("", WithBearerToken("access-token"))
		if err != nil || client.bearerToken != "access-token" {
			t.Fatalf("NewClient() = (%#v, %v)", client, err)
		}
	})

	t.Run("rejects both credential types", func(t *testing.T) {
		client, err := NewClient(testAPIKey, WithBearerToken("access-token"))
		if err == nil || client != nil {
			t.Fatalf("NewClient() = (%#v, %v), want error", client, err)
		}
	})
}

func TestCreateTaskBuildsV2Request(t *testing.T) {
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/v2/task.create" {
			t.Fatalf("request = %s %s", r.Method, r.URL.Path)
		}
		if r.Header.Get("x-manus-api-key") != testAPIKey {
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

func TestBearerTokenAuthentication(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer access-token" || r.Header.Get("x-manus-api-key") != "" {
			t.Fatalf("unexpected authentication headers: %#v", r.Header)
		}
		_, _ = io.WriteString(w, `{"ok":true,"task":{"id":"task_123"}}`)
	}))
	defer server.Close()
	client, err := NewClient("", WithBearerToken("access-token"), WithBaseURL(server.URL))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := client.GetTask("task_123"); err != nil {
		t.Fatalf("GetTask() error = %v", err)
	}
}

func TestCreateTaskValidatesInput(t *testing.T) {
	client, err := NewClient(testAPIKey)
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

func TestTaskPayloadOptions(t *testing.T) {
	hideInTaskList := true
	enableAskUser := false
	payload, err := newTaskPayload("Build a report", &TaskOptions{
		AgentProfile:    AgentProfileManus16,
		Locale:          "ru-RU",
		HideInTaskList:  &hideInTaskList,
		ShareVisibility: "private",
		Title:           "Report",
		ProjectID:       "project_123",
		EnableAskUser:   &enableAskUser,
		Connectors:      []string{"github"},
		EnableSkills:    []string{"research"},
		ForceSkills:     []string{"browser"},
		Attachments:     []interface{}{NewAttachmentFromURL("https://example.com/report.pdf")},
	})
	if err != nil {
		t.Fatalf("newTaskPayload() error = %v", err)
	}
	if payload["locale"] != "ru-RU" || payload["hide_in_task_list"] != true || payload["interactive_mode"] != false {
		t.Fatalf("unexpected task options payload: %#v", payload)
	}
	message := payload[messageKey].(map[string]interface{})
	if len(message[messageContentKey].([]map[string]interface{})) != 2 || len(message["connectors"].([]string)) != 1 {
		t.Fatalf("unexpected task message: %#v", message)
	}
}

func TestClientRejectsInvalidInput(t *testing.T) {
	client, err := NewClient(testAPIKey)
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name string
		call func() error
	}{
		{"empty task ID", func() error { _, err := client.GetTask(" "); return err }},
		{"nil task update", func() error { _, err := client.UpdateTask("task_123", nil); return err }},
		{"empty task update", func() error { _, err := client.UpdateTask("task_123", &TaskUpdate{}); return err }},
		{"empty task to update", func() error { _, err := client.UpdateTask("", &TaskUpdate{}); return err }},
		{"empty task to delete", func() error { _, err := client.DeleteTask(""); return err }},
		{"empty filename", func() error { _, err := client.CreateFile(""); return err }},
		{"empty upload URL", func() error { return client.UploadFileContent("", nil, "") }},
		{"empty file ID", func() error { _, err := client.GetFile(""); return err }},
		{"empty file to delete", func() error { _, err := client.DeleteFile(""); return err }},
		{"nil webhook", func() error { _, err := client.CreateWebhook(nil); return err }},
		{"empty webhook URL", func() error { _, err := client.CreateWebhook(&WebhookConfig{}); return err }},
		{"empty webhook ID", func() error { return client.DeleteWebhook("") }},
		{"empty messages task", func() error { _, err := client.ListMessages("", 0, "", "", false); return err }},
		{"empty message task", func() error { _, err := client.SendMessage("", "message", nil); return err }},
		{"empty message", func() error { _, err := client.SendMessage("task_123", "", nil); return err }},
		{"empty stopped task", func() error { _, err := client.StopTask(""); return err }},
		{"empty confirmed task", func() error { _, err := client.ConfirmAction("", "event_123", nil); return err }},
		{"empty event ID", func() error { _, err := client.ConfirmAction("task_123", "", nil); return err }},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.call(); err == nil {
				t.Fatal("expected validation error")
			}
		})
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
		contentType := r.Header.Get("Content-Type")
		if r.Method != http.MethodPut || (contentType != "text/plain" && contentType != defaultContentType) {
			t.Fatalf("unexpected upload request: %s, %s", r.Method, r.Header.Get("Content-Type"))
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()
	client, _ := NewClient(testAPIKey)
	if err := client.UploadFileContent(server.URL, []byte("content"), "text/plain"); err != nil {
		t.Fatalf("UploadFileContent() error = %v", err)
	}
	if err := client.UploadFileContent(server.URL, []byte("content"), ""); err != nil {
		t.Fatalf("UploadFileContent() with default type error = %v", err)
	}

	failingUpload := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadGateway)
		_, _ = io.WriteString(w, "upstream failure")
	})
	if err := failingUpload.UploadFileContent(failingUpload.baseURL, nil, "text/plain"); err == nil || !strings.Contains(err.Error(), "502") {
		t.Fatalf("UploadFileContent() error = %v, want HTTP status error", err)
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

func TestRequestAndResponseFailurePaths(t *testing.T) {
	transportErrorClient, err := NewClient(testAPIKey, WithHTTPClient(&http.Client{
		Transport: roundTripperFunc(func(*http.Request) (*http.Response, error) {
			return nil, errors.New("network unavailable")
		}),
	}))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := transportErrorClient.GetTask("task_123"); err == nil {
		t.Fatal("GetTask() did not return a transport error")
	}
	if err := transportErrorClient.UploadFileContent("https://example.com/upload", nil, "text/plain"); err == nil {
		t.Fatal("UploadFileContent() did not return a transport error")
	}

	if err := transportErrorClient.request(http.MethodPost, "/v2/task.create", make(chan int), nil, nil); err == nil {
		t.Fatal("request() accepted an unsupported JSON value")
	}

	invalidURLClient, _ := NewClient(testAPIKey)
	invalidURLClient.baseURL = "://invalid"
	if err := invalidURLClient.request(http.MethodGet, "/v2/task.detail", nil, nil, nil); err == nil {
		t.Fatal("request() accepted an invalid URL")
	}

	malformedResponseClient := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		_, _ = io.WriteString(w, "not-json")
	})
	if _, err := malformedResponseClient.GetTask("task_123"); err == nil {
		t.Fatal("GetTask() accepted malformed JSON")
	}

	if _, err := readResponseBody(failingReader{}); err == nil {
		t.Fatal("readResponseBody() did not report a reader error")
	}
	tooLarge := bytes.Repeat([]byte("x"), maxResponseBodySize+1)
	if _, err := readResponseBody(bytes.NewReader(tooLarge)); err == nil {
		t.Fatal("readResponseBody() accepted an oversized body")
	}
}

func TestAdditionalOpenAPIResources(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		method   string
		response string
		call     func(*Client) error
	}{
		{"create project", "/v2/project.create", http.MethodPost, `{"ok":true,"project":{"id":"project_123","name":"Reports"}}`, func(c *Client) error { _, err := c.CreateProject("Reports", "Use Russian"); return err }},
		{"list projects", "/v2/project.list", http.MethodGet, `{"ok":true,"data":[]}`, func(c *Client) error { _, err := c.ListProjects(); return err }},
		{"list skills", "/v2/skill.list", http.MethodGet, `{"ok":true,"data":[]}`, func(c *Client) error { _, err := c.ListSkills("project_123"); return err }},
		{"list agents", "/v2/agent.list", http.MethodGet, `{"ok":true,"data":[]}`, func(c *Client) error { _, err := c.ListAgents(); return err }},
		{"get agent", "/v2/agent.detail", http.MethodGet, `{"ok":true,"agent":{"id":"agent_123"}}`, func(c *Client) error { _, err := c.GetAgent("agent_123"); return err }},
		{"update agent", "/v2/agent.update", http.MethodPost, `{"ok":true,"agent":{"id":"agent_123"}}`, func(c *Client) error { _, err := c.UpdateAgent("agent_123", "Researcher", "Finds sources"); return err }},
		{"list webhooks", "/v2/webhook.list", http.MethodGet, `{"ok":true,"data":[]}`, func(c *Client) error { _, err := c.ListWebhooks(); return err }},
		{"list browser clients", "/v2/browser.onlineList", http.MethodGet, `{"ok":true,"data":[]}`, func(c *Client) error { _, err := c.ListOnlineBrowserClients(); return err }},
		{"get webhook public key", "/v2/webhook.publicKey", http.MethodGet, `{"ok":true,"public_key":"key","algorithm":"RSA"}`, func(c *Client) error { _, err := c.GetWebhookPublicKey(); return err }},
		{"list connectors", "/v2/connector.list", http.MethodGet, `{"ok":true,"data":[]}`, func(c *Client) error { _, err := c.ListConnectors(); return err }},
		{"list usage", "/v2/usage.list", http.MethodGet, `{"ok":true,"data":[]}`, func(c *Client) error { _, err := c.ListUsage(20, "next"); return err }},
		{"team usage statistic", "/v2/usage.teamStatistic", http.MethodGet, `{"ok":true,"data":[]}`, func(c *Client) error { _, err := c.GetTeamUsageStatistic("2026-01-01", "2026-01-31"); return err }},
		{"team usage log", "/v2/usage.teamLog", http.MethodGet, `{"ok":true,"data":[]}`, func(c *Client) error {
			ascending := true
			_, err := c.ListTeamUsageLog(20, "next", "2026-01-01", "2026-01-31", "credits", &ascending)
			return err
		}},
		{"available credits", "/v2/usage.availableCredits", http.MethodGet, `{"ok":true,"data":{}}`, func(c *Client) error { _, err := c.GetAvailableCredits(); return err }},
		{"website status", "/v2/website.status", http.MethodGet, `{"ok":true,"website_id":"website_123"}`, func(c *Client) error { _, err := c.GetWebsiteStatus("task_123", "website_123"); return err }},
		{"website checkpoints", "/v2/website.listCheckpoints", http.MethodGet, `{"ok":true,"website_id":"website_123","data":[]}`, func(c *Client) error { _, err := c.ListWebsiteCheckpoints("task_123", "website_123"); return err }},
		{"publish website", "/v2/website.publish", http.MethodPost, `{"ok":true,"website_id":"website_123"}`, func(c *Client) error { _, err := c.PublishWebsite("task_123", "website_123", "public"); return err }},
		{"update website", "/v2/website.update", http.MethodPost, `{"ok":true}`, func(c *Client) error { return c.UpdateWebsite("task_123", "website_123", "Title", "") }},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
				if r.Method != tt.method || r.URL.Path != tt.path {
					t.Fatalf("request = %s %s, want %s %s", r.Method, r.URL.Path, tt.method, tt.path)
				}
				if tt.name == "list skills" && r.URL.Query().Get("project_id") != "project_123" {
					t.Fatalf("project_id = %q", r.URL.Query().Get("project_id"))
				}
				_, _ = io.WriteString(w, tt.response)
			})
			if err := tt.call(client); err != nil {
				t.Fatalf("request error = %v", err)
			}
		})
	}

	client, _ := NewClient(testAPIKey)
	if _, err := client.CreateProject("", ""); err == nil {
		t.Fatal("CreateProject() accepted an empty name")
	}
	if _, err := client.GetAgent(""); err == nil {
		t.Fatal("GetAgent() accepted an empty agent ID")
	}
	if _, err := client.UpdateAgent("agent_123", "", ""); err == nil {
		t.Fatal("UpdateAgent() accepted an empty update")
	}
}
