# Deploy testing harness

## Overview

This document describes a manual integration harness for ai-gram examples. It is intended for local smoke checks and for deploying the `examples/webhook_server` test bot to a server under systemd.

The harness does not contain real tokens, IP addresses, domains, or secrets. Runtime configuration is read from `.env.local` and generated values are stored in `.deploy/generated.env`. Both files must stay uncommitted.

## Minimal `.env.local`

Copy the template and fill the minimal values:

```bash
cp .env.example .env.local
chmod 600 .env.local
```

Recommended minimal configuration:

```bash
AIGRAM_BOT_TOKEN="123456:REAL_TOKEN"
AIGRAM_CHAT_ID="123456789"
AIGRAM_DEPLOY_SSH_TARGET="vk1"
```

`AIGRAM_DEPLOY_SSH_TARGET` should be an SSH alias from `~/.ssh/config`. Example:

```sshconfig
Host vk1
    HostName example.invalid
    User deploy
    IdentityFile ~/.ssh/id_ed25519
```

Do not paste the bot token into prompts or commit `.env.local`. Let Codex read the local env file when it needs to run the scripts.

## Auto-discovery

Run discovery first:

```bash
./scripts/discover_env.sh
```

The script:

1. Uses `AIGRAM_DEPLOY_SSH_TARGET` directly, or legacy `AIGRAM_DEPLOY_USER@AIGRAM_DEPLOY_HOST` as fallback.
2. Applies defaults:
   - `AIGRAM_DEPLOY_DIR=/opt/aigram-test`
   - `AIGRAM_SERVICE_NAME=aigram-webhook-test`
   - `AIGRAM_REMOTE_ENV_DIR=/etc/aigram`
   - `AIGRAM_LISTEN_ADDR=:8090`
3. Checks SSH connectivity.
4. Tries to discover a local Telegram Bot API server on the remote host:
   - `http://127.0.0.1:8081`
   - `http://127.0.0.1:8080`
   - `http://localhost:8081`
   - `http://localhost:8080`
   - `http://<remote_ip>:8081`
   - `http://<remote_ip>:8080`
5. Computes `AIGRAM_FILE_BASE_URL` as `<base>/file` when a local base URL is found.
6. Computes `AIGRAM_WEBHOOK_URL` for local Bot API mode.
7. Generates `AIGRAM_WEBHOOK_SECRET` if it is not set.
8. Writes `.deploy/generated.env` and prints a token-safe summary.

`.env.local` has priority over `.deploy/generated.env`. To override auto-discovery, set the variable explicitly in `.env.local` and rerun `./scripts/discover_env.sh`.

## Local Telegram Bot API server smoke

After discovery, run:

```bash
./scripts/smoke_local_api.sh
```

The script runs `examples/local_api_server`, which calls `GetMe` and `GetWebhookInfo`. If `AIGRAM_BASE_URL` is still unknown, it will invoke discovery and then fail with a clear message if no local Bot API server is found.

## Long polling smoke

Long polling can use the official Telegram API or the discovered local Bot API server. The example calls `DeleteWebhook` with `drop_pending_updates=true` before polling because Telegram webhook and long polling modes are mutually exclusive.

```bash
./scripts/smoke_longpoll.sh
```

The script runs `examples/inline_longpoll`, so `/start` should send an inline keyboard and callback taps should be acknowledged with `AnswerCallbackQuery`.

## Media smoke

Set at least one of `AIGRAM_MEDIA_PATH` or `AIGRAM_FILE_ID` in `.env.local`:

```bash
./scripts/smoke_media.sh
```

- `AIGRAM_MEDIA_PATH` uploads a local file through `SendDocument` and `FileUpload`.
- `AIGRAM_FILE_ID` calls `GetFile` and downloads the file into memory through `DownloadFile`.

## Deploy webhook example

With the minimal `.env.local`, discovery can usually prepare the rest:

```bash
./scripts/discover_env.sh
./scripts/deploy_webhook_example.sh
```

If `.deploy/generated.env` does not exist, `deploy_webhook_example.sh` runs discovery automatically.

The deploy script:

1. Builds `examples/webhook_server` for `linux/amd64`.
2. Uploads the binary to the server.
3. Writes `/etc/aigram/<service>.env` with `chmod 600`.
4. Installs a systemd service from `deploy/systemd/aigram-example.service.tmpl`.
5. Runs `systemctl daemon-reload`, `enable`, `restart`, `status`.
6. Prints the latest journal logs.

