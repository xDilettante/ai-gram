# Manual smoke testing

## Overview

This document describes manual smoke checks for ai-gram examples against either the official Telegram Bot API or a local Telegram Bot API server. The checks are intentionally manual: they require a real bot token and sometimes a chat with the bot, so they are not part of `go test ./...`.

All examples read configuration from environment variables and never hardcode bot tokens.

## Environment variables

| Variable | Purpose |
| --- | --- |
| `AIGRAM_BOT_TOKEN` | Telegram bot token. Required for all examples. |
| `AIGRAM_BASE_URL` | Optional Bot API base URL, for example `http://127.0.0.1:8081` for a local Bot API server. |
| `AIGRAM_FILE_BASE_URL` | Optional Bot API file base URL. Usually derived from `AIGRAM_BASE_URL` when using a local server. |
| `AIGRAM_CHAT_ID` | Chat ID or username used by examples that proactively send media/messages. |
| `AIGRAM_LISTEN_ADDR` | HTTP listen address for webhook examples. Defaults to `:8080`. |
| `AIGRAM_WEBHOOK_URL` | Webhook URL passed to Telegram in `SetWebhook`. |
| `AIGRAM_WEBHOOK_SECRET` | Optional secret token used both in `SetWebhook` and `webhook.Config`. |
| `AIGRAM_MEDIA_PATH` | Local file path for upload smoke testing. |
| `AIGRAM_FILE_ID` | Existing Telegram `file_id` for download smoke testing. |

## Long polling через api.telegram.org

```bash
export AIGRAM_BOT_TOKEN='123456:replace_me'
go run ./examples/echo_longpoll
```

Checklist:

- The example calls `DeleteWebhook(drop_pending_updates=true)` before starting polling.
- Send a text message to the bot.
- The bot replies with `echo: <your text>`.
- Stop the process with `Ctrl+C` and confirm graceful shutdown.

## Long polling через local Bot API server

Start your local Telegram Bot API server separately, then point ai-gram at it:

```bash
export AIGRAM_BOT_TOKEN='123456:replace_me'
export AIGRAM_BASE_URL='http://127.0.0.1:8081'
go run ./examples/local_api_server

go run ./examples/echo_longpoll
```

Checklist:

- `local_api_server` prints the bot username and webhook info.
- `echo_longpoll` starts without printing the token.
- Messages sent to the bot are echoed.

## Webhook через local Bot API server

For a local Bot API server running in `--local` mode, Telegram Bot API can accept local HTTP webhook URLs. ai-gram allows HTTP webhook URLs when a custom `AIGRAM_BASE_URL` is configured.

```bash
export AIGRAM_BOT_TOKEN='123456:replace_me'
export AIGRAM_BASE_URL='http://127.0.0.1:8081'
export AIGRAM_LISTEN_ADDR='127.0.0.1:8080'
export AIGRAM_WEBHOOK_URL='http://127.0.0.1:8080/webhook'
export AIGRAM_WEBHOOK_SECRET='local_secret_123'
go run ./examples/webhook_server
```

Checklist:

- The server listens on `/webhook`.
- `SetWebhook` succeeds.
- Sending `/start` or a text message to the bot reaches the local handler.
- `AIGRAM_WEBHOOK_SECRET` is the same value for `SetWebhook` and `webhook.Config`.

## Webhook через публичный HTTPS URL

Official `api.telegram.org` requires a public HTTPS webhook URL.

```bash
export AIGRAM_BOT_TOKEN='123456:replace_me'
export AIGRAM_LISTEN_ADDR=':8080'
export AIGRAM_WEBHOOK_URL='https://example.com/webhook'
export AIGRAM_WEBHOOK_SECRET='public_secret_123'
go run ./examples/webhook_server
```

Checklist:

- Your reverse proxy forwards `https://example.com/webhook` to the local `/webhook` handler.
- `SetWebhook` succeeds.
- `GetWebhookInfo` shows the expected URL and no recent error.
- Incoming messages receive replies.

## Inline keyboard callback checklist

```bash
export AIGRAM_BOT_TOKEN='123456:replace_me'
go run ./examples/inline_longpoll
```

Checklist:

- Send `/start` to the bot.
- The bot sends an inline keyboard with `Да` and `Нет`.
- Press `Да`: the client shows a toast from `AnswerCallbackQuery`, and the bot sends `Да подтверждено`.
- Press `Нет`: the client shows an alert from `AnswerCallbackQuery`.

## Media upload/download checklist

Upload a local file as a document:

```bash
export AIGRAM_BOT_TOKEN='123456:replace_me'
export AIGRAM_CHAT_ID='123456789'
export AIGRAM_MEDIA_PATH='/path/to/file.pdf'
go run ./examples/media_upload
```

Download an existing Telegram file by `file_id`:

```bash
export AIGRAM_BOT_TOKEN='123456:replace_me'
export AIGRAM_FILE_ID='existing_file_id'
go run ./examples/media_upload
```

Checklist:

- Upload prints a message ID.
- Download prints the downloaded byte count.
- The library never prints a full token-bearing download URL.
- For large files, use a streaming source/destination rather than keeping content in memory in your own application code.

## Troubleshooting

- Long polling does not receive updates: check whether a webhook is still configured with `GetWebhookInfo`; long polling and webhooks are mutually exclusive, so call `DeleteWebhook` before polling.
- Webhook updates do not arrive: verify the public/local URL, listen port, reverse proxy, secret token, and `GetWebhookInfo` last error fields.
- Webhook returns 401: make sure `AIGRAM_WEBHOOK_SECRET` matches the secret sent in `SetWebhook` and configured in `webhook.Config`.
- Media upload fails: check `AIGRAM_MEDIA_PATH`, file permissions, file size, and Telegram Bot API limits.
- File download fails: check that `AIGRAM_FILE_ID` is valid and that `GetFile` returns a non-empty `file_path`.
- Local Bot API server does not work: verify `AIGRAM_BASE_URL`, that the server process is reachable, and that the file base URL is correct for your local setup.

## Security notes

- Never commit a real bot token.
- Do not log bot tokens.
- Do not log full Bot API endpoint URLs because they contain the token.
- Do not log full Telegram file download URLs because they contain the token.
- Keep the webhook secret the same in `SetWebhook` and `webhook.Config`.
- Call `DeleteWebhook` before starting long polling.
- Official Telegram webhooks require a public HTTPS URL.
- A local Telegram Bot API server in `--local` mode can accept HTTP/local webhook URLs when `AIGRAM_BASE_URL` points to that local server.
