// Package hostapi is what a dashboard can ask flai serve for (ADR-0029):
// the table of named methods, each answered from the same code the CLI
// prints with. Arguments are data to be validated; no method takes a
// command line or a path outside the manifest's folders, and a method that
// is not in this table does not exist for the dashboard.
package hostapi

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"
	"regexp"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/channel"
	"github.com/bytepunx/system-flow/flai/internal/execx"
	"github.com/bytepunx/system-flow/flai/internal/manifest"
	"github.com/bytepunx/system-flow/flai/internal/metrics"
	"github.com/bytepunx/system-flow/flai/internal/perf"
	"github.com/bytepunx/system-flow/flai/internal/release"
	"github.com/bytepunx/system-flow/flai/internal/threads"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// ProjectInfo is what project.info answers: the project as its manifest
// names it, the settings the dashboard's server acts on, and the flai that
// serves it.
type ProjectInfo struct {
	Version     int                `json:"version"`
	Name        string             `json:"name"`
	Key         string             `json:"key"`
	Description string             `json:"description,omitempty"`
	Owner       string             `json:"owner,omitempty"`
	Repo        string             `json:"repo,omitempty"`
	Template    ProjectTemplate    `json:"template"`
	Layout      map[string]string  `json:"layout"`
	Projects    []manifest.Project `json:"projects"`
	Dashboard   ProjectDashboard   `json:"dashboard"`
	Flai        string             `json:"flai"`
	// FlaiOutdated is set while that flai is older than the newest flai
	// release in the project's history (S-0181): the host badge says so.
	FlaiOutdated *release.Outdated `json:"flai_outdated,omitempty"`
	// HostActions are the host actions a dashboard sees, and whether the
	// operator enabled each for this project (S-0078). Read-only to the
	// dashboard.
	HostActions map[string]bool `json:"host_actions"`
	// Agent is the project's default agent, which a story created now gets (S-0103).
	Agent *manifest.Agent `json:"agent,omitempty"`
}

// ProjectTemplate is which template the project was made from, which the
// overview page shows. The repository address is left out: a local template
// path would name a folder on the host.
type ProjectTemplate struct {
	Ref     string `json:"ref,omitempty"`
	Version string `json:"version,omitempty"`
	Applied string `json:"applied,omitempty"`
}

// ProjectDashboard is the part of the manifest's dashboard section that the
// dashboard's server reads.
type ProjectDashboard struct {
	NotifyURL  string `json:"notify_url,omitempty"`
	Autocommit bool   `json:"autocommit"`
}

// ItemWithChildren is what item.get answers: a story with tasks has its
// task plan too (S-0176), and an item with a forecast or an estimate its
// expected cost (ADR-0083).
type ItemWithChildren struct {
	Item         *workitem.Item        `json:"item"`
	Children     []*workitem.Item      `json:"children"`
	Plan         *workitem.Plan        `json:"plan,omitempty"`
	ExpectedCost *metrics.ExpectedCost `json:"expected_cost,omitempty"`
}

// ItemCount is what items.count answers.
type ItemCount struct {
	Active   int `json:"active"`
	Archived int `json:"archived"`
}

var (
	itemID   = regexp.MustCompile(`(?i)^[EST]-?\d{1,6}$`)
	itemType = map[string]bool{"": true, workitem.Epic: true, workitem.Story: true, workitem.Task: true}
	storyID  = regexp.MustCompile(`^S-\d{3,}$`)
	planID   = regexp.MustCompile(`^[ES]-\d{3,}$`)
)

// ErrNoAgent is what Host.AgentStream returns for a story flai serve has
// started no agent for, and Host.PlanStream, wrapped, for an item it has
// started no planner for.
var ErrNoAgent = errors.New("flai serve has started no agent for this story")

func bad(format string, a ...any) *channel.Error {
	return &channel.Error{Code: channel.CodeInvalidParams, Message: fmt.Sprintf(format, a...)}
}

func failed(err error) *channel.Error {
	return &channel.Error{Code: channel.CodeInternal, Message: err.Error()}
}

// Commands runs the flai commands behind the write methods; nil is this
// executable. Tests replace it before building the table.
var Commands Runner

