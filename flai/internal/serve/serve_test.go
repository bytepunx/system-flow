package serve

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/buildinfo"
	"github.com/bytepunx/system-flow/flai/internal/channel"
	"github.com/bytepunx/system-flow/flai/internal/channel/channeltest"
	"github.com/bytepunx/system-flow/flai/internal/hostapi"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

func scratchProject(t *testing.T, key string) (root, keyFile string) {
	t.Helper()
	root = t.TempDir()
	manifest := "version: 1\nname: " + key + "\nkey: " + key + "\nlayout:\n  design: design\n  docs: docs\n  wip: wip\n"
	if err := os.WriteFile(filepath.Join(root, "system-flow.yaml"), []byte(manifest), 0o644); err != nil {
		t.Fatal(err)
	}
	keyFile = filepath.Join(root, "agent.key")
	if err := os.WriteFile(keyFile, []byte("s3cret\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	return root, keyFile
}

func run(t *testing.T, dir Dir, dash *channeltest.Dashboard) {
	t.Helper()
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() {
		done <- Run(ctx, Options{Dir: dir, Version: "test", Every: 20 * time.Millisecond, WatchEvery: 20 * time.Millisecond,
			NewClient: func(e Entry, key []byte) *channel.Client {
				return &channel.Client{URL: e.URL, Key: key, Project: channel.Project{Key: e.Key, Name: e.Name, Root: e.Root},
					Methods: hostapi.Methods("test", nil), Version: "test", PingEvery: time.Hour /* the test dashboard answers no ping: none is sent while the test runs */, MinBackoff: 10 * time.Millisecond, MaxBackoff: 40 * time.Millisecond, ReleaseKeys: dash.Keys}
			}})
	}()
	t.Cleanup(func() {
		cancel()
		if err := <-done; err != nil {
			t.Errorf("Run: %v", err)
		}
	})
}

// S-0237: each client flai serve builds sends this build's release stamp in
// hello and verifies the dashboard's with the embedded release key.
func TestTheDefaultClientCarriesTheStampAndTheReleaseKey(t *testing.T) {
	statement, signature := buildinfo.ReleaseStatement, buildinfo.ReleaseSignature
	buildinfo.ReleaseStatement, buildinfo.ReleaseSignature = "flai 1.43.0 a1825d8", "c2lnbmF0dXJl"
	t.Cleanup(func() { buildinfo.ReleaseStatement, buildinfo.ReleaseSignature = statement, signature })
	want, err := buildinfo.ReleaseKeys()
	if err != nil {
		t.Fatal(err)
	}
	c := defaultClient(Options{Version: "1.43.0", Now: time.Now, Logger: slog.New(slog.DiscardHandler)}, Entry{Key: "harbour", URL: "http://127.0.0.1:4242"}, []byte("s3cret"))
	if c.Release != (channel.Release{Statement: "flai 1.43.0 a1825d8", Signature: "c2lnbmF0dXJl"}) {
		t.Errorf("the client sends %+v, want the build's stamp", c.Release)
	}
	if len(c.ReleaseKeys) != 1 || !c.ReleaseKeys[0].Equal(want[0]) {
		t.Errorf("the client verifies with %v, want the embedded release key", c.ReleaseKeys)
	}
	if c.AllowUnsigned != nil {
		t.Error("the client admits an unsigned dashboard by itself")
	}
}

func TestRegistryRoundTrip(t *testing.T) {
	dir := DirFor(filepath.Join(t.TempDir(), "config.json"))
	if got, err := dir.Projects(); err != nil || len(got) != 0 {
		t.Fatalf("empty registry: %v %v", got, err)
	}
	if err := dir.Register(Entry{Root: "/p"}); err == nil {
		t.Error("an entry without a dashboard address or credential was accepted")
	}
	a := Entry{Key: "a", Name: "A", Root: "/p/a", URL: "http://127.0.0.1:1", KeyFile: "/p/a/k"}
	b := Entry{Key: "b", Name: "B", Root: "/p/b", URL: "http://127.0.0.1:2", KeyFile: "/p/b/k"}
	for _, e := range []Entry{b, a, a} {
		if err := dir.Register(e); err != nil {
			t.Fatal(err)
		}
	}
	got, _ := dir.Projects()
	if len(got) != 2 || got[0] != a || got[1] != b {
		t.Errorf("registry: %+v", got)
	}
	a.URL = "http://127.0.0.1:9"
	_ = dir.Register(a)
	_ = dir.Unregister("/p/b")
	_ = dir.Unregister("/p/never")
	got, _ = dir.Projects()
	if len(got) != 1 || got[0].URL != "http://127.0.0.1:9" {
		t.Errorf("after replace and unregister: %+v", got)
	}
	if info, err := os.Stat(filepath.Join(string(dir), "projects.json")); err != nil || info.Mode().Perm() != 0o600 {
		t.Errorf("registry mode: %v %v", info, err)
	}
}

func TestServesRegisteredProjectsAndFollowsTheRegistry(t *testing.T) {
	dash := channeltest.New(t, "s3cret")
	dir := DirFor(filepath.Join(t.TempDir(), "config.json"))
	rootA, keyA := scratchProject(t, "harbour")
	if err := dir.Register(Entry{Key: "harbour", Name: "harbour", Root: rootA, URL: dash.URL, KeyFile: keyA}); err != nil {
		t.Fatal(err)
	}
	run(t, dir, dash)
	conn := dash.Wait(t)
	if conn.Project != "harbour" {
		t.Errorf("hello named %q", conn.Project)
	}
	res, rerr := conn.Ask(t, "project.info", `{"project":"harbour"}`)
	var info hostapi.ProjectInfo
	if rerr != nil || json.Unmarshal(res, &info) != nil || info.Key != "harbour" {
		t.Errorf("project.info: %v %s", rerr, res)
	}

	// registered while it runs
	rootB, keyB := scratchProject(t, "lighthouse")
	_ = dir.Register(Entry{Key: "lighthouse", Name: "lighthouse", Root: rootB, URL: dash.URL, KeyFile: keyB})
	if second := dash.Wait(t); second.Project != "lighthouse" {
		t.Errorf("second connection is %q", second.Project)
	}

	deadline := time.Now().Add(3 * time.Second)
	var st Status
	for time.Now().Before(deadline) {
		var alive bool
		if st, alive = dir.ReadStatus(time.Now()); alive && st.Connections[rootA].Connected && st.Connections[rootB].Connected {
			break
		}
		time.Sleep(20 * time.Millisecond)
	}
	if st.PID != os.Getpid() || !st.Connections[rootA].Connected || !st.Connections[rootB].Connected {
		t.Fatalf("status: %+v", st)
	}
	data, _ := os.ReadFile(filepath.Join(string(dir), "state.json"))
	if strings.Contains(string(data), "s3cret") {
		t.Error("the status file holds the credential")
	}

	// unregistered while it runs
	_ = dir.Unregister(rootA)
	deadline = time.Now().Add(3 * time.Second)
	for time.Now().Before(deadline) {
		if st, _ = dir.ReadStatus(time.Now()); len(st.Connections) == 1 {
			break
		}
		time.Sleep(20 * time.Millisecond)
	}
	if _, still := st.Connections[rootA]; still {
		t.Errorf("an unregistered project is still served: %+v", st.Connections)
	}
}

// I-0107: a publish that raises flai.minimum past the flai serving the
// project makes its manifest stop loading while the publish is still being
// answered. The publish finishes and answers what it did, and only then is
// the project dropped, said once, with why.
func TestAProjectThatStopsLoadingAnswersItsRequestInFlightBeforeItIsDropped(t *testing.T) {
	was := buildinfo.Version
	buildinfo.Version = "1.2.0"
	t.Cleanup(func() { buildinfo.Version = was }) // the first cleanup registered runs last, once Run has ended
	dash := channeltest.New(t, "s3cret")
	dir := DirFor(filepath.Join(t.TempDir(), "config.json"))
	root, key := scratchProject(t, "harbour")
	if err := dir.Register(Entry{Key: "harbour", Name: "harbour", Root: root, URL: dash.URL, KeyFile: key}); err != nil {
		t.Fatal(err)
	}
	logs := &syncBuffer{}
	dropped := func() []map[string]any {
		var out []map[string]any
		lines := bufio.NewScanner(strings.NewReader(logs.String()))
		for lines.Scan() {
			var ev map[string]any
			if json.Unmarshal(lines.Bytes(), &ev) == nil && ev["msg"] == "project dropped" && ev["root"] == root {
				out = append(out, ev)
			}
		}
		return out
	}

	var droppedBeforeAnswer atomic.Int32
	publish := func(ctx context.Context, _ channel.Project, _ json.RawMessage) (any, *channel.Error) {
		// what flai release does: commit a minimum above the flai that serves the project
		f, err := os.OpenFile(filepath.Join(root, "system-flow.yaml"), os.O_APPEND|os.O_WRONLY, 0)
		if err != nil {
			return nil, &channel.Error{Code: channel.CodeInternal, Message: err.Error()}
		}
		_, err = f.WriteString("flai:\n  minimum: 1.2.1\n")
		if cerr := f.Close(); err == nil {
			err = cerr
		}
		if err != nil {
			return nil, &channel.Error{Code: channel.CodeInternal, Message: err.Error()}
		}
		// and then push, which takes a while: flai serve sees the manifest
		// no longer loads, and a few looks more go by
		deadline := time.Now().Add(4 * time.Second)
		for st, _ := dir.ReadStatus(time.Now()); st.Unavailable[root] == ""; st, _ = dir.ReadStatus(time.Now()) {
			if time.Now().After(deadline) {
				return nil, &channel.Error{Code: channel.CodeInternal, Message: "flai serve never found the project unavailable"}
			}
			select {
			case <-ctx.Done():
				return nil, &channel.Error{Code: channel.CodeInternal, Message: "cancelled: " + ctx.Err().Error()}
			case <-time.After(10 * time.Millisecond):
			}
		}
		select {
		case <-ctx.Done():
			return nil, &channel.Error{Code: channel.CodeInternal, Message: "cancelled: " + ctx.Err().Error()}
		case <-time.After(100 * time.Millisecond):
		}
		droppedBeforeAnswer.Store(int32(len(dropped())))
		return map[string]string{"published": "1.2.1"}, nil
	}

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() {
		done <- Run(ctx, Options{Dir: dir, Version: "1.2.0", Every: 20 * time.Millisecond, WatchEvery: time.Hour, DrainGrace: 10 * time.Second,
			Logger: slog.New(slog.NewJSONHandler(logs, nil)),
			NewClient: func(e Entry, key []byte) *channel.Client {
				methods := hostapi.Methods("test", nil)
				methods["test.publish"] = publish
				// the test dashboard answers no ping: none is sent while the test runs
				return &channel.Client{URL: e.URL, Key: key, Project: channel.Project{Key: e.Key, Name: e.Name, Root: e.Root},
					Methods: methods, Version: "test", PingEvery: time.Hour, MinBackoff: 10 * time.Millisecond, MaxBackoff: 40 * time.Millisecond, ReleaseKeys: dash.Keys}
			}})
	}()
	t.Cleanup(func() {
		cancel()
		if err := <-done; err != nil {
			t.Errorf("Run: %v", err)
		}
	})

	conn := dash.Wait(t)
	res, rerr := conn.Ask(t, "test.publish", `{"project":"harbour"}`)
	if rerr != nil || string(res) != `{"published":"1.2.1"}` {
		t.Fatalf("the publish in flight answered %s %+v, want its own result", res, rerr)
	}
	if n := droppedBeforeAnswer.Load(); n != 0 {
		t.Errorf("the project was dropped %d times before the publish in flight answered", n)
	}
	if got := conn.Last(t); len(got) != 0 {
		t.Errorf("a project that stopped loading sent %v before its connection ended, want nothing: it is not removed", got)
	}
	waitFor(t, "the project is dropped", func() bool { return len(dropped()) > 0 })
	time.Sleep(100 * time.Millisecond) // a few looks more, which say nothing again
	said := dropped()
	if len(said) != 1 {
		t.Fatalf("said %d times that the project was dropped, want once:\n%s", len(said), logs.String())
	}
	if why, _ := said[0]["reason"].(string); !strings.Contains(why, "needs flai 1.2.1") || said[0]["unanswered"] != float64(0) {
		t.Errorf("the drop said %v, want why (the manifest needs flai 1.2.1) and no request unanswered", said[0])
	}
	if st, _ := dir.ReadStatus(time.Now()); st.Connections[root].Connected {
		t.Errorf("a project that no longer loads is still served: %+v", st.Connections)
	}
}

func TestASharedCredentialDoesNotLetOneProjectAnswerForAnother(t *testing.T) {
	// One flai serve, one shared credential, two projects (S-0080): the credential proves who flai
	// is, not which project a given connection speaks for. Each connection still only answers for
	// the one project it named in hello.
	dash := channeltest.New(t, "s3cret")
	dir := DirFor(filepath.Join(t.TempDir(), "config.json"))
	rootA, keyA := scratchProject(t, "harbour")
	rootB, keyB := scratchProject(t, "lighthouse")
	if err := dir.Register(Entry{Key: "harbour", Name: "harbour", Root: rootA, URL: dash.URL, KeyFile: keyA}); err != nil {
		t.Fatal(err)
	}
	if err := dir.Register(Entry{Key: "lighthouse", Name: "lighthouse", Root: rootB, URL: dash.URL, KeyFile: keyB}); err != nil {
		t.Fatal(err)
	}
	run(t, dir, dash)
	first, second := dash.Wait(t), dash.Wait(t)
	byProject := map[string]*channeltest.Conn{first.Project: first, second.Project: second}
	harbour, lighthouse := byProject["harbour"], byProject["lighthouse"]
	if harbour == nil || lighthouse == nil {
		t.Fatalf("expected harbour and lighthouse, got %q and %q", first.Project, second.Project)
	}

	if _, rerr := harbour.Ask(t, "project.info", `{"project":"harbour"}`); rerr != nil {
		t.Errorf("harbour's connection answering for harbour: %v", rerr)
	}
	if _, rerr := harbour.Ask(t, "project.info", `{"project":"lighthouse"}`); rerr == nil || rerr.Code != channel.CodeUnknownProject {
		t.Errorf("harbour's connection asked for lighthouse: %+v", rerr)
	}
	if _, rerr := lighthouse.Ask(t, "project.info", `{"project":"lighthouse"}`); rerr != nil {
		t.Errorf("lighthouse's connection answering for lighthouse: %v", rerr)
	}
	if _, rerr := lighthouse.Ask(t, "project.info", `{"project":"harbour"}`); rerr == nil || rerr.Code != channel.CodeUnknownProject {
		t.Errorf("lighthouse's connection asked for harbour: %+v", rerr)
	}
}

func TestASecondServeRefusesWhileTheFirstRuns(t *testing.T) {
	dir := DirFor(filepath.Join(t.TempDir(), "config.json"))
	// A status as another live process would have written it: this test's parent is alive and is not us.
	st := Status{PID: os.Getppid(), Version: "test", Started: "2026-09-20T00:00:00Z", Updated: time.Now().UTC().Format(time.RFC3339)}
	if err := dir.write(dir.status(), st); err != nil {
		t.Fatal(err)
	}
	err := Run(context.Background(), Options{Dir: dir, Version: "test"})
	if err == nil || !strings.Contains(err.Error(), "already running") {
		t.Errorf("second serve: %v", err)
	}
	// A stale status is a dead serve: Run starts.
	st.Updated = time.Now().Add(-time.Minute).UTC().Format(time.RFC3339)
	_ = dir.write(dir.status(), st)
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()
	if err := Run(ctx, Options{Dir: dir, Version: "test", Every: 10 * time.Millisecond}); err != nil {
		t.Errorf("after a stale status: %v", err)
	}
	if _, err := os.Stat(dir.status()); err == nil {
		t.Error("a serve that ended left its status behind")
	}
}

func TestAChangedFileReachesTheDashboardAsANotification(t *testing.T) {
	dash := channeltest.New(t, "s3cret")
	dir := DirFor(filepath.Join(t.TempDir(), "config.json"))
	root, key := scratchProject(t, "harbour")
	_ = os.MkdirAll(filepath.Join(root, "wip/kanban/stories"), 0o755)
	_ = os.MkdirAll(filepath.Join(root, "src"), 0o755)
	if err := dir.Register(Entry{Key: "harbour", Name: "harbour", Root: root, URL: dash.URL, KeyFile: key}); err != nil {
		t.Fatal(err)
	}
	run(t, dir, dash)
	conn := dash.Wait(t)
	time.Sleep(80 * time.Millisecond)                                                     // the watcher has its first look
	_ = os.WriteFile(filepath.Join(root, "src/main.go"), []byte("package main\n"), 0o644) // not the dashboard's business
	_ = os.WriteFile(filepath.Join(root, "wip/kanban/stories/S-0001-a.md"), []byte("a\n"), 0o644)
	method, params := conn.Next(t)
	var p struct{ Project, Path string }
	_ = json.Unmarshal(params, &p)
	if method != "change" || p.Project != "harbour" || p.Path != "wip/kanban/stories/S-0001-a.md" {
		t.Errorf("notification: %s %s", method, params)
	}
}

// S-0154: an agent starting or ending changes no file of the project, so flai
// serve tells the dashboard with the notification agent, naming the story.
func TestServeTellsTheDashboardWhenAStorysAgentStartsAndEnds(t *testing.T) {
	dash := channeltest.New(t, "s3cret")
	// The launcher records the ended agent and looks again after Run returns:
	// its folder is removed once that is done, not while it writes.
	host, err := os.MkdirTemp("", "serve-agent-changed")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		for range 50 {
			if os.RemoveAll(host) == nil {
				return
			}
			time.Sleep(20 * time.Millisecond)
		}
		t.Errorf("%s not removed", host)
	})
	dir := DirFor(filepath.Join(host, "config.json"))
	root, key := scratchProject(t, "harbour")
	if err := dir.Register(Entry{Key: "harbour", Name: "harbour", Root: root, URL: dash.URL, KeyFile: key}); err != nil {
		t.Fatal(err)
	}
	for _, d := range []string{"wip/kanban/epics", "wip/kanban/stories", "wip/kanban/tasks", "wip/agents"} {
		_ = os.MkdirAll(filepath.Join(root, d), 0o755)
	}
	stub := filepath.Join(t.TempDir(), "stub-agent")
	if err := os.WriteFile(stub, []byte("#!/bin/sh\nexit 0\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() {
		done <- Run(ctx, Options{Dir: dir, Version: "test", Every: 20 * time.Millisecond, WatchEvery: 20 * time.Millisecond,
			Agent: func(string) AgentConfig { return AgentConfig{Enabled: true, Command: []string{stub}, Name: "builder"} },
			NewClient: func(e Entry, key []byte) *channel.Client {
				return &channel.Client{URL: e.URL, Key: key, Project: channel.Project{Key: e.Key, Name: e.Name, Root: e.Root},
					Methods: hostapi.Methods("test", nil), Version: "test", PingEvery: time.Hour /* the test dashboard answers no ping: none is sent while the test runs */, MinBackoff: 10 * time.Millisecond, MaxBackoff: 40 * time.Millisecond, ReleaseKeys: dash.Keys}
			}})
	}()
	t.Cleanup(func() {
		cancel()
		if err := <-done; err != nil {
			t.Errorf("Run: %v", err)
		}
	})
	conn := dash.Wait(t)

	// made ready once the dashboard listens, so that the start is heard
	repo, err := workitem.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	now := time.Now()
	epic, err := repo.Create(workitem.NewOptions{Type: workitem.Epic, Title: "Epic", Owner: "alex", Now: now})
	if err != nil {
		t.Fatal(err)
	}
	st, err := repo.Create(workitem.NewOptions{Type: workitem.Story, Title: "Story", Parent: epic.ID, Owner: "alex", Touches: []string{"docs/story"}, Now: now})
	if err != nil {
		t.Fatal(err)
	}
	data, _ := os.ReadFile(st.Path)
	_ = os.WriteFile(st.Path, []byte(strings.Replace(string(data), "## Acceptance criteria\n", "## Acceptance criteria\n- [ ] works\n", 1)), 0o644)
	if st, err = repo.Get(st.ID); err != nil {
		t.Fatal(err)
	}
	if _, err := repo.Transition(st, workitem.Ready, "alex", "", now); err != nil {
		t.Fatal(err)
	}

	told := 0
	for told < 2 {
		method, params := conn.Next(t)
		if method != AgentChanged {
			continue
		}
		var p map[string]string
		_ = json.Unmarshal(params, &p)
		if p["project"] != "harbour" || p["story"] != st.ID {
			t.Fatalf("agent said %v, want the project harbour and the story %s", p, st.ID)
		}
		told++
	}
	waitFor(t, "the agent's end is recorded", func() bool {
		run := dir.AgentStates()[root].Stories[st.ID]
		return run != nil && !run.live()
	})
}

// lockedDocker is a fake docker that a test may swap while flai serve asks
// it, as a container restarted from another image would change its answers.
type lockedDocker struct {
	mu sync.Mutex
	d  *fakeDocker
}

func (l *lockedDocker) set(d *fakeDocker) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.d = d
}

