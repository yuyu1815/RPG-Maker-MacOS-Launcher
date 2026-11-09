package main

import (
	"flag"
	"os"
	"os/exec"
	"strings"
	"testing"
)

func TestMainDispatch(t *testing.T) {
	if os.Getenv("LAUNCHER_DISPATCH_TEST") == "1" {
		for i, arg := range os.Args {
			if arg == "--" {
				os.Args = append([]string{"launcher"}, os.Args[i+1:]...)
				flag.CommandLine = flag.NewFlagSet("launcher", flag.ExitOnError)
				main()
				return
			}
		}
		t.Fatal("missing argument separator")
	}
	for _, tc := range []struct {
		name string
		args []string
		want string
		code int
	}{
		{"gui help", []string{"--help"}, "-webui", 0},
		{"gui flags", []string{"--root", "/root", "--webui=/web", "--resources", "/resources", "--help"}, "-webui", 0},
		{"gui flag value launch", []string{"--root", "launch", "--help"}, "-webui", 0},
		{"cli help", []string{"launch", "--help"}, "-game", 0},
		{"cli invalid", []string{"launch", "--root", "/root"}, "--game and --resources must be absolute paths", 2},
		{"unknown flag", []string{"--game", "/game"}, "flag provided but not defined", 2},
	} {
		t.Run(tc.name, func(t *testing.T) {
			args := append([]string{"-test.run=^TestMainDispatch$", "--"}, tc.args...)
			cmd := exec.Command(os.Args[0], args...)
			cmd.Env = append(os.Environ(), "LAUNCHER_DISPATCH_TEST=1")
			output, err := cmd.CombinedOutput()
			code := 0
			if err != nil {
				if exit, ok := err.(*exec.ExitError); ok {
					code = exit.ExitCode()
				} else {
					t.Fatal(err)
				}
			}
			if code != tc.code || !strings.Contains(string(output), tc.want) {
				t.Fatalf("exit %d, output %q; want exit %d and %q", code, output, tc.code, tc.want)
			}
		})
	}
}
