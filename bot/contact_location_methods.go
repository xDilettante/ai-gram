package bot

import (
	"context"
	stderrors "errors"
	"math"

	"github.com/xDilettante/ai-gram/telegram"
)

// SendContactParams contains supported parameters for sendContact.
type SendContactParams struct {
	ChatID              ChatID                    `json:"chat_id"`
	MessageThreadID     int64                     `json:"message_thread_id,omitempty"`
	PhoneNumber         string                    `json:"phone_number"`
	FirstName           string                    `json:"first_name"`
	LastName            string                    `json:"last_name,omitempty"`
	VCard               string                    `json:"vcard,omitempty"`
	DisableNotification bool                      `json:"disable_notification,omitempty"`
	ProtectContent      bool                      `json:"protect_content,omitempty"`
	ReplyParameters     *telegram.ReplyParameters `json:"reply_parameters,omitempty"`
	ReplyMarkup         telegram.ReplyMarkup      `json:"reply_markup,omitempty"`
}

// SendLocationParams contains supported parameters for sendLocation.
type SendLocationParams struct {
	ChatID               ChatID                    `json:"chat_id"`
	MessageThreadID      int64                     `json:"message_thread_id,omitempty"`
	Latitude             float64                   `json:"latitude"`
	Longitude            float64                   `json:"longitude"`
	HorizontalAccuracy   float64                   `json:"horizontal_accuracy,omitempty"`
	LivePeriod           int                       `json:"live_period,omitempty"`
	Heading              int                       `json:"heading,omitempty"`
	ProximityAlertRadius int                       `json:"proximity_alert_radius,omitempty"`
	DisableNotification  bool                      `json:"disable_notification,omitempty"`
	ProtectContent       bool                      `json:"protect_content,omitempty"`
	ReplyParameters      *telegram.ReplyParameters `json:"reply_parameters,omitempty"`
	ReplyMarkup          telegram.ReplyMarkup      `json:"reply_markup,omitempty"`
}

// SendVenueParams contains supported parameters for sendVenue.
type SendVenueParams struct {
	ChatID              ChatID                    `json:"chat_id"`
	MessageThreadID     int64                     `json:"message_thread_id,omitempty"`
	Latitude            float64                   `json:"latitude"`
	Longitude           float64                   `json:"longitude"`
	Title               string                    `json:"title"`
	Address             string                    `json:"address"`
	FoursquareID        string                    `json:"foursquare_id,omitempty"`
	FoursquareType      string                    `json:"foursquare_type,omitempty"`
	GooglePlaceID       string                    `json:"google_place_id,omitempty"`
	GooglePlaceType     string                    `json:"google_place_type,omitempty"`
	DisableNotification bool                      `json:"disable_notification,omitempty"`
	ProtectContent      bool                      `json:"protect_content,omitempty"`
	ReplyParameters     *telegram.ReplyParameters `json:"reply_parameters,omitempty"`
	ReplyMarkup         telegram.ReplyMarkup      `json:"reply_markup,omitempty"`
}

// SendContact sends a phone contact.
func (b *Bot) SendContact(ctx context.Context, params SendContactParams) (*telegram.Message, error) {
	if err := params.validate(); err != nil {
		return nil, err
	}

	var message telegram.Message
	if err := b.call(ctx, "sendContact", params, &message); err != nil {
		return nil, err
	}

	return &message, nil
}

// SendLocation sends a point on the map.
func (b *Bot) SendLocation(ctx context.Context, params SendLocationParams) (*telegram.Message, error) {
	if err := params.validate(); err != nil {
		return nil, err
	}

	var message telegram.Message
	if err := b.call(ctx, "sendLocation", params, &message); err != nil {
		return nil, err
	}

	return &message, nil
}

// SendVenue sends information about a venue.
func (b *Bot) SendVenue(ctx context.Context, params SendVenueParams) (*telegram.Message, error) {
	if err := params.validate(); err != nil {
		return nil, err
	}

	var message telegram.Message
	if err := b.call(ctx, "sendVenue", params, &message); err != nil {
		return nil, err
	}

	return &message, nil
}

func (params SendContactParams) validate() error {
	if !params.ChatID.valid() {
		return stderrors.New("chat_id is required")
	}
	if err := validateMessageThreadID(params.MessageThreadID); err != nil {
		return err
	}
	if params.PhoneNumber == "" {
		return stderrors.New("phone_number is required")
	}
	if params.FirstName == "" {
		return stderrors.New("first_name is required")
	}
	if err := validateReplyParameters(params.ReplyParameters); err != nil {
		return err
	}
	if err := telegram.ValidateReplyMarkup(params.ReplyMarkup); err != nil {
		return err
	}

	return nil
}

func (params SendLocationParams) validate() error {
	if !params.ChatID.valid() {
		return stderrors.New("chat_id is required")
	}
	if err := validateMessageThreadID(params.MessageThreadID); err != nil {
		return err
	}
	if err := validateLatitude(params.Latitude); err != nil {
		return err
	}
	if err := validateLongitude(params.Longitude); err != nil {
		return err
	}
	if params.HorizontalAccuracy < 0 {
		return stderrors.New("horizontal_accuracy must not be negative")
	}
	if params.LivePeriod < 0 {
		return stderrors.New("live_period must not be negative")
	}
	if params.Heading < 0 {
		return stderrors.New("heading must not be negative")
	}
	if params.ProximityAlertRadius < 0 {
		return stderrors.New("proximity_alert_radius must not be negative")
	}
	if err := validateReplyParameters(params.ReplyParameters); err != nil {
		return err
	}
	if err := telegram.ValidateReplyMarkup(params.ReplyMarkup); err != nil {
		return err
	}

	return nil
}

func (params SendVenueParams) validate() error {
	if !params.ChatID.valid() {
		return stderrors.New("chat_id is required")
	}
	if err := validateMessageThreadID(params.MessageThreadID); err != nil {
		return err
	}
	if err := validateLatitude(params.Latitude); err != nil {
		return err
	}
	if err := validateLongitude(params.Longitude); err != nil {
		return err
	}
	if params.Title == "" {
		return stderrors.New("title is required")
	}
	if params.Address == "" {
		return stderrors.New("address is required")
	}
	if err := validateReplyParameters(params.ReplyParameters); err != nil {
		return err
	}
	if err := telegram.ValidateReplyMarkup(params.ReplyMarkup); err != nil {
		return err
	}

	return nil
}

func validateLatitude(latitude float64) error {
	if math.IsNaN(latitude) || math.IsInf(latitude, 0) || latitude < -90 || latitude > 90 {
		return stderrors.New("latitude must be between -90 and 90")
	}

	return nil
}

func validateLongitude(longitude float64) error {
	if math.IsNaN(longitude) || math.IsInf(longitude, 0) || longitude < -180 || longitude > 180 {
		return stderrors.New("longitude must be between -180 and 180")
	}

	return nil
}
