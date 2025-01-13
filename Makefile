.PHONY: all test coverage staticcheck autotests clean

all: test coverage statictest autotests

test:
	@echo "Running tests..."
	go test ./... -v "$@"

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

clean:
	@echo "Cleaning up..."
	rm -f cover.out
