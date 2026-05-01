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

func TestGetChatSendsPayloadAndDecodesResult(t *testing.T) {
	const token = "123:secret"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/bot"+token+"/getChat" {
			t.Fatalf("unexpected path: %q", r.URL.Path)
		}
		var payload map[string]any
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode payload: %v", err)
		}
		if payload["chat_id"] != float64(12345) {
			t.Fatalf("unexpected payload: %#v", payload)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true,"result":{"id":12345,"type":"supergroup","title":"Test chat","username":"testchat","description":"About","invite_link":"https://t.me/+redacted","pinned_message":{"message_id":99,"chat":{"id":12345,"type":"supergroup"},"date":1,"text":"Pinned"}}}`))
	}))
	defer server.Close()

	bot := newTestBot(t, token, server.URL, server.Client())
	chat, err := bot.GetChat(context.Background(), GetChatParams{ChatID: ChatIDInt(12345)})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if chat.ID != 12345 || chat.Type != "supergroup" || chat.Title != "Test chat" || chat.Username != "testchat" || chat.Description != "About" || chat.InviteLink == "" {
		t.Fatalf("unexpected chat: %#v", chat)
	}
	if chat.PinnedMessage == nil || chat.PinnedMessage.MessageID != 99 {
		t.Fatalf("unexpected pinned message: %#v", chat.PinnedMessage)
	}
}

func TestGetChatFullInfoSendsPayloadAndDecodesResult(t *testing.T) {
	const token = "123:secret"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/bot"+token+"/getChat" {
			t.Fatalf("unexpected path: %q", r.URL.Path)
		}
		var payload map[string]any
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode payload: %v", err)
		}
		if payload["chat_id"] != float64(12345) {
			t.Fatalf("unexpected payload: %#v", payload)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true,"result":{"id":12345,"type":"supergroup","title":"Full chat","is_forum":true,"accent_color_id":7,"max_reaction_count":10,"available_reactions":[{"type":"emoji","emoji":"👍"}],"accepted_gift_types":{"unlimited_gifts":true,"limited_gifts":true,"unique_gifts":true,"premium_subscription":true,"gifts_from_channels":true},"permissions":{"can_send_messages":true},"paid_message_star_count":5}}`))
	}))
	defer server.Close()

	bot := newTestBot(t, token, server.URL, server.Client())
	chat, err := bot.GetChatFullInfo(context.Background(), GetChatParams{ChatID: ChatIDInt(12345)})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if chat.ID != 12345 || chat.Type != "supergroup" || chat.Title != "Full chat" || !chat.IsForum || chat.AccentColorID != 7 || chat.MaxReactionCount != 10 || len(chat.AvailableReactions) != 1 || chat.Permissions == nil || !chat.Permissions.CanSendMessages || !chat.AcceptedGiftTypes.GiftsFromChannels || chat.PaidMessageStarCount != 5 {
		t.Fatalf("unexpected chat full info: %#v", chat)
	}
}

func TestGetChatValidation(t *testing.T) {
	const token = "123:secret"
	bot := newTestBot(t, token, "https://example.test", nil)
	chat, err := bot.GetChat(context.Background(), GetChatParams{})
	if err == nil {
		t.Fatal("expected error")
	}
	if chat != nil {
		t.Fatalf("expected nil chat, got %#v", chat)
	}
	assertNoToken(t, err, token)

	fullInfo, err := bot.GetChatFullInfo(context.Background(), GetChatParams{})
	if err == nil {
		t.Fatal("expected full info error")
	}
	if fullInfo != nil {
		t.Fatalf("expected nil chat full info, got %#v", fullInfo)
	}
	assertNoToken(t, err, token)
}

func TestGetChatMemberSendsPayloadAndDecodesResult(t *testing.T) {
	const token = "123:secret"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/bot"+token+"/getChatMember" {
			t.Fatalf("unexpected path: %q", r.URL.Path)
		}
		var payload map[string]any
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode payload: %v", err)
		}
		if payload["chat_id"] != float64(12345) || payload["user_id"] != float64(777) {
			t.Fatalf("unexpected payload: %#v", payload)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true,"result":{"status":"administrator","user":{"id":777,"is_bot":false,"first_name":"Admin","username":"admin"},"can_manage_chat":true,"can_delete_messages":true,"can_pin_messages":true,"custom_title":"Moderator"}}`))
	}))
	defer server.Close()

	bot := newTestBot(t, token, server.URL, server.Client())
	member, err := bot.GetChatMember(context.Background(), GetChatMemberParams{ChatID: ChatIDInt(12345), UserID: 777})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if member.Status != telegram.ChatMemberStatusAdministrator || member.User.ID != 777 || !member.CanManageChat || !member.CanDeleteMessages || !member.CanPinMessages || member.CustomTitle != "Moderator" {
		t.Fatalf("unexpected member: %#v", member)
	}
}

func TestGetChatMemberValidation(t *testing.T) {
	const token = "123:secret"
	bot := newTestBot(t, token, "https://example.test", nil)
	tests := []struct {
		name   string
		params GetChatMemberParams
	}{
		{name: "empty chat", params: GetChatMemberParams{UserID: 1}},
		{name: "zero user", params: GetChatMemberParams{ChatID: ChatIDInt(123)}},
		{name: "negative user", params: GetChatMemberParams{ChatID: ChatIDInt(123), UserID: -1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			member, err := bot.GetChatMember(context.Background(), tt.params)
			if err == nil {
				t.Fatal("expected error")
			}
			if member != nil {
				t.Fatalf("expected nil member, got %#v", member)
			}
			assertNoToken(t, err, token)
		})
	}
}

func TestGetChatAdministratorsSendsPayloadAndDecodesResult(t *testing.T) {
	const token = "123:secret"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/bot"+token+"/getChatAdministrators" {
			t.Fatalf("unexpected path: %q", r.URL.Path)
		}
		var payload map[string]any
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode payload: %v", err)
		}
		if payload["chat_id"] != float64(12345) {
			t.Fatalf("unexpected payload: %#v", payload)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true,"result":[{"status":"creator","user":{"id":1,"is_bot":false,"first_name":"Owner"},"is_anonymous":true},{"status":"administrator","user":{"id":2,"is_bot":false,"first_name":"Admin"},"can_manage_video_chats":true,"can_invite_users":true}]}`))
	}))
	defer server.Close()

	bot := newTestBot(t, token, server.URL, server.Client())
	admins, err := bot.GetChatAdministrators(context.Background(), GetChatAdministratorsParams{ChatID: ChatIDInt(12345)})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(admins) != 2 || admins[0].Status != telegram.ChatMemberStatusCreator || !admins[0].IsAnonymous || admins[1].User.ID != 2 || !admins[1].CanManageVideoChats || !admins[1].CanInviteUsers {
		t.Fatalf("unexpected administrators: %#v", admins)
	}
}

