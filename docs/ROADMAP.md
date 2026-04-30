# Roadmap

This roadmap is intentionally pragmatic: keep `ai-gram` useful and stable as a Go Telegram Bot API library without turning it into a large framework too early.

## v0.1 stabilization

- API polish
  - Review exported names and GoDoc for consistency.
  - Keep typed params/results and avoid `map[string]any` in public API.
  - Confirm token-safe error messages across all public methods.
  - Preserve compatibility of `ChatID`, `FileRef`, reply markup, reply parameters, and edit targets.
- Docs polish
  - Keep README concise and accurate.
  - Maintain `docs/API_COVERAGE.md`, `docs/MANUAL_TESTING.md`, `docs/DEPLOY_TESTING.md`, `docs/LIVE_SMOKE_MATRIX.md`, and `docs/RELEASE_CHECKLIST.md`.
  - Add a release checklist covering tests, vet, examples compilation, safe logs, and secret hygiene.
- Examples cleanup
  - Keep examples short and runnable.
  - Ensure examples read configuration only from env.
  - Keep webhook and long polling examples admin-protected by default.
  - Keep safe logs free of full message text, bot tokens, webhook secrets, and token-bearing URLs.
- Live smoke matrix
  - Maintain targeted smoke flows for local Bot API, long polling, webhook, media, callbacks, edit/delete, forward/copy, chat action, and chat info.
  - Keep destructive/admin methods out of automatic smoke.
  - Document when manual confirmation and dedicated test chats are required.
- Release checklist
  - `gofmt -w .`
  - `bash -n scripts/*.sh`
  - `go test ./...`
  - `go vet ./...`
  - `git diff --check`
  - Optional targeted live smoke only for safe flows.

## v0.2 Bot API coverage expansion

- Remaining send methods
  - `SendContact`
  - `SendLocation`
  - `SendVenue`
  - `SendDice`
  - `SendAnimation`
  - `SendVideoNote`
  - `SendSticker`
  - `SendMediaGroup`
- Polls
  - `SendPoll`
  - `StopPoll`
- Stickers
  - Basic sticker sending first.
  - Sticker set management only after the basic type surface is stable.
- Invite links
  - Create/edit/revoke invite links.
  - Keep admin/destructive behavior clearly documented.
- Join requests
  - Approve/decline methods.
  - Join request update helpers.

## v0.3 Advanced/admin features

- Promote/restrict/admin tools
  - `PromoteChatMember` and related admin methods.
  - `SetChatPermissions` and chat metadata setters.
  - Dedicated manual smoke docs for admin-only flows.
- Forum topics
  - Create/edit/close/reopen/delete topic methods.
  - General forum topic helpers.
- Bot commands/menu
  - Command scopes.
  - Command set/get/delete.
  - Menu button methods.
- Reactions
  - Message reaction send/set methods.
  - Reaction update/type coverage.

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
