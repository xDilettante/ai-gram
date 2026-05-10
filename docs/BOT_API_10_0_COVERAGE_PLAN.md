# Bot API 10.0 Coverage Plan

Telegram Bot API 10.0 was released on May 8, 2026. This document tracks the update from the current Bot API 9.6-complete implementation to Bot API 10.0 coverage.

Source of truth:

- [Official Telegram Bot API changelog](https://core.telegram.org/bots/api-changelog)
- [Official Telegram Bot API documentation](https://core.telegram.org/bots/api)

## Status

Bot API 10.0 code coverage is complete with documented architecture differences. The final audit is recorded in [`docs/BOT_API_10_0_FINAL_AUDIT.md`](BOT_API_10_0_FINAL_AUDIT.md).

Bot API 9.6 remains complete with documented architecture differences.

## New or changed areas

### Guest Mode

- [x] `User.supports_guest_queries`
- [x] `Message.guest_bot_caller_user`
- [x] `Message.guest_bot_caller_chat`
- [x] `Message.guest_query_id`
- [x] `Update.guest_message`
- [x] `SentGuestMessage`
- [x] `answerGuestQuery`

### Chat Management

- [x] `ChatMemberRestricted.can_react_to_messages`
- [x] `ChatPermissions.can_react_to_messages`
- [x] `getChatAdministrators.return_bots`
- [x] `deleteAllMessageReactions`
- [x] `deleteMessageReaction`

### Polls

- [x] `InputMediaSticker`
- [x] `InputMediaLocation`
- [x] `InputMediaVenue`
- [x] `PollMedia`
- [x] `Poll.media`
- [x] `Poll.explanation_media`
- [x] `PollOption.media`
- [x] `InputPollMedia`
- [x] `sendPoll.media`
- [x] `sendPoll.explanation_media`
- [x] `InputPollOptionMedia`
- [x] `InputPollOption.media`
- [x] `Poll.members_only`
- [x] `sendPoll.members_only`
- [x] `Poll.country_codes`
- [x] `sendPoll.country_codes`
- [x] Allow one-option polls in `SendPollParams` validation.
- [x] Multipart uploads for poll media file fields.

### Live Photos

- [x] `LivePhoto`
- [x] `InputMediaLivePhoto`
- [x] `Message.live_photo`
- [x] `ExternalReplyInfo.live_photo`
- [x] `sendLivePhoto`
- [x] `PaidMediaLivePhoto`
- [x] `InputPaidMediaLivePhoto`
- [x] Allow live photos in `sendMediaGroup`.
- [x] Allow live photos in `editMessageMedia`.

### Managed Bot Access And User Messages

- [x] `BotAccessSettings`
- [x] `getManagedBotAccessSettings`
- [x] `setManagedBotAccessSettings`
- [x] `getUserPersonalChatMessages`

### General Behavior Changes

- [x] Review server-side delivery of certain messages sent by other bots in groups for required type or validation changes.
- [x] Review Business Bots managing user accounts without Telegram Premium for required type or validation changes.
- [x] `sendMessageDraft` allows empty text.
- [x] Review bot-to-bot username sends and business bot replies for required type or validation changes.

## Implementation order

1. Add passive decode fields and tests first: new `User`, `Message`, `Update`, `Poll`, and permission fields.
2. Add read/delete method wrappers with httptest coverage.
3. Add guest mode method/types.
4. Add poll media input/output types and validation updates.
5. Add live photo send/edit/media-group support.
6. Add managed access settings and personal chat messages.
7. Update API coverage docs, README badges, release notes, and final audit.

## Verification

Each implementation slice should include focused unit or httptest coverage. Before marking this plan complete, run:

```bash
scripts/check.sh
go test -race ./bot ./dispatch ./middleware ./transport/longpoll ./transport/webhook ./internal/httpclient ./telegram
go test -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
```

Live smoke remains manual-only and requires explicit maintainer approval.
