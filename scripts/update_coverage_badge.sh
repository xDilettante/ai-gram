#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT_DIR"

coverage_file="${1:-coverage.out}"
badge_file="docs/assets/coverage.svg"

go test -coverprofile="$coverage_file" ./... >/tmp/aigram-coverage-test.log
coverage="$(go tool cover -func="$coverage_file" | awk '/^total:/ {print $3}')"
if [[ -z "$coverage" ]]; then
  echo "could not parse total coverage" >&2
  cat /tmp/aigram-coverage-test.log >&2
  exit 1
fi

coverage_value="${coverage%%%}"
color="#ef4444"
awk -v value="$coverage_value" 'BEGIN { exit !(value >= 70) }' && color="#22c55e" || true
awk -v value="$coverage_value" 'BEGIN { exit !(value >= 50 && value < 70) }' && color="#f59e0b" || true

python3 - "$coverage" "$color" "$badge_file" <<'PY'
import html
import sys
coverage, color, path = sys.argv[1:4]
label = "coverage"
value = coverage
svg = f'''<svg xmlns="http://www.w3.org/2000/svg" width="180" height="28" viewBox="0 0 180 28" role="img" aria-label="{html.escape(label)}: {html.escape(value)}">
  <title>{html.escape(label)}: {html.escape(value)}</title>
  <linearGradient id="s" x2="0" y2="100%">
    <stop offset="0" stop-color="#fff" stop-opacity=".18"/>
    <stop offset="1" stop-color="#000" stop-opacity=".08"/>
  </linearGradient>
  <clipPath id="r"><rect width="180" height="28" rx="4" fill="#fff"/></clipPath>
  <g clip-path="url(#r)">
    <rect width="92" height="28" fill="#374151"/>
    <rect x="92" width="88" height="28" fill="{color}"/>
    <rect width="180" height="28" fill="url(#s)"/>
  </g>
  <g fill="#fff" text-anchor="middle" font-family="Verdana,Geneva,DejaVu Sans,sans-serif" font-size="11" font-weight="700">
    <text x="46" y="18">coverage</text>
    <text x="136" y="18">{html.escape(value)}</text>
  </g>
</svg>
'''
with open(path, "w", encoding="utf-8") as f:
    f.write(svg)
PY

echo "wrote $badge_file ($coverage)"
