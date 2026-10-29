package serve

import (
	"cmp"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/config"
	"github.com/bytepunx/system-flow/flai/internal/harness"
	"github.com/bytepunx/system-flow/flai/internal/hostapi"
	"github.com/bytepunx/system-flow/flai/internal/manifest"
	"github.com/bytepunx/system-flow/flai/internal/usage"
)

// The host's budget state (S-0534, ADR-0209 §5). When a run of any role, or
// a story's agent, ends on the harness's usage limit, flai serve records a
// budget state for the provider the run went through, in budget.json in the
// serve folder: the limit the harness named, when it was first set, and the
// time it holds until, the reset the API reported, or else a backoff that
// starts at 5m and doubles, up to 1h, each time the limit outlasts it. While
// it holds, no run starts on that provider on any project the host serves.
// Once it has passed, one start goes through as the probe, and every other
// waits while it runs: its end on anything but the limit clears the state,
// and its end on the limit extends it.
//
// Every project's launcher in one flai serve reads and writes the file under
// budgetMu; the MCP server and the host API read it, from other processes,
// with Budgets. A change is announced in one place, launcher.budgetChanged,
// which journals it, logs it, and sends a set, an extension, or a clear to
// the operator through the channel of ADR-0205 as a BudgetNotice.

const (
	// budgetBackoffFirst is how long a limit with no reset time holds its
	// provider at first, and budgetBackoffMax the most a backoff doubles to.
	budgetBackoffFirst = 5 * time.Minute
	budgetBackoffMax   = time.Hour
	// budgetProbeSpawn is how long a probe claimed for a start that has not
	// recorded its process holds: a start that claimed it and then failed
	// without releasing it, or a flai serve that stopped between the two,
	// frees it after this.
	budgetProbeSpawn = 2 * time.Minute
	// budgetProbeServed is how long a probe's process may run before the
	// API is taken to serve its provider again: the harness ends a run the
	// limit refuses at its first call, within seconds, so a probe still
	// running after this clears the state rather than hold every start for
	// as long as a story's agent works.
	budgetProbeServed = 10 * time.Minute
)

// The kinds of change to a provider's budget state.
const (
	budgetSet      = "set"      // a run ended on the limit, and none was set
	budgetExtended = "extended" // a run ended on the limit, and it holds longer
	budgetProbing  = "probe"    // a start goes through as the probe
	budgetReleased = "released" // the probe's start could not be made
	budgetCleared  = "cleared"  // the probe ended on anything but the limit, or ran past budgetProbeServed
)

// Budget is one provider's budget state, as budget.json holds it and the MCP
// server and the host API read it.
type Budget struct {
	// Provider is what the state is kept under (budgetKey): the provider the
	// runs named, or their harness when they named none.
	Provider string `json:"provider"`
	// Limit is the error the harness named, such as usage_limit_reached, and
	// Status its HTTP status, such as 429; each empty when it named none.
	Limit  string `json:"limit,omitempty"`
	Status int    `json:"status,omitempty"`
	// Since is when the state was set.
	Since time.Time `json:"since"`
	// Until is when the hold ends and a probe may start: Reset when the API
	// reported one, else when Backoff ends.
	Until time.Time `json:"until"`
	// Reset is the reset time the API reported for the last refusal; zero
	// when it reported none.
	Reset time.Time `json:"reset,omitzero"`
	// Backoff is the hold a refusal with no reset time gets, as a Go
	// duration: 5m0s at first, doubled up to 1h0m0s each time the limit
	// outlasts it.
	Backoff string `json:"backoff"`
	// Probe is the run started as the probe once Until passed; nil while
	// none is.
	Probe *BudgetProbe `json:"probe,omitempty"`
}

// BudgetProbe is a run as the budget state knows it: the project it runs in,
// what it works, its session, its start, and, once its process is spawned,
// that process. A probe is matched to its run by the project, the work, and
// the start.
type BudgetProbe struct {
	Root string `json:"root"`
	// Work is what the run works: its story, the item it plans, or else its
	// agent's name, such as orchestrator.
	Work    string `json:"work"`
	Session string `json:"session,omitempty"`
	Started string `json:"started"`
	// Claimed is when the probe was let through.
	Claimed time.Time `json:"claimed,omitzero"`
	// PID and Start are its process, as AgentRun records them; zero until
	// it is spawned.
	PID   int   `json:"pid,omitempty"`
	Start int64 `json:"start,omitempty"`
}

