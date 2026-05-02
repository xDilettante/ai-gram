# Echo long polling example

This example runs a small long polling bot that demonstrates:

- `/start` and `/help` command handlers;
- text echo with reply parameters;
- inline keyboard buttons;
- callback query answers;
- editing a message from a callback handler;
- graceful shutdown with `Ctrl+C` or `SIGTERM`.

## Run

```bash
export AIGRAM_BOT_TOKEN="123456:replace-me"
go run ./examples/echo_longpoll
```

Optional variables:

- `AIGRAM_BASE_URL` for a local Telegram Bot API server.
- `AIGRAM_FILE_BASE_URL` for local file download URLs.

The example deletes the current webhook before starting long polling because Telegram does not deliver updates to long polling and webhook receivers at the same time.
