package bot

import (
	"context"
	stderrors "errors"
	"net/url"
	"strings"
	"unicode/utf8"
)

// AnswerCallbackQueryParams contains supported parameters for answerCallbackQuery.
type AnswerCallbackQueryParams struct {
	CallbackQueryID string `json:"callback_query_id"`
	Text            string `json:"text,omitempty"`
	ShowAlert       bool   `json:"show_alert,omitempty"`
	URL             string `json:"url,omitempty"`
	CacheTime       int    `json:"cache_time,omitempty"`
}

// AnswerCallbackQuery sends an answer to a callback query from an inline keyboard.
func (b *Bot) AnswerCallbackQuery(ctx context.Context, params AnswerCallbackQueryParams) (bool, error) {
	if err := params.validate(); err != nil {
		return false, err
	}

	var result bool
	if err := b.call(ctx, "answerCallbackQuery", params, &result); err != nil {
		return false, err
	}

	return result, nil
}

func (params AnswerCallbackQueryParams) validate() error {
	if strings.TrimSpace(params.CallbackQueryID) == "" {
		return stderrors.New("callback_query_id is required")
	}
	if utf8.RuneCountInString(params.Text) > 200 {
		return stderrors.New("callback query answer text must be at most 200 characters")
	}
	if params.CacheTime < 0 {
		return stderrors.New("cache_time must not be negative")
	}
	if params.URL != "" {
		if err := validateCallbackAnswerURL(params.URL); err != nil {
			return err
		}
	}

	return nil
}

func validateCallbackAnswerURL(rawURL string) error {
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return stderrors.New("callback query answer URL must be a valid HTTP(S) URL")
	}
	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return stderrors.New("callback query answer URL scheme must be http or https")
	}
	if parsed.Host == "" {
		return stderrors.New("callback query answer URL host is required")
	}

	return nil
}