// budgetGate is what the budget state answers a start: whether it goes,
// whether it goes as the probe, and, when it is held or probes, the state.
type budgetGate struct {
	Go    bool
	Probe bool
	State Budget
}

// budgetChange is what an operation changed: its kind, one of the budget
// kinds, empty when nothing changed; the provider's state after it, or as it
// was for a clear; the run whose start or end made it; and why, for a clear.
type budgetChange struct {
	Kind  string
	State Budget
	Run   BudgetProbe
	Why   string
}

// budgetMu orders the read-modify-write of budget.json by the launchers of
// the projects one flai serve serves, which start and end runs at once.
var budgetMu sync.Mutex

func (d Dir) budgetFile() string { return filepath.Join(string(d), "budget.json") }

// budgetKey is what a run's budget state is kept under: the provider it
// named, or its harness, such as claude-code, when it named none.
func budgetKey(provider, harness string) string {
	if provider != "" {
		return provider
	}
	return harness
}

// probeOf is run, in the project at root, as the budget state knows it.
func probeOf(root string, run *AgentRun) BudgetProbe {
	work := run.Story
	if work == "" {
		work = run.Item
	}
	if work == "" {
		work = run.Agent
	}
	return BudgetProbe{Root: root, Work: work, Session: run.Session, Started: run.Started, PID: run.PID, Start: run.Start}
}

// is says whether p and o are one run.
func (p *BudgetProbe) is(o BudgetProbe) bool {
	return p != nil && p.Root == o.Root && p.Work == o.Work && p.Started == o.Started
}

// running says whether the probe's run may still be going at now: its
// process is still the one spawned for it, or it has none yet and was
// claimed less than budgetProbeSpawn ago.
func (p *BudgetProbe) running(now time.Time) bool {
	if p == nil {
		return false
	}
	if p.PID == 0 {
		return now.Sub(p.Claimed) < budgetProbeSpawn
	}
	return Owns(p.PID, p.Start)
}

// served says whether the probe's process has run past budgetProbeServed at
// now, which the API's refusal would not have let it.
func (p *BudgetProbe) served(now time.Time) bool {
	return p != nil && p.PID != 0 && now.Sub(p.Claimed) >= budgetProbeServed && Owns(p.PID, p.Start)
}

// Held says whether b holds a start at now: until Until, and after it while
// the probe runs, up to budgetProbeServed.
func (b Budget) Held(now time.Time) bool {
	return now.Before(b.Until) || (b.Probe.running(now) && !b.Probe.served(now))
}

// backoff is b's backoff, the first for one not recorded or not valid.
func (b Budget) backoff() time.Duration {
	d, err := time.ParseDuration(b.Backoff)
	if err != nil || d <= 0 {
		return budgetBackoffFirst
	}
	return d
}

// budgets reads the budget state, by provider; none is an empty map.
func (d Dir) budgets() (map[string]Budget, error) {
	out := map[string]Budget{}
	data, err := os.ReadFile(d.budgetFile())
	if errors.Is(err, os.ErrNotExist) {
		return out, nil
	}
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(data, &out); err != nil {
		return nil, fmt.Errorf("%s: %w; remove it and flai serve holds no start until a run ends on a usage limit again", d.budgetFile(), err)
	}
	return out, nil
}

// Budgets is the budget state of each provider that has one, by provider; a
// missing file is none.
func (d Dir) Budgets() ([]Budget, error) {
	all, err := d.budgets()
	if err != nil {
		return nil, err
	}
	out := make([]Budget, 0, len(all))
	for _, b := range all {
		out = append(out, b)
	}
	slices.SortFunc(out, func(a, b Budget) int { return cmp.Compare(a.Provider, b.Provider) })
	return out, nil
}

// budgetHeld is the budget state of key, and whether it holds a start at now.
func (d Dir) budgetHeld(key string, now time.Time) (Budget, bool, error) {
	budgetMu.Lock()
	defer budgetMu.Unlock()
	all, err := d.budgets()
	if err != nil {
		return Budget{}, false, err
	}
	b, ok := all[key]
	return b, ok && b.Held(now), nil
}

