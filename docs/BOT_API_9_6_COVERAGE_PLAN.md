# Bot API 9.6 Coverage Plan

> **Historical note:** This plan is the working record for the Bot API 9.6 implementation stages. It is not the current public API contract. For current pre-v1 names and architecture choices, see [`PRE_V1_NOTES.md`](PRE_V1_NOTES.md), [`API_COVERAGE.md`](API_COVERAGE.md), and [`BOT_API_10_0_FINAL_AUDIT.md`](BOT_API_10_0_FINAL_AUDIT.md).

## Goal

Maintain complete Telegram Bot API 9.6 method, type, and update coverage with documented architecture differences while the public repository matures toward a future tagged release.

## Source of truth

- [Telegram Bot API documentation](https://core.telegram.org/bots/api)
- [Telegram Bot API changelog](https://core.telegram.org/bots/api-changelog), especially the April 3, 2026 Bot API 9.6 entry

This plan is a working checklist derived from the official documentation and changelog. Before each implementation slice, re-check the relevant official sections because the Bot API can change.

## Repository policy

- The public repository exists at <https://github.com/xDilettante/ai-gram>.
- Bot API 9.6 code coverage is complete with documented architecture differences.
- New tags and GitHub Releases require explicit maintainer approval.
- Local commits are still expected after verified local work.
- Final user reports should state whether push was intentionally skipped by project policy.


## Stage 88 audit result

Stage 88 compared the then-current local implementation with the official Telegram Bot API documentation and the April 3, 2026 Bot API 9.6 changelog. It found remaining gaps that were closed in Stages 89-99. The updated audit now lives in [`docs/BOT_API_9_6_AUDIT.md`](BOT_API_9_6_AUDIT.md).

Top pending method groups after Stage 100:

- No official method wrappers are missing (169/169 official methods are represented by exported `(*bot.Bot)` methods).
- No known method behavior blockers remain after `setWebhook.certificate` upload/multipart support was implemented in Stage 99.

Top pending type/field groups:

- `Message.giveaway` was corrected in Stage 98.
- Optional concrete chat member variants remain a possible compatibility-preserving refinement; the flat `telegram.ChatMember` struct has the audited official fields.



## Stage 98/99 final audit result

Stage 98 created [`docs/BOT_API_9_6_FINAL_AUDIT.md`](BOT_API_9_6_FINAL_AUDIT.md). The audit found all 169 official Bot API method wrappers present and no missing fields in the audited high-impact `User`, `Chat`, `ChatFullInfo`, `Update`, `Message`, `ReplyParameters`, `CallbackQuery`, `Video`, sticker, and keyboard field tables after adding `Message.giveaway`. Stage 99 added `SetWebhook` certificate upload support with `FileUpload`, JSON/multipart tests, secret redaction checks, and documentation updates. No known Bot API 9.6 code coverage blockers remain. Stage 100 added [`docs/maintainer/BOT_API_9_6_RELEASE_READINESS.md`](maintainer/BOT_API_9_6_RELEASE_READINESS.md) for verification and manual-only smoke planning. Later stages published `main` only after explicit approval; tags and GitHub Releases still require explicit maintainer approval.

## Stage 89 result

Stage 89 implemented lifecycle/profile read APIs: `logOut`, `close`, `getUserProfilePhotos`, `getUserProfileAudios`, and `getForumTopicIconStickers`. Stage 97 resolved the `getChat` result mismatch compatibly by keeping `GetChat` as `*telegram.Chat` and adding `GetChatFullInfo` for the official full result.

## Stage 90 result

Stage 90 implemented verification and user status APIs: `setUserEmojiStatus`, `verifyUser`, `verifyChat`, `removeUserVerification`, and `removeChatVerification`. These methods are state-changing and remain manual-only for live smoke.

## Stage 91 result

Stage 91 implemented chat member updates, chat boost updates, `getUserChatBoosts`, `setChatMemberTag`, `banChatSenderChat`, and `unbanChatSenderChat`. It extends the existing flat `telegram.ChatMember` struct with official tag/admin/restricted fields and keeps live checks manual-only.

## Stage 92 result

Stage 92 implemented Stars subscription invite links: `createChatSubscriptionInviteLink`, `editChatSubscriptionInviteLink`, and `ChatInviteLink.subscription_period` / `ChatInviteLink.subscription_price`. These methods are payment-related and state-changing, so live checks remain manual-only.

## Stage 93 result

Stage 93 implemented structured poll options (`InputPollOption` and `Poll.question_entities`), checklist message/service types, `sendChecklist`, `editMessageChecklist`, and `sendMessageDraft`. Checklist and draft flows mutate user-visible state, so live checks remain manual-only.

## Current implemented baseline

### Core

- `aigram.New`, `aigram.NewBot`, and `bot.New`.
- `bot.Config`, typed `bot.ChatID`, configurable base URL and HTTP client.
- Private token storage, token-safe diagnostics, and no public raw token accessor.
- Typed Telegram API error handling with `errors.APIError`.
- JSON and multipart request paths.

### Updates, webhook, and long polling

- `GetUpdates`.
- Managed long polling runner.
- Inbound webhook handler.
- Webhook management: `SetWebhook`, `DeleteWebhook`, `GetWebhookInfo`.
- Dispatcher predicates/routes for messages, commands, callbacks, and chat join requests.

### Send and media

- `SendMessage`.
- `SendPhoto`, `SendDocument`, `SendVideo`, `SendAudio`, `SendVoice`.
- `SendContact`, `SendLocation`, `SendVenue`.
- `SendPoll`, `StopPoll`, `SendDice`.
- `SendSticker`, `SendAnimation`, `SendVideoNote`.
- `SendMediaGroup` with photo/video/audio/document input media.

### Files

- `GetFile`.
- `DownloadFile`.
- `FileID`, `FileURL`, and `FileUpload` helpers.

### Reply markup, callback, edit, and delete

- Inline keyboard, reply keyboard, force reply, and remove keyboard types.
- `AnswerCallbackQuery`.
- `EditMessageText`, `EditMessageCaption`, `EditMessageReplyMarkup`.
- `DeleteMessage`.
- `ReplyParameters` and message thread IDs for implemented methods.

### Forward and copy

- `ForwardMessage`.
- `CopyMessage`.

### Chat actions and pinning

- `SendChatAction`.
- `PinChatMessage`, `UnpinChatMessage`, `UnpinAllChatMessages`.

### Chat and member info

- `GetChat`.
- `GetChatMember`.
- `GetChatAdministrators`.
- `GetChatMemberCount`.

### Moderation, chat management, and administration

- `BanChatMember`, `UnbanChatMember`, `RestrictChatMember`.
- `SetChatTitle`, `SetChatDescription`, `SetChatPhoto`, `DeleteChatPhoto`, and `LeaveChat`.
- `SetChatStickerSet` and `DeleteChatStickerSet`.
- Forum topic methods and service message types.
- `PromoteChatMember`.
- `SetChatAdministratorCustomTitle`.
- `SetChatPermissions`.
- `telegram.ChatPermissions` and minimal administrator/member types.

### Invite links and join requests

- `ExportChatInviteLink`.
- `CreateChatInviteLink`.
- `EditChatInviteLink`.
- `CreateChatSubscriptionInviteLink`.
- `EditChatSubscriptionInviteLink`.
- `RevokeChatInviteLink`.
- `ApproveChatJoinRequest`.
- `DeclineChatJoinRequest`.
- `telegram.ChatInviteLink` and `telegram.ChatJoinRequest`.

### Commands and menu

- `SetMyCommands`, `DeleteMyCommands`, `GetMyCommands`.
- `SetChatMenuButton`, `GetChatMenuButton`.
- `SetMyDefaultAdministratorRights`.
- Command scopes, menu buttons, Web App menu info, and default administrator rights.

### Access, examples, deploy, and smoke

- Access control middleware and admin-only examples.
- Webhook and long polling examples.
- Media upload and live smoke helpers.
- Deploy/manual smoke harness, safe logs, Telegram reports, and redaction rules.

## Bot API 9.6 coverage checklist

### Chat management

- [x] `SetChatTitle`
- [x] `SetChatDescription`
- [x] `SetChatPhoto`
- [x] `DeleteChatPhoto`
- [x] `LeaveChat`
- [x] `SetChatStickerSet`
- [x] `DeleteChatStickerSet`

### Forum topics

- [x] `CreateForumTopic`
- [x] `EditForumTopic`
- [x] `CloseForumTopic`
- [x] `ReopenForumTopic`
- [x] `DeleteForumTopic`
- [x] `UnpinAllForumTopicMessages`
- [x] `EditGeneralForumTopic`
- [x] `CloseGeneralForumTopic`
- [x] `ReopenGeneralForumTopic`
- [x] `HideGeneralForumTopic`
- [x] `UnhideGeneralForumTopic`
- [x] `UnpinAllGeneralForumTopicMessages`
- [x] `ForumTopic`
- [x] `ForumTopicCreated`
- [x] `ForumTopicEdited`
- [x] `ForumTopicClosed`
- [x] `ForumTopicReopened`
- [x] `GeneralForumTopicHidden`
- [x] `GeneralForumTopicUnhidden`
- [x] `GetForumTopicIconStickers` and topic icon sticker coverage

## Stage 94 result

Stage 94 implemented reply and message metadata coverage: `MessageOrigin` variants, `ExternalReplyInfo`, `TextQuote`, `InaccessibleMessage`, `MaybeInaccessibleMessage`, `ReplyParameters` cross-chat/quote/checklist fields, and high-impact `Message` metadata fields such as `forward_origin`, `reply_to_message`, `external_reply`, `quote`, `reply_to_story`, `direct_messages_topic`, `suggested_post_info`, `pinned_message`, sender metadata, caption/media flags, and star/effect metadata. Callback queries now preserve the existing `Message` field for accessible messages and expose `MaybeMessage` for inaccessible callback messages.

## Stage 95 result

Stage 95 implemented prepared inline message and reply-markup completion coverage: `SavePreparedInlineMessage`, `PreparedInlineMessage`, `LoginURL`, `SwitchInlineQueryChosenChat`, `CopyTextButton`, `KeyboardButtonPollType`, `KeyboardButton.request_poll`, `InlineKeyboardButton.pay`, and keyboard button `icon_custom_emoji_id`/`style` fields. Prepared inline and rich button flows remain manual-only for live testing.

### Reactions

- [x] `SetMessageReaction`
- [x] `ReactionTypeEmoji`
- [x] `ReactionTypeCustomEmoji`
- [x] `ReactionTypePaid`
- [x] `MessageReactionUpdated`
- [x] `MessageReactionCountUpdated`
- [x] update fields for reaction updates
- [x] Full reply/forward message metadata pass for reaction-adjacent message decoding (Stage 94)

### Inline mode

Stage 72 implements the first inline foundation: incoming inline query updates, chosen inline result updates, `AnswerInlineQuery`, article results, text input content, and dispatch helpers. Stage 73 adds non-media inline results and the remaining input message content variants, including inline invoice content. Stage 74 adds media and cached inline result variants. Stage 75 audits these types against the official Bot API 9.6 inline mode documentation; no additional inline-only result or input content variants remain pending.

- [x] `InlineQuery` update
- [x] `ChosenInlineResult` update
- [x] `AnswerInlineQuery`
- [x] `InlineQueryResultArticle`
- [x] `InlineQueryResultPhoto`
- [x] `InlineQueryResultGif`
- [x] `InlineQueryResultMpeg4Gif`
- [x] `InlineQueryResultVideo`
- [x] `InlineQueryResultAudio`
- [x] `InlineQueryResultVoice`
- [x] `InlineQueryResultDocument`
- [x] `InlineQueryResultLocation`
- [x] `InlineQueryResultVenue`
- [x] `InlineQueryResultContact`
- [x] `InlineQueryResultGame`
- [x] cached inline result variants
- [x] `InputTextMessageContent`
- [x] `InputLocationMessageContent`
- [x] `InputVenueMessageContent`
- [x] `InputContactMessageContent`
- [x] `InputInvoiceMessageContent`
- [x] inline mode dispatcher predicates/helpers

### Payments, invoices, stars, and paid media

- [x] `SendInvoice`
- [x] `CreateInvoiceLink`
- [x] `AnswerShippingQuery`
- [x] `AnswerPreCheckoutQuery`
- [x] `RefundStarPayment`
- [x] `SendPaidMedia`
- [x] `GetMyStarBalance`
- [x] gift methods: `GetAvailableGifts`, `SendGift`, and `GiftPremiumSubscription`
- [x] available gift and owned gift methods/types
- [x] `Invoice`
- [x] `SuccessfulPayment`
- [x] `ShippingQuery`
- [x] `PreCheckoutQuery`
- [x] `LabeledPrice`
- [x] `ShippingOption`
- [x] `RefundedPayment` message type
- [x] `PaidMediaInfo`
- [x] paid media input/result types
- [x] `StarTransaction` and revenue-related transaction partner types, including gift-specific partner payloads

### Stickers

- [x] `GetStickerSet`
- [x] `GetCustomEmojiStickers`
- [x] `UploadStickerFile`
- [x] `CreateNewStickerSet`
- [x] `AddStickerToSet`
- [x] `SetStickerPositionInSet`
- [x] `DeleteStickerFromSet`
- [x] `ReplaceStickerInSet`
- [x] `SetStickerEmojiList`
- [x] `SetStickerKeywords`
- [x] `SetStickerMaskPosition`
- [x] `SetStickerSetTitle`
- [x] `SetStickerSetThumbnail`
- [x] `SetCustomEmojiStickerSetThumbnail`
- [x] `DeleteStickerSet`
- [x] `StickerSet`
- [x] `InputSticker`
- [x] `MaskPosition`
- [x] minimal sticker, custom emoji, sticker set, mask position, and input sticker fields needed for Bot API 9.6 sticker set management

### WebApp, prepared buttons, and Mini App related coverage

Stages 80 and 95 implement the remaining Bot API WebApp / Mini App support that belongs to the core Bot API surface. Mini App JavaScript methods such as `requestChat` live in the client-side WebApp API and are not Go Bot API methods.

- [x] `AnswerWebAppQuery`
- [x] `SentWebAppMessage`
- [x] `WebAppData` / `Message.web_app_data`
- [x] `WriteAccessAllowed` / `Message.write_access_allowed`
- [x] `WebAppInfo` audit
- [x] Web App fields in `InlineKeyboardButton`, `KeyboardButton`, `MenuButtonWebApp`, and `InlineQueryResultsButton`
- [x] `PreparedKeyboardButton`
- [x] `SavePreparedKeyboardButton`
- [x] `KeyboardButtonRequestUsers`
- [x] `KeyboardButtonRequestChat`
- [x] `KeyboardButtonRequestManagedBot`
- [x] `KeyboardButtonPollType` / `KeyboardButton.request_poll`
- [x] `LoginURL` / `InlineKeyboardButton.login_url`
- [x] `SwitchInlineQueryChosenChat` / switch-inline button fields
- [x] `CopyTextButton` / `InlineKeyboardButton.copy_text`
- [x] `InlineKeyboardButton.pay`
- [x] button `icon_custom_emoji_id` and `style` fields
- [x] `PreparedInlineMessage`
- [x] `SavePreparedInlineMessage`

### Managed Bots, Bot API 9.6

Stage 78 implements Managed Bots 9.6 support. Token-returning methods are sensitive and remain manual-only for live checks.

- [x] `User.can_manage_bots`
- [x] `KeyboardButtonRequestManagedBot`
- [x] `KeyboardButton.request_managed_bot`
- [x] `ManagedBotCreated`
- [x] `Message.managed_bot_created`
- [x] `ManagedBotUpdated`
- [x] `Update.managed_bot`
- [x] `GetManagedBotToken`
- [x] `ReplaceManagedBotToken`
- [x] `PreparedKeyboardButton`
- [x] `SavePreparedKeyboardButton`

### Poll 9.6 updates

- [x] `Poll.correct_option_ids`
- [x] add plural `correct_option_ids` support while keeping old single `correct_option_id` for backward compatibility
- [x] `SendPoll.correct_option_ids`
- [x] quiz support with `allows_multiple_answers`, as allowed by Bot API 9.6
- [x] `Poll.allows_revoting`
- [x] `SendPoll.allows_revoting`
- [x] `SendPoll.shuffle_options`
- [x] `SendPoll.allow_adding_options`
- [x] `SendPoll.hide_results_until_closes`
- [x] `Poll.description` and `Poll.description_entities`
- [x] `SendPoll.description`, `description_parse_mode`, and `description_entities`
- [x] `PollOption.persistent_id`
- [x] `PollAnswer.option_persistent_ids`
- [x] `PollOption.added_by_user`, `added_by_chat`, and `addition_date`
- [x] `PollOptionAdded` and `Message.poll_option_added`
- [x] `PollOptionDeleted` and `Message.poll_option_deleted`
- [x] `ReplyParameters.poll_option_id`
- [x] `Message.reply_to_poll_option_id`
- [x] `InputPollOption` and `Poll.question_entities` entity-aware poll option/question coverage
- [x] audit all poll-related fields against official Bot API 9.6 docs for this slice

### Bot profile and metadata

- [x] `SetMyName`
- [x] `GetMyName`
- [x] `SetMyDescription`
- [x] `GetMyDescription`
- [x] `SetMyShortDescription`
- [x] `GetMyShortDescription`
- [x] `SetMyProfilePhoto` with Bot API 9.6 `InputProfilePhotoStatic`/`InputProfilePhotoAnimated` upload-only multipart payloads
- [x] `RemoveMyProfilePhoto`
- [x] `GetMyDefaultAdministratorRights`
- [x] typed bot name/description/short description objects

Live smoke for this slice is manual-only because set/remove operations change real bot profile state.

### Business APIs

Stage 81 implements the Business API foundation. Stage 82 adds business read, account profile, gift settings, story, and suggested post methods. Stage 83 adds gifts, business gifts, bot/business Star balances and transfers, Premium subscription gifts, and Stars subscription editing. Stage 84 adds `business_connection_id` support to currently implemented send/edit-style methods that expose it in official Bot API docs. Stage 85 adds `EditMessageMedia`, `EditMessageLiveLocation`, `StopMessageLiveLocation`, and `InputMediaAnimation`. Broader account metadata and other not-yet-implemented methods were handled in later slices or documented as compatibility choices.

- [x] `BusinessConnection`
- [x] `BusinessBotRights`
- [x] `BusinessMessagesDeleted`
- [x] `Update.business_connection`
- [x] `Update.business_message`
- [x] `Update.edited_business_message`
- [x] `Update.deleted_business_messages`
- [x] business message metadata fields: `business_connection_id`, `sender_business_bot`, `is_from_offline`
- [x] business dispatch helpers
- [x] `GetBusinessConnection`
- [x] `DeleteBusinessMessages`
- [x] `ReadBusinessMessage`
- [x] `SetBusinessAccountName`
- [x] `SetBusinessAccountUsername`
- [x] `SetBusinessAccountBio`
- [x] `SetBusinessAccountProfilePhoto`
- [x] `RemoveBusinessAccountProfilePhoto`
- [x] `AcceptedGiftTypes` and `SetBusinessAccountGiftSettings`
- [x] `Story`, `InputStoryContent`, `StoryArea`, `PostStory`, `EditStory`, and `DeleteStory`
- [x] `ApproveSuggestedPost`, `DeclineSuggestedPost`, and suggested post service message types
- [x] `business_connection_id` on current send methods: `SendMessage`, `SendPhoto`, `SendDocument`, `SendVideo`, `SendAudio`, `SendVoice`, `SendAnimation`, `SendVideoNote`, `SendSticker`, `SendPaidMedia`, `SendMediaGroup`, `SendContact`, `SendLocation`, `SendVenue`, `SendPoll`, and `SendDice`
- [x] `business_connection_id` on current chat action/pin methods: `SendChatAction`, `PinChatMessage`, and `UnpinChatMessage`
- [x] `business_connection_id` on current edit/stop methods: `EditMessageText`, `EditMessageCaption`, `EditMessageReplyMarkup`, `EditMessageMedia`, `EditMessageLiveLocation`, `StopMessageLiveLocation`, and `StopPoll`
- [x] `SendChecklist`, `EditMessageChecklist`, and `SendMessageDraft`
- [x] `GetBusinessAccountStarBalance`
- [x] `TransferBusinessAccountStars`
- [x] `Gift`, `Gifts`, `GiftInfo`, `UniqueGift`, `UniqueGiftInfo`, `OwnedGift`, and `OwnedGifts` types
- [x] `GetAvailableGifts`
- [x] `SendGift`
- [x] `GiftPremiumSubscription`
- [x] `GetBusinessAccountGifts`
- [x] `GetUserGifts`
- [x] `GetChatGifts`
- [x] `ConvertGiftToStars`
- [x] `UpgradeGift`
- [x] `TransferGift`
- [x] `GetMyStarBalance`
- [x] `EditUserStarSubscription`
- [x] `ChatFullInfo` business intro/location/hours metadata coverage

### Games

- [x] `SendGame`
- [x] `SetGameScore`
- [x] `GetGameHighScores`
- [x] `CallbackGame`
- [x] `Game`
- [x] `GameHighScore`

### Passport

- [x] `PassportData`
- [x] encrypted passport element types
- [x] passport credentials/files/errors
- [x] `SetPassportDataErrors`

Notes:

- Passport decryption helpers are intentionally out of scope for the typed Bot API wrapper unless a future product decision adds them.
- Passport live checks are manual-only and must never log encrypted payloads or user documents.

### Batch methods

- [x] `ForwardMessages`
- [x] `CopyMessages`
- [x] `DeleteMessages`

### Remaining message and edit methods

- [x] `EditMessageMedia`
- [x] `EditMessageLiveLocation`
- [x] `StopMessageLiveLocation`
- [x] `SendChecklist`
- [x] `EditMessageChecklist`
- [x] `SendMessageDraft`
- [x] `logOut` and `close`
- [x] `GetUserProfilePhotos`, `GetUserProfileAudios`, and `GetForumTopicIconStickers`
- [x] verification/status methods
- [x] chat boost/member update methods and sender-chat moderation


## Stage 96 result

Stage 96 implemented service/direct-message/story/media metadata completion: `RepostStory`, video cover/start/quality metadata, `SendVideo` thumbnail/cover/start/caption-placement serialization, shared user/chat service messages, chat backgrounds, video chat service messages, proximity alerts, auto-delete timer changes, giveaway service fields, and paid/direct message price-change service fields. `RepostStory` remains manual-only because it mutates business story state; service-message and media metadata coverage is verified through unit fixtures.

## Stage 97 result

Stage 97 implemented `ChatFullInfo` and `GetChatFullInfo`, completed official Bot API 9.6 `User` and lightweight `Chat` metadata fields, and added `Update.channel_post`, `Update.edited_channel_post`, `Update.poll`, effective helper support, and dispatch predicates/handler helpers. `GetChat` remains backward-compatible and returns `*telegram.Chat`; `GetChatFullInfo` decodes the official full `getChat` result. The flat `telegram.ChatMember` shape remains the compatibility strategy for official chat member variants.

## Implementation strategy

Completed implementation stages after the Stage 98 audit:

1. Stage 89 completed: lifecycle/profile read APIs (`logOut`, `close`, profile photos/audios, forum topic icon stickers).
2. Stage 90 completed: verification and user status APIs.
3. Stage 91 completed: chat boosts, chat-member updates/tags, and sender-chat moderation.
4. Stage 92 completed: subscription invite links.
5. Stage 93 completed: checklists, message drafts, and structured poll options.
6. Stage 94 completed: reply and message metadata types.
7. Stage 95 completed: prepared inline messages and reply-markup completion.
8. Stage 96 completed: service/direct-message/story/media metadata gaps.
9. Stage 97 completed: `ChatFullInfo`, full user/chat metadata, channel post/standalone poll update shape, and compatible flat chat member variant strategy.
10. Stage 98 completed: final official-doc audit and release-readiness blocker review.
11. Stage 99 completed: `SetWebhook` certificate upload / multipart support and focused final audit.
12. Stage 100 completed: release-readiness verification and manual-only smoke planning.

Future API maintenance stages should:

- re-check the relevant official Bot API 9.6 section;
- add typed params, result types, and minimal Telegram data types;
- add validation that is stable and not over-fitted to uncertain upper limits;
- add httptest/unit coverage for success, validation, API errors, invalid JSON, HTTP errors, cancelled context, and token leakage;
- update README, CHANGELOG, API coverage, and manual testing docs;
- create a local commit after checks pass;
- skip push/tag/release unless the user explicitly approves that exact publication action.

## Live smoke policy

- Safe methods can be live-smoked with generated data.
- State-changing/admin/destructive methods are manual-only.
- Payments, Passport, Business, Managed Bots, gifts, and Stars flows require explicit confirmation.
- No automatic destructive/admin smoke.
- Manual live smoke must use dedicated test chats, test accounts, or sandbox-like flows whenever possible.
- Logs and reports must not print bot tokens, webhook secrets, token-bearing URLs, full invite links, full payment payloads, or private message text.
