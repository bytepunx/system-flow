package workitem

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/mdlint"
)

var activityAt = time.Date(2026, 10, 3, 18, 0, 0, 0, time.UTC)

// repoLint is this repository's markdownlint configuration, which the
// template ships too.
func repoLint(t *testing.T) *mdlint.Config {
	t.Helper()
	cfg, err := mdlint.Load(filepath.Join("..", "..", ".."))
	if err != nil || cfg == nil {
		t.Fatalf("the repository's markdownlint configuration: %v", err)
	}
	return cfg
}

func mustLintClean(t *testing.T, path string) {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if f := repoLint(t).Lint(string(data)); len(f) != 0 {
		t.Errorf("%s is not lint clean: %v\n%s", path, f, data)
	}
}

// ADR-0079: the first activity creates the document, later ones append and
// accrue, and the log reads back as it was written.
func TestAppendActivity(t *testing.T) {
	r := lintProject(t)
	first := ActivityEntry{At: activityAt, Summary: "Drafted five stories\n  for E-0016.", Items: []string{"E-0016", " S-0230 ", ""}, Seconds: 723, Cost: 0.42134, Estimated: true}
	a, err := r.AppendActivity(ActivityPlanner, first)
	if err != nil {
		t.Fatal(err)
	}
	if a.TasksCompleted != 1 || a.AccruedSeconds != 723 || a.AccruedCost != 0.4213 || a.LastRun != "2026-10-03T18:00:00Z" {
		t.Errorf("totals after one: %+v", a)
	}
	data, _ := os.ReadFile(r.ActivityPath(ActivityPlanner))
	want := "### 2026-10-03T18:00:00Z\n\n- Summary: Drafted five stories for E-0016.\n- Items: E-0016, S-0230\n- Seconds: 723\n- Cost: 0.4213 USD, estimated\n"
	if !strings.HasPrefix(string(data), "---\nkind: planner\n") || !strings.Contains(string(data), "\n# Planner activity\n") || !strings.HasSuffix(string(data), want) {
		t.Errorf("document:\n%s", data)
	}

	second := ActivityEntry{At: activityAt.Add(-time.Hour), Summary: "Sized S-0231.", Seconds: 60, Cost: 0.1}
	if a, err = r.AppendActivity(ActivityPlanner, second); err != nil {
		t.Fatal(err)
	}
	if a.TasksCompleted != 2 || a.LastRun != "2026-10-03T18:00:00Z" {
		t.Errorf("totals after two (an earlier end keeps last_run): %+v", a)
	}
	third := ActivityEntry{At: activityAt, Summary: "Same second.", Items: []string{"S-0232"}, Seconds: 1, Cost: 0}
	if a, err = r.AppendActivity(ActivityPlanner, third); err != nil {
		t.Fatal(err)
	}
	if a.TasksCompleted != 3 || a.AccruedSeconds != 784 || a.AccruedCost != 0.5213 || a.LastRun != "2026-10-03T18:00:00Z" {
		t.Errorf("totals after three (an earlier end keeps last_run): %+v", a)
	}

	back, err := ReadActivity(r.ActivityPath(ActivityPlanner))
	if err != nil {
		t.Fatal(err)
	}
	wantEntries := []ActivityEntry{
		{At: activityAt, Summary: "Drafted five stories for E-0016.", Items: []string{"E-0016", "S-0230"}, Seconds: 723, Cost: 0.4213, Estimated: true},
		{At: activityAt.Add(-time.Hour), Summary: "Sized S-0231.", Seconds: 60, Cost: 0.1},
		{At: activityAt, Summary: "Same second.", Items: []string{"S-0232"}, Seconds: 1},
	}
	if !reflect.DeepEqual(back.Entries, wantEntries) || !reflect.DeepEqual(back.Entries, a.Entries) {
		t.Errorf("log read back:\n got %+v\nwant %+v", back.Entries, wantEntries)
	}
	if back.AccruedCost != a.AccruedCost || back.AccruedSeconds != a.AccruedSeconds || back.TasksCompleted != 3 || back.LastRun != a.LastRun {
		t.Errorf("front matter read back: %+v", back)
	}
	data, _ = os.ReadFile(r.ActivityPath(ActivityPlanner))
	if !strings.Contains(string(data), "\n### 2026-10-03T18:00:00Z (2)\n") {
		t.Errorf("a second entry in the same second needs a heading of its own:\n%s", data)
	}
	mustLintClean(t, r.ActivityPath(ActivityPlanner))
}

