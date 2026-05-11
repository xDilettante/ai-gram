package bot

import (
	"bytes"
	"context"
	"encoding/json"
	stderrors "errors"
	"strconv"
	"strings"

	"github.com/xDilettante/ai-gram/telegram"
)

// EditMessageResult contains the polymorphic result returned by Telegram edit message methods.
//
// Telegram returns an edited Message when the target is a regular chat message, or true when the
// target is an inline message. OK is true for both successful shapes.
type EditMessageResult struct {
	Message *telegram.Message
	OK      bool
}

// IsMessage reports whether the edit result contains an edited Message object.
func (r EditMessageResult) IsMessage() bool {
	return r.Message != nil
}

// IsOK reports whether Telegram returned a successful edit result.
func (r EditMessageResult) IsOK() bool {
	return r.OK
}

// UnmarshalJSON decodes Telegram edit results, which can be either a Message object or a boolean.
func (r *EditMessageResult) UnmarshalJSON(data []byte) error {
	trimmed := bytes.TrimSpace(data)
	if bytes.Equal(trimmed, []byte("true")) {
		r.OK = true
		r.Message = nil
		return nil
	}
	if bytes.Equal(trimmed, []byte("false")) {
		r.OK = false
		r.Message = nil
		return nil
	}
	if len(trimmed) > 0 && trimmed[0] == '{' {
		var message telegram.Message
		if err := json.Unmarshal(trimmed, &message); err != nil {
			return err
		}
		r.OK = true
		r.Message = &message
		return nil
	}

	return stderrors.New("edit message result must be a message object or boolean")
}

// EditMessageTarget identifies a message that can be edited, either by chat/message ID or inline message ID.
type EditMessageTarget struct {
	ChatID          ChatID `json:"chat_id,omitempty"`
	MessageID       int64  `json:"message_id,omitempty"`
	InlineMessageID string `json:"inline_message_id,omitempty"`
}

// EditTargetChat creates an edit target for a regular chat message.
func EditTargetChat(chatID ChatID, messageID int64) EditMessageTarget {
	return EditMessageTarget{ChatID: chatID, MessageID: messageID}
}

// EditTargetInline creates an edit target for an inline message.
func EditTargetInline(inlineMessageID string) EditMessageTarget {
	return EditMessageTarget{InlineMessageID: inlineMessageID}
}

// EditMessageTextParams contains supported parameters for editMessageText.
type EditMessageTextParams struct {
	BusinessConnectionID string                         `json:"business_connection_id,omitempty"`
	Target               EditMessageTarget              `json:"-"`
	Text                 string                         `json:"text"`
	ParseMode            string                         `json:"parse_mode,omitempty"`
	Entities             []telegram.MessageEntity       `json:"entities,omitempty"`
	LinkPreviewDisabled  bool                           `json:"disable_web_page_preview,omitempty"`
	ReplyMarkup          *telegram.InlineKeyboardMarkup `json:"reply_markup,omitempty"`
}

// EditMessageReplyMarkupParams contains supported parameters for editMessageReplyMarkup.
type EditMessageReplyMarkupParams struct {
	BusinessConnectionID string                         `json:"business_connection_id,omitempty"`
	Target               EditMessageTarget              `json:"-"`
	ReplyMarkup          *telegram.InlineKeyboardMarkup `json:"reply_markup,omitempty"`
}

// EditMessageCaptionParams contains supported parameters for editMessageCaption.
type EditMessageCaptionParams struct {
	BusinessConnectionID string                         `json:"business_connection_id,omitempty"`
	Target               EditMessageTarget              `json:"-"`
	Caption              string                         `json:"caption,omitempty"`
	ParseMode            string                         `json:"parse_mode,omitempty"`
	CaptionEntities      []telegram.MessageEntity       `json:"caption_entities,omitempty"`
	ReplyMarkup          *telegram.InlineKeyboardMarkup `json:"reply_markup,omitempty"`
}

// EditMessageMediaParams contains supported parameters for editMessageMedia.
type EditMessageMediaParams struct {
	BusinessConnectionID string                         `json:"business_connection_id,omitempty"`
	Target               EditMessageTarget              `json:"-"`
	Media                InputMedia                     `json:"media"`
	ReplyMarkup          *telegram.InlineKeyboardMarkup `json:"reply_markup,omitempty"`
}

