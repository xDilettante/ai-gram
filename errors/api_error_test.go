package errors

import (
	"context"
	stderrors "errors"
	"fmt"
	"net/url"
	"testing"
)

func TestAPIErrorErrorIsNotEmpty(t *testing.T) {
	err := (&APIError{Code: 400, Description: "bad request"}).Error()
	if err == "" {
		t.Fatal("expected non-empty error string")
	}
}

func TestAPIErrorHelpersWorkThroughWrappedErrors(t *testing.T) {
	err := fmt.Errorf("send message: %w", &APIError{
		Code:        429,
		Description: "Too Many Requests: retry after 17",
		Parameters:  &ResponseParameters{RetryAfter: 17},
	})

	apiErr, ok := AsAPIError(err)
	if !ok {
		t.Fatal("expected wrapped APIError")
	}
	if apiErr.Code != 429 {
		t.Fatalf("unexpected APIError code: %d", apiErr.Code)
	}
	if !IsAPIError(err) {
		t.Fatal("expected IsAPIError to report true")
	}
	if !IsRateLimited(err) {
		t.Fatal("expected IsRateLimited to report true")
	}
	retryAfter, ok := RetryAfter(err)
	if !ok || retryAfter != 17 {
		t.Fatalf("unexpected RetryAfter result: %d %v", retryAfter, ok)
	}
}

func TestMigrateToChatID(t *testing.T) {
	err := fmt.Errorf("send message: %w", &APIError{
		Code:       400,
		Parameters: &ResponseParameters{MigrateToChatID: -1001234567890},
	})

	chatID, ok := MigrateToChatID(err)
	if !ok || chatID != -1001234567890 {
		t.Fatalf("unexpected MigrateToChatID result: %d %v", chatID, ok)
	}
}

func TestAPIErrorStatusHelpers(t *testing.T) {
	tests := []struct {
		name      string
		err       error
		forbidden bool
		notFound  bool
	}{
		{
			name:      "forbidden code",
			err:       &APIError{Code: 403, Description: "bot was blocked by the user"},
			forbidden: true,
		},
		{
			name:      "forbidden description",
			err:       &APIError{Code: 400, Description: "Forbidden: bot is not a member"},
			forbidden: true,
		},
		{
			name:     "not found code",
			err:      &APIError{Code: 404, Description: "Not Found"},
			notFound: true,
		},
		{
			name:     "not found description",
			err:      &APIError{Code: 400, Description: "Bad Request: message to delete not found"},
			notFound: true,
		},
		{
			name: "unrelated API error",
			err:  &APIError{Code: 400, Description: "Bad Request: chat not modified"},
		},
		{
			name: "unrelated error",
			err:  stderrors.New("network unavailable"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsForbidden(tt.err); got != tt.forbidden {
				t.Fatalf("IsForbidden() = %v, want %v", got, tt.forbidden)
			}
			if got := IsNotFound(tt.err); got != tt.notFound {
				t.Fatalf("IsNotFound() = %v, want %v", got, tt.notFound)
			}
		})
	}
}

func TestAPIErrorHelpersIgnoreUnrelatedErrors(t *testing.T) {
	err := stderrors.New("plain error")

	if IsAPIError(err) {
		t.Fatal("expected IsAPIError to report false")
	}
	if IsRateLimited(err) {
		t.Fatal("expected IsRateLimited to report false")
	}
	if retryAfter, ok := RetryAfter(err); ok || retryAfter != 0 {
		t.Fatalf("unexpected RetryAfter result: %d %v", retryAfter, ok)
	}
	if chatID, ok := MigrateToChatID(err); ok || chatID != 0 {
		t.Fatalf("unexpected MigrateToChatID result: %d %v", chatID, ok)
	}
}

func TestNetworkAndContextErrorHelpers(t *testing.T) {
	networkErr := &url.Error{Op: "Get", URL: "https://api.telegram.org", Err: timeoutError{}}
	if !IsNetworkError(networkErr) {
		t.Fatal("expected network error")
	}

	if IsNetworkError(context.Canceled) {
		t.Fatal("expected context.Canceled to stay out of network classification")
	}
	if !IsContextCanceled(fmt.Errorf("handler stopped: %w", context.Canceled)) {
		t.Fatal("expected wrapped context.Canceled")
	}

	if IsNetworkError(context.DeadlineExceeded) {
		t.Fatal("expected context.DeadlineExceeded to stay out of network classification")
	}
	if !IsContextDeadlineExceeded(fmt.Errorf("request timeout: %w", context.DeadlineExceeded)) {
		t.Fatal("expected wrapped context.DeadlineExceeded")
	}

	if IsNetworkError(&APIError{Code: 500, Description: "Internal Server Error"}) {
		t.Fatal("expected Telegram API errors to stay out of network classification")
	}
}

type timeoutError struct{}

func (timeoutError) Error() string   { return "timeout" }
func (timeoutError) Timeout() bool   { return true }
func (timeoutError) Temporary() bool { return true }
