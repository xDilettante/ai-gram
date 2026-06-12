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

func TestSendRichMessageSendsPayloadAndDecodesMessage(t *testing.T) {
	const token = "123:secret"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/bot"+token+"/sendRichMessage" {
			t.Fatalf("unexpected path: %q", r.URL.Path)
		}
		var payload map[string]any
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode payload: %v", err)
		}
		if payload["business_connection_id"] != "bc-1" || payload["chat_id"] != float64(12345) || payload["message_thread_id"] != float64(7) || payload["direct_messages_topic_id"] != float64(8) {
			t.Fatalf("unexpected routing payload: %#v", payload)
		}
		if payload["disable_notification"] != true || payload["protect_content"] != true || payload["allow_paid_broadcast"] != true || payload["message_effect_id"] != "effect-id" {
			t.Fatalf("unexpected message options: %#v", payload)
		}
		rich, ok := payload["rich_message"].(map[string]any)
		if !ok || rich["html"] != "<b>Hello</b>" || rich["is_rtl"] != true || rich["skip_entity_detection"] != true {
			t.Fatalf("unexpected rich_message: %#v", payload["rich_message"])
		}
		reply := payload["reply_markup"].(map[string]any)
		if _, ok := reply["inline_keyboard"]; !ok {
			t.Fatalf("reply_markup missing inline keyboard: %#v", reply)
		}
		if _, ok := payload["reply_parameters"].(map[string]any); !ok {
			t.Fatalf("missing reply_parameters: %#v", payload)
		}
		if _, ok := payload["suggested_post_parameters"].(map[string]any); !ok {
			t.Fatalf("missing suggested_post_parameters: %#v", payload)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true,"result":{"message_id":10,"chat":{"id":12345,"type":"private"},"date":100,"rich_message":{"blocks":[{"type":"paragraph","text":"Hello"}]}}}`))
	}))
	defer server.Close()

	bot := newTestBot(t, token, server.URL, server.Client())
	markup := telegram.NewInlineKeyboard([]telegram.InlineKeyboardButton{telegram.InlineButtonCallback("OK", "ok")})
	message, err := bot.SendRichMessage(context.Background(), SendRichMessageParams{
		BusinessConnectionID:  "bc-1",
		ChatID:                ChatIDInt(12345),
		MessageThreadID:       7,
		DirectMessagesTopicID: 8,
		RichMessage:           InputRichMessage{HTML: "<b>Hello</b>", IsRTL: true, SkipEntityDetection: true},
		DisableNotification:   true,
		ProtectContent:        true,
		AllowPaidBroadcast:    true,
		MessageEffectID:       "effect-id",
		SuggestedPostParameters: &telegram.SuggestedPostParameters{
			Price:    &telegram.SuggestedPostPrice{Currency: "XTR", Amount: 1},
			SendDate: 123,
		},
		ReplyParameters: &telegram.ReplyParameters{MessageID: 9},
		ReplyMarkup:     markup,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if message == nil || message.RichMessage == nil || len(message.RichMessage.Blocks) != 1 {
		t.Fatalf("unexpected message: %+v", message)
	}
}

func TestSendRichMessageDraftSendsPayloadAndDecodesResult(t *testing.T) {
	const token = "123:secret"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/bot"+token+"/sendRichMessageDraft" {
			t.Fatalf("unexpected path: %q", r.URL.Path)
		}
		var payload map[string]any
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode payload: %v", err)
		}
		if payload["chat_id"] != float64(12345) || payload["message_thread_id"] != float64(7) || payload["draft_id"] != float64(99) {
			t.Fatalf("unexpected draft payload: %#v", payload)
		}
		rich, ok := payload["rich_message"].(map[string]any)
		if !ok || rich["markdown"] != "**Working**" {
			t.Fatalf("unexpected rich_message: %#v", payload["rich_message"])
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true,"result":true}`))
	}))
	defer server.Close()

	bot := newTestBot(t, token, server.URL, server.Client())
	ok, err := bot.SendRichMessageDraft(context.Background(), SendRichMessageDraftParams{
		ChatID:          ChatIDInt(12345),
		MessageThreadID: 7,
		DraftID:         99,
		RichMessage:     InputRichMessage{Markdown: "**Working**"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Fatal("expected true result")
	}
}

func TestRichMessageMethodValidation(t *testing.T) {
	const token = "123:secret"
	bot := newTestBot(t, token, "https://example.test", nil)
	markup := telegram.InlineKeyboardMarkup{}

	sendTests := []struct {
		name   string
		params SendRichMessageParams
	}{
		{name: "missing chat", params: SendRichMessageParams{RichMessage: InputRichHTML("<b>Hello</b>")}},
		{name: "missing rich content", params: SendRichMessageParams{ChatID: ChatIDInt(1)}},
		{name: "both rich formats", params: SendRichMessageParams{ChatID: ChatIDInt(1), RichMessage: InputRichMessage{HTML: "<b>Hello</b>", Markdown: "**Hello**"}}},
		{name: "blank html", params: SendRichMessageParams{ChatID: ChatIDInt(1), RichMessage: InputRichMessage{HTML: "   "}}},
		{name: "negative thread", params: SendRichMessageParams{ChatID: ChatIDInt(1), MessageThreadID: -1, RichMessage: InputRichHTML("<b>Hello</b>")}},
		{name: "negative direct topic", params: SendRichMessageParams{ChatID: ChatIDInt(1), DirectMessagesTopicID: -1, RichMessage: InputRichHTML("<b>Hello</b>")}},
		{name: "invalid reply params", params: SendRichMessageParams{ChatID: ChatIDInt(1), RichMessage: InputRichHTML("<b>Hello</b>"), ReplyParameters: &telegram.ReplyParameters{}}},
		{name: "invalid reply markup", params: SendRichMessageParams{ChatID: ChatIDInt(1), RichMessage: InputRichHTML("<b>Hello</b>"), ReplyMarkup: markup}},
		{name: "invalid suggested post", params: SendRichMessageParams{ChatID: ChatIDInt(1), RichMessage: InputRichHTML("<b>Hello</b>"), SuggestedPostParameters: &telegram.SuggestedPostParameters{Price: &telegram.SuggestedPostPrice{Amount: 1}}}},
	}
	for _, tt := range sendTests {
		t.Run("send "+tt.name, func(t *testing.T) {
			message, err := bot.SendRichMessage(context.Background(), tt.params)
			if err == nil {
				t.Fatal("expected error")
			}
			if message != nil {
				t.Fatalf("expected nil message, got %+v", message)
			}
			assertNoToken(t, err, token)
		})
	}

	draftTests := []struct {
		name   string
		params SendRichMessageDraftParams
	}{
		{name: "missing chat", params: SendRichMessageDraftParams{DraftID: 1, RichMessage: InputRichMarkdown("**Hello**")}},
		{name: "string chat", params: SendRichMessageDraftParams{ChatID: ChatIDString("@channel"), DraftID: 1, RichMessage: InputRichMarkdown("**Hello**")}},
		{name: "negative thread", params: SendRichMessageDraftParams{ChatID: ChatIDInt(1), MessageThreadID: -1, DraftID: 1, RichMessage: InputRichMarkdown("**Hello**")}},
		{name: "zero draft", params: SendRichMessageDraftParams{ChatID: ChatIDInt(1), RichMessage: InputRichMarkdown("**Hello**")}},
		{name: "missing rich content", params: SendRichMessageDraftParams{ChatID: ChatIDInt(1), DraftID: 1}},
	}
	for _, tt := range draftTests {
		t.Run("draft "+tt.name, func(t *testing.T) {
			ok, err := bot.SendRichMessageDraft(context.Background(), tt.params)
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

func TestRichMessageMethodsReturnAPIError(t *testing.T) {
	const token = "123:secret"
	tests := []struct {
		name   string
		method string
		call   func(*Bot) error
	}{
		{name: "send", method: "sendRichMessage", call: func(bot *Bot) error {
			_, err := bot.SendRichMessage(context.Background(), SendRichMessageParams{ChatID: ChatIDInt(1), RichMessage: InputRichHTML("<b>Hello</b>")})
			return err
		}},
		{name: "draft", method: "sendRichMessageDraft", call: func(bot *Bot) error {
			_, err := bot.SendRichMessageDraft(context.Background(), SendRichMessageDraftParams{ChatID: ChatIDInt(1), DraftID: 1, RichMessage: InputRichMarkdown("**Hello**")})
			return err
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
			err := tt.call(bot)
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

func TestInputRichMessageContentValidatesForInlineResults(t *testing.T) {
	if err := validateInputMessageContent(InputRichMessageContent{RichMessage: InputRichMarkdown("**Hello**")}); err != nil {
		t.Fatalf("unexpected valid rich content error: %v", err)
	}
	if err := validateInputMessageContent(InputRichMessageContent{}); err == nil {
		t.Fatal("expected invalid rich content error")
	}
}
