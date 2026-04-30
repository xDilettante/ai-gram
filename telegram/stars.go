package telegram

import (
	"bytes"
	"encoding/json"
	stderrors "errors"
	"fmt"
)

const (
	transactionPartnerUserType             = "user"
	transactionPartnerChatType             = "chat"
	transactionPartnerAffiliateProgramType = "affiliate_program"
	transactionPartnerFragmentType         = "fragment"
	transactionPartnerTelegramAdsType      = "telegram_ads"
	transactionPartnerTelegramAPIType      = "telegram_api"
	transactionPartnerOtherType            = "other"

	revenueWithdrawalStatePendingType   = "pending"
	revenueWithdrawalStateSucceededType = "succeeded"
	revenueWithdrawalStateFailedType    = "failed"
)

// StarAmount describes an amount of Telegram Stars.
type StarAmount struct {
	Amount         int `json:"amount"`
	NanostarAmount int `json:"nanostar_amount,omitempty"`
}

// StarTransactions contains a list of Telegram Star transactions.
type StarTransactions struct {
	Transactions []StarTransaction `json:"transactions"`
}

// StarTransaction describes a Telegram Star transaction.
type StarTransaction struct {
	ID             string             `json:"id"`
	Amount         int                `json:"amount"`
	NanostarAmount int                `json:"nanostar_amount,omitempty"`
	Date           int64              `json:"date"`
	Source         TransactionPartner `json:"source,omitempty"`
	Receiver       TransactionPartner `json:"receiver,omitempty"`
}

// TransactionPartner marks Telegram Star transaction partner objects.
type TransactionPartner interface {
	transactionPartner()
}

// TransactionPartnerUser describes a transaction with a user.
type TransactionPartnerUser struct {
	Type                        string         `json:"type"`
	TransactionType             string         `json:"transaction_type"`
	User                        User           `json:"user"`
	Affiliate                   *AffiliateInfo `json:"affiliate,omitempty"`
	InvoicePayload              string         `json:"invoice_payload,omitempty"`
	SubscriptionPeriod          int            `json:"subscription_period,omitempty"`
	PaidMedia                   []PaidMedia    `json:"paid_media,omitempty"`
	PaidMediaPayload            string         `json:"paid_media_payload,omitempty"`
	PremiumSubscriptionDuration int            `json:"premium_subscription_duration,omitempty"`
}

// TransactionPartnerChat describes a transaction with a chat.
type TransactionPartnerChat struct {
	Type string `json:"type"`
	Chat Chat   `json:"chat"`
}

// TransactionPartnerAffiliateProgram describes an affiliate program transaction.
type TransactionPartnerAffiliateProgram struct {
	Type               string `json:"type"`
	SponsorUser        *User  `json:"sponsor_user,omitempty"`
	CommissionPerMille int    `json:"commission_per_mille"`
}

// TransactionPartnerFragment describes a withdrawal transaction with Fragment.
type TransactionPartnerFragment struct {
	Type            string                 `json:"type"`
	WithdrawalState RevenueWithdrawalState `json:"withdrawal_state,omitempty"`
}

// TransactionPartnerTelegramAds describes a withdrawal transaction to Telegram Ads.
type TransactionPartnerTelegramAds struct {
	Type string `json:"type"`
}

// TransactionPartnerTelegramAPI describes paid broadcast request charges.
type TransactionPartnerTelegramAPI struct {
	Type         string `json:"type"`
	RequestCount int    `json:"request_count"`
}

// TransactionPartnerOther describes an unknown transaction source or recipient.
type TransactionPartnerOther struct {
	Type string `json:"type"`
}

// AffiliateInfo contains affiliate commission details for a Star transaction.
type AffiliateInfo struct {
	AffiliateUser      *User `json:"affiliate_user,omitempty"`
	AffiliateChat      *Chat `json:"affiliate_chat,omitempty"`
	CommissionPerMille int   `json:"commission_per_mille"`
	Amount             int   `json:"amount"`
	NanostarAmount     int   `json:"nanostar_amount,omitempty"`
}

// RevenueWithdrawalState marks withdrawal state objects.
type RevenueWithdrawalState interface {
	revenueWithdrawalState()
}

