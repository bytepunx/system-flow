package cmd

import (
	"bytes"
	"fmt"
	"maps"
	"os"
	"path"
	"path/filepath"
	"slices"
	"strings"

	"github.com/spf13/cobra"

	"github.com/bytepunx/system-flow/flai/internal/atomicfile"
	"github.com/bytepunx/system-flow/flai/internal/formulas"
	"github.com/bytepunx/system-flow/flai/internal/guard"
	"github.com/bytepunx/system-flow/flai/internal/itemnew"
	"github.com/bytepunx/system-flow/flai/internal/manifest"
	"github.com/bytepunx/system-flow/flai/internal/roles"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

func newFormulaCmd(a *app) *cobra.Command {
	c := &cobra.Command{
		Use:   "formula",
		Short: "Read the project's formulas, one file each under design/workflows/formulas, and make a story from one (ADR-0238)",
		Long: `A formula is a repeatable process as data: one markdown file, <name>.md, in
the design folder's workflows/formulas (ADR-0238). Its front matter names the
formula, the nature of the story it makes, its inputs, and its steps. Each
step runs a flai command, a script under scripts/, or a role with a prompt,
after the steps its after names, in the instance's worktree or, with in:
main, in the project's main checkout. An input is used as {{name}} in the title,
a step's title, and a step's run. flai check reports a formula file that
falls short as formulas.file.

flai formula run makes from a formula one ordinary story in backlog, with one
task per step, which a story's agent works as it works any story.`,
	}
	c.AddCommand(newFormulaListCmd(a), newFormulaShowCmd(a), newFormulaRunCmd(a))
	return c
}

func newFormulaListCmd(a *app) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List the project's formulas with their nature, step count, file, and title",
		Long: `List the formulas in this checkout's design folder's workflows/formulas, in
file order: each one's name, the nature of the story it makes, its number of
steps, and its file, and below it its title. A formula file that does not
load is left out and said on stderr; flai check reports it. --json gives an
array of {name, title, nature, status, steps, path}, steps the count, [] when
none.`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			_, list, err := a.loadFormulas()
			if err != nil {
				return err
			}
			type listed struct {
				Name   string `json:"name"`
				Title  string `json:"title"`
				Nature string `json:"nature"`
				Status string `json:"status"`
				Steps  int    `json:"steps"`
				Path   string `json:"path"`
			}
			out := make([]listed, len(list))
			for i, f := range list {
				out[i] = listed{Name: f.Name, Title: f.Title, Nature: f.Nature, Status: f.Status, Steps: len(f.Steps), Path: f.Path}
			}
			if a.jsonOut {
				return a.printJSON(out)
			}
			if len(out) == 0 {
				fmt.Fprintln(a.out, "no formulas")
				return nil
			}
			for _, f := range out {
				fmt.Fprintf(a.out, "%s  %s  %s  %s\n  %s\n", f.Name, f.Nature, plural(f.Steps, "step"), f.Path, f.Title)
			}
			return nil
		},
	}
}

