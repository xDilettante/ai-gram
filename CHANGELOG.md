# Changelog

## Unreleased

No unreleased changes.

## v0.6.0 - 2026-05-14

### Added

- Added a maintainer Bot API update checklist for future Telegram API compatibility audits.
- Added a Bot API 10.0 lightweight freshness audit against the current official Telegram Bot API documentation.
- Added a maintainer release-candidate checklist for the accumulated post-`v0.5.0` work.
- Added release notes for the `v0.6.0` release.
- Added a `callback` package for compact typed inline keyboard callback data, including encode/parse helpers, conventional confirm/cancel actions, pagination helpers, expiry checks, and callback button construction.
- Added `dispatch.CallbackAction` and `Dispatcher.OnCallbackActionFunc` for routing parsed typed callback data by namespace and action.
- Added `errors` helpers for classifying Telegram API errors, rate limits, chat migrations, forbidden/not-found responses, network errors, and context cancellation.
- Added `telegram.Actor` and identity helpers for message/update/callback actors, anonymous admin detection, and reply target actors.

### Examples

- Added shared example log masking for numeric private IDs and updated long polling/webhook examples to avoid raw `chat_id`, `from_user_id`, and `by_user_id` values in logs.
- Updated `examples/04_inline_keyboard` to use typed callback data instead of ad hoc callback strings.
- Added `examples/07_inline_panel` with typed callbacks, pagination, confirm/cancel actions, dispatcher routing, and long polling.
- Added `examples/08_retry_sender` with explicit retry/rate-limit-aware sending using the public error classification helpers.
- Added `examples/09_group_admin` with safe read-only group/admin identity commands using `telegram.Actor` helpers.
- Added `examples/10_moderation_skeleton` with dry-run reports, moderation previews, and join-request logging without destructive Bot API calls.
- Added `examples/11_transport_parity` with shared dispatcher/handlers and explicit long polling/webhook transport modes.

## v0.5.0 - 2026-05-11

### Changed

- Breaking pre-v1 cleanup: renamed `bot.BotConfig` and root `aigram.BotConfig` to `Config`.
- Breaking pre-v1 cleanup: renamed `telegram.LoginUrl` to `LoginURL`.

### Documentation

- Added missing GoDoc for inline result JSON encoders and chat-member helper methods.
- Added pre-v1 notes and marked Bot API 9.6 workstream documents as historical where later cleanup superseded their compatibility notes.
- Added public consumer smoke tooling and a manual GitHub Actions workflow for verifying external module consumption.
- Hardened maintainer smoke tooling with dry-run support and broader identifier sanitization.

### Verification

- Verified public `main` consumption from an external temporary module after the `Config` / `LoginURL` cleanup.
- Added repeatable public consumer smoke checks for `main` and release tags.

## v0.4.0 - 2026-05-11

### Added

- Telegram Bot API 10.0 coverage on `main`, including guest mode, reaction deletion, poll media, live photos, managed bot access settings, personal chat message reads, and empty-text message drafts.
- Release-readiness documentation for the planned `v0.4.0` pre-v1 milestone.

### Changed

- Breaking pre-v1 cleanup: `GetChat` now returns `telegram.ChatFullInfo`, matching the official `getChat` result.
- Breaking pre-v1 cleanup: `telegram.ChatMember` is now an interface implemented by official `ChatMember*` variants.
- Breaking pre-v1 cleanup: `CallbackQuery.Message` now uses `telegram.MaybeInaccessibleMessage` directly.
- Breaking pre-v1 cleanup: the root `aigram` package is now a compact quick-start facade instead of a broad re-export mirror.

### Documentation

- Updated README, architecture notes, API coverage, roadmap, and release notes for Bot API 10.0 readiness.
- Documented intentional architecture differences for `FileRef`/`FileUpload`, `ChatFullInfo`, `ChatMember`, `MaybeInaccessibleMessage`, and live-photo URL handling.

### Verification

- Verified public `main` consumption from an external temporary module with root quick-start imports and direct `bot` / `telegram` imports.

## v0.3.0 - 2026-05-02

### Added

