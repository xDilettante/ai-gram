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

func TestGetStickerSetSendsPayloadAndDecodesResult(t *testing.T) {
	const token = "123:secret"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/bot"+token+"/getStickerSet" {
			t.Fatalf("unexpected path: %q", r.URL.Path)
		}
		var payload map[string]any
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode payload: %v", err)
		}
		if payload["name"] != "animals_by_bot" {
			t.Fatalf("unexpected payload: %#v", payload)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true,"result":{"name":"animals_by_bot","title":"Animals","sticker_type":"regular","stickers":[{"file_id":"sticker-file","file_unique_id":"unique","type":"regular","width":512,"height":512,"is_animated":false,"is_video":false,"custom_emoji_id":"custom","needs_repainting":true,"mask_position":{"point":"eyes","x_shift":0.1,"y_shift":-0.2,"scale":1.5}}],"thumbnail":{"file_id":"thumb","file_unique_id":"thumb-unique","width":100,"height":100}}}`))
	}))
	defer server.Close()

	bot := newTestBot(t, token, server.URL, server.Client())
	set, err := bot.GetStickerSet(context.Background(), GetStickerSetParams{Name: "animals_by_bot"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if set == nil || set.Name != "animals_by_bot" || len(set.Stickers) != 1 || set.Stickers[0].MaskPosition == nil || !set.Stickers[0].NeedsRepainting {
		t.Fatalf("unexpected sticker set: %+v", set)
	}
}

func TestGetCustomEmojiStickersSendsPayloadAndDecodesResult(t *testing.T) {
	const token = "123:secret"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/bot"+token+"/getCustomEmojiStickers" {
			t.Fatalf("unexpected path: %q", r.URL.Path)
		}
		var payload map[string]any
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode payload: %v", err)
		}
		ids := payload["custom_emoji_ids"].([]any)
		if len(ids) != 2 || ids[0] != "emoji-1" || ids[1] != "emoji-2" {
			t.Fatalf("unexpected payload: %#v", payload)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true,"result":[{"file_id":"custom-file","file_unique_id":"custom-unique","type":"custom_emoji","width":512,"height":512,"is_animated":false,"is_video":false,"custom_emoji_id":"emoji-1"}]}`))
	}))
	defer server.Close()

	bot := newTestBot(t, token, server.URL, server.Client())
	stickers, err := bot.GetCustomEmojiStickers(context.Background(), GetCustomEmojiStickersParams{CustomEmojiIDs: []string{"emoji-1", "emoji-2"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(stickers) != 1 || stickers[0].CustomEmojiID != "emoji-1" {
		t.Fatalf("unexpected stickers: %+v", stickers)
	}
}

func TestUploadStickerFileSendsMultipartAndDecodesFile(t *testing.T) {
	const token = "123:secret"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/bot"+token+"/uploadStickerFile" {
			t.Fatalf("unexpected path: %q", r.URL.Path)
		}
		if got := r.Header.Get("Content-Type"); !strings.HasPrefix(got, "multipart/form-data") {
			t.Fatalf("unexpected content type: %q", got)
		}
		if err := r.ParseMultipartForm(2048); err != nil {
			t.Fatalf("parse multipart: %v", err)
		}
		assertMultipartValue(t, r, "user_id", "123")
		assertMultipartValue(t, r, "sticker_format", "static")
		content, header := readMultipartFile(t, r, "sticker")
		if header.Filename != "sticker.webp" || string(content) != "sticker-data" {
			t.Fatalf("unexpected upload: filename=%q content=%q", header.Filename, content)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true,"result":{"file_id":"uploaded-file","file_unique_id":"uploaded-unique","file_size":123}}`))
	}))
	defer server.Close()

	bot := newTestBot(t, token, server.URL, server.Client())
	file, err := bot.UploadStickerFile(context.Background(), UploadStickerFileParams{UserID: 123, Sticker: FileUpload(UploadFile{Name: "sticker.webp", Reader: strings.NewReader("sticker-data"), ContentType: "image/webp"}), StickerFormat: "static"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if file == nil || file.FileID != "uploaded-file" || file.FileSize != 123 {
		t.Fatalf("unexpected file: %+v", file)
	}
}

func TestCreateNewStickerSetSendsJSON(t *testing.T) {
	testStickerBoolJSONSuccess(t, "createNewStickerSet", func(bot *Bot) (bool, error) {
		return bot.CreateNewStickerSet(context.Background(), CreateNewStickerSetParams{UserID: 123, Name: "animals_by_bot", Title: "Animals", StickerType: "regular", NeedsRepainting: true, Stickers: []InputSticker{NewInputSticker(FileID("sticker-file"), "static", "😀")}})
	}, func(t *testing.T, payload map[string]any) {
		if payload["user_id"] != float64(123) || payload["name"] != "animals_by_bot" || payload["title"] != "Animals" || payload["sticker_type"] != "regular" || payload["needs_repainting"] != true {
			t.Fatalf("unexpected payload: %#v", payload)
		}
		stickers := payload["stickers"].([]any)
		sticker := stickers[0].(map[string]any)
		if sticker["sticker"] != "sticker-file" || sticker["format"] != "static" {
			t.Fatalf("unexpected sticker payload: %#v", sticker)
		}
	})
}

func TestCreateNewStickerSetSendsMultipartUpload(t *testing.T) {
	const token = "123:secret"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/bot"+token+"/createNewStickerSet" {
			t.Fatalf("unexpected path: %q", r.URL.Path)
		}
		if err := r.ParseMultipartForm(2048); err != nil {
			t.Fatalf("parse multipart: %v", err)
		}
		assertMultipartValue(t, r, "user_id", "123")
		assertMultipartValue(t, r, "name", "animals_by_bot")
		assertMultipartValue(t, r, "title", "Animals")
		var stickers []map[string]any
		if err := json.Unmarshal([]byte(r.MultipartForm.Value["stickers"][0]), &stickers); err != nil {
			t.Fatalf("decode stickers: %v", err)
		}
		if len(stickers) != 1 || stickers[0]["sticker"] != "attach://sticker0" || stickers[0]["format"] != "static" {
			t.Fatalf("unexpected stickers payload: %#v", stickers)
		}
		content, header := readMultipartFile(t, r, "sticker0")
		if header.Filename != "sticker.png" || string(content) != "png-data" {
			t.Fatalf("unexpected upload: filename=%q content=%q", header.Filename, content)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true,"result":true}`))
	}))
	defer server.Close()

	bot := newTestBot(t, token, server.URL, server.Client())
	ok, err := bot.CreateNewStickerSet(context.Background(), CreateNewStickerSetParams{UserID: 123, Name: "animals_by_bot", Title: "Animals", Stickers: []InputSticker{NewInputSticker(FileUpload(UploadFile{Name: "sticker.png", Reader: strings.NewReader("png-data")}), "static", "😀")}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Fatal("expected true result")
	}
}

