package issues

import (
	"os"
	"strings"
	"testing"
	"time"
)

// I-0011, second occurrence: two instances in one second gave two identical
// "### <timestamp>" headings, which failed the duplicate-heading rule on main.
func TestInstancesInTheSameSecondShareOneHeading(t *testing.T) {
	r := repo(t)
	is, err := New(r, NewOptions{Title: "Flaky thing", Class: "defect", Note: "first", Now: t0})
	if err != nil {
		t.Fatal(err)
	}
	// new and bump in the same second, then two bumps at one later instant
	if err := Bump(is, "", "second, the same second as the first", t0); err != nil {
		t.Fatal(err)
	}
	later := t0.Add(time.Hour)
	for _, note := range []string{"third", "fourth, the same second as the third"} {
		if err := Bump(is, "", note, later); err != nil {
			t.Fatal(err)
		}
	}
	if err := Bump(is, "", "fifth, a second later", later.Add(time.Second)); err != nil {
		t.Fatal(err)
	}
	raw, _ := os.ReadFile(is.Path)
	text := string(raw)
	seen := map[string]int{}
	for _, l := range strings.Split(text, "\n") {
		if strings.HasPrefix(l, "#") {
			seen[l]++
		}
	}
	for h, n := range seen {
		if n > 1 {
			t.Errorf("heading %q appears %d times", h, n)
		}
	}
	want := "## Instances\n\n### 2026-09-16T10:00:00Z\nfirst\n\nsecond, the same second as the first\n\n### 2026-09-16T11:00:00Z\nthird\n\nfourth, the same second as the third\n\n### 2026-09-16T11:00:01Z\nfifth, a second later\n\n## Remediation\n"
	if !strings.Contains(text, want) {
		t.Errorf("one heading per second, every note kept, in order:\n%s", text)
	}
	if is.Count != 5 {
		t.Errorf("every occurrence is counted: %d", is.Count)
	}
	if back, err := Read(is.Path); err != nil || back.Count != 5 {
		t.Errorf("the file still reads: %v", err)
	}
}
