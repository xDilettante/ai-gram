package bot

import (
	"context"
	stderrors "errors"

	"github.com/xDilettante/ai-gram/telegram"
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
	ChatID     ChatID `json:"chat_id"`
	ReturnBots bool   `json:"return_bots,omitempty"`
}

// GetChatMemberCountParams contains supported parameters for getChatMemberCount.
type GetChatMemberCountParams struct {
	ChatID ChatID `json:"chat_id"`
}

// GetChat returns full information about a chat.
func (b *Bot) GetChat(ctx context.Context, params GetChatParams) (*telegram.ChatFullInfo, error) {
	if err := params.validate(); err != nil {
		return nil, err
	}

	var chat telegram.ChatFullInfo
	if err := b.call(ctx, "getChat", params, &chat); err != nil {
		return nil, err
	}

	return &chat, nil
}

// GetChatFullInfo returns full information about a chat.
func (b *Bot) GetChatFullInfo(ctx context.Context, params GetChatParams) (*telegram.ChatFullInfo, error) {
	return b.GetChat(ctx, params)
}

// GetChatMember returns information about a chat member.
func (b *Bot) GetChatMember(ctx context.Context, params GetChatMemberParams) (telegram.ChatMember, error) {
	if err := params.validate(); err != nil {
		return nil, err
	}

	var member telegram.ChatMemberResult
	if err := b.call(ctx, "getChatMember", params, &member); err != nil {
		return nil, err
	}

	return member.ChatMember, nil
}

// GetChatAdministrators returns administrators of a chat.
func (b *Bot) GetChatAdministrators(ctx context.Context, params GetChatAdministratorsParams) ([]telegram.ChatMember, error) {
	if err := params.validate(); err != nil {
		return nil, err
	}

	var administrators []telegram.ChatMemberResult
	if err := b.call(ctx, "getChatAdministrators", params, &administrators); err != nil {
		return nil, err
	}

	members := make([]telegram.ChatMember, 0, len(administrators))
	for _, administrator := range administrators {
		if administrator.ChatMember != nil {
			members = append(members, administrator.ChatMember)
		}
	}

	return members, nil
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
