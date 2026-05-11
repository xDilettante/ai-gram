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

func TestEditMessageMediaJSONChatTargetSendsPayloadAndDecodesMessage(t *testing.T) {
	const token = "123:secret"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/bot"+token+"/editMessageMedia" {
			t.Fatalf("unexpected path: %q", r.URL.Path)
		}
		var payload map[string]any
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode payload: %v", err)
		}
		if payload["business_connection_id"] != "bc-123" || payload["chat_id"] != float64(12345) || payload["message_id"] != float64(77) {
			t.Fatalf("unexpected target payload: %#v", payload)
		}
		media := payload["media"].(map[string]any)
		if media["type"] != "photo" || media["media"] != "photo-file-id" || media["caption"] != "updated" || media["show_caption_above_media"] != true {
			t.Fatalf("unexpected media payload: %#v", media)
		}
		if _, ok := payload["reply_markup"].(map[string]any)["inline_keyboard"]; !ok {
			t.Fatalf("reply_markup missing inline keyboard: %#v", payload)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true,"result":{"message_id":77,"chat":{"id":12345,"type":"private"},"date":100,"photo":[{"file_id":"photo-result","file_unique_id":"unique","width":1,"height":1}]}}`))
	}))
	defer server.Close()

	bot := newTestBot(t, token, server.URL, server.Client())
	media := MediaPhoto(FileID("photo-file-id"))
	media.Caption = "updated"
	media.ShowCaptionAboveMedia = true
	markup := telegram.NewInlineKeyboard([]telegram.InlineKeyboardButton{telegram.InlineButtonCallback("OK", "ok")})
	result, err := bot.EditMessageMedia(context.Background(), EditMessageMediaParams{
		BusinessConnectionID: "bc-123",
		Target:               EditTargetChat(ChatIDInt(12345), 77),
		Media:                media,
		ReplyMarkup:          &markup,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result == nil || !result.IsOK() || !result.IsMessage() || len(result.Message.Photo) != 1 {
		t.Fatalf("unexpected result: %+v", result)
	}
}

func TestEditMessageMediaJSONInlineTargetDecodesTrue(t *testing.T) {
	const token = "123:secret"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/bot"+token+"/editMessageMedia" {
			t.Fatalf("unexpected path: %q", r.URL.Path)
		}
		var payload map[string]any
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode payload: %v", err)
		}
		if payload["inline_message_id"] != "inline-id" {
			t.Fatalf("unexpected inline_message_id: %#v", payload)
		}
		if _, ok := payload["chat_id"]; ok {
			t.Fatalf("chat_id should be omitted: %#v", payload)
		}
		media := payload["media"].(map[string]any)
		if media["type"] != "video" || media["media"] != "https://example.test/video.mp4" {
			t.Fatalf("unexpected media payload: %#v", media)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true,"result":true}`))
	}))
	defer server.Close()

	bot := newTestBot(t, token, server.URL, server.Client())
	result, err := bot.EditMessageMedia(context.Background(), EditMessageMediaParams{Target: EditTargetInline("inline-id"), Media: MediaVideo(FileURL("https://example.test/video.mp4"))})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result == nil || !result.IsOK() || result.IsMessage() {
		t.Fatalf("unexpected result: %+v", result)
	}
}

