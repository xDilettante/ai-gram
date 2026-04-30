# Bot API 9.6 Coverage Audit

## Source of truth

- [Official Telegram Bot API documentation](https://core.telegram.org/bots/api), fetched for this audit on 2026-04-30.
- [Official Telegram Bot API changelog](https://core.telegram.org/bots/api-changelog), especially the April 3, 2026 Bot API 9.6 entry.

The audit compares official method/type headings and high-impact object fields against the current local implementation. It is intentionally documentation-only: no new Bot API methods are implemented in this stage.

## Audit result

**Full coverage not yet reached.**

The current repository covers the large local Stage 66-87 workstream, including forum topics, reactions, inline mode, payments, paid media, Stars/gifts, Managed Bots 9.6, Poll 9.6, WebApp/Mini App, Business API foundation/account/story/suggested posts, games, and Passport. The remaining gaps are concentrated in legacy lifecycle methods, checklist/message-draft APIs, subscription invite links, chat boost/member-update/profile-audio surfaces, verification/emoji-status methods, and broad incoming message/type completeness.

## Implemented areas

- Core bot construction, token-safe HTTP calls, configurable base URL/client, JSON and multipart requests, typed API errors.
- Updates via `getUpdates`, managed long polling, inbound webhook handler, and JSON webhook management.
- Send/media methods for text, media, animation/sticker/video note, media groups, contact/location/venue, poll/dice/game, invoice, paid media, gifts, and business-enabled sends currently implemented.
- Edit/delete/forward/copy/batch methods, including `EditMessageMedia`, live-location edit/stop, batch forward/copy/delete, and business connection fields where currently supported.
- Chat management, moderation, admin methods, invite links, join requests, forum topics, reactions, commands/menu, bot profile/metadata, sticker set management.
- Inline mode query/result/input-content coverage for the implemented inline result families.
- Payments/invoices, paid media, Stars transaction/refund basics, gifts, business gifts, and Premium subscription gift methods implemented in the recent local stages.
- Managed Bots 9.6 types/methods and Poll 9.6 fields/service messages.
- WebApp/Mini App Bot API surface, Business API foundation/account/story/suggested post methods, games, and Telegram Passport types/error methods.
- Unit/httptest coverage for implemented method families and token/payload redaction checks in sensitive areas.

## Missing methods

| Official method name | Area | Risk level | Suggested implementation stage |
| --- | --- | --- | --- |
| `logOut` | Bot session lifecycle | state-changing | Stage 89: lifecycle/profile read APIs |
| `close` | Bot session lifecycle | state-changing | Stage 89: lifecycle/profile read APIs |
| `getUserProfilePhotos` | User/profile media | safe | Stage 89: lifecycle/profile read APIs |
| `getUserProfileAudios` | User/profile media | safe | Stage 89: lifecycle/profile read APIs |
| `setUserEmojiStatus` | User/profile status | state-changing | Stage 90: verification/status APIs |
| `getUserChatBoosts` | Chat boosts | safe | Stage 91: chat boosts/member updates |
| `setChatMemberTag` | Chat member tags | admin/state-changing | Stage 91: chat boosts/member updates |
| `banChatSenderChat` | Sender chat moderation | admin/destructive | Stage 91: chat boosts/member updates |
| `unbanChatSenderChat` | Sender chat moderation | admin/state-changing | Stage 91: chat boosts/member updates |
| `createChatSubscriptionInviteLink` | Paid subscription invite links | state-changing/payment-related | Stage 92: subscription invite links |
| `editChatSubscriptionInviteLink` | Paid subscription invite links | state-changing/payment-related | Stage 92: subscription invite links |
| `getForumTopicIconStickers` | Forum topic metadata | safe | Stage 89: lifecycle/profile read APIs |
| `verifyUser` | Verification | state-changing/sensitive | Stage 90: verification/status APIs |
| `verifyChat` | Verification | state-changing/sensitive | Stage 90: verification/status APIs |
| `removeUserVerification` | Verification | state-changing/sensitive | Stage 90: verification/status APIs |
| `removeChatVerification` | Verification | state-changing/sensitive | Stage 90: verification/status APIs |
| `sendChecklist` | Checklists | state-changing/send | Stage 93: checklists and drafts |
| `editMessageChecklist` | Checklists | state-changing/edit | Stage 93: checklists and drafts |
| `sendMessageDraft` | Message drafts | state-changing/send | Stage 93: checklists and drafts |
| `repostStory` | Business stories | state-changing/business | Stage 94: business story completion |
| `savePreparedInlineMessage` | Mini App / inline prepared messages | state-changing/sensitive identifier | Stage 95: prepared inline messages |

## Missing types and fields

| Official name | Parent type | Why it matters | Suggested stage |
| --- | --- | --- | --- |
| `ChatFullInfo` | `getChat` result | Official `getChat` returns the extended chat object; the current method returns minimal `Chat`, so many current chat metadata fields are unavailable. | Stage 89 / Stage 91 |
| `UserProfilePhotos`, `UserProfileAudios` | `getUserProfilePhotos`, `getUserProfileAudios` results | Required result types for the missing profile media methods. | Stage 89 |
| `User.language_code`, `is_premium`, `added_to_attachment_menu`, `can_join_groups`, `can_read_all_group_messages`, `supports_inline_queries`, `can_connect_to_business`, `has_main_web_app`, `has_topics_enabled`, `allows_users_to_create_topics` | `User` | Returned by `getMe`/user payloads and newer topic/business/profile capability checks. | Stage 89 |
| `Chat.is_forum`, `Chat.is_direct_messages` | `Chat` | Indicates forum and channel direct messages chats in lightweight chat payloads. | Stage 91 |
| `channel_post`, `edited_channel_post`, `poll`, `my_chat_member`, `chat_member`, `chat_boost`, `removed_chat_boost` | `Update` | Missing update entry points block channel posts, standalone poll updates, chat member changes, and chat boost updates. | Stage 91 |
| `ChatMemberUpdated`, concrete `ChatMember*` variants, `ChatBoost*`, `UserChatBoosts` | Chat member / boost types | Required for `my_chat_member`, `chat_member`, `getUserChatBoosts`, and boost updates. | Stage 91 |
| `ChatBoostAdded`, `ChatBackground`, `BackgroundFill*`, `BackgroundType*` | `Message` service messages | Needed to decode chat boost and background service messages. | Stage 91 |
| `MessageOrigin*`, `ExternalReplyInfo`, `TextQuote`, `MaybeInaccessibleMessage`, `InaccessibleMessage` | `Message` reply/forward fields | Current message decoding lacks official forward/reply metadata such as `forward_origin`, `external_reply`, `quote`, and inaccessible pinned messages. | Stage 96: message field completeness |
| `DirectMessagesTopic`, `SuggestedPostInfo` | `Message` | Required for channel direct messages and suggested post metadata. | Stage 94 / Stage 96 |
| `Message.direct_messages_topic`, `sender_chat`, `sender_boost_count`, `sender_tag`, `forward_origin`, `is_topic_message`, `is_automatic_forward`, `reply_to_message`, `external_reply`, `quote`, `reply_to_story`, `reply_to_checklist_task_id`, `via_bot`, `edit_date`, `has_protected_content`, `is_paid_post`, `media_group_id`, `author_signature`, `paid_star_count`, `link_preview_options`, `suggested_post_info`, `effect_id`, `story`, `show_caption_above_media`, `has_media_spoiler`, `reply_markup` | `Message` | High-impact official message fields are not decoded yet; several affect business/direct messages, captions, message effects, stars, and replies. | Stage 96 |
| `Checklist`, `ChecklistTask`, `InputChecklist`, `InputChecklistTask`, `ChecklistTasksDone`, `ChecklistTasksAdded` | Checklist API and message fields | Required for `sendChecklist`, `editMessageChecklist`, and checklist service messages. | Stage 93 |
| `Message.checklist`, `checklist_tasks_done`, `checklist_tasks_added` | `Message` | Required to decode checklist messages and service updates. | Stage 93 |
| `Poll.question_entities` | `Poll` | Poll 9.6 text entity support is almost complete, but question entities remain missing. | Stage 93 |
| `ReplyParameters.chat_id`, `quote`, `quote_parse_mode`, `quote_entities`, `quote_position`, `checklist_task_id` | `ReplyParameters` | Current reply parameter support is partial and misses quote/cross-chat/checklist reply fields. | Stage 96 |
| `InputPollOption` | `sendPoll` options | Official poll options are structured and can include text entities; current params still use strings. | Stage 93 |
| `KeyboardButton.icon_custom_emoji_id`, `style`, `request_poll`; `KeyboardButtonPollType` | Reply keyboard | Bot API 9.4/9.6 keyboard support is incomplete for custom emoji/style and poll request buttons. | Stage 95 |
| `InlineKeyboardButton.icon_custom_emoji_id`, `style`, `login_url`, `switch_inline_query*`, `copy_text`, `pay`; `LoginUrl`, `SwitchInlineQueryChosenChat`, `CopyTextButton` | Inline keyboard | Current inline keyboard support lacks several official button modes. | Stage 95 |
| `Message.users_shared`, `chat_shared`; `SharedUser`, `UsersShared`, `ChatShared` | Request keyboard service messages | Required to decode user/chat sharing responses from keyboard request buttons. | Stage 95 |
| `Video.cover`, `start_timestamp`, `qualities`; `VideoQuality` | `Video` | Official video metadata includes cover/start and alternative qualities. | Stage 96 |
| `Message.story` and incoming story/direct-message fields | `Message` / business stories | Business story methods and basic `Story`/story-area types exist, but incoming story message metadata is not fully decoded. | Stage 94 / Stage 96 |
| `Giveaway*` types and `Message.giveaway*` fields | Giveaway service messages | Giveaway messages/service states are not decoded. | Stage 97: giveaway/background service messages |
| `VideoChat*`, `ProximityAlertTriggered`, `MessageAutoDeleteTimerChanged` | Service messages | Legacy service-message coverage remains incomplete. | Stage 97 |
| `PaidMessagePriceChanged`, `DirectMessagePriceChanged`, `paid_star_count`, `is_paid_post` | Paid/direct message service fields | Needed for paid message and channel direct-message service state. | Stage 96 / Stage 97 |
| `PreparedInlineMessage` | `savePreparedInlineMessage` result | Required result type for the missing Mini App prepared inline method. | Stage 95 |
| `InputFile` official object | Upload parameters | The library intentionally uses `FileRef`/`FileUpload`; this is a naming/architecture mismatch to document, not necessarily a missing public type. | Needs verification |

## Potential mismatches / needs verification

- `getChat` currently returns `*telegram.Chat`; official docs return `ChatFullInfo`. Adding `ChatFullInfo` may require either a breaking signature change before stable release or a compatible new method/result strategy.
- `MessageId` is represented idiomatically as `telegram.MessageID`; this is acceptable but should be documented as a naming difference.
- `InputFile` is represented by `bot.FileID`, `bot.FileURL`, and `bot.FileUpload`; this is an intentional architecture difference, but future audit should ensure every official upload field is mapped.
- `SetWebhook` is JSON-only and does not support certificate upload; official `setWebhook` accepts an `InputFile` certificate.
- `sendPoll` still exposes the legacy singular `correct_option_id` for backward compatibility while official 9.6 replaced it with `correct_option_ids`; validation should continue rejecting ambiguous use.
- Poll options are still string-based in `SendPollParams`; official `InputPollOption` supports entity-aware option text.
- `ReactionType` and other polymorphic decoders should be rechecked when unknown official variants appear; current tests generally fail safely on unknown types.
- Inline mode was audited earlier, but inline keyboard button support still lacks `login_url`, switch-inline, copy-text, pay, icon, and style fields.
- Business story/account methods exist, but `repostStory` and some incoming story/direct-message fields remain pending.
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
- checklist/message draft methods once implemented because they mutate user-visible state.

## Recommended next stages

1. **Stage 89: lifecycle and profile read APIs** - `logOut`, `close`, `getUserProfilePhotos`, `getUserProfileAudios`, `getForumTopicIconStickers`, `ChatFullInfo` planning.
2. **Stage 90: verification and user status APIs** - `setUserEmojiStatus`, `verifyUser`, `verifyChat`, `removeUserVerification`, `removeChatVerification`.
3. **Stage 91: chat boosts, member updates, and sender-chat moderation** - `getUserChatBoosts`, `setChatMemberTag`, `banChatSenderChat`, `unbanChatSenderChat`, chat boost/member update types.
4. **Stage 92: subscription invite links** - `createChatSubscriptionInviteLink`, `editChatSubscriptionInviteLink`, invite link price/subscription fields audit.
5. **Stage 93: checklists, message drafts, and structured poll options** - `sendChecklist`, `editMessageChecklist`, `sendMessageDraft`, `InputPollOption`, checklist message/service types.
6. **Stage 94: business/direct-message story completion** - `repostStory`, direct message topic fields, incoming story metadata.
7. **Stage 95: prepared inline messages and reply-markup completion** - `savePreparedInlineMessage`, `PreparedInlineMessage`, LoginUrl/switch-inline/copy/pay/request-poll/icon/style button fields.
8. **Stage 96: message field completeness pass** - forward/reply origins, quote/external reply, `ReplyParameters` completion, video quality/cover/start metadata, paid/direct message metadata.
9. **Stage 97: service-message completeness pass** - giveaways, chat backgrounds, video chats, proximity alerts, auto-delete timers, shared users/chats, and remaining service messages.
10. **Final audit after Stage 97** - rerun official method/type/field comparison and only then reconsider push/tag/release readiness.
