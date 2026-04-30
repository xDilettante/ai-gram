// Example local_api_server checks connectivity to a local Telegram Bot API server.
//
// Required env: AIGRAM_BOT_TOKEN, AIGRAM_BASE_URL (for example "http://127.0.0.1:8081").
// Optional env: AIGRAM_FILE_BASE_URL.
// It calls getMe and getWebhookInfo once and exits.
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/xDilettante/ai-gram/examples/internal/exampleutil"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	if strings.TrimSpace(os.Getenv("AIGRAM_BASE_URL")) == "" {
		return fmt.Errorf("AIGRAM_BASE_URL is required for local_api_server")
	}
	b, err := exampleutil.NewBotFromEnv()
	if err != nil {
		return err
	}

	ctx := context.Background()
	me, err := b.GetMe(ctx)
	if err != nil {
		return fmt.Errorf("getMe: %w", err)
	}
	fmt.Println("bot username:", me.Username)

	info, err := b.GetWebhookInfo(ctx)
	if err != nil {
		return fmt.Errorf("getWebhookInfo: %w", err)
	}
	fmt.Printf("webhook url=%q pending_update_count=%d\n", info.URL, info.PendingUpdateCount)
	return nil
}
