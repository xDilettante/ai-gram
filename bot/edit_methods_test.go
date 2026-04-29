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

func TestEditMessageResultUnmarshal(t *testing.T) {
	tests := []struct {
		name      string
		json      string
		wantOK    bool
		wantMsg   bool
		wantError bool
	}{
		{name: "true", json: `true`, wantOK: true},
		{name: "message", json: `{"message_id":42,"chat":{"id":123,"type":"private"},"date":100,"text":"edited"}`, wantOK: true, wantMsg: true},
		{name: "false", json: `false`},
		{name: "invalid", json: `[]`, wantError: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result EditMessageResult
			err := json.Unmarshal([]byte(tt.json), &result)
			if tt.wantError {
				if err == nil {
					t.Fatal("expected error")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if result.IsOK() != tt.wantOK {
				t.Fatalf("unexpected OK: got %v want %v", result.IsOK(), tt.wantOK)
			}
			if result.IsMessage() != tt.wantMsg {
				t.Fatalf("unexpected IsMessage: got %v want %v", result.IsMessage(), tt.wantMsg)
			}
			if tt.wantMsg && result.Message.Text != "edited" {
				t.Fatalf("unexpected message: %+v", result.Message)
			}
		})
	}
}

func TestEditMessageTargetValidation(t *testing.T) {
	valid := []EditMessageTarget{
		EditTargetChat(ChatIDInt(123), 10),
		EditTargetInline("inline-id"),
	}
	for _, target := range valid {
		if err := target.validate(); err != nil {
			t.Fatalf("unexpected valid target error for %+v: %v", target, err)
		}
	}

	tests := []struct {
		name   string
		target EditMessageTarget
	}{
		{name: "empty", target: EditMessageTarget{}},
		{name: "chat without message", target: EditMessageTarget{ChatID: ChatIDInt(123)}},
		{name: "message without chat", target: EditMessageTarget{MessageID: 10}},
		{name: "chat and inline", target: EditMessageTarget{ChatID: ChatIDInt(123), MessageID: 10, InlineMessageID: "inline-id"}},
		{name: "inline empty", target: EditTargetInline("")},
		{name: "negative message", target: EditMessageTarget{ChatID: ChatIDInt(123), MessageID: -1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.target.validate(); err == nil {
				t.Fatal("expected error")
			}
		})
	}
}

func TestEditMessageTextChatTargetSendsPayloadAndDecodesMessage(t *testing.T) {
	const token = "123:secret"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/bot"+token+"/editMessageText" {
			t.Fatalf("unexpected path: %q", r.URL.Path)
		}
		var payload map[string]any
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode payload: %v", err)
		}
		if payload["chat_id"] != float64(12345) || payload["message_id"] != float64(77) {
			t.Fatalf("unexpected target payload: %#v", payload)
		}
		if payload["text"] != "edited" || payload["parse_mode"] != "HTML" || payload["disable_web_page_preview"] != true {
			t.Fatalf("unexpected text payload: %#v", payload)
		}
		reply := payload["reply_markup"].(map[string]any)
		if _, ok := reply["inline_keyboard"]; !ok {
			t.Fatalf("reply_markup missing inline keyboard: %#v", reply)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true,"result":{"message_id":77,"chat":{"id":12345,"type":"private"},"date":100,"text":"edited"}}`))
	}))
	defer server.Close()

	bot := newTestBot(t, token, server.URL, server.Client())
	markup := telegram.NewInlineKeyboard([]telegram.InlineKeyboardButton{telegram.InlineButtonCallback("OK", "ok")})
	result, err := bot.EditMessageText(context.Background(), EditMessageTextParams{
		Target:              EditTargetChat(ChatIDInt(12345), 77),
		Text:                "edited",
		ParseMode:           "HTML",
		LinkPreviewDisabled: true,
		ReplyMarkup:         &markup,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result == nil || !result.IsOK() || !result.IsMessage() || result.Message.Text != "edited" {
		t.Fatalf("unexpected result: %+v", result)
	}
}

func TestEditMessageTextInlineTargetDecodesTrue(t *testing.T) {
	const token = "123:secret"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/bot"+token+"/editMessageText" {
			t.Fatalf("unexpected path: %q", r.URL.Path)
		}
		var payload map[string]any
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode payload: %v", err)
		}
		if payload["inline_message_id"] != "inline-id" {
			t.Fatalf("unexpected inline_message_id: %#v", payload)
		}
		if _, ok := payload["chat_id"]; ok {
			t.Fatalf("chat_id should be omitted: %#v", payload)
		}
		if _, ok := payload["message_id"]; ok {
			t.Fatalf("message_id should be omitted: %#v", payload)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true,"result":true}`))
	}))
	defer server.Close()

	bot := newTestBot(t, token, server.URL, server.Client())
	result, err := bot.EditMessageText(context.Background(), EditMessageTextParams{Target: EditTargetInline("inline-id"), Text: "edited"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result == nil || !result.IsOK() || result.IsMessage() {
		t.Fatalf("unexpected result: %+v", result)
	}
}

func TestEditMessageTextValidation(t *testing.T) {
	const token = "123:secret"
	bot := newTestBot(t, token, "https://example.test", nil)
	markup := telegram.InlineKeyboardMarkup{}
	tests := []struct {
		name   string
		params EditMessageTextParams
	}{
		{name: "empty text", params: EditMessageTextParams{Target: EditTargetChat(ChatIDInt(123), 1)}},
		{name: "invalid target", params: EditMessageTextParams{Text: "edited"}},
		{name: "parse mode with entities", params: EditMessageTextParams{Target: EditTargetChat(ChatIDInt(123), 1), Text: "edited", ParseMode: "HTML", Entities: []telegram.MessageEntity{{Type: telegram.EntityBold, Offset: 0, Length: 1}}}},
		{name: "invalid markup", params: EditMessageTextParams{Target: EditTargetChat(ChatIDInt(123), 1), Text: "edited", ReplyMarkup: &markup}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := bot.EditMessageText(context.Background(), tt.params)
			if err == nil {
				t.Fatal("expected error")
			}
			if result != nil {
				t.Fatalf("expected nil result, got %+v", result)
			}
			assertNoToken(t, err, token)
		})
	}
}

