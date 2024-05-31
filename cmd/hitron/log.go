package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"
)

type LevelValue slog.Level

var _ flag.Getter = (*LevelValue)(nil)

// String implements flag.Value
func (l LevelValue) String() string {
	switch slog.Level(l) {
	case slog.LevelDebug:
		return "debug"
	case slog.LevelInfo:
		return "info"
	case slog.LevelWarn:
		return "warn"
	case slog.LevelError:
		return "error"
	default:
		return fmt.Sprintf("unknown(%d)", l)
	}
}

func (l *LevelValue) Set(s string) error {
	switch s {
	case "debug":
		*l = LevelValue(slog.LevelDebug)
	case "info":
		*l = LevelValue(slog.LevelInfo)
	case "warn":
		*l = LevelValue(slog.LevelWarn)
	case "error":
		*l = LevelValue(slog.LevelError)
	default:
		return fmt.Errorf("unrecognized log level %q", s)
	}

	return nil
}

func (l LevelValue) Get() any {
	return slog.Level(l)
}

// initLogger -
func initLogger(l LevelValue) {
	handler := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.Level(l)})

	slog.SetDefault(slog.New(handler))
}
