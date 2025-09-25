package log

import (
	"errors"
	"strings"
	"testing"
	"time"
)

func TestStructuredFieldTypes(t *testing.T) {
	buf, cleanup := setupTestLogger(t, InfoLevel)
	defer cleanup()

	testTime := time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC)
	testErr := errors.New("test error")

	// Test comprehensive field type handling in a single structured log
	InfoS(
		WithString("string", "test"),
		WithInt("int", 42),
		WithInt("int64", 1234567890),
		WithAny("uint", uint(42)),
		WithAny("uint64", uint64(9876543210)),
		WithFloat("float32", 3.14),
		WithFloat("float64", 2.718281828),
		WithBool("bool_true", true),
		WithBool("bool_false", false),
		WithError("error", testErr),
		WithError("nil_error", nil),
		WithDuration("duration", time.Millisecond*500),
		WithTime("time", testTime),
		WithAny("any_string", "any value"),
		WithAny("any_nil", nil),
	)

	output := buf.String()
	expectedFields := []string{
		`"string":"test"`,
		`"int":42`,
		`"int64":1234567890`,
		`"uint":42`,
		`"uint64":9876543210`,
		`"float32":3.14`,
		`"float64":2.718281828`,
		`"bool_true":true`,
		`"bool_false":false`,
		`"error":"test error"`,
		`"nil_error":null`,
		`"duration":500000000`,
		`"time":"2023-01-01T12:00:00Z"`,
		`"any_string":"any value"`,
		`"any_nil":null`,
	}

	for _, expected := range expectedFields {
		if !strings.Contains(output, expected) {
			t.Errorf("Expected %s in output: %s", expected, output)
		}
	}
}

func TestStructuredEdgeCases(t *testing.T) {
	buf, cleanup := setupTestLogger(t, InfoLevel)
	defer cleanup()

	tests := []struct {
		name   string
		testFn func()
		check  func(string) bool
	}{
		{
			"Empty fields",
			func() { logStructured(InfoLevel, nil) },
			func(output string) bool { return output == "" },
		},
		{
			"Empty slice",
			func() { logStructured(InfoLevel, []Data{}) },
			func(output string) bool { return output == "" },
		},
		{
			"Zero values",
			func() {
				InfoS(
					WithInt("zero", 0),
					WithFloat("zero_f", 0.0),
					WithString("empty", ""),
					WithDuration("zero_dur", 0),
				)
			},
			func(output string) bool {
				return strings.Contains(output, `"zero":0`) &&
					strings.Contains(output, `"zero_f":0`) &&
					strings.Contains(output, `"empty":""`) &&
					strings.Contains(output, `"zero_dur":0`)
			},
		},
		{
			"Negative values",
			func() {
				InfoS(
					WithInt("neg", -42),
					WithFloat("neg_f", -3.14),
					WithDuration("neg_dur", -time.Second),
				)
			},
			func(output string) bool {
				return strings.Contains(output, `"neg":-42`) &&
					strings.Contains(output, `"neg_f":-3.14`) &&
					strings.Contains(output, `"neg_dur":-1000000000`)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf.Reset()
			tt.testFn()
			if !tt.check(buf.String()) {
				t.Errorf("Test %s failed, output: %s", tt.name, buf.String())
			}
		})
	}
}

func TestStructuredLevelFiltering(t *testing.T) {
	buf, cleanup := setupTestLogger(t, WarnLevel)
	defer cleanup()

	// Test that lower levels are filtered
	logStructured(DebugLevel, []Data{WithString("debug", "test")})
	logStructured(InfoLevel, []Data{WithString("info", "test")})
	if buf.Len() > 0 {
		t.Error("Lower levels should be filtered out")
	}

	// Test that same/higher levels pass through
	logStructured(WarnLevel, []Data{WithString("warn", "test")})
	if buf.Len() == 0 {
		t.Error("Same level should not be filtered out")
	}
}

func TestUntypedFieldHandling(t *testing.T) {
	buf, cleanup := setupTestLogger(t, InfoLevel)
	defer cleanup()

	testErr := errors.New("test error")
	testTime := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)

	// Test manual field construction to cover different field types
	fields := []Data{
		{Key: "string", Type: StringType, String: "test"},
		{Key: "int", Type: IntType, Integer: 42},
		{Key: "float", Type: FloatType, Float: 1.23},
		{Key: "bool", Type: BoolType, Bool: true},
		{Key: "error", Type: ErrorType, String: testErr.Error(), Interface: testErr},
		{Key: "nil_error", Type: ErrorType, Interface: nil},
		{Key: "duration", Type: DurationType, Integer: int64(time.Second * 5)},
		{Key: "time", Type: TimeType, Interface: testTime},
		{Key: "unknown", Type: UnknownType, Interface: "unknown_value"},
	}

	InfoS(fields...)
	output := buf.String()

	expectedChecks := []string{
		`"string":"test"`,
		`"int":42`,
		`"float":1.23`,
		`"bool":true`,
		`"error":"test error"`,
		`"nil_error":null`,
		`"unknown":"unknown_value"`,
	}

	for _, check := range expectedChecks {
		if !strings.Contains(output, check) {
			t.Errorf("Expected %s in output: %s", check, output)
		}
	}
}

func TestConcurrentStructuredLogging(t *testing.T) {
	buf, cleanup := setupTestLogger(t, InfoLevel)
	defer cleanup()

	done := make(chan bool, 50)

	for i := range 50 {
		go func(id int) {
			InfoS(
				WithInt("goroutine", int64(id)),
				WithString("message", "concurrent test"),
				WithBool("active", true),
			)
			done <- true
		}(i)
	}

	for range 50 {
		<-done
	}

	output := buf.String()
	lines := strings.Split(strings.TrimSpace(output), "\n")

	if len(lines) != 50 {
		t.Errorf("Expected 50 log lines, got %d", len(lines))
	}

	// Verify that all lines contain expected structured data
	for _, line := range lines {
		if !strings.Contains(line, `"message":"concurrent test"`) ||
			!strings.Contains(line, `"active":true`) {
			t.Errorf("Malformed concurrent log line: %s", line)
		}
	}
}
