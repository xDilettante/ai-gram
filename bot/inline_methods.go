package bot

import (
	"context"
	"encoding/json"
	stderrors "errors"
	"fmt"
	"net/url"
	"reflect"
	"strings"

	"github.com/xDilettante/ai-gram/telegram"
)

const inlineQueryResultArticleType = "article"

// InputMessageContent marks Telegram input message content objects used by inline results.
type InputMessageContent interface {
	inputMessageContent()
}

// InputTextMessageContent describes text content to be sent when an inline result is chosen.
type InputTextMessageContent struct {
	MessageText        string                       `json:"message_text"`
	ParseMode          string                       `json:"parse_mode,omitempty"`
	Entities           []telegram.MessageEntity     `json:"entities,omitempty"`
	LinkPreviewOptions *telegram.LinkPreviewOptions `json:"link_preview_options,omitempty"`
}

func (InputTextMessageContent) inputMessageContent() {}

// InputText creates text content for an inline query result.
func InputText(message string) InputTextMessageContent {
	return InputTextMessageContent{MessageText: message}
}

// InlineQueryResult marks Telegram inline query result objects.
type InlineQueryResult interface {
	inlineQueryResult()
}

// InlineQueryResultArticle represents an article or web page inline query result.
type InlineQueryResultArticle struct {
	Type                string                         `json:"type"`
	ID                  string                         `json:"id"`
	Title               string                         `json:"title"`
	InputMessageContent InputMessageContent            `json:"input_message_content"`
	ReplyMarkup         *telegram.InlineKeyboardMarkup `json:"reply_markup,omitempty"`
	URL                 string                         `json:"url,omitempty"`
	HideURL             bool                           `json:"hide_url,omitempty"`
	Description         string                         `json:"description,omitempty"`
	ThumbnailURL        string                         `json:"thumbnail_url,omitempty"`
	ThumbnailWidth      int                            `json:"thumbnail_width,omitempty"`
	ThumbnailHeight     int                            `json:"thumbnail_height,omitempty"`
}

func (InlineQueryResultArticle) inlineQueryResult() {}

// InlineArticle creates an article inline query result with text or another input message content.
func InlineArticle(id string, title string, content InputMessageContent) InlineQueryResultArticle {
	return InlineQueryResultArticle{Type: inlineQueryResultArticleType, ID: id, Title: title, InputMessageContent: content}
}

// MarshalJSON encodes InlineQueryResultArticle with the required Telegram type field.
func (result InlineQueryResultArticle) MarshalJSON() ([]byte, error) {
	result.Type = inlineQueryResultArticleType
	type article InlineQueryResultArticle
	return json.Marshal(article(result))
}

// AnswerInlineQueryParams contains supported parameters for answerInlineQuery.
type AnswerInlineQueryParams struct {
	InlineQueryID string                             `json:"inline_query_id"`
	Results       []InlineQueryResult                `json:"results"`
	CacheTime     int                                `json:"cache_time,omitempty"`
	IsPersonal    bool                               `json:"is_personal,omitempty"`
	NextOffset    string                             `json:"next_offset,omitempty"`
	Button        *telegram.InlineQueryResultsButton `json:"button,omitempty"`
}

// MarshalJSON encodes AnswerInlineQueryParams while preserving an empty results array.
func (params AnswerInlineQueryParams) MarshalJSON() ([]byte, error) {
	type payload AnswerInlineQueryParams
	encoded := payload(params)
	if encoded.Results == nil {
		encoded.Results = []InlineQueryResult{}
	}
	return json.Marshal(encoded)
}

// AnswerInlineQuery sends answers to an inline query.
func (b *Bot) AnswerInlineQuery(ctx context.Context, params AnswerInlineQueryParams) (bool, error) {
	if err := params.validate(); err != nil {
		return false, err
	}

	var result bool
	if err := b.call(ctx, "answerInlineQuery", params, &result); err != nil {
		return false, err
	}

	return result, nil
}

func (params AnswerInlineQueryParams) validate() error {
	if strings.TrimSpace(params.InlineQueryID) == "" {
		return stderrors.New("inline_query_id is required")
	}
	if len(params.Results) > 50 {
		return stderrors.New("inline query results must contain at most 50 items")
	}
	for index, result := range params.Results {
		if err := validateInlineQueryResult(result); err != nil {
			return fmt.Errorf("results[%d] is invalid: %w", index, err)
		}
	}
	if params.CacheTime < 0 {
		return stderrors.New("cache_time must not be negative")
	}
	if len([]byte(params.NextOffset)) > 64 {
		return stderrors.New("next_offset must be at most 64 bytes")
	}
	if err := telegram.ValidateInlineQueryResultsButton(params.Button); err != nil {
		return err
	}

	return nil
}

func validateInlineQueryResult(result InlineQueryResult) error {
	if result == nil || isNilBotInterfaceValue(result) {
		return stderrors.New("inline query result must not be nil")
	}

	switch value := result.(type) {
	case InlineQueryResultArticle:
		return validateInlineQueryResultArticle(value)
	case *InlineQueryResultArticle:
		return validateInlineQueryResultArticle(*value)
	default:
		return stderrors.New("unsupported inline query result")
	}
}

func validateInlineQueryResultArticle(result InlineQueryResultArticle) error {
	if strings.TrimSpace(result.ID) == "" {
		return stderrors.New("inline query result article id is required")
	}
	if strings.TrimSpace(result.Title) == "" {
		return stderrors.New("inline query result article title is required")
	}
	if err := validateInputMessageContent(result.InputMessageContent); err != nil {
		return err
	}
	if result.ReplyMarkup != nil {
		if err := telegram.ValidateReplyMarkup(*result.ReplyMarkup); err != nil {
			return err
		}
	}
	if result.URL != "" {
		if err := validateInlineHTTPURL(result.URL, "inline query result article URL"); err != nil {
			return err
		}
	}
	if result.ThumbnailURL != "" {
		if err := validateInlineHTTPURL(result.ThumbnailURL, "inline query result article thumbnail_url"); err != nil {
			return err
		}
	}
	if result.ThumbnailWidth < 0 {
		return stderrors.New("thumbnail_width must not be negative")
	}
	if result.ThumbnailHeight < 0 {
		return stderrors.New("thumbnail_height must not be negative")
	}

	return nil
}

func validateInputMessageContent(content InputMessageContent) error {
	if content == nil || isNilBotInterfaceValue(content) {
		return stderrors.New("input_message_content is required")
	}

	switch value := content.(type) {
	case InputTextMessageContent:
		return validateInputTextMessageContent(value)
	case *InputTextMessageContent:
		return validateInputTextMessageContent(*value)
	default:
		return stderrors.New("unsupported input_message_content")
	}
}

func validateInputTextMessageContent(content InputTextMessageContent) error {
	if strings.TrimSpace(content.MessageText) == "" {
		return stderrors.New("message_text is required")
	}
	if err := validateEntityFormatting(content.ParseMode, content.Entities); err != nil {
		return err
	}
	if err := telegram.ValidateLinkPreviewOptions(content.LinkPreviewOptions); err != nil {
		return err
	}

	return nil
}

func validateInlineHTTPURL(rawURL string, field string) error {
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return stderrors.New(field + " must be a valid HTTP(S) URL")
	}
	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return stderrors.New(field + " scheme must be http or https")
	}
	if parsed.Host == "" {
		return stderrors.New(field + " host is required")
	}

	return nil
}

func isNilBotInterfaceValue(value any) bool {
	reflected := reflect.ValueOf(value)
	switch reflected.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Pointer, reflect.Slice:
		return reflected.IsNil()
	default:
		return false
	}
}
