package mcpserver

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/bytepunx/system-flow/flai/internal/guard"
	"github.com/bytepunx/system-flow/flai/internal/issues"
	"github.com/bytepunx/system-flow/flai/internal/manifest"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// The sub-projects (S-0491): projects in the manifest lists the repository's
// sub-projects. subprojects lists them for every agent; subproject_edit adds,
// changes, and removes one by running flai subproject, so that the folder a
// new sub-project needs is created by one code path, the command's, as the
// dashboard's settings.subproject creates it. Unlike the shared paths, the
// sub-projects are the project's to shape: a story's agent changes them on
// its branch, in its worktree, and flai guard refuses the edit only to
// sub-agents, the planner, the analyzer, and the orchestrator without its
// permission.

// subprojectWhere says which checkout both tools work in, as each tool says
// it.
const subprojectWhere = "The checkout is the calling session's story's worktree when the session has a story (FLAI_STORY, else the story in this agent's name of the form agent-S-nnnn, else the story branch checked out where this server runs) and the story has a worktree, so that the story's own branch carries the change and flai task done commits it; otherwise it is the project's main checkout, as for the orchestrator or the operator's own session; story names the story whose worktree it was, and is absent for the main checkout."

const subprojectsDescription = "The sub-projects (S-0491): each entry of projects in system-flow.yaml, in its order, as flai subproject list --json prints it: name, path (the folder, relative to the repository root), kind (such as go, node, python, rust, sveltekit, svelte, or template), tags ([] when there are none), folder (true when path names a folder in the checkout), and components and install as the manifest has them, or null when the entry leaves them out. " + subprojectWhere + " A read, open to every agent; flai subproject list does the same."

const subprojectEditDescription = "Add, change, or remove a sub-project, an entry of projects in system-flow.yaml (S-0491), by running flai subproject <op> --json in the checkout, so that the result is the command's. op add appends an entry named name, with kind (required) and the optional path (relative to the repository root; the name when left out) and tags; when the folder is missing it is created with a README.md naming the sub-project and its kind, so that git carries it. op set changes the entry named name: new_name renames it, path, kind, and tags replace theirs, and every other key, components and install among them, is kept; at least one of the four is needed, and tags given as an empty list removes the tags while tags left out keeps them; a folder missing after the change is created as for add, and set moves no files. op remove takes the entry out and leaves its folder alone; it takes name alone. Only the projects list is rewritten, the file's other keys and comments kept. The manifest is checked as it would be after the change, as flai check checks it: a name that is empty or another sub-project's, a path that is empty, the repository root, absolute, not clean, outside the repository, another sub-project's, or a file, and an empty kind are refused with each field and reason, and nothing is written, neither the manifest nor a folder; a name the manifest does not have is refused naming those it has. Returns action, project (the entry as written, or as it was before a remove, with components and install), from (the old name after a rename), created (the folder made, when one was), written (the files written, from the project root: system-flow.yaml and the README.md of a folder created), and story. " + subprojectWhere + " Nothing is committed: commit the files with the task they belong to, as task_done does. The sub-projects are the project's to shape, so a story's agent and the operator's own session call it; flai guard refuses it to every sub-agent, which puts the change in its final message instead, to the planner and the analyzer, and to the orchestrator. flai subproject add, set, and remove do the same."

// subprojectOps are the ops subproject_edit takes, as flai subproject's
// subcommands.
var subprojectOps = []string{string(manifest.ProjectAdd), string(manifest.ProjectSet), string(manifest.ProjectRemove)}

// flaiExecutable is the flai that subproject_edit runs: the one running,
// whose code this server is. A test gives the flai it built from this tree.
var flaiExecutable = os.Executable

// SubprojectsIn asks for the sub-projects.
type SubprojectsIn struct {
	Project string `json:"project,omitempty" jsonschema:"the project, by key or folder: needed only when the server serves more than one"`
}

func (in SubprojectsIn) project() string { return in.Project }

// Subproject is an entry of projects, as flai subproject list --json prints
// it.
type Subproject struct {
	Name       string               `json:"name"`
	Path       string               `json:"path" jsonschema:"the sub-project's folder, relative to the repository root"`
	Kind       string               `json:"kind" jsonschema:"what the sub-project is, such as go, node, python, rust, sveltekit, svelte, or template"`
	Tags       []string             `json:"tags" jsonschema:"the sub-project's tags; [] when there are none"`
	Folder     bool                 `json:"folder" jsonschema:"whether path names a folder in the checkout"`
	Components *manifest.Components `json:"components" jsonschema:"where a Svelte sub-project's components are found, as the manifest has it; null when the entry leaves it out"`
	Install    *manifest.Install    `json:"install" jsonschema:"how the sub-project installs its dependencies, as the manifest has it; null when the entry leaves it out"`
}

// SubprojectsOut is the sub-projects of one checkout.
type SubprojectsOut struct {
	Projects []Subproject `json:"projects" jsonschema:"projects in system-flow.yaml, in its order"`
	Story    string       `json:"story,omitempty" jsonschema:"the story whose worktree's manifest was read; absent for the main checkout"`
}

