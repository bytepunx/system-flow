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
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/channel"
	"github.com/bytepunx/system-flow/flai/internal/hostapi"
	"github.com/bytepunx/system-flow/flai/internal/manifest"
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

// Register adds the project or replaces the entry with its root.
func (d Dir) Register(e Entry) error {
	if e.Root == "" || e.URL == "" || e.KeyFile == "" || e.Key == "" {
		return errors.New("a project needs a root, a key, a dashboard address, and a credential file")
	}
	all, err := d.Projects()
	if err != nil {
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
	// NewClient lets tests shorten a client's timings.
	NewClient func(e Entry, key []byte) *channel.Client
}

type running struct {
	entry  Entry
	client *channel.Client
	stop   context.CancelFunc
	done   chan struct{}
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
	if o.NewClient == nil {
		o.NewClient = func(e Entry, key []byte) *channel.Client {
			return &channel.Client{URL: e.URL, Key: key, Project: channel.Project{Key: e.Key, Name: e.Name, Root: e.Root},
				Methods: hostapi.Methods(o.Version, o.Now), Version: o.Version, Logger: o.Logger}
		}
	}
	if st, alive := o.Dir.ReadStatus(o.Now()); alive && st.PID != os.Getpid() {
		return fmt.Errorf("flai serve is already running (pid %d, since %s)", st.PID, st.Started)
	}
	started := o.Now().UTC().Format(time.RFC3339)
	clients := map[string]*running{}
	var mu sync.Mutex
	defer func() {
		for _, r := range clients {
			r.stop()
			<-r.done
		}
		_ = os.Remove(o.Dir.status())
	}()

	reconcile := func() {
		entries, err := o.Dir.Projects()
		if err != nil {
			o.Logger.Warn("registry unreadable", "component", "serve", "err", err.Error())
			return
		}
		mu.Lock()
		defer mu.Unlock()
		want := map[string]Entry{}
		for _, e := range entries {
			want[e.Root] = e
		}
		for root, r := range clients {
			if e, ok := want[root]; !ok || e != r.entry {
				r.stop()
				<-r.done
				delete(clients, root)
				o.Logger.Info("project dropped", "component", "serve", "root", root)
			}
		}
		for root, e := range want {
			if _, ok := clients[root]; ok {
				continue
			}
			key, err := os.ReadFile(e.KeyFile)
			if err != nil || strings.TrimSpace(string(key)) == "" {
				o.Logger.Warn("agent credential unreadable", "component", "serve", "root", root, "file", e.KeyFile)
				continue
			}
			cctx, stop := context.WithCancel(ctx)
			r := &running{entry: e, client: o.NewClient(e, []byte(strings.TrimSpace(string(key)))), stop: stop, done: make(chan struct{})}
			clients[root] = r
			watcher := &watch.Watcher{Root: e.Root, Paths: watchedPaths(e.Root), Every: o.WatchEvery}
			go func() {
				// The dashboard hears of each changed file; what it shows comes from asking again.
				go watcher.Run(cctx, func(rel string) {
					r.client.Notify("change", map[string]string{"project": e.Key, "path": rel})
				})
				r.client.Run(cctx)
				close(r.done)
			}()
			o.Logger.Info("project served", "component", "serve", "root", root, "key", e.Key, "dashboard", e.URL)
		}
	}
	writeStatus := func() {
		mu.Lock()
		st := Status{PID: os.Getpid(), Version: o.Version, Started: started, Updated: o.Now().UTC().Format(time.RFC3339), Connections: map[string]channel.State{}}
		for root, r := range clients {
			st.Connections[root] = r.client.State()
		}
		mu.Unlock()
		if err := o.Dir.write(o.Dir.status(), st); err != nil {
			o.Logger.Warn("status not written", "component", "serve", "err", err.Error())
		}
	}

	tick := time.NewTicker(o.Every)
	defer tick.Stop()
	for {
		reconcile()
		writeStatus()
		select {
		case <-ctx.Done():
			return nil
		case <-tick.C:
		}
	}
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
