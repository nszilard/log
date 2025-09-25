package log

// Level represents the severity level of a log entry.
type Level uint8

const (
	PanicLevel Level = iota
	FatalLevel
	ErrorLevel
	WarnLevel
	InfoLevel
	DebugLevel
)

var levelNames = []string{
	"PANIC",
	"FATAL",
	"ERROR",
	"WARN",
	"INFO",
	"DEBUG",
}

// String returns the string representation of the level.
func (l Level) String() string {
	if int(l) >= len(levelNames) {
		return ""
	}
	return levelNames[l]
}
