package dialog

import (
	"os/exec"
	"strings"
)

// Warnings shows launch warnings in a native macOS dialog.
func Warnings(warnings []string) error {
	return warningCommand(warnings).Run()
}

func warningCommand(warnings []string) *exec.Cmd {
	return exec.Command("/usr/bin/osascript", "-e", "on run argv", "-e", `display dialog (item 1 of argv) with title "Game launch warnings" buttons {"OK"} default button 1`, "-e", "end run", "--", strings.Join(warnings, "\n"))
}
