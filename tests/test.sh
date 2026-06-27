#!/usr/bin/env bash
# Simple test harness for tinyfetch (both Shell and Go versions)
set -euo pipefail

failed=0

# Helper function for assertion
assert_contains() {
  local content="$1"
  local pattern="$2"
  local name="$3"
  if echo "$content" | grep -q "$pattern"; then
    echo "  [PASS] $name"
  else
    echo "  [FAIL] $name (Pattern '$pattern' not found)"
    failed=1
  fi
}

# shellcheck disable=SC2317
assert_eq() {
# shellcheck disable=SC2317
  local val1="$1"
  local val2="$2"
  local name="$3"
  if [ "$val1" -eq "$val2" ]; then
    echo "  [PASS] $name"
  else
    echo "  [FAIL] $name (Expected $val2, got $val1)"
    failed=1
  fi
}

run_suite() {
  local cmd="$1"
  local type="$2"
  echo "Running test suite for $type version: $cmd"

  # Test 1: Basic execution and exit code
  local out
  out=$($cmd)
  assert_contains "$out" "Host:" "$type: Outputs Host"
  assert_contains "$out" "OS:" "$type: Outputs OS"
  assert_contains "$out" "Kernel:" "$type: Outputs Kernel"
  assert_contains "$out" "Uptime:" "$type: Outputs Uptime"
  assert_contains "$out" "Shell:" "$type: Outputs Shell"
  assert_contains "$out" "CPU:" "$type: Outputs CPU"
  assert_contains "$out" "Memory:" "$type: Outputs Memory"
  assert_contains "$out" "Disk:" "$type: Outputs Disk"

  # Test 2: Help output
  local help_out
  help_out=$($cmd --help)
  assert_contains "$help_out" "Usage:" "$type: Help prints usage"

  # Test 3: --no-ascii flag output format (no logo, at least 8 lines of details)
  local no_ascii_out
  no_ascii_out=$($cmd --no-ascii)
  local line_count
  line_count=$(echo "$no_ascii_out" | grep -c -v "^$")
  if [ "$line_count" -ge 8 ]; then
    echo "  [PASS] $type: --no-ascii prints at least 8 lines ($line_count)"
  else
    echo "  [FAIL] $type: --no-ascii prints fewer than 8 lines ($line_count)"
    failed=1
  fi
  
  # Ensure no logo graphics in --no-ascii output
  if echo "$no_ascii_out" | grep -q -- "---"; then
    echo "  [FAIL] $type: --no-ascii output contains logo graphics"
    failed=1
  else
    echo "  [PASS] $type: --no-ascii does not contain logo graphics"
  fi
}

# Run Go version test (if compiled)
if [ -f "./tinyfetch" ]; then
  run_suite "./tinyfetch" "Go"
else
  echo "Error: tinyfetch binary not built"
  exit 1
fi

if [ "$failed" -eq 0 ]; then
  echo "ALL TESTS PASSED!"
  exit 0
else
  echo "SOME TESTS FAILED!"
  exit 1
fi