func (l *lockedDocker) Run(dir, name string, args ...string) (string, error) {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.d.Run(dir, name, args...)
}

func (l *lockedDocker) RunInput(dir, name, input string, args ...string) (string, error) {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.d.RunInput(dir, name, input, args...)
}

func (l *lockedDocker) LookPath(name string) (string, error) {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.d.LookPath(name)
}

// imageLab is flai serve judging the dashboard's image with docker and the
// verified lists of flaiover imgVersion, cached and signed by a key it trusts.
type imageLab struct {
	dir    Dir
	root   string
	dash   *channeltest.Dashboard
	docker *lockedDocker
	logs   *syncBuffer
}

// labSettings are the allowance an image lab runs with, and which sides of
// its connection carry a release stamp (S-0239).
type labSettings struct {
	host      func() bool // the host configuration's dashboard.allow_unsigned; nil is off
	manifest  string      // the project manifest's dashboard.allow_unsigned, "" for none
	signed    bool        // this flai sends a release stamp
	unstamped bool        // the dashboard answers hello without one
}

func newImageLab(t *testing.T, docker *fakeDocker) *imageLab {
	t.Helper()
	return newImageLabWith(t, docker, labSettings{})
}

func newImageLabWith(t *testing.T, docker *fakeDocker, set labSettings) *imageLab {
	t.Helper()
	key := imgKey(t)
	cache := imgCache(t, key)
	cachePut(t, cache, key)
	lab := &imageLab{dir: DirFor(filepath.Join(t.TempDir(), "config.json")), dash: channeltest.New(t, "s3cret"),
		docker: &lockedDocker{d: docker}, logs: &syncBuffer{}}
	root, keyFile := scratchProject(t, "harbour")
	lab.root = root
	if set.manifest != "" {
		f, err := os.OpenFile(filepath.Join(root, "system-flow.yaml"), os.O_APPEND|os.O_WRONLY, 0)
		if err != nil {
			t.Fatal(err)
		}
		_, err = f.WriteString("dashboard:\n  allow_unsigned: " + set.manifest + "\n")
		if cerr := f.Close(); err == nil {
			err = cerr
		}
		if err != nil {
			t.Fatal(err)
		}
	}
	if set.unstamped {
		lab.dash.Answer(nil)
	}
	var own channel.Release
	if set.signed {
		own = channel.Release{Statement: "flai 1.43.0 a1825d8", Signature: "c2lnbmF0dXJl"} // the test dashboard does not check it
	}
	if err := lab.dir.Register(Entry{Key: "harbour", Name: "harbour", Root: root, URL: lab.dash.URL, KeyFile: keyFile}); err != nil {
		t.Fatal(err)
	}
	logger := slog.New(slog.NewJSONHandler(lab.logs, nil))
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() {
		done <- Run(ctx, Options{Dir: lab.dir, Version: "test", Every: 20 * time.Millisecond, WatchEvery: time.Hour, Logger: logger,
			Images: &ImageLists{Cache: cache, Fetch: noFetch(t)}, Docker: lab.docker, AllowUnsigned: set.host,
			NewClient: func(e Entry, key []byte) *channel.Client {
				return &channel.Client{URL: e.URL, Key: key, Project: channel.Project{Key: e.Key, Name: e.Name, Root: e.Root}, Logger: logger,
					Methods: hostapi.Methods("test", nil), Version: "test", PingEvery: time.Hour /* the test dashboard answers no ping */, MinBackoff: 10 * time.Millisecond, MaxBackoff: 40 * time.Millisecond,
					HeldBackoff: 20 * time.Millisecond, ReleaseKeys: lab.dash.Keys, Release: own}
			}})
	}()
	t.Cleanup(func() {
		cancel()
		if err := <-done; err != nil {
			t.Errorf("Run: %v", err)
		}
	})
	return lab
}

