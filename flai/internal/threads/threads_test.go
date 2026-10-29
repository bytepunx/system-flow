package threads

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/mdlint"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

var t0 = time.Date(2026, 9, 17, 21, 0, 0, 0, time.UTC)

func project(t *testing.T) *workitem.Repo {
	t.Helper()
	root := t.TempDir()
	_ = os.WriteFile(filepath.Join(root, "system-flow.yaml"), []byte("version: 1\nname: t\nkey: t\nlayout:\n  design: design\n  docs: docs\n  wip: wip\n"), 0o644)
	for _, d := range []string{"design/system", "wip/kanban/epics", "wip/kanban/stories", "wip/kanban/tasks", "wip/agents", "wip/archive"} {
		_ = os.MkdirAll(filepath.Join(root, d), 0o755)
	}
	_ = os.WriteFile(filepath.Join(root, "design/system/plan.md"), []byte("---\ntitle: Plan\n---\n\n# Plan\n\n## Shape\ntext\n\n## Risks\nmore\n"), 0o644)
	r, err := workitem.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	return r
}

func TestThreadLifecycle(t *testing.T) {
	r := project(t)
	th, err := New(r, NewOptions{Title: "Is the shape right?", On: "design/system/plan.md", Heading: "shape", Author: "alex", Text: "I wonder.", Now: t0})
	if err != nil {
		t.Fatal(err)
	}
	if th.ID != "TH-0001" || th.Status != "open" || th.Anchor.Path != "design/system/plan.md" || th.Anchor.Heading != "shape" || th.Opener() != "alex" {
		t.Fatalf("new: %+v", th)
	}
	if _, err := os.Stat(filepath.Join(r.WipDir(), "threads", "TH-0001-is-the-shape-right.md")); err != nil {
		t.Fatal("file name")
	}
	if _, err := New(r, NewOptions{Title: "x", On: "design/system/plan.md", Heading: "nope", Author: "alex", Now: t0}); err == nil || !strings.Contains(err.Error(), "heading") {
		t.Errorf("missing heading should fail: %v", err)
	}
	if _, err := New(r, NewOptions{Title: "x", On: "docs/missing.md", Author: "alex", Now: t0}); err == nil {
		t.Error("missing path should fail")
	}
	if _, err := New(r, NewOptions{Title: "x", On: "S-0099", Author: "alex", Now: t0}); err == nil {
		t.Error("missing item should fail")
	}
	// reply from an agent answers; the opener's follow-up reopens; resolve closes
	th2, err := Reply(r, "th-1", "claude", "Yes, because.", t0.Add(time.Minute))
	if err != nil || th2.Status != "answered" || len(th2.Entries()) != 2 || th2.Entries()[1].Text != "Yes, because." || !contains(th2.Participants, "claude") {
		t.Fatalf("reply: %+v %v", th2, err)
	}
	th3, _ := Reply(r, "TH-0001", "alex", "Then also consider risks.", t0.Add(2*time.Minute))
	if th3.Status != "open" {
		t.Errorf("opener follow-up should reopen: %s", th3.Status)
	}
	th4, err := Resolve(r, "TH-0001", "alex", "settled in ADR-0021", t0.Add(3*time.Minute))
	if err != nil || th4.Status != "resolved" || !strings.Contains(th4.Body, "Resolved: settled in ADR-0021") || th4.Open() {
		t.Fatalf("resolve: %+v %v", th4, err)
	}
	// round trip through the file keeps everything
	back, err := Read(th4.Path)
	if err != nil || back.Marshal() != th4.Marshal() || len(back.Entries()) != 4 {
		t.Fatalf("round trip: %v\n%s", err, back.Marshal())
	}
	if err := back.Validate(); err != nil {
		t.Error(err)
	}
	// reply reopens a resolved thread
	th5, _ := Reply(r, "TH-0001", "claude", "One more thing.", t0.Add(4*time.Minute))
	if th5.Status != "answered" {
		t.Errorf("reply after resolve: %s", th5.Status)
	}
}

