package forge

import (
	"slices"
	"strings"
	"testing"
)

func TestParseMergeMethod(t *testing.T) {
	for _, s := range []string{"merge", "squash", "rebase"} {
		m, err := ParseMergeMethod(s)
		if err != nil {
			t.Errorf("ParseMergeMethod(%q): %v", s, err)
		}
		if string(m) != s || !m.Valid() {
			t.Errorf("ParseMergeMethod(%q) = %q, valid %v; want %q, valid", s, m, m.Valid(), s)
		}
	}
	for _, s := range []string{"", "Merge", "fast-forward", " squash"} {
		m, err := ParseMergeMethod(s)
		if err == nil {
			t.Errorf("ParseMergeMethod(%q) = %q, want an error", s, m)
			continue
		}
		if !strings.Contains(err.Error(), "merge, squash, rebase") {
			t.Errorf("ParseMergeMethod(%q) error %q does not name the methods", s, err)
		}
		if MergeMethod(s).Valid() {
			t.Errorf("MergeMethod(%q).Valid() = true, want false", s)
		}
	}
}

func TestMergeMethodsAreTheConstantsInOrder(t *testing.T) {
	want := []MergeMethod{MergeMerge, MergeSquash, MergeRebase}
	if got := MergeMethods(); !slices.Equal(got, want) {
		t.Errorf("MergeMethods() = %v, want %v", got, want)
	}
}
