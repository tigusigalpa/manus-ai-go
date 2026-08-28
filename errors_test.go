package manusai

import (
	"errors"
	"testing"
)

func TestErrors(t *testing.T) {
	cause := errors.New("cause")
	tests := []struct {
		name string
		err  interface {
			error
			Unwrap() error
		}
		want string
	}{
		{"general", &Error{Message: "failure", StatusCode: 500, Err: cause}, "manus-ai error (status 500): failure"},
		{"authentication", &AuthenticationError{Message: "denied", Err: cause}, "authentication error: denied"},
		{"validation", &ValidationError{Message: "invalid", StatusCode: 400, Err: cause}, "validation error (status 400): invalid"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.err.Error(); got != tt.want {
				t.Fatalf("Error() = %q, want %q", got, tt.want)
			}
			if !errors.Is(tt.err, cause) {
				t.Fatal("Unwrap() did not return the wrapped error")
			}
		})
	}
}
