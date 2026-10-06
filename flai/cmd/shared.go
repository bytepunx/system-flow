package cmd

import (
	"fmt"
	"io"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/spf13/cobra"

	"github.com/bytepunx/system-flow/flai/internal/manifest"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// sharedDialect is the glob dialect of claims.shared, as each command's help
// gives it (ADR-0096).
const sharedDialect = `Patterns are paths relative to the repository root, separated by /: * is any
characters within one segment, ** zero or more whole segments, ? one character,
and a plain path covers itself and everything below it.`

// newSharedCmd groups the commands that read and change the manifest's
// claims.shared, the paths whose overlaps hold no story (S-0295, ADR-0096).
func newSharedCmd(a *app) *cobra.Command {
	c := &cobra.Command{
		Use:   "shared",
		Short: "List, check, add, and remove the shared paths, whose overlaps hold no story",
		Long: `claims.shared in system-flow.yaml lists the paths many stories change in
separate sections or new files. An overlap of two stories' claims that lies
wholly inside one of them holds no ready story, and flai check does not report
it (ADR-0096).

` + sharedDialect,
		Example: `  flai shared list
  flai shared check docs/users/flai.md flai/cmd S-0295
  flai shared add design/adrs 'docs/users/*.md'
  flai shared remove design/adrs`,
	}
	c.AddCommand(newSharedListCmd(a), newSharedCheckCmd(a), newSharedAddCmd(a), newSharedRemoveCmd(a))
	return c
}

func newSharedListCmd(a *app) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "Print the shared paths' patterns, one per line",
		Long: `Print claims.shared's patterns in list order, one per line; with --json, a
JSON array of strings, [] when there is none. A pattern that is not valid is
printed too, frees nothing, and is warned about on stderr.

` + sharedDialect,
		Example: `  flai shared list
  flai shared list --json`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			repo, err := a.project()
			if err != nil {
				return err
			}
			claims := repo.Manifest.Claims
			a.warnInvalidShared(claims)
			if a.jsonOut {
				return a.printJSON(nonNil(claims.Shared))
			}
			if len(claims.Shared) == 0 {
				fmt.Fprintln(a.out, "no shared paths; add one with flai shared add <pattern>")
				return nil
			}
			for _, p := range claims.Shared {
				fmt.Fprintln(a.out, p)
			}
			return nil
		},
	}
}

func newSharedAddCmd(a *app) *cobra.Command {
	var autocommit bool
	var trailers []string
	c := &cobra.Command{
		Use:   "add <pattern>...",
		Short: "Add patterns to the shared paths",
		Long: `Add patterns to the end of claims.shared in system-flow.yaml, rewriting only that
list and keeping the file's other keys and comments. A pattern that is not
valid (empty, absolute, with a .. segment, or a malformed glob) or that is in
the list already is refused with the reason, and nothing is written. The change
is committed, system-flow.yaml alone, only with --autocommit and unless the
project sets dashboard.autocommit: false. Prints each pattern added; with --json, {"added": [...],
"shared": [...]}, shared being the list after the change. Run flai shared check
first to see what a pattern would free.

` + sharedDialect,
		Example: `  flai shared add design/adrs
  flai shared add 'docs/users/*.md' 'design/**/README.md'`,
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return a.editShared(args, manifest.AddShared, autocommit, "chore: add to the shared paths", trailers)
		},
	}
	sharedCommitFlags(c, &autocommit, &trailers)
	return c
}

func newSharedRemoveCmd(a *app) *cobra.Command {
	var autocommit bool
	var trailers []string
	c := &cobra.Command{
		Use:   "remove <pattern>...",
		Short: "Remove patterns from the shared paths",
		Long: `Remove patterns from claims.shared in system-flow.yaml, rewriting only that list
and keeping the file's other keys and comments. Name each pattern as flai
shared list prints it; one that is not in the list is refused, and nothing is
written. A pattern that is not valid may be removed. The change is committed,
system-flow.yaml alone, only with --autocommit and unless the project sets
dashboard.autocommit: false. Prints each pattern removed; with --json, {"removed": [...],
"shared": [...]}, shared being the list after the change.

` + sharedDialect,
		Example: `  flai shared remove design/adrs`,
		Args:    cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return a.editShared(args, manifest.RemoveShared, autocommit, "chore: remove from the shared paths", trailers)
		},
	}
	sharedCommitFlags(c, &autocommit, &trailers)
	return c
}

// sharedCommitFlags registers the flags that commit an edit of the shared
// paths, as flai agent set's do.
func sharedCommitFlags(c *cobra.Command, autocommit *bool, trailers *[]string) {
	c.Flags().BoolVar(autocommit, "autocommit", false, "commit system-flow.yaml on its own, unless dashboard.autocommit is false")
	c.Flags().StringArrayVar(trailers, "trailer", nil, "a trailer line for the commit (repeatable)")
}

// editShared changes claims.shared in the project's manifest with edit,
// commits it when autocommit asks and the project allows, and says what
// changed.
func (a *app) editShared(patterns []string, edit func(string, ...string) (manifest.Change, error), autocommit bool, msg string, trailers []string) error {
	repo, err := a.project()
	if err != nil {
		return err
	}
	ch, err := edit(filepath.Join(repo.Root, manifest.File), patterns...)
	if err != nil {
		return err
	}
	a.commitManifest(repo, autocommit, msg+": "+strings.Join(append(ch.Added, ch.Removed...), ", "), trailers)
	if a.jsonOut {
		return a.printJSON(ch)
	}
	for _, p := range ch.Added {
		fmt.Fprintf(a.out, "added %s\n", p)
	}
	for _, p := range ch.Removed {
		fmt.Fprintf(a.out, "removed %s\n", p)
	}
	return nil
}

