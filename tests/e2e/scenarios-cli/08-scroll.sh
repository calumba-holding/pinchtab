#!/bin/bash
# 08-scroll.sh — CLI scroll command

source "$(dirname "$0")/common.sh"

# ─────────────────────────────────────────────────────────────────
start_test "pinchtab scroll <pixels>"

pt_ok nav "${FIXTURES_URL}/table.html"

# Scroll down by pixels
pt_ok scroll 100
assert_output_contains "scrolled" "confirms scroll action"

end_test

# ─────────────────────────────────────────────────────────────────
start_test "pinchtab scroll down/up"

pt_ok scroll down
assert_output_contains "scrolled" "scroll down succeeded"

pt_ok scroll up
assert_output_contains "scrolled" "scroll up succeeded"

end_test
