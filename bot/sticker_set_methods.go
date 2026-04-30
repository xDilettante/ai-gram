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

const (
	stickerFormatStatic   = "static"
	stickerFormatAnimated = "animated"
	stickerFormatVideo    = "video"
)

// InputSticker describes a sticker to add to a sticker set.
type InputSticker struct {
	Sticker      FileRef                `json:"sticker"`
	Format       string                 `json:"format"`
	EmojiList    []string               `json:"emoji_list"`
	MaskPosition *telegram.MaskPosition `json:"mask_position,omitempty"`
	Keywords     []string               `json:"keywords,omitempty"`
}

// NewInputSticker creates an InputSticker with required fields.
func NewInputSticker(sticker FileRef, format string, emojiList ...string) InputSticker {
	return InputSticker{Sticker: sticker, Format: format, EmojiList: emojiList}
}

// GetStickerSetParams contains supported parameters for getStickerSet.
type GetStickerSetParams struct {
	Name string `json:"name"`
}

// GetCustomEmojiStickersParams contains supported parameters for getCustomEmojiStickers.
type GetCustomEmojiStickersParams struct {
	CustomEmojiIDs []string `json:"custom_emoji_ids"`
}

// UploadStickerFileParams contains supported parameters for uploadStickerFile.
type UploadStickerFileParams struct {
	UserID        int64   `json:"user_id"`
	Sticker       FileRef `json:"sticker"`
	StickerFormat string  `json:"sticker_format"`
}

// CreateNewStickerSetParams contains supported parameters for createNewStickerSet.
type CreateNewStickerSetParams struct {
	UserID          int64          `json:"user_id"`
	Name            string         `json:"name"`
	Title           string         `json:"title"`
	Stickers        []InputSticker `json:"stickers"`
	StickerType     string         `json:"sticker_type,omitempty"`
	NeedsRepainting bool           `json:"needs_repainting,omitempty"`
}

// AddStickerToSetParams contains supported parameters for addStickerToSet.
type AddStickerToSetParams struct {
	UserID  int64        `json:"user_id"`
	Name    string       `json:"name"`
	Sticker InputSticker `json:"sticker"`
}

// ReplaceStickerInSetParams contains supported parameters for replaceStickerInSet.
type ReplaceStickerInSetParams struct {
	UserID     int64        `json:"user_id"`
	Name       string       `json:"name"`
	OldSticker string       `json:"old_sticker"`
	Sticker    InputSticker `json:"sticker"`
}

// SetStickerPositionInSetParams contains supported parameters for setStickerPositionInSet.
type SetStickerPositionInSetParams struct {
	Sticker  string `json:"sticker"`
	Position int    `json:"position"`
}

// DeleteStickerFromSetParams contains supported parameters for deleteStickerFromSet.
type DeleteStickerFromSetParams struct {
	Sticker string `json:"sticker"`
}

// SetStickerEmojiListParams contains supported parameters for setStickerEmojiList.
type SetStickerEmojiListParams struct {
	Sticker   string   `json:"sticker"`
	EmojiList []string `json:"emoji_list"`
}

// SetStickerKeywordsParams contains supported parameters for setStickerKeywords.
type SetStickerKeywordsParams struct {
	Sticker  string   `json:"sticker"`
	Keywords []string `json:"keywords,omitempty"`
}

// SetStickerMaskPositionParams contains supported parameters for setStickerMaskPosition.
type SetStickerMaskPositionParams struct {
	Sticker      string                 `json:"sticker"`
	MaskPosition *telegram.MaskPosition `json:"mask_position,omitempty"`
}

// SetStickerSetTitleParams contains supported parameters for setStickerSetTitle.
type SetStickerSetTitleParams struct {
	Name  string `json:"name"`
	Title string `json:"title"`
}

// SetStickerSetThumbnailParams contains supported parameters for setStickerSetThumbnail.
type SetStickerSetThumbnailParams struct {
	Name      string  `json:"name"`
	UserID    int64   `json:"user_id"`
	Thumbnail FileRef `json:"thumbnail,omitempty"`
	Format    string  `json:"format"`
}

// SetCustomEmojiStickerSetThumbnailParams contains supported parameters for setCustomEmojiStickerSetThumbnail.
type SetCustomEmojiStickerSetThumbnailParams struct {
	Name          string `json:"name"`
	CustomEmojiID string `json:"custom_emoji_id"`
}

// DeleteStickerSetParams contains supported parameters for deleteStickerSet.
type DeleteStickerSetParams struct {
	Name string `json:"name"`
}

