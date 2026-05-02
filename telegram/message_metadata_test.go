package telegram

import (
	"encoding/json"
	"testing"
)

func TestMessageOriginUnmarshalVariants(t *testing.T) {
	tests := []struct {
		name string
		body string
		want func(MessageOrigin) bool
	}{
		{name: "user", body: `{"type":"user","date":10,"sender_user":{"id":1,"is_bot":false,"first_name":"Ada"}}`, want: func(origin MessageOrigin) bool {
			value, ok := origin.(MessageOriginUser)
			return ok && value.Type == "user" && value.SenderUser.ID == 1
		}},
		{name: "hidden user", body: `{"type":"hidden_user","date":11,"sender_user_name":"Hidden"}`, want: func(origin MessageOrigin) bool {
			value, ok := origin.(MessageOriginHiddenUser)
			return ok && value.Type == "hidden_user" && value.SenderUserName == "Hidden"
		}},
		{name: "chat", body: `{"type":"chat","date":12,"sender_chat":{"id":-100,"type":"supergroup","title":"Group"},"author_signature":"admin"}`, want: func(origin MessageOrigin) bool {
			value, ok := origin.(MessageOriginChat)
			return ok && value.Type == "chat" && value.SenderChat.ID == -100 && value.AuthorSignature == "admin"
		}},
		{name: "channel", body: `{"type":"channel","date":13,"chat":{"id":-200,"type":"channel","title":"News"},"message_id":55,"author_signature":"editor"}`, want: func(origin MessageOrigin) bool {
			value, ok := origin.(MessageOriginChannel)
			return ok && value.Type == "channel" && value.Chat.ID == -200 && value.MessageID == 55
		}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			origin, err := UnmarshalMessageOrigin([]byte(tt.body))
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !tt.want(origin) {
				t.Fatalf("unexpected origin: %#v", origin)
			}
		})
	}
}

func TestMessageOriginUnmarshalUnknownType(t *testing.T) {
	if _, err := UnmarshalMessageOrigin([]byte(`{"type":"unknown","date":1}`)); err == nil {
		t.Fatal("expected unsupported message origin type error")
	}
}

