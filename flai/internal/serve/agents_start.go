package serve

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"maps"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/config"
	"github.com/bytepunx/system-flow/flai/internal/harness"
	"github.com/bytepunx/system-flow/flai/internal/host"
	"github.com/bytepunx/system-flow/flai/internal/hostapi"
	"github.com/bytepunx/system-flow/flai/internal/inbox"
	"github.com/bytepunx/system-flow/flai/internal/manifest"
	"github.com/bytepunx/system-flow/flai/internal/messages"
	"github.com/bytepunx/system-flow/flai/internal/roles"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// agentEnv is flai serve's environment without the host's address and token
// and the config it was started with (S-0183, I-0044): a flai the agent runs
// finds its own config, and a flai serve it runs does not reach the
// operator's host. Nor does it carry what marks the session of whoever
// started it, such as an agent that ran flai plan (S-0208): each run sets
// FLAI_STORY, or FLAI_ROLE with FLAI_ITEM, FLAI_FOCUS (S-0223), or
// FLAI_SINCE (ADR-0196), for itself, and has the others unset. Nor the HTTP
// MCP server's URL and token flai serve may have been handed itself: a run
// gets them only when its adapter asks (mcpEnv, S-0506).
func agentEnv(env []string) []string {
	return slices.DeleteFunc(slices.Clone(env), func(kv string) bool {
		name, _, _ := strings.Cut(kv, "=")
		return name == host.URLEnv || name == host.TokenEnv || name == config.EnvVar || name == MCPURLEnv || name == MCPTokenEnv || slices.Contains(runMarks, name)
	})
}

// The project's HTTP MCP server as flai serve hands it to a run whose
// adapter asks for it (S-0506): MCPURLEnv and MCPTokenEnv are the run's
// variables that carry its URL and bearer token, MCPStateFile and
// MCPTokenFile the files under the project's .flai-cache where flai mcp http
// says where it listens and keeps its token, and MCPStaleAfter how long its
// state, rewritten every few seconds while it runs, stays true.
const (
	MCPURLEnv     = "FLAI_MCP_URL"
	MCPTokenEnv   = "FLAI_MCP_TOKEN"
	MCPStateFile  = "mcp-http.json"
	MCPTokenFile  = "mcp.token"
	MCPStaleAfter = 10 * time.Second
)

// MCPAsker is an adapter whose agent reaches flai's MCP server over HTTP
// rather than starting flai mcp on stdio, such as one that runs it in a
// container: flai serve hands each of its runs the project's server's URL
// and token as MCPURLEnv and MCPTokenEnv, in the run's environment alone,
// never its argv, a start request, or a log line (S-0506). An adapter that
// does not implement it, claude-code on the host among them, gets neither.
type MCPAsker interface {
	// MCPOverHTTP says whether the adapter's agent reaches the server over HTTP.
	MCPOverHTTP() bool
}

// mcpEnv is what a run of adapter, the harness name, is handed of the HTTP
// MCP server of the project at root: its URL and token, as MCPURLEnv and
// MCPTokenEnv, when the adapter asks for them (MCPAsker), and nothing
// otherwise. It refuses a start that asks while the server is not running
// at now or its token cannot be read, saying what to run and never the
// token; now is read only for an adapter that asks.
func mcpEnv(adapter harness.Adapter, name, root string, now func() time.Time) ([]string, error) {
	if asker, ok := adapter.(MCPAsker); !ok || !asker.MCPOverHTTP() {
		return nil, nil
	}
	repo, err := workitem.Open(root)
	if err != nil {
		return nil, fmt.Errorf("open the project at %s for its HTTP MCP server, which harness %s reaches flai through: %w", root, name, err)
	}
	url, ok := mcpServerURL(repo.CacheDir(), now())
	if !ok {
		return nil, fmt.Errorf("harness %s reaches flai's MCP server over HTTP, and the project's is not running: flai host starts it for every project flai serve serves, or run flai mcp start in the project, and start the agent again", name)
	}
	data, err := os.ReadFile(filepath.Join(repo.CacheDir(), MCPTokenFile))
	token := strings.TrimSpace(string(data))
	if err != nil || token == "" {
		return nil, fmt.Errorf("harness %s reaches flai's MCP server over HTTP, and the server's token in %s cannot be read: run flai mcp token in the project, and start the agent again", name, filepath.Join(repo.CacheDir(), MCPTokenFile))
	}
	return []string{MCPURLEnv + "=" + url, MCPTokenEnv + "=" + token}, nil
}

