# Manual smoke testing

## Overview

This document describes manual smoke checks for ai-gram examples against either the official Telegram Bot API or a local Telegram Bot API server. The checks are intentionally manual: they require a real bot token and sometimes a chat with the bot, so they are not part of `go test ./...`.

All examples read configuration from environment variables and never hardcode bot tokens.

The deploy/smoke scripts can send actionable Telegram notifications with the selected bot `@username`, a `t.me` link, exact commands/buttons to press, and a short note about which safe logs Codex will verify.

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


## Smoke notification modes

Deploy and smoke scripts support `AIGRAM_SMOKE_MODE`:

- `targeted` is the default. The deploy script only reports that the webhook example was deployed; perform manual actions only when Codex sends a separate targeted notification.
- `full` sends the full webhook regression checklist after deploy. Use it only when intentionally running a full manual regression.
- `none` disables deploy manual-action prompts. Final Codex reports may still be sent through the global notification helper.

When Codex asks for a targeted smoke, do only the listed steps, not the full checklist. Common targeted flows:

- Reply smoke: send one ordinary text message.
- Edit smoke: `/start` → `Edit message` → `Remove keyboard`.
- Caption smoke: `/start` → `Caption demo` → `Edit caption` → `Delete media message`.
- Delete smoke: `/start` → `Delete this message`.

## Inline keyboard callback checklist

```bash
export AIGRAM_BOT_TOKEN='123456:replace_me'
go run ./examples/inline_longpoll
```

Checklist:

- Send `/start` to the bot.
- The bot sends an inline keyboard with `Edit message` and `Remove keyboard`.
- Press `Edit message`: the client shows a toast from `AnswerCallbackQuery`, and the original message text changes to `Message edited by ai-gram`.
- Press `Remove keyboard`: the client shows a toast from `AnswerCallbackQuery`, and the inline keyboard disappears.
- For the deployed webhook example, inspect safe logs with `./scripts/remote_logs.sh`; logs should include `update_id`, `update_type`, `chat_id`, `from_user_id`, `command`, `has_text`, `has_media`, and only known short `demo:*` callback data.
- Successful webhook actions are logged explicitly as safe action lines, for example `action=answer_callback_query`, `action=edit_message_text`, `action=edit_message_reply_markup`, and `action=send_message`.

## Webhook DeleteMessage and EditMessageCaption checklist

The deployed `examples/webhook_server` also contains a live flow for deleting messages and editing media captions.

To check `DeleteMessage` on the `/start` message:

- Deploy or run `examples/webhook_server`.
- Send `/start` to the webhook bot.
- Press `Delete this message`.
- The message with the inline keyboard should disappear.
- Inspect safe logs with `./scripts/remote_logs.sh`; successful deletion is logged as `action=delete_message ok=true update_id=... chat_id=... message_id=...`.

To check `EditMessageCaption` on a media message, no media env is required. If `AIGRAM_FILE_ID` or `AIGRAM_MEDIA_PATH` is set, the example uses it; otherwise it uploads a generated in-memory text document named `aigram-caption-demo.txt`.

```bash
export AIGRAM_FILE_ID='existing_document_file_id'
# or:
export AIGRAM_MEDIA_PATH='/path/to/file.pdf'
```

Checklist:

- Deploy or run `examples/webhook_server`.
- Send `/start` to the webhook bot.
- Press `Caption demo`.
- The bot sends a document with caption `Original caption from ai-gram` and inline buttons.
- Press `Edit caption`; the media caption should change to `Caption edited by ai-gram`.
- Press `Delete media message`; the media message should disappear.
- Inspect safe logs with `./scripts/remote_logs.sh`; successful actions are logged as `action=send_media_caption_demo` with `source=generated_document`, `source=file_id`, or `source=media_path`, plus `action=edit_message_caption` and `action=delete_message`.

## Webhook ForwardMessage and CopyMessage checklist

The deployed `examples/webhook_server` contains a targeted flow for checking single-message forwarding and copying.

Checklist:

- Deploy or run `examples/webhook_server`.
- Send `/start` to the webhook bot.
- Press `Copy this message`; Telegram should create a copy of the current message in the same chat.
- Press `Forward this message`; Telegram should create a forwarded message in the same chat.
- Inspect safe logs with `./scripts/remote_logs.sh`; successful actions are logged as `action=copy_message ok=true update_id=... chat_id=... message_id=... copied_message_id=...` and `action=forward_message ok=true update_id=... chat_id=... message_id=... forwarded_message_id=...`.
- Do not log or paste full message text; safe logs are enough for verification.

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
