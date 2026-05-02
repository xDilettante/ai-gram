# 06 Webhook Basic

A minimal webhook bot skeleton for ai-gram beginners.

## What it does

- reads `AIGRAM_BOT_TOKEN`;
- reads `AIGRAM_WEBHOOK_URL` and calls `SetWebhook`;
- optionally uses `AIGRAM_WEBHOOK_SECRET`;
- listens on `AIGRAM_LISTEN_ADDR`, default `:8080`;
- serves `GET /healthz` for health checks;
- serves `POST /webhook` for Telegram updates;
- replies to `/start` and text messages.

## Important

For the official Telegram Bot API, `AIGRAM_WEBHOOK_URL` must be a public HTTPS URL that forwards to this server's `/webhook` path. Local HTTP URLs are only useful with a local Bot API server or a tunnel that provides public HTTPS.

## Run

```bash
export AIGRAM_BOT_TOKEN='123456:replace_me'
export AIGRAM_WEBHOOK_URL='https://your-public-host.example/webhook'
export AIGRAM_WEBHOOK_SECRET='replace_me_secret'
export AIGRAM_LISTEN_ADDR=':8080'
go run ./examples/06_webhook_basic
```

Check locally:

```bash
curl http://127.0.0.1:8080/healthz
```

Stop with `Ctrl+C`. The example does not delete the webhook automatically on shutdown.
