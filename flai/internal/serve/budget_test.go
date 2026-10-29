package serve

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"slices"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/config"
	"github.com/bytepunx/system-flow/flai/internal/hostapi"
	"github.com/bytepunx/system-flow/flai/internal/usage"
)

// budgetLab is one flai serve's serve folder with a launcher over a project,
// a clock, and the journal its record hook kept.
type budgetLab struct {
	t       *testing.T
	dir     Dir
	at      time.Time
	mu      sync.Mutex // guards journal, which launchers of two projects append to at once
	journal []hostapi.Entry
	logs    bytes.Buffer
	o       Options
	l       *launcher
}

// budgetAt is when the lab's clock starts, and budgetReset a reset time the
// API reports, an hour after it.
var (
	budgetAt    = time.Date(2026, 10, 28, 7, 0, 0, 0, time.UTC)
	budgetReset = budgetAt.Add(time.Hour)
)

func newBudgetLab(t *testing.T) *budgetLab {
	t.Helper()
	lab := &budgetLab{t: t, dir: Dir(filepath.Join(t.TempDir(), "serve")), at: budgetAt}
	lab.o = Options{Dir: lab.dir, Logger: slog.New(slog.NewTextHandler(&lab.logs, nil)), Now: func() time.Time { return lab.at },
		Host: hostapi.Host{Record: func(e hostapi.Entry) {
			lab.mu.Lock()
			defer lab.mu.Unlock()
			lab.journal = append(lab.journal, e)
		}}}
	lab.l = newLauncher(lab.o, Entry{Key: "t", Name: "t", Root: "/srv/t"})
	return lab
}

// run is a story's agent's run on harness claude-code with provider, ""
// for none, started at the lab's clock.
func (lab *budgetLab) run(story, provider string) *AgentRun {
	return &AgentRun{Story: story, Agent: "builder", Harness: "claude-code", Provider: provider, Session: "s-" + story, Started: lab.at.Format(time.RFC3339)}
}

// held is the budget state of key, and whether it holds a start at the lab's
// clock.
func (lab *budgetLab) held(key string) (Budget, bool) {
	lab.t.Helper()
	b, held, err := lab.dir.budgetHeld(key, lab.at)
	if err != nil {
		lab.t.Fatal(err)
	}
	return b, held
}

// limited is a refusal on the usage limit that names reset, zero for none.
func limited(reset time.Time) *usage.Limit {
	return &usage.Limit{Error: "usage_limit_reached", Status: 429, Resets: reset}
}

// spawned gives run a process of its own that runs until the test ends, as
// the launcher's spawn does, and records it on the probe.
func (lab *budgetLab) spawned(run *AgentRun) {
	lab.t.Helper()
	cmd := exec.Command("sleep", "60")
	Detach(cmd)
	if err := cmd.Start(); err != nil {
		lab.t.Fatal(err)
	}
	lab.t.Cleanup(func() { _ = cmd.Process.Kill(); _ = cmd.Wait() })
	run.PID, run.Start = cmd.Process.Pid, Started(cmd.Process.Pid)
	lab.l.budgetSpawned(run)
}

// A run that ends on the limit the harness's log names sets its provider's
// state until the reset the API reported, read from I-0212's recorded
// Claude Code result line.
func TestARefusalFromARecordedLogSetsTheStateUntilItsReset(t *testing.T) {
	lab := newBudgetLab(t)
	run := lab.run("S-0001", "")
	run.Log = limitLog(t, limitStream...)
	var w warned
	lim, _ := runLimit(run, w.warn)
	if lim == nil {
		t.Fatalf("no limit read from the log: %v", w)
	}
	lab.l.budgetEnded(run, lim)
	b, held := lab.held("claude-code")
	want, _ := time.Parse(time.RFC3339, limitResetsAt)
	if !held || !b.Until.Equal(want) || !b.Reset.Equal(want) || b.Limit != "usage_limit_reached" || b.Status != 429 || !b.Since.Equal(budgetAt) {
		t.Errorf("state = %+v, held %v; want held until %s", b, held, want)
	}
	lab.at = want.Add(-time.Second)
	if _, held := lab.held("claude-code"); !held {
		t.Error("not held a second before the reset")
	}
	lab.at = want
	if _, held := lab.held("claude-code"); held {
		t.Error("held at the reset")
	}
}

