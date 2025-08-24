package version

import (
	"encoding/json"
	"runtime"
	"runtime/debug"
	"time"

	"github.com/dp1140a/semver/pkg"
)

// These are set via -ldflags in release builds; default to empty/zero.
// Keep them exported if you want to keep using -X with the same names.
var (
	Version   string // release/semver (e.g., 1.2.3)
	Commit    string // git SHA
	Branch    string // git branch
	BuildTime string // RFC3339 if provided via ldflags

	// DevVersion is reported when Version is unset.
	DevVersion = "dev"
)

type VersionInfo struct {
	AppName   string `json:"appName"`
	Version   string `json:"version"`
	Branch    string `json:"branch,omitempty"`
	Commit    string `json:"commit,omitempty"`
	BuildTime string `json:"buildTime,omitempty"`
	Dirty     bool   `json:"dirty"`
	GoVersion string `json:"goVersion"`
	GOOS      string `json:"goos"`
	GOARCH    string `json:"goarch"`
}

// BuildVersion returns Version or DevVersion if unset.
func BuildVersion() string {
	if Version == "" {
		return DevVersion
	}
	return Version
}

// Info returns a fully-populated VersionInfo with ldflags if present,
// otherwise it falls back to runtime/debug build info when available.
func Info() VersionInfo {
	vi := VersionInfo{
		AppName:   pkg.APP_NAME,
		Version:   BuildVersion(),
		Branch:    Branch,
		Commit:    Commit,
		BuildTime: BuildTime,
		GoVersion: runtime.Version(),
		GOOS:      runtime.GOOS,
		GOARCH:    runtime.GOARCH,
	}

	// If ldflags didn’t populate fields, try to fill from build info.
	if bi, ok := debug.ReadBuildInfo(); ok {
		settings := map[string]string{}
		for _, s := range bi.Settings {
			settings[s.Key] = s.Value
		}

		// Version: for module builds, bi.Main.Version may be non-empty.
		// If ldflags didn't set Version and bi reports something useful, use it.
		if (vi.Version == "" || vi.Version == DevVersion) && bi.Main.Version != "" && bi.Main.Version != "(devel)" {
			vi.Version = bi.Main.Version
		}

		// VCS info (Go 1.18+ populates these when built with VCS available)
		if vi.Commit == "" {
			if rev := settings["vcs.revision"]; rev != "" {
				vi.Commit = rev
			}
		}
		if vi.Branch == "" {
			if br := settings["vcs.branch"]; br != "" {
				vi.Branch = br
			}
		}
		if vi.BuildTime == "" {
			if ts := settings["vcs.time"]; ts != "" {
				// Normalize to RFC3339 if parsable
				if t, err := time.Parse(time.RFC3339, ts); err == nil {
					vi.BuildTime = t.UTC().Format(time.RFC3339)
				} else {
					vi.BuildTime = ts
				}
			}
		}
		if dm := settings["vcs.modified"]; dm == "true" {
			vi.Dirty = true
		}
	}

	return vi
}

func NewVersionInfo() *VersionInfo { // backward-compatible constructor
	vi := Info()
	return &vi
}

func String() string {
	b, _ := json.Marshal(Info())
	return string(b)
}

func StringPretty() string {
	b, _ := json.MarshalIndent(Info(), "", "    ")
	return string(b)
}
