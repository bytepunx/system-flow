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

func ids(cards []BoardCard) []string {
	out := make([]string, len(cards))
	for i, c := range cards {
		out[i] = c.ID
	}
	return out
}
