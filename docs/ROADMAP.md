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

## v0.4 Bot API 10.0 complete

Strategic change: the small v0.3 release plan is superseded. Code coverage for Telegram Bot API 10.0 is complete with documented architecture differences. See [`docs/BOT_API_9_6_COVERAGE_PLAN.md`](BOT_API_9_6_COVERAGE_PLAN.md), [`docs/BOT_API_9_6_FINAL_AUDIT.md`](BOT_API_9_6_FINAL_AUDIT.md), [`docs/BOT_API_10_0_COVERAGE_PLAN.md`](BOT_API_10_0_COVERAGE_PLAN.md), [`docs/BOT_API_10_0_FINAL_AUDIT.md`](BOT_API_10_0_FINAL_AUDIT.md), and [`docs/maintainer/BOT_API_10_0_RELEASE_READINESS.md`](maintainer/BOT_API_10_0_RELEASE_READINESS.md).

Repository status for this workstream:

- The public repository is available at <https://github.com/xDilettante/ai-gram>.
- `v0.3.0`, `v0.4.0`, and `v0.5.0` are already published; `v0.6.0` is the current regular public release for the production-readiness helper and example work.
- Bot API 10.0 code coverage is complete on `main` with documented architecture differences.
- Bot API 10.0 final audit found no known code coverage blockers.
- The root `aigram` package is now a compact quick-start facade; the full Bot API surface lives in `bot` and `telegram`.
- Public `main` consumer smoke passed after the root facade cleanup.
- [`docs/releases/v0.5.0.md`](releases/v0.5.0.md) contains the release notes for the current public release.
- Do not create future tags or GitHub Releases until the user explicitly approves release work.

Stage 98/99 outcome: **Bot API 9.6 code coverage is complete with documented architecture differences**. The final audit found wrappers for all 169 official Bot API methods and no missing fields in the audited high-impact object tables after adding `Message.giveaway`; Stage 99 resolved the final `setWebhook.certificate` upload blocker.

Stage 100 outcome: local release-readiness verification and manual-only smoke planning are documented. Stage 105/109/111 published `main` only after explicit user approval; the `v0.3.0` tag exists from the previous release line. Bot API 10.0 follow-up slices and final audit are complete after the May 8, 2026 upstream release.

Next phase:

1. Keep CI and local verification green for any follow-up fixes.
2. Keep Bot API 10.0 coverage docs current when Telegram publishes future Bot API releases.
3. Keep sensitive/state-changing smoke manual-only and fixture-first by default.
4. Keep [`PRE_V1_NOTES.md`](PRE_V1_NOTES.md) and [`CHANGELOG.md`](../CHANGELOG.md) current for breaking pre-v1 cleanup.
5. Use [`docs/plans/2026-05-14-production-readiness.md`](plans/2026-05-14-production-readiness.md) as the working plan for callback helpers, error taxonomy, group identity helpers, production examples, and transport-mode parity.
6. Monitor `v0.6.0` feedback and keep future release work behind explicit approval.

## v0.6 production-readiness candidate

Status: completed and released as `v0.6.0`.

Completed scope since `v0.5.0`:

- Typed callback helper layer in `callback`.
- Dispatcher routes for parsed typed callback data.
- Error taxonomy helpers for Telegram API errors, rate limits, migrations, forbidden/not-found responses, network errors, and context cancellation.
- Telegram actor and reply-target identity helpers for group/admin workflows.
- Production-style examples for inline panels, retry-aware sends, group admin identity, dry-run moderation, and webhook/long polling parity.
- Safer public example logs with masked numeric private IDs.
- Bot API update checklist and Bot API 10.0 lightweight freshness audit.
- Release-candidate checklist, local gates, direct public consumer smoke, and `v0.6.0` release notes.

Release status:

- Local release-candidate gates passed.
- Direct public consumer smoke for `main` passed after the public Go proxy path timed out.
- Main CI passed for the release-prep commits.
- `docs/releases/v0.6.0.md` contains the release notes.
- The tag and GitHub Release were created only after explicit maintainer approval.

Live smoke policy:

- Safe/read-only flows may be live-smoked only when explicitly useful and explicitly requested.
- Admin/state-changing flows require a dedicated test chat and explicit user confirmation.
- Payments, Passport, Business, Managed Bots, gifts, Stars, lifecycle methods, and webhook certificate upload require explicit confirmation.
- Destructive/admin flows must not be auto-smoked.

## Later

- Passport decryption helpers
  - Intentionally out of scope for the typed Bot API wrapper unless a future product decision adds them.
- Codegen
  - Consider generation only after the hand-written public API shape is proven.
  - Do not introduce codegen just to inflate method count.