// A refusal with no reset holds its provider 5m, and each refused probe
// doubles the hold up to 1h.
func TestABackoffDoublesOnEachRefusedProbeUpToAnHour(t *testing.T) {
	lab := newBudgetLab(t)
	lab.l.budgetEnded(lab.run("S-0001", "gateway"), limited(time.Time{}))
	for i, want := range []time.Duration{5 * time.Minute, 10 * time.Minute, 20 * time.Minute, 40 * time.Minute, time.Hour, time.Hour} {
		b, held := lab.held("gateway")
		if !held || b.Until.Sub(lab.at) != want || b.backoff() != want || !b.Reset.IsZero() {
			t.Fatalf("hold %d: state = %+v, held %v; want held for %s", i, b, held, want)
		}
		lab.at = b.Until
		probe := lab.run("S-0002", "gateway")
		if gate := lab.l.budgetStart(probe); !gate.Go || !gate.Probe {
			t.Fatalf("hold %d: once it passed, the start is %+v, want the probe", i, gate)
		}
		lab.l.budgetEnded(probe, limited(time.Time{}))
	}
}

// Held answers false for another provider, and for a provider whose state
// was cleared; a provider is its harness when the run named none.
func TestHeldIsPerProviderAndEndsWhenCleared(t *testing.T) {
	lab := newBudgetLab(t)
	lab.l.budgetEnded(lab.run("S-0001", "anthropic"), limited(budgetReset))
	if _, held := lab.held("anthropic"); !held {
		t.Fatal("the provider that ended on the limit is not held")
	}
	if _, held := lab.held("gateway"); held {
		t.Error("another provider is held")
	}
	if _, held := lab.held("claude-code"); held {
		t.Error("the harness of runs that named no provider is held")
	}
	if got := budgetKey("", "claude-code"); got != "claude-code" {
		t.Errorf("a run with no provider is kept under %q, want its harness", got)
	}
	lab.at = budgetReset
	probe := lab.run("S-0002", "anthropic")
	if gate := lab.l.budgetStart(probe); !gate.Probe {
		t.Fatalf("start = %+v, want the probe", gate)
	}
	lab.l.budgetEnded(probe, nil)
	if b, held := lab.held("anthropic"); held || b.Provider != "" {
		t.Errorf("a cleared provider is %+v, held %v", b, held)
	}
	if all, err := lab.dir.Budgets(); err != nil || len(all) != 0 {
		t.Errorf("budgets = %+v, %v; want none", all, err)
	}
}

// Once until has passed, one start goes through as the probe and every other
// waits until it ends: a clean end clears the state and starts go again; an
// end on the limit extends it, and frees the probe.
func TestOneProbeGoesThroughAndTheRestWaitForItsEnd(t *testing.T) {
	if runtime.GOOS != "linux" {
		t.Skip("a probe's process is told from another by /proc")
	}
	for _, refused := range []bool{false, true} {
		lab := newBudgetLab(t)
		lab.l.budgetEnded(lab.run("S-0001", ""), limited(budgetReset))
		if gate := lab.l.budgetStart(lab.run("S-0002", "")); gate.Go || !gate.State.Until.Equal(budgetReset) {
			t.Fatalf("a start before the reset is %+v, want held until it", gate)
		}
		lab.at = budgetReset
		probe := lab.run("S-0002", "")
		if gate := lab.l.budgetStart(probe); !gate.Go || !gate.Probe {
			t.Fatalf("the first start after the reset is %+v, want the probe", gate)
		}
		lab.spawned(probe)
		lab.at = budgetReset.Add(budgetProbeServed - time.Second)
		for _, story := range []string{"S-0003", "S-0004"} {
			if gate := lab.l.budgetStart(lab.run(story, "")); gate.Go {
				t.Errorf("%s started while the probe runs: %+v", story, gate)
			}
		}
		if _, held := lab.held("claude-code"); !held {
			t.Error("not held while the probe runs")
		}
		if refused {
			lab.l.budgetEnded(probe, limited(time.Time{}))
			b, held := lab.held("claude-code")
			if !held || b.Probe != nil || b.Until.Sub(lab.at) != 10*time.Minute {
				t.Errorf("a refused probe left %+v, held %v; want held 10m, the backoff doubled, with no probe", b, held)
			}
			continue
		}
		lab.l.budgetEnded(probe, nil)
		if gate := lab.l.budgetStart(lab.run("S-0003", "")); !gate.Go || gate.Probe {
			t.Errorf("after the probe's clean end a start is %+v, want it to go", gate)
		}
	}
}

