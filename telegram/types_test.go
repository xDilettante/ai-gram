package telegram

import (
	"encoding/json"
	"testing"
)

func TestUpdateDecodesPracticalMessagePayloads(t *testing.T) {
	payload := []byte(`{
		"update_id": 100,
		"message": {
			"message_id": 10,
			"from": {"id": 1, "is_bot": false, "first_name": "Alice", "username": "alice"},
			"chat": {"id": 2, "type": "private", "first_name": "Alice"},
			"date": 123,
			"text": "/start payload",
			"entities": [{"type": "bot_command", "offset": 0, "length": 6}],
			"caption": "caption text",
			"caption_entities": [{"type": "bold", "offset": 0, "length": 7}],
			"photo": [{"file_id": "photo-id", "file_unique_id": "photo-unique", "width": 100, "height": 80, "file_size": 1234}],
			"document": {"file_id": "doc-id", "file_unique_id": "doc-unique", "file_name": "report.pdf", "mime_type": "application/pdf", "file_size": 2000, "thumbnail": {"file_id": "doc-thumb", "file_unique_id": "doc-thumb-unique", "width": 50, "height": 50}},
			"animation": {"file_id": "anim-id", "file_unique_id": "anim-unique", "width": 320, "height": 240, "duration": 3, "file_name": "a.gif", "mime_type": "image/gif", "file_size": 3000, "thumbnail": {"file_id": "anim-thumb", "file_unique_id": "anim-thumb-unique", "width": 50, "height": 40}},
			"audio": {"file_id": "audio-id", "file_unique_id": "audio-unique", "duration": 60, "performer": "Artist", "title": "Song", "file_name": "song.mp3", "mime_type": "audio/mpeg", "file_size": 4000, "thumbnail": {"file_id": "audio-thumb", "file_unique_id": "audio-thumb-unique", "width": 60, "height": 60}},
			"video": {"file_id": "video-id", "file_unique_id": "video-unique", "width": 640, "height": 480, "duration": 30, "file_name": "video.mp4", "mime_type": "video/mp4", "file_size": 5000, "thumbnail": {"file_id": "video-thumb", "file_unique_id": "video-thumb-unique", "width": 80, "height": 60}},
			"voice": {"file_id": "voice-id", "file_unique_id": "voice-unique", "duration": 5, "mime_type": "audio/ogg", "file_size": 6000},
			"sticker": {"file_id": "sticker-id", "file_unique_id": "sticker-unique", "type": "regular", "width": 512, "height": 512, "is_animated": false, "is_video": true, "emoji": "🙂", "set_name": "fun", "file_size": 7000, "thumbnail": {"file_id": "sticker-thumb", "file_unique_id": "sticker-thumb-unique", "width": 90, "height": 90}},
			"contact": {"phone_number": "+100", "first_name": "Bob", "last_name": "Smith", "user_id": 3, "vcard": "BEGIN:VCARD"},
			"location": {"longitude": 37.6, "latitude": 55.7, "horizontal_accuracy": 12.5, "live_period": 60, "heading": 90, "proximity_alert_radius": 100},
			"venue": {"location": {"longitude": 37.6, "latitude": 55.7}, "title": "Cafe", "address": "Main street", "foursquare_id": "fs", "foursquare_type": "food", "google_place_id": "gp", "google_place_type": "restaurant"}
		},
		"callback_query": {
			"id": "callback-id",
			"from": {"id": 4, "is_bot": false, "first_name": "Carol"},
			"message": {"message_id": 20, "chat": {"id": 5, "type": "private"}, "date": 124, "text": "button message"},
			"inline_message_id": "inline-id",
			"chat_instance": "chat-instance",
			"data": "button:data",
			"game_short_name": "game"
		}
	}`)

	var update Update
	if err := json.Unmarshal(payload, &update); err != nil {
		t.Fatalf("decode update: %v", err)
	}
	message := update.Message
	if update.UpdateID != 100 || message == nil {
		t.Fatalf("unexpected update: %+v", update)
	}
	if len(message.Entities) != 1 || message.Entities[0].Type != EntityBotCommand {
		t.Fatalf("unexpected entities: %#v", message.Entities)
	}
	if message.Caption != "caption text" || len(message.CaptionEntities) != 1 || message.CaptionEntities[0].Type != EntityBold {
		t.Fatalf("unexpected caption data: caption=%q entities=%#v", message.Caption, message.CaptionEntities)
	}
	if len(message.Photo) != 1 || message.Photo[0].FileID != "photo-id" || message.Photo[0].FileSize != 1234 {
		t.Fatalf("unexpected photo: %#v", message.Photo)
	}
	if message.Document == nil || message.Document.FileName != "report.pdf" || message.Document.Thumbnail == nil || message.Document.Thumbnail.FileID != "doc-thumb" {
		t.Fatalf("unexpected document: %+v", message.Document)
	}
	if message.Animation == nil || message.Animation.FileID != "anim-id" || message.Animation.Thumbnail == nil {
		t.Fatalf("unexpected animation: %+v", message.Animation)
	}
	if message.Audio == nil || message.Audio.Performer != "Artist" || message.Audio.Thumbnail == nil {
		t.Fatalf("unexpected audio: %+v", message.Audio)
	}
	if message.Video == nil || message.Video.Width != 640 || message.Video.Thumbnail == nil {
		t.Fatalf("unexpected video: %+v", message.Video)
	}
	if message.Voice == nil || message.Voice.MimeType != "audio/ogg" {
		t.Fatalf("unexpected voice: %+v", message.Voice)
	}
	if message.Sticker == nil || !message.Sticker.IsVideo || message.Sticker.Emoji != "🙂" || message.Sticker.Thumbnail == nil {
		t.Fatalf("unexpected sticker: %+v", message.Sticker)
	}
	if message.Contact == nil || message.Contact.PhoneNumber != "+100" || message.Contact.UserID != 3 {
		t.Fatalf("unexpected contact: %+v", message.Contact)
	}
	if message.Location == nil || message.Location.Longitude != 37.6 || message.Location.Heading != 90 {
		t.Fatalf("unexpected location: %+v", message.Location)
	}
	if message.Venue == nil || message.Venue.Title != "Cafe" || message.Venue.Location.Latitude != 55.7 {
		t.Fatalf("unexpected venue: %+v", message.Venue)
	}
	if update.CallbackQuery == nil || update.CallbackQuery.Message == nil || update.CallbackQuery.Data != "button:data" || update.CallbackQuery.ChatInstance != "chat-instance" || update.CallbackQuery.GameShortName != "game" {
		t.Fatalf("unexpected callback query: %+v", update.CallbackQuery)
	}
}