It does not delete webhook registration automatically. Use `DeleteWebhook` or another explicit maintenance step when you want to switch back to long polling.

## Manual overrides

Use these overrides when auto-discovery cannot infer the environment:

```bash
AIGRAM_BASE_URL="http://127.0.0.1:8081"
AIGRAM_FILE_BASE_URL="http://127.0.0.1:8081/file"
AIGRAM_WEBHOOK_URL="http://127.0.0.1:8090/webhook"
AIGRAM_WEBHOOK_SECRET="manual_secret_123"
AIGRAM_LISTEN_ADDR=":8090"
```

If no local Telegram Bot API server is used, the library falls back to official `api.telegram.org`. In that mode Telegram requires an explicit public HTTPS webhook URL:

```bash
AIGRAM_WEBHOOK_URL="https://example.com/telegram/webhook"
```

HTTP webhook URLs are acceptable only for local Telegram Bot API server mode.

Legacy SSH fallback is still supported when no alias is set:

```bash
AIGRAM_DEPLOY_HOST=
AIGRAM_DEPLOY_USER=
AIGRAM_DEPLOY_SSH_KEY=
```

## View remote logs

```bash
./scripts/remote_logs.sh
```

This shows:

```bash
journalctl -u "$AIGRAM_SERVICE_NAME" -n 120 --no-pager
```

## Stop remote service

```bash
./scripts/remote_stop.sh
```

This stops the systemd service only. It does not remove files, the remote env file, or webhook registration.

## Security notes

- Never commit `.env.local`, `.deploy/generated.env`, or real tokens.
- Do not paste the bot token into prompts; prefer giving Codex local access to `.env.local`.
- Do not log the bot token or full token-bearing Bot API/download URLs.
- Discovery checks Bot API candidates without printing `/bot<TOKEN>/getMe` URLs.
- The server token is stored in `/etc/aigram/*.env` with `chmod 600`.
- Prefer `AIGRAM_DEPLOY_SSH_TARGET` for existing SSH aliases; use explicit host/user/key only as a fallback.
- `AIGRAM_WEBHOOK_SECRET` must match both `SetWebhook` and `webhook.Config`; discovery writes one value used by both.
- Official Telegram webhook delivery requires a public HTTPS URL.
- A local Telegram Bot API server running in `--local` mode can use HTTP/local webhook URLs when `AIGRAM_BASE_URL` points to that server.

## Troubleshooting

### Discovery cannot find local Bot API

Set `AIGRAM_BASE_URL` manually, for example:

```bash
AIGRAM_BASE_URL="http://127.0.0.1:8081"
```

Then rerun:

```bash
./scripts/discover_env.sh
```

### Discovery cannot compute webhook URL

If no local Bot API server is found, set a public HTTPS URL manually:

```bash
AIGRAM_WEBHOOK_URL="https://example.com/telegram/webhook"
```

### Long polling does not receive updates

- Check that webhook is deleted before polling.
- Run `examples/local_api_server` or call `GetWebhookInfo` to see current webhook state.
- Ensure no other process is consuming updates with the same bot token.

### Webhook requests do not arrive

- Check `AIGRAM_WEBHOOK_URL` and ensure it points to the deployed `/webhook` route.
- Check server port binding and firewall rules.
- Check `AIGRAM_WEBHOOK_SECRET`; it must match Telegram and `webhook.Config`.
- Inspect systemd logs with `./scripts/remote_logs.sh`.

### systemd service does not start

- Run `./scripts/remote_logs.sh`.
- Check that the remote env file exists and has all required variables.
- Check that the deploy user can bind `AIGRAM_LISTEN_ADDR` and read the binary.

### Local Bot API server is unavailable

- Check `AIGRAM_BASE_URL` with curl before running examples.
- Ensure the local Bot API server is listening on the configured host/port.
- For local webhook checks, ensure the Bot API server was started with local mode enabled.

### Media upload/download fails

- Check `AIGRAM_MEDIA_PATH`, file permissions, and Telegram file size limits.
- Check that `AIGRAM_CHAT_ID` is valid for upload tests.
- Check that `AIGRAM_FILE_ID` belongs to the bot or is otherwise accessible to Telegram.