// status waits until flai serve's status satisfies ok, and returns it.
func (lab *imageLab) status(t *testing.T, what string, ok func(Status) bool) Status {
	t.Helper()
	var st Status
	waitFor(t, what, func() bool {
		st, _ = lab.dir.ReadStatus(time.Now())
		return ok(st)
	})
	return st
}

// refusedAttempts waits until the client has been refused at least n times
// with the image unlisted in the status, and returns the status.
func (lab *imageLab) refusedAttempts(t *testing.T, n int) Status {
	t.Helper()
	return lab.status(t, "the unlisted image is refused, again at each retry", func(st Status) bool {
		return st.Images[lab.root].Kind == ImageUnlisted && st.Connections[lab.root].Attempts >= n
	})
}

// S-0238: a dashboard whose container runs an image a verified digest list
// names is dialled, and the status says which release lists it.
func TestServeDialsADashboardRunningAListedImage(t *testing.T) {
	lab := newImageLab(t, dockerWith(t, imgVersion, []string{imgRepo + "@" + imgAMD64}))

	if conn := lab.dash.Wait(t); conn.Project != "harbour" {
		t.Fatalf("hello named %q", conn.Project)
	}
	st := lab.status(t, "the verdict is in the status", func(st Status) bool { return st.Images[lab.root].Kind != "" })
	if v := st.Images[lab.root]; v.Kind != ImageListed || v.Digest != imgAMD64 || v.Release != imgVersion || v.Image != imgID {
		t.Fatalf("verdict %+v, want listed at %s by flaiover %s", v, imgAMD64, imgVersion)
	}
}

