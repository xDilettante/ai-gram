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

// SendLivePhotoParams contains supported parameters for sendLivePhoto.
type SendLivePhotoParams = bot.SendLivePhotoParams

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

// SendChecklistParams contains supported parameters for sendChecklist.
type SendChecklistParams = bot.SendChecklistParams

// EditMessageChecklistParams contains supported parameters for editMessageChecklist.
type EditMessageChecklistParams = bot.EditMessageChecklistParams

// SendMessageDraftParams contains supported parameters for sendMessageDraft.
type SendMessageDraftParams = bot.SendMessageDraftParams

// InputPollOption describes one poll option to send.
type InputPollOption = telegram.InputPollOption

// PollMedia describes media attached to a poll, quiz explanation, or poll option.
type PollMedia = telegram.PollMedia

// PollOption describes one answer option in a Telegram poll.
type PollOption = telegram.PollOption

// Poll describes a native Telegram poll.
type Poll = telegram.Poll

// PollAnswer represents an answer of a user or anonymous voter in a non-anonymous poll.
type PollAnswer = telegram.PollAnswer

// InputChecklist describes a checklist to create.
type InputChecklist = telegram.InputChecklist

// InputChecklistTask describes a checklist task to create.
type InputChecklistTask = telegram.InputChecklistTask

// StopPollParams contains supported parameters for stopPoll.
type StopPollParams = bot.StopPollParams

// SendDiceParams contains supported parameters for sendDice.
type SendDiceParams = bot.SendDiceParams

// SendGameParams contains supported parameters for sendGame.
type SendGameParams = bot.SendGameParams

// SetGameScoreParams contains supported parameters for setGameScore.
type SetGameScoreParams = bot.SetGameScoreParams

// GetGameHighScoresParams contains supported parameters for getGameHighScores.
type GetGameHighScoresParams = bot.GetGameHighScoresParams

// SetPassportDataErrorsParams contains supported parameters for setPassportDataErrors.
type SetPassportDataErrorsParams = bot.SetPassportDataErrorsParams

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

// GetUserProfilePhotosParams contains supported parameters for getUserProfilePhotos.
type GetUserProfilePhotosParams = bot.GetUserProfilePhotosParams

// GetUserProfileAudiosParams contains supported parameters for getUserProfileAudios.
type GetUserProfileAudiosParams = bot.GetUserProfileAudiosParams

// SetUserEmojiStatusParams contains supported parameters for setUserEmojiStatus.
type SetUserEmojiStatusParams = bot.SetUserEmojiStatusParams

// VerifyUserParams contains supported parameters for verifyUser.
type VerifyUserParams = bot.VerifyUserParams

// VerifyChatParams contains supported parameters for verifyChat.
type VerifyChatParams = bot.VerifyChatParams

// RemoveUserVerificationParams contains supported parameters for removeUserVerification.
type RemoveUserVerificationParams = bot.RemoveUserVerificationParams

// RemoveChatVerificationParams contains supported parameters for removeChatVerification.
type RemoveChatVerificationParams = bot.RemoveChatVerificationParams

// GetFileParams contains supported parameters for getFile.
type GetFileParams = bot.GetFileParams

// SetWebhookParams contains supported parameters for setWebhook.
type SetWebhookParams = bot.SetWebhookParams

// DeleteWebhookParams contains supported parameters for deleteWebhook.
type DeleteWebhookParams = bot.DeleteWebhookParams

// AnswerInlineQueryParams contains supported parameters for answerInlineQuery.
type AnswerInlineQueryParams = bot.AnswerInlineQueryParams

// AnswerGuestQueryParams contains supported parameters for answerGuestQuery.
type AnswerGuestQueryParams = bot.AnswerGuestQueryParams

// AnswerWebAppQueryParams contains supported parameters for answerWebAppQuery.
type AnswerWebAppQueryParams = bot.AnswerWebAppQueryParams

// GetBusinessConnectionParams contains supported parameters for getBusinessConnection.
type GetBusinessConnectionParams = bot.GetBusinessConnectionParams

// DeleteBusinessMessagesParams contains supported parameters for deleteBusinessMessages.
type DeleteBusinessMessagesParams = bot.DeleteBusinessMessagesParams

// ReadBusinessMessageParams contains supported parameters for readBusinessMessage.
type ReadBusinessMessageParams = bot.ReadBusinessMessageParams

// SetBusinessAccountNameParams contains supported parameters for setBusinessAccountName.
type SetBusinessAccountNameParams = bot.SetBusinessAccountNameParams

// SetBusinessAccountUsernameParams contains supported parameters for setBusinessAccountUsername.
type SetBusinessAccountUsernameParams = bot.SetBusinessAccountUsernameParams

// SetBusinessAccountBioParams contains supported parameters for setBusinessAccountBio.
type SetBusinessAccountBioParams = bot.SetBusinessAccountBioParams

// SetBusinessAccountProfilePhotoParams contains supported parameters for setBusinessAccountProfilePhoto.
type SetBusinessAccountProfilePhotoParams = bot.SetBusinessAccountProfilePhotoParams

// RemoveBusinessAccountProfilePhotoParams contains supported parameters for removeBusinessAccountProfilePhoto.
type RemoveBusinessAccountProfilePhotoParams = bot.RemoveBusinessAccountProfilePhotoParams

// SetBusinessAccountGiftSettingsParams contains supported parameters for setBusinessAccountGiftSettings.
type SetBusinessAccountGiftSettingsParams = bot.SetBusinessAccountGiftSettingsParams

// InputStoryContent marks story content accepted by postStory and editStory.
type InputStoryContent = bot.InputStoryContent

// InputStoryContentPhoto describes a photo story content item.
type InputStoryContentPhoto = bot.InputStoryContentPhoto

// InputStoryContentVideo describes a video story content item.
type InputStoryContentVideo = bot.InputStoryContentVideo

// PostStoryParams contains supported parameters for postStory.
type PostStoryParams = bot.PostStoryParams

// EditStoryParams contains supported parameters for editStory.
type EditStoryParams = bot.EditStoryParams

