package bot

import (
	"context"
	stderrors "errors"
	"strings"

	"github.com/xDilettante/ai-gram/telegram"
)

// CreateForumTopicParams contains supported parameters for createForumTopic.
type CreateForumTopicParams struct {
	ChatID            ChatID `json:"chat_id"`
	Name              string `json:"name"`
	IconColor         int    `json:"icon_color,omitempty"`
	IconCustomEmojiID string `json:"icon_custom_emoji_id,omitempty"`
}

// EditForumTopicParams contains supported parameters for editForumTopic.
type EditForumTopicParams struct {
	ChatID            ChatID `json:"chat_id"`
	MessageThreadID   int64  `json:"message_thread_id"`
	Name              string `json:"name,omitempty"`
	IconCustomEmojiID string `json:"icon_custom_emoji_id,omitempty"`
}

// CloseForumTopicParams contains supported parameters for closeForumTopic.
type CloseForumTopicParams struct {
	ChatID          ChatID `json:"chat_id"`
	MessageThreadID int64  `json:"message_thread_id"`
}

// ReopenForumTopicParams contains supported parameters for reopenForumTopic.
type ReopenForumTopicParams struct {
	ChatID          ChatID `json:"chat_id"`
	MessageThreadID int64  `json:"message_thread_id"`
}

// DeleteForumTopicParams contains supported parameters for deleteForumTopic.
type DeleteForumTopicParams struct {
	ChatID          ChatID `json:"chat_id"`
	MessageThreadID int64  `json:"message_thread_id"`
}

// UnpinAllForumTopicMessagesParams contains supported parameters for unpinAllForumTopicMessages.
type UnpinAllForumTopicMessagesParams struct {
	ChatID          ChatID `json:"chat_id"`
	MessageThreadID int64  `json:"message_thread_id"`
}

// EditGeneralForumTopicParams contains supported parameters for editGeneralForumTopic.
type EditGeneralForumTopicParams struct {
	ChatID ChatID `json:"chat_id"`
	Name   string `json:"name"`
}

// CloseGeneralForumTopicParams contains supported parameters for closeGeneralForumTopic.
type CloseGeneralForumTopicParams struct {
	ChatID ChatID `json:"chat_id"`
}

// ReopenGeneralForumTopicParams contains supported parameters for reopenGeneralForumTopic.
type ReopenGeneralForumTopicParams struct {
	ChatID ChatID `json:"chat_id"`
}

// HideGeneralForumTopicParams contains supported parameters for hideGeneralForumTopic.
type HideGeneralForumTopicParams struct {
	ChatID ChatID `json:"chat_id"`
}

// UnhideGeneralForumTopicParams contains supported parameters for unhideGeneralForumTopic.
type UnhideGeneralForumTopicParams struct {
	ChatID ChatID `json:"chat_id"`
}

// UnpinAllGeneralForumTopicMessagesParams contains supported parameters for unpinAllGeneralForumTopicMessages.
type UnpinAllGeneralForumTopicMessagesParams struct {
	ChatID ChatID `json:"chat_id"`
}

// CreateForumTopic creates a new topic in a forum supergroup.
func (b *Bot) CreateForumTopic(ctx context.Context, params CreateForumTopicParams) (*telegram.ForumTopic, error) {
	if err := params.validate(); err != nil {
		return nil, err
	}

	var topic telegram.ForumTopic
	if err := b.call(ctx, "createForumTopic", params, &topic); err != nil {
		return nil, err
	}

	return &topic, nil
}

// EditForumTopic edits a topic in a forum supergroup.
func (b *Bot) EditForumTopic(ctx context.Context, params EditForumTopicParams) (bool, error) {
	if err := params.validate(); err != nil {
		return false, err
	}

	var result bool
	if err := b.call(ctx, "editForumTopic", params, &result); err != nil {
		return false, err
	}

	return result, nil
}

// CloseForumTopic closes a topic in a forum supergroup.
func (b *Bot) CloseForumTopic(ctx context.Context, params CloseForumTopicParams) (bool, error) {
	if err := params.validate(); err != nil {
		return false, err
	}

	var result bool
	if err := b.call(ctx, "closeForumTopic", params, &result); err != nil {
		return false, err
	}

	return result, nil
}

// ReopenForumTopic reopens a topic in a forum supergroup.
func (b *Bot) ReopenForumTopic(ctx context.Context, params ReopenForumTopicParams) (bool, error) {
	if err := params.validate(); err != nil {
		return false, err
	}

	var result bool
	if err := b.call(ctx, "reopenForumTopic", params, &result); err != nil {
		return false, err
	}

	return result, nil
}

// DeleteForumTopic deletes a topic in a forum supergroup.
func (b *Bot) DeleteForumTopic(ctx context.Context, params DeleteForumTopicParams) (bool, error) {
	if err := params.validate(); err != nil {
		return false, err
	}

	var result bool
	if err := b.call(ctx, "deleteForumTopic", params, &result); err != nil {
		return false, err
	}

	return result, nil
}

