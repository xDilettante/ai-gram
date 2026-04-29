#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
export AIGRAM_SKIP_GENERATED_ENV=1
# shellcheck disable=SC1091
source "${SCRIPT_DIR}/env.sh"
unset AIGRAM_SKIP_GENERATED_ENV

require_env AIGRAM_CHAT_ID >/dev/null
configure_deploy_ssh
configure_botapi_ssh
export_bot_token_for_role local

DEPLOY_SSH_CMD=(ssh "${DEPLOY_SSH_OPTS[@]}" "${DEPLOY_REMOTE}")
BOTAPI_SSH_CMD=(ssh "${BOTAPI_SSH_OPTS[@]}" "${BOTAPI_REMOTE}")

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
  "${DEPLOY_SSH_CMD[@]}" 'hostname -I 2>/dev/null' | tr -d '\r' | awk 'NF {print $1; exit}'
}

check_ssh() {
  echo "Checking deploy SSH target ${DEPLOY_REMOTE_LABEL}."
  "${DEPLOY_SSH_CMD[@]}" 'true' >/dev/null
  echo "Checking Bot API SSH target ${BOTAPI_REMOTE_LABEL}."
  "${BOTAPI_SSH_CMD[@]}" 'true' >/dev/null
}

check_botapi_outbound() {
  echo "checking Bot API host outbound HTTPS to api.telegram.org from ${BOTAPI_REMOTE_LABEL}"
  if "${BOTAPI_SSH_CMD[@]}" 'curl -4 -I --max-time 10 https://api.telegram.org >/dev/null 2>&1'; then
    echo "Bot API host can reach api.telegram.org over IPv4 HTTPS."
    return 0
  fi
  echo "warning: Bot API host cannot reach api.telegram.org over IPv4 HTTPS" >&2
  return 0
}

check_bot_api_candidate() {
  local candidate="$1"
  local remote_script
  remote_script='token="$(cat)"; curl -fsS --max-time 5 "${AIGRAM_CANDIDATE%/}/bot${token}/getMe" | grep -q "\"ok\"[[:space:]]*:[[:space:]]*true"'

  echo "checking Bot API candidate ${candidate} on ${BOTAPI_REMOTE_LABEL}"
  if "${BOTAPI_SSH_CMD[@]}" "AIGRAM_CANDIDATE=$(shell_quote "${candidate}") bash -c $(shell_quote "${remote_script}")" <<<"${AIGRAM_BOT_TOKEN}" >/dev/null 2>&1; then
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
check_botapi_outbound
REMOTE_IP="$(remote_ip || true)"
LISTEN_PORT="$(listen_port "${AIGRAM_LISTEN_ADDR}")"

BASE_URL="${AIGRAM_BASE_URL:-}"
if [ -z "${BASE_URL}" ]; then
  candidates=(
    "$(botapi_base_url_remote)"
    "http://127.0.0.1:8081"
    "http://127.0.0.1:8080"
    "http://localhost:8081"
    "http://localhost:8080"
  )

  for candidate in "${candidates[@]}"; do
    if check_bot_api_candidate "${candidate}"; then
      BASE_URL="${candidate}"
      break
    fi
  done
fi

if [ -z "${BASE_URL}" ]; then
  echo "Укажи AIGRAM_BOTAPI_SSH_TARGET или запусти telegram-bot-api на выбранном сервере." >&2
fi

FILE_BASE_URL="${AIGRAM_FILE_BASE_URL:-}"
if [ -z "${FILE_BASE_URL}" ] && [ -n "${BASE_URL}" ]; then
  FILE_BASE_URL="${BASE_URL%/}/file"
fi

WEBHOOK_URL="${AIGRAM_WEBHOOK_URL:-}"
WEBHOOK_URL_NOTE=""
if [ -z "${WEBHOOK_URL}" ]; then
  if [ -n "${BASE_URL}" ] && is_loopback_base_url "${BASE_URL}"; then
    if same_ssh_target "${BOTAPI_REMOTE}" "${DEPLOY_REMOTE}"; then
      WEBHOOK_URL="http://127.0.0.1:${LISTEN_PORT}/webhook"
    else
      WEBHOOK_URL_NOTE="manual required: Bot API server and webhook service are on different SSH targets"
      echo "warning: Bot API server и webhook service на разных серверах; AIGRAM_WEBHOOK_URL нужно задать вручную для deploy." >&2
    fi
  elif [ -n "${BASE_URL}" ] && [ -n "${REMOTE_IP}" ]; then
    if same_ssh_target "${BOTAPI_REMOTE}" "${DEPLOY_REMOTE}"; then
      WEBHOOK_URL="http://${REMOTE_IP}:${LISTEN_PORT}/webhook"
    else
      WEBHOOK_URL_NOTE="manual required: Bot API server and webhook service are on different SSH targets"
      echo "warning: Bot API server и webhook service на разных серверах; AIGRAM_WEBHOOK_URL нужно задать вручную для deploy." >&2
    fi
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
  write_generated_line AIGRAM_BOTAPI_SSH_TARGET "${BOTAPI_REMOTE}"
  write_generated_line AIGRAM_BASE_URL "${BASE_URL}"
  write_generated_line AIGRAM_FILE_BASE_URL "${FILE_BASE_URL}"
  write_generated_line AIGRAM_WEBHOOK_URL "${WEBHOOK_URL}"
  write_generated_line AIGRAM_WEBHOOK_SECRET "${WEBHOOK_SECRET}"
  write_generated_line AIGRAM_DEPLOY_DIR "${AIGRAM_DEPLOY_DIR}"
  write_generated_line AIGRAM_SERVICE_NAME "${AIGRAM_SERVICE_NAME}"
  write_generated_line AIGRAM_REMOTE_ENV_DIR "${AIGRAM_REMOTE_ENV_DIR}"
  write_generated_line AIGRAM_LISTEN_ADDR "${AIGRAM_LISTEN_ADDR}"
  write_generated_line AIGRAM_BOTAPI_PORT "$(botapi_port)"
  write_generated_line AIGRAM_BOTAPI_BIND_ADDR "$(botapi_bind_addr)"
} >"${GENERATED_ENV_FILE}"
chmod 600 "${GENERATED_ENV_FILE}"

cat <<SUMMARY
Generated ${GENERATED_ENV_FILE}
Summary:
- deploy ssh target: ${DEPLOY_REMOTE_LABEL}
- bot api ssh target: ${BOTAPI_REMOTE_LABEL}
- deploy dir: ${AIGRAM_DEPLOY_DIR}
- service name: ${AIGRAM_SERVICE_NAME}
- bot api remote base url: $(botapi_base_url_remote)
- base url: ${BASE_URL:-<official default>}
- file base url: ${FILE_BASE_URL:-<official default>}
- webhook url: ${WEBHOOK_URL:-<manual required>}
- listen addr: ${AIGRAM_LISTEN_ADDR}
- webhook note: ${WEBHOOK_URL_NOTE:-<none>}
- webhook secret: $(mask_secret "${WEBHOOK_SECRET}")
SUMMARY
