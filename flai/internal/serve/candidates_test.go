package serve

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/channel"
	"github.com/bytepunx/system-flow/flai/internal/channel/channeltest"
	"github.com/bytepunx/system-flow/flai/internal/hostapi"
)

func repoAt(t *testing.T, dir string, manifest bool) {
	t.Helper()
	if err := os.MkdirAll(filepath.Join(dir, ".git"), 0o755); err != nil {
		t.Fatal(err)
	}
	if manifest {
		if err := os.WriteFile(filepath.Join(dir, "system-flow.yaml"), []byte("version: 1\n"), 0o644); err != nil {
			t.Fatal(err)
		}
	}
}

// S-0098: the git repositories under the named folders that are not
// system-flow projects yet, each with a key of its own.
func TestFindCandidates(t *testing.T) {
	root := t.TempDir()
	repoAt(t, filepath.Join(root, "widget"), false)
	repoAt(t, filepath.Join(root, "org", "Web App"), false)
	repoAt(t, filepath.Join(root, "org", "widget"), false)         // same name, another folder
	repoAt(t, filepath.Join(root, "org", "flow"), true)            // a system-flow project already
	repoAt(t, filepath.Join(root, "served"), false)                // served already (a manifest elsewhere, say)
	repoAt(t, filepath.Join(root, "widget", "vendor-copy"), false) // inside a repository: not looked at
	repoAt(t, filepath.Join(root, ".hidden", "x"), false)
	repoAt(t, filepath.Join(root, "node_modules", "y"), false)
	repoAt(t, filepath.Join(root, "a", "b", "c", "too-deep"), false)
	if err := os.MkdirAll(filepath.Join(root, "just-a-folder"), 0o755); err != nil {
		t.Fatal(err)
	}

	got := FindCandidates([]string{root, root}, map[string]bool{filepath.Join(root, "served"): true}, map[string]bool{"web-app": true})
	want := []Candidate{
		{Root: filepath.Join(root, "org", "Web App"), Name: "Web App", Key: "import-web-app-2"},
		{Root: filepath.Join(root, "org", "widget"), Name: "widget", Key: "import-widget"},
		{Root: filepath.Join(root, "widget"), Name: "widget", Key: "import-widget-2"},
	}
	if len(got) != len(want) {
		t.Fatalf("got %+v", got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("%d: got %+v want %+v", i, got[i], want[i])
		}
	}

	// a named folder that is itself a repository is offered
	single := filepath.Join(t.TempDir(), "solo")
	repoAt(t, single, false)
	if got := FindCandidates([]string{single}, nil, nil); len(got) != 1 || got[0].Key != "import-solo" {
		t.Errorf("a repository named directly: %+v", got)
	}
	if got := FindCandidates([]string{filepath.Join(root, "missing")}, nil, nil); len(got) != 0 {
		t.Errorf("a missing folder: %+v", got)
	}
}

