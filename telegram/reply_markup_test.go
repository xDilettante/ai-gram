package telegram

import (
	"encoding/json"
	"testing"
)

func TestInlineKeyboardValidation(t *testing.T) {
	callback64 := "1234567890123456789012345678901234567890123456789012345678901234"

	valid := NewInlineKeyboard(
		[]InlineKeyboardButton{InlineButtonURL("Open", "https://example.com")},
		[]InlineKeyboardButton{InlineButtonCallback("Confirm", "confirm")},
		[]InlineKeyboardButton{InlineButtonWebApp("App", "https://example.com/app")},
	)
	if err := ValidateReplyMarkup(valid); err != nil {
		t.Fatalf("unexpected valid inline keyboard error: %v", err)
	}
	if got := InlineButtonURL("Open", "https://example.com"); got.URL != "https://example.com" || got.Text != "Open" {
		t.Fatalf("unexpected URL button: %+v", got)
	}
	if got := InlineButtonCallback("Confirm", "confirm"); got.CallbackData != "confirm" || got.Text != "Confirm" {
		t.Fatalf("unexpected callback button: %+v", got)
	}
	if got := InlineButtonWebApp("App", "https://example.com/app"); got.WebApp == nil || got.WebApp.URL != "https://example.com/app" || got.Text != "App" {
		t.Fatalf("unexpected web_app button: %+v", got)
	}
	if err := ValidateReplyMarkup(NewInlineKeyboard([]InlineKeyboardButton{InlineButtonCallback("OK", callback64)})); err != nil {
		t.Fatalf("64 byte callback_data should be valid: %v", err)
	}
	if err := ValidateReplyMarkup(NewInlineKeyboard([]InlineKeyboardButton{InlineButtonURL("Open", "http://example.com")})); err != nil {
		t.Fatalf("http URL should be valid: %v", err)
	}
	if err := ValidateReplyMarkup(NewInlineKeyboard([]InlineKeyboardButton{InlineButtonURL("Mention", "tg://user?id=123")})); err != nil {
		t.Fatalf("tg URL should be valid: %v", err)
	}
	if err := ValidateReplyMarkup(NewInlineKeyboard([]InlineKeyboardButton{InlineButtonLoginURL("Login", "https://example.com/login")})); err != nil {
		t.Fatalf("login_url should be valid: %v", err)
	}
	if err := ValidateReplyMarkup(NewInlineKeyboard([]InlineKeyboardButton{InlineButtonSwitchInlineQuery("Inline", "")})); err != nil {
		t.Fatalf("empty switch_inline_query should be valid: %v", err)
	}
	if err := ValidateReplyMarkup(NewInlineKeyboard([]InlineKeyboardButton{InlineButtonSwitchInlineQueryCurrentChat("Here", "")})); err != nil {
		t.Fatalf("empty switch_inline_query_current_chat should be valid: %v", err)
	}
	if err := ValidateReplyMarkup(NewInlineKeyboard([]InlineKeyboardButton{InlineButtonSwitchInlineQueryChosenChat("Choose", SwitchInlineQueryChosenChat{Query: "q", AllowUserChats: true})})); err != nil {
		t.Fatalf("switch_inline_query_chosen_chat should be valid: %v", err)
	}
	if err := ValidateReplyMarkup(NewInlineKeyboard([]InlineKeyboardButton{InlineButtonCopyText("Copy", "copy me")})); err != nil {
		t.Fatalf("copy_text should be valid: %v", err)
	}
	if err := ValidateReplyMarkup(NewInlineKeyboard([]InlineKeyboardButton{InlineButtonPay("Pay")})); err != nil {
		t.Fatalf("pay should be valid: %v", err)
	}
	styled := InlineButtonCallback("Styled", "style")
	styled.IconCustomEmojiID = "emoji-id"
	styled.Style = "primary"
	if err := ValidateReplyMarkup(NewInlineKeyboard([]InlineKeyboardButton{styled})); err != nil {
		t.Fatalf("button icon/style should be valid: %v", err)
	}

	tests := []struct {
		name   string
		markup ReplyMarkup
	}{
		{name: "empty text", markup: NewInlineKeyboard([]InlineKeyboardButton{InlineButtonCallback("", "x")})},
		{name: "empty callback", markup: NewInlineKeyboard([]InlineKeyboardButton{{Text: "x", CallbackData: ""}})},
		{name: "callback too long", markup: NewInlineKeyboard([]InlineKeyboardButton{InlineButtonCallback("x", callback64+"x")})},
		{name: "file url", markup: NewInlineKeyboard([]InlineKeyboardButton{InlineButtonURL("x", "file:///tmp/a")})},
		{name: "invalid web app url", markup: NewInlineKeyboard([]InlineKeyboardButton{InlineButtonWebApp("x", "ftp://example.com/app")})},
		{name: "login url must be https", markup: NewInlineKeyboard([]InlineKeyboardButton{InlineButtonLoginURL("x", "http://example.com/login")})},
		{name: "copy text empty", markup: NewInlineKeyboard([]InlineKeyboardButton{InlineButtonCopyText("x", "")})},
		{name: "invalid style", markup: NewInlineKeyboard([]InlineKeyboardButton{{Text: "x", CallbackData: "x", Style: "warning"}})},
		{name: "blank icon", markup: NewInlineKeyboard([]InlineKeyboardButton{{Text: "x", CallbackData: "x", IconCustomEmojiID: "   "}})},
		{name: "two actions", markup: NewInlineKeyboard([]InlineKeyboardButton{{Text: "x", URL: "https://example.com", CallbackData: "x"}})},
		{name: "login and switch actions", markup: NewInlineKeyboard([]InlineKeyboardButton{{Text: "x", LoginURL: &LoginUrl{URL: "https://example.com"}, SwitchInlineQueryChosenChat: &SwitchInlineQueryChosenChat{}}})},
		{name: "pay not first", markup: NewInlineKeyboard([]InlineKeyboardButton{InlineButtonCallback("Other", "x"), InlineButtonPay("Pay")})},
		{name: "no action", markup: NewInlineKeyboard([]InlineKeyboardButton{{Text: "x"}})},
		{name: "empty keyboard", markup: InlineKeyboardMarkup{}},
		{name: "empty row", markup: NewInlineKeyboard([]InlineKeyboardButton{})},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ValidateReplyMarkup(tt.markup); err == nil {
				t.Fatal("expected error")
			}
		})
	}
}

