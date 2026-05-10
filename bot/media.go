package bot

import (
	"context"
	"encoding/json"
	stderrors "errors"
	"io"
	"mime"
	"net/url"
	"strconv"
	"strings"

	"github.com/xDilettante/ai-gram/telegram"
)

type fileRefKind uint8

const (
	fileRefUnknown fileRefKind = iota
	fileRefID
	fileRefURL
	fileRefUpload
)

// UploadFile describes a file uploaded through multipart/form-data.
//
// Reader is consumed but not closed by ai-gram; callers keep ownership of the reader lifecycle.
type UploadFile struct {
	// Name is the multipart filename.
	Name string
	// Reader provides file contents and is not closed by the library.
	Reader io.Reader
	// ContentType is an optional media type for the file part.
	ContentType string
}

// FileRef identifies media by Telegram file_id, HTTP(S) URL, or multipart upload.
type FileRef struct {
	value  string
	kind   fileRefKind
	upload UploadFile
}

// FileID creates a file reference from an existing Telegram file_id.
func FileID(id string) FileRef {
	return FileRef{value: id, kind: fileRefID}
}

// FileURL creates a file reference from an HTTP(S) URL.
func FileURL(rawURL string) FileRef {
	return FileRef{value: rawURL, kind: fileRefURL}
}

// FileUpload creates a file reference from an UploadFile for multipart upload.
func FileUpload(file UploadFile) FileRef {
	return FileRef{kind: fileRefUpload, upload: file}
}

// MarshalJSON encodes ref as the Telegram Bot API string value for file_id or URL references.
func (ref FileRef) MarshalJSON() ([]byte, error) {
	if ref.kind == fileRefUpload {
		return nil, stderrors.New("upload file cannot be used in JSON request")
	}

	return json.Marshal(ref.value)
}

// SendPhotoParams contains supported parameters for sendPhoto.
type SendPhotoParams struct {
	BusinessConnectionID string                    `json:"business_connection_id,omitempty"`
	ChatID               ChatID                    `json:"chat_id"`
	MessageThreadID      int64                     `json:"message_thread_id,omitempty"`
	Photo                FileRef                   `json:"photo"`
	Caption              string                    `json:"caption,omitempty"`
	ParseMode            string                    `json:"parse_mode,omitempty"`
	CaptionEntities      []telegram.MessageEntity  `json:"caption_entities,omitempty"`
	DisableNotification  bool                      `json:"disable_notification,omitempty"`
	ProtectContent       bool                      `json:"protect_content,omitempty"`
	ReplyParameters      *telegram.ReplyParameters `json:"reply_parameters,omitempty"`
	ReplyMarkup          telegram.ReplyMarkup      `json:"reply_markup,omitempty"`
}

// SendDocumentParams contains supported parameters for sendDocument.
type SendDocumentParams struct {
	BusinessConnectionID        string                    `json:"business_connection_id,omitempty"`
	ChatID                      ChatID                    `json:"chat_id"`
	MessageThreadID             int64                     `json:"message_thread_id,omitempty"`
	Document                    FileRef                   `json:"document"`
	Caption                     string                    `json:"caption,omitempty"`
	ParseMode                   string                    `json:"parse_mode,omitempty"`
	CaptionEntities             []telegram.MessageEntity  `json:"caption_entities,omitempty"`
	DisableNotification         bool                      `json:"disable_notification,omitempty"`
	ProtectContent              bool                      `json:"protect_content,omitempty"`
	DisableContentTypeDetection bool                      `json:"disable_content_type_detection,omitempty"`
	ReplyParameters             *telegram.ReplyParameters `json:"reply_parameters,omitempty"`
	ReplyMarkup                 telegram.ReplyMarkup      `json:"reply_markup,omitempty"`
}

