//go:build !windows

package serve

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// The analyzer and the curator start on the trigger entries their role
// files name, through their own schedulers (ADR-0250 §4).

// roleTriggers writes the project at root's role file for the built-in role
// name: the template's, with its trigger list given as entries, one a line.
func roleTriggers(t *testing.T, root, name string, entries ...string) {
	t.Helper()
	data, err := os.ReadFile(filepath.Join("..", "..", "..", "template", "root", "design", "workflows", "roles", name+".md"))
	if err != nil {
		t.Fatal(err)
	}
	s := string(data)
	start, end := strings.Index(s, "\ntrigger:\n"), strings.Index(s, "\nitem:")
	if start < 0 || end < start {
		t.Fatalf("the template's %s.md has no trigger list before item:\n%s", name, s)
	}
	s = s[:start] + "\ntrigger:\n  - " + strings.Join(entries, "\n  - ") + s[end:]
	dir := filepath.Join(root, "design", "workflows", "roles")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, name+".md"), []byte(s), 0o644); err != nil {
		t.Fatal(err)
	}
}

// ADR-0250 §4: a project's analyzer file without a schedule entry starts no
// analysis when analysis.schedule comes round.
func TestAnAnalyzerFileWithoutAScheduleStartsNoAnalysis(t *testing.T) {
	lab := newScheduleLab(t, "daily", time.Date(2026, 10, 6, 10, 0, 0, 0, time.UTC))
	roleTriggers(t, lab.root, "analyzer", "asked")
	lab.lookAt(lab.clock)
	lab.lookAt(time.Date(2026, 10, 7, 0, 0, 30, 0, time.UTC)) // midnight came round
	if n := lab.starts(); n != 0 || lab.state().Analyzer != nil || strings.Contains(lab.logText(), "scheduled analyzer") {
		t.Fatalf("the schedule started the analyzer with no entry for it: %d runs, %+v\n%s", n, lab.state().Analyzer, lab.logText())
	}
}

// ADR-0250 §4: the analyzer follows the setting its file's schedule entry
// names, and starts on it with the schedule as its trigger, while
// analysis.schedule, which the file no longer names, starts nothing.
func TestTheAnalyzerFollowsTheScheduleItsFileNames(t *testing.T) {
	lab := newScheduleLab(t, "0 6 * * *", time.Date(2026, 10, 6, 10, 0, 0, 0, time.UTC))
	m := "version: 1\nname: t\nkey: t\nlayout:\n  design: design\n  docs: docs\n  wip: wip\nanalysis:\n  schedule: \"0 6 * * *\"\ndigest:\n  schedule: daily\n"
	if err := os.WriteFile(filepath.Join(lab.root, "system-flow.yaml"), []byte(m), 0o644); err != nil {
		t.Fatal(err)
	}
	roleTriggers(t, lab.root, "analyzer", "asked", "{schedule: digest.schedule}")
	lab.lookAt(lab.clock)
	lab.lookAt(time.Date(2026, 10, 6, 23, 59, 0, 0, time.UTC))
	if n := lab.starts(); n != 0 {
		t.Fatalf("started before the setting came round: %d runs", n)
	}
	lab.lookAt(time.Date(2026, 10, 7, 0, 0, 30, 0, time.UTC))
	if run := lab.state().Analyzer; lab.starts() != 1 || run == nil || run.Trigger != "schedule daily" || run.Focus != "all" {
		t.Fatalf("at midnight: %d runs, %+v; want one for all three focuses started by schedule daily", lab.starts(), run)
	}
	lab.ended()
	lab.lookAt(time.Date(2026, 10, 7, 6, 0, 30, 0, time.UTC)) // analysis.schedule came round
	if n := lab.starts(); n != 1 {
		t.Errorf("analysis.schedule, which the file does not name, started %d runs in all, want one", n)
	}
}

// ADR-0250 §4: a project's curator file without {event: accepted} starts
// nothing on an acceptance, while the digest.schedule it names still starts
// the curator, with the schedule as its trigger.
func TestACuratorFileWithoutAcceptedStartsOnlyOnItsSchedule(t *testing.T) {
	lab := newCurationLab(t, "daily", time.Date(2026, 10, 6, 10, 0, 0, 0, time.UTC))
	roleTriggers(t, lab.root, "curator", "asked", "{schedule: digest.schedule}")
	lab.lookAt(lab.clock)
	lab.curatorRuns("Took in one entry.", runStart)
	lab.accept("Accepted")
	lab.lookAt(lab.clock.Add(time.Second))
	if n := len(lab.starts()); n != 0 || len(lab.s.pending) != 0 {
		t.Fatalf("an acceptance with no entry for it: %d runs, pending %v", n, lab.s.pending)
	}
	lab.lookAt(time.Date(2026, 10, 7, 0, 0, 30, 0, time.UTC)) // midnight came round
	starts := lab.starts()
	if run := lab.state().Curator; len(starts) != 1 || run == nil || run.Trigger != "schedule daily" {
		t.Fatalf("at midnight: %d runs, %+v; want one started by schedule daily", len(starts), run)
	}
	lab.ended(1)
}

// ADR-0250 §4: a project's curator file without a schedule entry starts
// nothing when digest.schedule comes round, while an acceptance still starts
// it.
func TestACuratorFileWithoutAScheduleStartsOnlyOnAcceptances(t *testing.T) {
	lab := newCurationLab(t, "daily", time.Date(2026, 10, 6, 10, 0, 0, 0, time.UTC))
	roleTriggers(t, lab.root, "curator", "asked", "{event: accepted}")
	lab.lookAt(lab.clock)
	lab.curatorRuns("Took in one entry.", runStart)
	lab.lookAt(time.Date(2026, 10, 7, 0, 0, 30, 0, time.UTC)) // midnight came round
	if n := len(lab.starts()); n != 0 {
		t.Fatalf("the schedule started the curator with no entry for it: %d runs", n)
	}
	story := lab.accept("Accepted")
	lab.lookAt(lab.clock.Add(time.Second))
	if run := lab.state().Curator; len(lab.starts()) != 1 || run == nil || run.Trigger != "accepted "+story {
		t.Fatalf("after the acceptance: %d runs, %+v; want one started by accepted %s", len(lab.starts()), run, story)
	}
	lab.ended(1)
}
