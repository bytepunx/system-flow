package workitem

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

// shareAt is when the shares in these tests were made.
const shareAt = "2026-10-08T10:00:00Z"

func share(holder, held string, paths ...string) Share {
	return Share{Holder: holder, Held: held, Paths: paths, Split: "split", By: "agent", At: shareAt}
}

// ADR-0134: an overlap whose narrower path lies inside a path the holding
// story shared with the held one does not hold; one partly shared names only
// what still holds; a share with one story clears nothing held by another.
func TestSharesClearHolds(t *testing.T) {
	cases := []struct {
		name   string
		open   []*Item
		ready  *Item
		shares []Share
		why    string // empty: not held
	}{
		{"no share holds",
			[]*Item{claimed("S-0001", InProgress, "flai/cmd")}, claimed("S-0009", Ready, "flai/cmd/serve"), nil,
			"held (overlap): touches flai/cmd/serve, inside flai/cmd which S-0001 (in progress) touches; starts when S-0001 moves to review, is cancelled, or is sent back"},
		{"the overlap shared whole",
			[]*Item{claimed("S-0001", InProgress, "flai/cmd")}, claimed("S-0009", Ready, "flai/cmd/serve"),
			[]Share{share("S-1", "s-9", "flai/cmd/serve")}, ""},
		{"a folder around the overlap shared",
			[]*Item{claimed("S-0001", InProgress, "flai/cmd")}, claimed("S-0009", Ready, "flai/cmd/serve"),
			[]Share{share("S-0001", "S-0009", "flai/cmd/")}, ""},
		{"a path inside the narrower shares too little",
			[]*Item{claimed("S-0001", InProgress, "flai/cmd")}, claimed("S-0009", Ready, "flai/cmd/serve"),
			[]Share{share("S-0001", "S-0009", "flai/cmd/serve/x.go")},
			"held (overlap): touches flai/cmd/serve, inside flai/cmd which S-0001 (in progress) touches; starts when S-0001 moves to review, is cancelled, or is sent back"},
		{"partly shared names the path still held",
			[]*Item{claimed("S-0001", InProgress, "flai/cmd", "docs")}, claimed("S-0009", Ready, "flai/cmd/serve", "docs/users"),
			[]Share{share("S-0001", "S-0009", "flai/cmd")},
			"held (overlap): touches docs/users, inside docs which S-0001 (in progress) touches; starts when S-0001 moves to review, is cancelled, or is sent back"},
		{"a holder whose overlap is all shared is not named",
			[]*Item{claimed("S-0001", InProgress, "flai"), claimed("S-0002", InProgress, "docs")}, claimed("S-0009", Ready, "flai/cmd", "docs/users"),
			[]Share{share("S-0001", "S-0009", "flai/cmd")},
			"held (overlap): touches docs/users, inside docs which S-0002 (in progress) touches; starts when S-0002 moves to review, is cancelled, or is sent back"},
		{"a share by another holder clears nothing",
			[]*Item{claimed("S-0001", InProgress, "flai"), claimed("S-0002", InProgress, "docs")}, claimed("S-0009", Ready, "flai/cmd"),
			[]Share{share("S-0002", "S-0009", "flai/cmd")},
			"held (overlap): touches flai/cmd, inside flai which S-0001 (in progress) touches; starts when S-0001 moves to review, is cancelled, or is sent back"},
		{"a share with a third story clears nothing",
			[]*Item{claimed("S-0001", InProgress, "flai"), claimed("S-0008", Ready, "flai/cmd")}, claimed("S-0009", Ready, "flai/cmd"),
			[]Share{share("S-0001", "S-0008", "flai/cmd")},
			"held (overlap): touches flai/cmd, inside flai which S-0001 (in progress) touches; starts when S-0001 moves to review, is cancelled, or is sent back"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			h := NewHolds(append(c.open, c.ready), holdProjects).WithShares(c.shares).Of(c.ready)
			switch {
			case c.why == "" && h != nil:
				t.Fatalf("held: %s", h.Reason)
			case c.why == "":
			case h == nil:
				t.Fatal("not held")
			case h.Code != HoldOverlap || h.Reason != c.why:
				t.Errorf("hold = %s: %s\nwant   %s: %s", h.Code, h.Reason, HoldOverlap, c.why)
			}
		})
	}
}

