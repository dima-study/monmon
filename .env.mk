REPO := github.com/dima-study/monmon

PATH := $(shell go env GOPATH)/bin:$(PATH)

VERSION := develop
GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X $(REPO)/cmd/monmon-agent/build.Release="develop" -X $(REPO)/cmd/monmon-agent/build.Date=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X $(REPO)/cmd/monmon-agent/build.GitHash=$(GIT_HASH)

AGENT_BIN := "bin/monmon-agent"
AGENT_CONFIG := "config/monmon.yaml"

DOCKER_IMG_AGENT := "monmon-agent:$(VERSION)"
