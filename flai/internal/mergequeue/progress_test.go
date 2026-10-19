package mergequeue

import (
	"testing"
	"time"
)

// S-0400: the landing's next story is the first of its stories not landed,
// however each is written, and none once every one has.
func TestLandingNextIsTheFirstStoryNotLanded(t *testing.T) {
	var none *Landing
	if s, ok := none.Next(); ok || s != "" {
		t.Errorf("no landing: next %q, %v", s, ok)
	}
	l := &Landing{Stories: []string{"S-0001", "S-0002", "S-0003"}, Landed: []string{"S-1"}}
	if s, ok := l.Next(); !ok || s != "S-0002" {
		t.Errorf("next %q, %v; want S-0002", s, ok)
	}
	l.Landed = append(l.Landed, "S-0002", "S-0003")
	if s, ok := l.Next(); ok {
		t.Errorf("next %q with every story landed", s)
	}
}

// S-0400: a batch answers a story's heads, however it is written, and none
// for a story it kept none for.
func TestBatchHeadOf(t *testing.T) {
	b := &Batch{Heads: []Head{{Story: "S-0001", Verified: "v1", Rebased: "r1"}}}
	if h, ok := b.HeadOf("S-1"); !ok || h.Rebased != "r1" || h.Verified != "v1" {
		t.Errorf("head %+v, %v", h, ok)
	}
	if _, ok := b.HeadOf("S-0002"); ok {
		t.Error("a head for a story not in the batch")
	}
}

// S-0400: a batch's heads, bisection, and landing are kept through Update
// and read back as they were written.
func TestABatchsProgressIsKeptInTheState(t *testing.T) {
	dir := t.TempDir()
	at := time.Date(2026, 10, 19, 22, 0, 0, 0, time.UTC)
	bisect := &Bisect{
		Steps: []BisectStep{{Stories: []string{"S-0001"}, Commit: "r1", StartedAt: at, Result: &Result{Passed: true, FinishedAt: at}}},
		Good:  1, Bad: 2, Failing: "S-0002", Commit: "r2",
		Result: &Result{StoppedAt: "unit", Summary: "unit failed", FinishedAt: at},
	}
	if _, err := Update(dir, func(s *State) error {
		s.Batch = &Batch{ID: "b", Stories: []string{"S-0001", "S-0002"}, Step: StepFailed,
			Heads:   []Head{{Story: "S-0001", Verified: "v1", Rebased: "r1"}, {Story: "S-0002", Verified: "v2", Rebased: "r2"}},
			Bisect:  bisect,
			Landing: &Landing{Stories: []string{"S-0001"}, Landed: []string{}, Head: "base"}}
		return nil
	}); err != nil {
		t.Fatal(err)
	}
	st, err := ReadState(dir)
	if err != nil {
		t.Fatal(err)
	}
	b := st.Batch
	if b == nil || len(b.Heads) != 2 || b.Heads[1].Rebased != "r2" || !b.Bisect.Found() || b.Bisect.Failing != "S-0002" ||
		len(b.Bisect.Steps) != 1 || !b.Bisect.Steps[0].Result.Passed || b.Bisect.Result.StoppedAt != "unit" ||
		b.Landing == nil || b.Landing.Head != "base" {
		t.Errorf("batch read back %+v", b)
	}
	if (&Bisect{}).Found() || (*Bisect)(nil).Found() {
		t.Error("a bisection with no failing story is found")
	}
}
