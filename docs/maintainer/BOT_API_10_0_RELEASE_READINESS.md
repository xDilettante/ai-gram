# Bot API 10.0 Release Readiness

## Status

Code coverage for Telegram Bot API 10.0 is complete with documented architecture differences. The final audit is recorded in [`../BOT_API_10_0_FINAL_AUDIT.md`](../BOT_API_10_0_FINAL_AUDIT.md).

The public repository exists and `main` has been published only after explicit user approval. The Bot API 10.0 package was published as the annotated `v0.4.0` tag and GitHub pre-release after explicit maintainer approval. The current post-cleanup public release is `v0.5.0`, published as a regular GitHub Release.

## Verification

Final local verification for the Bot API 10.0 audit:

- `go test ./bot`
- `scripts/check.sh`
- `go test -race ./bot ./dispatch ./middleware ./transport/longpoll ./transport/webhook ./internal/httpclient ./telegram`
- `go test -coverprofile=coverage.out ./...`
- `go tool cover -func=coverage.out` -> total `63.4%`

Public `main` consumer verification after the root facade cleanup:

- external temporary module outside this repository;
- `go get github.com/xDilettante/ai-gram@main`;
- `go test ./...` with root-package quick-start usage and direct `bot` / `telegram` advanced usage;
- `go doc github.com/xDilettante/ai-gram`;
- `go doc github.com/xDilettante/ai-gram/bot.SendPollParams`;
- `go doc github.com/xDilettante/ai-gram/telegram.ChatMember`.

Post-cleanup public `main` consumer verification after the `Config` / `LoginURL` pre-v1 rename:

- external temporary module outside this repository;
- `GOPROXY=direct go get github.com/xDilettante/ai-gram@main` -> `v0.4.1-0.20260510214756-47907d32d24e`;
- `go test ./...` with root quick-start usage through `aigram.Config` and advanced direct imports through `bot.Config` and `telegram.LoginURL`;
- `go doc github.com/xDilettante/ai-gram`;
- `go doc github.com/xDilettante/ai-gram/bot.Config`;
- `go doc github.com/xDilettante/ai-gram/telegram.LoginURL`.

Public tag verification after publishing:

- external temporary module outside this repository;
- `go get github.com/xDilettante/ai-gram@v0.4.0`;
- `go list -m github.com/xDilettante/ai-gram`;
- `go doc github.com/xDilettante/ai-gram`.

Current release preparation:

- `v0.5.0` packages the post-`v0.4.0` `Config` / `LoginURL` cleanup and public consumer smoke tooling;
- the release notes source is [`../releases/v0.5.0.md`](../releases/v0.5.0.md);
- GitHub Release `v0.5.0` is intended to be a regular release, not a pre-release.

Post-release public adoption sweep:

- GitHub Release `v0.4.0` is public, non-draft, and marked as pre-release;
- pkg.go.dev serves `github.com/xDilettante/ai-gram@v0.4.0` with HTTP 200;
- external consumer compile smoke passed against `github.com/xDilettante/ai-gram@v0.4.0`;
- README coverage badge and release-notes link text were aligned with the published release state.

Coverage evidence:

- total statement coverage: 63.4%;
- no known Bot API 10.0 method, request parameter, result object, update field, or high-impact message/type field blockers remain;
- live smoke remains manual-only for sensitive and state-changing flows.

## Architecture Differences

- `FileRef` / `FileUpload` instead of official `InputFile`: upload-capable fields use typed Go helpers and multipart behavior.
- `GetChatFullInfo` remains as a same-result pre-v1 alias for `GetChat`; `GetChat` returns the official `ChatFullInfo` result shape.
- `telegram.ChatMember` is an interface implemented by official `ChatMember*` variants.
- `CallbackQuery.Message` uses the official `MaybeInaccessibleMessage` shape.
- Live-photo URL inputs are intentionally rejected where the implementation requires `file_id` or multipart upload for the official live-photo payload.

## Manual-Only Smoke Areas

These flows must not be run automatically and require explicit user approval plus dedicated test assets/accounts:

- payments, invoices, paid media, Stars, gifts, business gifts, subscription invite links, Premium subscription gifts, refunds, and subscription edits;
- Passport data and Passport error reporting;
- Business APIs, business account mutation, business messages, business stories, suggested posts, direct messages, and business gifts;
- managed bot token and managed bot access methods;
- guest mode flows;
- reaction deletion methods;
- admin/destructive chat methods, including bans, restrictions, promotions, invite links, join requests, chat profile changes, leave chat, mass unpin, sender-chat moderation, and deletion methods;
- sticker set mutation methods;
- games requiring BotFather game setup;
- inline mode, prepared inline messages, Web App, Mini App, and client-specific features requiring BotFather/client setup;
- lifecycle `logOut` and `close`;
- `setWebhook` certificate upload and webhook state changes.

## Publication Result

Published after explicit maintainer approval:

- annotated tag: `v0.4.0`;
- tag target: `2ce948b`;
- GitHub Release: <https://github.com/xDilettante/ai-gram/releases/tag/v0.4.0>;
- release type: pre-release;
- public module install: `go get github.com/xDilettante/ai-gram@v0.4.0`.

Current regular public release:

- annotated tag: `v0.5.0`;
- GitHub Release: <https://github.com/xDilettante/ai-gram/releases/tag/v0.5.0>;
- release type: regular release;
- public module install: `go get github.com/xDilettante/ai-gram@v0.5.0`.