func TestEditMessageMediaMultipartUpload(t *testing.T) {
	const token = "123:secret"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/bot"+token+"/editMessageMedia" {
			t.Fatalf("unexpected path: %q", r.URL.Path)
		}
		if got := r.Header.Get("Content-Type"); !strings.HasPrefix(got, "multipart/form-data") {
			t.Fatalf("unexpected content type: %q", got)
		}
		if r.ContentLength <= 0 {
			t.Fatalf("editMessageMedia upload must send content length, got %d", r.ContentLength)
		}
		if err := r.ParseMultipartForm(4096); err != nil {
			t.Fatalf("parse multipart: %v", err)
		}
		if r.MultipartForm.Value["business_connection_id"][0] != "bc-123" || r.MultipartForm.Value["chat_id"][0] != "12345" || r.MultipartForm.Value["message_id"][0] != "77" {
			t.Fatalf("unexpected multipart fields: %#v", r.MultipartForm.Value)
		}
		var media map[string]any
		if err := json.Unmarshal([]byte(r.MultipartForm.Value["media"][0]), &media); err != nil {
			t.Fatalf("decode media field: %v", err)
		}
		if media["type"] != "animation" || media["media"] != "attach://media0" || media["thumbnail"] != "attach://thumb0" || media["caption"] != "updated animation" {
			t.Fatalf("unexpected media field: %#v", media)
		}
		if len(r.MultipartForm.File["media0"]) != 1 || len(r.MultipartForm.File["thumb0"]) != 1 {
			t.Fatalf("expected media and thumbnail upload files: %#v", r.MultipartForm.File)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true,"result":{"message_id":77,"chat":{"id":12345,"type":"private"},"date":100,"animation":{"file_id":"anim","file_unique_id":"unique","width":1,"height":1,"duration":1}}}`))
	}))
	defer server.Close()

	bot := newTestBot(t, token, server.URL, server.Client())
	media := MediaAnimation(FileUpload(UploadFile{Name: "animation.gif", Reader: strings.NewReader("animation-data"), ContentType: "image/gif"}))
	media.Thumbnail = FileUpload(UploadFile{Name: "thumb.jpg", Reader: strings.NewReader("thumb-data"), ContentType: "image/jpeg"})
	media.Caption = "updated animation"
	result, err := bot.EditMessageMedia(context.Background(), EditMessageMediaParams{BusinessConnectionID: "bc-123", Target: EditTargetChat(ChatIDInt(12345), 77), Media: media})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result == nil || !result.IsOK() || !result.IsMessage() || result.Message.Animation == nil {
		t.Fatalf("unexpected result: %+v", result)
	}
}

func TestEditMessageMediaLivePhotoJSONAndMultipart(t *testing.T) {
	const token = "123:secret"

	t.Run("chat json", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/bot"+token+"/editMessageMedia" {
				t.Fatalf("unexpected path: %q", r.URL.Path)
			}
			var payload map[string]any
			if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
				t.Fatalf("decode payload: %v", err)
			}
			media := payload["media"].(map[string]any)
			if media["type"] != "live_photo" || media["media"] != "live-file-id" || media["photo"] != "photo-file-id" || media["caption"] != "updated live" || media["show_caption_above_media"] != true || media["has_spoiler"] != true {
				t.Fatalf("unexpected media payload: %#v", media)
			}
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"ok":true,"result":{"message_id":77,"chat":{"id":12345,"type":"private"},"date":100,"live_photo":{"photo":[{"file_id":"photo-result","file_unique_id":"photo-u","width":640,"height":480}],"file_id":"live-result","file_unique_id":"live-u","width":640,"height":480,"duration":3}}}`))
		}))
		defer server.Close()

		media := MediaLivePhoto(FileID("live-file-id"), FileID("photo-file-id"))
		media.Caption = "updated live"
		media.ShowCaptionAboveMedia = true
		media.HasSpoiler = true
		bot := newTestBot(t, token, server.URL, server.Client())
		result, err := bot.EditMessageMedia(context.Background(), EditMessageMediaParams{Target: EditTargetChat(ChatIDInt(12345), 77), Media: media})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result == nil || !result.IsOK() || !result.IsMessage() || result.Message.LivePhoto == nil || result.Message.LivePhoto.FileID != "live-result" {
			t.Fatalf("unexpected result: %+v", result)
		}
	})

	t.Run("chat multipart", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/bot"+token+"/editMessageMedia" {
				t.Fatalf("unexpected path: %q", r.URL.Path)
			}
			if got := r.Header.Get("Content-Type"); !strings.HasPrefix(got, "multipart/form-data") {
				t.Fatalf("unexpected content type: %q", got)
			}
			if err := r.ParseMultipartForm(4096); err != nil {
				t.Fatalf("parse multipart: %v", err)
			}
			var media map[string]any
			if err := json.Unmarshal([]byte(r.MultipartForm.Value["media"][0]), &media); err != nil {
				t.Fatalf("decode media field: %v", err)
			}
			if media["type"] != "live_photo" || media["media"] != "attach://media0" || media["photo"] != "attach://photo0" {
				t.Fatalf("unexpected media field: %#v", media)
			}
			content, header := readMultipartFile(t, r, "media0")
			if header.Filename != "live.mp4" || string(content) != "live-data" {
				t.Fatalf("unexpected media0 file: filename=%q content=%q", header.Filename, content)
			}
			content, header = readMultipartFile(t, r, "photo0")
			if header.Filename != "photo.jpg" || string(content) != "photo-data" {
				t.Fatalf("unexpected photo0 file: filename=%q content=%q", header.Filename, content)
			}
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"ok":true,"result":{"message_id":77,"chat":{"id":12345,"type":"private"},"date":100}}`))
		}))
		defer server.Close()

		media := MediaLivePhoto(
			FileUpload(UploadFile{Name: "live.mp4", Reader: strings.NewReader("live-data"), ContentType: "video/mp4"}),
			FileUpload(UploadFile{Name: "photo.jpg", Reader: strings.NewReader("photo-data"), ContentType: "image/jpeg"}),
		)
		bot := newTestBot(t, token, server.URL, server.Client())
		result, err := bot.EditMessageMedia(context.Background(), EditMessageMediaParams{Target: EditTargetChat(ChatIDInt(12345), 77), Media: media})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result == nil || !result.IsOK() || !result.IsMessage() {
			t.Fatalf("unexpected result: %+v", result)
		}
	})

	t.Run("inline file id", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var payload map[string]any
			if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
				t.Fatalf("decode payload: %v", err)
			}
			if payload["inline_message_id"] != "inline-id" {
				t.Fatalf("unexpected inline target: %#v", payload)
			}
			media := payload["media"].(map[string]any)
			if media["type"] != "live_photo" || media["media"] != "live-file-id" || media["photo"] != "photo-file-id" {
				t.Fatalf("unexpected media payload: %#v", media)
			}
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"ok":true,"result":true}`))
		}))
		defer server.Close()

		bot := newTestBot(t, token, server.URL, server.Client())
		result, err := bot.EditMessageMedia(context.Background(), EditMessageMediaParams{Target: EditTargetInline("inline-id"), Media: MediaLivePhoto(FileID("live-file-id"), FileID("photo-file-id"))})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result == nil || !result.IsOK() || result.IsMessage() {
			t.Fatalf("unexpected result: %+v", result)
		}
	})
}

