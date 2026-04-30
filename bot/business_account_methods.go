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

// InputStoryContent marks story content accepted by postStory and editStory.
type InputStoryContent interface {
	inputStoryContent()
}

// InputStoryContentPhoto describes a photo story content item.
type InputStoryContentPhoto struct {
	Type  string  `json:"type"`
	Photo FileRef `json:"photo"`
}

// InputStoryContentVideo describes a video story content item.
type InputStoryContentVideo struct {
	Type                string  `json:"type"`
	Video               FileRef `json:"video"`
	Duration            float64 `json:"duration,omitempty"`
	CoverFrameTimestamp float64 `json:"cover_frame_timestamp,omitempty"`
	IsAnimation         bool    `json:"is_animation,omitempty"`
}

func (InputStoryContentPhoto) inputStoryContent() {}
func (InputStoryContentVideo) inputStoryContent() {}

// StoryPhoto creates photo content for a business story.
func StoryPhoto(photo FileRef) InputStoryContentPhoto {
	return InputStoryContentPhoto{Type: "photo", Photo: photo}
}

// StoryVideo creates video content for a business story.
func StoryVideo(video FileRef) InputStoryContentVideo {
	return InputStoryContentVideo{Type: "video", Video: video}
}

// ReadBusinessMessageParams contains supported parameters for readBusinessMessage.
type ReadBusinessMessageParams struct {
	BusinessConnectionID string `json:"business_connection_id"`
	ChatID               ChatID `json:"chat_id"`
	MessageID            int64  `json:"message_id"`
}

// SetBusinessAccountNameParams contains supported parameters for setBusinessAccountName.
type SetBusinessAccountNameParams struct {
	BusinessConnectionID string `json:"business_connection_id"`
	FirstName            string `json:"first_name"`
	LastName             string `json:"last_name,omitempty"`
}

// SetBusinessAccountUsernameParams contains supported parameters for setBusinessAccountUsername.
type SetBusinessAccountUsernameParams struct {
	BusinessConnectionID string `json:"business_connection_id"`
	Username             string `json:"username,omitempty"`
}

// SetBusinessAccountBioParams contains supported parameters for setBusinessAccountBio.
type SetBusinessAccountBioParams struct {
	BusinessConnectionID string `json:"business_connection_id"`
	Bio                  string `json:"bio,omitempty"`
}

// SetBusinessAccountProfilePhotoParams contains supported parameters for setBusinessAccountProfilePhoto.
type SetBusinessAccountProfilePhotoParams struct {
	BusinessConnectionID string            `json:"business_connection_id"`
	Photo                InputProfilePhoto `json:"photo"`
	IsPublic             bool              `json:"is_public,omitempty"`
}

// RemoveBusinessAccountProfilePhotoParams contains supported parameters for removeBusinessAccountProfilePhoto.
type RemoveBusinessAccountProfilePhotoParams struct {
	BusinessConnectionID string `json:"business_connection_id"`
	IsPublic             bool   `json:"is_public,omitempty"`
}

// SetBusinessAccountGiftSettingsParams contains supported parameters for setBusinessAccountGiftSettings.
type SetBusinessAccountGiftSettingsParams struct {
	BusinessConnectionID string                     `json:"business_connection_id"`
	ShowGiftButton       bool                       `json:"show_gift_button"`
	AcceptedGiftTypes    telegram.AcceptedGiftTypes `json:"accepted_gift_types"`
}

// PostStoryParams contains supported parameters for postStory.
type PostStoryParams struct {
	BusinessConnectionID string                   `json:"business_connection_id"`
	Content              InputStoryContent        `json:"content"`
	ActivePeriod         int                      `json:"active_period"`
	Caption              string                   `json:"caption,omitempty"`
	ParseMode            string                   `json:"parse_mode,omitempty"`
	CaptionEntities      []telegram.MessageEntity `json:"caption_entities,omitempty"`
	Areas                []telegram.StoryArea     `json:"areas,omitempty"`
	PostToChatPage       bool                     `json:"post_to_chat_page,omitempty"`
	ProtectContent       bool                     `json:"protect_content,omitempty"`
}

