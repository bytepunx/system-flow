package serve

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"sort"
	"strings"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/atomicfile"
	"github.com/bytepunx/system-flow/flai/internal/harness"
	"github.com/bytepunx/system-flow/flai/internal/manifest"
	"github.com/bytepunx/system-flow/flai/internal/usage"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// What agents spent (S-0143, ADR-0051). flai serve keeps the log of every
// agent it starts, named for the project's key, the story, and the time it
// started. From those logs, read by the reader of the harness the story's
// agent ran with (ADR-0130), it measures a story whole, and each of the
// story's tasks from the calls of the sub-agents started for it and its
// share of the rest over the intervals it was in progress (S-0230), and
// writes their usage into the front matter: when an agent it started ends,
// and when a task of a story whose agent runs enters done without having
// been measured. flai serve agent usage does the same on the operator's word.

// AgentLogs are the logs flai serve keeps of the agents it started for
// story in the project named key, oldest first.
func (d Dir) AgentLogs(key, story string) ([]string, error) {
	if !StoryID.MatchString(story) {
		return nil, fmt.Errorf("%q is not a story's ID", story)
	}
	paths, err := filepath.Glob(filepath.Join(string(d), "agents", key+"-"+story+"-*.log"))
	sort.Strings(paths)
	return paths, err
}

// LoggedStories are the stories of the project named key that flai serve
// has kept a log for, in order of ID.
func (d Dir) LoggedStories(key string) ([]string, error) {
	paths, err := filepath.Glob(filepath.Join(string(d), "agents", key+"-S-*.log"))
	if err != nil {
		return nil, err
	}
	seen := map[string]bool{}
	var out []string
	for _, p := range paths {
		name := filepath.Base(p)[len(key)+1:]
		for i := 2; i < len(name); i++ {
			if name[i] == '-' {
				if id := name[:i]; StoryID.MatchString(id) && !seen[id] {
					seen[id] = true
					out = append(out, id)
				}
				break
			}
		}
	}
	sort.Slice(out, func(i, j int) bool { return workitem.CanonicalID(out[i]) < workitem.CanonicalID(out[j]) })
	return out, nil
}

// Rates are what each model's tokens cost across every log flai serve
// keeps, to price what no result reported.
func (d Dir) Rates() (usage.Rates, error) {
	paths, err := filepath.Glob(filepath.Join(string(d), "agents", "*.log"))
	if err != nil {
		return nil, err
	}
	// claude-code is the harness that reports what its calls cost
	return usage.ReadRates(harness.Adapters[harness.ClaudeCode].Reader(), paths...)
}

// logReader is the reader of the logs of story's agent: that of the harness
// its newest run recorded, or when none did, of the harness its agent names,
// as harness.For resolves it, the operator's command when it names none.
func logReader(d Dir, root string, story *workitem.Item) (usage.Reader, error) {
	a := story.Agent
	if run := d.AgentStates()[root].Stories[story.ID]; run != nil && run.Harness != "" {
		a = &manifest.Agent{Harness: run.Harness}
	}
	_, adapter, err := harness.For(a, true)
	if err != nil {
		return nil, fmt.Errorf("cannot read the logs of %s's agent: %w", story.ID, err)
	}
	return adapter.Reader(), nil
}

// Measured is what the logs say a story and its tasks spent.
type Measured struct {
	Story string                  `json:"story"`
	Logs  []string                `json:"logs"`
	Usage *usage.Usage            `json:"usage,omitempty"`
	Tasks map[string]*usage.Usage `json:"tasks"`
	// Changed are the items whose usage was written, their epic's among them.
	Changed []string `json:"changed"`
	// Unread are the runs through a provider whose gateway did not say what
	// it charged, priced by the harness or estimated instead (ADR-0132).
	Unread []Unread `json:"unread,omitempty"`
}

// Unread is a run through a provider whose cost its gateway's spend log did
// not give, and why; the next measurement asks again.
type Unread struct {
	Log      string `json:"log"`
	Provider string `json:"provider"`
	// Reason is Err's text, which names a variable at most, never a key.
	Reason string `json:"reason"`
	Err    error  `json:"-"`
}

