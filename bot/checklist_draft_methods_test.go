package bot

import (
	"context"
	"encoding/json"
	stderrors "errors"
	"net/http"
	"net/http/httptest"
	"testing"

	apierrors "github.com/xDilettante/ai-gram/errors"
	"github.com/xDilettante/ai-gram/telegram"
)

func TestSendChecklistSendsPayloadAndDecodesMessage(t *testing.T) {
	const token = "123:secret"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/bot"+token+"/sendChecklist" {
			t.Fatalf("unexpected path: %q", r.URL.Path)
		}
		var payload map[string]any
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode payload: %v", err)
		}
		if payload["business_connection_id"] != "bc-1" || payload["chat_id"] != float64(12345) || payload["disable_notification"] != true || payload["protect_content"] != true || payload["message_effect_id"] != "effect-id" {
			t.Fatalf("unexpected payload: %#v", payload)
		}
		checklist := payload["checklist"].(map[string]any)
		if checklist["title"] != "Release" || checklist["parse_mode"] != "HTML" || checklist["others_can_add_tasks"] != true || checklist["others_can_mark_tasks_as_done"] != true {
			t.Fatalf("unexpected checklist: %#v", checklist)
		}
		tasks := checklist["tasks"].([]any)
		task := tasks[0].(map[string]any)
		if task["id"] != float64(1) || task["text"] != "Build" || task["parse_mode"] != "MarkdownV2" {
			t.Fatalf("unexpected task: %#v", task)
		}
		reply := payload["reply_parameters"].(map[string]any)
		if reply["message_id"] != float64(9) {
			t.Fatalf("unexpected reply parameters: %#v", reply)
		}
		markup := payload["reply_markup"].(map[string]any)
		if _, ok := markup["inline_keyboard"]; !ok {
			t.Fatalf("missing inline keyboard: %#v", markup)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true,"result":{"message_id":10,"chat":{"id":12345,"type":"private"},"date":100,"checklist":{"title":"Release","tasks":[{"id":1,"text":"Build"}]}}}`))
	}))
	defer server.Close()

	markup := telegram.NewInlineKeyboard([]telegram.InlineKeyboardButton{telegram.InlineButtonCallback("OK", "ok")})
	bot := newTestBot(t, token, server.URL, server.Client())
	message, err := bot.SendChecklist(context.Background(), SendChecklistParams{
		BusinessConnectionID: "bc-1",
		ChatID:               ChatIDInt(12345),
		Checklist: telegram.InputChecklist{
			Title:                    "Release",
			ParseMode:                "HTML",
			Tasks:                    []telegram.InputChecklistTask{{ID: 1, Text: "Build", ParseMode: "MarkdownV2"}},
			OthersCanAddTasks:        true,
			OthersCanMarkTasksAsDone: true,
		},
		DisableNotification: true,
		ProtectContent:      true,
		MessageEffectID:     "effect-id",
		ReplyParameters:     &telegram.ReplyParameters{MessageID: 9},
		ReplyMarkup:         &markup,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if message == nil || message.Checklist == nil || message.Checklist.Title != "Release" {
		t.Fatalf("unexpected message: %+v", message)
	}
}

func TestEditMessageChecklistSendsPayloadAndDecodesMessage(t *testing.T) {
	const token = "123:secret"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/bot"+token+"/editMessageChecklist" {
			t.Fatalf("unexpected path: %q", r.URL.Path)
		}
		var payload map[string]any
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode payload: %v", err)
		}
		if payload["business_connection_id"] != "bc-1" || payload["chat_id"] != float64(12345) || payload["message_id"] != float64(10) {
			t.Fatalf("unexpected payload: %#v", payload)
		}
		checklist := payload["checklist"].(map[string]any)
		if checklist["title"] != "Release edited" {
			t.Fatalf("unexpected checklist: %#v", checklist)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true,"result":{"message_id":10,"chat":{"id":12345,"type":"private"},"date":100,"checklist":{"title":"Release edited","tasks":[{"id":1,"text":"Build"}]}}}`))
	}))
	defer server.Close()

	bot := newTestBot(t, token, server.URL, server.Client())
	message, err := bot.EditMessageChecklist(context.Background(), EditMessageChecklistParams{
		BusinessConnectionID: "bc-1",
		ChatID:               ChatIDInt(12345),
		MessageID:            10,
		Checklist: telegram.InputChecklist{
			Title: "Release edited",
			Tasks: []telegram.InputChecklistTask{{ID: 1, Text: "Build"}},
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if message == nil || message.Checklist == nil || message.Checklist.Title != "Release edited" {
		t.Fatalf("unexpected message: %+v", message)
	}
}

func TestSendMessageDraftSendsPayloadAndDecodesResult(t *testing.T) {
	const token = "123:secret"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/bot"+token+"/sendMessageDraft" {
			t.Fatalf("unexpected path: %q", r.URL.Path)
		}
		var payload map[string]any
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode payload: %v", err)
		}
		if payload["chat_id"] != float64(12345) || payload["message_thread_id"] != float64(7) || payload["draft_id"] != float64(99) || payload["text"] != "Generating" || payload["parse_mode"] != "HTML" {
			t.Fatalf("unexpected payload: %#v", payload)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true,"result":true}`))
	}))
	defer server.Close()

	bot := newTestBot(t, token, server.URL, server.Client())
	ok, err := bot.SendMessageDraft(context.Background(), SendMessageDraftParams{
		ChatID:          ChatIDInt(12345),
		MessageThreadID: 7,
		DraftID:         99,
		Text:            "Generating",
		ParseMode:       "HTML",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Fatal("expected true")
	}
}

func TestSendMessageDraftAllowsEmptyText(t *testing.T) {
	const token = "123:secret"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/bot"+token+"/sendMessageDraft" {
			t.Fatalf("unexpected path: %q", r.URL.Path)
		}
		var payload map[string]any
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode payload: %v", err)
		}
		if payload["chat_id"] != float64(12345) || payload["draft_id"] != float64(99) {
			t.Fatalf("unexpected payload: %#v", payload)
		}
		text, ok := payload["text"].(string)
		if !ok || text != "" {
			t.Fatalf("expected explicit empty text, got %#v", payload["text"])
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true,"result":true}`))
	}))
	defer server.Close()

	bot := newTestBot(t, token, server.URL, server.Client())
	ok, err := bot.SendMessageDraft(context.Background(), SendMessageDraftParams{ChatID: ChatIDInt(12345), DraftID: 99})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Fatal("expected true")
	}
}

