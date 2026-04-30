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
)

func TestSetChatTitleSendsPayloadAndDecodesResult(t *testing.T) {
	const token = "123:secret"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/bot"+token+"/setChatTitle" {
			t.Fatalf("unexpected path: %q", r.URL.Path)
		}
		var payload map[string]any
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode payload: %v", err)
		}
		if payload["chat_id"] != float64(12345) || payload["title"] != "New title" {
			t.Fatalf("unexpected payload: %#v", payload)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true,"result":true}`))
	}))
	defer server.Close()

	bot := newTestBot(t, token, server.URL, server.Client())
	ok, err := bot.SetChatTitle(context.Background(), SetChatTitleParams{ChatID: ChatIDInt(12345), Title: "New title"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Fatal("expected true result")
	}
}

func TestSetChatDescriptionSendsPayloadAndDecodesResult(t *testing.T) {
	const token = "123:secret"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/bot"+token+"/setChatDescription" {
			t.Fatalf("unexpected path: %q", r.URL.Path)
		}
		var payload map[string]any
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode payload: %v", err)
		}
		if payload["chat_id"] != float64(12345) || payload["description"] != "New description" {
			t.Fatalf("unexpected payload: %#v", payload)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true,"result":true}`))
	}))
	defer server.Close()

	bot := newTestBot(t, token, server.URL, server.Client())
	ok, err := bot.SetChatDescription(context.Background(), SetChatDescriptionParams{ChatID: ChatIDInt(12345), Description: "New description"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Fatal("expected true result")
	}
}