func TestUpdateDecodesChatJoinRequest(t *testing.T) {
	payload := []byte(`{
		"update_id": 101,
		"chat_join_request": {
			"chat": {"id": -100123, "type": "supergroup", "title": "Test group"},
			"from": {"id": 777, "is_bot": false, "first_name": "Joiner", "username": "joiner"},
			"user_chat_id": 888,
			"date": 1234567890,
			"bio": "hello",
			"invite_link": {
				"invite_link": "https://t.me/+redacted",
				"creator": {"id": 1, "is_bot": true, "first_name": "Bot"},
				"creates_join_request": true,
				"is_primary": false,
				"is_revoked": false,
				"name": "join requests"
			}
		}
	}`)

	var update Update
	if err := json.Unmarshal(payload, &update); err != nil {
		t.Fatalf("decode update: %v", err)
	}
	request := update.ChatJoinRequest
	if update.UpdateID != 101 || request == nil {
		t.Fatalf("unexpected update: %+v", update)
	}
	if request.Chat.ID != -100123 || request.Chat.Type != "supergroup" || request.From.ID != 777 || request.UserChatID != 888 || request.Date != 1234567890 || request.Bio != "hello" {
		t.Fatalf("unexpected join request: %+v", request)
	}
	if request.InviteLink == nil || request.InviteLink.Creator.ID != 1 || !request.InviteLink.CreatesJoinRequest || request.InviteLink.Name != "join requests" {
		t.Fatalf("unexpected invite link: %+v", request.InviteLink)
	}
}