// SendVideoParams contains supported parameters for sendVideo.
type SendVideoParams struct {
	BusinessConnectionID  string                    `json:"business_connection_id,omitempty"`
	ChatID                ChatID                    `json:"chat_id"`
	MessageThreadID       int64                     `json:"message_thread_id,omitempty"`
	Video                 FileRef                   `json:"video"`
	Thumbnail             FileRef                   `json:"-"`
	Cover                 FileRef                   `json:"-"`
	StartTimestamp        int                       `json:"start_timestamp,omitempty"`
	Duration              int                       `json:"duration,omitempty"`
	Width                 int                       `json:"width,omitempty"`
	Height                int                       `json:"height,omitempty"`
	Caption               string                    `json:"caption,omitempty"`
	ParseMode             string                    `json:"parse_mode,omitempty"`
	CaptionEntities       []telegram.MessageEntity  `json:"caption_entities,omitempty"`
	ShowCaptionAboveMedia bool                      `json:"show_caption_above_media,omitempty"`
	SupportsStreaming     bool                      `json:"supports_streaming,omitempty"`
	HasSpoiler            bool                      `json:"has_spoiler,omitempty"`
	DisableNotification   bool                      `json:"disable_notification,omitempty"`
	ProtectContent        bool                      `json:"protect_content,omitempty"`
	ReplyParameters       *telegram.ReplyParameters `json:"reply_parameters,omitempty"`
	ReplyMarkup           telegram.ReplyMarkup      `json:"reply_markup,omitempty"`
}

