package types

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/dp1140a/semver/pkg/util"
)

type Version struct {
	Major      uint16
	Minor      uint16
	Patch      uint16
	PreRelease string
	Build      string
	prefix     string
}

func NewVersion() Version {
	return Version{
		Major:      0,
		Minor:      0,
		Patch:      0,
		PreRelease: "",
		Build:      "",
	}
}

var semverRE = regexp.MustCompile(util.SemVerRegex)

var ErrInvalidVersion = errors.New("invalid semantic version")
var ErrNoBuildMetadata = errors.New("no build metadata to bump")

func NewVersionFromString(version string) Version {
	v, err := ParseVersion(version)
	if err != nil {
		return Version{}
	}
	return v
}

func ParseVersion(version string) (Version, error) {
	s := strings.TrimSpace(version)
	prefix := ""
	if len(s) > 0 && (s[0] == 'v' || s[0] == 'V') {
		prefix = s[:1]
		s = s[1:]
	}
	matches := semverRE.FindStringSubmatch(s)
	if matches == nil {
		return Version{}, ErrInvalidVersion
	}
	semver := Version{prefix: prefix}
	semver.Major = parseInt(matches[1])
	semver.Minor = parseInt(matches[2])
	semver.Patch = parseInt(matches[3])
	if len(matches) > 4 {
		semver.PreRelease = matches[4]
	}
	if len(matches) > 5 {
		semver.Build = matches[5]
	}
	return semver, nil
}

func (v *Version) IncrementMajor() {
	v.Major++
	v.Minor = 0
	v.Patch = 0
	v.PreRelease = ""
	v.Build = ""
}

func (v *Version) IncrementMinor() {
	v.Minor++
	v.Patch = 0
	v.PreRelease = ""
	v.Build = ""
}

func (v *Version) IncrementPatch() {
	v.Patch++
	v.PreRelease = ""
	v.Build = ""
}

func (v *Version) SetBuild(build string) {
	v.Build = build
}

func (v *Version) SetPre(pre string) {
	v.PreRelease = pre
}

func (v *Version) IncrementPre() {
	if v.PreRelease == "" {
		v.PreRelease = "alpha.1"
		v.Build = ""
		return
	}

	parts := strings.Split(v.PreRelease, ".")
	last := parts[len(parts)-1]
	if n, err := strconv.Atoi(last); err == nil {
		parts[len(parts)-1] = strconv.Itoa(n + 1)
		v.PreRelease = strings.Join(parts, ".")
	} else {
		v.PreRelease += ".1"
	}
	v.Build = ""
}

func (v *Version) IncrementBuild() error {
	if v.Build == "" {
		return ErrNoBuildMetadata
	}

	parts := strings.Split(v.Build, ".")
	last := parts[len(parts)-1]
	if n, err := strconv.Atoi(last); err == nil {
		parts[len(parts)-1] = strconv.Itoa(n + 1)
		v.Build = strings.Join(parts, ".")
		return nil
	}

	v.Build += ".1"
	return nil
}

func (v Version) IsValid() bool {
	return util.ValidVersionString(v.CanonicalString())
}

func (v Version) Prefix() string {
	return v.prefix
}

func (v *Version) SetPrefix(prefix string) {
	if prefix == "v" || prefix == "V" {
		v.prefix = prefix
		return
	}
	v.prefix = ""
}

func parseInt(s string) uint16 {
	num := 0
	for _, c := range s {
		num = num*10 + int(c-'0')
	}
	return uint16(num)
}

func (v Version) CanonicalString() string {
	suffix := ""
	if v.PreRelease != "" {
		suffix += fmt.Sprintf("-%v", v.PreRelease)
	}
	if v.Build != "" {
		suffix += fmt.Sprintf("+%v", v.Build)
	}
	return fmt.Sprintf("%v.%v.%v%v", v.Major, v.Minor, v.Patch, suffix)
}

func (v *Version) String() string {
	return v.prefix + v.CanonicalString()
}

func (v *Version) Json() string {
	json, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return ""
	}
	return string(json)
}

func (v *Version) PrettyPrint() string {
	return fmt.Sprintf("{Major: %v, Minor: %v, Patch: %v, PreRelease: \"%s\", Build: \"%s\"}", v.Major, v.Minor, v.Patch, v.PreRelease, v.Build)
}