// A probe claimed by a start that never recorded a process, or whose process
// is gone without its end being recorded, holds no start for ever: the next
// start after budgetProbeSpawn, or once the process is gone, probes in its
// place. A start that could not be made frees it at once.
func TestAProbeThatNeverEndsIsClaimedOver(t *testing.T) {
	lab := newBudgetLab(t)
	lab.l.budgetEnded(lab.run("S-0001", ""), limited(budgetReset))
	lab.at = budgetReset
	first := lab.run("S-0002", "")
	if gate := lab.l.budgetStart(first); !gate.Probe {
		t.Fatalf("start = %+v, want the probe", gate)
	}
	lab.l.budgetRelease(first)
	second := lab.run("S-0003", "")
	if gate := lab.l.budgetStart(second); !gate.Probe {
		t.Fatalf("after a release the start is %+v, want the probe", gate)
	}
	lab.at = budgetReset.Add(budgetProbeSpawn - time.Second)
	if gate := lab.l.budgetStart(lab.run("S-0004", "")); gate.Go {
		t.Errorf("a start while a probe is being spawned is %+v, want held", gate)
	}
	lab.at = budgetReset.Add(budgetProbeSpawn)
	third := lab.run("S-0004", "")
	if gate := lab.l.budgetStart(third); !gate.Probe {
		t.Fatalf("after a probe that recorded no process the start is %+v, want the probe", gate)
	}
	third.PID = 999999999 // a process that is gone
	lab.l.budgetSpawned(third)
	if b, held := lab.held("claude-code"); held || b.Probe == nil || b.Probe.PID != third.PID {
		t.Errorf("a probe whose process is gone is %+v, held %v", b.Probe, held)
	}
	if gate := lab.l.budgetStart(lab.run("S-0005", "")); !gate.Probe {
		t.Errorf("after a probe whose process is gone the start is %+v, want the probe", gate)
	}
}

// A probe that runs past budgetProbeServed was served by the API: the next
// start clears the state and goes.
func TestAProbeThatRunsOnClearsTheState(t *testing.T) {
	if runtime.GOOS != "linux" {
		t.Skip("a probe's process is told from another by /proc")
	}
	lab := newBudgetLab(t)
	lab.l.budgetEnded(lab.run("S-0001", ""), limited(time.Time{}))
	lab.at = lab.at.Add(budgetBackoffFirst)
	probe := lab.run("S-0002", "")
	lab.l.budgetStart(probe)
	lab.spawned(probe)
	lab.at = lab.at.Add(budgetProbeServed)
	if _, held := lab.held("claude-code"); held {
		t.Error("held by a probe that has run past budgetProbeServed")
	}
	if gate := lab.l.budgetStart(lab.run("S-0003", "")); !gate.Go || gate.Probe {
		t.Errorf("start = %+v, want it to go", gate)
	}
	if all, _ := lab.dir.Budgets(); len(all) != 0 {
		t.Errorf("budgets = %+v, want none", all)
	}
}

// The runs that were going when the state was set, refused within its hold,
// do not compound it; a later reset extends it.
func TestRefusalsWithinTheHoldDoNotCompoundIt(t *testing.T) {
	lab := newBudgetLab(t)
	lab.l.budgetEnded(lab.run("S-0001", ""), limited(time.Time{}))
	lab.at = lab.at.Add(time.Minute)
	lab.l.budgetEnded(lab.run("S-0002", ""), limited(time.Time{}))
	if b, _ := lab.held("claude-code"); !b.Until.Equal(budgetAt.Add(budgetBackoffFirst)) || b.backoff() != budgetBackoffFirst {
		t.Errorf("a second refusal within the hold made it %+v", b)
	}
	lab.l.budgetEnded(lab.run("S-0003", ""), limited(budgetReset))
	if b, _ := lab.held("claude-code"); !b.Until.Equal(budgetReset) || !b.Reset.Equal(budgetReset) {
		t.Errorf("a later reset left it %+v, want until %s", b, budgetReset)
	}
	lab.l.budgetEnded(lab.run("S-0004", ""), nil)
	if _, held := lab.held("claude-code"); !held {
		t.Error("a run that is not the probe, ending clean, cleared the state")
	}
	if n := len(lab.journal); n != 2 {
		t.Errorf("journal has %d entries, want the set and the extension: %+v", n, lab.journal)
	}
}

