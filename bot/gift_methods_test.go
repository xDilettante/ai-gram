package bot

import (
	"context"
	"encoding/json"
	stderrors "errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	apierrors "github.com/xDilettante/ai-gram/errors"
	"github.com/xDilettante/ai-gram/telegram"
)

func TestGiftAndStarMethodsSuccess(t *testing.T) {
	const token = "123:secret"
	seen := map[string]map[string]any{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		method := strings.TrimPrefix(r.URL.Path, "/bot"+token+"/")
		var payload map[string]any
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode %s payload: %v", method, err)
		}
		seen[method] = payload
		w.Header().Set("Content-Type", "application/json")
		switch method {
		case "getAvailableGifts":
			_, _ = w.Write([]byte(`{"ok":true,"result":{"gifts":[{"id":"gift-1","sticker":{"file_id":"s","file_unique_id":"u","type":"regular","width":1,"height":1,"is_animated":false,"is_video":false},"star_count":10}]}}`))
		case "getBusinessAccountStarBalance", "getMyStarBalance":
			_, _ = w.Write([]byte(`{"ok":true,"result":{"amount":42,"nanostar_amount":7}}`))
		case "getBusinessAccountGifts", "getUserGifts", "getChatGifts":
			_, _ = w.Write([]byte(`{"ok":true,"result":{"total_count":1,"gifts":[{"type":"regular","gift":{"id":"gift-1","sticker":{"file_id":"s","file_unique_id":"u","type":"regular","width":1,"height":1,"is_animated":false,"is_video":false},"star_count":10},"owned_gift_id":"owned-1","send_date":100}],"next_offset":"next"}}`))
		default:
			_, _ = w.Write([]byte(`{"ok":true,"result":true}`))
		}
	}))
	defer server.Close()

	bot := newTestBot(t, token, server.URL, server.Client())
	ctx := context.Background()

	gifts, err := bot.GetAvailableGifts(ctx, GetAvailableGiftsParams{})
	if err != nil || gifts == nil || len(gifts.Gifts) != 1 || gifts.Gifts[0].ID != "gift-1" {
		t.Fatalf("unexpected getAvailableGifts result: gifts=%+v err=%v", gifts, err)
	}
	if len(seen["getAvailableGifts"]) != 0 {
		t.Fatalf("unexpected getAvailableGifts payload: %#v", seen["getAvailableGifts"])
	}

	ok, err := bot.SendGift(ctx, SendGiftParams{UserID: 10, GiftID: "gift-1", PayForUpgrade: true, Text: "hello", TextParseMode: "Markdown"})
	if err != nil || !ok {
		t.Fatalf("SendGift failed: ok=%v err=%v", ok, err)
	}
	assertPayload(t, seen["sendGift"], map[string]any{"user_id": float64(10), "gift_id": "gift-1", "pay_for_upgrade": true, "text": "hello", "text_parse_mode": "Markdown"})
	if _, ok := seen["sendGift"]["chat_id"]; ok {
		t.Fatalf("chat_id should be omitted for user gift: %#v", seen["sendGift"])
	}

	ok, err = bot.GiftPremiumSubscription(ctx, GiftPremiumSubscriptionParams{UserID: 10, MonthCount: 3, StarCount: 1000, Text: "premium"})
	if err != nil || !ok {
		t.Fatalf("GiftPremiumSubscription failed: ok=%v err=%v", ok, err)
	}
	assertPayload(t, seen["giftPremiumSubscription"], map[string]any{"user_id": float64(10), "month_count": float64(3), "star_count": float64(1000), "text": "premium"})

	amount, err := bot.GetBusinessAccountStarBalance(ctx, GetBusinessAccountStarBalanceParams{BusinessConnectionID: "business-1"})
	if err != nil || amount == nil || amount.Amount != 42 || amount.NanostarAmount != 7 {
		t.Fatalf("GetBusinessAccountStarBalance failed: amount=%+v err=%v", amount, err)
	}
	assertPayload(t, seen["getBusinessAccountStarBalance"], map[string]any{"business_connection_id": "business-1"})

	ok, err = bot.TransferBusinessAccountStars(ctx, TransferBusinessAccountStarsParams{BusinessConnectionID: "business-1", StarCount: 100})
	if err != nil || !ok {
		t.Fatalf("TransferBusinessAccountStars failed: ok=%v err=%v", ok, err)
	}
	assertPayload(t, seen["transferBusinessAccountStars"], map[string]any{"business_connection_id": "business-1", "star_count": float64(100)})

	businessGifts, err := bot.GetBusinessAccountGifts(ctx, GetBusinessAccountGiftsParams{BusinessConnectionID: "business-1", ExcludeSaved: true, SortByPrice: true, Offset: "off", Limit: 10})
	if err != nil || businessGifts == nil || businessGifts.TotalCount != 1 || businessGifts.NextOffset != "next" {
		t.Fatalf("GetBusinessAccountGifts failed: gifts=%+v err=%v", businessGifts, err)
	}
	assertPayload(t, seen["getBusinessAccountGifts"], map[string]any{"business_connection_id": "business-1", "exclude_saved": true, "sort_by_price": true, "offset": "off", "limit": float64(10)})

	userGifts, err := bot.GetUserGifts(ctx, GetUserGiftsParams{UserID: 55, ExcludeUnique: true, Limit: 5})
	if err != nil || userGifts == nil || userGifts.TotalCount != 1 {
		t.Fatalf("GetUserGifts failed: gifts=%+v err=%v", userGifts, err)
	}
	assertPayload(t, seen["getUserGifts"], map[string]any{"user_id": float64(55), "exclude_unique": true, "limit": float64(5)})

	chatGifts, err := bot.GetChatGifts(ctx, GetChatGiftsParams{ChatID: ChatIDString("@channel"), ExcludeUnsaved: true, ExcludeFromBlockchain: true, Limit: 6})
	if err != nil || chatGifts == nil || chatGifts.TotalCount != 1 {
		t.Fatalf("GetChatGifts failed: gifts=%+v err=%v", chatGifts, err)
	}
	assertPayload(t, seen["getChatGifts"], map[string]any{"chat_id": "@channel", "exclude_unsaved": true, "exclude_from_blockchain": true, "limit": float64(6)})

	ok, err = bot.ConvertGiftToStars(ctx, ConvertGiftToStarsParams{BusinessConnectionID: "business-1", OwnedGiftID: "owned-1"})
	if err != nil || !ok {
		t.Fatalf("ConvertGiftToStars failed: ok=%v err=%v", ok, err)
	}
	assertPayload(t, seen["convertGiftToStars"], map[string]any{"business_connection_id": "business-1", "owned_gift_id": "owned-1"})

	ok, err = bot.UpgradeGift(ctx, UpgradeGiftParams{BusinessConnectionID: "business-1", OwnedGiftID: "owned-1", KeepOriginalDetails: true, StarCount: 50})
	if err != nil || !ok {
		t.Fatalf("UpgradeGift failed: ok=%v err=%v", ok, err)
	}
	assertPayload(t, seen["upgradeGift"], map[string]any{"business_connection_id": "business-1", "owned_gift_id": "owned-1", "keep_original_details": true, "star_count": float64(50)})

	ok, err = bot.TransferGift(ctx, TransferGiftParams{BusinessConnectionID: "business-1", OwnedGiftID: "owned-1", NewOwnerChatID: 123, StarCount: 25})
	if err != nil || !ok {
		t.Fatalf("TransferGift failed: ok=%v err=%v", ok, err)
	}
	assertPayload(t, seen["transferGift"], map[string]any{"business_connection_id": "business-1", "owned_gift_id": "owned-1", "new_owner_chat_id": float64(123), "star_count": float64(25)})

	botBalance, err := bot.GetMyStarBalance(ctx, GetMyStarBalanceParams{})
	if err != nil || botBalance == nil || botBalance.Amount != 42 {
		t.Fatalf("GetMyStarBalance failed: amount=%+v err=%v", botBalance, err)
	}
	if len(seen["getMyStarBalance"]) != 0 {
		t.Fatalf("unexpected getMyStarBalance payload: %#v", seen["getMyStarBalance"])
	}

	ok, err = bot.EditUserStarSubscription(ctx, EditUserStarSubscriptionParams{UserID: 10, TelegramPaymentChargeID: "charge-1", IsCanceled: true})
	if err != nil || !ok {
		t.Fatalf("EditUserStarSubscription failed: ok=%v err=%v", ok, err)
	}
	assertPayload(t, seen["editUserStarSubscription"], map[string]any{"user_id": float64(10), "telegram_payment_charge_id": "charge-1", "is_canceled": true})
}

