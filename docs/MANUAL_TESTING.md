# Manual smoke testing

## Overview

This document describes manual smoke checks for ai-gram examples against either the official Telegram Bot API or a local Telegram Bot API server. The checks are intentionally manual: they require a real bot token and sometimes a chat with the bot, so they are not part of `go test ./...`.

All examples read configuration from environment variables and never hardcode bot tokens.

The deploy/smoke scripts can send actionable Telegram notifications with the selected bot `@username`, a `t.me` link, exact commands/buttons to press, and a short note about which safe logs Codex will verify.

Smoke notifications prefer Telegram deep links such as `https://t.me/<bot>?start=access_panel` or `?start=smoke`. Telegram deep links can only pass a `/start` payload, so the target bot opens a control panel or smoke keyboard instead of receiving arbitrary commands from the notification bot.

Use [`LIVE_SMOKE_MATRIX.md`](LIVE_SMOKE_MATRIX.md) to decide which flows are safe to run automatically or manually, and which destructive/admin flows require an explicit isolated test setup.

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
| `AIGRAM_ACCESS_MODE` | Example access mode: `admin` (default), `public`, or `off`. |
| `AIGRAM_ADMIN_USER_IDS` | Comma-separated admin user IDs. Falls back to numeric `AIGRAM_CHAT_ID` when empty. |
| `AIGRAM_ALLOWED_USER_IDS` | Comma-separated user IDs allowed in admin mode. |
| `AIGRAM_ALLOWED_CHAT_IDS` | Comma-separated chat IDs allowed in admin mode. |

## Long polling through api.telegram.org

```bash
export AIGRAM_BOT_TOKEN='123456:replace_me'
go run ./examples/echo_longpoll
```

Checklist:

- The example calls `DeleteWebhook(drop_pending_updates=true)` before starting polling.
- Send a text message to the bot.
- The bot replies with `echo: <your text>`.
- Stop the process with `Ctrl+C` and confirm graceful shutdown.

## Long polling through a local Bot API server

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

## Webhook through a local Bot API server

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

## Webhook through a public HTTPS URL

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

## Access control

`examples/webhook_server` and `examples/inline_longpoll` are protected by access control middleware by default. This prevents random users who find a test bot username from running commands or pressing demo buttons.

Configuration:

```bash
export AIGRAM_ACCESS_MODE=admin
export AIGRAM_ADMIN_USER_IDS='123456789'
export AIGRAM_ALLOWED_USER_IDS=''
export AIGRAM_ALLOWED_CHAT_IDS=''
```

Rules:

- `admin` is the default mode.
- If `AIGRAM_ADMIN_USER_IDS` is empty, examples use numeric `AIGRAM_CHAT_ID` as both admin user/chat fallback for private smoke checks.
- `public` and `off` let all updates pass. Use them only for local development or a short controlled smoke.
- The examples do not disclose admin/allowed ID lists to denied users.

Runtime commands:

- `/access_status` — show current runtime mode.
- `/access_open` — switch runtime mode to `public`.
- `/access_close` — switch runtime mode back to `admin`.

Deep-link panel:

- Open `https://t.me/<bot_username>?start=access_panel`.
- The bot shows an inline control panel with `Access status`, `Open access`, `Close access`, and `Start smoke`.
- If the deep link does not work in your Telegram client, send `/start access_panel` manually.
- Safe logs include `action=start_payload payload=access_panel`, `action=access_panel_shown ok=true`, `action=access_status ok=true`, and optional `action=access_mode_changed ok=true mode=...`.

Only admin users can run `/access_*` commands. Even when access is open, a non-admin user cannot change the mode. In admin mode, unknown users receive `Access denied.` or are ignored depending on the update shape, and safe logs include `action=access_denied update_id=... chat_id=... from_user_id=...`.

Manual check:

