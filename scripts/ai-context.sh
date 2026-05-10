#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT_DIR"

section() {
  printf '\n==> %s\n' "$1"
}

section "Repository"
printf 'root: %s\n' "$ROOT_DIR"
if git rev-parse --is-inside-work-tree >/dev/null 2>&1; then
  git status --short --branch
else
  printf 'not a git repository\n'
fi

section "Project files"
find . -maxdepth 2 \
  -path './.git' -prune -o \
  -path './build' -prune -o \
  -path './coverage.out' -prune -o \
  -type f \
  \( -name 'AGENTS.md' -o -name 'README*' -o -name 'go.mod' -o -name 'go.sum' -o -name 'Makefile' -o -name 'justfile' -o -name '*.yml' -o -name '*.yaml' \) \
  -print | sort

section "Top-level tree"
find . -maxdepth 2 \
  -path './.git' -prune -o \
  -path './build' -prune -o \
  -print | sort | sed -n '1,160p'

section "Go module"
if [[ -f go.mod ]]; then
  sed -n '1,80p' go.mod
fi

section "Available local checks"
cat <<'EOF'
scripts/check.sh
gofmt -w .
go test ./...
go vet ./...
go build ./...
bash -n scripts/*.sh
git diff --check
EOF

section "Runnable examples"
find examples -maxdepth 2 -name main.go -print | sort | sed 's#^\./##'
