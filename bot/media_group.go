package bot

import (
	"context"
	"encoding/json"
	stderrors "errors"
	"fmt"
	"strings"

	"github.com/xDilettante/ai-gram/telegram"
)

// InputMedia describes one media item accepted by sendMediaGroup.
type InputMedia interface {
	inputMedia()
}

// InputMediaPhoto describes a photo item for sendMediaGroup.
type InputMediaPhoto struct {
	Type            string                   `json:"type"`
	Media           FileRef                  `json:"media"`
	Caption         string                   `json:"caption,omitempty"`
	ParseMode       string                   `json:"parse_mode,omitempty"`
	CaptionEntities []telegram.MessageEntity `json:"caption_entities,omitempty"`
	HasSpoiler      bool                     `json:"has_spoiler,omitempty"`
}

// InputMediaVideo describes a video item for sendMediaGroup.
type InputMediaVideo struct {
	Type              string                   `json:"type"`
	Media             FileRef                  `json:"media"`
	Thumbnail         FileRef                  `json:"-"`
	Caption           string                   `json:"caption,omitempty"`
	ParseMode         string                   `json:"parse_mode,omitempty"`
	CaptionEntities   []telegram.MessageEntity `json:"caption_entities,omitempty"`
	Width             int                      `json:"width,omitempty"`
	Height            int                      `json:"height,omitempty"`
	Duration          int                      `json:"duration,omitempty"`
	SupportsStreaming bool                     `json:"supports_streaming,omitempty"`
	HasSpoiler        bool                     `json:"has_spoiler,omitempty"`
}

// InputMediaAudio describes an audio item for sendMediaGroup.
type InputMediaAudio struct {
	Type            string                   `json:"type"`
	Media           FileRef                  `json:"media"`
	Thumbnail       FileRef                  `json:"-"`
	Caption         string                   `json:"caption,omitempty"`
	ParseMode       string                   `json:"parse_mode,omitempty"`
	CaptionEntities []telegram.MessageEntity `json:"caption_entities,omitempty"`
	Duration        int                      `json:"duration,omitempty"`
	Performer       string                   `json:"performer,omitempty"`
	Title           string                   `json:"title,omitempty"`
}

// InputMediaDocument describes a document item for sendMediaGroup.
type InputMediaDocument struct {
	Type                        string                   `json:"type"`
	Media                       FileRef                  `json:"media"`
	Thumbnail                   FileRef                  `json:"-"`
	Caption                     string                   `json:"caption,omitempty"`
	ParseMode                   string                   `json:"parse_mode,omitempty"`
	CaptionEntities             []telegram.MessageEntity `json:"caption_entities,omitempty"`
	DisableContentTypeDetection bool                     `json:"disable_content_type_detection,omitempty"`
}

func (InputMediaPhoto) inputMedia()    {}
func (InputMediaVideo) inputMedia()    {}
func (InputMediaAudio) inputMedia()    {}
func (InputMediaDocument) inputMedia() {}

// MediaPhoto creates a photo input media item.
func MediaPhoto(media FileRef) InputMediaPhoto {
	return InputMediaPhoto{Type: "photo", Media: media}
}

// MediaVideo creates a video input media item.
func MediaVideo(media FileRef) InputMediaVideo {
	return InputMediaVideo{Type: "video", Media: media}
}

// MediaAudio creates an audio input media item.
func MediaAudio(media FileRef) InputMediaAudio {
	return InputMediaAudio{Type: "audio", Media: media}
}

// MediaDocument creates a document input media item.
func MediaDocument(media FileRef) InputMediaDocument {
	return InputMediaDocument{Type: "document", Media: media}
}

// SendMediaGroupParams contains supported parameters for sendMediaGroup.
type SendMediaGroupParams struct {
	BusinessConnectionID string                    `json:"business_connection_id,omitempty"`
	ChatID               ChatID                    `json:"chat_id"`
	MessageThreadID      int64                     `json:"message_thread_id,omitempty"`
	Media                []InputMedia              `json:"media"`
	DisableNotification  bool                      `json:"disable_notification,omitempty"`
	ProtectContent       bool                      `json:"protect_content,omitempty"`
	ReplyParameters      *telegram.ReplyParameters `json:"reply_parameters,omitempty"`
}

