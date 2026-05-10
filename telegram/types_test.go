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
			"live_photo": {"photo": [{"file_id": "live-photo-id", "file_unique_id": "live-photo-unique", "width": 640, "height": 480}], "file_id": "live-id", "file_unique_id": "live-unique", "width": 640, "height": 480, "duration": 3, "mime_type": "video/mp4", "file_size": 8000},
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
	if message.LivePhoto == nil || message.LivePhoto.FileID != "live-id" || len(message.LivePhoto.Photo) != 1 || message.LivePhoto.Photo[0].FileID != "live-photo-id" {
		t.Fatalf("unexpected live photo: %+v", message.LivePhoto)
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
	if update.CallbackQuery == nil || update.CallbackQuery.Message == nil || !update.CallbackQuery.Message.IsAccessible() || update.CallbackQuery.Data != "button:data" || update.CallbackQuery.ChatInstance != "chat-instance" || update.CallbackQuery.GameShortName != "game" {
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
	update = &Update{CallbackQuery: &CallbackQuery{From: User{ID: 3, FirstName: "Carol"}, Message: &MaybeInaccessibleMessage{Message: callbackMessage, MessageID: callbackMessage.MessageID, Chat: callbackMessage.Chat}}}
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

	channelPost := &Message{MessageID: 13, Chat: Chat{ID: -100, Type: "channel"}, Text: "post"}
	update = &Update{ChannelPost: channelPost}
	if update.EffectiveMessage() != channelPost {
		t.Fatal("expected channel post as effective message")
	}
	if chat := update.EffectiveChat(); chat == nil || chat.ID != -100 {
		t.Fatalf("unexpected channel post effective chat: %+v", chat)
	}
	if user := update.EffectiveUser(); user != nil {
		t.Fatalf("channel post should not invent effective user: %+v", user)
	}

	editedChannelPost := &Message{MessageID: 14, Chat: Chat{ID: -101, Type: "channel"}, Text: "edited"}
	update = &Update{EditedChannelPost: editedChannelPost}
	if update.EffectiveMessage() != editedChannelPost {
		t.Fatal("expected edited channel post as effective message")
	}
	if chat := update.EffectiveChat(); chat == nil || chat.ID != -101 {
		t.Fatalf("unexpected edited channel post effective chat: %+v", chat)
	}

	guest := &Message{
		MessageID:          15,
		Chat:               Chat{ID: -200, Type: "supergroup"},
		GuestBotCallerUser: &User{ID: 77, FirstName: "Guest caller"},
		GuestQueryID:       "guest-query-id",
	}
	update = &Update{GuestMessage: guest}
	if update.EffectiveMessage() != guest {
		t.Fatal("expected guest message as effective message")
	}
	if chat := update.EffectiveChat(); chat == nil || chat.ID != -200 {
		t.Fatalf("unexpected guest effective chat: %+v", chat)
	}
	if user := update.EffectiveUser(); user == nil || user.ID != 77 {
		t.Fatalf("unexpected guest effective user: %+v", user)
	}
}

func TestUpdateDecodesGuestMessage(t *testing.T) {
	payload := []byte(`{
		"update_id": 900,
		"guest_message": {
			"message_id": 7,
			"chat": {"id": -200, "type": "supergroup", "title": "Guest chat"},
			"date": 1234567890,
			"text": "hello",
			"guest_bot_caller_user": {"id": 77, "is_bot": false, "first_name": "Caller"},
			"guest_bot_caller_chat": {"id": -300, "type": "channel", "title": "Caller chat"},
			"guest_query_id": "guest-query-id"
		}
	}`)

	var update Update
	if err := json.Unmarshal(payload, &update); err != nil {
		t.Fatalf("decode guest update: %v", err)
	}
	message := update.GuestMessage
	if update.UpdateID != 900 || message == nil || message.Text != "hello" || message.GuestQueryID != "guest-query-id" {
		t.Fatalf("unexpected guest update: %+v", update)
	}
	if message.GuestBotCallerUser == nil || message.GuestBotCallerUser.ID != 77 {
		t.Fatalf("unexpected guest caller user: %+v", message.GuestBotCallerUser)
	}
	if message.GuestBotCallerChat == nil || message.GuestBotCallerChat.ID != -300 {
		t.Fatalf("unexpected guest caller chat: %+v", message.GuestBotCallerChat)
	}
}

func TestChatPermissionsDecodesCanReactToMessages(t *testing.T) {
	var permissions ChatPermissions
	if err := json.Unmarshal([]byte(`{"can_send_messages":true,"can_react_to_messages":true}`), &permissions); err != nil {
		t.Fatalf("decode permissions: %v", err)
	}
	if !permissions.CanSendMessages || !permissions.CanReactToMessages {
		t.Fatalf("unexpected permissions: %+v", permissions)
	}
}

func TestUserAndChatMetadataDecode(t *testing.T) {
	var user User
	if err := json.Unmarshal([]byte(`{
		"id": 42,
		"is_bot": true,
		"first_name": "Bot",
		"language_code": "en",
		"is_premium": true,
		"added_to_attachment_menu": true,
		"can_join_groups": true,
		"can_read_all_group_messages": true,
		"supports_inline_queries": true,
		"can_connect_to_business": true,
		"has_main_web_app": true,
		"has_topics_enabled": true,
		"allows_users_to_create_topics": true,
		"can_manage_bots": true
	}`), &user); err != nil {
		t.Fatalf("decode user: %v", err)
	}
	if user.ID != 42 || !user.IsBot || user.LanguageCode != "en" || !user.IsPremium || !user.AddedToAttachmentMenu || !user.CanJoinGroups || !user.CanReadAllGroupMessages || !user.SupportsInlineQueries || !user.CanConnectToBusiness || !user.HasMainWebApp || !user.HasTopicsEnabled || !user.AllowsUsersToCreateTopics || !user.CanManageBots {
		t.Fatalf("unexpected user metadata: %+v", user)
	}

	var chat Chat
	if err := json.Unmarshal([]byte(`{"id":-100,"type":"supergroup","title":"Forum","is_forum":true,"is_direct_messages":true}`), &chat); err != nil {
		t.Fatalf("decode chat: %v", err)
	}
	if chat.ID != -100 || chat.Type != "supergroup" || chat.Title != "Forum" || !chat.IsForum || !chat.IsDirectMessages {
		t.Fatalf("unexpected chat metadata: %+v", chat)
	}
}

func TestChatFullInfoDecode(t *testing.T) {
	var info ChatFullInfo
	if err := json.Unmarshal([]byte(`{
		"id": -100,
		"type": "supergroup",
		"title": "Full Chat",
		"username": "fullchat",
		"is_forum": true,
		"is_direct_messages": true,
		"accent_color_id": 7,
		"max_reaction_count": 11,
		"photo": {
			"small_file_id": "small",
			"small_file_unique_id": "small-uniq",
			"big_file_id": "big",
			"big_file_unique_id": "big-uniq"
		},
		"active_usernames": ["fullchat", "fullchat2"],
		"birthdate": {"day": 1, "month": 2, "year": 2000},
		"business_intro": {"title": "Intro", "message": "Welcome"},
		"business_location": {"address": "Main street", "location": {"longitude": 37.6, "latitude": 55.7}},
		"business_opening_hours": {"time_zone_name": "Europe/Moscow", "opening_hours": [{"opening_minute": 60, "closing_minute": 120}]},
		"personal_chat": {"id": -200, "type": "channel", "title": "Personal"},
		"parent_chat": {"id": -201, "type": "channel", "title": "Parent"},
		"available_reactions": [{"type": "emoji", "emoji": "👍"}, {"type": "custom_emoji", "custom_emoji_id": "custom"}],
		"background_custom_emoji_id": "bg",
		"profile_accent_color_id": 8,
		"profile_background_custom_emoji_id": "profile-bg",
		"emoji_status_custom_emoji_id": "status",
		"emoji_status_expiration_date": 1893456000,
		"bio": "bio",
		"has_private_forwards": true,
		"has_restricted_voice_and_video_messages": true,
		"join_to_send_messages": true,
		"join_by_request": true,
		"description": "description",
		"invite_link": "https://t.me/+redacted",
		"pinned_message": {"message_id": 9, "chat": {"id": -100, "type": "supergroup"}, "date": 1, "text": "Pinned"},
		"permissions": {"can_send_messages": true, "can_edit_tag": true},
		"accepted_gift_types": {"unlimited_gifts": true, "limited_gifts": true, "unique_gifts": true, "premium_subscription": true, "gifts_from_channels": true},
		"can_send_paid_media": true,
		"slow_mode_delay": 10,
		"unrestrict_boost_count": 2,
		"message_auto_delete_time": 86400,
		"has_aggressive_anti_spam_enabled": true,
		"has_hidden_members": true,
		"has_protected_content": true,
		"has_visible_history": true,
		"sticker_set_name": "stickers",
		"can_set_sticker_set": true,
		"custom_emoji_sticker_set_name": "emoji",
		"linked_chat_id": -300,
		"location": {"location": {"longitude": 37.7, "latitude": 55.8}, "address": "Address"},
		"rating": {"level": 3, "rating": 100, "current_level_rating": 50, "next_level_rating": 150},
		"first_profile_audio": {"file_id": "audio", "file_unique_id": "audio-uniq", "duration": 12},
		"unique_gift_colors": {
			"model_custom_emoji_id": "model",
			"symbol_custom_emoji_id": "symbol",
			"light_theme_main_color": 1,
			"light_theme_other_colors": [2, 3],
			"dark_theme_main_color": 4,
			"dark_theme_other_colors": [5, 6]
		},
		"paid_message_star_count": 15
	}`), &info); err != nil {
		t.Fatalf("decode chat full info: %v", err)
	}
	if info.ID != -100 || info.Type != "supergroup" || info.Title != "Full Chat" || !info.IsForum || !info.IsDirectMessages || info.AccentColorID != 7 || info.MaxReactionCount != 11 {
		t.Fatalf("unexpected core chat full info: %+v", info)
	}
	if info.Photo == nil || info.Photo.SmallFileID != "small" || len(info.ActiveUsernames) != 2 || info.Birthdate == nil || info.Birthdate.Year != 2000 {
		t.Fatalf("unexpected profile chat full info: %+v", info)
	}
	if info.BusinessIntro == nil || info.BusinessIntro.Title != "Intro" || info.BusinessLocation == nil || info.BusinessLocation.Address != "Main street" || info.BusinessOpeningHours == nil || len(info.BusinessOpeningHours.OpeningHours) != 1 {
		t.Fatalf("unexpected business chat full info: %+v", info)
	}
	if info.PersonalChat == nil || info.PersonalChat.ID != -200 || info.ParentChat == nil || info.ParentChat.ID != -201 {
		t.Fatalf("unexpected related chats: %+v", info)
	}
	if len(info.AvailableReactions) != 2 {
		t.Fatalf("unexpected available reactions: %+v", info.AvailableReactions)
	}
	if info.PinnedMessage == nil || info.PinnedMessage.MessageID != 9 || info.Permissions == nil || !info.Permissions.CanEditTag || !info.AcceptedGiftTypes.GiftsFromChannels {
		t.Fatalf("unexpected permissions or pinned message: %+v", info)
	}
	if info.Location == nil || info.Location.Address != "Address" || info.Rating == nil || info.Rating.Level != 3 || info.FirstProfileAudio == nil || info.FirstProfileAudio.FileID != "audio" || info.UniqueGiftColors == nil || info.PaidMessageStarCount != 15 {
		t.Fatalf("unexpected remaining chat full info fields: %+v", info)
	}
}

func TestUpdateDecodesChannelPostsAndPoll(t *testing.T) {
	var channelPost Update
	if err := json.Unmarshal([]byte(`{"update_id":1,"channel_post":{"message_id":10,"chat":{"id":-100,"type":"channel","title":"Channel"},"date":1,"text":"post"}}`), &channelPost); err != nil {
		t.Fatalf("decode channel post update: %v", err)
	}
	if channelPost.ChannelPost == nil || channelPost.ChannelPost.Chat.Type != "channel" || channelPost.ChannelPost.Text != "post" {
		t.Fatalf("unexpected channel post update: %+v", channelPost)
	}

	var editedChannelPost Update
	if err := json.Unmarshal([]byte(`{"update_id":2,"edited_channel_post":{"message_id":11,"chat":{"id":-100,"type":"channel","title":"Channel"},"date":1,"edit_date":2,"text":"edited"}}`), &editedChannelPost); err != nil {
		t.Fatalf("decode edited channel post update: %v", err)
	}
	if editedChannelPost.EditedChannelPost == nil || editedChannelPost.EditedChannelPost.EditDate != 2 || editedChannelPost.EditedChannelPost.Text != "edited" {
		t.Fatalf("unexpected edited channel post update: %+v", editedChannelPost)
	}

	var pollUpdate Update
	if err := json.Unmarshal([]byte(`{"update_id":3,"poll":{"id":"poll-id","question":"Question?","options":[{"text":"A","voter_count":1}],"total_voter_count":1,"is_closed":false,"is_anonymous":true,"type":"regular","allows_multiple_answers":false}}`), &pollUpdate); err != nil {
		t.Fatalf("decode poll update: %v", err)
	}
	if pollUpdate.Poll == nil || pollUpdate.Poll.ID != "poll-id" || pollUpdate.Poll.Question != "Question?" {
		t.Fatalf("unexpected poll update: %+v", pollUpdate)
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

func TestPoll96FieldsDecode(t *testing.T) {
	message := mustDecodeMessage(t, `{
		"message_id": 200,
		"date": 1,
		"chat": {"id": 123, "type": "private"},
		"poll": {
			"id": "poll-id",
			"question": "Pick",
			"options": [{
				"persistent_id": "option-a",
				"text": "A",
				"text_entities": [{"type": "custom_emoji", "offset": 0, "length": 1, "custom_emoji_id": "emoji"}],
				"voter_count": 2,
				"added_by_user": {"id": 7, "is_bot": false, "first_name": "Alice"},
				"added_by_chat": {"id": -100, "type": "supergroup", "title": "Poll chat"},
				"addition_date": 1710000000
			}],
			"total_voter_count": 2,
			"is_closed": false,
			"is_anonymous": false,
			"type": "quiz",
			"allows_multiple_answers": true,
			"allows_revoting": true,
			"correct_option_ids": [0],
			"description": "Details",
			"description_entities": [{"type": "bold", "offset": 0, "length": 7}]
		}
	}`)
	if message.Poll == nil || len(message.Poll.CorrectOptionIDs) != 1 || message.Poll.CorrectOptionIDs[0] != 0 || !message.Poll.AllowsRevoting {
		t.Fatalf("unexpected poll 9.6 fields: %+v", message.Poll)
	}
	option := message.Poll.Options[0]
	if option.PersistentID != "option-a" || option.AddedByUser == nil || option.AddedByUser.ID != 7 || option.AddedByChat == nil || option.AddedByChat.ID != -100 || option.AdditionDate != 1710000000 || len(option.TextEntities) != 1 {
		t.Fatalf("unexpected poll option 9.6 fields: %+v", option)
	}
	if message.Poll.Description != "Details" || len(message.Poll.DescriptionEntities) != 1 || message.Poll.DescriptionEntities[0].Type != EntityBold {
		t.Fatalf("unexpected poll description fields: %+v", message.Poll)
	}
}

func TestPoll10MediaFieldsDecode(t *testing.T) {
	message := mustDecodeMessage(t, `{
		"message_id": 201,
		"date": 1,
		"chat": {"id": 123, "type": "private"},
		"poll": {
			"id": "poll-id",
			"question": "Pick",
			"options": [{
				"persistent_id": "option-a",
				"text": "A",
				"media": {
					"sticker": {
						"file_id": "sticker-file",
						"file_unique_id": "sticker-unique",
						"type": "regular",
						"width": 512,
						"height": 512,
						"is_animated": false,
						"is_video": false
					}
				},
				"voter_count": 2
			}],
			"total_voter_count": 2,
			"is_closed": false,
			"is_anonymous": false,
			"type": "quiz",
			"allows_multiple_answers": true,
			"allows_revoting": true,
			"members_only": true,
			"country_codes": ["US", "DE"],
			"correct_option_ids": [0],
			"explanation": "Because",
			"explanation_media": {
				"live_photo": {
					"photo": [{"file_id": "photo-file", "file_unique_id": "photo-unique", "width": 640, "height": 480}],
					"file_id": "live-file",
					"file_unique_id": "live-unique",
					"width": 640,
					"height": 480,
					"duration": 3,
					"mime_type": "video/mp4",
					"file_size": 123456
				}
			},
			"description": "Details",
			"media": {
				"photo": [{"file_id": "poll-photo", "file_unique_id": "poll-photo-unique", "width": 800, "height": 600}]
			}
		}
	}`)

	poll := message.Poll
	if poll == nil || !poll.MembersOnly || len(poll.CountryCodes) != 2 || poll.CountryCodes[0] != "US" || poll.CountryCodes[1] != "DE" {
		t.Fatalf("unexpected poll 10.0 fields: %+v", poll)
	}
	if poll.Media == nil || len(poll.Media.Photo) != 1 || poll.Media.Photo[0].FileID != "poll-photo" {
		t.Fatalf("unexpected poll media: %+v", poll.Media)
	}
	if poll.ExplanationMedia == nil || poll.ExplanationMedia.LivePhoto == nil || poll.ExplanationMedia.LivePhoto.FileID != "live-file" || len(poll.ExplanationMedia.LivePhoto.Photo) != 1 {
		t.Fatalf("unexpected explanation media: %+v", poll.ExplanationMedia)
	}
	option := poll.Options[0]
	if option.Media == nil || option.Media.Sticker == nil || option.Media.Sticker.FileID != "sticker-file" {
		t.Fatalf("unexpected poll option media: %+v", option.Media)
	}
}

func TestPollAnswerUpdateDecodeAndEffectiveHelpers(t *testing.T) {
	update := mustDecodeUpdate(t, `{
		"update_id": 201,
		"poll_answer": {
			"poll_id": "poll-id",
			"voter_chat": {"id": -100123, "type": "supergroup", "title": "Anonymous voters"},
			"user": {"id": 42, "is_bot": false, "first_name": "Alice"},
			"option_ids": [0, 2],
			"option_persistent_ids": ["option-a", "option-c"]
		}
	}`)
	answer := update.PollAnswer
	if answer == nil || answer.PollID != "poll-id" || len(answer.OptionIDs) != 2 || len(answer.OptionPersistentIDs) != 2 {
		t.Fatalf("unexpected poll answer: %+v", answer)
	}
	if user := update.EffectiveUser(); user == nil || user.ID != 42 {
		t.Fatalf("unexpected poll answer effective user: %+v", user)
	}
	if chat := update.EffectiveChat(); chat == nil || chat.ID != -100123 {
		t.Fatalf("unexpected poll answer effective chat: %+v", chat)
	}
}

func TestMessageDecodesPollOptionServiceFields(t *testing.T) {
	message := mustDecodeMessage(t, `{
		"message_id": 202,
		"date": 1,
		"chat": {"id": 123, "type": "private"},
		"reply_to_poll_option_id": "option-a",
		"poll_option_added": {
			"poll_message": {"message_id": 99, "chat": {"id": 123, "type": "private"}, "date": 1},
			"option_persistent_id": "option-a",
			"option_text": "Added",
			"option_text_entities": [{"type": "bold", "offset": 0, "length": 5}]
		},
		"poll_option_deleted": {
			"poll_message": {"message_id": 99, "chat": {"id": 123, "type": "private"}, "date": 1},
			"option_persistent_id": "option-b",
			"option_text": "Deleted",
			"option_text_entities": [{"type": "italic", "offset": 0, "length": 7}]
		}
	}`)
	if message.ReplyToPollOptionID != "option-a" {
		t.Fatalf("unexpected reply_to_poll_option_id: %q", message.ReplyToPollOptionID)
	}
	if message.PollOptionAdded == nil || message.PollOptionAdded.OptionPersistentID != "option-a" || message.PollOptionAdded.PollMessage == nil || len(message.PollOptionAdded.OptionTextEntities) != 1 {
		t.Fatalf("unexpected poll_option_added: %+v", message.PollOptionAdded)
	}
	if message.PollOptionDeleted == nil || message.PollOptionDeleted.OptionPersistentID != "option-b" || message.PollOptionDeleted.PollMessage == nil || len(message.PollOptionDeleted.OptionTextEntities) != 1 {
		t.Fatalf("unexpected poll_option_deleted: %+v", message.PollOptionDeleted)
	}
}
