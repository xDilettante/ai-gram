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
	"github.com/xDilettante/ai-gram/telegram"
)

func TestSavePreparedKeyboardButtonSendsPayloadAndDecodesResult(t *testing.T) {
	testBotProfileObjectSuccess(t, "savePreparedKeyboardButton", `{"id":"prepared-button-id"}`, func(bot *Bot) (any, error) {
		return bot.SavePreparedKeyboardButton(context.Background(), SavePreparedKeyboardButtonParams{
			UserID: 123,
			Button: telegram.KeyboardButtonManagedBot("Create bot", telegram.KeyboardButtonRequestManagedBot{RequestID: 42, SuggestedName: "Test Bot", SuggestedUsername: "test_bot"}),
		})
	}, func(t *testing.T, payload map[string]any) {
		if payload["user_id"] != float64(123) {
			t.Fatalf("unexpected payload: %#v", payload)
		}
		button, ok := payload["button"].(map[string]any)
		if !ok {
			t.Fatalf("missing button payload: %#v", payload)
		}
		request := button["request_managed_bot"].(map[string]any)
		if button["text"] != "Create bot" || request["request_id"] != float64(42) || request["suggested_name"] != "Test Bot" || request["suggested_username"] != "test_bot" {
			t.Fatalf("unexpected button payload: %#v", button)
		}
	}, func(t *testing.T, result any) {
		button := result.(*telegram.PreparedKeyboardButton)
		if button.ID != "prepared-button-id" {
			t.Fatalf("unexpected result: %+v", button)
		}
	})
}

func TestSavePreparedKeyboardButtonValidation(t *testing.T) {
	const token = "123:secret"
	bot := newTestBot(t, token, "https://example.test", nil)
	tests := []struct {
		name   string
		params SavePreparedKeyboardButtonParams
	}{
		{name: "invalid user", params: SavePreparedKeyboardButtonParams{UserID: 0, Button: telegram.KeyboardButtonManagedBot("Create", telegram.KeyboardButtonRequestManagedBot{RequestID: 1})}},
		{name: "empty button", params: SavePreparedKeyboardButtonParams{UserID: 1}},
		{name: "unsupported plain button", params: SavePreparedKeyboardButtonParams{UserID: 1, Button: telegram.KeyboardButtonText("Plain")}},
		{name: "invalid managed bot request", params: SavePreparedKeyboardButtonParams{UserID: 1, Button: telegram.KeyboardButtonManagedBot("Create", telegram.KeyboardButtonRequestManagedBot{})}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := bot.SavePreparedKeyboardButton(context.Background(), tt.params)
			if err == nil {
				t.Fatal("expected error")
			}
			if result != nil {
				t.Fatalf("expected nil result, got %+v", result)
			}
			assertNoToken(t, err, token)
		})
	}
}

func TestSavePreparedKeyboardButtonErrors(t *testing.T) {
	call := func(bot *Bot) (any, error) {
		return bot.SavePreparedKeyboardButton(context.Background(), validSavePreparedKeyboardButtonParams())
	}
	callWithContext := func(bot *Bot, ctx context.Context) (any, error) {
		return bot.SavePreparedKeyboardButton(ctx, validSavePreparedKeyboardButtonParams())
	}
	testBotProfileObjectMethodErrorCases(t, "savePreparedKeyboardButton", call, callWithContext)
}

func TestGetManagedBotTokenSendsPayloadAndDecodesToken(t *testing.T) {
	testManagedBotTokenSuccess(t, "getManagedBotToken", func(bot *Bot) (string, error) {
		return bot.GetManagedBotToken(context.Background(), GetManagedBotTokenParams{UserID: 77})
	})
}

func TestReplaceManagedBotTokenSendsPayloadAndDecodesToken(t *testing.T) {
	testManagedBotTokenSuccess(t, "replaceManagedBotToken", func(bot *Bot) (string, error) {
		return bot.ReplaceManagedBotToken(context.Background(), ReplaceManagedBotTokenParams{UserID: 77})
	})
}

