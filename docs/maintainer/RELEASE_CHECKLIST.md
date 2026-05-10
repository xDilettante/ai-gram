# Release Checklist

This checklist is for preparing a future `v0.4.0` tag. It does not create the tag or publish a release by itself.

## Current Status

Last verified: May 10, 2026, Bot API 10.0 final audit.

- [x] Bot API 10.0 code coverage is complete with documented architecture differences.
- [x] Final Bot API 10.0 audit is recorded in [`../BOT_API_10_0_FINAL_AUDIT.md`](../BOT_API_10_0_FINAL_AUDIT.md).
- [x] README, API coverage, roadmap, release notes, and live smoke matrix are aligned with the Bot API 10.0 status.
- [x] Sensitive and state-changing live smoke remains manual-only.
- [x] No `v0.4.0` tag or GitHub Release has been created.

## Required Local Checks

Run from the repository root before publishing:

```bash
scripts/check.sh
go test -race ./bot ./dispatch ./middleware ./transport/longpoll ./transport/webhook ./internal/httpclient ./telegram
go test -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
git status --short
```

Do not release if any command fails or if `git status --short` contains unintended changes.

## Documentation Gates

- Confirm `README.md` describes Bot API 10.0 coverage accurately.
- Confirm `docs/API_COVERAGE.md` matches current implemented methods and intentional architecture differences.
- Confirm `docs/ROADMAP.md` reflects the current release stage.
- Confirm `docs/releases/v0.4.0.md` is the release notes source for the tag.
- Confirm `docs/maintainer/LIVE_SMOKE_MATRIX.md` marks destructive, sensitive, and state-changing flows as manual-only.

## API And Test Gates

- Confirm no new Bot API method was added without unit or httptest coverage.
- Confirm all public exported declarations have useful GoDoc comments.
- Confirm all blocking/network operations accept or use `context.Context`.
- Confirm token redaction tests still cover API errors, invalid JSON, HTTP errors, and context cancellation paths where relevant.
- Confirm examples compile through `scripts/check.sh`.

## Security Gates

- No real token, webhook secret, private chat ID, payment payload, Passport payload, managed bot token, SSH detail, or token-bearing URL is committed.
- `.env.local`, `.deploy/`, generated env files, SSH keys, and private keys remain ignored.
- No docs or examples contain full `/bot<TOKEN>/...` endpoints.
- Logs and final reports do not include secrets or full private message text.
- Destructive/admin live smoke was not run automatically.

## Manual-Only Live Smoke Areas

Use [`LIVE_SMOKE_MATRIX.md`](LIVE_SMOKE_MATRIX.md) as the source of truth. The following areas require explicit approval and dedicated test assets:

- payments, Stars, gifts, paid media, refunds, and subscription flows;
- Passport data and Passport error reporting;
- Business APIs and business account mutation;
- managed bot token/access methods;
- guest mode flows;
- reaction deletion;
- admin/destructive chat methods;
- sticker set mutation;
- games, inline mode, Web App, Mini App, and prepared inline flows;
- lifecycle methods such as `logOut` and `close`;
- webhook certificate upload and webhook state changes.

## Versioning And Tag

- Expected release tag: `v0.4.0`.
- Ensure the intended release commit is clean and all checks passed.
- Create an annotated tag only after explicit maintainer approval:

```bash
git tag -a v0.4.0 -m "Release v0.4.0"
```

- Push only that tag:

```bash
git push origin v0.4.0
```

- Verify public module installation after the tag is available:

```bash
go get github.com/xDilettante/ai-gram@v0.4.0
```

- Create the GitHub Release using [`../releases/v0.4.0.md`](../releases/v0.4.0.md) as the release notes source.

Do not run `git push --tags`, create unrelated tags, or create a GitHub Release before explicit approval.

## Do Not Release If

- Any required local check fails.
- Examples do not compile through `scripts/check.sh`.
- A secret or token-bearing URL appears in tracked files, logs, or reports.
- README or coverage docs claim behavior that is not implemented.
- New public API was added without GoDoc and tests.
- Destructive/admin smoke was run against a non-test chat or without explicit confirmation.
- Access control examples default to public access.
