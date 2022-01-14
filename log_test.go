package log

import (
	"bytes"
	"fmt"
	"os"
	"regexp"
	"testing"
	"time"

	"github.com/nszilard/log/internal"
	"github.com/nszilard/log/models"
)

func TestNew(t *testing.T) {
	lis := &strLogListener{}
	logger := New(models.InfoLevel, lis, internal.DefaultLogLayout)
	logMsg := "test log"
	logger.Info(logMsg)
	now := time.Now()

	re := compileRegex(t, now, models.InfoLevel, logMsg)
	ok := re.MatchString(lis.log)
	if !ok {
		t.Errorf("pattern %q did not match with: %q", re, lis.log)
	}
}

func TestDefaultLogLevel(t *testing.T) {
	lis := &strLogListener{}
	logger := New(models.InfoLevel, lis, internal.DefaultLogLayout)
	logMsg := "test log"

	logger.Debug(logMsg)
	if len(lis.log) != 0 {
		t.Errorf("'DEBUG' level log should not be output on 'INFO' log level")
	}
}

func TestSetLevel(t *testing.T) {
	SetLevelDebug()
	if stdLogger.Level != models.DebugLevel {
		t.Errorf("expected level to be: %v, got: %v", models.DebugLevel, stdLogger.Level)
	}

	SetLevelInfo()
	if stdLogger.Level != models.InfoLevel {
		t.Errorf("expected level to be: %v, got: %v", models.InfoLevel, stdLogger.Level)
	}
}

func TestStandardLoggers(t *testing.T) {
	defer SetLevelInfo()
	logMsg := "hello world"
	cases := []struct {
		name  string
		level models.Level
		run   func(*testing.T)
	}{
		{
			name:  "Panic",
			level: models.PanicLevel,
			run: func(t *testing.T) {
				Panic(logMsg)
			},
		},
		{
			name:  "Panicf",
			level: models.PanicLevel,
			run: func(t *testing.T) {
				Panicf("%v", logMsg)
			}},
		{
			name:  "Fatal",
			level: models.FatalLevel,
			run: func(t *testing.T) {
				Fatal(logMsg)
			},
		},
		{
			name:  "Fatalf",
			level: models.FatalLevel,
			run: func(t *testing.T) {
				Fatalf("%v", logMsg)
			}},
		{
			name:  "Error",
			level: models.ErrorLevel,
			run: func(t *testing.T) {
				Error(logMsg)
			},
		},
		{
			name:  "Errorf",
			level: models.ErrorLevel,
			run: func(t *testing.T) {
				Errorf("%v", logMsg)
			}},
		{
			name:  "Warn",
			level: models.WarnLevel,
			run: func(t *testing.T) {
				Warn(logMsg)
			},
		},
		{
			name:  "Warnf",
			level: models.WarnLevel,
			run: func(t *testing.T) {
				Warnf("%v", logMsg)
			},
		},
		{
			name:  "Println",
			level: models.NoLevel,
			run: func(t *testing.T) {
				Println(logMsg)
			},
		},
		{
			name:  "Printf",
			level: models.NoLevel,
			run: func(t *testing.T) {
				Printf("%v", logMsg)
			},
		},
		{
			name:  "Info",
			level: models.InfoLevel,
			run: func(t *testing.T) {
				Info(logMsg)
			}},
		{
			name:  "Infof",
			level: models.InfoLevel,
			run: func(t *testing.T) {
				Infof("%v", logMsg)
			},
		},
		{
			name:  "Trace",
			level: models.TraceLevel,
			run: func(t *testing.T) {
				Trace(logMsg)
			}},
		{
			name:  "Tracef",
			level: models.TraceLevel,
			run: func(t *testing.T) {
				Tracef("%v", logMsg)
			}},
		{
			name:  "Debug",
			level: models.DebugLevel,
			run: func(t *testing.T) {
				Debug(logMsg)
			}},
		{
			name:  "Debugf",
			level: models.DebugLevel,
			run: func(t *testing.T) {
				Debugf("%v", logMsg)
			}},
	}

	for _, c := range cases {
		SetLevelDebug()
		now := time.Now()
		c.run(t)
		re := compileRegex(t, now, c.level, logMsg)
		ok := re.MatchString(string(stdLogger.Buf))
		if !ok {
			t.Errorf("%s logger returned unexpected log: %q", c.name, string(stdLogger.Buf))
		}
	}
}

func TestSetOutput(t *testing.T) {
	defer SetOutput(os.Stdout)

	now := time.Now()
	logMsg := "test"
	re := compileRegex(t, now, models.InfoLevel, logMsg)

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
	l := New(models.InfoLevel, &strLogListener{}, "")
	for i := 0; i < b.N; i++ {
		buf.Reset()
		l.Info(testString)
	}
}

func BenchmarkLogDefaultLayout(b *testing.B) {
	const testString = "test"
	var buf bytes.Buffer
	l := New(models.InfoLevel, &strLogListener{}, internal.DefaultLogLayout)
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

func compileRegex(t *testing.T, now time.Time, level models.Level, message string) *regexp.Regexp {
	dateStr := fmt.Sprintf("%04d", int(now.Year())) + "\\/" + fmt.Sprintf("%02d", int(now.Month())) + "\\/" + fmt.Sprintf("%02d", int(now.Day()))
	timeStr := fmt.Sprintf("%02d", int(now.Hour())) + ":" + fmt.Sprintf("%02d", int(now.Minute())) + ":" + fmt.Sprintf("%02d", int(now.Second()))

	pattern := "^" + dateStr + " " + timeStr + " " + "\\[" + models.LevelToString(level) + "\\] \\(\\w+\\.go:\\d+\\) ▶ " + message + "\n"

	re, err := regexp.Compile(pattern)
	if err != nil {
		t.Fatalf("pattern %q did not compile: %s", pattern, err)
	}

	return re
}
