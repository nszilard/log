package internal

import "time"

const (
	timestampFormat = "2006-01-02T15:04:05.000Z07:00"
)

// FieldType represents the type of a log field
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

// Data represents a typed log field
type Data struct {
	Key       string
	Type      FieldType
	String    string
	Integer   int64
	Float     float64
	Bool      bool
	Interface any
}

// Field creation helpers
func StringField(key, value string) Data {
	return Data{Key: key, Type: StringType, String: value}
}

func IntField(key string, value int64) Data {
	return Data{Key: key, Type: IntType, Integer: value}
}

func FloatField(key string, value float64) Data {
	return Data{Key: key, Type: FloatType, Float: value}
}

func BoolField(key string, value bool) Data {
	return Data{Key: key, Type: BoolType, Bool: value}
}

func ErrorField(key string, err error) Data {
	return Data{Key: key, Type: ErrorType, String: err.Error()}
}

func DurationField(key string, duration time.Duration) Data {
	return Data{Key: key, Type: DurationType, Integer: int64(duration)}
}

func TimeField(key string, t time.Time) Data {
	return Data{Key: key, Type: TimeType, Interface: t}
}
