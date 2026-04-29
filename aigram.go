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

// NewBot creates a Bot from config.
func NewBot(config BotConfig) (*Bot, error) {
	return bot.New(config)
}
