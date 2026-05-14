# 07 Inline Panel

A production-style long polling example with typed inline callback data.

## What it does

- reads `AIGRAM_BOT_TOKEN`;
- deletes any active webhook before starting long polling;
- uses `transport/longpoll` for update intake;
- uses `dispatch` for command and callback routing;
- uses `callback` for typed `callback_data`;
- shows a paginated inline panel;
- shows confirm/cancel actions without running real side effects.

## Run

```bash
export AIGRAM_BOT_TOKEN='123456:replace_me'
go run ./examples/07_inline_panel
```

Open the bot in Telegram and send `/panel` or `/start`.

The example only edits messages. Put real production side effects behind your own authorization, audit logging, idempotency, and retry policy.
