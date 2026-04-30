package bot

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/xDilettante/ai-gram/telegram"
)

func TestInlineMediaResultMarshalAndValidation(t *testing.T) {
	markup := telegram.NewInlineKeyboard([]telegram.InlineKeyboardButton{telegram.InlineButtonCallback("Open", "open")})
	content := InputText("fallback")
	tests := []struct {
		name       string
		result     InlineQueryResult
		wantType   string
		wantFields map[string]string
	}{
		{name: "photo", result: func() InlineQueryResultPhoto {
			r := InlinePhoto("photo-1", "https://example.com/photo.jpg", "https://example.com/thumb.jpg")
			r.Caption = "caption"
			r.ReplyMarkup = &markup
			r.InputMessageContent = content
			return r
		}(), wantType: "photo", wantFields: map[string]string{"photo_url": "https://example.com/photo.jpg", "thumbnail_url": "https://example.com/thumb.jpg"}},
		{name: "gif", result: func() InlineQueryResultGif {
			r := InlineGif("gif-1", "https://example.com/anim.gif", "https://example.com/thumb.jpg")
			r.GifWidth = 320
			r.GifHeight = 240
			r.GifDuration = 3
			r.Caption = "caption"
			r.ReplyMarkup = &markup
			r.InputMessageContent = content
			return r
		}(), wantType: "gif", wantFields: map[string]string{"gif_url": "https://example.com/anim.gif", "thumbnail_url": "https://example.com/thumb.jpg"}},
		{name: "mpeg4 gif", result: func() InlineQueryResultMpeg4Gif {
			r := InlineMpeg4Gif("mpeg4-1", "https://example.com/anim.mp4", "https://example.com/thumb.jpg")
			r.Mpeg4Width = 320
			r.Mpeg4Height = 240
			r.Mpeg4Duration = 3
			return r
		}(), wantType: "mpeg4_gif", wantFields: map[string]string{"mpeg4_url": "https://example.com/anim.mp4", "thumbnail_url": "https://example.com/thumb.jpg"}},
		{name: "video", result: func() InlineQueryResultVideo {
			r := InlineVideo("video-1", "https://example.com/video.mp4", "video/mp4", "https://example.com/thumb.jpg", "Video")
			r.VideoWidth = 320
			r.VideoHeight = 240
			r.VideoDuration = 5
			return r
		}(), wantType: "video", wantFields: map[string]string{"video_url": "https://example.com/video.mp4", "mime_type": "video/mp4", "thumbnail_url": "https://example.com/thumb.jpg", "title": "Video"}},
		{name: "audio", result: func() InlineQueryResultAudio {
			r := InlineAudio("audio-1", "https://example.com/audio.mp3", "Audio")
			r.AudioDuration = 5
			return r
		}(), wantType: "audio", wantFields: map[string]string{"audio_url": "https://example.com/audio.mp3", "title": "Audio"}},
		{name: "voice", result: func() InlineQueryResultVoice {
			r := InlineVoice("voice-1", "https://example.com/voice.ogg", "Voice")
			r.VoiceDuration = 5
			return r
		}(), wantType: "voice", wantFields: map[string]string{"voice_url": "https://example.com/voice.ogg", "title": "Voice"}},
		{name: "document", result: func() InlineQueryResultDocument {
			r := InlineDocument("doc-1", "Document", "https://example.com/file.pdf", "application/pdf")
			r.ThumbnailURL = "https://example.com/thumb.jpg"
			r.ThumbnailWidth = 100
			r.ThumbnailHeight = 50
			return r
		}(), wantType: "document", wantFields: map[string]string{"document_url": "https://example.com/file.pdf", "mime_type": "application/pdf", "title": "Document"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := json.Marshal(tt.result)
			if err != nil {
				t.Fatalf("marshal result: %v", err)
			}
			var got map[string]any
			if err := json.Unmarshal(data, &got); err != nil {
				t.Fatalf("decode result: %v", err)
			}
			if got["type"] != tt.wantType {
				t.Fatalf("unexpected type: %#v", got)
			}
			for field, want := range tt.wantFields {
				if got[field] != want {
					t.Fatalf("unexpected %s: %#v", field, got)
				}
			}
			if err := validateInlineQueryResult(tt.result); err != nil {
				t.Fatalf("valid result rejected: %v", err)
			}
		})
	}
}

