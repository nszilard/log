package log

import "io"

// NoopLogger implements Logger but discards all log messages and doesn't panic
type NoopLogger struct{}

func (NoopLogger) SetLevel(Level)                 {}
func (NoopLogger) SetOutput(io.Writer)            {}
func (NoopLogger) SetIncludeFileInfo(bool)        {}
func (NoopLogger) Debug(...any)                   {}
func (NoopLogger) Debugf(string, ...any)          {}
func (NoopLogger) DebugS(...Data)                 {}
func (NoopLogger) Info(...any)                    {}
func (NoopLogger) Infof(string, ...any)           {}
func (NoopLogger) InfoS(...Data)                  {}
func (NoopLogger) Warn(...any)                    {}
func (NoopLogger) Warnf(string, ...any)           {}
func (NoopLogger) WarnS(...Data)                  {}
func (NoopLogger) Error(...any)                   {}
func (NoopLogger) Errorf(string, ...any)          {}
func (NoopLogger) ErrorS(...Data)                 {}
func (NoopLogger) Fatal(...any)                   {}
func (NoopLogger) Fatalf(string, ...any)          {}
func (NoopLogger) FatalS(...Data)                 {}
func (NoopLogger) Panic(v ...any)                 {}
func (NoopLogger) Panicf(format string, v ...any) {}
func (NoopLogger) PanicS(fields ...Data)          {}
