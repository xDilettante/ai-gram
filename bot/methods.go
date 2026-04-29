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

// GetMe returns basic information about the bot.
func (b *Bot) GetMe(ctx context.Context) (*telegram.User, error) {
	var user telegram.User
	if err := b.call(ctx, "getMe", nil, &user); err != nil {
		return nil, err
	}

	return &user, nil
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
