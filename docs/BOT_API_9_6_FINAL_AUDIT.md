# Bot API 9.6 Final Coverage Audit

## Result

**Full coverage reached with documented architecture differences.**

The Stage 98 official-doc audit found that `ai-gram` has wrappers for all 169 official Telegram Bot API 9.6 methods and representative typed coverage for the major object/update/message families. Stage 99 resolved the remaining hard blocker by adding multipart certificate upload for the official `setWebhook.certificate` `InputFile` parameter.

A small decode gap found during the audit, `Message.giveaway`, was corrected in Stage 98. The remaining differences are documented architecture choices or naming/package differences.

## Source of truth

- [Official Telegram Bot API documentation](https://core.telegram.org/bots/api), fetched during Stage 98 and rechecked for `setWebhook` during Stage 99 on 2026-05-02.
- [Official Telegram Bot API changelog](https://core.telegram.org/bots/api-changelog), especially the April 3, 2026 Bot API 9.6 entry.

## Method coverage

Stage 98 compared every official method heading in the Telegram Bot API documentation with exported `(*bot.Bot)` method wrappers. Stage 99 rechecked the official `setWebhook` parameter table and certificate upload note.

Summary:

- Official Bot API methods audited: **169**.
- Implemented `(*bot.Bot)` wrappers: **169 / 169**.
- Missing method wrappers: **none known**.
- Needs verification / behavior blockers: **none known after Stage 99**.

Implemented method areas include:

- lifecycle and profile reads: `getMe`, `logOut`, `close`, profile photos/audios, forum topic icon stickers;
- updates and webhook management: `getUpdates`, `setWebhook` including upload-only certificate `InputFile`, `deleteWebhook`, `getWebhookInfo`;
- send/edit/delete/copy/forward/batch methods, including checklists, drafts, games, media groups, paid media, and live-location edits;
- chat/forum/admin/member/boost methods;
- regular and subscription invite links plus join requests;
- bot commands, menus, profile metadata, and default administrator rights;
- Business API foundation, account, story, suggested-post, gift, and business-message methods;
- gifts, Stars, paid media, invoices, payments, refunds, and subscription management;
- Managed Bots 9.6 methods and prepared inline/keyboard methods;
- inline mode, WebApp/Mini App, stickers, Passport, games, reactions, polls, verification, and user status methods.

## Type and field coverage

Stage 98 compared high-impact official object field tables with local JSON tags and package types.

Representative field comparison result after the Stage 98 `Message.giveaway` correction:

- `telegram.User`: no missing official fields found.
- `telegram.Chat`: no missing official fields found; `ChatFullInfo`-only fields remain in `telegram.ChatFullInfo`.
- `telegram.ChatFullInfo`: no missing official fields found in the audited official field table.
- `telegram.Update`: no missing official update fields found.
- `telegram.Message`: no missing official fields found in the audited official field table.
- `telegram.ReplyParameters`: no missing official fields found.
- `telegram.CallbackQuery`: no missing official fields found.
- `telegram.Video`: no missing official fields found.
- `telegram.Sticker` and `telegram.StickerSet`: no missing official fields found in the audited field tables.
- `telegram.InlineKeyboardButton` and `telegram.KeyboardButton`: no missing official fields found.

Official named types intentionally represented differently or needing a compatibility note:

- `MessageId` is represented as idiomatic `telegram.MessageID`.
- `ResponseParameters` lives in the public `errors` package and is re-exported from the root facade.
- `InputFile` is represented by `bot.FileRef` / `bot.FileUpload` rather than a direct `InputFile` type.
- `TransactionPartnerTelegramApi` is represented as `TransactionPartnerTelegramAPI` to preserve Go acronym style.
- Official concrete `ChatMemberOwner`, `ChatMemberAdministrator`, `ChatMemberMember`, `ChatMemberRestricted`, `ChatMemberLeft`, and `ChatMemberBanned` variant types are represented by the compatibility-preserving flat `telegram.ChatMember` struct. No official fields were found missing from the flat struct during this audit.

## Architecture differences

### `FileRef` / `FileUpload` vs official `InputFile`

The official Bot API uses `InputFile` for file uploads. `ai-gram` intentionally exposes `bot.FileRef` and `bot.FileUpload` so callers can pass Telegram `file_id`, HTTP(S) URLs, or multipart uploads through one Go type. This is acceptable as a public API design choice when every upload-capable field has JSON/multipart behavior.

Stage 98/99 verified the direct official method parameters containing `InputFile`:

- implemented JSON/multipart upload behavior: `sendPhoto.photo`, `sendAudio.audio`, `sendAudio.thumbnail`, `sendDocument.document`, `sendDocument.thumbnail`, `sendVideo.video`, `sendVideo.thumbnail`, `sendVideo.cover`, `sendAnimation.animation`, `sendAnimation.thumbnail`, `sendVoice.voice`, `sendVideoNote.video_note`, `sendVideoNote.thumbnail`, `setChatPhoto.photo`, `sendSticker.sticker`, `uploadStickerFile.sticker`, `setStickerSetThumbnail.thumbnail`, and upload-only `setWebhook.certificate`.

Polymorphic upload objects (`InputMedia*`, `InputPaidMedia*`, `InputProfilePhoto*`, `InputStoryContent*`, `InputSticker`) are implemented through typed Go structs and multipart helpers in their method families.

### `GetChat` vs `GetChatFullInfo`

Official `getChat` returns `ChatFullInfo`. `ai-gram` keeps the older `GetChat(ctx, GetChatParams) (*telegram.Chat, error)` for backward compatibility and adds `GetChatFullInfo(ctx, GetChatParams) (*telegram.ChatFullInfo, error)` as the full Bot API 9.6 result path. This split is intentional and acceptable before a deliberate breaking release.

### Flat `ChatMember` vs concrete `ChatMember*` variants

The official docs describe concrete `ChatMember*` variants. `ai-gram` keeps a flat `telegram.ChatMember` struct to avoid public API churn while decoding the union fields. Stage 98 did not find missing official fields in this flat struct. Concrete variants remain a future public API refinement, not a current decode blocker.

### `MaybeInaccessibleMessage` compatibility

`CallbackQuery.Message` remains available for accessible messages, while `CallbackQuery.MaybeMessage` preserves the official `MaybeInaccessibleMessage` shape and can decode inaccessible callback messages. This compatibility split is intentional.

### Webhook certificate upload

`SetWebhookParams` now exposes the official optional `certificate` parameter through `FileRef` and accepts only `FileUpload` for that upload-only official `InputFile` path. When `Certificate` is empty, `SetWebhook` keeps the existing JSON request path; when it is a `FileUpload`, `SetWebhook` sends multipart/form-data with the certificate part named `certificate`. File IDs and URLs are rejected because the official docs state that sending a string will not work for this parameter.

## Release-readiness blockers

No known Bot API 9.6 code coverage blockers remain after Stage 99. Stage 100 tracks release-readiness verification and manual-only smoke planning in [`docs/maintainer/BOT_API_9_6_RELEASE_READINESS.md`](maintainer/BOT_API_9_6_RELEASE_READINESS.md). Sensitive and state-changing live smoke remains manual-only.

Soft follow-up, not a release blocker if documented:

- optional concrete `ChatMember*` variant structs for users who want exact official union type names while preserving existing flat `ChatMember` compatibility.

## Manual-only live smoke areas

These areas must not be automatically live-smoked:

- payments, invoices, paid media, Stars, gifts, business gifts, subscription invite links, Premium subscription gifts, refunds, and subscription edits;
- Passport data and Passport error reporting;
- Managed bot token methods;
- Business APIs, business messages, business account profile changes, stories, suggested posts, direct messages, and business gifts;
- admin/destructive chat methods, including bans, restrictions, promotions, invite links, join requests, chat profile changes, leave chat, mass unpin, and sender-chat moderation;
- sticker set mutation methods;
- games requiring BotFather game setup;
- inline mode and prepared inline messages requiring BotFather/client setup;
- WebApp/Mini App flows;
- lifecycle `logOut` and `close`.

Decode/serialization-only areas such as reply metadata, service messages, direct-message metadata, `ChatFullInfo`, channel posts, and standalone poll updates should remain fixture-first unless a future manual check is explicitly requested.

## Recommended next step

Release-readiness documentation and manual-only smoke planning are complete. `main` was later published only after explicit user approval. The next publication phase is creating tags or GitHub Releases only after explicit maintainer approval.
