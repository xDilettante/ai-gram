package bot

import (
	"context"
	stderrors "errors"
)

// DeleteMessageParams contains supported parameters for deleteMessage.
type DeleteMessageParams struct {
	ChatID    ChatID `json:"chat_id"`
	MessageID int64  `json:"message_id"`
}

// DeleteMessage deletes a message from a chat.
func (b *Bot) DeleteMessage(ctx context.Context, params DeleteMessageParams) (bool, error) {
	if err := params.validate(); err != nil {
		return false, err
	}

	var result bool
	if err := b.call(ctx, "deleteMessage", params, &result); err != nil {
		return false, err
	}

	return result, nil
}

func (params DeleteMessageParams) validate() error {
	if !params.ChatID.valid() {
		return stderrors.New("chat_id is required")
	}
	if params.MessageID <= 0 {
		return stderrors.New("message_id must be greater than zero")
	}

	return nil
}
