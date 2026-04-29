package telegram

import (
	stderrors "errors"
	"net/url"
	"strings"
)

// ReplyMarkup marks Telegram reply markup objects that can be attached to supported send methods.
type ReplyMarkup interface {
	replyMarkup()
}

// InlineKeyboardMarkup represents an inline keyboard attached to a message.
type InlineKeyboardMarkup struct {
	InlineKeyboard [][]InlineKeyboardButton `json:"inline_keyboard"`
}

// InlineKeyboardButton represents one inline keyboard button.
type InlineKeyboardButton struct {
	Text         string `json:"text"`
	URL          string `json:"url,omitempty"`
	CallbackData string `json:"callback_data,omitempty"`
}

// ReplyKeyboardMarkup represents a custom reply keyboard.
type ReplyKeyboardMarkup struct {
	Keyboard              [][]KeyboardButton `json:"keyboard"`
	IsPersistent          bool               `json:"is_persistent,omitempty"`
	ResizeKeyboard        bool               `json:"resize_keyboard,omitempty"`
	OneTimeKeyboard       bool               `json:"one_time_keyboard,omitempty"`
	InputFieldPlaceholder string             `json:"input_field_placeholder,omitempty"`
	Selective             bool               `json:"selective,omitempty"`
}

// KeyboardButton represents one custom reply keyboard button.
type KeyboardButton struct {
	Text            string `json:"text"`
	RequestContact  bool   `json:"request_contact,omitempty"`
	RequestLocation bool   `json:"request_location,omitempty"`
}

// ReplyKeyboardRemove requests removal of a custom reply keyboard.
type ReplyKeyboardRemove struct {
	RemoveKeyboard bool `json:"remove_keyboard"`
	Selective      bool `json:"selective,omitempty"`
}

// ForceReply requests Telegram clients to show a reply interface for the message.
type ForceReply struct {
	ForceReply            bool   `json:"force_reply"`
	InputFieldPlaceholder string `json:"input_field_placeholder,omitempty"`
	Selective             bool   `json:"selective,omitempty"`
}

func (InlineKeyboardMarkup) replyMarkup() {}
func (ReplyKeyboardMarkup) replyMarkup()  {}
func (ReplyKeyboardRemove) replyMarkup()  {}
func (ForceReply) replyMarkup()           {}

// NewInlineKeyboard creates an InlineKeyboardMarkup from rows of buttons.
func NewInlineKeyboard(rows ...[]InlineKeyboardButton) InlineKeyboardMarkup {
	return InlineKeyboardMarkup{InlineKeyboard: rows}
}

// InlineButtonURL creates an inline keyboard button that opens an HTTP(S) URL.
func InlineButtonURL(text string, rawURL string) InlineKeyboardButton {
	return InlineKeyboardButton{Text: text, URL: rawURL}
}

// InlineButtonCallback creates an inline keyboard button with callback data.
func InlineButtonCallback(text string, data string) InlineKeyboardButton {
	return InlineKeyboardButton{Text: text, CallbackData: data}
}

// NewReplyKeyboard creates a ReplyKeyboardMarkup from rows of buttons.
func NewReplyKeyboard(rows ...[]KeyboardButton) ReplyKeyboardMarkup {
	return ReplyKeyboardMarkup{Keyboard: rows}
}

// KeyboardButtonText creates a plain text reply keyboard button.
func KeyboardButtonText(text string) KeyboardButton {
	return KeyboardButton{Text: text}
}

// KeyboardButtonContact creates a reply keyboard button that requests a contact.
func KeyboardButtonContact(text string) KeyboardButton {
	return KeyboardButton{Text: text, RequestContact: true}
}

// KeyboardButtonLocation creates a reply keyboard button that requests a location.
func KeyboardButtonLocation(text string) KeyboardButton {
	return KeyboardButton{Text: text, RequestLocation: true}
}

