// Example 10_moderation_skeleton shows a dry-run moderation workflow.
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

const (
	logComponent = "moderation_skeleton"
	dryRun       = true
)

type moderationPreview struct {
	Action          string
	Reason          string
	Chat            *telegram.Chat
	Reporter        telegram.Actor
	Target          telegram.Actor
	TargetMessageID int64
	DryRun          bool
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
		AllowedUpdates: []string{"message", "chat_join_request"},
		OnError: func(ctx context.Context, err error) {
			log.Printf("%s event=longpoll_error err=%v", logComponent, err)
		},
	})
	if err != nil {
		return err
	}

	log.Printf("%s event=started access_mode=%s dry_run=%t", logComponent, accessController.Mode(), dryRun)
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
	if err := dp.OnCommandFunc("mod_status", modStatusHandler(b, accessController)); err != nil {
		return err
	}
	if err := dp.OnCommandFunc("report", reportHandler(b)); err != nil {
		return err
	}
	if err := dp.OnCommandFunc("mod_preview", modPreviewHandler(b, accessController)); err != nil {
		return err
	}
	if err := dp.OnChatJoinRequestFunc(joinRequestHandler()); err != nil {
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
			"Moderation skeleton is running in dry-run mode.",
			"",
			"Commands:",
			"/report - reply to a message and create a dry-run report",
			"/mod_preview - admin-only dry-run moderation preview for a replied message",
			"/mod_status - admin-only status",
			"",
			"This example does not call ban, restrict, delete, approve, or decline methods.",
		}, "\n"))
	}
}

func modStatusHandler(b *aigram.Bot, accessController *exampleutil.AccessController) dispatch.HandlerFunc {
	return func(ctx context.Context, update telegram.Update) error {
		logSafeUpdate(update, "mod_status")
		if !accessController.IsAdmin(update) {
			return reply(ctx, b, update, "Admin-only command. Configure AIGRAM_ADMIN_USER_IDS for this example.")
		}
		return reply(ctx, b, update, strings.Join([]string{
			"Moderation status:",
			"dry_run=true",
			"access_mode=" + string(accessController.Mode()),
			"destructive_actions=disabled",
		}, "\n"))
	}
}

func reportHandler(b *aigram.Bot) dispatch.HandlerFunc {
	return func(ctx context.Context, update telegram.Update) error {
		logSafeUpdate(update, "report")
		message := update.EffectiveMessage()
		if message == nil {
			return nil
		}
		if update.ReplyTarget().IsZero() {
			return reply(ctx, b, update, "Reply to the message you want to report and send /report again.")
		}

		preview := previewFromUpdate(update, "report", reasonFromMessage(message, "user report"))
		logModerationPreview(update, preview)
		return reply(ctx, b, update, "Dry-run report recorded:\n"+previewText(preview))
	}
}

func modPreviewHandler(b *aigram.Bot, accessController *exampleutil.AccessController) dispatch.HandlerFunc {
	return func(ctx context.Context, update telegram.Update) error {
		logSafeUpdate(update, "mod_preview")
		message := update.EffectiveMessage()
		if message == nil {
			return nil
		}
		if !accessController.IsAdmin(update) {
			return reply(ctx, b, update, "Admin-only command. Configure AIGRAM_ADMIN_USER_IDS for this example.")
		}
		if update.ReplyTarget().IsZero() {
			return reply(ctx, b, update, "Reply to a message and send /mod_preview to inspect the dry-run moderation plan.")
		}

		preview := previewFromUpdate(update, "moderator_preview", reasonFromMessage(message, "manual review"))
		logModerationPreview(update, preview)
		return reply(ctx, b, update, "Dry-run moderation preview:\n"+previewText(preview))
	}
}