- Repository maturity pack: GitHub Actions CI, local coverage badge generation, architecture docs, public examples, contribution policy, security policy, issue forms, and PR template.
- Bot API 9.6 release-readiness documentation and local verification plan.
- SetWebhook certificate upload support for the official upload-only `certificate` InputFile parameter.
- Final Bot API 9.6 coverage audit documentation and `Message.giveaway` decoding coverage.
- ChatFullInfo, GetChatFullInfo, fuller User/Chat metadata, and channel post/standalone poll update shapes.
- Service/direct-message/story/media metadata completion: `RepostStory`, video cover/start/quality metadata, shared user/chat service messages, chat backgrounds, video chats, proximity alerts, auto-delete timers, giveaway service messages, and paid/direct message price-change service fields.
- Prepared inline messages and reply markup completion: `SavePreparedInlineMessage`, `PreparedInlineMessage`, `LoginUrl`, `SwitchInlineQueryChosenChat`, `CopyTextButton`, `KeyboardButtonPollType`, request-poll buttons, pay buttons, and icon/style button fields.
- Reply and message metadata types: `MessageOrigin` variants, `ExternalReplyInfo`, `TextQuote`, `MaybeInaccessibleMessage`, `InaccessibleMessage`, and expanded `ReplyParameters`.
- Checklist, message draft, and structured poll option support: `SendChecklist`, `EditMessageChecklist`, `SendMessageDraft`, `InputChecklist`, `Checklist`, and `InputPollOption`.
- Chat subscription invite link methods: `CreateChatSubscriptionInviteLink` and `EditChatSubscriptionInviteLink`.
- Chat member update, chat boost, and sender-chat moderation support: `ChatMemberUpdated`, `ChatBoostUpdated`, `ChatBoostRemoved`, `GetUserChatBoosts`, `SetChatMemberTag`, `BanChatSenderChat`, and `UnbanChatSenderChat`.
- Verification and user status methods: `SetUserEmojiStatus`, `VerifyUser`, `VerifyChat`, `RemoveUserVerification`, and `RemoveChatVerification`.
- Lifecycle and profile read APIs: `LogOut`, `Close`, `GetUserProfilePhotos`, `GetUserProfileAudios`, and `GetForumTopicIconStickers`.
- Sticker set management methods and `InputSticker` support.
- Chat management methods: `SetChatTitle`, `SetChatDescription`, `SetChatPhoto`, `DeleteChatPhoto`, `LeaveChat`, `SetChatStickerSet`, and `DeleteChatStickerSet`.
- Forum topic methods and forum topic service message types.
- Reaction types, reaction update decoding and dispatch helpers, and `SetMessageReaction`.
- Batch message methods: `ForwardMessages`, `CopyMessages`, and `DeleteMessages`.
- Bot profile and metadata methods: `SetMyName`, `GetMyName`, `SetMyDescription`, `GetMyDescription`, `SetMyShortDescription`, `GetMyShortDescription`, `GetMyDefaultAdministratorRights`, `SetMyProfilePhoto`, and `RemoveMyProfilePhoto`.
- Inline mode basics: `InlineQuery`, `ChosenInlineResult`, `AnswerInlineQuery`, `InlineQueryResultArticle`, and `InputTextMessageContent`.
- Additional inline mode result and input message content variants.
- Media and cached inline query result variants.
- Payments and invoice basics: `SendInvoice`, `CreateInvoiceLink`, `AnswerShippingQuery`, `AnswerPreCheckoutQuery`, and payment update/message types.
- Paid media and Stars basics: `SendPaidMedia`, `GetStarTransactions`, `RefundStarPayment`, paid media message/update types, and Star transaction decoding.
- Managed Bots 9.6 support: `User.can_manage_bots`, managed bot keyboard request/update types, `SavePreparedKeyboardButton`, `GetManagedBotToken`, and `ReplaceManagedBotToken`.
- Poll 9.6 fields and update support: plural correct options, revoting/options controls, poll option service messages, `PollAnswer`, and poll-option replies.
- WebApp / Mini App remaining support: `AnswerWebAppQuery`, `SentWebAppMessage`, `WebAppData`, `WriteAccessAllowed`, and Web App buttons.
- Business API foundation: `BusinessConnection`, `BusinessMessagesDeleted`, business update fields, `GetBusinessConnection`, and `DeleteBusinessMessages`.
- Business API account/read/story/suggested post methods: `ReadBusinessMessage`, account profile methods, gift settings, `PostStory`, `EditStory`, `DeleteStory`, `ApproveSuggestedPost`, and `DeclineSuggestedPost`.
- Gifts, business gifts, and remaining Stars methods: `GetAvailableGifts`, `SendGift`, `GiftPremiumSubscription`, `GetBusinessAccountStarBalance`, `TransferBusinessAccountStars`, `GetBusinessAccountGifts`, `GetUserGifts`, `GetChatGifts`, `ConvertGiftToStars`, `UpgradeGift`, `TransferGift`, `GetMyStarBalance`, and `EditUserStarSubscription`.
- `business_connection_id` support for supported send, edit, chat action, pin/unpin, media group, and poll-stop methods.
- Edit media and live location edit methods: `EditMessageMedia`, `EditMessageLiveLocation`, `StopMessageLiveLocation`, and `InputMediaAnimation`.
- Game methods and types: `SendGame`, `SetGameScore`, `GetGameHighScores`, `Game`, `CallbackGame`, and `GameHighScore`.
- Telegram Passport types and `SetPassportDataErrors`.

