package workitem

import (
	"encoding/json"
	"strings"
	"testing"
	"time"
)

// S-0087, found in T-0318's end-to-end review: acceptance archives a story
// the moment it merges, and the done column used to drop every archived
// item without exception, so a story could never be seen there before it
// was published, which is exactly the state the operator asked to be able
// to see. A done, archived story named in pendingPublish stays on the done
// column; one not named there, or archived under any other status, does not.
func TestArchivedDoneStoryStaysUntilPublished(t *testing.T) {
	items := []*Item{
		{ID: "S-0001", Type: Story, Title: "Waiting", Status: Done, Archived: true},
		{ID: "S-0002", Type: Story, Title: "Published", Status: Done, Archived: true},
		{ID: "S-0003", Type: Story, Title: "Cancelled but named anyway", Status: Cancelled, Archived: true},
		{ID: "S-0004", Type: Story, Title: "Not done", Status: InProgress},
	}
	board := &Board{}
	pending := map[string]bool{"S-0001": true, "S-0003": true}
	now := time.Date(2026, 9, 21, 12, 0, 0, 0, time.UTC)

	v := NewBoardView(items, board, now, false, pending, nil)
	got := ids(v.Columns[Done])
	if len(got) != 1 || got[0] != "S-0001" {
		t.Errorf("done column: %v", got)
	}
	if !v.Columns[Done][0].Archived {
		t.Error("the card says it is archived, so the dashboard can tell it is there only while pending (S-0174)")
	}
	if len(v.Columns[Cancelled]) != 0 {
		t.Error("pendingPublish names a cancelled item too, but only a done one may stay")
	}

	// nothing named as pending: the old behaviour, nothing archived shows
	if v := NewBoardView(items, board, now, false, nil, nil); len(v.Columns[Done]) != 0 || v.Unpublished != nil {
		t.Errorf("no pending map: %v, unpublished %v", ids(v.Columns[Done]), v.Unpublished)
	}
}

// S-0195, ADR-0067: the view lists what is accepted and not yet published,
// every ID pendingPublish names, in order, for whoever publishes.
func TestBoardViewListsWhatIsUnpublished(t *testing.T) {
	items := []*Item{{ID: "S-0002", Type: Story, Title: "Accepted", Status: Done, Archived: true}}
	pending := map[string]bool{"S-0009": true, "S-0002": true, "S-0005": false}
	v := NewBoardView(items, &Board{}, time.Now(), false, pending, nil)
	if got := v.Unpublished; len(got) != 2 || got[0] != "S-0002" || got[1] != "S-0009" {
		t.Errorf("unpublished: %v", got)
	}
}

// S-0201: a story an agent wrote and the operator has not finalized carries
// draft on its card, so the dashboard can mark it; any other card leaves the
// key out, and only a story is ever a draft.
func TestADraftStorysCardSaysSo(t *testing.T) {
	items := []*Item{
		{ID: "S-0001", Type: Story, Title: "Drafted", Status: Backlog, Draft: true},
		{ID: "S-0002", Type: Story, Title: "Finalized", Status: Backlog},
		{ID: "E-0001", Type: Epic, Title: "Not a story", Status: Backlog, Draft: true},
	}
	v := NewBoardView(items, &Board{}, time.Now(), true, nil, nil)
	cards := map[string]BoardCard{}
	for _, c := range v.Columns[Backlog] {
		cards[c.ID] = c
	}
	for _, tc := range []struct {
		id    string
		draft bool
	}{
		{"S-0001", true},
		{"S-0002", false},
		{"E-0001", false},
	} {
		c, ok := cards[tc.id]
		if !ok {
			t.Fatalf("%s: no card in %v", tc.id, ids(v.Columns[Backlog]))
		}
		if c.Draft != tc.draft {
			t.Errorf("%s: draft %v, want %v", tc.id, c.Draft, tc.draft)
		}
		data, err := json.Marshal(c)
		if err != nil {
			t.Fatalf("%s: marshal: %v", tc.id, err)
		}
		if got := strings.Contains(string(data), `"draft":true`); got != tc.draft {
			t.Errorf("%s: JSON %s", tc.id, data)
		}
		if !tc.draft && strings.Contains(string(data), `"draft"`) {
			t.Errorf("%s: a card that is not a draft has a draft key: %s", tc.id, data)
		}
	}
}

