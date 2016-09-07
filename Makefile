# Copyright 2016 NetApp, Inc. All Rights Reserved.
GOOS=linux
GOARCH=amd64 d
GOGC=""

GO_PATH_VOLUME=netappdvp_go_path
GO=docker run --rm \
	-e GOOS=$(GOOS) \
	-e GOARCH=$(GOARCH) \
	-e GOGC=$(GOGC) \
	-v $(GO_PATH_VOLUME):/go \ 
	-v "$(PWD)":/go/src/github.com/netapp/netappdvp \
	-w /go/src/github.com/netapp/netappdvp \
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
	@$(GO) build -x -o /go/src/github.com/netapp/netappdvp/bin/netappdvp

install: build
	@$(GO) install

test:
	@$(GO) test github.com/netapp/netappdvp/...
