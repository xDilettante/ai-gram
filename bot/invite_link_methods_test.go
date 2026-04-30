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

func TestExportChatInviteLinkSendsPayloadAndDecodesResult(t *testing.T) {
	const token = "123:secret"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/bot"+token+"/exportChatInviteLink" {
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
		_, _ = w.Write([]byte(`{"ok":true,"result":"https://t.me/+redacted"}`))
	}))
	defer server.Close()

	bot := newTestBot(t, token, server.URL, server.Client())
	inviteLink, err := bot.ExportChatInviteLink(context.Background(), ExportChatInviteLinkParams{ChatID: ChatIDInt(12345)})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if inviteLink != "https://t.me/+redacted" {
		t.Fatalf("unexpected invite link: %q", inviteLink)
	}
}

func TestCreateChatInviteLinkSendsPayloadAndDecodesResult(t *testing.T) {
	const token = "123:secret"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/bot"+token+"/createChatInviteLink" {
			t.Fatalf("unexpected path: %q", r.URL.Path)
		}
		var payload map[string]any
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode payload: %v", err)
		}
		if payload["chat_id"] != float64(12345) ||
			payload["name"] != "test link" ||
			payload["expire_date"] != float64(1234567890) ||
			payload["member_limit"] != float64(10) ||
			payload["creates_join_request"] != true {
			t.Fatalf("unexpected payload: %#v", payload)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true,"result":{"invite_link":"https://t.me/+redacted","creator":{"id":1,"is_bot":true,"first_name":"Bot"},"creates_join_request":true,"is_primary":false,"is_revoked":false,"name":"test link","expire_date":1234567890,"member_limit":10,"pending_join_request_count":2}}`))
	}))
	defer server.Close()

	bot := newTestBot(t, token, server.URL, server.Client())
	inviteLink, err := bot.CreateChatInviteLink(context.Background(), CreateChatInviteLinkParams{
		ChatID:             ChatIDInt(12345),
		Name:               "test link",
		ExpireDate:         1234567890,
		MemberLimit:        10,
		CreatesJoinRequest: true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if inviteLink.InviteLink == "" || inviteLink.Creator.ID != 1 || !inviteLink.CreatesJoinRequest || inviteLink.IsPrimary || inviteLink.IsRevoked || inviteLink.Name != "test link" || inviteLink.ExpireDate != 1234567890 || inviteLink.MemberLimit != 10 || inviteLink.PendingJoinRequestCount != 2 {
		t.Fatalf("unexpected invite link: %#v", inviteLink)
	}
}

func TestEditChatInviteLinkSendsPayloadAndDecodesResult(t *testing.T) {
	const token = "123:secret"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/bot"+token+"/editChatInviteLink" {
			t.Fatalf("unexpected path: %q", r.URL.Path)
		}
		var payload map[string]any
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode payload: %v", err)
		}
		if payload["chat_id"] != float64(12345) ||
			payload["invite_link"] != "https://t.me/+redacted" ||
			payload["name"] != "edited link" ||
			payload["expire_date"] != float64(1234567890) ||
			payload["member_limit"] != float64(5) ||
			payload["creates_join_request"] != true {
			t.Fatalf("unexpected payload: %#v", payload)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true,"result":{"invite_link":"https://t.me/+redacted","creator":{"id":1,"is_bot":true,"first_name":"Bot"},"creates_join_request":true,"is_primary":false,"is_revoked":false,"name":"edited link","expire_date":1234567890,"member_limit":5}}`))
	}))
	defer server.Close()

	bot := newTestBot(t, token, server.URL, server.Client())
	inviteLink, err := bot.EditChatInviteLink(context.Background(), EditChatInviteLinkParams{
		ChatID:             ChatIDInt(12345),
		InviteLink:         "https://t.me/+redacted",
		Name:               "edited link",
		ExpireDate:         1234567890,
		MemberLimit:        5,
		CreatesJoinRequest: true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if inviteLink.InviteLink == "" || inviteLink.Name != "edited link" || inviteLink.MemberLimit != 5 || !inviteLink.CreatesJoinRequest {
		t.Fatalf("unexpected invite link: %#v", inviteLink)
	}
}

func TestRevokeChatInviteLinkSendsPayloadAndDecodesResult(t *testing.T) {
	const token = "123:secret"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/bot"+token+"/revokeChatInviteLink" {
			t.Fatalf("unexpected path: %q", r.URL.Path)
		}
		var payload map[string]any
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode payload: %v", err)
		}
		if payload["chat_id"] != float64(12345) || payload["invite_link"] != "https://t.me/+redacted" {
			t.Fatalf("unexpected payload: %#v", payload)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true,"result":{"invite_link":"https://t.me/+redacted","creator":{"id":1,"is_bot":true,"first_name":"Bot"},"creates_join_request":false,"is_primary":false,"is_revoked":true}}`))
	}))
	defer server.Close()

	bot := newTestBot(t, token, server.URL, server.Client())
	inviteLink, err := bot.RevokeChatInviteLink(context.Background(), RevokeChatInviteLinkParams{
		ChatID:     ChatIDInt(12345),
		InviteLink: "https://t.me/+redacted",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if inviteLink.InviteLink == "" || !inviteLink.IsRevoked || inviteLink.IsPrimary || inviteLink.CreatesJoinRequest {
		t.Fatalf("unexpected invite link: %#v", inviteLink)
	}
}

func TestInviteLinkMethodsValidation(t *testing.T) {
	const token = "123:secret"
	bot := newTestBot(t, token, "https://example.test", nil)

	t.Run("export empty chat", func(t *testing.T) {
		inviteLink, err := bot.ExportChatInviteLink(context.Background(), ExportChatInviteLinkParams{})
		if err == nil {
			t.Fatal("expected error")
		}
		if inviteLink != "" {
			t.Fatalf("expected empty invite link, got %q", inviteLink)
		}
		assertNoToken(t, err, token)
	})

	createTests := []struct {
		name   string
		params CreateChatInviteLinkParams
	}{
		{name: "empty chat", params: CreateChatInviteLinkParams{}},
		{name: "negative expire date", params: CreateChatInviteLinkParams{ChatID: ChatIDInt(123), ExpireDate: -1}},
		{name: "negative member limit", params: CreateChatInviteLinkParams{ChatID: ChatIDInt(123), MemberLimit: -1}},
	}
	for _, tt := range createTests {
		t.Run("create "+tt.name, func(t *testing.T) {
			inviteLink, err := bot.CreateChatInviteLink(context.Background(), tt.params)
			if err == nil {
				t.Fatal("expected error")
			}
			if inviteLink != nil {
				t.Fatalf("expected nil invite link, got %#v", inviteLink)
			}
			assertNoToken(t, err, token)
		})
	}

	editTests := []struct {
		name   string
		params EditChatInviteLinkParams
	}{
		{name: "empty chat", params: EditChatInviteLinkParams{InviteLink: "https://t.me/+redacted"}},
		{name: "empty invite link", params: EditChatInviteLinkParams{ChatID: ChatIDInt(123)}},
		{name: "negative expire date", params: EditChatInviteLinkParams{ChatID: ChatIDInt(123), InviteLink: "https://t.me/+redacted", ExpireDate: -1}},
		{name: "negative member limit", params: EditChatInviteLinkParams{ChatID: ChatIDInt(123), InviteLink: "https://t.me/+redacted", MemberLimit: -1}},
	}
	for _, tt := range editTests {
		t.Run("edit "+tt.name, func(t *testing.T) {
			inviteLink, err := bot.EditChatInviteLink(context.Background(), tt.params)
			if err == nil {
				t.Fatal("expected error")
			}
			if inviteLink != nil {
				t.Fatalf("expected nil invite link, got %#v", inviteLink)
			}
			assertNoToken(t, err, token)
		})
	}

	revokeTests := []struct {
		name   string
		params RevokeChatInviteLinkParams
	}{
		{name: "empty chat", params: RevokeChatInviteLinkParams{InviteLink: "https://t.me/+redacted"}},
		{name: "empty invite link", params: RevokeChatInviteLinkParams{ChatID: ChatIDInt(123)}},
	}
	for _, tt := range revokeTests {
		t.Run("revoke "+tt.name, func(t *testing.T) {
			inviteLink, err := bot.RevokeChatInviteLink(context.Background(), tt.params)
			if err == nil {
				t.Fatal("expected error")
			}
			if inviteLink != nil {
				t.Fatalf("expected nil invite link, got %#v", inviteLink)
			}
			assertNoToken(t, err, token)
		})
	}
}

func TestInviteLinkMethodsReturnAPIError(t *testing.T) {
	const token = "123:secret"
	tests := []struct {
		name   string
		method string
		call   func(*Bot) (any, error)
	}{
		{name: "export", method: "exportChatInviteLink", call: func(bot *Bot) (any, error) {
			return bot.ExportChatInviteLink(context.Background(), ExportChatInviteLinkParams{ChatID: ChatIDInt(123)})
		}},
		{name: "create", method: "createChatInviteLink", call: func(bot *Bot) (any, error) {
			return bot.CreateChatInviteLink(context.Background(), CreateChatInviteLinkParams{ChatID: ChatIDInt(123)})
		}},
		{name: "edit", method: "editChatInviteLink", call: func(bot *Bot) (any, error) {
			return bot.EditChatInviteLink(context.Background(), EditChatInviteLinkParams{ChatID: ChatIDInt(123), InviteLink: "https://t.me/+redacted"})
		}},
		{name: "revoke", method: "revokeChatInviteLink", call: func(bot *Bot) (any, error) {
			return bot.RevokeChatInviteLink(context.Background(), RevokeChatInviteLinkParams{ChatID: ChatIDInt(123), InviteLink: "https://t.me/+redacted"})
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

func TestInviteLinkMethodsResponseAndContextErrors(t *testing.T) {
	const token = "123:secret"
	tests := []struct {
		name   string
		method string
		call   func(context.Context, *Bot) (any, error)
	}{
		{name: "export", method: "exportChatInviteLink", call: func(ctx context.Context, bot *Bot) (any, error) {
			return bot.ExportChatInviteLink(ctx, ExportChatInviteLinkParams{ChatID: ChatIDInt(123)})
		}},
		{name: "create", method: "createChatInviteLink", call: func(ctx context.Context, bot *Bot) (any, error) {
			return bot.CreateChatInviteLink(ctx, CreateChatInviteLinkParams{ChatID: ChatIDInt(123)})
		}},
		{name: "edit", method: "editChatInviteLink", call: func(ctx context.Context, bot *Bot) (any, error) {
			return bot.EditChatInviteLink(ctx, EditChatInviteLinkParams{ChatID: ChatIDInt(123), InviteLink: "https://t.me/+redacted"})
		}},
		{name: "revoke", method: "revokeChatInviteLink", call: func(ctx context.Context, bot *Bot) (any, error) {
			return bot.RevokeChatInviteLink(ctx, RevokeChatInviteLinkParams{ChatID: ChatIDInt(123), InviteLink: "https://t.me/+redacted"})
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
