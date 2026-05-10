package telegram

import (
	"encoding/json"
	"testing"
)

func TestUpdateDecodesChatMemberUpdated(t *testing.T) {
	payload := []byte(`{
		"update_id": 300,
		"my_chat_member": {
			"chat": {"id": -100123, "type": "supergroup", "title": "Test"},
			"from": {"id": 7, "is_bot": false, "first_name": "Admin"},
			"date": 1234567890,
			"old_chat_member": {"status": "member", "tag": "old", "user": {"id": 8, "is_bot": false, "first_name": "Member"}},
			"new_chat_member": {"status": "administrator", "user": {"id": 8, "is_bot": false, "first_name": "Member"}, "can_manage_chat": true, "can_manage_tags": true, "can_manage_direct_messages": true},
			"via_join_request": true,
			"via_chat_folder_invite_link": true
		}
	}`)

	var update Update
	if err := json.Unmarshal(payload, &update); err != nil {
		t.Fatalf("decode update: %v", err)
	}
	change := update.MyChatMember
	if change == nil || update.UpdateID != 300 {
		t.Fatalf("unexpected update: %+v", update)
	}
	if change.Chat.ID != -100123 || change.From.ID != 7 || change.Date != 1234567890 || !change.ViaJoinRequest || !change.ViaChatFolderInviteLink {
		t.Fatalf("unexpected chat member update: %+v", change)
	}
	if change.OldChatMember.Tag != "old" || change.NewChatMember.Status != ChatMemberStatusAdministrator || !change.NewChatMember.CanManageTags || !change.NewChatMember.CanManageDirectMessages {
		t.Fatalf("unexpected chat members: old=%+v new=%+v", change.OldChatMember, change.NewChatMember)
	}
	if chat := update.EffectiveChat(); chat == nil || chat.ID != -100123 {
		t.Fatalf("unexpected effective chat: %+v", chat)
	}
	if user := update.EffectiveUser(); user == nil || user.ID != 7 {
		t.Fatalf("unexpected effective user: %+v", user)
	}
}

func TestUpdateDecodesChatMember(t *testing.T) {
	payload := []byte(`{
		"update_id": 301,
		"chat_member": {
			"chat": {"id": -100124, "type": "supergroup", "title": "Test"},
			"from": {"id": 9, "is_bot": false, "first_name": "Admin"},
			"date": 1234567891,
			"old_chat_member": {"status": "restricted", "tag": "guest", "user": {"id": 10, "is_bot": false, "first_name": "Member"}, "is_member": true, "can_send_messages": true, "can_react_to_messages": true, "can_edit_tag": true},
			"new_chat_member": {"status": "kicked", "user": {"id": 10, "is_bot": false, "first_name": "Member"}, "until_date": 1234567999}
		}
	}`)

	var update Update
	if err := json.Unmarshal(payload, &update); err != nil {
		t.Fatalf("decode update: %v", err)
	}
	change := update.ChatMember
	if change == nil || change.Chat.ID != -100124 || change.From.ID != 9 {
		t.Fatalf("unexpected chat member update: %+v", update)
	}
	if !change.OldChatMember.IsMember || !change.OldChatMember.CanSendMessages || !change.OldChatMember.CanReactToMessages || !change.OldChatMember.CanEditTag || change.NewChatMember.UntilDate != 1234567999 {
		t.Fatalf("unexpected chat member fields: old=%+v new=%+v", change.OldChatMember, change.NewChatMember)
	}
}

