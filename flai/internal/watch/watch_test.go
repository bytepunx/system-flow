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

// snap is a snapshot of files by name; a file's size stands in for its whole
// signature, since the debounce only compares signatures for equality.
func snap(sizes map[string]int64) map[string]sig {
	out := map[string]sig{}
	for rel, size := range sizes {
		out[rel] = sig{mod: size, size: size}
	}
	return out
}

func expect(t *testing.T, step string, got []string, want ...string) {
	t.Helper()
	if len(got) != len(want) {
		t.Fatalf("%s: reported %v, want %v", step, got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("%s: reported %v, want %v", step, got, want)
		}
	}
}

// The debounce is fed snapshots, one per tick, so these tests pass or fail the
// same on a machine that stalls for seconds as on one that never does.

func TestAFileStillBeingWrittenIsReportedOnce(t *testing.T) {
	d := newDebounce(snap(map[string]int64{"wip/a.md": 1}))
	for size := int64(2); size <= 9; size++ { // a different size at every tick: still being written
		expect(t, "while written", d.tick(snap(map[string]int64{"wip/a.md": size})))
	}
	expect(t, "first tick it holds still", d.tick(snap(map[string]int64{"wip/a.md": 9})), "wip/a.md")
	expect(t, "afterwards", d.tick(snap(map[string]int64{"wip/a.md": 9})))
	expect(t, "and again", d.tick(snap(map[string]int64{"wip/a.md": 9})))
}

func TestAFileThatSettlesAndIsWrittenAgainIsReportedForEachSettle(t *testing.T) {
	d := newDebounce(snap(map[string]int64{"a": 1}))
	expect(t, "changed", d.tick(snap(map[string]int64{"a": 2})))
	expect(t, "settled", d.tick(snap(map[string]int64{"a": 2})), "a")
	expect(t, "changed again", d.tick(snap(map[string]int64{"a": 3})))
	expect(t, "settled again", d.tick(snap(map[string]int64{"a": 3})), "a")
}

func TestAFileThatChangesBackWhilePendingIsStillReported(t *testing.T) {
	d := newDebounce(snap(map[string]int64{"a": 1}))
	expect(t, "changed", d.tick(snap(map[string]int64{"a": 2})))
	expect(t, "changed back", d.tick(snap(map[string]int64{"a": 1})))
	expect(t, "settled at the old value", d.tick(snap(map[string]int64{"a": 1})), "a")
}

func TestAnAddedFileIsReportedOnceItHoldsStill(t *testing.T) {
	d := newDebounce(snap(nil))
	expect(t, "appeared", d.tick(snap(map[string]int64{"new.md": 4})))
	expect(t, "settled", d.tick(snap(map[string]int64{"new.md": 4})), "new.md")
}

func TestARemovedFileIsReportedAtOnce(t *testing.T) {
	d := newDebounce(snap(map[string]int64{"a": 1, "b": 1}))
	expect(t, "b removed", d.tick(snap(map[string]int64{"a": 1})), "b")
	expect(t, "afterwards", d.tick(snap(map[string]int64{"a": 1})))
}

func TestAFileRemovedWhileItsChangeIsPendingIsReportedOnce(t *testing.T) {
	d := newDebounce(snap(map[string]int64{"a": 1}))
	expect(t, "changed", d.tick(snap(map[string]int64{"a": 2})))
	expect(t, "removed", d.tick(snap(nil)), "a")
	expect(t, "afterwards", d.tick(snap(nil)))
}

func TestSeveralFilesAreReportedSortedAndIndependently(t *testing.T) {
	d := newDebounce(snap(map[string]int64{"b": 1, "a": 1, "c": 1}))
	expect(t, "all three change", d.tick(snap(map[string]int64{"b": 2, "a": 2, "c": 1})))
	expect(t, "b keeps changing, a settles, c untouched", d.tick(snap(map[string]int64{"b": 3, "a": 2, "c": 1})), "a")
	expect(t, "b settles", d.tick(snap(map[string]int64{"b": 3, "a": 2, "c": 1})), "b")
}
