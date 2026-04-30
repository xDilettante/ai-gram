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

func TestSetMessageReactionSendsPayloadAndDecodesResult(t *testing.T) {
	const token = "123:secret"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/bot"+token+"/setMessageReaction" {
			t.Fatalf("unexpected path: %q", r.URL.Path)
		}
		var payload map[string]any
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode payload: %v", err)
		}
		if payload["chat_id"] != float64(12345) || payload["message_id"] != float64(777) || payload["is_big"] != true {
			t.Fatalf("unexpected payload: %#v", payload)
		}
		reactions, ok := payload["reaction"].([]any)
		if !ok || len(reactions) != 2 {
			t.Fatalf("unexpected reactions: %#v", payload["reaction"])
		}
		first, _ := reactions[0].(map[string]any)
		second, _ := reactions[1].(map[string]any)
		if first["type"] != "emoji" || first["emoji"] != "👍" || second["type"] != "custom_emoji" || second["custom_emoji_id"] != "custom-id" {
			t.Fatalf("unexpected reaction payload: %#v", reactions)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true,"result":true}`))
	}))
	defer server.Close()

	bot := newTestBot(t, token, server.URL, server.Client())
	ok, err := bot.SetMessageReaction(context.Background(), SetMessageReactionParams{
		ChatID:    ChatIDInt(12345),
		MessageID: 777,
		Reaction: []telegram.ReactionType{
			telegram.NewReactionTypeEmoji("👍"),
			telegram.NewReactionTypeCustomEmoji("custom-id"),
		},
		IsBig: true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Fatal("expected true result")
	}
}

func TestSetMessageReactionAllowsEmptyReaction(t *testing.T) {
	const token = "123:secret"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/bot"+token+"/setMessageReaction" {
			t.Fatalf("unexpected path: %q", r.URL.Path)
		}
		var payload map[string]any
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode payload: %v", err)
		}
		if _, ok := payload["reaction"]; ok {
			t.Fatalf("empty reaction should be omitted: %#v", payload)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true,"result":true}`))
	}))
	defer server.Close()

	bot := newTestBot(t, token, server.URL, server.Client())
	ok, err := bot.SetMessageReaction(context.Background(), SetMessageReactionParams{ChatID: ChatIDInt(12345), MessageID: 777})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Fatal("expected true result")
	}
}

func TestSetMessageReactionValidation(t *testing.T) {
	const token = "123:secret"
	bot := newTestBot(t, token, "https://example.test", nil)
	var typedNil *telegram.ReactionTypeEmoji
	tests := []struct {
		name   string
		params SetMessageReactionParams
	}{
		{name: "empty chat", params: SetMessageReactionParams{MessageID: 1}},
		{name: "zero message", params: SetMessageReactionParams{ChatID: ChatIDInt(123)}},
		{name: "negative message", params: SetMessageReactionParams{ChatID: ChatIDInt(123), MessageID: -1}},
		{name: "empty emoji", params: SetMessageReactionParams{ChatID: ChatIDInt(123), MessageID: 1, Reaction: []telegram.ReactionType{telegram.NewReactionTypeEmoji("")}}},
		{name: "typed nil reaction", params: SetMessageReactionParams{ChatID: ChatIDInt(123), MessageID: 1, Reaction: []telegram.ReactionType{typedNil}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ok, err := bot.SetMessageReaction(context.Background(), tt.params)
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

func TestSetMessageReactionReturnsAPIError(t *testing.T) {
	const token = "123:secret"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/bot"+token+"/setMessageReaction" {
			t.Fatalf("unexpected path: %q", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":false,"error_code":400,"description":"Bad Request"}`))
	}))
	defer server.Close()

	bot := newTestBot(t, token, server.URL, server.Client())
	ok, err := bot.SetMessageReaction(context.Background(), SetMessageReactionParams{ChatID: ChatIDInt(123), MessageID: 1})
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
}

func TestSetMessageReactionResponseAndContextErrors(t *testing.T) {
	const token = "123:secret"
	t.Run("cancelled context", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			t.Fatal("request should not reach server")
		}))
		defer server.Close()
		bot := newTestBot(t, token, server.URL, server.Client())
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		ok, err := bot.SetMessageReaction(ctx, SetMessageReactionParams{ChatID: ChatIDInt(123), MessageID: 1})
		if err == nil {
			t.Fatal("expected error")
		}
		if ok {
			t.Fatal("expected false result")
		}
		assertNoToken(t, err, token)
	})

	t.Run("invalid json", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/bot"+token+"/setMessageReaction" {
				t.Fatalf("unexpected path: %q", r.URL.Path)
			}
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`not-json`))
		}))
		defer server.Close()
		bot := newTestBot(t, token, server.URL, server.Client())
		ok, err := bot.SetMessageReaction(context.Background(), SetMessageReactionParams{ChatID: ChatIDInt(123), MessageID: 1})
		if err == nil {
			t.Fatal("expected error")
		}
		if ok {
			t.Fatal("expected false result")
		}
		assertNoToken(t, err, token)
	})

	t.Run("http status", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/bot"+token+"/setMessageReaction" {
				t.Fatalf("unexpected path: %q", r.URL.Path)
			}
			http.Error(w, "server error", http.StatusInternalServerError)
		}))
		defer server.Close()
		bot := newTestBot(t, token, server.URL, server.Client())
		ok, err := bot.SetMessageReaction(context.Background(), SetMessageReactionParams{ChatID: ChatIDInt(123), MessageID: 1})
		if err == nil {
			t.Fatal("expected error")
		}
		if ok {
			t.Fatal("expected false result")
		}
		assertNoToken(t, err, token)
	})
}
