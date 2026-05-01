package bot

import (
	"context"
	stderrors "errors"
	"strings"

	"github.com/xDilettante/ai-gram/telegram"
)

// SendChecklistParams contains supported parameters for sendChecklist.
type SendChecklistParams struct {
	BusinessConnectionID string                         `json:"business_connection_id"`
	ChatID               ChatID                         `json:"chat_id"`
	Checklist            telegram.InputChecklist        `json:"checklist"`
	DisableNotification  bool                           `json:"disable_notification,omitempty"`
	ProtectContent       bool                           `json:"protect_content,omitempty"`
	MessageEffectID      string                         `json:"message_effect_id,omitempty"`
	ReplyParameters      *telegram.ReplyParameters      `json:"reply_parameters,omitempty"`
	ReplyMarkup          *telegram.InlineKeyboardMarkup `json:"reply_markup,omitempty"`
}

// EditMessageChecklistParams contains supported parameters for editMessageChecklist.
type EditMessageChecklistParams struct {
	BusinessConnectionID string                         `json:"business_connection_id"`
	ChatID               ChatID                         `json:"chat_id"`
	MessageID            int64                          `json:"message_id"`
	Checklist            telegram.InputChecklist        `json:"checklist"`
	ReplyMarkup          *telegram.InlineKeyboardMarkup `json:"reply_markup,omitempty"`
}

// SendMessageDraftParams contains supported parameters for sendMessageDraft.
type SendMessageDraftParams struct {
	ChatID          ChatID                   `json:"chat_id"`
	MessageThreadID int64                    `json:"message_thread_id,omitempty"`
	DraftID         int64                    `json:"draft_id"`
	Text            string                   `json:"text"`
	ParseMode       string                   `json:"parse_mode,omitempty"`
	Entities        []telegram.MessageEntity `json:"entities,omitempty"`
}

// SendChecklist sends a checklist on behalf of a connected business account.
func (b *Bot) SendChecklist(ctx context.Context, params SendChecklistParams) (*telegram.Message, error) {
	if err := params.validate(); err != nil {
		return nil, err
	}

	var message telegram.Message
	if err := b.call(ctx, "sendChecklist", params, &message); err != nil {
		return nil, err
	}

	return &message, nil
}

// EditMessageChecklist edits a checklist on behalf of a connected business account.
func (b *Bot) EditMessageChecklist(ctx context.Context, params EditMessageChecklistParams) (*telegram.Message, error) {
	if err := params.validate(); err != nil {
		return nil, err
	}

	var message telegram.Message
	if err := b.call(ctx, "editMessageChecklist", params, &message); err != nil {
		return nil, err
	}

	return &message, nil
}

// SendMessageDraft streams a partial message draft to a private chat.
func (b *Bot) SendMessageDraft(ctx context.Context, params SendMessageDraftParams) (bool, error) {
	if err := params.validate(); err != nil {
		return false, err
	}

	var result bool
	if err := b.call(ctx, "sendMessageDraft", params, &result); err != nil {
		return false, err
	}

	return result, nil
}

func (params SendChecklistParams) validate() error {
	if strings.TrimSpace(params.BusinessConnectionID) == "" {
		return stderrors.New("business_connection_id is required")
	}
	if !params.ChatID.valid() {
		return stderrors.New("chat_id is required")
	}
	if err := validateInputChecklist(params.Checklist); err != nil {
		return err
	}
	if err := validateReplyParameters(params.ReplyParameters); err != nil {
		return err
	}
	if params.ReplyMarkup != nil {
		if err := telegram.ValidateReplyMarkup(*params.ReplyMarkup); err != nil {
			return err
		}
	}

	return nil
}

func (params EditMessageChecklistParams) validate() error {
	if strings.TrimSpace(params.BusinessConnectionID) == "" {
		return stderrors.New("business_connection_id is required")
	}
	if !params.ChatID.valid() {
		return stderrors.New("chat_id is required")
	}
	if params.MessageID <= 0 {
		return stderrors.New("message_id must be greater than zero")
	}
	if err := validateInputChecklist(params.Checklist); err != nil {
		return err
	}
	if params.ReplyMarkup != nil {
		if err := telegram.ValidateReplyMarkup(*params.ReplyMarkup); err != nil {
			return err
		}
	}

	return nil
}

func (params SendMessageDraftParams) validate() error {
	if !params.ChatID.valid() {
		return stderrors.New("chat_id is required")
	}
	if err := validateMessageThreadID(params.MessageThreadID); err != nil {
		return err
	}
	if params.DraftID == 0 {
		return stderrors.New("draft_id is required")
	}
	if strings.TrimSpace(params.Text) == "" {
		return stderrors.New("text is required")
	}
	if err := validateEntityFormatting(params.ParseMode, params.Entities); err != nil {
		return err
	}

	return nil
}

func validateInputChecklist(checklist telegram.InputChecklist) error {
	if strings.TrimSpace(checklist.Title) == "" {
		return stderrors.New("checklist.title is required")
	}
	if err := validateEntityFormatting(checklist.ParseMode, checklist.TitleEntities); err != nil {
		return err
	}
	if len(checklist.Tasks) == 0 {
		return stderrors.New("checklist.tasks must contain at least one item")
	}

	seen := make(map[int64]struct{}, len(checklist.Tasks))
	for _, task := range checklist.Tasks {
		if err := validateInputChecklistTask(task); err != nil {
			return err
		}
		if _, ok := seen[task.ID]; ok {
			return stderrors.New("checklist.tasks.id values must be unique")
		}
		seen[task.ID] = struct{}{}
	}

	return nil
}

func validateInputChecklistTask(task telegram.InputChecklistTask) error {
	if task.ID <= 0 {
		return stderrors.New("checklist.tasks.id must be greater than zero")
	}
	if strings.TrimSpace(task.Text) == "" {
		return stderrors.New("checklist.tasks.text is required")
	}
	if err := validateEntityFormatting(task.ParseMode, task.TextEntities); err != nil {
		return err
	}

	return nil
}
