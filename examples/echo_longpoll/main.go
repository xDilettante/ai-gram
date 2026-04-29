// Example echo_longpoll runs a simple echo bot with long polling.
//
// Required env: AIGRAM_BOT_TOKEN.
// Optional env: AIGRAM_BASE_URL, AIGRAM_FILE_BASE_URL.
// It deletes any configured webhook before polling. Stop with Ctrl+C or SIGTERM.
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
	if err := dp.OnMessageFunc(func(ctx context.Context, update telegram.Update) error {
		message := update.EffectiveMessage()
		if message == nil || message.Text == "" {
			return nil
		}
		_, err := b.SendMessage(ctx, aigram.SendMessageParams{
			ChatID: aigram.ChatIDInt(message.Chat.ID),
			Text:   "echo: " + message.Text,
		})
		return err
	}); err != nil {
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
