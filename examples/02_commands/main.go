// Example 02_commands shows simple command routing with switch/case.
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

	log.Println("commands bot started; press Ctrl+C to stop")
	return poll(ctx, b, handleCommand)
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

func handleCommand(ctx context.Context, b *aigram.Bot, update aigram.Update) error {
	message := update.Message
	if message == nil || strings.TrimSpace(message.Text) == "" {
		return nil
	}

	var text string
	switch message.Text {
	case "/start":
		text = "Welcome! Try /help or /about."
	case "/help":
		text = "Available commands: /start, /help, /about."
	case "/about":
		text = "This bot is a small ai-gram command routing example."
	default:
		text = "Unknown command. Try /help."
	}

	_, err := b.SendMessage(ctx, aigram.SendMessageParams{
		ChatID: aigram.ChatIDInt(message.Chat.ID),
		Text:   text,
	})
	return err
}
