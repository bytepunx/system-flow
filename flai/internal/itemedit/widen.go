package itemedit

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// A command that commits files on a story's branch widens the touches with
// them, as flai touches --add would: flai task done with its commit's paths
// (ADR-0107), and the commands run with --commit in a story's worktree with
// the files they wrote (S-0275).

// Widen adds to item id's touches the paths among paths they do not cover,
// stamped by by at now, read once the item is (I-0183), with the edit notice
// flai touches leaves, and returns the paths added, none when they cover
// every path.
func Widen(repo *workitem.Repo, id string, paths []string, by string, now func() time.Time) ([]string, error) {
	it, err := repo.Get(id)
	if err != nil {
		return nil, fmt.Errorf("read %s to widen its touches: %w", id, err)
	}
	covered := workitem.NewHolds(nil, repo.Manifest.Projects).Claim(&workitem.Item{Touches: it.Touches})
	var missing []string
	for _, p := range paths {
		if !within(p, covered) {
			missing = append(missing, p)
		}
	}
	add, err := workitem.CleanTouches(missing)
	if err != nil {
		return nil, fmt.Errorf("widen %s's touches: %w; add the paths by hand with flai touches %s --add", id, err, id)
	}
	if len(add) == 0 {
		return []string{}, nil
	}
	it.Touches = append(slices.Clone(it.Touches), add...)
	it.Updated = now().UTC().Format(workitem.TimeFormat)
	if err := repo.Save(it); err != nil {
		return nil, fmt.Errorf("save %s's widened touches: %w", id, err)
	}
	Record(repo, Notice{At: it.Updated, By: by, ID: it.ID, Type: it.Type, Title: it.Title, Changed: []string{"touches"}})
	return add, nil
}

// WidenOptions says whose touches WidenStory widens, with what, as whom.
type WidenOptions struct {
	Repo  *workitem.Repo
	Story string
	// Task, when set, is a task of the story whose touches are widened too,
	// before the story's.
	Task string
	// Paths are relative to the project's root; those under the wip folder,
	// which is flai's, are left out.
	Paths []string
	By    string
	Now   func() time.Time // time.Now when nil
}

// Widened is what WidenStory added to the task's touches and to the story's,
// and the stories in progress the story's claim grew into.
type Widened struct {
	Task     []string      `json:"task"`
	Story    []string      `json:"story"`
	Overlaps []Overlapping `json:"overlaps"`
	// Untold is why the overlaps could not be found or recorded, nil when
	// they were. It is advisory, as flai touches has it: the touches stand.
	Untold error `json:"-"`
}

// WidenStory adds the paths outside the wip folder to the task's touches,
// when there is a task, and to the story's, each where it does not cover
// them, and reports the stories in progress the story's claim grew into.
func WidenStory(o WidenOptions) (Widened, error) {
	w := Widened{Task: []string{}, Story: []string{}, Overlaps: []Overlapping{}}
	if o.Repo == nil || o.Story == "" {
		return w, errors.New("itemedit.WidenStory needs a project and a story")
	}
	if o.Now == nil {
		o.Now = time.Now
	}
	wip := strings.TrimSuffix(o.Repo.Manifest.Layout["wip"], "/")
	if wip == "" {
		wip = "wip"
	}
	var paths []string
	for _, p := range o.Paths {
		if !within(p, []string{wip}) {
			paths = append(paths, p)
		}
	}
	watch := WatchClaim(o.Repo, o.Story)
	if o.Task != "" {
		added, err := Widen(o.Repo, o.Task, paths, o.By, o.Now)
		if err != nil {
			return w, err
		}
		w.Task = added
	}
	added, err := Widen(o.Repo, o.Story, paths, o.By, o.Now)
	if err != nil {
		return w, err
	}
	w.Story = added
	if len(w.Task)+len(w.Story) == 0 {
		return w, nil
	}
	told, err := watch.Grown(o.By, o.Now())
	w.Untold = err
	if told != nil {
		w.Overlaps = told
	}
	return w, nil
}

// NarrowOptions says whose touches NarrowStory narrows, against what, as
// whom (ADR-0284).
type NarrowOptions struct {
	Repo  *workitem.Repo
	Story string
	// Task is the task of the story being closed, whose touches are narrowed
	// first.
	Task string
	// Tree is the folder a touch is judged a folder in: the story's worktree.
	Tree string
	// Changed are the paths the story's branch changes against its merge base
	// with the main branch, relative to the project's root.
	Changed []string
	By      string
	Now     func() time.Time // time.Now when nil
}

// Narrowed is what NarrowStory removed from the task's touches and from the
// story's.
type Narrowed struct {
	Task  []string `json:"task"`
	Story []string `json:"story"`
}

