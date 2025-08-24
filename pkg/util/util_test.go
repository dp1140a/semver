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