func TestGetManagedBotAccessSettingsSendsPayloadAndDecodesResult(t *testing.T) {
	testBotProfileObjectSuccess(t, "getManagedBotAccessSettings", `{"is_access_restricted":true,"added_users":[{"id":7,"is_bot":false,"first_name":"Alice"}]}`, func(bot *Bot) (any, error) {
		return bot.GetManagedBotAccessSettings(context.Background(), GetManagedBotAccessSettingsParams{UserID: 77})
	}, func(t *testing.T, payload map[string]any) {
		if payload["user_id"] != float64(77) {
			t.Fatalf("unexpected payload: %#v", payload)
		}
	}, func(t *testing.T, result any) {
		settings := result.(*telegram.BotAccessSettings)
		if !settings.IsAccessRestricted || len(settings.AddedUsers) != 1 || settings.AddedUsers[0].ID != 7 {
			t.Fatalf("unexpected result: %+v", settings)
		}
	})
}

func TestSetManagedBotAccessSettingsSendsPayloadAndDecodesResult(t *testing.T) {
	testChatManagementSimpleBoolSuccess(t, "setManagedBotAccessSettings", func(bot *Bot) (bool, error) {
		return bot.SetManagedBotAccessSettings(context.Background(), SetManagedBotAccessSettingsParams{
			UserID:             77,
			IsAccessRestricted: true,
			AddedUserIDs:       []int64{7, 8},
		})
	}, func(t *testing.T, payload map[string]any) {
		if payload["user_id"] != float64(77) || payload["is_access_restricted"] != true {
			t.Fatalf("unexpected payload: %#v", payload)
		}
		userIDs, ok := payload["added_user_ids"].([]any)
		if !ok || len(userIDs) != 2 || userIDs[0] != float64(7) || userIDs[1] != float64(8) {
			t.Fatalf("unexpected added_user_ids payload: %#v", payload["added_user_ids"])
		}
	})
}

func TestSetManagedBotAccessSettingsIncludesFalseRestrictionFlag(t *testing.T) {
	testChatManagementSimpleBoolSuccess(t, "setManagedBotAccessSettings", func(bot *Bot) (bool, error) {
		return bot.SetManagedBotAccessSettings(context.Background(), SetManagedBotAccessSettingsParams{UserID: 77})
	}, func(t *testing.T, payload map[string]any) {
		if payload["user_id"] != float64(77) {
			t.Fatalf("unexpected payload: %#v", payload)
		}
		restricted, ok := payload["is_access_restricted"].(bool)
		if !ok || restricted {
			t.Fatalf("expected explicit false restriction flag, got %#v", payload)
		}
		if _, ok := payload["added_user_ids"]; ok {
			t.Fatalf("unexpected added_user_ids payload: %#v", payload)
		}
	})
}

func TestGetUserPersonalChatMessagesSendsPayloadAndDecodesMessages(t *testing.T) {
	testBotProfileObjectSuccess(t, "getUserPersonalChatMessages", `[{"message_id":10,"chat":{"id":123,"type":"private"},"date":111,"text":"hello"}]`, func(bot *Bot) (any, error) {
		return bot.GetUserPersonalChatMessages(context.Background(), GetUserPersonalChatMessagesParams{UserID: 77, Limit: 2})
	}, func(t *testing.T, payload map[string]any) {
		if payload["user_id"] != float64(77) || payload["limit"] != float64(2) {
			t.Fatalf("unexpected payload: %#v", payload)
		}
	}, func(t *testing.T, result any) {
		messages := result.([]telegram.Message)
		if len(messages) != 1 || messages[0].MessageID != 10 || messages[0].Text != "hello" {
			t.Fatalf("unexpected messages: %+v", messages)
		}
	})
}

func TestManagedBotTokenValidation(t *testing.T) {
	const token = "123:secret"
	bot := newTestBot(t, token, "https://example.test", nil)
	for _, tt := range []struct {
		name string
		call func() (string, error)
	}{
		{name: "get zero", call: func() (string, error) {
			return bot.GetManagedBotToken(context.Background(), GetManagedBotTokenParams{})
		}},
		{name: "get negative", call: func() (string, error) {
			return bot.GetManagedBotToken(context.Background(), GetManagedBotTokenParams{UserID: -1})
		}},
		{name: "replace zero", call: func() (string, error) {
			return bot.ReplaceManagedBotToken(context.Background(), ReplaceManagedBotTokenParams{})
		}},
		{name: "replace negative", call: func() (string, error) {
			return bot.ReplaceManagedBotToken(context.Background(), ReplaceManagedBotTokenParams{UserID: -1})
		}},
	} {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.call()
			if err == nil {
				t.Fatal("expected error")
			}
			if result != "" {
				t.Fatal("expected empty token on error")
			}
			assertNoToken(t, err, token)
		})
	}
}

