.PHONY: build test lint vet fmt check clean run tidy vuln staticcheck

BINARY := tttns
BUILD_ENV ?= dev
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_DATE := $(shell date -u +%Y-%m-%dT%H:%M:%SZ)
LDFLAGS := -X main.version=$(VERSION) \
	-X github.com/haoli000/tttns/generated/buildinfo.buildDate=$(BUILD_DATE) \
	-X github.com/haoli000/tttns/generated/buildinfo.buildVersion=$(VERSION) \
	-X github.com/haoli000/tttns/generated/buildinfo.commit=$(COMMIT)

build:
	go build -tags $(BUILD_ENV) -ldflags "$(LDFLAGS)" -o $(BINARY) .

run: build
	./$(BINARY)

test:
	go test ./...

tidy:
	go mod tidy

fmt:
	gofmt -w .

vet:
	go vet ./...

staticcheck:
	go run honnef.co/go/tools/cmd/staticcheck@latest ./...

vuln:
	go run golang.org/x/vuln/cmd/govulncheck@latest ./...

lint: fmt vet staticcheck

check: lint test vuln

clean:
	rm -f $(BINARY)
