# Roadmap

This roadmap is intentionally pragmatic: keep `ai-gram` useful and stable as a Go Telegram Bot API library without turning it into a large framework too early.

## v0.1 stabilization

Status: completed and released as v0.1.1 for canonical public Go module usage.

Completed scope:

- Typed Bot construction, configurable base URL and HTTP client, private token storage, and token-safe diagnostics.
- Typed Bot API errors and consistent `errors.As` support.
- Updates, webhook receiver, webhook management, and long polling runner.
- Dispatcher/router and essential middleware including recovery, timeout, observability, and access control.
- Text/media sends, reply markup, reply parameters, thread IDs, callbacks, edit/delete, forward/copy, chat actions, chat info, file upload/download, and safe examples.
- Manual/deploy/live smoke documentation, API coverage, roadmap, release checklist, and safe logs.

## v0.2 Bot API coverage expansion

Status: completed and released as v0.2.0.

Completed slices:

- Additional send methods:
  - `SendContact`
  - `SendLocation`
  - `SendVenue`
  - `SendDice`
  - `SendSticker`
  - `SendAnimation`
  - `SendVideoNote`
- Polls:
  - `SendPoll`
  - `StopPoll`
- Media groups:
  - `SendMediaGroup`
  - `InputMediaPhoto`
  - `InputMediaVideo`
  - `InputMediaAudio`
  - `InputMediaDocument`
- Bot commands/menu:
  - `SetMyCommands`
  - `DeleteMyCommands`
  - `GetMyCommands`
  - `SetChatMenuButton`
  - `GetChatMenuButton`
  - `SetMyDefaultAdministratorRights`
- Invite links:
  - `ExportChatInviteLink`
  - `CreateChatInviteLink`
  - `EditChatInviteLink`
  - `RevokeChatInviteLink`
- Join requests:
  - `ApproveChatJoinRequest`
  - `DeclineChatJoinRequest`
  - `telegram.ChatJoinRequest` update decoding
  - dispatch helpers for join request updates
- Admin management:
  - `PromoteChatMember`
  - `SetChatAdministratorCustomTitle`
  - `SetChatPermissions`
- Smoke/docs:
  - v0.2 send-method smoke script
  - `SendMediaGroup` smoke script with generated upload fallback
  - English public documentation cleanup
  - v0.2 checkpoint document

Verification status:

- Unit/httptest coverage exists for the implemented v0.2 API methods.
- Safe live smoke has covered contact/location/venue/poll/stop-poll/dice and `SendMediaGroup` generated upload fallback.
- State-changing/admin methods are intentionally documented as manual-only and not auto-smoked.

Remaining candidate slices after v0.2.0:

- Chat management methods: `setChatTitle`, `setChatDescription`, `setChatPhoto`, `deleteChatPhoto`, `leaveChat`.
- Forum topic methods.
- Reactions.
- Inline mode basics.
- Remaining sticker set methods.
- Bot profile methods.

Milestone outcome:

- v0.2.0 was released as the coherent expanded API milestone.
- Chat management, forum topics, reactions, and inline mode are now planned for v0.3 instead of extending the v0.2 boundary.

## v0.3 Group/chat administration and advanced interaction surfaces

Focus: extend the v0.2 group/admin foundation into practical chat administration and advanced interaction surfaces while keeping live verification safe and explicit. See [`docs/V0_3_PLAN.md`](V0_3_PLAN.md) for the detailed plan.

Planned slices:

- Chat management methods:
  - `SetChatTitle`
  - `SetChatDescription`
  - `SetChatPhoto`
  - `DeleteChatPhoto`
  - `LeaveChat`
- Forum topics:
  - `CreateForumTopic`
  - `EditForumTopic`
  - `CloseForumTopic`
  - `ReopenForumTopic`
  - `DeleteForumTopic`
  - `UnpinAllForumTopicMessages`
- Reactions:
  - `SetMessageReaction`
  - reaction type/update support as needed
- Inline mode basics:
  - `AnswerInlineQuery`
  - minimal `InlineQuery` update/type support
  - minimal `InlineQueryResult` variants

Recommended order:

1. Chat management
2. Forum topics
3. Reactions
4. Inline mode basics

Live smoke policy:

- Safe/read-only flows may be live-smoked.
- Admin/state-changing flows require a dedicated test chat and explicit user confirmation.
- Destructive/admin flows must not be auto-smoked.

## Later

- Payments
  - Invoice, shipping query, pre-checkout query, refunds, paid media.
- Passport
  - Passport data and error reporting methods.
- Games
  - Game methods and scores.
- Business APIs
  - Business connection/message features and related account metadata.
- Stars/gifts
  - Review against the current official Bot API before planning.
- Codegen
  - Consider generation only after the hand-written public API shape is proven.
  - Do not introduce codegen just to inflate method count.
