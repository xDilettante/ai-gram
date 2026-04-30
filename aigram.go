package aigram

import (
	"github.com/xDilettante/ai-gram/bot"
	telegramerrors "github.com/xDilettante/ai-gram/errors"
	"github.com/xDilettante/ai-gram/telegram"
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

// DeleteMessageParams contains supported parameters for deleteMessage.
type DeleteMessageParams = bot.DeleteMessageParams

// ForwardMessageParams contains supported parameters for forwardMessage.
type ForwardMessageParams = bot.ForwardMessageParams

// CopyMessageParams contains supported parameters for copyMessage.
type CopyMessageParams = bot.CopyMessageParams

// SendChatActionParams contains supported parameters for sendChatAction.
type SendChatActionParams = bot.SendChatActionParams

// PinChatMessageParams contains supported parameters for pinChatMessage.
type PinChatMessageParams = bot.PinChatMessageParams

// UnpinChatMessageParams contains supported parameters for unpinChatMessage.
type UnpinChatMessageParams = bot.UnpinChatMessageParams

// UnpinAllChatMessagesParams contains supported parameters for unpinAllChatMessages.
type UnpinAllChatMessagesParams = bot.UnpinAllChatMessagesParams

// BanChatMemberParams contains supported parameters for banChatMember.
type BanChatMemberParams = bot.BanChatMemberParams

// UnbanChatMemberParams contains supported parameters for unbanChatMember.
type UnbanChatMemberParams = bot.UnbanChatMemberParams

// RestrictChatMemberParams contains supported parameters for restrictChatMember.
type RestrictChatMemberParams = bot.RestrictChatMemberParams

// GetChatParams contains supported parameters for getChat.
type GetChatParams = bot.GetChatParams

// GetChatMemberParams contains supported parameters for getChatMember.
type GetChatMemberParams = bot.GetChatMemberParams

// GetChatAdministratorsParams contains supported parameters for getChatAdministrators.
type GetChatAdministratorsParams = bot.GetChatAdministratorsParams

// GetChatMemberCountParams contains supported parameters for getChatMemberCount.
type GetChatMemberCountParams = bot.GetChatMemberCountParams

// EditMessageResult contains the result returned by edit message methods.
type EditMessageResult = bot.EditMessageResult

// EditMessageTarget identifies a chat or inline message for edit methods.
type EditMessageTarget = bot.EditMessageTarget

// EditMessageTextParams contains supported parameters for editMessageText.
type EditMessageTextParams = bot.EditMessageTextParams

// EditMessageReplyMarkupParams contains supported parameters for editMessageReplyMarkup.
type EditMessageReplyMarkupParams = bot.EditMessageReplyMarkupParams

// EditMessageCaptionParams contains supported parameters for editMessageCaption.
type EditMessageCaptionParams = bot.EditMessageCaptionParams

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

// ChatMemberStatus identifies a user's membership state in a chat.
type ChatMemberStatus = telegram.ChatMemberStatus

// ChatMember describes a Telegram user's membership and permissions in a chat.
type ChatMember = telegram.ChatMember

// CallbackQuery represents an incoming callback query.
type CallbackQuery = telegram.CallbackQuery

// WebhookInfo describes current Telegram webhook status.
type WebhookInfo = telegram.WebhookInfo

// File represents a Telegram file metadata object.
type File = telegram.File

// ReplyMarkup marks Telegram reply markup objects.
type ReplyMarkup = telegram.ReplyMarkup

// ReplyParameters describes the message being replied to.
type ReplyParameters = telegram.ReplyParameters

// ChatPermissions describes actions a user is allowed to take in a chat.
type ChatPermissions = telegram.ChatPermissions

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

const (
	// ChatMemberStatusCreator means the user owns the chat.
	ChatMemberStatusCreator = telegram.ChatMemberStatusCreator
	// ChatMemberStatusAdministrator means the user is a chat administrator.
	ChatMemberStatusAdministrator = telegram.ChatMemberStatusAdministrator
	// ChatMemberStatusMember means the user is a regular chat member.
	ChatMemberStatusMember = telegram.ChatMemberStatusMember
	// ChatMemberStatusRestricted means the user is restricted in the chat.
	ChatMemberStatusRestricted = telegram.ChatMemberStatusRestricted
	// ChatMemberStatusLeft means the user is not currently a member.
	ChatMemberStatusLeft = telegram.ChatMemberStatusLeft
	// ChatMemberStatusKicked means the user was removed from the chat.
	ChatMemberStatusKicked = telegram.ChatMemberStatusKicked

	// ChatActionTyping tells Telegram clients that the bot is typing.
	ChatActionTyping = bot.ChatActionTyping
	// ChatActionUploadPhoto tells Telegram clients that the bot is uploading a photo.
	ChatActionUploadPhoto = bot.ChatActionUploadPhoto
	// ChatActionRecordVideo tells Telegram clients that the bot is recording a video.
	ChatActionRecordVideo = bot.ChatActionRecordVideo
	// ChatActionUploadVideo tells Telegram clients that the bot is uploading a video.
	ChatActionUploadVideo = bot.ChatActionUploadVideo
	// ChatActionRecordVoice tells Telegram clients that the bot is recording a voice message.
	ChatActionRecordVoice = bot.ChatActionRecordVoice
	// ChatActionUploadVoice tells Telegram clients that the bot is uploading a voice message.
	ChatActionUploadVoice = bot.ChatActionUploadVoice
	// ChatActionUploadDocument tells Telegram clients that the bot is uploading a document.
	ChatActionUploadDocument = bot.ChatActionUploadDocument
	// ChatActionChooseSticker tells Telegram clients that the bot is choosing a sticker.
	ChatActionChooseSticker = bot.ChatActionChooseSticker
	// ChatActionFindLocation tells Telegram clients that the bot is finding a location.
	ChatActionFindLocation = bot.ChatActionFindLocation
	// ChatActionRecordVideoNote tells Telegram clients that the bot is recording a video note.
	ChatActionRecordVideoNote = bot.ChatActionRecordVideoNote
	// ChatActionUploadVideoNote tells Telegram clients that the bot is uploading a video note.
	ChatActionUploadVideoNote = bot.ChatActionUploadVideoNote
)

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

// EditTargetChat creates an edit target for a regular chat message.
func EditTargetChat(chatID ChatID, messageID int64) EditMessageTarget {
	return bot.EditTargetChat(chatID, messageID)
}

// EditTargetInline creates an edit target for an inline message.
func EditTargetInline(inlineMessageID string) EditMessageTarget {
	return bot.EditTargetInline(inlineMessageID)
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
