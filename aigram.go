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

// SendStickerParams contains supported parameters for sendSticker.
type SendStickerParams = bot.SendStickerParams

// SendAnimationParams contains supported parameters for sendAnimation.
type SendAnimationParams = bot.SendAnimationParams

// SendVideoNoteParams contains supported parameters for sendVideoNote.
type SendVideoNoteParams = bot.SendVideoNoteParams

// SendMediaGroupParams contains supported parameters for sendMediaGroup.
type SendMediaGroupParams = bot.SendMediaGroupParams

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

// InputSticker describes a sticker accepted by sticker set management methods.
type InputSticker = bot.InputSticker

// GetStickerSetParams contains supported parameters for getStickerSet.
type GetStickerSetParams = bot.GetStickerSetParams

// GetCustomEmojiStickersParams contains supported parameters for getCustomEmojiStickers.
type GetCustomEmojiStickersParams = bot.GetCustomEmojiStickersParams

// UploadStickerFileParams contains supported parameters for uploadStickerFile.
type UploadStickerFileParams = bot.UploadStickerFileParams

// CreateNewStickerSetParams contains supported parameters for createNewStickerSet.
type CreateNewStickerSetParams = bot.CreateNewStickerSetParams

// AddStickerToSetParams contains supported parameters for addStickerToSet.
type AddStickerToSetParams = bot.AddStickerToSetParams

// ReplaceStickerInSetParams contains supported parameters for replaceStickerInSet.
type ReplaceStickerInSetParams = bot.ReplaceStickerInSetParams

// SetStickerPositionInSetParams contains supported parameters for setStickerPositionInSet.
type SetStickerPositionInSetParams = bot.SetStickerPositionInSetParams

// DeleteStickerFromSetParams contains supported parameters for deleteStickerFromSet.
type DeleteStickerFromSetParams = bot.DeleteStickerFromSetParams

// SetStickerEmojiListParams contains supported parameters for setStickerEmojiList.
type SetStickerEmojiListParams = bot.SetStickerEmojiListParams

// SetStickerKeywordsParams contains supported parameters for setStickerKeywords.
type SetStickerKeywordsParams = bot.SetStickerKeywordsParams

// SetStickerMaskPositionParams contains supported parameters for setStickerMaskPosition.
type SetStickerMaskPositionParams = bot.SetStickerMaskPositionParams

// SetStickerSetTitleParams contains supported parameters for setStickerSetTitle.
type SetStickerSetTitleParams = bot.SetStickerSetTitleParams

// SetStickerSetThumbnailParams contains supported parameters for setStickerSetThumbnail.
type SetStickerSetThumbnailParams = bot.SetStickerSetThumbnailParams

// SetCustomEmojiStickerSetThumbnailParams contains supported parameters for setCustomEmojiStickerSetThumbnail.
type SetCustomEmojiStickerSetThumbnailParams = bot.SetCustomEmojiStickerSetThumbnailParams

// DeleteStickerSetParams contains supported parameters for deleteStickerSet.
type DeleteStickerSetParams = bot.DeleteStickerSetParams

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

// DeleteMessagesParams contains supported parameters for deleteMessages.
type DeleteMessagesParams = bot.DeleteMessagesParams

// ForwardMessageParams contains supported parameters for forwardMessage.
type ForwardMessageParams = bot.ForwardMessageParams

// CopyMessageParams contains supported parameters for copyMessage.
type CopyMessageParams = bot.CopyMessageParams

// ForwardMessagesParams contains supported parameters for forwardMessages.
type ForwardMessagesParams = bot.ForwardMessagesParams

// CopyMessagesParams contains supported parameters for copyMessages.
type CopyMessagesParams = bot.CopyMessagesParams

// SendChatActionParams contains supported parameters for sendChatAction.
type SendChatActionParams = bot.SendChatActionParams

// PinChatMessageParams contains supported parameters for pinChatMessage.
type PinChatMessageParams = bot.PinChatMessageParams

// UnpinChatMessageParams contains supported parameters for unpinChatMessage.
type UnpinChatMessageParams = bot.UnpinChatMessageParams

// UnpinAllChatMessagesParams contains supported parameters for unpinAllChatMessages.
type UnpinAllChatMessagesParams = bot.UnpinAllChatMessagesParams

// SetMessageReactionParams contains supported parameters for setMessageReaction.
type SetMessageReactionParams = bot.SetMessageReactionParams

// BanChatMemberParams contains supported parameters for banChatMember.
type BanChatMemberParams = bot.BanChatMemberParams

// UnbanChatMemberParams contains supported parameters for unbanChatMember.
type UnbanChatMemberParams = bot.UnbanChatMemberParams

// RestrictChatMemberParams contains supported parameters for restrictChatMember.
type RestrictChatMemberParams = bot.RestrictChatMemberParams

// PromoteChatMemberParams contains supported parameters for promoteChatMember.
type PromoteChatMemberParams = bot.PromoteChatMemberParams

// SetChatAdministratorCustomTitleParams contains supported parameters for setChatAdministratorCustomTitle.
type SetChatAdministratorCustomTitleParams = bot.SetChatAdministratorCustomTitleParams

