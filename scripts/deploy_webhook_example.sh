#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
# shellcheck disable=SC1091
source "${SCRIPT_DIR}/env.sh"

if [ ! -f "${GENERATED_ENV_FILE}" ]; then
  "${SCRIPT_DIR}/discover_env.sh"
  load_generated_env_missing "${GENERATED_ENV_FILE}"
  apply_env_defaults
fi
export_bot_token_for_role webhook
require_env AIGRAM_WEBHOOK_URL >/dev/null

REMOTE_ENV_DIR="${AIGRAM_REMOTE_ENV_DIR}"
LISTEN_ADDR="${AIGRAM_LISTEN_ADDR}"
BINARY_PATH="${REPO_ROOT}/build/aigram-webhook-server"
TEMPLATE_PATH="${REPO_ROOT}/deploy/systemd/aigram-example.service.tmpl"
configure_deploy_ssh
REMOTE_ENV_FILE="${REMOTE_ENV_DIR}/${AIGRAM_SERVICE_NAME}.env"
REMOTE_EXEC_START="${AIGRAM_DEPLOY_DIR}/aigram-webhook-server"

if [[ ! "${AIGRAM_SERVICE_NAME}" =~ ^[A-Za-z0-9_.@-]+$ ]]; then
  echo "AIGRAM_SERVICE_NAME must contain only letters, digits, underscore, dot, at, or dash" >&2
  exit 1