// ADR-0084: a planner run's entry names what started it on a Trigger line
// after its summary, made one line, and reads back with it; an entry without
// one, as older entries and an agent's own are, reads back without it.
func TestActivityTrigger(t *testing.T) {
	r := lintProject(t)
	with := ActivityEntry{At: activityAt, Summary: "Planned S-0230.", Trigger: "edited goal by alex;\n  schedule daily", Items: []string{"S-0230"}, Seconds: 60, Cost: 0.1}
	without := ActivityEntry{At: activityAt.Add(time.Minute), Summary: "Logged by hand.", Seconds: 1}
	for _, e := range []ActivityEntry{with, without} {
		if _, err := r.AppendActivity(ActivityPlanner, e); err != nil {
			t.Fatal(err)
		}
	}
	data, _ := os.ReadFile(r.ActivityPath(ActivityPlanner))
	want := "### 2026-10-03T18:00:00Z\n\n- Summary: Planned S-0230.\n- Trigger: edited goal by alex; schedule daily\n- Items: S-0230\n- Seconds: 60\n- Cost: 0.1000 USD\n\n" +
		"### 2026-10-03T18:01:00Z\n\n- Summary: Logged by hand.\n- Items: none\n- Seconds: 1\n- Cost: 0.0000 USD\n"
	if !strings.HasSuffix(string(data), want) {
		t.Errorf("document:\n%s\nwant it to end\n%s", data, want)
	}
	back, err := ReadActivity(r.ActivityPath(ActivityPlanner))
	if err != nil {
		t.Fatal(err)
	}
	wantEntries := []ActivityEntry{
		{At: activityAt, Summary: "Planned S-0230.", Trigger: "edited goal by alex; schedule daily", Items: []string{"S-0230"}, Seconds: 60, Cost: 0.1},
		{At: activityAt.Add(time.Minute), Summary: "Logged by hand.", Seconds: 1},
	}
	if !reflect.DeepEqual(back.Entries, wantEntries) {
		t.Errorf("log read back:\n got %+v\nwant %+v", back.Entries, wantEntries)
	}
	mustLintClean(t, r.ActivityPath(ActivityPlanner))
}

// S-0218: a call flai guard refused the orchestrator is logged under
// ## Refusals, after the log, with its time, the call, and the permission
// that would allow it, and no seconds or cost; the totals count activities
// alone, and an activity logged later goes in the log above.
func TestAppendRefusal(t *testing.T) {
	r := lintProject(t)
	if _, err := r.AppendActivity(ActivityOrchestrator, ActivityEntry{At: activityAt, Summary: "Ordered the ready column by WSJF.", Seconds: 60, Cost: 0.1}); err != nil {
		t.Fatal(err)
	}
	for _, f := range []ActivityRefusal{
		{At: activityAt.Add(time.Minute), Call: "flai move S-0001 ready", Needs: "orchestration.permissions.promote_to_ready"},
		{At: activityAt.Add(time.Minute), Call: "Edit flai/cmd/x.go"},
		{At: activityAt.Add(2 * time.Minute), Call: "echo `date` &&\n  git   commit -m x"},
	} {
		if _, err := r.AppendRefusal(ActivityOrchestrator, f); err != nil {
			t.Fatal(err)
		}
	}
	a, err := r.AppendActivity(ActivityOrchestrator, ActivityEntry{At: activityAt.Add(3 * time.Minute), Summary: "Asked about S-0001 on a thread.", Items: []string{"S-0001"}, Seconds: 30, Cost: 0.05})
	if err != nil {
		t.Fatal(err)
	}
	if a.TasksCompleted != 2 || a.AccruedSeconds != 90 || a.AccruedCost != 0.15 || a.LastRun != "2026-10-03T18:03:00Z" {
		t.Errorf("totals count the activities alone: %+v", a)
	}
	data, _ := os.ReadFile(r.ActivityPath(ActivityOrchestrator))
	want := "- Summary: Asked about S-0001 on a thread.\n- Items: S-0001\n- Seconds: 30\n- Cost: 0.0500 USD\n\n## Refusals\n\n" +
		"### 2026-10-03T18:01:00Z\n\n- Call: `flai move S-0001 ready`\n- Needs: orchestration.permissions.promote_to_ready\n\n" +
		"### 2026-10-03T18:01:00Z (2)\n\n- Call: `Edit flai/cmd/x.go`\n- Needs: none\n\n" +
		"### 2026-10-03T18:02:00Z\n\n- Call: `echo 'date' && git commit -m x`\n- Needs: none\n"
	if !strings.HasSuffix(string(data), want) {
		t.Errorf("document:\n%s\nwant it to end\n%s", data, want)
	}
	back, err := ReadActivity(r.ActivityPath(ActivityOrchestrator))
	if err != nil {
		t.Fatal(err)
	}
	wantRefusals := []ActivityRefusal{
		{At: activityAt.Add(time.Minute), Call: "flai move S-0001 ready", Needs: "orchestration.permissions.promote_to_ready"},
		{At: activityAt.Add(time.Minute), Call: "Edit flai/cmd/x.go"},
		{At: activityAt.Add(2 * time.Minute), Call: "echo 'date' && git commit -m x"},
	}
	if !reflect.DeepEqual(back.Refusals, wantRefusals) || !reflect.DeepEqual(a.Refusals, wantRefusals) || len(back.Entries) != 2 {
		t.Errorf("read back:\n got %+v %+v\nwant %+v", back.Entries, back.Refusals, wantRefusals)
	}
	mustLintClean(t, r.ActivityPath(ActivityOrchestrator))
}

