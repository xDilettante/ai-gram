package bot

import (
	"context"
	"encoding/json"
	stderrors "errors"
	"net/http"
	"net/http/httptest"
	"testing"

	apierrors "github.com/xDilettante/ai-gram/errors"
	"github.com/xDilettante/ai-gram/telegram"
)

func TestInputTextMessageContentMarshalAndValidation(t *testing.T) {
	content := InputTextMessageContent{
		MessageText: "hello",
		ParseMode:   "HTML",
		LinkPreviewOptions: &telegram.LinkPreviewOptions{
			IsDisabled: true,
		},
	}
	data, err := json.Marshal(content)
	if err != nil {
		t.Fatalf("marshal content: %v", err)
	}
	var got map[string]any
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("decode content: %v", err)
	}
	if got["message_text"] != "hello" || got["parse_mode"] != "HTML" {
		t.Fatalf("unexpected content payload: %#v", got)
	}
	if _, ok := got["link_preview_options"].(map[string]any); !ok {
		t.Fatalf("expected link_preview_options object: %#v", got)
	}
	if err := validateInputMessageContent(content); err != nil {
		t.Fatalf("valid content rejected: %v", err)
	}

	tests := []struct {
		name    string
		content InputMessageContent
	}{
		{name: "empty text", content: InputTextMessageContent{}},
		{name: "parse mode and entities", content: InputTextMessageContent{MessageText: "hello", ParseMode: "HTML", Entities: []telegram.MessageEntity{{Type: telegram.EntityBold, Offset: 0, Length: 5}}}},
		{name: "invalid link preview", content: InputTextMessageContent{MessageText: "hello", LinkPreviewOptions: &telegram.LinkPreviewOptions{URL: "ftp://example.com"}}},
		{name: "unsupported content", content: unsupportedInputMessageContent{}},
	}
	var typedNil *InputTextMessageContent
	tests = append(tests, struct {
		name    string
		content InputMessageContent
	}{name: "typed nil", content: typedNil})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateInputMessageContent(tt.content); err == nil {
				t.Fatal("expected validation error")
			}
		})
	}
}

func TestInlineQueryResultArticleMarshalAndValidation(t *testing.T) {
	markup := telegram.NewInlineKeyboard([]telegram.InlineKeyboardButton{telegram.InlineButtonCallback("Open", "open")})
	result := InlineArticle("article-1", "Article", InputText("hello"))
	result.ReplyMarkup = &markup
	result.URL = "https://example.com/article"
	result.Description = "description"
	result.ThumbnailURL = "https://example.com/thumb.jpg"
	result.ThumbnailWidth = 100
	result.ThumbnailHeight = 50

	data, err := json.Marshal(result)
	if err != nil {
		t.Fatalf("marshal article: %v", err)
	}
	var got map[string]any
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("decode article: %v", err)
	}
	if got["type"] != "article" || got["id"] != "article-1" || got["title"] != "Article" {
		t.Fatalf("unexpected article payload: %#v", got)
	}
	content, ok := got["input_message_content"].(map[string]any)
	if !ok || content["message_text"] != "hello" {
		t.Fatalf("unexpected input_message_content: %#v", got["input_message_content"])
	}
	if err := validateInlineQueryResult(result); err != nil {
		t.Fatalf("valid article rejected: %v", err)
	}

	tests := []struct {
		name   string
		result InlineQueryResult
	}{
		{name: "empty id", result: InlineQueryResultArticle{Title: "Article", InputMessageContent: InputText("hello")}},
		{name: "empty title", result: InlineQueryResultArticle{ID: "article-1", InputMessageContent: InputText("hello")}},
		{name: "missing content", result: InlineQueryResultArticle{ID: "article-1", Title: "Article"}},
		{name: "negative thumbnail width", result: InlineQueryResultArticle{ID: "article-1", Title: "Article", InputMessageContent: InputText("hello"), ThumbnailWidth: -1}},
		{name: "negative thumbnail height", result: InlineQueryResultArticle{ID: "article-1", Title: "Article", InputMessageContent: InputText("hello"), ThumbnailHeight: -1}},
		{name: "invalid reply markup", result: InlineQueryResultArticle{ID: "article-1", Title: "Article", InputMessageContent: InputText("hello"), ReplyMarkup: &telegram.InlineKeyboardMarkup{}}},
		{name: "invalid url", result: InlineQueryResultArticle{ID: "article-1", Title: "Article", InputMessageContent: InputText("hello"), URL: "ftp://example.com"}},
		{name: "invalid thumbnail url", result: InlineQueryResultArticle{ID: "article-1", Title: "Article", InputMessageContent: InputText("hello"), ThumbnailURL: "ftp://example.com"}},
		{name: "unsupported result", result: unsupportedInlineQueryResult{}},
	}
	var typedNil *InlineQueryResultArticle
	tests = append(tests, struct {
		name   string
		result InlineQueryResult
	}{name: "typed nil", result: typedNil})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateInlineQueryResult(tt.result); err == nil {
				t.Fatal("expected validation error")
			}
		})
	}
}

