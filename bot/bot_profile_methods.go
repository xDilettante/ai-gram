package bot

import (
	"context"
	"encoding/json"
	stderrors "errors"
	"strings"

	"github.com/xDilettante/ai-gram/telegram"
)

// SetMyNameParams contains supported parameters for setMyName.
type SetMyNameParams struct {
	Name         string `json:"name,omitempty"`
	LanguageCode string `json:"language_code,omitempty"`
}

// GetMyNameParams contains supported parameters for getMyName.
type GetMyNameParams struct {
	LanguageCode string `json:"language_code,omitempty"`
}

// SetMyDescriptionParams contains supported parameters for setMyDescription.
type SetMyDescriptionParams struct {
	Description  string `json:"description,omitempty"`
	LanguageCode string `json:"language_code,omitempty"`
}

// GetMyDescriptionParams contains supported parameters for getMyDescription.
type GetMyDescriptionParams struct {
	LanguageCode string `json:"language_code,omitempty"`
}

// SetMyShortDescriptionParams contains supported parameters for setMyShortDescription.
type SetMyShortDescriptionParams struct {
	ShortDescription string `json:"short_description,omitempty"`
	LanguageCode     string `json:"language_code,omitempty"`
}

// GetMyShortDescriptionParams contains supported parameters for getMyShortDescription.
type GetMyShortDescriptionParams struct {
	LanguageCode string `json:"language_code,omitempty"`
}

// GetMyDefaultAdministratorRightsParams contains supported parameters for getMyDefaultAdministratorRights.
type GetMyDefaultAdministratorRightsParams struct {
	ForChannels bool `json:"for_channels,omitempty"`
}

// InputProfilePhoto describes a bot profile photo accepted by setMyProfilePhoto.
type InputProfilePhoto interface {
	inputProfilePhoto()
}

// InputProfilePhotoStatic describes a static JPG bot profile photo.
type InputProfilePhotoStatic struct {
	Type  string  `json:"type"`
	Photo FileRef `json:"photo"`
}

// InputProfilePhotoAnimated describes an animated MPEG4 bot profile photo.
type InputProfilePhotoAnimated struct {
	Type               string  `json:"type"`
	Animation          FileRef `json:"animation"`
	MainFrameTimestamp float64 `json:"main_frame_timestamp,omitempty"`
}

func (InputProfilePhotoStatic) inputProfilePhoto()   {}
func (InputProfilePhotoAnimated) inputProfilePhoto() {}

// ProfilePhotoStatic creates a static JPG input profile photo.
func ProfilePhotoStatic(photo FileRef) InputProfilePhotoStatic {
	return InputProfilePhotoStatic{Type: "static", Photo: photo}
}

// ProfilePhotoAnimated creates an animated MPEG4 input profile photo.
func ProfilePhotoAnimated(animation FileRef) InputProfilePhotoAnimated {
	return InputProfilePhotoAnimated{Type: "animated", Animation: animation}
}

// SetMyProfilePhotoParams contains supported parameters for setMyProfilePhoto.
type SetMyProfilePhotoParams struct {
	Photo InputProfilePhoto `json:"photo"`
}

// RemoveMyProfilePhotoParams contains supported parameters for removeMyProfilePhoto.
type RemoveMyProfilePhotoParams struct{}

// SetMyName changes the bot's name.
func (b *Bot) SetMyName(ctx context.Context, params SetMyNameParams) (bool, error) {
	var result bool
	if err := b.call(ctx, "setMyName", params, &result); err != nil {
		return false, err
	}
	return result, nil
}

