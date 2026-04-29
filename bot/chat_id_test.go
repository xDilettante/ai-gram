package bot

import (
	"encoding/json"
	"testing"
)

func TestChatIDIntMarshalsAsJSONNumber(t *testing.T) {
	got, err := json.Marshal(ChatIDInt(12345))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(got) != "12345" {
		t.Fatalf("unexpected JSON: %s", got)
	}
}

func TestChatIDStringMarshalsAsJSONString(t *testing.T) {
	got, err := json.Marshal(ChatIDString("@channel"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(got) != `"@channel"` {
		t.Fatalf("unexpected JSON: %s", got)
	}
}

func TestEmptyChatIDIsInvalid(t *testing.T) {
	if (ChatID{}).valid() {
		t.Fatal("expected zero ChatID to be invalid")
	}
	if ChatIDString("").valid() {
		t.Fatal("expected empty string ChatID to be invalid")
	}
	if ChatIDString(" \t").valid() {
		t.Fatal("expected blank string ChatID to be invalid")
	}
}
