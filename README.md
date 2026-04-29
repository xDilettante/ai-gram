# ai-gram

`ai-gram` is a Go library project for working with the Telegram Bot API.

The project is in an early architecture stage. It provides a minimal package skeleton, practical incoming update types, a foundational HTTP core, the first public Bot API methods, media sending by file_id, URL, or multipart upload, minimal file download support, webhook management methods, a managed long polling runner, an inbound webhook HTTP handler, a small update dispatcher/router, and helper middleware. It does not yet implement FSM, scenes, storage, media groups, thumbnails, sendVideo/sendAudio/sendVoice, or full Bot API coverage.

## Статус

- Minimal Go module: present.
- Root facade package `aigram`: present.
- Base Telegram data types: include practical incoming message fields for text entities, captions, media, contacts, locations, venues, and callback queries.
- Bot client package: scaffolded with token validation, private token storage, and an internal HTTP call core.
- Typed Telegram API errors: scaffolded.
- Dispatcher/router: supports predicates, message/command/callback routes, middleware, fallback, and error handling.
- Middleware helpers: recover, timeout, and hook-based observability are available.
- Long polling transport: managed runner is available. Webhook transport: inbound HTTP handler is available.
- Telegram Bot API method coverage: `GetMe`, `SendMessage`, `SendPhoto`, `SendDocument`, the manual `GetUpdates` API call, `GetFile`, `DownloadFile`, multipart upload for `SendPhoto`/`SendDocument`, and JSON-only webhook management methods (`SetWebhook`, `DeleteWebhook`, `GetWebhookInfo`) are implemented. The rest of the Bot API is not implemented yet.
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
- `transport/webhook` provides an inbound `net/http` handler for Telegram webhook updates.
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


Send media by `file_id`, URL, or multipart upload:

```go
photoMessage, err := b.SendPhoto(ctx, aigram.SendPhotoParams{
    ChatID:  aigram.ChatIDInt(123456789),
    Photo:   aigram.FileID("existing-photo-file-id"),
    Caption: "Photo from file_id",
})
if err != nil {
    return err
}
fmt.Println(photoMessage.MessageID)

photoByURL, err := b.SendPhoto(ctx, aigram.SendPhotoParams{
    ChatID: aigram.ChatIDInt(123456789),
    Photo:  aigram.FileURL("https://example.com/photo.jpg"),
})
if err != nil {
    return err
}
fmt.Println(photoByURL.MessageID)

documentMessage, err := b.SendDocument(ctx, aigram.SendDocumentParams{
    ChatID:   aigram.ChatIDInt(123456789),
    Document: aigram.FileID("existing-document-file-id"),
    Caption:  "Document from file_id",
})
if err != nil {
    return err
}
fmt.Println(documentMessage.MessageID)
```

Upload a photo from `os.File`:

```go
photoFile, err := os.Open("photo.jpg")
if err != nil {
    return err
}
defer photoFile.Close()

uploadedPhoto, err := b.SendPhoto(ctx, aigram.SendPhotoParams{
    ChatID: aigram.ChatIDInt(123456789),
    Photo: aigram.FileUpload(aigram.UploadFile{
        Name:        "photo.jpg",
        Reader:      photoFile,
        ContentType: "image/jpeg",
    }),
    Caption: "Uploaded photo",
})
if err != nil {
    return err
}
fmt.Println(uploadedPhoto.MessageID)
```

Upload a document from `bytes.Reader`:

```go
report := []byte("report contents")

uploadedDocument, err := b.SendDocument(ctx, aigram.SendDocumentParams{
    ChatID: aigram.ChatIDInt(123456789),
    Document: aigram.FileUpload(aigram.UploadFile{
        Name:        "report.txt",
        Reader:      bytes.NewReader(report),
        ContentType: "text/plain",
    }),
})
if err != nil {
    return err
}
fmt.Println(uploadedDocument.MessageID)
```