// mcpServerURL is the URL the HTTP MCP server whose state is in cache says it
// listens at, the address flai host gave it, and whether it is running: its
// state was rewritten within MCPStaleAfter of now and its process is alive.
func mcpServerURL(cache string, now time.Time) (string, bool) {
	var st struct {
		PID     int    `json:"pid"`
		URL     string `json:"url"`
		Updated string `json:"updated"`
	}
	data, err := os.ReadFile(filepath.Join(cache, MCPStateFile))
	if err != nil || json.Unmarshal(data, &st) != nil || st.URL == "" {
		return "", false
	}
	updated, err := time.Parse(time.RFC3339, st.Updated)
	if err != nil || now.Sub(updated) > MCPStaleAfter {
		return "", false
	}
	return st.URL, Alive(st.PID)
}

// withoutForgeTokens is env without the variables each forge's token is
// read from: the forge's default, such as GITHUB_TOKEN, and any the host
// names with flai serve forge (ADR-0165). No agent flai serve starts gets
// one, whatever the project's acceptance mode: a story reaches its forge
// only through flai serve's forge host action, on the host (S-0393).
func (l *launcher) withoutForgeTokens(env []string) []string {
	drop := map[string]bool{}
	for _, f := range config.KnownForges() {
		drop[config.DefaultForgeTokenEnv(f)] = true
	}
	if l.config != nil {
		for _, name := range l.config(l.entry.Root).ForgeTokenEnvs {
			drop[name] = true
		}
	}
	return slices.DeleteFunc(env, func(kv string) bool {
		name, _, _ := strings.Cut(kv, "=")
		return drop[name]
	})
}

// runMarks are the variables that say what a run is for.
var runMarks = []string{"FLAI_STORY", "FLAI_ROLE", "FLAI_ITEM", "FLAI_FOCUS", "FLAI_SINCE"}

