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

func TestForwardMessagesSendsPayloadAndDecodesMessageIDs(t *testing.T) {
	testBatchMessageIDsSuccess(t, "forwardMessages", func(bot *Bot) ([]int64, error) {
		messageIDs, err := bot.ForwardMessages(context.Background(), ForwardMessagesParams{
			ChatID:              ChatIDInt(12345),
			FromChatID:          ChatIDInt(67890),
			MessageThreadID:     11,
			MessageIDs:          []int64{77, 78},
			DisableNotification: true,
			ProtectContent:      true,
		})
		return collectMessageIDs(messageIDs), err
	}, func(t *testing.T, payload map[string]any) {
		assertBatchBasePayload(t, payload)
		if payload["disable_notification"] != true || payload["protect_content"] != true {
			t.Fatalf("unexpected option payload: %#v", payload)
		}
	})
}

func TestCopyMessagesSendsPayloadAndDecodesMessageIDs(t *testing.T) {
	testBatchMessageIDsSuccess(t, "copyMessages", func(bot *Bot) ([]int64, error) {
		messageIDs, err := bot.CopyMessages(context.Background(), CopyMessagesParams{
			ChatID:              ChatIDInt(12345),
			FromChatID:          ChatIDInt(67890),
			MessageThreadID:     11,
			MessageIDs:          []int64{77, 78},
			DisableNotification: true,
			ProtectContent:      true,
			RemoveCaption:       true,
		})
		return collectMessageIDs(messageIDs), err
	}, func(t *testing.T, payload map[string]any) {
		assertBatchBasePayload(t, payload)
		if payload["disable_notification"] != true || payload["protect_content"] != true || payload["remove_caption"] != true {
			t.Fatalf("unexpected option payload: %#v", payload)
		}
	})
}

