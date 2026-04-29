#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
# shellcheck disable=SC1091
source "${SCRIPT_DIR}/env.sh"

require_env AIGRAM_DEPLOY_HOST >/dev/null
require_env AIGRAM_DEPLOY_USER >/dev/null
require_env AIGRAM_DEPLOY_SSH_KEY >/dev/null
require_env AIGRAM_SERVICE_NAME >/dev/null

REMOTE="${AIGRAM_DEPLOY_USER}@${AIGRAM_DEPLOY_HOST}"
SSH_OPTS=(-i "${AIGRAM_DEPLOY_SSH_KEY}" -o IdentitiesOnly=yes -o StrictHostKeyChecking=accept-new)

ssh "${SSH_OPTS[@]}" "${REMOTE}" bash -s -- "${AIGRAM_SERVICE_NAME}" <<'REMOTE_SCRIPT'
set -euo pipefail
service_name="$1"
as_root() {
  if [ "$(id -u)" -eq 0 ]; then
    "$@"
  elif command -v sudo >/dev/null 2>&1; then
    sudo "$@"
  else
    echo "root privileges are required to stop systemd service, but sudo is unavailable" >&2
    return 1
  fi
}

as_root systemctl stop "${service_name}"
echo "Stopped ${service_name}. Files, env, and webhook registration were left unchanged."
REMOTE_SCRIPT
