# API Coverage

This document maps the current `ai-gram` implementation to Telegram Bot API areas. It is a project inventory, not a generated copy of the full upstream Bot API specification. Telegram adds methods over time, so expansion work should still be checked against the official Bot API docs before implementation.

> **Bot API 9.6 target:** The current coverage target is 100% Telegram Bot API 9.6. Track the local-only full coverage workstream in [`docs/BOT_API_9_6_COVERAGE_PLAN.md`](BOT_API_9_6_COVERAGE_PLAN.md). Not all Bot API 9.6 areas are implemented yet. Push, tag, and GitHub Release suggestions are frozen until the full 9.6 plan is complete.

## Implemented

### Core

| Public Go API | Telegram Bot API method | Tests | Notes |
| --- | --- | --- | --- |
| `aigram.New`, `aigram.NewBot`, `bot.New` | n/a | unit | Token validation, configurable base URL and HTTP client. Token is stored privately, not exposed by public accessors, and redacted from string output. |
| `(*bot.Bot).GetMe` | `getMe` | unit/httptest, live via smoke scripts | Basic identity check used by discovery and smoke helpers. |
| `errors.APIError`, `errors.ResponseParameters` | Bot API error envelope | unit | `ok:false` responses return typed errors; tests cover `errors.As`. |
| `bot.ChatID`, `ChatIDInt`, `ChatIDString` | `chat_id` parameter shape | unit | Supports numeric chat IDs and string IDs such as `@channelusername`. |

### Updates

| Public Go API | Telegram Bot API method | Tests | Notes |
| --- | --- | --- | --- |
| `(*bot.Bot).GetUpdates` | `getUpdates` | unit/httptest | Manual one-shot updates call. |
| `transport/longpoll.Runner` | `getUpdates` loop | unit, live via examples/scripts | Managed offset advancement, backoff, context cancellation, handler error reporting. |
| `telegram.Update`, `telegram.Message`, helpers | n/a | unit | Practical incoming update/message/callback/media decoding and helper methods. |
| `dispatch.Dispatcher` | n/a | unit, live via examples | Predicate routing for messages, commands, callbacks, middleware, fallback, error handling. |

### Webhook

| Public Go API | Telegram Bot API method | Tests | Notes |
| --- | --- | --- | --- |
| `(*bot.Bot).SetWebhook` | `setWebhook` | unit/httptest, live via deploy harness | JSON-only webhook registration with URL and secret token. Certificate upload is not implemented. |
| `(*bot.Bot).DeleteWebhook` | `deleteWebhook` | unit/httptest, manual/live harness | Supports `drop_pending_updates`; destructive use should be explicit. |
| `(*bot.Bot).GetWebhookInfo` | `getWebhookInfo` | unit/httptest, smoke scripts | Used for troubleshooting and local Bot API checks. |
| `transport/webhook.New` | inbound webhook handler | unit, live via deploy harness | Validates method, content type, optional secret token, JSON body, and handler errors. |

### Send methods