func (s *server) subprojects(_ context.Context, _ *mcp.CallToolRequest, _ SubprojectsIn) (*mcp.CallToolResult, SubprojectsOut, error) {
	r, story := s.subprojectCheckout()
	// the manifest as it is now: the one the server opened with may be older
	// than an edit made since, here or with flai subproject
	m, err := manifest.Load(filepath.Join(r.Root, manifest.File))
	if err != nil {
		return nil, SubprojectsOut{}, fmt.Errorf("cannot read the sub-projects in %s: %w; fix the manifest, then call subprojects again", s.rel(r.Root), err)
	}
	out := SubprojectsOut{Projects: make([]Subproject, len(m.Projects)), Story: story}
	for i, p := range m.Projects {
		tags := p.Tags
		if tags == nil {
			tags = []string{}
		}
		out.Projects[i] = Subproject{Name: p.Name, Path: p.Path, Kind: p.Kind, Tags: tags, Folder: isFolder(r.Root, p.Path), Components: p.Components, Install: p.Install}
	}
	return nil, out, nil
}

// isFolder says whether path, relative to root, names a folder below it.
func isFolder(root, path string) bool {
	if !filepath.IsLocal(path) {
		return false
	}
	st, err := os.Stat(filepath.Join(root, filepath.FromSlash(path)))
	return err == nil && st.IsDir()
}

// subprojectCheckout is the checkout the sub-project tools work in, and the
// story whose worktree it is: the calling session's story's worktree when it
// has one, else the project's main checkout and no story.
func (s *server) subprojectCheckout() (*workitem.Repo, string) {
	story := s.recordingStory("")
	r := issues.RepoFor(s.repo, story)
	if r == s.repo {
		return s.repo, ""
	}
	return r, workitem.CanonicalID(story)
}

// SubprojectEditIn adds, changes, or removes a sub-project, as the host
// API's settings.subproject takes it.
type SubprojectEditIn struct {
	Project string   `json:"project,omitempty" jsonschema:"the project, by key or folder: needed only when the server serves more than one"`
	Op      string   `json:"op" jsonschema:"add, set, or remove"`
	Name    string   `json:"name" jsonschema:"the sub-project's name: the new entry's for add, the entry to change or remove for set and remove"`
	NewName string   `json:"new_name,omitempty" jsonschema:"set only: the entry's new name"`
	Path    string   `json:"path,omitempty" jsonschema:"add and set: the sub-project's folder, relative to the repository root, created with a README.md when missing; add's default is the name"`
	Kind    string   `json:"kind,omitempty" jsonschema:"add (required) and set: what the sub-project is, such as go, node, python, rust, sveltekit, svelte, or template"`
	Tags    []string `json:"tags,omitempty" jsonschema:"add and set: the sub-project's tags, replacing those it has; on set an empty list removes them and leaving it out keeps them; a tag holds no comma"`
}

func (in SubprojectEditIn) project() string { return in.Project }

// SubprojectEditOut is what flai subproject add, set, or remove did, as
// --json prints it, with the files written from the project root.
type SubprojectEditOut struct {
	Action  string           `json:"action" jsonschema:"add, set, or remove"`
	Project manifest.Project `json:"project" jsonschema:"the entry as written, or as it was before a remove"`
	From    string           `json:"from,omitempty" jsonschema:"the entry's name before a set that renamed it"`
	Created string           `json:"created,omitempty" jsonschema:"the folder created, relative to the repository root, when one was"`
	Written []string         `json:"written" jsonschema:"the files written, from the project root: system-flow.yaml, and the README.md of a folder created"`
	Story   string           `json:"story,omitempty" jsonschema:"the story in whose worktree the change was written; absent for the main checkout"`
}

func (s *server) subprojectEdit(ctx context.Context, _ *mcp.CallToolRequest, in SubprojectEditIn) (*mcp.CallToolResult, SubprojectEditOut, error) {
	args, err := subprojectArgs(in)
	if err != nil {
		return nil, SubprojectEditOut{}, err
	}
	exe, err := flaiExecutable()
	if err != nil {
		return nil, SubprojectEditOut{}, fmt.Errorf("cannot find the flai to run flai subproject %s with: %w; run it on the host instead", in.Op, err)
	}
	r, story := s.subprojectCheckout()
	cmd := exec.CommandContext(ctx, exe, args...)
	cmd.Dir = r.Root
	cmd.Env = append(os.Environ(), "FLAI_AGENT="+s.agent, "FLAI_SESSION="+s.session, "LOG_FORMAT=json")
	var stdout, stderr bytes.Buffer
	cmd.Stdout, cmd.Stderr = &stdout, &stderr
	runErr := cmd.Run()
	var exit *exec.ExitError
	if runErr != nil && !errors.As(runErr, &exit) {
		return nil, SubprojectEditOut{}, fmt.Errorf("cannot run %s subproject %s: %w; nothing was written", exe, in.Op, runErr)
	}
	if runErr != nil {
		return nil, SubprojectEditOut{}, subprojectFailed(in, s.rel(r.Root), stdout.Bytes(), stderr.String())
	}
	var out SubprojectEditOut
	if err := json.Unmarshal(stdout.Bytes(), &out); err != nil || out.Action == "" {
		return nil, SubprojectEditOut{}, fmt.Errorf("flai subproject %s printed no result in %s%s; run flai subproject list to see what it changed", in.Op, s.rel(r.Root), said("printed", stdout.String()))
	}
	for i, w := range out.Written {
		out.Written[i] = s.rel(filepath.Join(r.Root, filepath.FromSlash(w)))
	}
	out.Story = story
	return nil, out, nil
}

