# This Makefile is meant to be used by people that do not usually work
# with Go source code. If you know what GOPATH is then you probably
# don't need to bother with make.

.PHONY: carrier all test clean

GOBIN = $(shell pwd)/build/bin
GO ?= latest
GPATH = $(shell go env GOPATH)
GORUN = env GO111MODULE=on GOPATH=$(GPATH) GOPROXY=https://goproxy.cn go run

carrier:
	$(GORUN) build/ci.go install ./cmd/carrier
	@echo "Done building."
	@echo "Run \"$(GOBIN)/carrier\" to launch carrier."

carrier-race:
	$(GORUN) build/ci.go install ./cmd/carrier --race
	@echo "Done building."
	@echo "Run \"$(GOBIN)/carrier\" to launch carrier."

kmstool:
	$(GORUN) build/ci.go install ./core/metispay/kms/kmstool
	@echo "Done building."
	@echo "Run \"$(GOBIN)/kmstool --help\" to get help."

all:
	$(GORUN) build/ci.go install ./cmd/
	$(GORUN) build/ci.go install ./core/metispay/kms/kmstool

test: all
	$(GORUN) build/ci.go test

clean:
	env GO111MODULE=on go clean -cache
	rm -fr build/_workspace/pkg/ $(GOBIN)/*