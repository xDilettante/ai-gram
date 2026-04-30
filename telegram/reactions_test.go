package telegram

import (
	"encoding/json"
	"testing"
)

func TestReactionTypeMarshal(t *testing.T) {
	tests := []struct {
		name     string
		reaction ReactionType
		want     map[string]any
	}{
		{name: "emoji", reaction: NewReactionTypeEmoji("👍"), want: map[string]any{"type": "emoji", "emoji": "👍"}},
		{name: "custom emoji", reaction: NewReactionTypeCustomEmoji("custom-id"), want: map[string]any{"type": "custom_emoji", "custom_emoji_id": "custom-id"}},
		{name: "paid", reaction: NewReactionTypePaid(), want: map[string]any{"type": "paid"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := json.Marshal(tt.reaction)
			if err != nil {
				t.Fatalf("marshal reaction: %v", err)
			}
			var got map[string]any
			if err := json.Unmarshal(data, &got); err != nil {
				t.Fatalf("decode marshaled reaction: %v", err)
			}
			for key, value := range tt.want {
				if got[key] != value {
					t.Fatalf("unexpected %s: got %#v want %#v in %#v", key, got[key], value, got)
				}
			}
		})
	}
}

func TestUnmarshalReactionType(t *testing.T) {
	tests := []struct {
		name    string
		payload string
		check   func(*testing.T, ReactionType)
	}{
		{name: "emoji", payload: `{"type":"emoji","emoji":"👍"}`, check: func(t *testing.T, reaction ReactionType) {
			value, ok := reaction.(ReactionTypeEmoji)
			if !ok || value.Emoji != "👍" || value.Type != "emoji" {
				t.Fatalf("unexpected reaction: %#v", reaction)
			}
		}},
		{name: "custom emoji", payload: `{"type":"custom_emoji","custom_emoji_id":"custom-id"}`, check: func(t *testing.T, reaction ReactionType) {
			value, ok := reaction.(ReactionTypeCustomEmoji)
			if !ok || value.CustomEmojiID != "custom-id" || value.Type != "custom_emoji" {
				t.Fatalf("unexpected reaction: %#v", reaction)
			}
		}},
		{name: "paid", payload: `{"type":"paid"}`, check: func(t *testing.T, reaction ReactionType) {
			value, ok := reaction.(ReactionTypePaid)
			if !ok || value.Type != "paid" {
				t.Fatalf("unexpected reaction: %#v", reaction)
			}
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reaction, err := UnmarshalReactionType([]byte(tt.payload))
			if err != nil {
				t.Fatalf("unmarshal reaction: %v", err)
			}
			tt.check(t, reaction)
		})
	}

	if _, err := UnmarshalReactionType([]byte(`{"type":"unknown"}`)); err == nil {
		t.Fatal("expected error for unknown reaction type")
	}
}

func TestValidateReactionTypes(t *testing.T) {
	if err := ValidateReactionTypes(nil); err != nil {
		t.Fatalf("nil reactions should be valid: %v", err)
	}
	if err := ValidateReactionTypes([]ReactionType{NewReactionTypeEmoji("👍"), NewReactionTypeCustomEmoji("custom-id"), NewReactionTypePaid()}); err != nil {
		t.Fatalf("valid reactions rejected: %v", err)
	}

	tests := []struct {
		name      string
		reactions []ReactionType
	}{
		{name: "nil reaction", reactions: []ReactionType{nil}},
		{name: "empty emoji", reactions: []ReactionType{NewReactionTypeEmoji("")}},
		{name: "empty custom emoji", reactions: []ReactionType{NewReactionTypeCustomEmoji("")}},
	}
	var typedNil *ReactionTypeEmoji
	tests = append(tests, struct {
		name      string
		reactions []ReactionType
	}{name: "typed nil", reactions: []ReactionType{typedNil}})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ValidateReactionTypes(tt.reactions); err == nil {
				t.Fatal("expected validation error")
			}
		})
	}
}

