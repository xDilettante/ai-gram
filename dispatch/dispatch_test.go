package dispatch

import (
	"context"
	stderrors "errors"
	"reflect"
	"testing"

	"github.com/xDilettante/ai-gram/telegram"
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
		{name: "bot username payload", text: "/start@BotName payload", want: true},
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

func TestInlineQueryRoutes(t *testing.T) {
	dispatcher := New()
	var inlineCalls int
	var chosenCalls int
	must(t, dispatcher.OnInlineQueryFunc(func(ctx context.Context, update telegram.Update) error {
		inlineCalls++
		if update.InlineQuery == nil || update.InlineQuery.From.ID != 777 {
			t.Fatalf("unexpected inline query update: %+v", update)
		}
		return nil
	}))
	must(t, dispatcher.OnChosenInlineResultFunc(func(ctx context.Context, update telegram.Update) error {
		chosenCalls++
		if update.ChosenInlineResult == nil || update.ChosenInlineResult.ResultID != "result-1" {
			t.Fatalf("unexpected chosen inline result update: %+v", update)
		}
		return nil
	}))

	if err := dispatcher.HandleUpdate(context.Background(), inlineQueryUpdate()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := dispatcher.HandleUpdate(context.Background(), chosenInlineResultUpdate()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := dispatcher.HandleUpdate(context.Background(), messageUpdate("hello")); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if inlineCalls != 1 || chosenCalls != 1 {
		t.Fatalf("unexpected calls: inline=%d chosen=%d", inlineCalls, chosenCalls)
	}
	if !InlineQuery()(inlineQueryUpdate()) {
		t.Fatal("inline query predicate should match")
	}
	if InlineQuery()(messageUpdate("hello")) {
		t.Fatal("inline query predicate should not match message updates")
	}
	if !ChosenInlineResult()(chosenInlineResultUpdate()) {
		t.Fatal("chosen inline result predicate should match")
	}
	if ChosenInlineResult()(messageUpdate("hello")) {
		t.Fatal("chosen inline result predicate should not match message updates")
	}
}

func TestChatJoinRequestRoute(t *testing.T) {
	dispatcher := New()
	var calls int
	must(t, dispatcher.OnChatJoinRequestFunc(func(ctx context.Context, update telegram.Update) error {
		calls++
		if update.ChatJoinRequest == nil || update.ChatJoinRequest.From.ID != 777 {
			t.Fatalf("unexpected join request update: %+v", update)
		}
		return nil
	}))

	if err := dispatcher.HandleUpdate(context.Background(), joinRequestUpdate()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := dispatcher.HandleUpdate(context.Background(), messageUpdate("hello")); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if calls != 1 {
		t.Fatalf("unexpected calls: %d", calls)
	}
	if !ChatJoinRequest()(joinRequestUpdate()) {
		t.Fatal("chat join request predicate should match")
	}
	if ChatJoinRequest()(messageUpdate("hello")) {
		t.Fatal("chat join request predicate should not match message updates")
	}
}

func TestMessageReactionRoutes(t *testing.T) {
	dispatcher := New()
	var reactionCalls int
	var countCalls int
	must(t, dispatcher.OnMessageReactionFunc(func(ctx context.Context, update telegram.Update) error {
		reactionCalls++
		if update.MessageReaction == nil || update.MessageReaction.User == nil || update.MessageReaction.User.ID != 777 {
			t.Fatalf("unexpected message reaction update: %+v", update)
		}
		return nil
	}))
	must(t, dispatcher.OnMessageReactionCountFunc(func(ctx context.Context, update telegram.Update) error {
		countCalls++
		if update.MessageReactionCount == nil || update.MessageReactionCount.MessageID != 456 {
			t.Fatalf("unexpected message reaction count update: %+v", update)
		}
		return nil
	}))

	if err := dispatcher.HandleUpdate(context.Background(), messageReactionUpdate()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := dispatcher.HandleUpdate(context.Background(), messageReactionCountUpdate()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := dispatcher.HandleUpdate(context.Background(), messageUpdate("hello")); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if reactionCalls != 1 || countCalls != 1 {
		t.Fatalf("unexpected calls: reaction=%d count=%d", reactionCalls, countCalls)
	}
	if !MessageReaction()(messageReactionUpdate()) {
		t.Fatal("message reaction predicate should match")
	}
	if MessageReaction()(messageUpdate("hello")) {
		t.Fatal("message reaction predicate should not match message updates")
	}
	if !MessageReactionCount()(messageReactionCountUpdate()) {
		t.Fatal("message reaction count predicate should match")
	}
	if MessageReactionCount()(messageUpdate("hello")) {
		t.Fatal("message reaction count predicate should not match message updates")
	}
}

func TestChatMemberRoutes(t *testing.T) {
	dispatcher := New()
	var myChatMemberCalls int
	var chatMemberCalls int
	must(t, dispatcher.OnMyChatMemberFunc(func(ctx context.Context, update telegram.Update) error {
		myChatMemberCalls++
		if update.MyChatMember == nil || update.MyChatMember.From.ID != 777 {
			t.Fatalf("unexpected my_chat_member update: %+v", update)
		}
		return nil
	}))
	must(t, dispatcher.OnChatMemberFunc(func(ctx context.Context, update telegram.Update) error {
		chatMemberCalls++
		if update.ChatMember == nil || update.ChatMember.From.ID != 778 {
			t.Fatalf("unexpected chat_member update: %+v", update)
		}
		return nil
	}))

	if err := dispatcher.HandleUpdate(context.Background(), myChatMemberUpdate()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := dispatcher.HandleUpdate(context.Background(), chatMemberUpdate()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := dispatcher.HandleUpdate(context.Background(), messageUpdate("hello")); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if myChatMemberCalls != 1 || chatMemberCalls != 1 {
		t.Fatalf("unexpected calls: my_chat_member=%d chat_member=%d", myChatMemberCalls, chatMemberCalls)
	}
	if !MyChatMember()(myChatMemberUpdate()) {
		t.Fatal("my_chat_member predicate should match")
	}
	if MyChatMember()(messageUpdate("hello")) {
		t.Fatal("my_chat_member predicate should not match message updates")
	}
	if !ChatMember()(chatMemberUpdate()) {
		t.Fatal("chat_member predicate should match")
	}
	if ChatMember()(messageUpdate("hello")) {
		t.Fatal("chat_member predicate should not match message updates")
	}
}

func TestChatBoostRoutes(t *testing.T) {
	dispatcher := New()
	var boostCalls int
	var removedCalls int
	must(t, dispatcher.OnChatBoostFunc(func(ctx context.Context, update telegram.Update) error {
		boostCalls++
		if update.ChatBoost == nil || update.ChatBoost.Boost.BoostID != "boost-1" {
			t.Fatalf("unexpected chat_boost update: %+v", update)
		}
		return nil
	}))
	must(t, dispatcher.OnRemovedChatBoostFunc(func(ctx context.Context, update telegram.Update) error {
		removedCalls++
		if update.RemovedChatBoost == nil || update.RemovedChatBoost.BoostID != "boost-2" {
			t.Fatalf("unexpected removed_chat_boost update: %+v", update)
		}
		return nil
	}))

	if err := dispatcher.HandleUpdate(context.Background(), chatBoostUpdate()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := dispatcher.HandleUpdate(context.Background(), removedChatBoostUpdate()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := dispatcher.HandleUpdate(context.Background(), messageUpdate("hello")); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if boostCalls != 1 || removedCalls != 1 {
		t.Fatalf("unexpected calls: boost=%d removed=%d", boostCalls, removedCalls)
	}
	if !ChatBoost()(chatBoostUpdate()) {
		t.Fatal("chat_boost predicate should match")
	}
	if ChatBoost()(messageUpdate("hello")) {
		t.Fatal("chat_boost predicate should not match message updates")
	}
	if !RemovedChatBoost()(removedChatBoostUpdate()) {
		t.Fatal("removed_chat_boost predicate should match")
	}
	if RemovedChatBoost()(messageUpdate("hello")) {
		t.Fatal("removed_chat_boost predicate should not match message updates")
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

func joinRequestUpdate() telegram.Update {
	return telegram.Update{
		UpdateID: 3,
		ChatJoinRequest: &telegram.ChatJoinRequest{
			Chat:       telegram.Chat{ID: -100123, Type: "supergroup"},
			From:       telegram.User{ID: 777, FirstName: "Joiner"},
			UserChatID: 888,
			Date:       1234567890,
		},
	}
}

func messageReactionUpdate() telegram.Update {
	return telegram.Update{
		UpdateID: 4,
		MessageReaction: &telegram.MessageReactionUpdated{
			Chat:      telegram.Chat{ID: -100123, Type: "supergroup"},
			MessageID: 123,
			User:      &telegram.User{ID: 777, FirstName: "Alice"},
			Date:      1234567890,
		},
	}
}

func messageReactionCountUpdate() telegram.Update {
	return telegram.Update{
		UpdateID: 5,
		MessageReactionCount: &telegram.MessageReactionCountUpdated{
			Chat:      telegram.Chat{ID: -100123, Type: "supergroup"},
			MessageID: 456,
			Date:      1234567891,
		},
	}
}

func myChatMemberUpdate() telegram.Update {
	return telegram.Update{
		UpdateID: 8,
		MyChatMember: &telegram.ChatMemberUpdated{
			Chat:          telegram.Chat{ID: -100123, Type: "supergroup"},
			From:          telegram.User{ID: 777, FirstName: "Admin"},
			Date:          1234567890,
			OldChatMember: telegram.ChatMember{Status: telegram.ChatMemberStatusLeft, User: telegram.User{ID: 1}},
			NewChatMember: telegram.ChatMember{Status: telegram.ChatMemberStatusMember, User: telegram.User{ID: 1}},
		},
	}
}

func chatMemberUpdate() telegram.Update {
	return telegram.Update{
		UpdateID: 9,
		ChatMember: &telegram.ChatMemberUpdated{
			Chat:          telegram.Chat{ID: -100123, Type: "supergroup"},
			From:          telegram.User{ID: 778, FirstName: "Admin"},
			Date:          1234567891,
			OldChatMember: telegram.ChatMember{Status: telegram.ChatMemberStatusMember, User: telegram.User{ID: 2}},
			NewChatMember: telegram.ChatMember{Status: telegram.ChatMemberStatusRestricted, User: telegram.User{ID: 2}},
		},
	}
}

func chatBoostUpdate() telegram.Update {
	return telegram.Update{
		UpdateID: 10,
		ChatBoost: &telegram.ChatBoostUpdated{
			Chat: telegram.Chat{ID: -100124, Type: "supergroup"},
			Boost: telegram.ChatBoost{
				BoostID: "boost-1",
				Source:  telegram.ChatBoostSourcePremium{Source: "premium", User: telegram.User{ID: 779, FirstName: "Booster"}},
			},
		},
	}
}

func removedChatBoostUpdate() telegram.Update {
	return telegram.Update{
		UpdateID: 11,
		RemovedChatBoost: &telegram.ChatBoostRemoved{
			Chat:    telegram.Chat{ID: -100124, Type: "supergroup"},
			BoostID: "boost-2",
			Source:  telegram.ChatBoostSourceGiftCode{Source: "gift_code", User: telegram.User{ID: 780, FirstName: "Booster"}},
		},
	}
}

func inlineQueryUpdate() telegram.Update {
	return telegram.Update{
		UpdateID: 6,
		InlineQuery: &telegram.InlineQuery{
			ID:     "inline-query",
			From:   telegram.User{ID: 777, FirstName: "Alice"},
			Query:  "hello",
			Offset: "",
		},
	}
}

func chosenInlineResultUpdate() telegram.Update {
	return telegram.Update{
		UpdateID: 7,
		ChosenInlineResult: &telegram.ChosenInlineResult{
			ResultID: "result-1",
			From:     telegram.User{ID: 778, FirstName: "Bob"},
			Query:    "hello",
		},
	}
}

func TestPaymentQueryRoutes(t *testing.T) {
	dispatcher := New()
	var shippingCalls int
	var preCheckoutCalls int
	must(t, dispatcher.OnShippingQueryFunc(func(ctx context.Context, update telegram.Update) error {
		shippingCalls++
		if update.ShippingQuery == nil || update.ShippingQuery.From.ID != 777 {
			t.Fatalf("unexpected shipping query update: %+v", update)
		}
		return nil
	}))
	must(t, dispatcher.OnPreCheckoutQueryFunc(func(ctx context.Context, update telegram.Update) error {
		preCheckoutCalls++
		if update.PreCheckoutQuery == nil || update.PreCheckoutQuery.From.ID != 778 {
			t.Fatalf("unexpected pre-checkout query update: %+v", update)
		}
		return nil
	}))

	if err := dispatcher.HandleUpdate(context.Background(), shippingQueryUpdate()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := dispatcher.HandleUpdate(context.Background(), preCheckoutQueryUpdate()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := dispatcher.HandleUpdate(context.Background(), messageUpdate("hello")); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if shippingCalls != 1 || preCheckoutCalls != 1 {
		t.Fatalf("unexpected calls: shipping=%d pre_checkout=%d", shippingCalls, preCheckoutCalls)
	}
	if !ShippingQuery()(shippingQueryUpdate()) {
		t.Fatal("shipping query predicate should match")
	}
	if ShippingQuery()(messageUpdate("hello")) {
		t.Fatal("shipping query predicate should not match message updates")
	}
	if !PreCheckoutQuery()(preCheckoutQueryUpdate()) {
		t.Fatal("pre-checkout query predicate should match")
	}
	if PreCheckoutQuery()(messageUpdate("hello")) {
		t.Fatal("pre-checkout query predicate should not match message updates")
	}
}

func shippingQueryUpdate() telegram.Update {
	return telegram.Update{ShippingQuery: &telegram.ShippingQuery{ID: "ship", From: telegram.User{ID: 777, FirstName: "Alice"}}}
}

func preCheckoutQueryUpdate() telegram.Update {
	return telegram.Update{PreCheckoutQuery: &telegram.PreCheckoutQuery{ID: "pre", From: telegram.User{ID: 778, FirstName: "Bob"}}}
}

func TestPaidMediaPurchasedPredicateAndHandler(t *testing.T) {
	dispatcher := New()
	called := false
	if err := dispatcher.OnPaidMediaPurchasedFunc(func(ctx context.Context, update telegram.Update) error {
		called = true
		if update.PurchasedPaidMedia == nil || update.PurchasedPaidMedia.From.ID != 7 {
			t.Fatalf("unexpected update: %+v", update)
		}
		return nil
	}); err != nil {
		t.Fatalf("register handler: %v", err)
	}
	update := telegram.Update{UpdateID: 1, PurchasedPaidMedia: &telegram.PaidMediaPurchased{From: telegram.User{ID: 7, FirstName: "Alice"}, PaidMediaPayload: "payload"}}
	if !PaidMediaPurchased()(update) {
		t.Fatal("predicate should match paid media purchase")
	}
	if err := dispatcher.HandleUpdate(context.Background(), update); err != nil {
		t.Fatalf("handle update: %v", err)
	}
	if !called {
		t.Fatal("expected handler call")
	}
	if PaidMediaPurchased()(telegram.Update{UpdateID: 2}) {
		t.Fatal("predicate should not match unrelated update")
	}
}

func TestManagedBotRoute(t *testing.T) {
	dispatcher := New()
	var calls int
	must(t, dispatcher.OnManagedBotFunc(func(ctx context.Context, update telegram.Update) error {
		calls++
		if update.ManagedBot == nil || update.ManagedBot.User.ID != 7 || update.ManagedBot.Bot.ID != 77 {
			t.Fatalf("unexpected managed bot update: %+v", update)
		}
		return nil
	}))

	managed := telegram.Update{UpdateID: 500, ManagedBot: &telegram.ManagedBotUpdated{User: telegram.User{ID: 7, FirstName: "Owner"}, Bot: telegram.User{ID: 77, IsBot: true, FirstName: "Child"}}}
	if err := dispatcher.HandleUpdate(context.Background(), managed); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := dispatcher.HandleUpdate(context.Background(), messageUpdate("hello")); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if calls != 1 {
		t.Fatalf("unexpected calls: %d", calls)
	}
	if !ManagedBot()(managed) {
		t.Fatal("managed bot predicate should match")
	}
	if ManagedBot()(messageUpdate("hello")) {
		t.Fatal("managed bot predicate should not match message updates")
	}
}

func TestBusinessRoutes(t *testing.T) {
	dispatcher := New()
	var connectionCalls int
	var messageCalls int
	var editedCalls int
	var deletedCalls int

	must(t, dispatcher.OnBusinessConnectionFunc(func(ctx context.Context, update telegram.Update) error {
		connectionCalls++
		if update.BusinessConnection == nil || update.BusinessConnection.ID != "bc-1" {
			t.Fatalf("unexpected business connection update: %+v", update)
		}
		return nil
	}))
	must(t, dispatcher.OnBusinessMessageFunc(func(ctx context.Context, update telegram.Update) error {
		messageCalls++
		if update.BusinessMessage == nil || update.BusinessMessage.MessageID != 11 {
			t.Fatalf("unexpected business message update: %+v", update)
		}
		return nil
	}))
	must(t, dispatcher.OnEditedBusinessMessageFunc(func(ctx context.Context, update telegram.Update) error {
		editedCalls++
		if update.EditedBusinessMessage == nil || update.EditedBusinessMessage.MessageID != 12 {
			t.Fatalf("unexpected edited business message update: %+v", update)
		}
		return nil
	}))
	must(t, dispatcher.OnDeletedBusinessMessagesFunc(func(ctx context.Context, update telegram.Update) error {
		deletedCalls++
		if update.DeletedBusinessMessages == nil || update.DeletedBusinessMessages.BusinessConnectionID != "bc-1" {
			t.Fatalf("unexpected deleted business messages update: %+v", update)
		}
		return nil
	}))

	if err := dispatcher.HandleUpdate(context.Background(), businessConnectionUpdate()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := dispatcher.HandleUpdate(context.Background(), businessMessageUpdate()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := dispatcher.HandleUpdate(context.Background(), editedBusinessMessageUpdate()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := dispatcher.HandleUpdate(context.Background(), deletedBusinessMessagesUpdate()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := dispatcher.HandleUpdate(context.Background(), messageUpdate("hello")); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if connectionCalls != 1 || messageCalls != 1 || editedCalls != 1 || deletedCalls != 1 {
		t.Fatalf("unexpected calls: connection=%d message=%d edited=%d deleted=%d", connectionCalls, messageCalls, editedCalls, deletedCalls)
	}
	if !BusinessConnection()(businessConnectionUpdate()) || BusinessConnection()(messageUpdate("hello")) {
		t.Fatal("business connection predicate mismatch")
	}
	if !BusinessMessage()(businessMessageUpdate()) || BusinessMessage()(messageUpdate("hello")) {
		t.Fatal("business message predicate mismatch")
	}
	if !EditedBusinessMessage()(editedBusinessMessageUpdate()) || EditedBusinessMessage()(messageUpdate("hello")) {
		t.Fatal("edited business message predicate mismatch")
	}
	if !DeletedBusinessMessages()(deletedBusinessMessagesUpdate()) || DeletedBusinessMessages()(messageUpdate("hello")) {
		t.Fatal("deleted business messages predicate mismatch")
	}
}

func TestPollAnswerRoute(t *testing.T) {
	dispatcher := New()
	var calls int
	must(t, dispatcher.OnPollAnswerFunc(func(ctx context.Context, update telegram.Update) error {
		calls++
		if update.PollAnswer == nil || update.PollAnswer.PollID != "poll-id" {
			t.Fatalf("unexpected poll answer update: %+v", update)
		}
		return nil
	}))

	pollAnswer := telegram.Update{UpdateID: 600, PollAnswer: &telegram.PollAnswer{PollID: "poll-id", User: &telegram.User{ID: 7, FirstName: "Alice"}, OptionIDs: []int{0}}}
	if err := dispatcher.HandleUpdate(context.Background(), pollAnswer); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := dispatcher.HandleUpdate(context.Background(), messageUpdate("hello")); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if calls != 1 {
		t.Fatalf("unexpected calls: %d", calls)
	}
	if !PollAnswer()(pollAnswer) {
		t.Fatal("poll answer predicate should match")
	}
	if PollAnswer()(messageUpdate("hello")) {
		t.Fatal("poll answer predicate should not match message updates")
	}
}

func businessConnectionUpdate() telegram.Update {
	return telegram.Update{UpdateID: 700, BusinessConnection: &telegram.BusinessConnection{ID: "bc-1", User: telegram.User{ID: 7, FirstName: "Business"}, UserChatID: 7000, Date: 123, IsEnabled: true}}
}

func businessMessageUpdate() telegram.Update {
	return telegram.Update{UpdateID: 701, BusinessMessage: &telegram.Message{MessageID: 11, Chat: telegram.Chat{ID: 1001, Type: "private"}, BusinessConnectionID: "bc-1"}}
}

func editedBusinessMessageUpdate() telegram.Update {
	return telegram.Update{UpdateID: 702, EditedBusinessMessage: &telegram.Message{MessageID: 12, Chat: telegram.Chat{ID: 1002, Type: "private"}, BusinessConnectionID: "bc-1"}}
}

func deletedBusinessMessagesUpdate() telegram.Update {
	return telegram.Update{UpdateID: 703, DeletedBusinessMessages: &telegram.BusinessMessagesDeleted{BusinessConnectionID: "bc-1", Chat: telegram.Chat{ID: 1003, Type: "private"}, MessageIDs: []int64{11, 12}}}
}
