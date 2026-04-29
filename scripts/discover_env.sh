#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
export AIGRAM_SKIP_GENERATED_ENV=1
# shellcheck disable=SC1091
source "${SCRIPT_DIR}/env.sh"
unset AIGRAM_SKIP_GENERATED_ENV

require_env AIGRAM_BOT_TOKEN >/dev/null
require_env AIGRAM_CHAT_ID >/dev/null
configure_deploy_ssh

SSH_CMD=(ssh "${DEPLOY_SSH_OPTS[@]}" "${DEPLOY_REMOTE}")

listen_port() {
  local listen="$1"
  local port="${listen##*:}"
  if [ "${port}" = "${listen}" ]; then
    port="${listen}"
  fi
  if [[ ! "${port}" =~ ^[0-9]+$ ]]; then
    echo "cannot determine port from AIGRAM_LISTEN_ADDR=${listen}; set AIGRAM_LISTEN_ADDR manually" >&2
    return 1
  fi
  printf '%s\n' "${port}"
}

generate_secret() {
  if command -v openssl >/dev/null 2>&1; then
    openssl rand -hex 24
    return 0
  fi
  if command -v uuidgen >/dev/null 2>&1 && command -v sha256sum >/dev/null 2>&1; then
    printf '%s%s' "$(uuidgen)" "$(date +%s%N)" | sha256sum | awk '{print substr($1,1,48)}'
    return 0
  fi
  echo "cannot generate AIGRAM_WEBHOOK_SECRET automatically: install openssl or set AIGRAM_WEBHOOK_SECRET manually" >&2
  return 1
}

remote_ip() {
  "${SSH_CMD[@]}" 'hostname -I 2>/dev/null' | tr -d '\r' | awk 'NF {print $1; exit}'
}

check_ssh() {
  echo "Checking SSH target ${DEPLOY_REMOTE_LABEL}."
  "${SSH_CMD[@]}" 'true' >/dev/null
}

check_bot_api_candidate() {
  local candidate="$1"
  local remote_script
  remote_script='token="$(cat)"; curl -fsS --max-time 3 "${AIGRAM_CANDIDATE%/}/bot${token}/getMe" | grep -q "\"ok\"[[:space:]]*:[[:space:]]*true"'

  echo "checking Bot API candidate ${candidate}"
  if "${SSH_CMD[@]}" "AIGRAM_CANDIDATE=$(shell_quote "${candidate}") bash -c $(shell_quote "${remote_script}")" <<<"${AIGRAM_BOT_TOKEN}" >/dev/null 2>&1; then
    return 0
  fi
  return 1
}

is_loopback_base_url() {
  case "$1" in
    http://127.0.0.1:*|https://127.0.0.1:*|http://localhost:*|https://localhost:*) return 0 ;;
    *) return 1 ;;
  esac
}

write_generated_line() {
  local name="$1"
  local value="$2"
  printf '%s=%s\n' "${name}" "$(shell_quote "${value}")"
}

check_ssh
REMOTE_IP="$(remote_ip || true)"
LISTEN_PORT="$(listen_port "${AIGRAM_LISTEN_ADDR}")"

BASE_URL="${AIGRAM_BASE_URL:-}"
if [ -z "${BASE_URL}" ]; then
  candidates=(
    "http://127.0.0.1:8081"
    "http://127.0.0.1:8080"
    "http://localhost:8081"
    "http://localhost:8080"
  )
  if [ -n "${REMOTE_IP}" ]; then
    candidates+=("http://${REMOTE_IP}:8081" "http://${REMOTE_IP}:8080")
  fi

  for candidate in "${candidates[@]}"; do
    if check_bot_api_candidate "${candidate}"; then
      BASE_URL="${candidate}"
      break
    fi
  done
fi

if [ -z "${BASE_URL}" ]; then
  echo "Не удалось автоматически найти local Telegram Bot API server. Укажи AIGRAM_BASE_URL вручную, например http://127.0.0.1:8081" >&2
fi

FILE_BASE_URL="${AIGRAM_FILE_BASE_URL:-}"
if [ -z "${FILE_BASE_URL}" ] && [ -n "${BASE_URL}" ]; then
  FILE_BASE_URL="${BASE_URL%/}/file"
fi

WEBHOOK_URL="${AIGRAM_WEBHOOK_URL:-}"
if [ -z "${WEBHOOK_URL}" ]; then
  if [ -n "${BASE_URL}" ] && is_loopback_base_url "${BASE_URL}"; then
    WEBHOOK_URL="http://127.0.0.1:${LISTEN_PORT}/webhook"
  elif [ -n "${BASE_URL}" ] && [ -n "${REMOTE_IP}" ]; then
    WEBHOOK_URL="http://${REMOTE_IP}:${LISTEN_PORT}/webhook"
  else
    echo "AIGRAM_WEBHOOK_URL is required when local Telegram Bot API server was not discovered; use a public https:// URL for official Telegram API" >&2
    exit 1
  fi
elif [ -z "${BASE_URL}" ] && [[ ! "${WEBHOOK_URL}" =~ ^https:// ]]; then
  echo "AIGRAM_WEBHOOK_URL must be https:// when using official Telegram API without AIGRAM_BASE_URL" >&2
  exit 1
fi

WEBHOOK_SECRET="${AIGRAM_WEBHOOK_SECRET:-}"
if [ -z "${WEBHOOK_SECRET}" ]; then
  WEBHOOK_SECRET="$(generate_secret)"
fi

mkdir -p "${REPO_ROOT}/.deploy"
chmod 700 "${REPO_ROOT}/.deploy"
{
  write_generated_line AIGRAM_BASE_URL "${BASE_URL}"
  write_generated_line AIGRAM_FILE_BASE_URL "${FILE_BASE_URL}"
  write_generated_line AIGRAM_WEBHOOK_URL "${WEBHOOK_URL}"
  write_generated_line AIGRAM_WEBHOOK_SECRET "${WEBHOOK_SECRET}"
  write_generated_line AIGRAM_DEPLOY_DIR "${AIGRAM_DEPLOY_DIR}"
  write_generated_line AIGRAM_SERVICE_NAME "${AIGRAM_SERVICE_NAME}"
  write_generated_line AIGRAM_REMOTE_ENV_DIR "${AIGRAM_REMOTE_ENV_DIR}"
  write_generated_line AIGRAM_LISTEN_ADDR "${AIGRAM_LISTEN_ADDR}"
} >"${GENERATED_ENV_FILE}"
chmod 600 "${GENERATED_ENV_FILE}"

cat <<SUMMARY
Generated ${GENERATED_ENV_FILE}
Summary:
- ssh target: ${DEPLOY_REMOTE_LABEL}
- deploy dir: ${AIGRAM_DEPLOY_DIR}
- service name: ${AIGRAM_SERVICE_NAME}
- base url: ${BASE_URL:-<official default>}
- file base url: ${FILE_BASE_URL:-<official default>}
- webhook url: ${WEBHOOK_URL}
- listen addr: ${AIGRAM_LISTEN_ADDR}
- webhook secret: $(mask_secret "${WEBHOOK_SECRET}")
SUMMARY