// start starts an agent for story, or starts again the one that ended asking
// in story.Asked, for story.Answer, or the one of story.Resume, for
// story.SentBack, and says whether it did. What its
// harness cannot do refuses it as startAllowed says (ADR-0131), as any
// failed start is refused. Every start of a story's agent comes here, and
// passes the host's budget state's gate once its harness and provider are
// known (S-0534): a start the gate holds records nothing, no run and no
// failure, and is kept in holding for its caller to say; one it lets through
// as the probe has its process recorded on the probe once spawned, and frees
// the probe when it cannot be made.
func (l *launcher) start(ctx context.Context, cfg AgentConfig, story readyStory) bool {
	now := l.now().UTC()
	name := cfg.Name
	if name == "" {
		name = "agent"
	}
	run := &AgentRun{Story: story.ID, Agent: name + "-" + story.ID, Started: now.Format(time.RFC3339), Session: newSession(), AutoRestarts: story.AutoRestarts}
	if story.Asked != nil {
		run.Agent, run.Answered = story.Asked.Agent, story.Answer
		if story.Asked.Session != "" {
			run.Session = story.Asked.Session
		}
	}
	if story.SentBack != nil {
		run.SentBack = story.SentBack.Thread
		if story.Resume != nil {
			run.Agent, run.Session = story.Resume.Agent, story.Resume.Session
		}
	}
	run.StoryAgent = &manifest.Agent{}
	if story.Agent != nil {
		run.Model = story.Agent.Model
		run.StoryAgent = story.Agent
	}
	run.Retry = story.AfterLimit
	entry := hostapi.Entry{At: run.Started, Action: hostapi.ActionAgent, Method: "serve.agent", Project: l.entry.Key, Root: l.entry.Root, By: "flai serve"}
	probe := false // the budget state's gate let this start through as its probe
	fail := func(err error) bool {
		if probe {
			l.budgetRelease(run)
		}
		run.Error, run.Ended, run.Outcome, run.Why = err.Error(), run.Started, OutcomeFailed, "could not be started: "+err.Error()
		l.told(run)
		what := run.Command
		if what == "" {
			what = "an agent"
		}
		entry.Outcome, entry.Detail = "failed", fmt.Sprintf("%s for %s could not be started: %s", what, story.ID, run.Error)
		if l.record != nil {
			l.record(entry)
		}
		l.log("agent could not be started", "story", story.ID, "harness", run.Harness, "err", run.Error)
		appendAgentEvent(l.entry.Root, run, inbox.AgentFailed, run.Why, l.warn)
		return false
	}
	if !StoryID.MatchString(story.ID) {
		return fail(fmt.Errorf("%q is not a story's ID", story.ID))
	}
	name, adapter, err := harness.For(story.Agent, cfg.host(harness.Command).Program != "")
	run.Harness = name
	if err != nil {
		return fail(err)
	}
	// The budget state holds a start on its provider before anything is
	// recorded (S-0534): a held start is no failed run, and its story waits
	// where it is. provide records the same provider once the gate is passed.
	if story.Agent != nil {
		run.Provider = story.Agent.Provider
	}
	gate := l.budgetStart(run)
	if !gate.Go {
		l.holdStart(story.ID, gate.State)
		return false
	}
	probe = gate.Probe
	unguarded, err := startAllowed(name, adapter.Capabilities(), story.Agent, cfg)
	if err != nil {
		return fail(err)
	}
	var shares []harness.Share
	if story.Asked == nil && story.Resume == nil && story.Commit == "" {
		shares = sharedWith(l.entry.Root, story.ID) // a resumed session knows them; a commit needs none
	}
	auto := autoRestart(story.AutoRestarts, cfg.AutoRestarts)
	if story.AfterLimit {
		auto = "" // the count it carries is its last run's, and this start is none of them
	}
	spec, keys, err := l.startOf(cfg, run, adapter, name, harness.Request{Story: story.ID, Root: l.entry.Root, Project: l.entry.Key, Agent: story.Agent, Name: run.Agent, Flai: cfg.Flai,
		Session: run.Session, Answered: run.Answered, SentBack: story.SentBack, Restart: story.Restart, AutoRestart: auto, AfterLimit: story.AfterLimit, Commit: story.Commit, Started: story.Started, Past: story.Past, Begun: story.Begun,
		Shares: shares, Formula: formulaOf(l.entry.Root, story)})
	if err != nil {
		return fail(err)
	}
	spend := l.spendAtStart(ctx, cfg, run) // before its first call is charged
	cmd, out, err := l.spawnWith(ctx, run, spec.Argv, story.ID, now, slices.Concat([]string{"FLAI_STORY=" + story.ID}, spec.Env, keys), spec.Stdin)
	if err != nil {
		return fail(err)
	}
	if probe {
		l.budgetSpawned(run)
	}
	l.keepSpend(run, spend)
	l.told(run)
	entry.Outcome, entry.Detail = "done", fmt.Sprintf("started %s (%s) for %s as %s (pid %d); log %s", run.Command, run.Harness, story.ID, run.Agent, run.PID, run.Log)
	if run.Answered != "" {
		what := run.Answered + " answered"
		if strings.HasPrefix(run.Answered, "MS-") {
			what = "a message on " + run.Answered
		}
		entry.Detail = fmt.Sprintf("started %s (%s) again for %s as %s, %s (pid %d); log %s", run.Command, run.Harness, story.ID, run.Agent, what, run.PID, run.Log)
	}
	if run.SentBack != "" {
		how := "in a new session"
		if story.Resume != nil {
			how = "in its session"
		}
		entry.Detail = fmt.Sprintf("started %s (%s) again for %s as %s %s, flai having sent it back from review on %s (pid %d); log %s", run.Command, run.Harness, story.ID, run.Agent, how, run.SentBack, run.PID, run.Log)
	}
	if story.Commit != "" {
		entry.Detail = fmt.Sprintf("started %s (%s) for %s as %s to commit what its worktree holds (pid %d); log %s", run.Command, run.Harness, story.ID, run.Agent, run.PID, run.Log)
	}
	if story.Started && len(story.Past) > 0 {
		entry.Detail = fmt.Sprintf("started %s (%s) for %s as %s on the operator's word, past %s (pid %d); log %s", run.Command, run.Harness, story.ID, run.Agent, strings.Join(story.Past, "; "), run.PID, run.Log)
	}
	if story.Begun != nil {
		entry.Detail = fmt.Sprintf("started %s (%s) for %s as %s on the operator's word, begun %s (pid %d); log %s", run.Command, run.Harness, story.ID, run.Agent, story.Begun.Said(), run.PID, run.Log)
	}
	switch {
	case story.AfterLimit:
		entry.Detail = fmt.Sprintf("started %s (%s) again for %s as %s on its own, its last agent having ended on the usage limit, which no longer holds %s; no automatic restart spent (pid %d); log %s", run.Command, run.Harness, story.ID, run.Agent, budgetKey(run.Provider, run.Harness), run.PID, run.Log)
	case story.AutoRestarts > 0:
		entry.Detail = fmt.Sprintf("started %s (%s) again for %s as %s on its own, automatic restart %d of %d, its last agent having ended in progress (pid %d); log %s", run.Command, run.Harness, story.ID, run.Agent, story.AutoRestarts, cfg.AutoRestarts, run.PID, run.Log)
	}
	if probe {
		entry.Detail += "; the probe of the budget state of " + budgetKey(run.Provider, run.Harness)
	}
	said := []any{"story", story.ID, "harness", run.Harness, "command", run.Command, "pid", run.PID}
	if unguarded != "" {
		entry.Detail += "; " + unguarded
		said = append(said, "roles", unguarded)
	}
	if l.record != nil {
		l.record(entry)
	}
	l.log("agent started", said...)
	appendAgentEvent(l.entry.Root, run, inbox.AgentStarted, startWhy(story, run, cfg.AutoRestarts), l.warn)
	l.await(cmd, out, run.PID, func(code int) {
		ended := *run
		ended.Ended, ended.Exit = l.now().UTC().Format(time.RFC3339), &code
		l.endRun(&ended, &code)
		l.log("agent ended", "story", story.ID, "pid", run.PID, "exit", code, "outcome", ended.Outcome)
		l.measureWith(story.ID, cfg.Providers)
	})
	return true
}

