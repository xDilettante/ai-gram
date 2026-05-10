package bot

import (
	"context"
	"encoding/json"
	stderrors "errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/xDilettante/ai-gram/telegram"
)

// InputPaidMedia marks paid media items accepted by sendPaidMedia.
type InputPaidMedia interface {
	inputPaidMedia()
}

// InputPaidMediaPhoto describes a paid photo to send.
type InputPaidMediaPhoto struct {
	Type  string  `json:"type"`
	Media FileRef `json:"media"`
}

// InputPaidMediaLivePhoto describes a paid live photo to send.
type InputPaidMediaLivePhoto struct {
	Type  string  `json:"type"`
	Media FileRef `json:"media"`
	Photo FileRef `json:"photo"`
}

// InputPaidMediaVideo describes a paid video to send.
type InputPaidMediaVideo struct {
	Type              string  `json:"type"`
	Media             FileRef `json:"media"`
	Thumbnail         FileRef `json:"-"`
	Cover             FileRef `json:"-"`
	StartTimestamp    int     `json:"start_timestamp,omitempty"`
	Width             int     `json:"width,omitempty"`
	Height            int     `json:"height,omitempty"`
	Duration          int     `json:"duration,omitempty"`
	SupportsStreaming bool    `json:"supports_streaming,omitempty"`
}

func (InputPaidMediaPhoto) inputPaidMedia()     {}
func (InputPaidMediaLivePhoto) inputPaidMedia() {}
func (InputPaidMediaVideo) inputPaidMedia()     {}

// PaidPhoto creates a paid photo input media item.
func PaidPhoto(media FileRef) InputPaidMediaPhoto {
	return InputPaidMediaPhoto{Type: "photo", Media: media}
}

// PaidLivePhoto creates a paid live photo input media item.
func PaidLivePhoto(media FileRef, photo FileRef) InputPaidMediaLivePhoto {
	return InputPaidMediaLivePhoto{Type: "live_photo", Media: media, Photo: photo}
}

// PaidVideo creates a paid video input media item.
func PaidVideo(media FileRef) InputPaidMediaVideo {
	return InputPaidMediaVideo{Type: "video", Media: media}
}

// SendPaidMediaParams contains supported parameters for sendPaidMedia.
type SendPaidMediaParams struct {
	BusinessConnectionID    string                            `json:"business_connection_id,omitempty"`
	ChatID                  ChatID                            `json:"chat_id"`
	MessageThreadID         int64                             `json:"message_thread_id,omitempty"`
	DirectMessagesTopicID   int64                             `json:"direct_messages_topic_id,omitempty"`
	StarCount               int                               `json:"star_count"`
	Media                   []InputPaidMedia                  `json:"media"`
	Payload                 string                            `json:"payload,omitempty"`
	Caption                 string                            `json:"caption,omitempty"`
	ParseMode               string                            `json:"parse_mode,omitempty"`
	CaptionEntities         []telegram.MessageEntity          `json:"caption_entities,omitempty"`
	ShowCaptionAboveMedia   bool                              `json:"show_caption_above_media,omitempty"`
	DisableNotification     bool                              `json:"disable_notification,omitempty"`
	ProtectContent          bool                              `json:"protect_content,omitempty"`
	AllowPaidBroadcast      bool                              `json:"allow_paid_broadcast,omitempty"`
	SuggestedPostParameters *telegram.SuggestedPostParameters `json:"suggested_post_parameters,omitempty"`
	ReplyParameters         *telegram.ReplyParameters         `json:"reply_parameters,omitempty"`
	ReplyMarkup             telegram.ReplyMarkup              `json:"reply_markup,omitempty"`
}

// GetStarTransactionsParams contains supported parameters for getStarTransactions.
type GetStarTransactionsParams struct {
	Offset int `json:"offset,omitempty"`
	Limit  int `json:"limit,omitempty"`
}

