# ai-gram

`ai-gram` is a Go library project for working with the Telegram Bot API.

The project is in an early architecture stage. It provides a minimal package skeleton and a few foundational contracts, but it does not yet implement Telegram Bot API methods, long polling, webhooks, or a production dispatcher.

## Статус

- Minimal Go module: present.
- Root facade package `aigram`: present.
- Base Telegram data types: started with a minimal subset.
- Bot client package: scaffolded with token validation and private token storage.
- Typed Telegram API errors: scaffolded.
- Dispatcher contracts and middleware composition: scaffolded.
- Long polling and webhook transports: placeholders only.
- Telegram Bot API method coverage: not implemented yet.
- Public API stability: not guaranteed before the first stable release.

## Планируемая архитектура

The library is split into small packages with clear responsibilities:

- `telegram` contains basic Telegram Bot API data contracts such as `Update`, `Message`, `User`, `Chat`, and `CallbackQuery`.
- `bot` contains the primary Bot API client type and configuration. It owns token handling and will later call Telegram methods.
- `internal/httpclient` contains low-level HTTP sending helpers. It is internal and must not leak into the public API.
- `errors` contains typed errors returned by Telegram Bot API responses.
- `dispatch` defines update handling, middleware, and dispatcher contracts without depending on HTTP details.
- `transport/longpoll` is reserved for long polling update delivery.
- `transport/webhook` is reserved for webhook update delivery.
- `aigram` is a lightweight root facade that re-exports the most important public types.

The intended dependency direction is data types first, then the Bot API client and transports, then dispatching and middleware. Transports deliver updates; dispatchers process already received updates; the API client does not know about dispatching.

## Development checks

```bash
gofmt -w .
go test ./...
go vet ./...
```
