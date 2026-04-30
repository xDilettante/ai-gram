package bot

import (
	"context"
	stderrors "errors"
	"fmt"
	"strings"

	"github.com/xDilettante/ai-gram/telegram"
)

// SendInvoiceParams contains supported parameters for sendInvoice.
type SendInvoiceParams struct {
	ChatID                    ChatID                            `json:"chat_id"`
	MessageThreadID           int64                             `json:"message_thread_id,omitempty"`
	DirectMessagesTopicID     int64                             `json:"direct_messages_topic_id,omitempty"`
	Title                     string                            `json:"title"`
	Description               string                            `json:"description"`
	Payload                   string                            `json:"payload"`
	ProviderToken             string                            `json:"provider_token,omitempty"`
	Currency                  string                            `json:"currency"`
	Prices                    []telegram.LabeledPrice           `json:"prices"`
	MaxTipAmount              int64                             `json:"max_tip_amount,omitempty"`
	SuggestedTipAmounts       []int64                           `json:"suggested_tip_amounts,omitempty"`
	StartParameter            string                            `json:"start_parameter,omitempty"`
	ProviderData              string                            `json:"provider_data,omitempty"`
	PhotoURL                  string                            `json:"photo_url,omitempty"`
	PhotoSize                 int64                             `json:"photo_size,omitempty"`
	PhotoWidth                int                               `json:"photo_width,omitempty"`
	PhotoHeight               int                               `json:"photo_height,omitempty"`
	NeedName                  bool                              `json:"need_name,omitempty"`
	NeedPhoneNumber           bool                              `json:"need_phone_number,omitempty"`
	NeedEmail                 bool                              `json:"need_email,omitempty"`
	NeedShippingAddress       bool                              `json:"need_shipping_address,omitempty"`
	SendPhoneNumberToProvider bool                              `json:"send_phone_number_to_provider,omitempty"`
	SendEmailToProvider       bool                              `json:"send_email_to_provider,omitempty"`
	IsFlexible                bool                              `json:"is_flexible,omitempty"`
	DisableNotification       bool                              `json:"disable_notification,omitempty"`
	ProtectContent            bool                              `json:"protect_content,omitempty"`
	AllowPaidBroadcast        bool                              `json:"allow_paid_broadcast,omitempty"`
	MessageEffectID           string                            `json:"message_effect_id,omitempty"`
	SuggestedPostParameters   *telegram.SuggestedPostParameters `json:"suggested_post_parameters,omitempty"`
	ReplyParameters           *telegram.ReplyParameters         `json:"reply_parameters,omitempty"`
	ReplyMarkup               *telegram.InlineKeyboardMarkup    `json:"reply_markup,omitempty"`
}

