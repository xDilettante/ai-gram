package telegram

import (
	"bytes"
	"encoding/json"
	stderrors "errors"
)

const (
	messageOriginUserType       = "user"
	messageOriginHiddenUserType = "hidden_user"
	messageOriginChatType       = "chat"
	messageOriginChannelType    = "channel"
)

// MessageOrigin marks Telegram message origin objects.
type MessageOrigin interface {
	messageOrigin()
}

// MessageOriginUser describes a message originally sent by a known user.
type MessageOriginUser struct {
	Type       string `json:"type"`
	Date       int64  `json:"date"`
	SenderUser User   `json:"sender_user"`
}

// MessageOriginHiddenUser describes a message originally sent by an unknown user.
type MessageOriginHiddenUser struct {
	Type           string `json:"type"`
	Date           int64  `json:"date"`
	SenderUserName string `json:"sender_user_name"`
}

// MessageOriginChat describes a message originally sent on behalf of a chat.
type MessageOriginChat struct {
	Type            string `json:"type"`
	Date            int64  `json:"date"`
	SenderChat      Chat   `json:"sender_chat"`
	AuthorSignature string `json:"author_signature,omitempty"`
}

// MessageOriginChannel describes a message originally sent to a channel chat.
type MessageOriginChannel struct {
	Type            string `json:"type"`
	Date            int64  `json:"date"`
	Chat            Chat   `json:"chat"`
	MessageID       int64  `json:"message_id"`
	AuthorSignature string `json:"author_signature,omitempty"`
}

func (MessageOriginUser) messageOrigin()       {}
func (MessageOriginHiddenUser) messageOrigin() {}
func (MessageOriginChat) messageOrigin()       {}
func (MessageOriginChannel) messageOrigin()    {}

// UnmarshalMessageOrigin decodes a polymorphic Telegram MessageOrigin object.
func UnmarshalMessageOrigin(data []byte) (MessageOrigin, error) {
	var meta struct {
		Type string `json:"type"`
	}
	if err := json.Unmarshal(data, &meta); err != nil {
		return nil, err
	}

	switch meta.Type {
	case messageOriginUserType:
		var origin MessageOriginUser
		if err := json.Unmarshal(data, &origin); err != nil {
			return nil, err
		}
		origin.Type = messageOriginUserType
		return origin, nil
	case messageOriginHiddenUserType:
		var origin MessageOriginHiddenUser
		if err := json.Unmarshal(data, &origin); err != nil {
			return nil, err
		}
		origin.Type = messageOriginHiddenUserType
		return origin, nil
	case messageOriginChatType:
		var origin MessageOriginChat
		if err := json.Unmarshal(data, &origin); err != nil {
			return nil, err
		}
		origin.Type = messageOriginChatType
		return origin, nil
	case messageOriginChannelType:
		var origin MessageOriginChannel
		if err := json.Unmarshal(data, &origin); err != nil {
			return nil, err
		}
		origin.Type = messageOriginChannelType
		return origin, nil
	default:
		return nil, stderrors.New("unsupported message origin type")
	}
}

// TextQuote contains information about the quoted part of a replied-to message.
type TextQuote struct {
	Text     string          `json:"text"`
	Entities []MessageEntity `json:"entities,omitempty"`
	Position int             `json:"position"`
	IsManual bool            `json:"is_manual,omitempty"`
}

