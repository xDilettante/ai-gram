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

func TestGetUserChatBoostsSendsPayloadAndDecodesResult(t *testing.T) {
	const token = "123:secret"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/bot"+token+"/getUserChatBoosts" {
			t.Fatalf("unexpected path: %q", r.URL.Path)
		}
		var payload map[string]any
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode payload: %v", err)
		}
		if payload["chat_id"] != float64(-100123) || payload["user_id"] != float64(777) {
			t.Fatalf("unexpected payload: %#v", payload)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true,"result":{"boosts":[{"boost_id":"boost-1","add_date":1,"expiration_date":2,"source":{"source":"premium","user":{"id":777,"is_bot":false,"first_name":"Booster"}}}]}}`))
	}))
	defer server.Close()

	bot := newTestBot(t, token, server.URL, server.Client())
	boosts, err := bot.GetUserChatBoosts(context.Background(), GetUserChatBoostsParams{ChatID: ChatIDInt(-100123), UserID: 777})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if boosts == nil || len(boosts.Boosts) != 1 || boosts.Boosts[0].BoostID != "boost-1" {
		t.Fatalf("unexpected boosts: %+v", boosts)
	}
	if source, ok := boosts.Boosts[0].Source.(telegram.ChatBoostSourcePremium); !ok || source.User.ID != 777 {
		t.Fatalf("unexpected boost source: %#v", boosts.Boosts[0].Source)
	}
}

func TestSetChatMemberTagSendsPayloadAndAllowsEmptyTag(t *testing.T) {
	const token = "123:secret"
	tests := []struct {
		name string
		tag  string
	}{
		{name: "set tag", tag: "vip"},
		{name: "clear tag", tag: ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path != "/bot"+token+"/setChatMemberTag" {
					t.Fatalf("unexpected path: %q", r.URL.Path)
				}
				var payload map[string]any
				if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
					t.Fatalf("decode payload: %v", err)
				}
				if payload["chat_id"] != float64(-100123) || payload["user_id"] != float64(777) || payload["tag"] != tt.tag {
					t.Fatalf("unexpected payload: %#v", payload)
				}
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(`{"ok":true,"result":true}`))
			}))
			defer server.Close()

			bot := newTestBot(t, token, server.URL, server.Client())
			ok, err := bot.SetChatMemberTag(context.Background(), SetChatMemberTagParams{ChatID: ChatIDInt(-100123), UserID: 777, Tag: tt.tag})
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !ok {
				t.Fatal("expected true result")
			}
		})
	}
}

func TestBanChatSenderChatSendsPayloadAndDecodesResult(t *testing.T) {
	testSenderChatMethod(t, "banChatSenderChat", func(ctx context.Context, bot *Bot) (bool, error) {
		return bot.BanChatSenderChat(ctx, BanChatSenderChatParams{ChatID: ChatIDInt(-100123), SenderChatID: -100777})
	})
}

func TestUnbanChatSenderChatSendsPayloadAndDecodesResult(t *testing.T) {
	testSenderChatMethod(t, "unbanChatSenderChat", func(ctx context.Context, bot *Bot) (bool, error) {
		return bot.UnbanChatSenderChat(ctx, UnbanChatSenderChatParams{ChatID: ChatIDInt(-100123), SenderChatID: -100777})
	})
}

func TestChatBoostMemberMethodValidation(t *testing.T) {
	const token = "123:secret"
	bot := newTestBot(t, token, "https://example.test", nil)
	tests := []struct {
		name string
		call func() error
	}{
		{name: "boosts empty chat", call: func() error {
			_, err := bot.GetUserChatBoosts(context.Background(), GetUserChatBoostsParams{UserID: 1})
			return err
		}},
		{name: "boosts invalid user", call: func() error {
			_, err := bot.GetUserChatBoosts(context.Background(), GetUserChatBoostsParams{ChatID: ChatIDInt(123)})
			return err
		}},
		{name: "tag empty chat", call: func() error {
			_, err := bot.SetChatMemberTag(context.Background(), SetChatMemberTagParams{UserID: 1, Tag: "tag"})
			return err
		}},
		{name: "tag invalid user", call: func() error {
			_, err := bot.SetChatMemberTag(context.Background(), SetChatMemberTagParams{ChatID: ChatIDInt(123), UserID: -1})
			return err
		}},
		{name: "ban sender empty chat", call: func() error {
			_, err := bot.BanChatSenderChat(context.Background(), BanChatSenderChatParams{SenderChatID: -1001})
			return err
		}},
		{name: "ban sender zero sender", call: func() error {
			_, err := bot.BanChatSenderChat(context.Background(), BanChatSenderChatParams{ChatID: ChatIDInt(123)})
			return err
		}},
		{name: "unban sender empty chat", call: func() error {
			_, err := bot.UnbanChatSenderChat(context.Background(), UnbanChatSenderChatParams{SenderChatID: -1001})
			return err
		}},
		{name: "unban sender zero sender", call: func() error {
			_, err := bot.UnbanChatSenderChat(context.Background(), UnbanChatSenderChatParams{ChatID: ChatIDInt(123)})
			return err
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.call()
			if err == nil {
				t.Fatal("expected error")
			}
			assertNoToken(t, err, token)
		})
	}
}

func TestChatBoostMemberMethodsReturnAPIError(t *testing.T) {
	const token = "123:secret"
	tests := chatBoostMemberCalls()
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
			_, err := tt.call(context.Background(), bot)
			if err == nil {
				t.Fatal("expected error")
			}
			var apiErr *apierrors.APIError
			if !stderrors.As(err, &apiErr) {
				t.Fatalf("expected APIError, got %T", err)
			}
			assertNoToken(t, err, token)
		})
	}
}

func TestChatBoostMemberMethodsResponseAndContextErrors(t *testing.T) {
	const token = "123:secret"
	for _, tt := range chatBoostMemberCalls() {
		t.Run(tt.name+" cancelled context", func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				t.Fatal("request should not reach server")
			}))
			defer server.Close()
			bot := newTestBot(t, token, server.URL, server.Client())
			ctx, cancel := context.WithCancel(context.Background())
			cancel()
			_, err := tt.call(ctx, bot)
			if err == nil {
				t.Fatal("expected error")
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
			_, err := tt.call(context.Background(), bot)
			if err == nil {
				t.Fatal("expected error")
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
			_, err := tt.call(context.Background(), bot)
			if err == nil {
				t.Fatal("expected error")
			}
			assertNoToken(t, err, token)
		})
	}
}

func testSenderChatMethod(t *testing.T, method string, call func(context.Context, *Bot) (bool, error)) {
	t.Helper()
	const token = "123:secret"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/bot"+token+"/"+method {
			t.Fatalf("unexpected path: %q", r.URL.Path)
		}
		var payload map[string]any
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode payload: %v", err)
		}
		if payload["chat_id"] != float64(-100123) || payload["sender_chat_id"] != float64(-100777) {
			t.Fatalf("unexpected payload: %#v", payload)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true,"result":true}`))
	}))
	defer server.Close()

	bot := newTestBot(t, token, server.URL, server.Client())
	ok, err := call(context.Background(), bot)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Fatal("expected true result")
	}
}

