# Root Facade Cleanup Plan

Date: 2026-05-10

## Goal

Reduce the root `aigram` package from a broad re-export mirror into a small convenience facade for quick-start usage.

## Scope

- Keep root-level convenience for creating a bot and the most common examples.
- Move advanced Bot API method params and Telegram contracts to explicit `bot` and `telegram` imports.
- Update examples, README, architecture notes, API coverage, roadmap, and changelog.
- Keep behavior and JSON encoding unchanged.

## Out Of Scope

- No new Bot API features.
- No live Telegram smoke.
- No `v0.4.0` tag or GitHub Release.
- No package renames or code generation.

## Planned Root Surface

- `Bot`, `BotConfig`, `New`, `NewBot`
- `ChatID`, `ChatIDInt`, `ChatIDString`
- `SendMessageParams`
- common update contracts: `Update`, `Message`, `CallbackQuery`
- common reply helpers: inline keyboard, reply keyboard, remove keyboard, `ReplyParameters`
- common file helpers only if needed by public beginner examples

Everything else should be imported from:

- `github.com/xDilettante/ai-gram/bot`
- `github.com/xDilettante/ai-gram/telegram`
- `github.com/xDilettante/ai-gram/dispatch`
- `github.com/xDilettante/ai-gram/middleware`

## Checks

- `scripts/check.sh`
- external consumer smoke from a temporary Go module

## Done Criteria

- `aigram.go` is small enough to read as a facade, not a generated-looking alias list.
- Public examples compile with explicit package imports where advanced APIs are used.
- README explains the package split clearly.
- Changes are committed atomically.
