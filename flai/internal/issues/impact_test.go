package issues

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

const report = "design/analysis/2026-10-06-bottlenecks.md"

// S-0224: the analyzer's finding is written with its Impact section, which
// impact reads back as the cost of delay inputs of the story made from it,
// and links the report that found it.
func TestANewIssueWritesItsImpactAndReport(t *testing.T) {
	r := repo(t)
	im := Impact{RevenuePerWeek: "1200", PenaltyPerWeek: "0", TimeLostPerCycle: "90m", Evidence: "12 stories waited 18h on average in review."}
	is, err := New(r, NewOptions{Title: "Review waits a day", Class: "efficiency", Story: "S-0224", Impact: im, Report: report, Now: t0})
	if err != nil {
		t.Fatal(err)
	}
	raw, _ := os.ReadFile(is.Path)
	want := "## Instances\n\n### 2026-09-16T10:00:00Z\nStory: S-0224.\nReport: " + report + ".\nFirst occurrence.\n\n" +
		"## Impact\n12 stories waited 18h on average in review.\n\n- revenue_per_week: 1200\n- penalty_per_week: 0\n- time_lost_per_cycle: 1h30m\n\n" +
		"## Remediation\n\nFound by the analysis in [2026-10-06-bottlenecks.md](../analysis/2026-10-06-bottlenecks.md).\n"
	if !strings.HasSuffix(string(raw), want) {
		t.Errorf("the issue should hold its impact and link its report:\n%s", raw)
	}
	lintIssue(t, is)
	back, err := Read(is.Path)
	if err != nil {
		t.Fatal(err)
	}
	in, carried, skipped := impact(back)
	if in.RevenuePerWeek == nil || *in.RevenuePerWeek != 1200 || in.PenaltyPerWeek == nil || *in.PenaltyPerWeek != 0 || in.TimeLostPerCycle != "1h30m" || len(skipped) != 0 {
		t.Errorf("impact reads back what New wrote: %+v %v %v", in, carried, skipped)
	}
	if got := ReportLinks(back); len(got) != 1 || got[0] != "../analysis/2026-10-06-bottlenecks.md" {
		t.Errorf("report links: %v", got)
	}
	d := ForStory(back, t0, cycle)
	if d.CostOfDelay == nil || d.CostOfDelay.Inputs.TimeLostPerCycle != "1h30m" || *d.CostOfDelay.Inputs.RevenuePerWeek != 1200 {
		t.Errorf("the story made from it carries the inputs: %+v", d.CostOfDelay)
	}

	// without impact or report the body is as it always was
	plain, err := New(r, NewOptions{Title: "Plain", Class: "defect", Now: t0})
	if err != nil {
		t.Fatal(err)
	}
	if plain.Body != "\n# I-0002 Plain\n\n## Description\nPlain\n\n## Instances\n\n### 2026-09-16T10:00:00Z\nFirst occurrence.\n\n## Remediation\n" {
		t.Errorf("an issue with no impact or report:\n%q", plain.Body)
	}
}

