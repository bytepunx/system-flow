package cmd

import (
	"fmt"
	"io"

	"github.com/spf13/cobra"

	"github.com/bytepunx/system-flow/flai/internal/docedit"
)

// newFileCmd groups the commands that read and save a text file that is not
// markdown, as flai doc does a document (S-0436, ADR-0147, ADR-0150).
func newFileCmd(a *app) *cobra.Command {
	c := &cobra.Command{
		Use:   "file",
		Short: "Read and save a text file that is not markdown",
		Long: fmt.Sprintf(`The save path for a text file that is not markdown, such as the design
token file, edited in the dashboard (ADR-0147). It follows flai doc's rules:
a content hash detects a concurrent change, the check validates the save,
and the save is committed. A file under the design or docs folder, and the
design token file and design stylesheet the manifest names, are reached;
the wip folder, a markdown file (flai doc show and flai doc save handle
those), a binary file, a file over %d bytes (1 MiB), and any other path are
refused, each with its reason. A save of the token file writes the design
stylesheet from the new tokens and commits both in one commit (ADR-0150).`, docedit.MaxFileSize),
	}
	c.AddCommand(newFileShowCmd(a), newFileSaveCmd(a))
	return c
}

func newFileShowCmd(a *app) *cobra.Command {
	var story string
	c := &cobra.Command{
		Use:   "show <path>",
		Short: "Print a text file that is not markdown, with its content hash",
		Long: fmt.Sprintf(`Prints the file's content as it is on standard output, and its path and
content hash on standard error, so that the content can be redirected to a
file and saved again with flai file save --hash. With --json, prints
{path, content, hash}.

Reaches a file under the design or docs folder, and the design token file
and design stylesheet the manifest names. Refuses, with exit status 1 and
the reason, a path outside those, a path under the wip folder, a markdown
file, which flai doc show reads, a binary file, and a file over %d bytes
(1 MiB).

With --story, reads the file in that story's worktree, as it is on its
branch, for flai file save --story (ADR-0145). A story without a worktree is
refused.`, docedit.MaxFileSize),
		Example: `  flai file show design/system/design-system.tokens.json --json
  flai file show web/src/theme.css > theme.css
  flai file show design/system/design-system.tokens.json --story S-0042 --json`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			repo, err := a.project()
			if err != nil {
				return err
			}
			f, err := docedit.FileShow(repo, args[0], story)
			if err != nil {
				return err
			}
			if a.jsonOut {
				return a.printJSON(f)
			}
			fmt.Fprint(a.out, f.Content)
			fmt.Fprintf(a.errOut, "%s\n  hash: %s\n", f.Path, f.Hash)
			return nil
		},
	}
	c.Flags().StringVar(&story, "story", "", "read the file in this story's worktree (S-nnnn)")
	return c
}

func newFileSaveCmd(a *app) *cobra.Command {
	var opt docedit.SaveOptions
	c := &cobra.Command{
		Use:   "save <path> [--hash <sha256>]",
		Short: "Save a text file that is not markdown from standard input: conflict check, flai check, commit",
		Long: fmt.Sprintf(`Reads the new content on standard input and saves it as it is. --hash is
the hash flai file show gave for the content the editor loaded; when the
file no longer has it the save fails as a conflict (exit 3) with the
current content, its hash, and a diff. Without --hash, the save creates a
file that is absent, with its folders, and removes them again when the
check refuses it; without --hash on a file that exists, the save fails as
a conflict (ADR-0164). Content flai will not save is
refused (exit 4) and the file is left as it was: binary content, content
over %d bytes (1 MiB), a token file flai cannot render as a stylesheet, and content
the check has a new finding with. A saved file is committed on its own, as
docs, unless --no-commit is given or the project sets dashboard.autocommit:
false. Nothing is pushed.

Reaches what flai file show reaches and refuses what it refuses, with exit
status 1 and the reason.

A save of the design token file, when the manifest names a design
stylesheet, writes the stylesheet from the new tokens as flai design css
--write does, checks with both in place, puts both back when the check
refuses, and commits both paths in one commit (ADR-0150). --json prints the
result with stylesheet, the stylesheet's path, set.

With --story, saves the file in that story's worktree: the hash and the
check are the branch's, and the commit is on story/S-nnnn, with the story
named in its subject; the main checkout is left as it was (ADR-0145). A
story without a worktree and a worktree not on its branch are refused.`, docedit.MaxFileSize),
		Example: `  flai file save design/system/design-system.tokens.json --hash "$h" --message "darker primary" < tokens.json
  flai file save design/system/palette.json --story S-0042 --hash "$h" < palette.json`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			repo, err := a.project()
			if err != nil {
				return err
			}
			// One byte past the limit is enough for FileSave to refuse it.
			data, err := io.ReadAll(io.LimitReader(cmd.InOrStdin(), docedit.MaxFileSize+1))
			if err != nil {
				return fmt.Errorf("read the content to save from standard input: %w", err)
			}
			opt.Now = a.now()
			res, err := docedit.FileSave(repo, a.runner, args[0], string(data), opt)
			if c, ok := docedit.IsConflict(err); ok {
				if a.jsonOut {
					_ = a.printJSON(map[string]any{"conflict": c})
				} else if c.Diff != "" {
					fmt.Fprintln(a.out, c.Diff)
				}
				return &exitError{code: exitDocConflict, msg: c.Error()}
			}
			if r, ok := docedit.IsRefused(err); ok {
				if a.jsonOut {
					_ = a.printJSON(map[string]any{"refused": r})
				} else {
					for _, f := range r.Findings {
						fmt.Fprintf(a.out, "%s:%d: %s: %s: %s\n", f.Path, f.Line, f.Level, f.Rule, f.Message)
					}
				}
				return &exitError{code: exitDocRefused, msg: r.Error()}
			}
			if err != nil {
				return err
			}
			if a.jsonOut {
				return a.printJSON(res)
			}
			saved := res.Path
			if res.Stylesheet != "" {
				saved += " and wrote " + res.Stylesheet + " from it"
			}
			switch {
			case res.Unchanged:
				fmt.Fprintf(a.out, "%s is unchanged\n", res.Path)
			case res.Committed:
				fmt.Fprintf(a.out, "saved %s, committed %s\n", saved, res.Commit)
			case res.CommitError != "":
				fmt.Fprintf(a.out, "saved %s, NOT committed (%s)\n", saved, firstLine(res.CommitError))
			default:
				fmt.Fprintf(a.out, "saved %s, not committed\n", saved)
			}
			for _, f := range res.Warnings {
				fmt.Fprintf(a.out, "%s:%d: %s: %s: %s\n", f.Path, f.Line, f.Level, f.Rule, f.Message)
			}
			return nil
		},
	}
	c.Flags().StringVar(&opt.Hash, "hash", "", "hash of the content the editor loaded (from flai file show)")
	c.Flags().StringVar(&opt.Message, "message", "", "commit subject after the type prefix (default: edit <path>)")
	c.Flags().StringArrayVar(&opt.Trailers, "trailer", nil, "line appended to the commit message (repeatable)")
	c.Flags().BoolVar(&opt.NoCommit, "no-commit", false, "save without committing")
	c.Flags().StringVar(&opt.Story, "story", "", "save the file in this story's worktree and commit it on its branch (S-nnnn)")
	return c
}
