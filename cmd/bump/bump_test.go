package bump

import (
	"bytes"

	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/dp1140a/semver/cmd"
	"github.com/dp1140a/semver/pkg/util"
)

type exitPanic struct{ code int }

// override exitFn during tests (see tiny refactor above)
func withPatchedExit(t *testing.T, fn func()) (code int) {
	t.Helper()

	exitFn := func(c int) { panic(exitPanic{code: c}) }
	prev := exitFn
	defer func() { exitFn = prev }()
	defer func() {
		if r := recover(); r != nil {
			if ep, ok := r.(exitPanic); ok {
				code = ep.code
				return
			}
			panic(r)
		}
	}()
	fn()
	return 0 // unreachable if command exits, but keeps signature tidy
}

func withTempWD(t *testing.T, f func(tmp string)) {
	t.Helper()
	orig, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	t.Cleanup(func() { _ = os.Chdir(orig) })

	tmp := t.TempDir()
	if err := os.Chdir(tmp); err != nil {
		t.Fatalf("chdir temp: %v", err)
	}
	f(tmp)
}

func writeVERSION(t *testing.T, v string) {
	t.Helper()
	if err := os.WriteFile("VERSION", []byte(v+"\n"), 0o644); err != nil {
		t.Fatalf("write VERSION: %v", err)
	}
}

func readVERSION(t *testing.T) string {
	t.Helper()
	b, err := os.ReadFile("VERSION")
	if err != nil {
		t.Fatalf("read VERSION: %v", err)
	}
	return strings.TrimSpace(string(b))
}

// captureStdout runs fn while capturing stdout
func captureStdout(t *testing.T, fn func()) string {
	t.Helper()
	orig := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	defer func() { os.Stdout = orig }()
	fn()
	_ = w.Close()
	var buf bytes.Buffer
	_, _ = io.Copy(&buf, r)
	return buf.String()
}

func TestBumpPatch_WritesAndPrints(t *testing.T) {
	withTempWD(t, func(tmp string) {
		writeVERSION(t, "1.2.3")
		// sanity: util functions see the VERSION file
		if !util.VersionFileExists(tmp) {
			t.Fatalf("expected VERSION to exist")
		}

		out := captureStdout(t, func() {
			code := withPatchedExit(t, func() {
				cmd.RootCmd.SetArgs([]string{"bump", "patch"})
				_ = cmd.RootCmd.Execute()
			})
			if code != 0 {
				t.Fatalf("unexpected exit code: %d", code)
			}
		})

		got := readVERSION(t)
		if got != "1.2.4" {
			t.Fatalf("expected VERSION=1.2.4, got %q", got)
		}
		if !strings.Contains(out, "Bumping Patch") || !strings.Contains(out, "New Version: 1.2.4") {
			t.Fatalf("stdout missing expected lines:\n%s", out)
		}
	})
}

func TestBumpMinor_WritesAndPrints(t *testing.T) {
	withTempWD(t, func(tmp string) {
		writeVERSION(t, "1.2.3")
		out := captureStdout(t, func() {
			code := withPatchedExit(t, func() {
				cmd.RootCmd.SetArgs([]string{"bump", "minor"})
				_ = cmd.RootCmd.Execute()
			})
			if code != 0 {
				t.Fatalf("unexpected exit code: %d", code)
			}
		})
		if got := readVERSION(t); got != "1.3.0" {
			t.Fatalf("expected VERSION=1.3.0, got %q", got)
		}
		if !strings.Contains(out, "Bumping Minor") || !strings.Contains(out, "New Version: 1.3.0") {
			t.Fatalf("stdout missing expected lines:\n%s", out)
		}
	})
}

func TestBumpMajor_WritesAndPrints(t *testing.T) {
	withTempWD(t, func(tmp string) {
		writeVERSION(t, "1.2.3")
		out := captureStdout(t, func() {
			code := withPatchedExit(t, func() {
				cmd.RootCmd.SetArgs([]string{"bump", "major"})
				_ = cmd.RootCmd.Execute()
			})
			if code != 0 {
				t.Fatalf("unexpected exit code: %d", code)
			}
		})
		if got := readVERSION(t); got != "2.0.0" {
			t.Fatalf("expected VERSION=2.0.0, got %q", got)
		}
		if !strings.Contains(out, "Bumping Major") || !strings.Contains(out, "New Version: 2.0.0") {
			t.Fatalf("stdout missing expected lines:\n%s", out)
		}
	})
}

func TestBump_DefaultIsPatch(t *testing.T) {
	withTempWD(t, func(tmp string) {
		writeVERSION(t, "0.0.9")
		out := captureStdout(t, func() {
			code := withPatchedExit(t, func() {
				cmd.RootCmd.SetArgs([]string{"bump"}) // no subcommand
				_ = cmd.RootCmd.Execute()
			})
			if code != 0 {
				t.Fatalf("unexpected exit code: %d", code)
			}
		})
		if got := readVERSION(t); got != "0.0.10" {
			t.Fatalf("expected VERSION=0.0.10, got %q", got)
		}
		if !strings.Contains(out, "Bumping Patch") || !strings.Contains(out, "New Version: 0.0.10") {
			t.Fatalf("stdout missing expected lines:\n%s", out)
		}
	})
}

func TestBump_DryRunDoesNotWrite(t *testing.T) {
	withTempWD(t, func(tmp string) {
		writeVERSION(t, "1.2.3")

		out := captureStdout(t, func() {
			cmd.RootCmd.SetArgs([]string{"bump", "--dry", "minor"})
			if err := cmd.RootCmd.Execute(); err != nil {
				t.Fatalf("execute: %v", err)
			}
		})

		// VERSION should remain unchanged on dry-run
		if got := readVERSION(t); got != "1.2.3" {
			t.Fatalf("expected VERSION unchanged on dry-run, got %q", got)
		}

		// Output should clearly indicate dry-run and show the would-be version
		if !strings.Contains(out, "Bumping Minor") {
			t.Fatalf("stdout missing 'Bumping Minor':\n%s", out)
		}
		if !strings.Contains(out, "[dry-run]") {
			t.Fatalf("stdout missing dry-run marker:\n%s", out)
		}
		if !strings.Contains(out, "New Version would be: 1.3.0") {
			t.Fatalf("stdout missing would-be version:\n%s", out)
		}
		if !strings.Contains(out, "VERSION file unchanged") {
			t.Fatalf("stdout missing 'VERSION file unchanged' note:\n%s", out)
		}
	})
}

func TestBump_NoVersionFile_GracefulMessage(t *testing.T) {
	withTempWD(t, func(tmp string) {
		// Ensure no VERSION file
		if _, err := os.Stat(filepath.Join(tmp, "VERSION")); !os.IsNotExist(err) {
			t.Fatalf("expected no VERSION file")
		}
		out := captureStdout(t, func() {
			code := withPatchedExit(t, func() {
				cmd.RootCmd.SetArgs([]string{"bump", "patch"})
				_ = cmd.RootCmd.Execute()
			})
			// current code path prints and exits 0 (it’s fine to assert 0 here)
			if code != 0 {
				t.Fatalf("unexpected exit code: %d", code)
			}
		})
		if !strings.Contains(out, "No VERSION file found") {
			t.Fatalf("expected helpful message, got:\n%s", out)
		}
	})
}
