#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
# shellcheck disable=SC1091
source "${SCRIPT_DIR}/env.sh"
trap cleanup_smoke_tunnels EXIT INT TERM

require_env AIGRAM_BOT_TOKEN >/dev/null
require_env AIGRAM_CHAT_ID >/dev/null
prepare_smoke_tunnel

if [ -z "${AIGRAM_MEDIA_PATH:-}" ] && [ -z "${AIGRAM_FILE_ID:-}" ]; then
  echo "Media smoke skipped: set AIGRAM_MEDIA_PATH or AIGRAM_FILE_ID." >&2
  notify_user "Media smoke пропущен: задай AIGRAM_MEDIA_PATH или AIGRAM_FILE_ID."
  exit 0
fi

cd "${REPO_ROOT}"
echo "Starting media upload/download smoke via examples/media_upload."
notify_user "Media smoke запускается. Проверь, что бот прислал/скачал файл."
run_sanitized go run ./examples/media_upload
