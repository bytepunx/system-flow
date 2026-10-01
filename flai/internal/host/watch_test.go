package host

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"
)

// S-0184: the watcher restarts after WatchMisses looks in a row, for a gone
// and for a not-answering dashboard, and never one that is off or answers.
func TestWatcherRestartsAfterMissesInARow(t *testing.T) {
	for _, miss := range []string{DashboardGone, DashboardNotAnswering} {
		w := newWatcher(30*time.Second, 5*time.Minute)
		now := time.Date(2026, 10, 1, 9, 0, 0, 0, time.UTC)
		if w.look(now, miss) {
			t.Errorf("%s: one miss is not enough; it may be a swap", miss)
		}
		if w.look(now.Add(15*time.Second), DashboardRunning) {
			t.Errorf("%s: an answer is never restarted", miss)
		}
		if w.look(now.Add(30*time.Second), miss) {
			t.Errorf("%s: an answer in between starts the count again", miss)
		}
		if !w.look(now.Add(45*time.Second), miss) {
			t.Errorf("%s: two misses in a row restart it", miss)
		}
	}
	w := newWatcher(time.Second, time.Minute)
	now := time.Date(2026, 10, 1, 9, 0, 0, 0, time.UTC)
	for i := range 10 {
		if w.look(now.Add(time.Duration(i)*time.Second), DashboardOff) {
			t.Fatal("a dashboard that is not watched is never restarted")
		}
	}
}

// S-0184: restarts wait a back-off that doubles up to the most, and is reset
// once the dashboard has answered for a minute after one.
func TestWatcherBacksOffBetweenRestarts(t *testing.T) {
	w := newWatcher(30*time.Second, 2*time.Minute)
	start := time.Date(2026, 10, 1, 9, 0, 0, 0, time.UTC)
	var restarts []time.Duration
	for s := 0; s <= 600; s += 15 { // a look every 15 seconds for ten minutes, always gone
		if w.look(start.Add(time.Duration(s)*time.Second), DashboardGone) {
			restarts = append(restarts, time.Duration(s)*time.Second)
		}
	}
	// first at the second miss, then 30s, 60s, 120s, 120s… apart, each
	// waiting for two more misses as well
	want := []time.Duration{15 * time.Second, 45 * time.Second, 105 * time.Second, 225 * time.Second, 345 * time.Second, 465 * time.Second, 585 * time.Second}
	if len(restarts) != len(want) {
		t.Fatalf("restarts at %v, want %v", restarts, want)
	}
	for i := range want {
		if restarts[i] != want[i] {
			t.Fatalf("restarts at %v, want %v", restarts, want)
		}
	}
	// it answers for a minute: the back-off is back to the first
	at := start.Add(600 * time.Second)
	w.look(at.Add(15*time.Second), DashboardRunning)
	w.look(at.Add(75*time.Second), DashboardRunning)
	w.look(at.Add(90*time.Second), DashboardGone)
	if !w.look(at.Add(105*time.Second), DashboardGone) {
		t.Fatal("the next restart comes at the second miss")
	}
	if w.look(at.Add(120*time.Second), DashboardGone) {
		t.Error("one miss after a restart is not enough")
	}
	if !w.look(at.Add(135*time.Second), DashboardGone) {
		t.Error("30s after the restart is the first back-off again, not the two minutes it had grown to")
	}
}

type fakeDashboard struct {
	mu       sync.Mutex
	states   []string // what each look finds; the last repeats
	looks    int
	restarts int
	err      error
}

func (f *fakeDashboard) look(context.Context) string {
	f.mu.Lock()
	defer f.mu.Unlock()
	i := min(f.looks, len(f.states)-1)
	f.looks++
	return f.states[i]
}

func (f *fakeDashboard) restart(context.Context) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.restarts++
	return f.err
}

func (f *fakeDashboard) count() int {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.restarts
}

// S-0184: a running host restarts a dashboard that is gone, says so in its
// status, and does not restart one that is off.
func TestHostRestartsAGoneDashboard(t *testing.T) {
	d := &fakeDashboard{states: []string{DashboardRunning, DashboardGone, DashboardGone, DashboardRunning}}
	r := start(t, Options{Dashboard: &Watch{Look: d.look, Restart: d.restart}, WatchEvery: 10 * time.Millisecond, WatchBackoff: time.Hour})
	st := r.until("dashboard restarted", func(st Status) bool { return st.Dashboard != nil && st.Dashboard.Restarts == 1 })
	if st.Dashboard.LastReason != DashboardGone || st.Dashboard.LastRestart == "" || st.Dashboard.LastError != "" {
		t.Errorf("status: %+v", *st.Dashboard)
	}
	r.until("dashboard running", func(st Status) bool { return st.Dashboard.State == DashboardRunning })
	if n := d.count(); n != 1 {
		t.Errorf("restarted %d times, want once", n)
	}

	off := &fakeDashboard{states: []string{DashboardOff}}
	r2 := start(t, Options{Dashboard: &Watch{Look: off.look, Restart: off.restart}, WatchEvery: 10 * time.Millisecond})
	r2.until("dashboard off", func(st Status) bool { return st.Dashboard != nil && st.Dashboard.State == DashboardOff })
	time.Sleep(50 * time.Millisecond)
	if n := off.count(); n != 0 {
		t.Errorf("a dashboard that is off was restarted %d times", n)
	}
}

// S-0184: a failed restart is in the status, and is tried again after the
// back-off.
func TestHostSaysWhyADashboardWasNotRestarted(t *testing.T) {
	d := &fakeDashboard{states: []string{DashboardNotAnswering}, err: errors.New("docker: no such image")}
	r := start(t, Options{Dashboard: &Watch{Look: d.look, Restart: d.restart}, WatchEvery: 10 * time.Millisecond, WatchBackoff: 20 * time.Millisecond})
	st := r.until("two tries", func(st Status) bool { return st.Dashboard != nil && st.Dashboard.Restarts >= 2 })
	if st.Dashboard.LastReason != DashboardNotAnswering || st.Dashboard.LastError != "docker: no such image" {
		t.Errorf("status: %+v", *st.Dashboard)
	}
}
