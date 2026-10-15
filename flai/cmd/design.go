package cmd

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/spf13/cobra"

	"github.com/bytepunx/system-flow/flai/internal/designsystem"
	"github.com/bytepunx/system-flow/flai/internal/manifest"
	"github.com/bytepunx/system-flow/flai/internal/topics"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// newDesignCmd groups the commands on a project's design system: its token
// file and its guidance (S-0381, ADR-0146, ADR-0148).
func newDesignCmd(a *app) *cobra.Command {
	c := &cobra.Command{
		Use:   "design",
		Short: "Write, show, and render the project's design system",
		Long: `A project's design system is its tokens, the values its UI uses, and its
guidance, the words that say when to use which (ADR-0146). The tokens are one
DTCG 2025.10 file the manifest's design.tokens names,
design/system/design-system.tokens.json when it is unset; the guidance is
design/system/design-system.md. flai check validates the token file as
design.tokens; flai design init writes both from starters, flai design show
reports them with each colour role pair's contrast, and flai design css
renders the tokens as CSS. flai design css --write writes the tokens to the
generated design stylesheet the manifest's design.stylesheet names, which
flai check reports as design.stale when it falls behind the tokens
(ADR-0159).

A story proposes a change to the design system on its branch: flai design
proposals lists the open stories whose branches change it, and flai design
show --story S-nnnn reads the design system a story proposes from its
worktree. Accepting the story merges its branch, and with it the proposal.`,
		Example: `  flai design init --topics flaiover
  flai design show
  flai design show --json
  flai design show --story S-0042
  flai design proposals
  flai design css --tailwind
  flai design css --write`,
	}
	c.AddCommand(newDesignInitCmd(a), newDesignShowCmd(a), newDesignProposalsCmd(a), newDesignCSSCmd(a))
	return c
}

func newDesignInitCmd(a *app) *cobra.Command {
	var list, projects []string
	c := &cobra.Command{
		Use:   "init",
		Short: "Write the token file and the guidance from flai's starters, and the design stylesheet when one is named",
		Long: `Write the token file the manifest's design.tokens names, and the guidance,
design/system/design-system.md, from flai's starters, creating their folders.

The token file is a DTCG 2025.10 resolver document: one set of light values
and the theme modifier, whose dark context overrides each colour role and
each shadow. It holds every required token: the colour roles, the sans and
mono font stacks (system stacks, no web fonts), the type scale xs to 4xl,
the spacing step, and the radii and shadows sm, md, and lg, at Tailwind v4's
default sizes. Every role pair passes its contrast threshold in light and
dark. The guidance has the headings a design system's words go under, each
saying what to write there, and its front matter's updated is today.

--topics sets the guidance's topics, all by default. The topics choose the
stories flai prime gives the guidance to, so name the UI sub-project, or its
kind, so that every story that touches the UI is given it.

--project names a sub-project the design system governs, one of the
manifest's projects; repeat it or separate names with commas. init writes
each name into the guidance's topics, the one list prime reads, so there is
no second record of the sub-projects governed (ADR-0152). With --project and
no --topics, the topics are the sub-projects alone, without all; with both,
they are the --topics values, then each sub-project not among them. A name
the manifest does not list is refused before anything is written.

When the manifest's design.stylesheet names a design stylesheet, init then
writes it from the token file it has just written, as flai design css --write
does, replacing what is there. When that fails, the token file and the
guidance stay, and the error says to run flai design css --write.

When the token file or the guidance is there already, init writes nothing and
says which, with exit status 1. It commits nothing. Prints the paths written,
the token file, the guidance, and the stylesheet when there is one; with
--json, {"tokens": ..., "guidance": ..., "topics": [...], "stylesheet": ...},
stylesheet left out when the manifest names none.`,
		Example: `  flai design init
  flai design init --topics flaiover,ui
  flai design init --project flaiover`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			for _, t := range list {
				if !topics.Valid(t) {
					return fmt.Errorf("--topics: %q is not a topic; a topic is one word of letters, digits, dot, dash, or underscore, such as flaiover", t)
				}
			}
			repo, err := a.project()
			if err != nil {
				return err
			}
			list, err = governedTopics(list, cmd.Flags().Changed("topics"), projects, repo.Manifest.Projects)
			if err != nil {
				return err
			}
			tokens, guidance := repo.Manifest.DesignTokens(), repo.Manifest.DesignGuidance()
			body, err := starterGuidance(tokens, a.now(), list)
			if err != nil {
				return err
			}
			if err := writeNewFiles(repo.Root, []newFile{
				{tokens, "token file", designsystem.StarterTokens},
				{guidance, "guidance", []byte(body)},
			}); err != nil {
				return err
			}
			a.logger().Info("design system written", "component", "cmd", "tokens", tokens, "guidance", guidance)
			stylesheet := repo.Manifest.DesignStylesheet()
			if stylesheet != "" {
				if err := initStylesheet(repo.Root, stylesheet, tokens); err != nil {
					return err
				}
				a.logger().Info("design stylesheet written", "component", "cmd", "stylesheet", stylesheet, "tokens", tokens)
			}
			if a.jsonOut {
				out := map[string]any{"tokens": tokens, "guidance": guidance, "topics": list}
				if stylesheet != "" {
					out["stylesheet"] = stylesheet
				}
				return a.printJSON(out)
			}
			fmt.Fprintf(a.out, "wrote %s\nwrote %s\n", tokens, guidance)
			if stylesheet != "" {
				fmt.Fprintf(a.out, "wrote %s\n", stylesheet)
			}
			return nil
		},
	}
	c.Flags().StringSliceVar(&list, "topics", []string{"all"}, "the guidance's topics, which choose the stories prime gives it to: name the UI sub-project")
	c.Flags().StringSliceVar(&projects, "project", nil, "a sub-project of the manifest the design system governs, written into the guidance's topics; repeatable")
	return c
}

