// Example webhook_server runs an HTTP webhook receiver and registers it with Telegram.
//
// Required env: AIGRAM_BOT_TOKEN, AIGRAM_WEBHOOK_URL.
// Optional env: AIGRAM_BASE_URL, AIGRAM_FILE_BASE_URL, AIGRAM_LISTEN_ADDR (default ":8080"),
// AIGRAM_WEBHOOK_SECRET, AIGRAM_FILE_ID or AIGRAM_MEDIA_PATH for the caption demo.
// It serves /webhook and does not delete the webhook on shutdown. Stop with Ctrl+C or SIGTERM.
package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
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
		if err != nil {
			return err
		}
		log.Printf("webhook action=send_message ok=true update_id=%d chat_id=%d", update.UpdateID, message.Chat.ID)
		return nil
	}); err != nil {
		return nil, err
	}
	if err := dp.OnMessageFunc(func(ctx context.Context, update telegram.Update) error {
		message := update.EffectiveMessage()
		if message == nil || message.Text == "" {
			return nil
		}
		logSafeUpdate(update, "message")
		ok, err := b.SendChatAction(ctx, aigram.SendChatActionParams{
			ChatID:          aigram.ChatIDInt(message.Chat.ID),
			MessageThreadID: message.MessageThreadID,
			Action:          aigram.ChatActionTyping,
		})
		if err != nil {
			return err
		}
		if !ok {
			return errors.New("send_chat_action returned false")
		}
		log.Printf("webhook action=send_chat_action ok=true update_id=%d chat_id=%d chat_action=%s", update.UpdateID, message.Chat.ID, aigram.ChatActionTyping)

		_, err = b.SendMessage(ctx, aigram.SendMessageParams{
			ChatID:          aigram.ChatIDInt(message.Chat.ID),
			MessageThreadID: message.MessageThreadID,
			Text:            "echo received",
			ReplyParameters: &aigram.ReplyParameters{MessageID: message.MessageID},
		})
		if err != nil {
			return err
		}
		log.Printf("webhook action=send_message ok=true update_id=%d chat_id=%d reply_to_message_id=%d", update.UpdateID, message.Chat.ID, message.MessageID)
		return nil
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
		log.Printf("webhook action=answer_callback_query ok=true update_id=%d callback_data=%s", update.UpdateID, safeCallbackData(update))
		if callback.Message == nil {
			return nil
		}

		removeKeyboard := removeKeyboardMarkup()
		_, err := b.EditMessageText(ctx, aigram.EditMessageTextParams{
			Target:      aigram.EditTargetChat(aigram.ChatIDInt(callback.Message.Chat.ID), callback.Message.MessageID),
			Text:        "Message edited by ai-gram",
			ReplyMarkup: &removeKeyboard,
		})
		if err != nil {
			return err
		}
		log.Printf("webhook action=edit_message_text ok=true update_id=%d chat_id=%d message_id=%d", update.UpdateID, callback.Message.Chat.ID, callback.Message.MessageID)
		return nil
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
		log.Printf("webhook action=answer_callback_query ok=true update_id=%d callback_data=%s", update.UpdateID, safeCallbackData(update))
		if callback.Message == nil {
			return nil
		}
		_, err := b.EditMessageReplyMarkup(ctx, aigram.EditMessageReplyMarkupParams{
			Target: aigram.EditTargetChat(aigram.ChatIDInt(callback.Message.Chat.ID), callback.Message.MessageID),
		})
		if err != nil {
			return err
		}
		log.Printf("webhook action=edit_message_reply_markup ok=true update_id=%d chat_id=%d message_id=%d", update.UpdateID, callback.Message.Chat.ID, callback.Message.MessageID)
		return nil
	}); err != nil {
		return nil, err
	}
	if err := dp.OnCallbackDataFunc("demo:caption", func(ctx context.Context, update telegram.Update) error {
		callback := update.CallbackQuery
		if callback == nil {
			return nil
		}
		logSafeUpdate(update, "callback_query")
		if _, err := b.AnswerCallbackQuery(ctx, aigram.AnswerCallbackQueryParams{CallbackQueryID: callback.ID, Text: "Sending caption demo"}); err != nil {
			return err
		}
		log.Printf("webhook action=answer_callback_query ok=true update_id=%d callback_data=%s", update.UpdateID, safeCallbackData(update))
		if callback.Message == nil {
			return nil
		}

		sent, source, err := sendCaptionDemo(ctx, b, callback.Message.Chat.ID)
		if err != nil {
			return err
		}
		if sent == nil {
			return errors.New("send_media_caption_demo returned nil message")
		}
		log.Printf("webhook action=send_media_caption_demo ok=true source=%s update_id=%d chat_id=%d message_id=%d", source, update.UpdateID, sent.Chat.ID, sent.MessageID)
		return nil
	}); err != nil {
		return nil, err
	}
	if err := dp.OnCallbackDataFunc("demo:edit_caption", func(ctx context.Context, update telegram.Update) error {
		callback := update.CallbackQuery
		if callback == nil {
			return nil
		}
		logSafeUpdate(update, "callback_query")
		if _, err := b.AnswerCallbackQuery(ctx, aigram.AnswerCallbackQueryParams{CallbackQueryID: callback.ID, Text: "Editing caption"}); err != nil {
			return err
		}
		log.Printf("webhook action=answer_callback_query ok=true update_id=%d callback_data=%s", update.UpdateID, safeCallbackData(update))
		if callback.Message == nil {
			return nil
		}

		deleteKeyboard := deleteMediaKeyboard()
		result, err := b.EditMessageCaption(ctx, aigram.EditMessageCaptionParams{
			Target:      aigram.EditTargetChat(aigram.ChatIDInt(callback.Message.Chat.ID), callback.Message.MessageID),
			Caption:     "Caption edited by ai-gram",
			ReplyMarkup: &deleteKeyboard,
		})
		if err != nil {
			return err
		}
		if !result.IsOK() {
			return errors.New("edit_message_caption returned non-ok result")
		}
		log.Printf("webhook action=edit_message_caption ok=true update_id=%d chat_id=%d message_id=%d", update.UpdateID, callback.Message.Chat.ID, callback.Message.MessageID)
		return nil
	}); err != nil {
		return nil, err
	}
	if err := dp.OnCallbackDataFunc("demo:copy", func(ctx context.Context, update telegram.Update) error {
		callback := update.CallbackQuery
		if callback == nil {
			return nil
		}
		logSafeUpdate(update, "callback_query")
		if _, err := b.AnswerCallbackQuery(ctx, aigram.AnswerCallbackQueryParams{CallbackQueryID: callback.ID, Text: "Copying message"}); err != nil {
			return err
		}
		log.Printf("webhook action=answer_callback_query ok=true update_id=%d callback_data=%s", update.UpdateID, safeCallbackData(update))
		if callback.Message == nil {
			return nil
		}

		copied, err := b.CopyMessage(ctx, aigram.CopyMessageParams{
			ChatID:     aigram.ChatIDInt(callback.Message.Chat.ID),
			FromChatID: aigram.ChatIDInt(callback.Message.Chat.ID),
			MessageID:  callback.Message.MessageID,
		})
		if err != nil {
			return err
		}
		log.Printf("webhook action=copy_message ok=true update_id=%d chat_id=%d message_id=%d copied_message_id=%d", update.UpdateID, callback.Message.Chat.ID, callback.Message.MessageID, copied.MessageID)
		return nil
	}); err != nil {
		return nil, err
	}
	if err := dp.OnCallbackDataFunc("demo:forward", func(ctx context.Context, update telegram.Update) error {
		callback := update.CallbackQuery
		if callback == nil {
			return nil
		}
		logSafeUpdate(update, "callback_query")
		if _, err := b.AnswerCallbackQuery(ctx, aigram.AnswerCallbackQueryParams{CallbackQueryID: callback.ID, Text: "Forwarding message"}); err != nil {
			return err
		}
		log.Printf("webhook action=answer_callback_query ok=true update_id=%d callback_data=%s", update.UpdateID, safeCallbackData(update))
		if callback.Message == nil {
			return nil
		}

		forwarded, err := b.ForwardMessage(ctx, aigram.ForwardMessageParams{
			ChatID:     aigram.ChatIDInt(callback.Message.Chat.ID),
			FromChatID: aigram.ChatIDInt(callback.Message.Chat.ID),
			MessageID:  callback.Message.MessageID,
		})
		if err != nil {
			return err
		}
		log.Printf("webhook action=forward_message ok=true update_id=%d chat_id=%d message_id=%d forwarded_message_id=%d", update.UpdateID, callback.Message.Chat.ID, callback.Message.MessageID, forwarded.MessageID)
		return nil
	}); err != nil {
		return nil, err
	}
	if err := dp.OnCallbackDataFunc("demo:delete_media", func(ctx context.Context, update telegram.Update) error {
		return deleteCallbackMessage(ctx, b, update, "Deleting message")
	}); err != nil {
		return nil, err
	}
	if err := dp.OnCallbackDataFunc("demo:delete", func(ctx context.Context, update telegram.Update) error {
		return deleteCallbackMessage(ctx, b, update, "Deleting this message")
	}); err != nil {
		return nil, err
	}
	return dp, nil
}

