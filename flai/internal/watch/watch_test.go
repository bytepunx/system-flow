package watch

import (
	"context"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"
)

type collector struct {
	mu   sync.Mutex
	seen []string
}

func (c *collector) emit(rel string) {
	c.mu.Lock()
	c.seen = append(c.seen, rel)
	c.mu.Unlock()
}

func (c *collector) wait(t *testing.T, want ...string) {
	t.Helper()
	deadline := time.Now().Add(3 * time.Second)
	for time.Now().Before(deadline) {
		c.mu.Lock()
		got := append([]string{}, c.seen...)
		c.mu.Unlock()
		ok := len(got) == len(want)
		for i := range want {
			ok = ok && i < len(got) && got[i] == want[i]
		}
		if ok {
			return
		}
		time.Sleep(10 * time.Millisecond)
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	t.Fatalf("seen %v, want %v", c.seen, want)
}

func (c *collector) reset() {
	c.mu.Lock()
	c.seen = nil
	c.mu.Unlock()
}

func write(t *testing.T, root, rel, content string) {
	t.Helper()
	full := filepath.Join(root, rel)
	if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(full, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func TestReportsChangedAddedAndRemovedFilesUnderTheWatchedPathsOnly(t *testing.T) {
	root := t.TempDir()
	write(t, root, "system-flow.yaml", "name: x\n")
	write(t, root, "wip/kanban/stories/S-0001-a.md", "a\n")
	write(t, root, "src/main.go", "package main\n")
	w := &Watcher{Root: root, Paths: []string{"design", "docs", "wip", "system-flow.yaml"}, Every: 20 * time.Millisecond}
	c := &collector{}
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})
	go func() { w.Run(ctx, c.emit); close(done) }()
	t.Cleanup(func() { cancel(); <-done })
	time.Sleep(60 * time.Millisecond)
	c.wait(t) // nothing is reported for what was already there

	write(t, root, "wip/kanban/stories/S-0001-a.md", "a changed\n")
	c.wait(t, "wip/kanban/stories/S-0001-a.md")
	c.reset()

	// a new file in a folder that did not exist, and a folder that was not there at the start
	write(t, root, "design/system/new/plan.md", "# plan\n")
	c.wait(t, "design/system/new/plan.md")
	c.reset()

	write(t, root, "system-flow.yaml", "name: y, longer\n")
	c.wait(t, "system-flow.yaml")
	c.reset()

	if err := os.Remove(filepath.Join(root, "wip/kanban/stories/S-0001-a.md")); err != nil {
		t.Fatal(err)
	}
	c.wait(t, "wip/kanban/stories/S-0001-a.md")
	c.reset()

	// outside the watched paths, and inside a dot folder: nothing
	write(t, root, "src/main.go", "package main // changed\n")
	write(t, root, "wip/.cache/x.md", "x\n")
	time.Sleep(120 * time.Millisecond)
	c.wait(t)
}

func TestAFileStillBeingWrittenIsReportedOnce(t *testing.T) {
	root := t.TempDir()
	write(t, root, "wip/a.md", "0\n")
	w := &Watcher{Root: root, Paths: []string{"wip"}, Every: 40 * time.Millisecond}
	c := &collector{}
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})
	go func() { w.Run(ctx, c.emit); close(done) }()
	t.Cleanup(func() { cancel(); <-done })
	time.Sleep(60 * time.Millisecond)
	content := "0\n"
	for i := 0; i < 8; i++ { // a write every 15 ms: faster than the tick
		content += "more\n"
		write(t, root, "wip/a.md", content)
		time.Sleep(15 * time.Millisecond)
	}
	c.wait(t, "wip/a.md")
	time.Sleep(150 * time.Millisecond)
	c.wait(t, "wip/a.md")
}
