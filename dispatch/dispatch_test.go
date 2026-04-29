package dispatch

import (
	"context"
	"reflect"
	"testing"

	"ai-gram/telegram"
)

func TestMiddlewareCanBeComposedAroundHandler(t *testing.T) {
	var calls []string

	handler := HandlerFunc(func(context.Context, telegram.Update) error {
		calls = append(calls, "handler")
		return nil
	})

	first := MiddlewareFunc(func(next Handler) Handler {
		return HandlerFunc(func(ctx context.Context, update telegram.Update) error {
			calls = append(calls, "first:before")
			err := next.HandleUpdate(ctx, update)
			calls = append(calls, "first:after")
			return err
		})
	})

	second := MiddlewareFunc(func(next Handler) Handler {
		return HandlerFunc(func(ctx context.Context, update telegram.Update) error {
			calls = append(calls, "second:before")
			err := next.HandleUpdate(ctx, update)
			calls = append(calls, "second:after")
			return err
		})
	})

	composed := Chain(handler, first, second)
	if err := composed.HandleUpdate(context.Background(), telegram.Update{UpdateID: 1}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := []string{"first:before", "second:before", "handler", "second:after", "first:after"}
	if !reflect.DeepEqual(calls, want) {
		t.Fatalf("unexpected call order: got %v, want %v", calls, want)
	}
}
