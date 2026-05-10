package bot

import (
	"encoding/json"
	stderrors "errors"
	"fmt"
	"strings"

	"github.com/xDilettante/ai-gram/telegram"
)

const (
	inlineQueryResultPhotoType    = "photo"
	inlineQueryResultGifType      = "gif"
	inlineQueryResultMpeg4GifType = "mpeg4_gif"
	inlineQueryResultVideoType    = "video"
	inlineQueryResultAudioType    = "audio"
	inlineQueryResultVoiceType    = "voice"
	inlineQueryResultDocumentType = "document"
	inlineQueryResultStickerType  = "sticker"
)

// InlineQueryResultPhoto represents a photo inline query result.
type InlineQueryResultPhoto struct {
	Type                  string                         `json:"type"`
	ID                    string                         `json:"id"`
	PhotoURL              string                         `json:"photo_url"`
	ThumbnailURL          string                         `json:"thumbnail_url"`
	PhotoWidth            int                            `json:"photo_width,omitempty"`
	PhotoHeight           int                            `json:"photo_height,omitempty"`
	Title                 string                         `json:"title,omitempty"`
	Description           string                         `json:"description,omitempty"`
	Caption               string                         `json:"caption,omitempty"`
	ParseMode             string                         `json:"parse_mode,omitempty"`
	CaptionEntities       []telegram.MessageEntity       `json:"caption_entities,omitempty"`
	ShowCaptionAboveMedia bool                           `json:"show_caption_above_media,omitempty"`
	ReplyMarkup           *telegram.InlineKeyboardMarkup `json:"reply_markup,omitempty"`
	InputMessageContent   InputMessageContent            `json:"input_message_content,omitempty"`
}

// InlineQueryResultGif represents a GIF inline query result.
type InlineQueryResultGif struct {
	Type                  string                         `json:"type"`
	ID                    string                         `json:"id"`
	GifURL                string                         `json:"gif_url"`
	GifWidth              int                            `json:"gif_width,omitempty"`
	GifHeight             int                            `json:"gif_height,omitempty"`
	GifDuration           int                            `json:"gif_duration,omitempty"`
	ThumbnailURL          string                         `json:"thumbnail_url"`
	ThumbnailMimeType     string                         `json:"thumbnail_mime_type,omitempty"`
	Title                 string                         `json:"title,omitempty"`
	Caption               string                         `json:"caption,omitempty"`
	ParseMode             string                         `json:"parse_mode,omitempty"`
	CaptionEntities       []telegram.MessageEntity       `json:"caption_entities,omitempty"`
	ShowCaptionAboveMedia bool                           `json:"show_caption_above_media,omitempty"`
	ReplyMarkup           *telegram.InlineKeyboardMarkup `json:"reply_markup,omitempty"`
	InputMessageContent   InputMessageContent            `json:"input_message_content,omitempty"`
}

// InlineQueryResultMpeg4Gif represents an MPEG-4 GIF inline query result.
type InlineQueryResultMpeg4Gif struct {
	Type                  string                         `json:"type"`
	ID                    string                         `json:"id"`
	Mpeg4URL              string                         `json:"mpeg4_url"`
	Mpeg4Width            int                            `json:"mpeg4_width,omitempty"`
	Mpeg4Height           int                            `json:"mpeg4_height,omitempty"`
	Mpeg4Duration         int                            `json:"mpeg4_duration,omitempty"`
	ThumbnailURL          string                         `json:"thumbnail_url"`
	ThumbnailMimeType     string                         `json:"thumbnail_mime_type,omitempty"`
	Title                 string                         `json:"title,omitempty"`
	Caption               string                         `json:"caption,omitempty"`
	ParseMode             string                         `json:"parse_mode,omitempty"`
	CaptionEntities       []telegram.MessageEntity       `json:"caption_entities,omitempty"`
	ShowCaptionAboveMedia bool                           `json:"show_caption_above_media,omitempty"`
	ReplyMarkup           *telegram.InlineKeyboardMarkup `json:"reply_markup,omitempty"`
	InputMessageContent   InputMessageContent            `json:"input_message_content,omitempty"`
}

