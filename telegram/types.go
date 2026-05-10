// Package telegram contains data types that mirror Telegram Bot API JSON objects.
package telegram

// Update represents an incoming Telegram update.
type Update struct {
	UpdateID                int64                        `json:"update_id"`
	Message                 *Message                     `json:"message,omitempty"`
	EditedMessage           *Message                     `json:"edited_message,omitempty"`
	ChannelPost             *Message                     `json:"channel_post,omitempty"`
	EditedChannelPost       *Message                     `json:"edited_channel_post,omitempty"`
	BusinessConnection      *BusinessConnection          `json:"business_connection,omitempty"`
	BusinessMessage         *Message                     `json:"business_message,omitempty"`
	EditedBusinessMessage   *Message                     `json:"edited_business_message,omitempty"`
	DeletedBusinessMessages *BusinessMessagesDeleted     `json:"deleted_business_messages,omitempty"`
	GuestMessage            *Message                     `json:"guest_message,omitempty"`
	CallbackQuery           *CallbackQuery               `json:"callback_query,omitempty"`
	InlineQuery             *InlineQuery                 `json:"inline_query,omitempty"`
	ChosenInlineResult      *ChosenInlineResult          `json:"chosen_inline_result,omitempty"`
	ChatJoinRequest         *ChatJoinRequest             `json:"chat_join_request,omitempty"`
	MessageReaction         *MessageReactionUpdated      `json:"message_reaction,omitempty"`
	MessageReactionCount    *MessageReactionCountUpdated `json:"message_reaction_count,omitempty"`
	ShippingQuery           *ShippingQuery               `json:"shipping_query,omitempty"`
	PreCheckoutQuery        *PreCheckoutQuery            `json:"pre_checkout_query,omitempty"`
	PurchasedPaidMedia      *PaidMediaPurchased          `json:"purchased_paid_media,omitempty"`
	Poll                    *Poll                        `json:"poll,omitempty"`
	ManagedBot              *ManagedBotUpdated           `json:"managed_bot,omitempty"`
	PollAnswer              *PollAnswer                  `json:"poll_answer,omitempty"`
	MyChatMember            *ChatMemberUpdated           `json:"my_chat_member,omitempty"`
	ChatMember              *ChatMemberUpdated           `json:"chat_member,omitempty"`
	ChatBoost               *ChatBoostUpdated            `json:"chat_boost,omitempty"`
	RemovedChatBoost        *ChatBoostRemoved            `json:"removed_chat_boost,omitempty"`
}

