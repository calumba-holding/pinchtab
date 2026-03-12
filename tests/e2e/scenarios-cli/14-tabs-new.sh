#!/bin/bash
# 14-tabs-new.sh — CLI tabs new command

source "$(dirname "$0")/common.sh"

# ─────────────────────────────────────────────────────────────────
start_test "pinchtab tabs new <url>"

pt_ok tabs new "${FIXTURES_URL}/buttons.html"
assert_output_json
assert_output_contains "tabId" "returns new tab ID"

# Save the tab ID for cleanup
NEW_TAB_ID=$(echo "$PT_OUT" | jq -r '.tabId')

# Verify it appears in tabs list
pt_ok tabs
assert_output_contains "$NEW_TAB_ID" "new tab appears in list"

# Cleanup
if [ -n "$NEW_TAB_ID" ] && [ "$NEW_TAB_ID" != "null" ]; then
  pt tabs close "$NEW_TAB_ID" > /dev/null 2>&1 || true
fi

end_test
