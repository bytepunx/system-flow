package itemedit

import (
	"errors"
	"fmt"
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
// stamped by by at now with the edit notice flai touches leaves, and returns
// the paths added, none when they cover every path.
func Widen(repo *workitem.Repo, id string, paths []string, by string, now time.Time) ([]string, error) {
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
	it.Updated = now.UTC().Format(workitem.TimeFormat)
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
		added, err := Widen(o.Repo, o.Task, paths, o.By, o.Now())
		if err != nil {
			return w, err
		}
		w.Task = added
	}
	added, err := Widen(o.Repo, o.Story, paths, o.By, o.Now())
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

// within says whether path is one of entries or lies below one.
func within(path string, entries []string) bool {
	for _, e := range entries {
		if path == e || strings.HasPrefix(path, e+"/") {
			return true
		}
	}
	return false
}
