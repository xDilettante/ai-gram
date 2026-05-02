# 03 Reply Keyboard

A long polling bot that shows a regular Telegram reply keyboard.

## What it does

- reads `AIGRAM_BOT_TOKEN`;
- sends a keyboard with `Help`, `About`, and `Remove keyboard` buttons;
- removes the keyboard with `ReplyKeyboardRemove`.

## Run

```bash
export AIGRAM_BOT_TOKEN='123456:replace_me'
go run ./examples/03_reply_keyboard
```

Send `/start` to show the keyboard, then press the buttons.
