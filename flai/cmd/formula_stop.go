package cmd

import (
	"cmp"
	"fmt"
	"io"
	"os"
	"slices"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"github.com/bytepunx/system-flow/flai/internal/formulas"
	"github.com/bytepunx/system-flow/flai/internal/guard"
	"github.com/bytepunx/system-flow/flai/internal/itemedit"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// stopOutputLines are the last lines of a stopping step's output that its
// task's Notes keep (ADR-0271 §2).
const stopOutputLines = 40

func newFormulaStopCmd(a *app) *cobra.Command {
	var step, output string
	var exit int
	c := &cobra.Command{
		Use:   "stop <S-nnnn> --step <id> --exit <n>",
		Short: "Record that a step of a formula's instance stopped its run cleanly with one of its stop exits, and skip the steps after it (ADR-0271)",
		Long: `Record that a step of a formula's instance exited with one of the exits its
stop lists, which stop the run cleanly rather than fail the step (ADR-0271).
The step's output comes from --output, a file, or standard input.

In one change it moves the step's task to done, unblocking it first when it
is blocked, with "Stopped the run with exit <n>:" and the output's last 40
lines in its Notes; moves every other step task of the instance not done to
cancelled with the reason "step <id> stopped the run with exit <n>"; ticks
the step's acceptance criterion, and ticks each cancelled step's with
"(skipped: step <id> stopped the run with exit <n>)" after its text; and
writes one line naming the tasks it skipped to the story's Notes and to its
narrative's log. It commits nothing: the instance's agent then closes the
story out and moves it to review, asking nothing.

It refuses, writing nothing, a story that is no formula's instance or is not
in progress; an instance whose formula file is gone or has a problem flai
check reports; a step the formula does not have; an exit of 0 or one the
step's stop does not list; a step whose task is not in progress, or waits
for a task not done; and an instance with a task done in a layer after the
step's, since the run went on past it. An agent flai serve started runs it
only for the instance it works. --json gives {story, step, exit, done,
cancelled, ticked, skipped, logged}.`,
		Example: `  flai release --evaluate --require-met > /tmp/evaluate.txt; flai formula stop S-0585 --step evaluate --exit 2 --output /tmp/evaluate.txt
  flai release --evaluate --require-met | flai formula stop S-0585 --step evaluate --exit 2`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if !cmd.Flags().Changed("exit") {
				return fmt.Errorf("give --exit, the exit the step stopped with")
			}
			var text []byte
			var err error
			if output != "" {
				text, err = os.ReadFile(output)
			} else {
				text, err = io.ReadAll(cmd.InOrStdin())
			}
			if err != nil {
				return fmt.Errorf("reading the step's output: %w", err)
			}
			return a.formulaStop(args[0], step, exit, string(text))
		},
	}
	c.Flags().StringVar(&step, "step", "", "the id of the step that stopped the run, as the formula names it")
	c.Flags().IntVar(&exit, "exit", 0, "the exit the step stopped with, one its stop lists")
	c.Flags().StringVar(&output, "output", "", "a file holding the step's output (default: standard input)")
	_ = c.MarkFlagRequired("step")
	return c
}

// formulaStopResult is what flai formula stop wrote: the stopping step's
// task, done; the step tasks cancelled; the criteria ticked, the skipped ones
// among them; and whether the narrative's log took the line.
type formulaStopResult struct {
	Story     string   `json:"story"`
	Step      string   `json:"step"`
	Exit      int      `json:"exit"`
	Done      string   `json:"done"`
	Cancelled []string `json:"cancelled"`
	Ticked    []int    `json:"ticked"`
	Skipped   []int    `json:"skipped"`
	Logged    bool     `json:"logged"`
}

// stoppedReason is the reason each step task a stop cancels is moved to
// cancelled with, and the words its criterion gains (ADR-0271 §2), which the
// orchestrator's acceptance reads (§4).
func stoppedReason(step string, exit int) string {
	return fmt.Sprintf("step %s stopped the run with exit %d", step, exit)
}

