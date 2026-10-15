package serve

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/channel"
	"github.com/bytepunx/system-flow/flai/internal/channel/channeltest"
	"github.com/bytepunx/system-flow/flai/internal/hostapi"
	"github.com/bytepunx/system-flow/flai/internal/manifest"
	"github.com/bytepunx/system-flow/flai/internal/threads"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

type syncBuffer struct {
	mu sync.Mutex
	b  bytes.Buffer
}

func (s *syncBuffer) Write(p []byte) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.b.Write(p)
}

func (s *syncBuffer) String() string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.b.String()
}

// taggedProject is a scratch project whose history holds flai/v1.27.0.
func taggedProject(t *testing.T, key string) (root, keyFile string) {
	t.Helper()
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not installed")
	}
	root, keyFile = scratchProject(t, key)
	for _, args := range [][]string{
		{"init", "-q", "-b", "main"}, {"add", "system-flow.yaml"},
		{"-c", "user.name=t", "-c", "user.email=t@t", "commit", "-q", "-m", "init"},
		{"-c", "user.name=t", "-c", "user.email=t@t", "tag", "-a", "flai/v1.27.0", "-m", "flai 1.27.0"},
	} {
		cmd := exec.Command("git", args...)
		cmd.Dir = root
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("git %v: %v %s", args, err, out)
		}
	}
	return root, keyFile
}

// serveOlder runs flai serve 1.26.4 on root, with opt's host and upgrade,
// until the test ends.
func serveOlder(t *testing.T, root, keyFile string, opt Options) *syncBuffer {
	t.Helper()
	dash := channeltest.New(t, "s3cret")
	dir := DirFor(filepath.Join(t.TempDir(), "config.json"))
	if err := dir.Register(Entry{Key: "harbour", Name: "harbour", Root: root, URL: dash.URL, KeyFile: keyFile}); err != nil {
		t.Fatal(err)
	}
	logs := &syncBuffer{}
	opt.Dir, opt.Version, opt.Every, opt.WatchEvery = dir, "1.26.4", 20*time.Millisecond, 20*time.Millisecond
	opt.Logger = slog.New(slog.NewTextHandler(logs, nil))
	opt.NewClient = func(e Entry, key []byte) *channel.Client {
		return &channel.Client{URL: e.URL, Key: key, Project: channel.Project{Key: e.Key, Name: e.Name, Root: e.Root},
			Methods: hostapi.Methods("1.26.4", nil), Version: "1.26.4", PingEvery: time.Hour /* the test dashboard answers no ping: none is sent while the test runs */, MinBackoff: 10 * time.Millisecond, MaxBackoff: 40 * time.Millisecond, ReleaseKeys: dash.Keys}
	}
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() { done <- Run(ctx, opt) }()
	t.Cleanup(func() {
		cancel()
		if err := <-done; err != nil {
			t.Errorf("Run: %v", err)
		}
	})
	dash.Wait(t)
	return logs
}

// S-0387, ADR-0158: flai serve, older than the newest flai release in the
// history of a project it serves, asks the host it runs under for that
// release when the project's auto-upgrade action is on.
func TestServeAsksItsHostToUpgradeWithAutoUpgradeOn(t *testing.T) {
	root, keyFile := taggedProject(t, "harbour")
	var mu sync.Mutex
	var asked []string
	logs := serveOlder(t, root, keyFile, Options{
		Host: hostapi.Host{Enabled: func(action, r string) bool { return action == hostapi.ActionAutoUpgrade && r == root }},
		Upgrade: func(_ context.Context, version string) (bool, error) {
			mu.Lock()
			defer mu.Unlock()
			asked = append(asked, version)
			return true, nil
		}})
	deadline := time.Now().Add(3 * time.Second)
	for time.Now().Before(deadline) && !strings.Contains(logs.String(), "flai upgrade asked") {
		time.Sleep(20 * time.Millisecond)
	}
	time.Sleep(100 * time.Millisecond) // a few more looks, which must not ask again within the minute
	mu.Lock()
	defer mu.Unlock()
	if strings.Join(asked, " ") != "1.27.0" {
		t.Errorf("the host was asked for %v, want 1.27.0 once; serve log:\n%s", asked, logs.String())
	}
}