// governedTopics is the guidance's topics once the sub-projects the design
// system governs are written into them (ADR-0152): the topics list, its
// default all dropped when explicit is false and there are sub-projects, then
// each sub-project not already in it. A name that is no sub-project of known
// is refused, naming the ones there are.
func governedTopics(list []string, explicit bool, projects []string, known []manifest.Project) ([]string, error) {
	if len(projects) == 0 {
		return list, nil
	}
	names := make([]string, 0, len(known))
	for _, p := range known {
		names = append(names, p.Name)
	}
	for _, p := range projects {
		if !slices.Contains(names, p) {
			if len(names) == 0 {
				return nil, fmt.Errorf("--project: %q is not a sub-project of this project, and %s lists none; add it under projects in %s, or leave --project out and name the topics with --topics", p, manifest.File, manifest.File)
			}
			return nil, fmt.Errorf("--project: %q is not a sub-project of this project; %s lists %s, so name one of those", p, manifest.File, strings.Join(names, ", "))
		}
	}
	var out []string
	if explicit {
		out = slices.Clone(list)
	}
	for _, p := range projects {
		if !slices.Contains(out, p) {
			out = append(out, p)
		}
	}
	return out, nil
}

// initStylesheet writes the design stylesheet from the token file flai design
// init has just written. The token file and the guidance stay when it fails,
// and the error says how to write the stylesheet again.
func initStylesheet(root, stylesheet, tokens string) error {
	sys, err := loadDesignTokens(root, tokens)
	if err == nil {
		err = writeStylesheet(root, stylesheet, sys, tokens)
	}
	if err != nil {
		return fmt.Errorf("flai design init wrote the token file and the guidance but not the design stylesheet %s: %w; run flai design css --write to write it", stylesheet, err)
	}
	return nil
}

