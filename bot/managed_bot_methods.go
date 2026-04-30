package bot

import (
	"context"
	stderrors "errors"

	"github.com/xDilettante/ai-gram/telegram"
)

// SavePreparedKeyboardButtonParams contains supported parameters for savePreparedKeyboardButton.
type SavePreparedKeyboardButtonParams struct {
	UserID int64                   `json:"user_id"`
	Button telegram.KeyboardButton `json:"button"`
}

// GetManagedBotTokenParams contains supported parameters for getManagedBotToken.
type GetManagedBotTokenParams struct {
	UserID int64 `json:"user_id"`
}

// ReplaceManagedBotTokenParams contains supported parameters for replaceManagedBotToken.
type ReplaceManagedBotTokenParams struct {
	UserID int64 `json:"user_id"`
}

// SavePreparedKeyboardButton stores a request keyboard button for use by a Mini App user.
func (b *Bot) SavePreparedKeyboardButton(ctx context.Context, params SavePreparedKeyboardButtonParams) (*telegram.PreparedKeyboardButton, error) {
	if err := params.validate(); err != nil {
		return nil, err
	}

	var result telegram.PreparedKeyboardButton
	if err := b.call(ctx, "savePreparedKeyboardButton", params, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetManagedBotToken returns the token of a managed bot.
func (b *Bot) GetManagedBotToken(ctx context.Context, params GetManagedBotTokenParams) (string, error) {
	if err := params.validate(); err != nil {
		return "", err
	}

	var result string
	if err := b.call(ctx, "getManagedBotToken", params, &result); err != nil {
		return "", err
	}
	return result, nil
}

// ReplaceManagedBotToken revokes a managed bot token and returns the replacement token.
func (b *Bot) ReplaceManagedBotToken(ctx context.Context, params ReplaceManagedBotTokenParams) (string, error) {
	if err := params.validate(); err != nil {
		return "", err
	}

	var result string
	if err := b.call(ctx, "replaceManagedBotToken", params, &result); err != nil {
		return "", err
	}
	return result, nil
}

func (params SavePreparedKeyboardButtonParams) validate() error {
	if params.UserID <= 0 {
		return stderrors.New("user_id must be positive")
	}
	if err := telegram.ValidatePreparedKeyboardButton(params.Button); err != nil {
		return err
	}
	return nil
}

func (params GetManagedBotTokenParams) validate() error {
	return validateManagedBotUserID(params.UserID)
}

func (params ReplaceManagedBotTokenParams) validate() error {
	return validateManagedBotUserID(params.UserID)
}

func validateManagedBotUserID(userID int64) error {
	if userID <= 0 {
		return stderrors.New("user_id must be positive")
	}
	return nil
}