func TestUpdateDecodesMessageReaction(t *testing.T) {
	payload := []byte(`{
		"update_id": 200,
		"message_reaction": {
			"chat": {"id": -100123, "type": "supergroup", "title": "Test"},
			"message_id": 10,
			"user": {"id": 7, "is_bot": false, "first_name": "Alice"},
			"date": 1234567890,
			"old_reaction": [{"type":"emoji","emoji":"👍"}],
			"new_reaction": [{"type":"custom_emoji","custom_emoji_id":"custom-id"},{"type":"paid"}]
		}
	}`)

	var update Update
	if err := json.Unmarshal(payload, &update); err != nil {
		t.Fatalf("decode update: %v", err)
	}
	change := update.MessageReaction
	if update.UpdateID != 200 || change == nil {
		t.Fatalf("unexpected update: %+v", update)
	}
	if change.Chat.ID != -100123 || change.MessageID != 10 || change.User == nil || change.User.ID != 7 || change.Date != 1234567890 {
		t.Fatalf("unexpected message reaction: %+v", change)
	}
	if _, ok := change.OldReaction[0].(ReactionTypeEmoji); !ok {
		t.Fatalf("unexpected old reaction: %#v", change.OldReaction)
	}
	if _, ok := change.NewReaction[0].(ReactionTypeCustomEmoji); !ok {
		t.Fatalf("unexpected first new reaction: %#v", change.NewReaction)
	}
	if _, ok := change.NewReaction[1].(ReactionTypePaid); !ok {
		t.Fatalf("unexpected second new reaction: %#v", change.NewReaction)
	}
	if chat := update.EffectiveChat(); chat == nil || chat.ID != -100123 {
		t.Fatalf("unexpected effective chat: %+v", chat)
	}
	if user := update.EffectiveUser(); user == nil || user.ID != 7 {
		t.Fatalf("unexpected effective user: %+v", user)
	}
}

func TestUpdateDecodesMessageReactionCount(t *testing.T) {
	payload := []byte(`{
		"update_id": 201,
		"message_reaction_count": {
			"chat": {"id": -100456, "type": "supergroup", "title": "Test"},
			"message_id": 11,
			"date": 1234567999,
			"reactions": [
				{"type":{"type":"emoji","emoji":"🔥"},"total_count":3},
				{"type":{"type":"paid"},"total_count":1}
			]
		}
	}`)

	var update Update
	if err := json.Unmarshal(payload, &update); err != nil {
		t.Fatalf("decode update: %v", err)
	}
	count := update.MessageReactionCount
	if update.UpdateID != 201 || count == nil {
		t.Fatalf("unexpected update: %+v", update)
	}
	if count.Chat.ID != -100456 || count.MessageID != 11 || count.Date != 1234567999 || len(count.Reactions) != 2 {
		t.Fatalf("unexpected reaction count update: %+v", count)
	}
	if reaction, ok := count.Reactions[0].Type.(ReactionTypeEmoji); !ok || reaction.Emoji != "🔥" || count.Reactions[0].TotalCount != 3 {
		t.Fatalf("unexpected first reaction count: %+v", count.Reactions[0])
	}
	if _, ok := count.Reactions[1].Type.(ReactionTypePaid); !ok || count.Reactions[1].TotalCount != 1 {
		t.Fatalf("unexpected second reaction count: %+v", count.Reactions[1])
	}
	if chat := update.EffectiveChat(); chat == nil || chat.ID != -100456 {
		t.Fatalf("unexpected effective chat: %+v", chat)
	}
	if user := update.EffectiveUser(); user != nil {
		t.Fatalf("reaction count update should not have effective user: %+v", user)
	}
}

func TestReactionCountRejectsUnknownType(t *testing.T) {
	var count ReactionCount
	if err := json.Unmarshal([]byte(`{"type":{"type":"unknown"},"total_count":1}`), &count); err == nil {
		t.Fatal("expected error for unknown reaction count type")
	}
}