// starterGuidance is the starter guidance with its front matter's updated
// set to today, its topics to list, and the token file's path in its body.
func starterGuidance(tokens string, now time.Time, list []string) (string, error) {
	fm, body, err := workitem.SplitFrontMatter(designsystem.StarterGuidance)
	if err != nil {
		return "", fmt.Errorf("read flai's starter guidance: %w; this flai is broken, reinstall it", err)
	}
	lines := strings.Split(fm, "\n")
	for i, l := range lines {
		if strings.HasPrefix(l, "updated:") {
			lines[i] = "updated: " + now.UTC().Format(time.DateOnly)
		}
	}
	fm = topics.SetInFrontMatter(strings.Join(lines, "\n"), list)
	return "---\n" + fm + "---\n" + strings.ReplaceAll(body, designsystem.StarterTokensMark, tokens), nil
}

// newFile is a file flai design init writes: its repository path, what it
// is, and its content.
type newFile struct {
	path, what string
	data       []byte
}

// writeNewFiles writes each file under root, creating its folders, or none
// of them when any is there already, naming it. A file that fails to write
// takes the ones written before it away again.
func writeNewFiles(root string, files []newFile) error {
	for _, f := range files {
		_, err := os.Stat(filepath.Join(root, filepath.FromSlash(f.path)))
		switch {
		case err == nil:
			return fmt.Errorf("the %s %s is there already, so flai design init writes neither file; edit it, or remove it and run flai design init again", f.what, f.path)
		case !errors.Is(err, fs.ErrNotExist):
			return fmt.Errorf("look for the %s %s: %w", f.what, f.path, err)
		}
	}
	var written []string
	for _, f := range files {
		abs := filepath.Join(root, filepath.FromSlash(f.path))
		if err := writeNew(abs, f.data); err != nil {
			for _, w := range written {
				_ = os.Remove(w) // best effort: the error below says what failed
			}
			return fmt.Errorf("write the %s %s: %w; nothing was written", f.what, f.path, err)
		}
		written = append(written, abs)
	}
	return nil
}

// writeNew creates the file at path with data, its folders too, and refuses
// to replace one that is there.
func writeNew(path string, data []byte) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o644)
	if err != nil {
		return err
	}
	if _, err := f.Write(data); err != nil {
		_ = f.Close()
		return err
	}
	return f.Close()
}

func newDesignShowCmd(a *app) *cobra.Command {
	var story string
	c := &cobra.Command{
		Use:   "show",
		Short: "Report the token file, the guidance, the problems, and each role pair's contrast",
		Long: `Report the project's design system: the token file's path, the guidance's
path and its headings, the number of tokens in each theme, each problem flai
check reports as design.tokens, as <file>:<line>: <token path>: <message>,
and, for light and dark, each colour role pair's contrast ratio against
ADR-0148's threshold: 4.5:1 for text, muted text, each state on background,
surface, and its tint, and on-primary and on-secondary on their fill; 3:1 for
border and focus on background. A pair whose oklch colour lies outside sRGB
is computed on the colour clipped to sRGB, and says so.

Contrast is a report, never a finding: a pair below its threshold is FAIL in
the table and the exit status is still 0. So is a token file with problems,
which flai check reports as findings; show prints them and exits 0. No token
file is no error either: show says so and how to write one.

With --story S-nnnn, show reads the design system a story proposes from the
story's worktree instead of the main checkout, at the token file and
guidance paths the worktree's manifest gives, and names the story and its
worktree first; the host method design.get does the same when given a story.
It works the same from the main checkout and from any worktree. It refuses,
with exit status 1, an ID that is no story of this project, a story that is
done or cancelled, which proposes nothing, and a story with no worktree,
which flai stream open S-nnnn opens. flai design proposals lists the stories
whose branches change the design system.

--json prints what the host method design.get answers: {path, exists,
no_token_file, guidance, stylesheet, headings, resolver, themes: [{name,
tokens: [{path, type, value, css, variables, ref, line}]}], problems, css,
tailwind, contrast: [{theme, pairs: [{theme, foreground, background,
foreground_hex, background_hex, ratio, threshold, passes, clipped}]}]}, the
same with --story. stylesheet is the manifest's design.stylesheet, left out
when it names none.`,
		Example: `  flai design show
  flai design show --json
  flai design show --story S-0042`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			repo, err := a.project()
			if err != nil {
				return err
			}
			var (
				shown *designsystem.Shown
				fsys  fs.FS
				from  string
			)
			if story != "" {
				var id string
				if shown, id, err = designsystem.ShowStory(repo, story); err != nil {
					return err
				}
				wt := repo.WorktreePath(id)
				fsys, from = os.DirFS(wt), fmt.Sprintf("%s, read from its worktree %s", id, wt)
			} else {
				fsys = os.DirFS(repo.Root)
				if shown, err = designsystem.Show(fsys, repo.Manifest.DesignTokens(), repo.Manifest.DesignGuidance()); err != nil {
					return err
				}
				shown.Stylesheet = repo.Manifest.DesignStylesheet()
			}
			if a.jsonOut {
				return a.printJSON(shown)
			}
			return printDesign(a.out, fsys, shown, from)
		},
	}
	c.Flags().StringVar(&story, "story", "", "read the design system the open story S-nnnn proposes, from its worktree, instead of the main checkout's")
	return c
}

