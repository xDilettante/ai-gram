package bot

import (
	"context"
	stderrors "errors"

	"github.com/xDilettante/ai-gram/telegram"
)

// SetPassportDataErrorsParams contains supported parameters for setPassportDataErrors.
type SetPassportDataErrorsParams struct {
	UserID int64                           `json:"user_id"`
	Errors []telegram.PassportElementError `json:"errors"`
}

// SetPassportDataErrors informs a user about errors in submitted Telegram Passport elements.
func (b *Bot) SetPassportDataErrors(ctx context.Context, params SetPassportDataErrorsParams) (bool, error) {
	if err := params.validate(); err != nil {
		return false, err
	}

	var result bool
	if err := b.call(ctx, "setPassportDataErrors", params, &result); err != nil {
		return false, err
	}

	return result, nil
}

func (params SetPassportDataErrorsParams) validate() error {
	if params.UserID <= 0 {
		return stderrors.New("user_id must be greater than zero")
	}
	if len(params.Errors) == 0 {
		return stderrors.New("errors must not be empty")
	}
	return telegram.ValidatePassportElementErrors(params.Errors)
}
