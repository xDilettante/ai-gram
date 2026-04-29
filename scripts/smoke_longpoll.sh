#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
# shellcheck disable=SC1091
source "${SCRIPT_DIR}/env.sh"
trap cleanup_smoke_tunnels EXIT INT TERM

export_bot_token_for_role local
prepare_smoke_tunnel

cd "${REPO_ROOT}"
echo "Starting inline long polling smoke. The example calls DeleteWebhook before getUpdates."
echo "Base URL: ${AIGRAM_BASE_URL:-https://api.telegram.org}"
print_bot_identity
BOT_USERNAME="$(bot_username_for_current_token)"
WAIT_SECONDS="${AIGRAM_SMOKE_WAIT_SECONDS:-120}"
notify_longpoll_smoke_ready "${BOT_USERNAME}" "${WAIT_SECONDS}"
set +e
if command -v timeout >/dev/null 2>&1; then
  timeout --foreground "${WAIT_SECONDS}s" go run ./examples/inline_longpoll 2>&1 | sanitize_stream
  status=${PIPESTATUS[0]}
  if [ "${status}" -eq 124 ]; then
    status=0
  fi
else
  run_sanitized go run ./examples/inline_longpoll
  status=$?
fi
set -e
notify_user "Long polling smoke завершён."
exit "${status}"
