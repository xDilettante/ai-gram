#!/usr/bin/env bash
# Shared environment helpers for ai-gram manual smoke scripts.

if [ -z "${BASH_VERSION:-}" ]; then
  echo "scripts/env.sh requires bash" >&2
  return 2 2>/dev/null || exit 2
fi

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "${SCRIPT_DIR}/.." && pwd)"
ENV_FILE="${REPO_ROOT}/.env.local"
GENERATED_ENV_FILE="${REPO_ROOT}/.deploy/generated.env"

load_env_file() {
  local file="$1"
  if [ -f "${file}" ]; then
    set -a
    # shellcheck disable=SC1090
    source "${file}"
    set +a
  fi
}

load_generated_env_missing() {
  local file="${1:-${GENERATED_ENV_FILE}}"
  local line name raw value

  [ -f "${file}" ] || return 0

  while IFS= read -r line || [ -n "${line}" ]; do
    case "${line}" in
      ''|'#'*) continue ;;
    esac
    name="${line%%=*}"
    raw="${line#*=}"
    if [[ ! "${name}" =~ ^[A-Za-z_][A-Za-z0-9_]*$ ]]; then
      continue
    fi
    if [ -z "${!name:-}" ]; then
      eval "value=${raw}"
      printf -v "${name}" '%s' "${value}"
      export "${name}"
    fi
  done <"${file}"
}

apply_env_defaults() {
  export AIGRAM_DEPLOY_DIR="${AIGRAM_DEPLOY_DIR:-/opt/aigram-test}"
  export AIGRAM_SERVICE_NAME="${AIGRAM_SERVICE_NAME:-aigram-webhook-test}"
  export AIGRAM_REMOTE_ENV_DIR="${AIGRAM_REMOTE_ENV_DIR:-/etc/aigram}"
  export AIGRAM_LISTEN_ADDR="${AIGRAM_LISTEN_ADDR:-:8090}"
}

load_env_file "${ENV_FILE}"
if [ "${AIGRAM_SKIP_GENERATED_ENV:-}" != "1" ]; then
  load_generated_env_missing "${GENERATED_ENV_FILE}"
fi
apply_env_defaults

require_env() {
  local name="$1"
  local value="${!name:-}"
  if [ -z "${value}" ]; then
    echo "required environment variable ${name} is not set; create .env.local from .env.example" >&2
    return 1
  fi
  printf '%s\n' "${value}"
}

optional_env() {
  local name="$1"
  local fallback="${2:-}"
  local value="${!name:-}"
  if [ -z "${value}" ]; then
    printf '%s\n' "${fallback}"
    return 0
  fi
  printf '%s\n' "${value}"
}

mask_secret() {
  local value="${1:-}"
  local length=${#value}
  if [ "${length}" -eq 0 ]; then
    printf '<empty>'
  elif [ "${length}" -le 4 ]; then
    printf '****'
  else
    printf '%s****%s' "${value:0:2}" "${value: -2}"
  fi
}

shell_quote() {
  printf '%s' "$1" | sed "s/'/'\\''/g; s/^/'/; s/$/'/"
}

deploy_ssh_target() {
  if [ -n "${AIGRAM_DEPLOY_SSH_TARGET:-}" ]; then
    printf '%s\n' "${AIGRAM_DEPLOY_SSH_TARGET}"
    return 0
  fi

  require_env AIGRAM_DEPLOY_HOST >/dev/null
  require_env AIGRAM_DEPLOY_USER >/dev/null
  printf '%s@%s\n' "${AIGRAM_DEPLOY_USER}" "${AIGRAM_DEPLOY_HOST}"
}

configure_deploy_ssh() {
  DEPLOY_SSH_OPTS=()
  DEPLOY_REMOTE="$(deploy_ssh_target)"

  if [ -n "${AIGRAM_DEPLOY_SSH_TARGET:-}" ]; then
    DEPLOY_REMOTE_LABEL="${AIGRAM_DEPLOY_SSH_TARGET}"
    return 0
  fi

  DEPLOY_REMOTE_LABEL="${AIGRAM_DEPLOY_HOST}"

  if [ -n "${AIGRAM_DEPLOY_SSH_KEY:-}" ]; then
    if [ ! -f "${AIGRAM_DEPLOY_SSH_KEY}" ]; then
      echo "AIGRAM_DEPLOY_SSH_KEY does not point to a readable private key file" >&2
      return 1
    fi
    DEPLOY_SSH_OPTS=(-i "${AIGRAM_DEPLOY_SSH_KEY}" -o IdentitiesOnly=yes)
  fi

  DEPLOY_SSH_OPTS+=(-o StrictHostKeyChecking=accept-new)
}

sanitize_stream() {
  local token="${AIGRAM_BOT_TOKEN:-}"
  local secret="${AIGRAM_WEBHOOK_SECRET:-}"
  python3 -c '
import os, re, sys
text = sys.stdin.read()
for value, repl in ((os.environ.get("AIGRAM_BOT_TOKEN", ""), "<BOT_TOKEN>"), (os.environ.get("AIGRAM_WEBHOOK_SECRET", ""), "<WEBHOOK_SECRET>")):
    if value:
        text = text.replace(value, repl)
text = re.sub(r"/bot[0-9]+:[A-Za-z0-9_-]+/", "/bot<TOKEN>/", text)
text = re.sub(r"bot[0-9]+:[A-Za-z0-9_-]+", "bot<TOKEN>", text)
sys.stdout.write(text)
'
}

is_loopback_url() {
  case "${1:-}" in
    http://127.0.0.1|http://127.0.0.1:*|http://127.0.0.1/*|http://127.0.0.1:*/*|\
    https://127.0.0.1|https://127.0.0.1:*|https://127.0.0.1/*|https://127.0.0.1:*/*|\
    http://localhost|http://localhost:*|http://localhost/*|http://localhost:*/*|\
    https://localhost|https://localhost:*|https://localhost/*|https://localhost:*/*)
      return 0
      ;;
    *) return 1 ;;
  esac
}

