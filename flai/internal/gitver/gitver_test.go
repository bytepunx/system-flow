package gitver

import "testing"

func TestParseAndCompare(t *testing.T) {
	for _, c := range []struct {
		out      string
		want     Version
		relative bool
	}{
		{"git version 2.47.3", Version{2, 47}, false},
		{"git version 2.48.0", Version{2, 48}, true},
		{"git version 2.54.0\n", Version{2, 54}, true},
		{"git version 2.39.3 (Apple Git-146)", Version{2, 39}, false},
		{"git version 3.0.1.windows.1", Version{3, 0}, true},
	} {
		got, err := Parse(c.out)
		if err != nil || got != c.want {
			t.Errorf("Parse(%q) = %v, %v; want %v", c.out, got, err, c.want)
		}
		if got.AtLeast(RelativeWorktrees) != c.relative {
			t.Errorf("%v AtLeast %v = %v; want %v", got, RelativeWorktrees, !c.relative, c.relative)
		}
	}
	if _, err := Parse("not git"); err == nil {
		t.Error("output with no version must be an error")
	}
}
