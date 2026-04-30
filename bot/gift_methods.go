package bot

import (
	"context"
	stderrors "errors"
	"fmt"
	"strings"

	"github.com/xDilettante/ai-gram/telegram"
)

// GetAvailableGiftsParams contains supported parameters for getAvailableGifts.
type GetAvailableGiftsParams struct{}

// SendGiftParams contains supported parameters for sendGift.
type SendGiftParams struct {
	UserID        int64                    `json:"user_id,omitempty"`
	ChatID        ChatID                   `json:"-"`
	GiftID        string                   `json:"gift_id"`
	PayForUpgrade bool                     `json:"pay_for_upgrade,omitempty"`
	Text          string                   `json:"text,omitempty"`
	TextParseMode string                   `json:"text_parse_mode,omitempty"`
	TextEntities  []telegram.MessageEntity `json:"text_entities,omitempty"`
}

// GiftPremiumSubscriptionParams contains supported parameters for giftPremiumSubscription.
type GiftPremiumSubscriptionParams struct {
	UserID        int64                    `json:"user_id"`
	MonthCount    int                      `json:"month_count"`
	StarCount     int                      `json:"star_count"`
	Text          string                   `json:"text,omitempty"`
	TextParseMode string                   `json:"text_parse_mode,omitempty"`
	TextEntities  []telegram.MessageEntity `json:"text_entities,omitempty"`
}

// GetBusinessAccountStarBalanceParams contains supported parameters for getBusinessAccountStarBalance.
type GetBusinessAccountStarBalanceParams struct {
	BusinessConnectionID string `json:"business_connection_id"`
}

// TransferBusinessAccountStarsParams contains supported parameters for transferBusinessAccountStars.
type TransferBusinessAccountStarsParams struct {
	BusinessConnectionID string `json:"business_connection_id"`
	StarCount            int    `json:"star_count"`
}

// GetBusinessAccountGiftsParams contains supported parameters for getBusinessAccountGifts.
type GetBusinessAccountGiftsParams struct {
	BusinessConnectionID        string `json:"business_connection_id"`
	ExcludeUnsaved              bool   `json:"exclude_unsaved,omitempty"`
	ExcludeSaved                bool   `json:"exclude_saved,omitempty"`
	ExcludeUnlimited            bool   `json:"exclude_unlimited,omitempty"`
	ExcludeLimitedUpgradable    bool   `json:"exclude_limited_upgradable,omitempty"`
	ExcludeLimitedNonUpgradable bool   `json:"exclude_limited_non_upgradable,omitempty"`
	ExcludeUnique               bool   `json:"exclude_unique,omitempty"`
	ExcludeFromBlockchain       bool   `json:"exclude_from_blockchain,omitempty"`
	SortByPrice                 bool   `json:"sort_by_price,omitempty"`
	Offset                      string `json:"offset,omitempty"`
	Limit                       int    `json:"limit,omitempty"`
}

// GetUserGiftsParams contains supported parameters for getUserGifts.
type GetUserGiftsParams struct {
	UserID                      int64  `json:"user_id"`
	ExcludeUnlimited            bool   `json:"exclude_unlimited,omitempty"`
	ExcludeLimitedUpgradable    bool   `json:"exclude_limited_upgradable,omitempty"`
	ExcludeLimitedNonUpgradable bool   `json:"exclude_limited_non_upgradable,omitempty"`
	ExcludeUnique               bool   `json:"exclude_unique,omitempty"`
	ExcludeFromBlockchain       bool   `json:"exclude_from_blockchain,omitempty"`
	SortByPrice                 bool   `json:"sort_by_price,omitempty"`
	Offset                      string `json:"offset,omitempty"`
	Limit                       int    `json:"limit,omitempty"`
}

