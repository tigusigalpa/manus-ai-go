# Contributing to Manus AI Go SDK

First off, thank you for considering contributing to Manus AI Go SDK! It's people like you that make this project better.

## Code of Conduct

This project and everyone participating in it is governed by our Code of Conduct. By participating, you are expected to uphold this code.

## How Can I Contribute?

### Reporting Bugs

Before creating bug reports, please check the existing issues as you might find out that you don't need to create one. When you are creating a bug report, please include as many details as possible:

* **Use a clear and descriptive title**
* **Describe the exact steps which reproduce the problem**
* **Provide specific examples to demonstrate the steps**
* **Describe the behavior you observed after following the steps**
* **Explain which behavior you expected to see instead and why**
* **Include Go version and OS details**

### Suggesting Enhancements

Enhancement suggestions are tracked as GitHub issues. When creating an enhancement suggestion, please include:

* **Use a clear and descriptive title**
* **Provide a step-by-step description of the suggested enhancement**
* **Provide specific examples to demonstrate the steps**
* **Describe the current behavior and explain which behavior you expected to see instead**
* **Explain why this enhancement would be useful**

### Pull Requests

* Fill in the required template
* Do not include issue numbers in the PR title
* Follow the Go coding style
* Include thoughtfully-worded, well-structured tests
* Document new code
* End all files with a newline

## Development Process

1. Fork the repo
2. Create a new branch from `main`
3. Make your changes
4. Add tests for your changes
5. Run tests: `go test -v ./...`
6. Run linter: `golangci-lint run`
7. Format code: `go fmt ./...`
8. Commit your changes
9. Push to your fork
10. Create a Pull Request

## Coding Standards

* Follow standard Go conventions and idioms
* Write clear, readable code with appropriate comments
* Keep functions focused and concise
* Use meaningful variable and function names
* Write comprehensive tests for new features
* Maintain backward compatibility when possible

## Testing

* Write unit tests for all new code
* Ensure all tests pass before submitting PR
* Aim for high test coverage
* Use table-driven tests where appropriate

```bash
# Run tests
go test -v ./...

# Run tests with coverage
go test -v -cover ./...

# Generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## Documentation

* Update README.md if you change functionality
* Add godoc comments for exported functions and types
* Include examples in documentation where helpful
* Update CHANGELOG.md for significant changes

## Commit Messages

* Use the present tense ("Add feature" not "Added feature")
* Use the imperative mood ("Move cursor to..." not "Moves cursor to...")
* Limit the first line to 72 characters or less
* Reference issues and pull requests liberally after the first line

## License

By contributing, you agree that your contributions will be licensed under the MIT License.