// S-0218: a refusal before any activity creates the document with an empty
// log and zero totals.
func TestAppendRefusalCreatesTheDocument(t *testing.T) {
	r := lintProject(t)
	a, err := r.AppendRefusal(ActivityOrchestrator, ActivityRefusal{At: activityAt, Call: "flai accept S-0001", Needs: "orchestration.permissions.accept_reviews"})
	if err != nil {
		t.Fatal(err)
	}
	if a.TasksCompleted != 0 || a.LastRun != "" || len(a.Entries) != 0 || len(a.Refusals) != 1 {
		t.Errorf("document: %+v", a)
	}
	data, _ := os.ReadFile(r.ActivityPath(ActivityOrchestrator))
	if !strings.HasSuffix(string(data), "## Log\n\n## Refusals\n\n### 2026-10-03T18:00:00Z\n\n- Call: `flai accept S-0001`\n- Needs: orchestration.permissions.accept_reviews\n") {
		t.Errorf("document:\n%s", data)
	}
	mustLintClean(t, r.ActivityPath(ActivityOrchestrator))
}

func TestAppendRefusalRefuses(t *testing.T) {
	r := newProject(t)
	for name, c := range map[string]struct {
		kind string
		f    ActivityRefusal
		want string
	}{
		"unknown kind": {"designer", ActivityRefusal{At: activityAt, Call: "x"}, "use one of planner, orchestrator, analyzer"},
		"no call":      {ActivityOrchestrator, ActivityRefusal{At: activityAt, Call: " \n "}, "needs the call"},
		"no time":      {ActivityOrchestrator, ActivityRefusal{Call: "x"}, "time of the call"},
		"bad needs":    {ActivityOrchestrator, ActivityRefusal{At: activityAt, Call: "x", Needs: "the publish one"}, "one word"},
	} {
		if _, err := r.AppendRefusal(c.kind, c.f); err == nil || !strings.Contains(err.Error(), c.want) {
			t.Errorf("%s: want an error saying %q, got %v", name, c.want, err)
		}
	}
	if m, _ := filepath.Glob(filepath.Join(r.AgentsDir(), "*.md")); len(m) != 0 {
		t.Errorf("nothing is written: %v", m)
	}
}

func TestReadActivityRefusalsAreStrict(t *testing.T) {
	r := newProject(t)
	if _, err := r.AppendRefusal(ActivityOrchestrator, ActivityRefusal{At: activityAt, Call: "flai push", Needs: "orchestration.permissions.publish"}); err != nil {
		t.Fatal(err)
	}
	path := r.ActivityPath(ActivityOrchestrator)
	good, _ := os.ReadFile(path)
	for name, c := range map[string]struct{ from, to, want string }{
		"stray line": {"- Needs: orchestration.permissions.publish\n", "- Needs: orchestration.permissions.publish\nby hand\n", "not one of Call or Needs"},
		"no needs":   {"- Needs: orchestration.permissions.publish\n", "", "Needs"},
		"bare call":  {"- Call: `flai push`\n", "- Call: flai push\n", "one code span"},
		"two calls":  {"- Call: `flai push`\n", "- Call: `flai push`\n- Call: `git push`\n", "two \"- Call:\" lines"},
		"no heading": {"### 2026-10-03T18:00:00Z\n", "", "under ## Refusals is not an entry heading"},
		"a cost":     {"- Needs: orchestration.permissions.publish\n", "- Needs: orchestration.permissions.publish\n- Cost: 0.0000 USD\n", "not one of Call or Needs"},
	} {
		if err := os.WriteFile(path, []byte(strings.Replace(string(good), c.from, c.to, 1)), 0o644); err != nil {
			t.Fatal(err)
		}
		if _, err := ReadActivity(path); err == nil || !strings.Contains(err.Error(), c.want) {
			t.Errorf("%s: want an error naming %q, got %v", name, c.want, err)
		}
	}
}

