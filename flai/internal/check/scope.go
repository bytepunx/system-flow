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
// clears, and it never names a story in progress (ADR-0122, I-0078). So is a
// board.wip-limit outside the story, advisory or not: it counts the main
// checkout's columns, which only the operator's acceptances, cancellations,
// or a raised limit clear (ADR-0133, I-0123). So is an issues.no-story
// outside the story: it names an issue the story did not touch, and only the
// operator's choice to make a story for it clears it (ADR-0136, I-0134). So
// is a threads.archived outside the story: it names an open or answered
// thread on an archived item the story did not touch, which only the operator
// or the thread's opener resolves, in the main checkout, which no story
// branch changes (ADR-0137, I-0135); one on a thread anchored on the story or
// one of its tasks is inside it. So is the note components.examples gives
// when it could not check examples for want of node or a sub-project's
// svelte, unless one of changed is under the examples folder it names, or a
// folder of changed holds it: only installing the tool clears it, which no
// story branch does; when one is, the story's own examples went unchecked,
// and the note is inside it (ADR-0171, I-0161). So is a markdown or
// narrative.state finding on the narrative of another story that is neither
// done nor cancelled: that story's own close-out finds it (ADR-0123, I-0096;
// ADR-0125, I-0109). So is a touches.folder on another story, and a
// wip.after-hub whose hub is another story, by the story in its Stories, not
// by the paths changed: only that story's tasks clear the one, and only the
// stories that name the hub clear the other; on the story itself both are
// kept, advisory (S-0468, ADR-0188, I-0152).
func ScopeToStory(res *Result, repo *workitem.Repo, story string, changed []string) error {
	return ScopeToStories(res, repo, []string{story}, changed)
}

