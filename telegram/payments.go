package telegram

// LabeledPrice represents a price component in the smallest units of a currency.
type LabeledPrice struct {
	Label  string `json:"label"`
	Amount int64  `json:"amount"`
}

// Invoice contains basic information about an invoice message.
type Invoice struct {
	Title          string `json:"title"`
	Description    string `json:"description"`
	StartParameter string `json:"start_parameter,omitempty"`
	Currency       string `json:"currency"`
	TotalAmount    int64  `json:"total_amount"`
}

// ShippingAddress contains a user's shipping address.
type ShippingAddress struct {
	CountryCode string `json:"country_code"`
	State       string `json:"state"`
	City        string `json:"city"`
	StreetLine1 string `json:"street_line1"`
	StreetLine2 string `json:"street_line2"`
	PostCode    string `json:"post_code"`
}

// OrderInfo contains information about an order.
type OrderInfo struct {
	Name            string           `json:"name,omitempty"`
	PhoneNumber     string           `json:"phone_number,omitempty"`
	Email           string           `json:"email,omitempty"`
	ShippingAddress *ShippingAddress `json:"shipping_address,omitempty"`
}

// ShippingOption represents one shipping option for an invoice.
type ShippingOption struct {
	ID     string         `json:"id"`
	Title  string         `json:"title"`
	Prices []LabeledPrice `json:"prices"`
}

// SuccessfulPayment contains basic information about a successful payment.
type SuccessfulPayment struct {
	Currency                   string     `json:"currency"`
	TotalAmount                int64      `json:"total_amount"`
	InvoicePayload             string     `json:"invoice_payload"`
	SubscriptionExpirationDate int64      `json:"subscription_expiration_date,omitempty"`
	IsRecurring                bool       `json:"is_recurring,omitempty"`
	IsFirstRecurring           bool       `json:"is_first_recurring,omitempty"`
	ShippingOptionID           string     `json:"shipping_option_id,omitempty"`
	OrderInfo                  *OrderInfo `json:"order_info,omitempty"`
	TelegramPaymentChargeID    string     `json:"telegram_payment_charge_id"`
	ProviderPaymentChargeID    string     `json:"provider_payment_charge_id"`
}

// RefundedPayment contains basic information about a refunded payment.
type RefundedPayment struct {
	Currency                string `json:"currency"`
	TotalAmount             int64  `json:"total_amount"`
	InvoicePayload          string `json:"invoice_payload"`
	TelegramPaymentChargeID string `json:"telegram_payment_charge_id"`
	ProviderPaymentChargeID string `json:"provider_payment_charge_id,omitempty"`
}

// ShippingQuery represents an incoming shipping query for a flexible invoice.
type ShippingQuery struct {
	ID              string          `json:"id"`
	From            User            `json:"from"`
	InvoicePayload  string          `json:"invoice_payload"`
	ShippingAddress ShippingAddress `json:"shipping_address"`
}

// PreCheckoutQuery represents an incoming pre-checkout query.
type PreCheckoutQuery struct {
	ID               string     `json:"id"`
	From             User       `json:"from"`
	Currency         string     `json:"currency"`
	TotalAmount      int64      `json:"total_amount"`
	InvoicePayload   string     `json:"invoice_payload"`
	ShippingOptionID string     `json:"shipping_option_id,omitempty"`
	OrderInfo        *OrderInfo `json:"order_info,omitempty"`
}

// SuggestedPostPrice describes the amount a channel direct message sender pays for a suggested post.
type SuggestedPostPrice struct {
	Currency string `json:"currency"`
	Amount   int64  `json:"amount"`
}

// SuggestedPostParameters describes suggested post metadata for direct messages.
type SuggestedPostParameters struct {
	Price    *SuggestedPostPrice `json:"price,omitempty"`
	SendDate int64               `json:"send_date,omitempty"`
}

// SuggestedPostApprovalFailed describes a service message about a failed suggested post approval.
type SuggestedPostApprovalFailed struct {
	SuggestedPostMessage *Message           `json:"suggested_post_message,omitempty"`
	Price                SuggestedPostPrice `json:"price"`
}

// SuggestedPostApproved describes a service message about a suggested post approval.
type SuggestedPostApproved struct {
	SuggestedPostMessage *Message            `json:"suggested_post_message,omitempty"`
	Price                *SuggestedPostPrice `json:"price,omitempty"`
	SendDate             int64               `json:"send_date"`
}

// SuggestedPostDeclined describes a service message about a suggested post decline.
type SuggestedPostDeclined struct {
	SuggestedPostMessage *Message `json:"suggested_post_message,omitempty"`
	Comment              string   `json:"comment,omitempty"`
}

// SuggestedPostPaid describes a service message about a suggested post payment.
type SuggestedPostPaid struct {
	SuggestedPostMessage *Message    `json:"suggested_post_message,omitempty"`
	Currency             string      `json:"currency"`
	Amount               int64       `json:"amount,omitempty"`
	StarAmount           *StarAmount `json:"star_amount,omitempty"`
}

// SuggestedPostRefunded describes a service message about a suggested post refund.
type SuggestedPostRefunded struct {
	SuggestedPostMessage *Message `json:"suggested_post_message,omitempty"`
	Reason               string   `json:"reason"`
}
