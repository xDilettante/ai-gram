package bot

import (
	"context"
	"encoding/json"
	stderrors "errors"
	"net/http"
	"net/http/httptest"
	"testing"

	apierrors "github.com/xDilettante/ai-gram/errors"
	"github.com/xDilettante/ai-gram/telegram"
)

func TestBanChatMemberSendsPayloadAndDecodesResult(t *testing.T) {
	const token = "123:secret"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/bot"+token+"/banChatMember" {
			t.Fatalf("unexpected path: %q", r.URL.Path)
		}
		var payload map[string]any
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode payload: %v", err)
		}
		if payload["chat_id"] != float64(12345) || payload["user_id"] != float64(777) || payload["until_date"] != float64(1700000000) || payload["revoke_messages"] != true {
			t.Fatalf("unexpected payload: %#v", payload)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true,"result":true}`))
	}))
	defer server.Close()

	bot := newTestBot(t, token, server.URL, server.Client())
	ok, err := bot.BanChatMember(context.Background(), BanChatMemberParams{ChatID: ChatIDInt(12345), UserID: 777, UntilDate: 1700000000, RevokeMessages: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Fatal("expected true result")
	}
}

func TestBanChatMemberValidation(t *testing.T) {
	const token = "123:secret"
	bot := newTestBot(t, token, "https://example.test", nil)
	tests := []struct {
		name   string
		params BanChatMemberParams
	}{
		{name: "empty chat", params: BanChatMemberParams{UserID: 1}},
		{name: "zero user", params: BanChatMemberParams{ChatID: ChatIDInt(123)}},
		{name: "negative user", params: BanChatMemberParams{ChatID: ChatIDInt(123), UserID: -1}},
		{name: "negative until", params: BanChatMemberParams{ChatID: ChatIDInt(123), UserID: 1, UntilDate: -1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ok, err := bot.BanChatMember(context.Background(), tt.params)
			if err == nil {
				t.Fatal("expected error")
			}
			if ok {
				t.Fatal("expected false result")
			}
			assertNoToken(t, err, token)
		})
	}
}

func TestUnbanChatMemberSendsPayloadAndDecodesResult(t *testing.T) {
	const token = "123:secret"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/bot"+token+"/unbanChatMember" {
			t.Fatalf("unexpected path: %q", r.URL.Path)
		}
		var payload map[string]any
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode payload: %v", err)
		}
		if payload["chat_id"] != float64(12345) || payload["user_id"] != float64(777) || payload["only_if_banned"] != true {
			t.Fatalf("unexpected payload: %#v", payload)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true,"result":true}`))
	}))
	defer server.Close()

	bot := newTestBot(t, token, server.URL, server.Client())
	ok, err := bot.UnbanChatMember(context.Background(), UnbanChatMemberParams{ChatID: ChatIDInt(12345), UserID: 777, OnlyIfBanned: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Fatal("expected true result")
	}
}

func TestUnbanChatMemberValidation(t *testing.T) {
	const token = "123:secret"
	bot := newTestBot(t, token, "https://example.test", nil)
	tests := []struct {
		name   string
		params UnbanChatMemberParams
	}{
		{name: "empty chat", params: UnbanChatMemberParams{UserID: 1}},
		{name: "zero user", params: UnbanChatMemberParams{ChatID: ChatIDInt(123)}},
		{name: "negative user", params: UnbanChatMemberParams{ChatID: ChatIDInt(123), UserID: -1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ok, err := bot.UnbanChatMember(context.Background(), tt.params)
			if err == nil {
				t.Fatal("expected error")
			}
			if ok {
				t.Fatal("expected false result")
			}
			assertNoToken(t, err, token)
		})
	}
}

func TestRestrictChatMemberSendsPayloadAndDecodesResult(t *testing.T) {
	const token = "123:secret"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/bot"+token+"/restrictChatMember" {
			t.Fatalf("unexpected path: %q", r.URL.Path)
		}
		var payload map[string]any
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode payload: %v", err)
		}
		permissions, ok := payload["permissions"].(map[string]any)
		if !ok {
			t.Fatalf("expected permissions object: %#v", payload)
		}
		if payload["chat_id"] != float64(12345) || payload["user_id"] != float64(777) || payload["use_independent_chat_permissions"] != true || payload["until_date"] != float64(1700000000) {
			t.Fatalf("unexpected payload: %#v", payload)
		}
		if permissions["can_send_messages"] != true || permissions["can_send_photos"] != true || permissions["can_manage_topics"] != true {
			t.Fatalf("unexpected permissions: %#v", permissions)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true,"result":true}`))
	}))
	defer server.Close()

	bot := newTestBot(t, token, server.URL, server.Client())
	ok, err := bot.RestrictChatMember(context.Background(), RestrictChatMemberParams{
		ChatID: ChatIDInt(12345),
		UserID: 777,
		Permissions: telegram.ChatPermissions{
			CanSendMessages: true,
			CanSendPhotos:   true,
			CanManageTopics: true,
		},
		UseIndependentChatPermissions: true,
		UntilDate:                     1700000000,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Fatal("expected true result")
	}
}

func TestRestrictChatMemberAllowsZeroPermissions(t *testing.T) {
	const token = "123:secret"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/bot"+token+"/restrictChatMember" {
			t.Fatalf("unexpected path: %q", r.URL.Path)
		}
		var payload map[string]any
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode payload: %v", err)
		}
		permissions, ok := payload["permissions"].(map[string]any)
		if !ok {
			t.Fatalf("expected permissions object: %#v", payload)
		}
		if len(permissions) != 0 {
			t.Fatalf("expected zero permissions object, got %#v", permissions)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true,"result":true}`))
	}))
	defer server.Close()

	bot := newTestBot(t, token, server.URL, server.Client())
	ok, err := bot.RestrictChatMember(context.Background(), RestrictChatMemberParams{ChatID: ChatIDInt(12345), UserID: 777})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Fatal("expected true result")
	}
}

