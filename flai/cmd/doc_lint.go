package cmd

import (
	"fmt"
	"io"

	"github.com/spf13/cobra"

	"github.com/bytepunx/system-flow/flai/internal/hostapi"
)

// Exit codes of flai doc lint, as markdownlint-cli2 and flai test have them:
// findings are 1, and a lint that could not run is 2, kept apart from them.
const (
	exitDocLintFindings = 1
	exitDocLintUnusable = 2
)

// newDocLintCmd lints a document's content as it is edited, with the
// project's markdownlint configuration (S-0378).
func newDocLintCmd(a *app) *cobra.Command {
	var stdin bool
	c := &cobra.Command{
		Use:   "lint <path> --stdin",
		Short: "Lint markdown given on standard input as the document at path, and print its findings",
		Long: `Reads the content on standard input and lints it, as the Markdown file at
path under the design, docs, or wip folder, with the project's markdownlint
configuration and the subset of its rules flai implements (ADR-0061). The
file need not exist and is not read: an editor sends what it holds. It is
what the dashboard's editor shows as the operator types, through the host
method doc.lint; the project's markdown tier stays the authority at
acceptance. A project without a markdownlint configuration has no findings.

Prints each finding as path:line rule/alias message, or with --json
{"path": ..., "findings": [{"line", "rule", "alias", "message"}]}, an empty
list for clean content. The exit status is 0 for clean content, 1 when it
has findings, and 2 when it could not be linted: no --stdin, a path outside
the three folders, or a configuration that cannot be read.`,
		Example: `  flai doc lint docs/users/index.md --stdin --json < edited.md`,
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if !stdin {
				return &exitError{code: exitDocLintUnusable, msg: "flai doc lint lints the content on standard input, not the file: give --stdin"}
			}
			repo, err := a.project()
			if err != nil {
				return &exitError{code: exitDocLintUnusable, msg: err.Error()}
			}
			data, err := io.ReadAll(cmd.InOrStdin())
			if err != nil {
				return &exitError{code: exitDocLintUnusable, msg: "read the content to lint from standard input: " + err.Error()}
			}
			res, e := hostapi.LintDoc(repo.Root, args[0], string(data))
			if e != nil {
				return &exitError{code: exitDocLintUnusable, msg: e.Message}
			}
			if a.jsonOut {
				if err := a.printJSON(res); err != nil {
					return err
				}
			} else {
				for _, f := range res.Findings {
					fmt.Fprintf(a.out, "%s:%d %s/%s %s\n", res.Path, f.Line, f.Rule, f.Alias, f.Message)
				}
			}
			if len(res.Findings) > 0 {
				return &exitError{code: exitDocLintFindings}
			}
			return nil
		},
	}
	c.Flags().BoolVar(&stdin, "stdin", false, "lint the content on standard input (required)")
	return c
}