// formulaStop records that step of the instance id stopped its run with exit
// (ADR-0271 §2): it checks everything first, then moves the step's task to
// done with the output's last lines in its Notes, cancels the other step
// tasks not done, ticks the criteria, notes the stop on the story, and logs
// it to the narrative when there is one.
func (a *app) formulaStop(id, step string, exit int, output string) error {
	repo, err := a.project()
	if err != nil {
		return err
	}
	story, err := repo.Get(id)
	if err != nil {
		return err
	}
	if story.Type != workitem.Story || story.Formula == "" {
		return fmt.Errorf("rule: %s is no formula's instance (it names no formula), so it has no step to stop: a step that exits other than 0 fails it, as for any task", story.ID)
	}
	if story.Status != workitem.InProgress {
		return fmt.Errorf("rule: %s is %s, not in progress: a step stops the run of an instance being worked", story.ID, story.Status)
	}
	if err := servedStopsOwn(story.ID); err != nil {
		return err
	}
	_, f, err := a.formula(story.Formula)
	if err != nil {
		return fmt.Errorf("rule: %s is an instance of the formula %s, which cannot be read: %w", story.ID, story.Formula, err)
	}
	if problems := formulas.Validate([]formulas.Formula{f}, formulaNames(repo)); len(problems) > 0 {
		return fmt.Errorf("rule: the formula %s has %s, so %s's run was not stopped; fix the file, as flai check reports it, and run it again: %s", f.Name, plural(len(problems), "problem"), story.ID, problems[0])
	}
	at := slices.IndexFunc(f.Steps, func(s formulas.Step) bool { return s.ID == step })
	if at < 0 {
		ids := make([]string, len(f.Steps))
		for i, s := range f.Steps {
			ids[i] = s.ID
		}
		return fmt.Errorf("rule: the formula %s has no step %q: give one of %s", f.Name, step, strings.Join(ids, ", "))
	}
	s := f.Steps[at]
	if !s.Stops(exit) {
		if len(s.Stop) == 0 {
			return fmt.Errorf("rule: step %s of the formula %s names no stop, so exit %d fails it: block its task and ask, as for a failed step", step, f.Name, exit)
		}
		return fmt.Errorf("rule: step %s of the formula %s stops the run only with exit %s, not %d: any other exit fails it, so block its task and ask, as for a failed step", step, f.Name, joinExits(s.Stop, " or "), exit)
	}
	all, err := repo.List(false)
	if err != nil {
		return err
	}
	tasks, err := stepTasks(f, workitem.Children(all, story.ID), story.ID)
	if err != nil {
		return err
	}
	criteria := itemedit.Criteria(story.Body)
	if len(criteria) != len(f.Steps) {
		return fmt.Errorf("rule: %s has %d acceptance criteria, not one per step of the formula %s (%d), so the steps' criteria cannot be told apart; put them back as flai formula run wrote them", story.ID, len(criteria), f.Name, len(f.Steps))
	}
	task := tasks[step]
	if task.Status != workitem.InProgress {
		return fmt.Errorf("rule: %s, the task of step %s, is %s, not in progress: stop the run from the step being run", task.ID, step, task.Status)
	}
	byID := map[string]*workitem.Item{}
	for _, t := range tasks {
		byID[t.ID] = t
	}
	for _, after := range task.After {
		if t := byID[workitem.CanonicalID(after)]; t == nil || t.Status != workitem.Done {
			return fmt.Errorf("rule: %s, the task of step %s, waits for %s, which is not done: run the steps in their layers", task.ID, step, workitem.CanonicalID(after))
		}
	}
	layer := map[string]int{}
	for i, l := range formulas.Layers(f) {
		for _, sid := range l {
			layer[sid] = i
		}
	}
	var skip []formulas.Step
	for _, other := range f.Steps {
		t := tasks[other.ID]
		switch {
		case other.ID == step:
		case t.Status == workitem.Done && layer[other.ID] > layer[step]:
			return fmt.Errorf("rule: %s, the task of step %s, is done, and its layer comes after step %s's: the run went on past step %s, so it cannot stop there", t.ID, other.ID, step, step)
		case t.Status == workitem.Done, t.Status == workitem.Cancelled:
		case !workitem.Follows(workitem.Task, t.Status, workitem.Cancelled):
			return fmt.Errorf("rule: %s, the task of step %s, is %s and cannot be cancelled", t.ID, other.ID, t.Status)
		default:
			skip = append(skip, other)
		}
	}
	logged := false
	if _, err := workitem.ReadNarrative(repo.NarrativePath(story.ID)); err == nil {
		logged = true
	}

	now, by := a.now(), a.movedBy("")
	reason := stoppedReason(step, exit)
	if task.IsBlocked() {
		if err := workitem.UnblockItem(task, now); err != nil {
			return err
		}
	}
	task.Body = appendNotes(task.Body, fmt.Sprintf("Stopped the run with exit %d:\n\n%s", exit, fenced(lastLines(output, stopOutputLines))))
	if _, err := repo.TransitionAll(task, workitem.Done, by, "", now, false); err != nil {
		return fmt.Errorf("moving %s, the task of step %s, to done: %w", task.ID, step, err)
	}
	res := formulaStopResult{Story: story.ID, Step: step, Exit: exit, Done: task.ID, Cancelled: []string{}, Ticked: []int{at + 1}, Skipped: []int{}}
	suffixes := map[int]string{}
	for _, other := range skip {
		t := tasks[other.ID]
		if _, err := repo.TransitionAll(t, workitem.Cancelled, by, reason, now, false); err != nil {
			return fmt.Errorf("moving %s, the task of step %s, to cancelled, after %s went to done: %w; cancel the rest by hand with the reason %q", t.ID, other.ID, task.ID, err, reason)
		}
		n := slices.IndexFunc(f.Steps, func(s formulas.Step) bool { return s.ID == other.ID }) + 1
		res.Cancelled = append(res.Cancelled, t.ID)
		res.Ticked = append(res.Ticked, n)
		res.Skipped = append(res.Skipped, n)
		suffixes[n] = " (skipped: " + reason + ")"
	}
	slices.Sort(res.Ticked)
	slices.Sort(res.Skipped)
	line := fmt.Sprintf("Step %s stopped the run with exit %d; skipped %s.", step, exit, cmp.Or(strings.Join(res.Cancelled, ", "), "no step"))
	// the moves wrote the story's file only through its tasks: read it again
	story, err = repo.Get(story.ID)
	if err != nil {
		return err
	}
	body, err := itemedit.Tick(story.Body, res.Ticked, nil)
	if err != nil {
		return err
	}
	story.Body = appendNotes(skipCriteria(body, suffixes), line)
	story.Updated = now.UTC().Format(workitem.TimeFormat)
	if err := repo.Save(story); err != nil {
		return fmt.Errorf("writing %s's criteria and notes, after its step tasks moved: %w", story.ID, err)
	}
	if logged {
		agent, session := agentIdentity()
		if _, err := repo.LogStream(story.ID, line, workitem.StreamOptions{Agent: agent, Session: session, Now: now}); err != nil {
			return fmt.Errorf("logging the stop to %s's narrative, after everything else was written: %w", story.ID, err)
		}
	}
	if err := a.refreshIndex(repo); err != nil {
		return err
	}
	res.Logged = logged
	if a.jsonOut {
		return a.printJSON(res)
	}
	fmt.Fprintf(a.out, "%s %s\n  %s → done (step %s, exit %d)\n", story.ID, line, task.ID, step, exit)
	for _, t := range res.Cancelled {
		fmt.Fprintf(a.out, "  %s → cancelled: %s\n", t, reason)
	}
	if !logged {
		fmt.Fprintf(a.out, "  not logged: %s has no narrative\n", story.ID)
	}
	return nil
}

