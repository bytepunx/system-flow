package serve

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"sort"
	"strings"
	"sync/atomic"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/channel"
	"github.com/bytepunx/system-flow/flai/internal/hostapi"
	"github.com/bytepunx/system-flow/flai/internal/manifest"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// Candidate is a git repository under a folder the operator named that is
// not a system-flow project yet, which the board offers to import (S-0098).
type Candidate struct {
	Root string `json:"root"`
	Name string `json:"name"`
	// Key is what the connection names it by: CandidatePrefix and the key the
	// project will have once imported.
	Key string `json:"key"`
}

// candidateDepth is how far below a named folder repositories are looked
// for: ~/git/repo and ~/git/org/repo, and one more.
const candidateDepth = 3

var skipScanning = map[string]bool{"node_modules": true, "vendor": true, "dist": true, "build": true, "target": true, "__pycache__": true}

// CandidatePrefix begins every candidate's key; what follows it is the key
// the project is imported with.
const CandidatePrefix = hostapi.CandidatePrefix

// FindCandidates lists the git repositories under roots with no
// system-flow.yaml of their own, other than those already served. A folder
// that is a repository is not looked into further, hidden folders and build
// output never are, and nothing is followed through a symbolic link. Each
// key, after CandidatePrefix, is the folder's name made safe and unique,
// taken keys (the served projects') included: it becomes the project's key.
func FindCandidates(roots []string, served, taken map[string]bool) []Candidate {
	var out []Candidate
	seen := map[string]bool{}
	var walk func(dir string, depth int)
	walk = func(dir string, depth int) {
		if seen[dir] {
			return
		}
		seen[dir] = true
		if isRepository(dir) {
			if !served[dir] && !exists(filepath.Join(dir, manifest.File)) {
				out = append(out, Candidate{Root: dir, Name: filepath.Base(dir)})
			}
			return
		}
		if depth >= candidateDepth {
			return
		}
		entries, err := os.ReadDir(dir)
		if err != nil {
			return
		}
		for _, e := range entries {
			if !e.IsDir() || strings.HasPrefix(e.Name(), ".") || skipScanning[e.Name()] {
				continue
			}
			walk(filepath.Join(dir, e.Name()), depth+1)
		}
	}
	for _, r := range roots {
		if abs, err := filepath.Abs(r); err == nil {
			walk(filepath.Clean(abs), 0)
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Root < out[j].Root })
	used := map[string]bool{}
	for k := range taken {
		used[k] = true
	}
	for i := range out {
		base := slug(out[i].Name)
		key := base
		for n := 2; used[key]; n++ {
			key = fmt.Sprintf("%s-%d", base, n)
		}
		used[key] = true
		out[i].Key = CandidatePrefix + key
	}
	return out
}

func isRepository(dir string) bool { return exists(filepath.Join(dir, ".git")) }

func exists(p string) bool {
	_, err := os.Stat(p)
	return err == nil
}

// slug is name in lower case, with anything but letters and digits a dash.
func slug(name string) string {
	var b strings.Builder
	dash := false
	for _, r := range strings.ToLower(name) {
		if r >= 'a' && r <= 'z' || r >= '0' && r <= '9' {
			b.WriteRune(r)
			dash = false
		} else if !dash && b.Len() > 0 {
			b.WriteByte('-')
			dash = true
		}
	}
	s := strings.TrimSuffix(b.String(), "-")
	if s == "" {
		return "repository"
	}
	return s
}

// scanEvery is how often the named folders are looked through again, when
// Options.ScanEvery is zero; an import asks for a look at once.
const scanEvery = 30 * time.Second

// offers keeps a connection open to each dashboard the served projects reach
// for every candidate (S-0098), answering only the import methods; once one
// is imported it is registered, and served as a project from then on.
type offers struct {
	o       Options
	running map[string]*running // by dashboard address and root
	found   []Candidate
	scanned time.Time
	rescan  atomic.Bool

	// the projects below Options.Folder (S-0102), found again every ScanEvery
	folderFound   []Entry
	folderScanned time.Time
}

func (f *offers) every() time.Duration {
	if f.o.ScanEvery > 0 {
		return f.o.ScanEvery
	}
	return scanEvery
}

// importRoots are the folders named for import, and the folder flai serve
// was started in when it is not a project (S-0102).
func (f *offers) importRoots() []string {
	var roots []string
	if f.o.ImportRoots != nil {
		roots = f.o.ImportRoots()
	}
	if f.o.Folder != "" && !slices.Contains(roots, f.o.Folder) {
		roots = append(roots, f.o.Folder)
	}
	return roots
}

