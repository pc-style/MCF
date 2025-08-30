# MCF TUI Testing Makefile
# Provides comprehensive testing automation for the MCF TUI project

.PHONY: help test test-unit test-integration test-e2e test-performance test-coverage
.PHONY: test-race test-verbose test-watch test-clean benchmark lint fmt vet
.PHONY: test-report test-ci setup-test-env clean-test-env mock-generate

# Default target
help: ## Show this help message
	@echo "MCF TUI Testing Automation"
	@echo "=========================="
	@echo ""
	@echo "Available targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-20s %s\n", $$1, $$2}'

# Go binary and flags
GO := go
GO_VERSION := 1.21
GO_TEST_FLAGS := -v -race -timeout=10m
GO_BUILD_FLAGS := -ldflags="-s -w"

# Test directories
TEST_DIRS := ./internal/app ./internal/commands ./internal/orchestration ./internal/config ./internal/e2e ./internal/testing
UNIT_TEST_DIRS := ./internal/app ./internal/commands ./internal/config ./internal/testing
INTEGRATION_TEST_DIRS := ./internal/orchestration ./internal/commands
E2E_TEST_DIRS := ./internal/e2e

# Coverage settings
COVERAGE_OUT := coverage.out
COVERAGE_HTML := coverage.html
COVERAGE_THRESHOLD := 80

# Test environment
TEST_CONFIG_FILE := test-config.json
TEST_LOG_LEVEL := debug

# Setup and cleanup
setup-test-env: ## Setup test environment
	@echo "Setting up test environment..."
	@mkdir -p ./testdata/configs
	@mkdir -p ./testdata/logs
	@mkdir -p ./testdata/temp
	@echo '{"test": true, "mcf": {"host": "localhost", "port": 8080}}' > ./testdata/configs/$(TEST_CONFIG_FILE)
	@echo "Test environment ready"

clean-test-env: ## Clean up test environment
	@echo "Cleaning test environment..."
	@rm -rf ./testdata
	@rm -f $(COVERAGE_OUT) $(COVERAGE_HTML)
	@rm -f *.test
	@echo "Test environment cleaned"

# Code generation
mock-generate: ## Generate mock files
	@echo "Generating mocks..."
	@$(GO) install github.com/golang/mock/mockgen@latest
	@$(GO) generate ./...
	@echo "Mocks generated"

# Linting and formatting
fmt: ## Format Go code
	@echo "Formatting Go code..."
	@$(GO) fmt ./...
	@echo "Code formatted"

vet: ## Run go vet
	@echo "Running go vet..."
	@$(GO) vet ./...
	@echo "Vet completed"

lint: ## Run golangci-lint
	@echo "Running golangci-lint..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run ./...; \
	else \
		echo "golangci-lint not found. Installing..."; \
		$(GO) install github.com/golangci/golangci-lint/cmd/golangci-lint@latest; \
		golangci-lint run ./...; \
	fi
	@echo "Linting completed"

# Testing targets
test: setup-test-env test-unit test-integration test-e2e ## Run all tests
	@echo "All tests completed successfully"

test-unit: ## Run unit tests
	@echo "Running unit tests..."
	@$(GO) test $(GO_TEST_FLAGS) $(UNIT_TEST_DIRS)
	@echo "Unit tests completed"

test-integration: ## Run integration tests
	@echo "Running integration tests..."
	@$(GO) test $(GO_TEST_FLAGS) -tags=integration $(INTEGRATION_TEST_DIRS)
	@echo "Integration tests completed"

test-e2e: ## Run end-to-end tests
	@echo "Running end-to-end tests..."
	@$(GO) test $(GO_TEST_FLAGS) -tags=e2e $(E2E_TEST_DIRS)
	@echo "End-to-end tests completed"

test-performance: ## Run performance tests
	@echo "Running performance tests..."
	@$(GO) test $(GO_TEST_FLAGS) -tags=performance -run=TestPerformance ./...
	@echo "Performance tests completed"

test-race: ## Run tests with race detection
	@echo "Running tests with race detection..."
	@$(GO) test -race -timeout=15m $(TEST_DIRS)
	@echo "Race detection tests completed"

test-verbose: ## Run tests with verbose output
	@echo "Running tests with verbose output..."
	@$(GO) test -v -timeout=15m $(TEST_DIRS)
	@echo "Verbose tests completed"