func TestItemAnchorsAndNarrativeMirror(t *testing.T) {
	r := project(t)
	epic, _ := r.Create(workitem.NewOptions{Type: workitem.Epic, Title: "E", Owner: "a", Now: t0})
	story, _ := r.Create(workitem.NewOptions{Type: workitem.Story, Title: "S", Parent: epic.ID, Owner: "a", Now: t0})
	task, _ := r.Create(workitem.NewOptions{Type: workitem.Task, Title: "T", Parent: story.ID, Owner: "a", Now: t0})
	if _, err := r.OpenStream(story, workitem.StreamOptions{Agent: "claude", Now: t0}); err != nil {
		t.Fatal(err)
	}
	th, err := New(r, NewOptions{Title: "Task question", On: task.ID, Author: "alex", Text: "?", Now: t0})
	if err != nil {
		t.Fatal(err)
	}
	if th.Anchor.Item != task.ID || !strings.HasPrefix(th.Anchor.Path, "wip/kanban/tasks/T-0001") {
		t.Fatalf("item anchor: %+v", th.Anchor)
	}
	if StoryOf(r, th) != story.ID {
		t.Errorf("story of task thread: %s", StoryOf(r, th))
	}
	list, _ := For(r, "t-1")
	if len(list) != 1 {
		t.Errorf("For by item in any padding: %d", len(list))
	}
	if err := MirrorNarrative(r, story.ID); err != nil {
		t.Fatal(err)
	}
	n, _ := os.ReadFile(r.NarrativePath(story.ID))
	if !strings.Contains(string(n), "## Open questions\n<!-- threads:start -->") || !strings.Contains(string(n), "- TH-0001 (open, alex, 2026-09-17): Task question") {
		t.Fatalf("mirror:\n%s", n)
	}
	if _, err := Resolve(r, th.ID, "claude", "", t0.Add(time.Minute)); err != nil {
		t.Fatal(err)
	}
	_ = MirrorNarrative(r, story.ID)
	n, _ = os.ReadFile(r.NarrativePath(story.ID))
	if strings.Contains(string(n), "threads:start") || strings.Contains(string(n), "TH-0001") {
		t.Errorf("resolved threads leave the mirror:\n%s", n)
	}
	if strings.Count(string(n), "## Open questions") != 1 {
		t.Error("heading must survive")
	}
}

// I-0231: a long title made a file name that Claude Code's Edit could not
// write once its temporary suffix was added; the slug is cut at 80 bytes.
func TestNewCutsALongTitlesFileName(t *testing.T) {
	r := project(t)
	title := "Claude Code's Edit tool cannot write a work item whose file name is near the 255-byte limit, since its temporary file name adds about 25 bytes"
	th, err := New(r, NewOptions{Title: title, On: "design/system/plan.md", Author: "alex", Text: "Why?", Now: t0})
	if err != nil {
		t.Fatal(err)
	}
	base := filepath.Base(th.Path)
	if max := len(th.ID) + 1 + 80 + len(".md"); len(base) > max || !strings.HasPrefix(base, th.ID+"-claude-code-s-edit-") || !strings.HasSuffix(base, ".md") {
		t.Errorf("the file name is %d bytes, want at most %d on a word: %s", len(base), max, base)
	}
	if _, err := os.Stat(th.Path); err != nil {
		t.Fatal(err)
	}
	if back, err := Read(th.Path); err != nil || back.Title != title {
		t.Errorf("the thread keeps its whole title: %v", err)
	}
}

// I-0043: a reply and a resolution in one second by one author share the
// entry's heading; another author in that second gets a heading of its own.
func TestSameSecondEntriesShareAHeading(t *testing.T) {
	r := project(t)
	if _, err := New(r, NewOptions{Title: "Which way", On: "design/system/plan.md", Author: "alex", Text: "Left or right?", Now: t0}); err != nil {
		t.Fatal(err)
	}
	at := t0.Add(time.Minute)
	if _, err := Reply(r, "TH-0001", "claude", "Left.", at); err != nil {
		t.Fatal(err)
	}
	th, err := Resolve(r, "TH-0001", "claude", "answered", at.Add(500*time.Millisecond))
	if err != nil {
		t.Fatal(err)
	}
	heading := "### 2026-09-17T21:01:00Z claude"
	if n := strings.Count(th.Body, heading); n != 1 {
		t.Fatalf("want one %q heading, got %d:\n%s", heading, n, th.Body)
	}
	e := th.Entries()
	if len(e) != 2 || e[1].Author != "claude" || e[1].Text != "Left.\n\nResolved: answered" {
		t.Fatalf("entries: %+v", e)
	}
	if th.Status != "resolved" {
		t.Errorf("status %s", th.Status)
	}
	th, err = Reply(r, "TH-0001", "alex", "Thanks.", at)
	if err != nil {
		t.Fatal(err)
	}
	if e := th.Entries(); len(e) != 3 || e[2].Author != "alex" || e[2].At != "2026-09-17T21:01:00Z" {
		t.Fatalf("another author in the same second has its own entry: %+v", e)
	}
	got, err := Read(th.Path)
	if err != nil {
		t.Fatal(err)
	}
	if got.Body != th.Body {
		t.Error("saved body differs from the returned one")
	}
}