// UnpinAllForumTopicMessages clears pinned messages in a forum topic.
func (b *Bot) UnpinAllForumTopicMessages(ctx context.Context, params UnpinAllForumTopicMessagesParams) (bool, error) {
	if err := params.validate(); err != nil {
		return false, err
	}

	var result bool
	if err := b.call(ctx, "unpinAllForumTopicMessages", params, &result); err != nil {
		return false, err
	}

	return result, nil
}

// EditGeneralForumTopic edits the General topic in a forum supergroup.
func (b *Bot) EditGeneralForumTopic(ctx context.Context, params EditGeneralForumTopicParams) (bool, error) {
	if err := params.validate(); err != nil {
		return false, err
	}

	var result bool
	if err := b.call(ctx, "editGeneralForumTopic", params, &result); err != nil {
		return false, err
	}

	return result, nil
}

// CloseGeneralForumTopic closes the General topic in a forum supergroup.
func (b *Bot) CloseGeneralForumTopic(ctx context.Context, params CloseGeneralForumTopicParams) (bool, error) {
	if err := params.validate(); err != nil {
		return false, err
	}

	var result bool
	if err := b.call(ctx, "closeGeneralForumTopic", params, &result); err != nil {
		return false, err
	}

	return result, nil
}

// ReopenGeneralForumTopic reopens the General topic in a forum supergroup.
func (b *Bot) ReopenGeneralForumTopic(ctx context.Context, params ReopenGeneralForumTopicParams) (bool, error) {
	if err := params.validate(); err != nil {
		return false, err
	}

	var result bool
	if err := b.call(ctx, "reopenGeneralForumTopic", params, &result); err != nil {
		return false, err
	}

	return result, nil
}

// HideGeneralForumTopic hides the General topic in a forum supergroup.
func (b *Bot) HideGeneralForumTopic(ctx context.Context, params HideGeneralForumTopicParams) (bool, error) {
	if err := params.validate(); err != nil {
		return false, err
	}

	var result bool
	if err := b.call(ctx, "hideGeneralForumTopic", params, &result); err != nil {
		return false, err
	}

	return result, nil
}

// UnhideGeneralForumTopic unhides the General topic in a forum supergroup.
func (b *Bot) UnhideGeneralForumTopic(ctx context.Context, params UnhideGeneralForumTopicParams) (bool, error) {
	if err := params.validate(); err != nil {
		return false, err
	}

	var result bool
	if err := b.call(ctx, "unhideGeneralForumTopic", params, &result); err != nil {
		return false, err
	}

	return result, nil
}

// UnpinAllGeneralForumTopicMessages clears pinned messages in the General forum topic.
func (b *Bot) UnpinAllGeneralForumTopicMessages(ctx context.Context, params UnpinAllGeneralForumTopicMessagesParams) (bool, error) {
	if err := params.validate(); err != nil {
		return false, err
	}

	var result bool
	if err := b.call(ctx, "unpinAllGeneralForumTopicMessages", params, &result); err != nil {
		return false, err
	}

	return result, nil
}

func (params CreateForumTopicParams) validate() error {
	if !params.ChatID.valid() {
		return stderrors.New("chat_id is required")
	}
	if strings.TrimSpace(params.Name) == "" {
		return stderrors.New("name is required")
	}
	if params.IconColor < 0 {
		return stderrors.New("icon_color must be non-negative")
	}
	return nil
}

func (params EditForumTopicParams) validate() error {
	if err := validateRequiredChatID(params.ChatID); err != nil {
		return err
	}
	return validateRequiredMessageThreadID(params.MessageThreadID)
}

func (params CloseForumTopicParams) validate() error {
	return validateForumTopicTarget(params.ChatID, params.MessageThreadID)
}

func (params ReopenForumTopicParams) validate() error {
	return validateForumTopicTarget(params.ChatID, params.MessageThreadID)
}

func (params DeleteForumTopicParams) validate() error {
	return validateForumTopicTarget(params.ChatID, params.MessageThreadID)
}

func (params UnpinAllForumTopicMessagesParams) validate() error {
	return validateForumTopicTarget(params.ChatID, params.MessageThreadID)
}

func (params EditGeneralForumTopicParams) validate() error {
	if err := validateRequiredChatID(params.ChatID); err != nil {
		return err
	}
	if strings.TrimSpace(params.Name) == "" {
		return stderrors.New("name is required")
	}
	return nil
}

func (params CloseGeneralForumTopicParams) validate() error {
	return validateRequiredChatID(params.ChatID)
}

func (params ReopenGeneralForumTopicParams) validate() error {
	return validateRequiredChatID(params.ChatID)
}

func (params HideGeneralForumTopicParams) validate() error {
	return validateRequiredChatID(params.ChatID)
}

func (params UnhideGeneralForumTopicParams) validate() error {
	return validateRequiredChatID(params.ChatID)
}

func (params UnpinAllGeneralForumTopicMessagesParams) validate() error {
	return validateRequiredChatID(params.ChatID)
}

func validateForumTopicTarget(chatID ChatID, messageThreadID int64) error {
	if err := validateRequiredChatID(chatID); err != nil {
		return err
	}
	return validateRequiredMessageThreadID(messageThreadID)
}

func validateRequiredMessageThreadID(messageThreadID int64) error {
	if messageThreadID <= 0 {
		return stderrors.New("message_thread_id must be greater than zero")
	}
	return nil
}
