package designsystem

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/bytepunx/system-flow/flai/internal/execx"
	"github.com/bytepunx/system-flow/flai/internal/storygit"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// Proposal is an open story whose branch, or whose worktree uncommitted,
// changes the design system, as flai design proposals --json and the host
// method design.proposals list it.
type Proposal struct {
	ID     string   `json:"id"`
	Title  string   `json:"title"`
	Status string   `json:"status"` // its state, never done or cancelled
	Paths  []string `json:"paths"`  // the token file, the guidance, or both, in that order
}

// ShowStory reads the design system from the worktree of the open story
// named story, at the paths its worktree's manifest gives, and answers it
// with the story's own ID. An ID that is no work item, an item that is not a
// story, a story that is done or cancelled, and a story without a worktree
// are errors.
func ShowStory(repo *workitem.Repo, story string) (*Shown, string, error) {
	id, in, err := storyWorktree(repo, story)
	if err != nil {
		return nil, "", err
	}
	shown, err := Show(os.DirFS(in.Root), in.Manifest.DesignTokens(), in.Manifest.DesignGuidance())
	if err != nil {
		return nil, "", fmt.Errorf("read the design system in %s's worktree %s: %w", id, in.Root, err)
	}
	shown.Stylesheet = in.Manifest.DesignStylesheet()
	return shown, id, nil
}

// Refused is ShowStory's refusal of a story whose design system it does not
// read: an item that is not a story, a story done or cancelled, or one
// without a worktree that is a checkout of the project.
type Refused struct{ msg string }

func (e *Refused) Error() string { return e.msg }

// refused is a Refused saying why, as fmt.Sprintf formats it.
func refused(format string, args ...any) error {
	return &Refused{msg: fmt.Sprintf(format, args...)}
}

// storyWorktree is story's canonical ID and the repository rooted at its
// worktree, refusing an item that is not a story, a story that is not open,
// and a story that has no worktree.
func storyWorktree(repo *workitem.Repo, story string) (string, *workitem.Repo, error) {
	it, err := repo.Get(story)
	if err != nil {
		return "", nil, fmt.Errorf("read the design system of %s: %w; name a story of this project", story, err)
	}
	if it.Type != workitem.Story {
		return "", nil, refused("read the design system of %s: it is a %s; name a story (S-nnnn), whose worktree holds the design system it proposes", it.ID, it.Type)
	}
	if it.Closed() {
		return "", nil, refused("read the design system of %s: it is %s, so it proposes nothing; name an open story, or read the main checkout's without a story", it.ID, it.Status)
	}
	wt := repo.WorktreePath(it.ID)
	if st, err := os.Stat(wt); err != nil || !st.IsDir() {
		return "", nil, refused("read the design system of %s: it has no worktree at %s; open it with flai stream open %s, or read the main checkout's without a story", it.ID, wt, it.ID)
	}
	in, err := workitem.Open(wt)
	if err != nil || filepath.Clean(in.Root) != filepath.Clean(wt) {
		return "", nil, refused("read the design system of %s: its worktree %s is not a checkout of the project; open it again with flai stream open %s", it.ID, wt, it.ID)
	}
	return it.ID, in, nil
}

// Proposals lists the open stories whose branches, or whose worktrees
// uncommitted, change the token file or the guidance at the paths the main
// checkout's manifest gives, in story ID order. None is an empty list, never
// nil.
func Proposals(r execx.Runner, repo *workitem.Repo) ([]Proposal, error) {
	items, err := repo.List(false)
	if err != nil {
		return nil, fmt.Errorf("list the stories that propose design system changes: %w; run flai check to find the item that does not read", err)
	}
	designPaths := []string{path.Clean(repo.Manifest.DesignTokens()), path.Clean(repo.Manifest.DesignGuidance())}
	out := []Proposal{}
	for _, it := range items {
		if it.Type != workitem.Story || it.Closed() {
			continue
		}
		changed, err := storygit.StoryChanges(r, repo, it.ID)
		if err != nil {
			return nil, fmt.Errorf("list the stories that propose design system changes: %w; run flai stream diff %s to see what git reports", err, it.ID)
		}
		if paths := changedOf(designPaths, changed); len(paths) > 0 {
			out = append(out, Proposal{ID: it.ID, Title: it.Title, Status: it.Status, Paths: paths})
		}
	}
	return out, nil
}

// changedOf is those of want that changed names, in want's order. A changed
// name ending in a slash is a folder git does not track yet, as git status
// lists it, and covers every path under it.
func changedOf(want, changed []string) []string {
	seen := make(map[string]bool, len(changed))
	var folders []string
	for _, p := range changed {
		p = filepath.ToSlash(p)
		if strings.HasSuffix(p, "/") {
			folders = append(folders, path.Clean(p)+"/")
			continue
		}
		seen[path.Clean(p)] = true
	}
	var out []string
	for _, p := range want {
		if seen[p] || under(p, folders) {
			out = append(out, p)
		}
	}
	return out
}

// under reports whether p lies in one of folders, each ending in a slash.
func under(p string, folders []string) bool {
	for _, f := range folders {
		if strings.HasPrefix(p, f) {
			return true
		}
	}
	return false
}
