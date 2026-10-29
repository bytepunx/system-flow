// Package serve is flai's process on the host that dashboards are reached
// through (ADR-0029): one per user, serving every project registered with
// it, holding one connection to each project's dashboard. Its registry and
// its state live beside flai's config file, so a test or a repository with
// its own config never touches the operator's.
package serve

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"maps"
	"os"
	"path/filepath"
	"slices"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/buildinfo"
	"github.com/bytepunx/system-flow/flai/internal/channel"
	"github.com/bytepunx/system-flow/flai/internal/config"
	"github.com/bytepunx/system-flow/flai/internal/conventions"
	"github.com/bytepunx/system-flow/flai/internal/dashboard"
	"github.com/bytepunx/system-flow/flai/internal/execx"
	"github.com/bytepunx/system-flow/flai/internal/hostapi"
	"github.com/bytepunx/system-flow/flai/internal/manifest"
	"github.com/bytepunx/system-flow/flai/internal/release"
	"github.com/bytepunx/system-flow/flai/internal/watch"
)

// Entry is one registered project: where it is, and the dashboard to dial.
type Entry struct {
	Key     string `json:"key"`
	Name    string `json:"name"`
	Root    string `json:"root"`
	URL     string `json:"url"`
	KeyFile string `json:"key_file"` // the agent credential, mode 0600
}

// Status is what a running flai serve last wrote about itself.
type Status struct {
	PID         int                      `json:"pid"`
	Version     string                   `json:"version"`
	Started     string                   `json:"started"`
	Updated     string                   `json:"updated"`
	Connections map[string]channel.State `json:"connections"` // by project root
	// Offered are the repositories offered for import (S-0098).
	Offered []Candidate `json:"offered,omitempty"`
	// Folder is the folder flai serve was started in when it is not a project,
	// and FolderProjects the projects below it that it serves (S-0102).
	Folder         string  `json:"folder,omitempty"`
	FolderProjects []Entry `json:"folder_projects,omitempty"`
	// ImportProjects are the git repositories with a system-flow.yaml below
	// the folders named for import that it serves, and Unserved the projects
	// below either kind of folder that it does not, with why (S-0120).
	ImportProjects []Entry `json:"import_projects,omitempty"`
	Unserved       []Found `json:"unserved,omitempty"`
	// Unavailable are the registered projects not served, by root, and why
	// (S-0121): the folder is gone or has no manifest, say.
	Unavailable map[string]string `json:"unavailable,omitempty"`
	// Images are each served project's last verdict on the image its
	// dashboard container runs, by root, as Connections are (S-0238): listed,
	// unlisted and so not dialled, or not measured, with why.
	Images map[string]ImageVerdict `json:"images,omitempty"`
}

// Dir is the directory that holds the registry and the state.
type Dir string

// DirFor is serve's directory for a config file: a folder beside it.
func DirFor(configPath string) Dir {
	return Dir(filepath.Join(filepath.Dir(configPath), "serve"))
}

func (d Dir) registry() string { return filepath.Join(string(d), "projects.json") }
func (d Dir) status() string   { return filepath.Join(string(d), "state.json") }

// Log is where a detached flai serve writes.
func (d Dir) Log() string { return filepath.Join(string(d), "serve.log") }

// Journal is where every host action asked for is written, one JSON object a line.
func (d Dir) Journal() string { return filepath.Join(string(d), "journal.jsonl") }

// Projects reads the registry; none is an empty list.
func (d Dir) Projects() ([]Entry, error) {
	data, err := os.ReadFile(d.registry())
	if errors.Is(err, os.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var out []Entry
	if err := json.Unmarshal(data, &out); err != nil {
		return nil, fmt.Errorf("%s: %w", d.registry(), err)
	}
	return out, nil
}

func (d Dir) write(path string, v any) error {
	if err := os.MkdirAll(string(d), 0o700); err != nil {
		return err
	}
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, append(data, '\n'), 0o600); err != nil {
		return err
	}
	return os.Rename(tmp, path)
}

