package telegram

import (
	"bytes"
	"encoding/json"
	stderrors "errors"
	"fmt"
)

const (
	ownedGiftRegularType = "regular"
	ownedGiftUniqueType  = "unique"
)

// GiftBackground describes the background of a regular gift.
type GiftBackground struct {
	CenterColor int `json:"center_color"`
	EdgeColor   int `json:"edge_color"`
	TextColor   int `json:"text_color"`
}

// Gift represents a gift that can be sent by the bot.
type Gift struct {
	ID                     string          `json:"id"`
	Sticker                Sticker         `json:"sticker"`
	StarCount              int             `json:"star_count"`
	UpgradeStarCount       int             `json:"upgrade_star_count,omitempty"`
	IsPremium              bool            `json:"is_premium,omitempty"`
	HasColors              bool            `json:"has_colors,omitempty"`
	TotalCount             int             `json:"total_count,omitempty"`
	RemainingCount         int             `json:"remaining_count,omitempty"`
	PersonalTotalCount     int             `json:"personal_total_count,omitempty"`
	PersonalRemainingCount int             `json:"personal_remaining_count,omitempty"`
	Background             *GiftBackground `json:"background,omitempty"`
	UniqueGiftVariantCount int             `json:"unique_gift_variant_count,omitempty"`
	PublisherChat          *Chat           `json:"publisher_chat,omitempty"`
}

// Gifts contains a list of gifts available to send.
type Gifts struct {
	Gifts []Gift `json:"gifts"`
}

// UniqueGiftModel describes the model of a unique gift.
type UniqueGiftModel struct {
	Name           string  `json:"name"`
	Sticker        Sticker `json:"sticker"`
	RarityPerMille int     `json:"rarity_per_mille"`
	Rarity         string  `json:"rarity,omitempty"`
}

// UniqueGiftSymbol describes the symbol shown on the pattern of a unique gift.
type UniqueGiftSymbol struct {
	Name           string  `json:"name"`
	Sticker        Sticker `json:"sticker"`
	RarityPerMille int     `json:"rarity_per_mille"`
}

// UniqueGiftBackdropColors describes backdrop colors of a unique gift.
type UniqueGiftBackdropColors struct {
	CenterColor int `json:"center_color"`
	EdgeColor   int `json:"edge_color"`
	SymbolColor int `json:"symbol_color"`
	TextColor   int `json:"text_color"`
}

// UniqueGiftBackdrop describes the backdrop of a unique gift.
type UniqueGiftBackdrop struct {
	Name           string                   `json:"name"`
	Colors         UniqueGiftBackdropColors `json:"colors"`
	RarityPerMille int                      `json:"rarity_per_mille"`
}

// UniqueGiftColors contains color scheme information derived from a unique gift.
type UniqueGiftColors struct {
	ModelCustomEmojiID    string `json:"model_custom_emoji_id"`
	SymbolCustomEmojiID   string `json:"symbol_custom_emoji_id"`
	LightThemeMainColor   int    `json:"light_theme_main_color"`
	LightThemeOtherColors []int  `json:"light_theme_other_colors"`
	DarkThemeMainColor    int    `json:"dark_theme_main_color"`
	DarkThemeOtherColors  []int  `json:"dark_theme_other_colors"`
}

// UniqueGift describes a unique gift upgraded from a regular gift.
type UniqueGift struct {
	GiftID           string             `json:"gift_id"`
	BaseName         string             `json:"base_name"`
	Name             string             `json:"name"`
	Number           int                `json:"number"`
	Model            UniqueGiftModel    `json:"model"`
	Symbol           UniqueGiftSymbol   `json:"symbol"`
	Backdrop         UniqueGiftBackdrop `json:"backdrop"`
	IsPremium        bool               `json:"is_premium,omitempty"`
	IsBurned         bool               `json:"is_burned,omitempty"`
	IsFromBlockchain bool               `json:"is_from_blockchain,omitempty"`
	Colors           *UniqueGiftColors  `json:"colors,omitempty"`
	PublisherChat    *Chat              `json:"publisher_chat,omitempty"`
}

// GiftInfo describes a service message about a regular gift.
type GiftInfo struct {
	Gift                    Gift            `json:"gift"`
	OwnedGiftID             string          `json:"owned_gift_id,omitempty"`
	ConvertStarCount        int             `json:"convert_star_count,omitempty"`
	PrepaidUpgradeStarCount int             `json:"prepaid_upgrade_star_count,omitempty"`
	IsUpgradeSeparate       bool            `json:"is_upgrade_separate,omitempty"`
	CanBeUpgraded           bool            `json:"can_be_upgraded,omitempty"`
	Text                    string          `json:"text,omitempty"`
	Entities                []MessageEntity `json:"entities,omitempty"`
	IsPrivate               bool            `json:"is_private,omitempty"`
	UniqueGiftNumber        int             `json:"unique_gift_number,omitempty"`
}

