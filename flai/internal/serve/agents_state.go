package serve

import (
	"encoding/json"
	"maps"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/bytepunx/system-flow/flai/internal/hostapi"
	"github.com/bytepunx/system-flow/flai/internal/roles"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// AgentState is what is known of agents for one project.
type AgentState struct {
	// Running is the newest run that has not ended, and Last the newest that
	// has: what dashboards before S-0104 show.
	Running *AgentRun `json:"running,omitempty"`
	Last    *AgentRun `json:"last,omitempty"`
	// Waiting says why a ready story has not had an agent started for it.
	Waiting string `json:"waiting,omitempty"`
	// Stories are each story's newest run, by ID (S-0104).
	Stories map[string]*AgentRun `json:"stories,omitempty"`
	// Plans are the newest planner run for each item it planned, by the
	// item's ID (S-0208). They are no story's runs: Running and Last never
	// name one.
	Plans map[string]*AgentRun `json:"plans,omitempty"`
	// Orchestrator is the project's newest orchestrator run (S-0218). It is
	// no story's run either.
	Orchestrator *AgentRun `json:"orchestrator,omitempty"`
	// Analyzer is the project's newest analyzer run (S-0223), no story's
	// run either.
	Analyzer *AgentRun `json:"analyzer,omitempty"`
	// Curator is the project's newest curator run (S-0431, ADR-0196), the
	// newest alone, no story's run either.
	Curator *AgentRun `json:"curator,omitempty"`
	// Roles are the newest run of each role the project defines, by the
	// role's name, then by the item it ran on, "" for a run on none (S-0417,
	// ADR-0207). They are no story's runs either.
	Roles map[string]map[string]*AgentRun `json:"roles,omitempty"`
	// PastOrchestrators and PastAnalyzers are the runs Orchestrator and
	// Analyzer named before, newest first, at most pastRuns of each
	// (S-0228): what the dashboard's pages list as past runs.
	PastOrchestrators []*AgentRun `json:"past_orchestrators,omitempty"`
	PastAnalyzers     []*AgentRun `json:"past_analyzers,omitempty"`
	// PastPlans and PastCurators are the planner runs a newer run on the same
	// item replaced in Plans, and the curator runs Curator named before, and
	// PastRoles the runs of each role the project defines that a newer run on
	// the same item replaced in Roles, by the role's name: each newest first,
	// at most pastRuns, kept as PastOrchestrators are, so that agent.status
	// lists every role's past runs (S-0420).
	PastPlans    []*AgentRun            `json:"past_plans,omitempty"`
	PastCurators []*AgentRun            `json:"past_curators,omitempty"`
	PastRoles    map[string][]*AgentRun `json:"past_roles,omitempty"`
}

// pastRuns is how many of a role's earlier runs are kept: the orchestrator's
// and the analyzer's (S-0228), and the planner's, the curator's, and each
// defined role's (S-0420).
const pastRuns = 20

// keep is past with was, the run a new one replaces, put first and the
// oldest let go beyond pastRuns; past as it is when was is nil or is that
// new run.
func keep(past []*AgentRun, was, r *AgentRun) []*AgentRun {
	if was == nil || was.same(r) {
		return past
	}
	out := append([]*AgentRun{was}, past...)
	if len(out) > pastRuns {
		out = out[:pastRuns]
	}
	return out
}

// RoleState is one role as agent.status carries it under roles, by its name,
// with its runs as flai serve records them (S-0420).
type RoleState = hostapi.RoleState[AgentRun]

// RoleStates are the roles files gives, by name, each as st records its
// runs, with enabled saying whether a host action is on for the project
// (S-0420): its runs under way newest first, Run the first of them, and its
// ended runs newest first, at most pastRuns; held on the orchestrator alone,
// while the operator holds it stopped (S-0228). A built-in role's runs are those st keeps for it apart, the
// planner's in Plans, the orchestrator's, the analyzer's, and the curator's
// each in its own field, with their past runs; any other role's are under
// Roles and PastRoles by its name.
func RoleStates(st AgentState, files []roles.Role, enabled func(action string) bool) map[string]RoleState {
	out := make(map[string]RoleState, len(files))
	for _, f := range files {
		r := RoleState{Name: f.Name, Role: f.Role, Path: f.Path, Item: slices.Clone(f.Item), Trigger: []string{}, Running: []*AgentRun{}, Runs: []*AgentRun{}}
		r.Action = hostapi.RoleAction{Name: f.Action, Enabled: f.Action != "" && enabled != nil && enabled(f.Action)}
		if r.Item == nil {
			r.Item = []string{}
		}
		for _, t := range f.Trigger {
			if !slices.Contains(r.Trigger, t.Kind) {
				r.Trigger = append(r.Trigger, t.Kind)
			}
		}
		for _, run := range st.roleRuns(f.Name) {
			if run.live() {
				r.Running = append(r.Running, run)
			} else if len(r.Runs) < pastRuns {
				r.Runs = append(r.Runs, run)
			}
		}
		if len(r.Running) > 0 {
			r.Run = r.Running[0]
		}
		r.Held = f.Name == roles.Orchestrator && st.Orchestrator != nil && st.Orchestrator.Held
		out[f.Name] = r
	}
	return out
}

// roleRuns are every run st records of the role name, newest first by when
// each started, then by item, each once: a built-in role's from the fields
// kept for it, any other's from Roles and PastRoles (S-0420).
func (s AgentState) roleRuns(name string) []*AgentRun {
	var newest, past []*AgentRun
	switch name {
	case workitem.ActivityPlanner:
		newest, past = slices.Collect(maps.Values(s.Plans)), s.PastPlans
	case workitem.ActivityOrchestrator:
		newest, past = []*AgentRun{s.Orchestrator}, s.PastOrchestrators
	case workitem.ActivityAnalyzer:
		newest, past = []*AgentRun{s.Analyzer}, s.PastAnalyzers
	case workitem.ActivityCurator:
		newest, past = []*AgentRun{s.Curator}, s.PastCurators
	default:
		newest, past = slices.Collect(maps.Values(s.Roles[name])), s.PastRoles[name]
	}
	var out []*AgentRun
	for _, r := range append(newest, past...) {
		if r != nil && !slices.ContainsFunc(out, r.same) {
			out = append(out, r)
		}
	}
	slices.SortStableFunc(out, func(a, b *AgentRun) int {
		if c := strings.Compare(b.Started, a.Started); c != 0 {
			return c
		}
		return strings.Compare(a.Item, b.Item)
	})
	return out
}

// of is the newest run recorded for what r is for: its role's on its item,
// for a role the project defines, its item's planner run, the
// orchestrator's, the analyzer's, the curator's, or its story's agent's run.
func (s *AgentState) of(r *AgentRun) *AgentRun {
	if r.defined() {
		return s.Roles[r.RoleName][r.Item]
	}
	if r.curates() {
		return s.Curator
	}
	if r.orchestrates() {
		return s.Orchestrator
	}
	if r.analyzes() {
		return s.Analyzer
	}
	if r.Item != "" {
		return s.Plans[r.Item]
	}
	return s.Stories[r.Story]
}

// put records a run as its story's newest, and keeps Running and Last true,
// a planner run as its item's newest, the orchestrator's, the analyzer's, or
// the curator's as the newest, or a defined role's as its newest on its
// item; the run it replaces, but a story's agent's, is kept among its role's
// past runs (S-0228, S-0420).
func (s *AgentState) put(r *AgentRun) {
	if r.defined() {
		if s.Roles == nil {
			s.Roles = map[string]map[string]*AgentRun{}
		}
		if s.Roles[r.RoleName] == nil {
			s.Roles[r.RoleName] = map[string]*AgentRun{}
		}
		if past := keep(s.PastRoles[r.RoleName], s.Roles[r.RoleName][r.Item], r); len(past) > 0 {
			if s.PastRoles == nil {
				s.PastRoles = map[string][]*AgentRun{}
			}
			s.PastRoles[r.RoleName] = past
		}
		s.Roles[r.RoleName][r.Item] = r
		return
	}
	if r.curates() {
		s.PastCurators = keep(s.PastCurators, s.Curator, r)
		s.Curator = r
		return
	}
	if r.orchestrates() {
		s.PastOrchestrators = keep(s.PastOrchestrators, s.Orchestrator, r)
		s.Orchestrator = r
		return
	}
	if r.analyzes() {
		s.PastAnalyzers = keep(s.PastAnalyzers, s.Analyzer, r)
		s.Analyzer = r
		return
	}
	if r.Item != "" {
		if s.Plans == nil {
			s.Plans = map[string]*AgentRun{}
		}
		s.PastPlans = keep(s.PastPlans, s.Plans[r.Item], r)
		s.Plans[r.Item] = r
		return
	}
	if s.Stories == nil {
		s.Stories = map[string]*AgentRun{}
	}
	s.Stories[r.Story] = r
	s.Running = nil
	for _, run := range s.Stories {
		if run.live() && (s.Running == nil || run.Started > s.Running.Started) {
			s.Running = run
		}
	}
	if !r.live() {
		s.Last = r
	}
}

func (d Dir) agents() string { return filepath.Join(string(d), "agents.json") }

// keptAgents is agents.json as last parsed or written, with the file's stat
// at that moment (S-0603).
type keptAgents struct {
	info   os.FileInfo
	states map[string]AgentState
}

var (
	// agentsMu guards kept, and is held only to read or replace an entry,
	// never across a read, a parse, or a write of the file (S-0603).
	agentsMu sync.Mutex
	kept     = map[string]*keptAgents{} // by agents.json's path
	// agentsWriteMu serializes this process's read-modify-write of
	// agents.json; the file lock serializes it with other processes.
	agentsWriteMu sync.Mutex
	// agentsParsed, when set, is called before each parse of agents.json:
	// the tests' count of parses.
	agentsParsed atomic.Pointer[func(path string)]
)

// AgentStates reads the state of every project's agents. It parses
// agents.json again only when the file's size, modification time, or file
// changed since it was last parsed or written, and keeps what it parsed
// (S-0603). What it returns is shared with every other caller: read it,
// never change it. A change goes through updateAgent, on a copy.
func (d Dir) AgentStates() map[string]AgentState {
	path := d.agents()
	info, err := os.Stat(path)
	if err != nil {
		forget(path)
		return map[string]AgentState{}
	}
	agentsMu.Lock()
	k := kept[path]
	agentsMu.Unlock()
	if k != nil && sameFile(k.info, info) {
		return k.states
	}
	// stat before the read: a write during it changes the stat, and the
	// next call parses again
	states, err := parseAgents(path)
	if err != nil {
		return states
	}
	agentsMu.Lock()
	if kept[path] == k { // not replaced meanwhile by a newer read or write
		kept[path] = &keptAgents{info: info, states: states}
	}
	agentsMu.Unlock()
	return states
}

// sameFile is whether two stats of agents.json name the same contents: the
// same file, size, and modification time. Every write flai makes renames a
// new file into place (Dir.write), so each is a new file whatever the
// clock's granularity.
func sameFile(a, b os.FileInfo) bool {
	return os.SameFile(a, b) && a.Size() == b.Size() && a.ModTime().Equal(b.ModTime())
}

func forget(path string) {
	agentsMu.Lock()
	delete(kept, path)
	agentsMu.Unlock()
}

// parseAgents reads and parses agents.json, empty when it cannot be read and
// as far as it parses when it is not valid, as before S-0603; the error says
// only that it could not be read.
func parseAgents(path string) (map[string]AgentState, error) {
	if parsed := agentsParsed.Load(); parsed != nil {
		(*parsed)(path)
	}
	out := map[string]AgentState{}
	data, err := os.ReadFile(path)
	if err != nil {
		return out, err
	}
	_ = json.Unmarshal(data, &out)
	return out, nil
}

func (d Dir) updateAgent(root string, change func(*AgentState)) {
	agentsWriteMu.Lock()
	defer agentsWriteMu.Unlock()
	path := d.agents()
	// flai serve agent restart writes here from a process of its own (S-0116).
	defer lockFile(path + ".lock")()
	all := d.statesToChange(path)
	st := all[root]
	change(&st)
	all[root] = st
	if err := d.write(path, all); err != nil {
		forget(path)
		return
	}
	info, err := os.Stat(path)
	if err != nil {
		forget(path)
		return
	}
	// a copy: change may have put in runs its caller goes on changing
	agentsMu.Lock()
	kept[path] = &keptAgents{info: info, states: cloneStates(all)}
	agentsMu.Unlock()
}

// statesToChange is agents.json for a change to work on: a copy of what is
// kept while the file is as it was kept, else parsed again. Read under the
// file lock, so no other process writes it meanwhile.
func (d Dir) statesToChange(path string) map[string]AgentState {
	agentsMu.Lock()
	k := kept[path]
	agentsMu.Unlock()
	if info, err := os.Stat(path); err == nil && k != nil && sameFile(k.info, info) {
		return cloneStates(k.states)
	}
	all, _ := parseAgents(path)
	return all
}

// cloneStates copies every project's state deep, so that a change to the
// copy, or to a run put in it, leaves the original as it was.
func cloneStates(all map[string]AgentState) map[string]AgentState {
	out := make(map[string]AgentState, len(all))
	for root, s := range all {
		out[root] = s.clone()
	}
	return out
}

func (s AgentState) clone() AgentState {
	c := s
	c.Running, c.Last = s.Running.clone(), s.Last.clone()
	c.Orchestrator, c.Analyzer, c.Curator = s.Orchestrator.clone(), s.Analyzer.clone(), s.Curator.clone()
	c.Stories, c.Plans = cloneRuns(s.Stories), cloneRuns(s.Plans)
	if s.Roles != nil {
		c.Roles = make(map[string]map[string]*AgentRun, len(s.Roles))
		for name, runs := range s.Roles {
			c.Roles[name] = cloneRuns(runs)
		}
	}
	c.PastOrchestrators, c.PastAnalyzers = cloneList(s.PastOrchestrators), cloneList(s.PastAnalyzers)
	c.PastPlans, c.PastCurators = cloneList(s.PastPlans), cloneList(s.PastCurators)
	if s.PastRoles != nil {
		c.PastRoles = make(map[string][]*AgentRun, len(s.PastRoles))
		for name, runs := range s.PastRoles {
			c.PastRoles[name] = cloneList(runs)
		}
	}
	return c
}

func cloneRuns(runs map[string]*AgentRun) map[string]*AgentRun {
	if runs == nil {
		return nil
	}
	out := make(map[string]*AgentRun, len(runs))
	for k, r := range runs {
		out[k] = r.clone()
	}
	return out
}

func cloneList(runs []*AgentRun) []*AgentRun {
	if runs == nil {
		return nil
	}
	out := make([]*AgentRun, len(runs))
	for i, r := range runs {
		out[i] = r.clone()
	}
	return out
}

// clone copies the run with every field it points to; a field of a kind
// that points is cloned here, which TestAgentRunCloneCoversEveryField holds
// to.
func (r *AgentRun) clone() *AgentRun {
	if r == nil {
		return nil
	}
	c := *r
	if r.Exit != nil {
		exit := *r.Exit
		c.Exit = &exit
	}
	c.Entries, c.Conversations, c.WaitsOn = slices.Clone(r.Entries), slices.Clone(r.Conversations), slices.Clone(r.WaitsOn)
	c.StoryAgent = r.StoryAgent.Clone()
	return &c
}
