package bot

import (
	"context"
	stderrors "errors"
)

const (
	// ChatActionTyping tells Telegram clients that the bot is typing.
	ChatActionTyping = "typing"
	// ChatActionUploadPhoto tells Telegram clients that the bot is uploading a photo.
	ChatActionUploadPhoto = "upload_photo"
	// ChatActionRecordVideo tells Telegram clients that the bot is recording a video.
	ChatActionRecordVideo = "record_video"
	// ChatActionUploadVideo tells Telegram clients that the bot is uploading a video.
	ChatActionUploadVideo = "upload_video"
	// ChatActionRecordVoice tells Telegram clients that the bot is recording a voice message.
	ChatActionRecordVoice = "record_voice"
	// ChatActionUploadVoice tells Telegram clients that the bot is uploading a voice message.
	ChatActionUploadVoice = "upload_voice"
	// ChatActionUploadDocument tells Telegram clients that the bot is uploading a document.
	ChatActionUploadDocument = "upload_document"
	// ChatActionChooseSticker tells Telegram clients that the bot is choosing a sticker.
	ChatActionChooseSticker = "choose_sticker"
	// ChatActionFindLocation tells Telegram clients that the bot is finding a location.
	ChatActionFindLocation = "find_location"
	// ChatActionRecordVideoNote tells Telegram clients that the bot is recording a video note.
	ChatActionRecordVideoNote = "record_video_note"
	// ChatActionUploadVideoNote tells Telegram clients that the bot is uploading a video note.
	ChatActionUploadVideoNote = "upload_video_note"
)

// SendChatActionParams contains supported parameters for sendChatAction.
type SendChatActionParams struct {
	BusinessConnectionID string `json:"business_connection_id,omitempty"`
	ChatID               ChatID `json:"chat_id"`
	MessageThreadID      int64  `json:"message_thread_id,omitempty"`
	Action               string `json:"action"`
}

// PinChatMessageParams contains supported parameters for pinChatMessage.
type PinChatMessageParams struct {
	BusinessConnectionID string `json:"business_connection_id,omitempty"`
	ChatID               ChatID `json:"chat_id"`
	MessageID            int64  `json:"message_id"`
	DisableNotification  bool   `json:"disable_notification,omitempty"`
}

// UnpinChatMessageParams contains supported parameters for unpinChatMessage.
type UnpinChatMessageParams struct {
	BusinessConnectionID string `json:"business_connection_id,omitempty"`
	ChatID               ChatID `json:"chat_id"`
	MessageID            int64  `json:"message_id,omitempty"`
}

// UnpinAllChatMessagesParams contains supported parameters for unpinAllChatMessages.
type UnpinAllChatMessagesParams struct {
	ChatID ChatID `json:"chat_id"`
}

// SendChatAction sends a chat action, such as typing or uploading a document.
func (b *Bot) SendChatAction(ctx context.Context, params SendChatActionParams) (bool, error) {
	if err := params.validate(); err != nil {
		return false, err
	}

	var result bool
	if err := b.call(ctx, "sendChatAction", params, &result); err != nil {
		return false, err
	}

	return result, nil
}

// PinChatMessage pins a message in a chat.
func (b *Bot) PinChatMessage(ctx context.Context, params PinChatMessageParams) (bool, error) {
	if err := params.validate(); err != nil {
		return false, err
	}

	var result bool
	if err := b.call(ctx, "pinChatMessage", params, &result); err != nil {
		return false, err
	}

	return result, nil
}

// UnpinChatMessage unpins a message in a chat, or the most recent pinned message when MessageID is omitted.
func (b *Bot) UnpinChatMessage(ctx context.Context, params UnpinChatMessageParams) (bool, error) {
	if err := params.validate(); err != nil {
		return false, err
	}

	var result bool
	if err := b.call(ctx, "unpinChatMessage", params, &result); err != nil {
		return false, err
	}

	return result, nil
}

// UnpinAllChatMessages unpins all pinned messages in a chat.
func (b *Bot) UnpinAllChatMessages(ctx context.Context, params UnpinAllChatMessagesParams) (bool, error) {
	if err := params.validate(); err != nil {
		return false, err
	}

	var result bool
	if err := b.call(ctx, "unpinAllChatMessages", params, &result); err != nil {
		return false, err
	}

	return result, nil
}

func (params SendChatActionParams) validate() error {
	if !params.ChatID.valid() {
		return stderrors.New("chat_id is required")
	}
	if err := validateMessageThreadID(params.MessageThreadID); err != nil {
		return err
	}
	if params.Action == "" {
		return stderrors.New("action is required")
	}
	if !validChatAction(params.Action) {
		return stderrors.New("action must be a known chat action")
	}

	return nil
}

func (params PinChatMessageParams) validate() error {
	if !params.ChatID.valid() {
		return stderrors.New("chat_id is required")
	}
	if params.MessageID <= 0 {
		return stderrors.New("message_id must be greater than zero")
	}

	return nil
}

func (params UnpinChatMessageParams) validate() error {
	if !params.ChatID.valid() {
		return stderrors.New("chat_id is required")
	}
	if params.MessageID < 0 {
		return stderrors.New("message_id must not be negative")
	}

	return nil
}

func (params UnpinAllChatMessagesParams) validate() error {
	if !params.ChatID.valid() {
		return stderrors.New("chat_id is required")
	}

	return nil
}

func validChatAction(action string) bool {
	switch action {
	case ChatActionTyping,
		ChatActionUploadPhoto,
		ChatActionRecordVideo,
		ChatActionUploadVideo,
		ChatActionRecordVoice,
		ChatActionUploadVoice,
		ChatActionUploadDocument,
		ChatActionChooseSticker,
		ChatActionFindLocation,
		ChatActionRecordVideoNote,
		ChatActionUploadVideoNote:
		return true
	default:
		return false
	}
}