// ExternalReplyInfo describes a message being replied to from another chat or topic.
type ExternalReplyInfo struct {
	Origin             MessageOrigin       `json:"origin"`
	Chat               *Chat               `json:"chat,omitempty"`
	MessageID          int64               `json:"message_id,omitempty"`
	LinkPreviewOptions *LinkPreviewOptions `json:"link_preview_options,omitempty"`
	Animation          *Animation          `json:"animation,omitempty"`
	Audio              *Audio              `json:"audio,omitempty"`
	Document           *Document           `json:"document,omitempty"`
	LivePhoto          *LivePhoto          `json:"live_photo,omitempty"`
	PaidMedia          *PaidMediaInfo      `json:"paid_media,omitempty"`
	Photo              []PhotoSize         `json:"photo,omitempty"`
	Sticker            *Sticker            `json:"sticker,omitempty"`
	Story              *Story              `json:"story,omitempty"`
	Video              *Video              `json:"video,omitempty"`
	VideoNote          *VideoNote          `json:"video_note,omitempty"`
	Voice              *Voice              `json:"voice,omitempty"`
	HasMediaSpoiler    bool                `json:"has_media_spoiler,omitempty"`
	Checklist          *Checklist          `json:"checklist,omitempty"`
	Contact            *Contact            `json:"contact,omitempty"`
	Dice               *Dice               `json:"dice,omitempty"`
	Game               *Game               `json:"game,omitempty"`
	Giveaway           *Giveaway           `json:"giveaway,omitempty"`
	GiveawayWinners    *GiveawayWinners    `json:"giveaway_winners,omitempty"`
	Invoice            *Invoice            `json:"invoice,omitempty"`
	Location           *Location           `json:"location,omitempty"`
	Poll               *Poll               `json:"poll,omitempty"`
	Venue              *Venue              `json:"venue,omitempty"`
}

// Giveaway represents a scheduled giveaway message.
type Giveaway struct {
	Chats                         []Chat   `json:"chats"`
	WinnersSelectionDate          int64    `json:"winners_selection_date"`
	WinnerCount                   int      `json:"winner_count"`
	OnlyNewMembers                bool     `json:"only_new_members,omitempty"`
	HasPublicWinners              bool     `json:"has_public_winners,omitempty"`
	PrizeDescription              string   `json:"prize_description,omitempty"`
	CountryCodes                  []string `json:"country_codes,omitempty"`
	PrizeStarCount                int      `json:"prize_star_count,omitempty"`
	PremiumSubscriptionMonthCount int      `json:"premium_subscription_month_count,omitempty"`
}

// GiveawayWinners represents a completed giveaway with public winners.
type GiveawayWinners struct {
	Chat                          Chat   `json:"chat"`
	GiveawayMessageID             int64  `json:"giveaway_message_id"`
	WinnersSelectionDate          int64  `json:"winners_selection_date"`
	WinnerCount                   int    `json:"winner_count"`
	Winners                       []User `json:"winners"`
	AdditionalChatCount           int    `json:"additional_chat_count,omitempty"`
	PrizeStarCount                int    `json:"prize_star_count,omitempty"`
	PremiumSubscriptionMonthCount int    `json:"premium_subscription_month_count,omitempty"`
	UnclaimedPrizeCount           int    `json:"unclaimed_prize_count,omitempty"`
	OnlyNewMembers                bool   `json:"only_new_members,omitempty"`
	WasRefunded                   bool   `json:"was_refunded,omitempty"`
	PrizeDescription              string `json:"prize_description,omitempty"`
}

// SuggestedPostInfo contains metadata about a suggested post message.
type SuggestedPostInfo struct {
	State    string              `json:"state"`
	Price    *SuggestedPostPrice `json:"price,omitempty"`
	SendDate int64               `json:"send_date,omitempty"`
}

// DirectMessagesTopic describes a topic of a channel direct messages chat.
type DirectMessagesTopic struct {
	TopicID int64 `json:"topic_id"`
	User    *User `json:"user,omitempty"`
}

// UnmarshalJSON decodes ExternalReplyInfo with a polymorphic message origin.
func (info *ExternalReplyInfo) UnmarshalJSON(data []byte) error {
	type externalReplyInfoAlias ExternalReplyInfo
	payload := struct {
		Origin json.RawMessage `json:"origin"`
		*externalReplyInfoAlias
	}{externalReplyInfoAlias: (*externalReplyInfoAlias)(info)}
	if err := json.Unmarshal(data, &payload); err != nil {
		return err
	}
	if len(payload.Origin) > 0 && !bytes.Equal(payload.Origin, []byte("null")) {
		origin, err := UnmarshalMessageOrigin(payload.Origin)
		if err != nil {
			return err
		}
		info.Origin = origin
	}
	return nil
}

