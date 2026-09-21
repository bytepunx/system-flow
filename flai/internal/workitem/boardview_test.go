package workitem

import (
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

	v := NewBoardView(items, board, now, false, pending)
	got := ids(v.Columns[Done])
	if len(got) != 1 || got[0] != "S-0001" {
		t.Errorf("done column: %v", got)
	}
	if len(v.Columns[Cancelled]) != 0 {
		t.Error("pendingPublish names a cancelled item too, but only a done one may stay")
	}

	// nothing named as pending: the old behaviour, nothing archived shows
	if v := NewBoardView(items, board, now, false, nil); len(v.Columns[Done]) != 0 {
		t.Errorf("no pending map: %v", ids(v.Columns[Done]))
	}
}

func ids(cards []BoardCard) []string {
	out := make([]string, len(cards))
	for i, c := range cards {
		out[i] = c.ID
	}
	return out
}
