package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/bytepunx/system-flow/flai/internal/check"
	"github.com/bytepunx/system-flow/flai/internal/gitver"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

func newCheckCmd(a *app) *cobra.Command {
	var strict bool
	c := &cobra.Command{
		Use:   "check [dir]",
		Short: "Validate the repository against the system-flow standard",
		Long: `Check the manifest, layout, every work item in kanban and archive, the
narratives and their index, the board, and documentation front matter.

Findings print as path:line: level: rule: message. Errors exit 1; with
--strict warnings do too, except two that only the operator clears: the
review column over its limit, which acceptance clears, and an epic behind
its stories (epic.lags-stories), which moving or accepting the epic clears.`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var repo *workitem.Repo
			var err error
			if len(args) == 1 {
				repo, err = workitem.Open(args[0])
			} else {
				repo, err = a.project()
			}
			if err != nil {
				return err
			}
			res, err := check.Run(repo, a.now())
			if err != nil {
				return err
			}
			// The clone may need a newer git than the one installed (ADR-0022).
			// No git on PATH means nothing here depends on its version.
			if v, err := gitver.Installed(a.runner); err == nil {
				check.GitCompat(res, repo, v)
			}
			if a.jsonOut {
				if err := a.printJSON(res); err != nil {
					return err
				}
			} else {
				for _, f := range res.Findings {
					fmt.Fprintf(a.out, "%s:%d: %s: %s: %s\n", f.Path, f.Line, f.Level, f.Rule, f.Message)
				}
				fmt.Fprintf(a.out, "%d items checked, %d errors, %d warnings", res.Items, res.Errors, res.Warnings)
				if res.Advisory > 0 {
					// S-0243: say which warnings --strict does not fail on.
					fmt.Fprintf(a.out, " (%d that --strict passes over: only the operator clears them, by accepting or by moving an epic)", res.Advisory)
				}
				fmt.Fprintln(a.out)
			}
			if !res.OK(strict) {
				return &exitError{code: 1}
			}
			return nil
		},
	}
	c.Flags().BoolVar(&strict, "strict", false, "treat warnings as failures, except the review column over its limit and an epic behind its stories")
	return c
}

// exitError signals a non-zero exit, with an optional message the boundary
// logs as the fatal event.
type exitError struct {
	code int
	msg  string
}

func (e *exitError) Error() string { return e.msg }
