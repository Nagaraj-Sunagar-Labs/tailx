#!/usr/bin/env bash
# setup.sh - prepare the TailX testing environment inside the container.
#
# Available at /opt/tailx/setup.sh (root of the testing env).
# Recommended usage, so PATH changes persist in your current shell:
#
#     . ./setup.sh        # or: source ./setup.sh
#
set -u

TAILX_BIN="${TAILX_BIN:-/opt/tailx/bin/tailx}"
# Fixtures are generated next to the test scripts in the test folder.
TESTDIR="${TAILX_TESTDIR:-/opt/tailx/tests}"

echo "== tailx test env setup =="

# 1. Verify the volume mounted binary.
if [ ! -x "$TAILX_BIN" ]; then
	echo "WARN: $TAILX_BIN is missing or not executable."
	echo "      Run 'make docker-shell' on the host so dist/ gets mounted."
else
	echo "OK:   binary at $TAILX_BIN ($("$TAILX_BIN" --version 2>/dev/null || echo 'version unknown'))"
fi

# 2. Put the binary's directory on PATH (idempotent).
case ":$PATH:" in
*":/opt/tailx/bin:"*) ;;
*)
	export PATH="/opt/tailx/bin:$PATH"
	echo "OK:   /opt/tailx/bin added to PATH"
	;;
esac

# 3. Generate sample fixtures in the test folder.
seq 1 20 >"$TESTDIR/sample.log"          # 20 numbered lines
printf 'alpha\nbeta\ngamma\n' >"$TESTDIR/tiny.log"
: >"$TESTDIR/empty.log"                   # empty file edge case
printf 'no newline at end' >"$TESTDIR/no_newline.log"

echo "OK:   fixtures ready in $TESTDIR (sample.log, tiny.log, empty.log, no_newline.log)"
echo "== ready. try: cd $TESTDIR && tailx -n 5 sample.log =="
