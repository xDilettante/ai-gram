# Support

This repository is maintained as a Go library for the Telegram Bot API. Support works best when reports are reproducible, scoped to a package or Bot API area, and safe to share publicly.

## Where to ask

- Use a **bug report** when a released version or `main` behaves incorrectly.
- Use a **question** when you need help choosing a package, transport, dispatcher pattern, or Bot API method.
- Use a **feature request** for new helpers, examples, or Telegram Bot API additions.
- Use the **security policy** for vulnerabilities or anything that requires private disclosure.

GitHub issue forms are available under the repository's **New issue** page. Prefer those forms because they request the version, Go version, Bot API area, expected behavior, actual behavior, and a minimal reproduction.

## What to include

For bugs, include:

- `ai-gram` version, tag, pseudo-version, or commit hash;
- Go version and operating system;
- affected package or Bot API method;
- a minimal safe reproduction or focused failing test;
- redacted logs or error output;
- official Telegram Bot API documentation links when the issue is about API coverage or behavior.

For feature requests, include:

- the real use case;
- the desired Go API shape when you have one;
- official Telegram Bot API links for Bot API surface changes;
- compatibility or migration concerns.

## Safety rules

Never include real bot tokens, webhook secrets, private chat IDs, payment payloads, Passport data, managed bot tokens, private keys, cookies, authorization headers, or token-bearing URLs in public issues, pull requests, screenshots, logs, examples, or test fixtures.

If a reproduction needs Telegram state, use a dedicated test bot and test chat, then redact identifiers before posting. Destructive/admin, payment, Passport, business, managed-token, sticker-set mutation, webhook-certificate, and lifecycle checks should stay manual unless a maintainer explicitly approves a targeted test plan.

## Scope

This project can help with:

- typed Bot API client behavior;
- Telegram data contract encoding and decoding;
- dispatcher, middleware, long polling, and webhook library patterns;
- examples and documentation in this repository;
- Bot API compatibility audits.

This project is not a hosted bot service. It cannot debug private production bots, recover bot tokens, inspect private Telegram chats, or provide official Telegram platform support.