func joinRequestHandler() dispatch.HandlerFunc {
	return func(ctx context.Context, update telegram.Update) error {
		request := update.ChatJoinRequest
		if request == nil {
			return nil
		}
		log.Printf("%s event=join_request_seen update_id=%d chat_id=%s actor_user_id=%s dry_run=%t",
			logComponent, update.UpdateID, exampleutil.MaskInt64(request.Chat.ID), exampleutil.MaskInt64(request.From.ID), dryRun)
		return nil
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

func previewFromUpdate(update telegram.Update, action, reason string) moderationPreview {
	message := update.EffectiveMessage()
	var chat *telegram.Chat
	targetMessageID := int64(0)
	if message != nil {
		chat = &message.Chat
		if message.ReplyToMessage != nil {
			targetMessageID = message.ReplyToMessage.MessageID
		}
	}
	return moderationPreview{
		Action:          action,
		Reason:          reason,
		Chat:            chat,
		Reporter:        update.Actor(),
		Target:          update.ReplyTarget(),
		TargetMessageID: targetMessageID,
		DryRun:          dryRun,
	}
}

func reasonFromMessage(message *telegram.Message, fallback string) string {
	reason := strings.TrimSpace(message.CommandArguments())
	if reason == "" {
		return fallback
	}
	return reason
}

func logModerationPreview(update telegram.Update, preview moderationPreview) {
	chatID := "0"
	reporterUserID := "0"
	reporterChatID := "0"
	targetUserID := "0"
	targetChatID := "0"
	if preview.Chat != nil {
		chatID = exampleutil.MaskInt64(preview.Chat.ID)
	}
	if preview.Reporter.User != nil {
		reporterUserID = exampleutil.MaskInt64(preview.Reporter.User.ID)
	}
	if preview.Reporter.Chat != nil {
		reporterChatID = exampleutil.MaskInt64(preview.Reporter.Chat.ID)
	}
	if preview.Target.User != nil {
		targetUserID = exampleutil.MaskInt64(preview.Target.User.ID)
	}
	if preview.Target.Chat != nil {
		targetChatID = exampleutil.MaskInt64(preview.Target.Chat.ID)
	}

	log.Printf("%s event=moderation_preview update_id=%d action=%s chat_id=%s reporter_user_id=%s reporter_chat_id=%s target_user_id=%s target_chat_id=%s target_message_id=%d dry_run=%t reason=%s",
		logComponent, update.UpdateID, preview.Action, chatID, reporterUserID, reporterChatID, targetUserID, targetChatID, preview.TargetMessageID, preview.DryRun, safeLogValue(preview.Reason))
}

func previewText(preview moderationPreview) string {
	return strings.Join([]string{
		"action=" + safeValue(preview.Action),
		"dry_run=" + fmt.Sprint(preview.DryRun),
		"reason=" + safeValue(preview.Reason),
		"chat=" + chatLabel(preview.Chat),
		"reporter=" + actorLabel(preview.Reporter),
		"target=" + actorLabel(preview.Target),
		fmt.Sprintf("target_message_id=%d", preview.TargetMessageID),
		"would_delete_message=false",
		"would_restrict_user=false",
		"would_ban_user=false",
		"would_approve_or_decline_join_request=false",
	}, "\n")
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

func actorLabel(actor telegram.Actor) string {
	if actor.User != nil {
		return "user:" + exampleutil.MaskInt64(actor.User.ID)
	}
	if actor.Chat != nil {
		return "chat:" + exampleutil.MaskInt64(actor.Chat.ID)
	}
	return "unknown"
}

func chatLabel(chat *telegram.Chat) string {
	if chat == nil {
		return "none"
	}
	return exampleutil.MaskInt64(chat.ID) + ":" + safeValue(chat.Type)
}

func safeLogValue(value string) string {
	value = safeValue(value)
	if len(value) > 80 {
		return value[:80] + "..."
	}
	return value
}

func safeValue(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return "<none>"
	}
	return strings.NewReplacer("\n", " ", "\r", " ", "\t", " ").Replace(value)
}
