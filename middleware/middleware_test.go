package middleware

import (
	"context"
	stderrors "errors"
	"reflect"
	"strings"
	"testing"
	"time"

	"ai-gram/dispatch"
	"ai-gram/telegram"
)

func TestRecoverConvertsPanicToPanicError(t *testing.T) {
	const recoveredValue = "boom"
	update := tokenUpdate("123:secret")
	var gotUpdate telegram.Update
	var gotRecovered any

	handler := Recover(func(ctx context.Context, update telegram.Update, recovered any) {
		gotUpdate = update
		gotRecovered = recovered
	})(dispatch.HandlerFunc(func(context.Context, telegram.Update) error {
		panic(recoveredValue)
	}))

	err := handler.HandleUpdate(context.Background(), update)
	if err == nil {
		t.Fatal("expected error")
	}
	var panicErr *PanicError
	if !stderrors.As(err, &panicErr) {
		t.Fatalf("expected PanicError, got %T", err)
	}
	if panicErr.Value != recoveredValue {
		t.Fatalf("unexpected recovered value: %#v", panicErr.Value)
	}
	if gotUpdate.UpdateID != update.UpdateID || gotRecovered != recoveredValue {
		t.Fatalf("unexpected onPanic args: update=%+v recovered=%#v", gotUpdate, gotRecovered)
	}
	if strings.Contains(err.Error(), "123:secret") {
		t.Fatalf("panic error leaked token: %q", err.Error())
	}
}