func TestInlineMediaResultValidation(t *testing.T) {
	markup := telegram.InlineKeyboardMarkup{}
	tests := []struct {
		name   string
		result InlineQueryResult
	}{
		{name: "photo missing id", result: InlineQueryResultPhoto{PhotoURL: "https://example.com/photo.jpg", ThumbnailURL: "https://example.com/thumb.jpg"}},
		{name: "photo missing url", result: InlinePhoto("photo-1", "", "https://example.com/thumb.jpg")},
		{name: "photo bad thumbnail", result: InlinePhoto("photo-1", "https://example.com/photo.jpg", "ftp://example.com/thumb.jpg")},
		{name: "photo negative width", result: InlineQueryResultPhoto{ID: "photo-1", PhotoURL: "https://example.com/photo.jpg", ThumbnailURL: "https://example.com/thumb.jpg", PhotoWidth: -1}},
		{name: "gif negative duration", result: InlineQueryResultGif{ID: "gif-1", GifURL: "https://example.com/anim.gif", ThumbnailURL: "https://example.com/thumb.jpg", GifDuration: -1}},
		{name: "mpeg4 negative height", result: InlineQueryResultMpeg4Gif{ID: "mpeg4-1", Mpeg4URL: "https://example.com/anim.mp4", ThumbnailURL: "https://example.com/thumb.jpg", Mpeg4Height: -1}},
		{name: "video missing mime", result: InlineQueryResultVideo{ID: "video-1", VideoURL: "https://example.com/video.mp4", ThumbnailURL: "https://example.com/thumb.jpg", Title: "Video"}},
		{name: "video missing title", result: InlineQueryResultVideo{ID: "video-1", VideoURL: "https://example.com/video.mp4", MimeType: "video/mp4", ThumbnailURL: "https://example.com/thumb.jpg"}},
		{name: "audio missing title", result: InlineQueryResultAudio{ID: "audio-1", AudioURL: "https://example.com/audio.mp3"}},
		{name: "audio negative duration", result: InlineQueryResultAudio{ID: "audio-1", AudioURL: "https://example.com/audio.mp3", Title: "Audio", AudioDuration: -1}},
		{name: "voice missing url", result: InlineQueryResultVoice{ID: "voice-1", Title: "Voice"}},
		{name: "voice negative duration", result: InlineQueryResultVoice{ID: "voice-1", VoiceURL: "https://example.com/voice.ogg", Title: "Voice", VoiceDuration: -1}},
		{name: "document missing mime", result: InlineQueryResultDocument{ID: "doc-1", Title: "Document", DocumentURL: "https://example.com/file.pdf"}},
		{name: "document negative thumbnail", result: InlineQueryResultDocument{ID: "doc-1", Title: "Document", DocumentURL: "https://example.com/file.pdf", MimeType: "application/pdf", ThumbnailWidth: -1}},
		{name: "caption conflict", result: InlineQueryResultPhoto{ID: "photo-1", PhotoURL: "https://example.com/photo.jpg", ThumbnailURL: "https://example.com/thumb.jpg", ParseMode: "HTML", CaptionEntities: []telegram.MessageEntity{{Type: telegram.EntityBold, Offset: 0, Length: 1}}}},
		{name: "invalid reply markup", result: InlineQueryResultPhoto{ID: "photo-1", PhotoURL: "https://example.com/photo.jpg", ThumbnailURL: "https://example.com/thumb.jpg", ReplyMarkup: &markup}},
		{name: "invalid input content", result: InlineQueryResultPhoto{ID: "photo-1", PhotoURL: "https://example.com/photo.jpg", ThumbnailURL: "https://example.com/thumb.jpg", InputMessageContent: InputTextMessageContent{}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateInlineQueryResult(tt.result); err == nil {
				t.Fatal("expected validation error")
			}
		})
	}
}

func TestInlineCachedMediaResultMarshalAndValidation(t *testing.T) {
	tests := []struct {
		name       string
		result     InlineQueryResult
		wantType   string
		wantFields map[string]string
	}{
		{name: "cached photo", result: InlineCachedPhoto("photo-1", "photo-file-id"), wantType: "photo", wantFields: map[string]string{"photo_file_id": "photo-file-id"}},
		{name: "cached gif", result: InlineCachedGif("gif-1", "gif-file-id"), wantType: "gif", wantFields: map[string]string{"gif_file_id": "gif-file-id"}},
		{name: "cached mpeg4", result: InlineCachedMpeg4Gif("mpeg4-1", "mpeg4-file-id"), wantType: "mpeg4_gif", wantFields: map[string]string{"mpeg4_file_id": "mpeg4-file-id"}},
		{name: "cached sticker", result: InlineCachedSticker("sticker-1", "sticker-file-id"), wantType: "sticker", wantFields: map[string]string{"sticker_file_id": "sticker-file-id"}},
		{name: "cached document", result: InlineCachedDocument("doc-1", "doc-file-id", "Document"), wantType: "document", wantFields: map[string]string{"document_file_id": "doc-file-id", "title": "Document"}},
		{name: "cached video", result: InlineCachedVideo("video-1", "video-file-id", "Video"), wantType: "video", wantFields: map[string]string{"video_file_id": "video-file-id", "title": "Video"}},
		{name: "cached voice", result: InlineCachedVoice("voice-1", "voice-file-id", "Voice"), wantType: "voice", wantFields: map[string]string{"voice_file_id": "voice-file-id", "title": "Voice"}},
		{name: "cached audio", result: InlineCachedAudio("audio-1", "audio-file-id"), wantType: "audio", wantFields: map[string]string{"audio_file_id": "audio-file-id"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := json.Marshal(tt.result)
			if err != nil {
				t.Fatalf("marshal result: %v", err)
			}
			var got map[string]any
			if err := json.Unmarshal(data, &got); err != nil {
				t.Fatalf("decode result: %v", err)
			}
			if got["type"] != tt.wantType {
				t.Fatalf("unexpected type: %#v", got)
			}
			for field, want := range tt.wantFields {
				if got[field] != want {
					t.Fatalf("unexpected %s: %#v", field, got)
				}
			}
			if err := validateInlineQueryResult(tt.result); err != nil {
				t.Fatalf("valid result rejected: %v", err)
			}
		})
	}
}

