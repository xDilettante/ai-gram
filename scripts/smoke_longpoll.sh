#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
# shellcheck disable=SC1091
source "${SCRIPT_DIR}/env.sh"

require_env AIGRAM_BOT_TOKEN >/dev/null

cd "${REPO_ROOT}"
echo "Starting inline long polling smoke. The example calls DeleteWebhook before getUpdates."
echo "Base URL: ${AIGRAM_BASE_URL:-https://api.telegram.org}"
exec go run ./examples/inline_longpoll