// The state survives a restart of flai serve: a fresh Dir of the same path
// reads it back from budget.json; a missing file is no state.
func TestTheStateIsReadBackAfterARestart(t *testing.T) {
	lab := newBudgetLab(t)
	if all, err := lab.dir.Budgets(); err != nil || len(all) != 0 {
		t.Fatalf("budgets before any = %+v, %v; want none", all, err)
	}
	lab.l.budgetEnded(lab.run("S-0001", "zeta"), limited(budgetReset))
	lab.l.budgetEnded(lab.run("S-0002", "alpha"), limited(time.Time{}))
	lab.at = lab.at.Add(budgetBackoffFirst)
	lab.l.budgetStart(lab.run("S-0003", "alpha"))
	fresh := Dir(string(lab.dir))
	all, err := fresh.Budgets()
	if err != nil {
		t.Fatal(err)
	}
	if len(all) != 2 || all[0].Provider != "alpha" || all[1].Provider != "zeta" || !all[1].Until.Equal(budgetReset) ||
		all[0].Probe == nil || all[0].Probe.Work != "S-0003" || all[0].Probe.Root != "/srv/t" {
		t.Errorf("budgets read back = %+v", all)
	}
	data, err := os.ReadFile(filepath.Join(string(lab.dir), "budget.json"))
	if err != nil {
		t.Fatal(err)
	}
	for _, field := range []string{`"provider"`, `"limit"`, `"status"`, `"since"`, `"until"`, `"backoff"`, `"probe"`} {
		if !strings.Contains(string(data), field) {
			t.Errorf("budget.json has no %s: %s", field, data)
		}
	}
}

// Each change appends one journal entry, by flai serve, naming the
// provider and until, and logs it; a start that goes, or is held, with
// nothing changed, journals nothing.
func TestEachChangeIsJournalledOnce(t *testing.T) {
	lab := newBudgetLab(t)
	lab.l.budgetStart(lab.run("S-0000", "anthropic"))
	lab.l.budgetEnded(lab.run("S-0001", "anthropic"), limited(time.Time{}))
	lab.l.budgetStart(lab.run("S-0002", "anthropic"))
	lab.at = lab.at.Add(budgetBackoffFirst)
	refused := lab.run("S-0003", "anthropic")
	lab.l.budgetStart(refused)
	lab.l.budgetEnded(refused, limited(time.Time{}))
	lab.at = lab.at.Add(2 * budgetBackoffFirst)
	released := lab.run("S-0004", "anthropic")
	lab.l.budgetStart(released)
	lab.l.budgetRelease(released)
	clean := lab.run("S-0005", "anthropic")
	lab.l.budgetStart(clean)
	lab.l.budgetEnded(clean, nil)
	until1 := budgetAt.Add(budgetBackoffFirst).Format(time.RFC3339)
	until2 := budgetAt.Add(3 * budgetBackoffFirst).Format(time.RFC3339)
	want := []struct{ starts, until string }{
		{"set the budget state", until1},
		{"started a probe", until1},
		{"extended the budget state", until2},
		{"started a probe", until2},
		{"freed the budget state's probe", until2},
		{"started a probe", until2},
		{"cleared the budget state", until2},
	}
	if len(lab.journal) != len(want) {
		t.Fatalf("journal has %d entries, want %d: %+v", len(lab.journal), len(want), lab.journal)
	}
	for i, w := range want {
		e := lab.journal[i]
		if !strings.HasPrefix(e.Detail, w.starts) || !strings.Contains(e.Detail, "provider anthropic") || !strings.Contains(e.Detail, "until "+w.until) ||
			!strings.Contains(e.Detail, "usage_limit_reached (429)") || e.Action != hostapi.ActionAgent || e.Method != "serve.budget" ||
			e.By != "flai serve" || e.Outcome != "done" || e.Root != "/srv/t" || e.Project != "t" {
			t.Errorf("entry %d = %+v, want %q naming the provider and until %s", i, e, w.starts, w.until)
		}
	}
	if n := strings.Count(lab.logs.String(), `msg="budget state `); n != len(want) {
		t.Errorf("logged %d changes, want %d:\n%s", n, len(want), lab.logs.String())
	}
}