func TestMessageDecodesReplyMetadata(t *testing.T) {
	var message Message
	body := []byte(`{
		"message_id":101,
		"message_thread_id":7,
		"direct_messages_topic":{"topic_id":1234567890123,"user":{"id":9,"is_bot":false,"first_name":"Topic"}},
		"chat":{"id":123,"type":"private"},
		"date":1000,
		"sender_chat":{"id":-100,"type":"supergroup","title":"Group"},
		"sender_boost_count":2,
		"sender_tag":"mod",
		"forward_origin":{"type":"user","date":900,"sender_user":{"id":5,"is_bot":false,"first_name":"Ada"}},
		"is_topic_message":true,
		"is_automatic_forward":true,
		"reply_to_message":{"message_id":100,"chat":{"id":123,"type":"private"},"date":999,"text":"original"},
		"external_reply":{"origin":{"type":"channel","date":901,"chat":{"id":-200,"type":"channel","title":"News"},"message_id":77},"chat":{"id":-200,"type":"channel","title":"News"},"message_id":77,"photo":[{"file_id":"p","file_unique_id":"pu","width":10,"height":10}],"has_media_spoiler":true},
		"quote":{"text":"quoted","position":4,"is_manual":true,"entities":[{"type":"bold","offset":0,"length":6}]},
		"reply_to_story":{"chat":{"id":123,"type":"private"},"id":44},
		"via_bot":{"id":10,"is_bot":true,"first_name":"Helper"},
		"edit_date":1001,
		"has_protected_content":true,
		"is_paid_post":true,
		"media_group_id":"album-1",
		"author_signature":"author",
		"paid_star_count":3,
		"link_preview_options":{"url":"https://example.test","show_above_text":true},
		"suggested_post_info":{"state":"pending","price":{"currency":"XTR","amount":5},"send_date":2000},
		"effect_id":"effect-1",
		"show_caption_above_media":true,
		"has_media_spoiler":true,
		"reply_to_checklist_task_id":8,
		"reply_to_poll_option_id":"poll-option-a",
		"pinned_message":{"message_id":99,"chat":{"id":123,"type":"private"},"date":0},
		"reply_markup":{"inline_keyboard":[[{"text":"OK","callback_data":"ok"}]]}
	}`)
	if err := json.Unmarshal(body, &message); err != nil {
		t.Fatalf("decode message: %v", err)
	}
	if message.DirectMessagesTopic == nil || message.DirectMessagesTopic.TopicID != 1234567890123 || message.DirectMessagesTopic.User.ID != 9 {
		t.Fatalf("unexpected direct messages topic: %#v", message.DirectMessagesTopic)
	}
	if message.SenderChat == nil || message.SenderChat.ID != -100 || message.SenderBoostCount != 2 || message.SenderTag != "mod" {
		t.Fatalf("unexpected sender metadata: %#v", message)
	}
	if _, ok := message.ForwardOrigin.(MessageOriginUser); !ok {
		t.Fatalf("unexpected forward origin: %#v", message.ForwardOrigin)
	}
	if message.ReplyToMessage == nil || message.ReplyToMessage.Text != "original" {
		t.Fatalf("unexpected reply_to_message: %#v", message.ReplyToMessage)
	}
	if message.ExternalReply == nil || message.ExternalReply.MessageID != 77 || len(message.ExternalReply.Photo) != 1 || !message.ExternalReply.HasMediaSpoiler {
		t.Fatalf("unexpected external reply: %#v", message.ExternalReply)
	}
	if _, ok := message.ExternalReply.Origin.(MessageOriginChannel); !ok {
		t.Fatalf("unexpected external origin: %#v", message.ExternalReply.Origin)
	}
	if message.Quote == nil || message.Quote.Text != "quoted" || message.Quote.Position != 4 || !message.Quote.IsManual || len(message.Quote.Entities) != 1 {
		t.Fatalf("unexpected quote: %#v", message.Quote)
	}
	if message.ReplyToStory == nil || message.ReplyToStory.ID != 44 || message.ViaBot == nil || message.ViaBot.ID != 10 {
		t.Fatalf("unexpected story/via bot metadata: %#v", message)
	}
	if !message.IsTopicMessage || !message.IsAutomaticForward || !message.HasProtectedContent || !message.IsPaidPost {
		t.Fatalf("unexpected boolean metadata: %#v", message)
	}
	if message.EditDate != 1001 || message.MediaGroupID != "album-1" || message.AuthorSignature != "author" || message.PaidStarCount != 3 || message.EffectID != "effect-1" {
		t.Fatalf("unexpected scalar metadata: %#v", message)
	}
	if message.LinkPreviewOptions == nil || message.LinkPreviewOptions.URL != "https://example.test" || !message.LinkPreviewOptions.ShowAboveText {
		t.Fatalf("unexpected link preview options: %#v", message.LinkPreviewOptions)
	}
	if message.SuggestedPostInfo == nil || message.SuggestedPostInfo.State != "pending" || message.SuggestedPostInfo.Price.Amount != 5 || message.SuggestedPostInfo.SendDate != 2000 {
		t.Fatalf("unexpected suggested post info: %#v", message.SuggestedPostInfo)
	}
	if !message.ShowCaptionAboveMedia || !message.HasMediaSpoiler {
		t.Fatalf("unexpected media display metadata: %#v", message)
	}
	if message.ReplyToChecklistTaskID != 8 || message.ReplyToPollOptionID != "poll-option-a" {
		t.Fatalf("unexpected reply ids: checklist=%d poll=%q", message.ReplyToChecklistTaskID, message.ReplyToPollOptionID)
	}
	if message.PinnedMessage == nil || message.PinnedMessage.InaccessibleMessage == nil || message.PinnedMessage.Message != nil || message.PinnedMessage.MessageID != 99 {
		t.Fatalf("unexpected pinned message: %#v", message.PinnedMessage)
	}
	if message.ReplyMarkup == nil || len(message.ReplyMarkup.InlineKeyboard) != 1 || message.ReplyMarkup.InlineKeyboard[0][0].CallbackData != "ok" {
		t.Fatalf("unexpected reply markup: %#v", message.ReplyMarkup)
	}
}

func TestMaybeInaccessibleMessageDecodesAccessibleAndInaccessible(t *testing.T) {
	var accessible MaybeInaccessibleMessage
	if err := json.Unmarshal([]byte(`{"message_id":5,"chat":{"id":1,"type":"private"},"date":100,"text":"ok"}`), &accessible); err != nil {
		t.Fatalf("decode accessible: %v", err)
	}
	if accessible.Message == nil || accessible.InaccessibleMessage != nil || accessible.Message.Text != "ok" {
		t.Fatalf("unexpected accessible message: %#v", accessible)
	}

	var inaccessible MaybeInaccessibleMessage
	if err := json.Unmarshal([]byte(`{"message_id":6,"chat":{"id":1,"type":"private"},"date":0}`), &inaccessible); err != nil {
		t.Fatalf("decode inaccessible: %v", err)
	}
	if inaccessible.InaccessibleMessage == nil || inaccessible.Message != nil || inaccessible.MessageID != 6 {
		t.Fatalf("unexpected inaccessible message: %#v", inaccessible)
	}
}

