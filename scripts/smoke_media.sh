#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
# shellcheck disable=SC1091
source "${SCRIPT_DIR}/env.sh"

require_env AIGRAM_BOT_TOKEN >/dev/null
require_env AIGRAM_CHAT_ID >/dev/null

if [ -z "${AIGRAM_MEDIA_PATH:-}" ] && [ -z "${AIGRAM_FILE_ID:-}" ]; then
  echo "set at least one of AIGRAM_MEDIA_PATH or AIGRAM_FILE_ID for media smoke" >&2
  exit 1
fi

cd "${REPO_ROOT}"
echo "Starting media upload/download smoke via examples/media_upload."
exec go run ./examples/media_upload
