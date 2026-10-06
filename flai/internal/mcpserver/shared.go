package mcpserver

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/bytepunx/system-flow/flai/internal/guard"
	"github.com/bytepunx/system-flow/flai/internal/manifest"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// The shared paths (S-0295, ADR-0096): claims.shared in the manifest lists
// glob patterns of the paths many stories change in separate sections or new
// files, and an overlap wholly inside one holds no story. shared_paths lists
// and checks them for every agent; shared_paths_edit changes them, which is
// the operator's alone, so flai guard refuses it to every session flai serve
// starts and the tool refuses itself in one.

// sharedDialect is the glob dialect of claims.shared, as each tool says it.
const sharedDialect = "Patterns are paths relative to the repository root, separated by /: * is any characters within one segment, ** zero or more whole segments, ? one character within a segment, and a plain path covers itself and everything below it."

const sharedPathsDescription = "The shared paths (ADR-0096): claims.shared in system-flow.yaml, glob patterns of the paths many stories change in separate sections or new files. An overlap of two stories' claims that lies wholly inside one holds no ready story. Returns shared, the patterns in list order, and invalid, each that is not valid and so frees nothing, with the reason. Given paths (paths or touches entries; a component's name or tag is read as its path) or story (a story's ID, whose claim is checked entry by entry as its hold reads it), entries says for each its entry, path (as matched), whether it is shared, the pattern it lies inside, and the story whose claim it is. An entry whose last segment has an extension, such as flai.md, is a file and lies inside a pattern that matches it; any other is read as a folder, which lies inside only a pattern that covers everything below it. " + sharedDialect + " A read, open to every agent; flai shared list and check do the same."

const sharedPathsEditDescription = "Change the shared paths, claims.shared in system-flow.yaml (ADR-0096): remove takes out patterns named as shared_paths lists them, then add appends patterns. Only that list is rewritten, the file's other keys and comments kept, and nothing is committed. A pattern that is not valid (empty, absolute, with a . or .. segment, or a malformed glob), one to add that is in the list already, or one to remove that is not in it is refused with the reason, and nothing is written. Returns added, removed, and shared, the list after the change. " + sharedDialect + " The list decides which overlaps hold a story, so it is the operator's alone: flai guard refuses this tool to every session flai serve starts (a story's agent and its sub-agents, the planner, the orchestrator, and the analyzer), and it refuses itself in one; ask the operator on a thread with thread_open instead. flai shared add and remove do the same."

// SharedPathsIn asks for the shared paths and, optionally, which of some
// entries lie inside them.
type SharedPathsIn struct {
	Project string   `json:"project,omitempty" jsonschema:"the project, by key or folder: needed only when the server serves more than one"`
	Paths   []string `json:"paths,omitempty" jsonschema:"paths or touches entries to check, relative to the repository root; a component's name or tag is read as its path"`
	Story   string   `json:"story,omitempty" jsonschema:"a story's ID (any zero padding): each entry of its claim, as its hold reads it, is checked"`
}

func (in SharedPathsIn) project() string { return in.Project }

// SharedEntry is one entry checked against the shared paths, as flai shared
// check --json prints it.
type SharedEntry struct {
	Entry   string `json:"entry" jsonschema:"the path or touches entry given, or an entry of the story's claim"`
	Path    string `json:"path" jsonschema:"the entry as it is matched: a component's name or tag read as its path, without a leading ./ or a trailing /"`
	Shared  bool   `json:"shared" jsonschema:"whether the entry lies wholly inside a pattern"`
	Pattern string `json:"pattern,omitempty" jsonschema:"the first pattern it lies inside, when it is shared"`
	Story   string `json:"story,omitempty" jsonschema:"the story whose claim the entry is, when a story was given"`
}

// SharedInvalid is a pattern of claims.shared that is not valid.
type SharedInvalid struct {
	Pattern string `json:"pattern"`
	Reason  string `json:"reason" jsonschema:"what is wrong and what to write instead"`
}

// SharedPathsOut is the shared paths and the entries checked against them.
type SharedPathsOut struct {
	Shared  []string        `json:"shared" jsonschema:"claims.shared's patterns in list order"`
	Invalid []SharedInvalid `json:"invalid,omitempty" jsonschema:"the patterns that are not valid, which free nothing"`
	Entries []SharedEntry   `json:"entries" jsonschema:"each path given, then each entry of the story's claim, in order; empty when neither is given"`
}

func (s *server) sharedPaths(_ context.Context, _ *mcp.CallToolRequest, in SharedPathsIn) (*mcp.CallToolResult, SharedPathsOut, error) {
	// the manifest as it is now: the one the server opened with may be older
	// than an edit made since, here or with flai shared
	m, err := manifest.Load(s.manifestFile())
	if err != nil {
		return nil, SharedPathsOut{}, fmt.Errorf("cannot read the shared paths: %w; fix the manifest, then call shared_paths again", err)
	}
	out := SharedPathsOut{Shared: m.Claims.Shared, Entries: []SharedEntry{}}
	if out.Shared == nil {
		out.Shared = []string{}
	}
	for _, e := range m.Claims.Errors() {
		out.Invalid = append(out.Invalid, SharedInvalid{Pattern: e.Pattern, Reason: e.Reason})
	}
	components := workitem.NewHolds(nil, m.Projects)
	for _, p := range in.Paths {
		entry := strings.TrimSpace(p)
		if entry == "" {
			return nil, SharedPathsOut{}, fmt.Errorf("an empty entry in paths names no path; give a path or a touches entry, such as docs/users/flai.md")
		}
		out.Entries = append(out.Entries, sharedEntry(m.Claims, entry, entryPath(components, entry), ""))
	}
	if in.Story != "" {
		claim, id, err := s.storyClaim(in.Story, m.Projects)
		if err != nil {
			return nil, SharedPathsOut{}, err
		}
		for _, e := range claim {
			out.Entries = append(out.Entries, sharedEntry(m.Claims, e, e, id))
		}
	}
	return nil, out, nil
}