// servedStopsOwn says why the session may not stop the run of the instance
// id, or nil when it may: an agent flai serve started stops the instance it
// works alone, and a strategic role none, as flai guard holds it to (ADR-0271
// §2). The operator's own session stops any.
func servedStopsOwn(id string) error {
	if os.Getenv(guard.StartedByEnv) != guard.StartedByServe {
		return nil
	}
	if role := os.Getenv("FLAI_ROLE"); role != "" {
		return fmt.Errorf("rule: the %s role does not stop a formula's run: the agent of the instance %s, or the operator, does (ADR-0271)", role, id)
	}
	own := recordingStory("", os.Getenv, func() string { return "" })
	if own == "" || workitem.CanonicalID(own) != id {
		return fmt.Errorf("rule: an agent flai serve started stops the run of the instance it works alone, and this one works %s, not %s: the agent of %s, or the operator, stops it (ADR-0271)", cmp.Or(own, "no story"), id, id)
	}
	return nil
}

// stepTasks are the instance's tasks by the id of the step each is for:
// flai formula run makes them in creationOrder, so they are the instance's
// tasks in ID order, one per step. Anything else is refused.
func stepTasks(f formulas.Formula, children []*workitem.Item, story string) (map[string]*workitem.Item, error) {
	order := creationOrder(f)
	if len(children) != len(order) {
		return nil, fmt.Errorf("rule: %s has %d tasks, not one per step of the formula %s (%d), so its steps' tasks cannot be told apart", story, len(children), f.Name, len(order))
	}
	sorted := slices.Clone(children)
	slices.SortFunc(sorted, func(x, y *workitem.Item) int { return cmp.Compare(idNumber(x.ID), idNumber(y.ID)) })
	out := map[string]*workitem.Item{}
	for i, s := range order {
		if !strings.Contains(s.Title, "{{") && workitem.CleanTitle(sorted[i].Title) != workitem.CleanTitle(s.Title) {
			return nil, fmt.Errorf("rule: %s, %s's task %d in ID order, is titled %q, not as step %s of the formula %s is (%q), so its steps' tasks cannot be told apart", sorted[i].ID, story, i+1, sorted[i].Title, s.ID, f.Name, s.Title)
		}
		out[s.ID] = sorted[i]
	}
	return out, nil
}

