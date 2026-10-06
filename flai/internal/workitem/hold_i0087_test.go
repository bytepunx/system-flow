package workitem

import (
	"path/filepath"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/manifest"
)

// I-0087's board of 2026-10-06, rebuilt from its instances. One story in
// progress and one in review held every ready story, and the board ran one
// story at a time under an in-progress limit of three. ADR-0096 frees them by
// three rules: a story in review holds nothing, an overlap inside the
// manifest's shared paths holds nothing, and a story's tasks narrow its folder
// touches in its claim.
//
// S-0285 is in progress, its touches the first guess I-0087 names, folders
// all; its tasks name the files it changes inside them, permission.go as
// S-0283's did. S-0220 waits in review with its folder-wide touches. The ready
// stories, in pull order:
//
//   - S-0284 changes permission.go, which S-0285 names: held under every rule.
//   - S-0278 meets S-0285 only at docs/users/flai.md, a shared path.
//   - S-0290, S-0284 without permission.go, meets both open stories only at
//     other files of flai/internal/mcpserver and flai/internal/harness.
//   - S-0223 meets S-0285 only at another file of docs/users, and writes a
//     new ADR in design/adrs.
//   - S-0291 meets only S-0220, the story in review, at other files of its
//     folders.

const i0087Limit = 3 // the in-progress limit on 2026-10-06

var i0087Order = []string{"S-0284", "S-0278", "S-0290", "S-0223", "S-0291"}

// i0087Old says which of ADR-0096's rules a board is built without.
type i0087Old struct {
	reviewHolds bool // rule 1 reverted: the story in review holds as in progress
	noShared    bool // rule 2 reverted: the shared list is empty
	noNarrowing bool // rule 3 reverted: folder touches stay whole
}

// i0087Board lays out the board with this repository's sub-projects and
// shared paths, read from its system-flow.yaml, the rules in old reverted.
func i0087Board(t *testing.T, old i0087Old) (*Holds, BoardView, map[string]*Item) {
	t.Helper()
	m, err := manifest.Load(filepath.Join("..", "..", "..", manifest.File))
	if err != nil {
		t.Fatalf("read this repository's manifest: %v", err)
	}
	for _, p := range []string{"docs/users/flai.md", "design/adrs"} {
		if _, ok := m.Claims.Covers(p); !ok {
			t.Fatalf("system-flow.yaml's claims.shared does not cover %s, which I-0087's holds went through", p)
		}
	}
	// Rule 3 reverted: ADR-0046 claimed a story's touches and its open
	// tasks', and these tasks name only files inside the story's touches. A
	// task with no touches leaves the folders whole, as that claim did.
	files := func(paths ...string) []string {
		if old.noNarrowing {
			return nil
		}
		return paths
	}
	inProgress := claimed("S-0285", InProgress, "flai/internal/mcpserver", "flai/internal/harness", "docs/users")
	review := claimed("S-0220", Review, "flai/internal/harness", "flai/internal/serve", "flai/internal/mcpserver", "flaiover/src/routes")
	items := []*Item{
		inProgress,
		claimTask("T-1001", "S-0285", Done, files("flai/internal/harness/harness.go")...),
		claimTask("T-1002", "S-0285", InProgress, files("flai/internal/mcpserver/permission.go", "flai/internal/mcpserver/permission_test.go")...),
		claimTask("T-1003", "S-0285", Ready, files("docs/users/flai.md")...),
		review,
		claimed("S-0284", Ready, "flai/internal/mcpserver/permission.go", "design/adrs", "docs/users/flai.md"),
		claimed("S-0278", Ready, "flai/cmd/stream_sync.go", "flai/internal/issues/issues.go", "design/adrs", "docs/users/flai.md", "docs/users/flai-reference.md"),
		claimed("S-0290", Ready, "flai/internal/mcpserver/plan.go", "flai/internal/harness/prompt.go"),
		claimed("S-0223", Ready, "docs/users/flaiover.md", "design/adrs"),
		claimed("S-0291", Ready, "flai/internal/serve/agents.go", "flaiover/src/routes/settings/+page.svelte"),
	}
	h := NewHolds(items, m.Projects)
	// Rule 2 reverted: no shared paths, as before claims.shared.
	if !old.noShared {
		h.WithShared(m.Claims)
	}
	// Rule 1 reverted: ADR-0046 opened a story in review as it opens one in
	// progress, so its claim held every ready story it overlapped.
	if old.reviewHolds {
		h.Open(review, "in review")
	}
	board := &Board{Order: i0087Order, WIPLimits: map[string]int{InProgress: i0087Limit}}
	byID := map[string]*Item{}
	for _, it := range items {
		byID[it.ID] = it
	}
	return h, NewBoardView(items, board, time.Now(), false, nil, h), byID
}

// i0087Held is the ready stories the view marks held, in pull order.
func i0087Held(view BoardView) []string {
	var out []string
	for _, c := range view.ReadyInPullOrder() {
		if c.Held != nil {
			out = append(out, c.ID)
		}
	}
	return out
}