func TestRestrictChatMemberValidation(t *testing.T) {
	const token = "123:secret"
	bot := newTestBot(t, token, "https://example.test", nil)
	tests := []struct {
		name   string
		params RestrictChatMemberParams
	}{
		{name: "empty chat", params: RestrictChatMemberParams{UserID: 1}},
		{name: "zero user", params: RestrictChatMemberParams{ChatID: ChatIDInt(123)}},
		{name: "negative user", params: RestrictChatMemberParams{ChatID: ChatIDInt(123), UserID: -1}},
		{name: "negative until", params: RestrictChatMemberParams{ChatID: ChatIDInt(123), UserID: 1, UntilDate: -1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ok, err := bot.RestrictChatMember(context.Background(), tt.params)
			if err == nil {
				t.Fatal("expected error")
			}
			if ok {
				t.Fatal("expected false result")
			}
			assertNoToken(t, err, token)
		})
	}
}

func TestModerationMethodsReturnAPIError(t *testing.T) {
	const token = "123:secret"
	tests := []struct {
		name   string
		method string
		call   func(*Bot) (bool, error)
	}{
		{name: "ban", method: "banChatMember", call: func(bot *Bot) (bool, error) {
			return bot.BanChatMember(context.Background(), BanChatMemberParams{ChatID: ChatIDInt(123), UserID: 1})
		}},
		{name: "unban", method: "unbanChatMember", call: func(bot *Bot) (bool, error) {
			return bot.UnbanChatMember(context.Background(), UnbanChatMemberParams{ChatID: ChatIDInt(123), UserID: 1})
		}},
		{name: "restrict", method: "restrictChatMember", call: func(bot *Bot) (bool, error) {
			return bot.RestrictChatMember(context.Background(), RestrictChatMemberParams{ChatID: ChatIDInt(123), UserID: 1})
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path != "/bot"+token+"/"+tt.method {
					t.Fatalf("unexpected path: %q", r.URL.Path)
				}
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(`{"ok":false,"error_code":400,"description":"Bad Request"}`))
			}))
			defer server.Close()

			bot := newTestBot(t, token, server.URL, server.Client())
			ok, err := tt.call(bot)
			if err == nil {
				t.Fatal("expected error")
			}
			if ok {
				t.Fatal("expected false result")
			}
			var apiErr *apierrors.APIError
			if !stderrors.As(err, &apiErr) {
				t.Fatalf("expected APIError, got %T", err)
			}
			assertNoToken(t, err, token)
		})
	}
}

func TestModerationMethodsResponseAndContextErrors(t *testing.T) {
	const token = "123:secret"
	tests := []struct {
		name   string
		method string
		call   func(context.Context, *Bot) (bool, error)
	}{
		{name: "ban", method: "banChatMember", call: func(ctx context.Context, bot *Bot) (bool, error) {
			return bot.BanChatMember(ctx, BanChatMemberParams{ChatID: ChatIDInt(123), UserID: 1})
		}},
		{name: "unban", method: "unbanChatMember", call: func(ctx context.Context, bot *Bot) (bool, error) {
			return bot.UnbanChatMember(ctx, UnbanChatMemberParams{ChatID: ChatIDInt(123), UserID: 1})
		}},
		{name: "restrict", method: "restrictChatMember", call: func(ctx context.Context, bot *Bot) (bool, error) {
			return bot.RestrictChatMember(ctx, RestrictChatMemberParams{ChatID: ChatIDInt(123), UserID: 1})
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name+" cancelled context", func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				t.Fatal("request should not reach server")
			}))
			defer server.Close()
			bot := newTestBot(t, token, server.URL, server.Client())
			ctx, cancel := context.WithCancel(context.Background())
			cancel()
			ok, err := tt.call(ctx, bot)
			if err == nil {
				t.Fatal("expected error")
			}
			if ok {
				t.Fatal("expected false result")
			}
			assertNoToken(t, err, token)
		})

		t.Run(tt.name+" invalid json", func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path != "/bot"+token+"/"+tt.method {
					t.Fatalf("unexpected path: %q", r.URL.Path)
				}
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(`not-json`))
			}))
			defer server.Close()
			bot := newTestBot(t, token, server.URL, server.Client())
			ok, err := tt.call(context.Background(), bot)
			if err == nil {
				t.Fatal("expected error")
			}
			if ok {
				t.Fatal("expected false result")
			}
			assertNoToken(t, err, token)
		})

		t.Run(tt.name+" http status", func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path != "/bot"+token+"/"+tt.method {
					t.Fatalf("unexpected path: %q", r.URL.Path)
				}
				http.Error(w, "server error", http.StatusInternalServerError)
			}))
			defer server.Close()
			bot := newTestBot(t, token, server.URL, server.Client())
			ok, err := tt.call(context.Background(), bot)
			if err == nil {
				t.Fatal("expected error")
			}
			if ok {
				t.Fatal("expected false result")
			}
			assertNoToken(t, err, token)
		})
	}
}
