#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
# shellcheck disable=SC1091
source "${SCRIPT_DIR}/env.sh"

configure_botapi_ssh
BOTAPI_SSH_CMD=(ssh "${BOTAPI_SSH_OPTS[@]}" "${BOTAPI_REMOTE}")
BOTAPI_SCP_CMD=(scp "${BOTAPI_SSH_OPTS[@]}")

BOTAPI_WORKDIR="${AIGRAM_BOTAPI_WORKDIR:-/opt/telegram-bot-api}"
BOTAPI_ENV_FILE="${AIGRAM_BOTAPI_ENV_FILE:-/etc/aigram/telegram-bot-api.env}"
BOTAPI_SERVICE_NAME="${AIGRAM_BOTAPI_SERVICE_NAME}"
BOTAPI_DEST_BINARY="${BOTAPI_WORKDIR%/}/bin/telegram-bot-api"

if [[ ! "${BOTAPI_SERVICE_NAME}" =~ ^[A-Za-z0-9_.@-]+$ ]]; then
  echo "AIGRAM_BOTAPI_SERVICE_NAME must contain only letters, digits, underscore, dot, at, or dash" >&2
  exit 1
fi
if [[ "${BOTAPI_WORKDIR}" != /* ]] || [[ "${BOTAPI_ENV_FILE}" != /* ]]; then
  echo "AIGRAM_BOTAPI_WORKDIR and AIGRAM_BOTAPI_ENV_FILE must be absolute paths" >&2
  exit 1
fi

find_remote_binary() {
  local probe
  probe='set -e
if [ -n "${AIGRAM_BOTAPI_BINARY:-}" ]; then
  [ -f "${AIGRAM_BOTAPI_BINARY}" ] && [ -x "${AIGRAM_BOTAPI_BINARY}" ] && printf "%s\n" "${AIGRAM_BOTAPI_BINARY}"
  exit 0
fi
for candidate in \
  "$HOME/telegram-bot-api" \
  "$HOME/telegram-bot-api/bin/telegram-bot-api" \
  "$HOME/telegram-bot-api/build/telegram-bot-api" \
  "$HOME/telegram-bot-api/telegram-bot-api" \
  "$HOME/bin/telegram-bot-api"; do
  if [ -f "$candidate" ] && [ -x "$candidate" ]; then
    printf "%s\n" "$candidate"
    exit 0
  fi
done
exit 1'
  "${BOTAPI_SSH_CMD[@]}" "AIGRAM_BOTAPI_BINARY=$(shell_quote "${AIGRAM_BOTAPI_BINARY:-}") bash -c $(shell_quote "${probe}")"
}

if ! REMOTE_BINARY="$(find_remote_binary 2>/dev/null)" || [ -z "${REMOTE_BINARY}" ]; then
  echo "telegram-bot-api binary was not found on ${BOTAPI_REMOTE_LABEL}; set AIGRAM_BOTAPI_BINARY" >&2
  exit 1
fi

cat <<PLAN
Bot API service setup plan (no changes unless AIGRAM_CONFIRM_SETUP_BOTAPI=1):
- target: ${BOTAPI_REMOTE_LABEL}
- source binary: ${REMOTE_BINARY}
- installed binary: ${BOTAPI_DEST_BINARY}
- workdir: ${BOTAPI_WORKDIR}
- env file: ${BOTAPI_ENV_FILE}
- service: ${BOTAPI_SERVICE_NAME}
- bind: $(botapi_bind_addr):$(botapi_port)
- mode: --local
- logOut/close: not called
PLAN

if [ "${AIGRAM_CONFIRM_SETUP_BOTAPI:-0}" != "1" ]; then
  echo "Dry run only. Set AIGRAM_CONFIRM_SETUP_BOTAPI=1 to apply this plan."
  exit 0
fi

require_env TELEGRAM_API_ID >/dev/null
require_env TELEGRAM_API_HASH >/dev/null

tmp_dir="$(mktemp -d)"
remote_tmp=""
cleanup() {
  rm -rf "${tmp_dir}"
  if [ -n "${remote_tmp}" ]; then
    "${BOTAPI_SSH_CMD[@]}" "rm -rf \"${remote_tmp}\"" >/dev/null 2>&1 || true
  fi
}
trap cleanup EXIT

cat >"${tmp_dir}/telegram-bot-api.env" <<ENVFILE
TELEGRAM_API_ID=$(shell_quote "${TELEGRAM_API_ID}")
TELEGRAM_API_HASH=$(shell_quote "${TELEGRAM_API_HASH}")
TELEGRAM_BOT_API_PORT=$(shell_quote "$(botapi_port)")
TELEGRAM_BOT_API_BIND_ADDR=$(shell_quote "$(botapi_bind_addr)")
ENVFILE
chmod 600 "${tmp_dir}/telegram-bot-api.env"

cat >"${tmp_dir}/telegram-bot-api.service" <<SERVICE
[Unit]
Description=Local Telegram Bot API server
After=network-online.target
Wants=network-online.target

[Service]
Type=simple
EnvironmentFile=${BOTAPI_ENV_FILE}
ExecStart=${BOTAPI_DEST_BINARY} --api-id=\${TELEGRAM_API_ID} --api-hash=\${TELEGRAM_API_HASH} --http-ip-address=\${TELEGRAM_BOT_API_BIND_ADDR} --http-port=\${TELEGRAM_BOT_API_PORT} --local --dir=${BOTAPI_WORKDIR%/}/data --temp-dir=${BOTAPI_WORKDIR%/}/tmp --log=${BOTAPI_WORKDIR%/}/telegram-bot-api.log --verbosity=2
Restart=on-failure
RestartSec=3
NoNewPrivileges=true
LimitNOFILE=65536

[Install]
WantedBy=multi-user.target
SERVICE

remote_tmp="$("${BOTAPI_SSH_CMD[@]}" 'mktemp -d /tmp/aigram-botapi-setup.XXXXXX')"
"${BOTAPI_SCP_CMD[@]}" "${tmp_dir}/telegram-bot-api.env" "${BOTAPI_REMOTE}:${remote_tmp}/telegram-bot-api.env" >/dev/null
"${BOTAPI_SCP_CMD[@]}" "${tmp_dir}/telegram-bot-api.service" "${BOTAPI_REMOTE}:${remote_tmp}/telegram-bot-api.service" >/dev/null

set +e
{
"${BOTAPI_SSH_CMD[@]}" bash -s -- "${remote_tmp}" "${REMOTE_BINARY}" "${BOTAPI_WORKDIR}" "${BOTAPI_ENV_FILE}" "${BOTAPI_SERVICE_NAME}" "${BOTAPI_DEST_BINARY}" <<'REMOTE_SCRIPT'
set -euo pipefail

tmp="$1"
source_binary="$2"
workdir="$3"
env_file="$4"
service_name="$5"
dest_binary="$6"
service_file="/etc/systemd/system/${service_name}.service"

as_root() {
  if [ "$(id -u)" -eq 0 ]; then
    "$@"
  elif command -v sudo >/dev/null 2>&1; then
    sudo "$@"
  else
    echo "root privileges are required for Bot API service setup, but sudo is unavailable" >&2
    return 1
  fi
}

as_root mkdir -p "${workdir}/bin" "${workdir}/data" "${workdir}/tmp" "$(dirname "${env_file}")"
as_root install -m 0755 "${source_binary}" "${dest_binary}"
as_root install -m 0600 "${tmp}/telegram-bot-api.env" "${env_file}"
as_root install -m 0644 "${tmp}/telegram-bot-api.service" "${service_file}"
as_root systemctl daemon-reload
as_root systemctl enable "${service_name}"
as_root systemctl restart "${service_name}"
as_root systemctl status "${service_name}" --no-pager || true
as_root journalctl -u "${service_name}" -n 80 --no-pager || true
rm -rf "${tmp}"
REMOTE_SCRIPT
} 2>&1 | sanitize_stream
status=${PIPESTATUS[0]}
set -e
exit "${status}"