// S-0238: an image no verified list names, or a local build with no
// RepoDigests, is not dialled: the refusal names the image and its digests,
// is logged once at error however often it is met again, and the verdict is
// in the status.
func TestServeRefusesADashboardRunningAnUnlistedImage(t *testing.T) {
	cases := []struct {
		name    string
		docker  *fakeDocker
		digests string
	}{
		{"unlisted digest", dockerWith(t, imgVersion, []string{imgRepo + "@" + imgOther}), imgRepo + "@" + imgOther},
		{"no RepoDigests", dockerWith(t, "0.0.0", nil), "no RepoDigests"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			lab := newImageLab(t, tc.docker)

			st := lab.refusedAttempts(t, 3)
			select {
			case conn := <-lab.dash.Proven:
				t.Fatalf("flai dialled a dashboard running an unlisted image, for %q", conn.Project)
			default:
			}
			want := "the dashboard runs an unsigned image " + imgID + " (" + tc.digests + "): "
			c := st.Connections[lab.root]
			if c.Connected || c.RefusedBy != channel.RefusedByFlai || !strings.HasPrefix(c.RefusedReason, want) {
				t.Fatalf("connection %+v, want refused by flai with a reason beginning %q", c, want)
			}
			if v := st.Images[lab.root]; v.Image != imgID || v.Container != imgContainer || v.Reason == "" {
				t.Fatalf("verdict %+v, want the image %s of the container %s and why", v, imgID, imgContainer)
			}
			var errs int
			lines := bufio.NewScanner(strings.NewReader(lab.logs.String()))
			for lines.Scan() {
				var ev map[string]any
				if json.Unmarshal(lines.Bytes(), &ev) == nil && ev["level"] == "ERROR" {
					errs++
					if why, _ := ev["reason"].(string); !strings.HasPrefix(why, want) {
						t.Errorf("the error event says %v, want the reason beginning %q", ev, want)
					}
				}
			}
			if errs != 1 {
				t.Fatalf("%d error events over %d refused attempts, want 1:\n%s", errs, c.Attempts, lab.logs.String())
			}
		})
	}
}

