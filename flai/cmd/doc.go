package cmd

import (
	"fmt"
	"io"

	"github.com/spf13/cobra"

	"github.com/bytepunx/system-flow/flai/internal/docedit"
)

// Exit codes an editor can tell apart from an ordinary failure (ADR-0023).
const (
	exitDocConflict = 3
	exitDocRefused  = 4
)

func newDocCmd(a *app) *cobra.Command {
	c := &cobra.Command{
		Use:   "doc",
		Short: "Read and save one markdown document for an editor",
		Long: `The save path for documents edited in the dashboard (ADR-0023). flai
decides what may be edited, detects a concurrent change by a content hash,
validates with the check, and commits the one path.`,
	}
	c.AddCommand(newDocShowCmd(a), newDocSaveCmd(a))
	return c
}

func newDocShowCmd(a *app) *cobra.Command {
	return &cobra.Command{
		Use:   "show <path>",
		Short: "Print a document with its content hash and what may be edited",
		Long: `Modes: full (body and front matter: design and docs), body (flai owns the
front matter: work items, narratives, board.md), none (generated files,
threads, issues, the archive), with the reason.`,
		Example: `  flai doc show design/system/overview.md --json`,
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			repo, err := a.project()
			if err != nil {
				return err
			}
			doc, err := docedit.Show(repo, args[0])
			if err != nil {
				return err
			}
			if a.jsonOut {
				return a.printJSON(doc)
			}
			fmt.Fprintf(a.out, "%s\n  mode: %s", doc.Path, doc.Mode)
			if doc.Reason != "" {
				fmt.Fprintf(a.out, " (%s)", doc.Reason)
			}
			fmt.Fprintf(a.out, "\n  hash: %s\n", doc.Hash)
			return nil
		},
	}
}

func newDocSaveCmd(a *app) *cobra.Command {
	var opt docedit.SaveOptions
	c := &cobra.Command{
		Use:   "save <path> --hash <sha256>",
		Short: "Save a document from standard input: conflict check, flai check, commit",
		Long: `Reads the new content on standard input. --hash is the hash flai doc show
gave for the content the editor loaded; when the file no longer has it the
save fails as a conflict (exit 3) with the current content, its hash, and a
diff. Content flai will not save (front matter it owns was changed, or the
check has an error with it) is refused (exit 4) and the file is left as it
was. A saved file is committed on its own unless --no-commit is given or
the project sets dashboard.autocommit: false. Nothing is pushed.`,
		Example: `  flai doc save docs/users/index.md --hash "$h" --message "clarify install" < new.md`,
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			repo, err := a.project()
			if err != nil {
				return err
			}
			if opt.Hash == "" {
				return fmt.Errorf("--hash is required: the hash flai doc show printed for the content you edited")
			}
			data, err := io.ReadAll(cmd.InOrStdin())
			if err != nil {
				return err
			}
			opt.Now = a.now()
			res, err := docedit.Save(repo, a.runner, args[0], string(data), opt)
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
			switch {
			case res.Unchanged:
				fmt.Fprintf(a.out, "%s is unchanged\n", res.Path)
			case res.Committed:
				fmt.Fprintf(a.out, "saved %s, committed %s\n", res.Path, res.Commit)
			case res.CommitError != "":
				fmt.Fprintf(a.out, "saved %s, NOT committed (%s)\n", res.Path, firstLine(res.CommitError))
			default:
				fmt.Fprintf(a.out, "saved %s, not committed\n", res.Path)
			}
			for _, f := range res.Warnings {
				fmt.Fprintf(a.out, "%s:%d: %s: %s: %s\n", f.Path, f.Line, f.Level, f.Rule, f.Message)
			}
			return nil
		},
	}
	c.Flags().StringVar(&opt.Hash, "hash", "", "hash of the content the editor loaded (from flai doc show)")
	c.Flags().StringVar(&opt.Message, "message", "", "commit subject after the type prefix (default: edit <path>)")
	c.Flags().StringArrayVar(&opt.Trailers, "trailer", nil, "line appended to the commit message (repeatable)")
	c.Flags().BoolVar(&opt.NoCommit, "no-commit", false, "save without committing")
	return c
}
