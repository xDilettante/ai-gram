// Package telegram contains data types that mirror Telegram Bot API JSON objects.
package telegram

// Update represents an incoming Telegram update.
type Update struct {
	UpdateID             int64                        `json:"update_id"`
	Message              *Message                     `json:"message,omitempty"`
	EditedMessage        *Message                     `json:"edited_message,omitempty"`
	CallbackQuery        *CallbackQuery               `json:"callback_query,omitempty"`
	InlineQuery          *InlineQuery                 `json:"inline_query,omitempty"`
	ChosenInlineResult   *ChosenInlineResult          `json:"chosen_inline_result,omitempty"`
	ChatJoinRequest      *ChatJoinRequest             `json:"chat_join_request,omitempty"`
	MessageReaction      *MessageReactionUpdated      `json:"message_reaction,omitempty"`
	MessageReactionCount *MessageReactionCountUpdated `json:"message_reaction_count,omitempty"`
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
	VideoNote *VideoNote  `json:"video_note,omitempty"`
	Voice     *Voice      `json:"voice,omitempty"`

	Contact  *Contact  `json:"contact,omitempty"`
	Dice     *Dice     `json:"dice,omitempty"`
	Location *Location `json:"location,omitempty"`
	Poll     *Poll     `json:"poll,omitempty"`
	Venue    *Venue    `json:"venue,omitempty"`

	ForumTopicCreated         *ForumTopicCreated         `json:"forum_topic_created,omitempty"`
	ForumTopicEdited          *ForumTopicEdited          `json:"forum_topic_edited,omitempty"`
	ForumTopicClosed          *ForumTopicClosed          `json:"forum_topic_closed,omitempty"`
	ForumTopicReopened        *ForumTopicReopened        `json:"forum_topic_reopened,omitempty"`
	GeneralForumTopicHidden   *GeneralForumTopicHidden   `json:"general_forum_topic_hidden,omitempty"`
	GeneralForumTopicUnhidden *GeneralForumTopicUnhidden `json:"general_forum_topic_unhidden,omitempty"`
}

// ForumTopic represents a forum topic in a Telegram supergroup.
type ForumTopic struct {
	MessageThreadID   int64  `json:"message_thread_id"`
	Name              string `json:"name"`
	IconColor         int    `json:"icon_color"`
	IconCustomEmojiID string `json:"icon_custom_emoji_id,omitempty"`
}

// ForumTopicCreated describes a service message about a newly created forum topic.
type ForumTopicCreated struct {
	Name              string `json:"name"`
	IconColor         int    `json:"icon_color"`
	IconCustomEmojiID string `json:"icon_custom_emoji_id,omitempty"`
}

// ForumTopicEdited describes a service message about an edited forum topic.
type ForumTopicEdited struct {
	Name              string `json:"name,omitempty"`
	IconCustomEmojiID string `json:"icon_custom_emoji_id,omitempty"`
}

// ForumTopicClosed describes a service message about a closed forum topic.
type ForumTopicClosed struct{}

// ForumTopicReopened describes a service message about a reopened forum topic.
type ForumTopicReopened struct{}

// GeneralForumTopicHidden describes a service message about the hidden General forum topic.
type GeneralForumTopicHidden struct{}

// GeneralForumTopicUnhidden describes a service message about the unhidden General forum topic.
type GeneralForumTopicUnhidden struct{}

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
	ID            int64    `json:"id"`
	Type          string   `json:"type"`
	Title         string   `json:"title,omitempty"`
	Username      string   `json:"username,omitempty"`
	FirstName     string   `json:"first_name,omitempty"`
	LastName      string   `json:"last_name,omitempty"`
	Description   string   `json:"description,omitempty"`
	InviteLink    string   `json:"invite_link,omitempty"`
	PinnedMessage *Message `json:"pinned_message,omitempty"`
}

// ChatInviteLink describes a Telegram chat invite link.
type ChatInviteLink struct {
	InviteLink              string `json:"invite_link"`
	Creator                 User   `json:"creator"`
	CreatesJoinRequest      bool   `json:"creates_join_request"`
	IsPrimary               bool   `json:"is_primary"`
	IsRevoked               bool   `json:"is_revoked"`
	Name                    string `json:"name,omitempty"`
	ExpireDate              int64  `json:"expire_date,omitempty"`
	MemberLimit             int    `json:"member_limit,omitempty"`
	PendingJoinRequestCount int    `json:"pending_join_request_count,omitempty"`
}

// ChatJoinRequest describes a request to join a chat.
type ChatJoinRequest struct {
	Chat       Chat            `json:"chat"`
	From       User            `json:"from"`
	UserChatID int64           `json:"user_chat_id"`
	Date       int64           `json:"date"`
	Bio        string          `json:"bio,omitempty"`
	InviteLink *ChatInviteLink `json:"invite_link,omitempty"`
}

// InlineQuery represents an incoming inline query.
type InlineQuery struct {
	ID       string    `json:"id"`
	From     User      `json:"from"`
	Query    string    `json:"query"`
	Offset   string    `json:"offset"`
	ChatType string    `json:"chat_type,omitempty"`
	Location *Location `json:"location,omitempty"`
}

// ChosenInlineResult represents a result chosen by a user from an inline query.
type ChosenInlineResult struct {
	ResultID        string    `json:"result_id"`
	From            User      `json:"from"`
	Location        *Location `json:"location,omitempty"`
	InlineMessageID string    `json:"inline_message_id,omitempty"`
	Query           string    `json:"query"`
}

// ChatMemberStatus identifies a user's membership state in a chat.
type ChatMemberStatus string

