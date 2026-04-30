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

func TestSetMyCommandsSendsDefaultScopeAndDecodesResult(t *testing.T) {
	const token = "123:secret"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/bot"+token+"/setMyCommands" {
			t.Fatalf("unexpected path: %q", r.URL.Path)
		}
		var payload map[string]any
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode payload: %v", err)
		}
		commands := payload["commands"].([]any)
		if len(commands) != 1 {
			t.Fatalf("unexpected commands: %#v", commands)
		}
		command := commands[0].(map[string]any)
		if command["command"] != "start" || command["description"] != "Start the bot" {
			t.Fatalf("unexpected command: %#v", command)
		}
		scope := payload["scope"].(map[string]any)
		if scope["type"] != "default" {
			t.Fatalf("unexpected scope: %#v", scope)
		}
		if payload["language_code"] != "en" {
			t.Fatalf("unexpected language_code: %#v", payload)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true,"result":true}`))
	}))
	defer server.Close()

	bot := newTestBot(t, token, server.URL, server.Client())
	ok, err := bot.SetMyCommands(context.Background(), SetMyCommandsParams{
		Commands:     []telegram.BotCommand{{Command: "start", Description: "Start the bot"}},
		Scope:        telegram.NewBotCommandScopeDefault(),
		LanguageCode: "en",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Fatal("expected true result")
	}
}

func TestSetMyCommandsSendsChatMemberScope(t *testing.T) {
	const token = "123:secret"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/bot"+token+"/setMyCommands" {
			t.Fatalf("unexpected path: %q", r.URL.Path)
		}
		var payload map[string]any
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode payload: %v", err)
		}
		scope := payload["scope"].(map[string]any)
		if scope["type"] != "chat_member" || scope["chat_id"] != "@group" || scope["user_id"] != float64(42) {
			t.Fatalf("unexpected scope: %#v", scope)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true,"result":true}`))
	}))
	defer server.Close()

	bot := newTestBot(t, token, server.URL, server.Client())
	ok, err := bot.SetMyCommands(context.Background(), SetMyCommandsParams{
		Commands: []telegram.BotCommand{{Command: "help", Description: "Show help"}},
		Scope:    telegram.NewBotCommandScopeChatMember("@group", 42),
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Fatal("expected true result")
	}
}

func TestSetMyCommandsValidation(t *testing.T) {
	const token = "123:secret"
	bot := newTestBot(t, token, "https://example.test", nil)
	valid := SetMyCommandsParams{Commands: []telegram.BotCommand{{Command: "start", Description: "Start"}}}
	tests := []struct {
		name   string
		mutate func(*SetMyCommandsParams)
	}{
		{name: "empty commands", mutate: func(p *SetMyCommandsParams) { p.Commands = nil }},
		{name: "empty command", mutate: func(p *SetMyCommandsParams) { p.Commands[0].Command = "" }},
		{name: "empty description", mutate: func(p *SetMyCommandsParams) { p.Commands[0].Description = "" }},
		{name: "invalid scope chat", mutate: func(p *SetMyCommandsParams) { p.Scope = telegram.NewBotCommandScopeChat("") }},
		{name: "invalid scope member", mutate: func(p *SetMyCommandsParams) { p.Scope = telegram.NewBotCommandScopeChatMember("@group", 0) }},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			params := valid
			params.Commands = append([]telegram.BotCommand(nil), valid.Commands...)
			tt.mutate(&params)
			ok, err := bot.SetMyCommands(context.Background(), params)
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

func TestDeleteMyCommandsSendsPayloadAndDecodesResult(t *testing.T) {
	const token = "123:secret"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/bot"+token+"/deleteMyCommands" {
			t.Fatalf("unexpected path: %q", r.URL.Path)
		}
		var payload map[string]any
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode payload: %v", err)
		}
		scope := payload["scope"].(map[string]any)
		if scope["type"] != "all_private_chats" || payload["language_code"] != "en" {
			t.Fatalf("unexpected payload: %#v", payload)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true,"result":true}`))
	}))
	defer server.Close()

	bot := newTestBot(t, token, server.URL, server.Client())
	ok, err := bot.DeleteMyCommands(context.Background(), DeleteMyCommandsParams{Scope: telegram.NewBotCommandScopeAllPrivateChats(), LanguageCode: "en"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Fatal("expected true result")
	}
}

func TestGetMyCommandsSendsPayloadAndDecodesCommands(t *testing.T) {
	const token = "123:secret"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/bot"+token+"/getMyCommands" {
			t.Fatalf("unexpected path: %q", r.URL.Path)
		}
		var payload map[string]any
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode payload: %v", err)
		}
		scope := payload["scope"].(map[string]any)
		if scope["type"] != "chat" || scope["chat_id"] != float64(12345) {
			t.Fatalf("unexpected scope: %#v", scope)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true,"result":[{"command":"start","description":"Start"},{"command":"help","description":"Help"}]}`))
	}))
	defer server.Close()

	bot := newTestBot(t, token, server.URL, server.Client())
	commands, err := bot.GetMyCommands(context.Background(), GetMyCommandsParams{Scope: telegram.NewBotCommandScopeChat(ChatIDInt(12345))})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(commands) != 2 || commands[0].Command != "start" || commands[1].Description != "Help" {
		t.Fatalf("unexpected commands: %+v", commands)
	}
}

