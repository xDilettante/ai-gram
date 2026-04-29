package bot

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"ai-gram/telegram"
)

func TestSendMessageReplyMarkupJSON(t *testing.T) {
	tests := []struct {
		name   string
		markup telegram.ReplyMarkup
		check  func(t *testing.T, reply map[string]any)
	}{
		{
			name:   "inline keyboard",
			markup: telegram.NewInlineKeyboard([]telegram.InlineKeyboardButton{telegram.InlineButtonCallback("Confirm", "confirm")}),
			check: func(t *testing.T, reply map[string]any) {
				keyboard := reply["inline_keyboard"].([]any)
				row := keyboard[0].([]any)
				button := row[0].(map[string]any)
				if button["text"] != "Confirm" || button["callback_data"] != "confirm" {
					t.Fatalf("unexpected inline button: %#v", button)
				}
			},
		},
		{
			name: "reply keyboard",
			markup: telegram.ReplyKeyboardMarkup{
				Keyboard:       [][]telegram.KeyboardButton{{telegram.KeyboardButtonText("Yes")}},
				ResizeKeyboard: true,
			},
			check: func(t *testing.T, reply map[string]any) {
				keyboard := reply["keyboard"].([]any)
				row := keyboard[0].([]any)
				button := row[0].(map[string]any)
				if button["text"] != "Yes" || reply["resize_keyboard"] != true {
					t.Fatalf("unexpected reply keyboard: %#v", reply)
				}
			},
		},
		{
			name:   "remove keyboard",
			markup: telegram.RemoveKeyboard(false),
			check: func(t *testing.T, reply map[string]any) {
				if reply["remove_keyboard"] != true {
					t.Fatalf("unexpected remove keyboard: %#v", reply)
				}
			},
		},
		{
			name:   "force reply",
			markup: telegram.NewForceReply(),
			check: func(t *testing.T, reply map[string]any) {
				if reply["force_reply"] != true {
					t.Fatalf("unexpected force reply: %#v", reply)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			const token = "123:secret"
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				var payload map[string]any
				if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
					t.Fatalf("decode payload: %v", err)
				}
				reply, ok := payload["reply_markup"].(map[string]any)
				if !ok {
					t.Fatalf("reply_markup missing or invalid: %#v", payload["reply_markup"])
				}
				tt.check(t, reply)
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(`{"ok":true,"result":{"message_id":30,"chat":{"id":12345,"type":"private"},"date":100,"text":"hello"}}`))
			}))
			defer server.Close()

			bot := newTestBot(t, token, server.URL, server.Client())
			message, err := bot.SendMessage(context.Background(), SendMessageParams{ChatID: ChatIDInt(12345), Text: "hello", ReplyMarkup: tt.markup})
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if message == nil || message.MessageID != 30 {
				t.Fatalf("unexpected message: %+v", message)
			}
		})
	}
}

func TestSendMessageInvalidReplyMarkupSkipsRequest(t *testing.T) {
	const token = "123:secret"
	called := false
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		t.Fatal("request should not be sent")
	}))
	defer server.Close()

	bot := newTestBot(t, token, server.URL, server.Client())
	message, err := bot.SendMessage(context.Background(), SendMessageParams{ChatID: ChatIDInt(12345), Text: "hello", ReplyMarkup: telegram.InlineKeyboardMarkup{}})
	if err == nil {
		t.Fatal("expected error")
	}
	if message != nil {
		t.Fatalf("expected nil message, got %+v", message)
	}
	if called {
		t.Fatal("HTTP request was sent")
	}
	assertNoToken(t, err, token)
}