// ScopeToStories is ScopeToStory for several stories, as a merge queue batch
// names them (I-0174): a finding is inside when it is inside any of them, and
// outside, a note or left out, only when it is outside every one. changed is
// the paths all of them change. A wip.overlap that names any of them is
// outside; one that names none is left out, and so is a touches.folder or a
// wip.after-hub on a story that is none of them. A markdown or narrative.state
// finding is left out on the narrative of an open story that is none of them.
func ScopeToStories(res *Result, repo *workitem.Repo, stories []string, changed []string) error {
	if len(stories) == 0 {
		return fmt.Errorf("scope the check to a story: none given; name a story (S-nnnn)")
	}
	sts := make([]*workitem.Item, 0, len(stories))
	ids := make([]string, 0, len(stories))
	for _, story := range stories {
		st, err := repo.Get(story)
		if err != nil {
			return fmt.Errorf("scope the check to %s: %w", story, err)
		}
		if st.Type != workitem.Story {
			return fmt.Errorf("scope the check to %s: it is a %s; name a story (S-nnnn)", st.ID, st.Type)
		}
		sts = append(sts, st)
		ids = append(ids, st.ID)
	}
	scope := strings.Join(ids, ", ")
	in, err := storyPaths(repo, sts)
	if err != nil {
		return fmt.Errorf("scope the check to %s: %w", scope, err)
	}
	others, err := openNarratives(repo, ids)
	if err != nil {
		return fmt.Errorf("scope the check to %s: %w", scope, err)
	}
	var folders, paths []string
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
			paths = append(paths, abs)
		}
	}
	resolve := func(path string) string {
		if !filepath.IsAbs(path) {
			path = filepath.Join(repo.Root, path)
		}
		return filepath.Clean(path)
	}
	inside := func(path string) bool {
		path = resolve(path)
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
	// under reports whether a changed path is folder or under it, or a
	// changed folder holds it: an unchecked note's folder is no changed path
	// itself (ADR-0171).
	under := func(folder string) bool {
		folder = resolve(folder)
		for _, p := range paths {
			if p == folder || strings.HasPrefix(p, folder+string(filepath.Separator)) {
				return true
			}
		}
		for _, f := range folders {
			if strings.HasPrefix(folder+string(filepath.Separator), f) {
				return true
			}
		}
		return false
	}
	kept := res.Findings[:0]
	for _, f := range res.Findings {
		if ((f.Rule == "wip.overlap" || f.Rule == "touches.folder" || f.Rule == "wip.after-hub") && !namesStory(f, ids)) ||
			((f.Rule == "item.archive" || f.Rule == "board.wip-limit" || f.Rule == "issues.no-story" || f.Rule == "threads.archived") && !inside(f.Path)) ||
			(f.unchecked && !under(f.Path)) ||
			((strings.HasPrefix(f.Rule, "markdown.") || f.Rule == "narrative.state") && others[resolve(f.Path)]) {
			res.drop(f)
			continue
		}
		if !f.Outside && !f.unchecked && (f.Rule == "wip.overlap" || !inside(f.Path)) {
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

// ScopeToItem scopes res, a check run with a write of item it in place, to
// the item's story, the story itself or a task's parent, with written, the
// paths the write made, relative to the checkout, as changed (ScopeToStories),
// and returns whether a finding is outside that story: what a checked edit
// and a checked creation refuse on and what they only warn of (ADR-0197,
// I-0182, I-0191). An epic has no story, and a task whose story does not read
// is held to the epic's rule: a finding is outside unless it is on a path
// written.
func ScopeToItem(res *Result, repo *workitem.Repo, it *workitem.Item, written []string) (func(Finding) bool, error) {
	story := ""
	switch it.Type {
	case workitem.Story:
		story = it.ID
	case workitem.Task:
		if p, err := repo.Get(it.Parent); err == nil && p.Type == workitem.Story {
			story = p.ID
		}
	}
	if story != "" {
		if err := ScopeToStories(res, repo, []string{story}, written); err != nil {
			return nil, err
		}
		return func(f Finding) bool { return f.Outside }, nil
	}
	in := map[string]bool{}
	for _, p := range written {
		for _, r := range []string{repo.Root, repo.MainRoot} {
			if r != "" {
				in[filepath.Join(r, filepath.FromSlash(p))] = true
			}
		}
	}
	return func(f Finding) bool {
		path := f.Path
		if !filepath.IsAbs(path) {
			path = filepath.Join(repo.Root, path)
		}
		return !in[filepath.Clean(path)]
	}, nil
}

// namesStory reports whether f, a wip.overlap, a touches.folder, or a
// wip.after-hub, names any of stories in its Stories (ADR-0115, ADR-0188);
// the message is never read.
func namesStory(f Finding, stories []string) bool {
	for _, id := range f.Stories {
		for _, story := range stories {
			if workitem.CanonicalID(id) == workitem.CanonicalID(story) {
				return true
			}
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

// openNarratives is the cleaned paths of the narratives of the stories other
// than stories whose status is neither done nor cancelled (ADR-0123,
// ADR-0125).
func openNarratives(repo *workitem.Repo, stories []string) (map[string]bool, error) {
	items, err := repo.List(true)
	if err != nil {
		return nil, err
	}
	scoped := map[string]bool{}
	for _, story := range stories {
		scoped[workitem.CanonicalID(story)] = true
	}
	open := map[string]bool{}
	for _, it := range items {
		if it.Type != workitem.Story || it.Status == workitem.Done || it.Status == workitem.Cancelled ||
			scoped[workitem.CanonicalID(it.ID)] {
			continue
		}
		open[filepath.Clean(repo.NarrativePath(it.ID))] = true
	}
	return open, nil
}

// storyPaths is the cleaned paths of the files that belong to any of sts:
// its item file and its tasks', its narrative, the threads anchored on any of
// them, and the conversations it is one of the two stories of (ADR-0120).
func storyPaths(repo *workitem.Repo, sts []*workitem.Item) (map[string]bool, error) {
	items, err := repo.List(true)
	if err != nil {
		return nil, err
	}
	in := map[string]bool{}
	ids := map[string]bool{}
	for _, st := range sts {
		in[filepath.Clean(st.Path)] = true
		in[filepath.Clean(repo.NarrativePath(st.ID))] = true
		ids[workitem.CanonicalID(st.ID)] = true
		for _, t := range workitem.Children(items, st.ID) {
			in[filepath.Clean(t.Path)] = true
			ids[workitem.CanonicalID(t.ID)] = true
		}
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
		for _, st := range sts {
			if c.Names(st.ID) {
				in[filepath.Clean(c.Path)] = true
			}
		}
	}
	return in, nil
}