func TestManagedBotAccessValidation(t *testing.T) {
	const token = "123:secret"
	bot := newTestBot(t, token, "https://example.test", nil)
	tests := []struct {
		name string
		call func() (any, error)
	}{
		{name: "get access zero user", call: func() (any, error) {
			return bot.GetManagedBotAccessSettings(context.Background(), GetManagedBotAccessSettingsParams{})
		}},
		{name: "set access zero user", call: func() (any, error) {
			return bot.SetManagedBotAccessSettings(context.Background(), SetManagedBotAccessSettingsParams{})
		}},
		{name: "set access too many added users", call: func() (any, error) {
			return bot.SetManagedBotAccessSettings(context.Background(), SetManagedBotAccessSettingsParams{UserID: 77, AddedUserIDs: makeInt64Range(11)})
		}},
		{name: "set access invalid added user", call: func() (any, error) {
			return bot.SetManagedBotAccessSettings(context.Background(), SetManagedBotAccessSettingsParams{UserID: 77, AddedUserIDs: []int64{7, 0}})
		}},
		{name: "personal messages zero user", call: func() (any, error) {
			return bot.GetUserPersonalChatMessages(context.Background(), GetUserPersonalChatMessagesParams{Limit: 1})
		}},
		{name: "personal messages missing limit", call: func() (any, error) {
			return bot.GetUserPersonalChatMessages(context.Background(), GetUserPersonalChatMessagesParams{UserID: 77})
		}},
		{name: "personal messages excessive limit", call: func() (any, error) {
			return bot.GetUserPersonalChatMessages(context.Background(), GetUserPersonalChatMessagesParams{UserID: 77, Limit: 21})
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.call()
			if err == nil {
				t.Fatal("expected error")
			}
			assertManagedBotAccessZeroResult(t, result)
			assertNoToken(t, err, token)
		})
	}
}

func TestManagedBotTokenMethodErrors(t *testing.T) {
	for _, tt := range []struct {
		name            string
		method          string
		call            func(*Bot) (string, error)
		callWithContext func(*Bot, context.Context) (string, error)
	}{
		{name: "get", method: "getManagedBotToken", call: func(bot *Bot) (string, error) {
			return bot.GetManagedBotToken(context.Background(), GetManagedBotTokenParams{UserID: 77})
		}, callWithContext: func(bot *Bot, ctx context.Context) (string, error) {
			return bot.GetManagedBotToken(ctx, GetManagedBotTokenParams{UserID: 77})
		}},
		{name: "replace", method: "replaceManagedBotToken", call: func(bot *Bot) (string, error) {
			return bot.ReplaceManagedBotToken(context.Background(), ReplaceManagedBotTokenParams{UserID: 77})
		}, callWithContext: func(bot *Bot, ctx context.Context) (string, error) {
			return bot.ReplaceManagedBotToken(ctx, ReplaceManagedBotTokenParams{UserID: 77})
		}},
	} {
		t.Run(tt.name, func(t *testing.T) {
			testManagedBotTokenMethodErrorCases(t, tt.method, tt.call, tt.callWithContext)
		})
	}
}

func TestManagedBotAccessMethodErrors(t *testing.T) {
	t.Run("get access settings", func(t *testing.T) {
		testBotProfileObjectMethodErrorCases(t, "getManagedBotAccessSettings", func(bot *Bot) (any, error) {
			return bot.GetManagedBotAccessSettings(context.Background(), GetManagedBotAccessSettingsParams{UserID: 77})
		}, func(bot *Bot, ctx context.Context) (any, error) {
			return bot.GetManagedBotAccessSettings(ctx, GetManagedBotAccessSettingsParams{UserID: 77})
		})
	})

	t.Run("set access settings", func(t *testing.T) {
		testBoolMethodErrorCases(t, "setManagedBotAccessSettings", func(bot *Bot) (bool, error) {
			return bot.SetManagedBotAccessSettings(context.Background(), SetManagedBotAccessSettingsParams{UserID: 77, IsAccessRestricted: true})
		}, func(bot *Bot, ctx context.Context) (bool, error) {
			return bot.SetManagedBotAccessSettings(ctx, SetManagedBotAccessSettingsParams{UserID: 77, IsAccessRestricted: true})
		})
	})

	t.Run("personal chat messages", func(t *testing.T) {
		testBotProfileObjectMethodErrorCases(t, "getUserPersonalChatMessages", func(bot *Bot) (any, error) {
			return bot.GetUserPersonalChatMessages(context.Background(), GetUserPersonalChatMessagesParams{UserID: 77, Limit: 2})
		}, func(bot *Bot, ctx context.Context) (any, error) {
			return bot.GetUserPersonalChatMessages(ctx, GetUserPersonalChatMessagesParams{UserID: 77, Limit: 2})
		})
	})
}

