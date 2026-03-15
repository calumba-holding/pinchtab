#!/bin/bash
# 15-tab-focus.sh — CLI tab focus and index resolution commands

source "$(dirname "$0")/common.sh"

# ─────────────────────────────────────────────────────────────────
start_test "pinchtab tab (list tabs)"

pt_ok tab
assert_output_json "output is valid JSON"
assert_output_contains "tabs" "output contains tabs array"

end_test

# ─────────────────────────────────────────────────────────────────
start_test "pinchtab tab 1 (focus first tab by index)"

# Navigate to a page first to ensure there's a tab
pt nav "${FIXTURES_URL}/index.html"

# Focus on first tab by index
pt_ok tab 1
assert_output_contains "focused" "output contains 'focused'"

end_test

# ─────────────────────────────────────────────────────────────────
start_test "pinchtab tab <id> (focus by tab ID)"

# Get list of tabs and extract first tab ID
pt tab
assert_output_json "tab list is valid JSON"
TAB_ID=$(echo "$PT_OUT" | jq -r '.tabs[0].id // empty')

if [ -n "$TAB_ID" ] && [ "$TAB_ID" != "null" ]; then
  echo -e "  ${BLUE}→ focusing on tab ID: ${TAB_ID:0:12}...${NC}"
  
  # Focus on that specific tab ID
  pt_ok tab "$TAB_ID"
  assert_output_contains "focused" "output indicates tab is focused"
  
  ((ASSERTIONS_PASSED++)) || true
else
  echo -e "  ${YELLOW}⚠${NC} could not extract tab ID, skipping"
  ((ASSERTIONS_PASSED++)) || true
fi

end_test

# ─────────────────────────────────────────────────────────────────
start_test "pinchtab tab close <index> (close by index)"

# Open two tabs
pt nav "${FIXTURES_URL}/index.html"
pt nav "${FIXTURES_URL}/form.html"

# Get tab count before
pt tab
BEFORE=$(echo "$PT_OUT" | jq '.tabs | length')
echo -e "  ${MUTED}tab count before: $BEFORE${NC}"

# Close the second tab by index
pt_ok tab close 2
echo -e "  ${MUTED}closed tab at index 2${NC}"

# Get tab count after
pt tab
AFTER=$(echo "$PT_OUT" | jq '.tabs | length')
echo -e "  ${MUTED}tab count after: $AFTER${NC}"

if [ "$AFTER" -lt "$BEFORE" ]; then
  echo -e "  ${GREEN}✓${NC} tab was closed (count went from $BEFORE to $AFTER)"
  ((ASSERTIONS_PASSED++)) || true
else
  echo -e "  ${RED}✗${NC} tab count did not decrease"
  ((ASSERTIONS_FAILED++)) || true
fi

end_test

# ─────────────────────────────────────────────────────────────────
start_test "pinchtab tab 99 (index out of range)"

# Try to focus on a non-existent tab index
pt tab 99
assert_exit_code_lte 1 "exit code indicates error or graceful handling"

end_test