// EditMessageLiveLocationParams contains supported parameters for editMessageLiveLocation.
type EditMessageLiveLocationParams struct {
	BusinessConnectionID string                         `json:"business_connection_id,omitempty"`
	Target               EditMessageTarget              `json:"-"`
	Latitude             float64                        `json:"latitude"`
	Longitude            float64                        `json:"longitude"`
	LivePeriod           int                            `json:"live_period,omitempty"`
	HorizontalAccuracy   float64                        `json:"horizontal_accuracy,omitempty"`
	Heading              int                            `json:"heading,omitempty"`
	ProximityAlertRadius int                            `json:"proximity_alert_radius,omitempty"`
	ReplyMarkup          *telegram.InlineKeyboardMarkup `json:"reply_markup,omitempty"`
}

// StopMessageLiveLocationParams contains supported parameters for stopMessageLiveLocation.
type StopMessageLiveLocationParams struct {
	BusinessConnectionID string                         `json:"business_connection_id,omitempty"`
	Target               EditMessageTarget              `json:"-"`
	ReplyMarkup          *telegram.InlineKeyboardMarkup `json:"reply_markup,omitempty"`
}

// EditMessageText edits text and inline keyboard markup for an existing chat or inline message.
func (b *Bot) EditMessageText(ctx context.Context, params EditMessageTextParams) (*EditMessageResult, error) {
	if err := params.validate(); err != nil {
		return nil, err
	}

	var result EditMessageResult
	if err := b.call(ctx, "editMessageText", params.payload(), &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// EditMessageReplyMarkup edits or removes inline keyboard markup for an existing chat or inline message.
func (b *Bot) EditMessageReplyMarkup(ctx context.Context, params EditMessageReplyMarkupParams) (*EditMessageResult, error) {
	if err := params.validate(); err != nil {
		return nil, err
	}

	var result EditMessageResult
	if err := b.call(ctx, "editMessageReplyMarkup", params.payload(), &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// EditMessageCaption edits the caption and inline keyboard markup for an existing chat or inline message.
func (b *Bot) EditMessageCaption(ctx context.Context, params EditMessageCaptionParams) (*EditMessageResult, error) {
	if err := params.validate(); err != nil {
		return nil, err
	}

	var result EditMessageResult
	if err := b.call(ctx, "editMessageCaption", params.payload(), &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// EditMessageMedia edits animation, audio, document, live photo, photo, or video media for an existing chat or inline message.
func (b *Bot) EditMessageMedia(ctx context.Context, params EditMessageMediaParams) (*EditMessageResult, error) {
	if err := params.validate(); err != nil {
		return nil, err
	}

	media, files, err := params.mediaPayload()
	if err != nil {
		return nil, err
	}
	if params.Target.isInline() && len(files) > 0 {
		return nil, stderrors.New("inline message media cannot use FileUpload")
	}

	var result EditMessageResult
	if len(files) > 0 {
		fields, err := params.multipartFields(media)
		if err != nil {
			return nil, err
		}
		if err := b.callMultipartBuffered(ctx, "editMessageMedia", fields, files, &result); err != nil {
			return nil, err
		}
		return &result, nil
	}

	if err := b.call(ctx, "editMessageMedia", params.payload(media), &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// EditMessageLiveLocation edits a live location for an existing chat or inline message.
func (b *Bot) EditMessageLiveLocation(ctx context.Context, params EditMessageLiveLocationParams) (*EditMessageResult, error) {
	if err := params.validate(); err != nil {
		return nil, err
	}

	var result EditMessageResult
	if err := b.call(ctx, "editMessageLiveLocation", params.payload(), &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// StopMessageLiveLocation stops updating a live location for an existing chat or inline message.
func (b *Bot) StopMessageLiveLocation(ctx context.Context, params StopMessageLiveLocationParams) (*EditMessageResult, error) {
	if err := params.validate(); err != nil {
		return nil, err
	}

	var result EditMessageResult
	if err := b.call(ctx, "stopMessageLiveLocation", params.payload(), &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (target EditMessageTarget) validate() error {
	inlineSet := strings.TrimSpace(target.InlineMessageID) != ""
	chatIDSet := target.ChatID.valid()
	messageIDSet := target.MessageID != 0
	chatMode := chatIDSet || messageIDSet

	if inlineSet && chatMode {
		return stderrors.New("edit target must use either chat_id/message_id or inline_message_id, not both")
	}
	if !inlineSet && !chatMode {
		return stderrors.New("edit target is required")
	}
	if inlineSet {
		return nil
	}
	if !chatIDSet {
		return stderrors.New("chat_id is required for chat edit target")
	}
	if target.MessageID <= 0 {
		return stderrors.New("message_id must be greater than zero for chat edit target")
	}

	return nil
}

func (params EditMessageTextParams) validate() error {
	if err := params.Target.validate(); err != nil {
		return err
	}
	if params.Text == "" {
		return stderrors.New("text is required")
	}
	if err := validateEntityFormatting(params.ParseMode, params.Entities); err != nil {
		return err
	}
	if params.ReplyMarkup != nil {
		if err := telegram.ValidateReplyMarkup(*params.ReplyMarkup); err != nil {
			return err
		}
	}

	return nil
}

func (params EditMessageReplyMarkupParams) validate() error {
	if err := params.Target.validate(); err != nil {
		return err
	}
	if params.ReplyMarkup != nil {
		if err := telegram.ValidateReplyMarkup(*params.ReplyMarkup); err != nil {
			return err
		}
	}

	return nil
}

func (params EditMessageCaptionParams) validate() error {
	if err := params.Target.validate(); err != nil {
		return err
	}
	if err := validateCaptionFormatting(params.ParseMode, params.CaptionEntities); err != nil {
		return err
	}
	if params.ReplyMarkup != nil {
		if err := telegram.ValidateReplyMarkup(*params.ReplyMarkup); err != nil {
			return err
		}
	}

	return nil
}

func (params EditMessageMediaParams) validate() error {
	if err := params.Target.validate(); err != nil {
		return err
	}
	if err := validateInputMediaForEdit(params.Media); err != nil {
		return err
	}
	if params.ReplyMarkup != nil {
		if err := telegram.ValidateReplyMarkup(*params.ReplyMarkup); err != nil {
			return err
		}
	}

	return nil
}

func (params EditMessageLiveLocationParams) validate() error {
	if err := params.Target.validate(); err != nil {
		return err
	}
	if err := validateLatitude(params.Latitude); err != nil {
		return err
	}
	if err := validateLongitude(params.Longitude); err != nil {
		return err
	}
	if params.LivePeriod < 0 {
		return stderrors.New("live_period must not be negative")
	}
	if params.HorizontalAccuracy < 0 {
		return stderrors.New("horizontal_accuracy must not be negative")
	}
	if params.Heading < 0 {
		return stderrors.New("heading must not be negative")
	}
	if params.ProximityAlertRadius < 0 {
		return stderrors.New("proximity_alert_radius must not be negative")
	}
	if params.ReplyMarkup != nil {
		if err := telegram.ValidateReplyMarkup(*params.ReplyMarkup); err != nil {
			return err
		}
	}

	return nil
}

func (params StopMessageLiveLocationParams) validate() error {
	if err := params.Target.validate(); err != nil {
		return err
	}
	if params.ReplyMarkup != nil {
		if err := telegram.ValidateReplyMarkup(*params.ReplyMarkup); err != nil {
			return err
		}
	}

	return nil
}

type editMessageTextPayload struct {
	BusinessConnectionID string                         `json:"business_connection_id,omitempty"`
	ChatID               *ChatID                        `json:"chat_id,omitempty"`
	MessageID            int64                          `json:"message_id,omitempty"`
	InlineMessageID      string                         `json:"inline_message_id,omitempty"`
	Text                 string                         `json:"text"`
	ParseMode            string                         `json:"parse_mode,omitempty"`
	Entities             []telegram.MessageEntity       `json:"entities,omitempty"`
	LinkPreviewDisabled  bool                           `json:"disable_web_page_preview,omitempty"`
	ReplyMarkup          *telegram.InlineKeyboardMarkup `json:"reply_markup,omitempty"`
}

type editMessageReplyMarkupPayload struct {
	BusinessConnectionID string                         `json:"business_connection_id,omitempty"`
	ChatID               *ChatID                        `json:"chat_id,omitempty"`
	MessageID            int64                          `json:"message_id,omitempty"`
	InlineMessageID      string                         `json:"inline_message_id,omitempty"`
	ReplyMarkup          *telegram.InlineKeyboardMarkup `json:"reply_markup,omitempty"`
}

type editMessageCaptionPayload struct {
	BusinessConnectionID string                         `json:"business_connection_id,omitempty"`
	ChatID               *ChatID                        `json:"chat_id,omitempty"`
	MessageID            int64                          `json:"message_id,omitempty"`
	InlineMessageID      string                         `json:"inline_message_id,omitempty"`
	Caption              string                         `json:"caption"`
	ParseMode            string                         `json:"parse_mode,omitempty"`
	CaptionEntities      []telegram.MessageEntity       `json:"caption_entities,omitempty"`
	ReplyMarkup          *telegram.InlineKeyboardMarkup `json:"reply_markup,omitempty"`
}

type editMessageMediaPayload struct {
	BusinessConnectionID string                         `json:"business_connection_id,omitempty"`
	ChatID               *ChatID                        `json:"chat_id,omitempty"`
	MessageID            int64                          `json:"message_id,omitempty"`
	InlineMessageID      string                         `json:"inline_message_id,omitempty"`
	Media                inputMediaPayload              `json:"media"`
	ReplyMarkup          *telegram.InlineKeyboardMarkup `json:"reply_markup,omitempty"`
}

type editMessageLiveLocationPayload struct {
	BusinessConnectionID string                         `json:"business_connection_id,omitempty"`
	ChatID               *ChatID                        `json:"chat_id,omitempty"`
	MessageID            int64                          `json:"message_id,omitempty"`
	InlineMessageID      string                         `json:"inline_message_id,omitempty"`
	Latitude             float64                        `json:"latitude"`
	Longitude            float64                        `json:"longitude"`
	LivePeriod           int                            `json:"live_period,omitempty"`
	HorizontalAccuracy   float64                        `json:"horizontal_accuracy,omitempty"`
	Heading              int                            `json:"heading,omitempty"`
	ProximityAlertRadius int                            `json:"proximity_alert_radius,omitempty"`
	ReplyMarkup          *telegram.InlineKeyboardMarkup `json:"reply_markup,omitempty"`
}

type stopMessageLiveLocationPayload struct {
	BusinessConnectionID string                         `json:"business_connection_id,omitempty"`
	ChatID               *ChatID                        `json:"chat_id,omitempty"`
	MessageID            int64                          `json:"message_id,omitempty"`
	InlineMessageID      string                         `json:"inline_message_id,omitempty"`
	ReplyMarkup          *telegram.InlineKeyboardMarkup `json:"reply_markup,omitempty"`
}

func (params EditMessageTextParams) payload() editMessageTextPayload {
	chatID, messageID, inlineMessageID := params.Target.payloadValues()
	payload := editMessageTextPayload{
		BusinessConnectionID: params.BusinessConnectionID,
		ChatID:               chatID,
		MessageID:            messageID,
		InlineMessageID:      inlineMessageID,
		Text:                 params.Text,
		ParseMode:            params.ParseMode,
		Entities:             params.Entities,
		LinkPreviewDisabled:  params.LinkPreviewDisabled,
		ReplyMarkup:          params.ReplyMarkup,
	}

	return payload
}

func (params EditMessageReplyMarkupParams) payload() editMessageReplyMarkupPayload {
	chatID, messageID, inlineMessageID := params.Target.payloadValues()
	payload := editMessageReplyMarkupPayload{
		BusinessConnectionID: params.BusinessConnectionID,
		ChatID:               chatID,
		MessageID:            messageID,
		InlineMessageID:      inlineMessageID,
		ReplyMarkup:          params.ReplyMarkup,
	}

	return payload
}

func (params EditMessageCaptionParams) payload() editMessageCaptionPayload {
	chatID, messageID, inlineMessageID := params.Target.payloadValues()
	return editMessageCaptionPayload{
		BusinessConnectionID: params.BusinessConnectionID,
		ChatID:               chatID,
		MessageID:            messageID,
		InlineMessageID:      inlineMessageID,
		Caption:              params.Caption,
		ParseMode:            params.ParseMode,
		CaptionEntities:      params.CaptionEntities,
		ReplyMarkup:          params.ReplyMarkup,
	}
}

func (params EditMessageMediaParams) mediaPayload() (inputMediaPayload, map[string]UploadFile, error) {
	files := make(map[string]UploadFile)
	media, err := buildInputMediaPayload(params.Media, 0, files)
	if err != nil {
		return inputMediaPayload{}, nil, err
	}
	return media, files, nil
}

func (params EditMessageMediaParams) payload(media inputMediaPayload) editMessageMediaPayload {
	chatID, messageID, inlineMessageID := params.Target.payloadValues()
	return editMessageMediaPayload{
		BusinessConnectionID: params.BusinessConnectionID,
		ChatID:               chatID,
		MessageID:            messageID,
		InlineMessageID:      inlineMessageID,
		Media:                media,
		ReplyMarkup:          params.ReplyMarkup,
	}
}

func (params EditMessageMediaParams) multipartFields(media inputMediaPayload) (map[string]string, error) {
	fields, err := params.Target.multipartFields()
	if err != nil {
		return nil, err
	}
	stringField(fields, "business_connection_id", params.BusinessConnectionID)
	if err := inlineReplyMarkupField(fields, params.ReplyMarkup); err != nil {
		return nil, err
	}
	body, err := json.Marshal(media)
	if err != nil {
		return nil, err
	}
	fields["media"] = string(body)
	return fields, nil
}

func (params EditMessageLiveLocationParams) payload() editMessageLiveLocationPayload {
	chatID, messageID, inlineMessageID := params.Target.payloadValues()
	return editMessageLiveLocationPayload{
		BusinessConnectionID: params.BusinessConnectionID,
		ChatID:               chatID,
		MessageID:            messageID,
		InlineMessageID:      inlineMessageID,
		Latitude:             params.Latitude,
		Longitude:            params.Longitude,
		LivePeriod:           params.LivePeriod,
		HorizontalAccuracy:   params.HorizontalAccuracy,
		Heading:              params.Heading,
		ProximityAlertRadius: params.ProximityAlertRadius,
		ReplyMarkup:          params.ReplyMarkup,
	}
}

func (params StopMessageLiveLocationParams) payload() stopMessageLiveLocationPayload {
	chatID, messageID, inlineMessageID := params.Target.payloadValues()
	return stopMessageLiveLocationPayload{
		BusinessConnectionID: params.BusinessConnectionID,
		ChatID:               chatID,
		MessageID:            messageID,
		InlineMessageID:      inlineMessageID,
		ReplyMarkup:          params.ReplyMarkup,
	}
}

func (target EditMessageTarget) payloadValues() (*ChatID, int64, string) {
	if target.ChatID.valid() {
		chatID := target.ChatID
		return &chatID, target.MessageID, ""
	}

	return nil, 0, target.InlineMessageID
}

func (target EditMessageTarget) isInline() bool {
	return strings.TrimSpace(target.InlineMessageID) != ""
}

func (target EditMessageTarget) multipartFields() (map[string]string, error) {
	fields := make(map[string]string)
	if target.ChatID.valid() {
		chatIDValue, err := target.ChatID.multipartValue()
		if err != nil {
			return nil, err
		}
		fields["chat_id"] = chatIDValue
		fields["message_id"] = strconv.FormatInt(target.MessageID, 10)
		return fields, nil
	}
	fields["inline_message_id"] = target.InlineMessageID
	return fields, nil
}

func inlineReplyMarkupField(fields map[string]string, replyMarkup *telegram.InlineKeyboardMarkup) error {
	if replyMarkup == nil {
		return nil
	}
	body, err := json.Marshal(replyMarkup)
	if err != nil {
		return err
	}
	fields["reply_markup"] = string(body)
	return nil
}

func validateEntityFormatting(parseMode string, entities []telegram.MessageEntity) error {
	if parseMode != "" && len(entities) > 0 {
		return stderrors.New("parse_mode and entities cannot be used together")
	}

	return nil
}
