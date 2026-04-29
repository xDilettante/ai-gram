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

// GetUpdatesParams contains supported parameters for getUpdates.
type GetUpdatesParams = bot.GetUpdatesParams

// GetFileParams contains supported parameters for getFile.
type GetFileParams = bot.GetFileParams

// SetWebhookParams contains supported parameters for setWebhook.
type SetWebhookParams = bot.SetWebhookParams

// DeleteWebhookParams contains supported parameters for deleteWebhook.
type DeleteWebhookParams = bot.DeleteWebhookParams

// ChatID identifies a Telegram chat by numeric ID or username string.
type ChatID = bot.ChatID

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
