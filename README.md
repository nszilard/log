# Log

Log is a flexible logger for Golang applications. Inspired by [xfxdev/xlog](https://github.com/xfxdev/xlog).

> ***NOTICE*** Neither `Fatal(f)` nor `Panic(f)` functions exit the running code.

## Installation

``` shell
go get github.com/nszilard/log
```

## Features

### • Level logging

Support 7 different log levels, with an additional one having no level set. This makes it super easy to migrate from the built-in `"log"` package.

``` go
const (
	PanicLevel Level = iota
	FatalLevel
	ErrorLevel
	WarnLevel
	NoLevel
	InfoLevel
	TraceLevel
	DebugLevel
)
```

### • Custom Log Layout

The layout can be modified according to the following table (you can use any combination):

| Key  | Maps to      | Example     |
| ----:|:------------ |:-----------:|
| `%y` | Year         |             |
| `%M` | Month        |             |
| `%d` | Day          |             |
| `%h` | Hours        |             |
| `%m` | Minutes      |             |
| `%s` | Seconds      |             |
| `%l` | Log message  |             |
| `%L` | Log level    |             |
| `%F` | File path    | /a/b/c/d.go |
| `%f` | File name    | d.go        |
| `%i` | Line in file |             |
| `%D` | `%y/%M/%d`   |             |
| `%T` | `%h:%m:%s`   |             |

> ***NOTICE:*** If you doesn't call 'SetLayout', it will use `%D %T %L (%f:%i) ▶ %l` by default.

### • Thread safety

Log is protected by a mutex, so you can output logs in multiple goroutines.