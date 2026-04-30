// Package telegram contains data types that mirror Telegram Bot API JSON objects.
package telegram

// Update represents an incoming Telegram update.
type Update struct {
	UpdateID      int64          `json:"update_id"`
	Message       *Message       `json:"message,omitempty"`
	EditedMessage *Message       `json:"edited_message,omitempty"`
	CallbackQuery *CallbackQuery `json:"callback_query,omitempty"`
}

// Message represents a Telegram message with the minimal fields needed by update handlers.
type Message struct {
	MessageID       int64 `json:"message_id"`
	MessageThreadID int64 `json:"message_thread_id,omitempty"`
	From            *User `json:"from,omitempty"`
	Chat            Chat  `json:"chat"`
	Date            int64 `json:"date"`

	Text     string          `json:"text,omitempty"`
	Entities []MessageEntity `json:"entities,omitempty"`

	Caption         string          `json:"caption,omitempty"`
	CaptionEntities []MessageEntity `json:"caption_entities,omitempty"`

	Animation *Animation  `json:"animation,omitempty"`
	Audio     *Audio      `json:"audio,omitempty"`
	Document  *Document   `json:"document,omitempty"`
	Photo     []PhotoSize `json:"photo,omitempty"`
	Sticker   *Sticker    `json:"sticker,omitempty"`
	Video     *Video      `json:"video,omitempty"`
	Voice     *Voice      `json:"voice,omitempty"`

	Contact  *Contact  `json:"contact,omitempty"`
	Location *Location `json:"location,omitempty"`
	Venue    *Venue    `json:"venue,omitempty"`
}

// MessageID contains the identifier of a Telegram message returned by methods that create a copy.
type MessageID struct {
	MessageID int64 `json:"message_id"`
}

// User represents a Telegram user or bot account.
type User struct {
	ID        int64  `json:"id"`
	IsBot     bool   `json:"is_bot"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name,omitempty"`
	Username  string `json:"username,omitempty"`
}

// Chat represents a Telegram chat.
type Chat struct {
	ID        int64  `json:"id"`
	Type      string `json:"type"`
	Title     string `json:"title,omitempty"`
	Username  string `json:"username,omitempty"`
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
}

// MessageEntity represents one special entity in a message text or caption.
type MessageEntity struct {
	Type          string `json:"type"`
	Offset        int    `json:"offset"`
	Length        int    `json:"length"`
	URL           string `json:"url,omitempty"`
	User          *User  `json:"user,omitempty"`
	Language      string `json:"language,omitempty"`
	CustomEmojiID string `json:"custom_emoji_id,omitempty"`
}

// ReplyParameters describes the message being replied to.
type ReplyParameters struct {
	MessageID                int64 `json:"message_id"`
	AllowSendingWithoutReply bool  `json:"allow_sending_without_reply,omitempty"`
}

const (
	// EntityMention marks an @username mention.
	EntityMention = "mention"
	// EntityHashtag marks a hashtag.
	EntityHashtag = "hashtag"
	// EntityCashtag marks a cashtag.
	EntityCashtag = "cashtag"
	// EntityBotCommand marks a bot command.
	EntityBotCommand = "bot_command"
	// EntityURL marks a URL.
	EntityURL = "url"
	// EntityEmail marks an email address.
	EntityEmail = "email"
	// EntityBold marks bold text.
	EntityBold = "bold"
	// EntityItalic marks italic text.
	EntityItalic = "italic"
	// EntityUnderline marks underlined text.
	EntityUnderline = "underline"
	// EntityStrikethrough marks strikethrough text.
	EntityStrikethrough = "strikethrough"
	// EntitySpoiler marks spoiler text.
	EntitySpoiler = "spoiler"
	// EntityCode marks inline code.
	EntityCode = "code"
	// EntityPre marks a preformatted code block.
	EntityPre = "pre"
	// EntityTextLink marks a text link.
	EntityTextLink = "text_link"
	// EntityTextMention marks a text mention of a user.
	EntityTextMention = "text_mention"
	// EntityCustomEmoji marks a custom emoji.
	EntityCustomEmoji = "custom_emoji"
)

// PhotoSize represents one available size of a photo-like file.
type PhotoSize struct {
	FileID       string `json:"file_id"`
	FileUniqueID string `json:"file_unique_id"`
	Width        int    `json:"width"`
	Height       int    `json:"height"`
	FileSize     int64  `json:"file_size,omitempty"`
}

// Animation represents an incoming animation file.
type Animation struct {
	FileID       string     `json:"file_id"`
	FileUniqueID string     `json:"file_unique_id"`
	Width        int        `json:"width"`
	Height       int        `json:"height"`
	Duration     int        `json:"duration"`
	Thumbnail    *PhotoSize `json:"thumbnail,omitempty"`
	FileName     string     `json:"file_name,omitempty"`
	MimeType     string     `json:"mime_type,omitempty"`
	FileSize     int64      `json:"file_size,omitempty"`
}

