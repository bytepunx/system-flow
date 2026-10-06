package workitem

import (
	"os"
	"reflect"
	"sort"
	"strings"
	"testing"
	"time"
)

// S-0219: board.md records who placed a story by hand and when, beside the
// order, and reads it back.
func TestPlacedRoundTrips(t *testing.T) {
	r := newProject(t)
	b, err := r.LoadBoard()
	if err != nil {
		t.Fatal(err)
	}
	b.Order = []string{"S-0010", "S-0002"}
	if err := b.Save("2026-09-15"); err != nil {
		t.Fatal(err)
	}
	data, _ := os.ReadFile(b.Path)
	if strings.Contains(string(data), "placed") {
		t.Errorf("a board with no placement writes no placed map:\n%s", data)
	}

	b.RecordPlacement("S-0010", "alex", t0)
	b.RecordPlacement("S-0002", "the operator: alex", t0.Add(time.Hour))
	b.RecordPlacement("S-0010", "flaiover", t0.Add(2*time.Hour))
	if err := b.Save("2026-09-15"); err != nil {
		t.Fatal(err)
	}
	data, _ = os.ReadFile(b.Path)
	want := "placed:\n  S-0002:\n    by: \"the operator: alex\"\n    at: 2026-09-15T21:00:00Z\n  S-0010:\n    by: flaiover\n    at: 2026-09-15T22:00:00Z\n---\n"
	if !strings.Contains(string(data), "order:\n  - S-0010\n  - S-0002\n"+want) {
		t.Errorf("placed follows order, by ID, the latest placement of each:\n%s", data)
	}
	back, err := r.LoadBoard()
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(back.Placed, b.Placed) {
		t.Errorf("read back: got %+v want %+v", back.Placed, b.Placed)
	}
}

// S-0219: a placement holds while its story stays in its column; flai move
// out of ready or backlog drops it.
func TestMoveOutOfColumnDropsPlacement(t *testing.T) {
	r := newProject(t)
	mustCreate(t, r, Epic, "E", "")
	var stories []*Item
	for _, title := range []string{"One", "Two", "Three", "Four"} {
		s := mustCreate(t, r, Story, title, "E-0001")
		s.Body = strings.Replace(s.Body, "## Acceptance criteria\n- [ ]\n", "## Acceptance criteria\n- [ ] works\n", 1)
		if err := r.Save(s); err != nil {
			t.Fatal(err)
		}
		stories = append(stories, s)
	}
	mustMove(t, r, stories[0], Ready, "")
	mustMove(t, r, stories[1], Ready, "")
	b, _ := r.LoadBoard()
	for _, s := range stories {
		b.RecordPlacement(s.ID, "alex", t0)
	}
	if err := b.Save("2026-09-15"); err != nil {
		t.Fatal(err)
	}

	mustMove(t, r, stories[0], InProgress, "")
	mustMove(t, r, stories[1], Backlog, "")
	mustMove(t, r, stories[2], Cancelled, "not needed")
	b, _ = r.LoadBoard()
	if got := keys(b.Placed); !reflect.DeepEqual(got, []string{"S-0004"}) {
		t.Errorf("out of ready to in-progress or backlog, and out of backlog to cancelled, drop the placement: left %v", got)
	}
}

func TestHandPlaced(t *testing.T) {
	b := &Board{}
	now := t0.Add(48 * time.Hour)
	b.RecordPlacement("S-0001", "alex", now.Add(-time.Hour))
	b.RecordPlacement("S-0002", "alex", now.Add(-48*time.Hour))
	b.RecordPlacement("S-0003", ActivityOrchestrator, now.Add(-time.Hour))
	b.RecordPlacement("S-0004", "flaiover", now.Add(-DefaultKeepPlaced))
	b.Placed["S-0005"] = Placed{By: "alex", At: "yesterday"}
	for _, c := range []struct {
		id   string
		keep time.Duration
		want bool
	}{
		{"S-0001", DefaultKeepPlaced, true},
		{"S-0001", 30 * time.Minute, false},
		{"S-0001", 0, false},
		{"S-0002", DefaultKeepPlaced, false},
		{"S-0003", DefaultKeepPlaced, false},
		{"S-0004", DefaultKeepPlaced, true},
		{"S-0005", DefaultKeepPlaced, false},
		{"S-0006", DefaultKeepPlaced, false},
	} {
		p, ok := b.HandPlaced(c.id, c.keep, now)
		if ok != c.want {
			t.Errorf("%s within %s: got %v want %v", c.id, c.keep, ok, c.want)
		}
		if ok && p != b.Placed[c.id] {
			t.Errorf("%s: got %+v want %+v", c.id, p, b.Placed[c.id])
		}
	}
}

func keys(m map[string]Placed) []string {
	out := []string{}
	for k := range m {
		out = append(out, k)
	}
	sort.Slice(out, func(i, j int) bool { return lessID(out[i], out[j]) })
	return out
}