func TestCallbackQueryDecodesMaybeInaccessibleMessage(t *testing.T) {
	var accessible CallbackQuery
	if err := json.Unmarshal([]byte(`{"id":"cb-1","from":{"id":1,"is_bot":false,"first_name":"Ada"},"message":{"message_id":5,"chat":{"id":1,"type":"private"},"date":100,"text":"ok"},"chat_instance":"ci","data":"payload"}`), &accessible); err != nil {
		t.Fatalf("decode accessible callback: %v", err)
	}
	if accessible.Message == nil || accessible.MaybeMessage == nil || accessible.MaybeMessage.Message == nil || accessible.Message.Text != "ok" {
		t.Fatalf("unexpected accessible callback message: %#v", accessible)
	}

	var inaccessible CallbackQuery
	if err := json.Unmarshal([]byte(`{"id":"cb-2","from":{"id":1,"is_bot":false,"first_name":"Ada"},"message":{"message_id":6,"chat":{"id":1,"type":"private"},"date":0}}`), &inaccessible); err != nil {
		t.Fatalf("decode inaccessible callback: %v", err)
	}
	if inaccessible.Message != nil || inaccessible.MaybeMessage == nil || inaccessible.MaybeMessage.InaccessibleMessage == nil || inaccessible.MaybeMessage.MessageID != 6 {
		t.Fatalf("unexpected inaccessible callback message: %#v", inaccessible)
	}
}

func TestExternalReplyInfoDecodesGiveawayPayloads(t *testing.T) {
	var info ExternalReplyInfo
	body := []byte(`{"origin":{"type":"chat","date":1,"sender_chat":{"id":-100,"type":"supergroup","title":"Group"}},"giveaway":{"chats":[{"id":-100,"type":"supergroup","title":"Group"}],"winners_selection_date":200,"winner_count":2,"prize_star_count":50},"giveaway_winners":{"chat":{"id":-100,"type":"supergroup","title":"Group"},"giveaway_message_id":10,"winners_selection_date":300,"winner_count":1,"winners":[{"id":7,"is_bot":false,"first_name":"Winner"}],"was_refunded":true}}`)
	if err := json.Unmarshal(body, &info); err != nil {
		t.Fatalf("decode external reply: %v", err)
	}
	if _, ok := info.Origin.(MessageOriginChat); !ok {
		t.Fatalf("unexpected origin: %#v", info.Origin)
	}
	if info.Giveaway == nil || info.Giveaway.PrizeStarCount != 50 || info.Giveaway.WinnerCount != 2 {
		t.Fatalf("unexpected giveaway: %#v", info.Giveaway)
	}
	if info.GiveawayWinners == nil || len(info.GiveawayWinners.Winners) != 1 || !info.GiveawayWinners.WasRefunded {
		t.Fatalf("unexpected giveaway winners: %#v", info.GiveawayWinners)
	}
}