// DeleteStoryParams contains supported parameters for deleteStory.
type DeleteStoryParams = bot.DeleteStoryParams

// RepostStoryParams contains supported parameters for repostStory.
type RepostStoryParams = bot.RepostStoryParams

// ApproveSuggestedPostParams contains supported parameters for approveSuggestedPost.
type ApproveSuggestedPostParams = bot.ApproveSuggestedPostParams

// DeclineSuggestedPostParams contains supported parameters for declineSuggestedPost.
type DeclineSuggestedPostParams = bot.DeclineSuggestedPostParams

// SendInvoiceParams contains supported parameters for sendInvoice.
type SendInvoiceParams = bot.SendInvoiceParams

// CreateInvoiceLinkParams contains supported parameters for createInvoiceLink.
type CreateInvoiceLinkParams = bot.CreateInvoiceLinkParams

// AnswerShippingQueryParams contains supported parameters for answerShippingQuery.
type AnswerShippingQueryParams = bot.AnswerShippingQueryParams

// AnswerPreCheckoutQueryParams contains supported parameters for answerPreCheckoutQuery.
type AnswerPreCheckoutQueryParams = bot.AnswerPreCheckoutQueryParams

// InputMessageContent marks Telegram input message content objects used by inline results.
type InputMessageContent = bot.InputMessageContent

// InputTextMessageContent describes text content for an inline query result.
type InputTextMessageContent = bot.InputTextMessageContent

// InputLocationMessageContent describes location content for an inline query result.
type InputLocationMessageContent = bot.InputLocationMessageContent

// InputVenueMessageContent describes venue content for an inline query result.
type InputVenueMessageContent = bot.InputVenueMessageContent

// InputContactMessageContent describes contact content for an inline query result.
type InputContactMessageContent = bot.InputContactMessageContent

// InputInvoiceMessageContent describes invoice content for an inline query result.
type InputInvoiceMessageContent = bot.InputInvoiceMessageContent

// InlineQueryResult marks Telegram inline query result objects.
type InlineQueryResult = bot.InlineQueryResult

// InlineQueryResultArticle represents an article inline query result.
type InlineQueryResultArticle = bot.InlineQueryResultArticle

// InlineQueryResultLocation represents a location inline query result.
type InlineQueryResultLocation = bot.InlineQueryResultLocation

// InlineQueryResultVenue represents a venue inline query result.
type InlineQueryResultVenue = bot.InlineQueryResultVenue

// InlineQueryResultContact represents a contact inline query result.
type InlineQueryResultContact = bot.InlineQueryResultContact

// InlineQueryResultGame represents a game inline query result.
type InlineQueryResultGame = bot.InlineQueryResultGame

// InlineQueryResultPhoto represents a photo inline query result.
type InlineQueryResultPhoto = bot.InlineQueryResultPhoto

// InlineQueryResultGif represents a GIF inline query result.
type InlineQueryResultGif = bot.InlineQueryResultGif

// InlineQueryResultMpeg4Gif represents an MPEG-4 GIF inline query result.
type InlineQueryResultMpeg4Gif = bot.InlineQueryResultMpeg4Gif

// InlineQueryResultVideo represents a video inline query result.
type InlineQueryResultVideo = bot.InlineQueryResultVideo

// InlineQueryResultAudio represents an audio inline query result.
type InlineQueryResultAudio = bot.InlineQueryResultAudio

// InlineQueryResultVoice represents a voice inline query result.
type InlineQueryResultVoice = bot.InlineQueryResultVoice

// InlineQueryResultDocument represents a document inline query result.
type InlineQueryResultDocument = bot.InlineQueryResultDocument

// InlineQueryResultCachedPhoto represents a cached photo inline query result.
type InlineQueryResultCachedPhoto = bot.InlineQueryResultCachedPhoto

// InlineQueryResultCachedGif represents a cached GIF inline query result.
type InlineQueryResultCachedGif = bot.InlineQueryResultCachedGif

// InlineQueryResultCachedMpeg4Gif represents a cached MPEG-4 GIF inline query result.
type InlineQueryResultCachedMpeg4Gif = bot.InlineQueryResultCachedMpeg4Gif

// InlineQueryResultCachedSticker represents a cached sticker inline query result.
type InlineQueryResultCachedSticker = bot.InlineQueryResultCachedSticker

// InlineQueryResultCachedDocument represents a cached document inline query result.
type InlineQueryResultCachedDocument = bot.InlineQueryResultCachedDocument

// InlineQueryResultCachedVideo represents a cached video inline query result.
type InlineQueryResultCachedVideo = bot.InlineQueryResultCachedVideo

// InlineQueryResultCachedVoice represents a cached voice inline query result.
type InlineQueryResultCachedVoice = bot.InlineQueryResultCachedVoice

// InlineQueryResultCachedAudio represents a cached audio inline query result.
type InlineQueryResultCachedAudio = bot.InlineQueryResultCachedAudio

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

// DeleteMessageReactionParams contains supported parameters for deleteMessageReaction.
type DeleteMessageReactionParams = bot.DeleteMessageReactionParams

// DeleteAllMessageReactionsParams contains supported parameters for deleteAllMessageReactions.
type DeleteAllMessageReactionsParams = bot.DeleteAllMessageReactionsParams

// GetUserChatBoostsParams contains supported parameters for getUserChatBoosts.
type GetUserChatBoostsParams = bot.GetUserChatBoostsParams

// SetChatMemberTagParams contains supported parameters for setChatMemberTag.
type SetChatMemberTagParams = bot.SetChatMemberTagParams

// ReactionType marks Telegram message reaction type objects.
type ReactionType = telegram.ReactionType

// ReactionTypeEmoji describes an emoji-based reaction.
type ReactionTypeEmoji = telegram.ReactionTypeEmoji

// ReactionTypeCustomEmoji describes a custom emoji reaction.
type ReactionTypeCustomEmoji = telegram.ReactionTypeCustomEmoji

// ReactionTypePaid describes a paid reaction.
type ReactionTypePaid = telegram.ReactionTypePaid

