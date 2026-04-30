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

func TestSendGameSendsPayloadAndDecodesMessage(t *testing.T) {
	const token = "123:secret"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/bot"+token+"/sendGame" {
			t.Fatalf("unexpected path: %q", r.URL.Path)
		}
		var payload map[string]any
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode payload: %v", err)
		}
		if payload["business_connection_id"] != "bc-123" || payload["chat_id"] != float64(12345) || payload["message_thread_id"] != float64(99) || payload["game_short_name"] != "chess" {
			t.Fatalf("unexpected payload: %#v", payload)
		}
		if payload["allow_paid_broadcast"] != true || payload["message_effect_id"] != "effect" || payload["disable_notification"] != true || payload["protect_content"] != true {
			t.Fatalf("missing optional flags: %#v", payload)
		}
		if payload["reply_parameters"].(map[string]any)["message_id"] != float64(7) {
			t.Fatalf("unexpected reply_parameters: %#v", payload)
		}
		keyboard := payload["reply_markup"].(map[string]any)["inline_keyboard"].([]any)
		firstRow := keyboard[0].([]any)
		if _, ok := firstRow[0].(map[string]any)["callback_game"]; !ok {
			t.Fatalf("callback_game missing: %#v", payload)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true,"result":{"message_id":77,"chat":{"id":12345,"type":"private"},"date":100,"game":{"title":"Chess","description":"A game","photo":[{"file_id":"photo","file_unique_id":"unique","width":1,"height":1}],"text":"Scores","text_entities":[{"type":"bold","offset":0,"length":6}],"animation":{"file_id":"anim","file_unique_id":"anim-unique","width":2,"height":2,"duration":3}}}}`))
	}))
	defer server.Close()

	bot := newTestBot(t, token, server.URL, server.Client())
	markup := telegram.NewInlineKeyboard([]telegram.InlineKeyboardButton{telegram.InlineButtonGame("Play")})
	message, err := bot.SendGame(context.Background(), SendGameParams{
		BusinessConnectionID: "bc-123",
		ChatID:               ChatIDInt(12345),
		MessageThreadID:      99,
		GameShortName:        "chess",
		DisableNotification:  true,
		ProtectContent:       true,
		AllowPaidBroadcast:   true,
		MessageEffectID:      "effect",
		ReplyParameters:      &telegram.ReplyParameters{MessageID: 7},
		ReplyMarkup:          &markup,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if message == nil || message.Game == nil || message.Game.Title != "Chess" || message.Game.Animation == nil {
		t.Fatalf("unexpected message: %+v", message)
	}
}

func TestSendGameValidation(t *testing.T) {
	const token = "123:secret"
	bot := newTestBot(t, token, "https://example.test", nil)
	invalidMarkup := telegram.InlineKeyboardMarkup{}
	tests := []struct {
		name   string
		params SendGameParams
	}{
		{name: "invalid chat", params: SendGameParams{GameShortName: "game"}},
		{name: "empty game", params: SendGameParams{ChatID: ChatIDInt(1)}},
		{name: "blank game", params: SendGameParams{ChatID: ChatIDInt(1), GameShortName: "   "}},
		{name: "negative thread", params: SendGameParams{ChatID: ChatIDInt(1), GameShortName: "game", MessageThreadID: -1}},
		{name: "invalid reply parameters", params: SendGameParams{ChatID: ChatIDInt(1), GameShortName: "game", ReplyParameters: &telegram.ReplyParameters{}}},
		{name: "invalid reply markup", params: SendGameParams{ChatID: ChatIDInt(1), GameShortName: "game", ReplyMarkup: &invalidMarkup}},
		{name: "reply markup missing callback game", params: SendGameParams{ChatID: ChatIDInt(1), GameShortName: "game", ReplyMarkup: inlineMarkupPtr(telegram.NewInlineKeyboard([]telegram.InlineKeyboardButton{telegram.InlineButtonCallback("Other", "other")}))}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			message, err := bot.SendGame(context.Background(), tt.params)
			if err == nil {
				t.Fatal("expected error")
			}
			if message != nil {
				t.Fatalf("expected nil message, got %+v", message)
			}
			assertNoToken(t, err, token)
		})
	}
}

func TestSetGameScoreChatTargetReturnsMessage(t *testing.T) {
	const token = "123:secret"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/bot"+token+"/setGameScore" {
			t.Fatalf("unexpected path: %q", r.URL.Path)
		}
		var payload map[string]any
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode payload: %v", err)
		}
		if payload["user_id"] != float64(42) || payload["score"] != float64(9001) || payload["chat_id"] != float64(12345) || payload["message_id"] != float64(77) {
			t.Fatalf("unexpected payload: %#v", payload)
		}
		if payload["force"] != true || payload["disable_edit_message"] != true {
			t.Fatalf("unexpected flags: %#v", payload)
		}
		if _, ok := payload["inline_message_id"]; ok {
			t.Fatalf("inline target should be omitted: %#v", payload)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true,"result":{"message_id":77,"chat":{"id":12345,"type":"private"},"date":100,"game":{"title":"Chess","description":"A game","photo":[{"file_id":"photo","file_unique_id":"unique","width":1,"height":1}]}}}`))
	}))
	defer server.Close()

	bot := newTestBot(t, token, server.URL, server.Client())
	result, err := bot.SetGameScore(context.Background(), SetGameScoreParams{UserID: 42, Score: 9001, Force: true, DisableEditMessage: true, ChatID: ChatIDInt(12345), MessageID: 77})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result == nil || !result.IsOK() || !result.IsMessage() || result.Message.Game == nil {
		t.Fatalf("unexpected result: %+v", result)
	}
}

