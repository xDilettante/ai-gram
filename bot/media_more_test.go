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

func TestSendVideoJSONAndURL(t *testing.T) {
	const token = "123:secret"
	request := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		request++
		if r.URL.Path != "/bot"+token+"/sendVideo" {
			t.Fatalf("unexpected path: %q", r.URL.Path)
		}
		if got := r.Header.Get("Content-Type"); got != "application/json" {
			t.Fatalf("unexpected content type: %q", got)
		}
		var payload map[string]any
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode payload: %v", err)
		}
		switch request {
		case 1:
			assertJSONValue(t, payload, "chat_id", float64(12345))
			assertJSONValue(t, payload, "video", "video-file-id")
			assertJSONValue(t, payload, "thumbnail", "https://example.com/thumb.jpg")
			assertJSONValue(t, payload, "cover", "https://example.com/cover.jpg")
			assertJSONValue(t, payload, "start_timestamp", float64(5))
			assertJSONValue(t, payload, "duration", float64(30))
			assertJSONValue(t, payload, "width", float64(640))
			assertJSONValue(t, payload, "height", float64(480))
			assertJSONValue(t, payload, "caption", "caption")
			assertJSONValue(t, payload, "parse_mode", "HTML")
			assertJSONValue(t, payload, "show_caption_above_media", true)
			assertJSONValue(t, payload, "supports_streaming", true)
			assertJSONValue(t, payload, "has_spoiler", true)
			assertJSONValue(t, payload, "disable_notification", true)
			assertJSONValue(t, payload, "protect_content", true)
		case 2:
			assertJSONValue(t, payload, "video", "https://example.com/video.mp4")
		default:
			t.Fatalf("unexpected request %d", request)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true,"result":{"message_id":20,"chat":{"id":12345,"type":"private"},"date":100,"video":{"file_id":"video-result","file_unique_id":"video-unique","width":640,"height":480,"duration":30}}}`))
	}))
	defer server.Close()

	bot := newTestBot(t, token, server.URL, server.Client())
	message, err := bot.SendVideo(context.Background(), SendVideoParams{
		ChatID:                ChatIDInt(12345),
		Video:                 FileID("video-file-id"),
		Thumbnail:             FileURL("https://example.com/thumb.jpg"),
		Cover:                 FileURL("https://example.com/cover.jpg"),
		StartTimestamp:        5,
		Duration:              30,
		Width:                 640,
		Height:                480,
		Caption:               "caption",
		ParseMode:             "HTML",
		ShowCaptionAboveMedia: true,
		SupportsStreaming:     true,
		HasSpoiler:            true,
		DisableNotification:   true,
		ProtectContent:        true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if message == nil || message.Video == nil || message.Video.FileID != "video-result" {
		t.Fatalf("unexpected message: %+v", message)
	}

	message, err = bot.SendVideo(context.Background(), SendVideoParams{ChatID: ChatIDInt(12345), Video: FileURL("https://example.com/video.mp4")})
	if err != nil {
		t.Fatalf("unexpected URL error: %v", err)
	}
	if message == nil || message.Video == nil {
		t.Fatalf("unexpected URL message: %+v", message)
	}
}

func TestSendVideoMultipartUpload(t *testing.T) {
	const token = "123:secret"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/bot"+token+"/sendVideo" {
			t.Fatalf("unexpected path: %q", r.URL.Path)
		}
		if got := r.Header.Get("Content-Type"); !strings.HasPrefix(got, "multipart/form-data") {
			t.Fatalf("unexpected content type: %q", got)
		}
		if err := r.ParseMultipartForm(2048); err != nil {
			t.Fatalf("parse multipart: %v", err)
		}
		assertMultipartValue(t, r, "chat_id", "12345")
		assertMultipartValue(t, r, "video", "attach://video")
		assertMultipartValue(t, r, "thumbnail", "https://example.com/thumb.jpg")
		assertMultipartValue(t, r, "cover", "attach://cover")
		assertMultipartValue(t, r, "start_timestamp", "5")
		assertMultipartValue(t, r, "duration", "30")
		assertMultipartValue(t, r, "width", "640")
		assertMultipartValue(t, r, "height", "480")
		assertMultipartValue(t, r, "show_caption_above_media", "true")
		assertMultipartValue(t, r, "supports_streaming", "true")
		assertMultipartValue(t, r, "has_spoiler", "true")
		var entities []telegram.MessageEntity
		if err := json.Unmarshal([]byte(r.MultipartForm.Value["caption_entities"][0]), &entities); err != nil {
			t.Fatalf("decode caption_entities: %v", err)
		}
		if len(entities) != 1 || entities[0].Type != telegram.EntityBold {
			t.Fatalf("unexpected caption_entities: %#v", entities)
		}
		content, header := readMultipartFile(t, r, "video")
		if header.Filename != "video.mp4" {
			t.Fatalf("unexpected filename: %q", header.Filename)
		}
		if string(content) != "video-data" {
			t.Fatalf("unexpected file content: %q", content)
		}
		content, header = readMultipartFile(t, r, "cover")
		if header.Filename != "cover.jpg" || string(content) != "cover-data" {
			t.Fatalf("unexpected cover file: %q %q", header.Filename, content)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true,"result":{"message_id":21,"chat":{"id":12345,"type":"private"},"date":100,"video":{"file_id":"uploaded-video","file_unique_id":"video-unique","width":640,"height":480,"duration":30}}}`))
	}))
	defer server.Close()

	bot := newTestBot(t, token, server.URL, server.Client())
	message, err := bot.SendVideo(context.Background(), SendVideoParams{
		ChatID:                ChatIDInt(12345),
		Video:                 FileUpload(UploadFile{Name: "video.mp4", Reader: strings.NewReader("video-data"), ContentType: "video/mp4"}),
		Thumbnail:             FileURL("https://example.com/thumb.jpg"),
		Cover:                 FileUpload(UploadFile{Name: "cover.jpg", Reader: strings.NewReader("cover-data"), ContentType: "image/jpeg"}),
		StartTimestamp:        5,
		Duration:              30,
		Width:                 640,
		Height:                480,
		CaptionEntities:       []telegram.MessageEntity{{Type: telegram.EntityBold, Offset: 0, Length: 5}},
		ShowCaptionAboveMedia: true,
		SupportsStreaming:     true,
		HasSpoiler:            true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if message == nil || message.Video == nil || message.Video.FileID != "uploaded-video" {
		t.Fatalf("unexpected message: %+v", message)
	}
}

func TestSendVideoValidation(t *testing.T) {
	const token = "123:secret"
	bot := newTestBot(t, token, "https://example.test", nil)
	tests := []struct {
		name   string
		params SendVideoParams
	}{
		{name: "empty chat", params: SendVideoParams{Video: FileID("video")}},
		{name: "empty video", params: SendVideoParams{ChatID: ChatIDInt(12345)}},
		{name: "negative duration", params: SendVideoParams{ChatID: ChatIDInt(12345), Video: FileID("video"), Duration: -1}},
		{name: "negative width", params: SendVideoParams{ChatID: ChatIDInt(12345), Video: FileID("video"), Width: -1}},
		{name: "negative height", params: SendVideoParams{ChatID: ChatIDInt(12345), Video: FileID("video"), Height: -1}},
		{name: "negative start timestamp", params: SendVideoParams{ChatID: ChatIDInt(12345), Video: FileID("video"), StartTimestamp: -1}},
		{name: "invalid thumbnail", params: SendVideoParams{ChatID: ChatIDInt(12345), Video: FileID("video"), Thumbnail: FileUpload(UploadFile{Name: "thumb.jpg"})}},
		{name: "invalid cover", params: SendVideoParams{ChatID: ChatIDInt(12345), Video: FileID("video"), Cover: FileUpload(UploadFile{Name: "cover.jpg"})}},
		{name: "parse mode and entities", params: SendVideoParams{ChatID: ChatIDInt(12345), Video: FileID("video"), ParseMode: "HTML", CaptionEntities: []telegram.MessageEntity{{Type: telegram.EntityBold, Offset: 0, Length: 1}}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			message, err := bot.SendVideo(context.Background(), tt.params)
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

func TestSendAudioJSONAndURL(t *testing.T) {
	const token = "123:secret"
	request := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		request++
		if r.URL.Path != "/bot"+token+"/sendAudio" {
			t.Fatalf("unexpected path: %q", r.URL.Path)
		}
		if got := r.Header.Get("Content-Type"); got != "application/json" {
			t.Fatalf("unexpected content type: %q", got)
		}
		var payload map[string]any
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode payload: %v", err)
		}
		switch request {
		case 1:
			assertJSONValue(t, payload, "chat_id", float64(12345))
			assertJSONValue(t, payload, "audio", "audio-file-id")
			assertJSONValue(t, payload, "duration", float64(120))
			assertJSONValue(t, payload, "performer", "Performer")
			assertJSONValue(t, payload, "title", "Title")
			assertJSONValue(t, payload, "caption", "caption")
			assertJSONValue(t, payload, "parse_mode", "MarkdownV2")
			assertJSONValue(t, payload, "disable_notification", true)
			assertJSONValue(t, payload, "protect_content", true)
		case 2:
			assertJSONValue(t, payload, "audio", "https://example.com/audio.mp3")
		default:
			t.Fatalf("unexpected request %d", request)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true,"result":{"message_id":22,"chat":{"id":12345,"type":"private"},"date":100,"audio":{"file_id":"audio-result","file_unique_id":"audio-unique","duration":120,"performer":"Performer","title":"Title"}}}`))
	}))
	defer server.Close()

	bot := newTestBot(t, token, server.URL, server.Client())
	message, err := bot.SendAudio(context.Background(), SendAudioParams{
		ChatID:              ChatIDInt(12345),
		Audio:               FileID("audio-file-id"),
		Duration:            120,
		Performer:           "Performer",
		Title:               "Title",
		Caption:             "caption",
		ParseMode:           "MarkdownV2",
		DisableNotification: true,
		ProtectContent:      true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if message == nil || message.Audio == nil || message.Audio.FileID != "audio-result" {
		t.Fatalf("unexpected message: %+v", message)
	}

	message, err = bot.SendAudio(context.Background(), SendAudioParams{ChatID: ChatIDInt(12345), Audio: FileURL("https://example.com/audio.mp3")})
	if err != nil {
		t.Fatalf("unexpected URL error: %v", err)
	}
	if message == nil || message.Audio == nil {
		t.Fatalf("unexpected URL message: %+v", message)
	}
}

func TestSendAudioMultipartUpload(t *testing.T) {
	const token = "123:secret"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/bot"+token+"/sendAudio" {
			t.Fatalf("unexpected path: %q", r.URL.Path)
		}
		if err := r.ParseMultipartForm(2048); err != nil {
			t.Fatalf("parse multipart: %v", err)
		}
		assertMultipartValue(t, r, "audio", "attach://audio")
		assertMultipartValue(t, r, "duration", "120")
		assertMultipartValue(t, r, "performer", "Performer")
		assertMultipartValue(t, r, "title", "Title")
		var entities []telegram.MessageEntity
		if err := json.Unmarshal([]byte(r.MultipartForm.Value["caption_entities"][0]), &entities); err != nil {
			t.Fatalf("decode caption_entities: %v", err)
		}
		content, header := readMultipartFile(t, r, "audio")
		if header.Filename != "audio.mp3" {
			t.Fatalf("unexpected filename: %q", header.Filename)
		}
		if string(content) != "audio-data" {
			t.Fatalf("unexpected file content: %q", content)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true,"result":{"message_id":23,"chat":{"id":12345,"type":"private"},"date":100,"audio":{"file_id":"uploaded-audio","file_unique_id":"audio-unique","duration":120}}}`))
	}))
	defer server.Close()

	bot := newTestBot(t, token, server.URL, server.Client())
	message, err := bot.SendAudio(context.Background(), SendAudioParams{
		ChatID:          ChatIDInt(12345),
		Audio:           FileUpload(UploadFile{Name: "audio.mp3", Reader: strings.NewReader("audio-data")}),
		Duration:        120,
		Performer:       "Performer",
		Title:           "Title",
		CaptionEntities: []telegram.MessageEntity{{Type: telegram.EntityItalic, Offset: 0, Length: 5}},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if message == nil || message.Audio == nil || message.Audio.FileID != "uploaded-audio" {
		t.Fatalf("unexpected message: %+v", message)
	}
}

func TestSendAudioValidation(t *testing.T) {
	const token = "123:secret"
	bot := newTestBot(t, token, "https://example.test", nil)
	tests := []struct {
		name   string
		params SendAudioParams
	}{
		{name: "empty chat", params: SendAudioParams{Audio: FileID("audio")}},
		{name: "empty audio", params: SendAudioParams{ChatID: ChatIDInt(12345)}},
		{name: "negative duration", params: SendAudioParams{ChatID: ChatIDInt(12345), Audio: FileID("audio"), Duration: -1}},
		{name: "parse mode and entities", params: SendAudioParams{ChatID: ChatIDInt(12345), Audio: FileID("audio"), ParseMode: "HTML", CaptionEntities: []telegram.MessageEntity{{Type: telegram.EntityBold, Offset: 0, Length: 1}}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			message, err := bot.SendAudio(context.Background(), tt.params)
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

func TestSendVoiceJSONAndURL(t *testing.T) {
	const token = "123:secret"
	request := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		request++
		if r.URL.Path != "/bot"+token+"/sendVoice" {
			t.Fatalf("unexpected path: %q", r.URL.Path)
		}
		if got := r.Header.Get("Content-Type"); got != "application/json" {
			t.Fatalf("unexpected content type: %q", got)
		}
		var payload map[string]any
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode payload: %v", err)
		}
		switch request {
		case 1:
			assertJSONValue(t, payload, "chat_id", float64(12345))
			assertJSONValue(t, payload, "voice", "voice-file-id")
			assertJSONValue(t, payload, "duration", float64(20))
			assertJSONValue(t, payload, "caption", "caption")
			assertJSONValue(t, payload, "parse_mode", "HTML")
			assertJSONValue(t, payload, "disable_notification", true)
			assertJSONValue(t, payload, "protect_content", true)
		case 2:
			assertJSONValue(t, payload, "voice", "https://example.com/voice.ogg")
		default:
			t.Fatalf("unexpected request %d", request)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true,"result":{"message_id":24,"chat":{"id":12345,"type":"private"},"date":100,"voice":{"file_id":"voice-result","file_unique_id":"voice-unique","duration":20}}}`))
	}))
	defer server.Close()

	bot := newTestBot(t, token, server.URL, server.Client())
	message, err := bot.SendVoice(context.Background(), SendVoiceParams{
		ChatID:              ChatIDInt(12345),
		Voice:               FileID("voice-file-id"),
		Duration:            20,
		Caption:             "caption",
		ParseMode:           "HTML",
		DisableNotification: true,
		ProtectContent:      true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if message == nil || message.Voice == nil || message.Voice.FileID != "voice-result" {
		t.Fatalf("unexpected message: %+v", message)
	}

	message, err = bot.SendVoice(context.Background(), SendVoiceParams{ChatID: ChatIDInt(12345), Voice: FileURL("https://example.com/voice.ogg")})
	if err != nil {
		t.Fatalf("unexpected URL error: %v", err)
	}
	if message == nil || message.Voice == nil {
		t.Fatalf("unexpected URL message: %+v", message)
	}
}

func TestSendVoiceMultipartUpload(t *testing.T) {
	const token = "123:secret"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/bot"+token+"/sendVoice" {
			t.Fatalf("unexpected path: %q", r.URL.Path)
		}
		if err := r.ParseMultipartForm(2048); err != nil {
			t.Fatalf("parse multipart: %v", err)
		}
		assertMultipartValue(t, r, "voice", "attach://voice")
		assertMultipartValue(t, r, "duration", "20")
		var entities []telegram.MessageEntity
		if err := json.Unmarshal([]byte(r.MultipartForm.Value["caption_entities"][0]), &entities); err != nil {
			t.Fatalf("decode caption_entities: %v", err)
		}
		content, header := readMultipartFile(t, r, "voice")
		if header.Filename != "voice.ogg" {
			t.Fatalf("unexpected filename: %q", header.Filename)
		}
		if string(content) != "voice-data" {
			t.Fatalf("unexpected file content: %q", content)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true,"result":{"message_id":25,"chat":{"id":12345,"type":"private"},"date":100,"voice":{"file_id":"uploaded-voice","file_unique_id":"voice-unique","duration":20}}}`))
	}))
	defer server.Close()

	bot := newTestBot(t, token, server.URL, server.Client())
	message, err := bot.SendVoice(context.Background(), SendVoiceParams{
		ChatID:          ChatIDInt(12345),
		Voice:           FileUpload(UploadFile{Name: "voice.ogg", Reader: strings.NewReader("voice-data")}),
		Duration:        20,
		CaptionEntities: []telegram.MessageEntity{{Type: telegram.EntityCode, Offset: 0, Length: 5}},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if message == nil || message.Voice == nil || message.Voice.FileID != "uploaded-voice" {
		t.Fatalf("unexpected message: %+v", message)
	}
}

func TestSendVoiceValidation(t *testing.T) {
	const token = "123:secret"
	bot := newTestBot(t, token, "https://example.test", nil)
	tests := []struct {
		name   string
		params SendVoiceParams
	}{
		{name: "empty chat", params: SendVoiceParams{Voice: FileID("voice")}},
		{name: "empty voice", params: SendVoiceParams{ChatID: ChatIDInt(12345)}},
		{name: "negative duration", params: SendVoiceParams{ChatID: ChatIDInt(12345), Voice: FileID("voice"), Duration: -1}},
		{name: "parse mode and entities", params: SendVoiceParams{ChatID: ChatIDInt(12345), Voice: FileID("voice"), ParseMode: "HTML", CaptionEntities: []telegram.MessageEntity{{Type: telegram.EntityBold, Offset: 0, Length: 1}}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			message, err := bot.SendVoice(context.Background(), tt.params)
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

func TestSendAdditionalMediaAPIErrors(t *testing.T) {
	const token = "123:secret"
	tests := []struct {
		name string
		call func(*Bot) (*telegram.Message, error)
	}{
		{name: "video", call: func(bot *Bot) (*telegram.Message, error) {
			return bot.SendVideo(context.Background(), SendVideoParams{ChatID: ChatIDInt(12345), Video: FileID("video")})
		}},
		{name: "audio", call: func(bot *Bot) (*telegram.Message, error) {
			return bot.SendAudio(context.Background(), SendAudioParams{ChatID: ChatIDInt(12345), Audio: FileID("audio")})
		}},
		{name: "voice", call: func(bot *Bot) (*telegram.Message, error) {
			return bot.SendVoice(context.Background(), SendVoiceParams{ChatID: ChatIDInt(12345), Voice: FileID("voice")})
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(`{"ok":false,"error_code":400,"description":"Bad Request"}`))
			}))
			defer server.Close()
			bot := newTestBot(t, token, server.URL, server.Client())
			message, err := tt.call(bot)
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
		})
	}
}

func TestSendAdditionalMediaResponseAndContextErrors(t *testing.T) {
	const token = "123:secret"

	t.Run("invalid json", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`not-json`))
		}))
		defer server.Close()
		bot := newTestBot(t, token, server.URL, server.Client())
		_, err := bot.SendVideo(context.Background(), SendVideoParams{ChatID: ChatIDInt(12345), Video: FileID("video")})
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
		_, err := bot.SendAudio(context.Background(), SendAudioParams{ChatID: ChatIDInt(12345), Audio: FileID("audio")})
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
		_, err := bot.SendVoice(ctx, SendVoiceParams{ChatID: ChatIDInt(12345), Voice: FileID("voice")})
		if err == nil {
			t.Fatal("expected error")
		}
		assertNoToken(t, err, token)
	})
}

func assertJSONValue(t *testing.T, payload map[string]any, name string, want any) {
	t.Helper()
	if got := payload[name]; got != want {
		t.Fatalf("payload[%s] = %#v, want %#v", name, got, want)
	}
}
