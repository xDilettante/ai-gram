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

func TestInputTextMessageContentMarshalAndValidation(t *testing.T) {
	content := InputTextMessageContent{
		MessageText: "hello",
		ParseMode:   "HTML",
		LinkPreviewOptions: &telegram.LinkPreviewOptions{
			IsDisabled: true,
		},
	}
	data, err := json.Marshal(content)
	if err != nil {
		t.Fatalf("marshal content: %v", err)
	}
	var got map[string]any
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("decode content: %v", err)
	}
	if got["message_text"] != "hello" || got["parse_mode"] != "HTML" {
		t.Fatalf("unexpected content payload: %#v", got)
	}
	if _, ok := got["link_preview_options"].(map[string]any); !ok {
		t.Fatalf("expected link_preview_options object: %#v", got)
	}
	if err := validateInputMessageContent(content); err != nil {
		t.Fatalf("valid content rejected: %v", err)
	}

	tests := []struct {
		name    string
		content InputMessageContent
	}{
		{name: "empty text", content: InputTextMessageContent{}},
		{name: "parse mode and entities", content: InputTextMessageContent{MessageText: "hello", ParseMode: "HTML", Entities: []telegram.MessageEntity{{Type: telegram.EntityBold, Offset: 0, Length: 5}}}},
		{name: "invalid link preview", content: InputTextMessageContent{MessageText: "hello", LinkPreviewOptions: &telegram.LinkPreviewOptions{URL: "ftp://example.com"}}},
		{name: "unsupported content", content: unsupportedInputMessageContent{}},
	}
	var typedNil *InputTextMessageContent
	tests = append(tests, struct {
		name    string
		content InputMessageContent
	}{name: "typed nil", content: typedNil})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateInputMessageContent(tt.content); err == nil {
				t.Fatal("expected validation error")
			}
		})
	}
}

func TestInputLocationMessageContentMarshalAndValidation(t *testing.T) {
	content := InputLocation(52.3676, 4.9041)
	content.HorizontalAccuracy = 12.5
	content.LivePeriod = 60
	content.Heading = 90
	content.ProximityAlertRadius = 100

	data, err := json.Marshal(content)
	if err != nil {
		t.Fatalf("marshal content: %v", err)
	}
	var got map[string]any
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("decode content: %v", err)
	}
	if got["latitude"] != 52.3676 || got["longitude"] != 4.9041 || got["heading"] != float64(90) {
		t.Fatalf("unexpected location content: %#v", got)
	}
	if err := validateInputMessageContent(content); err != nil {
		t.Fatalf("valid location content rejected: %v", err)
	}

	tests := []struct {
		name    string
		content InputMessageContent
	}{
		{name: "latitude too small", content: InputLocationMessageContent{Latitude: -91, Longitude: 4}},
		{name: "latitude too large", content: InputLocationMessageContent{Latitude: 91, Longitude: 4}},
		{name: "longitude too small", content: InputLocationMessageContent{Latitude: 52, Longitude: -181}},
		{name: "longitude too large", content: InputLocationMessageContent{Latitude: 52, Longitude: 181}},
		{name: "negative horizontal accuracy", content: InputLocationMessageContent{Latitude: 52, Longitude: 4, HorizontalAccuracy: -1}},
		{name: "negative live period", content: InputLocationMessageContent{Latitude: 52, Longitude: 4, LivePeriod: -1}},
		{name: "negative heading", content: InputLocationMessageContent{Latitude: 52, Longitude: 4, Heading: -1}},
		{name: "negative proximity alert radius", content: InputLocationMessageContent{Latitude: 52, Longitude: 4, ProximityAlertRadius: -1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateInputMessageContent(tt.content); err == nil {
				t.Fatal("expected validation error")
			}
		})
	}
}

