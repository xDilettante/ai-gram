// Example 09_group_admin shows safe group/admin identity helpers.
package main

import (
	"context"
	"fmt"
	"log"
	"strings"

	aigram "github.com/xDilettante/ai-gram"
	"github.com/xDilettante/ai-gram/dispatch"
	"github.com/xDilettante/ai-gram/examples/internal/exampleutil"
	"github.com/xDilettante/ai-gram/middleware"
	"github.com/xDilettante/ai-gram/telegram"
	"github.com/xDilettante/ai-gram/transport/longpoll"
)

const logComponent = "group_admin"

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

	accessConfig, err := exampleutil.AccessConfigFromEnv()
	if err != nil {
		return err
	}
	accessController := exampleutil.NewAccessController(accessConfig)

	dp := dispatch.New(dispatch.WithErrorHandler(func(ctx context.Context, update telegram.Update, err error) {
		log.Printf("%s event=handler_error update_id=%d err=%v", logComponent, update.UpdateID, err)
	}))
	dp.Use(middleware.AccessWithPolicy(accessController, accessDenyHandler(b)))

	if err := registerRoutes(dp, b, accessController); err != nil {
		return err
	}

	runner, err := longpoll.New(b, dp, longpoll.Config{
		Timeout:        30,
		AllowedUpdates: []string{"message"},
		OnError: func(ctx context.Context, err error) {
			log.Printf("%s event=longpoll_error err=%v", logComponent, err)
		},
	})
	if err != nil {
		return err
	}

	log.Printf("%s event=started access_mode=%s", logComponent, accessController.Mode())
	if err := runner.Run(ctx); err != nil && err != context.Canceled {
		return err
	}
	log.Printf("%s event=stopped", logComponent)
	return nil
}

func registerRoutes(dp *dispatch.Dispatcher, b *aigram.Bot, accessController *exampleutil.AccessController) error {
	if err := dp.OnCommandFunc("start", helpHandler(b)); err != nil {
		return err
	}
	if err := dp.OnCommandFunc("help", helpHandler(b)); err != nil {
		return err
	}
	if err := dp.OnCommandFunc("whoami", whoamiHandler(b)); err != nil {
		return err
	}
	if err := dp.OnCommandFunc("chat", chatHandler(b)); err != nil {
		return err
	}
	if err := dp.OnCommandFunc("replytarget", replyTargetHandler(b)); err != nil {
		return err
	}
	if err := dp.OnCommandFunc("admin_panel", adminPanelHandler(b, accessController)); err != nil {
		return err
	}
	return dp.OnMessageFunc(func(ctx context.Context, update telegram.Update) error {
		message := update.EffectiveMessage()
		if message == nil || strings.HasPrefix(strings.TrimSpace(message.Text), "/") {
			return nil
		}
		logSafeUpdate(update, "message_ignored")
		return nil
	})
}

func helpHandler(b *aigram.Bot) dispatch.HandlerFunc {
	return func(ctx context.Context, update telegram.Update) error {
		logSafeUpdate(update, "help")
		return reply(ctx, b, update, strings.Join([]string{
			"Group admin identity example.",
			"",
			"Commands:",
			"/whoami - show the user/chat actor Telegram exposed",
			"/chat - show the effective chat",
			"/replytarget - reply to a message and show its actor",
			"/admin_panel - admin-only read-only identity panel",
			"",
			"This example does not ban, restrict, delete, or approve users.",
		}, "\n"))
	}
}

func whoamiHandler(b *aigram.Bot) dispatch.HandlerFunc {
	return func(ctx context.Context, update telegram.Update) error {
		logSafeUpdate(update, "whoami")
		return reply(ctx, b, update, "Actor:\n"+actorText(update.Actor()))
	}
}

func chatHandler(b *aigram.Bot) dispatch.HandlerFunc {
	return func(ctx context.Context, update telegram.Update) error {
		logSafeUpdate(update, "chat")
		return reply(ctx, b, update, "Chat:\n"+chatText(update.EffectiveChat()))
	}
}

func replyTargetHandler(b *aigram.Bot) dispatch.HandlerFunc {
	return func(ctx context.Context, update telegram.Update) error {
		logSafeUpdate(update, "replytarget")
		target := update.ReplyTarget()
		if target.IsZero() {
			return reply(ctx, b, update, "Reply to another message and run /replytarget again.")
		}
		return reply(ctx, b, update, "Reply target:\n"+actorText(target))
	}
}

