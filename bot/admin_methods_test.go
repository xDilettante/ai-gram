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

func TestPromoteChatMemberSendsPayloadAndDecodesResult(t *testing.T) {
	const token = "123:secret"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/bot"+token+"/promoteChatMember" {
			t.Fatalf("unexpected path: %q", r.URL.Path)
		}
		var payload map[string]any
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode payload: %v", err)
		}
		if payload["chat_id"] != float64(12345) || payload["user_id"] != float64(777) {
			t.Fatalf("unexpected payload: %#v", payload)
		}
		wantTrueFields := []string{
			"is_anonymous",
			"can_manage_chat",
			"can_delete_messages",
			"can_manage_video_chats",
			"can_restrict_members",
			"can_promote_members",
			"can_change_info",
			"can_invite_users",
			"can_post_stories",
			"can_edit_stories",
			"can_delete_stories",
			"can_post_messages",
			"can_edit_messages",
			"can_pin_messages",
			"can_manage_topics",
			"can_manage_direct_messages",
			"can_manage_tags",
		}
		for _, field := range wantTrueFields {
			if payload[field] != true {
				t.Fatalf("expected %s=true in payload: %#v", field, payload)
			}
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true,"result":true}`))
	}))
	defer server.Close()

	bot := newTestBot(t, token, server.URL, server.Client())
	ok, err := bot.PromoteChatMember(context.Background(), PromoteChatMemberParams{
		ChatID:                  ChatIDInt(12345),
		UserID:                  777,
		IsAnonymous:             true,
		CanManageChat:           true,
		CanDeleteMessages:       true,
		CanManageVideoChats:     true,
		CanRestrictMembers:      true,
		CanPromoteMembers:       true,
		CanChangeInfo:           true,
		CanInviteUsers:          true,
		CanPostStories:          true,
		CanEditStories:          true,
		CanDeleteStories:        true,
		CanPostMessages:         true,
		CanEditMessages:         true,
		CanPinMessages:          true,
		CanManageTopics:         true,
		CanManageDirectMessages: true,
		CanManageTags:           true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Fatal("expected true result")
	}
}

func TestSetChatAdministratorCustomTitleSendsPayloadAndDecodesResult(t *testing.T) {
	const token = "123:secret"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/bot"+token+"/setChatAdministratorCustomTitle" {
			t.Fatalf("unexpected path: %q", r.URL.Path)
		}
		var payload map[string]any
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode payload: %v", err)
		}
		if payload["chat_id"] != float64(12345) || payload["user_id"] != float64(777) || payload["custom_title"] != "moderator" {
			t.Fatalf("unexpected payload: %#v", payload)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true,"result":true}`))
	}))
	defer server.Close()

	bot := newTestBot(t, token, server.URL, server.Client())
	ok, err := bot.SetChatAdministratorCustomTitle(context.Background(), SetChatAdministratorCustomTitleParams{ChatID: ChatIDInt(12345), UserID: 777, CustomTitle: "moderator"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Fatal("expected true result")
	}
}

func TestSetChatPermissionsSendsPayloadAndDecodesResult(t *testing.T) {
	const token = "123:secret"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/bot"+token+"/setChatPermissions" {
			t.Fatalf("unexpected path: %q", r.URL.Path)
		}
		var payload map[string]any
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode payload: %v", err)
		}
		permissions, ok := payload["permissions"].(map[string]any)
		if !ok {
			t.Fatalf("expected permissions object: %#v", payload)
		}
		if payload["chat_id"] != float64(12345) || payload["use_independent_chat_permissions"] != true {
			t.Fatalf("unexpected payload: %#v", payload)
		}
		if permissions["can_send_messages"] != true || permissions["can_send_photos"] != true || permissions["can_manage_topics"] != true {
			t.Fatalf("unexpected permissions: %#v", permissions)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true,"result":true}`))
	}))
	defer server.Close()

	bot := newTestBot(t, token, server.URL, server.Client())
	ok, err := bot.SetChatPermissions(context.Background(), SetChatPermissionsParams{
		ChatID: ChatIDInt(12345),
		Permissions: telegram.ChatPermissions{
			CanSendMessages: true,
			CanSendPhotos:   true,
			CanManageTopics: true,
		},
		UseIndependentChatPermissions: true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Fatal("expected true result")
	}
}

func TestSetChatPermissionsAllowsZeroPermissions(t *testing.T) {
	const token = "123:secret"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/bot"+token+"/setChatPermissions" {
			t.Fatalf("unexpected path: %q", r.URL.Path)
		}
		var payload map[string]any
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode payload: %v", err)
		}
		permissions, ok := payload["permissions"].(map[string]any)
		if !ok {
			t.Fatalf("expected permissions object: %#v", payload)
		}
		if len(permissions) != 0 {
			t.Fatalf("expected zero permissions object, got %#v", permissions)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true,"result":true}`))
	}))
	defer server.Close()

	bot := newTestBot(t, token, server.URL, server.Client())
	ok, err := bot.SetChatPermissions(context.Background(), SetChatPermissionsParams{ChatID: ChatIDInt(12345)})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Fatal("expected true result")
	}
}