// printDesign writes flai design show's summary of shown, read from fsys;
// from names the story whose worktree it was read from, or is empty for the
// checkout flai runs in.
func printDesign(w io.Writer, fsys fs.FS, shown *designsystem.Shown, from string) error {
	if !shown.Exists {
		if from != "" {
			fmt.Fprintf(w, "story %s\n", from)
		}
		fmt.Fprintf(w, "no token file at %s: run flai design init to write one and its guidance\n", shown.Path)
		return nil
	}
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	if from != "" {
		fmt.Fprintf(tw, "story\t%s\n", from)
	}
	layout := "a plain token file"
	if shown.Resolver {
		layout = "a resolver document"
	}
	fmt.Fprintf(tw, "token file\t%s, %s\n", shown.Path, layout)
	guidance := shown.Guidance
	if _, err := fs.Stat(fsys, guidance); errors.Is(err, fs.ErrNotExist) {
		fmt.Fprintf(tw, "guidance\t%s, not there\n", guidance)
	} else {
		fmt.Fprintf(tw, "guidance\t%s, %d headings\n", guidance, len(shown.Headings))
	}
	counts := make([]string, 0, len(shown.Themes))
	for _, th := range shown.Themes {
		counts = append(counts, fmt.Sprintf("%s %d", th.Name, len(th.Tokens)))
	}
	if len(counts) == 0 {
		counts = append(counts, "none read")
	}
	fmt.Fprintf(tw, "tokens\t%s\n", strings.Join(counts, ", "))
	if len(shown.Problems) == 0 {
		fmt.Fprintf(tw, "problems\tnone\n")
	} else {
		fmt.Fprintf(tw, "problems\t%d, which flai check reports as design.tokens\n", len(shown.Problems))
	}
	if err := tw.Flush(); err != nil {
		return err
	}
	for _, p := range shown.Problems {
		fmt.Fprintln(w, designProblem(shown.Path, p))
	}
	for _, c := range shown.Contrast {
		fmt.Fprintf(w, "\ncontrast in %s\n", c.Theme)
		tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
		fmt.Fprintln(tw, "  PAIR\tRATIO\tTHRESHOLD\tRESULT\tCLIPPED")
		for _, p := range c.Pairs {
			result, clipped := "pass", "no"
			if !p.Passes {
				result = "FAIL"
			}
			if p.Clipped {
				clipped = "yes"
			}
			fmt.Fprintf(tw, "  %s on %s\t%.2f\t%.1f\t%s\t%s\n", p.Foreground, p.Background, p.Ratio, p.Threshold, result, clipped)
		}
		if err := tw.Flush(); err != nil {
			return err
		}
	}
	return nil
}

// designProblem is one problem of the token file as
// <file>:<line>: <token path>: <message>, the token path left out when the
// problem is the whole file's.
func designProblem(file string, p designsystem.Problem) string {
	if p.Path == "" {
		return fmt.Sprintf("%s:%d: %s", file, p.Line, p.Message)
	}
	return fmt.Sprintf("%s:%d: %s: %s", file, p.Line, p.Path, p.Message)
}

