package telegram

import (
	"encoding/json"
	stderrors "errors"
	"fmt"
	"strings"
)

const (
	// StoryAreaTypeLocationType identifies StoryAreaTypeLocation payloads.
	StoryAreaTypeLocationType = "location"
	// StoryAreaTypeSuggestedReactionType identifies StoryAreaTypeSuggestedReaction payloads.
	StoryAreaTypeSuggestedReactionType = "suggested_reaction"
	// StoryAreaTypeLinkType identifies StoryAreaTypeLink payloads.
	StoryAreaTypeLinkType = "link"
	// StoryAreaTypeWeatherType identifies StoryAreaTypeWeather payloads.
	StoryAreaTypeWeatherType = "weather"
	// StoryAreaTypeUniqueGiftType identifies StoryAreaTypeUniqueGift payloads.
	StoryAreaTypeUniqueGiftType = "unique_gift"
)

// AcceptedGiftTypes describes gift types accepted by a user, chat, or business account.
type AcceptedGiftTypes struct {
	UnlimitedGifts      bool `json:"unlimited_gifts"`
	LimitedGifts        bool `json:"limited_gifts"`
	UniqueGifts         bool `json:"unique_gifts"`
	PremiumSubscription bool `json:"premium_subscription"`
	GiftsFromChannels   bool `json:"gifts_from_channels"`
}

// Story represents a Telegram story.
type Story struct {
	Chat Chat  `json:"chat"`
	ID   int64 `json:"id"`
}

// StoryArea describes a clickable area on a story media.
type StoryArea struct {
	Position StoryAreaPosition `json:"position"`
	Type     StoryAreaType     `json:"type"`
}

// StoryAreaPosition describes the position of a clickable area within a story.
type StoryAreaPosition struct {
	XPercentage            float64 `json:"x_percentage"`
	YPercentage            float64 `json:"y_percentage"`
	WidthPercentage        float64 `json:"width_percentage"`
	HeightPercentage       float64 `json:"height_percentage"`
	RotationAngle          float64 `json:"rotation_angle"`
	CornerRadiusPercentage float64 `json:"corner_radius_percentage"`
}

// StoryAreaType marks Telegram story area type objects.
type StoryAreaType interface {
	storyAreaType()
}

// LocationAddress describes the physical address of a location.
type LocationAddress struct {
	CountryCode string `json:"country_code"`
	State       string `json:"state,omitempty"`
	City        string `json:"city,omitempty"`
	Street      string `json:"street,omitempty"`
}

// StoryAreaTypeLocation describes a story area pointing to a location.
type StoryAreaTypeLocation struct {
	Type      string           `json:"type"`
	Latitude  float64          `json:"latitude"`
	Longitude float64          `json:"longitude"`
	Address   *LocationAddress `json:"address,omitempty"`
}

// StoryAreaTypeSuggestedReaction describes a story area pointing to a suggested reaction.
type StoryAreaTypeSuggestedReaction struct {
	Type         string       `json:"type"`
	ReactionType ReactionType `json:"reaction_type"`
	IsDark       bool         `json:"is_dark,omitempty"`
	IsFlipped    bool         `json:"is_flipped,omitempty"`
}

// StoryAreaTypeLink describes a story area pointing to a link.
type StoryAreaTypeLink struct {
	Type string `json:"type"`
	URL  string `json:"url"`
}

// StoryAreaTypeWeather describes a story area containing weather information.
type StoryAreaTypeWeather struct {
	Type            string  `json:"type"`
	Temperature     float64 `json:"temperature"`
	Emoji           string  `json:"emoji"`
	BackgroundColor int     `json:"background_color"`
}

// StoryAreaTypeUniqueGift describes a story area pointing to a unique gift.
type StoryAreaTypeUniqueGift struct {
	Type string `json:"type"`
	Name string `json:"name"`
}

