package log

import (
	"bytes"
	"fmt"
	"os"
	"regexp"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	lis := &strLogListener{}
	logger := New(InfoLevel, lis, defaultLogLayout)
	logMsg := "test log"
	logger.Info(logMsg)
	now := time.Now()

	re := compileRegex(t, now, InfoLevel, logMsg)
	ok := re.MatchString(lis.log)
	if !ok {
		t.Errorf("pattern %q did not match with: %q", re, lis.log)
	}
}

func TestDefaultLogLevel(t *testing.T) {
	lis := &strLogListener{}
	logger := New(InfoLevel, lis, defaultLogLayout)
	logMsg := "test log"

	logger.Debug(logMsg)
	if len(lis.log) != 0 {
		t.Errorf("'DEBUG' level log should not be output on 'INFO' log level.")
	}
}

func TestSetLevel(t *testing.T) {
	SetLevelDebug()
	if stdLogger.level != DebugLevel {
		t.Errorf("Expected level to be: %v, got: %v", DebugLevel, stdLogger.level)
	}

	SetLevelInfo()
	if stdLogger.level != InfoLevel {
		t.Errorf("Expected level to be: %v, got: %v", InfoLevel, stdLogger.level)
	}
}

func TestStandardLoggers(t *testing.T) {
	defer SetLevelInfo()
	logMsg := "hello world"
	cases := []struct {
		name  string
		level Level
		run   func(*testing.T)
	}{
		{
			name:  "Panic",
			level: PanicLevel,
			run: func(t *testing.T) {
				Panic(logMsg)
			},
		},
		{
			name:  "Panicf",
			level: PanicLevel,
			run: func(t *testing.T) {
				Panicf("%v", logMsg)
			}},
		{
			name:  "Fatal",
			level: FatalLevel,
			run: func(t *testing.T) {
				Fatal(logMsg)
			},
		},
		{
			name:  "Fatalf",
			level: FatalLevel,
			run: func(t *testing.T) {
				Fatalf("%v", logMsg)
			}},
		{
			name:  "Error",
			level: ErrorLevel,
			run: func(t *testing.T) {
				Error(logMsg)
			},
		},
		{
			name:  "Errorf",
			level: ErrorLevel,
			run: func(t *testing.T) {
				Errorf("%v", logMsg)
			}},
		{
			name:  "Warn",
			level: WarnLevel,
			run: func(t *testing.T) {
				Warn(logMsg)
			},
		},
		{
			name:  "Warnf",
			level: WarnLevel,
			run: func(t *testing.T) {
				Warnf("%v", logMsg)
			},
		},
		{
			name:  "Println",
			level: NoLevel,
			run: func(t *testing.T) {
				Println(logMsg)
			},
		},
		{
			name:  "Printf",
			level: NoLevel,
			run: func(t *testing.T) {
				Printf("%v", logMsg)
			},
		},
		{
			name:  "Info",
			level: InfoLevel,
			run: func(t *testing.T) {
				Info(logMsg)
			}},
		{
			name:  "Infof",
			level: InfoLevel,
			run: func(t *testing.T) {
				Infof("%v", logMsg)
			},
		},
		{
			name:  "Trace",
			level: TraceLevel,
			run: func(t *testing.T) {
				Trace(logMsg)
			}},
		{
			name:  "Tracef",
			level: TraceLevel,
			run: func(t *testing.T) {
				Tracef("%v", logMsg)
			}},
		{
			name:  "Debug",
			level: DebugLevel,
			run: func(t *testing.T) {
				Debug(logMsg)
			}},
		{
			name:  "Debugf",
			level: DebugLevel,
			run: func(t *testing.T) {
				Debugf("%v", logMsg)
			}},
	}

	for _, c := range cases {
		SetLevelDebug()
		now := time.Now()
		c.run(t)
		re := compileRegex(t, now, c.level, logMsg)
		ok := re.MatchString(string(stdLogger.buf))
		if !ok {
			t.Errorf("%s logger returned unexpected log: %q", c.name, string(stdLogger.buf))
		}
	}
}

func TestSetOutput(t *testing.T) {
	defer SetOutput(os.Stdout)

	now := time.Now()
	logMsg := "test"
	re := compileRegex(t, now, InfoLevel, logMsg)

	var out bytes.Buffer
	SetOutput(&out)
	Info(logMsg)

	ok := re.MatchString(out.String())
	if !ok {
		t.Error("failed to change logger output")
	}
}

// ------------------------------------------
// Benchmark
// ------------------------------------------

func BenchmarkLogNoFlags(b *testing.B) {
	const testString = "test"
	var buf bytes.Buffer
	l := New(InfoLevel, &strLogListener{}, "")
	for i := 0; i < b.N; i++ {
		buf.Reset()
		l.Info(testString)
	}
}

func BenchmarkLogDefaultLayout(b *testing.B) {
	const testString = "test"
	var buf bytes.Buffer
	l := New(InfoLevel, &strLogListener{}, defaultLogLayout)
	for i := 0; i < b.N; i++ {
		buf.Reset()
		l.Info(testString)
	}
}

// ------------------------------------------
// Helpers
// ------------------------------------------

type strLogListener struct {
	log string
}

func (lis *strLogListener) Write(p []byte) (n int, err error) {
	lis.log = string(p)
	return n, nil
}

func (lis *strLogListener) Close() error {
	return nil
}

func compileRegex(t *testing.T, now time.Time, level Level, message string) *regexp.Regexp {
	dateStr := fmt.Sprintf("%04d", int(now.Year())) + "\\/" + fmt.Sprintf("%02d", int(now.Month())) + "\\/" + fmt.Sprintf("%02d", int(now.Day()))
	timeStr := fmt.Sprintf("%02d", int(now.Hour())) + ":" + fmt.Sprintf("%02d", int(now.Minute())) + ":" + fmt.Sprintf("%02d", int(now.Second()))

	pattern := "^" + dateStr + " " + timeStr + " " + "\\[" + Level2Str[level] + "\\] \\(\\w+\\.go:\\d+\\) ▶ " + message + "\n"

	re, err := regexp.Compile(pattern)
	if err != nil {
		t.Fatalf("pattern %q did not compile: %s", pattern, err)
	}

	return re
}