// NotFound is the code for an item or thread that does not exist.
const NotFound = -32004

func params(raw json.RawMessage, into any) *channel.Error {
	if len(raw) == 0 {
		return nil
	}
	if err := json.Unmarshal(raw, into); err != nil {
		return bad("params: %s", err.Error())
	}
	return nil
}

// relative rewrites item paths to be relative to the repository, which is
// how the dashboard links to them and all it needs to know of the host.
func relative(root string, items []*workitem.Item, bodies bool) []*workitem.Item {
	out := make([]*workitem.Item, 0, len(items))
	for _, it := range items {
		c := *it
		if rel, err := filepath.Rel(root, it.Path); err == nil {
			c.Path = filepath.ToSlash(rel)
		}
		if !bodies {
			c.Body = ""
		}
		out = append(out, &c)
	}
	return out
}

// enabledActions says which host actions the operator enabled for a project,
// of those a dashboard sees. It is told to the dashboard so that it can say
// what will happen; nothing the dashboard can ask for changes it.
func enabledActions(host Host, root string) map[string]bool {
	out := map[string]bool{}
	for name := range Actions {
		if DashboardSees(name) {
			out[name] = host.enabled(name, root)
		}
	}
	return out
}

// Methods is the table flai serve offers, with no host action enabled.
func Methods(version string, now func() time.Time) map[string]channel.Method {
	return MethodsFor(version, now, Host{})
}