// CreateInvoiceLinkParams contains supported parameters for createInvoiceLink.
type CreateInvoiceLinkParams struct {
	BusinessConnectionID      string                  `json:"business_connection_id,omitempty"`
	Title                     string                  `json:"title"`
	Description               string                  `json:"description"`
	Payload                   string                  `json:"payload"`
	ProviderToken             string                  `json:"provider_token,omitempty"`
	Currency                  string                  `json:"currency"`
	Prices                    []telegram.LabeledPrice `json:"prices"`
	SubscriptionPeriod        int64                   `json:"subscription_period,omitempty"`
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

// AnswerShippingQueryParams contains supported parameters for answerShippingQuery.
type AnswerShippingQueryParams struct {
	ShippingQueryID string                    `json:"shipping_query_id"`
	OK              bool                      `json:"ok"`
	ShippingOptions []telegram.ShippingOption `json:"shipping_options,omitempty"`
	ErrorMessage    string                    `json:"error_message,omitempty"`
}

// AnswerPreCheckoutQueryParams contains supported parameters for answerPreCheckoutQuery.
type AnswerPreCheckoutQueryParams struct {
	PreCheckoutQueryID string `json:"pre_checkout_query_id"`
	OK                 bool   `json:"ok"`
	ErrorMessage       string `json:"error_message,omitempty"`
}

// SendInvoice sends an invoice message.
func (b *Bot) SendInvoice(ctx context.Context, params SendInvoiceParams) (*telegram.Message, error) {
	if err := params.validate(); err != nil {
		return nil, err
	}

	var message telegram.Message
	if err := b.call(ctx, "sendInvoice", params, &message); err != nil {
		return nil, err
	}

	return &message, nil
}

// CreateInvoiceLink creates a link for an invoice.
func (b *Bot) CreateInvoiceLink(ctx context.Context, params CreateInvoiceLinkParams) (string, error) {
	if err := params.validate(); err != nil {
		return "", err
	}

	var link string
	if err := b.call(ctx, "createInvoiceLink", params, &link); err != nil {
		return "", err
	}

	return link, nil
}

// AnswerShippingQuery answers a shipping query for a flexible invoice.
func (b *Bot) AnswerShippingQuery(ctx context.Context, params AnswerShippingQueryParams) (bool, error) {
	if err := params.validate(); err != nil {
		return false, err
	}

	var result bool
	if err := b.call(ctx, "answerShippingQuery", params, &result); err != nil {
		return false, err
	}

	return result, nil
}

// AnswerPreCheckoutQuery answers a pre-checkout query.
func (b *Bot) AnswerPreCheckoutQuery(ctx context.Context, params AnswerPreCheckoutQueryParams) (bool, error) {
	if err := params.validate(); err != nil {
		return false, err
	}

	var result bool
	if err := b.call(ctx, "answerPreCheckoutQuery", params, &result); err != nil {
		return false, err
	}

	return result, nil
}

func (params SendInvoiceParams) validate() error {
	if !params.ChatID.valid() {
		return stderrors.New("chat_id is required")
	}
	if err := validateMessageThreadID(params.MessageThreadID); err != nil {
		return err
	}
	if params.DirectMessagesTopicID < 0 {
		return stderrors.New("direct_messages_topic_id must not be negative")
	}
	if err := validateInvoiceCore(params.Title, params.Description, params.Payload, params.Currency, params.Prices); err != nil {
		return err
	}
	if err := validateInvoiceOptions(params.MaxTipAmount, params.SuggestedTipAmounts, params.PhotoURL, params.PhotoSize, params.PhotoWidth, params.PhotoHeight); err != nil {
		return err
	}
	if err := validateSuggestedPostParameters(params.SuggestedPostParameters); err != nil {
		return err
	}
	if err := validateReplyParameters(params.ReplyParameters); err != nil {
		return err
	}
	if params.ReplyMarkup != nil {
		if err := telegram.ValidateReplyMarkup(*params.ReplyMarkup); err != nil {
			return err
		}
	}
	return nil
}

func (params CreateInvoiceLinkParams) validate() error {
	if err := validateInvoiceCore(params.Title, params.Description, params.Payload, params.Currency, params.Prices); err != nil {
		return err
	}
	if params.SubscriptionPeriod < 0 {
		return stderrors.New("subscription_period must not be negative")
	}
	return validateInvoiceOptions(params.MaxTipAmount, params.SuggestedTipAmounts, params.PhotoURL, params.PhotoSize, params.PhotoWidth, params.PhotoHeight)
}

func (params AnswerShippingQueryParams) validate() error {
	if strings.TrimSpace(params.ShippingQueryID) == "" {
		return stderrors.New("shipping_query_id is required")
	}
	if params.OK {
		if strings.TrimSpace(params.ErrorMessage) != "" {
			return stderrors.New("error_message must be empty when ok is true")
		}
		if len(params.ShippingOptions) == 0 {
			return stderrors.New("shipping_options must not be empty when ok is true")
		}
		return validateShippingOptions(params.ShippingOptions)
	}
	if strings.TrimSpace(params.ErrorMessage) == "" {
		return stderrors.New("error_message is required when ok is false")
	}
	if len(params.ShippingOptions) != 0 {
		return stderrors.New("shipping_options must be empty when ok is false")
	}
	return nil
}

func (params AnswerPreCheckoutQueryParams) validate() error {
	if strings.TrimSpace(params.PreCheckoutQueryID) == "" {
		return stderrors.New("pre_checkout_query_id is required")
	}
	if params.OK && strings.TrimSpace(params.ErrorMessage) != "" {
		return stderrors.New("error_message must be empty when ok is true")
	}
	if !params.OK && strings.TrimSpace(params.ErrorMessage) == "" {
		return stderrors.New("error_message is required when ok is false")
	}
	return nil
}

func validateInvoiceCore(title string, description string, payload string, currency string, prices []telegram.LabeledPrice) error {
	if strings.TrimSpace(title) == "" {
		return stderrors.New("invoice title is required")
	}
	if strings.TrimSpace(description) == "" {
		return stderrors.New("invoice description is required")
	}
	if strings.TrimSpace(payload) == "" {
		return stderrors.New("invoice payload is required")
	}
	if strings.TrimSpace(currency) == "" {
		return stderrors.New("invoice currency is required")
	}
	return validateLabeledPrices(prices, "prices")
}

func validateInvoiceOptions(maxTipAmount int64, suggestedTipAmounts []int64, photoURL string, photoSize int64, photoWidth int, photoHeight int) error {
	if maxTipAmount < 0 {
		return stderrors.New("max_tip_amount must not be negative")
	}
	for index, amount := range suggestedTipAmounts {
		if amount < 0 {
			return fmt.Errorf("suggested_tip_amounts[%d] must not be negative", index)
		}
	}
	if photoURL != "" {
		if err := validateInlineHTTPURL(photoURL, "photo_url"); err != nil {
			return err
		}
	}
	if photoSize < 0 {
		return stderrors.New("photo_size must not be negative")
	}
	if photoWidth < 0 {
		return stderrors.New("photo_width must not be negative")
	}
	if photoHeight < 0 {
		return stderrors.New("photo_height must not be negative")
	}
	return nil
}

func validateLabeledPrices(prices []telegram.LabeledPrice, field string) error {
	if len(prices) == 0 {
		return fmt.Errorf("%s must not be empty", field)
	}
	for index, price := range prices {
		if strings.TrimSpace(price.Label) == "" {
			return fmt.Errorf("%s[%d].label is required", field, index)
		}
		if price.Amount < 0 {
			return fmt.Errorf("%s[%d].amount must not be negative", field, index)
		}
	}
	return nil
}

func validateShippingOptions(options []telegram.ShippingOption) error {
	for index, option := range options {
		if strings.TrimSpace(option.ID) == "" {
			return fmt.Errorf("shipping_options[%d].id is required", index)
		}
		if strings.TrimSpace(option.Title) == "" {
			return fmt.Errorf("shipping_options[%d].title is required", index)
		}
		if err := validateLabeledPrices(option.Prices, fmt.Sprintf("shipping_options[%d].prices", index)); err != nil {
			return err
		}
	}
	return nil
}

func validateSuggestedPostParameters(params *telegram.SuggestedPostParameters) error {
	if params == nil {
		return nil
	}
	if params.SendDate < 0 {
		return stderrors.New("suggested_post_parameters.send_date must not be negative")
	}
	if params.Price != nil {
		if strings.TrimSpace(params.Price.Currency) == "" {
			return stderrors.New("suggested_post_parameters.price.currency is required")
		}
		if params.Price.Amount < 0 {
			return stderrors.New("suggested_post_parameters.price.amount must not be negative")
		}
	}
	return nil
}
