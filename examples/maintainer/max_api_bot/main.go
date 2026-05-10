// Example max_api_bot runs a broad maintainer live-smoke bot for ai-gram.
//
// Required env for once mode: AIGRAM_BOT_TOKEN and AIGRAM_CHAT_ID.
// Optional env: AIGRAM_BASE_URL, AIGRAM_FILE_BASE_URL, AIGRAM_MAX_API_MODE,
// AIGRAM_MAX_API_ALLOW_COMMANDS, AIGRAM_MAX_API_ALLOW_DELETE,
// AIGRAM_MAX_API_DELETE_WEBHOOK, AIGRAM_MAX_API_RUN_SECONDS.
package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/xDilettante/ai-gram"
	"github.com/xDilettante/ai-gram/examples/internal/exampleutil"
	"github.com/xDilettante/ai-gram/telegram"
	"github.com/xDilettante/ai-gram/transport/longpoll"
)

const (
	modeOnce = "once"
	modePoll = "poll"

	eventStarted = "max_api_bot_started"
	eventStepOK  = "max_api_step_ok"
	eventStepErr = "max_api_step_error"
)

type config struct {
	mode           string
	chatID         aigram.ChatID
	runSeconds     int
	allowCommands  bool
	allowDelete    bool
	deleteWebhook  bool
	allowedUpdates []string
}

type app struct {
	bot    *aigram.Bot
	config config
}

func main() {
	log.SetFlags(0)
	if err := run(); err != nil {
		logJSON("error", "max_api_bot_failed", "max api bot failed", fields{"error": err.Error()})
		os.Exit(1)
	}
}

func run() error {
	b, err := exampleutil.NewBotFromEnv()
	if err != nil {
		return err
	}

	cfg, err := configFromEnv()
	if err != nil {
		return err
	}

	a := &app{bot: b, config: cfg}
	logJSON("info", eventStarted, "max api smoke bot started", fields{
		"mode":            cfg.mode,
		"run_seconds":     cfg.runSeconds,
		"allow_commands":  cfg.allowCommands,
		"allow_delete":    cfg.allowDelete,
		"delete_webhook":  cfg.deleteWebhook,
		"allowed_updates": strings.Join(cfg.allowedUpdates, ","),
		"configured_chat": chatIDForLog(cfg.chatID),
	})

	ctx := context.Background()
	switch cfg.mode {
	case modeOnce:
		ctx, cancel := context.WithTimeout(ctx, time.Duration(cfg.runSeconds)*time.Second)
		defer cancel()
		return a.runSuite(ctx, cfg.chatID, "once")
	case modePoll:
		return a.runPolling(ctx)
	default:
		return fmt.Errorf("AIGRAM_MAX_API_MODE must be %q or %q", modeOnce, modePoll)
	}
}

func configFromEnv() (config, error) {
	mode := strings.TrimSpace(exampleutil.OptionalEnv("AIGRAM_MAX_API_MODE", modeOnce))
	runSeconds, err := intEnv("AIGRAM_MAX_API_RUN_SECONDS", 120)
	if err != nil {
		return config{}, err
	}

	var chatID aigram.ChatID
	if raw := strings.TrimSpace(os.Getenv("AIGRAM_CHAT_ID")); raw != "" {
		parsed, err := exampleutil.ParseChatID(raw)
		if err != nil {
			return config{}, err
		}
		chatID = parsed
	}
	if mode == modeOnce && !chatIDValid(chatID) {
		return config{}, errors.New("AIGRAM_CHAT_ID is required in once mode")
	}

	return config{
		mode:          mode,
		chatID:        chatID,
		runSeconds:    runSeconds,
		allowCommands: boolEnv("AIGRAM_MAX_API_ALLOW_COMMANDS"),
		allowDelete:   boolEnv("AIGRAM_MAX_API_ALLOW_DELETE"),
		deleteWebhook: boolEnv("AIGRAM_MAX_API_DELETE_WEBHOOK"),
		allowedUpdates: []string{
			"message",
			"edited_message",
			"callback_query",
			"inline_query",
			"chosen_inline_result",
			"poll",
			"poll_answer",
			"my_chat_member",
			"chat_member",
			"chat_join_request",
			"message_reaction",
			"message_reaction_count",
			"chat_boost",
			"removed_chat_boost",
			"business_connection",
			"business_message",
			"edited_business_message",
			"deleted_business_messages",
			"purchased_paid_media",
		},
	}, nil
}

