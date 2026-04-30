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

func TestLifecycleMethodsSendNoParamsAndReturnTrue(t *testing.T) {
	const token = "123:secret"
	tests := []struct {
		name   string
		method string
		call   func(context.Context, *Bot) (bool, error)
	}{
		{name: "log out", method: "logOut", call: func(ctx context.Context, bot *Bot) (bool, error) {
			return bot.LogOut(ctx)
		}},
		{name: "close", method: "close", call: func(ctx context.Context, bot *Bot) (bool, error) {
			return bot.Close(ctx)
		}},
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
				if len(payload) != 0 {
					t.Fatalf("unexpected payload: %#v", payload)
				}
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(`{"ok":true,"result":true}`))
			}))
			defer server.Close()

			bot := newTestBot(t, token, server.URL, server.Client())
			ok, err := tt.call(context.Background(), bot)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !ok {
				t.Fatal("expected true result")
			}
		})
	}
}

func TestGetUserProfilePhotosSendsPayloadAndDecodesResult(t *testing.T) {
	const token = "123:secret"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/bot"+token+"/getUserProfilePhotos" {
			t.Fatalf("unexpected path: %q", r.URL.Path)
		}
		var payload map[string]any
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode payload: %v", err)
		}
		if payload["user_id"] != float64(123) || payload["offset"] != float64(2) || payload["limit"] != float64(3) {
			t.Fatalf("unexpected payload: %#v", payload)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true,"result":{"total_count":2,"photos":[[{"file_id":"photo-small","file_unique_id":"photo-small-unique","width":64,"height":64},{"file_id":"photo-large","file_unique_id":"photo-large-unique","width":640,"height":640,"file_size":12345}]]}}`))
	}))
	defer server.Close()

	bot := newTestBot(t, token, server.URL, server.Client())
	photos, err := bot.GetUserProfilePhotos(context.Background(), GetUserProfilePhotosParams{UserID: 123, Offset: 2, Limit: 3})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if photos == nil || photos.TotalCount != 2 || len(photos.Photos) != 1 || len(photos.Photos[0]) != 2 || photos.Photos[0][1].FileSize != 12345 {
		t.Fatalf("unexpected photos: %#v", photos)
	}
}

func TestGetUserProfileAudiosSendsPayloadAndDecodesResult(t *testing.T) {
	const token = "123:secret"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/bot"+token+"/getUserProfileAudios" {
			t.Fatalf("unexpected path: %q", r.URL.Path)
		}
		var payload map[string]any
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode payload: %v", err)
		}
		if payload["user_id"] != float64(123) || payload["offset"] != float64(1) || payload["limit"] != float64(4) {
			t.Fatalf("unexpected payload: %#v", payload)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true,"result":{"total_count":1,"audios":[{"file_id":"audio-file","file_unique_id":"audio-unique","duration":120,"performer":"Artist","title":"Track","file_size":456}]}}`))
	}))
	defer server.Close()

	bot := newTestBot(t, token, server.URL, server.Client())
	audios, err := bot.GetUserProfileAudios(context.Background(), GetUserProfileAudiosParams{UserID: 123, Offset: 1, Limit: 4})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if audios == nil || audios.TotalCount != 1 || len(audios.Audios) != 1 || audios.Audios[0].Title != "Track" || audios.Audios[0].FileSize != 456 {
		t.Fatalf("unexpected audios: %#v", audios)
	}
}

