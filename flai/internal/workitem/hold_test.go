package workitem

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/manifest"
)

var holdProjects = []manifest.Project{
	{Name: "flai", Path: "flai", Tags: []string{"cli"}},
	{Name: "flaiover", Path: "flaiover", Tags: []string{"dashboard"}},
}

func claimed(id, status string, touches ...string) *Item {
	return &Item{ID: id, Type: Story, Status: status, Touches: touches}
}

func claimTask(id, parent, status string, touches ...string) *Item {
	return &Item{ID: id, Type: Task, Parent: parent, Status: status, Touches: touches}
}

// S-0128: a claim is the story's touches and its open tasks', a component's
// name or tag read as its path, without duplicates; a done task's touch
// outside the story's is not claimed (ADR-0096).
func TestClaimTakesOpenTasksAndReadsComponentsAsPaths(t *testing.T) {
	s := claimed("S-0001", InProgress, "cli", "docs/")
	items := []*Item{s,
		claimTask("T-0001", "S-0001", InProgress, "flaiover/src", "flai"),
		claimTask("T-0002", "S-0001", Done, "design/system"),
		claimTask("T-0003", "S-0001", Cancelled, "template"),
		claimTask("T-0004", "S-0002", InProgress, "wip"),
		{ID: "T-0005", Type: Task, Parent: "S-0001", Status: InProgress, Touches: []string{"scripts"}, Archived: true},
	}
	got := strings.Join(NewHolds(items, holdProjects).Claim(s), ",")
	if got != "flai,docs,flaiover/src" {
		t.Errorf("claim = %s, want flai,docs,flaiover/src", got)
	}
}

// ADR-0096: a story's folder touch is narrowed to the touches its tasks name
// inside it, done ones included and cancelled ones not; a folder no task names
// inside, or one a task names whole, stays whole; a done task's touch outside
// every story touch is not added, an open one's is.
func TestClaimNarrowsAFolderToItsTasksTouches(t *testing.T) {
	cases := []struct {
		name    string
		touches []string
		tasks   []*Item
		want    string
	}{
		{"no task names inside",
			[]string{"flai/internal/mcpserver", "docs/users/flai.md"},
			[]*Item{claimTask("T-0001", "S-0001", InProgress, "design/system/x.md")},
			"flai/internal/mcpserver,docs/users/flai.md,design/system/x.md"},
		{"an open task narrows",
			[]string{"flai/internal/mcpserver"},
			[]*Item{claimTask("T-0001", "S-0001", InProgress, "flai/internal/mcpserver/permission.go")},
			"flai/internal/mcpserver/permission.go"},
		{"a done task's files stay claimed",
			[]string{"flai/internal/mcpserver"},
			[]*Item{
				claimTask("T-0001", "S-0001", Done, "flai/internal/mcpserver/permission.go"),
				claimTask("T-0002", "S-0001", Ready, "flai/internal/mcpserver/tools.go"),
			},
			"flai/internal/mcpserver/permission.go,flai/internal/mcpserver/tools.go"},
		{"a cancelled task's do not",
			[]string{"flai/internal/mcpserver"},
			[]*Item{
				claimTask("T-0001", "S-0001", Cancelled, "flai/internal/mcpserver/permission.go"),
				claimTask("T-0002", "S-0001", Done, "flai/internal/mcpserver/tools.go"),
			},
			"flai/internal/mcpserver/tools.go"},
		{"a folder only a cancelled task named inside stays whole",
			[]string{"flai/internal/mcpserver"},
			[]*Item{claimTask("T-0001", "S-0001", Cancelled, "flai/internal/mcpserver/permission.go")},
			"flai/internal/mcpserver"},
		{"a task that names the folder itself keeps it whole",
			[]string{"flai/internal/mcpserver"},
			[]*Item{
				claimTask("T-0001", "S-0001", InProgress, "flai/internal/mcpserver/permission.go"),
				claimTask("T-0002", "S-0001", Done, "flai/internal/mcpserver/"),
			},
			"flai/internal/mcpserver"},
		{"a task that names a folder around it keeps it whole",
			[]string{"flai/internal/mcpserver"},
			[]*Item{
				claimTask("T-0001", "S-0001", InProgress, "flai/internal/mcpserver/permission.go"),
				claimTask("T-0002", "S-0001", Done, "cli"),
			},
			"flai/internal/mcpserver"},
		{"a component narrows as its path",
			[]string{"cli", "docs"},
			[]*Item{claimTask("T-0001", "S-0001", Done, "flai/cmd/x.go", "docs/users/flai.md")},
			"flai/cmd/x.go,docs/users/flai.md"},
		{"a done task's touch outside is not added, an open one's is",
			[]string{"flai/cmd"},
			[]*Item{
				claimTask("T-0001", "S-0001", Done, "flai/cmd/x.go", "design/adrs/0096.md"),
				claimTask("T-0002", "S-0001", InProgress, "flai/cmd/y.go", "docs/users/flai.md"),
			},
			"flai/cmd/x.go,flai/cmd/y.go,docs/users/flai.md"},
		{"one folder narrows, another stays whole",
			[]string{"flai/internal/mcpserver", "docs/users"},
			[]*Item{claimTask("T-0001", "S-0001", InProgress, "flai/internal/mcpserver/permission.go")},
			"flai/internal/mcpserver/permission.go,docs/users"},
		{"no story touches: the open tasks' alone",
			nil,
			[]*Item{claimTask("T-0001", "S-0001", InProgress, "flai/cmd"), claimTask("T-0002", "S-0001", Done, "docs")},
			"flai/cmd"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			s := claimed("S-0001", InProgress, c.touches...)
			got := strings.Join(NewHolds(append([]*Item{s}, c.tasks...), holdProjects).Claim(s), ",")
			if got != c.want {
				t.Errorf("claim = %s\nwant    %s", got, c.want)
			}
		})
	}
}

