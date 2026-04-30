package bot

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/xDilettante/ai-gram/telegram"
)

func TestSendStickerSendsJSONFileIDAndDecodesMessage(t *testing.T) {
	const token = "123:secret"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/bot"+token+"/sendSticker" {
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
		if got := payload["sticker"]; got != "sticker-file-id" {
			t.Fatalf("unexpected sticker: %#v", got)
		}
		if got := payload["emoji"]; got != "😀" {
			t.Fatalf("unexpected emoji: %#v", got)
		}
		reply := payload["reply_parameters"].(map[string]any)
		if got := reply["message_id"]; got != float64(10) {
			t.Fatalf("unexpected reply message_id: %#v", got)
		}
		markup := payload["reply_markup"].(map[string]any)
		if _, ok := markup["inline_keyboard"]; !ok {
			t.Fatalf("reply_markup.inline_keyboard missing: %#v", markup)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true,"result":{"message_id":5,"chat":{"id":12345,"type":"private"},"date":100,"sticker":{"file_id":"sent-sticker","file_unique_id":"sticker-unique","type":"regular","width":512,"height":512,"is_animated":false,"is_video":false,"emoji":"😀"}}}`))
	}))
	defer server.Close()

	bot := newTestBot(t, token, server.URL, server.Client())
	message, err := bot.SendSticker(context.Background(), SendStickerParams{
		ChatID:          ChatIDInt(12345),
		MessageThreadID: 7,
		Sticker:         FileID("sticker-file-id"),
		Emoji:           "😀",
		ReplyParameters: &telegram.ReplyParameters{MessageID: 10},
		ReplyMarkup:     telegram.NewInlineKeyboard([]telegram.InlineKeyboardButton{telegram.InlineButtonCallback("OK", "ok")}),
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if message == nil || message.Sticker == nil || message.Sticker.FileID != "sent-sticker" || message.Sticker.Emoji != "😀" {
		t.Fatalf("unexpected message: %+v", message)
	}
}

func TestSendStickerSendsJSONFileURL(t *testing.T) {
	const token = "123:secret"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/bot"+token+"/sendSticker" {
			t.Fatalf("unexpected path: %q", r.URL.Path)
		}
		var payload map[string]any
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode payload: %v", err)
		}
		if got := payload["sticker"]; got != "https://example.com/sticker.webp" {
			t.Fatalf("unexpected sticker URL: %#v", got)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true,"result":{"message_id":6,"chat":{"id":12345,"type":"private"},"date":101,"sticker":{"file_id":"url-sticker","file_unique_id":"url-sticker-unique","type":"regular","width":512,"height":512,"is_animated":false,"is_video":false}}}`))
	}))
	defer server.Close()

	bot := newTestBot(t, token, server.URL, server.Client())
	message, err := bot.SendSticker(context.Background(), SendStickerParams{ChatID: ChatIDInt(12345), Sticker: FileURL("https://example.com/sticker.webp")})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if message == nil || message.Sticker == nil || message.Sticker.FileID != "url-sticker" {
		t.Fatalf("unexpected message: %+v", message)
	}
}

