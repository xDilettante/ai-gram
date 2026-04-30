package telegram

import (
	"encoding/json"
	"testing"
)

func TestGameMessageDecoding(t *testing.T) {
	message := mustDecodeMessage(t, `{
		"message_id": 10,
		"chat": {"id": 123, "type": "private"},
		"date": 100,
		"game": {
			"title": "Chess",
			"description": "A game",
			"photo": [{"file_id": "photo", "file_unique_id": "unique", "width": 64, "height": 64}],
			"text": "High scores",
			"text_entities": [{"type": "bold", "offset": 0, "length": 4}],
			"animation": {"file_id": "anim", "file_unique_id": "anim-unique", "width": 320, "height": 240, "duration": 5}
		}
	}`)

	if message.Game == nil {
		t.Fatal("expected game")
	}
	if message.Game.Title != "Chess" || message.Game.Description != "A game" || len(message.Game.Photo) != 1 {
		t.Fatalf("unexpected game: %+v", message.Game)
	}
	if message.Game.Text != "High scores" || len(message.Game.TextEntities) != 1 || message.Game.Animation == nil {
		t.Fatalf("unexpected optional game fields: %+v", message.Game)
	}
}

func TestCallbackGameInlineKeyboardButtonValidation(t *testing.T) {
	markup := NewInlineKeyboard(
		[]InlineKeyboardButton{InlineButtonGame("Play")},
		[]InlineKeyboardButton{InlineButtonCallback("Scores", "scores")},
	)
	if err := ValidateReplyMarkup(markup); err != nil {
		t.Fatalf("unexpected callback_game markup error: %v", err)
	}
	body, err := json.Marshal(markup)
	if err != nil {
		t.Fatalf("marshal markup: %v", err)
	}
	var decoded map[string][][]map[string]any
	if err := json.Unmarshal(body, &decoded); err != nil {
		t.Fatalf("unmarshal markup: %v", err)
	}
	if _, ok := decoded["inline_keyboard"][0][0]["callback_game"]; !ok {
		t.Fatalf("callback_game missing from payload: %s", string(body))
	}

	tests := []struct {
		name   string
		markup ReplyMarkup
	}{
		{name: "callback game with another action", markup: NewInlineKeyboard([]InlineKeyboardButton{{Text: "Play", CallbackData: "x", CallbackGame: &CallbackGame{}}})},
		{name: "callback game not first", markup: NewInlineKeyboard([]InlineKeyboardButton{InlineButtonCallback("Other", "x"), InlineButtonGame("Play")})},
		{name: "callback game not first row", markup: NewInlineKeyboard([]InlineKeyboardButton{InlineButtonCallback("Other", "x")}, []InlineKeyboardButton{InlineButtonGame("Play")})},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ValidateReplyMarkup(tt.markup); err == nil {
				t.Fatal("expected error")
			}
		})
	}
}
