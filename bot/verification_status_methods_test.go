package bot

import (
	"context"
	"encoding/json"
	stderrors "errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	apierrors "github.com/xDilettante/ai-gram/errors"
)

func TestSetUserEmojiStatusSendsPayloadAndReturnsTrue(t *testing.T) {
	const token = "123:secret"
	const customEmojiID = "emoji-status-id"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/bot"+token+"/setUserEmojiStatus" {
			t.Fatalf("unexpected path: %q", r.URL.Path)
		}
		var payload map[string]any
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode payload: %v", err)
		}
		if payload["user_id"] != float64(123) || payload["emoji_status_custom_emoji_id"] != customEmojiID || payload["emoji_status_expiration_date"] != float64(456) {
			t.Fatalf("unexpected payload: %#v", payload)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true,"result":true}`))
	}))
	defer server.Close()

	bot := newTestBot(t, token, server.URL, server.Client())
	ok, err := bot.SetUserEmojiStatus(context.Background(), SetUserEmojiStatusParams{UserID: 123, EmojiStatusCustomEmojiID: customEmojiID, EmojiStatusExpirationDate: 456})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Fatal("expected true result")
	}
}

func TestSetUserEmojiStatusAllowsEmptyEmojiID(t *testing.T) {
	const token = "123:secret"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/bot"+token+"/setUserEmojiStatus" {
			t.Fatalf("unexpected path: %q", r.URL.Path)
		}
		var payload map[string]any
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode payload: %v", err)
		}
		if payload["user_id"] != float64(123) {
			t.Fatalf("unexpected payload: %#v", payload)
		}
		if _, ok := payload["emoji_status_custom_emoji_id"]; ok {
			t.Fatalf("empty emoji status id should be omitted by JSON tag: %#v", payload)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true,"result":true}`))
	}))
	defer server.Close()

	bot := newTestBot(t, token, server.URL, server.Client())
	ok, err := bot.SetUserEmojiStatus(context.Background(), SetUserEmojiStatusParams{UserID: 123})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Fatal("expected true result")
	}
}

func TestVerifyMethodsSendPayloadAndReturnTrue(t *testing.T) {
	const token = "123:secret"
	tests := []struct {
		name       string
		method     string
		call       func(*Bot) (bool, error)
		assertBody func(*testing.T, map[string]any)
	}{
		{
			name:   "verify user",
			method: "verifyUser",
			call: func(bot *Bot) (bool, error) {
				return bot.VerifyUser(context.Background(), VerifyUserParams{UserID: 123, CustomDescription: "Test organization"})
			},
			assertBody: func(t *testing.T, payload map[string]any) {
				t.Helper()
				if payload["user_id"] != float64(123) || payload["custom_description"] != "Test organization" {
					t.Fatalf("unexpected payload: %#v", payload)
				}
			},
		},
		{
			name:   "verify chat",
			method: "verifyChat",
			call: func(bot *Bot) (bool, error) {
				return bot.VerifyChat(context.Background(), VerifyChatParams{ChatID: ChatIDString("@testchannel"), CustomDescription: "Test channel"})
			},
			assertBody: func(t *testing.T, payload map[string]any) {
				t.Helper()
				if payload["chat_id"] != "@testchannel" || payload["custom_description"] != "Test channel" {
					t.Fatalf("unexpected payload: %#v", payload)
				}
			},
		},
		{
			name:   "remove user verification",
			method: "removeUserVerification",
			call: func(bot *Bot) (bool, error) {
				return bot.RemoveUserVerification(context.Background(), RemoveUserVerificationParams{UserID: 123})
			},
			assertBody: func(t *testing.T, payload map[string]any) {
				t.Helper()
				if payload["user_id"] != float64(123) {
					t.Fatalf("unexpected payload: %#v", payload)
				}
			},
		},
		{
			name:   "remove chat verification",
			method: "removeChatVerification",
			call: func(bot *Bot) (bool, error) {
				return bot.RemoveChatVerification(context.Background(), RemoveChatVerificationParams{ChatID: ChatIDInt(-100123)})
			},
			assertBody: func(t *testing.T, payload map[string]any) {
				t.Helper()
				if payload["chat_id"] != float64(-100123) {
					t.Fatalf("unexpected payload: %#v", payload)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path != "/bot"+token+"/"+tt.method {
					t.Fatalf("unexpected path: %q", r.URL.Path)
				}
				var payload map[string]any
				if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
					t.Fatalf("decode payload: %v", err)
				}
				tt.assertBody(t, payload)
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(`{"ok":true,"result":true}`))
			}))
			defer server.Close()

			bot := newTestBot(t, token, server.URL, server.Client())
			ok, err := tt.call(bot)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !ok {
				t.Fatal("expected true result")
			}
		})
	}
}

