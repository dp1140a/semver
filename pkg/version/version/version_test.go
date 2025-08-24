// pkg/version/version_test.go
package version

import (
	"encoding/json"
	"testing"
)

func TestBuildVersion_DefaultsToDev(t *testing.T) {
	Version = "" // simulate no ldflags
	if got := BuildVersion(); got != DevVersion {
		t.Fatalf("BuildVersion=%q, want %q", got, DevVersion)
	}
}

func TestString_JSONShape(t *testing.T) {
	Version = "" // ensure no panic
	var v map[string]any
	if err := json.Unmarshal([]byte(String()), &v); err != nil {
		t.Fatalf("String() not valid JSON: %v", err)
	}
	if _, ok := v["version"]; !ok {
		t.Fatalf("missing 'version' in JSON: %v", v)
	}
}