// UnmarshalJSON decodes a Message with polymorphic metadata fields.
func (m *Message) UnmarshalJSON(data []byte) error {
	type messageAlias Message
	payload := struct {
		ForwardOrigin json.RawMessage `json:"forward_origin"`
		PinnedMessage json.RawMessage `json:"pinned_message"`
		*messageAlias
	}{messageAlias: (*messageAlias)(m)}
	if err := json.Unmarshal(data, &payload); err != nil {
		return err
	}
	if len(payload.ForwardOrigin) > 0 && !bytes.Equal(payload.ForwardOrigin, []byte("null")) {
		origin, err := UnmarshalMessageOrigin(payload.ForwardOrigin)
		if err != nil {
			return err
		}
		m.ForwardOrigin = origin
	}
	if len(payload.PinnedMessage) > 0 && !bytes.Equal(payload.PinnedMessage, []byte("null")) {
		var pinned MaybeInaccessibleMessage
		if err := json.Unmarshal(payload.PinnedMessage, &pinned); err != nil {
			return err
		}
		m.PinnedMessage = &pinned
	}
	return nil
}

// UnmarshalJSON decodes a MaybeInaccessibleMessage as either Message or InaccessibleMessage.
func (m *MaybeInaccessibleMessage) UnmarshalJSON(data []byte) error {
	var meta struct {
		MessageID int64 `json:"message_id"`
		Chat      Chat  `json:"chat"`
		Date      int64 `json:"date"`
	}
	if err := json.Unmarshal(data, &meta); err != nil {
		return err
	}
	m.MessageID = meta.MessageID
	m.Chat = meta.Chat
	m.Date = meta.Date
	m.Message = nil
	m.InaccessibleMessage = nil
	if meta.Date == 0 {
		var inaccessible InaccessibleMessage
		if err := json.Unmarshal(data, &inaccessible); err != nil {
			return err
		}
		m.InaccessibleMessage = &inaccessible
		return nil
	}
	var message Message
	if err := json.Unmarshal(data, &message); err != nil {
		return err
	}
	m.Message = &message
	m.MessageID = message.MessageID
	m.Chat = message.Chat
	m.Date = message.Date
	return nil
}

// MarshalJSON encodes the accessible or inaccessible message payload.
func (m MaybeInaccessibleMessage) MarshalJSON() ([]byte, error) {
	if m.Message != nil {
		return json.Marshal(m.Message)
	}
	if m.InaccessibleMessage != nil {
		return json.Marshal(m.InaccessibleMessage)
	}
	type maybeInaccessibleMessage MaybeInaccessibleMessage
	return json.Marshal(struct {
		Message *Message             `json:"-"`
		Hidden  *InaccessibleMessage `json:"-"`
		maybeInaccessibleMessage
	}{maybeInaccessibleMessage: maybeInaccessibleMessage(m)})
}

// UnmarshalJSON decodes CallbackQuery.message without breaking the legacy Message field.
func (q *CallbackQuery) UnmarshalJSON(data []byte) error {
	type callbackQueryAlias CallbackQuery
	payload := struct {
		Message json.RawMessage `json:"message"`
		*callbackQueryAlias
	}{callbackQueryAlias: (*callbackQueryAlias)(q)}
	if err := json.Unmarshal(data, &payload); err != nil {
		return err
	}
	q.Message = nil
	q.MaybeMessage = nil
	if len(payload.Message) > 0 && !bytes.Equal(payload.Message, []byte("null")) {
		var message MaybeInaccessibleMessage
		if err := json.Unmarshal(payload.Message, &message); err != nil {
			return err
		}
		q.MaybeMessage = &message
		q.Message = message.Message
	}
	return nil
}
