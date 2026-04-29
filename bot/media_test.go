package bot

import (
	"context"
	"encoding/json"
	stderrors "errors"
	"net/http"
	"net/http/httptest"
	"testing"

	apierrors "ai-gram/errors"
	"ai-gram/telegram"
)

func TestFileRefMarshalAndValidation(t *testing.T) {
	tests := []struct {
		name    string
		ref     FileRef
		want    string
		wantErr bool
	}{
		{name: "file id", ref: FileID("abc"), want: `"abc"`},
		{name: "url", ref: FileURL("https://example.com/a.jpg"), want: `"https://example.com/a.jpg"`},
		{name: "empty", ref: FileRef{}, wantErr: true},
		{name: "empty url", ref: FileURL(""), wantErr: true},
		{name: "file url", ref: FileURL("file:///tmp/a.jpg"), wantErr: true},
		{name: "relative url", ref: FileURL("images/a.jpg"), wantErr: true},
		{name: "url without host", ref: FileURL("https:///a.jpg"), wantErr: true},
		{name: "special file id", ref: FileID("AgACAgIAAxkBA+_-:/special"), want: `"AgACAgIAAxkBA+_-:/special"`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, err := json.Marshal(tt.ref)
			if err != nil {
				t.Fatalf("marshal: %v", err)
			}
			if !tt.wantErr && string(body) != tt.want {
				t.Fatalf("marshal = %s, want %s", body, tt.want)
			}

			err = tt.ref.validate("media")
			if tt.wantErr && err == nil {
				t.Fatal("expected error")
			}
			if !tt.wantErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func TestSendPhotoSendsPayloadAndDecodesMessage(t *testing.T) {
	const token = "123:secret"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/bot"+token+"/sendPhoto" {
			t.Fatalf("unexpected path: %q", r.URL.Path)
		}
		var payload map[string]any
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode payload: %v", err)
		}
		if got := payload["chat_id"]; got != float64(12345) {
			t.Fatalf("unexpected chat_id: %#v", got)
		}
		if got := payload["photo"]; got != "photo-file-id" {
			t.Fatalf("unexpected photo: %#v", got)
		}
		if got := payload["caption"]; got != "caption" {
			t.Fatalf("unexpected caption: %#v", got)
		}
		if got := payload["parse_mode"]; got != "HTML" {
			t.Fatalf("unexpected parse_mode: %#v", got)
		}
		if got := payload["disable_notification"]; got != true {
			t.Fatalf("unexpected disable_notification: %#v", got)
		}
		if got := payload["protect_content"]; got != true {
			t.Fatalf("unexpected protect_content: %#v", got)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true,"result":{"message_id":7,"chat":{"id":12345,"type":"private"},"date":100,"caption":"caption","photo":[{"file_id":"photo-file-id","file_unique_id":"photo-unique","width":640,"height":480}]}}`))
	}))
	defer server.Close()

	bot := newTestBot(t, token, server.URL, server.Client())
	message, err := bot.SendPhoto(context.Background(), SendPhotoParams{
		ChatID:              ChatIDInt(12345),
		Photo:               FileID("photo-file-id"),
		Caption:             "caption",
		ParseMode:           "HTML",
		DisableNotification: true,
		ProtectContent:      true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if message == nil || len(message.Photo) != 1 || message.Photo[0].FileID != "photo-file-id" {
		t.Fatalf("unexpected message: %+v", message)
	}
}

func TestSendPhotoAcceptsURL(t *testing.T) {
	const token = "123:secret"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var payload map[string]any
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode payload: %v", err)
		}
		if got := payload["photo"]; got != "https://example.com/photo.jpg" {
			t.Fatalf("unexpected photo: %#v", got)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true,"result":{"message_id":8,"chat":{"id":12345,"type":"private"},"date":100}}`))
	}))
	defer server.Close()

	bot := newTestBot(t, token, server.URL, server.Client())
	message, err := bot.SendPhoto(context.Background(), SendPhotoParams{ChatID: ChatIDInt(12345), Photo: FileURL("https://example.com/photo.jpg")})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if message == nil || message.MessageID != 8 {
		t.Fatalf("unexpected message: %+v", message)
	}
}