// holdStart keeps story, whose agent's start the budget state b held, in
// holding (S-0534).
func (l *launcher) holdStart(story string, b Budget) {
	if l.holding == nil {
		l.holding = map[string]Budget{}
	}
	l.holding[story] = b
}

// sayHeld logs, in one line, the stories whose agent's start the budget
// state held at this look outside the ready column, resumed, sent back, or
// started again, and why, when they are not those the last look logged; the
// ready column's are said with its other waits (say).
func (l *launcher) sayHeld() {
	stories := slices.Sorted(maps.Keys(l.holding))
	var whys []string
	for _, id := range stories {
		if why := budgetHeldWhy(l.holding[id], l.now()); !slices.Contains(whys, why) {
			whys = append(whys, why)
		}
	}
	said := strings.Join(stories, ", ")
	if said != "" && said != l.heldSaid {
		l.log("agents not started: the budget state holds them", "stories", said, "why", strings.Join(whys, "; "))
	}
	l.heldSaid = said
}

// startWhy is why story's agent, run, was started, for its started agent
// event (S-0532): the start's reason in a short form of the words its
// journal entry, or the agent's prompt, gives it, the one that wins in the
// journal winning here; limit is how many automatic restarts in a row the
// project allows.
func startWhy(story readyStory, run *AgentRun, limit int) string {
	last := func(why string) string {
		if story.Restart == "" {
			return why
		}
		return why + ": its last agent " + story.Restart
	}
	switch {
	case story.AfterLimit:
		return last("flai serve started it again on its own once the usage limit no longer held it, spending no automatic restart")
	case story.AutoRestarts > 0:
		return last(fmt.Sprintf("flai serve started it again on its own, automatic restart %d of %d", story.AutoRestarts, limit))
	case story.Begun != nil:
		return "the operator started it here, begun " + story.Begun.Said()
	case story.Started && len(story.Past) > 0:
		return "the operator started it, past " + strings.Join(story.Past, "; ")
	case story.Commit != "":
		return "the operator started it to commit what its worktree holds"
	case run.SentBack != "":
		how := "in a new session"
		if story.Resume != nil {
			how = "in its session"
		}
		return fmt.Sprintf("started again %s, flai having sent it back from review on %s", how, run.SentBack)
	case strings.HasPrefix(run.Answered, "MS-"):
		return "started again in its session, for a message on " + run.Answered
	case run.Answered != "":
		return "started again in its session, " + run.Answered + " answered"
	case story.Restart != "":
		return last("the operator restarted it")
	case story.Started:
		return "the operator started it"
	}
	return "it entered ready"
}