// GetStickerSet gets a sticker set by name.
func (b *Bot) GetStickerSet(ctx context.Context, params GetStickerSetParams) (*telegram.StickerSet, error) {
	if err := params.validate(); err != nil {
		return nil, err
	}
	var result telegram.StickerSet
	if err := b.call(ctx, "getStickerSet", params, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetCustomEmojiStickers gets custom emoji stickers by identifiers.
func (b *Bot) GetCustomEmojiStickers(ctx context.Context, params GetCustomEmojiStickersParams) ([]telegram.Sticker, error) {
	if err := params.validate(); err != nil {
		return nil, err
	}
	var result []telegram.Sticker
	if err := b.call(ctx, "getCustomEmojiStickers", params, &result); err != nil {
		return nil, err
	}
	return result, nil
}

// UploadStickerFile uploads a sticker file for later sticker set operations.
func (b *Bot) UploadStickerFile(ctx context.Context, params UploadStickerFileParams) (*telegram.File, error) {
	if err := params.validate(); err != nil {
		return nil, err
	}
	fields := map[string]string{
		"user_id":        strconv.FormatInt(params.UserID, 10),
		"sticker_format": params.StickerFormat,
	}
	files := map[string]UploadFile{"sticker": params.Sticker.upload}
	var result telegram.File
	if err := b.callMultipart(ctx, "uploadStickerFile", fields, files, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// CreateNewStickerSet creates a new sticker set owned by a user.
func (b *Bot) CreateNewStickerSet(ctx context.Context, params CreateNewStickerSetParams) (bool, error) {
	if err := params.validate(); err != nil {
		return false, err
	}
	payload, fields, files, err := params.payload()
	if err != nil {
		return false, err
	}
	var result bool
	if len(files) > 0 {
		if err := b.callMultipart(ctx, "createNewStickerSet", fields, files, &result); err != nil {
			return false, err
		}
		return result, nil
	}
	if err := b.call(ctx, "createNewStickerSet", payload, &result); err != nil {
		return false, err
	}
	return result, nil
}

// AddStickerToSet adds a sticker to a set created by the bot.
func (b *Bot) AddStickerToSet(ctx context.Context, params AddStickerToSetParams) (bool, error) {
	if err := params.validate(); err != nil {
		return false, err
	}
	payload, fields, files, err := params.payload()
	if err != nil {
		return false, err
	}
	var result bool
	if len(files) > 0 {
		if err := b.callMultipart(ctx, "addStickerToSet", fields, files, &result); err != nil {
			return false, err
		}
		return result, nil
	}
	if err := b.call(ctx, "addStickerToSet", payload, &result); err != nil {
		return false, err
	}
	return result, nil
}

// ReplaceStickerInSet replaces a sticker in a set created by the bot.
func (b *Bot) ReplaceStickerInSet(ctx context.Context, params ReplaceStickerInSetParams) (bool, error) {
	if err := params.validate(); err != nil {
		return false, err
	}
	payload, fields, files, err := params.payload()
	if err != nil {
		return false, err
	}
	var result bool
	if len(files) > 0 {
		if err := b.callMultipart(ctx, "replaceStickerInSet", fields, files, &result); err != nil {
			return false, err
		}
		return result, nil
	}
	if err := b.call(ctx, "replaceStickerInSet", payload, &result); err != nil {
		return false, err
	}
	return result, nil
}

// SetStickerPositionInSet moves a sticker to a zero-based position in its set.
func (b *Bot) SetStickerPositionInSet(ctx context.Context, params SetStickerPositionInSetParams) (bool, error) {
	if err := params.validate(); err != nil {
		return false, err
	}
	return b.callStickerBool(ctx, "setStickerPositionInSet", params)
}

// DeleteStickerFromSet deletes a sticker from a set created by the bot.
func (b *Bot) DeleteStickerFromSet(ctx context.Context, params DeleteStickerFromSetParams) (bool, error) {
	if err := params.validate(); err != nil {
		return false, err
	}
	return b.callStickerBool(ctx, "deleteStickerFromSet", params)
}

// SetStickerEmojiList changes the emoji list assigned to a sticker.
func (b *Bot) SetStickerEmojiList(ctx context.Context, params SetStickerEmojiListParams) (bool, error) {
	if err := params.validate(); err != nil {
		return false, err
	}
	return b.callStickerBool(ctx, "setStickerEmojiList", params)
}

// SetStickerKeywords changes search keywords assigned to a sticker.
func (b *Bot) SetStickerKeywords(ctx context.Context, params SetStickerKeywordsParams) (bool, error) {
	if err := params.validate(); err != nil {
		return false, err
	}
	return b.callStickerBool(ctx, "setStickerKeywords", params.payload())
}

// SetStickerMaskPosition changes or removes the mask position assigned to a mask sticker.
func (b *Bot) SetStickerMaskPosition(ctx context.Context, params SetStickerMaskPositionParams) (bool, error) {
	if err := params.validate(); err != nil {
		return false, err
	}
	return b.callStickerBool(ctx, "setStickerMaskPosition", params)
}

// SetStickerSetTitle changes a sticker set title.
func (b *Bot) SetStickerSetTitle(ctx context.Context, params SetStickerSetTitleParams) (bool, error) {
	if err := params.validate(); err != nil {
		return false, err
	}
	return b.callStickerBool(ctx, "setStickerSetTitle", params)
}

// SetStickerSetThumbnail changes or removes a regular or mask sticker set thumbnail.
func (b *Bot) SetStickerSetThumbnail(ctx context.Context, params SetStickerSetThumbnailParams) (bool, error) {
	if err := params.validate(); err != nil {
		return false, err
	}
	var result bool
	if params.Thumbnail.isUpload() {
		fields, files, err := params.multipart()
		if err != nil {
			return false, err
		}
		if err := b.callMultipart(ctx, "setStickerSetThumbnail", fields, files, &result); err != nil {
			return false, err
		}
		return result, nil
	}
	if err := b.call(ctx, "setStickerSetThumbnail", params.payload(), &result); err != nil {
		return false, err
	}
	return result, nil
}

// SetCustomEmojiStickerSetThumbnail changes or removes a custom emoji sticker set thumbnail.
func (b *Bot) SetCustomEmojiStickerSetThumbnail(ctx context.Context, params SetCustomEmojiStickerSetThumbnailParams) (bool, error) {
	if err := params.validate(); err != nil {
		return false, err
	}
	return b.callStickerBool(ctx, "setCustomEmojiStickerSetThumbnail", params)
}

// DeleteStickerSet deletes a sticker set created by the bot.
func (b *Bot) DeleteStickerSet(ctx context.Context, params DeleteStickerSetParams) (bool, error) {
	if err := params.validate(); err != nil {
		return false, err
	}
	return b.callStickerBool(ctx, "deleteStickerSet", params)
}

func (b *Bot) callStickerBool(ctx context.Context, method string, payload any) (bool, error) {
	var result bool
	if err := b.call(ctx, method, payload, &result); err != nil {
		return false, err
	}
	return result, nil
}

func (params GetStickerSetParams) validate() error {
	if strings.TrimSpace(params.Name) == "" {
		return stderrors.New("name is required")
	}
	return nil
}

func (params GetCustomEmojiStickersParams) validate() error {
	if len(params.CustomEmojiIDs) == 0 {
		return stderrors.New("custom_emoji_ids must not be empty")
	}
	for _, id := range params.CustomEmojiIDs {
		if strings.TrimSpace(id) == "" {
			return stderrors.New("custom_emoji_id is required")
		}
	}
	return nil
}

func (params UploadStickerFileParams) validate() error {
	if params.UserID <= 0 {
		return stderrors.New("user_id must be greater than zero")
	}
	if params.Sticker.kind != fileRefUpload {
		return stderrors.New("sticker must be uploaded with FileUpload")
	}
	if err := params.Sticker.validate("sticker"); err != nil {
		return err
	}
	return validateStickerFormat(params.StickerFormat, "sticker_format")
}

func (params CreateNewStickerSetParams) validate() error {
	if params.UserID <= 0 {
		return stderrors.New("user_id must be greater than zero")
	}
	if strings.TrimSpace(params.Name) == "" {
		return stderrors.New("name is required")
	}
	if strings.TrimSpace(params.Title) == "" {
		return stderrors.New("title is required")
	}
	if len(params.Stickers) == 0 {
		return stderrors.New("stickers must not be empty")
	}
	for i, sticker := range params.Stickers {
		if err := sticker.validate(); err != nil {
			return fmt.Errorf("stickers[%d]: %w", i, err)
		}
	}
	return validateOptionalStickerType(params.StickerType)
}

func (params AddStickerToSetParams) validate() error {
	if params.UserID <= 0 {
		return stderrors.New("user_id must be greater than zero")
	}
	if strings.TrimSpace(params.Name) == "" {
		return stderrors.New("name is required")
	}
	return params.Sticker.validate()
}

func (params ReplaceStickerInSetParams) validate() error {
	if params.UserID <= 0 {
		return stderrors.New("user_id must be greater than zero")
	}
	if strings.TrimSpace(params.Name) == "" {
		return stderrors.New("name is required")
	}
	if strings.TrimSpace(params.OldSticker) == "" {
		return stderrors.New("old_sticker is required")
	}
	return params.Sticker.validate()
}

func (params SetStickerPositionInSetParams) validate() error {
	if strings.TrimSpace(params.Sticker) == "" {
		return stderrors.New("sticker is required")
	}
	if params.Position < 0 {
		return stderrors.New("position must not be negative")
	}
	return nil
}

func (params DeleteStickerFromSetParams) validate() error {
	return validateStickerID(params.Sticker)
}

func (params SetStickerEmojiListParams) validate() error {
	if err := validateStickerID(params.Sticker); err != nil {
		return err
	}
	return validateNonEmptyStrings(params.EmojiList, "emoji_list", "emoji")
}

func (params SetStickerKeywordsParams) validate() error {
	if err := validateStickerID(params.Sticker); err != nil {
		return err
	}
	for _, keyword := range params.Keywords {
		if strings.TrimSpace(keyword) == "" {
			return stderrors.New("keyword is required")
		}
	}
	return nil
}

func (params SetStickerKeywordsParams) payload() map[string]any {
	payload := map[string]any{"sticker": params.Sticker}
	if params.Keywords != nil {
		payload["keywords"] = params.Keywords
	}
	return payload
}

func (params SetStickerMaskPositionParams) validate() error {
	if err := validateStickerID(params.Sticker); err != nil {
		return err
	}
	return validateMaskPosition(params.MaskPosition)
}

func (params SetStickerSetTitleParams) validate() error {
	if strings.TrimSpace(params.Name) == "" {
		return stderrors.New("name is required")
	}
	if strings.TrimSpace(params.Title) == "" {
		return stderrors.New("title is required")
	}
	return nil
}

func (params SetStickerSetThumbnailParams) validate() error {
	if strings.TrimSpace(params.Name) == "" {
		return stderrors.New("name is required")
	}
	if params.UserID <= 0 {
		return stderrors.New("user_id must be greater than zero")
	}
	if err := validateStickerFormat(params.Format, "format"); err != nil {
		return err
	}
	if params.Thumbnail.kind == fileRefUnknown && strings.TrimSpace(params.Thumbnail.value) == "" {
		return nil
	}
	if err := params.Thumbnail.validate("thumbnail"); err != nil {
		return err
	}
	if params.Format != stickerFormatStatic && (params.Thumbnail.kind == fileRefURL || strings.Contains(strings.TrimSpace(params.Thumbnail.value), "://")) {
		return stderrors.New("thumbnail URL is supported only for static sticker set thumbnails")
	}
	return nil
}

func (params SetStickerSetThumbnailParams) payload() map[string]any {
	payload := map[string]any{"name": params.Name, "user_id": params.UserID, "format": params.Format}
	if params.Thumbnail.kind != fileRefUnknown || strings.TrimSpace(params.Thumbnail.value) != "" {
		payload["thumbnail"] = params.Thumbnail
	}
	return payload
}

func (params SetStickerSetThumbnailParams) multipart() (map[string]string, map[string]UploadFile, error) {
	fields := map[string]string{"name": params.Name, "user_id": strconv.FormatInt(params.UserID, 10), "format": params.Format, "thumbnail": "attach://thumbnail"}
	return fields, map[string]UploadFile{"thumbnail": params.Thumbnail.upload}, nil
}

func (params SetCustomEmojiStickerSetThumbnailParams) validate() error {
	if strings.TrimSpace(params.Name) == "" {
		return stderrors.New("name is required")
	}
	return nil
}

func (params DeleteStickerSetParams) validate() error {
	if strings.TrimSpace(params.Name) == "" {
		return stderrors.New("name is required")
	}
	return nil
}

func (sticker InputSticker) validate() error {
	if err := sticker.Sticker.validate("sticker"); err != nil {
		return err
	}
	if err := validateStickerFormat(sticker.Format, "format"); err != nil {
		return err
	}
	if sticker.Format != stickerFormatStatic && (sticker.Sticker.kind == fileRefURL || strings.Contains(strings.TrimSpace(sticker.Sticker.value), "://")) {
		return stderrors.New("sticker URL is supported only for static stickers")
	}
	if err := validateNonEmptyStrings(sticker.EmojiList, "emoji_list", "emoji"); err != nil {
		return err
	}
	for _, keyword := range sticker.Keywords {
		if strings.TrimSpace(keyword) == "" {
			return stderrors.New("keyword is required")
		}
	}
	return validateMaskPosition(sticker.MaskPosition)
}

func (sticker InputSticker) payload(uploadName string) (map[string]any, map[string]UploadFile, error) {
	if err := sticker.validate(); err != nil {
		return nil, nil, err
	}
	payload := map[string]any{
		"format":     sticker.Format,
		"emoji_list": sticker.EmojiList,
	}
	files := map[string]UploadFile{}
	if sticker.Sticker.isUpload() {
		payload["sticker"] = "attach://" + uploadName
		files[uploadName] = sticker.Sticker.upload
	} else {
		payload["sticker"] = sticker.Sticker
	}
	if sticker.MaskPosition != nil {
		payload["mask_position"] = sticker.MaskPosition
	}
	if sticker.Keywords != nil {
		payload["keywords"] = sticker.Keywords
	}
	return payload, files, nil
}

func (params CreateNewStickerSetParams) payload() (map[string]any, map[string]string, map[string]UploadFile, error) {
	stickers := make([]map[string]any, 0, len(params.Stickers))
	files := map[string]UploadFile{}
	for i, sticker := range params.Stickers {
		payload, stickerFiles, err := sticker.payload("sticker" + strconv.Itoa(i))
		if err != nil {
			return nil, nil, nil, err
		}
		stickers = append(stickers, payload)
		for name, file := range stickerFiles {
			files[name] = file
		}
	}

	payload := map[string]any{
		"user_id":  params.UserID,
		"name":     params.Name,
		"title":    params.Title,
		"stickers": stickers,
	}
	if params.StickerType != "" {
		payload["sticker_type"] = params.StickerType
	}
	if params.NeedsRepainting {
		payload["needs_repainting"] = true
	}
	fields, err := stickerMultipartFields(payload, "stickers")
	return payload, fields, files, err
}

func (params AddStickerToSetParams) payload() (map[string]any, map[string]string, map[string]UploadFile, error) {
	sticker, files, err := params.Sticker.payload("sticker")
	if err != nil {
		return nil, nil, nil, err
	}
	payload := map[string]any{"user_id": params.UserID, "name": params.Name, "sticker": sticker}
	fields, err := stickerMultipartFields(payload, "sticker")
	return payload, fields, files, err
}

func (params ReplaceStickerInSetParams) payload() (map[string]any, map[string]string, map[string]UploadFile, error) {
	sticker, files, err := params.Sticker.payload("sticker")
	if err != nil {
		return nil, nil, nil, err
	}
	payload := map[string]any{"user_id": params.UserID, "name": params.Name, "old_sticker": params.OldSticker, "sticker": sticker}
	fields, err := stickerMultipartFields(payload, "sticker")
	return payload, fields, files, err
}

func stickerMultipartFields(payload map[string]any, jsonFields ...string) (map[string]string, error) {
	fields := map[string]string{}
	jsonFieldSet := map[string]bool{}
	for _, name := range jsonFields {
		jsonFieldSet[name] = true
	}
	for name, value := range payload {
		if jsonFieldSet[name] {
			body, err := json.Marshal(value)
			if err != nil {
				return nil, err
			}
			fields[name] = string(body)
			continue
		}
		switch typed := value.(type) {
		case string:
			fields[name] = typed
		case int64:
			fields[name] = strconv.FormatInt(typed, 10)
		case bool:
			if typed {
				fields[name] = "true"
			}
		}
	}
	return fields, nil
}

func validateStickerID(sticker string) error {
	if strings.TrimSpace(sticker) == "" {
		return stderrors.New("sticker is required")
	}
	return nil
}

func validateStickerFormat(format string, field string) error {
	switch format {
	case stickerFormatStatic, stickerFormatAnimated, stickerFormatVideo:
		return nil
	case "":
		return stderrors.New(field + " is required")
	default:
		return stderrors.New(field + " is unsupported")
	}
}

func validateOptionalStickerType(stickerType string) error {
	switch stickerType {
	case "", "regular", "mask", "custom_emoji":
		return nil
	default:
		return stderrors.New("sticker_type is unsupported")
	}
}

func validateNonEmptyStrings(values []string, field string, itemName string) error {
	if len(values) == 0 {
		return stderrors.New(field + " must not be empty")
	}
	for _, value := range values {
		if strings.TrimSpace(value) == "" {
			return stderrors.New(itemName + " is required")
		}
	}
	return nil
}

func validateMaskPosition(position *telegram.MaskPosition) error {
	if position == nil {
		return nil
	}
	if strings.TrimSpace(position.Point) == "" {
		return stderrors.New("mask_position point is required")
	}
	if position.Scale <= 0 {
		return stderrors.New("mask_position scale must be greater than zero")
	}
	return nil
}
