// Example 01_echo_bot runs the smallest useful long polling bot.
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

	// Long polling and webhooks are mutually exclusive for one bot.
	if _, err := b.DeleteWebhook(ctx, aigram.DeleteWebhookParams{DropPendingUpdates: true}); err != nil {
		return fmt.Errorf("delete webhook before polling: %w", err)
	}

	log.Println("echo bot started; press Ctrl+C to stop")
	return poll(ctx, b, handleUpdate)
}

func poll(ctx context.Context, b *aigram.Bot, handle func(context.Context, *aigram.Bot, aigram.Update) error) error {
	var offset int64
	for {
		updates, err := b.GetUpdates(ctx, aigram.GetUpdatesParams{Offset: offset, Timeout: 30, AllowedUpdates: []string{"message"}})
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
	message := update.Message
	if message == nil || strings.TrimSpace(message.Text) == "" {
		return nil
	}

	if message.Text == "/start" {
		_, err := b.SendMessage(ctx, aigram.SendMessageParams{
			ChatID: aigram.ChatIDInt(message.Chat.ID),
			Text:   "Hi! Send any text message and I will echo it back.",
		})
		return err
	}

	_, err := b.SendMessage(ctx, aigram.SendMessageParams{
		ChatID: aigram.ChatIDInt(message.Chat.ID),
		Text:   message.Text,
	})
	return err
}
