package sl

import (
	"log/slog"
	"time"
)

func Err(err error) slog.Attr {
	return slog.Attr{
		Key:   "error",
		Value: slog.StringValue(err.Error()),
	}
}

func Duration(key string, d time.Duration) slog.Attr {
	return slog.Attr{
		Key:   key,
		Value: slog.StringValue(d.String()),
	}
}

// func Int64(key string, v int64) slog.Attr {
// 	return slog.Attr{
// 		Key:   key,
// 		Value: slog.Int64Value(v),
// 	}
// }
