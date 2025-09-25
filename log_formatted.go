package log

import (
	"github.com/nszilard/log/internal"
)

// Panicf logs a formatted message at PanicLevel and then panics.
func Panicf(format string, v ...any) {
	msg := internal.Sprintf(format, v...)
	if PanicLevel <= std.currentLevel {
		std.internal.LogWithFileInfo(PanicLevel.String(), msg, std.includeFileInfo)
	}
	panic(msg)
}

// Fatalf logs a formatted message at FatalLevel.
func Fatalf(format string, v ...any) {
	if FatalLevel <= std.currentLevel {
		std.internal.LogWithFileInfo(FatalLevel.String(), internal.Sprintf(format, v...), std.includeFileInfo)
	}
}

// Errorf logs a formatted message at ErrorLevel.
func Errorf(format string, v ...any) {
	if ErrorLevel <= std.currentLevel {
		std.internal.LogWithFileInfo(ErrorLevel.String(), internal.Sprintf(format, v...), std.includeFileInfo)
	}
}

// Warnf logs a formatted message at WarnLevel.
func Warnf(format string, v ...any) {
	if WarnLevel <= std.currentLevel {
		std.internal.LogWithFileInfo(WarnLevel.String(), internal.Sprintf(format, v...), std.includeFileInfo)
	}
}

// Infof logs a formatted message at InfoLevel.
func Infof(format string, v ...any) {
	if InfoLevel <= std.currentLevel {
		std.internal.LogWithFileInfo(InfoLevel.String(), internal.Sprintf(format, v...), std.includeFileInfo)
	}
}

// Debugf logs a formatted message at DebugLevel.
func Debugf(format string, v ...any) {
	if DebugLevel <= std.currentLevel {
		std.internal.LogWithFileInfo(DebugLevel.String(), internal.Sprintf(format, v...), std.includeFileInfo)
	}
}
