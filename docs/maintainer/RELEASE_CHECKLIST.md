# Release Checklist

This checklist records completed regular releases and remains useful as the template for future pre-v1 releases. The `v0.6.0` release-candidate preparation pass is tracked in [`NEXT_RELEASE_CANDIDATE.md`](NEXT_RELEASE_CANDIDATE.md).

## Current Status

Last verified: May 14, 2026, Bot API 10.0 lightweight freshness audit, release-candidate local gates, direct public `main` consumer smoke, and public `v0.6.0` release preparation.

- [x] Bot API 10.0 code coverage is complete with documented architecture differences.
- [x] Final Bot API 10.0 audit is recorded in [`../BOT_API_10_0_FINAL_AUDIT.md`](../BOT_API_10_0_FINAL_AUDIT.md).
- [x] README, API coverage, roadmap, release notes, and live smoke matrix are aligned with the Bot API 10.0 status.
- [x] Public `main` consumer smoke passed after the root facade cleanup and after the later `Config` / `LoginURL` pre-v1 rename.
- [x] Sensitive and state-changing live smoke remains manual-only.
- [x] Annotated `v0.4.0` tag and GitHub pre-release were created after explicit maintainer approval.
- [x] `v0.5.0` is prepared as the regular public release for the post-`v0.4.0` cleanup.
- [x] `v0.6.0` is prepared as the regular public release for the post-`v0.5.0` production-readiness work.

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

Before tagging a future release, also verify the public branch from a temporary module outside this repository:

```bash
scripts/smoke_public_consumer.sh
```

If the public Go proxy has not yet observed the latest pushed commit, repeat the branch verification with `AIGRAM_CONSUMER_DIRECT=1 scripts/smoke_public_consumer.sh`.

The same check is available from GitHub Actions as the manually triggered `Public consumer smoke` workflow. Use the `version` input to verify `main` or a release tag.

## Documentation Gates

- Confirm `README.md` describes Bot API 10.0 coverage accurately.
- Confirm `docs/API_COVERAGE.md` matches current implemented methods and intentional architecture differences.
- Confirm Bot API update work followed [`BOT_API_UPDATE_CHECKLIST.md`](BOT_API_UPDATE_CHECKLIST.md) when upstream coverage changed.
- Confirm `docs/ROADMAP.md` reflects the current release stage.
- Confirm `docs/releases/v0.6.0.md` is updated to final release notes before the `v0.6.0` tag is created.
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

- Current release tag: `v0.6.0`.
- Previous regular release tag: `v0.5.0`.
- Previous Bot API 10.0 package tag: `v0.4.0`.
- `v0.6.0` release commit: the commit targeted by the annotated `v0.6.0` tag.
- `v0.5.0` release commit: the commit targeted by the annotated `v0.5.0` tag.
- `v0.4.0` release commit: `2ce948b`.
- Annotated tag creation used explicit maintainer approval:

```bash
git tag -a v0.6.0 -m "Release v0.6.0"
```

- Only that tag was pushed:

```bash
git push https://github.com/xDilettante/ai-gram.git refs/tags/v0.6.0:refs/tags/v0.6.0
```

- Public module installation was verified after the tag became available:

```bash
go get github.com/xDilettante/ai-gram@v0.6.0
```

- GitHub Release: <https://github.com/xDilettante/ai-gram/releases/tag/v0.6.0>.

Do not run `git push --tags`, create unrelated tags, or create future GitHub Releases before explicit approval.

## Do Not Release If

- Any required local check fails.
- Examples do not compile through `scripts/check.sh`.
- A secret or token-bearing URL appears in tracked files, logs, or reports.
- README or coverage docs claim behavior that is not implemented.
- New public API was added without GoDoc and tests.
- Destructive/admin smoke was run against a non-test chat or without explicit confirmation.
- Access control examples default to public access.
