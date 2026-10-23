package cmd

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"slices"
	"strings"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"github.com/bytepunx/system-flow/flai/internal/manifest"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// subprojectReadme is the file flai subproject add and set write in a
// sub-project's folder when they create it, so that git carries the folder.
const subprojectReadme = "README.md"

// newSubprojectCmd groups the commands that list and change the manifest's
// projects, the sub-projects of the repository, without editing
// system-flow.yaml by hand (S-0491).
func newSubprojectCmd(a *app) *cobra.Command {
	c := &cobra.Command{
		Use:   "subproject",
		Short: "List, add, change, and remove the sub-projects under projects in system-flow.yaml",
		Long: `projects in system-flow.yaml lists the repository's sub-projects: each one's
name, its folder (path, relative to the repository root), its kind, such as go
or sveltekit, and its tags. These commands read and change that list, keeping
every other line, the order of the keys, and the comments.

add and set create the sub-project's folder when it is missing, with a
README.md in it so that git carries it. remove leaves the folder alone.
flai components config writes a Svelte sub-project's components key; the
install key is written by hand.`,
		Example: `  flai subproject list
  flai subproject add web --kind sveltekit
  flai subproject set web --tags ui,app
  flai subproject remove web`,
	}
	c.AddCommand(newSubprojectListCmd(a), newSubprojectAddCmd(a), newSubprojectSetCmd(a), newSubprojectRemoveCmd(a))
	return c
}

// listedSubproject is an entry of projects as flai subproject list prints it.
type listedSubproject struct {
	Name string   `json:"name"`
	Path string   `json:"path"`
	Kind string   `json:"kind"`
	Tags []string `json:"tags"`
	// Folder reports whether path names a folder in the repository.
	Folder     bool                 `json:"folder"`
	Components *manifest.Components `json:"components"`
	Install    *manifest.Install    `json:"install"`
}

func newSubprojectListCmd(a *app) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "Print each sub-project: its name, path, kind, tags, and whether its folder is there",
		Long: `Print each entry of projects in system-flow.yaml, in its order, a line each:
its name, path, kind, and tags, whether its folder is present in the
repository, and whether it sets components and install.

--json prints {"projects": [{"name", "path", "kind", "tags", "folder",
"components", "install"}]}: tags [] when there are none, folder true when the
path names a folder, and components and install as the manifest has them, or
null when the entry leaves them out.`,
		Example: `  flai subproject list
  flai subproject list --json`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			repo, err := a.project()
			if err != nil {
				return err
			}
			listed := make([]listedSubproject, len(repo.Manifest.Projects))
			for i, p := range repo.Manifest.Projects {
				listed[i] = listedSubproject{
					Name: p.Name, Path: p.Path, Kind: p.Kind, Tags: nonNil(p.Tags),
					Folder:     filepath.IsLocal(p.Path) && isFolder(filepath.Join(repo.Root, filepath.FromSlash(p.Path))),
					Components: p.Components, Install: p.Install,
				}
			}
			if a.jsonOut {
				return a.printJSON(map[string]any{"projects": listed})
			}
			return a.printSubprojects(listed)
		},
	}
}

// printSubprojects writes flai subproject list's text: a header, then a row
// for each sub-project.
func (a *app) printSubprojects(listed []listedSubproject) error {
	if len(listed) == 0 {
		_, err := fmt.Fprintln(a.out, "no sub-projects in system-flow.yaml; add one with flai subproject add <name> --kind <kind>")
		return err
	}
	mark := func(set bool, yes string) string {
		if set {
			return yes
		}
		return "-"
	}
	tw := tabwriter.NewWriter(a.out, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "NAME\tPATH\tKIND\tTAGS\tFOLDER\tCOMPONENTS\tINSTALL")
	for _, p := range listed {
		tags := strings.Join(p.Tags, ",")
		if tags == "" {
			tags = "-"
		}
		folder := "missing"
		if p.Folder {
			folder = "present"
		}
		fmt.Fprintf(tw, "%s\t%s\t%s\t%s\t%s\t%s\t%s\n", p.Name, p.Path, p.Kind, tags, folder, mark(p.Components != nil, "set"), mark(p.Install != nil, "set"))
	}
	return tw.Flush()
}

