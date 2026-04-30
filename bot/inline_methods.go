package bot

import (
	"context"
	"encoding/json"
	stderrors "errors"
	"fmt"
	"net/url"
	"reflect"
	"strings"

	"github.com/xDilettante/ai-gram/telegram"
)

const (
	inlineQueryResultArticleType  = "article"
	inlineQueryResultLocationType = "location"
	inlineQueryResultVenueType    = "venue"
	inlineQueryResultContactType  = "contact"
	inlineQueryResultGameType     = "game"
)

// InputMessageContent marks Telegram input message content objects used by inline results.
type InputMessageContent interface {
	inputMessageContent()
}

// InputTextMessageContent describes text content to be sent when an inline result is chosen.
type InputTextMessageContent struct {
	MessageText        string                       `json:"message_text"`
	ParseMode          string                       `json:"parse_mode,omitempty"`
	Entities           []telegram.MessageEntity     `json:"entities,omitempty"`
	LinkPreviewOptions *telegram.LinkPreviewOptions `json:"link_preview_options,omitempty"`
}

func (InputTextMessageContent) inputMessageContent() {}

// InputText creates text content for an inline query result.
func InputText(message string) InputTextMessageContent {
	return InputTextMessageContent{MessageText: message}
}

// InputLocationMessageContent describes location content to be sent when an inline result is chosen.
type InputLocationMessageContent struct {
	Latitude             float64 `json:"latitude"`
	Longitude            float64 `json:"longitude"`
	HorizontalAccuracy   float64 `json:"horizontal_accuracy,omitempty"`
	LivePeriod           int     `json:"live_period,omitempty"`
	Heading              int     `json:"heading,omitempty"`
	ProximityAlertRadius int     `json:"proximity_alert_radius,omitempty"`
}

// InputVenueMessageContent describes venue content to be sent when an inline result is chosen.
type InputVenueMessageContent struct {
	Latitude        float64 `json:"latitude"`
	Longitude       float64 `json:"longitude"`
	Title           string  `json:"title"`
	Address         string  `json:"address"`
	FoursquareID    string  `json:"foursquare_id,omitempty"`
	FoursquareType  string  `json:"foursquare_type,omitempty"`
	GooglePlaceID   string  `json:"google_place_id,omitempty"`
	GooglePlaceType string  `json:"google_place_type,omitempty"`
}

// InputContactMessageContent describes contact content to be sent when an inline result is chosen.
type InputContactMessageContent struct {
	PhoneNumber string `json:"phone_number"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name,omitempty"`
	VCard       string `json:"vcard,omitempty"`
}

// InputInvoiceMessageContent describes invoice content to be sent when an inline result is chosen.
type InputInvoiceMessageContent struct {
	Title                     string                  `json:"title"`
	Description               string                  `json:"description"`
	Payload                   string                  `json:"payload"`
	ProviderToken             string                  `json:"provider_token,omitempty"`
	Currency                  string                  `json:"currency"`
	Prices                    []telegram.LabeledPrice `json:"prices"`
	MaxTipAmount              int64                   `json:"max_tip_amount,omitempty"`
	SuggestedTipAmounts       []int64                 `json:"suggested_tip_amounts,omitempty"`
	ProviderData              string                  `json:"provider_data,omitempty"`
	PhotoURL                  string                  `json:"photo_url,omitempty"`
	PhotoSize                 int64                   `json:"photo_size,omitempty"`
	PhotoWidth                int                     `json:"photo_width,omitempty"`
	PhotoHeight               int                     `json:"photo_height,omitempty"`
	NeedName                  bool                    `json:"need_name,omitempty"`
	NeedPhoneNumber           bool                    `json:"need_phone_number,omitempty"`
	NeedEmail                 bool                    `json:"need_email,omitempty"`
	NeedShippingAddress       bool                    `json:"need_shipping_address,omitempty"`
	SendPhoneNumberToProvider bool                    `json:"send_phone_number_to_provider,omitempty"`
	SendEmailToProvider       bool                    `json:"send_email_to_provider,omitempty"`
	IsFlexible                bool                    `json:"is_flexible,omitempty"`
}

