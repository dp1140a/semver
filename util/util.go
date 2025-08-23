package util

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var NOT_GIT_MSG = `The current directory is not a git project.
either change directories to a git project or first run:
$ git init`

const SemVerRegex = `^(0|[1-9]\d*)\.(0|[1-9]\d*)\.(0|[1-9]\d*)(?:-((?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*)(?:\.(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*))*))?(?:\+([0-9a-zA-Z-\.]+(?:\.[0-9a-zA-Z-]+)*))?$`

var semverRe = regexp.MustCompile(SemVerRegex)

// CleanVersion trims trailing/leading whitespace (incl. \r\n)
func CleanVersion(s string) string {
	return strings.TrimSpace(s)
}

func ValidVersionString(version string) bool {
	version = CleanVersion(version) // <-- ignore trailing newline/CR/spaces
	return semverRe.MatchString(version)
}

func WriteVersionFile(version string) error {
	fmt.Println("Writing VERSION file")
	err := os.WriteFile("VERSION", []byte(version), 0644)
	if err != nil {
		return err
	}
	return nil
}

func VersionFileExists(cwd string) bool {
	_, err := os.Stat("VERSION")
	if err != nil {
		return false
	} else {
		return true
	}
}

func InGitDir(cwd string) {
	//Look for .git directory.  If not exit
	if _, err := os.Stat(filepath.Join(cwd, ".git")); err != nil {
		if os.IsNotExist(err) {
			fmt.Println(NOT_GIT_MSG)
			os.Exit(0)
		}
	}
}
