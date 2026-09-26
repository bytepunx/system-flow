package serve

import (
	"context"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/channel"
	"github.com/bytepunx/system-flow/flai/internal/channel/channeltest"
	"github.com/bytepunx/system-flow/flai/internal/hostapi"
)

func TestRegisterRefusesAKeyAnotherRootHas(t *testing.T) {
	dir := DirFor(filepath.Join(t.TempDir(), "config.json"))
	a := Entry{Key: "harbour", Name: "harbour", Root: "/p/a", URL: "http://127.0.0.1:1", KeyFile: "/k"}
	if err := dir.Register(a); err != nil {
		t.Fatal(err)
	}
	b := a
	b.Root = "/p/b"
	err := dir.Register(b)
	if err == nil || !strings.Contains(err.Error(), "/p/a") || !strings.Contains(err.Error(), "harbour") {
		t.Errorf("a second root with the key: %v", err)
	}
	// the same root again, at another address, is a replacement and not a conflict
	a.URL = "http://127.0.0.1:2"
	if err := dir.Register(a); err != nil {
		t.Errorf("replacing the entry: %v", err)
	}
	if got, _ := dir.Projects(); len(got) != 1 || got[0].URL != a.URL {
		t.Errorf("registry: %+v", got)
	}
}

func TestUnavailableSaysWhyARegisteredProjectCannotBeServed(t *testing.T) {
	root, key := scratchProject(t, "harbour")
	e := Entry{Key: "harbour", Root: root, KeyFile: key}
	if why := e.Unavailable(); why != "" {
		t.Errorf("a whole project: %q", why)
	}
	for _, c := range []struct {
		name  string
		entry func() Entry
		want  string
	}{
		{"gone", func() Entry { x := e; x.Root = filepath.Join(root, "gone"); return x }, "the folder is gone"},
		{"no manifest", func() Entry { x := e; x.Root = t.TempDir(); return x }, "no system-flow.yaml"},
		{"no credential", func() Entry { x := e; x.KeyFile = filepath.Join(root, "none"); return x }, "credential"},
	} {
		if why := c.entry().Unavailable(); !strings.Contains(why, c.want) {
			t.Errorf("%s: %q, want %q", c.name, why, c.want)
		}
	}
	noKey := t.TempDir()
	_ = os.WriteFile(filepath.Join(noKey, "system-flow.yaml"), []byte("version: 1\nname: n\nlayout:\n  design: design\n  docs: docs\n  wip: wip\n"), 0o644)
	if why := (Entry{Root: noKey, KeyFile: key}).Unavailable(); why != "system-flow.yaml has no key" {
		t.Errorf("no key: %q", why)
	}
}

func TestAnUnavailableProjectIsReportedOnceAndServedWhenItComesBack(t *testing.T) {
	dash := channeltest.New(t, "s3cret")
	dir := DirFor(filepath.Join(t.TempDir(), "config.json"))
	good, key := scratchProject(t, "harbour")
	lost, _ := scratchProject(t, "lighthouse")
	gone := filepath.Join(t.TempDir(), "moved-away")
	for _, e := range []Entry{
		{Key: "harbour", Name: "harbour", Root: good, URL: dash.URL, KeyFile: key},
		{Key: "lighthouse", Name: "lighthouse", Root: lost, URL: dash.URL, KeyFile: key},
		{Key: "moved", Name: "moved", Root: gone, URL: dash.URL, KeyFile: key},
	} {
		if err := dir.Register(e); err != nil {
			t.Fatal(err)
		}
	}
	manifest := filepath.Join(lost, "system-flow.yaml")
	saved, _ := os.ReadFile(manifest)
	_ = os.Remove(manifest)

	logs := &lockedBuffer{}
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() {
		done <- Run(ctx, Options{Dir: dir, Version: "test", Every: 20 * time.Millisecond, WatchEvery: 20 * time.Millisecond,
			Logger: slog.New(slog.NewTextHandler(logs, nil)),
			NewClient: func(e Entry, key []byte) *channel.Client {
				return &channel.Client{URL: e.URL, Key: key, Project: channel.Project{Key: e.Key, Name: e.Name, Root: e.Root},
					Methods: hostapi.Methods("test", nil), Version: "test", PingEvery: 50 * time.Millisecond, MinBackoff: 10 * time.Millisecond, MaxBackoff: 40 * time.Millisecond}
			}})
	}()
	t.Cleanup(func() {
		cancel()
		if err := <-done; err != nil {
			t.Errorf("Run: %v", err)
		}
	})
	if conn := dash.Wait(t); conn.Project != "harbour" {
		t.Errorf("served %q", conn.Project)
	}
	time.Sleep(200 * time.Millisecond) // ten ticks
	st, _ := dir.ReadStatus(time.Now())
	if _, served := st.Connections[lost]; served || st.Unavailable[lost] != "it has no system-flow.yaml" || st.Unavailable[gone] != "the folder is gone" {
		t.Errorf("status: connections %v, unavailable %v", st.Connections, st.Unavailable)
	}
	count := func(msg string) int {
		logs.mu.Lock()
		defer logs.mu.Unlock()
		return strings.Count(logs.buf.String(), `msg="`+msg+`"`)
	}
	if n := count("registered project not served"); n != 2 {
		t.Errorf("said %d times that a registered project is not served, want once for each of two:\n%s", n, logs.buf.String())
	}

	_ = os.WriteFile(manifest, saved, 0o644)
	// the test dashboard answers no ping, so harbour comes and goes meanwhile
	for conn := dash.Wait(t); conn.Project != "lighthouse"; conn = dash.Wait(t) {
	}
	if n := count("registered project can be served again"); n != 1 {
		t.Errorf("said %d times that it came back", n)
	}
}
