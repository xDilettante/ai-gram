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

Status: current `main` contains the planned v0.2 expansion slices. Recommended next step: stabilize for v0.2.0 instead of adding more methods immediately, unless a specific missing Bot API area is required before release.

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

Remaining candidate slices before or after v0.2.0:

- Chat management methods: `setChatTitle`, `setChatDescription`, `setChatPhoto`, `deleteChatPhoto`, `leaveChat`.
- Forum topic methods.
- Reactions.
- Inline mode basics.
- Remaining sticker set methods.
- Bot profile methods.

Milestone recommendation:

- Prefer a v0.2.0 stabilization pass now if the goal is to publish a coherent expanded API milestone.
- Continue coverage before v0.2.0 only if chat management, forum topics, reactions, inline mode, or sticker set management is required for the release boundary.

## v0.3 Advanced coverage candidates

The v0.3 scope should be chosen after the v0.2.0 decision. Good candidates:

- Chat management methods:
  - title/description/photo/sticker-set/default permissions follow-ups
  - leave chat and related lifecycle methods
- Forum topics:
  - create/edit/close/reopen/delete topic methods
  - general forum topic helpers
- Reactions:
  - message reaction methods
  - reaction update/type coverage
- Inline mode:
  - inline query result types
  - `answerInlineQuery`
  - chosen inline result handling
- Sticker set management:
  - create/add/set/delete sticker set methods
  - custom emoji/sticker metadata methods

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