// InlineQueryResultVideo represents a video inline query result.
type InlineQueryResultVideo struct {
	Type                  string                         `json:"type"`
	ID                    string                         `json:"id"`
	VideoURL              string                         `json:"video_url"`
	MimeType              string                         `json:"mime_type"`
	ThumbnailURL          string                         `json:"thumbnail_url"`
	Title                 string                         `json:"title"`
	Caption               string                         `json:"caption,omitempty"`
	ParseMode             string                         `json:"parse_mode,omitempty"`
	CaptionEntities       []telegram.MessageEntity       `json:"caption_entities,omitempty"`
	ShowCaptionAboveMedia bool                           `json:"show_caption_above_media,omitempty"`
	VideoWidth            int                            `json:"video_width,omitempty"`
	VideoHeight           int                            `json:"video_height,omitempty"`
	VideoDuration         int                            `json:"video_duration,omitempty"`
	Description           string                         `json:"description,omitempty"`
	ReplyMarkup           *telegram.InlineKeyboardMarkup `json:"reply_markup,omitempty"`
	InputMessageContent   InputMessageContent            `json:"input_message_content,omitempty"`
}

// InlineQueryResultAudio represents an audio inline query result.
type InlineQueryResultAudio struct {
	Type                string                         `json:"type"`
	ID                  string                         `json:"id"`
	AudioURL            string                         `json:"audio_url"`
	Title               string                         `json:"title"`
	Caption             string                         `json:"caption,omitempty"`
	ParseMode           string                         `json:"parse_mode,omitempty"`
	CaptionEntities     []telegram.MessageEntity       `json:"caption_entities,omitempty"`
	Performer           string                         `json:"performer,omitempty"`
	AudioDuration       int                            `json:"audio_duration,omitempty"`
	ReplyMarkup         *telegram.InlineKeyboardMarkup `json:"reply_markup,omitempty"`
	InputMessageContent InputMessageContent            `json:"input_message_content,omitempty"`
}

// InlineQueryResultVoice represents a voice inline query result.
type InlineQueryResultVoice struct {
	Type                string                         `json:"type"`
	ID                  string                         `json:"id"`
	VoiceURL            string                         `json:"voice_url"`
	Title               string                         `json:"title"`
	Caption             string                         `json:"caption,omitempty"`
	ParseMode           string                         `json:"parse_mode,omitempty"`
	CaptionEntities     []telegram.MessageEntity       `json:"caption_entities,omitempty"`
	VoiceDuration       int                            `json:"voice_duration,omitempty"`
	ReplyMarkup         *telegram.InlineKeyboardMarkup `json:"reply_markup,omitempty"`
	InputMessageContent InputMessageContent            `json:"input_message_content,omitempty"`
}

// InlineQueryResultDocument represents a document inline query result.
type InlineQueryResultDocument struct {
	Type                string                         `json:"type"`
	ID                  string                         `json:"id"`
	Title               string                         `json:"title"`
	Caption             string                         `json:"caption,omitempty"`
	ParseMode           string                         `json:"parse_mode,omitempty"`
	CaptionEntities     []telegram.MessageEntity       `json:"caption_entities,omitempty"`
	DocumentURL         string                         `json:"document_url"`
	MimeType            string                         `json:"mime_type"`
	Description         string                         `json:"description,omitempty"`
	ReplyMarkup         *telegram.InlineKeyboardMarkup `json:"reply_markup,omitempty"`
	InputMessageContent InputMessageContent            `json:"input_message_content,omitempty"`
	ThumbnailURL        string                         `json:"thumbnail_url,omitempty"`
	ThumbnailWidth      int                            `json:"thumbnail_width,omitempty"`
	ThumbnailHeight     int                            `json:"thumbnail_height,omitempty"`
}

// InlineQueryResultCachedPhoto represents a cached photo inline query result.
type InlineQueryResultCachedPhoto struct {
	Type                  string                         `json:"type"`
	ID                    string                         `json:"id"`
	PhotoFileID           string                         `json:"photo_file_id"`
	Title                 string                         `json:"title,omitempty"`
	Description           string                         `json:"description,omitempty"`
	Caption               string                         `json:"caption,omitempty"`
	ParseMode             string                         `json:"parse_mode,omitempty"`
	CaptionEntities       []telegram.MessageEntity       `json:"caption_entities,omitempty"`
	ShowCaptionAboveMedia bool                           `json:"show_caption_above_media,omitempty"`
	ReplyMarkup           *telegram.InlineKeyboardMarkup `json:"reply_markup,omitempty"`
	InputMessageContent   InputMessageContent            `json:"input_message_content,omitempty"`
}

