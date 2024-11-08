ifneq (,$(wildcard ./.env.mk))
    include .env.mk
    export
endif


ifeq ($(OS),Windows_NT)     # is Windows_NT on XP, 2000, 7, Vista, 10...
    BUILDOS := Windows
else
    BUILDOS := $(shell uname | tr A-Z a-z)
endif


BUILDARCH ?= $(shell uname -m)

ifeq ($(BUILDARCH),x86_64)
	BUILDARCH := amd64
endif


TARGET_OS := $(BUILDOS)
TARGET_ARCH := $(BUILDARCH)


.PHONY: help
help:
	@printf "%-20s %s\n" "Target" "Description"
	@printf "%-20s %s\n" "------" "-----------"
	@make -pqR : 2>/dev/null \
		| awk -v RS= -F: '/# Files/,/^# Finished Make data base/ {if ($$1 !~ "^[#.]") {print $$1}}' \
    | sort \
    | grep -v -e '^[^[:alnum:]]' -e '^$@$$' -e '-meta$$' \
    | xargs -I _ sh -c 'printf "%-20s " _; make _ -nB | (grep -i "^# Help:" || echo "") | tail -1 | sed "s/^# Help: //g"'


#
# GENERATE
#
.PHONY: generate
generate: generate-proto
	@# Help: Do generate sources


.PHONY: generate-proto
generate-proto: install-deps-protoc
	@# Help: Generate from proto files
	protoc \
		-Iapi/vendor/proto/alta/protopatch/ \
		-Iapi/proto \
		--go-patch_out=plugin=go,paths=source_relative:pkg/api/proto \
		--go-patch_out=plugin=go-grpc,paths=source_relative,require_unimplemented_servers=true:pkg/api/proto \
		./api/proto/stats/v1/*.proto


#
# BUILD
#

.PHONY: build
build: $(addprefix build-sub-,$(ARCH_LIST))
	@# Help: Build multi-target monmon apps: ARCH_LIST="$(ARCH_LIST)", OS_LIST="$(OS_LIST)"
build-sub-%:
	$(MAKE) build-app-arch-meta TARGET_ARCH=$*

.PHONY: build-app-arch-meta
build-app-arch-meta: $(addprefix build-app-arch-sub-,$(OS_LIST))
build-app-arch-sub-%:
	$(MAKE) build-app TARGET_OS=$*

.PHONY: build-app
build-app: build-app-agent
	@# Help: Build monmon apps: TARGET_OS="$(TARGET_OS)", TARGET_ARCH="$(TARGET_ARCH)"


.PHONY: build-app-agent
build-app-agent:
	@# Help: Build monmon agent: TARGET_OS="$(TARGET_OS)", TARGET_ARCH="$(TARGET_ARCH)"
	$(MAKE) build-app-agent-meta
build-app-agent-meta: generate
	CGO_ENABLED=0 GOOS=$(TARGET_OS) GOARCH=$(TARGET_ARCH) go build -v -o $(AGENT_BIN) -ldflags "$(LDFLAGS)" ./cmd/monmon-agent


.PHONY: build-img-all
build-img-all: build-img-agent
	@# Help: Build all docker images


.PHONY: build-img-agent
build-img-agent:
	@# Help: Build monmon agent docker image $(DOCKER_IMG_AGENT)
	docker build \
		-t $(DOCKER_IMG_AGENT) \
		-f build/monmon-agent/Dockerfile .


#
# RUN
#

.PHONY: run
run: build-app
	@# Help: Run monmon agent with default config $(AGENT_CONFIG): TARGET_OS="$(TARGET_OS)", TARGET_ARCH="$(TARGET_ARCH)"
	$(MAKE) run-meta
run-meta:
	$(AGENT_BIN) start -config $(AGENT_CONFIG)


#
# Deps
#

.PHONY: install-deps-all
install-deps-all: install-deps-lint install-deps-protoc
	@# Help: Install all deps


.PHONY: install-deps-lint
install-deps-lint:
	@# Help: Install golangci linter
	@(which golangci-lint 1>/dev/null 2>&1) && echo "golangci-lint is already installed" \
		|| ( echo "install golangci-lint..." \
				&& curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.57.2 \
				&& echo "done" \
			)


.PHONY: install-deps-protoc
install-deps-protoc:
	@# Help: Install protoc generator deps
	@ (which protoc-gen-go 1>/dev/null 2>&1) && echo "protoc-gen-go is already installed" \
		|| ( echo "install protoc-gen-go..." \
				&& go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.35 \
				&& echo "done" \
			)
	@ (which protoc-gen-go-grpc 1>/dev/null 2>&1) && echo "protoc-gen-go-grpc is already installed" \
		|| ( echo "install protoc-gen-go-grpc..." \
				&& go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.5 \
				&& echo "done" \
			)
	@ (which protoc-gen-go-patch 1>/dev/null 2>&1) && echo "protoc-gen-go-patch is already installed" \
		|| ( echo "install protoc-gen-go-patch..." \
				&& go install github.com/alta/protopatch/cmd/protoc-gen-go-patch@v0.5 \
				&& echo "done" \
			)


#
# Util commands
#

.PHONY: test
test:
	@# Help: Run tests
	go test -race -count 100 ./...


.PHONY: lint
lint: install-deps-lint
	@# Help: Lint the project
	golangci-lint run --config=.golangci.yml ./...
