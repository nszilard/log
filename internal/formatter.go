package internal

import (
	"encoding/json"
	"runtime"
	"strconv"
	"time"
)

// getCaller returns the filename and line number of the caller
func getCaller(skip int) (filename string, lineNumber int) {
	_, file, line, _ := runtime.Caller(skip)
	if file == "" {
		return "???", line
	}
	for i := len(file) - 1; i > 0; i-- {
		if file[i] == '/' {
			return file[i+1:], line
		}
	}
	return file, line
}

// Sprint formats values similar to fmt.Sprint but optimized for logging
func Sprint(v ...any) string {
	if len(v) == 0 {
		return ""
	}
	if len(v) == 1 {
		if s, ok := v[0].(string); ok {
			return s
		}
	}
	result := make([]byte, 0, 64)
	for i, val := range v {
		if i > 0 {
			result = append(result, ' ')
		}
		result = appendAny(result, val)
	}
	return string(result)
}

// Sprintf formats a string with arguments similar to fmt.Sprintf but optimized for logging
func Sprintf(format string, args ...any) string {
	if len(args) == 0 {
		return format
	}
	result := make([]byte, 0, len(format)+64)
	argIndex := 0

	for i := 0; i < len(format); i++ {
		if format[i] == '%' && i+1 < len(format) && argIndex < len(args) {
			next := i + 1
			if format[next] == '.' {
				i = handleFloatFormat(&result, format, i, args, &argIndex)
				continue
			}
			i = handleFormatSpecifier(&result, format, i, args, &argIndex)
		} else {
			result = append(result, format[i])
		}
	}
	return string(result)
}

// Text formatting functions

// AppendTextHeader formats and appends a text log header to the buffer
func AppendTextHeader(buf []byte, timestamp, levelStr, file string, line int, includeFileInfo bool) []byte {
	buf = append(buf, timestamp...)
	buf = append(buf, " ["...)
	buf = append(buf, levelStr...)
	buf = append(buf, ']')
	if includeFileInfo {
		buf = append(buf, " ("...)
		buf = append(buf, file...)
		buf = append(buf, ':')
		buf = strconv.AppendInt(buf, int64(line), 10)
		buf = append(buf, ')')
	}
	buf = append(buf, " ▶ "...)
	return buf
}

// JSON formatting functions

// BuildStructuredHeader builds a JSON header for structured logging
func BuildStructuredHeader(buf *[]byte, now time.Time, levelStr string, includeFileInfo bool, file string, line int) {
	*buf = append(*buf, `{"timestamp":"`...)
	*buf = append(*buf, now.Format(timestampFormat)...)
	*buf = append(*buf, `","level":"`...)
	*buf = append(*buf, levelStr...)
	*buf = append(*buf, '"')

	if includeFileInfo {
		*buf = append(*buf, `,"caller":"`...)
		*buf = append(*buf, file...)
		*buf = append(*buf, ':')
		*buf = strconv.AppendInt(*buf, int64(line), 10)
		*buf = append(*buf, '"')
	}
}

// AppendJSONKey appends a JSON key to the buffer
func AppendJSONKey(buf []byte, key string) []byte {
	buf = append(buf, ',', '"')
	buf = append(buf, key...)
	buf = append(buf, `":`...)
	return buf
}

// AppendQuoted appends a quoted string to the buffer
func AppendQuoted(buf []byte, s string) []byte {
	buf = append(buf, '"')
	buf = append(buf, s...)
	buf = append(buf, '"')
	return buf
}

// AppendJSONValue appends any value as JSON to the buffer
func AppendJSONValue(buf []byte, v any) []byte {
	switch val := v.(type) {
	case string:
		return AppendQuoted(buf, val)
	case int, int32, int64:
		return strconv.AppendInt(buf, toInt64(val), 10)
	case uint, uint32, uint64:
		return strconv.AppendUint(buf, toUint64(val), 10)
	case float32:
		return strconv.AppendFloat(buf, float64(val), 'f', -1, 32)
	case float64:
		return strconv.AppendFloat(buf, val, 'f', -1, 64)
	case bool:
		if val {
			return append(buf, "true"...)
		}
		return append(buf, "false"...)
	case error:
		return AppendQuoted(buf, val.Error())
	case nil:
		return append(buf, "null"...)
	case []byte:
		if json.Valid(val) {
			return append(buf, val...)
		}
		return AppendQuoted(buf, string(val))
	default:
		jsonData, err := json.Marshal(val)
		if err != nil {
			return append(buf, `"<error>"`...)
		}
		return append(buf, jsonData...)
	}
}