func TestAddStickerToSetSendsMultipartUpload(t *testing.T) {
	testStickerMultipartStickerSuccess(t, "addStickerToSet", "sticker", func(bot *Bot) (bool, error) {
		return bot.AddStickerToSet(context.Background(), AddStickerToSetParams{UserID: 123, Name: "animals_by_bot", Sticker: NewInputSticker(FileUpload(UploadFile{Name: "sticker.webp", Reader: strings.NewReader("sticker-data")}), "static", "😀")})
	})
}

func TestReplaceStickerInSetSendsMultipartUpload(t *testing.T) {
	testStickerMultipartStickerSuccess(t, "replaceStickerInSet", "sticker", func(bot *Bot) (bool, error) {
		return bot.ReplaceStickerInSet(context.Background(), ReplaceStickerInSetParams{UserID: 123, Name: "animals_by_bot", OldSticker: "old-sticker", Sticker: NewInputSticker(FileUpload(UploadFile{Name: "sticker.webp", Reader: strings.NewReader("sticker-data")}), "static", "😀")})
	})
}

func TestSimpleStickerMutationMethodsSendJSON(t *testing.T) {
	tests := []struct {
		name   string
		method string
		call   func(*Bot) (bool, error)
		check  func(*testing.T, map[string]any)
	}{
		{name: "set position", method: "setStickerPositionInSet", call: func(bot *Bot) (bool, error) {
			return bot.SetStickerPositionInSet(context.Background(), SetStickerPositionInSetParams{Sticker: "sticker", Position: 2})
		}, check: func(t *testing.T, payload map[string]any) {
			if payload["sticker"] != "sticker" || payload["position"] != float64(2) {
				t.Fatalf("unexpected payload: %#v", payload)
			}
		}},
		{name: "delete sticker", method: "deleteStickerFromSet", call: func(bot *Bot) (bool, error) {
			return bot.DeleteStickerFromSet(context.Background(), DeleteStickerFromSetParams{Sticker: "sticker"})
		}, check: requireStickerPayload},
		{name: "set emoji list", method: "setStickerEmojiList", call: func(bot *Bot) (bool, error) {
			return bot.SetStickerEmojiList(context.Background(), SetStickerEmojiListParams{Sticker: "sticker", EmojiList: []string{"😀", "🙂"}})
		}, check: func(t *testing.T, payload map[string]any) {
			if payload["sticker"] != "sticker" || len(payload["emoji_list"].([]any)) != 2 {
				t.Fatalf("unexpected payload: %#v", payload)
			}
		}},
		{name: "set keywords", method: "setStickerKeywords", call: func(bot *Bot) (bool, error) {
			return bot.SetStickerKeywords(context.Background(), SetStickerKeywordsParams{Sticker: "sticker", Keywords: []string{"cat", "animal"}})
		}, check: func(t *testing.T, payload map[string]any) {
			if payload["sticker"] != "sticker" || len(payload["keywords"].([]any)) != 2 {
				t.Fatalf("unexpected payload: %#v", payload)
			}
		}},
		{name: "clear keywords", method: "setStickerKeywords", call: func(bot *Bot) (bool, error) {
			return bot.SetStickerKeywords(context.Background(), SetStickerKeywordsParams{Sticker: "sticker", Keywords: []string{}})
		}, check: func(t *testing.T, payload map[string]any) {
			if payload["sticker"] != "sticker" || len(payload["keywords"].([]any)) != 0 {
				t.Fatalf("unexpected payload: %#v", payload)
			}
		}},
		{name: "set mask position", method: "setStickerMaskPosition", call: func(bot *Bot) (bool, error) {
			return bot.SetStickerMaskPosition(context.Background(), SetStickerMaskPositionParams{Sticker: "sticker", MaskPosition: &telegram.MaskPosition{Point: "eyes", XShift: 0.1, YShift: -0.2, Scale: 1.5}})
		}, check: func(t *testing.T, payload map[string]any) {
			mask := payload["mask_position"].(map[string]any)
			if payload["sticker"] != "sticker" || mask["point"] != "eyes" || mask["scale"] != 1.5 {
				t.Fatalf("unexpected payload: %#v", payload)
			}
		}},
		{name: "set sticker set title", method: "setStickerSetTitle", call: func(bot *Bot) (bool, error) {
			return bot.SetStickerSetTitle(context.Background(), SetStickerSetTitleParams{Name: "animals_by_bot", Title: "Animals"})
		}, check: func(t *testing.T, payload map[string]any) {
			if payload["name"] != "animals_by_bot" || payload["title"] != "Animals" {
				t.Fatalf("unexpected payload: %#v", payload)
			}
		}},
		{name: "set custom emoji thumbnail", method: "setCustomEmojiStickerSetThumbnail", call: func(bot *Bot) (bool, error) {
			return bot.SetCustomEmojiStickerSetThumbnail(context.Background(), SetCustomEmojiStickerSetThumbnailParams{Name: "emoji_by_bot", CustomEmojiID: "emoji-id"})
		}, check: func(t *testing.T, payload map[string]any) {
			if payload["name"] != "emoji_by_bot" || payload["custom_emoji_id"] != "emoji-id" {
				t.Fatalf("unexpected payload: %#v", payload)
			}
		}},
		{name: "clear custom emoji thumbnail", method: "setCustomEmojiStickerSetThumbnail", call: func(bot *Bot) (bool, error) {
			return bot.SetCustomEmojiStickerSetThumbnail(context.Background(), SetCustomEmojiStickerSetThumbnailParams{Name: "emoji_by_bot"})
		}, check: func(t *testing.T, payload map[string]any) {
			if payload["name"] != "emoji_by_bot" || payload["custom_emoji_id"] != "" {
				t.Fatalf("unexpected payload: %#v", payload)
			}
		}},
		{name: "delete sticker set", method: "deleteStickerSet", call: func(bot *Bot) (bool, error) {
			return bot.DeleteStickerSet(context.Background(), DeleteStickerSetParams{Name: "animals_by_bot"})
		}, check: func(t *testing.T, payload map[string]any) {
			if payload["name"] != "animals_by_bot" {
				t.Fatalf("unexpected payload: %#v", payload)
			}
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testStickerBoolJSONSuccess(t, tt.method, tt.call, tt.check)
		})
	}
}

