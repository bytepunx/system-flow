package digest

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
	"time"
	"unicode/utf8"

	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// item is a work item file: its front matter, with a done transition at
// done when it is given, and a heading.
func item(id, typ, status, parent, done, extra string) string {
	transitions := "transitions: []\n"
	if done != "" {
		transitions = "transitions:\n  - to: in-progress\n    at: 2026-08-01T00:00:00Z\n    by: alex\n  - to: " + status + "\n    at: " + done + "\n    by: alex\n"
	}
	return fmt.Sprintf("---\nid: %s\ntype: %s\nnature: feature\ntitle: %s title\nstatus: %s\nparent: %s\nowner: alex\ncreated: 2026-08-01T00:00:00Z\nupdated: 2026-08-01T00:00:00Z\n%s%s---\n\n# %s %s title\n", id, typ, id, status, parent, transitions, extra, id, id)
}

// narrative is a story's narrative with the sections given after its
// Context.
func narrative(id, sections string) string {
	return "---\nstream: " + id + "\ntitle: t\nupdated: 2026-08-01T00:00:00Z\nagent: a\nsession: s\n---\n\n# " + id + " t\n\n## Context\n\nWhy.\n\n" + sections + "\n## Log\n"
}

// issue is an issue file with the status and body given.
func issue(id, status, body string) string {
	return "---\nid: " + id + "\ntitle: " + id + " went wrong\nclass: defect\nstatus: " + status + "\ncount: 1\nfirst_reported: 2026-08-01T00:00:00Z\nlast_reported: 2026-08-01T00:00:00Z\nupdated: 2026-08-02T00:00:00Z\n---\n\n# " + id + " went wrong\n" + body
}

// decisionsS0002 holds a bold lead-in, a bullet continued over indented
// lines and a blank line, an abbreviation, a long bullet, a paragraph that
// is no bullet, and the threads mirror.
var decisionsS0002 = "## Decisions\n\n" +
	"- **The stop test (T-0630).** The cause is a race in the test, not in `Stop. Here`. More text.\n" +
	"- Plan layers: 1. T-1807 first. 2. T-1808 second.\n" +
	"  It continues here\n" +
	"  - a nested point.\n" +
	"\n" +
	"  And after a blank line.\n" +
	"- Use a helper, e.g. The Helper, here. Then more.\n" +
	"- " + longWords + "\n" +
	"\n" +
	"A paragraph that is not a bullet.\n" +
	"\n" +
	"<!-- threads:start -->\n" +
	"- TH-0001 a mirrored question\n" +
	"<!-- threads:end -->\n"

// longWords is a bullet of 303 characters with no sentence end.
var longWords = strings.Repeat("word ", 60) + "end"

