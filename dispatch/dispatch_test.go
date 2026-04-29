package dispatch

import (
	"context"
	stderrors "errors"
	"reflect"
	"testing"

	"ai-gram/telegram"
)

func TestNewCreatesDispatcher(t *testing.T) {
	dispatcher := New()
	if dispatcher == nil {
		t.Fatal("expected dispatcher")
	}
}

func TestDispatcherImplementsHandleUpdateInterface(t *testing.T) {
	var _ interface {
		HandleUpdate(context.Context, telegram.Update) error
	} = New()
}

func TestHandleValidation(t *testing.T) {
	dispatcher := New()
	handler := HandlerFunc(func(context.Context, telegram.Update) error { return nil })

	if err := dispatcher.Handle(nil, handler); err == nil {
		t.Fatal("expected error for nil predicate")
	}
	if err := dispatcher.Handle(Any(), nil); err == nil {
		t.Fatal("expected error for nil handler")
	}
	if err := dispatcher.Handle(Any(), HandlerFunc(nil)); err == nil {
		t.Fatal("expected error for typed nil handler")
	}
}

func TestRoutesAreCheckedInRegistrationOrderAndFirstMatchWins(t *testing.T) {
	dispatcher := New()
	var calls []string

	must(t, dispatcher.HandleFunc(Any(), func(context.Context, telegram.Update) error {
		calls = append(calls, "first")
		return nil
	}))
	must(t, dispatcher.HandleFunc(Any(), func(context.Context, telegram.Update) error {
		calls = append(calls, "second")
		return nil
	}))

	if err := dispatcher.HandleUpdate(context.Background(), telegram.Update{UpdateID: 1}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got, want := calls, []string{"first"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("unexpected calls: got %v, want %v", got, want)
	}
}

func TestFallbackIsCalledWhenNoRouteMatches(t *testing.T) {
	var fallbackCalled bool
	dispatcher := New(WithFallback(HandlerFunc(func(context.Context, telegram.Update) error {
		fallbackCalled = true
		return nil
	})))
	must(t, dispatcher.OnMessageFunc(func(context.Context, telegram.Update) error {
		t.Fatal("message route should not be called")
		return nil
	}))

	if err := dispatcher.HandleUpdate(context.Background(), callbackUpdate("x")); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !fallbackCalled {
		t.Fatal("expected fallback")
	}
}

func TestNoMatchWithoutFallbackReturnsNil(t *testing.T) {
	dispatcher := New()
	must(t, dispatcher.OnMessageFunc(func(context.Context, telegram.Update) error {
		t.Fatal("message route should not be called")
		return nil
	}))

	if err := dispatcher.HandleUpdate(context.Background(), callbackUpdate("x")); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestOnMessageMatchesOnlyMessageUpdates(t *testing.T) {
	dispatcher := New()
	var calls int
	must(t, dispatcher.OnMessageFunc(func(context.Context, telegram.Update) error {
		calls++
		return nil
	}))

	if err := dispatcher.HandleUpdate(context.Background(), messageUpdate("hello")); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := dispatcher.HandleUpdate(context.Background(), callbackUpdate("x")); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if calls != 1 {
		t.Fatalf("unexpected calls: %d", calls)
	}
}

func TestOnCommandMatching(t *testing.T) {
	tests := []struct {
		name string
		text string
		want bool
	}{
		{name: "exact", text: "/start", want: true},
		{name: "payload", text: "/start payload", want: true},
		{name: "newline payload", text: "/start\npayload", want: true},
		{name: "prefix only", text: "/startx", want: false},
		{name: "plain text", text: "start", want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dispatcher := New()
			var calls int
			must(t, dispatcher.OnCommandFunc("start", func(context.Context, telegram.Update) error {
				calls++
				return nil
			}))

			if err := dispatcher.HandleUpdate(context.Background(), messageUpdate(tt.text)); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			got := calls == 1
			if got != tt.want {
				t.Fatalf("match = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOnCommandValidation(t *testing.T) {
	dispatcher := New()
	handler := HandlerFunc(func(context.Context, telegram.Update) error { return nil })

	for _, command := range []string{"", "/start", "start now", "start\tnow"} {
		if err := dispatcher.OnCommand(command, handler); err == nil {
			t.Fatalf("expected error for command %q", command)
		}
	}
}

func TestCallbackRoutes(t *testing.T) {
	dispatcher := New()
	var queryCalls int
	var dataCalls int
	must(t, dispatcher.OnCallbackQueryFunc(func(context.Context, telegram.Update) error {
		queryCalls++
		return nil
	}))
	must(t, dispatcher.OnCallbackDataFunc("x", func(context.Context, telegram.Update) error {
		dataCalls++
		return nil
	}))

	if err := dispatcher.HandleUpdate(context.Background(), callbackUpdate("x")); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if queryCalls != 1 || dataCalls != 0 {
		t.Fatalf("first matching route should win: query=%d data=%d", queryCalls, dataCalls)
	}

	dispatcher = New()
	must(t, dispatcher.OnCallbackDataFunc("x", func(context.Context, telegram.Update) error {
		dataCalls++
		return nil
	}))
	if err := dispatcher.HandleUpdate(context.Background(), callbackUpdate("y")); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if dataCalls != 0 {
		t.Fatalf("callback data route matched wrong data: %d", dataCalls)
	}

	if err := dispatcher.OnCallbackData("", HandlerFunc(func(context.Context, telegram.Update) error { return nil })); err == nil {
		t.Fatal("expected error for empty callback data")
	}
}

func TestMiddlewareAppliesToRouteInOrder(t *testing.T) {
	dispatcher := New()
	var calls []string
	dispatcher.Use(nil, middlewareRecorder("mw1", &calls), middlewareRecorder("mw2", &calls))
	must(t, dispatcher.OnMessageFunc(func(context.Context, telegram.Update) error {
		calls = append(calls, "handler")
		return nil
	}))

	if err := dispatcher.HandleUpdate(context.Background(), messageUpdate("hello")); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := []string{"mw1 before", "mw2 before", "handler", "mw2 after", "mw1 after"}
	if !reflect.DeepEqual(calls, want) {
		t.Fatalf("unexpected middleware order: got %v, want %v", calls, want)
	}
}

func TestMiddlewareAppliesToFallback(t *testing.T) {
	var calls []string
	dispatcher := New(WithFallback(HandlerFunc(func(context.Context, telegram.Update) error {
		calls = append(calls, "fallback")
		return nil
	})))
	dispatcher.Use(middlewareRecorder("mw1", &calls))

	if err := dispatcher.HandleUpdate(context.Background(), callbackUpdate("x")); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := []string{"mw1 before", "fallback", "mw1 after"}
	if !reflect.DeepEqual(calls, want) {
		t.Fatalf("unexpected calls: got %v, want %v", calls, want)
	}
}

func TestMiddlewareDoesNotAccumulateBetweenUpdates(t *testing.T) {
	dispatcher := New()
	var calls int
	dispatcher.Use(func(next Handler) Handler {
		return HandlerFunc(func(ctx context.Context, update telegram.Update) error {
			calls++
			return next.HandleUpdate(ctx, update)
		})
	})
	must(t, dispatcher.OnMessageFunc(func(context.Context, telegram.Update) error { return nil }))

	for i := 0; i < 2; i++ {
		if err := dispatcher.HandleUpdate(context.Background(), messageUpdate("hello")); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	}
	if calls != 2 {
		t.Fatalf("middleware should run once per update, got %d", calls)
	}
}

func TestHandlerErrorWithoutErrorHandlerIsReturned(t *testing.T) {
	want := stderrors.New("handler failed")
	dispatcher := New()
	must(t, dispatcher.OnMessageFunc(func(context.Context, telegram.Update) error { return want }))

	err := dispatcher.HandleUpdate(context.Background(), messageUpdate("hello"))
	if !stderrors.Is(err, want) {
		t.Fatalf("got %v, want %v", err, want)
	}
}

func TestErrorHandlerReceivesRouteError(t *testing.T) {
	want := stderrors.New("handler failed")
	var gotErr error
	var gotUpdate telegram.Update
	dispatcher := New(WithErrorHandler(func(ctx context.Context, update telegram.Update, err error) {
		gotErr = err
		gotUpdate = update
	}))
	must(t, dispatcher.OnMessageFunc(func(context.Context, telegram.Update) error { return want }))

	update := messageUpdate("hello")
	err := dispatcher.HandleUpdate(context.Background(), update)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !stderrors.Is(gotErr, want) {
		t.Fatalf("got error %v, want %v", gotErr, want)
	}
	if gotUpdate.UpdateID != update.UpdateID {
		t.Fatalf("unexpected update: %+v", gotUpdate)
	}
}

func TestErrorHandlerReceivesFallbackError(t *testing.T) {
	want := stderrors.New("fallback failed")
	var gotErr error
	dispatcher := New(
		WithFallback(HandlerFunc(func(context.Context, telegram.Update) error { return want })),
		WithErrorHandler(func(ctx context.Context, update telegram.Update, err error) { gotErr = err }),
	)

	err := dispatcher.HandleUpdate(context.Background(), callbackUpdate("x"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !stderrors.Is(gotErr, want) {
		t.Fatalf("got error %v, want %v", gotErr, want)
	}
}

func TestHandlersAndMiddlewareReceiveSameContext(t *testing.T) {
	type contextKey string
	const key contextKey = "key"
	ctx := context.WithValue(context.Background(), key, "value")
	dispatcher := New()
	var middlewareSawValue bool
	var handlerSawValue bool
	dispatcher.Use(func(next Handler) Handler {
		return HandlerFunc(func(ctx context.Context, update telegram.Update) error {
			middlewareSawValue = ctx.Value(key) == "value"
			return next.HandleUpdate(ctx, update)
		})
	})
	must(t, dispatcher.OnMessageFunc(func(ctx context.Context, update telegram.Update) error {
		handlerSawValue = ctx.Value(key) == "value"
		return nil
	}))

	if err := dispatcher.HandleUpdate(ctx, messageUpdate("hello")); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !middlewareSawValue || !handlerSawValue {
		t.Fatalf("context was not propagated: middleware=%v handler=%v", middlewareSawValue, handlerSawValue)
	}
}

func TestCanceledContextDoesNotCallMiddleware(t *testing.T) {
	dispatcher := New()
	var middlewareCalls int
	dispatcher.Use(func(next Handler) Handler {
		return HandlerFunc(func(ctx context.Context, update telegram.Update) error {
			middlewareCalls++
			return next.HandleUpdate(ctx, update)
		})
	})
	must(t, dispatcher.OnMessageFunc(func(context.Context, telegram.Update) error { return nil }))
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := dispatcher.HandleUpdate(ctx, messageUpdate("hello"))
	if !stderrors.Is(err, context.Canceled) {
		t.Fatalf("expected context.Canceled, got %v", err)
	}
	if middlewareCalls != 0 {
		t.Fatalf("middleware should not be called, calls=%d", middlewareCalls)
	}
}

func TestCanceledContextDoesNotCallRoutes(t *testing.T) {
	dispatcher := New()
	var calls int
	must(t, dispatcher.OnMessageFunc(func(context.Context, telegram.Update) error {
		calls++
		return nil
	}))
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := dispatcher.HandleUpdate(ctx, messageUpdate("hello"))
	if !stderrors.Is(err, context.Canceled) {
		t.Fatalf("expected context.Canceled, got %v", err)
	}
	if calls != 0 {
		t.Fatalf("route should not be called, calls=%d", calls)
	}
}

func TestPredicateHelpersAreSafeForInvalidInputs(t *testing.T) {
	update := messageUpdate("/start")
	if Command("")(update) {
		t.Fatal("empty command should not match")
	}
	if Command("/start")(update) {
		t.Fatal("slash command should not match")
	}
	if Command("start now")(update) {
		t.Fatal("space command should not match")
	}
	if CallbackData("")(callbackUpdate("")) {
		t.Fatal("empty callback data should not match")
	}
}

func TestChainSkipsNilMiddleware(t *testing.T) {
	var calls []string
	handler := HandlerFunc(func(context.Context, telegram.Update) error {
		calls = append(calls, "handler")
		return nil
	})
	chained := Chain(handler, nil, middlewareRecorder("mw", &calls))

	if err := chained.HandleUpdate(context.Background(), telegram.Update{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := []string{"mw before", "handler", "mw after"}
	if !reflect.DeepEqual(calls, want) {
		t.Fatalf("unexpected calls: got %v, want %v", calls, want)
	}
}

func must(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func middlewareRecorder(name string, calls *[]string) Middleware {
	return func(next Handler) Handler {
		return HandlerFunc(func(ctx context.Context, update telegram.Update) error {
			*calls = append(*calls, name+" before")
			err := next.HandleUpdate(ctx, update)
			*calls = append(*calls, name+" after")
			return err
		})
	}
}

func messageUpdate(text string) telegram.Update {
	return telegram.Update{UpdateID: 1, Message: &telegram.Message{Text: text}}
}

func callbackUpdate(data string) telegram.Update {
	return telegram.Update{UpdateID: 2, CallbackQuery: &telegram.CallbackQuery{ID: "cb", Data: data}}
}