| Public Go API | Telegram Bot API method | Tests | Notes |
| --- | --- | --- | --- |
| `(*bot.Bot).SendMessage` | `sendMessage` | unit/httptest, live examples | Supports text, parse mode/entities conflict validation, reply markup, `message_thread_id`, `reply_parameters`. |
| `(*bot.Bot).SendPhoto` | `sendPhoto` | unit/httptest, live examples | Supports `FileID`, `FileURL`, `FileUpload`, caption, reply markup, thread/reply params. |
| `(*bot.Bot).SendDocument` | `sendDocument` | unit/httptest, live examples | Supports `FileID`, `FileURL`, `FileUpload`, caption, reply markup, thread/reply params. |
| `(*bot.Bot).SendVideo` | `sendVideo` | unit/httptest | Supports `FileID`, `FileURL`, `FileUpload`, caption, duration, dimensions, streaming, thread/reply params. |
| `(*bot.Bot).SendAudio` | `sendAudio` | unit/httptest | Supports `FileID`, `FileURL`, `FileUpload`, caption, duration, performer/title, thread/reply params. |
| `(*bot.Bot).SendVoice` | `sendVoice` | unit/httptest | Supports `FileID`, `FileURL`, `FileUpload`, caption, duration, thread/reply params. |
| `(*bot.Bot).SendContact` | `sendContact` | unit/httptest, live v0.2 smoke | Supports contact phone/name/vCard fields, reply markup, `message_thread_id`, `reply_parameters`. |
| `(*bot.Bot).SendLocation` | `sendLocation` | unit/httptest, live v0.2 smoke | Supports latitude/longitude, live-location optional fields, reply markup, thread/reply params. |
| `(*bot.Bot).SendVenue` | `sendVenue` | unit/httptest, live v0.2 smoke | Supports venue coordinates, title/address, Foursquare/Google place fields, reply markup, thread/reply params. |
| `(*bot.Bot).SendPoll` | `sendPoll` | unit/httptest, live v0.2 smoke | Supports question/options, quiz fields, explanation formatting, reply markup, thread/reply params. |
| `(*bot.Bot).SendDice` | `sendDice` | unit/httptest, live v0.2 smoke | Supports known Telegram dice emoji, reply markup, thread/reply params. |
| `(*bot.Bot).SendSticker` | `sendSticker` | unit/httptest, optional live v0.2 smoke | Supports `FileID`, `FileURL`, `FileUpload`, emoji, reply markup, thread/reply params. |
| `(*bot.Bot).SendAnimation` | `sendAnimation` | unit/httptest, optional live v0.2 smoke | Supports `FileID`, `FileURL`, `FileUpload`, caption fields, thumbnail file ref/upload, spoiler, reply markup, thread/reply params. |
| `(*bot.Bot).SendVideoNote` | `sendVideoNote` | unit/httptest, optional live v0.2 smoke | Supports `FileID`, `FileUpload`, thumbnail file ref/upload, duration/length, reply markup, thread/reply params. HTTP URL is intentionally rejected for video notes. |
| `(*bot.Bot).SendMediaGroup` | `sendMediaGroup` | unit/httptest, live generated-upload smoke | Supports `InputMediaPhoto`, `InputMediaVideo`, `InputMediaAudio`, `InputMediaDocument`, JSON file IDs/URLs, multipart uploads, thumbnail uploads, thread/reply params. Does not support reply markup because Telegram does not accept it for media groups. |
| `telegram.ReplyParameters` | send/copy reply payload | unit | Minimal supported fields: `message_id`, `allow_sending_without_reply`. |
| `telegram.ReplyMarkup` implementations | send/edit reply markup | unit, live examples | Inline keyboard, reply keyboard, remove keyboard, force reply. Edit methods accept inline keyboard only. |

### Media/files

| Public Go API | Telegram Bot API method | Tests | Notes |
| --- | --- | --- | --- |
| `bot.FileID`, `bot.FileURL`, `bot.FileUpload` | file reference parameters | unit/httptest, live examples | Uploads use multipart `attach://`; callers own reader lifecycle. |
| `(*bot.Bot).GetFile` | `getFile` | unit/httptest, live media script | Gets `file_path` for later download. |
| `(*bot.Bot).DownloadFile` | file download endpoint | unit/httptest, live media script | Streams to caller-provided writer and does not expose token-bearing download URLs. |
| multipart helpers | n/a | unit/httptest | Covers media uploads and JSON string fields such as reply parameters. |


### Sticker set management