// noticeBox is a fake sender of a serve folder's budget notices: what it was
// sent, in order.
type noticeBox struct {
	mu   sync.Mutex
	sent []BudgetNotice
}

// notices has the lab's serve folder send its budget notices to a fake
// through ch, until the test ends, and registers two projects in it.
func (lab *budgetLab) notices(ch config.Notify) *noticeBox {
	lab.t.Helper()
	for _, e := range []Entry{{Key: "t", Name: "system-flow", Root: "/srv/t"}, {Key: "u", Root: "/srv/u"}} {
		e.URL, e.KeyFile = "ws://127.0.0.1:1/"+e.Key, "/srv/"+e.Key+".key"
		if err := lab.dir.Register(e); err != nil {
			lab.t.Fatal(err)
		}
	}
	box := &noticeBox{}
	lab.dir.notifyBudgetWith(budgetNotifier{channel: func() config.Notify { return ch }, send: func(_ context.Context, got config.Notify, n BudgetNotice) error {
		if got != ch {
			lab.t.Errorf("sent through %+v, want %+v", got, ch)
		}
		box.mu.Lock()
		defer box.mu.Unlock()
		box.sent = append(box.sent, n)
		return nil
	}})
	lab.t.Cleanup(func() { lab.dir.notifyBudget(nil) })
	return box
}

// states are the states of the notices sent, in order.
func (box *noticeBox) states() []string {
	box.mu.Lock()
	defer box.mu.Unlock()
	out := []string{}
	for _, n := range box.sent {
		out = append(out, n.State)
	}
	return out
}

// S-0534 (ADR-0209 §5): a set, an extension after a refused probe, and a clear
// each send one notice of kind budget with the provider, the limit, until, the
// projects the host serves, and the dashboard's root; a probe's start, a start
// held, and a probe freed send none.
func TestEachSetExtensionAndClearSendsOneNotice(t *testing.T) {
	lab := newBudgetLab(t)
	box := lab.notices(config.Notify{URL: "https://push.example/flai"})
	lab.l.budgetEnded(lab.run("S-0001", "anthropic"), limited(time.Time{}))
	lab.l.budgetStart(lab.run("S-0002", "anthropic"))
	lab.at = lab.at.Add(budgetBackoffFirst)
	refused := lab.run("S-0003", "anthropic")
	lab.l.budgetStart(refused)
	lab.l.budgetEnded(refused, limited(time.Time{}))
	lab.at = lab.at.Add(2 * budgetBackoffFirst)
	released := lab.run("S-0004", "anthropic")
	lab.l.budgetStart(released)
	lab.l.budgetRelease(released)
	clean := lab.run("S-0005", "anthropic")
	lab.l.budgetStart(clean)
	lab.l.budgetEnded(clean, nil)

	until1, until2 := budgetAt.Add(budgetBackoffFirst), budgetAt.Add(3*budgetBackoffFirst)
	want := []struct {
		state string
		until time.Time
		at    time.Time
	}{{"set", until1, budgetAt}, {"extended", until2, until1}, {"cleared", until2, until2}}
	if len(box.sent) != len(want) {
		t.Fatalf("sent %v, want one notice each of set, extended, and cleared", box.states())
	}
	for i, w := range want {
		n := box.sent[i]
		if n.State != w.state || n.Provider != "anthropic" || n.Limit != "usage_limit_reached" || n.Status != 429 || !n.Until.Equal(w.until) ||
			!n.Since.Equal(budgetAt) || !n.At.Equal(w.at) || n.Href != "/" || !slices.Equal(n.Projects, []string{"system-flow", "u"}) {
			t.Errorf("notice %d = %+v, want %s until %s at %s", i, n, w.state, w.until, w.at)
		}
		data, err := json.Marshal(n)
		if err != nil || !strings.Contains(string(data), `"kind":"budget"`) {
			t.Errorf("notice %d is not of kind budget: %s, %v", i, data, err)
		}
	}
	if n := strings.Count(lab.logs.String(), `msg="budget notice sent"`); n != len(want) {
		t.Errorf("logged %d notices sent, want %d:\n%s", n, len(want), lab.logs.String())
	}
}

