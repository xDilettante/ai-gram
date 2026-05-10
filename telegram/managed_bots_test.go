package telegram

import (
	"encoding/json"
	"testing"
)

func TestUserDecodesCanManageBots(t *testing.T) {
	var user User
	if err := json.Unmarshal([]byte(`{"id":1,"is_bot":true,"first_name":"Manager","can_manage_bots":true}`), &user); err != nil {
		t.Fatalf("decode user: %v", err)
	}
	if !user.CanManageBots {
		t.Fatalf("expected can_manage_bots: %+v", user)
	}
}

func TestMessageDecodesManagedBotCreated(t *testing.T) {
	payload := []byte(`{
		"message_id": 10,
		"chat": {"id": 1, "type": "private"},
		"date": 123,
		"managed_bot_created": {"bot": {"id": 77, "is_bot": true, "first_name": "ChildBot", "username": "child_bot"}}
	}`)

	var message Message
	if err := json.Unmarshal(payload, &message); err != nil {
		t.Fatalf("decode message: %v", err)
	}
	if message.ManagedBotCreated == nil || message.ManagedBotCreated.Bot.ID != 77 || !message.ManagedBotCreated.Bot.IsBot {
		t.Fatalf("unexpected managed_bot_created: %+v", message.ManagedBotCreated)
	}
}

func TestUpdateDecodesManagedBotAndEffectiveUser(t *testing.T) {
	payload := []byte(`{
		"update_id": 200,
		"managed_bot": {
			"user": {"id": 7, "is_bot": false, "first_name": "Owner"},
			"bot": {"id": 77, "is_bot": true, "first_name": "ChildBot"}
		}
	}`)

	var update Update
	if err := json.Unmarshal(payload, &update); err != nil {
		t.Fatalf("decode update: %v", err)
	}
	if update.ManagedBot == nil || update.ManagedBot.User.ID != 7 || update.ManagedBot.Bot.ID != 77 {
		t.Fatalf("unexpected managed_bot update: %+v", update.ManagedBot)
	}
	if user := update.EffectiveUser(); user == nil || user.ID != 7 {
		t.Fatalf("unexpected effective user: %+v", user)
	}
	if chat := update.EffectiveChat(); chat != nil {
		t.Fatalf("managed_bot update should not invent chat: %+v", chat)
	}
}

func TestBotAccessSettingsDecodesUsers(t *testing.T) {
	payload := []byte(`{
		"is_access_restricted": true,
		"added_users": [
			{"id": 7, "is_bot": false, "first_name": "Alice"},
			{"id": 8, "is_bot": false, "first_name": "Bob"}
		]
	}`)

	var settings BotAccessSettings
	if err := json.Unmarshal(payload, &settings); err != nil {
		t.Fatalf("decode access settings: %v", err)
	}
	if !settings.IsAccessRestricted || len(settings.AddedUsers) != 2 || settings.AddedUsers[0].ID != 7 || settings.AddedUsers[1].ID != 8 {
		t.Fatalf("unexpected access settings: %+v", settings)
	}
}

func TestValidateKeyboardButtonManagedBot(t *testing.T) {
	button := KeyboardButtonManagedBot("Create bot", KeyboardButtonRequestManagedBot{RequestID: 42, SuggestedName: "Test Bot", SuggestedUsername: "test_bot"})
	if err := ValidateKeyboardButton(button); err != nil {
		t.Fatalf("valid managed bot button rejected: %v", err)
	}
	body, err := json.Marshal(button)
	if err != nil {
		t.Fatalf("marshal button: %v", err)
	}
	var payload map[string]any
	if err := json.Unmarshal(body, &payload); err != nil {
		t.Fatalf("decode button: %v", err)
	}
	request := payload["request_managed_bot"].(map[string]any)
	if request["request_id"] != float64(42) || request["suggested_name"] != "Test Bot" || request["suggested_username"] != "test_bot" {
		t.Fatalf("unexpected request_managed_bot payload: %#v", request)
	}
}

func TestValidateKeyboardButtonRequests(t *testing.T) {
	premium := true
	if err := ValidateKeyboardButton(KeyboardButtonUsers("Pick", KeyboardButtonRequestUsers{RequestID: 1, UserIsPremium: &premium, MaxQuantity: 2})); err != nil {
		t.Fatalf("valid request_users rejected: %v", err)
	}
	forum := false
	if err := ValidateKeyboardButton(KeyboardButtonChat("Chat", KeyboardButtonRequestChat{RequestID: 2, ChatIsChannel: false, ChatIsForum: &forum})); err != nil {
		t.Fatalf("valid request_chat rejected: %v", err)
	}

	tests := []struct {
		name   string
		button KeyboardButton
	}{
		{name: "empty request id", button: KeyboardButtonManagedBot("Create", KeyboardButtonRequestManagedBot{})},
		{name: "too many actions", button: KeyboardButton{Text: "Bad", RequestContact: true, RequestManagedBot: &KeyboardButtonRequestManagedBot{RequestID: 1}}},
		{name: "negative max quantity", button: KeyboardButtonUsers("Pick", KeyboardButtonRequestUsers{RequestID: 1, MaxQuantity: -1})},
		{name: "too many users", button: KeyboardButtonUsers("Pick", KeyboardButtonRequestUsers{RequestID: 1, MaxQuantity: 11})},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ValidateKeyboardButton(tt.button); err == nil {
				t.Fatal("expected validation error")
			}
		})
	}
}

func TestValidatePreparedKeyboardButton(t *testing.T) {
	if err := ValidatePreparedKeyboardButton(KeyboardButtonManagedBot("Create", KeyboardButtonRequestManagedBot{RequestID: 1})); err != nil {
		t.Fatalf("valid prepared managed bot button rejected: %v", err)
	}
	if err := ValidatePreparedKeyboardButton(KeyboardButtonText("Plain")); err == nil {
		t.Fatal("expected plain text button to be rejected for prepared keyboard")
	}
	if err := ValidatePreparedKeyboardButton(KeyboardButtonContact("Phone")); err == nil {
		t.Fatal("expected contact button to be rejected for prepared keyboard")
	}
}
