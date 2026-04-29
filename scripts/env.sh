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
