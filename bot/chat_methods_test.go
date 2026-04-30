package bot

import (
	"context"
	"encoding/json"
	stderrors "errors"
	"net/http"
	"net/http/httptest"
	"testing"

	apierrors "github.com/xDilettante/ai-gram/errors"
)

func TestSendChatActionSendsPayloadAndDecodesResult(t *testing.T) {
	const token = "123:secret"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/bot"+token+"/sendChatAction" {
			t.Fatalf("unexpected path: %q", r.URL.Path)
		}
		var payload map[string]any
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode payload: %v", err)
		}
		if payload["chat_id"] != float64(12345) || payload["message_thread_id"] != float64(77) || payload["action"] != ChatActionTyping {
			t.Fatalf("unexpected payload: %#v", payload)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true,"result":true}`))
	}))
	defer server.Close()

	bot := newTestBot(t, token, server.URL, server.Client())
	ok, err := bot.SendChatAction(context.Background(), SendChatActionParams{ChatID: ChatIDInt(12345), MessageThreadID: 77, Action: ChatActionTyping})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Fatal("expected true result")
	}
}

func TestSendChatActionValidation(t *testing.T) {
	const token = "123:secret"
	bot := newTestBot(t, token, "https://example.test", nil)
	tests := []struct {
		name   string
		params SendChatActionParams
	}{
		{name: "empty chat", params: SendChatActionParams{Action: ChatActionTyping}},
		{name: "empty action", params: SendChatActionParams{ChatID: ChatIDInt(123)}},
		{name: "unknown action", params: SendChatActionParams{ChatID: ChatIDInt(123), Action: "unknown"}},
		{name: "negative thread", params: SendChatActionParams{ChatID: ChatIDInt(123), Action: ChatActionTyping, MessageThreadID: -1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ok, err := bot.SendChatAction(context.Background(), tt.params)
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

func TestPinChatMessageSendsPayloadAndDecodesResult(t *testing.T) {
	const token = "123:secret"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/bot"+token+"/pinChatMessage" {
			t.Fatalf("unexpected path: %q", r.URL.Path)
		}
		var payload map[string]any
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode payload: %v", err)
		}
		if payload["chat_id"] != float64(12345) || payload["message_id"] != float64(88) || payload["disable_notification"] != true {
			t.Fatalf("unexpected payload: %#v", payload)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true,"result":true}`))
	}))
	defer server.Close()

	bot := newTestBot(t, token, server.URL, server.Client())
	ok, err := bot.PinChatMessage(context.Background(), PinChatMessageParams{ChatID: ChatIDInt(12345), MessageID: 88, DisableNotification: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Fatal("expected true result")
	}
}

func TestPinChatMessageValidation(t *testing.T) {
	const token = "123:secret"
	bot := newTestBot(t, token, "https://example.test", nil)
	tests := []struct {
		name   string
		params PinChatMessageParams
	}{
		{name: "empty chat", params: PinChatMessageParams{MessageID: 1}},
		{name: "zero message", params: PinChatMessageParams{ChatID: ChatIDInt(123)}},
		{name: "negative message", params: PinChatMessageParams{ChatID: ChatIDInt(123), MessageID: -1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ok, err := bot.PinChatMessage(context.Background(), tt.params)
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

func TestUnpinChatMessageSendsPayloadAndDecodesResult(t *testing.T) {
	const token = "123:secret"
	tests := []struct {
		name          string
		params        UnpinChatMessageParams
		wantMessageID bool
	}{
		{name: "with message", params: UnpinChatMessageParams{ChatID: ChatIDInt(12345), MessageID: 88}, wantMessageID: true},
		{name: "without message", params: UnpinChatMessageParams{ChatID: ChatIDInt(12345)}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path != "/bot"+token+"/unpinChatMessage" {
					t.Fatalf("unexpected path: %q", r.URL.Path)
				}
				var payload map[string]any
				if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
					t.Fatalf("decode payload: %v", err)
				}
				if payload["chat_id"] != float64(12345) {
					t.Fatalf("unexpected chat_id: %#v", payload)
				}
				got, ok := payload["message_id"]
				if tt.wantMessageID {
					if !ok || got != float64(88) {
						t.Fatalf("unexpected message_id: %#v", payload)
					}
				} else if ok {
					t.Fatalf("message_id should be omitted: %#v", payload)
				}
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(`{"ok":true,"result":true}`))
			}))
			defer server.Close()

			bot := newTestBot(t, token, server.URL, server.Client())
			ok, err := bot.UnpinChatMessage(context.Background(), tt.params)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !ok {
				t.Fatal("expected true result")
			}
		})
	}
}

