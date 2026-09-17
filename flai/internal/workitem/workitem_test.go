package workitem

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"time"
)

var t0 = time.Date(2026, 9, 15, 20, 0, 0, 0, time.UTC)

// newProject writes a minimal conforming repo and returns it opened.
func newProject(t *testing.T) *Repo {
	t.Helper()
	root := t.TempDir()
	manifest := "version: 1\nname: t\nkey: t\nlayout:\n  design: design\n  docs: docs\n  wip: wip\n"
	if err := os.WriteFile(filepath.Join(root, "system-flow.yaml"), []byte(manifest), 0o644); err != nil {
		t.Fatal(err)
	}
	for _, d := range []string{"wip/kanban/epics", "wip/kanban/stories", "wip/kanban/tasks", "wip/agents", "wip/archive"} {
		_ = os.MkdirAll(filepath.Join(root, d), 0o755)
	}
	r, err := Open(filepath.Join(root, "wip"))
	if err != nil {
		t.Fatal(err)
	}
	return r
}

func mustCreate(t *testing.T, r *Repo, typ, title, parent string) *Item {
	t.Helper()
	it, err := r.Create(NewOptions{Type: typ, Title: title, Parent: parent, Now: t0})
	if err != nil {
		t.Fatalf("create %s %q: %v", typ, title, err)
	}
	return it
}

func mustMove(t *testing.T, r *Repo, it *Item, to string, reason string) {
	t.Helper()
	items, _ := r.List(false)
	board, _ := r.LoadBoard()
	if _, err := r.Move(it, to, MoveOptions{Now: t0, Items: items, Board: board, Reason: reason, By: "test"}); err != nil {
		t.Fatalf("move %s to %s: %v", it.ID, to, err)
	}
	if err := r.Save(it); err != nil {
		t.Fatal(err)
	}
	if err := board.Save("2026-09-15"); err != nil {
		t.Fatal(err)
	}
}

func TestRoundTripRepositoryItems(t *testing.T) {
	if testing.Short() {
		t.Skip("integration: reads the monorepo")
	}
	// Every hand-written item in the monorepo must survive parse and marshal
	// byte for byte, so flai edits never churn unrelated lines.
	root := filepath.Join("..", "..", "..")
	if _, err := os.Stat(filepath.Join(root, "system-flow.yaml")); err != nil {
		t.Skip("monorepo not present")
	}
	r, err := Open(root)
	if err != nil {
		t.Fatal(err)
	}
	items, err := r.List(true)
	if err != nil {
		t.Fatal(err)
	}
	if len(items) < 20 {
		t.Fatalf("expected the bootstrapped items, got %d", len(items))
	}
	for _, it := range items {
		if err := it.Validate(); err != nil {
			t.Errorf("%s: %v", it.Path, err)
		}
		orig, _ := os.ReadFile(it.Path)
		if got := it.Marshal(); got != string(orig) {
			t.Errorf("%s: marshal differs from file\n--- got ---\n%s\n--- want ---\n%s", it.Path, got, orig)
		}
	}
	for _, typ := range Types {
		id, _ := r.NextID(typ)
		if !regexp.MustCompile(`^` + strings.ToUpper(typ[:1]) + `-\d{3,}$`).MatchString(id) {
			t.Errorf("NextID(%s) = %s", typ, id)
		}
	}
}

func TestScalar(t *testing.T) {
	cases := map[string]string{
		"plain title":         "plain title",
		"":                    `""`,
		"Three: four":         `"Three: four"`,
		"true":                `"true"`,
		"12":                  `"12"`,
		"ends with colon:":    `"ends with colon:"`,
		"has # comment":       `"has # comment"`,
		"flai new from x.y":   "flai new from x.y",
		"Tab\there":           `"Tab\there"`,
		"- leading dash":      `"- leading dash"`,
		"quote \"inside\" it": `"quote \"inside\" it"`,
		"what? really! (yes)": "what? really! (yes)",
	}
	for in, want := range cases {
		if got := Scalar(in); got != want {
			t.Errorf("Scalar(%q) = %s, want %s", in, got, want)
		}
	}
}