func TestUpdateDecodesChatBoosts(t *testing.T) {
	payload := []byte(`{
		"update_id": 302,
		"chat_boost": {
			"chat": {"id": -100125, "type": "supergroup", "title": "Boosted"},
			"boost": {
				"boost_id": "boost-1",
				"add_date": 1234567892,
				"expiration_date": 1239999999,
				"source": {"source": "premium", "user": {"id": 11, "is_bot": false, "first_name": "Booster"}}
			}
		}
	}`)

	var update Update
	if err := json.Unmarshal(payload, &update); err != nil {
		t.Fatalf("decode update: %v", err)
	}
	boost := update.ChatBoost
	if boost == nil || boost.Chat.ID != -100125 || boost.Boost.BoostID != "boost-1" || boost.Boost.AddDate != 1234567892 || boost.Boost.ExpirationDate != 1239999999 {
		t.Fatalf("unexpected chat boost: %+v", update)
	}
	source, ok := boost.Boost.Source.(ChatBoostSourcePremium)
	if !ok || source.User.ID != 11 {
		t.Fatalf("unexpected boost source: %#v", boost.Boost.Source)
	}
	if chat := update.EffectiveChat(); chat == nil || chat.ID != -100125 {
		t.Fatalf("unexpected effective chat: %+v", chat)
	}
	if user := update.EffectiveUser(); user == nil || user.ID != 11 {
		t.Fatalf("unexpected effective user: %+v", user)
	}
}

func TestUpdateDecodesRemovedChatBoost(t *testing.T) {
	payload := []byte(`{
		"update_id": 303,
		"removed_chat_boost": {
			"chat": {"id": -100126, "type": "supergroup", "title": "Boosted"},
			"boost_id": "boost-2",
			"remove_date": 1234567893,
			"source": {"source": "giveaway", "giveaway_message_id": 99, "user": {"id": 12, "is_bot": false, "first_name": "Winner"}, "prize_star_count": 1000, "is_unclaimed": true}
		}
	}`)

	var update Update
	if err := json.Unmarshal(payload, &update); err != nil {
		t.Fatalf("decode update: %v", err)
	}
	removed := update.RemovedChatBoost
	if removed == nil || removed.Chat.ID != -100126 || removed.BoostID != "boost-2" || removed.RemoveDate != 1234567893 {
		t.Fatalf("unexpected removed chat boost: %+v", update)
	}
	source, ok := removed.Source.(ChatBoostSourceGiveaway)
	if !ok || source.GiveawayMessageID != 99 || source.User == nil || source.User.ID != 12 || source.PrizeStarCount != 1000 || !source.IsUnclaimed {
		t.Fatalf("unexpected removed boost source: %#v", removed.Source)
	}
}

func TestUserChatBoostsDecodesSourceVariants(t *testing.T) {
	payload := []byte(`{
		"boosts": [
			{"boost_id": "premium", "add_date": 1, "expiration_date": 2, "source": {"source": "premium", "user": {"id": 1, "is_bot": false, "first_name": "A"}}},
			{"boost_id": "gift", "add_date": 3, "expiration_date": 4, "source": {"source": "gift_code", "user": {"id": 2, "is_bot": false, "first_name": "B"}}},
			{"boost_id": "giveaway", "add_date": 5, "expiration_date": 6, "source": {"source": "giveaway", "giveaway_message_id": 7, "prize_star_count": 500}}
		]
	}`)

	var boosts UserChatBoosts
	if err := json.Unmarshal(payload, &boosts); err != nil {
		t.Fatalf("decode user chat boosts: %v", err)
	}
	if len(boosts.Boosts) != 3 {
		t.Fatalf("unexpected boost count: %d", len(boosts.Boosts))
	}
	if _, ok := boosts.Boosts[0].Source.(ChatBoostSourcePremium); !ok {
		t.Fatalf("unexpected premium source: %#v", boosts.Boosts[0].Source)
	}
	if _, ok := boosts.Boosts[1].Source.(ChatBoostSourceGiftCode); !ok {
		t.Fatalf("unexpected gift code source: %#v", boosts.Boosts[1].Source)
	}
	if _, ok := boosts.Boosts[2].Source.(ChatBoostSourceGiveaway); !ok {
		t.Fatalf("unexpected giveaway source: %#v", boosts.Boosts[2].Source)
	}
}

func TestChatBoostRejectsUnknownSource(t *testing.T) {
	var boost ChatBoost
	if err := json.Unmarshal([]byte(`{"boost_id":"x","source":{"source":"unknown"}}`), &boost); err == nil {
		t.Fatal("expected error for unknown chat boost source")
	}
}