func TestReplyKeyboardValidation(t *testing.T) {
	valid := NewReplyKeyboard([]KeyboardButton{KeyboardButtonText("Yes"), KeyboardButtonContact("Phone"), KeyboardButtonLocation("Place"), KeyboardButtonPoll("Poll", "quiz"), KeyboardButtonWebApp("App", "https://example.com/app")})
	if err := ValidateReplyMarkup(valid); err != nil {
		t.Fatalf("unexpected valid reply keyboard error: %v", err)
	}
	if got := KeyboardButtonText("Yes"); got.Text != "Yes" || got.RequestContact || got.RequestLocation {
		t.Fatalf("unexpected text button: %+v", got)
	}
	if got := KeyboardButtonContact("Phone"); got.Text != "Phone" || !got.RequestContact || got.RequestLocation {
		t.Fatalf("unexpected contact button: %+v", got)
	}
	if got := KeyboardButtonLocation("Place"); got.Text != "Place" || got.RequestContact || !got.RequestLocation {
		t.Fatalf("unexpected location button: %+v", got)
	}
	if got := KeyboardButtonWebApp("App", "https://example.com/app"); got.Text != "App" || got.WebApp == nil || got.WebApp.URL != "https://example.com/app" {
		t.Fatalf("unexpected web_app button: %+v", got)
	}
	if got := KeyboardButtonPoll("Poll", "regular"); got.Text != "Poll" || got.RequestPoll == nil || got.RequestPoll.Type != "regular" {
		t.Fatalf("unexpected request_poll button: %+v", got)
	}
	styled := KeyboardButtonText("Styled")
	styled.IconCustomEmojiID = "emoji-id"
	styled.Style = "success"
	if err := ValidateReplyMarkup(NewReplyKeyboard([]KeyboardButton{styled})); err != nil {
		t.Fatalf("button icon/style should be valid: %v", err)
	}

	tests := []struct {
		name   string
		markup ReplyMarkup
	}{
		{name: "contact and location", markup: NewReplyKeyboard([]KeyboardButton{{Text: "x", RequestContact: true, RequestLocation: true}})},
		{name: "contact and web app", markup: NewReplyKeyboard([]KeyboardButton{{Text: "x", RequestContact: true, WebApp: &WebAppInfo{URL: "https://example.com/app"}}})},
		{name: "poll and web app", markup: NewReplyKeyboard([]KeyboardButton{{Text: "x", RequestPoll: &KeyboardButtonPollType{}, WebApp: &WebAppInfo{URL: "https://example.com/app"}}})},
		{name: "invalid poll type", markup: NewReplyKeyboard([]KeyboardButton{KeyboardButtonPoll("Poll", "custom")})},
		{name: "invalid style", markup: NewReplyKeyboard([]KeyboardButton{{Text: "x", Style: "warning"}})},
		{name: "blank icon", markup: NewReplyKeyboard([]KeyboardButton{{Text: "x", IconCustomEmojiID: "   "}})},
		{name: "invalid web app", markup: NewReplyKeyboard([]KeyboardButton{KeyboardButtonWebApp("x", "ftp://example.com/app")})},
		{name: "empty text", markup: NewReplyKeyboard([]KeyboardButton{KeyboardButtonText("")})},
		{name: "empty keyboard", markup: ReplyKeyboardMarkup{}},
		{name: "empty row", markup: NewReplyKeyboard([]KeyboardButton{})},
		{name: "blank placeholder", markup: ReplyKeyboardMarkup{Keyboard: [][]KeyboardButton{{KeyboardButtonText("x")}}, InputFieldPlaceholder: "   "}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ValidateReplyMarkup(tt.markup); err == nil {
				t.Fatal("expected error")
			}
		})
	}
}