// Register adds the project or replaces the entry with its root. A key
// another root has is refused (S-0121).
func (d Dir) Register(e Entry) error {
	if e.Root == "" || e.URL == "" || e.KeyFile == "" || e.Key == "" {
		return errors.New("a project needs a root, a key, a dashboard address, and a credential file")
	}
	all, err := d.Projects()
	if err != nil {
		return err
	}
	if err := keyTaken(all, e.Key, e.Root); err != nil {
		return err
	}
	out := []Entry{}
	for _, x := range all {
		if x.Root != e.Root {
			out = append(out, x)
		}
	}
	out = append(out, e)
	sort.Slice(out, func(i, j int) bool { return out[i].Root < out[j].Root })
	return d.write(d.registry(), out)
}

// Unregister removes the project with this root; absent is not an error.
func (d Dir) Unregister(root string) error {
	all, err := d.Projects()
	if err != nil {
		return err
	}
	out := []Entry{}
	for _, x := range all {
		if x.Root != root {
			out = append(out, x)
		}
	}
	return d.write(d.registry(), out)
}

// StaleAfter is how old a status may be before its writer counts as gone.
const StaleAfter = 10 * time.Second

// ReadStatus returns the last status and whether its process is still
// running: the PID is alive and the status is fresh.
func (d Dir) ReadStatus(now time.Time) (Status, bool) {
	var st Status
	data, err := os.ReadFile(d.status())
	if err != nil || json.Unmarshal(data, &st) != nil {
		return Status{}, false
	}
	updated, err := time.Parse(time.RFC3339, st.Updated)
	if err != nil || now.Sub(updated) > StaleAfter {
		return st, false
	}
	return st, Alive(st.PID)
}

// Options configure Run.
type Options struct {
	Dir        Dir
	Version    string
	Logger     *slog.Logger
	Every      time.Duration // how often the registry is read and the status written
	Now        func() time.Time
	WatchEvery time.Duration // how often a project's files are looked at; the watcher's default when zero
	// Host is the operator's say over host actions and the journal of them
	// (S-0078); the zero value enables nothing.
	Host hostapi.Host
	// Agent is the operator's say about starting an agent when a story
	// becomes ready (S-0079); nil starts none.
	Agent func(root string) AgentConfig
	// StopGrace is how long Stop waits for an agent it asked to stop before
	// it kills it (S-0170); ten seconds when zero.
	StopGrace time.Duration
	// NewClient lets tests shorten a client's timings.
	NewClient func(e Entry, key []byte) *channel.Client
	// Upgrade asks the host flai serve runs under to install a newer flai
	// release and restart on it (ADR-0158), for a served project whose
	// history holds one while its auto-upgrade host action is on; nil when
	// it runs under no host, when it asks the operator on a thread instead.
	Upgrade HostUpgrade
	// Published says whether a flai release is published, asked before the
	// thread that names it is opened (ADR-0161); nil opens it unasked.
	Published ReleasePublished
	// MCP has each served project's HTTP MCP server kept by flai host
	// (S-0096, S-0106); nil keeps none. MCPEvery is how often the host is
	// told again when nothing changed.
	MCP      MCP
	MCPEvery time.Duration
	// ImportRoots are the folders the operator named for repositories to
	// import from the board (S-0098), asked at every look; nil offers none.
	// ScanEvery is how often they are looked through; ImportRun runs the
	// import (the flai that is running when nil).
	ImportRoots func() []string
	ScanEvery   time.Duration
	ImportRun   hostapi.Runner
	// Folder is the folder flai serve was started in when that is not a
	// project (S-0102): the projects below it are served, and it is looked
	// in for repositories to import, as if named. Empty otherwise.
	Folder string
	// Git asks git what a story's worktree holds, for flai serve agent
	// commit (S-0140); the system's git when nil.
	Git execx.Runner
	// DrainGrace is how long a project whose manifest stopped loading has to
	// answer the requests in flight before its connection ends (I-0107);
	// DefaultDrainGrace when zero.
	DrainGrace time.Duration
	// Images are the verified digest lists the image the dashboard container
	// runs is judged against before each dial of a project's dashboard
	// (S-0238); nil judges none. Docker runs docker for that check, the
	// system's docker when nil.
	Images *ImageLists
	Docker execx.Runner
	// AllowUnsigned is the host configuration's dashboard.allow_unsigned,
	// asked before each dial of a project's dashboard, which a project's
	// manifest overrides (S-0239); nil is off.
	AllowUnsigned func() bool
	// ForgeToken says how the token of a forge acceptance.forge names is
	// read (ADR-0165), for the push of each story in review under
	// pull-request acceptance (ADR-0166); nil reads none, so each such push
	// fails saying so. NewForge makes the forge's adapter, forge.New with
	// GitHub's API when nil.
	ForgeToken func(forge string) ForgeToken
	NewForge   NewForge
	// Notify is the host configuration's reminder channel (ADR-0205), asked
	// at each look of a project's thread reminders, so that the operator's
	// flai serve notify set is seen without a restart; nil is no channel,
	// with which a reminder is logged instead.
	Notify func() config.Notify
}

