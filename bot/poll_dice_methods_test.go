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

func TestSendPollSendsPayloadAndDecodesMessage(t *testing.T) {
	const token = "123:secret"
	isAnonymous := false
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/bot"+token+"/sendPoll" {
			t.Fatalf("unexpected path: %q", r.URL.Path)
		}
		var payload map[string]any
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode payload: %v", err)
		}
		if got := payload["chat_id"]; got != float64(12345) {
			t.Fatalf("unexpected chat_id: %#v", got)
		}
		if got := payload["message_thread_id"]; got != float64(7) {
			t.Fatalf("unexpected message_thread_id: %#v", got)
		}
		if got := payload["question"]; got != "Pick one" {
			t.Fatalf("unexpected question: %#v", got)
		}
		if got := payload["question_parse_mode"]; got != "HTML" {
			t.Fatalf("unexpected question_parse_mode: %#v", got)
		}
		options, ok := payload["options"].([]any)
		if !ok || len(options) != 2 || options[0] != "A" || options[1] != "B" {
			t.Fatalf("unexpected options: %#v", payload["options"])
		}
		if got := payload["is_anonymous"]; got != false {
			t.Fatalf("unexpected is_anonymous: %#v", got)
		}
		if got := payload["type"]; got != "quiz" {
			t.Fatalf("unexpected type: %#v", got)
		}
		if got := payload["allows_multiple_answers"]; got != true {
			t.Fatalf("unexpected allows_multiple_answers: %#v", got)
		}
		if got := payload["allows_revoting"]; got != true {
			t.Fatalf("unexpected allows_revoting: %#v", got)
		}
		if got := payload["shuffle_options"]; got != true {
			t.Fatalf("unexpected shuffle_options: %#v", got)
		}
		if got := payload["allow_adding_options"]; got != true {
			t.Fatalf("unexpected allow_adding_options: %#v", got)
		}
		if got := payload["hide_results_until_closes"]; got != true {
			t.Fatalf("unexpected hide_results_until_closes: %#v", got)
		}
		correctOptionIDs, ok := payload["correct_option_ids"].([]any)
		if !ok || len(correctOptionIDs) != 1 || correctOptionIDs[0] != float64(1) {
			t.Fatalf("unexpected correct_option_ids: %#v", payload["correct_option_ids"])
		}
		if _, ok := payload["correct_option_id"]; ok {
			t.Fatalf("legacy correct_option_id should be omitted when correct_option_ids is used: %#v", payload["correct_option_id"])
		}
		if got := payload["explanation"]; got != "Because" {
			t.Fatalf("unexpected explanation: %#v", got)
		}
		if got := payload["explanation_parse_mode"]; got != "HTML" {
			t.Fatalf("unexpected explanation_parse_mode: %#v", got)
		}
		if got := payload["description"]; got != "Details" {
			t.Fatalf("unexpected description: %#v", got)
		}
		if got := payload["description_parse_mode"]; got != "MarkdownV2" {
			t.Fatalf("unexpected description_parse_mode: %#v", got)
		}
		if got := payload["open_period"]; got != float64(60) {
			t.Fatalf("unexpected open_period: %#v", got)
		}
		if got := payload["close_date"]; got != float64(1234567890) {
			t.Fatalf("unexpected close_date: %#v", got)
		}
		if got := payload["is_closed"]; got != true {
			t.Fatalf("unexpected is_closed: %#v", got)
		}
		if got := payload["disable_notification"]; got != true {
			t.Fatalf("unexpected disable_notification: %#v", got)
		}
		if got := payload["protect_content"]; got != true {
			t.Fatalf("unexpected protect_content: %#v", got)
		}
		reply := payload["reply_parameters"].(map[string]any)
		if got := reply["message_id"]; got != float64(44) {
			t.Fatalf("unexpected reply message_id: %#v", got)
		}
		markup := payload["reply_markup"].(map[string]any)
		if _, ok := markup["inline_keyboard"]; !ok {
			t.Fatalf("reply_markup.inline_keyboard missing: %#v", markup)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true,"result":{"message_id":10,"chat":{"id":12345,"type":"private"},"date":100,"poll":{"id":"poll-id","question":"Pick one","question_entities":[{"type":"custom_emoji","offset":0,"length":1,"custom_emoji_id":"emoji-id"}],"options":[{"persistent_id":"a","text":"A","voter_count":1},{"persistent_id":"b","text":"B","voter_count":0}],"total_voter_count":1,"is_closed":false,"is_anonymous":false,"type":"quiz","allows_multiple_answers":true,"allows_revoting":true,"correct_option_ids":[1],"description":"Details","description_entities":[{"type":"bold","offset":0,"length":7}],"explanation":"Because","open_period":60,"close_date":1234567890}}}`))
	}))
	defer server.Close()

	bot := newTestBot(t, token, server.URL, server.Client())
	message, err := bot.SendPoll(context.Background(), SendPollParams{
		ChatID:                 ChatIDInt(12345),
		MessageThreadID:        7,
		Question:               "Pick one",
		QuestionParseMode:      "HTML",
		Options:                []string{"A", "B"},
		IsAnonymous:            &isAnonymous,
		Type:                   "quiz",
		AllowsMultipleAnswers:  true,
		AllowsRevoting:         true,
		ShuffleOptions:         true,
		AllowAddingOptions:     true,
		HideResultsUntilCloses: true,
		CorrectOptionIDs:       []int{1},
		Explanation:            "Because",
		ExplanationParseMode:   "HTML",
		Description:            "Details",
		DescriptionParseMode:   "MarkdownV2",
		OpenPeriod:             60,
		CloseDate:              1234567890,
		IsClosed:               true,
		DisableNotification:    true,
		ProtectContent:         true,
		ReplyParameters:        &telegram.ReplyParameters{MessageID: 44},
		ReplyMarkup:            telegram.NewInlineKeyboard([]telegram.InlineKeyboardButton{telegram.InlineButtonCallback("OK", "ok")}),
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if message == nil || message.Poll == nil || message.Poll.ID != "poll-id" || len(message.Poll.Options) != 2 || len(message.Poll.CorrectOptionIDs) != 1 || message.Poll.CorrectOptionIDs[0] != 1 {
		t.Fatalf("unexpected message: %+v", message)
	}
	if !message.Poll.AllowsRevoting || message.Poll.Description != "Details" || len(message.Poll.DescriptionEntities) != 1 || message.Poll.Options[0].PersistentID != "a" {
		t.Fatalf("unexpected decoded poll 9.6 fields: %+v", message.Poll)
	}
	if len(message.Poll.QuestionEntities) != 1 || message.Poll.QuestionEntities[0].Type != telegram.EntityCustomEmoji {
		t.Fatalf("unexpected question entities: %+v", message.Poll.QuestionEntities)
	}
}

func TestSendPollStructuredOptions(t *testing.T) {
	const token = "123:secret"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/bot"+token+"/sendPoll" {
			t.Fatalf("unexpected path: %q", r.URL.Path)
		}
		var payload map[string]any
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode payload: %v", err)
		}
		options, ok := payload["options"].([]any)
		if !ok || len(options) != 2 {
			t.Fatalf("unexpected options: %#v", payload["options"])
		}
		first, ok := options[0].(map[string]any)
		if !ok || first["text"] != "A" || first["text_parse_mode"] != "HTML" {
			t.Fatalf("unexpected first option: %#v", options[0])
		}
		second, ok := options[1].(map[string]any)
		if !ok || second["text"] != "B" {
			t.Fatalf("unexpected second option: %#v", options[1])
		}
		entities, ok := second["text_entities"].([]any)
		if !ok || len(entities) != 1 {
			t.Fatalf("unexpected option entities: %#v", second["text_entities"])
		}
		correctOptionIDs, ok := payload["correct_option_ids"].([]any)
		if !ok || len(correctOptionIDs) != 1 || correctOptionIDs[0] != float64(1) {
			t.Fatalf("unexpected correct_option_ids: %#v", payload["correct_option_ids"])
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true,"result":{"message_id":10,"chat":{"id":12345,"type":"private"},"date":100,"poll":{"id":"poll-id","question":"Pick one","options":[{"text":"A","voter_count":0},{"text":"B","voter_count":0}],"total_voter_count":0,"is_closed":false,"is_anonymous":true,"type":"quiz"}}}`))
	}))
	defer server.Close()

	bot := newTestBot(t, token, server.URL, server.Client())
	message, err := bot.SendPoll(context.Background(), SendPollParams{
		ChatID:   ChatIDInt(12345),
		Question: "Pick one",
		OptionObjects: []telegram.InputPollOption{
			{Text: "A", TextParseMode: "HTML"},
			{Text: "B", TextEntities: []telegram.MessageEntity{{Type: telegram.EntityCustomEmoji, Offset: 0, Length: 1, CustomEmojiID: "emoji-id"}}},
		},
		Type:             "quiz",
		CorrectOptionIDs: []int{1},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if message == nil || message.Poll == nil || message.Poll.ID != "poll-id" {
		t.Fatalf("unexpected message: %+v", message)
	}
}

func TestSendPollValidation(t *testing.T) {
	const token = "123:secret"
	valid := SendPollParams{ChatID: ChatIDInt(12345), Question: "Pick one", Options: []string{"A", "B"}}
	tests := []struct {
		name   string
		mutate func(*SendPollParams)
	}{
		{name: "empty chat", mutate: func(p *SendPollParams) { p.ChatID = ChatID{} }},
		{name: "empty question", mutate: func(p *SendPollParams) { p.Question = "" }},
		{name: "too few options", mutate: func(p *SendPollParams) { p.Options = []string{"A"} }},
		{name: "empty option", mutate: func(p *SendPollParams) { p.Options = []string{"A", ""} }},
		{name: "options and option objects", mutate: func(p *SendPollParams) {
			p.OptionObjects = []telegram.InputPollOption{{Text: "A"}, {Text: "B"}}
		}},
		{name: "empty structured option", mutate: func(p *SendPollParams) {
			p.Options = nil
			p.OptionObjects = []telegram.InputPollOption{{Text: "A"}, {}}
		}},
		{name: "structured option parse mode and entities", mutate: func(p *SendPollParams) {
			p.Options = nil
			p.OptionObjects = []telegram.InputPollOption{
				{Text: "A"},
				{Text: "B", TextParseMode: "HTML", TextEntities: []telegram.MessageEntity{{Type: telegram.EntityBold, Offset: 0, Length: 1}}},
			}
		}},
		{name: "question parse mode and entities", mutate: func(p *SendPollParams) {
			p.QuestionParseMode = "HTML"
			p.QuestionEntities = []telegram.MessageEntity{{Type: telegram.EntityCustomEmoji, Offset: 0, Length: 1, CustomEmojiID: "emoji-id"}}
		}},
		{name: "structured correct option out of range", mutate: func(p *SendPollParams) {
			p.Options = nil
			p.OptionObjects = []telegram.InputPollOption{{Text: "A"}, {Text: "B"}}
			p.CorrectOptionIDs = []int{2}
		}},
		{name: "negative thread", mutate: func(p *SendPollParams) { p.MessageThreadID = -1 }},
		{name: "negative open period", mutate: func(p *SendPollParams) { p.OpenPeriod = -1 }},
		{name: "negative close date", mutate: func(p *SendPollParams) { p.CloseDate = -1 }},
		{name: "negative correct option", mutate: func(p *SendPollParams) { id := -1; p.CorrectOptionID = &id }},
		{name: "out of range correct option", mutate: func(p *SendPollParams) { id := 2; p.CorrectOptionID = &id }},
		{name: "out of range correct option ids", mutate: func(p *SendPollParams) { p.CorrectOptionIDs = []int{2} }},
		{name: "negative correct option ids", mutate: func(p *SendPollParams) { p.CorrectOptionIDs = []int{-1} }},
		{name: "singular and plural correct options", mutate: func(p *SendPollParams) { id := 0; p.CorrectOptionID = &id; p.CorrectOptionIDs = []int{1} }},
		{name: "parse mode and entities", mutate: func(p *SendPollParams) {
			p.ExplanationParseMode = "HTML"
			p.ExplanationEntities = []telegram.MessageEntity{{Type: telegram.EntityBold, Offset: 0, Length: 1}}
		}},
		{name: "description parse mode and entities", mutate: func(p *SendPollParams) {
			p.DescriptionParseMode = "HTML"
			p.DescriptionEntities = []telegram.MessageEntity{{Type: telegram.EntityBold, Offset: 0, Length: 1}}
		}},
		{name: "invalid reply parameters", mutate: func(p *SendPollParams) { p.ReplyParameters = &telegram.ReplyParameters{} }},
		{name: "invalid reply markup", mutate: func(p *SendPollParams) { p.ReplyMarkup = telegram.InlineKeyboardMarkup{} }},
	}

	bot := newTestBot(t, token, "https://example.test", nil)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			params := valid
			tt.mutate(&params)
			message, err := bot.SendPoll(context.Background(), params)
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

func TestSendPollAPIAndTransportErrors(t *testing.T) {
	valid := SendPollParams{ChatID: ChatIDInt(12345), Question: "Pick one", Options: []string{"A", "B"}}
	testSendMethodErrorCases(t, "sendPoll", func(bot *Bot) (*telegram.Message, error) {
		return bot.SendPoll(context.Background(), valid)
	}, func(bot *Bot, ctx context.Context) (*telegram.Message, error) {
		return bot.SendPoll(ctx, valid)
	})
}

func TestStopPollSendsPayloadAndDecodesPoll(t *testing.T) {
	const token = "123:secret"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/bot"+token+"/stopPoll" {
			t.Fatalf("unexpected path: %q", r.URL.Path)
		}
		var payload map[string]any
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode payload: %v", err)
		}
		if got := payload["chat_id"]; got != float64(12345) {
			t.Fatalf("unexpected chat_id: %#v", got)
		}
		if got := payload["message_id"]; got != float64(99) {
			t.Fatalf("unexpected message_id: %#v", got)
		}
		markup := payload["reply_markup"].(map[string]any)
		if _, ok := markup["inline_keyboard"]; !ok {
			t.Fatalf("reply_markup.inline_keyboard missing: %#v", markup)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true,"result":{"id":"poll-id","question":"Pick one","options":[{"text":"A","voter_count":1},{"text":"B","voter_count":0}],"total_voter_count":1,"is_closed":true,"is_anonymous":true,"type":"regular"}}`))
	}))
	defer server.Close()

	markup := telegram.NewInlineKeyboard([]telegram.InlineKeyboardButton{telegram.InlineButtonCallback("OK", "ok")})
	bot := newTestBot(t, token, server.URL, server.Client())
	poll, err := bot.StopPoll(context.Background(), StopPollParams{ChatID: ChatIDInt(12345), MessageID: 99, ReplyMarkup: &markup})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if poll == nil || poll.ID != "poll-id" || !poll.IsClosed || len(poll.Options) != 2 {
		t.Fatalf("unexpected poll: %+v", poll)
	}
}

func TestStopPollValidation(t *testing.T) {
	const token = "123:secret"
	tests := []struct {
		name   string
		params StopPollParams
	}{
		{name: "empty chat", params: StopPollParams{MessageID: 1}},
		{name: "zero message", params: StopPollParams{ChatID: ChatIDInt(12345)}},
		{name: "negative message", params: StopPollParams{ChatID: ChatIDInt(12345), MessageID: -1}},
		{name: "invalid reply markup", params: StopPollParams{ChatID: ChatIDInt(12345), MessageID: 1, ReplyMarkup: &telegram.InlineKeyboardMarkup{}}},
	}

	bot := newTestBot(t, token, "https://example.test", nil)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			poll, err := bot.StopPoll(context.Background(), tt.params)
			if err == nil {
				t.Fatal("expected error")
			}
			if poll != nil {
				t.Fatalf("expected nil poll, got %+v", poll)
			}
			assertNoToken(t, err, token)
		})
	}
}

func TestStopPollAPIAndTransportErrors(t *testing.T) {
	const token = "123:secret"
	valid := StopPollParams{ChatID: ChatIDInt(12345), MessageID: 99}
	t.Run("api error", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/bot"+token+"/stopPoll" {
				t.Fatalf("unexpected path: %q", r.URL.Path)
			}
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"ok":false,"error_code":400,"description":"Bad Request"}`))
		}))
		defer server.Close()

		bot := newTestBot(t, token, server.URL, server.Client())
		poll, err := bot.StopPoll(context.Background(), valid)
		if err == nil {
			t.Fatal("expected error")
		}
		if poll != nil {
			t.Fatalf("expected nil poll, got %+v", poll)
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
		_, err := bot.StopPoll(context.Background(), valid)
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
		_, err := bot.StopPoll(context.Background(), valid)
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
		_, err := bot.StopPoll(ctx, valid)
		if err == nil {
			t.Fatal("expected error")
		}
		assertNoToken(t, err, token)
	})
}

