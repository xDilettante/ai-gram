package telegram

import "testing"

func TestInlineKeyboardValidation(t *testing.T) {
	callback64 := "1234567890123456789012345678901234567890123456789012345678901234"

	valid := NewInlineKeyboard(
		[]InlineKeyboardButton{InlineButtonURL("Open", "https://example.com")},
		[]InlineKeyboardButton{InlineButtonCallback("Confirm", "confirm")},
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
	if err := ValidateReplyMarkup(NewInlineKeyboard([]InlineKeyboardButton{InlineButtonCallback("OK", callback64)})); err != nil {
		t.Fatalf("64 byte callback_data should be valid: %v", err)
	}
	if err := ValidateReplyMarkup(NewInlineKeyboard([]InlineKeyboardButton{InlineButtonURL("Open", "http://example.com")})); err != nil {
		t.Fatalf("http URL should be valid: %v", err)
	}

	tests := []struct {
		name   string
		markup ReplyMarkup
	}{
		{name: "empty text", markup: NewInlineKeyboard([]InlineKeyboardButton{InlineButtonCallback("", "x")})},
		{name: "empty callback", markup: NewInlineKeyboard([]InlineKeyboardButton{{Text: "x", CallbackData: ""}})},
		{name: "callback too long", markup: NewInlineKeyboard([]InlineKeyboardButton{InlineButtonCallback("x", callback64+"x")})},
		{name: "file url", markup: NewInlineKeyboard([]InlineKeyboardButton{InlineButtonURL("x", "file:///tmp/a")})},
		{name: "two actions", markup: NewInlineKeyboard([]InlineKeyboardButton{{Text: "x", URL: "https://example.com", CallbackData: "x"}})},
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
	valid := NewReplyKeyboard([]KeyboardButton{KeyboardButtonText("Yes"), KeyboardButtonContact("Phone"), KeyboardButtonLocation("Place")})
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

	tests := []struct {
		name   string
		markup ReplyMarkup
	}{
		{name: "contact and location", markup: NewReplyKeyboard([]KeyboardButton{{Text: "x", RequestContact: true, RequestLocation: true}})},
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
