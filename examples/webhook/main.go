package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	manusai "github.com/tigusigalpa/manus-ai-go"
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

	fmt.Println("=== Creating Webhook ===")
	webhook := &manusai.WebhookConfig{
		URL:    "https://your-domain.com/webhook/manus-ai",
		Events: []string{"task_created", "task_stopped"},
	}

	webhookResult, err := client.CreateWebhook(webhook)
	if err != nil {
		log.Fatalf("Failed to create webhook: %v", err)
	}

	fmt.Printf("Webhook created successfully!\n")
	fmt.Printf("Webhook ID: %s\n", webhookResult.WebhookID)

	fmt.Println("\n=== Starting Webhook Server ===")
	fmt.Println("Server will listen on http://localhost:8080/webhook")
	fmt.Println("Press Ctrl+C to stop")

	http.HandleFunc("/webhook", handleWebhook)
	
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

func handleWebhook(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading request body: %v", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	payload, err := manusai.ParseWebhookPayload(body)
	if err != nil {
		log.Printf("Error parsing webhook payload: %v", err)
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}

	fmt.Printf("\n=== Webhook Event Received ===\n")
	fmt.Printf("Event Type: %s\n", payload.EventType)

	if manusai.IsTaskCreated(payload) {
		fmt.Println("Event: Task Created")
		taskDetail := manusai.GetTaskDetail(payload)
		if taskDetail != nil {
			fmt.Printf("Task Details: %+v\n", taskDetail)
		}
	}

	if manusai.IsTaskStopped(payload) {
		fmt.Println("Event: Task Stopped")
		
		if manusai.IsTaskCompleted(payload) {
			fmt.Println("Task completed successfully!")
			taskDetail := manusai.GetTaskDetail(payload)
			if taskDetail != nil {
				if taskID, ok := taskDetail["task_id"].(string); ok {
					fmt.Printf("Task ID: %s\n", taskID)
				}
				if message, ok := taskDetail["message"].(string); ok {
					fmt.Printf("Message: %s\n", message)
				}
			}

			attachments := manusai.GetAttachments(payload)
			if len(attachments) > 0 {
				fmt.Printf("Attachments (%d):\n", len(attachments))
				for i, att := range attachments {
					if attMap, ok := att.(map[string]interface{}); ok {
						fmt.Printf("  %d. ", i+1)
						if fileName, ok := attMap["file_name"].(string); ok {
							fmt.Printf("File: %s ", fileName)
						}
						if size, ok := attMap["size_bytes"].(float64); ok {
							fmt.Printf("(%d bytes) ", int64(size))
						}
						if url, ok := attMap["url"].(string); ok {
							fmt.Printf("URL: %s", url)
						}
						fmt.Println()
					}
				}
			}
		}

		if manusai.IsTaskAskingForInput(payload) {
			fmt.Println("Task is asking for user input!")
			taskDetail := manusai.GetTaskDetail(payload)
			if taskDetail != nil {
				if message, ok := taskDetail["message"].(string); ok {
					fmt.Printf("Input required: %s\n", message)
				}
			}
		}
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