`FileID` and `FileURL` are sent as JSON requests. `FileUpload` uses multipart/form-data and ai-gram generates the internal `attach://` value for the file field. The library consumes `UploadFile.Reader` but does not close it; the caller owns reader lifecycle. Thumbnails, media groups, `sendVideo`, `sendAudio`, and `sendVoice` are not implemented yet.

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


Handle common incoming message shapes:

```go
if err := d.OnMessageFunc(func(ctx context.Context, update telegram.Update) error {
    message := update.EffectiveMessage()
    if message == nil {
        return nil
    }

    switch {
    case message.IsCommand("start"):
        fmt.Println("command args:", message.CommandArguments())
    case message.HasPhoto():
        largest := message.Photo[len(message.Photo)-1]
        fmt.Println("photo:", largest.FileID)
    case message.HasDocument():
        fmt.Println("document:", message.Document.FileName)
    default:
        fmt.Println("text:", message.Text)
    }

    return nil
}); err != nil {
    return err
}

if err := d.OnCallbackDataFunc("confirm", func(ctx context.Context, update telegram.Update) error {
    fmt.Println("callback data:", update.CallbackQuery.Data)
    return nil
}); err != nil {
    return err
}
```

Telegram types currently support decoding incoming media and helper methods for handling them. Sending is currently limited to text, photo, and document methods; media groups and other media send methods will be added separately later.

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

The long polling runner fetches updates and calls a handler; `dispatch.Dispatcher` is one compatible handler implementation.


Download a file by `file_id` from an incoming document:

```go
if message.Document == nil {
    return nil
}

file, err := b.GetFile(ctx, aigram.GetFileParams{
    FileID: message.Document.FileID,
})
if err != nil {
    return err
}

var buf bytes.Buffer
if err := b.DownloadFile(ctx, file.FilePath, &buf); err != nil {
    return err
}
fmt.Println("downloaded bytes:", buf.Len())
```

For large files pass an `*os.File` or another streaming `io.Writer` instead of `bytes.Buffer`. Telegram download URLs contain the bot token; ai-gram builds them internally and does not expose them as a public API. The regular cloud Bot API has Telegram-side file download limits. Upload is currently implemented only for `SendPhoto` and `SendDocument`; download helpers never expose a full token-bearing download URL.

Serve inbound webhook updates with `net/http`:

```go
webhookHandler, err := webhook.New(d, webhook.Config{
    SecretToken: "your-secret-token",
})
if err != nil {
    return err
}

http.Handle("/telegram/webhook", webhookHandler)
if err := http.ListenAndServe(":8080", nil); err != nil {
    return err
}
```

Manage webhook registration through outbound Bot API methods:

```go
secret := "my_secret_123"

ok, err := b.SetWebhook(ctx, aigram.SetWebhookParams{
    URL:         "https://example.com/telegram/webhook",
    SecretToken: secret,
})
if err != nil {
    return err
}
if !ok {
    return fmt.Errorf("set webhook returned false")
}

webhookHandler, err := webhook.New(d, webhook.Config{
    SecretToken: secret,
})
if err != nil {
    return err
}
http.Handle("/telegram/webhook", webhookHandler)
```

The `SecretToken` passed to `SetWebhook` must match `transport/webhook.Config.SecretToken` so inbound requests can be verified.

Read current webhook status:

```go
info, err := b.GetWebhookInfo(ctx)
if err != nil {
    return err
}
fmt.Println(info.URL, info.PendingUpdateCount)
```

Delete a webhook:

```go
ok, err := b.DeleteWebhook(ctx, aigram.DeleteWebhookParams{
    DropPendingUpdates: true,
})
if err != nil {
    return err
}
fmt.Println(ok)
```

Webhook management is JSON-only for now. Webhook certificate upload, media groups, thumbnails, FSM, scenes, storage, dependency injection, and full Bot API coverage are not implemented yet.

## Development checks

```bash
gofmt -w .
go test ./...
go vet ./...
```