// BanChatMemberParams contains supported parameters for banChatMember.
type BanChatMemberParams = bot.BanChatMemberParams

// BanChatSenderChatParams contains supported parameters for banChatSenderChat.
type BanChatSenderChatParams = bot.BanChatSenderChatParams

// UnbanChatMemberParams contains supported parameters for unbanChatMember.
type UnbanChatMemberParams = bot.UnbanChatMemberParams

// UnbanChatSenderChatParams contains supported parameters for unbanChatSenderChat.
type UnbanChatSenderChatParams = bot.UnbanChatSenderChatParams

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

// CreateChatSubscriptionInviteLinkParams contains supported parameters for createChatSubscriptionInviteLink.
type CreateChatSubscriptionInviteLinkParams = bot.CreateChatSubscriptionInviteLinkParams

// EditChatInviteLinkParams contains supported parameters for editChatInviteLink.
type EditChatInviteLinkParams = bot.EditChatInviteLinkParams

// EditChatSubscriptionInviteLinkParams contains supported parameters for editChatSubscriptionInviteLink.
type EditChatSubscriptionInviteLinkParams = bot.EditChatSubscriptionInviteLinkParams

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

// EditMessageMediaParams contains supported parameters for editMessageMedia.
type EditMessageMediaParams = bot.EditMessageMediaParams

// EditMessageLiveLocationParams contains supported parameters for editMessageLiveLocation.
type EditMessageLiveLocationParams = bot.EditMessageLiveLocationParams

// StopMessageLiveLocationParams contains supported parameters for stopMessageLiveLocation.
type StopMessageLiveLocationParams = bot.StopMessageLiveLocationParams

// ChatID identifies a Telegram chat by numeric ID or username string.
type ChatID = bot.ChatID

// FileRef identifies media by Telegram file_id, HTTP(S) URL, or multipart upload.
type FileRef = bot.FileRef

// UploadFile describes a file uploaded through multipart/form-data.
type UploadFile = bot.UploadFile

// InputMedia describes one media item accepted by supported media methods.
type InputMedia = bot.InputMedia

// InputMediaPhoto describes a photo input media item.
type InputMediaPhoto = bot.InputMediaPhoto

// InputMediaVideo describes a video input media item.
type InputMediaVideo = bot.InputMediaVideo

// InputMediaAnimation describes an animation item for editMessageMedia.
type InputMediaAnimation = bot.InputMediaAnimation

// InputMediaAudio describes an audio input media item.
type InputMediaAudio = bot.InputMediaAudio

// InputMediaDocument describes a document input media item.
type InputMediaDocument = bot.InputMediaDocument

// InputPollMedia describes media accepted in sendPoll media fields.
type InputPollMedia = bot.InputPollMedia

// InputPollOptionMedia describes media accepted in InputPollOption.media.
type InputPollOptionMedia = bot.InputPollOptionMedia

// InputMediaLivePhoto describes a live photo input media item.
type InputMediaLivePhoto = bot.InputMediaLivePhoto

// InputMediaLocation describes a location input media item.
type InputMediaLocation = bot.InputMediaLocation

// InputMediaSticker describes a sticker input media item.
type InputMediaSticker = bot.InputMediaSticker

// InputMediaVenue describes a venue input media item.
type InputMediaVenue = bot.InputMediaVenue

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

// ChatFullInfo contains full information about a chat returned by getChat.
type ChatFullInfo = telegram.ChatFullInfo

// ChatPhoto represents a Telegram chat photo.
type ChatPhoto = telegram.ChatPhoto

// Birthdate describes the birthdate of a user.
type Birthdate = telegram.Birthdate

// BusinessIntro contains start page settings of a Telegram Business account.
type BusinessIntro = telegram.BusinessIntro

// BusinessLocation contains the location of a Telegram Business account.
type BusinessLocation = telegram.BusinessLocation

// BusinessOpeningHours describes the opening hours of a business.
type BusinessOpeningHours = telegram.BusinessOpeningHours

// BusinessOpeningHoursInterval describes one business opening interval.
type BusinessOpeningHoursInterval = telegram.BusinessOpeningHoursInterval

// ChatLocation represents a location to which a chat is connected.
type ChatLocation = telegram.ChatLocation

// UserRating describes a user's Telegram Stars purchase reliability rating.
type UserRating = telegram.UserRating

// ChatMemberStatus identifies a user's membership state in a chat.
type ChatMemberStatus = telegram.ChatMemberStatus

// ChatMember is one of the official Telegram chat member variants.
type ChatMember = telegram.ChatMember

// ChatMemberOwner describes the chat owner.
type ChatMemberOwner = telegram.ChatMemberOwner

// ChatMemberAdministrator describes an administrator in a chat.
type ChatMemberAdministrator = telegram.ChatMemberAdministrator

// ChatMemberMember describes a regular chat member.
type ChatMemberMember = telegram.ChatMemberMember

// ChatMemberRestricted describes a restricted chat member.
type ChatMemberRestricted = telegram.ChatMemberRestricted

// ChatMemberLeft describes a user who is not currently a chat member.
type ChatMemberLeft = telegram.ChatMemberLeft

// ChatMemberBanned describes a user banned from a chat.
type ChatMemberBanned = telegram.ChatMemberBanned

// ChatMemberUpdated represents changes in the status of a chat member.
type ChatMemberUpdated = telegram.ChatMemberUpdated

// ChatBoostUpdated represents a boost added to a chat or changed.
type ChatBoostUpdated = telegram.ChatBoostUpdated

// ChatBoostRemoved represents a boost removed from a chat.
type ChatBoostRemoved = telegram.ChatBoostRemoved

// ChatBoost contains information about a chat boost.
type ChatBoost = telegram.ChatBoost

// ChatBoostSource marks Telegram chat boost source objects.
type ChatBoostSource = telegram.ChatBoostSource

// ChatBoostSourcePremium describes a boost from Telegram Premium.
type ChatBoostSourcePremium = telegram.ChatBoostSourcePremium