// startOf is what adapter starts for req, the run of the harness name, with
// the provider the agent names resolved into req and recorded on run, and the
// variables that carry the provider's key into the child (ADR-0129), with the
// project's HTTP MCP server's URL and token when the adapter asks for them
// (mcpEnv, S-0506), the URL's host the container's own name for the host
// when the container harness's network is not the host's (S-0508), and the
// container's image the project's layer when its manifest names one
// (containerImageHost, S-0510), and the host's deny_protected for the
// harness given to the adapter, which writes its deny rules (ADR-0260): every
// start, a story's agent's and a strategic agent's, builds its command here.
// On the container harness it readies the host for the start too, a story's
// and a strategic agent's alike (S-0508, S-0509): containerAllowed before
// the adapter is asked, containerReachable once the MCP server's URL is
// known, refusing a loopback bind a native engine's bridge cannot reach
// (S-0619, I-0254), and containerReady once the adapter has taken the start.
// The keys and the token go to spawnWith's env and nowhere else, not run, a log
// line, the journal, or an error, so that a secret flai serve holds reaches
// the child alone. A refusal is the start's
// failure, as an adapter's refusal is.
func (l *launcher) startOf(cfg AgentConfig, run *AgentRun, adapter harness.Adapter, name string, req harness.Request) (harness.Start, []string, error) {
	if name == harness.Container {
		if err := l.containerAllowed(cfg, req.Agent); err != nil {
			return harness.Start{}, nil, err
		}
		// The adapter reads a story's worktree's git to mount it (ADR-0246),
		// so the worktree is opened before it is asked.
		if req.Story != "" {
			if err := l.containerWorktree(req.Story); err != nil {
				return harness.Start{}, nil, err
			}
		}
	}
	if err := l.provide(cfg, run, &req); err != nil {
		return harness.Start{}, nil, err
	}
	if err := l.builtinFile(&req); err != nil {
		return harness.Start{}, nil, err
	}
	host := cfg.host(name)
	host.DenyProtected = cfg.Lifts[name].DenyProtected
	if name == harness.Container {
		h, _, err := containerImageHost(l.entry.Root, cfg, innerOf(req.Agent))
		if err != nil {
			return harness.Start{}, nil, err
		}
		host = h
	}
	spec, err := adapter.Start(req, host)
	if err != nil {
		return harness.Start{}, nil, err
	}
	if c := spec.Container; c != nil {
		// Recorded so that agents.json says what the run ran once --rm has
		// removed its container (I-0274).
		run.Image, run.Container, run.Mounts = c.Image, c.Name, c.Mounts
	}
	keys, err := providerKeys(req, spec, os.LookupEnv)
	if err != nil {
		return harness.Start{}, nil, err
	}
	mcp, err := mcpEnv(adapter, name, l.entry.Root, l.now)
	if err != nil {
		return harness.Start{}, nil, err
	}
	if name == harness.Container {
		if len(mcp) > 0 {
			if err := l.containerReachable(cfg, mcp, host.Network); err != nil {
				return harness.Start{}, nil, err
			}
		}
		mcp = containerMCP(mcp, host.Network)
		if err := l.containerReady(cfg, run, req); err != nil {
			return harness.Start{}, nil, err
		}
	}
	return spec, append(keys, mcp...), nil
}

