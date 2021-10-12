package main

import (
	"context"
	"flag"
	"fmt"

	hitron "github.com/hairyhenderson/hitron_coda"
)

func cmdCM(ctx context.Context, cm *hitron.CableModem, f *flag.FlagSet, argv []string) error {
	f.Usage = func() {
		f.PrintDefaults()
		fmt.Fprintf(f.Output(), `
subcommands:
	version
		Print cable modem version information
	reboot
		Reboot the cable modem
	log
		Print cable modem logs
	clearLog
		Clear cable modem logs
`)
	}

	_ = f.Parse(argv)

	args := f.Args()
	if len(args) == 0 {
		f.Usage()

		return nil
	}

	cmds := map[string]func(ctx context.Context) (fmt.Stringer, error){
		"version": func(ctx context.Context) (fmt.Stringer, error) {
			return cm.CMVersion(ctx)
		},
		"reboot": func(ctx context.Context) (fmt.Stringer, error) {
			return cm.CMReboot(ctx)
		},
		"log": func(ctx context.Context) (fmt.Stringer, error) {
			return cm.CMLog(ctx)
		},
		"clearLog": func(ctx context.Context) (fmt.Stringer, error) {
			return cm.CMClearLog(ctx)
		},
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