// S-0238: a dashboard Docker here cannot inspect is dialled, the release
// stamp check alone admitting it, and the status says it was not measured.
func TestServeDialsADashboardItCannotInspectAndSaysItWasNotMeasured(t *testing.T) {
	lab := newImageLab(t, &fakeDocker{absent: true})

	lab.dash.Wait(t)
	st := lab.status(t, "the verdict is in the status", func(st Status) bool { return st.Images[lab.root].Kind != "" })
	if v := st.Images[lab.root]; v.Kind != ImageNotMeasured || !strings.Contains(v.Reason, "not on PATH") {
		t.Fatalf("verdict %+v, want not measured because docker is not on PATH", v)
	}
	data, err := os.ReadFile(filepath.Join(string(lab.dir), "state.json"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(data), `"kind": "not-measured"`) {
		t.Errorf("the status file does not say not-measured under kind:\n%s", data)
	}
}

// S-0238: a container restarted from a listed image is dialled at the next
// look after the refusal of the unlisted one it ran before.
func TestServeDialsOnceTheContainerRunsAListedImage(t *testing.T) {
	lab := newImageLab(t, dockerWith(t, imgVersion, []string{imgRepo + "@" + imgOther}))
	lab.refusedAttempts(t, 1)

	lab.docker.set(dockerWith(t, imgVersion, []string{imgRepo + "@" + imgARM64}))

	lab.dash.Wait(t)
	st := lab.status(t, "the listed verdict and the connection are in the status", func(st Status) bool {
		return st.Images[lab.root].Kind == ImageListed && st.Connections[lab.root].Connected
	})
	if c := st.Connections[lab.root]; c.RefusedBy != "" {
		t.Errorf("connection %+v still says it was refused", c)
	}
}

func on() bool { return true }

// unsignedWarns are the warn events that name a connection's unsigned sides.
func (lab *imageLab) unsignedWarns() []map[string]any {
	var out []map[string]any
	lines := bufio.NewScanner(strings.NewReader(lab.logs.String()))
	for lines.Scan() {
		var ev map[string]any
		if json.Unmarshal(lines.Bytes(), &ev) == nil && ev["level"] == "WARN" && ev["unsigned"] != nil {
			out = append(out, ev)
		}
	}
	return out
}

// S-0239: with the allowance in force for the project, by the host's setting
// or by the manifest's over the host's, a local build with no RepoDigests is
// dialled and still recorded as unlisted; a dashboard without a stamp is
// accepted; the unsigned sides reach the status; and one warn names them for
// each connection, again for a new one.
func TestServeDialsAnUnlistedImageWithTheAllowance(t *testing.T) {
	for _, tc := range []struct {
		name string
		set  labSettings
	}{
		{"the host's setting", labSettings{host: on, unstamped: true}},
		{"the manifest's over the host's", labSettings{manifest: "true", unstamped: true}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			lab := newImageLabWith(t, dockerWith(t, "0.0.0", nil), tc.set)
			conn := lab.dash.Wait(t)
			want := []string{channel.UnsignedFlai, channel.UnsignedDashboard}
			st := lab.status(t, "the connection and its unsigned sides are in the status", func(st Status) bool {
				return st.Connections[lab.root].Connected && len(st.Connections[lab.root].Unsigned) > 0
			})
			if c := st.Connections[lab.root]; !slices.Equal(c.Unsigned, want) || c.RefusedBy != "" || c.PeerSigned {
				t.Errorf("connection %+v, want unsigned %v", c, want)
			}
			if v := st.Images[lab.root]; v.Kind != ImageUnlisted {
				t.Errorf("verdict %+v, want unlisted though dialled", v)
			}
			data, err := os.ReadFile(filepath.Join(string(lab.dir), "state.json"))
			if err != nil {
				t.Fatal(err)
			}
			if !strings.Contains(string(data), `"unsigned": [`) {
				t.Errorf("the status file does not carry the unsigned sides:\n%s", data)
			}
			said := lab.unsignedWarns()
			if len(said) != 1 || said[0]["project"] != "harbour" || said[0]["component"] == nil || fmt.Sprint(said[0]["unsigned"]) != fmt.Sprint(want) {
				t.Fatalf("logged %v, want one warn naming %v", said, want)
			}

			conn.Close()
			lab.dash.Wait(t)
			waitFor(t, "a second warn for the second connection", func() bool { return len(lab.unsignedWarns()) == 2 })
		})
	}
}

