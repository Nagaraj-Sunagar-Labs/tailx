# TailX build & test pipeline
#
# Flow: build (host, cross-compile for linux) -> run lightweight linux
# container with the built binary volume mounted -> smoke tests inside.

GO        ?= go
APP       := tailx
MODULE    := github.com/Nagaraj-Sunagar-Labs/tailx
CMD       := ./cmd/tailx
DIST      := dist

# Target linux architecture (override: make build-linux GOARCH=arm64)
GOOS      ?= linux
GOARCH    ?= amd64

IMAGE     := tailx-test:latest
BIN       := $(DIST)/$(APP)

.PHONY: all build build-linux test vet fmt clean docker-image docker-test docker-shell

all: build

## build: compile the binary for the host platform
build:
	$(GO) build -o $(DIST)/$(APP) $(CMD)

## build-linux: cross-compile the binary for the target linux architecture
build-linux:
	CGO_ENABLED=0 GOOS=$(GOOS) GOARCH=$(GOARCH) $(GO) build -ldflags="-s -w" -o $(BIN) $(CMD)
	@echo "built $(BIN) for $(GOOS)/$(GOARCH)"

## test: run unit tests on the host
test:
	$(GO) test ./...

## vet / fmt: static checks
vet:
	$(GO) vet ./...

fmt:
	$(GO) fmt ./...

## clean: remove build artifacts
clean:
	rm -rf $(DIST)

## docker-image: build the lightweight linux test image
docker-image:
	docker build -t $(IMAGE) .

## docker-test: build for linux, then run smoke tests in the container
## with the binary volume mounted from $(DIST)
docker-test: build-linux docker-image
	docker run --rm \
		-v "$(CURDIR)/$(DIST):/opt/tailx/bin:ro" \
		$(IMAGE)

## docker-shell: interactive shell in the test container (binary mounted)
docker-shell: build-linux docker-image
	docker run --rm -it \
		-v "$(CURDIR)/$(DIST):/opt/tailx/bin:ro" \
		--entrypoint /bin/sh $(IMAGE)
