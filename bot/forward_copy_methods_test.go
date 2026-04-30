package bot

import (
	"context"
	"encoding/json"
	stderrors "errors"
	"net/http"
	"net/http/httptest"
	"testing"

	apierrors "ai-gram/errors"
	"ai-gram/telegram"
)

func TestForwardMessageSendsPayloadAndDecodesMessage(t *testing.T) {
	const token = "123:secret"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/bot"+token+"/forwardMessage" {
			t.Fatalf("unexpected path: %q", r.URL.Path)
		}
		var payload map[string]any
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode payload: %v", err)
		}
		if payload["chat_id"] != float64(12345) || payload["from_chat_id"] != float64(67890) || payload["message_id"] != float64(77) {
			t.Fatalf("unexpected message payload: %#v", payload)
		}
		if payload["message_thread_id"] != float64(11) || payload["disable_notification"] != true || payload["protect_content"] != true {
			t.Fatalf("unexpected option payload: %#v", payload)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true,"result":{"message_id":88,"chat":{"id":12345,"type":"private"},"date":100,"text":"forwarded"}}`))
	}))
	defer server.Close()

	bot := newTestBot(t, token, server.URL, server.Client())
	message, err := bot.ForwardMessage(context.Background(), ForwardMessageParams{
		ChatID:              ChatIDInt(12345),
		FromChatID:          ChatIDInt(67890),
		MessageID:           77,
		MessageThreadID:     11,
		DisableNotification: true,
		ProtectContent:      true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if message == nil || message.MessageID != 88 || message.Text != "forwarded" {
		t.Fatalf("unexpected message: %+v", message)
	}
}

func TestForwardMessageValidation(t *testing.T) {
	const token = "123:secret"
	bot := newTestBot(t, token, "https://example.test", nil)
	tests := []struct {
		name   string
		params ForwardMessageParams
	}{
		{name: "empty chat", params: ForwardMessageParams{FromChatID: ChatIDInt(1), MessageID: 1}},
		{name: "empty from chat", params: ForwardMessageParams{ChatID: ChatIDInt(1), MessageID: 1}},
		{name: "zero message", params: ForwardMessageParams{ChatID: ChatIDInt(1), FromChatID: ChatIDInt(2)}},
		{name: "negative message", params: ForwardMessageParams{ChatID: ChatIDInt(1), FromChatID: ChatIDInt(2), MessageID: -1}},
		{name: "negative thread", params: ForwardMessageParams{ChatID: ChatIDInt(1), FromChatID: ChatIDInt(2), MessageID: 1, MessageThreadID: -1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			message, err := bot.ForwardMessage(context.Background(), tt.params)
			if err == nil {
				t.Fatal("expected error")
			}
			if message != nil {
				t.Fatalf("expected nil message, got %+v", message)
			}
			assertNoToken(t, err, token)
		})
	}
}

func TestForwardMessageReturnsAPIError(t *testing.T) {
	const token = "123:secret"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":false,"error_code":400,"description":"Bad Request"}`))
	}))
	defer server.Close()

	bot := newTestBot(t, token, server.URL, server.Client())
	message, err := bot.ForwardMessage(context.Background(), ForwardMessageParams{ChatID: ChatIDInt(123), FromChatID: ChatIDInt(456), MessageID: 1})
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
	assertNoToken(t, err, token)
}

