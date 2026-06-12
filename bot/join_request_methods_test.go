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

func TestApproveChatJoinRequestSendsPayloadAndDecodesResult(t *testing.T) {
	const token = "123:secret"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/bot"+token+"/approveChatJoinRequest" {
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
		_, _ = w.Write([]byte(`{"ok":true,"result":true}`))
	}))
	defer server.Close()

	bot := newTestBot(t, token, server.URL, server.Client())
	ok, err := bot.ApproveChatJoinRequest(context.Background(), ApproveChatJoinRequestParams{ChatID: ChatIDInt(12345), UserID: 777})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Fatal("expected true result")
	}
}

func TestDeclineChatJoinRequestSendsPayloadAndDecodesResult(t *testing.T) {
	const token = "123:secret"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/bot"+token+"/declineChatJoinRequest" {
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
		_, _ = w.Write([]byte(`{"ok":true,"result":true}`))
	}))
	defer server.Close()

	bot := newTestBot(t, token, server.URL, server.Client())
	ok, err := bot.DeclineChatJoinRequest(context.Background(), DeclineChatJoinRequestParams{ChatID: ChatIDInt(12345), UserID: 777})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Fatal("expected true result")
	}
}

func TestAnswerChatJoinRequestQuerySendsPayloadAndDecodesResult(t *testing.T) {
	const token = "123:secret"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/bot"+token+"/answerChatJoinRequestQuery" {
			t.Fatalf("unexpected path: %q", r.URL.Path)
		}
		var payload map[string]any
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode payload: %v", err)
		}
		if payload["chat_join_request_query_id"] != "join-query-id" || payload["result"] != "queue" {
			t.Fatalf("unexpected payload: %#v", payload)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true,"result":true}`))
	}))
	defer server.Close()

	bot := newTestBot(t, token, server.URL, server.Client())
	ok, err := bot.AnswerChatJoinRequestQuery(context.Background(), AnswerChatJoinRequestQueryParams{
		ChatJoinRequestQueryID: "join-query-id",
		Result:                 ChatJoinRequestQueryQueue,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Fatal("expected true result")
	}
}

func TestSendChatJoinRequestWebAppSendsPayloadAndDecodesResult(t *testing.T) {
	const token = "123:secret"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/bot"+token+"/sendChatJoinRequestWebApp" {
			t.Fatalf("unexpected path: %q", r.URL.Path)
		}
		var payload map[string]any
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode payload: %v", err)
		}
		if payload["chat_join_request_query_id"] != "join-query-id" || payload["web_app_url"] != "https://example.com/check" {
			t.Fatalf("unexpected payload: %#v", payload)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true,"result":true}`))
	}))
	defer server.Close()

	bot := newTestBot(t, token, server.URL, server.Client())
	ok, err := bot.SendChatJoinRequestWebApp(context.Background(), SendChatJoinRequestWebAppParams{
		ChatJoinRequestQueryID: "join-query-id",
		WebAppURL:              "https://example.com/check",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Fatal("expected true result")
	}
}