// allowsUnsigned is whether the allowance is in force for the project at
// root: its manifest's dashboard.allow_unsigned over the host's. A root whose
// manifest does not load, a repository offered for import say, has the
// host's.
func (o Options) allowsUnsigned(root string) bool {
	host := o.AllowUnsigned != nil && o.AllowUnsigned()
	m, err := manifest.Load(filepath.Join(root, manifest.File))
	if err != nil {
		return host
	}
	return m.Dashboard.AllowUnsignedOr(host)
}

// guard has c ask before each dial whether the allowance is in force for e's
// project and, when flai serve judges images, whether the dashboard's image
// may be dialled.
func guard(o Options, images *imageVerdicts, e Entry, c *channel.Client) {
	c.AllowUnsigned = func() bool { return o.allowsUnsigned(e.Root) }
	if o.Images != nil && images != nil {
		c.BeforeDial = images.gate(o, e)
	}
}

// ImageLists are what flai serve judges the dashboard's image against: the
// verified digest lists cached beside its state and how a missing one is
// fetched, or Err, why it has none it can verify, which leaves every image
// not measured, with Err as the reason, and the release stamp check alone.
type ImageLists struct {
	Cache dashboard.Cache
	Fetch ListFetcher
	Err   error
}

// DefaultDrainGrace is how long a dropped project's requests in flight may
// run on. A publish commits, tags, and pushes in well under a minute, and
// it is the request that raises flai.minimum and so drops the project; five
// minutes covers a slow push without keeping a project that no longer loads
// half-served for long.
const DefaultDrainGrace = 5 * time.Minute

type running struct {
	entry  Entry
	client *channel.Client
	// stop ends everything the project runs; stopLooks only its looks and
	// watcher, leaving its connection open for a drain.
	stop      context.CancelFunc
	stopLooks context.CancelFunc
	done      chan struct{}
}

// halt stops the project's client and waits for it.
func (r *running) halt() {
	r.stop()
	<-r.done
}

// drain stops the project's looks at once, lets the requests in flight be
// answered for at most grace, then stops its client; it returns how many
// were left unanswered.
func (r *running) drain(ctx context.Context, grace time.Duration, why string) int {
	r.stopLooks()
	dctx, cancel := context.WithTimeout(ctx, grace)
	defer cancel()
	left := r.client.Drain(dctx, why)
	r.halt()
	return left
}

// gitOf is the runner that asks git, the system's when the options name none.
func gitOf(o Options) execx.Runner {
	if o.Git == nil {
		return execx.System{}
	}
	return o.Git
}

