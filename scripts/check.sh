#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT_DIR"

section() {
  printf '\n==> %s\n' "$1"
}

section "Checking Go formatting"
mapfile -d '' go_files < <(find . \
  -path './.git' -prune -o \
  -path './build' -prune -o \
  -path './vendor' -prune -o \
  -name '*.go' -type f -print0)
unformatted="$(gofmt -l "${go_files[@]}")"
if [[ -n "$unformatted" ]]; then
  printf 'Go files need gofmt:\n%s\n' "$unformatted" >&2
  exit 1
fi

section "Checking shell syntax"
bash -n scripts/*.sh

section "Running Go tests"
go test ./...

section "Running go vet"
go vet ./...

section "Building packages"
go build ./...

section "Listing packages"
go list ./...

section "Checking diff whitespace"
git diff --check

printf '\nAll checks passed.\n'
