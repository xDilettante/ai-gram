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

func TestSetWebhookSendsPayloadAndReturnsTrue(t *testing.T) {
	const token = "123:secret"
	const secret = "webhook_SECRET-123"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/bot"+token+"/setWebhook" {
			t.Fatalf("unexpected path: %q", r.URL.Path)
		}

		var payload map[string]any
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode payload: %v", err)
		}
		if got := payload["url"]; got != "https://example.com/telegram/webhook" {
			t.Fatalf("unexpected url: %#v", got)
		}
		if got := payload["max_connections"]; got != float64(40) {
			t.Fatalf("unexpected max_connections: %#v", got)
		}
		allowed, ok := payload["allowed_updates"].([]any)
		if !ok {
			t.Fatalf("unexpected allowed_updates type: %#v", payload["allowed_updates"])
		}
		if len(allowed) != 2 || allowed[0] != "message" || allowed[1] != "callback_query" {
			t.Fatalf("unexpected allowed_updates: %#v", allowed)
		}
		if got := payload["drop_pending_updates"]; got != true {
			t.Fatalf("unexpected drop_pending_updates: %#v", got)
		}
		if got := payload["secret_token"]; got != secret {
			t.Fatalf("unexpected secret_token: %#v", got)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true,"result":true}`))
	}))
	defer server.Close()

	bot := newTestBot(t, token, server.URL, server.Client())
	ok, err := bot.SetWebhook(context.Background(), SetWebhookParams{
		URL:                "https://example.com/telegram/webhook",
		MaxConnections:     40,
		AllowedUpdates:     []string{"message", "callback_query"},
		DropPendingUpdates: true,
		SecretToken:        secret,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Fatal("expected true result")
	}
}

func TestSetWebhookSendsMultipartCertificate(t *testing.T) {
	const token = "123:secret"
	const secret = "webhook_SECRET-123"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/bot"+token+"/setWebhook" {
			t.Fatalf("unexpected path: %q", r.URL.Path)
		}
		if got := r.Header.Get("Content-Type"); !strings.HasPrefix(got, "multipart/form-data") {
			t.Fatalf("unexpected content type: %q", got)
		}
		if err := r.ParseMultipartForm(1 << 20); err != nil {
			t.Fatalf("parse multipart: %v", err)
		}

		assertMultipartValue(t, r, "url", "https://example.com/telegram/webhook")
		assertMultipartValue(t, r, "ip_address", "203.0.113.10")
		assertMultipartValue(t, r, "max_connections", "40")
		assertMultipartValue(t, r, "allowed_updates", `["message","callback_query"]`)
		assertMultipartValue(t, r, "drop_pending_updates", "true")
		assertMultipartValue(t, r, "secret_token", secret)

		content, header := readMultipartFile(t, r, "certificate")
		if header.Filename != "public-cert.pem" {
			t.Fatalf("unexpected certificate filename: %q", header.Filename)
		}
		if header.Header.Get("Content-Type") != "application/x-pem-file" {
			t.Fatalf("unexpected certificate content type: %q", header.Header.Get("Content-Type"))
		}
		if string(content) != "certificate-data" {
			t.Fatalf("unexpected certificate content: %q", string(content))
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true,"result":true}`))
	}))
	defer server.Close()

	bot := newTestBot(t, token, server.URL, server.Client())
	ok, err := bot.SetWebhook(context.Background(), SetWebhookParams{
		URL:                "https://example.com/telegram/webhook",
		Certificate:        FileUpload(UploadFile{Name: "public-cert.pem", Reader: strings.NewReader("certificate-data"), ContentType: "application/x-pem-file"}),
		IPAddress:          "203.0.113.10",
		MaxConnections:     40,
		AllowedUpdates:     []string{"message", "callback_query"},
		DropPendingUpdates: true,
		SecretToken:        secret,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Fatal("expected true result")
	}
}