// InlineQueryResultCachedGif represents a cached GIF inline query result.
type InlineQueryResultCachedGif struct {
	Type                  string                         `json:"type"`
	ID                    string                         `json:"id"`
	GifFileID             string                         `json:"gif_file_id"`
	Title                 string                         `json:"title,omitempty"`
	Caption               string                         `json:"caption,omitempty"`
	ParseMode             string                         `json:"parse_mode,omitempty"`
	CaptionEntities       []telegram.MessageEntity       `json:"caption_entities,omitempty"`
	ShowCaptionAboveMedia bool                           `json:"show_caption_above_media,omitempty"`
	ReplyMarkup           *telegram.InlineKeyboardMarkup `json:"reply_markup,omitempty"`
	InputMessageContent   InputMessageContent            `json:"input_message_content,omitempty"`
}

// InlineQueryResultCachedMpeg4Gif represents a cached MPEG-4 GIF inline query result.
type InlineQueryResultCachedMpeg4Gif struct {
	Type                  string                         `json:"type"`
	ID                    string                         `json:"id"`
	Mpeg4FileID           string                         `json:"mpeg4_file_id"`
	Title                 string                         `json:"title,omitempty"`
	Caption               string                         `json:"caption,omitempty"`
	ParseMode             string                         `json:"parse_mode,omitempty"`
	CaptionEntities       []telegram.MessageEntity       `json:"caption_entities,omitempty"`
	ShowCaptionAboveMedia bool                           `json:"show_caption_above_media,omitempty"`
	ReplyMarkup           *telegram.InlineKeyboardMarkup `json:"reply_markup,omitempty"`
	InputMessageContent   InputMessageContent            `json:"input_message_content,omitempty"`
}

// InlineQueryResultCachedSticker represents a cached sticker inline query result.
type InlineQueryResultCachedSticker struct {
	Type                string                         `json:"type"`
	ID                  string                         `json:"id"`
	StickerFileID       string                         `json:"sticker_file_id"`
	ReplyMarkup         *telegram.InlineKeyboardMarkup `json:"reply_markup,omitempty"`
	InputMessageContent InputMessageContent            `json:"input_message_content,omitempty"`
}

// InlineQueryResultCachedDocument represents a cached document inline query result.
type InlineQueryResultCachedDocument struct {
	Type                string                         `json:"type"`
	ID                  string                         `json:"id"`
	Title               string                         `json:"title"`
	DocumentFileID      string                         `json:"document_file_id"`
	Description         string                         `json:"description,omitempty"`
	Caption             string                         `json:"caption,omitempty"`
	ParseMode           string                         `json:"parse_mode,omitempty"`
	CaptionEntities     []telegram.MessageEntity       `json:"caption_entities,omitempty"`
	ReplyMarkup         *telegram.InlineKeyboardMarkup `json:"reply_markup,omitempty"`
	InputMessageContent InputMessageContent            `json:"input_message_content,omitempty"`
}

// InlineQueryResultCachedVideo represents a cached video inline query result.
type InlineQueryResultCachedVideo struct {
	Type                  string                         `json:"type"`
	ID                    string                         `json:"id"`
	VideoFileID           string                         `json:"video_file_id"`
	Title                 string                         `json:"title"`
	Description           string                         `json:"description,omitempty"`
	Caption               string                         `json:"caption,omitempty"`
	ParseMode             string                         `json:"parse_mode,omitempty"`
	CaptionEntities       []telegram.MessageEntity       `json:"caption_entities,omitempty"`
	ShowCaptionAboveMedia bool                           `json:"show_caption_above_media,omitempty"`
	ReplyMarkup           *telegram.InlineKeyboardMarkup `json:"reply_markup,omitempty"`
	InputMessageContent   InputMessageContent            `json:"input_message_content,omitempty"`
}

// InlineQueryResultCachedVoice represents a cached voice inline query result.
type InlineQueryResultCachedVoice struct {
	Type                string                         `json:"type"`
	ID                  string                         `json:"id"`
	VoiceFileID         string                         `json:"voice_file_id"`
	Title               string                         `json:"title"`
	Caption             string                         `json:"caption,omitempty"`
	ParseMode           string                         `json:"parse_mode,omitempty"`
	CaptionEntities     []telegram.MessageEntity       `json:"caption_entities,omitempty"`
	ReplyMarkup         *telegram.InlineKeyboardMarkup `json:"reply_markup,omitempty"`
	InputMessageContent InputMessageContent            `json:"input_message_content,omitempty"`
}