// budgetStart is the gate a start of run on key passes before it spawns: it
// goes when no state is set; it is held while the state holds; once Until
// has passed it goes as the probe, claimed for it, unless a probe still runs.
// A probe whose process is gone, or that recorded none within
// budgetProbeSpawn, is claimed over; one that ran past budgetProbeServed
// clears the state, and the start goes.
func (d Dir) budgetStart(key string, run BudgetProbe, now time.Time) (budgetGate, budgetChange, error) {
	budgetMu.Lock()
	defer budgetMu.Unlock()
	all, err := d.budgets()
	if err != nil {
		return budgetGate{}, budgetChange{}, err
	}
	b, ok := all[key]
	switch {
	case !ok:
		return budgetGate{Go: true}, budgetChange{}, nil
	case now.Before(b.Until):
		return budgetGate{State: b}, budgetChange{}, nil
	case b.Probe.served(now):
		delete(all, key)
		ch := budgetChange{Kind: budgetCleared, State: b, Run: *b.Probe, Why: fmt.Sprintf("the probe ran past %s, which the limit's refusal would not have let it", budgetProbeServed)}
		return budgetGate{Go: true}, ch, d.write(d.budgetFile(), all)
	case b.Probe.running(now):
		return budgetGate{State: b}, budgetChange{}, nil
	}
	run.Claimed, run.PID, run.Start = now.UTC(), 0, 0
	b.Probe = &run
	all[key] = b
	return budgetGate{Go: true, Probe: true, State: b}, budgetChange{Kind: budgetProbing, State: b, Run: run}, d.write(d.budgetFile(), all)
}

// budgetSpawned records the process of run, the probe on key, once it is
// spawned; nothing when run is not the probe.
func (d Dir) budgetSpawned(key string, run BudgetProbe) error {
	budgetMu.Lock()
	defer budgetMu.Unlock()
	all, err := d.budgets()
	if err != nil {
		return err
	}
	b, ok := all[key]
	if !ok || !b.Probe.is(run) {
		return nil
	}
	b.Probe.PID, b.Probe.Start = run.PID, run.Start
	all[key] = b
	return d.write(d.budgetFile(), all)
}

// budgetRelease frees the probe on key for run, whose start could not be
// made, so that the next start probes; nothing when run is not the probe.
func (d Dir) budgetRelease(key string, run BudgetProbe) (budgetChange, error) {
	budgetMu.Lock()
	defer budgetMu.Unlock()
	all, err := d.budgets()
	if err != nil {
		return budgetChange{}, err
	}
	b, ok := all[key]
	if !ok || !b.Probe.is(run) {
		return budgetChange{}, nil
	}
	b.Probe = nil
	all[key] = b
	return budgetChange{Kind: budgetReleased, State: b, Run: run}, d.write(d.budgetFile(), all)
}

// budgetEnded records the end of run on key, on lim, the usage limit it
// ended on, or nil when it ended on none, at now. A run on the limit sets
// the state when none is set: until the reset lim names, when it is after
// now, else for the first backoff. When one is set, a later reset extends it
// to that reset, and a refusal after Until passed, the probe's included,
// doubles the backoff, up to budgetBackoffMax, and holds for it: the runs
// that were going when the state was set, refused within its hold, do not
// compound it. The probe's end frees the probe, and its end on anything but
// the limit clears the state. Any other end changes nothing.
func (d Dir) budgetEnded(key string, run BudgetProbe, lim *usage.Limit, now time.Time) (budgetChange, error) {
	budgetMu.Lock()
	defer budgetMu.Unlock()
	all, err := d.budgets()
	if err != nil {
		return budgetChange{}, err
	}
	now = now.UTC()
	b, ok := all[key]
	probe := ok && b.Probe.is(run)
	var ch budgetChange
	switch {
	case lim == nil && !probe:
		return budgetChange{}, nil
	case lim == nil:
		delete(all, key)
		ch = budgetChange{Kind: budgetCleared, State: b, Run: run, Why: "the probe ended on no usage limit"}
	case !ok:
		b = Budget{Provider: key, Limit: lim.Error, Status: lim.Status, Since: now, Until: now.Add(budgetBackoffFirst), Backoff: budgetBackoffFirst.String()}
		if lim.Resets.After(now) {
			b.Until, b.Reset = lim.Resets.UTC(), lim.Resets.UTC()
		}
		all[key] = b
		ch = budgetChange{Kind: budgetSet, State: b, Run: run}
	default:
		extended := probe
		if lim.Resets.After(now) {
			if lim.Resets.After(b.Until) {
				b.Until, b.Reset, extended = lim.Resets.UTC(), lim.Resets.UTC(), true
			}
		} else if !now.Before(b.Until) {
			backoff := min(2*b.backoff(), budgetBackoffMax)
			b.Until, b.Reset, b.Backoff, extended = now.Add(backoff), time.Time{}, backoff.String(), true
		}
		if !extended {
			return budgetChange{}, nil
		}
		if lim.Error != "" || lim.Status != 0 {
			b.Limit, b.Status = lim.Error, lim.Status
		}
		if probe {
			b.Probe = nil
		}
		all[key] = b
		ch = budgetChange{Kind: budgetExtended, State: b, Run: run}
	}
	return ch, d.write(d.budgetFile(), all)
}

