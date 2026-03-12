#!/bin/bash
# 03-actions.sh — CLI action commands (click, type, press)

source "$(dirname "$0")/common.sh"

# ─────────────────────────────────────────────────────────────────
start_test "pinchtab click --css <selector>"

pt_ok nav "${FIXTURES_URL}/buttons.html"
pt_ok click --css "#increment"

end_test

# ─────────────────────────────────────────────────────────────────
start_test "pinchtab fill <selector> <text>"

pt_ok nav "${FIXTURES_URL}/form.html"
pt_ok fill "#username" "testuser"
assert_output_contains "filled" "confirms text was filled"

end_test

# ─────────────────────────────────────────────────────────────────
start_test "pinchtab press <key>"

pt_ok press Escape
assert_output_contains "pressed" "confirms key was pressed"

end_test

# ─────────────────────────────────────────────────────────────────
start_test "pinchtab press Enter (does not type text)"

pt_ok nav "${FIXTURES_URL}/form.html"
pt_ok fill "#username" "hello"
pt_ok press Enter

# Verify via eval that username doesn't contain "Enter"
pt_ok eval 'document.getElementById("username").value'
assert_output_not_contains "Enter" "Enter key dispatched, not typed"
assert_output_contains "hello" "original text preserved"

end_test

# ─────────────────────────────────────────────────────────────────
start_test "pinchtab focus --css <selector>"

pt_ok nav "${FIXTURES_URL}/form.html"
pt_ok focus --css "#email"

end_test

# ─────────────────────────────────────────────────────────────────
start_test "pinchtab hover --css <selector>"

pt_ok nav "${FIXTURES_URL}/buttons.html"
pt_ok hover --css "#increment"

end_test