func TestCommandsAPIAndTransportErrors(t *testing.T) {
	validSet := SetMyCommandsParams{Commands: []telegram.BotCommand{{Command: "start", Description: "Start"}}}
	testBoolMethodErrorCases(t, "setMyCommands", func(bot *Bot) (bool, error) {
		return bot.SetMyCommands(context.Background(), validSet)
	}, func(bot *Bot, ctx context.Context) (bool, error) {
		return bot.SetMyCommands(ctx, validSet)
	})
	testBotCommandsMethodErrorCases(t, "getMyCommands", func(bot *Bot) ([]telegram.BotCommand, error) {
		return bot.GetMyCommands(context.Background(), GetMyCommandsParams{})
	}, func(bot *Bot, ctx context.Context) ([]telegram.BotCommand, error) {
		return bot.GetMyCommands(ctx, GetMyCommandsParams{})
	})
}

func TestDeleteMyCommandsValidation(t *testing.T) {
	const token = "123:secret"
	bot := newTestBot(t, token, "https://example.test", nil)
	ok, err := bot.DeleteMyCommands(context.Background(), DeleteMyCommandsParams{Scope: telegram.NewBotCommandScopeChatMember("@group", -1)})
	if err == nil {
		t.Fatal("expected error")
	}
	if ok {
		t.Fatal("expected false result")
	}
	assertNoToken(t, err, token)
}

func TestSetChatMenuButtonSendsPayloadAndDecodesResult(t *testing.T) {
	const token = "123:secret"
	tests := []struct {
		name     string
		params   SetChatMenuButtonParams
		wantType string
	}{
		{name: "commands", params: SetChatMenuButtonParams{ChatID: ChatIDInt(12345), MenuButton: telegram.NewMenuButtonCommands()}, wantType: "commands"},
		{name: "default", params: SetChatMenuButtonParams{MenuButton: telegram.NewMenuButtonDefault()}, wantType: "default"},
		{name: "web app", params: SetChatMenuButtonParams{ChatID: ChatIDInt(12345), MenuButton: telegram.NewMenuButtonWebApp("Open", "https://example.com/app")}, wantType: "web_app"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path != "/bot"+token+"/setChatMenuButton" {
					t.Fatalf("unexpected path: %q", r.URL.Path)
				}
				var payload map[string]any
				if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
					t.Fatalf("decode payload: %v", err)
				}
				if tt.params.ChatID.valid() && payload["chat_id"] != float64(12345) {
					t.Fatalf("unexpected chat_id: %#v", payload)
				}
				button := payload["menu_button"].(map[string]any)
				if button["type"] != tt.wantType {
					t.Fatalf("unexpected button: %#v", button)
				}
				if tt.wantType == "web_app" {
					if button["text"] != "Open" {
						t.Fatalf("unexpected web app text: %#v", button)
					}
					webApp := button["web_app"].(map[string]any)
					if webApp["url"] != "https://example.com/app" {
						t.Fatalf("unexpected web app: %#v", webApp)
					}
				}
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(`{"ok":true,"result":true}`))
			}))
			defer server.Close()

			bot := newTestBot(t, token, server.URL, server.Client())
			ok, err := bot.SetChatMenuButton(context.Background(), tt.params)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !ok {
				t.Fatal("expected true result")
			}
		})
	}
}