test-watch: ## Watch for changes and run tests
	@echo "Watching for changes and running tests..."
	@if command -v watchexec >/dev/null 2>&1; then \
		watchexec -e go -- make test-unit; \
	elif command -v fswatch >/dev/null 2>&1; then \
		fswatch -o . | xargs -n1 -I{} make test-unit; \
	else \
		echo "No file watcher found. Install watchexec or fswatch for watch functionality."; \
		echo "Falling back to single test run..."; \
		make test-unit; \
	fi

test-clean: clean-test-env ## Clean test artifacts and run fresh tests
	@$(GO) clean -testcache
	@make test

# Coverage targets
test-coverage: setup-test-env ## Run tests with coverage
	@echo "Running tests with coverage..."
	@$(GO) test -coverprofile=$(COVERAGE_OUT) -covermode=atomic $(TEST_DIRS)
	@$(GO) tool cover -html=$(COVERAGE_OUT) -o $(COVERAGE_HTML)
	@echo "Coverage report generated: $(COVERAGE_HTML)"
	@$(GO) tool cover -func=$(COVERAGE_OUT) | grep total | awk '{print "Total coverage: " $$3}'

test-coverage-threshold: test-coverage ## Check if coverage meets threshold
	@echo "Checking coverage threshold ($(COVERAGE_THRESHOLD)%)..."
	@COVERAGE=$$($(GO) tool cover -func=$(COVERAGE_OUT) | grep total | awk '{print $$3}' | sed 's/%//'); \
	if (( $$(echo "$$COVERAGE >= $(COVERAGE_THRESHOLD)" | bc -l) )); then \
		echo "✓ Coverage $$COVERAGE% meets threshold $(COVERAGE_THRESHOLD)%"; \
	else \
		echo "✗ Coverage $$COVERAGE% below threshold $(COVERAGE_THRESHOLD)%"; \
		exit 1; \
	fi

# Benchmark targets
benchmark: ## Run benchmark tests
	@echo "Running benchmark tests..."
	@$(GO) test -bench=. -benchmem -timeout=30m $(TEST_DIRS)
	@echo "Benchmarks completed"

benchmark-compare: ## Run benchmarks and compare with previous results
	@echo "Running benchmark comparison..."
	@if [ -f bench-old.txt ]; then \
		$(GO) test -bench=. -benchmem $(TEST_DIRS) > bench-new.txt; \
		if command -v benchcmp >/dev/null 2>&1; then \
			benchcmp bench-old.txt bench-new.txt; \
		else \
			echo "benchcmp not found. Install with: go install golang.org/x/tools/cmd/benchcmp@latest"; \
			echo "New benchmark results:"; \
			cat bench-new.txt; \
		fi; \
		mv bench-new.txt bench-old.txt; \
	else \
		$(GO) test -bench=. -benchmem $(TEST_DIRS) > bench-old.txt; \
		echo "Baseline benchmark results saved to bench-old.txt"; \
	fi

# Test reporting
test-report: ## Generate comprehensive test report
	@echo "Generating comprehensive test report..."
	@mkdir -p ./reports
	@echo "# MCF TUI Test Report" > ./reports/test-report.md
	@echo "" >> ./reports/test-report.md
	@echo "Generated at: $$(date)" >> ./reports/test-report.md
	@echo "" >> ./reports/test-report.md
	
	@echo "## Test Coverage" >> ./reports/test-report.md
	@$(GO) test -coverprofile=$(COVERAGE_OUT) $(TEST_DIRS) > /dev/null 2>&1 || true
	@if [ -f $(COVERAGE_OUT) ]; then \
		$(GO) tool cover -func=$(COVERAGE_OUT) | tail -1 >> ./reports/test-report.md; \
	fi
	@echo "" >> ./reports/test-report.md
	
	@echo "## Test Results" >> ./reports/test-report.md
	@$(GO) test -json $(TEST_DIRS) 2>/dev/null | jq -r '. | select(.Action=="pass" or .Action=="fail") | "\(.Action | ascii_upcase): \(.Package) - \(.Test // "Package")"' >> ./reports/test-report.md 2>/dev/null || echo "JSON parsing not available" >> ./reports/test-report.md
	
	@echo "Test report generated: ./reports/test-report.md"

# CI/CD specific targets
test-ci: ## Run tests for CI environment
	@echo "Running CI tests..."
	@$(GO) version
	@make lint
	@make test-coverage-threshold
	@make test-race
	@echo "CI tests completed successfully"

test-ci-fast: ## Run fast CI tests (unit only)
	@echo "Running fast CI tests..."
	@make fmt vet
	@make test-unit
	@echo "Fast CI tests completed"

