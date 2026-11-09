package host

import (
	"context"
	"time"
)

// How the dashboard stands when the host's watch looks at it (S-0184).
const (
	DashboardRunning      = "running"       // it answers
	DashboardNotAnswering = "not-answering" // its container runs and does not answer
	DashboardGone         = "gone"          // no container runs
	DashboardPortHeld     = "port-held"     // its container runs, docker publishes nothing for it, and another process holds its recorded port (ADR-0267)
	DashboardOff          = "off"           // not watched: none was started, it was stopped, or dashboard.no_restart is set
)

// WatchMisses is how many looks in a row must find the dashboard gone or not
// answering before the host restarts it: one miss can be a swap that flai
// dashboard restart or upgrade is making.
const WatchMisses = 2

// Watch is what the host looks after besides its children: the dashboard
// container flai dashboard started. Look says how it stands, DashboardOff
// when it is not watched, and a detail: for DashboardPortHeld what holds the
// port, "" otherwise. Restart starts it again.
type Watch struct {
	Look    func(ctx context.Context) (state, detail string)
	Restart func(ctx context.Context) error
}

// DashboardWatch is how the watch stands, in the host's status.
type DashboardWatch struct {
	State       string `json:"state"`
	Detail      string `json:"detail,omitempty"` // the last look's detail: for port-held, what holds the port
	Restarts    int    `json:"restarts"`
	LastRestart string `json:"last_restart,omitempty"`
	LastReason  string `json:"last_reason,omitempty"`
	LastError   string `json:"last_error,omitempty"`
}

// watcher decides, one look at a time, when the dashboard is restarted: after
// WatchMisses looks in a row that did not find it running, and no sooner than
// a back-off after the last restart, which doubles from first to max and is
// reset once the dashboard has answered at every look for a minute.
type watcher struct {
	first, max time.Duration
	backoff    time.Duration
	misses     int
	restarted  time.Time // when the last restart was made
	answering  time.Time // since when it has answered at every look, zero after a miss
	next       time.Time // no restart before this
}

func newWatcher(first, max time.Duration) *watcher {
	return &watcher{first: first, max: max, backoff: first}
}

// look takes in one look's state at now and says whether to restart now.
func (w *watcher) look(now time.Time, state string) bool {
	switch state {
	case DashboardOff:
		w.misses = 0
		return false
	case DashboardRunning:
		w.misses = 0
		if w.answering.IsZero() {
			w.answering = now
		}
		if !w.restarted.IsZero() && now.Sub(w.answering) >= time.Minute {
			w.backoff = w.first
		}
		return false
	case DashboardPortHeld:
		// a restart cannot win the port back from the process that holds it
		// (ADR-0267): not a miss, and the back-off stands
		w.misses = 0
		w.answering = time.Time{}
		return false
	}
	w.answering = time.Time{}
	w.misses++
	if w.misses < WatchMisses || now.Before(w.next) {
		return false
	}
	w.misses = 0
	w.restarted, w.next = now, now.Add(w.backoff)
	w.backoff = min(w.backoff*2, w.max)
	return true
}

// watchDashboard looks at the dashboard every WatchEvery until ctx ends, and
// restarts it as the watcher decides, logging each restart with why, and a
// port another process holds once each time the watch finds it so.
func (h *host) watchDashboard(ctx context.Context) {
	o := h.o
	w := newWatcher(o.WatchBackoff, o.WatchMaxBackoff)
	t := time.NewTicker(o.WatchEvery)
	defer t.Stop()
	last := ""
	for {
		select {
		case <-ctx.Done():
			return
		case <-t.C:
		}
		state, detail := o.Dashboard.Look(ctx)
		h.setWatch(func(d *DashboardWatch) { d.State, d.Detail = state, detail })
		if state == DashboardPortHeld && last != DashboardPortHeld {
			o.Logger.Error("dashboard port held by another process", "component", "host", "holder", detail)
		}
		last = state
		now := o.Now()
		if !w.look(now, state) {
			continue
		}
		err := o.Dashboard.Restart(ctx)
		h.setWatch(func(d *DashboardWatch) {
			d.Restarts++
			d.LastRestart, d.LastReason, d.LastError = now.UTC().Format(time.RFC3339), state, ""
			if err != nil {
				d.LastError = err.Error()
			}
		})
		if err != nil {
			o.Logger.Error("dashboard not restarted", "component", "host", "reason", state, "err", err.Error(), "next_try_after", w.next.Sub(now).String())
			continue
		}
		o.Logger.Warn("dashboard restarted", "component", "host", "reason", state, "next_restart_after", w.next.Sub(now).String())
	}
}

func (h *host) setWatch(f func(*DashboardWatch)) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if h.watch == nil {
		h.watch = &DashboardWatch{}
	}
	f(h.watch)
}
