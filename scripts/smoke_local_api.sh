#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
# shellcheck disable=SC1091
source "${SCRIPT_DIR}/env.sh"
trap cleanup_smoke_tunnels EXIT INT TERM

require_env AIGRAM_BOT_TOKEN >/dev/null
if [ -z "${AIGRAM_BASE_URL:-}" ]; then
  "${SCRIPT_DIR}/discover_env.sh"
  load_generated_env_missing "${GENERATED_ENV_FILE}"
  apply_env_defaults
fi
require_env AIGRAM_BASE_URL >/dev/null
prepare_smoke_tunnel

cd "${REPO_ROOT}"
echo "Starting local Telegram Bot API smoke via examples/local_api_server."
echo "Base URL: ${AIGRAM_BASE_URL}"
set +e
run_sanitized go run ./examples/local_api_server
status=$?
set -e
if [ "${status}" -eq 0 ]; then
  notify_user "Local Bot API smoke успешен."
fi
exit "${status}"