func TestSendGiftSendsChatTarget(t *testing.T) {
	const token = "123:secret"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var payload map[string]any
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode payload: %v", err)
		}
		assertPayload(t, payload, map[string]any{"chat_id": "@channel", "gift_id": "gift-1"})
		if _, ok := payload["user_id"]; ok {
			t.Fatalf("user_id should be omitted for chat gift: %#v", payload)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true,"result":true}`))
	}))
	defer server.Close()

	bot := newTestBot(t, token, server.URL, server.Client())
	ok, err := bot.SendGift(context.Background(), SendGiftParams{ChatID: ChatIDString("@channel"), GiftID: "gift-1"})
	if err != nil || !ok {
		t.Fatalf("SendGift failed: ok=%v err=%v", ok, err)
	}
}

func TestGiftAndStarValidation(t *testing.T) {
	const token = "123:secret"
	bot := newTestBot(t, token, "https://example.test", nil)
	ctx := context.Background()
	tests := []struct {
		name string
		call func() error
		want string
	}{
		{"send gift missing target", func() error { _, err := bot.SendGift(ctx, SendGiftParams{GiftID: "gift"}); return err }, "exactly one"},
		{"send gift target conflict", func() error {
			_, err := bot.SendGift(ctx, SendGiftParams{UserID: 1, ChatID: ChatIDInt(2), GiftID: "gift"})
			return err
		}, "exactly one"},
		{"send gift missing gift", func() error { _, err := bot.SendGift(ctx, SendGiftParams{UserID: 1}); return err }, "gift_id"},
		{"send gift text conflict", func() error {
			_, err := bot.SendGift(ctx, SendGiftParams{UserID: 1, GiftID: "gift", TextParseMode: "HTML", TextEntities: []telegram.MessageEntity{{Type: "bold"}}})
			return err
		}, "text_parse_mode"},
		{"premium invalid month", func() error {
			_, err := bot.GiftPremiumSubscription(ctx, GiftPremiumSubscriptionParams{UserID: 1, MonthCount: 1, StarCount: 1})
			return err
		}, "month_count"},
		{"premium invalid stars", func() error {
			_, err := bot.GiftPremiumSubscription(ctx, GiftPremiumSubscriptionParams{UserID: 1, MonthCount: 3, StarCount: 999})
			return err
		}, "star_count"},
		{"business balance missing id", func() error {
			_, err := bot.GetBusinessAccountStarBalance(ctx, GetBusinessAccountStarBalanceParams{})
			return err
		}, "business_connection_id"},
		{"business transfer invalid count", func() error {
			_, err := bot.TransferBusinessAccountStars(ctx, TransferBusinessAccountStarsParams{BusinessConnectionID: "business", StarCount: 0})
			return err
		}, "star_count"},
		{"business gifts invalid limit", func() error {
			_, err := bot.GetBusinessAccountGifts(ctx, GetBusinessAccountGiftsParams{BusinessConnectionID: "business", Limit: 101})
			return err
		}, "limit"},
		{"user gifts invalid user", func() error { _, err := bot.GetUserGifts(ctx, GetUserGiftsParams{}); return err }, "user_id"},
		{"chat gifts invalid chat", func() error { _, err := bot.GetChatGifts(ctx, GetChatGiftsParams{}); return err }, "chat_id"},
		{"convert missing owned gift", func() error {
			_, err := bot.ConvertGiftToStars(ctx, ConvertGiftToStarsParams{BusinessConnectionID: "business"})
			return err
		}, "owned_gift_id"},
		{"upgrade negative stars", func() error {
			_, err := bot.UpgradeGift(ctx, UpgradeGiftParams{BusinessConnectionID: "business", OwnedGiftID: "owned", StarCount: -1})
			return err
		}, "star_count"},
		{"transfer invalid new owner", func() error {
			_, err := bot.TransferGift(ctx, TransferGiftParams{BusinessConnectionID: "business", OwnedGiftID: "owned"})
			return err
		}, "new_owner_chat_id"},
		{"subscription missing charge", func() error {
			_, err := bot.EditUserStarSubscription(ctx, EditUserStarSubscriptionParams{UserID: 1})
			return err
		}, "telegram_payment_charge_id"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.call()
			if err == nil || !strings.Contains(err.Error(), tt.want) {
				t.Fatalf("expected %q error, got %v", tt.want, err)
			}
			assertNoToken(t, err, token)
		})
	}
}

func TestGiftMethodsReturnAPIError(t *testing.T) {
	const token = "123:secret"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":false,"error_code":400,"description":"Bad Request"}`))
	}))
	defer server.Close()

	bot := newTestBot(t, token, server.URL, server.Client())
	ctx := context.Background()
	tests := []struct {
		name string
		call func() error
	}{
		{"GetAvailableGifts", func() error { _, err := bot.GetAvailableGifts(ctx, GetAvailableGiftsParams{}); return err }},
		{"SendGift", func() error { _, err := bot.SendGift(ctx, SendGiftParams{UserID: 10, GiftID: "gift-1"}); return err }},
		{"GiftPremiumSubscription", func() error {
			_, err := bot.GiftPremiumSubscription(ctx, GiftPremiumSubscriptionParams{UserID: 10, MonthCount: 3, StarCount: 1000})
			return err
		}},
		{"GetBusinessAccountStarBalance", func() error {
			_, err := bot.GetBusinessAccountStarBalance(ctx, GetBusinessAccountStarBalanceParams{BusinessConnectionID: "business-1"})
			return err
		}},
		{"TransferBusinessAccountStars", func() error {
			_, err := bot.TransferBusinessAccountStars(ctx, TransferBusinessAccountStarsParams{BusinessConnectionID: "business-1", StarCount: 1})
			return err
		}},
		{"GetBusinessAccountGifts", func() error {
			_, err := bot.GetBusinessAccountGifts(ctx, GetBusinessAccountGiftsParams{BusinessConnectionID: "business-1"})
			return err
		}},
		{"GetUserGifts", func() error { _, err := bot.GetUserGifts(ctx, GetUserGiftsParams{UserID: 10}); return err }},
		{"GetChatGifts", func() error { _, err := bot.GetChatGifts(ctx, GetChatGiftsParams{ChatID: ChatIDInt(123)}); return err }},
		{"ConvertGiftToStars", func() error {
			_, err := bot.ConvertGiftToStars(ctx, ConvertGiftToStarsParams{BusinessConnectionID: "business-1", OwnedGiftID: "owned-1"})
			return err
		}},
		{"UpgradeGift", func() error {
			_, err := bot.UpgradeGift(ctx, UpgradeGiftParams{BusinessConnectionID: "business-1", OwnedGiftID: "owned-1"})
			return err
		}},
		{"TransferGift", func() error {
			_, err := bot.TransferGift(ctx, TransferGiftParams{BusinessConnectionID: "business-1", OwnedGiftID: "owned-1", NewOwnerChatID: 123})
			return err
		}},
		{"GetMyStarBalance", func() error { _, err := bot.GetMyStarBalance(ctx, GetMyStarBalanceParams{}); return err }},
		{"EditUserStarSubscription", func() error {
			_, err := bot.EditUserStarSubscription(ctx, EditUserStarSubscriptionParams{UserID: 10, TelegramPaymentChargeID: "charge-1"})
			return err
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.call()
			if err == nil {
				t.Fatal("expected API error")
			}
			var apiErr *apierrors.APIError
			if !stderrors.As(err, &apiErr) {
				t.Fatalf("expected APIError, got %T", err)
			}
			assertNoToken(t, err, token)
			if strings.Contains(err.Error(), "gift-1") || strings.Contains(err.Error(), "owned-1") || strings.Contains(err.Error(), "charge-1") {
				t.Fatalf("error leaked sensitive gift/payment identifier: %q", err.Error())
			}
		})
	}
}