// Message represents a Telegram message.
type Message struct {
	MessageID            int64                `json:"message_id"`
	MessageThreadID      int64                `json:"message_thread_id,omitempty"`
	DirectMessagesTopic  *DirectMessagesTopic `json:"direct_messages_topic,omitempty"`
	From                 *User                `json:"from,omitempty"`
	SenderChat           *Chat                `json:"sender_chat,omitempty"`
	SenderBoostCount     int                  `json:"sender_boost_count,omitempty"`
	SenderBusinessBot    *User                `json:"sender_business_bot,omitempty"`
	GuestBotCallerUser   *User                `json:"guest_bot_caller_user,omitempty"`
	GuestBotCallerChat   *Chat                `json:"guest_bot_caller_chat,omitempty"`
	GuestQueryID         string               `json:"guest_query_id,omitempty"`
	SenderTag            string               `json:"sender_tag,omitempty"`
	Chat                 Chat                 `json:"chat"`
	Date                 int64                `json:"date"`
	BusinessConnectionID string               `json:"business_connection_id,omitempty"`
	ForwardOrigin        MessageOrigin        `json:"forward_origin,omitempty"`
	IsTopicMessage       bool                 `json:"is_topic_message,omitempty"`
	IsAutomaticForward   bool                 `json:"is_automatic_forward,omitempty"`
	ReplyToMessage       *Message             `json:"reply_to_message,omitempty"`
	ExternalReply        *ExternalReplyInfo   `json:"external_reply,omitempty"`
	Quote                *TextQuote           `json:"quote,omitempty"`
	ReplyToStory         *Story               `json:"reply_to_story,omitempty"`
	ViaBot               *User                `json:"via_bot,omitempty"`
	EditDate             int64                `json:"edit_date,omitempty"`
	HasProtectedContent  bool                 `json:"has_protected_content,omitempty"`
	IsFromOffline        bool                 `json:"is_from_offline,omitempty"`
	IsPaidPost           bool                 `json:"is_paid_post,omitempty"`
	MediaGroupID         string               `json:"media_group_id,omitempty"`
	AuthorSignature      string               `json:"author_signature,omitempty"`
	PaidStarCount        int                  `json:"paid_star_count,omitempty"`

	Text               string              `json:"text,omitempty"`
	Entities           []MessageEntity     `json:"entities,omitempty"`
	LinkPreviewOptions *LinkPreviewOptions `json:"link_preview_options,omitempty"`
	SuggestedPostInfo  *SuggestedPostInfo  `json:"suggested_post_info,omitempty"`
	EffectID           string              `json:"effect_id,omitempty"`

	Caption               string          `json:"caption,omitempty"`
	CaptionEntities       []MessageEntity `json:"caption_entities,omitempty"`
	ShowCaptionAboveMedia bool            `json:"show_caption_above_media,omitempty"`
	HasMediaSpoiler       bool            `json:"has_media_spoiler,omitempty"`

	Animation *Animation  `json:"animation,omitempty"`
	Audio     *Audio      `json:"audio,omitempty"`
	Document  *Document   `json:"document,omitempty"`
	LivePhoto *LivePhoto  `json:"live_photo,omitempty"`
	Photo     []PhotoSize `json:"photo,omitempty"`
	Sticker   *Sticker    `json:"sticker,omitempty"`
	Story     *Story      `json:"story,omitempty"`
	Video     *Video      `json:"video,omitempty"`
	VideoNote *VideoNote  `json:"video_note,omitempty"`
	Voice     *Voice      `json:"voice,omitempty"`

	Contact  *Contact  `json:"contact,omitempty"`
	Dice     *Dice     `json:"dice,omitempty"`
	Game     *Game     `json:"game,omitempty"`
	Location *Location `json:"location,omitempty"`
	Poll     *Poll     `json:"poll,omitempty"`
	Venue    *Venue    `json:"venue,omitempty"`

	NewChatMembers        []User            `json:"new_chat_members,omitempty"`
	LeftChatMember        *User             `json:"left_chat_member,omitempty"`
	ChatOwnerLeft         *ChatOwnerLeft    `json:"chat_owner_left,omitempty"`
	ChatOwnerChanged      *ChatOwnerChanged `json:"chat_owner_changed,omitempty"`
	NewChatTitle          string            `json:"new_chat_title,omitempty"`
	NewChatPhoto          []PhotoSize       `json:"new_chat_photo,omitempty"`
	DeleteChatPhoto       bool              `json:"delete_chat_photo,omitempty"`
	GroupChatCreated      bool              `json:"group_chat_created,omitempty"`
	SupergroupChatCreated bool              `json:"supergroup_chat_created,omitempty"`
	ChannelChatCreated    bool              `json:"channel_chat_created,omitempty"`
	MigrateToChatID       int64             `json:"migrate_to_chat_id,omitempty"`
	MigrateFromChatID     int64             `json:"migrate_from_chat_id,omitempty"`

	Invoice           *Invoice           `json:"invoice,omitempty"`
	SuccessfulPayment *SuccessfulPayment `json:"successful_payment,omitempty"`
	RefundedPayment   *RefundedPayment   `json:"refunded_payment,omitempty"`
	PaidMedia         *PaidMediaInfo     `json:"paid_media,omitempty"`
	PassportData      *PassportData      `json:"passport_data,omitempty"`
	Gift              *GiftInfo          `json:"gift,omitempty"`
	UniqueGift        *UniqueGiftInfo    `json:"unique_gift,omitempty"`
	GiftUpgradeSent   *GiftInfo          `json:"gift_upgrade_sent,omitempty"`

	UsersShared      *UsersShared `json:"users_shared,omitempty"`
	ChatShared       *ChatShared  `json:"chat_shared,omitempty"`
	ConnectedWebsite string       `json:"connected_website,omitempty"`

	ForumTopicCreated         *ForumTopicCreated         `json:"forum_topic_created,omitempty"`
	ForumTopicEdited          *ForumTopicEdited          `json:"forum_topic_edited,omitempty"`
	ForumTopicClosed          *ForumTopicClosed          `json:"forum_topic_closed,omitempty"`
	ForumTopicReopened        *ForumTopicReopened        `json:"forum_topic_reopened,omitempty"`
	GeneralForumTopicHidden   *GeneralForumTopicHidden   `json:"general_forum_topic_hidden,omitempty"`
	GeneralForumTopicUnhidden *GeneralForumTopicUnhidden `json:"general_forum_topic_unhidden,omitempty"`

	PinnedMessage                 *MaybeInaccessibleMessage      `json:"pinned_message,omitempty"`
	ManagedBotCreated             *ManagedBotCreated             `json:"managed_bot_created,omitempty"`
	PollOptionAdded               *PollOptionAdded               `json:"poll_option_added,omitempty"`
	PollOptionDeleted             *PollOptionDeleted             `json:"poll_option_deleted,omitempty"`
	Checklist                     *Checklist                     `json:"checklist,omitempty"`
	ChecklistTasksDone            *ChecklistTasksDone            `json:"checklist_tasks_done,omitempty"`
	ChecklistTasksAdded           *ChecklistTasksAdded           `json:"checklist_tasks_added,omitempty"`
	ProximityAlertTriggered       *ProximityAlertTriggered       `json:"proximity_alert_triggered,omitempty"`
	BoostAdded                    *ChatBoostAdded                `json:"boost_added,omitempty"`
	ChatBackgroundSet             *ChatBackground                `json:"chat_background_set,omitempty"`
	MessageAutoDeleteTimerChanged *MessageAutoDeleteTimerChanged `json:"message_auto_delete_timer_changed,omitempty"`
	GiveawayCreated               *GiveawayCreated               `json:"giveaway_created,omitempty"`
	Giveaway                      *Giveaway                      `json:"giveaway,omitempty"`
	GiveawayWinners               *GiveawayWinners               `json:"giveaway_winners,omitempty"`
	GiveawayCompleted             *GiveawayCompleted             `json:"giveaway_completed,omitempty"`
	PaidMessagePriceChanged       *PaidMessagePriceChanged       `json:"paid_message_price_changed,omitempty"`
	DirectMessagePriceChanged     *DirectMessagePriceChanged     `json:"direct_message_price_changed,omitempty"`
	SuggestedPostApproved         *SuggestedPostApproved         `json:"suggested_post_approved,omitempty"`
	SuggestedPostApprovalFailed   *SuggestedPostApprovalFailed   `json:"suggested_post_approval_failed,omitempty"`
	SuggestedPostDeclined         *SuggestedPostDeclined         `json:"suggested_post_declined,omitempty"`
	SuggestedPostPaid             *SuggestedPostPaid             `json:"suggested_post_paid,omitempty"`
	SuggestedPostRefunded         *SuggestedPostRefunded         `json:"suggested_post_refunded,omitempty"`
	VideoChatScheduled            *VideoChatScheduled            `json:"video_chat_scheduled,omitempty"`
	VideoChatStarted              *VideoChatStarted              `json:"video_chat_started,omitempty"`
	VideoChatEnded                *VideoChatEnded                `json:"video_chat_ended,omitempty"`
	VideoChatParticipantsInvited  *VideoChatParticipantsInvited  `json:"video_chat_participants_invited,omitempty"`
	ReplyToChecklistTaskID        int64                          `json:"reply_to_checklist_task_id,omitempty"`
	ReplyToPollOptionID           string                         `json:"reply_to_poll_option_id,omitempty"`
	WebAppData                    *WebAppData                    `json:"web_app_data,omitempty"`
	WriteAccessAllowed            *WriteAccessAllowed            `json:"write_access_allowed,omitempty"`
	ReplyMarkup                   *InlineKeyboardMarkup          `json:"reply_markup,omitempty"`
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
	ID                        int64  `json:"id"`
	IsBot                     bool   `json:"is_bot"`
	FirstName                 string `json:"first_name"`
	LastName                  string `json:"last_name,omitempty"`
	Username                  string `json:"username,omitempty"`
	LanguageCode              string `json:"language_code,omitempty"`
	IsPremium                 bool   `json:"is_premium,omitempty"`
	AddedToAttachmentMenu     bool   `json:"added_to_attachment_menu,omitempty"`
	CanJoinGroups             bool   `json:"can_join_groups,omitempty"`
	CanReadAllGroupMessages   bool   `json:"can_read_all_group_messages,omitempty"`
	SupportsInlineQueries     bool   `json:"supports_inline_queries,omitempty"`
	SupportsGuestQueries      bool   `json:"supports_guest_queries,omitempty"`
	CanConnectToBusiness      bool   `json:"can_connect_to_business,omitempty"`
	HasMainWebApp             bool   `json:"has_main_web_app,omitempty"`
	HasTopicsEnabled          bool   `json:"has_topics_enabled,omitempty"`
	AllowsUsersToCreateTopics bool   `json:"allows_users_to_create_topics,omitempty"`
	CanManageBots             bool   `json:"can_manage_bots,omitempty"`
}