// dockerOf is the runner that asks docker, the system's when the options name none.
func dockerOf(o Options) execx.Runner {
	if o.Docker == nil {
		return execx.System{}
	}
	return o.Docker
}

// check judges the image the dashboard at url runs, or, when l has no list it
// can verify, says it was not measured and why.
func (l *ImageLists) check(ctx context.Context, docker execx.Runner, url string) ImageVerdict {
	if l.Err != nil {
		return notMeasured(ImageVerdict{}, "flai serve has no verified digest list to judge it by: "+l.Err.Error())
	}
	return CheckImage(ctx, ImageCheck{Runner: docker, URL: url, Cache: l.Cache, Fetch: l.Fetch})
}

// imageVerdicts are each project's last image verdict, by root, recorded by
// the gate before each dial and reported in the status.
type imageVerdicts struct {
	mu sync.Mutex
	by map[string]ImageVerdict
}

// gate is what a project's client asks before each dial, reconnections
// included (S-0238): it judges the image the dashboard container runs and
// records the verdict. An unlisted image is refused, and the channel logs and
// shows the refusal, unless the allowance is in force for the project, when
// it is dialled and recorded as unlisted and Allowed (S-0239); a listed one is
// dialled, and so is one not measured, whose dashboard the release stamp
// check alone then admits or refuses. Every retry of a refused dial looks
// again, so a container restarted from a listed image, or a setting turned
// on, is dialled at the next.
func (iv *imageVerdicts) gate(o Options, e Entry) func(context.Context) error {
	docker := dockerOf(o)
	return func(ctx context.Context) error {
		v := o.Images.check(ctx, docker, e.URL)
		v.Allowed = v.Kind == ImageUnlisted && o.allowsUnsigned(e.Root)
		if iv.record(e.Root, v) {
			logVerdict(o.Logger, e, v)
		}
		if v.Kind == ImageUnlisted && !v.Allowed {
			return errors.New(unsignedImage(v))
		}
		return nil
	}
}

// record keeps v as root's verdict and reports whether it differs from the
// last.
func (iv *imageVerdicts) record(root string, v ImageVerdict) bool {
	iv.mu.Lock()
	defer iv.mu.Unlock()
	was, ok := iv.by[root]
	iv.by[root] = v
	return !ok || !sameVerdict(was, v)
}

// of are the verdicts of roots that have one; nil when none has.
func (iv *imageVerdicts) of(roots []string) map[string]ImageVerdict {
	iv.mu.Lock()
	defer iv.mu.Unlock()
	var out map[string]ImageVerdict
	for _, root := range roots {
		if v, ok := iv.by[root]; ok {
			if out == nil {
				out = map[string]ImageVerdict{}
			}
			out[root] = v
		}
	}
	return out
}

func sameVerdict(a, b ImageVerdict) bool {
	return a.Kind == b.Kind && a.Container == b.Container && a.Image == b.Image && a.Version == b.Version &&
		slices.Equal(a.RepoDigests, b.RepoDigests) && a.Digest == b.Digest && a.Release == b.Release && a.Reason == b.Reason &&
		a.Allowed == b.Allowed
}

// logVerdict says a project's new verdict once. An unlisted one is said here
// only when the allowance dials it; otherwise the channel logs the refusal it
// causes, with this reason.
func logVerdict(log *slog.Logger, e Entry, v ImageVerdict) {
	switch v.Kind {
	case ImageUnlisted:
		if v.Allowed {
			log.Info("dashboard image is on no verified digest list; dialled as unsigned is allowed", "component", "serve", "root", e.Root, "key", e.Key,
				"container", v.Container, "image", v.Image, "reason", v.Reason)
		}
	case ImageListed:
		log.Info("dashboard image is on a verified digest list", "component", "serve", "root", e.Root, "key", e.Key,
			"container", v.Container, "image", v.Image, "digest", v.Digest, "release", v.Release)
	case ImageNotMeasured:
		log.Info("dashboard image not measured; the release stamp check alone admits it", "component", "serve", "root", e.Root, "key", e.Key,
			"reason", v.Reason)
	}
}