func TestChatJoinRequestMethodValidation(t *testing.T) {
	const token = "123:secret"
	bot := newTestBot(t, token, "https://example.test", nil)
	tests := []struct {
		name    string
		approve ApproveChatJoinRequestParams
		decline DeclineChatJoinRequestParams
	}{
		{name: "empty chat", approve: ApproveChatJoinRequestParams{UserID: 1}, decline: DeclineChatJoinRequestParams{UserID: 1}},
		{name: "zero user", approve: ApproveChatJoinRequestParams{ChatID: ChatIDInt(123)}, decline: DeclineChatJoinRequestParams{ChatID: ChatIDInt(123)}},
		{name: "negative user", approve: ApproveChatJoinRequestParams{ChatID: ChatIDInt(123), UserID: -1}, decline: DeclineChatJoinRequestParams{ChatID: ChatIDInt(123), UserID: -1}},
	}
	for _, tt := range tests {
		t.Run("approve "+tt.name, func(t *testing.T) {
			ok, err := bot.ApproveChatJoinRequest(context.Background(), tt.approve)
			if err == nil {
				t.Fatal("expected error")
			}
			if ok {
				t.Fatal("expected false result")
			}
			assertNoToken(t, err, token)
		})
		t.Run("decline "+tt.name, func(t *testing.T) {
			ok, err := bot.DeclineChatJoinRequest(context.Background(), tt.decline)
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

func TestChatJoinRequestQueryValidation(t *testing.T) {
	const token = "123:secret"
	bot := newTestBot(t, token, "https://example.test", nil)
	answerTests := []AnswerChatJoinRequestQueryParams{
		{Result: ChatJoinRequestQueryApprove},
		{ChatJoinRequestQueryID: "query", Result: "maybe"},
	}
	for _, params := range answerTests {
		ok, err := bot.AnswerChatJoinRequestQuery(context.Background(), params)
		if err == nil {
			t.Fatal("expected answer error")
		}
		if ok {
			t.Fatal("expected false answer result")
		}
		assertNoToken(t, err, token)
	}

	webAppTests := []SendChatJoinRequestWebAppParams{
		{WebAppURL: "https://example.com/check"},
		{ChatJoinRequestQueryID: "query"},
		{ChatJoinRequestQueryID: "query", WebAppURL: "ftp://example.com/check"},
	}
	for _, params := range webAppTests {
		ok, err := bot.SendChatJoinRequestWebApp(context.Background(), params)
		if err == nil {
			t.Fatal("expected web app error")
		}
		if ok {
			t.Fatal("expected false web app result")
		}
		assertNoToken(t, err, token)
	}
}

func TestChatJoinRequestMethodsReturnAPIError(t *testing.T) {
	const token = "123:secret"
	tests := []struct {
		name   string
		method string
		call   func(*Bot) (bool, error)
	}{
		{name: "approve", method: "approveChatJoinRequest", call: func(bot *Bot) (bool, error) {
			return bot.ApproveChatJoinRequest(context.Background(), ApproveChatJoinRequestParams{ChatID: ChatIDInt(123), UserID: 1})
		}},
		{name: "decline", method: "declineChatJoinRequest", call: func(bot *Bot) (bool, error) {
			return bot.DeclineChatJoinRequest(context.Background(), DeclineChatJoinRequestParams{ChatID: ChatIDInt(123), UserID: 1})
		}},
		{name: "answer query", method: "answerChatJoinRequestQuery", call: func(bot *Bot) (bool, error) {
			return bot.AnswerChatJoinRequestQuery(context.Background(), AnswerChatJoinRequestQueryParams{ChatJoinRequestQueryID: "query", Result: ChatJoinRequestQueryApprove})
		}},
		{name: "send web app", method: "sendChatJoinRequestWebApp", call: func(bot *Bot) (bool, error) {
			return bot.SendChatJoinRequestWebApp(context.Background(), SendChatJoinRequestWebAppParams{ChatJoinRequestQueryID: "query", WebAppURL: "https://example.com/check"})
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

func TestChatJoinRequestMethodsResponseAndContextErrors(t *testing.T) {
	const token = "123:secret"
	tests := []struct {
		name   string
		method string
		call   func(context.Context, *Bot) (bool, error)
	}{
		{name: "approve", method: "approveChatJoinRequest", call: func(ctx context.Context, bot *Bot) (bool, error) {
			return bot.ApproveChatJoinRequest(ctx, ApproveChatJoinRequestParams{ChatID: ChatIDInt(123), UserID: 1})
		}},
		{name: "decline", method: "declineChatJoinRequest", call: func(ctx context.Context, bot *Bot) (bool, error) {
			return bot.DeclineChatJoinRequest(ctx, DeclineChatJoinRequestParams{ChatID: ChatIDInt(123), UserID: 1})
		}},
		{name: "answer query", method: "answerChatJoinRequestQuery", call: func(ctx context.Context, bot *Bot) (bool, error) {
			return bot.AnswerChatJoinRequestQuery(ctx, AnswerChatJoinRequestQueryParams{ChatJoinRequestQueryID: "query", Result: ChatJoinRequestQueryApprove})
		}},
		{name: "send web app", method: "sendChatJoinRequestWebApp", call: func(ctx context.Context, bot *Bot) (bool, error) {
			return bot.SendChatJoinRequestWebApp(ctx, SendChatJoinRequestWebAppParams{ChatJoinRequestQueryID: "query", WebAppURL: "https://example.com/check"})
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