// I-0087, ADR-0096: S-0283 claimed flai/internal/mcpserver and held every
// ready story under it, though its one task named only permission.go. Its
// claim narrowed to that file, it holds a story touching plan.go no longer,
// and still holds one touching permission.go.
func TestNarrowedClaimHoldsOnlyTheFilesItsTasksName(t *testing.T) {
	open := claimed("S-0283", InProgress, "flai/internal/mcpserver")
	task := claimTask("T-0001", "S-0283", InProgress, "flai/internal/mcpserver/permission.go")
	plan := claimed("S-0290", Ready, "flai/internal/mcpserver/plan.go")
	perm := claimed("S-0291", Ready, "flai/internal/mcpserver/permission.go")
	h := NewHolds([]*Item{open, task, plan, perm}, holdProjects)
	if got := h.Of(plan); got != nil {
		t.Errorf("plan.go is held: %s", got.Reason)
	}
	want := "held (overlap): touches flai/internal/mcpserver/permission.go, which S-0283 (in progress) touches too; starts when S-0283 moves to review, is cancelled, or is sent back"
	if got := h.Of(perm); got == nil || got.Reason != want {
		t.Errorf("permission.go: hold = %+v\nwant reason %s", got, want)
	}

	task.Status = Done // the branch changed it, so it stays claimed
	if got := NewHolds([]*Item{open, task, plan, perm}, holdProjects).Of(perm); got == nil {
		t.Error("a done task's file no longer holds")
	}
	task.Status = Cancelled // nothing names inside: the folder is whole again
	if got := NewHolds([]*Item{open, task, plan, perm}, holdProjects).Of(plan); got == nil {
		t.Error("with its only task cancelled, the folder holds plan.go no longer")
	}
}

func TestPathsOverlapByPrefix(t *testing.T) {
	for _, c := range []struct {
		a, b string
		want bool
	}{
		{"flai/cmd", "flai/cmd", true},
		{"flai/cmd", "flai/cmd/serve", true},
		{"flai/cmd/serve/", "flai/cmd", true},
		{"flai/cmd", "flaiover", false},
		{"flai", "flaiover/src", false},
		{"docs/users", "docs/operators", false},
		{"", "", true},
		{"", "flai", false},
		{"", "/flai", true},
		{"flai/", "flai", true},
		{"flaiover", "flai", false},
		{"flai/cmd", "flai/", true},
		{"flai//", "flai/cmd", false},
	} {
		if got := PathsOverlap(c.a, c.b); got != c.want {
			t.Errorf("PathsOverlap(%q, %q) = %v", c.a, c.b, got)
		}
	}
}

// I-0155: Holds that share a memo judge as Holds without one, from one
// reading to the next, while a story's claim is forgotten whenever its tasks
// change; two stories with one ID are told apart.
func TestHoldsWithAMemoJudgeAsWithout(t *testing.T) {
	open := claimed("S-0001", InProgress, "flai/cmd", "docs/users")
	task := claimTask("T-0001", "S-0001", InProgress, "flai/cmd/serve.go")
	ready := []*Item{
		claimed("S-0009", Ready, "flai/cmd/board.go"),
		waiting("S-0001"),
		claimed("S-0010", Ready, "docs/users/flai.md"),
		claimed("S-0011", Ready, "docs/users/other.md", "flaiover"),
		claimed("S-0012", Ready),
	}
	shared := manifest.Claims{Shared: []string{"docs/users/flai.md"}}
	judge := func(h *Holds) []string {
		h.Open(ready[3], "its agent started")
		var out []string
		for _, s := range ready {
			if hold := h.Of(s); hold != nil {
				out = append(out, s.ID+": "+hold.Reason)
				continue
			}
			out = append(out, s.ID+": not held")
		}
		return out
	}
	memo := NewHoldsMemo()
	for _, status := range []string{InProgress, Done, Cancelled, Ready, Cancelled} {
		task.Status = status
		memo.Forget(task.Parent)
		items := append([]*Item{open, task}, ready...)
		got := judge(memo.NewHolds(items, holdProjects).WithShared(shared))
		want := judge(NewHolds(items, holdProjects).WithShared(shared))
		if strings.Join(got, "\n") != strings.Join(want, "\n") {
			t.Errorf("task %s, with a memo:\n%s\nwithout:\n%s", status, strings.Join(got, "\n"), strings.Join(want, "\n"))
		}
	}
}

