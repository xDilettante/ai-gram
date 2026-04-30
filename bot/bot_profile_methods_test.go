package bot

import (
	"context"
	"encoding/json"
	stderrors "errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	apierrors "github.com/xDilettante/ai-gram/errors"
	"github.com/xDilettante/ai-gram/telegram"
)

func TestSetMyNameSendsPayloadAndDecodesResult(t *testing.T) {
	testBotProfileBoolSuccess(t, "setMyName", func(bot *Bot) (bool, error) {
		return bot.SetMyName(context.Background(), SetMyNameParams{Name: "ai-gram", LanguageCode: "en"})
	}, func(t *testing.T, payload map[string]any) {
		if payload["name"] != "ai-gram" || payload["language_code"] != "en" {
			t.Fatalf("unexpected payload: %#v", payload)
		}
	})
}

func TestSetMyNameAllowsEmptyName(t *testing.T) {
	testBotProfileBoolSuccess(t, "setMyName", func(bot *Bot) (bool, error) {
		return bot.SetMyName(context.Background(), SetMyNameParams{LanguageCode: "en"})
	}, func(t *testing.T, payload map[string]any) {
		if _, ok := payload["name"]; ok {
			t.Fatalf("empty name should be omitted: %#v", payload)
		}
		if payload["language_code"] != "en" {
			t.Fatalf("unexpected payload: %#v", payload)
		}
	})
}

func TestGetMyNameSendsPayloadAndDecodesResult(t *testing.T) {
	testBotProfileObjectSuccess(t, "getMyName", `{"name":"ai-gram"}`, func(bot *Bot) (any, error) {
		return bot.GetMyName(context.Background(), GetMyNameParams{LanguageCode: "en"})
	}, func(t *testing.T, payload map[string]any) {
		if payload["language_code"] != "en" {
			t.Fatalf("unexpected payload: %#v", payload)
		}
	}, func(t *testing.T, result any) {
		name := result.(*telegram.BotName)
		if name.Name != "ai-gram" {
			t.Fatalf("unexpected result: %+v", name)
		}
	})
}

func TestSetMyDescriptionSendsPayloadAndDecodesResult(t *testing.T) {
	testBotProfileBoolSuccess(t, "setMyDescription", func(bot *Bot) (bool, error) {
		return bot.SetMyDescription(context.Background(), SetMyDescriptionParams{Description: "Bot description", LanguageCode: "en"})
	}, func(t *testing.T, payload map[string]any) {
		if payload["description"] != "Bot description" || payload["language_code"] != "en" {
			t.Fatalf("unexpected payload: %#v", payload)
		}
	})
}

func TestSetMyDescriptionAllowsEmptyDescription(t *testing.T) {
	testBotProfileBoolSuccess(t, "setMyDescription", func(bot *Bot) (bool, error) {
		return bot.SetMyDescription(context.Background(), SetMyDescriptionParams{LanguageCode: "en"})
	}, func(t *testing.T, payload map[string]any) {
		if _, ok := payload["description"]; ok {
			t.Fatalf("empty description should be omitted: %#v", payload)
		}
		if payload["language_code"] != "en" {
			t.Fatalf("unexpected payload: %#v", payload)
		}
	})
}

func TestGetMyDescriptionSendsPayloadAndDecodesResult(t *testing.T) {
	testBotProfileObjectSuccess(t, "getMyDescription", `{"description":"Bot description"}`, func(bot *Bot) (any, error) {
		return bot.GetMyDescription(context.Background(), GetMyDescriptionParams{LanguageCode: "en"})
	}, func(t *testing.T, payload map[string]any) {
		if payload["language_code"] != "en" {
			t.Fatalf("unexpected payload: %#v", payload)
		}
	}, func(t *testing.T, result any) {
		description := result.(*telegram.BotDescription)
		if description.Description != "Bot description" {
			t.Fatalf("unexpected result: %+v", description)
		}
	})
}

func TestSetMyShortDescriptionSendsPayloadAndDecodesResult(t *testing.T) {
	testBotProfileBoolSuccess(t, "setMyShortDescription", func(bot *Bot) (bool, error) {
		return bot.SetMyShortDescription(context.Background(), SetMyShortDescriptionParams{ShortDescription: "Short bot", LanguageCode: "en"})
	}, func(t *testing.T, payload map[string]any) {
		if payload["short_description"] != "Short bot" || payload["language_code"] != "en" {
			t.Fatalf("unexpected payload: %#v", payload)
		}
	})
}