func TestInputMediaAnimationValidationAndMarshal(t *testing.T) {
	media := MediaAnimation(FileID("animation-file-id"))
	media.Caption = "caption"
	media.ParseMode = "HTML"
	media.ShowCaptionAboveMedia = true
	media.Width = 320
	media.Height = 240
	media.Duration = 5
	if err := validateInputMediaForEdit(media); err != nil {
		t.Fatalf("unexpected validation error: %v", err)
	}
	payload, _, err := (EditMessageMediaParams{Target: EditTargetChat(ChatIDInt(1), 1), Media: media}).mediaPayload()
	if err != nil {
		t.Fatalf("unexpected payload error: %v", err)
	}
	body, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("marshal payload: %v", err)
	}
	var got map[string]any
	if err := json.Unmarshal(body, &got); err != nil {
		t.Fatalf("unmarshal payload: %v", err)
	}
	if got["type"] != "animation" || got["media"] != "animation-file-id" || got["show_caption_above_media"] != true {
		t.Fatalf("unexpected payload: %#v", got)
	}
}

func TestEditMessageMediaValidation(t *testing.T) {
	const token = "123:secret"
	bot := newTestBot(t, token, "https://example.test", nil)
	invalidMarkup := telegram.InlineKeyboardMarkup{}
	tests := []struct {
		name   string
		params EditMessageMediaParams
	}{
		{name: "invalid target", params: EditMessageMediaParams{Media: MediaPhoto(FileID("photo"))}},
		{name: "nil media", params: EditMessageMediaParams{Target: EditTargetChat(ChatIDInt(123), 1)}},
		{name: "invalid media", params: EditMessageMediaParams{Target: EditTargetChat(ChatIDInt(123), 1), Media: MediaPhoto(FileID(""))}},
		{name: "invalid reply markup", params: EditMessageMediaParams{Target: EditTargetChat(ChatIDInt(123), 1), Media: MediaPhoto(FileID("photo")), ReplyMarkup: &invalidMarkup}},
		{name: "inline upload", params: EditMessageMediaParams{Target: EditTargetInline("inline-id"), Media: MediaPhoto(FileUpload(UploadFile{Name: "photo.jpg", Reader: strings.NewReader("photo")}))}},
		{name: "inline live photo upload", params: EditMessageMediaParams{Target: EditTargetInline("inline-id"), Media: MediaLivePhoto(FileUpload(UploadFile{Name: "live.mp4", Reader: strings.NewReader("live")}), FileID("photo"))}},
		{name: "live photo media url", params: EditMessageMediaParams{Target: EditTargetChat(ChatIDInt(123), 1), Media: MediaLivePhoto(FileURL("https://example.test/live.mp4"), FileID("photo"))}},
		{name: "live photo preview url", params: EditMessageMediaParams{Target: EditTargetChat(ChatIDInt(123), 1), Media: MediaLivePhoto(FileID("live"), FileURL("https://example.test/photo.jpg"))}},
		{name: "invalid animation type", params: EditMessageMediaParams{Target: EditTargetChat(ChatIDInt(123), 1), Media: InputMediaAnimation{Type: "video", Media: FileID("animation")}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := bot.EditMessageMedia(context.Background(), tt.params)
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

func TestSendMediaGroupRejectsAnimation(t *testing.T) {
	const token = "123:secret"
	bot := newTestBot(t, token, "https://example.test", nil)
	_, err := bot.SendMediaGroup(context.Background(), SendMediaGroupParams{
		ChatID: ChatIDInt(12345),
		Media:  []InputMedia{MediaPhoto(FileID("photo")), MediaAnimation(FileID("animation"))},
	})
	if err == nil {
		t.Fatal("expected error")
	}
	assertNoToken(t, err, token)
}

func TestEditMessageMediaErrorCases(t *testing.T) {
	const token = "123:secret"
	params := EditMessageMediaParams{Target: EditTargetChat(ChatIDInt(123), 1), Media: MediaPhoto(FileID("photo"))}

	t.Run("api error", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"ok":false,"error_code":400,"description":"Bad Request"}`))
		}))
		defer server.Close()
		bot := newTestBot(t, token, server.URL, server.Client())
		result, err := bot.EditMessageMedia(context.Background(), params)
		if err == nil {
			t.Fatal("expected error")
		}
		if result != nil {
			t.Fatalf("expected nil result, got %+v", result)
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
		result, err := bot.EditMessageMedia(context.Background(), params)
		if err == nil {
			t.Fatal("expected error")
		}
		if result != nil {
			t.Fatalf("expected nil result, got %+v", result)
		}
		assertNoToken(t, err, token)
	})

	t.Run("http status", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "server error", http.StatusInternalServerError)
		}))
		defer server.Close()
		bot := newTestBot(t, token, server.URL, server.Client())
		result, err := bot.EditMessageMedia(context.Background(), params)
		if err == nil {
			t.Fatal("expected error")
		}
		if result != nil {
			t.Fatalf("expected nil result, got %+v", result)
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
		result, err := bot.EditMessageMedia(ctx, params)
		if err == nil {
			t.Fatal("expected error")
		}
		if result != nil {
			t.Fatalf("expected nil result, got %+v", result)
		}
		assertNoToken(t, err, token)
	})
}

func TestEditMessageLiveLocationSendsPayloadAndDecodesMessage(t *testing.T) {
	const token = "123:secret"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/bot"+token+"/editMessageLiveLocation" {
			t.Fatalf("unexpected path: %q", r.URL.Path)
		}
		var payload map[string]any
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode payload: %v", err)
		}
		if payload["business_connection_id"] != "bc-123" || payload["chat_id"] != float64(12345) || payload["message_id"] != float64(77) {
			t.Fatalf("unexpected target payload: %#v", payload)
		}
		if payload["latitude"] != 51.5 || payload["longitude"] != -0.12 || payload["live_period"] != float64(3600) || payload["horizontal_accuracy"] != 15.5 || payload["heading"] != float64(90) || payload["proximity_alert_radius"] != float64(100) {
			t.Fatalf("unexpected location payload: %#v", payload)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true,"result":{"message_id":77,"chat":{"id":12345,"type":"private"},"date":100,"location":{"latitude":51.5,"longitude":-0.12}}}`))
	}))
	defer server.Close()

	bot := newTestBot(t, token, server.URL, server.Client())
	markup := telegram.NewInlineKeyboard([]telegram.InlineKeyboardButton{telegram.InlineButtonCallback("OK", "ok")})
	result, err := bot.EditMessageLiveLocation(context.Background(), EditMessageLiveLocationParams{
		BusinessConnectionID: "bc-123",
		Target:               EditTargetChat(ChatIDInt(12345), 77),
		Latitude:             51.5,
		Longitude:            -0.12,
		LivePeriod:           3600,
		HorizontalAccuracy:   15.5,
		Heading:              90,
		ProximityAlertRadius: 100,
		ReplyMarkup:          &markup,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result == nil || !result.IsOK() || !result.IsMessage() || result.Message.Location == nil {
		t.Fatalf("unexpected result: %+v", result)
	}
}

func TestEditMessageLiveLocationInlineTargetDecodesTrue(t *testing.T) {
	const token = "123:secret"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var payload map[string]any
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode payload: %v", err)
		}
		if payload["inline_message_id"] != "inline-id" || payload["latitude"] != float64(52) || payload["longitude"] != float64(4) {
			t.Fatalf("unexpected payload: %#v", payload)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true,"result":true}`))
	}))
	defer server.Close()

	bot := newTestBot(t, token, server.URL, server.Client())
	result, err := bot.EditMessageLiveLocation(context.Background(), EditMessageLiveLocationParams{Target: EditTargetInline("inline-id"), Latitude: 52, Longitude: 4})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result == nil || !result.IsOK() || result.IsMessage() {
		t.Fatalf("unexpected result: %+v", result)
	}
}

func TestEditMessageLiveLocationValidation(t *testing.T) {
	const token = "123:secret"
	bot := newTestBot(t, token, "https://example.test", nil)
	valid := EditMessageLiveLocationParams{Target: EditTargetChat(ChatIDInt(123), 1), Latitude: 52, Longitude: 4}
	invalidMarkup := telegram.InlineKeyboardMarkup{}
	tests := []struct {
		name   string
		mutate func(*EditMessageLiveLocationParams)
	}{
		{name: "invalid target", mutate: func(p *EditMessageLiveLocationParams) { p.Target = EditMessageTarget{} }},
		{name: "latitude too small", mutate: func(p *EditMessageLiveLocationParams) { p.Latitude = -91 }},
		{name: "longitude too large", mutate: func(p *EditMessageLiveLocationParams) { p.Longitude = 181 }},
		{name: "negative live period", mutate: func(p *EditMessageLiveLocationParams) { p.LivePeriod = -1 }},
		{name: "negative horizontal accuracy", mutate: func(p *EditMessageLiveLocationParams) { p.HorizontalAccuracy = -1 }},
		{name: "negative heading", mutate: func(p *EditMessageLiveLocationParams) { p.Heading = -1 }},
		{name: "negative proximity", mutate: func(p *EditMessageLiveLocationParams) { p.ProximityAlertRadius = -1 }},
		{name: "invalid markup", mutate: func(p *EditMessageLiveLocationParams) { p.ReplyMarkup = &invalidMarkup }},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			params := valid
			tt.mutate(&params)
			result, err := bot.EditMessageLiveLocation(context.Background(), params)
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

func TestStopMessageLiveLocationSendsPayloadAndDecodesMessage(t *testing.T) {
	const token = "123:secret"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/bot"+token+"/stopMessageLiveLocation" {
			t.Fatalf("unexpected path: %q", r.URL.Path)
		}
		var payload map[string]any
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode payload: %v", err)
		}
		if payload["business_connection_id"] != "bc-123" || payload["chat_id"] != float64(12345) || payload["message_id"] != float64(77) {
			t.Fatalf("unexpected target payload: %#v", payload)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true,"result":{"message_id":77,"chat":{"id":12345,"type":"private"},"date":100,"location":{"latitude":51.5,"longitude":-0.12}}}`))
	}))
	defer server.Close()

	bot := newTestBot(t, token, server.URL, server.Client())
	markup := telegram.NewInlineKeyboard([]telegram.InlineKeyboardButton{telegram.InlineButtonCallback("OK", "ok")})
	result, err := bot.StopMessageLiveLocation(context.Background(), StopMessageLiveLocationParams{BusinessConnectionID: "bc-123", Target: EditTargetChat(ChatIDInt(12345), 77), ReplyMarkup: &markup})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result == nil || !result.IsOK() || !result.IsMessage() || result.Message.Location == nil {
		t.Fatalf("unexpected result: %+v", result)
	}
}