func (InputLocationMessageContent) inputMessageContent() {}
func (InputVenueMessageContent) inputMessageContent()    {}
func (InputContactMessageContent) inputMessageContent()  {}
func (InputInvoiceMessageContent) inputMessageContent()  {}

// InputLocation creates location content for an inline query result.
func InputLocation(latitude float64, longitude float64) InputLocationMessageContent {
	return InputLocationMessageContent{Latitude: latitude, Longitude: longitude}
}

// InputVenue creates venue content for an inline query result.
func InputVenue(latitude float64, longitude float64, title string, address string) InputVenueMessageContent {
	return InputVenueMessageContent{Latitude: latitude, Longitude: longitude, Title: title, Address: address}
}

// InputContact creates contact content for an inline query result.
func InputContact(phoneNumber string, firstName string) InputContactMessageContent {
	return InputContactMessageContent{PhoneNumber: phoneNumber, FirstName: firstName}
}

// InlineQueryResult marks Telegram inline query result objects.
type InlineQueryResult interface {
	inlineQueryResult()
}

// InlineQueryResultArticle represents an article or web page inline query result.
type InlineQueryResultArticle struct {
	Type                string                         `json:"type"`
	ID                  string                         `json:"id"`
	Title               string                         `json:"title"`
	InputMessageContent InputMessageContent            `json:"input_message_content"`
	ReplyMarkup         *telegram.InlineKeyboardMarkup `json:"reply_markup,omitempty"`
	URL                 string                         `json:"url,omitempty"`
	HideURL             bool                           `json:"hide_url,omitempty"`
	Description         string                         `json:"description,omitempty"`
	ThumbnailURL        string                         `json:"thumbnail_url,omitempty"`
	ThumbnailWidth      int                            `json:"thumbnail_width,omitempty"`
	ThumbnailHeight     int                            `json:"thumbnail_height,omitempty"`
}

func (InlineQueryResultArticle) inlineQueryResult() {}

// InlineQueryResultLocation represents a location inline query result.
type InlineQueryResultLocation struct {
	Type                 string                         `json:"type"`
	ID                   string                         `json:"id"`
	Latitude             float64                        `json:"latitude"`
	Longitude            float64                        `json:"longitude"`
	Title                string                         `json:"title"`
	HorizontalAccuracy   float64                        `json:"horizontal_accuracy,omitempty"`
	LivePeriod           int                            `json:"live_period,omitempty"`
	Heading              int                            `json:"heading,omitempty"`
	ProximityAlertRadius int                            `json:"proximity_alert_radius,omitempty"`
	ReplyMarkup          *telegram.InlineKeyboardMarkup `json:"reply_markup,omitempty"`
	InputMessageContent  InputMessageContent            `json:"input_message_content,omitempty"`
	ThumbnailURL         string                         `json:"thumbnail_url,omitempty"`
	ThumbnailWidth       int                            `json:"thumbnail_width,omitempty"`
	ThumbnailHeight      int                            `json:"thumbnail_height,omitempty"`
}

// InlineQueryResultVenue represents a venue inline query result.
type InlineQueryResultVenue struct {
	Type                string                         `json:"type"`
	ID                  string                         `json:"id"`
	Latitude            float64                        `json:"latitude"`
	Longitude           float64                        `json:"longitude"`
	Title               string                         `json:"title"`
	Address             string                         `json:"address"`
	FoursquareID        string                         `json:"foursquare_id,omitempty"`
	FoursquareType      string                         `json:"foursquare_type,omitempty"`
	GooglePlaceID       string                         `json:"google_place_id,omitempty"`
	GooglePlaceType     string                         `json:"google_place_type,omitempty"`
	ReplyMarkup         *telegram.InlineKeyboardMarkup `json:"reply_markup,omitempty"`
	InputMessageContent InputMessageContent            `json:"input_message_content,omitempty"`
	ThumbnailURL        string                         `json:"thumbnail_url,omitempty"`
	ThumbnailWidth      int                            `json:"thumbnail_width,omitempty"`
	ThumbnailHeight     int                            `json:"thumbnail_height,omitempty"`
}