// sourcesRepo is a project with a sub-project flai, an epic, five stories
// and their narratives, a task, and six issues:
//
//   - S-0001 archived, done 2026-09-01, one decision;
//   - S-0002 archived, done 2026-10-20, under E-0001, tagged cli, touching
//     the digest, with decisionsS0002 and an archived task T-0010;
//   - S-0003 done 2026-10-22 and not archived, one decision;
//   - S-0004 archived, done, a narrative with no Decisions;
//   - S-0005 in progress, a decision that counts for nothing yet.
func sourcesRepo(t *testing.T) *workitem.Repo {
	t.Helper()
	root := t.TempDir()
	write(t, root, "system-flow.yaml", "version: 1\nname: t\nkey: t\nlayout:\n  design: design\n  docs: docs\n  wip: wip\nprojects:\n  - name: flai\n    path: flai\n    kind: go\n    tags: [cli]\n")
	write(t, root, "wip/kanban/epics/E-0001-e.md", item("E-0001", "epic", "in-progress", "", "", "topics: [planning]\n"))
	write(t, root, "wip/archive/kanban/stories/S-0001-s.md", item("S-0001", "story", "done", "", "2026-09-01T10:00:00Z", ""))
	write(t, root, "wip/archive/agents/S-0001.md", narrative("S-0001", "## Decisions\n\n- An old decision.\n"))
	write(t, root, "wip/archive/kanban/stories/S-0002-s.md", item("S-0002", "story", "done", "E-0001", "2026-10-20T12:00:00Z", "tags: [cli]\ntopics: [digest]\ntouches: [flai/internal/digest/sources.go]\n"))
	write(t, root, "wip/archive/kanban/tasks/T-0010-t.md", item("T-0010", "task", "done", "S-0002", "2026-10-20T11:00:00Z", ""))
	write(t, root, "wip/archive/agents/S-0002.md", narrative("S-0002", decisionsS0002+"\n## Open questions\n\n- Not a decision.\n"))
	write(t, root, "wip/kanban/stories/S-0003-s.md", item("S-0003", "story", "done", "", "2026-10-22T08:00:00Z", ""))
	write(t, root, "wip/agents/S-0003.md", narrative("S-0003", "## Decisions\n\n- Not archived yet. Still learned.\n"))
	write(t, root, "wip/archive/kanban/stories/S-0004-s.md", item("S-0004", "story", "done", "", "2026-10-21T08:00:00Z", ""))
	write(t, root, "wip/archive/agents/S-0004.md", narrative("S-0004", "## Current state\n\nDone.\n"))
	write(t, root, "wip/kanban/stories/S-0005-s.md", item("S-0005", "story", "in-progress", "", "", ""))
	write(t, root, "wip/agents/S-0005.md", narrative("S-0005", "## Decisions\n\n- Still in progress.\n"))
	write(t, root, "wip/agents/index.md", "# Narratives\n")

	write(t, root, "design/issues/I-0001-a.md", issue("I-0001", "closed", "\n## Description\nTwo writers took\none number.\n\n## Instances\n\n### 2026-08-01T00:00:00Z\nStory: S-0001.\nFirst.\n\n## Remediation\nNumber from every branch.\nClosed 2026-10-10T00:00:00Z: a first fix\nClosed 2026-10-21T09:00:00Z: S-0002: numbers come from every branch.\n"))
	write(t, root, "design/issues/I-0002-b.md", issue("I-0002", "closed", "\n## Description\nThe push was refused.\n\n## Remediation\nRetry the push.\n\nStory S-0003 remediates this issue, created from it at 2026-10-20T00:00:00Z.\nClosed 2026-10-21T10:00:00Z:\n"))
	write(t, root, "design/issues/I-0003-c.md", issue("I-0003", "closed", "\n## Remediation\nClosed 2026-10-21T11:00:00Z: S-0002: fixed with no description.\n"))
	write(t, root, "design/issues/I-0004-d.md", issue("I-0004", "open", "\n## Description\nStill happens.\n\n## Remediation\nNone yet.\n"))
	write(t, root, "design/issues/I-0005-e.md", issue("I-0005", "closed", "\n## Description\nNo fix recorded.\n\n## Remediation\nClosed 2026-10-21T12:00:00Z:\n"))
	write(t, root, "design/issues/I-0006-f.md", issue("I-0006", "closed", "\n## Description\nAn old cause.\n\n## Remediation\nClosed 2026-09-01T00:00:00Z: no longer applies\n"))
	write(t, root, "design/issues/summary.md", "# Issues\n")
	r, err := workitem.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	return r
}

// The topics OfStory gives S-0002: its own, its epic's, what its tag and
// touches reach in flai, code, and all.
var s0002Topics = []string{"digest", "planning", "flai", "cli", "go", "code", "all"}

// ADR-0191: Lookup gives a story's topics and touches, archived or not; a
// task its story's; an epic its own worked out the same way; and an ID it
// cannot find, empty lists.
func TestLookupGivesAnItemsTopicsAndPaths(t *testing.T) {
	lookup := Lookup(sourcesRepo(t))
	for _, c := range []struct {
		id            string
		topics, paths []string
	}{
		{"S-0002", s0002Topics, []string{"flai/internal/digest/sources.go"}},
		{"s-2", s0002Topics, []string{"flai/internal/digest/sources.go"}},
		{"T-0010", s0002Topics, []string{"flai/internal/digest/sources.go"}},
		{"S-0003", []string{"all"}, []string{}},
		{"E-0001", []string{"planning", "all"}, []string{}},
		{"S-0999", []string{}, []string{}},
		{"TH-0001", []string{}, []string{}},
	} {
		topics, paths := lookup(c.id)
		if topics == nil || paths == nil || !reflect.DeepEqual(topics, c.topics) || !reflect.DeepEqual(paths, c.paths) {
			t.Errorf("Lookup(%s) = %#v, %#v; want %#v, %#v", c.id, topics, paths, c.topics, c.paths)
		}
	}
	topics, _ := lookup("S-0002")
	topics[0] = "changed"
	if again, _ := lookup("S-0002"); again[0] != "digest" {
		t.Errorf("a caller's change reached the cache: %v", again)
	}
}

