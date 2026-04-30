package bot

import (
	"context"
	stderrors "errors"

	"github.com/xDilettante/ai-gram/telegram"
)

// SetMessageReactionParams contains supported parameters for setMessageReaction.
type SetMessageReactionParams struct {
	ChatID    ChatID                  `json:"chat_id"`
	MessageID int64                   `json:"message_id"`
	Reaction  []telegram.ReactionType `json:"reaction,omitempty"`
	IsBig     bool                    `json:"is_big,omitempty"`
}

// SetMessageReaction changes reactions on a message.
func (b *Bot) SetMessageReaction(ctx context.Context, params SetMessageReactionParams) (bool, error) {
	if err := params.validate(); err != nil {
		return false, err
	}

	var result bool
	if err := b.call(ctx, "setMessageReaction", params, &result); err != nil {
		return false, err
	}

	return result, nil
}

func (params SetMessageReactionParams) validate() error {
	if !params.ChatID.valid() {
		return stderrors.New("chat_id is required")
	}
	if params.MessageID <= 0 {
		return stderrors.New("message_id must be greater than zero")
	}
	if err := telegram.ValidateReactionTypes(params.Reaction); err != nil {
		return err
	}

	return nil
}