func TestEditMessageTextReturnsAPIError(t *testing.T) {
	const token = "123:secret"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":false,"error_code":400,"description":"Bad Request"}`))
	}))
	defer server.Close()

	bot := newTestBot(t, token, server.URL, server.Client())
	result, err := bot.EditMessageText(context.Background(), EditMessageTextParams{Target: EditTargetChat(ChatIDInt(123), 1), Text: "edited"})
	if err == nil {
		t.Fatal("expected error")
	}
	if result != nil {
		t.Fatalf("expected nil result, got %+v", result)
	}
	var apiErr *apierrors.APIError
	if !stderrors.As(err, &apiErr) {
		t.Fatalf("expected APIError, got %T", err)
	}
	assertNoToken(t, err, token)
}

func TestEditMessageReplyMarkupChatTargetSendsPayloadAndDecodesMessage(t *testing.T) {
	const token = "123:secret"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/bot"+token+"/editMessageReplyMarkup" {
			t.Fatalf("unexpected path: %q", r.URL.Path)
		}
		var payload map[string]any
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode payload: %v", err)
		}
		if payload["chat_id"] != float64(12345) || payload["message_id"] != float64(77) {
			t.Fatalf("unexpected target payload: %#v", payload)
		}
		reply := payload["reply_markup"].(map[string]any)
		if _, ok := reply["inline_keyboard"]; !ok {
			t.Fatalf("reply_markup missing inline keyboard: %#v", reply)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true,"result":{"message_id":77,"chat":{"id":12345,"type":"private"},"date":100,"text":"edited"}}`))
	}))
	defer server.Close()

	bot := newTestBot(t, token, server.URL, server.Client())
	markup := telegram.NewInlineKeyboard([]telegram.InlineKeyboardButton{telegram.InlineButtonCallback("OK", "ok")})
	result, err := bot.EditMessageReplyMarkup(context.Background(), EditMessageReplyMarkupParams{Target: EditTargetChat(ChatIDInt(12345), 77), ReplyMarkup: &markup})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result == nil || !result.IsOK() || !result.IsMessage() {
		t.Fatalf("unexpected result: %+v", result)
	}
}

func TestEditMessageReplyMarkupNilMarkupOmitsReplyMarkup(t *testing.T) {
	const token = "123:secret"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var payload map[string]any
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode payload: %v", err)
		}
		if _, ok := payload["reply_markup"]; ok {
			t.Fatalf("reply_markup should be omitted: %#v", payload)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true,"result":{"message_id":77,"chat":{"id":12345,"type":"private"},"date":100}}`))
	}))
	defer server.Close()

	bot := newTestBot(t, token, server.URL, server.Client())
	result, err := bot.EditMessageReplyMarkup(context.Background(), EditMessageReplyMarkupParams{Target: EditTargetChat(ChatIDInt(12345), 77)})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result == nil || !result.IsOK() || !result.IsMessage() {
		t.Fatalf("unexpected result: %+v", result)
	}
}