// InlineQueryResultCachedAudio represents a cached audio inline query result.
type InlineQueryResultCachedAudio struct {
	Type                string                         `json:"type"`
	ID                  string                         `json:"id"`
	AudioFileID         string                         `json:"audio_file_id"`
	Caption             string                         `json:"caption,omitempty"`
	ParseMode           string                         `json:"parse_mode,omitempty"`
	CaptionEntities     []telegram.MessageEntity       `json:"caption_entities,omitempty"`
	ReplyMarkup         *telegram.InlineKeyboardMarkup `json:"reply_markup,omitempty"`
	InputMessageContent InputMessageContent            `json:"input_message_content,omitempty"`
}

func (InlineQueryResultPhoto) inlineQueryResult()          {}
func (InlineQueryResultGif) inlineQueryResult()            {}
func (InlineQueryResultMpeg4Gif) inlineQueryResult()       {}
func (InlineQueryResultVideo) inlineQueryResult()          {}
func (InlineQueryResultAudio) inlineQueryResult()          {}
func (InlineQueryResultVoice) inlineQueryResult()          {}
func (InlineQueryResultDocument) inlineQueryResult()       {}
func (InlineQueryResultCachedPhoto) inlineQueryResult()    {}
func (InlineQueryResultCachedGif) inlineQueryResult()      {}
func (InlineQueryResultCachedMpeg4Gif) inlineQueryResult() {}
func (InlineQueryResultCachedSticker) inlineQueryResult()  {}
func (InlineQueryResultCachedDocument) inlineQueryResult() {}
func (InlineQueryResultCachedVideo) inlineQueryResult()    {}
func (InlineQueryResultCachedVoice) inlineQueryResult()    {}
func (InlineQueryResultCachedAudio) inlineQueryResult()    {}

// InlinePhoto creates a photo inline query result.
func InlinePhoto(id string, photoURL string, thumbnailURL string) InlineQueryResultPhoto {
	return InlineQueryResultPhoto{Type: inlineQueryResultPhotoType, ID: id, PhotoURL: photoURL, ThumbnailURL: thumbnailURL}
}

// InlineGif creates a GIF inline query result.
func InlineGif(id string, gifURL string, thumbnailURL string) InlineQueryResultGif {
	return InlineQueryResultGif{Type: inlineQueryResultGifType, ID: id, GifURL: gifURL, ThumbnailURL: thumbnailURL}
}

// InlineMpeg4Gif creates an MPEG-4 GIF inline query result.
func InlineMpeg4Gif(id string, mpeg4URL string, thumbnailURL string) InlineQueryResultMpeg4Gif {
	return InlineQueryResultMpeg4Gif{Type: inlineQueryResultMpeg4GifType, ID: id, Mpeg4URL: mpeg4URL, ThumbnailURL: thumbnailURL}
}

// InlineVideo creates a video inline query result.
func InlineVideo(id string, videoURL string, mimeType string, thumbnailURL string, title string) InlineQueryResultVideo {
	return InlineQueryResultVideo{Type: inlineQueryResultVideoType, ID: id, VideoURL: videoURL, MimeType: mimeType, ThumbnailURL: thumbnailURL, Title: title}
}

// InlineAudio creates an audio inline query result.
func InlineAudio(id string, audioURL string, title string) InlineQueryResultAudio {
	return InlineQueryResultAudio{Type: inlineQueryResultAudioType, ID: id, AudioURL: audioURL, Title: title}
}

// InlineVoice creates a voice inline query result.
func InlineVoice(id string, voiceURL string, title string) InlineQueryResultVoice {
	return InlineQueryResultVoice{Type: inlineQueryResultVoiceType, ID: id, VoiceURL: voiceURL, Title: title}
}

// InlineDocument creates a document inline query result.
func InlineDocument(id string, title string, documentURL string, mimeType string) InlineQueryResultDocument {
	return InlineQueryResultDocument{Type: inlineQueryResultDocumentType, ID: id, Title: title, DocumentURL: documentURL, MimeType: mimeType}
}

// InlineCachedPhoto creates a cached photo inline query result.
func InlineCachedPhoto(id string, fileID string) InlineQueryResultCachedPhoto {
	return InlineQueryResultCachedPhoto{Type: inlineQueryResultPhotoType, ID: id, PhotoFileID: fileID}
}

// InlineCachedGif creates a cached GIF inline query result.
func InlineCachedGif(id string, fileID string) InlineQueryResultCachedGif {
	return InlineQueryResultCachedGif{Type: inlineQueryResultGifType, ID: id, GifFileID: fileID}
}

