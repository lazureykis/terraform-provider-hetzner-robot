TEST?=$$(go list ./... | grep -v 'vendor')
TESTPACKAGES=$(shell go list ./... | grep -v /vendor/)
VERSION=0.1.0

default: build

build:
	go build -o terraform-provider-hetznerrobot

install: build
	mkdir -p ~/.terraform.d/plugins/registry.terraform.io/lazureykis/hetznerrobot/${VERSION}/$$(go env GOOS)_$$(go env GOARCH)
	cp terraform-provider-hetznerrobot ~/.terraform.d/plugins/registry.terraform.io/lazureykis/hetznerrobot/${VERSION}/$$(go env GOOS)_$$(go env GOARCH)

test: 
	go test -i $(TEST) || exit 1                                                   
	echo $(TEST) | xargs -t -n4 go test $(TESTARGS) -timeout=30s -parallel=4                    

testacc: 
	TF_ACC=1 go test $(TESTPACKAGES) -v $(TESTARGS) -timeout 120m

fmt:
	gofmt -w $$(find . -type f -name '*.go' -not -path "./vendor/*")

generate:
	go generate ./...

.PHONY: build install test testacc fmt generate