- Start the bot as an admin and send `/access_status`; expect `Access mode: admin`.
- Send `/start`; expect the normal demo keyboard.
- Prefer the access panel deep link when it is available.
- Optional: send `/access_open`, test from another account, then send `/access_close`.
- Inspect logs with `./scripts/remote_logs.sh`; safe logs should include `action=access_status`, optional `action=access_mode_changed`, and no token, secret, or full message text.

## v0.2 send methods smoke

Use this script to verify the safe v0.2 send-method subset against a real bot without requiring user interaction:

```bash
export AIGRAM_BOT_TOKEN_MAIN='123456:replace_me'
export AIGRAM_CHAT_ID='123456789'
./scripts/smoke_v02_send_methods.sh
```

Required checks:

- `SendContact` with a fake test contact.
- `SendLocation` with neutral test coordinates.
- `SendVenue` with neutral test coordinates.
- `SendPoll`, followed by `StopPoll` for the sent poll message.
- `SendDice` with `🎲`.

Optional media checks are skipped without failing when media env is absent:

- `AIGRAM_STICKER_FILE_ID` enables `SendSticker`.
- `AIGRAM_ANIMATION_FILE_ID` or `AIGRAM_ANIMATION_PATH` enables `SendAnimation`.
- `AIGRAM_VIDEO_NOTE_FILE_ID` or `AIGRAM_VIDEO_NOTE_PATH` enables `SendVideoNote`.

Expected safe markers include `AIGRAM_V02_SMOKE_SEND_CONTACT_OK`, `AIGRAM_V02_SMOKE_STOP_POLL_OK`, optional `AIGRAM_V02_SMOKE_SEND_*_SKIPPED`, and final `AIGRAM_V02_SMOKE_OK`. The script must not print bot tokens, token-bearing endpoints, full file IDs, or private message text.

## SendMediaGroup smoke

Use this script to verify a real `SendMediaGroup` flow without requiring external media fixtures:

```bash
export AIGRAM_BOT_TOKEN_MAIN='123456:replace_me'
export AIGRAM_CHAT_ID='123456789'
./scripts/smoke_media_group.sh
```

Default behavior is self-contained: when no media-group file IDs or paths are configured, the smoke sends two generated small text documents as a media group through multipart upload.

Optional inputs:

- `AIGRAM_MEDIA_GROUP_CHAT_ID` overrides `AIGRAM_CHAT_ID`.
- `AIGRAM_MEDIA_GROUP_FILE_ID_1` and `AIGRAM_MEDIA_GROUP_FILE_ID_2` enable FileID mode. These should be file IDs suitable for document media groups; the script must not print them fully.
- `AIGRAM_MEDIA_GROUP_PATH_1` and `AIGRAM_MEDIA_GROUP_PATH_2` enable upload mode from local files. Output should use only safe basenames, not sensitive full paths.

Expected safe markers:

- `AIGRAM_MEDIA_GROUP_SMOKE_WAITING chat_id=... mode=...`
- one of `AIGRAM_MEDIA_GROUP_FILE_ID_OK`, `AIGRAM_MEDIA_GROUP_UPLOAD_OK`, or `AIGRAM_MEDIA_GROUP_GENERATED_UPLOAD_OK`
- final `AIGRAM_MEDIA_GROUP_OK`

The smoke sends real test documents/messages to the configured chat, but it requires no user action and does not run destructive/admin checks.

## Inline keyboard callback checklist

```bash
export AIGRAM_BOT_TOKEN='123456:replace_me'
go run ./examples/inline_longpoll
```

Checklist:

- Send `/start` to the bot.
- Or open `https://t.me/<bot_username>?start=smoke` for examples that support deep-link smoke payloads.
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

## Batch message methods checklist

`ForwardMessages`, `CopyMessages`, and `DeleteMessages` operate on up to 100 message IDs per call. `DeleteMessages` is destructive and is intentionally not part of automatic live smoke.

Manual checklist:

- Use a dedicated test chat.
- Create disposable test messages for the batch operation.
- Run only after explicit confirmation for the target chat.
- Use `ForwardMessages` and `CopyMessages` only with messages intended for testing.
- Use `DeleteMessages` only on disposable test messages created for the check.
- Log only method names, chat ID, message ID count, first returned message ID for copy/forward, and boolean result for delete.
- Do not paste bot tokens, token-bearing URLs, private message text, or full private chat content into logs or reports.

## Webhook SendChatAction and pin/unpin checklist

The default webhook smoke checks `SendChatAction` in the echo handler without requiring extra chat permissions. Pin and unpin methods are not enabled in the default webhook example because they require suitable admin rights in groups/channels and can be noisy.

Checklist for `SendChatAction`:

- Deploy or run `examples/webhook_server`.
- Send a regular text message to the webhook bot.
- The bot should reply with `echo received`.
- Inspect safe logs with `./scripts/remote_logs.sh`; successful chat action is logged as `action=send_chat_action ok=true update_id=... chat_id=... chat_action=typing`, followed by `action=send_message ok=true`.

Manual note for pin/unpin:

- Use a group/channel where the bot has permission to pin messages.
- Call `PinChatMessage`, `UnpinChatMessage`, or `UnpinAllChatMessages` from a small local probe or future targeted example.
- Do not treat pin/unpin failure in a private chat or insufficient-rights chat as a library error; check Telegram permissions first.


## Webhook chat info checklist

The webhook example includes an admin-only chat info action in the access panel. It uses `GetChat` and, when Telegram allows it for the current chat, `GetChatMemberCount`.

Checklist:

- Deploy or run `examples/webhook_server`.
- Open the access panel via the targeted notification or send `/start access_panel`.
- Press `Bot chat info`.
- The bot should send a short safe summary with chat ID, chat type, optional title/username, and member count when available.
- Inspect safe logs with `./scripts/remote_logs.sh`; successful `GetChat` is logged as `action=get_chat ok=true update_id=... chat_id=... chat_type=...`.
- In groups/channels, `GetChatMember`, `GetChatAdministrators`, and member count behavior depends on Telegram access and bot permissions; permission errors are not library errors.

## Moderation methods checklist

`BanChatMember`, `UnbanChatMember`, and `RestrictChatMember` are destructive/admin methods. They are intentionally not part of the default live smoke flow and the smoke/deploy scripts do not call them automatically.

Manual checklist for a dedicated test environment only:

- Create a test group or supergroup.
- Add the bot as an admin with the exact moderation rights you want to test.
- Use only a dedicated test user account as the moderation target.
- Run a small local probe with explicit `ChatID` and `UserID`; do not reuse production chats.
- Start with `RestrictChatMember` and a short `UntilDate` if possible.
- Use `UnbanChatMember` with `OnlyIfBanned: true` to restore the test user after a ban check.
- Do not paste bot tokens, token-bearing URLs, or private group content into logs or reports.

Treat Telegram permission errors in chats where the bot is not an admin as expected Bot API behavior, not as a library bug.

## Chat management methods checklist

`SetChatTitle`, `SetChatDescription`, `SetChatPhoto`, `DeleteChatPhoto`, `LeaveChat`, `SetChatStickerSet`, and `DeleteChatStickerSet` require bot admin rights where applicable and change real chat state. They are intentionally not part of automatic live smoke.

Manual checklist for a dedicated test environment only:

- Create a dedicated test group or supergroup.
- Add the bot as an admin with rights to change chat info; add sticker-set permissions when testing sticker-set methods.
- Do not test these methods on production groups/channels.
- Save the original title, description, photo, and sticker-set state before testing.
- Revert title, description, photo, and sticker-set changes after testing.
- Test `LeaveChat` only with a disposable group or bot instance that can be safely re-added.
- Use `SetChatPhoto` only with explicit test image uploads; do not paste token-bearing URLs or private paths into logs.
- Do not paste bot tokens, token-bearing URLs, private group content, or full file metadata into logs or reports.

Treat Telegram permission errors in chats where the bot is not an admin or lacks change-info rights as expected Bot API behavior, not as a library bug.

