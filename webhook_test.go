package manusai

import "testing"

func TestParseWebhookPayload(t *testing.T) {
	payload, err := ParseWebhookPayload([]byte(`{"event_type":"task_stopped","task_detail":{"stop_reason":"finish"}}`))
	if err != nil || !IsTaskStopped(payload) || !IsTaskCompleted(payload) || IsTaskAskingForInput(payload) {
		t.Fatalf("unexpected payload result: %#v, %v", payload, err)
	}

	for _, input := range [][]byte{[]byte(`{invalid}`), []byte(`{"task_detail":{}}`)} {
		if payload, err := ParseWebhookPayload(input); err == nil || payload != nil {
			t.Fatalf("ParseWebhookPayload(%s) = %#v, %v; want error", input, payload, err)
		}
	}
}

func TestWebhookHelpersHandleNilAndAttachments(t *testing.T) {
	if IsTaskCreated(nil) || IsTaskStopped(nil) || IsTaskCompleted(nil) || IsTaskAskingForInput(nil) || GetTaskDetail(nil) != nil || GetAttachments(nil) != nil {
		t.Fatal("webhook helpers must handle nil payloads")
	}

	attachments := []interface{}{map[string]interface{}{"file_name": "report.pdf"}}
	payload := &WebhookPayload{EventType: "task_stopped", TaskDetail: map[string]interface{}{
		"stop_reason": "ask",
		"attachments": attachments,
	}}
	if !IsTaskAskingForInput(payload) || IsTaskCompleted(payload) || len(GetAttachments(payload)) != 1 {
		t.Fatalf("unexpected webhook helpers result: %#v", payload)
	}
}
