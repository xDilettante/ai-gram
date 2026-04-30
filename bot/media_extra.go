package bot

import (
	"context"
	stderrors "errors"
	"strings"

	"github.com/xDilettante/ai-gram/telegram"
)

// SendStickerParams contains supported parameters for sendSticker.
type SendStickerParams struct {
	ChatID              ChatID                    `json:"chat_id"`
	MessageThreadID     int64                     `json:"message_thread_id,omitempty"`
	Sticker             FileRef                   `json:"sticker"`
	Emoji               string                    `json:"emoji,omitempty"`
	DisableNotification bool                      `json:"disable_notification,omitempty"`
	ProtectContent      bool                      `json:"protect_content,omitempty"`
	ReplyParameters     *telegram.ReplyParameters `json:"reply_parameters,omitempty"`
	ReplyMarkup         telegram.ReplyMarkup      `json:"reply_markup,omitempty"`
}

// SendAnimationParams contains supported parameters for sendAnimation.
type SendAnimationParams struct {
	ChatID                ChatID                    `json:"chat_id"`
	MessageThreadID       int64                     `json:"message_thread_id,omitempty"`
	Animation             FileRef                   `json:"animation"`
	Duration              int                       `json:"duration,omitempty"`
	Width                 int                       `json:"width,omitempty"`
	Height                int                       `json:"height,omitempty"`
	Thumbnail             FileRef                   `json:"-"`
	Caption               string                    `json:"caption,omitempty"`
	ParseMode             string                    `json:"parse_mode,omitempty"`
	CaptionEntities       []telegram.MessageEntity  `json:"caption_entities,omitempty"`
	ShowCaptionAboveMedia bool                      `json:"show_caption_above_media,omitempty"`
	HasSpoiler            bool                      `json:"has_spoiler,omitempty"`
	DisableNotification   bool                      `json:"disable_notification,omitempty"`
	ProtectContent        bool                      `json:"protect_content,omitempty"`
	ReplyParameters       *telegram.ReplyParameters `json:"reply_parameters,omitempty"`
	ReplyMarkup           telegram.ReplyMarkup      `json:"reply_markup,omitempty"`
}

// SendVideoNoteParams contains supported parameters for sendVideoNote.
type SendVideoNoteParams struct {
	ChatID              ChatID                    `json:"chat_id"`
	MessageThreadID     int64                     `json:"message_thread_id,omitempty"`
	VideoNote           FileRef                   `json:"video_note"`
	Duration            int                       `json:"duration,omitempty"`
	Length              int                       `json:"length,omitempty"`
	Thumbnail           FileRef                   `json:"-"`
	DisableNotification bool                      `json:"disable_notification,omitempty"`
	ProtectContent      bool                      `json:"protect_content,omitempty"`
	ReplyParameters     *telegram.ReplyParameters `json:"reply_parameters,omitempty"`
	ReplyMarkup         telegram.ReplyMarkup      `json:"reply_markup,omitempty"`
}

// SendSticker sends a sticker by Telegram file_id, HTTP(S) URL, or multipart upload.
func (b *Bot) SendSticker(ctx context.Context, params SendStickerParams) (*telegram.Message, error) {
	if err := params.validate(); err != nil {
		return nil, err
	}

	var message telegram.Message
	if params.Sticker.isUpload() {
		fields, files, err := params.multipart()
		if err != nil {
			return nil, err
		}
		if err := b.callMultipart(ctx, "sendSticker", fields, files, &message); err != nil {
			return nil, err
		}
		return &message, nil
	}

	if err := b.call(ctx, "sendSticker", params, &message); err != nil {
		return nil, err
	}

	return &message, nil
}

// SendAnimation sends an animation by Telegram file_id, HTTP(S) URL, or multipart upload.
func (b *Bot) SendAnimation(ctx context.Context, params SendAnimationParams) (*telegram.Message, error) {
	if err := params.validate(); err != nil {
		return nil, err
	}

	var message telegram.Message
	if params.requiresMultipart() {
		fields, files, err := params.multipart()
		if err != nil {
			return nil, err
		}
		if err := b.callMultipart(ctx, "sendAnimation", fields, files, &message); err != nil {
			return nil, err
		}
		return &message, nil
	}

	if err := b.call(ctx, "sendAnimation", params.payload(), &message); err != nil {
		return nil, err
	}

	return &message, nil
}

// SendVideoNote sends a video note by Telegram file_id or multipart upload.
func (b *Bot) SendVideoNote(ctx context.Context, params SendVideoNoteParams) (*telegram.Message, error) {
	if err := params.validate(); err != nil {
		return nil, err
	}

	var message telegram.Message
	if params.requiresMultipart() {
		fields, files, err := params.multipart()
		if err != nil {
			return nil, err
		}
		if err := b.callMultipart(ctx, "sendVideoNote", fields, files, &message); err != nil {
			return nil, err
		}
		return &message, nil
	}

	if err := b.call(ctx, "sendVideoNote", params.payload(), &message); err != nil {
		return nil, err
	}

	return &message, nil
}