// RevenueWithdrawalStatePending means the withdrawal is in progress.
type RevenueWithdrawalStatePending struct {
	Type string `json:"type"`
}

// RevenueWithdrawalStateSucceeded means the withdrawal succeeded.
type RevenueWithdrawalStateSucceeded struct {
	Type string `json:"type"`
	Date int64  `json:"date"`
	URL  string `json:"url"`
}

// RevenueWithdrawalStateFailed means the withdrawal failed and was refunded.
type RevenueWithdrawalStateFailed struct {
	Type string `json:"type"`
}

func (TransactionPartnerUser) transactionPartner()             {}
func (TransactionPartnerChat) transactionPartner()             {}
func (TransactionPartnerAffiliateProgram) transactionPartner() {}
func (TransactionPartnerFragment) transactionPartner()         {}
func (TransactionPartnerTelegramAds) transactionPartner()      {}
func (TransactionPartnerTelegramAPI) transactionPartner()      {}
func (TransactionPartnerOther) transactionPartner()            {}

func (RevenueWithdrawalStatePending) revenueWithdrawalState()   {}
func (RevenueWithdrawalStateSucceeded) revenueWithdrawalState() {}
func (RevenueWithdrawalStateFailed) revenueWithdrawalState()    {}

// UnmarshalTransactionPartner decodes a polymorphic TransactionPartner object.
func UnmarshalTransactionPartner(data []byte) (TransactionPartner, error) {
	var meta struct {
		Type string `json:"type"`
	}
	if err := json.Unmarshal(data, &meta); err != nil {
		return nil, err
	}

	switch meta.Type {
	case transactionPartnerUserType:
		var partner TransactionPartnerUser
		if err := json.Unmarshal(data, &partner); err != nil {
			return nil, err
		}
		partner.Type = transactionPartnerUserType
		return partner, nil
	case transactionPartnerChatType:
		var partner TransactionPartnerChat
		if err := json.Unmarshal(data, &partner); err != nil {
			return nil, err
		}
		partner.Type = transactionPartnerChatType
		return partner, nil
	case transactionPartnerAffiliateProgramType:
		var partner TransactionPartnerAffiliateProgram
		if err := json.Unmarshal(data, &partner); err != nil {
			return nil, err
		}
		partner.Type = transactionPartnerAffiliateProgramType
		return partner, nil
	case transactionPartnerFragmentType:
		return unmarshalTransactionPartnerFragment(data)
	case transactionPartnerTelegramAdsType:
		var partner TransactionPartnerTelegramAds
		if err := json.Unmarshal(data, &partner); err != nil {
			return nil, err
		}
		partner.Type = transactionPartnerTelegramAdsType
		return partner, nil
	case transactionPartnerTelegramAPIType:
		var partner TransactionPartnerTelegramAPI
		if err := json.Unmarshal(data, &partner); err != nil {
			return nil, err
		}
		partner.Type = transactionPartnerTelegramAPIType
		return partner, nil
	case transactionPartnerOtherType:
		var partner TransactionPartnerOther
		if err := json.Unmarshal(data, &partner); err != nil {
			return nil, err
		}
		partner.Type = transactionPartnerOtherType
		return partner, nil
	default:
		return nil, stderrors.New("unsupported transaction partner type")
	}
}

// UnmarshalRevenueWithdrawalState decodes a polymorphic RevenueWithdrawalState object.
func UnmarshalRevenueWithdrawalState(data []byte) (RevenueWithdrawalState, error) {
	var meta struct {
		Type string `json:"type"`
	}
	if err := json.Unmarshal(data, &meta); err != nil {
		return nil, err
	}

	switch meta.Type {
	case revenueWithdrawalStatePendingType:
		var state RevenueWithdrawalStatePending
		if err := json.Unmarshal(data, &state); err != nil {
			return nil, err
		}
		state.Type = revenueWithdrawalStatePendingType
		return state, nil
	case revenueWithdrawalStateSucceededType:
		var state RevenueWithdrawalStateSucceeded
		if err := json.Unmarshal(data, &state); err != nil {
			return nil, err
		}
		state.Type = revenueWithdrawalStateSucceededType
		return state, nil
	case revenueWithdrawalStateFailedType:
		var state RevenueWithdrawalStateFailed
		if err := json.Unmarshal(data, &state); err != nil {
			return nil, err
		}
		state.Type = revenueWithdrawalStateFailedType
		return state, nil
	default:
		return nil, stderrors.New("unsupported revenue withdrawal state type")
	}
}