// unsignedImage is why the gate refuses v, an unlisted image, in one line.
func unsignedImage(v ImageVerdict) string {
	digests := "no RepoDigests"
	if len(v.RepoDigests) > 0 {
		digests = strings.Join(v.RepoDigests, ", ")
	}
	return fmt.Sprintf("the dashboard runs an unsigned image %s (%s): %s", v.Image, digests, v.Reason)
}

// Run serves every registered project until ctx ends. The registry is read
// again every tick, so a project registered while it runs is picked up and
// one that was removed or changed is dropped.
func Run(ctx context.Context, o Options) error {
	if o.Every <= 0 {
		o.Every = time.Second
	}
	if o.Now == nil {
		o.Now = time.Now
	}
	if o.Logger == nil {
		o.Logger = slog.New(slog.DiscardHandler)
	}
	// A change of the host's budget state is sent through the reminder
	// channel, read at each change (S-0534, ADR-0205).
	o.Dir.notifyBudget(o.Notify)
	if o.DrainGrace <= 0 {
		o.DrainGrace = DefaultDrainGrace
	}
	if o.NewClient == nil {
		o.NewClient = func(e Entry, key []byte) *channel.Client { return defaultClient(o, e, key) }
	}
	if st, alive := o.Dir.ReadStatus(o.Now()); alive && st.PID != os.Getpid() {
		return fmt.Errorf("flai serve is already running (pid %d, since %s)", st.PID, st.Started)
	}
	started := o.Now().UTC().Format(time.RFC3339)
	clients := map[string]*running{}
	// A project dropped because it stopped loading answers its requests in
	// flight before its client stops (I-0107): it waits here, by root, and is
	// not served again until its drain is over.
	draining := map[string]*running{}
	var drains sync.WaitGroup
	unavailable := map[string]string{}
	images := &imageVerdicts{by: map[string]ImageVerdict{}}
	offered := &offers{o: o, running: map[string]*running{}, images: images}
	var mu sync.Mutex
	// Each Claude Code that flai serve may start an agent with is checked
	// against permission_prompt once (S-0286, ADR-0106): at the first look
	// after it starts, and at a later look once the claude it runs has
	// changed. Every look that may start an agent reads the operator's
	// say first, so the check begins there, beside the start, never before it.
	checker := newClaudeChecker(o, func() []Entry {
		mu.Lock()
		defer mu.Unlock()
		out := make([]Entry, 0, len(clients))
		for _, r := range clients {
			out = append(out, r.entry)
		}
		return out
	})
	defer checker.wait()
	// Each served project's history is looked at for a newer flai release
	// every minute, and at once when it is first served (ADR-0158).
	upgrades := newUpgrader(o)
	defer upgrades.wait()
	// Each served project's stories in review are pushed, and their pull
	// requests opened, once a minute under pull-request acceptance while its
	// forge host action is on (ADR-0166); a push in flight is waited for.
	var forgeLooks sync.WaitGroup
	defer forgeLooks.Wait()
	// Each served project's merge queue builds and verifies one batch at a
	// time, off the agents' turns, behind the merge-queue host action
	// (S-0399, ADR-0173); a batch running is cancelled and waited for.
	queues := newQueueRuns()
	defer queues.wait()
	// Each served project's threads that await the operator past
	// threads.remind_after are reminded through the host's channel, off the
	// project's loop (ADR-0205); a send in flight is waited for.
	var reminds sync.WaitGroup
	defer reminds.Wait()
	if agent := o.Agent; agent != nil {
		o.Agent = func(root string) AgentConfig {
			cfg := agent(root)
			checker.consider(ctx, root, cfg)
			return cfg
		}
	}
	defer func() {
		offered.halt()
		for _, r := range clients {
			r.halt()
		}
		drains.Wait() // each drain's bound is within ctx, which has ended
		_ = os.Remove(o.Dir.status())
	}()

	reconcile := func() {
		entries, err := o.Dir.Projects()
		if err != nil {
			o.Logger.Warn("registry unreadable", "component", "serve", "err", err.Error())
			return
		}
		entries = append(entries, offered.folderProjects(entries)...)
		mu.Lock()
		defer mu.Unlock()
		// A project that cannot be served is said once, when it goes and when
		// it comes back, not tried and logged at every tick (S-0121).
		want, gone := map[string]Entry{}, map[string]string{}
		for _, e := range entries {
			if why := e.Unavailable(); why != "" {
				gone[e.Root] = why
				if unavailable[e.Root] != why {
					o.Logger.Warn("registered project not served", "component", "serve", "root", e.Root, "key", e.Key, "reason", why)
				}
				continue
			}
			want[e.Root] = e
			if _, was := unavailable[e.Root]; was {
				o.Logger.Info("registered project can be served again", "component", "serve", "root", e.Root, "key", e.Key)
			}
		}
		unavailable = gone
		for root, r := range clients {
			e, ok := want[root]
			if ok && e == r.entry {
				continue
			}
			delete(clients, root)
			// A project that stopped loading, its manifest raised past this
			// flai by the very publish it is answering, say, takes no new
			// request and answers those in flight before it stops (I-0107).
			if why, lost := gone[root]; lost && !ok {
				draining[root] = r
				drains.Add(1)
				go func() {
					defer drains.Done()
					left := r.drain(ctx, o.DrainGrace, why)
					mu.Lock()
					delete(draining, root)
					mu.Unlock()
					level := slog.LevelInfo
					if left > 0 {
						level = slog.LevelWarn
					}
					o.Logger.Log(ctx, level, "project dropped", "component", "serve", "root", root, "reason", why, "unanswered", left)
				}()
				continue
			}
			// A project that left the registry, or is served under another
			// key, is said to be removed first, so that the dashboard drops
			// it rather than show it as waiting for flai (S-0121).
			why := "its registry entry changed"
			if !ok || e.Key != r.entry.Key {
				r.client.Notify(Removed, map[string]string{"project": r.entry.Key})
				why = "it left the registry"
				if ok {
					why = "it is served under another key"
				}
			}
			r.halt()
			o.Logger.Info("project dropped", "component", "serve", "root", root, "reason", why)
		}
		for root, r := range draining {
			// one removed while it drains stops at once, as any removed project does
			_, back := want[root]
			if _, still := gone[root]; !back && !still {
				r.client.Notify(Removed, map[string]string{"project": r.entry.Key})
				r.stop()
			}
		}
		for root, e := range want {
			if _, ok := clients[root]; ok {
				continue
			}
			if _, busy := draining[root]; busy {
				continue // served again once its drain is over
			}
			key, err := os.ReadFile(e.KeyFile)
			if err != nil || strings.TrimSpace(string(key)) == "" {
				o.Logger.Warn("agent credential unreadable", "component", "serve", "root", root, "file", e.KeyFile)
				continue
			}
			clientCtx, stop := context.WithCancel(ctx)
			cctx, stopLooks := context.WithCancel(clientCtx)
			r := &running{entry: e, client: o.NewClient(e, []byte(strings.TrimSpace(string(key)))), stop: stop, stopLooks: stopLooks, done: make(chan struct{})}
			// the allowance is read and the dashboard's image judged before each dial (S-0238, S-0239)
			guard(o, images, e, r.client)
			clients[root] = r
			watcher := &watch.Watcher{Root: e.Root, Paths: watchedPaths(e.Root), Every: o.WatchEvery}
			starter := newLauncher(o, e)
			// An agent starting or ending changes no file: the dashboard is told
			// so that a story page shows it at once (S-0154), a planner's
			// item page, with the item in place of the story (S-0208), and
			// the orchestrator's run, with its role (S-0218), the analyzer's,
			// with its (S-0223), the curator's, with its (ADR-0196), and a
			// defined role's, with its verb and name (S-0417), as agentAbout
			// says.
			starter.changed = func(run *AgentRun) { r.client.Notify(AgentChanged, agentAbout(e.Key, run)) }
			// The replanner looks after each look of the launcher: on the same
			// changes, when an agent ends and its queue can move on, and every
			// minute, which is when the schedule is seen to come round (S-0211).
			// So does the orchestrator, which starts its run again when it
			// ends and stops it once the action is off (S-0218), and the
			// analysis scheduler, which sees analysis.schedule come round as
			// the replanner sees planning.schedule (S-0223), and the curate
			// scheduler, which sees stories accepted, the curator's run end,
			// and digest.schedule come round (S-0432, ADR-0199).
			replanner := newReplanner(o, e, starter)
			orch := newOrchestrator(o, e, starter)
			analyst := newAnalysisScheduler(o, e, starter)
			curation := newCurateScheduler(o, e, starter)
			// The thread reminder looks on the same changes and every minute,
			// in the background, as a send may wait seconds (ADR-0205).
			reminding := newReminder(o, e, &reminds)
			// The confirmer looks on the same changes and every minute, in
			// the loop, as it only writes the project's files, and before the
			// reminder, which then reminds of no thread it confirmed; a
			// confirmation is a change under threads, on which the launcher
			// starts the asking agent again (ADR-0206).
			confirming := newConfirmer(o, e)
			starter.look(cctx, false) // starts what is ready and has had no agent since (S-0112)
			replanner.look(cctx)
			orch.look(cctx)
			analyst.look(cctx)
			curation.look(cctx)
			confirming.look(cctx)
			reminding.soon(cctx)
			newForgePass(o, e, &forgeLooks).start(cctx)
			// flai accept's requests to merge a pull request are taken every
			// second, and answered behind the forge host action (S-0397)
			newMergePass(o, e, &forgeLooks).start(cctx)
			queue := newQueuePass(o, e, queues)
			queue.start(cctx)
			go func() {
				for {
					select {
					case <-cctx.Done():
						return
					case <-starter.again:
						starter.look(cctx, true)
					case <-time.After(time.Minute):
						// an agent that outlived an earlier flai serve is not waited for
						starter.look(cctx, false)
					}
					replanner.look(cctx)
					orch.look(cctx)
					analyst.look(cctx)
					curation.look(cctx)
					confirming.look(cctx)
					reminding.soon(cctx)
				}
			}()
			go func() {
				// The dashboard hears of each changed file; what it shows comes from asking again.
				go watcher.Run(cctx, func(rel string) {
					r.client.Notify("change", map[string]string{"project": e.Key, "path": rel})
					// a story's state is in its file under kanban, and the pull order in the
					// board; an answer to an agent's question is in its thread
					if slashed := filepath.ToSlash(rel); strings.Contains(slashed, "/kanban/") || strings.Contains(slashed, "/threads/") {
						starter.look(cctx, false)
						replanner.look(cctx)
						orch.look(cctx)
						analyst.look(cctx)
						curation.look(cctx)
						confirming.look(cctx)
						reminding.soon(cctx)
						queue.wake()
					}
				})
				r.client.Run(clientCtx)
				close(r.done)
			}()
			o.Logger.Info("project served", "component", "serve", "root", root, "key", e.Key, "dashboard", e.URL)
			if out := release.FlaiOutdated(gitOf(o), root, o.Version); out != nil {
				o.Logger.Warn("flai is older than the project it serves", "component", "serve", "root", root,
					"running", out.Running, "newest", out.Newest, "detail", out.Message)
			}
		}
		offered.reconcile(ctx, entries)
		served := make([]Entry, 0, len(clients))
		for _, r := range clients {
			served = append(served, r.entry)
		}
		upgrades.soon(ctx, served)
	}
	writeStatus := func() {
		mu.Lock()
		st := Status{PID: os.Getpid(), Version: o.Version, Started: started, Updated: o.Now().UTC().Format(time.RFC3339), Connections: map[string]channel.State{}}
		roots := make([]string, 0, len(clients))
		for root, r := range clients {
			st.Connections[root] = r.client.State()
			roots = append(roots, root)
		}
		st.Images = images.of(roots)
		st.Offered = offered.found
		st.Folder = o.Folder
		st.FolderProjects, st.ImportProjects, st.Unserved = offered.placement()
		if len(unavailable) > 0 {
			st.Unavailable = maps.Clone(unavailable)
		}
		mu.Unlock()
		if err := o.Dir.write(o.Dir.status(), st); err != nil {
			o.Logger.Warn("status not written", "component", "serve", "err", err.Error())
		}
	}

	mcp := &mcpTeller{o: o}
	tellMCP := func() {
		mu.Lock()
		roots := make([]string, 0, len(clients))
		for root := range clients {
			roots = append(roots, root)
		}
		mu.Unlock()
		mcp.tell(roots)
	}

	tick := time.NewTicker(o.Every)
	defer tick.Stop()
	for {
		reconcile()
		writeStatus()
		tellMCP()
		select {
		case <-ctx.Done():
			return nil
		case <-tick.C:
		}
	}
}