func TestSetStickerSetThumbnailJSONAndMultipart(t *testing.T) {
	testStickerBoolJSONSuccess(t, "setStickerSetThumbnail", func(bot *Bot) (bool, error) {
		return bot.SetStickerSetThumbnail(context.Background(), SetStickerSetThumbnailParams{Name: "animals_by_bot", UserID: 123, Thumbnail: FileID("thumb-file"), Format: "static"})
	}, func(t *testing.T, payload map[string]any) {
		if payload["name"] != "animals_by_bot" || payload["user_id"] != float64(123) || payload["thumbnail"] != "thumb-file" || payload["format"] != "static" {
			t.Fatalf("unexpected payload: %#v", payload)
		}
	})

	const token = "123:secret"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/bot"+token+"/setStickerSetThumbnail" {
			t.Fatalf("unexpected path: %q", r.URL.Path)
		}
		if err := r.ParseMultipartForm(2048); err != nil {
			t.Fatalf("parse multipart: %v", err)
		}
		assertMultipartValue(t, r, "name", "animals_by_bot")
		assertMultipartValue(t, r, "user_id", "123")
		assertMultipartValue(t, r, "format", "static")
		assertMultipartValue(t, r, "thumbnail", "attach://thumbnail")
		content, header := readMultipartFile(t, r, "thumbnail")
		if header.Filename != "thumb.webp" || string(content) != "thumb-data" {
			t.Fatalf("unexpected upload: filename=%q content=%q", header.Filename, content)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true,"result":true}`))
	}))
	defer server.Close()

	bot := newTestBot(t, token, server.URL, server.Client())
	ok, err := bot.SetStickerSetThumbnail(context.Background(), SetStickerSetThumbnailParams{Name: "animals_by_bot", UserID: 123, Format: "static", Thumbnail: FileUpload(UploadFile{Name: "thumb.webp", Reader: strings.NewReader("thumb-data")})})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Fatal("expected true result")
	}
}

func TestStickerSetMethodValidation(t *testing.T) {
	const token = "123:secret"
	bot := newTestBot(t, token, "https://example.test", nil)
	tests := []struct {
		name string
		call func() error
	}{
		{name: "get sticker set empty name", call: func() error { _, err := bot.GetStickerSet(context.Background(), GetStickerSetParams{}); return err }},
		{name: "custom emoji empty list", call: func() error {
			_, err := bot.GetCustomEmojiStickers(context.Background(), GetCustomEmojiStickersParams{})
			return err
		}},
		{name: "custom emoji empty id", call: func() error {
			_, err := bot.GetCustomEmojiStickers(context.Background(), GetCustomEmojiStickersParams{CustomEmojiIDs: []string{""}})
			return err
		}},
		{name: "upload user", call: func() error {
			_, err := bot.UploadStickerFile(context.Background(), UploadStickerFileParams{Sticker: FileUpload(UploadFile{Name: "s.webp", Reader: strings.NewReader("s")}), StickerFormat: "static"})
			return err
		}},
		{name: "upload file id", call: func() error {
			_, err := bot.UploadStickerFile(context.Background(), UploadStickerFileParams{UserID: 1, Sticker: FileID("s"), StickerFormat: "static"})
			return err
		}},
		{name: "upload format", call: func() error {
			_, err := bot.UploadStickerFile(context.Background(), UploadStickerFileParams{UserID: 1, Sticker: FileUpload(UploadFile{Name: "s.webp", Reader: strings.NewReader("s")})})
			return err
		}},
		{name: "create name", call: func() error {
			_, err := bot.CreateNewStickerSet(context.Background(), CreateNewStickerSetParams{UserID: 1, Title: "T", Stickers: []InputSticker{NewInputSticker(FileID("s"), "static", "😀")}})
			return err
		}},
		{name: "create stickers", call: func() error {
			_, err := bot.CreateNewStickerSet(context.Background(), CreateNewStickerSetParams{UserID: 1, Name: "set", Title: "T"})
			return err
		}},
		{name: "input emoji", call: func() error {
			_, err := bot.AddStickerToSet(context.Background(), AddStickerToSetParams{UserID: 1, Name: "set", Sticker: NewInputSticker(FileID("s"), "static")})
			return err
		}},
		{name: "input animated URL", call: func() error {
			_, err := bot.AddStickerToSet(context.Background(), AddStickerToSetParams{UserID: 1, Name: "set", Sticker: NewInputSticker(FileURL("https://example.com/s.tgs"), "animated", "😀")})
			return err
		}},
		{name: "replace old sticker", call: func() error {
			_, err := bot.ReplaceStickerInSet(context.Background(), ReplaceStickerInSetParams{UserID: 1, Name: "set", Sticker: NewInputSticker(FileID("s"), "static", "😀")})
			return err
		}},
		{name: "position negative", call: func() error {
			_, err := bot.SetStickerPositionInSet(context.Background(), SetStickerPositionInSetParams{Sticker: "s", Position: -1})
			return err
		}},
		{name: "emoji list empty", call: func() error {
			_, err := bot.SetStickerEmojiList(context.Background(), SetStickerEmojiListParams{Sticker: "s"})
			return err
		}},
		{name: "keyword empty", call: func() error {
			_, err := bot.SetStickerKeywords(context.Background(), SetStickerKeywordsParams{Sticker: "s", Keywords: []string{""}})
			return err
		}},
		{name: "mask point", call: func() error {
			_, err := bot.SetStickerMaskPosition(context.Background(), SetStickerMaskPositionParams{Sticker: "s", MaskPosition: &telegram.MaskPosition{Scale: 1}})
			return err
		}},
		{name: "mask scale", call: func() error {
			_, err := bot.SetStickerMaskPosition(context.Background(), SetStickerMaskPositionParams{Sticker: "s", MaskPosition: &telegram.MaskPosition{Point: "eyes"}})
			return err
		}},
		{name: "title empty", call: func() error {
			_, err := bot.SetStickerSetTitle(context.Background(), SetStickerSetTitleParams{Name: "set"})
			return err
		}},
		{name: "thumbnail user", call: func() error {
			_, err := bot.SetStickerSetThumbnail(context.Background(), SetStickerSetThumbnailParams{Name: "set", Format: "static"})
			return err
		}},
		{name: "thumbnail format", call: func() error {
			_, err := bot.SetStickerSetThumbnail(context.Background(), SetStickerSetThumbnailParams{Name: "set", UserID: 1})
			return err
		}},
		{name: "animated thumbnail URL", call: func() error {
			_, err := bot.SetStickerSetThumbnail(context.Background(), SetStickerSetThumbnailParams{Name: "set", UserID: 1, Format: "animated", Thumbnail: FileURL("https://example.com/thumb.tgs")})
			return err
		}},
		{name: "custom thumbnail name", call: func() error {
			_, err := bot.SetCustomEmojiStickerSetThumbnail(context.Background(), SetCustomEmojiStickerSetThumbnailParams{})
			return err
		}},
		{name: "delete set name", call: func() error {
			_, err := bot.DeleteStickerSet(context.Background(), DeleteStickerSetParams{})
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

func TestStickerSetBoolMethodsAPIAndTransportErrors(t *testing.T) {
	tests := []struct {
		method          string
		call            func(*Bot) (bool, error)
		callWithContext func(*Bot, context.Context) (bool, error)
	}{
		{method: "createNewStickerSet", call: func(bot *Bot) (bool, error) {
			return bot.CreateNewStickerSet(context.Background(), validCreateStickerSet())
		}, callWithContext: func(bot *Bot, ctx context.Context) (bool, error) {
			return bot.CreateNewStickerSet(ctx, validCreateStickerSet())
		}},
		{method: "addStickerToSet", call: func(bot *Bot) (bool, error) { return bot.AddStickerToSet(context.Background(), validAddStickerToSet()) }, callWithContext: func(bot *Bot, ctx context.Context) (bool, error) {
			return bot.AddStickerToSet(ctx, validAddStickerToSet())
		}},
		{method: "replaceStickerInSet", call: func(bot *Bot) (bool, error) {
			return bot.ReplaceStickerInSet(context.Background(), validReplaceStickerInSet())
		}, callWithContext: func(bot *Bot, ctx context.Context) (bool, error) {
			return bot.ReplaceStickerInSet(ctx, validReplaceStickerInSet())
		}},
		{method: "setStickerPositionInSet", call: func(bot *Bot) (bool, error) {
			return bot.SetStickerPositionInSet(context.Background(), SetStickerPositionInSetParams{Sticker: "s", Position: 1})
		}, callWithContext: func(bot *Bot, ctx context.Context) (bool, error) {
			return bot.SetStickerPositionInSet(ctx, SetStickerPositionInSetParams{Sticker: "s", Position: 1})
		}},
		{method: "deleteStickerFromSet", call: func(bot *Bot) (bool, error) {
			return bot.DeleteStickerFromSet(context.Background(), DeleteStickerFromSetParams{Sticker: "s"})
		}, callWithContext: func(bot *Bot, ctx context.Context) (bool, error) {
			return bot.DeleteStickerFromSet(ctx, DeleteStickerFromSetParams{Sticker: "s"})
		}},
		{method: "setStickerEmojiList", call: func(bot *Bot) (bool, error) {
			return bot.SetStickerEmojiList(context.Background(), SetStickerEmojiListParams{Sticker: "s", EmojiList: []string{"😀"}})
		}, callWithContext: func(bot *Bot, ctx context.Context) (bool, error) {
			return bot.SetStickerEmojiList(ctx, SetStickerEmojiListParams{Sticker: "s", EmojiList: []string{"😀"}})
		}},
		{method: "setStickerKeywords", call: func(bot *Bot) (bool, error) {
			return bot.SetStickerKeywords(context.Background(), SetStickerKeywordsParams{Sticker: "s", Keywords: []string{"cat"}})
		}, callWithContext: func(bot *Bot, ctx context.Context) (bool, error) {
			return bot.SetStickerKeywords(ctx, SetStickerKeywordsParams{Sticker: "s", Keywords: []string{"cat"}})
		}},
		{method: "setStickerMaskPosition", call: func(bot *Bot) (bool, error) {
			return bot.SetStickerMaskPosition(context.Background(), SetStickerMaskPositionParams{Sticker: "s"})
		}, callWithContext: func(bot *Bot, ctx context.Context) (bool, error) {
			return bot.SetStickerMaskPosition(ctx, SetStickerMaskPositionParams{Sticker: "s"})
		}},
		{method: "setStickerSetTitle", call: func(bot *Bot) (bool, error) {
			return bot.SetStickerSetTitle(context.Background(), SetStickerSetTitleParams{Name: "set", Title: "Set"})
		}, callWithContext: func(bot *Bot, ctx context.Context) (bool, error) {
			return bot.SetStickerSetTitle(ctx, SetStickerSetTitleParams{Name: "set", Title: "Set"})
		}},
		{method: "setStickerSetThumbnail", call: func(bot *Bot) (bool, error) {
			return bot.SetStickerSetThumbnail(context.Background(), SetStickerSetThumbnailParams{Name: "set", UserID: 1, Format: "static", Thumbnail: FileID("thumb")})
		}, callWithContext: func(bot *Bot, ctx context.Context) (bool, error) {
			return bot.SetStickerSetThumbnail(ctx, SetStickerSetThumbnailParams{Name: "set", UserID: 1, Format: "static", Thumbnail: FileID("thumb")})
		}},
		{method: "setCustomEmojiStickerSetThumbnail", call: func(bot *Bot) (bool, error) {
			return bot.SetCustomEmojiStickerSetThumbnail(context.Background(), SetCustomEmojiStickerSetThumbnailParams{Name: "set", CustomEmojiID: "emoji"})
		}, callWithContext: func(bot *Bot, ctx context.Context) (bool, error) {
			return bot.SetCustomEmojiStickerSetThumbnail(ctx, SetCustomEmojiStickerSetThumbnailParams{Name: "set", CustomEmojiID: "emoji"})
		}},
		{method: "deleteStickerSet", call: func(bot *Bot) (bool, error) {
			return bot.DeleteStickerSet(context.Background(), DeleteStickerSetParams{Name: "set"})
		}, callWithContext: func(bot *Bot, ctx context.Context) (bool, error) {
			return bot.DeleteStickerSet(ctx, DeleteStickerSetParams{Name: "set"})
		}},
	}
	for _, tt := range tests {
		t.Run(tt.method, func(t *testing.T) {
			testBoolMethodErrorCases(t, tt.method, tt.call, tt.callWithContext)
		})
	}
}

func TestStickerSetObjectMethodsAPIAndTransportErrors(t *testing.T) {
	tests := []struct {
		method string
		call   func(context.Context, *Bot) (any, error)
	}{
		{method: "getStickerSet", call: func(ctx context.Context, bot *Bot) (any, error) {
			return bot.GetStickerSet(ctx, GetStickerSetParams{Name: "set"})
		}},
		{method: "getCustomEmojiStickers", call: func(ctx context.Context, bot *Bot) (any, error) {
			return bot.GetCustomEmojiStickers(ctx, GetCustomEmojiStickersParams{CustomEmojiIDs: []string{"emoji"}})
		}},
		{method: "uploadStickerFile", call: func(ctx context.Context, bot *Bot) (any, error) {
			return bot.UploadStickerFile(ctx, UploadStickerFileParams{UserID: 1, Sticker: FileUpload(UploadFile{Name: "s.webp", Reader: strings.NewReader("s")}), StickerFormat: "static"})
		}},
	}
	for _, tt := range tests {
		t.Run(tt.method, func(t *testing.T) {
			testStickerObjectErrors(t, tt.method, tt.call)
		})
	}
}

func testStickerBoolJSONSuccess(t *testing.T, method string, call func(*Bot) (bool, error), check func(*testing.T, map[string]any)) {
	t.Helper()
	const token = "123:secret"
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
		check(t, payload)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true,"result":true}`))
	}))
	defer server.Close()

	bot := newTestBot(t, token, server.URL, server.Client())
	ok, err := call(bot)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Fatal("expected true result")
	}
}

