package bot

import (
	"context"
	"testing"

	"github.com/xDilettante/ai-gram/telegram"
)

func TestSavePreparedInlineMessageSendsPayloadAndDecodesResult(t *testing.T) {
	testBotProfileObjectSuccess(t, "savePreparedInlineMessage", `{"id":"prepared-inline-id","expiration_date":1777606551}`, func(bot *Bot) (any, error) {
		article := InlineArticle("article-1", "Article", InputText("hello"))
		article.ReplyMarkup = &telegram.InlineKeyboardMarkup{InlineKeyboard: [][]telegram.InlineKeyboardButton{{
			telegram.InlineButtonCopyText("Copy", "copy me"),
		}}}
		return bot.SavePreparedInlineMessage(context.Background(), SavePreparedInlineMessageParams{
			UserID:            123,
			Result:            article,
			AllowUserChats:    true,
			AllowBotChats:     true,
			AllowGroupChats:   true,
			AllowChannelChats: true,
		})
	}, func(t *testing.T, payload map[string]any) {
		if payload["user_id"] != float64(123) || payload["allow_user_chats"] != true || payload["allow_bot_chats"] != true || payload["allow_group_chats"] != true || payload["allow_channel_chats"] != true {
			t.Fatalf("unexpected payload: %#v", payload)
		}
		result, ok := payload["result"].(map[string]any)
		if !ok {
			t.Fatalf("missing result payload: %#v", payload)
		}
		if result["type"] != "article" || result["id"] != "article-1" || result["title"] != "Article" {
			t.Fatalf("unexpected result payload: %#v", result)
		}
		content := result["input_message_content"].(map[string]any)
		if content["message_text"] != "hello" {
			t.Fatalf("unexpected input_message_content: %#v", content)
		}
		replyMarkup := result["reply_markup"].(map[string]any)
		button := replyMarkup["inline_keyboard"].([]any)[0].([]any)[0].(map[string]any)
		copyText := button["copy_text"].(map[string]any)
		if copyText["text"] != "copy me" {
			t.Fatalf("unexpected copy_text button: %#v", button)
		}
	}, func(t *testing.T, result any) {
		message := result.(*telegram.PreparedInlineMessage)
		if message.ID != "prepared-inline-id" || message.ExpirationDate != 1777606551 {
			t.Fatalf("unexpected result: %+v", message)
		}
	})
}

func TestSavePreparedInlineMessageValidation(t *testing.T) {
	const token = "123:secret"
	bot := newTestBot(t, token, "https://example.test", nil)
	tests := []struct {
		name   string
		params SavePreparedInlineMessageParams
	}{
		{name: "invalid user", params: SavePreparedInlineMessageParams{UserID: 0, Result: InlineArticle("article-1", "Article", InputText("hello"))}},
		{name: "nil result", params: SavePreparedInlineMessageParams{UserID: 1}},
		{name: "invalid result", params: SavePreparedInlineMessageParams{UserID: 1, Result: InlineQueryResultArticle{ID: "article-1", InputMessageContent: InputText("hello")}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := bot.SavePreparedInlineMessage(context.Background(), tt.params)
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

func TestSavePreparedInlineMessageErrors(t *testing.T) {
	call := func(bot *Bot) (any, error) {
		return bot.SavePreparedInlineMessage(context.Background(), validSavePreparedInlineMessageParams())
	}
	callWithContext := func(bot *Bot, ctx context.Context) (any, error) {
		return bot.SavePreparedInlineMessage(ctx, validSavePreparedInlineMessageParams())
	}
	testBotProfileObjectMethodErrorCases(t, "savePreparedInlineMessage", call, callWithContext)
}

func validSavePreparedInlineMessageParams() SavePreparedInlineMessageParams {
	return SavePreparedInlineMessageParams{
		UserID: 123,
		Result: InlineArticle("article-1", "Article", InputText("hello")),
	}
}
