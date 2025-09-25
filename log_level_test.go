package log

import "testing"

func TestLevelString(t *testing.T) {
	tests := map[Level]string{
		PanicLevel: "PANIC",
		FatalLevel: "FATAL",
		ErrorLevel: "ERROR",
		WarnLevel:  "WARN",
		InfoLevel:  "INFO",
		DebugLevel: "DEBUG",
		Level(99):  "", // Invalid level

	}

	for level, expected := range tests {
		if got := level.String(); got != expected {
			t.Errorf("Level(%d).String() = %q, want %q", level, got, expected)
		}
	}
}
