# Bot API 10.0 Lightweight Audit - 2026-05-14

This document records a lightweight current-state audit against the official Telegram Bot API documentation after adding the maintainer update checklist.

## Source

- Official source checked: <https://core.telegram.org/bots/api>
- Page state observed: Recent changes still list Bot API 10.0 dated 2026-05-08 as the latest Bot API release.
- Local baseline: [`docs/API_COVERAGE.md`](API_COVERAGE.md) and [`docs/BOT_API_10_0_FINAL_AUDIT.md`](BOT_API_10_0_FINAL_AUDIT.md).
- Scope: lightweight freshness check, not a new full field-by-field final audit.

## Result

No new Bot API version newer than 10.0 was observed. No missing official Bot API method wrapper was found in the local `bot` package during the method-level comparison.

No code changes are required from this audit.

## Method-Level Check

The official page method headings were compared with public `(*bot.Bot)` methods.

Observed result:

- official Bot API methods extracted: 176;
- local public `(*bot.Bot)` methods extracted: 180;
- missing official methods in local wrappers: none;
- local-only methods: `DownloadFile`, `GetChatFullInfo`, `GoString`, and `String`.

The local-only methods are expected:

- `DownloadFile` is an ai-gram helper for file download URLs.
- `GetChatFullInfo` is the local typed split for the full chat result shape.
- `String` and `GoString` are diagnostic/token-redaction helpers, not Bot API methods.

## Bot API 10.0 Change Groups

### Guest Mode

Official 10.0 guest-mode changes remain represented locally:

- `User.supports_guest_queries`;
- `Message.guest_bot_caller_user`;
- `Message.guest_bot_caller_chat`;
- `Message.guest_query_id`;
- `Update.guest_message`;
- `SentGuestMessage`;
- `answerGuestQuery`;
- guest-message dispatch routing.

Coverage references:

- `docs/API_COVERAGE.md` guest-mode rows;
- `docs/BOT_API_10_0_FINAL_AUDIT.md` guest-mode section;
- `bot/guest_methods.go`;
- `dispatch.OnGuestMessage`;
- `telegram.Update`, `telegram.Message`, and `telegram.User`.

### Chat Management

Official 10.0 chat-management changes remain represented locally:

- `ChatMemberRestricted.can_react_to_messages`;
- `ChatPermissions.can_react_to_messages`;
- `getChatAdministrators.return_bots`;
- `deleteMessageReaction`;
- `deleteAllMessageReactions`.

Coverage references:

- `docs/API_COVERAGE.md` chat management and reaction rows;
- `docs/BOT_API_10_0_FINAL_AUDIT.md` chat-management section;
- `bot/chat_info_methods.go`;
- `bot/reaction_methods.go`;
- `telegram.ChatPermissions` and chat member variants.

### Polls

Official 10.0 poll changes remain represented locally:

- incoming `PollMedia`, poll media fields, one-option polls, `members_only`, and `country_codes`;
- outgoing `InputPollMedia`, `InputPollOptionMedia`, poll media params, `members_only`, and `country_codes`;
- input media variants for animation, audio, document, live photo, location, photo, sticker, venue, and video where locally supported.

Coverage references:

- `docs/API_COVERAGE.md` poll rows;
- `docs/BOT_API_10_0_FINAL_AUDIT.md` poll section;
- `bot/poll_dice_methods.go`;
- `bot/poll_media.go`;
- `telegram.PollMedia`, `telegram.Poll`, `telegram.PollOption`, and `telegram.InputPollOption`.

### Live Photos

Official 10.0 live-photo changes remain represented locally:

- `telegram.LivePhoto`;
- `Message.live_photo`;
- `ExternalReplyInfo.live_photo`;
- `sendLivePhoto`;
- `InputMediaLivePhoto` for media groups and edit media;
- `PaidMediaLivePhoto`;
- `InputPaidMediaLivePhoto`.

Coverage references:

- `docs/API_COVERAGE.md` media and paid-media rows;
- `docs/BOT_API_10_0_FINAL_AUDIT.md` live-photo section;
- `bot/media.go`;
- `bot/media_group.go`;
- `bot/paid_media_methods.go`;
- `telegram.LivePhoto` and paid-media variants.

### General And Managed Bot Changes

Official 10.0 general and managed-bot changes remain represented locally:

- bot-username send targets are accepted through existing `ChatIDString` behavior;
- business bot replies are covered through existing reply and business connection parameters;
- empty `sendMessageDraft.text` is covered;
- `BotAccessSettings`;
- `getManagedBotAccessSettings`;
- `setManagedBotAccessSettings`;
- `getUserPersonalChatMessages`.

Coverage references:

- `docs/API_COVERAGE.md` managed-bot, draft, business, and send rows;
- `docs/BOT_API_10_0_FINAL_AUDIT.md` general and managed-bot sections;
- `bot/checklist_draft_methods.go`;
- `bot/managed_bot_methods.go`;
- `telegram.BotAccessSettings`.

## Follow-Up

- Keep using [`docs/maintainer/BOT_API_UPDATE_CHECKLIST.md`](maintainer/BOT_API_UPDATE_CHECKLIST.md) for future upstream Bot API changes.
- Do not introduce code generation from this audit alone.
- Re-run a full audit only when Telegram publishes a Bot API release newer than 10.0 or when a field-level discrepancy is suspected.
