# Requirements: git, go, vgo
NAME     := plant_erd
VERSION  := $(shell cat VERSION)
REVISION := $(shell git rev-parse --short HEAD)

SRCS    := $(shell find . -type f -name '*.go')
LDFLAGS := "-s -w -X \"main.Version=$(VERSION)\" -X \"main.Revision=$(REVISION)\""

.DEFAULT_GOAL := bin/$(NAME)

export GO111MODULE ?= on

bin/$(NAME): $(SRCS)
	go build -ldflags=$(LDFLAGS) -o bin/$(NAME)

.PHONY: gox
gox:
	gox -osarch="$${GOX_OSARCH}" -ldflags=$(LDFLAGS) -output="bin/$(NAME)_{{.OS}}_{{.Arch}}"

.PHONY: zip
zip:
	cd bin ; \
	for file in *; do \
		zip $${file}.zip $${file} ; \
	done

.PHONY: gox_with_zip
gox_with_zip: clean gox zip

.PHONY: clean
clean:
	rm -rf bin/*

.PHONY: tag
tag:
	git tag -a $(VERSION) -m "Release $(VERSION)"
	git push --tags

.PHONY: release
release: tag
	git push origin master

.PHONY: test
test:
	go test -count=1 $${TEST_ARGS} ./...

.PHONY: testrace
testrace:
	go test -count=1 $${TEST_ARGS} -race ./...

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: fmtci
fmtci:
	! gofmt -d . | grep '^'

.PHONY: lint
lint:
	golint -set_exit_status ./...

.PHONY: vet
vet:
	go vet ./...

.PHONY: integration_test
integration_test: bin/$(NAME)
	go test _integration/check_readme_test.go

.PHONY: test_all
test_all: test testrace fmt lint vet integration_test