// S-0179: a thread's title loses the punctuation its heading cannot end
// with, and an entry the project's lint rejects is refused, leaving the
// thread as it was.
func TestThreadWritesTheLintRejectsAreRefused(t *testing.T) {
	r := project(t)
	_ = os.WriteFile(filepath.Join(r.Root, ".markdownlint.yaml"), []byte("default: true\nMD013: false\nMD022:\n  lines_below: 0\nMD024:\n  siblings_only: true\nMD025:\n  front_matter_title: \"\"\nMD032: false\nMD041: false\n"), 0o644)
	th, err := New(r, NewOptions{Title: "Settled.", On: "design/system/plan.md", Author: "alex", Text: "Is it?", Now: t0})
	if err != nil {
		t.Fatal(err)
	}
	if th.Title != "Settled" || !strings.Contains(th.Body, "# TH-0001 Settled\n") {
		t.Errorf("title: %q\n%s", th.Title, th.Body)
	}
	if _, err := New(r, NewOptions{Title: "Bad", On: "design/system/plan.md", Author: "alex", Text: "**Bold as a heading**", Now: t0}); err == nil || !strings.Contains(err.Error(), "MD036") {
		t.Errorf("a first entry the lint rejects is refused: %v", err)
	}
	if m, _ := filepath.Glob(filepath.Join(Dir(r), "TH-0002-*")); len(m) != 0 {
		t.Errorf("nothing is written: %v", m)
	}
	was, _ := os.ReadFile(th.Path)
	_, err = Reply(r, th.ID, "claude", "Yes.\n\n\n\nIt is.", t0.Add(time.Minute))
	var le *mdlint.Error
	if !errors.As(err, &le) || le.Findings[0].Rule != "MD012" {
		t.Fatalf("a reply the lint rejects is refused with the rule: %v", err)
	}
	if got, _ := os.ReadFile(th.Path); string(got) != string(was) {
		t.Error("a refused reply leaves the thread as it was")
	}
	if _, err := Resolve(r, th.ID, "claude", "settled", t0.Add(time.Minute)); err != nil {
		t.Errorf("a clean resolution is kept: %v", err)
	}
}

// The mirror block, removed when its last thread resolves, leaves one blank
// line where it stood, even when a blank line was put above it by hand
// (MD012 on S-0175's and S-0180's narratives).
func TestMirrorRemovalLeavesOneBlankLine(t *testing.T) {
	r := project(t)
	story, _ := r.Create(workitem.NewOptions{Type: workitem.Story, Title: "S", Owner: "a", Now: t0})
	if _, err := r.OpenStream(story, workitem.StreamOptions{Agent: "claude", Now: t0}); err != nil {
		t.Fatal(err)
	}
	th, err := New(r, NewOptions{Title: "Q", On: story.ID, Author: "claude", Text: "?", Now: t0})
	if err != nil {
		t.Fatal(err)
	}
	_ = MirrorNarrative(r, story.ID)
	path := r.NarrativePath(story.ID)
	n, _ := os.ReadFile(path)
	// the agent rewrites the section with a blank line under the heading
	edited := strings.Replace(string(n), "## Open questions\n<!-- threads:start -->", "## Open questions\n\n<!-- threads:start -->", 1)
	if edited == string(n) {
		t.Fatalf("fixture:\n%s", n)
	}
	_ = os.WriteFile(path, []byte(edited), 0o644)
	if _, err := Resolve(r, th.ID, "alex", "", t0.Add(time.Minute)); err != nil {
		t.Fatal(err)
	}
	_ = MirrorNarrative(r, story.ID)
	n, _ = os.ReadFile(path)
	if strings.Contains(string(n), "\n\n\n") || !strings.Contains(string(n), "## Open questions\n\n## Log") {
		t.Errorf("one blank line between the sections:\n%s", n)
	}
}

// S-0181: a thread a newer flai wrote is listed, and a reply keeps the
// field this flai does not know.
func TestAThreadWithAnUnknownFieldIsListedAndKept(t *testing.T) {
	r := project(t)
	th, err := New(r, NewOptions{Title: "Shape?", On: "design/system/plan.md", Author: "alex", Text: "I wonder.", Now: t0})
	if err != nil {
		t.Fatal(err)
	}
	data, _ := os.ReadFile(th.Path)
	doc := strings.Replace(string(data), "\n---\n", "\npriority: high\n---\n", 1)
	if err := os.WriteFile(th.Path, []byte(doc), 0o644); err != nil {
		t.Fatal(err)
	}
	list, err := List(r)
	if err != nil || len(list) != 1 || len(list[0].Unknown) != 1 || list[0].Unknown[0].Name != "priority" {
		t.Fatalf("listing: %v %+v", err, list)
	}
	if _, err := Reply(r, th.ID, "claude", "Yes.", t0.Add(time.Minute)); err != nil {
		t.Fatal(err)
	}
	after, _ := os.ReadFile(th.Path)
	if !strings.Contains(string(after), "priority: high\n---\n") || !strings.Contains(string(after), "Yes.") {
		t.Errorf("a reply dropped the unknown field:\n%s", after)
	}
}

// wipLint writes the repository's markdown lint configuration into a test
// project, so that a refused write fails the test.
func wipLint(t *testing.T, r *workitem.Repo) {
	t.Helper()
	cfg := "default: true\nMD013: false\nMD022:\n  lines_below: 0\nMD024:\n  siblings_only: true\nMD025:\n  front_matter_title: \"\"\nMD032: false\nMD033: false\nMD041: false\nMD060: false\n"
	if err := os.WriteFile(filepath.Join(r.Root, ".markdownlint.yaml"), []byte(cfg), 0o644); err != nil {
		t.Fatal(err)
	}
}