// Bad amounts, durations, and report paths are refused before anything is
// written.
func TestImpactAndReportRefusals(t *testing.T) {
	r := repo(t)
	is, err := New(r, NewOptions{Title: "Kept", Class: "defect", Now: t0})
	if err != nil {
		t.Fatal(err)
	}
	before, _ := os.ReadFile(is.Path)
	for _, c := range []struct {
		name   string
		im     Impact
		report string
		want   string
	}{
		{"a word for an amount", Impact{RevenuePerWeek: "lots"}, "", "revenue_per_week \"lots\" is not an amount"},
		{"a negative amount", Impact{PenaltyPerWeek: "-5"}, "", "penalty_per_week \"-5\" is not an amount"},
		{"not a number", Impact{PenaltyPerWeek: "NaN"}, "", "not an amount"},
		{"a duration of zero", Impact{TimeLostPerCycle: "0s"}, "", "time_lost_per_cycle \"0s\" is not a duration longer than zero"},
		{"not a duration", Impact{TimeLostPerCycle: "a day"}, "", "not a duration"},
		{"a negative duration", Impact{TimeLostPerCycle: "-1h"}, "", "not a duration"},
		{"a report outside design/analysis", Impact{}, "design/issues/x.md", "not an analysis report"},
		{"a report that climbs out", Impact{}, "design/analysis/../x.md", "not an analysis report"},
		{"a report that is not markdown", Impact{}, "design/analysis/x.txt", "not an analysis report"},
		{"the reports' index", Impact{}, "design/analysis/README.md", "not an analysis report"},
		{"the folder itself", Impact{}, "design/analysis", "not an analysis report"},
		{"a report outside the project", Impact{}, "/elsewhere/design/analysis/x.md", "not an analysis report"},
	} {
		if _, err := New(r, NewOptions{Title: "Refused", Class: "defect", Impact: c.im, Report: c.report, Now: t0}); err == nil || !strings.Contains(err.Error(), c.want) {
			t.Errorf("new, %s: %v", c.name, err)
		}
		if _, _, err := NewOrBump(r, NewOptions{Title: "Kept", Class: "defect", Impact: c.im, Report: orReport(c.report), Now: t0}); err == nil || !strings.Contains(err.Error(), c.want) {
			t.Errorf("new or bump, %s: %v", c.name, err)
		}
		if err := BumpWith(r, is, BumpOptions{Impact: c.im, Report: c.report, Now: t0.Add(time.Hour)}); err == nil || !strings.Contains(err.Error(), c.want) {
			t.Errorf("bump, %s: %v", c.name, err)
		}
	}
	if is.Count != 1 {
		t.Errorf("a refused bump changed the issue: count %d", is.Count)
	}
	if after, _ := os.ReadFile(is.Path); string(after) != string(before) {
		t.Errorf("a refusal wrote the issue:\n%s", after)
	}
	if m, _ := filepath.Glob(filepath.Join(Dir(r), "I-*.md")); len(m) != 1 {
		t.Errorf("a refusal wrote an issue: %v", m)
	}
	abs := filepath.Join(r.Root, filepath.FromSlash(report))
	if got, err := reportPath(r, abs); err != nil || got != report {
		t.Errorf("an absolute path in the project is taken: %q %v", got, err)
	}
	if got, err := reportPath(r, "./"+report); err != nil || got != report {
		t.Errorf("a path from the root is cleaned: %q %v", got, err)
	}
}

// lintIssue fails the test when the issue's body is not lint-clean.
func lintIssue(t *testing.T, is *Issue) {
	t.Helper()
	doc := "---\nid: " + is.ID + "\n---\n" + is.Body
	if f := storyLint(t).Lint(doc); len(f) > 0 {
		t.Errorf("lint: %v\n%s", f, doc)
	}
}

// orReport keeps a refusal's report, or gives a good one, so that NewOrBump
// takes the path that bumps.
func orReport(r string) string {
	if r == "" {
		return report
	}
	return r
}

