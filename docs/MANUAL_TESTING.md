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
| `AIGRAM_ALLOWED_USER_IDS` | Optional comma-separated non-admin user IDs for protected examples. |
| `AIGRAM_ALLOWED_CHAT_IDS` | Optional comma-separated chat IDs for protected examples. |
| `AIGRAM_ACCESS_MODE` | Optional access mode for protected examples: `admin`, `public`, or `off`. Defaults to `admin`. |
| `AIGRAM_BASE_URL` | Optional custom Bot API base URL, for example a local Telegram Bot API server. |
| `AIGRAM_FILE_BASE_URL` | Optional file base URL for a custom Bot API endpoint. |
| `AIGRAM_LISTEN_ADDR` | HTTP listen address for webhook examples. Defaults to `:8080`. |
| `AIGRAM_WEBHOOK_URL` | Webhook URL passed to Telegram by the webhook example. |
| `AIGRAM_WEBHOOK_SECRET` | Optional secret token used by both `SetWebhook` and the webhook handler. |
| `AIGRAM_PHOTO_PATH` | Optional local photo path for the media upload example. |
| `AIGRAM_PHOTO_FILE_ID` | Optional Telegram photo `file_id` for the media upload example. |
| `AIGRAM_DOCUMENT_PATH` | Optional local document path for the media upload example. |
| `AIGRAM_DOCUMENT_FILE_ID` | Optional Telegram document `file_id` for the media upload example. |
| `AIGRAM_RETRY_TEXT` | Optional message text for the retry sender example. |
| `AIGRAM_RETRY_MAX_ATTEMPTS` | Optional retry sender attempt limit. Defaults to `4`, maximum `10`. |
| `AIGRAM_RETRY_ATTEMPT_TIMEOUT` | Optional per-attempt timeout for the retry sender example. Defaults to `10s`. |
| `AIGRAM_RETRY_BASE_DELAY` | Optional fallback retry delay. Defaults to `1s`. |
| `AIGRAM_RETRY_MAX_DELAY` | Optional maximum fallback retry delay. Defaults to `30s`. |

## Long polling echo bot

```bash
export AIGRAM_BOT_TOKEN='123456:replace_me'
go run ./examples/01_echo_bot
```

Checklist:

- The example deletes any active webhook before polling.
- Send a text message to the bot.
- The bot replies with the same text.
- `/start` returns the welcome message.
- Logs and smoke markers redact numeric chat/user IDs.
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
go run ./examples/06_webhook_basic
```

Checklist:

- Your reverse proxy or tunnel forwards the webhook URL to `/webhook` on the example server.
- `SetWebhook` succeeds.
- Incoming messages reach the handler and receive replies.
- Logs redact numeric chat/user IDs.

## Inline keyboard example

```bash
export AIGRAM_BOT_TOKEN='123456:replace_me'
go run ./examples/04_inline_keyboard
```

Checklist:

- Send `/start` to the bot.
- The bot shows an inline keyboard.
- Press demo buttons and confirm callback answers, message editing, and reply-markup removal work.
- Logs redact numeric chat/user IDs when advanced long-polling examples are used.

## Media upload example

```bash
export AIGRAM_BOT_TOKEN='123456:replace_me'
export AIGRAM_CHAT_ID='123456789'
export AIGRAM_DOCUMENT_PATH='./testdata/example.txt'
go run ./examples/05_media_upload
```

Optional file ID check:

```bash
export AIGRAM_DOCUMENT_FILE_ID='existing_document_file_id'
go run ./examples/05_media_upload
```

Checklist:

- Uploads use multipart when a local path is provided.
- If no document path or file ID is set, the example generates and uploads a small temporary text document.
- Logs do not print the full bot token or token-bearing URLs.

## Retry sender example

```bash
export AIGRAM_BOT_TOKEN='123456:replace_me'
export AIGRAM_CHAT_ID='123456789'
export AIGRAM_RETRY_TEXT='Retry sender smoke message'
go run ./examples/08_retry_sender
```

Checklist:

- The example sends one message and exits after success.
- `429 retry_after` responses are retried explicitly after the Telegram-provided delay.
- Network errors and per-attempt context deadlines use bounded exponential backoff.
- Forbidden, not-found, migrated-chat, and unrelated Telegram API errors are not retried by default.
- Logs redact numeric chat IDs and do not print the bot token.

## Group admin identity example

Use a private test group or supergroup. The example is read-only: it does not ban, restrict, delete, approve, or decline users.

```bash
export AIGRAM_BOT_TOKEN='123456:replace_me'
export AIGRAM_ADMIN_USER_IDS='123456789'
export AIGRAM_ALLOWED_CHAT_IDS='-1001234567890'
go run ./examples/09_group_admin
```

Checklist:

- Add the bot to a test group and send `/help`.
- `/whoami` shows the actor user or sender chat with redacted numeric IDs.
- `/chat` shows the effective chat with redacted numeric IDs.
- Reply to a message and send `/replytarget` to inspect the replied-to actor.
- `/admin_panel` works only for configured admin users and remains read-only.
- Anonymous admin and `sender_chat` messages are described as chat actors instead of invented users.
- Logs redact numeric chat/user IDs and do not print the bot token.

## Moderation skeleton example

Use a private test group or supergroup. The example is dry-run only: it does not call ban, restrict, delete, approve, or decline methods.

```bash
export AIGRAM_BOT_TOKEN='123456:replace_me'
export AIGRAM_ACCESS_MODE='public'
export AIGRAM_ADMIN_USER_IDS='123456789'
export AIGRAM_ALLOWED_CHAT_IDS='-1001234567890'
go run ./examples/10_moderation_skeleton
```

Checklist:

- Add the bot to a test group and send `/help`.
- Reply to a message and send `/report spam` to create a dry-run report.
- As a configured admin, reply to a message and send `/mod_preview spam` to inspect the dry-run moderation plan.
- `/mod_status` reports `dry_run=true` and `destructive_actions=disabled`.
- If join requests are enabled in the test group, incoming requests are logged as dry-run observations only.
- Logs and bot replies redact numeric chat/user IDs and do not print the bot token.

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
