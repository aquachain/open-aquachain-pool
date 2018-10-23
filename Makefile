# This Makefile is meant to be used by people that do not usually work
# with Go source code. If you know what GOPATH is then you probably
# don't need to bother with make.

.PHONY: all test clean

GOBIN = build/bin

all:
	go build -o aquapool

old-all:
	#go get -v -u -d gitlab.com/aquachain/aquachain
	GOBIN=${PWD}/build/bin go install -tags 'netgo osusergo static' -ldflags '-s -w -linkmode external -extldflags -static' -v ./...

test: all
	build/env.sh go test -v ./...

clean:
	rm -fr build/_workspace/pkg/ $(GOBIN)/*
	rm -fr build/_workspace/src/ $(GOBIN)/*
