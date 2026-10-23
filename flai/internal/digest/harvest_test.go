package digest

import (
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// threadFile is a thread file anchored to S-0001, whose entries are given as
// "<ts> <author>" headings each followed by its text.
func threadFile(id, title, status string, entries ...string) string {
	var b strings.Builder
	b.WriteString("---\nid: " + id + "\ntitle: " + title + "\nanchor:\n  path: wip/archive/kanban/stories/S-0001-s.md\n  item: S-0001\nstatus: " + status + "\nparticipants: [claude, alex]\ncreated: 2026-10-19T10:00:00Z\nupdated: 2026-10-23T10:00:00Z\n---\n\n# " + id + " " + title + "\n\n## Entries\n")
	for i := 0; i+1 < len(entries); i += 2 {
		b.WriteString("\n### " + entries[i] + "\n" + entries[i+1] + "\n")
	}
	return b.String()
}

// harvestRepo is a project owned by alex with one candidate source of each
// kind, learned on different days:
//
//   - TH-0001, answered, two answers by alex on 2026-10-19;
//   - S-0001, archived, done 2026-10-20T12:00:00Z, one decision;
//   - I-0001, closed 2026-10-21T09:00:00Z by S-0001;
//   - TH-0002, resolved, alex overruling a recommendation on 2026-10-22.
func harvestRepo(t *testing.T) *workitem.Repo {
	t.Helper()
	root := t.TempDir()
	write(t, root, "system-flow.yaml", "version: 1\nname: t\nkey: t\nowner: alex\nlayout:\n  design: design\n  docs: docs\n  wip: wip\n")
	write(t, root, "wip/archive/kanban/stories/S-0001-s.md", item("S-0001", "story", "done", "", "2026-10-20T12:00:00Z", ""))
	write(t, root, "wip/archive/agents/S-0001.md", narrative("S-0001", "## Decisions\n\n- Keep one file per entry. It merges cleanly.\n"))
	write(t, root, "wip/agents/index.md", "# Narratives\n")
	write(t, root, "design/issues/I-0001-a.md", issue("I-0001", "closed", "\n## Description\nTwo writers took one number.\n\n## Remediation\nClosed 2026-10-21T09:00:00Z: S-0001: numbers come from every branch.\n"))
	write(t, root, "wip/threads/TH-0001-shape.md", threadFile("TH-0001", "Which shape?", "answered",
		"2026-10-19T10:00:00Z claude", "Left or right?",
		"2026-10-19T10:05:00Z alex", "Left. It is cheaper.",
		"2026-10-19T10:07:00Z claude", "Noted.",
		"2026-10-19T10:10:00Z alex", "Narrow, too."))
	write(t, root, "wip/threads/TH-0002-store.md", threadFile("TH-0002", "Which store?", "resolved",
		"2026-10-22T08:00:00Z claude", "Files or a database?",
		"2026-10-22T08:01:00Z orchestrator (recommendation)", "A database.",
		"2026-10-22T08:02:00Z alex", "Files: git keeps them.",
		"2026-10-22T08:03:00Z alex", "Resolved."))
	r, err := workitem.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	return r
}

// found is a candidate's kind, source, and time, as one string to compare.
func found(list []Candidate) []string {
	out := []string{}
	for _, c := range list {
		out = append(out, c.Kind+" "+c.Source+" "+c.At)
	}
	return out
}

// Harvest joins every thread's answers, resolved threads included, the
// narratives' decisions, and the issues' causes, sorted by when each was
// learned, each with its date, text, topics, and paths.
func TestHarvestListsEveryKindSortedByDate(t *testing.T) {
	got, err := Harvest(harvestRepo(t), time.Time{}, nil)
	if err != nil {
		t.Fatal(err)
	}
	want := []string{
		"answer TH-0001 2026-10-19T10:05:00Z",
		"answer TH-0001 2026-10-19T10:10:00Z",
		"decision S-0001 2026-10-20T12:00:00Z",
		"cause I-0001 2026-10-21T09:00:00Z",
		"recommendation TH-0002 2026-10-22T08:02:00Z",
	}
	if !reflect.DeepEqual(found(got), want) {
		t.Fatalf("candidates %v\nwant %v", found(got), want)
	}
	for _, c := range got {
		if c.Title == "" || c.Text == "" || c.Date != c.At[:10] || c.Topics == nil || c.Paths == nil || c.CitedBy == nil || len(c.CitedBy) != 0 {
			t.Errorf("candidate %+v", c)
		}
	}
	if c := got[4]; c.Text != "Files: git keeps them." || c.Overruled != "A database." || c.Question != "Files or a database?" {
		t.Errorf("the overruled recommendation: %+v", c)
	}
	if c := got[3]; c.Fix != "S-0001: numbers come from every branch." || c.Text != "Two writers took one number." {
		t.Errorf("the cause: %+v", c)
	}
}

// Harvest leaves out what was learned before since, and keeps what was
// learned at it.
func TestHarvestLeavesOutWhatIsOlderThanSince(t *testing.T) {
	got, err := Harvest(harvestRepo(t), time.Date(2026, 10, 20, 12, 0, 0, 0, time.UTC), nil)
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"decision S-0001 2026-10-20T12:00:00Z", "cause I-0001 2026-10-21T09:00:00Z", "recommendation TH-0002 2026-10-22T08:02:00Z"}
	if !reflect.DeepEqual(found(got), want) {
		t.Errorf("since 2026-10-20T12:00:00Z: %v\nwant %v", found(got), want)
	}
}