### Documentation

- Polished public documentation before repository publication by simplifying the README and checking moved maintainer links.
- Cleaned the public release surface by minimizing `.env.example` and separating maintainer-only smoke/deploy docs and examples.
- Added Bot API 9.6 release-readiness documentation and manual-only smoke planning.
- Added full Bot API 9.6 coverage audit documentation.
- Added final Bot API 9.6 coverage audit documentation and release-readiness blockers.
- Completed inline mode audit against official Bot API 9.6 documentation.

### Planning

- Added v0.3 plan for chat management, forum topics, reactions, and inline mode basics.
- Added full Telegram Bot API 9.6 coverage plan and local-only push/tag/release freeze policy.

## v0.2.0 - 2026-04-30

### Added

- Service/direct-message/story/media metadata completion: `RepostStory`, video cover/start/quality metadata, shared user/chat service messages, chat backgrounds, video chats, proximity alerts, auto-delete timers, giveaway service messages, and paid/direct message price-change service fields.
#### Send methods

- `SendContact`.
- `SendLocation`.
- `SendVenue`.
- `SendPoll`.
- `StopPoll`.
- `SendDice`.
- `SendSticker`.
- `SendAnimation`.
- `SendVideoNote`.

#### Media groups

- `SendMediaGroup` with `InputMediaPhoto`, `InputMediaVideo`, `InputMediaAudio`, and `InputMediaDocument`.

#### Bot commands/menu

- Bot commands and menu methods: `SetMyCommands`, `DeleteMyCommands`, `GetMyCommands`, `SetChatMenuButton`, `GetChatMenuButton`, and `SetMyDefaultAdministratorRights`.

#### Invite links and join requests

- Chat invite link methods: `ExportChatInviteLink`, `CreateChatInviteLink`, `EditChatInviteLink`, and `RevokeChatInviteLink`.
- Chat join request methods: `ApproveChatJoinRequest`, `DeclineChatJoinRequest`, plus `telegram.ChatJoinRequest` update decoding and dispatch predicates.

#### Admin management

- Admin management methods: `PromoteChatMember`, `SetChatAdministratorCustomTitle`, and `SetChatPermissions`.

#### Smoke and release planning

- Targeted v0.2 send-method live smoke script.
- Targeted `SendMediaGroup` live smoke script with generated upload fallback.
- v0.2 checkpoint document for release decision-making.

## v0.1.1 - 2026-04-30

### Fixed

- Corrected Go module path to `github.com/xDilettante/ai-gram`.
- Verified canonical public install via `go get github.com/xDilettante/ai-gram@v0.1.1`.

## v0.1.0

Date: 2026-04-30

First working `ai-gram` milestone: an early but verifiable Go library for the Telegram Bot API with a typed core, transports, dispatcher/middleware, examples, and manual/live smoke tooling.

### Added

- Service/direct-message/story/media metadata completion: `RepostStory`, video cover/start/quality metadata, shared user/chat service messages, chat backgrounds, video chats, proximity alerts, auto-delete timers, giveaway service messages, and paid/direct message price-change service fields.
#### Core

- Bot construction/config through `aigram.New`, `aigram.NewBot`, and `bot.New`.
- Typed `ChatID` helpers for numeric chat IDs and string IDs such as `@channelusername`.
- Internal JSON Bot API calls with configurable base URL and HTTP client.
- Multipart upload support for implemented media send methods.
- Typed Telegram API errors through `errors.APIError` and response parameters.
- Safe token redaction in bot string output, scripts, logs, and reports.

#### Updates/transports

