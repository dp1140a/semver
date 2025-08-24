package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ReadVersion reads VERSION and returns the trimmed string.
// If the file doesn’t exist, returns ("", nil) so callers can print a helpful message.
func ReadVersion() (string, error) {
	b, err := os.ReadFile("VERSION")
	if err != nil {
		if os.IsNotExist(err) {
			return "", nil
		}
		return "", err
	}
	return strings.TrimSpace(string(b)), nil
}

// WriteVersion atomically writes VERSION with a trailing newline.
func WriteVersion(v string) error {
	tmp := filepath.Join(".", ".VERSION.tmp")
	if err := os.WriteFile(tmp, []byte(v+"\n"), 0o644); err != nil {
		return err
	}
	return os.Rename(tmp, "VERSION")
}

func PrintNoVersionMsg(cwd string) {
	fmt.Printf("No VERSION file found in %v.\nPlease either change directory or first run 'semver init'\n", cwd)
}

// RenderDry prints a standardized dry-run message.
func RenderDry(next string) {
	fmt.Printf("[dry-run] New Version would be: %s (VERSION file unchanged)\n", next)
}
