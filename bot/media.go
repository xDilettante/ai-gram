package bot

import (
	"context"
	"encoding/json"
	stderrors "errors"
	"io"
	"mime"
	"net/url"
	"strings"

	"ai-gram/telegram"
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
	ChatID              ChatID                   `json:"chat_id"`
	Photo               FileRef                  `json:"photo"`
	Caption             string                   `json:"caption,omitempty"`
	ParseMode           string                   `json:"parse_mode,omitempty"`
	CaptionEntities     []telegram.MessageEntity `json:"caption_entities,omitempty"`
	DisableNotification bool                     `json:"disable_notification,omitempty"`
	ProtectContent      bool                     `json:"protect_content,omitempty"`
}

// SendDocumentParams contains supported parameters for sendDocument.
type SendDocumentParams struct {
	ChatID                      ChatID                   `json:"chat_id"`
	Document                    FileRef                  `json:"document"`
	Caption                     string                   `json:"caption,omitempty"`
	ParseMode                   string                   `json:"parse_mode,omitempty"`
	CaptionEntities             []telegram.MessageEntity `json:"caption_entities,omitempty"`
	DisableNotification         bool                     `json:"disable_notification,omitempty"`
	ProtectContent              bool                     `json:"protect_content,omitempty"`
	DisableContentTypeDetection bool                     `json:"disable_content_type_detection,omitempty"`
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

func (params SendPhotoParams) validate() error {
	if !params.ChatID.valid() {
		return stderrors.New("chat_id is required")
	}
	if err := params.Photo.validate("photo"); err != nil {
		return err
	}
	if params.ParseMode != "" && len(params.CaptionEntities) > 0 {
		return stderrors.New("parse_mode and caption_entities cannot be used together")
	}

	return nil
}

func (params SendPhotoParams) multipart() (map[string]string, map[string]UploadFile, error) {
	fields, err := baseMediaFields(params.ChatID, params.Caption, params.ParseMode, params.CaptionEntities, params.DisableNotification, params.ProtectContent)
	if err != nil {
		return nil, nil, err
	}
	fields["photo"] = "attach://photo"
	return fields, map[string]UploadFile{"photo": params.Photo.upload}, nil
}

func (params SendDocumentParams) validate() error {
	if !params.ChatID.valid() {
		return stderrors.New("chat_id is required")
	}
	if err := params.Document.validate("document"); err != nil {
		return err
	}
	if params.ParseMode != "" && len(params.CaptionEntities) > 0 {
		return stderrors.New("parse_mode and caption_entities cannot be used together")
	}

	return nil
}

func (params SendDocumentParams) multipart() (map[string]string, map[string]UploadFile, error) {
	fields, err := baseMediaFields(params.ChatID, params.Caption, params.ParseMode, params.CaptionEntities, params.DisableNotification, params.ProtectContent)
	if err != nil {
		return nil, nil, err
	}
	fields["document"] = "attach://document"
	if params.DisableContentTypeDetection {
		fields["disable_content_type_detection"] = "true"
	}
	return fields, map[string]UploadFile{"document": params.Document.upload}, nil
}

func baseMediaFields(chatID ChatID, caption string, parseMode string, captionEntities []telegram.MessageEntity, disableNotification bool, protectContent bool) (map[string]string, error) {
	chatIDValue, err := chatID.multipartValue()
	if err != nil {
		return nil, err
	}
	fields := map[string]string{"chat_id": chatIDValue}
	if caption != "" {
		fields["caption"] = caption
	}
	if parseMode != "" {
		fields["parse_mode"] = parseMode
	}
	if len(captionEntities) > 0 {
		body, err := json.Marshal(captionEntities)
		if err != nil {
			return nil, err
		}
		fields["caption_entities"] = string(body)
	}
	if disableNotification {
		fields["disable_notification"] = "true"
	}
	if protectContent {
		fields["protect_content"] = "true"
	}

	return fields, nil
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
