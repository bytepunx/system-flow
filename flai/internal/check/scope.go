package check

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/bytepunx/system-flow/flai/internal/threads"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// ScopeToStory marks each finding of res outside story as a note the run
// passes over, errors as well as warnings (S-0249, I-0057). A finding is
// inside when its path is the story's item file or one of its tasks', its
// narrative, a thread anchored on the story or one of its tasks, or one of
// changed: paths relative to the project root, as git names them, where a
// path ending in / is a folder. A wip.overlap is always outside: clearing it
// is the pull hold's business and the other story's (ADR-0019, ADR-0046).
func ScopeToStory(res *Result, repo *workitem.Repo, story string, changed []string) error {
	st, err := repo.Get(story)
	if err != nil {
		return fmt.Errorf("scope the check to %s: %w", story, err)
	}
	if st.Type != workitem.Story {
		return fmt.Errorf("scope the check to %s: it is a %s; name a story (S-nnnn)", st.ID, st.Type)
	}
	in, err := storyPaths(repo, st)
	if err != nil {
		return fmt.Errorf("scope the check to %s: %w", st.ID, err)
	}
	var folders []string
	for _, p := range changed {
		for _, root := range []string{repo.Root, repo.MainRoot} {
			if root == "" {
				continue
			}
			abs := filepath.Join(root, filepath.FromSlash(p))
			if strings.HasSuffix(p, "/") {
				folders = append(folders, abs+string(filepath.Separator))
			}
			in[abs] = true
		}
	}
	inside := func(path string) bool {
		if !filepath.IsAbs(path) {
			path = filepath.Join(repo.Root, path)
		}
		path = filepath.Clean(path)
		if in[path] {
			return true
		}
		for _, f := range folders {
			if strings.HasPrefix(path, f) {
				return true
			}
		}
		return false
	}
	for i := range res.Findings {
		f := &res.Findings[i]
		if f.Outside || (f.Rule != "wip.overlap" && inside(f.Path)) {
			continue
		}
		f.Outside = true
		res.Outside++
		switch {
		case f.Level == Error:
			res.outsideErrors++
		case f.advisory:
			// counted once, as outside: Advisory keeps those inside
			res.Advisory--
			res.outsideWarnings++
		default:
			res.outsideWarnings++
		}
	}
	return nil
}

// storyPaths is the cleaned paths of the files that belong to st: its item
// file and its tasks', its narrative, and the threads anchored on any of them.
func storyPaths(repo *workitem.Repo, st *workitem.Item) (map[string]bool, error) {
	items, err := repo.List(true)
	if err != nil {
		return nil, err
	}
	in := map[string]bool{filepath.Clean(st.Path): true, filepath.Clean(repo.NarrativePath(st.ID)): true}
	ids := map[string]bool{workitem.CanonicalID(st.ID): true}
	for _, t := range workitem.Children(items, st.ID) {
		in[filepath.Clean(t.Path)] = true
		ids[workitem.CanonicalID(t.ID)] = true
	}
	// A thread folder that does not read is reported on the folder by the
	// threads rule, which then checks no thread: none is the story's.
	list, err := threads.List(repo)
	if err != nil {
		return in, nil //nolint:nilerr // the threads rule reports it; no thread is inside
	}
	for _, th := range list {
		if th.Anchor.Item != "" && ids[workitem.CanonicalID(th.Anchor.Item)] {
			in[filepath.Clean(th.Path)] = true
		}
	}
	return in, nil
}
