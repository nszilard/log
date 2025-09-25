package log

import (
	"bytes"
	"errors"
	"strings"
	"sync"
	"testing"
	"time"
)

func setupTestLogger(t *testing.T, level Level) (buffer *bytes.Buffer, cleanup func()) {
	t.Helper()
	var buf bytes.Buffer
	old := std
	SetOutput(&buf)
	SetLevel(level)
	return &buf, func() { std = old }
}

func TestPackageLevelFunctions(t *testing.T) {
	buf, cleanup := setupTestLogger(t, DebugLevel)
	defer cleanup()

	tests := []struct {
		name    string
		logFn   func()
		level   string
		message string
	}{
		{"Debug", func() { Debug("debug msg") }, "[DEBUG]", "debug msg"},
		{"Info", func() { Info("info msg") }, "[INFO]", "info msg"},
		{"Warn", func() { Warn("warn msg") }, "[WARN]", "warn msg"},
		{"Error", func() { Error("error msg") }, "[ERROR]", "error msg"},
		{"Fatal", func() { Fatal("fatal msg") }, "[FATAL]", "fatal msg"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf.Reset()
			tt.logFn()
			output := buf.String()
			if !strings.Contains(output, tt.level) || !strings.Contains(output, tt.message) {
				t.Errorf("Expected %s and %s in output: %s", tt.level, tt.message, output)
			}
		})
	}
}

func TestFormattedLogging(t *testing.T) {
	buf, cleanup := setupTestLogger(t, DebugLevel)
	defer cleanup()

	tests := []struct {
		name     string
		logFn    func()
		expected string
	}{
		{"Debugf", func() { Debugf("user %s has %d items", "john", 5) }, "user john has 5 items"},
		{"Infof", func() { Infof("error %d: %s", 404, "not found") }, "error 404: not found"},
		{"Warnf", func() { Warnf("value: %.2f", 3.14159) }, "value: 3.14"},
		{"Errorf", func() { Errorf("warning: %s", "test") }, "warning: test"},
		{"Fatalf", func() { Fatalf("fatal: %s", "error") }, "fatal: error"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf.Reset()
			tt.logFn()
			if !strings.Contains(buf.String(), tt.expected) {
				t.Errorf("Expected %s in output: %s", tt.expected, buf.String())
			}
		})
	}
}

func TestStructuredLogging(t *testing.T) {
	buf, cleanup := setupTestLogger(t, DebugLevel)
	defer cleanup()

	tests := []struct {
		name   string
		logFn  func()
		checks []string
	}{
		{
			"InfoS basic",
			func() { InfoS(WithString("user", "john"), WithInt("count", 42)) },
			[]string{`"user":"john"`, `"count":42`},
		},
		{
			"ErrorS mixed types",
			func() { ErrorS(WithBool("admin", true), WithFloat("score", 95.5)) },
			[]string{`"admin":true`, `"score":95.5`},
		},
		{
			"DebugS with error",
			func() { DebugS(WithError("error", errors.New("test error"))) },
			[]string{`"error":"test error"`},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf.Reset()
			tt.logFn()
			output := buf.String()
			for _, check := range tt.checks {
				if !strings.Contains(output, check) {
					t.Errorf("Expected %s in output: %s", check, output)
				}
			}
		})
	}
}

func TestFieldTypes(t *testing.T) {
	buf, cleanup := setupTestLogger(t, InfoLevel)
	defer cleanup()

	testTime := time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC)
	testErr := errors.New("test error")

	tests := []struct {
		name  string
		field Data
		check string
	}{
		{"String", WithString("str", "hello"), `"str":"hello"`},
		{"Int", WithInt("num", 42), `"num":42`},
		{"Float", WithFloat("float", 3.14), `"float":3.14`},
		{"Bool true", WithBool("flag", true), `"flag":true`},
		{"Bool false", WithBool("flag", false), `"flag":false`},
		{"Duration", WithDuration("dur", time.Second), `"dur":1000000000`},
		{"Time", WithTime("time", testTime), `"time":"2023-01-01T12:00:00Z"`},
		{"Error", WithError("err", testErr), `"err":"test error"`},
		{"Nil Error", WithError("err", nil), `"err":""`},
		{"Any", WithAny("any", "value"), `"any":"value"`},
		{"Nil Any", WithAny("nil", nil), `"nil":null`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf.Reset()
			InfoS(tt.field)
			if !strings.Contains(buf.String(), tt.check) {
				t.Errorf("Expected %s in output: %s", tt.check, buf.String())
			}
		})
	}
}

func TestLevelFiltering(t *testing.T) {
	buf, cleanup := setupTestLogger(t, WarnLevel)
	defer cleanup()

	tests := []struct {
		name      string
		logFn     func()
		shouldLog bool
	}{
		{"Debug filtered", func() { Debug("test") }, false},
		{"Info filtered", func() { Info("test") }, false},
		{"Warn allowed", func() { Warn("test") }, true},
		{"Error allowed", func() { Error("test") }, true},
		{"DebugS filtered", func() { DebugS(WithString("test", "debug")) }, false},
		{"InfoS filtered", func() { InfoS(WithString("test", "info")) }, false},
		{"WarnS allowed", func() { WarnS(WithString("test", "warn")) }, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf.Reset()
			tt.logFn()
			hasOutput := buf.Len() > 0
			if hasOutput != tt.shouldLog {
				t.Errorf("Expected shouldLog=%v, got=%v", tt.shouldLog, hasOutput)
			}
		})
	}
}