// Decisions lists each top-level bullet of a done story's Decisions, archived
// or not, in story order, and nothing from a narrative without the section
// or a story still open.
func TestDecisionsListsEachBulletOfADoneStory(t *testing.T) {
	got, err := Decisions(sourcesRepo(t), time.Time{}, nil)
	if err != nil {
		t.Fatal(err)
	}
	var sources []string
	for _, c := range got {
		sources = append(sources, c.Source)
	}
	if want := []string{"S-0001", "S-0002", "S-0002", "S-0002", "S-0002", "S-0003"}; !reflect.DeepEqual(sources, want) {
		t.Fatalf("sources %v, want %v", sources, want)
	}
	first := got[1]
	want := Candidate{
		Kind: "decision", Source: "S-0002", Ref: "wip/archive/agents/S-0002.md",
		Title:  "The stop test (T-0630): The cause is a race in the test, not in `Stop. Here`.",
		Topics: s0002Topics, Paths: []string{"flai/internal/digest/sources.go"},
		At: "2026-10-20T12:00:00Z", Date: "2026-10-20",
		Text:    "**The stop test (T-0630).** The cause is a race in the test, not in `Stop. Here`. More text.",
		CitedBy: []string{},
	}
	if !reflect.DeepEqual(first, want) {
		t.Errorf("first decision of S-0002:\n%#v\nwant:\n%#v", first, want)
	}
	if c := got[2]; c.Text != "Plan layers: 1. T-1807 first. 2. T-1808 second. It continues here - a nested point. And after a blank line." || c.Title != "Plan layers: 1. T-1807 first. 2. T-1808 second." {
		t.Errorf("a bullet over indented lines: title %q, text %q", c.Title, c.Text)
	}
	if c := got[3]; c.Title != "Use a helper, e.g. The Helper, here." {
		t.Errorf("an abbreviation ended the sentence: %q", c.Title)
	}
	if c := got[4]; c.Text != longWords || !strings.HasSuffix(c.Title, "word…") || utf8.RuneCountInString(c.Title) > titleLimit+1 {
		t.Errorf("a long bullet: title %q (%d characters)", c.Title, utf8.RuneCountInString(c.Title))
	}
	if c := got[5]; c.Ref != "wip/agents/S-0003.md" || c.Title != "Not archived yet." || c.Date != "2026-10-22" {
		t.Errorf("the unarchived narrative's decision: %+v", c)
	}
	for _, c := range got {
		if strings.Contains(c.Text, "TH-0001") || strings.Contains(c.Text, "paragraph") || strings.Contains(c.Text, "Not a decision") {
			t.Errorf("text from outside a Decisions bullet: %q", c.Text)
		}
	}
}

// Decisions leaves out a story completed before since.
func TestDecisionsLeavesOutAStoryCompletedBeforeSince(t *testing.T) {
	got, err := Decisions(sourcesRepo(t), time.Date(2026, 10, 1, 0, 0, 0, 0, time.UTC), nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 5 || got[0].Source != "S-0002" {
		t.Errorf("since 2026-10-01: %d decisions, first from %s; want 5, first from S-0002", len(got), got[0].Source)
	}
}

// Causes lists each closed issue with a cause and a fix: the fix is the last
// close's reason, or the Remediation when that reason is empty; the closing
// story the fix names, or the one made from the issue, gives its topics. An
// issue with no description or no fix, or still open, gives nothing.
func TestCausesListsEachClosedIssueWithACauseAndAFix(t *testing.T) {
	got, err := Causes(sourcesRepo(t), time.Time{}, nil)
	if err != nil {
		t.Fatal(err)
	}
	var sources []string
	for _, c := range got {
		sources = append(sources, c.Source)
	}
	if want := []string{"I-0001", "I-0002", "I-0006"}; !reflect.DeepEqual(sources, want) {
		t.Fatalf("sources %v, want %v", sources, want)
	}
	want := Candidate{
		Kind: "cause", Source: "I-0001", Title: "I-0001 went wrong",
		Topics: s0002Topics, Paths: []string{"flai/internal/digest/sources.go"},
		At: "2026-10-21T09:00:00Z", Date: "2026-10-21",
		Text: "Two writers took\none number.", Fix: "S-0002: numbers come from every branch.",
		CitedBy: []string{},
	}
	if !reflect.DeepEqual(got[0], want) {
		t.Errorf("I-0001:\n%#v\nwant:\n%#v", got[0], want)
	}
	if c := got[1]; c.Fix != "Retry the push." || !reflect.DeepEqual(c.Topics, []string{"all"}) || c.At != "2026-10-21T10:00:00Z" {
		t.Errorf("I-0002, fixed by its Remediation and made into S-0003: %+v", c)
	}
	if c := got[2]; c.Topics == nil || len(c.Topics) != 0 || c.Paths == nil || len(c.Paths) != 0 || c.Fix != "no longer applies" {
		t.Errorf("I-0006, closed by no story: %+v", c)
	}

	got, err = Causes(sourcesRepo(t), time.Date(2026, 10, 1, 0, 0, 0, 0, time.UTC), nil)
	if err != nil || len(got) != 2 || got[1].Source != "I-0002" {
		t.Errorf("since 2026-10-01: %+v, %v; want I-0001 and I-0002", got, err)
	}
}