- `GetUpdates` typed API call.
- Managed long polling runner with context cancellation, offset advancement, backoff, and handler error reporting.
- Inbound webhook HTTP handler with method/content-type checks, JSON decoding, optional secret-token validation, and handler error handling.
- Webhook management methods: `SetWebhook`, `DeleteWebhook`, and `GetWebhookInfo`.

#### Dispatch/middleware

- Dispatcher routes and predicates for messages, commands, callback queries, and exact callback data.
- Middleware chain for update handlers.
- Recovery middleware, per-update timeout middleware, and observability hooks.
- Access control middleware with admin/public/off modes and dynamic policy support.

#### Send/media/files

- `SendMessage`.
- `SendPhoto`.
- `SendDocument`.
- `SendVideo`.
- `SendAudio`.
- `SendVoice`.
- File references through `FileID`, `FileURL`, and `FileUpload`.
- `GetFile` and `DownloadFile` without exposing token-bearing download URLs.

#### Reply/callback/edit/delete

- Reply markup types for inline keyboards, reply keyboards, remove keyboard, and force reply.
- `AnswerCallbackQuery`.
- `EditMessageText`.
- `EditMessageCaption`.
- `EditMessageReplyMarkup`.
- `DeleteMessage`.
- `ReplyParameters` support for implemented send/copy methods.
- `MessageThreadID` support for implemented send methods and relevant operations.

#### Forward/copy

- `ForwardMessage`.
- `CopyMessage` with `telegram.MessageID` result.

#### Chat actions/pin

- `SendChatAction` with known action constants.
- `PinChatMessage`.
- `UnpinChatMessage`.
- `UnpinAllChatMessages`.

#### Chat/member/moderation

- `GetChat`.
- `GetChatMember`.
- `GetChatAdministrators`.
- `GetChatMemberCount`.
- `BanChatMember`.
- `UnbanChatMember`.
- `RestrictChatMember` with `telegram.ChatPermissions`.

#### Examples/testing/deploy

- Echo long polling example.
- Inline long polling example.
- Webhook server example.
- Media upload/download example.
- Local Bot API server smoke example.
- Deploy/manual smoke harness.
- Auto-discovery for deploy/manual smoke env.
- SSH alias support and separate local Bot API host support.
- Auto SSH tunnel support for remote loopback local Bot API servers.
- Telegram notifications for manual smoke/deploy checks.
- Deep-link smoke panels for examples.
- Safe logs that avoid tokens, webhook secrets, token-bearing URLs, and full private message text.
- Live smoke matrix and release checklist docs.

### Changed

- Removed raw public `Bot.Token()` access before `v0.1.0`.
- Updated README and docs to reflect implemented coverage, roadmap, safe smoke flows, and release readiness.
- Protected the webhook example with admin-only access mode by default.
- Made long polling smoke event-driven: it exits after the first successfully handled text update instead of waiting for the full timeout window.

### Security

- Raw bot token is not exposed by the public Bot API.
- Errors and logs avoid token leakage.
- Scripts redact bot tokens, webhook secrets, Telegram API hash values, and token-bearing Bot API URLs.
- `.env.local` and generated deploy env files are ignored.
- Examples use access control by default so test bots are not public by default.

### Live smoke verified

Safe live flows verified before the tag:

- Local Bot API smoke.
- Webhook deploy.
- Access panel including open/close access checks.
- Access denied behavior.
- Edit text and edit reply markup flow.
- Generated document caption edit flow.
- Delete message flow limited to bot-created test messages.
- Reply parameters flow.
- `SendChatAction` flow.
- Forward/copy flow.
- `GetChat` chat info flow.
- Long polling auto-exit flow after successful update and reply.

### Not included yet

`v0.1.0` is not full Telegram Bot API coverage. Notable deferred areas include:

- Remaining send methods: paid media.
- Invite links and join requests.
- Promote/admin management and chat metadata setters.
- Inline mode and `answerInlineQuery`.
- Payments, Passport, and Games.
- WebApp/LoginUrl button fields and helpers.
- Business APIs, Stars, gifts, and related recent Bot API areas.
- Full Bot API code generation/openapi tooling.

### Notes

- `v0.1.0` is an early API milestone, not a stable promise of full Telegram Bot API coverage.
- Some admin/destructive methods are implemented but intentionally not live-smoked automatically.
- Documentation contains API coverage, roadmap, manual testing, deploy testing, live smoke matrix, and release checklist.
- At the `v0.1.0` tag, the Go module path was `ai-gram`. Public import-path/module-path stabilization is prepared for the next patch release.