// SendMediaGroup sends an album of photos, videos, documents, or audio files.
func (b *Bot) SendMediaGroup(ctx context.Context, params SendMediaGroupParams) ([]telegram.Message, error) {
	if err := params.validate(); err != nil {
		return nil, err
	}

	media, files, err := params.mediaPayload()
	if err != nil {
		return nil, err
	}

	var messages []telegram.Message
	if len(files) > 0 {
		fields, err := params.multipartFields(media)
		if err != nil {
			return nil, err
		}
		if err := b.callMultipart(ctx, "sendMediaGroup", fields, files, &messages); err != nil {
			return nil, err
		}
		return messages, nil
	}

	payload := sendMediaGroupPayload{
		BusinessConnectionID: params.BusinessConnectionID,
		ChatID:               params.ChatID,
		MessageThreadID:      params.MessageThreadID,
		Media:                media,
		DisableNotification:  params.DisableNotification,
		ProtectContent:       params.ProtectContent,
		ReplyParameters:      params.ReplyParameters,
	}
	if err := b.call(ctx, "sendMediaGroup", payload, &messages); err != nil {
		return nil, err
	}

	return messages, nil
}

type sendMediaGroupPayload struct {
	BusinessConnectionID string                    `json:"business_connection_id,omitempty"`
	ChatID               ChatID                    `json:"chat_id"`
	MessageThreadID      int64                     `json:"message_thread_id,omitempty"`
	Media                []inputMediaPayload       `json:"media"`
	DisableNotification  bool                      `json:"disable_notification,omitempty"`
	ProtectContent       bool                      `json:"protect_content,omitempty"`
	ReplyParameters      *telegram.ReplyParameters `json:"reply_parameters,omitempty"`
}

type inputMediaPayload struct {
	Type                        string                   `json:"type"`
	Media                       string                   `json:"media"`
	Thumbnail                   string                   `json:"thumbnail,omitempty"`
	Caption                     string                   `json:"caption,omitempty"`
	ParseMode                   string                   `json:"parse_mode,omitempty"`
	CaptionEntities             []telegram.MessageEntity `json:"caption_entities,omitempty"`
	HasSpoiler                  bool                     `json:"has_spoiler,omitempty"`
	Width                       int                      `json:"width,omitempty"`
	Height                      int                      `json:"height,omitempty"`
	Duration                    int                      `json:"duration,omitempty"`
	SupportsStreaming           bool                     `json:"supports_streaming,omitempty"`
	Performer                   string                   `json:"performer,omitempty"`
	Title                       string                   `json:"title,omitempty"`
	DisableContentTypeDetection bool                     `json:"disable_content_type_detection,omitempty"`
}

func (params SendMediaGroupParams) validate() error {
	if !params.ChatID.valid() {
		return stderrors.New("chat_id is required")
	}
	if err := validateMessageThreadID(params.MessageThreadID); err != nil {
		return err
	}
	if len(params.Media) < 2 {
		return stderrors.New("media must contain at least two items")
	}
	if len(params.Media) > 10 {
		return stderrors.New("media must contain at most ten items")
	}
	for i, media := range params.Media {
		if err := validateInputMedia(media); err != nil {
			return fmt.Errorf("media[%d]: %w", i, err)
		}
	}
	if err := validateReplyParameters(params.ReplyParameters); err != nil {
		return err
	}

	return nil
}

func validateInputMedia(media InputMedia) error {
	switch item := media.(type) {
	case nil:
		return stderrors.New("media item is required")
	case InputMediaPhoto:
		return validateInputMediaPhoto(item)
	case *InputMediaPhoto:
		if item == nil {
			return stderrors.New("media item is required")
		}
		return validateInputMediaPhoto(*item)
	case InputMediaVideo:
		return validateInputMediaVideo(item)
	case *InputMediaVideo:
		if item == nil {
			return stderrors.New("media item is required")
		}
		return validateInputMediaVideo(*item)
	case InputMediaAudio:
		return validateInputMediaAudio(item)
	case *InputMediaAudio:
		if item == nil {
			return stderrors.New("media item is required")
		}
		return validateInputMediaAudio(*item)
	case InputMediaDocument:
		return validateInputMediaDocument(item)
	case *InputMediaDocument:
		if item == nil {
			return stderrors.New("media item is required")
		}
		return validateInputMediaDocument(*item)
	default:
		return stderrors.New("unsupported input media type")
	}
}