func TestSendPhotoValidation(t *testing.T) {
	const token = "123:secret"
	bot := newTestBot(t, token, "https://example.test", nil)
	tests := []struct {
		name   string
		params SendPhotoParams
	}{
		{name: "empty chat", params: SendPhotoParams{Photo: FileID("photo")}},
		{name: "empty photo", params: SendPhotoParams{ChatID: ChatIDInt(12345)}},
		{name: "parse mode and entities", params: SendPhotoParams{ChatID: ChatIDInt(12345), Photo: FileID("photo"), ParseMode: "HTML", CaptionEntities: []telegram.MessageEntity{{Type: telegram.EntityBold, Offset: 0, Length: 4}}}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			message, err := bot.SendPhoto(context.Background(), tt.params)
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

func TestSendPhotoReturnsAPIError(t *testing.T) {
	const token = "123:secret"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":false,"error_code":400,"description":"Bad Request"}`))
	}))
	defer server.Close()

	bot := newTestBot(t, token, server.URL, server.Client())
	message, err := bot.SendPhoto(context.Background(), SendPhotoParams{ChatID: ChatIDInt(12345), Photo: FileID("photo")})
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

func TestSendDocumentSendsPayloadAndDecodesMessage(t *testing.T) {
	const token = "123:secret"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/bot"+token+"/sendDocument" {
			t.Fatalf("unexpected path: %q", r.URL.Path)
		}
		var payload map[string]any
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode payload: %v", err)
		}
		if got := payload["chat_id"]; got != float64(12345) {
			t.Fatalf("unexpected chat_id: %#v", got)
		}
		if got := payload["document"]; got != "document-file-id" {
			t.Fatalf("unexpected document: %#v", got)
		}
		if got := payload["caption"]; got != "document caption" {
			t.Fatalf("unexpected caption: %#v", got)
		}
		if got := payload["parse_mode"]; got != "MarkdownV2" {
			t.Fatalf("unexpected parse_mode: %#v", got)
		}
		if got := payload["disable_notification"]; got != true {
			t.Fatalf("unexpected disable_notification: %#v", got)
		}
		if got := payload["protect_content"]; got != true {
			t.Fatalf("unexpected protect_content: %#v", got)
		}
		if got := payload["disable_content_type_detection"]; got != true {
			t.Fatalf("unexpected disable_content_type_detection: %#v", got)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true,"result":{"message_id":9,"chat":{"id":12345,"type":"private"},"date":101,"caption":"document caption","document":{"file_id":"document-file-id","file_unique_id":"document-unique","file_name":"report.pdf","mime_type":"application/pdf"}}}`))
	}))
	defer server.Close()

	bot := newTestBot(t, token, server.URL, server.Client())
	message, err := bot.SendDocument(context.Background(), SendDocumentParams{
		ChatID:                      ChatIDInt(12345),
		Document:                    FileID("document-file-id"),
		Caption:                     "document caption",
		ParseMode:                   "MarkdownV2",
		DisableNotification:         true,
		ProtectContent:              true,
		DisableContentTypeDetection: true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if message == nil || message.Document == nil || message.Document.FileID != "document-file-id" {
		t.Fatalf("unexpected message: %+v", message)
	}
}

func TestSendDocumentValidation(t *testing.T) {
	const token = "123:secret"
	bot := newTestBot(t, token, "https://example.test", nil)
	tests := []struct {
		name   string
		params SendDocumentParams
	}{
		{name: "empty chat", params: SendDocumentParams{Document: FileID("document")}},
		{name: "empty document", params: SendDocumentParams{ChatID: ChatIDInt(12345)}},
		{name: "parse mode and entities", params: SendDocumentParams{ChatID: ChatIDInt(12345), Document: FileID("document"), ParseMode: "HTML", CaptionEntities: []telegram.MessageEntity{{Type: telegram.EntityBold, Offset: 0, Length: 4}}}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			message, err := bot.SendDocument(context.Background(), tt.params)
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

func TestSendDocumentReturnsAPIError(t *testing.T) {
	const token = "123:secret"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":false,"error_code":400,"description":"Bad Request"}`))
	}))
	defer server.Close()

	bot := newTestBot(t, token, server.URL, server.Client())
	message, err := bot.SendDocument(context.Background(), SendDocumentParams{ChatID: ChatIDInt(12345), Document: FileID("document")})
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

func TestSendMediaContextAndInvalidJSONErrors(t *testing.T) {
	const token = "123:secret"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`not-json`))
	}))
	defer server.Close()

	bot := newTestBot(t, token, server.URL, server.Client())
	message, err := bot.SendPhoto(context.Background(), SendPhotoParams{ChatID: ChatIDInt(12345), Photo: FileID("photo")})
	if err == nil {
		t.Fatal("expected invalid JSON error")
	}
	if message != nil {
		t.Fatalf("expected nil message, got %+v", message)
	}
	assertNoToken(t, err, token)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	message, err = bot.SendDocument(ctx, SendDocumentParams{ChatID: ChatIDInt(12345), Document: FileID("document")})
	if err == nil {
		t.Fatal("expected context error")
	}
	if message != nil {
		t.Fatalf("expected nil message, got %+v", message)
	}
	assertNoToken(t, err, token)
}

func TestFileURLValidationDoesNotLeakToken(t *testing.T) {
	const token = "123:secret"
	bot := newTestBot(t, token, "https://example.test", nil)

	_, err := bot.SendPhoto(context.Background(), SendPhotoParams{ChatID: ChatIDInt(12345), Photo: FileURL("file:///tmp/photo.jpg")})
	if err == nil {
		t.Fatal("expected error")
	}
	assertNoToken(t, err, token)
}
