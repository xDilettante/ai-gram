// Example inline_longpoll demonstrates inline keyboards and callback answers.
//
// Required env: AIGRAM_BOT_TOKEN.
// Optional env: AIGRAM_BASE_URL, AIGRAM_FILE_BASE_URL.
// It deletes any configured webhook before polling. Send /start to the bot. Stop with Ctrl+C or SIGTERM.
package main

import (
	"context"
	"fmt"
	"log"

	"ai-gram"
	"ai-gram/dispatch"
	"ai-gram/examples/internal/exampleutil"
	"ai-gram/telegram"
	"ai-gram/transport/longpoll"
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

	dp := dispatch.New(dispatch.WithErrorHandler(func(ctx context.Context, update telegram.Update, err error) {
		log.Printf("handler error for update %d: %v", update.UpdateID, err)
	}))
	if err := dp.OnCommandFunc("start", func(ctx context.Context, update telegram.Update) error {
		message := update.EffectiveMessage()
		if message == nil {
			return nil
		}
		keyboard := aigram.NewInlineKeyboard([]aigram.InlineKeyboardButton{
			aigram.InlineButtonCallback("Да", "demo:yes"),
			aigram.InlineButtonCallback("Нет", "demo:no"),
		})
		_, err := b.SendMessage(ctx, aigram.SendMessageParams{
			ChatID:      aigram.ChatIDInt(message.Chat.ID),
			Text:        "Выберите вариант",
			ReplyMarkup: keyboard,
		})
		return err
	}); err != nil {
		return err
	}
	if err := dp.OnCallbackDataFunc("demo:yes", func(ctx context.Context, update telegram.Update) error {
		callback := update.CallbackQuery
		if callback == nil || callback.Message == nil {
			return nil
		}
		if _, err := b.AnswerCallbackQuery(ctx, aigram.AnswerCallbackQueryParams{
			CallbackQueryID: callback.ID,
			Text:            "Вы выбрали Да",
		}); err != nil {
			return err
		}
		_, err := b.SendMessage(ctx, aigram.SendMessageParams{
			ChatID: aigram.ChatIDInt(callback.Message.Chat.ID),
			Text:   "Да подтверждено",
		})
		return err
	}); err != nil {
		return err
	}
	if err := dp.OnCallbackDataFunc("demo:no", func(ctx context.Context, update telegram.Update) error {
		if update.CallbackQuery == nil {
			return nil
		}
		_, err := b.AnswerCallbackQuery(ctx, aigram.AnswerCallbackQueryParams{
			CallbackQueryID: update.CallbackQuery.ID,
			Text:            "Вы выбрали Нет",
			ShowAlert:       true,
		})
		return err
	}); err != nil {
		return err
	}

	runner, err := longpoll.New(b, dp, longpoll.Config{Timeout: 30})
	if err != nil {
		return err
	}
	log.Println("inline long polling started")
	if err := runner.Run(ctx); err != nil && err != context.Canceled {
		return err
	}
	log.Println("inline long polling stopped")
	return nil
}