func TestValidate(t *testing.T) {
	good := &Item{ID: "S-0001", Type: Story, Nature: "feature", Title: "x", Status: Ready, Parent: "E-0001", Owner: "a",
		Created: "2026-09-15T10:00:00Z", Updated: "2026-09-15T10:00:00Z", Transitions: []Transition{{To: Ready, At: "2026-09-15T10:00:00Z", By: "a"}}}
	if err := good.Validate(); err != nil {
		t.Fatalf("valid item rejected: %v", err)
	}
	bad := *good
	bad.Status = InProgress
	if err := bad.Validate(); err == nil || !strings.Contains(err.Error(), "does not match the last transition") {
		t.Errorf("status mismatch not caught: %v", err)
	}
	bad = *good
	bad.Nature = "bug"
	if err := bad.Validate(); err == nil || !strings.Contains(err.Error(), "nature") {
		t.Errorf("bad nature not caught: %v", err)
	}
	bad = *good
	bad.Parent = "T-0001"
	if err := bad.Validate(); err == nil || !strings.Contains(err.Error(), "parent must be an epic") {
		t.Errorf("wrong parent type not caught: %v", err)
	}
	bad = *good
	bad.Created = "yesterday"
	if err := bad.Validate(); err == nil || !strings.Contains(err.Error(), "created") {
		t.Errorf("bad timestamp not caught: %v", err)
	}
	bad = *good
	bad.Blocked = []Block{{From: "2026-09-15T10:00:00Z", Reason: "a"}, {From: "2026-09-15T11:00:00Z", Reason: "b"}}
	if err := bad.Validate(); err == nil || !strings.Contains(err.Error(), "more than one open") {
		t.Errorf("two open blocks not caught: %v", err)
	}
}

func TestCreateAndLinkParents(t *testing.T) {
	r := newProject(t)
	e := mustCreate(t, r, Epic, "Ship it: now", "")
	if e.ID != "E-0001" || e.Status != Backlog || !strings.HasSuffix(e.Path, "E-0001-ship-it-now.md") {
		t.Fatalf("epic: %+v", e)
	}
	s := mustCreate(t, r, Story, "First slice", "E-0001")
	tk := mustCreate(t, r, Task, "Do a thing", "S-0001")
	if tk.Stream != "S-0001" || tk.Parent != "S-0001" {
		t.Errorf("task: %+v", tk)
	}
	e2, _ := r.Get("E-0001")
	if !strings.Contains(e2.Body, "## Stories\n- S-0001 First slice\n") {
		t.Errorf("epic body not linked:\n%s", e2.Body)
	}
	s2, _ := r.Get("S-0001")
	if !strings.Contains(s2.Body, "## Tasks\n- T-0001 Do a thing\n") {
		t.Errorf("story body not linked:\n%s", s2.Body)
	}
	mustCreate(t, r, Task, "Second thing", "S-0001")
	s3, _ := r.Get("S-0001")
	if !strings.Contains(s3.Body, "- T-0001 Do a thing\n- T-0002 Second thing\n\n## Notes") {
		t.Errorf("second task not appended cleanly:\n%s", s3.Body)
	}
	orig, _ := os.ReadFile(s3.Path)
	if s3.Marshal() != string(orig) {
		t.Error("created story does not round-trip")
	}
	if _, err := r.Create(NewOptions{Type: Task, Title: "orphan", Parent: "E-0001", Now: t0}); err == nil {
		t.Error("task under epic should fail")
	}
	if _, err := r.Create(NewOptions{Type: Story, Title: "x", Now: t0}); err == nil {
		t.Error("story without parent should fail")
	}
	if _, err := r.Create(NewOptions{Type: Story, Title: "x", Parent: "E-0001", Nature: "bug", Now: t0}); err == nil {
		t.Error("bad nature should fail")
	}
	_ = s
}

