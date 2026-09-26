package workitem

import (
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
// name or tag read as its path, without duplicates.
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
	} {
		if got := PathsOverlap(c.a, c.b); got != c.want {
			t.Errorf("PathsOverlap(%q, %q) = %v", c.a, c.b, got)
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
			HoldOverlap, "held (overlap): touches flai/cmd/serve, inside flai/cmd which S-0001 (in progress) touches; starts when S-0001 is accepted, cancelled, or sent back"},
		{"the same path, in review",
			[]*Item{claimed("S-0001", Review, "docs")}, claimed("S-0009", Ready, "docs/"),
			HoldOverlap, "held (overlap): touches docs, which S-0001 (in review) touches too; starts when S-0001 is accepted, cancelled, or sent back"},
		{"around an open story's path",
			[]*Item{claimed("S-0001", InProgress, "flai/cmd")}, claimed("S-0009", Ready, "cli"),
			HoldOverlap, "held (overlap): touches flai, which holds flai/cmd that S-0001 (in progress) touches; starts when S-0001 is accepted, cancelled, or sent back"},
		{"a sibling path is not held",
			[]*Item{claimed("S-0001", InProgress, "flai")}, claimed("S-0009", Ready, "flaiover", "docs"),
			"", ""},
		{"backlog, done, and ready stories hold nothing",
			[]*Item{claimed("S-0001", Backlog, "flai"), claimed("S-0002", Done, "flai"), claimed("S-0003", Ready, "flai")}, claimed("S-0009", Ready, "flai"),
			"", ""},
		{"an open task widens the open story's claim",
			[]*Item{claimed("S-0001", InProgress, "docs"), claimTask("T-0001", "S-0001", InProgress, "dashboard")}, claimed("S-0009", Ready, "flaiover/src"),
			HoldOverlap, "held (overlap): touches flaiover/src, inside flaiover which S-0001 (in progress) touches; starts when S-0001 is accepted, cancelled, or sent back"},
		{"held by two",
			[]*Item{claimed("S-0002", Review, "docs"), claimed("S-0001", InProgress, "flai")}, claimed("S-0009", Ready, "flai/cmd", "docs/users"),
			HoldOverlap, "held (overlap): touches flai/cmd, inside flai which S-0001 (in progress) touches; touches docs/users, inside docs which S-0002 (in review) touches; starts when S-0001 and S-0002 are accepted, cancelled, or sent back"},
		{"a ready story with no touches",
			[]*Item{claimed("S-0001", InProgress, "flai"), claimed("S-0002", Review, "docs")}, claimed("S-0009", Ready),
			HoldNoTouches, "held (no-touches): declares no touches, so it may change what S-0001 (in progress) and S-0002 (in review) change; starts when it declares touches that overlap no open story's, or when S-0001 and S-0002 are accepted, cancelled, or sent back"},
		{"an open story with no touches",
			[]*Item{claimed("S-0001", InProgress)}, claimed("S-0009", Ready, "flai"),
			HoldNoTouches, "held (no-touches): S-0001 (in progress) declares no touches, so it may change anything; starts when S-0001 is accepted, cancelled, or sent back"},
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
	v := NewBoardView(items, board, time.Now(), false, nil, holdProjects)
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
	if first := NewBoardView(items, board, time.Now(), false, nil, holdProjects).FirstPullable(); first != nil {
		t.Errorf("every ready story is held, yet %s is pullable", first.ID)
	}
}