func TestEditMessageReplyMarkupInlineTargetDecodesTrue(t *testing.T) {
	const token = "123:secret"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var payload map[string]any
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode payload: %v", err)
		}
		if payload["inline_message_id"] != "inline-id" {
			t.Fatalf("unexpected payload: %#v", payload)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true,"result":true}`))
	}))
	defer server.Close()

	bot := newTestBot(t, token, server.URL, server.Client())
	result, err := bot.EditMessageReplyMarkup(context.Background(), EditMessageReplyMarkupParams{Target: EditTargetInline("inline-id")})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result == nil || !result.IsOK() || result.IsMessage() {
		t.Fatalf("unexpected result: %+v", result)
	}
}

func TestEditMessageReplyMarkupValidation(t *testing.T) {
	const token = "123:secret"
	bot := newTestBot(t, token, "https://example.test", nil)
	markup := telegram.InlineKeyboardMarkup{}
	tests := []struct {
		name   string
		params EditMessageReplyMarkupParams
	}{
		{name: "invalid target", params: EditMessageReplyMarkupParams{}},
		{name: "invalid markup", params: EditMessageReplyMarkupParams{Target: EditTargetChat(ChatIDInt(123), 1), ReplyMarkup: &markup}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := bot.EditMessageReplyMarkup(context.Background(), tt.params)
			if err == nil {
				t.Fatal("expected error")
			}
			if result != nil {
				t.Fatalf("expected nil result, got %+v", result)
			}
			assertNoToken(t, err, token)
		})
	}
}

func TestEditMessageReplyMarkupReturnsAPIError(t *testing.T) {
	const token = "123:secret"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":false,"error_code":400,"description":"Bad Request"}`))
	}))
	defer server.Close()

	bot := newTestBot(t, token, server.URL, server.Client())
	result, err := bot.EditMessageReplyMarkup(context.Background(), EditMessageReplyMarkupParams{Target: EditTargetChat(ChatIDInt(123), 1)})
	if err == nil {
		t.Fatal("expected error")
	}
	if result != nil {
		t.Fatalf("expected nil result, got %+v", result)
	}
	var apiErr *apierrors.APIError
	if !stderrors.As(err, &apiErr) {
		t.Fatalf("expected APIError, got %T", err)
	}
	assertNoToken(t, err, token)
}