// S-0128, ADR-0046: which ready stories are held, and what each reason says.
func TestHoldsByOverlapAndByEmptyClaims(t *testing.T) {
	cases := []struct {
		name  string
		open  []*Item
		ready *Item
		code  string // empty: not held
		why   string
	}{
		{"nothing open", nil, claimed("S-0009", Ready, "flai"), "", ""},
		{"inside an open story's path",
			[]*Item{claimed("S-0001", InProgress, "flai/cmd")}, claimed("S-0009", Ready, "flai/cmd/serve"),
			HoldOverlap, "held (overlap): touches flai/cmd/serve, inside flai/cmd which S-0001 (in progress) touches; starts when S-0001 moves to review, is cancelled, or is sent back"},
		{"the same path, in progress",
			[]*Item{claimed("S-0001", InProgress, "docs")}, claimed("S-0009", Ready, "docs/"),
			HoldOverlap, "held (overlap): touches docs, which S-0001 (in progress) touches too; starts when S-0001 moves to review, is cancelled, or is sent back"},
		{"the same path, in review, holds nothing",
			[]*Item{claimed("S-0001", Review, "docs")}, claimed("S-0009", Ready, "docs/"),
			"", ""},
		{"around an open story's path",
			[]*Item{claimed("S-0001", InProgress, "flai/cmd")}, claimed("S-0009", Ready, "cli"),
			HoldOverlap, "held (overlap): touches flai, which holds flai/cmd that S-0001 (in progress) touches; starts when S-0001 moves to review, is cancelled, or is sent back"},
		{"a sibling path is not held",
			[]*Item{claimed("S-0001", InProgress, "flai")}, claimed("S-0009", Ready, "flaiover", "docs"),
			"", ""},
		{"backlog, done, and ready stories hold nothing",
			[]*Item{claimed("S-0001", Backlog, "flai"), claimed("S-0002", Done, "flai"), claimed("S-0003", Ready, "flai")}, claimed("S-0009", Ready, "flai"),
			"", ""},
		{"an open task widens the open story's claim",
			[]*Item{claimed("S-0001", InProgress, "docs"), claimTask("T-0001", "S-0001", InProgress, "dashboard")}, claimed("S-0009", Ready, "flaiover/src"),
			HoldOverlap, "held (overlap): touches flaiover/src, inside flaiover which S-0001 (in progress) touches; starts when S-0001 moves to review, is cancelled, or is sent back"},
		{"held by two",
			[]*Item{claimed("S-0002", InProgress, "docs"), claimed("S-0001", InProgress, "flai")}, claimed("S-0009", Ready, "flai/cmd", "docs/users"),
			HoldOverlap, "held (overlap): touches flai/cmd, inside flai which S-0001 (in progress) touches; touches docs/users, inside docs which S-0002 (in progress) touches; starts when S-0001 and S-0002 move to review, are cancelled, or are sent back"},
		{"held by the one in progress, not the one in review",
			[]*Item{claimed("S-0002", Review, "docs"), claimed("S-0001", InProgress, "flai")}, claimed("S-0009", Ready, "flai/cmd", "docs/users"),
			HoldOverlap, "held (overlap): touches flai/cmd, inside flai which S-0001 (in progress) touches; starts when S-0001 moves to review, is cancelled, or is sent back"},
		{"a ready story with no touches",
			[]*Item{claimed("S-0001", InProgress, "flai"), claimed("S-0003", InProgress, "docs"), claimed("S-0002", Review, "docs")}, claimed("S-0009", Ready),
			HoldNoTouches, "held (no-touches): declares no touches, so it may change what S-0001 (in progress) and S-0003 (in progress) change; starts when it declares touches that overlap no story's in progress, or when S-0001 and S-0003 move to review, are cancelled, or are sent back"},
		{"a ready story with no touches and only a story in review",
			[]*Item{claimed("S-0002", Review, "docs")}, claimed("S-0009", Ready),
			"", ""},
		{"an open story with no touches",
			[]*Item{claimed("S-0001", InProgress)}, claimed("S-0009", Ready, "flai"),
			HoldNoTouches, "held (no-touches): S-0001 (in progress) declares no touches, so it may change anything; starts when S-0001 moves to review, is cancelled, or is sent back"},
		{"a story in review with no touches holds nothing",
			[]*Item{claimed("S-0001", Review)}, claimed("S-0009", Ready, "flai"),
			"", ""},
		{"no touches and nothing open", nil, claimed("S-0009", Ready), "", ""},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			h := NewHolds(append(c.open, c.ready), holdProjects).Of(c.ready)
			switch {
			case c.code == "" && h != nil:
				t.Fatalf("held: %s", h.Reason)
			case c.code == "":
			case h == nil:
				t.Fatal("not held")
			case h.Code != c.code || h.Reason != c.why:
				t.Errorf("hold = %s: %s\nwant   %s: %s", h.Code, h.Reason, c.code, c.why)
			}
		})
	}
}