| Public Go API | Telegram Bot API method / object | Tests | Notes |
| --- | --- | --- | --- |
| `(*bot.Bot).GetStickerSet` | `getStickerSet` | unit/httptest | Decodes `telegram.StickerSet` with minimal sticker metadata. |
| `(*bot.Bot).GetCustomEmojiStickers` | `getCustomEmojiStickers` | unit/httptest | Returns custom emoji `[]telegram.Sticker` by ID list. |
| `(*bot.Bot).UploadStickerFile` | `uploadStickerFile` | unit/httptest | Multipart-only upload for sticker set creation workflows; callers own reader lifecycle. |
| `(*bot.Bot).CreateNewStickerSet` | `createNewStickerSet` | unit/httptest | Supports `InputSticker` JSON mode and multipart upload mode with deterministic `attach://` names. Manual-only live smoke. |
| `(*bot.Bot).AddStickerToSet` | `addStickerToSet` | unit/httptest | Supports `InputSticker` JSON/multipart payloads. Manual-only live smoke. |
| `(*bot.Bot).ReplaceStickerInSet` | `replaceStickerInSet` | unit/httptest | Supports replacing a sticker with an `InputSticker` JSON/multipart payload. Manual-only live smoke. |
| `(*bot.Bot).SetStickerPositionInSet` | `setStickerPositionInSet` | unit/httptest | Moves a sticker within its set. Manual-only live smoke. |
| `(*bot.Bot).DeleteStickerFromSet` | `deleteStickerFromSet` | unit/httptest | Deletes a real sticker from a set. Manual-only live smoke. |
| `(*bot.Bot).SetStickerEmojiList` | `setStickerEmojiList` | unit/httptest | Sets the emoji list for a sticker. Manual-only live smoke. |
| `(*bot.Bot).SetStickerKeywords` | `setStickerKeywords` | unit/httptest | Sets or clears sticker search keywords. Manual-only live smoke. |
| `(*bot.Bot).SetStickerMaskPosition` | `setStickerMaskPosition` | unit/httptest | Sets or clears mask sticker position. Manual-only live smoke. |
| `(*bot.Bot).SetStickerSetTitle` | `setStickerSetTitle` | unit/httptest | Changes a real sticker set title. Manual-only live smoke. |
| `(*bot.Bot).SetStickerSetThumbnail` | `setStickerSetThumbnail` | unit/httptest | Supports optional thumbnail removal, file IDs, and multipart upload; animated/video thumbnail URLs are rejected. Manual-only live smoke. |
| `(*bot.Bot).SetCustomEmojiStickerSetThumbnail` | `setCustomEmojiStickerSetThumbnail` | unit/httptest | Sets or clears a custom emoji sticker set thumbnail. Manual-only live smoke. |
| `(*bot.Bot).DeleteStickerSet` | `deleteStickerSet` | unit/httptest | Deletes a real sticker set created by the bot. Manual-only live smoke. |
| `bot.InputSticker`, `telegram.StickerSet`, `telegram.MaskPosition` | related Bot API objects | unit through method payload/result tests | Minimal typed coverage for sticker set management and custom emoji sticker workflows. |

### Callback/edit/delete

| Public Go API | Telegram Bot API method | Tests | Notes |
| --- | --- | --- | --- |
| `(*bot.Bot).AnswerCallbackQuery` | `answerCallbackQuery` | unit/httptest, live examples | Supports toast/alert and URL/cache fields. |
| `(*bot.Bot).EditMessageText` | `editMessageText` | unit/httptest, live examples | Supports chat and inline targets; result decodes `Message` or `true`. |
| `(*bot.Bot).EditMessageCaption` | `editMessageCaption` | unit/httptest, live examples | Supports empty caption removal and inline keyboard. |
| `(*bot.Bot).EditMessageReplyMarkup` | `editMessageReplyMarkup` | unit/httptest, live examples | `nil` reply markup removes inline keyboard. |
| `bot.EditMessageTarget`, `bot.EditMessageResult` | edit helpers/result | unit | Validates chat-vs-inline target and handles `Message`/`true` return shape. |
| `(*bot.Bot).DeleteMessage` | `deleteMessage` | unit/httptest, live examples | Destructive; live example only deletes messages created during smoke. |
| `(*bot.Bot).DeleteMessages` | `deleteMessages` | unit/httptest | Destructive batch delete for 1-100 message IDs; manual-only live smoke. |
| `(*bot.Bot).StopPoll` | `stopPoll` | unit/httptest, live v0.2 smoke | Stops a poll sent by the bot and returns `telegram.Poll`. |
| `(*bot.Bot).AnswerInlineQuery` | `answerInlineQuery` | unit/httptest | Supports article results, text input message content, cache/pagination fields, and inline results button. Manual-only live smoke. |
| `telegram.InlineQuery` | `inline_query` update | unit | Incoming inline query decoding and `EffectiveUser` support. |
| `telegram.ChosenInlineResult` | `chosen_inline_result` update | unit | Incoming chosen inline result decoding and `EffectiveUser` support. |
| `dispatch.InlineQuery`, `dispatch.ChosenInlineResult` | dispatch predicates/helpers | unit | Includes `OnInlineQuery` and `OnChosenInlineResult` handler registration helpers. |
| `bot.InlineQueryResultArticle`, `bot.InputTextMessageContent` | inline payload objects | unit | Minimal inline result/content foundation; remaining result variants are pending. |

### Forward/copy