// NarrowStory removes from the task's touches each one that names a file the
// story's branch does not change, never a sub-project's name or tag, a folder
// in the tree, or a path under the wip folder, and removes each such touch
// from the story's touches too when the story names it and no other open
// task of the story does (ADR-0284). Each item it changes is stamped by by at
// now, with the edit notice flai touches leaves. A claim that shrinks tells
// no one.
func NarrowStory(o NarrowOptions) (Narrowed, error) {
	n := Narrowed{Task: []string{}, Story: []string{}}
	if o.Repo == nil || o.Story == "" || o.Task == "" {
		return n, errors.New("itemedit.NarrowStory needs a project, a story, and a task")
	}
	if o.Now == nil {
		o.Now = time.Now
	}
	task, err := o.Repo.Get(o.Task)
	if err != nil {
		return n, fmt.Errorf("read %s to narrow its touches: %w", o.Task, err)
	}
	wip := strings.TrimSuffix(o.Repo.Manifest.Layout["wip"], "/")
	if wip == "" {
		wip = "wip"
	}
	var keep []string
	for _, e := range task.Touches {
		if predicted(o, wip, e) {
			n.Task = append(n.Task, e)
		} else {
			keep = append(keep, e)
		}
	}
	if len(n.Task) == 0 {
		return n, nil
	}
	if err := saveTouches(o.Repo, task, keep, o.By, o.Now); err != nil {
		return n, err
	}
	story, err := o.Repo.Get(o.Story)
	if err != nil {
		return n, fmt.Errorf("read %s to narrow its touches: %w", o.Story, err)
	}
	tasks, err := o.Repo.ChildrenOf(story.ID, false)
	if err != nil {
		return n, fmt.Errorf("read the tasks of %s to narrow its touches: %w", story.ID, err)
	}
	named := map[string]bool{}
	for _, t := range tasks {
		if t.Type == workitem.Task && !t.Closed() && workitem.CanonicalID(t.ID) != workitem.CanonicalID(task.ID) {
			for _, e := range t.Touches {
				named[cleanTouch(e)] = true
			}
		}
	}
	drop := map[string]bool{}
	for _, e := range n.Task {
		if !named[cleanTouch(e)] {
			drop[cleanTouch(e)] = true
		}
	}
	keep = nil
	for _, e := range story.Touches {
		if drop[cleanTouch(e)] {
			n.Story = append(n.Story, e)
		} else {
			keep = append(keep, e)
		}
	}
	if len(n.Story) == 0 {
		return n, nil
	}
	return n, saveTouches(o.Repo, story, keep, o.By, o.Now)
}

// predicted says whether touch e names a file the story's branch does not
// change: neither a sub-project's name or tag, nor a path under the wip
// folder, nor a folder in the tree, nor a path the branch changes or one
// that holds a path it changes, such as a folder whose files it deleted.
func predicted(o NarrowOptions, wip, e string) bool {
	p := cleanTouch(e)
	if p == "" || strings.HasPrefix(p, "/") || slices.Contains(strings.Split(p, "/"), "..") || within(p, []string{wip}) {
		return false
	}
	for _, pr := range o.Repo.Manifest.Projects {
		if p == pr.Name || slices.Contains(pr.Tags, p) {
			return false
		}
	}
	for _, c := range o.Changed {
		if within(c, []string{p}) {
			return false // the branch changes it, or a file below it
		}
	}
	if fi, err := os.Stat(filepath.Join(o.Tree, filepath.FromSlash(p))); err == nil && fi.IsDir() {
		return false
	}
	return true
}

// cleanTouch is a touches entry as a path: no leading ./, no trailing /.
func cleanTouch(e string) string {
	return strings.TrimSuffix(strings.TrimPrefix(strings.TrimSpace(e), "./"), "/")
}

// saveTouches saves it with touches, stamped by by at now, with the edit
// notice flai touches leaves.
func saveTouches(repo *workitem.Repo, it *workitem.Item, touches []string, by string, now func() time.Time) error {
	if touches == nil {
		touches = []string{}
	}
	it.Touches = touches
	it.Updated = now().UTC().Format(workitem.TimeFormat)
	if err := repo.Save(it); err != nil {
		return fmt.Errorf("save %s's narrowed touches: %w", it.ID, err)
	}
	Record(repo, Notice{At: it.Updated, By: by, ID: it.ID, Type: it.Type, Title: it.Title, Changed: []string{"touches"}})
	return nil
}

// within says whether path is one of entries or lies below one.
func within(path string, entries []string) bool {
	for _, e := range entries {
		if path == e || strings.HasPrefix(path, e+"/") {
			return true
		}
	}
	return false
}