// A story the launcher has started counts as open for the rest of its look.
func TestOpenedStoryHoldsTheNext(t *testing.T) {
	a, b := claimed("S-0001", Ready, "flai"), claimed("S-0002", Ready, "flai/cmd")
	h := NewHolds([]*Item{a, b}, holdProjects)
	if h.Of(b) != nil {
		t.Fatal("held before anything was opened")
	}
	h.Open(a, "its agent started")
	if got := h.Of(b); got == nil || !strings.Contains(got.Reason, "S-0001 (its agent started)") {
		t.Errorf("hold = %+v", got)
	}
	if h.Of(a) != nil {
		t.Error("a story holds itself")
	}
}

// The board marks held ready stories, and FirstPullable skips them.
func TestBoardMarksHeldStories(t *testing.T) {
	items := []*Item{
		claimed("S-0001", InProgress, "flai/cmd"),
		claimed("S-0002", Ready, "flai"),
		claimed("S-0003", Ready, "docs"),
		claimed("S-0004", Backlog, "flai"),
	}
	board := &Board{Order: []string{"S-0002", "S-0003"}}
	v := NewBoardView(items, board, time.Now(), false, nil, NewHolds(items, holdProjects))
	ready := v.ReadyInPullOrder()
	if len(ready) != 2 || ready[0].Held == nil || ready[0].Held.Code != HoldOverlap || ready[1].Held != nil {
		t.Fatalf("ready = %+v", ready)
	}
	if first := v.FirstPullable(); first == nil || first.ID != "S-0003" {
		t.Errorf("first pullable = %+v", first)
	}
	for _, c := range append(v.Columns[Backlog], v.Columns[InProgress]...) {
		if c.Held != nil {
			t.Errorf("%s is %s and marked held", c.ID, c.Status)
		}
	}
	items[2].Touches = []string{"flai/cmd/board.go"}
	if first := NewBoardView(items, board, time.Now(), false, nil, NewHolds(items, holdProjects)).FirstPullable(); first != nil {
		t.Errorf("every ready story is held, yet %s is pullable", first.ID)
	}
}

// I-0087, ADR-0096: S-0220 waited in review for 2h45m with a claim that
// overlapped all ten ready stories, and no agent ran. A story in review holds
// nothing, so the first ready story in pull order is offered; the same story
// in progress still holds every one of them. Its claim stays whole for flai
// check, the trial merge at sync, and the notice at acceptance.
func TestStoryInReviewHoldsNoReadyStory(t *testing.T) {
	open := claimed("S-0220", Review, "flai/internal/harness", "flai/internal/serve", "flai/internal/mcpserver", "flaiover/src/routes")
	under := []string{"flai/internal/harness", "flai/internal/serve", "flai/internal/mcpserver", "flaiover/src/routes", "flai"}
	items := []*Item{open}
	board := &Board{}
	for i := range 10 {
		id := fmt.Sprintf("S-%04d", 230+i)
		items = append(items, claimed(id, Ready, under[i%len(under)]))
		board.Order = append(board.Order, id)
	}
	ready := NewBoardView(items, board, time.Now(), false, nil, NewHolds(items, holdProjects)).ReadyInPullOrder()
	for _, c := range ready {
		if c.Held != nil {
			t.Errorf("%s is held by a story in review: %s", c.ID, c.Held.Reason)
		}
	}
	if first := FirstClear(ready); first == nil || first.ID != "S-0230" {
		t.Errorf("first clear = %+v, want S-0230", first)
	}
	if got := strings.Join(NewHolds(items, holdProjects).Claim(open), ","); got != strings.Join(open.Touches, ",") {
		t.Errorf("the claim of a story in review = %s, want its touches", got)
	}

	open.Status = InProgress
	if first := FirstClear(NewBoardView(items, board, time.Now(), false, nil, NewHolds(items, holdProjects)).ReadyInPullOrder()); first != nil {
		t.Errorf("S-0220 in progress overlaps every ready story, yet %s is clear", first.ID)
	}
}

// i0087Shared is the shared paths the I-0087 holds went through.
var i0087Shared = manifest.Claims{Shared: []string{"docs/users/flai.md", "design/adrs"}}

