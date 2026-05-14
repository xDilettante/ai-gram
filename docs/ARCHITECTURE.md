# Architecture

`ai-gram` is organized as a library of small layers. Applications can use the low-level Bot API client directly, wire update transports into their own runtime, or add the dispatcher and middleware packages when they want routing.

The root `aigram` package is intentionally a compact convenience facade. It is not a full re-export mirror of every Bot API method parameter or Telegram object. Import `bot` for the full outgoing method surface and `telegram` for data contracts.

```mermaid
flowchart LR
    app[User code]
    root[aigram root facade]
    bot[bot package\nTyped Bot API client]
    types[telegram package\nBot API data contracts]
    callback[callback package\nCallback data helpers]
    lp[transport/longpoll\nUpdate source]
    wh[transport/webhook\nHTTP receiver]
    dispatch[dispatch package\nRoutes and predicates]
    mw[middleware package\nHandler wrappers]
    api[Telegram Bot API]

    app --> root
    root --> bot
    root --> types
    app --> lp
    app --> wh
    app --> dispatch
    app --> mw
    app --> callback
    bot --> types
    bot --> api
    api --> lp
    api --> wh
    lp --> dispatch
    wh --> dispatch
    dispatch --> mw
    mw --> app
```

## Outgoing requests

The `bot` package owns outgoing Telegram Bot API calls. It accepts typed parameter structs, validates required fields, encodes JSON or multipart requests, decodes typed results, and returns typed Telegram API errors where possible. A configurable HTTP client and base URL make the layer testable with `httptest` and usable with the official or local Telegram Bot API server.

## Incoming updates

`transport/longpoll` fetches updates through `getUpdates`. `transport/webhook` validates inbound webhook HTTP requests and decodes update JSON. Both transports can feed any handler implementing the small update handler interface, including the `dispatch.Dispatcher`.

`dispatch` routes already received `telegram.Update` values. It does not own the Bot API client, so application code can decide whether handlers call Telegram, enqueue work, or only observe updates. `middleware` wraps handlers for reusable concerns such as access control.

## Callback data

The `callback` package builds and parses compact inline keyboard `callback_data` values. It is intentionally stateless: applications own authorization, storage, and side effects, while the package provides typed namespace/action/ID/page/expiry fields, length validation, and helpers for common confirm/cancel or pagination flows.

`dispatch` can route typed callback data with namespace/action predicates and handlers that receive the parsed callback payload.

## File uploads

Telegram's official `InputFile` concept is represented by `FileRef` and `FileUpload` in the client layer. Existing file IDs and URLs stay in JSON requests when the method allows them. New uploads use multipart requests and deterministic `attach://` references for media, thumbnails, covers, webhook certificates, and other upload-capable fields.

## Public API shape

- The root `aigram` package keeps quick-start helpers and common message/reply markup types only.
- Advanced Bot API method params live in `bot`; Telegram objects live in `telegram`; typed callback payload helpers live in `callback`.
- Bot construction uses `aigram.Config` or `bot.Config`.
- `GetChat` returns the official `ChatFullInfo` result shape.
- `GetChatFullInfo` remains as a same-result alias while the project is pre-v1.
- `ChatMember` is an interface implemented by official `ChatMember*` variants.
- `CallbackQuery.Message` uses the official `MaybeInaccessibleMessage` shape, with helpers for accessible messages.

See [`PRE_V1_NOTES.md`](PRE_V1_NOTES.md) for current pre-v1 API cleanup notes.
