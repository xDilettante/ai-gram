package bot

import (
	"context"
	"encoding/json"
	stderrors "errors"
	"net/url"
	"strconv"
	"strings"

	apierrors "github.com/xDilettante/ai-gram/errors"
	"github.com/xDilettante/ai-gram/internal/telegramsecret"
	"github.com/xDilettante/ai-gram/telegram"
)

// SetWebhookParams contains supported parameters for setWebhook.
type SetWebhookParams struct {
	URL                string   `json:"url"`
	Certificate        FileRef  `json:"-"`
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
	if b == nil {
		return false, stderrors.New("bot is required")
	}
	if err := params.validateWithLocalHTTP(b.baseURL != defaultBaseURL); err != nil {
		return false, err
	}

	var result bool
	if params.Certificate.isSet() {
		fields, files, err := params.multipart()
		if err != nil {
			return false, err
		}
		if err := b.callMultipart(ctx, "setWebhook", fields, files, &result); err != nil {
			return false, redactSetWebhookError(err, params.SecretToken, params.URL)
		}

		return result, nil
	}

	if err := b.call(ctx, "setWebhook", params, &result); err != nil {
		return false, redactSetWebhookError(err, params.SecretToken, params.URL)
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
	return params.validateWithLocalHTTP(false)
}

func (params SetWebhookParams) validateWithLocalHTTP(allowHTTP bool) error {
	if params.URL == "" {
		return stderrors.New("webhook URL is required")
	}
	parsed, err := url.Parse(params.URL)
	if err != nil || parsed.Host == "" {
		return stderrors.New("webhook URL must be an absolute http(s) URL")
	}
	if parsed.Scheme != "https" && !(allowHTTP && parsed.Scheme == "http") {
		return stderrors.New("webhook URL must be an absolute https URL unless a custom Bot API base URL is used")
	}
	if params.MaxConnections < 0 {
		return stderrors.New("max_connections must not be negative")
	}
	if params.MaxConnections > 100 {
		return stderrors.New("max_connections must be between 1 and 100")
	}
	if params.Certificate.isSet() {
		if !params.Certificate.isUpload() {
			return stderrors.New("certificate must be uploaded with FileUpload")
		}
		if err := params.Certificate.validate("certificate"); err != nil {
			return err
		}
	}
	if err := telegramsecret.Validate(params.SecretToken); err != nil {
		return err
	}

	return nil
}

func (params SetWebhookParams) multipart() (map[string]string, map[string]UploadFile, error) {
	fields := map[string]string{"url": params.URL}
	if params.IPAddress != "" {
		fields["ip_address"] = params.IPAddress
	}
	if params.MaxConnections != 0 {
		fields["max_connections"] = strconv.Itoa(params.MaxConnections)
	}
	if params.AllowedUpdates != nil {
		body, err := json.Marshal(params.AllowedUpdates)
		if err != nil {
			return nil, nil, err
		}
		fields["allowed_updates"] = string(body)
	}
	if params.DropPendingUpdates {
		fields["drop_pending_updates"] = strconv.FormatBool(params.DropPendingUpdates)
	}
	if params.SecretToken != "" {
		fields["secret_token"] = params.SecretToken
	}

	return fields, map[string]UploadFile{"certificate": params.Certificate.upload}, nil
}

func redactSetWebhookError(err error, secret string, webhookURL string) error {
	if err == nil {
		return err
	}

	var apiErr *apierrors.APIError
	if stderrors.As(err, &apiErr) && apiErr != nil {
		if secret != "" {
			apiErr.Description = strings.ReplaceAll(apiErr.Description, secret, "[redacted]")
		}
		if webhookURL != "" {
			apiErr.Description = strings.ReplaceAll(apiErr.Description, webhookURL, "[redacted]")
		}
	}

	return err
}