func TestInputVenueMessageContentMarshalAndValidation(t *testing.T) {
	content := InputVenue(52.3676, 4.9041, "ai-gram venue", "test address")
	content.FoursquareID = "fs"
	content.GooglePlaceID = "gp"

	data, err := json.Marshal(content)
	if err != nil {
		t.Fatalf("marshal content: %v", err)
	}
	var got map[string]any
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("decode content: %v", err)
	}
	if got["title"] != "ai-gram venue" || got["address"] != "test address" || got["foursquare_id"] != "fs" || got["google_place_id"] != "gp" {
		t.Fatalf("unexpected venue content: %#v", got)
	}
	if err := validateInputMessageContent(content); err != nil {
		t.Fatalf("valid venue content rejected: %v", err)
	}

	tests := []struct {
		name    string
		content InputMessageContent
	}{
		{name: "invalid coordinates", content: InputVenueMessageContent{Latitude: 100, Longitude: 4, Title: "venue", Address: "address"}},
		{name: "empty title", content: InputVenueMessageContent{Latitude: 52, Longitude: 4, Address: "address"}},
		{name: "empty address", content: InputVenueMessageContent{Latitude: 52, Longitude: 4, Title: "venue"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateInputMessageContent(tt.content); err == nil {
				t.Fatal("expected validation error")
			}
		})
	}
}

func TestInputContactMessageContentMarshalAndValidation(t *testing.T) {
	content := InputContact("+10000000000", "ai-gram")
	content.LastName = "Smoke"
	content.VCard = "BEGIN:VCARD\nEND:VCARD"

	data, err := json.Marshal(content)
	if err != nil {
		t.Fatalf("marshal content: %v", err)
	}
	var got map[string]any
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("decode content: %v", err)
	}
	if got["phone_number"] != "+10000000000" || got["first_name"] != "ai-gram" || got["last_name"] != "Smoke" {
		t.Fatalf("unexpected contact content: %#v", got)
	}
	if err := validateInputMessageContent(content); err != nil {
		t.Fatalf("valid contact content rejected: %v", err)
	}

	tests := []struct {
		name    string
		content InputMessageContent
	}{
		{name: "empty phone", content: InputContactMessageContent{FirstName: "ai-gram"}},
		{name: "empty first name", content: InputContactMessageContent{PhoneNumber: "+10000000000"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateInputMessageContent(tt.content); err == nil {
				t.Fatal("expected validation error")
			}
		})
	}
}

func TestInputInvoiceMessageContentMarshalAndValidation(t *testing.T) {
	content := InputInvoiceMessageContent{
		Title:               "Product",
		Description:         "Description",
		Payload:             "payload",
		Currency:            "XTR",
		Prices:              []telegram.LabeledPrice{{Label: "Price", Amount: 100}},
		MaxTipAmount:        10,
		SuggestedTipAmounts: []int64{1, 5},
		PhotoURL:            "https://example.com/photo.jpg",
		PhotoSize:           1000,
		PhotoWidth:          320,
		PhotoHeight:         240,
	}

	data, err := json.Marshal(content)
	if err != nil {
		t.Fatalf("marshal content: %v", err)
	}
	var got map[string]any
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("decode content: %v", err)
	}
	if got["title"] != "Product" || got["currency"] != "XTR" {
		t.Fatalf("unexpected invoice content: %#v", got)
	}
	prices, ok := got["prices"].([]any)
	if !ok || len(prices) != 1 {
		t.Fatalf("unexpected prices: %#v", got["prices"])
	}
	if err := validateInputMessageContent(content); err != nil {
		t.Fatalf("valid invoice content rejected: %v", err)
	}

	tests := []struct {
		name    string
		content InputMessageContent
	}{
		{name: "empty title", content: InputInvoiceMessageContent{Description: "d", Payload: "p", Currency: "XTR", Prices: []telegram.LabeledPrice{{Label: "Price", Amount: 1}}}},
		{name: "empty description", content: InputInvoiceMessageContent{Title: "t", Payload: "p", Currency: "XTR", Prices: []telegram.LabeledPrice{{Label: "Price", Amount: 1}}}},
		{name: "empty payload", content: InputInvoiceMessageContent{Title: "t", Description: "d", Currency: "XTR", Prices: []telegram.LabeledPrice{{Label: "Price", Amount: 1}}}},
		{name: "empty currency", content: InputInvoiceMessageContent{Title: "t", Description: "d", Payload: "p", Prices: []telegram.LabeledPrice{{Label: "Price", Amount: 1}}}},
		{name: "empty prices", content: InputInvoiceMessageContent{Title: "t", Description: "d", Payload: "p", Currency: "XTR"}},
		{name: "empty price label", content: InputInvoiceMessageContent{Title: "t", Description: "d", Payload: "p", Currency: "XTR", Prices: []telegram.LabeledPrice{{Amount: 1}}}},
		{name: "negative price", content: InputInvoiceMessageContent{Title: "t", Description: "d", Payload: "p", Currency: "XTR", Prices: []telegram.LabeledPrice{{Label: "Price", Amount: -1}}}},
		{name: "negative max tip", content: InputInvoiceMessageContent{Title: "t", Description: "d", Payload: "p", Currency: "XTR", Prices: []telegram.LabeledPrice{{Label: "Price", Amount: 1}}, MaxTipAmount: -1}},
		{name: "negative suggested tip", content: InputInvoiceMessageContent{Title: "t", Description: "d", Payload: "p", Currency: "XTR", Prices: []telegram.LabeledPrice{{Label: "Price", Amount: 1}}, SuggestedTipAmounts: []int64{-1}}},
		{name: "invalid photo url", content: InputInvoiceMessageContent{Title: "t", Description: "d", Payload: "p", Currency: "XTR", Prices: []telegram.LabeledPrice{{Label: "Price", Amount: 1}}, PhotoURL: "ftp://example.com/photo.jpg"}},
		{name: "negative photo size", content: InputInvoiceMessageContent{Title: "t", Description: "d", Payload: "p", Currency: "XTR", Prices: []telegram.LabeledPrice{{Label: "Price", Amount: 1}}, PhotoSize: -1}},
		{name: "negative photo width", content: InputInvoiceMessageContent{Title: "t", Description: "d", Payload: "p", Currency: "XTR", Prices: []telegram.LabeledPrice{{Label: "Price", Amount: 1}}, PhotoWidth: -1}},
		{name: "negative photo height", content: InputInvoiceMessageContent{Title: "t", Description: "d", Payload: "p", Currency: "XTR", Prices: []telegram.LabeledPrice{{Label: "Price", Amount: 1}}, PhotoHeight: -1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateInputMessageContent(tt.content); err == nil {
				t.Fatal("expected validation error")
			}
		})
	}
}

