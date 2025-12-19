package manusai

import (
	"encoding/json"
	"fmt"
)

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

func IsTaskCreated(payload *WebhookPayload) bool {
	return payload.EventType == "task_created"
}

func IsTaskStopped(payload *WebhookPayload) bool {
	return payload.EventType == "task_stopped"
}

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

func GetTaskDetail(payload *WebhookPayload) map[string]interface{} {
	return payload.TaskDetail
}

func GetAttachments(payload *WebhookPayload) []interface{} {
	if payload.TaskDetail == nil {
		return nil
	}

	attachments, ok := payload.TaskDetail["attachments"].([]interface{})
	if !ok {
		return nil
	}

	return attachments
}