// idNumber is the number of an item ID such as T-0042, or 0 for none.
func idNumber(id string) int {
	_, num, _ := strings.Cut(id, "-")
	n, _ := strconv.Atoi(num)
	return n
}

// skipCriteria is body with the text of each acceptance criterion numbered in
// suffixes followed by its suffix, and every other byte as it was.
func skipCriteria(body string, suffixes map[int]string) string {
	if len(suffixes) == 0 {
		return body
	}
	lines := strings.Split(body, "\n")
	in, n := false, 0
	for i, l := range lines {
		switch {
		case strings.HasPrefix(l, "## "):
			in = strings.TrimSpace(l) == "## Acceptance criteria"
		case in && isCriterionBox(l):
			n++
			if s, ok := suffixes[n]; ok {
				lines[i] = strings.TrimRight(l, " \t") + s
			}
		}
	}
	return strings.Join(lines, "\n")
}

// isCriterionBox says whether line is a checkbox of the acceptance criteria,
// "- [ ] text" or "- [x] text", indented or not, as itemedit counts them.
func isCriterionBox(line string) bool {
	t := strings.TrimLeft(line, " \t")
	return len(t) > 6 && strings.HasPrefix(t, "- [") && strings.ContainsRune(" xX", rune(t[3])) && t[4:6] == "] " && strings.TrimSpace(t[6:]) != ""
}

// appendNotes is body with text added at the end of its ## Notes section, or
// under a new one at the end when it has none.
func appendNotes(body, text string) string {
	body = strings.TrimRight(body, "\n")
	notes := strings.Index(body, "\n## Notes")
	if notes < 0 {
		return body + "\n\n## Notes\n\n" + text + "\n"
	}
	end := len(body)
	if next := strings.Index(body[notes+1:], "\n## "); next >= 0 {
		end = notes + 1 + next
	}
	head, tail := strings.TrimRight(body[:end], "\n"), body[end:]
	if tail != "" {
		// tail starts at the newline before the next heading
		tail += "\n"
	}
	return head + "\n\n" + text + "\n" + tail
}

// lastLines are the last n lines of text, with no trailing newline.
func lastLines(text string, n int) string {
	lines := strings.Split(strings.TrimRight(text, "\n"), "\n")
	if len(lines) > n {
		lines = lines[len(lines)-n:]
	}
	return strings.Join(lines, "\n")
}

// fenced is text in a fenced code block whose fence is longer than any run of
// backticks in it.
func fenced(text string) string {
	longest, run := 0, 0
	for _, r := range text {
		if r == '`' {
			run++
			longest = max(longest, run)
		} else {
			run = 0
		}
	}
	fence := strings.Repeat("`", max(3, longest+1))
	if strings.TrimSpace(text) == "" {
		text = "(no output)"
	}
	return fence + "text\n" + text + "\n" + fence
}