func TestSetWebhookValidation(t *testing.T) {
	const secret = "bad secret"
	tests := []struct {
		name    string
		params  SetWebhookParams
		wantErr bool
	}{
		{name: "empty url", params: SetWebhookParams{}, wantErr: true},
		{name: "malformed url", params: SetWebhookParams{URL: "://bad"}, wantErr: true},
		{name: "http url", params: SetWebhookParams{URL: "http://example.com/hook"}, wantErr: true},
		{name: "url without host", params: SetWebhookParams{URL: "https:///hook"}, wantErr: true},
		{name: "max zero", params: SetWebhookParams{URL: "https://example.com/hook"}},
		{name: "max one", params: SetWebhookParams{URL: "https://example.com/hook", MaxConnections: 1}},
		{name: "max hundred", params: SetWebhookParams{URL: "https://example.com/hook", MaxConnections: 100}},
		{name: "max negative", params: SetWebhookParams{URL: "https://example.com/hook", MaxConnections: -1}, wantErr: true},
		{name: "max too high", params: SetWebhookParams{URL: "https://example.com/hook", MaxConnections: 101}, wantErr: true},
		{name: "invalid secret", params: SetWebhookParams{URL: "https://example.com/hook", SecretToken: secret}, wantErr: true},
		{name: "certificate file id", params: SetWebhookParams{URL: "https://example.com/hook", Certificate: FileID("cert-file-id")}, wantErr: true},
		{name: "certificate file url", params: SetWebhookParams{URL: "https://example.com/hook", Certificate: FileURL("https://example.com/cert.pem")}, wantErr: true},
		{name: "certificate missing filename", params: SetWebhookParams{URL: "https://example.com/hook", Certificate: FileUpload(UploadFile{Reader: strings.NewReader("cert")})}, wantErr: true},
		{name: "certificate missing reader", params: SetWebhookParams{URL: "https://example.com/hook", Certificate: FileUpload(UploadFile{Name: "cert.pem"})}, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.params.validate()
			if tt.wantErr && err == nil {
				t.Fatal("expected error")
			}
			if !tt.wantErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if err != nil {
				assertNoSecret(t, err, tt.params.SecretToken)
				assertNoSecret(t, err, secret)
			}
		})
	}
}

func TestSetWebhookAllowsHTTPWithCustomBaseURL(t *testing.T) {
	const token = "123:secret"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var payload map[string]any
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode payload: %v", err)
		}
		if got := payload["url"]; got != "http://127.0.0.1:8080/webhook" {
			t.Fatalf("unexpected url: %#v", got)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true,"result":true}`))
	}))
	defer server.Close()

	bot := newTestBot(t, token, server.URL, server.Client())
	ok, err := bot.SetWebhook(context.Background(), SetWebhookParams{URL: "http://127.0.0.1:8080/webhook"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Fatal("expected true result")
	}
}

func TestSetWebhookReturnsAPIErrorAndRedactsSecrets(t *testing.T) {
	const token = "123:secret"
	const secret = "webhook_SECRET-123"
	const webhookURL = "https://example.com/hook"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":false,"error_code":400,"description":"Bad token 123:secret, secret webhook_SECRET-123, url https://example.com/hook"}`))
	}))
	defer server.Close()

	bot := newTestBot(t, token, server.URL, server.Client())
	ok, err := bot.SetWebhook(context.Background(), SetWebhookParams{URL: webhookURL, SecretToken: secret})
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
	assertNoSecret(t, err, secret)
	if strings.Contains(err.Error(), webhookURL) {
		t.Fatalf("error leaked webhook URL: %q", err.Error())
	}
}

