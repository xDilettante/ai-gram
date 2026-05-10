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

// DeleteMessageReactionParams contains supported parameters for deleteMessageReaction.
type DeleteMessageReactionParams struct {
	ChatID      ChatID `json:"chat_id"`
	MessageID   int64  `json:"message_id"`
	UserID      int64  `json:"user_id,omitempty"`
	ActorChatID int64  `json:"actor_chat_id,omitempty"`
}

// DeleteAllMessageReactionsParams contains supported parameters for deleteAllMessageReactions.
type DeleteAllMessageReactionsParams struct {
	ChatID      ChatID `json:"chat_id"`
	UserID      int64  `json:"user_id,omitempty"`
	ActorChatID int64  `json:"actor_chat_id,omitempty"`
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

// DeleteMessageReaction removes a user's or chat's reaction from a message.
func (b *Bot) DeleteMessageReaction(ctx context.Context, params DeleteMessageReactionParams) (bool, error) {
	if err := params.validate(); err != nil {
		return false, err
	}

	var result bool
	if err := b.call(ctx, "deleteMessageReaction", params, &result); err != nil {
		return false, err
	}

	return result, nil
}

// DeleteAllMessageReactions removes recent reactions added by a user or chat.
func (b *Bot) DeleteAllMessageReactions(ctx context.Context, params DeleteAllMessageReactionsParams) (bool, error) {
	if err := params.validate(); err != nil {
		return false, err
	}

	var result bool
	if err := b.call(ctx, "deleteAllMessageReactions", params, &result); err != nil {
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

func (params DeleteMessageReactionParams) validate() error {
	if !params.ChatID.valid() {
		return stderrors.New("chat_id is required")
	}
	if params.MessageID <= 0 {
		return stderrors.New("message_id must be greater than zero")
	}
	return validateReactionActor(params.UserID, params.ActorChatID)
}

func (params DeleteAllMessageReactionsParams) validate() error {
	if !params.ChatID.valid() {
		return stderrors.New("chat_id is required")
	}
	return validateReactionActor(params.UserID, params.ActorChatID)
}

func validateReactionActor(userID int64, actorChatID int64) error {
	if userID <= 0 && actorChatID == 0 {
		return stderrors.New("user_id or actor_chat_id is required")
	}
	if userID < 0 {
		return stderrors.New("user_id must be greater than zero")
	}
	if actorChatID == 0 {
		return nil
	}
	if userID > 0 {
		return stderrors.New("user_id and actor_chat_id cannot both be set")
	}
	return nil
}
