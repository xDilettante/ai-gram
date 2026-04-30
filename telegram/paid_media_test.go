package telegram

import (
	"encoding/json"
	"testing"
)

func TestPaidMediaInfoDecoding(t *testing.T) {
	message := mustDecodeMessage(t, `{
		"message_id": 100,
		"date": 1,
		"chat": {"id": 123, "type": "private"},
		"paid_media": {
			"star_count": 5,
			"paid_media": [
				{"type":"preview","width":320,"height":240,"duration":10},
				{"type":"photo","photo":[{"file_id":"photo","file_unique_id":"photo-u","width":10,"height":10}]},
				{"type":"video","video":{"file_id":"video","file_unique_id":"video-u","width":640,"height":480,"duration":12}}
			]
		}
	}`)
	if message.PaidMedia == nil || message.PaidMedia.StarCount != 5 || len(message.PaidMedia.PaidMedia) != 3 {
		t.Fatalf("unexpected paid media: %+v", message.PaidMedia)
	}
	if preview, ok := message.PaidMedia.PaidMedia[0].(PaidMediaPreview); !ok || preview.Width != 320 || preview.Duration != 10 {
		t.Fatalf("unexpected preview: %#v", message.PaidMedia.PaidMedia[0])
	}
	if photo, ok := message.PaidMedia.PaidMedia[1].(PaidMediaPhoto); !ok || len(photo.Photo) != 1 || photo.Photo[0].FileID != "photo" {
		t.Fatalf("unexpected photo: %#v", message.PaidMedia.PaidMedia[1])
	}
	if video, ok := message.PaidMedia.PaidMedia[2].(PaidMediaVideo); !ok || video.Video.FileID != "video" {
		t.Fatalf("unexpected video: %#v", message.PaidMedia.PaidMedia[2])
	}
}

func TestPaidMediaUnknownTypeReturnsError(t *testing.T) {
	var info PaidMediaInfo
	if err := json.Unmarshal([]byte(`{"star_count":5,"paid_media":[{"type":"unknown"}]}`), &info); err == nil {
		t.Fatal("expected unsupported paid media type error")
	}
}

func TestPaidMediaPurchasedUpdateDecodingAndEffectiveUser(t *testing.T) {
	update := mustDecodeUpdate(t, `{
		"update_id": 3,
		"purchased_paid_media": {
			"from": {"id": 9, "is_bot": false, "first_name": "Alice"},
			"paid_media_payload": "payload"
		}
	}`)
	if update.PurchasedPaidMedia == nil || update.PurchasedPaidMedia.PaidMediaPayload != "payload" {
		t.Fatalf("unexpected paid media purchase: %+v", update.PurchasedPaidMedia)
	}
	if chat := update.EffectiveChat(); chat != nil {
		t.Fatalf("paid media purchase should not invent an effective chat: %+v", chat)
	}
	if user := update.EffectiveUser(); user == nil || user.ID != 9 {
		t.Fatalf("unexpected paid media purchase effective user: %+v", user)
	}
}
