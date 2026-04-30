package bot

import (
	"context"
	stderrors "errors"
)

// ApproveChatJoinRequestParams contains supported parameters for approveChatJoinRequest.
type ApproveChatJoinRequestParams struct {
	ChatID ChatID `json:"chat_id"`
	UserID int64  `json:"user_id"`
}

// DeclineChatJoinRequestParams contains supported parameters for declineChatJoinRequest.
type DeclineChatJoinRequestParams struct {
	ChatID ChatID `json:"chat_id"`
	UserID int64  `json:"user_id"`
}

// ApproveChatJoinRequest approves a pending join request for a chat.
func (b *Bot) ApproveChatJoinRequest(ctx context.Context, params ApproveChatJoinRequestParams) (bool, error) {
	if err := params.validate(); err != nil {
		return false, err
	}

	var result bool
	if err := b.call(ctx, "approveChatJoinRequest", params, &result); err != nil {
		return false, err
	}

	return result, nil
}

// DeclineChatJoinRequest declines a pending join request for a chat.
func (b *Bot) DeclineChatJoinRequest(ctx context.Context, params DeclineChatJoinRequestParams) (bool, error) {
	if err := params.validate(); err != nil {
		return false, err
	}

	var result bool
	if err := b.call(ctx, "declineChatJoinRequest", params, &result); err != nil {
		return false, err
	}

	return result, nil
}

func (params ApproveChatJoinRequestParams) validate() error {
	return validateChatJoinRequestParams(params.ChatID, params.UserID)
}

func (params DeclineChatJoinRequestParams) validate() error {
	return validateChatJoinRequestParams(params.ChatID, params.UserID)
}

func validateChatJoinRequestParams(chatID ChatID, userID int64) error {
	if !chatID.valid() {
		return stderrors.New("chat_id is required")
	}
	if userID <= 0 {
		return stderrors.New("user_id must be greater than zero")
	}
	return nil
}