func TestUpdateDecodesInlineQuery(t *testing.T) {
	payload := []byte(`{
		"update_id": 102,
		"inline_query": {
			"id": "inline-query-id",
			"from": {"id": 777, "is_bot": false, "first_name": "Inline"},
			"query": "hello",
			"offset": "next",
			"chat_type": "sender",
			"location": {"longitude": 4.9041, "latitude": 52.3676}
		}
	}`)

	var update Update
	if err := json.Unmarshal(payload, &update); err != nil {
		t.Fatalf("decode update: %v", err)
	}
	query := update.InlineQuery
	if update.UpdateID != 102 || query == nil {
		t.Fatalf("unexpected update: %+v", update)
	}
	if query.ID != "inline-query-id" || query.From.ID != 777 || query.Query != "hello" || query.Offset != "next" || query.ChatType != "sender" {
		t.Fatalf("unexpected inline query: %+v", query)
	}
	if query.Location == nil || query.Location.Latitude != 52.3676 || query.Location.Longitude != 4.9041 {
		t.Fatalf("unexpected inline query location: %+v", query.Location)
	}
}

func TestUpdateDecodesChosenInlineResult(t *testing.T) {
	payload := []byte(`{
		"update_id": 103,
		"chosen_inline_result": {
			"result_id": "article-1",
			"from": {"id": 778, "is_bot": false, "first_name": "Chooser"},
			"location": {"longitude": 4.9041, "latitude": 52.3676},
			"inline_message_id": "inline-message",
			"query": "hello"
		}
	}`)

	var update Update
	if err := json.Unmarshal(payload, &update); err != nil {
		t.Fatalf("decode update: %v", err)
	}
	chosen := update.ChosenInlineResult
	if update.UpdateID != 103 || chosen == nil {
		t.Fatalf("unexpected update: %+v", update)
	}
	if chosen.ResultID != "article-1" || chosen.From.ID != 778 || chosen.InlineMessageID != "inline-message" || chosen.Query != "hello" {
		t.Fatalf("unexpected chosen inline result: %+v", chosen)
	}
	if chosen.Location == nil || chosen.Location.Latitude != 52.3676 || chosen.Location.Longitude != 4.9041 {
		t.Fatalf("unexpected chosen inline result location: %+v", chosen.Location)
	}
}

func TestMessageDecodesForumTopicServiceFields(t *testing.T) {
	payload := []byte(`{
		"message_id": 10,
		"message_thread_id": 777,
		"chat": {"id": -100123, "type": "supergroup", "title": "Forum"},
		"date": 123,
		"forum_topic_created": {"name": "News", "icon_color": 7322096, "icon_custom_emoji_id": "emoji-create"},
		"forum_topic_edited": {"name": "Renamed", "icon_custom_emoji_id": "emoji-edit"},
		"forum_topic_closed": {},
		"forum_topic_reopened": {},
		"general_forum_topic_hidden": {},
		"general_forum_topic_unhidden": {}
	}`)

	var message Message
	if err := json.Unmarshal(payload, &message); err != nil {
		t.Fatalf("decode message: %v", err)
	}
	if message.MessageThreadID != 777 {
		t.Fatalf("unexpected message_thread_id: %d", message.MessageThreadID)
	}
	if message.ForumTopicCreated == nil || message.ForumTopicCreated.Name != "News" || message.ForumTopicCreated.IconColor != 7322096 || message.ForumTopicCreated.IconCustomEmojiID != "emoji-create" {
		t.Fatalf("unexpected forum_topic_created: %+v", message.ForumTopicCreated)
	}
	if message.ForumTopicEdited == nil || message.ForumTopicEdited.Name != "Renamed" || message.ForumTopicEdited.IconCustomEmojiID != "emoji-edit" {
		t.Fatalf("unexpected forum_topic_edited: %+v", message.ForumTopicEdited)
	}
	if message.ForumTopicClosed == nil || message.ForumTopicReopened == nil || message.GeneralForumTopicHidden == nil || message.GeneralForumTopicUnhidden == nil {
		t.Fatalf("expected all forum service message fields: %+v", message)
	}
}

