# This Makefile is meant to be used by people that do not usually work
# with Go source code. If you know what GOPATH is then you probably
# don't need to bother with make.
#

NAME := aquapool


.PHONY: all test clean deps
GOBIN = build/bin
OUTDIR = ${PWD}/${GOBIN}
all: deps ${GOBIN}/${NAME} frontend

.PHONY += all
deps:
	go get -d -v ./...
.PHONY += deps
${GOBIN}/${NAME}:
	#go get -v -u -d gitlab.com/aquachain/aquachain
	CGO_ENABLED=0 go build -tags 'netgo osusergo static' -ldflags '-s -w' -v -o $@

test: all
	build/env.sh go test -v ./...

clean:
	${RM} -rf build/_workspace/pkg/ $(GOBIN)/*
	${RM} -rf build/_workspace/src/ $(GOBIN)/*

clean-www:
	${RM} -rf www/dist 

frontend: ${OUTDIR}/frontend.tar.gz

${OUTDIR}/frontend.tar.gz: www/dist
	cd www/dist && tar czf $@ .


www/dist: $(wildcard www/app/*.* www/app/*/*.*)
	cd www && npm install -g ember-cli
	cd www && npm install -g bower
	cd www && npm install
	cd www && bower install
	cd www && ./build.sh