func TestMoveRules(t *testing.T) {
	r := newProject(t)
	mustCreate(t, r, Epic, "E", "")
	s := mustCreate(t, r, Story, "S", "E-0001")
	items, _ := r.List(false)
	board, _ := r.LoadBoard()
	if _, err := r.Move(s, Ready, MoveOptions{Now: t0, Items: items, Board: board}); err == nil || !strings.Contains(err.Error(), "needs at least one task") {
		t.Errorf("ready without tasks: %v", err)
	}
	mustCreate(t, r, Task, "T", "S-0001")
	items, _ = r.List(false)
	s, _ = r.Get("S-0001")
	if _, err := r.Move(s, Ready, MoveOptions{Now: t0, Items: items, Board: board}); err == nil || !strings.Contains(err.Error(), "Acceptance criteria") {
		t.Errorf("ready without criteria: %v", err)
	}
	s.Body = strings.Replace(s.Body, "## Acceptance criteria\n- [ ]\n", "## Acceptance criteria\n- [ ] works\n", 1)
	_ = r.Save(s)
	mustMove(t, r, s, Ready, "")
	board, _ = r.LoadBoard()
	if len(board.Order) != 1 || board.Order[0] != "S-0001" {
		t.Errorf("ready story not added to pull order: %v", board.Order)
	}
	if _, err := r.Move(s, Done, MoveOptions{Now: t0, Items: items}); err == nil {
		t.Error("ready to done must be refused")
	}
	if _, err := r.Move(s, Cancelled, MoveOptions{Now: t0, Items: items}); err == nil || !strings.Contains(err.Error(), "--reason") {
		t.Errorf("cancel without reason: %v", err)
	}
	mustMove(t, r, s, InProgress, "")
	mustMove(t, r, s, Review, "")
	items, _ = r.List(false)
	if _, err := r.Move(s, Done, MoveOptions{Now: t0, Items: items}); err == nil || !strings.Contains(err.Error(), "T-0001 is backlog") {
		t.Errorf("done with open task: %v", err)
	}
	tk, _ := r.Get("T-0001")
	mustMove(t, r, tk, Ready, "")
	mustMove(t, r, tk, InProgress, "")
	mustMove(t, r, tk, Done, "") // tasks may skip review
	items, _ = r.List(false)
	s, _ = r.Get("S-0001")
	if _, err := r.Move(s, Done, MoveOptions{Now: t0, Items: items}); err == nil || !strings.Contains(err.Error(), "unchecked") {
		t.Errorf("done with unchecked criteria: %v", err)
	}
	s.Body = strings.Replace(s.Body, "- [ ] works", "- [x] works", 1)
	mustMove(t, r, s, Done, "")
	if s.Status != Done || len(s.Transitions) != 4 || s.Transitions[3].By != "test" {
		t.Errorf("history: %+v", s.Transitions)
	}
	if err := s.Validate(); err != nil {
		t.Errorf("moved item invalid: %v", err)
	}
	// review -> in-progress needs a reason and records it in Notes
	mustCreate(t, r, Story, "S2", "E-0001")
	mustCreate(t, r, Task, "T2", "S-0002")
	s2, _ := r.Get("S-0002")
	s2.Body = strings.Replace(s2.Body, "- [ ]\n", "- [ ] ok\n", 1)
	_ = r.Save(s2)
	mustMove(t, r, s2, Ready, "")
	mustMove(t, r, s2, InProgress, "")
	mustMove(t, r, s2, Review, "")
	mustMove(t, r, s2, InProgress, "tests missing")
	if !strings.Contains(s2.Body, "moved to in-progress: tests missing") {
		t.Errorf("reason not recorded:\n%s", s2.Body)
	}
}

func TestWIPWarningAndBlocks(t *testing.T) {
	r := newProject(t)
	mustCreate(t, r, Epic, "E", "")
	var stories []*Item
	for i := 0; i < 3; i++ {
		s := mustCreate(t, r, Story, "S", "E-0001")
		mustCreate(t, r, Task, "T", s.ID)
		s, _ = r.Get(s.ID)
		s.Body = strings.Replace(s.Body, "- [ ]\n", "- [ ] ok\n", 1)
		_ = r.Save(s)
		mustMove(t, r, s, Ready, "")
		stories = append(stories, s)
	}
	var warned bool
	for _, s := range stories {
		items, _ := r.List(false)
		board, _ := r.LoadBoard()
		w, err := r.Move(s, InProgress, MoveOptions{Now: t0, Items: items, Board: board})
		if err != nil {
			t.Fatal(err)
		}
		_ = r.Save(s)
		if len(w) > 0 {
			warned = true
		}
	}
	if !warned {
		t.Error("third in-progress story should warn about the WIP limit of 2")
	}
	s := stories[0]
	if err := BlockItem(s, "waiting", t0); err != nil {
		t.Fatal(err)
	}
	if err := BlockItem(s, "again", t0); err == nil {
		t.Error("double block should fail")
	}
	if !s.IsBlocked() {
		t.Error("should be blocked")
	}
	if err := UnblockItem(s, t0.Add(time.Hour)); err != nil {
		t.Fatal(err)
	}
	if s.IsBlocked() || s.Blocked[0].Until == "" {
		t.Error("unblock did not close interval")
	}
	if err := UnblockItem(s, t0); err == nil {
		t.Error("unblock when not blocked should fail")
	}
	_ = r.Save(s)
	back, _ := r.Get(s.ID)
	if len(back.Blocked) != 1 || back.Blocked[0].Reason != "waiting" {
		t.Errorf("blocked not persisted: %+v", back.Blocked)
	}
}

