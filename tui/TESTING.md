# MCF TUI Testing Guide

This document provides comprehensive information about testing the MCF TUI application.

## Table of Contents

1. [Testing Overview](#testing-overview)
2. [Test Structure](#test-structure)
3. [Running Tests](#running-tests)
4. [Test Types](#test-types)
5. [Mock Usage](#mock-usage)
6. [Test Data Management](#test-data-management)
7. [Debugging Tests](#debugging-tests)
8. [CI/CD Pipeline](#cicd-pipeline)
9. [Best Practices](#best-practices)
10. [Troubleshooting](#troubleshooting)

## Testing Overview

The MCF TUI project employs a comprehensive testing strategy that includes:

- **Unit Tests**: Fast, isolated tests for individual components
- **Integration Tests**: Tests for component interactions and MCF integrations
- **End-to-End Tests**: Complete workflow tests simulating user interactions
- **Performance Tests**: Benchmarks and performance validation
- **Security Tests**: Vulnerability scanning and security validation

### Coverage Goals

- **Overall Coverage**: ≥ 80%
- **Critical Path Coverage**: ≥ 95%
- **Business Logic Coverage**: ≥ 90%

## Test Structure

```
tui/
├── internal/
│   ├── app/
│   │   ├── model_test.go          # App model unit tests
│   │   └── update_test.go         # Update logic unit tests
│   ├── commands/
│   │   └── adapter_test.go        # MCF command adapter tests
│   ├── orchestration/
│   │   └── orchestrator_test.go   # Agent orchestration tests
│   ├── config/
│   │   └── manager_test.go        # Configuration management tests
│   ├── e2e/
│   │   └── workflow_test.go       # End-to-end workflow tests
│   └── testing/
│       ├── utils.go               # Test utilities and helpers
│       └── mocks.go               # Mock implementations
├── testdata/                      # Test data and fixtures
├── Makefile                       # Test automation
└── TESTING.md                     # This document
```

## Running Tests

### Prerequisites

1. **Go 1.21+**: Required for running tests
2. **Make**: For test automation scripts
3. **Testing Tools**: Installed via `make install-tools`

```bash
# Install required tools
make install-tools
```

### Basic Commands

```bash
# Run all tests
make test

# Run specific test types
make test-unit           # Unit tests only
make test-integration    # Integration tests only
make test-e2e           # End-to-end tests only
make test-performance   # Performance tests only

# Run tests with coverage
make test-coverage

# Run tests with race detection
make test-race

# Quick development tests (unit only, no race detection)
make dev-test
```

### Advanced Testing

```bash
# Run tests with verbose output
make test-verbose

# Watch for changes and run tests automatically
make test-watch

# Run tests in Docker
make test-docker

# Generate test report
make test-report
```

## Test Types

### 1. Unit Tests

Unit tests focus on individual components in isolation.

**Example: Testing the MCF Model**

```go
func TestInitialModel(t *testing.T) {
    t.Run("should create model with default values", func(t *testing.T) {
        model := InitialModel()

        assert.Equal(t, DashboardView, model.currentView)
        assert.False(t, model.ready)
        assert.Equal(t, 0, model.width)
        assert.Equal(t, 0, model.height)
    })
}
```

**Running Unit Tests:**

```bash
# Run all unit tests
make test-unit

# Run specific package tests
go test -v ./internal/app

# Run specific test function
go test -v ./internal/app -run TestInitialModel
```

### 2. Integration Tests

Integration tests verify component interactions and MCF system integration.

**Example: Testing MCF Command Adapter**

```go
func TestMCFCommandAdapter_ExecuteCommand_Success(t *testing.T) {
    // Arrange
    mockClient := testutils.NewMockMCFClient()
    adapter := NewMCFCommandAdapter(mockClient, logger)

    command := "status"
    args := []string{"--verbose"}
    expectedResult := testutils.CommandResult{
        Command:  command,
        Args:     args,
        ExitCode: 0,
        Output:   "Service status: OK",
    }

    mockClient.On("ExecuteCommand", ctx, command, args).Return(expectedResult, nil)

    // Act
    result, err := adapter.ExecuteCommand(ctx, command, args)

    // Assert
    assert.NoError(t, err)
    assert.Equal(t, command, result.Command)
    assert.Equal(t, 0, result.ExitCode)
}
```

**Running Integration Tests:**

```bash
# Run all integration tests
make test-integration

# Run with integration tag
go test -v -tags=integration ./internal/commands
```

### 3. End-to-End Tests

E2E tests simulate complete user workflows using the Bubble Tea test runner.

**Example: Complete Application Workflow**

```go
func TestCompleteApplicationWorkflow(t *testing.T) {
    // Setup
    model := app.InitialModel()
    runner := testutils.NewTestProgramRunner(model)

    runner.Start(t)
    defer runner.Stop()

    // Simulate user interaction
    runner.SendKeys(":", "h", "e", "l", "p", "\r") // :help command
    runner.WaitForOutput(t, "Help", 2*time.Second)

    runner.SendKeys("\x1b", "q") // Escape then quit
}
```

**Running E2E Tests:**

```bash
# Run all e2e tests
make test-e2e

# Run with e2e tag
go test -v -tags=e2e ./internal/e2e
```

### 4. Performance Tests

Performance tests include benchmarks and load testing.

**Example: Performance Benchmark**

```go
func TestMCFModel_Performance(t *testing.T) {
    model := InitialModel()

    benchmark := testutils.NewPerformanceBenchmark("rapid_view_switching", func() error {
        views := []View{DashboardView, LogView, CommandBarView}
        for _, view := range views {
            model.SetView(view)
        }
        return nil
    }).WithIterations(1000)

    benchmark.Run(t)
}
```

**Running Performance Tests:**

```bash
# Run performance tests
make test-performance

# Run benchmarks
make benchmark

# Compare benchmark results
make benchmark-compare
```

## Mock Usage

The testing framework provides comprehensive mocking capabilities for MCF integrations.

### MCF Client Mocking

```go
// Create mock client
mockClient := testutils.NewMockMCFClient()

// Setup system health mock data
healthStatus := testutils.SystemHealthStatus{
    Status:  "healthy",
    Version: "1.0.0-test",
    Uptime:  time.Hour * 24,
}
mockClient.SetSystemHealth(healthStatus)

// Setup expectations
mockClient.On("GetSystemHealth", ctx).Return(healthStatus, nil)

// Use in tests
status, err := mockClient.GetSystemHealth(ctx)
assert.NoError(t, err)
assert.Equal(t, "healthy", status.Status)

// Verify expectations were met
mockClient.AssertExpectations(t)
```

### Agent Orchestration Mocking

```go
// Create mock orchestrator
orchestrator := testutils.NewMockAgentOrchestrator()

// Register mock agents
agent := &testutils.MockAgent{
    ID:           "test-agent",
    Capabilities: []string{"analysis", "reporting"},
    Status:       "active",
}
orchestrator.RegisterAgent(agent)

// Setup task submission expectations
task := testutils.Task{
    ID:   "test-task",
    Type: "analysis",
}
orchestrator.On("SubmitTask", ctx, task).Return(nil)
```

### Configuration Mocking

```go
// Create mock config manager
configManager := testutils.NewMockConfigManager()

// Setup configuration expectations
configManager.On("Get", "tui.theme").Return("dark", true)
configManager.On("Set", "tui.theme", "light").Return(nil)

// Use in tests
theme, exists := configManager.Get("tui.theme")
assert.True(t, exists)
assert.Equal(t, "dark", theme)
```

### Test Utilities

The testing package provides various utilities:

```go
// Test program runner for Bubble Tea applications
runner := testutils.NewTestProgramRunner(model)
runner.Start(t)
defer runner.Stop()

// Send keyboard input
runner.SendKeys(":", "h", "e", "l", "p", "\r")
runner.WaitForOutput(t, "expected text", 2*time.Second)

// Performance benchmarking
benchmark := testutils.NewPerformanceBenchmark("operation_name", operation)
benchmark.WithIterations(1000).WithTimeout(30*time.Second)
benchmark.Run(t)

// Test data builder
data := testutils.NewTestDataBuilder()
    .Set("key1", "value1")
    .Set("key2", 123)
    .Build()

// Mock logger
logger := testutils.NewTestLogger(t)
logger.Log("Test message: %s", "data")
logger.Error("Error occurred: %v", err)
```

## Test Data Management

### Test Data Structure

```
testdata/
├── configs/              # Configuration files
│   ├── default.json
│   ├── test.json
│   └── invalid.json
├── fixtures/             # Test fixtures
│   ├── agents.json
│   ├── services.json
│   └── logs.json
├── mocks/               # Mock data files
└── temp/                # Temporary test files
```

### Creating Test Data

```go
// Generate test data
func generateTestServices() []testutils.ServiceStatus {
    return []testutils.ServiceStatus{
        {
            Name:    "api-server",
            Status:  "running",
            Health:  "healthy",
            Port:    8080,
        },
        {
            Name:    "database",
            Status:  "running",
            Health:  "healthy",
            Port:    5432,
        },
    }
}

// Use test data
services := generateTestServices()
mockClient.SetServices(services)
```

### Test Data Factories

```go
// User factory
func createTestUser(overrides ...map[string]interface{}) TestUser {
    user := TestUser{
        ID:   "test-user-" + generateID(),
        Name: "Test User",
        Role: "user",
    }

    // Apply overrides
    if len(overrides) > 0 {
        // Apply override values
    }

    return user
}

// Usage
user := createTestUser(map[string]interface{}{
    "Name": "Admin User",
    "Role": "admin",
})
```

## Debugging Tests

### Debug Configuration

```bash
# Run tests with debug output
go test -v -debug ./...

# Run specific test with debug
go test -v ./internal/app -run TestSpecificFunction -debug

# Run with race detector and debug
go test -v -race -debug ./...
```

### Debug Environment Variables

```bash
# Set debug level
export MCF_TUI_DEBUG=true
export MCF_TUI_LOG_LEVEL=debug

# Enable test debug output
export MCF_TUI_TEST_DEBUG=true
```

### Debugging Test Failures

```go
func TestDebugExample(t *testing.T) {
    // Use test logger for debug output
    logger := testutils.NewTestLogger(t)
    logger.Log("Starting test with data: %+v", testData)

    // Add debug assertions
    assert.NotNil(t, result, "Result should not be nil - debug info: %+v", debugInfo)

    // Use require for early failure
    require.NoError(t, err, "Critical error occurred: %v", err)

    // Debug test state
    t.Logf("Test state: %+v", currentState)
}
```

### Common Debug Techniques

1. **Test Isolation**: Run single test to isolate issues
2. **Verbose Output**: Use `-v` flag for detailed output
3. **Test Logging**: Add strategic log statements
4. **State Inspection**: Dump object state at failure points
5. **Mock Verification**: Verify mock expectations are correct

## CI/CD Pipeline

The CI/CD pipeline runs comprehensive tests on every push and pull request.

### Pipeline Stages

1. **Validation**: Code formatting, linting, and basic checks
2. **Unit Tests**: Fast isolated tests across multiple Go versions and OS
3. **Integration Tests**: Component interaction tests
4. **E2E Tests**: Complete workflow validation
5. **Coverage**: Code coverage analysis and reporting
6. **Race Detection**: Concurrent execution safety tests
7. **Security**: Vulnerability and security scanning
8. **Performance**: Benchmark execution and comparison
9. **Quality Gate**: Overall test result validation

### Pipeline Configuration

The pipeline is defined in `.github/workflows/test.yml` and includes:

- **Matrix Testing**: Multiple Go versions (1.21, 1.22) and OS (Linux, macOS, Windows)
- **Parallel Execution**: Tests run in parallel for faster feedback
- **Caching**: Go module caching for improved performance
- **Artifacts**: Test reports and coverage data
- **Quality Gates**: Enforced coverage thresholds and test success

### Local CI Simulation

```bash
# Run CI tests locally
make test-ci

# Run fast CI tests (unit only)
make test-ci-fast

# Run tests in Docker (simulates CI environment)
make test-docker
```

## Best Practices

### Test Organization

1. **Naming Convention**: Use descriptive test names

   ```go
   func TestMCFModel_Update_ShouldHandleWindowSizeMessage(t *testing.T)
   ```

2. **Test Structure**: Follow Arrange-Act-Assert pattern

   ```go
   func TestExample(t *testing.T) {
       // Arrange
       model := InitialModel()
       expectedValue := "expected"

       // Act
       result := model.SomeMethod()

       // Assert
       assert.Equal(t, expectedValue, result)
   }
   ```

3. **Test Isolation**: Each test should be independent
   ```go
   func TestSuite(t *testing.T) {
       t.Run("test case 1", func(t *testing.T) {
           // Isolated test logic
       })

       t.Run("test case 2", func(t *testing.T) {
           // Isolated test logic
       })
   }
   ```

### Mock Best Practices

1. **Minimal Mocking**: Only mock external dependencies
2. **Clear Expectations**: Set explicit mock expectations
3. **Verification**: Always verify mock expectations
4. **Realistic Data**: Use realistic mock data

### Performance Testing

1. **Baseline Establishment**: Establish performance baselines
2. **Regression Detection**: Monitor for performance regressions
3. **Resource Monitoring**: Monitor memory and CPU usage
4. **Load Testing**: Test under realistic load conditions

### Error Testing

1. **Error Conditions**: Test both success and error paths
2. **Edge Cases**: Test boundary conditions and edge cases
3. **Resource Limits**: Test resource exhaustion scenarios
4. **Recovery**: Test error recovery mechanisms

## Troubleshooting

### Common Issues

#### Test Timeouts

```bash
# Increase test timeout
go test -timeout=30m ./...

# Or use Makefile with extended timeout
make test-verbose
```

#### Race Conditions

```bash
# Run with race detector
make test-race

# Debug race conditions
go test -race -v ./internal/app -run TestSpecificFunction
```

#### Mock Failures

```go
// Check mock expectations
mockClient.AssertExpectations(t)

// Debug mock calls
mockClient.AssertNumberOfCalls(t, "GetSystemHealth", 1)
mockClient.AssertCalled(t, "GetSystemHealth", mock.AnythingOfType("context.Context"))
```

#### Coverage Issues

```bash
# Generate detailed coverage report
make test-coverage

# View coverage by package
go tool cover -func=coverage.out

# View HTML coverage report
go tool cover -html=coverage.out
```

#### Performance Regression

```bash
# Compare benchmarks
make benchmark-compare

# Profile specific tests
go test -cpuprofile=cpu.prof -bench=BenchmarkSpecific ./...
go tool pprof cpu.prof
```

### Getting Help

1. **Check Logs**: Review test output and logs
2. **Run Individually**: Isolate failing tests
3. **Debug Mode**: Use debug flags and logging
4. **Mock Verification**: Verify mock setup and expectations
5. **Documentation**: Refer to this guide and code comments

### Environment Issues

```bash
# Clean test environment
make clean-test-env

# Reset test cache
go clean -testcache

# Reinstall tools
make install-tools
```

## Continuous Improvement

### Test Metrics

Monitor these key testing metrics:

- **Coverage Percentage**: Maintain ≥80% overall coverage
- **Test Execution Time**: Keep tests fast and efficient
- **Flaky Test Rate**: Minimize test flakiness
- **Test Maintenance Overhead**: Keep tests maintainable

### Regular Tasks

1. **Review Test Coverage**: Weekly coverage analysis
2. **Update Test Data**: Keep test data current
3. **Performance Baselines**: Update performance baselines
4. **Mock Maintenance**: Keep mocks synchronized with real APIs
5. **Documentation**: Update test documentation

---

For questions or issues with testing, please refer to the project documentation or create an issue in the repository.
