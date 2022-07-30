package main

import (
	"context"
	"flag"
	"fmt"

	hitron "github.com/hairyhenderson/hitron_coda"
)

func cmdRouter(ctx context.Context, cm *hitron.CableModem, f *flag.FlagSet, argv []string) error {
	f.Usage = func() {
		f.PrintDefaults()
		fmt.Fprintf(f.Output(), `
subcommands:
	capability
		Print cable modem capabilities
	location
		Print cable modem location
	sysInfo
		Print cable modem router system information
`)
	}

	_ = f.Parse(argv)

	args := f.Args()
	if len(args) == 0 {
		f.Usage()

		return nil
	}

	cmds := map[string]func(ctx context.Context) (fmt.Stringer, error){
		"capability": func(ctx context.Context) (fmt.Stringer, error) { return cm.RouterCapability(ctx) },
		"location":   func(ctx context.Context) (fmt.Stringer, error) { return cm.RouterLocation(ctx) },
		"sysInfo":    func(ctx context.Context) (fmt.Stringer, error) { return cm.RouterSysInfo(ctx) },
	}

	c, ok := cmds[args[0]]
	if !ok {
		return fmt.Errorf("unknown subcommand: %s", args[0])
	}

	if err := cm.Login(ctx); err != nil {
		return err
	}

	defer func() { _ = cm.Logout(ctx) }()

	if c != nil {
		out, err := c(ctx)
		fmt.Printf("%s", out)

		return err
	}

	return nil
}