| Public Go API | Telegram Bot API method | Tests | Notes |
| --- | --- | --- | --- |
| `(*bot.Bot).ForwardMessage` | `forwardMessage` | unit/httptest, live examples | Supports thread ID, disable notification, protect content. |
| `(*bot.Bot).ForwardMessages` | `forwardMessages` | unit/httptest | Batch forwards 1-100 message IDs and returns `[]telegram.MessageID`. Not auto-smoked. |
| `(*bot.Bot).CopyMessage` | `copyMessage` | unit/httptest, live examples | Returns `telegram.MessageID`; supports caption, reply parameters, reply markup, notification/protect flags. |
| `(*bot.Bot).CopyMessages` | `copyMessages` | unit/httptest | Batch copies 1-100 message IDs and returns `[]telegram.MessageID`; supports remove caption and notification/protect flags. Not auto-smoked. |

### Chat actions

| Public Go API | Telegram Bot API method | Tests | Notes |
| --- | --- | --- | --- |
| `(*bot.Bot).SendChatAction` | `sendChatAction` | unit/httptest, live examples | Validates known action constants; echo handler uses `typing` in live smoke. |
| `(*bot.Bot).PinChatMessage` | `pinChatMessage` | unit/httptest | Admin-required; not part of default live smoke. |
| `(*bot.Bot).UnpinChatMessage` | `unpinChatMessage` | unit/httptest | Admin-required; `message_id` optional. |
| `(*bot.Bot).UnpinAllChatMessages` | `unpinAllChatMessages` | unit/httptest | Admin/destructive; not part of default live smoke. |

### Chat/member info

| Public Go API | Telegram Bot API method | Tests | Notes |
| --- | --- | --- | --- |
| `(*bot.Bot).GetChat` | `getChat` | unit/httptest, live example access panel | Minimal `telegram.Chat` fields plus optional description/invite/pinned message. |
| `(*bot.Bot).GetChatMember` | `getChatMember` | unit/httptest | Minimal `telegram.ChatMember` fields and admin permission booleans. |
| `(*bot.Bot).GetChatAdministrators` | `getChatAdministrators` | unit/httptest | Returns `[]telegram.ChatMember`. |
| `(*bot.Bot).GetChatMemberCount` | `getChatMemberCount` | unit/httptest, optional live example | Safe read method; availability depends on chat permissions. |

### Moderation

| Public Go API | Telegram Bot API method | Tests | Notes |
| --- | --- | --- | --- |
| `(*bot.Bot).BanChatMember` | `banChatMember` | unit/httptest | Destructive/admin method; no automatic live smoke. |
| `(*bot.Bot).UnbanChatMember` | `unbanChatMember` | unit/httptest | Admin method; `OnlyIfBanned` supported. |
| `(*bot.Bot).RestrictChatMember` | `restrictChatMember` | unit/httptest | Destructive/admin method; zero `telegram.ChatPermissions` is valid and restricts all supported actions. |
| `telegram.ChatPermissions` | moderation permissions object | unit through method payload tests | Minimal supported permission fields for restriction and default chat permission payloads. |

### Chat management

| Public Go API | Telegram Bot API method | Tests | Notes |
| --- | --- | --- | --- |
| `(*bot.Bot).SetChatTitle` | `setChatTitle` | unit/httptest | Admin/state-changing method for chat titles. Not auto-smoked. |
| `(*bot.Bot).SetChatDescription` | `setChatDescription` | unit/httptest | Admin/state-changing method for chat descriptions; empty description is allowed to remove it. Not auto-smoked. |
| `(*bot.Bot).SetChatPhoto` | `setChatPhoto` | unit/httptest | Multipart-only upload method; `FileID` and `FileURL` are intentionally rejected. Not auto-smoked. |
| `(*bot.Bot).DeleteChatPhoto` | `deleteChatPhoto` | unit/httptest | Admin/state-changing method for chat photos. Not auto-smoked. |
| `(*bot.Bot).LeaveChat` | `leaveChat` | unit/httptest | Makes the bot leave a chat; manual-only and disposable-chat testing recommended. |
| `(*bot.Bot).SetChatStickerSet` | `setChatStickerSet` | unit/httptest | Admin/state-changing supergroup sticker-set method. Not auto-smoked. |
| `(*bot.Bot).DeleteChatStickerSet` | `deleteChatStickerSet` | unit/httptest | Admin/state-changing supergroup sticker-set method. Not auto-smoked. |

### Forum topics

