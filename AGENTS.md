# AGENTS.md

Project instructions for AI agents working on `ai-gram`.

These instructions apply to this repository. A specific user request may narrow the task, but it does not remove the safety and quality rules below.

## Project Shape

`ai-gram` is a Go library for the Telegram Bot API. It is library-first, not an application framework.

The repository provides:

- a typed Bot API client in `bot/`;
- Telegram data contracts in `telegram/`;
- update routing in `dispatch/`;
- reusable middleware in `middleware/`;
- long polling and webhook transports in `transport/`;
- runnable examples in `examples/`;
- maintainer smoke/deploy helpers in `scripts/`.

Keep these layers separate. The API client must not depend on dispatcher, middleware, long polling, or webhook packages. Telegram data types should stay predictable JSON-compatible contracts.

## Language Policy

Repository content must be English:

- README, docs, changelog, GoDoc, code comments;
- examples and user-visible example output;
- scripts, script output, and shell error messages;
- commit messages, release notes, issue/PR templates.

Russian is allowed only in direct private conversation with the user.

## Development Rules

- Read this file and inspect local conventions before changing code.
- Check `git status` before edits and do not overwrite unrelated user changes.
- Make the smallest focused change that satisfies the task.
- Prefer the Go standard library. Add external dependencies only with a clear reason.
- Keep public APIs idiomatic, documented, and stable unless the task explicitly asks for a breaking change.
- All blocking operations must accept or use a caller-provided `context.Context`.
- Do not add framework magic, global mutable state, hidden background behavior, or broad refactors without need.

## Security Rules

- Never commit real bot tokens, webhook secrets, private chat IDs, payment payloads, passport data, SSH details, or token-bearing URLs.
- Never log full Telegram Bot API URLs containing `/bot<TOKEN>/`.
- Preserve token redaction behavior in errors and diagnostics.
- Treat live smoke scripts as potentially stateful: they may interact with real Telegram bots, chats, webhooks, local Bot API servers, SSH, or systemd.
- Do not run destructive, payment, passport, managed-token, sticker-set, webhook-certificate, deploy, or live Telegram checks unless the user explicitly approves them.

## Testing And Checks

Prefer the unified local check:

```bash
scripts/check.sh
```

Useful focused commands:

```bash
gofmt -w .
go test ./...
go vet ./...
go build ./...
bash -n scripts/*.sh
git diff --check
```

For repository context before a task:

```bash
scripts/ai-context.sh
```

Every meaningful behavior change should include tests. Important test areas include request encoding, Telegram API errors, network errors, context cancellation, multipart uploads, long polling offsets, webhook validation, dispatcher routing, middleware order, and JSON compatibility.

## Documentation

Update documentation when user-facing behavior changes:

- `README.md` for public workflow or API changes;
- `docs/` for architecture, coverage, or maintainer workflow changes;
- examples when the intended usage changes;
- GoDoc for exported symbols.

Do not let README or examples promise behavior that is not implemented.

## Before Finishing

Run relevant checks and report exactly what ran. If a check is skipped, state why.

Default final report shape:

```text
Done:
- ...

Changed files:
- ...

Checks:
- ...

Notes:
- ...
```

Do not push, tag, create releases, or modify remotes unless the user explicitly asks for that exact action.