func TestChecklistAndDraftValidation(t *testing.T) {
	const token = "123:secret"
	bot := newTestBot(t, token, "https://example.test", nil)

	validChecklist := telegram.InputChecklist{Title: "Release", Tasks: []telegram.InputChecklistTask{{ID: 1, Text: "Build"}}}
	checklistTests := []struct {
		name string
		call func() (*telegram.Message, error)
	}{
		{name: "send missing business connection", call: func() (*telegram.Message, error) {
			return bot.SendChecklist(context.Background(), SendChecklistParams{ChatID: ChatIDInt(1), Checklist: validChecklist})
		}},
		{name: "send missing chat", call: func() (*telegram.Message, error) {
			return bot.SendChecklist(context.Background(), SendChecklistParams{BusinessConnectionID: "bc", Checklist: validChecklist})
		}},
		{name: "send empty checklist title", call: func() (*telegram.Message, error) {
			return bot.SendChecklist(context.Background(), SendChecklistParams{BusinessConnectionID: "bc", ChatID: ChatIDInt(1), Checklist: telegram.InputChecklist{Tasks: []telegram.InputChecklistTask{{ID: 1, Text: "Build"}}}})
		}},
		{name: "send title parse mode conflict", call: func() (*telegram.Message, error) {
			return bot.SendChecklist(context.Background(), SendChecklistParams{BusinessConnectionID: "bc", ChatID: ChatIDInt(1), Checklist: telegram.InputChecklist{Title: "Release", ParseMode: "HTML", TitleEntities: []telegram.MessageEntity{{Type: telegram.EntityBold, Offset: 0, Length: 1}}, Tasks: []telegram.InputChecklistTask{{ID: 1, Text: "Build"}}}})
		}},
		{name: "send empty tasks", call: func() (*telegram.Message, error) {
			return bot.SendChecklist(context.Background(), SendChecklistParams{BusinessConnectionID: "bc", ChatID: ChatIDInt(1), Checklist: telegram.InputChecklist{Title: "Release"}})
		}},
		{name: "send task invalid id", call: func() (*telegram.Message, error) {
			return bot.SendChecklist(context.Background(), SendChecklistParams{BusinessConnectionID: "bc", ChatID: ChatIDInt(1), Checklist: telegram.InputChecklist{Title: "Release", Tasks: []telegram.InputChecklistTask{{Text: "Build"}}}})
		}},
		{name: "send task duplicate id", call: func() (*telegram.Message, error) {
			return bot.SendChecklist(context.Background(), SendChecklistParams{BusinessConnectionID: "bc", ChatID: ChatIDInt(1), Checklist: telegram.InputChecklist{Title: "Release", Tasks: []telegram.InputChecklistTask{{ID: 1, Text: "Build"}, {ID: 1, Text: "Ship"}}}})
		}},
		{name: "send task empty text", call: func() (*telegram.Message, error) {
			return bot.SendChecklist(context.Background(), SendChecklistParams{BusinessConnectionID: "bc", ChatID: ChatIDInt(1), Checklist: telegram.InputChecklist{Title: "Release", Tasks: []telegram.InputChecklistTask{{ID: 1}}}})
		}},
		{name: "send task parse mode conflict", call: func() (*telegram.Message, error) {
			return bot.SendChecklist(context.Background(), SendChecklistParams{BusinessConnectionID: "bc", ChatID: ChatIDInt(1), Checklist: telegram.InputChecklist{Title: "Release", Tasks: []telegram.InputChecklistTask{{ID: 1, Text: "Build", ParseMode: "HTML", TextEntities: []telegram.MessageEntity{{Type: telegram.EntityBold, Offset: 0, Length: 1}}}}}})
		}},
		{name: "send invalid reply parameters", call: func() (*telegram.Message, error) {
			return bot.SendChecklist(context.Background(), SendChecklistParams{BusinessConnectionID: "bc", ChatID: ChatIDInt(1), Checklist: validChecklist, ReplyParameters: &telegram.ReplyParameters{}})
		}},
		{name: "send invalid reply markup", call: func() (*telegram.Message, error) {
			markup := telegram.InlineKeyboardMarkup{}
			return bot.SendChecklist(context.Background(), SendChecklistParams{BusinessConnectionID: "bc", ChatID: ChatIDInt(1), Checklist: validChecklist, ReplyMarkup: &markup})
		}},
		{name: "edit missing message id", call: func() (*telegram.Message, error) {
			return bot.EditMessageChecklist(context.Background(), EditMessageChecklistParams{BusinessConnectionID: "bc", ChatID: ChatIDInt(1), Checklist: validChecklist})
		}},
	}
	for _, tt := range checklistTests {
		t.Run(tt.name, func(t *testing.T) {
			message, err := tt.call()
			if err == nil {
				t.Fatal("expected error")
			}
			if message != nil {
				t.Fatalf("expected nil message, got %+v", message)
			}
			assertNoToken(t, err, token)
		})
	}

	draftTests := []struct {
		name   string
		params SendMessageDraftParams
	}{
		{name: "missing chat", params: SendMessageDraftParams{DraftID: 1, Text: "Generating"}},
		{name: "negative thread", params: SendMessageDraftParams{ChatID: ChatIDInt(1), MessageThreadID: -1, DraftID: 1, Text: "Generating"}},
		{name: "missing draft id", params: SendMessageDraftParams{ChatID: ChatIDInt(1), Text: "Generating"}},
		{name: "parse mode conflict", params: SendMessageDraftParams{ChatID: ChatIDInt(1), DraftID: 1, Text: "Generating", ParseMode: "HTML", Entities: []telegram.MessageEntity{{Type: telegram.EntityBold, Offset: 0, Length: 1}}}},
	}
	for _, tt := range draftTests {
		t.Run("draft "+tt.name, func(t *testing.T) {
			ok, err := bot.SendMessageDraft(context.Background(), tt.params)
			if err == nil {
				t.Fatal("expected error")
			}
			if ok {
				t.Fatal("expected false")
			}
			assertNoToken(t, err, token)
		})
	}
}