func TestAdminManagementMethodValidation(t *testing.T) {
	const token = "123:secret"
	bot := newTestBot(t, token, "https://example.test", nil)

	promoteTests := []struct {
		name   string
		params PromoteChatMemberParams
	}{
		{name: "empty chat", params: PromoteChatMemberParams{UserID: 1}},
		{name: "zero user", params: PromoteChatMemberParams{ChatID: ChatIDInt(123)}},
		{name: "negative user", params: PromoteChatMemberParams{ChatID: ChatIDInt(123), UserID: -1}},
	}
	for _, tt := range promoteTests {
		t.Run("promote "+tt.name, func(t *testing.T) {
			ok, err := bot.PromoteChatMember(context.Background(), tt.params)
			if err == nil {
				t.Fatal("expected error")
			}
			if ok {
				t.Fatal("expected false result")
			}
			assertNoToken(t, err, token)
		})
	}

	customTitleTests := []struct {
		name   string
		params SetChatAdministratorCustomTitleParams
	}{
		{name: "empty chat", params: SetChatAdministratorCustomTitleParams{UserID: 1, CustomTitle: "moderator"}},
		{name: "zero user", params: SetChatAdministratorCustomTitleParams{ChatID: ChatIDInt(123), CustomTitle: "moderator"}},
		{name: "negative user", params: SetChatAdministratorCustomTitleParams{ChatID: ChatIDInt(123), UserID: -1, CustomTitle: "moderator"}},
		{name: "empty custom title", params: SetChatAdministratorCustomTitleParams{ChatID: ChatIDInt(123), UserID: 1}},
	}
	for _, tt := range customTitleTests {
		t.Run("custom title "+tt.name, func(t *testing.T) {
			ok, err := bot.SetChatAdministratorCustomTitle(context.Background(), tt.params)
			if err == nil {
				t.Fatal("expected error")
			}
			if ok {
				t.Fatal("expected false result")
			}
			assertNoToken(t, err, token)
		})
	}

	ok, err := bot.SetChatPermissions(context.Background(), SetChatPermissionsParams{})
	if err == nil {
		t.Fatal("expected error")
	}
	if ok {
		t.Fatal("expected false result")
	}
	assertNoToken(t, err, token)
}

func TestAdminManagementMethodsReturnAPIError(t *testing.T) {
	const token = "123:secret"
	tests := []struct {
		name   string
		method string
		call   func(*Bot) (bool, error)
	}{
		{name: "promote", method: "promoteChatMember", call: func(bot *Bot) (bool, error) {
			return bot.PromoteChatMember(context.Background(), PromoteChatMemberParams{ChatID: ChatIDInt(123), UserID: 1})
		}},
		{name: "custom title", method: "setChatAdministratorCustomTitle", call: func(bot *Bot) (bool, error) {
			return bot.SetChatAdministratorCustomTitle(context.Background(), SetChatAdministratorCustomTitleParams{ChatID: ChatIDInt(123), UserID: 1, CustomTitle: "moderator"})
		}},
		{name: "permissions", method: "setChatPermissions", call: func(bot *Bot) (bool, error) {
			return bot.SetChatPermissions(context.Background(), SetChatPermissionsParams{ChatID: ChatIDInt(123)})
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
			ok, err := tt.call(bot)
			if err == nil {
				t.Fatal("expected error")
			}
			if ok {
				t.Fatal("expected false result")
			}
			var apiErr *apierrors.APIError
			if !stderrors.As(err, &apiErr) {
				t.Fatalf("expected APIError, got %T", err)
			}
			assertNoToken(t, err, token)
		})
	}
}

func TestAdminManagementMethodsResponseAndContextErrors(t *testing.T) {
	const token = "123:secret"
	tests := []struct {
		name   string
		method string
		call   func(context.Context, *Bot) (bool, error)
	}{
		{name: "promote", method: "promoteChatMember", call: func(ctx context.Context, bot *Bot) (bool, error) {
			return bot.PromoteChatMember(ctx, PromoteChatMemberParams{ChatID: ChatIDInt(123), UserID: 1})
		}},
		{name: "custom title", method: "setChatAdministratorCustomTitle", call: func(ctx context.Context, bot *Bot) (bool, error) {
			return bot.SetChatAdministratorCustomTitle(ctx, SetChatAdministratorCustomTitleParams{ChatID: ChatIDInt(123), UserID: 1, CustomTitle: "moderator"})
		}},
		{name: "permissions", method: "setChatPermissions", call: func(ctx context.Context, bot *Bot) (bool, error) {
			return bot.SetChatPermissions(ctx, SetChatPermissionsParams{ChatID: ChatIDInt(123)})
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
			ok, err := tt.call(ctx, bot)
			if err == nil {
				t.Fatal("expected error")
			}
			if ok {
				t.Fatal("expected false result")
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
			ok, err := tt.call(context.Background(), bot)
			if err == nil {
				t.Fatal("expected error")
			}
			if ok {
				t.Fatal("expected false result")
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
			ok, err := tt.call(context.Background(), bot)
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