// SetChatPermissionsParams contains supported parameters for setChatPermissions.
type SetChatPermissionsParams = bot.SetChatPermissionsParams

// SetChatTitleParams contains supported parameters for setChatTitle.
type SetChatTitleParams = bot.SetChatTitleParams

// SetChatDescriptionParams contains supported parameters for setChatDescription.
type SetChatDescriptionParams = bot.SetChatDescriptionParams

// SetChatPhotoParams contains supported parameters for setChatPhoto.
type SetChatPhotoParams = bot.SetChatPhotoParams

// DeleteChatPhotoParams contains supported parameters for deleteChatPhoto.
type DeleteChatPhotoParams = bot.DeleteChatPhotoParams

// LeaveChatParams contains supported parameters for leaveChat.
type LeaveChatParams = bot.LeaveChatParams

// SetChatStickerSetParams contains supported parameters for setChatStickerSet.
type SetChatStickerSetParams = bot.SetChatStickerSetParams

// DeleteChatStickerSetParams contains supported parameters for deleteChatStickerSet.
type DeleteChatStickerSetParams = bot.DeleteChatStickerSetParams

// CreateForumTopicParams contains supported parameters for createForumTopic.
type CreateForumTopicParams = bot.CreateForumTopicParams

// EditForumTopicParams contains supported parameters for editForumTopic.
type EditForumTopicParams = bot.EditForumTopicParams

// CloseForumTopicParams contains supported parameters for closeForumTopic.
type CloseForumTopicParams = bot.CloseForumTopicParams

// ReopenForumTopicParams contains supported parameters for reopenForumTopic.
type ReopenForumTopicParams = bot.ReopenForumTopicParams

// DeleteForumTopicParams contains supported parameters for deleteForumTopic.
type DeleteForumTopicParams = bot.DeleteForumTopicParams

// UnpinAllForumTopicMessagesParams contains supported parameters for unpinAllForumTopicMessages.
type UnpinAllForumTopicMessagesParams = bot.UnpinAllForumTopicMessagesParams

// EditGeneralForumTopicParams contains supported parameters for editGeneralForumTopic.
type EditGeneralForumTopicParams = bot.EditGeneralForumTopicParams

// CloseGeneralForumTopicParams contains supported parameters for closeGeneralForumTopic.
type CloseGeneralForumTopicParams = bot.CloseGeneralForumTopicParams

// ReopenGeneralForumTopicParams contains supported parameters for reopenGeneralForumTopic.
type ReopenGeneralForumTopicParams = bot.ReopenGeneralForumTopicParams

// HideGeneralForumTopicParams contains supported parameters for hideGeneralForumTopic.
type HideGeneralForumTopicParams = bot.HideGeneralForumTopicParams

// UnhideGeneralForumTopicParams contains supported parameters for unhideGeneralForumTopic.
type UnhideGeneralForumTopicParams = bot.UnhideGeneralForumTopicParams

// UnpinAllGeneralForumTopicMessagesParams contains supported parameters for unpinAllGeneralForumTopicMessages.
type UnpinAllGeneralForumTopicMessagesParams = bot.UnpinAllGeneralForumTopicMessagesParams

// GetChatParams contains supported parameters for getChat.
type GetChatParams = bot.GetChatParams

// GetChatMemberParams contains supported parameters for getChatMember.
type GetChatMemberParams = bot.GetChatMemberParams

// GetChatAdministratorsParams contains supported parameters for getChatAdministrators.
type GetChatAdministratorsParams = bot.GetChatAdministratorsParams

// GetChatMemberCountParams contains supported parameters for getChatMemberCount.
type GetChatMemberCountParams = bot.GetChatMemberCountParams

// SetMyCommandsParams contains supported parameters for setMyCommands.
type SetMyCommandsParams = bot.SetMyCommandsParams

// DeleteMyCommandsParams contains supported parameters for deleteMyCommands.
type DeleteMyCommandsParams = bot.DeleteMyCommandsParams

// GetMyCommandsParams contains supported parameters for getMyCommands.
type GetMyCommandsParams = bot.GetMyCommandsParams

// SetChatMenuButtonParams contains supported parameters for setChatMenuButton.
type SetChatMenuButtonParams = bot.SetChatMenuButtonParams

// GetChatMenuButtonParams contains supported parameters for getChatMenuButton.
type GetChatMenuButtonParams = bot.GetChatMenuButtonParams

// SetMyDefaultAdministratorRightsParams contains supported parameters for setMyDefaultAdministratorRights.
type SetMyDefaultAdministratorRightsParams = bot.SetMyDefaultAdministratorRightsParams

// SetMyNameParams contains supported parameters for setMyName.
type SetMyNameParams = bot.SetMyNameParams

// GetMyNameParams contains supported parameters for getMyName.
type GetMyNameParams = bot.GetMyNameParams

// SetMyDescriptionParams contains supported parameters for setMyDescription.
type SetMyDescriptionParams = bot.SetMyDescriptionParams

// GetMyDescriptionParams contains supported parameters for getMyDescription.
type GetMyDescriptionParams = bot.GetMyDescriptionParams

