package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"rpgm-launcher/nwjs-helper/internal/cli"
	"rpgm-launcher/nwjs-helper/internal/dialog"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	os.Exit(cli.Run(ctx, os.Args[1:], os.Stdout, os.Stderr, dialog.Warnings))
}