// I-0087, ADR-0096: S-0283, in progress, held S-0278 through
// docs/users/flai.md, and S-0285 held S-0278 through design/adrs and through
// docs/users, which its tasks narrow to docs/users/flai.md. With those paths
// shared, neither holds it; S-0284 is still held by S-0283 through
// permission.go, which both change and nothing shares. With the list empty,
// each holds as before.
func TestSharedPathsHoldNoReadyStory(t *testing.T) {
	s0278 := claimed("S-0278", Ready, "flai/cmd/stream_sync.go", "flai/internal/issues/issues.go", "design/adrs", "docs/users/flai.md", "docs/users/flai-reference.md")
	s0284 := claimed("S-0284", Ready, "flai/internal/mcpserver/permission.go", "design/adrs", "docs/users/flai.md")
	s0283 := []*Item{
		claimed("S-0283", InProgress, "flai/internal/mcpserver", "docs/users/flai.md"),
		claimTask("T-1001", "S-0283", InProgress, "flai/internal/mcpserver/permission.go", "flai/internal/mcpserver/permission_test.go"),
		claimTask("T-1002", "S-0283", Ready, "docs/users/flai.md"),
	}
	s0285 := []*Item{
		claimed("S-0285", InProgress, "flai/internal/harness", "flai/internal/mcpserver", "flai/internal/guard", "docs/users", "docs/operators", "design/adrs"),
		claimTask("T-1004", "S-0285", Done, "flai/internal/guard/guard.go", "flai/internal/mcpserver/events.go"),
		claimTask("T-1008", "S-0285", InProgress, "docs/users/flai.md", "docs/operators/serve.md", "design/adrs/0095-x.md"),
	}
	cases := []struct {
		name          string
		open          []*Item
		ready         *Item
		empty, inside string // the reason with the list empty, and with it set; "" when not held
	}{
		{"S-0283 through docs/users/flai.md", s0283, s0278,
			"held (overlap): touches docs/users/flai.md, which S-0283 (in progress) touches too; starts when S-0283 moves to review, is cancelled, or is sent back",
			""},
		{"S-0283 through permission.go, which is not shared", s0283, s0284,
			"held (overlap): touches flai/internal/mcpserver/permission.go, which S-0283 (in progress) touches too; starts when S-0283 moves to review, is cancelled, or is sent back",
			"held (overlap): touches flai/internal/mcpserver/permission.go, which S-0283 (in progress) touches too; starts when S-0283 moves to review, is cancelled, or is sent back"},
		{"S-0285 through design/adrs", s0285, s0278,
			"held (overlap): touches design/adrs, which holds design/adrs/0095-x.md that S-0285 (in progress) touches; starts when S-0285 moves to review, is cancelled, or is sent back",
			""},
		{"S-0285 through design/adrs, by narrowing not through mcpserver", s0285, s0284,
			"held (overlap): touches design/adrs, which holds design/adrs/0095-x.md that S-0285 (in progress) touches; starts when S-0285 moves to review, is cancelled, or is sent back",
			""},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			items := append(append([]*Item{}, c.open...), c.ready)
			for _, k := range []struct {
				claims manifest.Claims
				want   string
			}{{manifest.Claims{}, c.empty}, {i0087Shared, c.inside}} {
				h := NewHolds(items, holdProjects).WithShared(k.claims).Of(c.ready)
				switch {
				case k.want == "" && h != nil:
					t.Errorf("shared %v: held: %s", k.claims.Shared, h.Reason)
				case k.want != "" && (h == nil || h.Reason != k.want):
					t.Errorf("shared %v: hold = %+v\nwant reason %s", k.claims.Shared, h, k.want)
				}
			}
		})
	}

	// the board, which wait_for_work reads, offers S-0278 and not S-0284
	items := append(append([]*Item{}, s0283...), s0284, s0278)
	board := &Board{Order: []string{"S-0284", "S-0278"}}
	if first := NewBoardView(items, board, time.Now(), false, nil, NewHolds(items, holdProjects).WithShared(i0087Shared)).FirstPullable(); first == nil || first.ID != "S-0278" {
		t.Errorf("first pullable = %+v, want S-0278", first)
	}
	if first := NewBoardView(items, board, time.Now(), false, nil, NewHolds(items, holdProjects)).FirstPullable(); first != nil {
		t.Errorf("with no shared paths, %s is pullable", first.ID)
	}
}

// ADR-0096: the narrower entry of an overlapping pair decides, and a folder
// only partly inside a shared path still holds.
func TestOverlapsOutsideTheSharedPaths(t *testing.T) {
	h := NewHolds(nil, nil).WithShared(manifest.Claims{Shared: []string{"docs/users/flai.md", "design/adrs", "design/issues/*.md", "../bad"}})
	for _, c := range []struct {
		a, b string
		want bool
	}{
		{"docs/users/flai.md", "docs/users/flai.md", false},
		{"docs/users", "docs/users/flai.md", false},   // the narrower is the shared file
		{"docs", "docs/users/flai.md", false},         // at any depth
		{"docs/users", "docs/users/other.md", true},   // the narrower is not shared
		{"design", "design/adrs", false},              // the narrower is the shared folder
		{"design/adrs/0096.md", "design", false},      // either order
		{"design", "design/issues", true},             // a folder only partly inside *.md
		{"design/issues/summary.md", "design", false}, // a file inside *.md
		{"flai/cmd", "flai/cmd/x.go", true},
		{"flai/cmd", "docs", false}, // no overlap at all
		{"bad", "bad/x.md", true},   // a pattern that is not valid frees nothing
	} {
		if got := h.Overlaps(c.a, c.b); got != c.want {
			t.Errorf("Overlaps(%s, %s) = %v, want %v", c.a, c.b, got, c.want)
		}
	}
}

