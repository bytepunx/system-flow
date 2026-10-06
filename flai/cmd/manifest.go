package cmd

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/bytepunx/system-flow/flai/internal/manifest"
)

// newManifestCmd groups the commands that change system-flow.yaml's
// strategic agents' settings (S-0229, ADR-0039).
func newManifestCmd(a *app) *cobra.Command {
	c := &cobra.Command{
		Use:   "manifest",
		Short: "Change the strategic agents' settings in system-flow.yaml",
		Long: `The strategic agents' settings live in system-flow.yaml: what the orchestrator
may do (orchestration.permissions), how it orders the ready column and when a
release is due (orchestration.policy, orchestration.release), and the planner's
and the analyzer's agents and schedules (planning, analysis). flai manifest set
writes them, checked as flai check checks them; the dashboard's settings panels
run it.`,
		Example: `  flai manifest set orchestration.permissions.promote_to_ready=true
  flai manifest set planning.replan=agent planning.schedule='0 6 * * 1-5'`,
	}
	c.AddCommand(newManifestSetCmd(a))
	return c
}

func newManifestSetCmd(a *app) *cobra.Command {
	var unset, trailers []string
	var autocommit bool
	c := &cobra.Command{
		Use:   "set <key>=<value>...",
		Short: "Write the strategic agents' settings, each checked, or refuse with the field and the reason",
		Long: `Write each key=value into its block of system-flow.yaml, and remove each
--unset key so that its default applies, keeping every other key, the order,
and the comments. The keys:

` + settingsHelp() + `

A boolean is true or false; a number is written as digits, a count a whole
one; a duration is a Go duration such as 168h; a cron expression has five
fields in UTC, or is daily; an agent is JSON, such as
{"harness":"claude-code","model":"claude-sonnet-5","config":{"effort":"medium"}},
and {} unsets it. planning.currency is not among them: changing it
re-denominates every amount on the items, so it is edited by hand.

The manifest is checked whole, as it would be after the change, as flai check
checks it. A key outside the list, a value not of its key's kind, or a problem
the result would have is refused with the field and the reason, such as
orchestration.release.value: -1 is not an amount of zero or more, and nothing
is written; the exit status is 1, as for a workflow rule's refusal. With
--json a refusal is {"refused": [{"field": ..., "reason": ...}]} on standard
output, and a change {"set": [{"key": ..., "value": ...}], "unset": [...]},
with "commit" when it was committed.

The change is committed, system-flow.yaml alone, only with --autocommit and
unless the project sets dashboard.autocommit: false.`,
		Example: `  flai manifest set orchestration.permissions.promote_to_ready=true orchestration.policy=wsjf
  flai manifest set orchestration.release.policy=threshold orchestration.release.value=500
  flai manifest set planning.agent='{"model":"claude-sonnet-5"}' --unset planning.schedule
  flai manifest set analysis.schedule=daily --autocommit --json`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 && len(unset) == 0 {
				return fmt.Errorf("nothing to change: give key=value to set, or --unset key, such as flai manifest set planning.replan=agent")
			}
			var set []manifest.Assignment
			var malformed []manifest.Problem
			for _, arg := range args {
				k, v, ok := strings.Cut(arg, "=")
				if !ok {
					malformed = append(malformed, manifest.Problem{Field: strings.TrimSpace(arg), Reason: "is not key=value; write one such as planning.replan=agent, or --unset the key"})
					continue
				}
				set = append(set, manifest.Assignment{Key: strings.TrimSpace(k), Value: v})
			}
			if len(malformed) > 0 {
				return a.refusedSettings(&manifest.RefusedError{Problems: malformed})
			}
			repo, err := a.project()
			if err != nil {
				return err
			}
			if err := manifest.WriteSettings(filepath.Join(repo.Root, manifest.File), set, unset); err != nil {
				var refused *manifest.RefusedError
				if errors.As(err, &refused) {
					return a.refusedSettings(refused)
				}
				return err
			}
			keys := make([]string, 0, len(set)+len(unset))
			for _, s := range set {
				keys = append(keys, s.Key)
			}
			keys = append(keys, unset...)
			commit := a.commitManifest(repo, autocommit, "chore: change the manifest's settings: "+strings.Join(keys, ", "), trailers)
			return a.printSettingsChange(set, unset, commit)
		},
	}
	c.Flags().StringArrayVar(&unset, "unset", nil, "remove a key, so that its default applies (repeatable)")
	c.Flags().BoolVar(&autocommit, "autocommit", false, "commit system-flow.yaml on its own, unless dashboard.autocommit is false")
	c.Flags().StringArrayVar(&trailers, "trailer", nil, "a trailer line for the commit (repeatable)")
	return c
}

// settingsHelp lists the catalog's keys for the command's help: each with
// its kind, its values, and its default.
func settingsHelp() string {
	var b strings.Builder
	for i, s := range manifest.Settings() {
		if i > 0 {
			b.WriteString("\n")
		}
		fmt.Fprintf(&b, "  %s (%s", s.Key, s.Kind)
		if s.Whole {
			b.WriteString(", whole")
		}
		if s.Kind == manifest.KindChoice {
			fmt.Fprintf(&b, ": %s", strings.Join(s.Values, ", "))
		}
		if s.Default != "" {
			fmt.Fprintf(&b, "; default %s", s.Default)
		}
		b.WriteString(")")
	}
	return b.String()
}

// refusedSettings prints a refusal of flai manifest set, with --json as
// {"refused": [{field, reason}]}, and exits 1 as a workflow rule's refusal
// does, its message beginning rule:.
func (a *app) refusedSettings(r *manifest.RefusedError) error {
	if a.jsonOut {
		_ = a.printJSON(map[string]any{"refused": r.Problems})
	} else {
		for _, p := range r.Problems {
			fmt.Fprintln(a.out, p.String())
		}
	}
	return &exitError{code: 1, msg: "rule: " + r.Error()}
}

// printSettingsChange says what flai manifest set wrote, and the commit
// when there is one.
func (a *app) printSettingsChange(set []manifest.Assignment, unset []string, commit string) error {
	if a.jsonOut {
		if set == nil {
			set = []manifest.Assignment{}
		}
		out := map[string]any{"set": set, "unset": nonNil(unset)}
		if commit != "" {
			out["commit"] = commit
		}
		return a.printJSON(out)
	}
	for _, s := range set {
		fmt.Fprintf(a.out, "set %s = %s\n", s.Key, s.Value)
	}
	for _, k := range unset {
		fmt.Fprintf(a.out, "unset %s\n", k)
	}
	if commit != "" {
		fmt.Fprintf(a.out, "committed %s\n", commit)
	}
	return nil
}
