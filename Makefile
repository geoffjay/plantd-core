M := $(shell printf "\033[34;1m▶\033[0m")
VERSION := $(shell git describe 2>/dev/null || echo "undefined")
SHELL := /bin/bash
BUILD_ARGS := -ldflags "-X core.VERSION=$(VERSION)"
TEST_ARGS := $(shell if [ ! -z ${COVERAGE} ]; then echo "-race -coverprofile=coverage.txt -covermode=atomic"; fi)

all: build

-include init/setup.mk
-include logger/manage.mk

setup: deps hooks

deps: ; $(info $(M) Installing dependencies...)
	@./scripts/install-deps

hooks: ; $(info $(M) Installing commit hooks...)
	@./scripts/install-hooks

lint: ; $(info $(M) Lint projects...)
	@./scripts/utility go-lint core
	@./scripts/utility go-lint identity
	@./scripts/utility go-lint logger
	@./scripts/utility go-lint proxy
	@./scripts/utility go-lint state

build: build-pre build-client build-identity build-logger build-proxy build-state

build-pre: ; $(info $(M) Building projects...)
	@mkdir -p build/

build-client:
	@pushd client >/dev/null; \
	go build -o ../build/plant $(BUILD_ARGS) .; \
	popd >/dev/null

build-identity:
	@pushd identity >/dev/null; \
	go build -o ../build/plantd-identity $(BUILD_ARGS) .; \
	popd >/dev/null

build-logger:
	@pushd logger >/dev/null; \
	go build -o ../build/plantd-logger $(BUILD_ARGS) .; \
	popd >/dev/null

build-proxy:
	@pushd proxy >/dev/null; \
	go build -o ../build/plantd-proxy $(BUILD_ARGS) .; \
	popd >/dev/null

build-state:
	@pushd state >/dev/null; \
	go build -o ../build/plantd-state $(BUILD_ARGS) .; \
	popd >/dev/null

test: test-pre test-core test-state

test-pre: ; $(info $(M) Testing projects...)
	@mkdir -p coverage/

test-core:
	@pushd core >/dev/null; \
	go test $(TEST_ARGS) ./... -v; \
	if [[ -f coverage.txt ]]; then mv coverage.txt ../coverage/core.txt; fi; \
	popd >/dev/null

test-state:
	@pushd state >/dev/null; \
	go test $(TEST_ARGS) ./... -v; \
	if [[ -f coverage.txt ]]; then mv coverage.txt ../coverage/state.txt; fi; \
	popd >/dev/null

# live reload helpers
dev-state:
	@air -c state/.air.toml

dev-logger:
	@air -c logger/.air.toml

# notebooks for new ideas
jupyter:
	@mkdir -p notebooks
	@docker run -it -p 8888:8888 -v notebooks:/notebooks gopherdata/gophernotes:latest-ds

install: ; $(info $(M) Installing binaries...)
	@install build/plantd-* /usr/local/bin/

uninstall: ; $(info $(M) Uninstalling binaries...)
	@rm /usr/local/bin/plantd-*

clean: ; $(info $(M) Removing build files...)
	@rm -rf build/
	@rm -rf coverage/

.PHONY: all build clean
