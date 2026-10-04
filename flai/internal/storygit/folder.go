package storygit

import (
	"maps"
	"os"
	"path"
	"path/filepath"
	"slices"
	"strings"

	"github.com/bytepunx/system-flow/flai/internal/execx"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// FolderNames is the sorted names of the files directly under folder, a path
// from the project's root, in any worktree, on main, or on a story branch.
func FolderNames(r execx.Runner, repo *workitem.Repo, folder string) []string {
	root := repo.MainRoot
	if root == "" {
		root = repo.Root
	}
	folder = strings.TrimSuffix(folder, "/")
	names := map[string]bool{}
	readTree := func(dir string) {
		entries, _ := os.ReadDir(filepath.Join(dir, filepath.FromSlash(folder)))
		for _, e := range entries {
			if !e.IsDir() {
				names[e.Name()] = true
			}
		}
	}
	readTree(root)

	// A worktree is listed by its repository's top; the project may sit in a
	// folder of it.
	if prefix, err := r.Run(root, "git", "rev-parse", "--show-prefix"); err == nil {
		if out, err := r.Run(root, "git", "worktree", "list", "--porcelain"); err == nil {
			for _, l := range strings.Split(out, "\n") {
				if top, ok := strings.CutPrefix(l, "worktree "); ok {
					readTree(filepath.Join(top, filepath.FromSlash(prefix)))
				}
			}
		}
	}

	var refs []string
	if main, err := MainBranch(r, root); err == nil {
		refs = append(refs, main)
	}
	if out, err := r.Run(root, "git", "for-each-ref", "--format=%(refname)", "refs/heads/"+Prefix); err == nil {
		refs = append(refs, strings.Fields(out)...)
	}
	for _, ref := range refs {
		// ls-tree reads paths from the working directory, so a project in a
		// folder of its repository names its own; -z leaves names unquoted.
		out, err := r.Run(root, "git", "ls-tree", "-z", ref, "--", folder+"/")
		if err != nil {
			continue
		}
		for _, entry := range strings.Split(out, "\x00") {
			meta, name, ok := strings.Cut(entry, "\t")
			if ok && strings.Contains(meta, " blob ") {
				names[path.Base(name)] = true
			}
		}
	}
	return slices.Sorted(maps.Keys(names))
}
