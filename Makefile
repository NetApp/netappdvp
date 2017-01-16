# Copyright 2016 NetApp, Inc. All Rights Reserved.
GOOS=linux
GOARCH=amd64
GOGC=""

# BUILD_TAG is intended to be an optional environment variable that uniquely
# identifies the build instance. It is intended to enable concurrent builds
# on the same machine.
GO_PATH_VOLUME="netappdvp_go_path_$(BUILD_TAG)"
GO=docker run --rm \
	-e GOOS=$(GOOS) \
	-e GOARCH=$(GOARCH) \
	-e GOGC=$(GOGC) \
	-v $(GO_PATH_VOLUME):/go \
	-v "$(PWD)":/go/src/github.com/ebalduf/netappdvp \
	-w /go/src/github.com/ebalduf/netappdvp \
	golang:1.6 go

.PHONY=clean default fmt get install test

default: build

clean:
	rm -f $(PWD)/bin/netappdvp
	docker volume rm $(GO_PATH_VOLUME) || true

fmt:
	@$(GO) fmt ./...

get:
	@$(GO) get -v

build: get *.go
	@mkdir -p $(PWD)/bin
	@$(GO) build -x -o /go/src/github.com/ebalduf/netappdvp/bin/netappdvp

install: build
	@$(GO) install

test:
	@$(GO) test github.com/ebalduf/netappdvp/...
