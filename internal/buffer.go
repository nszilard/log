package internal

import "sync"

const (
	smallBufSize  = 512
	mediumBufSize = 4 * 1024
	largeBufSize  = 32 * 1024
	maxBufSize    = 64 * 1024
)

var (
	smallBufPool  = sync.Pool{New: func() any { buf := make([]byte, 0, smallBufSize); return &buf }}
	mediumBufPool = sync.Pool{New: func() any { buf := make([]byte, 0, mediumBufSize); return &buf }}
	largeBufPool  = sync.Pool{New: func() any { buf := make([]byte, 0, largeBufSize); return &buf }}
)

// getBuf returns a buffer from the appropriate pool based on the requested size
func getBuf(size int) *[]byte {
	switch {
	case size <= smallBufSize:
		return smallBufPool.Get().(*[]byte)
	case size <= mediumBufSize:
		return mediumBufPool.Get().(*[]byte)
	default:
		return largeBufPool.Get().(*[]byte)
	}
}

// putBuf returns a buffer to the appropriate pool
func putBuf(buf *[]byte) {
	if cap(*buf) > maxBufSize {
		return
	}
	*buf = (*buf)[:0]
	switch {
	case cap(*buf) <= smallBufSize:
		smallBufPool.Put(buf)
	case cap(*buf) <= mediumBufSize:
		mediumBufPool.Put(buf)
	case cap(*buf) <= largeBufSize:
		largeBufPool.Put(buf)
	}
}
