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

// S-0120: the git repositories with a system-flow.yaml below a folder named
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

	got := Place(found, registered, dashboards, nil)
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

	for _, f := range Place(found[1:2], nil, nil, nil) {
		if !strings.Contains(f.Reason, "no dashboard") {
			t.Errorf("with no dashboard: %+v", f)
		}
	}
}

// S-0123: a project the operator removed is not served from below a folder,
// whether named for import or the one flai serve was started in, and says
// so; a registered one is the registry's whatever the list says.
func TestPlaceLeavesOutWhatTheOperatorRemoved(t *testing.T) {
	found := []Found{
		{Entry: Entry{Root: "/a/blog", Key: "blog"}, Imported: true},
		{Entry: Entry{Root: "/f/notes", Key: "notes"}},
		{Entry: Entry{Root: "/a/registered", Key: "reg"}},
		{Entry: Entry{Root: "/a/kept", Key: "kept"}, Imported: true},
	}
	registered := []Entry{{Root: "/a/registered", Key: "reg"}}
	dashboards := map[string]Entry{"http://a": {URL: "http://a", KeyFile: "ka"}}
	removed := map[string]bool{"/a/blog": true, "/f/notes": true, "/a/registered": true}

	by := map[string]Found{}
	for _, f := range Place(found, registered, dashboards, removed) {
		by[f.Root] = f
	}
	for _, root := range []string{"/a/blog", "/f/notes"} {
		if f := by[root]; !f.Removed || f.URL != "" || !strings.Contains(f.Reason, "removed from the dashboard") || !strings.Contains(f.Reason, "flai serve project add "+root) {
			t.Errorf("%s: %+v", root, f)
		}
	}
	if _, ok := by["/a/registered"]; ok {
		t.Error("a registered project is the registry's")
	}
	if f := by["/a/kept"]; f.Removed || f.Reason != "" || f.URL != "http://a" {
		t.Errorf("kept: %+v", f)
	}
}

func TestTheRemovedListRoundTrips(t *testing.T) {
	dir := DirFor(filepath.Join(t.TempDir(), "config.json"))
	if got, err := dir.RemovedRoots(); err != nil || got != nil {
		t.Fatalf("no list: %v %v", got, err)
	}
	for _, r := range []string{"/b", "/a", "/b"} {
		if err := dir.RemoveRoot(r); err != nil {
			t.Fatal(err)
		}
	}
	if got, _ := dir.RemovedRoots(); strings.Join(got, ",") != "/a,/b" {
		t.Errorf("list %v, want /a,/b once each", got)
	}
	if was, err := dir.RestoreRoot("/a"); !was || err != nil {
		t.Errorf("restore /a: %v %v", was, err)
	}
	if was, err := dir.RestoreRoot("/a"); was || err != nil {
		t.Errorf("restore /a again: %v %v", was, err)
	}
	if set, _ := dir.RemovedSet(); len(set) != 1 || !set["/b"] {
		t.Errorf("set %v, want /b", set)
	}
	_ = os.WriteFile(filepath.Join(string(dir), "removed.json"), []byte("{"), 0o600)
	if set, err := dir.RemovedSet(); err == nil || len(set) != 0 {
		t.Errorf("a list that does not read: %v %v, want empty and an error", set, err)
	}
}

// S-0120: a repository imported on the command line, below a folder named
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
					Methods: hostapi.Methods("test", nil), Version: "test", MinBackoff: 10 * time.Millisecond, MaxBackoff: 40 * time.Millisecond, ReleaseKeys: dash.Keys}
			}})
	}()
	t.Cleanup(func() {
		cancel()
		if err := <-done; err != nil {
			t.Errorf("Run: %v", err)
		}
	})

	got, conns := map[string]string{}, map[string]*channeltest.Conn{}
	for range 2 {
		c := dash.Wait(t)
		got[c.Project], conns[c.Project] = c.Kind, c
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

	// S-0123: removed, it is said to be removed and listed as not served, with why
	blog := filepath.Join(named, "blog")
	if err := dir.RemoveRoot(blog); err != nil {
		t.Fatal(err)
	}
	if got := conns["blog"].Last(t); len(got) != 1 || got[0] != Removed {
		t.Errorf("a removed project sent %v, want [%s]", got, Removed)
	}
	waitFor(t, "the status lists it as removed", func() bool {
		st, _ := dir.ReadStatus(time.Now())
		for _, u := range st.Unserved {
			if u.Root == blog {
				return u.Removed && len(st.ImportProjects) == 0
			}
		}
		return false
	})
	// taken off the list, it is served again
	if _, err := dir.RestoreRoot(blog); err != nil {
		t.Fatal(err)
	}
	for conn := dash.Wait(t); conn.Project != "blog"; conn = dash.Wait(t) {
	}
}

