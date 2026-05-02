# Manual Testing

This page is the public manual-testing guide for ai-gram users. It focuses on small local checks that are useful when trying the library with your own bot. Maintainer-only deploy, multi-bot, notification, and destructive live-smoke procedures live under [`docs/maintainer/`](maintainer/).

All examples read configuration from environment variables and never hardcode bot tokens. Do not commit real tokens, webhook secrets, private chat IDs, payment payloads, or message contents.

## Environment variables

Start from the minimal public template:

```bash
cp .env.example .env.local
```

Common variables:

| Variable | Purpose |
| --- | --- |
| `AIGRAM_BOT_TOKEN` | Telegram bot token. Required for examples that call Telegram. |
| `AIGRAM_CHAT_ID` | Chat ID or username for examples that send messages proactively. |
| `AIGRAM_ADMIN_USER_IDS` | Optional comma-separated admin user IDs for protected examples. |
| `AIGRAM_BASE_URL` | Optional custom Bot API base URL, for example a local Telegram Bot API server. |
| `AIGRAM_FILE_BASE_URL` | Optional file base URL for a custom Bot API endpoint. |
| `AIGRAM_LISTEN_ADDR` | HTTP listen address for webhook examples. Defaults to `:8080`. |
| `AIGRAM_WEBHOOK_URL` | Webhook URL passed to Telegram by the webhook example. |
| `AIGRAM_WEBHOOK_SECRET` | Optional secret token used by both `SetWebhook` and the webhook handler. |
| `AIGRAM_MEDIA_PATH` | Optional local file path for media upload checks. |
| `AIGRAM_FILE_ID` | Optional Telegram `file_id` for file download checks. |

## Long polling echo bot

```bash
export AIGRAM_BOT_TOKEN='123456:replace_me'
go run ./examples/echo_longpoll
```

Checklist:

- The example deletes any active webhook before polling.
- Send a text message to the bot.
- The bot replies with `echo: <your text>`.
- Stop with `Ctrl+C` and confirm graceful shutdown.

## Local Bot API connectivity

If you run a local Telegram Bot API server, point ai-gram at it:

```bash
export AIGRAM_BOT_TOKEN='123456:replace_me'
export AIGRAM_BASE_URL='http://127.0.0.1:8081'
go run ./examples/local_api_server
```

Checklist:

- The example prints the bot username and webhook status.
- It does not print the bot token or token-bearing URLs.

## Webhook server

For official `api.telegram.org`, `AIGRAM_WEBHOOK_URL` must be a public HTTPS URL. A local Bot API server running in local mode may allow local HTTP URLs.

```bash
export AIGRAM_BOT_TOKEN='123456:replace_me'
export AIGRAM_LISTEN_ADDR=':8080'
export AIGRAM_WEBHOOK_URL='https://example.com/webhook'
export AIGRAM_WEBHOOK_SECRET='replace_me_secret'
go run ./examples/webhook_server
```

Checklist:

- Your reverse proxy forwards the webhook URL to `/webhook` on the example server.
- `SetWebhook` succeeds.
- `GetWebhookInfo` shows the expected webhook and no recent error.
- Incoming messages reach the handler and receive replies.

## Inline keyboard example

```bash
export AIGRAM_BOT_TOKEN='123456:replace_me'
export AIGRAM_ADMIN_USER_IDS='123456789'
go run ./examples/inline_longpoll
```

Checklist:

- Send `/start` to the bot.
- The bot shows the demo inline keyboard.
- Press demo buttons and confirm callback answers and message edits work.
- Access-control examples stay admin-only unless you intentionally configure a public mode.

## Media upload/download example

```bash
export AIGRAM_BOT_TOKEN='123456:replace_me'
export AIGRAM_CHAT_ID='123456789'
export AIGRAM_MEDIA_PATH='./testdata/example.txt'
go run ./examples/media_upload
```

Optional download check:

```bash
export AIGRAM_FILE_ID='existing_file_id'
go run ./examples/media_upload
```

Checklist:

- Uploads use multipart when a local path is provided.
- Downloads use `GetFile` and the configured file base URL.
- Logs do not print the full bot token or token-bearing URLs.

## Sensitive and state-changing areas

The following areas should remain manual-only and should be tested only with dedicated test bots, disposable chats, or explicit maintainer approval:

- payments, Stars, gifts, and paid media;
- business APIs and managed bot token methods;
- passport data;
- admin/destructive chat methods;
- sticker set mutation;
- games and inline mode features that require BotFather setup;
- lifecycle `LogOut`/`Close`;
- `SetWebhook` certificate upload.

Use unit tests and `httptest` coverage as the default verification for these areas. If you run live checks, log only safe metadata such as method names, result booleans, and redacted IDs.

## Maintainer-only live smoke

Maintainer smoke/deploy flows are intentionally separated from the public guide:

- [`docs/maintainer/LIVE_SMOKE_MATRIX.md`](maintainer/LIVE_SMOKE_MATRIX.md) classifies safe, sensitive, and destructive live checks.
- [`docs/maintainer/DEPLOY_TESTING.md`](maintainer/DEPLOY_TESTING.md) documents the private deploy and multi-bot smoke harness.
- [`docs/maintainer/ENV_SMOKE_TEMPLATE.md`](maintainer/ENV_SMOKE_TEMPLATE.md) lists the extended maintainer environment variables.