# Docker test targets
test-docker: ## Run tests in Docker container
	@echo "Running tests in Docker..."
	@docker run --rm -v $(PWD):/workspace -w /workspace golang:$(GO_VERSION) make test-ci
	@echo "Docker tests completed"

test-docker-build: ## Build Docker test image
	@echo "Building Docker test image..."
	@docker build -f Dockerfile.test -t mcf-tui-test .
	@echo "Docker test image built"

# Security and vulnerability testing
test-security: ## Run security tests
	@echo "Running security tests..."
	@if command -v gosec >/dev/null 2>&1; then \
		gosec ./...; \
	else \
		echo "gosec not found. Installing..."; \
		curl -sfL https://raw.githubusercontent.com/securego/gosec/master/install.sh | sh -s -- -b "$$(go env GOPATH)/bin" latest; \
		gosec ./...; \
	fi
	@echo "Security tests completed"

test-vuln: ## Check for known vulnerabilities
	@echo "Checking for vulnerabilities..."
	@$(GO) install golang.org/x/vuln/cmd/govulncheck@latest
	@govulncheck ./...
	@echo "Vulnerability check completed"
	@echo "Checking for vulnerabilities..."
	@$(GO) install golang.org/x/vuln/cmd/govulncheck@latest
	@govulncheck ./...
	@echo "Vulnerability check completed"

# Performance profiling
profile-cpu: ## Run CPU profiling
	@echo "Running CPU profiling..."
	@$(GO) test -cpuprofile=cpu.prof -bench=. $(TEST_DIRS)
	@echo "CPU profile saved to cpu.prof"
	@echo "View with: go tool pprof cpu.prof"

profile-mem: ## Run memory profiling
	@echo "Running memory profiling..."
	@$(GO) test -memprofile=mem.prof -bench=. $(TEST_DIRS)
	@echo "Memory profile saved to mem.prof"
	@echo "View with: go tool pprof mem.prof"

profile-block: ## Run blocking profiling
	@echo "Running blocking profiling..."
	@$(GO) test -blockprofile=block.prof -bench=. $(TEST_DIRS)
	@echo "Blocking profile saved to block.prof"
	@echo "View with: go tool pprof block.prof"

# Test data management
test-data-generate: ## Generate test data
	@echo "Generating test data..."
	@mkdir -p ./testdata/fixtures
	@$(GO) run ./scripts/generate-test-data.go
	@echo "Test data generated in ./testdata/fixtures"

test-data-clean: ## Clean test data
	@echo "Cleaning test data..."
	@rm -rf ./testdata/fixtures/*
	@echo "Test data cleaned"

# Quality gates
quality-gate: ## Run quality gate checks
	@echo "Running quality gate checks..."
	@make lint
	@make test-coverage-threshold
	@make test-security
	@make test-vuln
	@echo "✓ Quality gate passed"

# Development helpers
dev-test: ## Quick development tests (unit tests only, no race)
	@echo "Running development tests..."
	@$(GO) test -timeout=5m $(UNIT_TEST_DIRS)
	@echo "Development tests completed"

dev-watch: ## Development mode with file watching
	@echo "Starting development watch mode..."
	@if command -v watchexec >/dev/null 2>&1; then \
		watchexec -e go -i testdata -- make dev-test; \
	else \
		echo "watchexec not found. Please install it for watch functionality."; \
		make dev-test; \
	fi

# Installation and setup
install-tools: ## Install required testing tools
	@echo "Installing testing tools..."
	@$(GO) install github.com/golang/mock/mockgen@latest
	@$(GO) install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@$(GO) install golang.org/x/vuln/cmd/govulncheck@latest
	@$(GO) install golang.org/x/tools/cmd/benchcmp@latest
	@echo "Testing tools installed"

# Information
info: ## Show testing information
	@echo "MCF TUI Testing Information"
	@echo "=========================="
	@echo "Go version: $$($(GO) version)"
	@echo "Test directories: $(TEST_DIRS)"
	@echo "Coverage threshold: $(COVERAGE_THRESHOLD)%"
	@echo "Test config: $(TEST_CONFIG_FILE)"
	@echo ""
	@echo "Available test types:"
	@echo "  - Unit tests (fast, isolated)"
	@echo "  - Integration tests (component interaction)"
	@echo "  - E2E tests (complete workflows)"
	@echo "  - Performance tests (benchmarks and load)"
	@echo "  - Race condition detection"
	@echo "  - Security vulnerability scanning"
	@echo ""
	@echo "Use 'make help' to see all available targets"