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

func TestInputPaidMediaMarshalAndValidation(t *testing.T) {
	photo := PaidPhoto(FileID("photo-file-id"))
	body, err := json.Marshal(photo)
	if err != nil {
		t.Fatalf("marshal photo: %v", err)
	}
	assertJSONFields(t, body, map[string]any{"type": "photo", "media": "photo-file-id"})
	if err := validateInputPaidMedia(photo); err != nil {
		t.Fatalf("photo validation: %v", err)
	}

	video := PaidVideo(FileURL("https://example.test/video.mp4"))
	video.Width = 640
	video.Height = 480
	video.Duration = 12
	video.SupportsStreaming = true
	if err := validateInputPaidMedia(video); err != nil {
		t.Fatalf("video validation: %v", err)
	}

	invalid := video
	invalid.Width = -1
	if err := validateInputPaidMedia(invalid); err == nil {
		t.Fatal("expected invalid width error")
	}
}

func TestSendPaidMediaSendsJSONAndDecodesMessage(t *testing.T) {
	const token = "123:secret"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/bot"+token+"/sendPaidMedia" {
			t.Fatalf("unexpected path: %q", r.URL.Path)
		}
		if got := r.Header.Get("Content-Type"); got != "application/json" {
			t.Fatalf("unexpected content type: %q", got)
		}
		var payload map[string]any
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode payload: %v", err)
		}
		if payload["business_connection_id"] != "business" || payload["chat_id"] != float64(12345) || payload["message_thread_id"] != float64(7) || payload["direct_messages_topic_id"] != float64(8) {
			t.Fatalf("unexpected chat payload: %#v", payload)
		}
		if payload["star_count"] != float64(5) || payload["payload"] != "safe-payload" || payload["caption"] != "caption" || payload["parse_mode"] != "HTML" {
			t.Fatalf("unexpected paid media payload: %#v", payload)
		}
		if payload["show_caption_above_media"] != true || payload["disable_notification"] != true || payload["protect_content"] != true || payload["allow_paid_broadcast"] != true {
			t.Fatalf("unexpected send flags: %#v", payload)
		}
		media, ok := payload["media"].([]any)
		if !ok || len(media) != 2 {
			t.Fatalf("unexpected media: %#v", payload["media"])
		}
		photo := media[0].(map[string]any)
		video := media[1].(map[string]any)
		if photo["type"] != "photo" || photo["media"] != "photo-file-id" {
			t.Fatalf("unexpected photo media: %#v", photo)
		}
		if video["type"] != "video" || video["media"] != "https://example.test/video.mp4" || video["cover"] != "cover-file-id" || video["start_timestamp"] != float64(3) || video["width"] != float64(640) || video["height"] != float64(480) || video["duration"] != float64(12) || video["supports_streaming"] != true {
			t.Fatalf("unexpected video media: %#v", video)
		}
		if _, ok := payload["reply_parameters"].(map[string]any); !ok {
			t.Fatalf("missing reply_parameters: %#v", payload)
		}
		if _, ok := payload["reply_markup"].(map[string]any); !ok {
			t.Fatalf("missing reply_markup: %#v", payload)
		}
		if suggestedPost, ok := payload["suggested_post_parameters"].(map[string]any); !ok || suggestedPost["send_date"] != float64(123) {
			t.Fatalf("unexpected suggested_post_parameters: %#v", payload["suggested_post_parameters"])
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true,"result":{"message_id":55,"chat":{"id":12345,"type":"private"},"date":1,"paid_media":{"star_count":5,"paid_media":[{"type":"preview","width":320,"height":240,"duration":10},{"type":"photo","photo":[{"file_id":"photo","file_unique_id":"photo-u","width":10,"height":10}]},{"type":"video","video":{"file_id":"video","file_unique_id":"video-u","width":640,"height":480,"duration":12}}]}}}`))
	}))
	defer server.Close()

	markup := telegram.NewInlineKeyboard([]telegram.InlineKeyboardButton{telegram.InlineButtonCallback("Buy", "buy")})
	video := PaidVideo(FileURL("https://example.test/video.mp4"))
	video.Cover = FileID("cover-file-id")
	video.StartTimestamp = 3
	video.Width = 640
	video.Height = 480
	video.Duration = 12
	video.SupportsStreaming = true
	bot := newTestBot(t, token, server.URL, server.Client())
	message, err := bot.SendPaidMedia(context.Background(), SendPaidMediaParams{
		BusinessConnectionID:    "business",
		ChatID:                  ChatIDInt(12345),
		MessageThreadID:         7,
		DirectMessagesTopicID:   8,
		StarCount:               5,
		Media:                   []InputPaidMedia{PaidPhoto(FileID("photo-file-id")), video},
		Payload:                 "safe-payload",
		Caption:                 "caption",
		ParseMode:               "HTML",
		ShowCaptionAboveMedia:   true,
		DisableNotification:     true,
		ProtectContent:          true,
		AllowPaidBroadcast:      true,
		SuggestedPostParameters: &telegram.SuggestedPostParameters{Price: &telegram.SuggestedPostPrice{Currency: "XTR", Amount: 1}, SendDate: 123},
		ReplyParameters:         &telegram.ReplyParameters{MessageID: 1},
		ReplyMarkup:             markup,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if message == nil || message.PaidMedia == nil || message.PaidMedia.StarCount != 5 || len(message.PaidMedia.PaidMedia) != 3 {
		t.Fatalf("unexpected message: %+v", message)
	}
	if _, ok := message.PaidMedia.PaidMedia[0].(telegram.PaidMediaPreview); !ok {
		t.Fatalf("unexpected paid media preview: %#v", message.PaidMedia.PaidMedia[0])
	}
}

func TestSendPaidMediaMultipartUpload(t *testing.T) {
	const token = "123:secret"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/bot"+token+"/sendPaidMedia" {
			t.Fatalf("unexpected path: %q", r.URL.Path)
		}
		if got := r.Header.Get("Content-Type"); !strings.HasPrefix(got, "multipart/form-data") {
			t.Fatalf("unexpected content type: %q", got)
		}
		if err := r.ParseMultipartForm(4096); err != nil {
			t.Fatalf("parse multipart: %v", err)
		}
		assertMultipartValue(t, r, "chat_id", "12345")
		assertMultipartValue(t, r, "star_count", "5")
		assertMultipartValue(t, r, "payload", "safe-payload")
		var media []map[string]any
		if err := json.Unmarshal([]byte(r.MultipartForm.Value["media"][0]), &media); err != nil {
			t.Fatalf("decode media field: %v", err)
		}
		if len(media) != 2 || media[0]["media"] != "attach://media0" || media[1]["media"] != "attach://media1" || media[1]["thumbnail"] != "attach://thumb1" || media[1]["cover"] != "attach://cover1" {
			t.Fatalf("unexpected media field: %#v", media)
		}
		content, header := readMultipartFile(t, r, "media0")
		if header.Filename != "photo.jpg" || string(content) != "photo-data" {
			t.Fatalf("unexpected media0 file: filename=%q content=%q", header.Filename, content)
		}
		content, header = readMultipartFile(t, r, "media1")
		if header.Filename != "video.mp4" || string(content) != "video-data" {
			t.Fatalf("unexpected media1 file: filename=%q content=%q", header.Filename, content)
		}
		content, header = readMultipartFile(t, r, "thumb1")
		if header.Filename != "thumb.jpg" || string(content) != "thumb-data" {
			t.Fatalf("unexpected thumb file: filename=%q content=%q", header.Filename, content)
		}
		content, header = readMultipartFile(t, r, "cover1")
		if header.Filename != "cover.jpg" || string(content) != "cover-data" {
			t.Fatalf("unexpected cover file: filename=%q content=%q", header.Filename, content)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true,"result":{"message_id":56,"chat":{"id":12345,"type":"private"},"date":1}}`))
	}))
	defer server.Close()

	video := PaidVideo(FileUpload(UploadFile{Name: "video.mp4", Reader: strings.NewReader("video-data"), ContentType: "video/mp4"}))
	video.Thumbnail = FileUpload(UploadFile{Name: "thumb.jpg", Reader: strings.NewReader("thumb-data"), ContentType: "image/jpeg"})
	video.Cover = FileUpload(UploadFile{Name: "cover.jpg", Reader: strings.NewReader("cover-data"), ContentType: "image/jpeg"})
	bot := newTestBot(t, token, server.URL, server.Client())
	message, err := bot.SendPaidMedia(context.Background(), SendPaidMediaParams{
		ChatID:    ChatIDInt(12345),
		StarCount: 5,
		Media: []InputPaidMedia{
			PaidPhoto(FileUpload(UploadFile{Name: "photo.jpg", Reader: strings.NewReader("photo-data"), ContentType: "image/jpeg"})),
			video,
		},
		Payload: "safe-payload",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if message == nil || message.MessageID != 56 {
		t.Fatalf("unexpected message: %+v", message)
	}
}

func TestSendPaidMediaValidation(t *testing.T) {
	const token = "123:secret"
	bot := newTestBot(t, token, "https://example.test", nil)
	valid := validSendPaidMediaParams()
	tests := []struct {
		name   string
		mutate func(*SendPaidMediaParams)
	}{
		{name: "invalid chat", mutate: func(p *SendPaidMediaParams) { p.ChatID = ChatID{} }},
		{name: "negative message thread", mutate: func(p *SendPaidMediaParams) { p.MessageThreadID = -1 }},
		{name: "negative direct topic", mutate: func(p *SendPaidMediaParams) { p.DirectMessagesTopicID = -1 }},
		{name: "zero stars", mutate: func(p *SendPaidMediaParams) { p.StarCount = 0 }},
		{name: "empty media", mutate: func(p *SendPaidMediaParams) { p.Media = nil }},
		{name: "too many media", mutate: func(p *SendPaidMediaParams) {
			p.Media = make([]InputPaidMedia, 11)
			for i := range p.Media {
				p.Media[i] = PaidPhoto(FileID("photo"))
			}
		}},
		{name: "invalid media", mutate: func(p *SendPaidMediaParams) { p.Media = []InputPaidMedia{PaidPhoto(FileID(""))} }},
		{name: "caption conflict", mutate: func(p *SendPaidMediaParams) {
			p.ParseMode = "HTML"
			p.CaptionEntities = []telegram.MessageEntity{{Type: telegram.EntityBold, Offset: 0, Length: 1}}
		}},
		{name: "invalid reply parameters", mutate: func(p *SendPaidMediaParams) { p.ReplyParameters = &telegram.ReplyParameters{} }},
		{name: "invalid reply markup", mutate: func(p *SendPaidMediaParams) { p.ReplyMarkup = telegram.InlineKeyboardMarkup{} }},
		{name: "invalid suggested post", mutate: func(p *SendPaidMediaParams) {
			p.SuggestedPostParameters = &telegram.SuggestedPostParameters{Price: &telegram.SuggestedPostPrice{Amount: 1}}
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			params := valid
			params.Media = append([]InputPaidMedia(nil), valid.Media...)
			params.Payload = "secret-payment-payload"
			tt.mutate(&params)
			message, err := bot.SendPaidMedia(context.Background(), params)
			if err == nil {
				t.Fatal("expected error")
			}
			if message != nil {
				t.Fatalf("expected nil message, got %+v", message)
			}
			assertNoToken(t, err, token)
			if strings.Contains(err.Error(), params.Payload) {
				t.Fatalf("error leaked payload: %q", err.Error())
			}
		})
	}
}

func TestSendPaidMediaErrors(t *testing.T) {
	t.Run("api error", func(t *testing.T) {
		const token = "123:secret"
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"ok":false,"error_code":400,"description":"Bad token 123:secret"}`))
		}))
		defer server.Close()
		bot := newTestBot(t, token, server.URL, server.Client())
		_, err := bot.SendPaidMedia(context.Background(), validSendPaidMediaParams())
		if err == nil {
			t.Fatal("expected error")
		}
		var apiErr *apierrors.APIError
		if !stderrors.As(err, &apiErr) || apiErr.Code != 400 {
			t.Fatalf("expected APIError, got %T %[1]v", err)
		}
		assertNoToken(t, err, token)
	})
	t.Run("invalid json", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { _, _ = w.Write([]byte(`not-json`)) }))
		defer server.Close()
		bot := newTestBot(t, "123:secret", server.URL, server.Client())
		_, err := bot.SendPaidMedia(context.Background(), validSendPaidMediaParams())
		if err == nil {
			t.Fatal("expected error")
		}
	})
	t.Run("http 500", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { http.Error(w, "boom", http.StatusInternalServerError) }))
		defer server.Close()
		bot := newTestBot(t, "123:secret", server.URL, server.Client())
		_, err := bot.SendPaidMedia(context.Background(), validSendPaidMediaParams())
		if err == nil {
			t.Fatal("expected error")
		}
	})
	t.Run("cancelled context", func(t *testing.T) {
		bot := newTestBot(t, "123:secret", "https://example.test", nil)
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		_, err := bot.SendPaidMedia(ctx, validSendPaidMediaParams())
		if err == nil {
			t.Fatal("expected error")
		}
	})
}

func TestGetStarTransactionsSendsPayloadAndDecodesTransactions(t *testing.T) {
	const token = "123:secret"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/bot"+token+"/getStarTransactions" {
			t.Fatalf("unexpected path: %q", r.URL.Path)
		}
		var payload map[string]any
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode payload: %v", err)
		}
		if payload["offset"] != float64(2) || payload["limit"] != float64(10) {
			t.Fatalf("unexpected payload: %#v", payload)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true,"result":{"transactions":[{"id":"tx1","amount":5,"nanostar_amount":7,"date":100,"source":{"type":"user","transaction_type":"paid_media_payment","user":{"id":7,"is_bot":false,"first_name":"Alice"},"paid_media":[{"type":"preview","width":10}],"paid_media_payload":"payload"}},{"id":"tx2","amount":-1,"date":101,"receiver":{"type":"telegram_api","request_count":3}}]}}`))
	}))
	defer server.Close()

	bot := newTestBot(t, token, server.URL, server.Client())
	transactions, err := bot.GetStarTransactions(context.Background(), GetStarTransactionsParams{Offset: 2, Limit: 10})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if transactions == nil || len(transactions.Transactions) != 2 || transactions.Transactions[0].NanostarAmount != 7 {
		t.Fatalf("unexpected transactions: %+v", transactions)
	}
	if _, ok := transactions.Transactions[0].Source.(telegram.TransactionPartnerUser); !ok {
		t.Fatalf("unexpected source: %#v", transactions.Transactions[0].Source)
	}
	if _, ok := transactions.Transactions[1].Receiver.(telegram.TransactionPartnerTelegramAPI); !ok {
		t.Fatalf("unexpected receiver: %#v", transactions.Transactions[1].Receiver)
	}
}

