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
notify_user "Long polling smoke запускается. Отправь local test bot любое сообщение или /start."
set +e
run_sanitized go run ./examples/inline_longpoll
status=$?
set -e
notify_user "Long polling smoke завершён."
exit "${status}"
