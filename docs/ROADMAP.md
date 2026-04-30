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

## vNext Bot API 9.6 full coverage workstream

Strategic change: the small v0.3 release plan is superseded. The current goal remains full Telegram Bot API 9.6 method, type, and update coverage before the next push, tag, or GitHub Release. See [`docs/BOT_API_9_6_COVERAGE_PLAN.md`](BOT_API_9_6_COVERAGE_PLAN.md) for the working plan and [`docs/BOT_API_9_6_AUDIT.md`](BOT_API_9_6_AUDIT.md) for the Stage 88 official-doc audit.

Repository policy for this workstream:

- `v0.2.0` remains the latest public release.
- Continue local-only development with verified local commits.
- Do not suggest pushing main after each local commit.
- Do not create tags or GitHub Releases until full Bot API 9.6 coverage is complete.
- Do not run `git push` unless the user explicitly asks.

Stage 88 outcome: **full coverage is not yet reached**. The largest recent slices are implemented locally through Games and Passport, but the audit found missing official methods and type/field groups.

Recommended local-only implementation order after Stage 88:

1. Lifecycle/profile read APIs: `logOut`, `close`, profile photos/audios, forum topic icon stickers.
2. Verification and user status APIs.
3. Chat boosts, chat-member updates/tags, and sender-chat moderation.
4. Subscription invite links.
5. Checklists, message drafts, and structured poll options.
6. Business/direct-message story completion.
7. Prepared inline messages and reply-markup completion.
8. Message field completeness.
9. Service-message completeness.
10. Final official-doc audit and release-readiness review only after the missing checklist is empty.

Live smoke policy:

- Safe/read-only flows may be live-smoked only when explicitly useful.
- Admin/state-changing flows require a dedicated test chat and explicit user confirmation.
- Payments, Passport, Business, Managed Bots, gifts, and Stars flows require explicit confirmation.
- Destructive/admin flows must not be auto-smoked.

## Later

- Payments
  - Invoice, shipping query, pre-checkout query, refunds, paid media.
- Passport decryption helpers
  - Intentionally out of scope for the typed Bot API wrapper unless a future product decision adds them.
- Business APIs
  - Business connection/message features and related account metadata.
- Stars/gifts
  - Review against the current official Bot API before planning.
- Codegen
  - Consider generation only after the hand-written public API shape is proven.
  - Do not introduce codegen just to inflate method count.