// S-0218: writers of one activity document take turns under its lock, so
// none loses another's entry, and a lock left behind by a flai that died is
// taken over once it is old.
func TestActivityWritersTakeTurns(t *testing.T) {
	r := newProject(t)
	lock := filepath.Join(r.CacheDir(), "activity-orchestrator.lock")
	if err := os.MkdirAll(filepath.Dir(lock), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(lock, nil, 0o600); err != nil {
		t.Fatal(err)
	}
	old := time.Now().Add(-time.Minute)
	if err := os.Chtimes(lock, old, old); err != nil {
		t.Fatal(err)
	}
	var wg sync.WaitGroup
	for i := range 12 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			var err error
			if i%2 == 0 {
				_, err = r.AppendRefusal(ActivityOrchestrator, ActivityRefusal{At: activityAt, Call: fmt.Sprintf("flai accept S-%04d", i)})
			} else {
				_, err = r.AppendActivity(ActivityOrchestrator, ActivityEntry{At: activityAt, Summary: fmt.Sprintf("Activity %d.", i), Seconds: 1})
			}
			if err != nil {
				t.Error(err)
			}
		}()
	}
	wg.Wait()
	a, err := ReadActivity(r.ActivityPath(ActivityOrchestrator))
	if err != nil {
		t.Fatal(err)
	}
	if len(a.Entries) != 6 || len(a.Refusals) != 6 || a.TasksCompleted != 6 || a.AccruedSeconds != 6 {
		t.Errorf("an entry was lost: %d entries, %d refusals, %+v", len(a.Entries), len(a.Refusals), a)
	}
	if _, err := os.Stat(lock); !os.IsNotExist(err) {
		t.Errorf("the lock is left behind: %v", err)
	}
}

func TestActivityMissingIsEmpty(t *testing.T) {
	r := newProject(t)
	a, err := r.Activity(ActivityAnalyzer)
	if err != nil {
		t.Fatal(err)
	}
	if a.Kind != ActivityAnalyzer || a.Path != r.ActivityPath(ActivityAnalyzer) || a.TasksCompleted != 0 || len(a.Entries) != 0 {
		t.Errorf("missing document: %+v", a)
	}
	if list, err := r.Activities(); err != nil || len(list) != 0 {
		t.Errorf("no documents listed: %v %v", list, err)
	}
}

func TestAppendActivityRefuses(t *testing.T) {
	r := newProject(t)
	for name, c := range map[string]struct {
		kind string
		e    ActivityEntry
		want string
	}{
		"unknown kind":  {"designer", ActivityEntry{At: activityAt, Summary: "x"}, "use one of planner, orchestrator, analyzer"},
		"empty summary": {ActivityOrchestrator, ActivityEntry{At: activityAt, Summary: " \n "}, "needs a summary"},
		"no end":        {ActivityOrchestrator, ActivityEntry{Summary: "x"}, "time it ended"},
		"negative":      {ActivityOrchestrator, ActivityEntry{At: activityAt, Summary: "x", Seconds: -1}, "zero or more"},
		"bad item":      {ActivityOrchestrator, ActivityEntry{At: activityAt, Summary: "x", Items: []string{"S-1, S-2"}}, "no commas or spaces"},
	} {
		if _, err := r.AppendActivity(c.kind, c.e); err == nil || !strings.Contains(err.Error(), c.want) {
			t.Errorf("%s: want an error saying %q, got %v", name, c.want, err)
		}
	}
	if _, err := r.Activity("designer"); err == nil {
		t.Error("an unknown kind is refused when read")
	}
	if m, _ := filepath.Glob(filepath.Join(r.AgentsDir(), "*.md")); len(m) != 0 {
		t.Errorf("nothing is written: %v", m)
	}
}

