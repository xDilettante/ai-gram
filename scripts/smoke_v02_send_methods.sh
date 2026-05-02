#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
# shellcheck disable=SC1091
source "${SCRIPT_DIR}/env.sh"
trap cleanup_smoke_tunnels EXIT INT TERM

export_bot_token_for_role main
prepare_smoke_tunnel

CHAT_ID="${AIGRAM_V02_SMOKE_CHAT_ID:-${AIGRAM_CHAT_ID:-}}"
if [ -z "${CHAT_ID}" ]; then
  echo "AIGRAM_V02_SMOKE_CHAT_ID or AIGRAM_CHAT_ID is required" >&2
  exit 1
fi
export AIGRAM_V02_SMOKE_CHAT_ID="${CHAT_ID}"

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

cd "${REPO_ROOT}"
echo "Starting v0.2 send methods smoke."
echo "Chat ID: $(mask_chat_id "${CHAT_ID}")"
BOT_USERNAME="$(bot_username_for_current_token)"
if [ -n "${BOT_USERNAME}" ]; then
  BOT_LINE="@${BOT_USERNAME}"
else
  BOT_LINE="username unknown"
fi
notify_user "v0.2 send methods smoke is starting.

Bot: ${BOT_LINE}
Chat: $(mask_chat_id "${CHAT_ID}")

Codex will send test contact/location/venue/poll/dice messages.
No user action is required.

If optional media env is set, sticker/animation/video note will be checked." || true

run_sanitized go run ./examples/maintainer/v02_send_methods
