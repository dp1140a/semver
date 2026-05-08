package types

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestNewVersionFromString_Basic(t *testing.T) {
	v := NewVersionFromString("1.2.3")
	if v.Major != 1 || v.Minor != 2 || v.Patch != 3 {
		t.Fatalf("parsed wrong: %+v", v)
	}
	if v.PreRelease != "" || v.Build != "" {
		t.Fatalf("unexpected suffixes: pre=%q build=%q", v.PreRelease, v.Build)
	}
}

func TestNewVersionFromString_WithNewline(t *testing.T) {
	// This assumes NewVersionFromString trims whitespace (recommended).
	v := NewVersionFromString("1.2.3\n")
	if v.Major != 1 || v.Minor != 2 || v.Patch != 3 {
		t.Fatalf("expected 1.2.3 with newline trimmed, got: %+v", v)
	}
}

func TestNewVersionFromString_WithLeadingV(t *testing.T) {
	// This assumes NewVersionFromString accepts a leading 'v'/'V' (recommended).
	v := NewVersionFromString("v2.0.1")
	if v.Major != 2 || v.Minor != 0 || v.Patch != 1 {
		t.Fatalf("expected 2.0.1 from v2.0.1, got: %+v", v)
	}
	if got := v.String(); got != "v2.0.1" {
		t.Fatalf("expected prefix preserved in rendering, got %q", got)
	}
}

func TestNewVersionFromString_PreReleaseAndBuild(t *testing.T) {
	v := NewVersionFromString("1.2.3-alpha.1+build.99")
	if v.Major != 1 || v.Minor != 2 || v.Patch != 3 {
		t.Fatalf("parsed wrong core: %+v", v)
	}
	if v.PreRelease != "alpha.1" {
		t.Fatalf("expected prerelease 'alpha.1', got %q", v.PreRelease)
	}
	if v.Build != "build.99" {
		t.Fatalf("expected build 'build.99', got %q", v.Build)
	}
}

func TestNewVersionFromString_Invalid(t *testing.T) {
	v := NewVersionFromString("not-a-version")
	// Zero-value expected on parse failure
	if v.Major != 0 || v.Minor != 0 || v.Patch != 0 || v.PreRelease != "" || v.Build != "" {
		t.Fatalf("expected zero value on invalid parse, got: %+v", v)
	}
}

func TestParseVersion_Invalid(t *testing.T) {
	if _, err := ParseVersion("1.2.3-01"); err == nil {
		t.Fatal("expected parse error for invalid prerelease numeric identifier")
	}
}

func TestParseVersion_PrefixOptional(t *testing.T) {
	tests := []struct {
		input  string
		prefix string
		want   string
	}{
		{"1.0.0-rc.1", "", "1.0.0-rc.1"},
		{"v1.0.0-rc.1", "v", "v1.0.0-rc.1"},
		{"V1.0.0-beta.2", "V", "V1.0.0-beta.2"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			v, err := ParseVersion(tt.input)
			if err != nil {
				t.Fatalf("parse: %v", err)
			}
			if v.Prefix() != tt.prefix {
				t.Fatalf("expected prefix %q, got %q", tt.prefix, v.Prefix())
			}
			if got := v.String(); got != tt.want {
				t.Fatalf("expected %q, got %q", tt.want, got)
			}
		})
	}
}

func TestString_Rendering(t *testing.T) {
	v := Version{Major: 1, Minor: 2, Patch: 3}
	if got := v.String(); got != "1.2.3" {
		t.Fatalf("expected 1.2.3, got %q", got)
	}
	v.PreRelease = "rc.1"
	if got := v.String(); got != "1.2.3-rc.1" {
		t.Fatalf("expected 1.2.3-rc.1, got %q", got)
	}
	v.Build = "exp.sha"
	if got := v.String(); got != "1.2.3-rc.1+exp.sha" {
		t.Fatalf("expected 1.2.3-rc.1+exp.sha, got %q", got)
	}
}

func TestSetters_AffectString(t *testing.T) {
	v := Version{Major: 0, Minor: 1, Patch: 0}
	v.SetPre("beta")
	v.SetBuild("001")
	got := v.String()
	if got != "0.1.0-beta+001" {
		t.Fatalf("expected 0.1.0-beta+001, got %q", got)
	}
}

