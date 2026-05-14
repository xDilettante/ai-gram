# Production Readiness Plan

Date: 2026-05-14

This document records the current improvement plan for `ai-gram` after the regular `v0.5.0` public release. It is the working memory for what the project needs next, what is already done, and what remains.

## Goal

Make `ai-gram` more useful for real production bots without turning it into a heavyweight framework. The next improvements should reduce repeated application code around inline panels, group identity, retry behavior, and deployment-mode examples.

## Current State

- `v0.5.0` is published as a regular GitHub Release, not a pre-release.
- `main` has a local unpublished commit: `38a426e Use buffered multipart for edit media uploads`.
- Telegram Bot API 10.0 code coverage is complete with documented architecture differences.
- The project has separate library layers:
  - `bot` for outgoing Bot API calls;
  - `telegram` for data contracts;
  - `dispatch` for routing;
  - `middleware` for reusable middleware;
  - `transport/longpoll` and `transport/webhook` for update intake.
- Basic callback routing exists through exact callback data matching.
- `errors.APIError` already exposes Telegram response parameters such as `retry_after` and `migrate_to_chat_id`.
- `transport/longpoll` already has polling retry backoff.
- Public examples exist, but the higher-value production patterns are still thin.

## Done

- Published regular release `v0.5.0`.
- Added public consumer smoke tooling and workflow.
- Cleaned pre-v1 public names:
  - `aigram.Config`;
  - `bot.Config`;
  - `telegram.LoginURL`.
- Added `docs/PRE_V1_NOTES.md`.
- Kept live smoke and destructive Telegram flows manual-only.
- Added root-package facade documentation and current API split.
- Added buffered multipart behavior for edit media uploads in the local unpublished commit.

## Needed

1. Callback and inline keyboard helper layer.
2. Error taxonomy and retry/rate-limit helpers.
3. Group identity helpers.
4. Production examples that prove the helper layers.
5. Cleaner long polling/webhook parity for application structure.
6. Bot API compatibility discipline for future Telegram releases.

## Priority Order

### 1. Callback Layer

Add a small `callback` package for typed `callback_data` construction and parsing.

Required capabilities:

- compact format within Telegram's 64-byte `callback_data` limit;
- namespace and action fields;
- optional ID field;
- optional page field for pagination;
- optional expiry or TTL field;
- safe parser with typed errors;
- helpers for confirm/cancel flows;
- helpers for previous/next pagination buttons;
- dispatch predicates or adapters for namespace/action matching.

Acceptance criteria:

- callback data round-trips through encode/parse tests;
- invalid input returns typed parse errors;
- too-long payloads are rejected before sending;
- expired payloads can be detected;
- examples use typed callbacks instead of ad hoc string constants where appropriate.

### 2. Error Taxonomy And Retry Helpers

Build on the existing `errors.APIError` instead of replacing it.

Required capabilities:

- helper functions for Telegram API errors;
- helper functions for `retry_after`;
- helper functions for chat migration;
- helper functions for forbidden/not found style responses;
- network/transient error classification where it can be done safely;
- no hidden automatic retries by default.

Possible public helpers:

- `errors.IsAPIError(err error)`;
- `errors.AsAPIError(err error)`;
- `errors.IsRateLimited(err error)`;
- `errors.RetryAfter(err error)`;
- `errors.MigrateToChatID(err error)`;
- `errors.IsForbidden(err error)`;
- `errors.IsNotFound(err error)`;
- `errors.IsNetworkError(err error)`.

Acceptance criteria:

- existing `APIError` behavior remains compatible;
- helpers work through wrapped errors;
- tests cover Telegram error payloads, network errors, context cancellation, and unrelated errors.

### 3. Group Identity Helpers

Add helper functions around `telegram.Update`, `telegram.Message`, and callback updates. This should not become a moderation framework.

Required capabilities:

- identify the actor user when one exists;
- identify sender chat and anonymous-admin style messages;
- identify reply target user;
- expose safe helper names for group/admin examples;
- keep policy decisions in application code or middleware.

Acceptance criteria:

- helpers are documented;
- tests use synthetic update fixtures;
- middleware and examples can consume helpers without duplicating identity extraction logic.

### 4. Production Examples

Use examples to validate helper APIs.

Candidate examples:

- inline panel with typed callbacks, pagination, and confirm/cancel;
- group bot with admin commands;
- moderation action skeleton with explicit safety warnings;
- webhook service with graceful shutdown;
- long polling service with graceful shutdown;
- retry/rate-limit aware sender.

Acceptance criteria:

- examples compile through `scripts/check.sh`;
- examples avoid raw private IDs in public logs;
- examples avoid live destructive behavior by default;
- examples do not require secrets beyond documented environment variables.

### 5. Webhook And Long Polling Parity

Keep both modes first-class without hiding transport behavior.

Required capabilities:

- clear app structure that can swap long polling and webhook intake;
- shared dispatcher and middleware setup;
- graceful shutdown pattern;
- examples for systemd/nginx live under maintainer docs or production examples.

Acceptance criteria:

- users can compare long polling and webhook examples without rewriting business handlers;
- docs explain mode-specific tradeoffs and Telegram constraints.

### 6. Bot API Compatibility Discipline

Do not introduce code generation until the public API shape is more proven.

Required capabilities now:

- checklist for Telegram Bot API updates;
- audit workflow for new methods, fields, and result types;
- tests for request encoding, result decoding, and JSON compatibility;
- changelog and release-note discipline for breaking pre-v1 changes.

Future option:

- schema-assisted generation or audit tooling if manual tracking becomes the bottleneck.

## Next Slice

Start with the callback layer.

Planned first implementation slice:

1. Push or otherwise resolve the current local `38a426e` commit before stacking more work.
2. Add `callback` package with encode/parse and typed errors.
3. Add tests for roundtrip, invalid input, expiry, and length limit.
4. Add confirm/cancel and pagination helpers.
5. Update one small example first, likely `examples/04_inline_keyboard`.
6. Run `scripts/check.sh`.
7. Commit the slice atomically.

## Not Now

- No new release until the next meaningful feature set is complete and verified.
- No automatic live Telegram smoke for destructive/admin/payment/passport/business flows.
- No broad framework runner before callback/error/identity helpers prove the API shape.
- No code generation until the manual public API shape stabilizes further.

## Verification Commands

Use the standard verification wrapper:

```bash
scripts/check.sh
```

For callback-specific work, also use focused package tests once the package exists:

```bash
go test ./callback ./dispatch ./examples/...
```