// isFolder reports whether name is a folder; anything it cannot stat is not.
func isFolder(name string) bool {
	st, err := os.Stat(name)
	return err == nil && st.IsDir()
}

// subprojectChange is what flai subproject add, set, and remove did, as
// --json prints it.
type subprojectChange struct {
	// Action is add, set, or remove.
	Action manifest.ProjectAction `json:"action"`
	// Project is the entry as written, or as it was before a remove.
	Project manifest.Project `json:"project"`
	// From is the entry's name before a set that renamed it.
	From string `json:"from,omitempty"`
	// Created is the folder the command created, relative to the
	// repository root, when it created one.
	Created string `json:"created,omitempty"`
	// Written are the files written, relative to the repository root:
	// system-flow.yaml, and the README.md of a folder created.
	Written []string `json:"written"`
	// Commit is the commit that holds them, with --autocommit.
	Commit string `json:"commit,omitempty"`
}

// subprojectCommitFlags registers the flags that commit an edit of projects.
func subprojectCommitFlags(c *cobra.Command, autocommit *bool, trailers *[]string) {
	c.Flags().BoolVar(autocommit, "autocommit", false, "commit system-flow.yaml and the README.md written, and nothing else, unless dashboard.autocommit is false")
	c.Flags().StringArrayVar(trailers, "trailer", nil, "a trailer line for the commit (repeatable)")
}

// subprojectRefusals says how flai subproject add, set, and remove refuse,
// as each one's help gives it.
const subprojectRefusals = `The manifest is checked as it would be after the change, as flai check
checks it: a name that is empty or another sub-project's, a path that is
empty, the repository root, absolute, not clean, leaves the repository, is
another sub-project's, or names a file, and an empty kind are refused with
the field and the reason on standard output, and nothing is written, neither
the manifest nor a folder; the exit status is 1, as for a workflow rule's
refusal. With --json a refusal is {"refused": [{"field", "reason"}]}.`

// subprojectJSON says what flai subproject add, set, and remove print with
// --json, as each one's help gives it.
const subprojectJSON = `--json prints {"action", "project": {"name", "path", "kind", "tags",
"components", "install"}, "from", "created", "written", "commit"}: project
the entry as written, or as it was for remove; from the old name after a
rename; created the folder made, when one was; written the files written,
relative to the repository root; and commit only with --autocommit.`

func newSubprojectAddCmd(a *app) *cobra.Command {
	var p manifest.Project
	var autocommit bool
	var trailers []string
	c := &cobra.Command{
		Use:   "add <name>",
		Short: "Add a sub-project to projects, creating its folder when it is missing",
		Long: `Add an entry to the end of projects in system-flow.yaml: the name given, its
--path, its --kind, and its --tags. --path is relative to the repository root
and defaults to the name. When the folder is missing, add creates it with a
README.md naming the sub-project and its kind, so that git carries the
folder, and prints each path it wrote; a folder that is there already gets
no README.md. The change is committed, system-flow.yaml and that README.md
alone, only with --autocommit and unless the project sets
dashboard.autocommit: false.

Set the rest of the entry afterwards: flai components config writes a Svelte
sub-project's components key, and the install key is written by hand.

` + subprojectRefusals + `

` + subprojectJSON,
		Example: `  flai subproject add web --kind sveltekit
  flai subproject add api --path services/api --kind go --tags backend,http
  flai subproject add web --kind sveltekit --autocommit --trailer 'Co-Authored-By: someone <s@s>'`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			p.Name = args[0]
			if !cmd.Flags().Changed("path") {
				p.Path = p.Name
			}
			p.Tags = cleanTags(p.Tags)
			return a.editSubproject(manifest.ProjectOp{Action: manifest.ProjectAdd}, p, autocommit, trailers)
		},
	}
	f := c.Flags()
	f.StringVar(&p.Path, "path", "", "the sub-project's folder, relative to the repository root (default the name)")
	f.StringVar(&p.Kind, "kind", "", "what the sub-project is, such as go, node, python, rust, sveltekit, svelte, or template")
	f.StringSliceVar(&p.Tags, "tags", nil, "the sub-project's tags, separated by commas, such as ui,app")
	_ = c.MarkFlagRequired("kind")
	subprojectCommitFlags(c, &autocommit, &trailers)
	return c
}

