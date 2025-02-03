.PHONY: all build-server build-agent build test coverage statictest autotests lint gen godoc-server clean

all: gen statictest autotests coverage build

BUILD_COMMIT := $(shell git rev-parse --short HEAD)
BUILD_DATE := $(shell git log -1 --format=%cd --date=format:"%Y%m%d")
BUILD_VERSION := v1.0.0
FLAGS := -ldflags "-X main.buildVersion=$(BUILD_VERSION) -X main.buildDate=$(BUILD_DATE) -X main.buildCommit=$(BUILD_COMMIT)"

build-server:
	@echo "Building server..."
	go build -o ./cmd/server/server $(FLAGS) ./cmd/server

build-agent:
	@echo "Building agent..."
	go build -o ./cmd/agent/agent $(FLAGS) ./cmd/agent

build: build-server build-agent

test:
	@echo "Running tests..."
	go test ./...

coverage:
	@echo "Calculating test coverage..."
	go test ./... -coverprofile=cover.out
	go tool cover -func=cover.out | grep total

statictest:
	@echo "Running statictest..."
	./scripts/statictest.sh

autotests:
	@echo "Running autotests..."
	./scripts/autotests.sh

lint:
	@echo "Running staticlint..."
	go run cmd/staticlint/main.go ./...

gen:
	@echo "Generating all code..."
	go generate ./...

godoc-serve:
	@echo "Starting godoc server..."
	rm -rf /tmp/.monit-godoc; mkdir -p /tmp/.monit-godoc/src && cp -r . /tmp/.monit-godoc/src/monit && godoc -play -http=:6060 -goroot /tmp/.monit-godoc

clean:
	@echo "Cleaning up..."
	rm -f cover.out
	rm -rf /tmp/.monit-godoc
