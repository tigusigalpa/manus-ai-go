package manusai

import (
	"encoding/json"
	"fmt"
)

// ParseWebhookPayload decodes and validates a Manus webhook request body.
func ParseWebhookPayload(jsonPayload []byte) (*WebhookPayload, error) {
	var payload WebhookPayload
	if err := json.Unmarshal(jsonPayload, &payload); err != nil {
		return nil, fmt.Errorf("invalid JSON payload: %w", err)
	}

	if payload.EventType == "" {
		return nil, fmt.Errorf("missing event_type in webhook payload")
	}

	return &payload, nil
}

// IsTaskCreated reports whether payload represents a task_created event.
func IsTaskCreated(payload *WebhookPayload) bool {
	return payload != nil && payload.EventType == "task_created"
}

// IsTaskStopped reports whether payload represents a task_stopped event.
func IsTaskStopped(payload *WebhookPayload) bool {
	return payload != nil && payload.EventType == "task_stopped"
}

// IsTaskCompleted reports whether a stopped task finished successfully.
func IsTaskCompleted(payload *WebhookPayload) bool {
	if !IsTaskStopped(payload) {
		return false
	}

	if payload.TaskDetail == nil {
		return false
	}

	stopReason, ok := payload.TaskDetail["stop_reason"].(string)
	return ok && stopReason == "finish"
}

// IsTaskAskingForInput reports whether a stopped task is waiting for user input.
func IsTaskAskingForInput(payload *WebhookPayload) bool {
	if !IsTaskStopped(payload) {
		return false
	}

	if payload.TaskDetail == nil {
		return false
	}

	stopReason, ok := payload.TaskDetail["stop_reason"].(string)
	return ok && stopReason == "ask"
}

// GetTaskDetail returns the event task detail, or nil when payload is nil.
func GetTaskDetail(payload *WebhookPayload) map[string]interface{} {
	if payload == nil {
		return nil
	}
	return payload.TaskDetail
}

// GetAttachments returns attachments from the event task detail, or nil when absent.
func GetAttachments(payload *WebhookPayload) []interface{} {
	if payload == nil || payload.TaskDetail == nil {
		return nil
	}

	attachments, ok := payload.TaskDetail["attachments"].([]interface{})
	if !ok {
		return nil
	}

	return attachments
}
