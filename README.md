# ai-gram

`ai-gram` is a Go library project for working with the Telegram Bot API.

The project is in an early architecture stage. It provides a minimal package skeleton, practical incoming update types, a foundational HTTP core, the first public Bot API methods, media sending by file_id, URL, or multipart upload, minimal file download support, webhook management methods, a managed long polling runner, an inbound webhook HTTP handler, a small update dispatcher/router, and helper middleware. It does not yet implement FSM, scenes, storage, media groups, thumbnails, sendAnimation/sendVideoNote, or full Bot API coverage.

## Статус

- Minimal Go module: present.
- Root facade package `aigram`: present.
- Base Telegram data types: include practical incoming message fields for text entities, captions, media, contacts, locations, venues, and callback queries.
- Bot client package: scaffolded with token validation, private token storage, and an internal HTTP call core.
- Typed Telegram API errors: scaffolded.
- Dispatcher/router: supports predicates, message/command/callback routes, middleware, fallback, and error handling.
- Middleware helpers: recover, timeout, and hook-based observability are available.
- Long polling transport: managed runner is available. Webhook transport: inbound HTTP handler is available.
- Telegram Bot API method coverage: `GetMe`, `SendMessage`, `SendPhoto`, `SendDocument`, `SendVideo`, `SendAudio`, `SendVoice`, `AnswerCallbackQuery`, `EditMessageText`, `EditMessageCaption`, `EditMessageReplyMarkup`, `DeleteMessage`, reply markup for supported send and edit methods, the manual `GetUpdates` API call, `GetFile`, `DownloadFile`, multipart upload for media send methods, and JSON-only webhook management methods (`SetWebhook`, `DeleteWebhook`, `GetWebhookInfo`) are implemented. The rest of the Bot API is not implemented yet.
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

Reply to an incoming message and keep a forum topic/thread when Telegram provides one:

```go
msg := update.EffectiveMessage()
if msg == nil {
    return nil
}

_, err := b.SendMessage(ctx, aigram.SendMessageParams{
    ChatID:          aigram.ChatIDInt(msg.Chat.ID),
    MessageThreadID: msg.MessageThreadID,
    Text:            "Reply from ai-gram",
    ReplyParameters: &aigram.ReplyParameters{MessageID: msg.MessageID},
})
if err != nil {
    return err
}
```

Attach reply markup:

```go
inlineKeyboard := aigram.NewInlineKeyboard(
    []aigram.InlineKeyboardButton{
        aigram.InlineButtonCallback("Confirm", "confirm"),
        aigram.InlineButtonURL("Docs", "https://core.telegram.org/bots/api"),
    },
)

_, err = b.SendMessage(ctx, aigram.SendMessageParams{
    ChatID:      aigram.ChatIDInt(123456789),
    Text:        "Choose an action",
    ReplyMarkup: inlineKeyboard,
})
if err != nil {
    return err
}

d := dispatch.New()
if err := d.OnCallbackDataFunc("confirm", func(ctx context.Context, update telegram.Update) error {
    if update.CallbackQuery == nil {
        return nil
    }

    ok, err := b.AnswerCallbackQuery(ctx, aigram.AnswerCallbackQueryParams{
        CallbackQueryID: update.CallbackQuery.ID,
        Text:            "Готово",
    })
    if err != nil {
        return err
    }
    fmt.Println("callback answered:", ok)
    return nil
}); err != nil {
    return err
}

if err := d.OnCallbackDataFunc("danger", func(ctx context.Context, update telegram.Update) error {
    if update.CallbackQuery == nil {
        return nil
    }

    _, err := b.AnswerCallbackQuery(ctx, aigram.AnswerCallbackQueryParams{
        CallbackQueryID: update.CallbackQuery.ID,
        Text:            "Нужно подтверждение",
        ShowAlert:       true,
    })
    return err
}); err != nil {
    return err
}
```

Use a regular reply keyboard and remove it later:

```go
replyKeyboard := aigram.NewReplyKeyboard(
    []aigram.KeyboardButton{
        aigram.KeyboardButtonText("Help"),
        aigram.KeyboardButtonContact("Share phone"),
    },
)
replyKeyboard.ResizeKeyboard = true

_, err = b.SendMessage(ctx, aigram.SendMessageParams{
    ChatID:      aigram.ChatIDInt(123456789),
    Text:        "Pick an option",
    ReplyMarkup: replyKeyboard,
})
if err != nil {
    return err
}

_, err = b.SendMessage(ctx, aigram.SendMessageParams{
    ChatID:      aigram.ChatIDInt(123456789),
    Text:        "Keyboard removed",
    ReplyMarkup: aigram.RemoveKeyboard(false),
})
if err != nil {
    return err
}
```

Edit a message from a callback query and remove the inline keyboard:

```go
d := dispatch.New()

if err := d.OnCallbackDataFunc("confirm", func(ctx context.Context, update telegram.Update) error {
    if update.CallbackQuery == nil {
        return nil
    }

    if _, err := b.AnswerCallbackQuery(ctx, aigram.AnswerCallbackQueryParams{
        CallbackQueryID: update.CallbackQuery.ID,
        Text:            "Готово",
    }); err != nil {
        return err
    }

    msg := update.CallbackQuery.Message
    if msg == nil {
        return nil
    }

    if _, err := b.EditMessageText(ctx, aigram.EditMessageTextParams{
        Target: aigram.EditTargetChat(aigram.ChatIDInt(msg.Chat.ID), msg.MessageID),
        Text:   "Confirmed",
    }); err != nil {
        return err
    }

    _, err := b.EditMessageReplyMarkup(ctx, aigram.EditMessageReplyMarkupParams{
        Target: aigram.EditTargetChat(aigram.ChatIDInt(msg.Chat.ID), msg.MessageID),
        // nil ReplyMarkup removes the inline keyboard.
    })
    return err
}); err != nil {
    return err
}
```

Edit a media caption and delete a message:

```go
captionResult, err := b.EditMessageCaption(ctx, aigram.EditMessageCaptionParams{
    Target:  aigram.EditTargetChat(aigram.ChatIDInt(123456789), 42),
    Caption: "Updated caption",
})
if err != nil {
    return err
}
fmt.Println("caption edit ok:", captionResult.IsOK())

deleted, err := b.DeleteMessage(ctx, aigram.DeleteMessageParams{
    ChatID:    aigram.ChatIDInt(123456789),
    MessageID: 42,
})
if err != nil {
    return err
}
fmt.Println("deleted:", deleted)
```

Reply markup currently supports inline keyboards, reply keyboards, keyboard removal, and force reply for send methods. Edit methods intentionally accept only inline keyboard markup. `AnswerCallbackQuery` can acknowledge callback taps with a toast or alert. `editMessageMedia`, WebApp/LoginUrl buttons, payments, and a keyboard builder DSL will be added separately later.

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

Send video, audio, and voice messages:

```go
videoMessage, err := b.SendVideo(ctx, aigram.SendVideoParams{
    ChatID:            aigram.ChatIDInt(123456789),
    Video:             aigram.FileID("existing-video-file-id"),
    Caption:           "Video from file_id",
    SupportsStreaming: true,
})
if err != nil {
    return err
}
fmt.Println(videoMessage.MessageID)

audioMessage, err := b.SendAudio(ctx, aigram.SendAudioParams{
    ChatID:    aigram.ChatIDInt(123456789),
    Audio:     aigram.FileURL("https://example.com/audio.mp3"),
    Performer: "Example Artist",
    Title:     "Example Track",
})
if err != nil {
    return err
}
fmt.Println(audioMessage.MessageID)

voiceMessage, err := b.SendVoice(ctx, aigram.SendVoiceParams{
    ChatID: aigram.ChatIDInt(123456789),
    Voice:  aigram.FileID("existing-voice-file-id"),
})
if err != nil {
    return err
}
fmt.Println(voiceMessage.MessageID)
```