func TestAnswerInlineQuerySendsPayloadAndDecodesResult(t *testing.T) {
	const token = "123:secret"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/bot"+token+"/answerInlineQuery" {
			t.Fatalf("unexpected path: %q", r.URL.Path)
		}
		var payload map[string]any
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode payload: %v", err)
		}
		if payload["inline_query_id"] != "inline-query-id" || payload["cache_time"] != float64(10) || payload["is_personal"] != true || payload["next_offset"] != "next" {
			t.Fatalf("unexpected payload: %#v", payload)
		}
		results, ok := payload["results"].([]any)
		if !ok || len(results) != 1 {
			t.Fatalf("unexpected results: %#v", payload["results"])
		}
		article, _ := results[0].(map[string]any)
		content, _ := article["input_message_content"].(map[string]any)
		if article["type"] != "article" || article["id"] != "article-1" || article["title"] != "Article" || content["message_text"] != "hello" {
			t.Fatalf("unexpected article result: %#v", article)
		}
		button, _ := payload["button"].(map[string]any)
		if button["text"] != "Open" || button["start_parameter"] != "start_1" {
			t.Fatalf("unexpected button: %#v", button)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true,"result":true}`))
	}))
	defer server.Close()

	bot := newTestBot(t, token, server.URL, server.Client())
	ok, err := bot.AnswerInlineQuery(context.Background(), AnswerInlineQueryParams{
		InlineQueryID: "inline-query-id",
		Results:       []InlineQueryResult{InlineArticle("article-1", "Article", InputText("hello"))},
		CacheTime:     10,
		IsPersonal:    true,
		NextOffset:    "next",
		Button:        &telegram.InlineQueryResultsButton{Text: "Open", StartParameter: "start_1"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Fatal("expected true result")
	}
}

func TestAnswerInlineQueryAllowsEmptyResults(t *testing.T) {
	const token = "123:secret"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/bot"+token+"/answerInlineQuery" {
			t.Fatalf("unexpected path: %q", r.URL.Path)
		}
		var payload map[string]any
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode payload: %v", err)
		}
		results, ok := payload["results"].([]any)
		if !ok || len(results) != 0 {
			t.Fatalf("empty results should be encoded as an array: %#v", payload["results"])
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true,"result":true}`))
	}))
	defer server.Close()

	bot := newTestBot(t, token, server.URL, server.Client())
	ok, err := bot.AnswerInlineQuery(context.Background(), AnswerInlineQueryParams{InlineQueryID: "inline-query-id"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Fatal("expected true result")
	}
}

