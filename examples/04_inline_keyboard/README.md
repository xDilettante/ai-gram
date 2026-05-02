# 04 Inline Keyboard

A long polling bot that sends inline buttons and handles callback queries.

## What it does

- reads `AIGRAM_BOT_TOKEN`;
- sends an inline keyboard on `/start`;
- handles `callback_query` updates;
- calls `AnswerCallbackQuery`;
- edits a message or removes inline buttons.

## Run

```bash
export AIGRAM_BOT_TOKEN='123456:replace_me'
go run ./examples/04_inline_keyboard
```

Send `/start`, then press the inline buttons.