// ChatBoostSourceGiftCode describes a boost from Premium gift codes.
type ChatBoostSourceGiftCode = telegram.ChatBoostSourceGiftCode

// ChatBoostSourceGiveaway describes a boost from a Premium or Star giveaway.
type ChatBoostSourceGiveaway = telegram.ChatBoostSourceGiveaway

// UserChatBoosts represents a list of boosts added to a chat by a user.
type UserChatBoosts = telegram.UserChatBoosts

// CallbackQuery represents an incoming callback query.
type CallbackQuery = telegram.CallbackQuery

// MessageOrigin marks Telegram message origin objects.
type MessageOrigin = telegram.MessageOrigin

// MessageOriginUser describes a message originally sent by a known user.
type MessageOriginUser = telegram.MessageOriginUser

// MessageOriginHiddenUser describes a message originally sent by an unknown user.
type MessageOriginHiddenUser = telegram.MessageOriginHiddenUser

// MessageOriginChat describes a message originally sent on behalf of a chat.
type MessageOriginChat = telegram.MessageOriginChat

// MessageOriginChannel describes a message originally sent to a channel chat.
type MessageOriginChannel = telegram.MessageOriginChannel

// ExternalReplyInfo describes a message being replied to from another chat or topic.
type ExternalReplyInfo = telegram.ExternalReplyInfo

// TextQuote contains information about the quoted part of a replied-to message.
type TextQuote = telegram.TextQuote

// InaccessibleMessage describes a message that is inaccessible to the bot.
type InaccessibleMessage = telegram.InaccessibleMessage

// MaybeInaccessibleMessage describes a message that may be inaccessible to the bot.
type MaybeInaccessibleMessage = telegram.MaybeInaccessibleMessage

// DirectMessagesTopic describes a topic of a channel direct messages chat.
type DirectMessagesTopic = telegram.DirectMessagesTopic

// UsersShared contains users shared through a reply keyboard button.
type UsersShared = telegram.UsersShared

// SharedUser contains information about one user shared with the bot.
type SharedUser = telegram.SharedUser

// ChatShared contains a chat shared through a reply keyboard button.
type ChatShared = telegram.ChatShared

// ProximityAlertTriggered describes a proximity alert service message.
type ProximityAlertTriggered = telegram.ProximityAlertTriggered

// MessageAutoDeleteTimerChanged describes changed chat auto-delete settings.
type MessageAutoDeleteTimerChanged = telegram.MessageAutoDeleteTimerChanged

// ChatBoostAdded describes a chat boost service message.
type ChatBoostAdded = telegram.ChatBoostAdded

// ChatBackground represents a chat background service message payload.
type ChatBackground = telegram.ChatBackground

// BackgroundType marks chat background type objects.
type BackgroundType = telegram.BackgroundType

// BackgroundTypeFill describes a background filled by a gradient or solid color.
type BackgroundTypeFill = telegram.BackgroundTypeFill

// BackgroundTypeWallpaper describes a wallpaper background.
type BackgroundTypeWallpaper = telegram.BackgroundTypeWallpaper

// BackgroundTypePattern describes a patterned background.
type BackgroundTypePattern = telegram.BackgroundTypePattern

// BackgroundTypeChatTheme describes a built-in chat theme background.
type BackgroundTypeChatTheme = telegram.BackgroundTypeChatTheme

// BackgroundFill marks chat background fill objects.
type BackgroundFill = telegram.BackgroundFill

// BackgroundFillSolid describes a solid color background fill.
type BackgroundFillSolid = telegram.BackgroundFillSolid

// BackgroundFillGradient describes a two-color gradient background fill.
type BackgroundFillGradient = telegram.BackgroundFillGradient

// BackgroundFillFreeformGradient describes a freeform gradient background fill.
type BackgroundFillFreeformGradient = telegram.BackgroundFillFreeformGradient

// ChatOwnerLeft describes a service message about the chat owner leaving.
type ChatOwnerLeft = telegram.ChatOwnerLeft

// ChatOwnerChanged describes a service message about ownership transfer.
type ChatOwnerChanged = telegram.ChatOwnerChanged

// SuggestedPostInfo contains metadata about a suggested post message.
type SuggestedPostInfo = telegram.SuggestedPostInfo

// GiveawayCreated describes a scheduled giveaway creation service message.
type GiveawayCreated = telegram.GiveawayCreated

// Giveaway represents a scheduled giveaway message.
type Giveaway = telegram.Giveaway

// GiveawayWinners represents a completed giveaway with public winners.
type GiveawayWinners = telegram.GiveawayWinners

// GiveawayCompleted describes a completed giveaway service message.
type GiveawayCompleted = telegram.GiveawayCompleted

// PaidMessagePriceChanged describes a paid-message price change service message.
type PaidMessagePriceChanged = telegram.PaidMessagePriceChanged

// DirectMessagePriceChanged describes a direct-message price change service message.
type DirectMessagePriceChanged = telegram.DirectMessagePriceChanged

// VideoChatScheduled describes a scheduled video chat service message.
type VideoChatScheduled = telegram.VideoChatScheduled

// VideoChatStarted describes a started video chat service message.
type VideoChatStarted = telegram.VideoChatStarted

// VideoChatEnded describes an ended video chat service message.
type VideoChatEnded = telegram.VideoChatEnded

// VideoChatParticipantsInvited describes users invited to a video chat.
type VideoChatParticipantsInvited = telegram.VideoChatParticipantsInvited

// InlineQuery represents an incoming inline query.
type InlineQuery = telegram.InlineQuery

// ChosenInlineResult represents a chosen inline query result.
type ChosenInlineResult = telegram.ChosenInlineResult

// LinkPreviewOptions describes link preview generation options.
type LinkPreviewOptions = telegram.LinkPreviewOptions

// InlineQueryResultsButton represents a button shown above inline query results.
type InlineQueryResultsButton = telegram.InlineQueryResultsButton

// WebAppInfo describes a Telegram Web App URL.
type WebAppInfo = telegram.WebAppInfo

