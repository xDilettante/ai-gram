# Bot API 9.6 Coverage Plan

## Goal

Reach 100% Telegram Bot API 9.6 method, type, and update coverage before the next push, tag, or GitHub Release.

## Source of truth

- [Telegram Bot API documentation](https://core.telegram.org/bots/api)
- [Telegram Bot API changelog](https://core.telegram.org/bots/api-changelog), especially the April 3, 2026 Bot API 9.6 entry

This plan is a working checklist derived from the official documentation and changelog. Before each implementation slice, re-check the relevant official sections because the Bot API can change.

## Repository policy

- `v0.2.0` remains the latest public release.
- Local development continues toward full Bot API 9.6 coverage.
- Push, tag, and GitHub Release operations are frozen until full Bot API 9.6 coverage is reached.
- Local commits are still expected after verified local work.
- Final user reports may say: "local commit created; push intentionally skipped by project policy".

## Current implemented baseline

### Core

- `aigram.New`, `aigram.NewBot`, and `bot.New`.
- `bot.BotConfig`, typed `bot.ChatID`, configurable base URL and HTTP client.
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

## Missing Bot API 9.6 areas

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
- [ ] audit topic icon sticker fields against Bot API 9.6

### Reactions

- [x] `SetMessageReaction`
- [x] `ReactionTypeEmoji`
- [x] `ReactionTypeCustomEmoji`
- [x] `ReactionTypePaid`
- [x] `MessageReactionUpdated`
- [x] `MessageReactionCountUpdated`
- [x] update fields for reaction updates
- [ ] reaction count/list fields on messages, if missing

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
- [ ] `RefundStarPayment`, if present in Bot API 9.6 docs
- [ ] `SendPaidMedia`, if present in Bot API 9.6 docs
- [ ] `GetMyStarBalance`, if present in Bot API 9.6 docs
- [ ] gift methods such as `SendGift` and `GiftPremiumSubscription`, if present in Bot API 9.6 docs
- [ ] available gift and owned gift methods/types, if present in Bot API 9.6 docs
- [x] `Invoice`
- [x] `SuccessfulPayment`
- [x] `ShippingQuery`
- [x] `PreCheckoutQuery`
- [x] `LabeledPrice`
- [x] `ShippingOption`
- [x] `RefundedPayment` message type
- [ ] `PaidMediaInfo`
- [ ] paid media input/result types
- [ ] `StarTransaction` and revenue-related types, if present in Bot API 9.6 docs

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

- [ ] `AnswerWebAppQuery`
- [ ] `SentWebAppMessage`
- [ ] Web App fields in keyboard buttons and inline buttons, beyond current menu button support
- [ ] Web App `requestChat` support from Bot API 9.6
- [ ] Web App data message fields and helpers
- [ ] `PreparedKeyboardButton`
- [ ] `SavePreparedKeyboardButton`
- [ ] `KeyboardButtonRequestUsers`
- [ ] `KeyboardButtonRequestChat`
- [ ] `KeyboardButtonRequestManagedBot`, if present in Bot API 9.6 docs

### Managed Bots, Bot API 9.6

- [ ] `User.can_manage_bots`
- [ ] `KeyboardButtonRequestManagedBot`
- [ ] `KeyboardButton.request_managed_bot`
- [ ] `ManagedBotCreated`
- [ ] `Message.managed_bot_created`
- [ ] `ManagedBotUpdated`
- [ ] `Update.managed_bot`
- [ ] `GetManagedBotToken`
- [ ] `ReplaceManagedBotToken`
- [ ] `PreparedKeyboardButton`
- [ ] `SavePreparedKeyboardButton`

### Poll 9.6 updates

- [ ] `Poll.correct_option_ids`
- [ ] replace old single `correct_option_id` usage where Bot API 9.6 requires plural support
- [ ] `SendPoll.correct_option_ids`
- [ ] quiz support with `allows_multiple_answers`, as allowed by Bot API 9.6
- [ ] `Poll.allows_revoting`
- [ ] `SendPoll.allows_revoting`
- [ ] `SendPoll.shuffle_options`
- [ ] `SendPoll.allow_adding_options`
- [ ] `SendPoll.hide_results_until_closes`
- [ ] `Poll.description` and `Poll.description_entities`
- [ ] `SendPoll.description`, `description_parse_mode`, and `description_entities`
- [ ] `PollOption.persistent_id`
- [ ] `PollAnswer.option_persistent_ids`
- [ ] `PollOption.added_by_user`, `added_by_chat`, and `addition_date`
- [ ] `PollOptionAdded` and `Message.poll_option_added`
- [ ] `PollOptionDeleted` and `Message.poll_option_deleted`
- [ ] `ReplyParameters.poll_option_id`
- [ ] `Message.reply_to_poll_option_id`
- [ ] `date_time` entity support in poll/checklist/gift-related contexts where relevant
- [ ] audit all poll-related fields against official Bot API 9.6 docs

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

- [ ] `BusinessConnection`
- [ ] `BusinessMessagesDeleted`
- [ ] business-related update fields
- [ ] business connection helpers
- [ ] business send/edit/delete/read methods, if present in Bot API 9.6 docs
- [ ] business account profile/name/username/bio methods, if present in Bot API 9.6 docs
- [ ] business account Star balance and transfer methods, if present in Bot API 9.6 docs
- [ ] gifts/stars/business account gifts, if present in Bot API 9.6 docs
- [ ] business intro/location/hours/account metadata, if present in Bot API 9.6 docs

### Games

- [ ] `SendGame`
- [ ] `SetGameScore`
- [ ] `GetGameHighScores`
- [ ] `CallbackGame`
- [ ] `Game`
- [ ] `GameHighScore`

### Passport

- [ ] `PassportData`
- [ ] encrypted passport element types
- [ ] passport credentials/files/errors
- [ ] `SetPassportDataErrors`

### Batch methods

- [x] `ForwardMessages`
- [x] `CopyMessages`
- [x] `DeleteMessages`

### Remaining message and edit methods

- [ ] `EditMessageMedia`
- [ ] `EditMessageLiveLocation`
- [ ] `StopMessageLiveLocation`
- [ ] `SendChecklist`, if present in Bot API 9.6 docs
- [ ] `EditMessageChecklist`, if present in Bot API 9.6 docs
- [ ] `SendMessageDraft`, if present in Bot API 9.6 docs
- [ ] any missing methods discovered by a final official-doc audit

## Implementation strategy

Recommended local-only stages:

1. Chat management - implemented locally in Stage 66; manual-only live smoke.
2. Forum topics - implemented locally in Stage 67; manual-only live smoke.
3. Reactions - implemented locally in Stage 68; manual-only live smoke.
4. Inline mode basics
5. Sticker set management
6. Payments, stars, and paid media
7. WebApp and prepared buttons
8. Managed Bots 9.6
9. Poll 9.6 updates
10. Business APIs
11. Games and Passport
12. Batch methods - implemented locally in Stage 69; `DeleteMessages` manual-only live smoke.
13. Remaining message/edit methods
14. Final full coverage audit against official Bot API 9.6

Each stage should:

- re-check the relevant official Bot API 9.6 section;
- add typed params, result types, and minimal Telegram data types;
- add validation that is stable and not over-fitted to uncertain upper limits;
- add httptest/unit coverage for success, validation, API errors, invalid JSON, HTTP errors, cancelled context, and token leakage;
- update README, CHANGELOG, API coverage, and manual testing docs;
- create a local commit after checks pass;
- skip push/tag/release by policy.

## Live smoke policy

- Safe methods can be live-smoked with generated data.
- State-changing/admin/destructive methods are manual-only.
- Payments, Passport, Business, Managed Bots, gifts, and Stars flows require explicit confirmation.
- No automatic destructive/admin smoke.
- Manual live smoke must use dedicated test chats, test accounts, or sandbox-like flows whenever possible.
- Logs and reports must not print bot tokens, webhook secrets, token-bearing URLs, full invite links, full payment payloads, or private message text.
