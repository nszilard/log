package log

import "io"

// FieldType represents the type of a structured logging field.
type FieldType uint8

const (
	UnknownType FieldType = iota
	StringType
	IntType
	FloatType
	BoolType
	ErrorType
	DurationType
	TimeType
)

// Data represents a key-value pair for structured logging with type-specific storage
type Data struct {
	Key  string
	Type FieldType

	// Type-specific storage
	String    string
	Integer   int64
	Float     float64
	Bool      bool
	Interface any
}

// Logger provides leveled and structured logging.
// All methods are safe for concurrent use.
type Logger interface {
	// SetLevel sets the minimum level for logging.
	SetLevel(level Level)
	// SetOutput sets the output destination for logs.
	SetOutput(out io.Writer)
	// SetIncludeFileInfo sets whether to include file and line information in logs.
	SetIncludeFileInfo(include bool)

	// Panic
	Panic(v ...any)
	Panicf(format string, v ...any)
	PanicS(fields ...Data)

	// Fatal
	Fatal(v ...any)
	Fatalf(format string, v ...any)
	FatalS(fields ...Data)

	// Error
	Error(v ...any)
	Errorf(format string, v ...any)
	ErrorS(fields ...Data)

	// Warn
	Warn(v ...any)
	Warnf(format string, v ...any)
	WarnS(fields ...Data)

	// Info
	Info(v ...any)
	Infof(format string, v ...any)
	InfoS(fields ...Data)

	// Debug
	Debug(v ...any)
	Debugf(format string, v ...any)
	DebugS(fields ...Data)
}