// InlineQueryResultContact represents a contact inline query result.
type InlineQueryResultContact struct {
	Type                string                         `json:"type"`
	ID                  string                         `json:"id"`
	PhoneNumber         string                         `json:"phone_number"`
	FirstName           string                         `json:"first_name"`
	LastName            string                         `json:"last_name,omitempty"`
	VCard               string                         `json:"vcard,omitempty"`
	ReplyMarkup         *telegram.InlineKeyboardMarkup `json:"reply_markup,omitempty"`
	InputMessageContent InputMessageContent            `json:"input_message_content,omitempty"`
	ThumbnailURL        string                         `json:"thumbnail_url,omitempty"`
	ThumbnailWidth      int                            `json:"thumbnail_width,omitempty"`
	ThumbnailHeight     int                            `json:"thumbnail_height,omitempty"`
}

// InlineQueryResultGame represents a game inline query result.
type InlineQueryResultGame struct {
	Type          string                         `json:"type"`
	ID            string                         `json:"id"`
	GameShortName string                         `json:"game_short_name"`
	ReplyMarkup   *telegram.InlineKeyboardMarkup `json:"reply_markup,omitempty"`
}

func (InlineQueryResultLocation) inlineQueryResult() {}
func (InlineQueryResultVenue) inlineQueryResult()    {}
func (InlineQueryResultContact) inlineQueryResult()  {}
func (InlineQueryResultGame) inlineQueryResult()     {}

// InlineArticle creates an article inline query result with text or another input message content.
func InlineArticle(id string, title string, content InputMessageContent) InlineQueryResultArticle {
	return InlineQueryResultArticle{Type: inlineQueryResultArticleType, ID: id, Title: title, InputMessageContent: content}
}

// MarshalJSON encodes InlineQueryResultArticle with the required Telegram type field.
func (result InlineQueryResultArticle) MarshalJSON() ([]byte, error) {
	result.Type = inlineQueryResultArticleType
	type article InlineQueryResultArticle
	return json.Marshal(article(result))
}

// InlineLocation creates a location inline query result.
func InlineLocation(id string, latitude float64, longitude float64, title string) InlineQueryResultLocation {
	return InlineQueryResultLocation{Type: inlineQueryResultLocationType, ID: id, Latitude: latitude, Longitude: longitude, Title: title}
}

// InlineVenue creates a venue inline query result.
func InlineVenue(id string, latitude float64, longitude float64, title string, address string) InlineQueryResultVenue {
	return InlineQueryResultVenue{Type: inlineQueryResultVenueType, ID: id, Latitude: latitude, Longitude: longitude, Title: title, Address: address}
}

// InlineContact creates a contact inline query result.
func InlineContact(id string, phoneNumber string, firstName string) InlineQueryResultContact {
	return InlineQueryResultContact{Type: inlineQueryResultContactType, ID: id, PhoneNumber: phoneNumber, FirstName: firstName}
}

// InlineGame creates a game inline query result.
func InlineGame(id string, gameShortName string) InlineQueryResultGame {
	return InlineQueryResultGame{Type: inlineQueryResultGameType, ID: id, GameShortName: gameShortName}
}

// MarshalJSON encodes InlineQueryResultLocation with the required Telegram type field.
func (result InlineQueryResultLocation) MarshalJSON() ([]byte, error) {
	result.Type = inlineQueryResultLocationType
	type location InlineQueryResultLocation
	return json.Marshal(location(result))
}

// MarshalJSON encodes InlineQueryResultVenue with the required Telegram type field.
func (result InlineQueryResultVenue) MarshalJSON() ([]byte, error) {
	result.Type = inlineQueryResultVenueType
	type venue InlineQueryResultVenue
	return json.Marshal(venue(result))
}

// MarshalJSON encodes InlineQueryResultContact with the required Telegram type field.
func (result InlineQueryResultContact) MarshalJSON() ([]byte, error) {
	result.Type = inlineQueryResultContactType
	type contact InlineQueryResultContact
	return json.Marshal(contact(result))
}

// MarshalJSON encodes InlineQueryResultGame with the required Telegram type field.
func (result InlineQueryResultGame) MarshalJSON() ([]byte, error) {
	result.Type = inlineQueryResultGameType
	type game InlineQueryResultGame
	return json.Marshal(game(result))
}

