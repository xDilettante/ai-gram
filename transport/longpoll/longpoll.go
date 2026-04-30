// Package longpoll contains a managed long polling runner for Telegram updates.
package longpoll

import (
	"context"
	stderrors "errors"
	"reflect"
	"time"

	"github.com/xDilettante/ai-gram/bot"
	"github.com/xDilettante/ai-gram/telegram"
)

const (
	defaultBackoffInitial    = time.Second
	defaultBackoffMax        = 30 * time.Second
	defaultBackoffMultiplier = 2.0
)

// UpdateGetter fetches Telegram updates.
type UpdateGetter interface {
	GetUpdates(context.Context, bot.GetUpdatesParams) ([]telegram.Update, error)
}

// Handler handles one Telegram update.
type Handler interface {
	HandleUpdate(context.Context, telegram.Update) error
}

// HandlerFunc adapts a function to Handler.
type HandlerFunc func(context.Context, telegram.Update) error

// HandleUpdate calls f(ctx, update).
func (f HandlerFunc) HandleUpdate(ctx context.Context, update telegram.Update) error {
	return f(ctx, update)
}

// Config configures a Runner.
type Config struct {
	Offset         int64
	Limit          int
	Timeout        int
	AllowedUpdates []string

	Backoff BackoffConfig
	OnError func(context.Context, error)
}

// BackoffConfig configures retry delays after polling errors.
type BackoffConfig struct {
	Initial    time.Duration
	Max        time.Duration
	Multiplier float64
}

// Runner manages a long polling loop.
type Runner struct {
	getter  UpdateGetter
	handler Handler
	config  Config
	offset  int64
	backoff BackoffConfig
	sleep   func(context.Context, time.Duration) error
}

// New creates a Runner.
func New(getter UpdateGetter, handler Handler, config Config) (*Runner, error) {
	if isNil(getter) {
		return nil, stderrors.New("update getter is required")
	}
	if isNil(handler) {
		return nil, stderrors.New("handler is required")
	}
	if err := validateConfig(config); err != nil {
		return nil, err
	}

	return &Runner{
		getter:  getter,
		handler: handler,
		config:  config,
		offset:  config.Offset,
		backoff: normalizeBackoff(config.Backoff),
		sleep:   sleepContext,
	}, nil
}

// Run starts the long polling loop and blocks until ctx is canceled.
func (r *Runner) Run(ctx context.Context) error {
	if ctx == nil {
		return stderrors.New("context is required")
	}
	if r == nil {
		return stderrors.New("runner is required")
	}

	delay := r.backoff.Initial
	for {
		if err := ctx.Err(); err != nil {
			return err
		}

		updates, err := r.getter.GetUpdates(ctx, r.params())
		if err != nil {
			if ctxErr := ctx.Err(); ctxErr != nil {
				return ctxErr
			}
			r.reportError(ctx, err)
			if err := r.sleep(ctx, delay); err != nil {
				return err
			}
			delay = nextBackoffDelay(delay, r.backoff)
			continue
		}

		delay = r.backoff.Initial
		for _, update := range updates {
			if err := ctx.Err(); err != nil {
				return err
			}

			err := r.handler.HandleUpdate(ctx, update)
			r.advanceOffset(update.UpdateID)
			if err != nil {
				if ctxErr := ctx.Err(); ctxErr != nil {
					return ctxErr
				}
				r.reportError(ctx, err)
			}
		}
	}
}

func (r *Runner) params() bot.GetUpdatesParams {
	return bot.GetUpdatesParams{
		Offset:         r.offset,
		Limit:          r.config.Limit,
		Timeout:        r.config.Timeout,
		AllowedUpdates: r.config.AllowedUpdates,
	}
}

func (r *Runner) advanceOffset(updateID int64) {
	next := updateID + 1
	if next > r.offset {
		r.offset = next
	}
}

func (r *Runner) reportError(ctx context.Context, err error) {
	if r.config.OnError != nil {
		r.config.OnError(ctx, err)
	}
}

func isNil(value any) bool {
	if value == nil {
		return true
	}

	reflected := reflect.ValueOf(value)
	switch reflected.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Pointer, reflect.Slice:
		return reflected.IsNil()
	default:
		return false
	}
}

func validateConfig(config Config) error {
	if config.Limit < 0 {
		return stderrors.New("limit must not be negative")
	}
	if config.Limit > 100 {
		return stderrors.New("limit must be between 1 and 100")
	}
	if config.Timeout < 0 {
		return stderrors.New("timeout must not be negative")
	}

	return nil
}

func normalizeBackoff(config BackoffConfig) BackoffConfig {
	if config.Initial <= 0 {
		config.Initial = defaultBackoffInitial
	}
	if config.Max <= 0 {
		config.Max = defaultBackoffMax
	}
	if config.Multiplier < 1 {
		config.Multiplier = defaultBackoffMultiplier
	}
	if config.Initial > config.Max {
		config.Initial = config.Max
	}

	return config
}

func nextBackoffDelay(current time.Duration, config BackoffConfig) time.Duration {
	next := time.Duration(float64(current) * config.Multiplier)
	if next < current {
		return config.Max
	}
	if next > config.Max {
		return config.Max
	}

	return next
}

func sleepContext(ctx context.Context, delay time.Duration) error {
	if delay <= 0 {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			return nil
		}
	}

	timer := time.NewTimer(delay)
	defer timer.Stop()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}
