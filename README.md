# ai-gram

`ai-gram` is a Go library project for working with the Telegram Bot API.

The project is in an early architecture stage. It provides a minimal package skeleton, a foundational HTTP core, and the first public Bot API methods, but it does not yet implement long polling, webhooks, or a production dispatcher.

## Статус

- Minimal Go module: present.
- Root facade package `aigram`: present.
- Base Telegram data types: started with a minimal subset.
- Bot client package: scaffolded with token validation, private token storage, and an internal HTTP call core.
- Typed Telegram API errors: scaffolded.
- Dispatcher contracts and middleware composition: scaffolded.
- Long polling and webhook transports: placeholders only.
- Telegram Bot API method coverage: only `GetMe` and `SendMessage` are implemented. `GetUpdates` and the rest of the Bot API are not implemented yet.
- Public API stability: not guaranteed before the first stable release.

## Планируемая архитектура

The library is split into small packages with clear responsibilities:

- `telegram` contains basic Telegram Bot API data contracts such as `Update`, `Message`, `User`, `Chat`, and `CallbackQuery`.
- `bot` contains the primary Bot API client type and configuration. It owns token handling and an unexported JSON call core that will later power public Telegram methods.
- `internal/httpclient` contains low-level HTTP sending helpers, response body handling, and HTTP status checks. It is internal and must not leak into the public API.
- `errors` contains typed errors returned by Telegram Bot API responses.
- `dispatch` defines update handling, middleware, and dispatcher contracts without depending on HTTP details.
- `transport/longpoll` is reserved for long polling update delivery.
- `transport/webhook` is reserved for webhook update delivery.
- `aigram` is a lightweight root facade that re-exports the most important public types.

The intended dependency direction is data types first, then the Bot API client and transports, then dispatching and middleware. Transports deliver updates; dispatchers process already received updates; the API client does not know about dispatching.

## Usage examples

Create a bot and call `getMe`:

```go
ctx := context.Background()

b, err := aigram.New(aigram.BotConfig{Token: token})
if err != nil {
    return err
}

me, err := b.GetMe(ctx)
if err != nil {
    return err
}
fmt.Println(me.Username)
```

Send a text message:

```go
message, err := b.SendMessage(ctx, aigram.SendMessageParams{
    ChatID: aigram.ChatIDInt(123456789),
    Text:   "Hello from ai-gram",
})
if err != nil {
    return err
}
fmt.Println(message.MessageID)
```

Only these first methods are available for now; the library does not yet claim full Telegram Bot API coverage.

## Development checks

```bash
gofmt -w .
go test ./...
go vet ./...
```