// sharedEntry is whether a touches entry, read as path, lies inside a
// pattern of claims.
func sharedEntry(claims manifest.Claims, entry, path, story string) SharedEntry {
	p, ok := claims.Covers(path)
	return SharedEntry{Entry: entry, Path: path, Shared: ok, Pattern: p, Story: story}
}

// entryPath is a touches entry as a claim reads it: a component's name or
// tag becomes its path.
func entryPath(components *workitem.Holds, entry string) string {
	if c := components.Claim(&workitem.Item{Touches: []string{entry}}); len(c) == 1 {
		return c[0]
	}
	return entry
}

// storyClaim is the claim of the story id names, as its hold reads it with
// the manifest's components, and the story's ID as flai writes it.
func (s *server) storyClaim(id string, components []manifest.Project) ([]string, string, error) {
	if workitem.TypeOfID(workitem.CanonicalID(id)) != workitem.Story {
		return nil, "", fmt.Errorf("%s is not a story; give story a story's ID, or the paths to check in paths", workitem.CanonicalID(id))
	}
	story, err := s.repo.Get(id)
	if err != nil {
		return nil, "", err
	}
	items, err := s.repo.List(story.Archived)
	if err != nil {
		return nil, "", fmt.Errorf("cannot read the tasks of %s: %w; run flai check to see what is wrong", story.ID, err)
	}
	return workitem.NewHolds(items, components).Claim(story), story.ID, nil
}

// manifestFile is the project's system-flow.yaml, in the checkout the
// server opened, where flai shared reads and writes it too.
func (s *server) manifestFile() string { return filepath.Join(s.repo.Root, manifest.File) }

// SharedPathsEditIn adds patterns to the shared paths and removes them.
type SharedPathsEditIn struct {
	Project string   `json:"project,omitempty" jsonschema:"the project, by key or folder: needed only when the server serves more than one"`
	Add     []string `json:"add,omitempty" jsonschema:"patterns to append to claims.shared"`
	Remove  []string `json:"remove,omitempty" jsonschema:"patterns to take out of claims.shared, as shared_paths lists them; removed before add is applied"`
}

func (in SharedPathsEditIn) project() string { return in.Project }

func (s *server) sharedPathsEdit(_ context.Context, _ *mcp.CallToolRequest, in SharedPathsEditIn) (*mcp.CallToolResult, manifest.Change, error) {
	if s.served {
		return nil, manifest.Change{}, fmt.Errorf("a session flai serve started does not change the shared paths: they decide which overlaps hold a story, so the list is the operator's (ADR-0096); ask the operator on a thread with thread_open, naming the pattern and what it would free")
	}
	if len(in.Add) == 0 && len(in.Remove) == 0 {
		return nil, manifest.Change{}, fmt.Errorf("give add, remove, or both: the patterns to add to claims.shared or to remove from it")
	}
	file := s.manifestFile()
	data, err := os.ReadFile(file)
	if err != nil {
		return nil, manifest.Change{}, fmt.Errorf("cannot read the manifest to change the shared paths: %w", err)
	}
	var ch manifest.Change
	if len(in.Remove) > 0 {
		if data, ch, err = manifest.RemoveSharedBytes(data, in.Remove...); err != nil {
			return nil, manifest.Change{}, fmt.Errorf("%s: %w; nothing was written", manifest.File, err)
		}
	}
	if len(in.Add) > 0 {
		removed := ch.Removed
		if data, ch, err = manifest.AddSharedBytes(data, in.Add...); err != nil {
			return nil, manifest.Change{}, fmt.Errorf("%s: %w; nothing was written", manifest.File, err)
		}
		ch.Removed = removed
	}
	if err := os.WriteFile(file, data, 0o644); err != nil {
		return nil, manifest.Change{}, fmt.Errorf("cannot write the shared paths to %s: %w", file, err)
	}
	return nil, ch, nil
}

// servedSession says whether flai serve started the session this server
// runs in, whose environment it inherits: flai serve marks it with
// FLAI_STARTED_BY, and only flai serve sets a role or FLAI_STORY, as flai
// guard reads them.
func servedSession(role string) bool {
	return role != "" || os.Getenv("FLAI_STORY") != "" || os.Getenv(guard.StartedByEnv) == guard.StartedByServe
}

// addSharedTools registers the tools that list, check, and change the
// shared paths.
func addSharedTools(srv *mcp.Server, p projects) {
	mcp.AddTool(srv, &mcp.Tool{Name: "shared_paths", Description: sharedPathsDescription}, route(p, (*server).sharedPaths))
	mcp.AddTool(srv, &mcp.Tool{Name: guard.SharedPathsEdit, Description: sharedPathsEditDescription}, route(p, (*server).sharedPathsEdit))
}
