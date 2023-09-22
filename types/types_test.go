package types

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

var tests = []struct {
	versionStr string
	expected   Version
}{
	{"1.2.3+4-56", Version{Major: 1, Minor: 2, Patch: 3, PreRelease: "", Build: "4-56"}},
	{"1.0.0-alpha", Version{Major: 1, Minor: 0, Patch: 0, PreRelease: "alpha", Build: ""}},
	{"1.0.0-alpha.1", Version{Major: 1, Minor: 0, Patch: 0, PreRelease: "alpha.1", Build: ""}},
	{"1.0.0-0.3.7", Version{Major: 1, Minor: 0, Patch: 0, PreRelease: "0.3.7", Build: ""}},
	{"1.0.0-x.7.z.92", Version{Major: 1, Minor: 0, Patch: 0, PreRelease: "x.7.z.92", Build: ""}},
	{"1.0.0-x-y-z.--", Version{Major: 1, Minor: 0, Patch: 0, PreRelease: "x-y-z.--", Build: ""}},
	{"1.0.0-alpha+001", Version{Major: 1, Minor: 0, Patch: 0, PreRelease: "alpha", Build: "001"}},
	{"4.5.6-alpha+b1135a", Version{Major: 4, Minor: 5, Patch: 6, PreRelease: "alpha", Build: "b1135a"}},
	{"1.0.0+20130313144700", Version{Major: 1, Minor: 0, Patch: 0, PreRelease: "", Build: "20130313144700"}},
	{"1.0.0-beta+exp.sha.5114f85", Version{Major: 1, Minor: 0, Patch: 0, PreRelease: "beta", Build: "exp.sha.5114f85"}},
	{"1.0.0+21AF26D3----117B344092BD", Version{Major: 1, Minor: 0, Patch: 0, PreRelease: "", Build: "21AF26D3----117B344092BD"}},
	{"1.2.3+exp.123-123", Version{Major: 1, Minor: 2, Patch: 3, PreRelease: "", Build: "exp.123-123"}},
}

func TestNewVersion(t *testing.T) {
	tests := []struct {
		name     string
		expected Version
	}{

		{"New Version", Version{0, 0, 0, "", ""}},
	}
	for _, tt := range tests {
		assert.Equal(t, NewVersion(), tt.expected, "Should Equal")
	}
}

func TestNewVersionFromString(t *testing.T) {
	for _, tt := range tests {
		t.Run(tt.versionStr, func(t *testing.T) {
			got := NewVersionFromString(tt.versionStr)
			fmt.Printf("%s | %v\n", tt.versionStr, got)
			assert.Equal(t, tt.expected, got, "Should Equal")
		})
	}
}

func TestVersion_string(t *testing.T) {
	for _, tt := range tests {
		t.Run(tt.versionStr, func(t *testing.T) {
			v := NewVersionFromString(tt.versionStr)
			got := v.String()
			//fmt.Printf("%s | %s | %s | %s\n", tt.versionStr, got, v.PreRelease, v.Build)
			assert.Equalf(t, tt.versionStr, got, "Not Equal")
		})
	}
}

func generateTests() {
	for _, tt := range tests {
		v := NewVersionFromString(tt.versionStr)
		vStr := v.PrettyPrint()
		t := fmt.Sprintf("{\"%s\", Version%s},", tt.versionStr, vStr)
		fmt.Println(t)
	}
}
