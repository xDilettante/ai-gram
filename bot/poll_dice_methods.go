package bot

import (
	"context"
	stderrors "errors"

	"github.com/xDilettante/ai-gram/telegram"
)

// SendPollParams contains supported parameters for sendPoll.
type SendPollParams struct {
	BusinessConnectionID   string                    `json:"business_connection_id,omitempty"`
	ChatID                 ChatID                    `json:"chat_id"`
	MessageThreadID        int64                     `json:"message_thread_id,omitempty"`
	Question               string                    `json:"question"`
	Options                []string                  `json:"options"`
	IsAnonymous            *bool                     `json:"is_anonymous,omitempty"`
	Type                   string                    `json:"type,omitempty"`
	AllowsMultipleAnswers  bool                      `json:"allows_multiple_answers,omitempty"`
	AllowsRevoting         bool                      `json:"allows_revoting,omitempty"`
	ShuffleOptions         bool                      `json:"shuffle_options,omitempty"`
	AllowAddingOptions     bool                      `json:"allow_adding_options,omitempty"`
	HideResultsUntilCloses bool                      `json:"hide_results_until_closes,omitempty"`
	CorrectOptionID        *int                      `json:"correct_option_id,omitempty"`
	CorrectOptionIDs       []int                     `json:"correct_option_ids,omitempty"`
	Explanation            string                    `json:"explanation,omitempty"`
	ExplanationParseMode   string                    `json:"explanation_parse_mode,omitempty"`
	ExplanationEntities    []telegram.MessageEntity  `json:"explanation_entities,omitempty"`
	Description            string                    `json:"description,omitempty"`
	DescriptionParseMode   string                    `json:"description_parse_mode,omitempty"`
	DescriptionEntities    []telegram.MessageEntity  `json:"description_entities,omitempty"`
	OpenPeriod             int                       `json:"open_period,omitempty"`
	CloseDate              int64                     `json:"close_date,omitempty"`
	IsClosed               bool                      `json:"is_closed,omitempty"`
	DisableNotification    bool                      `json:"disable_notification,omitempty"`
	ProtectContent         bool                      `json:"protect_content,omitempty"`
	ReplyParameters        *telegram.ReplyParameters `json:"reply_parameters,omitempty"`
	ReplyMarkup            telegram.ReplyMarkup      `json:"reply_markup,omitempty"`
}

// StopPollParams contains supported parameters for stopPoll.
type StopPollParams struct {
	BusinessConnectionID string                         `json:"business_connection_id,omitempty"`
	ChatID               ChatID                         `json:"chat_id"`
	MessageID            int64                          `json:"message_id"`
	ReplyMarkup          *telegram.InlineKeyboardMarkup `json:"reply_markup,omitempty"`
}

// SendDiceParams contains supported parameters for sendDice.
type SendDiceParams struct {
	BusinessConnectionID string                    `json:"business_connection_id,omitempty"`
	ChatID               ChatID                    `json:"chat_id"`
	MessageThreadID      int64                     `json:"message_thread_id,omitempty"`
	Emoji                string                    `json:"emoji,omitempty"`
	DisableNotification  bool                      `json:"disable_notification,omitempty"`
	ProtectContent       bool                      `json:"protect_content,omitempty"`
	ReplyParameters      *telegram.ReplyParameters `json:"reply_parameters,omitempty"`
	ReplyMarkup          telegram.ReplyMarkup      `json:"reply_markup,omitempty"`
}

// SendPoll sends a native poll.
func (b *Bot) SendPoll(ctx context.Context, params SendPollParams) (*telegram.Message, error) {
	if err := params.validate(); err != nil {
		return nil, err
	}

	var message telegram.Message
	if err := b.call(ctx, "sendPoll", params, &message); err != nil {
		return nil, err
	}

	return &message, nil
}

// StopPoll stops a poll sent by the bot and returns the stopped poll.
func (b *Bot) StopPoll(ctx context.Context, params StopPollParams) (*telegram.Poll, error) {
	if err := params.validate(); err != nil {
		return nil, err
	}

	var poll telegram.Poll
	if err := b.call(ctx, "stopPoll", params, &poll); err != nil {
		return nil, err
	}

	return &poll, nil
}

// SendDice sends an animated dice emoji and returns the resulting message.
func (b *Bot) SendDice(ctx context.Context, params SendDiceParams) (*telegram.Message, error) {
	if err := params.validate(); err != nil {
		return nil, err
	}

	var message telegram.Message
	if err := b.call(ctx, "sendDice", params, &message); err != nil {
		return nil, err
	}

	return &message, nil
}

func (params SendPollParams) validate() error {
	if !params.ChatID.valid() {
		return stderrors.New("chat_id is required")
	}
	if err := validateMessageThreadID(params.MessageThreadID); err != nil {
		return err
	}
	if params.Question == "" {
		return stderrors.New("question is required")
	}
	if len(params.Options) < 2 {
		return stderrors.New("options must contain at least two items")
	}
	for _, option := range params.Options {
		if option == "" {
			return stderrors.New("options must not contain empty items")
		}
	}
	if params.CorrectOptionID != nil && len(params.CorrectOptionIDs) > 0 {
		return stderrors.New("correct_option_id and correct_option_ids cannot both be set")
	}
	if params.CorrectOptionID != nil && (*params.CorrectOptionID < 0 || *params.CorrectOptionID >= len(params.Options)) {
		return stderrors.New("correct_option_id must reference an existing option")
	}
	for _, id := range params.CorrectOptionIDs {
		if id < 0 || id >= len(params.Options) {
			return stderrors.New("correct_option_ids must reference existing options")
		}
	}
	if params.OpenPeriod < 0 {
		return stderrors.New("open_period must not be negative")
	}
	if params.CloseDate < 0 {
		return stderrors.New("close_date must not be negative")
	}
	if err := validateEntityFormatting(params.ExplanationParseMode, params.ExplanationEntities); err != nil {
		return err
	}
	if err := validateEntityFormatting(params.DescriptionParseMode, params.DescriptionEntities); err != nil {
		return err
	}
	if err := validateReplyParameters(params.ReplyParameters); err != nil {
		return err
	}
	if err := telegram.ValidateReplyMarkup(params.ReplyMarkup); err != nil {
		return err
	}

	return nil
}

func (params StopPollParams) validate() error {
	if !params.ChatID.valid() {
		return stderrors.New("chat_id is required")
	}
	if params.MessageID <= 0 {
		return stderrors.New("message_id must be greater than zero")
	}
	if params.ReplyMarkup != nil {
		if err := telegram.ValidateReplyMarkup(*params.ReplyMarkup); err != nil {
			return err
		}
	}

	return nil
}

func (params SendDiceParams) validate() error {
	if !params.ChatID.valid() {
		return stderrors.New("chat_id is required")
	}
	if err := validateMessageThreadID(params.MessageThreadID); err != nil {
		return err
	}
	if params.Emoji != "" && !validDiceEmoji(params.Emoji) {
		return stderrors.New("emoji must be one of 🎲, 🎯, 🏀, ⚽, 🎳, or 🎰")
	}
	if err := validateReplyParameters(params.ReplyParameters); err != nil {
		return err
	}
	if err := telegram.ValidateReplyMarkup(params.ReplyMarkup); err != nil {
		return err
	}

	return nil
}

func validDiceEmoji(emoji string) bool {
	switch emoji {
	case "🎲", "🎯", "🏀", "⚽", "🎳", "🎰":
		return true
	default:
		return false
	}
}
