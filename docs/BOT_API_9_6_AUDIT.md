# Bot API 9.6 Coverage Audit

## Source of truth

- [Official Telegram Bot API documentation](https://core.telegram.org/bots/api), fetched for the original audit on 2026-04-30 and rechecked through Stage 99 on 2026-05-02.
- [Official Telegram Bot API changelog](https://core.telegram.org/bots/api-changelog), especially the April 3, 2026 Bot API 9.6 entry.
- Latest sources for release-readiness status: [`docs/BOT_API_9_6_FINAL_AUDIT.md`](BOT_API_9_6_FINAL_AUDIT.md) and [`docs/maintainer/BOT_API_9_6_RELEASE_READINESS.md`](maintainer/BOT_API_9_6_RELEASE_READINESS.md).

The audit compares official method/type headings and high-impact object fields against the current local implementation. Stage notes below are updated as follow-up slices are implemented locally.

## Audit result

**Full coverage reached with documented architecture differences.**

Stage 98 found wrappers for all 169 official methods and no missing fields in the audited `User`, `Chat`, `ChatFullInfo`, `Update`, `Message`, `ReplyParameters`, `CallbackQuery`, `Video`, sticker, and keyboard field tables after correcting `Message.giveaway`. Stage 99 resolved the remaining hard blocker by adding `setWebhook.certificate` multipart upload support. Stage 100 records release-readiness verification and manual-only smoke planning. See [`docs/BOT_API_9_6_FINAL_AUDIT.md`](BOT_API_9_6_FINAL_AUDIT.md) and [`docs/maintainer/BOT_API_9_6_RELEASE_READINESS.md`](maintainer/BOT_API_9_6_RELEASE_READINESS.md) for the latest readiness status.

## Implemented areas

- Chat member update and chat boost support: `ChatMemberUpdated`, `Update.my_chat_member`, `Update.chat_member`, `ChatBoostUpdated`, `ChatBoostRemoved`, `UserChatBoosts`, `GetUserChatBoosts`, `SetChatMemberTag`, `BanChatSenderChat`, and `UnbanChatSenderChat`.

- Core bot construction, token-safe HTTP calls, configurable base URL/client, JSON and multipart requests, typed API errors.
- Updates via `getUpdates`, managed long polling, inbound webhook handler, and webhook management with JSON mode plus multipart certificate upload.
- Send/media methods for text, media, animation/sticker/video note, media groups, contact/location/venue, poll/dice/game, invoice, paid media, gifts, and business-enabled sends currently implemented.
- Edit/delete/forward/copy/batch methods, including `EditMessageMedia`, live-location edit/stop, batch forward/copy/delete, and business connection fields where currently supported.
- Chat management, moderation, admin methods, regular and subscription invite links, join requests, forum topics, reactions, commands/menu, bot profile/metadata, sticker set management.
- Inline mode query/result/input-content coverage for the implemented inline result families.
- Payments/invoices, paid media, Stars transaction/refund basics, gifts, business gifts, and Premium subscription gift methods implemented in the recent local stages.
- Managed Bots 9.6 types/methods and Poll 9.6 fields/service messages, including structured `InputPollOption` support.
- WebApp/Mini App Bot API surface, Business API foundation/account/story/suggested post methods, games, and Telegram Passport types/error methods.
- Lifecycle/profile read APIs: `logOut`, `close`, `getUserProfilePhotos`, `getUserProfileAudios`, and `getForumTopicIconStickers`.
- Verification/status APIs: `setUserEmojiStatus`, `verifyUser`, `verifyChat`, `removeUserVerification`, and `removeChatVerification`.
- Checklist/message draft APIs: `sendChecklist`, `editMessageChecklist`, `sendMessageDraft`, checklist message/service types, and manual-only safety documentation.
- Reply/message metadata: `MessageOrigin` variants, `ExternalReplyInfo`, `TextQuote`, `InaccessibleMessage`, `MaybeInaccessibleMessage`, `ReplyParameters` quote/cross-chat/checklist fields, and high-impact `Message` metadata such as `forward_origin`, `reply_to_message`, `external_reply`, `quote`, `reply_to_story`, `direct_messages_topic`, `suggested_post_info`, `pinned_message`, sender metadata, caption/media flags, star/effect fields, and `reply_markup`.
- Prepared inline and reply markup completion: `SavePreparedInlineMessage`, `PreparedInlineMessage`, `LoginUrl`, `SwitchInlineQueryChosenChat`, `CopyTextButton`, `KeyboardButtonPollType`, `KeyboardButton.request_poll`, `InlineKeyboardButton.pay`, and keyboard button `icon_custom_emoji_id`/`style` fields.
- Unit/httptest coverage for implemented method families and token/payload redaction checks in sensitive areas.

## Missing methods

No official method wrappers are missing after Stage 99 (169/169 official methods have exported `(*bot.Bot)` wrappers). No known method behavior blockers remain.

## Missing types and fields

| Official name | Parent type | Why it matters | Suggested stage |
| --- | --- | --- | --- |
| concrete `ChatMember*` variant structs | Chat member types | Stage 91/97 keep the existing flat `ChatMember` struct and extend it with official 9.6 fields instead of introducing a breaking polymorphic API. Dedicated concrete variants remain a possible future refinement, not a blocking decode gap. | Future refinement |
| `InputFile` official object | Upload parameters | The library intentionally uses `FileRef`/`FileUpload`; direct upload-capable methods are covered, including upload-only `setWebhook.certificate` via `FileUpload`. | Architecture note |

## Potential mismatches / needs verification

- `GetChat` remains backward-compatible and returns `*telegram.Chat`; `GetChatFullInfo` now calls the same official `getChat` method and decodes the full `ChatFullInfo` result.
- `MessageId` is represented idiomatically as `telegram.MessageID`; this is acceptable but should be documented as a naming difference.
- `InputFile` is represented by `bot.FileID`, `bot.FileURL`, and `bot.FileUpload`; this is an intentional architecture difference. Stage 99 confirmed every direct official upload field is mapped, with `setWebhook.certificate` restricted to upload-only `FileUpload` as required by official docs.
- `sendPoll` still exposes the legacy singular `correct_option_id` for backward compatibility while official 9.6 replaced it with `correct_option_ids`; validation should continue rejecting ambiguous use.
- `SendPollParams` keeps legacy `Options []string` while adding `OptionObjects []telegram.InputPollOption`; validation rejects ambiguous use and serializes both shapes through the official `options` field.
- `ReactionType` and other polymorphic decoders should be rechecked when unknown official variants appear; current tests generally fail safely on unknown types.
- Several validation rules intentionally avoid hardcoding Telegram upper limits; this is safer for forward compatibility but should be reviewed for methods with official hard limits.

## Manual-only live smoke areas

These areas must remain manual-only and require explicit user confirmation plus dedicated test assets/accounts:

- payments, invoices, paid media, Stars, gifts, business gifts, subscription invite links, Premium subscription gifts, and refunds;
- Passport data and Passport error reporting;
- Business APIs, business messages, business account profile changes, stories, suggested posts, and direct messages;
- Managed bot token methods and prepared button/inline methods that create reusable identifiers;
- admin/destructive chat methods including bans, restrictions, promotions, chat photo/title/description changes, leave chat, invite links, and join request decisions;
- sticker set mutation methods;
- forum topic mutation methods;
- reaction methods that change real message state;
- games requiring BotFather game setup;
- inline mode and Mini App flows that require BotFather/Mini App setup;
- checklist/message draft methods because they mutate user-visible state.

## Recommended next stages

1. **Stage 89 completed:** lifecycle and profile read APIs - `logOut`, `close`, `getUserProfilePhotos`, `getUserProfileAudios`, `getForumTopicIconStickers`.
2. **Stage 90 completed:** verification and user status APIs - `setUserEmojiStatus`, `verifyUser`, `verifyChat`, `removeUserVerification`, `removeChatVerification`.
3. **Stage 91 completed:** chat boosts, member updates, and sender-chat moderation - `getUserChatBoosts`, `setChatMemberTag`, `banChatSenderChat`, `unbanChatSenderChat`, chat boost/member update types.
4. **Stage 92 completed:** subscription invite links - `createChatSubscriptionInviteLink`, `editChatSubscriptionInviteLink`, invite link price/subscription fields.
5. **Stage 93 completed:** checklists, message drafts, and structured poll options - `sendChecklist`, `editMessageChecklist`, `sendMessageDraft`, `InputPollOption`, checklist message/service types.
6. **Stage 94 completed:** reply and message metadata types - `MessageOrigin*`, `ExternalReplyInfo`, `TextQuote`, `MaybeInaccessibleMessage`, `InaccessibleMessage`, `ReplyParameters` quote/cross-chat/checklist fields, and high-impact message metadata fields.
7. **Stage 95 completed:** prepared inline messages and reply-markup completion - `savePreparedInlineMessage`, `PreparedInlineMessage`, LoginUrl/switch-inline/copy/pay/request-poll/icon/style button fields.
8. **Stage 96 completed:** service/direct-message/story/media metadata - `repostStory`, video quality/cover/start metadata, shared user/chat service messages, chat backgrounds, video chats, proximity alerts, auto-delete timers, giveaway service messages, and paid/direct message price changes.
9. **Stage 97 completed:** ChatFullInfo/update shape strategy - `GetChatFullInfo`, fuller `User`/`Chat` metadata, channel post updates, standalone poll updates, and compatible flat `ChatMember` strategy.
10. **Stage 98 completed:** final official-doc audit found all 169 official method wrappers present, corrected `Message.giveaway`, and identified `setWebhook.certificate` upload as the remaining hard blocker.
11. **Stage 99 completed:** implemented `SetWebhook` certificate upload / multipart support and reran a short final audit. No known Bot API 9.6 code coverage blockers remain.
12. **Stage 100 completed:** release-readiness verification and manual-only smoke planning. Later stages published `main` only after explicit approval; tags and GitHub Releases still require explicit maintainer approval.