func TestMessageDecodesServiceDirectStoryAndMediaMetadata(t *testing.T) {
	var message Message
	body := []byte(`{
		"message_id":201,
		"chat":{"id":123,"type":"private"},
		"date":2000,
		"new_chat_members":[{"id":2,"is_bot":false,"first_name":"New"}],
		"left_chat_member":{"id":3,"is_bot":false,"first_name":"Left"},
		"chat_owner_left":{"new_owner":{"id":4,"is_bot":false,"first_name":"Return"}},
		"chat_owner_changed":{"new_owner":{"id":5,"is_bot":false,"first_name":"Owner"}},
		"new_chat_title":"New title",
		"new_chat_photo":[{"file_id":"photo","file_unique_id":"photo-u","width":10,"height":10}],
		"delete_chat_photo":true,
		"group_chat_created":true,
		"supergroup_chat_created":true,
		"channel_chat_created":true,
		"migrate_to_chat_id":-1001,
		"migrate_from_chat_id":-1000,
		"users_shared":{"request_id":7,"users":[{"user_id":1234567890123,"first_name":"Ada","username":"ada","photo":[{"file_id":"u-photo","file_unique_id":"u-photo-u","width":20,"height":20}]}]},
		"chat_shared":{"request_id":8,"chat_id":-1001234567890,"title":"Shared","username":"shared","photo":[{"file_id":"c-photo","file_unique_id":"c-photo-u","width":30,"height":30}]},
		"connected_website":"example.test",
		"proximity_alert_triggered":{"traveler":{"id":9,"is_bot":false,"first_name":"Traveler"},"watcher":{"id":10,"is_bot":false,"first_name":"Watcher"},"distance":42},
		"message_auto_delete_timer_changed":{"message_auto_delete_time":86400},
		"boost_added":{"boost_count":2},
		"chat_background_set":{"type":{"type":"fill","fill":{"type":"solid","color":1122867},"dark_theme_dimming":15}},
		"direct_message_price_changed":{"are_direct_messages_enabled":true,"direct_message_star_count":5},
		"giveaway_created":{"prize_star_count":100},
		"giveaway":{"chats":[{"id":-101,"type":"channel","title":"Giveaway channel"}],"winners_selection_date":3100,"winner_count":2,"prize_star_count":50},
		"giveaway_winners":{"chat":{"id":-100,"type":"supergroup","title":"Group"},"giveaway_message_id":55,"winners_selection_date":3000,"winner_count":1,"winners":[{"id":11,"is_bot":false,"first_name":"Winner"}],"prize_star_count":100},
		"giveaway_completed":{"winner_count":3,"unclaimed_prize_count":1,"is_star_giveaway":true},
		"paid_message_price_changed":{"paid_message_star_count":9},
		"video_chat_scheduled":{"start_date":4000},
		"video_chat_started":{},
		"video_chat_ended":{"duration":3600},
		"video_chat_participants_invited":{"users":[{"id":12,"is_bot":false,"first_name":"Guest"}]},
		"video":{"file_id":"video","file_unique_id":"video-u","width":640,"height":360,"duration":30,"cover":[{"file_id":"cover","file_unique_id":"cover-u","width":320,"height":180}],"start_timestamp":5,"qualities":[{"file_id":"q","file_unique_id":"q-u","width":1280,"height":720,"codec":"h264","file_size":1234}],"file_name":"clip.mp4","mime_type":"video/mp4"}
	}`)
	if err := json.Unmarshal(body, &message); err != nil {
		t.Fatalf("decode message: %v", err)
	}
	if len(message.NewChatMembers) != 1 || message.LeftChatMember == nil || message.ChatOwnerLeft == nil || message.ChatOwnerChanged == nil {
		t.Fatalf("unexpected chat membership service fields: %#v", message)
	}
	if message.NewChatTitle != "New title" || len(message.NewChatPhoto) != 1 || !message.DeleteChatPhoto || !message.GroupChatCreated || !message.SupergroupChatCreated || !message.ChannelChatCreated {
		t.Fatalf("unexpected chat service fields: %#v", message)
	}
	if message.MigrateToChatID != -1001 || message.MigrateFromChatID != -1000 {
		t.Fatalf("unexpected migration ids: %#v", message)
	}
	if message.UsersShared == nil || len(message.UsersShared.Users) != 1 || message.UsersShared.Users[0].UserID != 1234567890123 {
		t.Fatalf("unexpected users_shared: %#v", message.UsersShared)
	}
	if message.ChatShared == nil || message.ChatShared.ChatID != -1001234567890 || message.ConnectedWebsite != "example.test" {
		t.Fatalf("unexpected shared chat/website: %#v", message)
	}
	if message.ProximityAlertTriggered == nil || message.ProximityAlertTriggered.Distance != 42 || message.MessageAutoDeleteTimerChanged.MessageAutoDeleteTime != 86400 || message.BoostAdded.BoostCount != 2 {
		t.Fatalf("unexpected proximity/auto-delete/boost fields: %#v", message)
	}
	background, ok := message.ChatBackgroundSet.Type.(BackgroundTypeFill)
	if message.ChatBackgroundSet == nil || !ok {
		t.Fatalf("unexpected chat background: %#v", message.ChatBackgroundSet)
	}
	fill, ok := background.Fill.(BackgroundFillSolid)
	if !ok || fill.Color != 1122867 || background.DarkThemeDimming != 15 {
		t.Fatalf("unexpected background fill: %#v", background.Fill)
	}
	if message.DirectMessagePriceChanged == nil || !message.DirectMessagePriceChanged.AreDirectMessagesEnabled || message.DirectMessagePriceChanged.DirectMessageStarCount != 5 {
		t.Fatalf("unexpected direct message price: %#v", message.DirectMessagePriceChanged)
	}
	if message.GiveawayCreated == nil || message.GiveawayCreated.PrizeStarCount != 100 || message.Giveaway == nil || message.Giveaway.PrizeStarCount != 50 || message.GiveawayWinners == nil || len(message.GiveawayWinners.Winners) != 1 || message.GiveawayCompleted == nil || !message.GiveawayCompleted.IsStarGiveaway {
		t.Fatalf("unexpected giveaway service fields: %#v", message)
	}
	if message.PaidMessagePriceChanged == nil || message.PaidMessagePriceChanged.PaidMessageStarCount != 9 {
		t.Fatalf("unexpected paid message price: %#v", message.PaidMessagePriceChanged)
	}
	if message.VideoChatScheduled == nil || message.VideoChatScheduled.StartDate != 4000 || message.VideoChatStarted == nil || message.VideoChatEnded.Duration != 3600 || len(message.VideoChatParticipantsInvited.Users) != 1 {
		t.Fatalf("unexpected video chat fields: %#v", message)
	}
	if message.Video == nil || message.Video.StartTimestamp != 5 || len(message.Video.Cover) != 1 || len(message.Video.Qualities) != 1 || message.Video.Qualities[0].Codec != "h264" {
		t.Fatalf("unexpected video metadata: %#v", message.Video)
	}
}

