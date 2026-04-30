// Example v02_send_methods runs safe live smoke checks for v0.2 send methods.
//
// Required env: AIGRAM_BOT_TOKEN and AIGRAM_V02_SMOKE_CHAT_ID or AIGRAM_CHAT_ID.
// Optional env: AIGRAM_BASE_URL, AIGRAM_FILE_BASE_URL, AIGRAM_STICKER_FILE_ID,
// AIGRAM_ANIMATION_FILE_ID, AIGRAM_ANIMATION_PATH, AIGRAM_VIDEO_NOTE_FILE_ID,
// AIGRAM_VIDEO_NOTE_PATH.
package main

import (
	"context"
	"fmt"
	"log"
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
		log.Fatal(err)
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

	contactMessage, err := b.SendContact(ctx, aigram.SendContactParams{
		ChatID:      chatID,
		PhoneNumber: "+10000000000",
		FirstName:   "ai-gram",
		LastName:    "Smoke",
	})
	if err != nil {
		return fmt.Errorf("send contact: %w", err)
	}
	fmt.Printf("AIGRAM_V02_SMOKE_SEND_CONTACT_OK message_id=%d\n", contactMessage.MessageID)

	locationMessage, err := b.SendLocation(ctx, aigram.SendLocationParams{
		ChatID:    chatID,
		Latitude:  52.3676,
		Longitude: 4.9041,
	})
	if err != nil {
		return fmt.Errorf("send location: %w", err)
	}
	fmt.Printf("AIGRAM_V02_SMOKE_SEND_LOCATION_OK message_id=%d\n", locationMessage.MessageID)

	venueMessage, err := b.SendVenue(ctx, aigram.SendVenueParams{
		ChatID:    chatID,
		Latitude:  52.3676,
		Longitude: 4.9041,
		Title:     "ai-gram smoke venue",
		Address:   "ai-gram test address",
	})
	if err != nil {
		return fmt.Errorf("send venue: %w", err)
	}
	fmt.Printf("AIGRAM_V02_SMOKE_SEND_VENUE_OK message_id=%d\n", venueMessage.MessageID)

	pollMessage, err := b.SendPoll(ctx, aigram.SendPollParams{
		ChatID:   chatID,
		Question: "ai-gram smoke poll?",
		Options:  []string{"yes", "no"},
	})
	if err != nil {
		return fmt.Errorf("send poll: %w", err)
	}
	fmt.Printf("AIGRAM_V02_SMOKE_SEND_POLL_OK message_id=%d\n", pollMessage.MessageID)

	poll, err := b.StopPoll(ctx, aigram.StopPollParams{ChatID: chatID, MessageID: pollMessage.MessageID})
	if err != nil {
		return fmt.Errorf("stop poll: %w", err)
	}
	fmt.Printf("AIGRAM_V02_SMOKE_STOP_POLL_OK poll_id=%s\n", safeValue(poll.ID))

	diceMessage, err := b.SendDice(ctx, aigram.SendDiceParams{ChatID: chatID, Emoji: "🎲"})
	if err != nil {
		return fmt.Errorf("send dice: %w", err)
	}
	fmt.Printf("AIGRAM_V02_SMOKE_SEND_DICE_OK message_id=%d\n", diceMessage.MessageID)

	if err := smokeSticker(ctx, b, chatID); err != nil {
		return err
	}
	if err := smokeAnimation(ctx, b, chatID); err != nil {
		return err
	}
	if err := smokeVideoNote(ctx, b, chatID); err != nil {
		return err
	}

	fmt.Println("AIGRAM_V02_SMOKE_OK")
	return nil
}

func smokeChatID() (aigram.ChatID, error) {
	if value := strings.TrimSpace(os.Getenv("AIGRAM_V02_SMOKE_CHAT_ID")); value != "" {
		return exampleutil.ParseChatID(value)
	}
	return exampleutil.ParseChatID(os.Getenv("AIGRAM_CHAT_ID"))
}

