# Contributing

Thank you for considering a contribution to `ai-gram`.

## Development setup

Use the Go version declared in [`go.mod`](go.mod). Before opening a pull request, run:

```bash
scripts/check.sh
```

The unified check mirrors CI. When you need focused commands, use:

```bash
gofmt -w .
bash -n scripts/*.sh
go test ./...
go vet ./...
go build ./...
go list ./...
git diff --check
```

For a quick local repository overview before making changes:

```bash
scripts/ai-context.sh
```

## Project shape

Keep package responsibilities separate:

- `bot` contains outgoing Bot API calls and typed request params.
- `telegram` contains JSON-compatible Telegram data contracts and small data helpers.
- `dispatch` contains update routing helpers.
- `middleware` contains reusable dispatcher middleware.
- `transport/longpoll` and `transport/webhook` contain update intake.
- `callback`, `errors`, and other helper packages should stay focused and policy-light.
- `examples` should show usage patterns without hiding transport behavior or requiring unsafe live state by default.

## Code style

- Keep public APIs small, typed, and idiomatic Go.
- Prefer the standard library unless a dependency has a clear benefit.
- Preserve `context.Context` on blocking operations.
- Avoid framework-level globals, hidden background behavior, and broad refactors in focused PRs.
- Do not log bot tokens, webhook secrets, payment payloads, passport data, managed bot tokens, or token-bearing URLs.
- Public repository text, comments, examples, and commit messages must be English.

## Bot API changes

When adding or changing Telegram Bot API support:

- Use the official Telegram Bot API documentation as the source of truth.
- Follow the maintainer workflow in [`docs/maintainer/BOT_API_UPDATE_CHECKLIST.md`](docs/maintainer/BOT_API_UPDATE_CHECKLIST.md).
- Add typed params, result types, validation, and `httptest` coverage.
- Update [`docs/API_COVERAGE.md`](docs/API_COVERAGE.md) when method, object, or update coverage changes.
- Update README or docs when the public surface changes.
- Keep state-changing, payment-related, Passport, business, managed-token, admin, and webhook-certificate live checks manual-only unless maintainers explicitly approve a safe test plan.

## Tests and documentation

Behavior changes should include focused tests. Prefer synthetic fixtures and `httptest` over live Telegram checks.

Documentation changes should be factual and current. Do not claim support that is not implemented and tested. Update [`CHANGELOG.md`](CHANGELOG.md) for user-visible changes.

## Pull requests

A good pull request includes:

- tests for behavior changes;
- documentation updates for user-visible changes;
- `scripts/check.sh` results or the focused checks that were run;
- a note about any breaking change or compatibility risk;
- official Bot API documentation links for API surface changes;
- no secrets, private IDs, or token-bearing URLs in code, logs, fixtures, screenshots, or comments.