func TestForwardMessageResponseAndContextErrors(t *testing.T) {
	const token = "123:secret"

	t.Run("cancelled context", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			t.Fatal("request should not reach server")
		}))
		defer server.Close()
		bot := newTestBot(t, token, server.URL, server.Client())
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		message, err := bot.ForwardMessage(ctx, ForwardMessageParams{ChatID: ChatIDInt(123), FromChatID: ChatIDInt(456), MessageID: 1})
		if err == nil {
			t.Fatal("expected error")
		}
		if message != nil {
			t.Fatalf("expected nil message, got %+v", message)
		}
		assertNoToken(t, err, token)
	})

	t.Run("invalid json", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`not-json`))
		}))
		defer server.Close()
		bot := newTestBot(t, token, server.URL, server.Client())
		message, err := bot.ForwardMessage(context.Background(), ForwardMessageParams{ChatID: ChatIDInt(123), FromChatID: ChatIDInt(456), MessageID: 1})
		if err == nil {
			t.Fatal("expected error")
		}
		if message != nil {
			t.Fatalf("expected nil message, got %+v", message)
		}
		assertNoToken(t, err, token)
	})

	t.Run("http status", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "server error", http.StatusInternalServerError)
		}))
		defer server.Close()
		bot := newTestBot(t, token, server.URL, server.Client())
		message, err := bot.ForwardMessage(context.Background(), ForwardMessageParams{ChatID: ChatIDInt(123), FromChatID: ChatIDInt(456), MessageID: 1})
		if err == nil {
			t.Fatal("expected error")
		}
		if message != nil {
			t.Fatalf("expected nil message, got %+v", message)
		}
		assertNoToken(t, err, token)
	})
}

func TestCopyMessageSendsPayloadAndDecodesMessageID(t *testing.T) {
	const token = "123:secret"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/bot"+token+"/copyMessage" {
			t.Fatalf("unexpected path: %q", r.URL.Path)
		}
		var payload map[string]any
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode payload: %v", err)
		}
		if payload["chat_id"] != float64(12345) || payload["from_chat_id"] != float64(67890) || payload["message_id"] != float64(77) {
			t.Fatalf("unexpected message payload: %#v", payload)
		}
		if payload["message_thread_id"] != float64(11) || payload["caption"] != "copied" || payload["parse_mode"] != "HTML" {
			t.Fatalf("unexpected text payload: %#v", payload)
		}
		if payload["disable_notification"] != true || payload["protect_content"] != true {
			t.Fatalf("unexpected option payload: %#v", payload)
		}
		reply, ok := payload["reply_parameters"].(map[string]any)
		if !ok || reply["message_id"] != float64(42) || reply["allow_sending_without_reply"] != true {
			t.Fatalf("unexpected reply_parameters: %#v", payload["reply_parameters"])
		}
		markup, ok := payload["reply_markup"].(map[string]any)
		if !ok {
			t.Fatalf("reply_markup missing: %#v", payload["reply_markup"])
		}
		if _, ok := markup["inline_keyboard"]; !ok {
			t.Fatalf("reply_markup missing inline keyboard: %#v", markup)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true,"result":{"message_id":99}}`))
	}))
	defer server.Close()

	bot := newTestBot(t, token, server.URL, server.Client())
	markup := telegram.NewInlineKeyboard([]telegram.InlineKeyboardButton{telegram.InlineButtonCallback("OK", "ok")})
	messageID, err := bot.CopyMessage(context.Background(), CopyMessageParams{
		ChatID:              ChatIDInt(12345),
		FromChatID:          ChatIDInt(67890),
		MessageID:           77,
		MessageThreadID:     11,
		Caption:             "copied",
		ParseMode:           "HTML",
		DisableNotification: true,
		ProtectContent:      true,
		ReplyParameters:     &telegram.ReplyParameters{MessageID: 42, AllowSendingWithoutReply: true},
		ReplyMarkup:         markup,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if messageID == nil || messageID.MessageID != 99 {
		t.Fatalf("unexpected message id: %+v", messageID)
	}
}

func TestCopyMessageValidation(t *testing.T) {
	const token = "123:secret"
	bot := newTestBot(t, token, "https://example.test", nil)
	invalidMarkup := telegram.InlineKeyboardMarkup{}
	tests := []struct {
		name   string
		params CopyMessageParams
	}{
		{name: "empty chat", params: CopyMessageParams{FromChatID: ChatIDInt(1), MessageID: 1}},
		{name: "empty from chat", params: CopyMessageParams{ChatID: ChatIDInt(1), MessageID: 1}},
		{name: "zero message", params: CopyMessageParams{ChatID: ChatIDInt(1), FromChatID: ChatIDInt(2)}},
		{name: "negative message", params: CopyMessageParams{ChatID: ChatIDInt(1), FromChatID: ChatIDInt(2), MessageID: -1}},
		{name: "negative thread", params: CopyMessageParams{ChatID: ChatIDInt(1), FromChatID: ChatIDInt(2), MessageID: 1, MessageThreadID: -1}},
		{name: "parse mode with caption entities", params: CopyMessageParams{ChatID: ChatIDInt(1), FromChatID: ChatIDInt(2), MessageID: 1, ParseMode: "HTML", CaptionEntities: []telegram.MessageEntity{{Type: telegram.EntityBold, Offset: 0, Length: 1}}}},
		{name: "invalid reply parameters", params: CopyMessageParams{ChatID: ChatIDInt(1), FromChatID: ChatIDInt(2), MessageID: 1, ReplyParameters: &telegram.ReplyParameters{}}},
		{name: "invalid reply markup", params: CopyMessageParams{ChatID: ChatIDInt(1), FromChatID: ChatIDInt(2), MessageID: 1, ReplyMarkup: invalidMarkup}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			messageID, err := bot.CopyMessage(context.Background(), tt.params)
			if err == nil {
				t.Fatal("expected error")
			}
			if messageID != nil {
				t.Fatalf("expected nil message id, got %+v", messageID)
			}
			assertNoToken(t, err, token)
		})
	}
}

func TestCopyMessageReturnsAPIError(t *testing.T) {
	const token = "123:secret"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":false,"error_code":400,"description":"Bad Request"}`))
	}))
	defer server.Close()

	bot := newTestBot(t, token, server.URL, server.Client())
	messageID, err := bot.CopyMessage(context.Background(), CopyMessageParams{ChatID: ChatIDInt(123), FromChatID: ChatIDInt(456), MessageID: 1})
	if err == nil {
		t.Fatal("expected error")
	}
	if messageID != nil {
		t.Fatalf("expected nil message id, got %+v", messageID)
	}
	var apiErr *apierrors.APIError
	if !stderrors.As(err, &apiErr) {
		t.Fatalf("expected APIError, got %T", err)
	}
	assertNoToken(t, err, token)
}