Upload a video from `os.File`:

```go
videoFile, err := os.Open("video.mp4")
if err != nil {
    return err
}
defer videoFile.Close()

uploadedVideo, err := b.SendVideo(ctx, aigram.SendVideoParams{
    ChatID: aigram.ChatIDInt(123456789),
    Video: aigram.FileUpload(aigram.UploadFile{
        Name:        "video.mp4",
        Reader:      videoFile,
        ContentType: "video/mp4",
    }),
    SupportsStreaming: true,
})
if err != nil {
    return err
}
fmt.Println(uploadedVideo.MessageID)
```

`FileID` and `FileURL` are sent as JSON requests. `FileUpload` uses multipart/form-data and ai-gram generates the internal `attach://` value for the file field. The library consumes `UploadFile.Reader` but does not close it; the caller owns reader lifecycle. Thumbnail upload, media groups, `sendAnimation`, and `sendVideoNote` are not implemented yet.

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

Telegram types currently support decoding incoming media and helper methods for handling them. Sending is currently available for text, photo, document, video, audio, and voice methods; media groups and other media send methods will be added separately later.

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

For large files pass an `*os.File` or another streaming `io.Writer` instead of `bytes.Buffer`. Telegram download URLs contain the bot token; ai-gram builds them internally and does not expose them as a public API. The regular cloud Bot API has Telegram-side file download limits. Upload is currently implemented for `SendPhoto`, `SendDocument`, `SendVideo`, `SendAudio`, and `SendVoice`; download helpers never expose a full token-bearing download URL.

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

Webhook management is JSON-only for now. Webhook certificate upload, media groups, thumbnails, editMessageMedia, answerInlineQuery, WebApp/LoginUrl buttons, payments, sendAnimation, sendVideoNote, FSM, scenes, storage, dependency injection, and full Bot API coverage are not implemented yet.


## Examples

Runnable examples are available under `examples/`:

- `examples/echo_longpoll` — basic long polling echo bot.
- `examples/inline_longpoll` — inline keyboard callbacks with `AnswerCallbackQuery`.
- `examples/webhook_server` — inbound webhook server with `SetWebhook`, safe action logs, callback edit/delete flows, and caption edit smoke.
- `examples/media_upload` — document upload and file download smoke checks.
- `examples/local_api_server` — connectivity check for a local Telegram Bot API server.

## Manual testing

Manual smoke testing instructions are in [`docs/MANUAL_TESTING.md`](docs/MANUAL_TESTING.md). The examples require real environment variables at runtime, but they are written so `go test ./...` can compile them without a token.

Deployment-oriented manual integration checks are described in [`docs/DEPLOY_TESTING.md`](docs/DEPLOY_TESTING.md). The deploy harness can start from a minimal `.env.local` with bot token, chat ID, and SSH alias, then write discovered values to ignored `.deploy/generated.env`. Smoke scripts can open a temporary SSH tunnel when a discovered local Bot API server listens only on a remote loopback; the Bot API server may live on a separate SSH target from the webhook deploy target. Remote logs are redacted before printing.

The integration harness supports role-specific test bot tokens (`MAIN`, `LOCAL`, `WEBHOOK`, `NOTIFY`, and others) while preserving the legacy single-token `AIGRAM_BOT_TOKEN` mode. Set `AIGRAM_BOTAPI_SSH_TARGET` when the local Telegram Bot API server runs on a different SSH host than the webhook example.

Manual smoke scripts can also send actionable Telegram notifications with the target `@username`, `t.me` link, exact commands/buttons, and what Codex will verify in safe logs. See [`docs/DEPLOY_TESTING.md`](docs/DEPLOY_TESTING.md#telegram-notifications-during-smoke-checks).

## Development checks

```bash
gofmt -w .
go test ./...
go vet ./...
```