// GetChatGiftsParams contains supported parameters for getChatGifts.
type GetChatGiftsParams struct {
	ChatID                      ChatID `json:"chat_id"`
	ExcludeUnsaved              bool   `json:"exclude_unsaved,omitempty"`
	ExcludeSaved                bool   `json:"exclude_saved,omitempty"`
	ExcludeUnlimited            bool   `json:"exclude_unlimited,omitempty"`
	ExcludeLimitedUpgradable    bool   `json:"exclude_limited_upgradable,omitempty"`
	ExcludeLimitedNonUpgradable bool   `json:"exclude_limited_non_upgradable,omitempty"`
	ExcludeUnique               bool   `json:"exclude_unique,omitempty"`
	ExcludeFromBlockchain       bool   `json:"exclude_from_blockchain,omitempty"`
	SortByPrice                 bool   `json:"sort_by_price,omitempty"`
	Offset                      string `json:"offset,omitempty"`
	Limit                       int    `json:"limit,omitempty"`
}

// ConvertGiftToStarsParams contains supported parameters for convertGiftToStars.
type ConvertGiftToStarsParams struct {
	BusinessConnectionID string `json:"business_connection_id"`
	OwnedGiftID          string `json:"owned_gift_id"`
}

// UpgradeGiftParams contains supported parameters for upgradeGift.
type UpgradeGiftParams struct {
	BusinessConnectionID string `json:"business_connection_id"`
	OwnedGiftID          string `json:"owned_gift_id"`
	KeepOriginalDetails  bool   `json:"keep_original_details,omitempty"`
	StarCount            int    `json:"star_count,omitempty"`
}

// TransferGiftParams contains supported parameters for transferGift.
type TransferGiftParams struct {
	BusinessConnectionID string `json:"business_connection_id"`
	OwnedGiftID          string `json:"owned_gift_id"`
	NewOwnerChatID       int64  `json:"new_owner_chat_id"`
	StarCount            int    `json:"star_count,omitempty"`
}

// GetMyStarBalanceParams contains supported parameters for getMyStarBalance.
type GetMyStarBalanceParams struct{}

// EditUserStarSubscriptionParams contains supported parameters for editUserStarSubscription.
type EditUserStarSubscriptionParams struct {
	UserID                  int64  `json:"user_id"`
	TelegramPaymentChargeID string `json:"telegram_payment_charge_id"`
	IsCanceled              bool   `json:"is_canceled"`
}

// GetAvailableGifts returns gifts that can be sent by the bot.
func (b *Bot) GetAvailableGifts(ctx context.Context, params GetAvailableGiftsParams) (*telegram.Gifts, error) {
	var gifts telegram.Gifts
	if err := b.call(ctx, "getAvailableGifts", params, &gifts); err != nil {
		return nil, err
	}
	return &gifts, nil
}

// SendGift sends a gift to a user or channel chat.
func (b *Bot) SendGift(ctx context.Context, params SendGiftParams) (bool, error) {
	if err := params.validate(); err != nil {
		return false, err
	}
	var result bool
	if err := b.call(ctx, "sendGift", params.payload(), &result); err != nil {
		return false, err
	}
	return result, nil
}

// GiftPremiumSubscription gifts a Telegram Premium subscription to a user.
func (b *Bot) GiftPremiumSubscription(ctx context.Context, params GiftPremiumSubscriptionParams) (bool, error) {
	if err := params.validate(); err != nil {
		return false, err
	}
	var result bool
	if err := b.call(ctx, "giftPremiumSubscription", params, &result); err != nil {
		return false, err
	}
	return result, nil
}

// GetBusinessAccountStarBalance returns the Star balance of a business account.
func (b *Bot) GetBusinessAccountStarBalance(ctx context.Context, params GetBusinessAccountStarBalanceParams) (*telegram.StarAmount, error) {
	if err := params.validate(); err != nil {
		return nil, err
	}
	var amount telegram.StarAmount
	if err := b.call(ctx, "getBusinessAccountStarBalance", params, &amount); err != nil {
		return nil, err
	}
	return &amount, nil
}

// TransferBusinessAccountStars transfers Stars from a business account to the bot.
func (b *Bot) TransferBusinessAccountStars(ctx context.Context, params TransferBusinessAccountStarsParams) (bool, error) {
	if err := params.validate(); err != nil {
		return false, err
	}
	var result bool
	if err := b.call(ctx, "transferBusinessAccountStars", params, &result); err != nil {
		return false, err
	}
	return result, nil
}

