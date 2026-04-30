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

## Multiple test bot tokens

The harness supports role-specific bot tokens so different smoke checks do not fight over one bot state. A single `AIGRAM_BOT_TOKEN` still works as a legacy/default fallback, but role-specific tokens are safer for repeated integration checks.

Recommended roles:

```bash
# legacy fallback
AIGRAM_BOT_TOKEN=

# general manual checks, media, examples
AIGRAM_BOT_TOKEN_MAIN=

# official api.telegram.org checks
AIGRAM_BOT_TOKEN_CLOUD=

# local Telegram Bot API server and long polling checks
AIGRAM_BOT_TOKEN_LOCAL=

# webhook deploy/service checks
AIGRAM_BOT_TOKEN_WEBHOOK=

# logOut/close/migration checks that can trigger cooldowns
AIGRAM_BOT_TOKEN_MIGRATION=

# deleteWebhook/drop_pending_updates and experiments that can lose pending updates
AIGRAM_BOT_TOKEN_DESTRUCTIVE=

# stable bot for operator notifications
AIGRAM_BOT_TOKEN_NOTIFY=
```

Fallback rules:

- `main`: `AIGRAM_BOT_TOKEN_MAIN` -> `AIGRAM_BOT_TOKEN`
- `cloud`: `AIGRAM_BOT_TOKEN_CLOUD` -> `AIGRAM_BOT_TOKEN_MAIN` -> `AIGRAM_BOT_TOKEN`
- `local`: `AIGRAM_BOT_TOKEN_LOCAL` -> `AIGRAM_BOT_TOKEN_MAIN` -> `AIGRAM_BOT_TOKEN`
- `webhook`: `AIGRAM_BOT_TOKEN_WEBHOOK` -> `AIGRAM_BOT_TOKEN_MAIN` -> `AIGRAM_BOT_TOKEN`
- `notify`: `AIGRAM_BOT_TOKEN_NOTIFY` -> `AIGRAM_BOT_TOKEN_MAIN` -> `AIGRAM_BOT_TOKEN`

Migration and destructive roles intentionally do not fall back to the default token unless explicitly allowed:

```bash
AIGRAM_ALLOW_DEFAULT_TOKEN_FOR_MIGRATION=1
AIGRAM_ALLOW_DEFAULT_TOKEN_FOR_DESTRUCTIVE=1
```

Why this matters:

- Long polling and webhook modes are mutually exclusive for one bot token.
- `logOut`/`close`/migration-style checks can have cooldowns and side effects.
- Destructive checks can delete webhooks or drop pending updates.
- Notification delivery should stay stable even while test bots are being reconfigured.

`scripts/discover_env.sh` and local Bot API smoke use the `local` role. `scripts/deploy_webhook_example.sh` writes the `webhook` role token to the remote systemd env file. `.deploy/generated.env` never stores bot tokens; tokens stay in `.env.local` or in the remote service env file created during deploy.

## Auto-discovery

Run discovery first:

```bash
./scripts/discover_env.sh
```

The script:

1. Uses `AIGRAM_DEPLOY_SSH_TARGET` directly, or legacy `AIGRAM_DEPLOY_USER@AIGRAM_DEPLOY_HOST` as fallback.
2. Uses `AIGRAM_BOTAPI_SSH_TARGET` for the local Telegram Bot API server when it is set; otherwise it falls back to the deploy SSH target.
3. Applies defaults:
   - `AIGRAM_DEPLOY_DIR=/opt/aigram-test`
   - `AIGRAM_SERVICE_NAME=aigram-webhook-test`
   - `AIGRAM_REMOTE_ENV_DIR=/etc/aigram`
   - `AIGRAM_LISTEN_ADDR=:8090`
   - `AIGRAM_BOTAPI_PORT=8081`
   - `AIGRAM_BOTAPI_BIND_ADDR=127.0.0.1`
   - `AIGRAM_BOTAPI_SERVICE_NAME=telegram-bot-api`
4. Checks SSH connectivity to both deploy and Bot API targets.
5. Checks outbound HTTPS reachability to `api.telegram.org` from the Bot API target.
6. Tries to discover a local Telegram Bot API server on the Bot API target:
   - `http://<AIGRAM_BOTAPI_BIND_ADDR>:<AIGRAM_BOTAPI_PORT>`
   - `http://127.0.0.1:8081`
   - `http://127.0.0.1:8080`
   - `http://localhost:8081`
   - `http://localhost:8080`
