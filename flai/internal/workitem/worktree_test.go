package workitem

import (
	"os"
	"path/filepath"
	"testing"
)

func TestOpenFromLinkedWorktreeUsesMainWip(t *testing.T) {
	main := newProject(t)
	wt := filepath.Join(main.Root, ".flai-cache", "worktrees", "S-0001")
	_ = os.MkdirAll(wt, 0o755)
	data, _ := os.ReadFile(filepath.Join(main.Root, "system-flow.yaml"))
	_ = os.WriteFile(filepath.Join(wt, "system-flow.yaml"), data, 0o644)
	_ = os.WriteFile(filepath.Join(wt, ".git"), []byte("gitdir: "+filepath.Join(main.Root, ".git", "worktrees", "S-0001")+"\n"), 0o644)
	r, err := Open(wt)
	if err != nil {
		t.Fatal(err)
	}
	if r.Root != wt || r.MainRoot != main.Root || !r.IsWorktree() {
		t.Fatalf("roots: %+v", r)
	}
	if r.WipDir() != main.WipDir() || r.CacheDir() != filepath.Join(main.Root, ".flai-cache") {
		t.Errorf("wip must resolve to the main checkout: %s", r.WipDir())
	}
	if main.IsWorktree() || main.MainRoot != main.Root {
		t.Errorf("a normal checkout is its own main: %+v", main)
	}
}
