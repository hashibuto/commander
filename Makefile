UID ?= $(shell id --user)
GID ?= $(shell id --group)
PWD := $(shell pwd)
GO_IMAGE := "golang:1.22"
GO := docker run --rm -u $(UID):$(GID) -e HOME=$$HOME -v $$HOME:$$HOME -v $(PWD):/build -w /build $(GO_IMAGE) go

.PHONY:
test:
	$(GO) test