7. Computes `AIGRAM_FILE_BASE_URL` as `<base>/file` when a local base URL is found.
8. Computes `AIGRAM_WEBHOOK_URL` for local Bot API mode when deploy and Bot API targets are the same. If the targets differ, discovery keeps the Bot API values but reports that an explicit webhook URL is required for deploy.
9. Generates `AIGRAM_WEBHOOK_SECRET` if it is not set.
10. Writes `.deploy/generated.env` and prints a token-safe summary.

`.env.local` has priority over `.deploy/generated.env`. To override auto-discovery, set the variable explicitly in `.env.local` and rerun `./scripts/discover_env.sh`.


## Separate Bot API server host

Use this when the webhook example is deployed on one server, but the local Telegram Bot API server runs on another server. This is useful when the deploy host cannot reach `api.telegram.org`, while another host can run `telegram-bot-api` in local mode.

Example `.env.local`:

```bash
AIGRAM_DEPLOY_SSH_TARGET=vk1
AIGRAM_BOTAPI_SSH_TARGET=tgapi1
AIGRAM_BOTAPI_PORT=8081
AIGRAM_BOTAPI_BIND_ADDR=127.0.0.1
```

If `AIGRAM_BOTAPI_SSH_TARGET` is empty, the harness uses `AIGRAM_DEPLOY_SSH_TARGET` as the Bot API target. The local Bot API base URL is checked on the Bot API target, not necessarily on the deploy target. When the Bot API server listens on its remote loopback, local smoke scripts open the temporary SSH tunnel to `AIGRAM_BOTAPI_SSH_TARGET`. Discovery can still be used for local Bot API smoke without a webhook URL; deploy requires `AIGRAM_WEBHOOK_URL` when the targets differ.

Check a candidate Bot API host without changing it:

```bash
./scripts/check_botapi_host.sh
```

The check reports SSH reachability, DNS/HTTPS reachability to `api.telegram.org`, local listeners on ports `8081`/`8080`, default binary locations, and systemd status for `AIGRAM_BOTAPI_SERVICE_NAME`. It does not call `logOut`, `close`, or any migration operation.

To prepare a `telegram-bot-api` systemd service, first run the setup script without confirmation to see the plan:

```bash
./scripts/setup_botapi_service.sh
```

Apply the plan only after setting the required Telegram application credentials and explicit confirmation:

```bash
TELEGRAM_API_ID=... TELEGRAM_API_HASH=... AIGRAM_CONFIRM_SETUP_BOTAPI=1 ./scripts/setup_botapi_service.sh
```

The setup script stores credentials in an env file, defaults to `/etc/aigram/telegram-bot-api.env`, with `chmod 600`. It copies the discovered binary into the configured workdir, defaults to `/opt/telegram-bot-api`, creates a systemd service, and starts it. It never calls `logOut` or `close`.

When `AIGRAM_BOTAPI_SSH_TARGET` differs from `AIGRAM_DEPLOY_SSH_TARGET`, do not use `http://127.0.0.1:<port>/webhook` for `AIGRAM_WEBHOOK_URL`; that loopback would point at the Bot API server host, not the webhook service host. Set `AIGRAM_WEBHOOK_URL` explicitly to an HTTP/HTTPS URL reachable from the Bot API host. Also ensure `AIGRAM_BASE_URL` and `AIGRAM_FILE_BASE_URL` are reachable from the deployed webhook service if it needs to answer messages through the Bot API server.

## Local Telegram Bot API server smoke

After discovery, run:

```bash
./scripts/smoke_local_api.sh
```

The script runs `examples/local_api_server`, which calls `GetMe` and `GetWebhookInfo`. If `AIGRAM_BASE_URL` is still unknown, it will invoke discovery and then fail with a clear message if no local Bot API server is found.

If discovery selected a loopback base URL such as `http://127.0.0.1:8081` on the Bot API SSH target, the smoke script checks whether that URL is reachable locally. When it is not reachable, it opens a temporary SSH tunnel to the Bot API SSH target, rewrites `AIGRAM_BASE_URL`/`AIGRAM_FILE_BASE_URL` for the current process, and closes the tunnel on exit.

## Long polling smoke

Long polling can use the official Telegram API or the discovered local Bot API server. The example calls `DeleteWebhook` with `drop_pending_updates=true` before polling because Telegram webhook and long polling modes are mutually exclusive.

When `AIGRAM_BASE_URL` points to a remote loopback discovered on the Bot API SSH target, the script uses the same temporary SSH tunnel mechanism as `smoke_local_api.sh`.

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