func TestRecoverDoesNotSuppressHandlerErrorOrSuccess(t *testing.T) {
	want := stderrors.New("handler failed")
	handler := Recover(nil)(dispatch.HandlerFunc(func(context.Context, telegram.Update) error {
		return want
	}))
	if err := handler.HandleUpdate(context.Background(), telegram.Update{}); !stderrors.Is(err, want) {
		t.Fatalf("got %v, want %v", err, want)
	}

	handler = Recover(nil)(dispatch.HandlerFunc(func(context.Context, telegram.Update) error {
		return nil
	}))
	if err := handler.HandleUpdate(context.Background(), telegram.Update{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRecoverDoesNotPanicWhenOnPanicPanics(t *testing.T) {
	handler := Recover(func(context.Context, telegram.Update, any) {
		panic("observer failed")
	})(dispatch.HandlerFunc(func(context.Context, telegram.Update) error {
		panic("handler failed")
	}))

	var panicErr *PanicError
	if err := handler.HandleUpdate(context.Background(), telegram.Update{}); !stderrors.As(err, &panicErr) {
		t.Fatalf("expected PanicError, got %v", err)
	}
}

func TestPanicErrorErrorIsNonEmptyAndRedacted(t *testing.T) {
	panicErr := &PanicError{Value: "123:secret"}
	if panicErr.Error() == "" {
		t.Fatal("expected non-empty error")
	}
	if strings.Contains(panicErr.Error(), "123:secret") {
		t.Fatalf("panic error leaked value: %q", panicErr.Error())
	}
}

func TestTimeoutNoopForNonPositiveTimeout(t *testing.T) {
	var gotCtx context.Context
	parent := context.Background()
	handler := Timeout(0)(dispatch.HandlerFunc(func(ctx context.Context, update telegram.Update) error {
		gotCtx = ctx
		return nil
	}))

	if err := handler.HandleUpdate(parent, telegram.Update{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotCtx != parent {
		t.Fatal("expected no-op timeout to pass original context")
	}
}

func TestTimeoutPassesContextWithDeadline(t *testing.T) {
	var gotDeadline bool
	handler := Timeout(time.Second)(dispatch.HandlerFunc(func(ctx context.Context, update telegram.Update) error {
		_, gotDeadline = ctx.Deadline()
		return nil
	}))

	if err := handler.HandleUpdate(context.Background(), telegram.Update{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !gotDeadline {
		t.Fatal("expected handler context deadline")
	}
}

func TestTimeoutReturnsHandlerError(t *testing.T) {
	want := stderrors.New("handler failed")
	handler := Timeout(time.Second)(dispatch.HandlerFunc(func(ctx context.Context, update telegram.Update) error {
		return want
	}))

	if err := handler.HandleUpdate(context.Background(), telegram.Update{}); !stderrors.Is(err, want) {
		t.Fatalf("got %v, want %v", err, want)
	}
}

func TestTimeoutAllowsHandlerToReturnDeadlineExceeded(t *testing.T) {
	handler := Timeout(time.Nanosecond)(dispatch.HandlerFunc(func(ctx context.Context, update telegram.Update) error {
		<-ctx.Done()
		return ctx.Err()
	}))

	err := handler.HandleUpdate(context.Background(), telegram.Update{})
	if !stderrors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("expected context.DeadlineExceeded, got %v", err)
	}
}

func TestObserveNilObserverIsNoop(t *testing.T) {
	var called bool
	handler := Observe(nil)(dispatch.HandlerFunc(func(ctx context.Context, update telegram.Update) error {
		called = true
		return nil
	}))

	if err := handler.HandleUpdate(context.Background(), telegram.Update{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !called {
		t.Fatal("expected handler call")
	}
}

func TestObserveReportsStartFinishDurationAndError(t *testing.T) {
	want := stderrors.New("handler failed")
	observer := &recordingObserver{}
	handler := Observe(observer)(dispatch.HandlerFunc(func(ctx context.Context, update telegram.Update) error {
		observer.events = append(observer.events, "handler")
		return want
	}))

	err := handler.HandleUpdate(context.Background(), telegram.Update{UpdateID: 7})
	if !stderrors.Is(err, want) {
		t.Fatalf("got %v, want %v", err, want)
	}
	if got, wantEvents := observer.events, []string{"start", "handler", "finish"}; !reflect.DeepEqual(got, wantEvents) {
		t.Fatalf("unexpected events: got %v, want %v", got, wantEvents)
	}
	if !stderrors.Is(observer.finishErr, want) {
		t.Fatalf("finish error = %v, want %v", observer.finishErr, want)
	}
	if observer.duration < 0 {
		t.Fatalf("duration must be >= 0, got %v", observer.duration)
	}
	if observer.startUpdate.UpdateID != 7 || observer.finishUpdate.UpdateID != 7 {
		t.Fatalf("unexpected updates: start=%+v finish=%+v", observer.startUpdate, observer.finishUpdate)
	}
}

func TestMiddlewareIntegrationWithDispatch(t *testing.T) {
	observer := &recordingObserver{}
	dispatcher := dispatch.New()
	dispatcher.Use(
		Recover(nil),
		Timeout(time.Second),
		Observe(observer),
	)
	if err := dispatcher.OnMessageFunc(func(ctx context.Context, update telegram.Update) error {
		if _, ok := ctx.Deadline(); !ok {
			t.Fatal("expected deadline")
		}
		panic("boom")
	}); err != nil {
		t.Fatalf("unexpected route error: %v", err)
	}

	err := dispatcher.HandleUpdate(context.Background(), tokenUpdate("hello"))
	var panicErr *PanicError
	if !stderrors.As(err, &panicErr) {
		t.Fatalf("expected PanicError, got %v", err)
	}
	if got, want := observer.events, []string{"start"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("observer should see start before recovered panic, got %v", got)
	}
}

func TestMiddlewareOrderWithDispatch(t *testing.T) {
	var calls []string
	dispatcher := dispatch.New()
	dispatcher.Use(
		recordingMiddleware("mw1", &calls),
		recordingMiddleware("mw2", &calls),
		Observe(&eventObserver{events: &calls}),
	)
	if err := dispatcher.OnMessageFunc(func(ctx context.Context, update telegram.Update) error {
		calls = append(calls, "handler")
		return nil
	}); err != nil {
		t.Fatalf("unexpected route error: %v", err)
	}

	if err := dispatcher.HandleUpdate(context.Background(), tokenUpdate("hello")); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := []string{"mw1 before", "mw2 before", "start", "handler", "finish", "mw2 after", "mw1 after"}
	if !reflect.DeepEqual(calls, want) {
		t.Fatalf("unexpected order: got %v, want %v", calls, want)
	}
}

func tokenUpdate(text string) telegram.Update {
	return telegram.Update{UpdateID: 1, Message: &telegram.Message{Text: text}}
}

type recordingObserver struct {
	events       []string
	startUpdate  telegram.Update
	finishUpdate telegram.Update
	finishErr    error
	duration     time.Duration
}

func (o *recordingObserver) OnUpdateStart(ctx context.Context, update telegram.Update) {
	o.events = append(o.events, "start")
	o.startUpdate = update
}

func (o *recordingObserver) OnUpdateFinish(ctx context.Context, update telegram.Update, err error, duration time.Duration) {
	o.events = append(o.events, "finish")
	o.finishUpdate = update
	o.finishErr = err
	o.duration = duration
}

type eventObserver struct {
	events *[]string
}

func (o *eventObserver) OnUpdateStart(context.Context, telegram.Update) {
	*o.events = append(*o.events, "start")
}

func (o *eventObserver) OnUpdateFinish(context.Context, telegram.Update, error, time.Duration) {
	*o.events = append(*o.events, "finish")
}

func recordingMiddleware(name string, calls *[]string) dispatch.Middleware {
	return func(next dispatch.Handler) dispatch.Handler {
		return dispatch.HandlerFunc(func(ctx context.Context, update telegram.Update) error {
			*calls = append(*calls, name+" before")
			err := next.HandleUpdate(ctx, update)
			*calls = append(*calls, name+" after")
			return err
		})
	}
}