func (a *app) runPolling(ctx context.Context) error {
	ctx, stop := exampleutil.SignalContext()
	defer stop()

	if a.config.deleteWebhook {
		_, err := a.bot.DeleteWebhook(ctx, aigram.DeleteWebhookParams{DropPendingUpdates: false})
		if err != nil {
			return fmt.Errorf("delete webhook before polling: %w", err)
		}
		logJSON("info", "max_api_delete_webhook_ok", "webhook disabled before polling", nil)
	}

	runner, err := longpoll.New(a.bot, longpoll.HandlerFunc(a.handleUpdate), longpoll.Config{
		Timeout:        30,
		Limit:          50,
		AllowedUpdates: a.config.allowedUpdates,
		OnError: func(ctx context.Context, err error) {
			logJSON("error", "max_api_poll_error", "long polling error", fields{"error": err.Error()})
		},
	})
	if err != nil {
		return err
	}

	logJSON("info", "max_api_polling_started", "polling started", fields{
		"delete_webhook": a.config.deleteWebhook,
	})
	err = runner.Run(ctx)
	if errors.Is(err, context.Canceled) {
		logJSON("info", "max_api_polling_stopped", "polling stopped by context", nil)
		return nil
	}
	return err
}

func (a *app) handleUpdate(ctx context.Context, update aigram.Update) error {
	logUpdate(update)

	if query := update.CallbackQuery; query != nil {
		return a.handleCallback(ctx, update)
	}
	message := update.EffectiveMessage()
	if message == nil {
		return nil
	}

	switch strings.TrimSpace(message.Text) {
	case "/start", "/help":
		return a.sendHelp(ctx, aigram.ChatIDInt(message.Chat.ID))
	case "/smoke":
		return a.runSuite(ctx, aigram.ChatIDInt(message.Chat.ID), "command")
	case "/media":
		return a.runMediaSuite(ctx, aigram.ChatIDInt(message.Chat.ID), "command")
	case "/status":
		return a.runReadOnlySuite(ctx, aigram.ChatIDInt(message.Chat.ID), "command")
	case "/remove_keyboard":
		_, err := a.bot.SendMessage(ctx, aigram.SendMessageParams{
			ChatID:      aigram.ChatIDInt(message.Chat.ID),
			Text:        "Reply keyboard removed.",
			ReplyMarkup: aigram.RemoveKeyboard(false),
		})
		return err
	default:
		if strings.TrimSpace(message.Text) != "" {
			_, err := a.bot.SendMessage(ctx, aigram.SendMessageParams{
				ChatID: aigram.ChatIDInt(message.Chat.ID),
				Text:   "Logged your update. Send /smoke, /media, /status, or /help.",
				ReplyParameters: &aigram.ReplyParameters{
					MessageID: message.MessageID,
				},
			})
			return err
		}
		return nil
	}
}

func (a *app) handleCallback(ctx context.Context, update aigram.Update) error {
	query := update.CallbackQuery
	if query == nil {
		return nil
	}
	_, _ = a.bot.AnswerCallbackQuery(ctx, aigram.AnswerCallbackQueryParams{
		CallbackQueryID: query.ID,
		Text:            "Action received.",
	})

	chat := update.EffectiveChat()
	if chat == nil {
		return nil
	}
	chatID := aigram.ChatIDInt(chat.ID)
	switch query.Data {
	case "smoke:safe":
		return a.runSuite(ctx, chatID, "callback")
	case "smoke:media":
		return a.runMediaSuite(ctx, chatID, "callback")
	case "smoke:status":
		return a.runReadOnlySuite(ctx, chatID, "callback")
	case "smoke:help":
		return a.sendHelp(ctx, chatID)
	default:
		logJSON("warn", "max_api_unknown_callback", "unknown callback data", fields{"callback_data": query.Data})
		return nil
	}
}