// WebAppData describes data sent from a Web App to the bot.
type WebAppData = telegram.WebAppData

// WriteAccessAllowed describes a service message about Web App write access.
type WriteAccessAllowed = telegram.WriteAccessAllowed

// SentWebAppMessage describes an inline message sent by a Web App on behalf of a user.
type SentWebAppMessage = telegram.SentWebAppMessage

// SentGuestMessage describes an inline message sent by a guest bot.
type SentGuestMessage = telegram.SentGuestMessage

// BusinessBotRights represents the rights of a business bot.
type BusinessBotRights = telegram.BusinessBotRights

// BusinessConnection describes a bot connection with a business account.
type BusinessConnection = telegram.BusinessConnection

// BusinessMessagesDeleted describes deleted messages from a connected business account.
type BusinessMessagesDeleted = telegram.BusinessMessagesDeleted

// AcceptedGiftTypes describes gift types accepted by a business account, user, or chat.
type AcceptedGiftTypes = telegram.AcceptedGiftTypes

// Story represents a Telegram story.
type Story = telegram.Story

// StoryArea describes a clickable area on a story media.
type StoryArea = telegram.StoryArea

// StoryAreaPosition describes the position of a clickable story area.
type StoryAreaPosition = telegram.StoryAreaPosition

// StoryAreaType marks Telegram story area type objects.
type StoryAreaType = telegram.StoryAreaType

// LocationAddress describes the physical address of a location.
type LocationAddress = telegram.LocationAddress

// StoryAreaTypeLocation describes a story area pointing to a location.
type StoryAreaTypeLocation = telegram.StoryAreaTypeLocation

// StoryAreaTypeSuggestedReaction describes a story area pointing to a suggested reaction.
type StoryAreaTypeSuggestedReaction = telegram.StoryAreaTypeSuggestedReaction

// StoryAreaTypeLink describes a story area pointing to a link.
type StoryAreaTypeLink = telegram.StoryAreaTypeLink

// StoryAreaTypeWeather describes a story area containing weather information.
type StoryAreaTypeWeather = telegram.StoryAreaTypeWeather

// StoryAreaTypeUniqueGift describes a story area pointing to a unique gift.
type StoryAreaTypeUniqueGift = telegram.StoryAreaTypeUniqueGift

// SuggestedPostApprovalFailed describes a failed suggested post approval service message.
type SuggestedPostApprovalFailed = telegram.SuggestedPostApprovalFailed

// SuggestedPostApproved describes a suggested post approval service message.
type SuggestedPostApproved = telegram.SuggestedPostApproved

// SuggestedPostDeclined describes a suggested post decline service message.
type SuggestedPostDeclined = telegram.SuggestedPostDeclined

// SuggestedPostPaid describes a suggested post payment service message.
type SuggestedPostPaid = telegram.SuggestedPostPaid

// SuggestedPostRefunded describes a suggested post refund service message.
type SuggestedPostRefunded = telegram.SuggestedPostRefunded

// WebhookInfo describes current Telegram webhook status.
type WebhookInfo = telegram.WebhookInfo

// Sticker represents a Telegram sticker.
type Sticker = telegram.Sticker

// LivePhoto represents a Telegram live photo.
type LivePhoto = telegram.LivePhoto

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

// ReplyChatID marks chat identifiers accepted by ReplyParameters.ChatID.
type ReplyChatID = telegram.ReplyChatID

// ReplyChatIDInt identifies a chat by its numeric identifier.
type ReplyChatIDInt = telegram.ReplyChatIDInt

// ReplyChatIDUsername identifies a channel by its @username.
type ReplyChatIDUsername = telegram.ReplyChatIDUsername

// ChatPermissions describes actions a user is allowed to take in a chat.
type ChatPermissions = telegram.ChatPermissions

// BotName describes a localized bot name.
type BotName = telegram.BotName

// BotDescription describes a localized bot description.
type BotDescription = telegram.BotDescription

// BotShortDescription describes a localized bot short description.
type BotShortDescription = telegram.BotShortDescription

// LabeledPrice represents a price component in the smallest units of a currency.
type LabeledPrice = telegram.LabeledPrice

// PaidMediaInfo describes paid media attached to a message.
type PaidMediaInfo = telegram.PaidMediaInfo

// PaidMedia marks Telegram paid media objects.
type PaidMedia = telegram.PaidMedia

// StarTransactions contains a list of Telegram Star transactions.
type StarTransactions = telegram.StarTransactions

// StarTransaction describes a Telegram Star transaction.
type StarTransaction = telegram.StarTransaction

// StarAmount describes an amount of Telegram Stars.
type StarAmount = telegram.StarAmount

// GiftBackground describes the background of a regular gift.
type GiftBackground = telegram.GiftBackground

// Gift represents a gift that can be sent by the bot.
type Gift = telegram.Gift

// Gifts contains a list of gifts available to send.
type Gifts = telegram.Gifts

// UniqueGiftModel describes the model of a unique gift.
type UniqueGiftModel = telegram.UniqueGiftModel

// UniqueGiftSymbol describes the symbol shown on a unique gift.
type UniqueGiftSymbol = telegram.UniqueGiftSymbol

// UniqueGiftBackdropColors describes backdrop colors of a unique gift.
type UniqueGiftBackdropColors = telegram.UniqueGiftBackdropColors

// UniqueGiftBackdrop describes the backdrop of a unique gift.
type UniqueGiftBackdrop = telegram.UniqueGiftBackdrop

// UniqueGiftColors contains color scheme information derived from a unique gift.
type UniqueGiftColors = telegram.UniqueGiftColors

// UniqueGift describes a unique gift upgraded from a regular gift.
type UniqueGift = telegram.UniqueGift

// GiftInfo describes a service message about a regular gift.
type GiftInfo = telegram.GiftInfo

// UniqueGiftInfo describes a service message about a unique gift.
type UniqueGiftInfo = telegram.UniqueGiftInfo

// OwnedGift marks Telegram owned gift objects.
type OwnedGift = telegram.OwnedGift