const (
	// ChatMemberStatusCreator means the user owns the chat.
	ChatMemberStatusCreator ChatMemberStatus = "creator"
	// ChatMemberStatusAdministrator means the user is a chat administrator.
	ChatMemberStatusAdministrator ChatMemberStatus = "administrator"
	// ChatMemberStatusMember means the user is a regular chat member.
	ChatMemberStatusMember ChatMemberStatus = "member"
	// ChatMemberStatusRestricted means the user is restricted in the chat.
	ChatMemberStatusRestricted ChatMemberStatus = "restricted"
	// ChatMemberStatusLeft means the user is not currently a member.
	ChatMemberStatusLeft ChatMemberStatus = "left"
	// ChatMemberStatusKicked means the user was removed from the chat.
	ChatMemberStatusKicked ChatMemberStatus = "kicked"
)

// ChatMember describes a Telegram user's membership and relevant permissions in a chat.
type ChatMember struct {
	Status ChatMemberStatus `json:"status"`
	User   User             `json:"user"`

	IsAnonymous bool   `json:"is_anonymous,omitempty"`
	CustomTitle string `json:"custom_title,omitempty"`
	UntilDate   int64  `json:"until_date,omitempty"`

	CanBeEdited         bool `json:"can_be_edited,omitempty"`
	CanManageChat       bool `json:"can_manage_chat,omitempty"`
	CanDeleteMessages   bool `json:"can_delete_messages,omitempty"`
	CanManageVideoChats bool `json:"can_manage_video_chats,omitempty"`
	CanRestrictMembers  bool `json:"can_restrict_members,omitempty"`
	CanPromoteMembers   bool `json:"can_promote_members,omitempty"`
	CanChangeInfo       bool `json:"can_change_info,omitempty"`
	CanInviteUsers      bool `json:"can_invite_users,omitempty"`
	CanPinMessages      bool `json:"can_pin_messages,omitempty"`
	CanPostStories      bool `json:"can_post_stories,omitempty"`
	CanEditStories      bool `json:"can_edit_stories,omitempty"`
	CanDeleteStories    bool `json:"can_delete_stories,omitempty"`
}

// ChatPermissions describes actions a user is allowed to take in a chat.
type ChatPermissions struct {
	CanSendMessages       bool `json:"can_send_messages,omitempty"`
	CanSendAudios         bool `json:"can_send_audios,omitempty"`
	CanSendDocuments      bool `json:"can_send_documents,omitempty"`
	CanSendPhotos         bool `json:"can_send_photos,omitempty"`
	CanSendVideos         bool `json:"can_send_videos,omitempty"`
	CanSendVideoNotes     bool `json:"can_send_video_notes,omitempty"`
	CanSendVoiceNotes     bool `json:"can_send_voice_notes,omitempty"`
	CanSendPolls          bool `json:"can_send_polls,omitempty"`
	CanSendOtherMessages  bool `json:"can_send_other_messages,omitempty"`
	CanAddWebPagePreviews bool `json:"can_add_web_page_previews,omitempty"`
	CanChangeInfo         bool `json:"can_change_info,omitempty"`
	CanInviteUsers        bool `json:"can_invite_users,omitempty"`
	CanPinMessages        bool `json:"can_pin_messages,omitempty"`
	CanManageTopics       bool `json:"can_manage_topics,omitempty"`
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

// PollOption describes one answer option in a Telegram poll.
type PollOption struct {
	Text       string `json:"text"`
	VoterCount int    `json:"voter_count"`
}

// Poll describes a native Telegram poll.
type Poll struct {
	ID                    string          `json:"id"`
	Question              string          `json:"question"`
	Options               []PollOption    `json:"options"`
	TotalVoterCount       int             `json:"total_voter_count"`
	IsClosed              bool            `json:"is_closed"`
	IsAnonymous           bool            `json:"is_anonymous"`
	Type                  string          `json:"type"`
	AllowsMultipleAnswers bool            `json:"allows_multiple_answers,omitempty"`
	CorrectOptionID       int             `json:"correct_option_id,omitempty"`
	Explanation           string          `json:"explanation,omitempty"`
	ExplanationEntities   []MessageEntity `json:"explanation_entities,omitempty"`
	OpenPeriod            int             `json:"open_period,omitempty"`
	CloseDate             int64           `json:"close_date,omitempty"`
}

// Dice describes a Telegram dice message result.
type Dice struct {
	Emoji string `json:"emoji"`
	Value int    `json:"value"`
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

// VideoNote represents an incoming video note message.
type VideoNote struct {
	FileID       string     `json:"file_id"`
	FileUniqueID string     `json:"file_unique_id,omitempty"`
	Length       int        `json:"length"`
	Duration     int        `json:"duration"`
	Thumbnail    *PhotoSize `json:"thumbnail,omitempty"`
	FileSize     int64      `json:"file_size,omitempty"`
}

// Sticker represents an incoming sticker.
type Sticker struct {
	FileID           string        `json:"file_id"`
	FileUniqueID     string        `json:"file_unique_id"`
	Type             string        `json:"type"`
	Width            int           `json:"width"`
	Height           int           `json:"height"`
	IsAnimated       bool          `json:"is_animated"`
	IsVideo          bool          `json:"is_video"`
	Thumbnail        *PhotoSize    `json:"thumbnail,omitempty"`
	Emoji            string        `json:"emoji,omitempty"`
	SetName          string        `json:"set_name,omitempty"`
	PremiumAnimation *File         `json:"premium_animation,omitempty"`
	MaskPosition     *MaskPosition `json:"mask_position,omitempty"`
	CustomEmojiID    string        `json:"custom_emoji_id,omitempty"`
	NeedsRepainting  bool          `json:"needs_repainting,omitempty"`
	FileSize         int64         `json:"file_size,omitempty"`
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