func (a *app) sendHelp(ctx context.Context, chatID aigram.ChatID) error {
	keyboard := aigram.NewInlineKeyboard(
		[]aigram.InlineKeyboardButton{
			aigram.InlineButtonCallback("Run safe smoke", "smoke:safe"),
			aigram.InlineButtonCallback("Media smoke", "smoke:media"),
		},
		[]aigram.InlineKeyboardButton{
			aigram.InlineButtonCallback("Status", "smoke:status"),
			aigram.InlineButtonCopyText("Copy command", "/smoke"),
		},
		[]aigram.InlineKeyboardButton{
			aigram.InlineButtonSwitchInlineQueryCurrentChat("Inline query", "ai-gram"),
			aigram.InlineButtonURL("Project", "https://github.com/xDilettante/ai-gram"),
		},
	)
	replyKeyboard := aigram.NewReplyKeyboard(
		[]aigram.KeyboardButton{aigram.KeyboardButtonText("/smoke"), aigram.KeyboardButtonText("/media")},
		[]aigram.KeyboardButton{aigram.KeyboardButtonText("/status"), aigram.KeyboardButtonText("/remove_keyboard")},
		[]aigram.KeyboardButton{
			aigram.KeyboardButtonContact("Share contact"),
			aigram.KeyboardButtonLocation("Share location"),
			aigram.KeyboardButtonPoll("Create poll", "regular"),
		},
	)
	replyKeyboard.ResizeKeyboard = true
	replyKeyboard.InputFieldPlaceholder = "Try /smoke or press an inline button"

	if _, err := a.bot.SendMessage(ctx, aigram.SendMessageParams{
		ChatID:      chatID,
		Text:        "ai-gram max API smoke bot is ready. Use the inline buttons or commands.",
		ReplyMarkup: replyKeyboard,
	}); err != nil {
		return err
	}
	_, err := a.bot.SendMessage(ctx, aigram.SendMessageParams{
		ChatID:      chatID,
		Text:        "Control panel",
		ReplyMarkup: keyboard,
	})
	return err
}

