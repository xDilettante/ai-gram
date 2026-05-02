# Scripts

The `scripts/` directory contains helper scripts used by maintainers for local checks, deploy tests, and live smoke runs. Normal library users do not need these scripts to import or use ai-gram.

## Publicly useful helpers

- `update_coverage_badge.sh` regenerates `docs/assets/coverage.svg` from `go test -coverprofile` without using external coverage services.
- `remote_logs.sh` and related log helpers are useful when you intentionally run the webhook example on your own host.
- `smoke_v02_send_methods.sh` and `smoke_media_group.sh` are targeted live smoke helpers, but they require a real bot token and a disposable test chat.

## Maintainer-only helpers

Deploy, discovery, notification, local Bot API service, and multi-bot smoke scripts are maintainer-oriented. They may read `.env.local` and ignored `.deploy/generated.env` values, open temporary SSH tunnels, or interact with real Telegram state.

See [`../docs/maintainer/DEPLOY_TESTING.md`](../docs/maintainer/DEPLOY_TESTING.md), [`../docs/maintainer/LIVE_SMOKE_MATRIX.md`](../docs/maintainer/LIVE_SMOKE_MATRIX.md), and [`../docs/maintainer/ENV_SMOKE_TEMPLATE.md`](../docs/maintainer/ENV_SMOKE_TEMPLATE.md) before running them.

## Safety rules

- Do not commit real bot tokens, webhook secrets, SSH details, invite links, payment payloads, or private chat IDs.
- Do not print token-bearing Bot API URLs.
- Keep destructive, payment-related, business, passport, managed-token, sticker-set, lifecycle, and webhook-certificate checks manual-only.
- Prefer `go test ./...` and `httptest` coverage before running live smoke.