// builtinFile gives the start of a built-in role, the planner, the
// orchestrator, the analyzer, or the curator, the role file it runs from:
// the project's under its design folder, else the baseline flai carries
// (roles.Builtin), so that its definition's prompt is the file's body with
// the project's additions below the marker (S-0419, ADR-0207 §4). Nothing for
// a story's agent or a role the project defines, which carries its own file.
// A manifest or a role file that cannot be read refuses the start.
func (l *launcher) builtinFile(req *harness.Request) error {
	if req.Story != "" || req.RoleFile != nil || req.Role == "" {
		return nil
	}
	m, err := manifest.Load(filepath.Join(l.entry.Root, manifest.File))
	if err != nil {
		return fmt.Errorf("read the project's %s for the design folder its role files are under: %w; mend it (flai check says what is wrong) and start the agent again", manifest.File, err)
	}
	design := m.Layout["design"]
	if design == "" {
		design = "design"
	}
	b, ok := roles.BaselineFor(design, req.Role)
	if !ok {
		return nil
	}
	f, _, err := roles.Builtin(l.entry.Root, design, b.Name)
	if err != nil {
		return err
	}
	req.Builtin = &f
	return nil
}

// provide resolves the provider req's agent names, from the host's providers
// with the project's manifest's overrides laid over, into req, and records
// its name and its key's variable on run. It refuses a name the host has no
// entry for, an entry that is not whole, and a manifest that cannot be read:
// a start without its overrides could send the agent's model calls to
// another gateway than the project says. Nothing when the agent names none,
// and the manifest is not read.
func (l *launcher) provide(cfg AgentConfig, run *AgentRun, req *harness.Request) error {
	if req.Agent == nil || req.Agent.Provider == "" {
		return nil
	}
	name := req.Agent.Provider
	run.Provider = name
	m, err := manifest.Load(filepath.Join(l.entry.Root, manifest.File))
	if err != nil {
		return fmt.Errorf("read the project's %s for the overrides of provider %s: %w; mend it (flai check says what is wrong) and start the agent again", manifest.File, name, err)
	}
	p, err := manifest.ResolveProvider(name, cfg.Providers, m.Providers)
	if err != nil {
		return err
	}
	req.ProviderName, req.Provider = name, p
	run.KeyEnv = p.KeyEnv
	return nil
}

