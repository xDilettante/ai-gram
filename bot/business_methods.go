package bot

import (
	"context"
	stderrors "errors"
	"strings"

	"github.com/xDilettante/ai-gram/telegram"
)

// GetBusinessConnectionParams contains supported parameters for getBusinessConnection.
type GetBusinessConnectionParams struct {
	BusinessConnectionID string `json:"business_connection_id"`
}

// DeleteBusinessMessagesParams contains supported parameters for deleteBusinessMessages.
type DeleteBusinessMessagesParams struct {
	BusinessConnectionID string  `json:"business_connection_id"`
	MessageIDs           []int64 `json:"message_ids"`
}

// GetBusinessConnection returns information about a business account connection.
func (b *Bot) GetBusinessConnection(ctx context.Context, params GetBusinessConnectionParams) (*telegram.BusinessConnection, error) {
	if err := params.validate(); err != nil {
		return nil, err
	}

	var connection telegram.BusinessConnection
	if err := b.call(ctx, "getBusinessConnection", params, &connection); err != nil {
		return nil, err
	}
	return &connection, nil
}

// DeleteBusinessMessages deletes messages on behalf of a business account.
func (b *Bot) DeleteBusinessMessages(ctx context.Context, params DeleteBusinessMessagesParams) (bool, error) {
	if err := params.validate(); err != nil {
		return false, err
	}

	var result bool
	if err := b.call(ctx, "deleteBusinessMessages", params, &result); err != nil {
		return false, err
	}
	return result, nil
}

func (params GetBusinessConnectionParams) validate() error {
	if strings.TrimSpace(params.BusinessConnectionID) == "" {
		return stderrors.New("business_connection_id is required")
	}
	return nil
}

func (params DeleteBusinessMessagesParams) validate() error {
	if strings.TrimSpace(params.BusinessConnectionID) == "" {
		return stderrors.New("business_connection_id is required")
	}
	return validateBatchMessageIDs(params.MessageIDs)
}