// sharedCheck is one entry flai shared check reports.
type sharedCheck struct {
	// Entry is the argument as given, or an entry of a story's claim.
	Entry string `json:"entry"`
	// Path is the entry as it is matched: a component's name or tag read as
	// its path, without a leading ./ or a trailing /.
	Path string `json:"path"`
	// Shared reports whether the entry lies wholly inside a pattern.
	Shared bool `json:"shared"`
	// Pattern is the first pattern it lies inside, when it is shared.
	Pattern string `json:"pattern,omitempty"`
	// Story is the story whose claim the entry is, when a story was named.
	Story string `json:"story,omitempty"`
}

// itemIDArg is an argument read as a work item's ID rather than a path: a
// letter, a dash, and a number in any padding.
var itemIDArg = regexp.MustCompile(`^[EeSsTt]-\d+$`)

func newSharedCheckCmd(a *app) *cobra.Command {
	return &cobra.Command{
		Use:   "check <path-or-story>...",
		Short: "Say whether paths, touches entries, or a story's claim lie inside the shared paths",
		Long: `For each path or touches entry, say whether it lies wholly inside a pattern of
claims.shared and which one. A component's name or tag is read as its path. An
entry whose last segment has an extension, such as flai.md, is a file and lies
inside a pattern that matches it; any other is read as a folder, which lies
inside only a pattern that covers everything below it. Given a story ID
(S-0295, s-295), report each entry of the story's claim, as its hold reads it,
the same way: what adding a pattern would free. Nothing is written, and the
exit status is 0 whether or not an entry is shared.

With --json, an array with an object for each entry: entry (the argument, or
the claim's entry), path (the entry as matched), shared (true or false),
pattern (the pattern matched, only when shared), and story (only for a story's
claim entries).

` + sharedDialect,
		Example: `  flai shared check docs/users/flai.md design/adrs/0096-x.md flai/cmd
  flai shared check S-0295 --json`,
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			repo, err := a.project()
			if err != nil {
				return err
			}
			claims := repo.Manifest.Claims
			a.warnInvalidShared(claims)
			paths := workitem.NewHolds(nil, repo.Manifest.Projects)
			checks := []sharedCheck{}
			var empty []string
			for _, arg := range args {
				arg = strings.TrimSpace(arg)
				if arg == "" {
					return fmt.Errorf("an empty argument names no path; give a path, a touches entry, or a story ID")
				}
				if !itemIDArg.MatchString(arg) {
					checks = append(checks, checkShared(claims, arg, entryPath(paths, arg), ""))
					continue
				}
				claim, id, err := storyClaim(repo, arg)
				if err != nil {
					return err
				}
				if len(claim) == 0 {
					empty = append(empty, id)
				}
				for _, e := range claim {
					checks = append(checks, checkShared(claims, e, e, id))
				}
			}
			if a.jsonOut {
				return a.printJSON(checks)
			}
			printSharedChecks(a.out, checks, empty)
			return nil
		},
	}
}

// checkShared is whether a touches entry, read as path, lies inside a
// pattern of claims.
func checkShared(claims manifest.Claims, entry, path, story string) sharedCheck {
	p, ok := claims.Covers(path)
	return sharedCheck{Entry: entry, Path: path, Shared: ok, Pattern: p, Story: story}
}

// entryPath is a touches entry as a claim reads it: a component's name or
// tag becomes its path.
func entryPath(paths *workitem.Holds, entry string) string {
	if c := paths.Claim(&workitem.Item{Touches: []string{entry}}); len(c) == 1 {
		return c[0]
	}
	return entry
}

// storyClaim is the claim of the story id names, as its hold reads it, and
// the story's ID as flai writes it.
func storyClaim(repo *workitem.Repo, id string) ([]string, string, error) {
	if workitem.TypeOfID(workitem.CanonicalID(id)) != workitem.Story {
		return nil, "", fmt.Errorf("%s is not a story; flai shared check takes a story's ID, a path, or a touches entry", workitem.CanonicalID(id))
	}
	story, err := repo.Get(id)
	if err != nil {
		return nil, "", err
	}
	items, err := repo.List(story.Archived)
	if err != nil {
		return nil, "", fmt.Errorf("cannot read the tasks of %s: %w; run flai check to see what is wrong", story.ID, err)
	}
	return repo.Holds(items).Claim(story), story.ID, nil
}

// printSharedChecks writes a line for each entry checked, then one for each
// story whose claim is empty.
func printSharedChecks(w io.Writer, checks []sharedCheck, empty []string) {
	for _, c := range checks {
		name := c.Entry
		if c.Path != c.Entry {
			name += " (" + c.Path + ")"
		}
		if c.Story != "" {
			name = c.Story + " " + name
		}
		if c.Shared {
			fmt.Fprintf(w, "%s: shared, inside %s\n", name, c.Pattern)
		} else {
			fmt.Fprintf(w, "%s: not shared\n", name)
		}
	}
	for _, id := range empty {
		fmt.Fprintf(w, "%s claims nothing: it and its open tasks touch nothing\n", id)
	}
}

// warnInvalidShared warns on stderr of each pattern of claims.shared that is
// not valid, which frees nothing.
func (a *app) warnInvalidShared(claims manifest.Claims) {
	for _, e := range claims.Errors() {
		a.logger().Warn("shared path pattern not valid, frees nothing", "component", "manifest", "pattern", e.Pattern, "reason", e.Reason, "fix", "flai shared remove "+e.Pattern)
	}
}