func TestInlineQueryResultArticleMarshalAndValidation(t *testing.T) {
	markup := telegram.NewInlineKeyboard([]telegram.InlineKeyboardButton{telegram.InlineButtonCallback("Open", "open")})
	result := InlineArticle("article-1", "Article", InputText("hello"))
	result.ReplyMarkup = &markup
	result.URL = "https://example.com/article"
	result.Description = "description"
	result.ThumbnailURL = "https://example.com/thumb.jpg"
	result.ThumbnailWidth = 100
	result.ThumbnailHeight = 50

	data, err := json.Marshal(result)
	if err != nil {
		t.Fatalf("marshal article: %v", err)
	}
	var got map[string]any
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("decode article: %v", err)
	}
	if got["type"] != "article" || got["id"] != "article-1" || got["title"] != "Article" {
		t.Fatalf("unexpected article payload: %#v", got)
	}
	content, ok := got["input_message_content"].(map[string]any)
	if !ok || content["message_text"] != "hello" {
		t.Fatalf("unexpected input_message_content: %#v", got["input_message_content"])
	}
	if err := validateInlineQueryResult(result); err != nil {
		t.Fatalf("valid article rejected: %v", err)
	}

	tests := []struct {
		name   string
		result InlineQueryResult
	}{
		{name: "empty id", result: InlineQueryResultArticle{Title: "Article", InputMessageContent: InputText("hello")}},
		{name: "empty title", result: InlineQueryResultArticle{ID: "article-1", InputMessageContent: InputText("hello")}},
		{name: "missing content", result: InlineQueryResultArticle{ID: "article-1", Title: "Article"}},
		{name: "negative thumbnail width", result: InlineQueryResultArticle{ID: "article-1", Title: "Article", InputMessageContent: InputText("hello"), ThumbnailWidth: -1}},
		{name: "negative thumbnail height", result: InlineQueryResultArticle{ID: "article-1", Title: "Article", InputMessageContent: InputText("hello"), ThumbnailHeight: -1}},
		{name: "invalid reply markup", result: InlineQueryResultArticle{ID: "article-1", Title: "Article", InputMessageContent: InputText("hello"), ReplyMarkup: &telegram.InlineKeyboardMarkup{}}},
		{name: "invalid url", result: InlineQueryResultArticle{ID: "article-1", Title: "Article", InputMessageContent: InputText("hello"), URL: "ftp://example.com"}},
		{name: "invalid thumbnail url", result: InlineQueryResultArticle{ID: "article-1", Title: "Article", InputMessageContent: InputText("hello"), ThumbnailURL: "ftp://example.com"}},
		{name: "unsupported result", result: unsupportedInlineQueryResult{}},
	}
	var typedNil *InlineQueryResultArticle
	tests = append(tests, struct {
		name   string
		result InlineQueryResult
	}{name: "typed nil", result: typedNil})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateInlineQueryResult(tt.result); err == nil {
				t.Fatal("expected validation error")
			}
		})
	}
}

