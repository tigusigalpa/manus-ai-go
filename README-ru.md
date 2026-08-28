# Manus AI Go SDK

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

Небольшой и понятный Go-клиент для Manus API v2. Он берёт на себя HTTP-запросы, чтобы приложение могло создавать и вести задачи Manus, работать с файлами и принимать вебхуки обычным Go-кодом.

[English](README.md) · [Документация Manus API](https://open.manus.im/docs/v2/introduction) · [Документация пакета](https://pkg.go.dev/github.com/tigusigalpa/manus-ai-go/v2)

> Это независимый community SDK. Доступные профили агентов и допустимые значения опций определяет API Manus, поэтому перед обновлением интеграции сверяйтесь с официальной документацией.

## Возможности

- Создание, просмотр, обновление, остановка и удаление задач.
- Чтение сообщений задачи и продолжение диалога через `SendMessage`.
- Создание слота для файла, загрузка содержимого и добавление файла во вложения.
- Регистрация вебхуков и удобный разбор входящего JSON.
- Настройка адреса API, HTTP-клиента и таймаута.

## Требования и установка

Нужен Go 1.21 или новее.

```bash
go get github.com/tigusigalpa/manus-ai-go/v2
```

Суффикс `/v2` обязателен: он подключает версию библиотеки для Manus API v2.

## Быстрый старт

Не храните API-ключ в исходниках. Для локальной разработки удобно использовать переменную окружения:

```bash
export MANUS_AI_API_KEY="ваш-api-ключ"
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

	task, err := client.CreateTask("Напиши короткий релиз-ноут для Go-проекта.", &manusai.TaskOptions{
		AgentProfile: manusai.AgentProfileManus16,
		Title:        "Релиз-ноут",
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Создана задача %s: %s\n", task.TaskID, task.TaskURL)
}
```

По умолчанию запросы идут на `https://api.manus.ai`, а общий таймаут запроса равен 30 секундам. Конструктор сразу отклоняет пустой ключ, `nil` HTTP-клиент и некорректный базовый URL.

## Настройка клиента

В большинстве случаев достаточно `NewClient(key)`. Для тестового сервера, трассировки, прокси или другого таймаута передайте опции:

```go
httpClient := &http.Client{Timeout: 45 * time.Second}
client, err := manusai.NewClient(apiKey,
	manusai.WithHTTPClient(httpClient),
	manusai.WithBaseURL("https://api.manus.ai"),
)
```

`WithTimeout` меняет таймаут HTTP-клиента библиотеки. Если нужен собственный transport, повторные попытки или трассировка, создайте `http.Client` самостоятельно и передайте его через `WithHTTPClient`.

## Работа с задачами

### Создание

`TaskOptions` необязателен. Пустые поля не отправляются в API.

```go
hideFromList := true
task, err := client.CreateTask("Суммируй это архитектурное решение.", &manusai.TaskOptions{
	AgentProfile:    manusai.AgentProfileManus16,
	Locale:          "ru-RU",
	Title:           "Краткое решение",
	ShareVisibility: "private",
	HideInTaskList:  &hideFromList,
	ProjectID:       "project_123",
	Connectors:      []string{"connector_123"},
	EnableSkills:    []string{"skill_123"},
})
if err != nil {
	log.Fatal(err)
}
```

Для профилей есть готовые константы и справочные функции:

```go
profiles := manusai.RecommendedAgentProfiles()
deprecated := manusai.IsDeprecatedAgentProfile(manusai.AgentProfileSpeed)
```

### Просмотр, список и обновление

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

title := "Новое название"
_, err = client.UpdateTask(task.TaskID, &manusai.TaskUpdate{Title: &title})
```

Для постраничного списка передавайте `NextCursor` в `TaskFilters.Cursor`, пока `HasMore` равно `true`.

### Прогресс и продолжение диалога

```go
messages, err := client.ListMessages(task.TaskID, 50, "", "desc", true)
if err != nil {
	log.Fatal(err)
}

_, err = client.SendMessage(task.TaskID, "Сделай итог короче.", nil)
_, err = client.StopTask(task.TaskID)
```

Если Manus запросил подтверждение действия, используйте идентификаторы задачи и события вместе с ожидаемыми данными:

```go
_, err := client.ConfirmAction(task.TaskID, "event_123", map[string]interface{}{
	"confirmed": true,
})
```

## Файлы и вложения

Стандартный сценарий состоит из трёх шагов: получить URL для загрузки, загрузить байты и передать `file_id` при создании задачи.

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

_, err = client.CreateTask("Проанализируй приложенный отчёт.", &manusai.TaskOptions{
	Attachments: []interface{}{manusai.NewAttachmentFromFileID(file.FileID)},
})
```

Вложения можно создать из ID, URL, base64-строки или локального файла:

```go
fromID := manusai.NewAttachmentFromFileID("file_123")
fromURL := manusai.NewAttachmentFromURL("https://example.com/report.pdf")
fromBase64 := manusai.NewAttachmentFromBase64(encoded, "application/pdf")
fromPath, err := manusai.NewAttachmentFromFilePath("report.pdf")
```

`NewAttachmentFromFilePath` полностью читает файл в память, поэтому для больших файлов лучше использовать поток «создать слот → загрузить → приложить». Некорректные элементы `Attachments` теперь возвращают ошибку, а не исчезают из запроса молча.

```go
files, err := client.ListFiles(20, "")
if err != nil {
	log.Fatal(err)
}
for _, file := range files.Files {
	fmt.Println(file.FileID, file.Filename, file.Status)
}

_, err = client.DeleteFile("file_123")
```

## Вебхуки

Вебхуки полезнее частого polling, если вашему сервису нужно реагировать на события задач. URL должен быть доступен Manus извне. Также проверьте в документации Manus требования к аутентификации или проверке подписи.

```go
_, err := client.CreateWebhook(&manusai.WebhookConfig{
	URL:    "https://example.com/webhooks/manus",
	Events: []string{"task_created", "task_stopped"},
})
```

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
		fmt.Printf("Задача завершена: %#v\n", manusai.GetTaskDetail(payload))
	case manusai.IsTaskAskingForInput(payload):
		fmt.Printf("Нужен ответ пользователя: %#v\n", manusai.GetTaskDetail(payload))
	}
	w.WriteHeader(http.StatusOK)
}
```

`GetAttachments(payload)` возвращает вложения из деталей задачи или `nil`, если их нет. Ненужный вебхук удаляется вызовом `client.DeleteWebhook(webhookID)`.

## Ошибки

Ошибки — обычные Go errors, их можно безопасно различать через `errors.As`:

```go
var authErr *manusai.AuthenticationError
var validationErr *manusai.ValidationError