// RefundStarPaymentParams contains supported parameters for refundStarPayment.
type RefundStarPaymentParams struct {
	UserID                  int64  `json:"user_id"`
	TelegramPaymentChargeID string `json:"telegram_payment_charge_id"`
}

// SendPaidMedia sends paid media and returns the sent message.
func (b *Bot) SendPaidMedia(ctx context.Context, params SendPaidMediaParams) (*telegram.Message, error) {
	if err := params.validate(); err != nil {
		return nil, err
	}
	media, files, err := params.mediaPayload()
	if err != nil {
		return nil, err
	}

	var message telegram.Message
	if len(files) > 0 {
		fields, err := params.multipartFields(media)
		if err != nil {
			return nil, err
		}
		if err := b.callMultipart(ctx, "sendPaidMedia", fields, files, &message); err != nil {
			return nil, err
		}
		return &message, nil
	}

	payload := sendPaidMediaPayload{
		BusinessConnectionID:    params.BusinessConnectionID,
		ChatID:                  params.ChatID,
		MessageThreadID:         params.MessageThreadID,
		DirectMessagesTopicID:   params.DirectMessagesTopicID,
		StarCount:               params.StarCount,
		Media:                   media,
		Payload:                 params.Payload,
		Caption:                 params.Caption,
		ParseMode:               params.ParseMode,
		CaptionEntities:         params.CaptionEntities,
		ShowCaptionAboveMedia:   params.ShowCaptionAboveMedia,
		DisableNotification:     params.DisableNotification,
		ProtectContent:          params.ProtectContent,
		AllowPaidBroadcast:      params.AllowPaidBroadcast,
		SuggestedPostParameters: params.SuggestedPostParameters,
		ReplyParameters:         params.ReplyParameters,
		ReplyMarkup:             params.ReplyMarkup,
	}
	if err := b.call(ctx, "sendPaidMedia", payload, &message); err != nil {
		return nil, err
	}
	return &message, nil
}

// GetStarTransactions returns the bot's Telegram Star transactions.
func (b *Bot) GetStarTransactions(ctx context.Context, params GetStarTransactionsParams) (*telegram.StarTransactions, error) {
	if err := params.validate(); err != nil {
		return nil, err
	}
	var transactions telegram.StarTransactions
	if err := b.call(ctx, "getStarTransactions", params, &transactions); err != nil {
		return nil, err
	}
	return &transactions, nil
}

// RefundStarPayment refunds a successful Telegram Stars payment.
func (b *Bot) RefundStarPayment(ctx context.Context, params RefundStarPaymentParams) (bool, error) {
	if err := params.validate(); err != nil {
		return false, err
	}
	var result bool
	if err := b.call(ctx, "refundStarPayment", params, &result); err != nil {
		return false, err
	}
	return result, nil
}

type sendPaidMediaPayload struct {
	BusinessConnectionID    string                            `json:"business_connection_id,omitempty"`
	ChatID                  ChatID                            `json:"chat_id"`
	MessageThreadID         int64                             `json:"message_thread_id,omitempty"`
	DirectMessagesTopicID   int64                             `json:"direct_messages_topic_id,omitempty"`
	StarCount               int                               `json:"star_count"`
	Media                   []inputPaidMediaPayload           `json:"media"`
	Payload                 string                            `json:"payload,omitempty"`
	Caption                 string                            `json:"caption,omitempty"`
	ParseMode               string                            `json:"parse_mode,omitempty"`
	CaptionEntities         []telegram.MessageEntity          `json:"caption_entities,omitempty"`
	ShowCaptionAboveMedia   bool                              `json:"show_caption_above_media,omitempty"`
	DisableNotification     bool                              `json:"disable_notification,omitempty"`
	ProtectContent          bool                              `json:"protect_content,omitempty"`
	AllowPaidBroadcast      bool                              `json:"allow_paid_broadcast,omitempty"`
	SuggestedPostParameters *telegram.SuggestedPostParameters `json:"suggested_post_parameters,omitempty"`
	ReplyParameters         *telegram.ReplyParameters         `json:"reply_parameters,omitempty"`
	ReplyMarkup             telegram.ReplyMarkup              `json:"reply_markup,omitempty"`
}

