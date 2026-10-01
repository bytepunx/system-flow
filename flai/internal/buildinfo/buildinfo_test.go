package buildinfo

import "testing"

func TestBelow(t *testing.T) {
	for _, c := range []struct {
		version, other string
		below          bool
	}{
		{"1.26.4", "1.26.5", true},
		{"v1.26.4", "flai/v1.27.0", false}, // not a release: the caller strips the prefix
		{"1.26.4", "v1.27.0", true},
		{"1.9.0", "1.10.0", true},
		{"1.26.5", "1.26.5", false},
		{"1.27.0", "1.26.9", false},
		{"2.0.0-rc.1", "1.30.0", false},
		{"dev", "1.26.5", false},
		{"1.26.5", "", false},
	} {
		if got := Below(c.version, c.other); got != c.below {
			t.Errorf("Below(%q, %q) = %v, want %v", c.version, c.other, got, c.below)
		}
	}
}
