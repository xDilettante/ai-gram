package bot

import (
	"context"
	stderrors "errors"

	"github.com/xDilettante/ai-gram/telegram"
)

// GetUserProfilePhotosParams contains supported parameters for getUserProfilePhotos.
type GetUserProfilePhotosParams struct {
	UserID int64 `json:"user_id"`
	Offset int   `json:"offset,omitempty"`
	Limit  int   `json:"limit,omitempty"`
}

// GetUserProfileAudiosParams contains supported parameters for getUserProfileAudios.
type GetUserProfileAudiosParams struct {
	UserID int64 `json:"user_id"`
	Offset int   `json:"offset,omitempty"`
	Limit  int   `json:"limit,omitempty"`
}

// LogOut logs the bot out from the cloud Bot API server.
func (b *Bot) LogOut(ctx context.Context) (bool, error) {
	var result bool
	if err := b.call(ctx, "logOut", nil, &result); err != nil {
		return false, err
	}
	return result, nil
}

// Close closes the bot instance on a local Bot API server.
func (b *Bot) Close(ctx context.Context) (bool, error) {
	var result bool
	if err := b.call(ctx, "close", nil, &result); err != nil {
		return false, err
	}
	return result, nil
}

// GetUserProfilePhotos gets a list of profile pictures for a user.
func (b *Bot) GetUserProfilePhotos(ctx context.Context, params GetUserProfilePhotosParams) (*telegram.UserProfilePhotos, error) {
	if err := params.validate(); err != nil {
		return nil, err
	}

	var result telegram.UserProfilePhotos
	if err := b.call(ctx, "getUserProfilePhotos", params, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetUserProfileAudios gets a list of profile audios for a user.
func (b *Bot) GetUserProfileAudios(ctx context.Context, params GetUserProfileAudiosParams) (*telegram.UserProfileAudios, error) {
	if err := params.validate(); err != nil {
		return nil, err
	}

	var result telegram.UserProfileAudios
	if err := b.call(ctx, "getUserProfileAudios", params, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetForumTopicIconStickers gets custom emoji stickers available as forum topic icons.
func (b *Bot) GetForumTopicIconStickers(ctx context.Context) ([]telegram.Sticker, error) {
	var result []telegram.Sticker
	if err := b.call(ctx, "getForumTopicIconStickers", nil, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (params GetUserProfilePhotosParams) validate() error {
	return validateUserProfileMediaParams(params.UserID, params.Offset, params.Limit)
}

func (params GetUserProfileAudiosParams) validate() error {
	return validateUserProfileMediaParams(params.UserID, params.Offset, params.Limit)
}

func validateUserProfileMediaParams(userID int64, offset int, limit int) error {
	if userID <= 0 {
		return stderrors.New("user_id must be greater than zero")
	}
	if offset < 0 {
		return stderrors.New("offset must not be negative")
	}
	if limit < 0 {
		return stderrors.New("limit must not be negative")
	}
	return nil
}
