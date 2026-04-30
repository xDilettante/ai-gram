package bot

import (
	"context"
	stderrors "errors"
	"strings"

	"github.com/xDilettante/ai-gram/telegram"
)

// SendGameParams contains supported parameters for sendGame.
type SendGameParams struct {
	BusinessConnectionID string                         `json:"business_connection_id,omitempty"`
	ChatID               ChatID                         `json:"chat_id"`
	MessageThreadID      int64                          `json:"message_thread_id,omitempty"`
	GameShortName        string                         `json:"game_short_name"`
	DisableNotification  bool                           `json:"disable_notification,omitempty"`
	ProtectContent       bool                           `json:"protect_content,omitempty"`
	AllowPaidBroadcast   bool                           `json:"allow_paid_broadcast,omitempty"`
	MessageEffectID      string                         `json:"message_effect_id,omitempty"`
	ReplyParameters      *telegram.ReplyParameters      `json:"reply_parameters,omitempty"`
	ReplyMarkup          *telegram.InlineKeyboardMarkup `json:"reply_markup,omitempty"`
}

// SetGameScoreParams contains supported parameters for setGameScore.
type SetGameScoreParams struct {
	UserID             int64  `json:"user_id"`
	Score              int    `json:"score"`
	Force              bool   `json:"force,omitempty"`
	DisableEditMessage bool   `json:"disable_edit_message,omitempty"`
	ChatID             ChatID `json:"chat_id,omitempty"`
	MessageID          int64  `json:"message_id,omitempty"`
	InlineMessageID    string `json:"inline_message_id,omitempty"`
}

// GetGameHighScoresParams contains supported parameters for getGameHighScores.
type GetGameHighScoresParams struct {
	UserID          int64  `json:"user_id"`
	ChatID          ChatID `json:"chat_id,omitempty"`
	MessageID       int64  `json:"message_id,omitempty"`
	InlineMessageID string `json:"inline_message_id,omitempty"`
}

// SendGame sends a game message.
func (b *Bot) SendGame(ctx context.Context, params SendGameParams) (*telegram.Message, error) {
	if err := params.validate(); err != nil {
		return nil, err
	}

	var message telegram.Message
	if err := b.call(ctx, "sendGame", params, &message); err != nil {
		return nil, err
	}

	return &message, nil
}

// SetGameScore sets a user's score in a game message.
func (b *Bot) SetGameScore(ctx context.Context, params SetGameScoreParams) (*EditMessageResult, error) {
	if err := params.validate(); err != nil {
		return nil, err
	}

	var result EditMessageResult
	if err := b.call(ctx, "setGameScore", params.payload(), &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetGameHighScores returns high score rows around a user for a game message.
func (b *Bot) GetGameHighScores(ctx context.Context, params GetGameHighScoresParams) ([]telegram.GameHighScore, error) {
	if err := params.validate(); err != nil {
		return nil, err
	}

	var scores []telegram.GameHighScore
	if err := b.call(ctx, "getGameHighScores", params.payload(), &scores); err != nil {
		return nil, err
	}

	return scores, nil
}

func (params SendGameParams) validate() error {
	if !params.ChatID.valid() {
		return stderrors.New("chat_id is required")
	}
	if err := validateMessageThreadID(params.MessageThreadID); err != nil {
		return err
	}
	if strings.TrimSpace(params.GameShortName) == "" {
		return stderrors.New("game_short_name is required")
	}
	if err := validateReplyParameters(params.ReplyParameters); err != nil {
		return err
	}
	if params.ReplyMarkup != nil {
		if err := telegram.ValidateReplyMarkup(*params.ReplyMarkup); err != nil {
			return err
		}
		if len(params.ReplyMarkup.InlineKeyboard) > 0 && len(params.ReplyMarkup.InlineKeyboard[0]) > 0 && params.ReplyMarkup.InlineKeyboard[0][0].CallbackGame == nil {
			return stderrors.New("reply_markup first button must use callback_game")
		}
	}
	return nil
}

func (params SetGameScoreParams) validate() error {
	if params.UserID <= 0 {
		return stderrors.New("user_id must be greater than zero")
	}
	if params.Score < 0 {
		return stderrors.New("score must not be negative")
	}
	return validateGameMessageTarget(params.ChatID, params.MessageID, params.InlineMessageID)
}

func (params GetGameHighScoresParams) validate() error {
	if params.UserID <= 0 {
		return stderrors.New("user_id must be greater than zero")
	}
	return validateGameMessageTarget(params.ChatID, params.MessageID, params.InlineMessageID)
}

type setGameScorePayload struct {
	UserID             int64   `json:"user_id"`
	Score              int     `json:"score"`
	Force              bool    `json:"force,omitempty"`
	DisableEditMessage bool    `json:"disable_edit_message,omitempty"`
	ChatID             *ChatID `json:"chat_id,omitempty"`
	MessageID          int64   `json:"message_id,omitempty"`
	InlineMessageID    string  `json:"inline_message_id,omitempty"`
}

type getGameHighScoresPayload struct {
	UserID          int64   `json:"user_id"`
	ChatID          *ChatID `json:"chat_id,omitempty"`
	MessageID       int64   `json:"message_id,omitempty"`
	InlineMessageID string  `json:"inline_message_id,omitempty"`
}

func (params SetGameScoreParams) payload() setGameScorePayload {
	chatID, messageID, inlineMessageID := gameTargetPayloadValues(params.ChatID, params.MessageID, params.InlineMessageID)
	return setGameScorePayload{
		UserID:             params.UserID,
		Score:              params.Score,
		Force:              params.Force,
		DisableEditMessage: params.DisableEditMessage,
		ChatID:             chatID,
		MessageID:          messageID,
		InlineMessageID:    inlineMessageID,
	}
}

func (params GetGameHighScoresParams) payload() getGameHighScoresPayload {
	chatID, messageID, inlineMessageID := gameTargetPayloadValues(params.ChatID, params.MessageID, params.InlineMessageID)
	return getGameHighScoresPayload{
		UserID:          params.UserID,
		ChatID:          chatID,
		MessageID:       messageID,
		InlineMessageID: inlineMessageID,
	}
}

func validateGameMessageTarget(chatID ChatID, messageID int64, inlineMessageID string) error {
	inlineSet := strings.TrimSpace(inlineMessageID) != ""
	chatIDSet := chatID.valid()
	messageIDSet := messageID != 0
	chatMode := chatIDSet || messageIDSet

	if inlineSet && chatMode {
		return stderrors.New("game target must use either chat_id/message_id or inline_message_id, not both")
	}
	if !inlineSet && !chatMode {
		return stderrors.New("game target is required")
	}
	if inlineSet {
		return nil
	}
	if !chatIDSet {
		return stderrors.New("chat_id is required for game target")
	}
	if messageID <= 0 {
		return stderrors.New("message_id must be greater than zero for game target")
	}
	return nil
}

func gameTargetPayloadValues(chatID ChatID, messageID int64, inlineMessageID string) (*ChatID, int64, string) {
	if chatID.valid() {
		value := chatID
		return &value, messageID, ""
	}
	return nil, 0, inlineMessageID
}