func TestInlineQueryResultLocationMarshalAndValidation(t *testing.T) {
	markup := telegram.NewInlineKeyboard([]telegram.InlineKeyboardButton{telegram.InlineButtonCallback("Open", "open")})
	result := InlineLocation("location-1", 52.3676, 4.9041, "Amsterdam")
	result.ReplyMarkup = &markup
	result.InputMessageContent = InputLocation(52.3676, 4.9041)
	result.ThumbnailURL = "https://example.com/thumb.jpg"
	result.ThumbnailWidth = 100
	result.ThumbnailHeight = 50

	data, err := json.Marshal(result)
	if err != nil {
		t.Fatalf("marshal location: %v", err)
	}
	var got map[string]any
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("decode location: %v", err)
	}
	if got["type"] != "location" || got["id"] != "location-1" || got["title"] != "Amsterdam" {
		t.Fatalf("unexpected location payload: %#v", got)
	}
	if _, ok := got["input_message_content"].(map[string]any); !ok {
		t.Fatalf("expected input_message_content object: %#v", got)
	}
	if err := validateInlineQueryResult(result); err != nil {
		t.Fatalf("valid location rejected: %v", err)
	}

	tests := []struct {
		name   string
		result InlineQueryResult
	}{
		{name: "empty id", result: InlineQueryResultLocation{Latitude: 52, Longitude: 4, Title: "Location"}},
		{name: "invalid coordinates", result: InlineQueryResultLocation{ID: "location-1", Latitude: 91, Longitude: 4, Title: "Location"}},
		{name: "empty title", result: InlineQueryResultLocation{ID: "location-1", Latitude: 52, Longitude: 4}},
		{name: "negative thumbnail width", result: InlineQueryResultLocation{ID: "location-1", Latitude: 52, Longitude: 4, Title: "Location", ThumbnailWidth: -1}},
		{name: "negative thumbnail height", result: InlineQueryResultLocation{ID: "location-1", Latitude: 52, Longitude: 4, Title: "Location", ThumbnailHeight: -1}},
		{name: "invalid thumbnail url", result: InlineQueryResultLocation{ID: "location-1", Latitude: 52, Longitude: 4, Title: "Location", ThumbnailURL: "ftp://example.com/thumb.jpg"}},
		{name: "invalid reply markup", result: InlineQueryResultLocation{ID: "location-1", Latitude: 52, Longitude: 4, Title: "Location", ReplyMarkup: &telegram.InlineKeyboardMarkup{}}},
		{name: "invalid input content", result: InlineQueryResultLocation{ID: "location-1", Latitude: 52, Longitude: 4, Title: "Location", InputMessageContent: InputTextMessageContent{}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateInlineQueryResult(tt.result); err == nil {
				t.Fatal("expected validation error")
			}
		})
	}
}

func TestInlineQueryResultVenueMarshalAndValidation(t *testing.T) {
	result := InlineVenue("venue-1", 52.3676, 4.9041, "Venue", "Address")
	result.InputMessageContent = InputVenue(52.3676, 4.9041, "Venue", "Address")

	data, err := json.Marshal(result)
	if err != nil {
		t.Fatalf("marshal venue: %v", err)
	}
	var got map[string]any
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("decode venue: %v", err)
	}
	if got["type"] != "venue" || got["id"] != "venue-1" || got["title"] != "Venue" || got["address"] != "Address" {
		t.Fatalf("unexpected venue payload: %#v", got)
	}
	if err := validateInlineQueryResult(result); err != nil {
		t.Fatalf("valid venue rejected: %v", err)
	}

	tests := []struct {
		name   string
		result InlineQueryResult
	}{
		{name: "empty id", result: InlineQueryResultVenue{Latitude: 52, Longitude: 4, Title: "Venue", Address: "Address"}},
		{name: "invalid coordinates", result: InlineQueryResultVenue{ID: "venue-1", Latitude: 52, Longitude: 181, Title: "Venue", Address: "Address"}},
		{name: "empty title", result: InlineQueryResultVenue{ID: "venue-1", Latitude: 52, Longitude: 4, Address: "Address"}},
		{name: "empty address", result: InlineQueryResultVenue{ID: "venue-1", Latitude: 52, Longitude: 4, Title: "Venue"}},
		{name: "negative thumbnail width", result: InlineQueryResultVenue{ID: "venue-1", Latitude: 52, Longitude: 4, Title: "Venue", Address: "Address", ThumbnailWidth: -1}},
		{name: "invalid reply markup", result: InlineQueryResultVenue{ID: "venue-1", Latitude: 52, Longitude: 4, Title: "Venue", Address: "Address", ReplyMarkup: &telegram.InlineKeyboardMarkup{}}},
		{name: "invalid input content", result: InlineQueryResultVenue{ID: "venue-1", Latitude: 52, Longitude: 4, Title: "Venue", Address: "Address", InputMessageContent: InputVenueMessageContent{Title: "Venue", Address: "Address", Latitude: 91}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateInlineQueryResult(tt.result); err == nil {
				t.Fatal("expected validation error")
			}
		})
	}
}

