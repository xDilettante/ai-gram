# ai-gram

`ai-gram` is a typed Go library for building Telegram Bot API clients, update transports, dispatchers, and middleware.

## Status

- Latest public tag: `v0.2.0`.
- Current local `main`: Telegram Bot API 9.6 code coverage is complete with documented architecture differences.
- Publication is intentionally separate from local readiness: this repository has not been pushed, tagged, or released after the local-only Bot API 9.6 workstream in this checkout.
- Sensitive or state-changing live checks remain manual-only: payments, Stars, gifts, business APIs, managed bot tokens, passport data, admin/destructive chat methods, sticker set mutation, games, lifecycle `LogOut`/`Close`, and webhook certificate upload.

See [`docs/API_COVERAGE.md`](docs/API_COVERAGE.md) and [`docs/BOT_API_9_6_FINAL_AUDIT.md`](docs/BOT_API_9_6_FINAL_AUDIT.md) for the detailed coverage inventory.

## Install

```bash
go get github.com/xDilettante/ai-gram@v0.2.0
```

For local development from this checkout, use the module path `github.com/xDilettante/ai-gram` and a local `replace` directive if needed.

## Quick start

```go
package main

import (
    "context"
    "fmt"
    "log"
    "os"

    aigram "github.com/xDilettante/ai-gram"
)

func main() {
    token := os.Getenv("AIGRAM_BOT_TOKEN")
    if token == "" {
        log.Fatal("AIGRAM_BOT_TOKEN is required")
    }

    bot, err := aigram.New(aigram.BotConfig{Token: token})
    if err != nil {
        log.Fatal(err)
    }

    me, err := bot.GetMe(context.Background())
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println(me.Username)
}
```

To send messages, use `SendMessage` with typed parameters such as `aigram.SendMessageParams` and `aigram.ChatIDInt` or `aigram.ChatIDString`.

## What is included

- Typed Bot API client with token-safe configuration and configurable Bot API base URLs.
- JSON and multipart method calls, including `FileRef`/`FileUpload` helpers for upload-capable methods.
- Incoming update and message types for Bot API 9.6, including business, gifts, Stars, paid media, polls, checklists, Web Apps, games, passport, reply metadata, service messages, and chat boosts.
- Long polling transport and inbound webhook handler.
- Dispatcher/router with predicates, command routes, middleware, fallback, and error handling.
- Middleware helpers for panic recovery, timeouts, observability hooks, and access control.
- Testkit-style examples and `httptest`-friendly client configuration.

## Package overview

- `telegram` — Telegram Bot API data contracts.
- `bot` — primary Bot API client and typed method parameters.
- `transport/longpoll` — managed long polling update source.
- `transport/webhook` — inbound webhook HTTP handler.
- `dispatch` — update routing and handler execution.
- `middleware` — reusable dispatcher middleware.
- `errors` — typed Telegram API errors.
- root package `aigram` — convenience facade and common re-exports.

## Examples

Runnable examples are under [`examples/`](examples/):

- `examples/echo_longpoll` — basic long polling echo bot.
- `examples/inline_longpoll` — inline keyboard callbacks and access-control demo.
- `examples/webhook_server` — inbound webhook server example.
- `examples/media_upload` — document upload and file download checks.
- `examples/local_api_server` — connectivity check for a local Telegram Bot API server.

Maintainer-only smoke examples are separated under `examples/maintainer/` and are not needed for normal library use.

## Documentation

- [`docs/API_COVERAGE.md`](docs/API_COVERAGE.md) — Bot API method/type coverage inventory and architecture notes.
- [`docs/BOT_API_9_6_COVERAGE_PLAN.md`](docs/BOT_API_9_6_COVERAGE_PLAN.md) — local-only Bot API 9.6 coverage plan and freeze policy.
- [`docs/BOT_API_9_6_FINAL_AUDIT.md`](docs/BOT_API_9_6_FINAL_AUDIT.md) — final Bot API 9.6 coverage audit.
- [`docs/MANUAL_TESTING.md`](docs/MANUAL_TESTING.md) — public manual testing guide.
- [`docs/ROADMAP.md`](docs/ROADMAP.md) — stabilization and future work.
- [`CHANGELOG.md`](CHANGELOG.md) — project changelog.

Maintainer-only deploy, live-smoke, and release-readiness notes live under [`docs/maintainer/`](docs/maintainer/). They are useful for project maintainers, but they are intentionally separated from the public quick start.

## Development checks

```bash
gofmt -w .
go test ./...
go vet ./...
```

For script syntax checks:

```bash
bash -n scripts/*.sh
```

## Safety notes

- Never commit real bot tokens, webhook secrets, private chat IDs, payment payloads, passport data, managed bot tokens, or token-bearing URLs.
- `SetWebhook` supports the official upload-only certificate path via `FileUpload`; certificate live checks should use disposable test certificates only.
- `GetChat` remains a backward-compatible minimal chat decode. Use `GetChatFullInfo` for the full Bot API 9.6 `getChat` result shape.
- `ChatMember` keeps a flat compatibility shape while decoding official Bot API 9.6 fields.
- `CallbackQuery.Message` remains available for accessible messages; maybe-inaccessible callback message data is represented separately for compatibility.