// AppendTypedJSONValue appends a typed field value as JSON to the buffer
func AppendTypedJSONValue(buf []byte, field *Data) []byte {
	switch field.Type {
	case StringType, ErrorType:
		return AppendQuoted(buf, field.String)
	case IntType, DurationType:
		return strconv.AppendInt(buf, field.Integer, 10)
	case FloatType:
		return strconv.AppendFloat(buf, field.Float, 'f', -1, 64)
	case BoolType:
		if field.Bool {
			return append(buf, "true"...)
		}
		return append(buf, "false"...)
	case TimeType:
		if t, ok := field.Interface.(time.Time); ok {
			return AppendQuoted(buf, t.Format(time.RFC3339Nano))
		}
		return append(buf, "null"...)
	default:
		return AppendJSONValue(buf, field.Interface)
	}
}

// General purpose formatting functions

// appendAny appends any value to the buffer as a string representation
func appendAny(buf []byte, v any) []byte {
	switch val := v.(type) {
	case string:
		return append(buf, val...)
	case int, int32, int64:
		return strconv.AppendInt(buf, toInt64(val), 10)
	case uint, uint32, uint64:
		return strconv.AppendUint(buf, toUint64(val), 10)
	case float32:
		return strconv.AppendFloat(buf, float64(val), 'f', -1, 32)
	case float64:
		return strconv.AppendFloat(buf, val, 'f', -1, 64)
	case bool:
		if val {
			return append(buf, "true"...)
		}
		return append(buf, "false"...)
	case error:
		return append(buf, val.Error()...)
	case nil:
		return append(buf, "<nil>"...)
	default:
		return append(buf, "<unknown>"...)
	}
}

// Type conversion helpers

func toInt64(v any) int64 {
	switch i := v.(type) {
	case int:
		return int64(i)
	case int32:
		return int64(i)
	case int64:
		return i
	default:
		return 0
	}
}

func toUint64(v any) uint64 {
	switch i := v.(type) {
	case uint:
		return uint64(i)
	case uint32:
		return uint64(i)
	case uint64:
		return i
	default:
		return 0
	}
}

// Format specifier handlers for Sprintf

func handleFormatSpecifier(result *[]byte, format string, i int, args []any, argIndex *int) int {
	next := i + 1
	switch format[next] {
	case 'v', 's':
		*result = appendAny(*result, args[*argIndex])
	case 'd':
		if val, ok := args[*argIndex].(int); ok {
			*result = strconv.AppendInt(*result, int64(val), 10)
		} else {
			*result = appendAny(*result, args[*argIndex])
		}
	case 'f':
		if val, ok := args[*argIndex].(float64); ok {
			*result = strconv.AppendFloat(*result, val, 'f', -1, 64)
		} else {
			*result = appendAny(*result, args[*argIndex])
		}
	case '%':
		*result = append(*result, '%')
	default:
		*result = append(*result, format[i])
		return i
	}
	*argIndex++
	return next
}

func handleFloatFormat(result *[]byte, format string, i int, args []any, argIndex *int) int {
	precStart := i + 2
	precEnd := precStart
	for precEnd < len(format) && format[precEnd] >= '0' && format[precEnd] <= '9' {
		precEnd++
	}
	precision := -1
	if p, err := strconv.Atoi(format[precStart:precEnd]); err == nil {
		precision = p
	}
	if precEnd < len(format) && format[precEnd] == 'f' {
		if val, ok := args[*argIndex].(float64); ok {
			*result = strconv.AppendFloat(*result, val, 'f', precision, 64)
		} else {
			*result = appendAny(*result, args[*argIndex])
		}
		*argIndex++
		return precEnd
	}
	return i
}
