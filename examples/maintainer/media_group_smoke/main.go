// Example media_group_smoke runs a live SendMediaGroup smoke check.
//
// Required env: AIGRAM_BOT_TOKEN and AIGRAM_MEDIA_GROUP_CHAT_ID or AIGRAM_CHAT_ID.
// Optional env: AIGRAM_BASE_URL, AIGRAM_MEDIA_GROUP_FILE_ID_1,
// AIGRAM_MEDIA_GROUP_FILE_ID_2, AIGRAM_MEDIA_GROUP_PATH_1,
// AIGRAM_MEDIA_GROUP_PATH_2.
package main

import (
	"context"
	"fmt"
	"mime"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/xDilettante/ai-gram"
	"github.com/xDilettante/ai-gram/examples/internal/exampleutil"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "AIGRAM_MEDIA_GROUP_ERROR %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	b, err := exampleutil.NewBotFromEnv()
	if err != nil {
		return err
	}

	chatID, err := smokeChatID()
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 90*time.Second)
	defer cancel()

	media, mode, cleanup, err := smokeMedia()
	if err != nil {
		return err
	}
	defer cleanup()

	messages, err := b.SendMediaGroup(ctx, aigram.SendMediaGroupParams{
		ChatID: chatID,
		Media:  media,
	})
	if err != nil {
		return fmt.Errorf("send media group: %w", err)
	}
	if len(messages) < 2 {
		return fmt.Errorf("send media group returned %d messages, want at least 2", len(messages))
	}

	marker := "AIGRAM_MEDIA_GROUP_GENERATED_UPLOAD_OK"
	switch mode {
	case "file_id":
		marker = "AIGRAM_MEDIA_GROUP_FILE_ID_OK"
	case "upload":
		marker = "AIGRAM_MEDIA_GROUP_UPLOAD_OK"
	}
	fmt.Printf("%s count=%d first_message_id=%d\n", marker, len(messages), messages[0].MessageID)
	fmt.Printf("AIGRAM_MEDIA_GROUP_OK count=%d first_message_id=%d\n", len(messages), messages[0].MessageID)
	return nil
}

func smokeChatID() (aigram.ChatID, error) {
	if value := strings.TrimSpace(os.Getenv("AIGRAM_MEDIA_GROUP_CHAT_ID")); value != "" {
		return exampleutil.ParseChatID(value)
	}
	return exampleutil.ParseChatID(os.Getenv("AIGRAM_CHAT_ID"))
}

func smokeMedia() ([]aigram.InputMedia, string, func(), error) {
	fileID1 := strings.TrimSpace(os.Getenv("AIGRAM_MEDIA_GROUP_FILE_ID_1"))
	fileID2 := strings.TrimSpace(os.Getenv("AIGRAM_MEDIA_GROUP_FILE_ID_2"))
	path1 := strings.TrimSpace(os.Getenv("AIGRAM_MEDIA_GROUP_PATH_1"))
	path2 := strings.TrimSpace(os.Getenv("AIGRAM_MEDIA_GROUP_PATH_2"))

	if fileID1 != "" || fileID2 != "" {
		if fileID1 == "" || fileID2 == "" {
			return nil, "", nilCleanup, fmt.Errorf("both AIGRAM_MEDIA_GROUP_FILE_ID_1 and AIGRAM_MEDIA_GROUP_FILE_ID_2 are required for file_id mode")
		}
		first := aigram.MediaDocument(aigram.FileID(fileID1))
		first.Caption = "ai-gram media group file_id smoke"
		second := aigram.MediaDocument(aigram.FileID(fileID2))
		return []aigram.InputMedia{first, second}, "file_id", nilCleanup, nil
	}

	if path1 != "" || path2 != "" {
		if path1 == "" || path2 == "" {
			return nil, "", nilCleanup, fmt.Errorf("both AIGRAM_MEDIA_GROUP_PATH_1 and AIGRAM_MEDIA_GROUP_PATH_2 are required for upload mode")
		}
		first, closeFirst, err := fileUploadDocument(path1, "ai-gram media group upload smoke")
		if err != nil {
			return nil, "", nilCleanup, err
		}
		second, closeSecond, err := fileUploadDocument(path2, "")
		if err != nil {
			closeFirst()
			return nil, "", nilCleanup, err
		}
		cleanup := func() {
			closeSecond()
			closeFirst()
		}
		return []aigram.InputMedia{first, second}, "upload", cleanup, nil
	}

	first := aigram.MediaDocument(aigram.FileUpload(aigram.UploadFile{
		Name:        "aigram-media-group-1.txt",
		Reader:      strings.NewReader("ai-gram generated media group smoke document 1\n"),
		ContentType: "text/plain; charset=utf-8",
	}))
	first.Caption = "ai-gram media group generated upload smoke"
	second := aigram.MediaDocument(aigram.FileUpload(aigram.UploadFile{
		Name:        "aigram-media-group-2.txt",
		Reader:      strings.NewReader("ai-gram generated media group smoke document 2\n"),
		ContentType: "text/plain; charset=utf-8",
	}))
	return []aigram.InputMedia{first, second}, "generated_upload", nilCleanup, nil
}

func fileUploadDocument(path string, caption string) (aigram.InputMediaDocument, func(), error) {
	file, err := os.Open(path)
	if err != nil {
		return aigram.InputMediaDocument{}, nilCleanup, fmt.Errorf("open %s: %s", filepath.Base(path), safePathError(err))
	}
	document := aigram.MediaDocument(aigram.FileUpload(aigram.UploadFile{
		Name:        filepath.Base(path),
		Reader:      file,
		ContentType: contentTypeForPath(path),
	}))
	document.Caption = caption
	return document, func() { _ = file.Close() }, nil
}

func contentTypeForPath(path string) string {
	if contentType := mime.TypeByExtension(filepath.Ext(path)); contentType != "" {
		return contentType
	}
	return "application/octet-stream"
}

func safePathError(err error) string {
	if pathErr, ok := err.(*os.PathError); ok && pathErr.Err != nil {
		return pathErr.Err.Error()
	}
	return "failed"
}

func nilCleanup() {}