type inputPaidMediaPayload struct {
	Type              string `json:"type"`
	Media             string `json:"media"`
	Photo             string `json:"photo,omitempty"`
	Thumbnail         string `json:"thumbnail,omitempty"`
	Cover             string `json:"cover,omitempty"`
	StartTimestamp    int    `json:"start_timestamp,omitempty"`
	Width             int    `json:"width,omitempty"`
	Height            int    `json:"height,omitempty"`
	Duration          int    `json:"duration,omitempty"`
	SupportsStreaming bool   `json:"supports_streaming,omitempty"`
}

func (params SendPaidMediaParams) validate() error {
	if !params.ChatID.valid() {
		return stderrors.New("chat_id is required")
	}
	if err := validateMessageThreadID(params.MessageThreadID); err != nil {
		return err
	}
	if params.DirectMessagesTopicID < 0 {
		return stderrors.New("direct_messages_topic_id must not be negative")
	}
	if params.StarCount <= 0 {
		return stderrors.New("star_count must be greater than zero")
	}
	if len(params.Media) == 0 {
		return stderrors.New("media must not be empty")
	}
	if len(params.Media) > 10 {
		return stderrors.New("media must contain at most ten items")
	}
	for index, media := range params.Media {
		if err := validateInputPaidMedia(media); err != nil {
			return fmt.Errorf("media[%d]: %w", index, err)
		}
	}
	if err := validateCaptionFormatting(params.ParseMode, params.CaptionEntities); err != nil {
		return err
	}
	if err := validateSuggestedPostParameters(params.SuggestedPostParameters); err != nil {
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

func (params GetStarTransactionsParams) validate() error {
	if params.Offset < 0 {
		return stderrors.New("offset must not be negative")
	}
	if params.Limit < 0 {
		return stderrors.New("limit must not be negative")
	}
	return nil
}

func (params RefundStarPaymentParams) validate() error {
	if params.UserID <= 0 {
		return stderrors.New("user_id must be greater than zero")
	}
	if strings.TrimSpace(params.TelegramPaymentChargeID) == "" {
		return stderrors.New("telegram_payment_charge_id is required")
	}
	return nil
}

func validateInputPaidMedia(media InputPaidMedia) error {
	switch item := media.(type) {
	case nil:
		return stderrors.New("paid media item is required")
	case InputPaidMediaPhoto:
		return validateInputPaidMediaPhoto(item)
	case *InputPaidMediaPhoto:
		if item == nil {
			return stderrors.New("paid media item is required")
		}
		return validateInputPaidMediaPhoto(*item)
	case InputPaidMediaLivePhoto:
		return validateInputPaidMediaLivePhoto(item)
	case *InputPaidMediaLivePhoto:
		if item == nil {
			return stderrors.New("paid media item is required")
		}
		return validateInputPaidMediaLivePhoto(*item)
	case InputPaidMediaVideo:
		return validateInputPaidMediaVideo(item)
	case *InputPaidMediaVideo:
		if item == nil {
			return stderrors.New("paid media item is required")
		}
		return validateInputPaidMediaVideo(*item)
	default:
		return stderrors.New("unsupported input paid media type")
	}
}

func validateInputPaidMediaPhoto(media InputPaidMediaPhoto) error {
	if err := validateInputMediaType(media.Type, "photo"); err != nil {
		return err
	}
	return media.Media.validate("media")
}

func validateInputPaidMediaLivePhoto(media InputPaidMediaLivePhoto) error {
	if err := validateInputMediaType(media.Type, "live_photo"); err != nil {
		return err
	}
	if err := validateLivePhotoFileRef(media.Media, "media"); err != nil {
		return err
	}
	return validateLivePhotoFileRef(media.Photo, "photo")
}

func validateInputPaidMediaVideo(media InputPaidMediaVideo) error {
	if err := validateInputMediaType(media.Type, "video"); err != nil {
		return err
	}
	if err := media.Media.validate("media"); err != nil {
		return err
	}
	if err := validateOptionalFileRef(media.Thumbnail, "thumbnail"); err != nil {
		return err
	}
	if err := validateOptionalFileRef(media.Cover, "cover"); err != nil {
		return err
	}
	if media.StartTimestamp < 0 {
		return stderrors.New("start_timestamp must not be negative")
	}
	if media.Width < 0 {
		return stderrors.New("width must not be negative")
	}
	if media.Height < 0 {
		return stderrors.New("height must not be negative")
	}
	if media.Duration < 0 {
		return stderrors.New("duration must not be negative")
	}
	return nil
}

func (params SendPaidMediaParams) mediaPayload() ([]inputPaidMediaPayload, map[string]UploadFile, error) {
	payload := make([]inputPaidMediaPayload, 0, len(params.Media))
	files := make(map[string]UploadFile)
	for index, media := range params.Media {
		item, err := buildInputPaidMediaPayload(media, index, files)
		if err != nil {
			return nil, nil, fmt.Errorf("media[%d]: %w", index, err)
		}
		payload = append(payload, item)
	}
	return payload, files, nil
}

func buildInputPaidMediaPayload(media InputPaidMedia, index int, files map[string]UploadFile) (inputPaidMediaPayload, error) {
	switch item := media.(type) {
	case InputPaidMediaPhoto:
		return buildInputPaidMediaPhotoPayload(item, index, files)
	case *InputPaidMediaPhoto:
		if item == nil {
			return inputPaidMediaPayload{}, stderrors.New("paid media item is required")
		}
		return buildInputPaidMediaPhotoPayload(*item, index, files)
	case InputPaidMediaLivePhoto:
		return buildInputPaidMediaLivePhotoPayload(item, index, files)
	case *InputPaidMediaLivePhoto:
		if item == nil {
			return inputPaidMediaPayload{}, stderrors.New("paid media item is required")
		}
		return buildInputPaidMediaLivePhotoPayload(*item, index, files)
	case InputPaidMediaVideo:
		return buildInputPaidMediaVideoPayload(item, index, files)
	case *InputPaidMediaVideo:
		if item == nil {
			return inputPaidMediaPayload{}, stderrors.New("paid media item is required")
		}
		return buildInputPaidMediaVideoPayload(*item, index, files)
	default:
		return inputPaidMediaPayload{}, stderrors.New("unsupported input paid media type")
	}
}

func buildInputPaidMediaPhotoPayload(media InputPaidMediaPhoto, index int, files map[string]UploadFile) (inputPaidMediaPayload, error) {
	mediaValue, err := paidMediaFileValue(media.Media, fmt.Sprintf("media%d", index), files)
	if err != nil {
		return inputPaidMediaPayload{}, err
	}
	return inputPaidMediaPayload{Type: mediaType(media.Type, "photo"), Media: mediaValue}, nil
}

func buildInputPaidMediaLivePhotoPayload(media InputPaidMediaLivePhoto, index int, files map[string]UploadFile) (inputPaidMediaPayload, error) {
	mediaValue, err := paidLivePhotoFileValue(media.Media, "media", fmt.Sprintf("media%d", index), files)
	if err != nil {
		return inputPaidMediaPayload{}, err
	}
	photoValue, err := paidLivePhotoFileValue(media.Photo, "photo", fmt.Sprintf("photo%d", index), files)
	if err != nil {
		return inputPaidMediaPayload{}, err
	}
	return inputPaidMediaPayload{Type: mediaType(media.Type, "live_photo"), Media: mediaValue, Photo: photoValue}, nil
}

func buildInputPaidMediaVideoPayload(media InputPaidMediaVideo, index int, files map[string]UploadFile) (inputPaidMediaPayload, error) {
	mediaValue, err := paidMediaFileValue(media.Media, fmt.Sprintf("media%d", index), files)
	if err != nil {
		return inputPaidMediaPayload{}, err
	}
	thumbnail, err := optionalPaidMediaFileValue(media.Thumbnail, fmt.Sprintf("thumb%d", index), files)
	if err != nil {
		return inputPaidMediaPayload{}, err
	}
	cover, err := optionalPaidMediaFileValue(media.Cover, fmt.Sprintf("cover%d", index), files)
	if err != nil {
		return inputPaidMediaPayload{}, err
	}
	return inputPaidMediaPayload{
		Type:              mediaType(media.Type, "video"),
		Media:             mediaValue,
		Thumbnail:         thumbnail,
		Cover:             cover,
		StartTimestamp:    media.StartTimestamp,
		Width:             media.Width,
		Height:            media.Height,
		Duration:          media.Duration,
		SupportsStreaming: media.SupportsStreaming,
	}, nil
}

func paidMediaFileValue(ref FileRef, name string, files map[string]UploadFile) (string, error) {
	if err := ref.validate("media"); err != nil {
		return "", err
	}
	if ref.isUpload() {
		files[name] = ref.upload
		return "attach://" + name, nil
	}
	return ref.value, nil
}

func paidLivePhotoFileValue(ref FileRef, field string, name string, files map[string]UploadFile) (string, error) {
	if err := validateLivePhotoFileRef(ref, field); err != nil {
		return "", err
	}
	if ref.isUpload() {
		files[name] = ref.upload
		return "attach://" + name, nil
	}
	return ref.value, nil
}

func optionalPaidMediaFileValue(ref FileRef, name string, files map[string]UploadFile) (string, error) {
	if !ref.isSet() {
		return "", nil
	}
	if err := ref.validate(name); err != nil {
		return "", err
	}
	if ref.isUpload() {
		files[name] = ref.upload
		return "attach://" + name, nil
	}
	return ref.value, nil
}

func (params SendPaidMediaParams) multipartFields(media []inputPaidMediaPayload) (map[string]string, error) {
	chatIDValue, err := params.ChatID.multipartValue()
	if err != nil {
		return nil, err
	}
	fields := map[string]string{"chat_id": chatIDValue}
	stringField(fields, "business_connection_id", params.BusinessConnectionID)
	int64Field(fields, "message_thread_id", params.MessageThreadID)
	int64Field(fields, "direct_messages_topic_id", params.DirectMessagesTopicID)
	fields["star_count"] = strconv.Itoa(params.StarCount)
	stringField(fields, "payload", params.Payload)
	stringField(fields, "caption", params.Caption)
	stringField(fields, "parse_mode", params.ParseMode)
	if err := captionEntitiesField(fields, params.CaptionEntities); err != nil {
		return nil, err
	}
	boolField(fields, "show_caption_above_media", params.ShowCaptionAboveMedia)
	boolField(fields, "disable_notification", params.DisableNotification)
	boolField(fields, "protect_content", params.ProtectContent)
	boolField(fields, "allow_paid_broadcast", params.AllowPaidBroadcast)
	if params.SuggestedPostParameters != nil {
		body, err := json.Marshal(params.SuggestedPostParameters)
		if err != nil {
			return nil, err
		}
		fields["suggested_post_parameters"] = string(body)
	}
	if err := replyParametersField(fields, params.ReplyParameters); err != nil {
		return nil, err
	}
	if err := replyMarkupField(fields, params.ReplyMarkup); err != nil {
		return nil, err
	}
	body, err := json.Marshal(media)
	if err != nil {
		return nil, err
	}
	fields["media"] = string(body)
	return fields, nil
}