// OwnedGiftRegular describes a regular owned gift.
type OwnedGiftRegular = telegram.OwnedGiftRegular

// OwnedGiftUnique describes a unique owned gift.
type OwnedGiftUnique = telegram.OwnedGiftUnique

// OwnedGifts contains gifts owned by a user or chat.
type OwnedGifts = telegram.OwnedGifts

// InputPaidMedia marks paid media items accepted by sendPaidMedia.
type InputPaidMedia = bot.InputPaidMedia

// InputPaidMediaPhoto describes a paid photo to send.
type InputPaidMediaPhoto = bot.InputPaidMediaPhoto

// InputPaidMediaLivePhoto describes a paid live photo to send.
type InputPaidMediaLivePhoto = bot.InputPaidMediaLivePhoto

// InputPaidMediaVideo describes a paid video to send.
type InputPaidMediaVideo = bot.InputPaidMediaVideo

// SendPaidMediaParams contains supported parameters for sendPaidMedia.
type SendPaidMediaParams = bot.SendPaidMediaParams

// GetStarTransactionsParams contains supported parameters for getStarTransactions.
type GetStarTransactionsParams = bot.GetStarTransactionsParams

// RefundStarPaymentParams contains supported parameters for refundStarPayment.
type RefundStarPaymentParams = bot.RefundStarPaymentParams

// GetAvailableGiftsParams contains supported parameters for getAvailableGifts.
type GetAvailableGiftsParams = bot.GetAvailableGiftsParams

// SendGiftParams contains supported parameters for sendGift.
type SendGiftParams = bot.SendGiftParams

// GiftPremiumSubscriptionParams contains supported parameters for giftPremiumSubscription.
type GiftPremiumSubscriptionParams = bot.GiftPremiumSubscriptionParams

// GetBusinessAccountStarBalanceParams contains supported parameters for getBusinessAccountStarBalance.
type GetBusinessAccountStarBalanceParams = bot.GetBusinessAccountStarBalanceParams

// TransferBusinessAccountStarsParams contains supported parameters for transferBusinessAccountStars.
type TransferBusinessAccountStarsParams = bot.TransferBusinessAccountStarsParams

// GetBusinessAccountGiftsParams contains supported parameters for getBusinessAccountGifts.
type GetBusinessAccountGiftsParams = bot.GetBusinessAccountGiftsParams

// GetUserGiftsParams contains supported parameters for getUserGifts.
type GetUserGiftsParams = bot.GetUserGiftsParams

// GetChatGiftsParams contains supported parameters for getChatGifts.
type GetChatGiftsParams = bot.GetChatGiftsParams

// ConvertGiftToStarsParams contains supported parameters for convertGiftToStars.
type ConvertGiftToStarsParams = bot.ConvertGiftToStarsParams

// UpgradeGiftParams contains supported parameters for upgradeGift.
type UpgradeGiftParams = bot.UpgradeGiftParams

// TransferGiftParams contains supported parameters for transferGift.
type TransferGiftParams = bot.TransferGiftParams

// GetMyStarBalanceParams contains supported parameters for getMyStarBalance.
type GetMyStarBalanceParams = bot.GetMyStarBalanceParams

// EditUserStarSubscriptionParams contains supported parameters for editUserStarSubscription.
type EditUserStarSubscriptionParams = bot.EditUserStarSubscriptionParams

// SavePreparedKeyboardButtonParams contains supported parameters for savePreparedKeyboardButton.
type SavePreparedKeyboardButtonParams = bot.SavePreparedKeyboardButtonParams

// SavePreparedInlineMessageParams contains supported parameters for savePreparedInlineMessage.
type SavePreparedInlineMessageParams = bot.SavePreparedInlineMessageParams

// GetManagedBotTokenParams contains supported parameters for getManagedBotToken.
type GetManagedBotTokenParams = bot.GetManagedBotTokenParams

// ReplaceManagedBotTokenParams contains supported parameters for replaceManagedBotToken.
type ReplaceManagedBotTokenParams = bot.ReplaceManagedBotTokenParams

// GetManagedBotAccessSettingsParams contains supported parameters for getManagedBotAccessSettings.
type GetManagedBotAccessSettingsParams = bot.GetManagedBotAccessSettingsParams

// SetManagedBotAccessSettingsParams contains supported parameters for setManagedBotAccessSettings.
type SetManagedBotAccessSettingsParams = bot.SetManagedBotAccessSettingsParams

// GetUserPersonalChatMessagesParams contains supported parameters for getUserPersonalChatMessages.
type GetUserPersonalChatMessagesParams = bot.GetUserPersonalChatMessagesParams

// InlineKeyboardMarkup represents an inline keyboard attached to a message.
type InlineKeyboardMarkup = telegram.InlineKeyboardMarkup

// InlineKeyboardButton represents one inline keyboard button.
type InlineKeyboardButton = telegram.InlineKeyboardButton

// LoginUrl represents an automatic Telegram login URL for an inline keyboard button.
type LoginUrl = telegram.LoginUrl

// SwitchInlineQueryChosenChat describes chat filters for switching to inline mode.
type SwitchInlineQueryChosenChat = telegram.SwitchInlineQueryChosenChat

// CopyTextButton describes text copied to the clipboard by an inline keyboard button.
type CopyTextButton = telegram.CopyTextButton

// ReplyKeyboardMarkup represents a custom reply keyboard.
type ReplyKeyboardMarkup = telegram.ReplyKeyboardMarkup

// KeyboardButton represents one custom reply keyboard button.
type KeyboardButton = telegram.KeyboardButton

// KeyboardButtonPollType represents a poll type requested by a reply keyboard button.
type KeyboardButtonPollType = telegram.KeyboardButtonPollType

// KeyboardButtonRequestUsers defines criteria for requesting users with a reply keyboard button.
type KeyboardButtonRequestUsers = telegram.KeyboardButtonRequestUsers

// KeyboardButtonRequestChat defines criteria for requesting a chat with a reply keyboard button.
type KeyboardButtonRequestChat = telegram.KeyboardButtonRequestChat