func TestSetGameScoreInlineTargetReturnsTrue(t *testing.T) {
	const token = "123:secret"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/bot"+token+"/setGameScore" {
			t.Fatalf("unexpected path: %q", r.URL.Path)
		}
		var payload map[string]any
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode payload: %v", err)
		}
		if payload["inline_message_id"] != "inline-id" || payload["user_id"] != float64(42) || payload["score"] != float64(10) {
			t.Fatalf("unexpected payload: %#v", payload)
		}
		if _, ok := payload["chat_id"]; ok {
			t.Fatalf("chat_id should be omitted: %#v", payload)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true,"result":true}`))
	}))
	defer server.Close()

	bot := newTestBot(t, token, server.URL, server.Client())
	result, err := bot.SetGameScore(context.Background(), SetGameScoreParams{UserID: 42, Score: 10, InlineMessageID: "inline-id"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result == nil || !result.IsOK() || result.IsMessage() {
		t.Fatalf("unexpected result: %+v", result)
	}
}

func TestSetGameScoreValidation(t *testing.T) {
	const token = "123:secret"
	bot := newTestBot(t, token, "https://example.test", nil)
	tests := []struct {
		name   string
		params SetGameScoreParams
	}{
		{name: "invalid user", params: SetGameScoreParams{UserID: 0, Score: 1, InlineMessageID: "inline"}},
		{name: "negative score", params: SetGameScoreParams{UserID: 1, Score: -1, InlineMessageID: "inline"}},
		{name: "missing target", params: SetGameScoreParams{UserID: 1, Score: 1}},
		{name: "ambiguous target", params: SetGameScoreParams{UserID: 1, Score: 1, ChatID: ChatIDInt(1), MessageID: 1, InlineMessageID: "inline"}},
		{name: "chat without message", params: SetGameScoreParams{UserID: 1, Score: 1, ChatID: ChatIDInt(1)}},
		{name: "message without chat", params: SetGameScoreParams{UserID: 1, Score: 1, MessageID: 1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := bot.SetGameScore(context.Background(), tt.params)
			if err == nil {
				t.Fatal("expected error")
			}
			if !isNilResult(result) {
				t.Fatalf("expected nil result, got %+v", result)
			}
			assertNoToken(t, err, token)
		})
	}
}

func TestGetGameHighScoresChatTarget(t *testing.T) {
	const token = "123:secret"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/bot"+token+"/getGameHighScores" {
			t.Fatalf("unexpected path: %q", r.URL.Path)
		}
		var payload map[string]any
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode payload: %v", err)
		}
		if payload["user_id"] != float64(42) || payload["chat_id"] != float64(12345) || payload["message_id"] != float64(77) {
			t.Fatalf("unexpected payload: %#v", payload)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true,"result":[{"position":1,"user":{"id":42,"is_bot":false,"first_name":"Ada"},"score":9001}]}`))
	}))
	defer server.Close()

	bot := newTestBot(t, token, server.URL, server.Client())
	scores, err := bot.GetGameHighScores(context.Background(), GetGameHighScoresParams{UserID: 42, ChatID: ChatIDInt(12345), MessageID: 77})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(scores) != 1 || scores[0].Position != 1 || scores[0].User.ID != 42 || scores[0].Score != 9001 {
		t.Fatalf("unexpected scores: %+v", scores)
	}
}

