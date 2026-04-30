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
  export AIGRAM_BOTAPI_PORT="${AIGRAM_BOTAPI_PORT:-8081}"
  export AIGRAM_BOTAPI_BIND_ADDR="${AIGRAM_BOTAPI_BIND_ADDR:-127.0.0.1}"
  export AIGRAM_BOTAPI_SERVICE_NAME="${AIGRAM_BOTAPI_SERVICE_NAME:-telegram-bot-api}"
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

first_non_empty_env() {
  local name value
  for name in "$@"; do
    value="${!name:-}"
    if [ -n "${value}" ]; then
      printf '%s\n' "${value}"
      return 0
    fi
  done
  return 1
}

bot_token_for_role() {
  local role="${1:-}"
  case "${role}" in
    main)
      first_non_empty_env AIGRAM_BOT_TOKEN_MAIN AIGRAM_BOT_TOKEN
      ;;
    cloud)
      first_non_empty_env AIGRAM_BOT_TOKEN_CLOUD AIGRAM_BOT_TOKEN_MAIN AIGRAM_BOT_TOKEN
      ;;
    local)
      first_non_empty_env AIGRAM_BOT_TOKEN_LOCAL AIGRAM_BOT_TOKEN_MAIN AIGRAM_BOT_TOKEN
      ;;
    webhook)
      first_non_empty_env AIGRAM_BOT_TOKEN_WEBHOOK AIGRAM_BOT_TOKEN_MAIN AIGRAM_BOT_TOKEN
      ;;
    notify)
      first_non_empty_env AIGRAM_BOT_TOKEN_NOTIFY AIGRAM_BOT_TOKEN_MAIN AIGRAM_BOT_TOKEN
      ;;
    migration)
      if [ -n "${AIGRAM_BOT_TOKEN_MIGRATION:-}" ]; then
        printf '%s\n' "${AIGRAM_BOT_TOKEN_MIGRATION}"
      elif [ "${AIGRAM_ALLOW_DEFAULT_TOKEN_FOR_MIGRATION:-0}" = "1" ]; then
        first_non_empty_env AIGRAM_BOT_TOKEN_MAIN AIGRAM_BOT_TOKEN
      else
        echo "bot token for role migration is not set; set AIGRAM_BOT_TOKEN_MIGRATION or AIGRAM_ALLOW_DEFAULT_TOKEN_FOR_MIGRATION=1" >&2
        return 1
      fi
      ;;
    destructive)
      if [ -n "${AIGRAM_BOT_TOKEN_DESTRUCTIVE:-}" ]; then
        printf '%s\n' "${AIGRAM_BOT_TOKEN_DESTRUCTIVE}"
      elif [ "${AIGRAM_ALLOW_DEFAULT_TOKEN_FOR_DESTRUCTIVE:-0}" = "1" ]; then
        first_non_empty_env AIGRAM_BOT_TOKEN_MAIN AIGRAM_BOT_TOKEN
      else
        echo "bot token for role destructive is not set; set AIGRAM_BOT_TOKEN_DESTRUCTIVE or AIGRAM_ALLOW_DEFAULT_TOKEN_FOR_DESTRUCTIVE=1" >&2
        return 1
      fi
      ;;
    *)
      echo "unknown bot token role: ${role}" >&2
      return 1
      ;;
  esac
}

export_bot_token_for_role() {
  local role="$1"
  local token
  if ! token="$(bot_token_for_role "${role}")" || [ -z "${token}" ]; then
    echo "bot token for role ${role} is not set" >&2
    return 1
  fi
  export AIGRAM_BOT_TOKEN="${token}"
  echo "Using bot token role ${role}: $(mask_secret "${token}")"
}

