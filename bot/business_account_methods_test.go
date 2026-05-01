package bot

import (
	"context"
	"encoding/json"
	stderrors "errors"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	apierrors "github.com/xDilettante/ai-gram/errors"
	"github.com/xDilettante/ai-gram/telegram"
)

func TestBusinessAccountJSONMethodsSendPayloadAndDecodeResult(t *testing.T) {
	const token = "123:secret"
	tests := []struct {
		name   string
		method string
		call   func(*Bot) (bool, error)
		check  func(*testing.T, map[string]any)
	}{
		{name: "read message", method: "readBusinessMessage", call: func(bot *Bot) (bool, error) {
			return bot.ReadBusinessMessage(context.Background(), ReadBusinessMessageParams{BusinessConnectionID: "bc-1", ChatID: ChatIDInt(123), MessageID: 9})
		}, check: func(t *testing.T, payload map[string]any) {
			if payload["business_connection_id"] != "bc-1" || payload["chat_id"] != float64(123) || payload["message_id"] != float64(9) {
				t.Fatalf("unexpected payload: %#v", payload)
			}
		}},
		{name: "set name", method: "setBusinessAccountName", call: func(bot *Bot) (bool, error) {
			return bot.SetBusinessAccountName(context.Background(), SetBusinessAccountNameParams{BusinessConnectionID: "bc-1", FirstName: "Ada", LastName: "Lovelace"})
		}, check: func(t *testing.T, payload map[string]any) {
			if payload["business_connection_id"] != "bc-1" || payload["first_name"] != "Ada" || payload["last_name"] != "Lovelace" {
				t.Fatalf("unexpected payload: %#v", payload)
			}
		}},
		{name: "set username empty allowed", method: "setBusinessAccountUsername", call: func(bot *Bot) (bool, error) {
			return bot.SetBusinessAccountUsername(context.Background(), SetBusinessAccountUsernameParams{BusinessConnectionID: "bc-1"})
		}, check: func(t *testing.T, payload map[string]any) {
			if payload["business_connection_id"] != "bc-1" {
				t.Fatalf("unexpected payload: %#v", payload)
			}
			if _, ok := payload["username"]; ok {
				t.Fatalf("empty username should be omitted: %#v", payload)
			}
		}},
		{name: "set bio", method: "setBusinessAccountBio", call: func(bot *Bot) (bool, error) {
			return bot.SetBusinessAccountBio(context.Background(), SetBusinessAccountBioParams{BusinessConnectionID: "bc-1", Bio: "Business bio"})
		}, check: func(t *testing.T, payload map[string]any) {
			if payload["business_connection_id"] != "bc-1" || payload["bio"] != "Business bio" {
				t.Fatalf("unexpected payload: %#v", payload)
			}
		}},
		{name: "remove profile photo", method: "removeBusinessAccountProfilePhoto", call: func(bot *Bot) (bool, error) {
			return bot.RemoveBusinessAccountProfilePhoto(context.Background(), RemoveBusinessAccountProfilePhotoParams{BusinessConnectionID: "bc-1", IsPublic: true})
		}, check: func(t *testing.T, payload map[string]any) {
			if payload["business_connection_id"] != "bc-1" || payload["is_public"] != true {
				t.Fatalf("unexpected payload: %#v", payload)
			}
		}},
		{name: "set gift settings", method: "setBusinessAccountGiftSettings", call: func(bot *Bot) (bool, error) {
			return bot.SetBusinessAccountGiftSettings(context.Background(), SetBusinessAccountGiftSettingsParams{BusinessConnectionID: "bc-1", ShowGiftButton: true, AcceptedGiftTypes: telegram.AcceptedGiftTypes{UnlimitedGifts: true, UniqueGifts: true}})
		}, check: func(t *testing.T, payload map[string]any) {
			gifts, ok := payload["accepted_gift_types"].(map[string]any)
			if payload["business_connection_id"] != "bc-1" || payload["show_gift_button"] != true || !ok || gifts["unlimited_gifts"] != true || gifts["unique_gifts"] != true {
				t.Fatalf("unexpected payload: %#v", payload)
			}
		}},
		{name: "delete story", method: "deleteStory", call: func(bot *Bot) (bool, error) {
			return bot.DeleteStory(context.Background(), DeleteStoryParams{BusinessConnectionID: "bc-1", StoryID: 77})
		}, check: func(t *testing.T, payload map[string]any) {
			if payload["business_connection_id"] != "bc-1" || payload["story_id"] != float64(77) {
				t.Fatalf("unexpected payload: %#v", payload)
			}
		}},
		{name: "approve suggested post", method: "approveSuggestedPost", call: func(bot *Bot) (bool, error) {
			return bot.ApproveSuggestedPost(context.Background(), ApproveSuggestedPostParams{ChatID: ChatIDInt(123), MessageID: 55, SendDate: 999})
		}, check: func(t *testing.T, payload map[string]any) {
			if payload["chat_id"] != float64(123) || payload["message_id"] != float64(55) || payload["send_date"] != float64(999) {
				t.Fatalf("unexpected payload: %#v", payload)
			}
		}},
		{name: "decline suggested post", method: "declineSuggestedPost", call: func(bot *Bot) (bool, error) {
			return bot.DeclineSuggestedPost(context.Background(), DeclineSuggestedPostParams{ChatID: ChatIDInt(123), MessageID: 55, Comment: "Not now"})
		}, check: func(t *testing.T, payload map[string]any) {
			if payload["chat_id"] != float64(123) || payload["message_id"] != float64(55) || payload["comment"] != "Not now" {
				t.Fatalf("unexpected payload: %#v", payload)
			}
		}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodPost {
					t.Fatalf("unexpected method: %s", r.Method)
				}
				if r.URL.Path != "/bot"+token+"/"+tt.method {
					t.Fatalf("unexpected path: %q", r.URL.Path)
				}
				var payload map[string]any
				if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
					t.Fatalf("decode payload: %v", err)
				}
				tt.check(t, payload)
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(`{"ok":true,"result":true}`))
			}))
			defer server.Close()

			bot := newTestBot(t, token, server.URL, server.Client())
			ok, err := tt.call(bot)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !ok {
				t.Fatal("expected true result")
			}
		})
	}
}

