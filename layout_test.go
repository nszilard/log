package log

import (
	"fmt"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestLayout(t *testing.T) {
	now := time.Now()
	compareInt(t, &layoutYear{}, now, now.Year())
	compareInt(t, &layoutMonth{}, now, int(now.Month()))
	compareInt(t, &layoutDay{}, now, now.Day())
	compareInt(t, &layoutHour{}, now, now.Hour())
	compareInt(t, &layoutMinute{}, now, now.Minute())
	compareInt(t, &layoutSecond{}, now, now.Second())

	dateStr := fmt.Sprintf("%04d", int(now.Year())) + "/" + fmt.Sprintf("%02d", int(now.Month())) + "/" + fmt.Sprintf("%02d", int(now.Day()))
	timeStr := fmt.Sprintf("%02d", int(now.Hour())) + ":" + fmt.Sprintf("%02d", int(now.Minute())) + ":" + fmt.Sprintf("%02d", int(now.Second()))
	compareStringDate(t, &layoutDate{}, now, dateStr)
	compareStringDate(t, &layoutTime{}, now, timeStr)

	compareString(t, &layoutFile{}, now, "test/FileLayout.go", "test/FileLayout.go", 0)
	compareString(t, &layoutShortFile{}, now, "testFileLayout.go", "testFileLayout.go", 32)
}

// ------------------------------------------
// Benchmark
// ------------------------------------------

func BenchmarkItoa(b *testing.B) {
	dst := make([]byte, 0, 64)
	for i := 0; i < b.N; i++ {
		dst = dst[0:0]
		itoa(&dst, 2015, 4)   // year
		itoa(&dst, 1, 2)      // month
		itoa(&dst, 30, 2)     // day
		itoa(&dst, 12, 2)     // hour
		itoa(&dst, 56, 2)     // minute
		itoa(&dst, 0, 2)      // second
		itoa(&dst, 987654, 6) // microsecond
	}
}

// ------------------------------------------
// Helpers
// ------------------------------------------

func compareInt(t *testing.T, layout Layout, time time.Time, rv int) {
	var buf []byte
	layout.layout(&buf, DebugLevel, "", time, "", 0)

	lv, err := strconv.Atoi(string(buf))
	if err != nil {
		t.Error(err)
		return
	}
	if lv != rv {
		t.Errorf("layout[%T] failed! expected: %v, got: %v", layout, rv, lv)
	}
}

func compareStringDate(t *testing.T, layout Layout, time time.Time, rv string) {
	var buf []byte
	layout.layout(&buf, DebugLevel, "", time, "", 0)
	if strings.Compare(string(buf), rv) != 0 {
		t.Errorf("layout[%T] failed! expected: %q, got: %q", layout, rv, string(buf))
	}
}

func compareString(t *testing.T, layout Layout, time time.Time, input, expected string, lineNumber int) {
	var buf []byte
	layout.layout(&buf, DebugLevel, "", time, input, lineNumber)
	if strings.Compare(string(buf), input) != 0 {
		t.Errorf("layout[%T] failed! expected: %q, got: %q", layout, expected, string(buf))
	}
}
