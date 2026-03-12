#!/bin/bash
# 02-snap.sh — CLI snapshot commands

source "$(dirname "$0")/common.sh"

# ─────────────────────────────────────────────────────────────────
start_test "pinchtab snap (JSON format)"

pt_ok nav "${FIXTURES_URL}/form.html"
pt_ok snap
assert_output_json
assert_output_contains "nodes" "returns nodes array"
assert_output_contains "title" "returns page title"

end_test

# ─────────────────────────────────────────────────────────────────
start_test "pinchtab snap --text"

pt_ok snap --text
assert_output_contains "e0" "contains element refs"

end_test

# ─────────────────────────────────────────────────────────────────
start_test "pinchtab snap --interactive"

pt_ok snap --interactive
assert_output_json
assert_output_contains "textbox" "contains form inputs"
assert_output_contains "button" "contains buttons"

end_test

# ─────────────────────────────────────────────────────────────────
start_test "pinchtab snap --tab <id>"

pt_ok nav "${FIXTURES_URL}/buttons.html"
TAB_ID=$(echo "$PT_OUT" | jq -r '.tabId')

pt_ok snap --tab "$TAB_ID"
assert_output_contains "buttons.html" "snapshot from correct tab"

end_test
