M := $(shell printf "\033[34;1m▶\033[0m")
VERSION := $(shell git describe 2>/dev/null || echo "undefined")
BUILD_ARGS := -ldflags "-X core.VERSION=$(VERSION)"

all: build

-include init/setup.mk

setup: deps hooks

deps: ; $(info $(M) Installing dependencies...)
	@./scripts/install-deps

hooks: ; $(info $(M) Installing commit hooks...)
	@./scripts/install-hooks

lint: ; $(info $(M) Lint projects...)
	@./scripts/utility go-lint identity
	@./scripts/utility go-lint proxy
	@./scripts/utility go-lint state

build: build-pre build-identity build-proxy build-state

build-pre: ; $(info $(M) Building projects...)
	@mkdir -p build/

build-identity:
	@pushd identity >/dev/null; \
	go build -o ../build/plantd-identity $(BUILD_ARGS) .; \
	popd >/dev/null

build-proxy:
	@pushd proxy >/dev/null; \
	go build -o ../build/plantd-proxy $(BUILD_ARGS) .; \
	popd >/dev/null

build-state:
	@pushd state >/dev/null; \
	go build -o ../build/plantd-state $(BUILD_ARGS) .; \
	popd >/dev/null

test: test-pre test-state

test-pre: ; $(info $(M) Testing projects...)

test-state:
	@pushd state >/dev/null; \
	go get -t ./...; \
	go test ./... -v; \
	popd >/dev/null

install: ; $(info $(M) Installing binaries...)
	@install build/plantd-* /usr/local/bin/

uninstall: ; $(info $(M) Uninstalling binaries...)
	@rm /usr/local/bin/plantd-*

clean: ; $(info $(M) Removing build files...)
	@rm -rf build/

.PHONY: all build clean