// InlineCachedMpeg4Gif creates a cached MPEG-4 GIF inline query result.
func InlineCachedMpeg4Gif(id string, fileID string) InlineQueryResultCachedMpeg4Gif {
	return InlineQueryResultCachedMpeg4Gif{Type: inlineQueryResultMpeg4GifType, ID: id, Mpeg4FileID: fileID}
}

// InlineCachedSticker creates a cached sticker inline query result.
func InlineCachedSticker(id string, fileID string) InlineQueryResultCachedSticker {
	return InlineQueryResultCachedSticker{Type: inlineQueryResultStickerType, ID: id, StickerFileID: fileID}
}

// InlineCachedDocument creates a cached document inline query result.
func InlineCachedDocument(id string, fileID string, title string) InlineQueryResultCachedDocument {
	return InlineQueryResultCachedDocument{Type: inlineQueryResultDocumentType, ID: id, DocumentFileID: fileID, Title: title}
}

// InlineCachedVideo creates a cached video inline query result.
func InlineCachedVideo(id string, fileID string, title string) InlineQueryResultCachedVideo {
	return InlineQueryResultCachedVideo{Type: inlineQueryResultVideoType, ID: id, VideoFileID: fileID, Title: title}
}

// InlineCachedVoice creates a cached voice inline query result.
func InlineCachedVoice(id string, fileID string, title string) InlineQueryResultCachedVoice {
	return InlineQueryResultCachedVoice{Type: inlineQueryResultVoiceType, ID: id, VoiceFileID: fileID, Title: title}
}

// InlineCachedAudio creates a cached audio inline query result.
func InlineCachedAudio(id string, fileID string) InlineQueryResultCachedAudio {
	return InlineQueryResultCachedAudio{Type: inlineQueryResultAudioType, ID: id, AudioFileID: fileID}
}

// MarshalJSON encodes result with the official photo inline result type.
func (result InlineQueryResultPhoto) MarshalJSON() ([]byte, error) {
	result.Type = inlineQueryResultPhotoType
	type payload InlineQueryResultPhoto
	return json.Marshal(payload(result))
}

// MarshalJSON encodes result with the official GIF inline result type.
func (result InlineQueryResultGif) MarshalJSON() ([]byte, error) {
	result.Type = inlineQueryResultGifType
	type payload InlineQueryResultGif
	return json.Marshal(payload(result))
}

// MarshalJSON encodes result with the official MPEG-4 GIF inline result type.
func (result InlineQueryResultMpeg4Gif) MarshalJSON() ([]byte, error) {
	result.Type = inlineQueryResultMpeg4GifType
	type payload InlineQueryResultMpeg4Gif
	return json.Marshal(payload(result))
}

// MarshalJSON encodes result with the official video inline result type.
func (result InlineQueryResultVideo) MarshalJSON() ([]byte, error) {
	result.Type = inlineQueryResultVideoType
	type payload InlineQueryResultVideo
	return json.Marshal(payload(result))
}

// MarshalJSON encodes result with the official audio inline result type.
func (result InlineQueryResultAudio) MarshalJSON() ([]byte, error) {
	result.Type = inlineQueryResultAudioType
	type payload InlineQueryResultAudio
	return json.Marshal(payload(result))
}

// MarshalJSON encodes result with the official voice inline result type.
func (result InlineQueryResultVoice) MarshalJSON() ([]byte, error) {
	result.Type = inlineQueryResultVoiceType
	type payload InlineQueryResultVoice
	return json.Marshal(payload(result))
}

// MarshalJSON encodes result with the official document inline result type.
func (result InlineQueryResultDocument) MarshalJSON() ([]byte, error) {
	result.Type = inlineQueryResultDocumentType
	type payload InlineQueryResultDocument
	return json.Marshal(payload(result))
}

// MarshalJSON encodes result with the official cached photo inline result type.
func (result InlineQueryResultCachedPhoto) MarshalJSON() ([]byte, error) {
	result.Type = inlineQueryResultPhotoType
	type payload InlineQueryResultCachedPhoto
	return json.Marshal(payload(result))
}

// MarshalJSON encodes result with the official cached GIF inline result type.
func (result InlineQueryResultCachedGif) MarshalJSON() ([]byte, error) {
	result.Type = inlineQueryResultGifType
	type payload InlineQueryResultCachedGif
	return json.Marshal(payload(result))
}

