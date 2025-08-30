package testing

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"strings"
	"sync"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestProgramRunner provides utilities for testing Bubble Tea programs
type TestProgramRunner struct {
	program *tea.Program
	output  *bytes.Buffer
	input   *bytes.Buffer
	done    chan struct{}
	mu      sync.Mutex
}

// NewTestProgramRunner creates a new test runner for Bubble Tea programs
func NewTestProgramRunner(model tea.Model) *TestProgramRunner {
	output := &bytes.Buffer{}
	input := &bytes.Buffer{}

	program := tea.NewProgram(
		model,
		tea.WithInput(input),
		tea.WithOutput(output),
		tea.WithoutSignalHandler(),
	)

	return &TestProgramRunner{
		program: program,
		output:  output,
		input:   input,
		done:    make(chan struct{}),
	}
}

// Start runs the program in a separate goroutine
func (r *TestProgramRunner) Start(t *testing.T) {
	go func() {
		defer close(r.done)
		if err := r.program.Start(); err != nil {
			t.Errorf("Program failed to start: %v", err)
		}
	}()
	time.Sleep(50 * time.Millisecond) // Allow program to initialize
}

// Stop gracefully stops the program
func (r *TestProgramRunner) Stop() {
	r.program.Quit()
	<-r.done
}

// SendKeypress simulates a key press
func (r *TestProgramRunner) SendKeypress(key string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.input.WriteString(key)
}

// SendKeys sends multiple keypresses
func (r *TestProgramRunner) SendKeys(keys ...string) {
	for _, key := range keys {
		r.SendKeypress(key)
		time.Sleep(10 * time.Millisecond)
	}
}

// GetOutput returns the current program output
func (r *TestProgramRunner) GetOutput() string {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.output.String()
}

// WaitForOutput waits for specific text to appear in output
func (r *TestProgramRunner) WaitForOutput(t *testing.T, expected string, timeout time.Duration) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	ticker := time.NewTicker(10 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			t.Fatalf("Timeout waiting for output: %s\nActual output: %s", expected, r.GetOutput())
		case <-ticker.C:
			if strings.Contains(r.GetOutput(), expected) {
				return
			}
		}
	}
}

// AssertOutputContains checks if output contains expected text
func (r *TestProgramRunner) AssertOutputContains(t *testing.T, expected string) {
	output := r.GetOutput()
	assert.Contains(t, output, expected, "Output should contain expected text")
}

// AssertOutputNotContains checks if output does not contain text
func (r *TestProgramRunner) AssertOutputNotContains(t *testing.T, notExpected string) {
	output := r.GetOutput()
	assert.NotContains(t, output, notExpected, "Output should not contain text")
}

// MockTicker provides a controllable ticker for testing time-based features
type MockTicker struct {
	ch      chan time.Time
	stopped bool
	mu      sync.Mutex
}

// NewMockTicker creates a new mock ticker
func NewMockTicker() *MockTicker {
	return &MockTicker{
		ch: make(chan time.Time, 1),
	}
}

// C returns the ticker channel
func (t *MockTicker) C() <-chan time.Time {
	return t.ch
}

// Tick sends a tick event
func (t *MockTicker) Tick() {
	t.mu.Lock()
	defer t.mu.Unlock()
	if !t.stopped {
		select {
		case t.ch <- time.Now():
		default:
		}
	}
}

// Stop stops the ticker
func (t *MockTicker) Stop() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.stopped = true
	close(t.ch)
}

// TestDataBuilder helps build test data structures
type TestDataBuilder struct {
	data map[string]interface{}
}

// NewTestDataBuilder creates a new test data builder
func NewTestDataBuilder() *TestDataBuilder {
	return &TestDataBuilder{
		data: make(map[string]interface{}),
	}
}

// Set adds a key-value pair to the test data
func (b *TestDataBuilder) Set(key string, value interface{}) *TestDataBuilder {
	b.data[key] = value
	return b
}

// Get retrieves a value from the test data
func (b *TestDataBuilder) Get(key string) interface{} {
	return b.data[key]
}

// GetString retrieves a string value
func (b *TestDataBuilder) GetString(key string) string {
	if val, ok := b.data[key].(string); ok {
		return val
	}
	return ""
}

// GetInt retrieves an int value
func (b *TestDataBuilder) GetInt(key string) int {
	if val, ok := b.data[key].(int); ok {
		return val
	}
	return 0
}