func TestReplyMarkupCompletionMarshal(t *testing.T) {
	query := ""
	markup := NewInlineKeyboard([]InlineKeyboardButton{
		{
			Text:                         "Switch",
			IconCustomEmojiID:            "emoji-id",
			Style:                        "primary",
			SwitchInlineQueryCurrentChat: &query,
		},
		InlineButtonLoginURL("Login", "https://example.com/login"),
		InlineButtonSwitchInlineQueryChosenChat("Choose", SwitchInlineQueryChosenChat{Query: "q", AllowGroupChats: true}),
		InlineButtonCopyText("Copy", "copy me"),
		InlineButtonPay("Pay"),
	})
	body, err := json.Marshal(markup)
	if err != nil {
		t.Fatalf("marshal inline keyboard: %v", err)
	}
	var decoded map[string][][]map[string]any
	if err := json.Unmarshal(body, &decoded); err != nil {
		t.Fatalf("decode inline keyboard: %v", err)
	}
	buttons := decoded["inline_keyboard"][0]
	if buttons[0]["switch_inline_query_current_chat"] != "" || buttons[0]["icon_custom_emoji_id"] != "emoji-id" || buttons[0]["style"] != "primary" {
		t.Fatalf("unexpected switch button: %#v", buttons[0])
	}
	if loginURL := buttons[1]["login_url"].(map[string]any); loginURL["url"] != "https://example.com/login" {
		t.Fatalf("unexpected login_url button: %#v", buttons[1])
	}
	if chosenChat := buttons[2]["switch_inline_query_chosen_chat"].(map[string]any); chosenChat["allow_group_chats"] != true {
		t.Fatalf("unexpected chosen chat button: %#v", buttons[2])
	}
	if copyText := buttons[3]["copy_text"].(map[string]any); copyText["text"] != "copy me" {
		t.Fatalf("unexpected copy_text button: %#v", buttons[3])
	}
	if buttons[4]["pay"] != true {
		t.Fatalf("unexpected pay button: %#v", buttons[4])
	}

	reply := NewReplyKeyboard([]KeyboardButton{{
		Text:              "Poll",
		IconCustomEmojiID: "emoji-id",
		Style:             "success",
		RequestPoll:       &KeyboardButtonPollType{Type: "quiz"},
	}})
	body, err = json.Marshal(reply)
	if err != nil {
		t.Fatalf("marshal reply keyboard: %v", err)
	}
	var decodedReply map[string][][]map[string]any
	if err := json.Unmarshal(body, &decodedReply); err != nil {
		t.Fatalf("decode reply keyboard: %v", err)
	}
	replyButton := decodedReply["keyboard"][0][0]
	requestPoll := replyButton["request_poll"].(map[string]any)
	if requestPoll["type"] != "quiz" || replyButton["icon_custom_emoji_id"] != "emoji-id" || replyButton["style"] != "success" {
		t.Fatalf("unexpected request_poll button: %#v", replyButton)
	}
}

func TestReplyKeyboardRemoveValidation(t *testing.T) {
	if err := ValidateReplyMarkup(RemoveKeyboard(false)); err != nil {
		t.Fatalf("unexpected remove keyboard error: %v", err)
	}
	if err := ValidateReplyMarkup(ReplyKeyboardRemove{RemoveKeyboard: false}); err == nil {
		t.Fatal("expected manual remove keyboard error")
	}
}

func TestForceReplyValidation(t *testing.T) {
	if err := ValidateReplyMarkup(NewForceReply()); err != nil {
		t.Fatalf("unexpected force reply error: %v", err)
	}
	if err := ValidateReplyMarkup(ForceReply{ForceReply: false}); err == nil {
		t.Fatal("expected manual force reply error")
	}
	if err := ValidateReplyMarkup(ForceReply{ForceReply: true, InputFieldPlaceholder: "   "}); err == nil {
		t.Fatal("expected blank placeholder error")
	}
}

func TestValidateReplyMarkupNil(t *testing.T) {
	if err := ValidateReplyMarkup(nil); err != nil {
		t.Fatalf("nil markup should be valid: %v", err)
	}
}
