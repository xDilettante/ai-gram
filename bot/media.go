package bot

import (
	"context"
	"encoding/json"
	stderrors "errors"
	"net/url"
	"strings"

	"ai-gram/telegram"
)

type fileRefKind uint8

const (
	fileRefUnknown fileRefKind = iota
	fileRefID
	fileRefURL
)

// FileRef identifies media by Telegram file_id or by an HTTP(S) URL.
type FileRef struct {
	value string
	kind  fileRefKind
}

// FileID creates a file reference from an existing Telegram file_id.
func FileID(id string) FileRef {
	return FileRef{value: id, kind: fileRefID}
}

// FileURL creates a file reference from an HTTP(S) URL.
func FileURL(rawURL string) FileRef {
	return FileRef{value: rawURL, kind: fileRefURL}
}

// MarshalJSON encodes ref as the Telegram Bot API string value.
func (ref FileRef) MarshalJSON() ([]byte, error) {
	return json.Marshal(ref.value)
}

// SendPhotoParams contains supported parameters for sendPhoto without upload.
type SendPhotoParams struct {
	ChatID              ChatID                   `json:"chat_id"`
	Photo               FileRef                  `json:"photo"`
	Caption             string                   `json:"caption,omitempty"`
	ParseMode           string                   `json:"parse_mode,omitempty"`
	CaptionEntities     []telegram.MessageEntity `json:"caption_entities,omitempty"`
	DisableNotification bool                     `json:"disable_notification,omitempty"`
	ProtectContent      bool                     `json:"protect_content,omitempty"`
}

// SendDocumentParams contains supported parameters for sendDocument without upload.
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

// SendPhoto sends a photo by Telegram file_id or HTTP(S) URL.
func (b *Bot) SendPhoto(ctx context.Context, params SendPhotoParams) (*telegram.Message, error) {
	if err := params.validate(); err != nil {
		return nil, err
	}

	var message telegram.Message
	if err := b.call(ctx, "sendPhoto", params, &message); err != nil {
		return nil, err
	}

	return &message, nil
}

// SendDocument sends a document by Telegram file_id or HTTP(S) URL.
func (b *Bot) SendDocument(ctx context.Context, params SendDocumentParams) (*telegram.Message, error) {
	if err := params.validate(); err != nil {
		return nil, err
	}

	var message telegram.Message
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

func (ref FileRef) validate(field string) error {
	value := strings.TrimSpace(ref.value)
	if value == "" {
		return stderrors.New(field + " is required")
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
