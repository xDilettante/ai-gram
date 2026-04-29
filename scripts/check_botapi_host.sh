#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
# shellcheck disable=SC1091
source "${SCRIPT_DIR}/env.sh"

configure_botapi_ssh
BOTAPI_SSH_CMD=(ssh "${BOTAPI_SSH_OPTS[@]}" "${BOTAPI_REMOTE}")

remote_binary_probe='set -e
if [ -n "${AIGRAM_BOTAPI_BINARY:-}" ]; then
  if [ -x "${AIGRAM_BOTAPI_BINARY}" ]; then
    printf "%s\n" "${AIGRAM_BOTAPI_BINARY}"
    exit 0
  fi
  exit 1
fi
for candidate in \
  "$HOME/telegram-bot-api" \
  "$HOME/telegram-bot-api/bin/telegram-bot-api" \
  "$HOME/telegram-bot-api/build/telegram-bot-api" \
  "$HOME/telegram-bot-api/telegram-bot-api" \
  "$HOME/bin/telegram-bot-api"; do
  if [ -x "$candidate" ]; then
    printf "%s\n" "$candidate"
    exit 0
  fi
done
exit 1'

echo "Checking Bot API host: ${BOTAPI_REMOTE_LABEL}"
echo "Configured remote base URL: $(botapi_base_url_remote)"

"${BOTAPI_SSH_CMD[@]}" 'hostname && whoami && uname -a'

echo "--- DNS api.telegram.org ---"
"${BOTAPI_SSH_CMD[@]}" 'getent hosts api.telegram.org || true; getent ahosts api.telegram.org || true'

echo "--- HTTPS api.telegram.org ---"
"${BOTAPI_SSH_CMD[@]}" 'curl -4 -I --max-time 10 https://api.telegram.org || true; echo ---; curl -6 -I --max-time 10 https://api.telegram.org || true' 2>&1 | sanitize_stream

echo "--- local listeners ---"
"${BOTAPI_SSH_CMD[@]}" 'ss -tulpn | grep -E ":(8080|8081)\b" || true'

echo "--- local Bot API root probes ---"
"${BOTAPI_SSH_CMD[@]}" 'curl -fsS --max-time 5 http://127.0.0.1:8081/ || true; echo; curl -fsS --max-time 5 http://127.0.0.1:8080/ || true; echo' 2>&1 | sanitize_stream

echo "--- telegram-bot-api binary ---"
if binary_path=$("${BOTAPI_SSH_CMD[@]}" "AIGRAM_BOTAPI_BINARY=$(shell_quote "${AIGRAM_BOTAPI_BINARY:-}") bash -c $(shell_quote "${remote_binary_probe}")" 2>/dev/null); then
  echo "binary found: ${binary_path}"
else
  echo "binary not found in default locations; set AIGRAM_BOTAPI_BINARY"
fi

echo "--- service status ---"
"${BOTAPI_SSH_CMD[@]}" "systemctl status $(shell_quote "${AIGRAM_BOTAPI_SERVICE_NAME}") --no-pager || true" 2>&1 | sanitize_stream

cat <<SUMMARY
Summary:
- bot api ssh target: ${BOTAPI_REMOTE_LABEL}
- bind addr: $(botapi_bind_addr)
- port: $(botapi_port)
- remote base url: $(botapi_base_url_remote)
- setup service script is not run by this check
SUMMARY
