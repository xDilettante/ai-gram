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

If discovery selected a loopback base URL such as `http://127.0.0.1:8081` on `vk1`, the smoke script checks whether that URL is reachable locally. When it is not reachable, it opens a temporary SSH tunnel to `vk1`, rewrites `AIGRAM_BASE_URL`/`AIGRAM_FILE_BASE_URL` for the current process, and closes the tunnel on exit.

## Long polling smoke

Long polling can use the official Telegram API or the discovered local Bot API server. The example calls `DeleteWebhook` with `drop_pending_updates=true` before polling because Telegram webhook and long polling modes are mutually exclusive.

When `AIGRAM_BASE_URL` points to a remote loopback discovered on `vk1`, the script uses the same temporary SSH tunnel mechanism as `smoke_local_api.sh`.

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

When `AIGRAM_BASE_URL`/`AIGRAM_FILE_BASE_URL` point to a remote loopback discovered on `vk1`, the script opens a temporary SSH tunnel and rewrites those URLs only for the smoke process.

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

This shows the latest systemd journal lines for the service:

```bash
journalctl -u "$AIGRAM_SERVICE_NAME" -n 120 --no-pager
```

The script filters output before printing it locally: configured bot token, webhook secret, and generic `/bot<TOKEN>/...` endpoints are redacted. The deploy script applies the same filter to `systemctl status` and journal output printed after restart.

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

### TUN xray and local-vs-remote networking

The user's local machine may route traffic through TUN xray. Because of that, local network checks can differ from checks executed through `ssh vk1`. Do not treat that difference as a library bug by itself.

Important distinctions:

- `127.0.0.1` and `localhost` always mean the machine where the command runs.
- Local `127.0.0.1` is the user's workstation; `ssh vk1 'curl http://127.0.0.1:8081/...'` is the `vk1` loopback.
- A local Telegram Bot API server can be reachable on `vk1` as `http://127.0.0.1:8081` while being unreachable directly from the user's machine.
- Smoke scripts run locally. When discovery selects a remote loopback `AIGRAM_BASE_URL`, the scripts now open a temporary SSH tunnel automatically, for example local `127.0.0.1:18081 -> vk1:127.0.0.1:8081`.
- A webhook URL must be reachable by the component that sends webhook requests: official Telegram, the local Telegram Bot API server, or a synthetic local probe.

Recommended checks:

```bash
ssh vk1 'ss -tulpn | grep -E ":(8080|8081|8090)\b" || true'
ssh vk1 'systemctl status telegram-bot-api --no-pager || true'
ssh vk1 'systemctl status aigram-webhook-test --no-pager || true'
```

Check local access separately:

```bash
curl -sS --max-time 2 http://127.0.0.1:8081/ || true
```

If `vk1` can reach the Bot API but the local workstation cannot, the smoke scripts should create and clean up a temporary SSH tunnel automatically. For manual debugging, use a tunnel like this and stop it when done:

```bash
ssh -N -L 127.0.0.1:18081:127.0.0.1:8081 vk1
```

When a network check fails, first record where it ran: locally or through `ssh vk1`. Then check routing/TUN, tunnel state, firewall/listening ports, systemd logs, and `getWebhookInfo` before changing library code. Never print the bot token, webhook secret, token-bearing URLs, or full `/bot<TOKEN>/...` endpoints while debugging.

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
