package bot

import (
	"context"
	"encoding/json"
	stderrors "errors"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
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
		{name: "attach file id", ref: FileID("attach://photo"), wantErr: true},
		{name: "attach url", ref: FileURL("attach://photo"), wantErr: true},
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

func TestFileUploadValidation(t *testing.T) {
	valid := FileUpload(UploadFile{Name: "photo.jpg", Reader: strings.NewReader("data"), ContentType: "image/jpeg"})
	if err := valid.validate("photo"); err != nil {
		t.Fatalf("unexpected valid upload error: %v", err)
	}
	if _, err := json.Marshal(valid); err == nil {
		t.Fatal("upload FileRef should not marshal as JSON")
	}

	tests := []struct {
		name string
		file UploadFile
	}{
		{name: "nil reader", file: UploadFile{Name: "photo.jpg"}},
		{name: "empty name", file: UploadFile{Reader: strings.NewReader("data")}},
		{name: "path traversal", file: UploadFile{Name: "../x.jpg", Reader: strings.NewReader("data")}},
		{name: "slash", file: UploadFile{Name: "a/b.jpg", Reader: strings.NewReader("data")}},
		{name: "backslash", file: UploadFile{Name: `a\b.jpg`, Reader: strings.NewReader("data")}},
		{name: "nul", file: UploadFile{Name: "a\x00b.jpg", Reader: strings.NewReader("data")}},
		{name: "invalid content type", file: UploadFile{Name: "a.jpg", Reader: strings.NewReader("data"), ContentType: "bad content type"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := FileUpload(tt.file).validate("photo"); err == nil {
				t.Fatal("expected error")
			}
		})
	}
}

func TestSendPhotoMultipartUpload(t *testing.T) {
	const token = "123:secret"
	photoReader := &chunkReader{data: []byte("photo-data"), chunk: 2}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/bot"+token+"/sendPhoto" {
			t.Fatalf("unexpected path: %q", r.URL.Path)
		}
		contentType := r.Header.Get("Content-Type")
		if strings.HasPrefix(contentType, "application/json") || !strings.HasPrefix(contentType, "multipart/form-data") {
			t.Fatalf("unexpected content type: %q", contentType)
		}
		if err := r.ParseMultipartForm(1024); err != nil {
			t.Fatalf("parse multipart: %v", err)
		}
		assertMultipartValue(t, r, "chat_id", "12345")
		assertMultipartValue(t, r, "photo", "attach://photo")
		assertMultipartValue(t, r, "caption", "caption")
		assertMultipartValue(t, r, "disable_notification", "true")
		assertMultipartValue(t, r, "protect_content", "true")
		var entities []telegram.MessageEntity
		if err := json.Unmarshal([]byte(r.MultipartForm.Value["caption_entities"][0]), &entities); err != nil {
			t.Fatalf("decode caption_entities: %v", err)
		}
		if len(entities) != 1 || entities[0].Type != telegram.EntityBold {
			t.Fatalf("unexpected caption_entities: %#v", entities)
		}
		content, header := readMultipartFile(t, r, "photo")
		if header.Filename != "photo.jpg" {
			t.Fatalf("unexpected filename: %q", header.Filename)
		}
		if got := header.Header.Get("Content-Type"); got != "image/jpeg" {
			t.Fatalf("unexpected part content type: %q", got)
		}
		if string(content) != "photo-data" {
			t.Fatalf("unexpected file content: %q", content)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true,"result":{"message_id":10,"chat":{"id":12345,"type":"private"},"date":100,"photo":[{"file_id":"uploaded-photo","file_unique_id":"photo-unique","width":10,"height":20}]}}`))
	}))
	defer server.Close()

	bot := newTestBot(t, token, server.URL, server.Client())
	message, err := bot.SendPhoto(context.Background(), SendPhotoParams{
		ChatID:              ChatIDInt(12345),
		Photo:               FileUpload(UploadFile{Name: "photo.jpg", Reader: photoReader, ContentType: "image/jpeg"}),
		Caption:             "caption",
		CaptionEntities:     []telegram.MessageEntity{{Type: telegram.EntityBold, Offset: 0, Length: 7}},
		DisableNotification: true,
		ProtectContent:      true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if message == nil || len(message.Photo) != 1 || message.Photo[0].FileID != "uploaded-photo" {
		t.Fatalf("unexpected message: %+v", message)
	}
	if photoReader.reads < 2 {
		t.Fatalf("expected chunked streaming reads, got %d", photoReader.reads)
	}
}

func TestSendDocumentMultipartUpload(t *testing.T) {
	const token = "123:secret"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/bot"+token+"/sendDocument" {
			t.Fatalf("unexpected path: %q", r.URL.Path)
		}
		if err := r.ParseMultipartForm(1024); err != nil {
			t.Fatalf("parse multipart: %v", err)
		}
		assertMultipartValue(t, r, "document", "attach://document")
		assertMultipartValue(t, r, "disable_content_type_detection", "true")
		content, header := readMultipartFile(t, r, "document")
		if header.Filename != "report.pdf" {
			t.Fatalf("unexpected filename: %q", header.Filename)
		}
		if string(content) != "document-data" {
			t.Fatalf("unexpected file content: %q", content)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true,"result":{"message_id":11,"chat":{"id":12345,"type":"private"},"date":101,"document":{"file_id":"uploaded-document","file_unique_id":"document-unique","file_name":"report.pdf"}}}`))
	}))
	defer server.Close()

	bot := newTestBot(t, token, server.URL, server.Client())
	message, err := bot.SendDocument(context.Background(), SendDocumentParams{
		ChatID:                      ChatIDInt(12345),
		Document:                    FileUpload(UploadFile{Name: "report.pdf", Reader: strings.NewReader("document-data")}),
		DisableContentTypeDetection: true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if message == nil || message.Document == nil || message.Document.FileID != "uploaded-document" {
		t.Fatalf("unexpected message: %+v", message)
	}
}

func TestSendDocumentAcceptsURL(t *testing.T) {
	const token = "123:secret"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.Header.Get("Content-Type"); got != "application/json" {
			t.Fatalf("unexpected content type: %q", got)
		}
		var payload map[string]any
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode payload: %v", err)
		}
		if got := payload["document"]; got != "https://example.com/report.pdf" {
			t.Fatalf("unexpected document: %#v", got)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true,"result":{"message_id":12,"chat":{"id":12345,"type":"private"},"date":102}}`))
	}))
	defer server.Close()

	bot := newTestBot(t, token, server.URL, server.Client())
	message, err := bot.SendDocument(context.Background(), SendDocumentParams{ChatID: ChatIDInt(12345), Document: FileURL("https://example.com/report.pdf")})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if message == nil || message.MessageID != 12 {
		t.Fatalf("unexpected message: %+v", message)
	}
}

func TestMultipartUploadErrors(t *testing.T) {
	const token = "123:secret"

	t.Run("api error", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, _ = io.Copy(io.Discard, r.Body)
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"ok":false,"error_code":400,"description":"Bad Request"}`))
		}))
		defer server.Close()
		bot := newTestBot(t, token, server.URL, server.Client())
		_, err := bot.SendPhoto(context.Background(), SendPhotoParams{ChatID: ChatIDInt(12345), Photo: FileUpload(UploadFile{Name: "photo.jpg", Reader: strings.NewReader("data")})})
		if err == nil {
			t.Fatal("expected error")
		}
		var apiErr *apierrors.APIError
		if !stderrors.As(err, &apiErr) {
			t.Fatalf("expected APIError, got %T", err)
		}
		assertNoToken(t, err, token)
	})

	t.Run("invalid json", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, _ = io.Copy(io.Discard, r.Body)
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`not-json`))
		}))
		defer server.Close()
		bot := newTestBot(t, token, server.URL, server.Client())
		_, err := bot.SendPhoto(context.Background(), SendPhotoParams{ChatID: ChatIDInt(12345), Photo: FileUpload(UploadFile{Name: "photo.jpg", Reader: strings.NewReader("data")})})
		if err == nil {
			t.Fatal("expected error")
		}
		assertNoToken(t, err, token)
	})

	t.Run("http status", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, _ = io.Copy(io.Discard, r.Body)
			http.Error(w, "server error", http.StatusInternalServerError)
		}))
		defer server.Close()
		bot := newTestBot(t, token, server.URL, server.Client())
		_, err := bot.SendPhoto(context.Background(), SendPhotoParams{ChatID: ChatIDInt(12345), Photo: FileUpload(UploadFile{Name: "photo.jpg", Reader: strings.NewReader("data")})})
		if err == nil {
			t.Fatal("expected error")
		}
		assertNoToken(t, err, token)
		if strings.Contains(err.Error(), server.URL) {
			t.Fatalf("error leaked URL: %q", err.Error())
		}
	})

	t.Run("cancelled context", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			t.Fatal("request should not reach server")
		}))
		defer server.Close()
		bot := newTestBot(t, token, server.URL, server.Client())
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		_, err := bot.SendPhoto(ctx, SendPhotoParams{ChatID: ChatIDInt(12345), Photo: FileUpload(UploadFile{Name: "photo.jpg", Reader: strings.NewReader("data")})})
		if err == nil {
			t.Fatal("expected error")
		}
		assertNoToken(t, err, token)
	})

	t.Run("reader error", func(t *testing.T) {
		client := &http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			_, err := io.ReadAll(req.Body)
			return nil, err
		})}
		bot := newTestBot(t, token, "https://example.test", client)
		_, err := bot.SendPhoto(context.Background(), SendPhotoParams{ChatID: ChatIDInt(12345), Photo: FileUpload(UploadFile{Name: "photo.jpg", Reader: errorReader{}})})
		if err == nil {
			t.Fatal("expected error")
		}
		assertNoToken(t, err, token)
	})
}

