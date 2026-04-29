// Package middleware contains ready-to-use middleware helpers for dispatch handlers.
package middleware

import (
	"context"
	stderrors "errors"
	"time"

	"ai-gram/dispatch"
	"ai-gram/telegram"
)

// PanicError reports a panic recovered from a user handler.
type PanicError struct {
	// Value is the recovered panic value. It is not included in Error output.
	Value any
}

// Error returns a redacted panic error message.
func (e *PanicError) Error() string {
	return "panic recovered while handling update"
}

// Observer receives hook callbacks around update handling.
type Observer interface {
	// OnUpdateStart is called before the next handler runs.
	OnUpdateStart(context.Context, telegram.Update)
	// OnUpdateFinish is called after the next handler returns.
	OnUpdateFinish(context.Context, telegram.Update, error, time.Duration)
}

// Recover catches panics from the next handler and returns a PanicError.
func Recover(onPanic func(context.Context, telegram.Update, any)) dispatch.Middleware {
	return func(next dispatch.Handler) dispatch.Handler {
		if next == nil {
			return dispatch.HandlerFunc(func(context.Context, telegram.Update) error {
				return stderrors.New("handler is required")
			})
		}

		return dispatch.HandlerFunc(func(ctx context.Context, update telegram.Update) (err error) {
			defer func() {
				if recovered := recover(); recovered != nil {
					callOnPanic(onPanic, ctx, update, recovered)
					err = &PanicError{Value: recovered}
				}
			}()

			return next.HandleUpdate(ctx, update)
		})
	}
}

// Timeout passes a child context with timeout to the next handler.
func Timeout(timeout time.Duration) dispatch.Middleware {
	return func(next dispatch.Handler) dispatch.Handler {
		if timeout <= 0 {
			return next
		}
		if next == nil {
			return dispatch.HandlerFunc(func(context.Context, telegram.Update) error {
				return stderrors.New("handler is required")
			})
		}

		return dispatch.HandlerFunc(func(ctx context.Context, update telegram.Update) error {
			childCtx, cancel := context.WithTimeout(ctx, timeout)
			defer cancel()

			return next.HandleUpdate(childCtx, update)
		})
	}
}

// Observe reports update handling start and finish events to observer.
func Observe(observer Observer) dispatch.Middleware {
	return func(next dispatch.Handler) dispatch.Handler {
		if observer == nil {
			return next
		}
		if next == nil {
			return dispatch.HandlerFunc(func(context.Context, telegram.Update) error {
				return stderrors.New("handler is required")
			})
		}

		return dispatch.HandlerFunc(func(ctx context.Context, update telegram.Update) error {
			startedAt := time.Now()
			observer.OnUpdateStart(ctx, update)
			err := next.HandleUpdate(ctx, update)
			observer.OnUpdateFinish(ctx, update, err, time.Since(startedAt))

			return err
		})
	}
}

func callOnPanic(onPanic func(context.Context, telegram.Update, any), ctx context.Context, update telegram.Update, recovered any) {
	if onPanic == nil {
		return
	}
	defer func() {
		_ = recover()
	}()

	onPanic(ctx, update, recovered)
}