// dashboards are the dashboards served projects dial, and those flai
// dashboard started outside any project (S-0101), by address.
func (f *offers) dashboards(entries []Entry) map[string]Entry {
	out := map[string]Entry{}
	for _, e := range entries {
		out[e.URL] = e
	}
	if recorded, err := f.o.Dir.Dashboards(); err == nil {
		for _, db := range recorded {
			if _, ok := out[db.URL]; !ok {
				out[db.URL] = Entry{URL: db.URL, KeyFile: db.KeyFile}
			}
		}
	}
	return out
}

// folderProjects are the system-flow projects below the folder flai serve was
// started in (S-0102), served as registered ones are, for as long as it runs
// and without being written to the registry. Each is served for the first
// dashboard, by address, that the registered projects or a recorded one
// reach; with none there is nowhere to serve them yet. A project already
// registered, or whose key another has, is left to the registry.
func (f *offers) folderProjects(entries []Entry) []Entry {
	if f.o.Folder == "" {
		return nil
	}
	if now := f.o.Now(); f.folderScanned.IsZero() || now.Sub(f.folderScanned) >= f.every() || f.rescan.Load() {
		f.folderFound = nil
		for _, root := range workitem.FindProjects(f.o.Folder) {
			m, err := manifest.Load(filepath.Join(root, manifest.File))
			if err != nil || m.Key == "" {
				continue
			}
			f.folderFound = append(f.folderFound, Entry{Key: m.Key, Name: m.Name, Root: root})
		}
		f.folderScanned = now
	}
	dashboards := f.dashboards(entries)
	if len(dashboards) == 0 {
		return nil
	}
	urls := make([]string, 0, len(dashboards))
	for u := range dashboards {
		urls = append(urls, u)
	}
	sort.Strings(urls)
	d := dashboards[urls[0]]
	roots, keys := map[string]bool{}, map[string]bool{}
	for _, e := range entries {
		roots[e.Root], keys[e.Key] = true, true
	}
	var out []Entry
	for _, e := range f.folderFound {
		if roots[e.Root] || keys[e.Key] {
			continue
		}
		e.URL, e.KeyFile = d.URL, d.KeyFile
		keys[e.Key] = true
		out = append(out, e)
	}
	return out
}

func (f *offers) reconcile(ctx context.Context, entries []Entry) {
	if f.o.ImportRoots == nil && f.o.Folder == "" {
		return
	}
	if now := f.o.Now(); f.scanned.IsZero() || now.Sub(f.scanned) >= f.every() || f.rescan.Swap(false) {
		served, taken := map[string]bool{}, map[string]bool{}
		for _, e := range entries {
			served[e.Root], taken[e.Key] = true, true
		}
		f.found = FindCandidates(f.importRoots(), served, taken)
		f.scanned = now
	}
	// a candidate has no dashboard of its own
	dashboards := f.dashboards(entries)
	want := map[string]Entry{}
	for url, d := range dashboards {
		for _, c := range f.found {
			want[url+"|"+c.Root] = Entry{Key: c.Key, Name: c.Name, Root: c.Root, URL: url, KeyFile: d.KeyFile}
		}
	}
	for k, r := range f.running {
		if e, ok := want[k]; !ok || e != r.entry {
			r.halt()
			delete(f.running, k)
		}
	}
	for k, e := range want {
		if _, ok := f.running[k]; ok {
			continue
		}
		key, err := os.ReadFile(e.KeyFile)
		if err != nil || strings.TrimSpace(string(key)) == "" {
			continue
		}
		cctx, stop := context.WithCancel(ctx)
		c := f.o.NewClient(e, []byte(strings.TrimSpace(string(key))))
		c.Kind = channel.KindCandidate
		c.Methods = hostapi.ImportMethods(f.o.ImportRun, f.o.Now, f.o.Host, f.imported(e))
		r := &running{entry: e, client: c, stop: stop, done: make(chan struct{})}
		f.running[k] = r
		go func() {
			c.Run(cctx)
			close(r.done)
		}()
		f.o.Logger.Info("repository offered for import", "component", "serve", "root", e.Root, "key", e.Key, "dashboard", e.URL)
	}
}

// imported registers a repository just imported, as flai dashboard would
// have: with the manifest's key and name, for the dashboard it was offered on.
func (f *offers) imported(e Entry) func(root string) {
	return func(root string) {
		m, err := manifest.Load(filepath.Join(root, manifest.File))
		if err != nil || m.Key == "" {
			f.o.Logger.Warn("imported repository not served: no key in its manifest", "component", "serve", "root", root)
			return
		}
		if err := f.o.Dir.Register(Entry{Key: m.Key, Name: m.Name, Root: root, URL: e.URL, KeyFile: e.KeyFile}); err != nil {
			f.o.Logger.Warn("imported repository not registered", "component", "serve", "root", root, "err", err.Error())
			return
		}
		f.rescan.Store(true)
		f.o.Logger.Info("imported repository is served", "component", "serve", "root", root, "key", m.Key)
	}
}

func (f *offers) halt() {
	for _, r := range f.running {
		r.halt()
	}
}