func TestStopMessageLiveLocationInlineTargetDecodesTrue(t *testing.T) {
	const token = "123:secret"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var payload map[string]any
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode payload: %v", err)
		}
		if payload["inline_message_id"] != "inline-id" {
			t.Fatalf("unexpected payload: %#v", payload)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true,"result":true}`))
	}))
	defer server.Close()

	bot := newTestBot(t, token, server.URL, server.Client())
	result, err := bot.StopMessageLiveLocation(context.Background(), StopMessageLiveLocationParams{Target: EditTargetInline("inline-id")})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result == nil || !result.IsOK() || result.IsMessage() {
		t.Fatalf("unexpected result: %+v", result)
	}
}

func TestStopMessageLiveLocationValidation(t *testing.T) {
	const token = "123:secret"
	bot := newTestBot(t, token, "https://example.test", nil)
	invalidMarkup := telegram.InlineKeyboardMarkup{}
	tests := []StopMessageLiveLocationParams{
		{},
		{Target: EditTargetChat(ChatIDInt(123), 1), ReplyMarkup: &invalidMarkup},
	}
	for _, params := range tests {
		result, err := bot.StopMessageLiveLocation(context.Background(), params)
		if err == nil {
			t.Fatal("expected error")
		}
		if result != nil {
			t.Fatalf("expected nil result, got %+v", result)
		}
		assertNoToken(t, err, token)
	}
}

func TestLiveLocationEditErrorCases(t *testing.T) {
	const token = "123:secret"
	tests := []struct {
		name string
		call func(context.Context, *Bot) (*EditMessageResult, error)
	}{
		{name: "edit live", call: func(ctx context.Context, bot *Bot) (*EditMessageResult, error) {
			return bot.EditMessageLiveLocation(ctx, EditMessageLiveLocationParams{Target: EditTargetChat(ChatIDInt(123), 1), Latitude: 52, Longitude: 4})
		}},
		{name: "stop live", call: func(ctx context.Context, bot *Bot) (*EditMessageResult, error) {
			return bot.StopMessageLiveLocation(ctx, StopMessageLiveLocationParams{Target: EditTargetChat(ChatIDInt(123), 1)})
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
			result, err := tt.call(context.Background(), bot)
			if err == nil {
				t.Fatal("expected error")
			}
			if result != nil {
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
			result, err := tt.call(context.Background(), bot)
			if err == nil {
				t.Fatal("expected error")
			}
			if result != nil {
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
			result, err := tt.call(context.Background(), bot)
			if err == nil {
				t.Fatal("expected error")
			}
			if result != nil {
				t.Fatalf("expected nil result, got %+v", result)
			}
			assertNoToken(t, err, token)
		})

		t.Run(tt.name+" cancelled context", func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				t.Fatal("request should not reach server")
			}))
			defer server.Close()
			bot := newTestBot(t, token, server.URL, server.Client())
			ctx, cancel := context.WithCancel(context.Background())
			cancel()
			result, err := tt.call(ctx, bot)
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