func TestChecklistMethodsErrorCases(t *testing.T) {
	validChecklist := telegram.InputChecklist{Title: "Release", Tasks: []telegram.InputChecklistTask{{ID: 1, Text: "Build"}}}
	testSendMethodErrorCases(t, "sendChecklist", func(bot *Bot) (*telegram.Message, error) {
		return bot.SendChecklist(context.Background(), SendChecklistParams{BusinessConnectionID: "bc", ChatID: ChatIDInt(12345), Checklist: validChecklist})
	}, func(bot *Bot, ctx context.Context) (*telegram.Message, error) {
		return bot.SendChecklist(ctx, SendChecklistParams{BusinessConnectionID: "bc", ChatID: ChatIDInt(12345), Checklist: validChecklist})
	})
	testSendMethodErrorCases(t, "editMessageChecklist", func(bot *Bot) (*telegram.Message, error) {
		return bot.EditMessageChecklist(context.Background(), EditMessageChecklistParams{BusinessConnectionID: "bc", ChatID: ChatIDInt(12345), MessageID: 10, Checklist: validChecklist})
	}, func(bot *Bot, ctx context.Context) (*telegram.Message, error) {
		return bot.EditMessageChecklist(ctx, EditMessageChecklistParams{BusinessConnectionID: "bc", ChatID: ChatIDInt(12345), MessageID: 10, Checklist: validChecklist})
	})
}

