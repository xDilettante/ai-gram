package telegram

import (
	"encoding/json"
	"testing"
)

func TestUpdateDecodesBusinessConnectionAndEffectiveUser(t *testing.T) {
	payload := []byte(`{
		"update_id": 800,
		"business_connection": {
			"id": "bc-1",
			"user": {"id": 7, "is_bot": false, "first_name": "Business"},
			"user_chat_id": 7000,
			"date": 123456,
			"rights": {
				"can_reply": true,
				"can_read_messages": true,
				"can_delete_sent_messages": true,
				"can_delete_all_messages": true,
				"can_edit_name": true,
				"can_edit_bio": true,
				"can_edit_profile_photo": true,
				"can_edit_username": true,
				"can_change_gift_settings": true,
				"can_view_gifts_and_stars": true,
				"can_convert_gifts_to_stars": true,
				"can_transfer_and_upgrade_gifts": true,
				"can_transfer_stars": true,
				"can_manage_stories": true
			},
			"is_enabled": true
		}
	}`)

	var update Update
	if err := json.Unmarshal(payload, &update); err != nil {
		t.Fatalf("decode update: %v", err)
	}
	connection := update.BusinessConnection
	if connection == nil || connection.ID != "bc-1" || connection.User.ID != 7 || connection.UserChatID != 7000 || connection.Date != 123456 || !connection.IsEnabled {
		t.Fatalf("unexpected business connection: %+v", connection)
	}
	if connection.Rights == nil || !connection.Rights.CanReply || !connection.Rights.CanReadMessages || !connection.Rights.CanDeleteSentMessages || !connection.Rights.CanDeleteAllMessages || !connection.Rights.CanEditName || !connection.Rights.CanEditBio || !connection.Rights.CanEditProfilePhoto || !connection.Rights.CanEditUsername || !connection.Rights.CanChangeGiftSettings || !connection.Rights.CanViewGiftsAndStars || !connection.Rights.CanConvertGiftsToStars || !connection.Rights.CanTransferAndUpgradeGifts || !connection.Rights.CanTransferStars || !connection.Rights.CanManageStories {
		t.Fatalf("unexpected business rights: %+v", connection.Rights)
	}
	if user := update.EffectiveUser(); user == nil || user.ID != 7 {
		t.Fatalf("unexpected effective user: %+v", user)
	}
	if chat := update.EffectiveChat(); chat != nil {
		t.Fatalf("business connection update should not invent chat: %+v", chat)
	}
}

func TestUpdateDecodesBusinessMessagesAndEffectiveChat(t *testing.T) {
	businessMessagePayload := []byte(`{
		"update_id": 801,
		"business_message": {
			"message_id": 11,
			"from": {"id": 7, "is_bot": false, "first_name": "Business"},
			"sender_business_bot": {"id": 77, "is_bot": true, "first_name": "HelperBot"},
			"chat": {"id": 1001, "type": "private"},
			"date": 123,
			"business_connection_id": "bc-1",
			"is_from_offline": true,
			"text": "business message"
		}
	}`)

	var businessUpdate Update
	if err := json.Unmarshal(businessMessagePayload, &businessUpdate); err != nil {
		t.Fatalf("decode business message update: %v", err)
	}
	message := businessUpdate.BusinessMessage
	if message == nil || message.BusinessConnectionID != "bc-1" || message.SenderBusinessBot == nil || message.SenderBusinessBot.ID != 77 || !message.IsFromOffline {
		t.Fatalf("unexpected business message: %+v", message)
	}
	if chat := businessUpdate.EffectiveChat(); chat == nil || chat.ID != 1001 {
		t.Fatalf("unexpected effective chat for business message: %+v", chat)
	}
	if user := businessUpdate.EffectiveUser(); user == nil || user.ID != 7 {
		t.Fatalf("unexpected effective user for business message: %+v", user)
	}

	editedPayload := []byte(`{
		"update_id": 802,
		"edited_business_message": {
			"message_id": 12,
			"chat": {"id": 1002, "type": "private"},
			"date": 124,
			"business_connection_id": "bc-2",
			"text": "edited business message"
		}
	}`)
	var editedUpdate Update
	if err := json.Unmarshal(editedPayload, &editedUpdate); err != nil {
		t.Fatalf("decode edited business message update: %v", err)
	}
	if editedUpdate.EditedBusinessMessage == nil || editedUpdate.EditedBusinessMessage.BusinessConnectionID != "bc-2" {
		t.Fatalf("unexpected edited business message: %+v", editedUpdate.EditedBusinessMessage)
	}
	if chat := editedUpdate.EffectiveChat(); chat == nil || chat.ID != 1002 {
		t.Fatalf("unexpected effective chat for edited business message: %+v", chat)
	}
}

func TestUpdateDecodesDeletedBusinessMessagesAndEffectiveChat(t *testing.T) {
	payload := []byte(`{
		"update_id": 803,
		"deleted_business_messages": {
			"business_connection_id": "bc-1",
			"chat": {"id": 1003, "type": "private"},
			"message_ids": [11, 12]
		}
	}`)

	var update Update
	if err := json.Unmarshal(payload, &update); err != nil {
		t.Fatalf("decode deleted business messages: %v", err)
	}
	deleted := update.DeletedBusinessMessages
	if deleted == nil || deleted.BusinessConnectionID != "bc-1" || deleted.Chat.ID != 1003 || len(deleted.MessageIDs) != 2 || deleted.MessageIDs[0] != 11 || deleted.MessageIDs[1] != 12 {
		t.Fatalf("unexpected deleted business messages: %+v", deleted)
	}
	if chat := update.EffectiveChat(); chat == nil || chat.ID != 1003 {
		t.Fatalf("unexpected effective chat for deleted business messages: %+v", chat)
	}
}