func TestCopyMessageResponseAndContextErrors(t *testing.T) {
	const token = "123:secret"

	t.Run("cancelled context", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			t.Fatal("request should not reach server")
		}))
		defer server.Close()
		bot := newTestBot(t, token, server.URL, server.Client())
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		messageID, err := bot.CopyMessage(ctx, CopyMessageParams{ChatID: ChatIDInt(123), FromChatID: ChatIDInt(456), MessageID: 1})
		if err == nil {
			t.Fatal("expected error")
		}
		if messageID != nil {
			t.Fatalf("expected nil message id, got %+v", messageID)
		}
		assertNoToken(t, err, token)
	})

	t.Run("invalid json", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`not-json`))
		}))
		defer server.Close()
		bot := newTestBot(t, token, server.URL, server.Client())
		messageID, err := bot.CopyMessage(context.Background(), CopyMessageParams{ChatID: ChatIDInt(123), FromChatID: ChatIDInt(456), MessageID: 1})
		if err == nil {
			t.Fatal("expected error")
		}
		if messageID != nil {
			t.Fatalf("expected nil message id, got %+v", messageID)
		}
		assertNoToken(t, err, token)
	})

	t.Run("http status", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "server error", http.StatusInternalServerError)
		}))
		defer server.Close()
		bot := newTestBot(t, token, server.URL, server.Client())
		messageID, err := bot.CopyMessage(context.Background(), CopyMessageParams{ChatID: ChatIDInt(123), FromChatID: ChatIDInt(456), MessageID: 1})
		if err == nil {
			t.Fatal("expected error")
		}
		if messageID != nil {
			t.Fatalf("expected nil message id, got %+v", messageID)
		}
		assertNoToken(t, err, token)
	})
}
