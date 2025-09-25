package log

import (
	"time"
	"unsafe"

	"github.com/nszilard/log/internal"
)

func logStructured(level Level, fields []Data) {
	if level <= std.currentLevel && len(fields) > 0 {
		allTyped := true
		for _, f := range fields {
			if f.Type == UnknownType {
				allTyped = false
				break
			}
		}

		if allTyped {
			std.internal.LogStructuredTypedWithFileInfo(level.String(), std.includeFileInfo, *(*[]internal.Data)(unsafe.Pointer(&fields)))
		} else {
			kv := make([]any, 0, len(fields)*2)
			for _, f := range fields {
				kv = append(kv, f.Key)

				switch f.Type {
				case StringType:
					kv = append(kv, f.String)
				case IntType:
					kv = append(kv, f.Integer)
				case FloatType:
					kv = append(kv, f.Float)
				case BoolType:
					kv = append(kv, f.Bool)
				case ErrorType:
					if f.Interface == nil {
						kv = append(kv, (*error)(nil))
					} else {
						kv = append(kv, f.String)
					}
				case DurationType:
					kv = append(kv, time.Duration(f.Integer))
				case TimeType:
					kv = append(kv, f.Interface)
				default:
					kv = append(kv, f.Interface)
				}
			}
			std.internal.LogStructuredWithFileInfo(level.String(), std.includeFileInfo, kv...)
		}
	}
}
