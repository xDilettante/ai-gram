package bot

import (
	"context"
	"encoding/json"
	stderrors "errors"
	"net/http"
	"net/http/httptest"
	"testing"

	apierrors "github.com/xDilettante/ai-gram/errors"
	"github.com/xDilettante/ai-gram/telegram"
)

func TestSendContactSendsPayloadAndDecodesMessage(t *testing.T) {
	const token = "123:secret"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/bot"+token+"/sendContact" {
			t.Fatalf("unexpected path: %q", r.URL.Path)
		}
		var payload map[string]any
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode payload: %v", err)
		}
		if got := payload["chat_id"]; got != float64(12345) {
			t.Fatalf("unexpected chat_id: %#v", got)
		}
		if got := payload["phone_number"]; got != "+15551234567" {
			t.Fatalf("unexpected phone_number: %#v", got)
		}
		if got := payload["first_name"]; got != "Ada" {
			t.Fatalf("unexpected first_name: %#v", got)
		}
		if got := payload["last_name"]; got != "Lovelace" {
			t.Fatalf("unexpected last_name: %#v", got)
		}
		if got := payload["vcard"]; got != "BEGIN:VCARD" {
			t.Fatalf("unexpected vcard: %#v", got)
		}
		if got := payload["message_thread_id"]; got != float64(77) {
			t.Fatalf("unexpected message_thread_id: %#v", got)
		}
		reply := payload["reply_parameters"].(map[string]any)
		if got := reply["message_id"]; got != float64(10) {
			t.Fatalf("unexpected reply message_id: %#v", got)
		}
		markup := payload["reply_markup"].(map[string]any)
		if _, ok := markup["inline_keyboard"]; !ok {
			t.Fatalf("reply_markup.inline_keyboard missing: %#v", markup)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true,"result":{"message_id":5,"chat":{"id":12345,"type":"private"},"date":100,"contact":{"phone_number":"+15551234567","first_name":"Ada","last_name":"Lovelace","vcard":"BEGIN:VCARD"}}}`))
	}))
	defer server.Close()

	bot := newTestBot(t, token, server.URL, server.Client())
	message, err := bot.SendContact(context.Background(), SendContactParams{
		ChatID:          ChatIDInt(12345),
		MessageThreadID: 77,
		PhoneNumber:     "+15551234567",
		FirstName:       "Ada",
		LastName:        "Lovelace",
		VCard:           "BEGIN:VCARD",
		ReplyParameters: &telegram.ReplyParameters{MessageID: 10},
		ReplyMarkup:     telegram.NewInlineKeyboard([]telegram.InlineKeyboardButton{telegram.InlineButtonCallback("OK", "ok")}),
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if message == nil || message.Contact == nil || message.Contact.PhoneNumber != "+15551234567" || message.Contact.FirstName != "Ada" {
		t.Fatalf("unexpected message: %+v", message)
	}
}

func TestSendContactValidation(t *testing.T) {
	const token = "123:secret"
	tests := []struct {
		name   string
		params SendContactParams
	}{
		{name: "empty chat", params: SendContactParams{PhoneNumber: "+1", FirstName: "Ada"}},
		{name: "empty phone", params: SendContactParams{ChatID: ChatIDInt(12345), FirstName: "Ada"}},
		{name: "empty first name", params: SendContactParams{ChatID: ChatIDInt(12345), PhoneNumber: "+1"}},
		{name: "negative thread", params: SendContactParams{ChatID: ChatIDInt(12345), MessageThreadID: -1, PhoneNumber: "+1", FirstName: "Ada"}},
		{name: "invalid reply parameters", params: SendContactParams{ChatID: ChatIDInt(12345), PhoneNumber: "+1", FirstName: "Ada", ReplyParameters: &telegram.ReplyParameters{}}},
		{name: "invalid reply markup", params: SendContactParams{ChatID: ChatIDInt(12345), PhoneNumber: "+1", FirstName: "Ada", ReplyMarkup: telegram.InlineKeyboardMarkup{}}},
	}

	bot := newTestBot(t, token, "https://example.test", nil)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			message, err := bot.SendContact(context.Background(), tt.params)
			if err == nil {
				t.Fatal("expected error")
			}
			if message != nil {
				t.Fatalf("expected nil message, got %+v", message)
			}
			assertNoToken(t, err, token)
		})
	}
}

func TestSendContactAPIAndTransportErrors(t *testing.T) {
	testSendContactErrorCases(t)
}

func TestSendLocationSendsPayloadAndDecodesMessage(t *testing.T) {
	const token = "123:secret"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/bot"+token+"/sendLocation" {
			t.Fatalf("unexpected path: %q", r.URL.Path)
		}
		var payload map[string]any
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode payload: %v", err)
		}
		if got := payload["latitude"]; got != 51.5 {
			t.Fatalf("unexpected latitude: %#v", got)
		}
		if got := payload["longitude"]; got != -0.12 {
			t.Fatalf("unexpected longitude: %#v", got)
		}
		if got := payload["horizontal_accuracy"]; got != 20.5 {
			t.Fatalf("unexpected horizontal_accuracy: %#v", got)
		}
		if got := payload["live_period"]; got != float64(60) {
			t.Fatalf("unexpected live_period: %#v", got)
		}
		if got := payload["heading"]; got != float64(90) {
			t.Fatalf("unexpected heading: %#v", got)
		}
		if got := payload["proximity_alert_radius"]; got != float64(100) {
			t.Fatalf("unexpected proximity_alert_radius: %#v", got)
		}
		if got := payload["message_thread_id"]; got != float64(8) {
			t.Fatalf("unexpected message_thread_id: %#v", got)
		}
		reply := payload["reply_parameters"].(map[string]any)
		if got := reply["message_id"]; got != float64(11) {
			t.Fatalf("unexpected reply message_id: %#v", got)
		}
		markup := payload["reply_markup"].(map[string]any)
		if _, ok := markup["inline_keyboard"]; !ok {
			t.Fatalf("reply_markup.inline_keyboard missing: %#v", markup)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true,"result":{"message_id":6,"chat":{"id":12345,"type":"private"},"date":100,"location":{"latitude":51.5,"longitude":-0.12,"horizontal_accuracy":20.5,"live_period":60,"heading":90,"proximity_alert_radius":100}}}`))
	}))
	defer server.Close()

	bot := newTestBot(t, token, server.URL, server.Client())
	message, err := bot.SendLocation(context.Background(), SendLocationParams{
		ChatID:               ChatIDInt(12345),
		MessageThreadID:      8,
		Latitude:             51.5,
		Longitude:            -0.12,
		HorizontalAccuracy:   20.5,
		LivePeriod:           60,
		Heading:              90,
		ProximityAlertRadius: 100,
		ReplyParameters:      &telegram.ReplyParameters{MessageID: 11},
		ReplyMarkup:          telegram.NewInlineKeyboard([]telegram.InlineKeyboardButton{telegram.InlineButtonCallback("OK", "ok")}),
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if message == nil || message.Location == nil || message.Location.Latitude != 51.5 || message.Location.Longitude != -0.12 {
		t.Fatalf("unexpected message: %+v", message)
	}
}

func TestSendLocationValidation(t *testing.T) {
	const token = "123:secret"
	valid := SendLocationParams{ChatID: ChatIDInt(12345), Latitude: 10, Longitude: 20}
	tests := []struct {
		name   string
		mutate func(*SendLocationParams)
	}{
		{name: "empty chat", mutate: func(p *SendLocationParams) { p.ChatID = ChatID{} }},
		{name: "latitude too small", mutate: func(p *SendLocationParams) { p.Latitude = -91 }},
		{name: "latitude too large", mutate: func(p *SendLocationParams) { p.Latitude = 91 }},
		{name: "longitude too small", mutate: func(p *SendLocationParams) { p.Longitude = -181 }},
		{name: "longitude too large", mutate: func(p *SendLocationParams) { p.Longitude = 181 }},
		{name: "negative horizontal accuracy", mutate: func(p *SendLocationParams) { p.HorizontalAccuracy = -1 }},
		{name: "negative live period", mutate: func(p *SendLocationParams) { p.LivePeriod = -1 }},
		{name: "negative heading", mutate: func(p *SendLocationParams) { p.Heading = -1 }},
		{name: "negative proximity alert radius", mutate: func(p *SendLocationParams) { p.ProximityAlertRadius = -1 }},
		{name: "negative thread", mutate: func(p *SendLocationParams) { p.MessageThreadID = -1 }},
		{name: "invalid reply parameters", mutate: func(p *SendLocationParams) { p.ReplyParameters = &telegram.ReplyParameters{} }},
		{name: "invalid reply markup", mutate: func(p *SendLocationParams) { p.ReplyMarkup = telegram.InlineKeyboardMarkup{} }},
	}

	bot := newTestBot(t, token, "https://example.test", nil)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			params := valid
			tt.mutate(&params)
			message, err := bot.SendLocation(context.Background(), params)
			if err == nil {
				t.Fatal("expected error")
			}
			if message != nil {
				t.Fatalf("expected nil message, got %+v", message)
			}
			assertNoToken(t, err, token)
		})
	}
}

