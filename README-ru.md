# Manus AI Go SDK

![Manus AI Golang SDK](https://github.com/user-attachments/assets/1249e90c-a860-4f86-9a77-2d048f94854d)

Go-клиент для API [Manus AI](https://manus.ai). Задачи, загрузка файлов, вебхуки.

**Package:** [pkg.go.dev/github.com/tigusigalpa/manus-ai-go](https://pkg.go.dev/github.com/tigusigalpa/manus-ai-go)

[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/Go-%3E%3D1.21-blue)](https://golang.org/)

Русский | [English](README.md)

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

- Полная поддержка Manus AI API
- Создание и управление задачами
- Загрузка файлов и вложения
- Вебхуки
- Кастомные типы ошибок
- Типобезопасные интерфейсы
- Покрытие тестами
- Идиоматичный Go

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

## Автор

**Igor Sazonov**

- GitHub: [@tigusigalpa](https://github.com/tigusigalpa)
- Email: sovletig@gmail.com

Также см. PHP SDK: [manus-ai-php](https://github.com/tigusigalpa/manus-ai-php)