shell_quote() {
  printf "'"
  printf '%s' "${1:-}" | sed "s/'/'\\\\''/g"
  printf "'"
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

botapi_ssh_target() {
  if [ -n "${AIGRAM_BOTAPI_SSH_TARGET:-}" ]; then
    printf '%s\n' "${AIGRAM_BOTAPI_SSH_TARGET}"
    return 0
  fi
  deploy_ssh_target
}

botapi_port() {
  printf '%s\n' "${AIGRAM_BOTAPI_PORT:-8081}"
}

botapi_bind_addr() {
  printf '%s\n' "${AIGRAM_BOTAPI_BIND_ADDR:-127.0.0.1}"
}

botapi_base_url_remote() {
  printf 'http://%s:%s\n' "$(botapi_bind_addr)" "$(botapi_port)"
}

botapi_base_url_for_local_client() {
  botapi_base_url_remote
}

configure_botapi_ssh() {
  BOTAPI_SSH_OPTS=()
  BOTAPI_REMOTE="$(botapi_ssh_target)"
  BOTAPI_REMOTE_LABEL="${BOTAPI_REMOTE}"

  if [ -z "${AIGRAM_BOTAPI_SSH_TARGET:-}" ]; then
    configure_deploy_ssh
    BOTAPI_SSH_OPTS=("${DEPLOY_SSH_OPTS[@]}")
    BOTAPI_REMOTE="${DEPLOY_REMOTE}"
    BOTAPI_REMOTE_LABEL="${DEPLOY_REMOTE_LABEL}"
    return 0
  fi

  BOTAPI_SSH_OPTS=(-o StrictHostKeyChecking=accept-new)
}

same_ssh_target() {
  [ "${1:-}" = "${2:-}" ]
}

sanitize_stream() {
  python3 -c '
import os, re, sys
text = sys.stdin.read()
token_names = (
    "AIGRAM_BOT_TOKEN",
    "AIGRAM_BOT_TOKEN_MAIN",
    "AIGRAM_BOT_TOKEN_CLOUD",
    "AIGRAM_BOT_TOKEN_LOCAL",
    "AIGRAM_BOT_TOKEN_WEBHOOK",
    "AIGRAM_BOT_TOKEN_MIGRATION",
    "AIGRAM_BOT_TOKEN_DESTRUCTIVE",
    "AIGRAM_BOT_TOKEN_NOTIFY",
    "TELEGRAM_API_HASH",
)
for name in token_names:
    value = os.environ.get(name, "")
    if value:
        text = text.replace(value, "<BOT_TOKEN>")
for value, repl in ((os.environ.get("AIGRAM_WEBHOOK_SECRET", ""), "<WEBHOOK_SECRET>"),):
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

  configure_botapi_ssh
  remote_port="$(url_port "${base}")"
  if [ "${remote_port}" -lt 1000 ]; then
    preferred=$((18000 + remote_port))
  else
    preferred=$((10000 + remote_port))
  fi
  local_port="$(find_free_local_port "${preferred}")"
  suffix="$(url_path_suffix "${base}")"
  local_base="$(url_scheme "${base}")://127.0.0.1:${local_port}${suffix}"

  echo "Local ${base} is not reachable; opening temporary SSH tunnel: local 127.0.0.1:${local_port} -> ${BOTAPI_REMOTE_LABEL}:127.0.0.1:${remote_port}."
  ssh "${BOTAPI_SSH_OPTS[@]}" -o ExitOnForwardFailure=yes -N -L "127.0.0.1:${local_port}:127.0.0.1:${remote_port}" "${BOTAPI_REMOTE}" >/tmp/aigram-smoke-tunnel.log 2>&1 &
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

print_bot_identity() {
  local status
  set +e
  (
    cd "${REPO_ROOT}"
    run_sanitized go run ./examples/internal/botidentity
  )
  status=$?
  set -e
  if [ "${status}" -ne 0 ]; then
    echo "warning: could not read bot username for current token role" >&2
    return 0
  fi
  return 0
}

bot_username_for_current_token() {
  local output status username tmp
  tmp="$(mktemp)"
  set +e
  (
    cd "${REPO_ROOT}"
    go run ./examples/internal/botidentity
  ) 2>&1 | sanitize_stream >"${tmp}"
  status=${PIPESTATUS[0]}
  set -e
  output="$(cat "${tmp}")"
  rm -f "${tmp}"
  if [ "${status}" -ne 0 ]; then
    printf '\n'
    return 0
  fi
  username="$(printf '%s\n' "${output}" | sed -n 's/^bot username: @//p' | head -n 1)"
  printf '%s\n' "${username}"
}

notify_user() {
  local message="$1"
  local strict="${AIGRAM_NOTIFY_STRICT:-0}"
  local notify_token=""
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

  if ! notify_token="$(bot_token_for_role notify 2>/dev/null)"; then
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

notify_manual_action() {
  local title="$1"
  local body="$2"
  notify_user "${title}

${body}"
}

smoke_mode() {
  local mode="${AIGRAM_SMOKE_MODE:-targeted}"
  case "${mode}" in
    full|targeted|none)
      printf '%s\n' "${mode}"
      ;;
    *)
      echo "warning: unsupported AIGRAM_SMOKE_MODE=${mode}; using targeted" >&2
      printf 'targeted\n'
      ;;
  esac
}

notify_targeted_action() {
  local title="$1"
  local body="$2"
  notify_manual_action "${title}" "${body}"
}

urlencode_payload() {
  python3 -c 'import sys, urllib.parse; print(urllib.parse.quote(sys.argv[1], safe="ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789_-"))' "${1:-}"
}

telegram_bot_start_link() {
  local username="${1:-}"
  local payload="${2:-}"
  username="${username#@}"
  if [ -z "${username}" ]; then
    return 1
  fi
  if [ -z "${payload}" ]; then
    printf 'https://t.me/%s\n' "${username}"
    return 0
  fi
  printf 'https://t.me/%s?start=%s\n' "${username}" "$(urlencode_payload "${payload}")"
}

targeted_smoke() {
  local stage="${AIGRAM_TARGETED_SMOKE:-none}"
  case "${stage}" in
    access|reply|edit|caption|forward_copy|full|none)
      printf '%s\n' "${stage}"
      ;;
    *)
      echo "warning: unsupported AIGRAM_TARGETED_SMOKE=${stage}; using none" >&2
      printf 'none\n'
      ;;
  esac
}

