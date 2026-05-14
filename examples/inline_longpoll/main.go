// Example inline_longpoll demonstrates inline keyboards and callback answers.
//
// Required env: AIGRAM_BOT_TOKEN.
// Optional env: AIGRAM_BASE_URL, AIGRAM_FILE_BASE_URL, AIGRAM_ACCESS_MODE,
// AIGRAM_ADMIN_USER_IDS, AIGRAM_ALLOWED_USER_IDS, AIGRAM_ALLOWED_CHAT_IDS.
// It deletes any configured webhook before polling. Send /start to the bot. Stop with Ctrl+C or SIGTERM.
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/xDilettante/ai-gram"
	"github.com/xDilettante/ai-gram/dispatch"
	"github.com/xDilettante/ai-gram/examples/internal/exampleutil"
	"github.com/xDilettante/ai-gram/middleware"
	"github.com/xDilettante/ai-gram/telegram"
	"github.com/xDilettante/ai-gram/transport/longpoll"
)

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
	chatID := "0"
	userID := "0"
	if chat != nil {
		chatID = exampleutil.MaskInt64(chat.ID)
	}
	if user != nil {
		userID = exampleutil.MaskInt64(user.ID)
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
		log.Printf("longpoll update_id=%d update_type=%s matched=%s chat_id=%s from_user_id=%s command=%s has_text=%t has_media=%t callback_data=%s", update.UpdateID, updateType, matched, chatID, userID, command, hasText, hasMedia, callbackData)
		return
	}
	log.Printf("longpoll update_id=%d update_type=%s matched=%s chat_id=%s from_user_id=%s has_text=%t has_media=%t callback_data=%s", update.UpdateID, updateType, matched, chatID, userID, hasText, hasMedia, callbackData)
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

func smokeExitAfterUpdate() bool {
	return os.Getenv("AIGRAM_SMOKE_EXIT_AFTER_UPDATE") == "1"
}

func smokeIDs(update telegram.Update) (string, string) {
	chatID := "0"
	userID := "0"
	if chat := update.EffectiveChat(); chat != nil {
		chatID = exampleutil.MaskInt64(chat.ID)
	}
	if user := update.EffectiveUser(); user != nil {
		userID = exampleutil.MaskInt64(user.ID)
	}
	return chatID, userID
}

func registerAccessCommands(dp *dispatch.Dispatcher, b *aigram.Bot, controller *exampleutil.AccessController, logPrefix string) error {
	if err := dp.OnCommandFunc("access_status", func(ctx context.Context, update telegram.Update) error {
		if !controller.IsAdmin(update) {
			return accessDenyHandler(b, logPrefix)(ctx, update)
		}
		message := update.EffectiveMessage()
		if message == nil {
			return nil
		}
		mode := controller.Mode()
		logSafeUpdate(update, "command")
		_, err := b.SendMessage(ctx, aigram.SendMessageParams{
			ChatID:          aigram.ChatIDInt(message.Chat.ID),
			MessageThreadID: message.MessageThreadID,
			Text:            fmt.Sprintf("Access mode: %s", mode),
			ReplyParameters: &aigram.ReplyParameters{MessageID: message.MessageID},
		})
		if err != nil {
			return err
		}
		userID := "0"
		if user := update.EffectiveUser(); user != nil {
			userID = exampleutil.MaskInt64(user.ID)
		}
		log.Printf("%s action=access_status mode=%s update_id=%d by_user_id=%s", logPrefix, mode, update.UpdateID, userID)
		return nil
	}); err != nil {
		return err
	}
	if err := dp.OnCommandFunc("access_open", func(ctx context.Context, update telegram.Update) error {
		return setAccessMode(ctx, b, controller, logPrefix, update, middleware.AccessModePublic)
	}); err != nil {
		return err
	}
	if err := dp.OnCommandFunc("access_close", func(ctx context.Context, update telegram.Update) error {
		return setAccessMode(ctx, b, controller, logPrefix, update, middleware.AccessModeAdmin)
	}); err != nil {
		return err
	}
	return nil
}

func setAccessMode(ctx context.Context, b *aigram.Bot, controller *exampleutil.AccessController, logPrefix string, update telegram.Update, mode middleware.AccessMode) error {
	if !controller.IsAdmin(update) {
		return accessDenyHandler(b, logPrefix)(ctx, update)
	}
	message := update.EffectiveMessage()
	if message == nil {
		return nil
	}
	controller.SetMode(mode)
	logSafeUpdate(update, "command")
	_, err := b.SendMessage(ctx, aigram.SendMessageParams{
		ChatID:          aigram.ChatIDInt(message.Chat.ID),
		MessageThreadID: message.MessageThreadID,
		Text:            fmt.Sprintf("Access mode changed: %s", mode),
		ReplyParameters: &aigram.ReplyParameters{MessageID: message.MessageID},
	})
	if err != nil {
		return err
	}
	userID := "0"
	if user := update.EffectiveUser(); user != nil {
		userID = exampleutil.MaskInt64(user.ID)
	}
	log.Printf("%s action=access_mode_changed mode=%s update_id=%d by_user_id=%s", logPrefix, mode, update.UpdateID, userID)
	return nil
}

