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

func TestAnswerCallbackQuerySendsPayloadAndDecodesResult(t *testing.T) {
	const token = "123:secret"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/bot"+token+"/answerCallbackQuery" {
			t.Fatalf("unexpected path: %q", r.URL.Path)
		}
		var payload map[string]any
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode payload: %v", err)
		}
		if got := payload["callback_query_id"]; got != "callback-id" {
			t.Fatalf("unexpected callback_query_id: %#v", got)
		}
		if got := payload["text"]; got != "Готово" {
			t.Fatalf("unexpected text: %#v", got)
		}
		if got := payload["show_alert"]; got != true {
			t.Fatalf("unexpected show_alert: %#v", got)
		}
		if got := payload["url"]; got != "https://example.com/callback" {
			t.Fatalf("unexpected url: %#v", got)
		}
		if got := payload["cache_time"]; got != float64(30) {
			t.Fatalf("unexpected cache_time: %#v", got)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true,"result":true}`))
	}))
	defer server.Close()

	bot := newTestBot(t, token, server.URL, server.Client())
	ok, err := bot.AnswerCallbackQuery(context.Background(), AnswerCallbackQueryParams{
		CallbackQueryID: "callback-id",
		Text:            "Готово",
		ShowAlert:       true,
		URL:             "https://example.com/callback",
		CacheTime:       30,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Fatal("expected true result")
	}
}

func TestAnswerCallbackQueryMinimalSuccess(t *testing.T) {
	const token = "123:secret"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var payload map[string]any
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode payload: %v", err)
		}
		if len(payload) != 1 || payload["callback_query_id"] != "callback-id" {
			t.Fatalf("unexpected payload: %#v", payload)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true,"result":true}`))
	}))
	defer server.Close()

	bot := newTestBot(t, token, server.URL, server.Client())
	ok, err := bot.AnswerCallbackQuery(context.Background(), AnswerCallbackQueryParams{CallbackQueryID: "callback-id"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Fatal("expected true result")
	}
}

func TestAnswerCallbackQueryValidation(t *testing.T) {
	const token = "123:secret"
	bot := newTestBot(t, token, "https://example.test", nil)
	text200 := strings.Repeat("я", 200)
	text201 := strings.Repeat("я", 201)

	valid := []AnswerCallbackQueryParams{
		{CallbackQueryID: "id", Text: text200},
		{CallbackQueryID: "id", CacheTime: 0},
		{CallbackQueryID: "id", CacheTime: 1},
		{CallbackQueryID: "id", URL: "http://example.com"},
		{CallbackQueryID: "id", URL: "https://example.com"},
	}
	for _, params := range valid {
		if err := params.validate(); err != nil {
			t.Fatalf("unexpected valid params error for %+v: %v", params, err)
		}
	}

	tests := []struct {
		name   string
		params AnswerCallbackQueryParams
	}{
		{name: "empty callback id", params: AnswerCallbackQueryParams{}},
		{name: "long text", params: AnswerCallbackQueryParams{CallbackQueryID: "id", Text: text201}},
		{name: "negative cache time", params: AnswerCallbackQueryParams{CallbackQueryID: "id", CacheTime: -1}},
		{name: "url without host", params: AnswerCallbackQueryParams{CallbackQueryID: "id", URL: "https:///path"}},
		{name: "file url", params: AnswerCallbackQueryParams{CallbackQueryID: "id", URL: "file:///tmp/a"}},
		{name: "malformed url", params: AnswerCallbackQueryParams{CallbackQueryID: "id", URL: "://bad"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ok, err := bot.AnswerCallbackQuery(context.Background(), tt.params)
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

func TestAnswerCallbackQueryReturnsAPIError(t *testing.T) {
	const token = "123:secret"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":false,"error_code":400,"description":"Bad Request"}`))
	}))
	defer server.Close()

	bot := newTestBot(t, token, server.URL, server.Client())
	ok, err := bot.AnswerCallbackQuery(context.Background(), AnswerCallbackQueryParams{CallbackQueryID: "callback-id"})
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

func TestAnswerCallbackQueryResponseAndContextErrors(t *testing.T) {
	const token = "123:secret"

	t.Run("cancelled context", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			t.Fatal("request should not reach server")
		}))
		defer server.Close()
		bot := newTestBot(t, token, server.URL, server.Client())
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		ok, err := bot.AnswerCallbackQuery(ctx, AnswerCallbackQueryParams{CallbackQueryID: "callback-id"})
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
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`not-json`))
		}))
		defer server.Close()
		bot := newTestBot(t, token, server.URL, server.Client())
		ok, err := bot.AnswerCallbackQuery(context.Background(), AnswerCallbackQueryParams{CallbackQueryID: "callback-id"})
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
			http.Error(w, "server error", http.StatusInternalServerError)
		}))
		defer server.Close()
		bot := newTestBot(t, token, server.URL, server.Client())
		ok, err := bot.AnswerCallbackQuery(context.Background(), AnswerCallbackQueryParams{CallbackQueryID: "callback-id"})
		if err == nil {
			t.Fatal("expected error")
		}
		if ok {
			t.Fatal("expected false result")
		}
		assertNoToken(t, err, token)
	})
}