func newDesignProposalsCmd(a *app) *cobra.Command {
	return &cobra.Command{
		Use:   "proposals",
		Short: "List the open stories whose branches change the design system",
		Long: `List the proposals: the open stories whose branch, or whose worktree
uncommitted, changes the design system, that is the token file, the
guidance, or both, at the paths the main checkout's manifest gives. A
story's branch is compared with its merge base with the main branch, as
flai stream diff compares it, and git is read in the main checkout, so the
list is the same from the main checkout and from any worktree. A story that
is done or cancelled proposes nothing and is left out; so is every story of
a project that is not a git repository, which has no branches.

The Design page lists the same through the host method design.proposals,
and flai design show --story S-nnnn reads the design system a proposal
holds from its story's worktree. A proposal needs no step of its own to
land: accepting the story merges its branch, and with it the main checkout's
design system becomes the proposal's.

Prints a table of ID, STATE, TITLE, and CHANGES, the paths the story
changes, in story ID order, or a line saying no open story's branch changes
the design system, naming the token file and the guidance. With --json, the
list as the host method design.proposals answers it: [{id, title, status,
paths}], [] when there are none.`,
		Example: `  flai design proposals
  flai design proposals --json`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			repo, err := a.project()
			if err != nil {
				return err
			}
			if repo, err = mainCheckoutRepo(repo); err != nil {
				return err
			}
			list, err := designsystem.Proposals(a.runner, repo)
			if err != nil {
				return err
			}
			if a.jsonOut {
				return a.printJSON(list)
			}
			if len(list) == 0 {
				fmt.Fprintf(a.out, "no open story's branch changes the design system, the token file %s or the guidance %s\n", repo.Manifest.DesignTokens(), repo.Manifest.DesignGuidance())
				return nil
			}
			tw := tabwriter.NewWriter(a.out, 0, 0, 2, ' ', 0)
			fmt.Fprintln(tw, "ID\tSTATE\tTITLE\tCHANGES")
			for _, p := range list {
				fmt.Fprintf(tw, "%s\t%s\t%s\t%s\n", p.ID, p.Status, p.Title, strings.Join(p.Paths, ", "))
			}
			return tw.Flush()
		},
	}
}

// mainCheckoutRepo is the repository at the main checkout when repo is a
// story's worktree, so that the design system's paths are the ones the main
// checkout's manifest gives, and repo otherwise.
func mainCheckoutRepo(repo *workitem.Repo) (*workitem.Repo, error) {
	if !repo.IsWorktree() {
		return repo, nil
	}
	main, err := workitem.Open(repo.MainRoot)
	if err != nil {
		return nil, fmt.Errorf("open the main checkout %s of the worktree %s: %w; run the command in the main checkout", repo.MainRoot, repo.Root, err)
	}
	main.Git, main.TemplateDir = repo.Git, repo.TemplateDir
	return main, nil
}

