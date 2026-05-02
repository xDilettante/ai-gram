# Bot API 9.6 Release Readiness

## Status

Local code coverage for Telegram Bot API 9.6 is complete with documented architecture differences. No repository publication, push, tag, GitHub Release, or GitHub repository creation has been performed for this local-only workstream.

The latest public release remains `v0.2.0` until the user explicitly approves publication work.

## Verification

Stage 100 local verification checks:

- `gofmt -w .`
- `bash -n scripts/*.sh`
- `go test ./...`
- `go vet ./...`
- `go list ./...`
- `git diff --check`
- `git status --short`

Coverage evidence from the Stage 98/99 audit:

- method audit: 169/169 official Bot API 9.6 wrappers are present;
- `setWebhook.certificate` blocker resolved with upload-only `FileUpload` multipart support;
- representative high-impact field audit found no remaining missing fields after the Stage 98 `Message.giveaway` correction.

## Architecture differences

- `FileRef` / `FileUpload` instead of official `InputFile`: upload-capable fields use typed Go helpers and multipart behavior. The upload-only `setWebhook.certificate` parameter accepts `FileUpload` and rejects file IDs or URLs.
- `GetChat` + `GetChatFullInfo` compatibility split: `GetChat` remains backward-compatible and returns `*telegram.Chat`; `GetChatFullInfo` decodes the official Bot API 9.6 `ChatFullInfo` result.
- Flat `ChatMember` compatibility strategy: official chat member variant fields are decoded into the existing flat `telegram.ChatMember` shape. Optional concrete `ChatMember*` variant structs remain a possible future refinement.
- `CallbackQuery.Message` + `MaybeMessage` compatibility strategy: accessible callback messages keep the existing `Message` field, while `MaybeMessage` preserves official `MaybeInaccessibleMessage` decoding.

## Manual-only smoke areas

These flows must not be run automatically and require explicit user approval plus dedicated test assets/accounts:

- payments, invoices, paid media, Stars, gifts, business gifts, subscription invite links, Premium subscription gifts, refunds, and subscription edits;
- Passport data and Passport error reporting;
- Business APIs, business account mutation, business messages, business stories, suggested posts, direct messages, and business gifts;
- managed bot token methods;
- admin/destructive chat methods, including bans, restrictions, promotions, invite links, join requests, chat profile changes, leave chat, mass unpin, sender-chat moderation, and deletion methods;
- sticker set mutation methods;
- games requiring BotFather game setup;
- inline mode, prepared inline messages, and client-specific inline features requiring BotFather/client setup;
- lifecycle `logOut` and `close`;
- `setWebhook` certificate upload.

## Safe smoke candidates

These are candidates for future manual or explicitly requested smoke checks. They are not run automatically in Stage 100.

- local Bot API `getMe` / `getWebhookInfo` checks;
- webhook basic message flow with a disposable bot/chat and safe logs;
- long polling basic update/reply flow with a disposable bot/chat;
- safe send/contact/location/venue/poll/dice flow;
- media group generated upload flow;
- `GetChatFullInfo` read-only check if a test chat is configured.

## Do-not-release gates

Do not publish, tag, release, or create a repository while any of these are true:

- tests fail;
- `go vet` fails;
- docs claim unsupported behavior;
- any token, secret, webhook URL, payment payload, business payload, Passport payload, managed bot token, or private message payload leak is found;
- uncommitted changes remain after the intended local commit;
- the user has not explicitly approved repository creation, push, tag, or GitHub Release.

## Publication plan

No automatic publication is allowed.

If the user explicitly approves later:

1. create or configure the GitHub repository;
2. push `main`;
3. create the local tag;
4. push only that tag;
5. verify public `go get`;
6. create the GitHub Release.
