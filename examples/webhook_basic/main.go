// Example webhook_basic runs a minimal Telegram webhook HTTP server.
//
// Required env: AIGRAM_BOT_TOKEN, AIGRAM_WEBHOOK_URL.
// Optional env: AIGRAM_WEBHOOK_SECRET, AIGRAM_LISTEN_ADDR, AIGRAM_BASE_URL, AIGRAM_FILE_BASE_URL.
// It serves /webhook and does not delete the webhook on shutdown. Stop with Ctrl+C or SIGTERM.
package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/xDilettante/ai-gram"
	"github.com/xDilettante/ai-gram/dispatch"
	"github.com/xDilettante/ai-gram/examples/internal/exampleutil"
	"github.com/xDilettante/ai-gram/telegram"
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

	dp, err := newDispatcher(b)
	if err != nil {
		return err
	}
	receiver, err := webhook.New(dp, webhook.Config{
		SecretToken: secret,
		OnError: func(ctx context.Context, update *telegram.Update, err error) {
			if update != nil {
				log.Printf("webhook handler error update_id=%d err=%v", update.UpdateID, err)
				return
			}
			log.Printf("webhook handler error err=%v", err)
		},
	})
	if err != nil {
		return err
	}

	server := &http.Server{Addr: listenAddr, Handler: webhookMux(receiver)}
	serverErr := make(chan error, 1)
	go func() {
		log.Printf("listening on %s", listenAddr)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErr <- err
			return
		}
		serverErr <- nil
	}()

	if _, err := b.SetWebhook(ctx, aigram.SetWebhookParams{
		URL:                webhookURL,
		SecretToken:        secret,
		DropPendingUpdates: true,
	}); err != nil {
		shutdownServer(server)
		return fmt.Errorf("set webhook: %w", err)
	}
	log.Println("webhook registered")

	select {
	case <-ctx.Done():
		shutdownServer(server)
		<-serverErr
		log.Println("webhook server stopped; webhook was not deleted automatically")
		return nil
	case err := <-serverErr:
		return err
	}
}

func webhookMux(receiver http.Handler) http.Handler {
	mux := http.NewServeMux()
	mux.Handle("/webhook", receiver)
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok\n"))
	})
	return mux
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
			Text:   "Webhook bot is online. Send a text message to receive a reply.",
		})
		return err
	}); err != nil {
		return nil, err
	}

	if err := dp.OnMessageFunc(func(ctx context.Context, update telegram.Update) error {
		message := update.EffectiveMessage()
		if message == nil || message.Text == "" {
			return nil
		}
		_, err := b.SendMessage(ctx, aigram.SendMessageParams{
			ChatID:          aigram.ChatIDInt(message.Chat.ID),
			Text:            "received via webhook",
			ReplyParameters: &aigram.ReplyParameters{MessageID: message.MessageID},
		})
		return err
	}); err != nil {
		return nil, err
	}

	return dp, nil
}

func shutdownServer(server *http.Server) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Printf("server shutdown error: %v", err)
	}
}