func TestSetBusinessAccountProfilePhotoSendsMultipartAndDecodesResult(t *testing.T) {
	const token = "123:secret"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/bot"+token+"/setBusinessAccountProfilePhoto" {
			t.Fatalf("unexpected path: %q", r.URL.Path)
		}
		if got := r.Header.Get("Content-Type"); !strings.HasPrefix(got, "multipart/form-data") {
			t.Fatalf("unexpected content type: %q", got)
		}
		if err := r.ParseMultipartForm(4096); err != nil {
			t.Fatalf("parse multipart: %v", err)
		}
		assertMultipartValue(t, r, "business_connection_id", "bc-1")
		assertMultipartValue(t, r, "is_public", "true")
		var photo map[string]any
		if err := json.Unmarshal([]byte(r.MultipartForm.Value["photo"][0]), &photo); err != nil {
			t.Fatalf("decode photo: %v", err)
		}
		if photo["type"] != "static" || photo["photo"] != "attach://photo" {
			t.Fatalf("unexpected photo payload: %#v", photo)
		}
		content, header := readMultipartFile(t, r, "photo")
		if header.Filename != "profile.jpg" || string(content) != "profile-photo" {
			t.Fatalf("unexpected file: %q %q", header.Filename, content)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true,"result":true}`))
	}))
	defer server.Close()

	bot := newTestBot(t, token, server.URL, server.Client())
	ok, err := bot.SetBusinessAccountProfilePhoto(context.Background(), SetBusinessAccountProfilePhotoParams{
		BusinessConnectionID: "bc-1",
		Photo:                ProfilePhotoStatic(FileUpload(UploadFile{Name: "profile.jpg", Reader: strings.NewReader("profile-photo"), ContentType: "image/jpeg"})),
		IsPublic:             true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Fatal("expected true result")
	}
}

func TestPostAndEditStorySendMultipartAndDecodeStory(t *testing.T) {
	const token = "123:secret"
	tests := []struct {
		name          string
		method        string
		call          func(*Bot) (*telegram.Story, error)
		expectStoryID int64
		check         func(*testing.T, *http.Request)
	}{
		{name: "post", method: "postStory", expectStoryID: 7, call: func(bot *Bot) (*telegram.Story, error) {
			return bot.PostStory(context.Background(), PostStoryParams{
				BusinessConnectionID: "bc-1",
				Content:              StoryPhoto(FileUpload(UploadFile{Name: "story.jpg", Reader: strings.NewReader("story-photo"), ContentType: "image/jpeg"})),
				ActivePeriod:         86400,
				Caption:              "caption",
				ParseMode:            "HTML",
				Areas: []telegram.StoryArea{{
					Position: telegram.StoryAreaPosition{XPercentage: 50, YPercentage: 50, WidthPercentage: 20, HeightPercentage: 20},
					Type:     telegram.NewStoryAreaTypeLink("https://example.test/story"),
				}},
				PostToChatPage: true,
				ProtectContent: true,
			})
		}, check: func(t *testing.T, r *http.Request) {
			assertMultipartValue(t, r, "active_period", "86400")
			assertMultipartValue(t, r, "post_to_chat_page", "true")
			assertMultipartValue(t, r, "protect_content", "true")
			content, header := readMultipartFile(t, r, "story_photo")
			if header.Filename != "story.jpg" || string(content) != "story-photo" {
				t.Fatalf("unexpected story photo: %q %q", header.Filename, content)
			}
		}},
		{name: "edit", method: "editStory", expectStoryID: 8, call: func(bot *Bot) (*telegram.Story, error) {
			video := StoryVideo(FileUpload(UploadFile{Name: "story.mp4", Reader: strings.NewReader("story-video"), ContentType: "video/mp4"}))
			video.Duration = 4.5
			video.CoverFrameTimestamp = 1.25
			video.IsAnimation = true
			return bot.EditStory(context.Background(), EditStoryParams{BusinessConnectionID: "bc-1", StoryID: 8, Content: video})
		}, check: func(t *testing.T, r *http.Request) {
			assertMultipartValue(t, r, "story_id", "8")
			content, header := readMultipartFile(t, r, "story_video")
			if header.Filename != "story.mp4" || string(content) != "story-video" {
				t.Fatalf("unexpected story video: %q %q", header.Filename, content)
			}
			var storyContent map[string]any
			if err := json.Unmarshal([]byte(r.MultipartForm.Value["content"][0]), &storyContent); err != nil {
				t.Fatalf("decode content: %v", err)
			}
			if storyContent["type"] != "video" || storyContent["video"] != "attach://story_video" || storyContent["duration"] != 4.5 || storyContent["cover_frame_timestamp"] != 1.25 || storyContent["is_animation"] != true {
				t.Fatalf("unexpected story content: %#v", storyContent)
			}
		}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path != "/bot"+token+"/"+tt.method {
					t.Fatalf("unexpected path: %q", r.URL.Path)
				}
				if err := r.ParseMultipartForm(8192); err != nil {
					t.Fatalf("parse multipart: %v", err)
				}
				assertMultipartValue(t, r, "business_connection_id", "bc-1")
				var content map[string]any
				if err := json.Unmarshal([]byte(r.MultipartForm.Value["content"][0]), &content); err != nil {
					t.Fatalf("decode content: %v", err)
				}
				if content["type"] == "" {
					t.Fatalf("missing story content type: %#v", content)
				}
				tt.check(t, r)
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(`{"ok":true,"result":{"chat":{"id":123,"type":"private"},"id":` + strconvInt64(tt.expectStoryID) + `}}`))
			}))
			defer server.Close()

			bot := newTestBot(t, token, server.URL, server.Client())
			story, err := tt.call(bot)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if story == nil || story.ID != tt.expectStoryID || story.Chat.ID != 123 {
				t.Fatalf("unexpected story: %+v", story)
			}
		})
	}
}

func TestRepostStorySendsJSONAndDecodesStory(t *testing.T) {
	const token = "123:secret"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/bot"+token+"/repostStory" {
			t.Fatalf("unexpected path: %q", r.URL.Path)
		}
		if got := r.Header.Get("Content-Type"); got != "application/json" {
			t.Fatalf("unexpected content type: %q", got)
		}
		var payload map[string]any
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode payload: %v", err)
		}
		if payload["business_connection_id"] != "bc-1" || payload["from_chat_id"] != float64(123) || payload["from_story_id"] != float64(7) || payload["active_period"] != float64(86400) || payload["post_to_chat_page"] != true || payload["protect_content"] != true {
			t.Fatalf("unexpected payload: %#v", payload)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true,"result":{"chat":{"id":123,"type":"private"},"id":9}}`))
	}))
	defer server.Close()

	bot := newTestBot(t, token, server.URL, server.Client())
	story, err := bot.RepostStory(context.Background(), RepostStoryParams{
		BusinessConnectionID: "bc-1",
		FromChatID:           123,
		FromStoryID:          7,
		ActivePeriod:         86400,
		PostToChatPage:       true,
		ProtectContent:       true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if story == nil || story.ID != 9 || story.Chat.ID != 123 {
		t.Fatalf("unexpected story: %+v", story)
	}
}

func TestBusinessAccountMethodsValidation(t *testing.T) {
	const token = "123:secret"
	bot := newTestBot(t, token, "https://example.test", nil)
	tests := []struct {
		name string
		call func(*Bot) error
	}{
		{name: "read missing connection", call: func(bot *Bot) error {
			_, err := bot.ReadBusinessMessage(context.Background(), ReadBusinessMessageParams{ChatID: ChatIDInt(1), MessageID: 1})
			return err
		}},
		{name: "read missing chat", call: func(bot *Bot) error {
			_, err := bot.ReadBusinessMessage(context.Background(), ReadBusinessMessageParams{BusinessConnectionID: "bc", MessageID: 1})
			return err
		}},
		{name: "read invalid message", call: func(bot *Bot) error {
			_, err := bot.ReadBusinessMessage(context.Background(), ReadBusinessMessageParams{BusinessConnectionID: "bc", ChatID: ChatIDInt(1)})
			return err
		}},
		{name: "name missing first", call: func(bot *Bot) error {
			_, err := bot.SetBusinessAccountName(context.Background(), SetBusinessAccountNameParams{BusinessConnectionID: "bc"})
			return err
		}},
		{name: "username missing connection", call: func(bot *Bot) error {
			_, err := bot.SetBusinessAccountUsername(context.Background(), SetBusinessAccountUsernameParams{})
			return err
		}},
		{name: "bio missing connection", call: func(bot *Bot) error {
			_, err := bot.SetBusinessAccountBio(context.Background(), SetBusinessAccountBioParams{})
			return err
		}},
		{name: "profile missing photo", call: func(bot *Bot) error {
			_, err := bot.SetBusinessAccountProfilePhoto(context.Background(), SetBusinessAccountProfilePhotoParams{BusinessConnectionID: "bc"})
			return err
		}},
		{name: "profile file id rejected", call: func(bot *Bot) error {
			_, err := bot.SetBusinessAccountProfilePhoto(context.Background(), SetBusinessAccountProfilePhotoParams{BusinessConnectionID: "bc", Photo: ProfilePhotoStatic(FileID("file"))})
			return err
		}},
		{name: "remove photo missing connection", call: func(bot *Bot) error {
			_, err := bot.RemoveBusinessAccountProfilePhoto(context.Background(), RemoveBusinessAccountProfilePhotoParams{})
			return err
		}},
		{name: "gift settings missing accepted types", call: func(bot *Bot) error {
			_, err := bot.SetBusinessAccountGiftSettings(context.Background(), SetBusinessAccountGiftSettingsParams{BusinessConnectionID: "bc"})
			return err
		}},
		{name: "post invalid active period", call: func(bot *Bot) error {
			_, err := bot.PostStory(context.Background(), PostStoryParams{BusinessConnectionID: "bc", Content: StoryPhoto(FileUpload(UploadFile{Name: "s.jpg", Reader: strings.NewReader("x")})), ActivePeriod: 1})
			return err
		}},
		{name: "post missing content", call: func(bot *Bot) error {
			_, err := bot.PostStory(context.Background(), PostStoryParams{BusinessConnectionID: "bc", ActivePeriod: 86400})
			return err
		}},
		{name: "post caption conflict", call: func(bot *Bot) error {
			_, err := bot.PostStory(context.Background(), PostStoryParams{BusinessConnectionID: "bc", Content: StoryPhoto(FileUpload(UploadFile{Name: "s.jpg", Reader: strings.NewReader("x")})), ActivePeriod: 86400, ParseMode: "HTML", CaptionEntities: []telegram.MessageEntity{{Type: telegram.EntityBold, Offset: 0, Length: 1}}})
			return err
		}},
		{name: "post invalid area", call: func(bot *Bot) error {
			_, err := bot.PostStory(context.Background(), PostStoryParams{BusinessConnectionID: "bc", Content: StoryPhoto(FileUpload(UploadFile{Name: "s.jpg", Reader: strings.NewReader("x")})), ActivePeriod: 86400, Areas: []telegram.StoryArea{{Type: telegram.NewStoryAreaTypeLink("")}}})
			return err
		}},
		{name: "edit invalid story", call: func(bot *Bot) error {
			_, err := bot.EditStory(context.Background(), EditStoryParams{BusinessConnectionID: "bc", Content: StoryPhoto(FileUpload(UploadFile{Name: "s.jpg", Reader: strings.NewReader("x")}))})
			return err
		}},
		{name: "delete invalid story", call: func(bot *Bot) error {
			_, err := bot.DeleteStory(context.Background(), DeleteStoryParams{BusinessConnectionID: "bc"})
			return err
		}},
		{name: "repost missing connection", call: func(bot *Bot) error {
			_, err := bot.RepostStory(context.Background(), RepostStoryParams{FromChatID: 1, FromStoryID: 1, ActivePeriod: 86400})
			return err
		}},
		{name: "repost missing source chat", call: func(bot *Bot) error {
			_, err := bot.RepostStory(context.Background(), RepostStoryParams{BusinessConnectionID: "bc", FromStoryID: 1, ActivePeriod: 86400})
			return err
		}},
		{name: "repost invalid story", call: func(bot *Bot) error {
			_, err := bot.RepostStory(context.Background(), RepostStoryParams{BusinessConnectionID: "bc", FromChatID: 1, ActivePeriod: 86400})
			return err
		}},
		{name: "repost invalid active period", call: func(bot *Bot) error {
			_, err := bot.RepostStory(context.Background(), RepostStoryParams{BusinessConnectionID: "bc", FromChatID: 1, FromStoryID: 1, ActivePeriod: 1})
			return err
		}},
		{name: "approve invalid chat", call: func(bot *Bot) error {
			_, err := bot.ApproveSuggestedPost(context.Background(), ApproveSuggestedPostParams{MessageID: 1})
			return err
		}},
		{name: "approve negative send date", call: func(bot *Bot) error {
			_, err := bot.ApproveSuggestedPost(context.Background(), ApproveSuggestedPostParams{ChatID: ChatIDInt(1), MessageID: 1, SendDate: -1})
			return err
		}},
		{name: "decline invalid message", call: func(bot *Bot) error {
			_, err := bot.DeclineSuggestedPost(context.Background(), DeclineSuggestedPostParams{ChatID: ChatIDInt(1)})
			return err
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.call(bot); err == nil {
				t.Fatal("expected error")
			} else {
				assertNoToken(t, err, token)
			}
		})
	}
}

func TestBusinessAccountMethodsReturnAPIError(t *testing.T) {
	const token = "123:secret"
	tests := []struct {
		name   string
		method string
		call   func(*Bot) error
	}{
		{name: "read", method: "readBusinessMessage", call: func(bot *Bot) error {
			_, err := bot.ReadBusinessMessage(context.Background(), ReadBusinessMessageParams{BusinessConnectionID: "bc", ChatID: ChatIDInt(1), MessageID: 1})
			return err
		}},
		{name: "profile", method: "setBusinessAccountProfilePhoto", call: func(bot *Bot) error {
			_, err := bot.SetBusinessAccountProfilePhoto(context.Background(), SetBusinessAccountProfilePhotoParams{BusinessConnectionID: "bc", Photo: ProfilePhotoStatic(FileUpload(UploadFile{Name: "p.jpg", Reader: strings.NewReader("p")}))})
			return err
		}},
		{name: "post story", method: "postStory", call: func(bot *Bot) error {
			_, err := bot.PostStory(context.Background(), PostStoryParams{BusinessConnectionID: "bc", Content: StoryPhoto(FileUpload(UploadFile{Name: "s.jpg", Reader: strings.NewReader("s")})), ActivePeriod: 86400})
			return err
		}},
		{name: "repost story", method: "repostStory", call: func(bot *Bot) error {
			_, err := bot.RepostStory(context.Background(), RepostStoryParams{BusinessConnectionID: "bc", FromChatID: 1, FromStoryID: 1, ActivePeriod: 86400})
			return err
		}},
		{name: "approve", method: "approveSuggestedPost", call: func(bot *Bot) error {
			_, err := bot.ApproveSuggestedPost(context.Background(), ApproveSuggestedPostParams{ChatID: ChatIDInt(1), MessageID: 1})
			return err
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
			err := tt.call(bot)
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

func TestBusinessAccountMethodsResponseAndContextErrors(t *testing.T) {
	const token = "123:secret"
	tests := []struct {
		name   string
		method string
		call   func(context.Context, *Bot) error
	}{
		{name: "read", method: "readBusinessMessage", call: func(ctx context.Context, bot *Bot) error {
			_, err := bot.ReadBusinessMessage(ctx, ReadBusinessMessageParams{BusinessConnectionID: "bc", ChatID: ChatIDInt(1), MessageID: 1})
			return err
		}},
		{name: "delete story", method: "deleteStory", call: func(ctx context.Context, bot *Bot) error {
			_, err := bot.DeleteStory(ctx, DeleteStoryParams{BusinessConnectionID: "bc", StoryID: 1})
			return err
		}},
		{name: "post story", method: "postStory", call: func(ctx context.Context, bot *Bot) error {
			_, err := bot.PostStory(ctx, PostStoryParams{BusinessConnectionID: "bc", Content: StoryPhoto(FileUpload(UploadFile{Name: "s.jpg", Reader: strings.NewReader("s")})), ActivePeriod: 86400})
			return err
		}},
		{name: "repost story", method: "repostStory", call: func(ctx context.Context, bot *Bot) error {
			_, err := bot.RepostStory(ctx, RepostStoryParams{BusinessConnectionID: "bc", FromChatID: 1, FromStoryID: 1, ActivePeriod: 86400})
			return err
		}},
	}
	for _, tt := range tests {
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
			if err := tt.call(context.Background(), bot); err == nil {
				t.Fatal("expected error")
			} else {
				assertNoToken(t, err, token)
			}
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
			if err := tt.call(context.Background(), bot); err == nil {
				t.Fatal("expected error")
			} else {
				assertNoToken(t, err, token)
			}
		})
		t.Run(tt.name+" cancelled context", func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { t.Fatal("request should not reach server") }))
			defer server.Close()
			bot := newTestBot(t, token, server.URL, server.Client())
			ctx, cancel := context.WithCancel(context.Background())
			cancel()
			if err := tt.call(ctx, bot); err == nil {
				t.Fatal("expected error")
			} else {
				assertNoToken(t, err, token)
			}
		})
	}
}

func strconvInt64(v int64) string {
	return strconv.FormatInt(v, 10)
}
