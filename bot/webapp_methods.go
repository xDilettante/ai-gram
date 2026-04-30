package bot

import (
	"context"
	stderrors "errors"
	"fmt"
	"strings"

	"github.com/xDilettante/ai-gram/telegram"
)

// AnswerWebAppQueryParams contains supported parameters for answerWebAppQuery.
type AnswerWebAppQueryParams struct {
	WebAppQueryID string            `json:"web_app_query_id"`
	Result        InlineQueryResult `json:"result"`
}

// AnswerWebAppQuery sets the result of a Web App interaction.
func (b *Bot) AnswerWebAppQuery(ctx context.Context, params AnswerWebAppQueryParams) (*telegram.SentWebAppMessage, error) {
	if err := params.validate(); err != nil {
		return nil, err
	}

	var message telegram.SentWebAppMessage
	if err := b.call(ctx, "answerWebAppQuery", params, &message); err != nil {
		return nil, err
	}

	return &message, nil
}

func (params AnswerWebAppQueryParams) validate() error {
	if strings.TrimSpace(params.WebAppQueryID) == "" {
		return stderrors.New("web_app_query_id is required")
	}
	if err := validateInlineQueryResult(params.Result); err != nil {
		return fmt.Errorf("result is invalid: %w", err)
	}
	return nil
}