func (StoryAreaTypeLocation) storyAreaType()          {}
func (StoryAreaTypeSuggestedReaction) storyAreaType() {}
func (StoryAreaTypeLink) storyAreaType()              {}
func (StoryAreaTypeWeather) storyAreaType()           {}
func (StoryAreaTypeUniqueGift) storyAreaType()        {}

// NewStoryAreaTypeLocation creates a location story area type.
func NewStoryAreaTypeLocation(latitude, longitude float64) StoryAreaTypeLocation {
	return StoryAreaTypeLocation{Type: StoryAreaTypeLocationType, Latitude: latitude, Longitude: longitude}
}

// NewStoryAreaTypeSuggestedReaction creates a suggested reaction story area type.
func NewStoryAreaTypeSuggestedReaction(reaction ReactionType) StoryAreaTypeSuggestedReaction {
	return StoryAreaTypeSuggestedReaction{Type: StoryAreaTypeSuggestedReactionType, ReactionType: reaction}
}

// NewStoryAreaTypeLink creates a link story area type.
func NewStoryAreaTypeLink(url string) StoryAreaTypeLink {
	return StoryAreaTypeLink{Type: StoryAreaTypeLinkType, URL: url}
}

// NewStoryAreaTypeWeather creates a weather story area type.
func NewStoryAreaTypeWeather(temperature float64, emoji string, backgroundColor int) StoryAreaTypeWeather {
	return StoryAreaTypeWeather{Type: StoryAreaTypeWeatherType, Temperature: temperature, Emoji: emoji, BackgroundColor: backgroundColor}
}

// NewStoryAreaTypeUniqueGift creates a unique gift story area type.
func NewStoryAreaTypeUniqueGift(name string) StoryAreaTypeUniqueGift {
	return StoryAreaTypeUniqueGift{Type: StoryAreaTypeUniqueGiftType, Name: name}
}

// MarshalJSON encodes StoryAreaTypeLocation with the required Telegram type field.
func (s StoryAreaTypeLocation) MarshalJSON() ([]byte, error) {
	s.Type = StoryAreaTypeLocationType
	type storyArea StoryAreaTypeLocation
	return json.Marshal(storyArea(s))
}

// MarshalJSON encodes StoryAreaTypeSuggestedReaction with the required Telegram type field.
func (s StoryAreaTypeSuggestedReaction) MarshalJSON() ([]byte, error) {
	s.Type = StoryAreaTypeSuggestedReactionType
	type storyArea StoryAreaTypeSuggestedReaction
	return json.Marshal(storyArea(s))
}

// MarshalJSON encodes StoryAreaTypeLink with the required Telegram type field.
func (s StoryAreaTypeLink) MarshalJSON() ([]byte, error) {
	s.Type = StoryAreaTypeLinkType
	type storyArea StoryAreaTypeLink
	return json.Marshal(storyArea(s))
}

// MarshalJSON encodes StoryAreaTypeWeather with the required Telegram type field.
func (s StoryAreaTypeWeather) MarshalJSON() ([]byte, error) {
	s.Type = StoryAreaTypeWeatherType
	type storyArea StoryAreaTypeWeather
	return json.Marshal(storyArea(s))
}

// MarshalJSON encodes StoryAreaTypeUniqueGift with the required Telegram type field.
func (s StoryAreaTypeUniqueGift) MarshalJSON() ([]byte, error) {
	s.Type = StoryAreaTypeUniqueGiftType
	type storyArea StoryAreaTypeUniqueGift
	return json.Marshal(storyArea(s))
}

// ValidateAcceptedGiftTypes checks gift type settings before sending them to Telegram.
func ValidateAcceptedGiftTypes(types AcceptedGiftTypes) error {
	if !types.UnlimitedGifts && !types.LimitedGifts && !types.UniqueGifts && !types.PremiumSubscription && !types.GiftsFromChannels {
		return stderrors.New("accepted_gift_types must enable at least one gift type")
	}
	return nil
}

