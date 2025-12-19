/*
Package manusai provides a Go client library for the Manus AI API.

Manus AI is an AI agent platform that allows you to create and manage AI-powered tasks,
upload files, and receive real-time notifications via webhooks.

# Installation

	go get github.com/tigusigalpa/manus-ai-go

# Quick Start

	import manusai "github.com/tigusigalpa/manus-ai-go"

	// Create a client
	client, err := manusai.NewClient("your-api-key")
	if err != nil {
		log.Fatal(err)
	}

	// Create a task
	task, err := client.CreateTask("Write a poem about Go", &manusai.TaskOptions{
		AgentProfile: manusai.AgentProfileManus16,
		TaskMode:     "chat",
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Task created: %s\n", task.TaskID)

# Features

- Task creation and management
- File upload and attachment handling
- Webhook integration for real-time updates
- Comprehensive error handling
- Type-safe interfaces
- Full test coverage

# API Documentation

For detailed API documentation, visit: https://open.manus.ai/docs

# Examples

See the examples/ directory for complete working examples:
  - examples/basic/ - Basic task creation and management
  - examples/file-upload/ - File upload with attachments
  - examples/webhook/ - Webhook setup and handling
*/
package manusai
