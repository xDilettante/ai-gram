// Package exampleutil contains small helpers shared by ai-gram examples.
package exampleutil

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"ai-gram"
)

// RequiredEnv returns a non-empty environment variable value.
func RequiredEnv(name string) (string, error) {
	value := strings.TrimSpace(os.Getenv(name))
	if value == "" {
		return "", fmt.Errorf("%s is required", name)
	}
	return value, nil
}

// OptionalEnv returns an environment variable value or fallback when it is empty.
func OptionalEnv(name string, fallback string) string {
	value := strings.TrimSpace(os.Getenv(name))
	if value == "" {
		return fallback
	}
	return value
}

// NewBotFromEnv creates an ai-gram bot using AIGRAM_BOT_TOKEN and optional base URLs.
func NewBotFromEnv() (*aigram.Bot, error) {
	token, err := RequiredEnv("AIGRAM_BOT_TOKEN")
	if err != nil {
		return nil, err
	}

	return aigram.New(aigram.BotConfig{
		Token:       token,
		BaseURL:     strings.TrimSpace(os.Getenv("AIGRAM_BASE_URL")),
		FileBaseURL: strings.TrimSpace(os.Getenv("AIGRAM_FILE_BASE_URL")),
	})
}

// SignalContext returns a context cancelled by SIGINT or SIGTERM.
func SignalContext() (context.Context, context.CancelFunc) {
	return signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
}

// ParseChatID parses a numeric chat ID or username string.
func ParseChatID(raw string) (aigram.ChatID, error) {
	value := strings.TrimSpace(raw)
	if value == "" {
		return aigram.ChatID{}, errors.New("chat ID is required")
	}
	if id, err := strconv.ParseInt(value, 10, 64); err == nil {
		return aigram.ChatIDInt(id), nil
	}
	return aigram.ChatIDString(value), nil
}