| Public Go API | Telegram Bot API method | Tests | Notes |
| --- | --- | --- | --- |
| `(*bot.Bot).CreateForumTopic` | `createForumTopic` | unit/httptest | Admin/state-changing forum supergroup method; returns `telegram.ForumTopic`. Not auto-smoked. |
| `(*bot.Bot).EditForumTopic` | `editForumTopic` | unit/httptest | Edits a real forum topic name or icon. Not auto-smoked. |
| `(*bot.Bot).CloseForumTopic` | `closeForumTopic` | unit/httptest | Closes a real forum topic. Not auto-smoked. |
| `(*bot.Bot).ReopenForumTopic` | `reopenForumTopic` | unit/httptest | Reopens a real forum topic. Not auto-smoked. |
| `(*bot.Bot).DeleteForumTopic` | `deleteForumTopic` | unit/httptest | Deletes a real forum topic. Not auto-smoked. |
| `(*bot.Bot).UnpinAllForumTopicMessages` | `unpinAllForumTopicMessages` | unit/httptest | Clears pinned messages in a real forum topic. Not auto-smoked. |
| `(*bot.Bot).EditGeneralForumTopic` | `editGeneralForumTopic` | unit/httptest | Edits the General forum topic name. Not auto-smoked. |
| `(*bot.Bot).CloseGeneralForumTopic` | `closeGeneralForumTopic` | unit/httptest | Closes the General forum topic. Not auto-smoked. |
| `(*bot.Bot).ReopenGeneralForumTopic` | `reopenGeneralForumTopic` | unit/httptest | Reopens the General forum topic. Not auto-smoked. |
| `(*bot.Bot).HideGeneralForumTopic` | `hideGeneralForumTopic` | unit/httptest | Hides the General forum topic. Not auto-smoked. |
| `(*bot.Bot).UnhideGeneralForumTopic` | `unhideGeneralForumTopic` | unit/httptest | Unhides the General forum topic. Not auto-smoked. |
| `(*bot.Bot).UnpinAllGeneralForumTopicMessages` | `unpinAllGeneralForumTopicMessages` | unit/httptest | Clears pinned messages in the General forum topic. Not auto-smoked. |
| `telegram.ForumTopic` and forum topic service message types | related Bot API objects | unit | Minimal topic result and service message decoding for created, edited, closed, reopened, hidden, and unhidden topic events. |

### Reactions

| Public Go API | Telegram Bot API method / object | Tests | Notes |
| --- | --- | --- | --- |
| `(*bot.Bot).SetMessageReaction` | `setMessageReaction` | unit/httptest | Changes real message reaction state. Manual-only live smoke. |
| `telegram.ReactionTypeEmoji` | `ReactionTypeEmoji` | unit | Polymorphic reaction marshal/unmarshal with required `type: "emoji"`. |
| `telegram.ReactionTypeCustomEmoji` | `ReactionTypeCustomEmoji` | unit | Polymorphic reaction marshal/unmarshal with required `type: "custom_emoji"`. |
| `telegram.ReactionTypePaid` | `ReactionTypePaid` | unit | Polymorphic paid reaction support from Bot API 9.6. |
| `telegram.MessageReactionUpdated`, `telegram.Update.MessageReaction` | `message_reaction` update | unit | Decodes old/new reaction lists and supports `EffectiveChat`/`EffectiveUser`. |
| `telegram.MessageReactionCountUpdated`, `telegram.Update.MessageReactionCount` | `message_reaction_count` update | unit | Decodes anonymous reaction count updates and supports `EffectiveChat`. |
| `dispatch.MessageReaction`, `dispatch.MessageReactionCount` | n/a | unit | Predicate and route helpers for reaction updates. |

### Admin management

| Public Go API | Telegram Bot API method | Tests | Notes |
| --- | --- | --- | --- |
| `(*bot.Bot).PromoteChatMember` | `promoteChatMember` | unit/httptest | Admin/state-changing method for promoting or demoting users. Not auto-smoked. |
| `(*bot.Bot).SetChatAdministratorCustomTitle` | `setChatAdministratorCustomTitle` | unit/httptest | Admin/state-changing method for custom administrator titles. Not auto-smoked. |
| `(*bot.Bot).SetChatPermissions` | `setChatPermissions` | unit/httptest | Admin/state-changing method for default chat permissions; zero `telegram.ChatPermissions` is valid. Not auto-smoked. |

### Bot commands/menu

