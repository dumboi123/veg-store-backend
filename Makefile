# -v : show details log per test
# -cover : enable coverage mode
# -count=1 : no cache test result
TEST_FLAGS=-v -cover -count=1
PKG=./... # All sub packages in root
PKG_INTERNAL=./internal/...
PKG_TEST=./test/unit/...

.PHONY: test test-coverage test-one lint tidy

# Run all tests with coverage info
# Usage:
#   make test PKG=./test/unit/handler/rest_test
test:
	@echo "Running all unit tests..."
	@go test $(PKG) $(TEST_FLAGS)

# Run all tests and generate a coverage report
coverage:
	@echo "Running tests with coverage report..."
	@mkdir -p ./test/report
	@go test ./test/unit/... \
		-coverpkg=./internal/... \
		-coverprofile=./test/report/coverage.out \
		-covermode=atomic
	@go tool cover -func=./test/report/coverage.out | tail -n 1
	@go tool cover -html=./test/report/coverage.out -o ./test/report/index.html
	@echo "Coverage report saved to ./test/report/index.html"

# Run a specific test
# Usage:
#   make test-one TEST=TestUserHandler/TestHello_success
#   make test-one PKG=./test/unit/handler/rest_test TEST=TestHello_success
test-one:
	@echo "Running specific test: $(TEST)"
	@go test $(PKG) -v -run $(TEST)

# Run go vet and golangci-lint
lint:
	@echo "Running go vet..."
	@go vet $(PKG)
	@echo "Checking golangci-lint..."
	@if ! command -v golangci-lint >/dev/null 2>&1; then \
		echo "'golangci-lint' not found. Please run 'make prepare' first"; \
	fi
	@golangci-lint run

# Tidy up modules
tidy:
	@echo "Tidying Go modules..."
	@go mod tidy

prepare:
	@echo "Preparing external packages..."
	@go install -v github.com/air-verse/air@latest
	@go install -v github.com/nicksnyder/go-i18n/v2/goi18n@latest
	@go install -v github.com/swaggo/swag/cmd/swag@latest
	@#go install -v github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@echo "Done"

start:
	@echo "Starting..."
	@sudo chown -R $(shell id -u):$(shell id -g) .
	@sudo chmod +x scripts/*.sh
	@./scripts/start.sh

restart:
	@echo "Restarting..."
	@./scripts/restart.sh

stop:
	@echo "Stopping..."
	@./scripts/stop.sh

run:
	@echo "Running..."
	@docker exec -it -uroot veg-store-backend bash -c 'cd /app && go run cmd/server/main.go'

run-dev:
	@echo "Running with Hot reload..."
	@docker exec -it -uroot veg-store-backend bash -c 'cd /app && air -c .air.toml'

swagger:
	@echo "Renewing swagger schema..."
	@if ! command -v swag >/dev/null 2>&1; then \
		echo "'swag' not found. Please run 'make prepare' first"; \
    fi
	@swag init -g main.go --parseDependency --parseInternal --dir ./cmd/server,./internal/application/dto,./internal/api -o docs

force-download:
	@echo "Cleaning Go module cache..."
	@go clean -modcache

	@echo "Disabling SSL verification for Git temporarily..."
	@git config --global http.sslVerify false

	@echo "Configuring Go environment for insecure download..."
	@go env -w GOINSECURE=github.com,go.googlesource.com,golang.org,go.uber.org,google.golang.org,sigs.k8s.io,rsc.io
	@go env -w GOSUMDB=off
	@go env -w GOPROXY=direct

	@echo "Done! Go environment variables updated."