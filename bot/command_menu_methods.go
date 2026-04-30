package bot

import (
	"context"
	"encoding/json"
	stderrors "errors"
	"strings"

	"github.com/xDilettante/ai-gram/telegram"
)

// SetMyCommandsParams contains supported parameters for setMyCommands.
type SetMyCommandsParams struct {
	Commands     []telegram.BotCommand    `json:"commands"`
	Scope        telegram.BotCommandScope `json:"scope,omitempty"`
	LanguageCode string                   `json:"language_code,omitempty"`
}

// DeleteMyCommandsParams contains supported parameters for deleteMyCommands.
type DeleteMyCommandsParams struct {
	Scope        telegram.BotCommandScope `json:"scope,omitempty"`
	LanguageCode string                   `json:"language_code,omitempty"`
}

// GetMyCommandsParams contains supported parameters for getMyCommands.
type GetMyCommandsParams struct {
	Scope        telegram.BotCommandScope `json:"scope,omitempty"`
	LanguageCode string                   `json:"language_code,omitempty"`
}

// SetChatMenuButtonParams contains supported parameters for setChatMenuButton.
type SetChatMenuButtonParams struct {
	ChatID     ChatID              `json:"chat_id,omitempty"`
	MenuButton telegram.MenuButton `json:"menu_button,omitempty"`
}

// GetChatMenuButtonParams contains supported parameters for getChatMenuButton.
type GetChatMenuButtonParams struct {
	ChatID ChatID `json:"chat_id,omitempty"`
}

// SetMyDefaultAdministratorRightsParams contains supported parameters for setMyDefaultAdministratorRights.
type SetMyDefaultAdministratorRightsParams struct {
	Rights      *telegram.ChatAdministratorRights `json:"rights,omitempty"`
	ForChannels bool                              `json:"for_channels,omitempty"`
}

// SetMyCommands changes the list of the bot's commands.
func (b *Bot) SetMyCommands(ctx context.Context, params SetMyCommandsParams) (bool, error) {
	if err := params.validate(); err != nil {
		return false, err
	}

	var result bool
	if err := b.call(ctx, "setMyCommands", params, &result); err != nil {
		return false, err
	}

	return result, nil
}

// DeleteMyCommands deletes the bot's commands for the given scope and language.
func (b *Bot) DeleteMyCommands(ctx context.Context, params DeleteMyCommandsParams) (bool, error) {
	if err := params.validate(); err != nil {
		return false, err
	}

	var result bool
	if err := b.call(ctx, "deleteMyCommands", params, &result); err != nil {
		return false, err
	}

	return result, nil
}

// GetMyCommands gets the bot's commands for the given scope and language.
func (b *Bot) GetMyCommands(ctx context.Context, params GetMyCommandsParams) ([]telegram.BotCommand, error) {
	if err := params.validate(); err != nil {
		return nil, err
	}

	var result []telegram.BotCommand
	if err := b.call(ctx, "getMyCommands", params, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// SetChatMenuButton changes the bot's menu button in a private chat or the default menu button.
func (b *Bot) SetChatMenuButton(ctx context.Context, params SetChatMenuButtonParams) (bool, error) {
	if err := params.validate(); err != nil {
		return false, err
	}

	var result bool
	if err := b.call(ctx, "setChatMenuButton", params.payload(), &result); err != nil {
		return false, err
	}

	return result, nil
}

// GetChatMenuButton gets the bot's current menu button in a private chat or the default menu button.
func (b *Bot) GetChatMenuButton(ctx context.Context, params GetChatMenuButtonParams) (telegram.MenuButton, error) {
	var result menuButtonResult
	if err := b.call(ctx, "getChatMenuButton", params.payload(), &result); err != nil {
		return nil, err
	}

	return result.Button, nil
}

// SetMyDefaultAdministratorRights changes the default administrator rights requested by the bot.
func (b *Bot) SetMyDefaultAdministratorRights(ctx context.Context, params SetMyDefaultAdministratorRightsParams) (bool, error) {
	var result bool
	if err := b.call(ctx, "setMyDefaultAdministratorRights", params, &result); err != nil {
		return false, err
	}

	return result, nil
}

func (params SetMyCommandsParams) validate() error {
	if len(params.Commands) == 0 {
		return stderrors.New("commands must not be empty")
	}
	for _, command := range params.Commands {
		if strings.TrimSpace(command.Command) == "" {
			return stderrors.New("command is required")
		}
		if strings.TrimSpace(command.Description) == "" {
			return stderrors.New("command description is required")
		}
	}
	if err := telegram.ValidateBotCommandScope(params.Scope); err != nil {
		return err
	}
	return nil
}

func (params DeleteMyCommandsParams) validate() error {
	return telegram.ValidateBotCommandScope(params.Scope)
}

func (params GetMyCommandsParams) validate() error {
	return telegram.ValidateBotCommandScope(params.Scope)
}

func (params SetChatMenuButtonParams) validate() error {
	return telegram.ValidateMenuButton(params.MenuButton)
}

func (params SetChatMenuButtonParams) payload() map[string]any {
	payload := make(map[string]any)
	if params.ChatID.valid() {
		payload["chat_id"] = params.ChatID
	}
	if params.MenuButton != nil {
		payload["menu_button"] = params.MenuButton
	}
	return payload
}

func (params GetChatMenuButtonParams) payload() map[string]any {
	payload := make(map[string]any)
	if params.ChatID.valid() {
		payload["chat_id"] = params.ChatID
	}
	return payload
}

type menuButtonResult struct {
	Button telegram.MenuButton
}

func (result *menuButtonResult) UnmarshalJSON(data []byte) error {
	button, err := telegram.UnmarshalMenuButton(data)
	if err != nil {
		return err
	}
	result.Button = button
	return nil
}

var _ json.Unmarshaler = (*menuButtonResult)(nil)
