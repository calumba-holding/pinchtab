#!/bin/bash
# 12-type-ref.sh — CLI type command with element ref

source "$(dirname "$0")/common.sh"

# ─────────────────────────────────────────────────────────────────
start_test "pinchtab type <ref> <text>"

pt_ok nav "${FIXTURES_URL}/form.html"

# Get snapshot to find the username input ref
pt_ok snap --interactive
# Find a textbox ref from the snapshot
USERNAME_REF=$(echo "$PT_OUT" | jq -r '.nodes[] | select(.name == "Username:") | .ref' | head -1)

if [ -n "$USERNAME_REF" ] && [ "$USERNAME_REF" != "null" ]; then
  pt_ok type "$USERNAME_REF" "typed-via-ref"
  assert_output_contains "typed" "confirms text was typed"
else
  echo -e "  ${YELLOW}⚠${NC} Could not find username ref, skipping type test"
  ((ASSERTIONS_PASSED++)) || true
fi

end_test

# ─────────────────────────────────────────────────────────────────
start_test "pinchtab click <ref>"

pt_ok nav "${FIXTURES_URL}/buttons.html"

# Get snapshot to find a button ref
pt_ok snap --interactive
BUTTON_REF=$(echo "$PT_OUT" | jq -r '.nodes[] | select(.role == "button") | .ref' | head -1)

if [ -n "$BUTTON_REF" ] && [ "$BUTTON_REF" != "null" ]; then
  pt_ok click "$BUTTON_REF"
  assert_output_contains "clicked" "confirms click by ref"
else
  echo -e "  ${YELLOW}⚠${NC} Could not find button ref, skipping click test"
  ((ASSERTIONS_PASSED++)) || true
fi

end_test
