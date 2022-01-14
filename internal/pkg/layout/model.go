package layout

import (
	"time"

	"github.com/nszilard/log/models"
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
	Format(buf *[]byte, lev models.Level, msg string, t time.Time, file string, line int)
}