func TestMessageHelpers(t *testing.T) {
	var nilMessage *Message
	if nilMessage.IsText() || nilMessage.IsCommand("start") || nilMessage.Command() != "" || nilMessage.CommandArguments() != "" || nilMessage.HasPhoto() || nilMessage.HasDocument() || nilMessage.HasMedia() {
		t.Fatal("nil message helpers should return zero values")
	}

	if !(&Message{Text: "hello"}).IsText() {
		t.Fatal("expected text message")
	}
	if (&Message{}).IsText() {
		t.Fatal("empty text should not be text")
	}

	commandTests := []struct {
		text string
		cmd  string
		args string
	}{
		{text: "/start", cmd: "start", args: ""},
		{text: "/start payload text", cmd: "start", args: "payload text"},
		{text: "/start@BotName payload", cmd: "start", args: "payload"},
		{text: "/start\npayload", cmd: "start", args: "payload"},
	}
	for _, tt := range commandTests {
		message := &Message{Text: tt.text}
		if got := message.Command(); got != tt.cmd {
			t.Fatalf("Command(%q)=%q, want %q", tt.text, got, tt.cmd)
		}
		if got := message.CommandArguments(); got != tt.args {
			t.Fatalf("CommandArguments(%q)=%q, want %q", tt.text, got, tt.args)
		}
	}

	if !(&Message{Text: "/start"}).IsCommand("start") {
		t.Fatal("expected /start to match")
	}
	if !(&Message{Text: "/start payload"}).IsCommand("start") {
		t.Fatal("expected /start payload to match")
	}
	if !(&Message{Text: "/start@BotName payload"}).IsCommand("start") {
		t.Fatal("expected /start@BotName payload to match")
	}
	if (&Message{Text: "/startx"}).IsCommand("start") {
		t.Fatal("/startx must not match start")
	}
	if (&Message{Text: "/start"}).IsCommand("") {
		t.Fatal("empty command must not match")
	}
	if (&Message{Text: "/start"}).IsCommand("/start") {
		t.Fatal("command with slash must not match")
	}
}

func TestMessageMediaHelpers(t *testing.T) {
	if !(&Message{Photo: []PhotoSize{{FileID: "p"}}}).HasPhoto() {
		t.Fatal("expected photo")
	}
	if !(&Message{Document: &Document{FileID: "d"}}).HasDocument() {
		t.Fatal("expected document")
	}

	mediaMessages := []*Message{
		{Photo: []PhotoSize{{FileID: "p"}}},
		{Document: &Document{FileID: "d"}},
		{Animation: &Animation{FileID: "a"}},
		{Audio: &Audio{FileID: "a"}},
		{Video: &Video{FileID: "v"}},
		{Voice: &Voice{FileID: "v"}},
		{Sticker: &Sticker{FileID: "s"}},
	}
	for _, message := range mediaMessages {
		if !message.HasMedia() {
			t.Fatalf("expected media for %+v", message)
		}
	}
	if (&Message{Text: "plain"}).HasMedia() {
		t.Fatal("plain text should not be media")
	}
}

