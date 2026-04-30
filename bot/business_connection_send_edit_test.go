package bot

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/xDilettante/ai-gram/telegram"
)

func TestSendMessageBusinessConnectionID(t *testing.T) {
	const token = "123:secret"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/bot"+token+"/sendMessage" {
			t.Fatalf("unexpected path: %q", r.URL.Path)
		}
		var payload map[string]any
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode payload: %v", err)
		}
		if got := payload["business_connection_id"]; got != "business-1" {
			t.Fatalf("unexpected business_connection_id: %#v", got)
		}
		if got := payload["chat_id"]; got != float64(12345) {
			t.Fatalf("unexpected chat_id: %#v", got)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true,"result":{"message_id":1,"chat":{"id":12345,"type":"private"},"date":1,"text":"hello"}}`))
	}))
	defer server.Close()

	bot := newTestBot(t, token, server.URL, server.Client())
	message, err := bot.SendMessage(context.Background(), SendMessageParams{BusinessConnectionID: "business-1", ChatID: ChatIDInt(12345), Text: "hello"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if message == nil || message.MessageID != 1 {
		t.Fatalf("unexpected message: %+v", message)
	}
}

func TestSendMessageOmitsEmptyBusinessConnectionID(t *testing.T) {
	const token = "123:secret"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var payload map[string]any
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode payload: %v", err)
		}
		if _, ok := payload["business_connection_id"]; ok {
			t.Fatalf("business_connection_id should be omitted: %#v", payload)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true,"result":{"message_id":2,"chat":{"id":12345,"type":"private"},"date":1,"text":"hello"}}`))
	}))
	defer server.Close()

	bot := newTestBot(t, token, server.URL, server.Client())
	message, err := bot.SendMessage(context.Background(), SendMessageParams{ChatID: ChatIDInt(12345), Text: "hello"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if message == nil || message.MessageID != 2 {
		t.Fatalf("unexpected message: %+v", message)
	}
}

func TestSendPhotoMultipartBusinessConnectionID(t *testing.T) {
	const token = "123:secret"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/bot"+token+"/sendPhoto" {
			t.Fatalf("unexpected path: %q", r.URL.Path)
		}
		if got := r.Header.Get("Content-Type"); !strings.HasPrefix(got, "multipart/form-data") {
			t.Fatalf("unexpected content type: %q", got)
		}
		if err := r.ParseMultipartForm(4096); err != nil {
			t.Fatalf("parse multipart: %v", err)
		}
		assertMultipartValue(t, r, "business_connection_id", "business-1")
		assertMultipartValue(t, r, "chat_id", "12345")
		assertMultipartValue(t, r, "photo", "attach://photo")
		content, header := readMultipartFile(t, r, "photo")
		if header.Filename != "photo.jpg" || string(content) != "photo-data" {
			t.Fatalf("unexpected photo file: filename=%q content=%q", header.Filename, content)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true,"result":{"message_id":3,"chat":{"id":12345,"type":"private"},"date":1}}`))
	}))
	defer server.Close()

	bot := newTestBot(t, token, server.URL, server.Client())
	message, err := bot.SendPhoto(context.Background(), SendPhotoParams{
		BusinessConnectionID: "business-1",
		ChatID:               ChatIDInt(12345),
		Photo:                FileUpload(UploadFile{Name: "photo.jpg", Reader: strings.NewReader("photo-data"), ContentType: "image/jpeg"}),
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if message == nil || message.MessageID != 3 {
		t.Fatalf("unexpected message: %+v", message)
	}
}

func TestSendMediaGroupBusinessConnectionIDJSONAndMultipart(t *testing.T) {
	t.Run("json", func(t *testing.T) {
		const token = "123:secret"
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var payload map[string]any
			if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
				t.Fatalf("decode payload: %v", err)
			}
			if got := payload["business_connection_id"]; got != "business-1" {
				t.Fatalf("unexpected business_connection_id: %#v", got)
			}
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"ok":true,"result":[{"message_id":4,"chat":{"id":12345,"type":"private"},"date":1},{"message_id":5,"chat":{"id":12345,"type":"private"},"date":1}]}`))
		}))
		defer server.Close()

		bot := newTestBot(t, token, server.URL, server.Client())
		messages, err := bot.SendMediaGroup(context.Background(), SendMediaGroupParams{
			BusinessConnectionID: "business-1",
			ChatID:               ChatIDInt(12345),
			Media:                []InputMedia{MediaPhoto(FileID("photo-1")), MediaPhoto(FileID("photo-2"))},
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(messages) != 2 {
			t.Fatalf("unexpected messages: %+v", messages)
		}
	})

	t.Run("multipart", func(t *testing.T) {
		const token = "123:secret"
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if err := r.ParseMultipartForm(4096); err != nil {
				t.Fatalf("parse multipart: %v", err)
			}
			assertMultipartValue(t, r, "business_connection_id", "business-1")
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"ok":true,"result":[{"message_id":6,"chat":{"id":12345,"type":"private"},"date":1},{"message_id":7,"chat":{"id":12345,"type":"private"},"date":1}]}`))
		}))
		defer server.Close()

		bot := newTestBot(t, token, server.URL, server.Client())
		messages, err := bot.SendMediaGroup(context.Background(), SendMediaGroupParams{
			BusinessConnectionID: "business-1",
			ChatID:               ChatIDInt(12345),
			Media: []InputMedia{
				MediaPhoto(FileUpload(UploadFile{Name: "photo.jpg", Reader: strings.NewReader("photo-data"), ContentType: "image/jpeg"})),
				MediaDocument(FileID("document-id")),
			},
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(messages) != 2 {
			t.Fatalf("unexpected messages: %+v", messages)
		}
	})
}

func TestEditMessageBusinessConnectionID(t *testing.T) {
	const token = "123:secret"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var payload map[string]any
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode payload: %v", err)
		}
		if got := payload["business_connection_id"]; got != "business-1" {
			t.Fatalf("unexpected business_connection_id: %#v", got)
		}
		if payload["chat_id"] != float64(12345) || payload["message_id"] != float64(77) {
			t.Fatalf("unexpected target payload: %#v", payload)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true,"result":true}`))
	}))
	defer server.Close()

	bot := newTestBot(t, token, server.URL, server.Client())
	result, err := bot.EditMessageText(context.Background(), EditMessageTextParams{BusinessConnectionID: "business-1", Target: EditTargetChat(ChatIDInt(12345), 77), Text: "edited"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result == nil || !result.IsOK() {
		t.Fatalf("unexpected result: %+v", result)
	}
}

func TestEditMessageReplyMarkupBusinessConnectionID(t *testing.T) {
	const token = "123:secret"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var payload map[string]any
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode payload: %v", err)
		}
		if got := payload["business_connection_id"]; got != "business-1" {
			t.Fatalf("unexpected business_connection_id: %#v", got)
		}
		if _, ok := payload["reply_markup"].(map[string]any); !ok {
			t.Fatalf("missing reply_markup: %#v", payload)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true,"result":true}`))
	}))
	defer server.Close()

	markup := telegram.NewInlineKeyboard([]telegram.InlineKeyboardButton{telegram.InlineButtonCallback("OK", "ok")})
	bot := newTestBot(t, token, server.URL, server.Client())
	result, err := bot.EditMessageReplyMarkup(context.Background(), EditMessageReplyMarkupParams{BusinessConnectionID: "business-1", Target: EditTargetChat(ChatIDInt(12345), 77), ReplyMarkup: &markup})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result == nil || !result.IsOK() {
		t.Fatalf("unexpected result: %+v", result)
	}
}

func TestStopPollBusinessConnectionID(t *testing.T) {
	const token = "123:secret"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var payload map[string]any
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode payload: %v", err)
		}
		if got := payload["business_connection_id"]; got != "business-1" {
			t.Fatalf("unexpected business_connection_id: %#v", got)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true,"result":{"id":"poll","question":"q","options":[],"total_voter_count":0,"is_closed":true,"is_anonymous":true,"type":"regular","allows_multiple_answers":false}}`))
	}))
	defer server.Close()

	bot := newTestBot(t, token, server.URL, server.Client())
	poll, err := bot.StopPoll(context.Background(), StopPollParams{BusinessConnectionID: "business-1", ChatID: ChatIDInt(12345), MessageID: 77})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if poll == nil || poll.ID != "poll" {
		t.Fatalf("unexpected poll: %+v", poll)
	}
}
