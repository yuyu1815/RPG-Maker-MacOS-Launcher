package launcher

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestSharedExport(t *testing.T) {
	for _, engine := range []GameEngine{NWJS, EasyRPG} {
		t.Run(string(engine), func(t *testing.T) {
			m := engineContractManager(t)
			g := engineContractGame(t, engine)
			s := Settings{LastSelectedVersion: "v1.2.3"}
			engineContractRuntime(t, m, engine)
			engineContractWrite(t, filepath.Join(m.ResourcesDir, "rpgm-cli"), []byte("#!/bin/sh\nprintf '%s\\n' \"$@\" > \"$TEST_ARGS\"\n"), 0700)
			out, err := m.ExportShared(context.Background(), g, s, t.TempDir(), "Game", nil)
			if err != nil {
				t.Fatal(err)
			}
			if _, err := os.Stat(filepath.Join(out, "Contents/Resources/game/save/keep")); err != nil {
				t.Fatal(err)
			}
			if _, err := os.Stat(filepath.Join(out, "Contents/Resources/EasyRPG Player.app")); !os.IsNotExist(err) {
				t.Fatal("runtime was copied", err)
			}
			args := filepath.Join(t.TempDir(), "args")
			cmd := exec.Command(filepath.Join(out, "Contents/MacOS/launch"))
			cmd.Env = append(os.Environ(), "TEST_ARGS="+args)
			if b, err := cmd.CombinedOutput(); err != nil {
				t.Fatalf("%s: %v", b, err)
			}
			b, err := os.ReadFile(args)
			if err != nil {
				t.Fatal(err)
			}
			for _, value := range []string{"launch", m.Root, filepath.Join(out, "Contents/Resources/game"), "--wait", "--dialog-warnings", "v1.2.3"} {
				if !strings.Contains(string(b), value+"\n") {
					t.Fatalf("missing %q in %s", value, b)
				}
			}
		})
	}
}

func TestSharedExportRejectsMKXPAndMissingRuntime(t *testing.T) {
	m := engineContractManager(t)
	g := engineContractGame(t, MKXPZ)
	engineContractRuntime(t, m, MKXPZ)
	if _, err := m.ExportShared(context.Background(), g, Settings{}, t.TempDir(), "Game", nil); err == nil || !strings.Contains(err.Error(), "MKXP") {
		t.Fatalf("%v", err)
	}
	g = engineContractGame(t, EasyRPG)
	if _, err := m.ExportShared(context.Background(), g, Settings{}, t.TempDir(), "Game", nil); err == nil {
		t.Fatal("missing runtime accepted")
	}
}
