# Manus AI Golang SDK

![Manus API Golang SDK](https://i.postimg.cc/NfsBb6Td/manus-ai-golang-hero-banner.jpg)

[![Tests](https://github.com/tigusigalpa/manus-ai-go/actions/workflows/test.yml/badge.svg?branch=main)](https://github.com/tigusigalpa/manus-ai-go/actions/workflows/test.yml)
[![CodeQL](https://github.com/tigusigalpa/manus-ai-go/actions/workflows/codeql.yml/badge.svg?branch=main)](https://github.com/tigusigalpa/manus-ai-go/actions/workflows/codeql.yml)
[![Codecov](https://codecov.io/gh/tigusigalpa/manus-ai-go/branch/main/graph/badge.svg)](https://codecov.io/gh/tigusigalpa/manus-ai-go)
[![Go Reference](https://pkg.go.dev/badge/github.com/tigusigalpa/manus-ai-go/v2.svg)](https://pkg.go.dev/github.com/tigusigalpa/manus-ai-go/v2)
[![Go Version](https://img.shields.io/github/go-mod/go-version/tigusigalpa/manus-ai-go?logo=go)](go.mod)
[![Latest Release](https://img.shields.io/github/v/release/tigusigalpa/manus-ai-go?display_name=tag&logo=github)](https://github.com/tigusigalpa/manus-ai-go/releases)
[![License](https://img.shields.io/github/license/tigusigalpa/manus-ai-go?color=blue)](LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/tigusigalpa/manus-ai-go)](https://goreportcard.com/report/github.com/tigusigalpa/manus-ai-go)
[![Last Commit](https://img.shields.io/github/last-commit/tigusigalpa/manus-ai-go?logo=git)](https://github.com/tigusigalpa/manus-ai-go/commits/main)
[![Commit Activity](https://img.shields.io/github/commit-activity/m/tigusigalpa/manus-ai-go?logo=github)](https://github.com/tigusigalpa/manus-ai-go/pulse)
[![Open Issues](https://img.shields.io/github/issues/tigusigalpa/manus-ai-go?logo=github)](https://github.com/tigusigalpa/manus-ai-go/issues)
[![Pull Requests](https://img.shields.io/github/issues-pr/tigusigalpa/manus-ai-go?logo=github)](https://github.com/tigusigalpa/manus-ai-go/pulls)
[![Stars](https://img.shields.io/github/stars/tigusigalpa/manus-ai-go?style=flat&logo=github)](https://github.com/tigusigalpa/manus-ai-go/stargazers)
[![Forks](https://img.shields.io/github/forks/tigusigalpa/manus-ai-go?style=flat&logo=github)](https://github.com/tigusigalpa/manus-ai-go/network/members)
[![Repo Size](https://img.shields.io/github/repo-size/tigusigalpa/manus-ai-go?logo=github)](https://github.com/tigusigalpa/manus-ai-go)

A small, idiomatic Go client for the Manus API v2. It helps an application create and manage Manus tasks, continue a conversation, upload files, and process webhooks without hand-writing HTTP requests.

[Русская версия](README-ru.md) · [Manus API documentation](https://open.manus.im/docs/v2/introduction) · [Package documentation](https://pkg.go.dev/github.com/tigusigalpa/manus-ai-go/v2)

> This is an independent community SDK. Manus API behaviour, available agent profiles, and accepted option values are controlled by Manus; check their documentation when upgrading an integration.

## What you can do

- Create, inspect, update, stop, and delete tasks.
- Read task messages and send a follow-up message to an existing task.
- Create upload slots, upload bytes, attach files, list files, and delete files.
- Register webhooks and turn an incoming JSON payload into a useful Go value.
- Configure the base URL, HTTP client, or request timeout for tests and production environments.

## Requirements and installation

The module requires Go 1.21 or newer.

```bash
go get github.com/tigusigalpa/manus-ai-go/v2
```

The `/v2` suffix is important: it selects the API-v2-compatible version of the SDK.

## Quick start

Store the API key outside source code. For local work, an environment variable is a convenient choice:

```bash
export MANUS_AI_API_KEY="your-api-key"
```

```go
package main

import (
	"fmt"
	"log"
	"os"

	manusai "github.com/tigusigalpa/manus-ai-go/v2"
)

func main() {
	client, err := manusai.NewClient(os.Getenv("MANUS_AI_API_KEY"))
	if err != nil {
		log.Fatal(err)
	}

	task, err := client.CreateTask("Write a short release note for my Go project.", &manusai.TaskOptions{
		AgentProfile: manusai.AgentProfileManus16,
		Title:        "Release note",
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Created task %s: %s\n", task.TaskID, task.TaskURL)
}
```

`NewClient` validates an empty API key, a nil HTTP client, and an invalid base URL before making any request. By default it sends requests to `https://api.manus.ai` and uses a 30-second overall request timeout.

## Client configuration

Most programs only need `NewClient(key)`. Options are useful for a custom timeout, an instrumented client, or a local test server.

```go
httpClient := &http.Client{Timeout: 45 * time.Second}
client, err := manusai.NewClient(apiKey,
	manusai.WithHTTPClient(httpClient),
	manusai.WithBaseURL("https://api.manus.ai"),
)
```

`WithTimeout` changes the timeout on the SDK's HTTP client. If you need a custom transport, retry policy, proxy, or tracing, create an `http.Client` yourself and pass it with `WithHTTPClient`.

For OAuth Open Apps, create the client with an empty API-key argument and `WithBearerToken(accessToken)`. The SDK sends either `Authorization: Bearer …` or `x-manus-api-key`, never both.

## Working with tasks

### Create a task

`TaskOptions` is optional. Use only the options that are relevant to the task; empty values are not sent.

```go
hideFromList := true
task, err := client.CreateTask("Summarize this design decision.", &manusai.TaskOptions{
	AgentProfile:     manusai.AgentProfileManus16,
	Locale:           "en-US",
	Title:            "Design summary",
	ShareVisibility:  "private",
	HideInTaskList:   &hideFromList,
	ProjectID:        "project_123",
	Connectors:       []string{"connector_123"},
	EnableSkills:     []string{"skill_123"},
	ForceSkills:      []string{"skill_456"},
})
if err != nil {
	log.Fatal(err)
}
```

The SDK ships profile constants and small discovery helpers:

```go
profiles := manusai.RecommendedAgentProfiles()
if manusai.IsDeprecatedAgentProfile(manusai.AgentProfileSpeed) {
	// Choose a current profile instead.
}
```

### Inspect, list, and update tasks

```go
detail, err := client.GetTask(task.TaskID)
if err != nil {
	log.Fatal(err)
}
fmt.Println(detail.AgentStatus)

tasks, err := client.GetTasks(&manusai.TaskFilters{
	Limit:     20,
	Order:     "desc",
	ProjectID: "project_123",
})
if err != nil {
	log.Fatal(err)
}
for _, item := range tasks.Tasks {
	fmt.Printf("%s — %s\n", item.ID, item.AgentStatus)
}

title := "Revised title"
_, err = client.UpdateTask(task.TaskID, &manusai.TaskUpdate{Title: &title})
```

List responses use cursor pagination. Pass `NextCursor` back as `TaskFilters.Cursor` when `HasMore` is true.

### Follow progress or continue a task

Use `ListMessages` to poll a task, `SendMessage` to continue the conversation, and `StopTask` to request cancellation.

```go
messages, err := client.ListMessages(task.TaskID, 50, "", "desc", true)
if err != nil {
	log.Fatal(err)
}

_, err = client.SendMessage(task.TaskID, "Please make the summary more concise.", nil)
if err != nil {
	log.Fatal(err)
}

_, err = client.StopTask(task.TaskID)
```

When Manus asks the user to confirm an action, pass the task ID, event ID, and the input required by the API:

```go
_, err := client.ConfirmAction(task.TaskID, "event_123", map[string]interface{}{
	"confirmed": true,
})
```

## Files and attachments

The usual file flow has three stages: request an upload URL, upload bytes to that URL, then include the returned file ID in a task.

```go
contents, err := os.ReadFile("report.pdf")
if err != nil {
	log.Fatal(err)
}

file, err := client.CreateFile("report.pdf")
if err != nil {
	log.Fatal(err)
}
if err := client.UploadFileContent(file.UploadURL, contents, "application/pdf"); err != nil {
	log.Fatal(err)
}

task, err := client.CreateTask("Review the attached report.", &manusai.TaskOptions{
	Attachments: []interface{}{manusai.NewAttachmentFromFileID(file.FileID)},
})
```

There are helpers for every supported attachment source:

```go
fromFileID := manusai.NewAttachmentFromFileID("file_123")
fromURL := manusai.NewAttachmentFromURL("https://example.com/report.pdf")
fromBase64 := manusai.NewAttachmentFromBase64(encoded, "application/pdf")
fromPath, err := manusai.NewAttachmentFromFilePath("report.pdf")
```

`NewAttachmentFromFilePath` reads the complete local file into memory. Prefer the upload flow for large files. `CreateTask` and `SendMessage` accept attachments created by these helpers (or equivalent `map[string]interface{}` values); invalid attachment values are rejected rather than silently dropped.

The v2 API does not expose a file-list endpoint. Keep the file ID returned by `CreateFile` if you need to inspect or delete a file later:

```go
_, err = client.DeleteFile("file_123")
```

## Webhooks

Webhooks are preferable to frequent polling when your application needs to react to task events. Your endpoint must be publicly reachable by Manus and should validate any authentication or signature mechanism required by the Manus webhook documentation.

```go
_, err := client.CreateWebhook(&manusai.WebhookConfig{
	URL: "https://example.com/webhooks/manus",
})
```

The helper below parses the body and distinguishes a finished task from a task waiting for the user:

```go
func handleWebhook(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	defer r.Body.Close()

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

	switch {
	case manusai.IsTaskCompleted(payload):
		fmt.Printf("Task completed: %#v\n", manusai.GetTaskDetail(payload))
	case manusai.IsTaskAskingForInput(payload):
		fmt.Printf("Task needs input: %#v\n", manusai.GetTaskDetail(payload))
	}
	w.WriteHeader(http.StatusOK)
}
```

`GetAttachments(payload)` returns the attachments supplied with a webhook task detail, or `nil` when none are present. Delete a webhook with `client.DeleteWebhook(webhookID)` when it is no longer needed.

## Errors

The SDK returns ordinary Go errors. Its API errors can be inspected with `errors.As`:

```go
_, err := client.GetTask("missing-task")
if err != nil {
	var authErr *manusai.AuthenticationError
	var validationErr *manusai.ValidationError
	switch {
	case errors.As(err, &authErr):
		log.Printf("check the API key: %s", authErr.Message)
	case errors.As(err, &validationErr):
		log.Printf("fix the request: %s", validationErr.Message)
	default:
		log.Printf("Manus request failed: %v", err)
	}
}
```

`AuthenticationError` represents 401 and 403 responses; `ValidationError` represents 400 responses; other HTTP and transport errors are represented by `Error`. Each error includes `StatusCode` when an HTTP response was received.

> **Breaking change:** the general error type was renamed from `ManusAIError` to `Error` to follow Go naming conventions. Update type assertions and `errors.As` targets accordingly:
>
> ```go
> // Before
> var apiErr *manusai.ManusAIError
>
> // Now
> var apiErr *manusai.Error
> ```

## API summary

| Area | Methods |
| --- | --- |
| Tasks | `CreateTask`, `GetTasks`, `GetTask`, `UpdateTask`, `DeleteTask`, `ListMessages`, `SendMessage`, `StopTask`, `ConfirmAction` |
| Files | `CreateFile`, `UploadFileContent`, `GetFile`, `DeleteFile` |
| Projects and discovery | `CreateProject`, `ListProjects`, `ListSkills`, `ListAgents`, `GetAgent`, `UpdateAgent`, `ListConnectors` |
| Webhooks and browser | `CreateWebhook`, `ListWebhooks`, `DeleteWebhook`, `GetWebhookPublicKey`, `ListOnlineBrowserClients`, `ParseWebhookPayload` and predicate helpers |
| Usage | `ListUsage`, `GetTeamUsageStatistic`, `ListTeamUsageLog`, `GetAvailableCredits` |
| Websites | `GetWebsiteStatus`, `ListWebsiteCheckpoints`, `PublishWebsite`, `UpdateWebsite` |

The `examples/` directory contains complete programs for [basic task management](examples/basic), [file upload](examples/file-upload), and [webhook handling](examples/webhook).

## Development

Run the complete local check suite:

```bash
gofmt -w .
go vet ./...
go test ./...
```

Or, when `make` is available:

```bash
make check
```

## Migration from v1

Version 2 uses Manus API v2. Imports must end in `/v2`; requests use the `x-manus-api-key` header; endpoints use names such as `/v2/task.create`; and task creation now sends a structured `message.content` payload.

> **OpenAPI alignment (v2.1.8):** `interactive_mode` and `enable_visible_in_task_list` are the current API fields. `EnableAskUser` and `HideInTaskList` remain as compatibility aliases. `ListFiles` is deprecated because the official v2 API has no `/v2/file.list` endpoint. The legacy `speed` and `quality` agent profiles are likewise no longer sent as valid v2 profiles.

## License

Released under the [MIT License](LICENSE).