func (a *app) runSuite(ctx context.Context, chatID aigram.ChatID, source string) error {
	logJSON("info", "max_api_suite_started", "safe smoke suite started", fields{"source": source, "chat_id": chatIDForLog(chatID)})
	summary := newSummary()

	a.step(ctx, summary, "get_me", func(ctx context.Context) error {
		me, err := a.bot.GetMe(ctx)
		if err != nil {
			return err
		}
		logJSON("info", eventStepOK, "getMe succeeded", fields{"step": "get_me", "bot_id": me.ID, "username": safeString(me.Username)})
		return nil
	})
	a.step(ctx, summary, "read_only", func(ctx context.Context) error {
		return a.runReadOnlySuite(ctx, chatID, source)
	})
	a.step(ctx, summary, "send_chat_action", func(ctx context.Context) error {
		ok, err := a.bot.SendChatAction(ctx, aigram.SendChatActionParams{ChatID: chatID, Action: aigram.ChatActionTyping})
		logJSON("info", eventStepOK, "sendChatAction result", fields{"step": "send_chat_action", "ok": ok})
		return err
	})

	var baseMessageID int64
	a.step(ctx, summary, "send_message_inline_keyboard", func(ctx context.Context) error {
		message, err := a.bot.SendMessage(ctx, aigram.SendMessageParams{
			ChatID:    chatID,
			Text:      "Safe smoke message. This message will be edited, copied, and forwarded.",
			ParseMode: "HTML",
			ReplyMarkup: aigram.NewInlineKeyboard([]aigram.InlineKeyboardButton{
				aigram.InlineButtonCallback("Status", "smoke:status"),
				aigram.InlineButtonCallback("Media", "smoke:media"),
			}),
		})
		if err != nil {
			return err
		}
		baseMessageID = message.MessageID
		logJSON("info", eventStepOK, "sendMessage succeeded", fields{"step": "send_message_inline_keyboard", "message_id": baseMessageID})
		return nil
	})
	if baseMessageID > 0 {
		a.step(ctx, summary, "edit_message_text", func(ctx context.Context) error {
			result, err := a.bot.EditMessageText(ctx, aigram.EditMessageTextParams{
				Target: aigram.EditTargetChat(chatID, baseMessageID),
				Text:   "Safe smoke message edited successfully.",
				ReplyMarkup: &aigram.InlineKeyboardMarkup{InlineKeyboard: [][]aigram.InlineKeyboardButton{{
					aigram.InlineButtonCallback("Run again", "smoke:safe"),
				}}},
			})
			if err != nil {
				return err
			}
			logJSON("info", eventStepOK, "editMessageText succeeded", fields{"step": "edit_message_text", "ok": result.IsOK(), "message_result": result.IsMessage()})
			return nil
		})
		a.step(ctx, summary, "edit_reply_markup", func(ctx context.Context) error {
			result, err := a.bot.EditMessageReplyMarkup(ctx, aigram.EditMessageReplyMarkupParams{
				Target:      aigram.EditTargetChat(chatID, baseMessageID),
				ReplyMarkup: nil,
			})
			if err != nil {
				return err
			}
			logJSON("info", eventStepOK, "editMessageReplyMarkup succeeded", fields{"step": "edit_reply_markup", "ok": result.IsOK()})
			return nil
		})
		a.step(ctx, summary, "copy_message", func(ctx context.Context) error {
			messageID, err := a.bot.CopyMessage(ctx, aigram.CopyMessageParams{
				ChatID:     chatID,
				FromChatID: chatID,
				MessageID:  baseMessageID,
				Caption:    "Copied by ai-gram smoke.",
			})
			if err != nil {
				return err
			}
			logJSON("info", eventStepOK, "copyMessage succeeded", fields{"step": "copy_message", "copied_message_id": messageID.MessageID})
			return nil
		})
		a.step(ctx, summary, "forward_message", func(ctx context.Context) error {
			message, err := a.bot.ForwardMessage(ctx, aigram.ForwardMessageParams{
				ChatID:     chatID,
				FromChatID: chatID,
				MessageID:  baseMessageID,
			})
			if err != nil {
				return err
			}
			logJSON("info", eventStepOK, "forwardMessage succeeded", fields{"step": "forward_message", "message_id": message.MessageID})
			return nil
		})
	}

	a.step(ctx, summary, "contact_location_venue", func(ctx context.Context) error {
		if _, err := a.bot.SendContact(ctx, aigram.SendContactParams{ChatID: chatID, PhoneNumber: "+10000000000", FirstName: "ai-gram", LastName: "Smoke"}); err != nil {
			return fmt.Errorf("send contact: %w", err)
		}
		if _, err := a.bot.SendLocation(ctx, aigram.SendLocationParams{ChatID: chatID, Latitude: 52.3676, Longitude: 4.9041}); err != nil {
			return fmt.Errorf("send location: %w", err)
		}
		if _, err := a.bot.SendVenue(ctx, aigram.SendVenueParams{ChatID: chatID, Latitude: 52.3676, Longitude: 4.9041, Title: "ai-gram smoke venue", Address: "ai-gram test address"}); err != nil {
			return fmt.Errorf("send venue: %w", err)
		}
		logJSON("info", eventStepOK, "contact/location/venue methods succeeded", fields{"step": "contact_location_venue"})
		return nil
	})

	a.step(ctx, summary, "poll_and_dice", func(ctx context.Context) error {
		falseValue := false
		poll, err := a.bot.SendPoll(ctx, aigram.SendPollParams{
			ChatID:      chatID,
			Question:    "ai-gram smoke poll?",
			Options:     []string{"yes", "no"},
			IsAnonymous: &falseValue,
		})
		if err != nil {
			return fmt.Errorf("send poll: %w", err)
		}
		stopped, err := a.bot.StopPoll(ctx, aigram.StopPollParams{ChatID: chatID, MessageID: poll.MessageID})
		if err != nil {
			return fmt.Errorf("stop poll: %w", err)
		}
		dice, err := a.bot.SendDice(ctx, aigram.SendDiceParams{ChatID: chatID, Emoji: "🎲"})
		if err != nil {
			return fmt.Errorf("send dice: %w", err)
		}
		logJSON("info", eventStepOK, "poll and dice methods succeeded", fields{"step": "poll_and_dice", "poll_id": redactID(stopped.ID), "dice_message_id": dice.MessageID})
		return nil
	})

	a.step(ctx, summary, "media", func(ctx context.Context) error {
		return a.runMediaSuite(ctx, chatID, source)
	})
	if a.config.allowCommands {
		a.step(ctx, summary, "set_get_my_commands", func(ctx context.Context) error {
			ok, err := a.bot.SetMyCommands(ctx, aigram.SetMyCommandsParams{Commands: []telegram.BotCommand{
				{Command: "start", Description: "Show smoke controls"},
				{Command: "smoke", Description: "Run safe smoke"},
				{Command: "media", Description: "Run media smoke"},
				{Command: "status", Description: "Run read-only checks"},
			}})
			if err != nil {
				return err
			}
			commands, err := a.bot.GetMyCommands(ctx, aigram.GetMyCommandsParams{})
			if err != nil {
				return err
			}
			logJSON("info", eventStepOK, "set/get my commands succeeded", fields{"step": "set_get_my_commands", "ok": ok, "command_count": len(commands)})
			return nil
		})
	} else {
		logJSON("info", "max_api_step_skipped", "command mutation skipped", fields{"step": "set_get_my_commands", "reason": "AIGRAM_MAX_API_ALLOW_COMMANDS is not enabled"})
	}

	_, sendErr := a.bot.SendMessage(ctx, aigram.SendMessageParams{
		ChatID: chatID,
		Text:   fmt.Sprintf("ai-gram smoke summary: passed=%d failed=%d", summary.passed, summary.failed),
	})
	logJSON("info", "max_api_suite_finished", "safe smoke suite finished", fields{"source": source, "passed": summary.passed, "failed": summary.failed})
	if sendErr != nil {
		return sendErr
	}
	if summary.failed > 0 {
		return fmt.Errorf("safe smoke suite had %d failed steps", summary.failed)
	}
	return nil
}

