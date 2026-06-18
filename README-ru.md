# Manus AI Go SDK

![Manus AI Golang SDK](https://i.postimg.cc/6pm3pLcK/manus-ai-api-go-sdk.png)

Go-клиент для API v2 [Manus AI](https://manus.ai). Задачи, мультитёрн диалоги, загрузка файлов, вебхуки.

**Package:** [pkg.go.dev/github.com/tigusigalpa/manus-ai-go](https://pkg.go.dev/github.com/tigusigalpa/manus-ai-go)

[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/Go-%3E%3D1.21-blue)](https://golang.org/)

Русский | [English](README.md)

> **⚠️ Критические изменения:** Версия 2.0+ использует Manus API v2 со значительными изменениями.
> См. [Руководство по миграции](#миграция-с-v1) ниже.

## Содержание

- [Возможности](#возможности)
- [Требования](#требования)
- [Установка](#установка)
- [Конфигурация](#конфигурация)
- [Использование](#использование)
- [Примеры](#примеры)
- [Тестирование](#тестирование)
- [Лицензия](#лицензия)

## Возможности

- **Полная поддержка Manus AI API v2**
- Создание и управление задачами с новым форматом сообщений
- Мультитёрн диалоги с `SendMessage`
- Управление жизненным циклом задач (`ListMessages`, `StopTask`, `ConfirmAction`)
- Загрузка файлов и вложения (file_id, file_url, file_data)
- Вебхуки
- Поддержка проектов и навыков
- Интеграция с коннекторами
- Кастомные типы ошибок с детальными ответами
- Типобезопасные интерфейсы
- Покрытие тестами
- Идиоматичный Go

### Новое в v2

- **Message-based API**: Задачи используют структурированный формат контента
- **Polling задач**: Используйте `ListMessages` для отслеживания прогресса
- **Интерактивные задачи**: Подтверждайте действия через `ConfirmAction`
- **Мультитёрн диалоги**: Продолжайте беседы с `SendMessage`
- **Расширенные метаданные**: `agent_status`, `share_visibility`, временные метки
- **Cursor-пагинация**: Эффективная пагинация для задач и файлов
- **Навыки и коннекторы**: Включайте конкретные навыки и коннекторы для каждой задачи
- **Проекты**: Группируйте связанные задачи

## Требования

- Go 1.21 или выше

## Установка

```bash
go get github.com/tigusigalpa/manus-ai-go
```

## Конфигурация

### Получение API ключа

1. Зарегистрируйтесь на [Manus AI](https://manus.im)
2. Получите API ключ в [настройках интеграции API](http://manus.im/app?show_settings=integrations&app_name=api)

### Базовая конфигурация

```go
import manusai "github.com/tigusigalpa/manus-ai-go"

client, err := manusai.NewClient("ваш-api-ключ")
if err != nil {
    log.Fatal(err)
}
```

## Использование

### Создание задачи

```go
task, err := client.CreateTask("Напиши стихотворение о Go", &manusai.TaskOptions{
    AgentProfile: manusai.AgentProfileManus16,
    TaskMode:     "chat",
})
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Задача создана: %s\n", task.TaskID)
fmt.Printf("Ссылка: %s\n", task.TaskURL)
```

### Работа с файлами

```go
// Создание записи файла
fileResult, err := client.CreateFile("document.pdf")

// Загрузка содержимого
fileContent, _ := os.ReadFile("/path/to/document.pdf")
err = client.UploadFileContent(fileResult.UploadURL, fileContent, "application/pdf")

// Использование в задаче
attachment := manusai.NewAttachmentFromFileID(fileResult.ID)
task, err := client.CreateTask("Проанализируй документ", &manusai.TaskOptions{
    Attachments: []interface{}{attachment},
})
```

### Вебхуки

```go
webhook := &manusai.WebhookConfig{
    URL:    "https://your-domain.com/webhook/manus-ai",
    Events: []string{"task_created", "task_stopped"},
}

result, err := client.CreateWebhook(webhook)
```

## Примеры

См. директорию `examples/`:

- `examples/basic/` - Базовое создание и управление задачами
- `examples/file-upload/` - Загрузка файлов с вложениями
- `examples/webhook/` - Настройка и обработка вебхуков

## Тестирование

```bash
go test -v ./...
```

С отчётом покрытия:

```bash
go test -v -cover ./...
```

## Лицензия

MIT License — см. [LICENSE](LICENSE).

## Миграция с v1

### Критические изменения

1. **API Endpoints**: Все endpoints изменены с `/v1/` на `/v2/` с новыми названиями (например, `/v2/task.create`)

2. **Заголовок аутентификации**: Изменён с `Authorization` на `x-manus-api-key`

3. **Структура запросов**: Задачи теперь используют формат сообщений:
   ```go
   // v1
   payload := map[string]interface{}{
       "prompt": "Привет",
       "agentProfile": "manus-1.6",
   }
   
   // v2
   payload := map[string]interface{}{
       "message": map[string]interface{}{
           "content": []map[string]interface{}{
               {"type": "text", "text": "Привет"},
           },
       },
       "agent_profile": "manus-1.6",
   }
   ```

4. **Формат ответов**: Все ответы теперь включают поля `ok` и `request_id`

5. **Названия полей**: Snake_case вместо camelCase:
    - `agentProfile` → `agent_profile`
    - `hideInTaskList` → `hide_in_task_list`
    - `createShareableLink` → `share_visibility`

6. **Статус задачи**: `Status` → `AgentStatus`

7. **Временные метки**: Изменены со строк на int64 (Unix миллисекунды)

8. **Вложения**: Новая структура с `file_id`, `file_url`, `file_data`

9. **Удалённые поля**:
    - `TaskMode` (больше не нужен)
    - `CreateShareableLink` (заменён на `ShareVisibility`)

10. **Новые методы**:
    - `ListMessages()` - Отслеживание прогресса задачи
    - `SendMessage()` - Продолжение диалогов
    - `StopTask()` - Остановка выполняющихся задач
    - `ConfirmAction()` - Подтверждение ожидающих действий

### Пример миграции

```go
// v1
task, err := client.CreateTask("Привет", &manusai.TaskOptions{
    AgentProfile: "manus-1.6",
    TaskMode:     "agent",
})

// v2
task, err := client.CreateTask("Привет", &manusai.TaskOptions{
    AgentProfile:    manusai.AgentProfileManus16,
    ShareVisibility: "private",
})

// v2: Отслеживание прогресса
messages, _ := client.ListMessages(task.TaskID, 50, "", "desc", false)
```

## Ссылки

- [Manus AI](https://manus.ai)
- [Документация API v2](https://open.manus.im/docs/v2/introduction)
- [GitHub Repository](https://github.com/tigusigalpa/manus-ai-go)
- [Issues](https://github.com/tigusigalpa/manus-ai-go/issues)

## Автор

**Igor Sazonov**

- GitHub: [@tigusigalpa](https://github.com/tigusigalpa)
- Email: sovletig@gmail.com

Также см. PHP SDK: [manus-ai-php](https://github.com/tigusigalpa/manus-ai-php)
