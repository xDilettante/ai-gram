# Next Release Candidate Checklist

This document tracks the next release-candidate pass for the post-`v0.5.0` work on `main`.

It is a preparation checklist only. Do not create tags, GitHub Releases, or release notes uploads from this document without explicit maintainer approval.

## Candidate Shape

- Current latest release: `v0.5.0`.
- Candidate version class: minor release, likely `v0.6.0`.
- Reason: `main` contains new public helper APIs, dispatcher APIs, production examples, and maintainer compatibility workflow since `v0.5.0`.
- Current release-candidate base: commits after `v0.5.0` through the commit that adds this document.

## Included Since `v0.5.0`

### Public API

- `callback` package for compact typed callback data, parsing, expiry checks, pagination helpers, confirm/cancel actions, and callback button construction.
- Dispatcher integration for typed callback routing:
  - `dispatch.CallbackAction`;
  - `Dispatcher.OnCallbackAction`;
  - `Dispatcher.OnCallbackActionFunc`;
  - `dispatch.CallbackDataHandler`.
- Error taxonomy helpers on top of `errors.APIError`, including rate-limit, migration, forbidden/not-found, network, cancellation, and deadline classification.
- Telegram identity helpers:
  - `telegram.Actor`;
  - message, callback, and update actor helpers;
  - anonymous admin detection;
  - reply target helpers.

### Examples

- `examples/07_inline_panel` for typed callbacks, pagination, and confirm/cancel flows.
- `examples/08_retry_sender` for explicit retry/rate-limit-aware sending.
- `examples/09_group_admin` for safe group/admin identity commands.
- `examples/10_moderation_skeleton` for dry-run moderation reporting and admin previews.
- `examples/11_transport_parity` for shared handlers running through long polling or webhook intake.
- Public example logs now mask numeric private IDs more consistently.

### Compatibility And Maintainer Workflow

- Production readiness plan for the post-`v0.5.0` workstream.
- Bot API update checklist for future Telegram Bot API release audits.
- Bot API 10.0 lightweight freshness audit dated 2026-05-14; no newer official Bot API release or missing official method wrappers were observed.

### Internal Behavior

- Buffered multipart behavior for edit media uploads.
- No code generation added.
- No automatic destructive, payment, passport, managed-token, sticker-set, webhook-certificate, or lifecycle live smoke added.

## Release Candidate Gates

Run these gates before tagging a release candidate or final release.

## Gate Results

### 2026-05-14 Local Gates

Passed on `main` after adding this checklist:

- `scripts/check.sh`;
- `go test -race ./bot ./callback ./dispatch ./errors ./middleware ./transport/longpoll ./transport/webhook ./internal/httpclient ./telegram`;
- `go test -coverprofile=coverage.out ./...`;
- `go tool cover -func=coverage.out`;
- `git diff --check`;
- `git status --short`.

Coverage summary:

- total statement coverage: `60.0%`.

Post-run cleanup:

- `coverage.out` was removed after collecting the summary.
- Working tree was clean before recording this result.

### 2026-05-14 Public Consumer Smoke

Direct public consumer smoke passed for `main`:

```bash
AIGRAM_CONSUMER_DIRECT=1 scripts/smoke_public_consumer.sh
```

Resolved module version:

```text
github.com/xDilettante/ai-gram v0.5.1-0.20260514192257-29f0f64349df
```

The default proxy-backed run was attempted first and failed with a TLS handshake timeout from `proxy.golang.org`. The fallback direct mode is the documented path for stale or unavailable public proxy checks.

### Local Gates

```bash
scripts/check.sh
go test -race ./bot ./callback ./dispatch ./errors ./middleware ./transport/longpoll ./transport/webhook ./internal/httpclient ./telegram
go test -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
git diff --check
git status --short
```

Expected result:

- all commands pass;
- `git status --short` is clean except intentional release-prep edits before commit;
- no generated `coverage.out` is committed.

### Public Consumer Gates

After the release-candidate branch or commit is pushed:

```bash
scripts/smoke_public_consumer.sh
```

If the public Go proxy is stale:

```bash
AIGRAM_CONSUMER_DIRECT=1 scripts/smoke_public_consumer.sh
```

Also run the manual GitHub Actions workflow `Public consumer smoke` against `main` before tagging. After a tag exists, run it again against that tag.

### Documentation Gates

- `CHANGELOG.md` has a clear unreleased section for the candidate.
- `docs/API_COVERAGE.md` matches the implementation and does not claim unsupported behavior.
- `docs/BOT_API_10_0_LIGHTWEIGHT_AUDIT_2026_05_14.md` remains current enough for the release date, or a newer audit exists.
- `docs/maintainer/BOT_API_UPDATE_CHECKLIST.md` has been followed if Telegram published a newer Bot API release.
- `docs/maintainer/LIVE_SMOKE_MATRIX.md` still marks sensitive and destructive areas as manual-only.
- Release notes exist under `docs/releases/` before a final release is created.
- README examples list matches the shipped examples.

### Security Gates

- No real bot token, webhook secret, private chat ID, payment payload, Passport payload, managed bot token, SSH detail, or token-bearing URL is committed.
- No public example logs print raw numeric private IDs where masking is expected.
- No full `/bot<TOKEN>/...` endpoint is present in tracked docs, examples, scripts, or tests except safe placeholder text.
- No `.env.local`, `.deploy/`, generated env file, private key, or local smoke artifact is staged.
- No live destructive/admin/payment/passport/business/managed-token/webhook-certificate smoke was run without explicit maintainer approval.

### API Compatibility Gates

- New exported APIs have GoDoc.
- New public helpers have focused unit tests.
- Existing `v0.5.0` public APIs were not renamed or removed without explicit breaking-change notes.
- Any pre-v1 breaking change is called out in `CHANGELOG.md` and, if needed, `docs/PRE_V1_NOTES.md`.
- Blocking operations use caller-provided `context.Context`.

## Candidate Release Notes Draft

Use this as the starting point for `docs/releases/v0.6.0.md` if the candidate proceeds as `v0.6.0`.

### Added

- Typed callback data helpers and dispatcher routing for callback actions.
- Error classification helpers for Telegram API errors, rate limits, migrations, forbidden/not-found responses, network errors, and context cancellation.
- Telegram actor and reply-target identity helpers for group/admin workflows.
- Production-style examples for inline panels, retry-aware sending, group admin identity, dry-run moderation, and transport parity.
- Maintainer Bot API update checklist and a lightweight Bot API 10.0 freshness audit.

### Changed

- Public examples mask numeric private IDs in logs more consistently.
- Edit media multipart handling uses buffered multipart behavior.

### Compatibility

- No newer official Bot API release than 10.0 was observed during the 2026-05-14 lightweight audit.
- No missing official Bot API method wrappers were found in the method-level audit.
- No code generation was introduced.

### Manual-Only Areas

- Destructive admin flows, payments, Passport, managed tokens, sticker-set mutation, webhook certificates, and lifecycle methods remain manual-only.

## Do Not Proceed If

- `scripts/check.sh` fails.
- Public consumer smoke fails.
- GitHub Actions for `main` is failing.
- The official Bot API page shows a newer release that has not been audited.
- `CHANGELOG.md`, release notes, and `docs/API_COVERAGE.md` disagree.
- A secret or token-bearing URL appears in tracked files.
- The maintainer has not explicitly approved tagging or publishing.