// budgetSaid names b's provider, the limit it is on, and when it holds
// until, such as "provider anthropic, on usage_limit_reached (429), held
// until 2026-10-28T08:00:00Z, the reset the API reported".
func budgetSaid(b Budget) string {
	limit := b.Limit
	if limit == "" {
		limit = "the harness's usage limit"
	}
	if b.Status != 0 {
		limit = fmt.Sprintf("%s (%d)", limit, b.Status)
	}
	held := "the reset the API reported"
	if b.Reset.IsZero() || !b.Reset.Equal(b.Until) {
		held = "a backoff of " + b.backoff().String() + ", the API naming no reset"
	}
	return fmt.Sprintf("provider %s, on %s, held until %s, %s", b.Provider, limit, b.Until.UTC().Format(time.RFC3339), held)
}

// detail says what ch changed, for the journal.
func (ch budgetChange) detail() string {
	run := fmt.Sprintf("%s's run started at %s in %s", ch.Run.Work, ch.Run.Started, ch.Run.Root)
	switch ch.Kind {
	case budgetSet:
		return fmt.Sprintf("set the budget state: %s ended on the usage limit; every start on %s waits; %s", run, ch.State.Provider, budgetSaid(ch.State))
	case budgetExtended:
		return fmt.Sprintf("extended the budget state: %s ended on the usage limit again; %s", run, budgetSaid(ch.State))
	case budgetProbing:
		return fmt.Sprintf("started a probe of the budget state: %s goes through, every other start waits until it ends; %s", run, budgetSaid(ch.State))
	case budgetReleased:
		return fmt.Sprintf("freed the budget state's probe: %s could not be started, so the next start probes; %s", run, budgetSaid(ch.State))
	}
	return fmt.Sprintf("cleared the budget state: %s; starts on %s go again; it was %s", ch.Why, ch.State.Provider, budgetSaid(ch.State))
}

// budgetChanged announces ch: it journals it, logs it, and, for a set, an
// extension, or a clear, sends its notice (budgetNotice). Nothing when ch
// changed nothing.
func (l *launcher) budgetChanged(ch budgetChange) {
	if ch.Kind == "" {
		return
	}
	if l.record != nil {
		l.record(hostapi.Entry{At: l.now().UTC().Format(time.RFC3339), Action: hostapi.ActionAgent, Method: "serve.budget", Project: l.entry.Key, Root: l.entry.Root,
			By: "flai serve", Outcome: "done", Detail: ch.detail()})
	}
	l.log("budget state "+ch.Kind, "provider", ch.State.Provider, "limit", ch.State.Limit, "status", ch.State.Status,
		"until", ch.State.Until.UTC().Format(time.RFC3339), "backoff", ch.State.Backoff, "run", ch.Run.Work, "started", ch.Run.Started)
	l.budgetNotice(ch)
}