func TestSendMediaReplyMarkupJSON(t *testing.T) {
	tests := []struct {
		name     string
		method   string
		call     func(*Bot) (*telegram.Message, error)
		checkKey string
	}{
		{
			name:   "photo inline",
			method: "sendPhoto",
			call: func(bot *Bot) (*telegram.Message, error) {
				return bot.SendPhoto(context.Background(), SendPhotoParams{ChatID: ChatIDInt(12345), Photo: FileID("photo"), ReplyMarkup: telegram.NewInlineKeyboard([]telegram.InlineKeyboardButton{telegram.InlineButtonCallback("x", "x")})})
			},
			checkKey: "inline_keyboard",
		},
		{
			name:   "document reply keyboard",
			method: "sendDocument",
			call: func(bot *Bot) (*telegram.Message, error) {
				return bot.SendDocument(context.Background(), SendDocumentParams{ChatID: ChatIDInt(12345), Document: FileURL("https://example.com/doc.pdf"), ReplyMarkup: telegram.NewReplyKeyboard([]telegram.KeyboardButton{telegram.KeyboardButtonText("OK")})})
			},
			checkKey: "keyboard",
		},
		{
			name:   "video inline",
			method: "sendVideo",
			call: func(bot *Bot) (*telegram.Message, error) {
				return bot.SendVideo(context.Background(), SendVideoParams{ChatID: ChatIDInt(12345), Video: FileID("video"), ReplyMarkup: telegram.NewInlineKeyboard([]telegram.InlineKeyboardButton{telegram.InlineButtonCallback("x", "x")})})
			},
			checkKey: "inline_keyboard",
		},
		{
			name:   "audio remove",
			method: "sendAudio",
			call: func(bot *Bot) (*telegram.Message, error) {
				return bot.SendAudio(context.Background(), SendAudioParams{ChatID: ChatIDInt(12345), Audio: FileID("audio"), ReplyMarkup: telegram.RemoveKeyboard(false)})
			},
			checkKey: "remove_keyboard",
		},
		{
			name:   "voice force",
			method: "sendVoice",
			call: func(bot *Bot) (*telegram.Message, error) {
				return bot.SendVoice(context.Background(), SendVoiceParams{ChatID: ChatIDInt(12345), Voice: FileID("voice"), ReplyMarkup: telegram.NewForceReply()})
			},
			checkKey: "force_reply",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			const token = "123:secret"
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path != "/bot"+token+"/"+tt.method {
					t.Fatalf("unexpected path: %q", r.URL.Path)
				}
				if got := r.Header.Get("Content-Type"); got != "application/json" {
					t.Fatalf("unexpected content type: %q", got)
				}
				var payload map[string]any
				if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
					t.Fatalf("decode payload: %v", err)
				}
				reply, ok := payload["reply_markup"].(map[string]any)
				if !ok {
					t.Fatalf("reply_markup missing or invalid: %#v", payload["reply_markup"])
				}
				if _, ok := reply[tt.checkKey]; !ok {
					t.Fatalf("reply_markup missing %s: %#v", tt.checkKey, reply)
				}
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(`{"ok":true,"result":{"message_id":31,"chat":{"id":12345,"type":"private"},"date":100}}`))
			}))
			defer server.Close()

			bot := newTestBot(t, token, server.URL, server.Client())
			message, err := tt.call(bot)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if message == nil || message.MessageID != 31 {
				t.Fatalf("unexpected message: %+v", message)
			}
		})
	}
}

