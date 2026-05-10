package bot

import (
	"context"
	"encoding/json"
	stderrors "errors"
	"fmt"
	"reflect"
	"strconv"

	"github.com/xDilettante/ai-gram/telegram"
)

// SendPollParams contains supported parameters for sendPoll.
type SendPollParams struct {
	BusinessConnectionID   string                     `json:"business_connection_id,omitempty"`
	ChatID                 ChatID                     `json:"chat_id"`
	MessageThreadID        int64                      `json:"message_thread_id,omitempty"`
	Question               string                     `json:"question"`
	QuestionParseMode      string                     `json:"question_parse_mode,omitempty"`
	QuestionEntities       []telegram.MessageEntity   `json:"question_entities,omitempty"`
	Options                []string                   `json:"options"`
	OptionObjects          []telegram.InputPollOption `json:"-"`
	IsAnonymous            *bool                      `json:"is_anonymous,omitempty"`
	Type                   string                     `json:"type,omitempty"`
	AllowsMultipleAnswers  bool                       `json:"allows_multiple_answers,omitempty"`
	AllowsRevoting         bool                       `json:"allows_revoting,omitempty"`
	ShuffleOptions         bool                       `json:"shuffle_options,omitempty"`
	AllowAddingOptions     bool                       `json:"allow_adding_options,omitempty"`
	HideResultsUntilCloses bool                       `json:"hide_results_until_closes,omitempty"`
	MembersOnly            bool                       `json:"members_only,omitempty"`
	CountryCodes           []string                   `json:"country_codes,omitempty"`
	CorrectOptionID        *int                       `json:"correct_option_id,omitempty"`
	CorrectOptionIDs       []int                      `json:"correct_option_ids,omitempty"`
	Explanation            string                     `json:"explanation,omitempty"`
	ExplanationParseMode   string                     `json:"explanation_parse_mode,omitempty"`
	ExplanationEntities    []telegram.MessageEntity   `json:"explanation_entities,omitempty"`
	ExplanationMedia       InputPollMedia             `json:"explanation_media,omitempty"`
	Description            string                     `json:"description,omitempty"`
	DescriptionParseMode   string                     `json:"description_parse_mode,omitempty"`
	DescriptionEntities    []telegram.MessageEntity   `json:"description_entities,omitempty"`
	Media                  InputPollMedia             `json:"media,omitempty"`
	OpenPeriod             int                        `json:"open_period,omitempty"`
	CloseDate              int64                      `json:"close_date,omitempty"`
	IsClosed               bool                       `json:"is_closed,omitempty"`
	DisableNotification    bool                       `json:"disable_notification,omitempty"`
	ProtectContent         bool                       `json:"protect_content,omitempty"`
	ReplyParameters        *telegram.ReplyParameters  `json:"reply_parameters,omitempty"`
	ReplyMarkup            telegram.ReplyMarkup       `json:"reply_markup,omitempty"`
}

type sendPollPayload struct {
	BusinessConnectionID   string                    `json:"business_connection_id,omitempty"`
	ChatID                 ChatID                    `json:"chat_id"`
	MessageThreadID        int64                     `json:"message_thread_id,omitempty"`
	Question               string                    `json:"question"`
	QuestionParseMode      string                    `json:"question_parse_mode,omitempty"`
	QuestionEntities       []telegram.MessageEntity  `json:"question_entities,omitempty"`
	Options                any                       `json:"options"`
	IsAnonymous            *bool                     `json:"is_anonymous,omitempty"`
	Type                   string                    `json:"type,omitempty"`
	AllowsMultipleAnswers  bool                      `json:"allows_multiple_answers,omitempty"`
	AllowsRevoting         bool                      `json:"allows_revoting,omitempty"`
	ShuffleOptions         bool                      `json:"shuffle_options,omitempty"`
	AllowAddingOptions     bool                      `json:"allow_adding_options,omitempty"`
	HideResultsUntilCloses bool                      `json:"hide_results_until_closes,omitempty"`
	MembersOnly            bool                      `json:"members_only,omitempty"`
	CountryCodes           []string                  `json:"country_codes,omitempty"`
	CorrectOptionID        *int                      `json:"correct_option_id,omitempty"`
	CorrectOptionIDs       []int                     `json:"correct_option_ids,omitempty"`
	Explanation            string                    `json:"explanation,omitempty"`
	ExplanationParseMode   string                    `json:"explanation_parse_mode,omitempty"`
	ExplanationEntities    []telegram.MessageEntity  `json:"explanation_entities,omitempty"`
	ExplanationMedia       any                       `json:"explanation_media,omitempty"`
	Description            string                    `json:"description,omitempty"`
	DescriptionParseMode   string                    `json:"description_parse_mode,omitempty"`
	DescriptionEntities    []telegram.MessageEntity  `json:"description_entities,omitempty"`
	Media                  any                       `json:"media,omitempty"`
	OpenPeriod             int                       `json:"open_period,omitempty"`
	CloseDate              int64                     `json:"close_date,omitempty"`
	IsClosed               bool                      `json:"is_closed,omitempty"`
	DisableNotification    bool                      `json:"disable_notification,omitempty"`
	ProtectContent         bool                      `json:"protect_content,omitempty"`
	ReplyParameters        *telegram.ReplyParameters `json:"reply_parameters,omitempty"`
	ReplyMarkup            telegram.ReplyMarkup      `json:"reply_markup,omitempty"`
}