func TestGetStarTransactionsValidationAndErrors(t *testing.T) {
	bot := newTestBot(t, "123:secret", "https://example.test", nil)
	for _, params := range []GetStarTransactionsParams{{Offset: -1}, {Limit: -1}} {
		if _, err := bot.GetStarTransactions(context.Background(), params); err == nil {
			t.Fatalf("expected validation error for %+v", params)
		}
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":false,"error_code":400,"description":"Bad token 123:secret"}`))
	}))
	defer server.Close()
	bot = newTestBot(t, "123:secret", server.URL, server.Client())
	_, err := bot.GetStarTransactions(context.Background(), GetStarTransactionsParams{})
	if err == nil {
		t.Fatal("expected error")
	}
	var apiErr *apierrors.APIError
	if !stderrors.As(err, &apiErr) || apiErr.Code != 400 {
		t.Fatalf("expected APIError, got %T %[1]v", err)
	}
	assertNoToken(t, err, "123:secret")
}

func TestGetStarTransactionsTransportErrors(t *testing.T) {
	tests := []struct {
		name    string
		handler http.HandlerFunc
	}{
		{name: "invalid json", handler: func(w http.ResponseWriter, r *http.Request) { _, _ = w.Write([]byte(`not-json`)) }},
		{name: "http 500", handler: func(w http.ResponseWriter, r *http.Request) { http.Error(w, "boom", http.StatusInternalServerError) }},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(tt.handler)
			defer server.Close()
			bot := newTestBot(t, "123:secret", server.URL, server.Client())
			if transactions, err := bot.GetStarTransactions(context.Background(), GetStarTransactionsParams{}); err == nil || transactions != nil {
				t.Fatalf("expected error, got transactions=%+v err=%v", transactions, err)
			}
		})
	}

	t.Run("cancelled context", func(t *testing.T) {
		bot := newTestBot(t, "123:secret", "https://example.test", nil)
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		if transactions, err := bot.GetStarTransactions(ctx, GetStarTransactionsParams{}); err == nil || transactions != nil {
			t.Fatalf("expected error, got transactions=%+v err=%v", transactions, err)
		}
	})
}

