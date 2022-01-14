// Package log implements a simple logging package.
package log

import (
	"fmt"
	"io"
	"os"

	"github.com/nszilard/log/internal"
	"github.com/nszilard/log/internal/pkg/logger"
	"github.com/nszilard/log/models"
)

var stdLogger = New(models.InfoLevel, os.Stdout, internal.DefaultLogLayout)

// New creates a new Logger.
func New(lev models.Level, listener io.Writer, layout string) *logger.Logger {
	logger := &logger.Logger{
		Level: lev,
		Out:   listener,
		Depth: 2,
	}
	logger.SetLayout(layout)
	return logger
}

// SetLevelInfo will set the log level to Info
func SetLevelInfo() {
	stdLogger.SetLevel(models.InfoLevel)
}

// SetLevelDebug will set the log level to Debug
func SetLevelDebug() {
	stdLogger.SetLevel(models.DebugLevel)
}

// SetOutput will set the output to the specified
func SetOutput(listener io.Writer) {
	stdLogger.SetOutput(listener)
}

// Panic is equivalent to Logger.Panic.
func Panic(v ...interface{}) {
	s := fmt.Sprint(v...)
	stdLogger.Log(models.PanicLevel, s)
}

// Panicf is equivalent to Logger.Panicf.
func Panicf(format string, v ...interface{}) {
	s := fmt.Sprintf(format, v...)
	stdLogger.Log(models.PanicLevel, s)
}

// Fatal is equivalent to Logger.Fatal.
func Fatal(v ...interface{}) {
	stdLogger.Log(models.FatalLevel, fmt.Sprint(v...))
}

// Fatalf is equivalent to Logger.Fatalf.
func Fatalf(format string, v ...interface{}) {
	stdLogger.Log(models.FatalLevel, fmt.Sprintf(format, v...))
}

// Error is equivalent to Logger.Error.
func Error(v ...interface{}) {
	stdLogger.Log(models.ErrorLevel, fmt.Sprint(v...))
}

// Errorf is equivalent to Logger.Errorf.
func Errorf(format string, v ...interface{}) {
	stdLogger.Log(models.ErrorLevel, fmt.Sprintf(format, v...))
}

// Warn is equivalent to Logger.Warn.
func Warn(v ...interface{}) {
	stdLogger.Log(models.WarnLevel, fmt.Sprint(v...))
}

// Warnf is equivalent to Logger.Warnf.
func Warnf(format string, v ...interface{}) {
	stdLogger.Log(models.WarnLevel, fmt.Sprintf(format, v...))
}

// Info is equivalent to Logger.Info.
func Info(v ...interface{}) {
	stdLogger.Log(models.InfoLevel, fmt.Sprint(v...))
}

// Infof is equivalent to Logger.Infof.
func Infof(format string, v ...interface{}) {
	stdLogger.Log(models.InfoLevel, fmt.Sprintf(format, v...))
}

// Trace is equivalent to Logger.Debug.
func Trace(v ...interface{}) {
	stdLogger.Log(models.TraceLevel, fmt.Sprint(v...))
}

// Tracef is equivalent to Logger.Debugf.
func Tracef(format string, v ...interface{}) {
	stdLogger.Log(models.TraceLevel, fmt.Sprintf(format, v...))
}

// Debug is equivalent to Logger.Debug.
func Debug(v ...interface{}) {
	stdLogger.Log(models.DebugLevel, fmt.Sprint(v...))
}

// Debugf is equivalent to Logger.Debugf.
func Debugf(format string, v ...interface{}) {
	stdLogger.Log(models.DebugLevel, fmt.Sprintf(format, v...))
}

// Println is equivalent to Logger.Log.
func Println(v ...interface{}) {
	stdLogger.Log(models.NoLevel, fmt.Sprint(v...))
}

// Printf is equivalent to Logger.Logf.
func Printf(format string, v ...interface{}) {
	stdLogger.Log(models.NoLevel, fmt.Sprintf(format, v...))
}