url_scheme() {
  local url="$1"
  printf '%s\n' "${url%%://*}"
}

url_port() {
  local url="$1"
  local rest hostport port scheme
  scheme="${url%%://*}"
  rest="${url#*://}"
  hostport="${rest%%/*}"
  if [[ "${hostport}" == *:* ]]; then
    port="${hostport##*:}"
  elif [ "${scheme}" = "https" ]; then
    port="443"
  else
    port="80"
  fi
  if [[ ! "${port}" =~ ^[0-9]+$ ]]; then
    echo "cannot parse port from URL: ${url}" >&2
    return 1
  fi
  printf '%s\n' "${port}"
}

url_path_suffix() {
  local url="$1"
  local rest="${url#*://}"
  if [[ "${rest}" == */* ]]; then
    printf '/%s\n' "${rest#*/}"
  else
    printf '\n'
  fi
}

local_http_reachable() {
  local url="$1"
  local code
  code=$(curl -sS --max-time 2 -o /dev/null -w '%{http_code}' "${url}" 2>/dev/null || true)
  [ -n "${code}" ] && [ "${code}" != "000" ]
}

find_free_local_port() {
  local preferred="$1"
  local port
  for port in "${preferred}" 18081 18080 18090 18181 18180 18190 19081 19080 19090; do
    if [[ "${port}" =~ ^[0-9]+$ ]] && ! ss -ltn 2>/dev/null | awk '{print $4}' | grep -Eq "(^|:)${port}$"; then
      printf '%s\n' "${port}"
      return 0
    fi
  done
  echo "cannot find a free local port for SSH tunnel" >&2
  return 1
}

cleanup_smoke_tunnels() {
  local pid
  for pid in ${AIGRAM_SMOKE_TUNNEL_PIDS:-}; do
    kill "${pid}" >/dev/null 2>&1 || true
    wait "${pid}" >/dev/null 2>&1 || true
  done
  AIGRAM_SMOKE_TUNNEL_PIDS=""
}

prepare_smoke_tunnel() {
  local base="${AIGRAM_BASE_URL:-}"
  local remote_port local_port local_base suffix preferred wait_attempt

  [ -n "${base}" ] || return 0
  is_loopback_url "${base}" || return 0

  if local_http_reachable "${base}"; then
    echo "Loopback Bot API base URL is reachable locally: ${base}"
    return 0
  fi

  configure_deploy_ssh
  remote_port="$(url_port "${base}")"
  if [ "${remote_port}" -lt 1000 ]; then
    preferred=$((18000 + remote_port))
  else
    preferred=$((10000 + remote_port))
  fi
  local_port="$(find_free_local_port "${preferred}")"
  suffix="$(url_path_suffix "${base}")"
  local_base="$(url_scheme "${base}")://127.0.0.1:${local_port}${suffix}"

  echo "Local ${base} is not reachable; opening temporary SSH tunnel via ${DEPLOY_REMOTE_LABEL}: 127.0.0.1:${local_port} -> 127.0.0.1:${remote_port}."
  ssh "${DEPLOY_SSH_OPTS[@]}" -o ExitOnForwardFailure=yes -N -L "127.0.0.1:${local_port}:127.0.0.1:${remote_port}" "${DEPLOY_REMOTE}" >/tmp/aigram-smoke-tunnel.log 2>&1 &
  local pid=$!
  AIGRAM_SMOKE_TUNNEL_PIDS="${AIGRAM_SMOKE_TUNNEL_PIDS:-} ${pid}"
  export AIGRAM_SMOKE_TUNNEL_PIDS

  for wait_attempt in 1 2 3 4 5; do
    if local_http_reachable "${local_base}"; then
      export AIGRAM_BASE_URL="${local_base}"
      if [ -z "${AIGRAM_FILE_BASE_URL:-}" ] || is_loopback_url "${AIGRAM_FILE_BASE_URL}"; then
        export AIGRAM_FILE_BASE_URL="${local_base%/}/file"
      fi
      echo "Using tunneled Bot API base URL: ${AIGRAM_BASE_URL}"
      return 0
    fi
    sleep 0.4
  done

  echo "SSH tunnel started but ${local_base} is not reachable; see /tmp/aigram-smoke-tunnel.log" >&2
  return 1
}

run_sanitized() {
  local status
  set +e
  "$@" 2>&1 | sanitize_stream
  status=${PIPESTATUS[0]}
  set -e
  return "${status}"
}

notify_enabled() {
  [ "${AIGRAM_NOTIFY_ENABLED:-1}" != "0" ]
}

notify_user() {
  local message="$1"
  local strict="${AIGRAM_NOTIFY_STRICT:-0}"
  local notify_token="${AIGRAM_BOT_TOKEN_NOTIFY:-${AIGRAM_BOT_TOKEN_MAIN:-${AIGRAM_BOT_TOKEN:-}}}"
  local status
  local previous_bot_token_set=0
  local previous_bot_token="${AIGRAM_BOT_TOKEN:-}"
  local previous_notify_text_set=0
  local previous_notify_text="${AIGRAM_NOTIFY_TEXT:-}"

  if ! notify_enabled; then
    return 0
  fi

  if [ -z "${AIGRAM_CHAT_ID:-}" ]; then
    echo "warning: AIGRAM_CHAT_ID is not set; Telegram notification skipped" >&2
    if [ "${strict}" = "1" ]; then
      return 1
    fi
    return 0
  fi

  if [ -z "${notify_token}" ]; then
    echo "warning: notification bot token is not set; Telegram notification skipped" >&2
    if [ "${strict}" = "1" ]; then
      return 1
    fi
    return 0
  fi

  if [ -z "${message}" ]; then
    echo "warning: notification message is empty; Telegram notification skipped" >&2
    if [ "${strict}" = "1" ]; then
      return 1
    fi
    return 0
  fi

  if ! prepare_smoke_tunnel; then
    echo "warning: could not prepare Telegram notification transport; notification skipped" >&2
    if [ "${strict}" = "1" ]; then
      return 1
    fi
    return 0
  fi

  if [ -n "${AIGRAM_BOT_TOKEN+x}" ]; then
    previous_bot_token_set=1
  fi
  if [ -n "${AIGRAM_NOTIFY_TEXT+x}" ]; then
    previous_notify_text_set=1
  fi

  export AIGRAM_BOT_TOKEN="${notify_token}"
  export AIGRAM_NOTIFY_TEXT="${message}"

  set +e
  (
    cd "${REPO_ROOT}"
    run_sanitized go run ./examples/notify_user
  )
  status=$?
  set -e

  if [ "${previous_bot_token_set}" -eq 1 ]; then
    export AIGRAM_BOT_TOKEN="${previous_bot_token}"
  else
    unset AIGRAM_BOT_TOKEN
  fi
  if [ "${previous_notify_text_set}" -eq 1 ]; then
    export AIGRAM_NOTIFY_TEXT="${previous_notify_text}"
  else
    unset AIGRAM_NOTIFY_TEXT
  fi

  if [ "${status}" -ne 0 ]; then
    echo "warning: Telegram notification failed" >&2
    if [ "${strict}" = "1" ]; then
      return "${status}"
    fi
    return 0
  fi

  return 0
}