func TestSendMediaReplyMarkupMultipart(t *testing.T) {
	tests := []struct {
		name     string
		method   string
		call     func(*Bot) (*telegram.Message, error)
		checkKey string
	}{
		{
			name:   "photo inline",
			method: "sendPhoto",
			call: func(bot *Bot) (*telegram.Message, error) {
				return bot.SendPhoto(context.Background(), SendPhotoParams{ChatID: ChatIDInt(12345), Photo: FileUpload(UploadFile{Name: "photo.jpg", Reader: strings.NewReader("photo")}), ReplyMarkup: telegram.NewInlineKeyboard([]telegram.InlineKeyboardButton{telegram.InlineButtonCallback("x", "x")})})
			},
			checkKey: "inline_keyboard",
		},
		{
			name:   "document remove",
			method: "sendDocument",
			call: func(bot *Bot) (*telegram.Message, error) {
				return bot.SendDocument(context.Background(), SendDocumentParams{ChatID: ChatIDInt(12345), Document: FileUpload(UploadFile{Name: "doc.txt", Reader: strings.NewReader("doc")}), ReplyMarkup: telegram.RemoveKeyboard(false)})
			},
			checkKey: "remove_keyboard",
		},
		{
			name:   "video inline",
			method: "sendVideo",
			call: func(bot *Bot) (*telegram.Message, error) {
				return bot.SendVideo(context.Background(), SendVideoParams{ChatID: ChatIDInt(12345), Video: FileUpload(UploadFile{Name: "video.mp4", Reader: strings.NewReader("video")}), ReplyMarkup: telegram.NewInlineKeyboard([]telegram.InlineKeyboardButton{telegram.InlineButtonCallback("x", "x")})})
			},
			checkKey: "inline_keyboard",
		},
		{
			name:   "audio force",
			method: "sendAudio",
			call: func(bot *Bot) (*telegram.Message, error) {
				return bot.SendAudio(context.Background(), SendAudioParams{ChatID: ChatIDInt(12345), Audio: FileUpload(UploadFile{Name: "audio.mp3", Reader: strings.NewReader("audio")}), ReplyMarkup: telegram.NewForceReply()})
			},
			checkKey: "force_reply",
		},
		{
			name:   "voice reply keyboard",
			method: "sendVoice",
			call: func(bot *Bot) (*telegram.Message, error) {
				return bot.SendVoice(context.Background(), SendVoiceParams{ChatID: ChatIDInt(12345), Voice: FileUpload(UploadFile{Name: "voice.ogg", Reader: strings.NewReader("voice")}), ReplyMarkup: telegram.NewReplyKeyboard([]telegram.KeyboardButton{telegram.KeyboardButtonText("OK")})})
			},
			checkKey: "keyboard",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			const token = "123:secret"
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path != "/bot"+token+"/"+tt.method {
					t.Fatalf("unexpected path: %q", r.URL.Path)
				}
				if got := r.Header.Get("Content-Type"); !strings.HasPrefix(got, "multipart/form-data") {
					t.Fatalf("unexpected content type: %q", got)
				}
				if err := r.ParseMultipartForm(2048); err != nil {
					t.Fatalf("parse multipart: %v", err)
				}
				values := r.MultipartForm.Value["reply_markup"]
				if len(values) != 1 {
					t.Fatalf("reply_markup field count = %d", len(values))
				}
				var reply map[string]any
				if err := json.Unmarshal([]byte(values[0]), &reply); err != nil {
					t.Fatalf("decode reply_markup: %v", err)
				}
				if _, ok := reply[tt.checkKey]; !ok {
					t.Fatalf("reply_markup missing %s: %#v", tt.checkKey, reply)
				}
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(`{"ok":true,"result":{"message_id":32,"chat":{"id":12345,"type":"private"},"date":100}}`))
			}))
			defer server.Close()

			bot := newTestBot(t, token, server.URL, server.Client())
			message, err := tt.call(bot)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if message == nil || message.MessageID != 32 {
				t.Fatalf("unexpected message: %+v", message)
			}
		})
	}
}

func TestSendMediaInvalidReplyMarkupSkipsRequest(t *testing.T) {
	const token = "123:secret"
	called := false
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		t.Fatal("request should not be sent")
	}))
	defer server.Close()

	bot := newTestBot(t, token, server.URL, server.Client())
	message, err := bot.SendPhoto(context.Background(), SendPhotoParams{ChatID: ChatIDInt(12345), Photo: FileUpload(UploadFile{Name: "photo.jpg", Reader: strings.NewReader("photo")}), ReplyMarkup: telegram.InlineKeyboardMarkup{}})
	if err == nil {
		t.Fatal("expected error")
	}
	if message != nil {
		t.Fatalf("expected nil message, got %+v", message)
	}
	if called {
		t.Fatal("HTTP request was sent")
	}
	assertNoToken(t, err, token)
}
