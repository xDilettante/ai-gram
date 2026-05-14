package main

import (
	"strings"
	"testing"

	"github.com/xDilettante/ai-gram/telegram"
)

func TestPreviewTextIsDryRunAndMasksIDs(t *testing.T) {
	preview := moderationPreview{
		Action: "moderator_preview",
		Reason: "spam",
		Chat:   &telegram.Chat{ID: -1001234567890, Type: "supergroup"},
		Reporter: telegram.Actor{
			User: &telegram.User{ID: 123456789, FirstName: "Admin"},
		},
		Target: telegram.Actor{
			User: &telegram.User{ID: 987654321, FirstName: "Target"},
		},
		TargetMessageID: 44,
		DryRun:          true,
	}

	text := previewText(preview)
	for _, raw := range []string{"-1001234567890", "123456789", "987654321"} {
		if strings.Contains(text, raw) {
			t.Fatalf("previewText leaked raw ID %s in %q", raw, text)
		}
	}
	for _, want := range []string{
		"dry_run=true",
		"chat=-10***890:supergroup",
		"reporter=user:123***789",
		"target=user:987***321",
		"would_delete_message=false",
		"would_restrict_user=false",
		"would_ban_user=false",
	} {
		if !strings.Contains(text, want) {
			t.Fatalf("previewText missing %q in %q", want, text)
		}
	}
}

func TestPreviewFromUpdateUsesReplyTarget(t *testing.T) {
	update := telegram.Update{
		Message: &telegram.Message{
			MessageID: 10,
			From:      &telegram.User{ID: 123456789, FirstName: "Reporter"},
			Chat:      telegram.Chat{ID: -1001234567890, Type: "supergroup"},
			Text:      "/report spam",
			ReplyToMessage: &telegram.Message{
				MessageID: 9,
				From:      &telegram.User{ID: 987654321, FirstName: "Target"},
				Chat:      telegram.Chat{ID: -1001234567890, Type: "supergroup"},
			},
		},
	}

	preview := previewFromUpdate(update, "report", "spam")
	if preview.Reporter.User == nil || preview.Reporter.User.ID != 123456789 {
		t.Fatalf("unexpected reporter: %+v", preview.Reporter)
	}
	if preview.Target.User == nil || preview.Target.User.ID != 987654321 {
		t.Fatalf("unexpected target: %+v", preview.Target)
	}
	if preview.TargetMessageID != 9 {
		t.Fatalf("unexpected target message ID: %d", preview.TargetMessageID)
	}
}

func TestSafeValueFlattensWhitespace(t *testing.T) {
	if got := safeValue("spam\nwith\ttabs"); got != "spam with tabs" {
		t.Fatalf("safeValue() = %q", got)
	}
	if got := safeValue(" "); got != "<none>" {
		t.Fatalf("empty safeValue() = %q", got)
	}
}