// subprojectArgs are the flai arguments that make the edit in asks for, or
// why in is not one: every value is given as --flag=value, and the name
// after --, so that none is read as a flag.
func subprojectArgs(in SubprojectEditIn) ([]string, error) {
	op := strings.TrimSpace(in.Op)
	name := strings.TrimSpace(in.Name)
	switch {
	case !slices.Contains(subprojectOps, op):
		return nil, fmt.Errorf("op %q is not one subproject_edit takes: give add, set, or remove", in.Op)
	case name == "":
		return nil, fmt.Errorf("give name, the sub-project to %s, as subprojects lists it, or the new one's for add", op)
	}
	for _, t := range in.Tags {
		if strings.ContainsAny(t, ",\"") {
			return nil, fmt.Errorf("tag %q holds a comma or a double quote, which a tag does not: give each tag as its own entry of tags, such as [\"ui\", \"app\"]; nothing was written", t)
		}
	}
	args := []string{"subproject", op, "--json"}
	switch manifest.ProjectAction(op) {
	case manifest.ProjectAdd:
		if in.NewName != "" {
			return nil, fmt.Errorf("add takes no new_name: name is the new sub-project's name; nothing was written")
		}
		args = append(args, "--kind="+in.Kind)
		if in.Path != "" {
			args = append(args, "--path="+in.Path)
		}
		if len(in.Tags) > 0 {
			args = append(args, "--tags="+strings.Join(in.Tags, ","))
		}
	case manifest.ProjectSet:
		if in.NewName == "" && in.Path == "" && in.Kind == "" && in.Tags == nil {
			return nil, fmt.Errorf("nothing to change for sub-project %q: give new_name, path, kind, or tags (an empty list removes them); nothing was written", name)
		}
		if in.NewName != "" {
			args = append(args, "--name="+in.NewName)
		}
		if in.Path != "" {
			args = append(args, "--path="+in.Path)
		}
		if in.Kind != "" {
			args = append(args, "--kind="+in.Kind)
		}
		if in.Tags != nil {
			args = append(args, "--tags="+strings.Join(in.Tags, ","))
		}
	case manifest.ProjectRemove:
		if in.NewName != "" || in.Path != "" || in.Kind != "" || in.Tags != nil {
			return nil, fmt.Errorf("remove takes name alone, not new_name, path, kind, or tags: it takes the entry out and leaves its folder alone; nothing was written")
		}
	}
	return append(args, "--", name), nil
}

// subprojectFailed is the error for a flai subproject that exited non-zero:
// its refusal's fields and reasons, which it prints on standard output, else
// the error it logged on standard error, else what it said.
func subprojectFailed(in SubprojectEditIn, where string, stdout []byte, stderr string) error {
	var refused struct {
		Refused []manifest.Problem `json:"refused"`
	}
	if json.Unmarshal(stdout, &refused) == nil && len(refused.Refused) > 0 {
		parts := make([]string, len(refused.Refused))
		for i, p := range refused.Refused {
			parts[i] = p.String()
		}
		return fmt.Errorf("flai subproject %s refused %q in %s: %s; nothing was written", in.Op, strings.TrimSpace(in.Name), where, strings.Join(parts, "; "))
	}
	sc := bufio.NewScanner(strings.NewReader(stderr))
	sc.Buffer(make([]byte, 64*1024), 1024*1024)
	failed := ""
	for sc.Scan() {
		var ev struct {
			Err string `json:"err"`
		}
		if json.Unmarshal(sc.Bytes(), &ev) == nil && ev.Err != "" {
			failed = ev.Err
		}
	}
	if failed != "" {
		return fmt.Errorf("flai subproject %s %q in %s: %s", in.Op, strings.TrimSpace(in.Name), where, failed)
	}
	return fmt.Errorf("flai subproject %s %q failed in %s%s%s", in.Op, strings.TrimSpace(in.Name), where, said("printed", string(stdout)), said("said", stderr))
}

// addSubprojectTools registers the tools that list and change the
// sub-projects.
func addSubprojectTools(srv *mcp.Server, p projects) {
	mcp.AddTool(srv, &mcp.Tool{Name: guard.Subprojects, Description: subprojectsDescription}, route(p, (*server).subprojects))
	mcp.AddTool(srv, &mcp.Tool{Name: guard.SubprojectEdit, Description: subprojectEditDescription}, route(p, (*server).subprojectEdit))
}
