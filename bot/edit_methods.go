package bot

import (
	"bytes"
	"context"
	"encoding/json"
	stderrors "errors"
	"strings"

	"ai-gram/telegram"
)

// EditMessageResult contains the polymorphic result returned by Telegram edit message methods.
//
// Telegram returns an edited Message when the target is a regular chat message, or true when the
// target is an inline message. OK is true for both successful shapes.
type EditMessageResult struct {
	Message *telegram.Message
	OK      bool
}

// IsMessage reports whether the edit result contains an edited Message object.
func (r EditMessageResult) IsMessage() bool {
	return r.Message != nil
}

// IsOK reports whether Telegram returned a successful edit result.
func (r EditMessageResult) IsOK() bool {
	return r.OK
}

// UnmarshalJSON decodes Telegram edit results, which can be either a Message object or a boolean.
func (r *EditMessageResult) UnmarshalJSON(data []byte) error {
	trimmed := bytes.TrimSpace(data)
	if bytes.Equal(trimmed, []byte("true")) {
		r.OK = true
		r.Message = nil
		return nil
	}
	if bytes.Equal(trimmed, []byte("false")) {
		r.OK = false
		r.Message = nil
		return nil
	}
	if len(trimmed) > 0 && trimmed[0] == '{' {
		var message telegram.Message
		if err := json.Unmarshal(trimmed, &message); err != nil {
			return err
		}
		r.OK = true
		r.Message = &message
		return nil
	}

	return stderrors.New("edit message result must be a message object or boolean")
}

// EditMessageTarget identifies a message that can be edited, either by chat/message ID or inline message ID.
type EditMessageTarget struct {
	ChatID          ChatID `json:"chat_id,omitempty"`
	MessageID       int64  `json:"message_id,omitempty"`
	InlineMessageID string `json:"inline_message_id,omitempty"`
}

// EditTargetChat creates an edit target for a regular chat message.
func EditTargetChat(chatID ChatID, messageID int64) EditMessageTarget {
	return EditMessageTarget{ChatID: chatID, MessageID: messageID}
}

// EditTargetInline creates an edit target for an inline message.
func EditTargetInline(inlineMessageID string) EditMessageTarget {
	return EditMessageTarget{InlineMessageID: inlineMessageID}
}

// EditMessageTextParams contains supported parameters for editMessageText.
type EditMessageTextParams struct {
	Target              EditMessageTarget              `json:"-"`
	Text                string                         `json:"text"`
	ParseMode           string                         `json:"parse_mode,omitempty"`
	Entities            []telegram.MessageEntity       `json:"entities,omitempty"`
	LinkPreviewDisabled bool                           `json:"disable_web_page_preview,omitempty"`
	ReplyMarkup         *telegram.InlineKeyboardMarkup `json:"reply_markup,omitempty"`
}

// EditMessageReplyMarkupParams contains supported parameters for editMessageReplyMarkup.
type EditMessageReplyMarkupParams struct {
	Target      EditMessageTarget              `json:"-"`
	ReplyMarkup *telegram.InlineKeyboardMarkup `json:"reply_markup,omitempty"`
}

// EditMessageText edits text and inline keyboard markup for an existing chat or inline message.
func (b *Bot) EditMessageText(ctx context.Context, params EditMessageTextParams) (*EditMessageResult, error) {
	if err := params.validate(); err != nil {
		return nil, err
	}

	var result EditMessageResult
	if err := b.call(ctx, "editMessageText", params.payload(), &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// EditMessageReplyMarkup edits or removes inline keyboard markup for an existing chat or inline message.
func (b *Bot) EditMessageReplyMarkup(ctx context.Context, params EditMessageReplyMarkupParams) (*EditMessageResult, error) {
	if err := params.validate(); err != nil {
		return nil, err
	}

	var result EditMessageResult
	if err := b.call(ctx, "editMessageReplyMarkup", params.payload(), &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (target EditMessageTarget) validate() error {
	inlineSet := strings.TrimSpace(target.InlineMessageID) != ""
	chatIDSet := target.ChatID.valid()
	messageIDSet := target.MessageID != 0
	chatMode := chatIDSet || messageIDSet

	if inlineSet && chatMode {
		return stderrors.New("edit target must use either chat_id/message_id or inline_message_id, not both")
	}
	if !inlineSet && !chatMode {
		return stderrors.New("edit target is required")
	}
	if inlineSet {
		return nil
	}
	if !chatIDSet {
		return stderrors.New("chat_id is required for chat edit target")
	}
	if target.MessageID <= 0 {
		return stderrors.New("message_id must be greater than zero for chat edit target")
	}

	return nil
}

func (params EditMessageTextParams) validate() error {
	if err := params.Target.validate(); err != nil {
		return err
	}
	if params.Text == "" {
		return stderrors.New("text is required")
	}
	if err := validateEntityFormatting(params.ParseMode, params.Entities); err != nil {
		return err
	}
	if params.ReplyMarkup != nil {
		if err := telegram.ValidateReplyMarkup(*params.ReplyMarkup); err != nil {
			return err
		}
	}

	return nil
}

func (params EditMessageReplyMarkupParams) validate() error {
	if err := params.Target.validate(); err != nil {
		return err
	}
	if params.ReplyMarkup != nil {
		if err := telegram.ValidateReplyMarkup(*params.ReplyMarkup); err != nil {
			return err
		}
	}

	return nil
}

type editMessageTextPayload struct {
	ChatID              *ChatID                        `json:"chat_id,omitempty"`
	MessageID           int64                          `json:"message_id,omitempty"`
	InlineMessageID     string                         `json:"inline_message_id,omitempty"`
	Text                string                         `json:"text"`
	ParseMode           string                         `json:"parse_mode,omitempty"`
	Entities            []telegram.MessageEntity       `json:"entities,omitempty"`
	LinkPreviewDisabled bool                           `json:"disable_web_page_preview,omitempty"`
	ReplyMarkup         *telegram.InlineKeyboardMarkup `json:"reply_markup,omitempty"`
}

type editMessageReplyMarkupPayload struct {
	ChatID          *ChatID                        `json:"chat_id,omitempty"`
	MessageID       int64                          `json:"message_id,omitempty"`
	InlineMessageID string                         `json:"inline_message_id,omitempty"`
	ReplyMarkup     *telegram.InlineKeyboardMarkup `json:"reply_markup,omitempty"`
}

func (params EditMessageTextParams) payload() editMessageTextPayload {
	payload := editMessageTextPayload{
		InlineMessageID:     params.Target.InlineMessageID,
		Text:                params.Text,
		ParseMode:           params.ParseMode,
		Entities:            params.Entities,
		LinkPreviewDisabled: params.LinkPreviewDisabled,
		ReplyMarkup:         params.ReplyMarkup,
	}
	if params.Target.ChatID.valid() {
		chatID := params.Target.ChatID
		payload.ChatID = &chatID
		payload.MessageID = params.Target.MessageID
	}

	return payload
}

func (params EditMessageReplyMarkupParams) payload() editMessageReplyMarkupPayload {
	payload := editMessageReplyMarkupPayload{
		InlineMessageID: params.Target.InlineMessageID,
		ReplyMarkup:     params.ReplyMarkup,
	}
	if params.Target.ChatID.valid() {
		chatID := params.Target.ChatID
		payload.ChatID = &chatID
		payload.MessageID = params.Target.MessageID
	}

	return payload
}

func validateEntityFormatting(parseMode string, entities []telegram.MessageEntity) error {
	if parseMode != "" && len(entities) > 0 {
		return stderrors.New("parse_mode and entities cannot be used together")
	}

	return nil
}