func TestSendLocationAPIAndTransportErrors(t *testing.T) {
	testSendLocationErrorCases(t)
}

func TestSendVenueSendsPayloadAndDecodesMessage(t *testing.T) {
	const token = "123:secret"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/bot"+token+"/sendVenue" {
			t.Fatalf("unexpected path: %q", r.URL.Path)
		}
		var payload map[string]any
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode payload: %v", err)
		}
		if got := payload["latitude"]; got != 51.5 {
			t.Fatalf("unexpected latitude: %#v", got)
		}
		if got := payload["longitude"]; got != -0.12 {
			t.Fatalf("unexpected longitude: %#v", got)
		}
		if got := payload["title"]; got != "Venue title" {
			t.Fatalf("unexpected title: %#v", got)
		}
		if got := payload["address"]; got != "Venue address" {
			t.Fatalf("unexpected address: %#v", got)
		}
		if got := payload["foursquare_id"]; got != "fs-id" {
			t.Fatalf("unexpected foursquare_id: %#v", got)
		}
		if got := payload["foursquare_type"]; got != "arts/default" {
			t.Fatalf("unexpected foursquare_type: %#v", got)
		}
		if got := payload["google_place_id"]; got != "google-id" {
			t.Fatalf("unexpected google_place_id: %#v", got)
		}
		if got := payload["google_place_type"]; got != "restaurant" {
			t.Fatalf("unexpected google_place_type: %#v", got)
		}
		if got := payload["message_thread_id"]; got != float64(9) {
			t.Fatalf("unexpected message_thread_id: %#v", got)
		}
		reply := payload["reply_parameters"].(map[string]any)
		if got := reply["message_id"]; got != float64(12) {
			t.Fatalf("unexpected reply message_id: %#v", got)
		}
		markup := payload["reply_markup"].(map[string]any)
		if _, ok := markup["inline_keyboard"]; !ok {
			t.Fatalf("reply_markup.inline_keyboard missing: %#v", markup)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true,"result":{"message_id":7,"chat":{"id":12345,"type":"private"},"date":100,"venue":{"location":{"latitude":51.5,"longitude":-0.12},"title":"Venue title","address":"Venue address","foursquare_id":"fs-id","foursquare_type":"arts/default","google_place_id":"google-id","google_place_type":"restaurant"}}}`))
	}))
	defer server.Close()

	bot := newTestBot(t, token, server.URL, server.Client())
	message, err := bot.SendVenue(context.Background(), SendVenueParams{
		ChatID:          ChatIDInt(12345),
		MessageThreadID: 9,
		Latitude:        51.5,
		Longitude:       -0.12,
		Title:           "Venue title",
		Address:         "Venue address",
		FoursquareID:    "fs-id",
		FoursquareType:  "arts/default",
		GooglePlaceID:   "google-id",
		GooglePlaceType: "restaurant",
		ReplyParameters: &telegram.ReplyParameters{MessageID: 12},
		ReplyMarkup:     telegram.NewInlineKeyboard([]telegram.InlineKeyboardButton{telegram.InlineButtonCallback("OK", "ok")}),
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if message == nil || message.Venue == nil || message.Venue.Title != "Venue title" || message.Venue.Location.Latitude != 51.5 {
		t.Fatalf("unexpected message: %+v", message)
	}
}

func TestSendVenueValidation(t *testing.T) {
	const token = "123:secret"
	valid := SendVenueParams{ChatID: ChatIDInt(12345), Latitude: 10, Longitude: 20, Title: "Venue", Address: "Address"}
	tests := []struct {
		name   string
		mutate func(*SendVenueParams)
	}{
		{name: "empty chat", mutate: func(p *SendVenueParams) { p.ChatID = ChatID{} }},
		{name: "latitude too small", mutate: func(p *SendVenueParams) { p.Latitude = -91 }},
		{name: "longitude too large", mutate: func(p *SendVenueParams) { p.Longitude = 181 }},
		{name: "empty title", mutate: func(p *SendVenueParams) { p.Title = "" }},
		{name: "empty address", mutate: func(p *SendVenueParams) { p.Address = "" }},
		{name: "negative thread", mutate: func(p *SendVenueParams) { p.MessageThreadID = -1 }},
		{name: "invalid reply parameters", mutate: func(p *SendVenueParams) { p.ReplyParameters = &telegram.ReplyParameters{} }},
		{name: "invalid reply markup", mutate: func(p *SendVenueParams) { p.ReplyMarkup = telegram.InlineKeyboardMarkup{} }},
	}

	bot := newTestBot(t, token, "https://example.test", nil)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			params := valid
			tt.mutate(&params)
			message, err := bot.SendVenue(context.Background(), params)
			if err == nil {
				t.Fatal("expected error")
			}
			if message != nil {
				t.Fatalf("expected nil message, got %+v", message)
			}
			assertNoToken(t, err, token)
		})
	}
}

func TestSendVenueAPIAndTransportErrors(t *testing.T) {
	testSendVenueErrorCases(t)
}

func testSendContactErrorCases(t *testing.T) {
	t.Helper()
	valid := SendContactParams{ChatID: ChatIDInt(12345), PhoneNumber: "+1", FirstName: "Ada"}
	testSendMethodErrorCases(t, "sendContact", func(bot *Bot) (*telegram.Message, error) {
		return bot.SendContact(context.Background(), valid)
	}, func(bot *Bot, ctx context.Context) (*telegram.Message, error) {
		return bot.SendContact(ctx, valid)
	})
}

func testSendLocationErrorCases(t *testing.T) {
	t.Helper()
	valid := SendLocationParams{ChatID: ChatIDInt(12345), Latitude: 10, Longitude: 20}
	testSendMethodErrorCases(t, "sendLocation", func(bot *Bot) (*telegram.Message, error) {
		return bot.SendLocation(context.Background(), valid)
	}, func(bot *Bot, ctx context.Context) (*telegram.Message, error) {
		return bot.SendLocation(ctx, valid)
	})
}

func testSendVenueErrorCases(t *testing.T) {
	t.Helper()
	valid := SendVenueParams{ChatID: ChatIDInt(12345), Latitude: 10, Longitude: 20, Title: "Venue", Address: "Address"}
	testSendMethodErrorCases(t, "sendVenue", func(bot *Bot) (*telegram.Message, error) {
		return bot.SendVenue(context.Background(), valid)
	}, func(bot *Bot, ctx context.Context) (*telegram.Message, error) {
		return bot.SendVenue(ctx, valid)
	})
}

func testSendMethodErrorCases(t *testing.T, method string, call func(*Bot) (*telegram.Message, error), callWithContext func(*Bot, context.Context) (*telegram.Message, error)) {
	t.Helper()
	const token = "123:secret"
	t.Run("api error", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/bot"+token+"/"+method {
				t.Fatalf("unexpected path: %q", r.URL.Path)
			}
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"ok":false,"error_code":400,"description":"Bad Request"}`))
		}))
		defer server.Close()

		bot := newTestBot(t, token, server.URL, server.Client())
		message, err := call(bot)
		if err == nil {
			t.Fatal("expected error")
		}
		if message != nil {
			t.Fatalf("expected nil message, got %+v", message)
		}
		var apiErr *apierrors.APIError
		if !stderrors.As(err, &apiErr) {
			t.Fatalf("expected APIError, got %T", err)
		}
		assertNoToken(t, err, token)
	})

	t.Run("invalid json", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`not-json`))
		}))
		defer server.Close()

		bot := newTestBot(t, token, server.URL, server.Client())
		_, err := call(bot)
		if err == nil {
			t.Fatal("expected error")
		}
		assertNoToken(t, err, token)
	})

	t.Run("http status", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "server error", http.StatusInternalServerError)
		}))
		defer server.Close()

		bot := newTestBot(t, token, server.URL, server.Client())
		_, err := call(bot)
		if err == nil {
			t.Fatal("expected error")
		}
		assertNoToken(t, err, token)
	})

	t.Run("cancelled context", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			t.Fatal("request should not reach server")
		}))
		defer server.Close()

		bot := newTestBot(t, token, server.URL, server.Client())
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		_, err := callWithContext(bot, ctx)
		if err == nil {
			t.Fatal("expected error")
		}
		assertNoToken(t, err, token)
	})
}
