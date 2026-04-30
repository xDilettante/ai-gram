package telegram

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestGiftTypesDecode(t *testing.T) {
	data := []byte(`{
		"gifts":[{
			"id":"gift-1",
			"sticker":{"file_id":"sticker","file_unique_id":"unique","type":"regular","width":128,"height":128,"is_animated":false,"is_video":false},
			"star_count":25,
			"upgrade_star_count":50,
			"has_colors":true,
			"background":{"center_color":1,"edge_color":2,"text_color":3}
		}]
	}`)
	var gifts Gifts
	if err := json.Unmarshal(data, &gifts); err != nil {
		t.Fatalf("decode gifts: %v", err)
	}
	if len(gifts.Gifts) != 1 || gifts.Gifts[0].ID != "gift-1" || gifts.Gifts[0].Background.CenterColor != 1 {
		t.Fatalf("unexpected gifts: %+v", gifts)
	}
}

func TestOwnedGiftsDecodePolymorphicItems(t *testing.T) {
	data := []byte(`{
		"total_count":2,
		"next_offset":"next",
		"gifts":[
			{"type":"regular","gift":{"id":"regular","sticker":{"file_id":"s1","file_unique_id":"u1","type":"regular","width":1,"height":1,"is_animated":false,"is_video":false},"star_count":10},"owned_gift_id":"owned-1","send_date":100,"text":"hello"},
			{"type":"unique","gift":{"gift_id":"regular","base_name":"Base","name":"UniqueName","number":7,"model":{"name":"model","sticker":{"file_id":"s2","file_unique_id":"u2","type":"regular","width":1,"height":1,"is_animated":false,"is_video":false},"rarity_per_mille":5},"symbol":{"name":"symbol","sticker":{"file_id":"s3","file_unique_id":"u3","type":"regular","width":1,"height":1,"is_animated":false,"is_video":false},"rarity_per_mille":6},"backdrop":{"name":"backdrop","colors":{"center_color":1,"edge_color":2,"symbol_color":3,"text_color":4},"rarity_per_mille":7}},"owned_gift_id":"owned-2","send_date":101,"can_be_transferred":true}
		]
	}`)
	var gifts OwnedGifts
	if err := json.Unmarshal(data, &gifts); err != nil {
		t.Fatalf("decode owned gifts: %v", err)
	}
	if gifts.TotalCount != 2 || gifts.NextOffset != "next" || len(gifts.Gifts) != 2 {
		t.Fatalf("unexpected owned gifts: %+v", gifts)
	}
	regular, ok := gifts.Gifts[0].(OwnedGiftRegular)
	if !ok || regular.OwnedGiftID != "owned-1" || regular.Text != "hello" {
		t.Fatalf("unexpected regular gift: %#v", gifts.Gifts[0])
	}
	unique, ok := gifts.Gifts[1].(OwnedGiftUnique)
	if !ok || unique.Gift.Name != "UniqueName" || !unique.CanBeTransferred {
		t.Fatalf("unexpected unique gift: %#v", gifts.Gifts[1])
	}
}

func TestOwnedGiftsRejectUnknownType(t *testing.T) {
	var gifts OwnedGifts
	err := json.Unmarshal([]byte(`{"total_count":1,"gifts":[{"type":"mystery"}]}`), &gifts)
	if err == nil || !strings.Contains(err.Error(), "unsupported owned gift type") {
		t.Fatalf("expected unsupported owned gift error, got %v", err)
	}
}

func TestMessageGiftServiceFieldsDecode(t *testing.T) {
	data := []byte(`{
		"message_id":1,
		"chat":{"id":123,"type":"private"},
		"date":10,
		"gift":{"gift":{"id":"gift-1","sticker":{"file_id":"s1","file_unique_id":"u1","type":"regular","width":1,"height":1,"is_animated":false,"is_video":false},"star_count":10},"owned_gift_id":"owned-1","can_be_upgraded":true},
		"unique_gift":{"gift":{"gift_id":"gift-1","base_name":"Base","name":"UniqueName","number":7,"model":{"name":"model","sticker":{"file_id":"s2","file_unique_id":"u2","type":"regular","width":1,"height":1,"is_animated":false,"is_video":false},"rarity_per_mille":5},"symbol":{"name":"symbol","sticker":{"file_id":"s3","file_unique_id":"u3","type":"regular","width":1,"height":1,"is_animated":false,"is_video":false},"rarity_per_mille":6},"backdrop":{"name":"backdrop","colors":{"center_color":1,"edge_color":2,"symbol_color":3,"text_color":4},"rarity_per_mille":7}},"origin":"upgrade","owned_gift_id":"owned-2"},
		"gift_upgrade_sent":{"gift":{"id":"gift-2","sticker":{"file_id":"s4","file_unique_id":"u4","type":"regular","width":1,"height":1,"is_animated":false,"is_video":false},"star_count":20},"prepaid_upgrade_star_count":30}
	}`)
	var msg Message
	if err := json.Unmarshal(data, &msg); err != nil {
		t.Fatalf("decode message: %v", err)
	}
	if msg.Gift == nil || msg.Gift.OwnedGiftID != "owned-1" || !msg.Gift.CanBeUpgraded {
		t.Fatalf("unexpected gift field: %+v", msg.Gift)
	}
	if msg.UniqueGift == nil || msg.UniqueGift.Origin != "upgrade" || msg.UniqueGift.Gift.Name != "UniqueName" {
		t.Fatalf("unexpected unique gift field: %+v", msg.UniqueGift)
	}
	if msg.GiftUpgradeSent == nil || msg.GiftUpgradeSent.PrepaidUpgradeStarCount != 30 {
		t.Fatalf("unexpected gift upgrade field: %+v", msg.GiftUpgradeSent)
	}
}

func TestTransactionPartnersDecodeGiftFields(t *testing.T) {
	var userPartner TransactionPartnerUser
	if err := json.Unmarshal([]byte(`{"type":"user","transaction_type":"gift_purchase","user":{"id":1,"is_bot":false,"first_name":"Ann"},"gift":{"id":"gift-1","sticker":{"file_id":"s","file_unique_id":"u","type":"regular","width":1,"height":1,"is_animated":false,"is_video":false},"star_count":10}}`), &userPartner); err != nil {
		t.Fatalf("decode user partner: %v", err)
	}
	if userPartner.Gift == nil || userPartner.Gift.ID != "gift-1" {
		t.Fatalf("unexpected user partner gift: %+v", userPartner.Gift)
	}

	partner, err := UnmarshalTransactionPartner([]byte(`{"type":"chat","chat":{"id":-100,"type":"channel"},"gift":{"id":"gift-2","sticker":{"file_id":"s","file_unique_id":"u","type":"regular","width":1,"height":1,"is_animated":false,"is_video":false},"star_count":20}}`))
	if err != nil {
		t.Fatalf("decode chat partner: %v", err)
	}
	chatPartner, ok := partner.(TransactionPartnerChat)
	if !ok || chatPartner.Gift == nil || chatPartner.Gift.ID != "gift-2" {
		t.Fatalf("unexpected chat partner: %#v", partner)
	}
}