// MarshalJSON serializes either legacy string poll options or structured InputPollOption values as "options".
func (params SendPollParams) MarshalJSON() ([]byte, error) {
	payload, err := params.payload(nil)
	if err != nil {
		return nil, err
	}
	return json.Marshal(payload)
}

func (params SendPollParams) payload(files map[string]UploadFile) (sendPollPayload, error) {
	options := any(params.Options)
	if len(params.OptionObjects) > 0 {
		optionPayloads := make([]inputPollOptionPayload, 0, len(params.OptionObjects))
		for index, option := range params.OptionObjects {
			payload, err := buildInputPollOptionPayloadForSend(option, fmt.Sprintf("options[%d]", index), files)
			if err != nil {
				return sendPollPayload{}, fmt.Errorf("options[%d]: %w", index, err)
			}
			optionPayloads = append(optionPayloads, payload)
		}
		options = optionPayloads
	}
	var media any
	if params.Media != nil && !isNilPollMediaInterface(params.Media) {
		payload, err := buildInputPollMediaPayload(params.Media, "media", files)
		if err != nil {
			return sendPollPayload{}, err
		}
		media = payload
	}
	var explanationMedia any
	if params.ExplanationMedia != nil && !isNilPollMediaInterface(params.ExplanationMedia) {
		payload, err := buildInputPollMediaPayload(params.ExplanationMedia, "explanation_media", files)
		if err != nil {
			return sendPollPayload{}, err
		}
		explanationMedia = payload
	}

	return sendPollPayload{
		BusinessConnectionID:   params.BusinessConnectionID,
		ChatID:                 params.ChatID,
		MessageThreadID:        params.MessageThreadID,
		Question:               params.Question,
		QuestionParseMode:      params.QuestionParseMode,
		QuestionEntities:       params.QuestionEntities,
		Options:                options,
		IsAnonymous:            params.IsAnonymous,
		Type:                   params.Type,
		AllowsMultipleAnswers:  params.AllowsMultipleAnswers,
		AllowsRevoting:         params.AllowsRevoting,
		ShuffleOptions:         params.ShuffleOptions,
		AllowAddingOptions:     params.AllowAddingOptions,
		HideResultsUntilCloses: params.HideResultsUntilCloses,
		MembersOnly:            params.MembersOnly,
		CountryCodes:           params.CountryCodes,
		CorrectOptionID:        params.CorrectOptionID,
		CorrectOptionIDs:       params.CorrectOptionIDs,
		Explanation:            params.Explanation,
		ExplanationParseMode:   params.ExplanationParseMode,
		ExplanationEntities:    params.ExplanationEntities,
		ExplanationMedia:       explanationMedia,
		Description:            params.Description,
		DescriptionParseMode:   params.DescriptionParseMode,
		DescriptionEntities:    params.DescriptionEntities,
		Media:                  media,
		OpenPeriod:             params.OpenPeriod,
		CloseDate:              params.CloseDate,
		IsClosed:               params.IsClosed,
		DisableNotification:    params.DisableNotification,
		ProtectContent:         params.ProtectContent,
		ReplyParameters:        params.ReplyParameters,
		ReplyMarkup:            params.ReplyMarkup,
	}, nil
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
	files := make(map[string]UploadFile)
	payload, err := params.payload(files)
	if err != nil {
		return nil, err
	}

	var message telegram.Message
	if len(files) > 0 {
		fields, err := params.multipartFields(payload)
		if err != nil {
			return nil, err
		}
		if err := b.callMultipart(ctx, "sendPoll", fields, files, &message); err != nil {
			return nil, err
		}
		return &message, nil
	}

	if err := b.call(ctx, "sendPoll", payload, &message); err != nil {
		return nil, err
	}

	return &message, nil
}

