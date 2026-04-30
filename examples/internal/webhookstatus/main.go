// Package main prints token-safe webhook state for smoke diagnostics.
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"ai-gram/examples/internal/exampleutil"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	b, err := exampleutil.NewBotFromEnv()
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	info, err := b.GetWebhookInfo(ctx)
	if err != nil {
		return fmt.Errorf("getWebhookInfo: %w", err)
	}
	fmt.Printf("webhook_configured=%t pending_update_count=%d\n", info.URL != "", info.PendingUpdateCount)
	return nil
}
