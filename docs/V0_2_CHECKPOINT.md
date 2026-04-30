# v0.2 Checkpoint

This checkpoint reviews the current `main` branch after the v0.2 coverage expansion published after v0.1.1. It is a release-readiness note, not a tag or release announcement.

## Implemented since v0.1.1

### Send and media methods

- `SendContact`
- `SendLocation`
- `SendVenue`
- `SendPoll`
- `StopPoll`
- `SendDice`
- `SendSticker`
- `SendAnimation`
- `SendVideoNote`
- `SendMediaGroup`

### Bot commands and menu

- `SetMyCommands`
- `DeleteMyCommands`
- `GetMyCommands`
- `SetChatMenuButton`
- `GetChatMenuButton`
- `SetMyDefaultAdministratorRights`
- Minimal typed support for command scopes, menu buttons, Web App menu info, and default administrator rights.

### Chat access and administration

- Invite link methods:
  - `ExportChatInviteLink`
  - `CreateChatInviteLink`
  - `EditChatInviteLink`
  - `RevokeChatInviteLink`
- Chat join request methods:
  - `ApproveChatJoinRequest`
  - `DeclineChatJoinRequest`
- `telegram.ChatJoinRequest` update decoding and dispatch helpers.
- Admin management methods:
  - `PromoteChatMember`
  - `SetChatAdministratorCustomTitle`
  - `SetChatPermissions`

### Smoke and documentation

- Targeted v0.2 send-method smoke script for contact/location/venue/poll/dice plus optional sticker/animation/video note checks.
- Targeted `SendMediaGroup` smoke script with generated upload fallback.
- Public docs language cleanup to English.
- Updated API coverage, live smoke matrix, manual testing notes, README examples, and changelog entries.

## Live-smoked flows

The safe v0.2 live smoke surface has been exercised without destructive/admin actions:

- `SendContact`
- `SendLocation`
- `SendVenue`
- `SendPoll`
- `StopPoll`
- `SendDice`
- `SendMediaGroup` generated upload fallback

Optional media sends for `SendSticker`, `SendAnimation`, and `SendVideoNote` remain supported by the smoke helper when matching media env is configured; missing optional media is treated as a skip, not a failure.

## Not live-smoked intentionally

The following areas are implemented with unit/httptest coverage but should not be run automatically because they change bot/chat/admin state or require elevated permissions:

- Bot commands/menu state:
  - `SetMyCommands`
  - `DeleteMyCommands`
  - `SetChatMenuButton`
  - `SetMyDefaultAdministratorRights`
- Invite links:
  - `ExportChatInviteLink`
  - `CreateChatInviteLink`
  - `EditChatInviteLink`
  - `RevokeChatInviteLink`
- Join requests:
  - `ApproveChatJoinRequest`
  - `DeclineChatJoinRequest`
- Admin management:
  - `PromoteChatMember`
  - `SetChatAdministratorCustomTitle`
  - `SetChatPermissions`

Manual checks for these methods require a dedicated test group/channel, a test user when relevant, explicit confirmation, and rollback of any created links, changed rights, or changed permissions.

## Remaining high-value Bot API areas

Potential next slices after this checkpoint:

- Chat management methods:
  - `setChatTitle`
  - `setChatDescription`
  - `setChatPhoto`
  - `deleteChatPhoto`
  - `leaveChat`
- Forum topic methods.
- Reactions.
- Inline mode basics:
  - inline query result types
  - `answerInlineQuery`
  - chosen inline result handling
- Remaining sticker set methods.
- Bot profile methods.
- Payments, Passport, Games, Business, Stars/gifts, and full codegen remain better suited for later milestones.

## Release recommendation

### Option A: stabilize and release v0.2.0 soon

Pros:

- The current scope is already a meaningful step beyond v0.1.1.
- The release contains practical new safe send methods, media group support, command/menu support, invite/join request support, and admin management coverage.
- Safe v0.2 live smoke already covers the highest-value non-admin send flows, including `SendMediaGroup`.
- Admin/state-changing methods are intentionally covered by unit/httptest tests and documented as manual-only.
- Cutting v0.2.0 now gives downstream users a stable tag for the current expanded public API.

Cons:

- Some high-value Bot API areas remain missing, especially forum topics, reactions, inline mode, sticker set management, and chat metadata setters.
- Bot commands/menu and admin methods have not been live-smoked because they mutate real bot/chat state.
- Optional sticker/animation/video note live smoke depends on configured media fixtures and may still be skipped in default environments.

Recommendation: choose Option A if the goal is to publish a coherent v0.2 milestone with the current expanded coverage and then continue incremental v0.3 work.

### Option B: continue coverage before v0.2.0

Pros:

- A broader v0.2.0 could include more of the remaining everyday Bot API surface.
- Chat management, forum topics, reactions, and inline mode could make the milestone feel more complete.

Cons:

- The release boundary will keep moving and delay a usable v0.2 tag.
- Additional admin/forum/inline methods add more stateful behavior and more manual-only verification surface.
- Bigger batches increase review and regression risk.

Suggested next slices if choosing Option B:

- Chat management methods: `setChatTitle`, `setChatDescription`, `setChatPhoto`, `deleteChatPhoto`, `leaveChat`.
- Forum topic methods.
- Reactions.
- Inline mode basics.
- Remaining sticker set methods.

Recommendation: choose Option B only if one of those slices is required before users should consume a v0.2 tag.

## Final recommendation

The v0.2.0 release flow is complete: main, tag, public `go get`, and GitHub Release are published. Follow-up planning has moved to [`docs/V0_3_PLAN.md`](V0_3_PLAN.md).