func TestSendDiceSendsPayloadAndDecodesMessage(t *testing.T) {
	const token = "123:secret"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/bot"+token+"/sendDice" {
			t.Fatalf("unexpected path: %q", r.URL.Path)
		}
		var payload map[string]any
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode payload: %v", err)
		}
		if got := payload["chat_id"]; got != float64(12345) {
			t.Fatalf("unexpected chat_id: %#v", got)
		}
		if got := payload["emoji"]; got != "🎯" {
			t.Fatalf("unexpected emoji: %#v", got)
		}
		if got := payload["message_thread_id"]; got != float64(5) {
			t.Fatalf("unexpected message_thread_id: %#v", got)
		}
		reply := payload["reply_parameters"].(map[string]any)
		if got := reply["message_id"]; got != float64(13) {
			t.Fatalf("unexpected reply message_id: %#v", got)
		}
		markup := payload["reply_markup"].(map[string]any)
		if _, ok := markup["inline_keyboard"]; !ok {
			t.Fatalf("reply_markup.inline_keyboard missing: %#v", markup)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true,"result":{"message_id":11,"chat":{"id":12345,"type":"private"},"date":100,"dice":{"emoji":"🎯","value":6}}}`))
	}))
	defer server.Close()

	bot := newTestBot(t, token, server.URL, server.Client())
	message, err := bot.SendDice(context.Background(), SendDiceParams{
		ChatID:          ChatIDInt(12345),
		MessageThreadID: 5,
		Emoji:           "🎯",
		ReplyParameters: &telegram.ReplyParameters{MessageID: 13},
		ReplyMarkup:     telegram.NewInlineKeyboard([]telegram.InlineKeyboardButton{telegram.InlineButtonCallback("OK", "ok")}),
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if message == nil || message.Dice == nil || message.Dice.Emoji != "🎯" || message.Dice.Value != 6 {
		t.Fatalf("unexpected message: %+v", message)
	}
}

func TestSendDiceValidation(t *testing.T) {
	const token = "123:secret"
	valid := SendDiceParams{ChatID: ChatIDInt(12345)}
	tests := []struct {
		name   string
		mutate func(*SendDiceParams)
	}{
		{name: "empty chat", mutate: func(p *SendDiceParams) { p.ChatID = ChatID{} }},
		{name: "negative thread", mutate: func(p *SendDiceParams) { p.MessageThreadID = -1 }},
		{name: "invalid emoji", mutate: func(p *SendDiceParams) { p.Emoji = "😀" }},
		{name: "invalid reply parameters", mutate: func(p *SendDiceParams) { p.ReplyParameters = &telegram.ReplyParameters{} }},
		{name: "invalid reply markup", mutate: func(p *SendDiceParams) { p.ReplyMarkup = telegram.InlineKeyboardMarkup{} }},
	}

	bot := newTestBot(t, token, "https://example.test", nil)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			params := valid
			tt.mutate(&params)
			message, err := bot.SendDice(context.Background(), params)
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

func TestSendDiceAPIAndTransportErrors(t *testing.T) {
	valid := SendDiceParams{ChatID: ChatIDInt(12345), Emoji: "🎲"}
	testSendMethodErrorCases(t, "sendDice", func(bot *Bot) (*telegram.Message, error) {
		return bot.SendDice(context.Background(), valid)
	}, func(bot *Bot, ctx context.Context) (*telegram.Message, error) {
		return bot.SendDice(ctx, valid)
	})
}