func TestConfiguration(t *testing.T) {
	old := std
	defer func() { std = old }()

	t.Run("SetOutput", func(t *testing.T) {
		var buf1, buf2 bytes.Buffer
		SetOutput(&buf1)
		SetLevel(InfoLevel)
		Info("msg1")
		if buf1.Len() == 0 {
			t.Error("Should write to first buffer")
		}

		SetOutput(&buf2)
		Info("msg2")
		if buf2.Len() == 0 {
			t.Error("Should write to second buffer")
		}
	})

	t.Run("SetLevel", func(t *testing.T) {
		var buf bytes.Buffer
		SetOutput(&buf)

		SetLevel(InfoLevel)
		Debug("debug")
		if buf.Len() > 0 {
			t.Error("Debug should be filtered at Info level")
		}

		SetLevel(DebugLevel)
		Debug("debug")
		if buf.Len() == 0 {
			t.Error("Debug should not be filtered at Debug level")
		}
	})

	t.Run("SetIncludeFileInfo", func(t *testing.T) {
		var buf bytes.Buffer
		SetOutput(&buf)
		SetLevel(InfoLevel)

		SetIncludeFileInfo(true)
		Info("with file")
		if !strings.Contains(buf.String(), "log_test.go") {
			t.Error("Expected file info when enabled")
		}

		buf.Reset()
		SetIncludeFileInfo(false)
		Info("without file")
		if strings.Contains(buf.String(), "log_test.go") {
			t.Error("Expected no file info when disabled")
		}
	})
}

func TestPanicFunctions(t *testing.T) {
	_, cleanup := setupTestLogger(t, DebugLevel)
	defer cleanup()

	t.Run("Panic", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("Expected panic")
			}
		}()
		Panic("test panic")
	})

	t.Run("Panicf", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("Expected panic")
			}
		}()
		Panicf("test panic: %s", "formatted")
	})

	t.Run("PanicS", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("Expected panic")
			}
		}()
		PanicS(WithString("test", "panic"))
	})
}

func TestNoopLogger(t *testing.T) {
	noop := NoopLogger{}
	// These should not panic or cause any issues
	noop.SetLevel(InfoLevel)
	noop.SetIncludeFileInfo(true)
	noop.SetOutput(&bytes.Buffer{})

	// Test all methods
	noop.Debug("test")
	noop.Debugf("test %s", "debug")
	noop.DebugS(WithString("test", "debug"))
	noop.Info("test")
	noop.Infof("test %s", "info")
	noop.InfoS(WithString("test", "info"))
	noop.Warn("test")
	noop.Warnf("test %s", "warn")
	noop.WarnS(WithString("test", "warn"))
	noop.Error("test")
	noop.Errorf("test %s", "error")
	noop.ErrorS(WithString("test", "error"))
	noop.Fatal("test")
	noop.Fatalf("test %s", "fatal")
	noop.FatalS(WithString("test", "fatal"))
	noop.Panic("test")
	noop.Panicf("test %s", "panic")
	noop.PanicS(WithString("test", "panic"))
}

func TestEdgeCases(t *testing.T) {
	buf, cleanup := setupTestLogger(t, InfoLevel)
	defer cleanup()

	t.Run("Empty structured fields", func(t *testing.T) {
		InfoS()
		if buf.Len() > 0 {
			t.Error("Empty structured log should not produce output")
		}
	})

	t.Run("Mixed field types", func(t *testing.T) {
		buf.Reset()
		InfoS(
			WithString("typed", "field"),
			WithAny("untyped", "field"),
			WithInt("number", 42),
		)
		output := buf.String()
		checks := []string{`"typed":"field"`, `"untyped":"field"`, `"number":42`}
		for _, check := range checks {
			if !strings.Contains(output, check) {
				t.Errorf("Expected %s in output: %s", check, output)
			}
		}
	})

	t.Run("Zero and negative values", func(t *testing.T) {
		tests := []struct {
			field Data
			check string
		}{
			{WithInt("zero", 0), `"zero":0`},
			{WithInt("neg", -42), `"neg":-42`},
			{WithFloat("zero_f", 0.0), `"zero_f":0`},
			{WithFloat("neg_f", -3.14), `"neg_f":-3.14`},
			{WithString("empty", ""), `"empty":""`},
		}

		for _, tt := range tests {
			buf.Reset()
			InfoS(tt.field)
			if !strings.Contains(buf.String(), tt.check) {
				t.Errorf("Expected %s in output: %s", tt.check, buf.String())
			}
		}
	})
}

func TestConcurrentLogging(t *testing.T) {
	buf, cleanup := setupTestLogger(t, InfoLevel)
	defer cleanup()

	var wg sync.WaitGroup
	numGoroutines := 10
	messagesPerGoroutine := 10

	wg.Add(numGoroutines)
	for i := range numGoroutines {
		go func(id int) {
			defer wg.Done()
			for j := range messagesPerGoroutine {
				InfoS(WithInt("goroutine", int64(id)), WithInt("message", int64(j)))
			}
		}(i)
	}

	wg.Wait()
	output := buf.String()
	lines := strings.Split(strings.TrimSpace(output), "\n")
	expectedLines := numGoroutines * messagesPerGoroutine
	if len(lines) != expectedLines {
		t.Errorf("Expected %d log lines, got %d", expectedLines, len(lines))
	}
}