// Audio represents an incoming audio file.
type Audio struct {
	FileID       string     `json:"file_id"`
	FileUniqueID string     `json:"file_unique_id"`
	Duration     int        `json:"duration"`
	Performer    string     `json:"performer,omitempty"`
	Title        string     `json:"title,omitempty"`
	FileName     string     `json:"file_name,omitempty"`
	MimeType     string     `json:"mime_type,omitempty"`
	FileSize     int64      `json:"file_size,omitempty"`
	Thumbnail    *PhotoSize `json:"thumbnail,omitempty"`
}

// Document represents an incoming general file.
type Document struct {
	FileID       string     `json:"file_id"`
	FileUniqueID string     `json:"file_unique_id"`
	Thumbnail    *PhotoSize `json:"thumbnail,omitempty"`
	FileName     string     `json:"file_name,omitempty"`
	MimeType     string     `json:"mime_type,omitempty"`
	FileSize     int64      `json:"file_size,omitempty"`
}

// Video represents an incoming video file.
type Video struct {
	FileID       string     `json:"file_id"`
	FileUniqueID string     `json:"file_unique_id"`
	Width        int        `json:"width"`
	Height       int        `json:"height"`
	Duration     int        `json:"duration"`
	Thumbnail    *PhotoSize `json:"thumbnail,omitempty"`
	FileName     string     `json:"file_name,omitempty"`
	MimeType     string     `json:"mime_type,omitempty"`
	FileSize     int64      `json:"file_size,omitempty"`
}

// Voice represents an incoming voice message.
type Voice struct {
	FileID       string `json:"file_id"`
	FileUniqueID string `json:"file_unique_id"`
	Duration     int    `json:"duration"`
	MimeType     string `json:"mime_type,omitempty"`
	FileSize     int64  `json:"file_size,omitempty"`
}

// Sticker represents an incoming sticker.
type Sticker struct {
	FileID       string     `json:"file_id"`
	FileUniqueID string     `json:"file_unique_id"`
	Type         string     `json:"type"`
	Width        int        `json:"width"`
	Height       int        `json:"height"`
	IsAnimated   bool       `json:"is_animated"`
	IsVideo      bool       `json:"is_video"`
	Thumbnail    *PhotoSize `json:"thumbnail,omitempty"`
	Emoji        string     `json:"emoji,omitempty"`
	SetName      string     `json:"set_name,omitempty"`
	FileSize     int64      `json:"file_size,omitempty"`
}

// Contact represents a shared phone contact.
type Contact struct {
	PhoneNumber string `json:"phone_number"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name,omitempty"`
	UserID      int64  `json:"user_id,omitempty"`
	VCard       string `json:"vcard,omitempty"`
}

// Location represents a point on the map.
type Location struct {
	Longitude            float64 `json:"longitude"`
	Latitude             float64 `json:"latitude"`
	HorizontalAccuracy   float64 `json:"horizontal_accuracy,omitempty"`
	LivePeriod           int     `json:"live_period,omitempty"`
	Heading              int     `json:"heading,omitempty"`
	ProximityAlertRadius int     `json:"proximity_alert_radius,omitempty"`
}

// Venue represents a venue with a location and address.
type Venue struct {
	Location        Location `json:"location"`
	Title           string   `json:"title"`
	Address         string   `json:"address"`
	FoursquareID    string   `json:"foursquare_id,omitempty"`
	FoursquareType  string   `json:"foursquare_type,omitempty"`
	GooglePlaceID   string   `json:"google_place_id,omitempty"`
	GooglePlaceType string   `json:"google_place_type,omitempty"`
}

// CallbackQuery represents an incoming callback query from an inline keyboard.
type CallbackQuery struct {
	ID              string   `json:"id"`
	From            User     `json:"from"`
	Message         *Message `json:"message,omitempty"`
	InlineMessageID string   `json:"inline_message_id,omitempty"`
	ChatInstance    string   `json:"chat_instance,omitempty"`
	Data            string   `json:"data,omitempty"`
	GameShortName   string   `json:"game_short_name,omitempty"`
}

// File represents a Telegram file metadata object.
type File struct {
	FileID       string `json:"file_id"`
	FileUniqueID string `json:"file_unique_id"`
	FileSize     int64  `json:"file_size,omitempty"`
	FilePath     string `json:"file_path,omitempty"`
}

// WebhookInfo describes current Telegram webhook status.
type WebhookInfo struct {
	URL                          string   `json:"url"`
	HasCustomCertificate         bool     `json:"has_custom_certificate"`
	PendingUpdateCount           int      `json:"pending_update_count"`
	IPAddress                    string   `json:"ip_address,omitempty"`
	LastErrorDate                int64    `json:"last_error_date,omitempty"`
	LastErrorMessage             string   `json:"last_error_message,omitempty"`
	LastSynchronizationErrorDate int64    `json:"last_synchronization_error_date,omitempty"`
	MaxConnections               int      `json:"max_connections,omitempty"`
	AllowedUpdates               []string `json:"allowed_updates,omitempty"`
}