func (a *app) runReadOnlySuite(ctx context.Context, chatID aigram.ChatID, source string) error {
	summary := newSummary()
	logJSON("info", "max_api_readonly_started", "read-only suite started", fields{"source": source, "chat_id": chatIDForLog(chatID)})
	a.step(ctx, summary, "get_webhook_info", func(ctx context.Context) error {
		info, err := a.bot.GetWebhookInfo(ctx)
		if err != nil {
			return err
		}
		logJSON("info", eventStepOK, "getWebhookInfo succeeded", fields{"step": "get_webhook_info", "has_url": info.URL != "", "pending_update_count": info.PendingUpdateCount})
		return nil
	})
	a.step(ctx, summary, "get_chat", func(ctx context.Context) error {
		chat, err := a.bot.GetChat(ctx, aigram.GetChatParams{ChatID: chatID})
		if err != nil {
			return err
		}
		logJSON("info", eventStepOK, "getChat succeeded", fields{"step": "get_chat", "chat_id": maskInt64(chat.ID), "chat_type": chat.Type})
		return nil
	})
	a.step(ctx, summary, "get_chat_member_count", func(ctx context.Context) error {
		count, err := a.bot.GetChatMemberCount(ctx, aigram.GetChatMemberCountParams{ChatID: chatID})
		if err != nil {
			return err
		}
		logJSON("info", eventStepOK, "getChatMemberCount succeeded", fields{"step": "get_chat_member_count", "count": count})
		return nil
	})
	a.step(ctx, summary, "get_profile_metadata", func(ctx context.Context) error {
		name, err := a.bot.GetMyName(ctx, aigram.GetMyNameParams{})
		if err != nil {
			return fmt.Errorf("get my name: %w", err)
		}
		description, err := a.bot.GetMyDescription(ctx, aigram.GetMyDescriptionParams{})
		if err != nil {
			return fmt.Errorf("get my description: %w", err)
		}
		shortDescription, err := a.bot.GetMyShortDescription(ctx, aigram.GetMyShortDescriptionParams{})
		if err != nil {
			return fmt.Errorf("get my short description: %w", err)
		}
		logJSON("info", eventStepOK, "profile metadata read succeeded", fields{
			"step":                  "get_profile_metadata",
			"has_name":              name.Name != "",
			"description_len":       len(description.Description),
			"short_description_len": len(shortDescription.ShortDescription),
		})
		return nil
	})
	logJSON("info", "max_api_readonly_finished", "read-only suite finished", fields{"source": source, "passed": summary.passed, "failed": summary.failed})
	if summary.failed > 0 {
		return fmt.Errorf("read-only suite had %d failed steps", summary.failed)
	}
	return nil
}