func validateInputMediaPhoto(media InputMediaPhoto) error {
	if err := validateInputMediaType(media.Type, "photo"); err != nil {
		return err
	}
	if err := media.Media.validate("media"); err != nil {
		return err
	}
	return validateCaptionFormatting(media.ParseMode, media.CaptionEntities)
}

func validateInputMediaVideo(media InputMediaVideo) error {
	if err := validateInputMediaType(media.Type, "video"); err != nil {
		return err
	}
	if err := media.Media.validate("media"); err != nil {
		return err
	}
	if err := validateOptionalFileRef(media.Thumbnail, "thumbnail"); err != nil {
		return err
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
	return validateCaptionFormatting(media.ParseMode, media.CaptionEntities)
}

func validateInputMediaAudio(media InputMediaAudio) error {
	if err := validateInputMediaType(media.Type, "audio"); err != nil {
		return err
	}
	if err := media.Media.validate("media"); err != nil {
		return err
	}
	if err := validateOptionalFileRef(media.Thumbnail, "thumbnail"); err != nil {
		return err
	}
	if media.Duration < 0 {
		return stderrors.New("duration must not be negative")
	}
	return validateCaptionFormatting(media.ParseMode, media.CaptionEntities)
}

func validateInputMediaDocument(media InputMediaDocument) error {
	if err := validateInputMediaType(media.Type, "document"); err != nil {
		return err
	}
	if err := media.Media.validate("media"); err != nil {
		return err
	}
	if err := validateOptionalFileRef(media.Thumbnail, "thumbnail"); err != nil {
		return err
	}
	return validateCaptionFormatting(media.ParseMode, media.CaptionEntities)
}

func validateInputMediaType(value string, expected string) error {
	if strings.TrimSpace(value) == "" || value == expected {
		return nil
	}
	return fmt.Errorf("type must be %q", expected)
}

func (params SendMediaGroupParams) mediaPayload() ([]inputMediaPayload, map[string]UploadFile, error) {
	payload := make([]inputMediaPayload, 0, len(params.Media))
	files := make(map[string]UploadFile)
	for i, media := range params.Media {
		item, err := buildInputMediaPayload(media, i, files)
		if err != nil {
			return nil, nil, fmt.Errorf("media[%d]: %w", i, err)
		}
		payload = append(payload, item)
	}

	return payload, files, nil
}

func buildInputMediaPayload(media InputMedia, index int, files map[string]UploadFile) (inputMediaPayload, error) {
	switch item := media.(type) {
	case InputMediaPhoto:
		return buildInputMediaPhotoPayload(item, index, files)
	case *InputMediaPhoto:
		if item == nil {
			return inputMediaPayload{}, stderrors.New("media item is required")
		}
		return buildInputMediaPhotoPayload(*item, index, files)
	case InputMediaVideo:
		return buildInputMediaVideoPayload(item, index, files)
	case *InputMediaVideo:
		if item == nil {
			return inputMediaPayload{}, stderrors.New("media item is required")
		}
		return buildInputMediaVideoPayload(*item, index, files)
	case InputMediaAudio:
		return buildInputMediaAudioPayload(item, index, files)
	case *InputMediaAudio:
		if item == nil {
			return inputMediaPayload{}, stderrors.New("media item is required")
		}
		return buildInputMediaAudioPayload(*item, index, files)
	case InputMediaDocument:
		return buildInputMediaDocumentPayload(item, index, files)
	case *InputMediaDocument:
		if item == nil {
			return inputMediaPayload{}, stderrors.New("media item is required")
		}
		return buildInputMediaDocumentPayload(*item, index, files)
	default:
		return inputMediaPayload{}, stderrors.New("unsupported input media type")
	}
}

func buildInputMediaPhotoPayload(media InputMediaPhoto, index int, files map[string]UploadFile) (inputMediaPayload, error) {
	mediaValue, err := mediaFileValue(media.Media, fmt.Sprintf("media%d", index), files)
	if err != nil {
		return inputMediaPayload{}, err
	}
	return inputMediaPayload{
		Type:            mediaType(media.Type, "photo"),
		Media:           mediaValue,
		Caption:         media.Caption,
		ParseMode:       media.ParseMode,
		CaptionEntities: media.CaptionEntities,
		HasSpoiler:      media.HasSpoiler,
	}, nil
}

func buildInputMediaVideoPayload(media InputMediaVideo, index int, files map[string]UploadFile) (inputMediaPayload, error) {
	mediaValue, err := mediaFileValue(media.Media, fmt.Sprintf("media%d", index), files)
	if err != nil {
		return inputMediaPayload{}, err
	}
	thumbnail, err := optionalMediaFileValue(media.Thumbnail, fmt.Sprintf("thumb%d", index), files)
	if err != nil {
		return inputMediaPayload{}, err
	}
	return inputMediaPayload{
		Type:              mediaType(media.Type, "video"),
		Media:             mediaValue,
		Thumbnail:         thumbnail,
		Caption:           media.Caption,
		ParseMode:         media.ParseMode,
		CaptionEntities:   media.CaptionEntities,
		Width:             media.Width,
		Height:            media.Height,
		Duration:          media.Duration,
		SupportsStreaming: media.SupportsStreaming,
		HasSpoiler:        media.HasSpoiler,
	}, nil
}

func buildInputMediaAudioPayload(media InputMediaAudio, index int, files map[string]UploadFile) (inputMediaPayload, error) {
	mediaValue, err := mediaFileValue(media.Media, fmt.Sprintf("media%d", index), files)
	if err != nil {
		return inputMediaPayload{}, err
	}
	thumbnail, err := optionalMediaFileValue(media.Thumbnail, fmt.Sprintf("thumb%d", index), files)
	if err != nil {
		return inputMediaPayload{}, err
	}
	return inputMediaPayload{
		Type:            mediaType(media.Type, "audio"),
		Media:           mediaValue,
		Thumbnail:       thumbnail,
		Caption:         media.Caption,
		ParseMode:       media.ParseMode,
		CaptionEntities: media.CaptionEntities,
		Duration:        media.Duration,
		Performer:       media.Performer,
		Title:           media.Title,
	}, nil
}

func buildInputMediaDocumentPayload(media InputMediaDocument, index int, files map[string]UploadFile) (inputMediaPayload, error) {
	mediaValue, err := mediaFileValue(media.Media, fmt.Sprintf("media%d", index), files)
	if err != nil {
		return inputMediaPayload{}, err
	}
	thumbnail, err := optionalMediaFileValue(media.Thumbnail, fmt.Sprintf("thumb%d", index), files)
	if err != nil {
		return inputMediaPayload{}, err
	}
	return inputMediaPayload{
		Type:                        mediaType(media.Type, "document"),
		Media:                       mediaValue,
		Thumbnail:                   thumbnail,
		Caption:                     media.Caption,
		ParseMode:                   media.ParseMode,
		CaptionEntities:             media.CaptionEntities,
		DisableContentTypeDetection: media.DisableContentTypeDetection,
	}, nil
}

func mediaType(value string, fallback string) string {
	if strings.TrimSpace(value) == "" {
		return fallback
	}
	return value
}

func mediaFileValue(ref FileRef, name string, files map[string]UploadFile) (string, error) {
	if err := ref.validate("media"); err != nil {
		return "", err
	}
	if ref.isUpload() {
		files[name] = ref.upload
		return "attach://" + name, nil
	}
	return ref.value, nil
}

func optionalMediaFileValue(ref FileRef, name string, files map[string]UploadFile) (string, error) {
	if !ref.isSet() {
		return "", nil
	}
	if err := ref.validate("thumbnail"); err != nil {
		return "", err
	}
	if ref.isUpload() {
		files[name] = ref.upload
		return "attach://" + name, nil
	}
	return ref.value, nil
}

func (params SendMediaGroupParams) multipartFields(media []inputMediaPayload) (map[string]string, error) {
	chatIDValue, err := params.ChatID.multipartValue()
	if err != nil {
		return nil, err
	}
	fields := map[string]string{"chat_id": chatIDValue}
	stringField(fields, "business_connection_id", params.BusinessConnectionID)
	int64Field(fields, "message_thread_id", params.MessageThreadID)
	boolField(fields, "disable_notification", params.DisableNotification)
	boolField(fields, "protect_content", params.ProtectContent)
	if err := replyParametersField(fields, params.ReplyParameters); err != nil {
		return nil, err
	}
	body, err := json.Marshal(media)
	if err != nil {
		return nil, err
	}
	fields["media"] = string(body)

	return fields, nil
}