type chatBoostMemberCall struct {
	name   string
	method string
	call   func(context.Context, *Bot) (any, error)
}

func chatBoostMemberCalls() []chatBoostMemberCall {
	return []chatBoostMemberCall{
		{name: "get boosts", method: "getUserChatBoosts", call: func(ctx context.Context, bot *Bot) (any, error) {
			return bot.GetUserChatBoosts(ctx, GetUserChatBoostsParams{ChatID: ChatIDInt(123), UserID: 1})
		}},
		{name: "set tag", method: "setChatMemberTag", call: func(ctx context.Context, bot *Bot) (any, error) {
			return bot.SetChatMemberTag(ctx, SetChatMemberTagParams{ChatID: ChatIDInt(123), UserID: 1, Tag: "tag"})
		}},
		{name: "ban sender chat", method: "banChatSenderChat", call: func(ctx context.Context, bot *Bot) (any, error) {
			return bot.BanChatSenderChat(ctx, BanChatSenderChatParams{ChatID: ChatIDInt(123), SenderChatID: -1001})
		}},
		{name: "unban sender chat", method: "unbanChatSenderChat", call: func(ctx context.Context, bot *Bot) (any, error) {
			return bot.UnbanChatSenderChat(ctx, UnbanChatSenderChatParams{ChatID: ChatIDInt(123), SenderChatID: -1001})
		}},
	}
}
