package issues

import (
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/mdlint"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// S-0198: each instance names the story it was recorded for, first under its
// heading, and the issue lists those stories once each, in order.
func TestAnInstanceNamesItsStory(t *testing.T) {
	r := repo(t)
	is, err := New(r, NewOptions{Title: "Flaky thing", Class: "defect", Note: "first", Story: "s-198", Now: t0})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(is.Body, "### 2026-09-16T10:00:00Z\nStory: S-0198.\nfirst\n") {
		t.Errorf("new names its story first:\n%s", is.Body)
	}
	if err := Bump(is, "S-0200", "", "second", t0.Add(time.Hour)); err != nil {
		t.Fatal(err)
	}
	if err := Bump(is, "", "", "third, for no story", t0.Add(2*time.Hour)); err != nil {
		t.Fatal(err)
	}
	if err := Bump(is, "S-0198", "", "fourth", t0.Add(3*time.Hour)); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(is.Body, "### 2026-09-16T11:00:00Z\nStory: S-0200.\nsecond\n\n### 2026-09-16T12:00:00Z\nthird, for no story\n\n") {
		t.Errorf("bump names its story, or none:\n%s", is.Body)
	}
	back, err := Read(is.Path)
	if err != nil {
		t.Fatal(err)
	}
	if got := strings.Join(Stories(back), " "); got != "S-0198 S-0200" {
		t.Errorf("stories in order, once each: %q", got)
	}
	if _, err := New(r, NewOptions{Title: "x", Class: "defect", Story: "T-0001", Now: t0}); err == nil {
		t.Error("a task is not a story")
	}
	if err := Bump(is, "nonsense", "", "", t0); err == nil || is.Count != 4 {
		t.Errorf("a bad story is refused before anything changes: %v, count %d", err, is.Count)
	}
}

// S-0198 with I-0011: an instance joined to one in the same second names its
// story only when that instance does not already.
func TestTheSameSecondJoinNamesAStoryOnce(t *testing.T) {
	r := repo(t)
	is, err := New(r, NewOptions{Title: "Flaky thing", Class: "defect", Note: "first", Story: "S-0198", Now: t0})
	if err != nil {
		t.Fatal(err)
	}
	for _, b := range []struct{ story, note string }{{"S-0198", "second"}, {"S-0200", "third"}, {"", "fourth"}} {
		if err := Bump(is, b.story, "", b.note, t0); err != nil {
			t.Fatal(err)
		}
	}
	want := "### 2026-09-16T10:00:00Z\nStory: S-0198.\nfirst\n\nsecond\n\nStory: S-0200.\nthird\n\nfourth\n\n## Remediation\n"
	if !strings.Contains(is.Body, want) {
		t.Errorf("one heading, each story named once:\n%s", is.Body)
	}
	if strings.Count(is.Body, "### ") != 1 {
		t.Errorf("one heading for the second:\n%s", is.Body)
	}
}

func TestRecordedByListsTheIssuesAStoryRecorded(t *testing.T) {
	r := repo(t)
	a, _ := New(r, NewOptions{Title: "A", Class: "defect", Story: "S-0198", Now: t0})
	b, _ := New(r, NewOptions{Title: "B", Class: "defect", Now: t0})
	c, _ := New(r, NewOptions{Title: "C", Class: "defect", Story: "S-0200", Now: t0})
	if err := Bump(c, "S-0198", "", "again", t0.Add(time.Hour)); err != nil {
		t.Fatal(err)
	}
	var got []string
	for _, is := range RecordedBy([]*Issue{a, b, c}, "s-198") {
		got = append(got, is.ID)
	}
	if strings.Join(got, " ") != "I-0001 I-0003" {
		t.Errorf("recorded by S-0198: %v", got)
	}
	if len(RecordedBy([]*Issue{a, b, c}, "S-0999")) != 0 {
		t.Error("a story that recorded nothing")
	}
}

// storyLint is the markdown lint a story's or an issue's file must pass.
func storyLint(t *testing.T) *mdlint.Config {
	t.Helper()
	lint, err := mdlint.Parse([]byte("default: true\nMD013: false\nMD022:\n  lines_below: 0\nMD032: false\nMD041: false\n"), false, false)
	if err != nil {
		t.Fatal(err)
	}
	return lint
}

// lintStory fails the test when the drafted story's file is not lint-clean.
func lintStory(t *testing.T, lint *mdlint.Config, d StoryDraft) {
	t.Helper()
	doc := "---\nid: S-0301\n---\n\n# S-0301 " + d.Title + "\n\n" + d.Body
	if f := lint.Lint(doc); len(f) > 0 {
		t.Errorf("lint: %v\n%s", f, doc)
	}
}

var cycle = 168 * time.Hour

func TestForStory(t *testing.T) {
	lint := storyLint(t)
	for _, c := range []struct{ class, nature string }{
		{"defect", "remediation"}, {"blocker", "remediation"}, {"efficiency", "improvement"}, {"impression", "improvement"},
	} {
		for _, fix := range []string{"", "Pin the linter in the Makefile.\n\n- and say so in design/tech"} {
			is := &Issue{ID: "I-0056", Title: "Lint: version mismatch", Class: c.class, Path: "/x/design/issues/I-0056-lint-version-mismatch.md",
				Body: "\n# I-0056 Lint: version mismatch\n\n## Instances\n\n### 2026-09-16T10:00:00Z\nfirst\n\n## Remediation\n" + fix + "\n"}
			d := ForStory(is, t0, cycle)
			if d.Title != is.Title || d.Nature != c.nature || !d.Draft || d.CostOfDelay != nil {
				t.Errorf("%s: a draft with no cost of delay: %+v", c.class, d)
			}
			if !strings.HasPrefix(d.Body, "## Goal\n\nThis story remediates [I-0056](../../../design/issues/I-0056-lint-version-mismatch.md), \"Lint: version mismatch\".") {
				t.Errorf("goal links the issue:\n%s", d.Body)
			}
			if fix != "" && !strings.Contains(d.Body, "solution:\n\n"+fix+"\n\n## Acceptance criteria\n") {
				t.Errorf("goal recommends the remediation:\n%s", d.Body)
			}
			if fix == "" && !strings.Contains(d.Body, "propose one from its instances") {
				t.Errorf("goal asks for a solution:\n%s", d.Body)
			}
			criteria := d.Body[strings.Index(d.Body, "## Acceptance criteria\n"):strings.Index(d.Body, "## Tasks")]
			lines := strings.Split(strings.TrimSpace(criteria), "\n")
			if len(lines) != 3 || !strings.HasPrefix(lines[1], "- [ ] The cause I-0056 describes no longer occurs") || lines[2] != "- [ ] I-0056 is closed with `flai issue close I-0056 --reason` saying what fixed it" {
				t.Errorf("criteria, closing the issue last:\n%s", criteria)
			}
			if !strings.HasSuffix(d.Body, "\n\n## Tasks\n\n## Notes\n") {
				t.Errorf("tasks and empty notes:\n%s", d.Body)
			}
			lintStory(t, lint, d)
		}
	}
}

// TestForStoryLeavesOutStoriesMadeBefore: a story made again, after the
// first was cancelled, recommends the solution, not the line naming the first.
func TestForStoryLeavesOutStoriesMadeBefore(t *testing.T) {
	dir := t.TempDir()
	is := &Issue{ID: "I-0056", Title: "Lint version mismatch", Class: "defect", Status: "open", Count: 1,
		FirstReported: "2026-09-16T10:00:00Z", LastReported: "2026-09-16T10:00:00Z", Updated: "2026-09-16T10:00:00Z",
		Path: filepath.Join(dir, "I-0056-lint-version-mismatch.md"),
		Body: "\n# I-0056 Lint version mismatch\n\n## Remediation\n\nPin the linter.\n"}
	if err := LinkStory(is, "S-0301", t0); err != nil {
		t.Fatal(err)
	}
	d := ForStory(is, t0, cycle)
	if strings.Contains(d.Body, "S-0301") || !strings.Contains(d.Body, "solution:\n\nPin the linter.\n\n## Acceptance criteria\n") {
		t.Errorf("the solution without the earlier story:\n%s", d.Body)
	}
	is.Body = "\n# I-0056 Lint version mismatch\n\n## Remediation\n"
	if err := LinkStory(is, "S-0301", t0); err != nil {
		t.Fatal(err)
	}
	if d := ForStory(is, t0, cycle); !strings.Contains(d.Body, "propose one from its instances") {
		t.Errorf("no solution yet:\n%s", d.Body)
	}
}

// costIssue is I-0007 with a cost, a count, a first report that long before
// t0, and an Impact section when impact is not empty.
func costIssue(cost string, count int, before time.Duration, impact string) *Issue {
	body := "\n# I-0007 Slow builds\n\n## Instances\n\n### 2026-09-16T10:00:00Z\nfirst\n\n"
	if impact != "" {
		body += "## Impact\n\n" + impact + "\n"
	}
	return &Issue{ID: "I-0007", Title: "Slow builds", Class: "efficiency", Count: count, Cost: cost,
		FirstReported: t0.Add(-before).Format(workitem.TimeFormat), Path: "/x/design/issues/I-0007-slow-builds.md",
		Body: body + "## Remediation\n"}
}

// notes is the text of a drafted story's Notes section.
func notes(d StoryDraft) string {
	_, after, _ := strings.Cut(d.Body, "\n## Notes\n")
	return strings.TrimSpace(after)
}

// S-0203: a story made from an issue is a draft, and its time lost per cycle
// is the issue's cost × count ÷ the cycles since it was first reported, at
// least one, by flai at the time the story was made.
func TestForStoryDerivesTheTimeLostPerCycle(t *testing.T) {
	lint := storyLint(t)
	day := 24 * time.Hour
	for _, c := range []struct {
		cost   string
		count  int
		before time.Duration
		want   string
	}{
		{"7m", 5, 26*time.Hour + 24*time.Minute, "35m"},
		{"7m", 5, 0, "35m"},
		{"10m", 6, 21 * day, "20m"},
		{"7m", 5, 17*day + 12*time.Hour, "14m"},
		{"1h10m", 2, 7 * day, "2h20m"},
		{"10s", 1, 21 * day, "3s"},
		{"1s", 1, 70 * day, "1s"},
		{"25s", 3, 0, "1m"},
		{"20s", 2, 0, "40s"},
		{"40s", 3, 0, "2m"},
	} {
		is := costIssue(c.cost, c.count, c.before, "")
		d := ForStory(is, t0, cycle)
		if !d.Draft {
			t.Errorf("%s × %d: not a draft", c.cost, c.count)
		}
		cod := d.CostOfDelay
		if cod == nil || cod.Inputs == nil {
			t.Fatalf("%s × %d: no cost of delay", c.cost, c.count)
		}
		if cod.Inputs.TimeLostPerCycle != c.want || cod.Inputs.RevenuePerWeek != nil || cod.Inputs.PenaltyPerWeek != nil || cod.Value != nil {
			t.Errorf("%s × %d over %v: inputs %+v, want time lost %s", c.cost, c.count, c.before, *cod.Inputs, c.want)
		}
		if lost, err := time.ParseDuration(cod.Inputs.TimeLostPerCycle); err != nil || lost <= 0 {
			t.Errorf("time lost %q is not a duration longer than zero: %v", cod.Inputs.TimeLostPerCycle, err)
		}
		if cod.Inputs.By != "flai" || cod.Inputs.At != "2026-09-16T10:00:00Z" || cod.By != "" || cod.At != "" {
			t.Errorf("the inputs set by flai at t0, and no value's stamp (ADR-0080): %+v %+v", cod, *cod.Inputs)
		}
		lintStory(t, lint, d)
	}
	within := ForStory(costIssue("7m", 5, 26*time.Hour+24*time.Minute, ""), t0, cycle)
	want := "Cost of delay inputs set by flai from I-0007. time_lost_per_cycle 35m: 7m per occurrence × 5 occurrences ÷ 1 cycle of 168h (first reported 2026-09-15T07:36:00Z, 1.1 days before this story; under one cycle counts as one)."
	if got := notes(within); got != want {
		t.Errorf("notes within one cycle:\n got %s\nwant %s", got, want)
	}
	several := ForStory(costIssue("10m", 6, 21*day, ""), t0, cycle)
	want = "Cost of delay inputs set by flai from I-0007. time_lost_per_cycle 20m: 10m per occurrence × 6 occurrences ÷ 3 cycles of 168h (first reported 2026-08-26T10:00:00Z, 21 days before this story)."
	if got := notes(several); got != want {
		t.Errorf("notes over several cycles:\n got %s\nwant %s", got, want)
	}
	half := ForStory(costIssue("7m", 1, 17*day+12*time.Hour, ""), t0, cycle)
	if got := notes(half); !strings.Contains(got, "7m per occurrence × 1 occurrence ÷ 2.5 cycles of 168h") {
		t.Errorf("one occurrence, two and a half cycles: %s", got)
	}
	unknown := costIssue("7m", 5, 0, "")
	unknown.FirstReported = "yesterday"
	d := ForStory(unknown, t0, cycle)
	if d.CostOfDelay == nil || d.CostOfDelay.Inputs.TimeLostPerCycle != "35m" || !strings.Contains(notes(d), `(first reported "yesterday" is not a timestamp, so it counts as one cycle).`) {
		t.Errorf("an unparsable first report counts as one cycle: %+v\n%s", d.CostOfDelay, notes(d))
	}
}

// S-0203: an issue with no cost and no Impact inputs gives no cost of delay,
// and the story's Notes stay empty.
func TestForStoryWithNoCostGivesNoCostOfDelay(t *testing.T) {
	for _, is := range []*Issue{
		costIssue("", 5, 0, ""),
		costIssue("", 5, 0, "Builds wait on the cache every morning.\n"),
	} {
		d := ForStory(is, t0, cycle)
		if !d.Draft || d.CostOfDelay != nil || !strings.HasSuffix(d.Body, "\n\n## Tasks\n\n## Notes\n") {
			t.Errorf("a draft with no cost of delay and empty notes: %+v", d)
		}
	}
}

// S-0203: the amounts and time lost an issue's Impact section gives are
// carried over, and its time lost wins over the one derived from cost and
// count.
func TestForStoryCarriesTheImpactOver(t *testing.T) {
	impact := "Three stories waited on it.\n\n- revenue_per_week: 1200\n* penalty_per_week: 300.5\ntime_lost_per_cycle: 4h0m\n- revenue_per_week: 9999\n"
	d := ForStory(costIssue("7m", 5, 0, impact), t0, cycle)
	cod := d.CostOfDelay
	if cod == nil || cod.Inputs == nil {
		t.Fatalf("no cost of delay:\n%s", d.Body)
	}
	in := cod.Inputs
	if in.RevenuePerWeek == nil || *in.RevenuePerWeek != 1200 || in.PenaltyPerWeek == nil || *in.PenaltyPerWeek != 300.5 || in.TimeLostPerCycle != "4h" {
		t.Errorf("inputs carried over, the first of each: %+v", *in)
	}
	if in.By != "flai" || in.At != "2026-09-16T10:00:00Z" || cod.By != "" || cod.At != "" {
		t.Errorf("the inputs set by flai at t0, and no value's stamp (ADR-0080): %+v %+v", cod, *in)
	}
	want := "Cost of delay inputs set by flai from I-0007. revenue_per_week 1200, penalty_per_week 300.5 and time_lost_per_cycle 4h carried over from I-0007's Impact section. Its Impact time_lost_per_cycle was taken rather than the 35m derived from its cost and count."
	if got := notes(d); got != want {
		t.Errorf("notes:\n got %s\nwant %s", got, want)
	}
	lintStory(t, storyLint(t), d)

	amounts := ForStory(costIssue("7m", 5, 0, "- penalty_per_week: 0\n"), t0, cycle)
	if in := amounts.CostOfDelay.Inputs; in.PenaltyPerWeek == nil || *in.PenaltyPerWeek != 0 || in.TimeLostPerCycle != "35m" {
		t.Errorf("an amount of zero carried over beside the derived time lost: %+v", *in)
	}
	if got := notes(amounts); !strings.Contains(got, "time_lost_per_cycle 35m: 7m per occurrence") || !strings.HasSuffix(got, " penalty_per_week 0 carried over from I-0007's Impact section.") {
		t.Errorf("derivation, then the amount carried over: %s", got)
	}
}

// S-0203: an Impact value that does not parse is left out, and the story's
// Notes say so.
func TestForStoryLeavesABadImpactValueOut(t *testing.T) {
	impact := "- revenue_per_week: lots\n- penalty_per_week: -5\n- time_lost_per_cycle: 0s\n- penalty_per_week: 300\n"
	d := ForStory(costIssue("", 1, 0, impact), t0, cycle)
	if d.CostOfDelay == nil {
		t.Fatalf("the good penalty is set:\n%s", d.Body)
	}
	if in := d.CostOfDelay.Inputs; in.RevenuePerWeek != nil || in.PenaltyPerWeek == nil || *in.PenaltyPerWeek != 300 || in.TimeLostPerCycle != "" {
		t.Errorf("bad values left out, the first good one used: %+v", *in)
	}
	want := `Cost of delay inputs set by flai from I-0007. penalty_per_week 300 carried over from I-0007's Impact section. ` +
		`I-0007's Impact gives revenue_per_week "lots", which is not an amount of zero or more, so it was left out. ` +
		`I-0007's Impact gives penalty_per_week "-5", which is not an amount of zero or more, so it was left out. ` +
		`I-0007's Impact gives time_lost_per_cycle "0s", which is not a duration longer than zero, so it was left out.`
	if got := notes(d); got != want {
		t.Errorf("notes:\n got %s\nwant %s", got, want)
	}
	lintStory(t, storyLint(t), d)

	none := ForStory(costIssue("", 1, 0, "- revenue_per_week: NaN\n"), t0, cycle)
	if none.CostOfDelay != nil || notes(none) != `I-0007's Impact gives revenue_per_week "NaN", which is not an amount of zero or more, so it was left out.` {
		t.Errorf("nothing set, the value left out said so: %+v\n%s", none.CostOfDelay, none.Body)
	}
}

// S-0203: the issue's Remediation section names the story made from it, as
// its last paragraph, and the file reads back with it.
func TestLinkStoryNamesTheStoryUnderRemediation(t *testing.T) {
	lint := storyLint(t)
	head := "\n# I-0007 Slow builds\n\n## Instances\n\n### 2026-09-16T10:00:00Z\nfirst\n"
	line := "Story S-0301 remediates this issue, created from it at 2026-09-16T11:00:00Z.\n"
	for _, c := range []struct{ name, body, want string }{
		{"with a solution", head + "\n## Remediation\nCache the modules.\n\n- in CI too\n\n", head + "\n## Remediation\n\nCache the modules.\n\n- in CI too\n\n" + line},
		{"empty", head + "\n## Remediation\n", head + "\n## Remediation\n\n" + line},
		{"missing", head, head + "\n## Remediation\n\n" + line},
		{"missing, no final newline", strings.TrimSuffix(head, "\n"), head + "\n## Remediation\n\n" + line},
		{"followed by a section", head + "\n## Remediation\nCache the modules.\n\n## Related\nI-0003\n", head + "\n## Remediation\n\nCache the modules.\n\n" + line + "\n## Related\nI-0003\n"},
	} {
		is := &Issue{ID: "I-0007", Title: "Slow builds", Class: "efficiency", Status: "open", Count: 1,
			FirstReported: "2026-09-16T10:00:00Z", LastReported: "2026-09-16T10:00:00Z", Updated: "2026-09-16T10:00:00Z",
			Path: filepath.Join(t.TempDir(), "I-0007-slow-builds.md"), Body: c.body}
		if err := LinkStory(is, "s-301", t0.Add(time.Hour)); err != nil {
			t.Fatalf("%s: %v", c.name, err)
		}
		back, err := Read(is.Path)
		if err != nil {
			t.Fatalf("%s: %v", c.name, err)
		}
		if back.Body != c.want {
			t.Errorf("%s:\n got %q\nwant %q", c.name, back.Body, c.want)
		}
		if back.Updated != "2026-09-16T11:00:00Z" {
			t.Errorf("%s: updated %s", c.name, back.Updated)
		}
		if f := lint.Lint("---\nid: I-0007\n---\n" + back.Body); len(f) > 0 {
			t.Errorf("%s: lint: %v\n%s", c.name, f, back.Body)
		}
	}
	is := &Issue{ID: "I-0007", Path: filepath.Join(t.TempDir(), "I-0007-slow-builds.md"), Body: head}
	for _, story := range []string{"T-0001", " "} {
		if err := LinkStory(is, story, t0); err == nil || is.Body != head {
			t.Errorf("story %q is refused before anything changes: %v", story, err)
		}
	}
}

func TestLinks(t *testing.T) {
	is := &Issue{ID: "I-0056", Path: "/x/design/issues/I-0056-lint-version-mismatch.md"}
	for body, want := range map[string]bool{
		"Remediates I-0056.": true,
		"I-0056":             true,
		"see [it](../../../design/issues/I-0056-lint-version-mismatch.md)": true,
		"Remediates I-00561.": false,
		"Remediates XI-0056.": false,
		"Remediates I-005.":   false,
		"nothing here":        false,
	} {
		if got := Links(body, is); got != want {
			t.Errorf("Links(%q) = %v", body, got)
		}
	}
	legacy := &Issue{ID: "I-012"}
	if !Links("fixes I-0012", legacy) || !Links("fixes I-012", legacy) {
		t.Error("a three-digit issue is named in either form")
	}
}

func TestNoStory(t *testing.T) {
	open := func(id string) *Issue { return &Issue{ID: id, Status: "open"} }
	list := []*Issue{open("I-0001"), open("I-0002"), open("I-0003"), open("I-0004"), {ID: "I-0005", Status: "closed"}}
	items := []*workitem.Item{
		{ID: "S-0010", Type: workitem.Story, Status: workitem.Backlog, Body: "Remediates I-0001."},
		{ID: "S-0011", Type: workitem.Story, Status: workitem.Review, Body: "and I-0002"},
		{ID: "S-0012", Type: workitem.Story, Status: workitem.Done, Body: "Remediated I-0003."},
		{ID: "S-0013", Type: workitem.Story, Status: workitem.Cancelled, Body: "I-0003"},
		{ID: "T-0014", Type: workitem.Task, Status: workitem.InProgress, Body: "I-0004"},
		{ID: "E-0015", Type: workitem.Epic, Status: workitem.InProgress, Body: "I-0004"},
	}
	var got []string
	for _, is := range NoStory(list, items) {
		got = append(got, is.ID)
	}
	if strings.Join(got, " ") != "I-0003 I-0004" {
		t.Errorf("open issues no open story links: %v", got)
	}
}

// The story an issue's file names round-trips through the file.
func TestTheStoryLineSurvivesReading(t *testing.T) {
	r := repo(t)
	if _, err := New(r, NewOptions{Title: "A", Class: "defect", Story: "S-0198", Now: t0}); err != nil {
		t.Fatal(err)
	}
	list, err := List(r)
	if err != nil || len(list) != 1 || filepath.Base(list[0].Path) != "I-0001-a.md" {
		t.Fatalf("list: %v", err)
	}
	if got := Stories(list[0]); len(got) != 1 || got[0] != "S-0198" {
		t.Errorf("stories: %v", got)
	}
}

// S-0224: the story made from an issue the analyzer filed links the issue and
// the report that found it, says in its Notes that the analyzer found it, and
// carries the issue's impact over as cost of delay inputs set by flai. The
// report's link under the issue's Remediation, which is from the issue's
// folder, is not carried into the goal as a solution.
func TestForStoryFromAnAnalyzerIssue(t *testing.T) {
	r := repo(t)
	im := Impact{RevenuePerWeek: "1200", PenaltyPerWeek: "0", TimeLostPerCycle: "90m", Evidence: "12 stories waited 18h on average in review."}
	is, err := New(r, NewOptions{Title: "Review waits a day", Class: "efficiency", Impact: im, Report: report, Now: t0})
	if err != nil {
		t.Fatal(err)
	}
	back, err := Read(is.Path)
	if err != nil {
		t.Fatal(err)
	}
	if got := Reports(back); len(got) != 1 || got[0] != report {
		t.Errorf("reports, the instance's and the link's once: %v", got)
	}
	d := ForStory(back, t0.Add(time.Hour), cycle)
	want := "## Goal\n\nThis story remediates [I-0001](../../../design/issues/I-0001-review-waits-a-day.md), \"Review waits a day\", " +
		"which the analysis in [2026-10-06-bottlenecks.md](../../../design/analysis/2026-10-06-bottlenecks.md) found. " +
		"The issue recommends no solution yet: propose one from its instances before building it.\n\n" +
		"## Acceptance criteria\n- [ ] The cause I-0001 describes no longer occurs, with a test that reproduces it where one fits\n" +
		"- [ ] I-0001 is closed with `flai issue close I-0001 --reason` saying what fixed it\n\n## Tasks\n\n## Notes\n\n" +
		"The analyzer found I-0001: the analysis the goal links gives its evidence.\n\n" +
		"Cost of delay inputs set by flai from I-0001. revenue_per_week 1200, penalty_per_week 0 and time_lost_per_cycle 1h30m carried over from I-0001's Impact section.\n"
	if d.Body != want {
		t.Errorf("body:\n got %q\nwant %q", d.Body, want)
	}
	if d.Nature != "improvement" || !d.Draft || !Links(d.Body, back) {
		t.Errorf("a draft improvement that links its issue: %+v", d)
	}
	cod := d.CostOfDelay
	if cod == nil || cod.Inputs == nil {
		t.Fatalf("no cost of delay:\n%s", d.Body)
	}
	in := cod.Inputs
	if in.RevenuePerWeek == nil || *in.RevenuePerWeek != 1200 || in.PenaltyPerWeek == nil || *in.PenaltyPerWeek != 0 || in.TimeLostPerCycle != "1h30m" {
		t.Errorf("every figure of the impact carried over: %+v", *in)
	}
	if in.By != "flai" || in.At != "2026-09-16T11:00:00Z" || cod.Value != nil || cod.By != "" {
		t.Errorf("the inputs set by flai when the story was made (ADR-0080): %+v %+v", cod, *in)
	}
	lintStory(t, storyLint(t), d)
}

// S-0224: a duplicate finding bumps the open issue; the story made from it
// links both reports, in the order they found it, and carries the latest
// impact, which the bump replaced, beside the figures the bump left.
func TestForStoryFromABumpedAnalyzerIssue(t *testing.T) {
	r := repo(t)
	later := "design/analysis/2026-10-13-bottlenecks.md"
	if _, _, err := NewOrBump(r, NewOptions{Title: "Review waits a day", Class: "defect", Cost: "1h",
		Impact: Impact{RevenuePerWeek: "1200", TimeLostPerCycle: "4h", Evidence: "Twelve stories waited."}, Report: report, Now: t0}); err != nil {
		t.Fatal(err)
	}
	is, outcome, err := NewOrBump(r, NewOptions{Title: "Review waits a day", Class: "defect", Cost: "3h", Note: "Still waiting.",
		Impact: Impact{TimeLostPerCycle: "6h", Evidence: "Fifteen stories waited."}, Report: later, Now: t0.Add(time.Hour)})
	if err != nil || outcome != Bumped {
		t.Fatalf("the duplicate is bumped: %v %s", err, outcome)
	}
	if err := LinkStory(is, "S-0301", t0.Add(2*time.Hour)); err != nil {
		t.Fatal(err)
	}
	back, err := Read(is.Path)
	if err != nil {
		t.Fatal(err)
	}
	if got := strings.Join(Reports(back), " "); got != report+" "+later {
		t.Errorf("both reports, once each, in order: %s", got)
	}
	d := ForStory(back, t0.Add(2*time.Hour), cycle)
	goal := "## Goal\n\nThis story remediates [I-0001](../../../design/issues/I-0001-review-waits-a-day.md), \"Review waits a day\", " +
		"which the analyses in [2026-10-06-bottlenecks.md](../../../design/analysis/2026-10-06-bottlenecks.md) and " +
		"[2026-10-13-bottlenecks.md](../../../design/analysis/2026-10-13-bottlenecks.md) found. " +
		"The issue recommends no solution yet: propose one from its instances before building it.\n\n## Acceptance criteria\n"
	if !strings.HasPrefix(d.Body, goal) || d.Nature != "remediation" {
		t.Errorf("goal links the issue and both reports, and no report link as a solution:\n%s", d.Body)
	}
	want := "The analyzer found I-0001: the analyses the goal links give its evidence.\n\n" +
		"Cost of delay inputs set by flai from I-0001. revenue_per_week 1200 and time_lost_per_cycle 6h carried over from I-0001's Impact section. " +
		"Its Impact time_lost_per_cycle was taken rather than the 4h derived from its cost and count."
	if got := notes(d); got != want {
		t.Errorf("notes:\n got %s\nwant %s", got, want)
	}
	in := d.CostOfDelay.Inputs
	if in.RevenuePerWeek == nil || *in.RevenuePerWeek != 1200 || in.PenaltyPerWeek != nil || in.TimeLostPerCycle != "6h" || in.By != "flai" {
		t.Errorf("the latest impact wins, the figure it left is kept: %+v", *in)
	}
	lintStory(t, storyLint(t), d)
}

// S-0224: the solution an issue recommends is carried over without the lines
// linking its reports, wherever they fall in the section, and a report linked
// only under Remediation is still linked from the goal.
func TestForStoryLeavesTheReportLinksOutOfTheSolution(t *testing.T) {
	is := &Issue{ID: "I-0007", Title: "Slow builds", Class: "efficiency", Path: "/x/design/issues/I-0007-slow-builds.md",
		Body: "\n# I-0007 Slow builds\n\n## Instances\n\n### 2026-09-16T10:00:00Z\nfirst\n\n## Remediation\n\nCache the modules.\n\n" +
			"Found by the analysis in [a.md](../analysis/a.md).\n\nAnd in CI too.\n\n" +
			"Story S-0301 remediates this issue, created from it at 2026-09-16T11:00:00Z.\n\n" +
			"Found by the analysis in [b.md](../analysis/b.md).\n"}
	d := ForStory(is, t0, cycle)
	if !strings.Contains(d.Body, "\"Slow builds\", which the analyses in [a.md](../../../design/analysis/a.md) and [b.md](../../../design/analysis/b.md) found. The issue recommends this solution:\n\nCache the modules.\n\nAnd in CI too.\n\n## Acceptance criteria\n") {
		t.Errorf("the solution without the report links or the earlier story:\n%s", d.Body)
	}
	if strings.Contains(d.Body, "Found by") || strings.Contains(d.Body, "(../analysis/") {
		t.Errorf("a link from the issue's folder was carried over:\n%s", d.Body)
	}
	lintStory(t, storyLint(t), d)

	// a Report line that leaves the project names no report
	is.Body = "\n# I-0007 Slow builds\n\n## Instances\n\n### 2026-09-16T10:00:00Z\nReport: ../elsewhere.md.\nfirst\n\n## Remediation\n"
	if got := Reports(is); len(got) != 0 {
		t.Errorf("a path out of the project: %v", got)
	}
}