notify_longpoll_smoke_ready() {
  local username="$1"
  local wait_seconds="${2:-120}"
  local bot_line="Бот: username unknown"
  local open_line="Открой бота вручную."
  if [ -n "${username}" ]; then
    bot_line="Бот: @${username}"
    open_line="Открой: $(telegram_bot_start_link "${username}" "smoke")"
  fi

  notify_manual_action "Long polling smoke запущен." "${bot_line}
${open_line}

Сделай:
1. Отправь одно обычное текстовое сообщение.

Smoke завершится автоматически сразу после получения update и отправки reply.
Timeout: ${wait_seconds} секунд.

Codex проверит safe output сам."
}

webhook_bot_lines() {
  local username="$1"
  local payload="${2:-smoke}"
  local bot_line="Бот: username unknown"
  local open_line="Открой webhook bot вручную."
  if [ -n "${username}" ]; then
    bot_line="Бот: @${username}"
    open_line="Открой: $(telegram_bot_start_link "${username}" "${payload}")"
  fi
  printf '%s\n%s\n' "${bot_line}" "${open_line}"
}

notify_full_webhook_smoke() {
  local username="$1"
  local lines
  lines="$(webhook_bot_lines "${username}" "smoke")"

  notify_manual_action "Webhook full smoke готов." "${lines}

Нажми ссылку: бот откроет smoke keyboard.
Дальше пройди полный regression flow кнопками в боте.

Codex проверит safe action logs сам."
}