func assertMultipartValue(t *testing.T, r *http.Request, name string, want string) {
	t.Helper()
	values := r.MultipartForm.Value[name]
	if len(values) != 1 || values[0] != want {
		t.Fatalf("multipart field %s = %#v, want %q", name, values, want)
	}
}

func readMultipartFile(t *testing.T, r *http.Request, name string) ([]byte, *multipart.FileHeader) {
	t.Helper()
	files := r.MultipartForm.File[name]
	if len(files) != 1 {
		t.Fatalf("multipart file %s count = %d", name, len(files))
	}
	file, err := files[0].Open()
	if err != nil {
		t.Fatalf("open multipart file: %v", err)
	}
	defer file.Close()
	content, err := io.ReadAll(file)
	if err != nil {
		t.Fatalf("read multipart file: %v", err)
	}
	return content, files[0]
}

type chunkReader struct {
	data  []byte
	chunk int
	reads int
}

func (r *chunkReader) Read(p []byte) (int, error) {
	r.reads++
	if len(r.data) == 0 {
		return 0, io.EOF
	}
	limit := r.chunk
	if limit <= 0 || limit > len(r.data) {
		limit = len(r.data)
	}
	if limit > len(p) {
		limit = len(p)
	}
	copy(p, r.data[:limit])
	r.data = r.data[limit:]
	return limit, nil
}

type errorReader struct{}

func (errorReader) Read([]byte) (int, error) {
	return 0, stderrors.New("reader failed")
}
