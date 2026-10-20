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

func TestBare(t *testing.T) {
	for v, want := range map[string]bool{
		"1.2.3": true, "0.39.0": true, "10.0.12": true,
		"v1.2.3": false, "1.2": false, "1.2.3.4": false, "1.2.3-rc.1": false,
		"01.2.3": false, "1.2.3\n": false, "latest": false, "": false,
	} {
		if got := Bare(v); got != want {
			t.Errorf("Bare(%q) = %v, want %v", v, got, want)
		}
	}
}

// The embedded task-done rules number parses and reaches Get (ADR-0180).
func TestTaskDoneRules(t *testing.T) {
	if TaskDoneRules < 1 {
		t.Errorf("TaskDoneRules = %d, want a positive integer from %s", TaskDoneRules, TaskDoneRulesPath)
	}
	if got := Get().TaskDoneRules; got != TaskDoneRules {
		t.Errorf("Get().TaskDoneRules = %d, want %d", got, TaskDoneRules)
	}
}

func TestParseTaskDoneRules(t *testing.T) {
	for _, c := range []struct {
		text string
		want int
	}{{"1\n", 1}, {" 12 ", 12}, {"0", 0}, {"-1", 0}, {"two", 0}, {"", 0}, {"1.5", 0}} {
		n, ok := ParseTaskDoneRules(c.text)
		if n != c.want || ok != (c.want > 0) {
			t.Errorf("ParseTaskDoneRules(%q) = %d, %v, want %d, %v", c.text, n, ok, c.want, c.want > 0)
		}
	}
}