// Measure measures story in the project at root, named key, from the logs
// of the agents flai serve started for it: the story from all of them, each
// task from its sub-agents' calls and its share of the rest over the
// intervals it was in progress. A run through a provider is priced by its
// gateway's spend log, read with the key of the host's entry in providers,
// as price says (ADR-0132). With write, it records what differs from what
// the items say, and sums the story's epic again. A task that has been done
// without anything measured for it is written as having spent nothing, so
// that it is not measured at every look.
func Measure(d Dir, root, key, story string, providers map[string]manifest.Provider, write bool) (*Measured, error) {
	logs, err := d.AgentLogs(key, story)
	if err != nil {
		return nil, err
	}
	out := &Measured{Story: story, Logs: logs, Tasks: map[string]*usage.Usage{}, Changed: []string{}}
	if logs == nil {
		out.Logs = []string{}
	}
	repo, err := workitem.Open(root)
	if err != nil {
		return nil, err
	}
	it, err := repo.Get(story)
	if err != nil {
		return nil, err
	}
	if len(logs) == 0 {
		return out, nil
	}
	r, err := logReader(d, root, it)
	if err != nil {
		return nil, err
	}
	rec, err := usage.Read(r, logs...)
	if err != nil {
		return nil, err
	}
	rates, err := d.Rates()
	if err != nil {
		return nil, err
	}
	out.Usage = rec.Total(rates)
	p, err := price(rec, rates, d.AgentStates()[root].Stories[story], providers, repo.Manifest.Providers)
	if err != nil {
		return nil, err
	}
	out.Unread = p.unread
	ratio := p.apply(out.Usage)
	items, err := repo.List(true)
	if err != nil {
		return nil, err
	}
	tasks := workitem.Children(items, story)
	spans := map[string][]usage.Span{}
	for _, t := range tasks {
		spans[t.ID] = inProgress(t)
	}
	measured := rec.Tasks(spans, rates)
	for _, u := range measured {
		share(u, ratio, out.Usage.Priced())
	}
	for _, t := range tasks {
		u := measured[t.ID]
		if u == nil && t.Status == workitem.Done && len(spans[t.ID]) > 0 {
			u = &usage.Usage{Source: usage.SourceLog, Models: []usage.Model{}}
		}
		if u != nil {
			out.Tasks[t.ID] = u
		}
	}
	if !write {
		return out, nil
	}
	put := func(id string, u *usage.Usage) error {
		// read again just before writing: an agent may have changed it
		fresh, err := repo.Get(id)
		if err != nil {
			return err
		}
		if u == nil {
			return nil
		}
		// what strategic agents spent on it is no agent's: keep it (ADR-0083)
		u = usage.WithStrategic(u, fresh.Usage)
		if fresh.Usage != nil && usage.Same(fresh.Usage, u) {
			return nil
		}
		fresh.Usage = u
		if err := repo.Save(fresh); err != nil {
			return err
		}
		out.Changed = append(out.Changed, id)
		return nil
	}
	ids := make([]string, 0, len(out.Tasks))
	for id := range out.Tasks {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	for _, id := range ids {
		if err := put(id, out.Tasks[id]); err != nil {
			return out, err
		}
	}
	if err := put(story, out.Usage); err != nil {
		return out, err
	}
	it, err = repo.Get(story)
	if err != nil {
		return out, err
	}
	up, err := repo.RollUp(it)
	out.Changed = append(out.Changed, up...)
	return out, err
}

// inProgress are the intervals an item was in progress, the last open when
// it still is.
func inProgress(it *workitem.Item) []usage.Span {
	var out []usage.Span
	for i, tr := range it.Transitions {
		if tr.To != workitem.InProgress {
			continue
		}
		from, err := time.Parse(workitem.TimeFormat, tr.At)
		if err != nil {
			continue
		}
		s := usage.Span{From: from}
		if i+1 < len(it.Transitions) {
			s.To, _ = time.Parse(workitem.TimeFormat, it.Transitions[i+1].At)
		}
		out = append(out, s)
	}
	return out
}

// measure measures story and says so in flai serve's log, with each run
// through a provider whose gateway did not say what it charged.
func (l *launcher) measure(story string) {
	var providers map[string]manifest.Provider
	if l.config != nil {
		providers = l.config(l.entry.Root).Providers
	}
	l.measureWith(story, providers)
}

// measureWith is measure with the host's providers given: an agent's end
// prices its runs with those its start read, as its own goroutine does not
// read the configuration again.
func (l *launcher) measureWith(story string, providers map[string]manifest.Provider) {
	m, err := Measure(l.dir, l.entry.Root, l.entry.Key, story, providers, true)
	if err != nil {
		l.warn("usage not measured", "story", story, "err", err)
		return
	}
	for _, u := range m.Unread {
		l.warn("spend not read", "story", story, "provider", u.Provider, "log", filepath.Base(u.Log), "err", u.Err)
	}
	if len(m.Changed) > 0 {
		l.log("usage measured", "story", story, "logs", len(m.Logs), "tokens", m.Usage.Tokens(), "cost", fmt.Sprintf("%.2f", m.Usage.Cost()), "priced_by", m.Usage.Priced(), "changed", len(m.Changed))
	}
}

// measureDone measures the story of each agent that runs when one of its
// tasks is done and has not been measured.
func (l *launcher) measureDone(st AgentState) {
	var running []string
	for id, run := range st.Stories {
		if run.live() {
			running = append(running, id)
		}
	}
	if len(running) == 0 {
		return
	}
	repo, err := workitem.Open(l.entry.Root)
	if err != nil {
		return
	}
	items, err := repo.List(false)
	if err != nil {
		return
	}
	sort.Strings(running)
	for _, id := range running {
		for _, t := range workitem.Children(items, id) {
			if t.Status == workitem.Done && t.Usage == nil {
				l.measure(id)
				break
			}
		}
	}
}

// What a provider's gateway charged (S-0357, ADR-0132). A story's agent's run
// through a provider keeps a spend record beside its log: the provider's
// name and, where its gateway's reader joins a run by the key's usage, that
// usage when the run started. A measurement asks the gateway for each such
// run once it has ended, keeps its answer in the record, and asks no more;
// a gateway that does not answer leaves the run priced by the harness, or
// estimated, and is asked again at the next measurement.

// spendWait bounds what flai serve waits for a gateway's spend log: the
// key's usage at a run's start, and a run's cost when it measures it.
const spendWait = 10 * time.Second

// spendFor is the spend reader of a gateway at baseURL, read with key:
// usage.For's, which tests replace.
var spendFor = func(baseURL, key string) usage.Spend { return usage.For(baseURL, key, nil) }

// spendRecord is what flai serve keeps of a run through a provider, beside
// its log. It never holds a key.
type spendRecord struct {
	// Provider is the name of the provider the run went through.
	Provider string `json:"provider"`
	// KeyBefore is the key's usage in US dollars when the run started, read
	// where the gateway's reader is a usage.KeyReader.
	KeyBefore *float64 `json:"key_before,omitempty"`
	// Cost, By, PricedBy, and At are set once the gateway priced the run:
	// what it charged in US dollars, how its reader joined it to the run (a
	// usage.SpentBy constant), usage.PricedGateway, and when.
	Cost     *float64 `json:"cost,omitempty"`
	By       string   `json:"by,omitempty"`
	PricedBy string   `json:"priced_by,omitempty"`
	At       string   `json:"at,omitempty"`
}

// priced says the gateway priced the run.
func (r *spendRecord) priced() bool { return r != nil && r.Cost != nil && r.PricedBy != "" }

// spendPath is the spend record of the run whose log is at log: beside it,
// named so that the agents folder's *.log does not read it as a log.
func spendPath(log string) string { return log + ".spend.json" }

// readSpend is the spend record of the run whose log is at log; nil when it
// has none.
func readSpend(log string) (*spendRecord, error) {
	data, err := os.ReadFile(spendPath(log))
	if errors.Is(err, os.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var r spendRecord
	if err := json.Unmarshal(data, &r); err != nil {
		return nil, fmt.Errorf("read the spend record %s: %w; remove it, and a run that names a provider in flai serve's state is priced again", spendPath(log), err)
	}
	return &r, nil
}

// writeSpend writes r as the spend record of the run whose log is at log.
func writeSpend(log string, r *spendRecord) error {
	data, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		return err
	}
	return atomicfile.WriteFile(spendPath(log), append(data, '\n'), 0o600)
}

// spendOf is the spend reader of the provider name, resolved from the host's
// providers with the project's overrides as a start resolves it, read with
// the key in the variable its spend_key_env names, else in its key_env's;
// and the SHA-256 in hex of the agent's key, its key_env's, empty when that
// is unset. The reader is nil for a gateway flai reads no spend log of. It
// refuses a spend key's variable that is unset or empty, naming it.
func spendOf(name string, host, project map[string]manifest.Provider) (usage.Spend, string, error) {
	p, err := manifest.ResolveProvider(name, host, project)
	if err != nil {
		return nil, "", err
	}
	v := p.SpendKeyEnv
	if v == "" {
		v = p.KeyEnv
	}
	key, ok := os.LookupEnv(v)
	if !ok || key == "" {
		return nil, "", fmt.Errorf("provider %[1]s's spend log is read with the key in %[2]s, which is unset or empty in flai serve's environment: set %[2]s where flai serve runs, or name the variable of a key that can read the gateway's spend log with flai serve agent provider %[1]s --spend-key-env <variable>", name, v)
	}
	var hash string
	if agent, ok := os.LookupEnv(p.KeyEnv); ok && agent != "" {
		sum := sha256.Sum256([]byte(agent))
		hash = hex.EncodeToString(sum[:])
	}
	return spendFor(p.BaseURL, key), hash, nil
}

// spendAtStart is the spend record of run, a story's agent's run whose
// provider startOf resolved, with the key's usage read now where its
// gateway's reader joins a run by it; nil when run names no provider. A read
// that fails is warned of, and the record goes without it.
func (l *launcher) spendAtStart(ctx context.Context, cfg AgentConfig, run *AgentRun) *spendRecord {
	if run.Provider == "" {
		return nil
	}
	rec := &spendRecord{Provider: run.Provider}
	m, err := manifest.Load(filepath.Join(l.entry.Root, manifest.File))
	if err != nil {
		l.warn("key usage not read", "story", run.Story, "provider", run.Provider, "err", err)
		return rec
	}
	spend, _, err := spendOf(run.Provider, cfg.Providers, m.Providers)
	if err != nil {
		l.warn("key usage not read", "story", run.Story, "provider", run.Provider, "err", err)
		return rec
	}
	keys, ok := spend.(usage.KeyReader)
	if !ok {
		return rec
	}
	ctx, cancel := context.WithTimeout(ctx, spendWait)
	defer cancel()
	v, err := keys.KeyUsage(ctx)
	if err != nil {
		l.warn("key usage not read", "story", run.Story, "provider", run.Provider, "err", err)
		return rec
	}
	rec.KeyBefore = &v
	return rec
}

// keepSpend writes rec as the spend record of run, once its log is known;
// nothing for no record. A write that fails is warned of: the run is then
// priced as flai serve's state names its provider, without the key's usage
// at its start.
func (l *launcher) keepSpend(run *AgentRun, rec *spendRecord) {
	if rec == nil || run.Log == "" {
		return
	}
	if err := writeSpend(run.Log, rec); err != nil {
		l.warn("spend record not written", "story", run.Story, "provider", rec.Provider, "err", err)
	}
}

// pricing is what a story's runs were priced at.
type pricing struct {
	// cost is each model's cost over the runs: the gateway's where it priced
	// a run, the harness's figures otherwise.
	cost map[string]float64
	// gateway says a gateway priced some run.
	gateway bool
	// pricedBy is the least certain source of the runs' costs, and estimated
	// says some of the cost was estimated or apportioned.
	pricedBy  string
	estimated bool
	unread    []Unread
}

// price prices each run of rec on its own, from what it added to the
// record's total (usage.Record.RunTotals): a run that spent nothing is not
// priced; one whose spend record names a provider, or the newest, run,
// when it is the run's and names one, is priced by the provider's gateway,
// its cost apportioned over its models by their share of its cost, or of
// their tokens when that is nothing; any other, and one whose gateway did
// not answer, by the harness, or estimated when any of its cost was. The run
// newest names is not asked for while it goes on. It fails only when a
// gateway's price cannot be kept.
func price(rec *usage.Record, rates usage.Rates, newest *AgentRun, host, project map[string]manifest.Provider) (*pricing, error) {
	p := &pricing{cost: map[string]float64{}}
	totals := rec.RunTotals(rates)
	for i, run := range rec.Runs {
		u := totals[i]
		weights, sum := shares(u)
		if sum == 0 {
			continue
		}
		priced, err := gatewayOf(run, newest, host, project)
		var unread *unreadError
		switch {
		case errors.As(err, &unread):
			p.unread = append(p.unread, unread.Unread)
		case err != nil:
			return nil, err
		}
		if priced == nil {
			for _, m := range u.Models {
				p.cost[m.Model] += m.Cost
			}
			by := usage.PricedHarness
			if u.Estimated {
				by = usage.PricedEstimate
			}
			p.pricedBy = usage.LeastCertain(p.pricedBy, by)
			p.estimated = p.estimated || u.Estimated
			continue
		}
		over := 0
		for name, w := range weights {
			if w > 0 {
				p.cost[name] += *priced.Cost * w / sum
				over++
			}
		}
		p.gateway = true
		p.pricedBy = usage.LeastCertain(p.pricedBy, usage.PricedGateway)
		p.estimated = p.estimated || over > 1
	}
	return p, nil
}

// shares are what a run's gateway cost is apportioned by over its models:
// each model's cost, or its tokens when the run's cost is nothing; and
// their sum, nothing when the run spent nothing.
func shares(u *usage.Usage) (map[string]float64, float64) {
	for _, tokens := range []bool{false, true} {
		w, sum := map[string]float64{}, 0.0
		for _, m := range u.Models {
			v := m.Cost
			if tokens {
				v = float64(m.Tokens())
			}
			if v > 0 {
				w[m.Model] = v
				sum += v
			}
		}
		if sum > 0 {
			return w, sum
		}
	}
	return nil, 0
}

// gatewayOf is the spend record of run once its gateway priced it, asking
// the gateway when the record names a provider and has no price, and keeping
// its answer; nil when the run went through no provider, still goes on, or
// its gateway did not answer, which its error, an *unreadError, says.
func gatewayOf(run usage.Run, newest *AgentRun, host, project map[string]manifest.Provider) (*spendRecord, error) {
	r, err := readSpend(run.Path)
	if err != nil {
		return nil, &unreadError{Unread{Log: run.Path, Reason: err.Error(), Err: err}}
	}
	ours := newest != nil && newest.Log == run.Path
	if r == nil && ours && newest.Provider != "" {
		// started after S-0350 and before its run kept a spend record
		r = &spendRecord{Provider: newest.Provider}
	}
	switch {
	case r == nil:
		return nil, nil
	case r.priced():
		return r, nil
	case ours && newest.live():
		return nil, nil
	}
	spent, err := askSpend(r, run, host, project)
	if errors.Is(err, errNoSpendLog) {
		// nothing to read again: the harness's figure is the price
		return nil, nil
	}
	if err != nil {
		return nil, &unreadError{Unread{Log: run.Path, Provider: r.Provider, Reason: err.Error(), Err: err}}
	}
	r.Cost, r.By, r.PricedBy, r.At = &spent.Cost, spent.By, usage.PricedGateway, time.Now().UTC().Format(time.RFC3339)
	if err := writeSpend(run.Path, r); err != nil {
		return nil, fmt.Errorf("keep what provider %s's gateway charged for the run of %s: %w", r.Provider, filepath.Base(run.Path), err)
	}
	return r, nil
}

// unreadError is a run's gateway not saying what it charged: the run is
// priced without it.
type unreadError struct{ Unread }

func (e *unreadError) Error() string { return e.Reason }

func (e *unreadError) Unwrap() error { return e.Err }

// errNoSpendLog is a provider whose gateway keeps no spend log flai reads,
// such as Anthropic's own API: its runs are priced by the harness, quietly.
var errNoSpendLog = errors.New("the gateway keeps no spend log flai reads")

// askSpend asks the gateway of r's provider what it charged for run.
func askSpend(r *spendRecord, run usage.Run, host, project map[string]manifest.Provider) (*usage.Spent, error) {
	spend, hash, err := spendOf(r.Provider, host, project)
	if err != nil {
		return nil, err
	}
	if spend == nil {
		return nil, fmt.Errorf("provider %s: %w", r.Provider, errNoSpendLog)
	}
	ctx, cancel := context.WithTimeout(context.Background(), spendWait)
	defer cancel()
	spent, err := spend.Run(ctx, usage.RunRef{Sessions: run.Sessions, From: run.Start, To: run.End, KeyBefore: r.KeyBefore, KeyHash: hash})
	if err == nil && spent == nil {
		err = fmt.Errorf("provider %s's gateway answered nothing for the run: %w", r.Provider, usage.ErrNoSpend)
	}
	return spent, err
}

// apply prices u, the story's usage, as its runs were: each model's cost
// what the runs came to where a gateway priced any, and where its cost came
// from. It is the ratio of each model's priced cost to the harness's, where
// the harness's is something, for the tasks' shares; nil when no gateway
// priced a run.
func (p *pricing) apply(u *usage.Usage) map[string]float64 {
	if u == nil {
		return nil
	}
	if !p.gateway {
		u.PricedBy = usage.PricedHarness
		if u.Estimated {
			u.PricedBy = usage.PricedEstimate
		}
		return nil
	}
	ratio := map[string]float64{}
	for i := range u.Models {
		m := &u.Models[i]
		if m.Cost > 0 {
			ratio[m.Model] = p.cost[m.Model] / m.Cost
		}
		m.Cost = p.cost[m.Model]
	}
	u.PricedBy, u.Estimated = p.pricedBy, p.estimated
	u.Tidy()
	return ratio
}

// share prices a task's share of its story as the story was priced: each
// model's cost scaled by ratio's for it, and its cost's source the story's.
// It stays estimated, being apportioned.
func share(u *usage.Usage, ratio map[string]float64, pricedBy string) {
	for i := range u.Models {
		if r, ok := ratio[u.Models[i].Model]; ok {
			u.Models[i].Cost *= r
		}
	}
	u.PricedBy = pricedBy
	u.Tidy()
}

// withoutSpendKeys is env without the variables the host's providers name
// as spend_key_env, which hold keys that read a gateway's spend log, often
// its master key: the agents flai serve starts never get one. A variable
// that is also a provider's key_env stays, since an agent reads its key
// from it.
func (l *launcher) withoutSpendKeys(env []string) []string {
	if l.config == nil {
		return env
	}
	providers := l.config(l.entry.Root).Providers
	spend := map[string]bool{}
	for _, p := range providers {
		if p.SpendKeyEnv != "" {
			spend[p.SpendKeyEnv] = true
		}
	}
	for _, p := range providers {
		delete(spend, p.KeyEnv)
	}
	if len(spend) == 0 {
		return env
	}
	return slices.DeleteFunc(env, func(kv string) bool {
		name, _, _ := strings.Cut(kv, "=")
		return spend[name]
	})
}