func TestVerificationStatusValidation(t *testing.T) {
	const token = "123:secret"
	bot := newTestBot(t, token, "https://example.test", nil)
	tests := []struct {
		name string
		call func() error
	}{
		{name: "emoji status invalid user", call: func() error {
			_, err := bot.SetUserEmojiStatus(context.Background(), SetUserEmojiStatusParams{})
			return err
		}},
		{name: "emoji status negative expiration", call: func() error {
			_, err := bot.SetUserEmojiStatus(context.Background(), SetUserEmojiStatusParams{UserID: 1, EmojiStatusExpirationDate: -1})
			return err
		}},
		{name: "verify user invalid user", call: func() error {
			_, err := bot.VerifyUser(context.Background(), VerifyUserParams{})
			return err
		}},
		{name: "verify chat invalid chat", call: func() error {
			_, err := bot.VerifyChat(context.Background(), VerifyChatParams{})
			return err
		}},
		{name: "remove user verification invalid user", call: func() error {
			_, err := bot.RemoveUserVerification(context.Background(), RemoveUserVerificationParams{})
			return err
		}},
		{name: "remove chat verification invalid chat", call: func() error {
			_, err := bot.RemoveChatVerification(context.Background(), RemoveChatVerificationParams{})
			return err
		}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.call()
			if err == nil {
				t.Fatal("expected error")
			}
			assertNoToken(t, err, token)
		})
	}
}

func TestVerificationStatusMethodsReturnAPIError(t *testing.T) {
	const token = "123:secret"
	tests := verificationStatusErrorCases()
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
			_, err := tt.call(context.Background(), bot)
			if err == nil {
				t.Fatal("expected error")
			}
			var apiErr *apierrors.APIError
			if !stderrors.As(err, &apiErr) {
				t.Fatalf("expected APIError, got %T", err)
			}
			assertNoToken(t, err, token)
			assertNoSensitiveEmojiID(t, err)
		})
	}
}

func TestVerificationStatusMethodsResponseAndContextErrors(t *testing.T) {
	const token = "123:secret"
	tests := verificationStatusErrorCases()
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
			assertNoSensitiveEmojiID(t, err)
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
			assertNoSensitiveEmojiID(t, err)
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
			assertNoSensitiveEmojiID(t, err)
		})
	}
}

type verificationStatusErrorCase struct {
	name   string
	method string
	call   func(context.Context, *Bot) (bool, error)
}

func verificationStatusErrorCases() []verificationStatusErrorCase {
	return []verificationStatusErrorCase{
		{name: "set user emoji status", method: "setUserEmojiStatus", call: func(ctx context.Context, bot *Bot) (bool, error) {
			return bot.SetUserEmojiStatus(ctx, SetUserEmojiStatusParams{UserID: 123, EmojiStatusCustomEmojiID: "sensitive-custom-emoji-id"})
		}},
		{name: "verify user", method: "verifyUser", call: func(ctx context.Context, bot *Bot) (bool, error) {
			return bot.VerifyUser(ctx, VerifyUserParams{UserID: 123, CustomDescription: "Test organization"})
		}},
		{name: "verify chat", method: "verifyChat", call: func(ctx context.Context, bot *Bot) (bool, error) {
			return bot.VerifyChat(ctx, VerifyChatParams{ChatID: ChatIDInt(-100123), CustomDescription: "Test channel"})
		}},
		{name: "remove user verification", method: "removeUserVerification", call: func(ctx context.Context, bot *Bot) (bool, error) {
			return bot.RemoveUserVerification(ctx, RemoveUserVerificationParams{UserID: 123})
		}},
		{name: "remove chat verification", method: "removeChatVerification", call: func(ctx context.Context, bot *Bot) (bool, error) {
			return bot.RemoveChatVerification(ctx, RemoveChatVerificationParams{ChatID: ChatIDInt(-100123)})
		}},
	}
}

func assertNoSensitiveEmojiID(t *testing.T, err error) {
	t.Helper()
	if strings.Contains(err.Error(), "sensitive-custom-emoji-id") {
		t.Fatalf("error leaked custom emoji ID: %q", err.Error())
	}
}
