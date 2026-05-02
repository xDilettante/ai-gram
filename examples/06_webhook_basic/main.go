// Example 06_webhook_basic runs a minimal Telegram webhook server.
package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	aigram "github.com/xDilettante/ai-gram"
	"github.com/xDilettante/ai-gram/examples/internal/exampleutil"
	"github.com/xDilettante/ai-gram/transport/webhook"
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
	webhookURL, err := exampleutil.RequiredEnv("AIGRAM_WEBHOOK_URL")
	if err != nil {
		return err
	}
	secret := exampleutil.OptionalEnv("AIGRAM_WEBHOOK_SECRET", "")
	listenAddr := exampleutil.OptionalEnv("AIGRAM_LISTEN_ADDR", ":8080")

	receiver, err := webhook.New(webhook.HandlerFunc(func(ctx context.Context, update aigram.Update) error {
		return handleUpdate(ctx, b, update)
	}), webhook.Config{SecretToken: secret})
	if err != nil {
		return err
	}

	server := &http.Server{Addr: listenAddr, Handler: mux(receiver)}
	serverErr := make(chan error, 1)
	go func() {
		log.Printf("listening on %s", listenAddr)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErr <- err
			return
		}
		serverErr <- nil
	}()

	if _, err := b.SetWebhook(ctx, aigram.SetWebhookParams{URL: webhookURL, SecretToken: secret}); err != nil {
		shutdown(server)
		return fmt.Errorf("set webhook: %w", err)
	}
	log.Println("webhook registered")

	select {
	case <-ctx.Done():
		shutdown(server)
		<-serverErr
		return nil
	case err := <-serverErr:
		return err
	}
}

func mux(receiver http.Handler) http.Handler {
	mux := http.NewServeMux()
	mux.Handle("/webhook", receiver)
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok\n"))
	})
	return mux
}

func handleUpdate(ctx context.Context, b *aigram.Bot, update aigram.Update) error {
	message := update.Message
	if message == nil || strings.TrimSpace(message.Text) == "" {
		return nil
	}
	text := "Received via webhook."
	if message.Text == "/start" {
		text = "Webhook bot is online. Send any text message."
	}
	_, err := b.SendMessage(ctx, aigram.SendMessageParams{
		ChatID:          aigram.ChatIDInt(message.Chat.ID),
		Text:            text,
		ReplyParameters: &aigram.ReplyParameters{MessageID: message.MessageID},
	})
	return err
}

func shutdown(server *http.Server) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Printf("server shutdown error: %v", err)
	}
}