func TestInlineCachedMediaResultValidation(t *testing.T) {
	tests := []struct {
		name   string
		result InlineQueryResult
	}{
		{name: "cached photo missing id", result: InlineQueryResultCachedPhoto{PhotoFileID: "photo-file-id"}},
		{name: "cached photo missing file", result: InlineCachedPhoto("photo-1", "")},
		{name: "cached gif missing file", result: InlineCachedGif("gif-1", "")},
		{name: "cached mpeg4 missing file", result: InlineCachedMpeg4Gif("mpeg4-1", "")},
		{name: "cached sticker missing file", result: InlineCachedSticker("sticker-1", "")},
		{name: "cached document missing title", result: InlineQueryResultCachedDocument{ID: "doc-1", DocumentFileID: "doc-file-id"}},
		{name: "cached document missing file", result: InlineCachedDocument("doc-1", "", "Document")},
		{name: "cached video missing title", result: InlineQueryResultCachedVideo{ID: "video-1", VideoFileID: "video-file-id"}},
		{name: "cached voice missing title", result: InlineQueryResultCachedVoice{ID: "voice-1", VoiceFileID: "voice-file-id"}},
		{name: "cached audio missing file", result: InlineCachedAudio("audio-1", "")},
		{name: "caption conflict", result: InlineQueryResultCachedPhoto{ID: "photo-1", PhotoFileID: "photo-file-id", ParseMode: "HTML", CaptionEntities: []telegram.MessageEntity{{Type: telegram.EntityBold, Offset: 0, Length: 1}}}},
		{name: "invalid input content", result: InlineQueryResultCachedPhoto{ID: "photo-1", PhotoFileID: "photo-file-id", InputMessageContent: InputTextMessageContent{}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateInlineQueryResult(tt.result); err == nil {
				t.Fatal("expected validation error")
			}
		})
	}
}

func TestAnswerInlineQueryWithMixedMediaResults(t *testing.T) {
	const token = "123:secret"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/bot"+token+"/answerInlineQuery" {
			t.Fatalf("unexpected path: %q", r.URL.Path)
		}
		var payload map[string]any
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode payload: %v", err)
		}
		results, ok := payload["results"].([]any)
		if !ok || len(results) != 15 {
			t.Fatalf("unexpected results: %#v", payload["results"])
		}
		wantTypes := []string{"photo", "gif", "mpeg4_gif", "video", "audio", "voice", "document", "photo", "gif", "mpeg4_gif", "sticker", "document", "video", "voice", "audio"}
		for index, want := range wantTypes {
			result, _ := results[index].(map[string]any)
			if result["type"] != want {
				t.Fatalf("unexpected result %d type: %#v", index, result)
			}
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true,"result":true}`))
	}))
	defer server.Close()

	bot := newTestBot(t, token, server.URL, server.Client())
	ok, err := bot.AnswerInlineQuery(context.Background(), AnswerInlineQueryParams{
		InlineQueryID: "inline-query-id",
		Results: []InlineQueryResult{
			InlinePhoto("photo-1", "https://example.com/photo.jpg", "https://example.com/thumb.jpg"),
			InlineGif("gif-1", "https://example.com/anim.gif", "https://example.com/thumb.jpg"),
			InlineMpeg4Gif("mpeg4-1", "https://example.com/anim.mp4", "https://example.com/thumb.jpg"),
			InlineVideo("video-1", "https://example.com/video.mp4", "video/mp4", "https://example.com/thumb.jpg", "Video"),
			InlineAudio("audio-1", "https://example.com/audio.mp3", "Audio"),
			InlineVoice("voice-1", "https://example.com/voice.ogg", "Voice"),
			InlineDocument("doc-1", "Document", "https://example.com/file.pdf", "application/pdf"),
			InlineCachedPhoto("cached-photo-1", "photo-file-id"),
			InlineCachedGif("cached-gif-1", "gif-file-id"),
			InlineCachedMpeg4Gif("cached-mpeg4-1", "mpeg4-file-id"),
			InlineCachedSticker("cached-sticker-1", "sticker-file-id"),
			InlineCachedDocument("cached-doc-1", "doc-file-id", "Document"),
			InlineCachedVideo("cached-video-1", "video-file-id", "Video"),
			InlineCachedVoice("cached-voice-1", "voice-file-id", "Voice"),
			InlineCachedAudio("cached-audio-1", "audio-file-id"),
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Fatal("expected true result")
	}
}
