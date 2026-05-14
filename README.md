<p align="center">
  <img src="docs/assets/readme-header.png" alt="ai-gram — Telegram Bot API library for Go" width="100%">
</p>

<h1 align="center">ai-gram</h1>

<p align="center">
  <strong>A production-minded Telegram Bot API library for Go, built primarily with ChatGPT + Codex.</strong>
</p>

<p align="center">
  <a href="https://github.com/xDilettante/ai-gram/actions/workflows/ci.yml"><img alt="CI" src="https://img.shields.io/github/actions/workflow/status/xDilettante/ai-gram/ci.yml?branch=main&style=for-the-badge&label=CI&logo=github"></a>
  <img alt="Telegram Bot API 10.0 complete" src="https://img.shields.io/badge/Telegram%20Bot%20API-10.0%20complete-26A5E4?style=for-the-badge&logo=telegram&logoColor=white">
  <img alt="Go" src="https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white">
  <a href="LICENSE"><img alt="License: MIT" src="https://img.shields.io/badge/license-MIT-F59E0B?style=for-the-badge"></a>
</p>

<p align="center">
  <img alt="Coverage: 63.4%" src="https://img.shields.io/badge/coverage-63.4%25-F59E0B?style=for-the-badge">
  <img alt="API coverage: 10.0 complete" src="https://img.shields.io/badge/API%20coverage-10.0%20complete-7C3AED?style=for-the-badge">
  <img alt="Release: v0.5.0" src="https://img.shields.io/badge/release-v0.5.0-16A34A?style=for-the-badge">
  <img alt="Built with ChatGPT and Codex" src="https://img.shields.io/badge/built%20with-ChatGPT%20%2B%20Codex-8B5CF6?style=for-the-badge">
</p>

<p align="center">
  <a href="#quick-start">Quick start</a> ·
  <a href="docs/ARCHITECTURE.md">Architecture</a> ·
  <a href="docs/API_COVERAGE.md">API coverage</a> ·
  <a href="docs/MANUAL_TESTING.md">Manual testing</a> ·
  <a href="CONTRIBUTING.md">Contributing</a> ·
  <a href="SECURITY.md">Security</a> ·
  <a href="docs/ROADMAP.md">Roadmap</a>
</p>

`ai-gram` is a typed Go library for building Telegram Bot API clients, update transports, dispatchers, and middleware. It is designed as an AI-native open-source project: the implementation is built almost entirely with ChatGPT and Codex, while architecture, review, scope, and release decisions stay under human maintainer control.

The library focuses on a clear public API, token-safe HTTP behavior, replaceable transports, and testable building blocks instead of framework magic. It is suitable for low-level Bot API calls as well as production bot foundations that need long polling, webhooks, routing, middleware, and typed Telegram data contracts.

> **Compatibility:** Telegram Bot API 10.0 code coverage is complete with documented architecture differences; see [`docs/BOT_API_10_0_FINAL_AUDIT.md`](docs/BOT_API_10_0_FINAL_AUDIT.md). `ai-gram` is still a pre-v1 Go module, so public APIs may evolve before v1.0.

## Highlights

- Typed Bot API method parameters, result types, and Telegram update/message contracts.
- JSON and multipart method calls with `FileRef`/`FileUpload` helpers for upload-capable methods.
- Compact typed callback data helpers for inline keyboard flows.
- Long polling transport, inbound webhook handler, dispatcher/router, predicates, middleware, fallback, and error handling.
- Token-safe configuration with configurable Bot API base URLs for official or local Bot API servers.
- Broad test coverage built around unit tests and `httptest`-friendly client configuration.
- Practical public examples, with advanced maintainer tooling kept separate from the user-facing quick start.

## Quick start

Install the module once the repository or tag you need is available to your Go toolchain:

```bash
go get github.com/xDilettante/ai-gram
```

Create a bot client and call a typed Bot API method:

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

    bot, err := aigram.New(aigram.Config{Token: token})
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

To send messages, use `SendMessage` with typed parameters such as `aigram.SendMessageParams` and `aigram.ChatIDInt` or `aigram.ChatIDString`. For the full Telegram Bot API surface, import `github.com/xDilettante/ai-gram/bot` for method params and `github.com/xDilettante/ai-gram/telegram` for Telegram data contracts.

## Why ai-gram

`ai-gram` keeps the library layers separate:

- `telegram` contains Telegram Bot API data contracts.
- `bot` contains the primary Bot API client and typed method parameters.
- `transport/longpoll` provides a managed long polling update source.
- `transport/webhook` provides an inbound webhook HTTP handler.
- `dispatch` routes updates to handlers.
- `middleware` provides reusable dispatcher middleware.
- `callback` builds and parses compact typed inline keyboard callback data.
- `errors` exposes typed Telegram API errors.
- The root package `aigram` provides a small convenience facade for quick-start code; advanced methods and contracts are intentionally used through `bot` and `telegram`.

The project intentionally keeps AI-assisted development visible without turning it into marketing noise: ChatGPT and Codex produce most of the code and documentation, while the maintainer directs requirements, validates behavior, and decides what is safe to ship.

## Beginner-friendly examples

Start with the numbered examples under [`examples/`](examples/):