func newSubprojectSetCmd(a *app) *cobra.Command {
	var to manifest.Project
	var autocommit bool
	var trailers []string
	c := &cobra.Command{
		Use:   "set <name>",
		Short: "Change a sub-project's name, path, kind, or tags, keeping the rest of its entry",
		Long: `Change the entry of projects in system-flow.yaml that the name given names:
each flag given replaces its key, and every other key of the entry, its
components and install among them, is kept. --name renames the entry;
--tags "" removes its tags. At least one flag is needed.

When the entry's folder is missing after the change, set creates it with a
README.md, as flai subproject add does, so that a sub-project's folder is
always there. A folder the entry named before is left as it is: set moves no
files. The change is committed, system-flow.yaml and that README.md alone,
only with --autocommit and unless the project sets dashboard.autocommit:
false. A name the manifest does not have is refused, naming those it has.

` + subprojectRefusals + `

` + subprojectJSON,
		Example: `  flai subproject set web --tags ui,app
  flai subproject set web --path apps/web
  flai subproject set web --name site --json`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := strings.TrimSpace(args[0])
			f := cmd.Flags()
			if !f.Changed("name") && !f.Changed("path") && !f.Changed("kind") && !f.Changed("tags") {
				return fmt.Errorf("nothing to change for sub-project %q: give --name, --path, --kind, or --tags, such as flai subproject set %s --tags ui", name, name)
			}
			repo, err := a.project()
			if err != nil {
				return err
			}
			// Start from the entry as it is; one the manifest lacks is
			// refused by the edit, naming those it has.
			p := manifest.Project{Name: name}
			if i := slices.IndexFunc(repo.Manifest.Projects, func(q manifest.Project) bool { return strings.TrimSpace(q.Name) == name }); i >= 0 {
				p = repo.Manifest.Projects[i]
			}
			// nil keeps the components and install keys as the manifest
			// writes them.
			p.Components, p.Install = nil, nil
			if f.Changed("name") {
				p.Name = to.Name
			}
			if f.Changed("path") {
				p.Path = to.Path
			}
			if f.Changed("kind") {
				p.Kind = to.Kind
			}
			if f.Changed("tags") {
				p.Tags = cleanTags(to.Tags)
			}
			return a.editSubprojectIn(repo, manifest.ProjectOp{Action: manifest.ProjectSet, Name: name}, p, autocommit, trailers)
		},
	}
	f := c.Flags()
	f.StringVar(&to.Name, "name", "", "the sub-project's new name")
	f.StringVar(&to.Path, "path", "", "the sub-project's folder, relative to the repository root; created when missing")
	f.StringVar(&to.Kind, "kind", "", "what the sub-project is, such as go, node, python, rust, sveltekit, svelte, or template")
	f.StringSliceVar(&to.Tags, "tags", nil, `the sub-project's tags, separated by commas, replacing those it has; "" removes them`)
	subprojectCommitFlags(c, &autocommit, &trailers)
	return c
}

func newSubprojectRemoveCmd(a *app) *cobra.Command {
	var autocommit bool
	var trailers []string
	c := &cobra.Command{
		Use:   "remove <name>",
		Short: "Remove a sub-project from projects, leaving its folder alone",
		Long: `Remove the entry of projects in system-flow.yaml that the name given names,
with its components and install keys. The sub-project's folder and what is
in it are left as they are; delete them yourself if you mean to. Removing
the last entry leaves projects: []. A name the manifest does not have is
refused, naming those it has, and nothing is written. The change is
committed, system-flow.yaml alone, only with --autocommit and unless the
project sets dashboard.autocommit: false.

` + subprojectJSON,
		Example: `  flai subproject remove web
  flai subproject remove web --autocommit --json`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return a.editSubproject(manifest.ProjectOp{Action: manifest.ProjectRemove}, manifest.Project{Name: args[0]}, autocommit, trailers)
		},
	}
	subprojectCommitFlags(c, &autocommit, &trailers)
	return c
}

// cleanTags are the tags given, each trimmed, without the empty ones.
func cleanTags(tags []string) []string {
	var out []string
	for _, t := range tags {
		if t = strings.TrimSpace(t); t != "" {
			out = append(out, t)
		}
	}
	return out
}

