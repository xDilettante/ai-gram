package bot

import (
	"context"
	stderrors "errors"

	"ai-gram/telegram"
)

// BanChatMemberParams contains supported parameters for banChatMember.
type BanChatMemberParams struct {
	ChatID         ChatID `json:"chat_id"`
	UserID         int64  `json:"user_id"`
	UntilDate      int64  `json:"until_date,omitempty"`
	RevokeMessages bool   `json:"revoke_messages,omitempty"`
}

// UnbanChatMemberParams contains supported parameters for unbanChatMember.
type UnbanChatMemberParams struct {
	ChatID       ChatID `json:"chat_id"`
	UserID       int64  `json:"user_id"`
	OnlyIfBanned bool   `json:"only_if_banned,omitempty"`
}

// RestrictChatMemberParams contains supported parameters for restrictChatMember.
type RestrictChatMemberParams struct {
	ChatID                        ChatID                   `json:"chat_id"`
	UserID                        int64                    `json:"user_id"`
	Permissions                   telegram.ChatPermissions `json:"permissions"`
	UseIndependentChatPermissions bool                     `json:"use_independent_chat_permissions,omitempty"`
	UntilDate                     int64                    `json:"until_date,omitempty"`
}

// BanChatMember bans a user from a chat.
func (b *Bot) BanChatMember(ctx context.Context, params BanChatMemberParams) (bool, error) {
	if err := params.validate(); err != nil {
		return false, err
	}

	var result bool
	if err := b.call(ctx, "banChatMember", params, &result); err != nil {
		return false, err
	}

	return result, nil
}

// UnbanChatMember unbans a previously banned user from a chat.
func (b *Bot) UnbanChatMember(ctx context.Context, params UnbanChatMemberParams) (bool, error) {
	if err := params.validate(); err != nil {
		return false, err
	}

	var result bool
	if err := b.call(ctx, "unbanChatMember", params, &result); err != nil {
		return false, err
	}

	return result, nil
}

// RestrictChatMember restricts a user's permissions in a chat.
func (b *Bot) RestrictChatMember(ctx context.Context, params RestrictChatMemberParams) (bool, error) {
	if err := params.validate(); err != nil {
		return false, err
	}

	var result bool
	if err := b.call(ctx, "restrictChatMember", params, &result); err != nil {
		return false, err
	}

	return result, nil
}

func (params BanChatMemberParams) validate() error {
	if !params.ChatID.valid() {
		return stderrors.New("chat_id is required")
	}
	if params.UserID <= 0 {
		return stderrors.New("user_id must be greater than zero")
	}
	if params.UntilDate < 0 {
		return stderrors.New("until_date must not be negative")
	}
	return nil
}

func (params UnbanChatMemberParams) validate() error {
	if !params.ChatID.valid() {
		return stderrors.New("chat_id is required")
	}
	if params.UserID <= 0 {
		return stderrors.New("user_id must be greater than zero")
	}
	return nil
}

func (params RestrictChatMemberParams) validate() error {
	if !params.ChatID.valid() {
		return stderrors.New("chat_id is required")
	}
	if params.UserID <= 0 {
		return stderrors.New("user_id must be greater than zero")
	}
	if params.UntilDate < 0 {
		return stderrors.New("until_date must not be negative")
	}
	return nil
}
