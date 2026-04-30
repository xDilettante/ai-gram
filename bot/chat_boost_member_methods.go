package bot

import (
	"context"
	stderrors "errors"

	"github.com/xDilettante/ai-gram/telegram"
)

// GetUserChatBoostsParams contains supported parameters for getUserChatBoosts.
type GetUserChatBoostsParams struct {
	ChatID ChatID `json:"chat_id"`
	UserID int64  `json:"user_id"`
}

// SetChatMemberTagParams contains supported parameters for setChatMemberTag.
type SetChatMemberTagParams struct {
	ChatID ChatID `json:"chat_id"`
	UserID int64  `json:"user_id"`
	Tag    string `json:"tag"`
}

// BanChatSenderChatParams contains supported parameters for banChatSenderChat.
type BanChatSenderChatParams struct {
	ChatID       ChatID `json:"chat_id"`
	SenderChatID int64  `json:"sender_chat_id"`
}

// UnbanChatSenderChatParams contains supported parameters for unbanChatSenderChat.
type UnbanChatSenderChatParams struct {
	ChatID       ChatID `json:"chat_id"`
	SenderChatID int64  `json:"sender_chat_id"`
}

// GetUserChatBoosts gets the list of boosts added to a chat by a user.
func (b *Bot) GetUserChatBoosts(ctx context.Context, params GetUserChatBoostsParams) (*telegram.UserChatBoosts, error) {
	if err := params.validate(); err != nil {
		return nil, err
	}

	var result telegram.UserChatBoosts
	if err := b.call(ctx, "getUserChatBoosts", params, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// SetChatMemberTag sets a tag for a regular member in a group or supergroup.
func (b *Bot) SetChatMemberTag(ctx context.Context, params SetChatMemberTagParams) (bool, error) {
	if err := params.validate(); err != nil {
		return false, err
	}

	var result bool
	if err := b.call(ctx, "setChatMemberTag", params, &result); err != nil {
		return false, err
	}
	return result, nil
}

// BanChatSenderChat bans a sender channel in a supergroup or channel.
func (b *Bot) BanChatSenderChat(ctx context.Context, params BanChatSenderChatParams) (bool, error) {
	if err := params.validate(); err != nil {
		return false, err
	}

	var result bool
	if err := b.call(ctx, "banChatSenderChat", params, &result); err != nil {
		return false, err
	}
	return result, nil
}

// UnbanChatSenderChat unbans a previously banned sender channel.
func (b *Bot) UnbanChatSenderChat(ctx context.Context, params UnbanChatSenderChatParams) (bool, error) {
	if err := params.validate(); err != nil {
		return false, err
	}

	var result bool
	if err := b.call(ctx, "unbanChatSenderChat", params, &result); err != nil {
		return false, err
	}
	return result, nil
}

func (params GetUserChatBoostsParams) validate() error {
	if !params.ChatID.valid() {
		return stderrors.New("chat_id is required")
	}
	if params.UserID <= 0 {
		return stderrors.New("user_id must be greater than zero")
	}
	return nil
}

func (params SetChatMemberTagParams) validate() error {
	if !params.ChatID.valid() {
		return stderrors.New("chat_id is required")
	}
	if params.UserID <= 0 {
		return stderrors.New("user_id must be greater than zero")
	}
	return nil
}

func (params BanChatSenderChatParams) validate() error {
	if !params.ChatID.valid() {
		return stderrors.New("chat_id is required")
	}
	if params.SenderChatID == 0 {
		return stderrors.New("sender_chat_id is required")
	}
	return nil
}

func (params UnbanChatSenderChatParams) validate() error {
	if !params.ChatID.valid() {
		return stderrors.New("chat_id is required")
	}
	if params.SenderChatID == 0 {
		return stderrors.New("sender_chat_id is required")
	}
	return nil
}
