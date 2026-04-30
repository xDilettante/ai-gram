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
- [x] `RefundStarPayment`
- [x] `SendPaidMedia`
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
- [x] `PaidMediaInfo`
- [x] paid media input/result types
- [x] `StarTransaction` and basic revenue-related transaction partner types
- [ ] Advanced gift-specific transaction partner payloads remain pending for the gifts/business gifts slice

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

Stage 80 implements the remaining Bot API WebApp / Mini App support that belongs to the core Bot API surface. Mini App JavaScript methods such as `requestChat` live in the client-side WebApp API and are not Go Bot API methods.

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
- [ ] `date_time` entity support in poll/checklist/gift-related contexts where relevant
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

Stage 81 implements the Business API foundation only. Business send/edit/read/account-profile methods and business gifts remain separate local-only slices.

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
- [ ] business send/edit/read methods beyond `DeleteBusinessMessages`, if present in Bot API 9.6 docs
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
7. WebApp and prepared buttons - implemented locally in Stage 80; Mini App live checks manual-only.
8. Managed Bots 9.6 - implemented locally in Stage 78; token-returning methods manual-only.
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