func testStickerMultipartStickerSuccess(t *testing.T, method string, uploadName string, call func(*Bot) (bool, error)) {
	t.Helper()
	const token = "123:secret"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/bot"+token+"/"+method {
			t.Fatalf("unexpected path: %q", r.URL.Path)
		}
		if err := r.ParseMultipartForm(2048); err != nil {
			t.Fatalf("parse multipart: %v", err)
		}
		assertMultipartValue(t, r, "user_id", "123")
		assertMultipartValue(t, r, "name", "animals_by_bot")
		var sticker map[string]any
		if err := json.Unmarshal([]byte(r.MultipartForm.Value["sticker"][0]), &sticker); err != nil {
			t.Fatalf("decode sticker: %v", err)
		}
		if sticker["sticker"] != "attach://"+uploadName || sticker["format"] != "static" {
			t.Fatalf("unexpected sticker payload: %#v", sticker)
		}
		content, header := readMultipartFile(t, r, uploadName)
		if header.Filename != "sticker.webp" || string(content) != "sticker-data" {
			t.Fatalf("unexpected upload: filename=%q content=%q", header.Filename, content)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true,"result":true}`))
	}))
	defer server.Close()

	bot := newTestBot(t, token, server.URL, server.Client())
	ok, err := call(bot)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Fatal("expected true result")
	}
}

func requireStickerPayload(t *testing.T, payload map[string]any) {
	t.Helper()
	if payload["sticker"] != "sticker" {
		t.Fatalf("unexpected payload: %#v", payload)
	}
}

func validCreateStickerSet() CreateNewStickerSetParams {
	return CreateNewStickerSetParams{UserID: 1, Name: "set_by_bot", Title: "Set", Stickers: []InputSticker{NewInputSticker(FileID("sticker"), "static", "😀")}}
}

func validAddStickerToSet() AddStickerToSetParams {
	return AddStickerToSetParams{UserID: 1, Name: "set_by_bot", Sticker: NewInputSticker(FileID("sticker"), "static", "😀")}
}

func validReplaceStickerInSet() ReplaceStickerInSetParams {
	return ReplaceStickerInSetParams{UserID: 1, Name: "set_by_bot", OldSticker: "old", Sticker: NewInputSticker(FileID("sticker"), "static", "😀")}
}

func testStickerObjectErrors(t *testing.T, method string, call func(context.Context, *Bot) (any, error)) {
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
		result, err := call(context.Background(), bot)
		if err == nil {
			t.Fatal("expected error")
		}
		if result != nil && !isNilResult(result) {
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
		_, err := call(context.Background(), bot)
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
		_, err := call(context.Background(), bot)
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
		_, err := call(ctx, bot)
		if err == nil {
			t.Fatal("expected error")
		}
		assertNoToken(t, err, token)
	})
}
