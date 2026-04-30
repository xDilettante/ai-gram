package telegram

import (
	"encoding/json"
	"testing"
)

func TestWebAppServiceMessagesDecode(t *testing.T) {
	var update Update
	data := []byte(`{
		"update_id":100,
		"message":{
			"message_id":1,
			"date":10,
			"chat":{"id":123,"type":"private"},
			"web_app_data":{"data":"opaque-web-app-data","button_text":"Open app"},
			"write_access_allowed":{"from_request":true,"web_app_name":"Test App","from_attachment_menu":true}
		}
	}`)
	if err := json.Unmarshal(data, &update); err != nil {
		t.Fatalf("decode update: %v", err)
	}
	message := update.Message
	if message == nil {
		t.Fatal("expected message")
	}
	if message.WebAppData == nil || message.WebAppData.Data != "opaque-web-app-data" || message.WebAppData.ButtonText != "Open app" {
		t.Fatalf("unexpected web_app_data: %+v", message.WebAppData)
	}
	if message.WriteAccessAllowed == nil ||
		!message.WriteAccessAllowed.FromRequest ||
		message.WriteAccessAllowed.WebAppName != "Test App" ||
		!message.WriteAccessAllowed.FromAttachmentMenu {
		t.Fatalf("unexpected write_access_allowed: %+v", message.WriteAccessAllowed)
	}
}

func TestSentWebAppMessageDecode(t *testing.T) {
	var message SentWebAppMessage
	if err := json.Unmarshal([]byte(`{"inline_message_id":"inline-id"}`), &message); err != nil {
		t.Fatalf("decode sent web app message: %v", err)
	}
	if message.InlineMessageID != "inline-id" {
		t.Fatalf("unexpected sent web app message: %+v", message)
	}
}