// MarshalJSON encodes result with the official cached MPEG-4 GIF inline result type.
func (result InlineQueryResultCachedMpeg4Gif) MarshalJSON() ([]byte, error) {
	result.Type = inlineQueryResultMpeg4GifType
	type payload InlineQueryResultCachedMpeg4Gif
	return json.Marshal(payload(result))
}

// MarshalJSON encodes result with the official cached sticker inline result type.
func (result InlineQueryResultCachedSticker) MarshalJSON() ([]byte, error) {
	result.Type = inlineQueryResultStickerType
	type payload InlineQueryResultCachedSticker
	return json.Marshal(payload(result))
}

// MarshalJSON encodes result with the official cached document inline result type.
func (result InlineQueryResultCachedDocument) MarshalJSON() ([]byte, error) {
	result.Type = inlineQueryResultDocumentType
	type payload InlineQueryResultCachedDocument
	return json.Marshal(payload(result))
}

// MarshalJSON encodes result with the official cached video inline result type.
func (result InlineQueryResultCachedVideo) MarshalJSON() ([]byte, error) {
	result.Type = inlineQueryResultVideoType
	type payload InlineQueryResultCachedVideo
	return json.Marshal(payload(result))
}

// MarshalJSON encodes result with the official cached voice inline result type.
func (result InlineQueryResultCachedVoice) MarshalJSON() ([]byte, error) {
	result.Type = inlineQueryResultVoiceType
	type payload InlineQueryResultCachedVoice
	return json.Marshal(payload(result))
}

// MarshalJSON encodes result with the official cached audio inline result type.
func (result InlineQueryResultCachedAudio) MarshalJSON() ([]byte, error) {
	result.Type = inlineQueryResultAudioType
	type payload InlineQueryResultCachedAudio
	return json.Marshal(payload(result))
}

func validateInlineQueryResultPhoto(result InlineQueryResultPhoto) error {
	if err := validateInlineID(result.ID, "inline query result photo id"); err != nil {
		return err
	}
	if err := validateRequiredInlineHTTPURL(result.PhotoURL, "photo_url"); err != nil {
		return err
	}
	if err := validateRequiredInlineHTTPURL(result.ThumbnailURL, "thumbnail_url"); err != nil {
		return err
	}
	if err := validateInlineCaptionFields(result.ParseMode, result.CaptionEntities); err != nil {
		return err
	}
	if err := validateInlineDimensions(map[string]int{"photo_width": result.PhotoWidth, "photo_height": result.PhotoHeight}); err != nil {
		return err
	}
	return validateOptionalInlineResultFields(result.ReplyMarkup, result.InputMessageContent, "", 0, 0)
}

func validateInlineQueryResultGif(result InlineQueryResultGif) error {
	if err := validateInlineID(result.ID, "inline query result gif id"); err != nil {
		return err
	}
	if err := validateRequiredInlineHTTPURL(result.GifURL, "gif_url"); err != nil {
		return err
	}
	if err := validateRequiredInlineHTTPURL(result.ThumbnailURL, "thumbnail_url"); err != nil {
		return err
	}
	if err := validateInlineCaptionFields(result.ParseMode, result.CaptionEntities); err != nil {
		return err
	}
	if err := validateInlineDimensions(map[string]int{"gif_width": result.GifWidth, "gif_height": result.GifHeight, "gif_duration": result.GifDuration}); err != nil {
		return err
	}
	return validateOptionalInlineResultFields(result.ReplyMarkup, result.InputMessageContent, "", 0, 0)
}

func validateInlineQueryResultMpeg4Gif(result InlineQueryResultMpeg4Gif) error {
	if err := validateInlineID(result.ID, "inline query result mpeg4_gif id"); err != nil {
		return err
	}
	if err := validateRequiredInlineHTTPURL(result.Mpeg4URL, "mpeg4_url"); err != nil {
		return err
	}
	if err := validateRequiredInlineHTTPURL(result.ThumbnailURL, "thumbnail_url"); err != nil {
		return err
	}
	if err := validateInlineCaptionFields(result.ParseMode, result.CaptionEntities); err != nil {
		return err
	}
	if err := validateInlineDimensions(map[string]int{"mpeg4_width": result.Mpeg4Width, "mpeg4_height": result.Mpeg4Height, "mpeg4_duration": result.Mpeg4Duration}); err != nil {
		return err
	}
	return validateOptionalInlineResultFields(result.ReplyMarkup, result.InputMessageContent, "", 0, 0)
}

