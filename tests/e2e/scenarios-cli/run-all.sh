#!/bin/bash
# Run all CLI E2E test scenarios

set -uo pipefail

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"

# Source common utilities (initializes counters)
source "$SCRIPT_DIR/common.sh"

echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "  🦀 PinchTab CLI E2E Tests"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
echo "  Server: $PINCHTAB_URL"
echo "  Fixtures: $FIXTURES_URL"
echo ""

# Wait for server AND instance to be ready
echo "Waiting for pinchtab server..."
for i in $(seq 1 30); do
  if curl -s "$PINCHTAB_URL/health" > /dev/null 2>&1; then
    echo "Server responding..."
    break
  fi
  sleep 1
done

echo "Waiting for Chrome instance to be ready..."
for i in $(seq 1 60); do
  # Try a simple snapshot - if it works, Chrome is ready
  if PINCHTAB_URL="$PINCHTAB_URL" pinchtab health 2>/dev/null | grep -q '"ready":true'; then
    echo "Instance ready!"
    break
  fi
  echo "  waiting... ($i/60)"
  sleep 2
done

# Verify pinchtab CLI is available
if ! command -v pinchtab &> /dev/null; then
  echo "ERROR: pinchtab CLI not found in PATH"
  exit 1
fi

echo ""
echo "Running CLI tests..."
echo ""

# Find and run all test scripts in order
for script in "$SCRIPT_DIR"/[0-9][0-9]-*.sh; do
  if [ -f "$script" ]; then
    echo -e "${YELLOW}Running: $(basename "$script")${NC}"
    echo ""
    source "$script"
    echo ""
  fi
done

print_summary
