package manusai

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewClient(t *testing.T) {
	t.Run("valid API key", func(t *testing.T) {
		client, err := NewClient("test-api-key")
		require.NoError(t, err)
		assert.NotNil(t, client)
		assert.Equal(t, "test-api-key", client.apiKey)
		assert.Equal(t, DefaultBaseURL, client.baseURL)
	})

	t.Run("empty API key", func(t *testing.T) {
		client, err := NewClient("")
		assert.Error(t, err)
		assert.Nil(t, client)
		assert.IsType(t, &AuthenticationError{}, err)
	})

	t.Run("with custom base URL", func(t *testing.T) {
		client, err := NewClient("test-key", WithBaseURL("https://custom.api.com"))
		require.NoError(t, err)
		assert.Equal(t, "https://custom.api.com", client.baseURL)
	})
}

func TestCreateTask(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "/v1/tasks", r.URL.Path)
		assert.Equal(t, "test-api-key", r.Header.Get("Authorization"))

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"task_id":"task_123","task_title":"Test Task","task_url":"https://manus.ai/task/123"}`))
	}))
	defer server.Close()

	client, _ := NewClient("test-api-key", WithBaseURL(server.URL))

	t.Run("successful task creation", func(t *testing.T) {
		result, err := client.CreateTask("Test prompt", nil)
		require.NoError(t, err)
		assert.Equal(t, "task_123", result.TaskID)
		assert.Equal(t, "Test Task", result.TaskTitle)
		assert.Equal(t, "https://manus.ai/task/123", result.TaskURL)
	})

	t.Run("empty prompt", func(t *testing.T) {
		result, err := client.CreateTask("", nil)
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.IsType(t, &ValidationError{}, err)
	})
}

func TestGetTask(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, "/v1/tasks/task_123", r.URL.Path)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"id":"task_123","title":"Test","status":"completed","credit_usage":1.5}`))
	}))
	defer server.Close()

	client, _ := NewClient("test-api-key", WithBaseURL(server.URL))

	t.Run("successful get task", func(t *testing.T) {
		result, err := client.GetTask("task_123")
		require.NoError(t, err)
		assert.Equal(t, "task_123", result.ID)
		assert.Equal(t, "Test", result.Title)
		assert.Equal(t, "completed", result.Status)
		assert.Equal(t, 1.5, result.CreditUsage)
	})

	t.Run("empty task ID", func(t *testing.T) {
		result, err := client.GetTask("")
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.IsType(t, &ValidationError{}, err)
	})
}

func TestDeleteTask(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "DELETE", r.Method)
		assert.Equal(t, "/v1/tasks/task_123", r.URL.Path)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"deleted":true}`))
	}))
	defer server.Close()

	client, _ := NewClient("test-api-key", WithBaseURL(server.URL))

	result, err := client.DeleteTask("task_123")
	require.NoError(t, err)
	assert.True(t, result.Deleted)
}

func TestCreateFile(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "/v1/files", r.URL.Path)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"id":"file_123","filename":"test.pdf","upload_url":"https://s3.example.com/upload","status":"pending"}`))
	}))
	defer server.Close()

	client, _ := NewClient("test-api-key", WithBaseURL(server.URL))

	result, err := client.CreateFile("test.pdf")
	require.NoError(t, err)
	assert.Equal(t, "file_123", result.ID)
	assert.Equal(t, "test.pdf", result.Filename)
	assert.NotEmpty(t, result.UploadURL)
}

func TestCreateWebhook(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "/v1/webhooks", r.URL.Path)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"webhook_id":"webhook_123"}`))
	}))
	defer server.Close()

	client, _ := NewClient("test-api-key", WithBaseURL(server.URL))

	webhook := &WebhookConfig{
		URL:    "https://example.com/webhook",
		Events: []string{"task_created", "task_stopped"},
	}

	result, err := client.CreateWebhook(webhook)
	require.NoError(t, err)
	assert.Equal(t, "webhook_123", result.WebhookID)
}

func TestErrorHandling(t *testing.T) {
	t.Run("authentication error", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Invalid API key"))
		}))
		defer server.Close()

		client, _ := NewClient("invalid-key", WithBaseURL(server.URL))
		_, err := client.GetTask("task_123")
		assert.Error(t, err)
		assert.IsType(t, &AuthenticationError{}, err)
	})

	t.Run("validation error", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Invalid request"))
		}))
		defer server.Close()

		client, _ := NewClient("test-key", WithBaseURL(server.URL))
		_, err := client.GetTask("task_123")
		assert.Error(t, err)
		assert.IsType(t, &ValidationError{}, err)
	})
}