// UnmarshalJSON decodes a StarTransaction with polymorphic transaction partners.
func (transaction *StarTransaction) UnmarshalJSON(data []byte) error {
	var payload struct {
		ID             string          `json:"id"`
		Amount         int             `json:"amount"`
		NanostarAmount int             `json:"nanostar_amount,omitempty"`
		Date           int64           `json:"date"`
		Source         json.RawMessage `json:"source,omitempty"`
		Receiver       json.RawMessage `json:"receiver,omitempty"`
	}
	if err := json.Unmarshal(data, &payload); err != nil {
		return err
	}

	source, err := optionalTransactionPartner(payload.Source)
	if err != nil {
		return fmt.Errorf("source is invalid: %w", err)
	}
	receiver, err := optionalTransactionPartner(payload.Receiver)
	if err != nil {
		return fmt.Errorf("receiver is invalid: %w", err)
	}

	transaction.ID = payload.ID
	transaction.Amount = payload.Amount
	transaction.NanostarAmount = payload.NanostarAmount
	transaction.Date = payload.Date
	transaction.Source = source
	transaction.Receiver = receiver
	return nil
}

// UnmarshalJSON decodes TransactionPartnerUser with polymorphic paid media entries.
func (partner *TransactionPartnerUser) UnmarshalJSON(data []byte) error {
	var payload struct {
		Type                        string            `json:"type"`
		TransactionType             string            `json:"transaction_type"`
		User                        User              `json:"user"`
		Affiliate                   *AffiliateInfo    `json:"affiliate,omitempty"`
		InvoicePayload              string            `json:"invoice_payload,omitempty"`
		SubscriptionPeriod          int               `json:"subscription_period,omitempty"`
		PaidMedia                   []json.RawMessage `json:"paid_media,omitempty"`
		PaidMediaPayload            string            `json:"paid_media_payload,omitempty"`
		PremiumSubscriptionDuration int               `json:"premium_subscription_duration,omitempty"`
	}
	if err := json.Unmarshal(data, &payload); err != nil {
		return err
	}
	paidMedia, err := unmarshalPaidMediaItems(payload.PaidMedia)
	if err != nil {
		return err
	}

	partner.Type = transactionPartnerUserType
	partner.TransactionType = payload.TransactionType
	partner.User = payload.User
	partner.Affiliate = payload.Affiliate
	partner.InvoicePayload = payload.InvoicePayload
	partner.SubscriptionPeriod = payload.SubscriptionPeriod
	partner.PaidMedia = paidMedia
	partner.PaidMediaPayload = payload.PaidMediaPayload
	partner.PremiumSubscriptionDuration = payload.PremiumSubscriptionDuration
	return nil
}

func unmarshalTransactionPartnerFragment(data []byte) (TransactionPartner, error) {
	var payload struct {
		Type            string          `json:"type"`
		WithdrawalState json.RawMessage `json:"withdrawal_state,omitempty"`
	}
	if err := json.Unmarshal(data, &payload); err != nil {
		return nil, err
	}
	state, err := optionalRevenueWithdrawalState(payload.WithdrawalState)
	if err != nil {
		return nil, err
	}
	return TransactionPartnerFragment{Type: transactionPartnerFragmentType, WithdrawalState: state}, nil
}

func optionalTransactionPartner(raw json.RawMessage) (TransactionPartner, error) {
	if len(raw) == 0 || bytes.Equal(raw, []byte("null")) {
		return nil, nil
	}
	partner, err := UnmarshalTransactionPartner(raw)
	if err != nil {
		return nil, err
	}
	return partner, nil
}

func optionalRevenueWithdrawalState(raw json.RawMessage) (RevenueWithdrawalState, error) {
	if len(raw) == 0 || bytes.Equal(raw, []byte("null")) {
		return nil, nil
	}
	state, err := UnmarshalRevenueWithdrawalState(raw)
	if err != nil {
		return nil, err
	}
	return state, nil
}