// GetBusinessAccountGifts returns gifts owned by a business account.
func (b *Bot) GetBusinessAccountGifts(ctx context.Context, params GetBusinessAccountGiftsParams) (*telegram.OwnedGifts, error) {
	if err := params.validate(); err != nil {
		return nil, err
	}
	var gifts telegram.OwnedGifts
	if err := b.call(ctx, "getBusinessAccountGifts", params, &gifts); err != nil {
		return nil, err
	}
	return &gifts, nil
}

// GetUserGifts returns gifts owned and hosted by a user.
func (b *Bot) GetUserGifts(ctx context.Context, params GetUserGiftsParams) (*telegram.OwnedGifts, error) {
	if err := params.validate(); err != nil {
		return nil, err
	}
	var gifts telegram.OwnedGifts
	if err := b.call(ctx, "getUserGifts", params, &gifts); err != nil {
		return nil, err
	}
	return &gifts, nil
}

// GetChatGifts returns gifts owned by a chat.
func (b *Bot) GetChatGifts(ctx context.Context, params GetChatGiftsParams) (*telegram.OwnedGifts, error) {
	if err := params.validate(); err != nil {
		return nil, err
	}
	var gifts telegram.OwnedGifts
	if err := b.call(ctx, "getChatGifts", params, &gifts); err != nil {
		return nil, err
	}
	return &gifts, nil
}

// ConvertGiftToStars converts a regular business gift to Telegram Stars.
func (b *Bot) ConvertGiftToStars(ctx context.Context, params ConvertGiftToStarsParams) (bool, error) {
	return b.businessGiftBool(ctx, "convertGiftToStars", params, params.validate)
}

// UpgradeGift upgrades a regular business gift to a unique gift.
func (b *Bot) UpgradeGift(ctx context.Context, params UpgradeGiftParams) (bool, error) {
	if err := params.validate(); err != nil {
		return false, err
	}
	var result bool
	if err := b.call(ctx, "upgradeGift", params, &result); err != nil {
		return false, err
	}
	return result, nil
}

// TransferGift transfers a unique business gift to another owner.
func (b *Bot) TransferGift(ctx context.Context, params TransferGiftParams) (bool, error) {
	if err := params.validate(); err != nil {
		return false, err
	}
	var result bool
	if err := b.call(ctx, "transferGift", params, &result); err != nil {
		return false, err
	}
	return result, nil
}

// GetMyStarBalance returns the bot's current Telegram Stars balance.
func (b *Bot) GetMyStarBalance(ctx context.Context, params GetMyStarBalanceParams) (*telegram.StarAmount, error) {
	var amount telegram.StarAmount
	if err := b.call(ctx, "getMyStarBalance", params, &amount); err != nil {
		return nil, err
	}
	return &amount, nil
}

// EditUserStarSubscription cancels or re-enables a Telegram Stars subscription.
func (b *Bot) EditUserStarSubscription(ctx context.Context, params EditUserStarSubscriptionParams) (bool, error) {
	if err := params.validate(); err != nil {
		return false, err
	}
	var result bool
	if err := b.call(ctx, "editUserStarSubscription", params, &result); err != nil {
		return false, err
	}
	return result, nil
}

type sendGiftPayload struct {
	UserID        int64                    `json:"user_id,omitempty"`
	ChatID        *ChatID                  `json:"chat_id,omitempty"`
	GiftID        string                   `json:"gift_id"`
	PayForUpgrade bool                     `json:"pay_for_upgrade,omitempty"`
	Text          string                   `json:"text,omitempty"`
	TextParseMode string                   `json:"text_parse_mode,omitempty"`
	TextEntities  []telegram.MessageEntity `json:"text_entities,omitempty"`
}

func (params SendGiftParams) payload() sendGiftPayload {
	payload := sendGiftPayload{
		UserID:        params.UserID,
		GiftID:        params.GiftID,
		PayForUpgrade: params.PayForUpgrade,
		Text:          params.Text,
		TextParseMode: params.TextParseMode,
		TextEntities:  params.TextEntities,
	}
	if params.ChatID.valid() {
		chatID := params.ChatID
		payload.ChatID = &chatID
	}
	return payload
}

