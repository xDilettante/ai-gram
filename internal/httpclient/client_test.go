package httpclient

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestDoReadsAndClosesBody(t *testing.T) {
	body := &closeTrackingBody{Reader: strings.NewReader("ok")}
	client := New(doerFunc(func(*http.Request) (*http.Response, error) {
		return &http.Response{StatusCode: http.StatusOK, Body: body}, nil
	}))
	req, err := http.NewRequest(http.MethodPost, "https://example.test", nil)
	if err != nil {
		t.Fatalf("unexpected request error: %v", err)
	}

	got, err := client.Do(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(got) != "ok" {
		t.Fatalf("unexpected body: %q", got)
	}
	if !body.closed {
		t.Fatal("expected response body to be closed")
	}
}

func TestDoReturnsStatusError(t *testing.T) {
	client := New(doerFunc(func(*http.Request) (*http.Response, error) {
		return &http.Response{StatusCode: http.StatusInternalServerError, Body: io.NopCloser(strings.NewReader("fail"))}, nil
	}))
	req, err := http.NewRequest(http.MethodPost, "https://example.test", nil)
	if err != nil {
		t.Fatalf("unexpected request error: %v", err)
	}

	err = func() error {
		_, err := client.Do(context.Background(), req)
		return err
	}()
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "500") {
		t.Fatalf("expected status in error, got %q", err.Error())
	}
}

type closeTrackingBody struct {
	*strings.Reader
	closed bool
}

func (b *closeTrackingBody) Close() error {
	b.closed = true
	return nil
}

type doerFunc func(*http.Request) (*http.Response, error)

func (f doerFunc) Do(req *http.Request) (*http.Response, error) {
	return f(req)
}

func TestNewWithNilHTTPClientUsesDefault(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"ok":true}`))
	}))
	defer server.Close()

	client := New((*http.Client)(nil))
	req, err := http.NewRequest(http.MethodGet, server.URL, nil)
	if err != nil {
		t.Fatalf("new request: %v", err)
	}
	body, err := client.Do(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(body) != `{"ok":true}` {
		t.Fatalf("unexpected body: %s", body)
	}
}

func TestDoPassesThroughTelegramErrorPayload(t *testing.T) {
	payload := `{"ok":false,"error_code":400,"description":"Bad Request: message is not modified"}`
	client := New(doerFunc(func(*http.Request) (*http.Response, error) {
		return &http.Response{StatusCode: http.StatusBadRequest, Body: io.NopCloser(strings.NewReader(payload))}, nil
	}))
	req, err := http.NewRequest(http.MethodPost, "https://example.test/bot123:secret/editMessageText", nil)
	if err != nil {
		t.Fatalf("unexpected request error: %v", err)
	}

	got, err := client.Do(context.Background(), req)
	if err != nil {
		t.Fatalf("expected telegram error payload to pass through, got error: %v", err)
	}
	if string(got) != payload {
		t.Fatalf("unexpected body: %q", got)
	}
}

func TestDoNonJSONErrorIncludesMethodWithoutToken(t *testing.T) {
	client := New(doerFunc(func(*http.Request) (*http.Response, error) {
		return &http.Response{StatusCode: http.StatusBadGateway, Body: io.NopCloser(strings.NewReader("<html>bad gateway</html>"))}, nil
	}))
	req, err := http.NewRequest(http.MethodPost, "https://example.test/bot123:secret/sendMessage", nil)
	if err != nil {
		t.Fatalf("unexpected request error: %v", err)
	}

	_, err = client.Do(context.Background(), req)
	if err == nil {
		t.Fatal("expected error")
	}
	message := err.Error()
	if !strings.Contains(message, "sendMessage") || !strings.Contains(message, "502") {
		t.Fatalf("expected method and status in error, got %q", message)
	}
	if strings.Contains(message, "123:secret") {
		t.Fatalf("token leaked into error: %q", message)
	}
}

func TestCopyReturnsStatusErrorOnFailure(t *testing.T) {
	client := New(doerFunc(func(*http.Request) (*http.Response, error) {
		return &http.Response{StatusCode: http.StatusNotFound, Body: io.NopCloser(strings.NewReader("missing"))}, nil
	}))
	req, err := http.NewRequest(http.MethodGet, "https://example.test/file/bot123:secret/path/to/file.bin", nil)
	if err != nil {
		t.Fatalf("unexpected request error: %v", err)
	}

	err = client.Copy(context.Background(), req, io.Discard)
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "404") {
		t.Fatalf("expected status in error, got %q", err.Error())
	}
}
