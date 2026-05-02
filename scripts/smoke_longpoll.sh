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

print_longpoll_timeout_diagnostics() {
  echo "AIGRAM_SMOKE_TIMEOUT timeout_seconds=${WAIT_SECONDS} username=${BOT_USERNAME:+@${BOT_USERNAME}}" >&2
  if [ -n "${AIGRAM_BASE_URL:-}" ]; then
    if local_http_reachable "${AIGRAM_BASE_URL}"; then
      echo "AIGRAM_SMOKE_DIAGNOSTIC base_url_reachable=true" >&2
    else
      echo "AIGRAM_SMOKE_DIAGNOSTIC base_url_reachable=false" >&2
    fi
  fi
  echo "AIGRAM_SMOKE_DIAGNOSTIC webhook_state=safe_probe_start" >&2
  set +e
  go run ./examples/internal/webhookstatus 2>&1 | sanitize_stream >&2
  local diag_status=${PIPESTATUS[0]}
  set -e
  if [ "${diag_status}" -ne 0 ]; then
    echo "AIGRAM_SMOKE_DIAGNOSTIC webhook_state=failed" >&2
  fi
}

export AIGRAM_SMOKE_EXIT_AFTER_UPDATE=1
echo "AIGRAM_SMOKE_WAITING username=${BOT_USERNAME:+@${BOT_USERNAME}} timeout_seconds=${WAIT_SECONDS}"
set +e
if command -v timeout >/dev/null 2>&1; then
  timeout --foreground "${WAIT_SECONDS}s" go run ./examples/inline_longpoll 2>&1 | sanitize_stream
  status=${PIPESTATUS[0]}
  if [ "${status}" -eq 124 ]; then
    status=1
    set -e
    print_longpoll_timeout_diagnostics
    notify_user "Long polling smoke timeout: no update received within ${WAIT_SECONDS} seconds."
    exit "${status}"
  fi
else
  run_sanitized go run ./examples/inline_longpoll
  status=$?
fi
set -e
if [ "${status}" -eq 0 ]; then
  notify_user "Long polling smoke completed successfully: update received and reply sent."
else
  notify_user "Long polling smoke failed. Check terminal output."
fi
exit "${status}"