When `AIGRAM_BASE_URL`/`AIGRAM_FILE_BASE_URL` point to a remote loopback discovered on the Bot API SSH target, the script opens a temporary SSH tunnel to that target and rewrites those URLs only for the smoke process.

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
AIGRAM_BOTAPI_SSH_TARGET="tgapi1"
AIGRAM_BOTAPI_PORT="8081"
AIGRAM_BOTAPI_BIND_ADDR="127.0.0.1"
```

If no local Telegram Bot API server is used, the library falls back to official `api.telegram.org`. In that mode Telegram requires an explicit public HTTPS webhook URL:

```bash
AIGRAM_WEBHOOK_URL="https://example.com/telegram/webhook"
```

HTTP webhook URLs are acceptable only for local Telegram Bot API server mode. If the Bot API target and deploy target are different, the webhook URL must be explicitly reachable from the Bot API target and must not be a deploy-host loopback URL.

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

## Telegram notifications during smoke checks

Smoke/deploy scripts can send short Telegram notifications to the operator when manual action is needed. Notifications are best-effort by default: a failed notification prints a safe warning but does not fail the smoke script.

Optional env:

```bash
AIGRAM_BOT_TOKEN_NOTIFY=
AIGRAM_BOT_TOKEN_MAIN=
AIGRAM_NOTIFY_ENABLED=1
AIGRAM_NOTIFY_STRICT=0
AIGRAM_SMOKE_MODE=targeted
```

Token selection order:

1. `AIGRAM_BOT_TOKEN_NOTIFY`
2. `AIGRAM_BOT_TOKEN_MAIN`
3. `AIGRAM_BOT_TOKEN`

`AIGRAM_CHAT_ID` is required for notifications. Set `AIGRAM_NOTIFY_ENABLED=0` to disable all script notifications. Set `AIGRAM_NOTIFY_STRICT=1` when notification delivery should fail the script.

The scripts currently notify about these manual checkpoints:

- `smoke_longpoll.sh`: sends the selected local bot `@username`, a `https://t.me/<username>` link, asks you to send any message or `/start` within `AIGRAM_SMOKE_WAIT_SECONDS` seconds (default `120`), then sends a completion notification.
- `smoke_local_api.sh`: reports successful local Bot API smoke.
- `smoke_media.sh`: sends the selected main bot `@username` and a `https://t.me/<username>` link before asking you to check media delivery/download, or reports that media smoke was skipped until `AIGRAM_MEDIA_PATH` or `AIGRAM_FILE_ID` is set.
- `deploy_webhook_example.sh`: after successful webhook deploy, uses `AIGRAM_SMOKE_MODE` to decide whether to send a full manual checklist, a deploy-done FYI, or no manual-action notification.

`AIGRAM_SMOKE_MODE` controls webhook deploy notifications:

- `targeted` (default): deploy script sends only an FYI: webhook example is deployed and no action is required unless the current Codex stage sends a separate targeted smoke request.
- `full`: deploy script sends the full manual regression checklist (`/start`, `Edit message`, `Remove keyboard`, `Caption demo`, `Edit caption`, `Delete media message`, `Delete this message`). Use this only when you intentionally want a full manual regression.
- `none`: deploy script does not send manual-action notifications after successful deploy. The global Codex final report may still be sent separately through `~/.codex/bin/codex-notify`.

For targeted checks, Codex or a script should send a separate targeted notification. Examples:

- Reply smoke: send one ordinary text message to the bot.
- Edit smoke: `/start` → `Edit message` → `Remove keyboard`.
- Caption smoke: `/start` → `Caption demo` → `Edit caption` → `Delete media message`.
- Delete smoke: `/start` → `Delete this message`.

If Codex asks for a targeted action, do only the listed steps, not the full checklist.

Notifications are sent through `examples/notify_user`, which uses the ai-gram `SendMessage` method. Scripts do not print bot tokens or `/bot<TOKEN>/sendMessage` URLs.

## Actionable Telegram notifications

Manual smoke notifications are intentionally actionable: they include the bot username, a `t.me` link, the command to send, the buttons to press, and what Codex will verify in safe logs. The operator should not need to search for the bot or remember the smoke sequence from this document during a live check.

Notifications now prefer deep links:

- `https://t.me/<bot>?start=smoke` opens the normal smoke keyboard.
- `https://t.me/<bot>?start=access_panel` opens the access-control panel.
- Telegram deep links only pass a `/start` payload; the notify bot does not and cannot send commands to another bot.
- If a deep link does not work in the Telegram client, send the fallback command manually, for example `/start access_panel` or `/start smoke`.

