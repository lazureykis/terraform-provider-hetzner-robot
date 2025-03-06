TEST?=$$(go list ./... | grep -v 'vendor')
TESTPACKAGES=$(shell go list ./... | grep -v /vendor/)
VERSION=0.1.0

default: build

build:
	go build -o terraform-provider-hetzner-robot

install: build
	mkdir -p ~/.terraform.d/plugins/registry.terraform.io/lazureykis/hetzner_robot/${VERSION}/$$(go env GOOS)_$$(go env GOARCH)
	cp terraform-provider-hetzner-robot ~/.terraform.d/plugins/registry.terraform.io/lazureykis/hetzner_robot/${VERSION}/$$(go env GOOS)_$$(go env GOARCH)

test:
	go test -i $(TEST) || exit 1
	echo $(TEST) | xargs -t -n4 go test $(TESTARGS) -timeout=30s -parallel=4

fmt:
	gofmt -w $$(find . -type f -name '*.go' -not -path "./vendor/*")
