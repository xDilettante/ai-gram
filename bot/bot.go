// Package bot contains the primary Telegram Bot API client facade.
package bot

import (
	stderrors "errors"
	"net/http"
	"strings"

	"ai-gram/internal/httpclient"
)

const defaultBaseURL = "https://api.telegram.org"

var errEmptyToken = stderrors.New("bot token is required")

// Bot is the primary client for Telegram Bot API operations.
//
// Bot is safe to share between goroutines after construction. The token is stored privately and is not included in errors.
type Bot struct {
	token   string
	baseURL string
	client  *httpclient.Client
}

// BotConfig configures a Bot.
type BotConfig struct {
	// Token is the Telegram bot token. It is required and is stored privately by Bot.
	Token string
	// BaseURL is an optional Telegram Bot API base URL override for tests or compatible servers.
	BaseURL string
	// HTTPClient is an optional HTTP client used for future API calls.
	HTTPClient *http.Client
}

// New creates a Bot from config.
func New(config BotConfig) (*Bot, error) {
	if strings.TrimSpace(config.Token) == "" {
		return nil, errEmptyToken
	}

	baseURL := config.BaseURL
	if baseURL == "" {
		baseURL = defaultBaseURL
	}

	return &Bot{
		token:   config.Token,
		baseURL: baseURL,
		client:  httpclient.New(config.HTTPClient),
	}, nil
}

// Token returns the configured Telegram bot token.
func (b *Bot) Token() string {
	if b == nil {
		return ""
	}

	return b.token
}