func validateInlineQueryResultVideo(result InlineQueryResultVideo) error {
	if err := validateInlineID(result.ID, "inline query result video id"); err != nil {
		return err
	}
	if err := validateRequiredInlineHTTPURL(result.VideoURL, "video_url"); err != nil {
		return err
	}
	if strings.TrimSpace(result.MimeType) == "" {
		return stderrors.New("mime_type is required")
	}
	if err := validateRequiredInlineHTTPURL(result.ThumbnailURL, "thumbnail_url"); err != nil {
		return err
	}
	if strings.TrimSpace(result.Title) == "" {
		return stderrors.New("title is required")
	}
	if err := validateInlineCaptionFields(result.ParseMode, result.CaptionEntities); err != nil {
		return err
	}
	if err := validateInlineDimensions(map[string]int{"video_width": result.VideoWidth, "video_height": result.VideoHeight, "video_duration": result.VideoDuration}); err != nil {
		return err
	}
	return validateOptionalInlineResultFields(result.ReplyMarkup, result.InputMessageContent, "", 0, 0)
}

func validateInlineQueryResultAudio(result InlineQueryResultAudio) error {
	if err := validateInlineID(result.ID, "inline query result audio id"); err != nil {
		return err
	}
	if err := validateRequiredInlineHTTPURL(result.AudioURL, "audio_url"); err != nil {
		return err
	}
	if strings.TrimSpace(result.Title) == "" {
		return stderrors.New("title is required")
	}
	if err := validateInlineCaptionFields(result.ParseMode, result.CaptionEntities); err != nil {
		return err
	}
	if err := validateInlineDimensions(map[string]int{"audio_duration": result.AudioDuration}); err != nil {
		return err
	}
	return validateOptionalInlineResultFields(result.ReplyMarkup, result.InputMessageContent, "", 0, 0)
}

func validateInlineQueryResultVoice(result InlineQueryResultVoice) error {
	if err := validateInlineID(result.ID, "inline query result voice id"); err != nil {
		return err
	}
	if err := validateRequiredInlineHTTPURL(result.VoiceURL, "voice_url"); err != nil {
		return err
	}
	if strings.TrimSpace(result.Title) == "" {
		return stderrors.New("title is required")
	}
	if err := validateInlineCaptionFields(result.ParseMode, result.CaptionEntities); err != nil {
		return err
	}
	if err := validateInlineDimensions(map[string]int{"voice_duration": result.VoiceDuration}); err != nil {
		return err
	}
	return validateOptionalInlineResultFields(result.ReplyMarkup, result.InputMessageContent, "", 0, 0)
}

func validateInlineQueryResultDocument(result InlineQueryResultDocument) error {
	if err := validateInlineID(result.ID, "inline query result document id"); err != nil {
		return err
	}
	if strings.TrimSpace(result.Title) == "" {
		return stderrors.New("title is required")
	}
	if err := validateRequiredInlineHTTPURL(result.DocumentURL, "document_url"); err != nil {
		return err
	}
	if strings.TrimSpace(result.MimeType) == "" {
		return stderrors.New("mime_type is required")
	}
	if err := validateInlineCaptionFields(result.ParseMode, result.CaptionEntities); err != nil {
		return err
	}
	return validateOptionalInlineResultFields(result.ReplyMarkup, result.InputMessageContent, result.ThumbnailURL, result.ThumbnailWidth, result.ThumbnailHeight)
}

func validateInlineQueryResultCachedPhoto(result InlineQueryResultCachedPhoto) error {
	if err := validateInlineID(result.ID, "inline query result cached photo id"); err != nil {
		return err
	}
	if err := validateFileID(result.PhotoFileID, "photo_file_id"); err != nil {
		return err
	}
	if err := validateInlineCaptionFields(result.ParseMode, result.CaptionEntities); err != nil {
		return err
	}
	return validateOptionalInlineResultFields(result.ReplyMarkup, result.InputMessageContent, "", 0, 0)
}

func validateInlineQueryResultCachedGif(result InlineQueryResultCachedGif) error {
	if err := validateInlineID(result.ID, "inline query result cached gif id"); err != nil {
		return err
	}
	if err := validateFileID(result.GifFileID, "gif_file_id"); err != nil {
		return err
	}
	if err := validateInlineCaptionFields(result.ParseMode, result.CaptionEntities); err != nil {
		return err
	}
	return validateOptionalInlineResultFields(result.ReplyMarkup, result.InputMessageContent, "", 0, 0)
}