func TestIncrementMajor_ResetsLowerAndSuffixes(t *testing.T) {
	v := Version{Major: 1, Minor: 2, Patch: 3, PreRelease: "rc.1", Build: "exp"}
	v.IncrementMajor()
	if v.Major != 2 || v.Minor != 0 || v.Patch != 0 {
		t.Fatalf("expected 2.0.0 after major bump, got: %+v", v)
	}
	if v.PreRelease != "" || v.Build != "" {
		t.Fatalf("expected suffixes cleared on major bump, got pre=%q build=%q", v.PreRelease, v.Build)
	}
}

func TestIncrementMinor_ResetsPatchAndSuffixes(t *testing.T) {
	v := Version{Major: 1, Minor: 2, Patch: 3, PreRelease: "rc.1", Build: "exp"}
	v.IncrementMinor()
	if v.Major != 1 || v.Minor != 3 || v.Patch != 0 {
		t.Fatalf("expected 1.3.0 after minor bump, got: %+v", v)
	}
	if v.PreRelease != "" || v.Build != "" {
		t.Fatalf("expected suffixes cleared on minor bump, got pre=%q build=%q", v.PreRelease, v.Build)
	}
}

func TestIncrementPatch_IncrementsOnlyPatchAndClearsSuffixes(t *testing.T) {
	v := Version{Major: 1, Minor: 2, Patch: 3, PreRelease: "rc.1", Build: "exp"}
	v.IncrementPatch()
	if v.Major != 1 || v.Minor != 2 || v.Patch != 4 {
		t.Fatalf("expected 1.2.4 after patch bump, got: %+v", v)
	}
	if v.PreRelease != "" || v.Build != "" {
		t.Fatalf("expected suffixes cleared on patch bump, got pre=%q build=%q", v.PreRelease, v.Build)
	}
}

func TestIncrementPre_StartsAlpha1WhenMissing(t *testing.T) {
	v := Version{Major: 1, Minor: 0, Patch: 1}
	v.IncrementPre()
	if got := v.String(); got != "1.0.1-alpha.1" {
		t.Fatalf("expected 1.0.1-alpha.1, got %q", got)
	}
}

func TestIncrementPre_IncrementsNumericSuffix(t *testing.T) {
	v := Version{Major: 1, Minor: 0, Patch: 1, PreRelease: "beta.2", Build: "ci.7"}
	v.IncrementPre()
	if got := v.String(); got != "1.0.1-beta.3" {
		t.Fatalf("expected 1.0.1-beta.3, got %q", got)
	}
}

func TestIncrementBuild_IncrementsNumericSuffix(t *testing.T) {
	v := Version{Major: 1, Minor: 0, Patch: 1, PreRelease: "rc.1", Build: "ci.7"}
	if err := v.IncrementBuild(); err != nil {
		t.Fatalf("increment build: %v", err)
	}
	if got := v.String(); got != "1.0.1-rc.1+ci.8" {
		t.Fatalf("expected 1.0.1-rc.1+ci.8, got %q", got)
	}
}

func TestIncrementBuild_RequiresExistingMetadata(t *testing.T) {
	v := Version{Major: 1, Minor: 0, Patch: 1}
	if err := v.IncrementBuild(); err == nil {
		t.Fatal("expected missing build metadata error")
	}
}

func TestJson_RoundTripShape(t *testing.T) {
	v := Version{Major: 9, Minor: 8, Patch: 7, PreRelease: "alpha", Build: "42"}
	s := v.Json()
	// Just ensure JSON contains expected keys/values
	if !strings.Contains(s, `"Major": 9`) ||
		!strings.Contains(s, `"Minor": 8`) ||
		!strings.Contains(s, `"Patch": 7`) ||
		!strings.Contains(s, `"PreRelease": "alpha"`) ||
		!strings.Contains(s, `"Build": "42"`) {
		t.Fatalf("unexpected JSON: %s", s)
	}

	// Optional: unmarshal to a struct/map and check fields
	var got Version
	if err := json.Unmarshal([]byte(s), &got); err != nil {
		t.Fatalf("unmarshal json: %v", err)
	}
	if got != v {
		t.Fatalf("round-trip mismatch: want %+v got %+v", v, got)
	}
}
