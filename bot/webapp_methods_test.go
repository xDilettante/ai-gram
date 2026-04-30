package bot

import (
	"context"
	"encoding/json"
	stderrors "errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	apierrors "github.com/xDilettante/ai-gram/errors"
)

func TestAnswerWebAppQuerySendsPayloadAndDecodesResult(t *testing.T) {
	const token = "123:secret"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/bot"+token+"/answerWebAppQuery" {
			t.Fatalf("unexpected path: %q", r.URL.Path)
		}
		var payload map[string]any
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode payload: %v", err)
		}
		if payload["web_app_query_id"] != "query-id" {
			t.Fatalf("unexpected web_app_query_id: %#v", payload)
		}
		result, ok := payload["result"].(map[string]any)
		if !ok {
			t.Fatalf("expected result object: %#v", payload["result"])
		}
		if result["type"] != "article" || result["id"] != "article-id" || result["title"] != "Title" {
			t.Fatalf("unexpected result: %#v", result)
		}
		content, ok := result["input_message_content"].(map[string]any)
		if !ok || content["message_text"] != "Web App response" {
			t.Fatalf("unexpected input_message_content: %#v", result["input_message_content"])
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true,"result":{"inline_message_id":"inline-id"}}`))
	}))
	defer server.Close()

	bot := newTestBot(t, token, server.URL, server.Client())
	message, err := bot.AnswerWebAppQuery(context.Background(), AnswerWebAppQueryParams{
		WebAppQueryID: "query-id",
		Result:        InlineArticle("article-id", "Title", InputText("Web App response")),
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if message == nil || message.InlineMessageID != "inline-id" {
		t.Fatalf("unexpected sent message: %+v", message)
	}
}

func TestAnswerWebAppQueryValidation(t *testing.T) {
	const token = "123:secret"
	bot := newTestBot(t, token, "https://example.test", nil)
	tests := []struct {
		name   string
		params AnswerWebAppQueryParams
	}{
		{name: "missing query id", params: AnswerWebAppQueryParams{Result: InlineArticle("article-id", "Title", InputText("ok"))}},
		{name: "nil result", params: AnswerWebAppQueryParams{WebAppQueryID: "query-id"}},
		{name: "invalid result", params: AnswerWebAppQueryParams{WebAppQueryID: "query-id", Result: InlineArticle("", "Title", InputText("ok"))}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			message, err := bot.AnswerWebAppQuery(context.Background(), tt.params)
			if err == nil {
				t.Fatal("expected error")
			}
			if message != nil {
				t.Fatalf("expected nil message, got %+v", message)
			}
			assertNoToken(t, err, token)
			assertNoWebAppPayload(t, err)
		})
	}
}

func TestAnswerWebAppQueryReturnsAPIError(t *testing.T) {
	const token = "123:secret"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/bot"+token+"/answerWebAppQuery" {
			t.Fatalf("unexpected path: %q", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":false,"error_code":400,"description":"Bad Request"}`))
	}))
	defer server.Close()

	bot := newTestBot(t, token, server.URL, server.Client())
	message, err := bot.AnswerWebAppQuery(context.Background(), validAnswerWebAppQueryParams())
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
	assertNoWebAppPayload(t, err)
}

func TestAnswerWebAppQueryResponseAndContextErrors(t *testing.T) {
	const token = "123:secret"
	tests := []struct {
		name     string
		response func(http.ResponseWriter)
	}{
		{name: "invalid json", response: func(w http.ResponseWriter) {
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`not-json`))
		}},
		{name: "http status", response: func(w http.ResponseWriter) {
			http.Error(w, "server error", http.StatusInternalServerError)
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path != "/bot"+token+"/answerWebAppQuery" {
					t.Fatalf("unexpected path: %q", r.URL.Path)
				}
				tt.response(w)
			}))
			defer server.Close()

			bot := newTestBot(t, token, server.URL, server.Client())
			_, err := bot.AnswerWebAppQuery(context.Background(), validAnswerWebAppQueryParams())
			if err == nil {
				t.Fatal("expected error")
			}
			assertNoToken(t, err, token)
			assertNoWebAppPayload(t, err)
		})
	}

	t.Run("cancelled context", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			t.Fatal("request should not reach server")
		}))
		defer server.Close()

		bot := newTestBot(t, token, server.URL, server.Client())
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		_, err := bot.AnswerWebAppQuery(ctx, validAnswerWebAppQueryParams())
		if err == nil {
			t.Fatal("expected error")
		}
		assertNoToken(t, err, token)
		assertNoWebAppPayload(t, err)
	})
}

func validAnswerWebAppQueryParams() AnswerWebAppQueryParams {
	return AnswerWebAppQueryParams{
		WebAppQueryID: "query-id",
		Result:        InlineArticle("article-id", "Title", InputText("Web App response")),
	}
}

func assertNoWebAppPayload(t *testing.T, err error) {
	t.Helper()
	if strings.Contains(err.Error(), "Web App response") || strings.Contains(err.Error(), "opaque-web-app-data") {
		t.Fatalf("error leaked Web App payload: %q", err.Error())
	}
}