// SendLivePhotoParams contains supported parameters for sendLivePhoto.
type SendLivePhotoParams struct {
	BusinessConnectionID  string                    `json:"business_connection_id,omitempty"`
	ChatID                ChatID                    `json:"chat_id"`
	MessageThreadID       int64                     `json:"message_thread_id,omitempty"`
	DirectMessagesTopicID int64                     `json:"direct_messages_topic_id,omitempty"`
	LivePhoto             FileRef                   `json:"live_photo"`
	Photo                 FileRef                   `json:"photo"`
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

// SendAudioParams contains supported parameters for sendAudio.
type SendAudioParams struct {
	BusinessConnectionID string                    `json:"business_connection_id,omitempty"`
	ChatID               ChatID                    `json:"chat_id"`
	MessageThreadID      int64                     `json:"message_thread_id,omitempty"`
	Audio                FileRef                   `json:"audio"`
	Duration             int                       `json:"duration,omitempty"`
	Performer            string                    `json:"performer,omitempty"`
	Title                string                    `json:"title,omitempty"`
	Caption              string                    `json:"caption,omitempty"`
	ParseMode            string                    `json:"parse_mode,omitempty"`
	CaptionEntities      []telegram.MessageEntity  `json:"caption_entities,omitempty"`
	DisableNotification  bool                      `json:"disable_notification,omitempty"`
	ProtectContent       bool                      `json:"protect_content,omitempty"`
	ReplyParameters      *telegram.ReplyParameters `json:"reply_parameters,omitempty"`
	ReplyMarkup          telegram.ReplyMarkup      `json:"reply_markup,omitempty"`
}

// SendVoiceParams contains supported parameters for sendVoice.
type SendVoiceParams struct {
	BusinessConnectionID string                    `json:"business_connection_id,omitempty"`
	ChatID               ChatID                    `json:"chat_id"`
	MessageThreadID      int64                     `json:"message_thread_id,omitempty"`
	Voice                FileRef                   `json:"voice"`
	Caption              string                    `json:"caption,omitempty"`
	ParseMode            string                    `json:"parse_mode,omitempty"`
	CaptionEntities      []telegram.MessageEntity  `json:"caption_entities,omitempty"`
	Duration             int                       `json:"duration,omitempty"`
	DisableNotification  bool                      `json:"disable_notification,omitempty"`
	ProtectContent       bool                      `json:"protect_content,omitempty"`
	ReplyParameters      *telegram.ReplyParameters `json:"reply_parameters,omitempty"`
	ReplyMarkup          telegram.ReplyMarkup      `json:"reply_markup,omitempty"`
}

// SendPhoto sends a photo by Telegram file_id, HTTP(S) URL, or multipart upload.
func (b *Bot) SendPhoto(ctx context.Context, params SendPhotoParams) (*telegram.Message, error) {
	if err := params.validate(); err != nil {
		return nil, err
	}

	var message telegram.Message
	if params.Photo.isUpload() {
		fields, files, err := params.multipart()
		if err != nil {
			return nil, err
		}
		if err := b.callMultipart(ctx, "sendPhoto", fields, files, &message); err != nil {
			return nil, err
		}
		return &message, nil
	}

	if err := b.call(ctx, "sendPhoto", params, &message); err != nil {
		return nil, err
	}

	return &message, nil
}

// SendDocument sends a document by Telegram file_id, HTTP(S) URL, or multipart upload.
func (b *Bot) SendDocument(ctx context.Context, params SendDocumentParams) (*telegram.Message, error) {
	if err := params.validate(); err != nil {
		return nil, err
	}

	var message telegram.Message
	if params.Document.isUpload() {
		fields, files, err := params.multipart()
		if err != nil {
			return nil, err
		}
		if err := b.callMultipart(ctx, "sendDocument", fields, files, &message); err != nil {
			return nil, err
		}
		return &message, nil
	}

	if err := b.call(ctx, "sendDocument", params, &message); err != nil {
		return nil, err
	}

	return &message, nil
}

// SendVideo sends a video by Telegram file_id, HTTP(S) URL, or multipart upload.
func (b *Bot) SendVideo(ctx context.Context, params SendVideoParams) (*telegram.Message, error) {
	if err := params.validate(); err != nil {
		return nil, err
	}

	var message telegram.Message
	if params.requiresMultipart() {
		fields, files, err := params.multipart()
		if err != nil {
			return nil, err
		}
		if err := b.callMultipart(ctx, "sendVideo", fields, files, &message); err != nil {
			return nil, err
		}
		return &message, nil
	}

	if err := b.call(ctx, "sendVideo", params.payload(), &message); err != nil {
		return nil, err
	}

	return &message, nil
}

// SendLivePhoto sends a live photo by Telegram file_id or multipart upload.
func (b *Bot) SendLivePhoto(ctx context.Context, params SendLivePhotoParams) (*telegram.Message, error) {
	if err := params.validate(); err != nil {
		return nil, err
	}

	var message telegram.Message
	if params.requiresMultipart() {
		fields, files, err := params.multipart()
		if err != nil {
			return nil, err
		}
		if err := b.callMultipart(ctx, "sendLivePhoto", fields, files, &message); err != nil {
			return nil, err
		}
		return &message, nil
	}

	if err := b.call(ctx, "sendLivePhoto", params, &message); err != nil {
		return nil, err
	}

	return &message, nil
}

// SendAudio sends an audio file by Telegram file_id, HTTP(S) URL, or multipart upload.
func (b *Bot) SendAudio(ctx context.Context, params SendAudioParams) (*telegram.Message, error) {
	if err := params.validate(); err != nil {
		return nil, err
	}

	var message telegram.Message
	if params.Audio.isUpload() {
		fields, files, err := params.multipart()
		if err != nil {
			return nil, err
		}
		if err := b.callMultipart(ctx, "sendAudio", fields, files, &message); err != nil {
			return nil, err
		}
		return &message, nil
	}

	if err := b.call(ctx, "sendAudio", params, &message); err != nil {
		return nil, err
	}

	return &message, nil
}

// SendVoice sends a voice message by Telegram file_id, HTTP(S) URL, or multipart upload.
func (b *Bot) SendVoice(ctx context.Context, params SendVoiceParams) (*telegram.Message, error) {
	if err := params.validate(); err != nil {
		return nil, err
	}

	var message telegram.Message
	if params.Voice.isUpload() {
		fields, files, err := params.multipart()
		if err != nil {
			return nil, err
		}
		if err := b.callMultipart(ctx, "sendVoice", fields, files, &message); err != nil {
			return nil, err
		}
		return &message, nil
	}

	if err := b.call(ctx, "sendVoice", params, &message); err != nil {
		return nil, err
	}

	return &message, nil
}

func (params SendPhotoParams) validate() error {
	if !params.ChatID.valid() {
		return stderrors.New("chat_id is required")
	}
	if err := validateMessageThreadID(params.MessageThreadID); err != nil {
		return err
	}
	if err := params.Photo.validate("photo"); err != nil {
		return err
	}
	if err := validateCaptionFormatting(params.ParseMode, params.CaptionEntities); err != nil {
		return err
	}
	if err := telegram.ValidateReplyMarkup(params.ReplyMarkup); err != nil {
		return err
	}
	if err := validateReplyParameters(params.ReplyParameters); err != nil {
		return err
	}

	return nil
}

func (params SendPhotoParams) multipart() (map[string]string, map[string]UploadFile, error) {
	fields, err := baseMediaFields(params.ChatID, params.MessageThreadID, params.Caption, params.ParseMode, params.CaptionEntities, params.DisableNotification, params.ProtectContent, params.ReplyParameters, params.ReplyMarkup)
	if err != nil {
		return nil, nil, err
	}
	stringField(fields, "business_connection_id", params.BusinessConnectionID)
	fields["photo"] = "attach://photo"
	return fields, map[string]UploadFile{"photo": params.Photo.upload}, nil
}

func (params SendDocumentParams) validate() error {
	if !params.ChatID.valid() {
		return stderrors.New("chat_id is required")
	}
	if err := validateMessageThreadID(params.MessageThreadID); err != nil {
		return err
	}
	if err := params.Document.validate("document"); err != nil {
		return err
	}
	if err := validateCaptionFormatting(params.ParseMode, params.CaptionEntities); err != nil {
		return err
	}
	if err := telegram.ValidateReplyMarkup(params.ReplyMarkup); err != nil {
		return err
	}
	if err := validateReplyParameters(params.ReplyParameters); err != nil {
		return err
	}

	return nil
}

func (params SendDocumentParams) multipart() (map[string]string, map[string]UploadFile, error) {
	fields, err := baseMediaFields(params.ChatID, params.MessageThreadID, params.Caption, params.ParseMode, params.CaptionEntities, params.DisableNotification, params.ProtectContent, params.ReplyParameters, params.ReplyMarkup)
	if err != nil {
		return nil, nil, err
	}
	stringField(fields, "business_connection_id", params.BusinessConnectionID)
	fields["document"] = "attach://document"
	if params.DisableContentTypeDetection {
		fields["disable_content_type_detection"] = "true"
	}
	return fields, map[string]UploadFile{"document": params.Document.upload}, nil
}

func (params SendVideoParams) validate() error {
	if !params.ChatID.valid() {
		return stderrors.New("chat_id is required")
	}
	if err := validateMessageThreadID(params.MessageThreadID); err != nil {
		return err
	}
	if err := params.Video.validate("video"); err != nil {
		return err
	}
	if err := validateOptionalFileRef(params.Thumbnail, "thumbnail"); err != nil {
		return err
	}
	if err := validateOptionalFileRef(params.Cover, "cover"); err != nil {
		return err
	}
	if params.StartTimestamp < 0 {
		return stderrors.New("start_timestamp must not be negative")
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
	if err := validateCaptionFormatting(params.ParseMode, params.CaptionEntities); err != nil {
		return err
	}
	if err := telegram.ValidateReplyMarkup(params.ReplyMarkup); err != nil {
		return err
	}
	if err := validateReplyParameters(params.ReplyParameters); err != nil {
		return err
	}

	return nil
}

type sendVideoPayload struct {
	BusinessConnectionID  string                    `json:"business_connection_id,omitempty"`
	ChatID                ChatID                    `json:"chat_id"`
	MessageThreadID       int64                     `json:"message_thread_id,omitempty"`
	Video                 FileRef                   `json:"video"`
	Thumbnail             *FileRef                  `json:"thumbnail,omitempty"`
	Cover                 *FileRef                  `json:"cover,omitempty"`
	StartTimestamp        int                       `json:"start_timestamp,omitempty"`
	Duration              int                       `json:"duration,omitempty"`
	Width                 int                       `json:"width,omitempty"`
	Height                int                       `json:"height,omitempty"`
	Caption               string                    `json:"caption,omitempty"`
	ParseMode             string                    `json:"parse_mode,omitempty"`
	CaptionEntities       []telegram.MessageEntity  `json:"caption_entities,omitempty"`
	ShowCaptionAboveMedia bool                      `json:"show_caption_above_media,omitempty"`
	SupportsStreaming     bool                      `json:"supports_streaming,omitempty"`
	HasSpoiler            bool                      `json:"has_spoiler,omitempty"`
	DisableNotification   bool                      `json:"disable_notification,omitempty"`
	ProtectContent        bool                      `json:"protect_content,omitempty"`
	ReplyParameters       *telegram.ReplyParameters `json:"reply_parameters,omitempty"`
	ReplyMarkup           telegram.ReplyMarkup      `json:"reply_markup,omitempty"`
}

func (params SendVideoParams) payload() sendVideoPayload {
	payload := sendVideoPayload{
		BusinessConnectionID:  params.BusinessConnectionID,
		ChatID:                params.ChatID,
		MessageThreadID:       params.MessageThreadID,
		Video:                 params.Video,
		StartTimestamp:        params.StartTimestamp,
		Duration:              params.Duration,
		Width:                 params.Width,
		Height:                params.Height,
		Caption:               params.Caption,
		ParseMode:             params.ParseMode,
		CaptionEntities:       params.CaptionEntities,
		ShowCaptionAboveMedia: params.ShowCaptionAboveMedia,
		SupportsStreaming:     params.SupportsStreaming,
		HasSpoiler:            params.HasSpoiler,
		DisableNotification:   params.DisableNotification,
		ProtectContent:        params.ProtectContent,
		ReplyParameters:       params.ReplyParameters,
		ReplyMarkup:           params.ReplyMarkup,
	}
	if params.Thumbnail.isSet() {
		payload.Thumbnail = &params.Thumbnail
	}
	if params.Cover.isSet() {
		payload.Cover = &params.Cover
	}
	return payload
}

func (params SendVideoParams) requiresMultipart() bool {
	return params.Video.isUpload() || params.Thumbnail.isUpload() || params.Cover.isUpload()
}

func (params SendVideoParams) multipart() (map[string]string, map[string]UploadFile, error) {
	fields, err := baseMediaFields(params.ChatID, params.MessageThreadID, params.Caption, params.ParseMode, params.CaptionEntities, params.DisableNotification, params.ProtectContent, params.ReplyParameters, params.ReplyMarkup)
	if err != nil {
		return nil, nil, err
	}
	stringField(fields, "business_connection_id", params.BusinessConnectionID)
	files := make(map[string]UploadFile)
	if err := fileRefMultipartField(fields, files, "video", params.Video); err != nil {
		return nil, nil, err
	}
	if params.Thumbnail.isSet() {
		if err := fileRefMultipartField(fields, files, "thumbnail", params.Thumbnail); err != nil {
			return nil, nil, err
		}
	}
	if params.Cover.isSet() {
		if err := fileRefMultipartField(fields, files, "cover", params.Cover); err != nil {
			return nil, nil, err
		}
	}
	intField(fields, "start_timestamp", params.StartTimestamp)
	intField(fields, "duration", params.Duration)
	intField(fields, "width", params.Width)
	intField(fields, "height", params.Height)
	boolField(fields, "show_caption_above_media", params.ShowCaptionAboveMedia)
	boolField(fields, "supports_streaming", params.SupportsStreaming)
	boolField(fields, "has_spoiler", params.HasSpoiler)
	return fields, files, nil
}

func (params SendLivePhotoParams) validate() error {
	if !params.ChatID.valid() {
		return stderrors.New("chat_id is required")
	}
	if err := validateMessageThreadID(params.MessageThreadID); err != nil {
		return err
	}
	if params.DirectMessagesTopicID < 0 {
		return stderrors.New("direct_messages_topic_id must not be negative")
	}
	if err := validateLivePhotoFileRef(params.LivePhoto, "live_photo"); err != nil {
		return err
	}
	if err := validateLivePhotoFileRef(params.Photo, "photo"); err != nil {
		return err
	}
	if err := validateCaptionFormatting(params.ParseMode, params.CaptionEntities); err != nil {
		return err
	}
	if err := telegram.ValidateReplyMarkup(params.ReplyMarkup); err != nil {
		return err
	}
	if err := validateReplyParameters(params.ReplyParameters); err != nil {
		return err
	}

	return nil
}

func (params SendLivePhotoParams) requiresMultipart() bool {
	return params.LivePhoto.isUpload() || params.Photo.isUpload()
}

func (params SendLivePhotoParams) multipart() (map[string]string, map[string]UploadFile, error) {
	fields, err := baseMediaFields(params.ChatID, params.MessageThreadID, params.Caption, params.ParseMode, params.CaptionEntities, params.DisableNotification, params.ProtectContent, params.ReplyParameters, params.ReplyMarkup)
	if err != nil {
		return nil, nil, err
	}
	stringField(fields, "business_connection_id", params.BusinessConnectionID)
	int64Field(fields, "direct_messages_topic_id", params.DirectMessagesTopicID)
	boolField(fields, "show_caption_above_media", params.ShowCaptionAboveMedia)
	boolField(fields, "has_spoiler", params.HasSpoiler)
	files := make(map[string]UploadFile)
	if err := fileRefMultipartField(fields, files, "live_photo", params.LivePhoto); err != nil {
		return nil, nil, err
	}
	if err := fileRefMultipartField(fields, files, "photo", params.Photo); err != nil {
		return nil, nil, err
	}
	return fields, files, nil
}

func validateLivePhotoFileRef(ref FileRef, field string) error {
	if err := ref.validate(field); err != nil {
		return err
	}
	if ref.kind == fileRefURL || strings.Contains(ref.value, "://") {
		return stderrors.New(field + " must be a file_id or FileUpload; URLs are not supported")
	}
	return nil
}

func (params SendAudioParams) validate() error {
	if !params.ChatID.valid() {
		return stderrors.New("chat_id is required")
	}
	if err := validateMessageThreadID(params.MessageThreadID); err != nil {
		return err
	}
	if err := params.Audio.validate("audio"); err != nil {
		return err
	}
	if params.Duration < 0 {
		return stderrors.New("duration must not be negative")
	}
	if err := validateCaptionFormatting(params.ParseMode, params.CaptionEntities); err != nil {
		return err
	}
	if err := telegram.ValidateReplyMarkup(params.ReplyMarkup); err != nil {
		return err
	}
	if err := validateReplyParameters(params.ReplyParameters); err != nil {
		return err
	}

	return nil
}

func (params SendAudioParams) multipart() (map[string]string, map[string]UploadFile, error) {
	fields, err := baseMediaFields(params.ChatID, params.MessageThreadID, params.Caption, params.ParseMode, params.CaptionEntities, params.DisableNotification, params.ProtectContent, params.ReplyParameters, params.ReplyMarkup)
	if err != nil {
		return nil, nil, err
	}
	stringField(fields, "business_connection_id", params.BusinessConnectionID)
	fields["audio"] = "attach://audio"
	intField(fields, "duration", params.Duration)
	stringField(fields, "performer", params.Performer)
	stringField(fields, "title", params.Title)
	return fields, map[string]UploadFile{"audio": params.Audio.upload}, nil
}

func (params SendVoiceParams) validate() error {
	if !params.ChatID.valid() {
		return stderrors.New("chat_id is required")
	}
	if err := validateMessageThreadID(params.MessageThreadID); err != nil {
		return err
	}
	if err := params.Voice.validate("voice"); err != nil {
		return err
	}
	if params.Duration < 0 {
		return stderrors.New("duration must not be negative")
	}
	if err := validateCaptionFormatting(params.ParseMode, params.CaptionEntities); err != nil {
		return err
	}
	if err := telegram.ValidateReplyMarkup(params.ReplyMarkup); err != nil {
		return err
	}
	if err := validateReplyParameters(params.ReplyParameters); err != nil {
		return err
	}

	return nil
}

func (params SendVoiceParams) multipart() (map[string]string, map[string]UploadFile, error) {
	fields, err := baseMediaFields(params.ChatID, params.MessageThreadID, params.Caption, params.ParseMode, params.CaptionEntities, params.DisableNotification, params.ProtectContent, params.ReplyParameters, params.ReplyMarkup)
	if err != nil {
		return nil, nil, err
	}
	stringField(fields, "business_connection_id", params.BusinessConnectionID)
	fields["voice"] = "attach://voice"
	intField(fields, "duration", params.Duration)
	return fields, map[string]UploadFile{"voice": params.Voice.upload}, nil
}

func baseMediaFields(chatID ChatID, messageThreadID int64, caption string, parseMode string, captionEntities []telegram.MessageEntity, disableNotification bool, protectContent bool, replyParameters *telegram.ReplyParameters, replyMarkup telegram.ReplyMarkup) (map[string]string, error) {
	chatIDValue, err := chatID.multipartValue()
	if err != nil {
		return nil, err
	}
	fields := map[string]string{"chat_id": chatIDValue}
	int64Field(fields, "message_thread_id", messageThreadID)
	stringField(fields, "caption", caption)
	stringField(fields, "parse_mode", parseMode)
	if err := captionEntitiesField(fields, captionEntities); err != nil {
		return nil, err
	}
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

func validateCaptionFormatting(parseMode string, captionEntities []telegram.MessageEntity) error {
	if parseMode != "" && len(captionEntities) > 0 {
		return stderrors.New("parse_mode and caption_entities cannot be used together")
	}

	return nil
}

func boolField(fields map[string]string, name string, value bool) {
	if value {
		fields[name] = "true"
	}
}

func intField(fields map[string]string, name string, value int) {
	if value > 0 {
		fields[name] = strconv.Itoa(value)
	}
}

func int64Field(fields map[string]string, name string, value int64) {
	if value > 0 {
		fields[name] = strconv.FormatInt(value, 10)
	}
}

func stringField(fields map[string]string, name string, value string) {
	if value != "" {
		fields[name] = value
	}
}

func captionEntitiesField(fields map[string]string, captionEntities []telegram.MessageEntity) error {
	if len(captionEntities) == 0 {
		return nil
	}
	body, err := json.Marshal(captionEntities)
	if err != nil {
		return err
	}
	fields["caption_entities"] = string(body)
	return nil
}

func replyParametersField(fields map[string]string, replyParameters *telegram.ReplyParameters) error {
	if replyParameters == nil {
		return nil
	}
	body, err := json.Marshal(replyParameters)
	if err != nil {
		return err
	}
	fields["reply_parameters"] = string(body)
	return nil
}

func replyMarkupField(fields map[string]string, replyMarkup telegram.ReplyMarkup) error {
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

func (ref FileRef) isUpload() bool {
	return ref.kind == fileRefUpload
}

func (ref FileRef) validate(field string) error {
	if ref.kind == fileRefUpload {
		return ref.upload.validate(field)
	}

	value := strings.TrimSpace(ref.value)
	if value == "" {
		return stderrors.New(field + " is required")
	}
	if strings.HasPrefix(value, "attach://") {
		return stderrors.New(field + " must use FileUpload instead of attach URL")
	}
	if strings.HasPrefix(value, "/") {
		return stderrors.New(field + " must be a file_id or HTTP(S) URL, not a local path")
	}
	if ref.kind == fileRefURL || strings.Contains(value, "://") {
		parsed, err := url.Parse(value)
		if err != nil {
			return stderrors.New(field + " must be a file_id or valid HTTP(S) URL")
		}
		if parsed.Scheme != "http" && parsed.Scheme != "https" {
			return stderrors.New(field + " URL scheme must be http or https")
		}
		if parsed.Host == "" {
			return stderrors.New(field + " URL host is required")
		}
	}

	return nil
}

func (file UploadFile) validate(field string) error {
	if strings.TrimSpace(file.Name) == "" {
		return stderrors.New(field + " upload filename is required")
	}
	if file.Reader == nil {
		return stderrors.New(field + " upload reader is required")
	}
	if strings.ContainsAny(file.Name, "/\\\x00") || strings.Contains(file.Name, "..") {
		return stderrors.New(field + " upload filename must not contain path separators or traversal")
	}
	if strings.TrimSpace(file.ContentType) != "" {
		if _, _, err := mime.ParseMediaType(file.ContentType); err != nil {
			return stderrors.New(field + " upload content type is invalid")
		}
	}

	return nil
}