// UniqueGiftInfo describes a service message about a unique gift.
type UniqueGiftInfo struct {
	Gift               UniqueGift `json:"gift"`
	Origin             string     `json:"origin"`
	LastResaleCurrency string     `json:"last_resale_currency,omitempty"`
	LastResaleAmount   int        `json:"last_resale_amount,omitempty"`
	OwnedGiftID        string     `json:"owned_gift_id,omitempty"`
	TransferStarCount  int        `json:"transfer_star_count,omitempty"`
	NextTransferDate   int64      `json:"next_transfer_date,omitempty"`
}

// OwnedGift marks Telegram owned gift objects.
type OwnedGift interface {
	ownedGift()
}

// OwnedGiftRegular describes a regular gift owned by a user or chat.
type OwnedGiftRegular struct {
	Type                    string          `json:"type"`
	Gift                    Gift            `json:"gift"`
	OwnedGiftID             string          `json:"owned_gift_id,omitempty"`
	SenderUser              *User           `json:"sender_user,omitempty"`
	SendDate                int64           `json:"send_date"`
	Text                    string          `json:"text,omitempty"`
	Entities                []MessageEntity `json:"entities,omitempty"`
	IsPrivate               bool            `json:"is_private,omitempty"`
	IsSaved                 bool            `json:"is_saved,omitempty"`
	CanBeUpgraded           bool            `json:"can_be_upgraded,omitempty"`
	WasRefunded             bool            `json:"was_refunded,omitempty"`
	ConvertStarCount        int             `json:"convert_star_count,omitempty"`
	PrepaidUpgradeStarCount int             `json:"prepaid_upgrade_star_count,omitempty"`
	IsUpgradeSeparate       bool            `json:"is_upgrade_separate,omitempty"`
	UniqueGiftNumber        int             `json:"unique_gift_number,omitempty"`
}

// OwnedGiftUnique describes a unique gift owned by a user or chat.
type OwnedGiftUnique struct {
	Type              string     `json:"type"`
	Gift              UniqueGift `json:"gift"`
	OwnedGiftID       string     `json:"owned_gift_id,omitempty"`
	SenderUser        *User      `json:"sender_user,omitempty"`
	SendDate          int64      `json:"send_date"`
	IsSaved           bool       `json:"is_saved,omitempty"`
	CanBeTransferred  bool       `json:"can_be_transferred,omitempty"`
	TransferStarCount int        `json:"transfer_star_count,omitempty"`
	NextTransferDate  int64      `json:"next_transfer_date,omitempty"`
}

// OwnedGifts contains gifts owned by a user or chat.
type OwnedGifts struct {
	TotalCount int         `json:"total_count"`
	Gifts      []OwnedGift `json:"gifts"`
	NextOffset string      `json:"next_offset,omitempty"`
}

func (OwnedGiftRegular) ownedGift() {}
func (OwnedGiftUnique) ownedGift()  {}

// UnmarshalOwnedGift decodes a polymorphic OwnedGift object.
func UnmarshalOwnedGift(data []byte) (OwnedGift, error) {
	var meta struct {
		Type string `json:"type"`
	}
	if err := json.Unmarshal(data, &meta); err != nil {
		return nil, err
	}

	switch meta.Type {
	case ownedGiftRegularType:
		var gift OwnedGiftRegular
		if err := json.Unmarshal(data, &gift); err != nil {
			return nil, err
		}
		gift.Type = ownedGiftRegularType
		return gift, nil
	case ownedGiftUniqueType:
		var gift OwnedGiftUnique
		if err := json.Unmarshal(data, &gift); err != nil {
			return nil, err
		}
		gift.Type = ownedGiftUniqueType
		return gift, nil
	default:
		return nil, stderrors.New("unsupported owned gift type")
	}
}

// UnmarshalJSON decodes OwnedGifts with polymorphic gift entries.
func (gifts *OwnedGifts) UnmarshalJSON(data []byte) error {
	var payload struct {
		TotalCount int               `json:"total_count"`
		Gifts      []json.RawMessage `json:"gifts"`
		NextOffset string            `json:"next_offset,omitempty"`
	}
	if err := json.Unmarshal(data, &payload); err != nil {
		return err
	}
	owned, err := unmarshalOwnedGiftItems(payload.Gifts)
	if err != nil {
		return err
	}
	gifts.TotalCount = payload.TotalCount
	gifts.Gifts = owned
	gifts.NextOffset = payload.NextOffset
	return nil
}

func unmarshalOwnedGiftItems(rawItems []json.RawMessage) ([]OwnedGift, error) {
	if rawItems == nil {
		return nil, nil
	}
	gifts := make([]OwnedGift, 0, len(rawItems))
	for index, raw := range rawItems {
		if len(raw) == 0 || bytes.Equal(raw, []byte("null")) {
			return nil, fmt.Errorf("gifts[%d] is required", index)
		}
		gift, err := UnmarshalOwnedGift(raw)
		if err != nil {
			return nil, fmt.Errorf("gifts[%d]: %w", index, err)
		}
		gifts = append(gifts, gift)
	}
	return gifts, nil
}