func (params SendGiftParams) validate() error {
	hasUser := params.UserID > 0
	hasChat := params.ChatID.valid()
	if hasUser == hasChat {
		return stderrors.New("exactly one of user_id or chat_id is required")
	}
	if params.UserID < 0 {
		return stderrors.New("user_id must be greater than zero")
	}
	if strings.TrimSpace(params.GiftID) == "" {
		return stderrors.New("gift_id is required")
	}
	return validateTextFormatting(params.TextParseMode, params.TextEntities)
}

func (params GiftPremiumSubscriptionParams) validate() error {
	if params.UserID <= 0 {
		return stderrors.New("user_id must be greater than zero")
	}
	expectedStars := map[int]int{3: 1000, 6: 1500, 12: 2500}
	stars, ok := expectedStars[params.MonthCount]
	if !ok {
		return stderrors.New("month_count must be one of 3, 6, or 12")
	}
	if params.StarCount != stars {
		return fmt.Errorf("star_count must be %d for %d months", stars, params.MonthCount)
	}
	return validateTextFormatting(params.TextParseMode, params.TextEntities)
}

func (params GetBusinessAccountStarBalanceParams) validate() error {
	return validateBusinessConnectionID(params.BusinessConnectionID)
}

func (params TransferBusinessAccountStarsParams) validate() error {
	if err := validateBusinessConnectionID(params.BusinessConnectionID); err != nil {
		return err
	}
	if params.StarCount < 1 || params.StarCount > 10000 {
		return stderrors.New("star_count must be between 1 and 10000")
	}
	return nil
}

func (params GetBusinessAccountGiftsParams) validate() error {
	if err := validateBusinessConnectionID(params.BusinessConnectionID); err != nil {
		return err
	}
	return validateGiftListLimit(params.Limit)
}

func (params GetUserGiftsParams) validate() error {
	if params.UserID <= 0 {
		return stderrors.New("user_id must be greater than zero")
	}
	return validateGiftListLimit(params.Limit)
}

func (params GetChatGiftsParams) validate() error {
	if !params.ChatID.valid() {
		return stderrors.New("chat_id is required")
	}
	return validateGiftListLimit(params.Limit)
}

func (params ConvertGiftToStarsParams) validate() error {
	return validateBusinessOwnedGift(params.BusinessConnectionID, params.OwnedGiftID)
}

func (params UpgradeGiftParams) validate() error {
	if err := validateBusinessOwnedGift(params.BusinessConnectionID, params.OwnedGiftID); err != nil {
		return err
	}
	if params.StarCount < 0 {
		return stderrors.New("star_count must not be negative")
	}
	return nil
}

func (params TransferGiftParams) validate() error {
	if err := validateBusinessOwnedGift(params.BusinessConnectionID, params.OwnedGiftID); err != nil {
		return err
	}
	if params.NewOwnerChatID == 0 {
		return stderrors.New("new_owner_chat_id is required")
	}
	if params.StarCount < 0 {
		return stderrors.New("star_count must not be negative")
	}
	return nil
}

func (params EditUserStarSubscriptionParams) validate() error {
	if params.UserID <= 0 {
		return stderrors.New("user_id must be greater than zero")
	}
	if strings.TrimSpace(params.TelegramPaymentChargeID) == "" {
		return stderrors.New("telegram_payment_charge_id is required")
	}
	return nil
}

func (b *Bot) businessGiftBool(ctx context.Context, method string, params any, validate func() error) (bool, error) {
	if err := validate(); err != nil {
		return false, err
	}
	var result bool
	if err := b.call(ctx, method, params, &result); err != nil {
		return false, err
	}
	return result, nil
}

func validateTextFormatting(parseMode string, entities []telegram.MessageEntity) error {
	if parseMode != "" && len(entities) > 0 {
		return stderrors.New("text_parse_mode and text_entities cannot be used together")
	}
	return nil
}

func validateGiftListLimit(limit int) error {
	if limit < 0 {
		return stderrors.New("limit must not be negative")
	}
	if limit > 100 {
		return stderrors.New("limit must be between 1 and 100")
	}
	return nil
}

func validateBusinessOwnedGift(businessConnectionID, ownedGiftID string) error {
	if err := validateBusinessConnectionID(businessConnectionID); err != nil {
		return err
	}
	if strings.TrimSpace(ownedGiftID) == "" {
		return stderrors.New("owned_gift_id is required")
	}
	return nil
}