// EditStoryParams contains supported parameters for editStory.
type EditStoryParams struct {
	BusinessConnectionID string                   `json:"business_connection_id"`
	StoryID              int64                    `json:"story_id"`
	Content              InputStoryContent        `json:"content"`
	Caption              string                   `json:"caption,omitempty"`
	ParseMode            string                   `json:"parse_mode,omitempty"`
	CaptionEntities      []telegram.MessageEntity `json:"caption_entities,omitempty"`
	Areas                []telegram.StoryArea     `json:"areas,omitempty"`
}

// DeleteStoryParams contains supported parameters for deleteStory.
type DeleteStoryParams struct {
	BusinessConnectionID string `json:"business_connection_id"`
	StoryID              int64  `json:"story_id"`
}

// ApproveSuggestedPostParams contains supported parameters for approveSuggestedPost.
type ApproveSuggestedPostParams struct {
	ChatID    ChatID `json:"chat_id"`
	MessageID int64  `json:"message_id"`
	SendDate  int64  `json:"send_date,omitempty"`
}

// DeclineSuggestedPostParams contains supported parameters for declineSuggestedPost.
type DeclineSuggestedPostParams struct {
	ChatID    ChatID `json:"chat_id"`
	MessageID int64  `json:"message_id"`
	Comment   string `json:"comment,omitempty"`
}

// ReadBusinessMessage marks an incoming business message as read.
func (b *Bot) ReadBusinessMessage(ctx context.Context, params ReadBusinessMessageParams) (bool, error) {
	if err := params.validate(); err != nil {
		return false, err
	}
	var result bool
	if err := b.call(ctx, "readBusinessMessage", params, &result); err != nil {
		return false, err
	}
	return result, nil
}

// SetBusinessAccountName changes the first and last name of a business account.
func (b *Bot) SetBusinessAccountName(ctx context.Context, params SetBusinessAccountNameParams) (bool, error) {
	if err := params.validate(); err != nil {
		return false, err
	}
	var result bool
	if err := b.call(ctx, "setBusinessAccountName", params, &result); err != nil {
		return false, err
	}
	return result, nil
}

// SetBusinessAccountUsername changes or removes the username of a business account.
func (b *Bot) SetBusinessAccountUsername(ctx context.Context, params SetBusinessAccountUsernameParams) (bool, error) {
	if err := params.validate(); err != nil {
		return false, err
	}
	var result bool
	if err := b.call(ctx, "setBusinessAccountUsername", params, &result); err != nil {
		return false, err
	}
	return result, nil
}

// SetBusinessAccountBio changes or removes the bio of a business account.
func (b *Bot) SetBusinessAccountBio(ctx context.Context, params SetBusinessAccountBioParams) (bool, error) {
	if err := params.validate(); err != nil {
		return false, err
	}
	var result bool
	if err := b.call(ctx, "setBusinessAccountBio", params, &result); err != nil {
		return false, err
	}
	return result, nil
}

// SetBusinessAccountProfilePhoto uploads and sets a profile photo of a business account.
func (b *Bot) SetBusinessAccountProfilePhoto(ctx context.Context, params SetBusinessAccountProfilePhotoParams) (bool, error) {
	if err := params.validate(); err != nil {
		return false, err
	}
	fields, files, err := params.multipart()
	if err != nil {
		return false, err
	}
	var result bool
	if err := b.callMultipart(ctx, "setBusinessAccountProfilePhoto", fields, files, &result); err != nil {
		return false, err
	}
	return result, nil
}

// RemoveBusinessAccountProfilePhoto removes a profile photo of a business account.
func (b *Bot) RemoveBusinessAccountProfilePhoto(ctx context.Context, params RemoveBusinessAccountProfilePhotoParams) (bool, error) {
	if err := params.validate(); err != nil {
		return false, err
	}
	var result bool
	if err := b.call(ctx, "removeBusinessAccountProfilePhoto", params, &result); err != nil {
		return false, err
	}
	return result, nil
}

