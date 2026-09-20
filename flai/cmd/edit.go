package cmd

import (
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/spf13/cobra"

	"github.com/bytepunx/system-flow/flai/internal/docedit"
	"github.com/bytepunx/system-flow/flai/internal/itemedit"
)

// flai edit changes a work item after it was made (S-0085).
func newEditCmd(a *app) *cobra.Command {
	var title, nature, parent, hash, message, byFlag string
	var tags, touches, trailers []string
	var clearTags, clearTouches, bodyStdin, autocommit, show bool
	c := &cobra.Command{
		Use:   "edit <id>",
		Short: "Change an item's title, nature, tags, touches, parent, or body, checked and in one step",
		Long: `Change what an item says about itself. Any of the fields and the body can
change together. What is the item's state stays flai's and is changed by its
own commands: status by flai move, blocking by flai block, never here.

A retitle keeps everything that carries the title in step: the front matter,
the heading, the file's name, the line in the parent's list, the story's
narrative, and links to the old file name in design, docs, and wip. A new
parent must be an open item of the right type; the item leaves the old
parent's list and joins the new one. An archived or closed item is refused.

--body-stdin reads what lies below the heading; the heading is the ID and the
title, and flai writes it. With --hash, the hash flai edit --show printed, a
change someone made meanwhile is a conflict (exit 3) and nothing is written.
flai check runs with the change in place: if it reports anything the change
introduces, every file is put back and the findings are printed (exit 4).
--autocommit commits every file the edit touched in one commit, unless the
project sets dashboard.autocommit: false. Nothing is pushed.

Agents connected over MCP are told of an edit someone else made, and what of
the item changed.`,
		Example: `  flai edit S-0085 --show
  flai edit S-0085 --title "Items are editable from the dashboard" --autocommit
  flai edit S-0085 --tag dashboard --tag cli --touches flaiover/src
  flai edit S-0085 --parent E-0004
  flai edit S-0085 --body-stdin --hash 3f0c... < body.md`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			repo, err := a.project()
			if err != nil {
				return err
			}
			if show {
				v, err := itemedit.Show(repo, args[0])
				if err != nil {
					return err
				}
				if a.jsonOut {
					return a.printJSON(v)
				}
				fmt.Fprintf(a.out, "%s %s\n  %s, %s, %s\n  tags: %s\n", v.ID, v.Title, v.Type, v.Nature, v.Status, strings.Join(v.Tags, ", "))
				if v.Type != "epic" {
					fmt.Fprintf(a.out, "  touches: %s\n  parent: %s\n", strings.Join(v.Touches, ", "), v.Parent)
				}
				fmt.Fprintf(a.out, "  file: %s\n  hash: %s\n", v.Path, v.Hash)
				if !v.Editable {
					fmt.Fprintf(a.out, "  not editable: %s\n", v.Reason)
				}
				fmt.Fprintf(a.out, "\n%s", v.Body)
				return nil
			}
			var ch itemedit.Change
			f := cmd.Flags()
			if f.Changed("title") {
				ch.Title = &title
			}
			if f.Changed("nature") {
				ch.Nature = &nature
			}
			if f.Changed("parent") {
				ch.Parent = &parent
			}
			switch {
			case clearTags && f.Changed("tag"):
				return fmt.Errorf("--tag and --clear-tags contradict each other")
			case clearTags:
				ch.Tags = &[]string{}
			case f.Changed("tag"):
				ch.Tags = &tags
			}
			switch {
			case clearTouches && f.Changed("touches"):
				return fmt.Errorf("--touches and --clear-touches contradict each other")
			case clearTouches:
				ch.Touches = &[]string{}
			case f.Changed("touches"):
				ch.Touches = &touches
			}
			if bodyStdin {
				data, err := io.ReadAll(cmd.InOrStdin())
				if err != nil {
					return err
				}
				body := string(data)
				ch.Body = &body
			}
			if ch == (itemedit.Change{}) {
				return fmt.Errorf("nothing to change: give --title, --nature, --tag, --touches, --parent, or --body-stdin (flai edit %s --show prints what is there)", args[0])
			}
			by, _ := agentIdentity()
			if cfg, _, err := a.loadConfig(); err == nil && by == "agent" && cfg.Author != "" {
				by = cfg.Author
			}
			if byFlag != "" {
				by = byFlag
			}
			res, err := itemedit.Apply(repo, a.runner, args[0], ch, itemedit.Options{Hash: hash, By: by, Message: message, Trailers: trailers, NoCommit: !autocommit, Now: a.now()})
			if c, ok := docedit.IsConflict(err); ok {
				if a.jsonOut {
					_ = a.printJSON(map[string]any{"conflict": c})
				}
				return &exitError{code: exitDocConflict, msg: c.Error()}
			}
			if r, ok := docedit.IsRefused(err); ok {
				if a.jsonOut {
					_ = a.printJSON(map[string]any{"refused": r})
				} else {
					for _, fd := range r.Findings {
						fmt.Fprintf(a.out, "%s:%d: %s: %s: %s\n", fd.Path, fd.Line, fd.Level, fd.Rule, fd.Message)
					}
				}
				return &exitError{code: exitDocRefused, msg: r.Error()}
			}
			var inv *itemedit.InvalidError
			if errors.As(err, &inv) {
				// a rule, as flai move says its refusals: the caller's mistake, which a
				// dashboard answers with 400 and not with 500
				return fmt.Errorf("rule: %s", inv.Error())
			}
			if err != nil {
				return err
			}
			if a.jsonOut {
				return a.printJSON(res)
			}
			if res.Unchanged {
				fmt.Fprintf(a.out, "%s is unchanged\n", res.ID)
				return nil
			}
			fmt.Fprintf(a.out, "%s: changed %s\n  %s\n", res.ID, strings.Join(res.Changed, ", "), res.Path)
			if res.Renamed != "" {
				fmt.Fprintf(a.out, "  renamed from %s\n", res.Renamed)
			}
			switch {
			case res.Committed:
				fmt.Fprintf(a.out, "  committed %s (%d file(s))\n", res.Commit, len(res.Files))
			case res.CommitError != "":
				fmt.Fprintf(a.out, "  changed, not committed: %s\n", res.CommitError)
			}
			return nil
		},
	}
	f := c.Flags()
	f.BoolVar(&show, "show", false, "print the fields, the body, the hash, and whether the item may be edited")
	f.StringVar(&title, "title", "", "the new title")
	f.StringVar(&nature, "nature", "", "one of feature, improvement, remediation, research, experiment")
	f.StringSliceVar(&tags, "tag", nil, "the tags, replacing the ones there (repeatable or comma separated)")
	f.BoolVar(&clearTags, "clear-tags", false, "remove every tag")
	f.StringSliceVar(&touches, "touches", nil, "the paths or components the work changes, replacing the ones there")
	f.BoolVar(&clearTouches, "clear-touches", false, "remove the list")
	f.StringVar(&parent, "parent", "", "the new parent: an epic for a story, a story for a task")
	f.BoolVar(&bodyStdin, "body-stdin", false, "read the body below the heading from standard input")
	f.StringVar(&hash, "hash", "", "the hash flai edit --show printed; a change made meanwhile is then a conflict")
	f.StringVar(&byFlag, "by", "", "who edits, as agents are told (default: FLAI_AGENT, then the config author)")
	f.StringVar(&message, "message", "", "commit subject after the prefix (default names what changed)")
	f.BoolVar(&autocommit, "autocommit", false, "commit every file the edit touched, unless dashboard.autocommit is false")
	f.StringArrayVar(&trailers, "trailer", nil, "trailer line for the commit (repeatable)")
	return c
}