// KeyboardButtonRequestManagedBot defines parameters for creating and sharing a managed bot.
type KeyboardButtonRequestManagedBot = telegram.KeyboardButtonRequestManagedBot

// PreparedKeyboardButton describes a saved keyboard button for Mini App use.
type PreparedKeyboardButton = telegram.PreparedKeyboardButton

// BotAccessSettings describes access settings of a bot.
type BotAccessSettings = telegram.BotAccessSettings

// PreparedInlineMessage describes an inline message saved for a Mini App user.
type PreparedInlineMessage = telegram.PreparedInlineMessage

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

// MediaAnimation creates an animation input media item for editMessageMedia.
func MediaAnimation(media FileRef) InputMediaAnimation {
	return bot.MediaAnimation(media)
}

// MediaAudio creates an audio input media item.
func MediaAudio(media FileRef) InputMediaAudio {
	return bot.MediaAudio(media)
}

// MediaDocument creates a document input media item.
func MediaDocument(media FileRef) InputMediaDocument {
	return bot.MediaDocument(media)
}

// MediaLivePhoto creates a live photo input media item.
func MediaLivePhoto(media FileRef, photo FileRef) InputMediaLivePhoto {
	return bot.MediaLivePhoto(media, photo)
}

// MediaLocation creates a location input media item.
func MediaLocation(latitude float64, longitude float64) InputMediaLocation {
	return bot.MediaLocation(latitude, longitude)
}

// MediaSticker creates a sticker input media item.
func MediaSticker(media FileRef) InputMediaSticker {
	return bot.MediaSticker(media)
}

// MediaVenue creates a venue input media item.
func MediaVenue(latitude float64, longitude float64, title string, address string) InputMediaVenue {
	return bot.MediaVenue(latitude, longitude, title, address)
}

// PaidPhoto creates a paid photo input media item.
func PaidPhoto(media FileRef) InputPaidMediaPhoto {
	return bot.PaidPhoto(media)
}

// PaidLivePhoto creates a paid live photo input media item.
func PaidLivePhoto(media FileRef, photo FileRef) InputPaidMediaLivePhoto {
	return bot.PaidLivePhoto(media, photo)
}

// PaidVideo creates a paid video input media item.
func PaidVideo(media FileRef) InputPaidMediaVideo {
	return bot.PaidVideo(media)
}

// ProfilePhotoStatic creates a static JPG input profile photo.
func ProfilePhotoStatic(photo FileRef) InputProfilePhotoStatic {
	return bot.ProfilePhotoStatic(photo)
}

// ProfilePhotoAnimated creates an animated MPEG4 input profile photo.
func ProfilePhotoAnimated(animation FileRef) InputProfilePhotoAnimated {
	return bot.ProfilePhotoAnimated(animation)
}

// StoryPhoto creates photo content for a business story.
func StoryPhoto(photo FileRef) InputStoryContentPhoto {
	return bot.StoryPhoto(photo)
}

// StoryVideo creates video content for a business story.
func StoryVideo(video FileRef) InputStoryContentVideo {
	return bot.StoryVideo(video)
}

// NewStoryAreaTypeLocation creates a location story area type.
func NewStoryAreaTypeLocation(latitude float64, longitude float64) StoryAreaTypeLocation {
	return telegram.NewStoryAreaTypeLocation(latitude, longitude)
}

// NewStoryAreaTypeSuggestedReaction creates a suggested reaction story area type.
func NewStoryAreaTypeSuggestedReaction(reaction ReactionType) StoryAreaTypeSuggestedReaction {
	return telegram.NewStoryAreaTypeSuggestedReaction(reaction)
}

// NewStoryAreaTypeLink creates a link story area type.
func NewStoryAreaTypeLink(url string) StoryAreaTypeLink {
	return telegram.NewStoryAreaTypeLink(url)
}

// NewStoryAreaTypeWeather creates a weather story area type.
func NewStoryAreaTypeWeather(temperature float64, emoji string, backgroundColor int) StoryAreaTypeWeather {
	return telegram.NewStoryAreaTypeWeather(temperature, emoji, backgroundColor)
}

// NewStoryAreaTypeUniqueGift creates a unique gift story area type.
func NewStoryAreaTypeUniqueGift(name string) StoryAreaTypeUniqueGift {
	return telegram.NewStoryAreaTypeUniqueGift(name)
}

// InputText creates text content for an inline query result.
func InputText(message string) InputTextMessageContent {
	return bot.InputText(message)
}

// InputLocation creates location content for an inline query result.
func InputLocation(latitude float64, longitude float64) InputLocationMessageContent {
	return bot.InputLocation(latitude, longitude)
}

// InputVenue creates venue content for an inline query result.
func InputVenue(latitude float64, longitude float64, title string, address string) InputVenueMessageContent {
	return bot.InputVenue(latitude, longitude, title, address)
}

// InputContact creates contact content for an inline query result.
func InputContact(phoneNumber string, firstName string) InputContactMessageContent {
	return bot.InputContact(phoneNumber, firstName)
}

// InlineArticle creates an article inline query result.
func InlineArticle(id string, title string, content InputMessageContent) InlineQueryResultArticle {
	return bot.InlineArticle(id, title, content)
}

// InlineLocation creates a location inline query result.
func InlineLocation(id string, latitude float64, longitude float64, title string) InlineQueryResultLocation {
	return bot.InlineLocation(id, latitude, longitude, title)
}

// InlineVenue creates a venue inline query result.
func InlineVenue(id string, latitude float64, longitude float64, title string, address string) InlineQueryResultVenue {
	return bot.InlineVenue(id, latitude, longitude, title, address)
}

// InlineContact creates a contact inline query result.
func InlineContact(id string, phoneNumber string, firstName string) InlineQueryResultContact {
	return bot.InlineContact(id, phoneNumber, firstName)
}

// InlineGame creates a game inline query result.
func InlineGame(id string, gameShortName string) InlineQueryResultGame {
	return bot.InlineGame(id, gameShortName)
}