func accessDenyHandler(b *aigram.Bot, logPrefix string) func(context.Context, telegram.Update) error {
	return func(ctx context.Context, update telegram.Update) error {
		chatID := int64(0)
		userID := int64(0)
		messageThreadID := int64(0)
		if chat := update.EffectiveChat(); chat != nil {
			chatID = chat.ID
		}
		if user := update.EffectiveUser(); user != nil {
			userID = user.ID
		}
		if message := update.EffectiveMessage(); message != nil {
			messageThreadID = message.MessageThreadID
		}
		log.Printf("%s action=access_denied update_id=%d chat_id=%s from_user_id=%s", logPrefix, update.UpdateID, exampleutil.MaskInt64(chatID), exampleutil.MaskInt64(userID))
		if chatID == 0 {
			return nil
		}
		_, err := b.SendMessage(ctx, aigram.SendMessageParams{
			ChatID:          aigram.ChatIDInt(chatID),
			MessageThreadID: messageThreadID,
			Text:            "Access denied.",
		})
		return err
	}
}

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
		return fmt.Errorf("delete webhook before long polling: %w", err)
	}

	dp := dispatch.New(dispatch.WithErrorHandler(func(ctx context.Context, update telegram.Update, err error) {
		log.Printf("handler error for update %d: %v", update.UpdateID, err)
	}))
	accessConfig, err := exampleutil.AccessConfigFromEnv()
	if err != nil {
		return err
	}
	accessController := exampleutil.NewAccessController(accessConfig)
	dp.Use(middleware.AccessWithPolicy(accessController, accessDenyHandler(b, "longpoll")))

	if err := registerAccessCommands(dp, b, accessController, "longpoll"); err != nil {
		return err
	}
	if err := dp.OnCommandFunc("start", func(ctx context.Context, update telegram.Update) error {
		message := update.EffectiveMessage()
		if message == nil {
			return nil
		}
		logSafeUpdate(update, "command")
		_, err := b.SendMessage(ctx, aigram.SendMessageParams{
			ChatID:      aigram.ChatIDInt(message.Chat.ID),
			Text:        "Long polling bot is running. Choose an action:",
			ReplyMarkup: demoKeyboard(),
		})
		if err == nil {
			log.Printf("longpoll action=send_message update_id=%d ok=true", update.UpdateID)
		}
		return err
	}); err != nil {
		return err
	}
	if err := dp.OnMessageFunc(func(ctx context.Context, update telegram.Update) error {
		message := update.EffectiveMessage()
		if message == nil || message.Text == "" || message.Command() != "" {
			return nil
		}
		logSafeUpdate(update, "message")
		chatID, userID := smokeIDs(update)
		if smokeExitAfterUpdate() {
			fmt.Printf("AIGRAM_SMOKE_UPDATE_RECEIVED update_id=%d chat_id=%s from_user_id=%s\n", update.UpdateID, chatID, userID)
		}
		_, err := b.SendMessage(ctx, aigram.SendMessageParams{
			ChatID:          aigram.ChatIDInt(message.Chat.ID),
			MessageThreadID: message.MessageThreadID,
			Text:            "echo received",
			ReplyParameters: &aigram.ReplyParameters{MessageID: message.MessageID},
		})
		if err != nil {
			return err
		}
		log.Printf("longpoll action=send_message ok=true update_id=%d chat_id=%s reply_to_message_id=%d", update.UpdateID, exampleutil.MaskInt64(message.Chat.ID), message.MessageID)
		if smokeExitAfterUpdate() {
			fmt.Printf("AIGRAM_SMOKE_REPLY_SENT update_id=%d chat_id=%s reply_to_message_id=%d\n", update.UpdateID, chatID, message.MessageID)
			fmt.Printf("AIGRAM_SMOKE_OK update_id=%d chat_id=%s from_user_id=%s\n", update.UpdateID, chatID, userID)
			stop()
		}
		return nil
	}); err != nil {
		return err
	}
	if err := dp.OnCallbackDataFunc("demo:edit", func(ctx context.Context, update telegram.Update) error {
		callback := update.CallbackQuery
		if callback == nil {
			return nil
		}
		logSafeUpdate(update, "callback_query")
		if _, err := b.AnswerCallbackQuery(ctx, aigram.AnswerCallbackQueryParams{
			CallbackQueryID: callback.ID,
			Text:            "Editing message",
		}); err != nil {
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
		if err == nil {
			log.Printf("longpoll action=edit_message_text update_id=%d ok=true", update.UpdateID)
		}
		return err
	}); err != nil {
		return err
	}
	if err := dp.OnCallbackDataFunc("demo:remove_keyboard", func(ctx context.Context, update telegram.Update) error {
		callback := update.CallbackQuery
		if callback == nil {
			return nil
		}
		logSafeUpdate(update, "callback_query")
		if _, err := b.AnswerCallbackQuery(ctx, aigram.AnswerCallbackQueryParams{
			CallbackQueryID: callback.ID,
			Text:            "Removing keyboard",
		}); err != nil {
			return err
		}
		if callback.Message == nil {
			return nil
		}
		_, err := b.EditMessageReplyMarkup(ctx, aigram.EditMessageReplyMarkupParams{
			Target: aigram.EditTargetChat(aigram.ChatIDInt(callback.Message.Chat.ID), callback.Message.MessageID),
		})
		if err == nil {
			log.Printf("longpoll action=edit_message_reply_markup update_id=%d ok=true", update.UpdateID)
		}
		return err
	}); err != nil {
		return err
	}

	runner, err := longpoll.New(b, dp, longpoll.Config{Timeout: 30})
	if err != nil {
		return err
	}
	log.Println("inline long polling started")
	if err := runner.Run(ctx); err != nil && err != context.Canceled {
		return err
	}
	log.Println("inline long polling stopped")
	return nil
}
