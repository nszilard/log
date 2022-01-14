package models

import (
	"testing"
)

func TestLevelToString(t *testing.T) {
	cases := []struct {
		name     string
		level    Level
		expected string
	}{
		{
			name:     "Panic",
			level:    PanicLevel,
			expected: "PANIC",
		},
		{
			name:     "Fatal",
			level:    FatalLevel,
			expected: "FATAL",
		},
		{
			name:     "Error",
			level:    ErrorLevel,
			expected: "ERROR",
		},
		{
			name:     "Warn",
			level:    WarnLevel,
			expected: "WARN",
		},
		{
			name:     "No level set",
			level:    NoLevel,
			expected: "-",
		},
		{
			name:     "Info",
			level:    InfoLevel,
			expected: "INFO",
		},
		{
			name:     "Trace",
			level:    TraceLevel,
			expected: "TRACE",
		},
		{
			name:     "Debug",
			level:    DebugLevel,
			expected: "DEBUG",
		},
	}

	for _, c := range cases {
		actual := LevelToString(c.level)
		if actual != c.expected {
			t.Errorf("%q: expected %q but got %q", c.name, c.expected, actual)
		}
	}
}
