# Bot API Update Checklist

This checklist is the maintainer workflow for keeping ai-gram aligned with new Telegram Bot API releases.

Use it when Telegram publishes a new Bot API version, changes an existing method/object, or clarifies request/response behavior. The goal is disciplined manual coverage before adding code generation. Do not treat this file as a copy of the upstream specification; always compare against the official Telegram Bot API documentation at update time.

## Scope

This workflow covers:

- new Bot API methods;
- new request parameters on existing methods;
- changed parameter requirements or accepted value ranges;
- new result fields on Telegram objects;
- new update/message/service-message variants;
- changed response result shapes such as `Message` or `true`;
- multipart upload behavior, including upload-only fields;
- Bot API error response parameters such as `retry_after` and `migrate_to_chat_id`;
- examples, docs, changelog, and release notes affected by the update.

This workflow does not authorize live destructive smoke, release tagging, GitHub Releases, or broad public API rewrites.

## Sources To Compare

Use these sources in order:

1. Official Telegram Bot API changelog and method/object documentation.
2. Current ai-gram implementation in `bot/`, `telegram/`, `dispatch/`, `middleware/`, and `transport/`.
3. Current inventory in [`../API_COVERAGE.md`](../API_COVERAGE.md).
4. Sensitive/live-smoke classification in [`LIVE_SMOKE_MATRIX.md`](LIVE_SMOKE_MATRIX.md).
5. Release gate in [`RELEASE_CHECKLIST.md`](RELEASE_CHECKLIST.md).

When a source disagrees with the code, update the code or document the intentional architecture difference. Do not let `API_COVERAGE.md` claim support that is not implemented and tested.

## Audit Steps

### 1. Record The Upstream Delta

- Identify the Bot API version or upstream change date.
- List every new method.
- List every changed method parameter.
- List every new or changed object field.
- List every new update type, message subtype, service message, or enum-like value.
- List every changed response shape, especially methods that can return `Message` or `true`.
- List every new sensitive, destructive, payment, passport, business, managed-token, webhook, or admin capability.

Keep this as a short audit document when the delta is large, for example `docs/BOT_API_X_Y_FINAL_AUDIT.md`.

### 2. Map Public API Placement

Use the repository layers consistently:

- outgoing Bot API methods and request params belong in `bot/`;
- Telegram JSON data contracts belong in `telegram/`;
- update routing helpers belong in `dispatch/`;
- reusable access/recovery/timeout behavior belongs in `middleware/`;
- long polling and inbound webhook behavior belongs in `transport/`;
- runnable usage patterns belong in `examples/`;
- maintainer-only live or deploy workflow belongs in `docs/maintainer/` and `scripts/`.

Avoid adding framework-level globals or hidden background behavior to support a Bot API update.

### 3. Implement Methods And Request Params

For each new or changed method:

- add or update the public method wrapper;
- add or update the params struct with JSON tags that match Telegram names;
- preserve caller-provided `context.Context`;
- support optional fields with `omitempty` where Telegram treats omission differently from zero values;
- validate only rules that are stable and valuable locally;
- preserve server-side validation for Telegram-owned policy rules that may change;
- support `business_connection_id`, `message_thread_id`, `direct_messages_topic_id`, reply parameters, and reply markup consistently where Telegram accepts them;
- use existing `ChatID`, `FileID`, `FileURL`, `FileUpload`, and result helper patterns instead of inventing parallel types.

If Telegram accepts multiple result shapes, add or reuse a typed decoder instead of relying on `interface{}`.

### 4. Implement Telegram Objects And Updates

For each new or changed object:

- add JSON-compatible fields to the matching `telegram` type;
- keep field names idiomatic Go while preserving exact JSON tags;
- use pointers or slices where absence is meaningful;
- add helper methods only when they remove repeated application code or match an existing local pattern;
- add synthetic decode fixtures for new update/message/service-message variants;
- avoid policy decisions in data-contract helpers.