// SetBusinessAccountGiftSettings changes incoming gift privacy settings of a business account.
func (b *Bot) SetBusinessAccountGiftSettings(ctx context.Context, params SetBusinessAccountGiftSettingsParams) (bool, error) {
	if err := params.validate(); err != nil {
		return false, err
	}
	var result bool
	if err := b.call(ctx, "setBusinessAccountGiftSettings", params, &result); err != nil {
		return false, err
	}
	return result, nil
}

// PostStory posts a story on behalf of a business account.
func (b *Bot) PostStory(ctx context.Context, params PostStoryParams) (*telegram.Story, error) {
	if err := params.validate(); err != nil {
		return nil, err
	}
	fields, files, err := params.multipart()
	if err != nil {
		return nil, err
	}
	var story telegram.Story
	if err := b.callMultipart(ctx, "postStory", fields, files, &story); err != nil {
		return nil, err
	}
	return &story, nil
}

// EditStory edits a story previously posted by the bot on behalf of a business account.
func (b *Bot) EditStory(ctx context.Context, params EditStoryParams) (*telegram.Story, error) {
	if err := params.validate(); err != nil {
		return nil, err
	}
	fields, files, err := params.multipart()
	if err != nil {
		return nil, err
	}
	var story telegram.Story
	if err := b.callMultipart(ctx, "editStory", fields, files, &story); err != nil {
		return nil, err
	}
	return &story, nil
}

// DeleteStory deletes a story previously posted by the bot on behalf of a business account.
func (b *Bot) DeleteStory(ctx context.Context, params DeleteStoryParams) (bool, error) {
	if err := params.validate(); err != nil {
		return false, err
	}
	var result bool
	if err := b.call(ctx, "deleteStory", params, &result); err != nil {
		return false, err
	}
	return result, nil
}

// ApproveSuggestedPost approves a suggested post in a direct messages chat.
func (b *Bot) ApproveSuggestedPost(ctx context.Context, params ApproveSuggestedPostParams) (bool, error) {
	if err := params.validate(); err != nil {
		return false, err
	}
	var result bool
	if err := b.call(ctx, "approveSuggestedPost", params, &result); err != nil {
		return false, err
	}
	return result, nil
}

// DeclineSuggestedPost declines a suggested post in a direct messages chat.
func (b *Bot) DeclineSuggestedPost(ctx context.Context, params DeclineSuggestedPostParams) (bool, error) {
	if err := params.validate(); err != nil {
		return false, err
	}
	var result bool
	if err := b.call(ctx, "declineSuggestedPost", params, &result); err != nil {
		return false, err
	}
	return result, nil
}

func (params ReadBusinessMessageParams) validate() error {
	if err := validateBusinessConnectionID(params.BusinessConnectionID); err != nil {
		return err
	}
	if !params.ChatID.valid() {
		return stderrors.New("chat_id is required")
	}
	if params.MessageID <= 0 {
		return stderrors.New("message_id must be greater than zero")
	}
	return nil
}

func (params SetBusinessAccountNameParams) validate() error {
	if err := validateBusinessConnectionID(params.BusinessConnectionID); err != nil {
		return err
	}
	if strings.TrimSpace(params.FirstName) == "" {
		return stderrors.New("first_name is required")
	}
	return nil
}

func (params SetBusinessAccountUsernameParams) validate() error {
	return validateBusinessConnectionID(params.BusinessConnectionID)
}

func (params SetBusinessAccountBioParams) validate() error {
	return validateBusinessConnectionID(params.BusinessConnectionID)
}

func (params SetBusinessAccountProfilePhotoParams) validate() error {
	if err := validateBusinessConnectionID(params.BusinessConnectionID); err != nil {
		return err
	}
	_, _, err := profilePhotoPayloadAndUpload(params.Photo)
	return err
}

func (params SetBusinessAccountProfilePhotoParams) multipart() (map[string]string, map[string]UploadFile, error) {
	payload, files, err := profilePhotoPayloadAndUpload(params.Photo)
	if err != nil {
		return nil, nil, err
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, nil, err
	}
	fields := map[string]string{
		"business_connection_id": params.BusinessConnectionID,
		"photo":                  string(body),
	}
	boolField(fields, "is_public", params.IsPublic)
	return fields, files, nil
}

