package bot

import (
	"context"
	stderrors "errors"
	"net/url"
	"strings"

	apierrors "ai-gram/errors"
	"ai-gram/internal/telegramsecret"
	"ai-gram/telegram"
)

// SetWebhookParams contains supported parameters for setWebhook.
type SetWebhookParams struct {
	URL                string   `json:"url"`
	IPAddress          string   `json:"ip_address,omitempty"`
	MaxConnections     int      `json:"max_connections,omitempty"`
	AllowedUpdates     []string `json:"allowed_updates,omitempty"`
	DropPendingUpdates bool     `json:"drop_pending_updates,omitempty"`
	SecretToken        string   `json:"secret_token,omitempty"`
}

// DeleteWebhookParams contains supported parameters for deleteWebhook.
type DeleteWebhookParams struct {
	DropPendingUpdates bool `json:"drop_pending_updates,omitempty"`
}

// SetWebhook sets the bot webhook URL.
func (b *Bot) SetWebhook(ctx context.Context, params SetWebhookParams) (bool, error) {
	if err := params.validate(); err != nil {
		return false, err
	}

	var result bool
	if err := b.call(ctx, "setWebhook", params, &result); err != nil {
		return false, redactSecretInError(err, params.SecretToken)
	}

	return result, nil
}

// DeleteWebhook deletes the bot webhook.
func (b *Bot) DeleteWebhook(ctx context.Context, params DeleteWebhookParams) (bool, error) {
	var result bool
	if err := b.call(ctx, "deleteWebhook", params, &result); err != nil {
		return false, err
	}

	return result, nil
}

// GetWebhookInfo returns current webhook status.
func (b *Bot) GetWebhookInfo(ctx context.Context) (*telegram.WebhookInfo, error) {
	var info telegram.WebhookInfo
	if err := b.call(ctx, "getWebhookInfo", nil, &info); err != nil {
		return nil, err
	}

	return &info, nil
}

func (params SetWebhookParams) validate() error {
	if params.URL == "" {
		return stderrors.New("webhook URL is required")
	}
	parsed, err := url.Parse(params.URL)
	if err != nil || parsed.Scheme != "https" || parsed.Host == "" {
		return stderrors.New("webhook URL must be an absolute https URL")
	}
	if params.MaxConnections < 0 {
		return stderrors.New("max_connections must not be negative")
	}
	if params.MaxConnections > 100 {
		return stderrors.New("max_connections must be between 1 and 100")
	}
	if err := telegramsecret.Validate(params.SecretToken); err != nil {
		return err
	}

	return nil
}

func redactSecretInError(err error, secret string) error {
	if secret == "" || err == nil {
		return err
	}

	var apiErr *apierrors.APIError
	if stderrors.As(err, &apiErr) && apiErr != nil {
		apiErr.Description = strings.ReplaceAll(apiErr.Description, secret, "[redacted]")
	}

	return err
}
