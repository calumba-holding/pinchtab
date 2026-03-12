#!/bin/bash
# 10-quick.sh — CLI quick command (nav + snap combined)

source "$(dirname "$0")/common.sh"

# ─────────────────────────────────────────────────────────────────
start_test "pinchtab quick <url>"

pt_ok quick "${FIXTURES_URL}/form.html"
# Quick command navigates and returns snapshot
assert_output_contains "nodes" "returns snapshot nodes"
assert_output_contains "form.html" "navigated to correct page"

end_test
