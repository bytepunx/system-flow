package mcpserver

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/bytepunx/system-flow/flai/internal/execx"
	"github.com/bytepunx/system-flow/flai/internal/storygit"
	"github.com/bytepunx/system-flow/flai/internal/verify"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

const testDescription = "Run the project's test and lint tiers (S-0273): the tiers the manifest's tests declare, those the paths select, cheapest first, stopping at the first that fails. paths are files, folders, or Go package directories relative to the repository root (a folder stands for every file below it); with none, the paths are what the checkout changed against the main branch, its commits and what it has not committed; all runs every tier for the whole checkout, whatever paths says. It runs in story's worktree, as its branch has the manifest and the files, or, with no story, in the project's main checkout; a story with no worktree is refused. Answers passed, the paths the tiers were selected for, and each tier with its command, its state (passed, failed, or not-reached), its exit code, and its duration; the tier that failed carries its first failures as findings, each a name, a path, a line, and a message, never the whole log, at most max across the run (5 by default; omitted counts the rest). A failing tier is an answer, not an error; a manifest whose tests are not valid, or a path not in the checkout, is refused naming what to fix. Call it between tasks instead of running go test, vitest, golangci-lint, or gofmt by hand and reading their logs. flai test on the command line answers the same, and flai test --json prints this result."

// TestIn runs the project's test tiers for paths in a checkout (S-0273).
type TestIn struct {
	Project string   `json:"project,omitempty" jsonschema:"the project, by key or folder: needed only when the server serves more than one"`
	Story   string   `json:"story,omitempty" jsonschema:"the story whose worktree the tiers run in (a task's ID runs in the task's own worktree, when it has one under .flai-cache/worktrees); with none, the project's main checkout"`
	Paths   []string `json:"paths,omitempty" jsonschema:"files, folders, or Go package directories relative to the repository root; with none, what the checkout changed against the main branch, uncommitted included"`
	All     bool     `json:"all,omitempty" jsonschema:"run every tier for the whole checkout, whatever paths says"`
	Max     int      `json:"max,omitempty" jsonschema:"the most findings to answer across the run; 5 when left out"`
}

func (in TestIn) project() string { return in.Project }

// test answers what flai test --json prints for the checkout the call names:
// the story's worktree or the main checkout.
func (s *server) test(ctx context.Context, _ *mcp.CallToolRequest, in TestIn) (*mcp.CallToolResult, verify.Result, error) {
	if in.Max < 0 {
		return nil, verify.Result{}, fmt.Errorf("max %d is not a number of findings; give 1 or more, or leave it out for %d", in.Max, verify.DefaultMax)
	}
	root, err := s.testRoot(in.Story)
	if err != nil {
		return nil, verify.Result{}, err
	}
	tiers, err := verify.CheckoutTiers(root)
	if err != nil {
		return nil, verify.Result{}, err
	}
	git := s.runner
	if git == nil {
		git = execx.System{}
	}
	git = execx.Timed(ctx, git)
	var base string
	if len(in.Paths) == 0 && !in.All {
		base, err = storygit.MainBranch(git, projectRoot(s.repo))
		if err != nil {
			return nil, verify.Result{}, fmt.Errorf("find the main branch to count the changes against; name the paths to test instead: %w", err)
		}
	}
	args := make([]string, 0, len(in.Paths))
	for _, p := range in.Paths {
		args = append(args, inCheckout(root, p))
	}
	res, err := verify.Test(ctx, verify.Options{Root: root, Tiers: tiers, Args: args, All: in.All, Base: base, Git: git, RunOptions: verify.RunOptions{Max: in.Max}})
	if err != nil {
		return nil, verify.Result{}, err
	}
	return nil, res, nil
}

// testRoot is the checkout a test call runs in: the worktree of the item
// named, which must exist, or the main checkout.
func (s *server) testRoot(id string) (string, error) {
	if strings.TrimSpace(id) == "" {
		return projectRoot(s.repo), nil
	}
	it, err := s.repo.Get(id)
	if err != nil {
		if errors.Is(err, workitem.ErrNotFound) {
			return "", fmt.Errorf("%s is not a work item here; name a story with a worktree, or leave story out to run in the main checkout", id)
		}
		return "", err
	}
	wt := s.repo.WorktreePath(it.ID)
	if st, err := os.Stat(wt); err != nil || !st.IsDir() {
		return "", fmt.Errorf("%s has no worktree at %s; open it with flai stream open %s on the host, or leave story out to run in the main checkout", it.ID, s.rel(wt), it.ID)
	}
	return wt, nil
}

// inCheckout is a path as verify takes it, relative to the checkout's root:
// an absolute path inside root is made relative to it, and any other path
// is left as given, for verify to accept or refuse.
func inCheckout(root, p string) string {
	if !filepath.IsAbs(p) {
		return p
	}
	rel, err := filepath.Rel(root, p)
	if err != nil || rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		return p
	}
	return filepath.ToSlash(rel)
}