// RemoveKeyboard creates a ReplyKeyboardRemove markup.
func RemoveKeyboard(selective bool) ReplyKeyboardRemove {
	return ReplyKeyboardRemove{RemoveKeyboard: true, Selective: selective}
}

// NewForceReply creates a ForceReply markup.
func NewForceReply() ForceReply {
	return ForceReply{ForceReply: true}
}

// ValidateReplyMarkup checks whether markup can be sent to Telegram.
func ValidateReplyMarkup(markup ReplyMarkup) error {
	if markup == nil {
		return nil
	}

	switch value := markup.(type) {
	case InlineKeyboardMarkup:
		return validateInlineKeyboard(value)
	case ReplyKeyboardMarkup:
		return validateReplyKeyboard(value)
	case ReplyKeyboardRemove:
		return validateReplyKeyboardRemove(value)
	case ForceReply:
		return validateForceReply(value)
	default:
		return stderrors.New("unsupported reply markup")
	}
}

func validateInlineKeyboard(markup InlineKeyboardMarkup) error {
	if len(markup.InlineKeyboard) == 0 {
		return stderrors.New("inline keyboard must not be empty")
	}
	for _, row := range markup.InlineKeyboard {
		if len(row) == 0 {
			return stderrors.New("inline keyboard row must not be empty")
		}
		for _, button := range row {
			if err := validateInlineKeyboardButton(button); err != nil {
				return err
			}
		}
	}

	return nil
}

func validateInlineKeyboardButton(button InlineKeyboardButton) error {
	if strings.TrimSpace(button.Text) == "" {
		return stderrors.New("inline keyboard button text is required")
	}
	actions := 0
	if button.URL != "" {
		actions++
		if err := validateHTTPURL(button.URL, "inline keyboard button URL"); err != nil {
			return err
		}
	}
	if button.CallbackData != "" {
		actions++
		if len([]byte(button.CallbackData)) > 64 {
			return stderrors.New("inline keyboard callback_data must be at most 64 bytes")
		}
	}
	if actions != 1 {
		return stderrors.New("inline keyboard button must have exactly one action")
	}

	return nil
}

func validateReplyKeyboard(markup ReplyKeyboardMarkup) error {
	if len(markup.Keyboard) == 0 {
		return stderrors.New("reply keyboard must not be empty")
	}
	if markup.InputFieldPlaceholder != "" && strings.TrimSpace(markup.InputFieldPlaceholder) == "" {
		return stderrors.New("input field placeholder must not be blank")
	}
	for _, row := range markup.Keyboard {
		if len(row) == 0 {
			return stderrors.New("reply keyboard row must not be empty")
		}
		for _, button := range row {
			if err := validateKeyboardButton(button); err != nil {
				return err
			}
		}
	}

	return nil
}

func validateKeyboardButton(button KeyboardButton) error {
	if strings.TrimSpace(button.Text) == "" {
		return stderrors.New("keyboard button text is required")
	}
	if button.RequestContact && button.RequestLocation {
		return stderrors.New("keyboard button cannot request contact and location together")
	}

	return nil
}

func validateReplyKeyboardRemove(markup ReplyKeyboardRemove) error {
	if !markup.RemoveKeyboard {
		return stderrors.New("remove_keyboard must be true")
	}

	return nil
}

func validateForceReply(markup ForceReply) error {
	if !markup.ForceReply {
		return stderrors.New("force_reply must be true")
	}
	if markup.InputFieldPlaceholder != "" && strings.TrimSpace(markup.InputFieldPlaceholder) == "" {
		return stderrors.New("input field placeholder must not be blank")
	}

	return nil
}

func validateHTTPURL(rawURL string, field string) error {
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return stderrors.New(field + " must be a valid HTTP(S) URL")
	}
	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return stderrors.New(field + " scheme must be http or https")
	}
	if parsed.Host == "" {
		return stderrors.New(field + " host is required")
	}

	return nil
}
