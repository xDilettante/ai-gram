package bot

import (
	"context"
	"encoding/json"
	stderrors "errors"
	"net/http"
	"net/http/httptest"
	"testing"

	apierrors "github.com/xDilettante/ai-gram/errors"
)

func TestAnswerGuestQuerySendsPayloadAndDecodesResult(t *testing.T) {
	const token = "123:secret"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/bot"+token+"/answerGuestQuery" {
			t.Fatalf("unexpected path: %q", r.URL.Path)
		}
		var payload map[string]any
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode payload: %v", err)
		}
		if payload["guest_query_id"] != "guest-query-id" {
			t.Fatalf("unexpected guest_query_id: %#v", payload)
		}
		result, ok := payload["result"].(map[string]any)
		if !ok || result["type"] != "article" || result["id"] != "article-1" || result["title"] != "Article" {
			t.Fatalf("unexpected result: %#v", payload["result"])
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true,"result":{"inline_message_id":"inline-message-id"}}`))
	}))
	defer server.Close()

	bot := newTestBot(t, token, server.URL, server.Client())
	message, err := bot.AnswerGuestQuery(context.Background(), AnswerGuestQueryParams{
		GuestQueryID: "guest-query-id",
		Result:       InlineArticle("article-1", "Article", InputText("hello")),
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if message == nil || message.InlineMessageID != "inline-message-id" {
		t.Fatalf("unexpected message: %+v", message)
	}
}

func TestAnswerGuestQueryValidation(t *testing.T) {
	const token = "123:secret"
	bot := newTestBot(t, token, "https://example.test", nil)
	tests := []struct {
		name   string
		params AnswerGuestQueryParams
	}{
		{name: "empty query id", params: AnswerGuestQueryParams{Result: InlineArticle("article-1", "Article", InputText("hello"))}},
		{name: "invalid result", params: AnswerGuestQueryParams{GuestQueryID: "guest-query-id", Result: InlineQueryResultArticle{Title: "Article", InputMessageContent: InputText("hello")}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			message, err := bot.AnswerGuestQuery(context.Background(), tt.params)
			if err == nil {
				t.Fatal("expected error")
			}
			if message != nil {
				t.Fatalf("expected nil message, got %+v", message)
			}
			assertNoToken(t, err, token)
		})
	}
}

func TestAnswerGuestQueryReturnsAPIError(t *testing.T) {
	const token = "123:secret"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/bot"+token+"/answerGuestQuery" {
			t.Fatalf("unexpected path: %q", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":false,"error_code":400,"description":"Bad Request"}`))
	}))
	defer server.Close()

	bot := newTestBot(t, token, server.URL, server.Client())
	message, err := bot.AnswerGuestQuery(context.Background(), AnswerGuestQueryParams{
		GuestQueryID: "guest-query-id",
		Result:       InlineArticle("article-1", "Article", InputText("hello")),
	})
	if err == nil {
		t.Fatal("expected error")
	}
	if message != nil {
		t.Fatalf("expected nil message, got %+v", message)
	}
	var apiErr *apierrors.APIError
	if !stderrors.As(err, &apiErr) {
		t.Fatalf("expected APIError, got %T", err)
	}
	assertNoToken(t, err, token)
}

func TestAnswerGuestQueryResponseAndContextErrors(t *testing.T) {
	const token = "123:secret"
	t.Run("cancelled context", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			t.Fatal("request should not reach server")
		}))
		defer server.Close()
		bot := newTestBot(t, token, server.URL, server.Client())
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		message, err := bot.AnswerGuestQuery(ctx, AnswerGuestQueryParams{
			GuestQueryID: "guest-query-id",
			Result:       InlineArticle("article-1", "Article", InputText("hello")),
		})
		if err == nil {
			t.Fatal("expected error")
		}
		if message != nil {
			t.Fatalf("expected nil message, got %+v", message)
		}
		assertNoToken(t, err, token)
	})

	t.Run("invalid json", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/bot"+token+"/answerGuestQuery" {
				t.Fatalf("unexpected path: %q", r.URL.Path)
			}
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`not-json`))
		}))
		defer server.Close()
		bot := newTestBot(t, token, server.URL, server.Client())
		message, err := bot.AnswerGuestQuery(context.Background(), AnswerGuestQueryParams{
			GuestQueryID: "guest-query-id",
			Result:       InlineArticle("article-1", "Article", InputText("hello")),
		})
		if err == nil {
			t.Fatal("expected error")
		}
		if message != nil {
			t.Fatalf("expected nil message, got %+v", message)
		}
		assertNoToken(t, err, token)
	})
}