By default, deploy notifications are not full checklists. `AIGRAM_SMOKE_MODE=targeted` means the deploy script only says that the service was deployed; a separate targeted notification tells you exactly which current-stage action to perform. Set `AIGRAM_SMOKE_MODE=full` only for full manual regression, or `AIGRAM_SMOKE_MODE=none` to suppress deploy manual-action prompts.

`AIGRAM_TARGETED_SMOKE` can make deploy send a stage-specific deep-link notification:

- `access` — open `access_panel`.
- `reply`, `edit`, `caption`, `forward_copy` — open the smoke keyboard and describe only the current-stage buttons.
- `full` — open the smoke keyboard for a full manual regression.
- `none` — default deploy-done FYI without asking for manual actions.

If username discovery fails, scripts still send a notification with `username unknown` and continue without exposing the token. In that case, check the selected token role and `GetMe` connectivity.

## Security/access

Test bot usernames can be discovered or shared accidentally, so deployed examples run in admin-only mode by default.

Environment:

```bash
AIGRAM_ACCESS_MODE=admin
AIGRAM_ADMIN_USER_IDS=123456789
AIGRAM_ALLOWED_USER_IDS=
AIGRAM_ALLOWED_CHAT_IDS=
```

Rules:

- `AIGRAM_ACCESS_MODE=admin` is the default for `examples/webhook_server` and `examples/inline_longpoll`.
- If `AIGRAM_ADMIN_USER_IDS` is empty, examples fall back to numeric `AIGRAM_CHAT_ID` as admin user/chat for private smoke checks.
- `AIGRAM_ALLOWED_USER_IDS` and `AIGRAM_ALLOWED_CHAT_IDS` are temporary allow lists for admin mode.
- Do not use `public` or `off` for long-running public test bots. If access must be opened temporarily, use `/access_open` and then `/access_close`.
- `/access_status`, `/access_open`, and `/access_close` are admin-only commands even while runtime mode is public.
- Prefer the deep-link access panel in live smoke: `https://t.me/<bot>?start=access_panel`.
- Safe logs can include `action=access_status`, `action=access_mode_changed`, and `action=access_denied`, but they must not include tokens, webhook secrets, full message text, or admin lists.

## Security notes

- Never commit `.env.local`, `.deploy/generated.env`, or real tokens.
- Do not paste the bot token into prompts; prefer giving Codex local access to `.env.local`.
- Do not log the bot token or full token-bearing Bot API/download URLs.
- Discovery checks Bot API candidates without printing `/bot<TOKEN>/getMe` URLs.
- The server token is stored in `/etc/aigram/*.env` with `chmod 600`.
- Prefer `AIGRAM_DEPLOY_SSH_TARGET` for existing SSH aliases; use explicit host/user/key only as a fallback.
- `AIGRAM_WEBHOOK_SECRET` must match both `SetWebhook` and `webhook.Config`; discovery writes one value used by both.
- Official Telegram webhook delivery requires a public HTTPS URL.
- A local Telegram Bot API server running in `--local` mode can use HTTP/local webhook URLs when `AIGRAM_BASE_URL` points to that server and the webhook URL is reachable from the Bot API host.

## Troubleshooting

### TUN xray and local-vs-remote networking

The user's local machine may route traffic through TUN xray. Because of that, local network checks can differ from checks executed through `ssh vk1`. Do not treat that difference as a library bug by itself.

Important distinctions:

- `127.0.0.1` and `localhost` always mean the machine where the command runs.
- Local `127.0.0.1` is the user's workstation; `ssh vk1 'curl http://127.0.0.1:8081/...'` is the `vk1` loopback.
- A local Telegram Bot API server can be reachable on `vk1` as `http://127.0.0.1:8081` while being unreachable directly from the user's machine.
- Smoke scripts run locally. When discovery selects a remote loopback `AIGRAM_BASE_URL`, the scripts now open a temporary SSH tunnel automatically to the Bot API SSH target, for example local `127.0.0.1:18081 -> tgapi1:127.0.0.1:8081`.
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

If the Bot API SSH target can reach the Bot API but the local workstation cannot, the smoke scripts should create and clean up a temporary SSH tunnel automatically. For manual debugging, use a tunnel like this and stop it when done:

```bash
ssh -N -L 127.0.0.1:18081:127.0.0.1:8081 tgapi1
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