| Public Go API | Telegram Bot API method | Tests | Notes |
| --- | --- | --- | --- |
| `(*bot.Bot).SetMyCommands` | `setMyCommands` | unit/httptest | Supports command lists, scope objects, and language code. Changes bot-level command state, so no automatic live smoke. |
| `(*bot.Bot).DeleteMyCommands` | `deleteMyCommands` | unit/httptest | Deletes commands for a scope/language. Changes bot-level command state, so no automatic live smoke. |
| `(*bot.Bot).GetMyCommands` | `getMyCommands` | unit/httptest | Decodes command lists for a scope/language. |
| `(*bot.Bot).SetChatMenuButton` | `setChatMenuButton` | unit/httptest | Supports commands/default/web_app menu buttons; changes menu state, so no automatic live smoke. |
| `(*bot.Bot).GetChatMenuButton` | `getChatMenuButton` | unit/httptest | Decodes polymorphic commands/default/web_app menu buttons. |
| `(*bot.Bot).SetMyDefaultAdministratorRights` | `setMyDefaultAdministratorRights` | unit/httptest | Sets or clears default administrator rights requested by the bot; no automatic live smoke. |
| `telegram.BotCommandScope`, `telegram.MenuButton`, `telegram.ChatAdministratorRights` | related Bot API objects | unit through method payload tests | Hand-written minimal object coverage for command scopes, menu buttons, Web App info, and admin rights. |

### Bot profile and metadata

| Public Go API | Telegram Bot API method | Tests | Notes |
| --- | --- | --- | --- |
| `(*bot.Bot).SetMyName` | `setMyName` | unit/httptest | Changes localized bot name metadata; manual-only live smoke. |
| `(*bot.Bot).GetMyName` | `getMyName` | unit/httptest | Decodes `telegram.BotName`. |
| `(*bot.Bot).SetMyDescription` | `setMyDescription` | unit/httptest | Changes localized bot description metadata; manual-only live smoke. |
| `(*bot.Bot).GetMyDescription` | `getMyDescription` | unit/httptest | Decodes `telegram.BotDescription`. |
| `(*bot.Bot).SetMyShortDescription` | `setMyShortDescription` | unit/httptest | Changes localized bot short description metadata; manual-only live smoke. |
| `(*bot.Bot).GetMyShortDescription` | `getMyShortDescription` | unit/httptest | Decodes `telegram.BotShortDescription`. |
| `(*bot.Bot).GetMyDefaultAdministratorRights` | `getMyDefaultAdministratorRights` | unit/httptest | Reads default administrator rights requested by the bot. |
| `(*bot.Bot).SetMyProfilePhoto` | `setMyProfilePhoto` | unit/httptest | Supports Bot API 9.6 `InputProfilePhotoStatic` and `InputProfilePhotoAnimated`; upload-only multipart because profile photos cannot be reused. Manual-only live smoke. |
| `(*bot.Bot).RemoveMyProfilePhoto` | `removeMyProfilePhoto` | unit/httptest | Removes the bot profile photo; manual-only live smoke. |
| `telegram.BotName`, `telegram.BotDescription`, `telegram.BotShortDescription` | related Bot API objects | unit through method result tests | Minimal localized bot profile metadata objects. |

### Invite links

| Public Go API | Telegram Bot API method | Tests | Notes |
| --- | --- | --- | --- |
| `(*bot.Bot).ExportChatInviteLink` | `exportChatInviteLink` | unit/httptest | Admin-required method that generates a new primary invite link and revokes the previous primary link. Not auto-smoked. |
| `(*bot.Bot).CreateChatInviteLink` | `createChatInviteLink` | unit/httptest | Creates a real additional invite link; supports name, expire date, member limit, and join-request flag. Not auto-smoked. |
| `(*bot.Bot).EditChatInviteLink` | `editChatInviteLink` | unit/httptest | Edits a non-primary invite link created by the bot. Not auto-smoked. |
| `(*bot.Bot).RevokeChatInviteLink` | `revokeChatInviteLink` | unit/httptest | Revokes a real invite link created by the bot. Not auto-smoked. |
| `telegram.ChatInviteLink` | invite link object | unit through method result tests | Minimal invite link metadata with creator, primary/revoked flags, limits, expiry, and pending join request count. |

### Chat join requests

