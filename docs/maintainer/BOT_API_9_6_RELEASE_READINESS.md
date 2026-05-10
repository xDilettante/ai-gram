# Bot API 9.6 Release Readiness

## Status

This document is the historical Bot API 9.6 readiness record. It has been superseded by the Bot API 10.0 final audit in [`../BOT_API_10_0_FINAL_AUDIT.md`](../BOT_API_10_0_FINAL_AUDIT.md), current pre-v1 notes in [`../PRE_V1_NOTES.md`](../PRE_V1_NOTES.md), and the current release checklist in [`RELEASE_CHECKLIST.md`](RELEASE_CHECKLIST.md).

Code coverage for Telegram Bot API 10.0 is complete with documented architecture differences. The public repository exists and `main` has been published only after explicit user approval. The `v0.4.0` tag and GitHub pre-release were published after explicit maintainer approval.

## Verification

Stage 100 verification checks:

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

## Historical Architecture Differences

The notes below describe the Bot API 9.6 release-readiness state. Later pre-v1 cleanup changed some public shapes, including `GetChat`, `ChatMember`, and `CallbackQuery.Message`.

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

Do not tag or create a GitHub Release while any of these are true:

- tests fail;
- `go vet` fails;
- docs claim unsupported behavior or current upstream completeness that is not backed by the Bot API 10.0 final audit;
- any token, secret, webhook URL, payment payload, business payload, Passport payload, managed bot token, or private message payload leak is found;
- uncommitted changes remain after the intended local commit;
- the user has not explicitly approved the exact tag or GitHub Release action.

## Publication plan

No automatic tag or GitHub Release creation is allowed.

If the user explicitly approves later:

1. verify `main` is clean and green;
2. create the local tag;
3. push only that tag;
4. verify public `go get`;
5. create the GitHub Release.
