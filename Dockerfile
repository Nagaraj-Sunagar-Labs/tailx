# Lightweight linux test environment for TailX.
#
# The binary is NOT copied into this image. It is volume mounted from the
# host at runtime (see Makefile target `docker-test`), after it has been
# cross-compiled with `make build-linux GOARCH=<target arch>`.

FROM alpine:3.20

# Basic linux functionality: a shell, coreutils and a few utilities
# commonly needed when exercising a tail-like tool.
RUN apk add --no-cache bash coreutils grep diffutils

# Where the built tailx binary gets mounted (read-only).
VOLUME ["/opt/tailx/bin"]

# Smoke test scripts baked into the image.
COPY test/ /opt/tailx/tests/

# Environment setup script at the root of the testing env.
COPY test/setup.sh /opt/tailx/setup.sh

WORKDIR /opt/tailx

ENV TAILX_BIN=/opt/tailx/bin/tailx

# Run the smoke tests against the mounted binary.
CMD ["/opt/tailx/tests/smoke.sh"]