// ValidateStoryAreas checks story areas before sending them to Telegram.
func ValidateStoryAreas(areas []StoryArea) error {
	for index, area := range areas {
		if err := ValidateStoryArea(area); err != nil {
			return fmt.Errorf("areas[%d]: %w", index, err)
		}
	}
	return nil
}

// ValidateStoryArea checks a story area before sending it to Telegram.
func ValidateStoryArea(area StoryArea) error {
	if area.Position.WidthPercentage < 0 || area.Position.HeightPercentage < 0 || area.Position.CornerRadiusPercentage < 0 {
		return stderrors.New("story area position dimensions must not be negative")
	}
	if area.Type == nil || isNilInterfaceValue(area.Type) {
		return stderrors.New("story area type is required")
	}
	return ValidateStoryAreaType(area.Type)
}

// ValidateStoryAreaType checks a story area type before sending it to Telegram.
func ValidateStoryAreaType(areaType StoryAreaType) error {
	if areaType == nil || isNilInterfaceValue(areaType) {
		return stderrors.New("story area type is required")
	}
	switch value := areaType.(type) {
	case StoryAreaTypeLocation:
		return validateStoryAreaTypeLocation(value)
	case *StoryAreaTypeLocation:
		if value == nil {
			return stderrors.New("story area type is required")
		}
		return validateStoryAreaTypeLocation(*value)
	case StoryAreaTypeSuggestedReaction:
		return validateStoryAreaTypeSuggestedReaction(value)
	case *StoryAreaTypeSuggestedReaction:
		if value == nil {
			return stderrors.New("story area type is required")
		}
		return validateStoryAreaTypeSuggestedReaction(*value)
	case StoryAreaTypeLink:
		return validateStoryAreaTypeLink(value)
	case *StoryAreaTypeLink:
		if value == nil {
			return stderrors.New("story area type is required")
		}
		return validateStoryAreaTypeLink(*value)
	case StoryAreaTypeWeather:
		return validateStoryAreaTypeWeather(value)
	case *StoryAreaTypeWeather:
		if value == nil {
			return stderrors.New("story area type is required")
		}
		return validateStoryAreaTypeWeather(*value)
	case StoryAreaTypeUniqueGift:
		return validateStoryAreaTypeUniqueGift(value)
	case *StoryAreaTypeUniqueGift:
		if value == nil {
			return stderrors.New("story area type is required")
		}
		return validateStoryAreaTypeUniqueGift(*value)
	default:
		return stderrors.New("unsupported story area type")
	}
}

func validateStoryAreaTypeLocation(areaType StoryAreaTypeLocation) error {
	if areaType.Latitude < -90 || areaType.Latitude > 90 {
		return stderrors.New("story area latitude must be between -90 and 90")
	}
	if areaType.Longitude < -180 || areaType.Longitude > 180 {
		return stderrors.New("story area longitude must be between -180 and 180")
	}
	if areaType.Address != nil && strings.TrimSpace(areaType.Address.CountryCode) == "" {
		return stderrors.New("story area location address country_code is required")
	}
	return nil
}

func validateStoryAreaTypeSuggestedReaction(areaType StoryAreaTypeSuggestedReaction) error {
	return ValidateReactionType(areaType.ReactionType)
}

func validateStoryAreaTypeLink(areaType StoryAreaTypeLink) error {
	if strings.TrimSpace(areaType.URL) == "" {
		return stderrors.New("story area url is required")
	}
	return nil
}

func validateStoryAreaTypeWeather(areaType StoryAreaTypeWeather) error {
	if strings.TrimSpace(areaType.Emoji) == "" {
		return stderrors.New("story area weather emoji is required")
	}
	return nil
}

func validateStoryAreaTypeUniqueGift(areaType StoryAreaTypeUniqueGift) error {
	if strings.TrimSpace(areaType.Name) == "" {
		return stderrors.New("story area unique gift name is required")
	}
	return nil
}