// Chat represents a Telegram chat.
type Chat struct {
	ID               int64    `json:"id"`
	Type             string   `json:"type"`
	Title            string   `json:"title,omitempty"`
	Username         string   `json:"username,omitempty"`
	FirstName        string   `json:"first_name,omitempty"`
	LastName         string   `json:"last_name,omitempty"`
	IsForum          bool     `json:"is_forum,omitempty"`
	IsDirectMessages bool     `json:"is_direct_messages,omitempty"`
	Description      string   `json:"description,omitempty"`
	InviteLink       string   `json:"invite_link,omitempty"`
	PinnedMessage    *Message `json:"pinned_message,omitempty"`
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
	SubscriptionPeriod      int    `json:"subscription_period,omitempty"`
	SubscriptionPrice       int    `json:"subscription_price,omitempty"`
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

// ChatMember is one of the official Telegram chat member variants.
type ChatMember interface {
	ChatMemberStatus() ChatMemberStatus
	ChatMemberUser() User
	isChatMember()
}

// ChatMemberOwner describes the chat owner.
type ChatMemberOwner struct {
	Status ChatMemberStatus `json:"status"`
	User   User             `json:"user"`

	IsAnonymous bool   `json:"is_anonymous,omitempty"`
	CustomTitle string `json:"custom_title,omitempty"`
}

// ChatMemberAdministrator describes an administrator in a chat.
type ChatMemberAdministrator struct {
	Status ChatMemberStatus `json:"status"`
	User   User             `json:"user"`

	CanBeEdited             bool   `json:"can_be_edited,omitempty"`
	IsAnonymous             bool   `json:"is_anonymous,omitempty"`
	CanManageChat           bool   `json:"can_manage_chat,omitempty"`
	CanDeleteMessages       bool   `json:"can_delete_messages,omitempty"`
	CanManageVideoChats     bool   `json:"can_manage_video_chats,omitempty"`
	CanRestrictMembers      bool   `json:"can_restrict_members,omitempty"`
	CanPromoteMembers       bool   `json:"can_promote_members,omitempty"`
	CanChangeInfo           bool   `json:"can_change_info,omitempty"`
	CanInviteUsers          bool   `json:"can_invite_users,omitempty"`
	CanPostStories          bool   `json:"can_post_stories,omitempty"`
	CanEditStories          bool   `json:"can_edit_stories,omitempty"`
	CanDeleteStories        bool   `json:"can_delete_stories,omitempty"`
	CanPostMessages         bool   `json:"can_post_messages,omitempty"`
	CanEditMessages         bool   `json:"can_edit_messages,omitempty"`
	CanPinMessages          bool   `json:"can_pin_messages,omitempty"`
	CanManageTopics         bool   `json:"can_manage_topics,omitempty"`
	CanManageDirectMessages bool   `json:"can_manage_direct_messages,omitempty"`
	CanManageTags           bool   `json:"can_manage_tags,omitempty"`
	CanEditTag              bool   `json:"can_edit_tag,omitempty"`
	CustomTitle             string `json:"custom_title,omitempty"`
	Tag                     string `json:"tag,omitempty"`
}

// ChatMemberMember describes a regular chat member.
type ChatMemberMember struct {
	Status    ChatMemberStatus `json:"status"`
	User      User             `json:"user"`
	Tag       string           `json:"tag,omitempty"`
	UntilDate int64            `json:"until_date,omitempty"`
}

// ChatMemberRestricted describes a restricted chat member.
type ChatMemberRestricted struct {
	Status ChatMemberStatus `json:"status"`
	User   User             `json:"user"`

	IsMember              bool   `json:"is_member"`
	CanSendMessages       bool   `json:"can_send_messages,omitempty"`
	CanSendAudios         bool   `json:"can_send_audios,omitempty"`
	CanSendDocuments      bool   `json:"can_send_documents,omitempty"`
	CanSendPhotos         bool   `json:"can_send_photos,omitempty"`
	CanSendVideos         bool   `json:"can_send_videos,omitempty"`
	CanSendVideoNotes     bool   `json:"can_send_video_notes,omitempty"`
	CanSendVoiceNotes     bool   `json:"can_send_voice_notes,omitempty"`
	CanSendPolls          bool   `json:"can_send_polls,omitempty"`
	CanSendOtherMessages  bool   `json:"can_send_other_messages,omitempty"`
	CanReactToMessages    bool   `json:"can_react_to_messages,omitempty"`
	CanAddWebPagePreviews bool   `json:"can_add_web_page_previews,omitempty"`
	CanChangeInfo         bool   `json:"can_change_info,omitempty"`
	CanInviteUsers        bool   `json:"can_invite_users,omitempty"`
	CanPinMessages        bool   `json:"can_pin_messages,omitempty"`
	CanManageTopics       bool   `json:"can_manage_topics,omitempty"`
	CanEditTag            bool   `json:"can_edit_tag,omitempty"`
	Tag                   string `json:"tag,omitempty"`
	UntilDate             int64  `json:"until_date"`
}

// ChatMemberLeft describes a user who is not currently a chat member.
type ChatMemberLeft struct {
	Status ChatMemberStatus `json:"status"`
	User   User             `json:"user"`
}

// ChatMemberBanned describes a user banned from a chat.
type ChatMemberBanned struct {
	Status    ChatMemberStatus `json:"status"`
	User      User             `json:"user"`
	UntilDate int64            `json:"until_date"`
}

// ChatMemberUpdated represents changes in the status of a chat member.
type ChatMemberUpdated struct {
	Chat                    Chat            `json:"chat"`
	From                    User            `json:"from"`
	Date                    int64           `json:"date"`
	OldChatMember           ChatMember      `json:"old_chat_member"`
	NewChatMember           ChatMember      `json:"new_chat_member"`
	InviteLink              *ChatInviteLink `json:"invite_link,omitempty"`
	ViaJoinRequest          bool            `json:"via_join_request,omitempty"`
	ViaChatFolderInviteLink bool            `json:"via_chat_folder_invite_link,omitempty"`
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
	CanReactToMessages    bool `json:"can_react_to_messages,omitempty"`
	CanAddWebPagePreviews bool `json:"can_add_web_page_previews,omitempty"`
	CanChangeInfo         bool `json:"can_change_info,omitempty"`
	CanInviteUsers        bool `json:"can_invite_users,omitempty"`
	CanPinMessages        bool `json:"can_pin_messages,omitempty"`
	CanManageTopics       bool `json:"can_manage_topics,omitempty"`
	CanEditTag            bool `json:"can_edit_tag,omitempty"`
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
	MessageID                int64           `json:"message_id"`
	ChatID                   ReplyChatID     `json:"chat_id,omitempty"`
	AllowSendingWithoutReply bool            `json:"allow_sending_without_reply,omitempty"`
	Quote                    string          `json:"quote,omitempty"`
	QuoteParseMode           string          `json:"quote_parse_mode,omitempty"`
	QuoteEntities            []MessageEntity `json:"quote_entities,omitempty"`
	QuotePosition            int             `json:"quote_position,omitempty"`
	ChecklistTaskID          int64           `json:"checklist_task_id,omitempty"`
	PollOptionID             string          `json:"poll_option_id,omitempty"`
}

// ReplyChatID marks chat identifiers accepted by ReplyParameters.ChatID.
type ReplyChatID interface {
	replyChatID()
}

// ReplyChatIDInt identifies a chat by its numeric identifier.
type ReplyChatIDInt int64

// ReplyChatIDUsername identifies a channel by its @username.
type ReplyChatIDUsername string

func (ReplyChatIDInt) replyChatID()      {}
func (ReplyChatIDUsername) replyChatID() {}

// InputPollOption describes one poll option to send.
type InputPollOption struct {
	Text          string          `json:"text"`
	TextParseMode string          `json:"text_parse_mode,omitempty"`
	TextEntities  []MessageEntity `json:"text_entities,omitempty"`
	Media         any             `json:"media,omitempty"`
}

// PollMedia describes media attached to a poll, quiz explanation, or poll option.
type PollMedia struct {
	Animation *Animation  `json:"animation,omitempty"`
	Audio     *Audio      `json:"audio,omitempty"`
	Document  *Document   `json:"document,omitempty"`
	LivePhoto *LivePhoto  `json:"live_photo,omitempty"`
	Location  *Location   `json:"location,omitempty"`
	Photo     []PhotoSize `json:"photo,omitempty"`
	Sticker   *Sticker    `json:"sticker,omitempty"`
	Venue     *Venue      `json:"venue,omitempty"`
	Video     *Video      `json:"video,omitempty"`
}

// PollOption describes one answer option in a Telegram poll.
type PollOption struct {
	PersistentID string          `json:"persistent_id,omitempty"`
	Text         string          `json:"text"`
	TextEntities []MessageEntity `json:"text_entities,omitempty"`
	Media        *PollMedia      `json:"media,omitempty"`
	VoterCount   int             `json:"voter_count"`
	AddedByUser  *User           `json:"added_by_user,omitempty"`
	AddedByChat  *Chat           `json:"added_by_chat,omitempty"`
	AdditionDate int64           `json:"addition_date,omitempty"`
}

// Poll describes a native Telegram poll.
type Poll struct {
	ID                    string          `json:"id"`
	Question              string          `json:"question"`
	QuestionEntities      []MessageEntity `json:"question_entities,omitempty"`
	Options               []PollOption    `json:"options"`
	TotalVoterCount       int             `json:"total_voter_count"`
	IsClosed              bool            `json:"is_closed"`
	IsAnonymous           bool            `json:"is_anonymous"`
	Type                  string          `json:"type"`
	AllowsMultipleAnswers bool            `json:"allows_multiple_answers,omitempty"`
	AllowsRevoting        bool            `json:"allows_revoting,omitempty"`
	MembersOnly           bool            `json:"members_only,omitempty"`
	CountryCodes          []string        `json:"country_codes,omitempty"`
	CorrectOptionID       int             `json:"correct_option_id,omitempty"`
	CorrectOptionIDs      []int           `json:"correct_option_ids,omitempty"`
	Explanation           string          `json:"explanation,omitempty"`
	ExplanationEntities   []MessageEntity `json:"explanation_entities,omitempty"`
	ExplanationMedia      *PollMedia      `json:"explanation_media,omitempty"`
	OpenPeriod            int             `json:"open_period,omitempty"`
	CloseDate             int64           `json:"close_date,omitempty"`
	Description           string          `json:"description,omitempty"`
	DescriptionEntities   []MessageEntity `json:"description_entities,omitempty"`
	Media                 *PollMedia      `json:"media,omitempty"`
}

// PollAnswer represents an answer of a user or anonymous voter in a non-anonymous poll.
type PollAnswer struct {
	PollID              string   `json:"poll_id"`
	VoterChat           *Chat    `json:"voter_chat,omitempty"`
	User                *User    `json:"user,omitempty"`
	OptionIDs           []int    `json:"option_ids"`
	OptionPersistentIDs []string `json:"option_persistent_ids,omitempty"`
}

// InaccessibleMessage describes a message that is inaccessible to the bot.
type InaccessibleMessage struct {
	MessageID int64 `json:"message_id"`
	Chat      Chat  `json:"chat"`
	Date      int64 `json:"date"`
}

// MaybeInaccessibleMessage describes a message that may be inaccessible to the bot.
type MaybeInaccessibleMessage struct {
	MessageID           int64                `json:"message_id"`
	Chat                Chat                 `json:"chat"`
	Date                int64                `json:"date"`
	Message             *Message             `json:"-"`
	InaccessibleMessage *InaccessibleMessage `json:"-"`
}

// PollOptionAdded describes a service message about an option added to a poll.
type PollOptionAdded struct {
	PollMessage        *MaybeInaccessibleMessage `json:"poll_message,omitempty"`
	OptionPersistentID string                    `json:"option_persistent_id"`
	OptionText         string                    `json:"option_text"`
	OptionTextEntities []MessageEntity           `json:"option_text_entities,omitempty"`
}

// PollOptionDeleted describes a service message about an option deleted from a poll.
type PollOptionDeleted struct {
	PollMessage        *MaybeInaccessibleMessage `json:"poll_message,omitempty"`
	OptionPersistentID string                    `json:"option_persistent_id"`
	OptionText         string                    `json:"option_text"`
	OptionTextEntities []MessageEntity           `json:"option_text_entities,omitempty"`
}

// ChecklistTask describes a task in a checklist.
type ChecklistTask struct {
	ID              int64           `json:"id"`
	Text            string          `json:"text"`
	TextEntities    []MessageEntity `json:"text_entities,omitempty"`
	CompletedByUser *User           `json:"completed_by_user,omitempty"`
	CompletedByChat *Chat           `json:"completed_by_chat,omitempty"`
	CompletionDate  int64           `json:"completion_date,omitempty"`
}

// Checklist describes a Telegram checklist.
type Checklist struct {
	Title                    string          `json:"title"`
	TitleEntities            []MessageEntity `json:"title_entities,omitempty"`
	Tasks                    []ChecklistTask `json:"tasks"`
	OthersCanAddTasks        bool            `json:"others_can_add_tasks,omitempty"`
	OthersCanMarkTasksAsDone bool            `json:"others_can_mark_tasks_as_done,omitempty"`
}

// InputChecklistTask describes a checklist task to create.
type InputChecklistTask struct {
	ID           int64           `json:"id"`
	Text         string          `json:"text"`
	ParseMode    string          `json:"parse_mode,omitempty"`
	TextEntities []MessageEntity `json:"text_entities,omitempty"`
}

// InputChecklist describes a checklist to create.
type InputChecklist struct {
	Title                    string               `json:"title"`
	ParseMode                string               `json:"parse_mode,omitempty"`
	TitleEntities            []MessageEntity      `json:"title_entities,omitempty"`
	Tasks                    []InputChecklistTask `json:"tasks"`
	OthersCanAddTasks        bool                 `json:"others_can_add_tasks,omitempty"`
	OthersCanMarkTasksAsDone bool                 `json:"others_can_mark_tasks_as_done,omitempty"`
}

// ChecklistTasksDone describes a service message about checklist tasks marked done or not done.
type ChecklistTasksDone struct {
	ChecklistMessage       *Message `json:"checklist_message,omitempty"`
	MarkedAsDoneTaskIDs    []int64  `json:"marked_as_done_task_ids,omitempty"`
	MarkedAsNotDoneTaskIDs []int64  `json:"marked_as_not_done_task_ids,omitempty"`
}

// ChecklistTasksAdded describes a service message about tasks added to a checklist.
type ChecklistTasksAdded struct {
	ChecklistMessage *Message        `json:"checklist_message,omitempty"`
	Tasks            []ChecklistTask `json:"tasks"`
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

// UserProfilePhotos represents a user's profile pictures.
type UserProfilePhotos struct {
	TotalCount int           `json:"total_count"`
	Photos     [][]PhotoSize `json:"photos"`
}

// UserProfileAudios represents the audios displayed on a user's profile.
type UserProfileAudios struct {
	TotalCount int     `json:"total_count"`
	Audios     []Audio `json:"audios"`
}

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

// LivePhoto represents an incoming live photo.
type LivePhoto struct {
	Photo        []PhotoSize `json:"photo,omitempty"`
	FileID       string      `json:"file_id"`
	FileUniqueID string      `json:"file_unique_id"`
	Width        int         `json:"width"`
	Height       int         `json:"height"`
	Duration     int         `json:"duration"`
	MimeType     string      `json:"mime_type,omitempty"`
	FileSize     int64       `json:"file_size,omitempty"`
}

// Video represents an incoming video file.
type Video struct {
	FileID         string         `json:"file_id"`
	FileUniqueID   string         `json:"file_unique_id"`
	Width          int            `json:"width"`
	Height         int            `json:"height"`
	Duration       int            `json:"duration"`
	Thumbnail      *PhotoSize     `json:"thumbnail,omitempty"`
	Cover          []PhotoSize    `json:"cover,omitempty"`
	StartTimestamp int            `json:"start_timestamp,omitempty"`
	Qualities      []VideoQuality `json:"qualities,omitempty"`
	FileName       string         `json:"file_name,omitempty"`
	MimeType       string         `json:"mime_type,omitempty"`
	FileSize       int64          `json:"file_size,omitempty"`
}

// VideoQuality represents a specific available quality of a video.
type VideoQuality struct {
	FileID       string `json:"file_id"`
	FileUniqueID string `json:"file_unique_id"`
	Width        int    `json:"width"`
	Height       int    `json:"height"`
	Codec        string `json:"codec"`
	FileSize     int64  `json:"file_size,omitempty"`
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
	ID              string                    `json:"id"`
	From            User                      `json:"from"`
	Message         *MaybeInaccessibleMessage `json:"message,omitempty"`
	InlineMessageID string                    `json:"inline_message_id,omitempty"`
	ChatInstance    string                    `json:"chat_instance,omitempty"`
	Data            string                    `json:"data,omitempty"`
	GameShortName   string                    `json:"game_short_name,omitempty"`
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
