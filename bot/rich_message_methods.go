package bot

import (
	"context"
	stderrors "errors"
	"strings"

	"github.com/xDilettante/ai-gram/telegram"
)

// InputRichMessage describes a rich message to be sent.
type InputRichMessage struct {
	HTML                string `json:"html,omitempty"`
	Markdown            string `json:"markdown,omitempty"`
	IsRTL               bool   `json:"is_rtl,omitempty"`
	SkipEntityDetection bool   `json:"skip_entity_detection,omitempty"`
}

// InputRichMessageContent describes rich message content used by inline results.
type InputRichMessageContent struct {
	RichMessage InputRichMessage `json:"rich_message"`
}

// SendRichMessageParams contains supported parameters for sendRichMessage.
type SendRichMessageParams struct {
	BusinessConnectionID    string                            `json:"business_connection_id,omitempty"`
	ChatID                  ChatID                            `json:"chat_id"`
	MessageThreadID         int64                             `json:"message_thread_id,omitempty"`
	DirectMessagesTopicID   int64                             `json:"direct_messages_topic_id,omitempty"`
	RichMessage             InputRichMessage                  `json:"rich_message"`
	DisableNotification     bool                              `json:"disable_notification,omitempty"`
	ProtectContent          bool                              `json:"protect_content,omitempty"`
	AllowPaidBroadcast      bool                              `json:"allow_paid_broadcast,omitempty"`
	MessageEffectID         string                            `json:"message_effect_id,omitempty"`
	SuggestedPostParameters *telegram.SuggestedPostParameters `json:"suggested_post_parameters,omitempty"`
	ReplyParameters         *telegram.ReplyParameters         `json:"reply_parameters,omitempty"`
	ReplyMarkup             telegram.ReplyMarkup              `json:"reply_markup,omitempty"`
}

// SendRichMessageDraftParams contains supported parameters for sendRichMessageDraft.
type SendRichMessageDraftParams struct {
	ChatID          ChatID           `json:"chat_id"`
	MessageThreadID int64            `json:"message_thread_id,omitempty"`
	DraftID         int64            `json:"draft_id"`
	RichMessage     InputRichMessage `json:"rich_message"`
}

func (InputRichMessageContent) inputMessageContent() {}

// InputRichHTML creates an HTML rich message input.
func InputRichHTML(html string) InputRichMessage {
	return InputRichMessage{HTML: html}
}

// InputRichMarkdown creates a Markdown rich message input.
func InputRichMarkdown(markdown string) InputRichMessage {
	return InputRichMessage{Markdown: markdown}
}

// SendRichMessage sends a rich message.
func (b *Bot) SendRichMessage(ctx context.Context, params SendRichMessageParams) (*telegram.Message, error) {
	if err := params.validate(); err != nil {
		return nil, err
	}

	var message telegram.Message
	if err := b.call(ctx, "sendRichMessage", params, &message); err != nil {
		return nil, err
	}

	return &message, nil
}

// SendRichMessageDraft streams a partial rich message draft to a private chat.
func (b *Bot) SendRichMessageDraft(ctx context.Context, params SendRichMessageDraftParams) (bool, error) {
	if err := params.validate(); err != nil {
		return false, err
	}

	var result bool
	if err := b.call(ctx, "sendRichMessageDraft", params, &result); err != nil {
		return false, err
	}

	return result, nil
}

func (params SendRichMessageParams) validate() error {
	if !params.ChatID.valid() {
		return stderrors.New("chat_id is required")
	}
	if err := validateMessageThreadID(params.MessageThreadID); err != nil {
		return err
	}
	if params.DirectMessagesTopicID < 0 {
		return stderrors.New("direct_messages_topic_id must not be negative")
	}
	if err := validateInputRichMessage(params.RichMessage); err != nil {
		return err
	}
	if err := validateSuggestedPostParameters(params.SuggestedPostParameters); err != nil {
		return err
	}
	if err := validateReplyParameters(params.ReplyParameters); err != nil {
		return err
	}
	if err := telegram.ValidateReplyMarkup(params.ReplyMarkup); err != nil {
		return err
	}
	return nil
}

func (params SendRichMessageDraftParams) validate() error {
	if params.ChatID.kind != chatIDInt || params.ChatID.intID == 0 {
		return stderrors.New("chat_id must be a non-zero integer")
	}
	if err := validateMessageThreadID(params.MessageThreadID); err != nil {
		return err
	}
	if params.DraftID == 0 {
		return stderrors.New("draft_id must be non-zero")
	}
	return validateInputRichMessage(params.RichMessage)
}

func validateInputRichMessage(message InputRichMessage) error {
	htmlSet := strings.TrimSpace(message.HTML) != ""
	markdownSet := strings.TrimSpace(message.Markdown) != ""
	switch {
	case htmlSet && markdownSet:
		return stderrors.New("rich_message html and markdown are mutually exclusive")
	case !htmlSet && !markdownSet:
		return stderrors.New("rich_message html or markdown is required")
	default:
		return nil
	}
}

func (message InputRichMessage) hasContent() bool {
	return strings.TrimSpace(message.HTML) != "" || strings.TrimSpace(message.Markdown) != ""
}