// S-0239: without the allowance in force, the host's off and none in the
// manifest, or the manifest's off over the host's on, an unlisted image is
// refused as before.
func TestServeRefusesAnUnlistedImageWithoutTheAllowance(t *testing.T) {
	for _, tc := range []struct {
		name string
		set  labSettings
	}{
		{"none", labSettings{unstamped: true}},
		{"the manifest's off over the host's on", labSettings{host: on, manifest: "false", unstamped: true}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			lab := newImageLabWith(t, dockerWith(t, "0.0.0", nil), tc.set)
			st := lab.refusedAttempts(t, 3)
			select {
			case conn := <-lab.dash.Proven:
				t.Fatalf("flai dialled an unlisted image without the allowance, for %q", conn.Project)
			default:
			}
			if c := st.Connections[lab.root]; c.Connected || c.RefusedBy != channel.RefusedByFlai || c.Unsigned != nil {
				t.Errorf("connection %+v, want refused by flai", c)
			}
			if said := lab.unsignedWarns(); len(said) != 0 {
				t.Errorf("logged %v", said)
			}
		})
	}
}

// S-0239: the host's setting is read at each dial: turned on, the next dial
// of a refused unlisted image goes ahead.
func TestTheAllowanceTakesEffectAtTheNextDial(t *testing.T) {
	var allow atomic.Bool
	lab := newImageLabWith(t, dockerWith(t, "0.0.0", nil), labSettings{host: allow.Load, unstamped: true})
	lab.refusedAttempts(t, 1)
	allow.Store(true)
	lab.dash.Wait(t)
}

// S-0239: a signed pair with the allowance on records and logs nothing of it:
// the allowance permits, it does not downgrade.
func TestASignedPairWithTheAllowanceShowsNothing(t *testing.T) {
	lab := newImageLabWith(t, dockerWith(t, imgVersion, []string{imgRepo + "@" + imgAMD64}), labSettings{host: on, signed: true})
	lab.dash.Wait(t)
	st := lab.status(t, "connected", func(st Status) bool { return st.Connections[lab.root].Connected })
	if c := st.Connections[lab.root]; c.Unsigned != nil || !c.PeerSigned {
		t.Errorf("connection %+v, want signed with no unsigned side", c)
	}
	if v := st.Images[lab.root]; v.Kind != ImageListed {
		t.Errorf("verdict %+v, want listed", v)
	}
	if said := lab.unsignedWarns(); len(said) != 0 {
		t.Errorf("a signed pair logged %v", said)
	}
}
