#!/usr/bin/env bash
# Smoke tests for the tailx binary running inside the linux test container.
# The binary is volume mounted at /opt/tailx/bin/tailx.
set -u

TAILX="${TAILX_BIN:-/opt/tailx/bin/tailx}"

pass=0
fail=0

check() {
	local desc="$1"; shift
	if "$@" >/dev/null 2>&1; then
		echo "PASS: $desc"
		pass=$((pass + 1))
	else
		echo "FAIL: $desc"
		fail=$((fail + 1))
	fi
}

echo "== tailx smoke tests =="

# 1. Binary exists and is executable.
check "binary is executable" test -x "$TAILX"

# 2. --version prints something.
check "--version runs" "$TAILX" --version

# 3. Basic tail behaviour: last 2 lines of a 10 line file.
workdir="$(mktemp -d)"
seq 1 10 > "$workdir/input.txt"
expected=$'9\n10'
actual="$("$TAILX" -n 2 "$workdir/input.txt" 2>/dev/null)"
check "outputs last 2 lines with -n 2" test "$actual" = "$expected"

rm -rf "$workdir"

echo "== results: $pass passed, $fail failed =="
[ "$fail" -eq 0 ]