// ADR-0134: a share is in force while both stories are ready, in progress, or
// in review and neither has moved to backlog, done, or cancelled since it was
// made.
func TestSharesInForce(t *testing.T) {
	moved := func(s *Item, to, at string) *Item {
		s.Transitions = append(s.Transitions, Transition{To: to, At: at, By: "test"})
		return s
	}
	archived := claimed("S-0009", Ready, "flai/cmd")
	archived.Archived = true
	cases := []struct {
		name   string
		holder *Item
		held   *Item
		at     string
		want   bool
	}{
		{"both open", claimed("S-0001", InProgress, "flai"), claimed("S-0009", Ready, "flai/cmd"), shareAt, true},
		{"the holder in review", claimed("S-0001", Review, "flai"), claimed("S-0009", Ready, "flai/cmd"), shareAt, true},
		{"back in backlog before the share",
			claimed("S-0001", InProgress, "flai"), moved(claimed("S-0009", Ready, "flai/cmd"), Backlog, "2026-10-08T09:00:00Z"), shareAt, true},
		{"the held story back in backlog after the share",
			claimed("S-0001", InProgress, "flai"), moved(claimed("S-0009", Ready, "flai/cmd"), Backlog, "2026-10-08T11:00:00Z"), shareAt, false},
		{"moved to backlog the moment it was made",
			claimed("S-0001", InProgress, "flai"), moved(claimed("S-0009", Ready, "flai/cmd"), Backlog, shareAt), shareAt, false},
		{"the held story in backlog", claimed("S-0001", InProgress, "flai"), claimed("S-0009", Backlog, "flai/cmd"), shareAt, false},
		{"the holder done", moved(claimed("S-0001", Done, "flai"), Done, "2026-10-08T11:00:00Z"), claimed("S-0009", Ready, "flai/cmd"), shareAt, false},
		{"the holder cancelled and sent back after the share",
			moved(claimed("S-0001", InProgress, "flai"), Cancelled, "2026-10-08T11:00:00Z"), claimed("S-0009", Ready, "flai/cmd"), shareAt, false},
		{"the held story archived", claimed("S-0001", InProgress, "flai"), archived, shareAt, false},
		{"a time that does not parse", claimed("S-0001", InProgress, "flai"), claimed("S-0009", Ready, "flai/cmd"), "today", false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			s := share("S-0001", "S-0009", "flai/cmd")
			s.At = c.at
			h := NewHolds([]*Item{c.holder, c.held}, holdProjects).WithShares([]Share{s})
			if got := len(h.shares) == 1; got != c.want {
				t.Errorf("in force = %v, want %v", got, c.want)
			}
		})
	}
	if h := NewHolds([]*Item{claimed("S-0001", InProgress, "flai")}, holdProjects).WithShares([]Share{share("S-0001", "S-0009", "flai")}); len(h.shares) != 0 {
		t.Error("a share with a story that does not exist is in force")
	}
}

// writeConversation stores a conversation with the given front matter under
// r's wip/messages.
func writeConversation(t *testing.T, r *Repo, name, frontMatter string) {
	t.Helper()
	dir := filepath.Join(r.WipDir(), messagesFolder)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, name), []byte("---\n"+frontMatter+"---\n\n# Body\n"), 0o644); err != nil {
		t.Fatal(err)
	}
}

// ADR-0134: a repository's holds read the shares of the conversations stored
// open, skip a closed one and one that does not parse, and hold as before
// with no messages folder.
func TestRepoHoldsReadTheShares(t *testing.T) {
	r := newProject(t)
	items := []*Item{claimed("S-0001", InProgress, "flai/cmd"), claimed("S-0009", Ready, "flai/cmd/serve")}
	if got := r.Shares(); got != nil {
		t.Errorf("shares with no messages folder = %+v", got)
	}
	if r.Holds(items).Of(items[1]) == nil {
		t.Fatal("with no share, S-0001 should hold S-0009")
	}
	shares := "shares:\n  - holder: S-0001\n    held: S-0009\n    paths: [flai/cmd/serve]\n    split: \"S-0009 adds serve.go; S-0001 keeps root.go\"\n    by: agent\n    at: " + shareAt + "\n"
	writeConversation(t, r, "MS-0002-closed.md", "id: MS-0002\nstatus: closed\n"+shares)
	writeConversation(t, r, "MS-0003-bad.md", "id: MS-0003\nstatus: open\nshares: [: :\n")
	writeConversation(t, r, "notes.md", "status: open\n"+shares)
	if got := r.Shares(); got != nil {
		t.Errorf("shares of closed, unparsable, and stray files = %+v", got)
	}
	if r.Holds(items).Of(items[1]) == nil {
		t.Fatal("a closed conversation's share clears the hold")
	}
	writeConversation(t, r, "MS-0001-open.md", "id: MS-0001\nfrom: S-0009\nto: S-0001\nstatus: open\nowner: x\n"+shares)
	want := []Share{{Holder: "S-0001", Held: "S-0009", Paths: []string{"flai/cmd/serve"}, Split: "S-0009 adds serve.go; S-0001 keeps root.go", By: "agent", At: shareAt}}
	if got := r.Shares(); !reflect.DeepEqual(got, want) {
		t.Errorf("shares = %+v\nwant     %+v", got, want)
	}
	if h := r.Holds(items).Of(items[1]); h != nil {
		t.Errorf("shared in an open conversation, yet held: %s", h.Reason)
	}
}

