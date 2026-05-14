// Example 08_retry_sender shows explicit retry handling for outgoing messages.
package main

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	aigram "github.com/xDilettante/ai-gram"
	apierrors "github.com/xDilettante/ai-gram/errors"
	"github.com/xDilettante/ai-gram/examples/internal/exampleutil"
)

const (
	defaultText           = "ai-gram retry sender example"
	defaultMaxAttempts    = 4
	defaultAttemptTimeout = 10 * time.Second
	defaultBaseDelay      = time.Second
	defaultMaxDelay       = 30 * time.Second
)

type senderConfig struct {
	ChatID         aigram.ChatID
	ChatLabel      string
	Text           string
	MaxAttempts    int
	AttemptTimeout time.Duration
	BaseDelay      time.Duration
	MaxDelay       time.Duration
}

type retryDecision struct {
	Retry  bool
	Reason string
	Delay  time.Duration
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	ctx, stop := exampleutil.SignalContext()
	defer stop()

	b, err := exampleutil.NewBotFromEnv()
	if err != nil {
		return err
	}
	config, err := senderConfigFromEnv()
	if err != nil {
		return err
	}

	log.Printf("retry_sender_started chat_id=%s max_attempts=%d attempt_timeout=%s base_delay=%s max_delay=%s",
		config.ChatLabel, config.MaxAttempts, config.AttemptTimeout, config.BaseDelay, config.MaxDelay)
	if err := sendWithRetry(ctx, b, config); err != nil {
		return err
	}
	log.Printf("retry_sender_finished chat_id=%s status=sent", config.ChatLabel)
	return nil
}

func senderConfigFromEnv() (senderConfig, error) {
	rawChatID, err := exampleutil.RequiredEnv("AIGRAM_CHAT_ID")
	if err != nil {
		return senderConfig{}, err
	}
	chatID, err := exampleutil.ParseChatID(rawChatID)
	if err != nil {
		return senderConfig{}, err
	}

	maxAttempts, err := intEnv("AIGRAM_RETRY_MAX_ATTEMPTS", defaultMaxAttempts, 1, 10)
	if err != nil {
		return senderConfig{}, err
	}
	attemptTimeout, err := durationEnv("AIGRAM_RETRY_ATTEMPT_TIMEOUT", defaultAttemptTimeout)
	if err != nil {
		return senderConfig{}, err
	}
	baseDelay, err := durationEnv("AIGRAM_RETRY_BASE_DELAY", defaultBaseDelay)
	if err != nil {
		return senderConfig{}, err
	}
	maxDelay, err := durationEnv("AIGRAM_RETRY_MAX_DELAY", defaultMaxDelay)
	if err != nil {
		return senderConfig{}, err
	}
	if maxDelay < baseDelay {
		return senderConfig{}, fmt.Errorf("AIGRAM_RETRY_MAX_DELAY must be greater than or equal to AIGRAM_RETRY_BASE_DELAY")
	}

	return senderConfig{
		ChatID:         chatID,
		ChatLabel:      safeChatLabel(rawChatID),
		Text:           exampleutil.OptionalEnv("AIGRAM_RETRY_TEXT", defaultText),
		MaxAttempts:    maxAttempts,
		AttemptTimeout: attemptTimeout,
		BaseDelay:      baseDelay,
		MaxDelay:       maxDelay,
	}, nil
}