func adminPanelHandler(b *aigram.Bot, accessController *exampleutil.AccessController) dispatch.HandlerFunc {
	return func(ctx context.Context, update telegram.Update) error {
		logSafeUpdate(update, "admin_panel")
		if !accessController.IsAdmin(update) {
			return reply(ctx, b, update, "Admin-only command. Configure AIGRAM_ADMIN_USER_IDS for this example.")
		}

		text := strings.Join([]string{
			"Admin identity panel:",
			"access_mode=" + string(accessController.Mode()),
			"",
			"Actor:",
			actorText(update.Actor()),
			"",
			"Chat:",
			chatText(update.EffectiveChat()),
			"",
			"No moderation actions are performed by this example.",
		}, "\n")
		return reply(ctx, b, update, text)
	}
}

func accessDenyHandler(b *aigram.Bot) func(context.Context, telegram.Update) error {
	return func(ctx context.Context, update telegram.Update) error {
		logSafeUpdate(update, "access_denied")
		message := update.EffectiveMessage()
		if message == nil {
			return nil
		}
		return sendMessage(ctx, b, message, "Access denied. Configure AIGRAM_ACCESS_MODE or AIGRAM_ADMIN_USER_IDS for this example.")
	}
}

func reply(ctx context.Context, b *aigram.Bot, update telegram.Update, text string) error {
	message := update.EffectiveMessage()
	if message == nil {
		return nil
	}
	return sendMessage(ctx, b, message, text)
}

func sendMessage(ctx context.Context, b *aigram.Bot, message *telegram.Message, text string) error {
	_, err := b.SendMessage(ctx, aigram.SendMessageParams{
		ChatID:          aigram.ChatIDInt(message.Chat.ID),
		MessageThreadID: message.MessageThreadID,
		Text:            text,
		ReplyParameters: &aigram.ReplyParameters{MessageID: message.MessageID},
	})
	return err
}

func logSafeUpdate(update telegram.Update, action string) {
	chat := update.EffectiveChat()
	actor := update.Actor()
	message := update.EffectiveMessage()

	chatID := "0"
	actorUserID := "0"
	actorChatID := "0"
	command := ""
	messageID := int64(0)
	if chat != nil {
		chatID = exampleutil.MaskInt64(chat.ID)
	}
	if actor.User != nil {
		actorUserID = exampleutil.MaskInt64(actor.User.ID)
	}
	if actor.Chat != nil {
		actorChatID = exampleutil.MaskInt64(actor.Chat.ID)
	}
	if message != nil {
		command = message.Command()
		messageID = message.MessageID
	}

	log.Printf("%s event=update action=%s update_id=%d message_id=%d chat_id=%s actor_user_id=%s actor_chat_id=%s anonymous_admin=%t command=%s",
		logComponent, action, update.UpdateID, messageID, chatID, actorUserID, actorChatID, actor.AnonymousAdmin, command)
}

func actorText(actor telegram.Actor) string {
	if actor.IsZero() {
		return "type=unknown"
	}

	var lines []string
	if actor.User != nil {
		lines = append(lines, "user_id="+exampleutil.MaskInt64(actor.User.ID))
		lines = append(lines, "user_is_bot="+fmt.Sprint(actor.User.IsBot))
		lines = append(lines, "user_username="+safeUsername(actor.User.Username))
	}
	if actor.Chat != nil {
		lines = append(lines, "chat_id="+exampleutil.MaskInt64(actor.Chat.ID))
		lines = append(lines, "chat_type="+safeValue(actor.Chat.Type))
		lines = append(lines, "chat_username="+safeUsername(actor.Chat.Username))
	}
	lines = append(lines, "anonymous_admin="+fmt.Sprint(actor.AnonymousAdmin))
	return strings.Join(lines, "\n")
}

func chatText(chat *telegram.Chat) string {
	if chat == nil {
		return "chat=none"
	}
	return strings.Join([]string{
		"chat_id=" + exampleutil.MaskInt64(chat.ID),
		"chat_type=" + safeValue(chat.Type),
		"chat_username=" + safeUsername(chat.Username),
		"chat_is_forum=" + fmt.Sprint(chat.IsForum),
	}, "\n")
}

func safeUsername(username string) string {
	if strings.TrimSpace(username) == "" {
		return "<none>"
	}
	return "@" + strings.TrimPrefix(strings.TrimSpace(username), "@")
}

func safeValue(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return "<none>"
	}
	return value
}
