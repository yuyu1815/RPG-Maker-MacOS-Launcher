package dialog

import (
	"strings"
	"testing"
)

func TestWarningsPassedAsData(t *testing.T) {
	warnings := []string{"protection not applied", "quote \"; $(touch /tmp/evil)\nsecond line"}
	args := warningCommand(warnings).Args
	if args[len(args)-2] != "--" || args[len(args)-1] != strings.Join(warnings, "\n") {
		t.Fatalf("warnings not passed as data: %q", args)
	}
}