// lintClean fails the test when the thread's file has a lint finding.
func lintClean(t *testing.T, r *workitem.Repo, th *Thread) {
	t.Helper()
	c, err := mdlint.Load(r.Root)
	if err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(th.Path)
	if err != nil {
		t.Fatal(err)
	}
	if f := c.Lint(string(data)); len(f) > 0 {
		t.Fatalf("lint findings %+v in:\n%s", f, data)
	}
}

// S-0220, ADR-0090: a recommendation leaves the status as it was, reads back
// with its marks, and lints clean; an answer citing a source answers.
func TestARecommendationLeavesTheStatusAndReadsBack(t *testing.T) {
	r := project(t)
	wipLint(t, r)
	if _, err := New(r, NewOptions{Title: "Which shape", On: "design/system/plan.md", Author: "claude", Text: "Left or right?", Now: t0}); err != nil {
		t.Fatal(err)
	}
	at := t0.Add(time.Minute)
	th, err := ReplyWith(r, "TH-0001", "orchestrator", "Left, as the plan says.", at, Marks{Recommendation: true, Source: Source{Path: "design/system/plan.md", Heading: "Shape"}})
	if err != nil {
		t.Fatal(err)
	}
	if th.Status != "open" {
		t.Errorf("a recommendation leaves the status open: %s", th.Status)
	}
	if !strings.Contains(th.Body, "\n### 2026-09-17T21:01:00Z orchestrator (recommendation)\nLeft, as the plan says.\n\nSource: design/system/plan.md § Shape\n") {
		t.Errorf("on disk:\n%s", th.Body)
	}
	lintClean(t, r, th)
	back, err := Read(th.Path)
	if err != nil {
		t.Fatal(err)
	}
	e := back.Entries()
	want := Source{Path: "design/system/plan.md", Heading: "Shape"}
	if len(e) != 2 || !e[1].Recommendation || e[1].Author != "orchestrator" || e[1].Source == nil || *e[1].Source != want || e[0].Recommendation || e[0].Source != nil {
		t.Fatalf("entries: %+v", e)
	}
	if !contains(back.Participants, "orchestrator") || back.Opener() != "claude" {
		t.Errorf("participants %v, opener %s", back.Participants, back.Opener())
	}
	p, ok := View(r, back)["pending_recommendation"].(*Entry)
	if !ok || p == nil || p.At != "2026-09-17T21:01:00Z" || p.Author != "orchestrator" {
		t.Fatalf("pending recommendation: %+v", View(r, back)["pending_recommendation"])
	}
	// the opener's follow-up leaves the recommendation pending
	th, err = Reply(r, "TH-0001", "claude", "Also, why?", at.Add(time.Minute))
	if err != nil || th.Status != "open" || th.PendingRecommendation() == nil {
		t.Fatalf("follow-up: %v %s %+v", err, th.Status, th.PendingRecommendation())
	}

	// a second recommendation on an answered thread keeps it answered
	if _, err := Reply(r, "TH-0001", "alex", "Not sure.", at.Add(2*time.Minute)); err != nil {
		t.Fatal(err)
	}
	th, err = ReplyWith(r, "TH-0001", "orchestrator", "Left.", at.Add(3*time.Minute), Marks{Recommendation: true})
	if err != nil || th.Status != "answered" {
		t.Fatalf("recommendation on an answered thread: %v %s", err, th.Status)
	}
	if p := th.PendingRecommendation(); p == nil || p.Source != nil {
		t.Errorf("a recommendation without a source is pending: %+v", p)
	}
	lintClean(t, r, th)

	// refused, writing nothing: the opener's own, a missing source, a resolved thread
	was, _ := os.ReadFile(th.Path)
	if _, err := ReplyWith(r, "TH-0001", "claude", "Right.", at.Add(4*time.Minute), Marks{Recommendation: true}); err == nil || !strings.Contains(err.Error(), "opened") {
		t.Errorf("the opener's recommendation is refused: %v", err)
	}
	if _, err := ReplyWith(r, "TH-0001", "orchestrator", "Right.", at.Add(4*time.Minute), Marks{Source: Source{Path: "design/adrs/0099-none.md"}}); err == nil || !strings.Contains(err.Error(), "source") {
		t.Errorf("a missing source is refused: %v", err)
	}
	if _, err := ReplyWith(r, "TH-0001", "orchestrator", "Right.", at.Add(4*time.Minute), Marks{Source: Source{Path: "design/system/plan.md", Heading: "Nope"}}); err == nil || !strings.Contains(err.Error(), "heading") {
		t.Errorf("a missing source heading is refused: %v", err)
	}
	if got, _ := os.ReadFile(th.Path); string(got) != string(was) {
		t.Error("a refused reply leaves the thread as it was")
	}
	if _, err := Resolve(r, "TH-0001", "claude", "", at.Add(5*time.Minute)); err != nil {
		t.Fatal(err)
	}
	if _, err := ReplyWith(r, "TH-0001", "orchestrator", "Left.", at.Add(6*time.Minute), Marks{Recommendation: true}); err == nil || !strings.Contains(err.Error(), "resolved") {
		t.Errorf("a recommendation on a resolved thread is refused: %v", err)
	}

	// an answer that cites a source answers, and reads back its source
	if _, err := New(r, NewOptions{Title: "Which risk", On: "design/system/plan.md", Author: "claude", Text: "Which?", Now: t0}); err != nil {
		t.Fatal(err)
	}
	th, err = ReplyWith(r, "TH-0002", "orchestrator", "The one in the plan.", at, Marks{Source: Source{Path: "design/system/plan.md"}})
	if err != nil || th.Status != "answered" {
		t.Fatalf("an answer with a source: %v %+v", err, th)
	}
	e = th.Entries()
	if e[1].Recommendation || e[1].Source == nil || *e[1].Source != (Source{Path: "design/system/plan.md"}) || th.PendingRecommendation() != nil {
		t.Errorf("answer entry: %+v", e[1])
	}
	lintClean(t, r, th)
}

