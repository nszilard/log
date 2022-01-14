// Package logger provides the logging object
package logger

import (
	"fmt"
	"io"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/nszilard/log/internal/pkg/layout"
	"github.com/nszilard/log/models"
)

// A Logger represents an active logging object that generates lines of
// output to log listeners. A Logger can be used simultaneously from
// multiple goroutines; it guarantees to serialize access to the Writer.
type Logger struct {
	mu             sync.Mutex      // Ensures atomic writes; protects the following fields
	Level          models.Level    // Holds the log level
	LayoutFormats  []layout.Layout // Holds log message layout format
	Out            io.Writer       // Destination for the log output
	Buf            []byte          // Dor accumulating text to write
	NeedCallerInfo bool            // Flag of caller info need or not
	Depth          int             // Determines how far back in the stacktrace to look
}

// SetLevel set the log level for Logger.
func (l *Logger) SetLevel(lev models.Level) {
	if lev != l.Level {
		l.mu.Lock()
		l.Level = lev
		l.mu.Unlock()
	}
}

// SetLayout set the layout of log message.
func (l *Logger) SetLayout(logLayout string) {
	layoutFlags := logLayout

	l.mu.Lock()
	defer l.mu.Unlock()

	// clear before set.
	l.LayoutFormats = nil

	for {
		i := strings.IndexByte(layoutFlags, '%')
		if i != -1 {
			if i != 0 {
				l.LayoutFormats = append(l.LayoutFormats, &layout.Placeholder{
					Placeholder: layoutFlags[:i],
				})
			}

			if i+2 > len(layoutFlags) {
				break
			}

			f := layoutFlags[i : i+2]
			layouter := layout.Map(f)
			if layouter != nil {
				l.LayoutFormats = append(l.LayoutFormats, layouter)
				switch layouter.(type) {
				case *layout.File, *layout.ShortFile, *layout.Line:
					l.NeedCallerInfo = true
				}
			} else {
				l.LayoutFormats = append(l.LayoutFormats, &layout.Placeholder{
					Placeholder: f,
				})
			}

			layoutFlags = layoutFlags[i+2:]
		} else {
			if len(layoutFlags) > 0 {
				l.LayoutFormats = append(l.LayoutFormats, &layout.Placeholder{
					Placeholder: layoutFlags,
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
	l.Out = w
}

// Panic print a PanicLevel message to the logger.
// Arguments are handled in the manner of fmt.Print.
func (l *Logger) Panic(v ...interface{}) {
	s := fmt.Sprint(v...)
	l.Log(models.PanicLevel, s)
}

// Panicf print a PanicLevel message to the logger.
// Arguments are handled in the manner of fmt.Printf.
func (l *Logger) Panicf(format string, v ...interface{}) {
	s := fmt.Sprintf(format, v...)
	l.Log(models.PanicLevel, s)
}

// Fatal print a FatalLevel message to the logger.
// Arguments are handled in the manner of fmt.Print.
func (l *Logger) Fatal(v ...interface{}) {
	l.Log(models.FatalLevel, fmt.Sprint(v...))
}

// Fatalf print a FatalLevel message to the logger.
// Arguments are handled in the manner of fmt.Printf.
func (l *Logger) Fatalf(format string, v ...interface{}) {
	l.Log(models.FatalLevel, fmt.Sprintf(format, v...))
}

// Error print an ErrorLevel message to the logger.
// Arguments are handled in the manner of fmt.Print.
func (l *Logger) Error(v ...interface{}) {
	l.Log(models.ErrorLevel, fmt.Sprint(v...))
}

// Errorf print an ErrorLevel message to the logger.
// Arguments are handled in the manner of fmt.Printf.
func (l *Logger) Errorf(format string, v ...interface{}) {
	l.Log(models.ErrorLevel, fmt.Sprintf(format, v...))
}

// Warn print a WarnLevel message to the logger.
// Arguments are handled in the manner of fmt.Print.
func (l *Logger) Warn(v ...interface{}) {
	l.Log(models.WarnLevel, fmt.Sprint(v...))
}

// Warnf print a WarnLevel message to the logger.
// Arguments are handled in the manner of fmt.Printf.
func (l *Logger) Warnf(format string, v ...interface{}) {
	l.Log(models.WarnLevel, fmt.Sprintf(format, v...))
}

// Info print an InfoLevel message to the logger.
// Arguments are handled in the manner of fmt.Print.
func (l *Logger) Info(v ...interface{}) {
	l.Log(models.InfoLevel, fmt.Sprint(v...))
}

// Infof print an InfoLevel message to the logger.
// Arguments are handled in the manner of fmt.Printf.
func (l *Logger) Infof(format string, v ...interface{}) {
	l.Log(models.InfoLevel, fmt.Sprintf(format, v...))
}

// Println print a NoLevel message to the logger.
// Arguments are handled in the manner of fmt.Print.
func (l *Logger) Println(v ...interface{}) {
	l.Log(models.NoLevel, fmt.Sprint(v...))
}

// Printf print a NoLevel message to the logger.
// Arguments are handled in the manner of fmt.Printf.
func (l *Logger) Printf(format string, v ...interface{}) {
	l.Log(models.NoLevel, fmt.Sprintf(format, v...))
}

// Trace print a TraceLevel message to the logger.
// Arguments are handled in the manner of fmt.Print.
func (l *Logger) Trace(v ...interface{}) {
	l.Log(models.TraceLevel, fmt.Sprint(v...))
}

// Tracef print a TraceLevel message to the logger.
// Arguments are handled in the manner of fmt.Printf.
func (l *Logger) Tracef(format string, v ...interface{}) {
	l.Log(models.TraceLevel, fmt.Sprintf(format, v...))
}

// Debug print a DebugLevel message to the logger.
// Arguments are handled in the manner of fmt.Print.
func (l *Logger) Debug(v ...interface{}) {
	l.Log(models.DebugLevel, fmt.Sprint(v...))
}

// Debugf print a DebugLevel message to the logger.
// Arguments are handled in the manner of fmt.Printf.
func (l *Logger) Debugf(format string, v ...interface{}) {
	l.Log(models.DebugLevel, fmt.Sprintf(format, v...))
}

// Log print a leveled message to the logger.
// Arguments are handled in the manner of fmt.Print.
func (l *Logger) Log(level models.Level, msg string) {
	now := time.Now() // get this early.
	var file string
	var line int
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.Level >= level {
		l.Buf = l.Buf[:0]

		if l.NeedCallerInfo {
			// release lock while getting caller info - it's expensive.
			l.mu.Unlock()
			var ok bool
			_, file, line, ok = runtime.Caller(l.Depth)
			if !ok {
				file = "???"
				line = 0
			}
			// re-lock
			l.mu.Lock()
		}

		if l.LayoutFormats == nil {
			l.Buf = append(l.Buf, msg...)
		} else {
			for _, layouter := range l.LayoutFormats {
				layouter.Format(&l.Buf, level, msg, now, file, line)
			}
		}

		// Ensure log message ends with a new line
		if len(l.Buf) == 0 || l.Buf[len(l.Buf)-1] != '\n' {
			l.Buf = append(l.Buf, '\n')
		}

		l.Out.Write(l.Buf)
	}
}
