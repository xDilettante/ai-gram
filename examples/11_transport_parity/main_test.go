package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/xDilettante/ai-gram/middleware"
	"github.com/xDilettante/ai-gram/telegram"
)

func TestParseTransport(t *testing.T) {
	tests := []struct {
		raw     string
		want    string
		wantErr bool
	}{
		{raw: "", want: transportPolling},
		{raw: "polling", want: transportPolling},
		{raw: "POLLING", want: transportPolling},
		{raw: "webhook", want: transportWebhook},
		{raw: "bad", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.raw, func(t *testing.T) {
			got, err := parseTransport(tt.raw)
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Fatalf("parseTransport() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestNormalizePath(t *testing.T) {
	tests := map[string]string{
		"":        "/webhook",
		"hook":    "/hook",
		"/hook":   "/hook",
		" /hook ": "/hook",
	}
	for raw, want := range tests {
		if got := normalizePath(raw); got != want {
			t.Fatalf("normalizePath(%q) = %q, want %q", raw, got, want)
		}
	}
}

func TestStatusTextMasksIDs(t *testing.T) {
	update := telegram.Update{Message: &telegram.Message{
		MessageID: 10,
		From:      &telegram.User{ID: 123456789, FirstName: "Alice"},
		Chat:      telegram.Chat{ID: -1001234567890, Type: "supergroup"},
		Text:      "/status",
	}}

	text := statusText(update, middleware.AccessModePublic)
	for _, raw := range []string{"123456789", "-1001234567890"} {
		if strings.Contains(text, raw) {
			t.Fatalf("statusText leaked raw ID %s in %q", raw, text)
		}
	}
	for _, want := range []string{"access_mode=public", "chat=-10***890:supergroup", "actor=user:123***789"} {
		if !strings.Contains(text, want) {
			t.Fatalf("statusText missing %q in %q", want, text)
		}
	}
}

func TestMuxHealthz(t *testing.T) {
	handler := mux(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusAccepted)
	}), "/telegram")

	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/healthz", nil))
	if recorder.Code != http.StatusOK {
		t.Fatalf("healthz status = %d", recorder.Code)
	}
	if recorder.Body.String() != "ok\n" {
		t.Fatalf("healthz body = %q", recorder.Body.String())
	}

	recorder = httptest.NewRecorder()
	handler.ServeHTTP(recorder, httptest.NewRequest(http.MethodPost, "/telegram", nil))
	if recorder.Code != http.StatusAccepted {
		t.Fatalf("webhook status = %d", recorder.Code)
	}
}

func TestSafeValueFlattensWhitespace(t *testing.T) {
	if got := safeValue("hello\nworld\tok"); got != "hello world ok" {
		t.Fatalf("safeValue() = %q", got)
	}
	if got := safeValue(" "); got != "<none>" {
		t.Fatalf("empty safeValue() = %q", got)
	}
}
