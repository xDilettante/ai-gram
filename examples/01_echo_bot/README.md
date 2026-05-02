# 01 Echo Bot

A minimal long polling bot for first-time ai-gram users.

## What it does

- reads `AIGRAM_BOT_TOKEN`;
- deletes any active webhook before starting long polling;
- handles `/start`;
- replies to every text message with the same text.

## Run

```bash
export AIGRAM_BOT_TOKEN='123456:replace_me'
go run ./examples/01_echo_bot
```

Open the bot in Telegram, send `/start`, then send any text message. The bot should send the same text back.

Stop with `Ctrl+C`.
