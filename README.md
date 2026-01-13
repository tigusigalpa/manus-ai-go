# Manus AI Go SDK

![Manus AI Golang SDK](https://github.com/user-attachments/assets/1249e90c-a860-4f86-9a77-2d048f94854d)

🚀 Complete Golang library for integration with [Manus AI](https://manus.ai) API. Easily integrate Manus AI agent into
your Go applications.

**📦 Package:** [pkg.go.dev/github.com/tigusigalpa/manus-ai-go](https://pkg.go.dev/github.com/tigusigalpa/manus-ai-go)

[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/Go-%3E%3D1.21-blue)](https://golang.org/)
[![Go Report Card](https://goreportcard.com/badge/github.com/tigusigalpa/manus-ai-go)](https://goreportcard.com/report/github.com/tigusigalpa/manus-ai-go)

**🌐 Language:** English | [Русский](README-ru.md)

## 📋 Table of Contents

- [Features](#features)
- [Requirements](#requirements)
- [Installation](#installation)
- [Configuration](#configuration)
- [Usage](#usage)
    - [Basic Usage](#basic-usage)
    - [Task Management](#task-management)
    - [File Management](#file-management)
    - [Webhooks](#webhooks)
- [API Reference](#api-reference)
- [Examples](#examples)
- [Testing](#testing)
- [Contributing](#contributing)
- [License](#license)

## ✨ Features

- ✅ Full support for Manus AI API
- ✅ Task creation and management
- ✅ File upload and attachment handling
- ✅ Webhook integration for real-time updates
- ✅ Comprehensive error handling with custom error types
- ✅ Type-safe interfaces
- ✅ Full test coverage
- ✅ Detailed documentation and examples
- ✅ Idiomatic Go code

## 📦 Requirements

- Go 1.21 or higher

## 🔧 Installation

Install the package using `go get`:

```bash
go get github.com/tigusigalpa/manus-ai-go
```

## ⚙️ Configuration

### Getting Your API Key

1. Sign up at [Manus AI](https://manus.im)
2. Get your API key from the [API Integration settings](http://manus.im/app?show_settings=integrations&app_name=api)

### Basic Configuration

```go
import manusai "github.com/tigusigalpa/manus-ai-go"

client, err := manusai.NewClient("your-api-key-here")
if err != nil {
    log.Fatal(err)
}
```

### Custom Configuration

```go
import (
    "time"
    manusai "github.com/tigusigalpa/manus-ai-go"
)

client, err := manusai.NewClient(
    "your-api-key",
    manusai.WithBaseURL("https://custom.api.com"),
    manusai.WithTimeout(60 * time.Second),
)
```

## 🚀 Usage

### Basic Usage

```go
package main

import (
    "fmt"
    "log"
    
    manusai "github.com/tigusigalpa/manus-ai-go"
)

func main() {
    client, err := manusai.NewClient("your-api-key")
    if err != nil {
        log.Fatal(err)
    }

    // Create a task
    task, err := client.CreateTask("Write a poem about Go programming", &manusai.TaskOptions{
        AgentProfile: manusai.AgentProfileManus16,
        TaskMode:     "chat",
    })
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Task created: %s\n", task.TaskID)
    fmt.Printf("View at: %s\n", task.TaskURL)
}
```

### Task Management

Tasks are the core of Manus AI - they represent AI agent work items that can perform complex operations, answer
questions, or automate workflows.

**API Documentation:** [Tasks API Reference](https://open.manus.ai/docs/api-reference/create-task)

#### Create a Task

```go
task, err := client.CreateTask("Your task prompt here", &manusai.TaskOptions{
    AgentProfile:        manusai.AgentProfileManus16,
    TaskMode:            "agent",  // "chat", "adaptive", or "agent"
    Locale:              "en-US",
    HideInTaskList:      &falseVal,
    CreateShareableLink: &trueVal,
})
if err != nil {
    log.Fatal(err)
}
```

**Available Agent Profiles:**

- `AgentProfileManus16` - Latest and most capable model (recommended)
- `AgentProfileManus16Lite` - Faster, lightweight version
- `AgentProfileManus16Max` - Maximum capability version
- `AgentProfileSpeed` - ⚠️ Deprecated, use `AgentProfileManus16Lite` instead
- `AgentProfileQuality` - ⚠️ Deprecated, use `AgentProfileManus16` instead

```go
// Check if a profile is valid
if manusai.IsValidAgentProfile("manus-1.6") {
    fmt.Println("Valid profile")
}

// Check if deprecated
if manusai.IsDeprecatedAgentProfile(manusai.AgentProfileSpeed) {
    fmt.Println("This profile is deprecated")
}

// Get all recommended profiles
profiles := manusai.RecommendedAgentProfiles()
```

#### Get Task Details

```go
task, err := client.GetTask("task_id")
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Status: %s\n", task.Status)
fmt.Printf("Credits used: %.2f\n", task.CreditUsage)

// Access output messages
for _, message := range task.Output {
    fmt.Printf("[%s]: %s\n", message.Role, message.Content)
}
```

#### List Tasks

```go
tasks, err := client.GetTasks(&manusai.TaskFilters{
    Limit:   10,
    Order:   "desc",
    OrderBy: "created_at",
    Status:  []string{"completed", "running"},
})
if err != nil {
    log.Fatal(err)
}

for _, task := range tasks.Data {
    fmt.Printf("Task %s: %s\n", task.ID, task.Status)
}
```

#### Update Task

```go
newTitle := "New Task Title"
enableShared := true

updated, err := client.UpdateTask("task_id", &manusai.TaskUpdate{
    Title:                   &newTitle,
    EnableShared:            &enableShared,
    EnableVisibleInTaskList: &enableShared,
})
if err != nil {
    log.Fatal(err)
}
```

#### Delete Task

```go
result, err := client.DeleteTask("task_id")
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Deleted: %v\n", result.Deleted)
```

### File Management

Manus AI supports file attachments to provide context for your tasks. The file upload process uses a two-step approach:
first, create a file record to get a secure presigned URL, then upload your content directly to cloud storage.

**API Documentation:** [Files API Reference](https://open.manus.ai/docs/api-reference/create-file)

#### Upload a File

```go
// 1. Create file record
fileResult, err := client.CreateFile("document.pdf")
if err != nil {
    log.Fatal(err)
}

// 2. Upload file content
fileContent, _ := os.ReadFile("/path/to/document.pdf")
err = client.UploadFileContent(
    fileResult.UploadURL,
    fileContent,
    "application/pdf",
)
if err != nil {
    log.Fatal(err)
}

// 3. Use file in task
attachment := manusai.NewAttachmentFromFileID(fileResult.ID)

task, err := client.CreateTask("Analyze this document", &manusai.TaskOptions{
    Attachments: []interface{}{attachment},
})
```

#### Different Attachment Types

```go
// From file ID
attachment1 := manusai.NewAttachmentFromFileID("file_123")

// From URL
attachment2 := manusai.NewAttachmentFromURL("https://example.com/image.jpg")

// From base64
attachment3 := manusai.NewAttachmentFromBase64(base64Data, "image/png")

// From local file path
attachment4, err := manusai.NewAttachmentFromFilePath("/path/to/file.pdf")
if err != nil {
    log.Fatal(err)
}
```

#### List Files

```go
files, err := client.ListFiles()
if err != nil {
    log.Fatal(err)
}

for _, file := range files.Data {
    fmt.Printf("%s - %s\n", file.Filename, file.Status)
}
```

#### Delete File

```go
result, err := client.DeleteFile("file_id")
if err != nil {
    log.Fatal(err)
}
```

### Webhooks

Webhooks enable real-time notifications about your task lifecycle events. Instead of polling for updates, Manus AI will
send HTTP POST requests to your specified endpoint whenever important events occur.

**API Documentation:** [Webhooks Guide](https://open.manus.ai/docs/webhooks/index)

#### Create Webhook

```go
webhook := &manusai.WebhookConfig{
    URL:    "https://your-domain.com/webhook/manus-ai",
    Events: []string{"task_created", "task_stopped"},
}

result, err := client.CreateWebhook(webhook)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Webhook ID: %s\n", result.WebhookID)
```

#### Handle Webhook Events

```go
import (
    "io"
    "net/http"
    
    manusai "github.com/tigusigalpa/manus-ai-go"
)

func handleWebhook(w http.ResponseWriter, r *http.Request) {
    body, _ := io.ReadAll(r.Body)
    defer r.Body.Close()

    payload, err := manusai.ParseWebhookPayload(body)
    if err != nil {
        http.Error(w, "Invalid payload", http.StatusBadRequest)
        return
    }

    if manusai.IsTaskCompleted(payload) {
        taskDetail := manusai.GetTaskDetail(payload)
        attachments := manusai.GetAttachments(payload)
        
        fmt.Printf("Task completed: %v\n", taskDetail["task_id"])
        fmt.Printf("Message: %v\n", taskDetail["message"])
        
        // Download attachments
        for _, att := range attachments {
            if attMap, ok := att.(map[string]interface{}); ok {
                fmt.Printf("File: %v\n", attMap["file_name"])
                fmt.Printf("URL: %v\n", attMap["url"])
            }
        }
    }

    if manusai.IsTaskAskingForInput(payload) {
        taskDetail := manusai.GetTaskDetail(payload)
        fmt.Printf("Input required: %v\n", taskDetail["message"])
    }

    w.WriteHeader(http.StatusOK)
}
```

#### Delete Webhook

```go
err := client.DeleteWebhook("webhook_id")
if err != nil {
    log.Fatal(err)
}
```

## 📚 API Reference

### Client Methods

#### Task Methods

- `CreateTask(prompt string, options *TaskOptions) (*TaskResponse, error)`
- `GetTasks(filters *TaskFilters) (*TaskListResponse, error)`
- `GetTask(taskID string) (*TaskDetail, error)`
- `UpdateTask(taskID string, updates *TaskUpdate) (*TaskDetail, error)`
- `DeleteTask(taskID string) (*DeleteResponse, error)`

#### File Methods

- `CreateFile(filename string) (*FileResponse, error)`
- `UploadFileContent(uploadURL string, fileContent []byte, contentType string) error`
- `ListFiles() (*FileListResponse, error)`
- `GetFile(fileID string) (*FileDetail, error)`
- `DeleteFile(fileID string) (*DeleteResponse, error)`

#### Webhook Methods

- `CreateWebhook(webhook *WebhookConfig) (*WebhookResponse, error)`
- `DeleteWebhook(webhookID string) error`

### Helper Functions

#### Agent Profile

- `AllAgentProfiles() []string` - Get all available profiles
- `RecommendedAgentProfiles() []string` - Get recommended profiles
- `IsValidAgentProfile(profile string) bool` - Check if profile is valid
- `IsDeprecatedAgentProfile(profile string) bool` - Check if profile is deprecated

#### Attachments

- `NewAttachmentFromFileID(fileID string) map[string]interface{}`
- `NewAttachmentFromURL(url string) map[string]interface{}`
- `NewAttachmentFromBase64(base64Data, mimeType string) map[string]interface{}`
- `NewAttachmentFromFilePath(filePath string) (map[string]interface{}, error)`

#### Webhook Handlers

- `ParseWebhookPayload(jsonPayload []byte) (*WebhookPayload, error)`
- `IsTaskCreated(payload *WebhookPayload) bool`
- `IsTaskStopped(payload *WebhookPayload) bool`
- `IsTaskCompleted(payload *WebhookPayload) bool`
- `IsTaskAskingForInput(payload *WebhookPayload) bool`
- `GetTaskDetail(payload *WebhookPayload) map[string]interface{}`
- `GetAttachments(payload *WebhookPayload) []interface{}`

### Error Types

The SDK provides custom error types for better error handling:

- `ManusAIError` - General API errors
- `AuthenticationError` - Authentication/authorization failures
- `ValidationError` - Request validation errors

```go
_, err := client.GetTask("invalid_id")
if err != nil {
    switch e := err.(type) {
    case *manusai.AuthenticationError:
        fmt.Println("Authentication failed:", e.Message)
    case *manusai.ValidationError:
        fmt.Println("Validation error:", e.Message)
    case *manusai.ManusAIError:
        fmt.Println("API error:", e.Message)
    default:
        fmt.Println("Unknown error:", err)
    }
}
```

## 💡 Examples

See the `examples/` directory for complete working examples:

- `examples/basic/` - Basic task creation and management
- `examples/file-upload/` - File upload with attachments
- `examples/webhook/` - Webhook setup and handling

To run an example:

```bash
export MANUS_AI_API_KEY=your-api-key
cd examples/basic
go run main.go
```

## 🧪 Testing

Run the test suite:

```bash
go test -v ./...
```

Run with coverage:

```bash
go test -v -cover ./...
```

Generate coverage report:

```bash
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🔗 Links

- [Manus AI Official Website](https://manus.ai)
- [API Documentation](https://open.manus.ai/docs)
- [GitHub Repository](https://github.com/tigusigalpa/manus-ai-go)
- [Issue Tracker](https://github.com/tigusigalpa/manus-ai-go/issues)

## 👤 Author

**Igor Sazonov**

- GitHub: [@tigusigalpa](https://github.com/tigusigalpa)
- Email: sovletig@gmail.com

## 🙏 Acknowledgments

- Thanks to the Manus AI team for providing an excellent AI agent platform
- Inspired by the PHP SDK: [manus-ai-php](https://github.com/tigusigalpa/manus-ai-php)

---

Made with ❤️ by [Igor Sazonov](https://github.com/tigusigalpa)
