// Package layout defines various layouts to use when formatting a log message
package layout

import (
	"time"

	"github.com/nszilard/log/models"
)

var (
	mapLayouts = map[string]Layout{
		"%y": &Year{},
		"%M": &Month{},
		"%d": &Day{},
		"%h": &Hour{},
		"%m": &Minute{},
		"%s": &Second{},
		"%l": &Msg{},
		"%L": &Level{},
		"%F": &File{},
		"%f": &ShortFile{},
		"%i": &Line{},
		"%D": &Date{},
		"%T": &Time{},
	}
)

// Map maps a string to a layout
func Map(s string) Layout {
	return mapLayouts[s]
}

type Year struct{}
type Month struct{}
type Day struct{}
type Hour struct{}
type Minute struct{}
type Second struct{}
type Msg struct{}
type Level struct{}
type File struct{}
type ShortFile struct{}
type Line struct{}
type Date struct{}
type Time struct{}
type Placeholder struct {
	Placeholder string
}

// Cheap integer to fixed-width decimal ASCII. Give a negative width to avoid zero-padding.
func itoa(buf *[]byte, i int, wid int) {
	// Assemble decimal in reverse order.
	var b [20]byte
	bp := len(b) - 1
	for i >= 10 || wid > 1 {
		wid--
		q := i / 10
		b[bp] = byte('0' + i - q*10)
		bp--
		i = q
	}
	// i < 10
	b[bp] = byte('0' + i)
	*buf = append(*buf, b[bp:]...)
}

func (l *Year) Format(buf *[]byte, lev models.Level, msg string, t time.Time, file string, line int) {
	itoa(buf, t.Year(), 4)
}

func (l *Month) Format(buf *[]byte, lev models.Level, msg string, t time.Time, file string, line int) {
	itoa(buf, int(t.Month()), 2)
}

func (l *Day) Format(buf *[]byte, lev models.Level, msg string, t time.Time, file string, line int) {
	itoa(buf, t.Day(), 2)
}

func (l *Hour) Format(buf *[]byte, lev models.Level, msg string, t time.Time, file string, line int) {
	itoa(buf, t.Hour(), 2)
}

func (l *Minute) Format(buf *[]byte, lev models.Level, msg string, t time.Time, file string, line int) {
	itoa(buf, t.Minute(), 2)
}

func (l *Second) Format(buf *[]byte, lev models.Level, msg string, t time.Time, file string, line int) {
	itoa(buf, t.Second(), 2)
}

func (l *Msg) Format(buf *[]byte, lev models.Level, msg string, t time.Time, file string, line int) {
	*buf = append(*buf, msg...)
}

func (l *Level) Format(buf *[]byte, lev models.Level, msg string, t time.Time, file string, line int) {
	*buf = append(*buf, '[')
	*buf = append(*buf, models.LevelToString(lev)...)
	*buf = append(*buf, ']')
}

func (l *File) Format(buf *[]byte, lev models.Level, msg string, t time.Time, file string, line int) {
	*buf = append(*buf, file...)
}

func (l *ShortFile) Format(buf *[]byte, lev models.Level, msg string, t time.Time, file string, line int) {
	short := file
	for i := len(file) - 1; i > 0; i-- {
		if file[i] == '/' {
			short = file[i+1:]
			break
		}
	}
	*buf = append(*buf, short...)
}

func (l *Line) Format(buf *[]byte, lev models.Level, msg string, t time.Time, file string, line int) {
	itoa(buf, line, -1)
}

func (l *Date) Format(buf *[]byte, lev models.Level, msg string, t time.Time, file string, line int) {
	year, month, day := t.Date()
	itoa(buf, year, 4)
	*buf = append(*buf, '/')
	itoa(buf, int(month), 2)
	*buf = append(*buf, '/')
	itoa(buf, day, 2)
}

func (l *Time) Format(buf *[]byte, lev models.Level, msg string, t time.Time, file string, line int) {
	hour, min, sec := t.Clock()
	itoa(buf, hour, 2)
	*buf = append(*buf, ':')
	itoa(buf, min, 2)
	*buf = append(*buf, ':')
	itoa(buf, sec, 2)
}

func (l *Placeholder) Format(buf *[]byte, lev models.Level, msg string, t time.Time, file string, line int) {
	*buf = append(*buf, l.Placeholder...)
}
