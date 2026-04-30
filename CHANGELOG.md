# Changelog

## Unreleased

### Added

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

### Documentation

- Completed inline mode audit against official Bot API 9.6 documentation.

### Planning

- Added v0.3 plan for chat management, forum topics, reactions, and inline mode basics.
- Added full Telegram Bot API 9.6 coverage plan and local-only push/tag/release freeze policy.

## v0.2.0 - 2026-04-30

### Added

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