// AnswerInlineQueryParams contains supported parameters for answerInlineQuery.
type AnswerInlineQueryParams struct {
	InlineQueryID string                             `json:"inline_query_id"`
	Results       []InlineQueryResult                `json:"results"`
	CacheTime     int                                `json:"cache_time,omitempty"`
	IsPersonal    bool                               `json:"is_personal,omitempty"`
	NextOffset    string                             `json:"next_offset,omitempty"`
	Button        *telegram.InlineQueryResultsButton `json:"button,omitempty"`
}

// MarshalJSON encodes AnswerInlineQueryParams while preserving an empty results array.
func (params AnswerInlineQueryParams) MarshalJSON() ([]byte, error) {
	type payload AnswerInlineQueryParams
	encoded := payload(params)
	if encoded.Results == nil {
		encoded.Results = []InlineQueryResult{}
	}
	return json.Marshal(encoded)
}

// AnswerInlineQuery sends answers to an inline query.
func (b *Bot) AnswerInlineQuery(ctx context.Context, params AnswerInlineQueryParams) (bool, error) {
	if err := params.validate(); err != nil {
		return false, err
	}

	var result bool
	if err := b.call(ctx, "answerInlineQuery", params, &result); err != nil {
		return false, err
	}

	return result, nil
}

func (params AnswerInlineQueryParams) validate() error {
	if strings.TrimSpace(params.InlineQueryID) == "" {
		return stderrors.New("inline_query_id is required")
	}
	if len(params.Results) > 50 {
		return stderrors.New("inline query results must contain at most 50 items")
	}
	for index, result := range params.Results {
		if err := validateInlineQueryResult(result); err != nil {
			return fmt.Errorf("results[%d] is invalid: %w", index, err)
		}
	}
	if params.CacheTime < 0 {
		return stderrors.New("cache_time must not be negative")
	}
	if len([]byte(params.NextOffset)) > 64 {
		return stderrors.New("next_offset must be at most 64 bytes")
	}
	if err := telegram.ValidateInlineQueryResultsButton(params.Button); err != nil {
		return err
	}

	return nil
}

