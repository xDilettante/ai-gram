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