// providerKeys are the variables that carry the key of req's provider into
// the child: each of spec.EnvFrom, child=value, its value that of the
// variable of flai serve's own environment it names, as lookup reads it, in
// the order of the child's variables. It refuses a start whose provider's
// key_env, or a variable EnvFrom names, is unset or empty, naming the
// variable and the fix, never a value. A command adapter's child reads the
// key from the variable key_env names itself, so the variable is checked
// whatever the harness.
func providerKeys(req harness.Request, spec harness.Start, lookup func(string) (string, bool)) ([]string, error) {
	var vars []string
	if req.Provider != nil && req.Provider.KeyEnv != "" {
		vars = append(vars, req.Provider.KeyEnv)
	}
	children := slices.Sorted(maps.Keys(spec.EnvFrom))
	for _, child := range children {
		if v := spec.EnvFrom[child]; !slices.Contains(vars, v) {
			vars = append(vars, v)
		}
	}
	for _, v := range vars {
		if value, ok := lookup(v); !ok || value == "" {
			return nil, fmt.Errorf("provider %[1]s's key is read from %[2]s, which is unset or empty in flai serve's environment: set %[2]s in the environment flai host or flai serve runs in and start the agent again, or name the variable that holds the key with flai serve agent provider %[1]s --key-env <variable>", req.ProviderName, v)
		}
	}
	keys := make([]string, 0, len(children))
	for _, child := range children {
		value, _ := lookup(spec.EnvFrom[child])
		keys = append(keys, child+"="+value)
	}
	return keys, nil
}

// sharedWith are the shares in force on which story is the held side, each
// with the conversation it was made on, for its agent's prompt (ADR-0134):
// kept in a conversation that reads as open, and in force as
// workitem.Holds.WithShares judges one. None when they cannot be read: the
// prompt then names no share, and the agent's first inbox still lists the
// conversation.
func sharedWith(root, story string) []harness.Share {
	repo, err := workitem.Open(root)
	if err != nil {
		return nil
	}
	cs, err := messages.For(repo, story)
	if err != nil {
		return nil
	}
	holds := repo.Holds(nil)
	var out []harness.Share
	for _, c := range cs {
		if closed, _ := c.Closed(repo); closed {
			continue
		}
		for _, s := range c.Shares {
			if workitem.CanonicalID(s.Held) == story && holds.InForce(s) {
				out = append(out, harness.Share{Conversation: c.ID, Holder: workitem.CanonicalID(s.Holder), Paths: s.Paths, Split: s.Split})
			}
		}
	}
	return out
}

// openLog opens the log of a run for subject, a story or a strategic
// agent's kind, a built-in's or a defined role's, started at: beside the
// others in serve/agents, named for the project's key, the subject, and the
// time, which Dir.AgentLogs and Dir.ActivityLogs read by. A story's runs
// append to the log of the second they started in. A strategic agent's runs,
// whose logs only the time tells apart, each get a log of their own
// (S-0208): named for the first second from their start that no other run's
// log is, within a minute.
func (l *launcher) openLog(subject string, at time.Time) (*os.File, error) {
	dir := filepath.Join(string(l.dir), "agents")
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return nil, err
	}
	name := func(t time.Time) string {
		return filepath.Join(dir, fmt.Sprintf("%s-%s-%s.log", l.entry.Key, subject, t.Format("20060102T150405Z")))
	}
	if workitem.TypeOfID(subject) == workitem.Story {
		return os.OpenFile(name(at), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o600)
	}
	for i := 0; ; i++ {
		f, err := os.OpenFile(name(at.Add(time.Duration(i)*time.Second)), os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o600)
		if !errors.Is(err, os.ErrExist) || i == 59 {
			return f, err
		}
	}
}