func TestGetForumTopicIconStickersDecodesResult(t *testing.T) {
	const token = "123:secret"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/bot"+token+"/getForumTopicIconStickers" {
			t.Fatalf("unexpected path: %q", r.URL.Path)
		}
		var payload map[string]any
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode payload: %v", err)
		}
		if len(payload) != 0 {
			t.Fatalf("unexpected payload: %#v", payload)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true,"result":[{"file_id":"sticker-file","file_unique_id":"sticker-unique","type":"custom_emoji","width":512,"height":512,"is_animated":false,"is_video":false,"custom_emoji_id":"topic-icon"}]}`))
	}))
	defer server.Close()

	bot := newTestBot(t, token, server.URL, server.Client())
	stickers, err := bot.GetForumTopicIconStickers(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(stickers) != 1 || stickers[0].CustomEmojiID != "topic-icon" {
		t.Fatalf("unexpected stickers: %#v", stickers)
	}
}

func TestProfileReadValidation(t *testing.T) {
	const token = "123:secret"
	bot := newTestBot(t, token, "https://example.test", nil)
	tests := []struct {
		name string
		call func() error
	}{
		{name: "photos user zero", call: func() error {
			_, err := bot.GetUserProfilePhotos(context.Background(), GetUserProfilePhotosParams{})
			return err
		}},
		{name: "photos negative offset", call: func() error {
			_, err := bot.GetUserProfilePhotos(context.Background(), GetUserProfilePhotosParams{UserID: 1, Offset: -1})
			return err
		}},
		{name: "photos negative limit", call: func() error {
			_, err := bot.GetUserProfilePhotos(context.Background(), GetUserProfilePhotosParams{UserID: 1, Limit: -1})
			return err
		}},
		{name: "audios user zero", call: func() error {
			_, err := bot.GetUserProfileAudios(context.Background(), GetUserProfileAudiosParams{})
			return err
		}},
		{name: "audios negative offset", call: func() error {
			_, err := bot.GetUserProfileAudios(context.Background(), GetUserProfileAudiosParams{UserID: 1, Offset: -1})
			return err
		}},
		{name: "audios negative limit", call: func() error {
			_, err := bot.GetUserProfileAudios(context.Background(), GetUserProfileAudiosParams{UserID: 1, Limit: -1})
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

func TestLifecycleAndProfileMethodsReturnAPIError(t *testing.T) {
	const token = "123:secret"
	tests := []struct {
		name   string
		method string
		call   func(*Bot) (any, error)
	}{
		{name: "log out", method: "logOut", call: func(bot *Bot) (any, error) { return bot.LogOut(context.Background()) }},
		{name: "close", method: "close", call: func(bot *Bot) (any, error) { return bot.Close(context.Background()) }},
		{name: "get user profile photos", method: "getUserProfilePhotos", call: func(bot *Bot) (any, error) {
			return bot.GetUserProfilePhotos(context.Background(), GetUserProfilePhotosParams{UserID: 123})
		}},
		{name: "get user profile audios", method: "getUserProfileAudios", call: func(bot *Bot) (any, error) {
			return bot.GetUserProfileAudios(context.Background(), GetUserProfileAudiosParams{UserID: 123})
		}},
		{name: "get forum topic icon stickers", method: "getForumTopicIconStickers", call: func(bot *Bot) (any, error) {
			return bot.GetForumTopicIconStickers(context.Background())
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
			_, err := tt.call(bot)
			if err == nil {
				t.Fatal("expected error")
			}
			var apiErr *apierrors.APIError
			if !stderrors.As(err, &apiErr) {
				t.Fatalf("expected APIError, got %T", err)
			}
			assertNoToken(t, err, token)
		})
	}
}

func TestLifecycleAndProfileMethodsResponseAndContextErrors(t *testing.T) {
	const token = "123:secret"
	tests := []struct {
		name   string
		method string
		call   func(context.Context, *Bot) (any, error)
	}{
		{name: "log out", method: "logOut", call: func(ctx context.Context, bot *Bot) (any, error) { return bot.LogOut(ctx) }},
		{name: "close", method: "close", call: func(ctx context.Context, bot *Bot) (any, error) { return bot.Close(ctx) }},
		{name: "get user profile photos", method: "getUserProfilePhotos", call: func(ctx context.Context, bot *Bot) (any, error) {
			return bot.GetUserProfilePhotos(ctx, GetUserProfilePhotosParams{UserID: 123})
		}},
		{name: "get user profile audios", method: "getUserProfileAudios", call: func(ctx context.Context, bot *Bot) (any, error) {
			return bot.GetUserProfileAudios(ctx, GetUserProfileAudiosParams{UserID: 123})
		}},
		{name: "get forum topic icon stickers", method: "getForumTopicIconStickers", call: func(ctx context.Context, bot *Bot) (any, error) {
			return bot.GetForumTopicIconStickers(ctx)
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
			_, err := tt.call(ctx, bot)
			if err == nil {
				t.Fatal("expected error")
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
			_, err := tt.call(context.Background(), bot)
			if err == nil {
				t.Fatal("expected error")
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
			_, err := tt.call(context.Background(), bot)
			if err == nil {
				t.Fatal("expected error")
			}
			assertNoToken(t, err, token)
		})
	}
}
