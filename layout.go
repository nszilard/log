package log

import (
	"time"
)

// Layout used to format log message.
//   %y : year
//   %M : month
//   %d : day
//   %h : hour
//   %m : min
//   %s : second
//   %l : log msg
//   %L : log level
//   %F : file			eg: /a/b/c/d.go
//	 %f : short file	eg: d.go
//   %i : line
//   %D : %y/%M/%d
//   %T : %h:%m:%s
type Layout interface {
	layout(buf *[]byte, lev Level, msg string, t time.Time, file string, line int)
}

var (
	mapLayouts = map[string]Layout{
		"%y": &layoutYear{},
		"%M": &layoutMonth{},
		"%d": &layoutDay{},
		"%h": &layoutHour{},
		"%m": &layoutMinute{},
		"%s": &layoutSecond{},
		"%l": &layoutMsg{},
		"%L": &layoutLevel{},
		"%F": &layoutFile{},
		"%f": &layoutShortFile{},
		"%i": &layoutLine{},
		"%D": &layoutDate{},
		"%T": &layoutTime{},
	}
)

type layoutYear struct{}
type layoutMonth struct{}
type layoutDay struct{}
type layoutHour struct{}
type layoutMinute struct{}
type layoutSecond struct{}
type layoutMsg struct{}
type layoutLevel struct{}
type layoutFile struct{}
type layoutShortFile struct{}
type layoutLine struct{}
type layoutDate struct{}
type layoutTime struct{}
type layoutPlaceholder struct {
	placeholder string
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

func (l *layoutYear) layout(buf *[]byte, lev Level, msg string, t time.Time, file string, line int) {
	itoa(buf, t.Year(), 4)
}

func (l *layoutMonth) layout(buf *[]byte, lev Level, msg string, t time.Time, file string, line int) {
	itoa(buf, int(t.Month()), 2)
}

func (l *layoutDay) layout(buf *[]byte, lev Level, msg string, t time.Time, file string, line int) {
	itoa(buf, t.Day(), 2)
}

func (l *layoutHour) layout(buf *[]byte, lev Level, msg string, t time.Time, file string, line int) {
	itoa(buf, t.Hour(), 2)
}

func (l *layoutMinute) layout(buf *[]byte, lev Level, msg string, t time.Time, file string, line int) {
	itoa(buf, t.Minute(), 2)
}

func (l *layoutSecond) layout(buf *[]byte, lev Level, msg string, t time.Time, file string, line int) {
	itoa(buf, t.Second(), 2)
}

func (l *layoutMsg) layout(buf *[]byte, lev Level, msg string, t time.Time, file string, line int) {
	*buf = append(*buf, msg...)
}

func (l *layoutLevel) layout(buf *[]byte, lev Level, msg string, t time.Time, file string, line int) {
	*buf = append(*buf, '[')
	*buf = append(*buf, Level2Str[lev]...)
	*buf = append(*buf, ']')
}

func (l *layoutFile) layout(buf *[]byte, lev Level, msg string, t time.Time, file string, line int) {
	*buf = append(*buf, file...)
}

func (l *layoutShortFile) layout(buf *[]byte, lev Level, msg string, t time.Time, file string, line int) {
	short := file
	for i := len(file) - 1; i > 0; i-- {
		if file[i] == '/' {
			short = file[i+1:]
			break
		}
	}
	*buf = append(*buf, short...)
}

func (l *layoutLine) layout(buf *[]byte, lev Level, msg string, t time.Time, file string, line int) {
	itoa(buf, line, -1)
}

func (l *layoutDate) layout(buf *[]byte, lev Level, msg string, t time.Time, file string, line int) {
	year, month, day := t.Date()
	itoa(buf, year, 4)
	*buf = append(*buf, '/')
	itoa(buf, int(month), 2)
	*buf = append(*buf, '/')
	itoa(buf, day, 2)
}

func (l *layoutTime) layout(buf *[]byte, lev Level, msg string, t time.Time, file string, line int) {
	hour, min, sec := t.Clock()
	itoa(buf, hour, 2)
	*buf = append(*buf, ':')
	itoa(buf, min, 2)
	*buf = append(*buf, ':')
	itoa(buf, sec, 2)
}

func (l *layoutPlaceholder) layout(buf *[]byte, lev Level, msg string, t time.Time, file string, line int) {
	*buf = append(*buf, l.placeholder...)
}
