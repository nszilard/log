package internal

import (
	"io"
	"sync"
	"time"
)

// Logger provides thread-safe logging functionality
type Logger struct {
	mu  sync.Mutex
	out io.Writer
}

// New creates a new Logger that writes to the given io.Writer
func New(out io.Writer) *Logger {
	return &Logger{out: out}
}

// SetOutput changes the output destination for the logger
func (l *Logger) SetOutput(out io.Writer) {
	l.mu.Lock()
	l.out = out
	l.mu.Unlock()
}

// LogWithFileInfo logs a simple text message with optional file information
func (l *Logger) LogWithFileInfo(levelStr, msg string, includeFileInfo bool) {
	now := time.Now().UTC()
	file, line := "", 0
	if includeFileInfo {
		file, line = getCaller(3)
	}

	buf := getBuf(len(msg) + 200)
	defer putBuf(buf)
	*buf = (*buf)[:0]

	*buf = AppendTextHeader(*buf, now.Format(timestampFormat), levelStr, file, line, includeFileInfo)
	*buf = append(*buf, msg...)
	if len(*buf) == 0 || (*buf)[len(*buf)-1] != '\n' {
		*buf = append(*buf, '\n')
	}

	l.write(*buf)
}

// LogStructuredTypedWithFileInfo logs structured data using typed fields
func (l *Logger) LogStructuredTypedWithFileInfo(levelStr string, includeFileInfo bool, fields []Data) {
	now := time.Now().UTC()
	file, line := "", 0
	if includeFileInfo {
		file, line = getCaller(4)
	}

	buf := getBuf(200 + len(fields)*50)
	defer putBuf(buf)
	*buf = (*buf)[:0]

	BuildStructuredHeader(buf, now, levelStr, includeFileInfo, file, line)
	for _, field := range fields {
		*buf = AppendJSONKey(*buf, field.Key)
		*buf = AppendTypedJSONValue(*buf, &field)
	}
	*buf = append(*buf, "}\n"...)

	l.write(*buf)
}

// LogStructuredWithFileInfo logs structured data using key-value pairs
func (l *Logger) LogStructuredWithFileInfo(levelStr string, includeFileInfo bool, keyValuePairs ...any) {
	now := time.Now().UTC()
	file, line := "", 0
	if includeFileInfo {
		file, line = getCaller(4)
	}

	buf := getBuf(200 + len(keyValuePairs)*50)
	defer putBuf(buf)
	*buf = (*buf)[:0]

	BuildStructuredHeader(buf, now, levelStr, includeFileInfo, file, line)
	for i := 0; i < len(keyValuePairs)-1; i += 2 {
		if key, ok := keyValuePairs[i].(string); ok {
			*buf = AppendJSONKey(*buf, key)
			*buf = AppendJSONValue(*buf, keyValuePairs[i+1])
		}
	}
	*buf = append(*buf, "}\n"...)

	l.write(*buf)
}

// write is a helper method that handles thread-safe writing to the output
func (l *Logger) write(data []byte) {
	l.mu.Lock()
	_, _ = l.out.Write(data)
	l.mu.Unlock()
}
