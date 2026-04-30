package telegram

import "testing"

func TestValidateLinkPreviewOptions(t *testing.T) {
	if err := ValidateLinkPreviewOptions(nil); err != nil {
		t.Fatalf("nil options should be valid: %v", err)
	}
	if err := ValidateLinkPreviewOptions(&LinkPreviewOptions{URL: "https://example.com", PreferSmallMedia: true}); err != nil {
		t.Fatalf("valid options rejected: %v", err)
	}

	tests := []struct {
		name    string
		options *LinkPreviewOptions
	}{
		{name: "bad url", options: &LinkPreviewOptions{URL: "ftp://example.com"}},
		{name: "conflicting media preference", options: &LinkPreviewOptions{PreferSmallMedia: true, PreferLargeMedia: true}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ValidateLinkPreviewOptions(tt.options); err == nil {
				t.Fatal("expected validation error")
			}
		})
	}
}

func TestValidateInlineQueryResultsButton(t *testing.T) {
	if err := ValidateInlineQueryResultsButton(nil); err != nil {
		t.Fatalf("nil button should be valid: %v", err)
	}
	if err := ValidateInlineQueryResultsButton(&InlineQueryResultsButton{Text: "Open", StartParameter: "start_1"}); err != nil {
		t.Fatalf("valid start parameter button rejected: %v", err)
	}
	if err := ValidateInlineQueryResultsButton(&InlineQueryResultsButton{Text: "Open", WebApp: &WebAppInfo{URL: "https://example.com/app"}}); err != nil {
		t.Fatalf("valid web app button rejected: %v", err)
	}

	tests := []struct {
		name   string
		button *InlineQueryResultsButton
	}{
		{name: "empty text", button: &InlineQueryResultsButton{StartParameter: "start"}},
		{name: "missing action", button: &InlineQueryResultsButton{Text: "Open"}},
		{name: "two actions", button: &InlineQueryResultsButton{Text: "Open", WebApp: &WebAppInfo{URL: "https://example.com"}, StartParameter: "start"}},
		{name: "empty web app url", button: &InlineQueryResultsButton{Text: "Open", WebApp: &WebAppInfo{}}},
		{name: "bad web app url", button: &InlineQueryResultsButton{Text: "Open", WebApp: &WebAppInfo{URL: "ftp://example.com"}}},
		{name: "bad start parameter", button: &InlineQueryResultsButton{Text: "Open", StartParameter: "bad value"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ValidateInlineQueryResultsButton(tt.button); err == nil {
				t.Fatal("expected validation error")
			}
		})
	}
}