func newDesignCSSCmd(a *app) *cobra.Command {
	var tailwind, write bool
	c := &cobra.Command{
		Use:   "css",
		Short: "Print the tokens as CSS custom properties or a Tailwind v4 @theme block, or write the design stylesheet",
		Long: `Print the token file's tokens as CSS custom properties with Tailwind v4's
names, such as --color-primary and --text-sm, one block per theme, each with
every token: :root for light, then [data-theme='dark'] for dark. With
--tailwind, print them as a Tailwind v4 stylesheet: a @custom-variant rule
that makes Tailwind's dark: variant follow [data-theme='dark'], a @theme
block of the light values, then a [data-theme='dark'] block of the dark
values that differ from them, so that the dark values and the dark: variant
follow the same attribute (ADR-0159).

With --write, write the design stylesheet instead of printing: the file the
manifest's design.stylesheet names, holding what --tailwind prints with a
header ahead of it that says the file is generated from the token file by
flai design css --write and is not to be edited. It creates the file's
folders and replaces what is there, since the file is generated, and prints
the path written; with --json, {"stylesheet": ..., "tokens": ...}. When the
manifest names no stylesheet, --write writes nothing and exits with status 1:
set design.stylesheet in system-flow.yaml to the stylesheet's repository
path, such as web/src/theme.css. Once the stylesheet is there, flai check
reports design.stale whenever it falls behind the tokens; run flai design css
--write again to bring it up to date. --write and --tailwind are mutually
exclusive.

No token file is an error, exit status 1, naming the file and flai design
init. So is a token file with problems: css renders and writes nothing rather
than CSS that silently leaves out or misreads a token, and the error lists
each problem, as flai design show and flai check report them; fix them and
run it again.`,
		Example: `  flai design css
  flai design css --tailwind
  flai design css --write`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			repo, err := a.project()
			if err != nil {
				return err
			}
			stylesheet := repo.Manifest.DesignStylesheet()
			if write && stylesheet == "" {
				return fmt.Errorf("the manifest names no design stylesheet, so flai design css --write writes nothing; set design.stylesheet in %s to the stylesheet's repository path, such as web/src/theme.css, and run it again", manifest.File)
			}
			name := repo.Manifest.DesignTokens()
			sys, err := loadDesignTokens(repo.Root, name)
			if err != nil {
				return err
			}
			switch {
			case write:
				if err := writeStylesheet(repo.Root, stylesheet, sys, name); err != nil {
					return err
				}
				a.logger().Info("design stylesheet written", "component", "cmd", "stylesheet", stylesheet, "tokens", name)
				if a.jsonOut {
					return a.printJSON(map[string]any{"stylesheet": stylesheet, "tokens": name})
				}
				fmt.Fprintf(a.out, "wrote %s\n", stylesheet)
			case tailwind:
				fmt.Fprint(a.out, designsystem.Tailwind(sys))
			default:
				fmt.Fprint(a.out, designsystem.CSS(sys))
			}
			return nil
		},
	}
	c.Flags().BoolVar(&tailwind, "tailwind", false, "print a Tailwind v4 stylesheet, the dark: variant, a @theme block, and the dark values, instead of one block of custom properties per theme; --write writes it")
	c.Flags().BoolVar(&write, "write", false, "write the Tailwind v4 stylesheet, with a generated-file header, to the file the manifest's design.stylesheet names")
	c.MarkFlagsMutuallyExclusive("tailwind", "write")
	return c
}

// loadDesignTokens reads the token file name under root, and refuses when
// there is none or when it has problems, listing each, so that nothing is
// rendered that silently leaves out or misreads a token.
func loadDesignTokens(root, name string) (*designsystem.System, error) {
	sys, err := designsystem.Load(os.DirFS(root), name)
	switch {
	case errors.Is(err, fs.ErrNotExist):
		return nil, fmt.Errorf("there is no token file at %s; run flai design init to write one and its guidance", name)
	case err != nil:
		return nil, err
	}
	if len(sys.Problems) > 0 {
		lines := make([]string, 0, len(sys.Problems))
		for _, p := range sys.Problems {
			lines = append(lines, designProblem(name, p))
		}
		return nil, fmt.Errorf("flai design css renders nothing from the token file %s while it has problems; fix them and run it again: %s", name, strings.Join(lines, "; "))
	}
	return sys, nil
}

// writeStylesheet writes the generated design stylesheet, the repository
// path stylesheet under root, from sys, read from the token file tokens,
// creating its folders and replacing what is there (S-0385, ADR-0159).
func writeStylesheet(root, stylesheet string, sys *designsystem.System, tokens string) error {
	abs := filepath.Join(root, filepath.FromSlash(stylesheet))
	if err := os.MkdirAll(filepath.Dir(abs), 0o755); err != nil {
		return fmt.Errorf("write the design stylesheet %s: %w", stylesheet, err)
	}
	if err := os.WriteFile(abs, []byte(designsystem.Stylesheet(sys, tokens)), 0o644); err != nil {
		return fmt.Errorf("write the design stylesheet %s: %w", stylesheet, err)
	}
	return nil
}
