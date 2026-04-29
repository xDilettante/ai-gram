// Package main sends a short Telegram notification for manual smoke scripts.
//
// Required env:
//   - AIGRAM_BOT_TOKEN: notification bot token;
//   - AIGRAM_CHAT_ID: target chat ID or username;
//   - AIGRAM_NOTIFY_TEXT: notification text.
//
// Optional env:
//   - AIGRAM_BASE_URL / AIGRAM_FILE_BASE_URL for local Bot API server.
package main

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"ai-gram"
	"ai-gram/examples/internal/exampleutil"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	chatRaw, err := exampleutil.RequiredEnv("AIGRAM_CHAT_ID")
	if err != nil {
		return err
	}
	chatID, err := exampleutil.ParseChatID(chatRaw)
	if err != nil {
		return err
	}

	text, err := exampleutil.RequiredEnv("AIGRAM_NOTIFY_TEXT")
	if err != nil {
		return err
	}
	if strings.TrimSpace(text) == "" {
		return fmt.Errorf("AIGRAM_NOTIFY_TEXT is required")
	}

	b, err := exampleutil.NewBotFromEnv()
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	_, err = b.SendMessage(ctx, aigram.SendMessageParams{
		ChatID: chatID,
		Text:   text,
	})
	if err != nil {
		return fmt.Errorf("send Telegram notification: %w", err)
	}

	log.Println("Telegram notification sent")
	return nil
}
