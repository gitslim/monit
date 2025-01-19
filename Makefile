.PHONY: all test coverage staticcheck autotests clean

all: gen statictest autotests coverage

test:
	@echo "Running tests..."
	go test ./... "$@"

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
