.PHONY: build clean package generate fmt serve
# VERSION := $(shell git describe --always)
VERSION := 0.1
GOOS ?= darwin
GOARCH ?= amd64

build: generate fmt
	@echo "Compiling source for $(GOOS) $(GOARCH)"
	@mkdir -p build
	@GOOS=$(GOOS) GOARCH=$(GOARCH) go build -ldflags "-X main.version=$(VERSION)" -o build/trinity$(BINEXT)

clean:
	@echo "Cleaning up workspace"
	@rm -rf build
	@rm -rf dist/$(VERSION)

package: clean build
	@echo "Creating package for $(GOOS) $(GOARCH)"
	@mkdir -p dist/$(VERSION)
	@cp build/* dist/$(VERSION)
	@cd dist/$(VERSION) && tar -pczf ../trinity_$(VERSION)_$(GOOS)_$(GOARCH).tar.gz .
	@rm -rf dist/$(VERSION)

generate:
	@echo "Running go generate"
	@go generate main.go

fmt:
	@echo "Formating needle project"
	@go fmt github.com/k4jt/trinity/...

serve: build
	@echo "Starting Trinity DB UI"
	./build/trinity