func TestUnpinChatMessageValidation(t *testing.T) {
	const token = "123:secret"
	bot := newTestBot(t, token, "https://example.test", nil)
	tests := []struct {
		name   string
		params UnpinChatMessageParams
	}{
		{name: "empty chat", params: UnpinChatMessageParams{MessageID: 1}},
		{name: "negative message", params: UnpinChatMessageParams{ChatID: ChatIDInt(123), MessageID: -1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ok, err := bot.UnpinChatMessage(context.Background(), tt.params)
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

func TestUnpinAllChatMessagesSendsPayloadAndDecodesResult(t *testing.T) {
	const token = "123:secret"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/bot"+token+"/unpinAllChatMessages" {
			t.Fatalf("unexpected path: %q", r.URL.Path)
		}
		var payload map[string]any
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode payload: %v", err)
		}
		if len(payload) != 1 || payload["chat_id"] != float64(12345) {
			t.Fatalf("unexpected payload: %#v", payload)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true,"result":true}`))
	}))
	defer server.Close()

	bot := newTestBot(t, token, server.URL, server.Client())
	ok, err := bot.UnpinAllChatMessages(context.Background(), UnpinAllChatMessagesParams{ChatID: ChatIDInt(12345)})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Fatal("expected true result")
	}
}

func TestUnpinAllChatMessagesValidation(t *testing.T) {
	const token = "123:secret"
	bot := newTestBot(t, token, "https://example.test", nil)
	ok, err := bot.UnpinAllChatMessages(context.Background(), UnpinAllChatMessagesParams{})
	if err == nil {
		t.Fatal("expected error")
	}
	if ok {
		t.Fatal("expected false result")
	}
	assertNoToken(t, err, token)
}

func TestChatBoolMethodsReturnAPIError(t *testing.T) {
	const token = "123:secret"
	tests := []struct {
		name   string
		method string
		call   func(*Bot) (bool, error)
	}{
		{name: "send chat action", method: "sendChatAction", call: func(bot *Bot) (bool, error) {
			return bot.SendChatAction(context.Background(), SendChatActionParams{ChatID: ChatIDInt(123), Action: ChatActionTyping})
		}},
		{name: "pin", method: "pinChatMessage", call: func(bot *Bot) (bool, error) {
			return bot.PinChatMessage(context.Background(), PinChatMessageParams{ChatID: ChatIDInt(123), MessageID: 1})
		}},
		{name: "unpin", method: "unpinChatMessage", call: func(bot *Bot) (bool, error) {
			return bot.UnpinChatMessage(context.Background(), UnpinChatMessageParams{ChatID: ChatIDInt(123), MessageID: 1})
		}},
		{name: "unpin all", method: "unpinAllChatMessages", call: func(bot *Bot) (bool, error) {
			return bot.UnpinAllChatMessages(context.Background(), UnpinAllChatMessagesParams{ChatID: ChatIDInt(123)})
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

func TestChatBoolMethodsResponseAndContextErrors(t *testing.T) {
	const token = "123:secret"
	tests := []struct {
		name   string
		method string
		call   func(context.Context, *Bot) (bool, error)
	}{
		{name: "send chat action", method: "sendChatAction", call: func(ctx context.Context, bot *Bot) (bool, error) {
			return bot.SendChatAction(ctx, SendChatActionParams{ChatID: ChatIDInt(123), Action: ChatActionTyping})
		}},
		{name: "pin", method: "pinChatMessage", call: func(ctx context.Context, bot *Bot) (bool, error) {
			return bot.PinChatMessage(ctx, PinChatMessageParams{ChatID: ChatIDInt(123), MessageID: 1})
		}},
		{name: "unpin", method: "unpinChatMessage", call: func(ctx context.Context, bot *Bot) (bool, error) {
			return bot.UnpinChatMessage(ctx, UnpinChatMessageParams{ChatID: ChatIDInt(123), MessageID: 1})
		}},
		{name: "unpin all", method: "unpinAllChatMessages", call: func(ctx context.Context, bot *Bot) (bool, error) {
			return bot.UnpinAllChatMessages(ctx, UnpinAllChatMessagesParams{ChatID: ChatIDInt(123)})
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
