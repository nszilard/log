// Package log provides high-performance leveled and structured logging.
//
// The logger supports both formatted logging (like fmt.Printf) and structured
// logging with JSON output. File and line information can be optionally included
// for debugging (enabled by default).
package log

import (
	"io"
	"os"

	"github.com/nszilard/log/internal"
)

type logger struct {
	internal        *internal.Logger
	currentLevel    Level
	includeFileInfo bool
}

// std is the default logger instance.
var std = &logger{
	internal:        internal.New(os.Stdout),
	currentLevel:    InfoLevel,
	includeFileInfo: true,
}

// SetLevel sets the minimum level for the default logger.
func SetLevel(level Level) {
	std.currentLevel = level
}

// SetOutput sets the output destination for the default logger.
func SetOutput(out io.Writer) {
	std.internal.SetOutput(out)
}

// SetIncludeFileInfo sets whether to include file and line information in logs.
func SetIncludeFileInfo(include bool) {
	std.includeFileInfo = include
}

// Panic logs a message at PanicLevel and then panics.
func Panic(v ...any) {
	msg := internal.Sprint(v...)
	if PanicLevel <= std.currentLevel {
		std.internal.LogWithFileInfo(PanicLevel.String(), msg, std.includeFileInfo)
	}
	panic(msg)
}

// Fatal logs a message at FatalLevel.
func Fatal(v ...any) {
	if FatalLevel <= std.currentLevel {
		std.internal.LogWithFileInfo(FatalLevel.String(), internal.Sprint(v...), std.includeFileInfo)
	}
}

// Error logs a message at ErrorLevel.
func Error(v ...any) {
	if ErrorLevel <= std.currentLevel {
		std.internal.LogWithFileInfo(ErrorLevel.String(), internal.Sprint(v...), std.includeFileInfo)
	}
}

// Warn logs a message at WarnLevel.
func Warn(v ...any) {
	if WarnLevel <= std.currentLevel {
		std.internal.LogWithFileInfo(WarnLevel.String(), internal.Sprint(v...), std.includeFileInfo)
	}
}

// Info logs a message at InfoLevel.
func Info(v ...any) {
	if InfoLevel <= std.currentLevel {
		std.internal.LogWithFileInfo(InfoLevel.String(), internal.Sprint(v...), std.includeFileInfo)
	}
}

// Debug logs a message at DebugLevel.
func Debug(v ...any) {
	if DebugLevel <= std.currentLevel {
		std.internal.LogWithFileInfo(DebugLevel.String(), internal.Sprint(v...), std.includeFileInfo)
	}
}
