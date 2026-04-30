package telegram

import (
	"bytes"
	"encoding/json"
	stderrors "errors"
	"fmt"
	"strings"
)

const (
	reactionTypeEmojiType       = "emoji"
	reactionTypeCustomEmojiType = "custom_emoji"
	reactionTypePaidType        = "paid"
)

// ReactionType marks Telegram message reaction type objects.
type ReactionType interface {
	reactionType()
}

// ReactionTypeEmoji describes an emoji-based reaction.
type ReactionTypeEmoji struct {
	Type  string `json:"type"`
	Emoji string `json:"emoji"`
}

// ReactionTypeCustomEmoji describes a custom emoji reaction.
type ReactionTypeCustomEmoji struct {
	Type          string `json:"type"`
	CustomEmojiID string `json:"custom_emoji_id"`
}

// ReactionTypePaid describes a paid reaction.
type ReactionTypePaid struct {
	Type string `json:"type"`
}

func (ReactionTypeEmoji) reactionType()       {}
func (ReactionTypeCustomEmoji) reactionType() {}
func (ReactionTypePaid) reactionType()        {}

// NewReactionTypeEmoji creates an emoji reaction type.
func NewReactionTypeEmoji(emoji string) ReactionTypeEmoji {
	return ReactionTypeEmoji{Type: reactionTypeEmojiType, Emoji: emoji}
}

// NewReactionTypeCustomEmoji creates a custom emoji reaction type.
func NewReactionTypeCustomEmoji(customEmojiID string) ReactionTypeCustomEmoji {
	return ReactionTypeCustomEmoji{Type: reactionTypeCustomEmojiType, CustomEmojiID: customEmojiID}
}

// NewReactionTypePaid creates a paid reaction type.
func NewReactionTypePaid() ReactionTypePaid {
	return ReactionTypePaid{Type: reactionTypePaidType}
}

// MarshalJSON encodes ReactionTypeEmoji with the required Telegram type field.
func (r ReactionTypeEmoji) MarshalJSON() ([]byte, error) {
	r.Type = reactionTypeEmojiType
	type reaction ReactionTypeEmoji
	return json.Marshal(reaction(r))
}

// MarshalJSON encodes ReactionTypeCustomEmoji with the required Telegram type field.
func (r ReactionTypeCustomEmoji) MarshalJSON() ([]byte, error) {
	r.Type = reactionTypeCustomEmojiType
	type reaction ReactionTypeCustomEmoji
	return json.Marshal(reaction(r))
}

// MarshalJSON encodes ReactionTypePaid with the required Telegram type field.
func (r ReactionTypePaid) MarshalJSON() ([]byte, error) {
	r.Type = reactionTypePaidType
	type reaction ReactionTypePaid
	return json.Marshal(reaction(r))
}

// ValidateReactionTypes checks whether reactions can be encoded for Telegram.
func ValidateReactionTypes(reactions []ReactionType) error {
	for index, reaction := range reactions {
		if err := ValidateReactionType(reaction); err != nil {
			return fmt.Errorf("reaction[%d] is invalid: %w", index, err)
		}
	}
	return nil
}

// ValidateReactionType checks whether reaction can be encoded for Telegram.
func ValidateReactionType(reaction ReactionType) error {
	if reaction == nil || isNilInterfaceValue(reaction) {
		return stderrors.New("reaction type must not be nil")
	}

	switch value := reaction.(type) {
	case ReactionTypeEmoji:
		return validateReactionTypeEmoji(value)
	case *ReactionTypeEmoji:
		return validateReactionTypeEmoji(*value)
	case ReactionTypeCustomEmoji:
		return validateReactionTypeCustomEmoji(value)
	case *ReactionTypeCustomEmoji:
		return validateReactionTypeCustomEmoji(*value)
	case ReactionTypePaid, *ReactionTypePaid:
		return nil
	default:
		return stderrors.New("unsupported reaction type")
	}
}

// UnmarshalReactionType decodes a polymorphic Telegram ReactionType object.
func UnmarshalReactionType(data []byte) (ReactionType, error) {
	var meta struct {
		Type string `json:"type"`
	}
	if err := json.Unmarshal(data, &meta); err != nil {
		return nil, err
	}

	switch meta.Type {
	case reactionTypeEmojiType:
		var reaction ReactionTypeEmoji
		if err := json.Unmarshal(data, &reaction); err != nil {
			return nil, err
		}
		reaction.Type = reactionTypeEmojiType
		return reaction, nil
	case reactionTypeCustomEmojiType:
		var reaction ReactionTypeCustomEmoji
		if err := json.Unmarshal(data, &reaction); err != nil {
			return nil, err
		}
		reaction.Type = reactionTypeCustomEmojiType
		return reaction, nil
	case reactionTypePaidType:
		var reaction ReactionTypePaid
		if err := json.Unmarshal(data, &reaction); err != nil {
			return nil, err
		}
		reaction.Type = reactionTypePaidType
		return reaction, nil
	default:
		return nil, stderrors.New("unsupported reaction type")
	}
}