// S-0220, ADR-0090: the operator confirms a pending recommendation with one
// call, which makes it the answer; with none pending, Confirm is refused.
func TestConfirmMakesTheRecommendationTheAnswer(t *testing.T) {
	r := project(t)
	wipLint(t, r)
	if _, err := New(r, NewOptions{Title: "Which shape", On: "design/system/plan.md", Author: "claude", Text: "Left or right?", Now: t0}); err != nil {
		t.Fatal(err)
	}
	if _, err := Confirm(r, "TH-0001", "alex", t0.Add(time.Minute)); err == nil || !strings.Contains(err.Error(), "no recommendation") {
		t.Errorf("confirming with none pending is refused: %v", err)
	}
	rec := t0.Add(2 * time.Minute)
	if _, err := ReplyWith(r, "TH-0001", "orchestrator", "Left.", rec, Marks{Recommendation: true, Source: Source{Path: "design/system/plan.md", Heading: "Shape"}}); err != nil {
		t.Fatal(err)
	}
	if _, err := Confirm(r, "TH-0001", "orchestrator", t0.Add(3*time.Minute)); err == nil || !strings.Contains(err.Error(), "cannot confirm") {
		t.Errorf("the recommendation's author cannot confirm it: %v", err)
	}
	th, err := Confirm(r, "th-1", "alex", t0.Add(3*time.Minute))
	if err != nil {
		t.Fatal(err)
	}
	if th.Status != "answered" || th.PendingRecommendation() != nil || !contains(th.Participants, "alex") {
		t.Fatalf("confirmed: %s %+v %v", th.Status, th.PendingRecommendation(), th.Participants)
	}
	e := th.Entries()
	last := e[len(e)-1]
	if len(e) != 3 || last.Author != "alex" || last.Recommendation || last.Text != "Confirmed the recommendation of 2026-09-17T21:02:00Z orchestrator.\n\nSource: design/system/plan.md § Shape" || last.Source == nil {
		t.Fatalf("confirming entry: %+v", e)
	}
	lintClean(t, r, th)
	if _, err := Confirm(r, "TH-0001", "alex", t0.Add(4*time.Minute)); err == nil {
		t.Error("a confirmed recommendation is not confirmed twice")
	}

	// an answer after a recommendation takes its place: nothing is pending
	if _, err := ReplyWith(r, "TH-0001", "orchestrator", "Right, then.", t0.Add(5*time.Minute), Marks{Recommendation: true}); err != nil {
		t.Fatal(err)
	}
	th, err = Reply(r, "TH-0001", "alex", "No: up.", t0.Add(6*time.Minute))
	if err != nil || th.PendingRecommendation() != nil {
		t.Fatalf("an answer overrides the recommendation: %v %+v", err, th.PendingRecommendation())
	}
	if _, err := Confirm(r, "TH-0001", "alex", t0.Add(7*time.Minute)); err == nil {
		t.Error("confirming an overridden recommendation is refused")
	}

	// the opener may confirm a recommendation on their own thread
	if _, err := New(r, NewOptions{Title: "Which risk", On: "design/system/plan.md", Author: "alex", Text: "Which?", Now: t0}); err != nil {
		t.Fatal(err)
	}
	if _, err := ReplyWith(r, "TH-0002", "orchestrator", "The first.", rec, Marks{Recommendation: true}); err != nil {
		t.Fatal(err)
	}
	th, err = Confirm(r, "TH-0002", "alex", t0.Add(3*time.Minute))
	if err != nil || th.Status != "answered" || th.PendingRecommendation() != nil {
		t.Fatalf("the opener confirms: %v %s %+v", err, th.Status, th.PendingRecommendation())
	}
}

