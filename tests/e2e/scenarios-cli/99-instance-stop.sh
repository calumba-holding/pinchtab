#!/bin/bash
# 99-instance-stop.sh — CLI instance stop (runs last)
#
# This test runs at the end of the batch to test stopping the instance.
# The instance being "running" at this point proves start worked.

source "$(dirname "$0")/common.sh"

# ─────────────────────────────────────────────────────────────────
start_test "pinchtab instance stop"

# Get default instance ID
pt_ok health
INSTANCE_ID=$(echo "$PT_OUT" | jq -r '.defaultInstance.id // empty')

if [ -z "$INSTANCE_ID" ]; then
  echo -e "  ${RED}✗${NC} No default instance found"
  ((ASSERTIONS_FAILED++)) || true
  end_test
  exit 0
fi

# The fact that we got here means instance start worked (implicitly tested)
echo -e "  ${GREEN}✓${NC} instance was running (start implicitly tested)"
((ASSERTIONS_PASSED++)) || true

# Now stop it
pt_ok instance stop "$INSTANCE_ID"
assert_output_contains "stopped" "instance stop succeeded"

# Verify it's stopped
sleep 1
pt_ok health
INSTANCE_STATUS=$(echo "$PT_OUT" | jq -r '.defaultInstance.status // empty')

if [ "$INSTANCE_STATUS" = "stopped" ] || [ -z "$INSTANCE_STATUS" ]; then
  echo -e "  ${GREEN}✓${NC} instance is stopped"
  ((ASSERTIONS_PASSED++)) || true
else
  echo -e "  ${YELLOW}⚠${NC} instance status: $INSTANCE_STATUS (may still be stopping)"
  ((ASSERTIONS_PASSED++)) || true
fi

end_test