// editSubproject makes the edit of projects op says in the project's
// manifest, as editSubprojectIn does.
func (a *app) editSubproject(op manifest.ProjectOp, p manifest.Project, autocommit bool, trailers []string) error {
	repo, err := a.project()
	if err != nil {
		return err
	}
	return a.editSubprojectIn(repo, op, p, autocommit, trailers)
}

// editSubprojectIn makes the edit of projects op says in repo's manifest,
// with p the entry to add or set. An add or set first checks the edit
// without the folders, so that nothing is created for an edit refused for
// its name or path; then creates the entry's folder with a README.md when it
// is missing; then checks the edit against the repository's folders, and
// removes what it created when that refuses. It commits what it wrote when
// autocommit asks and the project allows, and says what it did.
func (a *app) editSubprojectIn(repo *workitem.Repo, op manifest.ProjectOp, p manifest.Project, autocommit bool, trailers []string) error {
	file := filepath.Join(repo.Root, manifest.File)
	data, err := os.ReadFile(file)
	if err != nil {
		return fmt.Errorf("cannot read %s to %s sub-project %q: %w; run flai from the repository, or restore the file", manifest.File, op.Action, strings.TrimSpace(p.Name), err)
	}
	change := subprojectChange{Action: op.Action, Project: p, Written: []string{manifest.File}}
	name := strings.TrimSpace(op.Name)
	if name == "" {
		name = strings.TrimSpace(p.Name)
	}
	before := slices.IndexFunc(repo.Manifest.Projects, func(q manifest.Project) bool { return strings.TrimSpace(q.Name) == name })
	var out []byte
	var top string
	if op.Action == manifest.ProjectRemove {
		// A remove creates no folder and cannot leave one missing, so it is
		// checked without them: another entry's missing folder does not
		// stop it.
		if out, err = manifest.EditProjects(data, nil, op, p); err != nil {
			return a.subprojectRefused(err)
		}
		if before >= 0 {
			change.Project = repo.Manifest.Projects[before]
		}
	} else {
		if _, err = manifest.EditProjects(data, nil, op, p); err != nil {
			return a.subprojectRefused(err)
		}
		var readme string
		if top, readme, err = makeSubprojectFolder(repo.Root, p, op.Action); err != nil {
			return err
		}
		if out, err = manifest.EditProjects(data, os.DirFS(repo.Root), op, p); err != nil {
			removeCreated(top)
			return a.subprojectRefused(err)
		}
		if readme != "" {
			change.Created = strings.TrimSpace(p.Path)
			change.Written = append(change.Written, readme)
		}
		if op.Action == manifest.ProjectSet && before >= 0 {
			was := repo.Manifest.Projects[before]
			change.Project.Components, change.Project.Install = was.Components, was.Install
			if strings.TrimSpace(was.Name) != strings.TrimSpace(p.Name) {
				change.From = was.Name
			}
		}
	}
	change.Project.Name, change.Project.Path, change.Project.Kind = strings.TrimSpace(change.Project.Name), strings.TrimSpace(change.Project.Path), strings.TrimSpace(change.Project.Kind)
	if err = os.WriteFile(file, out, 0o644); err != nil {
		removeCreated(top)
		return fmt.Errorf("cannot write %s to %s sub-project %q: %w; check that the file is writable, then run this again", manifest.File, op.Action, name, err)
	}
	change.Commit = a.commitSubproject(repo, autocommit, fmt.Sprintf("chore: %s sub-project %s", op.Action, change.Project.Name), trailers, change.Written)
	return a.printSubprojectChange(change)
}

// subprojectRefused prints a refusal of the edit as refusedSettings does, or
// returns any other error as it is.
func (a *app) subprojectRefused(err error) error {
	var refused *manifest.RefusedError
	if errors.As(err, &refused) {
		return a.refusedSettings(refused)
	}
	return err
}

