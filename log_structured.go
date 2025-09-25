package log

import "time"

// PanicS logs a structured message at PanicLevel and then panics.
func PanicS(fields ...Data) {
	logStructured(PanicLevel, fields)
	panic("panic")
}

// FatalS logs a structured message at FatalLevel.
func FatalS(fields ...Data) {
	logStructured(FatalLevel, fields)
}

// ErrorS logs a structured message at ErrorLevel.
func ErrorS(fields ...Data) {
	logStructured(ErrorLevel, fields)
}

// WarnS logs a structured message at WarnLevel.
func WarnS(fields ...Data) {
	logStructured(WarnLevel, fields)
}

// InfoS logs a structured message at InfoLevel.
func InfoS(fields ...Data) {
	logStructured(InfoLevel, fields)
}

// DebugS logs a structured message at DebugLevel.
func DebugS(fields ...Data) {
	logStructured(DebugLevel, fields)
}

// WithString adds a string key-value pair to the structured logger
func WithString(key, val string) Data {
	return Data{Key: key, Type: StringType, String: val}
}

// WithInt adds an int64 key-value pair to the structured logger
func WithInt(key string, val int64) Data {
	return Data{Key: key, Type: IntType, Integer: val}
}

// WithFloat adds a float64 key-value pair to the structured logger
func WithFloat(key string, val float64) Data {
	return Data{Key: key, Type: FloatType, Float: val}
}

// WithBool adds a bool key-value pair to the structured logger
func WithBool(key string, val bool) Data {
	return Data{Key: key, Type: BoolType, Bool: val}
}

// WithError adds an error key-value pair to the structured logger
func WithError(key string, val error) Data {
	if val == nil {
		return Data{Key: key, Type: ErrorType, Interface: nil}
	}
	return Data{Key: key, Type: ErrorType, String: val.Error(), Interface: val}
}

// WithDuration adds a time.Duration key-value pair to the structured logger
func WithDuration(key string, val time.Duration) Data {
	return Data{Key: key, Type: DurationType, Integer: int64(val)}
}

// WithTime adds a time.Time key-value pair to the structured logger
func WithTime(key string, val time.Time) Data {
	return Data{Key: key, Type: TimeType, Interface: val}
}

// WithAny adds an any key-value pair to the structured logger
func WithAny(key string, val any) Data {
	return Data{Key: key, Type: UnknownType, Interface: val}
}
