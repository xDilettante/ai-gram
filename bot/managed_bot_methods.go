package bot

import (
	"context"
	stderrors "errors"
	"fmt"

	"github.com/xDilettante/ai-gram/telegram"
)

const maxManagedBotAccessUserIDs = 10

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

// GetManagedBotAccessSettingsParams contains supported parameters for getManagedBotAccessSettings.
type GetManagedBotAccessSettingsParams struct {
	UserID int64 `json:"user_id"`
}

// SetManagedBotAccessSettingsParams contains supported parameters for setManagedBotAccessSettings.
type SetManagedBotAccessSettingsParams struct {
	UserID             int64   `json:"user_id"`
	IsAccessRestricted bool    `json:"is_access_restricted"`
	AddedUserIDs       []int64 `json:"added_user_ids,omitempty"`
}

// GetUserPersonalChatMessagesParams contains supported parameters for getUserPersonalChatMessages.
type GetUserPersonalChatMessagesParams struct {
	UserID int64 `json:"user_id"`
	Limit  int   `json:"limit"`
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

// GetManagedBotAccessSettings returns the access settings of a managed bot.
func (b *Bot) GetManagedBotAccessSettings(ctx context.Context, params GetManagedBotAccessSettingsParams) (*telegram.BotAccessSettings, error) {
	if err := params.validate(); err != nil {
		return nil, err
	}

	var result telegram.BotAccessSettings
	if err := b.call(ctx, "getManagedBotAccessSettings", params, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// SetManagedBotAccessSettings changes the access settings of a managed bot.
func (b *Bot) SetManagedBotAccessSettings(ctx context.Context, params SetManagedBotAccessSettingsParams) (bool, error) {
	if err := params.validate(); err != nil {
		return false, err
	}

	var result bool
	if err := b.call(ctx, "setManagedBotAccessSettings", params, &result); err != nil {
		return false, err
	}
	return result, nil
}

// GetUserPersonalChatMessages gets the last messages from a user's personal chat.
func (b *Bot) GetUserPersonalChatMessages(ctx context.Context, params GetUserPersonalChatMessagesParams) ([]telegram.Message, error) {
	if err := params.validate(); err != nil {
		return nil, err
	}

	var result []telegram.Message
	if err := b.call(ctx, "getUserPersonalChatMessages", params, &result); err != nil {
		return nil, err
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

func (params GetManagedBotAccessSettingsParams) validate() error {
	return validateManagedBotUserID(params.UserID)
}

func (params SetManagedBotAccessSettingsParams) validate() error {
	if err := validateManagedBotUserID(params.UserID); err != nil {
		return err
	}
	if len(params.AddedUserIDs) > maxManagedBotAccessUserIDs {
		return fmt.Errorf("added_user_ids must contain at most %d user IDs", maxManagedBotAccessUserIDs)
	}
	for index, userID := range params.AddedUserIDs {
		if userID <= 0 {
			return fmt.Errorf("added_user_ids[%d] must be positive", index)
		}
	}
	return nil
}

func (params GetUserPersonalChatMessagesParams) validate() error {
	if params.UserID <= 0 {
		return stderrors.New("user_id must be positive")
	}
	if params.Limit < 1 || params.Limit > 20 {
		return stderrors.New("limit must be between 1 and 20")
	}
	return nil
}

func validateManagedBotUserID(userID int64) error {
	if userID <= 0 {
		return stderrors.New("user_id must be positive")
	}
	return nil
}
