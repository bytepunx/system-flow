package check

import (
	"slices"
	"strings"
	"testing"

	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// S-0551: an open story gets one advisory wip.after-reason on its after:
// for each story there that its Notes do not name as a whole word, in any
// padding; a done story and a task's after: are passed over.
func TestAfterReason(t *testing.T) {
	repo := overlapProject(t)
	e := mustItem(t, repo, workitem.Epic, "E", "")
	reads := mustItem(t, repo, workitem.Story, "Reads", e.ID)
	waits := mustItem(t, repo, workitem.Story, "Waits", e.ID)
	tasks := []*workitem.Item{mustItem(t, repo, workitem.Task, "First", reads.ID), mustItem(t, repo, workitem.Task, "Second", reads.ID)}
	tasks[1].After = []string{tasks[0].ID}
	waits.Status, waits.After = workitem.Ready, []string{reads.ID}
	save := func(it *workitem.Item) {
		t.Helper()
		if err := repo.Save(it); err != nil {
			t.Fatal(err)
		}
	}
	save(tasks[1])
	reasons := func(notes string) []Finding {
		t.Helper()
		waits.Body = "\n## Goal\n\nWait on S-0001.\n\n## Notes\n\n" + notes + "\n\n## Tasks\n\nS-0001 again.\n"
		save(waits)
		res, err := Run(repo, now)
		if err != nil {
			t.Fatal(err)
		}
		var out []Finding
		for _, f := range res.Findings {
			if f.Rule == "wip.after-reason" {
				out = append(out, f)
			}
		}
		return out
	}

	for _, notes := range []string{"Reads the parser S-0001 adds.", "### Planning\n\nReads [S-001](S-001-reads.md)'s parser.", "Reads s1's parser."} {
		if got := reasons(notes); len(got) != 0 {
			t.Errorf("Notes %q name S-0001, yet wip.after-reason = %+v", notes, got)
		}
	}

	want := "S-0002 names S-0001 in after with no reason in its Notes: say there what S-0002 reads from S-0001, or drop it from after (the after rule)"
	for _, notes := range []string{"No reason.", "Reads S-0010 and S-00011.", "Reads S-0001a."} {
		got := reasons(notes)
		if len(got) != 1 || got[0].Message != want || got[0].Level != Warning || !got[0].advisory ||
			!slices.Equal(got[0].Stories, []string{"S-0002"}) || !strings.HasSuffix(got[0].Path, "S-0002-waits.md") || got[0].Line != keyLine(waits.Path, "after") {
			t.Errorf("Notes %q do not name S-0001 as a whole word: wip.after-reason = %+v\nwant one advisory on S-0002: %s", notes, got, want)
		}
	}

	waits.Status = workitem.Done
	if got := reasons("No reason."); len(got) != 0 {
		t.Errorf("a done story is passed over: %+v", got)
	}
}

// S-0551: scoped to a story, a wip.after-reason on another story is left out,
// even when the story's branch changes that story's file, and the story's own
// is kept, inside it and advisory, as a wip.after-hub is.
func TestScopedCheckLeavesOutAnotherStorysAfterReason(t *testing.T) {
	root := scopeProject(t)
	edit(t, root, "wip/kanban/board.md", "in-progress: 2\n", "in-progress: 5\n")
	addStory(t, root, "S-006", "E-001", "")
	addStory(t, root, "S-007", "E-001", "after: [S-006]\n")
	edit(t, root, "wip/kanban/stories/S-004-four.md", "tags: []\n", "tags: []\nafter: [S-006]\n")
	reasons := func(res *Result) map[string]Finding {
		out := map[string]Finding{}
		for _, f := range res.Findings {
			if f.Rule == "wip.after-reason" {
				out[strings.Join(f.Stories, ",")] = f
			}
		}
		return out
	}
	whole, scoped := runScoped(t, root, []string{"wip/kanban/stories/S-007-six.md"})
	if all := reasons(whole); len(all) != 2 || !all["S-004"].advisory || !all["S-007"].advisory {
		t.Fatalf("unscoped, the run should report an advisory wip.after-reason on S-004 and S-007: %+v", whole.Findings)
	}
	kept := reasons(scoped)
	own, ok := kept["S-004"]
	if len(kept) != 1 || !ok {
		t.Fatalf("scoped to S-004, only its own wip.after-reason should be kept: %+v", kept)
	}
	if own.Outside || !own.advisory {
		t.Errorf("S-004's own wip.after-reason should be inside it and advisory: %+v", own)
	}
}