// ADR-0090: a thread an older flai wrote reads as before, with no marks and
// nothing pending; an older flai reads a recommendation's mark as part of
// its author.
func TestAnOlderThreadReadsAsBefore(t *testing.T) {
	r := project(t)
	doc := "---\nid: TH-0001\ntitle: Which shape\nanchor:\n  path: design/system/plan.md\nstatus: answered\nparticipants: [claude, alex]\ncreated: 2026-09-17T21:00:00Z\nupdated: 2026-09-17T21:01:00Z\n---\n\n# TH-0001 Which shape\n\nOn design/system/plan.md.\n\n## Entries\n\n### 2026-09-17T21:00:00Z claude\nLeft or right?\n\n### 2026-09-17T21:01:00Z alex\nLeft.\n\nSee the plan.\n"
	_ = os.MkdirAll(Dir(r), 0o755)
	if err := os.WriteFile(filepath.Join(Dir(r), "TH-0001-which-shape.md"), []byte(doc), 0o644); err != nil {
		t.Fatal(err)
	}
	th, err := Get(r, "TH-0001")
	if err != nil {
		t.Fatal(err)
	}
	e := th.Entries()
	if len(e) != 2 || e[0].Author != "claude" || e[1].Author != "alex" || e[1].Text != "Left.\n\nSee the plan." {
		t.Fatalf("entries: %+v", e)
	}
	for _, x := range e {
		if x.Recommendation || x.Source != nil {
			t.Errorf("an older entry has no marks: %+v", x)
		}
	}
	if p := View(r, th)["pending_recommendation"].(*Entry); p != nil || th.Marshal() != doc {
		t.Errorf("nothing is pending (%+v), and the file round-trips:\n%s", p, th.Marshal())
	}
	old := entryHeading.FindStringSubmatch("### 2026-09-17T21:01:00Z orchestrator (recommendation)")
	if old == nil || old[2] != "orchestrator (recommendation)" {
		t.Errorf("an older flai's heading pattern reads the mark into the author: %q", old)
	}
}

// ADR-0189: a thread is the orchestrator's own while the orchestrator opened
// it and wrote every entry on it; one entry by the operator, or a thread an
// agent opened, is not.
func TestAThreadIsTheOrchestratorsOwnUntilAnotherWritesOnIt(t *testing.T) {
	r := project(t)
	open := func(author, title string) *Thread {
		t.Helper()
		th, err := New(r, NewOptions{Title: title, On: "design/system/plan.md", Author: author, Text: "Is the gate built?", Now: t0})
		if err != nil {
			t.Fatal(err)
		}
		return th
	}
	alone := open("orchestrator", "The gate")
	alone, err := Reply(r, alone.ID, "orchestrator", "Still not built.", t0.Add(time.Minute))
	if err != nil {
		t.Fatal(err)
	}
	if !alone.OrchestratorsOwn() {
		t.Errorf("a thread the orchestrator alone wrote is its own: %s", alone.Body)
	}

	replied := open("orchestrator", "The flag")
	if replied, err = Reply(r, replied.ID, "alex", "It is, in cli/feature.go.", t0.Add(time.Minute)); err != nil {
		t.Fatal(err)
	}
	if replied.OrchestratorsOwn() {
		t.Errorf("a thread the operator replied on is not the orchestrator's own: %s", replied.Body)
	}

	agents := open("claude", "The shape")
	if agents, err = Reply(r, agents.ID, "orchestrator", "Left.", t0.Add(time.Minute)); err != nil {
		t.Fatal(err)
	}
	if agents.OrchestratorsOwn() {
		t.Errorf("a thread an agent opened is not the orchestrator's own: %s", agents.Body)
	}
}

