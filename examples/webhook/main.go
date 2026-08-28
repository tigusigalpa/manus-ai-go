package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	manusai "github.com/tigusigalpa/manus-ai-go/v2"
)

func main() {
	apiKey := os.Getenv("MANUS_AI_API_KEY")
	if apiKey == "" {
		log.Fatal("MANUS_AI_API_KEY environment variable is required")
	}

	client, err := manusai.NewClient(apiKey)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	webhook, err := client.CreateWebhook(&manusai.WebhookConfig{
		URL:    "https://your-domain.com/webhook/manus-ai",
		Events: []string{"task_created", "task_stopped"},
	})
	if err != nil {
		log.Fatalf("Failed to create webhook: %v", err)
	}

	fmt.Printf("Webhook created: %s\n", webhook.WebhookID)
	http.HandleFunc("/webhook", handleWebhook)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleWebhook(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	defer func() { _ = r.Body.Close() }()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "could not read request", http.StatusBadRequest)
		return
	}
	payload, err := manusai.ParseWebhookPayload(body)
	if err != nil {
		http.Error(w, "invalid payload", http.StatusBadRequest)
		return
	}

	logWebhookEvent(payload)
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("OK"))
}

func logWebhookEvent(payload *manusai.WebhookPayload) {
	if manusai.IsTaskCreated(payload) {
		fmt.Printf("Task created: %+v\n", manusai.GetTaskDetail(payload))
		return
	}
	if !manusai.IsTaskStopped(payload) {
		return
	}

	fmt.Println("Task stopped")
	if manusai.IsTaskCompleted(payload) {
		logCompletedTask(payload)
	}
	if manusai.IsTaskAskingForInput(payload) {
		logInputRequest(payload)
	}
}

func logCompletedTask(payload *manusai.WebhookPayload) {
	detail := manusai.GetTaskDetail(payload)
	if taskID, ok := detail["task_id"].(string); ok {
		fmt.Printf("Task completed: %s\n", taskID)
	}
	if message, ok := detail["message"].(string); ok {
		fmt.Printf("Message: %s\n", message)
	}
	logAttachments(manusai.GetAttachments(payload))
}

func logAttachments(attachments []interface{}) {
	for _, attachment := range attachments {
		attachmentMap, ok := attachment.(map[string]interface{})
		if !ok {
			continue
		}
		fmt.Printf("Attachment: %v\n", attachmentMap)
	}
}

func logInputRequest(payload *manusai.WebhookPayload) {
	if message, ok := manusai.GetTaskDetail(payload)["message"].(string); ok {
		fmt.Printf("Input required: %s\n", message)
	}
}
