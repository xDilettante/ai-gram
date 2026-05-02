// Example 05_media_upload sends files by file_id or local upload.
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	aigram "github.com/xDilettante/ai-gram"
	"github.com/xDilettante/ai-gram/examples/internal/exampleutil"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	ctx := context.Background()
	b, err := exampleutil.NewBotFromEnv()
	if err != nil {
		return err
	}
	chatID, err := exampleutil.ParseChatID(os.Getenv("AIGRAM_CHAT_ID"))
	if err != nil {
		return fmt.Errorf("AIGRAM_CHAT_ID: %w", err)
	}

	if err := sendPhoto(ctx, b, chatID); err != nil {
		return err
	}
	if err := sendDocument(ctx, b, chatID); err != nil {
		return err
	}
	return nil
}

func sendPhoto(ctx context.Context, b *aigram.Bot, chatID aigram.ChatID) error {
	if fileID := strings.TrimSpace(os.Getenv("AIGRAM_PHOTO_FILE_ID")); fileID != "" {
		_, err := b.SendPhoto(ctx, aigram.SendPhotoParams{ChatID: chatID, Photo: aigram.FileID(fileID), Caption: "Photo sent by file_id"})
		return err
	}
	path := strings.TrimSpace(os.Getenv("AIGRAM_PHOTO_PATH"))
	if path == "" {
		log.Println("skip photo: set AIGRAM_PHOTO_FILE_ID or AIGRAM_PHOTO_PATH")
		return nil
	}
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("open photo: %w", err)
	}
	defer file.Close()
	_, err = b.SendPhoto(ctx, aigram.SendPhotoParams{
		ChatID:  chatID,
		Photo:   aigram.FileUpload(aigram.UploadFile{Name: filepath.Base(path), Reader: file}),
		Caption: "Photo uploaded from a local file",
	})
	return err
}

func sendDocument(ctx context.Context, b *aigram.Bot, chatID aigram.ChatID) error {
	if fileID := strings.TrimSpace(os.Getenv("AIGRAM_DOCUMENT_FILE_ID")); fileID != "" {
		_, err := b.SendDocument(ctx, aigram.SendDocumentParams{ChatID: chatID, Document: aigram.FileID(fileID), Caption: "Document sent by file_id"})
		return err
	}
	path := strings.TrimSpace(os.Getenv("AIGRAM_DOCUMENT_PATH"))
	if path != "" {
		return uploadDocument(ctx, b, chatID, path)
	}

	file, err := os.CreateTemp("", "aigram-example-*.txt")
	if err != nil {
		return err
	}
	path = file.Name()
	defer os.Remove(path)
	if _, err := file.WriteString("Hello from ai-gram document upload example.\n"); err != nil {
		file.Close()
		return err
	}
	if _, err := file.Seek(0, 0); err != nil {
		file.Close()
		return err
	}
	defer file.Close()

	_, err = b.SendDocument(ctx, aigram.SendDocumentParams{
		ChatID:   chatID,
		Document: aigram.FileUpload(aigram.UploadFile{Name: "aigram-example.txt", Reader: file, ContentType: "text/plain"}),
		Caption:  "Generated text document uploaded from a temp file",
	})
	return err
}

func uploadDocument(ctx context.Context, b *aigram.Bot, chatID aigram.ChatID, path string) error {
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("open document: %w", err)
	}
	defer file.Close()
	_, err = b.SendDocument(ctx, aigram.SendDocumentParams{
		ChatID:   chatID,
		Document: aigram.FileUpload(aigram.UploadFile{Name: filepath.Base(path), Reader: file}),
		Caption:  "Document uploaded from a local file",
	})
	return err
}
