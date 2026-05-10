// Package aigram provides a small convenience facade over the ai-gram library.
//
// Use the root package for quick-start code and common bot flows. Import
// github.com/xDilettante/ai-gram/bot and github.com/xDilettante/ai-gram/telegram
// directly for the full Telegram Bot API surface.
package aigram

import (
	"github.com/xDilettante/ai-gram/bot"
	"github.com/xDilettante/ai-gram/telegram"
)

// Bot is the primary Telegram Bot API client.
type Bot = bot.Bot

// BotConfig configures a Bot.
type BotConfig = bot.BotConfig

// ChatID identifies a Telegram chat by integer ID or @username.
type ChatID = bot.ChatID

// FileRef identifies media by Telegram file_id, HTTP(S) URL, or multipart upload.
type FileRef = bot.FileRef

// UploadFile describes a file uploaded through multipart/form-data.
type UploadFile = bot.UploadFile

// Update represents an incoming Telegram update.
type Update = telegram.Update

// Message represents a Telegram message.
type Message = telegram.Message

// CallbackQuery represents an incoming callback query.
type CallbackQuery = telegram.CallbackQuery

// User represents a Telegram user or bot account.
type User = telegram.User

// Chat represents a Telegram chat.
type Chat = telegram.Chat

// ReplyParameters describes reply metadata for outgoing messages.
type ReplyParameters = telegram.ReplyParameters

// InlineKeyboardMarkup represents an inline keyboard.
type InlineKeyboardMarkup = telegram.InlineKeyboardMarkup

// InlineKeyboardButton represents one inline keyboard button.
type InlineKeyboardButton = telegram.InlineKeyboardButton

// ReplyKeyboardMarkup represents a custom reply keyboard.
type ReplyKeyboardMarkup = telegram.ReplyKeyboardMarkup

// KeyboardButton represents one reply keyboard button.
type KeyboardButton = telegram.KeyboardButton

// ReplyKeyboardRemove requests removal of a custom reply keyboard.
type ReplyKeyboardRemove = telegram.ReplyKeyboardRemove

// ForceReply requests Telegram clients to show a reply interface for the message.
type ForceReply = telegram.ForceReply

// SendMessageParams contains supported parameters for sendMessage.
type SendMessageParams = bot.SendMessageParams

// GetUpdatesParams contains supported parameters for getUpdates.
type GetUpdatesParams = bot.GetUpdatesParams

// DeleteWebhookParams contains supported parameters for deleteWebhook.
type DeleteWebhookParams = bot.DeleteWebhookParams

// SetWebhookParams contains supported parameters for setWebhook.
type SetWebhookParams = bot.SetWebhookParams

// AnswerCallbackQueryParams contains supported parameters for answerCallbackQuery.
type AnswerCallbackQueryParams = bot.AnswerCallbackQueryParams

// EditMessageTextParams contains supported parameters for editMessageText.
type EditMessageTextParams = bot.EditMessageTextParams

// EditMessageReplyMarkupParams contains supported parameters for editMessageReplyMarkup.
type EditMessageReplyMarkupParams = bot.EditMessageReplyMarkupParams

// EditMessageCaptionParams contains supported parameters for editMessageCaption.
type EditMessageCaptionParams = bot.EditMessageCaptionParams

// EditMessageTarget identifies a message to edit.
type EditMessageTarget = bot.EditMessageTarget

// SendPhotoParams contains supported parameters for sendPhoto.
type SendPhotoParams = bot.SendPhotoParams

// SendDocumentParams contains supported parameters for sendDocument.
type SendDocumentParams = bot.SendDocumentParams

// GetFileParams contains supported parameters for getFile.
type GetFileParams = bot.GetFileParams

// SendMediaGroupParams contains supported parameters for sendMediaGroup.
type SendMediaGroupParams = bot.SendMediaGroupParams

// InputMedia marks media accepted by sendMediaGroup.
type InputMedia = bot.InputMedia

// InputMediaDocument describes a document media group item.
type InputMediaDocument = bot.InputMediaDocument

// DeleteMessageParams contains supported parameters for deleteMessage.
type DeleteMessageParams = bot.DeleteMessageParams

// CopyMessageParams contains supported parameters for copyMessage.
type CopyMessageParams = bot.CopyMessageParams

