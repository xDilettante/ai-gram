package telegram

import (
	"encoding/json"
	"testing"
)

func TestStarTransactionsDecoding(t *testing.T) {
	var transactions StarTransactions
	if err := json.Unmarshal([]byte(`{
		"transactions": [
			{
				"id": "tx1",
				"amount": 5,
				"nanostar_amount": 7,
				"date": 100,
				"source": {
					"type": "user",
					"transaction_type": "paid_media_payment",
					"user": {"id": 7, "is_bot": false, "first_name": "Alice"},
					"affiliate": {"affiliate_user": {"id": 8, "is_bot": true, "first_name": "Affiliate"}, "commission_per_mille": 100, "amount": 1, "nanostar_amount": 2},
					"paid_media": [{"type":"preview","width":10}],
					"paid_media_payload": "payload"
				}
			},
			{
				"id": "tx2",
				"amount": -2,
				"date": 101,
				"receiver": {"type":"fragment","withdrawal_state":{"type":"succeeded","date":102,"url":"https://fragment.com/tx"}}
			},
			{
				"id": "tx3",
				"amount": -1,
				"date": 103,
				"receiver": {"type":"telegram_api","request_count":3}
			}
		]
	}`), &transactions); err != nil {
		t.Fatalf("decode transactions: %v", err)
	}
	if len(transactions.Transactions) != 3 || transactions.Transactions[0].NanostarAmount != 7 {
		t.Fatalf("unexpected transactions: %+v", transactions)
	}
	user, ok := transactions.Transactions[0].Source.(TransactionPartnerUser)
	if !ok || user.TransactionType != "paid_media_payment" || len(user.PaidMedia) != 1 || user.Affiliate == nil {
		t.Fatalf("unexpected user source: %#v", transactions.Transactions[0].Source)
	}
	fragment, ok := transactions.Transactions[1].Receiver.(TransactionPartnerFragment)
	if !ok {
		t.Fatalf("unexpected fragment receiver: %#v", transactions.Transactions[1].Receiver)
	}
	if state, ok := fragment.WithdrawalState.(RevenueWithdrawalStateSucceeded); !ok || state.Date != 102 {
		t.Fatalf("unexpected withdrawal state: %#v", fragment.WithdrawalState)
	}
	telegramAPI, ok := transactions.Transactions[2].Receiver.(TransactionPartnerTelegramAPI)
	if !ok || telegramAPI.RequestCount != 3 {
		t.Fatalf("unexpected telegram_api receiver: %#v", transactions.Transactions[2].Receiver)
	}
}

func TestStarTransactionUnknownPartnerReturnsError(t *testing.T) {
	var transaction StarTransaction
	if err := json.Unmarshal([]byte(`{"id":"tx","amount":1,"date":1,"source":{"type":"unknown"}}`), &transaction); err == nil {
		t.Fatal("expected unsupported transaction partner error")
	}
}

func TestRevenueWithdrawalStateUnknownTypeReturnsError(t *testing.T) {
	if _, err := UnmarshalRevenueWithdrawalState([]byte(`{"type":"unknown"}`)); err == nil {
		t.Fatal("expected unsupported withdrawal state error")
	}
}