func TestSetChatMenuButtonValidation(t *testing.T) {
	const token = "123:secret"
	bot := newTestBot(t, token, "https://example.test", nil)
	tests := []struct {
		name   string
		button telegram.MenuButton
	}{
		{name: "empty web app text", button: telegram.NewMenuButtonWebApp("", "https://example.com/app")},
		{name: "empty web app url", button: telegram.NewMenuButtonWebApp("Open", "")},
		{name: "invalid web app url", button: telegram.NewMenuButtonWebApp("Open", "ftp://example.com/app")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ok, err := bot.SetChatMenuButton(context.Background(), SetChatMenuButtonParams{MenuButton: tt.button})
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

func TestGetChatMenuButtonDecodesPolymorphicResult(t *testing.T) {
	const token = "123:secret"
	tests := []struct {
		name     string
		result   string
		assertFn func(*testing.T, telegram.MenuButton)
	}{
		{name: "commands", result: `{"type":"commands"}`, assertFn: func(t *testing.T, button telegram.MenuButton) {
			if _, ok := button.(telegram.MenuButtonCommands); !ok {
				t.Fatalf("expected MenuButtonCommands, got %T", button)
			}
		}},
		{name: "default", result: `{"type":"default"}`, assertFn: func(t *testing.T, button telegram.MenuButton) {
			if _, ok := button.(telegram.MenuButtonDefault); !ok {
				t.Fatalf("expected MenuButtonDefault, got %T", button)
			}
		}},
		{name: "web app", result: `{"type":"web_app","text":"Open","web_app":{"url":"https://example.com/app"}}`, assertFn: func(t *testing.T, button telegram.MenuButton) {
			webApp, ok := button.(telegram.MenuButtonWebApp)
			if !ok {
				t.Fatalf("expected MenuButtonWebApp, got %T", button)
			}
			if webApp.Text != "Open" || webApp.WebApp.URL != "https://example.com/app" {
				t.Fatalf("unexpected web app button: %+v", webApp)
			}
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path != "/bot"+token+"/getChatMenuButton" {
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
				_, _ = w.Write([]byte(`{"ok":true,"result":` + tt.result + `}`))
			}))
			defer server.Close()

			bot := newTestBot(t, token, server.URL, server.Client())
			button, err := bot.GetChatMenuButton(context.Background(), GetChatMenuButtonParams{ChatID: ChatIDInt(12345)})
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			tt.assertFn(t, button)
		})
	}
}

func TestMenuButtonAPIAndTransportErrors(t *testing.T) {
	testBoolMethodErrorCases(t, "setChatMenuButton", func(bot *Bot) (bool, error) {
		return bot.SetChatMenuButton(context.Background(), SetChatMenuButtonParams{MenuButton: telegram.NewMenuButtonCommands()})
	}, func(bot *Bot, ctx context.Context) (bool, error) {
		return bot.SetChatMenuButton(ctx, SetChatMenuButtonParams{MenuButton: telegram.NewMenuButtonCommands()})
	})
	testMenuButtonMethodErrorCases(t, "getChatMenuButton", func(bot *Bot) (telegram.MenuButton, error) {
		return bot.GetChatMenuButton(context.Background(), GetChatMenuButtonParams{})
	}, func(bot *Bot, ctx context.Context) (telegram.MenuButton, error) {
		return bot.GetChatMenuButton(ctx, GetChatMenuButtonParams{})
	})
}

func TestSetMyDefaultAdministratorRightsSendsPayloadAndDecodesResult(t *testing.T) {
	const token = "123:secret"
	tests := []struct {
		name       string
		params     SetMyDefaultAdministratorRightsParams
		wantRights bool
	}{
		{name: "with rights", params: SetMyDefaultAdministratorRightsParams{Rights: &telegram.ChatAdministratorRights{CanManageChat: true, CanDeleteMessages: true, CanPinMessages: true}, ForChannels: true}, wantRights: true},
		{name: "nil rights", params: SetMyDefaultAdministratorRightsParams{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path != "/bot"+token+"/setMyDefaultAdministratorRights" {
					t.Fatalf("unexpected path: %q", r.URL.Path)
				}
				var payload map[string]any
				if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
					t.Fatalf("decode payload: %v", err)
				}
				if tt.wantRights {
					rights := payload["rights"].(map[string]any)
					if rights["can_manage_chat"] != true || rights["can_delete_messages"] != true || rights["can_pin_messages"] != true {
						t.Fatalf("unexpected rights: %#v", rights)
					}
					if payload["for_channels"] != true {
						t.Fatalf("unexpected for_channels: %#v", payload)
					}
				} else if _, ok := payload["rights"]; ok {
					t.Fatalf("rights should be omitted: %#v", payload)
				}
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(`{"ok":true,"result":true}`))
			}))
			defer server.Close()

			bot := newTestBot(t, token, server.URL, server.Client())
			ok, err := bot.SetMyDefaultAdministratorRights(context.Background(), tt.params)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !ok {
				t.Fatal("expected true result")
			}
		})
	}
}

func TestSetMyDefaultAdministratorRightsAPIAndTransportErrors(t *testing.T) {
	testBoolMethodErrorCases(t, "setMyDefaultAdministratorRights", func(bot *Bot) (bool, error) {
		return bot.SetMyDefaultAdministratorRights(context.Background(), SetMyDefaultAdministratorRightsParams{})
	}, func(bot *Bot, ctx context.Context) (bool, error) {
		return bot.SetMyDefaultAdministratorRights(ctx, SetMyDefaultAdministratorRightsParams{})
	})
}

