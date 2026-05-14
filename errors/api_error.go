// Package errors contains typed errors returned by ai-gram.
package errors

import (
	"context"
	stderrors "errors"
	"fmt"
	"net"
	"strings"
)

// ResponseParameters describes optional Telegram Bot API error parameters.
type ResponseParameters struct {
	MigrateToChatID int64 `json:"migrate_to_chat_id,omitempty"`
	RetryAfter      int   `json:"retry_after,omitempty"`
}

// APIError represents a Telegram Bot API response with ok=false.
type APIError struct {
	Code        int                 `json:"error_code"`
	Description string              `json:"description"`
	Parameters  *ResponseParameters `json:"parameters,omitempty"`
}

// Error returns a human-readable Telegram API error message.
func (e *APIError) Error() string {
	if e == nil {
		return "telegram API error"
	}
	if e.Description == "" {
		return fmt.Sprintf("telegram API error: code %d", e.Code)
	}
	if e.Code == 0 {
		return "telegram API error: " + e.Description
	}

	return fmt.Sprintf("telegram API error: code %d: %s", e.Code, e.Description)
}

// AsAPIError returns the Telegram API error in err, if any.
func AsAPIError(err error) (*APIError, bool) {
	var apiErr *APIError
	if stderrors.As(err, &apiErr) {
		return apiErr, true
	}
	return nil, false
}

// IsAPIError reports whether err wraps a Telegram API error.
func IsAPIError(err error) bool {
	_, ok := AsAPIError(err)
	return ok
}

// IsRateLimited reports whether err is a Telegram rate-limit response.
func IsRateLimited(err error) bool {
	apiErr, ok := AsAPIError(err)
	if !ok {
		return false
	}
	if apiErr.Code == 429 {
		return true
	}
	return apiErr.Parameters != nil && apiErr.Parameters.RetryAfter > 0
}

// RetryAfter returns the Telegram retry_after value in seconds, if present.
func RetryAfter(err error) (int, bool) {
	apiErr, ok := AsAPIError(err)
	if !ok || apiErr.Parameters == nil || apiErr.Parameters.RetryAfter <= 0 {
		return 0, false
	}
	return apiErr.Parameters.RetryAfter, true
}

// MigrateToChatID returns the Telegram migrate_to_chat_id value, if present.
func MigrateToChatID(err error) (int64, bool) {
	apiErr, ok := AsAPIError(err)
	if !ok || apiErr.Parameters == nil || apiErr.Parameters.MigrateToChatID == 0 {
		return 0, false
	}
	return apiErr.Parameters.MigrateToChatID, true
}

// IsForbidden reports whether err is a Telegram forbidden-style API response.
func IsForbidden(err error) bool {
	apiErr, ok := AsAPIError(err)
	if !ok {
		return false
	}
	return apiErr.Code == 403 || containsFold(apiErr.Description, "forbidden")
}

// IsNotFound reports whether err is a Telegram not-found-style API response.
func IsNotFound(err error) bool {
	apiErr, ok := AsAPIError(err)
	if !ok {
		return false
	}
	return apiErr.Code == 404 || containsFold(apiErr.Description, "not found")
}

// IsNetworkError reports whether err is a non-context network error.
func IsNetworkError(err error) bool {
	if err == nil || IsAPIError(err) || IsContextCanceled(err) || IsContextDeadlineExceeded(err) {
		return false
	}

	var netErr net.Error
	return stderrors.As(err, &netErr)
}

// IsContextCanceled reports whether err wraps context.Canceled.
func IsContextCanceled(err error) bool {
	return stderrors.Is(err, context.Canceled)
}

// IsContextDeadlineExceeded reports whether err wraps context.DeadlineExceeded.
func IsContextDeadlineExceeded(err error) bool {
	return stderrors.Is(err, context.DeadlineExceeded)
}

func containsFold(s, substr string) bool {
	return strings.Contains(strings.ToLower(s), substr)
}
