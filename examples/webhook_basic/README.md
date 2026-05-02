# Basic webhook example

This example runs a minimal webhook bot skeleton with:

- `/webhook` for Telegram updates;
- `/healthz` for local health checks;
- optional Telegram webhook secret validation;
- `/start` command handling;
- text message replies;
- graceful HTTP server shutdown.

## Run

Expose your local server through HTTPS using your preferred tunnel or deployment environment, then run:

```bash
export AIGRAM_BOT_TOKEN="123456:replace-me"
export AIGRAM_WEBHOOK_URL="https://example.com/webhook"
export AIGRAM_WEBHOOK_SECRET="replace-with-a-random-secret"
export AIGRAM_LISTEN_ADDR=":8080"
go run ./examples/webhook_basic
```

Optional variables:

- `AIGRAM_BASE_URL` for a local Telegram Bot API server.
- `AIGRAM_FILE_BASE_URL` for local file download URLs.

The example registers the webhook on startup and intentionally does not delete it on shutdown. Delete or replace the webhook explicitly when you switch back to long polling.
