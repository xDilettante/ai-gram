package bot

import (
	"context"
	"encoding/json"
	stderrors "errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	apierrors "ai-gram/errors"
)

func TestGetMeDecodesUser(t *testing.T) {
	const token = "123:secret"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/bot"+token+"/getMe" {
			t.Fatalf("unexpected path: %q", r.URL.Path)
		}
		var payload map[string]any
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode payload: %v", err)
		}
		if len(payload) != 0 {
			t.Fatalf("unexpected payload: %#v", payload)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true,"result":{"id":42,"is_bot":true,"first_name":"AiGram","username":"ai_gram_bot"}}`))
	}))
	defer server.Close()

	bot := newTestBot(t, token, server.URL, server.Client())
	user, err := bot.GetMe(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if user == nil {
		t.Fatal("expected user")
	}
	if user.ID != 42 || !user.IsBot || user.FirstName != "AiGram" || user.Username != "ai_gram_bot" {
		t.Fatalf("unexpected user: %+v", user)
	}
}

func TestGetMeReturnsAPIError(t *testing.T) {
	const token = "123:secret"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":false,"error_code":401,"description":"Unauthorized"}`))
	}))
	defer server.Close()

	bot := newTestBot(t, token, server.URL, server.Client())
	user, err := bot.GetMe(context.Background())
	if err == nil {
		t.Fatal("expected error")
	}
	if user != nil {
		t.Fatalf("expected nil user, got %+v", user)
	}

	var apiErr *apierrors.APIError
	if !stderrors.As(err, &apiErr) {
		t.Fatalf("expected APIError, got %T", err)
	}
	if apiErr.Code != 401 {
		t.Fatalf("unexpected APIError code: %d", apiErr.Code)
	}
	assertNoToken(t, err, token)
}

func TestSendMessageSendsNumericChatIDAndDecodesMessage(t *testing.T) {
	const token = "123:secret"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/bot"+token+"/sendMessage" {
			t.Fatalf("unexpected path: %q", r.URL.Path)
		}

		var payload map[string]any
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode payload: %v", err)
		}
		if got := payload["chat_id"]; got != float64(12345) {
			t.Fatalf("unexpected chat_id: %#v", got)
		}
		if got := payload["text"]; got != "hello" {
			t.Fatalf("unexpected text: %#v", got)
		}
		if _, ok := payload["disable_notification"]; ok {
			t.Fatal("disable_notification should be omitted when false")
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true,"result":{"message_id":7,"chat":{"id":12345,"type":"private"},"date":100,"text":"hello"}}`))
	}))
	defer server.Close()

	bot := newTestBot(t, token, server.URL, server.Client())
	message, err := bot.SendMessage(context.Background(), SendMessageParams{ChatID: ChatIDInt(12345), Text: "hello"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if message == nil {
		t.Fatal("expected message")
	}
	if message.MessageID != 7 || message.Chat.ID != 12345 || message.Date != 100 || message.Text != "hello" {
		t.Fatalf("unexpected message: %+v", message)
	}
}

func TestSendMessageSendsStringChatID(t *testing.T) {
	const token = "123:secret"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var payload map[string]any
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode payload: %v", err)
		}
		if got := payload["chat_id"]; got != "@channel" {
			t.Fatalf("unexpected chat_id: %#v", got)
		}
		if got := payload["parse_mode"]; got != "HTML" {
			t.Fatalf("unexpected parse_mode: %#v", got)
		}
		if got := payload["disable_notification"]; got != true {
			t.Fatalf("unexpected disable_notification: %#v", got)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true,"result":{"message_id":8,"chat":{"id":-100,"type":"channel","username":"channel"},"date":101,"text":"hello"}}`))
	}))
	defer server.Close()

	bot := newTestBot(t, token, server.URL, server.Client())
	message, err := bot.SendMessage(context.Background(), SendMessageParams{
		ChatID:              ChatIDString("@channel"),
		Text:                "hello",
		ParseMode:           "HTML",
		DisableNotification: true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if message == nil || message.MessageID != 8 {
		t.Fatalf("unexpected message: %+v", message)
	}
}

func TestSendMessageRejectsEmptyChatID(t *testing.T) {
	const token = "123:secret"
	bot := newTestBot(t, token, "https://example.test", nil)

	message, err := bot.SendMessage(context.Background(), SendMessageParams{Text: "hello"})
	if err == nil {
		t.Fatal("expected error")
	}
	if message != nil {
		t.Fatalf("expected nil message, got %+v", message)
	}
	if !strings.Contains(err.Error(), "chat_id") {
		t.Fatalf("unexpected error: %q", err.Error())
	}
	assertNoToken(t, err, token)
}

func TestSendMessageRejectsEmptyText(t *testing.T) {
	const token = "123:secret"
	bot := newTestBot(t, token, "https://example.test", nil)

	message, err := bot.SendMessage(context.Background(), SendMessageParams{ChatID: ChatIDInt(12345)})
	if err == nil {
		t.Fatal("expected error")
	}
	if message != nil {
		t.Fatalf("expected nil message, got %+v", message)
	}
	if !strings.Contains(err.Error(), "text") {
		t.Fatalf("unexpected error: %q", err.Error())
	}
	assertNoToken(t, err, token)
}

func TestSendMessageReturnsAPIError(t *testing.T) {
	const token = "123:secret"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":false,"error_code":400,"description":"Bad Request"}`))
	}))
	defer server.Close()

	bot := newTestBot(t, token, server.URL, server.Client())
	message, err := bot.SendMessage(context.Background(), SendMessageParams{ChatID: ChatIDInt(12345), Text: "hello"})
	if err == nil {
		t.Fatal("expected error")
	}
	if message != nil {
		t.Fatalf("expected nil message, got %+v", message)
	}

	var apiErr *apierrors.APIError
	if !stderrors.As(err, &apiErr) {
		t.Fatalf("expected APIError, got %T", err)
	}
	if apiErr.Code != 400 || apiErr.Description != "Bad Request" {
		t.Fatalf("unexpected APIError: %+v", apiErr)
	}
	assertNoToken(t, err, token)
}

func TestSendMessageRedactsTokenFromAPIErrorDescription(t *testing.T) {
	const token = "123:secret"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":false,"error_code":400,"description":"Bad token 123:secret"}`))
	}))
	defer server.Close()

	bot := newTestBot(t, token, server.URL, server.Client())
	_, err := bot.SendMessage(context.Background(), SendMessageParams{ChatID: ChatIDInt(12345), Text: "hello"})
	if err == nil {
		t.Fatal("expected error")
	}
	assertNoToken(t, err, token)
}