// i0087Pull walks the ready stories in pull order as flai serve's launcher
// does: it skips a held story, starts a clear one while the in-progress limit
// leaves room, and opens each it starts, so that it holds the rest of the walk.
func i0087Pull(h *Holds, view BoardView, byID map[string]*Item) []string {
	free := max(0, view.WIPLimits[InProgress]-view.Counts[InProgress])
	var started []string
	for _, c := range view.ReadyInPullOrder() {
		it := byID[c.ID]
		if len(started) >= free {
			break
		}
		if h.Of(it) != nil {
			continue
		}
		started = append(started, c.ID)
		h.Open(it, agentStartedLabel)
	}
	return started
}

// agentStartedLabel is how the launcher names a story it has started.
const agentStartedLabel = "its agent started"

// I-0087, ADR-0096: with the three rules, only the ready story that changes a
// file the story in progress names is held; the others are offered in pull
// order, two start in the two slots the limit of three leaves, and the rest
// stay clear for the next slot.
func TestI0087ReadyStoriesRunBesideTheOpenOne(t *testing.T) {
	h, view, byID := i0087Board(t, i0087Old{})
	if got := i0087Held(view); !slices.Equal(got, []string{"S-0284"}) {
		t.Fatalf("held = %v, want only S-0284, which changes permission.go", got)
	}
	want := "held (overlap): touches flai/internal/mcpserver/permission.go, which S-0285 (in progress) touches too; starts when S-0285 moves to review, is cancelled, or is sent back"
	if hold := view.ReadyInPullOrder()[0].Held; hold.Reason != want {
		t.Errorf("S-0284: %s\nwant      %s", hold.Reason, want)
	}
	var clear []string
	for _, c := range view.ReadyInPullOrder() {
		if c.Held == nil {
			clear = append(clear, c.ID)
		}
	}
	// at least three could run beside S-0285, given the slots
	if !slices.Equal(clear, []string{"S-0278", "S-0290", "S-0223", "S-0291"}) {
		t.Errorf("clear = %v, want S-0278, S-0290, S-0223 and S-0291", clear)
	}
	if first := view.FirstPullable(); first == nil || first.ID != "S-0278" {
		t.Errorf("first pullable = %+v, want S-0278", first)
	}
	if !view.CanPull() {
		t.Fatalf("no pull: %s", view.PullHold())
	}
	if got := i0087Pull(h, view, byID); !slices.Equal(got, []string{"S-0278", "S-0290"}) {
		t.Errorf("started = %v, want S-0278 and S-0290 in the two free slots", got)
	}
	// What the walk started holds neither of the others: they meet S-0278
	// only at design/adrs and docs/users/flai.md, which are shared.
	for _, id := range []string{"S-0223", "S-0291"} {
		if hold := h.Of(byID[id]); hold != nil {
			t.Errorf("%s is held once S-0278 and S-0290 start: %s", id, hold.Reason)
		}
	}
}

// I-0087: under ADR-0046 as it stood, with the list empty, review holding,
// and no narrowing, every ready story is held, and nothing starts though the
// limit leaves two slots.
func TestI0087OldRulesHoldEveryReadyStory(t *testing.T) {
	h, view, byID := i0087Board(t, i0087Old{reviewHolds: true, noShared: true, noNarrowing: true})
	if got := i0087Held(view); !slices.Equal(got, i0087Order) {
		t.Errorf("held = %v, want every ready story, %v", got, i0087Order)
	}
	for _, c := range view.ReadyInPullOrder() {
		if c.Held != nil && c.Held.Code != HoldOverlap {
			t.Errorf("%s is held (%s), want overlap", c.ID, c.Held.Code)
		}
	}
	if first := view.FirstPullable(); first != nil {
		t.Errorf("%s is pullable", first.ID)
	}
	if got := i0087Pull(h, view, byID); len(got) != 0 {
		t.Errorf("started %v", got)
	}
}

// I-0087, ADR-0096: each rule on its own frees a ready story, so reverting
// any one holds that story again.
func TestI0087EachRuleFreesAStory(t *testing.T) {
	cases := []struct {
		name  string
		old   i0087Old
		held  []string // the ready stories held with the rule reverted
		freed string   // the one of them only this rule frees
		by    string   // what its hold names with the rule reverted
	}{
		{"a story in review holds nothing", i0087Old{reviewHolds: true},
			[]string{"S-0284", "S-0290", "S-0291"}, "S-0291", "S-0220 (in review)"},
		{"an overlap inside a shared path holds nothing", i0087Old{noShared: true},
			[]string{"S-0284", "S-0278"}, "S-0278", "touches docs/users/flai.md, which S-0285 (in progress) touches too"},
		{"tasks narrow a folder touch", i0087Old{noNarrowing: true},
			[]string{"S-0284", "S-0290", "S-0223"}, "S-0223", "touches docs/users/flaiover.md, inside docs/users which S-0285 (in progress) touches"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			_, view, _ := i0087Board(t, c.old)
			if got := i0087Held(view); !slices.Equal(got, c.held) {
				t.Errorf("held with the rule reverted = %v, want %v", got, c.held)
			}
			for _, card := range view.ReadyInPullOrder() {
				if card.ID == c.freed && (card.Held == nil || !strings.Contains(card.Held.Reason, c.by)) {
					t.Errorf("%s: hold = %+v, want one naming %q", c.freed, card.Held, c.by)
				}
			}
			_, view, _ = i0087Board(t, i0087Old{})
			if got := i0087Held(view); slices.Contains(got, c.freed) {
				t.Errorf("with the rule, %s is still held", c.freed)
			}
		})
	}
}
