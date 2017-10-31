go ?= go
GOFILES = $(shell find . -name '*.go' -not -path './vendor/*')


_builds:
	mkdir -p _builds/{linux,osx} || echo "Directory $@ already there"

getdeps:
	go get -u ./...

_builds/linux/logrotate-sync: _builds $(GOFILES)
	GOOS=linux $(go) build -o _builds/linux/logrotate-sync .
	upx --brute _builds/linux/logrotate-sync

build-linux: _builds/linux/logrotate-sync

test:
	$(go) test -v ./...

.PHONY: build-linux
