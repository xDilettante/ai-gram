package aigram

import (
	"ai-gram/bot"
	telegramerrors "ai-gram/errors"
	"ai-gram/telegram"
)

// Bot is the primary Telegram Bot API client.
type Bot = bot.Bot

// BotConfig configures a Bot.
type BotConfig = bot.BotConfig

// SendMessageParams contains supported parameters for sendMessage.
type SendMessageParams = bot.SendMessageParams

// SendPhotoParams contains supported parameters for sendPhoto.
type SendPhotoParams = bot.SendPhotoParams

// SendDocumentParams contains supported parameters for sendDocument.
type SendDocumentParams = bot.SendDocumentParams

// SendVideoParams contains supported parameters for sendVideo.
type SendVideoParams = bot.SendVideoParams

// SendAudioParams contains supported parameters for sendAudio.
type SendAudioParams = bot.SendAudioParams

// SendVoiceParams contains supported parameters for sendVoice.
type SendVoiceParams = bot.SendVoiceParams

// GetUpdatesParams contains supported parameters for getUpdates.
type GetUpdatesParams = bot.GetUpdatesParams

// GetFileParams contains supported parameters for getFile.
type GetFileParams = bot.GetFileParams

// SetWebhookParams contains supported parameters for setWebhook.
type SetWebhookParams = bot.SetWebhookParams

// DeleteWebhookParams contains supported parameters for deleteWebhook.
type DeleteWebhookParams = bot.DeleteWebhookParams

// AnswerCallbackQueryParams contains supported parameters for answerCallbackQuery.
type AnswerCallbackQueryParams = bot.AnswerCallbackQueryParams

// ChatID identifies a Telegram chat by numeric ID or username string.
type ChatID = bot.ChatID

// FileRef identifies media by Telegram file_id, HTTP(S) URL, or multipart upload.
type FileRef = bot.FileRef

// UploadFile describes a file uploaded through multipart/form-data.
type UploadFile = bot.UploadFile

// APIError represents a Telegram Bot API response with ok=false.
type APIError = telegramerrors.APIError

// ResponseParameters describes optional Telegram Bot API error parameters.
type ResponseParameters = telegramerrors.ResponseParameters

// Update represents an incoming Telegram update.
type Update = telegram.Update

// Message represents a Telegram message.
type Message = telegram.Message

// User represents a Telegram user or bot account.
type User = telegram.User

// Chat represents a Telegram chat.
type Chat = telegram.Chat

// CallbackQuery represents an incoming callback query.
type CallbackQuery = telegram.CallbackQuery

// WebhookInfo describes current Telegram webhook status.
type WebhookInfo = telegram.WebhookInfo

// File represents a Telegram file metadata object.
type File = telegram.File

// ReplyMarkup marks Telegram reply markup objects.
type ReplyMarkup = telegram.ReplyMarkup

// InlineKeyboardMarkup represents an inline keyboard attached to a message.
type InlineKeyboardMarkup = telegram.InlineKeyboardMarkup

// InlineKeyboardButton represents one inline keyboard button.
type InlineKeyboardButton = telegram.InlineKeyboardButton

// ReplyKeyboardMarkup represents a custom reply keyboard.
type ReplyKeyboardMarkup = telegram.ReplyKeyboardMarkup

// KeyboardButton represents one custom reply keyboard button.
type KeyboardButton = telegram.KeyboardButton

// ReplyKeyboardRemove requests removal of a custom reply keyboard.
type ReplyKeyboardRemove = telegram.ReplyKeyboardRemove

// ForceReply requests Telegram clients to show a reply interface for the message.
type ForceReply = telegram.ForceReply

// New creates a Bot from config.
func New(config BotConfig) (*Bot, error) {
	return bot.New(config)
}

// NewBot creates a Bot from config.
func NewBot(config BotConfig) (*Bot, error) {
	return New(config)
}

// ChatIDInt creates a numeric chat ID.
func ChatIDInt(id int64) ChatID {
	return bot.ChatIDInt(id)
}

// ChatIDString creates a string chat ID, such as a channel username.
func ChatIDString(id string) ChatID {
	return bot.ChatIDString(id)
}

// FileID creates a file reference from an existing Telegram file_id.
func FileID(id string) FileRef {
	return bot.FileID(id)
}

// FileURL creates a file reference from an HTTP(S) URL.
func FileURL(rawURL string) FileRef {
	return bot.FileURL(rawURL)
}

// FileUpload creates a file reference from an UploadFile for multipart upload.
func FileUpload(file UploadFile) FileRef {
	return bot.FileUpload(file)
}

// NewInlineKeyboard creates an InlineKeyboardMarkup from rows of buttons.
func NewInlineKeyboard(rows ...[]InlineKeyboardButton) InlineKeyboardMarkup {
	return telegram.NewInlineKeyboard(rows...)
}

// InlineButtonURL creates an inline keyboard button that opens an HTTP(S) URL.
func InlineButtonURL(text string, rawURL string) InlineKeyboardButton {
	return telegram.InlineButtonURL(text, rawURL)
}

// InlineButtonCallback creates an inline keyboard button with callback data.
func InlineButtonCallback(text string, data string) InlineKeyboardButton {
	return telegram.InlineButtonCallback(text, data)
}

// NewReplyKeyboard creates a ReplyKeyboardMarkup from rows of buttons.
func NewReplyKeyboard(rows ...[]KeyboardButton) ReplyKeyboardMarkup {
	return telegram.NewReplyKeyboard(rows...)
}

// KeyboardButtonText creates a plain text reply keyboard button.
func KeyboardButtonText(text string) KeyboardButton {
	return telegram.KeyboardButtonText(text)
}

// KeyboardButtonContact creates a reply keyboard button that requests a contact.
func KeyboardButtonContact(text string) KeyboardButton {
	return telegram.KeyboardButtonContact(text)
}

// KeyboardButtonLocation creates a reply keyboard button that requests a location.
func KeyboardButtonLocation(text string) KeyboardButton {
	return telegram.KeyboardButtonLocation(text)
}

// RemoveKeyboard creates a ReplyKeyboardRemove markup.
func RemoveKeyboard(selective bool) ReplyKeyboardRemove {
	return telegram.RemoveKeyboard(selective)
}

// NewForceReply creates a ForceReply markup.
func NewForceReply() ForceReply {
	return telegram.NewForceReply()
}

// ValidateReplyMarkup checks whether markup can be sent to Telegram.
func ValidateReplyMarkup(markup ReplyMarkup) error {
	return telegram.ValidateReplyMarkup(markup)
}
