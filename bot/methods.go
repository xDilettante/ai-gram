package bot

import (
	"context"
	stderrors "errors"

	"ai-gram/telegram"
)

// SendMessageParams contains supported parameters for sendMessage.
type SendMessageParams struct {
	ChatID              ChatID `json:"chat_id"`
	Text                string `json:"text"`
	ParseMode           string `json:"parse_mode,omitempty"`
	DisableNotification bool   `json:"disable_notification,omitempty"`
}

// GetUpdatesParams contains supported parameters for getUpdates.
type GetUpdatesParams struct {
	Offset         int64    `json:"offset,omitempty"`
	Limit          int      `json:"limit,omitempty"`
	Timeout        int      `json:"timeout,omitempty"`
	AllowedUpdates []string `json:"allowed_updates,omitempty"`
}

// GetMe returns basic information about the bot.
func (b *Bot) GetMe(ctx context.Context) (*telegram.User, error) {
	var user telegram.User
	if err := b.call(ctx, "getMe", nil, &user); err != nil {
		return nil, err
	}

	return &user, nil
}

// GetUpdates fetches updates with one getUpdates API call.
func (b *Bot) GetUpdates(ctx context.Context, params GetUpdatesParams) ([]telegram.Update, error) {
	if err := params.validate(); err != nil {
		return nil, err
	}

	var updates []telegram.Update
	if err := b.call(ctx, "getUpdates", params, &updates); err != nil {
		return nil, err
	}

	return updates, nil
}

func (params GetUpdatesParams) validate() error {
	if params.Limit < 0 {
		return stderrors.New("limit must not be negative")
	}
	if params.Limit > 100 {
		return stderrors.New("limit must be between 1 and 100")
	}
	if params.Timeout < 0 {
		return stderrors.New("timeout must not be negative")
	}

	return nil
}

// SendMessage sends a text message.
func (b *Bot) SendMessage(ctx context.Context, params SendMessageParams) (*telegram.Message, error) {
	if !params.ChatID.valid() {
		return nil, stderrors.New("chat_id is required")
	}
	if params.Text == "" {
		return nil, stderrors.New("text is required")
	}

	var message telegram.Message
	if err := b.call(ctx, "sendMessage", params, &message); err != nil {
		return nil, err
	}

	return &message, nil
}
