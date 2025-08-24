package types

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/dp1140a/semver/pkg/util"
)

type Version struct {
	Major      uint16
	Minor      uint16
	Patch      uint16
	PreRelease string
	Build      string
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

func NewVersionFromString(version string) Version {
	s := strings.TrimSpace(version)
	if len(s) > 0 && (s[0] == 'v' || s[0] == 'V') {
		s = s[1:]
	}
	matches := semverRE.FindStringSubmatch(s)
	if matches == nil {
		return Version{} // no fmt.Println side-effect
	}
	semver := Version{}
	semver.Major = parseInt(matches[1])
	semver.Minor = parseInt(matches[2])
	semver.Patch = parseInt(matches[3])
	if len(matches) > 4 {
		semver.PreRelease = matches[4]
	}
	if len(matches) > 5 {
		semver.Build = matches[5]
	}
	return semver
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

func parseInt(s string) uint16 {
	num := 0
	for _, c := range s {
		num = num*10 + int(c-'0')
	}
	return uint16(num)
}

func (v *Version) String() string {
	suffix := ""
	if v.PreRelease != "" {
		suffix += fmt.Sprintf("-%v", v.PreRelease)
	}
	if v.Build != "" {
		suffix += fmt.Sprintf("+%v", v.Build)
	}
	return fmt.Sprintf("%v.%v.%v%v", v.Major, v.Minor, v.Patch, suffix)
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
