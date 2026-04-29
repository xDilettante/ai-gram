// Package errors contains typed errors returned by ai-gram.
package errors

import "fmt"

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