// S-0605: ImportRootOf answers from a root's path what FindBelow's walk of
// the folders named for import finds, for every folder of a tree that holds
// each case the walk treats apart, and for folders that overlap in either
// order.
func TestImportRootOfAgreesWithTheWalk(t *testing.T) {
	tmp := t.TempDir()
	a, b := filepath.Join(tmp, "a"), filepath.Join(tmp, "b")
	for _, dir := range []string{
		filepath.Join(a, "p1"),                 // depth 1
		filepath.Join(a, "org", "p2"),          // depth 2
		filepath.Join(a, "x", "y", "p3"),       // depth 3, as deep as the walk looks
		filepath.Join(a, "x", "y", "z", "p4"),  // depth 4, past it
		filepath.Join(a, ".hidden", "p5"),      // hidden
		filepath.Join(a, "node_modules", "p6"), // skipped
		filepath.Join(a, "outer"),              // a repository
		filepath.Join(a, "outer", "inner"),     // inside a repository
		filepath.Join(tmp, "elsewhere", "p7"),  // only reached through a link
		b,                                      // an import folder that is a repository
		filepath.Join(b, "sub"),                // inside it
	} {
		repoAt(t, dir, true)
	}
	repoAt(t, filepath.Join(a, "bare"), false) // a repository with no system-flow.yaml
	if err := os.Symlink(filepath.Join(tmp, "elsewhere"), filepath.Join(a, "link")); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(filepath.Join(tmp, "elsewhere", "p7"), filepath.Join(a, "p7")); err != nil {
		t.Fatal(err)
	}
	var dirs []string
	_ = filepath.WalkDir(tmp, func(p string, d os.DirEntry, err error) error {
		if err == nil && (d.IsDir() || d.Type()&os.ModeSymlink != 0) && d.Name() != ".git" {
			dirs = append(dirs, p)
		}
		return nil
	})
	dirs = append(dirs, filepath.Join(a, "link", "p7"), "relative/p1", a+"/p1/")

	cases := []struct {
		roots []string
		found []string // what must be found, so the agreement is not empty
		not   []string
	}{
		{[]string{a, b}, []string{filepath.Join(a, "p1"), filepath.Join(a, "org", "p2"), filepath.Join(a, "x", "y", "p3"), filepath.Join(a, "outer"), b},
			[]string{filepath.Join(a, "x", "y", "z", "p4"), filepath.Join(a, ".hidden", "p5"), filepath.Join(a, "node_modules", "p6"),
				filepath.Join(a, "outer", "inner"), filepath.Join(a, "bare"), filepath.Join(a, "link", "p7"), filepath.Join(a, "p7"), filepath.Join(b, "sub")}},
		{[]string{b, a}, []string{filepath.Join(a, "p1"), b}, nil},
		// a/x walked first reaches p4 at depth 3; a's walk then stops at a/x
		{[]string{filepath.Join(a, "x"), a}, []string{filepath.Join(a, "x", "y", "z", "p4"), filepath.Join(a, "p1")}, nil},
		// a walked first enters a/x and stops short of p4; a/x's walk adds nothing
		{[]string{a, filepath.Join(a, "x")}, []string{filepath.Join(a, "x", "y", "p3")}, []string{filepath.Join(a, "x", "y", "z", "p4")}},
	}
	for _, c := range cases {
		listed := map[string]bool{}
		for _, r := range repositories(c.roots) {
			listed[r] = exists(filepath.Join(r, "system-flow.yaml"))
		}
		for _, dir := range dirs {
			want := ""
			if listed[dir] {
				for _, r := range c.roots {
					if within(r, dir) {
						want = r
						break
					}
				}
			}
			if got := ImportRootOf(c.roots, dir); got != want {
				t.Errorf("roots %v, %s: got %q, the walk says %q", c.roots, dir, got, want)
			}
		}
		for _, dir := range c.found {
			if ImportRootOf(c.roots, dir) == "" {
				t.Errorf("roots %v: %s not found", c.roots, dir)
			}
		}
		for _, dir := range c.not {
			if got := ImportRootOf(c.roots, dir); got != "" {
				t.Errorf("roots %v: %s found below %s", c.roots, dir, got)
			}
		}
	}
}
