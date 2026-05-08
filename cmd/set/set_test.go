package set

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/dp1140a/semver/cmd"
)

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

func TestSetVersion_WritesAndPrints(t *testing.T) {
	withTempWD(t, func(tmp string) {
		writeVERSION(t, "1.2.3")
		out := captureStdout(t, func() {
			cmd.RootCmd.SetArgs([]string{"set", "2.0.0"})
			if err := cmd.RootCmd.Execute(); err != nil {
				t.Fatalf("execute: %v", err)
			}
		})
		if got := readVERSION(t); got != "2.0.0" {
			t.Fatalf("expected VERSION=2.0.0, got %q", got)
		}
		if !strings.Contains(out, "Current Version: 1.2.3") ||
			!strings.Contains(out, "Setting Version") ||
			!strings.Contains(out, "New Version: 2.0.0") {
			t.Fatalf("unexpected stdout:\n%s", out)
		}
	})
}

func TestSetVersion_DryRun(t *testing.T) {
	withTempWD(t, func(tmp string) {
		writeVERSION(t, "1.2.3")
		out := captureStdout(t, func() {
			cmd.RootCmd.SetArgs([]string{"set", "--dry", "2.1.0"})
			if err := cmd.RootCmd.Execute(); err != nil {
				t.Fatalf("execute: %v", err)
			}
		})
		// Unchanged
		if got := readVERSION(t); got != "1.2.3" {
			t.Fatalf("expected VERSION unchanged, got %q", got)
		}
		if !strings.Contains(out, "Current Version: 1.2.3") ||
			!strings.Contains(out, "Setting Version") ||
			!strings.Contains(out, "[dry-run] New Version would be: 2.1.0") ||
			!strings.Contains(out, "VERSION file unchanged") {
			t.Fatalf("unexpected stdout:\n%s", out)
		}
	})
}

func TestSetVersion_PreservesExistingVPrefixWhenArgOmitsIt(t *testing.T) {
	withTempWD(t, func(tmp string) {
		writeVERSION(t, "v1.0.0-alpha.3")
		out := captureStdout(t, func() {
			cmd.RootCmd.SetArgs([]string{"set", "--dry=false", "1.0.0-beta.1"})
			if err := cmd.RootCmd.Execute(); err != nil {
				t.Fatalf("execute: %v", err)
			}
		})
		if got := readVERSION(t); got != "v1.0.0-beta.1" {
			t.Fatalf("expected VERSION=v1.0.0-beta.1, got %q", got)
		}
		if !strings.Contains(out, "New Version: v1.0.0-beta.1") {
			t.Fatalf("unexpected stdout:\n%s", out)
		}
	})
}

func TestSetVersion_ExplicitPrefixCanSwitchStyle(t *testing.T) {
	withTempWD(t, func(tmp string) {
		writeVERSION(t, "1.0.0-rc.1")
		out := captureStdout(t, func() {
			cmd.RootCmd.SetArgs([]string{"set", "--dry=false", "v1.0.0"})
			if err := cmd.RootCmd.Execute(); err != nil {
				t.Fatalf("execute: %v", err)
			}
		})
		if got := readVERSION(t); got != "v1.0.0" {
			t.Fatalf("expected VERSION=v1.0.0, got %q", got)
		}
		if !strings.Contains(out, "New Version: v1.0.0") {
			t.Fatalf("unexpected stdout:\n%s", out)
		}
	})
}

func TestSetPre_SetAndClear(t *testing.T) {
	withTempWD(t, func(tmp string) {
		writeVERSION(t, "1.2.3")
		// set pre
		out := captureStdout(t, func() {
			cmd.RootCmd.SetArgs([]string{"set", "pre", "rc.1", "--clear=false", "--dry=false"})
			if err := cmd.RootCmd.Execute(); err != nil {
				t.Fatalf("execute: %v", err)
			}
		})
		if got := readVERSION(t); got != "1.2.3-rc.1" {
			t.Fatalf("expected VERSION=1.2.3-rc.1, got %q", got)
		}
		if !strings.Contains(out, "Setting Prerelease") ||
			!strings.Contains(out, "New Version: 1.2.3-rc.1") {
			t.Fatalf("unexpected stdout:\n%s", out)
		}

		// clear pre
		out = captureStdout(t, func() {
			cmd.RootCmd.SetArgs([]string{"set", "pre", "--clear", "--value=", "--dry=false"})
			if err := cmd.RootCmd.Execute(); err != nil {
				t.Fatalf("execute: %v", err)
			}
		})
		if got := readVERSION(t); got != "1.2.3" {
			t.Fatalf("expected VERSION=1.2.3 after clear, got %q", got)
		}
		if !strings.Contains(out, "Setting Prerelease") ||
			!strings.Contains(out, "New Version: 1.2.3") {
			t.Fatalf("unexpected stdout:\n%s", out)
		}
	})
}

func TestSetPre_PreservesExistingVPrefix(t *testing.T) {
	withTempWD(t, func(tmp string) {
		writeVERSION(t, "v1.0.0-alpha.1")
		out := captureStdout(t, func() {
			cmd.RootCmd.SetArgs([]string{"set", "pre", "alpha.2", "--clear=false", "--dry=false"})
			if err := cmd.RootCmd.Execute(); err != nil {
				t.Fatalf("execute: %v", err)
			}
		})
		if got := readVERSION(t); got != "v1.0.0-alpha.2" {
			t.Fatalf("expected VERSION=v1.0.0-alpha.2, got %q", got)
		}
		if !strings.Contains(out, "New Version: v1.0.0-alpha.2") {
			t.Fatalf("unexpected stdout:\n%s", out)
		}
	})
}

