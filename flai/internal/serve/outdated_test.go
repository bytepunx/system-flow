package serve

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/channel"
	"github.com/bytepunx/system-flow/flai/internal/channel/channeltest"
	"github.com/bytepunx/system-flow/flai/internal/hostapi"
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

// S-0181: flai serve says when the flai serving a project is older than the
// newest flai release in the project's history, in its log when it starts
// serving it and in project.info for the dashboard's host badge.
func TestServeSaysTheFlaiIsOlderThanTheProject(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not installed")
	}
	root, keyFile := scratchProject(t, "harbour")
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
					Methods: hostapi.Methods("1.26.4", nil), Version: "1.26.4", PingEvery: 50 * time.Millisecond, MinBackoff: 10 * time.Millisecond, MaxBackoff: 40 * time.Millisecond}
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
