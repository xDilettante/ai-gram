// Package bot contains the primary Telegram Bot API client facade.
package bot

import (
	"bytes"
	"context"
	"encoding/json"
	stderrors "errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	apierrors "github.com/xDilettante/ai-gram/errors"
	"github.com/xDilettante/ai-gram/internal/httpclient"
)

const (
	defaultBaseURL     = "https://api.telegram.org"
	defaultFileBaseURL = "https://api.telegram.org/file"
)

var errEmptyToken = stderrors.New("bot token is required")

// Bot is the primary client for Telegram Bot API operations.
//
// Bot is safe to share between goroutines after construction. The token is stored privately and is not included in errors.
type Bot struct {
	token       string
	baseURL     string
	fileBaseURL string
	client      *httpclient.Client
}

// BotConfig configures a Bot.
type BotConfig struct {
	// Token is the Telegram bot token. It is required and is stored privately by Bot.
	Token string
	// BaseURL is an optional Telegram Bot API base URL override for tests or compatible servers.
	BaseURL string
	// FileBaseURL is an optional Telegram Bot API file download base URL override.
	FileBaseURL string
	// HTTPClient is an optional HTTP client used for future API calls.
	HTTPClient *http.Client
}

// New creates a Bot from config.
func New(config BotConfig) (*Bot, error) {
	if strings.TrimSpace(config.Token) == "" {
		return nil, errEmptyToken
	}

	baseURL := strings.TrimRight(config.BaseURL, "/")
	baseURLExplicit := baseURL != ""
	if baseURL == "" {
		baseURL = defaultBaseURL
	}
	if err := validateBaseURL(baseURL); err != nil {
		return nil, err
	}

	fileBaseURL := strings.TrimRight(config.FileBaseURL, "/")
	if fileBaseURL == "" {
		if baseURLExplicit {
			fileBaseURL = baseURL + "/file"
		} else {
			fileBaseURL = defaultFileBaseURL
		}
	}
	if err := validateBaseURL(fileBaseURL); err != nil {
		return nil, err
	}

	return &Bot{
		token:       config.Token,
		baseURL:     baseURL,
		fileBaseURL: fileBaseURL,
		client:      httpclient.New(config.HTTPClient),
	}, nil
}

type telegramResponse struct {
	OK          bool                          `json:"ok"`
	Result      json.RawMessage               `json:"result,omitempty"`
	ErrorCode   int                           `json:"error_code,omitempty"`
	Description string                        `json:"description,omitempty"`
	Parameters  *apierrors.ResponseParameters `json:"parameters,omitempty"`
}

func (b *Bot) call(ctx context.Context, method string, payload any, result any) error {
	if b == nil {
		return stderrors.New("bot is required")
	}
	if ctx == nil {
		return stderrors.New("context is required")
	}
	if strings.TrimSpace(method) == "" {
		return stderrors.New("telegram method is required")
	}

	body, err := encodePayload(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, b.endpoint(method), bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("create telegram API request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	responseBody, err := b.client.Do(ctx, req)
	if err != nil {
		return err
	}

	var response telegramResponse
	if err := json.Unmarshal(responseBody, &response); err != nil {
		return fmt.Errorf("decode telegram API response: %w", err)
	}
	if !response.OK {
		return &apierrors.APIError{
			Code:        response.ErrorCode,
			Description: b.redactToken(response.Description),
			Parameters:  response.Parameters,
		}
	}
	if result == nil || len(response.Result) == 0 {
		return nil
	}
	if err := json.Unmarshal(response.Result, result); err != nil {
		return fmt.Errorf("decode telegram API result: %w", err)
	}

	return nil
}

func (b *Bot) endpoint(method string) string {
	return b.baseURL + "/bot" + b.token + "/" + strings.TrimLeft(method, "/")
}

func encodePayload(payload any) ([]byte, error) {
	if payload == nil {
		return []byte("{}"), nil
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("encode telegram API request: %w", err)
	}

	return body, nil
}

func (b *Bot) redactToken(message string) string {
	if b == nil || b.token == "" || message == "" {
		return message
	}

	return strings.ReplaceAll(message, b.token, "[redacted]")
}

func validateBaseURL(baseURL string) error {
	parsed, err := url.Parse(baseURL)
	if err != nil {
		return fmt.Errorf("invalid telegram API base URL: %w", err)
	}
	if parsed.Scheme == "" || parsed.Host == "" {
		return stderrors.New("invalid telegram API base URL")
	}

	return nil
}

// String returns a redacted Bot representation safe for fmt.Stringer use.
func (b *Bot) String() string {
	return "bot.Bot{token:[redacted]}"
}

// GoString returns a redacted Bot representation safe for %#v formatting.
func (b *Bot) GoString() string {
	return b.String()
}
