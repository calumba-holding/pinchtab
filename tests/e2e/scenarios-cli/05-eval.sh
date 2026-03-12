#!/bin/bash
# 05-eval.sh — CLI evaluate command

source "$(dirname "$0")/common.sh"

# ─────────────────────────────────────────────────────────────────
start_test "pinchtab eval <expression>"

pt_ok nav "${FIXTURES_URL}/index.html"
pt_ok eval "1 + 1"
assert_output_contains "2" "evaluates simple expression"

end_test

# ─────────────────────────────────────────────────────────────────
start_test "pinchtab eval (DOM query)"

pt_ok nav "${FIXTURES_URL}/form.html"
pt_ok eval "document.title"
assert_output_contains "Form" "returns page title"

end_test

# ─────────────────────────────────────────────────────────────────
start_test "pinchtab eval (JSON result)"

pt_ok eval 'JSON.stringify({a: 1, b: 2})'
# Output is {"result": "{\"a\":1,\"b\":2}"} - escaped JSON
assert_output_contains 'a' "returns JSON object"
assert_output_contains 'b' "contains both keys"

end_test

# ─────────────────────────────────────────────────────────────────
start_test "pinchtab tabs eval <tabId> <expression>"

pt_ok nav "${FIXTURES_URL}/buttons.html"
TAB_ID=$(echo "$PT_OUT" | jq -r '.tabId')

pt_ok tabs eval "$TAB_ID" "document.title"
assert_output_contains "Button" "evaluates in correct tab"

end_test