// The budget state is the host's: launchers of two projects under one serve
// folder, ending and starting runs at once, send one notice per change, not
// one per project.
func TestTwoProjectsUnderOneServeFolderSendOneNoticePerChange(t *testing.T) {
	lab := newBudgetLab(t)
	box := lab.notices(config.Notify{Command: "notify-send flai"})
	launchers := []*launcher{lab.l, newLauncher(lab.o, Entry{Key: "u", Name: "u", Root: "/srv/u"})}
	runs := make([]*AgentRun, len(launchers))
	gates := make([]budgetGate, len(launchers))
	// both runs each launcher's do at once, on its own run of story-<i>
	both := func(story string, do func(l *launcher, run *AgentRun) budgetGate) {
		var wg sync.WaitGroup
		for i, l := range launchers {
			runs[i] = lab.run(fmt.Sprintf("%s-%d", story, i), "anthropic")
			wg.Add(1)
			go func() {
				defer wg.Done()
				gates[i] = do(l, runs[i])
			}()
		}
		wg.Wait()
	}
	end := func(lim *usage.Limit) func(*launcher, *AgentRun) budgetGate {
		return func(l *launcher, run *AgentRun) budgetGate { l.budgetEnded(run, lim); return budgetGate{} }
	}
	start := func(l *launcher, run *AgentRun) budgetGate { return l.budgetStart(run) }
	// probe is the index of the one start of both that probes
	probe := func() int {
		t.Helper()
		if gates[0].Probe == gates[1].Probe {
			t.Fatalf("starts at once = %+v, want one probe", gates)
		}
		if gates[0].Probe {
			return 0
		}
		return 1
	}

	both("S-0001", end(limited(time.Time{})))
	lab.at = lab.at.Add(budgetBackoffFirst)
	both("S-0002", start)
	p := probe()
	launchers[p].budgetEnded(runs[p], limited(time.Time{}))
	both("S-0003", end(limited(time.Time{})))
	b, _ := lab.held("anthropic")
	lab.at = b.Until
	both("S-0004", start)
	p = probe()
	launchers[p].budgetEnded(runs[p], nil)

	if got := box.states(); !slices.Equal(got, []string{"set", "extended", "cleared"}) {
		t.Errorf("sent %v, want one notice each of set, extended, and cleared", got)
	}
	if n := len(lab.journal); n != 5 {
		t.Errorf("journal has %d entries, want the set, two probes, the extension, and the clear: %+v", n, lab.journal)
	}
}

// With no channel set, or none registered for the serve folder, no notice is
// sent, each change is logged, and that no channel is set is warned of once
// in the run, whichever project's launcher made the change.
func TestNoChannelSendsNothingAndWarnsOnce(t *testing.T) {
	for _, registered := range []bool{true, false} {
		lab := newBudgetLab(t)
		box := &noticeBox{}
		if registered {
			box = lab.notices(config.Notify{Command: "  "})
		}
		other := newLauncher(lab.o, Entry{Key: "u", Name: "u", Root: "/srv/u"})
		lab.l.budgetEnded(lab.run("S-0001", "anthropic"), limited(time.Time{}))
		lab.at = lab.at.Add(budgetBackoffFirst)
		refused := lab.run("S-0002", "anthropic")
		other.budgetStart(refused)
		other.budgetEnded(refused, limited(time.Time{}))
		lab.at = lab.at.Add(2 * budgetBackoffFirst)
		clean := lab.run("S-0003", "anthropic")
		lab.l.budgetStart(clean)
		lab.l.budgetEnded(clean, nil)
		logs := lab.logs.String()
		if len(box.sent) != 0 {
			t.Errorf("registered %v: sent %v with no channel set", registered, box.states())
		}
		if n := strings.Count(logs, `msg="budget notices have no channel"`); n != 1 {
			t.Errorf("registered %v: warned %d times that no channel is set, want once:\n%s", registered, n, logs)
		}
		for _, kind := range []string{"set", "extended", "cleared"} {
			if !strings.Contains(logs, `msg="budget state `+kind+`" `) {
				t.Errorf("registered %v: the %s is not logged:\n%s", registered, kind, logs)
			}
		}
		if strings.Contains(logs, "budget notice sent") {
			t.Errorf("registered %v: a notice is logged as sent:\n%s", registered, logs)
		}
	}
}
