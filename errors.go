package manusai

import "fmt"

type ManusAIError struct {
	Message    string
	StatusCode int
	Err        error
}

func (e *ManusAIError) Error() string {
	if e.StatusCode > 0 {
		return fmt.Sprintf("manus-ai error (status %d): %s", e.StatusCode, e.Message)
	}
	return fmt.Sprintf("manus-ai error: %s", e.Message)
}

func (e *ManusAIError) Unwrap() error {
	return e.Err
}

type AuthenticationError struct {
	Message    string
	StatusCode int
	Err        error
}

func (e *AuthenticationError) Error() string {
	if e.StatusCode > 0 {
		return fmt.Sprintf("authentication error (status %d): %s", e.StatusCode, e.Message)
	}
	return fmt.Sprintf("authentication error: %s", e.Message)
}

func (e *AuthenticationError) Unwrap() error {
	return e.Err
}

type ValidationError struct {
	Message    string
	StatusCode int
	Err        error
}

func (e *ValidationError) Error() string {
	if e.StatusCode > 0 {
		return fmt.Sprintf("validation error (status %d): %s", e.StatusCode, e.Message)
	}
	return fmt.Sprintf("validation error: %s", e.Message)
}

func (e *ValidationError) Unwrap() error {
	return e.Err
}
