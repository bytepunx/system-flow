package storygit

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/bytepunx/system-flow/flai/internal/execx"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// Limits on what a diff carries, so a large branch cannot swamp a browser.
const (
	DiffFileLimit  = 64 * 1024
	DiffTotalLimit = 768 * 1024
)

// DiffFile is one file a story branch changed.
type DiffFile struct {
	Path      string `json:"path"`
	OldPath   string `json:"old_path,omitempty"` // for a rename
	Status    string `json:"status"`             // added, modified, deleted, renamed
	Additions int    `json:"additions"`
	Deletions int    `json:"deletions"`
	Binary    bool   `json:"binary"`
	Truncated bool   `json:"truncated"` // the patch was cut at the limit
	Patch     string `json:"patch"`     // unified hunks, empty for a binary file
}

// Diff is a story branch against the main branch, as flai stream diff
// prints it.
type Diff struct {
	Story     string     `json:"story"`
	Branch    string     `json:"branch"`
	Base      string     `json:"base"` // the merge base, abbreviated
	Commits   int        `json:"commits"`
	Files     []DiffFile `json:"files"`
	Additions int        `json:"additions"`
	Deletions int        `json:"deletions"`
	Truncated bool       `json:"truncated"` // some patches were cut or left out
}

// StoryChanges is the paths the story's branch changes against the main
// branch, both sides of a rename, and those uncommitted in its worktree, as
// flai check --story counts them inside the story (S-0249). A project outside
// version control, or a story without a branch or a worktree, changes none; a
// branch or a worktree that cannot be read is an error.
func StoryChanges(r execx.Runner, repo *workitem.Repo, id string) ([]string, error) {
	if repo.MainRoot == "" || !InWorkTree(r, repo.MainRoot) {
		return nil, nil
	}
	var paths []string
	if BranchExists(r, repo.MainRoot, Branch(id)) {
		d, err := StoryDiff(r, repo, id)
		if err != nil {
			return nil, fmt.Errorf("read what %s changes: %w", id, err)
		}
		for _, f := range d.Files {
			paths = append(paths, f.Path)
			if f.OldPath != "" {
				paths = append(paths, f.OldPath)
			}
		}
	}
	wt := repo.WorktreePath(id)
	if st, err := os.Stat(wt); err == nil && st.IsDir() {
		dirty, err := Uncommitted(r, wt)
		if err != nil {
			return nil, fmt.Errorf("read what %s has uncommitted in %s: %w", id, wt, err)
		}
		paths = append(paths, dirty...)
	}
	return paths, nil
}

// StoryDiff is the story's branch against its merge base with the main
// branch, read with git in the main checkout, so it does not depend on the
// worktree.
func StoryDiff(r execx.Runner, repo *workitem.Repo, id string) (*Diff, error) {
	root := repo.MainRoot
	if root == "" || !InWorkTree(r, root) {
		return nil, fmt.Errorf("%s has no branch to compare: this project is not a git repository", id)
	}
	branch := Branch(id)
	if !BranchExists(r, root, branch) {
		return nil, fmt.Errorf("%s has no branch %s: it was never opened with flai stream open, or it has been merged and removed", id, branch)
	}
	main, err := MainBranch(r, root)
	if err != nil {
		return nil, err
	}
	base, err := r.Run(root, "git", "merge-base", main, branch)
	if err != nil {
		return nil, err
	}
	d := &Diff{Story: id, Branch: branch, Files: []DiffFile{}}
	if short, err := r.Run(root, "git", "rev-parse", "--short", base); err == nil {
		d.Base = short
	}
	if n, err := r.Run(root, "git", "rev-list", "--count", base+".."+branch); err == nil {
		d.Commits, _ = strconv.Atoi(n)
	}
	status, err := r.Run(root, "git", "diff", "--name-status", "-M", "-z", base, branch)
	if err != nil {
		return nil, err
	}
	byPath := map[string]*DiffFile{}
	fields := strings.Split(strings.TrimRight(status, "\x00"), "\x00")
	for i := 0; i < len(fields) && fields[i] != ""; i++ {
		f := DiffFile{}
		switch code := fields[i]; {
		case strings.HasPrefix(code, "R") && i+2 < len(fields):
			f.Status, f.OldPath, f.Path = "renamed", fields[i+1], fields[i+2]
			i += 2
		case i+1 < len(fields):
			f.Status = map[byte]string{'A': "added", 'D': "deleted", 'M': "modified"}[code[0]]
			if f.Status == "" {
				f.Status = "modified"
			}
			f.Path = fields[i+1]
			i++
		default:
			continue
		}
		d.Files = append(d.Files, f)
	}
	for i := range d.Files {
		byPath[d.Files[i].Path] = &d.Files[i]
	}
	numstat, err := r.Run(root, "git", "diff", "--numstat", "-M", "-z", base, branch)
	if err != nil {
		return nil, err
	}
	// records are "add\tdel\tpath\0", or for a rename "add\tdel\t\0old\0new\0"
	recs := strings.Split(strings.TrimRight(numstat, "\x00"), "\x00")
	for i := 0; i < len(recs); i++ {
		parts := strings.SplitN(recs[i], "\t", 3)
		if len(parts) != 3 {
			continue
		}
		path := parts[2]
		if path == "" && i+2 < len(recs) {
			path = recs[i+2]
			i += 2
		}
		f := byPath[path]
		if f == nil {
			continue
		}
		if parts[0] == "-" {
			f.Binary = true
			continue
		}
		f.Additions, _ = strconv.Atoi(parts[0])
		f.Deletions, _ = strconv.Atoi(parts[1])
		d.Additions += f.Additions
		d.Deletions += f.Deletions
	}
	total := 0
	for i := range d.Files {
		f := &d.Files[i]
		if f.Binary {
			continue
		}
		if total >= DiffTotalLimit {
			f.Truncated, d.Truncated = true, true
			continue
		}
		paths := []string{f.Path}
		if f.OldPath != "" {
			paths = append(paths, f.OldPath)
		}
		patch, err := r.Run(root, "git", append([]string{"diff", "--no-color", "-M", base, branch, "--"}, paths...)...)
		if err != nil {
			return nil, err
		}
		// keep the hunks; the header lines repeat what the fields say
		if at := strings.Index(patch, "\n@@"); at >= 0 {
			patch = patch[at+1:]
		} else if !strings.HasPrefix(patch, "@@") {
			patch = ""
		}
		if len(patch) > DiffFileLimit {
			patch = patch[:strings.LastIndex(patch[:DiffFileLimit], "\n")+1]
			f.Truncated, d.Truncated = true, true
		}
		f.Patch = patch
		total += len(patch)
	}
	return d, nil
}