func TestGetChatAdministratorsValidation(t *testing.T) {
	const token = "123:secret"
	bot := newTestBot(t, token, "https://example.test", nil)
	admins, err := bot.GetChatAdministrators(context.Background(), GetChatAdministratorsParams{})
	if err == nil {
		t.Fatal("expected error")
	}
	if admins != nil {
		t.Fatalf("expected nil administrators, got %#v", admins)
	}
	assertNoToken(t, err, token)
}

func TestGetChatMemberCountSendsPayloadAndDecodesResult(t *testing.T) {
	const token = "123:secret"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/bot"+token+"/getChatMemberCount" {
			t.Fatalf("unexpected path: %q", r.URL.Path)
		}
		var payload map[string]any
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode payload: %v", err)
		}
		if payload["chat_id"] != float64(12345) {
			t.Fatalf("unexpected payload: %#v", payload)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true,"result":42}`))
	}))
	defer server.Close()

	bot := newTestBot(t, token, server.URL, server.Client())
	count, err := bot.GetChatMemberCount(context.Background(), GetChatMemberCountParams{ChatID: ChatIDInt(12345)})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if count != 42 {
		t.Fatalf("unexpected count: %d", count)
	}
}

func TestGetChatMemberCountValidation(t *testing.T) {
	const token = "123:secret"
	bot := newTestBot(t, token, "https://example.test", nil)
	count, err := bot.GetChatMemberCount(context.Background(), GetChatMemberCountParams{})
	if err == nil {
		t.Fatal("expected error")
	}
	if count != 0 {
		t.Fatalf("expected zero count, got %d", count)
	}
	assertNoToken(t, err, token)
}

func TestChatInfoMethodsReturnAPIError(t *testing.T) {
	const token = "123:secret"
	tests := []struct {
		name   string
		method string
		call   func(*Bot) (any, error)
	}{
		{name: "get chat", method: "getChat", call: func(bot *Bot) (any, error) {
			return bot.GetChat(context.Background(), GetChatParams{ChatID: ChatIDInt(123)})
		}},
		{name: "get chat full info", method: "getChat", call: func(bot *Bot) (any, error) {
			return bot.GetChatFullInfo(context.Background(), GetChatParams{ChatID: ChatIDInt(123)})
		}},
		{name: "get chat member", method: "getChatMember", call: func(bot *Bot) (any, error) {
			return bot.GetChatMember(context.Background(), GetChatMemberParams{ChatID: ChatIDInt(123), UserID: 1})
		}},
		{name: "get administrators", method: "getChatAdministrators", call: func(bot *Bot) (any, error) {
			return bot.GetChatAdministrators(context.Background(), GetChatAdministratorsParams{ChatID: ChatIDInt(123)})
		}},
		{name: "get member count", method: "getChatMemberCount", call: func(bot *Bot) (any, error) {
			return bot.GetChatMemberCount(context.Background(), GetChatMemberCountParams{ChatID: ChatIDInt(123)})
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
			_, err := tt.call(bot)
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

func TestChatInfoMethodsResponseAndContextErrors(t *testing.T) {
	const token = "123:secret"
	tests := []struct {
		name   string
		method string
		call   func(context.Context, *Bot) (any, error)
	}{
		{name: "get chat", method: "getChat", call: func(ctx context.Context, bot *Bot) (any, error) {
			return bot.GetChat(ctx, GetChatParams{ChatID: ChatIDInt(123)})
		}},
		{name: "get chat full info", method: "getChat", call: func(ctx context.Context, bot *Bot) (any, error) {
			return bot.GetChatFullInfo(ctx, GetChatParams{ChatID: ChatIDInt(123)})
		}},
		{name: "get chat member", method: "getChatMember", call: func(ctx context.Context, bot *Bot) (any, error) {
			return bot.GetChatMember(ctx, GetChatMemberParams{ChatID: ChatIDInt(123), UserID: 1})
		}},
		{name: "get administrators", method: "getChatAdministrators", call: func(ctx context.Context, bot *Bot) (any, error) {
			return bot.GetChatAdministrators(ctx, GetChatAdministratorsParams{ChatID: ChatIDInt(123)})
		}},
		{name: "get member count", method: "getChatMemberCount", call: func(ctx context.Context, bot *Bot) (any, error) {
			return bot.GetChatMemberCount(ctx, GetChatMemberCountParams{ChatID: ChatIDInt(123)})
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
			_, err := tt.call(ctx, bot)
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
			_, err := tt.call(context.Background(), bot)
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
			_, err := tt.call(context.Background(), bot)
			if err == nil {
				t.Fatal("expected error")
			}
			assertNoToken(t, err, token)
		})
	}
}
