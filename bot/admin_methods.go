package bot

import (
	"context"
	stderrors "errors"
	"strings"

	"github.com/xDilettante/ai-gram/telegram"
)

// PromoteChatMemberParams contains supported parameters for promoteChatMember.
type PromoteChatMemberParams struct {
	ChatID                  ChatID `json:"chat_id"`
	UserID                  int64  `json:"user_id"`
	IsAnonymous             bool   `json:"is_anonymous,omitempty"`
	CanManageChat           bool   `json:"can_manage_chat,omitempty"`
	CanDeleteMessages       bool   `json:"can_delete_messages,omitempty"`
	CanManageVideoChats     bool   `json:"can_manage_video_chats,omitempty"`
	CanRestrictMembers      bool   `json:"can_restrict_members,omitempty"`
	CanPromoteMembers       bool   `json:"can_promote_members,omitempty"`
	CanChangeInfo           bool   `json:"can_change_info,omitempty"`
	CanInviteUsers          bool   `json:"can_invite_users,omitempty"`
	CanPostStories          bool   `json:"can_post_stories,omitempty"`
	CanEditStories          bool   `json:"can_edit_stories,omitempty"`
	CanDeleteStories        bool   `json:"can_delete_stories,omitempty"`
	CanPostMessages         bool   `json:"can_post_messages,omitempty"`
	CanEditMessages         bool   `json:"can_edit_messages,omitempty"`
	CanPinMessages          bool   `json:"can_pin_messages,omitempty"`
	CanManageTopics         bool   `json:"can_manage_topics,omitempty"`
	CanManageDirectMessages bool   `json:"can_manage_direct_messages,omitempty"`
	CanManageTags           bool   `json:"can_manage_tags,omitempty"`
}

// SetChatAdministratorCustomTitleParams contains supported parameters for setChatAdministratorCustomTitle.
type SetChatAdministratorCustomTitleParams struct {
	ChatID      ChatID `json:"chat_id"`
	UserID      int64  `json:"user_id"`
	CustomTitle string `json:"custom_title"`
}

// SetChatPermissionsParams contains supported parameters for setChatPermissions.
type SetChatPermissionsParams struct {
	ChatID                        ChatID                   `json:"chat_id"`
	Permissions                   telegram.ChatPermissions `json:"permissions"`
	UseIndependentChatPermissions bool                     `json:"use_independent_chat_permissions,omitempty"`
}

// PromoteChatMember promotes or demotes a user in a supergroup or channel.
func (b *Bot) PromoteChatMember(ctx context.Context, params PromoteChatMemberParams) (bool, error) {
	if err := params.validate(); err != nil {
		return false, err
	}

	var result bool
	if err := b.call(ctx, "promoteChatMember", params, &result); err != nil {
		return false, err
	}

	return result, nil
}

// SetChatAdministratorCustomTitle sets a custom title for an administrator promoted by the bot.
func (b *Bot) SetChatAdministratorCustomTitle(ctx context.Context, params SetChatAdministratorCustomTitleParams) (bool, error) {
	if err := params.validate(); err != nil {
		return false, err
	}

	var result bool
	if err := b.call(ctx, "setChatAdministratorCustomTitle", params, &result); err != nil {
		return false, err
	}

	return result, nil
}

// SetChatPermissions sets default permissions for all chat members.
func (b *Bot) SetChatPermissions(ctx context.Context, params SetChatPermissionsParams) (bool, error) {
	if err := params.validate(); err != nil {
		return false, err
	}

	var result bool
	if err := b.call(ctx, "setChatPermissions", params, &result); err != nil {
		return false, err
	}

	return result, nil
}

func (params PromoteChatMemberParams) validate() error {
	if !params.ChatID.valid() {
		return stderrors.New("chat_id is required")
	}
	if params.UserID <= 0 {
		return stderrors.New("user_id must be greater than zero")
	}
	return nil
}

func (params SetChatAdministratorCustomTitleParams) validate() error {
	if !params.ChatID.valid() {
		return stderrors.New("chat_id is required")
	}
	if params.UserID <= 0 {
		return stderrors.New("user_id must be greater than zero")
	}
	if strings.TrimSpace(params.CustomTitle) == "" {
		return stderrors.New("custom_title is required")
	}
	return nil
}

func (params SetChatPermissionsParams) validate() error {
	if !params.ChatID.valid() {
		return stderrors.New("chat_id is required")
	}
	return nil
}