func TestAnswerInlineQueryValidation(t *testing.T) {
	const token = "123:secret"
	bot := newTestBot(t, token, "https://example.test", nil)
	tooManyResults := make([]InlineQueryResult, 51)
	for i := range tooManyResults {
		tooManyResults[i] = InlineArticle("article", "Article", InputText("hello"))
	}
	tests := []struct {
		name   string
		params AnswerInlineQueryParams
	}{
		{name: "empty inline query id", params: AnswerInlineQueryParams{}},
		{name: "too many results", params: AnswerInlineQueryParams{InlineQueryID: "inline", Results: tooManyResults}},
		{name: "invalid result", params: AnswerInlineQueryParams{InlineQueryID: "inline", Results: []InlineQueryResult{InlineQueryResultArticle{Title: "Article", InputMessageContent: InputText("hello")}}}},
		{name: "negative cache time", params: AnswerInlineQueryParams{InlineQueryID: "inline", CacheTime: -1}},
		{name: "long next offset", params: AnswerInlineQueryParams{InlineQueryID: "inline", NextOffset: string(make([]byte, 65))}},
		{name: "invalid button", params: AnswerInlineQueryParams{InlineQueryID: "inline", Button: &telegram.InlineQueryResultsButton{Text: "Open"}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ok, err := bot.AnswerInlineQuery(context.Background(), tt.params)
			if err == nil {
				t.Fatal("expected error")
			}
			if ok {
				t.Fatal("expected false result")
			}
			assertNoToken(t, err, token)
		})
	}
}

func TestAnswerInlineQueryReturnsAPIError(t *testing.T) {
	const token = "123:secret"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/bot"+token+"/answerInlineQuery" {
			t.Fatalf("unexpected path: %q", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":false,"error_code":400,"description":"Bad Request"}`))
	}))
	defer server.Close()

	bot := newTestBot(t, token, server.URL, server.Client())
	ok, err := bot.AnswerInlineQuery(context.Background(), AnswerInlineQueryParams{InlineQueryID: "inline-query-id"})
	if err == nil {
		t.Fatal("expected error")
	}
	if ok {
		t.Fatal("expected false result")
	}
	var apiErr *apierrors.APIError
	if !stderrors.As(err, &apiErr) {
		t.Fatalf("expected APIError, got %T", err)
	}
	assertNoToken(t, err, token)
}

func TestAnswerInlineQueryResponseAndContextErrors(t *testing.T) {
	const token = "123:secret"
	t.Run("cancelled context", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			t.Fatal("request should not reach server")
		}))
		defer server.Close()
		bot := newTestBot(t, token, server.URL, server.Client())
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		ok, err := bot.AnswerInlineQuery(ctx, AnswerInlineQueryParams{InlineQueryID: "inline-query-id"})
		if err == nil {
			t.Fatal("expected error")
		}
		if ok {
			t.Fatal("expected false result")
		}
		assertNoToken(t, err, token)
	})

	t.Run("invalid json", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/bot"+token+"/answerInlineQuery" {
				t.Fatalf("unexpected path: %q", r.URL.Path)
			}
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`not-json`))
		}))
		defer server.Close()
		bot := newTestBot(t, token, server.URL, server.Client())
		ok, err := bot.AnswerInlineQuery(context.Background(), AnswerInlineQueryParams{InlineQueryID: "inline-query-id"})
		if err == nil {
			t.Fatal("expected error")
		}
		if ok {
			t.Fatal("expected false result")
		}
		assertNoToken(t, err, token)
	})

	t.Run("http status", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/bot"+token+"/answerInlineQuery" {
				t.Fatalf("unexpected path: %q", r.URL.Path)
			}
			http.Error(w, "server error", http.StatusInternalServerError)
		}))
		defer server.Close()
		bot := newTestBot(t, token, server.URL, server.Client())
		ok, err := bot.AnswerInlineQuery(context.Background(), AnswerInlineQueryParams{InlineQueryID: "inline-query-id"})
		if err == nil {
			t.Fatal("expected error")
		}
		if ok {
			t.Fatal("expected false result")
		}
		assertNoToken(t, err, token)
	})
}

type unsupportedInputMessageContent struct{}

func (unsupportedInputMessageContent) inputMessageContent() {}

type unsupportedInlineQueryResult struct{}

func (unsupportedInlineQueryResult) inlineQueryResult() {}
