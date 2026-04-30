package bot

import (
	"context"
	stderrors "errors"

	"github.com/xDilettante/ai-gram/telegram"
)

// ForwardMessageParams contains supported parameters for forwardMessage.
type ForwardMessageParams struct {
	ChatID              ChatID `json:"chat_id"`
	FromChatID          ChatID `json:"from_chat_id"`
	MessageID           int64  `json:"message_id"`
	MessageThreadID     int64  `json:"message_thread_id,omitempty"`
	DisableNotification bool   `json:"disable_notification,omitempty"`
	ProtectContent      bool   `json:"protect_content,omitempty"`
}

// CopyMessageParams contains supported parameters for copyMessage.
type CopyMessageParams struct {
	ChatID              ChatID                    `json:"chat_id"`
	FromChatID          ChatID                    `json:"from_chat_id"`
	MessageID           int64                     `json:"message_id"`
	MessageThreadID     int64                     `json:"message_thread_id,omitempty"`
	Caption             string                    `json:"caption,omitempty"`
	ParseMode           string                    `json:"parse_mode,omitempty"`
	CaptionEntities     []telegram.MessageEntity  `json:"caption_entities,omitempty"`
	DisableNotification bool                      `json:"disable_notification,omitempty"`
	ProtectContent      bool                      `json:"protect_content,omitempty"`
	ReplyParameters     *telegram.ReplyParameters `json:"reply_parameters,omitempty"`
	ReplyMarkup         telegram.ReplyMarkup      `json:"reply_markup,omitempty"`
}

// ForwardMessage forwards a message from one chat to another.
func (b *Bot) ForwardMessage(ctx context.Context, params ForwardMessageParams) (*telegram.Message, error) {
	if err := params.validate(); err != nil {
		return nil, err
	}

	var message telegram.Message
	if err := b.call(ctx, "forwardMessage", params, &message); err != nil {
		return nil, err
	}

	return &message, nil
}

// CopyMessage copies a message without a forward header and returns the new message ID.
func (b *Bot) CopyMessage(ctx context.Context, params CopyMessageParams) (*telegram.MessageID, error) {
	if err := params.validate(); err != nil {
		return nil, err
	}

	var messageID telegram.MessageID
	if err := b.call(ctx, "copyMessage", params, &messageID); err != nil {
		return nil, err
	}

	return &messageID, nil
}

func (params ForwardMessageParams) validate() error {
	if !params.ChatID.valid() {
		return stderrors.New("chat_id is required")
	}
	if !params.FromChatID.valid() {
		return stderrors.New("from_chat_id is required")
	}
	if params.MessageID <= 0 {
		return stderrors.New("message_id must be greater than zero")
	}
	if err := validateMessageThreadID(params.MessageThreadID); err != nil {
		return err
	}

	return nil
}

func (params CopyMessageParams) validate() error {
	if !params.ChatID.valid() {
		return stderrors.New("chat_id is required")
	}
	if !params.FromChatID.valid() {
		return stderrors.New("from_chat_id is required")
	}
	if params.MessageID <= 0 {
		return stderrors.New("message_id must be greater than zero")
	}
	if err := validateMessageThreadID(params.MessageThreadID); err != nil {
		return err
	}
	if err := validateCaptionFormatting(params.ParseMode, params.CaptionEntities); err != nil {
		return err
	}
	if err := validateReplyParameters(params.ReplyParameters); err != nil {
		return err
	}
	if err := telegram.ValidateReplyMarkup(params.ReplyMarkup); err != nil {
		return err
	}

	return nil
}
