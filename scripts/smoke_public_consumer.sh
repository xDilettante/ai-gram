#!/usr/bin/env bash
set -euo pipefail

VERSION="${AIGRAM_CONSUMER_VERSION:-main}"
DIRECT="${AIGRAM_CONSUMER_DIRECT:-0}"
KEEP_TEMP="${AIGRAM_CONSUMER_KEEP_TEMP:-0}"

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "${SCRIPT_DIR}/.." && pwd)"
MODULE="github.com/xDilettante/ai-gram"
WORK_DIR="$(mktemp -d -t aigram-public-consumer-XXXXXX)"

log() {
  printf '%s\n' "$*" >&2
}

hint_direct() {
  if [ "${DIRECT}" = "0" ]; then
    log "If the commit was pushed recently or the public Go proxy is stale/unavailable, retry with AIGRAM_CONSUMER_DIRECT=1."
  fi
}

cleanup() {
  if [ "${KEEP_TEMP}" = "1" ]; then
    log "Keeping temporary consumer module: ${WORK_DIR}"
    return
  fi
  rm -rf "${WORK_DIR}"
}
trap cleanup EXIT

case "${VERSION}" in
  ''|*[[:space:]]*)
    log "AIGRAM_CONSUMER_VERSION must be a non-empty Go module version without whitespace"
    exit 2
    ;;
esac

if [ "${DIRECT}" = "1" ]; then
  export GOPROXY=direct
elif [ "${DIRECT}" != "0" ]; then
  log "AIGRAM_CONSUMER_DIRECT must be 0 or 1"
  exit 2
fi

log "Repository: ${REPO_ROOT}"
log "Temporary consumer module: ${WORK_DIR}"
log "Version: ${VERSION}"
if [ "${DIRECT}" = "1" ]; then
  log "GOPROXY: direct"
fi

cd "${WORK_DIR}"

log "Initializing temporary Go module."
go mod init example.com/aigram-public-consumer >/dev/null

log "Fetching ${MODULE}@${VERSION}."
if ! go get "${MODULE}@${VERSION}"; then
  hint_direct
  exit 1
fi
log "Resolved module version:"
go list -m "${MODULE}" >&2

cat >main_test.go <<'GOEOF'
package consumer

import (
	"testing"

	aigram "github.com/xDilettante/ai-gram"
	botpkg "github.com/xDilettante/ai-gram/bot"
	"github.com/xDilettante/ai-gram/telegram"
)

func TestRootQuickStartAPICompiles(t *testing.T) {
	client, err := aigram.New(aigram.Config{Token: "123456:TEST_TOKEN"})
	if err != nil {
		t.Fatalf("New returned error: %v", err)
	}
	if client == nil {
		t.Fatal("expected bot client")
	}

	_ = aigram.SendMessageParams{
		ChatID: aigram.ChatIDInt(12345),
		Text:   "hello",
	}
}

func TestAdvancedImportsCompile(t *testing.T) {
	client, err := botpkg.New(botpkg.Config{Token: "123456:TEST_TOKEN"})
	if err != nil {
		t.Fatalf("New returned error: %v", err)
	}
	if client == nil {
		t.Fatal("expected bot client")
	}

	loginButton := telegram.InlineKeyboardButton{
		Text:     "Login",
		LoginURL: &telegram.LoginURL{URL: "https://example.com/login"},
	}
	markup := telegram.NewInlineKeyboard([]telegram.InlineKeyboardButton{loginButton})
	if err := telegram.ValidateReplyMarkup(markup); err != nil {
		t.Fatalf("reply markup should be valid: %v", err)
	}

	_ = botpkg.SendMessageParams{
		ChatID:      botpkg.ChatIDInt(12345),
		Text:        "hello",
		ReplyMarkup: markup,
	}
}
GOEOF

log "Running external consumer tests."
if ! go test ./...; then
  hint_direct
  exit 1
fi

log "Checking public Go documentation."
if ! go doc "${MODULE}" >/dev/null; then
  hint_direct
  exit 1
fi
if ! go doc "${MODULE}/bot.Config" >/dev/null; then
  hint_direct
  exit 1
fi
if ! go doc "${MODULE}/telegram.LoginURL" >/dev/null; then
  hint_direct
  exit 1
fi

printf 'Public consumer smoke passed for %s@%s\n' "${MODULE}" "${VERSION}"
