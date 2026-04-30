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

func TestSendMediaGroupSendsJSONAndDecodesMessages(t *testing.T) {
	const token = "123:secret"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/bot"+token+"/sendMediaGroup" {
			t.Fatalf("unexpected path: %q", r.URL.Path)
		}
		if got := r.Header.Get("Content-Type"); got != "application/json" {
			t.Fatalf("unexpected content type: %q", got)
		}
		var payload map[string]any
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode payload: %v", err)
		}
		if got := payload["chat_id"]; got != float64(12345) {
			t.Fatalf("unexpected chat_id: %#v", got)
		}
		media, ok := payload["media"].([]any)
		if !ok || len(media) != 2 {
			t.Fatalf("unexpected media: %#v", payload["media"])
		}
		photo := media[0].(map[string]any)
		if photo["type"] != "photo" || photo["media"] != "photo-file-id" || photo["caption"] != "caption" || photo["has_spoiler"] != true {
			t.Fatalf("unexpected photo payload: %#v", photo)
		}
		video := media[1].(map[string]any)
		if video["type"] != "video" || video["media"] != "https://example.com/video.mp4" || video["width"] != float64(640) || video["height"] != float64(480) || video["duration"] != float64(12) || video["supports_streaming"] != true {
			t.Fatalf("unexpected video payload: %#v", video)
		}
		reply := payload["reply_parameters"].(map[string]any)
		if got := reply["message_id"]; got != float64(9) {
			t.Fatalf("unexpected reply message_id: %#v", got)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true,"result":[{"message_id":5,"chat":{"id":12345,"type":"private"},"date":100,"photo":[{"file_id":"sent-photo","file_unique_id":"photo-u","width":10,"height":10}]},{"message_id":6,"chat":{"id":12345,"type":"private"},"date":101,"video":{"file_id":"sent-video","file_unique_id":"video-u","width":640,"height":480,"duration":12}}]}`))
	}))
	defer server.Close()

	bot := newTestBot(t, token, server.URL, server.Client())
	photo := MediaPhoto(FileID("photo-file-id"))
	photo.Caption = "caption"
	photo.HasSpoiler = true
	video := MediaVideo(FileURL("https://example.com/video.mp4"))
	video.Width = 640
	video.Height = 480
	video.Duration = 12
	video.SupportsStreaming = true
	messages, err := bot.SendMediaGroup(context.Background(), SendMediaGroupParams{
		ChatID:          ChatIDInt(12345),
		Media:           []InputMedia{photo, video},
		ReplyParameters: &telegram.ReplyParameters{MessageID: 9},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(messages) != 2 || len(messages[0].Photo) != 1 || messages[1].Video == nil || messages[1].Video.FileID != "sent-video" {
		t.Fatalf("unexpected messages: %+v", messages)
	}
}

func TestSendMediaGroupMultipartUpload(t *testing.T) {
	const token = "123:secret"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/bot"+token+"/sendMediaGroup" {
			t.Fatalf("unexpected path: %q", r.URL.Path)
		}
		if got := r.Header.Get("Content-Type"); !strings.HasPrefix(got, "multipart/form-data") {
			t.Fatalf("unexpected content type: %q", got)
		}
		if err := r.ParseMultipartForm(4096); err != nil {
			t.Fatalf("parse multipart: %v", err)
		}
		assertMultipartValue(t, r, "chat_id", "12345")
		assertMultipartValue(t, r, "message_thread_id", "7")
		assertMultipartValue(t, r, "reply_parameters", `{"message_id":11}`)
		var media []map[string]any
		if err := json.Unmarshal([]byte(r.MultipartForm.Value["media"][0]), &media); err != nil {
			t.Fatalf("decode media field: %v", err)
		}
		if len(media) != 2 || media[0]["media"] != "attach://media0" || media[1]["media"] != "attach://media1" {
			t.Fatalf("unexpected media field: %#v", media)
		}
		content, header := readMultipartFile(t, r, "media0")
		if header.Filename != "photo.jpg" || string(content) != "photo-data" {
			t.Fatalf("unexpected media0 file: filename=%q content=%q", header.Filename, content)
		}
		content, header = readMultipartFile(t, r, "media1")
		if header.Filename != "document.txt" || string(content) != "document-data" {
			t.Fatalf("unexpected media1 file: filename=%q content=%q", header.Filename, content)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true,"result":[{"message_id":7,"chat":{"id":12345,"type":"private"},"date":100},{"message_id":8,"chat":{"id":12345,"type":"private"},"date":101}]}`))
	}))
	defer server.Close()

	bot := newTestBot(t, token, server.URL, server.Client())
	messages, err := bot.SendMediaGroup(context.Background(), SendMediaGroupParams{
		ChatID:          ChatIDInt(12345),
		MessageThreadID: 7,
		Media: []InputMedia{
			MediaPhoto(FileUpload(UploadFile{Name: "photo.jpg", Reader: strings.NewReader("photo-data"), ContentType: "image/jpeg"})),
			MediaDocument(FileUpload(UploadFile{Name: "document.txt", Reader: strings.NewReader("document-data"), ContentType: "text/plain"})),
		},
		ReplyParameters: &telegram.ReplyParameters{MessageID: 11},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(messages) != 2 || messages[0].MessageID != 7 || messages[1].MessageID != 8 {
		t.Fatalf("unexpected messages: %+v", messages)
	}
}

func TestSendMediaGroupThumbnailUploadUsesSeparateAttachmentName(t *testing.T) {
	const token = "123:secret"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.Header.Get("Content-Type"); !strings.HasPrefix(got, "multipart/form-data") {
			t.Fatalf("unexpected content type: %q", got)
		}
		if err := r.ParseMultipartForm(4096); err != nil {
			t.Fatalf("parse multipart: %v", err)
		}
		var media []map[string]any
		if err := json.Unmarshal([]byte(r.MultipartForm.Value["media"][0]), &media); err != nil {
			t.Fatalf("decode media field: %v", err)
		}
		if media[0]["media"] != "attach://media0" || media[0]["thumbnail"] != "attach://thumb0" {
			t.Fatalf("unexpected first media attachments: %#v", media[0])
		}
		if media[1]["media"] != "attach://media1" || media[1]["thumbnail"] != "attach://thumb1" {
			t.Fatalf("unexpected second media attachments: %#v", media[1])
		}
		_, mediaHeader := readMultipartFile(t, r, "media0")
		_, thumbHeader := readMultipartFile(t, r, "thumb0")
		if mediaHeader.Filename == thumbHeader.Filename {
			t.Fatalf("expected distinct media and thumbnail files, got %q", mediaHeader.Filename)
		}
		readMultipartFile(t, r, "media1")
		readMultipartFile(t, r, "thumb1")

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true,"result":[{"message_id":9,"chat":{"id":12345,"type":"private"},"date":100},{"message_id":10,"chat":{"id":12345,"type":"private"},"date":101}]}`))
	}))
	defer server.Close()

	bot := newTestBot(t, token, server.URL, server.Client())
	video := MediaVideo(FileUpload(UploadFile{Name: "video.mp4", Reader: strings.NewReader("video-data")}))
	video.Thumbnail = FileUpload(UploadFile{Name: "video-thumb.jpg", Reader: strings.NewReader("video-thumb")})
	document := MediaDocument(FileUpload(UploadFile{Name: "doc.pdf", Reader: strings.NewReader("doc-data")}))
	document.Thumbnail = FileUpload(UploadFile{Name: "doc-thumb.jpg", Reader: strings.NewReader("doc-thumb")})
	messages, err := bot.SendMediaGroup(context.Background(), SendMediaGroupParams{
		ChatID: ChatIDInt(12345),
		Media:  []InputMedia{video, document},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(messages) != 2 {
		t.Fatalf("unexpected messages: %+v", messages)
	}
}

func TestSendMediaGroupValidation(t *testing.T) {
	const token = "123:secret"
	valid := SendMediaGroupParams{
		ChatID: ChatIDInt(12345),
		Media:  []InputMedia{MediaPhoto(FileID("photo1")), MediaPhoto(FileID("photo2"))},
	}
	many := make([]InputMedia, 11)
	for i := range many {
		many[i] = MediaPhoto(FileID("photo"))
	}
	tests := []struct {
		name   string
		mutate func(*SendMediaGroupParams)
	}{
		{name: "empty chat", mutate: func(p *SendMediaGroupParams) { p.ChatID = ChatID{} }},
		{name: "too few media", mutate: func(p *SendMediaGroupParams) { p.Media = p.Media[:1] }},
		{name: "too many media", mutate: func(p *SendMediaGroupParams) { p.Media = many }},
		{name: "nil media", mutate: func(p *SendMediaGroupParams) { p.Media[1] = nil }},
		{name: "invalid file ref", mutate: func(p *SendMediaGroupParams) { p.Media[0] = MediaPhoto(FileRef{}) }},
		{name: "negative thread", mutate: func(p *SendMediaGroupParams) { p.MessageThreadID = -1 }},
		{name: "negative duration", mutate: func(p *SendMediaGroupParams) {
			video := MediaVideo(FileID("video"))
			video.Duration = -1
			p.Media[0] = video
		}},
		{name: "negative width", mutate: func(p *SendMediaGroupParams) {
			video := MediaVideo(FileID("video"))
			video.Width = -1
			p.Media[0] = video
		}},
		{name: "negative height", mutate: func(p *SendMediaGroupParams) {
			video := MediaVideo(FileID("video"))
			video.Height = -1
			p.Media[0] = video
		}},
		{name: "parse mode and entities", mutate: func(p *SendMediaGroupParams) {
			photo := MediaPhoto(FileID("photo"))
			photo.ParseMode = "HTML"
			photo.CaptionEntities = []telegram.MessageEntity{{Type: telegram.EntityBold, Offset: 0, Length: 1}}
			p.Media[0] = photo
		}},
		{name: "invalid reply parameters", mutate: func(p *SendMediaGroupParams) { p.ReplyParameters = &telegram.ReplyParameters{} }},
		{name: "invalid thumbnail", mutate: func(p *SendMediaGroupParams) {
			video := MediaVideo(FileID("video"))
			video.Thumbnail = FileUpload(UploadFile{Name: "thumb.jpg"})
			p.Media[0] = video
		}},
	}

	bot := newTestBot(t, token, "https://example.test", nil)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			params := valid
			params.Media = append([]InputMedia(nil), valid.Media...)
			tt.mutate(&params)
			messages, err := bot.SendMediaGroup(context.Background(), params)
			if err == nil {
				t.Fatal("expected error")
			}
			if messages != nil {
				t.Fatalf("expected nil messages, got %+v", messages)
			}
			assertNoToken(t, err, token)
		})
	}
}

func TestSendMediaGroupAPIAndTransportErrors(t *testing.T) {
	const token = "123:secret"
	valid := SendMediaGroupParams{ChatID: ChatIDInt(12345), Media: []InputMedia{MediaPhoto(FileID("photo1")), MediaPhoto(FileID("photo2"))}}
	t.Run("api error", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/bot"+token+"/sendMediaGroup" {
				t.Fatalf("unexpected path: %q", r.URL.Path)
			}
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"ok":false,"error_code":400,"description":"Bad Request"}`))
		}))
		defer server.Close()
		bot := newTestBot(t, token, server.URL, server.Client())
		messages, err := bot.SendMediaGroup(context.Background(), valid)
		if err == nil {
			t.Fatal("expected error")
		}
		if messages != nil {
			t.Fatalf("expected nil messages, got %+v", messages)
		}
		var apiErr *apierrors.APIError
		if !stderrors.As(err, &apiErr) {
			t.Fatalf("expected APIError, got %T", err)
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
		_, err := bot.SendMediaGroup(context.Background(), valid)
		if err == nil {
			t.Fatal("expected error")
		}
		assertNoToken(t, err, token)
	})

	t.Run("http status", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "server error", http.StatusInternalServerError)
		}))
		defer server.Close()
		bot := newTestBot(t, token, server.URL, server.Client())
		_, err := bot.SendMediaGroup(context.Background(), valid)
		if err == nil {
			t.Fatal("expected error")
		}
		assertNoToken(t, err, token)
	})

	t.Run("cancelled context", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			t.Fatal("request should not reach server")
		}))
		defer server.Close()
		bot := newTestBot(t, token, server.URL, server.Client())
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		_, err := bot.SendMediaGroup(ctx, valid)
		if err == nil {
			t.Fatal("expected error")
		}
		assertNoToken(t, err, token)
	})
}
