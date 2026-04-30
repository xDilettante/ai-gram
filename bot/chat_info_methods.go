package bot

import (
	"context"
	stderrors "errors"

	"ai-gram/telegram"
)

// GetChatParams contains supported parameters for getChat.
type GetChatParams struct {
	ChatID ChatID `json:"chat_id"`
}

// GetChatMemberParams contains supported parameters for getChatMember.
type GetChatMemberParams struct {
	ChatID ChatID `json:"chat_id"`
	UserID int64  `json:"user_id"`
}

// GetChatAdministratorsParams contains supported parameters for getChatAdministrators.
type GetChatAdministratorsParams struct {
	ChatID ChatID `json:"chat_id"`
}

// GetChatMemberCountParams contains supported parameters for getChatMemberCount.
type GetChatMemberCountParams struct {
	ChatID ChatID `json:"chat_id"`
}

// GetChat returns information about a chat.
func (b *Bot) GetChat(ctx context.Context, params GetChatParams) (*telegram.Chat, error) {
	if err := params.validate(); err != nil {
		return nil, err
	}

	var chat telegram.Chat
	if err := b.call(ctx, "getChat", params, &chat); err != nil {
		return nil, err
	}

	return &chat, nil
}

// GetChatMember returns information about a chat member.
func (b *Bot) GetChatMember(ctx context.Context, params GetChatMemberParams) (*telegram.ChatMember, error) {
	if err := params.validate(); err != nil {
		return nil, err
	}

	var member telegram.ChatMember
	if err := b.call(ctx, "getChatMember", params, &member); err != nil {
		return nil, err
	}

	return &member, nil
}

// GetChatAdministrators returns administrators of a chat.
func (b *Bot) GetChatAdministrators(ctx context.Context, params GetChatAdministratorsParams) ([]telegram.ChatMember, error) {
	if err := params.validate(); err != nil {
		return nil, err
	}

	var administrators []telegram.ChatMember
	if err := b.call(ctx, "getChatAdministrators", params, &administrators); err != nil {
		return nil, err
	}

	return administrators, nil
}

// GetChatMemberCount returns the number of members in a chat.
func (b *Bot) GetChatMemberCount(ctx context.Context, params GetChatMemberCountParams) (int, error) {
	if err := params.validate(); err != nil {
		return 0, err
	}

	var count int
	if err := b.call(ctx, "getChatMemberCount", params, &count); err != nil {
		return 0, err
	}

	return count, nil
}

func (params GetChatParams) validate() error {
	if !params.ChatID.valid() {
		return stderrors.New("chat_id is required")
	}
	return nil
}

func (params GetChatMemberParams) validate() error {
	if !params.ChatID.valid() {
		return stderrors.New("chat_id is required")
	}
	if params.UserID <= 0 {
		return stderrors.New("user_id must be greater than zero")
	}
	return nil
}

func (params GetChatAdministratorsParams) validate() error {
	if !params.ChatID.valid() {
		return stderrors.New("chat_id is required")
	}
	return nil
}

func (params GetChatMemberCountParams) validate() error {
	if !params.ChatID.valid() {
		return stderrors.New("chat_id is required")
	}
	return nil
}
