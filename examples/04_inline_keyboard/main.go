// Example 04_inline_keyboard shows inline keyboard callbacks.
package main

import (
	"context"
	"fmt"
	"log"
	"strings"

	aigram "github.com/xDilettante/ai-gram"
	"github.com/xDilettante/ai-gram/examples/internal/exampleutil"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	ctx, stop := exampleutil.SignalContext()
	defer stop()

	b, err := exampleutil.NewBotFromEnv()
	if err != nil {
		return err
	}
	if _, err := b.DeleteWebhook(ctx, aigram.DeleteWebhookParams{DropPendingUpdates: true}); err != nil {
		return fmt.Errorf("delete webhook before polling: %w", err)
	}

	log.Println("inline keyboard bot started; press Ctrl+C to stop")
	return poll(ctx, b, handleUpdate)
}

func poll(ctx context.Context, b *aigram.Bot, handle func(context.Context, *aigram.Bot, aigram.Update) error) error {
	var offset int64
	for {
		updates, err := b.GetUpdates(ctx, aigram.GetUpdatesParams{Offset: offset, Timeout: 30, AllowedUpdates: []string{"message", "callback_query"}})
		if err != nil {
			if ctx.Err() != nil {
				return nil
			}
			return err
		}
		for _, update := range updates {
			offset = update.UpdateID + 1
			if err := handle(ctx, b, update); err != nil {
				log.Printf("handle update %d: %v", update.UpdateID, err)
			}
		}
	}
}

func handleUpdate(ctx context.Context, b *aigram.Bot, update aigram.Update) error {
	if update.CallbackQuery != nil {
		return handleCallback(ctx, b, update.CallbackQuery)
	}
	message := update.Message
	if message == nil || !strings.HasPrefix(message.Text, "/start") {
		return nil
	}
	_, err := b.SendMessage(ctx, aigram.SendMessageParams{
		ChatID: aigram.ChatIDInt(message.Chat.ID),
		Text:   "Press an inline button:",
		ReplyMarkup: aigram.NewInlineKeyboard([]aigram.InlineKeyboardButton{
			aigram.InlineButtonCallback("Edit this message", "edit"),
			aigram.InlineButtonCallback("Remove buttons", "remove"),
		}),
	})
	return err
}

func handleCallback(ctx context.Context, b *aigram.Bot, callback *aigram.CallbackQuery) error {
	if _, err := b.AnswerCallbackQuery(ctx, aigram.AnswerCallbackQueryParams{
		CallbackQueryID: callback.ID,
		Text:            "Callback received",
	}); err != nil {
		return err
	}
	if callback.Message == nil {
		return nil
	}

	target := aigram.EditTargetChat(aigram.ChatIDInt(callback.Message.Chat.ID), callback.Message.MessageID)
	switch callback.Data {
	case "edit":
		_, err := b.EditMessageText(ctx, aigram.EditMessageTextParams{
			Target: target,
			Text:   "The message text was edited after a callback.",
		})
		return err
	case "remove":
		_, err := b.EditMessageReplyMarkup(ctx, aigram.EditMessageReplyMarkupParams{Target: target})
		return err
	default:
		return nil
	}
}
