package mcpserver

import (
	"context"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/bytepunx/system-flow/flai/internal/adr"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

const adrNewDescription = "Record an architecture decision as flai adr new does (S-0275): decision is the decision as one sentence, its title. The number is one more than the highest in design/adrs in this checkout, on main, in every story's worktree, and on every story branch, so parallel stories do not take the same number; the file is NNNN-slug.md with its front matter (id, title, status, date, supersedes, superseded_by, and refines when given), its row is added to design/adrs/README.md, and each ADR it supersedes gets superseded_by set, the one edit allowed to an accepted ADR. status is proposed (the default) or accepted; supersedes and refines name ADRs as 7, 0007, or ADR-0007. body is the text below the heading, the project's 0000-template.md sections when empty. flai check runs with everything in place and refuses the ADR, putting every file back, if it reports anything the ADR introduces. story is the story the decision belongs to; by default FLAI_STORY, else the story in this agent's name of the form agent-S-nnnn, else the story branch checked out where this server runs, else none. The ADR is written in that story's worktree when it has one, for that story to commit, and in the project otherwise, so that it never lands in the main checkout while a story's worktree is its place. Returns the ADR's id, title, status, path (from the project root), changed (every file written, from the project root, the new ADR first), and the story whose worktree it was written in. With commit the subject is <ADR-nnnn> <decision>." + commitWhat

// AdrNewIn records an architecture decision (S-0275).
type AdrNewIn struct {
	Project    string   `json:"project,omitempty" jsonschema:"the project, by key or folder: needed only when the server serves more than one"`
	Decision   string   `json:"decision" jsonschema:"the decision, as one sentence: the ADR's title"`
	Status     string   `json:"status,omitempty" jsonschema:"proposed (the default) or accepted"`
	Supersedes []string `json:"supersedes,omitempty" jsonschema:"the ADRs this one supersedes, each as 7, 0007, or ADR-0007; each gets superseded_by set"`
	Refines    []string `json:"refines,omitempty" jsonschema:"the ADRs this one refines, each as 7, 0007, or ADR-0007"`
	Body       string   `json:"body,omitempty" jsonschema:"the markdown below the heading: Context, Decision, Consequences, Alternatives considered; the project's 0000-template.md sections when empty"`
	Story      string   `json:"story,omitempty" jsonschema:"the story the decision belongs to, in whose worktree it is written when it has one (default: FLAI_STORY, else the story in the agent's name agent-S-nnnn, else the story branch checked out)"`
	CommitIn
}

func (in AdrNewIn) project() string { return in.Project }

// AdrNewOut is the ADR adr_new recorded.
type AdrNewOut struct {
	ID      string   `json:"id"`
	Title   string   `json:"title"`
	Status  string   `json:"status"`
	Path    string   `json:"path" jsonschema:"from the project root; under the story's worktree when the ADR was written there"`
	Changed []string `json:"changed" jsonschema:"every file written, from the project root, the new ADR first: the index and any ADR it supersedes follow"`
	Story   string   `json:"story,omitempty" jsonschema:"the story in whose worktree the ADR was written"`
	CommittedOut
}

// adrNumbers reads ADR numbers given as 27, 0027, or ADR-0027.
func adrNumbers(field string, in []string) ([]int, error) {
	var out []int
	for _, s := range in {
		n, err := strconv.Atoi(strings.TrimPrefix(strings.ToUpper(strings.TrimSpace(s)), "ADR-"))
		if err != nil || n <= 0 {
			return nil, fmt.Errorf("adr_new %s %q is not an ADR number; write 27, 0027, or ADR-0027; nothing was written", field, s)
		}
		out = append(out, n)
	}
	return out, nil
}

func (s *server) adrNew(_ context.Context, _ *mcp.CallToolRequest, in AdrNewIn) (*mcp.CallToolResult, AdrNewOut, error) {
	supersedes, err := adrNumbers("supersedes", in.Supersedes)
	if err != nil {
		return nil, AdrNewOut{}, err
	}
	refines, err := adrNumbers("refines", in.Refines)
	if err != nil {
		return nil, AdrNewOut{}, err
	}
	story := s.recordingStory(in.Story)
	r, err := s.storyCheckout(story, in.Commit, "adr_new")
	if err != nil {
		return nil, AdrNewOut{}, err
	}
	if r == s.repo {
		// written in the project, not in a story's worktree
		story = ""
	}
	res, err := adr.New(r, s.gitRunner(), adr.Options{Title: in.Decision, Status: in.Status, Supersedes: supersedes, Refines: refines, Body: in.Body, Now: s.now()})
	if err != nil {
		return nil, AdrNewOut{}, fmt.Errorf("recording the decision %q in %s: %w", in.Decision, s.rel(r.Root), refusal(err))
	}
	out := AdrNewOut{ID: res.ID, Title: res.Title, Status: res.Status, Path: s.rel(filepath.Join(r.Root, res.Path)), Changed: []string{}}
	for _, p := range res.Changed {
		out.Changed = append(out.Changed, s.rel(filepath.Join(r.Root, p)))
	}
	if story != "" {
		out.Story = workitem.CanonicalID(story)
	}
	if in.Commit {
		out.CommittedOut, err = s.commitWritten(r, story, res.Changed, res.ID+" "+res.Title, in.Trailers)
	}
	return nil, out, err
}
