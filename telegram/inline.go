package telegram

import (
	stderrors "errors"
	"strings"
)

// LinkPreviewOptions describes link preview generation options.
type LinkPreviewOptions struct {
	IsDisabled       bool   `json:"is_disabled,omitempty"`
	URL              string `json:"url,omitempty"`
	PreferSmallMedia bool   `json:"prefer_small_media,omitempty"`
	PreferLargeMedia bool   `json:"prefer_large_media,omitempty"`
	ShowAboveText    bool   `json:"show_above_text,omitempty"`
}

// InlineQueryResultsButton represents a button shown above inline query results.
type InlineQueryResultsButton struct {
	Text           string      `json:"text"`
	WebApp         *WebAppInfo `json:"web_app,omitempty"`
	StartParameter string      `json:"start_parameter,omitempty"`
}

// ValidateLinkPreviewOptions checks whether options can be sent to Telegram.
func ValidateLinkPreviewOptions(options *LinkPreviewOptions) error {
	if options == nil {
		return nil
	}
	if options.URL != "" {
		if err := validateHTTPURL(options.URL, "link preview URL"); err != nil {
			return err
		}
	}
	if options.PreferSmallMedia && options.PreferLargeMedia {
		return stderrors.New("link preview cannot prefer both small and large media")
	}

	return nil
}

// ValidateInlineQueryResultsButton checks whether button can be sent to Telegram.
func ValidateInlineQueryResultsButton(button *InlineQueryResultsButton) error {
	if button == nil {
		return nil
	}
	if strings.TrimSpace(button.Text) == "" {
		return stderrors.New("inline query results button text is required")
	}

	actions := 0
	if button.WebApp != nil {
		actions++
		if err := validateWebAppInfo(*button.WebApp, "inline query results button web_app"); err != nil {
			return err
		}
	}
	if button.StartParameter != "" {
		actions++
		if err := validateInlineQueryStartParameter(button.StartParameter); err != nil {
			return err
		}
	}
	if actions != 1 {
		return stderrors.New("inline query results button must have exactly one action")
	}

	return nil
}

func validateInlineQueryStartParameter(value string) error {
	if len(value) == 0 || len(value) > 64 {
		return stderrors.New("inline query results button start_parameter must be 1-64 bytes")
	}
	for _, char := range []byte(value) {
		if char >= 'A' && char <= 'Z' {
			continue
		}
		if char >= 'a' && char <= 'z' {
			continue
		}
		if char >= '0' && char <= '9' {
			continue
		}
		if char == '_' || char == '-' {
			continue
		}
		return stderrors.New("inline query results button start_parameter contains unsupported characters")
	}

	return nil
}
