# Maintainer Smoke Environment Template

This template is for ai-gram maintainers who run the private live-smoke and deploy harness.
It is intentionally not part of the public `.env.example` because normal library users do not need
multiple bot roles, deploy hosts, notification bots, or local Bot API service controls.

Copy only the variables you need into `.env.local` or a private environment file. Never commit real
bot tokens, webhook secrets, payment payloads, invite links, private chat IDs, or SSH details.

```dotenv
# Common smoke targets
AIGRAM_CHAT_ID=
AIGRAM_DEPLOY_SSH_TARGET=

# Bot tokens
# AIGRAM_BOT_TOKEN is a legacy fallback. Prefer role-specific tokens for integration checks.
AIGRAM_BOT_TOKEN=
AIGRAM_BOT_TOKEN_MAIN=
AIGRAM_BOT_TOKEN_CLOUD=
AIGRAM_BOT_TOKEN_LOCAL=
AIGRAM_BOT_TOKEN_WEBHOOK=

# Keep migration/destructive checks on separate bots because they can trigger cooldowns or drop updates.
AIGRAM_BOT_TOKEN_MIGRATION=
AIGRAM_BOT_TOKEN_DESTRUCTIVE=
AIGRAM_BOT_TOKEN_NOTIFY=
AIGRAM_ALLOW_DEFAULT_TOKEN_FOR_MIGRATION=0
AIGRAM_ALLOW_DEFAULT_TOKEN_FOR_DESTRUCTIVE=0

# Optional smoke controls
AIGRAM_NOTIFY_ENABLED=1
AIGRAM_NOTIFY_STRICT=0
AIGRAM_SMOKE_WAIT_SECONDS=120
AIGRAM_SMOKE_MODE=targeted
AIGRAM_TARGETED_SMOKE=none

# Access control for examples. Default is admin-only.
# AIGRAM_ADMIN_USER_IDS, AIGRAM_ALLOWED_USER_IDS, and AIGRAM_ALLOWED_CHAT_IDS are comma-separated int64 IDs.
# If AIGRAM_ADMIN_USER_IDS is empty, examples fall back to AIGRAM_CHAT_ID as admin user/chat for private smoke.
AIGRAM_ACCESS_MODE=admin
AIGRAM_ADMIN_USER_IDS=
AIGRAM_ALLOWED_USER_IDS=
AIGRAM_ALLOWED_CHAT_IDS=

# Bot API endpoint and webhook settings
AIGRAM_BASE_URL=
AIGRAM_FILE_BASE_URL=
AIGRAM_WEBHOOK_URL=
AIGRAM_LISTEN_ADDR=
AIGRAM_WEBHOOK_SECRET=
AIGRAM_MEDIA_PATH=
AIGRAM_FILE_ID=

# Optional v0.2 send methods smoke inputs
AIGRAM_V02_SMOKE_CHAT_ID=
AIGRAM_STICKER_FILE_ID=
AIGRAM_ANIMATION_PATH=
AIGRAM_ANIMATION_FILE_ID=
AIGRAM_VIDEO_NOTE_PATH=
AIGRAM_VIDEO_NOTE_FILE_ID=

# Optional SendMediaGroup smoke inputs
AIGRAM_MEDIA_GROUP_CHAT_ID=
AIGRAM_MEDIA_GROUP_FILE_ID_1=
AIGRAM_MEDIA_GROUP_FILE_ID_2=
AIGRAM_MEDIA_GROUP_PATH_1=
AIGRAM_MEDIA_GROUP_PATH_2=

# Deploy harness settings
AIGRAM_DEPLOY_DIR=
AIGRAM_SERVICE_NAME=
AIGRAM_REMOTE_ENV_DIR=

# Optional separate local Bot API server host
AIGRAM_BOTAPI_SSH_TARGET=
AIGRAM_BOTAPI_PORT=8081
AIGRAM_BOTAPI_BIND_ADDR=127.0.0.1
AIGRAM_BOTAPI_BINARY=
AIGRAM_BOTAPI_WORKDIR=
AIGRAM_BOTAPI_SERVICE_NAME=telegram-bot-api
AIGRAM_BOTAPI_ENV_FILE=

# Required only if setting up or restarting a telegram-bot-api service.
TELEGRAM_API_ID=
TELEGRAM_API_HASH=

# Legacy SSH fallback when AIGRAM_DEPLOY_SSH_TARGET is not used.
AIGRAM_DEPLOY_HOST=
AIGRAM_DEPLOY_USER=
AIGRAM_DEPLOY_SSH_KEY=
```

Use `docs/maintainer/DEPLOY_TESTING.md` and `docs/maintainer/LIVE_SMOKE_MATRIX.md` for maintainer-only procedures.