func (a *app) runMediaSuite(ctx context.Context, chatID aigram.ChatID, source string) error {
	summary := newSummary()
	logJSON("info", "max_api_media_started", "media suite started", fields{"source": source, "chat_id": chatIDForLog(chatID)})
	var documentFileID string

	a.step(ctx, summary, "send_document_upload", func(ctx context.Context) error {
		message, err := a.bot.SendDocument(ctx, aigram.SendDocumentParams{
			ChatID: chatID,
			Document: aigram.FileUpload(aigram.UploadFile{
				Name:        "aigram-smoke.txt",
				Reader:      strings.NewReader("ai-gram generated smoke document\n"),
				ContentType: "text/plain; charset=utf-8",
			}),
			Caption: "Generated document upload.",
		})
		if err != nil {
			return err
		}
		if message.Document != nil {
			documentFileID = message.Document.FileID
		}
		logJSON("info", eventStepOK, "sendDocument upload succeeded", fields{"step": "send_document_upload", "message_id": message.MessageID, "has_file_id": documentFileID != ""})
		return nil
	})
	if documentFileID != "" {
		a.step(ctx, summary, "get_file", func(ctx context.Context) error {
			file, err := a.bot.GetFile(ctx, aigram.GetFileParams{FileID: documentFileID})
			if err != nil {
				return err
			}
			logJSON("info", eventStepOK, "getFile succeeded", fields{"step": "get_file", "file_id": redactID(file.FileID), "file_size": file.FileSize, "has_path": file.FilePath != ""})
			return nil
		})
	}
	a.step(ctx, summary, "send_photo_upload", func(ctx context.Context) error {
		data, err := base64.StdEncoding.DecodeString(tinyPNGBase64)
		if err != nil {
			return err
		}
		message, err := a.bot.SendPhoto(ctx, aigram.SendPhotoParams{
			ChatID: chatID,
			Photo: aigram.FileUpload(aigram.UploadFile{
				Name:        "aigram-smoke.png",
				Reader:      bytes.NewReader(data),
				ContentType: "image/png",
			}),
			Caption: "Generated PNG upload.",
		})
		if err != nil {
			return err
		}
		logJSON("info", eventStepOK, "sendPhoto upload succeeded", fields{"step": "send_photo_upload", "message_id": message.MessageID, "photo_sizes": len(message.Photo)})
		return nil
	})
	a.step(ctx, summary, "send_media_group_upload", func(ctx context.Context) error {
		messages, err := a.bot.SendMediaGroup(ctx, aigram.SendMediaGroupParams{
			ChatID: chatID,
			Media: []aigram.InputMedia{
				func() aigram.InputMediaDocument {
					item := aigram.MediaDocument(aigram.FileUpload(aigram.UploadFile{Name: "aigram-group-1.txt", Reader: strings.NewReader("media group item 1\n"), ContentType: "text/plain"}))
					item.Caption = "Generated media group"
					return item
				}(),
				aigram.MediaDocument(aigram.FileUpload(aigram.UploadFile{Name: "aigram-group-2.txt", Reader: strings.NewReader("media group item 2\n"), ContentType: "text/plain"})),
			},
		})
		if err != nil {
			return err
		}
		logJSON("info", eventStepOK, "sendMediaGroup upload succeeded", fields{"step": "send_media_group_upload", "message_count": len(messages)})
		return nil
	})

	if a.config.allowDelete {
		a.step(ctx, summary, "delete_test_message", func(ctx context.Context) error {
			message, err := a.bot.SendMessage(ctx, aigram.SendMessageParams{ChatID: chatID, Text: "This disposable smoke message will be deleted."})
			if err != nil {
				return err
			}
			ok, err := a.bot.DeleteMessage(ctx, aigram.DeleteMessageParams{ChatID: chatID, MessageID: message.MessageID})
			if err != nil {
				return err
			}
			logJSON("info", eventStepOK, "deleteMessage succeeded", fields{"step": "delete_test_message", "ok": ok, "message_id": message.MessageID})
			return nil
		})
	} else {
		logJSON("info", "max_api_step_skipped", "deleteMessage skipped", fields{"step": "delete_test_message", "reason": "AIGRAM_MAX_API_ALLOW_DELETE is not enabled"})
	}

	logJSON("info", "max_api_media_finished", "media suite finished", fields{"source": source, "passed": summary.passed, "failed": summary.failed})
	if summary.failed > 0 {
		return fmt.Errorf("media suite had %d failed steps", summary.failed)
	}
	return nil
}

