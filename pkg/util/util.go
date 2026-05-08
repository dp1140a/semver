package util

import (
	"errors"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var NOT_GIT_MSG = `The current directory is not a git project.
either change directories to a git project or first run:
$ git init`

const SemVerRegex = `^(0|[1-9]\d*)\.(0|[1-9]\d*)\.(0|[1-9]\d*)(?:-((?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*)(?:\.(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*))*))?(?:\+([0-9a-zA-Z-\.]+(?:\.[0-9a-zA-Z-]+)*))?$`
const prereleaseRegex = `^(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*)(?:\.(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*))*$`
const buildRegex = `^[0-9A-Za-z-]+(?:\.[0-9A-Za-z-]+)*$`

var semverRe = regexp.MustCompile(SemVerRegex)
var prereleaseRe = regexp.MustCompile(prereleaseRegex)
var buildRe = regexp.MustCompile(buildRegex)

// CleanVersion trims trailing/leading whitespace (incl. \r\n)
func CleanVersion(s string) string {
	return strings.TrimSpace(s)
}

func ValidVersionString(version string) bool {
	version = CleanVersion(version) // <-- ignore trailing newline/CR/spaces
	if len(version) > 0 && (version[0] == 'v' || version[0] == 'V') {
		version = version[1:]
	}
	return semverRe.MatchString(version)
}

func ValidPrereleaseString(prerelease string) bool {
	prerelease = CleanVersion(prerelease)
	if prerelease == "" {
		return false
	}
	return prereleaseRe.MatchString(prerelease)
}

func ValidBuildString(build string) bool {
	build = CleanVersion(build)
	if build == "" {
		return false
	}
	return buildRe.MatchString(build)
}

func WriteVersionFile(version string) error {
	// remove the println; add newline; or just migrate callers to cli.WriteVersion
	return os.WriteFile("VERSION", []byte(version+"\n"), 0644)
}

func VersionFileExists(cwd string) bool {
	_, err := os.Stat("VERSION")
	if err != nil {
		return false
	} else {
		return true
	}
}

func InGitDir(cwd string) error {
	if _, err := os.Stat(filepath.Join(cwd, ".git")); err != nil {
		if os.IsNotExist(err) {
			return errors.New(NOT_GIT_MSG)
		}
		return err
	}
	return nil
}