func TestSetBuild_SetValue(t *testing.T) {
	withTempWD(t, func(tmp string) {
		writeVERSION(t, "1.2.3")
		out := captureStdout(t, func() {
			cmd.RootCmd.SetArgs([]string{"set", "build", "exp.7", "--git=false", "--clear=false", "--dry=false"})
			if err := cmd.RootCmd.Execute(); err != nil {
				t.Fatalf("execute: %v", err)
			}
		})
		if got := readVERSION(t); got != "1.2.3+exp.7" {
			t.Fatalf("expected VERSION=1.2.3+exp.7, got %q", got)
		}
		if !strings.Contains(out, "Setting Build Metadata") ||
			!strings.Contains(out, "New Version: 1.2.3+exp.7") {
			t.Fatalf("unexpected stdout:\n%s", out)
		}
	})
}

func TestSetBuild_Clear(t *testing.T) {
	withTempWD(t, func(tmp string) {
		writeVERSION(t, "1.2.3+exp.7")
		out := captureStdout(t, func() {
			cmd.RootCmd.SetArgs([]string{"set", "build", "--clear", "--value=", "--git=false", "--dry=false"})
			if err := cmd.RootCmd.Execute(); err != nil {
				t.Fatalf("execute: %v", err)
			}
		})
		if got := readVERSION(t); got != "1.2.3" {
			t.Fatalf("expected VERSION=1.2.3 after clear, got %q", got)
		}
		if !strings.Contains(out, "Setting Build Metadata") ||
			!strings.Contains(out, "New Version: 1.2.3") {
			t.Fatalf("unexpected stdout:\n%s", out)
		}
	})
}

func TestSetPre_InvalidValueRejected(t *testing.T) {
	withTempWD(t, func(tmp string) {
		writeVERSION(t, "1.2.3")
		cmd.RootCmd.SetArgs([]string{"set", "pre", "01", "--clear=false", "--dry=false"})
		err := cmd.RootCmd.Execute()
		if err == nil {
			t.Fatal("expected invalid prerelease error")
		}
		if !strings.Contains(err.Error(), "invalid prerelease value") {
			t.Fatalf("unexpected error: %v", err)
		}
		if got := readVERSION(t); got != "1.2.3" {
			t.Fatalf("expected VERSION unchanged, got %q", got)
		}
	})
}

func TestSetBuild_InvalidValueRejected(t *testing.T) {
	withTempWD(t, func(tmp string) {
		writeVERSION(t, "1.2.3")
		cmd.RootCmd.SetArgs([]string{"set", "build", "exp+sha", "--git=false", "--clear=false", "--dry=false"})
		err := cmd.RootCmd.Execute()
		if err == nil {
			t.Fatal("expected invalid build metadata error")
		}
		if !strings.Contains(err.Error(), "invalid build metadata") {
			t.Fatalf("unexpected error: %v", err)
		}
		if got := readVERSION(t); got != "1.2.3" {
			t.Fatalf("expected VERSION unchanged, got %q", got)
		}
	})
}

func TestSetVersion_InvalidRejected(t *testing.T) {
	withTempWD(t, func(tmp string) {
		writeVERSION(t, "1.2.3")
		cmd.RootCmd.SetArgs([]string{"set", "1.2.3-01"})
		err := cmd.RootCmd.Execute()
		if err == nil {
			t.Fatal("expected invalid version error")
		}
		if !strings.Contains(err.Error(), "invalid semantic version") {
			t.Fatalf("unexpected error: %v", err)
		}
		if got := readVERSION(t); got != "1.2.3" {
			t.Fatalf("expected VERSION unchanged, got %q", got)
		}
	})
}

func TestSetPre_InvalidVersionFileRejected(t *testing.T) {
	withTempWD(t, func(tmp string) {
		writeVERSION(t, "1.2")
		cmd.RootCmd.SetArgs([]string{"set", "pre", "rc.1", "--clear=false", "--dry=false"})
		err := cmd.RootCmd.Execute()
		if err == nil {
			t.Fatal("expected invalid VERSION contents error")
		}
		if !strings.Contains(err.Error(), "invalid VERSION contents") {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}

func TestSetPre_ArgumentAndFlagConflictRejected(t *testing.T) {
	withTempWD(t, func(tmp string) {
		writeVERSION(t, "1.2.3")
		cmd.RootCmd.SetArgs([]string{"set", "pre", "rc.1", "--value=rc.2", "--clear=false", "--dry=false"})
		err := cmd.RootCmd.Execute()
		if err == nil {
			t.Fatal("expected argument/flag conflict")
		}
	})
}

func TestSetBuild_ArgumentAndFlagConflictRejected(t *testing.T) {
	withTempWD(t, func(tmp string) {
		writeVERSION(t, "1.2.3")
		cmd.RootCmd.SetArgs([]string{"set", "build", "ci.7", "--value=ci.8", "--git=false", "--clear=false", "--dry=false"})
		err := cmd.RootCmd.Execute()
		if err == nil {
			t.Fatal("expected argument/flag conflict")
		}
	})
}

func TestSet_NoVersionFile_Message(t *testing.T) {
	withTempWD(t, func(tmp string) {
		// Ensure no VERSION
		if _, err := os.Stat(filepath.Join(tmp, "VERSION")); !os.IsNotExist(err) {
			t.Fatalf("expected no VERSION file")
		}
		out := captureStdout(t, func() {
			cmd.RootCmd.SetArgs([]string{"set", "2.0.0"})
			_ = cmd.RootCmd.Execute() // prints message, returns nil
		})
		if !strings.Contains(out, "No VERSION file found") {
			t.Fatalf("expected helpful message, got:\n%s", out)
		}
	})
}