func smokeSticker(ctx context.Context, b *aigram.Bot, chatID aigram.ChatID) error {
	fileID := strings.TrimSpace(os.Getenv("AIGRAM_STICKER_FILE_ID"))
	if fileID == "" {
		fmt.Println("AIGRAM_V02_SMOKE_SEND_STICKER_SKIPPED reason=no_file_id")
		return nil
	}
	message, err := b.SendSticker(ctx, aigram.SendStickerParams{ChatID: chatID, Sticker: aigram.FileID(fileID)})
	if err != nil {
		return fmt.Errorf("send sticker: %w", err)
	}
	fmt.Printf("AIGRAM_V02_SMOKE_SEND_STICKER_OK message_id=%d\n", message.MessageID)
	return nil
}

func smokeAnimation(ctx context.Context, b *aigram.Bot, chatID aigram.ChatID) error {
	fileID := strings.TrimSpace(os.Getenv("AIGRAM_ANIMATION_FILE_ID"))
	path := strings.TrimSpace(os.Getenv("AIGRAM_ANIMATION_PATH"))
	if fileID == "" && path == "" {
		fmt.Println("AIGRAM_V02_SMOKE_SEND_ANIMATION_SKIPPED reason=no_media")
		return nil
	}

	params := aigram.SendAnimationParams{ChatID: chatID, Caption: "ai-gram animation smoke"}
	closeFile, err := setFileRef(&params.Animation, fileID, path)
	if err != nil {
		return fmt.Errorf("prepare animation media: %w", err)
	}
	if closeFile != nil {
		defer closeFile()
	}

	message, err := b.SendAnimation(ctx, params)
	if err != nil {
		return fmt.Errorf("send animation: %w", err)
	}
	fmt.Printf("AIGRAM_V02_SMOKE_SEND_ANIMATION_OK message_id=%d\n", message.MessageID)
	return nil
}

func smokeVideoNote(ctx context.Context, b *aigram.Bot, chatID aigram.ChatID) error {
	fileID := strings.TrimSpace(os.Getenv("AIGRAM_VIDEO_NOTE_FILE_ID"))
	path := strings.TrimSpace(os.Getenv("AIGRAM_VIDEO_NOTE_PATH"))
	if fileID == "" && path == "" {
		fmt.Println("AIGRAM_V02_SMOKE_SEND_VIDEO_NOTE_SKIPPED reason=no_media")
		return nil
	}

	params := aigram.SendVideoNoteParams{ChatID: chatID}
	closeFile, err := setFileRef(&params.VideoNote, fileID, path)
	if err != nil {
		return fmt.Errorf("prepare video note media: %w", err)
	}
	if closeFile != nil {
		defer closeFile()
	}

	message, err := b.SendVideoNote(ctx, params)
	if err != nil {
		return fmt.Errorf("send video note: %w", err)
	}
	fmt.Printf("AIGRAM_V02_SMOKE_SEND_VIDEO_NOTE_OK message_id=%d\n", message.MessageID)
	return nil
}

func setFileRef(dst *aigram.FileRef, fileID string, path string) (func(), error) {
	if strings.TrimSpace(fileID) != "" {
		*dst = aigram.FileID(fileID)
		return nil, nil
	}
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	*dst = aigram.FileUpload(aigram.UploadFile{
		Name:        filepath.Base(path),
		Reader:      file,
		ContentType: contentTypeForPath(path),
	})
	return func() { _ = file.Close() }, nil
}

func contentTypeForPath(path string) string {
	if contentType := mime.TypeByExtension(filepath.Ext(path)); contentType != "" {
		return contentType
	}
	return "application/octet-stream"
}

func safeValue(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return "unknown"
	}
	if len(value) <= 16 {
		return value
	}
	return value[:8] + "..." + value[len(value)-4:]
}
