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

func TestGetBusinessConnectionSendsPayloadAndDecodesResult(t *testing.T) {
	const token = "123:secret"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/bot"+token+"/getBusinessConnection" {
			t.Fatalf("unexpected path: %q", r.URL.Path)
		}
		var payload map[string]any
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode payload: %v", err)
		}
		if payload["business_connection_id"] != "bc-1" {
			t.Fatalf("unexpected payload: %#v", payload)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true,"result":{"id":"bc-1","user":{"id":7,"is_bot":false,"first_name":"Business"},"user_chat_id":7000,"date":123,"rights":{"can_reply":true,"can_delete_all_messages":true},"is_enabled":true}}`))
	}))
	defer server.Close()

	bot := newTestBot(t, token, server.URL, server.Client())
	connection, err := bot.GetBusinessConnection(context.Background(), GetBusinessConnectionParams{BusinessConnectionID: "bc-1"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if connection == nil || connection.ID != "bc-1" || connection.User.ID != 7 || connection.UserChatID != 7000 || connection.Rights == nil || !connection.Rights.CanReply || !connection.Rights.CanDeleteAllMessages || !connection.IsEnabled {
		t.Fatalf("unexpected connection: %+v", connection)
	}
}

func TestDeleteBusinessMessagesSendsPayloadAndDecodesResult(t *testing.T) {
	const token = "123:secret"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/bot"+token+"/deleteBusinessMessages" {
			t.Fatalf("unexpected path: %q", r.URL.Path)
		}
		var payload map[string]any
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode payload: %v", err)
		}
		if payload["business_connection_id"] != "bc-1" {
			t.Fatalf("unexpected business_connection_id: %#v", payload)
		}
		assertMessageIDsPayload(t, payload["message_ids"])
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true,"result":true}`))
	}))
	defer server.Close()

	bot := newTestBot(t, token, server.URL, server.Client())
	ok, err := bot.DeleteBusinessMessages(context.Background(), DeleteBusinessMessagesParams{BusinessConnectionID: "bc-1", MessageIDs: []int64{77, 78}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Fatal("expected true result")
	}
}

func TestBusinessMethodsValidation(t *testing.T) {
	const token = "123:secret"
	bot := newTestBot(t, token, "https://example.test", nil)

	if connection, err := bot.GetBusinessConnection(context.Background(), GetBusinessConnectionParams{}); err == nil {
		t.Fatal("expected getBusinessConnection validation error")
	} else {
		if connection != nil {
			t.Fatalf("expected nil connection, got %+v", connection)
		}
		assertNoToken(t, err, token)
	}

	tests := []struct {
		name   string
		params DeleteBusinessMessagesParams
	}{
		{name: "missing connection id", params: DeleteBusinessMessagesParams{MessageIDs: []int64{1}}},
		{name: "empty ids", params: DeleteBusinessMessagesParams{BusinessConnectionID: "bc-1"}},
		{name: "too many ids", params: DeleteBusinessMessagesParams{BusinessConnectionID: "bc-1", MessageIDs: makeMessageIDs(101)}},
		{name: "zero id", params: DeleteBusinessMessagesParams{BusinessConnectionID: "bc-1", MessageIDs: []int64{1, 0}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ok, err := bot.DeleteBusinessMessages(context.Background(), tt.params)
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

func TestBusinessMethodsReturnAPIError(t *testing.T) {
	const token = "123:secret"
	tests := []struct {
		name   string
		method string
		call   func(*Bot) error
	}{
		{name: "get", method: "getBusinessConnection", call: func(bot *Bot) error {
			_, err := bot.GetBusinessConnection(context.Background(), GetBusinessConnectionParams{BusinessConnectionID: "bc-1"})
			return err
		}},
		{name: "delete", method: "deleteBusinessMessages", call: func(bot *Bot) error {
			_, err := bot.DeleteBusinessMessages(context.Background(), DeleteBusinessMessagesParams{BusinessConnectionID: "bc-1", MessageIDs: []int64{1}})
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

func TestBusinessMethodsResponseAndContextErrors(t *testing.T) {
	const token = "123:secret"
	tests := []struct {
		name   string
		method string
		call   func(context.Context, *Bot) error
	}{
		{name: "get", method: "getBusinessConnection", call: func(ctx context.Context, bot *Bot) error {
			_, err := bot.GetBusinessConnection(ctx, GetBusinessConnectionParams{BusinessConnectionID: "bc-1"})
			return err
		}},
		{name: "delete", method: "deleteBusinessMessages", call: func(ctx context.Context, bot *Bot) error {
			_, err := bot.DeleteBusinessMessages(ctx, DeleteBusinessMessagesParams{BusinessConnectionID: "bc-1", MessageIDs: []int64{1}})
			return err
		}},
	}
	for _, tt := range tests {
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
			if err := tt.call(context.Background(), bot); err == nil {
				t.Fatal("expected error")
			} else {
				assertNoToken(t, err, token)
			}
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
			if err := tt.call(context.Background(), bot); err == nil {
				t.Fatal("expected error")
			} else {
				assertNoToken(t, err, token)
			}
		})

		t.Run(tt.name+" cancelled context", func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				t.Fatal("request should not reach server")
			}))
			defer server.Close()
			bot := newTestBot(t, token, server.URL, server.Client())
			ctx, cancel := context.WithCancel(context.Background())
			cancel()
			if err := tt.call(ctx, bot); err == nil {
				t.Fatal("expected error")
			} else {
				assertNoToken(t, err, token)
			}
		})
	}
}