func TestInlineQueryResultContactMarshalAndValidation(t *testing.T) {
	result := InlineContact("contact-1", "+10000000000", "ai-gram")
	result.InputMessageContent = InputContact("+10000000000", "ai-gram")

	data, err := json.Marshal(result)
	if err != nil {
		t.Fatalf("marshal contact: %v", err)
	}
	var got map[string]any
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("decode contact: %v", err)
	}
	if got["type"] != "contact" || got["id"] != "contact-1" || got["phone_number"] != "+10000000000" || got["first_name"] != "ai-gram" {
		t.Fatalf("unexpected contact payload: %#v", got)
	}
	if err := validateInlineQueryResult(result); err != nil {
		t.Fatalf("valid contact rejected: %v", err)
	}

	tests := []struct {
		name   string
		result InlineQueryResult
	}{
		{name: "empty id", result: InlineQueryResultContact{PhoneNumber: "+10000000000", FirstName: "ai-gram"}},
		{name: "empty phone", result: InlineQueryResultContact{ID: "contact-1", FirstName: "ai-gram"}},
		{name: "empty first name", result: InlineQueryResultContact{ID: "contact-1", PhoneNumber: "+10000000000"}},
		{name: "negative thumbnail width", result: InlineQueryResultContact{ID: "contact-1", PhoneNumber: "+10000000000", FirstName: "ai-gram", ThumbnailWidth: -1}},
		{name: "invalid reply markup", result: InlineQueryResultContact{ID: "contact-1", PhoneNumber: "+10000000000", FirstName: "ai-gram", ReplyMarkup: &telegram.InlineKeyboardMarkup{}}},
		{name: "invalid input content", result: InlineQueryResultContact{ID: "contact-1", PhoneNumber: "+10000000000", FirstName: "ai-gram", InputMessageContent: InputContactMessageContent{FirstName: "ai-gram"}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateInlineQueryResult(tt.result); err == nil {
				t.Fatal("expected validation error")
			}
		})
	}
}

func TestInlineQueryResultGameMarshalAndValidation(t *testing.T) {
	markup := telegram.NewInlineKeyboard([]telegram.InlineKeyboardButton{telegram.InlineButtonCallback("Play", "play")})
	result := InlineGame("game-1", "short-name")
	result.ReplyMarkup = &markup

	data, err := json.Marshal(result)
	if err != nil {
		t.Fatalf("marshal game: %v", err)
	}
	var got map[string]any
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("decode game: %v", err)
	}
	if got["type"] != "game" || got["id"] != "game-1" || got["game_short_name"] != "short-name" {
		t.Fatalf("unexpected game payload: %#v", got)
	}
	if err := validateInlineQueryResult(result); err != nil {
		t.Fatalf("valid game rejected: %v", err)
	}

	tests := []struct {
		name   string
		result InlineQueryResult
	}{
		{name: "empty id", result: InlineQueryResultGame{GameShortName: "short-name"}},
		{name: "empty game short name", result: InlineQueryResultGame{ID: "game-1"}},
		{name: "invalid reply markup", result: InlineQueryResultGame{ID: "game-1", GameShortName: "short-name", ReplyMarkup: &telegram.InlineKeyboardMarkup{}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateInlineQueryResult(tt.result); err == nil {
				t.Fatal("expected validation error")
			}
		})
	}
}

