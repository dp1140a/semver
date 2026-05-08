package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidVersionString(t *testing.T) {
	tests := []struct {
		versionStr string
	}{
		{"1.2.3+4-56"},
		{"1.0.0-alpha"},
		{"1.0.0-alpha.1"},
		{"1.0.0-0.3.7"},
		{"1.0.0-x.7.z.92"},
		{"1.0.0-x-y-z.--"},
		{"1.0.0-alpha+001"},
		{"1.0.0+20130313144700"},
		{"1.0.0-beta+exp.sha.5114f85"},
		{"1.0.0+21AF26D3----117B344092BD"},
		{"1.2.3+exp.123-123"},
		{`1.2.3
`},
	}
	for _, tt := range tests {
		t.Run(tt.versionStr, func(t *testing.T) {
			assert.True(t, ValidVersionString(tt.versionStr), "Should Be True")
		})
	}
}

func TestValidPrereleaseString(t *testing.T) {
	tests := []struct {
		value string
		want  bool
	}{
		{"alpha", true},
		{"alpha.1", true},
		{"0.3.7", true},
		{"x-y-z.--", true},
		{"", false},
		{"alpha..1", false},
		{"01", false},
		{"rc+1", false},
	}

	for _, tt := range tests {
		t.Run(tt.value, func(t *testing.T) {
			assert.Equal(t, tt.want, ValidPrereleaseString(tt.value))
		})
	}
}

func TestValidBuildString(t *testing.T) {
	tests := []struct {
		value string
		want  bool
	}{
		{"001", true},
		{"exp.sha.5114f85", true},
		{"21AF26D3----117B344092BD", true},
		{"", false},
		{"exp..sha", false},
		{"build+1", false},
	}

	for _, tt := range tests {
		t.Run(tt.value, func(t *testing.T) {
			assert.Equal(t, tt.want, ValidBuildString(tt.value))
		})
	}
}
