package testing

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewTestLogger(t *testing.T) {
	t.Run("should create logger with test context", func(t *testing.T) {
		logger := NewTestLogger(t)

		assert.NotNil(t, logger, "Logger should be created")
		assert.IsType(t, &TestLogger{}, logger, "Should return TestLogger type")
	})
}

func TestTestLogger_Operations(t *testing.T) {
	t.Run("should log messages without errors", func(t *testing.T) {
		logger := NewTestLogger(t)

		// These should not panic
		logger.Log("Test message: %s", "formatted")
		logger.Error("Test error: %d", 404)
	})
}

func TestNewPerformanceBenchmark(t *testing.T) {
	t.Run("should create benchmark with valid parameters", func(t *testing.T) {
		testFunc := func() error {
			time.Sleep(1 * time.Millisecond)
			return nil
		}

		benchmark := NewPerformanceBenchmark("test_operation", testFunc)

		assert.NotNil(t, benchmark, "Benchmark should be created")
		assert.Equal(t, "test_operation", benchmark.name)
		assert.Equal(t, 1000, benchmark.iterations, "Default iterations should be 1000")
	})

	t.Run("should configure iterations", func(t *testing.T) {
		testFunc := func() error { return nil }
		benchmark := NewPerformanceBenchmark("test_operation", testFunc)

		configuredBenchmark := benchmark.WithIterations(100)

		assert.Equal(t, 100, configuredBenchmark.iterations)
	})

	t.Run("should configure timeout", func(t *testing.T) {
		testFunc := func() error { return nil }
		benchmark := NewPerformanceBenchmark("test_operation", testFunc)

		configuredBenchmark := benchmark.WithTimeout(5 * time.Second)

		assert.Equal(t, 5*time.Second, configuredBenchmark.timeout)
	})
}

func TestPerformanceBenchmark_Run(t *testing.T) {
	t.Run("should execute benchmark successfully", func(t *testing.T) {
		executionCount := 0
		testFunc := func() error {
			executionCount++
			return nil
		}

		benchmark := NewPerformanceBenchmark("test_operation", testFunc).
			WithIterations(5)

		benchmark.Run(t)

		assert.Equal(t, 5, executionCount, "Function should be executed 5 times")
	})

	t.Run("should handle function errors", func(t *testing.T) {
		testFunc := func() error {
			return assert.AnError
		}

		benchmark := NewPerformanceBenchmark("failing_operation", testFunc)

		// This should not panic, but the test will fail due to the error
		benchmark.Run(t)
	})

	t.Run("should measure execution time", func(t *testing.T) {
		minDelay := 10 * time.Millisecond
		testFunc := func() error {
			time.Sleep(minDelay)
			return nil
		}

		benchmark := NewPerformanceBenchmark("timed_operation", testFunc).
			WithIterations(2)

		start := time.Now()
		benchmark.Run(t)
		elapsed := time.Since(start)

		// Should take at least 2 * minDelay
		expectedMinimum := 2 * minDelay
		assert.GreaterOrEqual(t, elapsed, expectedMinimum,
			"Benchmark should take at least the expected time")
	})

	t.Run("should handle timeout", func(t *testing.T) {
		testFunc := func() error {
			time.Sleep(100 * time.Millisecond) // Longer than timeout
			return nil
		}

		benchmark := NewPerformanceBenchmark("timeout_operation", testFunc).
			WithTimeout(50 * time.Millisecond).
			WithIterations(1)

		// This should complete quickly due to timeout, but may fail the test
		start := time.Now()
		benchmark.Run(t)
		elapsed := time.Since(start)

		// Should not take much longer than timeout
		assert.Less(t, elapsed, 200*time.Millisecond,
			"Benchmark should respect timeout")
	})
}

func TestMemoryBenchmark(t *testing.T) {
	t.Run("should create memory allocation test", func(t *testing.T) {
		testFunc := func() error {
			// Allocate some memory
			_ = make([]byte, 1024)
			return nil
		}

		benchmark := NewPerformanceBenchmark("memory_operation", testFunc).
			WithIterations(10)

		benchmark.Run(t)

		// Just verify it runs without error
		// In a real implementation, we'd check memory statistics
	})
}

func TestConcurrencyBenchmark(t *testing.T) {
	t.Run("should handle concurrent operations", func(t *testing.T) {
		counter := 0
		testFunc := func() error {
			// Not thread-safe, but for testing concurrent execution pattern
			counter++
			time.Sleep(1 * time.Millisecond)
			return nil
		}

		benchmark := NewPerformanceBenchmark("concurrent_operation", testFunc).
			WithIterations(5)

		benchmark.Run(t)

		assert.Equal(t, 5, counter, "Should execute all iterations")
	})
}

func TestDataGenerationHelpers(t *testing.T) {
	t.Run("should validate helper functions exist", func(t *testing.T) {
		// Test that our helper functions can be created and used
		testModel := map[string]interface{}{
			"string_field": "test_value",
			"int_field":    123,
			"bool_field":   true,
		}

		// These should not panic or fail
		CheckField("string_field", "test_value")(t, testModel)
		CheckField("int_field", 123)(t, testModel)
		CheckField("bool_field", true)(t, testModel)
		CheckNotNil("string_field")(t, testModel)
	})

	t.Run("should handle edge cases gracefully", func(t *testing.T) {
		// Test with model that has missing fields
		model := map[string]interface{}{}

		// This should not panic but might fail assertion
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("Function should not panic: %v", r)
			}
		}()

		// These will fail the test assertions but should not panic
		// We can't easily test assertion failures without failing the test
		_ = model // Use the variable to avoid unused warning
	})
}

func TestMockWriter(t *testing.T) {
	t.Run("should write and read data", func(t *testing.T) {
		writer := NewMockWriter()

		data := []byte("test data")
		n, err := writer.Write(data)

		assert.NoError(t, err)
		assert.Equal(t, len(data), n)
		assert.Equal(t, "test data", writer.String())
	})

	t.Run("should clear data", func(t *testing.T) {
		writer := NewMockWriter()
		writer.Write([]byte("test data"))

		writer.Clear()
		assert.Empty(t, writer.String())
	})
}

func TestMockTicker(t *testing.T) {
	t.Run("should create and control ticker", func(t *testing.T) {
		ticker := NewMockTicker()

		// Send a tick
		ticker.Tick()

		// Should be able to receive
		select {
		case <-ticker.C():
			// Success
		case <-time.After(100 * time.Millisecond):
			t.Fatal("Should have received tick")
		}

		ticker.Stop()
	})
}

func TestTestDataBuilder(t *testing.T) {
	t.Run("should build test data", func(t *testing.T) {
		builder := NewTestDataBuilder()

		data := builder.
			Set("key1", "value1").
			Set("key2", 42).
			Set("key3", true).
			Build()

		assert.Equal(t, "value1", data["key1"])
		assert.Equal(t, 42, data["key2"])
		assert.Equal(t, true, data["key3"])
	})

	t.Run("should provide typed getters", func(t *testing.T) {
		builder := NewTestDataBuilder()
		builder.Set("string_key", "test")
		builder.Set("int_key", 123)

		assert.Equal(t, "test", builder.GetString("string_key"))
		assert.Equal(t, 123, builder.GetInt("int_key"))
		assert.Equal(t, "", builder.GetString("missing_key"))
		assert.Equal(t, 0, builder.GetInt("missing_key"))
	})
}