fi
if [[ "${AIGRAM_DEPLOY_DIR}" != /* ]] || [[ "${REMOTE_ENV_DIR}" != /* ]]; then
  echo "AIGRAM_DEPLOY_DIR and AIGRAM_REMOTE_ENV_DIR must be absolute paths" >&2
  exit 1
fi

SSH_CMD=(ssh "${DEPLOY_SSH_OPTS[@]}" "${DEPLOY_REMOTE}")
SCP_CMD=(scp "${DEPLOY_SSH_OPTS[@]}")

shell_quote() {
  printf '%s' "$1" | sed "s/'/'\\''/g; s/^/'/; s/$/'/"
}

sed_escape() {
  printf '%s' "$1" | sed -e 's/[&|\\]/\\&/g'
}

write_env_line() {
  local name="$1"
  local value="$2"
  if [ -n "${value}" ]; then
    printf '%s=%s\n' "${name}" "$(shell_quote "${value}")"
  fi
}

tmp_dir="$(mktemp -d)"
remote_tmp=""
DEPLOY_SUCCEEDED=0
cleanup() {
  rm -rf "${tmp_dir}"
  if [ -n "${remote_tmp}" ]; then
    "${SSH_CMD[@]}" "rm -rf \"${remote_tmp}\"" >/dev/null 2>&1 || true
  fi
  cleanup_smoke_tunnels
}
on_exit() {
  local status=$?
  cleanup
  if [ "${status}" -ne 0 ] && [ "${DEPLOY_SUCCEEDED}" != "1" ]; then
    notify_user "Webhook deploy упал. Проверь terminal output и remote logs." || true
  fi
  exit "${status}"
}
trap on_exit EXIT

mkdir -p "${REPO_ROOT}/build"
echo "Building linux/amd64 webhook example binary."
(
  cd "${REPO_ROOT}"
  GOOS=linux GOARCH=amd64 go build -o "${BINARY_PATH}" ./examples/webhook_server
)

{
  write_env_line AIGRAM_BOT_TOKEN "${AIGRAM_BOT_TOKEN}"
  write_env_line AIGRAM_BASE_URL "${AIGRAM_BASE_URL:-}"
  write_env_line AIGRAM_FILE_BASE_URL "${AIGRAM_FILE_BASE_URL:-}"
  write_env_line AIGRAM_LISTEN_ADDR "${LISTEN_ADDR}"
  write_env_line AIGRAM_WEBHOOK_URL "${AIGRAM_WEBHOOK_URL}"
  write_env_line AIGRAM_WEBHOOK_SECRET "${AIGRAM_WEBHOOK_SECRET:-}"
} >"${tmp_dir}/service.env"
chmod 600 "${tmp_dir}/service.env"

REMOTE_SERVICE_USER="${AIGRAM_DEPLOY_USER:-$("${SSH_CMD[@]}" 'id -un')}"

sed \
  -e "s|__SERVICE_DESCRIPTION__|$(sed_escape "ai-gram webhook smoke example")|g" \
  -e "s|__ENV_FILE__|$(sed_escape "${REMOTE_ENV_FILE}")|g" \
  -e "s|__EXEC_START__|$(sed_escape "${REMOTE_EXEC_START}")|g" \
  -e "s|__WORKING_DIRECTORY__|$(sed_escape "${AIGRAM_DEPLOY_DIR}")|g" \
  -e "s|__USER__|$(sed_escape "${REMOTE_SERVICE_USER}")|g" \
  "${TEMPLATE_PATH}" >"${tmp_dir}/service.service"

remote_tmp="$("${SSH_CMD[@]}" 'mktemp -d /tmp/aigram-deploy.XXXXXX')"

echo "Uploading webhook binary, systemd unit, and redacted environment file to remote temp directory."
"${SCP_CMD[@]}" "${BINARY_PATH}" "${DEPLOY_REMOTE}:${remote_tmp}/aigram-webhook-server" >/dev/null
"${SCP_CMD[@]}" "${tmp_dir}/service.env" "${DEPLOY_REMOTE}:${remote_tmp}/service.env" >/dev/null
"${SCP_CMD[@]}" "${tmp_dir}/service.service" "${DEPLOY_REMOTE}:${remote_tmp}/service.service" >/dev/null

echo "Installing service ${AIGRAM_SERVICE_NAME} on ${DEPLOY_REMOTE_LABEL}."
set +e
{
"${SSH_CMD[@]}" bash -s -- "${remote_tmp}" "${AIGRAM_DEPLOY_DIR}" "${REMOTE_ENV_DIR}" "${AIGRAM_SERVICE_NAME}" <<'REMOTE_SCRIPT'
set -euo pipefail

tmp="$1"
deploy_dir="$2"
env_dir="$3"
service_name="$4"
env_file="${env_dir}/${service_name}.env"
service_file="/etc/systemd/system/${service_name}.service"

as_root() {
  if [ "$(id -u)" -eq 0 ]; then
    "$@"
  elif command -v sudo >/dev/null 2>&1; then
    sudo "$@"
  else
    echo "root privileges are required for systemd install, but sudo is unavailable" >&2
    return 1
  fi
}

as_root mkdir -p "${deploy_dir}" "${env_dir}"
as_root install -m 0755 "${tmp}/aigram-webhook-server" "${deploy_dir}/aigram-webhook-server"
as_root install -m 0600 "${tmp}/service.env" "${env_file}"
as_root install -m 0644 "${tmp}/service.service" "${service_file}"
as_root systemctl daemon-reload
as_root systemctl enable "${service_name}"

restart_rc=0
as_root systemctl restart "${service_name}" || restart_rc=$?
status_rc=0
as_root systemctl status "${service_name}" --no-pager || status_rc=$?
as_root journalctl -u "${service_name}" -n 80 --no-pager || true
rm -rf "${tmp}"

if [ "${restart_rc}" -ne 0 ]; then
  exit "${restart_rc}"
fi
exit "${status_rc}"
REMOTE_SCRIPT
} 2>&1 | sanitize_stream
remote_status=${PIPESTATUS[0]}
set -e
if [ "${remote_status}" -ne 0 ]; then
  exit "${remote_status}"
fi
remote_tmp=""

echo "Deploy finished. Remote env file: ${REMOTE_ENV_FILE}; webhook secret: $(mask_secret "${AIGRAM_WEBHOOK_SECRET:-}")"
echo "Use ./scripts/remote_logs.sh for logs and ./scripts/remote_stop.sh to stop the service."
DEPLOY_SUCCEEDED=1
notify_user "Webhook example задеплоен на vk1. Отправь /start webhook test bot и затем проверь логи."
