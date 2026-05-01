package telegram

import (
	"bytes"
	"encoding/json"
	stderrors "errors"
)

const (
	backgroundTypeFillType             = "fill"
	backgroundTypeWallpaperType        = "wallpaper"
	backgroundTypePatternType          = "pattern"
	backgroundTypeChatThemeType        = "chat_theme"
	backgroundFillSolidType            = "solid"
	backgroundFillGradientType         = "gradient"
	backgroundFillFreeformGradientType = "freeform_gradient"
)

// UsersShared contains users shared through a KeyboardButtonRequestUsers button.
type UsersShared struct {
	RequestID int          `json:"request_id"`
	Users     []SharedUser `json:"users"`
}

// SharedUser contains information about one user shared with the bot.
type SharedUser struct {
	UserID    int64       `json:"user_id"`
	FirstName string      `json:"first_name,omitempty"`
	LastName  string      `json:"last_name,omitempty"`
	Username  string      `json:"username,omitempty"`
	Photo     []PhotoSize `json:"photo,omitempty"`
}

// ChatShared contains a chat shared through a KeyboardButtonRequestChat button.
type ChatShared struct {
	RequestID int         `json:"request_id"`
	ChatID    int64       `json:"chat_id"`
	Title     string      `json:"title,omitempty"`
	Username  string      `json:"username,omitempty"`
	Photo     []PhotoSize `json:"photo,omitempty"`
}

// ChatOwnerLeft describes a service message about the chat owner leaving.
type ChatOwnerLeft struct {
	NewOwner *User `json:"new_owner,omitempty"`
}

// ChatOwnerChanged describes a service message about the chat owner changing.
type ChatOwnerChanged struct {
	NewOwner User `json:"new_owner"`
}

// ProximityAlertTriggered describes a proximity alert service message.
type ProximityAlertTriggered struct {
	Traveler User `json:"traveler"`
	Watcher  User `json:"watcher"`
	Distance int  `json:"distance"`
}

// MessageAutoDeleteTimerChanged describes changed chat auto-delete settings.
type MessageAutoDeleteTimerChanged struct {
	MessageAutoDeleteTime int `json:"message_auto_delete_time"`
}

// ChatBoostAdded describes a service message about a user boosting the chat.
type ChatBoostAdded struct {
	BoostCount int `json:"boost_count"`
}

// ChatBackground represents a chat background service message payload.
type ChatBackground struct {
	Type BackgroundType `json:"type"`
}

// BackgroundType marks Telegram chat background type objects.
type BackgroundType interface {
	backgroundType()
}

// BackgroundTypeFill describes an automatically filled background.
type BackgroundTypeFill struct {
	Type             string         `json:"type"`
	Fill             BackgroundFill `json:"fill"`
	DarkThemeDimming int            `json:"dark_theme_dimming"`
}

// BackgroundTypeWallpaper describes a wallpaper background.
type BackgroundTypeWallpaper struct {
	Type             string   `json:"type"`
	Document         Document `json:"document"`
	DarkThemeDimming int      `json:"dark_theme_dimming"`
	IsBlurred        bool     `json:"is_blurred,omitempty"`
	IsMoving         bool     `json:"is_moving,omitempty"`
}

// BackgroundTypePattern describes a pattern background.
type BackgroundTypePattern struct {
	Type       string         `json:"type"`
	Document   Document       `json:"document"`
	Fill       BackgroundFill `json:"fill"`
	Intensity  int            `json:"intensity"`
	IsInverted bool           `json:"is_inverted,omitempty"`
	IsMoving   bool           `json:"is_moving,omitempty"`
}

// BackgroundTypeChatTheme describes a built-in chat theme background.
type BackgroundTypeChatTheme struct {
	Type      string `json:"type"`
	ThemeName string `json:"theme_name"`
}

func (BackgroundTypeFill) backgroundType()      {}
func (BackgroundTypeWallpaper) backgroundType() {}
func (BackgroundTypePattern) backgroundType()   {}
func (BackgroundTypeChatTheme) backgroundType() {}

// BackgroundFill marks Telegram background fill objects.
type BackgroundFill interface {
	backgroundFill()
}

// BackgroundFillSolid describes a solid background fill.
type BackgroundFillSolid struct {
	Type  string `json:"type"`
	Color int    `json:"color"`
}

// BackgroundFillGradient describes a two-color gradient fill.
type BackgroundFillGradient struct {
	Type          string `json:"type"`
	TopColor      int    `json:"top_color"`
	BottomColor   int    `json:"bottom_color"`
	RotationAngle int    `json:"rotation_angle"`
}

// BackgroundFillFreeformGradient describes a freeform gradient fill.
type BackgroundFillFreeformGradient struct {
	Type   string `json:"type"`
	Colors []int  `json:"colors"`
}

func (BackgroundFillSolid) backgroundFill()            {}
func (BackgroundFillGradient) backgroundFill()         {}
func (BackgroundFillFreeformGradient) backgroundFill() {}

// GiveawayCreated describes a service message about a scheduled giveaway creation.
type GiveawayCreated struct {
	PrizeStarCount int `json:"prize_star_count,omitempty"`
}

// GiveawayCompleted describes a service message about a completed giveaway without public winners.
type GiveawayCompleted struct {
	WinnerCount         int      `json:"winner_count"`
	UnclaimedPrizeCount int      `json:"unclaimed_prize_count,omitempty"`
	GiveawayMessage     *Message `json:"giveaway_message,omitempty"`
	IsStarGiveaway      bool     `json:"is_star_giveaway,omitempty"`
}

