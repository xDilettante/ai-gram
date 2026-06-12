package bot

import (
	"context"
	stderrors "errors"
	"strings"
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

// ChatJoinRequestQueryResult identifies an answerChatJoinRequestQuery result.
type ChatJoinRequestQueryResult string

const (
	// ChatJoinRequestQueryApprove approves the user join request.
	ChatJoinRequestQueryApprove ChatJoinRequestQueryResult = "approve"
	// ChatJoinRequestQueryDecline declines the user join request.
	ChatJoinRequestQueryDecline ChatJoinRequestQueryResult = "decline"
	// ChatJoinRequestQueryQueue leaves the decision to other administrators.
	ChatJoinRequestQueryQueue ChatJoinRequestQueryResult = "queue"
)

// AnswerChatJoinRequestQueryParams contains supported parameters for answerChatJoinRequestQuery.
type AnswerChatJoinRequestQueryParams struct {
	ChatJoinRequestQueryID string                     `json:"chat_join_request_query_id"`
	Result                 ChatJoinRequestQueryResult `json:"result"`
}

// SendChatJoinRequestWebAppParams contains supported parameters for sendChatJoinRequestWebApp.
type SendChatJoinRequestWebAppParams struct {
	ChatJoinRequestQueryID string `json:"chat_join_request_query_id"`
	WebAppURL              string `json:"web_app_url"`
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

// AnswerChatJoinRequestQuery processes a received chat join request query.
func (b *Bot) AnswerChatJoinRequestQuery(ctx context.Context, params AnswerChatJoinRequestQueryParams) (bool, error) {
	if err := params.validate(); err != nil {
		return false, err
	}

	var result bool
	if err := b.call(ctx, "answerChatJoinRequestQuery", params, &result); err != nil {
		return false, err
	}

	return result, nil
}

// SendChatJoinRequestWebApp processes a chat join request query by opening a Mini App.
func (b *Bot) SendChatJoinRequestWebApp(ctx context.Context, params SendChatJoinRequestWebAppParams) (bool, error) {
	if err := params.validate(); err != nil {
		return false, err
	}

	var result bool
	if err := b.call(ctx, "sendChatJoinRequestWebApp", params, &result); err != nil {
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

func (params AnswerChatJoinRequestQueryParams) validate() error {
	if strings.TrimSpace(params.ChatJoinRequestQueryID) == "" {
		return stderrors.New("chat_join_request_query_id is required")
	}
	switch params.Result {
	case ChatJoinRequestQueryApprove, ChatJoinRequestQueryDecline, ChatJoinRequestQueryQueue:
		return nil
	default:
		return stderrors.New("result must be approve, decline, or queue")
	}
}

func (params SendChatJoinRequestWebAppParams) validate() error {
	if strings.TrimSpace(params.ChatJoinRequestQueryID) == "" {
		return stderrors.New("chat_join_request_query_id is required")
	}
	if err := validateRequiredInlineHTTPURL(params.WebAppURL, "web_app_url"); err != nil {
		return err
	}
	return nil
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