func TestGetGameHighScoresInlineTarget(t *testing.T) {
	const token = "123:secret"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/bot"+token+"/getGameHighScores" {
			t.Fatalf("unexpected path: %q", r.URL.Path)
		}
		var payload map[string]any
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode payload: %v", err)
		}
		if payload["inline_message_id"] != "inline-id" || payload["user_id"] != float64(42) {
			t.Fatalf("unexpected payload: %#v", payload)
		}
		if _, ok := payload["chat_id"]; ok {
			t.Fatalf("chat_id should be omitted: %#v", payload)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true,"result":[]}`))
	}))
	defer server.Close()

	bot := newTestBot(t, token, server.URL, server.Client())
	scores, err := bot.GetGameHighScores(context.Background(), GetGameHighScoresParams{UserID: 42, InlineMessageID: "inline-id"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(scores) != 0 {
		t.Fatalf("unexpected scores: %+v", scores)
	}
}

func TestGetGameHighScoresValidation(t *testing.T) {
	const token = "123:secret"
	bot := newTestBot(t, token, "https://example.test", nil)
	tests := []struct {
		name   string
		params GetGameHighScoresParams
	}{
		{name: "invalid user", params: GetGameHighScoresParams{UserID: 0, InlineMessageID: "inline"}},
		{name: "missing target", params: GetGameHighScoresParams{UserID: 1}},
		{name: "ambiguous target", params: GetGameHighScoresParams{UserID: 1, ChatID: ChatIDInt(1), MessageID: 1, InlineMessageID: "inline"}},
		{name: "chat without message", params: GetGameHighScoresParams{UserID: 1, ChatID: ChatIDInt(1)}},
		{name: "message without chat", params: GetGameHighScoresParams{UserID: 1, MessageID: 1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scores, err := bot.GetGameHighScores(context.Background(), tt.params)
			if err == nil {
				t.Fatal("expected error")
			}
			if scores != nil {
				t.Fatalf("expected nil scores, got %+v", scores)
			}
			assertNoToken(t, err, token)
		})
	}
}

func TestGameMethodErrorCases(t *testing.T) {
	const token = "123:secret"
	tests := []struct {
		name string
		call func(*Bot) (any, error)
	}{
		{name: "send game", call: func(bot *Bot) (any, error) { return bot.SendGame(context.Background(), validSendGameParams()) }},
		{name: "set score", call: func(bot *Bot) (any, error) { return bot.SetGameScore(context.Background(), validSetGameScoreParams()) }},
		{name: "get scores", call: func(bot *Bot) (any, error) {
			return bot.GetGameHighScores(context.Background(), validGetGameHighScoresParams())
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name+" api error", func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(`{"ok":false,"error_code":400,"description":"Bad Request"}`))
			}))
			defer server.Close()
			bot := newTestBot(t, token, server.URL, server.Client())
			result, err := tt.call(bot)
			if err == nil {
				t.Fatal("expected error")
			}
			if !isNilResult(result) {
				t.Fatalf("expected nil result, got %+v", result)
			}
			var apiErr *apierrors.APIError
			if !stderrors.As(err, &apiErr) {
				t.Fatalf("expected APIError, got %T", err)
			}
			assertNoToken(t, err, token)
		})
		t.Run(tt.name+" invalid json", func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(`not-json`))
			}))
			defer server.Close()
			bot := newTestBot(t, token, server.URL, server.Client())
			result, err := tt.call(bot)
			if err == nil {
				t.Fatal("expected error")
			}
			if !isNilResult(result) {
				t.Fatalf("expected nil result, got %+v", result)
			}
			assertNoToken(t, err, token)
		})
		t.Run(tt.name+" http status", func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				http.Error(w, "server error", http.StatusInternalServerError)
			}))
			defer server.Close()
			bot := newTestBot(t, token, server.URL, server.Client())
			result, err := tt.call(bot)
			if err == nil {
				t.Fatal("expected error")
			}
			if !isNilResult(result) {
				t.Fatalf("expected nil result, got %+v", result)
			}
			assertNoToken(t, err, token)
		})
		t.Run(tt.name+" cancelled context", func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			cancel()
			bot := newTestBot(t, token, "https://example.test", nil)
			var result any
			var err error
			switch tt.name {
			case "send game":
				result, err = bot.SendGame(ctx, validSendGameParams())
			case "set score":
				result, err = bot.SetGameScore(ctx, validSetGameScoreParams())
			case "get scores":
				result, err = bot.GetGameHighScores(ctx, validGetGameHighScoresParams())
			}
			if err == nil {
				t.Fatal("expected error")
			}
			if !isNilResult(result) {
				t.Fatalf("expected nil result, got %+v", result)
			}
			assertNoToken(t, err, token)
		})
	}
}

func validSendGameParams() SendGameParams {
	return SendGameParams{ChatID: ChatIDInt(12345), GameShortName: "chess"}
}

func validSetGameScoreParams() SetGameScoreParams {
	return SetGameScoreParams{UserID: 42, Score: 10, ChatID: ChatIDInt(12345), MessageID: 77}
}

func validGetGameHighScoresParams() GetGameHighScoresParams {
	return GetGameHighScoresParams{UserID: 42, ChatID: ChatIDInt(12345), MessageID: 77}
}

func inlineMarkupPtr(markup telegram.InlineKeyboardMarkup) *telegram.InlineKeyboardMarkup {
	return &markup
}