- [`examples/01_echo_bot`](examples/01_echo_bot) — minimal long polling echo bot with `/start`.
- [`examples/02_commands`](examples/02_commands) — command routing with `/start`, `/help`, and `/about`.
- [`examples/03_reply_keyboard`](examples/03_reply_keyboard) — regular reply keyboard and keyboard removal.
- [`examples/04_inline_keyboard`](examples/04_inline_keyboard) — inline keyboard callbacks, `AnswerCallbackQuery`, and message editing.
- [`examples/05_media_upload`](examples/05_media_upload) — photo/document sends by `file_id` or local upload.
- [`examples/06_webhook_basic`](examples/06_webhook_basic) — simple webhook server with `/healthz` and `/webhook`.

Additional examples remain available for advanced scenarios:

- [`examples/echo_longpoll`](examples/echo_longpoll) — long polling bot with commands, echo, inline keyboard, callbacks, and message editing.
- [`examples/webhook_basic`](examples/webhook_basic) — minimal inbound webhook server skeleton.
- [`examples/inline_longpoll`](examples/inline_longpoll) — inline keyboard callbacks and access-control demo.
- [`examples/webhook_server`](examples/webhook_server) — extended webhook server example.
- [`examples/media_upload`](examples/media_upload) — document upload and file download checks.
- [`examples/local_api_server`](examples/local_api_server) — connectivity check for a local Telegram Bot API server.

Maintainer-only smoke examples are separated under [`examples/maintainer/`](examples/maintainer/) and are not needed for normal library use.

## Documentation

- [`docs/ARCHITECTURE.md`](docs/ARCHITECTURE.md) — short package and update-flow architecture diagram.
- [`docs/API_COVERAGE.md`](docs/API_COVERAGE.md) — Bot API method/type coverage inventory and architecture notes.
- [`docs/releases/v0.5.0.md`](docs/releases/v0.5.0.md) — release notes for the current public release.
- [`docs/PRE_V1_NOTES.md`](docs/PRE_V1_NOTES.md) — current pre-v1 API shape and breaking-change notes.
- [`docs/BOT_API_10_0_FINAL_AUDIT.md`](docs/BOT_API_10_0_FINAL_AUDIT.md) — final Bot API 10.0 coverage audit.
- [`docs/releases/v0.4.0.md`](docs/releases/v0.4.0.md) — release notes for the Bot API 10.0 package milestone.
- [`docs/maintainer/BOT_API_10_0_RELEASE_READINESS.md`](docs/maintainer/BOT_API_10_0_RELEASE_READINESS.md) — maintainer release-readiness notes for Bot API 10.0.
- [`docs/BOT_API_10_0_COVERAGE_PLAN.md`](docs/BOT_API_10_0_COVERAGE_PLAN.md) — Bot API 10.0 update plan.
- [`docs/MANUAL_TESTING.md`](docs/MANUAL_TESTING.md) — public manual testing guide.
- [`docs/ROADMAP.md`](docs/ROADMAP.md) — stabilization and future work.
- [`CONTRIBUTING.md`](CONTRIBUTING.md) — contribution guide and PR expectations.
- [`SECURITY.md`](SECURITY.md) — security reporting policy.
- [`CHANGELOG.md`](CHANGELOG.md) — project changelog.
- [`LICENSE`](LICENSE) — MIT license.

Historical Bot API 9.6 workstream notes remain available in [`docs/BOT_API_9_6_FINAL_AUDIT.md`](docs/BOT_API_9_6_FINAL_AUDIT.md), [`docs/BOT_API_9_6_AUDIT.md`](docs/BOT_API_9_6_AUDIT.md), and [`docs/BOT_API_9_6_COVERAGE_PLAN.md`](docs/BOT_API_9_6_COVERAGE_PLAN.md). Treat them as historical records when they mention older compatibility decisions that were changed during later pre-v1 cleanup.

Maintainer-only deploy, live-smoke, and release-readiness notes live under [`docs/maintainer/`](docs/maintainer/). They are useful for project maintainers, but intentionally separated from the public quick start.

## Development checks

Use the unified local check before submitting changes:

```bash
scripts/check.sh
```

The script mirrors the main CI checks. Individual commands are:

```bash
find . -path './.git' -prune -o -path './build' -prune -o -path './vendor' -prune -o -name '*.go' -type f -print0 | xargs -0 gofmt -l
bash -n scripts/*.sh
go test ./...
go vet ./...
go build ./...
go list ./...
git diff --check
```

For a quick repository overview before making changes:

```bash
scripts/ai-context.sh
```


## Safety notes

- Never commit real bot tokens, webhook secrets, private chat IDs, payment payloads, passport data, managed bot tokens, or token-bearing URLs.
- `SetWebhook` supports the official upload-only certificate path via `FileUpload`; certificate live checks should use disposable test certificates only.
- `GetChat` returns the official `ChatFullInfo` result shape; `GetChatFullInfo` remains as a same-result alias during pre-v1 cleanup.
- `ChatMember` is decoded as official `ChatMember*` variants through the `telegram.ChatMember` interface.
- `CallbackQuery.Message` uses the official `MaybeInaccessibleMessage` shape with helpers for accessible messages.
- Pre-v1 API cleanup notes live in [`docs/PRE_V1_NOTES.md`](docs/PRE_V1_NOTES.md).