func TestSetWebhookMultipartErrorsDoNotLeakSensitiveValues(t *testing.T) {
	const token = "123:secret"
	const secret = "webhook_SECRET-123"
	const webhookURL = "https://example.com/hook"
	certificate := FileUpload(UploadFile{Name: "public-cert.pem", Reader: strings.NewReader("certificate-data")})

	t.Run("api error", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"ok":false,"error_code":400,"description":"Bad token 123:secret, secret webhook_SECRET-123, url https://example.com/hook"}`))
		}))
		defer server.Close()

		bot := newTestBot(t, token, server.URL, server.Client())
		ok, err := bot.SetWebhook(context.Background(), SetWebhookParams{URL: webhookURL, Certificate: certificate, SecretToken: secret})
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
		assertNoSecret(t, err, secret)
		if strings.Contains(err.Error(), webhookURL) || strings.Contains(err.Error(), "certificate-data") {
			t.Fatalf("error leaked sensitive value: %q", err.Error())
		}
	})

	t.Run("invalid json", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`not-json`))
		}))
		defer server.Close()

		bot := newTestBot(t, token, server.URL, server.Client())
		ok, err := bot.SetWebhook(context.Background(), SetWebhookParams{URL: webhookURL, Certificate: FileUpload(UploadFile{Name: "public-cert.pem", Reader: strings.NewReader("certificate-data")}), SecretToken: secret})
		if err == nil {
			t.Fatal("expected invalid JSON error")
		}
		if ok {
			t.Fatal("expected false result")
		}
		assertNoToken(t, err, token)
		assertNoSecret(t, err, secret)
		if strings.Contains(err.Error(), webhookURL) || strings.Contains(err.Error(), "certificate-data") {
			t.Fatalf("error leaked sensitive value: %q", err.Error())
		}
	})

	t.Run("http status", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "server error", http.StatusInternalServerError)
		}))
		defer server.Close()

		bot := newTestBot(t, token, server.URL, server.Client())
		ok, err := bot.SetWebhook(context.Background(), SetWebhookParams{URL: webhookURL, Certificate: FileUpload(UploadFile{Name: "public-cert.pem", Reader: strings.NewReader("certificate-data")}), SecretToken: secret})
		if err == nil {
			t.Fatal("expected HTTP status error")
		}
		if ok {
			t.Fatal("expected false result")
		}
		assertNoToken(t, err, token)
		assertNoSecret(t, err, secret)
		if strings.Contains(err.Error(), webhookURL) || strings.Contains(err.Error(), "certificate-data") {
			t.Fatalf("error leaked sensitive value: %q", err.Error())
		}
	})

	t.Run("cancelled context", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			t.Fatal("request should not reach server with canceled context")
		}))
		defer server.Close()

		bot := newTestBot(t, token, server.URL, server.Client())
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		ok, err := bot.SetWebhook(ctx, SetWebhookParams{URL: webhookURL, Certificate: FileUpload(UploadFile{Name: "public-cert.pem", Reader: strings.NewReader("certificate-data")}), SecretToken: secret})
		if err == nil {
			t.Fatal("expected context error")
		}
		if ok {
			t.Fatal("expected false result")
		}
		assertNoToken(t, err, token)
		assertNoSecret(t, err, secret)
		if strings.Contains(err.Error(), webhookURL) || strings.Contains(err.Error(), "certificate-data") {
			t.Fatalf("error leaked sensitive value: %q", err.Error())
		}
	})
}

func TestDeleteWebhookSendsPayloadAndReturnsTrue(t *testing.T) {
	const token = "123:secret"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/bot"+token+"/deleteWebhook" {
			t.Fatalf("unexpected path: %q", r.URL.Path)
		}

		var payload map[string]any
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode payload: %v", err)
		}
		if got := payload["drop_pending_updates"]; got != true {
			t.Fatalf("unexpected drop_pending_updates: %#v", got)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true,"result":true}`))
	}))
	defer server.Close()

	bot := newTestBot(t, token, server.URL, server.Client())
	ok, err := bot.DeleteWebhook(context.Background(), DeleteWebhookParams{DropPendingUpdates: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Fatal("expected true result")
	}
}

func TestDeleteWebhookZeroValueParamsAreValid(t *testing.T) {
	const token = "123:secret"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/bot"+token+"/deleteWebhook" {
			t.Fatalf("unexpected path: %q", r.URL.Path)
		}
		var payload map[string]any
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode payload: %v", err)
		}
		if _, ok := payload["drop_pending_updates"]; ok {
			t.Fatalf("drop_pending_updates should be omitted, payload=%#v", payload)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true,"result":true}`))
	}))
	defer server.Close()

	bot := newTestBot(t, token, server.URL, server.Client())
	ok, err := bot.DeleteWebhook(context.Background(), DeleteWebhookParams{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Fatal("expected true result")
	}
}

func TestDeleteWebhookReturnsAPIError(t *testing.T) {
	const token = "123:secret"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":false,"error_code":409,"description":"Conflict"}`))
	}))
	defer server.Close()

	bot := newTestBot(t, token, server.URL, server.Client())
	ok, err := bot.DeleteWebhook(context.Background(), DeleteWebhookParams{})
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
	if apiErr.Code != 409 {
		t.Fatalf("unexpected APIError code: %d", apiErr.Code)
	}
	assertNoToken(t, err, token)
}

