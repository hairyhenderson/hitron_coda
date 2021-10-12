package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	hitron "github.com/hairyhenderson/hitron_coda"
)

type debugLogAdapter struct {
	log.Logger
}

func (l debugLogAdapter) Logf(format string, args ...interface{}) {
	_ = l.Log("msg", fmt.Sprintf(format, args...))
}

type opts struct {
	host     string
	username string
	password string
	logLevel Level
}

func flags(args []string, o *opts) ([]string, error) {
	prog := args[0]

	fs := flag.NewFlagSet("root", flag.ContinueOnError)

	// default to info
	_ = o.logLevel.Set("info")

	fs.StringVar(&o.host, "host", "192.168.0.1", "hostname or IP address of the cable modem")
	fs.StringVar(&o.username, "username", "cusadmin", "username for the cable modem")
	fs.StringVar(&o.password, "password", os.Getenv("HITRON_CODA_PASSWORD"), "password for the cable modem")
	fs.Var(&o.logLevel, "log.level", "log messages with the given severity or above. Valid levels: [debug, info, warn, error]")

	fs.Usage = func() {
		fmt.Fprintf(fs.Output(), `usage: %s <flags> <command> ...
		flags:
		`, prog)
		fs.PrintDefaults()
		fmt.Fprintf(fs.Output(), `
		commands:
		cm <flags>
		Cable Modem subcommands
		router <flags>
		Router subcommands
		
		Run %s <command> -h for more information
		`, prog)
	}

	err := fs.Parse(args[1:])
	if err != nil {
		return nil, err
	}

	fsArgs := fs.Args()
	if len(fsArgs) == 0 {
		fs.Usage()

		return nil, fmt.Errorf("no subcommand specified")
	}

	return fsArgs, err
}

func main() {
	err := run(os.Args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run(args []string) error {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)
	defer stop()

	o := &opts{}

	fsArgs, err := flags(args, o)
	if err != nil {
		return err
	}

	logger := NewLogger(o.logLevel)
	debugLogger := debugLogAdapter{level.Debug(logger)}
	ctx = hitron.ContextWithDebugLogger(ctx, debugLogger)

	cm, err := hitron.New(o.host, o.username, o.password)
	if err != nil {
		return err
	}

	switch fsArgs[0] {
	case "cm":
		return cmdCM(ctx, cm, flag.NewFlagSet("cm", flag.ExitOnError), fsArgs[1:])
	case "router":
		return cmdRouter(ctx, cm, flag.NewFlagSet("router", flag.ExitOnError), fsArgs[1:])
	default:
		return fmt.Errorf("invalid subcommand %q", fsArgs[0])
	}
}