// Harvest lists the kinds it is given only, and refuses one that is none of
// Kinds, naming them.
func TestHarvestFiltersByKind(t *testing.T) {
	r := harvestRepo(t)
	got, err := Harvest(r, time.Time{}, []string{"cause", "answer"})
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"answer TH-0001 2026-10-19T10:05:00Z", "answer TH-0001 2026-10-19T10:10:00Z", "cause I-0001 2026-10-21T09:00:00Z"}
	if !reflect.DeepEqual(found(got), want) {
		t.Errorf("answers and causes: %v\nwant %v", found(got), want)
	}
	if got, err := Harvest(r, time.Time{}, []string{"recommendation"}); err != nil || !reflect.DeepEqual(found(got), []string{"recommendation TH-0002 2026-10-22T08:02:00Z"}) {
		t.Errorf("recommendations: %v, %v", found(got), err)
	}
	if _, err := Harvest(r, time.Time{}, []string{"answer", "rumour"}); err == nil || !strings.Contains(err.Error(), `"rumour"`) || !strings.Contains(err.Error(), "decision, answer, recommendation, cause") {
		t.Errorf("an unknown kind: %v", err)
	}
}

// Harvest names in CitedBy the entries whose sources name a candidate's
// source in any padding, in ID order, and none for a source no entry cites.
func TestHarvestNamesTheEntriesThatCiteASource(t *testing.T) {
	r := harvestRepo(t)
	write(t, r.Root, "design/digest/D-0002-two.md", strings.NewReplacer("D-0003", "D-0002", "sources: [S-0428, TH-0476]", "sources: [I-0001, S-0001]").Replace(good))
	write(t, r.Root, "design/digest/D-0001-one.md", strings.NewReplacer("D-0003", "D-0001", "sources: [S-0428, TH-0476]", "sources: [S-0001, th-1]").Replace(good))
	got, err := Harvest(r, time.Time{}, nil)
	if err != nil {
		t.Fatal(err)
	}
	var cited []string
	for _, c := range got {
		cited = append(cited, c.Source+" "+strings.Join(c.CitedBy, ","))
	}
	want := []string{"TH-0001 D-0001", "TH-0001 D-0001", "S-0001 D-0001,D-0002", "I-0001 D-0002", "TH-0002 "}
	if !reflect.DeepEqual(cited, want) {
		t.Errorf("cited by %q\nwant %q", cited, want)
	}
}

// A candidate's At is compared as a time whichever form it takes: a day
// alone is that day's midnight.
func TestLearnedReadsEitherFormOfAt(t *testing.T) {
	day, at := learned(Candidate{At: "2026-10-20"}), learned(Candidate{At: "2026-10-20T00:00:01Z"})
	if !day.Equal(time.Date(2026, 10, 20, 0, 0, 0, 0, time.UTC)) || !day.Before(at) {
		t.Errorf("day %v, time %v", day, at)
	}
	if !learned(Candidate{At: "soon"}).IsZero() {
		t.Error("an At in neither form is not the zero time")
	}
}