// budgetNotifier is how one serve folder's budget notices are sent: the
// channel of ADR-0205, read at each change so that the operator's flai serve
// notify set is seen without a restart, nil being no channel, and the
// sender. The notice is the host's, so the notifier is the serve folder's,
// not a project's: the flai serve of the folder sets it once with
// Dir.notifyBudget, and whichever project's launcher makes a change sends
// through it.
type budgetNotifier struct {
	channel func() config.Notify
	send    func(ctx context.Context, ch config.Notify, n BudgetNotice) error
}

// budgetNotifiers are the notifiers set in this process, by serve folder,
// and budgetUnchannelled the serve folders whose budget notices this process
// has warned have no channel, so that it warns once per run.
var (
	budgetNotifiersMu  sync.Mutex
	budgetNotifiers    = map[Dir]budgetNotifier{}
	budgetUnchannelled = map[Dir]bool{}
)

// notifyBudget has d's budget notices sent through channel with SendBudget;
// nil sends none, and they are logged instead.
func (d Dir) notifyBudget(channel func() config.Notify) {
	d.notifyBudgetWith(budgetNotifier{channel: channel, send: SendBudget})
}

// notifyBudgetWith has d's budget notices sent by n; a nil channel sends
// none.
func (d Dir) notifyBudgetWith(n budgetNotifier) {
	budgetNotifiersMu.Lock()
	defer budgetNotifiersMu.Unlock()
	if n.channel == nil {
		delete(budgetNotifiers, d)
		return
	}
	budgetNotifiers[d] = n
}

// budgetChannel is d's notifier and the channel it reads now, and, when the
// channel is not set, whether this is the first time this process finds so.
func (d Dir) budgetChannel() (n budgetNotifier, ch config.Notify, firstUnset bool) {
	budgetNotifiersMu.Lock()
	n, ok := budgetNotifiers[d]
	budgetNotifiersMu.Unlock()
	if ok {
		ch = n.channel()
	}
	if ch.URL != "" || strings.TrimSpace(ch.Command) != "" {
		return n, ch, false
	}
	budgetNotifiersMu.Lock()
	defer budgetNotifiersMu.Unlock()
	firstUnset = !budgetUnchannelled[d]
	budgetUnchannelled[d] = true
	return n, ch, firstUnset
}

// budgetNoticeStates are the changes a notice is sent of; a probe's start
// and its release are not.
var budgetNoticeStates = map[string]bool{budgetSet: true, budgetExtended: true, budgetCleared: true}

// budgetNotice sends the notice of ch, a set, an extension, or a clear,
// through the serve folder's channel, once: the one launcher whose run made
// the change sends it, whichever project it serves. One attempt within the
// reminder's timeouts, its failure warned of. With no channel set the change
// is logged (budgetChanged) and that none is set is warned of once per run.
func (l *launcher) budgetNotice(ch budgetChange) {
	if !budgetNoticeStates[ch.Kind] {
		return
	}
	n, channel, firstUnset := l.dir.budgetChannel()
	if channel.URL == "" && strings.TrimSpace(channel.Command) == "" {
		if firstUnset {
			l.warn("budget notices have no channel", "hint", "set one in a shell on the host with flai serve notify set --url <webhook> or --command <command>")
		}
		return
	}
	notice := BudgetNotice{State: ch.Kind, Provider: ch.State.Provider, Limit: ch.State.Limit, Status: ch.State.Status,
		Since: ch.State.Since, Until: ch.State.Until, At: l.now(), Projects: l.budgetProjects(), Href: BudgetHref}
	attrs := []any{"state", notice.State, "provider", notice.Provider, "until", notice.Until.UTC().Format(time.RFC3339)}
	if err := n.send(context.Background(), channel, notice); err != nil {
		l.warn("budget notice not sent", append(attrs, "err", err.Error())...)
		return
	}
	l.log("budget notice sent", attrs...)
}

// budgetProjects are the names of the projects in the serve folder's
// registry, by key when one has none; none when it cannot be read, which is
// warned of.
func (l *launcher) budgetProjects() []string {
	all, err := l.dir.Projects()
	if err != nil {
		l.warn("budget notice names no project: the registry is not read", "err", err)
		return nil
	}
	out := make([]string, 0, len(all))
	for _, e := range all {
		out = append(out, cmp.Or(e.Name, e.Key))
	}
	return out
}