func TestRefundStarPaymentSendsPayloadAndDecodesBool(t *testing.T) {
	const token = "123:secret"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/bot"+token+"/refundStarPayment" {
			t.Fatalf("unexpected path: %q", r.URL.Path)
		}
		var payload map[string]any
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode payload: %v", err)
		}
		if payload["user_id"] != float64(7) || payload["telegram_payment_charge_id"] != "charge-id" {
			t.Fatalf("unexpected payload: %#v", payload)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true,"result":true}`))
	}))
	defer server.Close()

	bot := newTestBot(t, token, server.URL, server.Client())
	ok, err := bot.RefundStarPayment(context.Background(), RefundStarPaymentParams{UserID: 7, TelegramPaymentChargeID: "charge-id"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Fatal("expected true")
	}
}

func TestRefundStarPaymentValidationAndErrors(t *testing.T) {
	bot := newTestBot(t, "123:secret", "https://example.test", nil)
	for _, params := range []RefundStarPaymentParams{{UserID: 0, TelegramPaymentChargeID: "charge"}, {UserID: -1, TelegramPaymentChargeID: "charge"}, {UserID: 7}} {
		if ok, err := bot.RefundStarPayment(context.Background(), params); err == nil || ok {
			t.Fatalf("expected validation error for %+v", params)
		}
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":false,"error_code":400,"description":"Bad token 123:secret"}`))
	}))
	defer server.Close()
	bot = newTestBot(t, "123:secret", server.URL, server.Client())
	ok, err := bot.RefundStarPayment(context.Background(), RefundStarPaymentParams{UserID: 7, TelegramPaymentChargeID: "charge-id"})
	if err == nil || ok {
		t.Fatal("expected error")
	}
	var apiErr *apierrors.APIError
	if !stderrors.As(err, &apiErr) || apiErr.Code != 400 {
		t.Fatalf("expected APIError, got %T %[1]v", err)
	}
	assertNoToken(t, err, "123:secret")
}

