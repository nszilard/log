package log

import (
	"bytes"
	"fmt"
	"os"
	"testing"
	"time"
)

func TestLoggers(t *testing.T) {
	lis := &strLogListener{}
	logger := New(DebugLevel, lis, defaultLogLayout)
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
				logger.Panic(logMsg)
			},
		},
		{
			name:  "Panicf",
			level: PanicLevel,
			run: func(t *testing.T) {
				logger.Panicf("%v", logMsg)
			}},
		{
			name:  "Fatal",
			level: FatalLevel,
			run: func(t *testing.T) {
				logger.Fatal(logMsg)
			},
		},
		{
			name:  "Fatalf",
			level: FatalLevel,
			run: func(t *testing.T) {
				logger.Fatalf("%v", logMsg)
			}},
		{
			name:  "Error",
			level: ErrorLevel,
			run: func(t *testing.T) {
				logger.Error(logMsg)
			},
		},
		{
			name:  "Errorf",
			level: ErrorLevel,
			run: func(t *testing.T) {
				logger.Errorf("%v", logMsg)
			}},
		{
			name:  "Warn",
			level: WarnLevel,
			run: func(t *testing.T) {
				logger.Warn(logMsg)
			},
		},
		{
			name:  "Warnf",
			level: WarnLevel,
			run: func(t *testing.T) {
				logger.Warnf("%v", logMsg)
			},
		},
		{
			name:  "Println",
			level: NoLevel,
			run: func(t *testing.T) {
				logger.Println(logMsg)
			},
		},
		{
			name:  "Printf",
			level: NoLevel,
			run: func(t *testing.T) {
				logger.Printf("%v", logMsg)
			},
		},
		{
			name:  "Info",
			level: InfoLevel,
			run: func(t *testing.T) {
				logger.Info(logMsg)
			}},
		{
			name:  "Infof",
			level: InfoLevel,
			run: func(t *testing.T) {
				logger.Infof("%v", logMsg)
			},
		},
		{
			name:  "Trace",
			level: TraceLevel,
			run: func(t *testing.T) {
				logger.Trace(logMsg)
			}},
		{
			name:  "Tracef",
			level: TraceLevel,
			run: func(t *testing.T) {
				logger.Tracef("%v", logMsg)
			}},
		{
			name:  "Debug",
			level: DebugLevel,
			run: func(t *testing.T) {
				logger.Debug(logMsg)
			}},
		{
			name:  "Debugf",
			level: DebugLevel,
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

func TestSetUnrecognizedLayout(t *testing.T) {
	var out bytes.Buffer
	testMessage := "test"

	l := New(InfoLevel, &out, "")
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

	l := New(InfoLevel, &out, "")
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
	l := New(InfoLevel, os.Stdout, "")
	l.Info("original")
	l.SetOutput(&out)
	l.Info(testMessage)

	expected := fmt.Sprintf("%v\n", testMessage)
	if out.String() != expected {
		t.Errorf("expected: %q, got: %q", expected, out.String())
	}
}