func sendWithRetry(ctx context.Context, b *aigram.Bot, config senderConfig) error {
	for attempt := 1; attempt <= config.MaxAttempts; attempt++ {
		attemptCtx, cancel := context.WithTimeout(ctx, config.AttemptTimeout)
		started := time.Now()
		_, err := b.SendMessage(attemptCtx, aigram.SendMessageParams{
			ChatID: config.ChatID,
			Text:   config.Text,
		})
		cancel()

		duration := time.Since(started)
		if err == nil {
			log.Printf("send_attempt_ok attempt=%d duration_ms=%d chat_id=%s", attempt, duration.Milliseconds(), config.ChatLabel)
			return nil
		}

		decision := classifyRetry(err, config, attempt)
		shouldRetry := decision.Retry && attempt < config.MaxAttempts
		log.Printf("send_attempt_error attempt=%d duration_ms=%d chat_id=%s retry=%t reason=%s err=%v",
			attempt, duration.Milliseconds(), config.ChatLabel, shouldRetry, decision.Reason, err)
		if !shouldRetry {
			return fmt.Errorf("send message after %d attempt(s): %w", attempt, err)
		}
		if err := waitBeforeRetry(ctx, decision.Delay, attempt, config.ChatLabel, decision.Reason); err != nil {
			return err
		}
	}

	return nil
}

func classifyRetry(err error, config senderConfig, attempt int) retryDecision {
	if retryAfter, ok := apierrors.RetryAfter(err); ok {
		return retryDecision{
			Retry:  true,
			Reason: "rate_limited",
			Delay:  time.Duration(retryAfter) * time.Second,
		}
	}
	if apierrors.IsRateLimited(err) {
		return retryDecision{Retry: true, Reason: "rate_limited", Delay: backoffDelay(config, attempt)}
	}
	if _, ok := apierrors.MigrateToChatID(err); ok {
		return retryDecision{Reason: "chat_migrated"}
	}
	if apierrors.IsForbidden(err) {
		return retryDecision{Reason: "forbidden"}
	}
	if apierrors.IsNotFound(err) {
		return retryDecision{Reason: "not_found"}
	}
	if apierrors.IsContextCanceled(err) {
		return retryDecision{Reason: "context_canceled"}
	}
	if apierrors.IsContextDeadlineExceeded(err) {
		return retryDecision{Retry: true, Reason: "attempt_timeout", Delay: backoffDelay(config, attempt)}
	}
	if apierrors.IsNetworkError(err) {
		return retryDecision{Retry: true, Reason: "network_error", Delay: backoffDelay(config, attempt)}
	}
	if apiErr, ok := apierrors.AsAPIError(err); ok {
		return retryDecision{Reason: fmt.Sprintf("telegram_api_%d", apiErr.Code)}
	}
	return retryDecision{Reason: "unclassified"}
}

func waitBeforeRetry(ctx context.Context, delay time.Duration, attempt int, chatLabel, reason string) error {
	log.Printf("send_attempt_retry attempt=%d delay=%s chat_id=%s reason=%s", attempt, delay, chatLabel, reason)
	timer := time.NewTimer(delay)
	defer timer.Stop()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}

func backoffDelay(config senderConfig, attempt int) time.Duration {
	delay := config.BaseDelay
	for i := 1; i < attempt; i++ {
		delay *= 2
		if delay >= config.MaxDelay {
			return config.MaxDelay
		}
	}
	return delay
}

func intEnv(name string, fallback, minValue, maxValue int) (int, error) {
	raw := strings.TrimSpace(exampleutil.OptionalEnv(name, ""))
	if raw == "" {
		return fallback, nil
	}
	value, err := strconv.Atoi(raw)
	if err != nil {
		return 0, fmt.Errorf("%s must be an integer", name)
	}
	if value < minValue || value > maxValue {
		return 0, fmt.Errorf("%s must be between %d and %d", name, minValue, maxValue)
	}
	return value, nil
}

func durationEnv(name string, fallback time.Duration) (time.Duration, error) {
	raw := strings.TrimSpace(exampleutil.OptionalEnv(name, ""))
	if raw == "" {
		return fallback, nil
	}
	value, err := time.ParseDuration(raw)
	if err != nil {
		return 0, fmt.Errorf("%s must be a Go duration such as 2s or 500ms", name)
	}
	if value <= 0 {
		return 0, fmt.Errorf("%s must be positive", name)
	}
	return value, nil
}

func safeChatLabel(raw string) string {
	if id, err := strconv.ParseInt(strings.TrimSpace(raw), 10, 64); err == nil {
		return exampleutil.MaskInt64(id)
	}
	return "username"
}
