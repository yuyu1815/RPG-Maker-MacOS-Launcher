package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"rpgm-launcher/nwjs-helper/internal/cli"
	"rpgm-launcher/nwjs-helper/internal/dialog"
	"rpgm-launcher/nwjs-helper/internal/launcher"
	"rpgm-launcher/nwjs-helper/internal/ui"
)

func main() {
	if len(os.Args) > 1 && os.Args[1] == "launch" {
		ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
		code := cli.Run(ctx, os.Args[1:], os.Stdout, os.Stderr, dialog.Warnings)
		stop()
		os.Exit(code)
	}
	web := flag.String("webui", "", "absolute source webui directory (ignored in .app)")
	resources := flag.String("resources", "", "absolute source JavaScript resources directory (ignored in .app)")
	root := flag.String("root", "", "absolute managed runtime/settings root (default: Application Support/RPGM-Launcher)")
	flag.Parse()
	if flag.NArg() != 0 {
		log.Fatal("unexpected positional arguments")
	}
	if err := launcher.EnsureSupportedHost(runtime.GOOS, runtime.GOARCH); err != nil {
		log.Fatal(err)
	}
	exe, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	paths, err := ui.ResolvePaths(exe, *web, *resources)
	if err != nil {
		log.Fatal(err)
	}
	if *root == "" {
		home, e := os.UserHomeDir()
		if e != nil {
			log.Fatal(e)
		}
		*root = launcher.DefaultRoot(home)
	}
	service, err := ui.New(*root, paths.Resources)
	if err != nil {
		log.Fatal(err)
	}
	if err = ui.RunApplication(service, paths); err != nil {
		log.Fatal(err)
	}
}