func TestReadActivityIsStrict(t *testing.T) {
	r := newProject(t)
	if _, err := r.AppendActivity(ActivityAnalyzer, ActivityEntry{At: activityAt, Summary: "Read the stats."}); err != nil {
		t.Fatal(err)
	}
	path := r.ActivityPath(ActivityAnalyzer)
	good, _ := os.ReadFile(path)
	for name, c := range map[string]struct{ from, to, want string }{
		"unknown key":   {"kind: analyzer\n", "kind: analyzer\nmood: fine\n", "mood"},
		"wrong kind":    {"kind: analyzer\n", "kind: planner\n", "file's name"},
		"bad last_run":  {"last_run: 2026-10-03T18:00:00Z\n", "last_run: yesterday\n", "last_run"},
		"stray line":    {"- Seconds: 0\n", "- Seconds: 0\nby hand\n", "not one of"},
		"missing field": {"- Seconds: 0\n", "", "Seconds"},
		"two triggers":  {"- Items: none\n", "- Trigger: asked\n- Trigger: reordered\n- Items: none\n", "two \"- Trigger:\" lines"},
		"bad cost":      {"- Cost: 0.0000 USD\n", "- Cost: free\n", "dollars"},
	} {
		if err := os.WriteFile(path, []byte(strings.Replace(string(good), c.from, c.to, 1)), 0o644); err != nil {
			t.Fatal(err)
		}
		if _, err := ReadActivity(path); err == nil || !strings.Contains(err.Error(), c.want) {
			t.Errorf("%s: want an error naming %q, got %v", name, c.want, err)
		}
		if _, err := r.AppendActivity(ActivityAnalyzer, ActivityEntry{At: activityAt, Summary: "x"}); err == nil {
			t.Errorf("%s: appending to a document that does not parse is refused", name)
		}
	}
}

func TestIsActivityDocument(t *testing.T) {
	for name, want := range map[string]bool{
		"planner.md": true, "wip/agents/orchestrator.md": true, "analyzer.md": true,
		"S-0206.md": false, "index.md": false, "README.md": false, "planner": false,
	} {
		if IsActivityDocument(name) != want {
			t.Errorf("IsActivityDocument(%q) = %v", name, !want)
		}
	}
}

// The index lists the activity documents that exist under their own
// heading, and the narratives list leaves them out.
func TestIndexListsActivityDocuments(t *testing.T) {
	r := lintProject(t)
	index := filepath.Join(r.AgentsDir(), "index.md")
	if err := r.WriteIndex(nil, activityAt); err != nil {
		t.Fatal(err)
	}
	if data, _ := os.ReadFile(index); strings.Contains(string(data), "Strategic agents") {
		t.Errorf("no documents, no heading:\n%s", data)
	}

	s, err := r.Create(NewOptions{Type: Story, Title: "Stream me", Owner: "alex", Now: activityAt})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := r.OpenStream(s, StreamOptions{Agent: "bot", Now: activityAt}); err != nil {
		t.Fatal(err)
	}
	for _, kind := range []string{ActivityOrchestrator, ActivityPlanner} {
		if _, err := r.AppendActivity(kind, ActivityEntry{At: activityAt, Summary: "Did one thing.", Seconds: 90, Cost: 0.25}); err != nil {
			t.Fatal(err)
		}
	}
	if ns, err := r.ActiveNarratives(); err != nil || len(ns) != 1 {
		t.Fatalf("narratives are the stories' only: %v %v", ns, err)
	}
	items, _ := r.List(false)
	if err := r.WriteIndex(items, activityAt); err != nil {
		t.Fatal(err)
	}
	data, _ := os.ReadFile(index)
	want := "\n## Strategic agents\n\n| Agent | Activities | Cost | Seconds | Last run |\n|-------|------------|------|---------|----------|\n" +
		"| [planner](planner.md) | 1 | 0.2500 USD | 90 | 2026-10-03T18:00:00Z |\n" +
		"| [orchestrator](orchestrator.md) | 1 | 0.2500 USD | 90 | 2026-10-03T18:00:00Z |\n"
	if !strings.HasSuffix(string(data), want) || !strings.Contains(string(data), "| [S-0001](S-0001.md) |") {
		t.Errorf("index:\n%s", data)
	}
	mustLintClean(t, index)
	mustLintClean(t, r.ActivityPath(ActivityOrchestrator))
}

// A document that does not parse is listed as unreadable, and the index is
// still written: every move rewrites it.
func TestIndexListsAnUnreadableActivityDocument(t *testing.T) {
	r := lintProject(t)
	if err := os.MkdirAll(r.AgentsDir(), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(r.ActivityPath(ActivityAnalyzer), []byte("---\nkind: planner\n---\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := r.WriteIndex(nil, activityAt); err != nil {
		t.Fatal(err)
	}
	data, _ := os.ReadFile(filepath.Join(r.AgentsDir(), "index.md"))
	if !strings.Contains(string(data), "| [analyzer](analyzer.md) | unreadable, see flai check | | | |\n") {
		t.Errorf("index:\n%s", data)
	}
	mustLintClean(t, filepath.Join(r.AgentsDir(), "index.md"))
}