// SetMyShortDescriptionParams contains supported parameters for setMyShortDescription.
type SetMyShortDescriptionParams = bot.SetMyShortDescriptionParams

// GetMyShortDescriptionParams contains supported parameters for getMyShortDescription.
type GetMyShortDescriptionParams = bot.GetMyShortDescriptionParams

// GetMyDefaultAdministratorRightsParams contains supported parameters for getMyDefaultAdministratorRights.
type GetMyDefaultAdministratorRightsParams = bot.GetMyDefaultAdministratorRightsParams

// SetMyProfilePhotoParams contains supported parameters for setMyProfilePhoto.
type SetMyProfilePhotoParams = bot.SetMyProfilePhotoParams

// RemoveMyProfilePhotoParams contains supported parameters for removeMyProfilePhoto.
type RemoveMyProfilePhotoParams = bot.RemoveMyProfilePhotoParams

// InputProfilePhoto describes a bot profile photo accepted by setMyProfilePhoto.
type InputProfilePhoto = bot.InputProfilePhoto

// InputProfilePhotoStatic describes a static JPG bot profile photo.
type InputProfilePhotoStatic = bot.InputProfilePhotoStatic

// InputProfilePhotoAnimated describes an animated MPEG4 bot profile photo.
type InputProfilePhotoAnimated = bot.InputProfilePhotoAnimated

// ExportChatInviteLinkParams contains supported parameters for exportChatInviteLink.
type ExportChatInviteLinkParams = bot.ExportChatInviteLinkParams

// CreateChatInviteLinkParams contains supported parameters for createChatInviteLink.
type CreateChatInviteLinkParams = bot.CreateChatInviteLinkParams

// EditChatInviteLinkParams contains supported parameters for editChatInviteLink.
type EditChatInviteLinkParams = bot.EditChatInviteLinkParams

// RevokeChatInviteLinkParams contains supported parameters for revokeChatInviteLink.
type RevokeChatInviteLinkParams = bot.RevokeChatInviteLinkParams

// ApproveChatJoinRequestParams contains supported parameters for approveChatJoinRequest.
type ApproveChatJoinRequestParams = bot.ApproveChatJoinRequestParams

// DeclineChatJoinRequestParams contains supported parameters for declineChatJoinRequest.
type DeclineChatJoinRequestParams = bot.DeclineChatJoinRequestParams

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

// InputMedia describes one media item accepted by sendMediaGroup.
type InputMedia = bot.InputMedia

// InputMediaPhoto describes a photo item for sendMediaGroup.
type InputMediaPhoto = bot.InputMediaPhoto

// InputMediaVideo describes a video item for sendMediaGroup.
type InputMediaVideo = bot.InputMediaVideo

// InputMediaAudio describes an audio item for sendMediaGroup.
type InputMediaAudio = bot.InputMediaAudio

// InputMediaDocument describes a document item for sendMediaGroup.
type InputMediaDocument = bot.InputMediaDocument

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

// Sticker represents a Telegram sticker.
type Sticker = telegram.Sticker

// StickerSet represents a Telegram sticker set.
type StickerSet = telegram.StickerSet

// MaskPosition describes where a mask sticker should be placed on faces.
type MaskPosition = telegram.MaskPosition

// File represents a Telegram file metadata object.
type File = telegram.File

// ReplyMarkup marks Telegram reply markup objects.
type ReplyMarkup = telegram.ReplyMarkup

// ReplyParameters describes the message being replied to.
type ReplyParameters = telegram.ReplyParameters

// ChatPermissions describes actions a user is allowed to take in a chat.
type ChatPermissions = telegram.ChatPermissions

// BotName describes a localized bot name.
type BotName = telegram.BotName

// BotDescription describes a localized bot description.
type BotDescription = telegram.BotDescription

// BotShortDescription describes a localized bot short description.
type BotShortDescription = telegram.BotShortDescription

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

// NewInputSticker creates an InputSticker with required fields.
func NewInputSticker(sticker FileRef, format string, emojiList ...string) InputSticker {
	return bot.NewInputSticker(sticker, format, emojiList...)
}

// MediaPhoto creates a photo input media item.
func MediaPhoto(media FileRef) InputMediaPhoto {
	return bot.MediaPhoto(media)
}

// MediaVideo creates a video input media item.
func MediaVideo(media FileRef) InputMediaVideo {
	return bot.MediaVideo(media)
}

// MediaAudio creates an audio input media item.
func MediaAudio(media FileRef) InputMediaAudio {
	return bot.MediaAudio(media)
}

// MediaDocument creates a document input media item.
func MediaDocument(media FileRef) InputMediaDocument {
	return bot.MediaDocument(media)
}

// ProfilePhotoStatic creates a static JPG input profile photo.
func ProfilePhotoStatic(photo FileRef) InputProfilePhotoStatic {
	return bot.ProfilePhotoStatic(photo)
}

// ProfilePhotoAnimated creates an animated MPEG4 input profile photo.
func ProfilePhotoAnimated(animation FileRef) InputProfilePhotoAnimated {
	return bot.ProfilePhotoAnimated(animation)
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