func TestDeleteMessagesSendsPayloadAndDecodesResult(t *testing.T) {
	const token = "123:secret"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/bot"+token+"/deleteMessages" {
			t.Fatalf("unexpected path: %q", r.URL.Path)
		}
		var payload map[string]any
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode payload: %v", err)
		}
		if payload["chat_id"] != float64(12345) {
			t.Fatalf("unexpected chat payload: %#v", payload)
		}
		assertMessageIDsPayload(t, payload["message_ids"])
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true,"result":true}`))
	}))
	defer server.Close()

	bot := newTestBot(t, token, server.URL, server.Client())
	ok, err := bot.DeleteMessages(context.Background(), DeleteMessagesParams{ChatID: ChatIDInt(12345), MessageIDs: []int64{77, 78}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Fatal("expected true result")
	}
}

func TestForwardMessagesValidation(t *testing.T) {
	const token = "123:secret"
	bot := newTestBot(t, token, "https://example.test", nil)
	tests := []struct {
		name   string
		params ForwardMessagesParams
	}{
		{name: "empty chat", params: ForwardMessagesParams{FromChatID: ChatIDInt(2), MessageIDs: []int64{1}}},
		{name: "empty from chat", params: ForwardMessagesParams{ChatID: ChatIDInt(1), MessageIDs: []int64{1}}},
		{name: "empty ids", params: ForwardMessagesParams{ChatID: ChatIDInt(1), FromChatID: ChatIDInt(2)}},
		{name: "too many ids", params: ForwardMessagesParams{ChatID: ChatIDInt(1), FromChatID: ChatIDInt(2), MessageIDs: makeMessageIDs(101)}},
		{name: "zero id", params: ForwardMessagesParams{ChatID: ChatIDInt(1), FromChatID: ChatIDInt(2), MessageIDs: []int64{1, 0}}},
		{name: "negative thread", params: ForwardMessagesParams{ChatID: ChatIDInt(1), FromChatID: ChatIDInt(2), MessageIDs: []int64{1}, MessageThreadID: -1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			messageIDs, err := bot.ForwardMessages(context.Background(), tt.params)
			if err == nil {
				t.Fatal("expected error")
			}
			if messageIDs != nil {
				t.Fatalf("expected nil message ids, got %+v", messageIDs)
			}
			assertNoToken(t, err, token)
		})
	}
}

func TestCopyMessagesValidation(t *testing.T) {
	const token = "123:secret"
	bot := newTestBot(t, token, "https://example.test", nil)
	tests := []struct {
		name   string
		params CopyMessagesParams
	}{
		{name: "empty chat", params: CopyMessagesParams{FromChatID: ChatIDInt(2), MessageIDs: []int64{1}}},
		{name: "empty from chat", params: CopyMessagesParams{ChatID: ChatIDInt(1), MessageIDs: []int64{1}}},
		{name: "empty ids", params: CopyMessagesParams{ChatID: ChatIDInt(1), FromChatID: ChatIDInt(2)}},
		{name: "too many ids", params: CopyMessagesParams{ChatID: ChatIDInt(1), FromChatID: ChatIDInt(2), MessageIDs: makeMessageIDs(101)}},
		{name: "zero id", params: CopyMessagesParams{ChatID: ChatIDInt(1), FromChatID: ChatIDInt(2), MessageIDs: []int64{1, 0}}},
		{name: "negative thread", params: CopyMessagesParams{ChatID: ChatIDInt(1), FromChatID: ChatIDInt(2), MessageIDs: []int64{1}, MessageThreadID: -1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			messageIDs, err := bot.CopyMessages(context.Background(), tt.params)
			if err == nil {
				t.Fatal("expected error")
			}
			if messageIDs != nil {
				t.Fatalf("expected nil message ids, got %+v", messageIDs)
			}
			assertNoToken(t, err, token)
		})
	}
}

func TestDeleteMessagesValidation(t *testing.T) {
	const token = "123:secret"
	bot := newTestBot(t, token, "https://example.test", nil)
	tests := []struct {
		name   string
		params DeleteMessagesParams
	}{
		{name: "empty chat", params: DeleteMessagesParams{MessageIDs: []int64{1}}},
		{name: "empty ids", params: DeleteMessagesParams{ChatID: ChatIDInt(1)}},
		{name: "too many ids", params: DeleteMessagesParams{ChatID: ChatIDInt(1), MessageIDs: makeMessageIDs(101)}},
		{name: "zero id", params: DeleteMessagesParams{ChatID: ChatIDInt(1), MessageIDs: []int64{1, 0}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ok, err := bot.DeleteMessages(context.Background(), tt.params)
			if err == nil {
				t.Fatal("expected error")
			}
			if ok {
				t.Fatal("expected false result")
			}
			assertNoToken(t, err, token)
		})
	}
}

func TestBatchMessageMethodsReturnAPIError(t *testing.T) {
	const token = "123:secret"
	tests := []struct {
		name   string
		method string
		call   func(*Bot) error
	}{
		{name: "forward", method: "forwardMessages", call: func(bot *Bot) error {
			_, err := bot.ForwardMessages(context.Background(), ForwardMessagesParams{ChatID: ChatIDInt(123), FromChatID: ChatIDInt(456), MessageIDs: []int64{1}})
			return err
		}},
		{name: "copy", method: "copyMessages", call: func(bot *Bot) error {
			_, err := bot.CopyMessages(context.Background(), CopyMessagesParams{ChatID: ChatIDInt(123), FromChatID: ChatIDInt(456), MessageIDs: []int64{1}})
			return err
		}},
		{name: "delete", method: "deleteMessages", call: func(bot *Bot) error {
			_, err := bot.DeleteMessages(context.Background(), DeleteMessagesParams{ChatID: ChatIDInt(123), MessageIDs: []int64{1}})
			return err
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path != "/bot"+token+"/"+tt.method {
					t.Fatalf("unexpected path: %q", r.URL.Path)
				}
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(`{"ok":false,"error_code":400,"description":"Bad Request"}`))
			}))
			defer server.Close()

			bot := newTestBot(t, token, server.URL, server.Client())
			err := tt.call(bot)
			if err == nil {
				t.Fatal("expected error")
			}
			var apiErr *apierrors.APIError
			if !stderrors.As(err, &apiErr) {
				t.Fatalf("expected APIError, got %T", err)
			}
			assertNoToken(t, err, token)
		})
	}
}

func TestBatchMessageMethodsResponseAndContextErrors(t *testing.T) {
	const token = "123:secret"
	tests := []struct {
		name   string
		method string
		call   func(context.Context, *Bot) error
	}{
		{name: "forward", method: "forwardMessages", call: func(ctx context.Context, bot *Bot) error {
			_, err := bot.ForwardMessages(ctx, ForwardMessagesParams{ChatID: ChatIDInt(123), FromChatID: ChatIDInt(456), MessageIDs: []int64{1}})
			return err
		}},
		{name: "copy", method: "copyMessages", call: func(ctx context.Context, bot *Bot) error {
			_, err := bot.CopyMessages(ctx, CopyMessagesParams{ChatID: ChatIDInt(123), FromChatID: ChatIDInt(456), MessageIDs: []int64{1}})
			return err
		}},
		{name: "delete", method: "deleteMessages", call: func(ctx context.Context, bot *Bot) error {
			_, err := bot.DeleteMessages(ctx, DeleteMessagesParams{ChatID: ChatIDInt(123), MessageIDs: []int64{1}})
			return err
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name+" cancelled context", func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				t.Fatal("request should not reach server")
			}))
			defer server.Close()
			bot := newTestBot(t, token, server.URL, server.Client())
			ctx, cancel := context.WithCancel(context.Background())
			cancel()
			err := tt.call(ctx, bot)
			if err == nil {
				t.Fatal("expected error")
			}
			assertNoToken(t, err, token)
		})

		t.Run(tt.name+" invalid json", func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path != "/bot"+token+"/"+tt.method {
					t.Fatalf("unexpected path: %q", r.URL.Path)
				}
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(`not-json`))
			}))
			defer server.Close()
			bot := newTestBot(t, token, server.URL, server.Client())
			err := tt.call(context.Background(), bot)
			if err == nil {
				t.Fatal("expected error")
			}
			assertNoToken(t, err, token)
		})

		t.Run(tt.name+" http status", func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path != "/bot"+token+"/"+tt.method {
					t.Fatalf("unexpected path: %q", r.URL.Path)
				}
				http.Error(w, "server error", http.StatusInternalServerError)
			}))
			defer server.Close()
			bot := newTestBot(t, token, server.URL, server.Client())
			err := tt.call(context.Background(), bot)
			if err == nil {
				t.Fatal("expected error")
			}
			assertNoToken(t, err, token)
		})
	}
}

func testBatchMessageIDsSuccess(t *testing.T, method string, call func(*Bot) ([]int64, error), checkPayload func(*testing.T, map[string]any)) {
	t.Helper()
	const token = "123:secret"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/bot"+token+"/"+method {
			t.Fatalf("unexpected path: %q", r.URL.Path)
		}
		var payload map[string]any
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode payload: %v", err)
		}
		checkPayload(t, payload)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true,"result":[{"message_id":88},{"message_id":89}]}`))
	}))
	defer server.Close()

	bot := newTestBot(t, token, server.URL, server.Client())
	messageIDs, err := call(bot)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(messageIDs) != 2 || messageIDs[0] != 88 || messageIDs[1] != 89 {
		t.Fatalf("unexpected message IDs: %#v", messageIDs)
	}
}

func assertBatchBasePayload(t *testing.T, payload map[string]any) {
	t.Helper()
	if payload["chat_id"] != float64(12345) || payload["from_chat_id"] != float64(67890) || payload["message_thread_id"] != float64(11) {
		t.Fatalf("unexpected base payload: %#v", payload)
	}
	assertMessageIDsPayload(t, payload["message_ids"])
}

func assertMessageIDsPayload(t *testing.T, value any) {
	t.Helper()
	items, ok := value.([]any)
	if !ok || len(items) != 2 || items[0] != float64(77) || items[1] != float64(78) {
		t.Fatalf("unexpected message_ids: %#v", value)
	}
}

func collectMessageIDs(messageIDs []telegram.MessageID) []int64 {
	if messageIDs == nil {
		return nil
	}
	ids := make([]int64, len(messageIDs))
	for i, messageID := range messageIDs {
		ids[i] = messageID.MessageID
	}
	return ids
}

func makeMessageIDs(count int) []int64 {
	ids := make([]int64, count)
	for i := range ids {
		ids[i] = int64(i + 1)
	}
	return ids
}
