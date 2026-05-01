package bot

import (
	"context"
	stderrors "errors"
	"strings"

	"github.com/xDilettante/ai-gram/telegram"
)

// ExportChatInviteLinkParams contains supported parameters for exportChatInviteLink.
type ExportChatInviteLinkParams struct {
	ChatID ChatID `json:"chat_id"`
}

// CreateChatInviteLinkParams contains supported parameters for createChatInviteLink.
type CreateChatInviteLinkParams struct {
	ChatID             ChatID `json:"chat_id"`
	Name               string `json:"name,omitempty"`
	ExpireDate         int64  `json:"expire_date,omitempty"`
	MemberLimit        int    `json:"member_limit,omitempty"`
	CreatesJoinRequest bool   `json:"creates_join_request,omitempty"`
}

// CreateChatSubscriptionInviteLinkParams contains supported parameters for createChatSubscriptionInviteLink.
type CreateChatSubscriptionInviteLinkParams struct {
	ChatID             ChatID `json:"chat_id"`
	Name               string `json:"name,omitempty"`
	SubscriptionPeriod int    `json:"subscription_period"`
	SubscriptionPrice  int    `json:"subscription_price"`
}

// EditChatInviteLinkParams contains supported parameters for editChatInviteLink.
type EditChatInviteLinkParams struct {
	ChatID             ChatID `json:"chat_id"`
	InviteLink         string `json:"invite_link"`
	Name               string `json:"name,omitempty"`
	ExpireDate         int64  `json:"expire_date,omitempty"`
	MemberLimit        int    `json:"member_limit,omitempty"`
	CreatesJoinRequest bool   `json:"creates_join_request,omitempty"`
}

// EditChatSubscriptionInviteLinkParams contains supported parameters for editChatSubscriptionInviteLink.
type EditChatSubscriptionInviteLinkParams struct {
	ChatID     ChatID `json:"chat_id"`
	InviteLink string `json:"invite_link"`
	Name       string `json:"name,omitempty"`
}

// RevokeChatInviteLinkParams contains supported parameters for revokeChatInviteLink.
type RevokeChatInviteLinkParams struct {
	ChatID     ChatID `json:"chat_id"`
	InviteLink string `json:"invite_link"`
}

// ExportChatInviteLink exports a new primary invite link for a chat.
func (b *Bot) ExportChatInviteLink(ctx context.Context, params ExportChatInviteLinkParams) (string, error) {
	if err := params.validate(); err != nil {
		return "", err
	}

	var inviteLink string
	if err := b.call(ctx, "exportChatInviteLink", params, &inviteLink); err != nil {
		return "", err
	}

	return inviteLink, nil
}

// CreateChatInviteLink creates an additional invite link for a chat.
func (b *Bot) CreateChatInviteLink(ctx context.Context, params CreateChatInviteLinkParams) (*telegram.ChatInviteLink, error) {
	if err := params.validate(); err != nil {
		return nil, err
	}

	var inviteLink telegram.ChatInviteLink
	if err := b.call(ctx, "createChatInviteLink", params, &inviteLink); err != nil {
		return nil, err
	}

	return &inviteLink, nil
}

// CreateChatSubscriptionInviteLink creates a subscription invite link for a channel chat.
func (b *Bot) CreateChatSubscriptionInviteLink(ctx context.Context, params CreateChatSubscriptionInviteLinkParams) (*telegram.ChatInviteLink, error) {
	if err := params.validate(); err != nil {
		return nil, err
	}

	var inviteLink telegram.ChatInviteLink
	if err := b.call(ctx, "createChatSubscriptionInviteLink", params, &inviteLink); err != nil {
		return nil, err
	}

	return &inviteLink, nil
}

// EditChatInviteLink edits a non-primary invite link created by the bot.
func (b *Bot) EditChatInviteLink(ctx context.Context, params EditChatInviteLinkParams) (*telegram.ChatInviteLink, error) {
	if err := params.validate(); err != nil {
		return nil, err
	}

	var inviteLink telegram.ChatInviteLink
	if err := b.call(ctx, "editChatInviteLink", params, &inviteLink); err != nil {
		return nil, err
	}

	return &inviteLink, nil
}

// EditChatSubscriptionInviteLink edits a subscription invite link created by the bot.
func (b *Bot) EditChatSubscriptionInviteLink(ctx context.Context, params EditChatSubscriptionInviteLinkParams) (*telegram.ChatInviteLink, error) {
	if err := params.validate(); err != nil {
		return nil, err
	}

	var inviteLink telegram.ChatInviteLink
	if err := b.call(ctx, "editChatSubscriptionInviteLink", params, &inviteLink); err != nil {
		return nil, err
	}

	return &inviteLink, nil
}

// RevokeChatInviteLink revokes an invite link created by the bot.
func (b *Bot) RevokeChatInviteLink(ctx context.Context, params RevokeChatInviteLinkParams) (*telegram.ChatInviteLink, error) {
	if err := params.validate(); err != nil {
		return nil, err
	}

	var inviteLink telegram.ChatInviteLink
	if err := b.call(ctx, "revokeChatInviteLink", params, &inviteLink); err != nil {
		return nil, err
	}

	return &inviteLink, nil
}

func (params ExportChatInviteLinkParams) validate() error {
	if !params.ChatID.valid() {
		return stderrors.New("chat_id is required")
	}
	return nil
}

func (params CreateChatInviteLinkParams) validate() error {
	if !params.ChatID.valid() {
		return stderrors.New("chat_id is required")
	}
	if params.ExpireDate < 0 {
		return stderrors.New("expire_date must not be negative")
	}
	if params.MemberLimit < 0 {
		return stderrors.New("member_limit must not be negative")
	}
	return nil
}

func (params CreateChatSubscriptionInviteLinkParams) validate() error {
	if !params.ChatID.valid() {
		return stderrors.New("chat_id is required")
	}
	if params.SubscriptionPeriod <= 0 {
		return stderrors.New("subscription_period must be positive")
	}
	if params.SubscriptionPrice <= 0 {
		return stderrors.New("subscription_price must be positive")
	}
	return nil
}

func (params EditChatInviteLinkParams) validate() error {
	if !params.ChatID.valid() {
		return stderrors.New("chat_id is required")
	}
	if strings.TrimSpace(params.InviteLink) == "" {
		return stderrors.New("invite_link is required")
	}
	if params.ExpireDate < 0 {
		return stderrors.New("expire_date must not be negative")
	}
	if params.MemberLimit < 0 {
		return stderrors.New("member_limit must not be negative")
	}
	return nil
}

func (params EditChatSubscriptionInviteLinkParams) validate() error {
	if !params.ChatID.valid() {
		return stderrors.New("chat_id is required")
	}
	if strings.TrimSpace(params.InviteLink) == "" {
		return stderrors.New("invite_link is required")
	}
	return nil
}

func (params RevokeChatInviteLinkParams) validate() error {
	if !params.ChatID.valid() {
		return stderrors.New("chat_id is required")
	}
	if strings.TrimSpace(params.InviteLink) == "" {
		return stderrors.New("invite_link is required")
	}
	return nil
}