// With a report, an open issue of the same title is bumped, not duplicated;
// its Impact lines given are replaced and the others kept, and the report is
// linked once. Filing the same report again changes nothing.
func TestNewOrBumpBumpsTheOpenIssueOfTheSameTitle(t *testing.T) {
	r := repo(t)
	first, outcome, err := NewOrBump(r, NewOptions{Title: "Review waits a day", Class: "efficiency", Cost: "1h",
		Impact: Impact{RevenuePerWeek: "1200", TimeLostPerCycle: "4h", Evidence: "Twelve stories waited."}, Report: report, Now: t0})
	if err != nil || outcome != Opened || first.ID != "I-0001" {
		t.Fatalf("first: %v %s %+v", err, outcome, first)
	}
	later := "design/analysis/2026-10-13-bottlenecks.md"
	is, outcome, err := NewOrBump(r, NewOptions{Title: "Review waits a day", Class: "efficiency", Cost: "3h", Note: "Still waiting.", Story: "S-0300",
		Impact: Impact{TimeLostPerCycle: "6h", Evidence: "Fifteen stories waited."}, Report: later, Now: t0.Add(time.Hour)})
	if err != nil || outcome != Bumped || is.ID != "I-0001" || is.Count != 2 || is.Cost != "2h" {
		t.Fatalf("same title: %v %s %+v", err, outcome, is)
	}
	if m, _ := filepath.Glob(filepath.Join(Dir(r), "I-*.md")); len(m) != 1 {
		t.Errorf("a duplicate was made: %v", m)
	}
	raw, _ := os.ReadFile(is.Path)
	want := "### 2026-09-16T11:00:00Z\nStory: S-0300.\nReport: " + later + ".\nStill waiting.\n\n" +
		"## Impact\nTwelve stories waited.\n\nFifteen stories waited.\n\n- revenue_per_week: 1200\n- time_lost_per_cycle: 6h\n\n" +
		"## Remediation\n\nFound by the analysis in [2026-10-06-bottlenecks.md](../analysis/2026-10-06-bottlenecks.md).\n\n" +
		"Found by the analysis in [2026-10-13-bottlenecks.md](../analysis/2026-10-13-bottlenecks.md).\n"
	if !strings.HasSuffix(string(raw), want) {
		t.Errorf("the bump should add its instance before Impact, replace time lost, and link the second report:\n%s", raw)
	}
	again, outcome, err := NewOrBump(r, NewOptions{Title: "Review waits a day", Class: "efficiency", Report: later, Now: t0.Add(2 * time.Hour)})
	if err != nil || outcome != Already || again.Count != 2 {
		t.Errorf("the same report again: %v %s %+v", err, outcome, again)
	}
	if after, _ := os.ReadFile(is.Path); string(after) != string(raw) {
		t.Errorf("filing the same report again wrote the issue:\n%s", after)
	}
	// an explicit bump with a report it links already does not link it twice
	if err := BumpWith(r, is, BumpOptions{Report: later, Impact: Impact{PenaltyPerWeek: "300"}, Now: t0.Add(3 * time.Hour)}); err != nil {
		t.Fatal(err)
	}
	if got := ReportLinks(is); len(got) != 2 {
		t.Errorf("a report linked twice: %v", got)
	}
	in, _, _ := impact(is)
	if *in.RevenuePerWeek != 1200 || *in.PenaltyPerWeek != 300 || in.TimeLostPerCycle != "6h" {
		t.Errorf("the impact after the bump: %+v", in)
	}
	lintIssue(t, is)
	// a new title is a new issue, and without a report the same title is too
	if is, outcome, err := NewOrBump(r, NewOptions{Title: "Another", Class: "defect", Report: report, Now: t0}); err != nil || outcome != Opened || is.ID != "I-0002" {
		t.Errorf("another title: %v %s %+v", err, outcome, is)
	}
	if is, outcome, err := NewOrBump(r, NewOptions{Title: "Another", Class: "defect", Now: t0}); err != nil || outcome != Opened || is.ID != "I-0003" {
		t.Errorf("no report: %v %s %+v", err, outcome, is)
	}
}

// A bump of an issue with no Impact section adds one before Remediation, and
// an instance recorded after it still goes in Instances.
func TestBumpAddsAnImpactSection(t *testing.T) {
	r := repo(t)
	is, err := New(r, NewOptions{Title: "Old", Class: "defect", Now: t0})
	if err != nil {
		t.Fatal(err)
	}
	if err := BumpWith(r, is, BumpOptions{Impact: Impact{PenaltyPerWeek: "99.5"}, Now: t0.Add(time.Hour)}); err != nil {
		t.Fatal(err)
	}
	if err := Bump(is, "", "", "third", t0.Add(2*time.Hour)); err != nil {
		t.Fatal(err)
	}
	want := "### 2026-09-16T11:00:00Z\nOccurred again.\n\n### 2026-09-16T12:00:00Z\nthird\n\n## Impact\n- penalty_per_week: 99.5\n\n## Remediation\n"
	if !strings.HasSuffix(is.Body, want) {
		t.Errorf("impact added before Remediation, instances kept in Instances:\n%s", is.Body)
	}
	if err := BumpWith(nil, is, BumpOptions{Report: report, Now: t0}); err == nil {
		t.Error("a report with no project to find it in should be refused")
	}
}