func validateInlineQueryResult(result InlineQueryResult) error {
	if result == nil || isNilBotInterfaceValue(result) {
		return stderrors.New("inline query result must not be nil")
	}

	switch value := result.(type) {
	case InlineQueryResultArticle:
		return validateInlineQueryResultArticle(value)
	case *InlineQueryResultArticle:
		return validateInlineQueryResultArticle(*value)
	case InlineQueryResultLocation:
		return validateInlineQueryResultLocation(value)
	case *InlineQueryResultLocation:
		return validateInlineQueryResultLocation(*value)
	case InlineQueryResultVenue:
		return validateInlineQueryResultVenue(value)
	case *InlineQueryResultVenue:
		return validateInlineQueryResultVenue(*value)
	case InlineQueryResultContact:
		return validateInlineQueryResultContact(value)
	case *InlineQueryResultContact:
		return validateInlineQueryResultContact(*value)
	case InlineQueryResultGame:
		return validateInlineQueryResultGame(value)
	case *InlineQueryResultGame:
		return validateInlineQueryResultGame(*value)
	case InlineQueryResultPhoto:
		return validateInlineQueryResultPhoto(value)
	case *InlineQueryResultPhoto:
		return validateInlineQueryResultPhoto(*value)
	case InlineQueryResultGif:
		return validateInlineQueryResultGif(value)
	case *InlineQueryResultGif:
		return validateInlineQueryResultGif(*value)
	case InlineQueryResultMpeg4Gif:
		return validateInlineQueryResultMpeg4Gif(value)
	case *InlineQueryResultMpeg4Gif:
		return validateInlineQueryResultMpeg4Gif(*value)
	case InlineQueryResultVideo:
		return validateInlineQueryResultVideo(value)
	case *InlineQueryResultVideo:
		return validateInlineQueryResultVideo(*value)
	case InlineQueryResultAudio:
		return validateInlineQueryResultAudio(value)
	case *InlineQueryResultAudio:
		return validateInlineQueryResultAudio(*value)
	case InlineQueryResultVoice:
		return validateInlineQueryResultVoice(value)
	case *InlineQueryResultVoice:
		return validateInlineQueryResultVoice(*value)
	case InlineQueryResultDocument:
		return validateInlineQueryResultDocument(value)
	case *InlineQueryResultDocument:
		return validateInlineQueryResultDocument(*value)
	case InlineQueryResultCachedPhoto:
		return validateInlineQueryResultCachedPhoto(value)
	case *InlineQueryResultCachedPhoto:
		return validateInlineQueryResultCachedPhoto(*value)
	case InlineQueryResultCachedGif:
		return validateInlineQueryResultCachedGif(value)
	case *InlineQueryResultCachedGif:
		return validateInlineQueryResultCachedGif(*value)
	case InlineQueryResultCachedMpeg4Gif:
		return validateInlineQueryResultCachedMpeg4Gif(value)
	case *InlineQueryResultCachedMpeg4Gif:
		return validateInlineQueryResultCachedMpeg4Gif(*value)
	case InlineQueryResultCachedSticker:
		return validateInlineQueryResultCachedSticker(value)
	case *InlineQueryResultCachedSticker:
		return validateInlineQueryResultCachedSticker(*value)
	case InlineQueryResultCachedDocument:
		return validateInlineQueryResultCachedDocument(value)
	case *InlineQueryResultCachedDocument:
		return validateInlineQueryResultCachedDocument(*value)
	case InlineQueryResultCachedVideo:
		return validateInlineQueryResultCachedVideo(value)
	case *InlineQueryResultCachedVideo:
		return validateInlineQueryResultCachedVideo(*value)
	case InlineQueryResultCachedVoice:
		return validateInlineQueryResultCachedVoice(value)
	case *InlineQueryResultCachedVoice:
		return validateInlineQueryResultCachedVoice(*value)
	case InlineQueryResultCachedAudio:
		return validateInlineQueryResultCachedAudio(value)
	case *InlineQueryResultCachedAudio:
		return validateInlineQueryResultCachedAudio(*value)
	default:
		return stderrors.New("unsupported inline query result")
	}
}

func validateInlineQueryResultArticle(result InlineQueryResultArticle) error {
	if strings.TrimSpace(result.ID) == "" {
		return stderrors.New("inline query result article id is required")
	}
	if strings.TrimSpace(result.Title) == "" {
		return stderrors.New("inline query result article title is required")
	}
	if err := validateInputMessageContent(result.InputMessageContent); err != nil {
		return err
	}
	if result.ReplyMarkup != nil {
		if err := telegram.ValidateReplyMarkup(*result.ReplyMarkup); err != nil {
			return err
		}
	}
	if result.URL != "" {
		if err := validateInlineHTTPURL(result.URL, "inline query result article URL"); err != nil {
			return err
		}
	}
	if result.ThumbnailURL != "" {
		if err := validateInlineHTTPURL(result.ThumbnailURL, "inline query result article thumbnail_url"); err != nil {
			return err
		}
	}
	if result.ThumbnailWidth < 0 {
		return stderrors.New("thumbnail_width must not be negative")
	}
	if result.ThumbnailHeight < 0 {
		return stderrors.New("thumbnail_height must not be negative")
	}

	return nil
}

func validateInlineQueryResultLocation(result InlineQueryResultLocation) error {
	if strings.TrimSpace(result.ID) == "" {
		return stderrors.New("inline query result location id is required")
	}
	if err := validateCoordinates(result.Latitude, result.Longitude); err != nil {
		return err
	}
	if strings.TrimSpace(result.Title) == "" {
		return stderrors.New("inline query result location title is required")
	}
	if err := validateLocationOptions(result.HorizontalAccuracy, result.LivePeriod, result.Heading, result.ProximityAlertRadius); err != nil {
		return err
	}
	if err := validateOptionalInlineResultFields(result.ReplyMarkup, result.InputMessageContent, result.ThumbnailURL, result.ThumbnailWidth, result.ThumbnailHeight); err != nil {
		return err
	}
	return nil
}

