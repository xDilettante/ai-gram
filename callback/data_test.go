package callback

import (
	"errors"
	"testing"
	"time"

	"github.com/xDilettante/ai-gram/telegram"
)

func TestEncodeParseRoundTrip(t *testing.T) {
	expiresAt := time.Unix(1_800_000_000, 0)
	original := New("panel", "open").WithID("item-42").WithPage(3).WithExpiry(expiresAt)

	encoded, err := original.Encode()
	if err != nil {
		t.Fatalf("Encode() error = %v", err)
	}
	if encoded != "panel:open:i=item-42:p=3:e=1800000000" {
		t.Fatalf("encoded = %q", encoded)
	}

	parsed, err := Parse(encoded)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}
	if parsed.Namespace != original.Namespace ||
		parsed.Action != original.Action ||
		parsed.ID != original.ID ||
		parsed.Page != original.Page ||
		parsed.HasPage != original.HasPage ||
		!parsed.ExpiresAt.Equal(original.ExpiresAt) {
		t.Fatalf("parsed = %+v, want %+v", parsed, original)
	}
}

func TestEncodeRejectsInvalidData(t *testing.T) {
	tests := []struct {
		name string
		data Data
		want error
	}{
		{name: "missing namespace", data: New("", "open"), want: ErrInvalidFormat},
		{name: "missing action", data: New("panel", ""), want: ErrInvalidFormat},
		{name: "invalid namespace", data: New("bad:value", "open"), want: ErrInvalidSegment},
		{name: "invalid id", data: New("panel", "open").WithID("bad value"), want: ErrInvalidSegment},
		{name: "negative page", data: New("panel", "open").WithPage(-1), want: ErrInvalidPage},
		{name: "invalid expiry", data: New("panel", "open").WithExpiry(time.Unix(0, 0)), want: ErrInvalidExpiry},
		{name: "too long", data: New("very-long-namespace", "very-long-action").WithID("very-long-item-identifier-that-does-not-fit"), want: ErrTooLong},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.data.Encode()
			if !errors.Is(err, tt.want) {
				t.Fatalf("Encode() error = %v, want %v", err, tt.want)
			}
		})
	}
}

func TestParseRejectsInvalidData(t *testing.T) {
	tests := []struct {
		name string
		raw  string
		want error
	}{
		{name: "empty", raw: "", want: ErrEmptyData},
		{name: "missing action", raw: "panel", want: ErrInvalidFormat},
		{name: "invalid segment", raw: "pan el:open", want: ErrInvalidSegment},
		{name: "unknown field", raw: "panel:open:x=1", want: ErrInvalidFormat},
		{name: "empty field value", raw: "panel:open:i=", want: ErrInvalidFormat},
		{name: "duplicate id", raw: "panel:open:i=1:i=2", want: ErrInvalidFormat},
		{name: "bad page", raw: "panel:open:p=abc", want: ErrInvalidPage},
		{name: "negative page", raw: "panel:open:p=-1", want: ErrInvalidPage},
		{name: "bad expiry", raw: "panel:open:e=abc", want: ErrInvalidExpiry},
		{name: "too long", raw: "panel:open:i=very-long-item-identifier-that-does-not-fit-into-telegram-callback-data", want: ErrTooLong},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Parse(tt.raw)
			if !errors.Is(err, tt.want) {
				t.Fatalf("Parse() error = %v, want %v", err, tt.want)
			}
		})
	}
}

func TestHelpers(t *testing.T) {
	now := time.Unix(1_700_000_000, 0)
	data := Confirm("panel", "item-1").WithPage(2).WithTTL(now, time.Minute)
	if !data.Match("panel", ActionConfirm) {
		t.Fatal("expected confirm callback to match")
	}
	if data.Expired(now.Add(59 * time.Second)) {
		t.Fatal("callback should not be expired before expiry")
	}
	if !data.Expired(now.Add(time.Minute)) {
		t.Fatal("callback should be expired at expiry")
	}

	next := data.NextPage()
	if next.Page != 3 || !next.HasPage {
		t.Fatalf("NextPage() = %+v", next)
	}
	previous := next.PreviousPage()
	if previous.Page != 2 || !previous.HasPage {
		t.Fatalf("PreviousPage() = %+v", previous)
	}
}

func TestButton(t *testing.T) {
	button, err := Button("Open", New("panel", "open").WithID("42"))
	if err != nil {
		t.Fatalf("Button() error = %v", err)
	}
	if button != (telegram.InlineKeyboardButton{Text: "Open", CallbackData: "panel:open:i=42"}) {
		t.Fatalf("button = %+v", button)
	}
}
