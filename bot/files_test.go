package bot

import (
	"bytes"
	"context"
	"encoding/json"
	stderrors "errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	apierrors "github.com/xDilettante/ai-gram/errors"
)

func TestGetFileSendsPayloadAndDecodesFile(t *testing.T) {
	const token = "123:secret"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/bot"+token+"/getFile" {
			t.Fatalf("unexpected path: %q", r.URL.Path)
		}

		var payload map[string]string
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode payload: %v", err)
		}
		if got := payload["file_id"]; got != "file-id" {
			t.Fatalf("unexpected file_id: %q", got)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true,"result":{"file_id":"file-id","file_unique_id":"unique-id","file_size":12345678901,"file_path":"photos/file_123.jpg"}}`))
	}))
	defer server.Close()

	bot := newTestBot(t, token, server.URL, server.Client())
	file, err := bot.GetFile(context.Background(), GetFileParams{FileID: "file-id"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if file == nil {
		t.Fatal("expected file")
	}
	if file.FileID != "file-id" || file.FileUniqueID != "unique-id" || file.FileSize != 12345678901 || file.FilePath != "photos/file_123.jpg" {
		t.Fatalf("unexpected file: %+v", file)
	}
}

func TestGetFileRejectsEmptyFileID(t *testing.T) {
	const token = "123:secret"
	bot := newTestBot(t, token, "https://example.test", nil)

	file, err := bot.GetFile(context.Background(), GetFileParams{})
	if err == nil {
		t.Fatal("expected error")
	}
	if file != nil {
		t.Fatalf("expected nil file, got %+v", file)
	}
	if !strings.Contains(err.Error(), "file_id") {
		t.Fatalf("unexpected error: %q", err.Error())
	}
	assertNoToken(t, err, token)
}

func TestGetFileReturnsAPIError(t *testing.T) {
	const token = "123:secret"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":false,"error_code":400,"description":"Bad Request"}`))
	}))
	defer server.Close()

	bot := newTestBot(t, token, server.URL, server.Client())
	file, err := bot.GetFile(context.Background(), GetFileParams{FileID: "file-id"})
	if err == nil {
		t.Fatal("expected error")
	}
	if file != nil {
		t.Fatalf("expected nil file, got %+v", file)
	}

	var apiErr *apierrors.APIError
	if !stderrors.As(err, &apiErr) {
		t.Fatalf("expected APIError, got %T", err)
	}
	if apiErr.Code != 400 {
		t.Fatalf("unexpected APIError code: %d", apiErr.Code)
	}
	assertNoToken(t, err, token)
}

func TestDownloadFileStreamsBody(t *testing.T) {
	const token = "123:secret"
	body := bytes.Repeat([]byte("downloaded-file-body"), 8192)
	var sawRequest bool

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sawRequest = true
		if r.Method != http.MethodGet {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/file/bot"+token+"/photos/file_123.jpg" {
			t.Fatalf("unexpected path: %q", r.URL.Path)
		}
		if strings.Contains(strings.TrimPrefix(r.URL.Path, "/file"), "//") {
			t.Fatalf("download path contains double slash: %q", r.URL.Path)
		}
		_, _ = w.Write(body)
	}))
	defer server.Close()

	bot := newTestFileBot(t, token, server.URL, server.URL+"/file/", server.Client())
	var dst bytes.Buffer
	if err := bot.DownloadFile(context.Background(), "photos/file_123.jpg", &dst); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !sawRequest {
		t.Fatal("expected request")
	}
	if !bytes.Equal(dst.Bytes(), body) {
		t.Fatalf("unexpected body length: %d", dst.Len())
	}
}

func TestDownloadFileUsesCustomFileBaseURL(t *testing.T) {
	const token = "123:secret"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/custom-files/bot"+token+"/documents/report%20final.pdf" {
			t.Fatalf("unexpected path: %q", r.URL.Path)
		}
		_, _ = w.Write([]byte("ok"))
	}))
	defer server.Close()

	bot := newTestFileBot(t, token, "https://api.example.test", server.URL+"/custom-files/", server.Client())
	var dst bytes.Buffer
	if err := bot.DownloadFile(context.Background(), "documents/report final.pdf", &dst); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if dst.String() != "ok" {
		t.Fatalf("unexpected body: %q", dst.String())
	}
}

