#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
# shellcheck disable=SC1091
source "${SCRIPT_DIR}/env.sh"
trap cleanup_smoke_tunnels EXIT INT TERM

export_bot_token_for_role main
prepare_smoke_tunnel

CHAT_ID="${AIGRAM_MEDIA_GROUP_CHAT_ID:-${AIGRAM_CHAT_ID:-}}"
if [ -z "${CHAT_ID}" ]; then
  echo "AIGRAM_MEDIA_GROUP_CHAT_ID or AIGRAM_CHAT_ID is required" >&2
  exit 1
fi
export AIGRAM_MEDIA_GROUP_CHAT_ID="${CHAT_ID}"

mask_chat_id() {
  local value="${1:-}"
  local length=${#value}
  if [ "${length}" -eq 0 ]; then
    printf 'unknown'
  elif [ "${length}" -le 6 ]; then
    printf '***'
  else
    printf '%s***%s' "${value:0:3}" "${value: -3}"
  fi
}

media_group_mode() {
  if [ -n "${AIGRAM_MEDIA_GROUP_FILE_ID_1:-}" ] || [ -n "${AIGRAM_MEDIA_GROUP_FILE_ID_2:-}" ]; then
    printf 'file_id'
  elif [ -n "${AIGRAM_MEDIA_GROUP_PATH_1:-}" ] || [ -n "${AIGRAM_MEDIA_GROUP_PATH_2:-}" ]; then
    printf 'upload'
  else
    printf 'generated_upload'
  fi
}

cd "${REPO_ROOT}"
echo "Starting SendMediaGroup smoke."
echo "AIGRAM_MEDIA_GROUP_SMOKE_WAITING chat_id=$(mask_chat_id "${CHAT_ID}") mode=$(media_group_mode)"
BOT_USERNAME="$(bot_username_for_current_token)"
if [ -n "${BOT_USERNAME}" ]; then
  BOT_LINE="@${BOT_USERNAME}"
else
  BOT_LINE="username unknown"
fi
notify_user "SendMediaGroup smoke запускается.

Бот: ${BOT_LINE}
Чат: $(mask_chat_id "${CHAT_ID}")

Codex отправит тестовую media group.
От тебя действий не требуется.

Если AIGRAM_MEDIA_GROUP_FILE_ID_* или AIGRAM_MEDIA_GROUP_PATH_* не заданы, будет использован generated upload fallback." || true

TMP_OUTPUT="$(mktemp)"
cleanup_output() {
  rm -f "${TMP_OUTPUT}"
}
trap 'cleanup_output; cleanup_smoke_tunnels' EXIT INT TERM

set +e
run_sanitized go run ./examples/maintainer/media_group_smoke | tee "${TMP_OUTPUT}"
status=${PIPESTATUS[0]}
set -e
if [ "${status}" -ne 0 ]; then
  exit "${status}"
fi
if ! grep -q '^AIGRAM_MEDIA_GROUP_OK ' "${TMP_OUTPUT}"; then
  echo "AIGRAM_MEDIA_GROUP_OK marker not found" >&2
  exit 1
fi
