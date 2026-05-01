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

func TestSendMessageSendsMessageThreadID(t *testing.T) {
	const token = "123:secret"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var payload map[string]any
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode payload: %v", err)
		}
		if got := payload["message_thread_id"]; got != float64(777) {
			t.Fatalf("unexpected message_thread_id: %#v", got)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true,"result":{"message_id":1,"message_thread_id":777,"chat":{"id":12345,"type":"supergroup"},"date":100,"text":"hello"}}`))
	}))
	defer server.Close()

	bot := newTestBot(t, token, server.URL, server.Client())
	message, err := bot.SendMessage(context.Background(), SendMessageParams{ChatID: ChatIDInt(12345), MessageThreadID: 777, Text: "hello"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if message.MessageThreadID != 777 {
		t.Fatalf("unexpected decoded message_thread_id: %d", message.MessageThreadID)
	}
}

func TestSendMessageSendsReplyParameters(t *testing.T) {
	const token = "123:secret"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var payload map[string]any
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode payload: %v", err)
		}
		reply, ok := payload["reply_parameters"].(map[string]any)
		if !ok {
			t.Fatalf("reply_parameters missing: %#v", payload["reply_parameters"])
		}
		if got := reply["message_id"]; got != float64(42) {
			t.Fatalf("unexpected reply message_id: %#v", got)
		}
		if got := reply["allow_sending_without_reply"]; got != true {
			t.Fatalf("unexpected allow_sending_without_reply: %#v", got)
		}
		if got := reply["poll_option_id"]; got != "option-a" {
			t.Fatalf("unexpected poll_option_id: %#v", got)
		}
		if got := reply["chat_id"]; got != float64(777) {
			t.Fatalf("unexpected reply chat_id: %#v", got)
		}
		if got := reply["quote"]; got != "quoted" {
			t.Fatalf("unexpected quote: %#v", got)
		}
		if got := reply["quote_parse_mode"]; got != "HTML" {
			t.Fatalf("unexpected quote_parse_mode: %#v", got)
		}
		if got := reply["quote_position"]; got != float64(3) {
			t.Fatalf("unexpected quote_position: %#v", got)
		}
		if got := reply["checklist_task_id"]; got != float64(9) {
			t.Fatalf("unexpected checklist_task_id: %#v", got)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true,"result":{"message_id":2,"chat":{"id":12345,"type":"private"},"date":100,"text":"hello"}}`))
	}))
	defer server.Close()

	bot := newTestBot(t, token, server.URL, server.Client())
	_, err := bot.SendMessage(context.Background(), SendMessageParams{
		ChatID:          ChatIDInt(12345),
		Text:            "hello",
		ReplyParameters: &telegram.ReplyParameters{MessageID: 42, ChatID: telegram.ReplyChatIDInt(777), AllowSendingWithoutReply: true, Quote: "quoted", QuoteParseMode: "HTML", QuotePosition: 3, ChecklistTaskID: 9, PollOptionID: "option-a"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestReplyParametersMarshalUsernameChatAndQuoteEntities(t *testing.T) {
	reply := telegram.ReplyParameters{
		MessageID:     42,
		ChatID:        telegram.ReplyChatIDUsername("@channel"),
		Quote:         "quoted",
		QuoteEntities: []telegram.MessageEntity{{Type: telegram.EntityBold, Offset: 0, Length: 6}},
	}

	data, err := json.Marshal(reply)
	if err != nil {
		t.Fatalf("marshal reply parameters: %v", err)
	}

	var payload map[string]any
	if err := json.Unmarshal(data, &payload); err != nil {
		t.Fatalf("decode reply parameters: %v", err)
	}
	if payload["chat_id"] != "@channel" {
		t.Fatalf("unexpected chat_id: %#v", payload["chat_id"])
	}
	if entities, ok := payload["quote_entities"].([]any); !ok || len(entities) != 1 {
		t.Fatalf("unexpected quote_entities: %#v", payload["quote_entities"])
	}
}

func TestSendMediaReplyParametersJSON(t *testing.T) {
	tests := []struct {
		name   string
		method string
		call   func(*Bot) (*telegram.Message, error)
		check  func(t *testing.T, payload map[string]any)
	}{
		{
			name:   "photo reply parameters",
			method: "sendPhoto",
			call: func(bot *Bot) (*telegram.Message, error) {
				return bot.SendPhoto(context.Background(), SendPhotoParams{ChatID: ChatIDInt(12345), Photo: FileID("photo"), ReplyParameters: &telegram.ReplyParameters{MessageID: 55}})
			},
			check: func(t *testing.T, payload map[string]any) {
				reply := payload["reply_parameters"].(map[string]any)
				if got := reply["message_id"]; got != float64(55) {
					t.Fatalf("unexpected reply message_id: %#v", got)
				}
			},
		},
		{
			name:   "document thread id",
			method: "sendDocument",
			call: func(bot *Bot) (*telegram.Message, error) {
				return bot.SendDocument(context.Background(), SendDocumentParams{ChatID: ChatIDInt(12345), MessageThreadID: 999, Document: FileURL("https://example.com/doc.pdf")})
			},
			check: func(t *testing.T, payload map[string]any) {
				if got := payload["message_thread_id"]; got != float64(999) {
					t.Fatalf("unexpected message_thread_id: %#v", got)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			const token = "123:secret"
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path != "/bot"+token+"/"+tt.method {
					t.Fatalf("unexpected path: %q", r.URL.Path)
				}
				var payload map[string]any
				if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
					t.Fatalf("decode payload: %v", err)
				}
				tt.check(t, payload)
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(`{"ok":true,"result":{"message_id":3,"chat":{"id":12345,"type":"private"},"date":100}}`))
			}))
			defer server.Close()

			bot := newTestBot(t, token, server.URL, server.Client())
			if _, err := tt.call(bot); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func TestSendMediaReplyParametersMultipart(t *testing.T) {
	tests := []struct {
		name   string
		method string
		call   func(*Bot) (*telegram.Message, error)
		check  func(t *testing.T, r *http.Request)
	}{
		{
			name:   "photo thread id",
			method: "sendPhoto",
			call: func(bot *Bot) (*telegram.Message, error) {
				return bot.SendPhoto(context.Background(), SendPhotoParams{ChatID: ChatIDInt(12345), MessageThreadID: 123, Photo: FileUpload(UploadFile{Name: "photo.jpg", Reader: strings.NewReader("photo")})})
			},
			check: func(t *testing.T, r *http.Request) {
				assertMultipartValue(t, r, "message_thread_id", "123")
				readMultipartFile(t, r, "photo")
			},
		},
		{
			name:   "document reply parameters",
			method: "sendDocument",
			call: func(bot *Bot) (*telegram.Message, error) {
				return bot.SendDocument(context.Background(), SendDocumentParams{ChatID: ChatIDInt(12345), Document: FileUpload(UploadFile{Name: "doc.txt", Reader: strings.NewReader("doc")}), ReplyParameters: &telegram.ReplyParameters{MessageID: 77, AllowSendingWithoutReply: true}})
			},
			check: func(t *testing.T, r *http.Request) {
				values := r.MultipartForm.Value["reply_parameters"]
				if len(values) != 1 {
					t.Fatalf("reply_parameters field count = %d", len(values))
				}
				var reply map[string]any
				if err := json.Unmarshal([]byte(values[0]), &reply); err != nil {
					t.Fatalf("decode reply_parameters: %v", err)
				}
				if got := reply["message_id"]; got != float64(77) {
					t.Fatalf("unexpected reply message_id: %#v", got)
				}
				if got := reply["allow_sending_without_reply"]; got != true {
					t.Fatalf("unexpected allow_sending_without_reply: %#v", got)
				}
				readMultipartFile(t, r, "document")
			},
		},
		{
			name:   "video regression",
			method: "sendVideo",
			call: func(bot *Bot) (*telegram.Message, error) {
				return bot.SendVideo(context.Background(), SendVideoParams{ChatID: ChatIDInt(12345), Video: FileUpload(UploadFile{Name: "video.mp4", Reader: strings.NewReader("video")}), ReplyParameters: &telegram.ReplyParameters{MessageID: 10}})
			},
			check: func(t *testing.T, r *http.Request) { readMultipartFile(t, r, "video") },
		},
		{
			name:   "audio regression",
			method: "sendAudio",
			call: func(bot *Bot) (*telegram.Message, error) {
				return bot.SendAudio(context.Background(), SendAudioParams{ChatID: ChatIDInt(12345), Audio: FileUpload(UploadFile{Name: "audio.mp3", Reader: strings.NewReader("audio")}), MessageThreadID: 11})
			},
			check: func(t *testing.T, r *http.Request) { readMultipartFile(t, r, "audio") },
		},
		{
			name:   "voice regression",
			method: "sendVoice",
			call: func(bot *Bot) (*telegram.Message, error) {
				return bot.SendVoice(context.Background(), SendVoiceParams{ChatID: ChatIDInt(12345), Voice: FileUpload(UploadFile{Name: "voice.ogg", Reader: strings.NewReader("voice")}), ReplyParameters: &telegram.ReplyParameters{MessageID: 12}})
			},
			check: func(t *testing.T, r *http.Request) { readMultipartFile(t, r, "voice") },
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
				tt.check(t, r)
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(`{"ok":true,"result":{"message_id":4,"chat":{"id":12345,"type":"private"},"date":100}}`))
			}))
			defer server.Close()

			bot := newTestBot(t, token, server.URL, server.Client())
			if _, err := tt.call(bot); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func TestSendReplyParameterValidation(t *testing.T) {
	const token = "123:secret"
	tests := []struct {
		name string
		call func(*Bot) (*telegram.Message, error)
	}{
		{
			name: "send message negative thread",
			call: func(bot *Bot) (*telegram.Message, error) {
				return bot.SendMessage(context.Background(), SendMessageParams{ChatID: ChatIDInt(12345), MessageThreadID: -1, Text: "hello"})
			},
		},
		{
			name: "send message zero reply message",
			call: func(bot *Bot) (*telegram.Message, error) {
				return bot.SendMessage(context.Background(), SendMessageParams{ChatID: ChatIDInt(12345), Text: "hello", ReplyParameters: &telegram.ReplyParameters{}})
			},
		},
		{
			name: "send photo negative reply message",
			call: func(bot *Bot) (*telegram.Message, error) {
				return bot.SendPhoto(context.Background(), SendPhotoParams{ChatID: ChatIDInt(12345), Photo: FileID("photo"), ReplyParameters: &telegram.ReplyParameters{MessageID: -1}})
			},
		},
		{
			name: "send message invalid reply chat id",
			call: func(bot *Bot) (*telegram.Message, error) {
				return bot.SendMessage(context.Background(), SendMessageParams{ChatID: ChatIDInt(12345), Text: "hello", ReplyParameters: &telegram.ReplyParameters{MessageID: 1, ChatID: telegram.ReplyChatIDInt(0)}})
			},
		},
		{
			name: "send message quote parse conflict",
			call: func(bot *Bot) (*telegram.Message, error) {
				return bot.SendMessage(context.Background(), SendMessageParams{ChatID: ChatIDInt(12345), Text: "hello", ReplyParameters: &telegram.ReplyParameters{MessageID: 1, QuoteParseMode: "HTML", QuoteEntities: []telegram.MessageEntity{{Type: telegram.EntityBold, Offset: 0, Length: 1}}}})
			},
		},
		{
			name: "send message negative quote position",
			call: func(bot *Bot) (*telegram.Message, error) {
				return bot.SendMessage(context.Background(), SendMessageParams{ChatID: ChatIDInt(12345), Text: "hello", ReplyParameters: &telegram.ReplyParameters{MessageID: 1, QuotePosition: -1}})
			},
		},
		{
			name: "send message negative checklist task",
			call: func(bot *Bot) (*telegram.Message, error) {
				return bot.SendMessage(context.Background(), SendMessageParams{ChatID: ChatIDInt(12345), Text: "hello", ReplyParameters: &telegram.ReplyParameters{MessageID: 1, ChecklistTaskID: -1}})
			},
		},
		{
			name: "send document negative thread",
			call: func(bot *Bot) (*telegram.Message, error) {
				return bot.SendDocument(context.Background(), SendDocumentParams{ChatID: ChatIDInt(12345), MessageThreadID: -1, Document: FileID("document")})
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			called := false
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				called = true
				t.Fatal("request should not be sent")
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
			if called {
				t.Fatal("HTTP request was sent")
			}
			assertNoToken(t, err, token)
		})
	}
}
