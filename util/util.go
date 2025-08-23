package util

import (
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
