#!/bin/bash
# 01-nav.sh — CLI navigation commands

source "$(dirname "$0")/common.sh"

# ─────────────────────────────────────────────────────────────────
start_test "pinchtab nav <url>"

pt_ok nav "${FIXTURES_URL}/index.html"
assert_output_json
assert_output_contains "tabId" "returns tab ID"
assert_output_contains "title" "returns page title"

end_test

# ─────────────────────────────────────────────────────────────────
start_test "pinchtab nav (invalid URL)"

pt_fail nav "not-a-valid-url"

end_test

# ─────────────────────────────────────────────────────────────────
start_test "pinchtab tabs navigate <tabId> <url>"

# First navigate to get a tab
pt_ok nav "${FIXTURES_URL}/index.html"
TAB_ID=$(echo "$PT_OUT" | jq -r '.tabId')

# Navigate same tab to different page using tabs subcommand
pt_ok tabs navigate "$TAB_ID" "${FIXTURES_URL}/form.html"
assert_output_contains "form.html" "navigated to form.html"

end_test