// S-0338: flai board --json gives a ready card held on overlap alone the
// conversation in which flai asked its holding story's agent about the hold,
// and leaves asked out of a held card no one was asked about and of one asked
// only before it last entered ready, or in a closed conversation; it gives a
// card started on a share the shared paths, the split, and the conversation,
// and none once the share has ended.
func TestBoardCardsNameTheAskAndTheShare(t *testing.T) {
	r := newProject(t)
	readyAt := func(id, at string, touches ...string) *Item {
		s := claimed(id, Ready, touches...)
		s.Transitions = []Transition{{To: Ready, At: at, By: "test"}}
		return s
	}
	items := []*Item{
		claimed("S-0001", InProgress, "flai/cmd"),
		claimed("S-0002", InProgress, "docs"),
		claimed("S-0003", InProgress, "flaiover/src/lib", "flaiover/src/routes"),
		readyAt("S-0007", "2026-10-08T09:00:00Z", "flai/cmd/serve"),        // asked
		readyAt("S-0008", "2026-10-08T09:00:00Z", "docs/users"),            // asked before it last entered ready
		readyAt("S-0009", "2026-10-08T09:00:00Z", "flai/cmd/board.go"),     // not asked
		readyAt("S-0006", "2026-10-08T09:00:00Z", "flaiover/src/lib/x.ts"), // asked in a closed conversation
		claimed("S-0005", InProgress, "flaiover/src/routes"),               // started on a share
	}
	writeConversationBody(t, r, "MS-0001-ask.md", "id: MS-0001\nfrom: S-0007\nto: S-0001\nstatus: open\n",
		askEntry("2026-10-08T08:00:00Z", "S-0007")+askEntry("2026-10-08T10:00:00Z", "S-0007"))
	writeConversationBody(t, r, "MS-0002-ask.md", "id: MS-0002\nfrom: S-0002\nto: S-0008\nstatus: open\n",
		askEntry("2026-10-08T08:00:00Z", "S-0008"))
	writeConversationBody(t, r, "MS-0003-closed.md", "id: MS-0003\nfrom: S-0006\nto: S-0003\nstatus: closed\n",
		askEntry("2026-10-08T10:00:00Z", "S-0006"))
	shares := "shares:\n  - holder: S-0003\n    held: S-0005\n    paths: [flaiover/src/routes]\n    split: \"S-0005 adds the page; S-0003 keeps the layout\"\n    by: agent\n    at: " + shareAt + "\n"
	writeConversationBody(t, r, "MS-0004-share.md", "id: MS-0004\nfrom: S-0005\nto: S-0003\nstatus: open\n"+shares,
		askEntry("2026-10-08T09:30:00Z", "S-0005"))

	cards := func() map[string]BoardCard {
		v := NewBoardView(items, &Board{}, time.Now(), false, nil, r.Holds(items))
		out := map[string]BoardCard{}
		for _, col := range v.Columns {
			for _, c := range col {
				out[c.ID] = c
			}
		}
		return out
	}
	asJSON := func(v any) string {
		data, err := json.Marshal(v)
		if err != nil {
			t.Fatal(err)
		}
		return string(data)
	}
	got := cards()
	for id, want := range map[string]string{
		"S-0007": `[{"by":"S-0001","conversation":"MS-0001","at":"2026-10-08T10:00:00Z"}]`,
		"S-0008": "null",
		"S-0009": "null",
		"S-0006": "null",
	} {
		c := got[id]
		if c.Held == nil || c.Held.Code != HoldOverlap {
			t.Fatalf("%s: held = %+v", id, c.Held)
		}
		if a := asJSON(c.Held.Asked); a != want {
			t.Errorf("%s: asked = %s\nwant     %s", id, a, want)
		}
		if want == "null" && strings.Contains(asJSON(c), `"asked"`) {
			t.Errorf("%s: a card no one was asked about since it entered ready has an asked key: %s", id, asJSON(c))
		}
	}
	want := `[{"holder":"S-0003","held":"S-0005","paths":["flaiover/src/routes"],"split":"S-0005 adds the page; S-0003 keeps the layout","by":"agent","at":"` + shareAt + `","conversation":"MS-0004"}]`
	if s := asJSON(got["S-0005"].Shared); s != want {
		t.Errorf("shared = %s\nwant     %s", s, want)
	}
	for id, c := range got {
		if id != "S-0005" && strings.Contains(asJSON(c), `"shared"`) {
			t.Errorf("%s: shared on a card no share was made for: %s", id, asJSON(c))
		}
	}

	// the share ends when the story it was made for goes back to backlog
	// after it, and the card returned to ready names no share
	items[7].Status = Ready
	items[7].Transitions = []Transition{
		{To: Backlog, At: "2026-10-08T11:00:00Z", By: "test"},
		{To: Ready, At: "2026-10-08T11:30:00Z", By: "test"},
	}
	if c := cards()["S-0005"]; c.Shared != nil || strings.Contains(asJSON(c), `"shared"`) {
		t.Errorf("a share that ended is shown: %s", asJSON(c))
	}
}

func ids(cards []BoardCard) []string {
	out := make([]string, len(cards))
	for i, c := range cards {
		out[i] = c.ID
	}
	return out
}