// makeSubprojectFolder creates p's folder under root with a README.md in it,
// when the folder is missing, and says the topmost folder it created, to
// remove should the edit be refused, and the README.md's path relative to
// root; both empty when it created nothing. A path that is there, a file
// among it, is left to the edit's check, which refuses what is no folder.
func makeSubprojectFolder(root string, p manifest.Project, action manifest.ProjectAction) (top, readme string, err error) {
	rel := path.Clean(strings.TrimSpace(p.Path))
	dir := filepath.Join(root, filepath.FromSlash(rel))
	if _, err := os.Stat(dir); !errors.Is(err, fs.ErrNotExist) {
		return "", "", nil
	}
	parts := strings.Split(rel, "/")
	for i := range parts {
		at := filepath.Join(root, filepath.FromSlash(path.Join(parts[:i+1]...)))
		if _, err := os.Lstat(at); errors.Is(err, fs.ErrNotExist) {
			top = at
			break
		}
	}
	if top == "" {
		return "", "", fmt.Errorf("cannot create sub-project %q's folder %s: a link there names nothing; remove the link, or give --path another folder", strings.TrimSpace(p.Name), rel)
	}
	readme = path.Join(rel, subprojectReadme)
	body := fmt.Sprintf("# %s\n\n`flai subproject %s` created this folder for the sub-project %s, of kind `%s`.\n", strings.TrimSpace(p.Name), action, strings.TrimSpace(p.Name), strings.TrimSpace(p.Kind))
	if err := os.MkdirAll(dir, 0o755); err != nil {
		removeCreated(top)
		return "", "", fmt.Errorf("cannot create sub-project %q's folder %s: %w; check that its parent is writable, or give --path another folder", strings.TrimSpace(p.Name), rel, err)
	}
	if err := os.WriteFile(filepath.Join(dir, subprojectReadme), []byte(body), 0o644); err != nil {
		removeCreated(top)
		return "", "", fmt.Errorf("cannot write %s: %w; check that the folder is writable, then run this again", readme, err)
	}
	return top, readme, nil
}

// removeCreated removes the folder top and everything in it, which flai
// subproject created; an empty top removes nothing.
func removeCreated(top string) {
	if top != "" {
		_ = os.RemoveAll(top)
	}
}

// commitSubproject commits the paths written, and nothing else, when asked
// and the project commits what flai writes for the dashboard; it says the
// commit, or nothing, warning on standard error when the commit fails.
func (a *app) commitSubproject(repo *workitem.Repo, autocommit bool, msg string, trailers, paths []string) string {
	if !autocommit || !repo.Manifest.Autocommit() {
		return ""
	}
	if len(trailers) > 0 {
		msg += "\n\n" + strings.Join(trailers, "\n")
	}
	fail := func(err error) string {
		a.logger().Warn("sub-project change written but not committed", "component", "subproject", "paths", strings.Join(paths, ","), "error", err, "fix", "commit the paths by hand")
		return ""
	}
	if _, err := a.runner.Run(repo.Root, "git", append([]string{"add", "--"}, paths...)...); err != nil {
		return fail(err)
	}
	if _, err := a.runner.Run(repo.Root, "git", append([]string{"commit", "-q", "-m", msg, "--"}, paths...)...); err != nil {
		return fail(err)
	}
	sha, err := a.runner.Run(repo.Root, "git", "rev-parse", "--short", "HEAD")
	if err != nil {
		return fail(err)
	}
	return strings.TrimSpace(sha)
}

// printSubprojectChange says what flai subproject add, set, or remove did:
// the entry, the paths written, and the commit when there is one.
func (a *app) printSubprojectChange(c subprojectChange) error {
	if a.jsonOut {
		return a.printJSON(c)
	}
	p := c.Project
	switch {
	case c.Action == manifest.ProjectRemove:
		fmt.Fprintf(a.out, "removed sub-project %s; its folder %s is left as it is\n", p.Name, p.Path)
	case c.From != "":
		fmt.Fprintf(a.out, "renamed sub-project %s to %s (%s, %s)\n", c.From, p.Name, p.Path, p.Kind)
	default:
		verb := "added"
		if c.Action == manifest.ProjectSet {
			verb = "set"
		}
		fmt.Fprintf(a.out, "%s sub-project %s (%s, %s)\n", verb, p.Name, p.Path, p.Kind)
	}
	if c.Created != "" {
		fmt.Fprintf(a.out, "created folder %s\n", c.Created)
	}
	for _, w := range c.Written {
		fmt.Fprintf(a.out, "wrote %s\n", w)
	}
	if c.Commit != "" {
		fmt.Fprintf(a.out, "committed %s\n", c.Commit)
	}
	return nil
}