| Public Go API | Telegram Bot API method | Tests | Notes |
| --- | --- | --- | --- |
| `(*bot.Bot).ApproveChatJoinRequest` | `approveChatJoinRequest` | unit/httptest | Admin-required method that approves a real pending join request. Not auto-smoked. |
| `(*bot.Bot).DeclineChatJoinRequest` | `declineChatJoinRequest` | unit/httptest | Admin-required method that declines a real pending join request. Not auto-smoked. |
| `telegram.ChatJoinRequest`, `telegram.Update.ChatJoinRequest` | `chat_join_request` update | unit | Decodes join request updates and supports `EffectiveChat`/`EffectiveUser`. |
| `dispatch.ChatJoinRequest`, `(*dispatch.Dispatcher).OnChatJoinRequest` | n/a | unit | Predicate and route helper for join request updates. |

### Access/example infrastructure

| Public Go API / artifact | Tests | Notes |
| --- | --- | --- |
| `middleware.Access`, `AccessWithPolicy`, `AccessConfig` | unit, live examples | Admin/public/off access control for dispatcher handlers without importing `bot`. |
| `middleware.Recover`, `Timeout`, `Observe` | unit | Handler safety and instrumentation hooks. |
| `examples/echo_longpoll` | compile via `go test ./...`, manual smoke | Basic long polling echo. |
| `examples/inline_longpoll` | compile via `go test ./...`, manual smoke | Inline callbacks, edit flow, access commands. |
| `examples/webhook_server` | compile via `go test ./...`, live deploy smoke | Webhook, access panel, safe logs, callback/edit/delete/copy/forward/chat-info flows. |
| `examples/media_upload` | compile via `go test ./...`, manual smoke | Upload/download smoke without committing tokens. |
| `examples/media_group_smoke` | compile via `go test ./...`, live SendMediaGroup smoke | Self-contained generated upload fallback plus optional FileID/path modes without printing full file IDs or sensitive paths. |
| `examples/local_api_server` | compile via `go test ./...`, smoke scripts | Local Telegram Bot API server checks. |
| `scripts/*.sh`, `deploy/systemd/*.tmpl` | `bash -n`, live/manual smoke | Discovery, auto SSH tunnel, deploy, logs, stop, notifications, separate Bot API host support. |
| `docs/MANUAL_TESTING.md`, `docs/DEPLOY_TESTING.md` | review/manual | Manual smoke, deploy harness, TUN/Xray caveats, security notes. |

## Not implemented yet

### Remaining send methods

- `sendPaidMedia`
- `sendGame`
- `sendInvoice`

### Invite links

- subscription invite link methods

### Join requests

- subscription-specific join request flows beyond basic approval/decline

### Payments

- invoice, shipping query, pre-checkout query, refund and paid media/star payment flows

### Passport

- Telegram Passport data types and error methods

### Games

- game sending and score methods

### Inline mode

- Remaining `InlineQueryResult` variants beyond `InlineQueryResultArticle`.
- Remaining `InputMessageContent` variants beyond `InputTextMessageContent`.
- Inline result media/file variants and cached variants.

### WebApp/LoginUrl

- WebApp/LoginUrl fields outside the implemented menu button support
- LoginUrl button fields and validation
- web app data helpers

### Business features

- business connection/message types and methods
- business intro/location/hours/account features

### Stars/gifts if applicable

- Stars transaction/gift methods and related update types should be reviewed against the current official Bot API before planning.

### Codegen/openapi tooling if applicable

- No code generation is in use. Full Bot API codegen/openapi tooling is intentionally deferred until the hand-written public API shape stabilizes.

## Risk classification

### Safe read-only

- `GetMe`
- `GetFile`
- `GetWebhookInfo`
- `GetChat`
- `GetChatMember`
- `GetChatAdministrators`
- `GetChatMemberCount`
- `DownloadFile` when the file path comes from `GetFile` and the destination is controlled by the caller

### Safe send

- `SendMessage`
- `SendPhoto`
- `SendDocument`
- `SendVideo`
- `SendAudio`
- `SendVoice`
- `SendContact`
- `SendLocation`
- `SendVenue`
- `SendPoll`
- `StopPoll`
- `SendDice`
- `SendSticker`
- `SendAnimation`
- `SendVideoNote`
- `SendMediaGroup`
- `AnswerCallbackQuery`
- `EditMessageText`
- `EditMessageCaption`
- `EditMessageReplyMarkup`
- `ForwardMessage`
- `ForwardMessages`
- `CopyMessage`
- `CopyMessages`
- `SendChatAction`

These still require real credentials and may notify users, but they are not destructive when used in dedicated test chats.

### Admin required