func TestSetMyShortDescriptionAllowsEmptyShortDescription(t *testing.T) {
	testBotProfileBoolSuccess(t, "setMyShortDescription", func(bot *Bot) (bool, error) {
		return bot.SetMyShortDescription(context.Background(), SetMyShortDescriptionParams{LanguageCode: "en"})
	}, func(t *testing.T, payload map[string]any) {
		if _, ok := payload["short_description"]; ok {
			t.Fatalf("empty short_description should be omitted: %#v", payload)
		}
		if payload["language_code"] != "en" {
			t.Fatalf("unexpected payload: %#v", payload)
		}
	})
}

func TestGetMyShortDescriptionSendsPayloadAndDecodesResult(t *testing.T) {
	testBotProfileObjectSuccess(t, "getMyShortDescription", `{"short_description":"Short bot"}`, func(bot *Bot) (any, error) {
		return bot.GetMyShortDescription(context.Background(), GetMyShortDescriptionParams{LanguageCode: "en"})
	}, func(t *testing.T, payload map[string]any) {
		if payload["language_code"] != "en" {
			t.Fatalf("unexpected payload: %#v", payload)
		}
	}, func(t *testing.T, result any) {
		shortDescription := result.(*telegram.BotShortDescription)
		if shortDescription.ShortDescription != "Short bot" {
			t.Fatalf("unexpected result: %+v", shortDescription)
		}
	})
}

func TestGetMyDefaultAdministratorRightsSendsPayloadAndDecodesResult(t *testing.T) {
	testBotProfileObjectSuccess(t, "getMyDefaultAdministratorRights", `{"can_manage_chat":true,"can_delete_messages":true,"can_post_messages":true}`, func(bot *Bot) (any, error) {
		return bot.GetMyDefaultAdministratorRights(context.Background(), GetMyDefaultAdministratorRightsParams{ForChannels: true})
	}, func(t *testing.T, payload map[string]any) {
		if payload["for_channels"] != true {
			t.Fatalf("unexpected payload: %#v", payload)
		}
	}, func(t *testing.T, result any) {
		rights := result.(*telegram.ChatAdministratorRights)
		if !rights.CanManageChat || !rights.CanDeleteMessages || !rights.CanPostMessages {
			t.Fatalf("unexpected result: %+v", rights)
		}
	})
}