func TestSendStickerMultipartUpload(t *testing.T) {
	const token = "123:secret"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/bot"+token+"/sendSticker" {
			t.Fatalf("unexpected path: %q", r.URL.Path)
		}
		if got := r.Header.Get("Content-Type"); !strings.HasPrefix(got, "multipart/form-data") {
			t.Fatalf("unexpected content type: %q", got)
		}
		if err := r.ParseMultipartForm(2048); err != nil {
			t.Fatalf("parse multipart: %v", err)
		}
		assertMultipartValue(t, r, "chat_id", "12345")
		assertMultipartValue(t, r, "message_thread_id", "8")
		assertMultipartValue(t, r, "sticker", "attach://sticker")
		assertMultipartValue(t, r, "emoji", "🙂")
		content, header := readMultipartFile(t, r, "sticker")
		if header.Filename != "sticker.webp" {
			t.Fatalf("unexpected filename: %q", header.Filename)
		}
		if string(content) != "sticker-data" {
			t.Fatalf("unexpected file content: %q", content)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true,"result":{"message_id":7,"chat":{"id":12345,"type":"private"},"date":102,"sticker":{"file_id":"uploaded-sticker","file_unique_id":"uploaded-sticker-unique","type":"regular","width":512,"height":512,"is_animated":false,"is_video":false}}}`))
	}))
	defer server.Close()

	bot := newTestBot(t, token, server.URL, server.Client())
	message, err := bot.SendSticker(context.Background(), SendStickerParams{
		ChatID:          ChatIDInt(12345),
		MessageThreadID: 8,
		Sticker:         FileUpload(UploadFile{Name: "sticker.webp", Reader: strings.NewReader("sticker-data"), ContentType: "image/webp"}),
		Emoji:           "🙂",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if message == nil || message.Sticker == nil || message.Sticker.FileID != "uploaded-sticker" {
		t.Fatalf("unexpected message: %+v", message)
	}
}

func TestSendStickerValidation(t *testing.T) {
	const token = "123:secret"
	valid := SendStickerParams{ChatID: ChatIDInt(12345), Sticker: FileID("sticker")}
	tests := []struct {
		name   string
		mutate func(*SendStickerParams)
	}{
		{name: "empty chat", mutate: func(p *SendStickerParams) { p.ChatID = ChatID{} }},
		{name: "empty sticker", mutate: func(p *SendStickerParams) { p.Sticker = FileRef{} }},
		{name: "negative thread", mutate: func(p *SendStickerParams) { p.MessageThreadID = -1 }},
		{name: "invalid reply parameters", mutate: func(p *SendStickerParams) { p.ReplyParameters = &telegram.ReplyParameters{} }},
		{name: "invalid reply markup", mutate: func(p *SendStickerParams) { p.ReplyMarkup = telegram.InlineKeyboardMarkup{} }},
	}

	bot := newTestBot(t, token, "https://example.test", nil)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			params := valid
			tt.mutate(&params)
			message, err := bot.SendSticker(context.Background(), params)
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

func TestSendStickerAPIAndTransportErrors(t *testing.T) {
	valid := SendStickerParams{ChatID: ChatIDInt(12345), Sticker: FileID("sticker")}
	testSendMethodErrorCases(t, "sendSticker", func(bot *Bot) (*telegram.Message, error) {
		return bot.SendSticker(context.Background(), valid)
	}, func(bot *Bot, ctx context.Context) (*telegram.Message, error) {
		return bot.SendSticker(ctx, valid)
	})
}

func TestSendAnimationSendsJSONFileIDAndDecodesMessage(t *testing.T) {
	const token = "123:secret"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/bot"+token+"/sendAnimation" {
			t.Fatalf("unexpected path: %q", r.URL.Path)
		}
		var payload map[string]any
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode payload: %v", err)
		}
		if got := payload["animation"]; got != "animation-file-id" {
			t.Fatalf("unexpected animation: %#v", got)
		}
		if got := payload["thumbnail"]; got != "https://example.com/thumb.jpg" {
			t.Fatalf("unexpected thumbnail: %#v", got)
		}
		if got := payload["duration"]; got != float64(10) {
			t.Fatalf("unexpected duration: %#v", got)
		}
		if got := payload["width"]; got != float64(320) {
			t.Fatalf("unexpected width: %#v", got)
		}
		if got := payload["height"]; got != float64(240) {
			t.Fatalf("unexpected height: %#v", got)
		}
		if got := payload["caption"]; got != "caption" {
			t.Fatalf("unexpected caption: %#v", got)
		}
		if got := payload["parse_mode"]; got != "HTML" {
			t.Fatalf("unexpected parse_mode: %#v", got)
		}
		if got := payload["show_caption_above_media"]; got != true {
			t.Fatalf("unexpected show_caption_above_media: %#v", got)
		}
		if got := payload["has_spoiler"]; got != true {
			t.Fatalf("unexpected has_spoiler: %#v", got)
		}
		reply := payload["reply_parameters"].(map[string]any)
		if got := reply["message_id"]; got != float64(20) {
			t.Fatalf("unexpected reply message_id: %#v", got)
		}
		markup := payload["reply_markup"].(map[string]any)
		if _, ok := markup["inline_keyboard"]; !ok {
			t.Fatalf("reply_markup.inline_keyboard missing: %#v", markup)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true,"result":{"message_id":8,"chat":{"id":12345,"type":"private"},"date":103,"animation":{"file_id":"sent-animation","file_unique_id":"animation-unique","width":320,"height":240,"duration":10,"file_name":"a.gif","mime_type":"image/gif"}}}`))
	}))
	defer server.Close()

	bot := newTestBot(t, token, server.URL, server.Client())
	message, err := bot.SendAnimation(context.Background(), SendAnimationParams{
		ChatID:                ChatIDInt(12345),
		Animation:             FileID("animation-file-id"),
		Duration:              10,
		Width:                 320,
		Height:                240,
		Thumbnail:             FileURL("https://example.com/thumb.jpg"),
		Caption:               "caption",
		ParseMode:             "HTML",
		ShowCaptionAboveMedia: true,
		HasSpoiler:            true,
		ReplyParameters:       &telegram.ReplyParameters{MessageID: 20},
		ReplyMarkup:           telegram.NewInlineKeyboard([]telegram.InlineKeyboardButton{telegram.InlineButtonCallback("OK", "ok")}),
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if message == nil || message.Animation == nil || message.Animation.FileID != "sent-animation" || message.Animation.Duration != 10 {
		t.Fatalf("unexpected message: %+v", message)
	}
}

func TestSendAnimationSendsJSONFileURL(t *testing.T) {
	const token = "123:secret"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/bot"+token+"/sendAnimation" {
			t.Fatalf("unexpected path: %q", r.URL.Path)
		}
		var payload map[string]any
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode payload: %v", err)
		}
		if got := payload["animation"]; got != "https://example.com/animation.gif" {
			t.Fatalf("unexpected animation URL: %#v", got)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true,"result":{"message_id":9,"chat":{"id":12345,"type":"private"},"date":104,"animation":{"file_id":"url-animation","file_unique_id":"url-animation-unique","width":320,"height":240,"duration":10}}}`))
	}))
	defer server.Close()

	bot := newTestBot(t, token, server.URL, server.Client())
	message, err := bot.SendAnimation(context.Background(), SendAnimationParams{ChatID: ChatIDInt(12345), Animation: FileURL("https://example.com/animation.gif")})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if message == nil || message.Animation == nil || message.Animation.FileID != "url-animation" {
		t.Fatalf("unexpected message: %+v", message)
	}
}

func TestSendAnimationMultipartUpload(t *testing.T) {
	const token = "123:secret"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/bot"+token+"/sendAnimation" {
			t.Fatalf("unexpected path: %q", r.URL.Path)
		}
		if got := r.Header.Get("Content-Type"); !strings.HasPrefix(got, "multipart/form-data") {
			t.Fatalf("unexpected content type: %q", got)
		}
		if err := r.ParseMultipartForm(4096); err != nil {
			t.Fatalf("parse multipart: %v", err)
		}
		assertMultipartValue(t, r, "chat_id", "12345")
		assertMultipartValue(t, r, "message_thread_id", "9")
		assertMultipartValue(t, r, "animation", "attach://animation")
		assertMultipartValue(t, r, "thumbnail", "attach://thumbnail")
		assertMultipartValue(t, r, "duration", "11")
		assertMultipartValue(t, r, "width", "640")
		assertMultipartValue(t, r, "height", "480")
		assertMultipartValue(t, r, "caption", "caption")
		assertMultipartValue(t, r, "show_caption_above_media", "true")
		assertMultipartValue(t, r, "has_spoiler", "true")
		content, header := readMultipartFile(t, r, "animation")
		if header.Filename != "animation.gif" {
			t.Fatalf("unexpected animation filename: %q", header.Filename)
		}
		if string(content) != "animation-data" {
			t.Fatalf("unexpected animation content: %q", content)
		}
		thumbContent, thumbHeader := readMultipartFile(t, r, "thumbnail")
		if thumbHeader.Filename != "thumb.jpg" {
			t.Fatalf("unexpected thumbnail filename: %q", thumbHeader.Filename)
		}
		if string(thumbContent) != "thumbnail-data" {
			t.Fatalf("unexpected thumbnail content: %q", thumbContent)
		}
		var reply telegram.ReplyParameters
		if err := json.Unmarshal([]byte(r.MultipartForm.Value["reply_parameters"][0]), &reply); err != nil {
			t.Fatalf("decode reply_parameters: %v", err)
		}
		if reply.MessageID != 21 {
			t.Fatalf("unexpected reply parameters: %+v", reply)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true,"result":{"message_id":10,"chat":{"id":12345,"type":"private"},"date":105,"animation":{"file_id":"uploaded-animation","file_unique_id":"uploaded-animation-unique","width":640,"height":480,"duration":11}}}`))
	}))
	defer server.Close()

	bot := newTestBot(t, token, server.URL, server.Client())
	message, err := bot.SendAnimation(context.Background(), SendAnimationParams{
		ChatID:                ChatIDInt(12345),
		MessageThreadID:       9,
		Animation:             FileUpload(UploadFile{Name: "animation.gif", Reader: strings.NewReader("animation-data"), ContentType: "image/gif"}),
		Thumbnail:             FileUpload(UploadFile{Name: "thumb.jpg", Reader: strings.NewReader("thumbnail-data"), ContentType: "image/jpeg"}),
		Duration:              11,
		Width:                 640,
		Height:                480,
		Caption:               "caption",
		ShowCaptionAboveMedia: true,
		HasSpoiler:            true,
		ReplyParameters:       &telegram.ReplyParameters{MessageID: 21},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if message == nil || message.Animation == nil || message.Animation.FileID != "uploaded-animation" {
		t.Fatalf("unexpected message: %+v", message)
	}
}

func TestSendAnimationValidation(t *testing.T) {
	const token = "123:secret"
	valid := SendAnimationParams{ChatID: ChatIDInt(12345), Animation: FileID("animation")}
	tests := []struct {
		name   string
		mutate func(*SendAnimationParams)
	}{
		{name: "empty chat", mutate: func(p *SendAnimationParams) { p.ChatID = ChatID{} }},
		{name: "empty animation", mutate: func(p *SendAnimationParams) { p.Animation = FileRef{} }},
		{name: "negative duration", mutate: func(p *SendAnimationParams) { p.Duration = -1 }},
		{name: "negative width", mutate: func(p *SendAnimationParams) { p.Width = -1 }},
		{name: "negative height", mutate: func(p *SendAnimationParams) { p.Height = -1 }},
		{name: "negative thread", mutate: func(p *SendAnimationParams) { p.MessageThreadID = -1 }},
		{name: "invalid thumbnail", mutate: func(p *SendAnimationParams) { p.Thumbnail = FileUpload(UploadFile{Name: "thumb.jpg"}) }},
		{name: "parse mode and entities", mutate: func(p *SendAnimationParams) {
			p.ParseMode = "HTML"
			p.CaptionEntities = []telegram.MessageEntity{{Type: telegram.EntityBold, Offset: 0, Length: 1}}
		}},
		{name: "invalid reply parameters", mutate: func(p *SendAnimationParams) { p.ReplyParameters = &telegram.ReplyParameters{} }},
		{name: "invalid reply markup", mutate: func(p *SendAnimationParams) { p.ReplyMarkup = telegram.InlineKeyboardMarkup{} }},
	}

	bot := newTestBot(t, token, "https://example.test", nil)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			params := valid
			tt.mutate(&params)
			message, err := bot.SendAnimation(context.Background(), params)
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

func TestSendAnimationAPIAndTransportErrors(t *testing.T) {
	valid := SendAnimationParams{ChatID: ChatIDInt(12345), Animation: FileID("animation")}
	testSendMethodErrorCases(t, "sendAnimation", func(bot *Bot) (*telegram.Message, error) {
		return bot.SendAnimation(context.Background(), valid)
	}, func(bot *Bot, ctx context.Context) (*telegram.Message, error) {
		return bot.SendAnimation(ctx, valid)
	})
}

func TestSendVideoNoteSendsJSONFileIDAndDecodesMessage(t *testing.T) {
	const token = "123:secret"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/bot"+token+"/sendVideoNote" {
			t.Fatalf("unexpected path: %q", r.URL.Path)
		}
		var payload map[string]any
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode payload: %v", err)
		}
		if got := payload["video_note"]; got != "video-note-file-id" {
			t.Fatalf("unexpected video_note: %#v", got)
		}
		if got := payload["thumbnail"]; got != "thumbnail-file-id" {
			t.Fatalf("unexpected thumbnail: %#v", got)
		}
		if got := payload["duration"]; got != float64(12) {
			t.Fatalf("unexpected duration: %#v", got)
		}
		if got := payload["length"]; got != float64(240) {
			t.Fatalf("unexpected length: %#v", got)
		}
		reply := payload["reply_parameters"].(map[string]any)
		if got := reply["message_id"]; got != float64(22) {
			t.Fatalf("unexpected reply message_id: %#v", got)
		}
		markup := payload["reply_markup"].(map[string]any)
		if _, ok := markup["inline_keyboard"]; !ok {
			t.Fatalf("reply_markup.inline_keyboard missing: %#v", markup)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true,"result":{"message_id":11,"chat":{"id":12345,"type":"private"},"date":106,"video_note":{"file_id":"sent-video-note","file_unique_id":"video-note-unique","length":240,"duration":12}}}`))
	}))
	defer server.Close()

	bot := newTestBot(t, token, server.URL, server.Client())
	message, err := bot.SendVideoNote(context.Background(), SendVideoNoteParams{
		ChatID:          ChatIDInt(12345),
		VideoNote:       FileID("video-note-file-id"),
		Thumbnail:       FileID("thumbnail-file-id"),
		Duration:        12,
		Length:          240,
		ReplyParameters: &telegram.ReplyParameters{MessageID: 22},
		ReplyMarkup:     telegram.NewInlineKeyboard([]telegram.InlineKeyboardButton{telegram.InlineButtonCallback("OK", "ok")}),
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if message == nil || message.VideoNote == nil || message.VideoNote.FileID != "sent-video-note" || message.VideoNote.Length != 240 {
		t.Fatalf("unexpected message: %+v", message)
	}
}

func TestSendVideoNoteMultipartUpload(t *testing.T) {
	const token = "123:secret"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/bot"+token+"/sendVideoNote" {
			t.Fatalf("unexpected path: %q", r.URL.Path)
		}
		if got := r.Header.Get("Content-Type"); !strings.HasPrefix(got, "multipart/form-data") {
			t.Fatalf("unexpected content type: %q", got)
		}
		if err := r.ParseMultipartForm(4096); err != nil {
			t.Fatalf("parse multipart: %v", err)
		}
		assertMultipartValue(t, r, "chat_id", "12345")
		assertMultipartValue(t, r, "message_thread_id", "10")
		assertMultipartValue(t, r, "video_note", "attach://video_note")
		assertMultipartValue(t, r, "thumbnail", "attach://thumbnail")
		assertMultipartValue(t, r, "duration", "13")
		assertMultipartValue(t, r, "length", "320")
		content, header := readMultipartFile(t, r, "video_note")
		if header.Filename != "video-note.mp4" {
			t.Fatalf("unexpected video_note filename: %q", header.Filename)
		}
		if string(content) != "video-note-data" {
			t.Fatalf("unexpected video_note content: %q", content)
		}
		thumbContent, thumbHeader := readMultipartFile(t, r, "thumbnail")
		if thumbHeader.Filename != "thumb.jpg" {
			t.Fatalf("unexpected thumbnail filename: %q", thumbHeader.Filename)
		}
		if string(thumbContent) != "thumbnail-data" {
			t.Fatalf("unexpected thumbnail content: %q", thumbContent)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true,"result":{"message_id":12,"chat":{"id":12345,"type":"private"},"date":107,"video_note":{"file_id":"uploaded-video-note","file_unique_id":"uploaded-video-note-unique","length":320,"duration":13}}}`))
	}))
	defer server.Close()

	bot := newTestBot(t, token, server.URL, server.Client())
	message, err := bot.SendVideoNote(context.Background(), SendVideoNoteParams{
		ChatID:          ChatIDInt(12345),
		MessageThreadID: 10,
		VideoNote:       FileUpload(UploadFile{Name: "video-note.mp4", Reader: strings.NewReader("video-note-data"), ContentType: "video/mp4"}),
		Thumbnail:       FileUpload(UploadFile{Name: "thumb.jpg", Reader: strings.NewReader("thumbnail-data"), ContentType: "image/jpeg"}),
		Duration:        13,
		Length:          320,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if message == nil || message.VideoNote == nil || message.VideoNote.FileID != "uploaded-video-note" {
		t.Fatalf("unexpected message: %+v", message)
	}
}

func TestSendVideoNoteValidation(t *testing.T) {
	const token = "123:secret"
	valid := SendVideoNoteParams{ChatID: ChatIDInt(12345), VideoNote: FileID("video-note")}
	tests := []struct {
		name   string
		mutate func(*SendVideoNoteParams)
	}{
		{name: "empty chat", mutate: func(p *SendVideoNoteParams) { p.ChatID = ChatID{} }},
		{name: "empty video note", mutate: func(p *SendVideoNoteParams) { p.VideoNote = FileRef{} }},
		{name: "file url unsupported", mutate: func(p *SendVideoNoteParams) { p.VideoNote = FileURL("https://example.com/video-note.mp4") }},
		{name: "negative duration", mutate: func(p *SendVideoNoteParams) { p.Duration = -1 }},
		{name: "negative length", mutate: func(p *SendVideoNoteParams) { p.Length = -1 }},
		{name: "negative thread", mutate: func(p *SendVideoNoteParams) { p.MessageThreadID = -1 }},
		{name: "invalid thumbnail", mutate: func(p *SendVideoNoteParams) { p.Thumbnail = FileUpload(UploadFile{Name: "thumb.jpg"}) }},
		{name: "invalid reply parameters", mutate: func(p *SendVideoNoteParams) { p.ReplyParameters = &telegram.ReplyParameters{} }},
		{name: "invalid reply markup", mutate: func(p *SendVideoNoteParams) { p.ReplyMarkup = telegram.InlineKeyboardMarkup{} }},
	}

	bot := newTestBot(t, token, "https://example.test", nil)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			params := valid
			tt.mutate(&params)
			message, err := bot.SendVideoNote(context.Background(), params)
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

func TestSendVideoNoteAPIAndTransportErrors(t *testing.T) {
	valid := SendVideoNoteParams{ChatID: ChatIDInt(12345), VideoNote: FileID("video-note")}
	testSendMethodErrorCases(t, "sendVideoNote", func(bot *Bot) (*telegram.Message, error) {
		return bot.SendVideoNote(context.Background(), valid)
	}, func(bot *Bot, ctx context.Context) (*telegram.Message, error) {
		return bot.SendVideoNote(ctx, valid)
	})
}