// ADR-0206: a recommendation marked for the operator alone carries the mark's
// paragraph above its source line, and reads back with both; the mark needs
// a recommendation.
func TestARecommendationMarkedForTheOperatorAloneReadsBack(t *testing.T) {
	r := project(t)
	wipLint(t, r)
	if _, err := New(r, NewOptions{Title: "Which host", On: "design/system/plan.md", Author: "agent-S-0001", Text: "Buy a second host?", Now: t0}); err != nil {
		t.Fatal(err)
	}
	if _, err := ReplyWith(r, "TH-0001", "orchestrator", "Buy it.", t0.Add(time.Minute), Marks{OperatorOnly: true}); err == nil || !strings.Contains(err.Error(), "only a recommendation") {
		t.Errorf("the mark without a recommendation is refused: %v", err)
	}
	th, err := ReplyWith(r, "TH-0001", "orchestrator", "Buy it.", t0.Add(time.Minute), Marks{Recommendation: true, OperatorOnly: true, Source: Source{Path: "design/system/plan.md", Heading: "Shape"}})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(th.Body, "\n### 2026-09-17T21:01:00Z orchestrator (recommendation)\nBuy it.\n\nFor the operator alone: flai serve does not confirm this after orchestration.permissions.confirm_after.\n\nSource: design/system/plan.md § Shape\n") {
		t.Errorf("on disk:\n%s", th.Body)
	}
	lintClean(t, r, th)
	back, err := Read(th.Path)
	if err != nil {
		t.Fatal(err)
	}
	e := back.Entries()
	if len(e) != 2 || !e[1].OperatorOnly || !e[1].Recommendation || e[1].Source == nil || *e[1].Source != (Source{Path: "design/system/plan.md", Heading: "Shape"}) || e[0].OperatorOnly {
		t.Fatalf("entries: %+v", e)
	}
	if p := back.PendingRecommendation(); p == nil || !p.OperatorOnly {
		t.Errorf("the pending recommendation carries the mark: %+v", p)
	}

	// without a source the mark is the last paragraph, and still read
	th, err = ReplyWith(r, "TH-0001", "orchestrator", "Or lease one.", t0.Add(2*time.Minute), Marks{Recommendation: true, OperatorOnly: true})
	if err != nil {
		t.Fatal(err)
	}
	if e := th.Entries(); !e[2].OperatorOnly || e[2].Source != nil {
		t.Errorf("the mark without a source: %+v", e[2])
	}
	// a plain recommendation, and a plain entry quoting the words, carry none
	if _, err := ReplyWith(r, "TH-0001", "orchestrator", "Lease one.", t0.Add(3*time.Minute), Marks{Recommendation: true}); err != nil {
		t.Fatal(err)
	}
	th, err = Reply(r, "TH-0001", "agent-S-0001", "For the operator alone: I would ask.", t0.Add(4*time.Minute))
	if err != nil {
		t.Fatal(err)
	}
	if e := th.Entries(); e[3].OperatorOnly || e[4].OperatorOnly {
		t.Errorf("unmarked entries: %+v", e[3:])
	}
}

// ADR-0206: WroteSince tells whether an author the predicate accepts wrote
// after an entry, an entry of the same second by another author included.
func TestWroteSinceLooksAfterTheEntry(t *testing.T) {
	r := project(t)
	if _, err := New(r, NewOptions{Title: "Which shape", On: "design/system/plan.md", Author: "alex", Text: "Left or right?", Now: t0}); err != nil {
		t.Fatal(err)
	}
	th, err := ReplyWith(r, "TH-0001", "orchestrator", "Left.", t0.Add(time.Minute), Marks{Recommendation: true})
	if err != nil {
		t.Fatal(err)
	}
	isAlex := func(name string) bool { return name == "alex" }
	rec := th.Entries()[1]
	if th.WroteSince(rec, isAlex) {
		t.Error("nothing follows the recommendation")
	}
	if !th.WroteSince(th.Entries()[0], func(name string) bool { return name == "orchestrator" }) {
		t.Error("the orchestrator wrote after the opening entry")
	}
	if th, err = Reply(r, "TH-0001", "alex", "Hm.", t0.Add(time.Minute)); err != nil {
		t.Fatal(err)
	}
	if !th.WroteSince(rec, isAlex) {
		t.Errorf("alex wrote after the recommendation, in the same second:\n%s", th.Body)
	}
	if th.WroteSince(Entry{At: "2026-09-17T21:09:00Z", Author: "nobody"}, isAlex) {
		t.Error("an entry the thread lacks has nothing after it")
	}
}