func TestSetMyProfilePhotoStaticSendsMultipartAndDecodesResult(t *testing.T) {
	const token = "123:secret"
	photoReader := &chunkReader{data: []byte("profile-photo"), chunk: 3}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/bot"+token+"/setMyProfilePhoto" {
			t.Fatalf("unexpected path: %q", r.URL.Path)
		}
		if got := r.Header.Get("Content-Type"); !strings.HasPrefix(got, "multipart/form-data") {
			t.Fatalf("unexpected content type: %q", got)
		}
		if err := r.ParseMultipartForm(1024); err != nil {
			t.Fatalf("parse multipart: %v", err)
		}
		var photo map[string]any
		if err := json.Unmarshal([]byte(r.MultipartForm.Value["photo"][0]), &photo); err != nil {
			t.Fatalf("decode photo field: %v", err)
		}
		if photo["type"] != "static" || photo["photo"] != "attach://photo" {
			t.Fatalf("unexpected photo payload: %#v", photo)
		}
		content, header := readMultipartFile(t, r, "photo")
		if header.Filename != "profile.jpg" {
			t.Fatalf("unexpected filename: %q", header.Filename)
		}
		if got := header.Header.Get("Content-Type"); got != "image/jpeg" {
			t.Fatalf("unexpected part content type: %q", got)
		}
		if string(content) != "profile-photo" {
			t.Fatalf("unexpected file content: %q", content)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true,"result":true}`))
	}))
	defer server.Close()

	bot := newTestBot(t, token, server.URL, server.Client())
	ok, err := bot.SetMyProfilePhoto(context.Background(), SetMyProfilePhotoParams{
		Photo: ProfilePhotoStatic(FileUpload(UploadFile{Name: "profile.jpg", Reader: photoReader, ContentType: "image/jpeg"})),
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Fatal("expected true result")
	}
	if photoReader.reads < 2 {
		t.Fatalf("expected chunked streaming reads, got %d", photoReader.reads)
	}
}

func TestSetMyProfilePhotoAnimatedSendsMultipartAndDecodesResult(t *testing.T) {
	const token = "123:secret"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/bot"+token+"/setMyProfilePhoto" {
			t.Fatalf("unexpected path: %q", r.URL.Path)
		}
		if err := r.ParseMultipartForm(1024); err != nil {
			t.Fatalf("parse multipart: %v", err)
		}
		var photo map[string]any
		if err := json.Unmarshal([]byte(r.MultipartForm.Value["photo"][0]), &photo); err != nil {
			t.Fatalf("decode photo field: %v", err)
		}
		if photo["type"] != "animated" || photo["animation"] != "attach://animation" || photo["main_frame_timestamp"] != 1.5 {
			t.Fatalf("unexpected photo payload: %#v", photo)
		}
		content, header := readMultipartFile(t, r, "animation")
		if header.Filename != "profile.mp4" {
			t.Fatalf("unexpected filename: %q", header.Filename)
		}
		if string(content) != "profile-animation" {
			t.Fatalf("unexpected file content: %q", content)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true,"result":true}`))
	}))
	defer server.Close()

	bot := newTestBot(t, token, server.URL, server.Client())
	ok, err := bot.SetMyProfilePhoto(context.Background(), SetMyProfilePhotoParams{
		Photo: InputProfilePhotoAnimated{
			Type:               "animated",
			Animation:          FileUpload(UploadFile{Name: "profile.mp4", Reader: strings.NewReader("profile-animation"), ContentType: "video/mp4"}),
			MainFrameTimestamp: 1.5,
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Fatal("expected true result")
	}
}

func TestRemoveMyProfilePhotoSendsPayloadAndDecodesResult(t *testing.T) {
	testBotProfileBoolSuccess(t, "removeMyProfilePhoto", func(bot *Bot) (bool, error) {
		return bot.RemoveMyProfilePhoto(context.Background(), RemoveMyProfilePhotoParams{})
	}, func(t *testing.T, payload map[string]any) {
		if len(payload) != 0 {
			t.Fatalf("unexpected payload: %#v", payload)
		}
	})
}

func TestSetMyProfilePhotoValidation(t *testing.T) {
	const token = "123:secret"
	bot := newTestBot(t, token, "https://example.test", nil)
	tests := []struct {
		name   string
		params SetMyProfilePhotoParams
	}{
		{name: "nil photo"},
		{name: "static file id", params: SetMyProfilePhotoParams{Photo: ProfilePhotoStatic(FileID("photo-file-id"))}},
		{name: "static file url", params: SetMyProfilePhotoParams{Photo: ProfilePhotoStatic(FileURL("https://example.com/photo.jpg"))}},
		{name: "static missing upload", params: SetMyProfilePhotoParams{Photo: ProfilePhotoStatic(FileUpload(UploadFile{}))}},
		{name: "static wrong type", params: SetMyProfilePhotoParams{Photo: InputProfilePhotoStatic{Type: "animated", Photo: FileUpload(UploadFile{Name: "profile.jpg", Reader: strings.NewReader("photo")})}}},
		{name: "animated file id", params: SetMyProfilePhotoParams{Photo: ProfilePhotoAnimated(FileID("animation-file-id"))}},
		{name: "animated file url", params: SetMyProfilePhotoParams{Photo: ProfilePhotoAnimated(FileURL("https://example.com/profile.mp4"))}},
		{name: "animated missing upload", params: SetMyProfilePhotoParams{Photo: ProfilePhotoAnimated(FileUpload(UploadFile{}))}},
		{name: "animated negative timestamp", params: SetMyProfilePhotoParams{Photo: InputProfilePhotoAnimated{Animation: FileUpload(UploadFile{Name: "profile.mp4", Reader: strings.NewReader("animation")}), MainFrameTimestamp: -1}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ok, err := bot.SetMyProfilePhoto(context.Background(), tt.params)
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

func TestBotProfileBoolMethodsAPIAndTransportErrors(t *testing.T) {
	tests := []struct {
		name            string
		method          string
		call            func(*Bot) (bool, error)
		callWithContext func(*Bot, context.Context) (bool, error)
	}{
		{name: "set name", method: "setMyName", call: func(bot *Bot) (bool, error) {
			return bot.SetMyName(context.Background(), SetMyNameParams{Name: "ai-gram"})
		}, callWithContext: func(bot *Bot, ctx context.Context) (bool, error) {
			return bot.SetMyName(ctx, SetMyNameParams{Name: "ai-gram"})
		}},
		{name: "set description", method: "setMyDescription", call: func(bot *Bot) (bool, error) {
			return bot.SetMyDescription(context.Background(), SetMyDescriptionParams{Description: "description"})
		}, callWithContext: func(bot *Bot, ctx context.Context) (bool, error) {
			return bot.SetMyDescription(ctx, SetMyDescriptionParams{Description: "description"})
		}},
		{name: "set short description", method: "setMyShortDescription", call: func(bot *Bot) (bool, error) {
			return bot.SetMyShortDescription(context.Background(), SetMyShortDescriptionParams{ShortDescription: "short"})
		}, callWithContext: func(bot *Bot, ctx context.Context) (bool, error) {
			return bot.SetMyShortDescription(ctx, SetMyShortDescriptionParams{ShortDescription: "short"})
		}},
		{name: "set profile photo", method: "setMyProfilePhoto", call: func(bot *Bot) (bool, error) {
			return bot.SetMyProfilePhoto(context.Background(), SetMyProfilePhotoParams{Photo: ProfilePhotoStatic(FileUpload(UploadFile{Name: "profile.jpg", Reader: strings.NewReader("photo")}))})
		}, callWithContext: func(bot *Bot, ctx context.Context) (bool, error) {
			return bot.SetMyProfilePhoto(ctx, SetMyProfilePhotoParams{Photo: ProfilePhotoStatic(FileUpload(UploadFile{Name: "profile.jpg", Reader: strings.NewReader("photo")}))})
		}},
		{name: "remove profile photo", method: "removeMyProfilePhoto", call: func(bot *Bot) (bool, error) {
			return bot.RemoveMyProfilePhoto(context.Background(), RemoveMyProfilePhotoParams{})
		}, callWithContext: func(bot *Bot, ctx context.Context) (bool, error) {
			return bot.RemoveMyProfilePhoto(ctx, RemoveMyProfilePhotoParams{})
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testBoolMethodErrorCases(t, tt.method, tt.call, tt.callWithContext)
		})
	}
}

func TestBotProfileObjectMethodsAPIAndTransportErrors(t *testing.T) {
	tests := []struct {
		name            string
		method          string
		call            func(*Bot) (any, error)
		callWithContext func(*Bot, context.Context) (any, error)
	}{
		{name: "get name", method: "getMyName", call: func(bot *Bot) (any, error) {
			return bot.GetMyName(context.Background(), GetMyNameParams{})
		}, callWithContext: func(bot *Bot, ctx context.Context) (any, error) {
			return bot.GetMyName(ctx, GetMyNameParams{})
		}},
		{name: "get description", method: "getMyDescription", call: func(bot *Bot) (any, error) {
			return bot.GetMyDescription(context.Background(), GetMyDescriptionParams{})
		}, callWithContext: func(bot *Bot, ctx context.Context) (any, error) {
			return bot.GetMyDescription(ctx, GetMyDescriptionParams{})
		}},
		{name: "get short description", method: "getMyShortDescription", call: func(bot *Bot) (any, error) {
			return bot.GetMyShortDescription(context.Background(), GetMyShortDescriptionParams{})
		}, callWithContext: func(bot *Bot, ctx context.Context) (any, error) {
			return bot.GetMyShortDescription(ctx, GetMyShortDescriptionParams{})
		}},
		{name: "get default administrator rights", method: "getMyDefaultAdministratorRights", call: func(bot *Bot) (any, error) {
			return bot.GetMyDefaultAdministratorRights(context.Background(), GetMyDefaultAdministratorRightsParams{})
		}, callWithContext: func(bot *Bot, ctx context.Context) (any, error) {
			return bot.GetMyDefaultAdministratorRights(ctx, GetMyDefaultAdministratorRightsParams{})
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testBotProfileObjectMethodErrorCases(t, tt.method, tt.call, tt.callWithContext)
		})
	}
}

func testBotProfileBoolSuccess(t *testing.T, method string, call func(*Bot) (bool, error), checkPayload func(*testing.T, map[string]any)) {
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
		checkPayload(t, payload)
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

func testBotProfileObjectSuccess(t *testing.T, method string, resultJSON string, call func(*Bot) (any, error), checkPayload func(*testing.T, map[string]any), checkResult func(*testing.T, any)) {
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
		checkPayload(t, payload)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true,"result":` + resultJSON + `}`))
	}))
	defer server.Close()

	bot := newTestBot(t, token, server.URL, server.Client())
	result, err := call(bot)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	checkResult(t, result)
}

func testBotProfileObjectMethodErrorCases(t *testing.T, method string, call func(*Bot) (any, error), callWithContext func(*Bot, context.Context) (any, error)) {
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
		result, err := call(bot)
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

func isNilResult(result any) bool {
	value := reflect.ValueOf(result)
	switch value.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Pointer, reflect.Slice:
		return value.IsNil()
	default:
		return false
	}
}