func TestSetChatDescriptionAllowsEmptyDescription(t *testing.T) {
	const token = "123:secret"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var payload map[string]any
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode payload: %v", err)
		}
		if payload["chat_id"] != float64(12345) {
			t.Fatalf("unexpected payload: %#v", payload)
		}
		if _, ok := payload["description"]; ok {
			t.Fatalf("empty description should be omitted: %#v", payload)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true,"result":true}`))
	}))
	defer server.Close()

	bot := newTestBot(t, token, server.URL, server.Client())
	ok, err := bot.SetChatDescription(context.Background(), SetChatDescriptionParams{ChatID: ChatIDInt(12345)})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Fatal("expected true result")
	}
}

func TestSetChatPhotoSendsMultipartAndDecodesResult(t *testing.T) {
	const token = "123:secret"
	photoReader := &chunkReader{data: []byte("photo-data"), chunk: 2}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/bot"+token+"/setChatPhoto" {
			t.Fatalf("unexpected path: %q", r.URL.Path)
		}
		if got := r.Header.Get("Content-Type"); !strings.HasPrefix(got, "multipart/form-data") {
			t.Fatalf("unexpected content type: %q", got)
		}
		if err := r.ParseMultipartForm(1024); err != nil {
			t.Fatalf("parse multipart: %v", err)
		}
		assertMultipartValue(t, r, "chat_id", "12345")
		assertMultipartValue(t, r, "photo", "attach://photo")
		content, header := readMultipartFile(t, r, "photo")
		if header.Filename != "chat.jpg" {
			t.Fatalf("unexpected filename: %q", header.Filename)
		}
		if got := header.Header.Get("Content-Type"); got != "image/jpeg" {
			t.Fatalf("unexpected part content type: %q", got)
		}
		if string(content) != "photo-data" {
			t.Fatalf("unexpected file content: %q", content)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true,"result":true}`))
	}))
	defer server.Close()

	bot := newTestBot(t, token, server.URL, server.Client())
	ok, err := bot.SetChatPhoto(context.Background(), SetChatPhotoParams{
		ChatID: ChatIDInt(12345),
		Photo:  FileUpload(UploadFile{Name: "chat.jpg", Reader: photoReader, ContentType: "image/jpeg"}),
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

func TestDeleteChatPhotoSendsPayloadAndDecodesResult(t *testing.T) {
	testChatManagementSimpleBoolSuccess(t, "deleteChatPhoto", func(bot *Bot) (bool, error) {
		return bot.DeleteChatPhoto(context.Background(), DeleteChatPhotoParams{ChatID: ChatIDInt(12345)})
	}, func(t *testing.T, payload map[string]any) {
		if payload["chat_id"] != float64(12345) {
			t.Fatalf("unexpected payload: %#v", payload)
		}
	})
}

func TestLeaveChatSendsPayloadAndDecodesResult(t *testing.T) {
	testChatManagementSimpleBoolSuccess(t, "leaveChat", func(bot *Bot) (bool, error) {
		return bot.LeaveChat(context.Background(), LeaveChatParams{ChatID: ChatIDInt(12345)})
	}, func(t *testing.T, payload map[string]any) {
		if payload["chat_id"] != float64(12345) {
			t.Fatalf("unexpected payload: %#v", payload)
		}
	})
}

func TestSetChatStickerSetSendsPayloadAndDecodesResult(t *testing.T) {
	testChatManagementSimpleBoolSuccess(t, "setChatStickerSet", func(bot *Bot) (bool, error) {
		return bot.SetChatStickerSet(context.Background(), SetChatStickerSetParams{ChatID: ChatIDInt(12345), StickerSetName: "test_by_bot"})
	}, func(t *testing.T, payload map[string]any) {
		if payload["chat_id"] != float64(12345) || payload["sticker_set_name"] != "test_by_bot" {
			t.Fatalf("unexpected payload: %#v", payload)
		}
	})
}

func TestDeleteChatStickerSetSendsPayloadAndDecodesResult(t *testing.T) {
	testChatManagementSimpleBoolSuccess(t, "deleteChatStickerSet", func(bot *Bot) (bool, error) {
		return bot.DeleteChatStickerSet(context.Background(), DeleteChatStickerSetParams{ChatID: ChatIDInt(12345)})
	}, func(t *testing.T, payload map[string]any) {
		if payload["chat_id"] != float64(12345) {
			t.Fatalf("unexpected payload: %#v", payload)
		}
	})
}

func TestChatManagementMethodValidation(t *testing.T) {
	const token = "123:secret"
	bot := newTestBot(t, token, "https://example.test", nil)
	tests := []struct {
		name string
		call func() (bool, error)
	}{
		{name: "set title empty chat", call: func() (bool, error) {
			return bot.SetChatTitle(context.Background(), SetChatTitleParams{Title: "title"})
		}},
		{name: "set title empty title", call: func() (bool, error) {
			return bot.SetChatTitle(context.Background(), SetChatTitleParams{ChatID: ChatIDInt(123)})
		}},
		{name: "set description empty chat", call: func() (bool, error) {
			return bot.SetChatDescription(context.Background(), SetChatDescriptionParams{Description: "description"})
		}},
		{name: "set photo empty chat", call: func() (bool, error) {
			return bot.SetChatPhoto(context.Background(), SetChatPhotoParams{Photo: FileUpload(UploadFile{Name: "chat.jpg", Reader: strings.NewReader("photo")})})
		}},
		{name: "set photo file id", call: func() (bool, error) {
			return bot.SetChatPhoto(context.Background(), SetChatPhotoParams{ChatID: ChatIDInt(123), Photo: FileID("photo-file-id")})
		}},
		{name: "set photo file url", call: func() (bool, error) {
			return bot.SetChatPhoto(context.Background(), SetChatPhotoParams{ChatID: ChatIDInt(123), Photo: FileURL("https://example.com/photo.jpg")})
		}},
		{name: "set photo missing upload", call: func() (bool, error) {
			return bot.SetChatPhoto(context.Background(), SetChatPhotoParams{ChatID: ChatIDInt(123), Photo: FileUpload(UploadFile{})})
		}},
		{name: "delete photo empty chat", call: func() (bool, error) { return bot.DeleteChatPhoto(context.Background(), DeleteChatPhotoParams{}) }},
		{name: "leave empty chat", call: func() (bool, error) { return bot.LeaveChat(context.Background(), LeaveChatParams{}) }},
		{name: "set sticker set empty chat", call: func() (bool, error) {
			return bot.SetChatStickerSet(context.Background(), SetChatStickerSetParams{StickerSetName: "set"})
		}},
		{name: "set sticker set empty name", call: func() (bool, error) {
			return bot.SetChatStickerSet(context.Background(), SetChatStickerSetParams{ChatID: ChatIDInt(123)})
		}},
		{name: "delete sticker set empty chat", call: func() (bool, error) {
			return bot.DeleteChatStickerSet(context.Background(), DeleteChatStickerSetParams{})
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ok, err := tt.call()
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

func TestChatManagementMethodsReturnAPIError(t *testing.T) {
	const token = "123:secret"
	tests := chatManagementErrorCases()
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
			ok, err := tt.call(context.Background(), bot)
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

func TestChatManagementMethodsResponseAndContextErrors(t *testing.T) {
	const token = "123:secret"
	tests := chatManagementErrorCases()
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

func testChatManagementSimpleBoolSuccess(t *testing.T, method string, call func(*Bot) (bool, error), checkPayload func(*testing.T, map[string]any)) {
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

type chatManagementErrorCase struct {
	name   string
	method string
	call   func(context.Context, *Bot) (bool, error)
}

func chatManagementErrorCases() []chatManagementErrorCase {
	return []chatManagementErrorCase{
		{name: "set title", method: "setChatTitle", call: func(ctx context.Context, bot *Bot) (bool, error) {
			return bot.SetChatTitle(ctx, SetChatTitleParams{ChatID: ChatIDInt(123), Title: "title"})
		}},
		{name: "set description", method: "setChatDescription", call: func(ctx context.Context, bot *Bot) (bool, error) {
			return bot.SetChatDescription(ctx, SetChatDescriptionParams{ChatID: ChatIDInt(123), Description: "description"})
		}},
		{name: "set photo", method: "setChatPhoto", call: func(ctx context.Context, bot *Bot) (bool, error) {
			return bot.SetChatPhoto(ctx, SetChatPhotoParams{ChatID: ChatIDInt(123), Photo: FileUpload(UploadFile{Name: "chat.jpg", Reader: strings.NewReader("photo")})})
		}},
		{name: "delete photo", method: "deleteChatPhoto", call: func(ctx context.Context, bot *Bot) (bool, error) {
			return bot.DeleteChatPhoto(ctx, DeleteChatPhotoParams{ChatID: ChatIDInt(123)})
		}},
		{name: "leave chat", method: "leaveChat", call: func(ctx context.Context, bot *Bot) (bool, error) {
			return bot.LeaveChat(ctx, LeaveChatParams{ChatID: ChatIDInt(123)})
		}},
		{name: "set sticker set", method: "setChatStickerSet", call: func(ctx context.Context, bot *Bot) (bool, error) {
			return bot.SetChatStickerSet(ctx, SetChatStickerSetParams{ChatID: ChatIDInt(123), StickerSetName: "set_by_bot"})
		}},
		{name: "delete sticker set", method: "deleteChatStickerSet", call: func(ctx context.Context, bot *Bot) (bool, error) {
			return bot.DeleteChatStickerSet(ctx, DeleteChatStickerSetParams{ChatID: ChatIDInt(123)})
		}},
	}
}