// ADR-0206, I-0195: flai serve confirms a recommendation left unanswered past
// confirm_after with an entry of its own, which ends the pending
// recommendation and reads as a confirmation by timeout; it is refused with
// none pending, before confirm_after, on one marked for the operator alone,
// and when the operator wrote since.
func TestConfirmByTimeoutConfirmsAnUnansweredRecommendation(t *testing.T) {
	r := project(t)
	wipLint(t, r)
	after := 48 * time.Hour
	if _, err := New(r, NewOptions{Title: "Which shape", On: "design/system/plan.md", Author: "planner-S-0001", Text: "Left or right?", Now: t0}); err != nil {
		t.Fatal(err)
	}
	if _, err := ConfirmByTimeout(r, "TH-0001", after, t0.Add(after)); err == nil || !strings.Contains(err.Error(), "no recommendation") {
		t.Errorf("with none pending it is refused: %v", err)
	}
	rec := t0.Add(time.Minute)
	if _, err := ReplyWith(r, "TH-0001", "orchestrator", "Left.", rec, Marks{Recommendation: true, Source: Source{Path: "design/system/plan.md", Heading: "Shape"}}); err != nil {
		t.Fatal(err)
	}
	if _, err := ConfirmByTimeout(r, "TH-0001", after, rec.Add(after-time.Second)); err == nil || !strings.Contains(err.Error(), "has waited 47h59m59s, less than confirm_after 48h") {
		t.Errorf("before confirm_after it is refused: %v", err)
	}
	if _, err := ConfirmByTimeout(r, "TH-0001", 0, rec.Add(after)); err == nil || !strings.Contains(err.Error(), "longer than zero") {
		t.Errorf("a zero confirm_after is refused: %v", err)
	}
	th, err := ConfirmByTimeout(r, "th-1", after, rec.Add(after))
	if err != nil {
		t.Fatal(err)
	}
	if th.Status != "answered" || th.PendingRecommendation() != nil || !contains(th.Participants, TimeoutAuthor) {
		t.Fatalf("confirmed: %s %+v %v", th.Status, th.PendingRecommendation(), th.Participants)
	}
	e := th.Entries()
	last := e[len(e)-1]
	if len(e) != 3 || last.Author != "flai serve" || last.Text != "Confirmed the recommendation of 2026-09-17T21:01:00Z orchestrator after 48h with no answer from the operator (orchestration.permissions.confirm_after).\n\nSource: design/system/plan.md § Shape" || last.Source == nil {
		t.Fatalf("confirming entry: %+v", e)
	}
	if !last.Confirms() || !last.ConfirmsByTimeout() || e[1].ConfirmsByTimeout() {
		t.Errorf("a confirmation by timeout: %+v", last)
	}
	lintClean(t, r, th)
	if _, err := ConfirmByTimeout(r, "TH-0001", after, rec.Add(2*after)); err == nil {
		t.Error("a confirmed recommendation is not confirmed twice")
	}
	// the operator's confirmation is no confirmation by timeout
	if _, err := ReplyWith(r, "TH-0001", "orchestrator", "Right, then.", rec.Add(3*after), Marks{Recommendation: true}); err != nil {
		t.Fatal(err)
	}
	if th, err = Confirm(r, "TH-0001", "designer", rec.Add(3*after+time.Minute)); err != nil {
		t.Fatal(err)
	}
	if e := th.Entries(); !e[len(e)-1].Confirms() || e[len(e)-1].ConfirmsByTimeout() {
		t.Errorf("the operator's confirmation: %+v", e[len(e)-1])
	}

	// marked for the operator alone: refused, writing nothing
	if _, err := New(r, NewOptions{Title: "Which host", On: "design/system/plan.md", Author: "agent-S-0001", Text: "Buy one?", Now: t0}); err != nil {
		t.Fatal(err)
	}
	th, err = ReplyWith(r, "TH-0002", "orchestrator", "Buy it.", rec, Marks{Recommendation: true, OperatorOnly: true})
	if err != nil {
		t.Fatal(err)
	}
	was, _ := os.ReadFile(th.Path)
	if _, err := ConfirmByTimeout(r, "TH-0002", after, rec.Add(after)); err == nil || !strings.Contains(err.Error(), "operator alone") {
		t.Errorf("one marked for the operator alone is refused: %v", err)
	}
	if got, _ := os.ReadFile(th.Path); string(got) != string(was) {
		t.Error("a refused confirmation leaves the thread as it was")
	}

	// the operator, here the opener, wrote after the recommendation: refused
	if _, err := New(r, NewOptions{Title: "Which risk", On: "design/system/plan.md", Author: "designer", Text: "Which?", Now: t0}); err != nil {
		t.Fatal(err)
	}
	if _, err := ReplyWith(r, "TH-0003", "orchestrator", "The first.", rec, Marks{Recommendation: true}); err != nil {
		t.Fatal(err)
	}
	th, err = Reply(r, "TH-0003", "designer", "Let me think.", rec.Add(time.Hour))
	if err != nil || th.PendingRecommendation() == nil {
		t.Fatalf("the opener's follow-up leaves it pending: %v", err)
	}
	if _, err := ConfirmByTimeout(r, "TH-0003", after, rec.Add(after)); err == nil || !strings.Contains(err.Error(), "the operator wrote on TH-0003") {
		t.Errorf("after the operator wrote it is refused: %v", err)
	}
}

func TestReadableDurationDropsZeroUnits(t *testing.T) {
	for d, want := range map[time.Duration]string{48 * time.Hour: "48h", 90 * time.Minute: "1h30m", 30 * time.Minute: "30m", 90 * time.Second: "1m30s", 2*time.Hour + 10*time.Minute: "2h10m", time.Hour + 5*time.Second: "1h0m5s"} {
		if got := readableDuration(d); got != want {
			t.Errorf("%s: %q, want %q", d, got, want)
		}
	}
}

// ADR-0206: the openers flai serve may confirm on are a story's agent and a
// planner.
func TestOpenedByAgentIsAStorysAgentOrAPlanner(t *testing.T) {
	for opener, want := range map[string]bool{"agent-S-0001": true, "planner-S-0002": true, "planner-E-0003": true, "alex": false, "orchestrator": false, "flai serve": false, "agent": false, "agent-S-0001-x": false, "planner": false} {
		th := &Thread{Body: "\n## Entries\n\n### 2026-09-17T21:00:00Z " + opener + "\nWhich?\n"}
		if got := th.OpenedByAgent(); got != want {
			t.Errorf("opened by %q: %v, want %v", opener, got, want)
		}
	}
}