func TestStreamsAndIndex(t *testing.T) {
	r := newProject(t)
	mustCreate(t, r, Epic, "E", "")
	s := mustCreate(t, r, Story, "Stream me", "E-0001")
	n, err := r.OpenStream(s, StreamOptions{Agent: "bot", Session: "abc", Now: t0})
	if err != nil {
		t.Fatal(err)
	}
	if n.Stream != "S-0001" || n.Agent != "bot" || !strings.Contains(n.Body, "# S-0001 Stream me") || !strings.Contains(n.Body, "## Log\n\n### 2026-09-15T20:00:00Z\nStream opened.") {
		t.Errorf("narrative: %+v\n%s", n, n.Body)
	}
	if _, err := r.OpenStream(s, StreamOptions{Now: t0}); err == nil {
		t.Error("second open should fail")
	}
	later := t0.Add(30 * time.Minute)
	n2, err := r.LogStream("S-0001", "did a thing", StreamOptions{Agent: "bot2", Now: later})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasSuffix(n2.Body, "### 2026-09-15T20:30:00Z\ndid a thing\n") || n2.Updated != "2026-09-15T20:30:00Z" || n2.Agent != "bot2" {
		t.Errorf("log: %+v\n%s", n2, n2.Body)
	}
	n3, err := r.LogStream("S-0001", "same second", StreamOptions{Now: later})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasSuffix(n3.Body, "### 2026-09-15T20:30:00Z\ndid a thing\n\nsame second\n") || strings.Count(n3.Body, "### 2026-09-15T20:30:00Z") != 1 {
		t.Errorf("same-second log should share the heading:\n%s", n3.Body)
	}
	if _, err := r.LogStream("S-0009", "x", StreamOptions{Now: t0}); err == nil {
		t.Error("log to missing stream should fail")
	}
	items, _ := r.List(false)
	if err := r.WriteIndex(items, later); err != nil {
		t.Fatal(err)
	}
	idx, _ := os.ReadFile(filepath.Join(r.AgentsDir(), "index.md"))
	if !strings.Contains(string(idx), "| [S-0001](S-0001.md) | Stream me | backlog | bot2 | 2026-09-15T20:30:00Z |") {
		t.Errorf("index:\n%s", idx)
	}
}

func TestArchive(t *testing.T) {
	r := newProject(t)
	mustCreate(t, r, Epic, "E", "")
	s := mustCreate(t, r, Story, "S", "E-0001")
	mustCreate(t, r, Task, "T", "S-0001")
	_, _ = r.OpenStream(s, StreamOptions{Now: t0})
	items, _ := r.List(false)
	if _, err := r.PlanArchive(items, []string{"S-0001"}); err == nil {
		t.Error("archiving an open story should fail")
	}
	plan, _ := r.PlanArchive(items, nil)
	if len(plan.Items) != 0 {
		t.Errorf("nothing should be eligible yet: %v", plan.Items)
	}
	s, _ = r.Get("S-0001")
	s.Body = strings.Replace(s.Body, "- [ ]\n", "- [x] ok\n", 1)
	_ = r.Save(s)
	for _, st := range []string{Ready, InProgress, Done} {
		tk, _ := r.Get("T-0001")
		mustMove(t, r, tk, st, "")
	}
	for _, st := range []string{Ready, InProgress, Review, Done} {
		s, _ = r.Get("S-0001")
		mustMove(t, r, s, st, "")
	}
	items, _ = r.List(false)
	plan, err := r.PlanArchive(items, nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(plan.Items) != 2 || len(plan.Narratives) != 1 {
		t.Fatalf("plan: %d items %d narratives", len(plan.Items), len(plan.Narratives))
	}
	if err := r.Archive(plan); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(r.ArchiveDir(), "kanban", "stories", "S-0001-s.md")); err != nil {
		entries, _ := os.ReadDir(filepath.Join(r.ArchiveDir(), "kanban", "stories"))
		names := []string{}
		for _, e := range entries {
			names = append(names, e.Name())
		}
		ids := []string{}
		for _, it := range plan.Items {
			ids = append(ids, it.ID+"@"+it.Path)
		}
		t.Errorf("story not archived; archive has %v, plan was %v", names, ids)
	}
	if _, err := os.Stat(filepath.Join(r.ArchiveDir(), "agents", "S-0001.md")); err != nil {
		t.Error("narrative not archived")
	}
	got, _ := r.Get("S-0001")
	if !got.Archived {
		t.Error("Get should find archived items")
	}
	for _, typ := range Types {
		if got := Folder(typ); got != map[string]string{Epic: "epics", Story: "stories", Task: "tasks"}[typ] {
			t.Errorf("Folder(%s) = %s", typ, got)
		}
	}
	if id, _ := r.NextID(Story); id != "S-0002" {
		t.Errorf("NextID must count archived items: %s", id)
	}
	items, _ = r.List(false)
	if _, err := r.PlanArchive(items, []string{"E-0001"}); err == nil {
		t.Error("open epic should not archive")
	}
}