func TestUpdateHelpers(t *testing.T) {
	var nilUpdate *Update
	if nilUpdate.EffectiveMessage() != nil || nilUpdate.EffectiveChat() != nil || nilUpdate.EffectiveUser() != nil {
		t.Fatal("nil update helpers should return nil")
	}

	from := &User{ID: 1, FirstName: "Alice"}
	message := &Message{MessageID: 10, From: from, Chat: Chat{ID: 20, Type: "private"}, Text: "hello"}
	update := &Update{Message: message}
	if update.EffectiveMessage() != message {
		t.Fatal("expected message as effective message")
	}
	if chat := update.EffectiveChat(); chat == nil || chat.ID != 20 {
		t.Fatalf("unexpected effective chat: %+v", chat)
	}
	if user := update.EffectiveUser(); user != from {
		t.Fatalf("unexpected effective user: %+v", user)
	}

	edited := &Message{MessageID: 11, From: &User{ID: 2, FirstName: "Bob"}, Chat: Chat{ID: 21, Type: "private"}}
	update = &Update{EditedMessage: edited}
	if update.EffectiveMessage() != edited {
		t.Fatal("expected edited message as effective message")
	}
	if user := update.EffectiveUser(); user == nil || user.ID != 2 {
		t.Fatalf("unexpected edited effective user: %+v", user)
	}

	callbackMessage := &Message{MessageID: 12, Chat: Chat{ID: 22, Type: "private"}}
	update = &Update{CallbackQuery: &CallbackQuery{From: User{ID: 3, FirstName: "Carol"}, Message: callbackMessage}}
	if update.EffectiveMessage() != callbackMessage {
		t.Fatal("expected callback message as effective message")
	}
	if chat := update.EffectiveChat(); chat == nil || chat.ID != 22 {
		t.Fatalf("unexpected callback effective chat: %+v", chat)
	}
	if user := update.EffectiveUser(); user == nil || user.ID != 3 {
		t.Fatalf("unexpected callback effective user: %+v", user)
	}

	update = &Update{ChatJoinRequest: &ChatJoinRequest{Chat: Chat{ID: 23, Type: "supergroup"}, From: User{ID: 4, FirstName: "Dave"}}}
	if update.EffectiveMessage() != nil {
		t.Fatal("chat join request should not have an effective message")
	}
	if chat := update.EffectiveChat(); chat == nil || chat.ID != 23 {
		t.Fatalf("unexpected join request effective chat: %+v", chat)
	}
	if user := update.EffectiveUser(); user == nil || user.ID != 4 {
		t.Fatalf("unexpected join request effective user: %+v", user)
	}

	update = &Update{InlineQuery: &InlineQuery{From: User{ID: 5, FirstName: "Eve"}}}
	if update.EffectiveMessage() != nil {
		t.Fatal("inline query should not have an effective message")
	}
	if chat := update.EffectiveChat(); chat != nil {
		t.Fatalf("inline query should not invent an effective chat: %+v", chat)
	}
	if user := update.EffectiveUser(); user == nil || user.ID != 5 {
		t.Fatalf("unexpected inline query effective user: %+v", user)
	}

	update = &Update{ChosenInlineResult: &ChosenInlineResult{From: User{ID: 6, FirstName: "Frank"}}}
	if update.EffectiveMessage() != nil {
		t.Fatal("chosen inline result should not have an effective message")
	}
	if chat := update.EffectiveChat(); chat != nil {
		t.Fatalf("chosen inline result should not invent an effective chat: %+v", chat)
	}
	if user := update.EffectiveUser(); user == nil || user.ID != 6 {
		t.Fatalf("unexpected chosen inline result effective user: %+v", user)
	}
}

