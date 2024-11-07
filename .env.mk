REPO := github.com/dima-study/monmon

ARCH_LIST := amd64
OS_LIST := linux windows

PATH := $(shell go env GOPATH)/bin:$(PATH)

VERSION := develop
GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X $(REPO)/cmd/monmon-agent/build.Release="develop" -X $(REPO)/cmd/monmon-agent/build.Date=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X $(REPO)/cmd/monmon-agent/build.GitHash=$(GIT_HASH)

AGENT_BIN_windows_EXT := .exe

AGENT_BIN_EXT := $(AGENT_BIN_$(TARGET_OS)_EXT)
AGENT_BIN := "bin/monmon-agent.$(TARGET_OS)-$(TARGET_ARCH)$(AGENT_BIN_EXT)"
AGENT_CONFIG := "config/monmon.yaml"

DOCKER_IMG_AGENT := "monmon-agent:$(VERSION)"
