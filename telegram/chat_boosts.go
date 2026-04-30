package telegram

import (
	"bytes"
	"encoding/json"
	stderrors "errors"
	"fmt"
)

const (
	chatBoostSourcePremiumType  = "premium"
	chatBoostSourceGiftCodeType = "gift_code"
	chatBoostSourceGiveawayType = "giveaway"
)

// ChatBoostUpdated represents a boost added to a chat or changed.
type ChatBoostUpdated struct {
	Chat  Chat      `json:"chat"`
	Boost ChatBoost `json:"boost"`
}

// ChatBoostRemoved represents a boost removed from a chat.
type ChatBoostRemoved struct {
	Chat       Chat            `json:"chat"`
	BoostID    string          `json:"boost_id"`
	RemoveDate int64           `json:"remove_date"`
	Source     ChatBoostSource `json:"source"`
}

// ChatBoost contains information about a chat boost.
type ChatBoost struct {
	BoostID        string          `json:"boost_id"`
	AddDate        int64           `json:"add_date"`
	ExpirationDate int64           `json:"expiration_date"`
	Source         ChatBoostSource `json:"source"`
}

// ChatBoostSource marks Telegram chat boost source objects.
type ChatBoostSource interface {
	chatBoostSource()
}

// ChatBoostSourcePremium describes a boost from Telegram Premium.
type ChatBoostSourcePremium struct {
	Source string `json:"source"`
	User   User   `json:"user"`
}

// ChatBoostSourceGiftCode describes a boost from Premium gift codes.
type ChatBoostSourceGiftCode struct {
	Source string `json:"source"`
	User   User   `json:"user"`
}

// ChatBoostSourceGiveaway describes a boost from a Premium or Star giveaway.
type ChatBoostSourceGiveaway struct {
	Source            string `json:"source"`
	GiveawayMessageID int64  `json:"giveaway_message_id"`
	User              *User  `json:"user,omitempty"`
	PrizeStarCount    int    `json:"prize_star_count,omitempty"`
	IsUnclaimed       bool   `json:"is_unclaimed,omitempty"`
}

// UserChatBoosts represents a list of boosts added to a chat by a user.
type UserChatBoosts struct {
	Boosts []ChatBoost `json:"boosts"`
}

func (ChatBoostSourcePremium) chatBoostSource()  {}
func (ChatBoostSourceGiftCode) chatBoostSource() {}
func (ChatBoostSourceGiveaway) chatBoostSource() {}

// UnmarshalChatBoostSource decodes a polymorphic Telegram ChatBoostSource object.
func UnmarshalChatBoostSource(data []byte) (ChatBoostSource, error) {
	var meta struct {
		Source string `json:"source"`
	}
	if err := json.Unmarshal(data, &meta); err != nil {
		return nil, err
	}

	switch meta.Source {
	case chatBoostSourcePremiumType:
		var source ChatBoostSourcePremium
		if err := json.Unmarshal(data, &source); err != nil {
			return nil, err
		}
		source.Source = chatBoostSourcePremiumType
		return source, nil
	case chatBoostSourceGiftCodeType:
		var source ChatBoostSourceGiftCode
		if err := json.Unmarshal(data, &source); err != nil {
			return nil, err
		}
		source.Source = chatBoostSourceGiftCodeType
		return source, nil
	case chatBoostSourceGiveawayType:
		var source ChatBoostSourceGiveaway
		if err := json.Unmarshal(data, &source); err != nil {
			return nil, err
		}
		source.Source = chatBoostSourceGiveawayType
		return source, nil
	default:
		return nil, stderrors.New("unsupported chat boost source")
	}
}

// UnmarshalJSON decodes a ChatBoost with a polymorphic boost source.
func (boost *ChatBoost) UnmarshalJSON(data []byte) error {
	var payload struct {
		BoostID        string          `json:"boost_id"`
		AddDate        int64           `json:"add_date"`
		ExpirationDate int64           `json:"expiration_date"`
		Source         json.RawMessage `json:"source"`
	}
	if err := json.Unmarshal(data, &payload); err != nil {
		return err
	}
	if len(payload.Source) == 0 || bytes.Equal(payload.Source, []byte("null")) {
		return stderrors.New("chat boost source is required")
	}
	source, err := UnmarshalChatBoostSource(payload.Source)
	if err != nil {
		return fmt.Errorf("source is invalid: %w", err)
	}

	boost.BoostID = payload.BoostID
	boost.AddDate = payload.AddDate
	boost.ExpirationDate = payload.ExpirationDate
	boost.Source = source
	return nil
}

// UnmarshalJSON decodes a ChatBoostRemoved with a polymorphic boost source.
func (removed *ChatBoostRemoved) UnmarshalJSON(data []byte) error {
	var payload struct {
		Chat       Chat            `json:"chat"`
		BoostID    string          `json:"boost_id"`
		RemoveDate int64           `json:"remove_date"`
		Source     json.RawMessage `json:"source"`
	}
	if err := json.Unmarshal(data, &payload); err != nil {
		return err
	}
	if len(payload.Source) == 0 || bytes.Equal(payload.Source, []byte("null")) {
		return stderrors.New("chat boost source is required")
	}
	source, err := UnmarshalChatBoostSource(payload.Source)
	if err != nil {
		return fmt.Errorf("source is invalid: %w", err)
	}

	removed.Chat = payload.Chat
	removed.BoostID = payload.BoostID
	removed.RemoveDate = payload.RemoveDate
	removed.Source = source
	return nil
}