New object fields are usually compatible. Renames, removals, or type changes are breaking pre-v1 changes and must be called out in `CHANGELOG.md`.

### 5. Implement Multipart And File Semantics

For every new file-like parameter:

- confirm whether Telegram accepts file IDs, URLs, uploads, or upload-only values;
- use `FileID`, `FileURL`, and `FileUpload` only where Telegram accepts them;
- reject unsupported file reference modes before sending;
- use deterministic `attach://` names in multipart requests;
- add tests for JSON mode and multipart mode where both are supported;
- add tests for upload-only rejection of file IDs and URLs where applicable.

Do not log local file paths, token-bearing URLs, or large payloads.

### 6. Classify Safety And Live Testing

Before adding examples or smoke scripts, classify the new capability:

- safe read;
- safe send to a dedicated test chat;
- sensitive;
- state-changing;
- destructive/admin;
- payment/Stars/gift/refund;
- passport;
- business or managed-token;
- webhook lifecycle;
- local Bot API lifecycle.

Update [`LIVE_SMOKE_MATRIX.md`](LIVE_SMOKE_MATRIX.md) when the new capability changes manual/live testing boundaries. Destructive, payment, passport, managed-token, sticker-set mutation, webhook-certificate, and lifecycle checks stay manual-only unless the maintainer explicitly approves a targeted run.

## Required Tests

Every meaningful Bot API update should have focused tests before documentation claims support.

Minimum coverage:

- request encoding for new methods and parameters;
- result decoding for new response shapes;
- JSON decoding for new Telegram object fields;
- multipart encoding for new upload fields;
- validation/rejection tests for unsupported local file reference modes;
- API error behavior when a method has special response parameters;
- context cancellation or timeout behavior for new blocking operations;
- dispatcher predicate or update routing tests for new update variants;
- examples compile through `scripts/check.sh` when examples are touched.

Prefer `httptest` and synthetic fixtures over live Telegram checks. Live checks should use dedicated test bots/chats and must not print tokens, private IDs, webhook secrets, payment payloads, Passport payloads, managed bot tokens, or token-bearing URLs.

## Documentation Updates

For each update slice, revise the smallest useful set:

- `docs/API_COVERAGE.md` for implemented methods, object areas, tests, and intentional gaps;
- `CHANGELOG.md` for user-facing additions, changes, and breaking pre-v1 cleanup;
- `README.md` only when public workflow or headline capability changes;
- `docs/PRE_V1_NOTES.md` when public naming or compatibility changes;
- `docs/maintainer/LIVE_SMOKE_MATRIX.md` when manual/live safety boundaries change;
- release notes under `docs/releases/` when preparing a release;
- examples and `docs/MANUAL_TESTING.md` when intended usage changes.

Keep documentation factual. Do not describe future work as implemented support.

## Completion Gate

Before merging or releasing a Bot API update slice:

```bash
scripts/check.sh
git diff --check
git status --short
```

For large method/object updates, also run targeted package tests, for example:

```bash
go test ./bot ./telegram ./dispatch ./transport/longpoll ./transport/webhook
```

The slice is complete only when:

- all new upstream changes are implemented or explicitly listed as intentional gaps;
- `docs/API_COVERAGE.md` matches the implementation;
- new public APIs have GoDoc;
- request/result/object tests cover the new behavior;
- sensitive and destructive areas are classified;
- no real secrets, private IDs, payment data, Passport data, managed bot tokens, or token-bearing URLs are committed;
- `CHANGELOG.md` records user-visible changes;
- release publication remains unstarted unless explicitly approved.

## When To Consider Generation

Stay manual for now. Revisit schema-assisted generation only if at least one of these becomes true:

- manual audits repeatedly miss fields or methods;
- request/result tests become mostly mechanical;
- Telegram publishes a stable machine-readable schema that matches production behavior;
- generated code can preserve the existing public API shape and documentation quality.

Generated code must not bypass tests, security rules, or manual safety classification.