func newFormulaShowCmd(a *app) *cobra.Command {
	return &cobra.Command{
		Use:   "show <name>",
		Short: "Print a formula: its inputs, its steps and what each runs, and the steps' layers",
		Long: `Print one formula by its name: its nature, file, title, status, updated, and
description; its inputs, each required or with its default, and what it is;
its steps, each with its title, what it runs, the steps it is after, and
in: main when it runs in the project's main checkout; and
its layers, the steps grouped by the longest chain of after before each, so
that the tasks of a layer can run at once once those of the layers before it
are done. The problems flai check would report in it are said on stderr.
--json gives the formula, each step's kind among its fields, and its layers.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			repo, f, err := a.formula(args[0])
			if err != nil {
				return err
			}
			for _, p := range formulas.Validate([]formulas.Formula{f}, formulaNames(repo)) {
				fmt.Fprintln(a.errOut, p)
			}
			layers := formulas.Layers(f)
			if a.jsonOut {
				type step struct {
					formulas.Step
					Kind string `json:"kind"`
				}
				steps := make([]step, len(f.Steps))
				for i, s := range f.Steps {
					steps[i] = step{Step: s, Kind: s.Kind()}
				}
				return a.printJSON(struct {
					formulas.Formula
					Steps  []step     `json:"steps"`
					Layers [][]string `json:"layers"`
				}{f, steps, layers})
			}
			a.printFormula(f, layers)
			return nil
		},
	}
}

// printFormula prints what flai formula show says of a formula.
func (a *app) printFormula(f formulas.Formula, layers [][]string) {
	fmt.Fprintf(a.out, "%s  %s  %s\n", f.Name, f.Nature, f.Path)
	for _, kv := range [][2]string{{"title", f.Title}, {"status", f.Status}, {"updated", f.Updated}, {"description", f.Description}} {
		fmt.Fprintf(a.out, "  %-12s %s\n", kv[0]+":", kv[1])
	}
	if len(f.Inputs) == 0 {
		fmt.Fprintln(a.out, "inputs: none")
	} else {
		fmt.Fprintln(a.out, "inputs:")
	}
	for _, in := range f.Inputs {
		given := "required"
		if !in.Required() {
			given = fmt.Sprintf("default %q", *in.Default)
		}
		fmt.Fprintf(a.out, "  %s  %s  %s\n", in.Name, given, in.Description)
	}
	fmt.Fprintln(a.out, "steps:")
	for _, s := range f.Steps {
		fmt.Fprintf(a.out, "  %s: %s\n", s.ID, s.Title)
		switch s.Kind() {
		case formulas.KindFlai:
			fmt.Fprintf(a.out, "    flai:   flai %s\n", s.Run.Flai)
		case formulas.KindScript:
			fmt.Fprintf(a.out, "    script: %s\n", s.Run.Script)
		case formulas.KindRole:
			fmt.Fprintf(a.out, "    role:   %s\n    prompt: %s\n", s.Run.Role, s.Run.Prompt)
		default:
			fmt.Fprintln(a.out, "    runs:   none of flai, script, or role alone; flai check reports it")
		}
		if len(s.After) > 0 {
			fmt.Fprintf(a.out, "    after:  %s\n", strings.Join(s.After, ", "))
		}
		// S-0423: a step that runs in the main checkout says so; one in the
		// worktree, the default, prints as before.
		if s.In != "" && s.In != formulas.InWorktree {
			fmt.Fprintf(a.out, "    in:     %s\n", s.In)
		}
	}
	fmt.Fprintln(a.out, "layers:")
	for i, l := range layers {
		fmt.Fprintf(a.out, "  %d  %s\n", i+1, strings.Join(l, ", "))
	}
}

func newFormulaRunCmd(a *app) *cobra.Command {
	var inputs, trailers []string
	var epic string
	var autocommit, ready bool
	c := &cobra.Command{
		Use:   "run <name>",
		Short: "Make a story from a formula, with one task per step, in backlog or, with --ready, in ready",
		Long: `Make the formula's instance (ADR-0238 §7): one story in backlog, as flai
story new makes it, and one task of it per step. Each --input, name=value,
gives one of the formula's inputs; an input not given takes its default.

Before it writes anything it refuses a formula flai check reports a problem
in, an --input that is not name=value or is given twice, an input the
formula does not declare, and a required input not given.

The story has the formula's nature, formula: <name> in its front matter, the
formula's title with the inputs filled as its title, the description as its
goal, one acceptance criterion per step, and notes naming the formula file
and the inputs the run took. --epic makes it a story of that epic; a story
need not belong to one. Each task has the step's title, a Work section
naming what the step runs, and for a step with in: main that it runs from
the project's main checkout, by its absolute path, rather than the
instance's worktree, a Done when section that is its exit (the
command or script exits 0, or the role's run ends with its work done), and
the tasks of the step's after as its after, so that the tasks' layers are
the steps'. The tasks are made in file order, each after those it waits for.

Each item is checked as flai story new --body-stdin checks one: a finding it
brings refuses it, and the story and the tasks made before it are removed
(exit 4). --ready then moves the story to ready as flai move does, by
FLAI_AGENT when set, its epic following and the board placing it last, so
that flai serve starts its agent; a move the workflow rules refuse removes
what it made too. Without --ready the story stays in backlog. --autocommit
commits the story, its tasks, its epic, and, with --ready, the board and the
stream index the move wrote, in one commit, unless the project sets
dashboard.autocommit: false. Nothing is pushed. It prints the story's ID,
path, and tasks; --json gives {story, path, status, tasks: [{id, step}],
formula, inputs}, status backlog or ready, with followed when the epic moved
with it and commit when it committed.`,
		Example: `  flai formula list
  flai formula show release-check
  flai formula run release-check --input version=1.4.0 --epic E-0007
  flai formula run publish --ready --autocommit`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			given, err := formulaInputs(inputs)
			if err != nil {
				return err
			}
			return a.formulaRun(args[0], given, epic, formulaRunOptions{ready: ready, autocommit: autocommit, trailers: trailers})
		},
	}
	c.Flags().StringArrayVar(&inputs, "input", nil, "one of the formula's inputs, as name=value (repeatable)")
	c.Flags().StringVar(&epic, "epic", "", "the epic the story belongs to (optional: a story need not belong to one)")
	c.Flags().BoolVar(&ready, "ready", false, "move the story to ready once it and its tasks are made, so that flai serve starts its agent")
	c.Flags().BoolVar(&autocommit, "autocommit", false, "commit the story, its tasks, its epic, and the move to ready in one commit, unless dashboard.autocommit is false")
	c.Flags().StringArrayVar(&trailers, "trailer", nil, "trailer line for the commit (repeatable)")
	return c
}

// formulaInputs are the inputs --input gives, name=value each; one with no =
// or no name, or a name given twice, is refused.
func formulaInputs(given []string) (map[string]string, error) {
	out := map[string]string{}
	for _, g := range given {
		name, value, ok := strings.Cut(g, "=")
		if name = strings.TrimSpace(name); !ok || name == "" {
			return nil, fmt.Errorf("--input %q is not name=value: give each input as its name, =, and its value, such as --input version=1.4.0", g)
		}
		if _, twice := out[name]; twice {
			return nil, fmt.Errorf("--input %s is given twice: give each input once", name)
		}
		out[name] = value
	}
	return out, nil
}

// formulaRunOptions are how flai formula run leaves the instance: moved to
// ready or not, and committed or not, with the trailers.
type formulaRunOptions struct {
	ready, autocommit bool
	trailers          []string
}

// formulaRun makes the instance of the formula name with the inputs given,
// in epic when it is set: the story, touching its own story file so that no
// story in progress holds it (I-0282), then a task per step, each checked as
// it is made, then, when asked, the story's move to ready, and on a refusal
// removes what it made. It commits them when asked, and prints them.
func (a *app) formulaRun(name string, given map[string]string, epic string, opt formulaRunOptions) error {
	if os.Getenv("FLAI_ROLE") == guard.RoleOrchestrate {
		if err := a.orchestratorRunsFormula(name); err != nil {
			return err
		}
	}
	repo, f, err := a.formula(name)
	if err != nil {
		return err
	}
	if problems := formulas.Validate([]formulas.Formula{f}, formulaNames(repo)); len(problems) > 0 {
		lines := make([]string, len(problems))
		for i, p := range problems {
			lines[i] = "  " + p.String()
		}
		return fmt.Errorf("formula %s has %s, so nothing was made; fix the file and run it again:\n%s", f.Name, plural(len(problems), "problem"), strings.Join(lines, "\n"))
	}
	values, err := f.Resolve(given)
	if err != nil {
		return err
	}
	filled := f.Filled(values)
	now, owner := a.now(), a.author()
	var epicPath string
	var epicWas []byte
	if epic != "" {
		if e, err := repo.Get(epic); err == nil {
			epicPath = e.Path
			epicWas, _ = os.ReadFile(e.Path)
		}
	}
	storyBody, err := repo.TemplateBody(workitem.Story)
	if err != nil {
		return err
	}
	story, err := itemnew.Create(repo, a.runner, itemnew.Options{New: workitem.NewOptions{
		Type: workitem.Story, Title: filled.Title, Nature: f.Nature, Parent: epic, Owner: owner, Formula: f.Name, Now: now,
		Body: fillSections(storyBody, formulaStorySections(filled, values)),
	}})
	if err != nil {
		if refused, ok := a.refusedADR(err); ok {
			return refused
		}
		return err
	}
	made := []string{story.Item.Path}
	undo := func(cause error) error {
		for _, p := range made {
			if err := os.Remove(p); err != nil && !os.IsNotExist(err) {
				return fmt.Errorf("%w; removing what it made failed too, delete %s: %w", cause, strings.Join(made, ", "), err)
			}
		}
		if epicPath != "" && epicWas != nil {
			if err := atomicfile.WriteFile(epicPath, epicWas, 0o644); err != nil {
				return fmt.Errorf("%w; restoring %s failed too, take %s out of its stories: %w", cause, epicPath, story.Item.ID, err)
			}
		}
		return cause
	}
	// The instance claims its own story file (I-0282): a story with no
	// touches is held by every story in progress (ADR-0046), the one that
	// started it too, and a claim on its own file overlaps no other story's.
	story.Item.Touches = []string{story.Path}
	if err := repo.Save(story.Item); err != nil {
		return undo(fmt.Errorf("giving the story %s of formula %s its own file as its touch, so it was removed: %w", story.Item.ID, f.Name, err))
	}
	taskBody, err := repo.TemplateBody(workitem.Task)
	if err != nil {
		return undo(err)
	}
	ids := map[string]string{}
	run := formulaRunResult{Story: story.Item.ID, Path: story.Path, Status: workitem.Backlog, Tasks: []formulaTask{}, Formula: f.Name, Inputs: values}
	after := map[string][]string{}
	for _, s := range creationOrder(filled) {
		for _, step := range s.After {
			after[s.ID] = append(after[s.ID], ids[step])
		}
		res, err := itemnew.Create(repo, a.runner, itemnew.Options{New: workitem.NewOptions{
			Type: workitem.Task, Title: s.Title, Nature: f.Nature, Parent: story.Item.ID, Owner: owner, After: after[s.ID], Now: now,
			Body: fillSections(taskBody, formulaTaskSections(s, formulaMainCheckout(repo))),
		}})
		if err != nil {
			err = undo(fmt.Errorf("making the task for step %s of formula %s, so the story %s and the tasks made before it were removed: %w", s.ID, f.Name, story.Item.ID, err))
			if refused, ok := a.refusedADR(err); ok {
				return refused
			}
			return err
		}
		made = append(made, res.Item.Path)
		ids[s.ID] = res.Item.ID
		run.Tasks = append(run.Tasks, formulaTask{ID: res.Item.ID, Step: s.ID})
	}
	paths := slices.Clone(made)
	if epicPath != "" {
		paths = append(paths, epicPath)
	}
	msg := fmt.Sprintf("chore: [%s] create story: %s, from formula %s", story.Item.ID, story.Item.Title, f.Name)
	if opt.ready {
		wrote, err := a.formulaReady(repo, story.Item.ID, &run)
		if err != nil {
			return undo(fmt.Errorf("moving the story %s of formula %s to ready, so it and its tasks were removed: %w", story.Item.ID, f.Name, err))
		}
		for _, p := range wrote {
			if !slices.Contains(paths, p) {
				paths = append(paths, p)
			}
		}
		msg += ", ready"
	}
	run.Commit, run.CommitError = a.commitFormulaRun(repo, opt.autocommit, msg, opt.trailers, paths)
	if a.jsonOut {
		return a.printJSON(run)
	}
	fmt.Fprintf(a.out, "%s %s\n  %s\n", story.Item.ID, story.Item.Title, story.Path)
	for _, t := range run.Tasks {
		line := fmt.Sprintf("  %s %s (%s)", t.ID, workitem.CleanTitle(filledStep(filled, t.Step).Title), t.Step)
		if len(after[t.Step]) > 0 {
			line += ", after " + strings.Join(after[t.Step], ", ")
		}
		fmt.Fprintln(a.out, line)
	}
	if run.Status == workitem.Ready {
		fmt.Fprintf(a.out, "  %s → %s\n", run.Story, run.Status)
		a.printFollowed(run.Followed)
	}
	switch {
	case run.Commit != "":
		fmt.Fprintf(a.out, "  committed %s\n", run.Commit)
	case run.CommitError != "":
		fmt.Fprintf(a.out, "  NOT committed (%s)\n", firstLine(run.CommitError))
	}
	return nil
}

// publishFormula is the formula the orchestrator runs to publish, under
// orchestration.permissions.publish (ADR-0240).
const publishFormula = "publish"

// orchestratorRunsFormula says why the orchestrator may not run the formula
// name, or nil when it may (S-0423, ADR-0240): flai holds the run itself, as
// it holds the orchestrator's other writes, past flai guard. It runs the
// publish formula alone, and that only while orchestration.permissions.publish
// is on. The refusal changes nothing.
func (a *app) orchestratorRunsFormula(name string) error {
	if name != publishFormula {
		return fmt.Errorf("rule: the orchestrator runs only the formula %s, under orchestration.permissions.%s, and no other: running %s is the operator's or a story's agent's; ask the operator with thread_open on the item it concerns (ADR-0240)", publishFormula, manifest.PermitPublish, name)
	}
	repo, err := a.project()
	if err != nil {
		return err
	}
	if err := repo.OrchestratorPermits(manifest.PermitPublish, "runs the publish formula", "the most recently accepted story"); err != nil {
		return fmt.Errorf("rule: %w", err)
	}
	return nil
}

// formulaRunResult is what flai formula run made: the story, its file, its
// status, a task per step, the formula, the inputs the run took, the epic's
// move with the story's to ready, and its commit.
type formulaRunResult struct {
	Story       string             `json:"story"`
	Path        string             `json:"path"`
	Status      string             `json:"status"`
	Tasks       []formulaTask      `json:"tasks"`
	Formula     string             `json:"formula"`
	Inputs      map[string]string  `json:"inputs"`
	Followed    *workitem.Followed `json:"followed,omitempty"`
	Commit      string             `json:"commit,omitempty"`
	CommitError string             `json:"commit_error,omitempty"`
}

// formulaReady moves the instance story id to ready as flai move does, by
// FLAI_AGENT when set, under the workflow rules, its epic following and the
// board placing it last, and records the move in run. It returns the files
// the move wrote beside the story: the board, the epic when it followed, and
// the stream index when its text changed. A refused move writes nothing; a
// move that fails part way puts the board and the index back as they were.
func (a *app) formulaReady(repo *workitem.Repo, id string, run *formulaRunResult) ([]string, error) {
	it, err := repo.Get(id)
	if err != nil {
		return nil, err
	}
	board, index := filepath.Join(repo.KanbanDir(), workitem.BoardFile), filepath.Join(repo.AgentsDir(), "index.md")
	boardWas, boardErr := os.ReadFile(board)
	indexWas, indexErr := os.ReadFile(index)
	res, err := repo.TransitionAll(it, workitem.Ready, a.movedBy(""), "", a.now(), false)
	if err != nil {
		for _, f := range []struct {
			path string
			was  []byte
			err  error
		}{{board, boardWas, boardErr}, {index, indexWas, indexErr}} {
			if f.err == nil {
				_ = atomicfile.WriteFile(f.path, f.was, 0o644)
			} else if os.IsNotExist(f.err) {
				_ = os.Remove(f.path)
			}
		}
		return nil, err
	}
	for _, w := range res.Warnings {
		a.logger().Warn("workflow policy warning", "component", "workitem", "item", it.ID, "detail", w)
	}
	run.Status, run.Followed = it.Status, res.Followed
	wrote := []string{board}
	if res.Followed != nil {
		if e, err := repo.Get(res.Followed.ID); err == nil {
			wrote = append(wrote, e.Path)
		}
	}
	if now, err := os.ReadFile(index); err == nil && (indexErr != nil || !bytes.Equal(now, indexWas)) {
		wrote = append(wrote, index)
	}
	return wrote, nil
}

// formulaTask is a task flai formula run made and the step it is for.
type formulaTask struct {
	ID   string `json:"id"`
	Step string `json:"step"`
}

// filledStep is the step of f with the id id.
func filledStep(f formulas.Formula, id string) formulas.Step {
	i := slices.IndexFunc(f.Steps, func(s formulas.Step) bool { return s.ID == id })
	return f.Steps[i]
}

// commitFormulaRun commits the paths, absolute, and nothing else, in one
// commit when asked and the project commits what flai writes for the
// dashboard, as itemnew commits one item; it returns the commit, or why the
// commit failed, the files being written and valid either way.
func (a *app) commitFormulaRun(repo *workitem.Repo, autocommit bool, msg string, trailers, paths []string) (string, string) {
	if !autocommit || !repo.Manifest.Autocommit() {
		return "", ""
	}
	root := repo.MainRoot
	if root == "" {
		root = repo.Root
	}
	rel := make([]string, len(paths))
	for i, p := range paths {
		rel[i] = relPath(root, p)
	}
	if len(trailers) > 0 {
		msg += "\n\n" + strings.Join(trailers, "\n")
	}
	if _, err := a.runner.Run(root, "git", append([]string{"add", "--"}, rel...)...); err != nil {
		return "", err.Error()
	}
	if _, err := a.runner.Run(root, "git", append([]string{"commit", "--quiet", "-m", msg, "--"}, rel...)...); err != nil {
		return "", err.Error()
	}
	sha, err := a.runner.Run(root, "git", "rev-parse", "--short", "HEAD")
	if err != nil {
		return "", err.Error()
	}
	return strings.TrimSpace(sha), ""
}

// creationOrder is the formula's steps in file order, save that a step comes
// after each step its after names: the first step in file order whose after
// steps all came before comes next. The formula is a valid one, its after
// naming its steps and forming no cycle.
func creationOrder(f formulas.Formula) []formulas.Step {
	left := slices.Clone(f.Steps)
	placed := map[string]bool{}
	out := make([]formulas.Step, 0, len(left))
	for len(left) > 0 {
		i := slices.IndexFunc(left, func(s formulas.Step) bool {
			return !slices.ContainsFunc(s.After, func(id string) bool { return !placed[id] })
		})
		if i < 0 {
			// a cycle Validate refuses before this runs; the steps on it stay out
			return out
		}
		placed[left[i].ID] = true
		out = append(out, left[i])
		left = slices.Delete(left, i, i+1)
	}
	return out
}

// formulaStorySections are the sections the instance story's body takes
// from the formula f, filled: the goal, a criterion per step, and notes
// naming the formula's file and the inputs values the run took.
func formulaStorySections(f formulas.Formula, values map[string]string) [][2]string {
	criteria := make([]string, len(f.Steps))
	for i, s := range f.Steps {
		criteria[i] = "- [ ] " + workitem.CleanTitle(s.Title)
	}
	took := "no inputs"
	if len(values) > 0 {
		var given []string
		for _, k := range slices.Sorted(maps.Keys(values)) {
			given = append(given, codeSpan(k+"="+values[k]))
		}
		took = "the inputs " + strings.Join(given, ", ")
	}
	notes := fmt.Sprintf("Made by %s from %s, with %s.", codeSpan("flai formula run "+f.Name), codeSpan(f.Path), took)
	return [][2]string{{"Goal", f.Description}, {"Acceptance criteria", strings.Join(criteria, "\n")}, {"Notes", notes}}
}

// formulaTaskSections are the sections a step's task takes, filled: Work,
// what the step runs, and Done when, its exit. A step in main (S-0423) runs
// from the project's main checkout, main, an absolute path, rather than the
// instance's worktree, and its Work says so.
func formulaTaskSections(s formulas.Step, main string) [][2]string {
	there := " from the project's main checkout, " + codeSpan(main) + ", not from this story's worktree"
	var work, done string
	switch s.Kind() {
	case formulas.KindFlai:
		line := codeSpan("flai " + s.Run.Flai)
		work, done = "Run "+line+".", line+" exits 0."
		if s.InMain() {
			work, done = "Run "+line+there+".", line+" exits 0 in the main checkout."
		}
	case formulas.KindScript:
		line := codeSpan(s.Run.Script)
		work, done = "Run "+line+" from the repository root.", line+" exits 0."
		if s.InMain() {
			work, done = "Run "+line+there+".", line+" exits 0 in the main checkout."
		}
	case formulas.KindRole:
		quoted := strings.Split(strings.TrimSpace(s.Run.Prompt), "\n")
		for i, l := range quoted {
			quoted[i] = strings.TrimRight("> "+l, " ")
		}
		work = "Run the role " + codeSpan(s.Run.Role) + " with this prompt:\n\n" + strings.Join(quoted, "\n")
		if s.InMain() {
			work += "\n\nRun it" + there + "."
		}
		done = "The role's run ends with its work done."
	}
	return [][2]string{{"Work", work}, {"Done when", done}}
}

// formulaMainCheckout is the project's main checkout, absolute, where a step
// in main runs (S-0423): the main checkout of a story worktree, the one its
// .git file names when the main checkout is not mounted here, or the root.
func formulaMainCheckout(repo *workitem.Repo) string {
	switch {
	case repo.MainAbsent != "":
		return repo.MainAbsent
	case repo.MainRoot != "":
		return repo.MainRoot
	}
	return repo.Root
}

// codeSpan is s as markdown inline code, its fence longer than any run of
// backticks in s.
func codeSpan(s string) string {
	fence := "`"
	for strings.Contains(s, fence) {
		fence += "`"
	}
	if strings.HasPrefix(s, "`") || strings.HasSuffix(s, "`") {
		return fence + " " + s + " " + fence
	}
	return fence + s + fence
}

// fillSections is body, an item template's sections below its heading,
// with the text under each "## <name>" sections names replaced by its text,
// and a section body lacks added at the end, one blank line around each
// heading.
func fillSections(body string, sections [][2]string) string {
	type section struct {
		heading string
		lines   []string
	}
	var pre []string
	var all []section
	for _, l := range strings.Split(body, "\n") {
		if h, ok := strings.CutPrefix(l, "## "); ok {
			all = append(all, section{heading: strings.TrimSpace(h)})
			continue
		}
		if len(all) == 0 {
			pre = append(pre, l)
		} else {
			all[len(all)-1].lines = append(all[len(all)-1].lines, l)
		}
	}
	for _, s := range sections {
		i := slices.IndexFunc(all, func(x section) bool { return x.heading == s[0] })
		if i < 0 {
			all = append(all, section{heading: s[0]})
			i = len(all) - 1
		}
		all[i].lines = strings.Split(s[1], "\n")
	}
	var parts []string
	if p := strings.TrimSpace(strings.Join(pre, "\n")); p != "" {
		parts = append(parts, p)
	}
	for _, s := range all {
		block := "## " + s.heading
		if text := strings.TrimSpace(strings.Join(s.lines, "\n")); text != "" {
			block += "\n\n" + text
		}
		parts = append(parts, block)
	}
	return strings.Join(parts, "\n\n") + "\n"
}

// formula is the project and its formula named name, or an error that says
// why the file of that name does not load, or names the formulas there are.
func (a *app) formula(name string) (*workitem.Repo, formulas.Formula, error) {
	repo, err := a.project()
	if err != nil {
		return nil, formulas.Formula{}, err
	}
	list, problems, err := formulas.Load(repo.Root, repo.Manifest.Layout["design"])
	if err != nil {
		return nil, formulas.Formula{}, err
	}
	for _, p := range problems {
		if strings.TrimSuffix(path.Base(p.Path), ".md") == name {
			return nil, formulas.Formula{}, fmt.Errorf("formula %s does not load: %s; fix the file (flai check reports it)", name, p)
		}
	}
	f, err := formulas.Find(list, name)
	return repo, f, err
}

// loadFormulas reads the project's formulas and says each formula file that
// does not load on stderr.
func (a *app) loadFormulas() (*workitem.Repo, []formulas.Formula, error) {
	repo, err := a.project()
	if err != nil {
		return nil, nil, err
	}
	list, problems, err := formulas.Load(repo.Root, repo.Manifest.Layout["design"])
	if err != nil {
		return nil, nil, err
	}
	for _, p := range problems {
		fmt.Fprintf(a.errOut, "%s (left out)\n", p)
	}
	return repo, list, nil
}

// formulaNames are what a formula's steps are checked against, as flai
// check checks them: the roles the project's role files define, and the
// scripts on disk.
func formulaNames(repo *workitem.Repo) formulas.Names {
	n := formulas.Names{Script: formulas.ScriptOnDisk(repo.Root)}
	if list, _, err := roles.Load(repo.Root, repo.Manifest.Layout["design"]); err == nil {
		for _, r := range list {
			if r.Name != "" {
				n.Roles = append(n.Roles, r.Name)
			}
		}
	}
	return n
}