func TestEditMessageResponseAndContextErrors(t *testing.T) {
	const token = "123:secret"

	t.Run("cancelled context", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			t.Fatal("request should not reach server")
		}))
		defer server.Close()
		bot := newTestBot(t, token, server.URL, server.Client())
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		result, err := bot.EditMessageText(ctx, EditMessageTextParams{Target: EditTargetChat(ChatIDInt(123), 1), Text: "edited"})
		if err == nil {
			t.Fatal("expected error")
		}
		if result != nil {
			t.Fatalf("expected nil result, got %+v", result)
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
		result, err := bot.EditMessageText(context.Background(), EditMessageTextParams{Target: EditTargetChat(ChatIDInt(123), 1), Text: "edited"})
		if err == nil {
			t.Fatal("expected error")
		}
		if result != nil {
			t.Fatalf("expected nil result, got %+v", result)
		}
		assertNoToken(t, err, token)
	})

	t.Run("http status", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "server error", http.StatusInternalServerError)
		}))
		defer server.Close()
		bot := newTestBot(t, token, server.URL, server.Client())
		result, err := bot.EditMessageReplyMarkup(context.Background(), EditMessageReplyMarkupParams{Target: EditTargetChat(ChatIDInt(123), 1)})
		if err == nil {
			t.Fatal("expected error")
		}
		if result != nil {
			t.Fatalf("expected nil result, got %+v", result)
		}
		assertNoToken(t, err, token)
	})
}

func TestEditMessageCaptionChatTargetSendsPayloadAndDecodesMessage(t *testing.T) {
	const token = "123:secret"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/bot"+token+"/editMessageCaption" {
			t.Fatalf("unexpected path: %q", r.URL.Path)
		}
		var payload map[string]any
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode payload: %v", err)
		}
		if payload["chat_id"] != float64(12345) || payload["message_id"] != float64(77) {
			t.Fatalf("unexpected target payload: %#v", payload)
		}
		if payload["caption"] != "edited caption" || payload["parse_mode"] != "HTML" {
			t.Fatalf("unexpected caption payload: %#v", payload)
		}
		reply := payload["reply_markup"].(map[string]any)
		if _, ok := reply["inline_keyboard"]; !ok {
			t.Fatalf("reply_markup missing inline keyboard: %#v", reply)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true,"result":{"message_id":77,"chat":{"id":12345,"type":"private"},"date":100,"caption":"edited caption"}}`))
	}))
	defer server.Close()

	bot := newTestBot(t, token, server.URL, server.Client())
	markup := telegram.NewInlineKeyboard([]telegram.InlineKeyboardButton{telegram.InlineButtonCallback("OK", "ok")})
	result, err := bot.EditMessageCaption(context.Background(), EditMessageCaptionParams{
		Target:      EditTargetChat(ChatIDInt(12345), 77),
		Caption:     "edited caption",
		ParseMode:   "HTML",
		ReplyMarkup: &markup,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result == nil || !result.IsOK() || !result.IsMessage() || result.Message.Caption != "edited caption" {
		t.Fatalf("unexpected result: %+v", result)
	}
}

func TestEditMessageCaptionInlineTargetDecodesTrue(t *testing.T) {
	const token = "123:secret"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/bot"+token+"/editMessageCaption" {
			t.Fatalf("unexpected path: %q", r.URL.Path)
		}
		var payload map[string]any
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode payload: %v", err)
		}
		if payload["inline_message_id"] != "inline-id" {
			t.Fatalf("unexpected inline_message_id: %#v", payload)
		}
		if _, ok := payload["chat_id"]; ok {
			t.Fatalf("chat_id should be omitted: %#v", payload)
		}
		if _, ok := payload["message_id"]; ok {
			t.Fatalf("message_id should be omitted: %#v", payload)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true,"result":true}`))
	}))
	defer server.Close()

	bot := newTestBot(t, token, server.URL, server.Client())
	result, err := bot.EditMessageCaption(context.Background(), EditMessageCaptionParams{Target: EditTargetInline("inline-id"), Caption: "edited caption"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result == nil || !result.IsOK() || result.IsMessage() {
		t.Fatalf("unexpected result: %+v", result)
	}
}

func TestEditMessageCaptionAllowsEmptyCaption(t *testing.T) {
	const token = "123:secret"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var payload map[string]any
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode payload: %v", err)
		}
		if got, ok := payload["caption"]; !ok || got != "" {
			t.Fatalf("empty caption should be sent to remove caption: %#v", payload)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true,"result":{"message_id":77,"chat":{"id":12345,"type":"private"},"date":100}}`))
	}))
	defer server.Close()

	bot := newTestBot(t, token, server.URL, server.Client())
	result, err := bot.EditMessageCaption(context.Background(), EditMessageCaptionParams{Target: EditTargetChat(ChatIDInt(12345), 77)})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result == nil || !result.IsOK() || !result.IsMessage() {
		t.Fatalf("unexpected result: %+v", result)
	}
}

func TestEditMessageCaptionValidation(t *testing.T) {
	const token = "123:secret"
	bot := newTestBot(t, token, "https://example.test", nil)
	markup := telegram.InlineKeyboardMarkup{}
	tests := []struct {
		name   string
		params EditMessageCaptionParams
	}{
		{name: "invalid target", params: EditMessageCaptionParams{Caption: "caption"}},
		{name: "parse mode with caption entities", params: EditMessageCaptionParams{Target: EditTargetChat(ChatIDInt(123), 1), Caption: "caption", ParseMode: "HTML", CaptionEntities: []telegram.MessageEntity{{Type: telegram.EntityBold, Offset: 0, Length: 1}}}},
		{name: "invalid markup", params: EditMessageCaptionParams{Target: EditTargetChat(ChatIDInt(123), 1), Caption: "caption", ReplyMarkup: &markup}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := bot.EditMessageCaption(context.Background(), tt.params)
			if err == nil {
				t.Fatal("expected error")
			}
			if result != nil {
				t.Fatalf("expected nil result, got %+v", result)
			}
			assertNoToken(t, err, token)
		})
	}
}

func TestEditMessageCaptionReturnsAPIError(t *testing.T) {
	const token = "123:secret"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":false,"error_code":400,"description":"Bad Request"}`))
	}))
	defer server.Close()

	bot := newTestBot(t, token, server.URL, server.Client())
	result, err := bot.EditMessageCaption(context.Background(), EditMessageCaptionParams{Target: EditTargetChat(ChatIDInt(123), 1), Caption: "caption"})
	if err == nil {
		t.Fatal("expected error")
	}
	if result != nil {
		t.Fatalf("expected nil result, got %+v", result)
	}
	var apiErr *apierrors.APIError
	if !stderrors.As(err, &apiErr) {
		t.Fatalf("expected APIError, got %T", err)
	}
	assertNoToken(t, err, token)
}