func (a *app) step(ctx context.Context, summary *suiteSummary, name string, fn func(context.Context) error) {
	start := time.Now()
	logJSON("info", "max_api_step_started", "smoke step started", fields{"step": name})
	err := fn(ctx)
	duration := time.Since(start).Milliseconds()
	if err != nil {
		summary.failed++
		logJSON("error", eventStepErr, "smoke step failed", fields{"step": name, "duration_ms": duration, "error": err.Error()})
		return
	}
	summary.passed++
	logJSON("info", eventStepOK, "smoke step finished", fields{"step": name, "duration_ms": duration})
}

type suiteSummary struct {
	passed int
	failed int
}

func newSummary() *suiteSummary {
	return &suiteSummary{}
}

type fields map[string]any

func logJSON(level string, event string, message string, extra fields) {
	record := fields{
		"ts":      time.Now().UTC().Format(time.RFC3339Nano),
		"level":   level,
		"event":   event,
		"message": message,
	}
	for key, value := range extra {
		record[key] = value
	}
	data, err := json.Marshal(record)
	if err != nil {
		log.Printf(`{"level":"error","event":"max_api_log_encode_error","message":"could not encode log record"}`)
		return
	}
	log.Println(string(data))
}

func logUpdate(update aigram.Update) {
	chat := update.EffectiveChat()
	user := update.EffectiveUser()
	message := update.EffectiveMessage()
	data := fields{
		"update_id":   update.UpdateID,
		"update_type": updateType(update),
	}
	if chat != nil {
		data["chat_id"] = maskInt64(chat.ID)
		data["chat_type"] = chat.Type
	}
	if user != nil {
		data["user_id"] = maskInt64(user.ID)
		data["is_bot"] = user.IsBot
	}
	if message != nil {
		data["message_id"] = message.MessageID
		data["message_kind"] = messageKind(message)
		data["text_len"] = len(message.Text)
	}
	if update.CallbackQuery != nil {
		data["callback_data"] = update.CallbackQuery.Data
	}
	logJSON("info", "max_api_update_received", "update received", data)
}

