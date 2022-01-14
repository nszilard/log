package logger

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"regexp"
	"testing"
	"time"

	"github.com/nszilard/log/internal"
	"github.com/nszilard/log/models"
)

func TestLogger_SetLevel(t *testing.T) {
	lis := &strLogListener{}
	logger := new(models.DebugLevel, lis, internal.DefaultLogLayout, 2)

	logger.SetLevel(models.DebugLevel)
	if logger.Level != models.DebugLevel {
		t.Errorf("expected level to be: %v, got: %v", models.DebugLevel, logger.Level)
	}

	logger.SetLevel(models.InfoLevel)
	if logger.Level != models.InfoLevel {
		t.Errorf("expected level to be: %v, got: %v", models.InfoLevel, logger.Level)
	}
}

func TestLoggers(t *testing.T) {
	lis := &strLogListener{}
	logger := new(models.DebugLevel, lis, internal.DefaultLogLayout, 2)
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
				logger.Panic(logMsg)
			},
		},
		{
			name:  "Panicf",
			level: models.PanicLevel,
			run: func(t *testing.T) {
				logger.Panicf("%v", logMsg)
			}},
		{
			name:  "Fatal",
			level: models.FatalLevel,
			run: func(t *testing.T) {
				logger.Fatal(logMsg)
			},
		},
		{
			name:  "Fatalf",
			level: models.FatalLevel,
			run: func(t *testing.T) {
				logger.Fatalf("%v", logMsg)
			}},
		{
			name:  "Error",
			level: models.ErrorLevel,
			run: func(t *testing.T) {
				logger.Error(logMsg)
			},
		},
		{
			name:  "Errorf",
			level: models.ErrorLevel,
			run: func(t *testing.T) {
				logger.Errorf("%v", logMsg)
			}},
		{
			name:  "Warn",
			level: models.WarnLevel,
			run: func(t *testing.T) {
				logger.Warn(logMsg)
			},
		},
		{
			name:  "Warnf",
			level: models.WarnLevel,
			run: func(t *testing.T) {
				logger.Warnf("%v", logMsg)
			},
		},
		{
			name:  "Println",
			level: models.NoLevel,
			run: func(t *testing.T) {
				logger.Println(logMsg)
			},
		},
		{
			name:  "Printf",
			level: models.NoLevel,
			run: func(t *testing.T) {
				logger.Printf("%v", logMsg)
			},
		},
		{
			name:  "Info",
			level: models.InfoLevel,
			run: func(t *testing.T) {
				logger.Info(logMsg)
			}},
		{
			name:  "Infof",
			level: models.InfoLevel,
			run: func(t *testing.T) {
				logger.Infof("%v", logMsg)
			},
		},
		{
			name:  "Trace",
			level: models.TraceLevel,
			run: func(t *testing.T) {
				logger.Trace(logMsg)
			}},
		{
			name:  "Tracef",
			level: models.TraceLevel,
			run: func(t *testing.T) {
				logger.Tracef("%v", logMsg)
			}},
		{
			name:  "Debug",
			level: models.DebugLevel,
			run: func(t *testing.T) {
				logger.Debug(logMsg)
			}},
		{
			name:  "Debugf",
			level: models.DebugLevel,
			run: func(t *testing.T) {
				logger.Debugf("%v", logMsg)
			}},
	}

	for _, c := range cases {
		now := time.Now()
		c.run(t)
		re := compileRegex(t, now, c.level, logMsg)
		ok := re.MatchString(lis.log)
		if !ok {
			t.Errorf("%q logger returned unexpected log: %q", c.name, lis.log)
		}
	}
}

func TestUnreachableStack(t *testing.T) {
	var out bytes.Buffer
	testMessage := "test"

	l := new(models.InfoLevel, &out, "(%f:%i) ▶ %l", 100)
	l.Info(testMessage)

	expected := fmt.Sprintf("(???:0) ▶ %v\n", testMessage)
	if out.String() != expected {
		t.Errorf("expected: %q, got: %q", expected, out.String())
	}
}

func TestSetUnrecognizedLayout(t *testing.T) {
	var out bytes.Buffer
	testMessage := "test"

	l := new(models.InfoLevel, &out, "", 2)
	l.SetLayout("%as %la")
	l.Info(testMessage)

	expected := fmt.Sprintf("%vas %va\n", "%", testMessage)
	if out.String() != expected {
		t.Errorf("expected: %q, got: %q", expected, out.String())
	}
}

func TestSetLayoutUnfinishedLayout(t *testing.T) {
	var out bytes.Buffer
	testMessage := "test"

	l := new(models.InfoLevel, &out, "", 2)
	l.SetLayout("%")
	l.Info(testMessage)

	expected := fmt.Sprintf("%v\n", testMessage)
	if out.String() != expected {
		t.Errorf("expected: %q, got: %q", expected, out.String())
	}
}

func TestLoggerSetOutput(t *testing.T) {
	var out bytes.Buffer

	testMessage := "changed"
	l := new(models.InfoLevel, os.Stdout, "", 2)
	l.Info("original")
	l.SetOutput(&out)
	l.Info(testMessage)

	expected := fmt.Sprintf("%v\n", testMessage)
	if out.String() != expected {
		t.Errorf("expected: %q, got: %q", expected, out.String())
	}
}

// ------------------------------------------
// Helpers
// ------------------------------------------

func new(lev models.Level, listener io.Writer, layout string, depth int) *Logger {
	l := &Logger{
		Level: lev,
		Out:   listener,
		Depth: depth,
	}
	l.SetLayout(layout)
	return l
}

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
