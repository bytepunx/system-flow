package serve

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/channel"
	"github.com/bytepunx/system-flow/flai/internal/channel/channeltest"
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

func run(t *testing.T, dir Dir) {
	t.Helper()
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() {
		done <- Run(ctx, Options{Dir: dir, Version: "test", Every: 20 * time.Millisecond,
			NewClient: func(e Entry, key []byte) *channel.Client {
				return &channel.Client{URL: e.URL, Key: key, Project: channel.Project{Key: e.Key, Name: e.Name, Root: e.Root},
					Methods: channel.Methods("test"), Version: "test", PingEvery: 50 * time.Millisecond, MinBackoff: 10 * time.Millisecond, MaxBackoff: 40 * time.Millisecond}
			}})
	}()
	t.Cleanup(func() {
		cancel()
		if err := <-done; err != nil {
			t.Errorf("Run: %v", err)
		}
	})
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
	run(t, dir)
	conn := dash.Wait(t)
	if conn.Project != "harbour" {
		t.Errorf("hello named %q", conn.Project)
	}
	res, rerr := conn.Ask(t, "project.info", `{"project":"harbour"}`)
	var info channel.ProjectInfo
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