// GetMyName gets the bot's current name for a language.
func (b *Bot) GetMyName(ctx context.Context, params GetMyNameParams) (*telegram.BotName, error) {
	var result telegram.BotName
	if err := b.call(ctx, "getMyName", params, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// SetMyDescription changes the bot's description.
func (b *Bot) SetMyDescription(ctx context.Context, params SetMyDescriptionParams) (bool, error) {
	var result bool
	if err := b.call(ctx, "setMyDescription", params, &result); err != nil {
		return false, err
	}
	return result, nil
}

// GetMyDescription gets the bot's current description for a language.
func (b *Bot) GetMyDescription(ctx context.Context, params GetMyDescriptionParams) (*telegram.BotDescription, error) {
	var result telegram.BotDescription
	if err := b.call(ctx, "getMyDescription", params, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// SetMyShortDescription changes the bot's short description.
func (b *Bot) SetMyShortDescription(ctx context.Context, params SetMyShortDescriptionParams) (bool, error) {
	var result bool
	if err := b.call(ctx, "setMyShortDescription", params, &result); err != nil {
		return false, err
	}
	return result, nil
}

// GetMyShortDescription gets the bot's current short description for a language.
func (b *Bot) GetMyShortDescription(ctx context.Context, params GetMyShortDescriptionParams) (*telegram.BotShortDescription, error) {
	var result telegram.BotShortDescription
	if err := b.call(ctx, "getMyShortDescription", params, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetMyDefaultAdministratorRights gets the default administrator rights requested by the bot.
func (b *Bot) GetMyDefaultAdministratorRights(ctx context.Context, params GetMyDefaultAdministratorRightsParams) (*telegram.ChatAdministratorRights, error) {
	var result telegram.ChatAdministratorRights
	if err := b.call(ctx, "getMyDefaultAdministratorRights", params, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// SetMyProfilePhoto uploads and sets a new bot profile photo.
func (b *Bot) SetMyProfilePhoto(ctx context.Context, params SetMyProfilePhotoParams) (bool, error) {
	if err := params.validate(); err != nil {
		return false, err
	}

	fields, files, err := params.multipart()
	if err != nil {
		return false, err
	}

	var result bool
	if err := b.callMultipart(ctx, "setMyProfilePhoto", fields, files, &result); err != nil {
		return false, err
	}
	return result, nil
}

// RemoveMyProfilePhoto removes the bot's profile photo.
func (b *Bot) RemoveMyProfilePhoto(ctx context.Context, params RemoveMyProfilePhotoParams) (bool, error) {
	var result bool
	if err := b.call(ctx, "removeMyProfilePhoto", params, &result); err != nil {
		return false, err
	}
	return result, nil
}

func (params SetMyProfilePhotoParams) validate() error {
	_, _, err := params.profilePhotoPayloadAndUpload()
	return err
}

func (params SetMyProfilePhotoParams) multipart() (map[string]string, map[string]UploadFile, error) {
	payload, files, err := params.profilePhotoPayloadAndUpload()
	if err != nil {
		return nil, nil, err
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, nil, err
	}
	return map[string]string{"photo": string(body)}, files, nil
}

func (params SetMyProfilePhotoParams) profilePhotoPayloadAndUpload() (map[string]any, map[string]UploadFile, error) {
	return profilePhotoPayloadAndUpload(params.Photo)
}

func profilePhotoPayloadAndUpload(photo InputProfilePhoto) (map[string]any, map[string]UploadFile, error) {
	if photo == nil {
		return nil, nil, stderrors.New("photo is required")
	}

	switch photo := photo.(type) {
	case InputProfilePhotoStatic:
		return staticProfilePhotoPayload(photo)
	case *InputProfilePhotoStatic:
		if photo == nil {
			return nil, nil, stderrors.New("photo is required")
		}
		return staticProfilePhotoPayload(*photo)
	case InputProfilePhotoAnimated:
		return animatedProfilePhotoPayload(photo)
	case *InputProfilePhotoAnimated:
		if photo == nil {
			return nil, nil, stderrors.New("photo is required")
		}
		return animatedProfilePhotoPayload(*photo)
	default:
		return nil, nil, stderrors.New("profile photo type is unsupported")
	}
}

func staticProfilePhotoPayload(photo InputProfilePhotoStatic) (map[string]any, map[string]UploadFile, error) {
	if strings.TrimSpace(photo.Type) != "" && photo.Type != "static" {
		return nil, nil, stderrors.New("profile photo type must be static")
	}
	if err := validateProfilePhotoUpload(photo.Photo, "photo"); err != nil {
		return nil, nil, err
	}
	payload := map[string]any{"type": "static", "photo": "attach://photo"}
	files := map[string]UploadFile{"photo": photo.Photo.upload}
	return payload, files, nil
}

func animatedProfilePhotoPayload(photo InputProfilePhotoAnimated) (map[string]any, map[string]UploadFile, error) {
	if strings.TrimSpace(photo.Type) != "" && photo.Type != "animated" {
		return nil, nil, stderrors.New("profile photo type must be animated")
	}
	if err := validateProfilePhotoUpload(photo.Animation, "animation"); err != nil {
		return nil, nil, err
	}
	if photo.MainFrameTimestamp < 0 {
		return nil, nil, stderrors.New("main_frame_timestamp must not be negative")
	}
	payload := map[string]any{"type": "animated", "animation": "attach://animation"}
	if photo.MainFrameTimestamp > 0 {
		payload["main_frame_timestamp"] = photo.MainFrameTimestamp
	}
	files := map[string]UploadFile{"animation": photo.Animation.upload}
	return payload, files, nil
}

func validateProfilePhotoUpload(ref FileRef, field string) error {
	if ref.kind != fileRefUpload {
		return stderrors.New(field + " must be uploaded with FileUpload")
	}
	return ref.validate(field)
}