func TestEditMessageCaptionResponseAndContextErrors(t *testing.T) {
	const token = "123:secret"

	t.Run("cancelled context", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			t.Fatal("request should not reach server")
		}))
		defer server.Close()
		bot := newTestBot(t, token, server.URL, server.Client())
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		result, err := bot.EditMessageCaption(ctx, EditMessageCaptionParams{Target: EditTargetChat(ChatIDInt(123), 1), Caption: "caption"})
		if err == nil {
			t.Fatal("expected error")
		}
		if result != nil {
			t.Fatalf("expected nil result, got %+v", result)
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
		result, err := bot.EditMessageCaption(context.Background(), EditMessageCaptionParams{Target: EditTargetChat(ChatIDInt(123), 1), Caption: "caption"})
		if err == nil {
			t.Fatal("expected error")
		}
		if result != nil {
			t.Fatalf("expected nil result, got %+v", result)
		}
		assertNoToken(t, err, token)
	})

	t.Run("http status", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "server error", http.StatusInternalServerError)
		}))
		defer server.Close()
		bot := newTestBot(t, token, server.URL, server.Client())
		result, err := bot.EditMessageCaption(context.Background(), EditMessageCaptionParams{Target: EditTargetChat(ChatIDInt(123), 1), Caption: "caption"})
		if err == nil {
			t.Fatal("expected error")
		}
		if result != nil {
			t.Fatalf("expected nil result, got %+v", result)
		}
		assertNoToken(t, err, token)
	})
}