// S-0098: a repository under a named folder is offered to the dashboard the
// served projects reach, answers the import methods, and, once imported, is
// served as a project of its own.
func TestARepositoryIsOfferedImportedAndThenServed(t *testing.T) {
	dash := channeltest.New(t, "s3cret")
	dir := DirFor(filepath.Join(t.TempDir(), "config.json"))
	harbour, key := scratchProject(t, "harbour")
	if err := dir.Register(Entry{Key: "harbour", Name: "harbour", Root: harbour, URL: dash.URL, KeyFile: key}); err != nil {
		t.Fatal(err)
	}
	folder := t.TempDir()
	widget := filepath.Join(folder, "widget")
	repoAt(t, widget, false)

	var ranIn []string
	importRun := func(_ context.Context, r hostapi.Run) (hostapi.Ran, error) {
		ranIn = append(ranIn, r.Dir)
		// what flai import would leave: a manifest with the key it was given
		m := "version: 1\nname: widget\nkey: widget\nlayout:\n  design: design\n  docs: docs\n  wip: wip\n"
		if err := os.WriteFile(filepath.Join(r.Dir, "system-flow.yaml"), []byte(m), 0o644); err != nil {
			return hostapi.Ran{}, err
		}
		return hostapi.Ran{Stdout: []byte(`{"key":"widget","commit":{"committed":true,"commit":"abc1234"}}`)}, nil
	}
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() {
		done <- Run(ctx, Options{Dir: dir, Version: "test", Every: 20 * time.Millisecond, WatchEvery: 20 * time.Millisecond,
			ImportRoots: func() []string { return []string{folder} }, ScanEvery: time.Hour, ImportRun: importRun,
			NewClient: func(e Entry, key []byte) *channel.Client {
				return &channel.Client{URL: e.URL, Key: key, Project: channel.Project{Key: e.Key, Name: e.Name, Root: e.Root},
					Methods: hostapi.Methods("test", nil), Version: "test", MinBackoff: 10 * time.Millisecond, MaxBackoff: 40 * time.Millisecond}
			}})
	}()
	t.Cleanup(func() {
		cancel()
		if err := <-done; err != nil {
			t.Errorf("Run: %v", err)
		}
	})

	byKey := map[string]*channeltest.Conn{}
	for range 2 {
		c := dash.Wait(t)
		byKey[c.Project] = c
	}
	offered := byKey["import-widget"]
	if byKey["harbour"] == nil || byKey["harbour"].Kind != "" || offered == nil || offered.Kind != channel.KindCandidate {
		t.Fatalf("connections: %+v", byKey)
	}
	// it answers the import methods, and nothing of a project's
	if _, e := offered.Ask(t, "board.get", `{"project":"import-widget"}`); e == nil || e.Code != channel.CodeMethodNotFound {
		t.Errorf("board.get on a candidate: %+v", e)
	}
	res, e := offered.Ask(t, "import.run", `{"project":"import-widget","request_id":"req-00000001"}`)
	if e != nil || !strings.Contains(string(res), "abc1234") || len(ranIn) != 1 || ranIn[0] != widget {
		t.Fatalf("import.run: %s %+v %v", res, e, ranIn)
	}

	// registered, and served as a project from the next look, at once
	served := dash.Wait(t)
	if served.Project != "widget" || served.Kind != "" {
		t.Errorf("after the import: %q %q", served.Project, served.Kind)
	}
	projects, _ := dir.Projects()
	var found bool
	for _, p := range projects {
		found = found || (p.Root == widget && p.Key == "widget" && p.URL == dash.URL && p.KeyFile == key)
	}
	if !found {
		t.Errorf("registry: %+v", projects)
	}
	waitFor(t, "no longer offered", func() bool {
		st, _ := dir.ReadStatus(time.Now())
		return len(st.Offered) == 0
	})
}

