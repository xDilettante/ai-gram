// Example webhook_server runs an HTTP webhook receiver and registers it with Telegram.
//
// Required env: AIGRAM_BOT_TOKEN, AIGRAM_WEBHOOK_URL.
// Optional env: AIGRAM_BASE_URL, AIGRAM_FILE_BASE_URL, AIGRAM_LISTEN_ADDR (default ":8080"), AIGRAM_WEBHOOK_SECRET.
// It serves /webhook and does not delete the webhook on shutdown. Stop with Ctrl+C or SIGTERM.
package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"ai-gram"
	"ai-gram/dispatch"
	"ai-gram/examples/internal/exampleutil"
	"ai-gram/telegram"
	"ai-gram/transport/webhook"
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
	listenAddr := exampleutil.OptionalEnv("AIGRAM_LISTEN_ADDR", ":8080")
	secret := exampleutil.OptionalEnv("AIGRAM_WEBHOOK_SECRET", "")

	dp, err := newDispatcher(b)
	if err != nil {
		return err
	}
	handler, err := webhook.New(dp, webhook.Config{
		SecretToken: secret,
		OnError: func(ctx context.Context, update *telegram.Update, err error) {
			if update != nil {
				log.Printf("webhook handler error update_id=%d: %v", update.UpdateID, err)
				return
			}
			log.Printf("webhook handler error: %v", err)
		},
	})
	if err != nil {
		return err
	}

	mux := http.NewServeMux()
	mux.Handle("/webhook", handler)
	server := &http.Server{Addr: listenAddr, Handler: mux}
	serverErr := make(chan error, 1)
	go func() {
		log.Printf("webhook server listening on %s", listenAddr)
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
		if err != nil {
			return err
		}
		return nil
	}
}

func newDispatcher(b *aigram.Bot) (*dispatch.Dispatcher, error) {
	dp := dispatch.New(dispatch.WithErrorHandler(func(ctx context.Context, update telegram.Update, err error) {
		log.Printf("handler error for update %d: %v", update.UpdateID, err)
	}))
	if err := dp.OnCommandFunc("start", func(ctx context.Context, update telegram.Update) error {
		message := update.EffectiveMessage()
		if message == nil {
			return nil
		}
		logSafeUpdate(update, "command")
		_, err := b.SendMessage(ctx, aigram.SendMessageParams{
			ChatID:      aigram.ChatIDInt(message.Chat.ID),
			Text:        "Webhook bot is running. Choose an action:",
			ReplyMarkup: demoKeyboard(),
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
		logSafeUpdate(update, "message")
		_, err := b.SendMessage(ctx, aigram.SendMessageParams{ChatID: aigram.ChatIDInt(message.Chat.ID), Text: "echo received"})
		return err
	}); err != nil {
		return nil, err
	}
	if err := dp.OnCallbackDataFunc("demo:edit", func(ctx context.Context, update telegram.Update) error {
		callback := update.CallbackQuery
		if callback == nil {
			return nil
		}
		logSafeUpdate(update, "callback_query")
		if _, err := b.AnswerCallbackQuery(ctx, aigram.AnswerCallbackQueryParams{CallbackQueryID: callback.ID, Text: "Editing message"}); err != nil {
			return err
		}
		if callback.Message == nil {
			_, err := b.AnswerCallbackQuery(ctx, aigram.AnswerCallbackQueryParams{CallbackQueryID: callback.ID, Text: "Inline message cannot be edited by this demo", ShowAlert: true})
			return err
		}

		removeKeyboard := removeKeyboardMarkup()
		_, err := b.EditMessageText(ctx, aigram.EditMessageTextParams{
			Target:      aigram.EditTargetChat(aigram.ChatIDInt(callback.Message.Chat.ID), callback.Message.MessageID),
			Text:        "Message edited by ai-gram",
			ReplyMarkup: &removeKeyboard,
		})
		return err
	}); err != nil {
		return nil, err
	}
	if err := dp.OnCallbackDataFunc("demo:remove_keyboard", func(ctx context.Context, update telegram.Update) error {
		callback := update.CallbackQuery
		if callback == nil {
			return nil
		}
		logSafeUpdate(update, "callback_query")
		if _, err := b.AnswerCallbackQuery(ctx, aigram.AnswerCallbackQueryParams{CallbackQueryID: callback.ID, Text: "Removing keyboard"}); err != nil {
			return err
		}
		if callback.Message == nil {
			return nil
		}
		_, err := b.EditMessageReplyMarkup(ctx, aigram.EditMessageReplyMarkupParams{
			Target: aigram.EditTargetChat(aigram.ChatIDInt(callback.Message.Chat.ID), callback.Message.MessageID),
		})
		return err
	}); err != nil {
		return nil, err
	}
	return dp, nil
}

func demoKeyboard() aigram.InlineKeyboardMarkup {
	return aigram.NewInlineKeyboard([]aigram.InlineKeyboardButton{
		aigram.InlineButtonCallback("Edit message", "demo:edit"),
		aigram.InlineButtonCallback("Remove keyboard", "demo:remove_keyboard"),
	})
}

func removeKeyboardMarkup() aigram.InlineKeyboardMarkup {
	return aigram.NewInlineKeyboard([]aigram.InlineKeyboardButton{
		aigram.InlineButtonCallback("Remove keyboard", "demo:remove_keyboard"),
	})
}

func logSafeUpdate(update telegram.Update, matched string) {
	message := update.EffectiveMessage()
	chat := update.EffectiveChat()
	user := update.EffectiveUser()
	chatID := int64(0)
	userID := int64(0)
	if chat != nil {
		chatID = chat.ID
	}
	if user != nil {
		userID = user.ID
	}

	updateType := "unknown"
	command := ""
	hasText := false
	hasMedia := false
	if update.CallbackQuery != nil {
		updateType = "callback_query"
	} else if update.Message != nil {
		updateType = "message"
	} else if update.EditedMessage != nil {
		updateType = "edited_message"
	}
	if message != nil {
		hasText = message.Text != ""
		hasMedia = message.HasMedia()
		command = message.Command()
	}

	callbackData := safeCallbackData(update)
	if command != "" {
		log.Printf("webhook update_id=%d update_type=%s matched=%s chat_id=%d from_user_id=%d command=%s has_text=%t has_media=%t callback_data=%s", update.UpdateID, updateType, matched, chatID, userID, command, hasText, hasMedia, callbackData)
		return
	}
	log.Printf("webhook update_id=%d update_type=%s matched=%s chat_id=%d from_user_id=%d has_text=%t has_media=%t callback_data=%s", update.UpdateID, updateType, matched, chatID, userID, hasText, hasMedia, callbackData)
}

func safeCallbackData(update telegram.Update) string {
	if update.CallbackQuery == nil {
		return ""
	}
	switch update.CallbackQuery.Data {
	case "demo:edit", "demo:remove_keyboard":
		return update.CallbackQuery.Data
	default:
		return "<redacted>"
	}
}

func shutdownServer(server *http.Server) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Printf("webhook server shutdown error: %v", err)
	}
}