func validateInlineQueryResultVenue(result InlineQueryResultVenue) error {
	if strings.TrimSpace(result.ID) == "" {
		return stderrors.New("inline query result venue id is required")
	}
	if err := validateCoordinates(result.Latitude, result.Longitude); err != nil {
		return err
	}
	if strings.TrimSpace(result.Title) == "" {
		return stderrors.New("inline query result venue title is required")
	}
	if strings.TrimSpace(result.Address) == "" {
		return stderrors.New("inline query result venue address is required")
	}
	if err := validateOptionalInlineResultFields(result.ReplyMarkup, result.InputMessageContent, result.ThumbnailURL, result.ThumbnailWidth, result.ThumbnailHeight); err != nil {
		return err
	}
	return nil
}

func validateInlineQueryResultContact(result InlineQueryResultContact) error {
	if strings.TrimSpace(result.ID) == "" {
		return stderrors.New("inline query result contact id is required")
	}
	if strings.TrimSpace(result.PhoneNumber) == "" {
		return stderrors.New("inline query result contact phone_number is required")
	}
	if strings.TrimSpace(result.FirstName) == "" {
		return stderrors.New("inline query result contact first_name is required")
	}
	if err := validateOptionalInlineResultFields(result.ReplyMarkup, result.InputMessageContent, result.ThumbnailURL, result.ThumbnailWidth, result.ThumbnailHeight); err != nil {
		return err
	}
	return nil
}

func validateInlineQueryResultGame(result InlineQueryResultGame) error {
	if strings.TrimSpace(result.ID) == "" {
		return stderrors.New("inline query result game id is required")
	}
	if strings.TrimSpace(result.GameShortName) == "" {
		return stderrors.New("inline query result game_short_name is required")
	}
	if result.ReplyMarkup != nil {
		if err := telegram.ValidateReplyMarkup(*result.ReplyMarkup); err != nil {
			return err
		}
	}
	return nil
}

func validateOptionalInlineResultFields(markup *telegram.InlineKeyboardMarkup, content InputMessageContent, thumbnailURL string, thumbnailWidth int, thumbnailHeight int) error {
	if markup != nil {
		if err := telegram.ValidateReplyMarkup(*markup); err != nil {
			return err
		}
	}
	if content != nil {
		if err := validateInputMessageContent(content); err != nil {
			return err
		}
	}
	if thumbnailURL != "" {
		if err := validateInlineHTTPURL(thumbnailURL, "inline query result thumbnail_url"); err != nil {
			return err
		}
	}
	if thumbnailWidth < 0 {
		return stderrors.New("thumbnail_width must not be negative")
	}
	if thumbnailHeight < 0 {
		return stderrors.New("thumbnail_height must not be negative")
	}
	return nil
}

func validateInputMessageContent(content InputMessageContent) error {
	if content == nil || isNilBotInterfaceValue(content) {
		return stderrors.New("input_message_content is required")
	}

	switch value := content.(type) {
	case InputTextMessageContent:
		return validateInputTextMessageContent(value)
	case *InputTextMessageContent:
		return validateInputTextMessageContent(*value)
	case InputLocationMessageContent:
		return validateInputLocationMessageContent(value)
	case *InputLocationMessageContent:
		return validateInputLocationMessageContent(*value)
	case InputVenueMessageContent:
		return validateInputVenueMessageContent(value)
	case *InputVenueMessageContent:
		return validateInputVenueMessageContent(*value)
	case InputContactMessageContent:
		return validateInputContactMessageContent(value)
	case *InputContactMessageContent:
		return validateInputContactMessageContent(*value)
	case InputInvoiceMessageContent:
		return validateInputInvoiceMessageContent(value)
	case *InputInvoiceMessageContent:
		return validateInputInvoiceMessageContent(*value)
	default:
		return stderrors.New("unsupported input_message_content")
	}
}

func validateInputTextMessageContent(content InputTextMessageContent) error {
	if strings.TrimSpace(content.MessageText) == "" {
		return stderrors.New("message_text is required")
	}
	if err := validateEntityFormatting(content.ParseMode, content.Entities); err != nil {
		return err
	}
	if err := telegram.ValidateLinkPreviewOptions(content.LinkPreviewOptions); err != nil {
		return err
	}

	return nil
}