// ADR-0096: a repository's holds read the shared paths from system-flow.yaml
// as it is now, so an edit made after it was opened, as flai mcp keeps it open
// while flai shared changes the list, frees a story without a restart.
func TestRepoHoldsReadTheSharedPathsAgain(t *testing.T) {
	r := newProject(t)
	items := []*Item{claimed("S-0001", InProgress, "docs/users/flai.md"), claimed("S-0002", Ready, "docs/users/flai.md")}
	if r.Holds(items).Of(items[1]) == nil {
		t.Fatal("with no shared paths, S-0001 should hold S-0002")
	}
	file := filepath.Join(r.Root, manifest.File)
	data, err := os.ReadFile(file)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(file, append(data, "claims:\n  shared: [docs/users/flai.md]\n"...), 0o644); err != nil {
		t.Fatal(err)
	}
	if h := r.Holds(items).Of(items[1]); h != nil {
		t.Errorf("shared after the repo was opened, yet held: %s", h.Reason)
	}
	if err := os.Remove(file); err != nil {
		t.Fatal(err)
	}
	if r.Holds(items).Of(items[1]) == nil {
		t.Error("with the manifest unreadable, the list it was opened with (none) should apply")
	}
}

// waiting is S-0009, ready, naming stories in after:.
func waiting(names ...string) *Item {
	s := claimed("S-0009", Ready, "flai/cmd")
	s.After = names
	return s
}

// S-0130, ADR-0046: a ready story waits for every story it names in after:
// that is not done; a cancelled or missing one keeps the hold and says so.
func TestHoldsByAfter(t *testing.T) {
	done := claimed("S-0002", Done, "docs")
	archived := claimed("S-0003", Done, "docs")
	archived.Archived = true
	cancelled := claimed("S-0004", Cancelled, "docs")
	items := []*Item{
		claimed("S-0001", Backlog, "docs"), done, archived, cancelled,
		claimed("S-0005", Review, "docs"), claimed("S-0006", Ready, "docs"),
	}
	cases := []struct {
		name  string
		ready *Item
		code  string
		why   string
	}{
		{"nothing named", waiting(), "", ""},
		{"done and archived done", waiting("S-2", "S-0003"), "", ""},
		{"the story itself is left to flai check", waiting("S-0009"), "", ""},
		{"one in backlog",
			waiting("S-1"),
			HoldAfter, "held (after): waits for S-0001 (in backlog); starts when S-0001 is done"},
		{"two, one done, repeats once",
			waiting("S-0005", "S-0002", "S-0006", "S-5"),
			HoldAfter, "held (after): waits for S-0005 (in review) and S-0006 (ready); starts when S-0005 and S-0006 are done"},
		{"cancelled keeps the hold",
			waiting("S-0004"),
			HoldAfter, "held (after): waits for S-0004 (cancelled); starts when S-0004 is done; S-0004 was cancelled, so drop it from after: if S-0009 no longer needs it"},
		{"no such story keeps the hold",
			waiting("S-0077"),
			HoldAfter, "held (after): waits for S-0077 (no such story); starts when S-0077 is done; S-0077 names no story, so fix after:"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			h := NewHolds(append(items, c.ready), holdProjects).Of(c.ready)
			switch {
			case c.code == "" && h != nil:
				t.Fatalf("held: %+v", h)
			case c.code != "" && h == nil:
				t.Fatal("not held")
			case h != nil && (h.Code != c.code || h.Reason != c.why):
				t.Errorf("hold = %s: %s\nwant   %s: %s", h.Code, h.Reason, c.code, c.why)
			}
		})
	}
}

// S-0130: after: and an overlap together are held (after), and the reason
// says both, so that "starts when" is true.
func TestHoldsByAfterAndOverlap(t *testing.T) {
	open := claimed("S-0001", InProgress, "flai")
	ready := waiting("S-0001")
	h := NewHolds([]*Item{open, ready}, holdProjects).Of(ready)
	want := "held (after): waits for S-0001 (in progress); starts when S-0001 is done; also held (overlap): touches flai/cmd, inside flai which S-0001 (in progress) touches; starts when S-0001 moves to review, is cancelled, or is sent back"
	if h == nil || h.Code != HoldAfter || h.Reason != want {
		t.Errorf("hold = %+v\nwant reason %s", h, want)
	}
}