func TestAnswerInlineQuerySendsPayloadAndDecodesResult(t *testing.T) {
	const token = "123:secret"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/bot"+token+"/answerInlineQuery" {
			t.Fatalf("unexpected path: %q", r.URL.Path)
		}
		var payload map[string]any
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode payload: %v", err)
		}
		if payload["inline_query_id"] != "inline-query-id" || payload["cache_time"] != float64(10) || payload["is_personal"] != true || payload["next_offset"] != "next" {
			t.Fatalf("unexpected payload: %#v", payload)
		}
		results, ok := payload["results"].([]any)
		if !ok || len(results) != 5 {
			t.Fatalf("unexpected results: %#v", payload["results"])
		}
		article, _ := results[0].(map[string]any)
		content, _ := article["input_message_content"].(map[string]any)
		if article["type"] != "article" || article["id"] != "article-1" || article["title"] != "Article" || content["message_text"] != "hello" {
			t.Fatalf("unexpected article result: %#v", article)
		}
		location, _ := results[1].(map[string]any)
		locationContent, _ := location["input_message_content"].(map[string]any)
		if location["type"] != "location" || location["id"] != "location-1" || locationContent["latitude"] != 52.3676 {
			t.Fatalf("unexpected location result: %#v", location)
		}
		venue, _ := results[2].(map[string]any)
		venueContent, _ := venue["input_message_content"].(map[string]any)
		if venue["type"] != "venue" || venue["id"] != "venue-1" || venueContent["title"] != "Venue" {
			t.Fatalf("unexpected venue result: %#v", venue)
		}
		contact, _ := results[3].(map[string]any)
		contactContent, _ := contact["input_message_content"].(map[string]any)
		if contact["type"] != "contact" || contact["id"] != "contact-1" || contactContent["phone_number"] != "+10000000000" {
			t.Fatalf("unexpected contact result: %#v", contact)
		}
		game, _ := results[4].(map[string]any)
		if game["type"] != "game" || game["id"] != "game-1" || game["game_short_name"] != "game" {
			t.Fatalf("unexpected game result: %#v", game)
		}
		button, _ := payload["button"].(map[string]any)
		if button["text"] != "Open" || button["start_parameter"] != "start_1" {
			t.Fatalf("unexpected button: %#v", button)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true,"result":true}`))
	}))
	defer server.Close()

	bot := newTestBot(t, token, server.URL, server.Client())
	ok, err := bot.AnswerInlineQuery(context.Background(), AnswerInlineQueryParams{
		InlineQueryID: "inline-query-id",
		Results: []InlineQueryResult{
			InlineArticle("article-1", "Article", InputText("hello")),
			func() InlineQueryResultLocation {
				result := InlineLocation("location-1", 52.3676, 4.9041, "Location")
				result.InputMessageContent = InputLocation(52.3676, 4.9041)
				return result
			}(),
			func() InlineQueryResultVenue {
				result := InlineVenue("venue-1", 52.3676, 4.9041, "Venue", "Address")
				result.InputMessageContent = InputVenue(52.3676, 4.9041, "Venue", "Address")
				return result
			}(),
			func() InlineQueryResultContact {
				result := InlineContact("contact-1", "+10000000000", "ai-gram")
				result.InputMessageContent = InputContact("+10000000000", "ai-gram")
				return result
			}(),
			InlineGame("game-1", "game"),
		},
		CacheTime:  10,
		IsPersonal: true,
		NextOffset: "next",
		Button:     &telegram.InlineQueryResultsButton{Text: "Open", StartParameter: "start_1"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Fatal("expected true result")
	}
}

func TestAnswerInlineQueryAllowsEmptyResults(t *testing.T) {
	const token = "123:secret"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/bot"+token+"/answerInlineQuery" {
			t.Fatalf("unexpected path: %q", r.URL.Path)
		}
		var payload map[string]any
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode payload: %v", err)
		}
		results, ok := payload["results"].([]any)
		if !ok || len(results) != 0 {
			t.Fatalf("empty results should be encoded as an array: %#v", payload["results"])
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true,"result":true}`))
	}))
	defer server.Close()

	bot := newTestBot(t, token, server.URL, server.Client())
	ok, err := bot.AnswerInlineQuery(context.Background(), AnswerInlineQueryParams{InlineQueryID: "inline-query-id"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Fatal("expected true result")
	}
}

