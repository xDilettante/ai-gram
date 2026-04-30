package longpoll

import (
	"context"
	stderrors "errors"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/xDilettante/ai-gram/bot"
	"github.com/xDilettante/ai-gram/telegram"
)

func TestNewValidation(t *testing.T) {
	getter := &fakeGetter{}
	handler := HandlerFunc(func(context.Context, telegram.Update) error { return nil })

	tests := []struct {
		name    string
		getter  UpdateGetter
		handler Handler
		config  Config
	}{
		{name: "nil getter", handler: handler},
		{name: "nil handler", getter: getter},
		{name: "typed nil handler", getter: getter, handler: HandlerFunc(nil)},
		{name: "negative limit", getter: getter, handler: handler, config: Config{Limit: -1}},
		{name: "too high limit", getter: getter, handler: handler, config: Config{Limit: 101}},
		{name: "negative timeout", getter: getter, handler: handler, config: Config{Timeout: -1}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			runner, err := New(tt.getter, tt.handler, tt.config)
			if err == nil {
				t.Fatal("expected error")
			}
			if runner != nil {
				t.Fatalf("expected nil runner, got %+v", runner)
			}
			assertNoToken(t, err, "123:secret")
		})
	}
}

func TestNewValidConfigCreatesRunner(t *testing.T) {
	getter := &fakeGetter{}
	handler := HandlerFunc(func(context.Context, telegram.Update) error { return nil })

	runner, err := New(getter, handler, Config{
		Offset:         10,
		Limit:          100,
		Timeout:        30,
		AllowedUpdates: []string{"message"},
		Backoff: BackoffConfig{
			Initial:    -1,
			Max:        -1,
			Multiplier: 0,
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if runner == nil {
		t.Fatal("expected runner")
	}
	if runner.offset != 10 {
		t.Fatalf("unexpected offset: %d", runner.offset)
	}
	if runner.backoff.Initial != defaultBackoffInitial || runner.backoff.Max != defaultBackoffMax || runner.backoff.Multiplier != defaultBackoffMultiplier {
		t.Fatalf("unexpected default backoff: %+v", runner.backoff)
	}
}

func TestRunSuccessHandlesUpdatesAndAdvancesOffset(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	getter := &fakeGetter{responses: []fakeGetUpdatesResponse{
		{updates: []telegram.Update{{UpdateID: 10}, {UpdateID: 11}, {UpdateID: 12}}},
		{updates: []telegram.Update{}, after: cancel},
	}}
	var handled []int64
	runner := newRunnerForTest(t, getter, HandlerFunc(func(ctx context.Context, update telegram.Update) error {
		handled = append(handled, update.UpdateID)
		return nil
	}), Config{
		Offset:         5,
		Limit:          50,
		Timeout:        20,
		AllowedUpdates: []string{"message", "callback_query"},
	})

	err := runner.Run(ctx)
	if !stderrors.Is(err, context.Canceled) {
		t.Fatalf("expected context.Canceled, got %v", err)
	}
	if got, want := handled, []int64{10, 11, 12}; !reflect.DeepEqual(got, want) {
		t.Fatalf("unexpected handled updates: got %v, want %v", got, want)
	}
	if len(getter.calls) != 2 {
		t.Fatalf("unexpected call count: %d", len(getter.calls))
	}
	first := getter.calls[0]
	if first.Offset != 5 || first.Limit != 50 || first.Timeout != 20 || !reflect.DeepEqual(first.AllowedUpdates, []string{"message", "callback_query"}) {
		t.Fatalf("unexpected first params: %+v", first)
	}
	if got := getter.calls[1].Offset; got != 13 {
		t.Fatalf("unexpected second offset: %d", got)
	}
}

func TestRunEmptyUpdatesContinuesWithoutHandler(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	getter := &fakeGetter{responses: []fakeGetUpdatesResponse{
		{updates: []telegram.Update{}},
		{updates: []telegram.Update{}, after: cancel},
	}}
	var handlerCalls int
	runner := newRunnerForTest(t, getter, HandlerFunc(func(ctx context.Context, update telegram.Update) error {
		handlerCalls++
		return nil
	}), Config{})

	err := runner.Run(ctx)
	if !stderrors.Is(err, context.Canceled) {
		t.Fatalf("expected context.Canceled, got %v", err)
	}
	if handlerCalls != 0 {
		t.Fatalf("unexpected handler calls: %d", handlerCalls)
	}
	if len(getter.calls) != 2 {
		t.Fatalf("expected loop to continue, call count %d", len(getter.calls))
	}
}

func TestRunCanceledContextReturnsContextError(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	getter := &fakeGetter{}
	runner := newRunnerForTest(t, getter, HandlerFunc(func(context.Context, telegram.Update) error { return nil }), Config{})

	err := runner.Run(ctx)
	if !stderrors.Is(err, context.Canceled) {
		t.Fatalf("expected context.Canceled, got %v", err)
	}
	if len(getter.calls) != 0 {
		t.Fatalf("unexpected getter calls: %d", len(getter.calls))
	}
}

func TestRunCanceledDuringBackoffReturnsContextError(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	getter := &fakeGetter{responses: []fakeGetUpdatesResponse{{err: stderrors.New("temporary failure")}}}
	runner := newRunnerForTest(t, getter, HandlerFunc(func(context.Context, telegram.Update) error { return nil }), Config{})
	runner.sleep = func(ctx context.Context, delay time.Duration) error {
		cancel()
		<-ctx.Done()
		return ctx.Err()
	}

	err := runner.Run(ctx)
	if !stderrors.Is(err, context.Canceled) {
		t.Fatalf("expected context.Canceled, got %v", err)
	}
}

func TestRunGetUpdatesErrorBacksOffAndContinues(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	getter := &fakeGetter{responses: []fakeGetUpdatesResponse{
		{err: stderrors.New("temporary failure")},
		{updates: []telegram.Update{{UpdateID: 1}}, after: cancel},
	}}
	var errorsSeen []error
	var delays []time.Duration
	runner := newRunnerForTest(t, getter, HandlerFunc(func(context.Context, telegram.Update) error { return nil }), Config{
		Backoff: BackoffConfig{Initial: time.Millisecond, Max: 10 * time.Millisecond, Multiplier: 2},
		OnError: func(ctx context.Context, err error) {
			errorsSeen = append(errorsSeen, err)
		},
	})
	runner.sleep = func(ctx context.Context, delay time.Duration) error {
		delays = append(delays, delay)
		return nil
	}

	err := runner.Run(ctx)
	if !stderrors.Is(err, context.Canceled) {
		t.Fatalf("expected context.Canceled, got %v", err)
	}
	if len(errorsSeen) != 1 {
		t.Fatalf("unexpected OnError calls: %d", len(errorsSeen))
	}
	if got, want := delays, []time.Duration{time.Millisecond}; !reflect.DeepEqual(got, want) {
		t.Fatalf("unexpected delays: got %v, want %v", got, want)
	}
	if len(getter.calls) != 2 {
		t.Fatalf("expected retry, call count %d", len(getter.calls))
	}
}

func TestRunBackoffResetsAfterSuccess(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	getter := &fakeGetter{responses: []fakeGetUpdatesResponse{
		{err: stderrors.New("first failure")},
		{updates: []telegram.Update{}},
		{err: stderrors.New("second failure")},
	}}
	var delays []time.Duration
	runner := newRunnerForTest(t, getter, HandlerFunc(func(context.Context, telegram.Update) error { return nil }), Config{
		Backoff: BackoffConfig{Initial: time.Millisecond, Max: 10 * time.Millisecond, Multiplier: 2},
	})
	runner.sleep = func(ctx context.Context, delay time.Duration) error {
		delays = append(delays, delay)
		if len(delays) == 2 {
			cancel()
			return ctx.Err()
		}
		return nil
	}

	err := runner.Run(ctx)
	if !stderrors.Is(err, context.Canceled) {
		t.Fatalf("expected context.Canceled, got %v", err)
	}
	if got, want := delays, []time.Duration{time.Millisecond, time.Millisecond}; !reflect.DeepEqual(got, want) {
		t.Fatalf("expected backoff reset after success: got %v, want %v", got, want)
	}
}

func TestRunHandlerErrorReportsAndContinues(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	getter := &fakeGetter{responses: []fakeGetUpdatesResponse{
		{updates: []telegram.Update{{UpdateID: 10}, {UpdateID: 11}}},
		{updates: []telegram.Update{}, after: cancel},
	}}
	var handled []int64
	var errorsSeen []error
	runner := newRunnerForTest(t, getter, HandlerFunc(func(ctx context.Context, update telegram.Update) error {
		handled = append(handled, update.UpdateID)
		if update.UpdateID == 10 {
			return stderrors.New("handler failed")
		}
		return nil
	}), Config{OnError: func(ctx context.Context, err error) {
		errorsSeen = append(errorsSeen, err)
	}})

	err := runner.Run(ctx)
	if !stderrors.Is(err, context.Canceled) {
		t.Fatalf("expected context.Canceled, got %v", err)
	}
	if got, want := handled, []int64{10, 11}; !reflect.DeepEqual(got, want) {
		t.Fatalf("unexpected handled updates: got %v, want %v", got, want)
	}
	if len(errorsSeen) != 1 {
		t.Fatalf("unexpected OnError calls: %d", len(errorsSeen))
	}
	if got := getter.calls[1].Offset; got != 12 {
		t.Fatalf("handler error should not keep old offset, got %d", got)
	}
}

func TestRunOffsetMonotonic(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	getter := &fakeGetter{responses: []fakeGetUpdatesResponse{
		{updates: []telegram.Update{{UpdateID: 10}, {UpdateID: 19}}},
		{updates: []telegram.Update{}, after: cancel},
	}}
	runner := newRunnerForTest(t, getter, HandlerFunc(func(context.Context, telegram.Update) error { return nil }), Config{Offset: 20})

	err := runner.Run(ctx)
	if !stderrors.Is(err, context.Canceled) {
		t.Fatalf("expected context.Canceled, got %v", err)
	}
	if got := getter.calls[1].Offset; got != 20 {
		t.Fatalf("offset decreased: %d", got)
	}
}

func TestLongpollErrorsDoNotContainToken(t *testing.T) {
	const token = "123:secret"
	_, err := New(nil, nil, Config{Limit: -1, Timeout: -1})
	if err == nil {
		t.Fatal("expected error")
	}
	if strings.Contains(err.Error(), token) {
		t.Fatalf("error leaked token: %q", err.Error())
	}
}

func assertNoToken(t *testing.T, err error, token string) {
	t.Helper()

	if strings.Contains(err.Error(), token) {
		t.Fatalf("error leaked token: %q", err.Error())
	}
}

func newRunnerForTest(t *testing.T, getter UpdateGetter, handler Handler, config Config) *Runner {
	t.Helper()

	runner, err := New(getter, handler, config)
	if err != nil {
		t.Fatalf("unexpected New error: %v", err)
	}
	runner.sleep = func(context.Context, time.Duration) error { return nil }

	return runner
}

type fakeGetUpdatesResponse struct {
	updates []telegram.Update
	err     error
	after   func()
}

type fakeGetter struct {
	calls     []bot.GetUpdatesParams
	responses []fakeGetUpdatesResponse
}

func (g *fakeGetter) GetUpdates(ctx context.Context, params bot.GetUpdatesParams) ([]telegram.Update, error) {
	g.calls = append(g.calls, params)
	if len(g.responses) == 0 {
		return nil, stderrors.New("unexpected GetUpdates call")
	}

	response := g.responses[0]
	g.responses = g.responses[1:]
	if response.after != nil {
		defer response.after()
	}

	return response.updates, response.err
}
