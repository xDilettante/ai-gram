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
	"github.com/xDilettante/ai-gram/telegram"
)

func TestSetPassportDataErrorsSendsPayload(t *testing.T) {
	const token = "123:secret"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/bot"+token+"/setPassportDataErrors" {
			t.Fatalf("unexpected path: %q", r.URL.Path)
		}
		var payload map[string]any
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode payload: %v", err)
		}
		if payload["user_id"] != float64(42) {
			t.Fatalf("unexpected user_id: %#v", payload)
		}
		errorsPayload, ok := payload["errors"].([]any)
		if !ok || len(errorsPayload) != 3 {
			t.Fatalf("unexpected errors payload: %#v", payload)
		}
		first := errorsPayload[0].(map[string]any)
		if first["source"] != "data" || first["type"] != "personal_details" || first["field_name"] != "first_name" || first["data_hash"] != "data-hash" || first["message"] != "Fix first name" {
			t.Fatalf("unexpected data field error: %#v", first)
		}
		second := errorsPayload[1].(map[string]any)
		if second["source"] != "front_side" || second["file_hash"] != "front-hash" {
			t.Fatalf("unexpected front side error: %#v", second)
		}
		third := errorsPayload[2].(map[string]any)
		hashes := third["file_hashes"].([]any)
		if third["source"] != "translation_files" || len(hashes) != 2 || hashes[0] != "translation-hash-1" {
			t.Fatalf("unexpected translation files error: %#v", third)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true,"result":true}`))
	}))
	defer server.Close()

	bot := newTestBot(t, token, server.URL, server.Client())
	ok, err := bot.SetPassportDataErrors(context.Background(), SetPassportDataErrorsParams{
		UserID: 42,
		Errors: []telegram.PassportElementError{
			telegram.PassportElementErrorDataField{Type: "personal_details", FieldName: "first_name", DataHash: "data-hash", Message: "Fix first name"},
			telegram.PassportElementErrorFrontSide{Type: "passport", FileHash: "front-hash", Message: "Upload a sharper scan"},
			telegram.PassportElementErrorTranslationFiles{Type: "passport", FileHashes: []string{"translation-hash-1", "translation-hash-2"}, Message: "Upload readable translations"},
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Fatal("expected true result")
	}
}

func TestSetPassportDataErrorsValidation(t *testing.T) {
	const token = "123:secret"
	bot := newTestBot(t, token, "https://example.test", nil)
	tests := []struct {
		name   string
		params SetPassportDataErrorsParams
	}{
		{name: "invalid user", params: SetPassportDataErrorsParams{Errors: validPassportErrors()}},
		{name: "empty errors", params: SetPassportDataErrorsParams{UserID: 42}},
		{name: "nil error", params: SetPassportDataErrorsParams{UserID: 42, Errors: []telegram.PassportElementError{nil}}},
		{name: "invalid error", params: SetPassportDataErrorsParams{UserID: 42, Errors: []telegram.PassportElementError{telegram.PassportElementErrorDataField{Type: "personal_details", FieldName: "first_name", Message: "Fix"}}}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ok, err := bot.SetPassportDataErrors(context.Background(), tt.params)
			if err == nil {
				t.Fatal("expected error")
			}
			if ok {
				t.Fatal("expected false result")
			}
			assertNoToken(t, err, token)
			assertNoPassportPayload(t, err)
		})
	}
}

func TestSetPassportDataErrorsReturnsAPIError(t *testing.T) {
	const token = "123:secret"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":false,"error_code":400,"description":"Bad Request: passport data errors invalid"}`))
	}))
	defer server.Close()

	bot := newTestBot(t, token, server.URL, server.Client())
	ok, err := bot.SetPassportDataErrors(context.Background(), SetPassportDataErrorsParams{UserID: 42, Errors: validPassportErrors()})
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
	if apiErr.Code != 400 {
		t.Fatalf("unexpected APIError code: %d", apiErr.Code)
	}
	assertNoToken(t, err, token)
	assertNoPassportPayload(t, err)
}

func TestSetPassportDataErrorsInvalidJSON(t *testing.T) {
	const token = "123:secret"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`not-json`))
	}))
	defer server.Close()

	bot := newTestBot(t, token, server.URL, server.Client())
	ok, err := bot.SetPassportDataErrors(context.Background(), SetPassportDataErrorsParams{UserID: 42, Errors: validPassportErrors()})
	if err == nil {
		t.Fatal("expected error")
	}
	if ok {
		t.Fatal("expected false result")
	}
	assertNoToken(t, err, token)
	assertNoPassportPayload(t, err)
}

func TestSetPassportDataErrorsHTTP500(t *testing.T) {
	const token = "123:secret"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "server error", http.StatusInternalServerError)
	}))
	defer server.Close()

	bot := newTestBot(t, token, server.URL, server.Client())
	ok, err := bot.SetPassportDataErrors(context.Background(), SetPassportDataErrorsParams{UserID: 42, Errors: validPassportErrors()})
	if err == nil {
		t.Fatal("expected error")
	}
	if ok {
		t.Fatal("expected false result")
	}
	assertNoToken(t, err, token)
	assertNoPassportPayload(t, err)
}

func TestSetPassportDataErrorsCancelledContext(t *testing.T) {
	const token = "123:secret"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("request should not reach server")
	}))
	defer server.Close()

	bot := newTestBot(t, token, server.URL, server.Client())
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	ok, err := bot.SetPassportDataErrors(ctx, SetPassportDataErrorsParams{UserID: 42, Errors: validPassportErrors()})
	if err == nil {
		t.Fatal("expected error")
	}
	if ok {
		t.Fatal("expected false result")
	}
	assertNoToken(t, err, token)
	assertNoPassportPayload(t, err)
}

func validPassportErrors() []telegram.PassportElementError {
	return []telegram.PassportElementError{
		telegram.PassportElementErrorDataField{Type: "personal_details", FieldName: "first_name", DataHash: "data-hash", Message: "Fix first name"},
	}
}

func assertNoPassportPayload(t *testing.T, err error) {
	t.Helper()
	if err == nil {
		return
	}
	message := err.Error()
	for _, sensitive := range []string{"enc-data", "cred-secret", "data-hash", "front-hash", "translation-hash", "element-hash"} {
		if strings.Contains(message, sensitive) {
			t.Fatalf("passport payload leaked in error: %q", message)
		}
	}
}
