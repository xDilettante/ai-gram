# 02 Commands

A small long polling bot that routes commands with a plain Go `switch`.

## What it does

- reads `AIGRAM_BOT_TOKEN`;
- handles `/start`, `/help`, and `/about`;
- sends text messages with `SendMessage`.

## Run

```bash
export AIGRAM_BOT_TOKEN='123456:replace_me'
go run ./examples/02_commands
```

Send `/start`, `/help`, or `/about` to the bot.