func TestSendMessageDraftErrorCases(t *testing.T) {
	const token = "123:secret"
	valid := SendMessageDraftParams{ChatID: ChatIDInt(12345), DraftID: 99, Text: "Generating"}
	tests := []struct {
		name   string
		server func(http.ResponseWriter, *http.Request)
		check  func(error)
	}{
		{name: "api error", server: func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"ok":false,"error_code":400,"description":"Bad Request"}`))
		}, check: func(err error) {
			var apiErr *apierrors.APIError
			if !stderrors.As(err, &apiErr) {
				t.Fatalf("expected APIError, got %T", err)
			}
		}},
		{name: "invalid json", server: func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`not-json`))
		}},
		{name: "http status", server: func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "server error", http.StatusInternalServerError)
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path != "/bot"+token+"/sendMessageDraft" {
					t.Fatalf("unexpected path: %q", r.URL.Path)
				}
				tt.server(w, r)
			}))
			defer server.Close()
			bot := newTestBot(t, token, server.URL, server.Client())
			ok, err := bot.SendMessageDraft(context.Background(), valid)
			if err == nil {
				t.Fatal("expected error")
			}
			if ok {
				t.Fatal("expected false")
			}
			if tt.check != nil {
				tt.check(err)
			}
			assertNoToken(t, err, token)
		})
	}

	t.Run("cancelled context", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			t.Fatal("request should not reach server")
		}))
		defer server.Close()
		bot := newTestBot(t, token, server.URL, server.Client())
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		ok, err := bot.SendMessageDraft(ctx, valid)
		if err == nil {
			t.Fatal("expected error")
		}
		if ok {
			t.Fatal("expected false")
		}
		assertNoToken(t, err, token)
	})
}