// agentAbout is what the agent notification of run, in the project named
// key, says it is about: the project, then a defined role's verb as role,
// its name as name, and its item when it ran on one (S-0417); a planner's
// item (S-0208); the orchestrator's, the analyzer's, or the curator's role
// (S-0218, S-0223, ADR-0196); or a story's agent's story (S-0154).
func agentAbout(key string, run *AgentRun) map[string]string {
	switch {
	case run.defined():
		about := map[string]string{"project": key, "role": run.Role, "name": run.RoleName}
		if run.Item != "" {
			about["item"] = run.Item
		}
		return about
	case run.Item != "":
		return map[string]string{"project": key, "item": run.Item}
	case run.orchestrates():
		return map[string]string{"project": key, "role": conventions.RoleOrchestrate}
	case run.analyzes():
		return map[string]string{"project": key, "role": conventions.RoleAnalyze}
	case run.curates():
		return map[string]string{"project": key, "role": conventions.RoleCurate}
	}
	return map[string]string{"project": key, "story": run.Story}
}

// defaultClient is the client for one project when the options name no
// other: it sends this build's release stamp in hello, none for a
// development build, and verifies the dashboard's with the embedded release
// key (S-0237). Run gives it the allowance and the image gate with guard.
func defaultClient(o Options, e Entry, key []byte) *channel.Client {
	keys, err := buildinfo.ReleaseKeys()
	if err != nil {
		keys = nil // the client reads them again and logs why it cannot
	}
	return &channel.Client{URL: e.URL, Key: key, Project: channel.Project{Key: e.Key, Name: e.Name, Root: e.Root},
		Methods: hostapi.MethodsFor(o.Version, o.Now, o.Host), Version: o.Version, Logger: o.Logger,
		Release:     channel.Release{Statement: buildinfo.ReleaseStatement, Signature: buildinfo.ReleaseSignature},
		ReleaseKeys: keys}
}

// watchedPaths are the manifest's three folders and the manifest itself:
// what the dashboard shows, and nothing else of the clone.
func watchedPaths(root string) []string {
	out := []string{manifest.File}
	m, err := manifest.Load(filepath.Join(root, manifest.File))
	if err != nil {
		return append(out, "design", "docs", "wip")
	}
	for _, k := range []string{"design", "docs", "wip"} {
		if dir := m.Layout[k]; dir != "" && !filepath.IsAbs(dir) && !strings.HasPrefix(filepath.Clean(dir), "..") {
			out = append(out, dir)
		}
	}
	return out
}