// ADR-0134: flai serve asks each open story that holds a ready story on
// overlap alone, about the narrower path of each pair still held.
func TestOverlapsBy(t *testing.T) {
	after := claimed("S-0009", Ready, "flai/cmd")
	after.After = []string{"S-0001"}
	cases := []struct {
		name   string
		open   []*Item
		ready  *Item
		shares []Share
		want   []Overlap
	}{
		{"not held", []*Item{claimed("S-0001", InProgress, "docs")}, claimed("S-0009", Ready, "flai"), nil, nil},
		{"one holder, the narrower of each pair",
			[]*Item{claimed("S-0001", InProgress, "flai/cmd", "flai/internal/x.go", "docs")}, claimed("S-0009", Ready, "flai", "docs/users"), nil,
			[]Overlap{{By: "S-0001", Paths: []string{"flai/cmd", "flai/internal/x.go", "docs/users"}}}},
		{"each holder in ID order, a path once",
			[]*Item{claimed("S-0002", InProgress, "flai/cmd"), claimed("S-0001", InProgress, "flai", "flai/cmd")}, claimed("S-0009", Ready, "flai/cmd"), nil,
			[]Overlap{{By: "S-0001", Paths: []string{"flai/cmd"}}, {By: "S-0002", Paths: []string{"flai/cmd"}}}},
		{"a share clears part",
			[]*Item{claimed("S-0001", InProgress, "flai/cmd", "docs")}, claimed("S-0009", Ready, "flai/cmd/serve", "docs/users"),
			[]Share{share("S-0001", "S-0009", "flai/cmd")},
			[]Overlap{{By: "S-0001", Paths: []string{"docs/users"}}}},
		{"a holder all shared is left out",
			[]*Item{claimed("S-0001", InProgress, "flai"), claimed("S-0002", InProgress, "docs")}, claimed("S-0009", Ready, "flai/cmd", "docs/users"),
			[]Share{share("S-0001", "S-0009", "flai/cmd")},
			[]Overlap{{By: "S-0002", Paths: []string{"docs/users"}}}},
		{"held after", []*Item{claimed("S-0001", InProgress, "flai")}, after, nil, nil},
		{"the ready story has no touches", []*Item{claimed("S-0001", InProgress, "flai")}, claimed("S-0009", Ready), nil, nil},
		{"an open story has no touches",
			[]*Item{claimed("S-0001", InProgress, "flai"), claimed("S-0002", InProgress)}, claimed("S-0009", Ready, "flai/cmd"), nil, nil},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := NewHolds(append(c.open, c.ready), holdProjects).WithShares(c.shares).OverlapsBy(c.ready)
			if !reflect.DeepEqual(got, c.want) {
				t.Errorf("overlaps = %+v\nwant       %+v", got, c.want)
			}
		})
	}

	// a story whose agent the launcher started holds as one in progress does
	a, b := claimed("S-0001", Ready, "flai"), claimed("S-0002", Ready, "flai/cmd")
	h := NewHolds([]*Item{a, b}, holdProjects)
	h.Open(a, "its agent started")
	if got, want := h.OverlapsBy(b), []Overlap{{By: "S-0001", Paths: []string{"flai/cmd"}}}; !reflect.DeepEqual(got, want) {
		t.Errorf("overlaps = %+v, want %+v", got, want)
	}
}