// ForwardMessageParams contains supported parameters for forwardMessage.
type ForwardMessageParams = bot.ForwardMessageParams

// SendChatActionParams contains supported parameters for sendChatAction.
type SendChatActionParams = bot.SendChatActionParams

// SendContactParams contains supported parameters for sendContact.
type SendContactParams = bot.SendContactParams

// SendLocationParams contains supported parameters for sendLocation.
type SendLocationParams = bot.SendLocationParams

// SendVenueParams contains supported parameters for sendVenue.
type SendVenueParams = bot.SendVenueParams

// SendPollParams contains supported parameters for sendPoll.
type SendPollParams = bot.SendPollParams

// StopPollParams contains supported parameters for stopPoll.
type StopPollParams = bot.StopPollParams

// SendDiceParams contains supported parameters for sendDice.
type SendDiceParams = bot.SendDiceParams

// SendStickerParams contains supported parameters for sendSticker.
type SendStickerParams = bot.SendStickerParams

// SendAnimationParams contains supported parameters for sendAnimation.
type SendAnimationParams = bot.SendAnimationParams

// SendVideoNoteParams contains supported parameters for sendVideoNote.
type SendVideoNoteParams = bot.SendVideoNoteParams

// GetChatParams contains supported parameters for getChat.
type GetChatParams = bot.GetChatParams

// GetChatMemberCountParams contains supported parameters for getChatMemberCount.
type GetChatMemberCountParams = bot.GetChatMemberCountParams

// SetMyCommandsParams contains supported parameters for setMyCommands.
type SetMyCommandsParams = bot.SetMyCommandsParams

// GetMyCommandsParams contains supported parameters for getMyCommands.
type GetMyCommandsParams = bot.GetMyCommandsParams

// GetMyNameParams contains supported parameters for getMyName.
type GetMyNameParams = bot.GetMyNameParams

// GetMyDescriptionParams contains supported parameters for getMyDescription.
type GetMyDescriptionParams = bot.GetMyDescriptionParams

// GetMyShortDescriptionParams contains supported parameters for getMyShortDescription.
type GetMyShortDescriptionParams = bot.GetMyShortDescriptionParams

// ChatActionTyping tells Telegram clients that the bot is typing.
const ChatActionTyping = bot.ChatActionTyping

// New creates a Bot.
func New(config BotConfig) (*Bot, error) {
	return bot.New(config)
}

// NewBot creates a Bot.
func NewBot(config BotConfig) (*Bot, error) {
	return bot.New(config)
}

// ChatIDInt creates a chat ID from an integer Telegram chat ID.
func ChatIDInt(id int64) ChatID {
	return bot.ChatIDInt(id)
}

// ChatIDString creates a chat ID from @username or another string identifier.
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

// EditTargetChat targets a chat message for editing.
func EditTargetChat(chatID ChatID, messageID int64) EditMessageTarget {
	return bot.EditTargetChat(chatID, messageID)
}

// EditTargetInline targets an inline message for editing.
func EditTargetInline(inlineMessageID string) EditMessageTarget {
	return bot.EditTargetInline(inlineMessageID)
}

// MediaDocument creates a document item for sendMediaGroup.
func MediaDocument(media FileRef) InputMediaDocument {
	return bot.MediaDocument(media)
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

// InlineButtonSwitchInlineQueryCurrentChat creates an inline keyboard button that switches inline mode in the current chat.
func InlineButtonSwitchInlineQueryCurrentChat(text string, query string) InlineKeyboardButton {
	return telegram.InlineButtonSwitchInlineQueryCurrentChat(text, query)
}

// InlineButtonCopyText creates an inline keyboard button that copies text to the clipboard.
func InlineButtonCopyText(text string, copyText string) InlineKeyboardButton {
	return telegram.InlineButtonCopyText(text, copyText)
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

// KeyboardButtonPoll creates a reply keyboard button that requests a poll.
func KeyboardButtonPoll(text string, pollType string) KeyboardButton {
	return telegram.KeyboardButtonPoll(text, pollType)
}

// RemoveKeyboard creates a ReplyKeyboardRemove markup.
func RemoveKeyboard(selective bool) ReplyKeyboardRemove {
	return telegram.RemoveKeyboard(selective)
}
