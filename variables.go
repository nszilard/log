package log

import "os"

const defaultLogLayout = "%D %T %L (%f:%i) ▶ %l"

// logging levels
const (
	PanicLevel Level = iota
	FatalLevel
	ErrorLevel
	WarnLevel
	NoLevel
	InfoLevel
	TraceLevel
	DebugLevel
)

// Level2Str used to convert Level value to string.
var Level2Str = []string{
	"PANIC", // PanicLevel
	"FATAL", // FatalLevel
	"ERROR", // ErrorLevel
	"WARN",  // WarnLevel
	"-",     // NoLevel
	"INFO",  // InfoLevel
	"TRACE", // TraceLevel
	"DEBUG", // DebugLevel
}

var stdLogger = New(InfoLevel, os.Stdout, defaultLogLayout)