func demoKeyboard() aigram.InlineKeyboardMarkup {
	return aigram.NewInlineKeyboard(
		[]aigram.InlineKeyboardButton{
			aigram.InlineButtonCallback("Edit message", "demo:edit"),
			aigram.InlineButtonCallback("Remove keyboard", "demo:remove_keyboard"),
		},
		[]aigram.InlineKeyboardButton{
			aigram.InlineButtonCallback("Caption demo", "demo:caption"),
			aigram.InlineButtonCallback("Delete this message", "demo:delete"),
		},
		[]aigram.InlineKeyboardButton{
			aigram.InlineButtonCallback("Copy this message", "demo:copy"),
			aigram.InlineButtonCallback("Forward this message", "demo:forward"),
		},
	)
}

func removeKeyboardMarkup() aigram.InlineKeyboardMarkup {
	return aigram.NewInlineKeyboard([]aigram.InlineKeyboardButton{
		aigram.InlineButtonCallback("Remove keyboard", "demo:remove_keyboard"),
	})
}

func captionDemoKeyboard() aigram.InlineKeyboardMarkup {
	return aigram.NewInlineKeyboard([]aigram.InlineKeyboardButton{
		aigram.InlineButtonCallback("Edit caption", "demo:edit_caption"),
		aigram.InlineButtonCallback("Delete media message", "demo:delete_media"),
	})
}