_, err := client.GetTask("missing-task")
switch {
case errors.As(err, &authErr):
	log.Printf("проверьте API-ключ: %s", authErr.Message)
case errors.As(err, &validationErr):
	log.Printf("исправьте запрос: %s", validationErr.Message)
case err != nil:
	log.Printf("ошибка запроса к Manus: %v", err)
}
```

`AuthenticationError` соответствует ответам 401/403, `ValidationError` — 400, остальные HTTP- и транспортные ошибки представлены типом `Error`. Когда сервер успел ответить, `StatusCode` содержит HTTP-код.

> **Критическое изменение:** общий тип ошибки переименован с `ManusAIError` в `Error`, чтобы соответствовать соглашениям Go. Обновите type assertions и цели `errors.As`:
>
> ```go
> // Было
> var apiErr *manusai.ManusAIError
>
> // Стало
> var apiErr *manusai.Error
> ```

## Краткий API

| Область | Методы |
| --- | --- |
| Задачи | `CreateTask`, `GetTasks`, `GetTask`, `UpdateTask`, `DeleteTask`, `ListMessages`, `SendMessage`, `StopTask`, `ConfirmAction` |
| Файлы | `CreateFile`, `UploadFileContent`, `ListFiles(limit, cursor)`, `GetFile`, `DeleteFile` |
| Вебхуки | `CreateWebhook`, `DeleteWebhook`, `ParseWebhookPayload` и функции-предикаты |

Готовые программы находятся в [examples/basic](examples/basic), [examples/file-upload](examples/file-upload) и [examples/webhook](examples/webhook).

## Разработка

```bash
gofmt -w .
go vet ./...
go test ./...
```

Если доступен `make`, можно запустить `make check`.

## Переход с v1

Версия 2 использует Manus API v2. Импорт должен оканчиваться на `/v2`; ключ передаётся в заголовке `x-manus-api-key`; вызовы идут на endpoints наподобие `/v2/task.create`; создание задачи передаёт структурированный `message.content`. Поля v1 `TaskMode`, `TaskDetail.Status`, `TaskListResponse.Data`, `FileResponse.ID` и `FileListResponse.Data` больше не применимы. Используйте соответственно `AgentStatus`, `Tasks`, `FileID` и `Files`.

## Лицензия

[MIT](LICENSE).