func validateInlineQueryResultCachedMpeg4Gif(result InlineQueryResultCachedMpeg4Gif) error {
	if err := validateInlineID(result.ID, "inline query result cached mpeg4_gif id"); err != nil {
		return err
	}
	if err := validateFileID(result.Mpeg4FileID, "mpeg4_file_id"); err != nil {
		return err
	}
	if err := validateInlineCaptionFields(result.ParseMode, result.CaptionEntities); err != nil {
		return err
	}
	return validateOptionalInlineResultFields(result.ReplyMarkup, result.InputMessageContent, "", 0, 0)
}

func validateInlineQueryResultCachedSticker(result InlineQueryResultCachedSticker) error {
	if err := validateInlineID(result.ID, "inline query result cached sticker id"); err != nil {
		return err
	}
	if err := validateFileID(result.StickerFileID, "sticker_file_id"); err != nil {
		return err
	}
	return validateOptionalInlineResultFields(result.ReplyMarkup, result.InputMessageContent, "", 0, 0)
}

func validateInlineQueryResultCachedDocument(result InlineQueryResultCachedDocument) error {
	if err := validateInlineID(result.ID, "inline query result cached document id"); err != nil {
		return err
	}
	if strings.TrimSpace(result.Title) == "" {
		return stderrors.New("title is required")
	}
	if err := validateFileID(result.DocumentFileID, "document_file_id"); err != nil {
		return err
	}
	if err := validateInlineCaptionFields(result.ParseMode, result.CaptionEntities); err != nil {
		return err
	}
	return validateOptionalInlineResultFields(result.ReplyMarkup, result.InputMessageContent, "", 0, 0)
}

func validateInlineQueryResultCachedVideo(result InlineQueryResultCachedVideo) error {
	if err := validateInlineID(result.ID, "inline query result cached video id"); err != nil {
		return err
	}
	if strings.TrimSpace(result.Title) == "" {
		return stderrors.New("title is required")
	}
	if err := validateFileID(result.VideoFileID, "video_file_id"); err != nil {
		return err
	}
	if err := validateInlineCaptionFields(result.ParseMode, result.CaptionEntities); err != nil {
		return err
	}
	return validateOptionalInlineResultFields(result.ReplyMarkup, result.InputMessageContent, "", 0, 0)
}

func validateInlineQueryResultCachedVoice(result InlineQueryResultCachedVoice) error {
	if err := validateInlineID(result.ID, "inline query result cached voice id"); err != nil {
		return err
	}
	if strings.TrimSpace(result.Title) == "" {
		return stderrors.New("title is required")
	}
	if err := validateFileID(result.VoiceFileID, "voice_file_id"); err != nil {
		return err
	}
	if err := validateInlineCaptionFields(result.ParseMode, result.CaptionEntities); err != nil {
		return err
	}
	return validateOptionalInlineResultFields(result.ReplyMarkup, result.InputMessageContent, "", 0, 0)
}

func validateInlineQueryResultCachedAudio(result InlineQueryResultCachedAudio) error {
	if err := validateInlineID(result.ID, "inline query result cached audio id"); err != nil {
		return err
	}
	if err := validateFileID(result.AudioFileID, "audio_file_id"); err != nil {
		return err
	}
	if err := validateInlineCaptionFields(result.ParseMode, result.CaptionEntities); err != nil {
		return err
	}
	return validateOptionalInlineResultFields(result.ReplyMarkup, result.InputMessageContent, "", 0, 0)
}

func validateInlineID(id string, field string) error {
	if strings.TrimSpace(id) == "" {
		return stderrors.New(field + " is required")
	}
	return nil
}

func validateRequiredInlineHTTPURL(rawURL string, field string) error {
	if strings.TrimSpace(rawURL) == "" {
		return stderrors.New(field + " is required")
	}
	return validateInlineHTTPURL(rawURL, field)
}

func validateFileID(fileID string, field string) error {
	if strings.TrimSpace(fileID) == "" {
		return stderrors.New(field + " is required")
	}
	return nil
}

func validateInlineCaptionFields(parseMode string, captionEntities []telegram.MessageEntity) error {
	if err := validateEntityFormatting(parseMode, captionEntities); err != nil {
		return err
	}
	return nil
}

func validateInlineDimensions(fields map[string]int) error {
	for field, value := range fields {
		if value < 0 {
			return fmt.Errorf("%s must not be negative", field)
		}
	}
	return nil
}
