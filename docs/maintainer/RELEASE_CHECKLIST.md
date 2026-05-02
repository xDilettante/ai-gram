# Release Checklist

This checklist is for preparing a future `v0.1.0` tag. It does not create the tag or publish a release by itself.


## Current pre-v0.1 status

Last verified: 2026-04-30, Stage 30 pre-release verification.

- [x] Static checks passed: `gofmt -w .`, `bash -n scripts/*.sh`, `go test ./...`, `go vet ./...`, `git diff --check`, and clean `git status --short` before live/docs updates.
- [x] Docs checked: README, API coverage, roadmap, manual testing, deploy testing, live smoke matrix, and this checklist are present and aligned with current scope.
- [x] Examples compile through `go test ./...` and keep admin-only access control enabled by default.
- [x] API coverage checked: raw `Bot.Token()` accessor is removed; chat info, chat action/pin, and moderation methods are listed; destructive/admin methods are marked as not auto-smoked.
- [x] Safe live smoke subset passed: local Bot API smoke, webhook deploy, access panel/status, chat info, edit text/reply markup, generated-document caption edit/delete, delete bot-created message, reply/sendChatAction, and forward/copy.
- [x] Security/secrets checks passed: no tracked token-like strings found, `.env.local` and `.deploy/generated.env` are ignored, raw token accessor search is clean, and logs/reports used safe excerpts only.
- [x] Known limitations remain documented; no `v0.1.0` tag or GitHub release has been created.

## Pre-release checks

- Confirm `docs/API_COVERAGE.md` matches current implemented methods.
- Confirm `docs/ROADMAP.md` reflects the next planned stage.
- Review exported API for obvious pre-v0.1 naming issues before the tag.
- Token exposure decision resolved before v0.1: raw bot token is intentionally not exposed through public `Bot` methods; use `GetMe` for bot identity and redacted `fmt.Stringer` output for diagnostics.
- Confirm no new Bot API method was added without unit/httptest coverage.
- Confirm all public exported declarations have useful GoDoc comments.
- Confirm all blocking/network operations accept `context.Context` where applicable.

## Tests

Run from the repository root:

```bash
gofmt -w .
bash -n scripts/*.sh
go test ./...
go vet ./...
git diff --check
```

Do not release if any command fails or if `git status --short` contains unintended changes.

## Docs

- README describes the project as early-stage and does not promise full Bot API coverage.
- README links to:
  - `docs/API_COVERAGE.md`
  - `docs/ROADMAP.md`
  - `docs/MANUAL_TESTING.md`
  - `docs/maintainer/DEPLOY_TESTING.md`
  - `docs/maintainer/LIVE_SMOKE_MATRIX.md`
  - `docs/maintainer/RELEASE_CHECKLIST.md`
- Manual testing docs describe access control/admin-only mode.
- Deploy testing docs describe deep-link smoke panels, role-specific tokens, separate Bot API host support, and TUN/Xray/vk1 caveats.
- Coverage docs classify destructive/admin methods and methods that should not be smoke-tested automatically.

## Examples

- `go test ./...` compiles all `examples/*` packages without runtime env.
- Examples read env only inside `main` or runtime helpers, not package init.
- Examples do not log bot tokens, webhook secrets, token-bearing URLs, `.env.local`, or full private message text.
- `examples/webhook_server` defaults to admin-only access control.
- `examples/inline_longpoll` uses access control and runtime access commands.
- Webhook safe logs include action/update/chat/message IDs where useful and avoid full text/callback payloads except known demo callback data.

## Live smoke safe flows

Use `docs/maintainer/LIVE_SMOKE_MATRIX.md` as the source of truth. Recommended safe subset before `v0.1.0`:

- local Bot API smoke;
- webhook `/start`;
- access panel status/open/close with immediate close;
- edit text/reply markup flow;
- caption flow using generated document;
- delete flow limited to bot-created test message;
- reply flow;
- forward/copy flow;
- sendChatAction flow.

Record only safe log excerpts such as `update_id`, `action`, `matched`, `chat_id`, `from_user_id`, and `message_id`.

## Security/secrets

- No real token or secret is committed.
- `.env.local`, `.deploy/`, generated env files, SSH keys, and private keys remain ignored.
- No docs or examples contain token-bearing URLs or full `/bot<TOKEN>/...` endpoints.
- Public API does not expose raw bot token accessors; diagnostics use redacted string output.
- Telegram notifications and final reports do not include secrets or full private message text.
- Webhook secret matches both `SetWebhook` and `webhook.Config` in examples/deploy env.
- Destructive/admin methods are not run automatically.

## Versioning/tag

- Decide final module version and tag name, expected `v0.1.0`.
- Ensure release commit is clean and all checks passed.
- Create an annotated tag only after the release checklist is complete:

```bash
git tag -a v0.1.0 -m "v0.1.0"
```

- Do not create the tag or GitHub release during ordinary stabilization tasks unless explicitly requested.

## Known limitations

- Full Telegram Bot API coverage is not implemented.
- No code generation/openapi pipeline is used.
- Media groups, thumbnails, animation/video note sending, polls, stickers, invite links, join requests, payments, passport, games, inline mode, WebApp/LoginUrl, bot commands/menu, forum topics, business APIs, Stars/gifts, and many admin setters remain deferred.
- Some Telegram types are intentionally minimal and should be expanded only when required by implemented methods or update handling.
- Live smoke depends on real credentials and the local/network environment; local TUN/Xray and remote local Bot API routing can change observed behavior.

## Do not release if

- Any required local check fails.
- Examples do not compile with `go test ./...`.
- A secret or token-bearing URL appears in tracked files, logs, or reports.
- README or coverage docs claim methods that are not implemented.
- New public API was added without GoDoc and tests.
- Destructive/admin smoke was run against a non-test chat or without explicit confirmation.
- Access control examples default to public access.
