# This Makefile is meant to be used by people that do not usually work
# with Go source code. If you know what GOPATH is then you probably
# don't need to bother with make.
#

NAME := aquapool
#export GO111MODULE=on


.PHONY: all test clean deps
GOBIN = build/bin
OUTDIR = ${PWD}/${GOBIN}
all: ${GOBIN}/${NAME}
	@echo "Build complete"
	@file ${GOBIN}/${NAME} || true
	@echo "Consider running 'make deps' to bump dependencies"

.PHONY += all
deps:
	CGO_ENABLED=0 go get -v -u gitlab.com/aquachain/aquachain@master
	CGO_ENABLED=0 go get -v ./...
.PHONY += deps
${GOBIN}/${NAME}: *.go */*.go */*/*.go
	#go get -v -u -d gitlab.com/aquachain/aquachain
	CGO_ENABLED=0 go build -trimpath -tags 'netgo osusergo static' -ldflags '-s -w' -v -o $@

release: ${GOBIN}/${NAME} ${GOBIN}/frontend.tar.gz
	rm -rf ./release
	mkdir -p ./release
	cp ${GOBIN}/${NAME} ./release/
	cp ${GOBIN}/frontend.tar.gz ./release/
	cd release && tar czf ${NAME}.tar.gz ${NAME}

test: all
	CGO_ENABLED=0 go test -v ./...

clean:
	${RM} -rf build/_workspace/pkg/
	${RM} -rf build/_workspace/src/

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

