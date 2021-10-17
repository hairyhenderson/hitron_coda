package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
)

type Level struct {
	lev level.Option
	s   string
}

var _ flag.Value = (*Level)(nil)

// String implements flag.Value
func (l *Level) String() string {
	return l.s
}

// Set implements flag.Value
func (l *Level) Set(s string) error {
	l.s = s

	switch s {
	case "debug":
		l.lev = level.AllowDebug()
	case "info":
		l.lev = level.AllowInfo()
	case "warn":
		l.lev = level.AllowWarn()
	case "error":
		l.lev = level.AllowError()
	default:
		return fmt.Errorf("unrecognized log level %q", s)
	}

	return nil
}

// NewLogger -
func NewLogger(l Level) log.Logger {
	logger := log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))

	// add common labels (displayed in order)
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)
	logger = level.NewFilter(logger, l.lev)
	logger = log.With(logger, "caller", log.DefaultCaller)

	return logger
}