// InlinePhoto creates a photo inline query result.
func InlinePhoto(id string, photoURL string, thumbnailURL string) InlineQueryResultPhoto {
	return bot.InlinePhoto(id, photoURL, thumbnailURL)
}

// InlineGif creates a GIF inline query result.
func InlineGif(id string, gifURL string, thumbnailURL string) InlineQueryResultGif {
	return bot.InlineGif(id, gifURL, thumbnailURL)
}

// InlineMpeg4Gif creates an MPEG-4 GIF inline query result.
func InlineMpeg4Gif(id string, mpeg4URL string, thumbnailURL string) InlineQueryResultMpeg4Gif {
	return bot.InlineMpeg4Gif(id, mpeg4URL, thumbnailURL)
}

// InlineVideo creates a video inline query result.
func InlineVideo(id string, videoURL string, mimeType string, thumbnailURL string, title string) InlineQueryResultVideo {
	return bot.InlineVideo(id, videoURL, mimeType, thumbnailURL, title)
}

// InlineAudio creates an audio inline query result.
func InlineAudio(id string, audioURL string, title string) InlineQueryResultAudio {
	return bot.InlineAudio(id, audioURL, title)
}

// InlineVoice creates a voice inline query result.
func InlineVoice(id string, voiceURL string, title string) InlineQueryResultVoice {
	return bot.InlineVoice(id, voiceURL, title)
}

// InlineDocument creates a document inline query result.
func InlineDocument(id string, title string, documentURL string, mimeType string) InlineQueryResultDocument {
	return bot.InlineDocument(id, title, documentURL, mimeType)
}

// InlineCachedPhoto creates a cached photo inline query result.
func InlineCachedPhoto(id string, fileID string) InlineQueryResultCachedPhoto {
	return bot.InlineCachedPhoto(id, fileID)
}

// InlineCachedGif creates a cached GIF inline query result.
func InlineCachedGif(id string, fileID string) InlineQueryResultCachedGif {
	return bot.InlineCachedGif(id, fileID)
}

// InlineCachedMpeg4Gif creates a cached MPEG-4 GIF inline query result.
func InlineCachedMpeg4Gif(id string, fileID string) InlineQueryResultCachedMpeg4Gif {
	return bot.InlineCachedMpeg4Gif(id, fileID)
}

// InlineCachedSticker creates a cached sticker inline query result.
func InlineCachedSticker(id string, fileID string) InlineQueryResultCachedSticker {
	return bot.InlineCachedSticker(id, fileID)
}

// InlineCachedDocument creates a cached document inline query result.
func InlineCachedDocument(id string, fileID string, title string) InlineQueryResultCachedDocument {
	return bot.InlineCachedDocument(id, fileID, title)
}

// InlineCachedVideo creates a cached video inline query result.
func InlineCachedVideo(id string, fileID string, title string) InlineQueryResultCachedVideo {
	return bot.InlineCachedVideo(id, fileID, title)
}

// InlineCachedVoice creates a cached voice inline query result.
func InlineCachedVoice(id string, fileID string, title string) InlineQueryResultCachedVoice {
	return bot.InlineCachedVoice(id, fileID, title)
}

// InlineCachedAudio creates a cached audio inline query result.
func InlineCachedAudio(id string, fileID string) InlineQueryResultCachedAudio {
	return bot.InlineCachedAudio(id, fileID)
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

// InlineButtonWebApp creates an inline keyboard button that opens a Web App.
func InlineButtonWebApp(text string, url string) InlineKeyboardButton {
	return telegram.InlineButtonWebApp(text, url)
}

// InlineButtonLoginURL creates an inline keyboard button with Telegram login authorization.
func InlineButtonLoginURL(text string, rawURL string) InlineKeyboardButton {
	return telegram.InlineButtonLoginURL(text, rawURL)
}

// InlineButtonSwitchInlineQuery creates an inline keyboard button that switches to inline mode.
func InlineButtonSwitchInlineQuery(text string, query string) InlineKeyboardButton {
	return telegram.InlineButtonSwitchInlineQuery(text, query)
}

// InlineButtonSwitchInlineQueryCurrentChat creates an inline keyboard button that switches inline mode in the current chat.
func InlineButtonSwitchInlineQueryCurrentChat(text string, query string) InlineKeyboardButton {
	return telegram.InlineButtonSwitchInlineQueryCurrentChat(text, query)
}

// InlineButtonSwitchInlineQueryChosenChat creates an inline keyboard button that switches inline mode in a chosen chat.
func InlineButtonSwitchInlineQueryChosenChat(text string, options SwitchInlineQueryChosenChat) InlineKeyboardButton {
	return telegram.InlineButtonSwitchInlineQueryChosenChat(text, options)
}

// InlineButtonCopyText creates an inline keyboard button that copies text to the clipboard.
func InlineButtonCopyText(text string, copyText string) InlineKeyboardButton {
	return telegram.InlineButtonCopyText(text, copyText)
}

// InlineButtonPay creates an inline keyboard button that pays an invoice.
func InlineButtonPay(text string) InlineKeyboardButton {
	return telegram.InlineButtonPay(text)
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

// KeyboardButtonUsers creates a reply keyboard button that requests users.
func KeyboardButtonUsers(text string, request KeyboardButtonRequestUsers) KeyboardButton {
	return telegram.KeyboardButtonUsers(text, request)
}

// KeyboardButtonChat creates a reply keyboard button that requests a chat.
func KeyboardButtonChat(text string, request KeyboardButtonRequestChat) KeyboardButton {
	return telegram.KeyboardButtonChat(text, request)
}

// KeyboardButtonManagedBot creates a reply keyboard button that requests a managed bot.
func KeyboardButtonManagedBot(text string, request KeyboardButtonRequestManagedBot) KeyboardButton {
	return telegram.KeyboardButtonManagedBot(text, request)
}

// KeyboardButtonWebApp creates a reply keyboard button that opens a Web App.
func KeyboardButtonWebApp(text string, url string) KeyboardButton {
	return telegram.KeyboardButtonWebApp(text, url)
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