notify_webhook_deploy_done() {
  local username="$1"
  local lines
  lines="$(webhook_bot_lines "${username}" "smoke")"

  notify_manual_action "Webhook example задеплоен." "${lines}

Действий от пользователя не требуется, если текущий этап отдельно не запросил targeted smoke.

Если Codex прислал targeted notification, выполняй только указанные там шаги, не полный checklist."
}

notify_webhook_smoke_ready() {
  local username="$1"
  case "$(smoke_mode)" in
    full)
      notify_full_webhook_smoke "${username}"
      ;;
    targeted)
      case "$(targeted_smoke)" in
        access)
          notify_access_panel_ready "${username}"
          ;;
        reply)
          notify_reply_smoke_ready "${username}"
          ;;
        edit)
          notify_edit_smoke_ready "${username}"
          ;;
        caption)
          notify_caption_smoke_ready "${username}"
          ;;
        forward_copy)
          notify_forward_copy_smoke_ready "${username}"
          ;;
        full)
          notify_full_webhook_smoke "${username}"
          ;;
        none)
          notify_webhook_deploy_done "${username}"
          ;;
      esac
      ;;
    none)
      return 0
      ;;
  esac
}

notify_access_panel_ready() {
  local username="$1"
  local lines
  lines="$(webhook_bot_lines "${username}" "access_panel")"

  notify_targeted_action "Access control smoke готов." "${lines}

Нажми ссылку: бот откроет control panel.
Дальше просто нажимай кнопки в боте:
- Access status
- Bot chat info
- Start smoke
- Open access
- Close access

Codex проверит safe logs сам."
}

notify_reply_smoke_ready() {
  local username="$1"
  local lines
  lines="$(webhook_bot_lines "${username}" "smoke")"

  notify_targeted_action "Reply smoke готов." "${lines}

Нажми ссылку и отправь обычное текстовое сообщение.

Codex проверит safe logs: matched=message и action=send_message с reply_to_message_id."
}

notify_edit_smoke_ready() {
  local username="$1"
  local lines
  lines="$(webhook_bot_lines "${username}" "smoke")"

  notify_targeted_action "Edit smoke готов." "${lines}

Нажми ссылку: бот откроет smoke keyboard.
Дальше нажми:
- Edit message
- Remove keyboard

Codex проверит action=edit_message_text и action=edit_message_reply_markup."
}

notify_caption_smoke_ready() {
  local username="$1"
  local lines
  lines="$(webhook_bot_lines "${username}" "smoke")"

  notify_targeted_action "Caption smoke готов." "${lines}

Нажми ссылку: бот откроет smoke keyboard.
Дальше нажми:
- Caption demo
- Edit caption
- Delete media message

Codex проверит action=send_media_caption_demo, action=edit_message_caption и action=delete_message."
}

notify_delete_smoke_ready() {
  local username="$1"
  local lines
  lines="$(webhook_bot_lines "${username}" "smoke")"

  notify_targeted_action "Delete smoke готов." "${lines}

Нажми ссылку: бот откроет smoke keyboard.
Дальше нажми "Delete this message".

Codex проверит action=delete_message."
}

notify_forward_copy_smoke_ready() {
  local username="$1"
  local lines
  lines="$(webhook_bot_lines "${username}" "smoke")"

  notify_targeted_action "Forward/copy smoke готов." "${lines}

Нажми ссылку: бот откроет smoke keyboard.
Дальше нажми:
- Copy this message
- Forward this message

Codex проверит action=copy_message и action=forward_message."
}

notify_media_smoke_ready() {
  local username="$1"
  local bot_line="Бот: username unknown"
  local open_line="Проверь чат с media smoke bot вручную."
  if [ -n "${username}" ]; then
    bot_line="Бот: @${username}"
    open_line="Открой: https://t.me/${username}"
  fi

  notify_manual_action "Media smoke запускается." "${bot_line}
${open_line}

Проверь, что бот прислал или скачал файл.
Codex проверит output/logs без token-bearing URLs."
}