// S-0101: with no project served, a dashboard flai dashboard recorded
// outside any project still has the repositories offered to it.
func TestARepositoryIsOfferedToARecordedDashboardWithNoProjectServed(t *testing.T) {
	dash := channeltest.New(t, "s3cret")
	dir := DirFor(filepath.Join(t.TempDir(), "config.json"))
	keyFile := filepath.Join(t.TempDir(), "agent.key")
	if err := os.WriteFile(keyFile, []byte("s3cret\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := dir.AddDashboard(Dashboard{URL: dash.URL, KeyFile: keyFile}); err != nil {
		t.Fatal(err)
	}
	folder := t.TempDir()
	repoAt(t, filepath.Join(folder, "widget"), false)
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() {
		done <- Run(ctx, Options{Dir: dir, Version: "test", Every: 20 * time.Millisecond,
			ImportRoots: func() []string { return []string{folder} }, ScanEvery: time.Hour,
			NewClient: func(e Entry, key []byte) *channel.Client {
				return &channel.Client{URL: e.URL, Key: key, Project: channel.Project{Key: e.Key, Name: e.Name, Root: e.Root},
					Version: "test", MinBackoff: 10 * time.Millisecond, MaxBackoff: 40 * time.Millisecond}
			}})
	}()
	t.Cleanup(func() {
		cancel()
		if err := <-done; err != nil {
			t.Errorf("Run: %v", err)
		}
	})
	if c := dash.Wait(t); c.Project != "import-widget" || c.Kind != channel.KindCandidate {
		t.Errorf("offered: %q %q", c.Project, c.Kind)
	}
	// and a dashboard forgotten (stopped) is no longer offered anything new
	if err := dir.ForgetDashboard(dash.URL); err != nil {
		t.Fatal(err)
	}
	if dbs, _ := dir.Dashboards(); len(dbs) != 0 {
		t.Errorf("dashboards after forgetting: %+v", dbs)
	}
}

// S-0102: flai serve started in a folder that is not a project serves the
// projects below it, for the dashboard the registered ones reach, and offers
// the folder's other repositories for import, as if the folder were named;
// nothing of the folder is written to the registry.
func TestServeStartedInAFolderServesAndOffersWhatIsBelowIt(t *testing.T) {
	dash := channeltest.New(t, "s3cret")
	dir := DirFor(filepath.Join(t.TempDir(), "config.json"))
	harbour, key := scratchProject(t, "harbour")
	if err := dir.Register(Entry{Key: "harbour", Name: "harbour", Root: harbour, URL: dash.URL, KeyFile: key}); err != nil {
		t.Fatal(err)
	}
	folder := t.TempDir()
	flow := filepath.Join(folder, "org", "flow")
	repoAt(t, flow, false)
	m := "version: 1\nname: Flow\nkey: flow\nlayout:\n  design: design\n  docs: docs\n  wip: wip\n"
	if err := os.WriteFile(filepath.Join(flow, "system-flow.yaml"), []byte(m), 0o644); err != nil {
		t.Fatal(err)
	}
	repoAt(t, filepath.Join(folder, "plain"), false)

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() {
		done <- Run(ctx, Options{Dir: dir, Version: "test", Every: 20 * time.Millisecond, WatchEvery: 20 * time.Millisecond,
			Folder: folder, ScanEvery: time.Hour,
			NewClient: func(e Entry, key []byte) *channel.Client {
				return &channel.Client{URL: e.URL, Key: key, Project: channel.Project{Key: e.Key, Name: e.Name, Root: e.Root},
					Methods: hostapi.Methods("test", nil), Version: "test", MinBackoff: 10 * time.Millisecond, MaxBackoff: 40 * time.Millisecond}
			}})
	}()
	t.Cleanup(func() {
		cancel()
		if err := <-done; err != nil {
			t.Errorf("Run: %v", err)
		}
	})

	got := map[string]string{}
	for range 3 {
		c := dash.Wait(t)
		got[c.Project] = c.Kind
	}
	if kind, ok := got["flow"]; !ok || kind != "" {
		t.Errorf("the folder's project is served: %v", got)
	}
	if kind := got["import-plain"]; kind != channel.KindCandidate {
		t.Errorf("the folder's repository is offered: %v", got)
	}
	if _, ok := got["harbour"]; !ok {
		t.Errorf("the registered project: %v", got)
	}
	projects, _ := dir.Projects()
	if len(projects) != 1 {
		t.Errorf("the folder's projects were written to the registry: %+v", projects)
	}
	waitFor(t, "the status names the folder", func() bool {
		st, _ := dir.ReadStatus(time.Now())
		return st.Folder == folder && len(st.FolderProjects) == 1 && st.FolderProjects[0].Key == "flow"
	})
}

// With no dashboard known, there is nowhere to serve the folder's projects.
func TestAFolderWithNoDashboardServesNothing(t *testing.T) {
	f := &offers{o: Options{Dir: DirFor(filepath.Join(t.TempDir(), "config.json")), Folder: t.TempDir(), Now: time.Now}}
	if got := f.folderProjects(nil); len(got) != 0 {
		t.Errorf("served: %+v", got)
	}
}
