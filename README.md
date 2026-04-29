# ai-gram

`ai-gram` is a Go library project for working with the Telegram Bot API.

The project is in an early architecture stage. It provides a minimal package skeleton, a foundational HTTP core, the first public Bot API methods, a managed long polling runner, a small update dispatcher/router, and helper middleware. It does not yet implement webhooks, FSM, scenes, or storage.

## Статус

- Minimal Go module: present.
- Root facade package `aigram`: present.
- Base Telegram data types: started with a minimal subset.
- Bot client package: scaffolded with token validation, private token storage, and an internal HTTP call core.
- Typed Telegram API errors: scaffolded.
- Dispatcher/router: supports predicates, message/command/callback routes, middleware, fallback, and error handling.
- Middleware helpers: recover, timeout, and hook-based observability are available.
- Long polling transport: managed runner is available. Webhook transport: placeholder only.
- Telegram Bot API method coverage: `GetMe`, `SendMessage`, and the manual `GetUpdates` API call are implemented. The rest of the Bot API is not implemented yet.
- Public API stability: not guaranteed before the first stable release.

## Планируемая архитектура

The library is split into small packages with clear responsibilities:

- `telegram` contains basic Telegram Bot API data contracts such as `Update`, `Message`, `User`, `Chat`, and `CallbackQuery`.
- `bot` contains the primary Bot API client type and configuration. It owns token handling and an unexported JSON call core that will later power public Telegram methods.
- `internal/httpclient` contains low-level HTTP sending helpers, response body handling, and HTTP status checks. It is internal and must not leak into the public API.
- `errors` contains typed errors returned by Telegram Bot API responses.
- `dispatch` defines update routing, middleware, fallback handling, and error handling without depending on HTTP details.
- `middleware` provides reusable dispatch middleware helpers for panic recovery, per-update timeout contexts, and hook-based observability.
- `transport/longpoll` provides a managed runner that repeatedly calls `GetUpdates` and passes updates to a handler.
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

Fetch updates manually with one `getUpdates` API call:

```go
updates, err := b.GetUpdates(ctx, aigram.GetUpdatesParams{
    Limit:   10,
    Timeout: 0,
})
if err != nil {
    return err
}
for _, update := range updates {
    if update.Message != nil {
        fmt.Println(update.Message.Text)
    }
}
```

Create a small dispatcher:

```go
d := dispatch.New(dispatch.WithErrorHandler(func(ctx context.Context, update telegram.Update, err error) {
    fmt.Println("handler error:", err)
}))

if err := d.OnCommandFunc("start", func(ctx context.Context, update telegram.Update) error {
    fmt.Println("start command")
    return nil
}); err != nil {
    return err
}

if err := d.OnMessageFunc(func(ctx context.Context, update telegram.Update) error {
    fmt.Println(update.Message.Text)
    return nil
}); err != nil {
    return err
}
```

Add helper middleware:

```go
type Observer struct{}

func (Observer) OnUpdateStart(ctx context.Context, update telegram.Update) {}
func (Observer) OnUpdateFinish(ctx context.Context, update telegram.Update, err error, duration time.Duration) {
    fmt.Println("handled update in", duration, "error:", err)
}

d.Use(
    middleware.Recover(nil),
    middleware.Timeout(5*time.Second),
    middleware.Observe(Observer{}),
)
```

Observability is hook-based only for now; the library does not include Prometheus, OpenTelemetry, or a logger.

Run managed long polling with the dispatcher:

```go
runner, err := longpoll.New(b, d, longpoll.Config{
    Timeout: 30,
})
if err != nil {
    return err
}

if err := runner.Run(ctx); err != nil {
    return err
}
```

The long polling runner fetches updates and calls a handler; `dispatch.Dispatcher` is one compatible handler implementation. FSM, scenes, storage, dependency injection, full Bot API coverage, and webhook support are not implemented yet.

## Development checks

```bash
gofmt -w .
go test ./...
go vet ./...
```