func testBoolMethodErrorCases(t *testing.T, method string, call func(*Bot) (bool, error), callWithContext func(*Bot, context.Context) (bool, error)) {
	t.Helper()
	const token = "123:secret"
	t.Run("api error", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/bot"+token+"/"+method {
				t.Fatalf("unexpected path: %q", r.URL.Path)
			}
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"ok":false,"error_code":400,"description":"Bad Request"}`))
		}))
		defer server.Close()

		bot := newTestBot(t, token, server.URL, server.Client())
		ok, err := call(bot)
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

	t.Run("invalid json", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`not-json`))
		}))
		defer server.Close()

		bot := newTestBot(t, token, server.URL, server.Client())
		_, err := call(bot)
		if err == nil {
			t.Fatal("expected error")
		}
		assertNoToken(t, err, token)
	})

	t.Run("http status", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "server error", http.StatusInternalServerError)
		}))
		defer server.Close()

		bot := newTestBot(t, token, server.URL, server.Client())
		_, err := call(bot)
		if err == nil {
			t.Fatal("expected error")
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
		_, err := callWithContext(bot, ctx)
		if err == nil {
			t.Fatal("expected error")
		}
		assertNoToken(t, err, token)
	})
}

func testBotCommandsMethodErrorCases(t *testing.T, method string, call func(*Bot) ([]telegram.BotCommand, error), callWithContext func(*Bot, context.Context) ([]telegram.BotCommand, error)) {
	t.Helper()
	const token = "123:secret"
	t.Run("api error", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/bot"+token+"/"+method {
				t.Fatalf("unexpected path: %q", r.URL.Path)
			}
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"ok":false,"error_code":400,"description":"Bad Request"}`))
		}))
		defer server.Close()

		bot := newTestBot(t, token, server.URL, server.Client())
		commands, err := call(bot)
		if err == nil {
			t.Fatal("expected error")
		}
		if commands != nil {
			t.Fatalf("expected nil commands, got %+v", commands)
		}
		var apiErr *apierrors.APIError
		if !stderrors.As(err, &apiErr) {
			t.Fatalf("expected APIError, got %T", err)
		}
		assertNoToken(t, err, token)
	})

	t.Run("invalid json", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`not-json`))
		}))
		defer server.Close()

		bot := newTestBot(t, token, server.URL, server.Client())
		_, err := call(bot)
		if err == nil {
			t.Fatal("expected error")
		}
		assertNoToken(t, err, token)
	})

	t.Run("http status", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "server error", http.StatusInternalServerError)
		}))
		defer server.Close()

		bot := newTestBot(t, token, server.URL, server.Client())
		_, err := call(bot)
		if err == nil {
			t.Fatal("expected error")
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
		_, err := callWithContext(bot, ctx)
		if err == nil {
			t.Fatal("expected error")
		}
		assertNoToken(t, err, token)
	})
}

func testMenuButtonMethodErrorCases(t *testing.T, method string, call func(*Bot) (telegram.MenuButton, error), callWithContext func(*Bot, context.Context) (telegram.MenuButton, error)) {
	t.Helper()
	const token = "123:secret"
	t.Run("api error", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/bot"+token+"/"+method {
				t.Fatalf("unexpected path: %q", r.URL.Path)
			}
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"ok":false,"error_code":400,"description":"Bad Request"}`))
		}))
		defer server.Close()

		bot := newTestBot(t, token, server.URL, server.Client())
		button, err := call(bot)
		if err == nil {
			t.Fatal("expected error")
		}
		if button != nil {
			t.Fatalf("expected nil button, got %+v", button)
		}
		var apiErr *apierrors.APIError
		if !stderrors.As(err, &apiErr) {
			t.Fatalf("expected APIError, got %T", err)
		}
		assertNoToken(t, err, token)
	})

	t.Run("invalid json", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`not-json`))
		}))
		defer server.Close()

		bot := newTestBot(t, token, server.URL, server.Client())
		_, err := call(bot)
		if err == nil {
			t.Fatal("expected error")
		}
		assertNoToken(t, err, token)
	})

	t.Run("http status", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "server error", http.StatusInternalServerError)
		}))
		defer server.Close()

		bot := newTestBot(t, token, server.URL, server.Client())
		_, err := call(bot)
		if err == nil {
			t.Fatal("expected error")
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
		_, err := callWithContext(bot, ctx)
		if err == nil {
			t.Fatal("expected error")
		}
		assertNoToken(t, err, token)
	})
}