func (params SendStickerParams) validate() error {
	if !params.ChatID.valid() {
		return stderrors.New("chat_id is required")
	}
	if err := validateMessageThreadID(params.MessageThreadID); err != nil {
		return err
	}
	if err := params.Sticker.validate("sticker"); err != nil {
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

func (params SendAnimationParams) validate() error {
	if !params.ChatID.valid() {
		return stderrors.New("chat_id is required")
	}
	if err := validateMessageThreadID(params.MessageThreadID); err != nil {
		return err
	}
	if err := params.Animation.validate("animation"); err != nil {
		return err
	}
	if params.Duration < 0 {
		return stderrors.New("duration must not be negative")
	}
	if params.Width < 0 {
		return stderrors.New("width must not be negative")
	}
	if params.Height < 0 {
		return stderrors.New("height must not be negative")
	}
	if err := validateOptionalFileRef(params.Thumbnail, "thumbnail"); err != nil {
		return err
	}
	if err := validateCaptionFormatting(params.ParseMode, params.CaptionEntities); err != nil {
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

func (params SendVideoNoteParams) validate() error {
	if !params.ChatID.valid() {
		return stderrors.New("chat_id is required")
	}
	if err := validateMessageThreadID(params.MessageThreadID); err != nil {
		return err
	}
	if err := params.VideoNote.validate("video_note"); err != nil {
		return err
	}
	if params.VideoNote.kind == fileRefURL || strings.Contains(strings.TrimSpace(params.VideoNote.value), "://") {
		return stderrors.New("video_note must be a file_id or FileUpload")
	}
	if params.Duration < 0 {
		return stderrors.New("duration must not be negative")
	}
	if params.Length < 0 {
		return stderrors.New("length must not be negative")
	}
	if err := validateOptionalFileRef(params.Thumbnail, "thumbnail"); err != nil {
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

func (params SendStickerParams) multipart() (map[string]string, map[string]UploadFile, error) {
	fields, err := baseNonCaptionMediaFields(params.ChatID, params.MessageThreadID, params.DisableNotification, params.ProtectContent, params.ReplyParameters, params.ReplyMarkup)
	if err != nil {
		return nil, nil, err
	}
	stringField(fields, "emoji", params.Emoji)
	fields["sticker"] = "attach://sticker"
	return fields, map[string]UploadFile{"sticker": params.Sticker.upload}, nil
}

type sendAnimationPayload struct {
	ChatID                ChatID                    `json:"chat_id"`
	MessageThreadID       int64                     `json:"message_thread_id,omitempty"`
	Animation             FileRef                   `json:"animation"`
	Duration              int                       `json:"duration,omitempty"`
	Width                 int                       `json:"width,omitempty"`
	Height                int                       `json:"height,omitempty"`
	Thumbnail             *FileRef                  `json:"thumbnail,omitempty"`
	Caption               string                    `json:"caption,omitempty"`
	ParseMode             string                    `json:"parse_mode,omitempty"`
	CaptionEntities       []telegram.MessageEntity  `json:"caption_entities,omitempty"`
	ShowCaptionAboveMedia bool                      `json:"show_caption_above_media,omitempty"`
	HasSpoiler            bool                      `json:"has_spoiler,omitempty"`
	DisableNotification   bool                      `json:"disable_notification,omitempty"`
	ProtectContent        bool                      `json:"protect_content,omitempty"`
	ReplyParameters       *telegram.ReplyParameters `json:"reply_parameters,omitempty"`
	ReplyMarkup           telegram.ReplyMarkup      `json:"reply_markup,omitempty"`
}

func (params SendAnimationParams) payload() sendAnimationPayload {
	payload := sendAnimationPayload{
		ChatID:                params.ChatID,
		MessageThreadID:       params.MessageThreadID,
		Animation:             params.Animation,
		Duration:              params.Duration,
		Width:                 params.Width,
		Height:                params.Height,
		Caption:               params.Caption,
		ParseMode:             params.ParseMode,
		CaptionEntities:       params.CaptionEntities,
		ShowCaptionAboveMedia: params.ShowCaptionAboveMedia,
		HasSpoiler:            params.HasSpoiler,
		DisableNotification:   params.DisableNotification,
		ProtectContent:        params.ProtectContent,
		ReplyParameters:       params.ReplyParameters,
		ReplyMarkup:           params.ReplyMarkup,
	}
	if params.Thumbnail.isSet() {
		payload.Thumbnail = &params.Thumbnail
	}

	return payload
}

func (params SendAnimationParams) requiresMultipart() bool {
	return params.Animation.isUpload() || params.Thumbnail.isUpload()
}

func (params SendAnimationParams) multipart() (map[string]string, map[string]UploadFile, error) {
	fields, err := baseMediaFields(params.ChatID, params.MessageThreadID, params.Caption, params.ParseMode, params.CaptionEntities, params.DisableNotification, params.ProtectContent, params.ReplyParameters, params.ReplyMarkup)
	if err != nil {
		return nil, nil, err
	}
	intField(fields, "duration", params.Duration)
	intField(fields, "width", params.Width)
	intField(fields, "height", params.Height)
	boolField(fields, "show_caption_above_media", params.ShowCaptionAboveMedia)
	boolField(fields, "has_spoiler", params.HasSpoiler)
	files := make(map[string]UploadFile)
	if err := fileRefMultipartField(fields, files, "animation", params.Animation); err != nil {
		return nil, nil, err
	}
	if params.Thumbnail.isSet() {
		if err := fileRefMultipartField(fields, files, "thumbnail", params.Thumbnail); err != nil {
			return nil, nil, err
		}
	}

	return fields, files, nil
}

type sendVideoNotePayload struct {
	ChatID              ChatID                    `json:"chat_id"`
	MessageThreadID     int64                     `json:"message_thread_id,omitempty"`
	VideoNote           FileRef                   `json:"video_note"`
	Duration            int                       `json:"duration,omitempty"`
	Length              int                       `json:"length,omitempty"`
	Thumbnail           *FileRef                  `json:"thumbnail,omitempty"`
	DisableNotification bool                      `json:"disable_notification,omitempty"`
	ProtectContent      bool                      `json:"protect_content,omitempty"`
	ReplyParameters     *telegram.ReplyParameters `json:"reply_parameters,omitempty"`
	ReplyMarkup         telegram.ReplyMarkup      `json:"reply_markup,omitempty"`
}

func (params SendVideoNoteParams) payload() sendVideoNotePayload {
	payload := sendVideoNotePayload{
		ChatID:              params.ChatID,
		MessageThreadID:     params.MessageThreadID,
		VideoNote:           params.VideoNote,
		Duration:            params.Duration,
		Length:              params.Length,
		DisableNotification: params.DisableNotification,
		ProtectContent:      params.ProtectContent,
		ReplyParameters:     params.ReplyParameters,
		ReplyMarkup:         params.ReplyMarkup,
	}
	if params.Thumbnail.isSet() {
		payload.Thumbnail = &params.Thumbnail
	}

	return payload
}

func (params SendVideoNoteParams) requiresMultipart() bool {
	return params.VideoNote.isUpload() || params.Thumbnail.isUpload()
}

func (params SendVideoNoteParams) multipart() (map[string]string, map[string]UploadFile, error) {
	fields, err := baseNonCaptionMediaFields(params.ChatID, params.MessageThreadID, params.DisableNotification, params.ProtectContent, params.ReplyParameters, params.ReplyMarkup)
	if err != nil {
		return nil, nil, err
	}
	intField(fields, "duration", params.Duration)
	intField(fields, "length", params.Length)
	files := make(map[string]UploadFile)
	if err := fileRefMultipartField(fields, files, "video_note", params.VideoNote); err != nil {
		return nil, nil, err
	}
	if params.Thumbnail.isSet() {
		if err := fileRefMultipartField(fields, files, "thumbnail", params.Thumbnail); err != nil {
			return nil, nil, err
		}
	}

	return fields, files, nil
}

func baseNonCaptionMediaFields(chatID ChatID, messageThreadID int64, disableNotification bool, protectContent bool, replyParameters *telegram.ReplyParameters, replyMarkup telegram.ReplyMarkup) (map[string]string, error) {
	chatIDValue, err := chatID.multipartValue()
	if err != nil {
		return nil, err
	}
	fields := map[string]string{"chat_id": chatIDValue}
	int64Field(fields, "message_thread_id", messageThreadID)
	if err := replyParametersField(fields, replyParameters); err != nil {
		return nil, err
	}
	if err := replyMarkupField(fields, replyMarkup); err != nil {
		return nil, err
	}
	boolField(fields, "disable_notification", disableNotification)
	boolField(fields, "protect_content", protectContent)

	return fields, nil
}

func (ref FileRef) isSet() bool {
	return ref.kind != fileRefUnknown || ref.value != "" || ref.upload.Name != "" || ref.upload.Reader != nil
}

func validateOptionalFileRef(ref FileRef, field string) error {
	if !ref.isSet() {
		return nil
	}

	return ref.validate(field)
}

func fileRefMultipartField(fields map[string]string, files map[string]UploadFile, name string, ref FileRef) error {
	if ref.isUpload() {
		fields[name] = "attach://" + name
		files[name] = ref.upload
		return nil
	}
	if err := ref.validate(name); err != nil {
		return err
	}
	fields[name] = ref.value
	return nil
}
