package bot

import (
	"context"
	stderrors "errors"
	"fmt"

	"github.com/xDilettante/ai-gram/telegram"
)

const maxBatchMessageIDs = 100

// ForwardMessagesParams contains supported parameters for forwardMessages.
type ForwardMessagesParams struct {
	ChatID              ChatID  `json:"chat_id"`
	MessageThreadID     int64   `json:"message_thread_id,omitempty"`
	FromChatID          ChatID  `json:"from_chat_id"`
	MessageIDs          []int64 `json:"message_ids"`
	DisableNotification bool    `json:"disable_notification,omitempty"`
	ProtectContent      bool    `json:"protect_content,omitempty"`
}

// CopyMessagesParams contains supported parameters for copyMessages.
type CopyMessagesParams struct {
	ChatID              ChatID  `json:"chat_id"`
	MessageThreadID     int64   `json:"message_thread_id,omitempty"`
	FromChatID          ChatID  `json:"from_chat_id"`
	MessageIDs          []int64 `json:"message_ids"`
	DisableNotification bool    `json:"disable_notification,omitempty"`
	ProtectContent      bool    `json:"protect_content,omitempty"`
	RemoveCaption       bool    `json:"remove_caption,omitempty"`
}

// DeleteMessagesParams contains supported parameters for deleteMessages.
type DeleteMessagesParams struct {
	ChatID     ChatID  `json:"chat_id"`
	MessageIDs []int64 `json:"message_ids"`
}

// ForwardMessages forwards multiple messages from one chat to another.
func (b *Bot) ForwardMessages(ctx context.Context, params ForwardMessagesParams) ([]telegram.MessageID, error) {
	if err := params.validate(); err != nil {
		return nil, err
	}

	var messageIDs []telegram.MessageID
	if err := b.call(ctx, "forwardMessages", params, &messageIDs); err != nil {
		return nil, err
	}

	return messageIDs, nil
}

// CopyMessages copies multiple messages without forward headers and returns new message IDs.
func (b *Bot) CopyMessages(ctx context.Context, params CopyMessagesParams) ([]telegram.MessageID, error) {
	if err := params.validate(); err != nil {
		return nil, err
	}

	var messageIDs []telegram.MessageID
	if err := b.call(ctx, "copyMessages", params, &messageIDs); err != nil {
		return nil, err
	}

	return messageIDs, nil
}

// DeleteMessages deletes multiple messages from a chat.
func (b *Bot) DeleteMessages(ctx context.Context, params DeleteMessagesParams) (bool, error) {
	if err := params.validate(); err != nil {
		return false, err
	}

	var result bool
	if err := b.call(ctx, "deleteMessages", params, &result); err != nil {
		return false, err
	}

	return result, nil
}

func (params ForwardMessagesParams) validate() error {
	if !params.ChatID.valid() {
		return stderrors.New("chat_id is required")
	}
	if !params.FromChatID.valid() {
		return stderrors.New("from_chat_id is required")
	}
	if err := validateMessageThreadID(params.MessageThreadID); err != nil {
		return err
	}
	return validateBatchMessageIDs(params.MessageIDs)
}

func (params CopyMessagesParams) validate() error {
	if !params.ChatID.valid() {
		return stderrors.New("chat_id is required")
	}
	if !params.FromChatID.valid() {
		return stderrors.New("from_chat_id is required")
	}
	if err := validateMessageThreadID(params.MessageThreadID); err != nil {
		return err
	}
	return validateBatchMessageIDs(params.MessageIDs)
}

func (params DeleteMessagesParams) validate() error {
	if !params.ChatID.valid() {
		return stderrors.New("chat_id is required")
	}
	return validateBatchMessageIDs(params.MessageIDs)
}

func validateBatchMessageIDs(messageIDs []int64) error {
	if len(messageIDs) == 0 {
		return stderrors.New("message_ids must contain at least one message_id")
	}
	if len(messageIDs) > maxBatchMessageIDs {
		return fmt.Errorf("message_ids must contain at most %d message_ids", maxBatchMessageIDs)
	}
	for index, messageID := range messageIDs {
		if messageID <= 0 {
			return fmt.Errorf("message_ids[%d] must be greater than zero", index)
		}
	}
	return nil
}
