package telegram

import (
	"encoding/json"
	"fmt"
)

// ChatFullInfo contains full information about a chat returned by getChat.
type ChatFullInfo struct {
	ID                                 int64                 `json:"id"`
	Type                               string                `json:"type"`
	Title                              string                `json:"title,omitempty"`
	Username                           string                `json:"username,omitempty"`
	FirstName                          string                `json:"first_name,omitempty"`
	LastName                           string                `json:"last_name,omitempty"`
	IsForum                            bool                  `json:"is_forum,omitempty"`
	IsDirectMessages                   bool                  `json:"is_direct_messages,omitempty"`
	AccentColorID                      int                   `json:"accent_color_id"`
	MaxReactionCount                   int                   `json:"max_reaction_count"`
	Photo                              *ChatPhoto            `json:"photo,omitempty"`
	ActiveUsernames                    []string              `json:"active_usernames,omitempty"`
	Birthdate                          *Birthdate            `json:"birthdate,omitempty"`
	BusinessIntro                      *BusinessIntro        `json:"business_intro,omitempty"`
	BusinessLocation                   *BusinessLocation     `json:"business_location,omitempty"`
	BusinessOpeningHours               *BusinessOpeningHours `json:"business_opening_hours,omitempty"`
	PersonalChat                       *Chat                 `json:"personal_chat,omitempty"`
	ParentChat                         *Chat                 `json:"parent_chat,omitempty"`
	AvailableReactions                 []ReactionType        `json:"available_reactions,omitempty"`
	BackgroundCustomEmojiID            string                `json:"background_custom_emoji_id,omitempty"`
	ProfileAccentColorID               int                   `json:"profile_accent_color_id,omitempty"`
	ProfileBackgroundCustomEmojiID     string                `json:"profile_background_custom_emoji_id,omitempty"`
	EmojiStatusCustomEmojiID           string                `json:"emoji_status_custom_emoji_id,omitempty"`
	EmojiStatusExpirationDate          int64                 `json:"emoji_status_expiration_date,omitempty"`
	Bio                                string                `json:"bio,omitempty"`
	HasPrivateForwards                 bool                  `json:"has_private_forwards,omitempty"`
	HasRestrictedVoiceAndVideoMessages bool                  `json:"has_restricted_voice_and_video_messages,omitempty"`
	JoinToSendMessages                 bool                  `json:"join_to_send_messages,omitempty"`
	JoinByRequest                      bool                  `json:"join_by_request,omitempty"`
	Description                        string                `json:"description,omitempty"`
	InviteLink                         string                `json:"invite_link,omitempty"`
	PinnedMessage                      *Message              `json:"pinned_message,omitempty"`
	Permissions                        *ChatPermissions      `json:"permissions,omitempty"`
	AcceptedGiftTypes                  AcceptedGiftTypes     `json:"accepted_gift_types"`
	CanSendPaidMedia                   bool                  `json:"can_send_paid_media,omitempty"`
	SlowModeDelay                      int                   `json:"slow_mode_delay,omitempty"`
	UnrestrictBoostCount               int                   `json:"unrestrict_boost_count,omitempty"`
	MessageAutoDeleteTime              int                   `json:"message_auto_delete_time,omitempty"`
	HasAggressiveAntiSpamEnabled       bool                  `json:"has_aggressive_anti_spam_enabled,omitempty"`
	HasHiddenMembers                   bool                  `json:"has_hidden_members,omitempty"`
	HasProtectedContent                bool                  `json:"has_protected_content,omitempty"`
	HasVisibleHistory                  bool                  `json:"has_visible_history,omitempty"`
	StickerSetName                     string                `json:"sticker_set_name,omitempty"`
	CanSetStickerSet                   bool                  `json:"can_set_sticker_set,omitempty"`
	CustomEmojiStickerSetName          string                `json:"custom_emoji_sticker_set_name,omitempty"`
	LinkedChatID                       int64                 `json:"linked_chat_id,omitempty"`
	Location                           *ChatLocation         `json:"location,omitempty"`
	Rating                             *UserRating           `json:"rating,omitempty"`
	FirstProfileAudio                  *Audio                `json:"first_profile_audio,omitempty"`
	UniqueGiftColors                   *UniqueGiftColors     `json:"unique_gift_colors,omitempty"`
	PaidMessageStarCount               int                   `json:"paid_message_star_count,omitempty"`
}

// ChatPhoto represents a chat photo.
type ChatPhoto struct {
	SmallFileID       string `json:"small_file_id"`
	SmallFileUniqueID string `json:"small_file_unique_id"`
	BigFileID         string `json:"big_file_id"`
	BigFileUniqueID   string `json:"big_file_unique_id"`
}

// Birthdate describes the birthdate of a user.
type Birthdate struct {
	Day   int `json:"day"`
	Month int `json:"month"`
	Year  int `json:"year,omitempty"`
}

// BusinessIntro contains start page settings of a Telegram Business account.
type BusinessIntro struct {
	Title   string   `json:"title,omitempty"`
	Message string   `json:"message,omitempty"`
	Sticker *Sticker `json:"sticker,omitempty"`
}

// BusinessLocation contains the location of a Telegram Business account.
type BusinessLocation struct {
	Address  string    `json:"address"`
	Location *Location `json:"location,omitempty"`
}

// BusinessOpeningHours describes the opening hours of a business.
type BusinessOpeningHours struct {
	TimeZoneName string                         `json:"time_zone_name"`
	OpeningHours []BusinessOpeningHoursInterval `json:"opening_hours"`
}

// BusinessOpeningHoursInterval describes one business opening interval.
type BusinessOpeningHoursInterval struct {
	OpeningMinute int `json:"opening_minute"`
	ClosingMinute int `json:"closing_minute"`
}

// ChatLocation represents a location to which a chat is connected.
type ChatLocation struct {
	Location Location `json:"location"`
	Address  string   `json:"address"`
}

// UserRating describes a user's Telegram Stars purchase reliability rating.
type UserRating struct {
	Level              int `json:"level"`
	Rating             int `json:"rating"`
	CurrentLevelRating int `json:"current_level_rating"`
	NextLevelRating    int `json:"next_level_rating,omitempty"`
}

// UnmarshalJSON decodes ChatFullInfo with polymorphic available reactions.
func (info *ChatFullInfo) UnmarshalJSON(data []byte) error {
	type chatFullInfoAlias ChatFullInfo
	payload := struct {
		AvailableReactions []json.RawMessage `json:"available_reactions"`
		*chatFullInfoAlias
	}{chatFullInfoAlias: (*chatFullInfoAlias)(info)}
	if err := json.Unmarshal(data, &payload); err != nil {
		return err
	}
	if len(payload.AvailableReactions) > 0 {
		reactions, err := unmarshalReactionTypes(payload.AvailableReactions)
		if err != nil {
			return fmt.Errorf("available_reactions is invalid: %w", err)
		}
		info.AvailableReactions = reactions
	}
	return nil
}
