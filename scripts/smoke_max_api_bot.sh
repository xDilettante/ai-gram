#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
# shellcheck disable=SC1091
source "${SCRIPT_DIR}/env.sh"
trap cleanup_smoke_tunnels EXIT INT TERM

MODE="${AIGRAM_MAX_API_MODE:-once}"
DRY_RUN="${AIGRAM_MAX_API_DRY_RUN:-0}"
case "${MODE}" in
  once|poll) ;;
  *)
    echo "AIGRAM_MAX_API_MODE must be once or poll" >&2
    exit 1
    ;;
esac

cd "${REPO_ROOT}"

if [ "${DRY_RUN}" = "1" ]; then
  echo "Dry run: compiling max API smoke bot without Telegram API calls."
  echo "Mode: ${MODE}"
  go test ./examples/maintainer/max_api_bot
  exit 0
fi

export_bot_token_for_role main
prepare_smoke_tunnel

if [ "${MODE}" = "once" ] && [ -z "${AIGRAM_CHAT_ID:-}" ]; then
  echo "AIGRAM_CHAT_ID is required when AIGRAM_MAX_API_MODE=once" >&2
  exit 1
fi

echo "Starting max API smoke bot."
echo "Mode: ${MODE}"
if [ "${MODE}" = "poll" ]; then
  echo "Polling mode may require AIGRAM_MAX_API_DELETE_WEBHOOK=1 when a webhook is currently set."
fi

mkdir -p build/logs
LOG_FILE="${AIGRAM_MAX_API_LOG_FILE:-build/logs/max-api-bot-$(date -u +%Y%m%dT%H%M%SZ).log}"
mkdir -p "$(dirname "${LOG_FILE}")"
echo "Log file: ${LOG_FILE}"
run_sanitized go run ./examples/maintainer/max_api_bot | tee "${LOG_FILE}"