func (params RemoveBusinessAccountProfilePhotoParams) validate() error {
	return validateBusinessConnectionID(params.BusinessConnectionID)
}

func (params SetBusinessAccountGiftSettingsParams) validate() error {
	if err := validateBusinessConnectionID(params.BusinessConnectionID); err != nil {
		return err
	}
	return telegram.ValidateAcceptedGiftTypes(params.AcceptedGiftTypes)
}

func (params PostStoryParams) validate() error {
	if err := validateBusinessConnectionID(params.BusinessConnectionID); err != nil {
		return err
	}
	if err := validateStoryActivePeriod(params.ActivePeriod); err != nil {
		return err
	}
	return validateStoryPayload(params.Content, params.ParseMode, params.CaptionEntities, params.Areas)
}

func (params PostStoryParams) multipart() (map[string]string, map[string]UploadFile, error) {
	content, files, err := storyContentPayloadAndUpload(params.Content)
	if err != nil {
		return nil, nil, err
	}
	fields, err := storyMultipartFields(params.BusinessConnectionID, content, params.Caption, params.ParseMode, params.CaptionEntities, params.Areas)
	if err != nil {
		return nil, nil, err
	}
	fields["active_period"] = strconv.Itoa(params.ActivePeriod)
	boolField(fields, "post_to_chat_page", params.PostToChatPage)
	boolField(fields, "protect_content", params.ProtectContent)
	return fields, files, nil
}

func (params EditStoryParams) validate() error {
	if err := validateBusinessConnectionID(params.BusinessConnectionID); err != nil {
		return err
	}
	if params.StoryID <= 0 {
		return stderrors.New("story_id must be greater than zero")
	}
	return validateStoryPayload(params.Content, params.ParseMode, params.CaptionEntities, params.Areas)
}

func (params EditStoryParams) multipart() (map[string]string, map[string]UploadFile, error) {
	content, files, err := storyContentPayloadAndUpload(params.Content)
	if err != nil {
		return nil, nil, err
	}
	fields, err := storyMultipartFields(params.BusinessConnectionID, content, params.Caption, params.ParseMode, params.CaptionEntities, params.Areas)
	if err != nil {
		return nil, nil, err
	}
	fields["story_id"] = strconv.FormatInt(params.StoryID, 10)
	return fields, files, nil
}

func (params DeleteStoryParams) validate() error {
	if err := validateBusinessConnectionID(params.BusinessConnectionID); err != nil {
		return err
	}
	if params.StoryID <= 0 {
		return stderrors.New("story_id must be greater than zero")
	}
	return nil
}

func (params ApproveSuggestedPostParams) validate() error {
	if !params.ChatID.valid() {
		return stderrors.New("chat_id is required")
	}
	if params.MessageID <= 0 {
		return stderrors.New("message_id must be greater than zero")
	}
	if params.SendDate < 0 {
		return stderrors.New("send_date must not be negative")
	}
	return nil
}

func (params DeclineSuggestedPostParams) validate() error {
	if !params.ChatID.valid() {
		return stderrors.New("chat_id is required")
	}
	if params.MessageID <= 0 {
		return stderrors.New("message_id must be greater than zero")
	}
	return nil
}

func validateBusinessConnectionID(id string) error {
	if strings.TrimSpace(id) == "" {
		return stderrors.New("business_connection_id is required")
	}
	return nil
}

func validateStoryActivePeriod(period int) error {
	switch period {
	case 6 * 3600, 12 * 3600, 86400, 2 * 86400:
		return nil
	default:
		return stderrors.New("active_period must be one of 21600, 43200, 86400, or 172800")
	}
}

func validateStoryPayload(content InputStoryContent, parseMode string, captionEntities []telegram.MessageEntity, areas []telegram.StoryArea) error {
	if err := validateInputStoryContent(content); err != nil {
		return err
	}
	if err := validateCaptionFormatting(parseMode, captionEntities); err != nil {
		return err
	}
	return telegram.ValidateStoryAreas(areas)
}

