package main

import (
	"context"
	"testing"
	"time"

	apierrors "github.com/xDilettante/ai-gram/errors"
)

func TestClassifyRetry(t *testing.T) {
	config := senderConfig{BaseDelay: time.Second, MaxDelay: 5 * time.Second}

	tests := []struct {
		name       string
		err        error
		wantRetry  bool
		wantReason string
		wantDelay  time.Duration
	}{
		{
			name: "retry after",
			err: &apierrors.APIError{
				Code:       429,
				Parameters: &apierrors.ResponseParameters{RetryAfter: 7},
			},
			wantRetry:  true,
			wantReason: "rate_limited",
			wantDelay:  7 * time.Second,
		},
		{
			name: "rate limit without retry after",
			err: &apierrors.APIError{
				Code: 429,
			},
			wantRetry:  true,
			wantReason: "rate_limited",
			wantDelay:  time.Second,
		},
		{
			name: "migrated chat",
			err: &apierrors.APIError{
				Code:       400,
				Parameters: &apierrors.ResponseParameters{MigrateToChatID: -1001234567890},
			},
			wantReason: "chat_migrated",
		},
		{
			name:       "forbidden",
			err:        &apierrors.APIError{Code: 403, Description: "Forbidden: bot was blocked"},
			wantReason: "forbidden",
		},
		{
			name:       "not found",
			err:        &apierrors.APIError{Code: 404, Description: "Not Found"},
			wantReason: "not_found",
		},
		{
			name:       "attempt timeout",
			err:        context.DeadlineExceeded,
			wantRetry:  true,
			wantReason: "attempt_timeout",
			wantDelay:  time.Second,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := classifyRetry(tt.err, config, 1)
			if got.Retry != tt.wantRetry {
				t.Fatalf("Retry = %v, want %v", got.Retry, tt.wantRetry)
			}
			if got.Reason != tt.wantReason {
				t.Fatalf("Reason = %q, want %q", got.Reason, tt.wantReason)
			}
			if got.Delay != tt.wantDelay {
				t.Fatalf("Delay = %s, want %s", got.Delay, tt.wantDelay)
			}
		})
	}
}

func TestBackoffDelayCapsAtMaxDelay(t *testing.T) {
	config := senderConfig{BaseDelay: time.Second, MaxDelay: 3 * time.Second}
	if got := backoffDelay(config, 4); got != 3*time.Second {
		t.Fatalf("backoffDelay() = %s, want 3s", got)
	}
}

func TestSafeChatLabel(t *testing.T) {
	if got := safeChatLabel("-1001234567890"); got != "-10***890" {
		t.Fatalf("safeChatLabel(numeric) = %q", got)
	}
	if got := safeChatLabel("@example_channel"); got != "username" {
		t.Fatalf("safeChatLabel(username) = %q", got)
	}
}
