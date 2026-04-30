package bot

import (
	"context"
	stderrors "errors"
)

// SetUserEmojiStatusParams contains supported parameters for setUserEmojiStatus.
type SetUserEmojiStatusParams struct {
	UserID                    int64  `json:"user_id"`
	EmojiStatusCustomEmojiID  string `json:"emoji_status_custom_emoji_id,omitempty"`
	EmojiStatusExpirationDate int64  `json:"emoji_status_expiration_date,omitempty"`
}

// VerifyUserParams contains supported parameters for verifyUser.
type VerifyUserParams struct {
	UserID            int64  `json:"user_id"`
	CustomDescription string `json:"custom_description,omitempty"`
}

// VerifyChatParams contains supported parameters for verifyChat.
type VerifyChatParams struct {
	ChatID            ChatID `json:"chat_id"`
	CustomDescription string `json:"custom_description,omitempty"`
}

// RemoveUserVerificationParams contains supported parameters for removeUserVerification.
type RemoveUserVerificationParams struct {
	UserID int64 `json:"user_id"`
}

// RemoveChatVerificationParams contains supported parameters for removeChatVerification.
type RemoveChatVerificationParams struct {
	ChatID ChatID `json:"chat_id"`
}

// SetUserEmojiStatus changes the emoji status for a user that granted access to the bot.
func (b *Bot) SetUserEmojiStatus(ctx context.Context, params SetUserEmojiStatusParams) (bool, error) {
	if err := params.validate(); err != nil {
		return false, err
	}

	var result bool
	if err := b.call(ctx, "setUserEmojiStatus", params, &result); err != nil {
		return false, err
	}
	return result, nil
}

// VerifyUser verifies a user on behalf of the organization represented by the bot.
func (b *Bot) VerifyUser(ctx context.Context, params VerifyUserParams) (bool, error) {
	if err := params.validate(); err != nil {
		return false, err
	}

	var result bool
	if err := b.call(ctx, "verifyUser", params, &result); err != nil {
		return false, err
	}
	return result, nil
}

// VerifyChat verifies a chat on behalf of the organization represented by the bot.
func (b *Bot) VerifyChat(ctx context.Context, params VerifyChatParams) (bool, error) {
	if err := params.validate(); err != nil {
		return false, err
	}

	var result bool
	if err := b.call(ctx, "verifyChat", params, &result); err != nil {
		return false, err
	}
	return result, nil
}

// RemoveUserVerification removes verification from a user verified by the bot's organization.
func (b *Bot) RemoveUserVerification(ctx context.Context, params RemoveUserVerificationParams) (bool, error) {
	if err := params.validate(); err != nil {
		return false, err
	}

	var result bool
	if err := b.call(ctx, "removeUserVerification", params, &result); err != nil {
		return false, err
	}
	return result, nil
}

// RemoveChatVerification removes verification from a chat verified by the bot's organization.
func (b *Bot) RemoveChatVerification(ctx context.Context, params RemoveChatVerificationParams) (bool, error) {
	if err := params.validate(); err != nil {
		return false, err
	}

	var result bool
	if err := b.call(ctx, "removeChatVerification", params, &result); err != nil {
		return false, err
	}
	return result, nil
}

func (params SetUserEmojiStatusParams) validate() error {
	if params.UserID <= 0 {
		return stderrors.New("user_id must be greater than zero")
	}
	if params.EmojiStatusExpirationDate < 0 {
		return stderrors.New("emoji_status_expiration_date must not be negative")
	}
	return nil
}

func (params VerifyUserParams) validate() error {
	if params.UserID <= 0 {
		return stderrors.New("user_id must be greater than zero")
	}
	return nil
}

func (params VerifyChatParams) validate() error {
	if !params.ChatID.valid() {
		return stderrors.New("chat_id is required")
	}
	return nil
}

func (params RemoveUserVerificationParams) validate() error {
	if params.UserID <= 0 {
		return stderrors.New("user_id must be greater than zero")
	}
	return nil
}

func (params RemoveChatVerificationParams) validate() error {
	if !params.ChatID.valid() {
		return stderrors.New("chat_id is required")
	}
	return nil
}