func validateInputStoryContent(content InputStoryContent) error {
	switch item := content.(type) {
	case nil:
		return stderrors.New("content is required")
	case InputStoryContentPhoto:
		return validateInputStoryContentPhoto(item)
	case *InputStoryContentPhoto:
		if item == nil {
			return stderrors.New("content is required")
		}
		return validateInputStoryContentPhoto(*item)
	case InputStoryContentVideo:
		return validateInputStoryContentVideo(item)
	case *InputStoryContentVideo:
		if item == nil {
			return stderrors.New("content is required")
		}
		return validateInputStoryContentVideo(*item)
	default:
		return stderrors.New("unsupported story content type")
	}
}

func validateInputStoryContentPhoto(content InputStoryContentPhoto) error {
	if err := validateInputMediaType(content.Type, "photo"); err != nil {
		return err
	}
	return validateProfilePhotoUpload(content.Photo, "photo")
}

func validateInputStoryContentVideo(content InputStoryContentVideo) error {
	if err := validateInputMediaType(content.Type, "video"); err != nil {
		return err
	}
	if err := validateProfilePhotoUpload(content.Video, "video"); err != nil {
		return err
	}
	if content.Duration < 0 {
		return stderrors.New("duration must not be negative")
	}
	if content.CoverFrameTimestamp < 0 {
		return stderrors.New("cover_frame_timestamp must not be negative")
	}
	return nil
}

func storyContentPayloadAndUpload(content InputStoryContent) (map[string]any, map[string]UploadFile, error) {
	switch item := content.(type) {
	case InputStoryContentPhoto:
		return storyPhotoPayload(item)
	case *InputStoryContentPhoto:
		if item == nil {
			return nil, nil, stderrors.New("content is required")
		}
		return storyPhotoPayload(*item)
	case InputStoryContentVideo:
		return storyVideoPayload(item)
	case *InputStoryContentVideo:
		if item == nil {
			return nil, nil, stderrors.New("content is required")
		}
		return storyVideoPayload(*item)
	default:
		return nil, nil, stderrors.New("unsupported story content type")
	}
}

func storyPhotoPayload(content InputStoryContentPhoto) (map[string]any, map[string]UploadFile, error) {
	if err := validateInputStoryContentPhoto(content); err != nil {
		return nil, nil, err
	}
	return map[string]any{"type": "photo", "photo": "attach://story_photo"}, map[string]UploadFile{"story_photo": content.Photo.upload}, nil
}

func storyVideoPayload(content InputStoryContentVideo) (map[string]any, map[string]UploadFile, error) {
	if err := validateInputStoryContentVideo(content); err != nil {
		return nil, nil, err
	}
	payload := map[string]any{"type": "video", "video": "attach://story_video"}
	if content.Duration > 0 {
		payload["duration"] = content.Duration
	}
	if content.CoverFrameTimestamp > 0 {
		payload["cover_frame_timestamp"] = content.CoverFrameTimestamp
	}
	if content.IsAnimation {
		payload["is_animation"] = true
	}
	return payload, map[string]UploadFile{"story_video": content.Video.upload}, nil
}

func storyMultipartFields(businessConnectionID string, content map[string]any, caption string, parseMode string, captionEntities []telegram.MessageEntity, areas []telegram.StoryArea) (map[string]string, error) {
	contentBody, err := json.Marshal(content)
	if err != nil {
		return nil, err
	}
	fields := map[string]string{
		"business_connection_id": businessConnectionID,
		"content":                string(contentBody),
	}
	stringField(fields, "caption", caption)
	stringField(fields, "parse_mode", parseMode)
	if err := captionEntitiesField(fields, captionEntities); err != nil {
		return nil, err
	}
	if len(areas) > 0 {
		body, err := json.Marshal(areas)
		if err != nil {
			return nil, fmt.Errorf("marshal story areas: %w", err)
		}
		fields["areas"] = string(body)
	}
	return fields, nil
}