- `PinChatMessage`
- `UnpinChatMessage`
- `UnpinAllChatMessages`
- `BanChatMember`
- `UnbanChatMember`
- `RestrictChatMember`
- `PromoteChatMember`
- `SetChatAdministratorCustomTitle`
- `SetChatPermissions`
- some chat/member info methods depending on chat type and bot permissions
- forum topic methods and service message types
- chat invite link methods (`ExportChatInviteLink`, `CreateChatInviteLink`, `EditChatInviteLink`, `RevokeChatInviteLink`)
- chat join request methods (`ApproveChatJoinRequest`, `DeclineChatJoinRequest`)
- reaction methods (`SetMessageReaction`) when used outside isolated test messages, because they change real message reaction state

### Destructive

- `DeleteMessage`
- `DeleteMessages`
- `DeleteWebhook` when `drop_pending_updates=true`
- `UnpinAllChatMessages`
- `BanChatMember`
- `UnbanChatMember`
- `RestrictChatMember`
- chat management methods (`SetChatTitle`, `SetChatDescription`, `SetChatPhoto`, `DeleteChatPhoto`, `LeaveChat`, `SetChatStickerSet`, `DeleteChatStickerSet`) when used outside isolated test chats, because they change real chat state
- forum topic methods when used outside isolated test forum supergroups, because they create, edit, close, reopen, delete, hide, or unpin real forum topic state
- admin management methods (`PromoteChatMember`, `SetChatAdministratorCustomTitle`, `SetChatPermissions`) when used outside isolated test chats, because they change real chat/admin state
- invite link methods when used outside isolated test chats, because they create or revoke real access links
- chat join request methods, because they approve or decline real users waiting to join
- reaction methods when used outside isolated test messages, because they change real message reaction state
- batch delete methods when used outside disposable test messages
- future moderation/admin methods

### Requires upload/multipart

- `SendPhoto` with `FileUpload`
- `SendDocument` with `FileUpload`
- `SendVideo` with `FileUpload`
- `SendAudio` with `FileUpload`
- `SendVoice` with `FileUpload`
- `SendSticker`, `SendAnimation`, and `SendVideoNote` with `FileUpload`
- `SendMediaGroup` with media or thumbnail `FileUpload`
- `SetChatPhoto` with `FileUpload`
- future upload methods such as remaining thumbnails, certificates

### Requires live credentials

- All real Bot API calls against Telegram or a local Telegram Bot API server
- long polling and webhook delivery
- file download from Telegram
- deploy/smoke scripts

Unit and httptest suites do not require tokens.

### Should not be smoke-tested automatically

- `BanChatMember`
- `UnbanChatMember`
- `RestrictChatMember`
- `PromoteChatMember`, `SetChatAdministratorCustomTitle`, and `SetChatPermissions`
- `PinChatMessage`, `UnpinChatMessage`, `UnpinAllChatMessages` outside a dedicated test group
- `DeleteWebhook` with `drop_pending_updates=true`
- bot commands/menu setters because they change bot-level command/menu state
- invite link and chat join request methods
- future migration methods such as `logOut`/`close`
- future payment/passport/gift methods

## v0.1 recommendation

### Must-have for v0.1

- Stable core `Bot` construction, token handling, base URL and HTTP client configuration.
- Typed `APIError` and consistent token-safe errors.
- Updates, webhook receiver, webhook management, long polling runner.
- Dispatcher/router and essential middleware including recovery, timeout, observability, and access control.
- Send text and implemented media methods with reply markup, reply parameters, and thread IDs.
- Callback query, edit text/caption/reply markup, delete message.
- File upload/download support for implemented media methods.
- Forward/copy and basic chat action/chat info support.
- Examples for long polling, inline callbacks, webhook, local Bot API, media, deploy harness.
- Admin-only examples and safe logs for live smoke.
- Documentation: README, manual testing, deploy testing, API coverage, roadmap, live smoke matrix, release checklist, security notes.
- Release checklist with no-token/no-secret verification and explicit no-auto-smoke rules for destructive/admin methods.

### Nice-to-have before v0.1

- Remaining high-risk/advanced Bot API coverage such as payments, business APIs, and advanced inline variants.
- Bot command and menu methods.
- A small release checklist document if not folded into existing docs.
- README tightening to avoid overpromising unimplemented Bot API areas.

### Defer after v0.1

- Payments and paid media.
- Passport.
- Games.
- Business APIs.
- Stars/gifts.
- Full Bot API codegen or openapi tooling.
- Broad admin/promote/forum management surface beyond methods already implemented.