// budgetStart is the gate run's start passes once its harness and provider
// are resolved and before it spawns (budgetStart on the Dir), announcing
// what it changed. A budget state that cannot be read is warned of and the
// start goes: a file in error must not hold every start on the host.
func (l *launcher) budgetStart(run *AgentRun) budgetGate {
	gate, ch, err := l.dir.budgetStart(budgetKey(run.Provider, run.Harness), probeOf(l.entry.Root, run), l.now())
	if err != nil {
		l.warn("budget state not read before a start; it goes", "run", probeOf(l.entry.Root, run).Work, "err", err)
		return budgetGate{Go: true}
	}
	l.budgetChanged(ch)
	return gate
}

// budgetSpawned records run's process on the probe, when run is the probe,
// once it is spawned; an error is warned of.
func (l *launcher) budgetSpawned(run *AgentRun) {
	if err := l.dir.budgetSpawned(budgetKey(run.Provider, run.Harness), probeOf(l.entry.Root, run)); err != nil {
		l.warn("budget state's probe not recorded with its process", "run", probeOf(l.entry.Root, run).Work, "err", err)
	}
}

// budgetRelease frees the probe for run, whose start could not be made,
// when run is the probe, and announces it; an error is warned of.
func (l *launcher) budgetRelease(run *AgentRun) {
	ch, err := l.dir.budgetRelease(budgetKey(run.Provider, run.Harness), probeOf(l.entry.Root, run))
	if err != nil {
		l.warn("budget state's probe not freed", "run", probeOf(l.entry.Root, run).Work, "err", err)
		return
	}
	l.budgetChanged(ch)
}

// budgetEnded records run's end, on lim, the usage limit it ended on, or nil
// for none, in the budget state (budgetEnded on the Dir), and announces what
// it changed; an error is warned of.
func (l *launcher) budgetEnded(run *AgentRun, lim *usage.Limit) {
	ch, err := l.dir.budgetEnded(budgetKey(run.Provider, run.Harness), probeOf(l.entry.Root, run), lim, l.now())
	if err != nil {
		l.warn("budget state not recorded at a run's end", "run", probeOf(l.entry.Root, run).Work, "err", err)
		return
	}
	l.budgetChanged(ch)
}

// budgetKeyFor is the key a start of agent under cfg is held by (budgetKey):
// the provider agent names, or the harness it would run on, as harness.For
// picks it; "" when no harness can start it, which the start refuses on its
// own.
func budgetKeyFor(cfg AgentConfig, agent *manifest.Agent) string {
	name, _, err := harness.For(agent, cfg.host(harness.Command).Program != "")
	if err != nil {
		return ""
	}
	provider := ""
	if agent != nil {
		provider = agent.Provider
	}
	return budgetKey(provider, name)
}

// budgetHolds is the budget state of key, and whether it holds a start at
// the launcher's clock; a state that cannot be read is warned of and holds
// nothing, as at a start (budgetStart), and an empty key holds nothing.
func (l *launcher) budgetHolds(key string) (Budget, bool) {
	if key == "" {
		return Budget{}, false
	}
	b, held, err := l.dir.budgetHeld(key, l.now())
	if err != nil {
		l.warn("budget state not read before a start; it goes", "provider", key, "err", err)
		return Budget{}, false
	}
	return b, held
}

// budgetHeldWhy says why a start on b's provider is held, as a start the
// operator asks for by hand is refused and a held start is logged: such as
// "held: claude-code usage limit (usage_limit_reached) until
// 2026-10-28T09:00:00Z; flai serve starts it when the limit resets", or,
// once until has passed, that the probe runs.
func budgetHeldWhy(b Budget, now time.Time) string {
	limit := b.Limit
	if limit == "" && b.Status != 0 {
		limit = fmt.Sprintf("HTTP %d", b.Status)
	}
	if limit != "" {
		limit = " (" + limit + ")"
	}
	until := b.Until.UTC().Format(time.RFC3339)
	if !now.Before(b.Until) {
		return fmt.Sprintf("held: %s usage limit%s, a probe running since %s to learn whether it has reset; flai serve starts it when the probe ends clean", b.Provider, limit, until)
	}
	return fmt.Sprintf("held: %s usage limit%s until %s; flai serve starts it when the limit resets", b.Provider, limit, until)
}