// MethodsFor is the table with the host's say over host actions.
func MethodsFor(version string, now func() time.Time, host Host) map[string]channel.Method {
	if now == nil {
		now = time.Now
	}
	open := func(ctx context.Context, p channel.Project) (*workitem.Repo, *channel.Error) {
		defer perf.Track(ctx, "repo.open")()
		repo, err := workitem.Open(p.Root)
		if err != nil {
			return nil, failed(err)
		}
		return repo, nil
	}
	table := map[string]channel.Method{
		"project.info": func(ctx context.Context, p channel.Project, _ json.RawMessage) (any, *channel.Error) {
			done := perf.Track(ctx, "manifest.load")
			m, err := manifest.Load(filepath.Join(p.Root, manifest.File))
			done()
			if err != nil {
				return nil, failed(err)
			}
			projects := m.Projects
			if projects == nil {
				projects = []manifest.Project{}
			}
			return ProjectInfo{Version: m.Version, Name: m.Name, Key: m.Key, Description: m.Description, Owner: m.Owner, Repo: m.Repo,
				Template: ProjectTemplate{Ref: m.Template.Ref, Version: m.Template.Version, Applied: m.Template.Applied},
				Layout:   m.Layout, Projects: projects, Dashboard: ProjectDashboard{NotifyURL: m.Dashboard.NotifyURL, Autocommit: m.Autocommit()}, Flai: version,
				HostActions: enabledActions(host, p.Root), Agent: m.Agent,
				FlaiOutdated: release.FlaiOutdated(execx.Timed(ctx, execx.System{}), p.Root, version)}, nil
		},

		// agent.status: whether flai serve starts an agent when a story becomes
		// ready here (S-0079), what it started, and why a ready story waits. It
		// is read-only, like everything about host actions: the dashboard shows
		// it, and cannot start, stop, or configure anything.
		// settings.get: the host's settings as they apply to this project, and
		// whether the dashboard may change them here and host-wide (S-0105).
		// Tokens are never in it: a rotation answers with the new one once.
		"settings.get": func(_ context.Context, p channel.Project, raw json.RawMessage) (any, *channel.Error) {
			if e := params(raw, &struct{}{}); e != nil {
				return nil, e
			}
			out := map[string]any{
				"here": host.enabled(ActionSettings, p.Root), "everywhere": host.enabledEverywhere(ActionSettings),
				"enable": EnableCommand(ActionSettings), "enable_everywhere": EnableEverywhere(ActionSettings),
			}
			if host.Settings != nil {
				out["host"] = host.Settings(p.Root)
			}
			return out, nil
		},

		"agent.status": func(ctx context.Context, p channel.Project, raw json.RawMessage) (any, *channel.Error) {
			if e := params(raw, &struct{}{}); e != nil {
				return nil, e
			}
			out := map[string]any{"enabled": host.enabled(ActionAgent, p.Root)}
			if host.Agent != nil {
				done := perf.Track(ctx, "agent.state")
				out["state"] = host.Agent(p.Root)
				done()
			}
			return out, nil
		},

		// agent.stream: what the newest agent flai serve started for a story
		// said and did, read from its log from the byte offset after, or its
		// tail when after is absent (S-0142); with plan in place of story, the
		// newest planner flai serve started for that epic or story (S-0259);
		// with orchestrator true, the project's newest orchestrator (S-0218).
		// Read-only, like agent.status.
		"agent.stream": func(ctx context.Context, p channel.Project, raw json.RawMessage) (any, *channel.Error) {
			var in struct {
				Story        string `json:"story"`
				Plan         string `json:"plan"`
				Orchestrator bool   `json:"orchestrator"`
				After        *int64 `json:"after"`
			}
			if e := params(raw, &in); e != nil {
				return nil, e
			}
			given := 0
			for _, g := range []bool{in.Story != "", in.Plan != "", in.Orchestrator} {
				if g {
					given++
				}
			}
			if given != 1 {
				return nil, bad("give one of story, for a story's agent, plan, for the planner of an epic or a story, and orchestrator, for the project's orchestrator")
			}
			id, read := in.Story, host.AgentStream
			switch {
			case in.Orchestrator:
				read = nil
				if host.OrchestratorStream != nil {
					read = func(root, _ string, after int64) (any, error) { return host.OrchestratorStream(root, after) }
				}
			case in.Plan != "":
				if !planID.MatchString(in.Plan) {
					return nil, bad("%q is not an epic's or a story's ID", in.Plan)
				}
				id, read = in.Plan, host.PlanStream
			case !storyID.MatchString(in.Story):
				return nil, bad("%q is not a story's ID", in.Story)
			}
			after := int64(-1)
			if in.After != nil {
				if *in.After < 0 {
					return nil, bad("after must be an offset, 0 or more, not %d", *in.After)
				}
				after = *in.After
			}
			if read == nil {
				return nil, &channel.Error{Code: NotFound, Message: "flai serve starts no agents here, so there is no stream to read"}
			}
			done := perf.Track(ctx, "agent.stream")
			out, err := read(p.Root, id, after)
			done()
			if errors.Is(err, ErrNoAgent) {
				return nil, &channel.Error{Code: NotFound, Message: err.Error()}
			}
			if err != nil {
				return nil, failed(err)
			}
			return out, nil
		},

		// board.get: the board as flai board --json gives it. all adds epics and tasks.
		"board.get": func(ctx context.Context, p channel.Project, raw json.RawMessage) (any, *channel.Error) {
			var in struct {
				All bool `json:"all"`
			}
			if e := params(raw, &in); e != nil {
				return nil, e
			}
			repo, e := open(ctx, p)
			if e != nil {
				return nil, e
			}
			done := perf.Track(ctx, "repo.list")
			items, err := repo.List(true)
			done()
			if err != nil {
				return nil, failed(err)
			}
			done = perf.Track(ctx, "board.load")
			board, err := repo.LoadBoard()
			done()
			if err != nil {
				return nil, failed(err)
			}
			done = perf.Track(ctx, "release.pending")
			pending := release.PendingIDs(execx.Timed(ctx, execx.System{}), repo.Root, repo.Manifest, repo)
			done()
			defer perf.Track(ctx, "board.view")()
			return workitem.NewBoardView(items, board, now(), in.All, pending, repo.Manifest.Projects), nil
		},

		// items.list: by type and state; the archive and the bodies on request.
		"items.list": func(ctx context.Context, p channel.Project, raw json.RawMessage) (any, *channel.Error) {
			var in struct {
				Type     string `json:"type"`
				Status   string `json:"status"`
				Archived bool   `json:"archived"`
				Bodies   bool   `json:"bodies"`
			}
			if e := params(raw, &in); e != nil {
				return nil, e
			}
			if !itemType[in.Type] {
				return nil, bad("type must be epic, story, or task, not %q", in.Type)
			}
			if in.Status != "" && !isState(in.Status) {
				return nil, bad("%q is not a state", in.Status)
			}
			repo, e := open(ctx, p)
			if e != nil {
				return nil, e
			}
			done := perf.Track(ctx, "repo.list")
			items, err := repo.List(in.Archived)
			done()
			if err != nil {
				return nil, failed(err)
			}
			keep := items[:0:0]
			for _, it := range items {
				if (in.Type == "" || it.Type == in.Type) && (in.Status == "" || it.Status == in.Status) {
					keep = append(keep, it)
				}
			}
			return relative(repo.MainRoot, keep, in.Bodies), nil
		},

		// items.count: how many items are active and how many archived, for a
		// page that shows the archive's size and not the archive (S-0162).
		"items.count": func(ctx context.Context, p channel.Project, raw json.RawMessage) (any, *channel.Error) {
			if e := params(raw, &struct{}{}); e != nil {
				return nil, e
			}
			repo, e := open(ctx, p)
			if e != nil {
				return nil, e
			}
			defer perf.Track(ctx, "repo.list")()
			items, err := repo.List(true)
			if err != nil {
				return nil, failed(err)
			}
			var n ItemCount
			for _, it := range items {
				if it.Archived {
					n.Archived++
				} else {
					n.Active++
				}
			}
			return n, nil
		},

		// item.get: one item in any padding, with its children, archive included.
		"item.get": func(ctx context.Context, p channel.Project, raw json.RawMessage) (any, *channel.Error) {
			var in struct {
				ID string `json:"id"`
			}
			if e := params(raw, &in); e != nil {
				return nil, e
			}
			if !itemID.MatchString(in.ID) {
				return nil, bad("%q is not a work item ID", in.ID)
			}
			repo, e := open(ctx, p)
			if e != nil {
				return nil, e
			}
			done := perf.Track(ctx, "repo.get")
			it, err := repo.Get(in.ID)
			done()
			if err != nil {
				return nil, &channel.Error{Code: NotFound, Message: err.Error()}
			}
			done = perf.Track(ctx, "repo.list")
			all, err := repo.List(true)
			done()
			if err != nil {
				return nil, failed(err)
			}
			one := relative(repo.MainRoot, []*workitem.Item{it}, true)[0]
			children := workitem.Children(all, it.ID)
			out := ItemWithChildren{Item: one, Children: relative(repo.MainRoot, children, true)}
			if it.Type == workitem.Story {
				out.Plan = workitem.PlanOf(children, it.ID)
			}
			out.ExpectedCost = metrics.ExpectedCostOf(it, metrics.CostPerAgentHour(all))
			return out, nil
		},

		// threads.list: as flai thread list --json; on is a path or an item, all adds the resolved.
		"threads.list": func(ctx context.Context, p channel.Project, raw json.RawMessage) (any, *channel.Error) {
			var in struct {
				On  string `json:"on"`
				All bool   `json:"all"`
			}
			if e := params(raw, &in); e != nil {
				return nil, e
			}
			repo, e := open(ctx, p)
			if e != nil {
				return nil, e
			}
			var list []*threads.Thread
			var err error
			done := perf.Track(ctx, "threads.read")
			if in.On != "" {
				list, err = threads.For(repo, in.On)
			} else {
				list, err = threads.List(repo)
			}
			done()
			if err != nil {
				return nil, failed(err)
			}
			defer perf.Track(ctx, "threads.view")()
			out := []map[string]any{}
			for _, th := range list {
				if in.All || th.Open() {
					out = append(out, threads.View(repo, th))
				}
			}
			return out, nil
		},
	}
	for _, more := range []map[string]channel.Method{docMethods(), peopleMethods(now), searchMethods(), readMethods(now), writeMethods(Commands, now, host)} {
		for name, m := range more {
			table[name] = m
		}
	}
	return table
}

func isState(s string) bool {
	for _, st := range workitem.States {
		if st == s {
			return true
		}
	}
	return false
}
