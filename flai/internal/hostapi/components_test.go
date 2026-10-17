package hostapi

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/channel"
	"github.com/bytepunx/system-flow/flai/internal/components"
	"github.com/bytepunx/system-flow/flai/internal/manifest"
)

// fakeServers is a supervisor that starts nothing: it remembers what it was
// asked and answers as a running dev server would.
type fakeServers struct {
	mu       sync.Mutex
	running  map[string]bool
	started  []string // root|project|sub of each Start
	statuses []string // sub of each Status
	startErr error
}

func (f *fakeServers) Start(_ context.Context, root, project string, sub manifest.Project) (ComponentServer, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.started = append(f.started, root+"|"+project+"|"+sub.Name)
	if f.startErr != nil {
		return ComponentServer{}, f.startErr
	}
	already := f.running[sub.Name]
	f.running[sub.Name] = true
	return ComponentServer{Project: sub.Name, Running: true, AlreadyRunning: already, URL: "http://127.0.0.1:5173/__flai/components/", Port: 5173, PID: 42, Since: "2026-10-16T08:00:00Z"}, nil
}

func (f *fakeServers) Stop(_, sub string) (ComponentStopped, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	was := f.running[sub]
	delete(f.running, sub)
	if !was {
		return ComponentStopped{Project: sub}, nil
	}
	return ComponentStopped{Project: sub, Stopped: true, PID: 42}, nil
}

func (f *fakeServers) Status(root, sub string) ComponentServer {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.statuses = append(f.statuses, sub)
	if !f.running[sub] {
		return ComponentServer{Project: sub, LastError: "the Vite dev server in " + filepath.Join(root, "ui") + " exited"}
	}
	return ComponentServer{Project: sub, Running: true, URL: "http://127.0.0.1:5173/__flai/components/", Port: 5173, PID: 42, Since: "2026-10-16T08:00:00Z"}
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
	if len(servers.started) != 1 || servers.started[0] != p.Root+"|Shipyard|ui" {
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
}

func (n *listNode) LookPath(name string) (string, error) {
	if n.missing {
		return "", errors.New("not found")
	}
	return "/usr/bin/" + name, nil
}

func (n *listNode) Run(context.Context, string, string, string, ...string) ([]byte, []byte, error) {
	return []byte(n.stdout), nil, nil
}

// withNode has components.list run r for node until the test ends.
func withNode(t *testing.T, r components.Runner) {
	t.Helper()
	was := nodeRunner
	nodeRunner = r
	t.Cleanup(func() { nodeRunner = was })
}
