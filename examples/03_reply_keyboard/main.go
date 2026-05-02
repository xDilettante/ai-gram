// Example 03_reply_keyboard shows a regular Telegram reply keyboard.
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

	log.Println("reply keyboard bot started; press Ctrl+C to stop")
	return poll(ctx, b, handleMessage)
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

func handleMessage(ctx context.Context, b *aigram.Bot, update aigram.Update) error {
	message := update.Message
	if message == nil || strings.TrimSpace(message.Text) == "" {
		return nil
	}

	params := aigram.SendMessageParams{ChatID: aigram.ChatIDInt(message.Chat.ID)}
	switch message.Text {
	case "/start":
		params.Text = "Choose a button below."
		params.ReplyMarkup = aigram.NewReplyKeyboard([]aigram.KeyboardButton{
			aigram.KeyboardButtonText("Help"),
			aigram.KeyboardButtonText("About"),
		}, []aigram.KeyboardButton{
			aigram.KeyboardButtonText("Remove keyboard"),
		})
	case "Help":
		params.Text = "Use /start to show the keyboard again."
	case "About":
		params.Text = "Reply keyboards send normal text messages when users press buttons."
	case "Remove keyboard":
		params.Text = "Keyboard removed. Send /start to show it again."
		params.ReplyMarkup = aigram.RemoveKeyboard(false)
	default:
		params.Text = "Send /start to open the reply keyboard."
	}

	_, err := b.SendMessage(ctx, params)
	return err
}
