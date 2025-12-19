package manusai

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseWebhookPayload(t *testing.T) {
	t.Run("valid payload", func(t *testing.T) {
		jsonPayload := []byte(`{"event_type":"task_created","task_detail":{"task_id":"123"}}`)
		payload, err := ParseWebhookPayload(jsonPayload)
		require.NoError(t, err)
		assert.Equal(t, "task_created", payload.EventType)
		assert.NotNil(t, payload.TaskDetail)
	})

	t.Run("invalid JSON", func(t *testing.T) {
		jsonPayload := []byte(`{invalid json}`)
		payload, err := ParseWebhookPayload(jsonPayload)
		assert.Error(t, err)
		assert.Nil(t, payload)
	})

	t.Run("missing event_type", func(t *testing.T) {
		jsonPayload := []byte(`{"task_detail":{"task_id":"123"}}`)
		payload, err := ParseWebhookPayload(jsonPayload)
		assert.Error(t, err)
		assert.Nil(t, payload)
	})
}

func TestIsTaskCreated(t *testing.T) {
	payload := &WebhookPayload{EventType: "task_created"}
	assert.True(t, IsTaskCreated(payload))

	payload = &WebhookPayload{EventType: "task_stopped"}
	assert.False(t, IsTaskCreated(payload))
}

func TestIsTaskStopped(t *testing.T) {
	payload := &WebhookPayload{EventType: "task_stopped"}
	assert.True(t, IsTaskStopped(payload))

	payload = &WebhookPayload{EventType: "task_created"}
	assert.False(t, IsTaskStopped(payload))
}

func TestIsTaskCompleted(t *testing.T) {
	t.Run("completed task", func(t *testing.T) {
		payload := &WebhookPayload{
			EventType: "task_stopped",
			TaskDetail: map[string]interface{}{
				"stop_reason": "finish",
			},
		}
		assert.True(t, IsTaskCompleted(payload))
	})

	t.Run("task asking for input", func(t *testing.T) {
		payload := &WebhookPayload{
			EventType: "task_stopped",
			TaskDetail: map[string]interface{}{
				"stop_reason": "ask",
			},
		}
		assert.False(t, IsTaskCompleted(payload))
	})

	t.Run("task created event", func(t *testing.T) {
		payload := &WebhookPayload{EventType: "task_created"}
		assert.False(t, IsTaskCompleted(payload))
	})
}

func TestIsTaskAskingForInput(t *testing.T) {
	t.Run("task asking for input", func(t *testing.T) {
		payload := &WebhookPayload{
			EventType: "task_stopped",
			TaskDetail: map[string]interface{}{
				"stop_reason": "ask",
			},
		}
		assert.True(t, IsTaskAskingForInput(payload))
	})

	t.Run("completed task", func(t *testing.T) {
		payload := &WebhookPayload{
			EventType: "task_stopped",
			TaskDetail: map[string]interface{}{
				"stop_reason": "finish",
			},
		}
		assert.False(t, IsTaskAskingForInput(payload))
	})
}

func TestGetTaskDetail(t *testing.T) {
	detail := map[string]interface{}{
		"task_id": "123",
		"message": "Test message",
	}
	payload := &WebhookPayload{
		EventType:  "task_stopped",
		TaskDetail: detail,
	}

	result := GetTaskDetail(payload)
	assert.Equal(t, detail, result)
}

func TestGetAttachments(t *testing.T) {
	t.Run("with attachments", func(t *testing.T) {
		attachments := []interface{}{
			map[string]interface{}{"file_name": "test.pdf"},
		}
		payload := &WebhookPayload{
			EventType: "task_stopped",
			TaskDetail: map[string]interface{}{
				"attachments": attachments,
			},
		}

		result := GetAttachments(payload)
		assert.Equal(t, attachments, result)
	})

	t.Run("without attachments", func(t *testing.T) {
		payload := &WebhookPayload{
			EventType:  "task_stopped",
			TaskDetail: map[string]interface{}{},
		}

		result := GetAttachments(payload)
		assert.Nil(t, result)
	})
}
