# Deploy testing harness

## Overview

This document describes a manual integration harness for ai-gram examples. It is intended for local smoke checks and for deploying the `examples/webhook_server` test bot to a server under systemd.

The harness does not contain real tokens, IP addresses, domains, or secrets. Runtime configuration is read from `.env.local`, which must stay uncommitted.

## Create `.env.local`

Copy the template and fill only the values needed for your smoke check:

```bash
cp .env.example .env.local
chmod 600 .env.local
```

Do not paste the bot token into prompts or commit `.env.local`. Let Codex read the local env file when it needs to run the scripts.

## Local Telegram Bot API server smoke

Use this when `AIGRAM_BASE_URL` points to a local Telegram Bot API server, for example `http://127.0.0.1:8081`:

```bash
./scripts/smoke_local_api.sh
```

The script runs `examples/local_api_server`, which calls `GetMe` and `GetWebhookInfo`.

## Long polling smoke

Use this for normal `getUpdates` delivery. The long polling example calls `DeleteWebhook` with `drop_pending_updates=true` before polling because Telegram webhook and long polling modes are mutually exclusive.

```bash
./scripts/smoke_longpoll.sh
```

The script runs `examples/inline_longpoll`, so `/start` should send an inline keyboard and callback taps should be acknowledged with `AnswerCallbackQuery`.

## Media smoke

Set `AIGRAM_CHAT_ID` and at least one of `AIGRAM_MEDIA_PATH` or `AIGRAM_FILE_ID`:

```bash
./scripts/smoke_media.sh
```

- `AIGRAM_MEDIA_PATH` uploads a local file through `SendDocument` and `FileUpload`.
- `AIGRAM_FILE_ID` calls `GetFile` and downloads the file into memory through `DownloadFile`.

## Deploy webhook example

The recommended deploy target configuration is an SSH alias from `~/.ssh/config`. Example:

```sshconfig
Host vk1
    HostName example.invalid
    User deploy
    IdentityFile ~/.ssh/id_ed25519
```

Then `.env.local` can reference only the alias:

```bash
AIGRAM_BOT_TOKEN=
AIGRAM_WEBHOOK_URL=
AIGRAM_DEPLOY_SSH_TARGET=vk1
AIGRAM_DEPLOY_DIR=/opt/aigram-test
AIGRAM_SERVICE_NAME=aigram-webhook-test
AIGRAM_REMOTE_ENV_DIR=/etc/aigram
AIGRAM_LISTEN_ADDR=:8090
AIGRAM_WEBHOOK_SECRET=
```

The fallback mode still supports explicit host/user/key variables when no alias is set:

```bash
AIGRAM_DEPLOY_HOST=
AIGRAM_DEPLOY_USER=
AIGRAM_DEPLOY_SSH_KEY=
```

If `AIGRAM_DEPLOY_SSH_TARGET` is set, the deploy scripts run `ssh <target>` and `scp ... <target>:...` directly and do not require `AIGRAM_DEPLOY_HOST`, `AIGRAM_DEPLOY_USER`, or `AIGRAM_DEPLOY_SSH_KEY`.

Then run:

```bash
./scripts/deploy_webhook_example.sh
```

The deploy script:

1. Builds `examples/webhook_server` for `linux/amd64`.
2. Uploads the binary to the server.
3. Writes `/etc/aigram/<service>.env` with `chmod 600`.
4. Installs a systemd service from `deploy/systemd/aigram-example.service.tmpl`.
5. Runs `systemctl daemon-reload`, `enable`, `restart`, `status`.
6. Prints the latest journal logs.

It does not delete webhook registration automatically. Use `DeleteWebhook` or another explicit maintenance step when you want to switch back to long polling.

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

- Never commit `.env.local` or real tokens.
- Do not paste the bot token into prompts; prefer giving Codex local access to `.env.local`.
- Do not log the bot token or full token-bearing Bot API/download URLs.
- The server token is stored in `/etc/aigram/*.env` with `chmod 600`.
- Prefer `AIGRAM_DEPLOY_SSH_TARGET` for existing SSH aliases; use explicit host/user/key only as a fallback.
- `AIGRAM_WEBHOOK_SECRET` must match both `SetWebhook` and `webhook.Config`.
- Official Telegram webhook delivery requires a public HTTPS URL.
- A local Telegram Bot API server running in `--local` mode can use HTTP/local webhook URLs when `AIGRAM_BASE_URL` points to that server.

## Troubleshooting

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
