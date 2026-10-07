package check

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/bytepunx/system-flow/flai/internal/messages"
	"github.com/bytepunx/system-flow/flai/internal/threads"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// ScopeToStory marks each finding of res outside story as a note the run
// passes over, errors as well as warnings (S-0249, I-0057). A finding is
// inside when its path is the story's item file or one of its tasks', its
// narrative, a thread anchored on the story or one of its tasks, or one of
// changed: paths relative to the project root, as git names them, where a
// path ending in / is a folder. A wip.overlap that names the story in its
// Stories is outside it, a note: clearing it is the pull hold's business and
// the other story's (ADR-0019, ADR-0046). One that does not name the story is
// left out of the scoped result, and the counts it added are taken back
// (ADR-0115). So is an item.archive outside the story: the item waits for
// the operator's flai archive in the main checkout, which no story branch
// clears, and it never names a story in progress (ADR-0122, I-0078).
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
	kept := res.Findings[:0]
	for _, f := range res.Findings {
		if (f.Rule == "wip.overlap" && !namesStory(f, st.ID)) || (f.Rule == "item.archive" && !inside(f.Path)) {
			res.drop(f)
			continue
		}
		if !f.Outside && (f.Rule == "wip.overlap" || !inside(f.Path)) {
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
		kept = append(kept, f)
	}
	clear(res.Findings[len(kept):])
	res.Findings = kept
	return nil
}

// namesStory reports whether f, a wip.overlap, names story in its Stories
// (ADR-0115); the message is never read.
func namesStory(f Finding, story string) bool {
	for _, id := range f.Stories {
		if workitem.CanonicalID(id) == workitem.CanonicalID(story) {
			return true
		}
	}
	return false
}

// drop takes back the counts f added to r, as left out of the result.
func (r *Result) drop(f Finding) {
	if f.Level == Error {
		r.Errors--
	} else {
		r.Warnings--
	}
	switch {
	case f.Outside:
		r.Outside--
		if f.Level == Error {
			r.outsideErrors--
		} else {
			r.outsideWarnings--
		}
	case f.advisory:
		r.Advisory--
	}
}

// storyPaths is the cleaned paths of the files that belong to st: its item
// file and its tasks', its narrative, the threads anchored on any of them,
// and the conversations it is one of the two stories of (ADR-0120).
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
	// A conversation that does not read is reported by the messages rule on
	// its own file; none is the story's then.
	convs, err := messages.List(repo)
	if err != nil {
		return in, nil //nolint:nilerr // the messages rule reports it; no conversation is inside
	}
	for _, c := range convs {
		if c.Names(st.ID) {
			in[filepath.Clean(c.Path)] = true
		}
	}
	return in, nil
}
