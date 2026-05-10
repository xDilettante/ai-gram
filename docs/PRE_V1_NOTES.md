# Pre-v1 Notes

`ai-gram` is a pre-v1 Go module. Public APIs can still change before v1.0 when a cleanup makes the library clearer, more idiomatic, or closer to the official Telegram Bot API shape.

## Current Public Shape

- Use the root `aigram` package for quick-start code and common helpers.
- Use `github.com/xDilettante/ai-gram/bot` for the full outgoing Bot API method surface.
- Use `github.com/xDilettante/ai-gram/telegram` for Telegram data contracts.
- Bot construction uses `aigram.Config` or `bot.Config`.
- Inline login buttons use `telegram.LoginURL`.
- `GetChat` returns the official `telegram.ChatFullInfo` result shape.
- `GetChatFullInfo` remains as a same-result pre-v1 alias for code that already uses the explicit name.
- `telegram.ChatMember` is an interface implemented by the official `ChatMember*` variants.
- `CallbackQuery.Message` uses the official `MaybeInaccessibleMessage` shape with helpers for accessible messages.

## Breaking Changes

Breaking pre-v1 changes are recorded in [`../CHANGELOG.md`](../CHANGELOG.md) under `Unreleased` or the release where they shipped.

Recent unreleased breaking cleanups:

- `bot.BotConfig` and root `aigram.BotConfig` were renamed to `Config`.
- `telegram.LoginUrl` was renamed to `LoginURL`.

## Historical Documents

Older Bot API 9.6 audit and coverage-plan documents record the implementation state at the time they were written. They are useful for archaeology, but the current API shape is defined by the source code, README, architecture notes, API coverage inventory, this document, and the changelog.

When historical documents mention older compatibility choices such as a flat `ChatMember` shape or `GetChat` returning `telegram.Chat`, treat those notes as superseded by the current Bot API 10.0 documentation and source code.