// spawnWith starts argv as run's process, in the project's directory, its
// output in the log openLog opens for subject, and its environment flai
// serve's without the host's address, token, and config, the providers'
// spend keys, and the forges' tokens, with FLAI_AGENT,
// the session started at, FLAI_STARTED_BY, and env. It sets run's command,
// log, PID, and process start, and has the launcher wait for it unless it
// hands its agents over. It records nothing: the caller does, as started or
// as failed with the error. env may carry a provider's key (startOf), so it
// goes into the child's environment and nowhere else: not the log, run, or
// an error. stdin is written to the child's standard input, then closed: the
// container harness's start request for pilot (S-0508). With no stdin the
// child's standard input is the null device.
func (l *launcher) spawnWith(ctx context.Context, run *AgentRun, argv []string, subject string, at time.Time, env []string, stdin []byte) (*exec.Cmd, *os.File, error) {
	run.Command = filepath.Base(argv[0])
	out, err := l.openLog(subject, at)
	if err != nil {
		return nil, nil, err
	}
	run.Log = out.Name()
	// It outlives the look that started it and flai serve itself: a session
	// at work is not ended because its starter is restarted.
	cmd := exec.CommandContext(context.WithoutCancel(ctx), argv[0], argv[1:]...)
	cmd.Dir = l.entry.Root
	cmd.Stdout, cmd.Stderr = out, out
	cmd.Env = append(l.withoutForgeTokens(l.withoutSpendKeys(agentEnv(os.Environ()))), "FLAI_AGENT="+run.Agent, "FLAI_SESSION="+at.Format("20060102T150405"), "FLAI_STARTED_BY=flai-serve")
	cmd.Env = append(cmd.Env, env...)
	Detach(cmd)
	var feed *os.File
	if len(stdin) > 0 {
		r, w, err := os.Pipe()
		if err != nil {
			_ = out.Close()
			return nil, nil, fmt.Errorf("a pipe for the start request on the agent's standard input: %w", err)
		}
		cmd.Stdin, feed = r, w
		defer func() { _ = r.Close() }() // the child holds its own copy once started
	}
	if err := cmd.Start(); err != nil {
		_ = out.Close()
		if feed != nil {
			_ = feed.Close()
		}
		var ee *exec.Error
		if errors.As(err, &ee) {
			err = fmt.Errorf("%s: %w", ee.Name, ee.Err)
		}
		return nil, nil, err
	}
	if feed != nil {
		if err := writeStdin(feed, stdin); err != nil {
			// a child without its whole request cannot start the agent
			_ = KillGroup(cmd.Process.Pid)
			_ = cmd.Wait()
			_ = out.Close()
			return nil, nil, fmt.Errorf("write the start request to %s's standard input: %w", run.Command, err)
		}
	}
	run.PID = cmd.Process.Pid
	run.Start = Started(run.PID)
	if !l.handOver {
		l.waiting[run.PID] = true
	}
	return cmd, out, nil
}

// stdinWait bounds the write of a start request to a child's standard
// input: a request fits the pipe's buffer, so only a child that never reads
// a larger one would hold the write.
const stdinWait = 10 * time.Second

// writeStdin writes data to w, the write end of a child's standard input,
// within stdinWait, and closes it, so that the child reads to its end. A
// look, and a start on the operator's word that exits once the agent runs,
// never wait on the child longer than that.
func writeStdin(w *os.File, data []byte) error {
	defer func() { _ = w.Close() }()
	// a pipe that takes no deadline, as on Windows, is written without one
	_ = w.SetWriteDeadline(time.Now().Add(stdinWait))
	_, err := w.Write(data)
	return err
}

// await waits in the background for the process spawnWith started, has end
// record how it ended from its exit code, and has the launcher look again;
// a launcher that hands its agents over records nothing and leaves the run
// for the serving flai to settle once it is gone, but still reaps the
// process in the background, so that a parent that outlives it, as flai mcp
// over stdio does, keeps no zombie the serving flai would count as running
// (I-0299).
func (l *launcher) await(cmd *exec.Cmd, out *os.File, pid int, end func(code int)) {
	if l.handOver {
		_ = out.Close()
		go func() { _ = cmd.Wait() }()
		return
	}
	go func() {
		err := cmd.Wait()
		_ = out.Close()
		code := 0
		var xe *exec.ExitError
		if errors.As(err, &xe) {
			code = xe.ExitCode()
		} else if err != nil {
			code = -1
		}
		end(code)
		l.mu.Lock()
		delete(l.waiting, pid)
		l.mu.Unlock()
		select {
		case l.again <- struct{}{}:
		default:
		}
	}()
}
