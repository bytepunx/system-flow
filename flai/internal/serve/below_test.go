package serve

import (
	"context"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/channel"
	"github.com/bytepunx/system-flow/flai/internal/channel/channeltest"
	"github.com/bytepunx/system-flow/flai/internal/hostapi"
)

// projectAt writes a system-flow.yaml with key (none when empty) in dir, a
// git repository when git is true.
func projectAt(t *testing.T, dir, key string, git bool) {
	t.Helper()
	if git {
		repoAt(t, dir, false)
	} else if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	m := "version: 1\nname: " + filepath.Base(dir) + "\n"
	if key != "" {
		m += "key: " + key + "\n"
	}
	m += "layout:\n  design: design\n  docs: docs\n  wip: wip\n"
	if err := os.WriteFile(filepath.Join(dir, "system-flow.yaml"), []byte(m), 0o644); err != nil {
		t.Fatal(err)
	}
}

// S-0117: the git repositories with a system-flow.yaml below a folder named
// for import are found, as the projects below the folder flai serve was
// started in are, each once and with why its manifest cannot serve it.
func TestFindBelow(t *testing.T) {
	folder, named := t.TempDir(), t.TempDir()
	projectAt(t, filepath.Join(folder, "flow"), "flow", false) // below the folder: no git needed
	projectAt(t, filepath.Join(named, "blog"), "blog", true)
	projectAt(t, filepath.Join(named, "org", "nokey"), "", true)
	projectAt(t, filepath.Join(named, "fixture"), "fx", false) // not a repository: a test fixture, say
	repoAt(t, filepath.Join(named, "plain"), false)            // a candidate, not a project
	repoAt(t, filepath.Join(named, "broken"), false)
	if err := os.WriteFile(filepath.Join(named, "broken", "system-flow.yaml"), []byte("version: [\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	got := FindBelow(folder, []string{named, folder})
	roots := make([]string, len(got))
	for i, f := range got {
		roots[i] = f.Root
	}
	want := []string{filepath.Join(named, "blog"), filepath.Join(named, "broken"), filepath.Join(named, "org", "nokey"), filepath.Join(folder, "flow")}
	sort.Strings(want)
	if strings.Join(roots, "|") != strings.Join(want, "|") {
		t.Fatalf("roots:\n got %v\nwant %v", roots, want)
	}
	by := map[string]Found{}
	for _, f := range got {
		by[filepath.Base(f.Root)] = f
	}
	if f := by["blog"]; f.Key != "blog" || !f.Imported || f.Reason != "" {
		t.Errorf("blog: %+v", f)
	}
	if f := by["flow"]; f.Key != "flow" || f.Imported {
		t.Errorf("flow: %+v", f)
	}
	if f := by["broken"]; !strings.Contains(f.Reason, "system-flow.yaml does not load") {
		t.Errorf("broken: %+v", f)
	}
}

// Place serves what it can for the first dashboard and says why it does not
// serve the rest.
func TestPlace(t *testing.T) {
	found := []Found{
		{Entry: Entry{Root: "/a/registered", Key: "reg"}},
		{Entry: Entry{Root: "/a/blog", Key: "blog"}, Imported: true},
		{Entry: Entry{Root: "/a/nokey"}, Imported: true},
		{Entry: Entry{Root: "/a/taken", Key: "harbour"}, Imported: true},
		{Entry: Entry{Root: "/a/twin", Key: "blog"}, Imported: true},
		{Entry: Entry{Root: "/a/broken"}, Reason: "its system-flow.yaml does not load: x"},
	}
	registered := []Entry{{Root: "/a/registered", Key: "reg"}, {Root: "/h/harbour", Key: "harbour", URL: "http://b", KeyFile: "kb"}}
	dashboards := map[string]Entry{"http://b": {URL: "http://b", KeyFile: "kb"}, "http://a": {URL: "http://a", KeyFile: "ka"}}

	got := Place(found, registered, dashboards)
	reasons := map[string]string{}
	for _, f := range got {
		reasons[f.Root] = f.Reason
		if f.Root == "/a/blog" && (f.URL != "http://a" || f.KeyFile != "ka") {
			t.Errorf("blog served for %s %s, want the first dashboard", f.URL, f.KeyFile)
		}
	}
	if _, ok := reasons["/a/registered"]; ok {
		t.Error("a registered project is the registry's")
	}
	for root, want := range map[string]string{
		"/a/blog":   "",
		"/a/nokey":  "has no key",
		"/a/taken":  "its key harbour is served already, for /h/harbour",
		"/a/twin":   "its key blog is served already, for /a/blog",
		"/a/broken": "does not load",
	} {
		r, ok := reasons[root]
		if !ok || (want == "" && r != "") || !strings.Contains(r, want) {
			t.Errorf("%s: reason %q, want %q", root, r, want)
		}
	}

	for _, f := range Place(found[1:2], nil, nil) {
		if !strings.Contains(f.Reason, "no dashboard") {
			t.Errorf("with no dashboard: %+v", f)
		}
	}
}

// S-0117: a repository imported on the command line, below a folder named
// for import and not registered, is served, and the status says so, and
// says why a project there is not.
func TestServeServesTheProjectsBelowTheImportRoots(t *testing.T) {
	dash := channeltest.New(t, "s3cret")
	dir := DirFor(filepath.Join(t.TempDir(), "config.json"))
	harbour, key := scratchProject(t, "harbour")
	if err := dir.Register(Entry{Key: "harbour", Name: "harbour", Root: harbour, URL: dash.URL, KeyFile: key}); err != nil {
		t.Fatal(err)
	}
	named := t.TempDir()
	projectAt(t, filepath.Join(named, "blog"), "blog", true)
	projectAt(t, filepath.Join(named, "nokey"), "", true)

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() {
		done <- Run(ctx, Options{Dir: dir, Version: "test", Every: 20 * time.Millisecond, WatchEvery: 20 * time.Millisecond,
			ImportRoots: func() []string { return []string{named} }, ScanEvery: time.Hour,
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
	for range 2 {
		c := dash.Wait(t)
		got[c.Project] = c.Kind
	}
	if kind, ok := got["blog"]; !ok || kind != "" {
		t.Errorf("the imported repository is served as a project: %v", got)
	}
	if projects, _ := dir.Projects(); len(projects) != 1 {
		t.Errorf("written to the registry: %+v", projects)
	}
	waitFor(t, "the status lists the served and the unserved", func() bool {
		st, _ := dir.ReadStatus(time.Now())
		return len(st.ImportProjects) == 1 && st.ImportProjects[0].Key == "blog" &&
			len(st.Unserved) == 1 && st.Unserved[0].Root == filepath.Join(named, "nokey") && strings.Contains(st.Unserved[0].Reason, "no key") &&
			len(st.Offered) == 0
	})
}
