#!/usr/bin/env bash
set -euo pipefail

usage() {
  cat <<'EOF'
Usage:
  ./dev bench pinchtab [options...]
  ./dev bench agent-browser [options...]

Examples:
  ./dev bench pinchtab --dry-run
  ANTHROPIC_API_KEY=... ./dev bench pinchtab --model claude-haiku-4-5-20251001
  ANTHROPIC_API_KEY=... ./dev bench pinchtab --groups 0,1,2,3
  ANTHROPIC_API_KEY=... ./dev bench agent-browser --max-turns 150

For optimization baseline (no API keys required):
  ./dev opt baseline
EOF
}

if [[ $# -lt 1 ]]; then
  usage
  exit 1
fi

mode="$1"
shift

case "${mode}" in
  pinchtab|agent-browser)
    exec go run ./tests/tools/runner --lane "${mode}" --finalize --terse-summary "$@"
    ;;
  -h|--help|help)
    usage
    ;;
  *)
    echo "ERROR: unknown benchmark mode: ${mode}" >&2
    usage
    exit 1
    ;;
esac
