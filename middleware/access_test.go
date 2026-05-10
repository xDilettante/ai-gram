package middleware

import (
	"context"
	stderrors "errors"
	"strings"
	"testing"

	"github.com/xDilettante/ai-gram/dispatch"
	"github.com/xDilettante/ai-gram/telegram"
)

func TestAccessModeOffAllowsAnyUpdate(t *testing.T) {
	called := false
	handler := Access(AccessConfig{Mode: AccessModeOff})(dispatch.HandlerFunc(func(context.Context, telegram.Update) error {
		called = true
		return nil
	}))

	if err := handler.HandleUpdate(context.Background(), telegram.Update{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !called {
		t.Fatal("expected handler to be called")
	}
}

func TestAccessModePublicAllowsAnyUpdate(t *testing.T) {
	called := false
	handler := Access(AccessConfig{Mode: AccessModePublic})(dispatch.HandlerFunc(func(context.Context, telegram.Update) error {
		called = true
		return nil
	}))

	if err := handler.HandleUpdate(context.Background(), telegram.Update{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !called {
		t.Fatal("expected handler to be called")
	}
}

func TestAccessModeAdminAllowsConfiguredPrincipals(t *testing.T) {
	tests := []struct {
		name   string
		config AccessConfig
		update telegram.Update
	}{
		{
			name:   "admin user",
			config: AccessConfig{Mode: AccessModeAdmin, AdminUserIDs: []int64{11}},
			update: accessMessageUpdate(11, 101),
		},
		{
			name:   "allowed user",
			config: AccessConfig{Mode: AccessModeAdmin, AllowedUserIDs: []int64{12}},
			update: accessMessageUpdate(12, 101),
		},
		{
			name:   "allowed chat",
			config: AccessConfig{Mode: AccessModeAdmin, AllowedChatIDs: []int64{202}},
			update: accessMessageUpdate(99, 202),
		},
		{
			name:   "empty mode is admin",
			config: AccessConfig{AdminUserIDs: []int64{11}},
			update: accessCallbackUpdate(11, 101),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			called := false
			handler := Access(tt.config)(dispatch.HandlerFunc(func(context.Context, telegram.Update) error {
				called = true
				return nil
			}))
			if err := handler.HandleUpdate(context.Background(), tt.update); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !called {
				t.Fatal("expected handler to be called")
			}
		})
	}
}

func TestAccessModeAdminDeniesUnknownOrMissingPrincipals(t *testing.T) {
	tests := []struct {
		name   string
		update telegram.Update
	}{
		{name: "unknown user", update: accessMessageUpdate(99, 101)},
		{name: "no user or chat", update: telegram.Update{UpdateID: 3}},
		{name: "nil message safe", update: telegram.Update{UpdateID: 4, CallbackQuery: &telegram.CallbackQuery{ID: "callback", From: telegram.User{ID: 99}}}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			denied := false
			called := false
			handler := Access(AccessConfig{
				Mode:         AccessModeAdmin,
				AdminUserIDs: []int64{11},
				OnDeny: func(context.Context, telegram.Update) error {
					denied = true
					return nil
				},
			})(dispatch.HandlerFunc(func(context.Context, telegram.Update) error {
				called = true
				return nil
			}))

			if err := handler.HandleUpdate(context.Background(), tt.update); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if called {
				t.Fatal("handler must not be called")
			}
			if !denied {
				t.Fatal("expected OnDeny to be called")
			}
		})
	}
}

func TestAccessOnDenyErrorIsReturned(t *testing.T) {
	want := stderrors.New("denied")
	handler := Access(AccessConfig{
		Mode: AccessModeAdmin,
		OnDeny: func(context.Context, telegram.Update) error {
			return want
		},
	})(dispatch.HandlerFunc(func(context.Context, telegram.Update) error {
		t.Fatal("handler must not be called")
		return nil
	}))

	if err := handler.HandleUpdate(context.Background(), accessMessageUpdate(99, 101)); !stderrors.Is(err, want) {
		t.Fatalf("got %v, want %v", err, want)
	}
}

func TestAccessSilentlyIgnoresDeniedUpdateWithoutOnDeny(t *testing.T) {
	called := false
	handler := Access(AccessConfig{Mode: AccessModeAdmin})(dispatch.HandlerFunc(func(context.Context, telegram.Update) error {
		called = true
		return nil
	}))

	if err := handler.HandleUpdate(context.Background(), accessMessageUpdate(99, 101)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if called {
		t.Fatal("handler must not be called")
	}
}

func TestAccessWithPolicy(t *testing.T) {
	allow := true
	denied := false
	called := 0
	handler := AccessWithPolicy(AccessFunc(func(context.Context, telegram.Update) bool {
		return allow
	}), func(context.Context, telegram.Update) error {
		denied = true
		return nil
	})(dispatch.HandlerFunc(func(context.Context, telegram.Update) error {
		called++
		return nil
	}))

	if err := handler.HandleUpdate(context.Background(), accessMessageUpdate(1, 1)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	allow = false
	if err := handler.HandleUpdate(context.Background(), accessMessageUpdate(1, 1)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if called != 1 {
		t.Fatalf("handler call count = %d, want 1", called)
	}
	if !denied {
		t.Fatal("expected onDeny call")
	}
}

func TestAccessWithNilPolicyDeniesSafely(t *testing.T) {
	denied := false
	handler := AccessWithPolicy(nil, func(context.Context, telegram.Update) error {
		denied = true
		return nil
	})(dispatch.HandlerFunc(func(context.Context, telegram.Update) error {
		t.Fatal("handler must not be called")
		return nil
	}))

	if err := handler.HandleUpdate(context.Background(), accessMessageUpdate(1, 1)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !denied {
		t.Fatal("expected deny callback")
	}
}

func TestAccessErrorsDoNotLeakTokenLikeValues(t *testing.T) {
	handler := AccessWithPolicy(nil, func(context.Context, telegram.Update) error {
		return stderrors.New("access denied")
	})(dispatch.HandlerFunc(func(context.Context, telegram.Update) error { return nil }))

	err := handler.HandleUpdate(context.Background(), accessMessageUpdate(123, 456))
	if err == nil {
		t.Fatal("expected error")
	}
	if strings.Contains(err.Error(), "123:secret") {
		t.Fatalf("error leaked token-like value: %q", err.Error())
	}
}

func accessMessageUpdate(userID int64, chatID int64) telegram.Update {
	return telegram.Update{
		UpdateID: 1,
		Message: &telegram.Message{
			MessageID: 10,
			From:      &telegram.User{ID: userID, FirstName: "User"},
			Chat:      telegram.Chat{ID: chatID, Type: "private"},
			Date:      1,
			Text:      "hello",
		},
	}
}

func accessCallbackUpdate(userID int64, chatID int64) telegram.Update {
	return telegram.Update{
		UpdateID: 2,
		CallbackQuery: &telegram.CallbackQuery{
			ID:   "callback",
			From: telegram.User{ID: userID, FirstName: "User"},
			Data: "demo:edit",
			Message: &telegram.MaybeInaccessibleMessage{Message: &telegram.Message{
				MessageID: 11,
				Chat:      telegram.Chat{ID: chatID, Type: "private"},
				Date:      1,
			}, MessageID: 11, Chat: telegram.Chat{ID: chatID, Type: "private"}, Date: 1},
		},
	}
}