func TestRefundStarPaymentTransportErrors(t *testing.T) {
	tests := []struct {
		name    string
		handler http.HandlerFunc
	}{
		{name: "invalid json", handler: func(w http.ResponseWriter, r *http.Request) { _, _ = w.Write([]byte(`not-json`)) }},
		{name: "http 500", handler: func(w http.ResponseWriter, r *http.Request) { http.Error(w, "boom", http.StatusInternalServerError) }},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(tt.handler)
			defer server.Close()
			bot := newTestBot(t, "123:secret", server.URL, server.Client())
			if ok, err := bot.RefundStarPayment(context.Background(), RefundStarPaymentParams{UserID: 7, TelegramPaymentChargeID: "charge-id"}); err == nil || ok {
				t.Fatalf("expected error, got ok=%v err=%v", ok, err)
			}
		})
	}

	t.Run("cancelled context", func(t *testing.T) {
		bot := newTestBot(t, "123:secret", "https://example.test", nil)
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		if ok, err := bot.RefundStarPayment(ctx, RefundStarPaymentParams{UserID: 7, TelegramPaymentChargeID: "charge-id"}); err == nil || ok {
			t.Fatalf("expected error, got ok=%v err=%v", ok, err)
		}
	})
}

func validSendPaidMediaParams() SendPaidMediaParams {
	return SendPaidMediaParams{
		ChatID:    ChatIDInt(12345),
		StarCount: 5,
		Media:     []InputPaidMedia{PaidPhoto(FileID("photo-file-id"))},
	}
}

func assertJSONFields(t *testing.T, body []byte, want map[string]any) {
	t.Helper()
	var got map[string]any
	if err := json.Unmarshal(body, &got); err != nil {
		t.Fatalf("decode JSON: %v", err)
	}
	for key, value := range want {
		if got[key] != value {
			t.Fatalf("%s = %#v, want %#v in %s", key, got[key], value, body)
		}
	}
}