func updateType(update aigram.Update) string {
	switch {
	case update.Message != nil:
		return "message"
	case update.EditedMessage != nil:
		return "edited_message"
	case update.CallbackQuery != nil:
		return "callback_query"
	case update.InlineQuery != nil:
		return "inline_query"
	case update.ChosenInlineResult != nil:
		return "chosen_inline_result"
	case update.Poll != nil:
		return "poll"
	case update.PollAnswer != nil:
		return "poll_answer"
	case update.MyChatMember != nil:
		return "my_chat_member"
	case update.ChatMember != nil:
		return "chat_member"
	case update.ChatJoinRequest != nil:
		return "chat_join_request"
	case update.MessageReaction != nil:
		return "message_reaction"
	case update.MessageReactionCount != nil:
		return "message_reaction_count"
	case update.ChatBoost != nil:
		return "chat_boost"
	case update.RemovedChatBoost != nil:
		return "removed_chat_boost"
	case update.BusinessConnection != nil:
		return "business_connection"
	case update.BusinessMessage != nil:
		return "business_message"
	case update.EditedBusinessMessage != nil:
		return "edited_business_message"
	case update.DeletedBusinessMessages != nil:
		return "deleted_business_messages"
	case update.PurchasedPaidMedia != nil:
		return "purchased_paid_media"
	default:
		return "unknown"
	}
}

func messageKind(message *aigram.Message) string {
	switch {
	case message.Text != "":
		return "text"
	case message.Document != nil:
		return "document"
	case len(message.Photo) > 0:
		return "photo"
	case message.Contact != nil:
		return "contact"
	case message.Location != nil:
		return "location"
	case message.Venue != nil:
		return "venue"
	case message.Poll != nil:
		return "poll"
	case message.Dice != nil:
		return "dice"
	default:
		return "other"
	}
}

func intEnv(name string, fallback int) (int, error) {
	raw := strings.TrimSpace(os.Getenv(name))
	if raw == "" {
		return fallback, nil
	}
	value, err := strconv.Atoi(raw)
	if err != nil || value <= 0 {
		return 0, fmt.Errorf("%s must be a positive integer", name)
	}
	return value, nil
}

func boolEnv(name string) bool {
	switch strings.ToLower(strings.TrimSpace(os.Getenv(name))) {
	case "1", "true", "yes", "on":
		return true
	default:
		return false
	}
}

func chatIDValid(chatID aigram.ChatID) bool {
	_, err := json.Marshal(chatID)
	return err == nil
}

func chatIDForLog(chatID aigram.ChatID) string {
	data, err := json.Marshal(chatID)
	if err != nil {
		return "unset"
	}
	value := strings.Trim(string(data), `"`)
	if len(value) <= 6 {
		return "***"
	}
	return value[:3] + "***" + value[len(value)-3:]
}

func safeString(value string) string {
	if strings.TrimSpace(value) == "" {
		return "unknown"
	}
	return value
}

func redactID(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return ""
	}
	if len(value) <= 12 {
		return "***"
	}
	return value[:6] + "..." + value[len(value)-4:]
}

func maskInt64(value int64) string {
	raw := strconv.FormatInt(value, 10)
	if len(raw) <= 6 {
		return "***"
	}
	return raw[:3] + "***" + raw[len(raw)-3:]
}

const tinyPNGBase64 = "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAQAAAC1HAwCAAAAC0lEQVR42mP8/x8AAwMCAO+/p9sAAAAASUVORK5CYII="