// Build returns the built data
func (b *TestDataBuilder) Build() map[string]interface{} {
	result := make(map[string]interface{})
	for k, v := range b.data {
		result[k] = v
	}
	return result
}

// AssertModelState provides utilities for testing model state
func AssertModelState(t *testing.T, model interface{}, checks ...StateCheck) {
	for _, check := range checks {
		check(t, model)
	}
}

// StateCheck is a function that validates model state
type StateCheck func(t *testing.T, model interface{})

// CheckField validates a specific field value
func CheckField(fieldName string, expected interface{}) StateCheck {
	return func(t *testing.T, model interface{}) {
		// Use reflection or type assertion to check field values
		// This is a simplified version - in practice, you'd use reflection
		switch m := model.(type) {
		case map[string]interface{}:
			actual, exists := m[fieldName]
			require.True(t, exists, "Field %s should exist", fieldName)
			assert.Equal(t, expected, actual, "Field %s should match expected value", fieldName)
		default:
			t.Errorf("Unsupported model type for field checking: %T", model)
		}
	}
}

// CheckNotNil validates that a field is not nil
func CheckNotNil(fieldName string) StateCheck {
	return func(t *testing.T, model interface{}) {
		switch m := model.(type) {
		case map[string]interface{}:
			actual, exists := m[fieldName]
			require.True(t, exists, "Field %s should exist", fieldName)
			assert.NotNil(t, actual, "Field %s should not be nil", fieldName)
		default:
			t.Errorf("Unsupported model type for nil checking: %T", model)
		}
	}
}

// MockWriter provides a writer for testing output
type MockWriter struct {
	data []byte
	mu   sync.Mutex
}

// NewMockWriter creates a new mock writer
func NewMockWriter() *MockWriter {
	return &MockWriter{}
}

// Write implements the io.Writer interface
func (w *MockWriter) Write(p []byte) (n int, err error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.data = append(w.data, p...)
	return len(p), nil
}

// String returns the written data as a string
func (w *MockWriter) String() string {
	w.mu.Lock()
	defer w.mu.Unlock()
	return string(w.data)
}

// Clear clears the written data
func (w *MockWriter) Clear() {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.data = nil
}

// TestLogger provides a test-safe logger
type TestLogger struct {
	t      *testing.T
	output io.Writer
}

// NewTestLogger creates a new test logger
func NewTestLogger(t *testing.T) *TestLogger {
	return &TestLogger{
		t:      t,
		output: NewMockWriter(),
	}
}

// Log logs a message
func (l *TestLogger) Log(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	l.t.Logf("TestLogger: %s", msg)
	if l.output != nil {
		fmt.Fprintf(l.output, "%s\n", msg)
	}
}

// Error logs an error message
func (l *TestLogger) Error(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	l.t.Errorf("TestLogger Error: %s", msg)
	if l.output != nil {
		fmt.Fprintf(l.output, "ERROR: %s\n", msg)
	}
}

// SetOutput sets the output writer
func (l *TestLogger) SetOutput(w io.Writer) {
	l.output = w
}

// Performance testing utilities
type PerformanceBenchmark struct {
	name       string
	operation  func() error
	iterations int
	timeout    time.Duration
}

// NewPerformanceBenchmark creates a new performance benchmark
func NewPerformanceBenchmark(name string, operation func() error) *PerformanceBenchmark {
	return &PerformanceBenchmark{
		name:       name,
		operation:  operation,
		iterations: 1000,
		timeout:    10 * time.Second,
	}
}

// WithIterations sets the number of iterations
func (b *PerformanceBenchmark) WithIterations(iterations int) *PerformanceBenchmark {
	b.iterations = iterations
	return b
}

// WithTimeout sets the timeout
func (b *PerformanceBenchmark) WithTimeout(timeout time.Duration) *PerformanceBenchmark {
	b.timeout = timeout
	return b
}

// Run executes the benchmark
func (b *PerformanceBenchmark) Run(t *testing.T) {
	start := time.Now()

	for i := 0; i < b.iterations; i++ {
		if err := b.operation(); err != nil {
			t.Fatalf("Benchmark %s failed at iteration %d: %v", b.name, i, err)
		}

		if time.Since(start) > b.timeout {
			t.Fatalf("Benchmark %s timed out after %d iterations", b.name, i)
		}
	}

	duration := time.Since(start)
	avgDuration := duration / time.Duration(b.iterations)

	t.Logf("Benchmark %s completed: %d iterations in %v (avg: %v per iteration)",
		b.name, b.iterations, duration, avgDuration)
}
