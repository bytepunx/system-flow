package workitem

import (
	"os"
	"path/filepath"
	"strings"
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

func TestOpenFromLinkedWorktreeWithMainAbsentRecordsIt(t *testing.T) {
	main := newProject(t)
	absent := filepath.Join(t.TempDir(), "system-flow")
	wt := filepath.Join(absent, ".flai-cache", "worktrees", "S-0001")
	if err := os.MkdirAll(wt, 0o755); err != nil {
		t.Fatal(err)
	}
	data, _ := os.ReadFile(filepath.Join(main.Root, "system-flow.yaml"))
	_ = os.WriteFile(filepath.Join(wt, "system-flow.yaml"), data, 0o644)
	_ = os.WriteFile(filepath.Join(wt, ".git"), []byte("gitdir: "+filepath.Join(absent, ".git", "worktrees", "S-0001")+"\n"), 0o644)
	r, err := Open(wt)
	if err != nil {
		t.Fatal(err)
	}
	if r.Root != wt || r.MainRoot != wt || r.IsWorktree() || r.MainAbsent != absent {
		t.Fatalf("roots: %+v", r)
	}
	if r.WipDir() != filepath.Join(wt, "wip") || r.CacheDir() != filepath.Join(wt, ".flai-cache") {
		t.Errorf("wip and cache keep today's values: %s %s", r.WipDir(), r.CacheDir())
	}
	err = r.RequireMain("flai verify")
	if err == nil || !strings.Contains(err.Error(), "flai verify") || !strings.Contains(err.Error(), absent) || !strings.Contains(err.Error(), "MCP tool verify") {
		t.Errorf("RequireMain must name the command, the main checkout, and the MCP tool: %v", err)
	}
	if got := r.WorktreePath("S-0001"); got != wt {
		t.Errorf("the worktree's own story is Root: %s", got)
	}
	if got := r.WorktreePath("S-0002"); got != filepath.Join(wt, ".flai-cache", "worktrees", "S-0002") {
		t.Errorf("another story's worktree is unchanged: %s", got)
	}
	if err := main.RequireMain("flai verify"); err != nil {
		t.Errorf("a normal checkout has its main: %v", err)
	}
	present := filepath.Join(main.Root, ".flai-cache", "worktrees", "S-0003")
	_ = os.MkdirAll(present, 0o755)
	_ = os.WriteFile(filepath.Join(present, "system-flow.yaml"), data, 0o644)
	_ = os.WriteFile(filepath.Join(present, ".git"), []byte("gitdir: "+filepath.Join(main.Root, ".git", "worktrees", "S-0003")+"\n"), 0o644)
	p, err := Open(present)
	if err != nil {
		t.Fatal(err)
	}
	if p.MainAbsent != "" || p.RequireMain("flai verify") != nil || p.WorktreePath("S-0003") != present {
		t.Errorf("a worktree with its main present needs nothing more: %+v", p)
	}
}
