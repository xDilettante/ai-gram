# Bot API 10.0 Release Readiness

## Status

Code coverage for Telegram Bot API 10.0 is complete with documented architecture differences. The final audit is recorded in [`../BOT_API_10_0_FINAL_AUDIT.md`](../BOT_API_10_0_FINAL_AUDIT.md).

The public repository exists and `main` has been published only after explicit user approval. The `v0.3.0` tag already exists, so Bot API 10.0 release preparation targets the next `v0.4.0` tag. No `v0.4.0` tag or GitHub Release has been created yet.

## Verification

Final local verification for the Bot API 10.0 audit:

- `go test ./bot`
- `scripts/check.sh`
- `go test -race ./bot ./dispatch ./middleware ./transport/longpoll ./transport/webhook ./internal/httpclient ./telegram`
- `go test -coverprofile=coverage.out ./...`
- `go tool cover -func=coverage.out`

Public `main` consumer verification after the root facade cleanup:

- external temporary module outside this repository;
- `go get github.com/xDilettante/ai-gram@main`;
- `go test ./...` with root-package quick-start usage and direct `bot` / `telegram` advanced usage;
- `go doc github.com/xDilettante/ai-gram`;
- `go doc github.com/xDilettante/ai-gram/bot.SendPollParams`;
- `go doc github.com/xDilettante/ai-gram/telegram.ChatMember`.

Coverage evidence:

- total statement coverage: 63.1%;
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

## Publication Plan

No automatic tag or GitHub Release creation is allowed.

If the user explicitly approves later:

1. verify `main` is clean and green;
2. create the local annotated `v0.4.0` tag;
3. push only that tag;
4. verify public `go get github.com/xDilettante/ai-gram@v0.4.0`;
5. create the GitHub Release from [`../releases/v0.4.0.md`](../releases/v0.4.0.md).
