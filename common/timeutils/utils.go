// Package timeutils is a wrapper around the go standard time library.
package timeutils

import (
	"time"
)

// Since returns the duration since t.
func Since(t time.Time) time.Duration {
	return Now().Sub(t)
}

// Until returns the duration until t.
func Until(t time.Time) time.Duration {
	return t.Sub(Now())
}

// Now returns the current local time.
func Now() time.Time {
	return time.Now()
}

func UnixMsec() int64 {
	return time.Now().UnixNano() / 1e6
}

func BeforeYearUnixMsec() int64 {
	return time.Now().AddDate(-1, 0, 0).UnixNano() / 1e6
}