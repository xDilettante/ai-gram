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

func TestSendInvoiceSendsPayloadAndDecodesMessage(t *testing.T) {
	const token = "123:secret"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/bot"+token+"/sendInvoice" {
			t.Fatalf("unexpected path: %q", r.URL.Path)
		}
		var payload map[string]any
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode payload: %v", err)
		}
		if payload["chat_id"] != float64(12345) || payload["message_thread_id"] != float64(7) || payload["direct_messages_topic_id"] != float64(8) {
			t.Fatalf("unexpected chat payload: %#v", payload)
		}
		if payload["title"] != "Title" || payload["description"] != "Description" || payload["payload"] != "safe-payload" || payload["currency"] != "XTR" {
			t.Fatalf("unexpected invoice payload: %#v", payload)
		}
		if payload["provider_token"] != "test-provider-token" || payload["start_parameter"] != "start" || payload["provider_data"] != "{}" {
			t.Fatalf("unexpected provider payload: %#v", payload)
		}
		prices, ok := payload["prices"].([]any)
		if !ok || len(prices) != 1 {
			t.Fatalf("unexpected prices: %#v", payload["prices"])
		}
		price, _ := prices[0].(map[string]any)
		if price["label"] != "Price" || price["amount"] != float64(150) {
			t.Fatalf("unexpected price payload: %#v", price)
		}
		tips, ok := payload["suggested_tip_amounts"].([]any)
		if !ok || len(tips) != 2 || tips[0] != float64(1) || tips[1] != float64(5) || payload["max_tip_amount"] != float64(10) {
			t.Fatalf("unexpected tip payload: %#v", payload)
		}
		if payload["photo_url"] != "https://example.test/invoice.png" || payload["photo_size"] != float64(10) || payload["photo_width"] != float64(20) || payload["photo_height"] != float64(30) {
			t.Fatalf("unexpected photo payload: %#v", payload)
		}
		if payload["need_name"] != true || payload["need_phone_number"] != true || payload["need_email"] != true || payload["need_shipping_address"] != true {
			t.Fatalf("unexpected required info payload: %#v", payload)
		}
		if payload["send_phone_number_to_provider"] != true || payload["send_email_to_provider"] != true || payload["is_flexible"] != true {
			t.Fatalf("unexpected provider flags: %#v", payload)
		}
		if payload["disable_notification"] != true || payload["protect_content"] != true || payload["allow_paid_broadcast"] != true || payload["message_effect_id"] != "effect" {
			t.Fatalf("unexpected send flags: %#v", payload)
		}
		if _, ok := payload["reply_parameters"].(map[string]any); !ok {
			t.Fatalf("missing reply_parameters: %#v", payload)
		}
		if _, ok := payload["reply_markup"].(map[string]any); !ok {
			t.Fatalf("missing reply_markup: %#v", payload)
		}
		if suggestedPost, ok := payload["suggested_post_parameters"].(map[string]any); !ok || suggestedPost["send_date"] != float64(123) {
			t.Fatalf("unexpected suggested_post_parameters: %#v", payload["suggested_post_parameters"])
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true,"result":{"message_id":55,"chat":{"id":12345,"type":"private"},"date":1,"invoice":{"title":"Title","description":"Description","currency":"XTR","total_amount":150}}}`))
	}))
	defer server.Close()

	markup := telegram.NewInlineKeyboard([]telegram.InlineKeyboardButton{telegram.InlineButtonCallback("Pay", "pay")})
	bot := newTestBot(t, token, server.URL, server.Client())
	message, err := bot.SendInvoice(context.Background(), SendInvoiceParams{
		ChatID:                    ChatIDInt(12345),
		MessageThreadID:           7,
		DirectMessagesTopicID:     8,
		Title:                     "Title",
		Description:               "Description",
		Payload:                   "safe-payload",
		ProviderToken:             "test-provider-token",
		Currency:                  "XTR",
		Prices:                    []telegram.LabeledPrice{{Label: "Price", Amount: 150}},
		MaxTipAmount:              10,
		SuggestedTipAmounts:       []int64{1, 5},
		StartParameter:            "start",
		ProviderData:              "{}",
		PhotoURL:                  "https://example.test/invoice.png",
		PhotoSize:                 10,
		PhotoWidth:                20,
		PhotoHeight:               30,
		NeedName:                  true,
		NeedPhoneNumber:           true,
		NeedEmail:                 true,
		NeedShippingAddress:       true,
		SendPhoneNumberToProvider: true,
		SendEmailToProvider:       true,
		IsFlexible:                true,
		DisableNotification:       true,
		ProtectContent:            true,
		AllowPaidBroadcast:        true,
		MessageEffectID:           "effect",
		SuggestedPostParameters:   &telegram.SuggestedPostParameters{Price: &telegram.SuggestedPostPrice{Currency: "XTR", Amount: 1}, SendDate: 123},
		ReplyParameters:           &telegram.ReplyParameters{MessageID: 1},
		ReplyMarkup:               &markup,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if message == nil || message.Invoice == nil || message.Invoice.TotalAmount != 150 {
		t.Fatalf("unexpected message: %+v", message)
	}
}

func TestSendInvoiceAllowsEmptyProviderToken(t *testing.T) {
	const token = "123:secret"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/bot"+token+"/sendInvoice" {
			t.Fatalf("unexpected path: %q", r.URL.Path)
		}
		var payload map[string]any
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode payload: %v", err)
		}
		if _, ok := payload["provider_token"]; ok {
			t.Fatalf("empty provider token should be omitted: %#v", payload)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true,"result":{"message_id":1,"chat":{"id":123,"type":"private"},"date":1}}`))
	}))
	defer server.Close()

	bot := newTestBot(t, token, server.URL, server.Client())
	_, err := bot.SendInvoice(context.Background(), validSendInvoiceParams())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCreateInvoiceLinkSendsPayloadAndDecodesLink(t *testing.T) {
	const token = "123:secret"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/bot"+token+"/createInvoiceLink" {
			t.Fatalf("unexpected path: %q", r.URL.Path)
		}
		var payload map[string]any
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode payload: %v", err)
		}
		if payload["business_connection_id"] != "business" || payload["title"] != "Title" || payload["subscription_period"] != float64(2592000) {
			t.Fatalf("unexpected payload: %#v", payload)
		}
		if _, ok := payload["provider_token"]; ok {
			t.Fatalf("empty provider token should be omitted: %#v", payload)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true,"result":"https://t.me/invoice/link"}`))
	}))
	defer server.Close()

	bot := newTestBot(t, token, server.URL, server.Client())
	link, err := bot.CreateInvoiceLink(context.Background(), CreateInvoiceLinkParams{
		BusinessConnectionID: "business",
		Title:                "Title",
		Description:          "Description",
		Payload:              "safe-payload",
		Currency:             "XTR",
		Prices:               []telegram.LabeledPrice{{Label: "Price", Amount: 150}},
		SubscriptionPeriod:   2592000,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if link != "https://t.me/invoice/link" {
		t.Fatalf("unexpected link: %q", link)
	}
}

func TestAnswerShippingQuerySendsOKAndErrorResponses(t *testing.T) {
	testBoolPaymentMethod(t, "answerShippingQuery", func(bot *Bot) (bool, error) {
		return bot.AnswerShippingQuery(context.Background(), AnswerShippingQueryParams{
			ShippingQueryID: "ship-id",
			OK:              true,
			ShippingOptions: []telegram.ShippingOption{{ID: "standard", Title: "Standard", Prices: []telegram.LabeledPrice{{Label: "Shipping", Amount: 25}}}},
		})
	}, func(t *testing.T, payload map[string]any) {
		if payload["shipping_query_id"] != "ship-id" || payload["ok"] != true {
			t.Fatalf("unexpected payload: %#v", payload)
		}
		options, ok := payload["shipping_options"].([]any)
		if !ok || len(options) != 1 {
			t.Fatalf("unexpected shipping options: %#v", payload["shipping_options"])
		}
	})

	testBoolPaymentMethod(t, "answerShippingQuery", func(bot *Bot) (bool, error) {
		return bot.AnswerShippingQuery(context.Background(), AnswerShippingQueryParams{ShippingQueryID: "ship-id", OK: false, ErrorMessage: "No shipping"})
	}, func(t *testing.T, payload map[string]any) {
		if payload["shipping_query_id"] != "ship-id" || payload["ok"] != false || payload["error_message"] != "No shipping" {
			t.Fatalf("unexpected payload: %#v", payload)
		}
	})
}

func TestAnswerPreCheckoutQuerySendsOKAndErrorResponses(t *testing.T) {
	testBoolPaymentMethod(t, "answerPreCheckoutQuery", func(bot *Bot) (bool, error) {
		return bot.AnswerPreCheckoutQuery(context.Background(), AnswerPreCheckoutQueryParams{PreCheckoutQueryID: "pre-id", OK: true})
	}, func(t *testing.T, payload map[string]any) {
		if payload["pre_checkout_query_id"] != "pre-id" || payload["ok"] != true {
			t.Fatalf("unexpected payload: %#v", payload)
		}
	})

	testBoolPaymentMethod(t, "answerPreCheckoutQuery", func(bot *Bot) (bool, error) {
		return bot.AnswerPreCheckoutQuery(context.Background(), AnswerPreCheckoutQueryParams{PreCheckoutQueryID: "pre-id", OK: false, ErrorMessage: "Cannot process"})
	}, func(t *testing.T, payload map[string]any) {
		if payload["pre_checkout_query_id"] != "pre-id" || payload["ok"] != false || payload["error_message"] != "Cannot process" {
			t.Fatalf("unexpected payload: %#v", payload)
		}
	})
}

func TestPaymentMethodsValidation(t *testing.T) {
	const token = "123:secret"
	bot := newTestBot(t, token, "https://example.test", nil)
	tests := []struct {
		name string
		call func() error
	}{
		{name: "send missing chat", call: func() error {
			params := validSendInvoiceParams()
			params.ChatID = ChatID{}
			_, err := bot.SendInvoice(context.Background(), params)
			return err
		}},
		{name: "send missing title", call: func() error {
			params := validSendInvoiceParams()
			params.Title = ""
			_, err := bot.SendInvoice(context.Background(), params)
			return err
		}},
		{name: "send missing description", call: func() error {
			params := validSendInvoiceParams()
			params.Description = ""
			_, err := bot.SendInvoice(context.Background(), params)
			return err
		}},
		{name: "send missing payload", call: func() error {
			params := validSendInvoiceParams()
			params.Payload = ""
			_, err := bot.SendInvoice(context.Background(), params)
			return err
		}},
		{name: "send missing currency", call: func() error {
			params := validSendInvoiceParams()
			params.Currency = ""
			_, err := bot.SendInvoice(context.Background(), params)
			return err
		}},
		{name: "send missing prices", call: func() error {
			params := validSendInvoiceParams()
			params.Prices = nil
			_, err := bot.SendInvoice(context.Background(), params)
			return err
		}},
		{name: "send empty price label", call: func() error {
			params := validSendInvoiceParams()
			params.Prices = []telegram.LabeledPrice{{Amount: 1}}
			_, err := bot.SendInvoice(context.Background(), params)
			return err
		}},
		{name: "send negative price", call: func() error {
			params := validSendInvoiceParams()
			params.Prices = []telegram.LabeledPrice{{Label: "Price", Amount: -1}}
			_, err := bot.SendInvoice(context.Background(), params)
			return err
		}},
		{name: "send negative message thread", call: func() error {
			params := validSendInvoiceParams()
			params.MessageThreadID = -1
			_, err := bot.SendInvoice(context.Background(), params)
			return err
		}},
		{name: "send negative direct topic", call: func() error {
			params := validSendInvoiceParams()
			params.DirectMessagesTopicID = -1
			_, err := bot.SendInvoice(context.Background(), params)
			return err
		}},
		{name: "send negative max tip", call: func() error {
			params := validSendInvoiceParams()
			params.MaxTipAmount = -1
			_, err := bot.SendInvoice(context.Background(), params)
			return err
		}},
		{name: "send negative suggested tip", call: func() error {
			params := validSendInvoiceParams()
			params.SuggestedTipAmounts = []int64{-1}
			_, err := bot.SendInvoice(context.Background(), params)
			return err
		}},
		{name: "send invalid photo url", call: func() error {
			params := validSendInvoiceParams()
			params.PhotoURL = "ftp://example.test/photo.png"
			_, err := bot.SendInvoice(context.Background(), params)
			return err
		}},
		{name: "send negative photo size", call: func() error {
			params := validSendInvoiceParams()
			params.PhotoSize = -1
			_, err := bot.SendInvoice(context.Background(), params)
			return err
		}},
		{name: "send invalid reply parameters", call: func() error {
			params := validSendInvoiceParams()
			params.ReplyParameters = &telegram.ReplyParameters{}
			_, err := bot.SendInvoice(context.Background(), params)
			return err
		}},
		{name: "send invalid reply markup", call: func() error {
			params := validSendInvoiceParams()
			invalid := telegram.InlineKeyboardMarkup{}
			params.ReplyMarkup = &invalid
			_, err := bot.SendInvoice(context.Background(), params)
			return err
		}},
		{name: "send invalid suggested post", call: func() error {
			params := validSendInvoiceParams()
			params.SuggestedPostParameters = &telegram.SuggestedPostParameters{Price: &telegram.SuggestedPostPrice{Amount: 1}}
			_, err := bot.SendInvoice(context.Background(), params)
			return err
		}},
		{name: "link missing title", call: func() error {
			params := validCreateInvoiceLinkParams()
			params.Title = ""
			_, err := bot.CreateInvoiceLink(context.Background(), params)
			return err
		}},
		{name: "link negative subscription", call: func() error {
			params := validCreateInvoiceLinkParams()
			params.SubscriptionPeriod = -1
			_, err := bot.CreateInvoiceLink(context.Background(), params)
			return err
		}},
		{name: "shipping missing id", call: func() error {
			_, err := bot.AnswerShippingQuery(context.Background(), AnswerShippingQueryParams{OK: true, ShippingOptions: []telegram.ShippingOption{{ID: "id", Title: "Title", Prices: []telegram.LabeledPrice{{Label: "Price", Amount: 1}}}}})
			return err
		}},
		{name: "shipping ok false missing error", call: func() error {
			_, err := bot.AnswerShippingQuery(context.Background(), AnswerShippingQueryParams{ShippingQueryID: "id"})
			return err
		}},
		{name: "shipping ok true error conflict", call: func() error {
			_, err := bot.AnswerShippingQuery(context.Background(), AnswerShippingQueryParams{ShippingQueryID: "id", OK: true, ErrorMessage: "bad", ShippingOptions: []telegram.ShippingOption{{ID: "id", Title: "Title", Prices: []telegram.LabeledPrice{{Label: "Price", Amount: 1}}}}})
			return err
		}},
		{name: "shipping ok true missing options", call: func() error {
			_, err := bot.AnswerShippingQuery(context.Background(), AnswerShippingQueryParams{ShippingQueryID: "id", OK: true})
			return err
		}},
		{name: "shipping invalid option", call: func() error {
			_, err := bot.AnswerShippingQuery(context.Background(), AnswerShippingQueryParams{ShippingQueryID: "id", OK: true, ShippingOptions: []telegram.ShippingOption{{Title: "Title", Prices: []telegram.LabeledPrice{{Label: "Price", Amount: 1}}}}})
			return err
		}},
		{name: "shipping ok false options conflict", call: func() error {
			_, err := bot.AnswerShippingQuery(context.Background(), AnswerShippingQueryParams{ShippingQueryID: "id", OK: false, ErrorMessage: "bad", ShippingOptions: []telegram.ShippingOption{{ID: "id", Title: "Title", Prices: []telegram.LabeledPrice{{Label: "Price", Amount: 1}}}}})
			return err
		}},
		{name: "precheckout missing id", call: func() error {
			_, err := bot.AnswerPreCheckoutQuery(context.Background(), AnswerPreCheckoutQueryParams{OK: true})
			return err
		}},
		{name: "precheckout ok false missing error", call: func() error {
			_, err := bot.AnswerPreCheckoutQuery(context.Background(), AnswerPreCheckoutQueryParams{PreCheckoutQueryID: "id"})
			return err
		}},
		{name: "precheckout ok true error conflict", call: func() error {
			_, err := bot.AnswerPreCheckoutQuery(context.Background(), AnswerPreCheckoutQueryParams{PreCheckoutQueryID: "id", OK: true, ErrorMessage: "bad"})
			return err
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.call()
			if err == nil {
				t.Fatal("expected error")
			}
			assertNoToken(t, err, token)
		})
	}
}

func TestPaymentMethodsReturnAPIError(t *testing.T) {
	const token = "123:secret"
	tests := []struct {
		name   string
		method string
		call   func(*Bot) error
	}{
		{name: "send invoice", method: "sendInvoice", call: func(bot *Bot) error {
			_, err := bot.SendInvoice(context.Background(), validSendInvoiceParams())
			return err
		}},
		{name: "create invoice link", method: "createInvoiceLink", call: func(bot *Bot) error {
			_, err := bot.CreateInvoiceLink(context.Background(), validCreateInvoiceLinkParams())
			return err
		}},
		{name: "answer shipping", method: "answerShippingQuery", call: func(bot *Bot) error {
			_, err := bot.AnswerShippingQuery(context.Background(), validAnswerShippingQueryParams(true))
			return err
		}},
		{name: "answer precheckout", method: "answerPreCheckoutQuery", call: func(bot *Bot) error {
			_, err := bot.AnswerPreCheckoutQuery(context.Background(), AnswerPreCheckoutQueryParams{PreCheckoutQueryID: "pre", OK: true})
			return err
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path != "/bot"+token+"/"+tt.method {
					t.Fatalf("unexpected path: %q", r.URL.Path)
				}
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(`{"ok":false,"error_code":400,"description":"Bad Request"}`))
			}))
			defer server.Close()

			bot := newTestBot(t, token, server.URL, server.Client())
			err := tt.call(bot)
			if err == nil {
				t.Fatal("expected error")
			}
			var apiErr *apierrors.APIError
			if !stderrors.As(err, &apiErr) {
				t.Fatalf("expected APIError, got %T", err)
			}
			assertNoToken(t, err, token)
		})
	}
}

func TestPaymentMethodsResponseAndContextErrors(t *testing.T) {
	const token = "123:secret"
	tests := []struct {
		name   string
		method string
		call   func(context.Context, *Bot) error
	}{
		{name: "send invoice", method: "sendInvoice", call: func(ctx context.Context, bot *Bot) error {
			_, err := bot.SendInvoice(ctx, validSendInvoiceParams())
			return err
		}},
		{name: "create invoice link", method: "createInvoiceLink", call: func(ctx context.Context, bot *Bot) error {
			_, err := bot.CreateInvoiceLink(ctx, validCreateInvoiceLinkParams())
			return err
		}},
		{name: "answer shipping", method: "answerShippingQuery", call: func(ctx context.Context, bot *Bot) error {
			_, err := bot.AnswerShippingQuery(ctx, validAnswerShippingQueryParams(true))
			return err
		}},
		{name: "answer precheckout", method: "answerPreCheckoutQuery", call: func(ctx context.Context, bot *Bot) error {
			_, err := bot.AnswerPreCheckoutQuery(ctx, AnswerPreCheckoutQueryParams{PreCheckoutQueryID: "pre", OK: true})
			return err
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name+" cancelled context", func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				t.Fatal("request should not reach server")
			}))
			defer server.Close()
			bot := newTestBot(t, token, server.URL, server.Client())
			ctx, cancel := context.WithCancel(context.Background())
			cancel()
			err := tt.call(ctx, bot)
			if err == nil {
				t.Fatal("expected error")
			}
			assertNoToken(t, err, token)
		})
		t.Run(tt.name+" invalid json", func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path != "/bot"+token+"/"+tt.method {
					t.Fatalf("unexpected path: %q", r.URL.Path)
				}
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(`not-json`))
			}))
			defer server.Close()
			bot := newTestBot(t, token, server.URL, server.Client())
			err := tt.call(context.Background(), bot)
			if err == nil {
				t.Fatal("expected error")
			}
			assertNoToken(t, err, token)
		})
		t.Run(tt.name+" http status", func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path != "/bot"+token+"/"+tt.method {
					t.Fatalf("unexpected path: %q", r.URL.Path)
				}
				http.Error(w, "server error", http.StatusInternalServerError)
			}))
			defer server.Close()
			bot := newTestBot(t, token, server.URL, server.Client())
			err := tt.call(context.Background(), bot)
			if err == nil {
				t.Fatal("expected error")
			}
			assertNoToken(t, err, token)
		})
	}
}

func validSendInvoiceParams() SendInvoiceParams {
	return SendInvoiceParams{
		ChatID:      ChatIDInt(123),
		Title:       "Title",
		Description: "Description",
		Payload:     "safe-payload",
		Currency:    "XTR",
		Prices:      []telegram.LabeledPrice{{Label: "Price", Amount: 1}},
	}
}

func validCreateInvoiceLinkParams() CreateInvoiceLinkParams {
	return CreateInvoiceLinkParams{
		Title:       "Title",
		Description: "Description",
		Payload:     "safe-payload",
		Currency:    "XTR",
		Prices:      []telegram.LabeledPrice{{Label: "Price", Amount: 1}},
	}
}

func validAnswerShippingQueryParams(ok bool) AnswerShippingQueryParams {
	if !ok {
		return AnswerShippingQueryParams{ShippingQueryID: "ship", OK: false, ErrorMessage: "No shipping"}
	}
	return AnswerShippingQueryParams{
		ShippingQueryID: "ship",
		OK:              true,
		ShippingOptions: []telegram.ShippingOption{{ID: "standard", Title: "Standard", Prices: []telegram.LabeledPrice{{Label: "Shipping", Amount: 1}}}},
	}
}

func testBoolPaymentMethod(t *testing.T, method string, call func(*Bot) (bool, error), assertPayload func(*testing.T, map[string]any)) {
	t.Helper()
	const token = "123:secret"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/bot"+token+"/"+method {
			t.Fatalf("unexpected path: %q", r.URL.Path)
		}
		var payload map[string]any
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode payload: %v", err)
		}
		assertPayload(t, payload)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true,"result":true}`))
	}))
	defer server.Close()

	bot := newTestBot(t, token, server.URL, server.Client())
	ok, err := call(bot)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Fatal("expected true result")
	}
}

func TestPaymentValidationDoesNotLeakPayloadValues(t *testing.T) {
	const token = "123:secret"
	bot := newTestBot(t, token, "https://example.test", nil)
	_, err := bot.SendInvoice(context.Background(), SendInvoiceParams{
		ChatID:      ChatIDInt(123),
		Title:       "Title",
		Description: "Description",
		Payload:     "sensitive-payload-value",
		Currency:    "XTR",
		Prices:      []telegram.LabeledPrice{{Label: "", Amount: 1}},
	})
	if err == nil {
		t.Fatal("expected error")
	}
	if strings.Contains(err.Error(), "sensitive-payload-value") {
		t.Fatalf("error leaked payload: %v", err)
	}
}
