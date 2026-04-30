// Example media_upload sends a local file and downloads a Telegram file into memory.
//
// Required env: AIGRAM_BOT_TOKEN.
// Optional env: AIGRAM_BASE_URL, AIGRAM_FILE_BASE_URL, AIGRAM_CHAT_ID, AIGRAM_MEDIA_PATH, AIGRAM_FILE_ID.
// Set AIGRAM_MEDIA_PATH and AIGRAM_CHAT_ID to upload a document. Set AIGRAM_FILE_ID to download a file.
package main

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/xDilettante/ai-gram"
	"github.com/xDilettante/ai-gram/examples/internal/exampleutil"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	b, err := exampleutil.NewBotFromEnv()
	if err != nil {
		return err
	}

	ctx := context.Background()
	mediaPath := strings.TrimSpace(os.Getenv("AIGRAM_MEDIA_PATH"))
	fileID := strings.TrimSpace(os.Getenv("AIGRAM_FILE_ID"))
	if mediaPath == "" && fileID == "" {
		return fmt.Errorf("set AIGRAM_MEDIA_PATH with AIGRAM_CHAT_ID to upload, or AIGRAM_FILE_ID to download")
	}

	if mediaPath != "" {
		chatID, err := exampleutil.ParseChatID(os.Getenv("AIGRAM_CHAT_ID"))
		if err != nil {
			return fmt.Errorf("AIGRAM_CHAT_ID is required for upload: %w", err)
		}
		file, err := os.Open(mediaPath)
		if err != nil {
			return fmt.Errorf("open media file: %w", err)
		}
		defer file.Close()

		message, err := b.SendDocument(ctx, aigram.SendDocumentParams{
			ChatID: chatID,
			Document: aigram.FileUpload(aigram.UploadFile{
				Name:   filepath.Base(mediaPath),
				Reader: file,
			}),
			Caption: "Uploaded with ai-gram",
		})
		if err != nil {
			return fmt.Errorf("send document: %w", err)
		}
		fmt.Println("uploaded document message_id:", message.MessageID)
	}

	if fileID != "" {
		file, err := b.GetFile(ctx, aigram.GetFileParams{FileID: fileID})
		if err != nil {
			return fmt.Errorf("get file: %w", err)
		}
		var buf bytes.Buffer
		if err := b.DownloadFile(ctx, file.FilePath, &buf); err != nil {
			return fmt.Errorf("download file: %w", err)
		}
		fmt.Println("downloaded bytes:", buf.Len())
	}

	return nil
}