func TestAnswerInlineQueryValidation(t *testing.T) {
	const token = "123:secret"
	bot := newTestBot(t, token, "https://example.test", nil)
	tooManyResults := make([]InlineQueryResult, 51)
	for i := range tooManyResults {
		tooManyResults[i] = InlineArticle("article", "Article", InputText("hello"))
	}
	tests := []struct {
		name   string
		params AnswerInlineQueryParams
	}{
		{name: "empty inline query id", params: AnswerInlineQueryParams{}},
		{name: "too many results", params: AnswerInlineQueryParams{InlineQueryID: "inline", Results: tooManyResults}},
		{name: "invalid result", params: AnswerInlineQueryParams{InlineQueryID: "inline", Results: []InlineQueryResult{InlineQueryResultArticle{Title: "Article", InputMessageContent: InputText("hello")}}}},
		{name: "negative cache time", params: AnswerInlineQueryParams{InlineQueryID: "inline", CacheTime: -1}},
		{name: "long next offset", params: AnswerInlineQueryParams{InlineQueryID: "inline", NextOffset: string(make([]byte, 65))}},
		{name: "invalid button", params: AnswerInlineQueryParams{InlineQueryID: "inline", Button: &telegram.InlineQueryResultsButton{Text: "Open"}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ok, err := bot.AnswerInlineQuery(context.Background(), tt.params)
			if err == nil {
				t.Fatal("expected error")
			}
			if ok {
				t.Fatal("expected false result")
			}
			assertNoToken(t, err, token)
		})
	}
}

func TestAnswerInlineQueryReturnsAPIError(t *testing.T) {
	const token = "123:secret"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/bot"+token+"/answerInlineQuery" {
			t.Fatalf("unexpected path: %q", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":false,"error_code":400,"description":"Bad Request"}`))
	}))
	defer server.Close()

	bot := newTestBot(t, token, server.URL, server.Client())
	ok, err := bot.AnswerInlineQuery(context.Background(), AnswerInlineQueryParams{InlineQueryID: "inline-query-id"})
	if err == nil {
		t.Fatal("expected error")
	}
	if ok {
		t.Fatal("expected false result")
	}
	var apiErr *apierrors.APIError
	if !stderrors.As(err, &apiErr) {
		t.Fatalf("expected APIError, got %T", err)
	}
	assertNoToken(t, err, token)
}

func TestAnswerInlineQueryResponseAndContextErrors(t *testing.T) {
	const token = "123:secret"
	t.Run("cancelled context", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			t.Fatal("request should not reach server")
		}))
		defer server.Close()
		bot := newTestBot(t, token, server.URL, server.Client())
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		ok, err := bot.AnswerInlineQuery(ctx, AnswerInlineQueryParams{InlineQueryID: "inline-query-id"})
		if err == nil {
			t.Fatal("expected error")
		}
		if ok {
			t.Fatal("expected false result")
		}
		assertNoToken(t, err, token)
	})

	t.Run("invalid json", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/bot"+token+"/answerInlineQuery" {
				t.Fatalf("unexpected path: %q", r.URL.Path)
			}
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`not-json`))
		}))
		defer server.Close()
		bot := newTestBot(t, token, server.URL, server.Client())
		ok, err := bot.AnswerInlineQuery(context.Background(), AnswerInlineQueryParams{InlineQueryID: "inline-query-id"})
		if err == nil {
			t.Fatal("expected error")
		}
		if ok {
			t.Fatal("expected false result")
		}
		assertNoToken(t, err, token)
	})

	t.Run("http status", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/bot"+token+"/answerInlineQuery" {
				t.Fatalf("unexpected path: %q", r.URL.Path)
			}
			http.Error(w, "server error", http.StatusInternalServerError)
		}))
		defer server.Close()
		bot := newTestBot(t, token, server.URL, server.Client())
		ok, err := bot.AnswerInlineQuery(context.Background(), AnswerInlineQueryParams{InlineQueryID: "inline-query-id"})
		if err == nil {
			t.Fatal("expected error")
		}
		if ok {
			t.Fatal("expected false result")
		}
		assertNoToken(t, err, token)
	})
}

type unsupportedInputMessageContent struct{}

func (unsupportedInputMessageContent) inputMessageContent() {}

type unsupportedInlineQueryResult struct{}

func (unsupportedInlineQueryResult) inlineQueryResult() {}