func deleteMediaKeyboard() aigram.InlineKeyboardMarkup {
	return aigram.NewInlineKeyboard([]aigram.InlineKeyboardButton{
		aigram.InlineButtonCallback("Delete media message", "demo:delete_media"),
	})
}

func sendCaptionDemo(ctx context.Context, b *aigram.Bot, chatID int64) (*telegram.Message, string, error) {
	fileID := exampleutil.OptionalEnv("AIGRAM_FILE_ID", "")
	mediaPath := exampleutil.OptionalEnv("AIGRAM_MEDIA_PATH", "")

	markup := captionDemoKeyboard()
	if fileID != "" {
		message, err := b.SendDocument(ctx, aigram.SendDocumentParams{
			ChatID:      aigram.ChatIDInt(chatID),
			Document:    aigram.FileID(fileID),
			Caption:     "Original caption from ai-gram",
			ReplyMarkup: markup,
		})
		return message, "file_id", err
	}

	if mediaPath != "" {
		file, err := os.Open(mediaPath)
		if err != nil {
			return nil, "", fmt.Errorf("open AIGRAM_MEDIA_PATH: %w", err)
		}
		defer file.Close()

		name := filepath.Base(mediaPath)
		if name == "." || name == string(filepath.Separator) || name == "" {
			name = "upload.bin"
		}
		message, err := b.SendDocument(ctx, aigram.SendDocumentParams{
			ChatID: aigram.ChatIDInt(chatID),
			Document: aigram.FileUpload(aigram.UploadFile{
				Name:   name,
				Reader: file,
			}),
			Caption:     "Original caption from ai-gram",
			ReplyMarkup: markup,
		})
		return message, "media_path", err
	}

	message, err := b.SendDocument(ctx, aigram.SendDocumentParams{
		ChatID: aigram.ChatIDInt(chatID),
		Document: aigram.FileUpload(aigram.UploadFile{
			Name:   "aigram-caption-demo.txt",
			Reader: strings.NewReader("ai-gram caption demo\n"),
		}),
		Caption:     "Original caption from ai-gram",
		ReplyMarkup: markup,
	})
	return message, "generated_document", err
}

func deleteCallbackMessage(ctx context.Context, b *aigram.Bot, update telegram.Update, text string) error {
	callback := update.CallbackQuery
	if callback == nil {
		return nil
	}
	logSafeUpdate(update, "callback_query")
	if _, err := b.AnswerCallbackQuery(ctx, aigram.AnswerCallbackQueryParams{CallbackQueryID: callback.ID, Text: text}); err != nil {
		return err
	}
	log.Printf("webhook action=answer_callback_query ok=true update_id=%d callback_data=%s", update.UpdateID, safeCallbackData(update))
	if callback.Message == nil {
		return nil
	}

	ok, err := b.DeleteMessage(ctx, aigram.DeleteMessageParams{
		ChatID:    aigram.ChatIDInt(callback.Message.Chat.ID),
		MessageID: callback.Message.MessageID,
	})
	if err != nil {
		return err
	}
	if !ok {
		return errors.New("delete_message returned false")
	}
	log.Printf("webhook action=delete_message ok=true update_id=%d chat_id=%d message_id=%d", update.UpdateID, callback.Message.Chat.ID, callback.Message.MessageID)
	return nil
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
	case "demo:edit", "demo:remove_keyboard", "demo:caption", "demo:edit_caption", "demo:delete_media", "demo:delete", "demo:copy", "demo:forward":
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