func validateInputLocationMessageContent(content InputLocationMessageContent) error {
	if err := validateCoordinates(content.Latitude, content.Longitude); err != nil {
		return err
	}
	return validateLocationOptions(content.HorizontalAccuracy, content.LivePeriod, content.Heading, content.ProximityAlertRadius)
}

func validateInputVenueMessageContent(content InputVenueMessageContent) error {
	if err := validateCoordinates(content.Latitude, content.Longitude); err != nil {
		return err
	}
	if strings.TrimSpace(content.Title) == "" {
		return stderrors.New("venue title is required")
	}
	if strings.TrimSpace(content.Address) == "" {
		return stderrors.New("venue address is required")
	}
	return nil
}

func validateInputContactMessageContent(content InputContactMessageContent) error {
	if strings.TrimSpace(content.PhoneNumber) == "" {
		return stderrors.New("contact phone_number is required")
	}
	if strings.TrimSpace(content.FirstName) == "" {
		return stderrors.New("contact first_name is required")
	}
	return nil
}

func validateInputInvoiceMessageContent(content InputInvoiceMessageContent) error {
	if strings.TrimSpace(content.Title) == "" {
		return stderrors.New("invoice title is required")
	}
	if strings.TrimSpace(content.Description) == "" {
		return stderrors.New("invoice description is required")
	}
	if strings.TrimSpace(content.Payload) == "" {
		return stderrors.New("invoice payload is required")
	}
	if strings.TrimSpace(content.Currency) == "" {
		return stderrors.New("invoice currency is required")
	}
	if len(content.Prices) == 0 {
		return stderrors.New("invoice prices must not be empty")
	}
	for index, price := range content.Prices {
		if strings.TrimSpace(price.Label) == "" {
			return fmt.Errorf("invoice prices[%d] label is required", index)
		}
		if price.Amount < 0 {
			return fmt.Errorf("invoice prices[%d] amount must not be negative", index)
		}
	}
	if content.MaxTipAmount < 0 {
		return stderrors.New("max_tip_amount must not be negative")
	}
	for index, amount := range content.SuggestedTipAmounts {
		if amount < 0 {
			return fmt.Errorf("suggested_tip_amounts[%d] must not be negative", index)
		}
	}
	if content.PhotoURL != "" {
		if err := validateInlineHTTPURL(content.PhotoURL, "invoice photo_url"); err != nil {
			return err
		}
	}
	if content.PhotoSize < 0 {
		return stderrors.New("photo_size must not be negative")
	}
	if content.PhotoWidth < 0 {
		return stderrors.New("photo_width must not be negative")
	}
	if content.PhotoHeight < 0 {
		return stderrors.New("photo_height must not be negative")
	}
	return nil
}

func validateCoordinates(latitude float64, longitude float64) error {
	if latitude < -90 || latitude > 90 {
		return stderrors.New("latitude must be between -90 and 90")
	}
	if longitude < -180 || longitude > 180 {
		return stderrors.New("longitude must be between -180 and 180")
	}
	return nil
}

func validateLocationOptions(horizontalAccuracy float64, livePeriod int, heading int, proximityAlertRadius int) error {
	if horizontalAccuracy < 0 {
		return stderrors.New("horizontal_accuracy must not be negative")
	}
	if livePeriod < 0 {
		return stderrors.New("live_period must not be negative")
	}
	if heading < 0 {
		return stderrors.New("heading must not be negative")
	}
	if proximityAlertRadius < 0 {
		return stderrors.New("proximity_alert_radius must not be negative")
	}
	return nil
}

func validateInlineHTTPURL(rawURL string, field string) error {
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return stderrors.New(field + " must be a valid HTTP(S) URL")
	}
	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return stderrors.New(field + " scheme must be http or https")
	}
	if parsed.Host == "" {
		return stderrors.New(field + " host is required")
	}

	return nil
}

func isNilBotInterfaceValue(value any) bool {
	reflected := reflect.ValueOf(value)
	switch reflected.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Pointer, reflect.Slice:
		return reflected.IsNil()
	default:
		return false
	}
}
