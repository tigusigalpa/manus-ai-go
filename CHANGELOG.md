# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

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

[Unreleased]: https://github.com/tigusigalpa/manus-ai-go/compare/v1.0.0...HEAD
[1.0.0]: https://github.com/tigusigalpa/manus-ai-go/releases/tag/v1.0.0
