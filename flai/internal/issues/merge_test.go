package issues

import (
	"errors"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// baseIssue is an open issue with one instance, at t0, costing 10m, as New
// writes it.
func baseIssue(t *testing.T) string {
	t.Helper()
	is, err := New(repo(t), NewOptions{Title: "Flaky sync", Class: "defect", Cost: "10m", Story: "S-0001", Note: "First occurrence.", Now: t0})
	if err != nil {
		t.Fatal(err)
	}
	return is.Marshal()
}

// version is text as an issue that saves to a file of its own, so that Bump
// and Close can change it.
func version(t *testing.T, text string) *Issue {
	t.Helper()
	is, err := Parse(text)
	if err != nil {
		t.Fatal(err)
	}
	is.Path = filepath.Join(t.TempDir(), "issue.md")
	return is
}

// bumped is text bumped once.
func bumped(t *testing.T, text, story, cost, note string, at time.Time) string {
	t.Helper()
	is := version(t, text)
	if err := Bump(is, story, cost, note, at); err != nil {
		t.Fatal(err)
	}
	return is.Marshal()
}

// merged merges the three versions, failing the test when they do not merge
// or the result does not read back as a valid issue unchanged.
func merged(t *testing.T, base, ours, theirs string) string {
	t.Helper()
	out, err := Merge(base, ours, theirs)
	if err != nil {
		t.Fatal(err)
	}
	is, err := Parse(out)
	if err != nil {
		t.Fatalf("the merge does not parse: %v\n%s", err, out)
	}
	if err := is.Validate(); err != nil {
		t.Errorf("the merge is not a valid issue: %v\n%s", err, out)
	}
	if is.Marshal() != out {
		t.Errorf("the merge does not round-trip:\n%s", out)
	}
	return out
}

// I-0112, second instance: two branches bumped one issue, so the rebase
// stopped on its file. The merge is the file both bumps give one after the
// other, in timestamp order.
func TestMergeKeepsTwoBumpsInTimestampOrder(t *testing.T) {
	base := baseIssue(t)
	ours := bumped(t, base, "S-0002", "20m", "Ours, later.", t0.Add(2*time.Hour))
	theirs := bumped(t, base, "S-0003", "30m", "Theirs, earlier.", t0.Add(time.Hour))
	got := merged(t, base, ours, theirs)
	want := bumped(t, bumped(t, base, "S-0003", "30m", "Theirs, earlier.", t0.Add(time.Hour)), "S-0002", "20m", "Ours, later.", t0.Add(2*time.Hour))
	if got != want {
		t.Errorf("merge:\n%s\nwant:\n%s", got, want)
	}
	is, _ := Parse(got)
	if is.Count != 3 || is.Cost != "20m" || is.LastReported != "2026-09-16T12:00:00Z" || is.Updated != "2026-09-16T12:00:00Z" {
		t.Errorf("count 3, cost (10m+30m+20m)/3, the later last_reported and updated: %+v", is)
	}
	if again := merged(t, base, theirs, ours); again != got {
		t.Errorf("the merge depends on which side is ours:\n%s", again)
	}
}

func TestMergeABumpAndACloseIsAClosedIssueHoldingTheBump(t *testing.T) {
	base := baseIssue(t)
	ours := bumped(t, base, "S-0002", "20m", "Occurred again.", t0.Add(time.Hour))
	closing := version(t, base)
	if err := Close(closing, "fixed by S-0009", t0.Add(2*time.Hour)); err != nil {
		t.Fatal(err)
	}
	got := merged(t, base, ours, closing.Marshal())
	want := version(t, ours)
	if err := Close(want, "fixed by S-0009", t0.Add(2*time.Hour)); err != nil {
		t.Fatal(err)
	}
	if got != want.Marshal() {
		t.Errorf("merge:\n%s\nwant:\n%s", got, want.Marshal())
	}
	is, _ := Parse(got)
	if is.Status != "closed" || is.Count != 2 || is.LastReported != "2026-09-16T11:00:00Z" || is.Updated != "2026-09-16T12:00:00Z" {
		t.Errorf("closed, counting the bump: %+v", is)
	}
}

func TestMergeTakesAFieldOneSideChanged(t *testing.T) {
	base := baseIssue(t)
	ours := strings.Replace(base, "class: defect", "class: efficiency", 1)
	theirs := bumped(t, base, "S-0002", "", "Occurred again.", t0.Add(time.Hour))
	is, _ := Parse(merged(t, base, ours, theirs))
	if is.Class != "efficiency" || is.Count != 2 || is.Cost != "10m" {
		t.Errorf("ours's class and theirs's bump: %+v", is)
	}
}

func TestMergeRefusesWhatBothSidesChangedDifferently(t *testing.T) {
	base := baseIssue(t)
	describe := func(text string) string {
		return strings.Replace(base, "## Description\nFlaky sync\n", "## Description\n"+text+"\n", 1)
	}
	for name, c := range map[string]struct {
		ours, theirs, names string
	}{
		"description": {describe("Ours."), describe("Theirs."), "the body before its instances"},
		"class":       {strings.Replace(base, "class: defect", "class: efficiency", 1), strings.Replace(base, "class: defect", "class: blocker", 1), "class"},
		"instance":    {strings.Replace(base, "First occurrence.", "Edited.", 1), bumped(t, base, "", "", "Again.", t0.Add(time.Hour)), "an instance the base recorded"},
	} {
		t.Run(name, func(t *testing.T) {
			out, err := Merge(base, c.ours, c.theirs)
			if !errors.Is(err, ErrNotMerged) || !strings.Contains(err.Error(), c.names) {
				t.Errorf("want ErrNotMerged naming %s, got %v:\n%s", c.names, err, out)
			}
		})
	}
}

func TestMergeRefusesWhatIsNotAnIssue(t *testing.T) {
	base := baseIssue(t)
	for name, c := range map[string][3]string{
		"no base":       {"", base, base},
		"not an issue":  {base, "# no front matter\n", base},
		"invalid issue": {base, base, strings.Replace(base, "status: open", "status: ajar", 1)},
	} {
		t.Run(name, func(t *testing.T) {
			if _, err := Merge(c[0], c[1], c[2]); !errors.Is(err, ErrNotMerged) {
				t.Errorf("want ErrNotMerged, got %v", err)
			}
		})
	}
}

func TestMergeKeepsAnInstanceBothSidesAddedOnce(t *testing.T) {
	base := baseIssue(t)
	side := bumped(t, base, "S-0002", "20m", "Occurred again.", t0.Add(time.Hour))
	got := merged(t, base, side, side)
	if got != side {
		t.Errorf("the identical bump once:\n%s", got)
	}
	ours := bumped(t, side, "S-0003", "", "Ours.", t0.Add(2*time.Hour))
	is, _ := Parse(merged(t, base, ours, side))
	if is.Count != 3 || strings.Count(is.Body, "Occurred again.") != 1 {
		t.Errorf("the shared instance once, counted once: %+v\n%s", is, is.Body)
	}
}

// Instances the two sides added in one second share a heading, as
// insertInstance writes them, since sibling headings must differ.
func TestMergeJoinsInstancesOfOneSecondUnderOneHeading(t *testing.T) {
	base := baseIssue(t)
	at := t0.Add(time.Hour)
	got := merged(t, base, bumped(t, base, "S-0002", "", "Ours.", at), bumped(t, base, "S-0002", "", "Theirs.", at))
	if want := "### 2026-09-16T11:00:00Z\nStory: S-0002.\nOurs.\n\nTheirs.\n\n## Remediation"; !strings.Contains(got, want) {
		t.Errorf("want %q in:\n%s", want, got)
	}
}

func TestMergedTallyWithCostsNotKnown(t *testing.T) {
	unknown := func(n int) tally { return tally{count: n} }
	cost := func(n int, avg time.Duration) tally { return tally{count: n}.at(avg) }
	for name, c := range map[string]struct {
		base, ours, theirs tally
		count              int
		cost               string
	}{
		"none known":                      {unknown(1), unknown(2), unknown(2), 3, ""},
		"all known":                       {cost(1, 10*time.Minute), cost(2, 20*time.Minute), cost(2, 25*time.Minute), 3, "27m"},
		"one side knows":                  {unknown(1), cost(2, 20*time.Minute), unknown(2), 3, "20m"},
		"a side not known takes the base": {cost(1, 10*time.Minute), unknown(2), cost(2, 40*time.Minute), 3, "30m"},
		"the base takes the sides' mean":  {unknown(1), cost(2, 20*time.Minute), cost(2, 40*time.Minute), 3, "30m"},
	} {
		t.Run(name, func(t *testing.T) {
			got := mergedTally(c.base, c.ours, c.theirs, 0)
			if got.count != c.count || got.cost() != c.cost {
				t.Errorf("count %d cost %q, want %d %q", got.count, got.cost(), c.count, c.cost)
			}
		})
	}
}

func TestAddOccurrencesFoldsOneIssueIntoAnother(t *testing.T) {
	r := repo(t)
	into, err := New(r, NewOptions{Title: "Flaky sync", Class: "defect", Cost: "10m", Story: "S-0001", Note: "Main's.", Now: t0.Add(time.Hour)})
	if err != nil {
		t.Fatal(err)
	}
	from, err := New(r, NewOptions{Title: "Flaky sync", Class: "defect", Cost: "20m", Story: "S-0002", Note: "Branch's first.", Now: t0})
	if err != nil {
		t.Fatal(err)
	}
	if err := Bump(from, "S-0002", "20m", "Branch's second.", t0.Add(2*time.Hour)); err != nil {
		t.Fatal(err)
	}
	if err := addOccurrences(into, from); err != nil {
		t.Fatal(err)
	}
	if into.Count != 3 || into.Cost != "17m" || into.FirstReported != "2026-09-16T10:00:00Z" || into.LastReported != "2026-09-16T12:00:00Z" {
		t.Errorf("count 3, cost (10m+40m)/3, the earlier first and later last reported: %+v", into)
	}
	want := "## Instances\n\n### 2026-09-16T10:00:00Z\nStory: S-0002.\nBranch's first.\n\n### 2026-09-16T11:00:00Z\nStory: S-0001.\nMain's.\n\n### 2026-09-16T12:00:00Z\nStory: S-0002.\nBranch's second.\n\n## Remediation\n"
	if !strings.Contains(into.Body, want) || into.ID != "I-0001" {
		t.Errorf("every instance in timestamp order:\n%s", into.Body)
	}
	if err := into.Validate(); err != nil {
		t.Error(err)
	}
}
