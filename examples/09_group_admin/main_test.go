package main

import (
	"strings"
	"testing"

	"github.com/xDilettante/ai-gram/telegram"
)

func TestActorTextMasksNumericIDs(t *testing.T) {
	text := actorText(telegram.Actor{
		User:           &telegram.User{ID: 123456789, Username: "alice"},
		Chat:           &telegram.Chat{ID: -1001234567890, Type: "supergroup", Username: "group"},
		AnonymousAdmin: true,
	})

	for _, raw := range []string{"123456789", "-1001234567890"} {
		if strings.Contains(text, raw) {
			t.Fatalf("actorText leaked raw ID %s in %q", raw, text)
		}
	}
	for _, want := range []string{"user_id=123***789", "chat_id=-10***890", "anonymous_admin=true"} {
		if !strings.Contains(text, want) {
			t.Fatalf("actorText missing %q in %q", want, text)
		}
	}
}

func TestChatTextMasksNumericID(t *testing.T) {
	text := chatText(&telegram.Chat{ID: -1001234567890, Type: "supergroup", IsForum: true})
	if strings.Contains(text, "-1001234567890") {
		t.Fatalf("chatText leaked raw ID in %q", text)
	}
	for _, want := range []string{"chat_id=-10***890", "chat_type=supergroup", "chat_is_forum=true"} {
		if !strings.Contains(text, want) {
			t.Fatalf("chatText missing %q in %q", want, text)
		}
	}
}

func TestSafeUsernameAndValue(t *testing.T) {
	if got := safeUsername("alice"); got != "@alice" {
		t.Fatalf("safeUsername() = %q", got)
	}
	if got := safeUsername(""); got != "<none>" {
		t.Fatalf("empty safeUsername() = %q", got)
	}
	if got := safeValue(" "); got != "<none>" {
		t.Fatalf("empty safeValue() = %q", got)
	}
}