func TestDownloadFileValidation(t *testing.T) {
	const token = "123:secret"
	bot := newTestFileBot(t, token, "https://example.test", "https://example.test/file", nil)

	tests := []struct {
		name     string
		filePath string
		writer   io.Writer
	}{
		{name: "empty", writer: io.Discard},
		{name: "nil writer", filePath: "photos/file_123.jpg"},
		{name: "http url", filePath: "http://example.com/file.jpg", writer: io.Discard},
		{name: "https url", filePath: "https://example.com/file.jpg", writer: io.Discard},
		{name: "leading slash", filePath: "/photos/file_123.jpg", writer: io.Discard},
		{name: "parent segment", filePath: "photos/../secret.jpg", writer: io.Discard},
		{name: "query", filePath: "photos/file.jpg?x=1", writer: io.Discard},
		{name: "fragment", filePath: "photos/file.jpg#x", writer: io.Discard},
		{name: "empty segment", filePath: "photos//file.jpg", writer: io.Discard},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := bot.DownloadFile(context.Background(), tt.filePath, tt.writer)
			if err == nil {
				t.Fatal("expected error")
			}
			assertNoToken(t, err, token)
		})
	}
}

func TestDownloadFileReturnsStatusErrorWithoutURLOrToken(t *testing.T) {
	const token = "123:secret"
	const fullPath = "/file/bot" + token + "/photos/missing.jpg"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "not found", http.StatusNotFound)
	}))
	defer server.Close()

	bot := newTestFileBot(t, token, server.URL, server.URL+"/file", server.Client())
	var dst bytes.Buffer
	err := bot.DownloadFile(context.Background(), "photos/missing.jpg", &dst)
	if err == nil {
		t.Fatal("expected error")
	}
	assertNoToken(t, err, token)
	if strings.Contains(err.Error(), fullPath) || strings.Contains(err.Error(), server.URL) {
		t.Fatalf("error leaked download URL: %q", err.Error())
	}
}

func TestDownloadFileClosesResponseBodyOnStatusError(t *testing.T) {
	const token = "123:secret"
	body := &fileCloseTrackingBody{Reader: strings.NewReader("not found")}
	client := &http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
		return &http.Response{StatusCode: http.StatusNotFound, Body: body}, nil
	})}
	bot := newTestFileBot(t, token, "https://example.test", "https://example.test/file", client)

	err := bot.DownloadFile(context.Background(), "photos/missing.jpg", io.Discard)
	if err == nil {
		t.Fatal("expected error")
	}
	if !body.closed {
		t.Fatal("expected response body to be closed")
	}
	assertNoToken(t, err, token)
}

func TestDownloadFileReturnsContextError(t *testing.T) {
	const token = "123:secret"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("request should not reach server with canceled context")
	}))
	defer server.Close()

	bot := newTestFileBot(t, token, server.URL, server.URL+"/file", server.Client())
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := bot.DownloadFile(ctx, "photos/file_123.jpg", io.Discard)
	if err == nil {
		t.Fatal("expected error")
	}
	assertNoToken(t, err, token)
}

func TestFileBaseURLDefaultsToBaseURLFileForCustomBaseURL(t *testing.T) {
	const token = "123:secret"
	bot := newTestFileBot(t, token, "https://example.test/api/", "", nil)
	if bot.fileBaseURL != "https://example.test/api/file" {
		t.Fatalf("unexpected fileBaseURL: %q", bot.fileBaseURL)
	}
}

func newTestFileBot(t *testing.T, token string, baseURL string, fileBaseURL string, client *http.Client) *Bot {
	t.Helper()

	bot, err := New(BotConfig{Token: token, BaseURL: baseURL, FileBaseURL: fileBaseURL, HTTPClient: client})
	if err != nil {
		t.Fatalf("unexpected New error: %v", err)
	}

	return bot
}

type fileCloseTrackingBody struct {
	*strings.Reader
	closed bool
}

func (b *fileCloseTrackingBody) Close() error {
	b.closed = true
	return nil
}
