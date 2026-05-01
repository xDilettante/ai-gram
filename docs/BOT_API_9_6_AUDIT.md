# Bot API 9.6 Coverage Audit

## Source of truth

- [Official Telegram Bot API documentation](https://core.telegram.org/bots/api), fetched for the original audit on 2026-04-30 and rechecked for Stages 89-94 on 2026-05-01.
- [Official Telegram Bot API changelog](https://core.telegram.org/bots/api-changelog), especially the April 3, 2026 Bot API 9.6 entry.

The audit compares official method/type headings and high-impact object fields against the current local implementation. Stage notes below are updated as follow-up slices are implemented locally.

## Audit result

**Full coverage not yet reached.**

The current repository covers the large local Stage 66-96 workstream, including forum topics, reactions, inline mode, payments, paid media, Stars/gifts, subscription invite links, Managed Bots 9.6, Poll 9.6, WebApp/Mini App, Business API foundation/account/story/suggested posts/repost, games, Passport, lifecycle/profile read APIs, verification/user status APIs, chat member/boost updates, checklists, message drafts, structured poll options, reply/message metadata, prepared inline messages, reply markup completion, service-message metadata, and video metadata completion. The remaining gaps are concentrated in ChatFullInfo/full user-chat metadata and update shape parity.

## Implemented areas

- Chat member update and chat boost support: `ChatMemberUpdated`, `Update.my_chat_member`, `Update.chat_member`, `ChatBoostUpdated`, `ChatBoostRemoved`, `UserChatBoosts`, `GetUserChatBoosts`, `SetChatMemberTag`, `BanChatSenderChat`, and `UnbanChatSenderChat`.

- Core bot construction, token-safe HTTP calls, configurable base URL/client, JSON and multipart requests, typed API errors.
- Updates via `getUpdates`, managed long polling, inbound webhook handler, and JSON webhook management.
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

No Stage 88 missing method groups are currently tracked after Stage 96. Remaining work is concentrated in result/update shape parity and final official-doc verification.

## Missing types and fields

| Official name | Parent type | Why it matters | Suggested stage |
| --- | --- | --- | --- |
| `ChatFullInfo` | `getChat` result | Official `getChat` returns the extended chat object; the current method returns minimal `Chat`, so many current chat metadata fields are unavailable. Stage 89 kept the existing signature and documented the transition strategy instead of making an incidental breaking change. | Stage 97 |
| `User.language_code`, `is_premium`, `added_to_attachment_menu`, `can_join_groups`, `can_read_all_group_messages`, `supports_inline_queries`, `can_connect_to_business`, `has_main_web_app`, `has_topics_enabled`, `allows_users_to_create_topics` | `User` | Returned by `getMe`/user payloads and newer topic/business/profile capability checks. | Stage 97 |
| `Chat.is_forum`, `Chat.is_direct_messages` | `Chat` | Indicates forum and channel direct messages chats in lightweight chat payloads. | Stage 97 |
| `channel_post`, `edited_channel_post`, `poll` | `Update` | Missing update entry points still block channel posts and standalone poll updates. Stage 91 added chat member and chat boost updates. | Stage 97 |
| concrete `ChatMember*` variant structs | Chat member types | Stage 91 keeps the existing flat `ChatMember` struct and extends it with official 9.6 fields instead of introducing a breaking polymorphic API. Dedicated concrete variants remain a possible future refinement. | Stage 97 |
| `InputFile` official object | Upload parameters | The library intentionally uses `FileRef`/`FileUpload`; this is a naming/architecture mismatch to document, not necessarily a missing public type. | Needs verification |

## Potential mismatches / needs verification

- `getChat` currently returns `*telegram.Chat`; official docs return `ChatFullInfo`. Adding `ChatFullInfo` may require either a breaking signature change before stable release or a compatible new method/result strategy.
- `MessageId` is represented idiomatically as `telegram.MessageID`; this is acceptable but should be documented as a naming difference.
- `InputFile` is represented by `bot.FileID`, `bot.FileURL`, and `bot.FileUpload`; this is an intentional architecture difference, but future audit should ensure every official upload field is mapped.
- `SetWebhook` is JSON-only and does not support certificate upload; official `setWebhook` accepts an `InputFile` certificate.
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

1. **Stage 89 completed:** lifecycle and profile read APIs - `logOut`, `close`, `getUserProfilePhotos`, `getUserProfileAudios`, `getForumTopicIconStickers`; `ChatFullInfo` remains a documented getChat strategy mismatch.
2. **Stage 90 completed:** verification and user status APIs - `setUserEmojiStatus`, `verifyUser`, `verifyChat`, `removeUserVerification`, `removeChatVerification`.
3. **Stage 91 completed:** chat boosts, member updates, and sender-chat moderation - `getUserChatBoosts`, `setChatMemberTag`, `banChatSenderChat`, `unbanChatSenderChat`, chat boost/member update types.
4. **Stage 92 completed:** subscription invite links - `createChatSubscriptionInviteLink`, `editChatSubscriptionInviteLink`, invite link price/subscription fields.
5. **Stage 93 completed:** checklists, message drafts, and structured poll options - `sendChecklist`, `editMessageChecklist`, `sendMessageDraft`, `InputPollOption`, checklist message/service types.
6. **Stage 94 completed:** reply and message metadata types - `MessageOrigin*`, `ExternalReplyInfo`, `TextQuote`, `MaybeInaccessibleMessage`, `InaccessibleMessage`, `ReplyParameters` quote/cross-chat/checklist fields, and high-impact message metadata fields.
7. **Stage 95 completed:** prepared inline messages and reply-markup completion - `savePreparedInlineMessage`, `PreparedInlineMessage`, LoginUrl/switch-inline/copy/pay/request-poll/icon/style button fields.
8. **Stage 96 completed:** service/direct-message/story/media metadata - `repostStory`, video quality/cover/start metadata, shared user/chat service messages, chat backgrounds, video chats, proximity alerts, auto-delete timers, giveaway service messages, and paid/direct message price changes.
9. **Stage 97: ChatFullInfo/update shape strategy** - decide compatible `getChat` result strategy, channel post updates, standalone poll updates, fuller `User`/`Chat` metadata, and optional concrete chat member variants.
10. **Final audit after Stage 97** - rerun official method/type/field comparison and only then reconsider push/tag/release readiness.
