# Copyright 2016 NetApp, Inc. All Rights Reserved.
# Test change
GOOS=linux
GOARCH=amd64
GOGC=""
GITHASH?=`git rev-parse HEAD || echo unknown`
BUILD_TYPE?=custom
BUILD_TYPE_REV?=0
BUILD_TIME=`date`

# Go compiler flags need to be properly encapsulated with double quotes to handle spaces in values
BUILD_FLAGS="-X \"main.GitHash=$(GITHASH)\" -X \"main.BuildType=$(BUILD_TYPE)\" -X \"main.BuildTypeRev=$(BUILD_TYPE_REV)\" -X \"main.BuildTime=$(BUILD_TIME)\""

# BUILD_TAG is intended to be an optional environment variable that uniquely
# identifies the build instance. It is intended to enable concurrent builds
# on the same machine.
GO_PATH_VOLUME="netappdvp_go_path_$(BUILD_TAG)"
GO=docker run --rm \
	-e GOOS=$(GOOS) \
	-e GOARCH=$(GOARCH) \
	-e GOGC=$(GOGC) \
	-v $(GO_PATH_VOLUME):/go \
	-v "$(PWD)":/go/src/github.com/netapp/netappdvp \
	-w /go/src/github.com/netapp/netappdvp \
	golang:1.8 go

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
	@$(GO) build -ldflags $(BUILD_FLAGS) -x -o /go/src/github.com/netapp/netappdvp/bin/netappdvp

install: build
	@$(GO) install

test:
	@$(GO) test github.com/netapp/netappdvp/...