func TestPaymentMessageDecoding(t *testing.T) {
	message := mustDecodeMessage(t, `{
		"message_id": 100,
		"date": 1,
		"chat": {"id": 123, "type": "private"},
		"invoice": {"title": "Invoice", "description": "Description", "start_parameter": "start", "currency": "XTR", "total_amount": 150},
		"successful_payment": {
			"currency": "XTR",
			"total_amount": 150,
			"invoice_payload": "payload",
			"subscription_expiration_date": 1893456000,
			"is_recurring": true,
			"is_first_recurring": true,
			"shipping_option_id": "standard",
			"order_info": {"name": "Alice", "shipping_address": {"country_code": "US", "state": "CA", "city": "SF", "street_line1": "1 Main", "street_line2": "2", "post_code": "94105"}},
			"telegram_payment_charge_id": "tg-charge",
			"provider_payment_charge_id": "provider-charge"
		},
		"refunded_payment": {
			"currency": "XTR",
			"total_amount": 50,
			"invoice_payload": "payload",
			"telegram_payment_charge_id": "tg-refund",
			"provider_payment_charge_id": "provider-refund"
		}
	}`)
	if message.Invoice == nil || message.Invoice.TotalAmount != 150 || message.Invoice.StartParameter != "start" {
		t.Fatalf("unexpected invoice: %+v", message.Invoice)
	}
	if message.SuccessfulPayment == nil || !message.SuccessfulPayment.IsRecurring || message.SuccessfulPayment.OrderInfo == nil || message.SuccessfulPayment.OrderInfo.ShippingAddress == nil {
		t.Fatalf("unexpected successful payment: %+v", message.SuccessfulPayment)
	}
	if message.RefundedPayment == nil || message.RefundedPayment.ProviderPaymentChargeID != "provider-refund" {
		t.Fatalf("unexpected refunded payment: %+v", message.RefundedPayment)
	}
}

func TestPaymentUpdateDecodingAndEffectiveUser(t *testing.T) {
	shipping := mustDecodeUpdate(t, `{
		"update_id": 1,
		"shipping_query": {
			"id": "ship-id",
			"from": {"id": 7, "is_bot": false, "first_name": "Alice"},
			"invoice_payload": "payload",
			"shipping_address": {"country_code": "US", "state": "CA", "city": "SF", "street_line1": "1 Main", "street_line2": "2", "post_code": "94105"}
		}
	}`)
	if shipping.ShippingQuery == nil || shipping.ShippingQuery.ShippingAddress.CountryCode != "US" {
		t.Fatalf("unexpected shipping query: %+v", shipping.ShippingQuery)
	}
	if chat := shipping.EffectiveChat(); chat != nil {
		t.Fatalf("shipping query should not invent an effective chat: %+v", chat)
	}
	if user := shipping.EffectiveUser(); user == nil || user.ID != 7 {
		t.Fatalf("unexpected shipping effective user: %+v", user)
	}

	preCheckout := mustDecodeUpdate(t, `{
		"update_id": 2,
		"pre_checkout_query": {
			"id": "pre-id",
			"from": {"id": 8, "is_bot": false, "first_name": "Bob"},
			"currency": "XTR",
			"total_amount": 150,
			"invoice_payload": "payload",
			"shipping_option_id": "standard",
			"order_info": {"email": "user@example.test"}
		}
	}`)
	if preCheckout.PreCheckoutQuery == nil || preCheckout.PreCheckoutQuery.TotalAmount != 150 || preCheckout.PreCheckoutQuery.OrderInfo == nil {
		t.Fatalf("unexpected pre-checkout query: %+v", preCheckout.PreCheckoutQuery)
	}
	if chat := preCheckout.EffectiveChat(); chat != nil {
		t.Fatalf("pre-checkout query should not invent an effective chat: %+v", chat)
	}
	if user := preCheckout.EffectiveUser(); user == nil || user.ID != 8 {
		t.Fatalf("unexpected pre-checkout effective user: %+v", user)
	}
}

func mustDecodeMessage(t *testing.T, payload string) Message {
	t.Helper()
	var message Message
	if err := json.Unmarshal([]byte(payload), &message); err != nil {
		t.Fatalf("decode message: %v", err)
	}
	return message
}

func mustDecodeUpdate(t *testing.T, payload string) Update {
	t.Helper()
	var update Update
	if err := json.Unmarshal([]byte(payload), &update); err != nil {
		t.Fatalf("decode update: %v", err)
	}
	return update
}
