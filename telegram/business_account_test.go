package telegram

import (
	"encoding/json"
	"testing"
)

func TestStoryAreaTypesMarshalAndValidate(t *testing.T) {
	area := StoryArea{
		Position: StoryAreaPosition{XPercentage: 50, YPercentage: 50, WidthPercentage: 25, HeightPercentage: 25},
		Type:     NewStoryAreaTypeSuggestedReaction(NewReactionTypeEmoji("👍")),
	}
	if err := ValidateStoryArea(area); err != nil {
		t.Fatalf("validate story area: %v", err)
	}
	body, err := json.Marshal(area)
	if err != nil {
		t.Fatalf("marshal story area: %v", err)
	}
	var payload map[string]any
	if err := json.Unmarshal(body, &payload); err != nil {
		t.Fatalf("decode story area: %v", err)
	}
	areaType := payload["type"].(map[string]any)
	if areaType["type"] != StoryAreaTypeSuggestedReactionType {
		t.Fatalf("unexpected story area type: %#v", areaType)
	}
	if err := ValidateStoryArea(StoryArea{Type: NewStoryAreaTypeLink("")}); err == nil {
		t.Fatal("expected invalid link story area")
	}
	if err := ValidateAcceptedGiftTypes(AcceptedGiftTypes{}); err == nil {
		t.Fatal("expected invalid accepted gift types")
	}
}

func TestSuggestedPostServiceMessagesDecode(t *testing.T) {
	body := []byte(`{
		"message_id":10,
		"chat":{"id":123,"type":"channel"},
		"date":1,
		"suggested_post_approved":{"price":{"currency":"XTR","amount":5},"send_date":777},
		"suggested_post_approval_failed":{"price":{"currency":"TON","amount":10000000}},
		"suggested_post_declined":{"comment":"not now"},
		"suggested_post_paid":{"currency":"XTR","star_amount":{"amount":5,"nanostar_amount":1}},
		"suggested_post_refunded":{"reason":"payment_refunded"}
	}`)
	var message Message
	if err := json.Unmarshal(body, &message); err != nil {
		t.Fatalf("decode message: %v", err)
	}
	if message.SuggestedPostApproved == nil || message.SuggestedPostApproved.Price == nil || message.SuggestedPostApproved.Price.Amount != 5 || message.SuggestedPostApproved.SendDate != 777 {
		t.Fatalf("unexpected approved suggested post: %+v", message.SuggestedPostApproved)
	}
	if message.SuggestedPostApprovalFailed == nil || message.SuggestedPostApprovalFailed.Price.Currency != "TON" {
		t.Fatalf("unexpected failed suggested post: %+v", message.SuggestedPostApprovalFailed)
	}
	if message.SuggestedPostDeclined == nil || message.SuggestedPostDeclined.Comment != "not now" {
		t.Fatalf("unexpected declined suggested post: %+v", message.SuggestedPostDeclined)
	}
	if message.SuggestedPostPaid == nil || message.SuggestedPostPaid.StarAmount == nil || message.SuggestedPostPaid.StarAmount.Amount != 5 {
		t.Fatalf("unexpected paid suggested post: %+v", message.SuggestedPostPaid)
	}
	if message.SuggestedPostRefunded == nil || message.SuggestedPostRefunded.Reason != "payment_refunded" {
		t.Fatalf("unexpected refunded suggested post: %+v", message.SuggestedPostRefunded)
	}
}
