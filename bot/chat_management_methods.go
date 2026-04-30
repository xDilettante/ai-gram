package bot

import (
	"context"
	stderrors "errors"
	"strings"
)

// SetChatTitleParams contains supported parameters for setChatTitle.
type SetChatTitleParams struct {
	ChatID ChatID `json:"chat_id"`
	Title  string `json:"title"`
}

// SetChatDescriptionParams contains supported parameters for setChatDescription.
type SetChatDescriptionParams struct {
	ChatID      ChatID `json:"chat_id"`
	Description string `json:"description,omitempty"`
}

// SetChatPhotoParams contains supported parameters for setChatPhoto.
type SetChatPhotoParams struct {
	ChatID ChatID  `json:"chat_id"`
	Photo  FileRef `json:"photo"`
}

// DeleteChatPhotoParams contains supported parameters for deleteChatPhoto.
type DeleteChatPhotoParams struct {
	ChatID ChatID `json:"chat_id"`
}

// LeaveChatParams contains supported parameters for leaveChat.
type LeaveChatParams struct {
	ChatID ChatID `json:"chat_id"`
}

// SetChatStickerSetParams contains supported parameters for setChatStickerSet.
type SetChatStickerSetParams struct {
	ChatID         ChatID `json:"chat_id"`
	StickerSetName string `json:"sticker_set_name"`
}

// DeleteChatStickerSetParams contains supported parameters for deleteChatStickerSet.
type DeleteChatStickerSetParams struct {
	ChatID ChatID `json:"chat_id"`
}

// SetChatTitle changes the title of a chat.
func (b *Bot) SetChatTitle(ctx context.Context, params SetChatTitleParams) (bool, error) {
	if err := params.validate(); err != nil {
		return false, err
	}

	var result bool
	if err := b.call(ctx, "setChatTitle", params, &result); err != nil {
		return false, err
	}

	return result, nil
}

// SetChatDescription changes or removes the description of a group, supergroup, or channel.
func (b *Bot) SetChatDescription(ctx context.Context, params SetChatDescriptionParams) (bool, error) {
	if err := params.validate(); err != nil {
		return false, err
	}

	var result bool
	if err := b.call(ctx, "setChatDescription", params, &result); err != nil {
		return false, err
	}

	return result, nil
}

// SetChatPhoto uploads and sets a new profile photo for a chat.
func (b *Bot) SetChatPhoto(ctx context.Context, params SetChatPhotoParams) (bool, error) {
	if err := params.validate(); err != nil {
		return false, err
	}

	fields, files, err := params.multipart()
	if err != nil {
		return false, err
	}

	var result bool
	if err := b.callMultipart(ctx, "setChatPhoto", fields, files, &result); err != nil {
		return false, err
	}

	return result, nil
}

// DeleteChatPhoto deletes a chat profile photo.
func (b *Bot) DeleteChatPhoto(ctx context.Context, params DeleteChatPhotoParams) (bool, error) {
	if err := params.validate(); err != nil {
		return false, err
	}

	var result bool
	if err := b.call(ctx, "deleteChatPhoto", params, &result); err != nil {
		return false, err
	}

	return result, nil
}

// LeaveChat makes the bot leave a group, supergroup, or channel.
func (b *Bot) LeaveChat(ctx context.Context, params LeaveChatParams) (bool, error) {
	if err := params.validate(); err != nil {
		return false, err
	}

	var result bool
	if err := b.call(ctx, "leaveChat", params, &result); err != nil {
		return false, err
	}

	return result, nil
}

// SetChatStickerSet sets a new group sticker set for a supergroup.
func (b *Bot) SetChatStickerSet(ctx context.Context, params SetChatStickerSetParams) (bool, error) {
	if err := params.validate(); err != nil {
		return false, err
	}

	var result bool
	if err := b.call(ctx, "setChatStickerSet", params, &result); err != nil {
		return false, err
	}

	return result, nil
}

// DeleteChatStickerSet deletes the group sticker set from a supergroup.
func (b *Bot) DeleteChatStickerSet(ctx context.Context, params DeleteChatStickerSetParams) (bool, error) {
	if err := params.validate(); err != nil {
		return false, err
	}

	var result bool
	if err := b.call(ctx, "deleteChatStickerSet", params, &result); err != nil {
		return false, err
	}

	return result, nil
}

func (params SetChatTitleParams) validate() error {
	if !params.ChatID.valid() {
		return stderrors.New("chat_id is required")
	}
	if strings.TrimSpace(params.Title) == "" {
		return stderrors.New("title is required")
	}
	return nil
}

func (params SetChatDescriptionParams) validate() error {
	if !params.ChatID.valid() {
		return stderrors.New("chat_id is required")
	}
	return nil
}

func (params SetChatPhotoParams) validate() error {
	if !params.ChatID.valid() {
		return stderrors.New("chat_id is required")
	}
	if params.Photo.kind != fileRefUpload {
		return stderrors.New("photo must be uploaded with FileUpload")
	}
	return params.Photo.validate("photo")
}

func (params SetChatPhotoParams) multipart() (map[string]string, map[string]UploadFile, error) {
	chatIDValue, err := params.ChatID.multipartValue()
	if err != nil {
		return nil, nil, err
	}

	fields := map[string]string{
		"chat_id": chatIDValue,
		"photo":   "attach://photo",
	}
	files := map[string]UploadFile{
		"photo": params.Photo.upload,
	}

	return fields, files, nil
}

func (params DeleteChatPhotoParams) validate() error {
	return validateRequiredChatID(params.ChatID)
}

func (params LeaveChatParams) validate() error {
	return validateRequiredChatID(params.ChatID)
}

func (params SetChatStickerSetParams) validate() error {
	if !params.ChatID.valid() {
		return stderrors.New("chat_id is required")
	}
	if strings.TrimSpace(params.StickerSetName) == "" {
		return stderrors.New("sticker_set_name is required")
	}
	return nil
}

func (params DeleteChatStickerSetParams) validate() error {
	return validateRequiredChatID(params.ChatID)
}

func validateRequiredChatID(chatID ChatID) error {
	if !chatID.valid() {
		return stderrors.New("chat_id is required")
	}
	return nil
}
