package log

import (
	"fmt"
	"io"
	"runtime"
	"strings"
	"time"
)

// SetLevel set the log level for Logger.
func (l *Logger) SetLevel(lev Level) {
	if lev != l.level {
		l.mu.Lock()
		l.level = lev
		l.mu.Unlock()
	}
}

// SetLayout set the layout of log message.
func (l *Logger) SetLayout(logLayout string) {
	layout := logLayout

	l.mu.Lock()
	defer l.mu.Unlock()

	// clear before set.
	l.layouters = nil

	for {
		i := strings.IndexByte(layout, '%')
		if i != -1 {
			if i != 0 {
				l.layouters = append(l.layouters, &layoutPlaceholder{
					placeholder: layout[:i],
				})
			}

			if i+2 > len(layout) {
				break
			}

			f := layout[i : i+2]
			layouter := mapLayouts[f]
			if layouter != nil {
				l.layouters = append(l.layouters, layouter)
				switch layouter.(type) {
				case *layoutFile, *layoutShortFile, *layoutLine:
					l.needCallerInfo = true
				}
			} else {
				l.layouters = append(l.layouters, &layoutPlaceholder{
					placeholder: f,
				})
			}

			layout = layout[i+2:]
		} else {
			if len(layout) > 0 {
				l.layouters = append(l.layouters, &layoutPlaceholder{
					placeholder: layout,
				})
			}
			break
		}
	}
}

// SetOutput sets the output destination for the logger.
func (l *Logger) SetOutput(w io.Writer) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.out = w
}

// Panic print a PanicLevel message to the logger.
// Arguments are handled in the manner of fmt.Print.
func (l *Logger) Panic(v ...interface{}) {
	s := fmt.Sprint(v...)
	l.Log(PanicLevel, s)
}

// Panicf print a PanicLevel message to the logger.
// Arguments are handled in the manner of fmt.Printf.
func (l *Logger) Panicf(format string, v ...interface{}) {
	s := fmt.Sprintf(format, v...)
	l.Log(PanicLevel, s)
}

// Fatal print a FatalLevel message to the logger.
// Arguments are handled in the manner of fmt.Print.
func (l *Logger) Fatal(v ...interface{}) {
	l.Log(FatalLevel, fmt.Sprint(v...))
}

// Fatalf print a FatalLevel message to the logger.
// Arguments are handled in the manner of fmt.Printf.
func (l *Logger) Fatalf(format string, v ...interface{}) {
	l.Log(FatalLevel, fmt.Sprintf(format, v...))
}

// Error print an ErrorLevel message to the logger.
// Arguments are handled in the manner of fmt.Print.
func (l *Logger) Error(v ...interface{}) {
	l.Log(ErrorLevel, fmt.Sprint(v...))
}

// Errorf print an ErrorLevel message to the logger.
// Arguments are handled in the manner of fmt.Printf.
func (l *Logger) Errorf(format string, v ...interface{}) {
	l.Log(ErrorLevel, fmt.Sprintf(format, v...))
}

// Warn print a WarnLevel message to the logger.
// Arguments are handled in the manner of fmt.Print.
func (l *Logger) Warn(v ...interface{}) {
	l.Log(WarnLevel, fmt.Sprint(v...))
}

// Warnf print a WarnLevel message to the logger.
// Arguments are handled in the manner of fmt.Printf.
func (l *Logger) Warnf(format string, v ...interface{}) {
	l.Log(WarnLevel, fmt.Sprintf(format, v...))
}

// Info print an InfoLevel message to the logger.
// Arguments are handled in the manner of fmt.Print.
func (l *Logger) Info(v ...interface{}) {
	l.Log(InfoLevel, fmt.Sprint(v...))
}

// Infof print an InfoLevel message to the logger.
// Arguments are handled in the manner of fmt.Printf.
func (l *Logger) Infof(format string, v ...interface{}) {
	l.Log(InfoLevel, fmt.Sprintf(format, v...))
}

// Println print a NoLevel message to the logger.
// Arguments are handled in the manner of fmt.Print.
func (l *Logger) Println(v ...interface{}) {
	l.Log(NoLevel, fmt.Sprint(v...))
}

// Printf print a NoLevel message to the logger.
// Arguments are handled in the manner of fmt.Printf.
func (l *Logger) Printf(format string, v ...interface{}) {
	l.Log(NoLevel, fmt.Sprintf(format, v...))
}

// Trace print a TraceLevel message to the logger.
// Arguments are handled in the manner of fmt.Print.
func (l *Logger) Trace(v ...interface{}) {
	l.Log(TraceLevel, fmt.Sprint(v...))
}

// Tracef print a TraceLevel message to the logger.
// Arguments are handled in the manner of fmt.Printf.
func (l *Logger) Tracef(format string, v ...interface{}) {
	l.Log(TraceLevel, fmt.Sprintf(format, v...))
}

// Debug print a DebugLevel message to the logger.
// Arguments are handled in the manner of fmt.Print.
func (l *Logger) Debug(v ...interface{}) {
	l.Log(DebugLevel, fmt.Sprint(v...))
}

// Debugf print a DebugLevel message to the logger.
// Arguments are handled in the manner of fmt.Printf.
func (l *Logger) Debugf(format string, v ...interface{}) {
	l.Log(DebugLevel, fmt.Sprintf(format, v...))
}

// Log print a leveled message to the logger.
// Arguments are handled in the manner of fmt.Print.
func (l *Logger) Log(level Level, msg string) {
	now := time.Now() // get this early.
	var file string
	var line int
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.level >= level {
		l.buf = l.buf[:0]

		if l.needCallerInfo {
			// release lock while getting caller info - it's expensive.
			l.mu.Unlock()
			var ok bool
			_, file, line, ok = runtime.Caller(2)
			if !ok {
				file = "???"
				line = 0
			}
			// re-lock
			l.mu.Lock()
		}

		if l.layouters == nil {
			l.buf = append(l.buf, msg...)
		} else {
			for _, layouter := range l.layouters {
				layouter.layout(&l.buf, level, msg, now, file, line)
			}
		}

		// Ensure log message ends with a new line
		if len(l.buf) == 0 || l.buf[len(l.buf)-1] != '\n' {
			l.buf = append(l.buf, '\n')
		}

		l.out.Write(l.buf)
	}
}
