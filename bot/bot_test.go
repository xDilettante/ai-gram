package bot

import (
	"context"
	stderrors "errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"strings"
	"testing"

	apierrors "ai-gram/errors"
)

func TestNewRejectsEmptyToken(t *testing.T) {
	bot, err := New(BotConfig{})
	if err == nil {
		t.Fatal("expected error for empty token")
	}
	if bot != nil {
		t.Fatal("expected nil bot for empty token")
	}
}

func TestNewCreatesBotWithToken(t *testing.T) {
	bot, err := New(BotConfig{Token: "123:abc"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if bot == nil {
		t.Fatal("expected bot")
	}
}

func TestBotDoesNotExposeRawTokenMethod(t *testing.T) {
	if _, ok := reflect.TypeOf(&Bot{}).MethodByName("Token"); ok {
		t.Fatal("Bot must not expose raw token through a public Token method")
	}
}

func TestNewCreatesBotWithBaseURLAndHTTPClient(t *testing.T) {
	httpClient := &http.Client{}
	bot, err := New(BotConfig{
		Token:      "123:abc",
		BaseURL:    "https://example.test/",
		HTTPClient: httpClient,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if bot == nil {
		t.Fatal("expected bot")
	}
	if bot.baseURL != "https://example.test" {
		t.Fatalf("unexpected base URL: %q", bot.baseURL)
	}
	if bot.client == nil {
		t.Fatal("expected internal HTTP client")
	}
}

func TestCallDecodesResult(t *testing.T) {
	const token = "123:secret"
	var sawRequest bool

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sawRequest = true
		if r.Method != http.MethodPost {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/bot"+token+"/testMethod" {
			t.Fatalf("unexpected path: %q", r.URL.Path)
		}
		if strings.Contains(r.URL.Path, "//") {
			t.Fatalf("endpoint path contains double slash: %q", r.URL.Path)
		}
		if got := r.Header.Get("Content-Type"); got != "application/json" {
			t.Fatalf("unexpected content type: %q", got)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true,"result":{"id":42,"name":"decoded"}}`))
	}))
	defer server.Close()

	bot := newTestBot(t, token, server.URL+"/", server.Client())
	var result struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}

	err := bot.call(context.Background(), "testMethod", map[string]string{"hello": "world"}, &result)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !sawRequest {
		t.Fatal("expected request")
	}
	if result.ID != 42 || result.Name != "decoded" {
		t.Fatalf("unexpected result: %+v", result)
	}
}

func TestCallReturnsAPIError(t *testing.T) {
	const token = "123:secret"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":false,"error_code":400,"description":"Bad Request"}`))
	}))
	defer server.Close()

	bot := newTestBot(t, token, server.URL, server.Client())
	err := bot.call(context.Background(), "testMethod", nil, nil)
	if err == nil {
		t.Fatal("expected error")
	}

	var apiErr *apierrors.APIError
	if !stderrors.As(err, &apiErr) {
		t.Fatalf("expected APIError, got %T", err)
	}
	if apiErr.Code != 400 || apiErr.Description != "Bad Request" {
		t.Fatalf("unexpected APIError: %+v", apiErr)
	}
	assertNoToken(t, err, token)
}

func TestCallRejectsEmptyMethod(t *testing.T) {
	const token = "123:secret"
	bot := newTestBot(t, token, "https://example.test", nil)

	err := bot.call(context.Background(), "", nil, nil)
	if err == nil {
		t.Fatal("expected error")
	}
	assertNoToken(t, err, token)
}

func TestCallReturnsErrorForInvalidJSON(t *testing.T) {
	const token = "123:secret"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`not-json`))
	}))
	defer server.Close()

	bot := newTestBot(t, token, server.URL, server.Client())
	err := bot.call(context.Background(), "testMethod", nil, nil)
	if err == nil {
		t.Fatal("expected error")
	}
	assertNoToken(t, err, token)
}

func TestCallReturnsErrorForHTTPStatus500(t *testing.T) {
	const token = "123:secret"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "server error", http.StatusInternalServerError)
	}))
	defer server.Close()

	bot := newTestBot(t, token, server.URL, server.Client())
	err := bot.call(context.Background(), "testMethod", nil, nil)
	if err == nil {
		t.Fatal("expected error")
	}
	assertNoToken(t, err, token)
}

func TestCallRedactsTokenFromNetworkErrors(t *testing.T) {
	const token = "123:secret"
	client := &http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
		return nil, &url.Error{Op: "Post", URL: req.URL.String(), Err: stderrors.New("network down")}
	})}
	bot := newTestBot(t, token, "https://example.test", client)

	err := bot.call(context.Background(), "testMethod", nil, nil)
	if err == nil {
		t.Fatal("expected error")
	}
	assertNoToken(t, err, token)
	if strings.Contains(fmt.Sprintf("%+v", stderrors.Unwrap(err)), token) {
		t.Fatal("unwrapped error leaked token")
	}
}

func TestBotStringRedactsToken(t *testing.T) {
	const token = "123:secret"
	bot := newTestBot(t, token, "https://example.test", nil)

	for _, formatted := range []string{fmt.Sprint(bot), fmt.Sprintf("%#v", bot)} {
		if strings.Contains(formatted, token) {
			t.Fatalf("formatted bot leaked token: %q", formatted)
		}
	}
}

func newTestBot(t *testing.T, token string, baseURL string, client *http.Client) *Bot {
	t.Helper()

	bot, err := New(BotConfig{Token: token, BaseURL: baseURL, HTTPClient: client})
	if err != nil {
		t.Fatalf("unexpected New error: %v", err)
	}

	return bot
}

func assertNoToken(t *testing.T, err error, token string) {
	t.Helper()

	if strings.Contains(err.Error(), token) {
		t.Fatalf("error leaked token: %q", err.Error())
	}
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}
