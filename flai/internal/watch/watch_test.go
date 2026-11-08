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

// write puts a file in place in one step: the content goes to a scratch file
// beside the watched root (not under it, so the watcher cannot see it) and is
// renamed over the target. os.WriteFile truncates first, so a tick landing
// between the truncate and the write would see an empty file, and a test
// goroutine that stalls there makes that a report of its own.
func write(t *testing.T, root, rel, content string) {
	t.Helper()
	full := filepath.Join(root, rel)
	if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
		t.Fatal(err)
	}
	tmp, err := os.CreateTemp(filepath.Dir(root), "write-*")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := tmp.WriteString(content); err != nil {
		t.Fatal(err)
	}
	if err := tmp.Close(); err != nil {
		t.Fatal(err)
	}
	if err := os.Rename(tmp.Name(), full); err != nil {
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

// slowFixture is a project with a fast item folder and a slow archive below
// the fast wip folder, as flai serve watches them, on a clock the test moves
// and a stat that records each path the slow tier statted.
type slowFixture struct {
	root    string
	w       *Watcher
	clock   time.Time
	mu      sync.Mutex
	statted []string
}

func newSlowFixture(t *testing.T) *slowFixture {
	t.Helper()
	f := &slowFixture{root: t.TempDir(), clock: time.Date(2026, 11, 8, 0, 0, 0, 0, time.UTC)}
	write(t, f.root, "wip/kanban/stories/S-0001-a.md", "a\n")
	write(t, f.root, "wip/archive/stories/S-0002-b.md", "b\n")
	write(t, f.root, "wip/archive/tasks/T-0001-c.md", "c\n")
	// A folder's modification time can be as coarse as a clock tick, so a
	// change made just after the first look could leave it the same: the
	// folders start in the past.
	old := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	for _, dir := range []string{"wip/archive", "wip/archive/stories", "wip/archive/tasks"} {
		if err := os.Chtimes(filepath.Join(f.root, dir), old, old); err != nil {
			t.Fatal(err)
		}
	}
	f.w = &Watcher{Root: f.root, Paths: []string{"wip"}, Slow: []string{"wip/archive"}, SlowEvery: 10 * time.Second,
		now: func() time.Time { return f.clock },
		stat: func(path string) (os.FileInfo, error) {
			f.mu.Lock()
			f.statted = append(f.statted, path)
			f.mu.Unlock()
			return os.Lstat(path)
		}}
	return f
}

// filesStatted are the regular files the slow tier statted since the last call.
func (f *slowFixture) filesStatted(t *testing.T) []string {
	t.Helper()
	f.mu.Lock()
	defer f.mu.Unlock()
	var files []string
	for _, p := range f.statted {
		if info, err := os.Lstat(p); err == nil && !info.IsDir() {
			files = append(files, p)
		}
	}
	f.statted = nil
	return files
}

// editInPlace rewrites a file without touching its folder, and gives it a
// modification time of its own, as an editor saving in place does.
func editInPlace(t *testing.T, root, rel, content string, at time.Time) {
	t.Helper()
	full := filepath.Join(root, rel)
	if err := os.WriteFile(full, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.Chtimes(full, at, at); err != nil {
		t.Fatal(err)
	}
}

func TestASlowPathStatsOnlyItsFoldersBetweenSlowLooks(t *testing.T) {
	f := newSlowFixture(t)
	first := f.w.snapshot()
	for _, rel := range []string{"wip/kanban/stories/S-0001-a.md", "wip/archive/stories/S-0002-b.md", "wip/archive/tasks/T-0001-c.md"} {
		if _, ok := first[rel]; !ok {
			t.Fatalf("first look misses %s: %v", rel, first)
		}
	}
	if got := f.filesStatted(t); len(got) != 2 {
		t.Fatalf("the first look statted %v, want the two archived files", got)
	}
	f.clock = f.clock.Add(9 * time.Second)
	again := f.w.snapshot()
	if got := f.filesStatted(t); len(got) != 0 {
		t.Fatalf("a look between slow looks statted files %v", got)
	}
	if len(again) != len(first) {
		t.Fatalf("a look between slow looks has %v, want %v", again, first)
	}
	f.clock = f.clock.Add(time.Second)
	f.w.snapshot()
	if got := f.filesStatted(t); len(got) != 2 {
		t.Fatalf("the slow look statted %v, want the two archived files", got)
	}
}

func TestAFileMovedIntoASlowPathIsReportedWithinAFastTickOrTwo(t *testing.T) {
	f := newSlowFixture(t)
	d := newDebounce(f.w.snapshot())
	scratch := filepath.Join(t.TempDir(), "T-0002-d.md")
	if err := os.WriteFile(scratch, []byte("d\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.Rename(scratch, filepath.Join(f.root, "wip/archive/tasks/T-0002-d.md")); err != nil {
		t.Fatal(err)
	}
	expect(t, "moved in", d.tick(f.w.snapshot()))
	expect(t, "next fast tick", d.tick(f.w.snapshot()), "wip/archive/tasks/T-0002-d.md")

	if err := os.Remove(filepath.Join(f.root, "wip/archive/stories/S-0002-b.md")); err != nil {
		t.Fatal(err)
	}
	expect(t, "removed", d.tick(f.w.snapshot()), "wip/archive/stories/S-0002-b.md")
	expect(t, "afterwards", d.tick(f.w.snapshot()))
}

func TestAnEditInPlaceUnderASlowPathIsReportedAtTheSlowInterval(t *testing.T) {
	for _, tc := range []struct{ name, content string }{{"same size", "x\n"}, {"another size", "longer\n"}} {
		t.Run(tc.name, func(t *testing.T) {
			f := newSlowFixture(t)
			d := newDebounce(f.w.snapshot())
			editInPlace(t, f.root, "wip/archive/stories/S-0002-b.md", tc.content, time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC))
			// the fast item folder beside it is still seen at once
			editInPlace(t, f.root, "wip/kanban/stories/S-0001-a.md", tc.content, time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC))
			for i := range 5 {
				f.clock = f.clock.Add(time.Second)
				want := []string{}
				if i == 1 {
					want = []string{"wip/kanban/stories/S-0001-a.md"}
				}
				expect(t, "before the slow interval", d.tick(f.w.snapshot()), want...)
			}
			f.clock = f.clock.Add(5 * time.Second)
			expect(t, "at the slow interval", d.tick(f.w.snapshot()))
			expect(t, "next fast tick", d.tick(f.w.snapshot()), "wip/archive/stories/S-0002-b.md")
		})
	}
}

func TestASlowPathBelowAFastPathIsLeftOutOfTheFastWalk(t *testing.T) {
	f := newSlowFixture(t)
	f.w.snapshot()
	f.filesStatted(t)
	editInPlace(t, f.root, "wip/archive/tasks/T-0001-c.md", "changed, and longer\n", time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC))
	before, err := os.Stat(filepath.Join(f.root, "wip/archive/tasks/T-0001-c.md"))
	if err != nil {
		t.Fatal(err)
	}
	got := f.w.snapshot()["wip/archive/tasks/T-0001-c.md"]
	if got == sigOf(before) {
		t.Fatal("the fast walk statted a file below the slow path")
	}
	if files := f.filesStatted(t); len(files) != 0 {
		t.Fatalf("statted %v between slow looks", files)
	}
}

func TestRunReportsAFileMovedIntoAndEditedUnderASlowPath(t *testing.T) {
	root := t.TempDir()
	write(t, root, "wip/archive/stories/S-0001-a.md", "a\n")
	w := &Watcher{Root: root, Paths: []string{"wip"}, Slow: []string{"wip/archive"}, Every: 20 * time.Millisecond, SlowEvery: 200 * time.Millisecond}
	c := &collector{}
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})
	go func() { w.Run(ctx, c.emit); close(done) }()
	t.Cleanup(func() { cancel(); <-done })
	time.Sleep(60 * time.Millisecond)
	c.wait(t)

	write(t, root, "wip/archive/stories/S-0002-b.md", "b\n")
	c.wait(t, "wip/archive/stories/S-0002-b.md")
	c.reset()

	editInPlace(t, root, "wip/archive/stories/S-0001-a.md", "a, by hand\n", time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC))
	c.wait(t, "wip/archive/stories/S-0001-a.md")
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