func TestGiftMethodsHandleInvalidJSONHTTPAndCancelledContext(t *testing.T) {
	const token = "123:secret"

	invalidJSON := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true,"result":`))
	}))
	defer invalidJSON.Close()
	botInvalid := newTestBot(t, token, invalidJSON.URL, invalidJSON.Client())
	if _, err := botInvalid.GetAvailableGifts(context.Background(), GetAvailableGiftsParams{}); err == nil {
		t.Fatal("expected invalid JSON error")
	} else {
		assertNoToken(t, err, token)
	}

	http500 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "server error", http.StatusInternalServerError)
	}))
	defer http500.Close()
	botHTTP := newTestBot(t, token, http500.URL, http500.Client())
	if _, err := botHTTP.GetMyStarBalance(context.Background(), GetMyStarBalanceParams{}); err == nil {
		t.Fatal("expected HTTP 500 error")
	} else {
		assertNoToken(t, err, token)
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	botCanceled := newTestBot(t, token, "https://example.test", nil)
	if _, err := botCanceled.GetUserGifts(ctx, GetUserGiftsParams{UserID: 1}); err == nil {
		t.Fatal("expected cancelled context error")
	} else {
		assertNoToken(t, err, token)
	}
}

func assertPayload(t *testing.T, payload map[string]any, want map[string]any) {
	t.Helper()
	for key, value := range want {
		if got := payload[key]; got != value {
			t.Fatalf("payload[%s] = %#v, want %#v (payload=%#v)", key, got, value, payload)
		}
	}
}
