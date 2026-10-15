package cmd

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/spf13/cobra"

	"github.com/bytepunx/system-flow/flai/internal/designsystem"
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
renders the tokens as CSS.`,
		Example: `  flai design init --topics flaiover
  flai design show
  flai design show --json
  flai design css --tailwind`,
	}
	c.AddCommand(newDesignInitCmd(a), newDesignShowCmd(a), newDesignCSSCmd(a))
	return c
}

func newDesignInitCmd(a *app) *cobra.Command {
	var list []string
	c := &cobra.Command{
		Use:   "init",
		Short: "Write the token file and the guidance from flai's starters",
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

When either file is there already, init writes neither and says which, with
exit status 1. It commits nothing. Prints the two paths written; with --json,
{"tokens": ..., "guidance": ..., "topics": [...]}.`,
		Example: `  flai design init
  flai design init --topics flaiover,ui`,
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
			if a.jsonOut {
				return a.printJSON(map[string]any{"tokens": tokens, "guidance": guidance, "topics": list})
			}
			fmt.Fprintf(a.out, "wrote %s\nwrote %s\n", tokens, guidance)
			return nil
		},
	}
	c.Flags().StringSliceVar(&list, "topics", []string{"all"}, "the guidance's topics, which choose the stories prime gives it to: name the UI sub-project")
	return c
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
	return &cobra.Command{
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

--json prints what the host method design.get answers: {path, exists,
no_token_file, guidance, headings, resolver, themes: [{name, tokens: [{path,
type, value, css, variables, ref, line}]}], problems, css, tailwind,
contrast: [{theme, pairs: [{theme, foreground, background, foreground_hex,
background_hex, ratio, threshold, passes, clipped}]}]}.`,
		Example: `  flai design show
  flai design show --json`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			repo, err := a.project()
			if err != nil {
				return err
			}
			fsys := os.DirFS(repo.Root)
			shown, err := designsystem.Show(fsys, repo.Manifest.DesignTokens(), repo.Manifest.DesignGuidance())
			if err != nil {
				return err
			}
			if a.jsonOut {
				return a.printJSON(shown)
			}
			return printDesign(a.out, fsys, shown, repo.Manifest.DesignGuidance())
		},
	}
}

// printDesign writes flai design show's summary of shown.
func printDesign(w io.Writer, fsys fs.FS, shown *designsystem.Shown, guidance string) error {
	if !shown.Exists {
		fmt.Fprintf(w, "no token file at %s: run flai design init to write one and its guidance\n", shown.Path)
		return nil
	}
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	layout := "a plain token file"
	if shown.Resolver {
		layout = "a resolver document"
	}
	fmt.Fprintf(tw, "token file\t%s, %s\n", shown.Path, layout)
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

// designProblem is one problem of the token file file as
// <file>:<line>: <token path>: <message>, the token path left out when the
// problem is the whole file's.
func designProblem(file string, p designsystem.Problem) string {
	if p.Path == "" {
		return fmt.Sprintf("%s:%d: %s", file, p.Line, p.Message)
	}
	return fmt.Sprintf("%s:%d: %s: %s", file, p.Line, p.Path, p.Message)
}

func newDesignCSSCmd(a *app) *cobra.Command {
	var tailwind bool
	c := &cobra.Command{
		Use:   "css",
		Short: "Print the tokens as CSS custom properties, or as a Tailwind v4 @theme block",
		Long: `Print the token file's tokens as CSS custom properties with Tailwind v4's
names, such as --color-primary and --text-sm, one block per theme, each with
every token: :root for light, then [data-theme='dark'] for dark. With
--tailwind, print them as a Tailwind v4 @theme block of the light values,
then a [data-theme='dark'] block of the dark values that differ from them.

No token file is an error, exit status 1, naming the file and flai design
init. So is a token file with problems: css renders nothing rather than CSS
that silently leaves out or misreads a token, and the error lists each
problem, as flai design show and flai check report them; fix them and run it
again.`,
		Example: `  flai design css
  flai design css --tailwind > src/theme.css`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			repo, err := a.project()
			if err != nil {
				return err
			}
			name := repo.Manifest.DesignTokens()
			sys, err := designsystem.Load(os.DirFS(repo.Root), name)
			switch {
			case errors.Is(err, fs.ErrNotExist):
				return fmt.Errorf("there is no token file at %s; run flai design init to write one and its guidance", name)
			case err != nil:
				return err
			}
			if len(sys.Problems) > 0 {
				lines := make([]string, 0, len(sys.Problems))
				for _, p := range sys.Problems {
					lines = append(lines, designProblem(name, p))
				}
				return fmt.Errorf("flai design css renders nothing from the token file %s while it has problems; fix them and run it again: %s", name, strings.Join(lines, "; "))
			}
			if tailwind {
				fmt.Fprint(a.out, designsystem.Tailwind(sys))
			} else {
				fmt.Fprint(a.out, designsystem.CSS(sys))
			}
			return nil
		},
	}
	c.Flags().BoolVar(&tailwind, "tailwind", false, "print a Tailwind v4 @theme block instead of one block of custom properties per theme")
	return c
}
