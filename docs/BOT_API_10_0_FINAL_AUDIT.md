# Bot API 10.0 Final Audit

Audit date: May 10, 2026.

Source of truth:

- [Official Telegram Bot API changelog](https://core.telegram.org/bots/api-changelog)
- [Official Telegram Bot API documentation](https://core.telegram.org/bots/api)

## Result

Bot API 10.0 code coverage is complete with documented architecture differences.

No missing Bot API 10.0 method wrappers, request parameters, result objects, update fields, or high-impact message/type fields were found in the final audit.

Live Telegram smoke remains manual-only because several Bot API 10.0 areas can interact with real chats, managed bots, business accounts, reactions, or guest-mode flows.

## Audited Changes

### Guest Mode

Official additions:

- `User.supports_guest_queries`
- `Message.guest_bot_caller_user`
- `Message.guest_bot_caller_chat`
- `Message.guest_query_id`
- `Update.guest_message`
- `SentGuestMessage`
- `answerGuestQuery`

Local coverage:

- `telegram.User.SupportsGuestQueries`
- `telegram.Message.GuestBotCallerUser`
- `telegram.Message.GuestBotCallerChat`
- `telegram.Message.GuestQueryID`
- `telegram.Update.GuestMessage`
- `telegram.SentGuestMessage`
- `(*bot.Bot).AnswerGuestQuery`
- `dispatch.GuestMessage` and `OnGuestMessage`

Verification:

- Unit and httptest coverage for guest message decoding, effective helpers, dispatch routing, `SentGuestMessage`, and `answerGuestQuery`.

### Chat Management

Official additions:

- `ChatMemberRestricted.can_react_to_messages`
- `ChatPermissions.can_react_to_messages`
- `getChatAdministrators.return_bots`
- `deleteAllMessageReactions`
- `deleteMessageReaction`
- server-side visibility for certain messages sent by other bots in groups

Local coverage:

- `telegram.ChatMember.CanReactToMessages`
- `telegram.ChatPermissions.CanReactToMessages`
- `bot.GetChatAdministratorsParams.ReturnBots`
- `(*bot.Bot).DeleteAllMessageReactions`
- `(*bot.Bot).DeleteMessageReaction`

Notes:

- Seeing certain messages sent by other bots in groups is a server-side delivery behavior. Existing `telegram.Message` decoding accepts bot senders and does not need a new field or method.

Verification:

- Unit and httptest coverage for permission/member decoding, `return_bots` request payloads, reaction deletion request payloads, validation, and error paths.

### Polls

Official additions:

- `InputMediaSticker`
- `InputMediaLocation`
- `InputMediaVenue`
- `PollMedia`
- `Poll.media`
- `Poll.explanation_media`
- `PollOption.media`
- `InputPollMedia`
- `sendPoll.media`
- `sendPoll.explanation_media`
- `InputPollOptionMedia`
- `InputPollOption.media`
- `Poll.members_only`
- `sendPoll.members_only`
- `Poll.country_codes`
- `sendPoll.country_codes`

Local coverage:

- Incoming `telegram.PollMedia`, `Poll.Media`, `Poll.ExplanationMedia`, `PollOption.Media`, `Poll.MembersOnly`, and `Poll.CountryCodes`.
- Outgoing `bot.InputPollMedia`, `bot.InputPollOptionMedia`, poll media wrappers, `SendPollParams.Media`, `ExplanationMedia`, `MembersOnly`, and `CountryCodes`.
- JSON and multipart upload support for poll media, explanation media, and poll-option media file fields.

Verification:

- Unit and httptest coverage for incoming poll media decoding, outgoing JSON payloads, multipart `attach://` payloads, validation, and one-option polls.

### Live Photos

Official additions:

- `LivePhoto`
- `InputMediaLivePhoto`
- `Message.live_photo`
- `ExternalReplyInfo.live_photo`
- `sendLivePhoto`
- `PaidMediaLivePhoto`
- `InputPaidMediaLivePhoto`
- live photos in `sendMediaGroup`
- live photos in `editMessageMedia`

Local coverage:

- `telegram.LivePhoto`
- `telegram.Message.LivePhoto`
- `telegram.ExternalReplyInfo.LivePhoto`
- `(*bot.Bot).SendLivePhoto`
- `telegram.PaidMediaLivePhoto`
- `bot.InputPaidMediaLivePhoto`
- `bot.InputMediaLivePhoto` in `sendMediaGroup` and `editMessageMedia`

Notes:

- The public `FileRef` architecture represents official `InputFile` as `FileID`, `FileURL`, or `FileUpload` where Telegram supports them.
- Live-photo URL media is intentionally rejected in places where the implementation requires `file_id` or multipart upload for the live-photo payload.

Verification:

- Unit and httptest coverage for decode contracts, JSON payloads, multipart uploads, paid media payloads, media groups, edit media, validation, API errors, invalid JSON, HTTP errors, and canceled contexts.

### Managed Bot Access And User Messages

Official additions:

- `BotAccessSettings`
- `getManagedBotAccessSettings`
- `setManagedBotAccessSettings`
- `getUserPersonalChatMessages`

Local coverage:

- `telegram.BotAccessSettings`
- `(*bot.Bot).GetManagedBotAccessSettings`
- `(*bot.Bot).SetManagedBotAccessSettings`
- `(*bot.Bot).GetUserPersonalChatMessages`

Verification:

- Unit and httptest coverage for settings decoding, method payloads, result decoding, validation, API errors, invalid JSON, HTTP errors, canceled contexts, and token redaction.

### General Behavior Changes

Official changes:

- Business Bots can manage user accounts without a Telegram Premium subscription.
- Bots can send messages to other bots via username when both bots enabled bot-to-bot communication.
- Business bots can reply to other bots when bot-to-bot communication is enabled.
- `sendMessageDraft.text` can be empty.

Local coverage:

- No client-side Premium gate exists for Business Bot methods, so no code change was needed.
- `bot.ChatIDString` already supports bot usernames for `chat_id`; audit tests cover `@other_bot` payloads.
- `SendMessageParams.BusinessConnectionID` and `ReplyParameters` support the business reply payload shape without client-side rejection.
- `SendMessageDraftParams.Text` now allows the zero value and serializes an explicit empty string.

Verification:

- Unit and httptest coverage for bot-username send payloads, business reply payloads, and empty-text `sendMessageDraft`.

## Verification Commands

Final verification for this audit:

```bash
go test ./bot
scripts/check.sh
go test -race ./bot ./dispatch ./middleware ./transport/longpoll ./transport/webhook ./internal/httpclient ./telegram
go test -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
```

The total statement coverage after the audit is 63.1%.

## Remaining Notes

- No known Bot API 10.0 code coverage blockers remain.
- Live smoke for Guest Mode, managed bots, business accounts, reaction deletion, and other state-changing or sensitive flows remains manual-only and requires explicit maintainer approval.
- Public release tags and GitHub Releases still require explicit maintainer approval.
