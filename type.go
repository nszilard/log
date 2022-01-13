package log

import (
	"io"
	"sync"
)

// Level used to filter log message by the Logger.
type Level uint8

// A Logger represents an active logging object that generates lines of
// output to log listeners. A Logger can be used simultaneously from
// multiple goroutines; it guarantees to serialize access to the Writer.
type Logger struct {
	mu             sync.Mutex // Ensures atomic writes; protects the following fields
	level          Level      // Holds the log level
	layouters      []Layout   // Holds log message layout format
	out            io.Writer  // Destination for the log output
	buf            []byte     // Dor accumulating text to write
	needCallerInfo bool       // Flag of caller info need or not
}