// PaidMessagePriceChanged describes a service message about paid-message price changes.
type PaidMessagePriceChanged struct {
	PaidMessageStarCount int `json:"paid_message_star_count"`
}

// DirectMessagePriceChanged describes a service message about channel direct-message price changes.
type DirectMessagePriceChanged struct {
	AreDirectMessagesEnabled bool `json:"are_direct_messages_enabled"`
	DirectMessageStarCount   int  `json:"direct_message_star_count,omitempty"`
}

// VideoChatScheduled describes a video chat scheduled service message.
type VideoChatScheduled struct {
	StartDate int64 `json:"start_date"`
}

// VideoChatStarted describes a video chat started service message.
type VideoChatStarted struct{}

// VideoChatEnded describes a video chat ended service message.
type VideoChatEnded struct {
	Duration int `json:"duration"`
}

// VideoChatParticipantsInvited describes users invited to a video chat.
type VideoChatParticipantsInvited struct {
	Users []User `json:"users"`
}

// UnmarshalJSON decodes a ChatBackground with a polymorphic background type.
func (b *ChatBackground) UnmarshalJSON(data []byte) error {
	var payload struct {
		Type json.RawMessage `json:"type"`
	}
	if err := json.Unmarshal(data, &payload); err != nil {
		return err
	}
	if len(payload.Type) == 0 || bytes.Equal(payload.Type, []byte("null")) {
		return nil
	}
	backgroundType, err := UnmarshalBackgroundType(payload.Type)
	if err != nil {
		return err
	}
	b.Type = backgroundType
	return nil
}

// UnmarshalBackgroundType decodes a polymorphic Telegram BackgroundType object.
func UnmarshalBackgroundType(data []byte) (BackgroundType, error) {
	var meta struct {
		Type string `json:"type"`
	}
	if err := json.Unmarshal(data, &meta); err != nil {
		return nil, err
	}
	switch meta.Type {
	case backgroundTypeFillType:
		var value BackgroundTypeFill
		if err := json.Unmarshal(data, &value); err != nil {
			return nil, err
		}
		value.Type = backgroundTypeFillType
		return value, nil
	case backgroundTypeWallpaperType:
		var value BackgroundTypeWallpaper
		if err := json.Unmarshal(data, &value); err != nil {
			return nil, err
		}
		value.Type = backgroundTypeWallpaperType
		return value, nil
	case backgroundTypePatternType:
		var value BackgroundTypePattern
		if err := json.Unmarshal(data, &value); err != nil {
			return nil, err
		}
		value.Type = backgroundTypePatternType
		return value, nil
	case backgroundTypeChatThemeType:
		var value BackgroundTypeChatTheme
		if err := json.Unmarshal(data, &value); err != nil {
			return nil, err
		}
		value.Type = backgroundTypeChatThemeType
		return value, nil
	default:
		return nil, stderrors.New("unsupported background type")
	}
}

// UnmarshalJSON decodes a BackgroundTypeFill with a polymorphic fill object.
func (b *BackgroundTypeFill) UnmarshalJSON(data []byte) error {
	type backgroundTypeFillAlias BackgroundTypeFill
	payload := struct {
		Fill json.RawMessage `json:"fill"`
		*backgroundTypeFillAlias
	}{backgroundTypeFillAlias: (*backgroundTypeFillAlias)(b)}
	if err := json.Unmarshal(data, &payload); err != nil {
		return err
	}
	if len(payload.Fill) > 0 && !bytes.Equal(payload.Fill, []byte("null")) {
		fill, err := UnmarshalBackgroundFill(payload.Fill)
		if err != nil {
			return err
		}
		b.Fill = fill
	}
	return nil
}

// UnmarshalJSON decodes a BackgroundTypePattern with a polymorphic fill object.
func (b *BackgroundTypePattern) UnmarshalJSON(data []byte) error {
	type backgroundTypePatternAlias BackgroundTypePattern
	payload := struct {
		Fill json.RawMessage `json:"fill"`
		*backgroundTypePatternAlias
	}{backgroundTypePatternAlias: (*backgroundTypePatternAlias)(b)}
	if err := json.Unmarshal(data, &payload); err != nil {
		return err
	}
	if len(payload.Fill) > 0 && !bytes.Equal(payload.Fill, []byte("null")) {
		fill, err := UnmarshalBackgroundFill(payload.Fill)
		if err != nil {
			return err
		}
		b.Fill = fill
	}
	return nil
}

// UnmarshalBackgroundFill decodes a polymorphic Telegram BackgroundFill object.
func UnmarshalBackgroundFill(data []byte) (BackgroundFill, error) {
	var meta struct {
		Type string `json:"type"`
	}
	if err := json.Unmarshal(data, &meta); err != nil {
		return nil, err
	}
	switch meta.Type {
	case backgroundFillSolidType:
		var value BackgroundFillSolid
		if err := json.Unmarshal(data, &value); err != nil {
			return nil, err
		}
		value.Type = backgroundFillSolidType
		return value, nil
	case backgroundFillGradientType:
		var value BackgroundFillGradient
		if err := json.Unmarshal(data, &value); err != nil {
			return nil, err
		}
		value.Type = backgroundFillGradientType
		return value, nil
	case backgroundFillFreeformGradientType:
		var value BackgroundFillFreeformGradient
		if err := json.Unmarshal(data, &value); err != nil {
			return nil, err
		}
		value.Type = backgroundFillFreeformGradientType
		return value, nil
	default:
		return nil, stderrors.New("unsupported background fill type")
	}
}
