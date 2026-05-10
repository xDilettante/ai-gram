package bot

import (
	"context"
	stderrors "errors"
	"fmt"
	"strings"

	"github.com/xDilettante/ai-gram/telegram"
)

// AnswerGuestQueryParams contains supported parameters for answerGuestQuery.
type AnswerGuestQueryParams struct {
	GuestQueryID string            `json:"guest_query_id"`
	Result       InlineQueryResult `json:"result"`
}

// AnswerGuestQuery replies to a received guest message.
func (b *Bot) AnswerGuestQuery(ctx context.Context, params AnswerGuestQueryParams) (*telegram.SentGuestMessage, error) {
	if err := params.validate(); err != nil {
		return nil, err
	}

	var message telegram.SentGuestMessage
	if err := b.call(ctx, "answerGuestQuery", params, &message); err != nil {
		return nil, err
	}

	return &message, nil
}

func (params AnswerGuestQueryParams) validate() error {
	if strings.TrimSpace(params.GuestQueryID) == "" {
		return stderrors.New("guest_query_id is required")
	}
	if err := validateInlineQueryResult(params.Result); err != nil {
		return fmt.Errorf("result is invalid: %w", err)
	}
	return nil
}
