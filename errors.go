package manusai

import "fmt"

// Error represents a general API or transport error returned by the SDK.
type Error struct {
	Message    string
	StatusCode int
	Err        error
}

func (e *Error) Error() string {
	if e.StatusCode > 0 {
		return fmt.Sprintf("manus-ai error (status %d): %s", e.StatusCode, e.Message)
	}
	return fmt.Sprintf("manus-ai error: %s", e.Message)
}

func (e *Error) Unwrap() error {
	return e.Err
}

// AuthenticationError represents an authentication or authorization failure.
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

// ValidationError represents invalid client input or an API validation failure.
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