func TestGetWebhookInfoDecodesInfo(t *testing.T) {
	const token = "123:secret"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/bot"+token+"/getWebhookInfo" {
			t.Fatalf("unexpected path: %q", r.URL.Path)
		}
		var payload map[string]any
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode payload: %v", err)
		}
		if len(payload) != 0 {
			t.Fatalf("unexpected payload: %#v", payload)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true,"result":{"url":"https://example.com/hook","has_custom_certificate":false,"pending_update_count":3,"last_error_message":"temporary failure","max_connections":40,"allowed_updates":["message","callback_query"]}}`))
	}))
	defer server.Close()

	bot := newTestBot(t, token, server.URL, server.Client())
	info, err := bot.GetWebhookInfo(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if info == nil {
		t.Fatal("expected webhook info")
	}
	if info.URL != "https://example.com/hook" || info.PendingUpdateCount != 3 || info.LastErrorMessage != "temporary failure" || info.MaxConnections != 40 {
		t.Fatalf("unexpected webhook info: %+v", info)
	}
	if len(info.AllowedUpdates) != 2 || info.AllowedUpdates[0] != "message" || info.AllowedUpdates[1] != "callback_query" {
		t.Fatalf("unexpected allowed updates: %#v", info.AllowedUpdates)
	}
}

func TestGetWebhookInfoReturnsAPIError(t *testing.T) {
	const token = "123:secret"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":false,"error_code":401,"description":"Unauthorized"}`))
	}))
	defer server.Close()

	bot := newTestBot(t, token, server.URL, server.Client())
	info, err := bot.GetWebhookInfo(context.Background())
	if err == nil {
		t.Fatal("expected error")
	}
	if info != nil {
		t.Fatalf("expected nil info, got %+v", info)
	}

	var apiErr *apierrors.APIError
	if !stderrors.As(err, &apiErr) {
		t.Fatalf("expected APIError, got %T", err)
	}
	if apiErr.Code != 401 {
		t.Fatalf("unexpected APIError code: %d", apiErr.Code)
	}
	assertNoToken(t, err, token)
}

func TestWebhookMethodsReturnContextError(t *testing.T) {
	const token = "123:secret"
	const secret = "webhook_SECRET-123"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("request should not reach server with canceled context")
	}))
	defer server.Close()

	bot := newTestBot(t, token, server.URL, server.Client())
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	ok, err := bot.SetWebhook(ctx, SetWebhookParams{URL: "https://example.com/hook", SecretToken: secret})
	if err == nil {
		t.Fatal("expected error")
	}
	if ok {
		t.Fatal("expected false result")
	}
	assertNoToken(t, err, token)
	assertNoSecret(t, err, secret)
}

func TestWebhookMethodsReturnInvalidJSONError(t *testing.T) {
	const token = "123:secret"
	const secret = "webhook_SECRET-123"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`not-json`))
	}))
	defer server.Close()

	bot := newTestBot(t, token, server.URL, server.Client())
	ok, err := bot.SetWebhook(context.Background(), SetWebhookParams{URL: "https://example.com/hook", SecretToken: secret})
	if err == nil {
		t.Fatal("expected error")
	}
	if ok {
		t.Fatal("expected false result")
	}
	assertNoToken(t, err, token)
	assertNoSecret(t, err, secret)
}

func TestSecretTokenIsNotLeakedFromSetWebhookValidation(t *testing.T) {
	const token = "123:secret"
	const secret = "bad secret"
	bot := newTestBot(t, token, "https://example.test", nil)

	ok, err := bot.SetWebhook(context.Background(), SetWebhookParams{URL: "https://example.com/hook", SecretToken: secret})
	if err == nil {
		t.Fatal("expected error")
	}
	if ok {
		t.Fatal("expected false result")
	}
	assertNoToken(t, err, token)
	assertNoSecret(t, err, secret)
}

func assertNoSecret(t *testing.T, err error, secret string) {
	t.Helper()
	if secret != "" && strings.Contains(err.Error(), secret) {
		t.Fatalf("error leaked secret: %q", err.Error())
	}
}
