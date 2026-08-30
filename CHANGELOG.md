# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [2.2.2] - 2026-08-30

### Fixed
- Normalize Manus API v2 response timestamps to Unix milliseconds while accepting JSON numbers, numeric strings, and RFC3339/RFC3339Nano strings.
- Fix `task.detail` decoding when `task.created_at` or `task.updated_at` is returned as a string.
- Apply compatible timestamp decoding to task lists and messages, files, projects, skills, webhooks, usage, credits, and website response models.

### Tests
- Add real `task.detail` regression payloads covering RFC3339, numeric, numeric-string, and invalid timestamps.

## [2.2.1] - 2026-08-30

### Fixed
- Decode Manus API v2 `task.detail` responses from their nested `task` envelope and populate `TaskDetail.AgentStatus` from `task.status`.
- Add `TaskDetail.TaskURL`, `TaskDetail.TaskType`, and `TaskDetail.AgentProfile` without removing existing fields.
- Decode `task.listMessages` timestamps sent as JSON numbers, numeric strings, or RFC3339/RFC3339Nano strings. RFC3339 values are normalized to Unix milliseconds.
- Return a descriptive decoding error for invalid message timestamps instead of a raw JSON type mismatch.

### Tests
- Add regression coverage for production `task.detail` and `task.listMessages` response payloads.

## [1.0.0] - 2025-01-XX

### Added
- Initial release of Manus AI Go SDK
- Full support for Manus AI API
- Task creation and management (create, get, list, update, delete)
- File upload and attachment handling
- Webhook integration for real-time updates
- Agent profile constants and helper functions
- Attachment helper functions (file ID, URL, base64, file path)
- Webhook payload parsing and event detection
- Comprehensive error handling with custom error types
- Full test coverage with unit tests
- Complete documentation and examples
- GitHub Actions CI/CD workflow
- Makefile for common tasks

### Features
- Type-safe interfaces for all API operations
- Idiomatic Go code following best practices
- Support for custom HTTP clients and timeouts
- Multiple agent profile options (Manus 1.6, Lite, Max)
- Multiple attachment types (file ID, URL, base64, local file)
- Webhook event handlers (task created, stopped, completed, asking for input)
- Comprehensive examples (basic usage, file upload, webhooks)

[Unreleased]: https://github.com/tigusigalpa/manus-ai-go/compare/v2.2.2...HEAD
[2.2.2]: https://github.com/tigusigalpa/manus-ai-go/compare/v2.2.1...v2.2.2
[2.2.1]: https://github.com/tigusigalpa/manus-ai-go/compare/v2.2.0...v2.2.1
[1.0.0]: https://github.com/tigusigalpa/manus-ai-go/releases/tag/v1.0.0