func (params SendPollParams) multipartFields(payload sendPollPayload) (map[string]string, error) {
	chatIDValue, err := params.ChatID.multipartValue()
	if err != nil {
		return nil, err
	}
	fields := map[string]string{"chat_id": chatIDValue}
	stringField(fields, "business_connection_id", payload.BusinessConnectionID)
	int64Field(fields, "message_thread_id", payload.MessageThreadID)
	stringField(fields, "question", payload.Question)
	stringField(fields, "question_parse_mode", payload.QuestionParseMode)
	if err := jsonField(fields, "question_entities", payload.QuestionEntities); err != nil {
		return nil, err
	}
	if err := requiredJSONField(fields, "options", payload.Options); err != nil {
		return nil, err
	}
	if payload.IsAnonymous != nil {
		fields["is_anonymous"] = strconv.FormatBool(*payload.IsAnonymous)
	}
	stringField(fields, "type", payload.Type)
	boolField(fields, "allows_multiple_answers", payload.AllowsMultipleAnswers)
	boolField(fields, "allows_revoting", payload.AllowsRevoting)
	boolField(fields, "shuffle_options", payload.ShuffleOptions)
	boolField(fields, "allow_adding_options", payload.AllowAddingOptions)
	boolField(fields, "hide_results_until_closes", payload.HideResultsUntilCloses)
	boolField(fields, "members_only", payload.MembersOnly)
	if err := jsonField(fields, "country_codes", payload.CountryCodes); err != nil {
		return nil, err
	}
	if payload.CorrectOptionID != nil {
		fields["correct_option_id"] = strconv.Itoa(*payload.CorrectOptionID)
	}
	if err := jsonField(fields, "correct_option_ids", payload.CorrectOptionIDs); err != nil {
		return nil, err
	}
	stringField(fields, "explanation", payload.Explanation)
	stringField(fields, "explanation_parse_mode", payload.ExplanationParseMode)
	if err := jsonField(fields, "explanation_entities", payload.ExplanationEntities); err != nil {
		return nil, err
	}
	if err := jsonField(fields, "explanation_media", payload.ExplanationMedia); err != nil {
		return nil, err
	}
	stringField(fields, "description", payload.Description)
	stringField(fields, "description_parse_mode", payload.DescriptionParseMode)
	if err := jsonField(fields, "description_entities", payload.DescriptionEntities); err != nil {
		return nil, err
	}
	if err := jsonField(fields, "media", payload.Media); err != nil {
		return nil, err
	}
	intField(fields, "open_period", payload.OpenPeriod)
	int64Field(fields, "close_date", payload.CloseDate)
	boolField(fields, "is_closed", payload.IsClosed)
	boolField(fields, "disable_notification", payload.DisableNotification)
	boolField(fields, "protect_content", payload.ProtectContent)
	if err := replyParametersField(fields, payload.ReplyParameters); err != nil {
		return nil, err
	}
	if err := replyMarkupField(fields, payload.ReplyMarkup); err != nil {
		return nil, err
	}
	return fields, nil
}

func jsonField(fields map[string]string, name string, value any) error {
	if value == nil || isNilPollMediaInterface(value) {
		return nil
	}
	reflectValue := reflect.ValueOf(value)
	if reflectValue.Kind() == reflect.Slice && reflectValue.Len() == 0 {
		return nil
	}
	body, err := json.Marshal(value)
	if err != nil {
		return err
	}
	if string(body) == "null" {
		return nil
	}
	fields[name] = string(body)
	return nil
}

func requiredJSONField(fields map[string]string, name string, value any) error {
	body, err := json.Marshal(value)
	if err != nil {
		return err
	}
	fields[name] = string(body)
	return nil
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
	optionCount := len(params.Options)
	if len(params.OptionObjects) > 0 {
		if len(params.Options) > 0 {
			return stderrors.New("options and option_objects cannot both be set")
		}
		optionCount = len(params.OptionObjects)
		for _, option := range params.OptionObjects {
			if err := validateInputPollOption(option); err != nil {
				return err
			}
		}
	}
	if optionCount < 1 {
		return stderrors.New("options must contain at least one item")
	}
	for _, option := range params.Options {
		if option == "" {
			return stderrors.New("options must not contain empty items")
		}
	}
	if params.CorrectOptionID != nil && len(params.CorrectOptionIDs) > 0 {
		return stderrors.New("correct_option_id and correct_option_ids cannot both be set")
	}
	if params.CorrectOptionID != nil && (*params.CorrectOptionID < 0 || *params.CorrectOptionID >= optionCount) {
		return stderrors.New("correct_option_id must reference an existing option")
	}
	for _, id := range params.CorrectOptionIDs {
		if id < 0 || id >= optionCount {
			return stderrors.New("correct_option_ids must reference existing options")
		}
	}
	if params.OpenPeriod < 0 {
		return stderrors.New("open_period must not be negative")
	}
	if params.CloseDate < 0 {
		return stderrors.New("close_date must not be negative")
	}
	if err := validateEntityFormatting(params.QuestionParseMode, params.QuestionEntities); err != nil {
		return err
	}
	if err := validateEntityFormatting(params.ExplanationParseMode, params.ExplanationEntities); err != nil {
		return err
	}
	if err := validateInputPollMedia(params.ExplanationMedia, "explanation_media"); err != nil {
		return err
	}
	if err := validateEntityFormatting(params.DescriptionParseMode, params.DescriptionEntities); err != nil {
		return err
	}
	if err := validateInputPollMedia(params.Media, "media"); err != nil {
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

func validateInputPollOption(option telegram.InputPollOption) error {
	if option.Text == "" {
		return stderrors.New("options.text is required")
	}
	if err := validateEntityFormatting(option.TextParseMode, option.TextEntities); err != nil {
		return err
	}
	if err := validateInputPollOptionMedia(option.Media, "options.media"); err != nil {
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
