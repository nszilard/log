package models

// Level used to filter log message by the Logger.
type Level uint8

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
var level2Str = []string{
	"PANIC", // PanicLevel
	"FATAL", // FatalLevel
	"ERROR", // ErrorLevel
	"WARN",  // WarnLevel
	"-",     // NoLevel
	"INFO",  // InfoLevel
	"TRACE", // TraceLevel
	"DEBUG", // DebugLevel
}

func LevelToString(l Level) string {
	return level2Str[l]
}