## Admin management methods checklist

`PromoteChatMember`, `SetChatAdministratorCustomTitle`, and `SetChatPermissions` require bot admin rights and change real chat or admin state. They are intentionally not part of automatic live smoke.

Manual checklist for a dedicated test environment only:

- Create a dedicated test group or supergroup.
- Use only a dedicated test user account as the admin/permission target.
- Add the bot as an admin with the needed rights, and verify the bot can change those rights before testing.
- Do not test these methods on production groups/channels.
- Check permissions before and after each method call.
- Revert any promoted rights, custom titles, and default chat permission changes after testing.
- Do not paste bot tokens, token-bearing URLs, private group content, or admin lists into logs or reports.

Treat Telegram permission errors in chats where the bot is not an admin or lacks ownership-level rights as expected Bot API behavior, not as a library bug.

## Forum topic methods checklist

`CreateForumTopic`, `EditForumTopic`, `CloseForumTopic`, `ReopenForumTopic`, `DeleteForumTopic`, `UnpinAllForumTopicMessages`, `EditGeneralForumTopic`, `CloseGeneralForumTopic`, `ReopenGeneralForumTopic`, `HideGeneralForumTopic`, `UnhideGeneralForumTopic`, and `UnpinAllGeneralForumTopicMessages` require bot admin rights in a forum supergroup and change real forum topic state. They are intentionally not part of automatic live smoke.

Manual checklist:

- Create a dedicated test forum supergroup.
- Add the bot as an admin with topic-management rights needed for the methods under test.
- Do not run these checks on production groups.
- Create only clearly named test topics.
- Close, reopen, edit, unpin, and delete only test topics.
- If testing General topic methods, record the initial state first and restore it after testing.
- Do not paste bot tokens, token-bearing URLs, private group content, or full production identifiers into logs or reports.

## Reaction methods checklist

`SetMessageReaction` changes real message reaction state. `message_reaction` and `message_reaction_count` updates require allowed update configuration and, for many chats, bot administrator visibility. They are intentionally not part of automatic live smoke.

Manual checklist:

- Use a dedicated test chat and a clearly identified test message.
- Run only after explicit confirmation for the target chat and message.
- Set and clear only test reactions on test messages.
- Do not run on production chats without explicit confirmation.
- If testing updates, include `message_reaction` and `message_reaction_count` in allowed updates where needed.
- Log only method names, chat ID, message ID, reaction type labels, boolean result, and update IDs.
- Do not paste bot tokens, token-bearing URLs, private message text, or full private chat content into logs or reports.

## Invite link methods checklist

`ExportChatInviteLink`, `CreateChatInviteLink`, `EditChatInviteLink`, and `RevokeChatInviteLink` require bot admin rights and create or revoke real chat invite links. They are intentionally not part of automatic live smoke.

Manual checklist for a dedicated test environment only:

- Create a dedicated test group or channel.
- Add the bot as an admin with invite-user rights.
- Use only test invite links and avoid production groups/channels.
- Prefer `CreateChatInviteLink` for a named temporary test link, then `EditChatInviteLink` if needed.
- Revoke every created test link with `RevokeChatInviteLink` after testing.
- Treat Telegram permission errors in non-admin chats as expected Bot API behavior, not as a library bug.
- Do not paste bot tokens, token-bearing URLs, full invite links, or private group content into logs or reports.

## Chat join request methods checklist

`ApproveChatJoinRequest` and `DeclineChatJoinRequest` require bot admin rights with invite-user permission and change real pending join request state. They are intentionally not part of automatic live smoke.

Manual checklist for a dedicated test environment only:

- Create a dedicated test group or channel.
- Add the bot as an admin with invite-user rights.
- Create an invite link with `CreatesJoinRequest: true`.
- Use a dedicated test account to request to join through that link.
- Approve or decline only that test account's join request.
- Do not run this checklist on production groups/channels.
- Do not paste bot tokens, token-bearing URLs, full invite links, user private content, or production group details into logs or reports.

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