// ReactionCount represents a reaction added to a message and how many times it was added.
type ReactionCount struct {
	Type       ReactionType `json:"type"`
	TotalCount int          `json:"total_count"`
}

// MessageReactionUpdated represents a non-anonymous message reaction change update.
type MessageReactionUpdated struct {
	Chat        Chat           `json:"chat"`
	MessageID   int64          `json:"message_id"`
	User        *User          `json:"user,omitempty"`
	ActorChat   *Chat          `json:"actor_chat,omitempty"`
	Date        int64          `json:"date"`
	OldReaction []ReactionType `json:"old_reaction"`
	NewReaction []ReactionType `json:"new_reaction"`
}

// MessageReactionCountUpdated represents anonymous reaction count changes on a message.
type MessageReactionCountUpdated struct {
	Chat      Chat            `json:"chat"`
	MessageID int64           `json:"message_id"`
	Date      int64           `json:"date"`
	Reactions []ReactionCount `json:"reactions"`
}

// UnmarshalJSON decodes a ReactionCount with a polymorphic reaction type.
func (r *ReactionCount) UnmarshalJSON(data []byte) error {
	var payload struct {
		Type       json.RawMessage `json:"type"`
		TotalCount int             `json:"total_count"`
	}
	if err := json.Unmarshal(data, &payload); err != nil {
		return err
	}
	if len(payload.Type) == 0 || bytes.Equal(payload.Type, []byte("null")) {
		return stderrors.New("reaction count type is required")
	}
	reactionType, err := UnmarshalReactionType(payload.Type)
	if err != nil {
		return err
	}

	r.Type = reactionType
	r.TotalCount = payload.TotalCount
	return nil
}

// UnmarshalJSON decodes MessageReactionUpdated with polymorphic old and new reactions.
func (m *MessageReactionUpdated) UnmarshalJSON(data []byte) error {
	var payload struct {
		Chat        Chat              `json:"chat"`
		MessageID   int64             `json:"message_id"`
		User        *User             `json:"user,omitempty"`
		ActorChat   *Chat             `json:"actor_chat,omitempty"`
		Date        int64             `json:"date"`
		OldReaction []json.RawMessage `json:"old_reaction"`
		NewReaction []json.RawMessage `json:"new_reaction"`
	}
	if err := json.Unmarshal(data, &payload); err != nil {
		return err
	}

	oldReaction, err := unmarshalReactionTypes(payload.OldReaction)
	if err != nil {
		return fmt.Errorf("old_reaction is invalid: %w", err)
	}
	newReaction, err := unmarshalReactionTypes(payload.NewReaction)
	if err != nil {
		return fmt.Errorf("new_reaction is invalid: %w", err)
	}

	m.Chat = payload.Chat
	m.MessageID = payload.MessageID
	m.User = payload.User
	m.ActorChat = payload.ActorChat
	m.Date = payload.Date
	m.OldReaction = oldReaction
	m.NewReaction = newReaction
	return nil
}

func validateReactionTypeEmoji(reaction ReactionTypeEmoji) error {
	if strings.TrimSpace(reaction.Emoji) == "" {
		return stderrors.New("reaction emoji is required")
	}
	return nil
}

func validateReactionTypeCustomEmoji(reaction ReactionTypeCustomEmoji) error {
	if strings.TrimSpace(reaction.CustomEmojiID) == "" {
		return stderrors.New("reaction custom_emoji_id is required")
	}
	return nil
}

func unmarshalReactionTypes(rawItems []json.RawMessage) ([]ReactionType, error) {
	if rawItems == nil {
		return nil, nil
	}
	reactions := make([]ReactionType, 0, len(rawItems))
	for index, raw := range rawItems {
		reaction, err := UnmarshalReactionType(raw)
		if err != nil {
			return nil, fmt.Errorf("reaction[%d]: %w", index, err)
		}
		reactions = append(reactions, reaction)
	}
	return reactions, nil
}
