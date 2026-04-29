// Package dispatch defines the update handling contracts used by routers and middleware.
package dispatch

import (
	"context"

	"ai-gram/telegram"
)

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

// Middleware wraps a Handler with additional behavior.
type Middleware interface {
	Wrap(Handler) Handler
}

// MiddlewareFunc adapts a function to Middleware.
type MiddlewareFunc func(Handler) Handler

// Wrap calls f(next).
func (f MiddlewareFunc) Wrap(next Handler) Handler {
	return f(next)
}

// Dispatcher dispatches updates to handlers.
type Dispatcher interface {
	Dispatch(context.Context, telegram.Update) error
}

// Chain wraps handler with middleware in the order it is provided.
func Chain(handler Handler, middleware ...Middleware) Handler {
	for i := len(middleware) - 1; i >= 0; i-- {
		handler = middleware[i].Wrap(handler)
	}

	return handler
}
