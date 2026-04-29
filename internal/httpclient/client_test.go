package httpclient

import (
	"context"
	"io"
	"net/http"
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
