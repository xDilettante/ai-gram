package bot

import (
	"context"
	stderrors "errors"
	"fmt"

	"github.com/xDilettante/ai-gram/telegram"
)

// SavePreparedInlineMessageParams contains supported parameters for savePreparedInlineMessage.
type SavePreparedInlineMessageParams struct {
	UserID            int64             `json:"user_id"`
	Result            InlineQueryResult `json:"result"`
	AllowUserChats    bool              `json:"allow_user_chats,omitempty"`
	AllowBotChats     bool              `json:"allow_bot_chats,omitempty"`
	AllowGroupChats   bool              `json:"allow_group_chats,omitempty"`
	AllowChannelChats bool              `json:"allow_channel_chats,omitempty"`
}

// SavePreparedInlineMessage stores an inline message for use by a Mini App user.
func (b *Bot) SavePreparedInlineMessage(ctx context.Context, params SavePreparedInlineMessageParams) (*telegram.PreparedInlineMessage, error) {
	if err := params.validate(); err != nil {
		return nil, err
	}

	var result telegram.PreparedInlineMessage
	if err := b.call(ctx, "savePreparedInlineMessage", params, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (params SavePreparedInlineMessageParams) validate() error {
	if params.UserID <= 0 {
		return stderrors.New("user_id must be positive")
	}
	if err := validateInlineQueryResult(params.Result); err != nil {
		return fmt.Errorf("result is invalid: %w", err)
	}
	return nil
}