// S-0387, ADR-0158: with auto-upgrade off it asks the operator on one thread
// on the project's system-flow.yaml instead, as flai serve.
func TestServeAsksTheOperatorWithAutoUpgradeOff(t *testing.T) {
	root, keyFile := taggedProject(t, "harbour")
	logs := serveOlder(t, root, keyFile, Options{Upgrade: func(context.Context, string) (bool, error) {
		t.Error("the host was asked with auto-upgrade off")
		return false, nil
	}})
	deadline := time.Now().Add(3 * time.Second)
	for time.Now().Before(deadline) && !strings.Contains(logs.String(), "flai upgrade asked on a thread") {
		time.Sleep(20 * time.Millisecond)
	}
	time.Sleep(100 * time.Millisecond)
	repo, err := workitem.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	all, err := threads.List(repo)
	if err != nil {
		t.Fatal(err)
	}
	if len(all) != 1 || all[0].Anchor.Path != manifest.File || !strings.Contains(all[0].Title, "1.27.0") || !strings.Contains(all[0].Body, "flai host upgrade --version 1.27.0") || !slices.Contains(all[0].Participants, "flai serve") {
		t.Errorf("one thread on the manifest from flai serve: %+v; serve log:\n%s", all, logs.String())
	}
}

// S-0181: flai serve says when the flai serving a project is older than the
// newest flai release in the project's history, in its log when it starts
// serving it and in project.info for the dashboard's host badge.
func TestServeSaysTheFlaiIsOlderThanTheProject(t *testing.T) {
	root, keyFile := taggedProject(t, "harbour")
	dash := channeltest.New(t, "s3cret")
	dir := DirFor(filepath.Join(t.TempDir(), "config.json"))
	if err := dir.Register(Entry{Key: "harbour", Name: "harbour", Root: root, URL: dash.URL, KeyFile: keyFile}); err != nil {
		t.Fatal(err)
	}
	var logs syncBuffer
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() {
		done <- Run(ctx, Options{Dir: dir, Version: "1.26.4", Every: 20 * time.Millisecond, WatchEvery: 20 * time.Millisecond,
			Logger: slog.New(slog.NewTextHandler(&logs, nil)),
			NewClient: func(e Entry, key []byte) *channel.Client {
				return &channel.Client{URL: e.URL, Key: key, Project: channel.Project{Key: e.Key, Name: e.Name, Root: e.Root},
					Methods: hostapi.Methods("1.26.4", nil), Version: "1.26.4", PingEvery: time.Hour /* the test dashboard answers no ping: none is sent while the test runs */, MinBackoff: 10 * time.Millisecond, MaxBackoff: 40 * time.Millisecond, ReleaseKeys: dash.Keys}
			}})
	}()
	t.Cleanup(func() {
		cancel()
		if err := <-done; err != nil {
			t.Errorf("Run: %v", err)
		}
	})
	conn := dash.Wait(t)
	res, rerr := conn.Ask(t, "project.info", `{"project":"harbour"}`)
	var info hostapi.ProjectInfo
	if rerr != nil || json.Unmarshal(res, &info) != nil || info.FlaiOutdated == nil || info.FlaiOutdated.Newest != "1.27.0" || info.FlaiOutdated.Running != "1.26.4" {
		t.Errorf("project.info: %v %s", rerr, res)
	}
	deadline := time.Now().Add(3 * time.Second)
	for time.Now().Before(deadline) && !strings.Contains(logs.String(), "flai is older than the project it serves") {
		time.Sleep(20 * time.Millisecond)
	}
	if got := logs.String(); !strings.Contains(got, "level=WARN msg=\"flai is older than the project it serves\"") || !strings.Contains(got, "newest=1.27.0") || !strings.Contains(got, "flai host upgrade") {
		t.Errorf("serve log:\n%s", got)
	}
}
