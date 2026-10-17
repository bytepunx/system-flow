package hostapi

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/channel"
	"github.com/bytepunx/system-flow/flai/internal/components"
	"github.com/bytepunx/system-flow/flai/internal/manifest"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// fakeServers is a supervisor that starts nothing: it remembers what it was
// asked and answers as a running dev server would, one per sub-project and
// story.
type fakeServers struct {
	mu       sync.Mutex
	running  map[string]bool // story|sub
	started  []string        // root|project|sub|story|dir of each Start
	statuses []string        // sub, or story|sub, of each Status
	startErr error
}

// fakeKey is how fakeServers keys a dev server, and names it in statuses.
func fakeKey(sub, story string) string {
	if story == "" {
		return sub
	}
	return story + "|" + sub
}

func (f *fakeServers) Start(_ context.Context, root, project string, sub manifest.Project, story, dir string) (ComponentServer, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.started = append(f.started, root+"|"+project+"|"+sub.Name+"|"+story+"|"+dir)
	if f.startErr != nil {
		return ComponentServer{}, f.startErr
	}
	k := fakeKey(sub.Name, story)
	already := f.running[k]
	f.running[k] = true
	return ComponentServer{Project: sub.Name, Story: story, Running: true, AlreadyRunning: already, URL: "http://127.0.0.1:5173/__flai/components/", Port: 5173, PID: 42, Since: "2026-10-16T08:00:00Z"}, nil
}

func (f *fakeServers) Stop(_, sub, story string) (ComponentStopped, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	k := fakeKey(sub, story)
	was := f.running[k]
	delete(f.running, k)
	if !was {
		return ComponentStopped{Project: sub, Story: story}, nil
	}
	return ComponentStopped{Project: sub, Story: story, Stopped: true, PID: 42}, nil
}

func (f *fakeServers) Status(root, sub, story string) ComponentServer {
	f.mu.Lock()
	defer f.mu.Unlock()
	k := fakeKey(sub, story)
	f.statuses = append(f.statuses, k)
	if !f.running[k] {
		return ComponentServer{Project: sub, Story: story, LastError: "the Vite dev server in " + filepath.Join(root, "ui") + " exited"}
	}
	return ComponentServer{Project: sub, Story: story, Running: true, URL: "http://127.0.0.1:5173/__flai/components/", Port: 5173, PID: 42, Since: "2026-10-16T08:00:00Z"}
}

func (f *fakeServers) Idle() time.Duration { return 30 * time.Minute }