func TestChatBackgroundUnknownTypeFails(t *testing.T) {
	var background ChatBackground
	if err := json.Unmarshal([]byte(`{"type":{"type":"unknown"}}`), &background); err == nil {
		t.Fatal("expected unsupported background type error")
	}
}

func TestBackgroundTypesDecodeVariants(t *testing.T) {
	tests := []struct {
		name string
		body string
		want any
	}{
		{
			name: "wallpaper",
			body: `{"type":{"type":"wallpaper","document":{"file_id":"doc","file_unique_id":"doc-u"},"dark_theme_dimming":20,"is_blurred":true,"is_moving":true}}`,
			want: BackgroundTypeWallpaper{},
		},
		{
			name: "pattern gradient",
			body: `{"type":{"type":"pattern","document":{"file_id":"doc","file_unique_id":"doc-u"},"fill":{"type":"gradient","top_color":1,"bottom_color":2,"rotation_angle":45},"intensity":50,"is_inverted":true,"is_moving":true}}`,
			want: BackgroundTypePattern{},
		},
		{
			name: "chat theme",
			body: `{"type":{"type":"chat_theme","theme_name":"day"}}`,
			want: BackgroundTypeChatTheme{},
		},
		{
			name: "freeform gradient",
			body: `{"type":{"type":"fill","fill":{"type":"freeform_gradient","colors":[1,2,3]},"dark_theme_dimming":10}}`,
			want: BackgroundTypeFill{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var background ChatBackground
			if err := json.Unmarshal([]byte(tt.body), &background); err != nil {
				t.Fatalf("decode background: %v", err)
			}
			switch tt.want.(type) {
			case BackgroundTypeWallpaper:
				value, ok := background.Type.(BackgroundTypeWallpaper)
				if !ok || !value.IsBlurred || !value.IsMoving || value.DarkThemeDimming != 20 {
					t.Fatalf("unexpected wallpaper: %#v", background.Type)
				}
			case BackgroundTypePattern:
				value, ok := background.Type.(BackgroundTypePattern)
				if !ok || value.Intensity != 50 || !value.IsInverted || !value.IsMoving {
					t.Fatalf("unexpected pattern: %#v", background.Type)
				}
				fill, ok := value.Fill.(BackgroundFillGradient)
				if !ok || fill.TopColor != 1 || fill.BottomColor != 2 || fill.RotationAngle != 45 {
					t.Fatalf("unexpected pattern fill: %#v", value.Fill)
				}
			case BackgroundTypeChatTheme:
				value, ok := background.Type.(BackgroundTypeChatTheme)
				if !ok || value.ThemeName != "day" {
					t.Fatalf("unexpected chat theme: %#v", background.Type)
				}
			case BackgroundTypeFill:
				value, ok := background.Type.(BackgroundTypeFill)
				if !ok || value.DarkThemeDimming != 10 {
					t.Fatalf("unexpected fill background: %#v", background.Type)
				}
				fill, ok := value.Fill.(BackgroundFillFreeformGradient)
				if !ok || len(fill.Colors) != 3 {
					t.Fatalf("unexpected freeform fill: %#v", value.Fill)
				}
			}
		})
	}
}
