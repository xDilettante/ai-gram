package bot

import (
	"context"
	"encoding/json"
	stderrors "errors"
	"net/http"
	"net/http/httptest"
	"testing"

	apierrors "github.com/xDilettante/ai-gram/errors"
)

func TestCreateForumTopicSendsPayloadAndDecodesResult(t *testing.T) {
	const token = "123:secret"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/bot"+token+"/createForumTopic" {
			t.Fatalf("unexpected path: %q", r.URL.Path)
		}
		var payload map[string]any
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode payload: %v", err)
		}
		if payload["chat_id"] != float64(12345) || payload["name"] != "News" || payload["icon_color"] != float64(0x6fb9f0) || payload["icon_custom_emoji_id"] != "emoji-id" {
			t.Fatalf("unexpected payload: %#v", payload)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true,"result":{"message_thread_id":777,"name":"News","icon_color":7322096,"icon_custom_emoji_id":"emoji-id"}}`))
	}))
	defer server.Close()

	bot := newTestBot(t, token, server.URL, server.Client())
	topic, err := bot.CreateForumTopic(context.Background(), CreateForumTopicParams{
		ChatID:            ChatIDInt(12345),
		Name:              "News",
		IconColor:         0x6fb9f0,
		IconCustomEmojiID: "emoji-id",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if topic == nil || topic.MessageThreadID != 777 || topic.Name != "News" || topic.IconColor != 7322096 || topic.IconCustomEmojiID != "emoji-id" {
		t.Fatalf("unexpected topic: %+v", topic)
	}
}

func TestEditForumTopicSendsPayloadAndDecodesResult(t *testing.T) {
	testForumTopicBoolSuccess(t, "editForumTopic", func(bot *Bot) (bool, error) {
		return bot.EditForumTopic(context.Background(), EditForumTopicParams{
			ChatID:            ChatIDInt(12345),
			MessageThreadID:   777,
			Name:              "Renamed",
			IconCustomEmojiID: "emoji-id",
		})
	}, func(t *testing.T, payload map[string]any) {
		if payload["chat_id"] != float64(12345) || payload["message_thread_id"] != float64(777) || payload["name"] != "Renamed" || payload["icon_custom_emoji_id"] != "emoji-id" {
			t.Fatalf("unexpected payload: %#v", payload)
		}
	})
}

func TestForumTopicTargetMethodsSendPayloadAndDecodeResult(t *testing.T) {
	tests := []struct {
		name   string
		method string
		call   func(*Bot) (bool, error)
	}{
		{name: "close", method: "closeForumTopic", call: func(bot *Bot) (bool, error) {
			return bot.CloseForumTopic(context.Background(), CloseForumTopicParams{ChatID: ChatIDInt(12345), MessageThreadID: 777})
		}},
		{name: "reopen", method: "reopenForumTopic", call: func(bot *Bot) (bool, error) {
			return bot.ReopenForumTopic(context.Background(), ReopenForumTopicParams{ChatID: ChatIDInt(12345), MessageThreadID: 777})
		}},
		{name: "delete", method: "deleteForumTopic", call: func(bot *Bot) (bool, error) {
			return bot.DeleteForumTopic(context.Background(), DeleteForumTopicParams{ChatID: ChatIDInt(12345), MessageThreadID: 777})
		}},
		{name: "unpin all", method: "unpinAllForumTopicMessages", call: func(bot *Bot) (bool, error) {
			return bot.UnpinAllForumTopicMessages(context.Background(), UnpinAllForumTopicMessagesParams{ChatID: ChatIDInt(12345), MessageThreadID: 777})
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testForumTopicBoolSuccess(t, tt.method, tt.call, func(t *testing.T, payload map[string]any) {
				if payload["chat_id"] != float64(12345) || payload["message_thread_id"] != float64(777) {
					t.Fatalf("unexpected payload: %#v", payload)
				}
			})
		})
	}
}

func TestEditGeneralForumTopicSendsPayloadAndDecodesResult(t *testing.T) {
	testForumTopicBoolSuccess(t, "editGeneralForumTopic", func(bot *Bot) (bool, error) {
		return bot.EditGeneralForumTopic(context.Background(), EditGeneralForumTopicParams{ChatID: ChatIDInt(12345), Name: "General"})
	}, func(t *testing.T, payload map[string]any) {
		if payload["chat_id"] != float64(12345) || payload["name"] != "General" {
			t.Fatalf("unexpected payload: %#v", payload)
		}
	})
}

func TestGeneralForumTopicMethodsSendPayloadAndDecodeResult(t *testing.T) {
	tests := []struct {
		name   string
		method string
		call   func(*Bot) (bool, error)
	}{
		{name: "close", method: "closeGeneralForumTopic", call: func(bot *Bot) (bool, error) {
			return bot.CloseGeneralForumTopic(context.Background(), CloseGeneralForumTopicParams{ChatID: ChatIDInt(12345)})
		}},
		{name: "reopen", method: "reopenGeneralForumTopic", call: func(bot *Bot) (bool, error) {
			return bot.ReopenGeneralForumTopic(context.Background(), ReopenGeneralForumTopicParams{ChatID: ChatIDInt(12345)})
		}},
		{name: "hide", method: "hideGeneralForumTopic", call: func(bot *Bot) (bool, error) {
			return bot.HideGeneralForumTopic(context.Background(), HideGeneralForumTopicParams{ChatID: ChatIDInt(12345)})
		}},
		{name: "unhide", method: "unhideGeneralForumTopic", call: func(bot *Bot) (bool, error) {
			return bot.UnhideGeneralForumTopic(context.Background(), UnhideGeneralForumTopicParams{ChatID: ChatIDInt(12345)})
		}},
		{name: "unpin all", method: "unpinAllGeneralForumTopicMessages", call: func(bot *Bot) (bool, error) {
			return bot.UnpinAllGeneralForumTopicMessages(context.Background(), UnpinAllGeneralForumTopicMessagesParams{ChatID: ChatIDInt(12345)})
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testForumTopicBoolSuccess(t, tt.method, tt.call, func(t *testing.T, payload map[string]any) {
				if payload["chat_id"] != float64(12345) {
					t.Fatalf("unexpected payload: %#v", payload)
				}
			})
		})
	}
}

func TestForumTopicMethodValidation(t *testing.T) {
	const token = "123:secret"
	bot := newTestBot(t, token, "https://example.test", nil)
	tests := []struct {
		name string
		call func() (bool, error)
	}{
		{name: "create empty chat", call: func() (bool, error) {
			_, err := bot.CreateForumTopic(context.Background(), CreateForumTopicParams{Name: "Topic"})
			return false, err
		}},
		{name: "create empty name", call: func() (bool, error) {
			_, err := bot.CreateForumTopic(context.Background(), CreateForumTopicParams{ChatID: ChatIDInt(123)})
			return false, err
		}},
		{name: "create negative icon color", call: func() (bool, error) {
			_, err := bot.CreateForumTopic(context.Background(), CreateForumTopicParams{ChatID: ChatIDInt(123), Name: "Topic", IconColor: -1})
			return false, err
		}},
		{name: "edit empty chat", call: func() (bool, error) {
			return bot.EditForumTopic(context.Background(), EditForumTopicParams{MessageThreadID: 1})
		}},
		{name: "edit zero thread", call: func() (bool, error) {
			return bot.EditForumTopic(context.Background(), EditForumTopicParams{ChatID: ChatIDInt(123)})
		}},
		{name: "close empty chat", call: func() (bool, error) {
			return bot.CloseForumTopic(context.Background(), CloseForumTopicParams{MessageThreadID: 1})
		}},
		{name: "close zero thread", call: func() (bool, error) {
			return bot.CloseForumTopic(context.Background(), CloseForumTopicParams{ChatID: ChatIDInt(123)})
		}},
		{name: "reopen negative thread", call: func() (bool, error) {
			return bot.ReopenForumTopic(context.Background(), ReopenForumTopicParams{ChatID: ChatIDInt(123), MessageThreadID: -1})
		}},
		{name: "delete zero thread", call: func() (bool, error) {
			return bot.DeleteForumTopic(context.Background(), DeleteForumTopicParams{ChatID: ChatIDInt(123)})
		}},
		{name: "unpin zero thread", call: func() (bool, error) {
			return bot.UnpinAllForumTopicMessages(context.Background(), UnpinAllForumTopicMessagesParams{ChatID: ChatIDInt(123)})
		}},
		{name: "edit general empty chat", call: func() (bool, error) {
			return bot.EditGeneralForumTopic(context.Background(), EditGeneralForumTopicParams{Name: "General"})
		}},
		{name: "edit general empty name", call: func() (bool, error) {
			return bot.EditGeneralForumTopic(context.Background(), EditGeneralForumTopicParams{ChatID: ChatIDInt(123)})
		}},
		{name: "close general empty chat", call: func() (bool, error) {
			return bot.CloseGeneralForumTopic(context.Background(), CloseGeneralForumTopicParams{})
		}},
		{name: "reopen general empty chat", call: func() (bool, error) {
			return bot.ReopenGeneralForumTopic(context.Background(), ReopenGeneralForumTopicParams{})
		}},
		{name: "hide general empty chat", call: func() (bool, error) {
			return bot.HideGeneralForumTopic(context.Background(), HideGeneralForumTopicParams{})
		}},
		{name: "unhide general empty chat", call: func() (bool, error) {
			return bot.UnhideGeneralForumTopic(context.Background(), UnhideGeneralForumTopicParams{})
		}},
		{name: "unpin general empty chat", call: func() (bool, error) {
			return bot.UnpinAllGeneralForumTopicMessages(context.Background(), UnpinAllGeneralForumTopicMessagesParams{})
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ok, err := tt.call()
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

func TestForumTopicMethodsReturnAPIError(t *testing.T) {
	const token = "123:secret"
	for _, tt := range forumTopicErrorCases() {
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
			err := tt.call(context.Background(), bot)
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

func TestForumTopicMethodsResponseAndContextErrors(t *testing.T) {
	const token = "123:secret"
	for _, tt := range forumTopicErrorCases() {
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

func testForumTopicBoolSuccess(t *testing.T, method string, call func(*Bot) (bool, error), checkPayload func(*testing.T, map[string]any)) {
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
		_, _ = w.Write([]byte(`{"ok":true,"result":true}`))
	}))
	defer server.Close()

	bot := newTestBot(t, token, server.URL, server.Client())
	ok, err := call(bot)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Fatal("expected true result")
	}
}

type forumTopicErrorCase struct {
	name   string
	method string
	call   func(context.Context, *Bot) error
}

func forumTopicErrorCases() []forumTopicErrorCase {
	return []forumTopicErrorCase{
		{name: "create", method: "createForumTopic", call: func(ctx context.Context, bot *Bot) error {
			_, err := bot.CreateForumTopic(ctx, CreateForumTopicParams{ChatID: ChatIDInt(123), Name: "Topic"})
			return err
		}},
		{name: "edit", method: "editForumTopic", call: func(ctx context.Context, bot *Bot) error {
			_, err := bot.EditForumTopic(ctx, EditForumTopicParams{ChatID: ChatIDInt(123), MessageThreadID: 1})
			return err
		}},
		{name: "close", method: "closeForumTopic", call: func(ctx context.Context, bot *Bot) error {
			_, err := bot.CloseForumTopic(ctx, CloseForumTopicParams{ChatID: ChatIDInt(123), MessageThreadID: 1})
			return err
		}},
		{name: "reopen", method: "reopenForumTopic", call: func(ctx context.Context, bot *Bot) error {
			_, err := bot.ReopenForumTopic(ctx, ReopenForumTopicParams{ChatID: ChatIDInt(123), MessageThreadID: 1})
			return err
		}},
		{name: "delete", method: "deleteForumTopic", call: func(ctx context.Context, bot *Bot) error {
			_, err := bot.DeleteForumTopic(ctx, DeleteForumTopicParams{ChatID: ChatIDInt(123), MessageThreadID: 1})
			return err
		}},
		{name: "unpin all", method: "unpinAllForumTopicMessages", call: func(ctx context.Context, bot *Bot) error {
			_, err := bot.UnpinAllForumTopicMessages(ctx, UnpinAllForumTopicMessagesParams{ChatID: ChatIDInt(123), MessageThreadID: 1})
			return err
		}},
		{name: "edit general", method: "editGeneralForumTopic", call: func(ctx context.Context, bot *Bot) error {
			_, err := bot.EditGeneralForumTopic(ctx, EditGeneralForumTopicParams{ChatID: ChatIDInt(123), Name: "General"})
			return err
		}},
		{name: "close general", method: "closeGeneralForumTopic", call: func(ctx context.Context, bot *Bot) error {
			_, err := bot.CloseGeneralForumTopic(ctx, CloseGeneralForumTopicParams{ChatID: ChatIDInt(123)})
			return err
		}},
		{name: "reopen general", method: "reopenGeneralForumTopic", call: func(ctx context.Context, bot *Bot) error {
			_, err := bot.ReopenGeneralForumTopic(ctx, ReopenGeneralForumTopicParams{ChatID: ChatIDInt(123)})
			return err
		}},
		{name: "hide general", method: "hideGeneralForumTopic", call: func(ctx context.Context, bot *Bot) error {
			_, err := bot.HideGeneralForumTopic(ctx, HideGeneralForumTopicParams{ChatID: ChatIDInt(123)})
			return err
		}},
		{name: "unhide general", method: "unhideGeneralForumTopic", call: func(ctx context.Context, bot *Bot) error {
			_, err := bot.UnhideGeneralForumTopic(ctx, UnhideGeneralForumTopicParams{ChatID: ChatIDInt(123)})
			return err
		}},
		{name: "unpin general", method: "unpinAllGeneralForumTopicMessages", call: func(ctx context.Context, bot *Bot) error {
			_, err := bot.UnpinAllGeneralForumTopicMessages(ctx, UnpinAllGeneralForumTopicMessagesParams{ChatID: ChatIDInt(123)})
			return err
		}},
	}
}