func validSavePreparedKeyboardButtonParams() SavePreparedKeyboardButtonParams {
	return SavePreparedKeyboardButtonParams{
		UserID: 123,
		Button: telegram.KeyboardButtonManagedBot("Create", telegram.KeyboardButtonRequestManagedBot{RequestID: 1}),
	}
}

func makeInt64Range(count int) []int64 {
	values := make([]int64, count)
	for index := range values {
		values[index] = int64(index + 1)
	}
	return values
}

func assertManagedBotAccessZeroResult(t *testing.T, result any) {
	t.Helper()
	switch value := result.(type) {
	case nil:
	case bool:
		if value {
			t.Fatal("expected false result")
		}
	case []telegram.Message:
		if value != nil {
			t.Fatalf("expected nil messages, got %+v", value)
		}
	case *telegram.BotAccessSettings:
		if value != nil {
			t.Fatalf("expected nil settings, got %+v", value)
		}
	default:
		t.Fatalf("unexpected result type %T", result)
	}
}

func testManagedBotTokenSuccess(t *testing.T, method string, call func(*Bot) (string, error)) {
	t.Helper()
	const token = "123:secret"
	const returnedToken = "managed-token-redacted-for-test"
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
		if payload["user_id"] != float64(77) {
			t.Fatalf("unexpected payload: %#v", payload)
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{"ok": true, "result": returnedToken})
	}))
	defer server.Close()

	bot := newTestBot(t, token, server.URL, server.Client())
	result, err := call(bot)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != returnedToken {
		t.Fatal("unexpected returned token")
	}
}

func testManagedBotTokenMethodErrorCases(t *testing.T, method string, call func(*Bot) (string, error), callWithContext func(*Bot, context.Context) (string, error)) {
	t.Helper()
	const token = "123:secret"
	const returnedToken = "managed-token-redacted-for-test"

	t.Run("api error", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/bot"+token+"/"+method {
				t.Fatalf("unexpected path: %q", r.URL.Path)
			}
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"ok":false,"error_code":400,"description":"Bad Request 123:secret"}`))
		}))
		defer server.Close()

		bot := newTestBot(t, token, server.URL, server.Client())
		result, err := call(bot)
		if err == nil {
			t.Fatal("expected error")
		}
		if result != "" {
			t.Fatal("expected empty token on error")
		}
		var apiErr *apierrors.APIError
		if !stderrors.As(err, &apiErr) {
			t.Fatalf("expected APIError, got %T", err)
		}
		assertNoToken(t, err, token)
		if strings.Contains(err.Error(), returnedToken) {
			t.Fatal("error leaked managed bot token")
		}
	})

	t.Run("invalid json", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`not-json`))
		}))
		defer server.Close()

		bot := newTestBot(t, token, server.URL, server.Client())
		result, err := call(bot)
		if err == nil {
			t.Fatal("expected error")
		}
		if result != "" {
			t.Fatal("expected empty token on error")
		}
		assertNoToken(t, err, token)
	})

	t.Run("http status", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "server error", http.StatusInternalServerError)
		}))
		defer server.Close()

		bot := newTestBot(t, token, server.URL, server.Client())
		result, err := call(bot)
		if err == nil {
			t.Fatal("expected error")
		}
		if result != "" {
			t.Fatal("expected empty token on error")
		}
		assertNoToken(t, err, token)
	})

	t.Run("cancelled context", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			t.Fatal("request should not reach server")
		}))
		defer server.Close()

		bot := newTestBot(t, token, server.URL, server.Client())
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		result, err := callWithContext(bot, ctx)
		if err == nil {
			t.Fatal("expected error")
		}
		if result != "" {
			t.Fatal("expected empty token on error")
		}
		assertNoToken(t, err, token)
	})
}