// S-0130: a repository's holds find a named story in the archive when only
// the active items were read, as flai serve and moves read them.
func TestRepoHoldsLookInTheArchive(t *testing.T) {
	r := newProject(t)
	mustCreate(t, r, Epic, "E", "")
	archive := func(title, status string) {
		s := mustCreate(t, r, Story, title, "E-0001")
		if err := os.Remove(s.Path); err != nil {
			t.Fatal(err)
		}
		s.Status, s.Archived = status, true
		s.Path = filepath.Join(r.ItemDir(Story, true), filepath.Base(s.Path))
		if err := os.MkdirAll(filepath.Dir(s.Path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := r.Save(s); err != nil {
			t.Fatal(err)
		}
	}
	archive("Done first", Done)
	archive("Given up", Cancelled)
	ready := mustCreate(t, r, Story, "Waits", "E-0001")
	items, err := r.List(false)
	if err != nil {
		t.Fatal(err)
	}
	ready.After = []string{"S-0001"}
	if h := r.Holds(items).Of(ready); h != nil {
		t.Errorf("an archived done story holds: %s", h.Reason)
	}
	if h := NewHolds(items, nil).Of(ready); h == nil || !strings.Contains(h.Reason, "S-0001 (no such story)") {
		t.Errorf("without the archive, S-0001 should read as missing: %+v", h)
	}
	ready.After = []string{"S-0002"}
	if h := r.Holds(items).Of(ready); h == nil || !strings.Contains(h.Reason, "S-0002 (cancelled)") {
		t.Errorf("an archived cancelled story should hold: %+v", h)
	}
}

// writeTree writes an empty file at each path below root, making its folders.
func writeTree(t *testing.T, root string, paths ...string) {
	t.Helper()
	for _, p := range paths {
		full := filepath.Join(root, filepath.FromSlash(p))
		if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(full, nil, 0o644); err != nil {
			t.Fatal(err)
		}
	}
}

// S-0468: a touches entry is a folder touch when it names a folder of the
// tree holding a file, counted at any depth, a component's name or tag read
// as its path; a file, a path absent from the tree, an empty folder, and a
// folder wholly inside a shared pattern are not (ADR-0188).
func TestFolderTouchReadsTheTree(t *testing.T) {
	root := t.TempDir()
	writeTree(t, root, "flai/cmd/a.go", "flai/cmd/sub/b.go", "flai/go.mod", "flai/.git/HEAD", "flai/node_modules/x/y.js",
		"docs/users/flai.md", "docs/users/img/x.png", "design/adrs/0001-a.md", "conf.d/x/y", "Makefile")
	if err := os.MkdirAll(filepath.Join(root, "empty/inner"), 0o755); err != nil {
		t.Fatal(err)
	}
	shared := manifest.Claims{Shared: []string{"design/adrs", "docs/users/*.md", "*"}}
	for _, c := range []struct {
		entry  string
		files  int
		folder bool
	}{
		{"flai/cmd", 2, true},
		{"flai/cmd/", 2, true},
		{"./flai/cmd", 2, true},
		{"flai", 3, true}, // .git and node_modules are not counted
		{"cli", 3, true},  // a component's tag
		{"flaiover", 0, false},
		{"flai/cmd/a.go", 0, false},
		{"flai/cmd/new.go", 0, false},
		{"flai/internal/new", 0, false},
		{"empty", 0, false},
		{"design/adrs", 0, false}, // inside a shared pattern
		{"docs/users", 2, true},   // docs/users/*.md does not hold all of it
		{"conf.d", 1, true},       // * matches conf.d, but not all of it
		{"Makefile", 0, false},
		{"../outside", 0, false},
	} {
		files, folder := FolderTouch(root, holdProjects, shared, c.entry)
		if files != c.files || folder != c.folder {
			t.Errorf("FolderTouch(%q) = %d, %v; want %d, %v", c.entry, files, folder, c.files, c.folder)
		}
	}
}

// S-0468: the refusal of a task's folder touch names the task, the touch, the
// folder a component's name stands for, its file count, and what to do
// instead (ADR-0188).
func TestFolderTouchErrorSaysWhatToDo(t *testing.T) {
	one := (&FolderTouchError{ID: "T-0012", Touch: "flai/cmd", Folder: "flai/cmd", Files: 1}).Error()
	if want := "T-0012's touch flai/cmd is a folder holding 1 file: a task takes no folder touch (ADR-0188); name the files the task changes, or leave the folder to the story's touches with a ### Planning note"; one != want {
		t.Errorf("got  %s\nwant %s", one, want)
	}
	tag := (&FolderTouchError{ID: "T-0012", Touch: "cli", Folder: "flai", Files: 42}).Error()
	if !strings.HasPrefix(tag, "T-0012's touch cli is the folder flai, holding 42 files: ") {
		t.Errorf("a component's tag names its folder: %s", tag)
	}
}

// S-0468: the planner is the writer whose FLAI_ROLE is plan, or else whose
// name starts with planner- (ADR-0188).
func TestIsPlanner(t *testing.T) {
	for _, c := range []struct {
		role, author string
		want         bool
	}{
		{"plan", "", true},
		{"plan", "alex", true},
		{"", "planner-S-0001", true},
		{"", "planner-E-0016", true},
		{"", "agent-S-0468", false},
		{"orchestrate", "alex", false},
		{"", "my-planner-x", false},
		{"", "", false},
	} {
		if got := IsPlanner(c.role, c.author); got != c.want {
			t.Errorf("IsPlanner(%q, %q) = %v, want %v", c.role, c.author, got, c.want)
		}
	}
}

// S-0468: the Planning note is what Notes says under ### Planning, up to the
// next heading of that level; a Planning heading outside Notes is not it.
func TestPlanningNote(t *testing.T) {
	body := "# S-0001 Story\n\n## Goal\n\n### Planning\n\nnot this\n\n## Notes\n\nfirst\n\n### Planning\n\n`flai/cmd` because new commands land there.\n#### Detail\nmore\n\n### Later\n\nnot this either\n"
	if got, want := PlanningNote(body), "`flai/cmd` because new commands land there.\n#### Detail\nmore"; got != want {
		t.Errorf("PlanningNote = %q, want %q", got, want)
	}
	if got := PlanningNote("## Notes\n\nnothing planned\n"); got != "" {
		t.Errorf("no Planning heading: %q", got)
	}
}

// S-0468: the planner's folder touch on a story is refused unless the body's
// Planning note names the folder's path as a whole, backticked or not, or
// the component as given; a file, a path absent from the tree, a shared
// folder, and a folder the story already had are taken (ADR-0188).
func TestCheckStoryTouchesReadsThePlanningNote(t *testing.T) {
	r := newProject(t)
	r.Manifest.Projects = []manifest.Project{{Name: "flai", Path: "flai", Tags: []string{"cli"}}}
	writeTree(t, r.Root, "flai/cmd/a.go", "flai/cmd/b.go", "flai/c/x.go", "design/adrs/0001-a.md")
	if err := os.WriteFile(filepath.Join(r.Root, "system-flow.yaml"), []byte("version: 1\nname: t\nkey: t\nlayout:\n  design: design\n  docs: docs\n  wip: wip\nclaims:\n  shared: [design/adrs]\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	planned := func(note string) string { return "## Goal\n\ng\n\n## Notes\n\n### Planning\n\n" + note + "\n" }
	for _, c := range []struct {
		name    string
		body    string
		before  []string
		touches []string
		refused string // the touch refused, or "" when taken
	}{
		{"no note", "## Notes\n\nnone\n", nil, []string{"flai/cmd/a.go", "flai/cmd"}, "flai/cmd"},
		{"named in backticks", planned("`flai/cmd` because new commands land there"), nil, []string{"flai/cmd"}, ""},
		{"named bare at a full stop", planned("New commands land in flai/cmd."), nil, []string{"flai/cmd/"}, ""},
		{"named with a trailing slash", planned("flai/cmd/ takes them"), nil, []string{"flai/cmd"}, ""},
		{"a longer path does not name it", planned("`flai/cmd` only"), nil, []string{"flai/c"}, "flai/c"},
		{"a file inside does not name it", planned("`flai/cmd/a.go`"), nil, []string{"flai/cmd"}, "flai/cmd"},
		{"named outside Planning", "## Notes\n\n`flai/cmd`\n\n### Planning\n\nother\n", nil, []string{"flai/cmd"}, "flai/cmd"},
		{"a component named as given", planned("the cli, since any package may change"), nil, []string{"cli"}, ""},
		{"a component named by its path", planned("`flai` since any package may change"), nil, []string{"cli"}, ""},
		{"a component not named", planned("`flai/cmd`"), nil, []string{"cli"}, "cli"},
		{"a file, a new path, a shared folder", "## Notes\n", nil, []string{"flai/cmd/a.go", "flai/new", "design/adrs"}, ""},
		{"a folder the story had", "## Notes\n", []string{"flai/cmd/"}, []string{"flai/cmd", "flai/cmd/a.go"}, ""},
	} {
		err := r.CheckStoryTouches("S-0007", c.body, c.before, c.touches)
		var pe *PlanningTouchError
		switch {
		case c.refused == "" && err != nil:
			t.Errorf("%s: taken, got %v", c.name, err)
		case c.refused != "" && (!errors.As(err, &pe) || pe.Touch != c.refused || pe.ID != "S-0007"):
			t.Errorf("%s: %s refused, got %v", c.name, c.refused, err)
		}
	}
}

// S-0468: the refusal of the planner's story folder touch names the story,
// the touch, its folder, its file count, ADR-0188, and the note that admits it.
func TestPlanningTouchErrorSaysWhatToDo(t *testing.T) {
	one := (&PlanningTouchError{ID: "S-0012", Touch: "flai/cmd", Folder: "flai/cmd", Files: 1}).Error()
	if want := "S-0012's touch flai/cmd is a folder holding 1 file: the planner's folder touch on a story needs the story's Planning note to name it (ADR-0188); name the files the story changes, or name flai/cmd under ### Planning in the story's Notes with why files may be added there"; one != want {
		t.Errorf("got  %s\nwant %s", one, want)
	}
	tag := (&PlanningTouchError{ID: "S-0012", Touch: "cli", Folder: "flai", Files: 42}).Error()
	if !strings.HasPrefix(tag, "S-0012's touch cli is the folder flai, holding 42 files: ") {
		t.Errorf("a component's tag names its folder: %s", tag)
	}
}