// shipyard is a project whose manifest names a SvelteKit sub-project, ui,
// with one component, and a Go one, cli.
func shipyard(t *testing.T) channel.Project {
	t.Helper()
	root := t.TempDir()
	m := "version: 1\nname: Shipyard\nkey: shipyard\nowner: olive\nlayout:\n  design: design\n  docs: docs\n  wip: wip\n" +
		"projects:\n  - name: ui\n    path: ui\n    kind: sveltekit\n  - name: cli\n    path: cli\n    kind: go\n"
	if err := os.WriteFile(filepath.Join(root, manifest.File), []byte(m), 0o644); err != nil {
		t.Fatal(err)
	}
	lib := filepath.Join(root, "ui", "src", "lib")
	if err := os.MkdirAll(lib, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(lib, "Button.svelte"), []byte("<button>hi</button>\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	return channel.Project{Key: "shipyard", Name: "Shipyard", Root: root}
}

// componentsHost is a host whose components action is on when on says so,
// journalling into journal, with servers as its supervisor.
func componentsHost(on *bool, journal *[]Entry, servers ComponentServers) Host {
	return Host{
		Enabled:    func(action, _ string) bool { return *on && action == ActionComponents },
		Record:     func(e Entry) { *journal = append(*journal, e) },
		Components: servers,
	}
}

// ask calls method of the table host makes, with params, and decodes a
// successful answer into into, failing when it names a path on the host.
func ask(t *testing.T, host Host, p channel.Project, method, params string, into any) *channel.Error {
	t.Helper()
	m, ok := MethodsFor("test", func() time.Time { return t0 }, host)[method]
	if !ok {
		t.Fatalf("no method %s", method)
	}
	res, rerr := m(context.Background(), p, json.RawMessage(params))
	if rerr != nil {
		return rerr
	}
	data, err := json.Marshal(res)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(data), p.Root) {
		t.Errorf("%s names a path on the host:\n%s", method, data)
	}
	if into != nil {
		if err := json.Unmarshal(data, into); err != nil {
			t.Fatalf("%s: %v\n%s", method, err, data)
		}
	}
	return nil
}

func TestComponentsIsAHostActionTheDashboardSees(t *testing.T) {
	if _, ok := Actions[ActionComponents]; !ok || !DashboardSees(ActionComponents) {
		t.Errorf("components is a host action the dashboard sees and the operator enables with %s", EnableCommand(ActionComponents))
	}
	for _, name := range []string{"components.start", "components.stop"} {
		if sp := componentSpecs()[name]; sp.action != ActionComponents || sp.reads || sp.do == nil {
			t.Errorf("%s is a write under the components action, answered in flai serve: %+v", name, sp)
		}
	}
}

func TestComponentWritesAreRefusedWhileTheActionIsOffAndReadsStillAnswer(t *testing.T) {
	p := shipyard(t)
	on, journal, servers := false, []Entry{}, &fakeServers{running: map[string]bool{}}
	host := componentsHost(&on, &journal, servers)
	for _, method := range []string{"components.start", "components.stop"} {
		e := ask(t, host, p, method, `{"project":"ui",`+rid+`}`, nil)
		if e == nil || e.Code != Disabled || !strings.Contains(e.Message, "flai serve enable components") {
			t.Errorf("%s with the action off: %+v", method, e)
		}
	}
	if len(servers.started) != 0 {
		t.Errorf("nothing starts while the action is off: %v", servers.started)
	}
	if len(journal) != 2 || journal[0].Outcome != "disabled" || journal[0].Action != ActionComponents || journal[0].Method != "components.start" || journal[1].Method != "components.stop" {
		t.Errorf("each refusal is journalled: %+v", journal)
	}
	var st ComponentsStatus
	if e := ask(t, host, p, "components.status", `{}`, &st); e != nil {
		t.Fatal(e)
	}
	if st.Enabled || st.Enable != "flai serve enable components" || len(st.Servers) != 1 || st.Servers[0].Project != "ui" || st.Servers[0].Running {
		t.Errorf("components.status with the action off: %+v", st)
	}
	if st.Servers[0].LastError != "the Vite dev server in ui exited" {
		t.Errorf("the last error names paths relative to the repository: %q", st.Servers[0].LastError)
	}
	withNode(t, &listNode{stdout: `{"components":[{"path":"src/lib/Button.svelte","props":[{"name":"label","kind":"string","required":true}]}]}`})
	var list struct {
		Projects []ComponentsProject `json:"projects"`
	}
	if e := ask(t, host, p, "components.list", `{}`, &list); e != nil {
		t.Fatal(e)
	}
	if len(list.Projects) != 2 || !list.Projects[0].Supported || list.Projects[1].Supported || list.Projects[1].Kind != "go" ||
		len(list.Projects[0].Components) != 1 || list.Projects[0].Components[0].Path != "src/lib/Button.svelte" || list.Projects[0].Components[0].Props[0].Name != "label" {
		t.Errorf("components.list is the catalogue: %+v", list)
	}
}

func TestComponentsStartThenStatusThenStopAreJournalled(t *testing.T) {
	p := shipyard(t)
	on, journal, servers := true, []Entry{}, &fakeServers{running: map[string]bool{}}
	host := componentsHost(&on, &journal, servers)
	var w struct {
		Data ComponentServer `json:"data"`
	}
	// the one supported sub-project is the one meant when none is named
	if e := ask(t, host, p, "components.start", `{`+rid+`}`, &w); e != nil {
		t.Fatal(e)
	}
	if !w.Data.Running || w.Data.URL != "http://127.0.0.1:5173/__flai/components/" || w.Data.Port != 5173 || w.Data.Project != "ui" {
		t.Errorf("components.start answers the running server: %+v", w.Data)
	}
	if len(servers.started) != 1 || servers.started[0] != p.Root+"|Shipyard|ui||"+p.Root {
		t.Errorf("the supervisor is asked for ui of the project at its root, by the manifest's name: %v", servers.started)
	}
	if e := ask(t, host, p, "components.start", `{"project":"ui","request_id":"req-0002-again"}`, &w); e != nil {
		t.Fatal(e)
	}
	if !w.Data.AlreadyRunning {
		t.Errorf("a second start answers the running server: %+v", w.Data)
	}
	var st ComponentsStatus
	if e := ask(t, host, p, "components.status", `{"project":"ui"}`, &st); e != nil {
		t.Fatal(e)
	}
	if !st.Enabled || st.IdleSeconds != 1800 || len(st.Servers) != 1 || !st.Servers[0].Running || st.Servers[0].Port != 5173 || st.Servers[0].Since != "2026-10-16T08:00:00Z" {
		t.Errorf("components.status: %+v", st)
	}
	if len(servers.statuses) != 1 || servers.statuses[0] != "ui" {
		t.Errorf("components.status asks the supervisor, which marks the server used: %v", servers.statuses)
	}
	var stopped struct {
		Data ComponentStopped `json:"data"`
	}
	if e := ask(t, host, p, "components.stop", `{"project":"ui","request_id":"req-0003-stop"}`, &stopped); e != nil {
		t.Fatal(e)
	}
	if !stopped.Data.Stopped || stopped.Data.PID != 42 {
		t.Errorf("components.stop: %+v", stopped.Data)
	}
	if e := ask(t, host, p, "components.stop", `{"project":"ui","request_id":"req-0004-stop"}`, &stopped); e != nil || stopped.Data.Stopped {
		t.Errorf("a stop with none running says so: %+v %+v", stopped.Data, e)
	}
	want := []string{
		"started ui's component dev server at http://127.0.0.1:5173/__flai/components/ (pid 42)",
		"ui's component dev server was running at http://127.0.0.1:5173/__flai/components/ (pid 42)",
		"stopped ui's component dev server (pid 42)",
		"ui's component dev server was not running",
	}
	if len(journal) != len(want) {
		t.Fatalf("each start and stop is journalled: %+v", journal)
	}
	for i, e := range journal {
		if e.Action != ActionComponents || e.Outcome != "done" || e.Detail != want[i] || e.By != "olive" || e.Project != "shipyard" {
			t.Errorf("journal %d: %+v, want detail %q", i, e, want[i])
		}
	}
}

func TestComponentsStartRefusals(t *testing.T) {
	p := shipyard(t)
	on, journal := true, []Entry{}
	servers := &fakeServers{running: map[string]bool{}}
	host := componentsHost(&on, &journal, servers)
	for params, want := range map[string]string{
		`{"project":"cli",` + rid + `}`:   `sub-project "cli" is of kind go`,
		`{"project":"nope",` + rid + `}`:  `no sub-project "nope"`,
		`{"project":["ui"],` + rid + `}`:  "params",
		`{"project":"ui"}`:                "request_id",
		`{"project":"ui","request_id":1}`: "request_id",
	} {
		if e := ask(t, host, p, "components.start", params, nil); e == nil || e.Code != channel.CodeInvalidParams || !strings.Contains(e.Message, want) {
			t.Errorf("components.start %s: %+v, want invalid params saying %q", params, e, want)
		}
	}
	if len(servers.started) != 0 || len(journal) != 0 {
		t.Errorf("a refused call starts nothing and journals nothing: %v %+v", servers.started, journal)
	}

	servers.startErr = fmt.Errorf("cannot serve sub-project ui's components: %w in %s", components.ErrViteMissing, filepath.Join(p.Root, "ui"))
	e := ask(t, host, p, "components.start", `{"project":"ui",`+rid+`}`, nil)
	if e == nil || e.Code != Rule || strings.Contains(e.Message, p.Root) || !strings.Contains(e.Message, "vite is not installed in ui") {
		t.Errorf("a sub-project without vite is the operator's to fix: %+v", e)
	}
	if len(journal) != 1 || journal[0].Outcome != "failed" || !strings.Contains(journal[0].Detail, "vite is not installed") {
		t.Errorf("a failed start is journalled: %+v", journal)
	}

	none := componentsHost(&on, &journal, nil)
	if e := ask(t, none, p, "components.start", `{"project":"ui","request_id":"req-0005-none"}`, nil); e == nil || e.Code != Rule || !strings.Contains(e.Message, "runs no component dev servers") {
		t.Errorf("a flai that runs no dev servers refuses a start: %+v", e)
	}
	var st ComponentsStatus
	if e := ask(t, none, p, "components.status", `{}`, &st); e != nil || st.IdleSeconds != 0 || len(st.Servers) != 1 || st.Servers[0].Running {
		t.Errorf("and its status says none runs: %+v %+v", st, e)
	}
}

func TestComponentReadsRefuseWhatIsNotTheirs(t *testing.T) {
	p := shipyard(t)
	on, journal := true, []Entry{}
	host := componentsHost(&on, &journal, &fakeServers{running: map[string]bool{}})
	for method, params := range map[string][]string{
		"components.status": {`{"project":"cli"}`, `{"project":"nope"}`, `[]`, `{"project":1}`},
		"components.list":   {`{"project":"nope"}`, `[]`},
	} {
		for _, in := range params {
			if e := ask(t, host, p, method, in, nil); e == nil || e.Code != channel.CodeInvalidParams {
				t.Errorf("%s %s: %+v, want invalid params", method, in, e)
			}
		}
	}
	withNode(t, &listNode{missing: true})
	e := ask(t, host, p, "components.list", `{"project":"ui"}`, nil)
	if e == nil || e.Code != Rule || !strings.Contains(e.Message, "install Node.js") {
		t.Errorf("components.list without node says what to install: %+v", e)
	}
	var list struct {
		Projects []ComponentsProject `json:"projects"`
	}
	if e := ask(t, host, p, "components.list", `{"project":"cli"}`, &list); e != nil || len(list.Projects) != 1 || list.Projects[0].Supported || list.Projects[0].Components == nil {
		t.Errorf("an unsupported sub-project is listed by its kind, without node: %+v %+v", list, e)
	}
}

// listNode is node for components.list: it answers stdout, or is missing.
type listNode struct {
	missing bool
	stdout  string
	dirs    []string // the folder of each run
}

func (n *listNode) LookPath(name string) (string, error) {
	if n.missing {
		return "", errors.New("not found")
	}
	return "/usr/bin/" + name, nil
}

func (n *listNode) Run(_ context.Context, dir, _, _ string, _ ...string) ([]byte, []byte, error) {
	n.dirs = append(n.dirs, dir)
	return []byte(n.stdout), nil, nil
}

// withNode has components.list run r for node until the test ends.
func withNode(t *testing.T, r components.Runner) {
	t.Helper()
	was := nodeRunner
	nodeRunner = r
	t.Cleanup(func() { nodeRunner = was })
}

// shipyardStories is shipyard with work items: S-0001 in progress and
// S-0002 in review, each with a worktree, S-0003 in progress without one,
// S-0004 ready and S-0005 done, each with one, and tasks: T-0001 of
// S-0002 and T-0002 of S-0005, done, and T-0003 of S-0001. Each worktree's
// manifest names a second SvelteKit sub-project, kit, that the main
// checkout's does not.
func shipyardStories(t *testing.T) channel.Project {
	t.Helper()
	p := shipyard(t)
	for _, d := range []string{"wip/kanban/epics", "wip/kanban/stories", "wip/kanban/tasks", "wip/agents", "wip/archive", "wip/threads"} {
		if err := os.MkdirAll(filepath.Join(p.Root, d), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	repo, err := workitem.Open(p.Root)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := repo.Create(workitem.NewOptions{Type: workitem.Epic, Title: "Yard", Owner: "olive", Now: t0}); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(filepath.Join(p.Root, manifest.File))
	if err != nil {
		t.Fatal(err)
	}
	branch := string(data) + "  - name: kit\n    path: kit\n    kind: svelte\n"
	for _, st := range []struct {
		states   []string
		worktree bool
	}{
		{[]string{workitem.Ready, workitem.InProgress}, true},
		{[]string{workitem.Ready, workitem.InProgress, workitem.Review}, true},
		{[]string{workitem.Ready, workitem.InProgress}, false},
		{[]string{workitem.Ready}, true},
		{[]string{workitem.Ready, workitem.InProgress, workitem.Review, workitem.Done}, true},
	} {
		it, err := repo.Create(workitem.NewOptions{Type: workitem.Story, Title: "Crane", Parent: "E-0001", Owner: "olive", Now: t0})
		if err != nil {
			t.Fatal(err)
		}
		body, _ := os.ReadFile(it.Path)
		_ = os.WriteFile(it.Path, []byte(strings.Replace(string(body), "## Acceptance criteria\n- [ ]\n", "## Acceptance criteria\n- [x] works\n", 1)), 0o644)
		it, _ = repo.Get(it.ID)
		if slices.Contains(st.states, workitem.Review) {
			// a story goes to review with its tasks done
			task, err := repo.Create(workitem.NewOptions{Type: workitem.Task, Title: "Hoist", Parent: it.ID, Owner: "olive", Now: t0})
			if err != nil {
				t.Fatal(err)
			}
			for i, to := range []string{workitem.Ready, workitem.InProgress, workitem.Done} {
				if _, err := repo.Transition(task, to, "olive", "", t0.Add(time.Duration(i+1)*time.Minute)); err != nil {
					t.Fatalf("%s to %s: %v", task.ID, to, err)
				}
			}
		}
		for i, to := range st.states {
			if _, err := repo.Transition(it, to, "olive", "", t0.Add(time.Duration(i+1)*time.Minute)); err != nil {
				t.Fatalf("%s to %s: %v", it.ID, to, err)
			}
		}
		if st.worktree {
			wt := repo.WorktreePath(it.ID)
			writeFile(t, wt, manifest.File, []byte(branch))
			writeFile(t, wt, "ui/src/lib/Button.svelte", []byte("<button>hi</button>\n"))
			writeFile(t, wt, "kit/src/lib/Card.svelte", []byte("<p>card</p>\n"))
		}
	}
	if _, err := repo.Create(workitem.NewOptions{Type: workitem.Task, Title: "Rig", Parent: "S-0001", Owner: "olive", Now: t0}); err != nil {
		t.Fatal(err)
	}
	return p
}

// S-0442: with a story, components.start, components.status, and
// components.stop ask the supervisor for that story's worktree's dev
// server, beside main's, each answer carries the story, and the journal
// names it.
func TestComponentsWithAStoryServeItsWorktree(t *testing.T) {
	p := shipyardStories(t)
	on, journal, servers := true, []Entry{}, &fakeServers{running: map[string]bool{}}
	host := componentsHost(&on, &journal, servers)
	wt := filepath.Join(p.Root, ".flai-cache", "worktrees", "S-0001")
	var w struct {
		Data ComponentServer `json:"data"`
	}
	if e := ask(t, host, p, "components.start", `{"project":"ui","story":"S-1",`+rid+`}`, &w); e != nil {
		t.Fatal(e)
	}
	if w.Data.Story != "S-0001" || !w.Data.Running {
		t.Errorf("components.start with a story answers its server, with the story: %+v", w.Data)
	}
	// kit is only in the branch's manifest, so it is found in the worktree's
	if e := ask(t, host, p, "components.start", `{"project":"kit","story":"S-0002","request_id":"req-0002-kit"}`, nil); e != nil {
		t.Fatal(e)
	}
	if e := ask(t, host, p, "components.start", `{"project":"ui","request_id":"req-0003-main"}`, nil); e != nil {
		t.Fatal(e)
	}
	want := []string{
		p.Root + "|Shipyard|ui|S-0001|" + wt,
		p.Root + "|Shipyard|kit|S-0002|" + filepath.Join(p.Root, ".flai-cache", "worktrees", "S-0002"),
		p.Root + "|Shipyard|ui||" + p.Root,
	}
	if strings.Join(servers.started, "\n") != strings.Join(want, "\n") {
		t.Errorf("the supervisor is asked for the story's worktree, keyed by the project's root:\n%s\nwant\n%s", strings.Join(servers.started, "\n"), strings.Join(want, "\n"))
	}

	var st ComponentsStatus
	if e := ask(t, host, p, "components.status", `{"story":"S-0001"}`, &st); e != nil {
		t.Fatal(e)
	}
	if len(st.Servers) != 2 || st.Servers[0].Project != "ui" || st.Servers[0].Story != "S-0001" || !st.Servers[0].Running ||
		st.Servers[1].Project != "kit" || st.Servers[1].Story != "S-0001" || st.Servers[1].Running {
		t.Errorf("components.status with a story is its worktree's sub-projects' servers: %+v", st.Servers)
	}
	if st.Servers[1].LastError != "the Vite dev server in ui exited" {
		t.Errorf("a worktree's last error names paths relative to the checkout: %q", st.Servers[1].LastError)
	}
	var mainSt ComponentsStatus
	if e := ask(t, host, p, "components.status", `{}`, &mainSt); e != nil {
		t.Fatal(e)
	}
	if len(mainSt.Servers) != 1 || mainSt.Servers[0].Story != "" || !mainSt.Servers[0].Running {
		t.Errorf("components.status without one is main's: %+v", mainSt.Servers)
	}
	if got := strings.Join(servers.statuses, " "); got != "S-0001|ui S-0001|kit ui" {
		t.Errorf("each status asks for its own checkout's server: %s", got)
	}

	var stopped struct {
		Data ComponentStopped `json:"data"`
	}
	if e := ask(t, host, p, "components.stop", `{"project":"ui","story":"S-0001","request_id":"req-0004-stop"}`, &stopped); e != nil {
		t.Fatal(e)
	}
	if !stopped.Data.Stopped || stopped.Data.Story != "S-0001" || !servers.running["ui"] || servers.running["S-0001|ui"] {
		t.Errorf("components.stop with a story stops its server only: %+v %v", stopped.Data, servers.running)
	}
	details := make([]string, len(journal))
	for i, e := range journal {
		details[i] = e.Detail
	}
	wantJournal := []string{
		"started S-0001's ui component dev server at http://127.0.0.1:5173/__flai/components/ (pid 42)",
		"started S-0002's kit component dev server at http://127.0.0.1:5173/__flai/components/ (pid 42)",
		"started ui's component dev server at http://127.0.0.1:5173/__flai/components/ (pid 42)",
		"stopped S-0001's ui component dev server (pid 42)",
	}
	if strings.Join(details, "\n") != strings.Join(wantJournal, "\n") {
		t.Errorf("the journal names the story:\n%s", strings.Join(details, "\n"))
	}
}

// S-0442: with a story, components.list is the worktree's catalogue, by its
// manifest, read in the worktree.
func TestComponentsListWithAStoryReadsItsWorktree(t *testing.T) {
	p := shipyardStories(t)
	on, journal := true, []Entry{}
	host := componentsHost(&on, &journal, &fakeServers{running: map[string]bool{}})
	node := &listNode{stdout: `{"components":[{"path":"src/lib/Card.svelte","props":[]}]}`}
	withNode(t, node)
	var list struct {
		Projects []ComponentsProject `json:"projects"`
	}
	if e := ask(t, host, p, "components.list", `{"story":"S-0002"}`, &list); e != nil {
		t.Fatal(e)
	}
	if len(list.Projects) != 3 || list.Projects[2].Name != "kit" || !list.Projects[2].Supported {
		t.Errorf("the worktree's manifest names kit: %+v", list.Projects)
	}
	wt := filepath.Join(p.Root, ".flai-cache", "worktrees", "S-0002")
	if len(node.dirs) != 2 || node.dirs[0] != filepath.Join(wt, "ui") || node.dirs[1] != filepath.Join(wt, "kit") {
		t.Errorf("node runs in the worktree's sub-projects: %v", node.dirs)
	}
}

// S-0442, ADR-0145: a story that is not a work item is NotFound; a task, a
// story neither in progress nor in review, and one without a worktree are
// refused, saying how to open it, by every components method, and nothing
// is started or journalled.
func TestComponentsRefuseAStoryWithoutAnOpenWorktree(t *testing.T) {
	p := shipyardStories(t)
	on, journal := true, []Entry{}
	servers := &fakeServers{running: map[string]bool{}}
	host := componentsHost(&on, &journal, servers)
	withNode(t, &listNode{stdout: `{"components":[]}`})
	cases := []struct {
		name, story string
		code        int
		want        string
	}{
		{"no such item", "S-0099", NotFound, "S-0099 is not a work item here"},
		{"not an ID", "../S-0001", channel.CodeInvalidParams, "not an item ID"},
		{"a task", "T-0001", channel.CodeInvalidParams, "T-0001 is a task"},
		{"no worktree", "S-0003", channel.CodeInvalidParams, "S-0003 has no worktree; open it with flai stream open S-0003"},
		{"ready", "S-0004", channel.CodeInvalidParams, "S-0004 is ready: flai serves a story's components only while it is in progress or in review"},
		{"done", "S-0005", channel.CodeInvalidParams, "S-0005 is done"},
	}
	for _, c := range cases {
		for _, method := range []string{"components.start", "components.stop", "components.status", "components.list"} {
			e := ask(t, host, p, method, `{"project":"ui","story":"`+c.story+`",`+rid+`}`, nil)
			if e == nil || e.Code != c.code || !strings.Contains(e.Message, c.want) || strings.Contains(e.Message, p.Root) {
				t.Errorf("%s with %s (%s): %+v, want code %d saying %q", method, c.story, c.name, e, c.code, c.want)
			}
		}
	}
	if len(servers.started) != 0 || len(servers.statuses) != 0 || len(journal) != 0 {
		t.Errorf("a refused story starts, asks, and journals nothing: %v %v %+v", servers.started, servers.statuses, journal)
	}
}

// S-0442: a start beyond components.max_servers is the operator's to
// settle, a rule, with the supervisor's message, which names the ones
// running.
func TestComponentsStartBeyondTheCapIsARule(t *testing.T) {
	p := shipyardStories(t)
	on, journal := true, []Entry{}
	servers := &fakeServers{running: map[string]bool{}}
	host := componentsHost(&on, &journal, servers)
	msg := "cannot start S-0001's ui component dev server: flai runs 3 component dev servers for this project, components.max_servers: main's ui, S-0002's ui, S-0002's kit; stop one first"
	servers.startErr = fmt.Errorf("%s: %w", msg, ErrComponentServersFull)
	e := ask(t, host, p, "components.start", `{"project":"ui","story":"S-0001",`+rid+`}`, nil)
	if e == nil || e.Code != Rule || !strings.Contains(e.Message, msg) {
		t.Errorf("a start beyond the cap: %+v", e)
	}
	servers.startErr = fmt.Errorf("cannot serve sub-project ui's components: cannot install: %w in %s", components.ErrNoLockFile, filepath.Join(p.Root, ".flai-cache", "worktrees", "S-0001", "ui"))
	e = ask(t, host, p, "components.start", `{"project":"ui","story":"S-0001","request_id":"req-0002-lock"}`, nil)
	if e == nil || e.Code != Rule || !strings.Contains(e.Message, "no lock file in ui") {
		t.Errorf("a worktree with no lock file is a rule, named relative to the worktree: %+v", e)
	}
}
