package check

import (
	"path/filepath"
	"testing"

	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// S-0200: an open epic its stories put further along the board than its
// status is warned about on its status line, advisory so that --strict passes
// over it; an epic that matches its stories, one ahead of them, one with no
// story that counts, and a done or cancelled one are not.
func TestAnEpicThatLagsItsStoriesIsWarned(t *testing.T) {
	repo := planningProject(t, "")
	set := func(it *workitem.Item, status string) {
		t.Helper()
		it.Status = status
		if err := repo.Save(it); err != nil {
			t.Fatal(err)
		}
	}
	epic := func(status string, stories ...string) *workitem.Item {
		t.Helper()
		e := mustItem(t, repo, workitem.Epic, "Epic "+status, "")
		for _, st := range stories {
			set(mustItem(t, repo, workitem.Story, "Story "+st, e.ID), st)
		}
		set(e, status)
		return e
	}
	started := epic(workitem.Backlog, workitem.InProgress, workitem.Backlog)
	reviewed := epic(workitem.Backlog, workitem.Review, workitem.Done, workitem.Cancelled)
	finished := epic(workitem.InProgress, workitem.Done, workitem.Done, workitem.Cancelled)
	epic(workitem.Ready, workitem.Ready, workitem.Backlog)
	epic(workitem.Review, workitem.InProgress)
	epic(workitem.Backlog, workitem.Cancelled)
	epic(workitem.Backlog)
	epic(workitem.Done, workitem.Done)
	epic(workitem.Cancelled, workitem.InProgress)

	// lags runs the rule alone, so the stories' other findings (no tasks, no
	// narrative) do not decide OK.
	items, err := repo.List(true)
	if err != nil {
		t.Fatal(err)
	}
	c := &checker{repo: repo, now: now, items: items, byID: map[string]*workitem.Item{}, res: &Result{}}
	c.epicLags()
	why := " (an epic follows its stories since S-0200, and this one was moved before that, or by hand)"
	want := map[string]string{
		started.ID:  started.ID + " is backlog but its stories put it in in-progress; catch it up with flai move " + started.ID + " ready, then in-progress" + why,
		reviewed.ID: reviewed.ID + " is backlog but its stories put it in review; catch it up with flai move " + reviewed.ID + " ready, then in-progress, then review" + why,
		finished.ID: finished.ID + " is in-progress but its stories are all done; accept it with flai accept " + finished.ID + ", which moves it to done and archives it" + why,
	}
	byPath := map[string]*workitem.Item{}
	for _, e := range []*workitem.Item{started, reviewed, finished} {
		p, _ := filepath.Rel(repo.Root, e.Path)
		byPath[p] = e
	}
	if len(c.res.Findings) != len(want) {
		t.Errorf("want %d findings, on %s, %s, and %s: %+v", len(want), started.ID, reviewed.ID, finished.ID, c.res.Findings)
	}
	for _, f := range c.res.Findings {
		e, ok := byPath[f.Path]
		if !ok {
			t.Errorf("an epic that does not lag is warned: %+v", f)
			continue
		}
		if f.Level != Warning || f.Rule != "epic.lags-stories" || f.Line != keyLine(e.Path, "status") || f.Line == 1 || f.Message != want[e.ID] {
			t.Errorf("got %+v, want a warning on %s's status line saying %q", f, e.ID, want[e.ID])
		}
	}
	if c.res.Warnings != len(want) || c.res.Advisory != len(want) || !c.res.OK(true) {
		t.Errorf("the warnings must be advisory and pass --strict: %+v", c.res)
	}
	full, err := Run(repo, now)
	if err != nil {
		t.Fatal(err)
	}
	if full.Advisory != len(want) {
		t.Errorf("Run: advisory %d, want %d", full.Advisory, len(want))
	}
}
