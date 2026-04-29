#!/usr/bin/env bash
# Shared environment helpers for ai-gram manual smoke scripts.

if [ -z "${BASH_VERSION:-}" ]; then
  echo "scripts/env.sh requires bash" >&2
  return 2 2>/dev/null || exit 2
fi

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "${SCRIPT_DIR}/.." && pwd)"
ENV_FILE="${REPO_ROOT}/.env.local"

if [ -f "${ENV_FILE}" ]; then
  set -a
  # shellcheck disable=SC1090
  source "${ENV_FILE}"
  set +a
fi

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

configure_deploy_ssh() {
  DEPLOY_SSH_OPTS=()

  if [ -n "${AIGRAM_DEPLOY_SSH_TARGET:-}" ]; then
    DEPLOY_REMOTE="${AIGRAM_DEPLOY_SSH_TARGET}"
    DEPLOY_REMOTE_LABEL="${AIGRAM_DEPLOY_SSH_TARGET}"
    return 0
  fi

  require_env AIGRAM_DEPLOY_HOST >/dev/null
  require_env AIGRAM_DEPLOY_USER >/dev/null

  DEPLOY_REMOTE="${AIGRAM_DEPLOY_USER}@${AIGRAM_DEPLOY_HOST}"
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
