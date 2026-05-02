// Example echo_longpoll runs a small long polling bot with commands, echo, and inline callbacks.
//
// Required env: AIGRAM_BOT_TOKEN.
// Optional env: AIGRAM_BASE_URL, AIGRAM_FILE_BASE_URL.
// It deletes any configured webhook before polling. Stop with Ctrl+C or SIGTERM.
package main

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/xDilettante/ai-gram"
	"github.com/xDilettante/ai-gram/dispatch"
	"github.com/xDilettante/ai-gram/examples/internal/exampleutil"
	"github.com/xDilettante/ai-gram/telegram"
	"github.com/xDilettante/ai-gram/transport/longpoll"
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
		return fmt.Errorf("delete webhook before long polling: %w", err)
	}

	dp, err := newDispatcher(b)
	if err != nil {
		return err
	}

	runner, err := longpoll.New(b, dp, longpoll.Config{Timeout: 30})
	if err != nil {
		return err
	}

	log.Println("echo long polling started")
	if err := runner.Run(ctx); err != nil && err != context.Canceled {
		return err
	}
	log.Println("echo long polling stopped")
	return nil
}

func newDispatcher(b *aigram.Bot) (*dispatch.Dispatcher, error) {
	dp := dispatch.New(dispatch.WithErrorHandler(func(ctx context.Context, update telegram.Update, err error) {
		log.Printf("handler error update_id=%d err=%v", update.UpdateID, err)
	}))

	if err := dp.OnCommandFunc("start", func(ctx context.Context, update telegram.Update) error {
		message := update.EffectiveMessage()
		if message == nil {
			return nil
		}
		_, err := b.SendMessage(ctx, aigram.SendMessageParams{
			ChatID: aigram.ChatIDInt(message.Chat.ID),
			Text:   "Hi! Send any text and I will echo it. Use /help for commands.",
			ReplyMarkup: telegram.NewInlineKeyboard([]telegram.InlineKeyboardButton{
				telegram.InlineButtonCallback("Show time", "demo:time"),
				telegram.InlineButtonCallback("About", "demo:about"),
			}),
		})
		return err
	}); err != nil {
		return nil, err
	}

	if err := dp.OnCommandFunc("help", func(ctx context.Context, update telegram.Update) error {
		message := update.EffectiveMessage()
		if message == nil {
			return nil
		}
		_, err := b.SendMessage(ctx, aigram.SendMessageParams{
			ChatID: aigram.ChatIDInt(message.Chat.ID),
			Text:   "Commands: /start, /help. Send text to receive an echo response.",
		})
		return err
	}); err != nil {
		return nil, err
	}

	if err := dp.OnCallbackDataFunc("demo:time", func(ctx context.Context, update telegram.Update) error {
		callback := update.CallbackQuery
		if callback == nil {
			return nil
		}
		if _, err := b.AnswerCallbackQuery(ctx, aigram.AnswerCallbackQueryParams{
			CallbackQueryID: callback.ID,
			Text:            "Current UTC time sent to the chat.",
		}); err != nil {
			return err
		}
		if callback.Message == nil {
			return nil
		}
		_, err := b.SendMessage(ctx, aigram.SendMessageParams{
			ChatID: aigram.ChatIDInt(callback.Message.Chat.ID),
			Text:   "UTC time: " + time.Now().UTC().Format(time.RFC3339),
		})
		return err
	}); err != nil {
		return nil, err
	}

	if err := dp.OnCallbackDataFunc("demo:about", func(ctx context.Context, update telegram.Update) error {
		callback := update.CallbackQuery
		if callback == nil {
			return nil
		}
		if _, err := b.AnswerCallbackQuery(ctx, aigram.AnswerCallbackQueryParams{
			CallbackQueryID: callback.ID,
			Text:            "Editing the original message.",
		}); err != nil {
			return err
		}
		if callback.Message == nil {
			return nil
		}
		_, err := b.EditMessageText(ctx, aigram.EditMessageTextParams{
			Target: aigram.EditTargetChat(aigram.ChatIDInt(callback.Message.Chat.ID), callback.Message.MessageID),
			Text:   "ai-gram routes updates with typed handlers, predicates, and middleware.",
		})
		return err
	}); err != nil {
		return nil, err
	}

	if err := dp.OnMessageFunc(func(ctx context.Context, update telegram.Update) error {
		message := update.EffectiveMessage()
		if message == nil || strings.TrimSpace(message.Text) == "" || strings.HasPrefix(message.Text, "/") {
			return nil
		}
		_, err := b.SendMessage(ctx, aigram.SendMessageParams{
			ChatID:          aigram.ChatIDInt(message.Chat.ID),
			MessageThreadID: message.MessageThreadID,
			Text:            "echo: " + message.Text,
			ReplyParameters: &aigram.ReplyParameters{MessageID: message.MessageID},
		})
		return err
	}); err != nil {
		return nil, err
	}

	return dp, nil
}
