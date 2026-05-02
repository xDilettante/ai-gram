# Contributing

Thank you for considering a contribution to `ai-gram`.

## Development setup

Use the Go version declared in [`go.mod`](go.mod). Before opening a pull request, run:

```bash
gofmt -w .
bash -n scripts/*.sh
go test ./...
go vet ./...
go list ./...
git diff --check
```

## Code style

- Keep public APIs small, typed, and idiomatic Go.
- Prefer the standard library unless a dependency has a clear benefit.
- Preserve `context.Context` on blocking operations.
- Do not log bot tokens, webhook secrets, payment payloads, passport data, managed bot tokens, or token-bearing URLs.
- Public repository text, comments, examples, and commit messages must be English.

## Bot API changes

When adding or changing Telegram Bot API support:

- Use the official Telegram Bot API documentation as the source of truth.
- Add typed params, result types, validation, and `httptest` coverage.
- Update README or docs when the public surface changes.
- Keep state-changing, payment-related, Passport, business, managed-token, admin, and webhook-certificate live checks manual-only unless maintainers explicitly approve a safe test plan.

## Pull requests

A good pull request includes:

- tests for behavior changes;
- documentation updates for user-visible changes;
- a note about any breaking change or compatibility risk;
- official Bot API documentation links for API surface changes;
- no secrets, private IDs, or token-bearing URLs in code, logs, fixtures, screenshots, or comments.
