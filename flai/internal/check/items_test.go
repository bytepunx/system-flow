package check

import (
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/workflow"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// itemsProject is a project with a markdownlint configuration, an epic, and
// two stories under it, each with a forecast whose delivery is not a
// timestamp and a line that ends in a space: what Items is given and what it
// is not, both bad the same way.
func itemsProject(t *testing.T) (repo *workitem.Repo, given, other *workitem.Item) {
	t.Helper()
	repo = overlapProject(t)
	_ = os.WriteFile(filepath.Join(repo.Root, ".markdownlint.yaml"), []byte("default: true\nMD022:\n  lines_below: 0\nMD024:\n  siblings_only: true\nMD025:\n  front_matter_title: \"\"\n"), 0o644)
	_ = os.WriteFile(filepath.Join(repo.Root, ".gitignore"), []byte(".flai-cache/\n"), 0o644)
	epic := mustItem(t, repo, workitem.Epic, "E", "")
	spoil := func(it *workitem.Item) {
		data, err := os.ReadFile(it.Path)
		if err != nil {
			t.Fatal(err)
		}
		doc := strings.Replace(string(data), "\n---\n", "\nforecast:\n  duration: 1h\n  delivery: soon\n  by: t\n  at: 2026-09-01T12:00:00Z\n---\n", 1)
		doc = strings.Replace(doc, "## Notes\n", "## Notes\n\ntrailing \n", 1)
		if err := os.WriteFile(it.Path, []byte(doc), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	given = mustItem(t, repo, workitem.Story, "Given", epic.ID)
	other = mustItem(t, repo, workitem.Story, "Other", epic.ID)
	spoil(given)
	spoil(other)
	return repo, given, other
}

// listed is the project's items as Run lists them, and the one with id.
func listed(t *testing.T, repo *workitem.Repo, id string) ([]*workitem.Item, *workitem.Item) {
	t.Helper()
	all, err := repo.List(true)
	if err != nil {
		t.Fatal(err)
	}
	for _, it := range all {
		if it.ID == id {
			return all, it
		}
	}
	t.Fatalf("%s not listed", id)
	return nil, nil
}

// S-0599: Items reports a bad forecast and a markdown finding on the items it
// is given, and nothing on an item it is not given.
func TestItemsChecksOnlyTheItemsGiven(t *testing.T) {
	repo, given, other := itemsProject(t)
	all, it := listed(t, repo, given.ID)
	res := Items(repo, all, []*workitem.Item{it}, now)
	if res.Items != 1 {
		t.Errorf("Items = %d, want 1", res.Items)
	}
	var forecast, lint bool
	for _, f := range res.Findings {
		if strings.Contains(f.Path, other.ID) {
			t.Errorf("a finding on %s, which Items was not given: %+v", other.ID, f)
		}
		if !strings.Contains(f.Path, given.ID) {
			continue
		}
		forecast = forecast || (f.Rule == "item.front-matter" && f.Level == Error && strings.Contains(f.Message, "forecast.delivery"))
		lint = lint || (f.Rule == "markdown.MD009" && f.Level == Warning)
	}
	if !forecast || !lint {
		t.Errorf("want item.front-matter on forecast.delivery and markdown.MD009 on %s, got %+v", given.ID, res.Findings)
	}
	if res.Errors == 0 || res.Warnings == 0 {
		t.Errorf("Errors = %d, Warnings = %d, want both counted", res.Errors, res.Warnings)
	}
}

// S-0599: on the items given, Items finds what Run finds under the per-item
// rules and the markdown lint. Run's item rules that look across items,
// item.duplicate-id, story.component-tag, and story.unaccepted, are not
// Items' to give.
func TestItemsMatchesRunOnTheItemsGiven(t *testing.T) {
	repo, given, _ := itemsProject(t)
	all, it := listed(t, repo, given.ID)
	var epic *workitem.Item
	for _, x := range all {
		if x.ID == it.Parent {
			epic = x
		}
	}
	items := []*workitem.Item{epic, it}
	paths := map[string]bool{}
	for _, x := range items {
		rel, _ := filepath.Rel(repo.Root, x.Path)
		paths[rel] = true
	}
	across := []string{"item.duplicate-id", "story.component-tag", "story.unaccepted"}
	keep := func(fs []Finding) []string {
		var out []string
		for _, f := range fs {
			per := (strings.HasPrefix(f.Rule, "item.") || strings.HasPrefix(f.Rule, "story.") || strings.HasPrefix(f.Rule, "markdown.")) && !slices.Contains(across, f.Rule)
			if per && paths[f.Path] {
				out = append(out, f.Level+" "+f.Rule+" "+f.Path+" "+f.Message)
			}
		}
		slices.Sort(out)
		return out
	}
	run, err := Run(repo, now)
	if err != nil {
		t.Fatal(err)
	}
	want, got := keep(run.Findings), keep(Items(repo, all, items, now).Findings)
	if len(want) == 0 {
		t.Fatal("Run found nothing on the items given; the fixture should give it findings")
	}
	if !slices.Equal(got, want) {
		t.Errorf("Items findings differ from Run's on the items given:\ngot  %q\nwant %q", got, want)
	}
}

// I-0344: a task file another story's agent had written for a moment as
// PLACEHOLDER, with no front matter, made Run fail with the list's parse
// error, so flai check could not run and a close-out stopped at it. Run now
// reports each file that does not parse as one item.unreadable on its path,
// an error, and checks every other item as before: the forecasts broken on
// both stories are still found. A rule that looks the file's ID up passes
// over it, so a story under the epic that does not parse gets no
// item.parent-missing and a task whose after names the task that does not
// parse no task.after, and the markdown lint's trailing space on the epic's
// file is left out with its count.
func TestRunReportsAnItemFileThatDoesNotParseOnceAndChecksTheRest(t *testing.T) {
	repo, given, other := itemsProject(t)
	epic := mustItem(t, repo, workitem.Epic, "Lost", "")
	child := mustItem(t, repo, workitem.Story, "Child", epic.ID)
	broken := mustItem(t, repo, workitem.Task, "Broken", other.ID)
	waiter, err := repo.Create(workitem.NewOptions{Type: workitem.Task, Title: "Waiter", Parent: other.ID, After: []string{broken.ID}, Owner: "t", Now: now})
	if err != nil {
		t.Fatal(err)
	}
	listed(t, repo, broken.ID) // the folders are listed, and kept, before the damage
	for path, doc := range map[string]string{broken.Path: "PLACEHOLDER", epic.Path: "PLACEHOLDER \n"} {
		if err := os.WriteFile(path, []byte(doc), 0o644); err != nil {
			t.Fatal(err)
		}
		rel, _ := filepath.Rel(repo.Root, path)
		workitem.ItemChanged(repo.Root, rel)
	}

	res, err := Run(repo, now)
	if err != nil {
		t.Fatalf("flai check could not run: %v", err)
	}
	on := map[string][]Finding{}
	for _, f := range res.Findings {
		on[filepath.Base(f.Path)] = append(on[filepath.Base(f.Path)], f)
	}
	for _, path := range []string{broken.Path, epic.Path} {
		got := on[filepath.Base(path)]
		want := "does not parse: no front matter; flai leaves it out of every list of items until it parses: fix its front matter or remove it"
		if len(got) != 1 || got[0].Rule != "item.unreadable" || got[0].Level != Error || got[0].Line != 1 || got[0].Message != want {
			t.Errorf("findings on %s: %+v, want one item.unreadable %q", filepath.Base(path), got, want)
			continue
		}
		if rel, _ := filepath.Rel(repo.Root, path); got[0].Path != rel {
			t.Errorf("item.unreadable on %s, want the path relative as other item findings give it, %s", got[0].Path, rel)
		}
	}
	for _, it := range []*workitem.Item{given, other} {
		found := false
		for _, f := range on[filepath.Base(it.Path)] {
			found = found || (f.Rule == "item.front-matter" && strings.Contains(f.Message, "forecast.delivery"))
		}
		if !found {
			t.Errorf("%s's bad forecast is no longer found: %+v", it.ID, on[filepath.Base(it.Path)])
		}
	}
	for _, it := range []*workitem.Item{child, waiter} {
		for _, f := range on[filepath.Base(it.Path)] {
			if f.Rule == "item.parent-missing" || f.Rule == "task.after" {
				t.Errorf("%s gets %s for an ID whose file does not parse: %+v", it.ID, f.Rule, f)
			}
		}
	}
	if all, _ := repo.List(true); res.Items != len(all) || res.Items != 5 {
		t.Errorf("Items = %d, want the 5 items that parse", res.Items)
	}
	errs, warns := 0, 0
	for _, f := range res.Findings {
		if f.Level == Error {
			errs++
		} else {
			warns++
		}
	}
	if res.Errors != errs || res.Warnings != warns {
		t.Errorf("counts %d errors, %d warnings; the findings give %d, %d", res.Errors, res.Warnings, errs, warns)
	}
}

// I-0347: item.blocked-closed names the command that clears it, and once
// flai unblock has closed the interval at the item's close it is not found.
func TestBlockedClosedNamesFlaiUnblockAndUnblockClearsIt(t *testing.T) {
	repo := overlapProject(t)
	epic := mustItem(t, repo, workitem.Epic, "E", "")
	s := mustItem(t, repo, workitem.Story, "S", epic.ID)
	closedAt := now.Add(time.Hour).Format(workitem.TimeFormat)
	s.Blocked = []workitem.Block{{From: now.Format(workitem.TimeFormat), Reason: "pause for urgent fix"}}
	s.Transitions = append(s.Transitions, workitem.Transition{To: workitem.Cancelled, At: closedAt, By: "t"})
	s.Status, s.Updated = workitem.Cancelled, closedAt
	if err := repo.Save(s); err != nil {
		t.Fatal(err)
	}
	got := findings(t, repo, "item.blocked-closed")
	want := "cancelled item still has an open blocked interval; close it at the item's close with flai unblock " + s.ID
	if len(got) != 1 || got[0].Message != want {
		t.Fatalf("item.blocked-closed = %+v, want one saying %q", got, want)
	}
	if err := workitem.UnblockItem(s, now.Add(48*time.Hour)); err != nil {
		t.Fatal(err)
	}
	if err := repo.Save(s); err != nil {
		t.Fatal(err)
	}
	if s.Blocked[0].Until != closedAt {
		t.Errorf("until = %q, want the item's close %q", s.Blocked[0].Until, closedAt)
	}
	if got := findings(t, repo, "item.blocked-closed"); len(got) != 0 {
		t.Errorf("after flai unblock, still %+v", got)
	}
}

// kindsRenamedFlow is the baseline with review renamed waiting-on-operator,
// review kept among its aliases, so that items' history and the board's
// wip_limits that name review still read.
func kindsRenamedFlow() *workflow.Workflow {
	w := workflow.Baseline()
	rename := func(s string) string {
		if s == workitem.Review {
			return "waiting-on-operator"
		}
		return s
	}
	for i := range w.States {
		if w.States[i].Name == workitem.Review {
			w.States[i].Name, w.States[i].Aliases = "waiting-on-operator", []string{workitem.Review}
		}
	}
	for i := range w.Transitions {
		tr := &w.Transitions[i]
		for j := range tr.From {
			tr.From[j] = rename(tr.From[j])
		}
		tr.To = rename(tr.To)
	}
	for i := range w.Timers {
		for j := range w.Timers[i].Sum {
			w.Timers[i].Sum[j] = rename(w.Timers[i].Sum[j])
		}
	}
	return &w
}

// kindsVerifyFlow is the baseline with verify, a waiting state with a limit,
// added between in-progress and review: in progress moves to verify, which
// moves on to review or back to in-progress, and a parent's cancellation
// takes an item out of it as out of review.
func kindsVerifyFlow() *workflow.Workflow {
	w := workflow.Baseline()
	all := slices.Clone(workflow.Types)
	anyone := []string{workflow.ByAnyone}
	w.States = slices.Insert(w.States, 3, workflow.State{Name: "verify", Kind: workflow.KindWaiting, Title: "Verify", Types: all, Limit: true, Aliases: []string{}})
	for i := range w.Transitions {
		tr := &w.Transitions[i]
		if slices.Equal(tr.From, []string{workitem.InProgress}) && tr.To == workitem.Review {
			tr.To = "verify"
		}
		if tr.To == workitem.Cancelled && tr.Allows(workflow.ByCascade) {
			tr.From = append(tr.From, "verify")
		}
	}
	w.Transitions = append(w.Transitions,
		workflow.Transition{From: []string{"verify"}, To: workitem.Review, Types: all, By: anyone, Rules: []workflow.RuleRef{}, Effects: []string{}},
		workflow.Transition{From: []string{"verify"}, To: workitem.InProgress, Types: all, By: anyone, Rules: []workflow.RuleRef{}, Effects: []string{}})
	return &w
}

// kindsStory saves a story under epic, with its criteria ticked, in status,
// having moved through history an hour apart, and with a task unless bare.
func kindsStory(t *testing.T, repo *workitem.Repo, epic, title, status string, history []string, bare bool) *workitem.Item {
	t.Helper()
	s := mustItem(t, repo, workitem.Story, title, epic)
	if !bare {
		mustItem(t, repo, workitem.Task, "Task of "+title, s.ID)
	}
	s, err := repo.Get(s.ID)
	if err != nil {
		t.Fatal(err)
	}
	body := strings.Replace(s.Body, "## Acceptance criteria\n- [ ]\n", "## Acceptance criteria\n- [x] works\n", 1)
	if body == s.Body {
		t.Fatalf("%s's body has no empty criterion to tick:\n%s", s.ID, s.Body)
	}
	s.Body, s.Status = body, status
	for i, to := range history {
		at := now.Add(time.Duration(i+1) * time.Hour).Format(workitem.TimeFormat)
		s.Transitions = append(s.Transitions, workitem.Transition{To: to, At: at, By: "t"})
		s.Updated = at
	}
	if err := repo.Save(s); err != nil {
		t.Fatal(err)
	}
	return s
}

// kindsRules are the rules that read a state's kind or a transition's
// effect in the scenario of TestItemRulesReadKindsNotNames, and the
// workflow's own rules, which find nothing in a workflow it runs.
var kindsRules = []string{"story.tasks", "story.criteria", "story.draft", "story.unchecked", "story.unaccepted", "item.done-children", "item.sequence", "item.parent-cancelled", "narrative.missing", "board.wip-limit", "board.order"}

// S-0741: the item rules read a state's kind and a transition's effects from
// the workflow the repository runs, never a state's name. The same stories,
// in the states of the same kinds, get the same findings on the baseline, on
// the baseline with review renamed waiting-on-operator, and on the baseline
// with verify, a waiting state, added between in-progress and review: a
// story waiting with no task gets story.tasks, a story waiting or in progress
// has its narrative required, two stories waiting over the limit of the wait
// state are an advisory board.wip-limit, read by an alias in the renamed
// workflow, a waiting story in the pull order is board.order, a story done
// straight from in progress must go through the state the acceptance leaves,
// and is done with an open task and never accepted. Each finding prints the
// state's name as the workflow gives it. History reads by the workflow's
// transitions: review, an alias, follows in-progress in the renamed workflow,
// and with verify added a story moved from in-progress to review skipped it.
// openNarratives counts the narratives of the stories waiting and in
// progress as open, and not the done one's.
func TestItemRulesReadKindsNotNames(t *testing.T) {
	cases := []struct {
		name string
		flow *workflow.Workflow
		// wait is the waiting state the stories wait in, alias the name the
		// second one's history gives it, via the state the acceptance
		// leaves, and limit the wip_limits key the board gives wait.
		wait, alias, via, limit string
		// skipped is the item.sequence the story moved from in-progress to
		// review gets, empty when that move follows.
		skipped string
	}{
		{"the baseline", nil, workitem.Review, workitem.Review, workitem.Review, workitem.Review, ""},
		{"review renamed waiting-on-operator", kindsRenamedFlow(), "waiting-on-operator", workitem.Review, "waiting-on-operator", workitem.Review, ""},
		{"verify added before review", kindsVerifyFlow(), "verify", "verify", workitem.Review, "verify", "transitions[2]: review cannot follow in-progress"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			repo := overlapProject(t)
			repo.Workflow = tc.flow
			epic := mustItem(t, repo, workitem.Epic, "E", "")
			waiting := kindsStory(t, repo, epic.ID, "Waiting", tc.wait, []string{workitem.Ready, workitem.InProgress, tc.wait}, true)
			second := kindsStory(t, repo, epic.ID, "Second", tc.wait, []string{workitem.Ready, workitem.InProgress, tc.alias}, false)
			straight := kindsStory(t, repo, epic.ID, "Straight", workitem.Done, []string{workitem.Ready, workitem.InProgress, workitem.Done}, false)
			back := kindsStory(t, repo, epic.ID, "Back", workitem.InProgress, []string{workitem.Ready, workitem.InProgress, workitem.Review, workitem.InProgress}, false)
			writeFile(t, filepath.Join(repo.KanbanDir(), workitem.BoardFile), "---\ntitle: Board\nupdated: 2026-09-01\nstatus: active\nwip_limits:\n  "+tc.limit+": 1\norder:\n  - "+second.ID+"\n---\n\n# Board\n")

			res, err := Run(repo, now)
			if err != nil {
				t.Fatal(err)
			}
			var got []string
			for _, f := range res.Findings {
				if !slices.Contains(kindsRules, f.Rule) && !strings.HasPrefix(f.Rule, "workflow.") {
					continue
				}
				on := fileID(f.Path)
				if on == "" {
					on = filepath.Base(f.Path)
				}
				line := f.Level + " " + f.Rule + " " + on + ": " + f.Message
				if f.advisory {
					line += " (advisory)"
				}
				got = append(got, line)
			}
			want := []string{
				"warning board.order board.md: order lists " + second.ID + " which is " + tc.wait + "; only ready and backlog stories belong in the pull order",
				"warning board.wip-limit board.md: 2 stories in " + tc.wait + ", limit 1 (advisory)",
				"error item.done-children " + straight.ID + ": " + straight.ID + " is done but T-0002 is backlog",
				"error item.sequence " + straight.ID + ": transitions[2]: a story must go through " + tc.via + " before done",
				"error narrative.missing " + waiting.ID + ": " + waiting.ID + " is " + tc.wait + " but has no narrative; run flai stream open " + waiting.ID,
				"error narrative.missing " + second.ID + ": " + second.ID + " is " + tc.wait + " but has no narrative; run flai stream open " + second.ID,
				"error narrative.missing " + back.ID + ": " + back.ID + " is in-progress but has no narrative; run flai stream open " + back.ID,
				"error story.tasks " + waiting.ID + ": a " + tc.wait + " story needs at least one task",
				"warning story.unaccepted " + straight.ID + ": " + straight.ID + " is done but was never accepted (not archived, not released); run flai accept " + straight.ID,
			}
			if tc.skipped != "" {
				want = append(want, "error item.sequence "+back.ID+": "+tc.skipped)
			}
			slices.Sort(got)
			slices.Sort(want)
			if !slices.Equal(got, want) {
				t.Errorf("findings:\n%s\nwant:\n%s", strings.Join(got, "\n"), strings.Join(want, "\n"))
			}

			open, err := openNarratives(repo, nil)
			if err != nil {
				t.Fatal(err)
			}
			for _, s := range []*workitem.Item{waiting, second, back, straight} {
				if want := s != straight; open[filepath.Clean(repo.NarrativePath(s.ID))] != want {
					t.Errorf("openNarratives counts %s's narrative open: %v, want %v", s.ID, !want, want)
				}
			}
		})
	}
}
